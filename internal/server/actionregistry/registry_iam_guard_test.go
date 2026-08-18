package actionregistry

import (
	"testing"

	iamservice "vorpalstacks/internal/services/aws/iam"
)

// The registry's IAM action list and the IAM service's handler
// registrations must correspond exactly, in both directions: an action
// without a handler misroutes unsigned requests, and a handler missing
// from the registry falls back to weaker classification. This is a pure
// code-to-code comparison and reads no model files.
func TestIAMRegistryMatchesHandlers(t *testing.T) {
	registryActions := IAMActions()
	registered := iamservice.RegisteredIAMOperations()

	actionSet := make(map[string]bool, len(registryActions))
	for _, action := range registryActions {
		actionSet[action] = true
	}
	for _, op := range registered {
		if !actionSet[op] {
			t.Errorf("handler %q is registered by the IAM service but missing from the action registry", op)
		}
	}

	handlerSet := make(map[string]bool, len(registered))
	for _, op := range registered {
		handlerSet[op] = true
	}
	for _, action := range registryActions {
		if !handlerSet[action] {
			t.Errorf("action %q is in the registry but has no registered IAM handler", action)
		}
	}
}
