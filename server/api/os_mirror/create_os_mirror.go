package os_mirror

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type CreateOsMirrorRequest struct {
	Name string `json:"kuboardspray_os_mirror_name" binding:"required"`
	Type string `json:"kuboardspray_os_mirror_type" binding:"required"`
	Kind string `json:"kuboardspray_os_mirror_kind" binding:"required"`
	Url  string `json:"kuboardspray_os_mirror_url"`
}

func CreateOsMirror(c *gin.Context) {
	var req CreateOsMirrorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, http.StatusBadRequest, "请求参数出错", err)
		return
	}

	name := req.Type + "-" + req.Name

	logrus.Trace("create ", name)

	mirrorDir := constants.GET_DATA_MIRROR_DIR() + "/" + name

	_, err := ioutil.ReadDir(mirrorDir)
	if err == nil {
		common.HandleError(c, http.StatusConflict, "conflict with a existing OS Mirror", err)
		return
	}

	if err := common.CreateDirIfNotExists(mirrorDir); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to create folder", err)
		return
	}

	status := gin.H{
		"type": req.Type,
		"kind": req.Kind,
		"url":  req.Url,
		"name": req.Name,
	}

	inventoryObj := map[string]interface{}{}
	if req.Kind == "provision" {
		status["status"] = "created"

		inventoryObj := getInventoryTemplate()

		inventoryContent, err := yaml.Marshal(inventoryObj)

		if err != nil {
			common.HandleError(c, http.StatusInternalServerError, "failed to marshal inventory ", err)
			return
		}

		inventoryFilePath := mirrorDir + "/inventory.yaml"
		if err := ioutil.WriteFile(inventoryFilePath, inventoryContent, 0666); err != nil {
			common.HandleError(c, http.StatusInternalServerError, "failed to save "+inventoryFilePath, err)
			return
		}
	} else {
		status["status"] = "success"
	}

	str, err := yaml.Marshal(status)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to marshal status", err)
		return
	}

	statusFilePath := mirrorDir + "/status.yaml"
	if err := ioutil.WriteFile(statusFilePath, str, 0666); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to save "+statusFilePath, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"inventory": inventoryObj,
			"status":    status,
		},
	})
}

func getInventoryTemplate() gin.H {
	template := gin.H{
		"all": gin.H{
			"hosts": gin.H{
				"localhost": gin.H{
					"ansible_connection": "local",
				},
			},
			"children": gin.H{
				"target": gin.H{
					"hosts": gin.H{},
					"vars": gin.H{
						"apt_mirror_client":          false,
						"apt_mirror_apache_root_dir": "/var/www/html",
					},
				},
			},
		},
	}

	return template
}
