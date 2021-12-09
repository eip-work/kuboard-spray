package vue

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var kuboardSprayWebDir string = os.Getenv("KUBOARD_SPRAY_WEB_DIR")

type Version struct {
	Version string `uri:"version" binding:"required"`
}

// ServeVue ServeVue
func ServeVue(router *gin.Engine) {

	if kuboardSprayWebDir == "" {
		kuboardSprayWebDir = "/root/git/kuboard-spray/web/dist"
	}
	// authorized := router.Group("/", login.HTTPLoginRequired(), login.AuthRequired())

	static := router.Group("/kuboardspray/")

	router.LoadHTMLGlob(kuboardSprayWebDir + "/**.html")
	static.StaticFS("/static/fonts", http.Dir(kuboardSprayWebDir+"/fonts"))
	static.StaticFS("/static/js", http.Dir(kuboardSprayWebDir+"/js"))
	static.StaticFile("/static/favicon.ico", kuboardSprayWebDir+"/favicon.ico")

	static.GET("/:version/*path", func(c *gin.Context) {
		var version Version
		c.ShouldBindUri(&version)

		if c.Request.RequestURI == "/kuboardspray/"+version.Version+"/" {
			c.HTML(http.StatusOK, "index.html", gin.H{})
			return
		}

		path := strings.Replace(c.Request.RequestURI, version.Version, "static", 1)

		c.Redirect(http.StatusMovedPermanently, path)
	})
}
