package cognitoidentityprovider

import (
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"

	"github.com/google/uuid"
)

// cognitoIdpHost returns the Cognito IDP hostname for the given region,
// accounting for partition-specific suffixes (aws-cn uses .amazonaws.com.cn).
func cognitoIdpHost(region string) string {
	if strings.HasPrefix(region, "cn-") {
		return "cognito-idp." + region + ".amazonaws.com.cn"
	}
	return "cognito-idp." + region + ".amazonaws.com"
}

// cognitoImportHost returns the hostname of the user-import CSV upload
// endpoint, accounting for partition-specific suffixes (aws-cn uses
// .amazonaws.com.cn).
func cognitoImportHost(region string) string {
	if strings.HasPrefix(region, "cn-") {
		return "cognito-import." + region + ".amazonaws.com.cn"
	}
	return "cognito-import." + region + ".amazonaws.com"
}

func getBoolParam(req *request.ParsedRequest, key string) bool {
	lowerKey := strings.ToLower(key[:1]) + key[1:]

	if v, ok := req.Parameters[lowerKey]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	if v, ok := req.Parameters[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}

	val := request.GetParamLowerFirst(req.Parameters, key)
	return val == "true" || val == "True" || val == "1"
}

func getIntParam(req *request.ParsedRequest, key string) int {
	v, _ := parseIntParam(req, key)
	return v
}

func getIntParamOK(req *request.ParsedRequest, key string) (int, bool) {
	return parseIntParam(req, key)
}

func parseIntParam(req *request.ParsedRequest, key string) (int, bool) {
	tryKey := func(k string) (int, bool) {
		if v, ok := req.Parameters[k]; ok {
			switch n := v.(type) {
			case int:
				return n, true
			case int64:
				return int(n), true
			case float64:
				return int(n), true
			case string:
				if n != "" {
					return parseInt(n), true
				}
			}
		}
		return 0, false
	}
	if v, ok := tryKey(key); ok {
		return v, true
	}
	lowerKey := strings.ToLower(key[:1]) + key[1:]
	return tryKey(lowerKey)
}

func getUserPoolID(req *request.ParsedRequest) string {
	return req.GetParam("UserPoolId")
}

func getUsername(req *request.ParsedRequest) string {
	return req.GetParam("Username")
}

func getGroupName(req *request.ParsedRequest) string {
	return req.GetParam("GroupName")
}

func getPassword(req *request.ParsedRequest) string {
	return req.GetParam("Password")
}

func getNewPassword(req *request.ParsedRequest) string {
	if v := req.GetParam("NewPassword"); v != "" {
		return v
	}
	return req.GetParam("ProposedPassword")
}

func getPreviousPassword(req *request.ParsedRequest) string {
	return req.GetParam("PreviousPassword")
}

func getAccessToken(req *request.ParsedRequest) string {
	return req.GetParam("AccessToken")
}

func getConfirmationCode(req *request.ParsedRequest) string {
	return req.GetParam("ConfirmationCode")
}

func getClientId(req *request.ParsedRequest) string {
	return req.GetParam("ClientId")
}

func parseUserAttributes(req *request.ParsedRequest) map[string]string {
	return parseNamedAttributeList(req, "UserAttributes")
}

// parseValidationData extracts ValidationData from the request.
func parseValidationData(req *request.ParsedRequest) map[string]string {
	return parseNamedAttributeList(req, "ValidationData")
}

func parseNamedAttributeList(req *request.ParsedRequest, key string) map[string]string {
	attributes := make(map[string]string)

	if attrs, ok := req.Parameters[key].([]interface{}); ok {
		for _, attr := range attrs {
			if m, ok := attr.(map[string]interface{}); ok {
				name, _ := m["Name"].(string)
				value, _ := m["Value"].(string)
				if name != "" {
					attributes[name] = value
				}
			}
		}
		return attributes
	}

	for i := 1; ; i++ {
		idx := strconv.Itoa(i)
		nameKey := key + "." + idx + ".Name"
		valueKey := key + "." + idx + ".Value"
		name := req.GetParam(nameKey)
		if name == "" {
			break
		}
		value := req.GetParam(valueKey)
		attributes[name] = value
	}
	return attributes
}

// parseClientMetadata extracts ClientMetadata from the request.
func parseClientMetadata(req *request.ParsedRequest) map[string]string {
	metadata := make(map[string]string)
	if cm, ok := req.Parameters["ClientMetadata"].(map[string]interface{}); ok {
		for k, v := range cm {
			if vs, ok := v.(string); ok {
				metadata[k] = vs
			}
		}
	}
	for k := range req.Parameters {
		if strings.HasPrefix(k, "ClientMetadata.") {
			attrKey := strings.TrimPrefix(k, "ClientMetadata.")
			metadata[attrKey] = req.GetParam(k)
		}
	}
	return metadata
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var result int
	for _, c := range s {
		if c < '0' || c > '9' {
			return -1
		}
		result = result*10 + int(c-'0')
	}
	return result
}

func generateSessionID() string {
	return "SESSION_" + uuid.New().String()
}
