package cmd

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GET_DATA_DIR() string {
	dir, _ := os.Getwd()
	lastindex := strings.LastIndex(dir, "/")
	dir = dir[0:lastindex] + "/data"
	return dir
}

var userDir = GET_DATA_DIR() + "/user"
var userFilePath = userDir + "/user.yaml"

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "重置 admin 的密码",
	Long:  "重置 PangeeCluster 管理员 admin 的密码。",
	Run:   reset,
}

func reset(cmd *cobra.Command, args []string) {
	resetPassword()
}

func resetPassword() {
	logrus.Info("try to reset-password")

	userFileContent, err := os.ReadFile(userFilePath)
	if err != nil {
		logrus.Error("不能读取文件 "+userFilePath, err)
		os.Exit(-1)
	}

	result := map[string]*UserInfo{}
	if err := yaml.Unmarshal(userFileContent, result); err != nil {
		logrus.Error("解析 YAML 出错: "+userFilePath, err)
		os.Exit(-1)
	}

	password := os.Getenv("PANGEE_CLUSTER_DEFAULT_ADMIN_PASSWORD")

	admin := result["admin"]
	admin.Password = password

	content, err := yaml.Marshal(result)
	if err != nil {
		logrus.Info(string(content))
		logrus.Error("序列化 YAML 出错: ", err)
		os.Exit(-1)
	}

	err = os.WriteFile(userFilePath, content, 0666)
	if err != nil {
		logrus.Error("写入文件时出错: "+userFilePath, err)
		os.Exit(-1)
	}

	logrus.Info("已将 admin 的密码重置为 PangeeCluster123")

}
