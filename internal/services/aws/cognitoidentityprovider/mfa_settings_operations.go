package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// AdminSetUserMFAPreference sets the MFA preferences for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminSetUserMFAPreference.html
func (s *CognitoService) AdminSetUserMFAPreference(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	username := getUsername(req)

	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := validateMFAPrerequisites(req, user, pool); err != nil {
		return nil, err
	}

	applyMFAPreference(user, req)

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// SetUserMFAPreference sets the MFA preferences for the authenticated user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SetUserMFAPreference.html
func (s *CognitoService) SetUserMFAPreference(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	pool, err := store.GetUserPool(user.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := validateMFAPrerequisites(req, user, pool); err != nil {
		return nil, err
	}

	applyMFAPreference(user, req)

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// AdminSetUserSettings sets the legacy MFA settings for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminSetUserSettings.html
func (s *CognitoService) AdminSetUserSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	username := getUsername(req)

	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.MFAOptions = parseMFAOptions(req)

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// SetUserSettings sets the legacy MFA settings for the authenticated user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SetUserSettings.html
func (s *CognitoService) SetUserSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.MFAOptions = parseMFAOptions(req)

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// applyMFAPreference parses SMS/SoftwareToken/Email/WebAuthn MFA settings from
// the request and applies them to the user record.
func applyMFAPreference(user *cognitostore.User, req *request.ParsedRequest) {
	if sms, ok := req.Parameters["SMSMfaSettings"]; ok {
		if m, ok := sms.(map[string]interface{}); ok {
			enabled, _ := m["Enabled"].(bool)
			preferred, _ := m["PreferredMfa"].(bool)
			if enabled {
				user.SmsMfa = &cognitostore.SmsMfaSettings{Enabled: true, PreferredMfa: preferred}
			} else {
				user.SmsMfa = nil
			}
		}
	}

	if st, ok := req.Parameters["SoftwareTokenMfaSettings"]; ok {
		if m, ok := st.(map[string]interface{}); ok {
			enabled, _ := m["Enabled"].(bool)
			preferred, _ := m["PreferredMfa"].(bool)
			if enabled {
				if user.SoftwareTokenMfa == nil {
					user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{}
				}
				user.SoftwareTokenMfa.Enabled = true
				user.SoftwareTokenMfa.PreferredMfa = preferred
			} else {
				if user.SoftwareTokenMfa != nil {
					user.SoftwareTokenMfa.Enabled = false
					user.SoftwareTokenMfa.PreferredMfa = false
				}
			}
		}
	}

	if em, ok := req.Parameters["EmailMfaSettings"]; ok {
		if m, ok := em.(map[string]interface{}); ok {
			enabled, _ := m["Enabled"].(bool)
			preferred, _ := m["PreferredMfa"].(bool)
			if enabled {
				user.EmailMfa = &cognitostore.EmailMfaSettings{Enabled: true, PreferredMfa: preferred}
			} else {
				user.EmailMfa = nil
			}
		}
	}

	if wa, ok := req.Parameters["WebAuthnMfaSettings"]; ok {
		if m, ok := wa.(map[string]interface{}); ok {
			if enabled, _ := m["Enabled"].(bool); enabled {
				user.WebAuthnMfaEnabled = true
			} else {
				user.WebAuthnMfaEnabled = false
			}
		}
	}
}

// parseMFAOptions parses the legacy MFAOptions list from the request.
func parseMFAOptions(req *request.ParsedRequest) []*cognitostore.MFAOptionType {
	val, ok := req.Parameters["MFAOptions"]
	if !ok {
		return nil
	}
	slice, ok := val.([]interface{})
	if !ok {
		return nil
	}
	result := make([]*cognitostore.MFAOptionType, 0, len(slice))
	for _, v := range slice {
		if m, ok := v.(map[string]interface{}); ok {
			opt := &cognitostore.MFAOptionType{}
			if dm, ok := m["DeliveryMedium"].(string); ok {
				opt.DeliveryMedium = dm
			}
			if an, ok := m["AttributeName"].(string); ok {
				opt.AttributeName = an
			}
			result = append(result, opt)
		}
	}
	return result
}

// validateMFAPrerequisites checks that the requested MFA settings are
// compatible with the pool configuration and the user's current state.
//   - Pool MfaConfiguration "OFF" rejects all MFA preference settings.
//   - SMS MFA requires a verified phone_number attribute.
//   - SoftwareToken MFA requires the user to have an enrolled TOTP secret
//     (user.SoftwareTokenMfa != nil with a non-empty Secret).
//   - Email MFA requires a verified email attribute.
func validateMFAPrerequisites(req *request.ParsedRequest, user *cognitostore.User, pool *cognitostore.UserPool) error {
	if pool.MfaConfiguration == "OFF" {
		for _, key := range []string{"SMSMfaSettings", "SoftwareTokenMfaSettings", "EmailMfaSettings", "WebAuthnMfaSettings"} {
			if m, ok := req.Parameters[key].(map[string]interface{}); ok {
				if enabled, _ := m["Enabled"].(bool); enabled {
					return ErrInvalidParameter
				}
			}
		}
		return nil
	}

	if sms, ok := req.Parameters["SMSMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := sms["Enabled"].(bool); enabled {
			if user.Attributes["phone_number"] == "" {
				return ErrInvalidParameter
			}
		}
	}

	if st, ok := req.Parameters["SoftwareTokenMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := st["Enabled"].(bool); enabled {
			if user.SoftwareTokenMfa == nil || !user.SoftwareTokenMfa.Verified {
				return ErrInvalidParameter
			}
		}
	}

	if em, ok := req.Parameters["EmailMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := em["Enabled"].(bool); enabled {
			if user.Attributes["email"] == "" {
				return ErrInvalidParameter
			}
		}
	}

	return nil
}
