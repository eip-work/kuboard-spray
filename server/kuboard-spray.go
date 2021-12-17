package main

import (
	"fmt"
	"os"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/api/node"
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

	api.POST("/clusters/:cluster/install", cluster.InstallCluster)
	api.GET("/clusters/:cluster/history/:pid/tail/:file", command.TailFile)
	// router.GET("/kuboardspray/:kuboardsprayID/api/clusters/:cluster/history/:pid/tail/:file", command.TailFile)

	api.POST("/clusters/:cluster/facts/:node", node.GetNodeFacts)

	api.GET("/clusters/:cluster/private-keys", private_key.ListPrivateKey)
	api.GET("/clusters/:cluster/private-keys/:name", private_key.GetPrivateKey)
	api.POST("/clusters/:cluster/private-keys/:name", private_key.UploadPrivateKeyFile)
	api.DELETE("/clusters/:cluster/private-keys/:name", private_key.DeletePrivateKey)

	api.GET("/resources", resource.ListResources)
	api.GET("/resources/:name", resource.GetResource)

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
