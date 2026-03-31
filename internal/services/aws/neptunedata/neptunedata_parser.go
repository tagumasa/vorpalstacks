package neptunedata

import (
	"net/url"

	"vorpalstacks/internal/services/aws/common/request"
)

func getPathParam(req *request.ParsedRequest, key string) string {
	raw := request.GetStringParam(req.Parameters, key)
	if raw == "" {
		return ""
	}
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return raw
	}
	return decoded
}
