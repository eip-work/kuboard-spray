package fact

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

func nodefact_cached(req GetNodeFactRequest) (*gin.H, error) {

	factDir := constants.GET_DATA_DIR() + "/fact/" + req.NodeType + "/" + req.NodeOwner
	factPath := factDir + "/" + req.Node + "_" + req.Ip + "_" + req.Port

	fact_bytes, err := ioutil.ReadFile(factPath)
	if err != nil {
		return nil, err
	}

	fact := &gin.H{}
	if err := json.Unmarshal(fact_bytes, fact); err != nil {
		return nil, errors.New("Failed to Unmarshal result " + err.Error())
	}
	return fact, nil

}
