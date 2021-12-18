package common

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func MapGet(m map[string]interface{}, path string) interface{} {
	var result interface{}
	defer func() {
		if err := recover(); err != nil {
			logrus.Warn("error:", err)
			result = nil
		}
	}()
	p := strings.Split(path, ".")
	var temp map[string]interface{}
	temp = m
	for i := 0; i < len(p); i++ {
		if i == len(p)-1 {
			return temp[p[i]]
		}
		if temp[p[i]] != nil {
			temp = temp[p[i]].(map[string]interface{})
		}
	}
	return result
}

func MapSet(m map[string]interface{}, path string, value interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warn("error:", err)
		}
	}()
	p := strings.Split(path, ".")
	var temp map[string]interface{}
	temp = m
	for i := 0; i < len(p); i++ {
		if i == len(p)-1 {
			temp[p[i]] = value
			return
		}
		if temp[p[i]] == nil {
			temp[p[i]] = map[string]interface{}{}
		} else {
			temp = temp[p[i]].(map[string]interface{})
		}
	}
}
