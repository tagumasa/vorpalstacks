package cloudwatchlogs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The emitted-key vocabulary pin: every operation's response carries only
// the model's output members (emitted is a subset of the model's set) and
// every member the model marks required is present. The expected sets are
// transcribed from the vendored Smithy model (models/cloudwatch-logs
// /2014-03-28); an operation whose output targets Unit has no members.
func TestEmittedKeyVocabulary(t *testing.T) {
	svc, reqCtx := newVocabularyEnv(t)
	ctx := context.Background()
	now := time.Now().UnixMilli()

	// The pointer addresses a seeded JSON event with this test's own
	// clock — the record-resolution validates the pointer against the
	// store, so the addressed record must exist at exactly this
	// timestamp.
	vocabSeedPointerEvent(t, svc, reqCtx, ctx, now)
	ptr := func() string {
		return eventPointer("vocab-group", "s1", now-2000, `{"level": "INFO", "service": "vocab"}`)
	}

	cases := []struct {
		op       string
		invoke   func() (interface{}, error)
		model    []string
		required []string
	}{
		{"DescribeLogGroups", func() (interface{}, error) {
			return svc.DescribeLogGroups(ctx, reqCtx, vocabRequest(nil))
		}, []string{"logGroups", "nextToken"}, nil},
		{"ListLogGroups", func() (interface{}, error) {
			return svc.ListLogGroups(ctx, reqCtx, vocabRequest(nil))
		}, []string{"logGroups", "nextToken"}, nil},
		{"DescribeLogStreams", func() (interface{}, error) {
			return svc.DescribeLogStreams(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group"}))
		}, []string{"logStreams", "nextToken"}, nil},
		{"GetLogEvents", func() (interface{}, error) {
			return svc.GetLogEvents(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group", "logStreamName": "s1"}))
		}, []string{"events", "nextForwardToken", "nextBackwardToken"}, nil},
		{"FilterLogEvents", func() (interface{}, error) {
			return svc.FilterLogEvents(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group"}))
		}, []string{"events", "searchedLogStreams", "nextToken"}, nil},
		{"PutLogEvents", func() (interface{}, error) {
			return svc.PutLogEvents(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupName": "vocab-group", "logStreamName": "s1",
				"logEvents": []interface{}{map[string]interface{}{"timestamp": float64(now), "message": "vocab"}},
			}))
		}, []string{"nextSequenceToken", "rejectedLogEventsInfo", "rejectedEntityInfo"}, nil},
		{"DescribeMetricFilters", func() (interface{}, error) {
			return svc.DescribeMetricFilters(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group"}))
		}, []string{"metricFilters", "nextToken"}, nil},
		{"TestMetricFilter", func() (interface{}, error) {
			return svc.TestMetricFilter(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupIdentifier": "vocab-group",
				"filterPattern":      "{ level = INFO }",
				"logEventMessages":   []interface{}{`{"level":"INFO"}`},
			}))
		}, []string{"matches"}, nil},
		{"DescribeSubscriptionFilters", func() (interface{}, error) {
			return svc.DescribeSubscriptionFilters(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group"}))
		}, []string{"subscriptionFilters", "nextToken"}, nil},
		{"PutDestination", func() (interface{}, error) {
			return svc.PutDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"destinationName": "vocab-dest-2",
				"roleArn":         "arn:aws:iam::000000000000:role/vocab",
				"targetArn":       "arn:aws:kinesis:us-east-1:000000000000:stream/vocab-target",
			}))
		}, []string{"destination"}, nil},
		{"DescribeDestinations", func() (interface{}, error) {
			return svc.DescribeDestinations(ctx, reqCtx, vocabRequest(nil))
		}, []string{"destinations", "nextToken"}, nil},
		{"PutResourcePolicy", func() (interface{}, error) {
			return svc.PutResourcePolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"policyName":     "vocab-rp-2",
				"policyDocument": `{"Version":"2012-10-17"}`,
			}))
		}, []string{"resourcePolicy", "revisionId"}, nil},
		{"DescribeResourcePolicies", func() (interface{}, error) {
			return svc.DescribeResourcePolicies(ctx, reqCtx, vocabRequest(nil))
		}, []string{"resourcePolicies", "nextToken"}, nil},
		{"DescribeAccountPolicies", func() (interface{}, error) {
			return svc.DescribeAccountPolicies(ctx, reqCtx, vocabRequest(map[string]interface{}{"policyType": "DATA_PROTECTION_POLICY"}))
		}, []string{"accountPolicies", "nextToken"}, nil},
		{"GetDataProtectionPolicy", func() (interface{}, error) {
			return svc.GetDataProtectionPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupIdentifier": "vocab-group"}))
		}, []string{"logGroupIdentifier", "policyDocument", "lastUpdatedTime"}, nil},
		{"PutQueryDefinition", func() (interface{}, error) {
			return svc.PutQueryDefinition(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"name": "vocab-qd-2", "queryString": "fields @message",
			}))
		}, []string{"queryDefinitionId"}, nil},
		{"DescribeQueryDefinitions", func() (interface{}, error) {
			return svc.DescribeQueryDefinitions(ctx, reqCtx, vocabRequest(nil))
		}, []string{"queryDefinitions", "nextToken"}, nil},
		{"DeleteQueryDefinition", func() (interface{}, error) {
			return svc.DeleteQueryDefinition(ctx, reqCtx, vocabRequest(map[string]interface{}{"queryDefinitionId": vocabQueryDefinitionId(t, svc, reqCtx)}))
		}, []string{"success"}, nil},
		{"GetLogGroupFields", func() (interface{}, error) {
			return svc.GetLogGroupFields(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group"}))
		}, []string{"logGroupFields"}, nil},
		{"GetLogRecord", func() (interface{}, error) {
			return svc.GetLogRecord(ctx, reqCtx, vocabRequest(map[string]interface{}{"logRecordPointer": ptr()}))
		}, []string{"logRecord"}, nil},
		// GetLogObject is absent from this census: its response is the
		// modelled event stream, not a JSON key map — the stream's wire
		// contract is pinned by TestGetLogObjectFieldStreamUnion and the
		// SDK's GetLogObject_Stream.
		{"GetLogFields", func() (interface{}, error) {
			return svc.GetLogFields(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"dataSourceName": "vocab-group", "dataSourceType": "AWS::Logs::LogGroup",
			}))
		}, []string{"logFields"}, nil},
		{"CreateExportTask", func() (interface{}, error) {
			return svc.CreateExportTask(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupName": "vocab-group", "destination": "vocab-bucket",
				"from": now - 3600000, "to": now,
			}))
		}, []string{"taskId"}, nil},
		{"DescribeExportTasks", func() (interface{}, error) {
			return svc.DescribeExportTasks(ctx, reqCtx, vocabRequest(nil))
		}, []string{"exportTasks", "nextToken"}, nil},
		{"CreateImportTask", func() (interface{}, error) {
			return svc.CreateImportTask(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"importSourceArn": "arn:aws:s3:::vocab-import-bucket-2",
				"importRoleArn":   "arn:aws:iam::000000000000:role/vocab-import",
			}))
		}, []string{"importId", "importDestinationArn", "creationTime"}, nil},
		{"DescribeImportTasks", func() (interface{}, error) {
			return svc.DescribeImportTasks(ctx, reqCtx, vocabRequest(nil))
		}, []string{"imports", "nextToken"}, nil},
		{"CancelImportTask", func() (interface{}, error) {
			return svc.CancelImportTask(ctx, reqCtx, vocabRequest(map[string]interface{}{"importId": vocabSeedImport(t, svc, "vocab-cancel-import")}))
		}, []string{"importId", "importStatistics", "importStatus", "creationTime", "lastUpdatedTime"}, nil},
		{"DescribeImportTaskBatches", func() (interface{}, error) {
			importId := vocabSeedImport(t, svc, "vocab-batches-import")
			return svc.DescribeImportTaskBatches(ctx, reqCtx, vocabRequest(map[string]interface{}{"importId": importId}))
		}, []string{"importSourceArn", "importId", "importBatches", "nextToken"}, nil},
		{"StartQuery", func() (interface{}, error) {
			return svc.StartQuery(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"startTime": now/1000 - 3600, "endTime": now / 1000,
				"logGroupNames": []interface{}{"vocab-group"}, "queryString": "fields @message | limit 1",
			}))
		}, []string{"queryId"}, nil},
		{"StopQuery", func() (interface{}, error) {
			return svc.StopQuery(ctx, reqCtx, vocabRequest(map[string]interface{}{"queryId": vocabSeedQuery(t, svc, "vocab-stop")}))
		}, []string{"success"}, nil},
		{"DescribeQueries", func() (interface{}, error) {
			return svc.DescribeQueries(ctx, reqCtx, vocabRequest(nil))
		}, []string{"queries", "nextToken"}, nil},
		{"GetQueryResults", func() (interface{}, error) {
			return svc.GetQueryResults(ctx, reqCtx, vocabRequest(map[string]interface{}{"queryId": vocabCompletedQuery(t, svc)}))
		}, []string{"queryLanguage", "results", "statistics", "status", "encryptionKey", "nextToken"}, nil},
		{"PutStorageTierPolicy", func() (interface{}, error) {
			return svc.PutStorageTierPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{"storageTier": "STANDARD"}))
		}, []string{"storageTier", "lastUpdatedTime"}, nil},
		{"GetStorageTierPolicy", func() (interface{}, error) {
			return svc.GetStorageTierPolicy(ctx, reqCtx, vocabRequest(nil))
		}, []string{"storageTier", "lastUpdatedTime"}, nil},
		{"ListLogGroupsForQuery", func() (interface{}, error) {
			return svc.ListLogGroupsForQuery(ctx, reqCtx, vocabRequest(map[string]interface{}{"queryId": vocabCompletedQuery(t, svc)}))
		}, []string{"logGroupIdentifiers", "nextToken"}, nil},
		{"ListAggregateLogGroupSummaries", func() (interface{}, error) {
			return svc.ListAggregateLogGroupSummaries(ctx, reqCtx, vocabRequest(map[string]interface{}{"groupBy": "DATA_SOURCE_NAME_AND_TYPE"}))
		}, []string{"aggregateLogGroupSummaries", "nextToken"}, nil},
		{"CreateScheduledQuery", func() (interface{}, error) {
			return svc.CreateScheduledQuery(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"name": "vocab-sq-2", "queryString": "fields @message",
				"queryLanguage": "CWLI", "executionRoleArn": "arn:aws:iam::000000000000:role/vocab",
				"scheduleExpression": "rate(1 hour)",
			}))
		}, []string{"scheduledQueryArn", "state"}, nil},
		{"GetScheduledQuery", func() (interface{}, error) {
			return svc.GetScheduledQuery(ctx, reqCtx, vocabRequest(map[string]interface{}{"identifier": vocabScheduledQueryARN}))
		}, []string{"scheduledQueryArn", "name", "description", "queryLanguage", "queryString", "logGroupIdentifiers", "scheduleExpression", "timezone", "startTimeOffset", "endTimeOffset", "destinationConfiguration", "state", "scheduleType", "lastTriggeredTime", "lastExecutionStatus", "scheduleStartTime", "scheduleEndTime", "executionRoleArn", "creationTime", "lastUpdatedTime"}, nil},
		{"GetScheduledQueryHistory", func() (interface{}, error) {
			// Both window members are required on this operation; the
			// fixture supplies a wide window over the execution records.
			return svc.GetScheduledQueryHistory(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"identifier": vocabScheduledQueryARN,
				"startTime":  time.Now().Add(-24 * time.Hour).UnixMilli(),
				"endTime":    time.Now().Add(time.Hour).UnixMilli(),
			}))
		}, []string{"name", "scheduledQueryArn", "triggerHistory", "nextToken"}, nil},
		{"ListScheduledQueries", func() (interface{}, error) {
			return svc.ListScheduledQueries(ctx, reqCtx, vocabRequest(nil))
		}, []string{"nextToken", "scheduledQueries"}, nil},
		{"UpdateScheduledQuery", func() (interface{}, error) {
			return svc.UpdateScheduledQuery(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"identifier":         vocabScheduledQueryARN,
				"queryLanguage":      "CWLI",
				"queryString":        "fields @message",
				"scheduleExpression": "rate(1 hour)",
				"executionRoleArn":   "arn:aws:iam::000000000000:role/vocab",
				"description":        "vocab-updated",
			}))
		}, []string{"scheduledQueryArn", "name", "description", "queryLanguage", "queryString", "logGroupIdentifiers", "scheduleExpression", "timezone", "startTimeOffset", "endTimeOffset", "destinationConfiguration", "state", "scheduleType", "lastTriggeredTime", "lastExecutionStatus", "scheduleStartTime", "scheduleEndTime", "executionRoleArn", "creationTime", "lastUpdatedTime"}, nil},
		{"DeleteScheduledQuery", func() (interface{}, error) {
			delResp, err := svc.CreateScheduledQuery(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"name": "vocab-sq-del", "queryString": "fields @message", "queryLanguage": "CWLI",
				"executionRoleArn":   "arn:aws:iam::000000000000:role/vocab",
				"scheduleExpression": "rate(1 hour)",
			}))
			if err != nil {
				t.Fatalf("create throwaway scheduled query: %v", err)
			}
			return svc.DeleteScheduledQuery(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"identifier": delResp.(map[string]interface{})["scheduledQueryArn"],
			}))
		}, []string{}, nil},
		{"CreateLookupTable", func() (interface{}, error) {
			return svc.CreateLookupTable(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"lookupTableName": "vocab_lt_2",
				"tableBody":       "id,name\n1,alpha\n",
			}))
		}, []string{"lookupTableArn", "createdAt"}, nil},
		{"GetLookupTable", func() (interface{}, error) {
			return svc.GetLookupTable(ctx, reqCtx, vocabRequest(map[string]interface{}{"lookupTableArn": "vocab_lt"}))
		}, []string{"lookupTableArn", "lookupTableName", "description", "tableBody", "sizeBytes", "lastUpdatedTime", "kmsKeyId"}, nil},
		{"UpdateLookupTable", func() (interface{}, error) {
			return svc.UpdateLookupTable(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"lookupTableArn": "arn:aws:logs:us-east-1:000000000000:lookup-table:vocab_lt",
				"tableBody":      "id,name\n1,beta\n",
			}))
		}, []string{"lookupTableArn", "lastUpdatedTime"}, nil},
		{"DescribeLookupTables", func() (interface{}, error) {
			return svc.DescribeLookupTables(ctx, reqCtx, vocabRequest(nil))
		}, []string{"lookupTables", "nextToken"}, nil},
		{"ListTagsForResource", func() (interface{}, error) {
			return svc.ListTagsForResource(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"resourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group",
			}))
		}, []string{"tags"}, nil},
		{"ListTagsLogGroup", func() (interface{}, error) {
			return svc.ListTagsLogGroup(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group"}))
		}, []string{"tags"}, nil},
	}

	// The delivery-family fixtures: a source over the fixture group, a
	// CWL destination over a second group, the destination policy and a
	// delivery pairing them.
	if vstore, err := svc.getLogsStoreByRegion("us-east-1"); err != nil {
		t.Fatalf("vocab delivery store: %v", err)
	} else if err := vstore.CreateLogGroup(logsstore.NewLogGroup(deliveryTestDestGroup, "us-east-1", "000000000000")); err != nil {
		t.Fatalf("vocab delivery group: %v", err)
	}
	vocabDeliveryDest := putCWLTestDestination(t, svc)
	vocabDeliverySource := putDeliverySourceOn(t, svc, "vocab-group")
	vocabDeliveryDestArn := vocabDeliveryDest.Arn
	vocabDeliverySourceName := vocabDeliverySource.Name
	if _, err := svc.PutDeliveryDestinationPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"deliveryDestinationName":   "vended-dest",
		"deliveryDestinationPolicy": `{"Version":"2012-10-17","Statement":[]}`,
	})); err != nil {
		t.Fatalf("fixture delivery destination policy: %v", err)
	}
	vocabDeliveryCreate, err := svc.CreateDelivery(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"deliverySourceName":     vocabDeliverySourceName,
		"deliveryDestinationArn": vocabDeliveryDestArn,
	}))
	if err != nil {
		t.Fatalf("fixture delivery: %v", err)
	}
	vocabDeliveryId := vocabDeliveryCreate.(map[string]interface{})["delivery"].(map[string]interface{})["id"].(string)

	cases = append(cases, []struct {
		op       string
		invoke   func() (interface{}, error)
		model    []string
		required []string
	}{
		{"PutDeliverySource", func() (interface{}, error) {
			return svc.PutDeliverySource(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"name":        "vocab-source-2",
				"resourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group",
				"logType":     "APPLICATION_LOGS",
			}))
		}, []string{"deliverySource"}, nil},
		{"DescribeDeliverySources", func() (interface{}, error) {
			return svc.DescribeDeliverySources(ctx, reqCtx, vocabRequest(nil))
		}, []string{"deliverySources", "nextToken"}, nil},
		{"GetDeliverySource", func() (interface{}, error) {
			return svc.GetDeliverySource(ctx, reqCtx, vocabRequest(map[string]interface{}{"name": vocabDeliverySourceName}))
		}, []string{"deliverySource"}, nil},
		{"PutDeliveryDestination", func() (interface{}, error) {
			return svc.PutDeliveryDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"name": "vocab-dest-2",
				"deliveryDestinationConfiguration": map[string]interface{}{
					"destinationResourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group",
				},
			}))
		}, []string{"deliveryDestination"}, nil},
		{"DescribeDeliveryDestinations", func() (interface{}, error) {
			return svc.DescribeDeliveryDestinations(ctx, reqCtx, vocabRequest(nil))
		}, []string{"deliveryDestinations", "nextToken"}, nil},
		{"GetDeliveryDestination", func() (interface{}, error) {
			return svc.GetDeliveryDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{"name": "vended-dest"}))
		}, []string{"deliveryDestination"}, nil},
		{"PutDeliveryDestinationPolicy", func() (interface{}, error) {
			return svc.PutDeliveryDestinationPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"deliveryDestinationName":   "vended-dest",
				"deliveryDestinationPolicy": `{"Version":"2012-10-17","Statement":[]}`,
			}))
		}, []string{"policy"}, nil},
		{"GetDeliveryDestinationPolicy", func() (interface{}, error) {
			return svc.GetDeliveryDestinationPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"deliveryDestinationName": "vended-dest",
			}))
		}, []string{"policy"}, nil},
		{"CreateDelivery", func() (interface{}, error) {
			return svc.CreateDelivery(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"deliverySourceName":     "vocab-source-2",
				"deliveryDestinationArn": "arn:aws:logs:us-east-1:000000000000:delivery-destination:vocab-dest-2",
			}))
		}, []string{"delivery"}, nil},
		{"DescribeDeliveries", func() (interface{}, error) {
			return svc.DescribeDeliveries(ctx, reqCtx, vocabRequest(nil))
		}, []string{"deliveries", "nextToken"}, nil},
		{"GetDelivery", func() (interface{}, error) {
			return svc.GetDelivery(ctx, reqCtx, vocabRequest(map[string]interface{}{"id": vocabDeliveryId}))
		}, []string{"delivery"}, nil},
		{"DescribeConfigurationTemplates", func() (interface{}, error) {
			return svc.DescribeConfigurationTemplates(ctx, reqCtx, vocabRequest(nil))
		}, []string{"configurationTemplates", "nextToken"}, nil},
	}...)

	for _, tc := range cases {
		result, err := tc.invoke()
		if err != nil {
			t.Fatalf("%s: %v", tc.op, err)
		}
		m, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("%s: response type %T, want a map", tc.op, result)
		}
		assertVocabulary(t, tc.op, m, tc.model, tc.required)
	}
}

// TestEmittedKeyVocabularyUnitOutputs drives the Unit-output operations:
// each returns an empty body — the model types their outputs as Unit.
func TestEmittedKeyVocabularyUnitOutputs(t *testing.T) {
	svc, reqCtx := newVocabularyEnv(t)
	ctx := context.Background()
	now := time.Now().UnixMilli()

	unit := []struct {
		name   string
		invoke func() (interface{}, error)
	}{
		{"CreateLogGroup", func() (interface{}, error) {
			return svc.CreateLogGroup(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-unit-group"}))
		}},
		{"CreateLogStream", func() (interface{}, error) {
			return svc.CreateLogStream(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group", "logStreamName": "vocab-unit-stream"}))
		}},
		{"PutMetricFilter", func() (interface{}, error) {
			return svc.PutMetricFilter(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupName": "vocab-group", "filterName": "vocab-unit-mf", "filterPattern": "{ level = INFO }",
				"metricTransformations": []interface{}{map[string]interface{}{
					"metricName": "VocabUnitErrors", "metricNamespace": "Vocab", "metricValue": "1",
				}},
			}))
		}},
		{"PutSubscriptionFilter", func() (interface{}, error) {
			return svc.PutSubscriptionFilter(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupName": "vocab-group", "filterName": "vocab-unit-sf", "filterPattern": "",
				"destinationArn": "arn:aws:logs:us-east-1:000000000000:destination:vocab-target",
			}))
		}},
		{"PutRetentionPolicy", func() (interface{}, error) {
			return svc.PutRetentionPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-unit-group", "retentionInDays": 30}))
		}},
		{"DeleteRetentionPolicy", func() (interface{}, error) {
			return svc.DeleteRetentionPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-unit-group"}))
		}},
		{"AssociateKmsKey", func() (interface{}, error) {
			return svc.AssociateKmsKey(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupName": "vocab-unit-group", "kmsKeyId": "arn:aws:kms:us-east-1:000000000000:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			}))
		}},
		{"DisassociateKmsKey", func() (interface{}, error) {
			return svc.DisassociateKmsKey(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupName": "vocab-unit-group",
			}))
		}},
		{"PutLogGroupDeletionProtection", func() (interface{}, error) {
			return svc.PutLogGroupDeletionProtection(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupIdentifier": "vocab-unit-group", "deletionProtectionEnabled": true,
			}))
		}},
		{"PutBearerTokenAuthentication", func() (interface{}, error) {
			return svc.PutBearerTokenAuthentication(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupIdentifier": "vocab-unit-group", "bearerTokenAuthenticationEnabled": true,
			}))
		}},
		{"PutDestinationPolicy", func() (interface{}, error) {
			return svc.PutDestinationPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"destinationName": "vocab-dest", "accessPolicy": `{"Version":"2012-10-17"}`,
			}))
		}},
		{"TagResource", func() (interface{}, error) {
			return svc.TagResource(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"resourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group",
				"tags":        map[string]interface{}{"env": "vocab"},
			}))
		}},
		{"UntagResource", func() (interface{}, error) {
			return svc.UntagResource(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"resourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group",
				"tagKeys":     []interface{}{"env"},
			}))
		}},
		{"TagLogGroup", func() (interface{}, error) {
			return svc.TagLogGroup(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupName": "vocab-group", "tags": map[string]interface{}{"env": "vocab"},
			}))
		}},
		{"UntagLogGroup", func() (interface{}, error) {
			return svc.UntagLogGroup(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupName": "vocab-group", "tags": []interface{}{"env"},
			}))
		}},
		{"CancelExportTask", func() (interface{}, error) {
			return svc.CancelExportTask(ctx, reqCtx, vocabRequest(map[string]interface{}{"taskId": vocabSeedExport(t, svc)}))
		}},
		{"DeleteMetricFilter", func() (interface{}, error) {
			return svc.DeleteMetricFilter(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group", "filterName": "vocab-unit-mf"}))
		}},
		{"DeleteSubscriptionFilter", func() (interface{}, error) {
			return svc.DeleteSubscriptionFilter(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group", "filterName": "vocab-unit-sf"}))
		}},
		{"DeleteLogStream", func() (interface{}, error) {
			return svc.DeleteLogStream(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group", "logStreamName": "vocab-unit-stream"}))
		}},
		{"DeleteLogGroup", func() (interface{}, error) {
			if _, err := svc.CreateLogGroup(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-unit-group-2"})); err != nil {
				t.Fatalf("create throwaway group: %v", err)
			}
			return svc.DeleteLogGroup(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-unit-group-2"}))
		}},
		{"UpdateDeliveryConfiguration", func() (interface{}, error) {
			delivery := vocabSeedDelivery(t, svc, reqCtx)
			return svc.UpdateDeliveryConfiguration(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"id": delivery, "recordFields": []interface{}{"time", "message"},
			}))
		}},
		{"DeleteDeliveryDestinationPolicy", func() (interface{}, error) {
			if _, err := svc.PutDeliveryDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"name": "vocab-unit-policy-dest",
				"deliveryDestinationConfiguration": map[string]interface{}{
					"destinationResourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group",
				},
			})); err != nil {
				t.Fatalf("create throwaway delivery destination: %v", err)
			}
			if _, err := svc.PutDeliveryDestinationPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"deliveryDestinationName":   "vocab-unit-policy-dest",
				"deliveryDestinationPolicy": `{"Version":"2012-10-17","Statement":[]}`,
			})); err != nil {
				t.Fatalf("seed throwaway policy: %v", err)
			}
			return svc.DeleteDeliveryDestinationPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"deliveryDestinationName": "vocab-unit-policy-dest",
			}))
		}},
		{"DeleteDelivery", func() (interface{}, error) {
			return svc.DeleteDelivery(ctx, reqCtx, vocabRequest(map[string]interface{}{"id": vocabSeedDelivery(t, svc, reqCtx)}))
		}},
		{"DeleteDeliverySource", func() (interface{}, error) {
			if _, err := svc.PutDeliverySource(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"name": "vocab-unit-src", "resourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group", "logType": "APPLICATION_LOGS",
			})); err != nil {
				t.Fatalf("create throwaway delivery source: %v", err)
			}
			return svc.DeleteDeliverySource(ctx, reqCtx, vocabRequest(map[string]interface{}{"name": "vocab-unit-src"}))
		}},
		{"DeleteDeliveryDestination", func() (interface{}, error) {
			if _, err := svc.PutDeliveryDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"name": "vocab-unit-dest",
				"deliveryDestinationConfiguration": map[string]interface{}{
					"destinationResourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group",
				},
			})); err != nil {
				t.Fatalf("create throwaway delivery destination: %v", err)
			}
			return svc.DeleteDeliveryDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{"name": "vocab-unit-dest"}))
		}},
		{"DeleteDestination", func() (interface{}, error) {
			if _, err := svc.PutDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"destinationName": "vocab-unit-del-dest",
				"roleArn":         "arn:aws:iam::000000000000:role/vocab",
				"targetArn":       "arn:aws:kinesis:us-east-1:000000000000:stream/vocab-unit-del-target",
			})); err != nil {
				t.Fatalf("create throwaway destination: %v", err)
			}
			return svc.DeleteDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"destinationName": "vocab-unit-del-dest",
			}))
		}},
		{"DeleteResourcePolicy", func() (interface{}, error) {
			if _, err := svc.PutResourcePolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"policyName": "vocab-unit-del-rp", "policyDocument": `{"Version":"2012-10-17"}`,
			})); err != nil {
				t.Fatalf("create throwaway resource policy: %v", err)
			}
			return svc.DeleteResourcePolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"policyName": "vocab-unit-del-rp",
			}))
		}},
		{"DeleteAccountPolicy", func() (interface{}, error) {
			if _, err := svc.PutAccountPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"policyName": "vocab-unit-del-ap", "policyDocument": `{"Fields":["A"]}`,
				"policyType": "FIELD_INDEX_POLICY", "selectionCriteria": "LogGroupNamePrefix = vocab-unit-",
			})); err != nil {
				t.Fatalf("create throwaway account policy: %v", err)
			}
			return svc.DeleteAccountPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"policyName": "vocab-unit-del-ap", "policyType": "FIELD_INDEX_POLICY",
			}))
		}},
		{"DeleteDataProtectionPolicy", func() (interface{}, error) {
			if _, err := svc.CreateLogGroup(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-unit-dpp-group"})); err != nil {
				t.Fatalf("create throwaway group: %v", err)
			}
			if _, err := svc.PutDataProtectionPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupIdentifier": "vocab-unit-dpp-group",
				"policyDocument": `{"Name":"vocab-unit-dpp","Version":"2021-06-01","Statement":[` +
					`{"DataIdentifer":["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],"Operation":{"Audit":{"FindingsDestination":{"CloudWatchLogs":{"LogGroup":"vocab-unit-findings"}}}}},` +
					`{"DataIdentifer":["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],"Operation":{"Deidentify":{"MaskConfig":{}}}}]}`,
			})); err != nil {
				t.Fatalf("create throwaway data protection policy: %v", err)
			}
			return svc.DeleteDataProtectionPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"logGroupIdentifier": "vocab-unit-dpp-group",
			}))
		}},
		{"DeleteLookupTable", func() (interface{}, error) {
			if _, err := svc.CreateLookupTable(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"lookupTableName": "vocab_unit_del",
				"tableBody":       "id,name\n1,alpha\n",
			})); err != nil {
				t.Fatalf("create throwaway lookup table: %v", err)
			}
			return svc.DeleteLookupTable(ctx, reqCtx, vocabRequest(map[string]interface{}{
				"lookupTableArn": "vocab_unit_del",
			}))
		}},
	}
	for _, tc := range unit {
		result, err := tc.invoke()
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		m, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("%s: response type %T, want a map", tc.name, result)
		}
		if len(m) != 0 {
			t.Fatalf("%s: Unit-output response carries members %v, want none", tc.name, vocabKeysOf(m))
		}
	}
	_ = now
}

// TestEmittedKeyVocabularyNestedShapes pins the nested item shapes this
// programme reshaped: the FilteredLogEvent with its eventId, the
// LogFieldsListItem key, the ExportTask without its invented top-level
// creationTime, the QueryInfo without errorMessage, and the log record's
// @timestamp remaining the model's string form.
func TestEmittedKeyVocabularyNestedShapes(t *testing.T) {
	svc, reqCtx := newVocabularyEnv(t)
	ctx := context.Background()
	now := time.Now().UnixMilli()

	filterResp, err := svc.FilterLogEvents(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group"}))
	if err != nil {
		t.Fatalf("filter log events: %v", err)
	}
	events := filterResp.(map[string]interface{})["events"].([]map[string]interface{})
	if len(events) == 0 {
		t.Fatal("filter log events returned no events")
	}
	for k := range events[0] {
		switch k {
		case "timestamp", "message", "ingestionTime", "logStreamName", "eventId":
		default:
			t.Fatalf("FilteredLogEvent carries the unmodelled member %q", k)
		}
	}
	if events[0]["eventId"] == "" {
		t.Fatal("FilteredLogEvent omitted eventId")
	}

	// The GetLogEvents output shape is OutputLogEvent, whose members are
	// exactly {ingestionTime, message, timestamp} — logStreamName and
	// eventId belong to FilteredLogEvent alone.
	getResp, err := svc.GetLogEvents(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group", "logStreamName": "s1"}))
	if err != nil {
		t.Fatalf("get log events: %v", err)
	}
	getEvents := getResp.(map[string]interface{})["events"].([]map[string]interface{})
	if len(getEvents) == 0 {
		t.Fatal("get log events returned no events")
	}
	for k := range getEvents[0] {
		switch k {
		case "timestamp", "message", "ingestionTime":
		default:
			t.Fatalf("OutputLogEvent carries the unmodelled member %q", k)
		}
	}

	// ListLogGroups serves the LogGroupSummary shape: logGroupArn,
	// logGroupClass and logGroupName, nothing else.
	listResp, err := svc.ListLogGroups(ctx, reqCtx, vocabRequest(map[string]interface{}{}))
	if err != nil {
		t.Fatalf("list log groups: %v", err)
	}
	summaries := listResp.(map[string]interface{})["logGroups"].([]map[string]interface{})
	if len(summaries) == 0 {
		t.Fatal("list log groups returned no summaries")
	}
	for k := range summaries[0] {
		switch k {
		case "logGroupArn", "logGroupClass", "logGroupName":
		default:
			t.Fatalf("LogGroupSummary carries the unmodelled member %q", k)
		}
	}
	if summaries[0]["logGroupClass"] == "" {
		t.Fatal("LogGroupSummary omitted logGroupClass")
	}

	fieldsResp, err := svc.GetLogFields(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"dataSourceName": "vocab-group", "dataSourceType": "AWS::Logs::LogGroup",
	}))
	if err != nil {
		t.Fatalf("get log fields: %v", err)
	}
	items := fieldsResp.(map[string]interface{})["logFields"].([]map[string]interface{})
	for k := range items[0] {
		if k != "logFieldName" {
			t.Fatalf("LogFieldsListItem carries the unmodelled member %q", k)
		}
	}

	exportResp, err := svc.CreateExportTask(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"logGroupName": "vocab-group", "destination": "vocab-bucket",
		"from": now - 3600000, "to": now,
	}))
	if err != nil {
		t.Fatalf("create export task: %v", err)
	}
	taskId := exportResp.(map[string]interface{})["taskId"].(string)
	describeResp, err := svc.DescribeExportTasks(ctx, reqCtx, vocabRequest(map[string]interface{}{"taskId": taskId}))
	if err != nil {
		t.Fatalf("describe export tasks: %v", err)
	}
	tasks := describeResp.(map[string]interface{})["exportTasks"].([]map[string]interface{})
	if len(tasks) == 0 {
		t.Fatal("describe export tasks returned no tasks")
	}
	for k := range tasks[0] {
		switch k {
		case "taskId", "taskName", "logGroupName", "from", "to", "destination", "destinationPrefix", "status", "executionInfo":
		default:
			t.Fatalf("ExportTask carries the unmodelled member %q", k)
		}
	}

	vocabCompletedQuery(t, svc)
	queriesResp, err := svc.DescribeQueries(ctx, reqCtx, vocabRequest(nil))
	if err != nil {
		t.Fatalf("describe queries: %v", err)
	}
	queries := queriesResp.(map[string]interface{})["queries"].([]map[string]interface{})
	if len(queries) == 0 {
		t.Fatal("describe queries returned no queries")
	}
	for k := range queries[0] {
		switch k {
		case "queryId", "queryString", "status", "createTime", "logGroupName", "queryLanguage", "queryDuration", "bytesScanned", "userIdentity":
		default:
			t.Fatalf("QueryInfo carries the unmodelled member %q", k)
		}
	}

	vocabSeedPointerEvent(t, svc, reqCtx, ctx, now)
	ptr := eventPointer("vocab-group", "s1", now-2000, `{"level": "INFO", "service": "vocab"}`)
	recordResp, err := svc.GetLogRecord(ctx, reqCtx, vocabRequest(map[string]interface{}{"logRecordPointer": ptr}))
	if err != nil {
		t.Fatalf("get log record: %v", err)
	}
	record := recordResp.(map[string]interface{})["logRecord"].(map[string]interface{})
	if _, ok := record["@timestamp"].(string); !ok {
		t.Fatalf("logRecord @timestamp type %T, want the model's string form", record["@timestamp"])
	}
}

// --- fixtures ---

// vocabSeedPointerEvent seeds the record a GetLogRecord/GetLogObject
// pointer will address: the resolution validates the pointer against the
// store, so the addressed event must exist at exactly the caller's
// timestamp.
func vocabSeedPointerEvent(t *testing.T, svc *LogsService, reqCtx *request.RequestContext, ctx context.Context, now int64) {
	t.Helper()
	if _, err := svc.PutLogEvents(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"logGroupName": "vocab-group", "logStreamName": "s1",
		"logEvents": []interface{}{map[string]interface{}{"timestamp": float64(now - 2000), "message": `{"level": "INFO", "service": "vocab"}`}},
	})); err != nil {
		t.Fatalf("fixture pointer event: %v", err)
	}
}

func vocabRequest(params map[string]interface{}) *request.ParsedRequest {
	if params == nil {
		params = map[string]interface{}{}
	}
	return &request.ParsedRequest{Operation: "Vocab", Parameters: params}
}

func newVocabularyEnv(t *testing.T) (*LogsService, *request.RequestContext) {
	t.Helper()
	svc, _ := newTestService(t)
	reqCtx := request.NewRequestContext(context.Background(), svc.storageManager, "000000000000", "us-east-1")
	ctx := context.Background()
	now := time.Now().UnixMilli()

	// The fixture stack: group + stream + events, metric and subscription
	// filters, destination, data protection policy, resource policy,
	// query definition, lookup table and scheduled query.
	if _, err := svc.CreateLogGroup(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group"})); err != nil {
		t.Fatalf("fixture create group: %v", err)
	}
	if _, err := svc.CreateLogStream(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": "vocab-group", "logStreamName": "s1"})); err != nil {
		t.Fatalf("fixture create stream: %v", err)
	}
	if _, err := svc.PutLogEvents(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"logGroupName": "vocab-group", "logStreamName": "s1",
		"logEvents": []interface{}{
			// JSON-decoded wire numbers arrive as float64 — the parser's
			// expected form.
			map[string]interface{}{"timestamp": float64(now - 2000), "message": `{"level": "INFO", "service": "vocab"}`},
			map[string]interface{}{"timestamp": float64(now - 1000), "message": "plain vocab message"},
		},
	})); err != nil {
		t.Fatalf("fixture put events: %v", err)
	}
	if _, err := svc.PutMetricFilter(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"logGroupName": "vocab-group", "filterName": "vocab-mf", "filterPattern": "{ level = INFO }",
		"metricTransformations": []interface{}{map[string]interface{}{
			"metricName": "VocabErrors", "metricNamespace": "Vocab", "metricValue": "1",
		}},
	})); err != nil {
		t.Fatalf("fixture put metric filter: %v", err)
	}
	if _, err := svc.PutSubscriptionFilter(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"logGroupName": "vocab-group", "filterName": "vocab-sf", "filterPattern": "",
		"destinationArn": "arn:aws:logs:us-east-1:000000000000:destination:vocab-target",
	})); err != nil {
		t.Fatalf("fixture put subscription filter: %v", err)
	}
	if _, err := svc.PutDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"destinationName": "vocab-dest",
		"roleArn":         "arn:aws:iam::000000000000:role/vocab",
		"targetArn":       "arn:aws:kinesis:us-east-1:000000000000:stream/vocab-target",
	})); err != nil {
		t.Fatalf("fixture put destination: %v", err)
	}
	if _, err := svc.PutDataProtectionPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"logGroupIdentifier": "vocab-group",
		"policyDocument": `{"Name":"vocab-dpp","Version":"2021-06-01","Statement":[` +
			`{"DataIdentifer":["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],"Operation":{"Audit":{"FindingsDestination":{"CloudWatchLogs":{"LogGroup":"vocab-findings"}}}}},` +
			`{"DataIdentifer":["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],"Operation":{"Deidentify":{"MaskConfig":{}}}}]}`,
	})); err != nil {
		t.Fatalf("fixture put dpp: %v", err)
	}
	if _, err := svc.PutResourcePolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"policyName": "vocab-rp", "policyDocument": `{"Version":"2012-10-17"}`,
	})); err != nil {
		t.Fatalf("fixture put resource policy: %v", err)
	}
	if _, err := svc.PutQueryDefinition(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"name": "vocab-qd", "queryString": "fields @message",
	})); err != nil {
		t.Fatalf("fixture put query definition: %v", err)
	}
	if _, err := svc.CreateLookupTable(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"lookupTableName": "vocab_lt", "tableBody": "id,name\n1,alpha\n",
	})); err != nil {
		t.Fatalf("fixture create lookup table: %v", err)
	}
	sqResp, err := svc.CreateScheduledQuery(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"name": "vocab-sq", "queryString": "fields @message", "queryLanguage": "CWLI",
		"executionRoleArn":   "arn:aws:iam::000000000000:role/vocab",
		"scheduleExpression": "rate(1 hour)",
	}))
	if err != nil {
		t.Fatalf("fixture create scheduled query: %v", err)
	}
	vocabScheduledQueryARN = sqResp.(map[string]interface{})["scheduledQueryArn"].(string)
	return svc, reqCtx
}

// vocabScheduledQueryARN holds the fixture scheduled query's ARN (the
// identifier form the get/history paths resolve).
var vocabScheduledQueryARN string

// vocabQueryDefinitionId creates a definition and returns its generated
// id, for the DeleteQueryDefinition case.
func vocabQueryDefinitionId(t *testing.T, svc *LogsService, reqCtx *request.RequestContext) string {
	t.Helper()
	resp, err := svc.PutQueryDefinition(context.Background(), reqCtx, vocabRequest(map[string]interface{}{
		"name": "vocab-qd-del", "queryString": "fields @message",
	}))
	if err != nil {
		t.Fatalf("create throwaway query definition: %v", err)
	}
	return resp.(map[string]interface{})["queryDefinitionId"].(string)
}

// vocabSeedQuery stores a Running query state for StopQuery to cancel.
func vocabSeedQuery(t *testing.T, svc *LogsService, id string) string {
	t.Helper()
	qs := &queryState{
		queryId:       id,
		logGroupNames: []string{"vocab-group"},
		queryString:   "fields @message",
		queryLanguage: "CWLI",
		status:        queryStatusRunning,
		createdAt:     time.Now(),
	}
	svc.queries.Store(id, qs)
	return id
}

// vocabCompletedQuery starts a real query and waits for its completion.
func vocabCompletedQuery(t *testing.T, svc *LogsService) string {
	t.Helper()
	now := time.Now().UnixMilli()
	queryId, err := svc.startQueryCore(&StartQueryInput{
		StartTime:      now/1000 - 3600,
		EndTime:        now / 1000,
		QueryString:    "fields @message | limit 1",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		LogGroupNames:  []string{"vocab-group"},
		Region:         "us-east-1",
	})
	if err != nil {
		t.Fatalf("start query: %v", err)
	}
	return waitTerminalQuery(t, svc, queryId).queryId
}

// vocabSeedExport writes a RUNNING export task for CancelExportTask.
func vocabSeedExport(t *testing.T, svc *LogsService) string {
	t.Helper()
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	task := &logsstore.ExportTask{
		TaskId:       "vocab-export-cancel",
		LogGroupName: "vocab-group",
		Status:       "RUNNING",
		From:         time.Now().UnixMilli() - 3600000,
		To:           time.Now().UnixMilli(),
	}
	if err := store.PutExportTask(task); err != nil {
		t.Fatal(err)
	}
	return task.TaskId
}

// vocabSeedImport writes an IN_PROGRESS import task for the cancel and
// batches paths.
func vocabSeedImport(t *testing.T, svc *LogsService, importId string) string {
	t.Helper()
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	task := &logsstore.ImportTask{
		ImportId:        importId,
		ImportSourceArn: "arn:aws:s3:::vocab-import-bucket",
		ImportStatus:    "IN_PROGRESS",
		CreationTime:    time.Now().UnixMilli(),
		LastUpdatedTime: time.Now().UnixMilli(),
	}
	if err := store.PutImportTask(task); err != nil {
		t.Fatal(err)
	}
	return importId
}

func assertVocabulary(t *testing.T, name string, m map[string]interface{}, model, required []string) {
	t.Helper()
	modelSet := make(map[string]bool, len(model))
	for _, k := range model {
		modelSet[k] = true
	}
	for _, k := range vocabKeysOf(m) {
		if !modelSet[k] {
			t.Fatalf("%s emitted the unmodelled member %q", name, k)
		}
	}
	for _, k := range required {
		if _, ok := m[k]; !ok {
			t.Fatalf("%s omitted the required member %q", name, k)
		}
	}
}

func vocabKeysOf(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// TestEventPointerRoundTripWithPipeCharacters pins the pointer codec: a
// log stream name legally containing the delimiter (only ':' and '*' are
// forbidden) and a message containing it both survive the
// encode→splitPointer→unescape round trip without field shift.
func TestEventPointerRoundTripWithPipeCharacters(t *testing.T) {
	svc, reqCtx := newVocabularyEnv(t)
	ctx := context.Background()
	group := "vocab-pipe-group"
	msg := "before|middle|after"
	// The pointer resolution validates the addressed record against the
	// store, so each round trip seeds its stream with the real event.
	ts := time.Now().UnixMilli() - 2000
	if _, err := svc.CreateLogGroup(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": group})); err != nil {
		t.Fatalf("fixture create group: %v", err)
	}
	seed := func(stream string) {
		t.Helper()
		if _, err := svc.CreateLogStream(ctx, reqCtx, vocabRequest(map[string]interface{}{"logGroupName": group, "logStreamName": stream})); err != nil {
			t.Fatalf("fixture create stream %s: %v", stream, err)
		}
		if _, err := svc.PutLogEvents(ctx, reqCtx, vocabRequest(map[string]interface{}{
			"logGroupName": group, "logStreamName": stream,
			"logEvents": []interface{}{map[string]interface{}{"timestamp": float64(ts), "message": msg}},
		})); err != nil {
			t.Fatalf("fixture put events %s: %v", stream, err)
		}
	}

	stream := "pipe|stream|name"
	seed(stream)
	ptr := eventPointer(group, stream, ts, msg)

	resp, err := svc.GetLogRecord(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"logRecordPointer": ptr,
	}))
	if err != nil {
		t.Fatalf("get log record: %v", err)
	}
	record := resp.(map[string]interface{})["logRecord"].(map[string]interface{})
	if record["@logStream"] != stream {
		t.Fatalf("@logStream = %q, want %q", record["@logStream"], stream)
	}
	if record["@message"] != msg {
		t.Fatalf("@message = %q, want %q", record["@message"], msg)
	}
	if record["@log"] != "000000000000:"+group {
		t.Fatalf("@log = %q", record["@log"])
	}
	if record["@timestamp"] != fmt.Sprintf("%d", ts) {
		t.Fatalf("@timestamp = %q, want the epoch-millis string", record["@timestamp"])
	}

	// A pointer to a stream whose name itself carries an escape-worthy
	// backslash also round-trips.
	bsStream := `back\slash`
	seed(bsStream)
	ptr = eventPointer(group, bsStream, ts, msg)
	resp, err = svc.GetLogRecord(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"logRecordPointer": ptr,
	}))
	if err != nil {
		t.Fatalf("get log record (backslash): %v", err)
	}
	record = resp.(map[string]interface{})["logRecord"].(map[string]interface{})
	if record["@logStream"] != bsStream {
		t.Fatalf("@logStream = %q, want %q", record["@logStream"], bsStream)
	}
}

// TestGetLogGroupFieldsTimeWindow pins the modelled time semantics: the
// request's time member is epoch SECONDS centring a ±8-minute window, and
// an omitted time scans the most recent 15 minutes.
func TestGetLogGroupFieldsTimeWindow(t *testing.T) {
	svc, store := newReadTestService(t, "fields-group")
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "fields-group")); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	// One JSON event inside the last 15 minutes, one an hour old.
	if _, err := store.PutLogEvents("fields-group", "s1", []logsstore.LogEntry{
		{Timestamp: now - 60_000, Message: `{"recent_field": 1}`, IngestionTime: now},
		{Timestamp: now - 3_600_000, Message: `{"old_field": 1}`, IngestionTime: now - 3_600_000},
	}); err != nil {
		t.Fatal(err)
	}

	hasField := func(fields []map[string]interface{}, name string) bool {
		for _, f := range fields {
			if f["name"] == name {
				return true
			}
		}
		return false
	}

	// Omitted time: the most recent 15 minutes — the recent event only.
	fields, err := svc.getLogGroupFieldsCore(store, "fields-group", 0)
	if err != nil {
		t.Fatal(err)
	}
	if !hasField(fields, "recent_field") {
		t.Fatal("default 15-minute window missed the recent event's field")
	}
	if hasField(fields, "old_field") {
		t.Fatal("default 15-minute window included the hour-old event's field")
	}

	// time in SECONDS centred on the hour-old event: a ±8-minute window
	// around it sees only that event.
	fields, err = svc.getLogGroupFieldsCore(store, "fields-group", (now-3_600_000)/1000)
	if err != nil {
		t.Fatal(err)
	}
	if !hasField(fields, "old_field") {
		t.Fatal("seconds-centred ±8-minute window missed the old event's field")
	}
	if hasField(fields, "recent_field") {
		t.Fatal("seconds-centred ±8-minute window included the recent event's field")
	}
}

// vocabSeedDelivery pairs a fresh delivery source over the fixture group
// with the fixture destination and returns the new delivery's id.
func vocabSeedDelivery(t *testing.T, svc *LogsService, reqCtx *request.RequestContext) string {
	t.Helper()
	ctx := context.Background()
	seed := fmt.Sprintf("vocab-seed-%d", time.Now().UnixNano()%1e9)
	if _, err := svc.PutDeliveryDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"name": seed,
		"deliveryDestinationConfiguration": map[string]interface{}{
			"destinationResourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group",
		},
	})); err != nil {
		t.Fatalf("seed delivery destination: %v", err)
	}
	if _, err := svc.PutDeliverySource(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"name": seed, "resourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group", "logType": "APPLICATION_LOGS",
	})); err != nil {
		t.Fatalf("seed delivery source: %v", err)
	}
	resp, err := svc.CreateDelivery(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"deliverySourceName":     seed,
		"deliveryDestinationArn": "arn:aws:logs:us-east-1:000000000000:delivery-destination:" + seed,
	}))
	if err != nil {
		t.Fatalf("seed delivery: %v", err)
	}
	return resp.(map[string]interface{})["delivery"].(map[string]interface{})["id"].(string)
}

// TestEmittedKeyVocabularyDeliveryShapes pins the delivery family's
// nested record shapes: every emitted member of DeliverySource,
// DeliveryDestination, Delivery, the destination Policy and
// ConfigurationTemplate belongs to the model's member set for that
// shape.
func TestEmittedKeyVocabularyDeliveryShapes(t *testing.T) {
	svc, reqCtx := newVocabularyEnv(t)
	ctx := context.Background()

	if vstore, err := svc.getLogsStoreByRegion("us-east-1"); err != nil {
		t.Fatalf("delivery shape store: %v", err)
	} else if err := vstore.CreateLogGroup(logsstore.NewLogGroup(deliveryTestDestGroup, "us-east-1", "000000000000")); err != nil {
		t.Fatalf("delivery shape group: %v", err)
	}

	srcResp, err := svc.PutDeliverySource(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"name": "vended-source", "resourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:vocab-group",
		"logType": "APPLICATION_LOGS", "tags": map[string]interface{}{"env": "vocab"},
	}))
	if err != nil {
		t.Fatalf("put delivery source: %v", err)
	}
	source := srcResp.(map[string]interface{})["deliverySource"].(map[string]interface{})
	for k := range source {
		switch k {
		case "name", "arn", "resourceArns", "service", "logType", "tags", "deliverySourceConfiguration", "status", "statusReason":
		default:
			t.Fatalf("DeliverySource carries the unmodelled member %q", k)
		}
	}
	if source["status"] != "ACTIVE" {
		t.Fatalf("DeliverySource status: %v", source["status"])
	}

	destResp, err := svc.PutDeliveryDestination(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"name": "vended-dest", "outputFormat": "json",
		"deliveryDestinationConfiguration": map[string]interface{}{
			"destinationResourceArn": "arn:aws:logs:us-east-1:000000000000:log-group:" + deliveryTestDestGroup,
		},
	}))
	if err != nil {
		t.Fatalf("put delivery destination: %v", err)
	}
	dest := destResp.(map[string]interface{})["deliveryDestination"].(map[string]interface{})
	for k := range dest {
		switch k {
		case "name", "arn", "deliveryDestinationType", "outputFormat", "deliveryDestinationConfiguration", "tags":
		default:
			t.Fatalf("DeliveryDestination carries the unmodelled member %q", k)
		}
	}
	inner := dest["deliveryDestinationConfiguration"].(map[string]interface{})
	if _, ok := inner["destinationResourceArn"]; !ok {
		t.Fatal("DeliveryDestinationConfiguration omitted destinationResourceArn")
	}
	for k := range inner {
		if k != "destinationResourceArn" {
			t.Fatalf("DeliveryDestinationConfiguration carries the unmodelled member %q", k)
		}
	}

	deliveryResp, err := svc.CreateDelivery(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"deliverySourceName":      "vended-source",
		"deliveryDestinationArn":  dest["arn"],
		"recordFields":            []interface{}{"time", "message"},
		"s3DeliveryConfiguration": map[string]interface{}{},
	}))
	if err == nil {
		t.Fatal("s3DeliveryConfiguration accepted on a CWL delivery")
	}
	deliveryResp, err = svc.CreateDelivery(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"deliverySourceName":     "vended-source",
		"deliveryDestinationArn": dest["arn"],
		"recordFields":           []interface{}{"time", "message"},
	}))
	if err != nil {
		t.Fatalf("create delivery: %v", err)
	}
	delivery := deliveryResp.(map[string]interface{})["delivery"].(map[string]interface{})
	for k := range delivery {
		switch k {
		case "id", "arn", "deliverySourceName", "deliveryDestinationArn", "deliveryDestinationType", "recordFields", "fieldDelimiter", "s3DeliveryConfiguration", "tags":
		default:
			t.Fatalf("Delivery carries the unmodelled member %q", k)
		}
	}

	policyResp, err := svc.PutDeliveryDestinationPolicy(ctx, reqCtx, vocabRequest(map[string]interface{}{
		"deliveryDestinationName":   "vended-dest",
		"deliveryDestinationPolicy": `{"Version":"2012-10-17","Statement":[]}`,
	}))
	if err != nil {
		t.Fatalf("put delivery destination policy: %v", err)
	}
	policy := policyResp.(map[string]interface{})["policy"].(map[string]interface{})
	for k := range policy {
		if k != "deliveryDestinationPolicy" {
			t.Fatalf("Policy carries the unmodelled member %q", k)
		}
	}

	templatesResp, err := svc.DescribeConfigurationTemplates(ctx, reqCtx, vocabRequest(nil))
	if err != nil {
		t.Fatalf("describe configuration templates: %v", err)
	}
	rows := templatesResp.(map[string]interface{})["configurationTemplates"].([]map[string]interface{})
	if len(rows) == 0 {
		t.Fatal("configuration template catalog is empty")
	}
	for k := range rows[0] {
		switch k {
		case "service", "logType", "resourceType", "deliveryDestinationType", "defaultDeliveryConfigValues", "allowedFields", "allowedOutputFormats", "allowedActionForAllowVendedLogsDeliveryForResource", "allowedFieldDelimiters", "allowedSuffixPathFields", "deliverySourceConfiguration", "s3TablesIntegration":
		default:
			t.Fatalf("ConfigurationTemplate carries the unmodelled member %q", k)
		}
	}
	defaults := rows[0]["defaultDeliveryConfigValues"].(map[string]interface{})
	for k := range defaults {
		switch k {
		case "recordFields", "fieldDelimiter", "s3DeliveryConfiguration":
		default:
			t.Fatalf("ConfigurationTemplateDeliveryConfigValues carries the unmodelled member %q", k)
		}
	}
}

// GetLogGroupFields' percent is the sampled ratio — "The percentage of
// log events queried that contained the field" — not a constant: a
// field in three of four sampled events reports 75.
func TestGetLogGroupFieldsPercentRatio(t *testing.T) {
	svc, store := newReadTestService(t, "fields-ratio-group")
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "fields-ratio-group")); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	if _, err := store.PutLogEvents("fields-ratio-group", "s1", []logsstore.LogEntry{
		{Timestamp: now - 40_000, Message: `{"shared":1,"always":1}`},
		{Timestamp: now - 30_000, Message: `{"shared":2,"always":2}`},
		{Timestamp: now - 20_000, Message: `{"shared":3,"always":3}`},
		{Timestamp: now - 10_000, Message: `{"always":4}`},
	}); err != nil {
		t.Fatal(err)
	}

	fields, err := svc.getLogGroupFieldsCore(store, "fields-ratio-group", 0)
	if err != nil {
		t.Fatal(err)
	}
	percent := func(name string) int {
		for _, f := range fields {
			if f["name"] == name {
				p, _ := f["percent"].(int)
				return p
			}
		}
		return -1
	}
	if got := percent("shared"); got != 75 {
		t.Fatalf("shared percent = %d, want 75 (three of four)", got)
	}
	if got := percent("always"); got != 100 {
		t.Fatalf("always percent = %d, want 100", got)
	}
	if got := percent("@message"); got != 100 {
		t.Fatalf("core field percent = %d, want 100", got)
	}
}

// GetLogObject's fieldStream is the modelled @streaming union: the wire
// form is one awsJson1.1 event stream — an initial-response message,
// then one fields event whose payload is the FieldsData document with
// the record as its base64 data blob — served with the eventstream
// content type.
func TestGetLogObjectFieldStreamUnion(t *testing.T) {
	svc, store := newReadTestService(t, "log-object-group")
	reqCtx := request.NewRequestContext(context.Background(), svc.storageManager, "000000000000", "us-east-1")
	// The pointer resolution validates the addressed record against the
	// store, so the fixture seeds the real event.
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "log-object-group")); err != nil {
		t.Fatal(err)
	}
	ts := time.Now().UnixMilli() - 2000
	if _, err := store.PutLogEvents("log-object-group", "s1", []logsstore.LogEntry{{Timestamp: ts, Message: `{"k":"v"}`}}); err != nil {
		t.Fatal(err)
	}
	pointer := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("log-object-group|s1|%d|{\"k\":\"v\"}", ts)))

	resp, err := svc.GetLogObject(context.Background(), reqCtx, vocabRequest(map[string]interface{}{
		"logObjectPointer": pointer,
	}))
	if err != nil {
		t.Fatal(err)
	}
	stream, ok := resp.(*logObjectEventStream)
	if !ok {
		t.Fatalf("GetLogObject must answer with the event-stream response, got %T", resp)
	}
	if ct := stream.GetStreamHeaders().Get("Content-Type"); ct != "application/vnd.amazon.eventstream" {
		t.Fatalf("content type = %q, want the event-stream one", ct)
	}
	initial, err := readEventStreamFrame(stream.GetStream())
	if err != nil {
		t.Fatalf("read initial-response frame: %v", err)
	}
	if initial.headers[":event-type"] != "initial-response" {
		t.Fatalf("first frame = %v, want initial-response", initial.headers)
	}
	fieldsFrame, err := readEventStreamFrame(stream.GetStream())
	if err != nil {
		t.Fatalf("read fields frame: %v", err)
	}
	if fieldsFrame.headers[":event-type"] != "fields" {
		t.Fatalf("second frame = %v, want the fields event", fieldsFrame.headers)
	}
	var fieldsData struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(fieldsFrame.payload, &fieldsData); err != nil {
		t.Fatalf("fields payload is not the FieldsData document: %v (%s)", err, fieldsFrame.payload)
	}
	blob, err := base64.StdEncoding.DecodeString(fieldsData.Data)
	if err != nil {
		t.Fatalf("data blob is not base64: %v", err)
	}
	var record map[string]interface{}
	if err := json.Unmarshal(blob, &record); err != nil {
		t.Fatalf("data blob is not the record JSON: %v (%q)", err, string(blob))
	}
	if record["k"] != "v" || record["@message"] != `{"k":"v"}` {
		t.Fatalf("record = %v", record)
	}
	// The stream is finite: the fields event is its last message.
	if trailing, err := readEventStreamFrame(stream.GetStream()); err == nil {
		t.Fatalf("stream delivered a trailing frame %v past the fields event", trailing.headers)
	}
}
