package operation_v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
)

func CheckStepStatus(c *gin.Context) {
	var request command.CheckStepStatusRequest
	c.ShouldBindUri(&request)
	c.ShouldBindJSON(&request)

	stepStatus, err := command.CheckStepStatusExec(request)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    stepStatus,
	})

}
