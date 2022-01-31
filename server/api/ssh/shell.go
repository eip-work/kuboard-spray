package ssh

import (
	"net/http"
	"strconv"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ShellRequest struct {
	NodeOwnerType string `uri:"node_owner_type"`
	NodeOwner     string `uri:"node_owner"`
	Node          string `uri:"node"`
}

func ShellWs(c *gin.Context) {
	var err error
	cols := c.DefaultQuery("cols", "260")
	rows := c.DefaultQuery("rows", "80")
	col, _ := strconv.Atoi(cols)
	row, _ := strconv.Atoi(rows)
	terminal := TerminalSpec{
		Columns: uint32(col),
		Rows:    uint32(row),
	}

	var shellRequest ShellRequest
	c.ShouldBindUri(&shellRequest)
	sshClient, err := NewSSHClient(shellRequest)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "bad status", err)
		return
	}
	logrus.Trace(sshClient, "host[", sshClient.Host, "] password[", sshClient.Password, "] privateKeyPath[", sshClient.PrivateKeyPath, "]")
	if sshClient.Host == "" || sshClient.Password == "" && sshClient.PrivateKeyPath == "" {
		common.HandleError(c, http.StatusBadRequest, "IP 或 密码为空", nil)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}
	err = sshClient.GenerateClient()
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
		conn.Close()
		return
	}
	sshClient.RequestTerminal(terminal)
	sshClient.Connect(conn)
}
