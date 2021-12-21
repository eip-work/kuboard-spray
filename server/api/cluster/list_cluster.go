package cluster

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ListClusters(c *gin.Context) {

	files, err := ioutil.ReadDir(constants.GET_DATA_CLUSTER_DIR())
	if err != nil {

		err1 := common.CreateDirIfNotExists(constants.GET_DATA_DIR())
		if err1 != nil {
			common.HandleError(c, http.StatusInternalServerError, "cannot create folder: "+constants.GET_DATA_DIR(), err1)
			return
		}
		err2 := common.CreateDirIfNotExists(constants.GET_DATA_CLUSTER_DIR())
		if err2 != nil {
			common.HandleError(c, http.StatusInternalServerError, "cannot create folder: "+constants.GET_DATA_CLUSTER_DIR(), err2)
			return
		}
		logrus.Warning("已成功创建该文件夹: " + constants.GET_DATA_CLUSTER_DIR())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "success",
			"data":    []string{},
		})
		return
	}

	data := []string{}
	for _, file := range files {
		data = append(data, file.Name())
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    data,
	})

}
