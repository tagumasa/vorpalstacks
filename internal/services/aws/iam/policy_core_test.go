package iam

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// An omitted Scope defaults to All ("If it is not included, or if it is
// set to All, all policies are returned"), an invalid Scope is rejected,
// and the defaulting lives in the Core so both planes share it.
func TestListPoliciesCoreScopeDefaultAndValidation(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	store := iamstore.NewIAMStore(st, "123456789012")
	s := NewIAMService("123456789012")

	customer := &iamstore.Policy{
		Arn:        "arn:aws:iam::123456789012:policy/customer-policy",
		PolicyName: "customer-policy",
		AccountId:  "123456789012",
		Path:       "/",
	}
	awsManaged := &iamstore.Policy{
		Arn:        "arn:aws:iam::aws:policy/CataloguePolicy",
		PolicyName: "CataloguePolicy",
		AccountId:  "aws",
		Path:       "/",
	}
	require.NoError(t, store.Policies().Put(customer))
	require.NoError(t, store.Policies().Put(awsManaged))

	arns := func(result *iamstore.PolicyListResult) []string {
		list := make([]string, len(result.Policies))
		for i, policy := range result.Policies {
			list[i] = policy.Arn
		}
		return list
	}

	// The AWS managed catalogue is seeded with the store itself, so the
	// default listing is asserted by membership, not by exact list.
	defaultScoped, err := s.listPoliciesCore(store, "", "", "", false, 1000)
	require.NoError(t, err)
	defaultArns := arns(defaultScoped)
	assert.Contains(t, defaultArns, customer.Arn)
	assert.Contains(t, defaultArns, awsManaged.Arn)

	allScoped, err := s.listPoliciesCore(store, "All", "", "", false, 1000)
	require.NoError(t, err)
	assert.Equal(t, defaultArns, arns(allScoped), "an omitted Scope must behave exactly as All")

	localScoped, err := s.listPoliciesCore(store, "Local", "", "", false, 1000)
	require.NoError(t, err)
	assert.Equal(t, []string{customer.Arn}, arns(localScoped))

	awsScoped, err := s.listPoliciesCore(store, "AWS", "", "", false, 1000)
	require.NoError(t, err)
	assert.NotContains(t, arns(awsScoped), customer.Arn)
	assert.Contains(t, arns(awsScoped), awsManaged.Arn)

	_, err = s.listPoliciesCore(store, "Bogus", "", "", false, 100)
	require.Error(t, err)
	awsErr, ok := err.(*awserrors.AWSError)
	require.True(t, ok)
	assert.Equal(t, "InvalidInput", awsErr.Code)
	assert.Contains(t, awsErr.Message, "Scope")
}

// ListPolicyVersions returns the newest version first in numeric order.
// After cumulative churn past nine versions the live set {v6..v10} must
// list v10 through v6 — bucket key order is lexicographic and would place
// v10 before v6 — and marker pagination follows the same descending
// sequence.
func TestListPolicyVersionsNewestFirst(t *testing.T) {
	store := quotaTestStore(t)

	doc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"]}]}`
	policy, err := store.Policies().Create("ChurnPol", "/", "123456789012", doc, "", nil)
	require.NoError(t, err)
	arn := policy.Arn

	// Churn to the non-sequential live set {v6..v10} the way the
	// documented versioning model produces it: fill the version quota,
	// promote the newest, delete the older versions, create again —
	// identifiers never change and numbering keeps advancing.
	fill := func() {
		for {
			_, err := store.Policies().CreateVersion(arn, doc, false, iamstore.MaxPolicyVersions)
			if errors.Is(err, iamstore.ErrPolicyVersionLimitExceeded) {
				return
			}
			require.NoError(t, err)
		}
	}

	fill() // v2..v5, live set {v1..v5}
	require.NoError(t, store.Policies().SetDefaultVersion(arn, "v5"))
	for _, vid := range []string{"v1", "v2", "v3", "v4"} {
		require.NoError(t, store.Policies().DeleteVersion(arn, vid))
	}

	fill() // v6..v9, live set {v5..v9}
	require.NoError(t, store.Policies().SetDefaultVersion(arn, "v9"))
	require.NoError(t, store.Policies().DeleteVersion(arn, "v5"))
	_, err = store.Policies().CreateVersion(arn, doc, false, iamstore.MaxPolicyVersions)
	require.NoError(t, err) // v10, live set {v6..v10}

	var got []string
	page := ""
	for {
		res, err := store.Policies().ListVersions(arn, page, 100)
		require.NoError(t, err)
		for _, v := range res.Versions {
			got = append(got, v.VersionId)
		}
		if !res.IsTruncated {
			break
		}
		page = res.Marker
	}
	assert.Equal(t, []string{"v10", "v9", "v8", "v7", "v6"}, got)

	res, err := store.Policies().ListVersions(arn, "", 2)
	require.NoError(t, err)
	require.Len(t, res.Versions, 2)
	assert.Equal(t, "v10", res.Versions[0].VersionId)
	assert.Equal(t, "v9", res.Versions[1].VersionId)
	assert.True(t, res.IsTruncated)
	assert.Equal(t, "v9", res.Marker)

	res, err = store.Policies().ListVersions(arn, res.Marker, 2)
	require.NoError(t, err)
	require.Len(t, res.Versions, 2)
	assert.Equal(t, "v8", res.Versions[0].VersionId)
	assert.Equal(t, "v7", res.Versions[1].VersionId)
	assert.True(t, res.IsTruncated)

	res, err = store.Policies().ListVersions(arn, res.Marker, 2)
	require.NoError(t, err)
	require.Len(t, res.Versions, 1)
	assert.Equal(t, "v6", res.Versions[0].VersionId)
	assert.False(t, res.IsTruncated)
}

// The required-member rejection lives in the Core so that both protocol
// planes share it; an empty PolicyArn must surface as a ValidationError,
// not reach the store and masquerade as a missing resource.
func TestPolicyCoresRejectEmptyPolicyArn(t *testing.T) {
	store := quotaTestStore(t)
	s := NewIAMService("123456789012")

	calls := map[string]func() error{
		"getPolicyVersionCore": func() error {
			_, err := s.getPolicyVersionCore(store, "", "v1")
			return err
		},
		"deletePolicyVersionCore": func() error {
			return s.deletePolicyVersionCore(store, "", "v1")
		},
		"listPolicyVersionsCore": func() error {
			_, err := s.listPolicyVersionsCore(store, "", "", 100)
			return err
		},
		"setDefaultPolicyVersionCore": func() error {
			return s.setDefaultPolicyVersionCore(store, "", "v1")
		},
		"listEntitiesForPolicyCore": func() error {
			_, err := s.listEntitiesForPolicyCore(store, "", "", "", 100)
			return err
		},
	}
	for name, call := range calls {
		err := call()
		require.Errorf(t, err, "%s must reject an empty PolicyArn", name)
		awsErr, ok := err.(*awserrors.AWSError)
		require.Truef(t, ok, "%s: expected AWSError, got %T", name, err)
		assert.Equalf(t, "InvalidInput", awsErr.Code, "%s", name)
	}
}
