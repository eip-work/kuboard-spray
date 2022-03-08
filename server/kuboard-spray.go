package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/cluster/backup"
	"github.com/eip-work/kuboard-spray/api/cluster/cis_scan"
	"github.com/eip-work/kuboard-spray/api/cluster/cluster_access"
	"github.com/eip-work/kuboard-spray/api/cluster/health_check"
	"github.com/eip-work/kuboard-spray/api/cluster/operation"
	"github.com/eip-work/kuboard-spray/api/cluster/state"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/api/fact"
	"github.com/eip-work/kuboard-spray/api/os_mirror"
	"github.com/eip-work/kuboard-spray/api/private_key"
	"github.com/eip-work/kuboard-spray/api/resource"
	"github.com/eip-work/kuboard-spray/api/ssh"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/eip-work/kuboard-spray/log"
	"github.com/eip-work/kuboard-spray/login"
	"github.com/eip-work/kuboard-spray/vue"
	"github.com/gin-gonic/gin"
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

	api.POST("/clusters/:cluster/install_cluster", operation.InstallCluster)
	api.POST("/clusters/:cluster/remove_node", operation.RemoveNode)
	api.POST("/clusters/:cluster/add_node", operation.AddNode)
	api.POST("/clusters/:cluster/sync_nginx_config", operation.SyncNginxConfigActions)
	api.POST("/clusters/:cluster/sync_etcd_address", operation.SyncEtcdConfigActions)
	api.POST("/clusters/:cluster/install_addon", operation.InstallAddon)
	api.POST("/clusters/:cluster/remove_addon", operation.RemoveAddon)
	api.POST("/clusters/:cluster/change_resource_package_version", cluster.ChangeResourcePackageVersion)
	api.POST("/clusters/:cluster/upgrade_all_nodes", operation.UpgradeCluster)
	api.POST("/clusters/:cluster/upgrade_master_nodes", operation.UpgradeCluster)
	api.POST("/clusters/:cluster/upgrade_single_node", operation.UpgradeCluster)
	api.POST("/clusters/:cluster/upgrade_multi_nodes", operation.UpgradeCluster)
	api.POST("/clusters/:cluster/download_binaries", operation.DownloadBinaries)
	api.POST("/clusters/:cluster/drain_node", operation.DrainNode)
	api.POST("/clusters/:cluster/uncordon_node", operation.UncordonNode)
	api.POST("/clusters/:cluster/renew_cert", operation.RenewCert)
	api.POST("/clusters/:cluster/sync_container_engine_params", operation.SyncContainerEngineParams)
	api.POST("/clusters/:cluster/backup_etcd", operation.BackupEtcd)

	api.POST("/clusters/:cluster/cis_scan", cis_scan.CisScan)

	api.GET("/clusters/:cluster/access/kubeconfig", cluster_access.GetKubeConfig)

	api.GET("/clusters/:cluster/backup", backup.ListBackup)
	api.POST("/clusters/:cluster/backup/remove", backup.RemoveBackup)

	api.POST("/clusters/:cluster/state/ping", state.Ping)
	api.GET("/clusters/:cluster/state/nodes", state.GetNodes)
	api.GET("/clusters/:cluster/state/etcd_member_health", state.CheckEtcdEndpoints)
	api.GET("/clusters/:cluster/state/addons", state.CheckAddonStatus)
	api.GET("/clusters/:cluster/state/version", state.CheckClusterVersion)
	api.GET("/clusters/:cluster/state/pods_on_node/:node", state.GetPodsOnNode)
	api.GET("/clusters/:cluster/state/check_cert_expiration", state.CheckCertExpiration)

	api.GET("/clusters/:cluster/health_check/connectivity_check", health_check.CheckConnectivity)
	api.GET("/clusters/:cluster/health_check/details", health_check.CheckConnectivityDetails)

	api.GET("/execute/:owner_type/:owner_name/tail/:pid/:file", command.TailFile)
	api.DELETE("/execute/:owner_type/:owner_name/kill/:pid", command.ExecuteKill)

	api.POST("/facts/:node_owner_type/:node_owner/:node", fact.GetNodeFacts)
	api.GET("/ssh/:node_owner_type/:node_owner/:node", ssh.ShellWs)

	api.GET("/private-keys/:owner_type/:owner_name", private_key.ListPrivateKey)
	api.GET("/private-keys/:owner_type/:owner_name/:name", private_key.GetPrivateKey)
	api.POST("/private-keys/:owner_type/:owner_name/:name", private_key.UploadPrivateKeyFile)
	api.DELETE("/private-keys/:owner_type/:owner_name/:name", private_key.DeletePrivateKey)

	api.GET("/resources", resource.ListResources)
	api.GET("/resources/:name", resource.GetResource)
	api.GET("/resources/:name/release_note", resource.GetResourceReleaseNote)
	api.POST("/resources/:name/download", resource.CreateAndDownloadResource)
	api.POST("/resources/:name/reload", resource.ReloadResource)
	api.DELETE("/resources/:name", resource.DeleteResource)

	api.GET("/mirrors", os_mirror.ListOsMirrors)
	api.POST("/mirrors", os_mirror.CreateOsMirror)
	api.POST("/mirrors/:name/install", os_mirror.InstallOsMirror)
	api.GET("/mirrors/:name", os_mirror.GetMirror)
	api.PUT("/mirrors/:name", os_mirror.ModifyOsMirro)
	api.DELETE("/mirrors/:name", os_mirror.DeleteMirror)

	vue.ServeVue(router, root)

	return router
}

func main() {

	initLogrus()

	router := setupRouter()

	port := ":8006"

	logrus.Trace(runtime.GOARCH)
	if runtime.GOARCH == "aarch64" || runtime.GOARCH == "arm64" {
		port = ":8007"
	}

	router.Run(constants.GetEnvDefault("KUBOARD_SPRAY_PORT", port))
	// s := &http.Server{
	// 	Addr:         ":8006",
	// 	Handler:      router,
	// 	ReadTimeout:  120 * time.Second,
	// 	WriteTimeout: 120 * time.Second,
	// }
	// s.ListenAndServe()
}

func initLogrus() {
	logrus.SetFormatter(new(log.KuboardLogFormatter))
	// logrus.SetFormatter(&logrus.TextFormatter{})

	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	value := os.Getenv("KUBOARD_SPRAY_LOGRUS_LEVEL")
	if value == "" {
		value = "trace"
	}
	level, err := logrus.ParseLevel(value)
	if err == nil {
		fmt.Println("设置日志级别为 " + value)
		logrus.SetLevel(level)
	} else {
		fmt.Println("请检查 KUBOARD_SPRAY_LOGRUS_LEVEL 的值，可选的有 panic / fatal / error / warn / info / debug / trace ，当前为： " + value)
		logrus.SetLevel(logrus.InfoLevel)
	}

}
