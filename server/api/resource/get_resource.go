package resource

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"gopkg.in/yaml.v3"
)

type GetResourceRequest struct {
	Name string `uri:"name"`
}

func GetLocalResource(c *gin.Context) {
	var req GetResourceRequest
	c.ShouldBindUri(&req)

	packageContent, err := os.ReadFile(GET_RESOURCE_YAML_PATH(req.Name))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot open file: "+GET_RESOURCE_YAML_PATH(req.Name), err)
		return
	}

	res := gin.H{}
	yaml.Unmarshal(packageContent, res)

	history, err := command.ReadTaskHistory("resource", req.Name)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot read task history", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"package": res,
			"history": history,
		},
	})
}

func GET_RESOURCE_YAML_PATH(name string) string {
	return constants.GET_DATA_RESOURCE_DIR() + "/" + name + "/package.yaml"
}
