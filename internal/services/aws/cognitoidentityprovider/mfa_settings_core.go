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

	pool, err := store.GetUserPool(in.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := applyLegacyMFAOptions(user, pool, in.Params); err != nil {
		return nil, err
	}

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

	pool, err := store.GetUserPool(user.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := applyLegacyMFAOptions(user, pool, in.Params); err != nil {
		return nil, err
	}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// applyLegacyMFAOptions validates the legacy SetUserSettings MFAOptions
// list and applies it to the user record: the entry set is stored for the
// admin projections, and an SMS-delivery entry maps onto the modern SmsMfa
// factor the sign-in second-factor machinery challenges with — under the
// same prerequisites as the modern preference path (pool MFA on, the
// factor enabled in the pool's per-factor configuration, verified
// phone_number).
func applyLegacyMFAOptions(user *cognitostore.User, pool *cognitostore.UserPool, params map[string]interface{}) error {
	opts, err := parseMFAOptions(params)
	if err != nil {
		return err
	}
	_, poolSMS, _ := mfaFactorAvailability(pool)
	for _, opt := range opts {
		if opt.DeliveryMedium == "SMS" {
			if pool.MfaConfiguration == "OFF" || !poolSMS {
				return ErrInvalidParameter
			}
			if !isAttributeVerified(user.Attributes, "phone_number") {
				return ErrInvalidParameter
			}
			user.SmsMfa = &cognitostore.SmsMfaSettings{Enabled: true}
		}
	}
	user.MFAOptions = opts
	return nil
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
// parameters. A present member must be a well-formed list of entries
// carrying both a valid DeliveryMedium and an AttributeName of email or
// phone_number — a malformed member is rejected, not silently dropped.
func parseMFAOptions(params map[string]interface{}) ([]*cognitostore.MFAOptionType, error) {
	val, ok := params["MFAOptions"]
	if !ok {
		return nil, nil
	}
	slice, ok := val.([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	result := make([]*cognitostore.MFAOptionType, 0, len(slice))
	for _, v := range slice {
		m, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		opt := &cognitostore.MFAOptionType{}
		dm, ok := m["DeliveryMedium"].(string)
		if !ok || !validateMFADeliveryMedium(dm) {
			return nil, ErrInvalidParameter
		}
		an, ok := m["AttributeName"].(string)
		if !ok || (an != "email" && an != "phone_number") {
			return nil, ErrInvalidParameter
		}
		opt.DeliveryMedium = dm
		opt.AttributeName = an
		result = append(result, opt)
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

	// The pool's per-factor configuration is the administrator's explicit
	// switch for each second factor; enabling a factor the pool disabled is
	// rejected exactly as the OFF pool rejects every enable.
	poolSoftware, poolSMS, poolEmail := mfaFactorAvailability(pool)

	if sms, ok := params["SMSMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := sms["Enabled"].(bool); enabled {
			if !poolSMS {
				return ErrInvalidParameter
			}
			if !isAttributeVerified(user.Attributes, "phone_number") {
				return ErrInvalidParameter
			}
		}
	}

	if st, ok := params["SoftwareTokenMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := st["Enabled"].(bool); enabled {
			if !poolSoftware {
				return ErrInvalidParameter
			}
			if user.SoftwareTokenMfa == nil || !user.SoftwareTokenMfa.Verified {
				return ErrInvalidParameter
			}
		}
	}

	if em, ok := params["EmailMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := em["Enabled"].(bool); enabled {
			if !poolEmail {
				return ErrInvalidParameter
			}
			if !isAttributeVerified(user.Attributes, "email") {
				return ErrInvalidParameter
			}
		}
	}

	// WebAuthn MFA is exercisable only when the pool runs WebAuthn under a
	// configured relying party id and the user holds at least one
	// registered credential — enabling a factor that could never be
	// challenged or answered is rejected.
	if wa, ok := params["WebAuthnMfaSettings"].(map[string]interface{}); ok {
		if enabled, _ := wa["Enabled"].(bool); enabled {
			if pool.WebAuthnConfiguration == nil || pool.WebAuthnConfiguration.RelyingPartyId == "" {
				return ErrInvalidParameter
			}
			creds, _ := store.ListWebAuthnCredentialsPaginated(user.UserPoolID, user.ID, storecommon.ListOptions{})
			if creds == nil || len(creds.Items) == 0 {
				return ErrInvalidParameter
			}
		}
	}

	return nil
}
