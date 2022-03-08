package backup

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

type ListBackupRequest struct {
	Cluster string `uri:"cluster"`
}

type BackupInfo struct {
	Name       string `json:"backup_name"`
	Size       int64  `json:"size"`
	NodeName   string `json:"node_name"`
	MemberName string `json:"etcd_member_name"`
}

func ListBackup(c *gin.Context) {
	var req ListBackupRequest

	c.BindUri(&req)

	backupDir := constants.GET_DATA_CLUSTER_DIR() + "/" + req.Cluster + "/backup"

	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		c.JSON(http.StatusOK, common.KuboardSprayResponse{
			Code:    http.StatusOK,
			Message: "success",
			Data:    []BackupInfo{},
		})
		return
	}

	result := []BackupInfo{}

	for _, node_dir := range files {
		if !node_dir.IsDir() {
			continue
		}
		members, err := ioutil.ReadDir(backupDir + "/" + node_dir.Name())
		if err != nil {
			continue
		}
		for _, member := range members {
			backups, err := ioutil.ReadDir(backupDir + "/" + node_dir.Name() + "/" + member.Name())
			if err != nil {
				continue
			}
			for _, backup := range backups {
				if !strings.Contains(backup.Name(), ".tgz") || backup.IsDir() {
					continue
				}
				result = append(result, BackupInfo{Name: backup.Name(), Size: backup.Size(), MemberName: member.Name(), NodeName: node_dir.Name()})
			}
		}
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    result,
	})
}
