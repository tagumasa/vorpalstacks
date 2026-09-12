package cognitoidentityprovider

import (
	"strings"

	"vorpalstacks/internal/common/request"
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

// The service protocol is awsJson1_1: request members are PascalCase JSON
// keys carrying typed JSON values. The readers below accept exactly that
// shape — a query-string-style member (dotted keys, string-typed numbers or
// booleans, lower-cased names) is not a second supported format; it reads
// as absent and the required-member validation downstream rejects it.

func getBoolParam(req *request.ParsedRequest, key string) bool {
	if v, ok := req.Parameters[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// getBoolParamOK reports whether the boolean member is present in the
// request, and its value. Presence, not the value, decides whether the
// member is applied: an omitted boolean member must leave the stored
// configuration untouched, so callers that cannot distinguish false from
// absent silently flip stored flags on every update.
func getBoolParamOK(req *request.ParsedRequest, key string) (bool, bool) {
	if v, ok := req.Parameters[key]; ok {
		if b, ok := v.(bool); ok {
			return b, true
		}
	}
	return false, false
}

func getIntParam(req *request.ParsedRequest, key string) int {
	v, _ := parseIntParam(req, key)
	return v
}

func getIntParamOK(req *request.ParsedRequest, key string) (int, bool) {
	return parseIntParam(req, key)
}

func parseIntParam(req *request.ParsedRequest, key string) (int, bool) {
	if v, ok := req.Parameters[key]; ok {
		return parseJSONInt(v)
	}
	return 0, false
}

// parseJSONInt reads a JSON number at its decoded Go type. awsJson1_1
// carries integers as JSON numbers; a value of any other type is a
// wire-shape violation, not a number in disguise.
func parseJSONInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	}
	return 0, false
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

	attrs, ok := req.Parameters[key].([]interface{})
	if !ok {
		return attributes
	}
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

// parseClientMetadata extracts ClientMetadata from the request.
func parseClientMetadata(req *request.ParsedRequest) map[string]string {
	metadata := make(map[string]string)

	cm, ok := req.Parameters["ClientMetadata"].(map[string]interface{})
	if !ok {
		return metadata
	}
	for k, v := range cm {
		if vs, ok := v.(string); ok {
			metadata[k] = vs
		}
	}
	return metadata
}

// getStringParam reads a string member from a decoded JSON object; a member
// of any other type reads as absent (empty string). It is the package-wide
// reader for nested wire objects such as provider details and asset maps.
func getStringParam(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
