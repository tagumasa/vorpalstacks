package admin_auth

import (
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/pkg/vsjwt"
)

// UserAdapter adapts an IAM user to the JWTUser interface for token generation.
type UserAdapter struct {
	user     *iamstore.User
	groups   []string
	policies []string
}

// NewUserAdapter creates a new UserAdapter from an IAM user, groups, and policies.
func NewUserAdapter(user *iamstore.User, groups, policies []string) *UserAdapter {
	return &UserAdapter{
		user:     user,
		groups:   groups,
		policies: policies,
	}
}

// GetID returns the user's ID.
func (u *UserAdapter) GetID() string {
	return u.user.UserName
}

// GetUsername returns the user's username.
func (u *UserAdapter) GetUsername() string {
	return u.user.UserName
}

// GetEmail returns the user's email address.
func (u *UserAdapter) GetEmail() string {
	return ""
}

// GetGroups returns the groups the user belongs to.
func (u *UserAdapter) GetGroups() []string {
	return u.groups
}

// GetCustomClaims returns custom claims for JWT token.
func (u *UserAdapter) GetCustomClaims() map[string]interface{} {
	return map[string]interface{}{
		"arn":               u.user.Arn,
		"user_id":           u.user.ID,
		"attached_policies": u.policies,
		"account_id":        u.user.AccountId,
	}
}

var _ vsjwt.JWTUser = (*UserAdapter)(nil)
