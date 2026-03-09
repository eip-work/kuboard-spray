package filebrowser

import (
	"net/http"

	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
)

type ListRequest struct {
	RootPath string `form:"rootPath" binding:"required"`
}

type ListResponse struct {
	Name    string    `json:"name" binding:"required"`
	IsDir   bool      `json:"isDir" binding:"required"`
	Size    int64     `json:"size"`
	Mode    uint32    `json:"mode"`
	ModTime time.Time `json:"time"`
}

func List(c *gin.Context) {
	var req ListRequest
	c.ShouldBindQuery(&req)

	rootPath := constants.GET_DATA_DIR() + strings.Replace(req.RootPath, "data://", "", -1)

	entries, err := os.ReadDir(rootPath)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot read dir: "+rootPath, err)
		return
	}

	var result []ListResponse = make([]ListResponse, len(entries))

	for idx, entry := range entries {
		entry.Info()
		result[idx] = ListResponse{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		}
		fileInfo, err := entry.Info()
		if err == nil {
			result[idx].ModTime = fileInfo.ModTime()
			result[idx].Mode = uint32(fileInfo.Mode())
			result[idx].Size = fileInfo.Size()
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    result,
	})
}
