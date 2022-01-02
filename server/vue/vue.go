package vue

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var kuboardSprayWebDir string = os.Getenv("KUBOARD_SPRAY_WEB_DIR")

// ServeVue ServeVue
func ServeVue(router *gin.Engine, root *gin.RouterGroup) {

	if kuboardSprayWebDir == "" {
		dir, _ := os.Getwd()
		lastindex := strings.LastIndex(dir, "/")
		kuboardSprayWebDir = dir[0:lastindex] + "/web/dist"
	}

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
