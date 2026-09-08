package lambda

import (
	"context"
	"errors"
	"testing"

	"vorpalstacks/internal/common/iam"
)

// fakeRoleProvider fakes the two-method RolePolicyProvider the IAM validator
// consults: known roles return a trust policy for the Lambda service
// principal, unknown roles fail the lookup.
type fakeRoleProvider struct {
	trustDoc string
	missing  bool
}

func (f fakeRoleProvider) GetAssumeRolePolicyDocument(roleName string) (string, error) {
	if f.missing {
		return "", errors.New("role not found")
	}
	return f.trustDoc, nil
}

func (f fakeRoleProvider) RoleExists(roleName string) bool { return !f.missing }

// lambdaTrustDoc allows exactly the Lambda service principal to assume the
// role, the minimum a valid execution role needs.
const lambdaTrustDoc = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}`

func roleTestValidator(provider iam.RolePolicyProvider) *iam.IAMValidator {
	return iam.NewIAMValidator(provider, "000000000000")
}

func roleValidationError(t *testing.T, err error) {
	t.Helper()
	var ipve *iam.InvalidParameterValueError
	if !errors.As(err, &ipve) || ipve.Code != "InvalidParameterValueException" {
		t.Fatalf("expected InvalidParameterValueException role error, got %v", err)
	}
}

// TestCreateFunctionCoreValidatesRole pins that the execution-role trust
// validation runs inside createFunctionCore: a role unknown to the IAM
// provider is rejected even when a validator is injected, instead of the
// validation living only in the HTTP handlers.
func TestCreateFunctionCoreValidatesRole(t *testing.T) {
	t.Setenv("TEST_MODE", "")
	svc, stores := canonicalRuntimeService(t)

	_, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "role-create",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
		IAMValidator: roleTestValidator(fakeRoleProvider{missing: true}),
	})
	roleValidationError(t, err)
}

// TestCreateFunctionCoreAcceptsAssumableRole pins the positive path: a role
// whose trust policy allows the Lambda service principal passes the core's
// validation.
func TestCreateFunctionCoreAcceptsAssumableRole(t *testing.T) {
	t.Setenv("TEST_MODE", "")
	svc, stores := canonicalRuntimeService(t)

	if _, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "role-create-ok",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
		IAMValidator: roleTestValidator(fakeRoleProvider{trustDoc: lambdaTrustDoc}),
	}); err != nil {
		t.Fatalf("create with assumable role: %v", err)
	}
}

// TestCreateFunctionCoreSkipsRoleValidationInTestMode pins the TEST_MODE
// gate: the test-runner mode skips the trust-policy lookup so the regression
// suite can create functions without seeded IAM roles.
func TestCreateFunctionCoreSkipsRoleValidationInTestMode(t *testing.T) {
	t.Setenv("TEST_MODE", "true")
	svc, stores := canonicalRuntimeService(t)

	if _, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "role-create-testmode",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
		IAMValidator: roleTestValidator(fakeRoleProvider{missing: true}),
	}); err != nil {
		t.Fatalf("create under TEST_MODE with unknown role: %v", err)
	}
}

// TestUpdateFunctionConfigurationCoreValidatesRole pins the same core-owned
// validation on the update path: switching to an unknown role is rejected.
func TestUpdateFunctionConfigurationCoreValidatesRole(t *testing.T) {
	t.Setenv("TEST_MODE", "")
	svc, stores := canonicalRuntimeService(t)
	createPresenceFixture(t, svc, stores, "role-update")

	_, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
		&UpdateFunctionConfigurationInput{
			FunctionName: "role-update",
			Role:         "arn:aws:iam::000000000000:role/other",
			IAMValidator: roleTestValidator(fakeRoleProvider{missing: true}),
		})
	roleValidationError(t, err)
}

// TestUpdateFunctionConfigurationCoreSkipsRoleValidationInTestMode pins the
// TEST_MODE gate on the update path.
func TestUpdateFunctionConfigurationCoreSkipsRoleValidationInTestMode(t *testing.T) {
	t.Setenv("TEST_MODE", "true")
	svc, stores := canonicalRuntimeService(t)
	createPresenceFixture(t, svc, stores, "role-update-testmode")

	if _, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
		&UpdateFunctionConfigurationInput{
			FunctionName: "role-update-testmode",
			Role:         "arn:aws:iam::000000000000:role/other",
			IAMValidator: roleTestValidator(fakeRoleProvider{missing: true}),
		}); err != nil {
		t.Fatalf("update under TEST_MODE with unknown role: %v", err)
	}
}
