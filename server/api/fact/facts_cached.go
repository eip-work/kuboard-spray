package fact

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/constants"
)

func nodefact_cached(req GetNodeFactRequest) (*ansible_rpc.AnsibleResultNode, error) {

	factDir := constants.GET_DATA_DIR() + "/" + req.NodeOwnerType + "/" + req.NodeOwner + "/fact"
	factPath := factDir + "/" + req.Node + "_" + req.Ip + "_" + req.Port

	fact_bytes, err := ioutil.ReadFile(factPath)
	if err != nil {
		return nil, err
	}

	fact := &ansible_rpc.AnsibleResultNode{}
	if err := json.Unmarshal(fact_bytes, fact); err != nil {
		return nil, errors.New("Failed to Unmarshal result " + err.Error())
	}
	return fact, nil

}
