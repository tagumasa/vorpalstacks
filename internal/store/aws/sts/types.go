package sts

// Package sts provides STS (Security Token Service) data store implementations
// for vorpalstacks.

import (
	"time"
)

// Session represents an STS temporary session with credentials.
type Session struct { // #nosec G117
	SessionToken    string    `json:"session_token"`
	AccessKeyId     string    `json:"access_key_id"`
	SecretAccessKey string    `json:"secret_access_key"`
	Expiration      time.Time `json:"expiration"`
	PrincipalArn    string    `json:"principal_arn"`
	PrincipalType   string    `json:"principal_type"`
	PrincipalName   string    `json:"principal_name"`
	RoleArn         string    `json:"role_arn,omitempty"`
	RoleSessionName string    `json:"role_session_name,omitempty"`
}
