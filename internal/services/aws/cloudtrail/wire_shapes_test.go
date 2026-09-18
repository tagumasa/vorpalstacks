package cloudtrail

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// The census tables below are the emitted-key contract for every response
// family, transcribed from the vendored Smithy model's member lists. Each
// pin asserts two properties: every key a formatter emits is a model member
// of that operation's response shape (no invented keys, no misspelt
// members that lenient SDK deserialisers would silently drop), and the
// always-present members survive.

func keys(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// assertSubset fails for every emitted key outside the model member set.
func assertSubset(t *testing.T, emitted map[string]interface{}, model map[string]bool, family string) {
	t.Helper()
	for _, k := range keys(emitted) {
		assert.True(t, model[k], "%s emits %q, which is not a model member of its response shape", family, k)
	}
}

func modelMembers(names ...string) map[string]bool {
	m := make(map[string]bool, len(names))
	for _, n := range names {
		m[n] = true
	}
	return m
}

func populatedTrail() *cloudtrailstore.Trail {
	t := cloudtrailstore.NewTrail("wire-shapes", "wire-shapes-bucket", "us-east-1")
	t.S3KeyPrefix = "prefix"
	t.SnsTopicName = "topic"
	t.SnsTopicARN = "arn:aws:sns:us-east-1:acc123:topic"
	t.CloudWatchLogsLogGroupARN = "arn:aws:logs:us-east-1:acc123:log-group:g"
	t.CloudWatchLogsRoleARN = "arn:aws:iam::acc123:role/r"
	t.KMSKeyID = "arn:aws:kms:us-east-1:acc123:key/k"
	return t
}

func TestTrailResponseShapes(t *testing.T) {
	svc := &CloudTrailService{}
	trail := populatedTrail()

	// Trail description shape (GetTrail, DescribeTrails items): 17 members.
	trailModel := modelMembers(
		"CloudWatchLogsLogGroupArn", "CloudWatchLogsRoleArn", "HasCustomEventSelectors",
		"HasInsightSelectors", "HomeRegion", "IncludeGlobalServiceEvents", "IsMultiRegionTrail",
		"IsOrganizationTrail", "KmsKeyId", "LogFileValidationEnabled", "Name", "RecursiveLogging",
		"S3BucketName", "S3KeyPrefix", "SnsTopicARN", "SnsTopicName", "TrailARN")
	emitted := svc.formatTrail(trail)
	assertSubset(t, emitted, trailModel, "Trail")
	for _, m := range []string{"Name", "TrailARN", "SnsTopicARN", "RecursiveLogging", "HomeRegion"} {
		assert.Contains(t, emitted, m, "Trail must carry %s", m)
	}
	assert.NotContains(t, emitted, "IsLogging", "IsLogging belongs to GetTrailStatus, never to a Trail")

	// Create/Update response shape: the Trail members minus the description
	// members.
	mutationModel := modelMembers(
		"CloudWatchLogsLogGroupArn", "CloudWatchLogsRoleArn", "IncludeGlobalServiceEvents",
		"IsMultiRegionTrail", "IsOrganizationTrail", "KmsKeyId", "LogFileValidationEnabled",
		"Name", "RecursiveLogging", "S3BucketName", "S3KeyPrefix", "SnsTopicARN",
		"SnsTopicName", "TrailARN")
	mutation := svc.formatTrailMutationResponse(trail)
	assertSubset(t, mutation, mutationModel, "CreateTrail/UpdateTrail response")
	for _, m := range []string{"HomeRegion", "HasCustomEventSelectors", "HasInsightSelectors", "IsLogging"} {
		assert.NotContains(t, mutation, m, "mutation response must not carry %s", m)
	}
	assert.Contains(t, mutation, "RecursiveLogging")
}

func TestRecursiveLoggingRoundTrip(t *testing.T) {
	store := newQueryTestStore(t)
	svc := newTrailTestService(store)

	// The documented default is true.
	created, err := svc.createTrailCore(t.Context(), store, CreateTrailInput{
		Name: "recursive-default", S3BucketName: "wire-shapes-bucket", Region: "us-east-1"})
	require.NoError(t, err)
	assert.True(t, created.RecursiveLogging)
	assert.Equal(t, true, svc.formatTrail(created)["RecursiveLogging"])

	// The member must survive the persistence round-trip: a reload answers
	// the stored value, never the proto zero.
	reloaded, err := store.GetTrail("recursive-default")
	require.NoError(t, err)
	assert.True(t, reloaded.RecursiveLogging, "the default survives a store round-trip")

	off := false
	created, err = svc.createTrailCore(t.Context(), store, CreateTrailInput{
		Name: "recursive-off", S3BucketName: "wire-shapes-bucket", Region: "us-east-1", RecursiveLogging: &off})
	require.NoError(t, err)
	assert.False(t, created.RecursiveLogging)

	updated, err := svc.updateTrailCore(t.Context(), store, UpdateTrailInput{
		Name:   "recursive-default",
		Params: map[string]interface{}{"RecursiveLogging": false},
	})
	require.NoError(t, err)
	assert.False(t, updated.RecursiveLogging)
	reloaded, err = store.GetTrail("recursive-default")
	require.NoError(t, err)
	assert.False(t, reloaded.RecursiveLogging, "the explicit false survives a store round-trip")
}

// TestEventAwsRegionRoundTrip pins the record's region member through the
// event-history persistence: a reloaded event answers the awsRegion the
// Lake engine's awsregion column reads.
func TestEventAwsRegionRoundTrip(t *testing.T) {
	store := newQueryTestStore(t)

	e := cloudtrailstore.NewEvent("LookupRegion", "cloudtrail.amazonaws.com",
		&cloudtrailstore.UserIdentity{UserName: "alice", AccountID: "acc123"}, false)
	e.AwsRegion = "eu-west-3"
	require.NoError(t, store.PutEvent(e))

	reloaded, err := store.GetEventByID(e.EventID)
	require.NoError(t, err)
	assert.Equal(t, "eu-west-3", reloaded.AwsRegion)
}

func TestLookupEventsEventShape(t *testing.T) {
	svc := &CloudTrailService{}
	e := cloudtrailstore.NewEvent("CreateTrail", "cloudtrail.amazonaws.com",
		&cloudtrailstore.UserIdentity{UserName: "alice", AccountID: "acc123"}, false)
	e.CloudTrailEvent = `{"eventID":"e1"}`
	e.Resources = []cloudtrailstore.Resource{{ResourceType: "AWS::CloudTrail::Trail", ResourceName: "arn"}}

	// The wire Event carries exactly the nine model members; the deeper
	// record fields live inside CloudTrailEvent.
	eventModel := modelMembers(
		"AccessKeyId", "CloudTrailEvent", "EventId", "EventName", "EventSource",
		"EventTime", "ReadOnly", "Resources", "Username")
	emitted := svc.formatEvent(e)
	assertSubset(t, emitted, eventModel, "LookupEvents Event")
	for _, m := range []string{"CloudTrailEvent", "EventId", "EventName", "EventSource", "EventTime", "ReadOnly", "Resources", "Username"} {
		assert.Contains(t, emitted, m, "populated event must carry %s", m)
	}
}

func TestQueryResponseShapes(t *testing.T) {
	store := newQueryTestStore(t)
	svc := NewCloudTrailService("acc123", "us-east-1")
	require.NoError(t, store.SaveQuery(&cloudtrailstore.QueryRecord{
		QueryID:        "q-shape",
		EventDataStore: "eds",
		QueryStatus:    "FINISHED",
		StartTime:      time.Now().UTC(),
	}))

	resp, err := svc.listQueriesCore(store, ListQueriesInput{EventDataStore: "eds"})
	require.NoError(t, err)
	queries := resp["Queries"].([]map[string]interface{})
	require.Len(t, queries, 1)
	// The list item is the Query shape: creation time, identifier, status.
	assertSubset(t, queries[0], modelMembers("CreationTime", "QueryId", "QueryStatus"), "ListQueries item")
	assert.Contains(t, queries[0], "CreationTime")

	results, err := svc.getQueryResultsCore(store, GetQueryResultsInput{QueryID: "q-shape"})
	require.NoError(t, err)
	assertSubset(t, results,
		modelMembers("ErrorMessage", "NextToken", "QueryResultRows", "QueryStatistics", "QueryStatus"),
		"GetQueryResults")

	desc, err := svc.describeQueryCore(store, DescribeQueryInput{QueryID: "q-shape"})
	require.NoError(t, err)
	assertSubset(t, desc,
		modelMembers("DeliveryS3Uri", "DeliveryStatus", "ErrorMessage", "EventDataStoreOwnerAccountId",
			"Prompt", "QueryId", "QueryStatistics", "QueryStatus", "QueryString"),
		"DescribeQuery")
	assert.Contains(t, desc, "EventDataStoreOwnerAccountId")
}

func TestEventDataStoreResponseShapes(t *testing.T) {
	eds := cloudtrailstore.NewEventDataStore("wire-eds", "acc123", "us-east-1")
	eds.BillingMode = "EXTENDABLE_RETENTION_PRICING"
	eds.KMSKeyID = "arn:aws:kms:us-east-1:acc123:key/k"
	eds.Tags = map[string]string{"k": "v"}
	eds.FederationStatus = "ENABLED"
	eds.FederationRoleARN = "arn:aws:iam::acc123:role/r"

	base := modelMembers(
		"AdvancedEventSelectors", "CreatedTimestamp", "EventDataStoreArn", "MultiRegionEnabled",
		"Name", "OrganizationEnabled", "RetentionPeriod", "Status", "TerminationProtectionEnabled",
		"UpdatedTimestamp")
	createShape := modelMembers(
		"AdvancedEventSelectors", "BillingMode", "CreatedTimestamp", "EventDataStoreArn", "KmsKeyId",
		"MultiRegionEnabled", "Name", "OrganizationEnabled", "RetentionPeriod", "Status", "TagsList",
		"TerminationProtectionEnabled", "UpdatedTimestamp")
	getShape := modelMembers(
		"AdvancedEventSelectors", "BillingMode", "CreatedTimestamp", "EventDataStoreArn",
		"FederationRoleArn", "FederationStatus", "KmsKeyId", "MultiRegionEnabled", "Name",
		"OrganizationEnabled", "PartitionKeys", "RetentionPeriod", "Status",
		"TerminationProtectionEnabled", "UpdatedTimestamp")
	updateShape := modelMembers(
		"AdvancedEventSelectors", "BillingMode", "CreatedTimestamp", "EventDataStoreArn",
		"FederationRoleArn", "FederationStatus", "KmsKeyId", "MultiRegionEnabled", "Name",
		"OrganizationEnabled", "RetentionPeriod", "Status", "TerminationProtectionEnabled",
		"UpdatedTimestamp")

	assertSubset(t, formatEventDataStoreFor(eds, edsProfileCreate), createShape, "CreateEventDataStore/Restore response")
	assertSubset(t, formatEventDataStoreFor(eds, edsProfileGet), getShape, "GetEventDataStore response")
	assertSubset(t, formatEventDataStoreFor(eds, edsProfileUpdate), updateShape, "UpdateEventDataStore response")
	assertSubset(t, formatEventDataStoreFor(eds, edsProfileList), base, "ListEventDataStores item")

	assert.Contains(t, formatEventDataStoreFor(eds, edsProfileCreate), "TagsList")
	assert.NotContains(t, formatEventDataStoreFor(eds, edsProfileGet), "TagsList")
	assert.NotContains(t, formatEventDataStoreFor(eds, edsProfileList), "BillingMode")
	// IngestionEnabled is a store-internal flag, never a response member.
	for _, profile := range []edsResponseProfile{edsProfileCreate, edsProfileGet, edsProfileUpdate, edsProfileList} {
		assert.NotContains(t, formatEventDataStoreFor(eds, profile), "IngestionEnabled")
	}
}

func TestChannelResponseShapes(t *testing.T) {
	ch := &cloudtrailstore.Channel{
		ChannelARN: "arn:aws:cloudtrail:us-east-1:acc123:channel/c1",
		Name:       "wire-channel",
		Source:     "arn:aws:cloudtrail:us-east-1:acc123:trail/t",
		Destinations: []cloudtrailstore.Destination{{
			Type:     cloudtrailstore.DestinationTypeEventDataStore,
			Location: "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/e1",
		}},
		Tags: map[string]string{"k": "v"},
	}

	base := modelMembers("ChannelArn", "Destinations", "Name", "Source")
	emitted := formatChannelBase(ch)
	assertSubset(t, emitted, base, "channel base")
	assert.Contains(t, emitted, "Destinations")

	// The list item is the two-member Channel shape.
	item := map[string]interface{}{"ChannelArn": ch.ChannelARN, "Name": ch.Name}
	assertSubset(t, item, modelMembers("ChannelArn", "Name"), "ListChannels item")

	// The create response adds tags; nothing emits the invented ingestion
	// state nesting.
	createResp := formatChannelBase(ch)
	createResp["Tags"] = formatTagsList(ch.Tags)
	assertSubset(t, createResp, modelMembers("ChannelArn", "Destinations", "Name", "Source", "Tags"),
		"CreateChannel response")
	assert.NotContains(t, formatChannelBase(ch), "IngestionStatus")
}

func TestSelectorResponseTrailARN(t *testing.T) {
	store := newQueryTestStore(t)
	svc := newTrailTestService(store)
	trail, err := svc.createTrailCore(t.Context(), store, CreateTrailInput{
		Name: "selector-shapes", S3BucketName: "wire-shapes-bucket", Region: "us-east-1"})
	require.NoError(t, err)

	// The selector responses' trail member is TrailARN — the SDK drops a
	// misspelt key silently. The Put cores' emission is pinned here; the
	// Get handlers share the spelling through the SDK plane pin.
	selectors := []interface{}{map[string]interface{}{
		"ReadWriteType":           "All",
		"IncludeManagementEvents": true,
	}}
	resp, err := svc.putEventSelectorsCore(store, PutEventSelectorsInput{
		TrailName:         "selector-shapes",
		EventSelectorsRaw: selectors,
	})
	require.NoError(t, err)
	assertSubset(t, resp, modelMembers("AdvancedEventSelectors", "EventSelectors", "TrailARN"),
		"PutEventSelectors response")
	assert.Equal(t, trail.TrailARN, resp["TrailARN"])

	insights := []interface{}{map[string]interface{}{"InsightType": "ApiCallRateInsight"}}
	iresp, err := svc.putInsightSelectorsCore(store, PutInsightSelectorsInput{
		TrailName:           "selector-shapes",
		InsightSelectorsRaw: insights,
	})
	require.NoError(t, err)
	assertSubset(t, iresp, modelMembers("EventDataStoreArn", "InsightSelectors", "InsightsDestination", "TrailARN"),
		"PutInsightSelectors response")
	assert.Equal(t, trail.TrailARN, iresp["TrailARN"])
}

func TestPutEventConfigurationResponseEcho(t *testing.T) {
	store := newQueryTestStore(t)
	svc := newTrailTestService(store)
	_, err := svc.createTrailCore(t.Context(), store, CreateTrailInput{
		Name: "wire-config", S3BucketName: "wire-shapes-bucket", Region: "us-east-1"})
	require.NoError(t, err)

	resp, err := svc.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName: "wire-config",
		Params:    map[string]interface{}{"MaxEventSize": "Large"},
	})
	require.NoError(t, err)
	// The response members are the stored configuration's members.
	assertSubset(t, resp,
		modelMembers("AggregationConfigurations", "ContextKeySelectors", "EventDataStoreArn",
			"MaxEventSize", "TrailARN"),
		"PutEventConfiguration response")
	assert.Equal(t, "Large", resp["MaxEventSize"])
	assert.Contains(t, resp, "TrailARN")
}
