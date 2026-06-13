package request

import (
	"net/http"
	"strings"
)

type iotEventsRESTParser struct{}

func (p *iotEventsRESTParser) MatchPath(path string) bool {
	return strings.HasPrefix(path, "/detector-models") ||
		strings.HasPrefix(path, "/inputs") ||
		strings.HasPrefix(path, "/tags")
}

func (p *iotEventsRESTParser) ExtractOperation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")

	switch {
	case path == "/detector-models" && method == http.MethodPost:
		return "CreateDetectorModel"
	case path == "/detector-models" && method == http.MethodGet:
		return "ListDetectorModels"
	case strings.HasPrefix(path, "/detector-models/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "UpdateDetectorModel"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeDetectorModel"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteDetectorModel"
		}
	case path == "/inputs" && method == http.MethodPost:
		return "CreateInput"
	case path == "/inputs" && method == http.MethodGet:
		return "ListInputs"
	case strings.HasPrefix(path, "/inputs/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPut:
			return "UpdateInput"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeInput"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteInput"
		}
	case path == "/tags" && method == http.MethodPost:
		return "TagResource"
	case path == "/tags" && method == http.MethodDelete:
		return "UntagResource"
	case path == "/tags" && method == http.MethodGet:
		return "ListTagsForResource"
	}
	return ""
}

func (p *iotEventsRESTParser) ExtractPathParams(r *http.Request, params map[string]interface{}) {
	path := r.URL.Path
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")

	switch {
	case strings.HasPrefix(path, "/detector-models/") && len(parts) >= 2:
		params["detectorModelName"] = parts[1]
	case strings.HasPrefix(path, "/inputs/") && len(parts) >= 2:
		params["inputName"] = parts[1]
	}
}
