package testutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) runCloudTrailEventTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "LookupEvents", func() error {
		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.Events == nil {
			return fmt.Errorf("events list is nil")
		}
		for _, e := range resp.Events {
			if e.EventName == nil || *e.EventName == "" {
				return fmt.Errorf("event has empty EventName")
			}
			if e.EventTime == nil || e.EventTime.IsZero() {
				return fmt.Errorf("event has zero EventTime")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_WithTimeRange", func() error {
		now := time.Now().UTC()
		pastHour := now.Add(-1 * time.Hour)

		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			StartTime:  &pastHour,
			EndTime:    &now,
			MaxResults: aws.Int32(5),
		})
		if err != nil {
			return fmt.Errorf("lookup events with time range: %v", err)
		}
		if resp.Events == nil {
			return fmt.Errorf("events list is nil")
		}
		// Freshly recorded events must remain visible through the lookup
		// paths — every dispatched API call is recorded, so the past hour
		// always holds events, and the retention sweep keeps everything
		// younger than its bound.
		if len(resp.Events) == 0 {
			return fmt.Errorf("expected at least 1 event recorded within the past hour")
		}
		for _, e := range resp.Events {
			if e.EventTime != nil {
				if e.EventTime.Before(pastHour) || e.EventTime.After(now) {
					return fmt.Errorf("event time %v outside query range [%v, %v]", *e.EventTime, pastHour, now)
				}
			}
		}
		return nil
	}))

	// The events list is sorted by time, most recent first.
	results = append(results, r.RunTest("cloudtrail", "LookupEvents_DescendingOrder", func() error {
		now := time.Now().UTC()
		pastHour := now.Add(-1 * time.Hour)

		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			StartTime:  &pastHour,
			EndTime:    &now,
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("lookup events: %v", err)
		}
		if len(resp.Events) == 0 {
			return fmt.Errorf("expected at least 1 event recorded within the past hour")
		}
		for i := 1; i < len(resp.Events); i++ {
			prev, curr := resp.Events[i-1].EventTime, resp.Events[i].EventTime
			if prev != nil && curr != nil && curr.After(*prev) {
				return fmt.Errorf("events not sorted most-recent-first: %v at %d precedes %v",
					curr, i, prev)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_StartAfterEnd", func() error {
		now := time.Now().UTC()
		earlier := now.Add(-time.Minute)
		_, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			StartTime: &now,
			EndTime:   &earlier,
		})
		if err == nil {
			return fmt.Errorf("expected InvalidTimeRangeException for StartTime after EndTime")
		}
		var apiErr *types.InvalidTimeRangeException
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected InvalidTimeRangeException, got %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_MultipleAttributesRejected", func() error {
		_, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{AttributeKey: types.LookupAttributeKeyEventName, AttributeValue: aws.String("CreateTrail")},
				{AttributeKey: types.LookupAttributeKeyUsername, AttributeValue: aws.String("root")},
			},
		})
		if err == nil {
			return fmt.Errorf("expected InvalidLookupAttributesException for two attributes")
		}
		var apiErr *types.InvalidLookupAttributesException
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected InvalidLookupAttributesException, got %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_InvalidEventCategory", func() error {
		_, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			EventCategory: types.EventCategory("management"),
		})
		if err == nil {
			return fmt.Errorf("expected InvalidEventCategoryException for a non-insight category")
		}
		var apiErr *types.InvalidEventCategoryException
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected InvalidEventCategoryException, got %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_InvalidNextToken", func() error {
		_, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			NextToken: aws.String("not-a-token-this-service-issued"),
		})
		if err == nil {
			return fmt.Errorf("expected InvalidNextTokenException for an unusable token")
		}
		var apiErr *types.InvalidNextTokenException
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected InvalidNextTokenException, got %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_ByEventName", func() error {
		name := tc.uniqueName("evtname")
		defer tc.deleteTrail(name)

		_, err := tc.createTrail(name, "evtname-bucket")
		if err != nil {
			return fmt.Errorf("create trail: %v", err)
		}

		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyEventName,
					AttributeValue: aws.String("CreateTrail"),
				},
			},
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("lookup by EventName: %v", err)
		}
		if resp.Events == nil || len(resp.Events) == 0 {
			return fmt.Errorf("expected at least 1 event with EventName=CreateTrail")
		}
		for _, e := range resp.Events {
			if e.EventName == nil || *e.EventName != "CreateTrail" {
				return fmt.Errorf("expected EventName=CreateTrail, got %v", e.EventName)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_ByEventSource", func() error {
		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyEventSource,
					AttributeValue: aws.String("cloudtrail.amazonaws.com"),
				},
			},
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("lookup by EventSource: %v", err)
		}
		if resp.Events == nil || len(resp.Events) == 0 {
			return fmt.Errorf("expected at least 1 event with EventSource=cloudtrail.amazonaws.com")
		}
		for _, e := range resp.Events {
			if e.EventSource == nil || *e.EventSource != "cloudtrail.amazonaws.com" {
				return fmt.Errorf("expected EventSource=cloudtrail.amazonaws.com, got %v", e.EventSource)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_ByResourceName", func() error {
		name := tc.uniqueName("rnsearch")
		defer tc.deleteTrail(name)

		createResp, err := tc.createTrail(name, "rnsearch-bucket")
		if err != nil {
			return fmt.Errorf("create trail: %v", err)
		}
		if createResp.TrailARN == nil {
			return fmt.Errorf("trail ARN is nil")
		}

		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyResourceName,
					AttributeValue: createResp.TrailARN,
				},
			},
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("lookup by ResourceName: %v", err)
		}
		if resp.Events == nil || len(resp.Events) == 0 {
			return fmt.Errorf("expected at least 1 event for resource %s", *createResp.TrailARN)
		}
		for _, e := range resp.Events {
			if e.Resources == nil {
				continue
			}
			found := false
			for _, r := range e.Resources {
				if r.ResourceName != nil && *r.ResourceName == *createResp.TrailARN {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("event resources do not contain trail ARN %s", *createResp.TrailARN)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_ByUsername", func() error {
		// Username lookups resolve the authenticated principal's user
		// name, so the pin signs with a real IAM user's credentials: the
		// permissive plane records unresolvable callers as Unknown, with
		// no user name to find.
		cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   tc.region,
		})
		if err != nil {
			return fmt.Errorf("load config: %v", err)
		}
		iamClient := iam.NewFromConfig(cfg)

		userName := tc.uniqueName("usersearch")
		createResp, err := iamClient.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(userName)})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer iamClient.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(userName)})

		keyResp, err := iamClient.CreateAccessKey(tc.ctx, &iam.CreateAccessKeyInput{UserName: aws.String(userName)})
		if err != nil {
			return fmt.Errorf("create access key: %v", err)
		}
		defer iamClient.DeleteAccessKey(tc.ctx, &iam.DeleteAccessKeyInput{
			UserName:    aws.String(userName),
			AccessKeyId: keyResp.AccessKey.AccessKeyId,
		})

		userCfg := cfg
		userCfg.Credentials = credentials.NewStaticCredentialsProvider(
			*keyResp.AccessKey.AccessKeyId,
			*keyResp.AccessKey.SecretAccessKey,
			"",
		)
		userCloudTrail := cloudtrail.NewFromConfig(userCfg)

		trailName := tc.uniqueName("usersearch-trail")
		defer tc.deleteTrail(trailName)
		if err := tc.ensureTrailBucket("usersearch-bucket"); err != nil {
			return err
		}
		if _, err := userCloudTrail.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(trailName),
			S3BucketName: aws.String("usersearch-bucket"),
		}); err != nil {
			return fmt.Errorf("create trail as user: %v", err)
		}

		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyUsername,
					AttributeValue: aws.String(userName),
				},
			},
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("lookup by Username: %v", err)
		}
		if resp.Events == nil || len(resp.Events) == 0 {
			return fmt.Errorf("expected at least 1 event for Username=%s", userName)
		}
		var createEvent *types.Event
		for i := range resp.Events {
			if resp.Events[i].EventName != nil && *resp.Events[i].EventName == "CreateTrail" {
				createEvent = &resp.Events[i]
				break
			}
		}
		if createEvent == nil {
			return fmt.Errorf("no CreateTrail event recorded for Username=%s", userName)
		}
		for _, e := range resp.Events {
			if e.Username == nil || *e.Username != userName {
				return fmt.Errorf("expected Username=%s, got %v", userName, e.Username)
			}
		}
		if createEvent.CloudTrailEvent == nil || *createEvent.CloudTrailEvent == "" {
			return fmt.Errorf("CreateTrail event has no record JSON")
		}
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(*createEvent.CloudTrailEvent), &record); err != nil {
			return fmt.Errorf("unmarshal record: %v", err)
		}
		identity, _ := record["userIdentity"].(map[string]interface{})
		if identity == nil {
			return fmt.Errorf("record has no userIdentity object")
		}
		if typ, _ := identity["type"].(string); typ != "IAMUser" {
			return fmt.Errorf("record userIdentity.type = %v, want IAMUser", identity["type"])
		}
		if un, _ := identity["userName"].(string); un != userName {
			return fmt.Errorf("record userIdentity.userName = %v, want %s", identity["userName"], userName)
		}
		if arn, _ := identity["arn"].(string); createResp.User.Arn != nil && arn != *createResp.User.Arn {
			return fmt.Errorf("record userIdentity.arn = %v, want %s", identity["arn"], *createResp.User.Arn)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_ByEventId", func() error {
		listResp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			MaxResults: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("list events: %v", err)
		}
		if listResp.Events == nil || len(listResp.Events) == 0 {
			return fmt.Errorf("no events available for EventId lookup")
		}
		eventID := listResp.Events[0].EventId
		if eventID == nil {
			return fmt.Errorf("event has nil EventId")
		}

		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyEventId,
					AttributeValue: eventID,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("lookup by EventId: %v", err)
		}
		if resp.Events == nil || len(resp.Events) != 1 {
			return fmt.Errorf("expected exactly 1 event for EventId=%s, got %d", *eventID, len(resp.Events))
		}
		if resp.Events[0].EventId == nil || *resp.Events[0].EventId != *eventID {
			return fmt.Errorf("expected EventId=%s, got %v", *eventID, resp.Events[0].EventId)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_ByReadOnly", func() error {
		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyReadOnly,
					AttributeValue: aws.String("true"),
				},
			},
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("lookup by ReadOnly=true: %v", err)
		}
		if resp.Events == nil {
			return fmt.Errorf("events list is nil")
		}
		for _, e := range resp.Events {
			if e.ReadOnly == nil || *e.ReadOnly != "true" {
				return fmt.Errorf("expected ReadOnly=true, got %v", e.ReadOnly)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_ByAccessKeyId", func() error {
		listResp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			MaxResults: aws.Int32(5),
		})
		if err != nil {
			return fmt.Errorf("list events: %v", err)
		}
		if listResp.Events == nil || len(listResp.Events) == 0 {
			return fmt.Errorf("no events available for AccessKeyId lookup")
		}
		var accessKeyID string
		for _, e := range listResp.Events {
			if e.AccessKeyId != nil && *e.AccessKeyId != "" {
				accessKeyID = *e.AccessKeyId
				break
			}
		}
		if accessKeyID == "" {
			return fmt.Errorf("no event with non-empty AccessKeyId found")
		}

		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyAccessKeyId,
					AttributeValue: &accessKeyID,
				},
			},
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("lookup by AccessKeyId: %v", err)
		}
		if resp.Events == nil || len(resp.Events) == 0 {
			return fmt.Errorf("expected at least 1 event for AccessKeyId=%s", accessKeyID)
		}
		for _, e := range resp.Events {
			if e.AccessKeyId == nil || *e.AccessKeyId != accessKeyID {
				return fmt.Errorf("expected AccessKeyId=%s, got %v", accessKeyID, e.AccessKeyId)
			}
		}
		return nil
	}))

	return results
}
