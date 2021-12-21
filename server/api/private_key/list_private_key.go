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
	OwnerType string `uri:"owner_type" binding:"required"`
	OwnerName string `uri:"owner_name" binding:"required"`
}

func ListPrivateKey(c *gin.Context) {
	var req ClusterPrivateKeyListRequest
	c.ShouldBindUri(&req)

	files, err := ioutil.ReadDir(PrivateKeyDir(req.OwnerType, req.OwnerName))
	logrus.Info("["+req.OwnerType+"/"+req.OwnerName+"]", PrivateKeyDir(req.OwnerType, req.OwnerName))
	if err != nil {

		err2 := common.CreateDirIfNotExists(PrivateKeyDir(req.OwnerType, req.OwnerName))
		if err2 != nil {
			common.HandleError(c, http.StatusInternalServerError, "cannot create folder: "+PrivateKeyDir(req.OwnerType, req.OwnerName), err)
			return
		}
		logrus.Warning("已成功创建该文件夹: " + PrivateKeyDir(req.OwnerType, req.OwnerName))
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

func PrivateKeyDir(ownerType string, ownerName string) string {
	return constants.GET_DATA_DIR() + "/" + ownerType + "/" + ownerName + "/private-key"
}
