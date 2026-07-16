package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
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
		for _, e := range resp.Events {
			if e.EventTime != nil {
				if e.EventTime.Before(pastHour) || e.EventTime.After(now) {
					return fmt.Errorf("event time %v outside query range [%v, %v]", *e.EventTime, pastHour, now)
				}
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_ByEventName", func() error {
		name := fmt.Sprintf("evtname-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("evtname-bucket"),
		})
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
		name := fmt.Sprintf("rnsearch-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("rnsearch-bucket"),
		})
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
		name := fmt.Sprintf("usersearch-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("usersearch-bucket"),
		})
		if err != nil {
			return fmt.Errorf("create trail: %v", err)
		}

		listResp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			MaxResults: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("list events: %v", err)
		}
		if listResp.Events == nil || len(listResp.Events) == 0 || listResp.Events[0].Username == nil {
			return fmt.Errorf("no events with Username available")
		}
		username := *listResp.Events[0].Username

		resp, err := tc.client.LookupEvents(tc.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyUsername,
					AttributeValue: &username,
				},
			},
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("lookup by Username: %v", err)
		}
		if resp.Events == nil || len(resp.Events) == 0 {
			return fmt.Errorf("expected at least 1 event for Username=%s", username)
		}
		for _, e := range resp.Events {
			if e.Username == nil || *e.Username != username {
				return fmt.Errorf("expected Username=%s, got %v", username, e.Username)
			}
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
