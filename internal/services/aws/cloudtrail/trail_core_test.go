package cloudtrail

import (
	"context"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
)

// TestDeleteTrailCoreRejectsEmptyName pins the shared empty-name
// rejection: it fires before any store access, so both the HTTP API and
// the admin console answer InvalidTrailNameException — the error
// DeleteTrail declares — for an omitted trail name instead of a
// not-found.
func TestDeleteTrailCoreRejectsEmptyName(t *testing.T) {
	svc := &CloudTrailService{}
	err := svc.deleteTrailCore(nil, DeleteTrailInput{NameOrARN: ""})
	if err != ErrInvalidTrailName {
		t.Fatalf("expected ErrInvalidTrailName for an empty name, got %v", err)
	}
}

// fakeRoleProvider fakes the two-method RolePolicyProvider the IAM
// validator consults: every role resolves to the provider's trust document.
type fakeRoleProvider struct{ trustDoc string }

func (f fakeRoleProvider) GetAssumeRolePolicyDocument(roleName string) (string, error) {
	return f.trustDoc, nil
}

func (f fakeRoleProvider) Exists(roleName string) bool { return true }

// TestCreateTrailCoreValidatesCloudWatchLogsRole pins the dual-plane
// create contract: the CloudWatchLogsRoleArn trust check runs inside
// createTrailCore whenever a validator is injected — the HTTP request
// context supplies one, and the admin console builds the same validator
// from the service's role provider — while a nil validator leaves the
// member unvalidated, mirroring updateTrailCore. The refusal surfaces as
// the shared IAM validator's own error, the same shape the update path
// returns.
func TestCreateTrailCoreValidatesCloudWatchLogsRole(t *testing.T) {
	ctx := context.Background()
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	// The role exists but trusts only Lambda, so the CloudTrail principal
	// is refused by the trust policy.
	lambdaTrustDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
	validator := iam.NewIAMValidator(fakeRoleProvider{trustDoc: lambdaTrustDoc}, "123456789012")

	in := CreateTrailInput{
		Name:                  "cwl-role-trail",
		S3BucketName:          "cwl-role-trail-bucket",
		Region:                "us-east-1",
		CloudWatchLogsRoleARN: "arn:aws:iam::123456789012:role/untrusted",
		IAMValidator:          validator,
	}
	_, err := svc.createTrailCore(ctx, store, in)
	wire, ok := err.(interface{ GetAWSError() *awserrors.AWSError })
	if !ok {
		t.Fatalf("untrusted role must be rejected with an AWS error surface, got %T", err)
	}
	if got := wire.GetAWSError().GetCode(); got != "InvalidParameterValueException" {
		t.Fatalf("untrusted role: code = %q, want InvalidParameterValueException", got)
	}

	// The validation runs before persistence, so the trail does not exist
	// yet; the same input passes untouched with no validator injected.
	in.IAMValidator = nil
	if _, err := svc.createTrailCore(ctx, store, in); err != nil {
		t.Fatalf("nil validator must leave the role member unvalidated: %v", err)
	}
}
