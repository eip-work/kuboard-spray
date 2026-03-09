package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/cluster"
	"github.com/opencmit/pangee-cluster/api/cluster/backup"
	"github.com/opencmit/pangee-cluster/api/cluster/cluster_access"
	"github.com/opencmit/pangee-cluster/api/cluster/health_check"
	"github.com/opencmit/pangee-cluster/api/cluster/operation_v2"
	"github.com/opencmit/pangee-cluster/api/cluster/state"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/api/fact"
	"github.com/opencmit/pangee-cluster/api/filebrowser"
	"github.com/opencmit/pangee-cluster/api/private_key"
	"github.com/opencmit/pangee-cluster/api/resource"
	"github.com/opencmit/pangee-cluster/api/ssh"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/opencmit/pangee-cluster/log"
	"github.com/opencmit/pangee-cluster/login"
	"github.com/opencmit/pangee-cluster/vue"
	"github.com/sirupsen/logrus"
)

func setupRouter() *gin.Engine {

	router := gin.Default()

	root := router.Group("/")
	root.POST("/api/login", login.AuthHandler)
	root.POST("/api/validate_password", login.ValidatePassword)

	api := root.Group("/api", login.JWTAuthMiddleware())

	api.GET("/profile", login.GetProfile)
	api.POST("/profile/change_password", login.ChangePassword)

	api.GET("/clusters", cluster.ListClusters)
	api.POST("/clusters", cluster.CreateCluster)
	api.GET("/clusters/:cluster", cluster.GetCluster)
	api.PUT("/clusters/:cluster", cluster.ModifyCluster)
	api.DELETE("/clusters/:cluster", cluster.DeleteCluster)
	api.PATCH("/clusters/:cluster", cluster.RenameCluster)

	api.POST("/clusters/:cluster/change_resource_package_version", cluster.ChangeResourcePackageVersion)
	api.POST("/clusters/:cluster/upgrade_all_nodes", operation_v2.UpgradeCluster)
	api.POST("/clusters/:cluster/upgrade_master_nodes", operation_v2.UpgradeCluster)
	api.POST("/clusters/:cluster/upgrade_single_node", operation_v2.UpgradeCluster)
	api.POST("/clusters/:cluster/upgrade_multi_nodes", operation_v2.UpgradeCluster)
	api.POST("/clusters/:cluster/download_binaries", operation_v2.DownloadBinaries)
	api.POST("/clusters/:cluster/drain_node", operation_v2.DrainNode)
	api.POST("/clusters/:cluster/uncordon_node", operation_v2.UncordonNode)

	api.POST("/clusters/:cluster/operation/:operation/step/:step", operation_v2.ExecuteStep)
	api.GET("/clusters/:cluster/operation/:operation/step/:step", operation_v2.CheckStepStatus)
	api.GET("/clusters/:cluster/history/:operation/step/:step", operation_v2.ListOperationStepHistory)
	api.GET("/clusters/:cluster/history/:operation/step/:step/:time", operation_v2.GetHistoryStepResult)

	api.GET("/clusters/:cluster/access/kubeconfig", cluster_access.GetKubeConfig)

	api.GET("/clusters/:cluster/backup", backup.ListBackup)
	api.POST("/clusters/:cluster/backup/remove", backup.RemoveBackup)

	api.POST("/clusters/:cluster/state/ping", state.Ping)
	api.GET("/clusters/:cluster/state/nodes", state.GetNodes)
	api.GET("/clusters/:cluster/state/etcd_member_health", state.CheckEtcdEndpoints)
	api.GET("/clusters/:cluster/state/version", state.CheckClusterVersion)
	api.GET("/clusters/:cluster/state/pods_on_node/:node", state.GetPodsOnNode)
	api.GET("/clusters/:cluster/state/check_cert_expiration", state.CheckCertExpiration)

	api.GET("/clusters/:cluster/health_check/connectivity_check", health_check.CheckConnectivity)
	api.GET("/clusters/:cluster/health_check/details", health_check.CheckConnectivityDetails)

	api.GET("/execute/:owner_type/:owner_name/tail/:pid/:file", command.TailFile)
	api.GET("/execute/:owner_type/:owner_name/tail-v2/:operation/:step/:time/:file", command.TailFile) // operation/step/time
	api.DELETE("/execute/:owner_type/:owner_name/kill/:pid", command.ExecuteKill)
	api.DELETE("/execute/:owner_type/:owner_name/kill-v2/:operation/:step/:time", command.ExecuteKill)

	api.POST("/facts/:node_owner_type/:node_owner/:node", fact.GetNodeFacts)
	api.GET("/ssh/:node_owner_type/:node_owner/:node", ssh.ShellWs)

	api.GET("/private-keys/:owner_type/:owner_name", private_key.ListPrivateKey)
	api.GET("/private-keys/:owner_type/:owner_name/:name", private_key.GetPrivateKey)
	api.POST("/private-keys/:owner_type/:owner_name/:name", private_key.UploadPrivateKeyFile)
	api.DELETE("/private-keys/:owner_type/:owner_name/:name", private_key.DeletePrivateKey)

	api.GET("/resources", resource.ListResources)
	api.GET("/resources/local/:name", resource.GetLocalResource)
	api.POST("/resources/upload", resource.UploadResource)
	api.POST("/resources/:name/download", resource.DownloadResourceDependancy)
	api.DELETE("/resources/:name", resource.DeleteResource)

	api.GET("/filebrowser", filebrowser.List)

	vue.ServeVue(router, root)

	return router
}

func main() {

	initLogrus()

	router := setupRouter()

	port := "9080"

	logrus.Trace(runtime.GOARCH)
	if runtime.GOARCH == "aarch64" || runtime.GOARCH == "arm64" {
		port = "9081"
	}

	router.Run(":" + constants.GetEnvDefault("PANGEE_CLUSTER_PORT", port))
}

func initLogrus() {
	logrus.SetFormatter(new(log.KuboardLogFormatter))
	// logrus.SetFormatter(&logrus.TextFormatter{})

	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	value := os.Getenv("PANGEE_CLUSTER_LOGRUS_LEVEL")
	if value == "" {
		value = "trace"
	}
	level, err := logrus.ParseLevel(value)
	if err == nil {
		fmt.Println("设置日志级别为 " + value)
		logrus.SetLevel(level)
	} else {
		fmt.Println("请检查 PANGEE_CLUSTER_LOGRUS_LEVEL 的值，可选的有 panic / fatal / error / warn / info / debug / trace ，当前为： " + value)
		logrus.SetLevel(logrus.InfoLevel)
	}

}
