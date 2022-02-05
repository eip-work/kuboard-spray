package fact

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

type GetNodeFactRequest struct {
	ansible_rpc.AdhocCommandRequestWithIP
	Node         string `uri:"node"`
	FromCache    bool   `json:"from_cache"`
	GatherSubset string `json:"gather_subset"`
	Filter       string `json:"filter"`
}

func GetNodeFacts(c *gin.Context) {

	var req GetNodeFactRequest
	c.ShouldBindJSON(&req)
	c.ShouldBindUri(&req)

	var result *ansible_rpc.AnsibleResultNode
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

func nodefacts(req GetNodeFactRequest) (*ansible_rpc.AnsibleResultNode, error) {

	typeDir := constants.GET_DATA_DIR() + "/" + req.NodeOwnerType
	common.CreateDirIfNotExists(typeDir)

	ownerDir := typeDir + "/" + req.NodeOwner
	common.CreateDirIfNotExists(ownerDir)

	factDir := ownerDir + "/fact"
	common.CreateDirIfNotExists(factDir)

	args := []string{"-m", "setup"}
	if req.GatherSubset != "" {
		args = append(args, "-a", "gather_subset="+req.GatherSubset)
	}
	if req.Filter != "" {
		args = append(args, "-a", "filter="+req.Filter)
	}
	fact, err := ansible_rpc.ExecuteAdhocCommandWithIp(req.AdhocCommandRequestWithIP, args)
	if err != nil {
		return nil, err
	}

	stdout, _ := json.Marshal(fact[0])
	factPath := factDir + "/" + req.Node + "_" + req.Ip + "_" + req.Port
	ioutil.WriteFile(factPath, stdout, 0666)

	return &fact[0], nil
}
