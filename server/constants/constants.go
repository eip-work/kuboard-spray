package constants

import (
	"os"
	"strings"
)

func GET_DATA_CLUSTER_DIR() string {
	return GET_DATA_DIR() + "/cluster"
}

func GET_DATA_RESOURCE_DIR() string {
	return GET_DATA_DIR() + "/resource"
}

func GET_DATA_MIRROR_DIR() string {
	return GET_DATA_DIR() + "/mirror"
}

func GET_DATA_DIR() string {
	dir, _ := os.Getwd()
	lastindex := strings.LastIndex(dir, "/")
	dir = dir[0:lastindex] + "/data"
	return dir
}

func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

func GET_ADHOC_CFG_PATH() string {
	dir, _ := os.Getwd()
	return dir + "/ansible-script/ansible-adhoc.cfg"
}

func GET_ANSIBLE_SCRIPT_DIR() string {
	dir, _ := os.Getwd()
	return dir + "/ansible-script"
}
