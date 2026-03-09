package resource

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
)

func UploadResource(c *gin.Context) {
	// 创建临时目录保存上传的ZIP文件
	tempDir, err := os.MkdirTemp("", "resource-upload")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to create temp directory", err)
		return
	}
	tempFilePath := filepath.Join(tempDir, "resource-pack.zip")
	defer os.RemoveAll(tempDir)

	// 限制上传大小50MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 50<<20)

	// 处理文件上传或URL下载
	err = handleFileSource(c, tempFilePath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to get resource file", err)
		return
	}

	// 解压ZIP文件并获取解压目录
	unzipDir, err := unzipResource(tempFilePath, tempDir)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to unzip resource file", err)
		return
	}

	// 验证package.yaml并提取版本号
	version, err := validateAndExtractVersion(unzipDir)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "invalid resource package", err)
		return
	}

	// 移动文件到目标目录
	err = moveResourceFiles(unzipDir, version)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to move resource files", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"version": version, // 返回获取到的版本号
		},
	})
}

// handleFileSource 处理文件上传或从URL下载
func handleFileSource(c *gin.Context, tempFilePath string) error {
	file, err := c.FormFile("file")
	if err == nil { // 报文中包含压缩包
		if !strings.HasSuffix(strings.ToLower(file.Filename), ".zip") {
			return fmt.Errorf("please upload a zip file")
		}

		err = c.SaveUploadedFile(file, tempFilePath)
		if err != nil {
			return fmt.Errorf("failed to save file: %v", err)
		}
	} else if sourceUrl := c.Request.FormValue("sourceUrl"); sourceUrl != "" { // 报文中包含sourceUrl
		return downloadFromUrl(sourceUrl, tempFilePath, c)
	} else {
		return fmt.Errorf("either file or sourceUrl is required")
	}
	return nil
}

// downloadFromUrl 从指定URL下载文件
func downloadFromUrl(sourceUrl, tempFilePath string, c *gin.Context) error {
	timeoutStr := c.Request.FormValue("timeout")
	if timeoutStr == "" {
		timeoutStr = "60"
	}
	timeoutSec, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return fmt.Errorf("invalid timeout value: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   time.Duration(timeoutSec) * time.Second,
	}

	if proxyUrl := c.Request.FormValue("proxyUrl"); proxyUrl != "" { // 通过代理下载
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			return fmt.Errorf("failed to parse proxyUrl: %v", err)
		}
		logrus.Infof("通过代理 %s 下载 %s", proxy.String(), sourceUrl)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
			Timeout: time.Duration(timeoutSec) * time.Second,
		}
	}

	resp, err := client.Get(sourceUrl)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: server returned status %s", resp.Status)
	}

	out, err := os.Create(tempFilePath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	logrus.Infof("保存文件到 %s", tempFilePath)
	return nil
}

// unzipResource 解压资源文件
func unzipResource(tempFilePath, tempDir string) (string, error) {
	unzipDir := filepath.Join(tempDir, "unzip")
	err := os.MkdirAll(unzipDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create unzip directory: %v", err)
	}

	// 使用命令行解压
	unzipCmd := exec.Command("unzip", tempFilePath, "-d", unzipDir)
	if err := unzipCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to unzip file: %v", err)
	}

	entries, err := os.ReadDir(unzipDir)
	if err != nil {
		return "", fmt.Errorf("failed to read unzip directory: %v", err)
	}

	// 如果只有一个目录，则进入该目录
	if len(entries) == 1 && entries[0].IsDir() {
		unzipDir = filepath.Join(unzipDir, entries[0].Name())
	}

	return unzipDir, nil
}

// validateAndExtractVersion 验证package.yaml并提取版本号
func validateAndExtractVersion(unzipDir string) (string, error) {
	packageYamlPath := filepath.Join(unzipDir, "package.yaml")
	info, err := os.Stat(packageYamlPath)
	if err != nil || info.IsDir() {
		return "", fmt.Errorf("package.yaml not found in zip root directory")
	}

	// 使用yq获取版本号
	yqCmd := exec.Command("yq", "eval", ".metadata.version", packageYamlPath)
	output, err := yqCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get version from package.yaml: %v", err)
	}

	version := strings.TrimSpace(string(output))
	if version == "" {
		return "", fmt.Errorf("version not found in package.yaml")
	}

	// 验证版本号格式
	if !strings.Contains(version, ".") {
		return "", fmt.Errorf("invalid version format, expected semantic versioning")
	}

	return version, nil
}

// moveResourceFiles 移动资源文件到目标目录
func moveResourceFiles(unzipDir, version string) error {
	uploadDir := filepath.Join(constants.GET_DATA_RESOURCE_DIR(), version)
	contentDir := filepath.Join(uploadDir, "content")

	// 清理已存在的目录
	if err := os.RemoveAll(uploadDir); err != nil {
		logrus.Warnf("failed to remove existing directory: %v", err)
	}

	// 创建目录结构
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %v", err)
	}

	// 复制package.yaml
	packageYamlPath := filepath.Join(unzipDir, "package.yaml")
	destPackageYaml := filepath.Join(uploadDir, "package.yaml")
	cpCmd := exec.Command("cp", packageYamlPath, destPackageYaml)
	if err := cpCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy package.yaml: %v", err)
	}

	// 移动内容到content目录
	logrus.Infof("moving %s to %s", unzipDir+"/", contentDir+"/")

	// 先复制所有文件，再删除原文件
	cpCmd = exec.Command("cp", "-r", unzipDir+"/.", contentDir+"/")
	if err := cpCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy resource files: %v", err)
	}

	return nil
}
