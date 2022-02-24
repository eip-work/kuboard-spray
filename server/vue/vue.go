package vue

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

// ServeVue ServeVue
func ServeVue(router *gin.Engine, root *gin.RouterGroup) {

	kuboardSprayWebDir := constants.GetKuboardSprayWebDir()

	static := router.Group("/")

	router.LoadHTMLGlob(kuboardSprayWebDir + "/**.html")
	static.StaticFS("/fonts", http.Dir(kuboardSprayWebDir+"/fonts"))
	static.StaticFS("/js", http.Dir(kuboardSprayWebDir+"/js"))
	static.StaticFS("/img", http.Dir(kuboardSprayWebDir+"/img"))
	static.StaticFile("/static/favicon.ico", kuboardSprayWebDir+"/favicon.ico")
	static.StaticFile("/version.json", kuboardSprayWebDir+"/version.json")
	static.StaticFile("/index.html", kuboardSprayWebDir+"/index.html")

	static.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
}
