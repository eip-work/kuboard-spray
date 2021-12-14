package resource

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type GetResourceRequest struct {
	Name string `uri:"name"`
}

func GetResource(c *gin.Context) {
	var req GetResourceRequest
	c.ShouldBindUri(&req)

	packageContent, err := ioutil.ReadFile(GET_RESOURCE_YAML_PATH(req.Name))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot open file: "+GET_RESOURCE_YAML_PATH(req.Name), err)
		return
	}

	res := gin.H{}
	yaml.Unmarshal(packageContent, res)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"package": res,
		},
	})
}

func GET_RESOURCE_YAML_PATH(name string) string {
	return constants.GET_DATA_RESOURCE_DIR() + "/" + name + "/package.yaml"
}
