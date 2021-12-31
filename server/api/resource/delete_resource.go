package resource

import (
	"net/http"
	"os"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

func DeleteResource(c *gin.Context) {
	var req GetResourceRequest
	c.ShouldBindUri(&req)

	if err := os.RemoveAll(constants.GET_DATA_RESOURCE_DIR() + "/" + req.Name); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to remove resource dir.", err)
		return
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
	})
}
