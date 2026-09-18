package testutil

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runCloudTrailEDSTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	var edsID string

	// Create EDS for lifecycle tests.
	results = append(results, r.RunTest("cloudtrail", "CreateEventDataStore_Lifecycle", func() error {
		resp, err := tc.createEventDataStore("ct-eds-lifecycle", aws.Int32(90))
		if err != nil {
			return fmt.Errorf("CreateEventDataStore failed: %w", err)
		}
		if resp.EventDataStoreArn == nil {
			return fmt.Errorf("EventDataStoreArn is nil")
		}
		edsID = tc.edsIDFromARN(aws.ToString(resp.EventDataStoreArn))
		return nil
	}))

	// Get EDS and verify initial state.
	results = append(results, r.RunTest("cloudtrail", "GetEventDataStore_InitialState", func() error {
		resp, err := tc.client.GetEventDataStore(tc.ctx, &cloudtrail.GetEventDataStoreInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return fmt.Errorf("GetEventDataStore failed: %w", err)
		}
		if resp.Status != types.EventDataStoreStatusEnabled {
			return fmt.Errorf("expected ENABLED status, got %s", resp.Status)
		}
		return nil
	}))

	// An ARN-prefixed selector that is not an event data store ARN is the
	// declared EventDataStoreARNInvalidException, not a not-found.
	results = append(results, r.RunTest("cloudtrail", "GetEventDataStore_InvalidARN", func() error {
		_, err := tc.client.GetEventDataStore(tc.ctx, &cloudtrail.GetEventDataStoreInput{
			EventDataStore: aws.String("arn:aws:s3:::not-an-event-data-store"),
		})
		return AssertErrorContains(err, "EventDataStoreARNInvalidException")
	}))

	// Stop ingestion: the store lands in STOPPED_INGESTION, the status the
	// start direction requires.
	results = append(results, r.RunTest("cloudtrail", "StopEventDataStoreIngestion", func() error {
		_, err := tc.client.StopEventDataStoreIngestion(tc.ctx, &cloudtrail.StopEventDataStoreIngestionInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return fmt.Errorf("StopEventDataStoreIngestion failed: %w", err)
		}
		desc, err := tc.client.GetEventDataStore(tc.ctx, &cloudtrail.GetEventDataStoreInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return fmt.Errorf("GetEventDataStore after stop failed: %w", err)
		}
		if desc.Status != types.EventDataStoreStatusStoppedIngestion {
			return fmt.Errorf("after stop: status = %s, want STOPPED_INGESTION", desc.Status)
		}
		return nil
	}))

	// Start ingestion: the store returns to ENABLED.
	results = append(results, r.RunTest("cloudtrail", "StartEventDataStoreIngestion", func() error {
		_, err := tc.client.StartEventDataStoreIngestion(tc.ctx, &cloudtrail.StartEventDataStoreIngestionInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return fmt.Errorf("StartEventDataStoreIngestion failed: %w", err)
		}
		desc, err := tc.client.GetEventDataStore(tc.ctx, &cloudtrail.GetEventDataStoreInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return fmt.Errorf("GetEventDataStore after start failed: %w", err)
		}
		if desc.Status != types.EventDataStoreStatusEnabled {
			return fmt.Errorf("after start: status = %s, want ENABLED", desc.Status)
		}
		return nil
	}))

	// Delete EDS (soft delete -> PENDING_DELETION).
	results = append(results, r.RunTest("cloudtrail", "DeleteEventDataStore_SoftDelete", func() error {
		_, err := tc.client.DeleteEventDataStore(tc.ctx, &cloudtrail.DeleteEventDataStoreInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return fmt.Errorf("DeleteEventDataStore failed: %w", err)
		}
		return nil
	}))

	// Verify PENDING_DELETION status.
	results = append(results, r.RunTest("cloudtrail", "DeleteEDS_VerifyPending", func() error {
		resp, err := tc.client.GetEventDataStore(tc.ctx, &cloudtrail.GetEventDataStoreInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return err
		}
		if resp.Status != types.EventDataStoreStatusPendingDeletion {
			return fmt.Errorf("expected PENDING_DELETION, got %s", resp.Status)
		}
		return nil
	}))

	// Restore EDS.
	results = append(results, r.RunTest("cloudtrail", "RestoreEventDataStore", func() error {
		_, err := tc.client.RestoreEventDataStore(tc.ctx, &cloudtrail.RestoreEventDataStoreInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return fmt.Errorf("RestoreEventDataStore failed: %w", err)
		}
		return nil
	}))

	// Verify restored.
	results = append(results, r.RunTest("cloudtrail", "RestoreEDS_Verify", func() error {
		resp, err := tc.client.GetEventDataStore(tc.ctx, &cloudtrail.GetEventDataStoreInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return err
		}
		if resp.Status != types.EventDataStoreStatusEnabled {
			return fmt.Errorf("expected ENABLED after restore, got %s", resp.Status)
		}
		return nil
	}))

	// EDS name validation.
	results = append(results, r.RunTest("cloudtrail", "CreateEDS_InvalidName", func() error {
		_, err := tc.client.CreateEventDataStore(tc.ctx, &cloudtrail.CreateEventDataStoreInput{
			Name: aws.String("ab"), // too short (min 3)
		})
		if err == nil {
			return fmt.Errorf("expected error for short EDS name")
		}
		return nil
	}))

	// RetentionPeriod validation.
	results = append(results, r.RunTest("cloudtrail", "CreateEDS_InvalidRetention", func() error {
		_, err := tc.client.CreateEventDataStore(tc.ctx, &cloudtrail.CreateEventDataStoreInput{
			Name:            aws.String("ct-eds-bad-retention"),
			RetentionPeriod: aws.Int32(1), // too low (min 7)
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid retention period")
		}
		return nil
	}))

	// UpdateEventDataStore requires at least one optional parameter.
	results = append(results, r.RunTest("cloudtrail", "UpdateEventDataStore_NoParameters", func() error {
		resp, err := tc.createEventDataStore("ct-eds-noparam", nil)
		if err != nil {
			return err
		}
		defer tc.deleteEventDataStore(*resp.EventDataStoreArn)

		_, err = tc.client.UpdateEventDataStore(tc.ctx, &cloudtrail.UpdateEventDataStoreInput{
			EventDataStore: resp.EventDataStoreArn,
		})
		return AssertErrorContains(err, "InvalidParameterException")
	}))

	// The KMS key of an event data store cannot be changed once associated.
	results = append(results, r.RunTest("cloudtrail", "UpdateEventDataStore_KmsImmutable", func() error {
		resp, err := tc.client.CreateEventDataStore(tc.ctx, &cloudtrail.CreateEventDataStoreInput{
			Name:     aws.String(tc.uniqueName("ct-eds-kms")),
			KmsKeyId: aws.String("alias/MatrixKey"),
		})
		if err != nil {
			return fmt.Errorf("CreateEventDataStore with KMS key failed: %w", err)
		}
		defer tc.deleteEventDataStore(*resp.EventDataStoreArn)

		// Re-asserting the associated key is a no-op success.
		if _, err = tc.client.UpdateEventDataStore(tc.ctx, &cloudtrail.UpdateEventDataStoreInput{
			EventDataStore: resp.EventDataStoreArn,
			KmsKeyId:       aws.String("alias/MatrixKey"),
		}); err != nil {
			return fmt.Errorf("re-asserting the associated KMS key failed: %w", err)
		}

		_, err = tc.client.UpdateEventDataStore(tc.ctx, &cloudtrail.UpdateEventDataStoreInput{
			EventDataStore: resp.EventDataStoreArn,
			KmsKeyId:       aws.String("alias/OtherKey"),
		})
		return AssertErrorContains(err, "OperationNotPermittedException")
	}))

	// Clean up: delete the restored EDS.
	results = append(results, r.RunTest("cloudtrail", "DeleteEDS_Cleanup", func() error {
		return tc.deleteEventDataStore(edsID)
	}))

	return results
}

func (r *TestRunner) runCloudTrailChannelTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	var channelARN string

	// Create channel. The source is unique per run: a maximum of one
	// channel exists per source, so a leaked channel from an earlier run
	// must not hold the fixed source value. Destinations must resolve to
	// stored event data stores, so the group provisions one first.
	var destEDSARN string
	results = append(results, r.RunTest("cloudtrail", "Channel_Setup_EDS", func() error {
		resp, err := tc.createEventDataStore("ct-ch-dest", nil)
		if err != nil {
			return fmt.Errorf("create destination EDS: %w", err)
		}
		destEDSARN = *resp.EventDataStoreArn
		return nil
	}))
	defer func() {
		if destEDSARN != "" {
			_ = tc.deleteEventDataStore(destEDSARN)
		}
	}()

	channelSource := tc.uniqueName("ct-src")
	// The channel's name is unique per run: a channel leaked by a failed
	// earlier run must not poison the fixed-name uniqueness checks.
	channelName := tc.uniqueName("ct-channel")
	results = append(results, r.RunTest("cloudtrail", "CreateChannel_Success", func() error {
		resp, err := tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String(channelName),
			Source: aws.String(channelSource),
			Destinations: []types.Destination{
				{
					Type:     types.DestinationTypeEventDataStore,
					Location: aws.String(destEDSARN),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("CreateChannel failed: %w", err)
		}
		if resp.ChannelArn == nil {
			return fmt.Errorf("ChannelArn is nil")
		}
		channelARN = *resp.ChannelArn
		return nil
	}))

	// Channel destination validation: the SDK validates length traits
	// client-side (name length, destination count), so the server-side
	// matrix pin runs at unit level; the enum and pattern violations on a
	// well-sized destination list reach the server and are pinned here.
	results = append(results, r.RunTest("cloudtrail", "CreateChannel_Validation", func() error {
		_, err := tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String(tc.uniqueName("ct-ch-badtype")),
			Source: aws.String(tc.uniqueName("ct-src-bad")),
			Destinations: []types.Destination{
				{Type: types.DestinationType("KINESIS"), Location: aws.String(destEDSARN)},
			},
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return err
		}

		_, err = tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String(tc.uniqueName("ct-ch-badloc")),
			Source: aws.String(tc.uniqueName("ct-src-bad2")),
			Destinations: []types.Destination{
				{Type: types.DestinationTypeEventDataStore, Location: aws.String("bad location!")},
			},
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return err
		}

		// A syntactically valid location that is not an event data store
		// ARN this account owns answers the model's ARN-invalid shape.
		_, err = tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String(tc.uniqueName("ct-ch-notarn")),
			Source: aws.String(tc.uniqueName("ct-src-bad3")),
			Destinations: []types.Destination{
				{Type: types.DestinationTypeEventDataStore, Location: aws.String("not-an-arn")},
			},
		})
		if err := AssertErrorContains(err, "EventDataStoreARNInvalidException"); err != nil {
			return err
		}

		_, err = tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String(tc.uniqueName("ct-ch-missing")),
			Source: aws.String(tc.uniqueName("ct-src-bad4")),
			Destinations: []types.Destination{
				{Type: types.DestinationTypeEventDataStore,
					Location: aws.String(fmt.Sprintf("arn:aws:cloudtrail:%s:%s:eventdatastore/00000000-0000-0000-0000-000000000000",
						tc.region, tc.accountID))},
			},
		})
		return AssertErrorContains(err, "EventDataStoreNotFoundException")
	}))

	// A maximum of one channel is allowed per source.
	results = append(results, r.RunTest("cloudtrail", "CreateChannel_DuplicateSource", func() error {
		source := tc.uniqueName("ct-src-dup")
		first, err := tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String(tc.uniqueName("ct-ch-dup1")),
			Source: aws.String(source),
			Destinations: []types.Destination{
				{Type: types.DestinationTypeEventDataStore, Location: aws.String(destEDSARN)},
			},
		})
		if err != nil {
			return fmt.Errorf("first create failed: %w", err)
		}
		defer tc.client.DeleteChannel(tc.ctx, &cloudtrail.DeleteChannelInput{Channel: first.ChannelArn})

		_, err = tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String(tc.uniqueName("ct-ch-dup2")),
			Source: aws.String(source),
			Destinations: []types.Destination{
				{Type: types.DestinationTypeEventDataStore, Location: aws.String(destEDSARN)},
			},
		})
		return AssertErrorContains(err, "InvalidSourceException")
	}))

	// Get channel.
	results = append(results, r.RunTest("cloudtrail", "GetChannel_Success", func() error {
		resp, err := tc.client.GetChannel(tc.ctx, &cloudtrail.GetChannelInput{
			Channel: aws.String(channelARN),
		})
		if err != nil {
			return fmt.Errorf("GetChannel failed: %w", err)
		}
		if aws.ToString(resp.Name) != channelName {
			return fmt.Errorf("expected name %s, got %s", channelName, aws.ToString(resp.Name))
		}
		return nil
	}))

	// The channel is addressable by its bare UUID ("The ARN or UUID of a
	// channel", GetChannel).
	results = append(results, r.RunTest("cloudtrail", "GetChannel_ByUUID", func() error {
		uuid := channelARN[strings.LastIndex(channelARN, "channel/")+len("channel/"):]
		resp, err := tc.client.GetChannel(tc.ctx, &cloudtrail.GetChannelInput{
			Channel: aws.String(uuid),
		})
		if err != nil {
			return fmt.Errorf("GetChannel by UUID failed: %w", err)
		}
		if aws.ToString(resp.ChannelArn) != channelARN {
			return fmt.Errorf("expected ARN %s, got %s", channelARN, aws.ToString(resp.ChannelArn))
		}
		return nil
	}))

	// An ARN-prefixed selector that is not a channel ARN is the declared
	// ChannelARNInvalidException, not a not-found.
	results = append(results, r.RunTest("cloudtrail", "GetChannel_InvalidARN", func() error {
		_, err := tc.client.GetChannel(tc.ctx, &cloudtrail.GetChannelInput{
			Channel: aws.String("arn:aws:s3:::not-a-channel"),
		})
		return AssertErrorContains(err, "ChannelARNInvalidException")
	}))

	// List channels.
	results = append(results, r.RunTest("cloudtrail", "ListChannels_Success", func() error {
		resp, err := tc.client.ListChannels(tc.ctx, &cloudtrail.ListChannelsInput{})
		if err != nil {
			return fmt.Errorf("ListChannels failed: %w", err)
		}
		found := false
		for _, ch := range resp.Channels {
			if aws.ToString(ch.ChannelArn) == channelARN {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created channel not found in list")
		}
		return nil
	}))

	// Update channel.
	results = append(results, r.RunTest("cloudtrail", "UpdateChannel_Success", func() error {
		_, err := tc.client.UpdateChannel(tc.ctx, &cloudtrail.UpdateChannelInput{
			Channel: aws.String(channelARN),
			Name:    aws.String(tc.uniqueName("ct-channel-updated")),
		})
		if err != nil {
			return fmt.Errorf("UpdateChannel failed: %w", err)
		}
		return nil
	}))

	// CreateChannel with tags (regression: CreateChannel uses "Tags" not "TagsList").
	results = append(results, r.RunTest("cloudtrail", "CreateChannel_WithTags", func() error {
		tagCh := tc.uniqueName("ct-channel-tags")
		resp, err := tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String(tagCh),
			Source: aws.String(tc.uniqueName("ct-src-tags")),
			Destinations: []types.Destination{
				{
					Type:     types.DestinationTypeEventDataStore,
					Location: aws.String(destEDSARN),
				},
			},
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("production")},
				{Key: aws.String("Team"), Value: aws.String("cloud")},
			},
		})
		if err != nil {
			return fmt.Errorf("CreateChannel with tags failed: %w", err)
		}
		defer tc.client.DeleteChannel(tc.ctx, &cloudtrail.DeleteChannelInput{
			Channel: resp.ChannelArn,
		})

		// Verify tags appear in the CreateChannel response.
		if len(resp.Tags) != 2 {
			return fmt.Errorf("expected 2 tags in CreateChannel response, got %d", len(resp.Tags))
		}
		tagMap := tagListToMap(resp.Tags)
		if tagMap["Environment"] != "production" {
			return fmt.Errorf("expected Environment=production, got %s", tagMap["Environment"])
		}
		if tagMap["Team"] != "cloud" {
			return fmt.Errorf("expected Team=cloud, got %s", tagMap["Team"])
		}
		return nil
	}))

	// Delete channel.
	results = append(results, r.RunTest("cloudtrail", "DeleteChannel_Success", func() error {
		_, err := tc.client.DeleteChannel(tc.ctx, &cloudtrail.DeleteChannelInput{
			Channel: aws.String(channelARN),
		})
		if err != nil {
			return fmt.Errorf("DeleteChannel failed: %w", err)
		}
		return nil
	}))

	// Get deleted channel should fail.
	results = append(results, r.RunTest("cloudtrail", "GetChannel_Deleted", func() error {
		_, err := tc.client.GetChannel(tc.ctx, &cloudtrail.GetChannelInput{
			Channel: aws.String(channelARN),
		})
		if err == nil {
			return fmt.Errorf("expected error for deleted channel")
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runCloudTrailQueryTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	// Create an EDS for query tests.
	var edsID string
	results = append(results, r.RunTest("cloudtrail", "Query_CreateEDS", func() error {
		resp, err := tc.createEventDataStore("ct-query-eds", nil)
		if err != nil {
			return err
		}
		edsID = tc.edsIDFromARN(aws.ToString(resp.EventDataStoreArn))
		return nil
	}))

	// StartQuery.
	var queryID string
	results = append(results, r.RunTest("cloudtrail", "StartQuery_Success", func() error {
		stmt := fmt.Sprintf("SELECT eventID, eventTime, eventName FROM %s WHERE eventName = 'CreateEventDataStore'", edsID)
		resp, err := tc.client.StartQuery(tc.ctx, &cloudtrail.StartQueryInput{
			QueryStatement: aws.String(stmt),
		})
		if err != nil {
			return fmt.Errorf("StartQuery failed: %w", err)
		}
		if resp.QueryId == nil {
			return fmt.Errorf("QueryId is nil")
		}
		queryID = *resp.QueryId
		return nil
	}))

	// A syntactically invalid WHERE clause must be rejected at StartQuery
	// time instead of silently degrading to an unfiltered scan.
	results = append(results, r.RunTest("cloudtrail", "StartQuery_MalformedWhereRejected", func() error {
		stmt := fmt.Sprintf("SELECT eventID FROM %s WHERE eventName = 'unclosed", edsID)
		_, err := tc.client.StartQuery(tc.ctx, &cloudtrail.StartQueryInput{
			QueryStatement: aws.String(stmt),
		})
		if err == nil {
			return fmt.Errorf("expected error for malformed WHERE clause")
		}
		var apiErr *types.InvalidQueryStatementException
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected InvalidQueryStatementException, got %v", err)
		}
		return nil
	}))

	// DescribeQuery — poll for FINISHED since StartQuery is asynchronous.
	// A FINISHED query reports honest execution statistics: every event
	// the scan examined and the record bytes it read are non-zero.
	results = append(results, r.RunTest("cloudtrail", "DescribeQuery_Success", func() error {
		var resp *cloudtrail.DescribeQueryOutput
		var err error
		for i := 0; i < 10; i++ {
			resp, err = tc.client.DescribeQuery(tc.ctx, &cloudtrail.DescribeQueryInput{
				EventDataStore: aws.String(edsID),
				QueryId:        aws.String(queryID),
			})
			if err != nil {
				return fmt.Errorf("DescribeQuery failed: %w", err)
			}
			if resp.QueryStatus == types.QueryStatusFinished {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		if resp.QueryStatus != types.QueryStatusFinished {
			return fmt.Errorf("expected FINISHED status, got %s", resp.QueryStatus)
		}
		if resp.QueryStatistics == nil {
			return fmt.Errorf("QueryStatistics is nil")
		}
		if aws.ToInt64(resp.QueryStatistics.EventsScanned) <= 0 {
			return fmt.Errorf("EventsScanned not reported: %d", aws.ToInt64(resp.QueryStatistics.EventsScanned))
		}
		if aws.ToInt64(resp.QueryStatistics.BytesScanned) <= 0 {
			return fmt.Errorf("BytesScanned not reported: %d", aws.ToInt64(resp.QueryStatistics.BytesScanned))
		}
		return nil
	}))

	// GetQueryResults — query should be FINISHED after DescribeQuery poll.
	// The rows must carry VALUES: the requested columns resolve through the
	// Lake column vocabulary with the statement's own spellings echoed as
	// row keys.
	results = append(results, r.RunTest("cloudtrail", "GetQueryResults_Success", func() error {
		resp, err := tc.client.GetQueryResults(tc.ctx, &cloudtrail.GetQueryResultsInput{
			EventDataStore: aws.String(edsID),
			QueryId:        aws.String(queryID),
		})
		if err != nil {
			return fmt.Errorf("GetQueryResults failed: %w", err)
		}
		if len(resp.QueryResultRows) == 0 {
			return fmt.Errorf("expected result rows for the CreateEventDataStore query")
		}
		found := false
		for _, row := range resp.QueryResultRows {
			flat := map[string]string{}
			for _, entry := range row {
				for k, v := range entry {
					flat[k] = v
				}
			}
			if flat["eventID"] == "" {
				return fmt.Errorf("row has empty eventID value: %v", flat)
			}
			if _, err := time.Parse(time.RFC3339, flat["eventTime"]); err != nil {
				return fmt.Errorf("eventTime is not an RFC3339 timestamp: %q", flat["eventTime"])
			}
			if flat["eventName"] == "CreateEventDataStore" {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("no row carries eventName CreateEventDataStore")
		}
		return nil
	}))

	// GetQueryResults pagination — walk every page through NextToken until
	// exhaustion: the collected row count must equal the query's reported
	// total, pages must never repeat, and an unusable token answers with a
	// terminal empty page instead of restarting the walk. The pin starts
	// its own query whose WHERE bounds the result set to this suite's own
	// query-operation audit events — the store's event history grows with
	// every suite run (and, in a full regression, with every other
	// service's audit events), so an unfiltered query is volume-fragile,
	// while the query-operation event names are touched only by this
	// suite and stay well inside the walk's page budget.
	results = append(results, r.RunTest("cloudtrail", "GetQueryResults_Pagination", func() error {
		// Own the result-set size: every ListQueries call is itself audited
		// as a '%Query%' event, so three seeding calls guarantee a
		// multi-page walk at MaxQueryResults=2 even on a fresh store.
		// The eventTime bound keeps the set to this run's own events —
		// query-op events accumulate across suite runs inside the 24-hour
		// test window and would outgrow the page budget otherwise.
		bound := time.Now().UTC().Add(-2 * time.Second)
		for i := 0; i < 3; i++ {
			if _, err := tc.client.ListQueries(tc.ctx, &cloudtrail.ListQueriesInput{EventDataStore: aws.String(edsID)}); err != nil {
				return fmt.Errorf("seeding ListQueries failed: %w", err)
			}
		}
		startResp, err := tc.client.StartQuery(tc.ctx, &cloudtrail.StartQueryInput{
			QueryStatement: aws.String(fmt.Sprintf("SELECT eventID FROM %s WHERE eventName LIKE '%%Query%%' AND eventTime > %d", edsID, bound.Unix())),
		})
		if err != nil {
			return fmt.Errorf("pagination StartQuery failed: %w", err)
		}
		pageQueryID := aws.ToString(startResp.QueryId)
		if pageQueryID == "" {
			return fmt.Errorf("pagination QueryId is nil")
		}
		finished := false
		for i := 0; i < 10; i++ {
			desc, descErr := tc.client.DescribeQuery(tc.ctx, &cloudtrail.DescribeQueryInput{
				EventDataStore: aws.String(edsID),
				QueryId:        aws.String(pageQueryID),
			})
			if descErr != nil {
				return fmt.Errorf("pagination DescribeQuery failed: %w", descErr)
			}
			if desc.QueryStatus == types.QueryStatusFinished {
				finished = true
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		if !finished {
			return fmt.Errorf("pagination query did not reach FINISHED")
		}

		var expected int32 = -1
		totalRows := 0
		pages := 0
		var nextToken *string
		for i := 0; i < 60; i++ {
			resp, err := tc.client.GetQueryResults(tc.ctx, &cloudtrail.GetQueryResultsInput{
				EventDataStore:  aws.String(edsID),
				QueryId:         aws.String(pageQueryID),
				MaxQueryResults: aws.Int32(2),
				NextToken:       nextToken,
			})
			if err != nil {
				return fmt.Errorf("GetQueryResults page %d failed: %w", pages+1, err)
			}
			if resp.QueryStatistics != nil && resp.QueryStatistics.TotalResultsCount != nil {
				expected = *resp.QueryStatistics.TotalResultsCount
			}
			totalRows += len(resp.QueryResultRows)
			pages++
			nextToken = resp.NextToken
			if nextToken == nil || *nextToken == "" {
				break
			}
		}
		if nextToken != nil && *nextToken != "" {
			return fmt.Errorf("pagination did not terminate after %d pages", pages)
		}
		if expected >= 0 && int32(totalRows) != expected {
			return fmt.Errorf("walked %d rows, query reports %d", totalRows, expected)
		}
		if expected < 3 {
			return fmt.Errorf("expected at least the 3 seeded query events, query reports %d", expected)
		}
		if pages < 2 {
			return fmt.Errorf("expected multiple pages for %d results, got %d", expected, pages)
		}

		bad, err := tc.client.GetQueryResults(tc.ctx, &cloudtrail.GetQueryResultsInput{
			EventDataStore: aws.String(edsID),
			QueryId:        aws.String(pageQueryID),
			NextToken:      aws.String("not-a-numeric-offset"),
		})
		if err != nil {
			return fmt.Errorf("GetQueryResults with unusable token failed: %w", err)
		}
		if len(bad.QueryResultRows) != 0 {
			return fmt.Errorf("unusable token returned %d rows, want an empty page", len(bad.QueryResultRows))
		}
		if bad.NextToken != nil && *bad.NextToken != "" {
			return fmt.Errorf("unusable token produced NextToken %q, want a terminal page", *bad.NextToken)
		}
		return nil
	}))

	// ListQueries.
	results = append(results, r.RunTest("cloudtrail", "ListQueries_Success", func() error {
		resp, err := tc.client.ListQueries(tc.ctx, &cloudtrail.ListQueriesInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return fmt.Errorf("ListQueries failed: %w", err)
		}
		found := false
		for _, q := range resp.Queries {
			if aws.ToString(q.QueryId) == queryID {
				found = true
				// The list item is the Query shape: CreationTime is its
				// time member — a missing key deserialises to nil.
				if q.CreationTime == nil {
					return fmt.Errorf("ListQueries item carries no CreationTime")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created query not found in list")
		}
		return nil
	}))

	// ListQueries rejects a MaxResults above the model's 1000 bound with
	// the operation's declared InvalidMaxResultsException.
	results = append(results, r.RunTest("cloudtrail", "ListQueries_MaxResultsBound", func() error {
		_, err := tc.client.ListQueries(tc.ctx, &cloudtrail.ListQueriesInput{
			EventDataStore: aws.String(edsID),
			MaxResults:     aws.Int32(1001),
		})
		return AssertErrorContains(err, "InvalidMaxResultsException")
	}))

	// ListQueries pagination — a second query on the same event data store
	// makes MaxResults=1 span multiple pages; the walk must terminate and
	// never repeat a query.
	results = append(results, r.RunTest("cloudtrail", "ListQueries_Pagination", func() error {
		second, err := tc.client.StartQuery(tc.ctx, &cloudtrail.StartQueryInput{
			QueryStatement: aws.String(fmt.Sprintf("SELECT eventID FROM %s", edsID)),
		})
		if err != nil {
			return fmt.Errorf("second StartQuery failed: %w", err)
		}
		secondID := aws.ToString(second.QueryId)
		if secondID == "" {
			return fmt.Errorf("second QueryId is nil")
		}

		seen := make(map[string]int)
		var nextToken *string
		pages := 0
		for i := 0; i < 20; i++ {
			resp, err := tc.client.ListQueries(tc.ctx, &cloudtrail.ListQueriesInput{
				EventDataStore: aws.String(edsID),
				MaxResults:     aws.Int32(1),
				NextToken:      nextToken,
			})
			if err != nil {
				return fmt.Errorf("ListQueries page %d failed: %w", pages+1, err)
			}
			pages++
			for _, q := range resp.Queries {
				seen[aws.ToString(q.QueryId)]++
			}
			nextToken = resp.NextToken
			if nextToken == nil || *nextToken == "" {
				break
			}
		}
		if nextToken != nil && *nextToken != "" {
			return fmt.Errorf("pagination did not terminate after %d pages", pages)
		}
		if seen[queryID] != 1 {
			return fmt.Errorf("first query appears %d times across pages, want exactly once", seen[queryID])
		}
		if seen[secondID] != 1 {
			return fmt.Errorf("second query appears %d times across pages, want exactly once", seen[secondID])
		}
		if pages < 2 {
			return fmt.Errorf("expected multiple pages for two queries with MaxResults=1, got %d", pages)
		}
		return nil
	}))

	// The dashboard alias form of StartQuery cannot resolve a template
	// until the dashboard operations exist, and the statement form does
	// not mix with it.
	results = append(results, r.RunTest("cloudtrail", "StartQuery_AliasRejected", func() error {
		_, err := tc.client.StartQuery(tc.ctx, &cloudtrail.StartQueryInput{
			QueryAlias:      aws.String("dashboard-query"),
			QueryParameters: []string{edsID},
		})
		if err == nil {
			return fmt.Errorf("expected error for QueryAlias form")
		}
		var apiErr *types.UnsupportedOperationException
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected UnsupportedOperationException, got %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "StartQuery_StatementAndAliasRejected", func() error {
		_, err := tc.client.StartQuery(tc.ctx, &cloudtrail.StartQueryInput{
			QueryStatement: aws.String(fmt.Sprintf("SELECT eventID FROM %s", edsID)),
			QueryAlias:     aws.String("dashboard-query"),
		})
		if err == nil {
			return fmt.Errorf("expected error for QueryStatement with QueryAlias")
		}
		var apiErr *types.InvalidParameterException
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected InvalidParameterException, got %v", err)
		}
		return nil
	}))

	// An owner account that does not own the event data store cannot query
	// it; the response of a successful start names the owning account.
	results = append(results, r.RunTest("cloudtrail", "StartQuery_OwnerAccountMismatch", func() error {
		_, err := tc.client.StartQuery(tc.ctx, &cloudtrail.StartQueryInput{
			QueryStatement:               aws.String(fmt.Sprintf("SELECT eventID FROM %s", edsID)),
			EventDataStoreOwnerAccountId: aws.String("999999999999"),
		})
		if err == nil {
			return fmt.Errorf("expected error for foreign owner account")
		}
		var apiErr *types.EventDataStoreNotFoundException
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected EventDataStoreNotFoundException, got %v", err)
		}
		return nil
	}))

	// DeliveryS3Uri — the finished query's rows are delivered to the named
	// bucket as a gzip CSV plus a sign file whose hashes match the
	// compressed content and whose RSA signature verifies against the
	// ListPublicKeys key with the sign file's fingerprint.
	results = append(results, r.RunTest("cloudtrail", "StartQuery_DeliveryS3Uri", func() error {
		bucket := tc.uniqueName("ct-query-results")
		if _, err := tc.s3Client.CreateBucket(tc.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer tc.s3Client.DeleteBucket(tc.ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})

		start, err := tc.client.StartQuery(tc.ctx, &cloudtrail.StartQueryInput{
			QueryStatement: aws.String(fmt.Sprintf("SELECT eventID, eventName FROM %s WHERE eventName = 'CreateEventDataStore'", edsID)),
			DeliveryS3Uri:  aws.String("s3://" + bucket),
		})
		if err != nil {
			return fmt.Errorf("StartQuery with DeliveryS3Uri failed: %w", err)
		}
		deliveryQueryID := aws.ToString(start.QueryId)
		if deliveryQueryID == "" {
			return fmt.Errorf("delivery QueryId is nil")
		}
		if start.EventDataStoreOwnerAccountId == nil {
			return fmt.Errorf("EventDataStoreOwnerAccountId is nil")
		}

		// Poll to FINISHED with a SUCCESS delivery.
		var desc *cloudtrail.DescribeQueryOutput
		for i := 0; i < 20; i++ {
			desc, err = tc.client.DescribeQuery(tc.ctx, &cloudtrail.DescribeQueryInput{
				EventDataStore: aws.String(edsID),
				QueryId:        aws.String(deliveryQueryID),
			})
			if err != nil {
				return fmt.Errorf("DescribeQuery failed: %w", err)
			}
			if desc.QueryStatus == types.QueryStatusFinished && desc.DeliveryStatus == types.DeliveryStatusSuccess {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		if desc.QueryStatus != types.QueryStatusFinished {
			return fmt.Errorf("delivery query status = %s", desc.QueryStatus)
		}
		if desc.DeliveryStatus != types.DeliveryStatusSuccess {
			return fmt.Errorf("delivery status = %s", desc.DeliveryStatus)
		}
		if aws.ToString(desc.DeliveryS3Uri) != "s3://"+bucket {
			return fmt.Errorf("DeliveryS3Uri echo = %q", aws.ToString(desc.DeliveryS3Uri))
		}

		keys, err := tc.client.ListPublicKeys(tc.ctx, &cloudtrail.ListPublicKeysInput{})
		if err != nil {
			return fmt.Errorf("ListPublicKeys failed: %w", err)
		}

		// Locate the delivered objects under the documented path.
		var signBody, gzBody []byte
		list, err := tc.s3Client.ListObjectsV2(tc.ctx, &s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 failed: %w", err)
		}
		for _, obj := range list.Contents {
			key := aws.ToString(obj.Key)
			switch {
			case strings.HasSuffix(key, "/result_1.csv.gz"):
				gzBody, err = tc.getObjectBytes(bucket, key)
				if err != nil {
					return err
				}
			case strings.HasSuffix(key, "/result_sign.json"):
				signBody, err = tc.getObjectBytes(bucket, key)
				if err != nil {
					return err
				}
			}
		}
		if gzBody == nil || signBody == nil {
			return fmt.Errorf("delivered objects not found in bucket %s", bucket)
		}

		// The gzip CSV carries the header row and the matching event rows.
		zr, err := gzip.NewReader(bytes.NewReader(gzBody))
		if err != nil {
			return fmt.Errorf("result file is not gzip: %w", err)
		}
		csvText, err := io.ReadAll(zr)
		if err != nil {
			return fmt.Errorf("reading result file: %w", err)
		}
		if !strings.HasPrefix(string(csvText), "eventID,eventName") {
			return fmt.Errorf("CSV header row missing: %q", string(csvText)[:min(60, len(csvText))])
		}
		if !strings.Contains(string(csvText), "CreateEventDataStore") {
			return fmt.Errorf("CSV does not carry the queried rows")
		}

		// The sign file hashes and signature verify end to end.
		var sign struct {
			Version              string `json:"version"`
			Region               string `json:"region"`
			HashAlgorithm        string `json:"hashAlgorithm"`
			SignatureAlgorithm   string `json:"signatureAlgorithm"`
			QueryCompleteTime    string `json:"queryCompleteTime"`
			HashSignature        string `json:"hashSignature"`
			PublicKeyFingerprint string `json:"publicKeyFingerprint"`
			Files                []struct {
				FileHashValue string `json:"fileHashValue"`
				FileName      string `json:"fileName"`
			} `json:"files"`
		}
		if err := json.Unmarshal(signBody, &sign); err != nil {
			return fmt.Errorf("sign file is not JSON: %w", err)
		}
		if sign.HashAlgorithm != "SHA-256" || sign.SignatureAlgorithm != "SHA256withRSA" || len(sign.Files) != 1 {
			return fmt.Errorf("sign file fields wrong: %+v", sign)
		}
		sum := sha256.Sum256(gzBody)
		if hex.EncodeToString(sum[:]) != sign.Files[0].FileHashValue {
			return fmt.Errorf("fileHashValue does not match the delivered content")
		}

		sig, err := hex.DecodeString(sign.HashSignature)
		if err != nil {
			return fmt.Errorf("hashSignature is not hexadecimal: %w", err)
		}
		digest := sha256.Sum256([]byte(sign.Files[0].FileHashValue))
		verified := false
		for _, pk := range keys.PublicKeyList {
			if aws.ToString(pk.Fingerprint) != sign.PublicKeyFingerprint {
				continue
			}
			rsaPub, err := x509.ParsePKCS1PublicKey(pk.Value)
			if err != nil {
				return fmt.Errorf("public key value does not parse as PKCS#1 DER: %w", err)
			}
			if err := rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, digest[:], sig); err == nil {
				verified = true
			}
		}
		if !verified {
			return fmt.Errorf("hashSignature does not verify against the fingerprint-matched public key")
		}
		return nil
	}))

	// ListQueries time filters — a window covering the run finds the query,
	// a future window does not.
	results = append(results, r.RunTest("cloudtrail", "ListQueries_TimeFilters", func() error {
		now := time.Now()
		resp, err := tc.client.ListQueries(tc.ctx, &cloudtrail.ListQueriesInput{
			EventDataStore: aws.String(edsID),
			StartTime:      aws.Time(now.Add(-time.Hour)),
			EndTime:        aws.Time(now.Add(time.Hour)),
		})
		if err != nil {
			return fmt.Errorf("ListQueries with time window failed: %w", err)
		}
		found := false
		for _, q := range resp.Queries {
			if aws.ToString(q.QueryId) == queryID {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("bounded window does not include the created query")
		}

		resp, err = tc.client.ListQueries(tc.ctx, &cloudtrail.ListQueriesInput{
			EventDataStore: aws.String(edsID),
			StartTime:      aws.Time(now.Add(time.Hour)),
			EndTime:        aws.Time(now.Add(2 * time.Hour)),
		})
		if err != nil {
			return fmt.Errorf("ListQueries with future window failed: %w", err)
		}
		if len(resp.Queries) != 0 {
			return fmt.Errorf("future window returned %d queries, want none", len(resp.Queries))
		}
		return nil
	}))

	// Clean up EDS.
	results = append(results, r.RunTest("cloudtrail", "Query_CleanupEDS", func() error {
		return tc.deleteEventDataStore(edsID)
	}))

	return results
}

func (r *TestRunner) runCloudTrailConfigTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	// PutEventConfiguration for a trail. The trail must exist: the
	// operation answers TrailNotFoundException for an absent target, so a
	// negative put against a never-created trail is pinned as well.
	results = append(results, r.RunTest("cloudtrail", "PutEventConfiguration_Trail", func() error {
		if _, err := tc.createTrail("ct-event-config-trail", "ct-event-config-bucket"); err != nil {
			return err
		}
		defer tc.deleteTrail("ct-event-config-trail")

		_, err := tc.client.PutEventConfiguration(tc.ctx, &cloudtrail.PutEventConfigurationInput{
			TrailName:    aws.String("no-such-trail"),
			MaxEventSize: types.MaxEventSizeStandard,
		})
		if err == nil {
			return fmt.Errorf("expected PutEventConfiguration for an absent trail to fail")
		}

		_, err = tc.client.PutEventConfiguration(tc.ctx, &cloudtrail.PutEventConfigurationInput{
			TrailName:    aws.String("ct-event-config-trail"),
			MaxEventSize: types.MaxEventSizeStandard,
		})
		if err != nil {
			return fmt.Errorf("PutEventConfiguration failed: %w", err)
		}
		return nil
	}))

	// PutEventConfiguration's response is the persisted configuration —
	// the trail ARN and the applied size echo back. The pin provisions its
	// own trail: the sibling put test's deferred delete removes its trail
	// before any later test runs.
	results = append(results, r.RunTest("cloudtrail", "PutEventConfiguration_Response", func() error {
		if _, err := tc.createTrail("ct-event-config-echo", "ct-event-config-bucket"); err != nil {
			return err
		}
		defer tc.deleteTrail("ct-event-config-echo")

		resp, err := tc.client.PutEventConfiguration(tc.ctx, &cloudtrail.PutEventConfigurationInput{
			TrailName:    aws.String("ct-event-config-echo"),
			MaxEventSize: types.MaxEventSizeLarge,
		})
		if err != nil {
			return fmt.Errorf("PutEventConfiguration failed: %w", err)
		}
		if resp.TrailARN == nil || *resp.TrailARN == "" {
			return fmt.Errorf("PutEventConfiguration returned no TrailARN")
		}
		if resp.MaxEventSize != types.MaxEventSizeLarge {
			return fmt.Errorf("PutEventConfiguration MaxEventSize echo = %s", resp.MaxEventSize)
		}
		return nil
	}))

	// GetEventConfiguration. The pin provisions its own trail: the
	// configuration read must target a live trail (the operation resolves
	// the trail first), and a sibling test's deferred teardown must not
	// decide this test's fixture lifetime.
	results = append(results, r.RunTest("cloudtrail", "GetEventConfiguration_Trail", func() error {
		if _, err := tc.createTrail("ct-event-config-trail", "ct-event-config-bucket"); err != nil {
			return err
		}
		defer tc.deleteTrail("ct-event-config-trail")

		if _, err := tc.client.PutEventConfiguration(tc.ctx, &cloudtrail.PutEventConfigurationInput{
			TrailName:    aws.String("ct-event-config-trail"),
			MaxEventSize: types.MaxEventSizeStandard,
		}); err != nil {
			return fmt.Errorf("PutEventConfiguration failed: %w", err)
		}

		resp, err := tc.client.GetEventConfiguration(tc.ctx, &cloudtrail.GetEventConfigurationInput{
			TrailName: aws.String("ct-event-config-trail"),
		})
		if err != nil {
			return fmt.Errorf("GetEventConfiguration failed: %w", err)
		}
		if resp.MaxEventSize != types.MaxEventSizeStandard {
			return fmt.Errorf("expected MaxEventSize=Standard, got %s", resp.MaxEventSize)
		}
		return nil
	}))

	// RegisterOrganizationDelegatedAdmin. Accounts on this platform never
	// belong to an organization, so the operation answers its documented
	// non-member refusal.
	results = append(results, r.RunTest("cloudtrail", "RegisterDelegatedAdmin", func() error {
		_, err := tc.client.RegisterOrganizationDelegatedAdmin(tc.ctx, &cloudtrail.RegisterOrganizationDelegatedAdminInput{
			MemberAccountId: aws.String(tc.accountID),
		})
		if err == nil {
			return fmt.Errorf("RegisterOrganizationDelegatedAdmin succeeded; want OrganizationsNotInUseException")
		}
		var notInUse *types.OrganizationsNotInUseException
		if !errors.As(err, &notInUse) {
			return fmt.Errorf("RegisterOrganizationDelegatedAdmin: want OrganizationsNotInUseException, got %v", err)
		}
		return nil
	}))

	// DeregisterOrganizationDelegatedAdmin answers the same non-member
	// refusal.
	results = append(results, r.RunTest("cloudtrail", "DeregisterDelegatedAdmin", func() error {
		_, err := tc.client.DeregisterOrganizationDelegatedAdmin(tc.ctx, &cloudtrail.DeregisterOrganizationDelegatedAdminInput{
			DelegatedAdminAccountId: aws.String(tc.accountID),
		})
		if err == nil {
			return fmt.Errorf("DeregisterOrganizationDelegatedAdmin succeeded; want OrganizationsNotInUseException")
		}
		var notInUse *types.OrganizationsNotInUseException
		if !errors.As(err, &notInUse) {
			return fmt.Errorf("DeregisterOrganizationDelegatedAdmin: want OrganizationsNotInUseException, got %v", err)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runCloudTrailValidationTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	// CreateTrail without S3BucketName should fail.
	results = append(results, r.RunTest("cloudtrail", "CreateTrail_NoBucket_Error", func() error {
		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name: aws.String("ct-no-bucket-test"),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing S3BucketName")
		}
		return nil
	}))

	// UpdateTrail routes through the create validators: a bucket name the
	// create path rejects is rejected on update too.
	results = append(results, r.RunTest("cloudtrail", "UpdateTrail_InvalidBucketName", func() error {
		name := tc.uniqueName("ut-bad-bucket")
		defer tc.deleteTrail(name)

		if _, err := tc.createTrail(name, "ut-bad-bucket-store"); err != nil {
			return err
		}
		_, err := tc.client.UpdateTrail(tc.ctx, &cloudtrail.UpdateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("..invalid.."),
		})
		return AssertErrorContains(err, "InvalidS3BucketNameException")
	}))

	// The KmsKeyId alias form is a documented accepted value on create.
	results = append(results, r.RunTest("cloudtrail", "CreateTrail_KmsAliasAccepted", func() error {
		name := tc.uniqueName("kms-alias")
		defer tc.deleteTrail(name)

		if err := tc.ensureTrailBucket("kms-alias-bucket"); err != nil {
			return err
		}
		resp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("kms-alias-bucket"),
			KmsKeyId:     aws.String("alias/MyAliasName"),
		})
		if err != nil {
			return fmt.Errorf("CreateTrail with alias KmsKeyId failed: %w", err)
		}
		if aws.ToString(resp.KmsKeyId) != "alias/MyAliasName" {
			return fmt.Errorf("KmsKeyId echo = %q, want alias/MyAliasName", aws.ToString(resp.KmsKeyId))
		}
		return nil
	}))

	// UpdateTrail clears SnsTopicName by setting empty string. The topic is
	// a real one carrying the documented CloudTrail publish policy — the
	// create-time resolution verifies existence and policy before accepting
	// the destination, and resolves the name to the topic ARN.
	results = append(results, r.RunTest("cloudtrail", "UpdateTrail_ClearField", func() error {
		if err := tc.ensureTrailBucket("test-bucket"); err != nil {
			return err
		}
		topicName := tc.uniqueName("ct-clear-topic")
		createTopic, err := tc.snsClient.CreateTopic(tc.ctx, &sns.CreateTopicInput{Name: aws.String(topicName)})
		if err != nil {
			return fmt.Errorf("CreateTopic failed: %w", err)
		}
		defer func() {
			_, _ = tc.snsClient.DeleteTopic(tc.ctx, &sns.DeleteTopicInput{TopicArn: createTopic.TopicArn})
		}()
		topicPolicy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Sid":"AWSCloudTrailSNSPolicy20131101","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"SNS:Publish","Resource":"%s"}]}`,
			aws.ToString(createTopic.TopicArn))
		if _, err := tc.snsClient.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn:       createTopic.TopicArn,
			AttributeName:  aws.String("Policy"),
			AttributeValue: aws.String(topicPolicy),
		}); err != nil {
			return fmt.Errorf("SetTopicAttributes Policy failed: %w", err)
		}

		// Create trail with SnsTopicName.
		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String("ct-clear-field-test"),
			S3BucketName: aws.String("test-bucket"),
			SnsTopicName: aws.String(topicName),
		})
		if err != nil {
			return fmt.Errorf("CreateTrail failed: %w", err)
		}
		if aws.ToString(createResp.SnsTopicARN) != aws.ToString(createTopic.TopicArn) {
			return fmt.Errorf("create response SnsTopicARN = %q, want the resolved topic ARN %q",
				aws.ToString(createResp.SnsTopicARN), aws.ToString(createTopic.TopicArn))
		}

		// Clear SnsTopicName.
		_, err = tc.client.UpdateTrail(tc.ctx, &cloudtrail.UpdateTrailInput{
			Name:         aws.String("ct-clear-field-test"),
			SnsTopicName: aws.String(""),
		})
		if err != nil {
			return fmt.Errorf("UpdateTrail clear SnsTopicName failed: %w", err)
		}

		// Verify SnsTopicName is cleared.
		getResp, err := tc.client.GetTrail(tc.ctx, &cloudtrail.GetTrailInput{
			Name: aws.String("ct-clear-field-test"),
		})
		if err != nil {
			return err
		}
		if aws.ToString(getResp.Trail.SnsTopicName) != "" {
			return fmt.Errorf("expected empty SnsTopicName after clear, got %s",
				aws.ToString(getResp.Trail.SnsTopicName))
		}
		if aws.ToString(getResp.Trail.SnsTopicARN) != "" {
			return fmt.Errorf("expected empty SnsTopicARN after clear, got %s",
				aws.ToString(getResp.Trail.SnsTopicARN))
		}

		// Cleanup.
		_, _ = tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{
			Name: aws.String("ct-clear-field-test"),
		})
		return nil
	}))

	// DeleteTrail by ARN.
	results = append(results, r.RunTest("cloudtrail", "DeleteTrail_ByARN", func() error {
		if err := tc.ensureTrailBucket("test-bucket"); err != nil {
			return err
		}
		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String("ct-arn-delete-test"),
			S3BucketName: aws.String("test-bucket"),
		})
		if err != nil {
			return err
		}

		// Delete by ARN.
		_, err = tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{
			Name: createResp.TrailARN,
		})
		if err != nil {
			return fmt.Errorf("DeleteTrail by ARN failed: %w", err)
		}

		// Verify it's gone.
		_, err = tc.client.GetTrail(tc.ctx, &cloudtrail.GetTrailInput{
			Name: aws.String("ct-arn-delete-test"),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting trail by ARN")
		}
		return nil
	}))

	return results
}
