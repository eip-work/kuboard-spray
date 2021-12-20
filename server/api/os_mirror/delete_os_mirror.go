package os_mirror

import (
	"net/http"
	"os"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

func DeleteMirror(c *gin.Context) {

	var req GetMirrorRequest
	c.ShouldBindUri(&req)

	if err := os.RemoveAll(constants.GET_DATA_MIRROR_DIR() + "/" + req.Name); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to remove os_mirror", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})

}
