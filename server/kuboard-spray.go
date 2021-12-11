package main

import (
	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/resource"
	"github.com/eip-work/kuboard-spray/login"
	"github.com/eip-work/kuboard-spray/vue"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	router := gin.Default()

	root := router.Group("/kuboardspray/:kuboardsprayID")
	root.POST("/api/login", login.AuthHandler)

	api := root.Group("/api", login.JWTAuthMiddleware())
	api.GET("/clusters", cluster.ListClusters)
	api.GET("/clusters/:name", cluster.GetCluster)

	api.GET("/resources", resource.ListResources)
	api.GET("/resources/:name", resource.GetResource)

	vue.ServeVue(router, root)

	return router
}

func main() {
	router := setupRouter()

	router.Run(":8006")
}
