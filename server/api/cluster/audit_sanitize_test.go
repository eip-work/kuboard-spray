package cluster

import (
	"testing"

	"github.com/eip-work/kuboard-spray/common"
)

func baseInventory() map[string]interface{} {
	return map[string]interface{}{
		"all": map[string]interface{}{
			"children": map[string]interface{}{
				"target": map[string]interface{}{
					"children": map[string]interface{}{
						"k8s_cluster": map[string]interface{}{
							"vars": map[string]interface{}{},
						},
					},
				},
			},
		},
	}
}

func TestSanitizeAuditLogConfigEnabledDefaults(t *testing.T) {
	inv := baseInventory()
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.kubernetes_audit", true)

	sanitizeAuditLogConfig(inv)

	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_path"); v != "/var/log/audit/kube-apiserver-audit.log" {
		t.Fatalf("expected default audit_log_path, got %v", v)
	}
	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxage"); v != 30 {
		t.Fatalf("expected default audit_log_maxage=30, got %v", v)
	}
	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxbackups"); v != 10 {
		t.Fatalf("expected default audit_log_maxbackups=10, got %v", v)
	}
	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxsize"); v != 100 {
		t.Fatalf("expected default audit_log_maxsize=100, got %v", v)
	}
}

func TestSanitizeAuditLogConfigEnabledPreserveExisting(t *testing.T) {
	inv := baseInventory()
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.kubernetes_audit", true)
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_path", "/custom/audit.log")
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxage", 7)
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxbackups", 3)
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxsize", 50)

	sanitizeAuditLogConfig(inv)

	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_path"); v != "/custom/audit.log" {
		t.Fatalf("expected preserve audit_log_path, got %v", v)
	}
	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxage"); v != 7 {
		t.Fatalf("expected preserve audit_log_maxage=7, got %v", v)
	}
	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxbackups"); v != 3 {
		t.Fatalf("expected preserve audit_log_maxbackups=3, got %v", v)
	}
	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxsize"); v != 50 {
		t.Fatalf("expected preserve audit_log_maxsize=50, got %v", v)
	}
}

func TestSanitizeAuditLogConfigDisabledClears(t *testing.T) {
	inv := baseInventory()
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.kubernetes_audit", false)
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_path", "/custom/audit.log")
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxage", 7)
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxbackups", 3)
	common.MapSet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxsize", 50)

	sanitizeAuditLogConfig(inv)

	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_path"); v != nil {
		t.Fatalf("expected audit_log_path cleared, got %v", v)
	}
	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxage"); v != nil {
		t.Fatalf("expected audit_log_maxage cleared, got %v", v)
	}
	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxbackups"); v != nil {
		t.Fatalf("expected audit_log_maxbackups cleared, got %v", v)
	}
	if v := common.MapGet(inv, "all.children.target.children.k8s_cluster.vars.audit_log_maxsize"); v != nil {
		t.Fatalf("expected audit_log_maxsize cleared, got %v", v)
	}
}


