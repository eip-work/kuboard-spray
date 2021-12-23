package command

import (
	"errors"
	"os"
	"syscall"

	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

func LockOwner(owner_type string, owner_name string) (*os.File, error) {
	lockFilePath := constants.GET_DATA_DIR() + "/" + owner_type + "/" + owner_name + "/inventory.lastrun"
	logrus.Trace("lockFilePath: ", lockFilePath)
	lockFile, err := os.OpenFile(lockFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return lockFile, errors.New("Cannot open file " + lockFilePath + " : " + err.Error())
	}

	err = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	if err != nil {
		return nil, errors.New("Cannot lock file " + lockFilePath + " : " + err.Error())
	}

	return lockFile, nil
}

func UnlockOwner(lockFile *os.File) {
	syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
	lockFile.Close()
}
