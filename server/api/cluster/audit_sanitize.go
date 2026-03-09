package cluster

import "github.com/eip-work/kuboard-spray/common"

// sanitizeAuditLogConfig ensures audit log fields are present with safe defaults when enabled,
// and removed when disabled, to avoid passing empty flags to kube-apiserver.
func sanitizeAuditLogConfig(inventory map[string]interface{}) {
	const base = "all.children.target.children.k8s_cluster.vars"
	enabled, ok := common.MapGet(inventory, base+".kubernetes_audit").(bool)
	if !ok {
		return
	}
	if enabled {
		if common.MapGet(inventory, base+".audit_log_path") == nil {
			common.MapSet(inventory, base+".audit_log_path", "/var/log/audit/kube-apiserver-audit.log")
		}
		if common.MapGet(inventory, base+".audit_log_maxage") == nil {
			common.MapSet(inventory, base+".audit_log_maxage", 30)
		}
		if common.MapGet(inventory, base+".audit_log_maxbackups") == nil {
			common.MapSet(inventory, base+".audit_log_maxbackups", 10)
		}
		if common.MapGet(inventory, base+".audit_log_maxsize") == nil {
			common.MapSet(inventory, base+".audit_log_maxsize", 100)
		}
	} else {
		common.MapDelete(inventory, base+".audit_log_path")
		common.MapDelete(inventory, base+".audit_log_maxage")
		common.MapDelete(inventory, base+".audit_log_maxbackups")
		common.MapDelete(inventory, base+".audit_log_maxsize")
	}
}


