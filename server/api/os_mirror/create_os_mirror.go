package os_mirror

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type CreateOsMirrorRequest struct {
	Name   string `json:"kuboardspray_os_mirror_name" binding:"required"`
	Type   string `json:"kuboardspray_os_mirror_type" binding:"required"`
	Kind   string `json:"kuboardspray_os_mirror_kind" binding:"required"`
	Url    string `json:"kuboardspray_os_mirror_url"`
	Params gin.H  `json:"params"`
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
		"type":   req.Type,
		"kind":   req.Kind,
		"url":    req.Url,
		"name":   req.Name,
		"params": req.Params,
	}

	inventoryObj := map[string]interface{}{}
	if req.Kind == "provision" {
		status["status"] = "created"

		inventoryObj := map[string]interface{}{}

		if req.Type == "docker_ubuntu" || req.Type == "ubuntu" {
			template := getDockerUbuntuMirrorInventoryTemplate()
			yaml.Unmarshal([]byte(template), inventoryObj)
		} else if req.Type == "ubuntu" {
			template := getUbuntuMirrorInventoryTemplate()
			yaml.Unmarshal([]byte(template), inventoryObj)
		} else {
			common.HandleError(c, http.StatusInternalServerError, "not supported yet", nil)
			return
		}

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

func getDockerUbuntuMirrorInventoryTemplate() string {
	return `all:
  children:
    target:
      hosts:
        mirror_node: {}
      vars:
        apt_mirror_client: false
        apt_mirror_dir: /var/spool/apt-mirror
        kuboardspray_get_gpg: true
        apt_mirror_repos:
        - "deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} focal stable"
  hosts:
    localhost:
      ansible_connection: local
    mirror_node: {}
`
}

func getUbuntuMirrorInventoryTemplate() string {
	return `all:
  children:
    target:
      hosts:
        mirror_node: {}
      vars:
        apt_mirror_client: false
        apt_mirror_dir: /var/spool/apt-mirror
        apt_mirror_repos:
          - deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} focal main restricted universe multiverse
          - deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} focal-backports main restricted universe multiverse
          - deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} focal-security main restricted universe multiverse
          - deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} focal-updates main restricted universe multiverse
          - deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} focal-proposed main restricted universe multiverse
  hosts:
    localhost:
      ansible_connection: local
    mirror_node: {}
`
}
