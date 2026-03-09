package vue

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
)

// ServeVue ServeVue
func ServeVue(router *gin.Engine, root *gin.RouterGroup) {

	pangeeClusterWebDir := constants.GetPangeeClusterWebDir()

	static := router.Group("/")

	resourceDir := constants.GET_DATA_RESOURCE_DIR()
	logrus.Info("resourceDir: ", resourceDir)
	static.StaticFS("/resource-package", http.Dir(resourceDir))

	router.LoadHTMLGlob(pangeeClusterWebDir + "/**.html")
	static.StaticFS("/assets", http.Dir(pangeeClusterWebDir+"/assets"))
	static.StaticFile("/static/favicon.ico", pangeeClusterWebDir+"/favicon.ico")
	static.StaticFile("/version.json", pangeeClusterWebDir+"/version.json")
	static.StaticFile("/index.html", pangeeClusterWebDir+"/index.html")

	static.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

}
