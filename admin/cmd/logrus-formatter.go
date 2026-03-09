package cmd

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// PangeeClusterLogFormatter Custom log format definition
type PangeeClusterLogFormatter struct{}

// Format 格式化日志输出
func (s *PangeeClusterLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05.000")

	msg := fmt.Sprintf("%s | %5v | %s\n", timestamp, entry.Level.String(), entry.Message)
	return []byte(msg), nil
}
