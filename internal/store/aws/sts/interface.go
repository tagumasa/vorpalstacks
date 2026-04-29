package sts

import "vorpalstacks/internal/common/auth"

// SessionStoreInterface defines operations for managing STS sessions.
type SessionStoreInterface interface {
	auth.SessionResolver
	Create(principalType, principalName, principalArn, roleArn, roleSessionName string, durationSeconds int) (*Session, error)
	Get(sessionToken string) (*Session, error)
	GetByAccessKeyId(accessKeyId string) (*Session, error)
	Delete(sessionToken string) error
}

var _ SessionStoreInterface = (*SessionStore)(nil)
