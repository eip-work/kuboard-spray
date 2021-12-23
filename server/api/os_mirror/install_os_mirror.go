package os_mirror

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type InstallOsMirrorRequest struct {
	Name    string `uri:"name" binding:"required"`
	Verbose bool   `json:"verbose"`
	VVV     bool   `json:"vvv"`
}

func InstallOsMirror(c *gin.Context) {

	var req InstallOsMirrorRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		success := status.Success
		var message string
		if success {
			message = "\033[32m[ " + "OS Mirror has been installed successfully, please go back to the os mirror page for information about how to access the OS Mirror." + " ]\033[0m \n"
			message += "\033[32m[ " + "镜像源已成功安装，请回到详情页面查看如何访问该镜像源。" + " ]\033[0m \n"
		} else {
			message = "\033[31m\033[01m\033[05m[" + "Failed to install OS Mirror. Please review the logs and fix the problem." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "镜像源安装失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
		}
		return "\n" + message, nil
	}

	env := []string{}
	if req.Verbose {
		env = append(env, "ANSIBLE_DISPLAY_ARGS_TO_STDOUT=True")
	}
	logrus.Trace(req.Verbose, env)

	dir, err := os.Getwd()
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot find current working dir.", err)
		return
	}
	command := command.Execute{
		OwnerType: "mirror",
		OwnerName: req.Name,
		Cmd:       "ansible-playbook",
		Args: func(run_dir string) []string {
			if req.VVV {
				return []string{"-i", run_dir + "/inventory.yaml", "mirror_ubuntu.yaml", "-vvv"}
			}
			return []string{"-i", run_dir + "/inventory.yaml", "mirror_ubuntu.yaml"}
		},
		Dir:  dir + "/ansible-script",
		Type: "install",
		PreExec: func(run_dir string) error {
			inventoryPath := MirrorInventoryPath(req.Name)
			inventory, err := ioutil.ReadFile(inventoryPath)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(run_dir+"/inventory.yaml", inventory, 0655)
		},
		PostExec: postExec,
		Env:      env,
	}

	if err := command.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to InstallCluster. ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"pid": command.R_Pid,
		},
	})

}
