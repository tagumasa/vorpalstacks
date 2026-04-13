package request

import (
	"net/http"
	"strings"
)

func extractLambdaOperation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method

	if !strings.HasPrefix(path, "/20") {
		return ""
	}

	switch {
	case strings.HasPrefix(path, "/2015-03-31/functions"):
		return extractLambdaFunctionOperation(path, method)
	case strings.HasPrefix(path, "/2018-10-31/layers"):
		return extractLambdaLayerOperation(path, method)
	case strings.HasPrefix(path, "/2015-03-31/event-source-mappings"):
		return extractLambdaEventSourceOperation(path, method)
	case strings.HasPrefix(path, "/2017-03-31/tags/"):
		return extractLambdaTagOperation(path, method)
	case strings.HasPrefix(path, "/2017-10-31/functions/") && strings.HasSuffix(path, "/concurrency"):
		return extractLambdaConcurrencyOperation(path, method)
	case strings.HasPrefix(path, "/2019-09-30/functions/") && strings.HasSuffix(path, "/concurrency"):
		if method == "GET" {
			return "GetFunctionConcurrency"
		}
	case strings.HasPrefix(path, "/2019-09-30/functions/") && strings.Contains(path, "/provisioned-concurrency"):
		return extractLambdaProvisionedConcurrencyOperation(path, method, r)
	case strings.HasPrefix(path, "/2019-09-25/functions/") && strings.Contains(path, "/event-invoke-config"):
		return extractLambdaEventInvokeConfigOperation(path, method)
	case strings.HasPrefix(path, "/2014-11-13/functions/") && strings.HasSuffix(path, "/invoke-async"):
		if method == "POST" {
			return "InvokeAsync"
		}
	case strings.HasPrefix(path, "/2021-11-15/functions/") && strings.Contains(path, "/response-streaming-invocations"):
		if method == "POST" {
			return "InvokeWithResponseStream"
		}
	case strings.HasPrefix(path, "/2021-10-31/functions/"):
		return extractLambdaUrlOperation(path, method)
	case path == "/2016-08-19/account-settings" && method == "GET":
		return "GetAccountSettings"
	}

	return ""
}

func extractLambdaFunctionOperation(path, method string) string {
	path = strings.TrimPrefix(path, "/2015-03-31/functions")

	if path == "" || path == "/" {
		switch method {
		case "GET":
			return "ListFunctions"
		case "POST":
			return "CreateFunction"
		}
		return ""
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	if len(parts) == 1 {
		switch method {
		case "GET":
			return "GetFunction"
		case "DELETE":
			return "DeleteFunction"
		}
		return ""
	}

	subResource := parts[1]
	switch subResource {
	case "invocations":
		if method == "POST" {
			return "Invoke"
		}
	case "invoke-async":
		if method == "POST" {
			return "InvokeAsync"
		}
	case "configuration":
		switch method {
		case "GET":
			return "GetFunctionConfiguration"
		case "PUT":
			return "UpdateFunctionConfiguration"
		}
	case "code":
		if method == "PUT" {
			return "UpdateFunctionCode"
		}
	case "versions":
		switch method {
		case "GET":
			return "ListVersionsByFunction"
		case "POST":
			return "PublishVersion"
		}
	case "aliases":
		if len(parts) == 2 {
			switch method {
			case "GET":
				return "ListAliases"
			case "POST":
				return "CreateAlias"
			}
		} else if len(parts) >= 3 {
			switch method {
			case "GET":
				return "GetAlias"
			case "PUT":
				return "UpdateAlias"
			case "DELETE":
				return "DeleteAlias"
			}
		}
	case "policy":
		if len(parts) == 2 {
			switch method {
			case "GET":
				return "GetPolicy"
			case "POST":
				return "AddPermission"
			}
		} else if len(parts) >= 3 {
			if method == "DELETE" {
				return "RemovePermission"
			}
		}
	}

	return ""
}

func extractLambdaLayerOperation(path, method string) string {
	path = strings.TrimPrefix(path, "/2018-10-31/layers")

	if path == "" || path == "/" {
		if method == "GET" {
			return "ListLayers"
		}
		return ""
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return ""
	}

	if parts[1] == "versions" {
		if len(parts) == 2 {
			switch method {
			case "GET":
				return "ListLayerVersions"
			case "POST":
				return "PublishLayerVersion"
			}
		} else if len(parts) >= 3 {
			switch method {
			case "GET":
				return "GetLayerVersion"
			case "DELETE":
				return "DeleteLayerVersion"
			}
		}
	}

	return ""
}

func extractLambdaEventSourceOperation(path, method string) string {
	path = strings.TrimPrefix(path, "/2015-03-31/event-source-mappings")

	if path == "" || path == "/" {
		switch method {
		case "GET":
			return "ListEventSourceMappings"
		case "POST":
			return "CreateEventSourceMapping"
		}
		return ""
	}

	switch method {
	case "GET":
		return "GetEventSourceMapping"
	case "PUT":
		return "UpdateEventSourceMapping"
	case "DELETE":
		return "DeleteEventSourceMapping"
	}

	return ""
}

func extractLambdaTagOperation(path, method string) string {
	switch method {
	case "POST":
		return "TagResource"
	case "GET":
		return "ListTags"
	case "DELETE":
		return "UntagResource"
	}
	return ""
}

func extractLambdaConcurrencyOperation(path, method string) string {
	switch method {
	case "PUT":
		return "PutFunctionConcurrency"
	case "DELETE":
		return "DeleteFunctionConcurrency"
	}
	return ""
}

func extractLambdaUrlOperation(path, method string) string {
	if !strings.Contains(path, "/url") {
		return ""
	}

	if strings.HasSuffix(path, "/urls") && method == "GET" {
		return "ListFunctionUrlConfigs"
	}

	if strings.Contains(path, "/url") {
		switch method {
		case "POST":
			return "CreateFunctionUrlConfig"
		case "GET":
			return "GetFunctionUrlConfig"
		case "PUT":
			return "UpdateFunctionUrlConfig"
		case "DELETE":
			return "DeleteFunctionUrlConfig"
		}
	}

	return ""
}

func extractLambdaProvisionedConcurrencyOperation(path string, method string, r *http.Request) string {
	path = strings.TrimPrefix(path, "/2019-09-30/functions/")

	if strings.HasSuffix(path, "/provisioned-concurrency") {
		if method == "GET" && r.URL.Query().Get("List") != "" {
			return "ListProvisionedConcurrencyConfigs"
		}
		switch method {
		case "PUT":
			return "PutProvisionedConcurrencyConfig"
		case "GET":
			return "GetProvisionedConcurrencyConfig"
		case "DELETE":
			return "DeleteProvisionedConcurrencyConfig"
		}
	}

	return ""
}

func extractLambdaEventInvokeConfigOperation(path, method string) string {
	path = strings.TrimPrefix(path, "/2019-09-25/functions/")

	if strings.HasSuffix(path, "/event-invoke-config/list") {
		if method == "GET" {
			return "ListFunctionEventInvokeConfigs"
		}
		return ""
	}

	if strings.HasSuffix(path, "/event-invoke-config") {
		switch method {
		case "PUT":
			return "PutFunctionEventInvokeConfig"
		case "GET":
			return "GetFunctionEventInvokeConfig"
		case "DELETE":
			return "DeleteFunctionEventInvokeConfig"
		}
	}

	return ""
}

func extractLambdaPathParams(path string, params map[string]interface{}) {
	switch {
	case strings.HasPrefix(path, "/2015-03-31/functions/"):
		extractFunctionPathParams(path, params)
	case strings.HasPrefix(path, "/2018-10-31/layers/"):
		extractLayerPathParams(path, params)
	case strings.HasPrefix(path, "/2015-03-31/event-source-mappings/"):
		extractEventSourcePathParams(path, params)
	case strings.HasPrefix(path, "/2017-03-31/tags/"):
		extractTagsPathParams(path, params)
	case strings.HasPrefix(path, "/2017-10-31/functions/"):
		extractConcurrencyPathParams(path, params)
	case strings.HasPrefix(path, "/2019-09-30/functions/"):
		if strings.Contains(path, "/provisioned-concurrency") {
			extractProvisionedConcurrencyPathParams(path, params)
		} else {
			extractConcurrencyPathParams(path, params)
		}
	case strings.HasPrefix(path, "/2019-09-25/functions/"):
		extractEventInvokeConfigPathParams(path, params)
	case strings.HasPrefix(path, "/2021-10-31/functions/"):
		extractUrlPathParams(path, params)
	case strings.HasPrefix(path, "/2014-11-13/functions/"):
		extractInvokeAsyncPathParams(path, params)
	case strings.HasPrefix(path, "/2021-11-15/functions/"):
		extractResponseStreamPathParams(path, params)
	}
}

func extractFunctionPathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/2015-03-31/functions/")
	parts := strings.Split(path, "/")

	if len(parts) >= 1 && parts[0] != "" {
		if _, ok := params["FunctionName"]; !ok {
			params["FunctionName"] = parts[0]
		}
	}

	if len(parts) >= 3 {
		switch parts[1] {
		case "aliases":
			if _, ok := params["Name"]; !ok {
				params["Name"] = parts[2]
			}
		case "policy":
			if _, ok := params["StatementId"]; !ok {
				params["StatementId"] = parts[2]
			}
		}
	}
}

func extractLayerPathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/2018-10-31/layers/")
	parts := strings.Split(path, "/")

	if len(parts) >= 1 && parts[0] != "" {
		if _, ok := params["LayerName"]; !ok {
			params["LayerName"] = parts[0]
		}
	}

	if len(parts) >= 3 && parts[1] == "versions" {
		if _, ok := params["VersionNumber"]; !ok {
			params["VersionNumber"] = parts[2]
		}
	}
}

func extractEventSourcePathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/2015-03-31/event-source-mappings/")
	if path != "" && !strings.Contains(path, "/") {
		if _, ok := params["UUID"]; !ok {
			params["UUID"] = path
		}
	}
}

func extractTagsPathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/2017-03-31/tags/")
	if path != "" {
		if _, ok := params["Resource"]; !ok {
			params["Resource"] = path
		}
	}
}

func extractConcurrencyPathParams(path string, params map[string]interface{}) {
	re := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/concurrency")
	parts := strings.Split(re, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && !isLambdaDatePath(parts[i]) {
			if _, ok := params["FunctionName"]; !ok {
				params["FunctionName"] = parts[i]
			}
			break
		}
	}
}

func extractUrlPathParams(path string, params map[string]interface{}) {
	re := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/url")
	re = strings.TrimSuffix(re, "/urls")
	parts := strings.Split(re, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && !isLambdaDatePath(parts[i]) {
			if _, ok := params["FunctionName"]; !ok {
				params["FunctionName"] = parts[i]
			}
			break
		}
	}
}

func isLambdaDatePath(s string) bool {
	if len(s) != 10 {
		return false
	}
	if !strings.HasPrefix(s, "20") {
		return false
	}
	return strings.Count(s, "-") == 2
}

func extractProvisionedConcurrencyPathParams(path string, params map[string]interface{}) {
	re := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/provisioned-concurrency")
	parts := strings.Split(re, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && !isLambdaDatePath(parts[i]) {
			if _, ok := params["FunctionName"]; !ok {
				params["FunctionName"] = parts[i]
			}
			break
		}
	}
}

func extractEventInvokeConfigPathParams(path string, params map[string]interface{}) {
	re := strings.TrimSuffix(path, "/event-invoke-config")
	re = strings.TrimSuffix(re, "/event-invoke-config/list")
	parts := strings.Split(strings.Trim(re, "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && !isLambdaDatePath(parts[i]) {
			if _, ok := params["FunctionName"]; !ok {
				params["FunctionName"] = parts[i]
			}
			break
		}
	}
}

func extractInvokeAsyncPathParams(path string, params map[string]interface{}) {
	re := strings.TrimSuffix(path, "/invoke-async")
	parts := strings.Split(strings.Trim(re, "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && !isLambdaDatePath(parts[i]) {
			if _, ok := params["FunctionName"]; !ok {
				params["FunctionName"] = parts[i]
			}
			break
		}
	}
}

func extractResponseStreamPathParams(path string, params map[string]interface{}) {
	re := strings.TrimSuffix(path, "/response-streaming-invocations")
	parts := strings.Split(strings.Trim(re, "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && !isLambdaDatePath(parts[i]) {
			if _, ok := params["FunctionName"]; !ok {
				params["FunctionName"] = parts[i]
			}
			break
		}
	}
}

func extractLambdaHeaders(r *http.Request, params map[string]interface{}) {
	headerMappings := map[string]string{
		"X-Amz-Invocation-Type": "InvocationType",
		"X-Amz-Log-Type":        "LogType",
		"X-Amz-Client-Context":  "ClientContext",
		"X-Amz-Function-Error":  "FunctionError",
		"X-Amz-Qualifier":       "Qualifier",
	}
	for header, param := range headerMappings {
		if v := r.Header.Get(header); v != "" {
			params[param] = v
		}
	}
}
