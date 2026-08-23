package cognitoidentityprovider

import (
	"crypto/rand"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// generateTemporaryPassword produces a random password that satisfies the
// pool's password policy, for admin-created users whose invitation is not
// suppressed. Length is the policy minimum (at least eight characters, per
// the AWS default policy floor) plus random padding up to that minimum.
func generateTemporaryPassword(policy *cognitostore.PasswordPolicy) (string, error) {
	length := 8
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
// that writes a native password (administrative set/reset, self-service
// change and confirmation, and the post-verification migration) makes the
// native hash authoritative, and a lingering imported-algorithm flag would
// send subsequent sign-ins down the imported-hash verification path
// against the wrong hash format.
func setNativePasswordCredentials(user *cognitostore.User, userPoolID, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	saltHex, verifierHex, err := computeSrpVerifier(userPoolID, username, password)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	user.PasswordHashAlgo = ""
	user.SrpSalt = saltHex
	user.SrpVerifier = verifierHex
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
