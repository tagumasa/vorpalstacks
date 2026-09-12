package cognitoidentityprovider

// Sign-in policy bounds. Source: the Smithy model's
// AllowedFirstAuthFactorsListType length trait (min 1, max 5). These are
// the single definitions; every other site must reference them.

const (
	// MinSignInPolicyFirstAuthFactors is the minimum number of permitted
	// first authentication factors a sign-in policy must name.
	MinSignInPolicyFirstAuthFactors = 1
	// MaxSignInPolicyFirstAuthFactors is the maximum number of permitted
	// first authentication factors a sign-in policy may name.
	MaxSignInPolicyFirstAuthFactors = 5
)
