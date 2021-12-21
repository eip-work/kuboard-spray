package command

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write the file to the client.
	//writeWait = 1 * time.Second
	writeWait = 100 * time.Millisecond

	// Time allowed to read the next pong message from the client.
	//pongWait = 24 * time.Hour
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 1 * time.Second
)

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
func tailFile(filePath string) *tail.Tail {

	tailfs, err := tail.TailFile(filePath, tail.Config{
		ReOpen: true, // 文件被移除或被打包，需要重新打开
		Follow: true, // 实时跟踪
		// Location:  &tail.SeekInfo{Offset: 10, Whence: 2}, // 如果程序出现异常，保存上次读取的位置，避免重新读取。
		MustExist: false, // 如果文件不存在，是否退出程序，false是不退出
		Poll:      true,
	})

	if err != nil {
		fmt.Println("tailf failed, err:", err)
		return nil
	}
	return tailfs
}
func writer(ws *websocket.Conn, filePath string) {
	tailfs := tailFile(filePath)
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case msg, ok := <-tailfs.Lines:
			if ok {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				// fmt.Printf("read file content： %s\n", msg)
				if err := ws.WriteMessage(websocket.TextMessage, []byte(msg.Text)); err != nil {
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

type TailFileRequest struct {
	Cluster string `uri:"cluster" binding:"required"`
	Pid     string `uri:"pid" binding:"required"`
	File    string `uri:"file" binding:"required"`
}

func TailFile(c *gin.Context) {

	var reqParams TailFileRequest
	c.ShouldBindUri(&reqParams)

	pid := reqParams.Pid
	if reqParams.Pid == "lastrun" {
		lockFilePath := constants.GET_DATA_CLUSTER_DIR() + "/" + reqParams.Cluster + "/inventory.lastrun"
		logrus.Trace("read pid from : ", lockFilePath)
		b, err := ioutil.ReadFile(lockFilePath)
		if err != nil {
			common.HandleError(c, http.StatusInternalServerError, "Cannot read file "+lockFilePath, err)
			return
		}
		pid = string(b)
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.HandleError(c, http.StatusUpgradeRequired, "failed to upgrade", err)
		return
	}
	defer ws.Close()

	filePath := constants.GET_DATA_CLUSTER_DIR() + "/" + reqParams.Cluster + "/history/" + pid + "/" + reqParams.File
	logrus.Trace("[", filePath, "]")
	go writer(ws, filePath)
	reader(ws)
}
