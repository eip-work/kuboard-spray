package resource

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func CreateAndDownloadResource(c *gin.Context) {
	templateMethod(c, false)
}

func ReloadResource(c *gin.Context) {
	templateMethod(c, true)
}

func templateMethod(c *gin.Context, canUseExisting bool) {
	var req GetResourceRequest
	c.ShouldBindUri(&req)

	buf, err := c.GetRawData()

	logrus.Trace(string(buf))

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read request", err)
		return
	}

	var downloadReq map[string]interface{}
	err = json.Unmarshal(buf, &downloadReq)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse request", err)
		return
	}

	version := common.MapGet(downloadReq, "package.metadata.version").(string)

	versionDir := constants.GET_DATA_RESOURCE_DIR() + "/" + version

	_, errExist := os.ReadDir(versionDir)
	if errExist == nil && !canUseExisting {
		common.HandleError(c, http.StatusConflict, "资源包已存在，不能重复创建", errExist)
		return
	}

	if err := common.CreateDirIfNotExists(versionDir); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Create dir failed. ", err)
		return
	}

	pkg, _ := yaml.Marshal(downloadReq["package"])

	pkgExists := false
	entries, _ := os.ReadDir(versionDir)
	for _, entry := range entries {
		if entry.Name() == "package.yaml" {
			pkgExists = true
		}
	}
	if !pkgExists {
		if err := ioutil.WriteFile(versionDir+"/package.yaml", pkg, 0655); err != nil {
			common.HandleError(c, http.StatusInternalServerError, "Write package.yaml failed. ", err)
			return
		}
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		success := status.Success
		var message string
		if success {
			message = "\033[32m[ " + "Kuboardspray resource package has been loaded successfully." + " ]\033[0m \n"
			message += "\033[32m[ " + "Kuboardspray 资源包已成功加载到本地，请回到资源包页面查看。" + " ]\033[0m \n"
		} else {
			message = "\033[31m\033[01m\033[05m[" + "Failed to download kuboardspray resource package, please reivew the log and fix the issues." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "下载资源包失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
		}
		return "\n" + message, nil
	}

	cmd := command.Execute{
		OwnerType: "resource",
		OwnerName: version,
		Cmd:       "./pull-resource-package.sh",
		Args: func(execute_dir string) []string {
			return []string{common.MapGet(downloadReq, "downloadFrom").(string) + ":" + version}
		},
		Type:     "download",
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to InstallCluster. ", err)
		return
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data: gin.H{
			"pid": cmd.R_Pid,
		},
	})
}
