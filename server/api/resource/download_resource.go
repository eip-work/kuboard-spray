package resource

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
)

func DownloadResourceDependancy(c *gin.Context) {
	buf, err := c.GetRawData()
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read request", err)
		return
	}

	logrus.Trace(string(buf))

	var downloadReq map[string]interface{}
	err = json.Unmarshal(buf, &downloadReq)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse request", err)
		return
	}

	version, ok := common.MapGet(downloadReq, "version").(string)
	if !ok || version == "" {
		common.HandleError(c, http.StatusBadRequest, "invalid version parameter", nil)
		return
	}

	versionDir := constants.GET_DATA_RESOURCE_DIR() + "/" + version

	_, errExist := os.ReadDir(versionDir)
	if errExist != nil {
		common.HandleError(c, http.StatusConflict, "资源包不存在", errExist)
		return
	}

	downloadScriptPath := versionDir + "/content/download-dependency.sh"
	if _, err := os.Stat(downloadScriptPath); os.IsNotExist(err) {
		common.HandleError(c, http.StatusInternalServerError, "download script not found", err)
		return
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {
		var message string
		if status.Success {
			message = "\033[32m[ " + "PangeeCluster resource package has been loaded successfully." + " ]\033[0m \n"
			message += "\033[32m[ " + "PangeeCluster 资源包已成功加载到本地，请回到资源包页面查看。" + " ]\033[0m \n"
		} else {
			message = "\033[31m\033[01m\033[05m[" + "Failed to download pangeecluster resource package, please reivew the log and fix the issues." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "下载资源包失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
		}
		return "\n" + message, nil
	}

	cmd := command.Execute{
		OwnerType: "resource",
		OwnerName: version,
		Cmd:       downloadScriptPath,
		Args: func(execute_dir string) []string {
			return []string{
				common.MapGetString(downloadReq, "retries"),
				common.MapGetString(downloadReq, "downloadArchitecture"),
				common.MapGetString(downloadReq, "enableProxyOnDownload"),
				common.MapGetString(downloadReq, "httpProxy"),
			}
		},
		Type:     "download",
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Failed to install cluster", err)
		return
	}

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data: gin.H{
			"pid": cmd.R_Pid,
		},
	})
}
