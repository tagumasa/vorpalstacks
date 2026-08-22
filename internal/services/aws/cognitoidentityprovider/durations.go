package cognitoidentityprovider

import "time"

// durations.go — named lifetimes and list bounds so every expiry and page
// size has a single definition in this package.

// Challenge-session lifetimes. Authentication challenges (custom auth,
// select-challenge, new-password-required, SRP password verifier, hosted-UI
// authorisation codes) are short-lived by design.
const (
	challengeSessionTTL        = 3 * time.Minute
	srpChallengeSessionTTL     = 5 * time.Minute
	hostedUIAuthCodeTTL        = 5 * time.Minute
	hostedUIAuthCodeSweepEvery = time.Minute
)

// verificationCodeTTL is the lifetime of sign-up confirmation and attribute
// verification codes.
const verificationCodeTTL = 24 * time.Hour

// defaultProvisionedLimit is the fallback value reported by
// GetProvisionedLimit when no explicit limit has been set.
const defaultProvisionedLimit = 400000

// TOTP parameters for software token MFA (RFC 6238 defaults).
const (
	totpCodeDigits   = 6
	totpTimeStepSec  = 30
	totpAllowedDrift = 1
)

// maxWebAuthnCredentialListLimit is the Smithy
// WebAuthnCredentialsQueryLimitType upper bound (range {0, 20}).
const maxWebAuthnCredentialListLimit = 20
