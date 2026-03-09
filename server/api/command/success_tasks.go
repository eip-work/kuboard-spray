package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/opencmit/pangee-cluster/constants"
)

type SuccessTask struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"succeedTime"`
	Pid       string    `json:"pid"`
}

type SuccessTasks []SuccessTask

func AddSuccessTask(ownerType string, ownerName string, task SuccessTask) error {
	tasks, err := ReadSuccessTasks(ownerType, ownerName)
	if err != nil {
		return err
	}
	tasks = append(tasks, task)
	// content, err := json.Marshal(tasks)
	// if err != nil {
	// 	return errors.New("failed tp marshal tasks: " + err.Error())
	// }
	// return os.WriteFile(successFilePath(ownerType, ownerName), content, 0666)

	file, err := os.Create(successFilePath(ownerType, ownerName))
	if err != nil {
		return errors.New("failed to create file " + successFilePath(ownerType, ownerName) + " : " + err.Error())
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tasks); err != nil {
		return errors.New("failed to marshal tasks to file " + successFilePath(ownerType, ownerName) + " : " + err.Error())
	}
	return nil
}

func ReadSuccessTasks(ownerType string, ownerName string) (SuccessTasks, error) {
	var tasks SuccessTasks
	content, err := os.ReadFile(successFilePath(ownerType, ownerName))
	if err != nil {
		content = []byte("[]")
	}
	contentReader := bytes.NewReader(content)
	decoder := json.NewDecoder(contentReader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&tasks); err != nil {
		return nil, errors.New("failed to unmarshal file " + successFilePath(ownerType, ownerName) + " : " + err.Error())
	}

	return tasks, nil
}

func successFilePath(ownerType string, ownerName string) string {
	return constants.GET_DATA_DIR() + "/" + ownerType + "/" + ownerName + "/success.json"
}
