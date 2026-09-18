package cloudtrail

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	tags "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// The pins below exercise the model-derived error vocabulary: every wire
// code, HTTP status, and query-error code comes from the vendored Smithy
// model, and the store sentinels reach the API surface through
// storeErrorMappings — never through inline literals.

func newErrorTestServiceStore(t *testing.T) *cloudtrailstore.CloudTrailStore {
	t.Helper()
	tmpDir := "./tmp/cloudtrail-errors-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	s, err := storage.Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open test storage: %v", err)
	}
	t.Cleanup(func() { s.Close() })

	return cloudtrailstore.NewCloudTrailStore(s, "acc123", "us-east-1")
}

// TestCloudTrailErrorVocabularyModelAlignment pins the wire code, HTTP
// status, and awsQueryError code of every produced error shape against the
// values the Smithy model declares. The not-found family carries 404 from
// the model's httpError trait, and the already-exists family 400.
func TestCloudTrailErrorVocabularyModelAlignment(t *testing.T) {
	cases := []struct {
		err       *awserrors.AWSError
		code      string
		status    int
		queryCode string
	}{
		{ErrTrailNotFound, "TrailNotFoundException", 404, "TrailNotFound"},
		{ErrTrailAlreadyExists, "TrailAlreadyExistsException", 400, "TrailAlreadyExists"},
		{ErrInvalidParameter, "InvalidParameterException", 400, "InvalidParameter"},
		{ErrInvalidTrailName, "InvalidTrailNameException", 400, "InvalidTrailName"},
		{ErrEventDataStoreNotFoundException, "EventDataStoreNotFoundException", 404, "EventDataStoreNotFound"},
		{ErrEventDataStoreAlreadyExists, "EventDataStoreAlreadyExistsException", 400, "EventDataStoreAlreadyExists"},
		{ErrChannelNotFound, "ChannelNotFoundException", 404, "ChannelNotFound"},
		{ErrChannelAlreadyExists, "ChannelAlreadyExistsException", 400, "ChannelAlreadyExists"},
		{ErrQueryIdNotFound, "QueryIdNotFoundException", 404, "QueryIdNotFound"},
		{ErrImportNotFound, "ImportNotFoundException", 404, "ImportNotFound"},
		{ErrResourcePolicyNotFound, "ResourcePolicyNotFoundException", 404, "ResourcePolicyNotFound"},
		{ErrOperationNotPermitted, "OperationNotPermittedException", 400, "OperationNotPermitted"},
		{ErrInvalidEventCategory, "InvalidEventCategoryException", 400, "InvalidEventCategory"},
		{newInvalidQueryStatementException("x"), "InvalidQueryStatementException", 400, "InvalidQueryStatement"},
		{newInvalidQueryStatusException("x"), "InvalidQueryStatusException", 400, "InvalidQueryStatus"},
		{newInvalidTagParameterException("x"), "InvalidTagParameterException", 400, "InvalidTagParameter"},
		{newInvalidImportSourceException("x"), "InvalidImportSourceException", 400, "InvalidImportSource"},
		{newInactiveEventDataStoreException("x"), "InactiveEventDataStoreException", 400, "InactiveEventDataStore"},
		{newInvalidLookupAttributesException("x"), "InvalidLookupAttributesException", 400, "InvalidLookupAttributes"},
		{newInvalidMaxResultsException("x"), "InvalidMaxResultsException", 400, "InvalidMaxResults"},
		{newInvalidTimeRangeException("x"), "InvalidTimeRangeException", 400, "InvalidTimeRange"},
		{newInvalidNextTokenException("x"), "InvalidNextTokenException", 400, "InvalidNextToken"},
		{newMaximumNumberOfTrailsExceededException("x"), "MaximumNumberOfTrailsExceededException", 403, "MaximumNumberOfTrailsExceeded"},
		{newEventDataStoreMaxLimitExceededException("x"), "EventDataStoreMaxLimitExceededException", 400, "EventDataStoreMaxLimitExceeded"},
		{newChannelMaxLimitExceededException("x"), "ChannelMaxLimitExceededException", 400, "ChannelMaxLimitExceeded"},
		{newOrganizationsNotInUseException("x"), "OrganizationsNotInUseException", 404, "OrganizationsNotInUse"},
		{ErrInternalError, "InternalFailure", 500, ""},
	}
	for _, tc := range cases {
		if got := tc.err.GetCode(); got != tc.code {
			t.Errorf("%s: code = %q, want %q", tc.code, got, tc.code)
		}
		if got := tc.err.GetHTTPStatusCode(); got != tc.status {
			t.Errorf("%s: status = %d, want %d", tc.code, got, tc.status)
		}
		if got := tc.err.QueryErrorCode; got != tc.queryCode {
			t.Errorf("%s: query code = %q, want %q", tc.code, got, tc.queryCode)
		}
	}
}

// TestDelegatedAdminAnswersOrganizationsNotInUse pins the delegated-
// administration pair: an account on this platform never belongs to an
// organization, so both operations answer the model's non-member refusal
// after their member validation.
func TestDelegatedAdminAnswersOrganizationsNotInUse(t *testing.T) {
	svc := &CloudTrailService{}

	_, err := svc.registerOrganizationDelegatedAdminCore(RegisterOrganizationDelegatedAdminInput{MemberAccountID: "acc123"})
	requireAWSCode(t, err, "OrganizationsNotInUseException", 404)

	err = svc.deregisterOrganizationDelegatedAdminCore(DeregisterOrganizationDelegatedAdminInput{DelegatedAdminAccountID: "acc123"})
	requireAWSCode(t, err, "OrganizationsNotInUseException", 404)

	if _, err := svc.registerOrganizationDelegatedAdminCore(RegisterOrganizationDelegatedAdminInput{}); err == nil {
		t.Fatal("expected a missing MemberAccountId to be rejected")
	}
}

// TestMapStoreErrorMasksUnknownFailuresAsInternal pins the fallthrough:
// a store error with no mapping entry must surface as InternalFailure,
// never as a raw store error.
func TestMapStoreErrorMasksUnknownFailuresAsInternal(t *testing.T) {
	svc := &CloudTrailService{}
	if err := svc.mapStoreError(errors.New("boom")); err != ErrInternalError {
		t.Fatalf("expected ErrInternalError for an unmapped store error, got %v", err)
	}
	if svc.mapStoreError(nil) != nil {
		t.Fatal("expected nil for a nil store error")
	}
}

// TestCoreNotFoundPathsPinMappedCodes drives the not-found paths through
// the real core functions against a real store: each store sentinel must
// arrive as its mapped, model-declared API error.
func TestCoreNotFoundPathsPinMappedCodes(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	cases := []struct {
		name string
		call func() error
		code string
	}{
		{"missing channel", func() error {
			_, err := svc.getChannelCore(store, ChannelInput{Channel: "arn:aws:cloudtrail:us-east-1:acc123:channel/nope"})
			return err
		}, "ChannelNotFoundException"},
		{"missing event data store", func() error {
			_, err := svc.getEventDataStoreCore(store, EventDataStoreIDInput{EventDataStore: "nope"})
			return err
		}, "EventDataStoreNotFoundException"},
		{"missing query", func() error {
			_, err := svc.getQueryResultsCore(store, GetQueryResultsInput{QueryID: "nope"})
			return err
		}, "QueryIdNotFoundException"},
		{"missing import", func() error {
			_, err := svc.getImportCore(store, ImportIDInput{ImportID: "nope"})
			return err
		}, "ImportNotFoundException"},
	}
	for _, tc := range cases {
		err := tc.call()
		if err == nil {
			t.Fatalf("%s: expected an error", tc.name)
		}
		if got := err.(*awserrors.AWSError).GetCode(); got != tc.code {
			t.Errorf("%s: code = %q, want %q", tc.name, got, tc.code)
		}
	}
}

// TestChannelNameUniquenessPinsRegisterSeven pins the wired
// ChannelAlreadyExistsException: creating a channel whose name is already
// taken, and renaming a channel onto another channel's name, are both
// rejected; keeping one's own name stays allowed. Each channel carries its
// own source — a maximum of one channel exists per source.
func TestChannelNameUniquenessPinsRegisterSeven(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	// The destination must resolve to a stored event data store.
	destEDS, err := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore("uniqueness-dest", store.GetAccountID(), store.GetRegion()))
	if err != nil {
		t.Fatalf("create destination EDS failed: %v", err)
	}
	dest := []interface{}{map[string]interface{}{
		"Type":     "EVENT_DATA_STORE",
		"Location": destEDS.EventDataStoreARN,
	}}

	first, err := svc.createChannelCore(store, CreateChannelInput{Name: "dup-name", Source: "Custom", DestinationsRaw: dest, DestinationsSet: true})
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err = svc.createChannelCore(store, CreateChannelInput{Name: "dup-name", Source: "other-source", DestinationsRaw: dest, DestinationsSet: true})
	if err == nil {
		t.Fatal("expected duplicate-name create to be rejected")
	}
	if got := err.(*awserrors.AWSError).GetCode(); got != "ChannelAlreadyExistsException" {
		t.Fatalf("duplicate create: code = %q, want ChannelAlreadyExistsException", got)
	}
	if got := err.(*awserrors.AWSError).GetHTTPStatusCode(); got != 400 {
		t.Fatalf("duplicate create: status = %d, want 400", got)
	}

	// A source already carried by another channel is rejected: "A maximum
	// of one channel is allowed per source" (CreateChannel Source).
	_, err = svc.createChannelCore(store, CreateChannelInput{Name: "source-dup", Source: "Custom", DestinationsRaw: dest, DestinationsSet: true})
	if err == nil {
		t.Fatal("expected duplicate-source create to be rejected")
	}
	if got := err.(*awserrors.AWSError).GetCode(); got != "InvalidSourceException" {
		t.Fatalf("duplicate source: code = %q, want InvalidSourceException", got)
	}

	second, err := svc.createChannelCore(store, CreateChannelInput{Name: "other-name", Source: "third-source", DestinationsRaw: dest, DestinationsSet: true})
	if err != nil {
		t.Fatalf("second create failed: %v", err)
	}

	_, err = svc.updateChannelCore(store, UpdateChannelInput{Channel: first["ChannelArn"].(string), Name: "other-name"})
	if err == nil {
		t.Fatal("expected rename onto a taken name to be rejected")
	}
	if got := err.(*awserrors.AWSError).GetCode(); got != "ChannelAlreadyExistsException" {
		t.Fatalf("rename collision: code = %q, want ChannelAlreadyExistsException", got)
	}

	// Keeping the channel's own name is not a collision.
	if _, err := svc.updateChannelCore(store, UpdateChannelInput{Channel: second["ChannelArn"].(string), Name: "other-name"}); err != nil {
		t.Fatalf("self-name update rejected: %v", err)
	}
}

// TestDuplicateCreateStatusIsBadRequest pins the register-one status fix:
// the model declares TrailAlreadyExistsException and
// EventDataStoreAlreadyExistsException at HTTP 400, not 409.
func TestDuplicateCreateStatusIsBadRequest(t *testing.T) {
	ctx := context.Background()
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	if _, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "dup-trail", S3BucketName: "dup-trail-bucket", Region: "us-east-1"}); err != nil {
		t.Fatalf("first trail create failed: %v", err)
	}
	_, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "dup-trail", S3BucketName: "dup-trail-bucket", Region: "us-east-1"})
	if err == nil {
		t.Fatal("expected duplicate trail create to be rejected")
	}
	if got := err.(*awserrors.AWSError).GetCode(); got != "TrailAlreadyExistsException" {
		t.Fatalf("duplicate trail: code = %q, want TrailAlreadyExistsException", got)
	}
	if got := err.(*awserrors.AWSError).GetHTTPStatusCode(); got != 400 {
		t.Fatalf("duplicate trail: status = %d, want 400", got)
	}

	if _, err := svc.createEventDataStoreCore(store, CreateEventDataStoreInput{Name: "dup-eds"}); err != nil {
		t.Fatalf("first EDS create failed: %v", err)
	}
	_, err = svc.createEventDataStoreCore(store, CreateEventDataStoreInput{Name: "dup-eds"})
	if err == nil {
		t.Fatal("expected duplicate EDS create to be rejected")
	}
	if got := err.(*awserrors.AWSError).GetCode(); got != "EventDataStoreAlreadyExistsException" {
		t.Fatalf("duplicate EDS: code = %q, want EventDataStoreAlreadyExistsException", got)
	}
	if got := err.(*awserrors.AWSError).GetHTTPStatusCode(); got != 400 {
		t.Fatalf("duplicate EDS: status = %d, want 400", got)
	}
}

// TestTrailEmptyNamePathsUseDeclaredCode pins the wrong-code-for-op fix:
// the single-trail operations declare InvalidTrailNameException, so an
// omitted name answers with it rather than the generic parameter error.
func TestTrailEmptyNamePathsUseDeclaredCode(t *testing.T) {
	// The empty-name rejections fire before any store or destination
	// access, so the zero service suffices.
	svc := &CloudTrailService{}

	calls := map[string]func() error{
		"updateTrailCore":  func() error { _, err := svc.updateTrailCore(nil, nil, UpdateTrailInput{}); return err },
		"resolveTrailCore": func() error { _, err := svc.resolveTrailCore(nil, ""); return err },
		"startLoggingCore": func() error { return svc.startLoggingCore(nil, TrailNameInput{}) },
		"stopLoggingCore":  func() error { return svc.stopLoggingCore(nil, TrailNameInput{}) },
	}
	for name, call := range calls {
		err := call()
		if err == nil {
			t.Fatalf("%s: expected an error", name)
		}
		awsErr, ok := err.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("%s: expected *awserrors.AWSError, got %T", name, err)
		}
		if got := awsErr.GetCode(); got != "InvalidTrailNameException" {
			t.Errorf("%s: code = %q, want InvalidTrailNameException", name, got)
		}
	}
}

// TestQueryAndTagValidationCodes pins the declared-error swaps on the
// existing validation paths: malformed statements answer
// InvalidQueryStatementException, a bad ListQueries status filter
// InvalidQueryStatusException, and every tag-validation failure the
// model's InvalidTagParameterException.
func TestQueryAndTagValidationCodes(t *testing.T) {
	if _, err := parseQueryStatement(""); err == nil {
		t.Fatal("expected empty statement to be rejected")
	} else if got := err.(*awserrors.AWSError).GetCode(); got != "InvalidQueryStatementException" {
		t.Fatalf("empty statement: code = %q, want InvalidQueryStatementException", got)
	}

	if _, err := parseQueryStatement("not a select statement"); err == nil {
		t.Fatal("expected malformed statement to be rejected")
	} else if got := err.(*awserrors.AWSError).GetCode(); got != "InvalidQueryStatementException" {
		t.Fatalf("malformed statement: code = %q, want InvalidQueryStatementException", got)
	}

	// The QueryStatement shape is @length 1-10000; a longer statement is
	// rejected before parsing (the typed SDK blocks it client-side, so this
	// guards the raw-HTTP path).
	long := "SELECT eventID FROM " + strings.Repeat("a", cloudtrailstore.MaxQueryStatementChars)
	if _, err := parseQueryStatement(long); err == nil {
		t.Fatal("expected over-length statement to be rejected")
	} else if got := err.(*awserrors.AWSError).GetCode(); got != "InvalidQueryStatementException" {
		t.Fatalf("over-length statement: code = %q, want InvalidQueryStatementException", got)
	}

	// The length counts characters, not bytes: a statement at the exact
	// character cap whose UTF-8 encoding exceeds it in bytes parses.
	prefix := "SELECT eventID FROM "
	multibyte := prefix + strings.Repeat("あ", cloudtrailstore.MaxQueryStatementChars-len(prefix))
	if utf8.RuneCountInString(multibyte) != cloudtrailstore.MaxQueryStatementChars {
		t.Fatalf("fixture must be exactly %d characters", cloudtrailstore.MaxQueryStatementChars)
	}
	if len(multibyte) <= cloudtrailstore.MaxQueryStatementChars {
		t.Fatal("fixture must exceed the cap in bytes")
	}
	if _, err := parseQueryStatement(multibyte); err != nil {
		t.Fatalf("character-capped statement must not be rejected on its byte length: %v", err)
	}

	if err := validateQueryStatus("BOGUS"); err == nil {
		t.Fatal("expected invalid query status to be rejected")
	} else if got := err.(*awserrors.AWSError).GetCode(); got != "InvalidQueryStatusException" {
		t.Fatalf("invalid status: code = %q, want InvalidQueryStatusException", got)
	}

	tagCases := map[string][]tags.Tag{
		"reserved key": {{Key: "aws:reserved", Value: "x"}},
		"key too long": {{Key: string(make([]byte, 129)), Value: "x"}},
	}
	for name, tagList := range tagCases {
		err := validateCloudTrailTags(tagList)
		if err == nil {
			t.Fatalf("%s: expected an error", name)
		}
		if got := err.(*awserrors.AWSError).GetCode(); got != "InvalidTagParameterException" {
			t.Errorf("%s: code = %q, want InvalidTagParameterException", name, got)
		}
	}
}

// TestEventConfigurationAbsentPathUsesDeclaredCodes pins the invented
// ConfigurationException removal: an absent event configuration answers
// with the selected resource type's declared not-found error, and a
// stored configuration still round-trips.
func TestEventConfigurationAbsentPathUsesDeclaredCodes(t *testing.T) {
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	_, err := svc.getEventConfigurationCore(store, EventConfigurationResourceInput{TrailName: "never-configured"})
	if err == nil {
		t.Fatal("expected absent trail configuration to be rejected")
	}
	if got := err.(*awserrors.AWSError).GetCode(); got != "TrailNotFoundException" {
		t.Fatalf("absent trail config: code = %q, want TrailNotFoundException", got)
	}

	_, err = svc.getEventConfigurationCore(store, EventConfigurationResourceInput{EventDataStore: "never-configured"})
	if err == nil {
		t.Fatal("expected absent EDS configuration to be rejected")
	}
	if got := err.(*awserrors.AWSError).GetCode(); got != "EventDataStoreNotFoundException" {
		t.Fatalf("absent EDS config: code = %q, want EventDataStoreNotFoundException", got)
	}

	trail, err := svc.createTrailCore(context.Background(), store, CreateTrailInput{Name: "configured-trail", S3BucketName: "configured-trail-bucket", Region: "us-east-1"})
	if err != nil {
		t.Fatalf("trail create failed: %v", err)
	}
	if _, err := svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName: "configured-trail",
		Params:    map[string]interface{}{"MaxEventSize": "Standard"},
	}); err != nil {
		t.Fatalf("put event configuration failed: %v", err)
	}
	_, err = svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName: "never-created",
	})
	if err == nil {
		t.Fatal("expected put for an absent trail to be rejected")
	}
	if got := err.(*awserrors.AWSError).GetCode(); got != "TrailNotFoundException" {
		t.Fatalf("absent trail put: code = %q, want TrailNotFoundException", got)
	}
	config, err := svc.getEventConfigurationCore(store, EventConfigurationResourceInput{TrailName: "configured-trail"})
	if err != nil {
		t.Fatalf("get event configuration failed: %v", err)
	}
	if config["MaxEventSize"] != "Standard" {
		t.Fatalf("expected MaxEventSize=Standard to round-trip, got %v", config["MaxEventSize"])
	}
	if config["TrailARN"] != trail.TrailARN {
		t.Fatalf("expected TrailARN=%q in the stored configuration, got %v", trail.TrailARN, config["TrailARN"])
	}
}

// TestRestoreNonPendingEDSAnswersOperationNotPermitted pins the mapping of
// the store's not-pending-deletion sentinel through
// storeErrorMappings.
func TestRestoreNonPendingEDSAnswersOperationNotPermitted(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	created, err := svc.createEventDataStoreCore(store, CreateEventDataStoreInput{Name: "active-eds"})
	if err != nil {
		t.Fatalf("EDS create failed: %v", err)
	}
	_, err = svc.restoreEventDataStoreCore(store, EventDataStoreIDInput{EventDataStore: created["EventDataStoreArn"].(string)})
	if err == nil {
		t.Fatal("expected restore of a non-pending EDS to be rejected")
	}
	if got := err.(*awserrors.AWSError).GetCode(); got != "OperationNotPermittedException" {
		t.Fatalf("restore non-pending: code = %q, want OperationNotPermittedException", got)
	}
}
