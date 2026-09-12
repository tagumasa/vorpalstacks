package cognitoidentityprovider

import (
	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"github.com/google/uuid"
)

// mfa_signin_core.go — second-factor enforcement at sign-in. The pool's
// MfaConfiguration and the user's enrolled factors decide, after the primary
// credentials have verified, whether the sign-in continues with an MFA
// challenge or mints tokens directly.

// mfaFactorAvailability resolves which second factors the pool's per-factor
// configuration leaves available. An absent per-factor member imposes no
// restriction (the coarse MfaConfiguration governs the pool); a present
// member is the administrator's explicit switch for that factor — a software
// token configuration with Enabled=false, an SMS configuration without its
// SnsCallerArn-carrying SmsConfiguration, and a present email configuration
// with no email factor to configure each disable their factor.
func mfaFactorAvailability(pool *cognitostore.UserPool) (software, sms, email bool) {
	software = pool.MfaConfigurationSoftwareToken == nil || pool.MfaConfigurationSoftwareToken.Enabled
	sms = pool.MfaConfigurationSms == nil || pool.MfaConfigurationSms.SmsConfiguration != nil
	email = pool.EmailMfaConfig == nil
	return software, sms, email
}

// mfaChallengeFor returns the challenge name the sign-in must continue with
// once the user's primary credentials have verified — "" lets the flow mint
// tokens directly. webauthnEnrolled reports whether the user's passkey
// factor is exercisable (the flag is set, the pool runs WebAuthn under a
// configured relying party id, and at least one registered credential
// exists); the caller resolves it against the store because the decision
// function itself is pure.
//
// With MfaConfiguration ON the second factor is mandatory: a user with no
// enrolled factor is sent to MFA_SETUP, a single enrolled factor receives
// that factor's challenge, and a user with several factors chooses through
// SELECT_MFA_TYPE. OPTIONAL applies the same selection whenever the user
// has at least one enrolled factor; with none, the sign-in completes
// without a second factor. The user's preferred factor, when set, enrolled
// and pool-enabled, always wins the selection. A factor the pool's
// per-factor configuration disabled is never challenged even when enrolled.
func mfaChallengeFor(pool *cognitostore.UserPool, user *cognitostore.User, webauthnEnrolled bool) string {
	switch pool.MfaConfiguration {
	case "ON", "OPTIONAL":
	default:
		return ""
	}

	poolSoftware, poolSMS, poolEmail := mfaFactorAvailability(pool)
	software := user.SoftwareTokenMfa != nil && user.SoftwareTokenMfa.Enabled && user.SoftwareTokenMfa.Verified && poolSoftware
	sms := user.SmsMfa != nil && user.SmsMfa.Enabled && poolSMS
	email := user.EmailMfa != nil && user.EmailMfa.Enabled && poolEmail

	if software && user.SoftwareTokenMfa.PreferredMfa {
		return "SOFTWARE_TOKEN_MFA"
	}
	if sms && user.SmsMfa.PreferredMfa {
		return "SMS_MFA"
	}
	if email && user.EmailMfa.PreferredMfa {
		return "EMAIL_OTP"
	}

	enrolled := 0
	for _, ok := range []bool{software, sms, email, webauthnEnrolled} {
		if ok {
			enrolled++
		}
	}
	switch {
	case enrolled == 0:
		if pool.MfaConfiguration == "ON" {
			return "MFA_SETUP"
		}
		return ""
	case enrolled == 1:
		switch {
		case software:
			return "SOFTWARE_TOKEN_MFA"
		case sms:
			return "SMS_MFA"
		case email:
			return "EMAIL_OTP"
		default:
			return "WEB_AUTHN"
		}
	default:
		return "SELECT_MFA_TYPE"
	}
}

// mfaChallengeResponse mints the challenge session the sign-in must
// continue with (generating the one-time code for OTP-typed factors) and
// renders the challenge response. It returns nil when the sign-in may mint
// tokens directly.
func (s *CognitoService) mfaChallengeResponse(reqCtx *request.RequestContext, store cognitostore.CognitoStoreInterface, pool *cognitostore.UserPool, clientID string, user *cognitostore.User) (map[string]interface{}, error) {
	webauthnCreds, werr := listUserWebAuthnCredentials(store, pool.ID, user.ID)
	if werr != nil {
		return nil, ErrInternalError
	}
	webauthnEnrolled := user.WebAuthnMfaEnabled &&
		pool.WebAuthnConfiguration != nil && pool.WebAuthnConfiguration.RelyingPartyId != "" &&
		len(webauthnCreds) > 0
	name := mfaChallengeFor(pool, user, webauthnEnrolled)
	if name == "" {
		return nil, nil
	}
	if name == "WEB_AUTHN" {
		return s.issueWebAuthnChallenge(store, pool, clientID, user, challengeRoleMFA)
	}
	session := newChallengeSession(pool.ID, clientID, user.Username, name, challengeSessionTTL)
	session.ChallengeRole = challengeRoleMFA
	if err := mintChallengeSessionState(store, session); err != nil {
		return nil, err
	}
	s.deliverChallengeCode(reqCtx, pool.ID, name)
	return map[string]interface{}{
		"ChallengeName": name,
		"Session":       session.SessionID,
		"ChallengeParameters": map[string]string{
			"USERNAME": user.Username,
		},
	}, nil
}

// nextChallengeAfterPrimary decides how the sign-in continues once the
// user's primary credentials have verified. Device authentication only ever
// replaces the MFA challenge: when the request names a remembered device by
// DEVICE_KEY in a device-tracking pool whose MFA configuration would issue a
// second-factor challenge, the sign-in continues with DEVICE_SRP_AUTH
// instead — the remembered device's SRP proof is the second factor. A device
// that is unknown, not remembered or holds no stored verifier leaves the
// regular MFA challenge in place.
func (s *CognitoService) nextChallengeAfterPrimary(reqCtx *request.RequestContext, store cognitostore.CognitoStoreInterface, pool *cognitostore.UserPool, clientID string, user *cognitostore.User, deviceKey string) (map[string]interface{}, error) {
	// A remembered device can stand in for the second factor only when the
	// pool's DeviceConfiguration asks for it: per the DeviceConfiguration
	// contract, "when true, a remembered device can sign in with device
	// authentication instead of SMS and time-based one-time password (TOTP)
	// factors for multi-factor authentication (MFA). Whether or not
	// ChallengeRequiredOnNewDevice is true, users who sign in with devices
	// that have not been confirmed or remembered must still provide a second
	// factor". Without the flag, the regular MFA challenge stays in place
	// however well the device is remembered.
	if deviceKey != "" && pool.DeviceConfiguration != nil && pool.DeviceConfiguration.ChallengeRequiredOnNewDevice {
		webauthnCreds, werr := listUserWebAuthnCredentials(store, pool.ID, user.ID)
		if werr != nil {
			return nil, ErrInternalError
		}
		webauthnEnrolled := user.WebAuthnMfaEnabled &&
			pool.WebAuthnConfiguration != nil && pool.WebAuthnConfiguration.RelyingPartyId != "" &&
			len(webauthnCreds) > 0
		if mfaChallengeFor(pool, user, webauthnEnrolled) != "" {
			if device, err := store.GetDevice(pool.ID, user.ID, deviceKey); err == nil &&
				device.DeviceRememberedStatus == "remembered" && device.DeviceSecretVerifierB != "" {
				session := newChallengeSession(pool.ID, clientID, user.Username, "DEVICE_SRP_AUTH", challengeSessionTTL)
				session.DeviceKey = deviceKey
				if err := store.SaveChallengeSession(session); err != nil {
					return nil, err
				}
				return map[string]interface{}{
					"ChallengeName": "DEVICE_SRP_AUTH",
					"Session":       session.SessionID,
					"ChallengeParameters": map[string]string{
						"USERNAME":   user.Username,
						"DEVICE_KEY": deviceKey,
					},
				}, nil
			}
		}
	}
	return s.mfaChallengeResponse(reqCtx, store, pool, clientID, user)
}

// rememberedDeviceForKey reports whether the named device is a remembered,
// verifier-backed device of the user — the state that suppresses the
// NewDeviceMetadata a sign-in otherwise mints for unidentified devices.
func rememberedDeviceForKey(store cognitostore.CognitoStoreInterface, poolID, userID, deviceKey string) bool {
	if deviceKey == "" {
		return false
	}
	device, err := store.GetDevice(poolID, userID, deviceKey)
	return err == nil && device.DeviceRememberedStatus == "remembered" && device.DeviceSecretVerifierB != ""
}

// withNewDeviceMetadata attaches the NewDeviceMetadata member to an
// AuthenticationResult when the pool tracks devices and the sign-in did not
// come from an identified remembered device: a fresh device key in the
// documented {{region}}_{{UUID}} format that the client confirms with
// ConfirmDevice, and the device group key (the pool-name suffix) that the
// client's device SRP secret derivation bakes into its verifier.
func (s *CognitoService) withNewDeviceMetadata(store cognitostore.CognitoStoreInterface, pool *cognitostore.UserPool, user *cognitostore.User, deviceKey string, result map[string]interface{}) map[string]interface{} {
	if pool.DeviceConfiguration == nil || rememberedDeviceForKey(store, pool.ID, user.ID, deviceKey) {
		return result
	}
	groupKey, ok := poolNameFromID(pool.ID)
	if !ok {
		return result
	}
	if authResult, ok := result["AuthenticationResult"].(map[string]interface{}); ok {
		authResult["NewDeviceMetadata"] = map[string]interface{}{
			"DeviceGroupKey": groupKey,
			"DeviceKey":      s.region + "_" + uuid.NewString(),
		}
	}
	return result
}
