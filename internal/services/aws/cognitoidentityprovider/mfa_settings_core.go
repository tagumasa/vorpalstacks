package cognitoidentityprovider

import (
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// SetUserMFAPreferenceInput carries the wire parameters of
// SetUserMFAPreference. Params holds the raw request parameter map; the
// *MfaSettings nested shapes are read from it inside the Core.
type SetUserMFAPreferenceInput struct {
	AccessToken string
	Params      map[string]interface{}
}

// AdminSetUserMFAPreferenceInput carries the wire parameters of
// AdminSetUserMFAPreference.
type AdminSetUserMFAPreferenceInput struct {
	UserPoolID string
	Username   string
	Params     map[string]interface{}
}

// SetUserSettingsInput carries the wire parameters of SetUserSettings.
type SetUserSettingsInput struct {
	AccessToken string
	Params      map[string]interface{}
}

// AdminSetUserSettingsInput carries the wire parameters of
// AdminSetUserSettings.
type AdminSetUserSettingsInput struct {
	UserPoolID string
	Username   string
	Params     map[string]interface{}
}

// adminSetUserMFAPreferenceCore sets the MFA preferences for a user.
func (s *CognitoService) adminSetUserMFAPreferenceCore(reqCtx *request.RequestContext, in AdminSetUserMFAPreferenceInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	pool, err := store.GetUserPool(in.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := validateMFAPrerequisites(in.Params, store, user, pool); err != nil {
		return nil, err
	}

	applyMFAPreference(user, in.Params)

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// setUserMFAPreferenceCore sets the MFA preferences for the access-token
// caller.
func (s *CognitoService) setUserMFAPreferenceCore(reqCtx *request.RequestContext, in SetUserMFAPreferenceInput) (interface{}, error) {
	if in.AccessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
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

	if err := validateMFAPrerequisites(in.Params, store, user, pool); err != nil {
		return nil, err
	}

	applyMFAPreference(user, in.Params)

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// adminSetUserSettingsCore sets the legacy MFA settings for a user.
func (s *CognitoService) adminSetUserSettingsCore(reqCtx *request.RequestContext, in AdminSetUserSettingsInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	opts, err := parseMFAOptions(in.Params)
	if err != nil {
		return nil, err
	}
	user.MFAOptions = opts

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// setUserSettingsCore sets the legacy MFA settings for the access-token
// caller.
func (s *CognitoService) setUserSettingsCore(reqCtx *request.RequestContext, in SetUserSettingsInput) (interface{}, error) {
	if in.AccessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
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

	opts, err := parseMFAOptions(in.Params)
	if err != nil {
		return nil, err
	}
	user.MFAOptions = opts

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// applyMFAPreference parses SMS/SoftwareToken/Email/WebAuthn MFA settings from
// the request parameters and applies them to the user record. Only one factor
// may be the preferred MFA at a time; setting a new preferred clears the
// previous.
func applyMFAPreference(user *cognitostore.User, params map[string]interface{}) {
	newPreferred := ""

	if sms, ok := params["SMSMfaSettings"]; ok {
		if m, ok := sms.(map[string]interface{}); ok {
			enabled, _ := m["Enabled"].(bool)
			preferred, _ := m["PreferredMfa"].(bool)
			if enabled {
				user.SmsMfa = &cognitostore.SmsMfaSettings{Enabled: true, PreferredMfa: preferred}
				if preferred {
					newPreferred = "SMS"
				}
			} else {
				user.SmsMfa = nil
			}
		}
	}

	if st, ok := params["SoftwareTokenMfaSettings"]; ok {
		if m, ok := st.(map[string]interface{}); ok {
			enabled, _ := m["Enabled"].(bool)
			preferred, _ := m["PreferredMfa"].(bool)
			if enabled {
				if user.SoftwareTokenMfa == nil {
					user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{}
				}
				user.SoftwareTokenMfa.Enabled = true
				user.SoftwareTokenMfa.PreferredMfa = preferred
				if preferred {
					newPreferred = "SoftwareToken"
				}
			} else {
				if user.SoftwareTokenMfa != nil {
					user.SoftwareTokenMfa.Enabled = false
					user.SoftwareTokenMfa.PreferredMfa = false
				}
			}
		}
	}

	if em, ok := params["EmailMfaSettings"]; ok {
		if m, ok := em.(map[string]interface{}); ok {
			enabled, _ := m["Enabled"].(bool)
			preferred, _ := m["PreferredMfa"].(bool)
			if enabled {
				user.EmailMfa = &cognitostore.EmailMfaSettings{Enabled: true, PreferredMfa: preferred}
				if preferred {
					newPreferred = "Email"
				}
			} else {
				user.EmailMfa = nil
			}
		}
	}

	if wa, ok := params["WebAuthnMfaSettings"]; ok {
		if m, ok := wa.(map[string]interface{}); ok {
			if enabled, _ := m["Enabled"].(bool); enabled {
				user.WebAuthnMfaEnabled = true
			} else {
				user.WebAuthnMfaEnabled = false
			}
		}
	}

	// Enforce single preferred MFA: clear all other factors' preferred flag
	if newPreferred != "" {
		clearPreferredMfaExcept(user, newPreferred)
	}
}

// clearPreferredMfaExcept sets PreferredMfa to false on every MFA factor
// except the one named by the keep parameter.
func clearPreferredMfaExcept(user *cognitostore.User, keep string) {
	if keep != "SMS" && user.SmsMfa != nil {
		user.SmsMfa.PreferredMfa = false
	}
	if keep != "SoftwareToken" && user.SoftwareTokenMfa != nil {
		user.SoftwareTokenMfa.PreferredMfa = false
	}
	if keep != "Email" && user.EmailMfa != nil {
		user.EmailMfa.PreferredMfa = false
	}
}

// parseMFAOptions parses the legacy MFAOptions list from the request
// parameters.
func parseMFAOptions(params map[string]interface{}) ([]*cognitostore.MFAOptionType, error) {
	val, ok := params["MFAOptions"]
	if !ok {
		return nil, nil
	}
	slice, ok := val.([]interface{})
	if !ok {
		return nil, nil
	}
	result := make([]*cognitostore.MFAOptionType, 0, len(slice))
	for _, v := range slice {
		if m, ok := v.(map[string]interface{}); ok {
			opt := &cognitostore.MFAOptionType{}
			if dm, ok := m["DeliveryMedium"].(string); ok {
				if !validateMFADeliveryMedium(dm) {
					return nil, ErrInvalidParameter
				}
				opt.DeliveryMedium = dm
			}
			if an, ok := m["AttributeName"].(string); ok {
				if an != "email" && an != "phone_number" {
					return nil, ErrInvalidParameter
				}
				opt.AttributeName = an
			}
			result = append(result, opt)
		}
	}
	return result, nil
}

// validateMFAPrerequisites checks that the requested MFA settings are
// compatible with the pool configuration and the user's current state.
//   - Pool MfaConfiguration "OFF" rejects all MFA preference settings.
//   - SMS MFA requires a verified phone_number attribute.
//   - SoftwareToken MFA requires the user to have an enrolled TOTP secret
//     (user.SoftwareTokenMfa != nil with a non-empty Secret).
//   - Email MFA requires a verified email attribute.
func validateMFAPrerequisites(params map[string]interface{}, store cognitostore.CognitoStoreInterface, user *cognitostore.User, pool *cognitostore.UserPool) error {
	if pool.MfaConfiguration == "OFF" {
		for _, key := range []string{"SMSMfaSettings", "SoftwareTokenMfaSettings", "EmailMfaSettings", "WebAuthnMfaSettings"} {
			if m, ok := params[key].(map[string]interface{}); ok {
				if enabled, _ := m["Enabled"].(bool); enabled {
					return ErrInvalidParameter
				}
			}
		}
		return nil
	}

	if sms, ok := params["SMSMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := sms["Enabled"].(bool); enabled {
			if !isAttributeVerified(user.Attributes, "phone_number") {
				return ErrInvalidParameter
			}
		}
	}

	if st, ok := params["SoftwareTokenMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := st["Enabled"].(bool); enabled {
			if user.SoftwareTokenMfa == nil || !user.SoftwareTokenMfa.Verified {
				return ErrInvalidParameter
			}
		}
	}

	if em, ok := params["EmailMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := em["Enabled"].(bool); enabled {
			if !isAttributeVerified(user.Attributes, "email") {
				return ErrInvalidParameter
			}
		}
	}

	// WebAuthn MFA requires at least one registered WebAuthn credential.
	if wa, ok := params["WebAuthnMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := wa["Enabled"].(bool); enabled {
			creds, _ := store.ListWebAuthnCredentialsPaginated(user.UserPoolID, user.ID, storecommon.ListOptions{})
			if creds == nil || len(creds.Items) == 0 {
				return ErrInvalidParameter
			}
		}
	}

	return nil
}
