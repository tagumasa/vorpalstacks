package iam

// RolePolicyProvider provides IAM role policy operations.
type RolePolicyProvider interface {
	GetAssumeRolePolicyDocument(roleName string) (string, error)
	RoleExists(roleName string) bool
}
