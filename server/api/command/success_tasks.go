package command

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/eip-work/kuboard-spray/constants"
)

type SuccessTask struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"succeedTime"`
	Pid       string    `json:"pid"`
}

type SuccessTasks []SuccessTask

func AddSuccessTask(cluster string, task SuccessTask) error {
	tasks, err := ReadSuccessTasks(cluster)
	if err != nil {
		return err
	}
	tasks = append(tasks, task)
	content, err := json.Marshal(tasks)
	if err != nil {
		return errors.New("failed tp marshal tasks: " + err.Error())
	}
	return ioutil.WriteFile(successFilePath(cluster), content, 0777)
}

func ReadSuccessTasks(cluster string) (SuccessTasks, error) {
	var tasks SuccessTasks
	content, err := ioutil.ReadFile(successFilePath(cluster))
	if err != nil {
		content = []byte("[]")
	}
	if err := json.Unmarshal(content, &tasks); err != nil {
		return nil, errors.New("failed to unmarshal file " + successFilePath(cluster) + " : " + err.Error())
	}

	return tasks, nil
}

func successFilePath(cluster string) string {
	return constants.GET_DATA_INVENTORY_DIR() + "/" + cluster + "/success.json"
}
