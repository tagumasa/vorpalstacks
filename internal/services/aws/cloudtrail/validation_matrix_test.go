package cloudtrail

import (
	"context"
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// The pins below cover the D2 validation matrix: every constraint the Core
// functions enforce is pinned with its model-declared error identity, and
// the boundary-adjacent valid values are pinned as accepted so the negative
// edge is exact.

func requireAWSCodeFromMatrix(t *testing.T, name string, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected %s, got success", name, code)
	}
	apiErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("%s: error %v is not an AWSError", name, err)
	}
	if apiErr.GetCode() != code {
		t.Fatalf("%s: code = %q, want %q", name, apiErr.GetCode(), code)
	}
}

func TestTrailNameRules(t *testing.T) {
	valid := []string{"my-trail", "My.Trail_1-x", "t12", "trail.name.2"}
	for _, name := range valid {
		if err := validateTrailName(name); err != nil {
			t.Errorf("validateTrailName(%q) rejected a valid name: %v", name, err)
		}
	}
	invalid := map[string]string{
		"my--trail":   "adjacent dashes",
		"my._trail":   "adjacent separators",
		"trail-":      "trailing separator",
		"-trail":      "leading separator",
		"1.2.3.4":     "IP address format",
		"192.168.5.4": "IP address format",
		"ab":          "too short",
		"a!b":         "invalid character",
	}
	for name, why := range invalid {
		err := validateTrailName(name)
		apiErr, ok := err.(*awserrors.AWSError)
		if !ok || apiErr.GetCode() != "InvalidTrailNameException" {
			t.Errorf("validateTrailName(%q) (%s) = %v, want InvalidTrailNameException", name, why, err)
		}
	}
	long := ""
	for i := 0; i < 129; i++ {
		long += "a"
	}
	if err := validateTrailName(long); err == nil {
		t.Error("validateTrailName accepted a 129-character name")
	}
}

func TestKMSKeyIDForms(t *testing.T) {
	valid := []string{
		"12345678-1234-1234-1234-123456789012",
		"arn:aws:kms:us-east-2:123456789012:key/12345678-1234-1234-1234-123456789012",
		"alias/MyAliasName",
		"arn:aws:kms:us-east-2:123456789012:alias/MyAliasName",
	}
	for _, id := range valid {
		if err := validateKMSKeyID(id); err != nil {
			t.Errorf("validateKMSKeyID(%q) rejected a documented form: %v", id, err)
		}
		if err := validateEventDataStoreKMSKeyID(id); err != nil {
			t.Errorf("validateEventDataStoreKMSKeyID(%q) rejected a documented form: %v", id, err)
		}
	}
	for _, id := range []string{"", "not-a-key", "arn:aws:kms:us-east-2:123456789012:key/not-a-uuid"} {
		if err := validateKMSKeyID(id); err == nil {
			t.Errorf("validateKMSKeyID(%q) accepted an invalid form", id)
		}
	}
	oversize := "k"
	for i := 0; i < 350; i++ {
		oversize += "k"
	}
	requireAWSCodeFromMatrix(t, "EDS KMS oversize", validateEventDataStoreKMSKeyID(oversize), "InvalidKmsKeyIdException")
}

func TestUpdateTrailRoutesThroughCreateValidators(t *testing.T) {
	ctx := context.Background()
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	created, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "matrix-trail", S3BucketName: "matrix-trail-bucket", Region: "us-east-1"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	cases := []struct {
		name   string
		params map[string]interface{}
		code   string
	}{
		{"bad bucket", map[string]interface{}{"S3BucketName": "..invalid.."}, "InvalidS3BucketNameException"},
		{"empty bucket", map[string]interface{}{"S3BucketName": ""}, "InvalidS3BucketNameException"},
		{"bad kms", map[string]interface{}{"KmsKeyId": "garbage"}, "InvalidKmsKeyIdException"},
		{"bad cwlogs group", map[string]interface{}{"CloudWatchLogsLogGroupArn": "not-an-arn"}, "InvalidCloudWatchLogsLogGroupArnException"},
		{"bad cwlogs role", map[string]interface{}{"CloudWatchLogsRoleArn": "not-an-arn"}, "InvalidCloudWatchLogsRoleArnException"},
		{"bad sns topic", map[string]interface{}{"SnsTopicName": "bad topic!"}, "InvalidSnsTopicNameException"},
	}
	for _, tc := range cases {
		_, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{Name: created.Name, Params: tc.params})
		requireAWSCodeFromMatrix(t, tc.name, err, tc.code)
	}

	// A valid update through the same validators still succeeds, and the
	// alias form of KmsKeyId is accepted.
	updated, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{Name: created.Name, Params: map[string]interface{}{
		"KmsKeyId": "alias/MatrixKey",
	}})
	if err != nil {
		t.Fatalf("valid update rejected: %v", err)
	}
	if updated.KMSKeyID != "alias/MatrixKey" {
		t.Fatalf("update did not apply the KMS key: %q", updated.KMSKeyID)
	}
}

func TestUpdateEventDataStoreValidationMatrix(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	created, err := svc.createEventDataStoreCore(store, CreateEventDataStoreInput{Name: "matrix-eds", KmsKeyId: "alias/k1"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	arn, _ := created["EventDataStoreArn"].(string)

	// "At least one optional parameter must be specified."
	_, err = svc.updateEventDataStoreCore(store, UpdateEventDataStoreInput{EventDataStore: arn})
	requireAWSCodeFromMatrix(t, "no parameters", err, "InvalidParameterException")

	// The KMS key cannot be changed once associated; re-asserting it stays
	// a no-op success.
	_, err = svc.updateEventDataStoreCore(store, UpdateEventDataStoreInput{EventDataStore: arn, KmsKeyId: "alias/k2"})
	requireAWSCodeFromMatrix(t, "kms change", err, "OperationNotPermittedException")
	if _, err = svc.updateEventDataStoreCore(store, UpdateEventDataStoreInput{EventDataStore: arn, KmsKeyId: "alias/k1"}); err != nil {
		t.Fatalf("kms re-assert rejected: %v", err)
	}

	_, err = svc.updateEventDataStoreCore(store, UpdateEventDataStoreInput{EventDataStore: arn, KmsKeyId: "bad key!"})
	requireAWSCodeFromMatrix(t, "kms form", err, "InvalidKmsKeyIdException")

	// The billing mode cannot change from EXTENDABLE to FIXED (the created
	// store defaults to EXTENDABLE_RETENTION_PRICING).
	_, err = svc.updateEventDataStoreCore(store, UpdateEventDataStoreInput{EventDataStore: arn, BillingMode: "FIXED_RETENTION_PRICING"})
	requireAWSCodeFromMatrix(t, "billing mode change", err, "OperationNotPermittedException")

	// A retention period above the FIXED_RETENTION_PRICING cap is
	// rejected on both create and update.
	_, err = svc.createEventDataStoreCore(store, CreateEventDataStoreInput{
		Name: "fixed-eds", BillingMode: "FIXED_RETENTION_PRICING",
		RetentionPeriodRaw: float64(3000), RetentionPeriodSet: true,
	})
	requireAWSCodeFromMatrix(t, "fixed retention create", err, "InvalidParameterException")
	_, err = svc.updateEventDataStoreCore(store, UpdateEventDataStoreInput{
		EventDataStore: arn, BillingMode: "EXTENDABLE_RETENTION_PRICING",
		RetentionPeriodRaw: float64(3000), RetentionPeriodSet: true,
	})
	if err != nil {
		t.Fatalf("extendable retention 3000 rejected: %v", err)
	}

	// A PENDING_DELETION store is inactive and rejects updates.
	dying, err := svc.createEventDataStoreCore(store, CreateEventDataStoreInput{Name: "dying-eds"})
	if err != nil {
		t.Fatalf("second create failed: %v", err)
	}
	noProtection := false
	_, err = svc.updateEventDataStoreCore(store, UpdateEventDataStoreInput{
		EventDataStore: dying["EventDataStoreArn"].(string), TerminationProtectionEnabled: &noProtection,
	})
	if err != nil {
		t.Fatalf("protection disable failed: %v", err)
	}
	if _, err = svc.deleteEventDataStoreCore(store, EventDataStoreIDInput{EventDataStore: dying["EventDataStoreArn"].(string)}); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	_, err = svc.updateEventDataStoreCore(store, UpdateEventDataStoreInput{
		EventDataStore: dying["EventDataStoreArn"].(string), Name: "renamed-dying",
	})
	requireAWSCodeFromMatrix(t, "pending deletion update", err, "InactiveEventDataStoreException")
}

func TestChannelValidationMatrix(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	// Destinations must resolve to stored event data stores, so the valid
	// create path targets a real one.
	destEDS, err := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore("channel-dest", store.GetAccountID(), store.GetRegion()))
	if err != nil {
		t.Fatalf("create destination EDS failed: %v", err)
	}
	dest := []interface{}{map[string]interface{}{"Type": "EVENT_DATA_STORE", "Location": destEDS.EventDataStoreARN}}

	_, err = svc.createChannelCore(store, CreateChannelInput{Name: "ab", Source: "Custom", DestinationsRaw: dest, DestinationsSet: true})
	requireAWSCodeFromMatrix(t, "short channel name", err, "InvalidParameterException")

	_, err = svc.createChannelCore(store, CreateChannelInput{Name: "valid-name", Source: "", DestinationsRaw: dest, DestinationsSet: true})
	requireAWSCodeFromMatrix(t, "empty source", err, "InvalidSourceException")

	_, err = svc.createChannelCore(store, CreateChannelInput{Name: "valid-name", Source: "Custom"})
	requireAWSCodeFromMatrix(t, "missing destinations", err, "InvalidParameterException")

	badType := []interface{}{map[string]interface{}{"Type": "KINESIS", "Location": "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/dest-a"}}
	_, err = svc.createChannelCore(store, CreateChannelInput{Name: "valid-name", Source: "Custom", DestinationsRaw: badType, DestinationsSet: true})
	requireAWSCodeFromMatrix(t, "bad destination type", err, "InvalidParameterException")

	badLocation := []interface{}{map[string]interface{}{"Type": "AWS_SERVICE", "Location": "x"}}
	_, err = svc.createChannelCore(store, CreateChannelInput{Name: "valid-name", Source: "Custom", DestinationsRaw: badLocation, DestinationsSet: true})
	requireAWSCodeFromMatrix(t, "bad destination location", err, "InvalidParameterException")

	created, err := svc.createChannelCore(store, CreateChannelInput{Name: "valid-name", Source: "Custom", DestinationsRaw: dest, DestinationsSet: true})
	if err != nil {
		t.Fatalf("valid create rejected: %v", err)
	}

	// "A maximum of one channel is allowed per source."
	_, err = svc.createChannelCore(store, CreateChannelInput{Name: "second-name", Source: "Custom", DestinationsRaw: dest, DestinationsSet: true})
	requireAWSCodeFromMatrix(t, "duplicate source", err, "InvalidSourceException")

	badRename := []interface{}{map[string]interface{}{"Type": "AWS_SERVICE", "Location": "s"}}
	_, err = svc.updateChannelCore(store, UpdateChannelInput{Channel: created["ChannelArn"].(string), DestinationsRaw: badRename, DestinationsSet: true})
	requireAWSCodeFromMatrix(t, "update bad destination", err, "InvalidParameterException")
}

func advancedSelectorRaw(fields ...map[string]interface{}) []interface{} {
	list := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		if f == nil {
			list = append(list, map[string]interface{}{"Name": "empty"})
		} else {
			list = append(list, map[string]interface{}{"Name": "sel", "FieldSelectors": []interface{}{f}})
		}
	}
	return list
}

func TestPutEventSelectorsValidationMatrix(t *testing.T) {
	ctx := context.Background()
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	created, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "sel-trail", S3BucketName: "sel-trail-bucket", Region: "us-east-1"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	basic := []interface{}{map[string]interface{}{"ReadWriteType": "All", "IncludeManagementEvents": true}}
	advanced := advancedSelectorRaw(map[string]interface{}{"Field": "eventCategory", "Equals": []interface{}{"Management"}})

	cases := []struct {
		name     string
		basic    interface{}
		advanced interface{}
		code     string
	}{
		{"both forms", basic, advanced, "InvalidEventSelectorsException"},
		{"neither form", nil, nil, "InvalidEventSelectorsException"},
		{"bad read write type", []interface{}{map[string]interface{}{"ReadWriteType": "read-only"}}, nil, "InvalidEventSelectorsException"},
		{"bad data resource type", []interface{}{map[string]interface{}{
			"ReadWriteType": "All",
			"DataResources": []interface{}{map[string]interface{}{"Type": "AWS::Sqs::Queue", "Values": []interface{}{"arn:aws:sqs:::q"}}},
		}}, nil, "InvalidEventSelectorsException"},
		{"advanced without field selectors", nil, []interface{}{map[string]interface{}{"Name": "empty"}}, "InvalidEventSelectorsException"},
		{"field selector without field", nil, advancedSelectorRaw(map[string]interface{}{"Equals": []interface{}{"Management"}}), "InvalidEventSelectorsException"},
		{"unknown advanced field", nil, advancedSelectorRaw(map[string]interface{}{
			"Field": "notAField", "Equals": []interface{}{"v"},
		}), "InvalidEventSelectorsException"},
		{"readOnly deselect operator", nil, advancedSelectorRaw(map[string]interface{}{
			"Field": "readOnly", "NotEquals": []interface{}{"true"},
		}), "InvalidEventSelectorsException"},
		{"eventCategory prefix operator", nil, advancedSelectorRaw(map[string]interface{}{
			"Field": "eventCategory", "StartsWith": []interface{}{"Man"},
		}), "InvalidEventSelectorsException"},
		{"resources.type suffix operator", nil, advancedSelectorRaw(map[string]interface{}{
			"Field": "resources.type", "EndsWith": []interface{}{"Object"},
		}), "InvalidEventSelectorsException"},
		{"console field prefix operator", nil, advancedSelectorRaw(map[string]interface{}{
			"Field": "sessionCredentialFromConsole", "StartsWith": []interface{}{"t"},
		}), "InvalidEventSelectorsException"},
		{"errorCode non-equals operator", nil, advancedSelectorRaw(map[string]interface{}{
			"Field": "errorCode", "NotEquals": []interface{}{"VpceAccessDenied"},
		}), "InvalidEventSelectorsException"},
		{"network activity eventSource non-equals operator", nil, []interface{}{map[string]interface{}{
			"FieldSelectors": []interface{}{
				map[string]interface{}{"Field": "eventCategory", "Equals": []interface{}{"NetworkActivity"}},
				map[string]interface{}{"Field": "eventSource", "NotEquals": []interface{}{"s3.amazonaws.com"}},
			},
		}}, "InvalidEventSelectorsException"},
		{"network activity without eventSource", nil, []interface{}{map[string]interface{}{
			"FieldSelectors": []interface{}{
				map[string]interface{}{"Field": "eventCategory", "Equals": []interface{}{"NetworkActivity"}},
				map[string]interface{}{"Field": "eventName", "Equals": []interface{}{"GetObject"}},
			},
		}}, "InvalidEventSelectorsException"},
		{"field selector without operator values", nil, advancedSelectorRaw(map[string]interface{}{
			"Field": "eventName",
		}), "InvalidEventSelectorsException"},
		{"field selector with an empty operator list", nil, advancedSelectorRaw(map[string]interface{}{
			"Field": "eventName", "Equals": []interface{}{},
		}), "InvalidEventSelectorsException"},
	}
	for _, tc := range cases {
		_, err := svc.putEventSelectorsCore(store, PutEventSelectorsInput{
			TrailName:                 created.Name,
			EventSelectorsRaw:         tc.basic,
			AdvancedEventSelectorsRaw: tc.advanced,
		})
		requireAWSCodeFromMatrix(t, tc.name, err, tc.code)
	}

	// Six basic selectors exceed the five-selector bound.
	six := make([]interface{}, 0, 6)
	for i := 0; i < 6; i++ {
		six = append(six, map[string]interface{}{"ReadWriteType": "All"})
	}
	_, err = svc.putEventSelectorsCore(store, PutEventSelectorsInput{TrailName: created.Name, EventSelectorsRaw: six})
	requireAWSCodeFromMatrix(t, "six basic selectors", err, "InvalidEventSelectorsException")

	// Six advanced selectors exceed the bound as well.
	sixAdvanced := make([]interface{}, 0, 6)
	for i := 0; i < 6; i++ {
		sixAdvanced = append(sixAdvanced, map[string]interface{}{
			"FieldSelectors": []interface{}{map[string]interface{}{"Field": "eventCategory", "Equals": []interface{}{"Management"}}},
		})
	}
	_, err = svc.putEventSelectorsCore(store, PutEventSelectorsInput{TrailName: created.Name, AdvancedEventSelectorsRaw: sixAdvanced})
	requireAWSCodeFromMatrix(t, "six advanced selectors", err, "InvalidEventSelectorsException")

	// 501 condition values exceed the 500-value trail bound: two selectors
	// with 251 Equals values each plus the valid single-value ones already
	// pinned above.
	manyValues := make([]interface{}, 0, 501)
	for i := 0; i < 501; i++ {
		manyValues = append(manyValues, fmt.Sprintf("v%d", i))
	}
	over := []interface{}{map[string]interface{}{
		"FieldSelectors": []interface{}{map[string]interface{}{"Field": "eventName", "Equals": manyValues}},
	}}
	_, err = svc.putEventSelectorsCore(store, PutEventSelectorsInput{TrailName: created.Name, AdvancedEventSelectorsRaw: over})
	requireAWSCodeFromMatrix(t, "501 advanced values", err, "InvalidEventSelectorsException")

	// One form at a time still succeeds, and the later form overwrites the
	// earlier one.
	if _, err = svc.putEventSelectorsCore(store, PutEventSelectorsInput{TrailName: created.Name, EventSelectorsRaw: basic}); err != nil {
		t.Fatalf("basic form rejected: %v", err)
	}
	if _, err = svc.putEventSelectorsCore(store, PutEventSelectorsInput{TrailName: created.Name, AdvancedEventSelectorsRaw: advanced}); err != nil {
		t.Fatalf("advanced form rejected: %v", err)
	}

	// A deselect-only selector on an unrestricted field is valid — the form
	// whose vacuous SELECT group the user guide documents.
	deselectOnly := advancedSelectorRaw(map[string]interface{}{
		"Field": "eventSource", "NotEquals": []interface{}{"s3.amazonaws.com"},
	})
	if _, err = svc.putEventSelectorsCore(store, PutEventSelectorsInput{TrailName: created.Name, AdvancedEventSelectorsRaw: deselectOnly}); err != nil {
		t.Fatalf("deselect-only selector rejected: %v", err)
	}

	// A well-formed network activity selector — eventCategory and the
	// required Equals-only eventSource — is accepted, as is an errorCode
	// filter under its Equals-only restriction.
	networkActivity := []interface{}{map[string]interface{}{
		"FieldSelectors": []interface{}{
			map[string]interface{}{"Field": "eventCategory", "Equals": []interface{}{"NetworkActivity"}},
			map[string]interface{}{"Field": "eventSource", "Equals": []interface{}{"s3.amazonaws.com"}},
			map[string]interface{}{"Field": "errorCode", "Equals": []interface{}{"VpceAccessDenied"}},
		},
	}}
	if _, err = svc.putEventSelectorsCore(store, PutEventSelectorsInput{TrailName: created.Name, AdvancedEventSelectorsRaw: networkActivity}); err != nil {
		t.Fatalf("network activity selector rejected: %v", err)
	}
}

// TestAdvancedSelectorVocabularyOnBothSurfaces pins the vocabulary gating
// on the event-data-store surface (the matrix above covers the trail
// surface): an unknown field and an operator outside a restricted field's
// documented set are InvalidEventSelectorsException on create, while a
// deselect-only selector on an unrestricted field is accepted.
func TestAdvancedSelectorVocabularyOnBothSurfaces(t *testing.T) {
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	unknownField := []interface{}{map[string]interface{}{
		"FieldSelectors": []interface{}{
			map[string]interface{}{"Field": "eventCategory", "Equals": []interface{}{"Data"}},
			map[string]interface{}{"Field": "resource.ARN", "Equals": []interface{}{"arn:aws:s3:::b/object"}},
		},
	}}
	_, err := svc.createEventDataStoreCore(store, CreateEventDataStoreInput{
		Name:                      "vocab-unknown",
		AdvancedEventSelectorsRaw: unknownField,
		AdvancedEventSelectorsSet: true,
	})
	requireAWSCodeFromMatrix(t, "eds unknown field", err, "InvalidEventSelectorsException")

	badOperator := []interface{}{map[string]interface{}{
		"FieldSelectors": []interface{}{map[string]interface{}{
			"Field": "readOnly", "NotEquals": []interface{}{"true"},
		}},
	}}
	_, err = svc.createEventDataStoreCore(store, CreateEventDataStoreInput{
		Name:                      "vocab-operator",
		AdvancedEventSelectorsRaw: badOperator,
		AdvancedEventSelectorsSet: true,
	})
	requireAWSCodeFromMatrix(t, "eds restricted operator", err, "InvalidEventSelectorsException")

	deselectOnly := []interface{}{map[string]interface{}{
		"FieldSelectors": []interface{}{map[string]interface{}{
			"Field": "eventSource", "NotEquals": []interface{}{"s3.amazonaws.com"},
		}},
	}}
	if _, err = svc.createEventDataStoreCore(store, CreateEventDataStoreInput{
		Name:                      "vocab-valid",
		AdvancedEventSelectorsRaw: deselectOnly,
		AdvancedEventSelectorsSet: true,
	}); err != nil {
		t.Fatalf("deselect-only selector rejected on create: %v", err)
	}
}

func TestPutInsightSelectorsValidationMatrix(t *testing.T) {
	ctx := context.Background()
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	created, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "ins-trail", S3BucketName: "ins-trail-bucket", Region: "us-east-1"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	valid := []interface{}{map[string]interface{}{"InsightType": "ApiCallRateInsight"}}

	_, err = svc.putInsightSelectorsCore(store, PutInsightSelectorsInput{
		TrailName:           created.Name,
		InsightSelectorsRaw: []interface{}{map[string]interface{}{}},
	})
	requireAWSCodeFromMatrix(t, "missing insight type", err, "InvalidInsightSelectorsException")

	_, err = svc.putInsightSelectorsCore(store, PutInsightSelectorsInput{
		TrailName:           created.Name,
		EventDataStore:      "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/x",
		InsightSelectorsRaw: valid,
	})
	requireAWSCodeFromMatrix(t, "trail with eds", err, "InvalidParameterCombinationException")

	_, err = svc.putInsightSelectorsCore(store, PutInsightSelectorsInput{
		EventDataStore:      "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/x",
		InsightSelectorsRaw: valid,
	})
	requireAWSCodeFromMatrix(t, "eds source without destination", err, "InvalidParameterCombinationException")

	_, err = svc.putInsightSelectorsCore(store, PutInsightSelectorsInput{
		EventDataStore:      "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/x",
		InsightsDestination: "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/y",
		InsightSelectorsRaw: valid,
	})
	requireAWSCodeFromMatrix(t, "eds form", err, "UnsupportedOperationException")

	_, err = svc.putInsightSelectorsCore(store, PutInsightSelectorsInput{
		InsightSelectorsRaw: valid,
	})
	requireAWSCodeFromMatrix(t, "no target", err, "InvalidParameterCombinationException")

	if _, err = svc.putInsightSelectorsCore(store, PutInsightSelectorsInput{TrailName: created.Name, InsightSelectorsRaw: valid}); err != nil {
		t.Fatalf("valid trail form rejected: %v", err)
	}
}

func TestQueryAndListMaxResultsBounds(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	eds, err := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore("max-eds", "acc123", "us-east-1"))
	if err != nil {
		t.Fatalf("eds create failed: %v", err)
	}
	if err := store.SaveQuery(&cloudtrailstore.QueryRecord{
		QueryID: "q-max", EventDataStore: eds.EventDataStoreID, QueryStatus: "FINISHED",
	}); err != nil {
		t.Fatalf("query seed failed: %v", err)
	}

	_, err = svc.getQueryResultsCore(store, GetQueryResultsInput{QueryID: "q-max", MaxQueryResults: 1001})
	requireAWSCodeFromMatrix(t, "get query results bound", err, "InvalidMaxResultsException")

	_, err = svc.listQueriesCore(store, ListQueriesInput{EventDataStore: eds.EventDataStoreARN, MaxResults: 1001})
	requireAWSCodeFromMatrix(t, "list queries bound", err, "InvalidMaxResultsException")

	_, err = svc.listEventDataStoresCore(store, ListEventDataStoresInput{MaxResults: 1001})
	requireAWSCodeFromMatrix(t, "list eds bound", err, "InvalidMaxResultsException")

	// ListChannels and ListImports declare no InvalidMaxResultsException;
	// their generic declared invalid-input error carries the rejection.
	_, err = svc.listChannelsCore(store, ListChannelsInput{MaxResults: 1001})
	requireAWSCodeFromMatrix(t, "list channels bound", err, "InvalidParameterException")

	_, err = svc.listImportsCore(store, ListImportsInput{MaxResults: 1001})
	requireAWSCodeFromMatrix(t, "list imports bound", err, "InvalidParameterException")

	_, err = svc.listImportsCore(store, ListImportsInput{ImportStatus: "RUNNING"})
	requireAWSCodeFromMatrix(t, "bad import status filter", err, "InvalidParameterException")

	_, err = svc.lookupEventsCore(store, LookupEventsInput{MaxResults: 51})
	requireAWSCodeFromMatrix(t, "lookup events bound", err, "InvalidMaxResultsException")

	_, err = svc.lookupEventsCore(store, LookupEventsInput{
		LookupAttributes: []interface{}{map[string]interface{}{
			"AttributeKey":   "EventName",
			"AttributeValue": string(make([]byte, 2001)),
		}},
	})
	requireAWSCodeFromMatrix(t, "lookup attribute value length", err, "InvalidLookupAttributesException")
}

func TestQueryAIValidation(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	if _, err := svc.searchSampleQueriesCore(SearchSampleQueriesInput{SearchPhrase: "x"}); err == nil {
		t.Error("short search phrase accepted")
	} else {
		requireAWSCodeFromMatrix(t, "short search phrase", err, "InvalidParameterException")
	}
	if _, err := svc.searchSampleQueriesCore(SearchSampleQueriesInput{SearchPhrase: "events", MaxResults: 51}); err == nil {
		t.Error("MaxResults 51 accepted")
	} else {
		requireAWSCodeFromMatrix(t, "search max results", err, "InvalidParameterException")
	}
	if _, err := svc.searchSampleQueriesCore(SearchSampleQueriesInput{SearchPhrase: "events"}); err != nil {
		t.Fatalf("valid search rejected: %v", err)
	}

	eds, err := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore("ai-eds", "acc123", "us-east-1"))
	if err != nil {
		t.Fatalf("eds create failed: %v", err)
	}
	_, err = svc.generateQueryCore(store, GenerateQueryInput{
		EventDataStores: []string{eds.EventDataStoreARN, eds.EventDataStoreARN}, Prompt: "show me all events",
	})
	requireAWSCodeFromMatrix(t, "two eds", err, "InvalidParameterException")

	_, err = svc.generateQueryCore(store, GenerateQueryInput{EventDataStores: []string{eds.EventDataStoreARN}, Prompt: "x"})
	requireAWSCodeFromMatrix(t, "short prompt", err, "InvalidParameterException")

	if _, err = svc.generateQueryCore(store, GenerateQueryInput{EventDataStores: []string{eds.EventDataStoreARN}, Prompt: "show me all events"}); err != nil {
		t.Fatalf("valid generate rejected: %v", err)
	}
}

func TestStartImportDestinationBound(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	_, err := svc.startImportCore(t.Context(), store, StartImportInput{
		Destinations: []string{
			"arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/a",
			"arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/b",
		},
	})
	requireAWSCodeFromMatrix(t, "two destinations", err, "InvalidParameterException")
}

func TestPutEventConfigurationValidationMatrix(t *testing.T) {
	ctx := context.Background()
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	trail, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "cfg-trail", S3BucketName: "cfg-trail-bucket", Region: "us-east-1"})
	if err != nil {
		t.Fatalf("trail create failed: %v", err)
	}

	_, err = svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName:      trail.Name,
		EventDataStore: "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/x",
	})
	requireAWSCodeFromMatrix(t, "both targets", err, "InvalidParameterCombinationException")

	_, err = svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		EventDataStore: "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/absent",
		Params:         map[string]interface{}{"MaxEventSize": "Standard"},
	})
	requireAWSCodeFromMatrix(t, "absent eds target", err, "EventDataStoreNotFoundException")

	_, err = svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName: trail.Name,
		Params: map[string]interface{}{"AggregationConfigurations": []interface{}{
			map[string]interface{}{"EventCategory": "Data"},
		}},
	})
	requireAWSCodeFromMatrix(t, "missing templates", err, "InvalidParameterException")

	_, err = svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName: trail.Name,
		Params: map[string]interface{}{"AggregationConfigurations": []interface{}{
			map[string]interface{}{"Templates": []interface{}{"NOT_A_TEMPLATE"}, "EventCategory": "Data"},
		}},
	})
	requireAWSCodeFromMatrix(t, "bad template value", err, "InvalidParameterException")

	_, err = svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName: trail.Name,
		Params: map[string]interface{}{"AggregationConfigurations": []interface{}{
			map[string]interface{}{"Templates": []interface{}{"API_ACTIVITY"}, "EventCategory": "management"},
		}},
	})
	requireAWSCodeFromMatrix(t, "bad event category", err, "InvalidParameterException")

	_, err = svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName: trail.Name,
		Params: map[string]interface{}{"ContextKeySelectors": []interface{}{
			map[string]interface{}{"Type": "WrongContext"},
		}},
	})
	requireAWSCodeFromMatrix(t, "bad context key type", err, "InvalidParameterException")

	_, err = svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName: trail.Name,
		Params: map[string]interface{}{"ContextKeySelectors": []interface{}{
			map[string]interface{}{"Type": "TagContext"},
		}},
	})
	requireAWSCodeFromMatrix(t, "missing equals", err, "InvalidParameterException")

	// The valid scalar form persists.
	if _, err = svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName: trail.Name,
		Params: map[string]interface{}{"AggregationConfigurations": []interface{}{
			map[string]interface{}{"Templates": []interface{}{"API_ACTIVITY"}, "EventCategory": "Data"},
		}},
	}); err != nil {
		t.Fatalf("valid aggregation rejected: %v", err)
	}
}

func TestDeleteTrailWhileLogging(t *testing.T) {
	ctx := context.Background()
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	created, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "del-logging", S3BucketName: "del-logging-bucket", Region: "us-east-1"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := store.StartLogging(created.Name); err != nil {
		t.Fatalf("start logging failed: %v", err)
	}
	// DeleteTrail documents no logging precondition: a logging trail
	// deletes like any other.
	if err := svc.deleteTrailCore(store, DeleteTrailInput{NameOrARN: created.Name}); err != nil {
		t.Fatalf("delete on a logging trail rejected: %v", err)
	}
}

func TestGetInsightSelectorsNotEnabled(t *testing.T) {
	ctx := context.Background()
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	created, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "ins-none", S3BucketName: "ins-none-bucket", Region: "us-east-1"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	_, err = svc.getInsightSelectorsCore(store, created.Name)
	requireAWSCodeFromMatrix(t, "not enabled", err, "InsightNotEnabledException")

	// After selectors are configured the read succeeds.
	if _, err := svc.putInsightSelectorsCore(store, PutInsightSelectorsInput{
		TrailName:           created.Name,
		InsightSelectorsRaw: []interface{}{map[string]interface{}{"InsightType": "ApiCallRateInsight"}},
	}); err != nil {
		t.Fatalf("put insight selectors failed: %v", err)
	}
	resp, err := svc.getInsightSelectorsCore(store, created.Name)
	if err != nil {
		t.Fatalf("get after put failed: %v", err)
	}
	if _, ok := resp["InsightSelectors"]; !ok {
		t.Fatal("get after put returned no InsightSelectors")
	}
}

func TestStartImportStrictTimes(t *testing.T) {
	// The source-bucket verification needs an S3 invoker; the permissive
	// fake answers every bucket as existing, keeping the pin's subject the
	// time parsing.
	svc := &CloudTrailService{}
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: newPermissiveS3Invoker()})
	store := newErrorTestServiceStore(t)

	// A present time value matching neither wire form is rejected instead
	// of silently dropping the bound.
	_, err := svc.startImportCore(t.Context(), store, StartImportInput{
		Destinations:         []string{"arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/a"},
		ImportSourceRaw:      map[string]interface{}{"S3": map[string]interface{}{"S3LocationUri": "s3://b/p", "S3BucketRegion": "us-east-1", "S3BucketAccessRoleArn": "arn:aws:iam::acc123:role/import-role"}},
		ImportSourceProvided: true,
		StartEventTimeStr:    "not-a-timestamp",
	})
	requireAWSCodeFromMatrix(t, "unparseable string", err, "InvalidParameterException")

	_, err = svc.startImportCore(t.Context(), store, StartImportInput{
		Destinations:         []string{"arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/a"},
		ImportSourceRaw:      map[string]interface{}{"S3": map[string]interface{}{"S3LocationUri": "s3://b/p", "S3BucketRegion": "us-east-1", "S3BucketAccessRoleArn": "arn:aws:iam::acc123:role/import-role"}},
		ImportSourceProvided: true,
		StartEventTimeRaw:    "garbage",
	})
	requireAWSCodeFromMatrix(t, "unparseable raw", err, "InvalidParameterException")

	// Both valid forms parse, and the epoch bound bounds the import.
	eds, edsErr := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore("import-times", "acc123", "us-east-1"))
	if edsErr != nil {
		t.Fatalf("eds create failed: %v", edsErr)
	}
	start := float64(1700000000)
	created, err := svc.startImportCore(t.Context(), store, StartImportInput{
		Destinations:         []string{eds.EventDataStoreARN},
		ImportSourceRaw:      map[string]interface{}{"S3": map[string]interface{}{"S3LocationUri": "s3://b/p", "S3BucketRegion": "us-east-1", "S3BucketAccessRoleArn": "arn:aws:iam::acc123:role/import-role"}},
		ImportSourceProvided: true,
		StartEventTimeRaw:    start,
		EndEventTimeStr:      "2023-11-14T22:13:21Z",
	})
	if err != nil {
		t.Fatalf("valid forms rejected: %v", err)
	}
	if created["StartEventTime"] == nil || created["EndEventTime"] == nil {
		t.Fatal("time bounds missing from the created import")
	}

	// EndEventTime before StartEventTime is an invalid range.
	_, err = svc.startImportCore(t.Context(), store, StartImportInput{
		Destinations:         []string{eds.EventDataStoreARN},
		ImportSourceRaw:      map[string]interface{}{"S3": map[string]interface{}{"S3LocationUri": "s3://b/p", "S3BucketRegion": "us-east-1", "S3BucketAccessRoleArn": "arn:aws:iam::acc123:role/import-role"}},
		ImportSourceProvided: true,
		StartEventTimeRaw:    start,
		EndEventTimeRaw:      start - 100,
	})
	requireAWSCodeFromMatrix(t, "end before start", err, "InvalidParameterException")
}

// TestPublicKeyPKCS1AndFingerprint pins the public-key wire contract: the
// served DER parses as PKCS#1 (the documented Value form) and not as the
// SPKI wrapping, and the fingerprint is a recomputable function of the
// served bytes — 32 lowercase hex characters, the length of the sample
// fingerprint in the digest file reference.
func TestPublicKeyPKCS1AndFingerprint(t *testing.T) {
	pk, _, err := cloudtrailstore.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	if _, err := x509.ParsePKCS1PublicKey(pk.Value); err != nil {
		t.Fatalf("served DER is not parseable as PKCS#1: %v", err)
	}
	if _, err := x509.ParsePKIXPublicKey(pk.Value); err == nil {
		t.Fatal("served DER parsed as SPKI — the PKIX encoding must never be served")
	}
	sum := md5.Sum(pk.Value)
	if got, want := pk.Fingerprint(), hex.EncodeToString(sum[:]); got != want {
		t.Fatalf("Fingerprint = %q, want the MD5 hex of the served DER %q", got, want)
	}
	if len(pk.Fingerprint()) != 32 {
		t.Fatalf("Fingerprint length = %d, want 32 hex characters", len(pk.Fingerprint()))
	}
}

// TestUpdateTrailEnablesValidationWithKey pins the update-path key
// generation: enabling log file validation provisions exactly one key (the
// transition CreateTrail already covers at creation), and re-asserting the
// flag on a validated trail adds nothing.
func TestUpdateTrailEnablesValidationWithKey(t *testing.T) {
	ctx := context.Background()
	store := newErrorTestServiceStore(t)
	svc := newTrailTestService(store)

	created, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "matrix-pk-trail", S3BucketName: "matrix-pk-bucket", Region: "us-east-1"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	before, err := store.ListPublicKeys(nil, nil)
	if err != nil {
		t.Fatalf("baseline list failed: %v", err)
	}

	updated, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{Name: created.Name, Params: map[string]interface{}{"EnableLogFileValidation": true}})
	if err != nil {
		t.Fatalf("enabling validation failed: %v", err)
	}
	if !updated.LogFileValidationEnabled {
		t.Fatal("update did not persist LogFileValidationEnabled")
	}

	after, err := store.ListPublicKeys(nil, nil)
	if err != nil {
		t.Fatalf("post-enable list failed: %v", err)
	}
	if len(after) != len(before)+1 {
		t.Fatalf("enabling validation generated %d keys, want exactly 1", len(after)-len(before))
	}

	if _, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{Name: created.Name, Params: map[string]interface{}{"EnableLogFileValidation": true}}); err != nil {
		t.Fatalf("re-asserting validation failed: %v", err)
	}
	final, err := store.ListPublicKeys(nil, nil)
	if err != nil {
		t.Fatalf("post-reassert list failed: %v", err)
	}
	if len(final) != len(after) {
		t.Fatalf("re-asserting the flag generated another key (%d → %d)", len(after), len(final))
	}
}

// TestListPublicKeysDefaultAndTimeRange pins the listing semantics: with no
// bounds both default to the current time (only currently-valid keys are
// returned), an explicit past window excludes a key generated now, and a
// start after the end is the declared InvalidTimeRangeException.
func TestListPublicKeysDefaultAndTimeRange(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	if _, err := store.GenerateAndStorePublicKey("matrix-region-key"); err != nil {
		t.Fatalf("key generation failed: %v", err)
	}

	res, err := svc.listPublicKeysCore(store, ListPublicKeysInput{})
	if err != nil {
		t.Fatalf("default listing failed: %v", err)
	}
	list, ok := res["PublicKeyList"].([]map[string]interface{})
	if !ok || len(list) == 0 {
		t.Fatalf("default listing returned no keys: %v", res["PublicKeyList"])
	}
	first := list[0]
	fp, _ := first["Fingerprint"].(string)
	if len(fp) != 32 {
		t.Fatalf("default listing fingerprint = %q, want 32 hex characters", fp)
	}

	pastStart := time.Now().UTC().Add(-2 * time.Hour)
	pastEnd := pastStart.Add(time.Minute)
	res, err = svc.listPublicKeysCore(store, ListPublicKeysInput{
		StartTimeStr: pastStart.Format(time.RFC3339),
		EndTimeStr:   pastEnd.Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("past-window listing failed: %v", err)
	}
	list, _ = res["PublicKeyList"].([]map[string]interface{})
	if len(list) != 0 {
		t.Fatalf("past window returned %d keys, want 0 (the key is valid now)", len(list))
	}

	future := time.Now().UTC().Add(time.Hour)
	_, err = svc.listPublicKeysCore(store, ListPublicKeysInput{
		StartTimeStr: future.Format(time.RFC3339),
		EndTimeStr:   pastStart.Format(time.RFC3339),
	})
	requireAWSCodeFromMatrix(t, "start after end", err, "InvalidTimeRangeException")
}

// TestCreateQuotaGuards pins the three resource quotas the fetched
// "Quotas in AWS CloudTrail" table documents: five trails per Region
// (MaximumNumberOfTrailsExceededException on the sixth), ten event data
// stores counting every lifecycle stage (EventDataStoreMaxLimitExceededException
// on the eleventh), and twenty-five channels
// (ChannelMaxLimitExceededException on the twenty-sixth).
func TestCreateQuotaGuards(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-quota-test"
	defer os.RemoveAll(tmpDir)
	st, err := storage.Open(tmpDir)
	if err != nil {
		t.Fatalf("storage.Open: %v", err)
	}
	defer st.Close()
	store := cloudtrailstore.NewCloudTrailStore(st, "acc123", "us-east-1")
	svc := newTrailTestService(store)
	// refEDS carries the first quota-occupying event data store's record
	// for the channel subtest's destination.
	var refEDS *cloudtrailstore.EventDataStore

	t.Run("sixth trail is rejected", func(t *testing.T) {
		for i := 0; i < cloudtrailstore.MaxTrailsPerRegion; i++ {
			_, err := svc.createTrailCore(context.Background(), store, CreateTrailInput{
				Name:         fmt.Sprintf("quota-trail-%d", i),
				S3BucketName: "quota-bucket",
			})
			if err != nil {
				t.Fatalf("trail %d creation failed: %v", i, err)
			}
		}
		_, err := svc.createTrailCore(context.Background(), store, CreateTrailInput{
			Name:         "quota-trail-overflow",
			S3BucketName: "quota-bucket",
		})
		requireAWSCodeFromMatrix(t, "sixth trail", err, "MaximumNumberOfTrailsExceededException")
	})

	t.Run("eleventh event data store is rejected", func(t *testing.T) {
		for i := 0; i < cloudtrailstore.MaxEventDataStoresPerRegion; i++ {
			eds := cloudtrailstore.NewEventDataStore(
				fmt.Sprintf("quota-eds-%d", i), "acc123", "us-east-1")
			if _, err := store.CreateEventDataStore(eds); err != nil {
				t.Fatalf("event data store %d creation failed: %v", i, err)
			}
			if i == 0 {
				// The channel subtest's destination: one of this run's
				// quota-occupying records, so no extra store is created.
				refEDS = eds
			}
		}
		_, err := svc.createEventDataStoreCore(store, CreateEventDataStoreInput{Name: "quota-eds-overflow"})
		requireAWSCodeFromMatrix(t, "eleventh event data store", err, "EventDataStoreMaxLimitExceededException")
	})

	t.Run("twenty-sixth channel is rejected", func(t *testing.T) {
		destinations := []interface{}{map[string]interface{}{
			"Type":     "EVENT_DATA_STORE",
			"Location": refEDS.EventDataStoreARN,
		}}
		for i := 0; i < cloudtrailstore.MaxChannelsPerRegion; i++ {
			_, err := svc.createChannelCore(store, CreateChannelInput{
				Name:            fmt.Sprintf("quota-channel-%d", i),
				Source:          fmt.Sprintf("custom/source-%d", i),
				DestinationsSet: true,
				DestinationsRaw: destinations,
			})
			if err != nil {
				t.Fatalf("channel %d creation failed: %v", i, err)
			}
		}
		_, err := svc.createChannelCore(store, CreateChannelInput{
			Name:            "quota-channel-overflow",
			Source:          "custom/source-overflow",
			DestinationsSet: true,
			DestinationsRaw: destinations,
		})
		requireAWSCodeFromMatrix(t, "twenty-sixth channel", err, "ChannelMaxLimitExceededException")
	})
}

// TestIngestionLifecycleTransitions pins the documented Status lifecycle:
// "To stop ingestion, the event data store Status must be ENABLED" and
// "To start ingestion, the event data store Status must be
// STOPPED_INGESTION" (StartEventDataStoreIngestion /
// StopEventDataStoreIngestion) — the stop lands STOPPED_INGESTION and the
// start returns ENABLED; a repeat stop is refused; a category outside the
// operations' set (a channel's ActivityAuditLog store) is refused; and
// federation on a PENDING_DELETION store is the declared
// InactiveEventDataStoreException in both directions.
func TestIngestionLifecycleTransitions(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	created, err := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore("ingestion-lifecycle", store.GetAccountID(), store.GetRegion()))
	if err != nil {
		t.Fatalf("create eds: %v", err)
	}

	if _, err := svc.stopEventDataStoreIngestionCore(store, EventDataStoreIDInput{EventDataStore: created.EventDataStoreARN}); err != nil {
		t.Fatalf("stop from ENABLED rejected: %v", err)
	}
	after, err := store.GetEventDataStore(created.EventDataStoreID)
	if err != nil {
		t.Fatalf("get eds: %v", err)
	}
	if after.Status != "STOPPED_INGESTION" || after.IngestionEnabled {
		t.Fatalf("after stop: status = %q, ingestionEnabled = %v", after.Status, after.IngestionEnabled)
	}

	_, err = svc.stopEventDataStoreIngestionCore(store, EventDataStoreIDInput{EventDataStore: created.EventDataStoreARN})
	requireAWSCode(t, err, "InvalidEventDataStoreStatusException", 400)

	if _, err := svc.startEventDataStoreIngestionCore(store, EventDataStoreIDInput{EventDataStore: created.EventDataStoreARN}); err != nil {
		t.Fatalf("start from STOPPED_INGESTION rejected: %v", err)
	}
	after, err = store.GetEventDataStore(created.EventDataStoreID)
	if err != nil {
		t.Fatalf("get eds: %v", err)
	}
	if after.Status != "ENABLED" || !after.IngestionEnabled {
		t.Fatalf("after start: status = %q, ingestionEnabled = %v", after.Status, after.IngestionEnabled)
	}

	_, err = svc.startEventDataStoreIngestionCore(store, EventDataStoreIDInput{EventDataStore: created.EventDataStoreARN})
	requireAWSCode(t, err, "InvalidEventDataStoreStatusException", 400)

	audit := cloudtrailstore.NewEventDataStore("audit-category", store.GetAccountID(), store.GetRegion())
	audit.AdvancedEventSelectors = []cloudtrailstore.AdvancedEventSelector{{
		Name: "channel destination",
		FieldSelectors: []cloudtrailstore.AdvancedFieldSelector{
			{Field: "eventCategory", Equals: []string{"ActivityAuditLog"}},
		},
	}}
	auditEDS, err := store.CreateEventDataStore(audit)
	if err != nil {
		t.Fatalf("create audit eds: %v", err)
	}
	_, err = svc.stopEventDataStoreIngestionCore(store, EventDataStoreIDInput{EventDataStore: auditEDS.EventDataStoreARN})
	requireAWSCode(t, err, "InvalidEventDataStoreCategoryException", 400)

	pending := cloudtrailstore.NewEventDataStore("federation-inactive", store.GetAccountID(), store.GetRegion())
	pendingEDS, err := store.CreateEventDataStore(pending)
	if err != nil {
		t.Fatalf("create pending eds: %v", err)
	}
	_, err = store.MutateEventDataStore(pendingEDS.EventDataStoreID, func(e *cloudtrailstore.EventDataStore) error {
		e.Status = "PENDING_DELETION"
		return nil
	})
	if err != nil {
		t.Fatalf("mark pending: %v", err)
	}
	_, err = svc.enableFederationCore(context.Background(), store, EnableFederationInput{
		EventDataStore:    pendingEDS.EventDataStoreARN,
		FederationRoleArn: "arn:aws:iam::acc123:role/federation",
	})
	requireAWSCode(t, err, "InactiveEventDataStoreException", 400)
	_, err = svc.disableFederationCore(store, DisableFederationInput{EventDataStore: pendingEDS.EventDataStoreARN})
	requireAWSCode(t, err, "InactiveEventDataStoreException", 400)
}
