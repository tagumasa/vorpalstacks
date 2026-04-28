package request

import "strings"

// ExtractRegionFromAuth extracts the region from an AWS V4 auth header.
func ExtractRegionFromAuth(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	if !strings.HasPrefix(authHeader, "AWS4-HMAC-SHA256 ") {
		return ""
	}
	authHeader = strings.TrimPrefix(authHeader, "AWS4-HMAC-SHA256 ")
	parts := strings.Split(authHeader, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "Credential=") {
			credential := strings.TrimPrefix(part, "Credential=")
			credParts := strings.Split(credential, "/")
			if len(credParts) >= 3 {
				return credParts[2]
			}
		}
	}
	return ""
}

// ExtractAccessKeyIDFromAuth extracts the access key ID from an AWS V4 auth header.
func ExtractAccessKeyIDFromAuth(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	if !strings.HasPrefix(authHeader, "AWS4-HMAC-SHA256 ") {
		return ""
	}
	authHeader = strings.TrimPrefix(authHeader, "AWS4-HMAC-SHA256 ")
	parts := strings.Split(authHeader, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "Credential=") {
			credential := strings.TrimPrefix(part, "Credential=")
			credParts := strings.Split(credential, "/")
			if len(credParts) >= 1 {
				return credParts[0]
			}
		}
	}
	return ""
}

// ExtractAccountIDFromAuth extracts the account ID from an AWS V4 auth header.
// It is an alias for ExtractAccessKeyIDFromAuth.
var ExtractAccountIDFromAuth = ExtractAccessKeyIDFromAuth
