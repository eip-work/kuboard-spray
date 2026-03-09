package command

import (
	"errors"
	"os"

	"github.com/opencmit/pangee-cluster/constants"
)

type TaskStatus struct {
	TaskType     string        `json:"task_type"`
	TaskName     string        `json:"task_name"`
	CurrentPid   string        `json:"current_pid"`
	Processing   bool          `json:"processing"`
	SuccessTasks []SuccessTask `json:"success_tasks"`
}

func ReadTaskHistory(taskType string, taskName string) (*TaskStatus, error) {

	taskStatus := TaskStatus{
		TaskType: taskType,
		TaskName: taskName,
	}

	successTasks, err := ReadSuccessTasks(taskType, taskName)
	if err != nil {
		return nil, err
	}
	taskStatus.SuccessTasks = successTasks

	taskStatus.Processing = false
	lockFile, err := LockOwner(taskType, taskName)
	if err != nil {
		taskStatus.Processing = true
	}
	lockFilePath := constants.GET_DATA_DIR() + "/" + taskType + "/" + taskName + "/inventory.lastrun"
	pid, err := os.ReadFile(lockFilePath)
	if err != nil {
		return nil, errors.New("cannot read pid: " + err.Error())
	}
	UnlockOwner(lockFile)
	taskStatus.CurrentPid = string(pid)

	return &taskStatus, nil
}
