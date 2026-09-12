package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// generateTemporaryPassword produces a random password that satisfies the
// pool's password policy, for admin-created users whose invitation is not
// suppressed. Length is the policy minimum (at least eight characters, per
// the AWS default policy floor) plus random padding up to that minimum.
func generateTemporaryPassword(policy *cognitostore.PasswordPolicy) (string, error) {
	length := cognitostore.DefaultPasswordMinimumLength
	requireUpper, requireLower, requireNumber, requireSymbol := true, true, true, true
	if policy != nil {
		if policy.MinimumLength > length {
			length = policy.MinimumLength
		}
		requireUpper = policy.RequireUppercase
		requireLower = policy.RequireLowercase
		requireNumber = policy.RequireNumbers
		requireSymbol = policy.RequireSymbols
	}

	var required []byte
	if requireUpper {
		c, err := pickRandomChar("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		if err != nil {
			return "", err
		}
		required = append(required, c)
	}
	if requireLower {
		c, err := pickRandomChar("abcdefghijklmnopqrstuvwxyz")
		if err != nil {
			return "", err
		}
		required = append(required, c)
	}
	if requireNumber {
		c, err := pickRandomChar("0123456789")
		if err != nil {
			return "", err
		}
		required = append(required, c)
	}
	if requireSymbol {
		c, err := pickRandomChar(passwordSpecialChars)
		if err != nil {
			return "", err
		}
		required = append(required, c)
	}

	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789" + passwordSpecialChars
	buf := make([]byte, 0, length)
	buf = append(buf, required...)
	for len(buf) < length {
		c, err := pickRandomChar(all)
		if err != nil {
			return "", err
		}
		buf = append(buf, c)
	}
	// Fisher-Yates shuffle so the required characters are not predictable
	// by position.
	for i := len(buf) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		buf[i], buf[j.Int64()] = buf[j.Int64()], buf[i]
	}
	return string(buf), nil
}

// setNativePasswordCredentials installs the native bcrypt+SRP credential
// pair for a password and clears the imported-hash format flag: every flow
// that writes a native password (administrative set, self-service change and
// reset, sign-up, the NEW_PASSWORD_REQUIRED challenge, and the
// post-verification migration) makes the native hash authoritative, and a
// lingering imported-algorithm flag would send subsequent sign-ins down the
// imported-hash verification path against the wrong hash format. The pool
// policy's PasswordHistorySize is enforced at this single password-write
// site: a password that reuses one of the remembered previous passwords is
// rejected, and the superseded set is retained for the next change.
func setNativePasswordCredentials(user *cognitostore.User, policy *cognitostore.PasswordPolicy, password string) error {
	historySize := 0
	if policy != nil {
		historySize = policy.PasswordHistorySize
	}
	// remembered is the reuse-forbidden set, most recent first: the current
	// native hash plus up to PasswordHistorySize-1 stored predecessors —
	// together the pool's configured number of previous passwords. Imported
	// CSV hashes predate the native credentials and never enter the set,
	// because only hashes this platform generated are known-comparable.
	var remembered []string
	if historySize > 0 {
		remembered = user.PasswordHistory
		if len(remembered) > historySize-1 {
			remembered = remembered[:historySize-1]
		}
		if user.PasswordHash != "" && user.PasswordHashAlgo == "" {
			remembered = append([]string{user.PasswordHash}, remembered...)
		}
		for _, prev := range remembered {
			if bcrypt.CompareHashAndPassword([]byte(prev), []byte(password)) == nil {
				return ErrPasswordHistoryViolation
			}
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	saltHex, verifierHex, err := computeSrpVerifier(user.UserPoolID, user.Username, password)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	user.PasswordHashAlgo = ""
	user.SrpSalt = saltHex
	user.SrpVerifier = verifierHex
	// With no history configured nothing is retained; otherwise the set
	// computed above is exactly what the next change must hold alongside
	// the hash it is about to replace.
	user.PasswordHistory = remembered
	return nil
}

// pickRandomChar picks one character uniformly at random from set. An empty
// set falls back to the documented special characters. A random-source
// failure is propagated: falling back to a fixed character would make the
// generated password partially predictable.
func pickRandomChar(set string) (byte, error) {
	if len(set) == 0 {
		set = passwordSpecialChars
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
	if err != nil {
		return 0, err
	}
	return set[n.Int64()], nil
}

// markAutoVerifiedAttributes sets the "<name>_verified" flag to "true" for
// every attribute listed in the pool's AutoVerifiedAttributes that the user
// actually possesses. Called after user confirmation flows (ConfirmSignUp,
// AdminConfirmSignUp) and for auto-confirmed admin-created users.
func markAutoVerifiedAttributes(user *cognitostore.User, pool *cognitostore.UserPool) {
	if pool == nil {
		return
	}
	for _, attrName := range pool.AutoVerifiedAttributes {
		if user.Attributes[attrName] != "" {
			user.Attributes[attrName+"_verified"] = "true"
		}
	}
}

// isAttributeVerified reports whether the user possesses the named attribute
// AND has it marked as verified.
func isAttributeVerified(attrs map[string]string, name string) bool {
	return attrs[name] != "" && attrs[name+"_verified"] == "true"
}

// issueAttributeUpdateVerification handles the verification round an
// email/phone_number attribute update enters: the new value's verified flag
// is dropped, a fresh code is minted for VerifyUserAttribute, and the
// CustomMessage UpdateUserAttribute trigger fires carrying that code. An
// attribute in the update set is treated as changed — the update writes the
// value, so it must be (re-)verified regardless of what was stored before.
func issueAttributeUpdateVerification(
	ctx context.Context,
	s *CognitoService,
	user *cognitostore.User,
	pool *cognitostore.UserPool,
	attrs map[string]string,
) error {
	for _, name := range []string{"email", "phone_number"} {
		if _, updated := attrs[name]; !updated || user.Attributes[name] == "" {
			continue
		}
		delete(user.Attributes, name+"_verified")
		code, err := generateConfirmationCode()
		if err != nil {
			return ErrInternalError
		}
		if user.AttributeVerificationCodes == nil {
			user.AttributeVerificationCodes = make(map[string]*cognitostore.AttributeVerification)
		}
		user.AttributeVerificationCodes[name] = &cognitostore.AttributeVerification{
			Code:   code,
			Expiry: time.Now().UTC().Add(verificationCodeTTL),
		}
		_, _ = invokeCustomMessage(ctx, s, CustomMessageUpdateUserAttribute, pool.ID, user.Username, "", pool.LambdaConfig, code, userAttributesMap(user), nil)
	}
	return nil
}

// userAttributeList is the single projection of a user's attributes as the
// wire AttributeType list (Name/Value pairs). The schema declares sub a
// required attribute of every user, so the user's immutable identifier is
// materialised into the list even for records whose stored attribute map
// predates that materialisation.
func userAttributeList(user *cognitostore.User) []map[string]string {
	attrs := make(map[string]string, len(user.Attributes)+1)
	for k, v := range user.Attributes {
		attrs[k] = v
	}
	if attrs["sub"] == "" {
		attrs["sub"] = user.ID
	}
	result := make([]map[string]string, 0, len(attrs))
	for k, v := range attrs {
		result = append(result, map[string]string{
			"Name":  k,
			"Value": v,
		})
	}
	return result
}

// passwordSpecialChars is the set of characters AWS documents as satisfying
// the symbol requirement of a password policy. Non-leading and non-trailing
// spaces also count; characters outside every class (for example non-basic-
// Latin letters) satisfy no requirement, matching the AWS contract that the
// required classes are basic Latin letters, numbers, and these specials.
const passwordSpecialChars = `^$*.[]{}()?"!@#%&/\,><':;|_~` + "`=+-"

// isPasswordSymbol reports whether the byte at position i of password is a
// character that satisfies the symbol requirement.
func isPasswordSymbol(password string, i int) bool {
	c := password[i]
	if c == ' ' {
		// Spaces count only when not leading or trailing.
		return i > 0 && i < len(password)-1
	}
	return strings.IndexByte(passwordSpecialChars, c) >= 0
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

	for i := 0; i < len(password); i++ {
		c := password[i]
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasNumber = true
		case isPasswordSymbol(password, i):
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

// determineDeliveryMedium picks the appropriate delivery medium and attribute
// name for confirmation codes, based on the pool's AutoVerifiedAttributes and
// the user's contact information.
func determineDeliveryMedium(pool *cognitostore.UserPool, user *cognitostore.User) (medium, attrName string) {
	// Check AutoVerifiedAttributes first — these take priority
	for _, attr := range pool.AutoVerifiedAttributes {
		if attr == "phone_number" {
			if phone, ok := user.Attributes["phone_number"]; ok && phone != "" {
				return "SMS", "phone_number"
			}
		}
		if attr == "email" {
			if email, ok := user.Attributes["email"]; ok && email != "" {
				return "EMAIL", "email"
			}
		}
	}
	// Fallback: use any available contact attribute
	if email, ok := user.Attributes["email"]; ok && email != "" {
		return "EMAIL", "email"
	}
	if phone, ok := user.Attributes["phone_number"]; ok && phone != "" {
		return "SMS", "phone_number"
	}
	// Default to EMAIL when the user has neither phone nor email — this
	// matches Cognito behaviour which always returns a delivery medium
	// rather than an empty CodeDeliveryDetails.
	return "EMAIL", "email"
}

// computeUserAuthFactors builds the auth factors response for a user, shared
// by GetUserAuthFactors (access-token based) and AdminGetUserAuthFactors
// (admin based). ConfiguredUserAuthFactors carries the AuthFactorType names;
// PreferredMfaSetting and UserMFASettingList use the activated-MFA names the
// model documents for UserMFASettingList.
func computeUserAuthFactors(user *cognitostore.User) map[string]interface{} {
	result := map[string]interface{}{
		"Username": user.Username,
	}

	var configuredFactors []string

	if user.PasswordHash != "" {
		configuredFactors = append(configuredFactors, "PASSWORD")
	}

	if user.SmsMfa != nil && user.SmsMfa.Enabled {
		configuredFactors = append(configuredFactors, "SMS_OTP")
	}

	if user.EmailMfa != nil && user.EmailMfa.Enabled {
		configuredFactors = append(configuredFactors, "EMAIL_OTP")
	}

	if user.SoftwareTokenMfa != nil && user.SoftwareTokenMfa.Enabled {
		configuredFactors = append(configuredFactors, "SOFTWARE_TOKEN")
	}

	if user.WebAuthnMfaEnabled {
		configuredFactors = append(configuredFactors, "WEB_AUTHN")
	}

	if len(user.MFAOptions) > 0 {
		for _, opt := range user.MFAOptions {
			if opt.DeliveryMedium == "SMS" {
				alreadyHas := false
				for _, f := range configuredFactors {
					if f == "SMS_OTP" {
						alreadyHas = true
						break
					}
				}
				if !alreadyHas {
					configuredFactors = append(configuredFactors, "SMS_OTP")
				}
			}
		}
	}

	preferred, settings := userMFAPreferences(user)
	if preferred != "" {
		result["PreferredMfaSetting"] = preferred
	}
	if len(settings) > 0 {
		result["UserMFASettingList"] = settings
	}

	if len(configuredFactors) == 0 {
		configuredFactors = []string{}
	}
	result["ConfiguredUserAuthFactors"] = configuredFactors

	return result
}
