package iam

// RolePolicyProvider provides IAM role policy operations. The method set
// mirrors the IAM role store's own surface, so the store satisfies this
// interface directly without an adapter.
type RolePolicyProvider interface {
	GetAssumeRolePolicyDocument(roleName string) (string, error)
	Exists(roleName string) bool
}
