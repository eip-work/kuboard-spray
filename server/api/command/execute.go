package command

import (
	"errors"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

type Execute struct {
	Cmd      string
	Args     func(string) []string
	Cluster  string
	mutex    *sync.Mutex
	Type     string
	PreExec  func(string) error
	PostExec func(string, string, *os.File) error
	Dir      string
	Env      []string
	R_Error  error
	R_Pid    string
}

func (execute *Execute) ToString(runDirPath string) string {
	result := "--" + execute.Cluster + "-- [" + execute.Cmd + ", "
	args := execute.Args(runDirPath)
	for index, arg := range args {
		result += arg
		if index < len(args)-1 {
			result += ", "
		}
	}
	result += "]"
	return result
}

func (execute *Execute) Exec() error {
	execute.mutex = &sync.Mutex{}
	execute.mutex.Lock()
	go execute.exec()
	logrus.Trace("geting lock in main thread...")
	execute.mutex.Lock()
	logrus.Trace("Got response from exec: ", execute.R_Error)
	execute.mutex.Unlock()
	return execute.R_Error
}

func (execute *Execute) exec() {

	lockFile, err := LockCluster(execute.Cluster)
	if err != nil {
		execute.R_Error = err
		execute.mutex.Unlock()
		return
	}

	defer UnlockCluster(lockFile)

	pid := time.Now().Format("2006-01-02_15-04-05.999") + "_" + execute.Type
	historyPath := constants.GET_DATA_INVENTORY_DIR() + "/" + execute.Cluster + "/history"
	if err := common.CreateDirIfNotExists(historyPath); err != nil {
		execute.R_Error = errors.New("cannot create historyDir : " + historyPath + " : " + err.Error())
		execute.mutex.Unlock()
		return
	}

	runDirPath := constants.GET_DATA_INVENTORY_DIR() + "/" + execute.Cluster + "/history/" + pid
	if err := common.CreateDirIfNotExists(runDirPath); err != nil {
		execute.R_Error = errors.New("cannot create runDir : " + runDirPath + " : " + err.Error())
		execute.mutex.Unlock()
		return
	}

	if execute.PreExec != nil {
		if err := execute.PreExec(runDirPath); err != nil {
			execute.R_Error = errors.New("failed to prepare for the task : " + err.Error())
			execute.mutex.Unlock()
			return
		}
	}

	logFilePath := runDirPath + "/command.log"
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		execute.R_Error = errors.New("cannot create logFile : " + logFilePath + " : " + err.Error())
		execute.mutex.Unlock()
		return
	}

	defer logFile.Sync()
	defer logFile.Close()

	cmd := exec.Command(execute.Cmd, execute.Args(runDirPath)...)

	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Dir = execute.Dir
	cmd.Env = execute.Env

	if err := cmd.Start(); err != nil {
		execute.R_Error = errors.New("failed to start command " + execute.ToString(runDirPath) + " : " + err.Error())
		os.Remove(runDirPath)
		execute.mutex.Unlock()
		return
	}

	logrus.Trace("started command "+execute.ToString(runDirPath), pid)

	if err := lockFile.Truncate(0); err != nil {
		execute.R_Error = errors.New("failed to truncate lockFile : " + err.Error())
	}
	_, err = lockFile.WriteString(pid)
	if err != nil {
		execute.R_Error = errors.New("failed to write pid " + err.Error())
	}

	execute.R_Pid = pid

	execute.mutex.Unlock()
	cmd.Wait()

	logFile.WriteString("\n\n\nKUBOARD SPRAY *****************************************************************\n")
	if execute.PostExec != nil {
		if err := execute.PostExec(pid, runDirPath, logFile); err != nil {
			logFile.WriteString("Error in PostExec: " + err.Error() + "\n")
		}
	}
}
