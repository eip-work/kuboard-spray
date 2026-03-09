package command

import (
	"errors"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
)

var runningProcesses = map[string]*os.Process{}

type Execute struct {
	Cmd       string
	Args      func(string) []string
	OwnerType string
	OwnerName string
	mutex     *sync.Mutex
	Type      string
	PreExec   func(string) error
	PostExec  func(ExecuteExitStatus) (string, error)
	Dir       string
	Env       []string
	Pid       string
	R_Error   error
	R_Pid     string
}

func (execute *Execute) ToString(execute_dir_path string, pid string) string {
	result := "---"
	result += "\nowner_type: " + execute.OwnerType
	result += "\nowner_name: " + execute.OwnerName
	result += "\nhistory: " + execute_dir_path
	result += "\npid: " + pid
	result += "\ndir: " + execute.Dir
	result += "\ncmd: " + execute.Cmd
	result += "\nargs:"
	args := execute.Args(execute_dir_path)
	for _, arg := range args {
		result += "\n  - " + arg
	}
	result += "\nenv:"
	for _, env := range execute.Env {
		result += "\n - " + env
	}
	result += "\nshell: cd " + execute.Dir
	for _, env := range execute.Env {
		result += " && export \"" + env + "\""
	}
	result += " && " + execute.Cmd
	for _, arg := range args {
		result += " " + arg
	}
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

	lockFile, err := LockOwner(execute.OwnerType, execute.OwnerName)
	if err != nil {
		execute.R_Error = err
		execute.mutex.Unlock()
		return
	}

	defer UnlockOwner(lockFile)
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("%+v", r)
			debug.PrintStack()
		}
	}()

	startTime := time.Now().Format(time.RFC3339Nano)
	pid := startTime + "_" + execute.Type
	if execute.Pid != "" {
		pid = execute.Pid
	}
	historyPath := constants.GET_DATA_DIR() + "/" + execute.OwnerType + "/" + execute.OwnerName + "/history"
	if err := common.CreateDirIfNotExists(historyPath); err != nil {
		execute.R_Error = errors.New("cannot create historyDir : " + historyPath + " : " + err.Error())
		execute.mutex.Unlock()
		return
	}

	execute_dir_path := historyPath + "/" + pid
	if err := common.CreateDirIfNotExists(execute_dir_path); err != nil {
		execute.R_Error = errors.New("cannot create runDir : " + execute_dir_path + " : " + err.Error())
		execute.mutex.Unlock()
		return
	}

	status := map[string]string{
		"startTime": startTime,
		"status":    "processing",
	}
	common.SaveYamlFile(execute_dir_path+"/status.yaml", status)

	common.Symlink(execute.Dir+"/group_vars", execute_dir_path+"/group_vars")
	// execute.symlink(execute.Dir+"/cache", execute_dir_path+"/cache")

	if execute.PreExec != nil {
		if err := execute.PreExec(execute_dir_path); err != nil {
			execute.R_Error = errors.New("failed to prepare for the task : " + err.Error())
			execute.mutex.Unlock()
			return
		}
	}

	logFilePath := execute_dir_path + "/execute.log"
	logFile, err := os.Create(logFilePath)
	if err != nil {
		execute.R_Error = errors.New("cannot create logFile : " + logFilePath + " : " + err.Error())
		execute.mutex.Unlock()
		return
	}

	defer logFile.Sync()
	defer logFile.Close()

	cmd := exec.Command(execute.Cmd, execute.Args(execute_dir_path)...)

	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Dir = execute.Dir
	cmd.Env = append(os.Environ(), execute.Env...)
	cmd.Env = append(cmd.Env, "ANSIBLE_CACHE_PLUGIN_CONNECTION="+constants.GET_DATA_DIR()+"/"+execute.OwnerType+"/"+execute.OwnerName+"/fact")

	if err := cmd.Start(); err != nil {
		execute.R_Error = errors.New("failed to start command " + cmd.String() + " : " + err.Error())
		err = os.Remove(execute_dir_path)
		if err != nil {
			logrus.Warn(err.Error())
		}
		execute.mutex.Unlock()
		return
	}

	runningProcesses[pid] = cmd.Process

	logrus.Trace("started command " + cmd.String())
	os.WriteFile(execute_dir_path+"/execute.command", []byte(cmd.String()), 0666)
	os.WriteFile(execute_dir_path+"/execute.yaml", []byte(execute.ToString(execute_dir_path, pid)), 0666)

	if err := lockFile.Truncate(0); err != nil {
		execute.R_Error = errors.New("failed to truncate lockFile : " + err.Error())
	}
	_, err = lockFile.WriteString(pid)
	if err != nil {
		execute.R_Error = errors.New("failed to write pid " + err.Error())
	}

	execute.R_Pid = pid

	execute.mutex.Unlock()
	err = cmd.Wait()
	if err != nil {
		logrus.Warn(err.Error())
	}

	delete(runningProcesses, pid)

	if execute.PostExec != nil {
		_, err := logFile.WriteString("\n\n\nPangee Cluster *****************************************************************\n")
		if err != nil {
			logrus.Warn(err.Error())
		}
		logs, err := os.ReadFile(logFilePath)
		if err != nil {
			return
		}
		logStr := string(logs)
		if strings.LastIndex(logStr, "PLAY RECAP *********************************************************************") < 0 {
			msg := "\033[31m\033[01m\033[05m  No PLAY RECAP found.\033[0m\n" + "\033[31m\033[01m\033[05m  执行出错。\033[0m\n"
			_, err = logFile.WriteString(msg)
			if err != nil {
				return
			}
			logrus.Warn("No ansbile-playbook recap.")
			return
		}
		recap := logStr[strings.LastIndex(logStr, "PLAY RECAP *********************************************************************")+81:]
		recap = recap[:strings.Index(recap, "\n\n")]

		lines := strings.Split(recap, "\n")
		status := []ExecuteExitNodeStatus{}
		for _, line := range lines {
			if strings.Index(line, "failed=") > 0 && strings.Index(line, "unreachable=") > 0 && strings.Index(line, "ok=") > 0 {
				status = append(status, parseAnsibleRecapLine(line))
			}
		}
		success := len(status) >= 0
		for _, nodestatus := range status {
			if nodestatus.Unreachable != "0" || nodestatus.Failed != "0" {
				success = false
			}
		}

		// 对于有operation的任务，检查执行结果
		if execute.OwnerType == "cluster" && strings.Contains(pid, "/") {
			splited := strings.Split(pid, "/")
			cssr := CheckStepStatusRequest{
				Cluster:   execute.OwnerName,
				Operation: splited[0],
				Step:      splited[1],
			}
			cssrResponse, err := CheckStepStatusExec(cssr)
			if err == nil {
				common.SaveYamlFile(execute_dir_path+"/result.yaml", cssrResponse)
			}
		}

		// 保存任务执行状态 status.yaml
		endTime := time.Now()
		st, _ := time.Parse(time.RFC3339Nano, startTime)
		statusTemp := map[string]any{
			"startTime": startTime,
			"endTime":   endTime.Format(time.RFC3339Nano),
			"duration":  endTime.Sub(st).Milliseconds(),
			"status":    "success",
		}
		if success {
			statusTemp["status"] = "success"
		} else {
			statusTemp["status"] = "failed"
		}
		common.SaveYamlFile(execute_dir_path+"/status.yaml", statusTemp)

		exitStatus := ExecuteExitStatus{
			Success:    success,
			NodeStatus: status,
			Pid:        pid,
			ExecuteDir: execute_dir_path,
		}

		// 执行 PostExec
		message, err := execute.PostExec(exitStatus)
		if err != nil {
			_, err = logFile.WriteString("Error in PostExec: " + err.Error() + "\n")
			if err != nil {
				logrus.Warn(err.Error())
			}
		}

		if success {
			task := SuccessTask{
				Type:      execute.Type,
				Timestamp: time.Now(),
				Pid:       pid,
			}
			if err := AddSuccessTask(execute.OwnerType, execute.OwnerName, task); err != nil {
				logrus.Warn("failed to add success task: ", err)
			}
		}

		_, err = logFile.WriteString(message)
		if err != nil {
			logrus.Warn(err.Error())
		}
	}
}

type ExecuteExitStatus struct {
	Success    bool
	NodeStatus []ExecuteExitNodeStatus
	Pid        string
	ExecuteDir string
}

type ExecuteExitNodeStatus struct {
	NodeName    string
	OK          string
	Changed     string
	Unreachable string
	Failed      string
	Skipped     string
	Rescued     string
	Ignored     string
}

func trimColoredText(text string) string {
	result := text
	if strings.Index(result, "[0;32m") > 0 {
		result = result[strings.Index(result, "[0;32m")+6:]
	}
	if strings.Index(result, "[0;33m") > 0 {
		result = result[strings.Index(result, "[0;33m")+6:]
	}
	if strings.Index(result, "\033[0m") > 0 {
		result = result[:strings.Index(result, "\033[0m")]
	}
	return result
}

func parseAnsibleRecapLine(line string) ExecuteExitNodeStatus {
	result := make(map[string]string)
	s1 := strings.Split(line, ":")

	result["node"] = strings.Trim(s1[0], " ")
	s2 := strings.Split(s1[1], " ")
	for _, kv := range s2 {
		s3 := strings.Split(kv, "=")
		if len(s3) == 1 {
			continue
		}
		k := strings.Trim(s3[0], " ")
		v := strings.Trim(s3[1], " ")
		result[k] = v
	}

	nodeStatus := ExecuteExitNodeStatus{
		NodeName:    trimColoredText(result["node"]),
		OK:          result["ok"],
		Changed:     trimColoredText(result["changed"]),
		Unreachable: result["unreachable"],
		Failed:      result["failed"],
		Skipped:     result["skipped"],
		Rescued:     result["rescured"],
		Ignored:     result["ignored"],
	}
	return nodeStatus
}
