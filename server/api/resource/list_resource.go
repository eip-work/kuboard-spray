package resource

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
)

func ListResources(c *gin.Context) {
	err1 := common.CreateDirIfNotExists(constants.GET_DATA_DIR())
	if err1 != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot create folder: "+constants.GET_DATA_DIR(), err1)
		return
	}
	err2 := common.CreateDirIfNotExists(constants.GET_DATA_RESOURCE_DIR())
	if err2 != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot create folder: "+constants.GET_DATA_RESOURCE_DIR(), err2)
		return
	}
	fileInfoList, err := os.ReadDir(constants.GET_DATA_RESOURCE_DIR())
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot read folder: "+constants.GET_DATA_RESOURCE_DIR(), err)
		return
	}

	result := []string{}
	for _, dir := range fileInfoList {
		if dir.Name()[0:1] != "." {
			result = append(result, dir.Name())
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    result,
	})
}
