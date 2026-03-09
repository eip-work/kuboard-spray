package operation_v2

import "strconv"

func appendCommonParams(result []string, req ExecuteStepRequest, skipLimitParam bool) []string {
	result = append(result, "--forks", strconv.Itoa(req.Fork))
	result = append(result, "-e", "pangeecluster_operation="+req.Operation)
	if req.ExcludeNodes != "" && !skipLimitParam {
		result = append(result, "--limit", req.ExcludeNodes)
	}
	if req.Verbose == "vvv" {
		result = append(result, "-vvv")
	}
	if req.Verbose == "v" {
		result = append(result, "-v")
	}
	return result
}
