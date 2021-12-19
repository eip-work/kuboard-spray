package common

import (
	"errors"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func ParseYamlFile(filePath string) (map[string]interface{}, error) {
	content, err := ioutil.ReadFile(filePath)
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
	str, err := yaml.Marshal(content)
	if err != nil {
		return errors.New("failed to marshal content : " + err.Error())
	}
	logrus.Trace("write yaml file: ", filePath)
	return ioutil.WriteFile(filePath, str, 0666)
}
