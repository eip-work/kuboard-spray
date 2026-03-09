package command

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Run struct {
	Cmd     string
	Args    []string
	Env     []string
	Dir     string
	Timeout int
}

func (run *Run) ToString() string {
	result := "[" + run.Cmd + " "
	for index, arg := range run.Args {
		if index > 1 && run.Args[index-1] == "-a" {
			result += "'"
		}
		result += arg
		if index > 1 && run.Args[index-1] == "-a" {
			result += "'"
		}
		if index < len(run.Args)-1 {
			result += " "
		}
	}
	result += "]"
	return result
}

func (run *Run) Run() ([]byte, []byte, error) {

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

	cmderr := cmd.Start()

	flagOut := false
	flagErr := false
	errContent := ""

	if cmderr != nil {
		errContent = cmderr.Error()
	}
	go func() {
		buf := new(strings.Builder)
		reader := bufio.NewReader(stderr)
		for {
			_, err := io.Copy(buf, reader)
			if err != nil || io.EOF == err {
				break
			}
		}
		errContent = buf.String()
		logrus.Trace("terminated err reader thread.")
		flagErr = true
	}()

	outContent := ""
	go func() {
		buf := new(strings.Builder)
		reader := bufio.NewReader(stdout)
		for {
			_, err := io.Copy(buf, reader)
			if err != nil || io.EOF == err {
				break
			}
		}
		outContent = buf.String()
		logrus.Trace("terminated out reader thread.")
		flagOut = true
	}()

	if run.Timeout > 0 {
		go func() {
			time.Sleep(time.Duration(time.Second * time.Duration(run.Timeout)))
			err = cmd.Process.Kill()
			if err != nil {
				logrus.Warn(err.Error())
			}
			outContent += "KilledByPangeeCluster[Timeout]"
		}()
	}

	err = cmd.Wait()
	if err != nil {
		logrus.Warn(err.Error())
	}
	for !flagErr || !flagOut {
		logrus.Trace("wait..")
		time.Sleep(time.Millisecond * 100)
	}
	return []byte(outContent), []byte(errContent), err
}
