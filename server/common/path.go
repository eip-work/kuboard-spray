package common

import "os"

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func CreateDirIfNotExists(path string) error {
	if PathExists(path) {
		return nil
	}
	return os.Mkdir(path, os.ModePerm)
}
