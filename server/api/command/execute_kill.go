package command

import (
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
)

type ExecuteKillRequest struct {
	OwnerType string `uri:"owner_type" binding:"required"`
	OwnerName string `uri:"owner_name" binding:"required"`
	Pid       string `uri:"pid"`
	Operation string `uri:"operation"`
	Step      string `uri:"step"`
	Time      string `uri:"time"`
}

func ExecuteKill(c *gin.Context) {
	var req ExecuteKillRequest
	c.ShouldBindUri(&req)

	pid := req.Pid
	if pid == "" {
		pid = req.Operation + "/" + req.Step + "/" + req.Time
	}

	process := runningProcesses[pid]
	if process == nil {
		common.HandleError(c, http.StatusNotFound, "process doesnot exist or is already killed. "+pid, nil)
		return
	}

	defer delete(runningProcesses, pid)

	if err := process.Kill(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to kill process "+pid, err)
		return
	}

	go func() {
		time.Sleep(time.Duration(1) * time.Second)

		logFilePath := constants.GET_DATA_DIR() + "/" + req.OwnerType + "/" + req.OwnerName + "/history/" + pid + "/execute.log"
		logrus.Trace("logFilePath: ", logFilePath)
		logFile, err := os.Create(logFilePath)
		if err != nil {
			logrus.Warn("cannot write logFile ", logFilePath, err)
			debug.PrintStack()
		} else {
			msg := "\n\n" +
				"\033[31m\033[01m\033[05m ###############" + "######################" + "###############\033[0m \n" +
				"\033[31m\033[01m\033[05m ###############" + " You killed the task. " + "###############\033[0m \n" +
				"\033[31m\033[01m\033[05m ###############" + " 您强制结束了该任务。 " + "###############\033[0m \n" +
				"\033[31m\033[01m\033[05m ###############" + "######################" + "###############\033[0m \n" +
				"\n\n"
			_, err = logFile.WriteString(msg)
			if err != nil {
				logrus.Warn(err.Error())
			}
		}
		defer logFile.Close()

		// 检查集群任务执行结果
		if req.OwnerType == "cluster" && req.Pid == "" {
			cssr := CheckStepStatusRequest{
				Cluster:   req.OwnerName,
				Operation: req.Operation,
				Step:      req.Step,
			}
			cssrResponse, err := CheckStepStatusExec(cssr)
			if err == nil {
				common.SaveYamlFile(constants.GET_DATA_DIR()+"/"+req.OwnerType+"/"+req.OwnerName+"/history/"+pid+"/result.yaml", cssrResponse)
			}
		}

		// 保存任务执行状态 status.yaml
		statusFilePath := constants.GET_DATA_DIR() + "/" + req.OwnerType + "/" + req.OwnerName + "/history/" + pid + "/status.yaml"
		statusTemp, err := common.ParseYamlFile(statusFilePath)
		if err == nil {
			endTime := time.Now()
			startTime, err := time.Parse(time.RFC3339Nano, statusTemp["startTime"].(string))
			if err == nil {
				statusTemp["status"] = "killed"
				statusTemp["endTime"] = endTime.Format(time.RFC3339Nano)
				statusTemp["duration"] = endTime.Sub(startTime).Milliseconds()
			}
			common.SaveYamlFile(statusFilePath, statusTemp)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
	})
}
