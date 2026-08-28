package testutil

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// uuidPattern is the v4-UUID shape the Smithy model places on TermsIdType
// and ManagedLoginBrandingIdType.
var cognitoUUIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[4][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

// cognitoReplicaTermsBrandingTests exercises the user-pool replica, terms
// and managed-login-branding families: the replica enum contract (a created
// replica reports ACTIVE/SECONDARY per ReplicaStatusType/ReplicaRoleType and
// the update requires its Status member), the terms required-member and
// UUID-id contract, and the branding required-ClientId and UUID-id contract.
func (r *TestRunner) cognitoReplicaTermsBrandingTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "UserPoolReplica_EnumConformance", func() error {
		return r.runReplicaEnumTest(tc)
	}))
	results = append(results, r.RunTest("cognito", "UserPoolReplica_UpdateWithoutStatusRejected", func() error {
		return r.runReplicaMissingStatusTest(tc)
	}))
	results = append(results, r.RunTest("cognito", "Terms_RequiredMembersAndUUIDId", func() error {
		return r.runTermsConformanceTest(tc)
	}))
	results = append(results, r.RunTest("cognito", "Terms_MissingRequiredMembersRejected", func() error {
		return r.runTermsNegativeTest(tc)
	}))
	results = append(results, r.RunTest("cognito", "ManagedLoginBranding_RequiredClientAndUUIDId", func() error {
		return r.runBrandingConformanceTest(tc)
	}))
	results = append(results, r.RunTest("cognito", "ManagedLoginBranding_MissingClientIdRejected", func() error {
		return r.runBrandingNegativeTest(tc)
	}))

	return results
}

func (r *TestRunner) runReplicaEnumTest(tc *cognitoIDPContext) error {
	region := "us-west-2"
	resp, err := tc.client.CreateUserPoolReplica(tc.ctx, &cognitoidentityprovider.CreateUserPoolReplicaInput{
		UserPoolId: aws.String(tc.userPoolID),
		RegionName: aws.String(region),
	})
	if err != nil {
		return err
	}
	defer tc.deleteReplica(region)

	replica := resp.UserPoolReplica
	if replica == nil {
		return fmt.Errorf("UserPoolReplica is nil")
	}
	if replica.Status != types.ReplicaStatusTypeActive {
		return fmt.Errorf("created replica Status: got %q, want ACTIVE", replica.Status)
	}
	if replica.Role != types.ReplicaRoleTypeSecondary {
		return fmt.Errorf("created replica Role: got %q, want SECONDARY", replica.Role)
	}
	if replica.UserPoolArn == nil || !strings.Contains(*replica.UserPoolArn, region) {
		return fmt.Errorf("replica UserPoolArn missing or wrong region: %v", replica.UserPoolArn)
	}

	updated, err := tc.client.UpdateUserPoolReplica(tc.ctx, &cognitoidentityprovider.UpdateUserPoolReplicaInput{
		UserPoolId: aws.String(tc.userPoolID),
		RegionName: aws.String(region),
		Status:     types.UpdateReplicaStatusTypeInactive,
	})
	if err != nil {
		return err
	}
	if updated.UserPoolReplica == nil || updated.UserPoolReplica.Status != types.ReplicaStatusTypeInactive {
		return fmt.Errorf("updated replica Status: got %+v, want INACTIVE", updated.UserPoolReplica)
	}
	return nil
}

func (r *TestRunner) runReplicaMissingStatusTest(tc *cognitoIDPContext) error {
	region := "us-east-1"
	_, err := tc.client.CreateUserPoolReplica(tc.ctx, &cognitoidentityprovider.CreateUserPoolReplicaInput{
		UserPoolId: aws.String(tc.userPoolID),
		RegionName: aws.String(region),
	})
	if err != nil {
		return err
	}
	defer tc.deleteReplica(region)

	_, err = tc.client.UpdateUserPoolReplica(tc.ctx, &cognitoidentityprovider.UpdateUserPoolReplicaInput{
		UserPoolId: aws.String(tc.userPoolID),
		RegionName: aws.String(region),
	})
	if err == nil {
		return fmt.Errorf("update without the required Status member succeeded")
	}
	// The typed SDK validates the required Status member client-side, so
	// the request never reaches the server; the server-side rejection is
	// pinned by unit tests in the service package.
	return nil
}

func (tc *cognitoIDPContext) deleteReplica(region string) {
	_, _ = tc.client.DeleteUserPoolReplica(tc.ctx, &cognitoidentityprovider.DeleteUserPoolReplicaInput{
		UserPoolId: aws.String(tc.userPoolID),
		RegionName: aws.String(region),
	})
}

func (r *TestRunner) runTermsConformanceTest(tc *cognitoIDPContext) error {
	clientID, cleanup, err := tc.createPoolClient(tc.userPoolID, tc.unique("terms-client"))
	if err != nil {
		return err
	}
	defer cleanup()

	resp, err := tc.client.CreateTerms(tc.ctx, &cognitoidentityprovider.CreateTermsInput{
		UserPoolId:  aws.String(tc.userPoolID),
		ClientId:    aws.String(clientID),
		TermsName:   aws.String("terms-of-use"),
		TermsSource: types.TermsSourceTypeLink,
		Enforcement: types.TermsEnforcementTypeNone,
	})
	if err != nil {
		return err
	}
	terms := resp.Terms
	if terms == nil {
		return fmt.Errorf("Terms is nil")
	}
	if terms.TermsId == nil || !cognitoUUIDPattern.MatchString(*terms.TermsId) {
		return fmt.Errorf("TermsId is not a v4 UUID: %v", terms.TermsId)
	}
	if terms.TermsSource != types.TermsSourceTypeLink {
		return fmt.Errorf("TermsSource: got %q, want LINK", terms.TermsSource)
	}
	if terms.Enforcement != types.TermsEnforcementTypeNone {
		return fmt.Errorf("Enforcement: got %q, want NONE", terms.Enforcement)
	}
	termsID := *terms.TermsId
	defer tc.deleteTerms(termsID)

	// Terms document names are unique to the app client: a second create
	// with the same name is rejected.
	_, err = tc.client.CreateTerms(tc.ctx, &cognitoidentityprovider.CreateTermsInput{
		UserPoolId:  aws.String(tc.userPoolID),
		ClientId:    aws.String(clientID),
		TermsName:   aws.String("terms-of-use"),
		TermsSource: types.TermsSourceTypeLink,
		Enforcement: types.TermsEnforcementTypeNone,
	})
	if err := expectAWSErrorCode(err, "TermsExistsException"); err != nil {
		return fmt.Errorf("duplicate CreateTerms: %w", err)
	}

	described, err := tc.client.DescribeTerms(tc.ctx, &cognitoidentityprovider.DescribeTermsInput{
		UserPoolId: aws.String(tc.userPoolID),
		TermsId:    aws.String(termsID),
	})
	if err != nil {
		return err
	}
	if described.Terms == nil || described.Terms.TermsId == nil || *described.Terms.TermsId != termsID {
		return fmt.Errorf("DescribeTerms did not return the created terms document")
	}

	updated, err := tc.client.UpdateTerms(tc.ctx, &cognitoidentityprovider.UpdateTermsInput{
		UserPoolId: aws.String(tc.userPoolID),
		TermsId:    aws.String(termsID),
		TermsName:  aws.String("privacy-policy"),
	})
	if err != nil {
		return err
	}
	if updated.Terms == nil || updated.Terms.TermsName == nil || *updated.Terms.TermsName != "privacy-policy" {
		return fmt.Errorf("UpdateTerms did not apply the new TermsName")
	}

	// ListTerms returns the smaller TermsDescriptionType summaries. The SDK
	// output shape has no UserPoolId/ClientId/TermsSource/Links fields, so
	// the member-set restriction is pinned at compile level; here the five
	// documented members must all be populated.
	listed, err := tc.client.ListTerms(tc.ctx, &cognitoidentityprovider.ListTermsInput{
		UserPoolId: aws.String(tc.userPoolID),
	})
	if err != nil {
		return err
	}
	foundListed := false
	for _, d := range listed.Terms {
		if d.TermsId != nil && *d.TermsId == termsID {
			foundListed = true
			if d.TermsName == nil || *d.TermsName != "privacy-policy" {
				return fmt.Errorf("ListTerms TermsName: got %v, want privacy-policy", d.TermsName)
			}
			if d.Enforcement != types.TermsEnforcementTypeNone {
				return fmt.Errorf("ListTerms Enforcement: got %q, want NONE", d.Enforcement)
			}
			if d.CreationDate == nil || d.LastModifiedDate == nil {
				return fmt.Errorf("ListTerms element is missing CreationDate/LastModifiedDate")
			}
		}
	}
	if !foundListed {
		return fmt.Errorf("created terms document not present in ListTerms")
	}

	// Renaming onto a name another document of the same app client holds
	// is likewise a duplicate.
	second, err := tc.client.CreateTerms(tc.ctx, &cognitoidentityprovider.CreateTermsInput{
		UserPoolId:  aws.String(tc.userPoolID),
		ClientId:    aws.String(clientID),
		TermsName:   aws.String("terms-of-use"),
		TermsSource: types.TermsSourceTypeLink,
		Enforcement: types.TermsEnforcementTypeNone,
	})
	if err != nil {
		return err
	}
	defer tc.deleteTerms(*second.Terms.TermsId)
	_, err = tc.client.UpdateTerms(tc.ctx, &cognitoidentityprovider.UpdateTermsInput{
		UserPoolId: aws.String(tc.userPoolID),
		TermsId:    aws.String(*second.Terms.TermsId),
		TermsName:  aws.String("privacy-policy"),
	})
	return expectAWSErrorCode(err, "TermsExistsException")
}

func (r *TestRunner) runTermsNegativeTest(tc *cognitoIDPContext) error {
	clientID, cleanup, err := tc.createPoolClient(tc.userPoolID, tc.unique("terms-neg-client"))
	if err != nil {
		return err
	}
	defer cleanup()

	_, err = tc.client.CreateTerms(tc.ctx, &cognitoidentityprovider.CreateTermsInput{
		UserPoolId: aws.String(tc.userPoolID),
		TermsName:  aws.String("terms-of-use"),
	})
	if err == nil {
		return fmt.Errorf("CreateTerms without ClientId/TermsSource/Enforcement succeeded")
	}

	// The typed SDK validates the required members client-side, so these
	// requests never reach the server; the server-side rejections are
	// pinned by unit tests in the service package.
	_, err = tc.client.CreateTerms(tc.ctx, &cognitoidentityprovider.CreateTermsInput{
		UserPoolId:  aws.String(tc.userPoolID),
		ClientId:    aws.String(clientID),
		TermsName:   aws.String("terms-of-use"),
		TermsSource: types.TermsSourceTypeLink,
	})
	if err == nil {
		return fmt.Errorf("CreateTerms without Enforcement succeeded")
	}
	return nil
}

func (r *TestRunner) runBrandingConformanceTest(tc *cognitoIDPContext) error {
	clientID, cleanup, err := tc.createPoolClient(tc.userPoolID, tc.unique("branding-client"))
	if err != nil {
		return err
	}
	defer cleanup()

	resp, err := tc.client.CreateManagedLoginBranding(tc.ctx, &cognitoidentityprovider.CreateManagedLoginBrandingInput{
		UserPoolId: aws.String(tc.userPoolID),
		ClientId:   aws.String(clientID),
	})
	if err != nil {
		return err
	}
	branding := resp.ManagedLoginBranding
	if branding == nil {
		return fmt.Errorf("ManagedLoginBranding is nil")
	}
	if branding.ManagedLoginBrandingId == nil || !cognitoUUIDPattern.MatchString(*branding.ManagedLoginBrandingId) {
		return fmt.Errorf("ManagedLoginBrandingId is not a v4 UUID: %v", branding.ManagedLoginBrandingId)
	}
	brandingID := *branding.ManagedLoginBrandingId
	defer tc.deleteBranding(brandingID)

	// An app client can only ever have one assigned branding style.
	_, err = tc.client.CreateManagedLoginBranding(tc.ctx, &cognitoidentityprovider.CreateManagedLoginBrandingInput{
		UserPoolId: aws.String(tc.userPoolID),
		ClientId:   aws.String(clientID),
	})
	if err := expectAWSErrorCode(err, "ManagedLoginBrandingExistsException"); err != nil {
		return fmt.Errorf("duplicate CreateManagedLoginBranding: %w", err)
	}

	byClient, err := tc.client.DescribeManagedLoginBrandingByClient(tc.ctx, &cognitoidentityprovider.DescribeManagedLoginBrandingByClientInput{
		UserPoolId: aws.String(tc.userPoolID),
		ClientId:   aws.String(clientID),
	})
	if err != nil {
		return err
	}
	if byClient.ManagedLoginBranding == nil || byClient.ManagedLoginBranding.ManagedLoginBrandingId == nil ||
		*byClient.ManagedLoginBranding.ManagedLoginBrandingId != brandingID {
		return fmt.Errorf("DescribeManagedLoginBrandingByClient did not return the created branding")
	}
	return nil
}

func (r *TestRunner) runBrandingNegativeTest(tc *cognitoIDPContext) error {
	_, err := tc.client.CreateManagedLoginBranding(tc.ctx, &cognitoidentityprovider.CreateManagedLoginBrandingInput{
		UserPoolId: aws.String(tc.userPoolID),
	})
	if err == nil {
		return fmt.Errorf("CreateManagedLoginBranding without ClientId succeeded")
	}
	// The typed SDK validates the required ClientId client-side; the
	// server-side rejection is pinned by unit tests in the service package.
	return nil
}

// deleteTerms removes a terms document created by a test (best effort).
func (tc *cognitoIDPContext) deleteTerms(termsID string) {
	_, _ = tc.client.DeleteTerms(tc.ctx, &cognitoidentityprovider.DeleteTermsInput{
		UserPoolId: aws.String(tc.userPoolID),
		TermsId:    aws.String(termsID),
	})
}

// deleteBranding removes a managed login branding created by a test (best
// effort).
func (tc *cognitoIDPContext) deleteBranding(brandingID string) {
	_, _ = tc.client.DeleteManagedLoginBranding(tc.ctx, &cognitoidentityprovider.DeleteManagedLoginBrandingInput{
		UserPoolId:             aws.String(tc.userPoolID),
		ManagedLoginBrandingId: aws.String(brandingID),
	})
}
