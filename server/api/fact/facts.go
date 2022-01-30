package fact

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GetNodeFactRequest struct {
	Node           string `uri:"node"`
	NodeOwnerType  string `uri:"node_owner_type"`
	NodeOwner      string `uri:"node_owner"`
	FromCache      bool   `json:"from_cache"`
	Ip             string `json:"ansible_host" binding:"required"`
	Port           string `json:"ansible_port" binding:"required"`
	User           string `json:"ansible_user" binding:"required"`
	Password       string `json:"ansible_password"`
	PrivateKeyFile string `json:"ansible_ssh_private_key_file"`
	Become         bool   `json:"ansible_become"`
	BecomeUser     string `json:"ansible_become_user"`
	BecomePassword string `json:"ansible_become_password"`
	SshCommonArgs  string `json:"ansible_ssh_common_args"`
	GatherSubset   string `json:"gather_subset"`
	Filter         string `json:"filter"`
}

func GetNodeFacts(c *gin.Context) {

	var req GetNodeFactRequest
	c.ShouldBindJSON(&req)
	c.ShouldBindUri(&req)

	var result *gin.H
	var err error

	if req.FromCache {
		result, err = nodefact_cached(req)
		if err != nil {
			common.HandleError(c, http.StatusExpectationFailed, "no cached facts.", err)
			return
		}
	} else {
		result, err = nodefacts(req)
		if err != nil {
			common.HandleError(c, http.StatusInternalServerError, "failed to get node facts.", err)
			return
		}
	}

	c.JSON(http.StatusOK, result)

}

func nodefacts(req GetNodeFactRequest) (*gin.H, error) {

	inventory := gin.H{
		"all": gin.H{
			"hosts": gin.H{
				req.Node: gin.H{
					"ansible_host":                 req.Ip,
					"ansible_port":                 req.Port,
					"ansible_user":                 req.User,
					"ansible_password":             req.Password,
					"ansible_ssh_private_key_file": req.PrivateKeyFile,
					"ansible_become":               req.Become,
					"ansible_become_user":          req.BecomeUser,
					"ansible_become_password":      req.BecomePassword,
					"ansible_ssh_common_args":      "-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=5 -o ConnectionAttempts=1 " + req.SshCommonArgs,
					"ansible_ssh_pipelining":       true,
					"kuboardspray_cluster_dir":     constants.GET_DATA_DIR() + "/" + req.NodeOwnerType + "/" + req.NodeOwner,
				},
			},
		},
	}

	inventoryBytes, err := json.Marshal(inventory)
	if err != nil {
		return nil, errors.New("failed to marshal inventory file: " + err.Error())
	}

	typeDir := constants.GET_DATA_DIR() + "/" + req.NodeOwnerType
	common.CreateDirIfNotExists(typeDir)

	ownerDir := typeDir + "/" + req.NodeOwner
	common.CreateDirIfNotExists(ownerDir)

	factDir := ownerDir + "/fact"
	common.CreateDirIfNotExists(factDir)

	inventoryPath := ownerDir + "/" + req.Node + "_" + time.Now().Format("2006-01-02_15-04-05.999") + ".json"

	err = ioutil.WriteFile(inventoryPath, inventoryBytes, 0666)
	if err != nil {
		return nil, errors.New("failed to create inventory file " + inventoryPath + err.Error())
	}
	logrus.Trace(string(inventoryBytes))

	defer os.Remove(inventoryPath)

	args := []string{req.Node, "-m", "setup", "-i", inventoryPath}
	if req.GatherSubset != "" {
		args = append(args, "-a", "gather_subset="+req.GatherSubset)
	}
	if req.Filter != "" {
		args = append(args, "-a", "filter="+req.Filter)
	}

	run := command.Run{
		Cmd:  "ansible",
		Args: args,
		Env:  []string{"ANSIBLE_CONFIG=" + constants.GET_ADHOC_CFG_PATH()},
		// Timeout: 5,
		// Dir:  dir + "/ansible-script",
	}

	stdout, _, err := run.Run()

	if strings.Contains(string(stdout), req.Node+" | ") {
		stdout = stdout[strings.Index(string(stdout), "{"):]
	} else if len(stdout) > 0 {
		return nil, errors.New(string(stdout))
	}
	if err != nil {
		return nil, errors.New("failed to get node facts " + err.Error())
	}

	factPath := factDir + "/" + req.Node + "_" + req.Ip + "_" + req.Port

	ioutil.WriteFile(factPath, stdout, 0666)

	fact := &gin.H{}
	if err := json.Unmarshal(stdout, fact); err != nil {
		return nil, errors.New("Failed to Unmarshal result " + err.Error())
	}
	return fact, nil
}
