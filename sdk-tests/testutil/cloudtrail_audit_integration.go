package testutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"vorpalstacks-sdk-tests/config"
)

const (
	auditSvc       = "integration"
	auditPollWait  = 500 * time.Millisecond
	auditPollRetry = 6
)

func isAuditEnabled() bool {
	if v := os.Getenv("CLOUDTRAIL_ENABLED"); v == "true" || v == "1" {
		return true
	}
	if v := os.Getenv("ALL_SERVICES_ENABLED"); v == "true" || v == "1" {
		return true
	}
	return false
}

type auditClients struct {
	cloudtrail *cloudtrail.Client
	s3         *s3.Client
	ctx        context.Context
}

// ensureTrailBucket creates the bucket when missing and attaches the
// CloudTrail delivery policy — CreateTrail rejects a bucket without it.
func (ac *auditClients) ensureTrailBucket(accountID, bucket string) error {
	if _, err := ac.s3.CreateBucket(ac.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) || apiErr.ErrorCode() != "BucketAlreadyOwnedByYou" {
			return fmt.Errorf("create bucket %s: %w", bucket, err)
		}
	}
	return putCloudTrailBucketPolicy(ac.ctx, ac.s3, accountID, bucket)
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
	const skipReason = "CLOUDTRAIL_ENABLED=true is required for audit integration tests"
	var results []TestResult

	if !isAuditEnabled() {
		results = append(results, r.SkipTest(auditSvc, "CloudTrailAudit_CreateTrail_VerifyEvent", skipReason))
		results = append(results, r.SkipTest(auditSvc, "CloudTrailAudit_S3_PutObject", skipReason))
		results = append(results, r.SkipTest(auditSvc, "CloudTrailAudit_CrossService_EventSource", skipReason))
		results = append(results, r.SkipTest(auditSvc, "CloudTrailAudit_RecordContent", skipReason))
		return results
	}

	ac, err := r.newAuditClients()
	if err != nil {
		return append(results, SetupFailResult(auditSvc, "Failed to create clients: %v", err))
	}

	results = append(results, r.runAuditCreateTrailVerifyEvent(ac))
	results = append(results, r.runAuditS3PutObject(ac))
	results = append(results, r.runAuditCrossServiceEventSource(ac))
	results = append(results, r.runAuditRecordContent(ac))

	return results
}

func (r *TestRunner) runAuditCreateTrailVerifyEvent(ac *auditClients) TestResult {
	return r.RunTest(auditSvc, "CloudTrailAudit_CreateTrail_VerifyEvent", func() error {
		name := fmt.Sprintf("audit-ct-%d", time.Now().UnixNano())
		defer ac.cloudtrail.DeleteTrail(ac.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		if err := ac.ensureTrailBucket(r.accountID, "audit-ct-bucket"); err != nil {
			return err
		}
		_, err := ac.cloudtrail.CreateTrail(ac.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("audit-ct-bucket"),
		})
		if err != nil {
			return fmt.Errorf("create trail: %v", err)
		}

		return r.pollAuditEvent("CreateTrail event", func(nextToken *string) (*cloudtrail.LookupEventsOutput, error) {
			return ac.cloudtrail.LookupEvents(ac.ctx, &cloudtrail.LookupEventsInput{
				LookupAttributes: []types.LookupAttribute{
					{
						AttributeKey:   types.LookupAttributeKeyEventName,
						AttributeValue: aws.String("CreateTrail"),
					},
				},
				MaxResults: aws.Int32(50),
				NextToken:  nextToken,
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

		return r.pollAuditEvent("S3 PutObject event", func(nextToken *string) (*cloudtrail.LookupEventsOutput, error) {
			return ac.cloudtrail.LookupEvents(ac.ctx, &cloudtrail.LookupEventsInput{
				LookupAttributes: []types.LookupAttribute{
					{
						AttributeKey:   types.LookupAttributeKeyEventSource,
						AttributeValue: aws.String("s3.amazonaws.com"),
					},
				},
				MaxResults: aws.Int32(50),
				NextToken:  nextToken,
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

// runAuditRecordContent pins the per-operation record content: a read-only
// call records readOnly=true with the caller's real user agent, a write
// records readOnly=false, and a failed call records the model error code
// and message instead of a success-shaped record.
func (r *TestRunner) runAuditRecordContent(ac *auditClients) TestResult {
	return r.RunTest(auditSvc, "CloudTrailAudit_RecordContent", func() error {
		if _, err := ac.cloudtrail.GetTrail(ac.ctx, &cloudtrail.GetTrailInput{
			Name: aws.String(fmt.Sprintf("audit-missing-%d", time.Now().UnixNano())),
		}); err == nil {
			return fmt.Errorf("GetTrail on a missing trail must fail")
		}

		name := fmt.Sprintf("audit-ro-%d", time.Now().UnixNano())
		defer ac.cloudtrail.DeleteTrail(ac.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})
		if err := ac.ensureTrailBucket(r.accountID, "audit-ro-bucket"); err != nil {
			return err
		}
		if _, err := ac.cloudtrail.CreateTrail(ac.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("audit-ro-bucket"),
		}); err != nil {
			return fmt.Errorf("create trail: %v", err)
		}

		if _, err := ac.cloudtrail.ListTrails(ac.ctx, &cloudtrail.ListTrailsInput{}); err != nil {
			return fmt.Errorf("list trails: %v", err)
		}

		waitFor := func(eventName, wantReadOnly string, check func(types.Event) error) error {
			return r.pollAuditEvent(eventName, func(nextToken *string) (*cloudtrail.LookupEventsOutput, error) {
				return ac.cloudtrail.LookupEvents(ac.ctx, &cloudtrail.LookupEventsInput{
					LookupAttributes: []types.LookupAttribute{
						{
							AttributeKey:   types.LookupAttributeKeyEventName,
							AttributeValue: aws.String(eventName),
						},
					},
					MaxResults: aws.Int32(50),
					NextToken:  nextToken,
				})
			}, func(events []types.Event) error {
				for _, e := range events {
					if e.ReadOnly == nil || *e.ReadOnly != wantReadOnly {
						continue
					}
					if err := check(e); err != nil {
						continue
					}
					return nil
				}
				return fmt.Errorf("no %s event with readOnly=%s satisfies the content check", eventName, wantReadOnly)
			})
		}

		// The wire Event shape carries nine members with no userAgent; the
		// caller's user agent lives inside the CloudTrailEvent record JSON,
		// which must parse as an AWS CloudTrail record: eventID spelling,
		// boolean readOnly, awsRegion, and the derived userIdentity.
		if err := waitFor("ListTrails", "true", func(e types.Event) error {
			if e.CloudTrailEvent == nil || *e.CloudTrailEvent == "" {
				return fmt.Errorf("record has no CloudTrailEvent JSON")
			}
			var record map[string]interface{}
			if err := json.Unmarshal([]byte(*e.CloudTrailEvent), &record); err != nil {
				return fmt.Errorf("unmarshal CloudTrailEvent: %v", err)
			}
			ua, _ := record["userAgent"].(string)
			if ua == "" || ua == "vorpalstacks-internal" {
				return fmt.Errorf("userAgent = %q, want the caller's real user agent", ua)
			}
			if id, _ := record["eventID"].(string); id == "" || id != *e.EventId {
				return fmt.Errorf("record eventID = %v, want the wire EventId %s", record["eventID"], *e.EventId)
			}
			if ro, ok := record["readOnly"].(bool); !ok || !ro {
				return fmt.Errorf("record readOnly = %#v, want boolean true", record["readOnly"])
			}
			if region, _ := record["awsRegion"].(string); region == "" {
				return fmt.Errorf("record awsRegion = %v, want the request region", record["awsRegion"])
			}
			if category, _ := record["eventCategory"].(string); category != "Management" {
				return fmt.Errorf("record eventCategory = %v, want Management", record["eventCategory"])
			}
			if me, ok := record["managementEvent"].(bool); !ok || !me {
				return fmt.Errorf("record managementEvent = %#v, want boolean true", record["managementEvent"])
			}
			identity, ok := record["userIdentity"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("record userIdentity = %#v, want an object", record["userIdentity"])
			}
			if typ, _ := identity["type"].(string); typ == "" {
				return fmt.Errorf("record userIdentity.type is empty")
			}
			if _, ok := record["eventId"]; ok {
				return fmt.Errorf("record carries eventId, want the record spelling eventID")
			}
			if resp, present := record["responseElements"]; !present || resp != nil {
				return fmt.Errorf("record responseElements = %#v (present=%v), want present null for a read-only API", resp, present)
			}
			return nil
		}); err != nil {
			return err
		}

		// A write records a boolean readOnly false, its trail ARN as a
		// record-shaped resource, and non-null response elements.
		if err := waitFor("CreateTrail", "false", func(e types.Event) error {
			if e.CloudTrailEvent == nil || *e.CloudTrailEvent == "" {
				return fmt.Errorf("record has no CloudTrailEvent JSON")
			}
			var record map[string]interface{}
			if err := json.Unmarshal([]byte(*e.CloudTrailEvent), &record); err != nil {
				return fmt.Errorf("unmarshal CloudTrailEvent: %v", err)
			}
			if ro, ok := record["readOnly"].(bool); !ok || ro {
				return fmt.Errorf("record readOnly = %#v, want boolean false", record["readOnly"])
			}
			if resp, ok := record["responseElements"].(map[string]interface{}); !ok || len(resp) == 0 {
				return fmt.Errorf("record responseElements = %#v, want the write's response object", record["responseElements"])
			}
			resources, ok := record["resources"].([]interface{})
			if !ok || len(resources) == 0 {
				return fmt.Errorf("record resources = %#v, want the created trail", record["resources"])
			}
			for _, raw := range resources {
				entry, ok := raw.(map[string]interface{})
				if !ok {
					continue
				}
				arn, _ := entry["ARN"].(string)
				if strings.HasSuffix(arn, ":trail/"+name) {
					if typ, _ := entry["type"].(string); typ != "AWS::CloudTrail::Trail" {
						return fmt.Errorf("trail resource type = %v, want AWS::CloudTrail::Trail", entry["type"])
					}
					return nil
				}
			}
			return fmt.Errorf("record resources carry no entry for trail %s", name)
		}); err != nil {
			return err
		}

		return waitFor("GetTrail", "true", func(e types.Event) error {
			if e.CloudTrailEvent == nil || *e.CloudTrailEvent == "" {
				return fmt.Errorf("record has no CloudTrailEvent JSON")
			}
			var record map[string]interface{}
			if err := json.Unmarshal([]byte(*e.CloudTrailEvent), &record); err != nil {
				return fmt.Errorf("unmarshal CloudTrailEvent: %v", err)
			}
			if code, _ := record["errorCode"].(string); code != "TrailNotFoundException" {
				return fmt.Errorf("errorCode = %v, want TrailNotFoundException", record["errorCode"])
			}
			if msg, _ := record["errorMessage"].(string); msg == "" {
				return fmt.Errorf("errorMessage is empty")
			}
			return nil
		})
	})
}

func (r *TestRunner) pollAuditEvent(label string, lookup func(nextToken *string) (*cloudtrail.LookupEventsOutput, error), verify func([]types.Event) error) error {
	for i := 0; i < auditPollRetry; i++ {
		var token *string
		for {
			resp, err := lookup(token)
			if err != nil {
				return fmt.Errorf("lookup: %v", err)
			}
			if resp.Events != nil && len(resp.Events) > 0 {
				if vErr := verify(resp.Events); vErr == nil {
					return nil
				}
			}
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			token = resp.NextToken
		}
		time.Sleep(auditPollWait)
	}
	return fmt.Errorf("timed out waiting for %s", label)
}
