package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

func (r *TestRunner) runCloudTrailEDSTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	var edsID string

	// Create EDS for lifecycle tests.
	results = append(results, r.RunTest("cloudtrail", "CreateEventDataStore_Lifecycle", func() error {
		resp, err := tc.client.CreateEventDataStore(tc.ctx, &cloudtrail.CreateEventDataStoreInput{
			Name:                         aws.String(fmt.Sprintf("ct-eds-lifecycle-%d", time.Now().UnixNano())),
			TerminationProtectionEnabled: aws.Bool(false),
			RetentionPeriod:              aws.Int32(90),
		})
		if err != nil {
			return fmt.Errorf("CreateEventDataStore failed: %w", err)
		}
		if resp.EventDataStoreArn == nil {
			return fmt.Errorf("EventDataStoreArn is nil")
		}
		if idx := strings.LastIndex(*resp.EventDataStoreArn, "/"); idx >= 0 {
			edsID = (*resp.EventDataStoreArn)[idx+1:]
		}
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

	// Stop ingestion.
	results = append(results, r.RunTest("cloudtrail", "StopEventDataStoreIngestion", func() error {
		_, err := tc.client.StopEventDataStoreIngestion(tc.ctx, &cloudtrail.StopEventDataStoreIngestionInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return fmt.Errorf("StopEventDataStoreIngestion failed: %w", err)
		}
		return nil
	}))

	// Start ingestion.
	results = append(results, r.RunTest("cloudtrail", "StartEventDataStoreIngestion", func() error {
		_, err := tc.client.StartEventDataStoreIngestion(tc.ctx, &cloudtrail.StartEventDataStoreIngestionInput{
			EventDataStore: aws.String(edsID),
		})
		if err != nil {
			return fmt.Errorf("StartEventDataStoreIngestion failed: %w", err)
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

	// Clean up: delete the restored EDS.
	results = append(results, r.RunTest("cloudtrail", "DeleteEDS_Cleanup", func() error {
		_, err := tc.client.DeleteEventDataStore(tc.ctx, &cloudtrail.DeleteEventDataStoreInput{
			EventDataStore: aws.String(edsID),
		})
		return err
	}))

	return results
}

func (r *TestRunner) runCloudTrailChannelTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	var channelARN string

	// Create channel.
	results = append(results, r.RunTest("cloudtrail", "CreateChannel_Success", func() error {
		resp, err := tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String("ct-channel-test"),
			Source: aws.String("Custom"),
			Destinations: []types.Destination{
				{
					Type:     types.DestinationTypeEventDataStore,
					Location: aws.String("test-eds-location"),
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

	// Get channel.
	results = append(results, r.RunTest("cloudtrail", "GetChannel_Success", func() error {
		resp, err := tc.client.GetChannel(tc.ctx, &cloudtrail.GetChannelInput{
			Channel: aws.String(channelARN),
		})
		if err != nil {
			return fmt.Errorf("GetChannel failed: %w", err)
		}
		if aws.ToString(resp.Name) != "ct-channel-test" {
			return fmt.Errorf("expected name ct-channel-test, got %s", aws.ToString(resp.Name))
		}
		return nil
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
			Name:    aws.String("ct-channel-updated"),
		})
		if err != nil {
			return fmt.Errorf("UpdateChannel failed: %w", err)
		}
		return nil
	}))

	// CreateChannel with tags (B1 regression: CreateChannel uses "Tags" not "TagsList").
	results = append(results, r.RunTest("cloudtrail", "CreateChannel_WithTags", func() error {
		tagCh := fmt.Sprintf("ct-channel-tags-%d", time.Now().UnixNano())
		resp, err := tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String(tagCh),
			Source: aws.String("Custom"),
			Destinations: []types.Destination{
				{
					Type:     types.DestinationTypeEventDataStore,
					Location: aws.String("test-eds-location"),
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
		tagMap := make(map[string]string)
		for _, t := range resp.Tags {
			tagMap[*t.Key] = *t.Value
		}
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
		resp, err := tc.client.CreateEventDataStore(tc.ctx, &cloudtrail.CreateEventDataStoreInput{
			Name:                         aws.String(fmt.Sprintf("ct-query-eds-%d", time.Now().UnixNano())),
			TerminationProtectionEnabled: aws.Bool(false),
		})
		if err != nil {
			return err
		}
		if idx := strings.LastIndex(*resp.EventDataStoreArn, "/"); idx >= 0 {
			edsID = (*resp.EventDataStoreArn)[idx+1:]
		}
		return nil
	}))

	// StartQuery.
	var queryID string
	results = append(results, r.RunTest("cloudtrail", "StartQuery_Success", func() error {
		stmt := fmt.Sprintf("SELECT eventID, eventTime FROM %s", edsID)
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

	// DescribeQuery — poll for FINISHED since StartQuery is asynchronous.
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
		return nil
	}))

	// GetQueryResults — query should be FINISHED after DescribeQuery poll.
	results = append(results, r.RunTest("cloudtrail", "GetQueryResults_Success", func() error {
		_, err := tc.client.GetQueryResults(tc.ctx, &cloudtrail.GetQueryResultsInput{
			EventDataStore: aws.String(edsID),
			QueryId:        aws.String(queryID),
		})
		if err != nil {
			return fmt.Errorf("GetQueryResults failed: %w", err)
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
				break
			}
		}
		if !found {
			return fmt.Errorf("created query not found in list")
		}
		return nil
	}))

	// Clean up EDS.
	results = append(results, r.RunTest("cloudtrail", "Query_CleanupEDS", func() error {
		_, err := tc.client.DeleteEventDataStore(tc.ctx, &cloudtrail.DeleteEventDataStoreInput{
			EventDataStore: aws.String(edsID),
		})
		return err
	}))

	return results
}

func (r *TestRunner) runCloudTrailConfigTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	// PutEventConfiguration for a trail.
	results = append(results, r.RunTest("cloudtrail", "PutEventConfiguration_Trail", func() error {
		_, err := tc.client.PutEventConfiguration(tc.ctx, &cloudtrail.PutEventConfigurationInput{
			TrailName:    aws.String("ct-event-config-trail"),
			MaxEventSize: types.MaxEventSizeStandard,
		})
		if err != nil {
			return fmt.Errorf("PutEventConfiguration failed: %w", err)
		}
		return nil
	}))

	// GetEventConfiguration.
	results = append(results, r.RunTest("cloudtrail", "GetEventConfiguration_Trail", func() error {
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

	// RegisterOrganizationDelegatedAdmin.
	results = append(results, r.RunTest("cloudtrail", "RegisterDelegatedAdmin", func() error {
		_, err := tc.client.RegisterOrganizationDelegatedAdmin(tc.ctx, &cloudtrail.RegisterOrganizationDelegatedAdminInput{
			MemberAccountId: aws.String(tc.accountID),
		})
		if err != nil {
			return fmt.Errorf("RegisterOrganizationDelegatedAdmin failed: %w", err)
		}
		return nil
	}))

	// DeregisterOrganizationDelegatedAdmin.
	results = append(results, r.RunTest("cloudtrail", "DeregisterDelegatedAdmin", func() error {
		_, err := tc.client.DeregisterOrganizationDelegatedAdmin(tc.ctx, &cloudtrail.DeregisterOrganizationDelegatedAdminInput{
			DelegatedAdminAccountId: aws.String(tc.accountID),
		})
		if err != nil {
			return fmt.Errorf("DeregisterOrganizationDelegatedAdmin failed: %w", err)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runCloudTrailBugFixTests(tc *cloudTrailTestContext) []TestResult {
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

	// UpdateTrail clears SnsTopicName by setting empty string.
	results = append(results, r.RunTest("cloudtrail", "UpdateTrail_ClearField", func() error {
		// Create trail with SnsTopicName.
		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String("ct-clear-field-test"),
			S3BucketName: aws.String("test-bucket"),
			SnsTopicName: aws.String("test-topic"),
		})
		if err != nil {
			return fmt.Errorf("CreateTrail failed: %w", err)
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

		// Cleanup.
		_, _ = tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{
			Name: aws.String("ct-clear-field-test"),
		})
		return nil
	}))

	// DeleteTrail by ARN.
	results = append(results, r.RunTest("cloudtrail", "DeleteTrail_ByARN", func() error {
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
