package private_key

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ClusterPrivateKeyListRequest struct {
	Cluster string `uri:"cluster" binding:"required"`
}

func ListPrivateKey(c *gin.Context) {
	var req ClusterPrivateKeyListRequest
	c.ShouldBindUri(&req)

	files, err := ioutil.ReadDir(ClusterPrivateKeyPath(req.Cluster))
	logrus.Info("["+req.Cluster+"]", ClusterPrivateKeyPath(req.Cluster))
	if err != nil {

		err2 := common.CreateDirIfNotExists(ClusterPrivateKeyPath(req.Cluster))
		if err2 != nil {
			common.HandleError(c, http.StatusInternalServerError, "cannot create folder: "+ClusterPrivateKeyPath(req.Cluster), err)
			return
		}
		logrus.Warning("已成功创建该文件夹: " + ClusterPrivateKeyPath(req.Cluster))
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

func ClusterPrivateKeyPath(clusterName string) string {
	return constants.GET_DATA_INVENTORY_DIR() + "/" + clusterName + "/private-key"
}
