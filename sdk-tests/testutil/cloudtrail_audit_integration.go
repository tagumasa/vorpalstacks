package testutil

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"vorpalstacks-sdk-tests/config"
)

const (
	auditSvc       = "integration"
	auditPollWait  = 500 * time.Millisecond
	auditPollRetry = 6
)

func isAuditEnabled() bool {
	v := os.Getenv("VS_AUDIT_ENABLED")
	return v == "true" || v == "1"
}

type auditClients struct {
	cloudtrail *cloudtrail.Client
	s3         *s3.Client
	ctx        context.Context
}

func (r *TestRunner) newAuditClients() (*auditClients, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, err
	}
	return &auditClients{
		cloudtrail: cloudtrail.NewFromConfig(cfg),
		s3:         s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true }),
		ctx:        context.Background(),
	}, nil
}

func (r *TestRunner) RunCloudTrailAuditIntegrationTests() []TestResult {
	const skipReason = "VS_AUDIT_ENABLED=true is required for audit integration tests"
	var results []TestResult

	if !isAuditEnabled() {
		results = append(results, r.SkipTest(auditSvc, "CloudTrailAudit_CreateTrail_VerifyEvent", skipReason))
		results = append(results, r.SkipTest(auditSvc, "CloudTrailAudit_S3_PutObject", skipReason))
		results = append(results, r.SkipTest(auditSvc, "CloudTrailAudit_CrossService_EventSource", skipReason))
		return results
	}

	ac, err := r.newAuditClients()
	if err != nil {
		return append(results, SetupFailResult(auditSvc, "Failed to create clients: %v", err))
	}

	results = append(results, r.runAuditCreateTrailVerifyEvent(ac))
	results = append(results, r.runAuditS3PutObject(ac))
	results = append(results, r.runAuditCrossServiceEventSource(ac))

	return results
}

func (r *TestRunner) runAuditCreateTrailVerifyEvent(ac *auditClients) TestResult {
	return r.RunTest(auditSvc, "CloudTrailAudit_CreateTrail_VerifyEvent", func() error {
		name := fmt.Sprintf("audit-ct-%d", time.Now().UnixNano())
		defer ac.cloudtrail.DeleteTrail(ac.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := ac.cloudtrail.CreateTrail(ac.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("audit-ct-bucket"),
		})
		if err != nil {
			return fmt.Errorf("create trail: %v", err)
		}

		return r.pollAuditEvent("CreateTrail event", func() (*cloudtrail.LookupEventsOutput, error) {
			return ac.cloudtrail.LookupEvents(ac.ctx, &cloudtrail.LookupEventsInput{
				LookupAttributes: []types.LookupAttribute{
					{
						AttributeKey:   types.LookupAttributeKeyEventName,
						AttributeValue: aws.String("CreateTrail"),
					},
				},
				MaxResults: aws.Int32(50),
			})
		}, func(events []types.Event) error {
			for _, e := range events {
				if e.CloudTrailEvent != nil && *e.CloudTrailEvent != "" {
					return nil
				}
			}
			return fmt.Errorf("no CreateTrail event with CloudTrailEvent content found")
		})
	})
}

func (r *TestRunner) runAuditS3PutObject(ac *auditClients) TestResult {
	return r.RunTest(auditSvc, "CloudTrailAudit_S3_PutObject", func() error {
		bucket := fmt.Sprintf("audit-s3-%d", time.Now().UnixNano())
		key := "audit-test-key.txt"

		_, err := ac.s3.CreateBucket(ac.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)})
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer ac.s3.DeleteBucket(ac.ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})

		_, err = ac.s3.PutObject(ac.ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   nil,
		})
		if err != nil {
			return fmt.Errorf("put object: %v", err)
		}

		return r.pollAuditEvent("S3 PutObject event", func() (*cloudtrail.LookupEventsOutput, error) {
			return ac.cloudtrail.LookupEvents(ac.ctx, &cloudtrail.LookupEventsInput{
				LookupAttributes: []types.LookupAttribute{
					{
						AttributeKey:   types.LookupAttributeKeyEventSource,
						AttributeValue: aws.String("s3.amazonaws.com"),
					},
				},
				MaxResults: aws.Int32(50),
			})
		}, func(events []types.Event) error {
			for _, e := range events {
				if e.EventName != nil && *e.EventName == "PutObject" {
					return nil
				}
			}
			return fmt.Errorf("no PutObject event from s3.amazonaws.com found")
		})
	})
}

func (r *TestRunner) runAuditCrossServiceEventSource(ac *auditClients) TestResult {
	return r.RunTest(auditSvc, "CloudTrailAudit_CrossService_EventSource", func() error {
		cloudtrailEvents, err := ac.cloudtrail.LookupEvents(ac.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyEventSource,
					AttributeValue: aws.String("cloudtrail.amazonaws.com"),
				},
			},
			MaxResults: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("lookup cloudtrail events: %v", err)
		}
		if cloudtrailEvents.Events == nil || len(cloudtrailEvents.Events) == 0 {
			return fmt.Errorf("no cloudtrail.amazonaws.com events found")
		}

		s3Events, err := ac.cloudtrail.LookupEvents(ac.ctx, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyEventSource,
					AttributeValue: aws.String("s3.amazonaws.com"),
				},
			},
			MaxResults: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("lookup s3 events: %v", err)
		}
		if s3Events.Events == nil || len(s3Events.Events) == 0 {
			return fmt.Errorf("no s3.amazonaws.com events found")
		}

		return nil
	})
}

func (r *TestRunner) pollAuditEvent(label string, lookup func() (*cloudtrail.LookupEventsOutput, error), verify func([]types.Event) error) error {
	for i := 0; i < auditPollRetry; i++ {
		resp, err := lookup()
		if err != nil {
			return fmt.Errorf("lookup: %v", err)
		}
		if resp.Events != nil && len(resp.Events) > 0 {
			if vErr := verify(resp.Events); vErr == nil {
				return nil
			}
		}
		time.Sleep(auditPollWait)
	}
	return fmt.Errorf("timed out waiting for %s", label)
}
