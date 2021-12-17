package command

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

type Run struct {
	Cmd     string
	Args    []string
	Cluster string
}

func (run *Run) ToString() string {
	result := "--" + run.Cluster + "-- [" + run.Cmd + ", "
	for index, arg := range run.Args {
		result += arg
		if index < len(run.Args)-1 {
			result += ", "
		}
	}
	result += "]"
	return result
}

func (run *Run) Run() ([]byte, []byte, error) {
	// create lockfile
	lockFilePath := constants.GET_DATA_INVENTORY_DIR() + "/" + run.Cluster + "/inventory.lastrun"
	logrus.Trace("lockFilePath: ", lockFilePath)
	lockFile, err := os.OpenFile(lockFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, nil, errors.New("Cannot open file " + lockFilePath + " : " + err.Error())
	}

	// get lock
	err = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	if err != nil {
		return nil, nil, errors.New("Cannot lock file " + lockFilePath + " : " + err.Error())
	}

	defer lockFile.Close()
	defer syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)

	logrus.Trace("run command: ", run.ToString())
	cmd := exec.Command(run.Cmd, run.Args...)
	cmd.Env = []string{"ANSIBLE_HOST_KEY_CHECKING=False"}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, errors.New("Cannot pip stderr: " + err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, errors.New("Cannot pip stdout: " + err.Error())
	}

	err = cmd.Start()

	errContent := ""
	go func() {
		reader := bufio.NewReader(stderr)
		for {
			line, _, err := reader.ReadLine()
			if err != nil || io.EOF == err {
				break
			}
			errContent += string(line) + "\n"
		}
		logrus.Trace("terminated err reader thread.")
	}()

	outContent := ""
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, _, err := reader.ReadLine()
			if err != nil || io.EOF == err {
				break
			}
			outContent += string(line) + "\n"
		}
		logrus.Trace("terminated out reader thread.")
	}()

	cmd.Wait()
	return []byte(outContent), []byte(errContent), err
}
