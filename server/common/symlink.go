package common

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Symlink(target string, link string) {
	if _, err := os.Lstat(link); os.IsNotExist(err) {
		err = os.Symlink(target, link)
		if err != nil {
			logrus.Warn("Failed to create Symlink: ", target, link)
		}
	}
}
