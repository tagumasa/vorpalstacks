package cognitoidentityprovider

import (
	"testing"

	"vorpalstacks/internal/common/handler"
)

// countingRegistrar records how many times each registration key is
// registered.
type countingRegistrar struct {
	counts map[string]int
}

func (m *countingRegistrar) RegisterHandler(operationName string, h handler.Handler) {
	m.count(":" + operationName)
}

func (m *countingRegistrar) RegisterHandlerForService(serviceName, operationName string, h handler.Handler) {
	m.count(serviceName + ":" + operationName)
}

func (m *countingRegistrar) count(key string) {
	if m.counts == nil {
		m.counts = make(map[string]int)
	}
	m.counts[key]++
}

// TestRegisterHandlersRegistersEachOperationOnce pins the registration
// census: every table row registers its operation exactly once under the
// cognito-idp service name and nothing else is registered. The dispatcher
// registry is last-write-wins, so a duplicated row would silently drop the
// earlier registration together with any WAF enforcement wrapper it
// carried.
func TestRegisterHandlersRegistersEachOperationOnce(t *testing.T) {
	seen := make(map[string]bool, len(userPoolsAPIOperations))
	for _, op := range userPoolsAPIOperations {
		if seen[op.operation] {
			t.Fatalf("operation %q appears more than once in the registration table", op.operation)
		}
		seen[op.operation] = true
	}

	svc := &CognitoService{}
	reg := &countingRegistrar{}
	svc.RegisterHandlers(reg)

	if len(reg.counts) != len(userPoolsAPIOperations) {
		t.Fatalf("registered %d distinct operations, table holds %d", len(reg.counts), len(userPoolsAPIOperations))
	}
	for _, op := range userPoolsAPIOperations {
		if got := reg.counts["cognito-idp:"+op.operation]; got != 1 {
			t.Fatalf("operation %q registered %d times, want exactly 1", op.operation, got)
		}
	}

	// The user pools API has no SignOut operation: client sign-out is the
	// hosted UI logout endpoint, and API sign-out is GlobalSignOut or
	// RevokeToken. Pin its absence so the invented operation cannot be
	// reintroduced.
	if reg.counts["cognito-idp:SignOut"] != 0 {
		t.Fatal("SignOut is registered but is not a user pools API operation")
	}
}

// TestWafInspectedOperationsMatchCredentialFreeSet pins the set of
// operations whose requests AWS WAF inspects — the credential-free
// operations of the Cognito WAF integration. Each of them is registered
// wrapped exactly once because the set is derived from the same
// registration table the wrapper loop reads.
func TestWafInspectedOperationsMatchCredentialFreeSet(t *testing.T) {
	credentialFree := map[string]bool{
		"SignUp":                           true,
		"ConfirmSignUp":                    true,
		"ResendConfirmationCode":           true,
		"InitiateAuth":                     true,
		"RespondToAuthChallenge":           true,
		"ForgotPassword":                   true,
		"ConfirmForgotPassword":            true,
		"GetUser":                          true,
		"UpdateUserAttributes":             true,
		"DeleteUser":                       true,
		"DeleteUserAttributes":             true,
		"GlobalSignOut":                    true,
		"ChangePassword":                   true,
		"GetUserAttributeVerificationCode": true,
		"VerifyUserAttribute":              true,
		"AssociateSoftwareToken":           true,
		"VerifySoftwareToken":              true,
		"SetUserMFAPreference":             true,
		"RevokeToken":                      true,
		"GetClientToken":                   true,
	}
	if len(wafInspectedOperations) != len(credentialFree) {
		t.Fatalf("wafInspectedOperations holds %d operations, want %d", len(wafInspectedOperations), len(credentialFree))
	}
	for name := range credentialFree {
		if !wafInspectedOperations[name] {
			t.Errorf("operation %q is credential-free but is not marked WAF-inspected in the registration table", name)
		}
	}
	for name := range wafInspectedOperations {
		if !credentialFree[name] {
			t.Errorf("operation %q is marked WAF-inspected but is not a credential-free operation", name)
		}
	}
}
