package os_mirror

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ListOsMirrorsRequest struct {
	Type string `form:"type"`
}

func ListOsMirrors(c *gin.Context) {

	var req ListOsMirrorsRequest

	c.ShouldBind(&req)

	files, err := ioutil.ReadDir(constants.GET_DATA_MIRROR_DIR())
	if err != nil {

		err2 := common.CreateDirIfNotExists(constants.GET_DATA_MIRROR_DIR())
		if err2 != nil {
			common.HandleError(c, http.StatusInternalServerError, "cannot create folder: "+constants.GET_DATA_MIRROR_DIR(), err2)
			return
		}
		logrus.Warning("已成功创建该文件夹: " + constants.GET_DATA_MIRROR_DIR())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "success",
			"data":    []string{},
		})
		return
	}

	data := []string{}
	for _, file := range files {
		if file.IsDir() && strings.Index(file.Name(), req.Type) == 0 {
			data = append(data, file.Name())
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    data,
	})
}
