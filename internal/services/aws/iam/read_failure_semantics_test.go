package iam

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// The policy-collection walks, the membership resolvers, and the credential
// report generation must never answer from partial input: a store read
// failure surfaces as an error (5xx, cause carried), and the report
// generation moves to a FAILED state. Only a stored document that fails to
// PARSE is skipped — fail-closed for that one policy, with a warning.

// Bucket names below mirror the unexported constants in the store package
// (inlinePolicyBucketName and friends); they are duplicated here because the
// store keeps them unexported.

// readFault selects which storage operations fail on the faulted bucket:
// failGet covers keyed reads, failScan covers ForEach — every paged and
// un-paged listing in the IAM stores routes through it. Writes pass
// unconditionally except Put on the single failPutKey, so test fixtures
// can be laid down before the fault matters while one targeted write
// still fails.
type readFault struct {
	failGet    bool
	failScan   bool
	failPutKey string
	err        error
	// failPutOnAttempt fails the Nth Put to the faulted bucket (1-based,
	// counting attempts regardless of key); 0 disables the counter. Used
	// where the test cannot predict the key, e.g. a job ID minted inside
	// the Core under test.
	failPutOnAttempt int
	putAttempts      int
}

type faultingBucket struct {
	inner storage.Bucket
	fault *readFault
}

func (b faultingBucket) Get(key []byte) ([]byte, error) {
	if b.fault.failGet {
		return nil, b.fault.err
	}
	return b.inner.Get(key)
}
func (b faultingBucket) Put(key, value []byte) error {
	if b.fault.failPutKey != "" && string(key) == b.fault.failPutKey {
		return b.fault.err
	}
	if b.fault.failPutOnAttempt > 0 {
		b.fault.putAttempts++
		if b.fault.putAttempts == b.fault.failPutOnAttempt {
			return b.fault.err
		}
	}
	return b.inner.Put(key, value)
}
func (b faultingBucket) Delete(key []byte) error { return b.inner.Delete(key) }
func (b faultingBucket) Has(key []byte) bool     { return b.inner.Has(key) }
func (b faultingBucket) ForEach(fn func(k, v []byte) error) error {
	if b.fault.failScan {
		return b.fault.err
	}
	return b.inner.ForEach(fn)
}
func (b faultingBucket) ScanPrefix(prefix []byte) storage.Iterator {
	return b.inner.ScanPrefix(prefix)
}
func (b faultingBucket) ScanPrefixReverse(prefix, before []byte) storage.Iterator {
	return b.inner.ScanPrefixReverse(prefix, before)
}
func (b faultingBucket) ScanRange(start, end []byte) storage.Iterator {
	return b.inner.ScanRange(start, end)
}
func (b faultingBucket) Count() int { return b.inner.Count() }

// faultingStorage delegates to a real storage and injects readFault into a
// single named bucket, so a Core walk runs past its earlier reads and fails
// exactly at the collection step under test.
type faultingStorage struct {
	inner      storage.BasicStorage
	bucketName string
	fault      *readFault
}

func (s faultingStorage) Close() error                   { return s.inner.Close() }
func (s faultingStorage) CreateBucket(name string) error { return s.inner.CreateBucket(name) }
func (s faultingStorage) DeleteBucket(name string) error { return s.inner.DeleteBucket(name) }
func (s faultingStorage) Bucket(name string) storage.Bucket {
	if name == s.bucketName {
		return faultingBucket{inner: s.inner.Bucket(name), fault: s.fault}
	}
	return s.inner.Bucket(name)
}

func faultTestStore(t *testing.T, bucketName string, fault *readFault) (*IAMService, *iamstore.IAMStore) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	store := iamstore.NewIAMStore(faultingStorage{inner: st, bucketName: bucketName, fault: fault}, "123456789012")
	return NewIAMService("123456789012"), store
}

func assertInternalFailure(t *testing.T, err error, causeSubstr string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error from a failing store read")
	}
	awsErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("expected *awserrors.AWSError, got %T", err)
	}
	assert.Equal(t, http.StatusInternalServerError, awsErr.GetHTTPStatusCode())
	assert.Contains(t, awsErr.Error(), "InternalFailure")
	assert.Contains(t, awsErr.Error(), causeSubstr)
}

// A storage failure during policy collection must fail SimulatePrincipalPolicy
// outright — the evaluated policy set is never silently narrowed, so a read
// failure can never come back as a partial implicitDeny answer.
func TestSimulatePrincipalPolicyErrorsOnStoreReadFailure(t *testing.T) {
	simInput := func() *SimulatePrincipalPolicyInput {
		return &SimulatePrincipalPolicyInput{
			PolicySourceArn: "arn:aws:iam::123456789012:user/simfail-user",
			ActionNames:     []string{"s3:ListBucket"},
		}
	}

	t.Run("inline policy listing failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_inline_policies", &readFault{
			failScan: true, err: errors.New("simulated inline listing failure"),
		})
		if _, err := store.Users().Create("simfail-user", "/", "123456789012", nil); err != nil {
			t.Fatalf("create user: %v", err)
		}
		_, err := s.simulatePrincipalPolicyCore(store, simInput())
		assertInternalFailure(t, err, "simulated inline listing failure")
	})

	t.Run("attached policy listing failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_attached_policies", &readFault{
			failScan: true, err: errors.New("simulated attachment listing failure"),
		})
		if _, err := store.Users().Create("simfail-user", "/", "123456789012", nil); err != nil {
			t.Fatalf("create user: %v", err)
		}
		_, err := s.simulatePrincipalPolicyCore(store, simInput())
		assertInternalFailure(t, err, "simulated attachment listing failure")
	})
}

// The attachment listing inside ListEntitiesForPolicy surfaces through the
// standard list-error wrapping, like every entity read around it.
func TestListEntitiesForPolicyCoreAttachmentListingFailure(t *testing.T) {
	s, store := faultTestStore(t, "iam_attached_policies", &readFault{
		failScan: true, err: errors.New("simulated attachment scan failure"),
	})
	policyArn := "arn:aws:iam::123456789012:policy/entities-fault"
	if err := store.Policies().Put(&iamstore.Policy{
		Arn: policyArn, PolicyName: "entities-fault", AccountId: "123456789012", Path: "/",
	}); err != nil {
		t.Fatalf("put policy: %v", err)
	}
	_, err := s.listEntitiesForPolicyCore(store, policyArn, "", "", 100)
	assertInternalFailure(t, err, "simulated attachment scan failure")
}

// A stored policy document that fails to parse is skipped — fail-closed for
// that policy — while the remaining policies still evaluate, and the
// simulation never errors on a corrupt document.
func TestSimulatePrincipalPolicySkipsUnparseableStoredDocument(t *testing.T) {
	s, store := simulationTestStore(t)
	if _, err := store.Users().Create("simparse-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	allowDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"]}]}`
	if err := store.InlinePolicies().Put(PrincipalTypeUser, "simparse-user", "CorruptDoc", "{not a policy document"); err != nil {
		t.Fatalf("put corrupt inline policy: %v", err)
	}
	if err := store.InlinePolicies().Put(PrincipalTypeUser, "simparse-user", "AllowListBucket", allowDoc); err != nil {
		t.Fatalf("put inline policy: %v", err)
	}

	input := &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simparse-user",
		ActionNames:     []string{"s3:ListBucket"},
	}
	result, err := s.simulatePrincipalPolicyCore(store, input)
	if err != nil {
		t.Fatalf("simulate with a corrupt document: %v", err)
	}
	if got := result.Evaluations[0].EvalDecision; got != "allowed" {
		t.Fatalf("decision with the valid policy alongside the corrupt one: got %s, want allowed", got)
	}

	// With only the corrupt document the decision is fail-closed: the
	// unparseable policy withholds rather than grants.
	if err := store.InlinePolicies().Delete(PrincipalTypeUser, "simparse-user", "AllowListBucket"); err != nil {
		t.Fatalf("delete inline policy: %v", err)
	}
	result, err = s.simulatePrincipalPolicyCore(store, input)
	if err != nil {
		t.Fatalf("simulate with only the corrupt document: %v", err)
	}
	if got := result.Evaluations[0].EvalDecision; got != "implicitDeny" {
		t.Fatalf("decision with only the corrupt document: got %s, want implicitDeny", got)
	}
}

// The membership resolvers fail the whole listing when a referenced entity
// cannot be read, instead of silently omitting it.
func TestMembershipResolversFailOnEntityReadFailure(t *testing.T) {
	t.Run("ListGroupsForUser group read failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_groups", &readFault{
			failGet: true, err: errors.New("simulated group read failure"),
		})
		if _, err := store.Users().Create("member-user", "/", "123456789012", nil); err != nil {
			t.Fatalf("create user: %v", err)
		}
		if _, err := store.Groups().Create("member-group", "/", "123456789012"); err != nil {
			t.Fatalf("create group: %v", err)
		}
		if err := store.UserGroups().AddUserToGroup("member-user", "member-group"); err != nil {
			t.Fatalf("add user to group: %v", err)
		}
		_, err := s.listGroupsForUserCore(store, "member-user")
		assertInternalFailure(t, err, "simulated group read failure")
	})

	t.Run("ListUsersInGroup user read failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_users", &readFault{
			failGet: true, err: errors.New("simulated user read failure"),
		})
		if _, err := store.Groups().Create("member-group", "/", "123456789012"); err != nil {
			t.Fatalf("create group: %v", err)
		}
		// Lay the membership down directly: creating the user through
		// UserStore is not possible with the users bucket's Get faulted.
		if err := store.UserGroups().AddUserToGroup("member-user", "member-group"); err != nil {
			t.Fatalf("add user to group: %v", err)
		}
		_, err := s.listUsersInGroupCore(store, "member-group")
		assertInternalFailure(t, err, "simulated user read failure")
	})
}

// ListEntitiesForPolicy fails when a listed attachment's entity cannot be
// read, instead of silently omitting the entity from the response.
func TestListEntitiesForPolicyCoreFailsOnEntityReadFailure(t *testing.T) {
	s, store := faultTestStore(t, "iam_users", &readFault{
		failGet: true, err: errors.New("simulated entity read failure"),
	})
	doc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"]}]}`
	policy, err := store.Policies().Create("EntPol", "/", "123456789012", doc, "", nil)
	if err != nil {
		t.Fatalf("create policy: %v", err)
	}
	if err := store.AttachedPolicies().Attach(PrincipalTypeUser, "entity-user", policy.Arn); err != nil {
		t.Fatalf("attach policy: %v", err)
	}
	_, err = s.listEntitiesForPolicyCore(store, policy.Arn, "", "", 100)
	assertInternalFailure(t, err, "simulated entity read failure")
}

// GetAccountAuthorizationDetails propagates the policy-listing failure of
// its per-entity builders instead of emitting truncated detail lists.
func TestGetAccountAuthorizationDetailsCoreFailsOnPolicyListingFailure(t *testing.T) {
	s, store := faultTestStore(t, "iam_inline_policies", &readFault{
		failScan: true, err: errors.New("simulated inline policy listing failure"),
	})
	if _, err := store.Users().Create("authz-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	_, err := s.getAccountAuthorizationDetailsCore(nil, store, &AccountAuthorizationDetailsInput{
		Filters: map[string]bool{"User": true},
	})
	assertInternalFailure(t, err, "simulated inline policy listing failure")
}

// Group deletion must fail on a membership-count storage error instead of
// reading the fault as an empty group and letting the deletion proceed.
func TestDeleteGroupCoreFailsOnUserCountReadFailure(t *testing.T) {
	s, store := faultTestStore(t, "iam_user_groups", &readFault{
		failScan: true, err: errors.New("simulated membership listing failure"),
	})
	if _, err := store.Groups().Create("countfail-group", "/", "123456789012"); err != nil {
		t.Fatalf("create group: %v", err)
	}
	err := s.deleteGroupCore(store, &DeleteGroupInput{GroupName: "countfail-group"})
	if err == nil {
		t.Fatal("expected the membership-count read failure to fail group deletion")
	}
	assert.Contains(t, err.Error(), "simulated membership listing failure")
	// The group record survives: the fault must not masquerade as the
	// empty-group state that permits deletion.
	if !store.Groups().Exists("countfail-group") {
		t.Fatal("group was deleted despite the membership-count read failure")
	}
}

// A failed keyed read during a tag operation must surface as an internal
// failure carrying the cause — the resource may exist while the read
// itself failed, so the answer must not be NoSuchEntity.
func TestTagCoresReadFailureIsInternalFailure(t *testing.T) {
	_, store := faultTestStore(t, "iam_users", &readFault{
		failGet: true,
		err:     errors.New("simulated tag read failure"),
	})

	err := tagResourceCore(store, userTagOps, &TagResourceInput{
		ResourceName: "tagged-user",
		Tags:         []tags.Tag{{Key: "Team", Value: "Ops"}},
	})
	assertInternalFailure(t, err, "simulated tag read failure")

	err = untagResourceCore(store, userTagOps, &UntagResourceInput{
		ResourceName: "tagged-user",
		TagKeys:      []string{"Team"},
	})
	assertInternalFailure(t, err, "simulated tag read failure")

	_, err = listResourceTagsCore(store, userTagOps, &ListResourceTagsInput{
		ResourceName: "tagged-user",
		MaxItems:     100,
	})
	assertInternalFailure(t, err, "simulated tag read failure")
}

// The sentinel path still maps a genuinely absent resource to the
// family's not-found error.
func TestTagCoreMissingEntityStillNotFound(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	store := iamstore.NewIAMStore(st, "123456789012")

	err = tagResourceCore(store, userTagOps, &TagResourceInput{
		ResourceName: "ghost-user",
		Tags:         []tags.Tag{{Key: "Team", Value: "Ops"}},
	})
	if err == nil {
		t.Fatal("expected NoSuchEntity for a missing user")
	}
	awsErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("expected *awserrors.AWSError, got %T", err)
	}
	assert.Equal(t, http.StatusNotFound, awsErr.GetHTTPStatusCode())
	assert.Contains(t, awsErr.Code, "NoSuch")
}

// The entity-resolving prechecks of the delete/list/version cores answer
// from resolved reads: a storage outage surfaces as a 5xx with the cause,
// never as the entity's not-found error.
func TestResolvedEntityReadsFailOutrightOnStoreOutage(t *testing.T) {
	t.Run("deleteUserCore user read failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_users", &readFault{
			failGet: true, err: errors.New("simulated user read failure"),
		})
		err := s.deleteUserCore(store, &DeleteUserInput{UserName: "outage-user"})
		assertInternalFailure(t, err, "simulated user read failure")
	})
	t.Run("deleteRoleCore role read failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_roles", &readFault{
			failGet: true, err: errors.New("simulated role read failure"),
		})
		err := s.deleteRoleCore(store, &DeleteRoleInput{RoleName: "outage-role"})
		assertInternalFailure(t, err, "simulated role read failure")
	})
	t.Run("deleteGroupCore group read failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_groups", &readFault{
			failGet: true, err: errors.New("simulated group read failure"),
		})
		err := s.deleteGroupCore(store, &DeleteGroupInput{GroupName: "outage-group"})
		assertInternalFailure(t, err, "simulated group read failure")
	})
	t.Run("setDefaultPolicyVersionCore policy read failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_policies", &readFault{
			failGet: true, err: errors.New("simulated policy read failure"),
		})
		err := s.setDefaultPolicyVersionCore(store, "arn:aws:iam::123456789012:policy/Outage", "v1")
		assertInternalFailure(t, err, "simulated policy read failure")
	})
	t.Run("simulation entity resolution failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_users", &readFault{
			failGet: true, err: errors.New("simulated source user read failure"),
		})
		_, err := s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
			PolicySourceArn: "arn:aws:iam::123456789012:user/outage-user",
			ActionNames:     []string{"s3:ListBucket"},
		})
		assertInternalFailure(t, err, "simulated source user read failure")
	})
	t.Run("listGroupsForUserCore user read failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_users", &readFault{
			failGet: true, err: errors.New("simulated member user read failure"),
		})
		_, err := s.listGroupsForUserCore(store, "outage-user")
		assertInternalFailure(t, err, "simulated member user read failure")
	})
}
