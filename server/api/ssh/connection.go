package ssh

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func dialSsh(node NodeInfo) (*ssh.ClientConfig, error) {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		config       ssh.Config
	)
	auth = make([]ssh.AuthMethod, 0)
	if node.Password != "" {
		auth = append(auth, ssh.Password(node.Password))
	}
	if node.PrivateKeyPath != "" {
		key, err := ioutil.ReadFile(node.PrivateKeyPath)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				auth = append(auth, ssh.PublicKeys(signer))
			} else {
				logrus.Warning("")
			}
		} else {
			logrus.Warning("cannot read private key: ", err)
		}
	}
	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	clientConfig = &ssh.ClientConfig{
		User:    node.User,
		Auth:    auth,
		Timeout: 5 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return clientConfig, nil
}

func (sshClient *SSHClient) GenerateClient() error {

	var client *ssh.Client
	var err error

	if sshClient.Bastion != nil {
		var bastionClient *ssh.Client
		bClientConfig, err := dialSsh(*sshClient.Bastion)
		if err != nil {
			return err
		}
		addr := fmt.Sprintf("%s:%d", sshClient.Bastion.Host, sshClient.Bastion.Port)
		if bastionClient, err = ssh.Dial("tcp", addr, bClientConfig); err != nil {
			return err
		}

		addr = fmt.Sprintf("%s:%d", sshClient.Host, sshClient.Port)
		conn, err := bastionClient.Dial("tcp", addr)

		if err != nil {
			return err
		}

		clientConfig, err := dialSsh(sshClient.NodeInfo)
		if err != nil {
			return err
		}

		ncc, chans, reqs, err := ssh.NewClientConn(conn, addr, clientConfig)
		if err != nil {
			return err
		}

		sClient := ssh.NewClient(ncc, chans, reqs)

		logrus.Trace("connected throw bastion.")

		sshClient.Client = sClient
	}

	clientConfig, err := dialSsh(sshClient.NodeInfo)
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d", sshClient.Host, sshClient.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err == nil {
		sshClient.Client = client
		return nil
	} else {
		return err
	}
}

func (sshClient *SSHClient) RequestTerminal(terminal TerminalSpec) *SSHClient {
	session, err := sshClient.Client.NewSession()
	if err != nil {
		logrus.Println(err)
		return nil
	}
	sshClient.Session = session
	channel, inRequests, err := sshClient.Client.OpenChannel("session", nil)
	if err != nil {
		logrus.Println(err)
		return nil
	}
	sshClient.channel = channel
	go func() {
		for req := range inRequests {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := ptyRequestMsg{
		Term:     "xterm",
		Columns:  terminal.Columns,
		Rows:     terminal.Rows,
		Width:    uint32(terminal.Columns * 8),
		Height:   uint32(terminal.Columns * 8),
		Modelist: string(modeList),
	}
	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		logrus.Println(err)
		return nil
	}
	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		logrus.Println(err)
		return nil
	}
	return sshClient
}

func (sshClient *SSHClient) Connect(ws *websocket.Conn) {

	//这里第一个协程获取用户的输入
	go func() {
		for {
			// p为用户输入
			_, p, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if string(p[0:1]) == "0" {
				_, err = sshClient.channel.Write(p[1:])
				if err != nil {
					return
				}
			} else if string(p[0:1]) == "1" {
				terminal := TerminalSpec{}
				json.Unmarshal(p[1:], &terminal)
				size := make([]byte, 16)
				binary.BigEndian.PutUint32(size, uint32(terminal.Columns))
				binary.BigEndian.PutUint32(size[4:], uint32(terminal.Rows))

				ok, err := sshClient.channel.SendRequest("window-change", false, size)
				if !ok || err != nil {
					logrus.Println(ok, err)
					// return
				}
			}
		}
	}()

	//第二个协程将远程主机的返回结果返回给用户
	go func() {
		br := bufio.NewReader(sshClient.channel)
		buf := []byte{}
		t := time.NewTimer(time.Microsecond * 100)
		defer t.Stop()
		// 构建一个信道, 一端将数据远程主机的数据写入, 一段读取数据写入ws
		r := make(chan rune)

		// 另起一个协程, 一个死循环不断的读取ssh channel的数据, 并传给r信道直到连接断开
		go func() {
			defer sshClient.Client.Close()
			defer sshClient.Session.Close()

			for {
				x, size, err := br.ReadRune()
				if err != nil {
					logrus.Println(err)
					ws.WriteMessage(1, []byte("\033[31m已经关闭连接!\033[0m"))
					ws.Close()
					return
				}
				if size > 0 {
					r <- x
				}
			}
		}()

		// 主循环
		for {
			select {
			// 每隔100微秒, 只要buf的长度不为0就将数据写入ws, 并重置时间和buf
			case <-t.C:
				if len(buf) != 0 {
					err := ws.WriteMessage(websocket.TextMessage, buf)
					buf = []byte{}
					if err != nil {
						logrus.Println(err)
						return
					}
				}
				t.Reset(time.Microsecond * 100)
			// 前面已经将ssh channel里读取的数据写入创建的通道r, 这里读取数据, 不断增加buf的长度, 在设定的 100 microsecond后由上面判定长度是否返送数据
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					buf = append(buf, []byte("@")...)
				}
			}
		}
	}()

	defer func() {
		if err := recover(); err != nil {
			logrus.Println(err)
		}
	}()
}
