package cognitoidentityprovider

import (
	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

func (s *CognitoService) setUserEnabled(reqCtx *request.RequestContext, userPoolID, username string, enabled bool) error {
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return err
	}
	user.Enabled = enabled
	return store.UpdateUser(user)
}

func formatUserAttributes(attrs map[string]string) []map[string]string {
	result := make([]map[string]string, 0)
	for k, v := range attrs {
		result = append(result, map[string]string{
			"Name":  k,
			"Value": v,
		})
	}
	return result
}

func validatePassword(password string, policy *cognitostore.PasswordPolicy) error {
	if policy == nil {
		return nil
	}

	if len(password) < policy.MinimumLength {
		return ErrPasswordPolicyViolation
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSymbol := false

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasNumber = true
		default:
			hasSymbol = true
		}
	}

	if policy.RequireUppercase && !hasUpper {
		return ErrPasswordPolicyViolation
	}
	if policy.RequireLowercase && !hasLower {
		return ErrPasswordPolicyViolation
	}
	if policy.RequireNumbers && !hasNumber {
		return ErrPasswordPolicyViolation
	}
	if policy.RequireSymbols && !hasSymbol {
		return ErrPasswordPolicyViolation
	}

	return nil
}
