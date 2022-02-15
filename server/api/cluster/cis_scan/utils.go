package cis_scan

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"gopkg.in/yaml.v3"
)

func GetCisCacheFilePath(cluster string, target string) string {
	result := constants.GET_DATA_CLUSTER_DIR() + "/" + cluster + "/cis_scan"
	temp := strings.ReplaceAll(target, "[", "_")
	temp = strings.ReplaceAll(temp, "]", "_")
	result += "/" + temp + ".yaml"
	return result
}

func GetCisCache(cluster string, target string) (*ansible_rpc.AnsibleResult, *time.Time, bool) {
	filePath := GetCisCacheFilePath(cluster, target)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, false
	}

	result := &ansible_rpc.AnsibleResult{}
	if err := yaml.Unmarshal(content, result); err != nil {
		return nil, nil, false
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, false
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, nil, false
	}
	modTime := fi.ModTime()

	return result, &modTime, true
}

func SaveCisCache(cluster string, target string, content *ansible_rpc.AnsibleResult) error {
	dir := constants.GET_DATA_CLUSTER_DIR() + "/" + cluster + "/cis_scan"
	common.CreateDirIfNotExists(dir)

	fileContent, err := yaml.Marshal(content)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(GetCisCacheFilePath(cluster, target), fileContent, 0655); err != nil {
		return err
	}

	return nil
}
