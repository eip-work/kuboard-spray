package os_mirror

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GetMirrorRequest struct {
	Name string `uri:"name"`
}

func GetMirror(c *gin.Context) {
	var req GetMirrorRequest
	c.ShouldBindUri(&req)

	inventory, err := common.ParseYamlFile(MirrorInventoryPath(req.Name))
	if err != nil {
		logrus.Trace("cannot parse file: ", err)
		// 不通过 kuboardspray 安装的软件源没有 inventory.yaml
	}
	statusFilePath := constants.GET_DATA_MIRROR_DIR() + "/" + req.Name + "/status.yaml"
	status, err := common.ParseYamlFile(statusFilePath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot pase file: "+statusFilePath, err)
	}

	history, err := command.ReadTaskHistory("mirror", req.Name)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot read status", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"inventory": inventory,
			"history":   history,
			"name":      req.Name,
			"status":    status,
		},
	})
}

func MirrorInventoryPath(name string) string {
	return constants.GET_DATA_MIRROR_DIR() + "/" + name + "/inventory.yaml"
}
