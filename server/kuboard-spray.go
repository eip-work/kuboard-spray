package main

import (
	"github.com/eip-work/kuboard-spray/vue"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	router := gin.Default()

	vue.ServeVue(router)

	return router
}

func main() {
	router := setupRouter()

	router.Run(":8006")
}
