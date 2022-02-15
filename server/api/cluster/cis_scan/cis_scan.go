package cis_scan

import (
	"net/http"
	"time"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type CisScanRequest struct {
	ClusterName string `uri:"cluster" binding:"required"`
	Target      string `json:"target"`
	CacheMode   string `json:"cache_mode"`
}

func CisScan(c *gin.Context) {
	var req CisScanRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	if req.CacheMode != "ignore_cache" {
		result, modTime, ok := GetCisCache(req.ClusterName, req.Target)
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"code":          http.StatusOK,
				"message":       "success",
				"last_run_time": modTime,
				"data":          result,
			})
			return
		}
	}

	if req.CacheMode == "force_cache" {
		common.HandleError(c, http.StatusInternalServerError, "no cache found", nil)
		return
	}

	cmd := "/root/kube-bench/kube-bench run"

	if req.Target == "kube_control_plane" {
		cmd += " --config-dir=/root/kube-bench/cfg --json --version=1.23"
		cmd += " --targets=master,controlplane,node,policies"
	} else if req.Target == "etcd" {
		cmd += " --config-dir=/root/kube-bench/cfg --json --version=1.23"
		cmd += " --targets=etcd,policies"
	} else if req.Target == "kube_node" {
		cmd += " --config-dir=/root/kube-bench/cfg --json --version=1.23"
		cmd += " --targets=node,policies"
	} else if req.Target == "kube_control_plane[0]" {
		cmd += " --config-dir=/root/kube-bench/cfg --json --version=1.23"
		cmd += " --targets=master,controlplane,node,policies,etcd"
	}

	commands := []ansible_rpc.AnsibleCommandsRequest{
		{
			Name:    "scan",
			Command: cmd,
		},
	}
	result, err := ansible_rpc.ExecuteShellCommandsWithStrategy("cluster", req.ClusterName, req.Target, commands, "linear")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Failed to execute cis scan", err)
		return
	}

	SaveCisCache(req.ClusterName, req.Target, result)

	c.JSON(http.StatusOK, gin.H{
		"code":          http.StatusOK,
		"message":       "success",
		"last_run_time": time.Now(),
		"data":          result,
	})
}
