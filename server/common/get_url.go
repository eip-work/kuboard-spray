package common

import (
	"fmt"
)

func GetRepoUrl(source string) (string, error) {
	url := ""

	// TODO: 修改地址
	switch source {
	case "github":
		url = "https://raw.githubusercontent.com/opencmit/pangee-cluster-resource-package"
	case "gitee":
		url = "https://gitee.com/opencmit/pangee-cluster-resource-package"
	default:
		return url, fmt.Errorf("unsupported source type: %s", source)
	}

	return url, nil
}

func GetDownloadUrl(source string) (string, error) {
	url := ""

	// TODO: 修改地址
	switch source {
	case "github":
		url = "https://github.com/opencmit/pangee-cluster-resource-package/releases/download"
	case "gitee":
		url = "https://gitee.com/opencmit/pangee-cluster-resource-package/releases/download"
	default:
		return url, fmt.Errorf("unsupported source type: %s", source)
	}

	return url, nil
}
