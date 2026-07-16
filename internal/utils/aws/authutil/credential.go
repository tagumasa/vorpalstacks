package authutil

import "strings"

// AWSCredential represents the fields extracted from an AWS SigV4
// Authorization header's Credential= component.
type AWSCredential struct {
	AccessKeyID    string
	Date           string
	Region         string
	SigningService string
	RequestType    string
}

// ParseAWSCredential extracts credential fields from an AWS Authorization
// header. Returns nil if the header is empty or does not contain a
// Credential= component with at least 5 slash-separated parts.
func ParseAWSCredential(authHeader string) *AWSCredential {
	if authHeader == "" || !strings.Contains(authHeader, "Credential=") {
		return nil
	}
	parts := strings.Split(authHeader, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "Credential="); idx >= 0 {
			cred := strings.TrimPrefix(part[idx:], "Credential=")
			credParts := strings.Split(cred, "/")
			if len(credParts) >= 5 {
				return &AWSCredential{
					AccessKeyID:    credParts[0],
					Date:           credParts[1],
					Region:         credParts[2],
					SigningService: credParts[3],
					RequestType:    credParts[4],
				}
			}
		}
	}
	return nil
}
