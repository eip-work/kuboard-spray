package state

import (
	"encoding/json"
	"fmt"

	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/sirupsen/logrus"
)

// 将kubectl get nodes中的name从IP改为设定名字
func IPLink(commandResp string, cluster string) (string, interface{}) {
	var commandRespObj interface{}
	err := json.Unmarshal([]byte(commandResp), &commandRespObj)
	if err != nil {
		fmt.Print(err)
	}
	respData := commandRespObj.(map[string]interface{})
	items := respData["items"].([]interface{})

	clusterMetadata, _ := cluster_common.ClusterMetadataByName(cluster)
	hosts := common.MapGet(clusterMetadata.Inventory, "all.hosts").(map[string]interface{})
	IPMap := map[string]string{}
	for name, host := range hosts {
		hostConfigMap := host.(map[string]interface{})
		IP, exist := hostConfigMap["ip"]
		if exist {
			IPMap[IP.(string)] = name
		}
	}
	// fmt.Println(IPMap)

	for _, item := range items {
		itemMap := item.(map[string]interface{})
		metadataMap := itemMap["metadata"].(map[string]interface{})
		IP := metadataMap["name"].(string)
		metadataMap["name"] = IPMap[IP]
	}
	var stdout string
	if err := json.Unmarshal([]byte(stdout), commandRespObj); err != nil {
		logrus.Warning("failed to parse", err)
	}
	return stdout, commandRespObj
}
