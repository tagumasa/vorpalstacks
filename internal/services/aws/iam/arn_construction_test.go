package iam

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// Every IAM entity ARN is minted by the store's ARNBuilder — one mechanism
// for the whole service. These units pin the exact wire strings the retired
// construction paths produced (the utils ARN builder for service-specific
// credentials, hand-rolled fmt.Sprintf formats for report principals), so
// the consolidation is proven byte-identical for the formats clients see.

const arnConstructionAccountID = "123456789012"

func arnConstructionTestStore(t *testing.T) *iamstore.IAMStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return iamstore.NewIAMStore(st, arnConstructionAccountID)
}

func TestReportPrincipalsMemberARNsComeFromStoreBuilder(t *testing.T) {
	store := arnConstructionTestStore(t)

	// A group report covers the member users, each principal ARN minted by
	// the builder's path-less user format.
	if err := store.UserGroups().AddUserToGroup("arncons-user", "arncons-group"); err != nil {
		t.Fatalf("add user to group: %v", err)
	}
	groupPrincipals := reportPrincipals(store, "arn:aws:iam::"+arnConstructionAccountID+":group/arncons-group")
	assert.Equal(t, []reportPrincipal{
		{arn: "arn:aws:iam::" + arnConstructionAccountID + ":user/arncons-user", name: "arncons-user"},
	}, groupPrincipals)

	// A policy report covers the attached users and roles.
	policyArn := "arn:aws:iam::" + arnConstructionAccountID + ":policy/arncons-policy"
	if err := store.AttachedPolicies().Attach(PrincipalTypeUser, "arncons-attached-user", policyArn); err != nil {
		t.Fatalf("attach policy to user: %v", err)
	}
	if err := store.AttachedPolicies().Attach(PrincipalTypeRole, "arncons-attached-role", policyArn); err != nil {
		t.Fatalf("attach policy to role: %v", err)
	}
	policyPrincipals := reportPrincipals(store, policyArn)
	arns := make(map[string]string, len(policyPrincipals))
	for _, p := range policyPrincipals {
		arns[p.arn] = p.name
	}
	assert.Equal(t, map[string]string{
		"arn:aws:iam::" + arnConstructionAccountID + ":user/arncons-attached-user": "arncons-attached-user",
		"arn:aws:iam::" + arnConstructionAccountID + ":role/arncons-attached-role": "arncons-attached-role",
	}, arns)

	// A user report covers the entity itself; its ARN is the request ARN,
	// never reconstructed.
	userArn := "arn:aws:iam::" + arnConstructionAccountID + ":user/arncons-user"
	userPrincipals := reportPrincipals(store, userArn)
	assert.Equal(t, []reportPrincipal{{arn: userArn, name: "arncons-user"}}, userPrincipals)
}

func TestServiceSpecificCredentialARNsComeFromStoreBuilder(t *testing.T) {
	store := arnConstructionTestStore(t)

	cred, err := store.ServiceSpecificCredentials().Create("arncons-user", "codecommit", 0)
	if err != nil {
		t.Fatalf("create service-specific credential: %v", err)
	}
	assert.Equal(t, "arn:aws:iam::"+arnConstructionAccountID+":user/arncons-user", cred.ServiceSpecificCredentialArn)

	if err := store.ServiceSpecificCredentials().MigrateUser("arncons-user", "arncons-user-renamed"); err != nil {
		t.Fatalf("migrate service-specific credentials on user rename: %v", err)
	}
	migrated, err := store.ServiceSpecificCredentials().Get(cred.ServiceSpecificCredentialId)
	if err != nil {
		t.Fatalf("get migrated credential: %v", err)
	}
	assert.Equal(t, "arncons-user-renamed", migrated.UserName)
	assert.Equal(t, "arn:aws:iam::"+arnConstructionAccountID+":user/arncons-user-renamed", migrated.ServiceSpecificCredentialArn)
}
