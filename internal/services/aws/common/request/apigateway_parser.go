package request

import (
	"net/http"
	"strings"
)

func extractApiGatewayOperation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method

	if strings.HasPrefix(path, "/restapis") {
		return extractRestApiOperation(path, method)
	}
	if strings.HasPrefix(path, "/apikeys") {
		return extractApiKeyOperation(path, method)
	}
	if strings.HasPrefix(path, "/usageplans") {
		return extractUsagePlanOperation(path, method)
	}
	if strings.HasPrefix(path, "/domainnames") {
		return extractDomainNameOperation(path, method)
	}
	if strings.HasPrefix(path, "/tags/") {
		switch method {
		case "POST", "PUT":
			return "TagResource"
		case "DELETE":
			return "UntagResource"
		case "GET":
			return "ListTagsForResource"
		}
	}

	return ""
}

func extractRestApiOperation(path, method string) string {
	path = strings.TrimPrefix(path, "/restapis")

	if path == "" || path == "/" {
		switch method {
		case "GET":
			return "GetRestApis"
		case "POST":
			return "CreateRestApi"
		}
		return ""
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}

	apiId := parts[0]
	_ = apiId

	if len(parts) == 1 {
		switch method {
		case "GET":
			return "GetRestApi"
		case "DELETE":
			return "DeleteRestApi"
		case "PATCH":
			return "UpdateRestApi"
		}
		return ""
	}

	subResource := parts[1]
	switch subResource {
	case "resources":
		return extractResourceOperation(parts[2:], method)
	case "deployments":
		return extractDeploymentOperation(parts[2:], method)
	case "stages":
		return extractStageOperation(parts[2:], method)
	case "requestvalidators":
		return extractRequestValidatorOperation(parts[2:], method)
	case "models":
		return extractModelOperation(parts[2:], method)
	case "authorizers":
		return extractAuthorizerOperation(parts[2:], method)
	case "tags":
		if method == "GET" {
			return "ListTagsForResource"
		}
	}

	return ""
}

func extractResourceOperation(parts []string, method string) string {
	if len(parts) == 0 || parts[0] == "" {
		switch method {
		case "GET":
			return "GetResources"
		case "POST":
			return "CreateResource"
		}
		return ""
	}

	if len(parts) == 1 {
		switch method {
		case "GET":
			return "GetResource"
		case "POST":
			return "CreateResource"
		case "DELETE":
			return "DeleteResource"
		case "PATCH":
			return "UpdateResource"
		}
		return ""
	}

	if parts[1] == "methods" && len(parts) >= 3 {
		if len(parts) >= 4 && parts[3] == "integration" {
			if len(parts) >= 6 && parts[4] == "responses" {
				switch method {
				case "GET":
					return "GetIntegrationResponse"
				case "PUT":
					return "PutIntegrationResponse"
				case "DELETE":
					return "DeleteIntegrationResponse"
				}
				return ""
			}
			switch method {
			case "GET":
				return "GetIntegration"
			case "PUT":
				return "PutIntegration"
			case "DELETE":
				return "DeleteIntegration"
			case "PATCH":
				return "UpdateIntegration"
			}
			return ""
		}
		if len(parts) >= 5 && parts[3] == "responses" {
			switch method {
			case "GET":
				return "GetMethodResponse"
			case "PUT":
				return "PutMethodResponse"
			case "DELETE":
				return "DeleteMethodResponse"
			}
			return ""
		}
		switch method {
		case "GET":
			return "GetMethod"
		case "PUT":
			return "PutMethod"
		case "DELETE":
			return "DeleteMethod"
		case "PATCH":
			return "UpdateMethod"
		}
	}

	return ""
}

func extractDeploymentOperation(parts []string, method string) string {
	if len(parts) == 0 || parts[0] == "" {
		switch method {
		case "GET":
			return "GetDeployments"
		case "POST":
			return "CreateDeployment"
		}
		return ""
	}

	switch method {
	case "GET":
		return "GetDeployment"
	case "DELETE":
		return "DeleteDeployment"
	}

	return ""
}

func extractStageOperation(parts []string, method string) string {
	if len(parts) == 0 || parts[0] == "" {
		switch method {
		case "GET":
			return "GetStages"
		case "POST":
			return "CreateStage"
		}
		return ""
	}

	switch method {
	case "GET":
		return "GetStage"
	case "DELETE":
		return "DeleteStage"
	case "PATCH":
		return "UpdateStage"
	}

	return ""
}

func extractRequestValidatorOperation(parts []string, method string) string {
	if len(parts) == 0 || parts[0] == "" {
		switch method {
		case "GET":
			return "GetRequestValidators"
		case "POST":
			return "CreateRequestValidator"
		}
		return ""
	}

	switch method {
	case "GET":
		return "GetRequestValidator"
	case "DELETE":
		return "DeleteRequestValidator"
	case "PATCH":
		return "UpdateRequestValidator"
	}

	return ""
}

func extractModelOperation(parts []string, method string) string {
	if len(parts) == 0 || parts[0] == "" {
		switch method {
		case "GET":
			return "GetModels"
		case "POST":
			return "CreateModel"
		}
		return ""
	}

	switch method {
	case "GET":
		return "GetModel"
	case "DELETE":
		return "DeleteModel"
	case "PATCH":
		return "UpdateModel"
	}

	return ""
}

func extractAuthorizerOperation(parts []string, method string) string {
	if len(parts) == 0 || parts[0] == "" {
		switch method {
		case "GET":
			return "GetAuthorizers"
		case "POST":
			return "CreateAuthorizer"
		}
		return ""
	}

	switch method {
	case "GET":
		return "GetAuthorizer"
	case "DELETE":
		return "DeleteAuthorizer"
	case "PATCH":
		return "UpdateAuthorizer"
	}

	return ""
}

func extractApiKeyOperation(path, method string) string {
	path = strings.TrimPrefix(path, "/apikeys")

	if path == "" || path == "/" {
		switch method {
		case "GET":
			return "GetApiKeys"
		case "POST":
			return "CreateApiKey"
		}
		return ""
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}

	switch method {
	case "GET":
		return "GetApiKey"
	case "DELETE":
		return "DeleteApiKey"
	case "PATCH":
		return "UpdateApiKey"
	}

	return ""
}

func extractUsagePlanOperation(path, method string) string {
	path = strings.TrimPrefix(path, "/usageplans")

	if path == "" || path == "/" {
		switch method {
		case "GET":
			return "GetUsagePlans"
		case "POST":
			return "CreateUsagePlan"
		}
		return ""
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}

	if len(parts) == 1 {
		switch method {
		case "GET":
			return "GetUsagePlan"
		case "DELETE":
			return "DeleteUsagePlan"
		case "PATCH":
			return "UpdateUsagePlan"
		}
		return ""
	}

	if parts[1] == "keys" && len(parts) >= 3 {
		switch method {
		case "GET":
			return "GetUsagePlanKey"
		case "POST":
			return "CreateUsagePlanKey"
		case "DELETE":
			return "DeleteUsagePlanKey"
		}
	}
	if parts[1] == "keys" && len(parts) == 2 {
		switch method {
		case "GET":
			return "GetUsagePlanKeys"
		case "POST":
			return "CreateUsagePlanKey"
		}
	}

	return ""
}

func extractDomainNameOperation(path, method string) string {
	path = strings.TrimPrefix(path, "/domainnames")

	if path == "" || path == "/" {
		switch method {
		case "GET":
			return "GetDomainNames"
		case "POST":
			return "CreateDomainName"
		}
		return ""
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}

	if len(parts) == 1 {
		switch method {
		case "GET":
			return "GetDomainName"
		case "DELETE":
			return "DeleteDomainName"
		case "PATCH":
			return "UpdateDomainName"
		}
		return ""
	}

	if parts[1] == "basepathmappings" && len(parts) >= 3 {
		switch method {
		case "GET":
			return "GetBasePathMapping"
		case "DELETE":
			return "DeleteBasePathMapping"
		case "PATCH":
			return "UpdateBasePathMapping"
		}
	}
	if parts[1] == "basepathmappings" && len(parts) == 2 {
		switch method {
		case "GET":
			return "GetBasePathMappings"
		case "POST":
			return "CreateBasePathMapping"
		}
	}

	return ""
}

func isApiGatewayPath(path string) bool {
	return strings.HasPrefix(path, "/restapis") ||
		strings.HasPrefix(path, "/apikeys") ||
		strings.HasPrefix(path, "/usageplans") ||
		strings.HasPrefix(path, "/domainnames") ||
		strings.HasPrefix(path, "/vpclinks") ||
		strings.HasPrefix(path, "/apis") ||
		strings.HasPrefix(path, "/authorizers") ||
		strings.HasPrefix(path, "/tags/")
}

func extractApiGatewayPathParams(path string, params map[string]interface{}) {
	if strings.HasPrefix(path, "/restapis/") {
		extractRestApiPathParams(path, params)
	} else if strings.HasPrefix(path, "/apikeys/") {
		extractApiKeyPathParams(path, params)
	} else if strings.HasPrefix(path, "/usageplans/") {
		extractUsagePlanPathParams(path, params)
	} else if strings.HasPrefix(path, "/domainnames/") {
		extractDomainNamePathParams(path, params)
	} else if strings.HasPrefix(path, "/tags/") {
		resourceArn := strings.TrimPrefix(path, "/tags/")
		if _, ok := params["resourceArn"]; !ok {
			params["resourceArn"] = resourceArn
		}
	}
}

func extractRestApiPathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/restapis/")
	parts := strings.Split(path, "/")

	if len(parts) >= 1 && parts[0] != "" {
		if _, ok := params["restApiId"]; !ok {
			params["restApiId"] = parts[0]
		}
	}

	if len(parts) >= 3 && parts[1] == "resources" {
		if _, ok := params["resourceId"]; !ok {
			params["resourceId"] = parts[2]
		}
		if len(parts) >= 5 && parts[3] == "methods" {
			if len(parts) >= 6 && parts[5] == "integration" {
				if existingMethod, ok := params["httpMethod"]; ok {
					if _, hasIntegrationMethod := params["integrationHttpMethod"]; !hasIntegrationMethod {
						params["integrationHttpMethod"] = existingMethod
					}
				}
				params["httpMethod"] = parts[4]
			} else {
				if _, ok := params["httpMethod"]; !ok {
					params["httpMethod"] = parts[4]
				}
			}
			if len(parts) >= 7 && parts[5] == "responses" {
				if _, ok := params["statusCode"]; !ok {
					params["statusCode"] = parts[6]
				}
			}
			if len(parts) >= 8 && parts[5] == "integration" && parts[6] == "responses" {
				if _, ok := params["statusCode"]; !ok {
					params["statusCode"] = parts[7]
				}
			}
		}
	}

	if len(parts) >= 3 && parts[1] == "deployments" {
		if _, ok := params["deploymentId"]; !ok {
			params["deploymentId"] = parts[2]
		}
	}

	if len(parts) >= 3 && parts[1] == "stages" {
		if _, ok := params["stageName"]; !ok {
			params["stageName"] = parts[2]
		}
	}

	if len(parts) >= 3 && parts[1] == "requestvalidators" {
		if _, ok := params["requestValidatorId"]; !ok {
			params["requestValidatorId"] = parts[2]
		}
	}

	if len(parts) >= 3 && parts[1] == "models" {
		if _, ok := params["modelName"]; !ok {
			params["modelName"] = parts[2]
		}
	}

	if len(parts) >= 3 && parts[1] == "authorizers" {
		if _, ok := params["authorizerId"]; !ok {
			params["authorizerId"] = parts[2]
		}
	}
}

func extractApiKeyPathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/apikeys/")
	if path != "" && !strings.Contains(path, "/") {
		if _, ok := params["apiKey"]; !ok {
			params["apiKey"] = path
		}
	}
}

func extractUsagePlanPathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/usageplans/")
	parts := strings.Split(path, "/")

	if len(parts) >= 1 && parts[0] != "" {
		if _, ok := params["usagePlanId"]; !ok {
			params["usagePlanId"] = parts[0]
		}
	}

	if len(parts) >= 3 && parts[1] == "keys" {
		if _, ok := params["keyId"]; !ok {
			params["keyId"] = parts[2]
		}
	}
}

func extractDomainNamePathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/domainnames/")
	parts := strings.Split(path, "/")

	if len(parts) >= 1 && parts[0] != "" {
		if _, ok := params["domainName"]; !ok {
			params["domainName"] = parts[0]
		}
	}

	if len(parts) >= 3 && parts[1] == "basepathmappings" {
		if _, ok := params["basePath"]; !ok {
			params["basePath"] = parts[2]
		}
	}
}
