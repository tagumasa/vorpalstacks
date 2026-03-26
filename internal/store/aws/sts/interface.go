package sts

// SessionStoreInterface defines operations for managing STS sessions.
type SessionStoreInterface interface {
	Create(principalType, principalName, principalArn, roleArn, roleSessionName string, durationSeconds int) (*Session, error)
	Get(sessionToken string) (*Session, error)
	Delete(sessionToken string) error
}

var _ SessionStoreInterface = (*SessionStore)(nil)
