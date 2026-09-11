package iam

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// The typed AWS SDK validates required members client-side, so the
// server-side validation of the simulation and granting-access parameters,
// and the context-entry and report-state machine behaviour, are pinned
// here as unit tests driving the Core functions directly.

const multiValueConditionPolicy = `{
  "Version": "2012-10-17",
  "Statement": [{
    "Sid": "AllowTeamListBucket",
    "Effect": "Allow",
    "Action": ["s3:ListBucket"],
    "Resource": ["*"],
    "Condition": {"ForAnyValue:StringEquals": {"simsample:Department": ["eng", "ops"]}}
  }]
}`

func simulationTestStore(t *testing.T) (*IAMService, *iamstore.IAMStore) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewIAMService("123456789012"), iamstore.NewIAMStore(st, "123456789012")
}

// A multi-valued context key must reach the evaluator with every supplied
// value: the set operators evaluate the whole value set, so a member beyond
// the first can carry the matching value.
func TestSimulatePrincipalPolicyCoreMultiValuedContextKey(t *testing.T) {
	s, store := simulationTestStore(t)

	if _, err := store.Users().Create("simctx-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	policy, err := store.Policies().Create("SimCtxPol", "/", "123456789012", multiValueConditionPolicy, "", nil)
	if err != nil {
		t.Fatalf("create policy: %v", err)
	}
	if err := store.AttachedPolicies().Attach(PrincipalTypeUser, "simctx-user", policy.Arn); err != nil {
		t.Fatalf("attach policy: %v", err)
	}

	base := &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simctx-user",
		ActionNames:     []string{"s3:ListBucket"},
	}

	// The matching value ("eng") is the SECOND member; with the historical
	// member.1-only parsing this evaluation returned implicitDeny.
	result, err := s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn: base.PolicySourceArn,
		ActionNames:     base.ActionNames,
		ContextEntries: []SimulationContextEntry{{
			ContextKeyName:   "simsample:Department",
			ContextKeyValues: []string{"sales", "eng"},
		}},
	})
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if len(result.Evaluations) != 1 {
		t.Fatalf("evaluations: got %d, want 1", len(result.Evaluations))
	}
	ev := result.Evaluations[0]
	if ev.EvalActionName != "s3:ListBucket" || ev.EvalResourceName != "*" {
		t.Fatalf("round-trip fields: got %+v", ev)
	}
	if ev.EvalDecision != "allowed" {
		t.Fatalf("decision with the second context value: got %s, want allowed", ev.EvalDecision)
	}
	// The matched statement resolves to the attached customer-managed
	// policy the simulator gathered.
	if len(ev.MatchedStatements) != 1 {
		t.Fatalf("matched statements: got %v, want one", ev.MatchedStatements)
	}
	assert.Equal(t, "user-managed", ev.MatchedStatements[0].PolicyType)
	assert.Contains(t, ev.MatchedStatements[0].PolicyId, "SimCtxPol")

	// A value set with no member in the policy's list denies.
	result, err = s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn: base.PolicySourceArn,
		ActionNames:     base.ActionNames,
		ContextEntries: []SimulationContextEntry{{
			ContextKeyName:   "simsample:Department",
			ContextKeyValues: []string{"sales", "marketing"},
		}},
	})
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if got := result.Evaluations[0].EvalDecision; got != "implicitDeny" {
		t.Fatalf("decision with disjoint values: got %s, want implicitDeny", got)
	}
}

func TestSimulatePrincipalPolicyCoreValidation(t *testing.T) {
	s, store := simulationTestStore(t)
	if _, err := store.Users().Create("simval-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}

	for name, input := range map[string]*SimulatePrincipalPolicyInput{
		"empty PolicySourceArn": {
			ActionNames: []string{"s3:ListBucket"},
		},
		"empty ActionNames": {
			PolicySourceArn: "arn:aws:iam::123456789012:user/simval-user",
		},
	} {
		_, err := s.simulatePrincipalPolicyCore(store, input)
		if err == nil {
			t.Fatalf("%s must be rejected", name)
		}
		awsErr, ok := err.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("%s: expected AWSError, got %T", name, err)
		}
		assert.Equal(t, http.StatusBadRequest, awsErr.GetHTTPStatusCode())
	}

	_, err := s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simval-user",
		ActionNames:     []string{"s3:ListBucket"},
		PolicyInputList: []string{"{not a policy document"},
	})
	if err == nil {
		t.Fatal("a malformed PolicyInputList document must be rejected")
	}
	assert.Contains(t, err.Error(), "PolicyInputList")

	_, err = s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn:                    "arn:aws:iam::123456789012:user/simval-user",
		ActionNames:                        []string{"s3:ListBucket"},
		PermissionsBoundaryPolicyInputList: []string{"{not a policy document"},
	})
	if err == nil {
		t.Fatal("a malformed boundary document must be rejected")
	}
	assert.Contains(t, err.Error(), "PermissionsBoundaryPolicyInputList")
}

// The permissions boundary is an upper bound: an identity-policy allow
// outside the boundary downgrades to an implicit deny, reported through the
// boundary detail fields.
func TestSimulatePrincipalPolicyCoreBoundaryOverride(t *testing.T) {
	s, store := simulationTestStore(t)
	if _, err := store.Users().Create("simbound-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}

	allowDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"]}]}`
	boundaryDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["sqs:*"],"Resource":["*"]}]}`

	result, err := s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn:                    "arn:aws:iam::123456789012:user/simbound-user",
		ActionNames:                        []string{"s3:ListBucket"},
		PolicyInputList:                    []string{allowDoc},
		PermissionsBoundaryPolicyInputList: []string{boundaryDoc},
	})
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	ev := result.Evaluations[0]
	if ev.EvalDecision != "implicitDeny" {
		t.Fatalf("decision outside the boundary: got %s, want implicitDeny", ev.EvalDecision)
	}
	if !ev.HasBoundary || ev.AllowedByBoundary {
		t.Fatalf("boundary detail: got HasBoundary=%v AllowedByBoundary=%v, want true/false", ev.HasBoundary, ev.AllowedByBoundary)
	}
	if len(ev.MatchedStatements) != 0 {
		t.Fatalf("boundary override must clear the matched statements, got %v", ev.MatchedStatements)
	}

	// Without a boundary the same input is allowed and carries no boundary
	// detail.
	result, err = s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simbound-user",
		ActionNames:     []string{"s3:ListBucket"},
		PolicyInputList: []string{allowDoc},
	})
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	ev = result.Evaluations[0]
	if ev.EvalDecision != "allowed" || ev.HasBoundary {
		t.Fatalf("decision without boundary: got decision=%s HasBoundary=%v, want allowed/false", ev.EvalDecision, ev.HasBoundary)
	}
}

// The ordered organisation policy list is a list of hierarchy levels,
// each carrying its own SCP documents; the reader flattens the levels to
// their documents in wire order.
func TestBuildOrderedOrganizationPolicies(t *testing.T) {
	params := map[string]interface{}{
		"OrderedOrganizationPolicyInputList.member.1.ServiceControlPolicyInputList.member.1": "root-scp",
		"OrderedOrganizationPolicyInputList.member.1.ServiceControlPolicyInputList.member.2": "root-scp-2",
		"OrderedOrganizationPolicyInputList.member.2.ServiceControlPolicyInputList.member.1": "ou-scp",
	}
	assert.Equal(t, []string{"root-scp", "root-scp-2", "ou-scp"}, buildOrderedOrganizationPolicies(params))

	if got := buildOrderedOrganizationPolicies(map[string]interface{}{}); got != nil {
		t.Fatalf("absent parameter: got %v, want nil", got)
	}
}

// The wire parser must collect every value of each context key, not only
// the first member, and skip entries without a usable key or values.
func TestBuildSimulationContextEntriesCollectsAllValues(t *testing.T) {
	params := map[string]interface{}{
		"ContextEntries.member.1.ContextKeyName":            "simsample:Department",
		"ContextEntries.member.1.ContextKeyValues.member.1": "sales",
		"ContextEntries.member.1.ContextKeyValues.member.2": "eng",
		"ContextEntries.member.2.ContextKeyName":            "simsample:Only",
		"ContextEntries.member.2.ContextKeyValues.member.1": "one",
		"ContextEntries.member.3.ContextKeyName":            "simsample:NoValues",
	}

	entries := buildSimulationContextEntries(params)
	if len(entries) != 2 {
		t.Fatalf("entries: got %d, want 2 (valueless entry skipped)", len(entries))
	}
	assert.Equal(t, "simsample:Department", entries[0].ContextKeyName)
	assert.Equal(t, []string{"sales", "eng"}, entries[0].ContextKeyValues)
	assert.Equal(t, "simsample:Only", entries[1].ContextKeyName)
	assert.Equal(t, []string{"one"}, entries[1].ContextKeyValues)
}

// gatedStorage wraps a real storage and, once armed, blocks every read
// until its release channel is closed — so a report generation running
// against it can only finish in the background. Before arming, all
// operations pass through (the IAM store constructor seeds AWS managed
// policies and must not be gated).
type gateState struct {
	armed chan struct{} // closed to start gating
	open  chan struct{} // closed to release blocked reads
}

func (g *gateState) await() {
	select {
	case <-g.armed:
		<-g.open
	default:
	}
}

type gatedBucket struct {
	inner storage.Bucket
	gate  *gateState
}

func (b gatedBucket) Get(key []byte) ([]byte, error) {
	b.gate.await()
	return b.inner.Get(key)
}
func (b gatedBucket) Put(key, value []byte) error {
	b.gate.await()
	return b.inner.Put(key, value)
}
func (b gatedBucket) Delete(key []byte) error {
	b.gate.await()
	return b.inner.Delete(key)
}
func (b gatedBucket) Has(key []byte) bool {
	b.gate.await()
	return b.inner.Has(key)
}
func (b gatedBucket) ForEach(fn func(k, v []byte) error) error {
	b.gate.await()
	return b.inner.ForEach(fn)
}
func (b gatedBucket) ScanPrefix(prefix []byte) storage.Iterator {
	b.gate.await()
	return b.inner.ScanPrefix(prefix)
}
func (b gatedBucket) ScanPrefixReverse(prefix, before []byte) storage.Iterator {
	b.gate.await()
	return b.inner.ScanPrefixReverse(prefix, before)
}
func (b gatedBucket) ScanRange(start, end []byte) storage.Iterator {
	b.gate.await()
	return b.inner.ScanRange(start, end)
}
func (b gatedBucket) Count() int {
	b.gate.await()
	return b.inner.Count()
}

type gatedStorage struct {
	inner storage.BasicStorage
	gate  *gateState
}

func (s gatedStorage) Close() error { return s.inner.Close() }
func (s gatedStorage) Bucket(name string) storage.Bucket {
	return gatedBucket{inner: s.inner.Bucket(name), gate: s.gate}
}
func (s gatedStorage) CreateBucket(name string) error { return s.inner.CreateBucket(name) }
func (s gatedStorage) DeleteBucket(name string) error { return s.inner.DeleteBucket(name) }

// GenerateCredentialReport must return STARTED before the report content
// has been generated: with the storage gated, the generation cannot make
// progress, yet the Core call still returns. While the generation is
// gated, retrieval must return the ReportInProgress fault instead of
// blocking until the build ends, and a second Generate must report the
// running task (INPROGRESS). The report only becomes retrievable after
// the gate opens and the background goroutine finishes.
func TestGenerateCredentialReportReturnsBeforeGeneration(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	gate := &gateState{armed: make(chan struct{}), open: make(chan struct{})}
	store := iamstore.NewIAMStore(gatedStorage{inner: st, gate: gate}, "123456789012")
	close(gate.armed)

	s := NewIAMService("123456789012")

	state, err := s.generateCredentialReportCore(store)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if state != "STARTED" {
		t.Fatalf("state: got %s, want STARTED while generation is gated", state)
	}

	getDone := make(chan error, 1)
	go func() {
		_, _, err := s.getCredentialReportCore()
		getDone <- err
	}()
	select {
	case err := <-getDone:
		if err != ErrReportInProgress {
			t.Fatalf("gated retrieval: got %v, want ErrReportInProgress", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("GetCredentialReport blocked during generation")
	}

	state, err = s.generateCredentialReportCore(store)
	if err != nil {
		t.Fatalf("second generate while gated: %v", err)
	}
	if state != "INPROGRESS" {
		t.Fatalf("second generate while gated: got %s, want INPROGRESS", state)
	}

	close(gate.open)
	s.WaitForReport()

	data, _, err := s.getCredentialReportCore()
	if err != nil {
		t.Fatalf("report after generation: %v", err)
	}
	if !strings.HasPrefix(data, "user,arn,") {
		t.Fatalf("report content: got %q, want the CSV header", data[:min(len(data), 40)])
	}
}

// The retrieval state machine maps each state to its wire error: no report
// yet to ReportNotPresent (410), a running generation to ReportInProgress,
// a failed generation to a 5xx carrying the generation failure reason.
func TestGetCredentialReportCoreStates(t *testing.T) {
	s := NewIAMService("123456789012")

	_, _, err := s.getCredentialReportCore()
	if err != ErrReportNotPresent {
		t.Fatalf("no report: got %v, want ErrReportNotPresent", err)
	}

	s.credentialReportState = "STARTED"
	_, _, err = s.getCredentialReportCore()
	if err != ErrReportInProgress {
		t.Fatalf("running generation: got %v, want ErrReportInProgress", err)
	}

	s.credentialReportState = "FAILED"
	s.credentialReportErr = "list users: simulated storage failure"
	_, _, err = s.getCredentialReportCore()
	if err == nil {
		t.Fatal("failed generation: expected an error")
	}
	awsErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("failed generation: expected *awserrors.AWSError, got %T", err)
	}
	assert.Equal(t, http.StatusInternalServerError, awsErr.GetHTTPStatusCode())
	assert.Contains(t, awsErr.Error(), "Credential report generation failed")
	assert.Contains(t, awsErr.Error(), "simulated storage failure")

	s.credentialReportState = "COMPLETE"
	_, _, err = s.getCredentialReportCore()
	if err != ErrReportNotPresent {
		t.Fatalf("empty completed report: got %v, want ErrReportNotPresent", err)
	}
}

func TestListPoliciesGrantingServiceAccessCore(t *testing.T) {
	s, store := simulationTestStore(t)

	if _, err := store.Users().Create("lpgsa-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	s3Doc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:*"],"Resource":["*"]}]}`
	policy, err := store.Policies().Create("LpgsaS3Pol", "/", "123456789012", s3Doc, "", nil)
	if err != nil {
		t.Fatalf("create policy: %v", err)
	}
	if err := store.AttachedPolicies().Attach(PrincipalTypeUser, "lpgsa-user", policy.Arn); err != nil {
		t.Fatalf("attach policy: %v", err)
	}
	dynamoDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["dynamodb:ListTables"],"Resource":["*"]}]}`
	if err := store.InlinePolicies().Put(PrincipalTypeUser, "lpgsa-user", "LpgsaInlineDynamo", dynamoDoc); err != nil {
		t.Fatalf("put inline policy: %v", err)
	}

	result, err := s.listPoliciesGrantingServiceAccessCore(store, &ListPoliciesGrantingServiceAccessInput{
		Arn:               "arn:aws:iam::123456789012:user/lpgsa-user",
		ServiceNamespaces: []string{"s3", "dynamodb", "sns"},
	})
	if err != nil {
		t.Fatalf("list granting access: %v", err)
	}
	if len(result.Entries) != 3 {
		t.Fatalf("entries: got %d, want 3", len(result.Entries))
	}

	s3Entry, dynamoEntry, snsEntry := result.Entries[0], result.Entries[1], result.Entries[2]
	if len(s3Entry.Policies) != 1 || !s3Entry.Policies[0].Managed || s3Entry.Policies[0].PolicyArn != policy.Arn {
		t.Fatalf("s3 entry: got %+v, want the attached managed policy", s3Entry.Policies)
	}
	if len(dynamoEntry.Policies) != 1 || dynamoEntry.Policies[0].Managed {
		t.Fatalf("dynamodb entry: got %+v, want the inline policy", dynamoEntry.Policies)
	}
	if dynamoEntry.Policies[0].EntityType != PrincipalTypeUser || dynamoEntry.Policies[0].EntityName != "lpgsa-user" {
		t.Fatalf("inline policy entity: got %+v, want the owning user", dynamoEntry.Policies[0])
	}
	if len(snsEntry.Policies) != 0 {
		t.Fatalf("sns entry: got %+v, want no granting policies", snsEntry.Policies)
	}

	// The wire enumeration maps the internal principal types.
	assert.Equal(t, "USER", principalWireEntityType(PrincipalTypeUser))
	assert.Equal(t, "GROUP", principalWireEntityType(PrincipalTypeGroup))
	assert.Equal(t, "ROLE", principalWireEntityType(PrincipalTypeRole))

	manyNamespaces := make([]string, 201)
	for i := range manyNamespaces {
		manyNamespaces[i] = fmt.Sprintf("ns%d", i)
	}
	for name, input := range map[string]*ListPoliciesGrantingServiceAccessInput{
		"empty Arn":               {ServiceNamespaces: []string{"s3"}},
		"empty namespaces":        {Arn: "arn:aws:iam::123456789012:user/lpgsa-user"},
		"too many namespaces":     {Arn: "arn:aws:iam::123456789012:user/lpgsa-user", ServiceNamespaces: manyNamespaces},
		"invalid namespace chars": {Arn: "arn:aws:iam::123456789012:user/lpgsa-user", ServiceNamespaces: []string{"s3;drop"}},
		"namespace too long":      {Arn: "arn:aws:iam::123456789012:user/lpgsa-user", ServiceNamespaces: []string{strings.Repeat("n", 65)}},
	} {
		_, err := s.listPoliciesGrantingServiceAccessCore(store, input)
		if err == nil {
			t.Fatalf("%s must be rejected", name)
		}
		assert.Contains(t, err.Error(), "InvalidInput")
	}
}

// The report rows are byte-identical with the single resolved login-profile
// read: a user with a profile reports TRUE, the changed timestamp and the
// rotation projection derived from that one read; a user without a profile
// reports FALSE/no_information/N/A/N/A.
func TestCredentialReportRowContent(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	store := iamstore.NewIAMStore(st, "123456789012")

	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	changed := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	users := []*iamstore.User{
		{
			ID: "AIDAXEXAMPLEWITHPROFILE", UserName: "alpha",
			Arn: "arn:aws:iam::123456789012:user/alpha", AccountId: "123456789012", CreateDate: created,
		},
		{
			ID: "AIDAXEXAMPLENOLOGIN", UserName: "beta",
			Arn: "arn:aws:iam::123456789012:user/beta", AccountId: "123456789012", CreateDate: created,
		},
	}
	for _, user := range users {
		if err := store.Users().Put(user); err != nil {
			t.Fatalf("put user %s: %v", user.UserName, err)
		}
	}
	if err := store.LoginProfiles().Put(&iamstore.LoginProfile{
		UserName: "alpha", CreateDate: created, PasswordChangedAt: changed,
	}); err != nil {
		t.Fatalf("put login profile: %v", err)
	}
	if err := store.PasswordPolicy().Put(&iamstore.AccountPasswordPolicy{MaxPasswordAge: 30}); err != nil {
		t.Fatalf("put password policy: %v", err)
	}

	report, err := generateReportContentFromStore(store)
	if err != nil {
		t.Fatalf("generate report: %v", err)
	}

	iso := func(ts time.Time) string { return ts.Format(timeutils.ISO8601SimpleFormat) }
	// No MFA devices, access keys or certificates exist: every remaining
	// column after password_next_rotation is FALSE or N/A.
	trailer := ",FALSE,FALSE,N/A,N/A,N/A,N/A,FALSE,N/A,N/A,N/A,N/A,FALSE,N/A,FALSE,N/A"
	alphaRow := "alpha,arn:aws:iam::123456789012:user/alpha," + iso(created) +
		",TRUE,no_information," + iso(changed) + "," + iso(changed.AddDate(0, 0, 30)) + trailer
	betaRow := "beta,arn:aws:iam::123456789012:user/beta," + iso(created) +
		",FALSE,no_information,N/A,N/A" + trailer
	header := "user,arn,user_creation_time,password_enabled,password_last_used," +
		"password_last_changed,password_next_rotation,mfa_active," +
		"access_key_1_active,access_key_1_last_rotated,access_key_1_last_used_date," +
		"access_key_1_last_used_region,access_key_1_last_used_service," +
		"access_key_2_active,access_key_2_last_rotated,access_key_2_last_used_date," +
		"access_key_2_last_used_region,access_key_2_last_used_service," +
		"cert_1_active,cert_1_last_rotated,cert_2_active,cert_2_last_rotated"

	want := header + "\n" + alphaRow + "\n" + betaRow
	if report != want {
		t.Fatalf("report content mismatch:\n got: %q\nwant: %q", report, want)
	}
}

// One evaluation per action regardless of resource count: each resource's
// own decision appears under ResourceSpecificResults, and the top-level
// decision is the most restrictive across resources with the union of the
// statements that determined it. An explicit deny is the only matched
// statement reported when one denies.
func TestSimulatePrincipalPolicyCorePerActionResourceSpecificResults(t *testing.T) {
	s, store := simulationTestStore(t)
	if _, err := store.Users().Create("simres-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}

	allowA := `{"Version":"2012-10-17","Statement":[{"Sid":"AllowA","Effect":"Allow","Action":["s3:ListBucket"],"Resource":["arn:aws:s3:::bucket-a"]}]}`
	allowAll := `{"Version":"2012-10-17","Statement":[{"Sid":"AllowAll","Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"]}]}`
	denyA := `{"Version":"2012-10-17","Statement":[{"Sid":"DenyA","Effect":"Deny","Action":["s3:ListBucket"],"Resource":["arn:aws:s3:::bucket-a"]}]}`

	// Allowed on bucket-a, implicitly denied on bucket-b: the aggregate is
	// the more restrictive implicitDeny, so no statement determined it.
	result, err := s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simres-user",
		ActionNames:     []string{"s3:ListBucket"},
		ResourceArns:    []string{"arn:aws:s3:::bucket-a", "arn:aws:s3:::bucket-b"},
		PolicyInputList: []string{allowA},
	})
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if len(result.Evaluations) != 1 {
		t.Fatalf("evaluations: got %d, want 1 (one per action)", len(result.Evaluations))
	}
	ev := result.Evaluations[0]
	if ev.EvalDecision != "implicitDeny" {
		t.Fatalf("aggregate decision: got %s, want implicitDeny", ev.EvalDecision)
	}
	if len(ev.MatchedStatements) != 0 {
		t.Fatalf("aggregate matched statements: got %v, want none", ev.MatchedStatements)
	}
	if len(ev.ResourceSpecificResults) != 2 {
		t.Fatalf("resource-specific results: got %d, want 2", len(ev.ResourceSpecificResults))
	}
	assert.Equal(t, "arn:aws:s3:::bucket-a", ev.ResourceSpecificResults[0].EvalResourceName)
	assert.Equal(t, "allowed", ev.ResourceSpecificResults[0].EvalResourceDecision)
	if len(ev.ResourceSpecificResults[0].MatchedStatements) != 1 || ev.ResourceSpecificResults[0].MatchedStatements[0].PolicyId != "PolicyInputList.1" {
		t.Fatalf("bucket-a matched statements: got %v, want PolicyInputList.1", ev.ResourceSpecificResults[0].MatchedStatements)
	}
	assert.Equal(t, "arn:aws:s3:::bucket-b", ev.ResourceSpecificResults[1].EvalResourceName)
	assert.Equal(t, "implicitDeny", ev.ResourceSpecificResults[1].EvalResourceDecision)

	// Allowed everywhere: the aggregate is allowed and carries the union of
	// the statements of every equally-deciding resource.
	result, err = s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simres-user",
		ActionNames:     []string{"s3:ListBucket"},
		ResourceArns:    []string{"arn:aws:s3:::bucket-a", "arn:aws:s3:::bucket-b"},
		PolicyInputList: []string{allowAll},
	})
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	ev = result.Evaluations[0]
	if ev.EvalDecision != "allowed" || len(ev.MatchedStatements) != 2 {
		t.Fatalf("aggregate over allowed resources: got decision=%s matched=%v, want allowed with 2 statements", ev.EvalDecision, ev.MatchedStatements)
	}

	// An explicit deny on one resource: the deny statement is the only
	// matched statement reported at the top level.
	result, err = s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simres-user",
		ActionNames:     []string{"s3:ListBucket"},
		ResourceArns:    []string{"arn:aws:s3:::bucket-a", "arn:aws:s3:::bucket-b"},
		PolicyInputList: []string{allowAll, denyA},
	})
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	ev = result.Evaluations[0]
	if ev.EvalDecision != "explicitDeny" {
		t.Fatalf("aggregate decision with a deny: got %s, want explicitDeny", ev.EvalDecision)
	}
	if len(ev.MatchedStatements) != 1 || ev.MatchedStatements[0].PolicyId != "PolicyInputList.2" {
		t.Fatalf("aggregate matched statements with a deny: got %v, want only PolicyInputList.2", ev.MatchedStatements)
	}
}

// MissingContextValues lists the condition keys the included policies
// require that neither the ContextEntries parameter nor the simulator's
// own materialised keys provide: at the top level when the simulated
// resource is "*", per resource when a resource list is supplied.
func TestSimulatePrincipalPolicyCoreMissingContextValues(t *testing.T) {
	s, store := simulationTestStore(t)
	if _, err := store.Users().Create("simmiss-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}

	conditional := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"],"Condition":{"StringEquals":{"simsample:Team":"eng"}}}]}`
	timeBound := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"],"Condition":{"DateGreaterThan":{"aws:CurrentTime":"2020-01-01T00:00:00Z"}}}]}`

	// With the wildcard resource the missing keys sit at the top level; a
	// key the simulator materialises (aws:CurrentTime) is never missing.
	result, err := s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simmiss-user",
		ActionNames:     []string{"s3:ListBucket"},
		PolicyInputList: []string{conditional, timeBound},
	})
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	ev := result.Evaluations[0]
	assert.Equal(t, []string{"simsample:Team"}, ev.MissingContextValues)
	if len(ev.ResourceSpecificResults) != 1 || len(ev.ResourceSpecificResults[0].MissingContextValues) != 0 {
		t.Fatalf("wildcard-resource missing keys must not repeat per resource: %+v", ev.ResourceSpecificResults)
	}

	// With a resource list the missing keys move into each resource entry.
	result, err = s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simmiss-user",
		ActionNames:     []string{"s3:ListBucket"},
		ResourceArns:    []string{"arn:aws:s3:::bucket-a"},
		PolicyInputList: []string{conditional},
	})
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	ev = result.Evaluations[0]
	if len(ev.MissingContextValues) != 0 {
		t.Fatalf("top-level missing keys with a resource list: got %v, want none", ev.MissingContextValues)
	}
	assert.Equal(t, []string{"simsample:Team"}, ev.ResourceSpecificResults[0].MissingContextValues)

	// Supplying the key through ContextEntries clears the missing list.
	result, err = s.simulatePrincipalPolicyCore(store, &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simmiss-user",
		ActionNames:     []string{"s3:ListBucket"},
		PolicyInputList: []string{conditional},
		ContextEntries: []SimulationContextEntry{{
			ContextKeyName:   "simsample:Team",
			ContextKeyValues: []string{"eng"},
		}},
	})
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if got := result.Evaluations[0].MissingContextValues; len(got) != 0 {
		t.Fatalf("missing keys with the context entry supplied: got %v, want none", got)
	}
}

// The PolicyExclusionList union validates its syntax and filters the
// gathered policy set: by type, by ARN with wildcards, and by inline
// policy name with a wildcarded entity name; a permission-boundary
// exclusion drops the caller-supplied boundary; identifiers that match
// nothing are ignored.
func TestSimulatePrincipalPolicyCorePolicyExclusions(t *testing.T) {
	s, store := simulationTestStore(t)
	if _, err := store.Users().Create("simexcl-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	policyArn, err := store.Policies().Create("SimExclPol", "/", "123456789012", `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"]}]}`, "", nil)
	if err != nil {
		t.Fatalf("create policy: %v", err)
	}
	if err := store.AttachedPolicies().Attach(PrincipalTypeUser, "simexcl-user", policyArn.Arn); err != nil {
		t.Fatalf("attach policy: %v", err)
	}
	if err := store.InlinePolicies().Put(PrincipalTypeUser, "simexcl-user", "InlineAllow", `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["sqs:ListQueues"],"Resource":["*"]}]}`); err != nil {
		t.Fatalf("put inline policy: %v", err)
	}

	base := func() *SimulatePrincipalPolicyInput {
		return &SimulatePrincipalPolicyInput{
			PolicySourceArn: "arn:aws:iam::123456789012:user/simexcl-user",
			ActionNames:     []string{"s3:ListBucket", "sqs:ListQueues"},
		}
	}

	// Unfiltered: both actions allowed.
	result, err := s.simulatePrincipalPolicyCore(store, base())
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if result.Evaluations[0].EvalDecision != "allowed" || result.Evaluations[1].EvalDecision != "allowed" {
		t.Fatalf("unfiltered decisions: got %s/%s, want allowed/allowed", result.Evaluations[0].EvalDecision, result.Evaluations[1].EvalDecision)
	}

	// Excluding the attached policy by ARN wildcard removes the s3 allow.
	input := base()
	input.PolicyExclusionList = []SimulationPolicyIdentifier{{PolicyArn: "arn:aws:iam::123456789012:policy/SimExcl*"}}
	result, err = s.simulatePrincipalPolicyCore(store, input)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if result.Evaluations[0].EvalDecision != "implicitDeny" || result.Evaluations[1].EvalDecision != "allowed" {
		t.Fatalf("arn-excluded decisions: got %s/%s, want implicitDeny/allowed", result.Evaluations[0].EvalDecision, result.Evaluations[1].EvalDecision)
	}

	// Excluding by inline identifier with a wildcarded entity removes the
	// sqs allow.
	input = base()
	input.PolicyExclusionList = []SimulationPolicyIdentifier{{InlinePolicyName: "InlineAllow", InlineAttachmentType: "user", InlineAttachmentName: "simexcl-*"}}
	result, err = s.simulatePrincipalPolicyCore(store, input)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if result.Evaluations[0].EvalDecision != "allowed" || result.Evaluations[1].EvalDecision != "implicitDeny" {
		t.Fatalf("inline-excluded decisions: got %s/%s, want allowed/implicitDeny", result.Evaluations[0].EvalDecision, result.Evaluations[1].EvalDecision)
	}

	// Excluding every inline policy by type.
	input = base()
	input.PolicyExclusionList = []SimulationPolicyIdentifier{{PolicyType: "inline"}}
	result, err = s.simulatePrincipalPolicyCore(store, input)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if result.Evaluations[1].EvalDecision != "implicitDeny" {
		t.Fatalf("type-excluded inline decision: got %s, want implicitDeny", result.Evaluations[1].EvalDecision)
	}

	// An identifier that matches nothing is ignored.
	input = base()
	input.PolicyExclusionList = []SimulationPolicyIdentifier{{PolicyArn: "arn:aws:iam::123456789012:policy/Other*"}}
	result, err = s.simulatePrincipalPolicyCore(store, input)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if result.Evaluations[0].EvalDecision != "allowed" {
		t.Fatalf("non-matching exclusion must be ignored, got %s", result.Evaluations[0].EvalDecision)
	}

	// Union syntax violations are rejected.
	for name, identifiers := range map[string][]SimulationPolicyIdentifier{
		"unknown type":       {{PolicyType: "managed"}},
		"no discriminator":   {{}},
		"two discriminators": {{PolicyType: "inline", PolicyArn: policyArn.Arn}},
		"malformed arn":      {{PolicyArn: "not-an-arn"}},
	} {
		input = base()
		input.PolicyExclusionList = identifiers
		if _, err := s.simulatePrincipalPolicyCore(store, input); err == nil {
			t.Fatalf("%s must be rejected", name)
		}
	}

	// A permission-boundary exclusion drops the caller-supplied boundary.
	allowDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"]}]}`
	boundaryDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["sqs:*"],"Resource":["*"]}]}`
	input = base()
	input.PolicyInputList = []string{allowDoc}
	input.PermissionsBoundaryPolicyInputList = []string{boundaryDoc}
	result, err = s.simulatePrincipalPolicyCore(store, input)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if result.Evaluations[0].EvalDecision != "implicitDeny" || !result.Evaluations[0].HasBoundary {
		t.Fatalf("boundary-capped decision: got %s HasBoundary=%v", result.Evaluations[0].EvalDecision, result.Evaluations[0].HasBoundary)
	}
	input.PolicyExclusionList = []SimulationPolicyIdentifier{{PolicyType: "permission-boundary"}}
	result, err = s.simulatePrincipalPolicyCore(store, input)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if result.Evaluations[0].EvalDecision != "allowed" || result.Evaluations[0].HasBoundary {
		t.Fatalf("boundary-excluded decision: got %s HasBoundary=%v, want allowed/false", result.Evaluations[0].EvalDecision, result.Evaluations[0].HasBoundary)
	}
}

// A resource policy participates per resource: same-account allows union
// with identity policies, cross-account allows intersect, and an explicit
// deny always wins. Cross-account simulations report the per-policy-type
// decisions; a role source cannot carry a resource policy.
func TestSimulatePrincipalPolicyCoreResourcePolicy(t *testing.T) {
	s, store := simulationTestStore(t)
	if _, err := store.Users().Create("simrp-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}

	rpDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"arn:aws:iam::123456789012:user/simrp-user"},"Action":["s3:PutObject"],"Resource":["arn:aws:s3:::mary/*"]}]}`
	identityAllow := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:*"],"Resource":["*"]}]}`
	input := func() *SimulatePrincipalPolicyInput {
		return &SimulatePrincipalPolicyInput{
			PolicySourceArn: "arn:aws:iam::123456789012:user/simrp-user",
			ActionNames:     []string{"s3:PutObject"},
			ResourceArns:    []string{"arn:aws:s3:::mary/Test"},
			ResourcePolicy:  rpDoc,
		}
	}

	// Same account (no ResourceOwner): the resource policy alone allows.
	result, err := s.simulatePrincipalPolicyCore(store, input())
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	ev := result.Evaluations[0]
	if ev.EvalDecision != "allowed" {
		t.Fatalf("same-account union: got %s, want allowed", ev.EvalDecision)
	}
	if len(ev.ResourceSpecificResults[0].MatchedStatements) != 1 || ev.ResourceSpecificResults[0].MatchedStatements[0].PolicyId != "ResourcePolicy" {
		t.Fatalf("same-account matched statements: got %v, want ResourcePolicy", ev.ResourceSpecificResults[0].MatchedStatements)
	}
	if ev.EvalDecisionDetails == nil || len(ev.EvalDecisionDetails) != 0 {
		t.Fatalf("same-account decision details must be present and empty, got %v", ev.EvalDecisionDetails)
	}

	// A same-account simulation over only the all-resources wildcard does
	// not return the decision-details member at all.
	wildcard := input()
	wildcard.ResourceArns = nil
	result, err = s.simulatePrincipalPolicyCore(store, wildcard)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if result.Evaluations[0].EvalDecisionDetails != nil {
		t.Fatalf("same-account all-resources decision details must be omitted, got %v", result.Evaluations[0].EvalDecisionDetails)
	}

	// Cross account (foreign ResourceOwner) without an identity allow:
	// intersection denies, the resource policy's statement still appears
	// per resource, and the decision details carry every policy type.
	cross := input()
	cross.ResourceOwner = "arn:aws:iam::123456789013:root"
	result, err = s.simulatePrincipalPolicyCore(store, cross)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	ev = result.Evaluations[0]
	if ev.EvalDecision != "implicitDeny" {
		t.Fatalf("cross-account intersection: got %s, want implicitDeny", ev.EvalDecision)
	}
	if len(ev.MatchedStatements) != 0 {
		t.Fatalf("implicit deny has no determining statements, got %v", ev.MatchedStatements)
	}
	if len(ev.ResourceSpecificResults[0].MatchedStatements) != 1 || ev.ResourceSpecificResults[0].MatchedStatements[0].PolicyId != "ResourcePolicy" {
		t.Fatalf("per-resource matched statements: got %v, want ResourcePolicy", ev.ResourceSpecificResults[0].MatchedStatements)
	}
	details := map[string]string{}
	for _, d := range ev.EvalDecisionDetails {
		details[d.PolicyType] = d.Decision
	}
	if details["IAM Policy"] != "implicitDeny" || details["Resource Policy"] != "allowed" {
		t.Fatalf("cross-account decision details: got %v", details)
	}

	// Cross-account aggregate: each policy type reports its most
	// restrictive decision across resources, even when the deciding
	// resources differ — the identity deny on bucket-a and the resource
	// policy's silence on bucket-b must both surface.
	identityDenyA := `{"Version":"2012-10-17","Statement":[{"Effect":"Deny","Action":["s3:PutObject"],"Resource":["arn:aws:s3:::bucket-a"]}]}`
	identityAllowB := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:PutObject"],"Resource":["arn:aws:s3:::bucket-b"]}]}`
	rpBucketA := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"arn:aws:iam::123456789012:user/simrp-user"},"Action":["s3:PutObject"],"Resource":["arn:aws:s3:::bucket-a"]}]}`
	split := input()
	split.ResourceArns = []string{"arn:aws:s3:::bucket-a", "arn:aws:s3:::bucket-b"}
	split.ResourceOwner = "123456789013"
	split.PolicyInputList = []string{identityDenyA, identityAllowB}
	split.ResourcePolicy = rpBucketA
	result, err = s.simulatePrincipalPolicyCore(store, split)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	ev = result.Evaluations[0]
	if ev.EvalDecision != "explicitDeny" {
		t.Fatalf("split aggregate decision: got %s, want explicitDeny", ev.EvalDecision)
	}
	details = map[string]string{}
	for _, d := range ev.EvalDecisionDetails {
		details[d.PolicyType] = d.Decision
	}
	if details["IAM Policy"] != "explicitDeny" || details["Resource Policy"] != "implicitDeny" {
		t.Fatalf("aggregated decision details: got %v, want IAM Policy=explicitDeny Resource Policy=implicitDeny", details)
	}

	// Cross account with both sides allowing: intersection allows.
	cross = input()
	cross.ResourceOwner = "123456789013"
	cross.PolicyInputList = []string{identityAllow}
	result, err = s.simulatePrincipalPolicyCore(store, cross)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if result.Evaluations[0].EvalDecision != "allowed" {
		t.Fatalf("cross-account both allow: got %s, want allowed", result.Evaluations[0].EvalDecision)
	}

	// An identity explicit deny overrides the resource-policy allow.
	denyDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Deny","Action":["s3:PutObject"],"Resource":["*"]}]}`
	denied := input()
	denied.PolicyInputList = []string{denyDoc}
	result, err = s.simulatePrincipalPolicyCore(store, denied)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	ev = result.Evaluations[0]
	if ev.EvalDecision != "explicitDeny" {
		t.Fatalf("identity deny override: got %s, want explicitDeny", ev.EvalDecision)
	}
	if len(ev.MatchedStatements) != 1 || ev.MatchedStatements[0].PolicyId != "PolicyInputList.1" {
		t.Fatalf("deny matched statements: got %v, want PolicyInputList.1", ev.MatchedStatements)
	}

	// A role source rejects the resource policy outright.
	if _, err := store.Roles().Create("simrp-role", "/", "123456789012", `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`, "", 0, nil, nil); err != nil {
		t.Fatalf("create role: %v", err)
	}
	roleInput := input()
	roleInput.PolicySourceArn = "arn:aws:iam::123456789012:role/simrp-role"
	if _, err := s.simulatePrincipalPolicyCore(store, roleInput); err == nil {
		t.Fatal("a role source with a resource policy must be rejected")
	}
}

// ResourceHandlingOption accepts only the documented scenario values and
// enforces each scenario's required resource types.
func TestSimulatePrincipalPolicyCoreHandlingOption(t *testing.T) {
	s, store := simulationTestStore(t)
	if _, err := store.Users().Create("simho-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	input := &SimulatePrincipalPolicyInput{
		PolicySourceArn: "arn:aws:iam::123456789012:user/simho-user",
		ActionNames:     []string{"ec2:RunInstances"},
	}

	input.ResourceHandlingOption = "NotAScenario"
	if _, err := s.simulatePrincipalPolicyCore(store, input); err == nil {
		t.Fatal("an unknown scenario value must be rejected")
	}

	input.ResourceHandlingOption = "EC2-VPC-EBS"
	input.ResourceArns = []string{"arn:aws:ec2:region:account:instance/i-1"}
	if _, err := s.simulatePrincipalPolicyCore(store, input); err == nil {
		t.Fatal("a scenario without its required resources must be rejected")
	}

	input.ResourceArns = []string{
		"arn:aws:ec2:region:account:instance/i-1",
		"arn:aws:ec2:region:account:image/ami-1",
		"arn:aws:ec2:region:account:security-group/sg-1",
		"arn:aws:ec2:region:account:network-interface/eni-1",
		"arn:aws:ec2:region:account:volume/vol-1",
	}
	if _, err := s.simulatePrincipalPolicyCore(store, input); err != nil {
		t.Fatalf("a satisfied scenario must pass: %v", err)
	}
}

// The custom-policy operation evaluates the supplied documents only, and
// the ordered organisation policy list caps allows like the boundary.
func TestSimulateCustomPolicyCore(t *testing.T) {
	s, _ := simulationTestStore(t)

	allowDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"]}]}`
	allowSQS := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["sqs:*"],"Resource":["*"]}]}`
	input := &SimulateCustomPolicyInput{
		PolicyInputList: []string{allowDoc},
		ActionNames:     []string{"s3:ListBucket", "sqs:ListQueues"},
	}

	if _, err := s.simulateCustomPolicyCore(nil, &SimulateCustomPolicyInput{ActionNames: []string{"s3:ListBucket"}}); err == nil {
		t.Fatal("an empty PolicyInputList must be rejected")
	}

	result, err := s.simulateCustomPolicyCore(nil, input)
	if err != nil {
		t.Fatalf("simulate custom: %v", err)
	}
	if result.Evaluations[0].EvalDecision != "allowed" || result.Evaluations[1].EvalDecision != "implicitDeny" {
		t.Fatalf("custom decisions: got %s/%s, want allowed/implicitDeny", result.Evaluations[0].EvalDecision, result.Evaluations[1].EvalDecision)
	}
	if got := result.Evaluations[0].MatchedStatements[0].PolicyId; got != "PolicyInputList.1" {
		t.Fatalf("custom matched statement id: got %s, want PolicyInputList.1", got)
	}

	// The organisation policy list is an upper bound: an action outside it
	// is implicitly denied even though the input policy allows everything,
	// and the organisation decision detail reports the cap while an action
	// inside the list reports allowance.
	input.OrderedOrganizationPolicyInputList = []string{allowSQS}
	result, err = s.simulateCustomPolicyCore(nil, input)
	if err != nil {
		t.Fatalf("simulate custom with org policies: %v", err)
	}
	if result.Evaluations[0].EvalDecision != "implicitDeny" {
		t.Fatalf("org-capped decision: got %s, want implicitDeny", result.Evaluations[0].EvalDecision)
	}
	if result.Evaluations[1].EvalDecision != "implicitDeny" {
		t.Fatalf("org-allowed decision: got %s, want implicitDeny (input policy denies sqs)", result.Evaluations[1].EvalDecision)
	}
	if result.Evaluations[0].AllowedByOrganizations || !result.Evaluations[0].HasOrgPolicies {
		t.Fatalf("org detail for the capped action: got HasOrgPolicies=%v AllowedByOrganizations=%v, want true/false", result.Evaluations[0].HasOrgPolicies, result.Evaluations[0].AllowedByOrganizations)
	}
	if !result.Evaluations[1].AllowedByOrganizations {
		t.Fatal("org detail for the in-list action: got AllowedByOrganizations=false, want true")
	}

	// The detail renders at the top level of each evaluation result.
	resp := simulationResponse(result.Evaluations, map[string]interface{}{})
	entries := resp["EvaluationResults"].([]interface{})
	orgCapped := entries[0].(map[string]interface{})["OrganizationsDecisionDetail"].(map[string]interface{})
	if orgCapped["AllowedByOrganizations"] != false {
		t.Fatalf("wire org detail for the capped action: got %v, want false", orgCapped["AllowedByOrganizations"])
	}
	orgInList := entries[1].(map[string]interface{})["OrganizationsDecisionDetail"].(map[string]interface{})
	if orgInList["AllowedByOrganizations"] != true {
		t.Fatalf("wire org detail for the in-list action: got %v, want true", orgInList["AllowedByOrganizations"])
	}

	// A caller ARN switches the evaluation principal for resource-policy
	// matching.
	input.OrderedOrganizationPolicyInputList = nil
	input.CallerArn = "arn:aws:iam::123456789012:user/custom-caller"
	input.ResourceArns = []string{"arn:aws:s3:::bucket"}
	input.ResourcePolicy = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"arn:aws:iam::123456789012:user/custom-caller"},"Action":["s3:*"],"Resource":["arn:aws:s3:::bucket"]}]}`
	result, err = s.simulateCustomPolicyCore(nil, input)
	if err != nil {
		t.Fatalf("simulate custom with caller: %v", err)
	}
	if result.Evaluations[0].EvalDecision != "allowed" {
		t.Fatalf("caller-matched resource policy: got %s, want allowed", result.Evaluations[0].EvalDecision)
	}
}

// Missing context values cover identity-based and resource-based policies
// only: a key referenced solely by a permissions boundary or an
// organisation policy is never reported as missing.
func TestSimulateCustomPolicyCoreMissingContextPolicyTypes(t *testing.T) {
	s, _ := simulationTestStore(t)

	identity := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"],"Condition":{"StringEquals":{"simsample:Team":"eng"}}}]}`
	boundary := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"],"Condition":{"StringEquals":{"simsample:BoundaryKey":"eng"}}}]}`
	orgPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"],"Condition":{"StringEquals":{"simsample:ScpKey":"eng"}}}]}`

	result, err := s.simulateCustomPolicyCore(nil, &SimulateCustomPolicyInput{
		PolicyInputList:                    []string{identity},
		PermissionsBoundaryPolicyInputList: []string{boundary},
		OrderedOrganizationPolicyInputList: []string{orgPolicy},
		ActionNames:                        []string{"s3:ListBucket"},
	})
	if err != nil {
		t.Fatalf("simulate custom: %v", err)
	}
	assert.Equal(t, []string{"simsample:Team"}, result.Evaluations[0].MissingContextValues)
}

// An SCP that denies a cross-account simulation ends the evaluation: the
// decision stands but no policy-type details are returned, at the top
// level or per resource.
func TestSimulateCustomPolicyCoreScpDenySuppressesDecisionDetails(t *testing.T) {
	s, _ := simulationTestStore(t)

	allowAll := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:*"],"Resource":["*"]}]}`
	allowSQS := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["sqs:*"],"Resource":["*"]}]}`
	rpDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"arn:aws:iam::123456789012:user/scp-caller"},"Action":["s3:PutObject"],"Resource":["arn:aws:s3:::bucket"]}]}`
	input := &SimulateCustomPolicyInput{
		PolicyInputList:                    []string{allowAll},
		OrderedOrganizationPolicyInputList: []string{allowSQS},
		ActionNames:                        []string{"s3:PutObject"},
		CallerArn:                          "arn:aws:iam::123456789012:user/scp-caller",
		ResourceArns:                       []string{"arn:aws:s3:::bucket"},
		ResourcePolicy:                     rpDoc,
		ResourceOwner:                      "123456789013",
	}

	result, err := s.simulateCustomPolicyCore(nil, input)
	if err != nil {
		t.Fatalf("simulate custom: %v", err)
	}
	ev := result.Evaluations[0]
	if ev.EvalDecision != "implicitDeny" {
		t.Fatalf("SCP-capped decision: got %s, want implicitDeny", ev.EvalDecision)
	}
	if ev.EvalDecisionDetails != nil {
		t.Fatalf("top-level details after an SCP deny must be omitted, got %v", ev.EvalDecisionDetails)
	}
	if ev.ResourceSpecificResults[0].EvalDecisionDetails != nil {
		t.Fatalf("per-resource details after an SCP deny must be omitted, got %v", ev.ResourceSpecificResults[0].EvalDecisionDetails)
	}

	// The same simulation without the SCP reports the details.
	input.OrderedOrganizationPolicyInputList = nil
	result, err = s.simulateCustomPolicyCore(nil, input)
	if err != nil {
		t.Fatalf("simulate custom: %v", err)
	}
	if len(result.Evaluations[0].EvalDecisionDetails) == 0 {
		t.Fatal("cross-account details without an SCP deny must be reported")
	}
}

// Pagination of the evaluation results is positional: a request that
// repeats an action name must not have an earlier duplicate re-delivered
// on the next page.
func TestSimulationResponsePaginatesDuplicateActions(t *testing.T) {
	evaluations := make([]SimulationEvaluation, 3)
	for i := range evaluations {
		evaluations[i] = SimulationEvaluation{
			EvalActionName: "s3:ListBucket",
			EvalDecision:   "implicitDeny",
		}
	}

	page1 := simulationResponse(evaluations, map[string]interface{}{"MaxItems": 2})
	if len(page1["EvaluationResults"].([]interface{})) != 2 || !page1["IsTruncated"].(bool) || page1["Marker"] != "1" {
		t.Fatalf("first page: got results=%v truncated=%v marker=%v, want 2 results truncated with marker 1", page1["EvaluationResults"], page1["IsTruncated"], page1["Marker"])
	}

	page2 := simulationResponse(evaluations, map[string]interface{}{"Marker": page1["Marker"], "MaxItems": 2})
	if results := page2["EvaluationResults"].([]interface{}); len(results) != 1 {
		t.Fatalf("second page: got %d results, want the single remaining duplicate", len(results))
	}
	if truncated, ok := page2["IsTruncated"].(bool); ok && truncated {
		t.Fatal("second page must close the walk")
	}
	if _, ok := page2["Marker"]; ok {
		t.Fatal("the final page carries no marker")
	}
}
