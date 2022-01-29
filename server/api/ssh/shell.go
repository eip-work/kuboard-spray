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

func ShellWs(c *gin.Context) {
	var err error
	msg := c.DefaultQuery("msg", "")
	cols := c.DefaultQuery("cols", "150")
	rows := c.DefaultQuery("rows", "35")
	col, _ := strconv.Atoi(cols)
	row, _ := strconv.Atoi(rows)
	terminal := Terminal{
		Columns: uint32(col),
		Rows:    uint32(row),
	}
	sshClient, err := DecodedMsgToSSHClient(msg)
	if err != nil {
		c.Error(err)
		return
	}
	logrus.Trace(sshClient)
	if sshClient.IpAddress == "" || sshClient.Password == "" {
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
