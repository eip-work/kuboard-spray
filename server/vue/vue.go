package vue

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var kuboardSprayWebDir string = os.Getenv("KUBOARD_SPRAY_WEB_DIR")

type KuboardSprayId struct {
	Id string `uri:"kuboardsprayID" binding:"required"`
}

// ServeVue ServeVue
func ServeVue(router *gin.Engine, root *gin.RouterGroup) {

	if kuboardSprayWebDir == "" {
		kuboardSprayWebDir = "/root/git/kuboard-spray/web/dist"
	}

	static := router.Group("/kuboardspray")

	router.LoadHTMLGlob(kuboardSprayWebDir + "/**.html")
	static.StaticFS("/static/fonts", http.Dir(kuboardSprayWebDir+"/fonts"))
	static.StaticFS("/static/js", http.Dir(kuboardSprayWebDir+"/js"))
	static.StaticFile("/static/favicon.ico", kuboardSprayWebDir+"/favicon.ico")

	static.GET("/kuboardspray/:kuboardsprayID/*path", func(c *gin.Context) {
		var id KuboardSprayId
		c.ShouldBindUri(&id)

		if c.Request.RequestURI == "/kuboardspray/"+id.Id+"/" {
			c.HTML(http.StatusOK, "index.html", gin.H{})
			return
		}

		path := strings.Replace(c.Request.RequestURI, id.Id, "static", 1)

		c.Redirect(http.StatusMovedPermanently, path)
	})
}
