package main

import (
	"fmt"
	"os"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/api/fact"
	"github.com/eip-work/kuboard-spray/api/os_mirror"
	"github.com/eip-work/kuboard-spray/api/private_key"
	"github.com/eip-work/kuboard-spray/api/resource"
	"github.com/eip-work/kuboard-spray/log"
	"github.com/eip-work/kuboard-spray/login"
	"github.com/eip-work/kuboard-spray/vue"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func setupRouter() *gin.Engine {

	router := gin.Default()

	root := router.Group("/kuboardspray/:kuboardsprayID")
	root.POST("/api/login", login.AuthHandler)

	api := root.Group("/api", login.JWTAuthMiddleware())
	api.GET("/clusters", cluster.ListClusters)
	api.POST("/clusters", cluster.CreateCluster)
	api.GET("/clusters/:cluster", cluster.GetCluster)
	api.PUT("/clusters/:cluster", cluster.ModifyCluster)
	api.DELETE("/clusters/:cluster", cluster.DeleteCluster)

	api.POST("/clusters/:cluster/install", cluster.InstallCluster)

	api.GET("/execute/:owner_type/:owner_name/tail/:pid/:file", command.TailFile)
	api.DELETE("/execute/:owner_type/:owner_name/kill/:pid", command.ExecuteKill)

	api.POST("/facts/:node_owner_type/:node_owner/:node", fact.GetNodeFacts)

	api.GET("/private-keys/:owner_type/:owner_name", private_key.ListPrivateKey)
	api.GET("/private-keys/:owner_type/:owner_name/:name", private_key.GetPrivateKey)
	api.POST("/private-keys/:owner_type/:owner_name/:name", private_key.UploadPrivateKeyFile)
	api.DELETE("/private-keys/:owner_type/:owner_name/:name", private_key.DeletePrivateKey)

	api.GET("/resources", resource.ListResources)
	api.GET("/resources/:name", resource.GetResource)

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

	router.Run(":8006")
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
