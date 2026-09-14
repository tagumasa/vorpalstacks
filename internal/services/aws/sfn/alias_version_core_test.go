package sfn

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"vorpalstacks/internal/common/request"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// newAliasTestStore provisions a state machine with two published versions
// for the alias idempotency tests.
func newAliasTestStore(t *testing.T) (*sfnstore.StepFunctionStore, string, string, string) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := sfnstore.NewStepFunctionStore(st, "000000000000", "us-east-1")

	ctx := context.Background()
	sm := &sfnstore.StateMachine{
		Name:       "alias-sm",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/sm",
	}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatal(err)
	}
	v1, err := store.PublishStateMachineVersion(ctx, sm.StateMachineArn, "")
	if err != nil {
		t.Fatal(err)
	}

	// A second version for routing-configuration conflicts: publishing
	// follows a revision change, so bump the revision before republishing.
	sm2, err := store.GetStateMachine(ctx, sm.StateMachineArn)
	if err != nil {
		t.Fatal(err)
	}
	sm2.RevisionId = "revision-two"
	if err := store.UpdateStateMachine(ctx, sm2); err != nil {
		t.Fatal(err)
	}
	v2, err := store.PublishStateMachineVersion(ctx, sm.StateMachineArn, "")
	if err != nil {
		t.Fatal(err)
	}
	return store, sm.StateMachineArn, v1.StateMachineVersionArn, v2.StateMachineVersionArn
}

func requireAWSCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected %s error, got nil", code)
	}
	var ae *awserrors.AWSError
	if !errors.As(err, &ae) {
		t.Fatalf("expected an AWSError, got %T: %v", err, err)
	}
	if ae.Code != code {
		t.Fatalf("expected %s, got %s (%s)", code, ae.Code, ae.Message)
	}
}

func TestCreateStateMachineAliasIdempotentRetry(t *testing.T) {
	store, _, v1, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	first, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	second, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	})
	if err != nil {
		t.Fatalf("identical retry must return an idempotent success: %v", err)
	}
	if first["stateMachineAliasArn"] != second["stateMachineAliasArn"] {
		t.Fatalf("identical retry returned a different alias ARN: %v vs %v", first["stateMachineAliasArn"], second["stateMachineAliasArn"])
	}
	if first["creationDate"] != second["creationDate"] {
		t.Fatalf("identical retry returned a different creationDate: %v vs %v", first["creationDate"], second["creationDate"])
	}
}

func TestCreateStateMachineAliasSameNameDifferentParametersConflict(t *testing.T) {
	store, _, v1, v2 := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	if _, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	}); err != nil {
		t.Fatalf("first create: %v", err)
	}

	// Same name with a different description conflicts.
	_, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		Description:   "different",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	})
	requireAWSCode(t, err, "ConflictException")

	// Same name routing to a different version conflicts.
	_, err = svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v2, Weight: 100}},
	})
	requireAWSCode(t, err, "ConflictException")

	// Same name with a split routing configuration conflicts.
	_, err = svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name: "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{
			{StateMachineVersionArn: v1, Weight: 50},
			{StateMachineVersionArn: v2, Weight: 50},
		},
	})
	requireAWSCode(t, err, "ConflictException")
}

func TestPublishStateMachineVersionIdempotentPerRevision(t *testing.T) {
	store, smArn, _, v2 := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	sm, err := store.GetStateMachine(ctx, smArn)
	if err != nil {
		t.Fatal(err)
	}

	// The current revision already has v2 (published by the helper);
	// republishing must return it rather than a new version.
	resp, err := svc.publishStateMachineVersionCore(ctx, store, PublishStateMachineVersionInput{
		StateMachineArn: smArn,
		RevisionId:      sm.RevisionId,
	})
	if err != nil {
		t.Fatalf("idempotent republish: %v", err)
	}
	if resp["stateMachineVersionArn"] != v2 {
		t.Fatalf("republish returned %v, want the existing %s", resp["stateMachineVersionArn"], v2)
	}
	count, err := store.CountStateMachineVersions(smArn)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("expected the two seeded versions, got %d", count)
	}
}

func TestPublishStateMachineVersionRevisionMismatchConflicts(t *testing.T) {
	store, smArn, _, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}

	_, err := svc.publishStateMachineVersionCore(context.Background(), store, PublishStateMachineVersionInput{
		StateMachineArn: smArn,
		RevisionId:      "stale-revision",
	})
	requireAWSCode(t, err, "ConflictException")
}

func TestPublishStateMachineVersionQuota(t *testing.T) {
	store, smArn, _, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	// Fill the quota with distinct revisions; the helper already published
	// two versions, so only the remainder is seeded here.
	sm, err := store.GetStateMachine(ctx, smArn)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < sfnstore.MaxVersionsPerStateMachine-2; i++ {
		sm.RevisionId = fmt.Sprintf("quota-revision-%d", i)
		if err := store.UpdateStateMachine(ctx, sm); err != nil {
			t.Fatal(err)
		}
		if _, err := store.PublishStateMachineVersion(ctx, smArn, ""); err != nil {
			t.Fatal(err)
		}
	}

	count, err := store.CountStateMachineVersions(smArn)
	if err != nil {
		t.Fatal(err)
	}
	if count != sfnstore.MaxVersionsPerStateMachine {
		t.Fatalf("expected %d seeded versions, got %d", sfnstore.MaxVersionsPerStateMachine, count)
	}

	// A fresh revision past the quota is rejected.
	sm.RevisionId = "quota-exceeding-revision"
	if err := store.UpdateStateMachine(ctx, sm); err != nil {
		t.Fatal(err)
	}
	_, err = svc.publishStateMachineVersionCore(ctx, store, PublishStateMachineVersionInput{
		StateMachineArn: smArn,
	})
	requireAWSCode(t, err, "ServiceQuotaExceededException")

	// An idempotent retry of the already-published current revision still
	// succeeds at the quota.
	sm.RevisionId = "quota-revision-0"
	if err := store.UpdateStateMachine(ctx, sm); err != nil {
		t.Fatal(err)
	}
	resp, err := svc.publishStateMachineVersionCore(ctx, store, PublishStateMachineVersionInput{
		StateMachineArn: smArn,
	})
	if err != nil {
		t.Fatalf("idempotent retry at the quota: %v", err)
	}
	if resp["stateMachineVersionArn"] == "" {
		t.Fatal("idempotent retry returned no version ARN")
	}
}

func TestUpdateStateMachineAliasRequiresDescriptionOrRouting(t *testing.T) {
	store, _, v1, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	if _, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          "prod",
		RoutingConfig: []sfnstore.RoutingConfiguration{{StateMachineVersionArn: v1, Weight: 100}},
	}); err != nil {
		t.Fatalf("create: %v", err)
	}
	smArn := smArnOfVersion(t, store, v1)
	alias, err := store.GetStateMachineAliasByName(ctx, smArn, "prod")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.updateStateMachineAliasCore(ctx, store, UpdateStateMachineAliasInput{
		StateMachineAliasArn: alias.StateMachineAliasArn,
	})
	requireAWSCode(t, err, "ValidationException")

	if _, err := svc.updateStateMachineAliasCore(ctx, store, UpdateStateMachineAliasInput{
		StateMachineAliasArn: alias.StateMachineAliasArn,
		DescriptionProvided:  true,
		Description:          "updated",
	}); err != nil {
		t.Fatalf("description-only update: %v", err)
	}
}

func smArnOfVersion(t *testing.T, store *sfnstore.StepFunctionStore, versionArn string) string {
	t.Helper()
	v, err := store.GetStateMachineVersion(context.Background(), versionArn)
	if err != nil {
		t.Fatal(err)
	}
	return v.StateMachineArn
}

func TestRoutingConfigsEqualMatchesEntriesRegardlessOfOrder(t *testing.T) {
	a := []sfnstore.RoutingConfiguration{
		{StateMachineVersionArn: "v1", Weight: 90},
		{StateMachineVersionArn: "v2", Weight: 10},
	}
	sameReordered := []sfnstore.RoutingConfiguration{
		{StateMachineVersionArn: "v2", Weight: 10},
		{StateMachineVersionArn: "v1", Weight: 90},
	}
	differentWeights := []sfnstore.RoutingConfiguration{
		{StateMachineVersionArn: "v1", Weight: 50},
		{StateMachineVersionArn: "v2", Weight: 50},
	}
	if !routingConfigsEqual(a, sameReordered) {
		t.Fatal("reordered entries must compare equal")
	}
	if routingConfigsEqual(a, differentWeights) {
		t.Fatal("different weights must not compare equal")
	}
	if routingConfigsEqual(a, a[:1]) {
		t.Fatal("different lengths must not compare equal")
	}
}

// TestParseRoutingConfigurationErrorVocabulary pins that routing-configuration
// problems report ValidationException — the routing configuration is an
// operation parameter, not a state machine definition, matching the sibling
// validation in the alias Core. Malformed shape is a type error here, not
// a silently degraded empty configuration.
func TestParseRoutingConfigurationErrorVocabulary(t *testing.T) {
	for name, raw := range map[string]interface{}{
		"non-JSON string":  "not json",
		"empty string":     "",
		"non-array object": map[string]interface{}{"weight": 1},
		"non-object entry": []interface{}{"arn:aws:states:us-east-1:000000000000:stateMachine:sm:1"},
		"mixed entry":      []interface{}{map[string]interface{}{"weight": 100}, 7},
	} {
		req := &request.ParsedRequest{Parameters: map[string]interface{}{"routingConfiguration": raw}}
		_, err := parseRoutingConfiguration(req)
		requireAWSCode(t, err, "ValidationException")
		_ = name
	}
}

// TestListStateMachineVersionsMaxResultsZeroUsesDefaultPage pins the
// Core-level page default: zero means unset on both protocol planes (the
// AWS SDK leaves maxResults absent; the admin console's optional field
// reads as zero), so an offset-paged listing with zero falls back to the
// documented page size — an unnormalised zero emits an empty page and a
// "0" nextToken that reproduces the empty page forever.
func TestListStateMachineVersionsMaxResultsZeroUsesDefaultPage(t *testing.T) {
	store, smArn, _, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	sm, err := store.GetStateMachine(ctx, smArn)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < sfnstore.DefaultPageSize+5; i++ {
		sm.RevisionId = fmt.Sprintf("page-revision-%d", i)
		if err := store.UpdateStateMachine(ctx, sm); err != nil {
			t.Fatal(err)
		}
		if _, err := store.PublishStateMachineVersion(ctx, smArn, ""); err != nil {
			t.Fatal(err)
		}
	}
	total := sfnstore.DefaultPageSize + 7 // published here plus the two seeded

	first, err := svc.listStateMachineVersionsCore(ctx, store, smArn, 0, "")
	if err != nil {
		t.Fatalf("list versions with maxResults=0: %v", err)
	}
	pageOne, _ := first["stateMachineVersions"].([]map[string]interface{})
	if len(pageOne) != sfnstore.DefaultPageSize {
		t.Fatalf("first page held %d versions, want the default page size %d", len(pageOne), sfnstore.DefaultPageSize)
	}
	nextToken, _ := first["nextToken"].(string)
	if nextToken != strconv.Itoa(sfnstore.DefaultPageSize) {
		t.Fatalf("first page nextToken = %q, want the default-page-size offset", nextToken)
	}

	second, err := svc.listStateMachineVersionsCore(ctx, store, smArn, 0, nextToken)
	if err != nil {
		t.Fatalf("list versions page two: %v", err)
	}
	pageTwo, _ := second["stateMachineVersions"].([]map[string]interface{})
	if len(pageTwo) != total-sfnstore.DefaultPageSize {
		t.Fatalf("second page held %d versions, want %d", len(pageTwo), total-sfnstore.DefaultPageSize)
	}
	if _, has := second["nextToken"]; has {
		t.Fatalf("exhausted listing still pages: %v", second["nextToken"])
	}
}

// TestListStateMachineAliasesMaxResultsZeroListsAll pins the same
// Core-level default for the alias listing: three aliases and maxResults=0
// list everything with no nextToken — an unnormalised zero would emit an
// empty page plus a "0" token.
func TestListStateMachineAliasesMaxResultsZeroListsAll(t *testing.T) {
	store, smArn, _, v1 := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	for _, name := range []string{"p", "q", "r"} {
		if err := store.CreateStateMachineAlias(ctx, &sfnstore.StateMachineAlias{
			StateMachineArn: smArn,
			Name:            name,
			RoutingConfiguration: []sfnstore.RoutingConfiguration{
				{StateMachineVersionArn: v1, Weight: 100},
			},
		}); err != nil {
			t.Fatal(err)
		}
	}

	resp, err := svc.listStateMachineAliasesCore(ctx, store, smArn, 0, "")
	if err != nil {
		t.Fatalf("list aliases with maxResults=0: %v", err)
	}
	aliases, _ := resp["stateMachineAliases"].([]map[string]interface{})
	if len(aliases) != 3 {
		t.Fatalf("listed %d aliases with maxResults=0, want all 3", len(aliases))
	}
	if _, has := resp["nextToken"]; has {
		t.Fatalf("exhausted alias listing still pages: %v", resp["nextToken"])
	}
}

// TestListStateMachineVersionsNewestFirst pins the listing order: "The
// results are sorted in descending order of the version creation time" —
// the key-lexicographic paging (1, 10, 11, 2, ...) must not leak through.
func TestListStateMachineVersionsNewestFirst(t *testing.T) {
	store, smArn, _, _ := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	sm, err := store.GetStateMachine(ctx, smArn)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		sm.RevisionId = fmt.Sprintf("order-revision-%d", i)
		if err := store.UpdateStateMachine(ctx, sm); err != nil {
			t.Fatal(err)
		}
		if _, err := store.PublishStateMachineVersion(ctx, smArn, ""); err != nil {
			t.Fatal(err)
		}
	}

	resp, err := svc.listStateMachineVersionsCore(ctx, store, smArn, 100, "")
	if err != nil {
		t.Fatalf("list versions: %v", err)
	}
	versions, _ := resp["stateMachineVersions"].([]map[string]interface{})
	if len(versions) != 12 {
		t.Fatalf("listed %d versions, want the 12 published", len(versions))
	}
	first, _ := versions[0]["stateMachineVersionArn"].(string)
	if !strings.HasSuffix(first, ":12") {
		t.Errorf("first version = %s, want the newest :12 (descending creation order)", first)
	}
	versionNumber := func(arn string) int {
		n, _ := strconv.Atoi(arn[strings.LastIndex(arn, ":")+1:])
		return n
	}
	for i := 1; i < len(versions); i++ {
		prev, _ := versions[i-1]["stateMachineVersionArn"].(string)
		cur, _ := versions[i]["stateMachineVersionArn"].(string)
		if versionNumber(prev) <= versionNumber(cur) {
			t.Fatalf("order broke at %d: %s then %s (want strictly descending)", i-1, prev, cur)
		}
	}
}

// TestListStateMachineAliasesNewestFirst pins the listing order: "Results
// are sorted by time, with the most recently created aliases listed
// first" — name-ascending storage order must not leak through.
func TestListStateMachineAliasesNewestFirst(t *testing.T) {
	store, smArn, _, v1 := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	for _, name := range []string{"alpha", "beta", "gamma"} {
		if err := store.CreateStateMachineAlias(ctx, &sfnstore.StateMachineAlias{
			StateMachineArn: smArn,
			Name:            name,
			RoutingConfiguration: []sfnstore.RoutingConfiguration{
				{StateMachineVersionArn: v1, Weight: 100},
			},
		}); err != nil {
			t.Fatal(err)
		}
	}

	resp, err := svc.listStateMachineAliasesCore(ctx, store, smArn, 100, "")
	if err != nil {
		t.Fatalf("list aliases: %v", err)
	}
	aliases, _ := resp["stateMachineAliases"].([]map[string]interface{})
	if len(aliases) != 3 {
		t.Fatalf("listed %d aliases, want 3", len(aliases))
	}
	var names []string
	for _, a := range aliases {
		arn, _ := a["stateMachineAliasArn"].(string)
		names = append(names, strings.TrimPrefix(arn, smArn+":"))
	}
	if strings.Join(names, ",") != "gamma,beta,alpha" {
		t.Errorf("alias order = %v, want gamma,beta,alpha (most recent first)", names)
	}

	// Offset paging walks the same newest-first order.
	page, err := svc.listStateMachineAliasesCore(ctx, store, smArn, 2, "2")
	if err != nil {
		t.Fatalf("second page: %v", err)
	}
	tail, _ := page["stateMachineAliases"].([]map[string]interface{})
	if len(tail) != 1 {
		t.Fatalf("second page has %d aliases, want the last one", len(tail))
	}
	if arn, _ := tail[0]["stateMachineAliasArn"].(string); !strings.HasSuffix(arn, ":alpha") {
		t.Errorf("second page = %v, want the oldest alias alpha", tail)
	}
}

// TestUpdateStateMachineAliasClearsDescription pins the Provided-flag
// semantics: an explicitly empty description clears the alias
// description, while an absent description member leaves it untouched.
func TestUpdateStateMachineAliasClearsDescription(t *testing.T) {
	store, smArn, _, v1 := newAliasTestStore(t)
	svc := &StepFunctionService{}
	ctx := context.Background()

	created, err := svc.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:        "described",
		Description: "original text",
		RoutingConfig: []sfnstore.RoutingConfiguration{
			{StateMachineVersionArn: v1, Weight: 100},
		},
	})
	if err != nil {
		t.Fatalf("create alias: %v", err)
	}
	aliasArn, _ := created["stateMachineAliasArn"].(string)
	if aliasArn == "" {
		t.Fatalf("create alias returned %v", created)
	}
	_ = smArn

	if _, err := svc.updateStateMachineAliasCore(ctx, store, UpdateStateMachineAliasInput{
		StateMachineAliasArn: aliasArn,
		Description:          "",
		DescriptionProvided:  true,
	}); err != nil {
		t.Fatalf("clear description: %v", err)
	}
	resp, err := svc.describeStateMachineAliasCore(ctx, store, aliasArn)
	if err != nil {
		t.Fatalf("describe alias: %v", err)
	}
	if _, present := resp["description"]; present {
		t.Errorf("cleared description survives: %v", resp["description"])
	}

	if _, err := svc.updateStateMachineAliasCore(ctx, store, UpdateStateMachineAliasInput{
		StateMachineAliasArn: aliasArn,
		Description:          "restored",
		DescriptionProvided:  true,
	}); err != nil {
		t.Fatalf("restore description: %v", err)
	}
	if _, err := svc.updateStateMachineAliasCore(ctx, store, UpdateStateMachineAliasInput{
		StateMachineAliasArn: aliasArn,
		RoutingConfig: []sfnstore.RoutingConfiguration{
			{StateMachineVersionArn: v1, Weight: 100},
		},
		RoutingProvided: true,
	}); err != nil {
		t.Fatalf("routing-only update: %v", err)
	}
	resp, err = svc.describeStateMachineAliasCore(ctx, store, aliasArn)
	if err != nil {
		t.Fatalf("describe alias again: %v", err)
	}
	if resp["description"] != "restored" {
		t.Errorf("routing-only update disturbed the description: %v", resp["description"])
	}
}
