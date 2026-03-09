package common

import (
	"bytes"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func ParseYamlFile(filePath string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	if err := yaml.Unmarshal(content, result); err != nil {
		return nil, err
	}
	return result, err
}

func SaveYamlFile(filePath string, content interface{}) error {

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)

	err := encoder.Encode(&content)

	if err != nil {
		return errors.New("failed to marshal content : " + err.Error())
	}
	logrus.Trace("write yaml file: ", filePath)
	return os.WriteFile(filePath, buffer.Bytes(), 0666)
}
