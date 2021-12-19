package command

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Run struct {
	Cmd     string
	Args    []string
	Cluster string
	Env     []string
	Sync    bool
	Dir     string
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
	if run.Sync {
		lockFile, err := LockCluster(run.Cluster)
		if err != nil {
			return nil, nil, err
		}

		defer UnlockCluster(lockFile)
	}

	logrus.Trace("run command: ", run.ToString())
	cmd := exec.Command(run.Cmd, run.Args...)
	cmd.Env = append(os.Environ(), run.Env...)
	cmd.Dir = run.Dir

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
