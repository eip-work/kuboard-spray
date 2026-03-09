package cluster

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"gopkg.in/yaml.v3"
)

type CreateClusterReq struct {
	Name            string `json:"name" binding:"required"`
	ResourcePackage string `json:"resource_package" binding:"required"`
}

func CreateCluster(c *gin.Context) {
	var req CreateClusterReq
	e := c.ShouldBindJSON(&req)

	if e != nil {
		common.HandleError(c, http.StatusBadRequest, "request should include field name and resourcePackage", e)
		return
	}

	clusterDir := constants.GET_DATA_CLUSTER_DIR() + "/" + req.Name

	_, err := os.ReadDir(clusterDir)
	if err == nil {
		common.HandleError(c, http.StatusConflict, "conflict with a existing cluster", err)
		return
	}
	err = common.CreateDirIfNotExists(clusterDir)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to create folder", err)
		return
	}

	inventoryFilePath := clusterDir + "/inventory.yaml"

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to create inventory file "+inventoryFilePath, err)
		return
	}

	resourcePackagePath := constants.GET_DATA_RESOURCE_DIR() + "/" + req.ResourcePackage + "/content"
	templateContent, err := os.ReadFile(resourcePackagePath + "/inventory-template.yaml")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot read file: "+resourcePackagePath+"/inventory-template.yaml", err)
		return
	}
	template := string(templateContent)
	template = strings.ReplaceAll(template, "PANGEECLUSTER_RESOURCE_PACKAGE", req.ResourcePackage)
	for strings.Contains(template, "GENERATE_UUID_PLACEHOLDER") {
		template = strings.Replace(template, "GENERATE_UUID_PLACEHOLDER", uuid.New().String(), 1)
	}
	for strings.Contains(template, "GENERATE_PASSWORD_PLACEHOLDER") {
		template = strings.Replace(template, "GENERATE_PASSWORD_PLACEHOLDER", common.GeneratePassword(12), 1)
	}

	inventory := map[string]interface{}{}
	yaml.Unmarshal([]byte(template), inventory)

	// resourcePackage, err := common.ParseYamlFile(resourcePackagePath + "/package.yaml")
	// if err != nil {
	// 	common.HandleError(c, http.StatusInternalServerError, "cannot parse package.yaml", err)
	// 	return
	// }

	common.MapSet(inventory, "all.vars.pangeecluster_resource_package_dir", resourcePackagePath)

	common.PopulatePangeeClusterVars(inventory, "cluster", req.Name)

	// addons := common.MapGet(resourcePackage, "data.addon").([]interface{})
	// for _, a := range addons {
	// 	addon := a.(map[string]interface{})
	// 	target := addon["target"].(string)
	// 	if addon["lifecycle"] != nil {
	// 		lifecycle := addon["lifecycle"].(map[string]interface{})
	// 		if lifecycle["install_by_default"] == true {
	// 			common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars."+target, true)
	// 		} else {
	// 			common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars."+target, false)
	// 		}
	// 	}
	// }

	// if runtime.GOARCH == "arm64" {
	// 	common.MapSet(inventory, "all.children.target.vars.container_manager", "docker")
	// }

	err = common.SaveYamlFile(inventoryFilePath, inventory)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to write inventory file: "+inventoryFilePath, err)
		return
	}

	data := gin.H{}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    data,
	})
}
