package testutil

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/smithy-go"
	"github.com/scritchley/orc"

	"vorpalstacks-sdk-tests/config"
)

// s3CreateReplicationDest creates a destination bucket for a replication
// test, enables versioning on it (a replication requirement), and returns
// its ARN together with a cleanup closure that empties and deletes it.
func s3CreateReplicationDest(ctx context.Context, client *s3.Client, name string) (string, func(), error) {
	if err := s3CreateBucket(ctx, client, name); err != nil {
		return "", nil, err
	}
	if err := s3EnableVersioning(ctx, client, name); err != nil {
		s3CleanupBucket(ctx, client, name)
		return "", nil, err
	}
	return "arn:aws:s3:::" + name, func() { s3CleanupBucket(ctx, client, name) }, nil
}

// waitForReplicationStatus polls HeadObject until the object's replication
// status reaches the wanted value, on a replication source or destination
// alike. On a destination the engine marks a replica REPLICA only after
// tags, ACL, and lock state have been applied, so waiting on the status
// also removes the race on those metadata assertions.
func waitForReplicationStatus(ctx context.Context, client *s3.Client, bucket, key string, want types.ReplicationStatus, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		head, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err == nil && head.ReplicationStatus == want {
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

func (r *TestRunner) s3BucketConfigTests(ctx context.Context, client *s3.Client, ts string, bucketName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("s3", "PutBucketTagging", func() error {
		_, err := client.PutBucketTagging(ctx, &s3.PutBucketTaggingInput{
			Bucket: aws.String(bucketName),
			Tagging: &types.Tagging{
				TagSet: []types.Tag{
					{Key: aws.String("Environment"), Value: aws.String("Test")},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketTagging failed: %w", err)
		}
		getResp, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketTagging failed: %w", err)
		}
		found := false
		for _, tag := range getResp.TagSet {
			if tag.Key != nil && *tag.Key == "Environment" && tag.Value != nil && *tag.Value == "Test" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("tag Environment=Test not found in TagSet")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketTagging", func() error {
		_, err := client.DeleteBucketTagging(ctx, &s3.DeleteBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("DeleteBucketTagging failed: %w", err)
		}
		getResp, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			if !strings.Contains(err.Error(), "NoSuchTagSet") {
				return fmt.Errorf("unexpected error after deleting tags: %v", err)
			}
			return nil
		}
		if len(getResp.TagSet) != 0 {
			return fmt.Errorf("expected empty TagSet after deletion, got %d tags", len(getResp.TagSet))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketAcl", func() error {
		_, err := client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
			Bucket: aws.String(bucketName),
			ACL:    types.BucketCannedACLPrivate,
		})
		if err != nil {
			return fmt.Errorf("PutBucketAcl failed: %w", err)
		}
		getResp, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketAcl failed: %w", err)
		}
		if getResp.Owner == nil {
			return fmt.Errorf("Owner is nil")
		}
		if getResp.Owner.DisplayName == nil || *getResp.Owner.DisplayName == "" {
			return fmt.Errorf("Owner.DisplayName is nil or empty")
		}
		if getResp.Owner.ID == nil || *getResp.Owner.ID == "" {
			return fmt.Errorf("Owner.ID is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetBucketAcl", func() error {
		resp, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketAcl failed: %w", err)
		}
		if resp.Owner == nil {
			return fmt.Errorf("Owner is nil")
		}
		if resp.Owner.DisplayName == nil || *resp.Owner.DisplayName == "" {
			return fmt.Errorf("Owner.DisplayName is nil or empty")
		}
		if resp.Owner.ID == nil || *resp.Owner.ID == "" {
			return fmt.Errorf("Owner.ID is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketPolicy", func() error {
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:GetObject","Resource":"arn:aws:s3:::` + bucketName + `/*"}]}`
		_, err := client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
			Bucket: aws.String(bucketName),
			Policy: aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("PutBucketPolicy failed: %w", err)
		}
		getResp, err := client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketPolicy failed: %w", err)
		}
		if getResp.Policy == nil {
			return fmt.Errorf("Policy is nil")
		}
		if !strings.Contains(*getResp.Policy, "Allow") {
			return fmt.Errorf("policy does not contain 'Allow'")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetBucketPolicy_PlainJSON", func() error {
		principalARN := fmt.Sprintf("arn:aws:iam::%s:root", r.accountID)
		policyDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"` + principalARN + `"},"Action":"s3:*","Resource":"arn:aws:s3:::` + bucketName + `/*"}]}`
		_, err := client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
			Bucket: aws.String(bucketName),
			Policy: aws.String(policyDoc),
		})
		if err != nil {
			return fmt.Errorf("PutBucketPolicy failed: %w", err)
		}
		getResp, err := client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketPolicy failed: %w", err)
		}
		if getResp.Policy == nil {
			return fmt.Errorf("Policy is nil")
		}
		policyStr := *getResp.Policy
		if strings.Contains(policyStr, "&quot;") || strings.Contains(policyStr, "&amp;") || strings.Contains(policyStr, "&lt;") {
			return fmt.Errorf("policy contains XML-escaped characters: %s", policyStr)
		}
		if policyStr[0] != '{' {
			return fmt.Errorf("policy does not start with '{': %s", policyStr[:min(50, len(policyStr))])
		}
		if !strings.Contains(policyStr, "\"Version\"") || !strings.Contains(policyStr, "\"Statement\"") {
			return fmt.Errorf("policy is not valid JSON: %s", policyStr[:min(100, len(policyStr))])
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketPolicy", func() error {
		_, err := client.DeleteBucketPolicy(ctx, &s3.DeleteBucketPolicyInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("DeleteBucketPolicy failed: %w", err)
		}
		_, err = client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting policy, got nil")
		}
		if !strings.Contains(err.Error(), "NoSuchBucketPolicy") {
			return fmt.Errorf("expected NoSuchBucketPolicy error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketPolicy_InvalidJSON", func() error {
		_, err := client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
			Bucket: aws.String(bucketName),
			Policy: aws.String("{invalid json"),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid JSON policy, got nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketCors", func() error {
		_, err := client.PutBucketCors(ctx, &s3.PutBucketCorsInput{
			Bucket: aws.String(bucketName),
			CORSConfiguration: &types.CORSConfiguration{
				CORSRules: []types.CORSRule{
					{
						AllowedMethods: []string{"GET", "PUT"},
						AllowedOrigins: []string{"https://example.com"},
						MaxAgeSeconds:  aws.Int32(3600),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketCors failed: %w", err)
		}
		getResp, err := client.GetBucketCors(ctx, &s3.GetBucketCorsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketCors failed: %w", err)
		}
		if len(getResp.CORSRules) == 0 {
			return fmt.Errorf("CORSRules is empty")
		}
		rule := getResp.CORSRules[0]
		if rule.MaxAgeSeconds == nil || *rule.MaxAgeSeconds != 3600 {
			return fmt.Errorf("expected MaxAgeSeconds=3600, got %v", rule.MaxAgeSeconds)
		}
		foundGet := false
		for _, m := range rule.AllowedMethods {
			if m == "GET" {
				foundGet = true
				break
			}
		}
		if !foundGet {
			return fmt.Errorf("AllowedMethods does not contain GET: %v", rule.AllowedMethods)
		}
		foundOrigin := false
		for _, o := range rule.AllowedOrigins {
			if o == "https://example.com" {
				foundOrigin = true
				break
			}
		}
		if !foundOrigin {
			return fmt.Errorf("AllowedOrigins does not contain https://example.com: %v", rule.AllowedOrigins)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketCors", func() error {
		_, err := client.DeleteBucketCors(ctx, &s3.DeleteBucketCorsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("DeleteBucketCors failed: %w", err)
		}
		_, err = client.GetBucketCors(ctx, &s3.GetBucketCorsInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting CORS, got nil")
		}
		if !strings.Contains(err.Error(), "NoSuchCORSConfiguration") {
			return fmt.Errorf("expected NoSuchCORSConfiguration error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketEncryption", func() error {
		_, err := client.PutBucketEncryption(ctx, &s3.PutBucketEncryptionInput{
			Bucket: aws.String(bucketName),
			ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
				Rules: []types.ServerSideEncryptionRule{
					{
						ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
							SSEAlgorithm: types.ServerSideEncryptionAes256,
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketEncryption failed: %w", err)
		}
		getResp, err := client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketEncryption failed: %w", err)
		}
		if getResp.ServerSideEncryptionConfiguration == nil {
			return fmt.Errorf("ServerSideEncryptionConfiguration is nil")
		}
		rules := getResp.ServerSideEncryptionConfiguration.Rules
		if len(rules) == 0 {
			return fmt.Errorf("Rules is empty")
		}
		if rules[0].ApplyServerSideEncryptionByDefault == nil {
			return fmt.Errorf("ApplyServerSideEncryptionByDefault is nil")
		}
		if rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm != types.ServerSideEncryptionAes256 {
			return fmt.Errorf("expected SSEAlgorithm aes256, got %s", rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketEncryption", func() error {
		_, err := client.DeleteBucketEncryption(ctx, &s3.DeleteBucketEncryptionInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("DeleteBucketEncryption failed: %w", err)
		}
		_, err = client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting encryption, got nil")
		}
		if !strings.Contains(err.Error(), "ServerSideEncryptionConfigurationNotFoundError") {
			return fmt.Errorf("expected ServerSideEncryptionConfigurationNotFoundError, got: %v", err)
		}
		return nil
	}))

	// Both versioning status transitions round-trip: the requested status
	// comes back from GetBucketVersioning. Rows run in order, Enabled
	// first, then Suspended.
	results = append(results, r.RunTest("s3", "PutBucketVersioning_StatusRoundTrip", func() error {
		for _, c := range []struct {
			name   string
			status types.BucketVersioningStatus
		}{
			{name: "enable versioning", status: types.BucketVersioningStatusEnabled},
			{name: "suspend versioning", status: types.BucketVersioningStatusSuspended},
		} {
			_, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
				Bucket: aws.String(bucketName),
				VersioningConfiguration: &types.VersioningConfiguration{
					Status: c.status,
				},
			})
			if err != nil {
				return fmt.Errorf("%s: PutBucketVersioning failed: %w", c.name, err)
			}
			getResp, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
				Bucket: aws.String(bucketName),
			})
			if err != nil {
				return fmt.Errorf("%s: GetBucketVersioning failed: %w", c.name, err)
			}
			if getResp.Status != c.status {
				return fmt.Errorf("%s: expected Status %s, got %s", c.name, c.status, getResp.Status)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketLifecycleConfiguration", func() error {
		_, err := client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucketName),
			LifecycleConfiguration: &types.BucketLifecycleConfiguration{
				Rules: []types.LifecycleRule{
					{
						ID:     aws.String("test-expire-rule"),
						Status: types.ExpirationStatusEnabled,
						Filter: &types.LifecycleRuleFilter{
							Prefix: aws.String("logs/"),
						},
						Expiration: &types.LifecycleExpiration{
							Days: aws.Int32(30),
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketLifecycleConfiguration failed: %w", err)
		}
		getResp, err := client.GetBucketLifecycleConfiguration(ctx, &s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketLifecycleConfiguration failed: %w", err)
		}
		if len(getResp.Rules) == 0 {
			return fmt.Errorf("Rules is empty")
		}
		rule := getResp.Rules[0]
		if rule.ID == nil || *rule.ID != "test-expire-rule" {
			return fmt.Errorf("expected ID test-expire-rule, got %v", rule.ID)
		}
		if rule.Expiration == nil || rule.Expiration.Days == nil || *rule.Expiration.Days != 30 {
			return fmt.Errorf("expected Expiration.Days=30, got %v", rule.Expiration)
		}
		if rule.Filter == nil || rule.Filter.Prefix == nil || *rule.Filter.Prefix != "logs/" {
			return fmt.Errorf("expected Filter prefix logs/, got %v", rule.Filter)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketLifecycleConfiguration", func() error {
		_, err := client.DeleteBucketLifecycle(ctx, &s3.DeleteBucketLifecycleInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("DeleteBucketLifecycleConfiguration failed: %w", err)
		}
		_, err = client.GetBucketLifecycleConfiguration(ctx, &s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting lifecycle, got nil")
		}
		if !strings.Contains(err.Error(), "NoSuchLifecycleConfiguration") {
			return fmt.Errorf("expected NoSuchLifecycleConfiguration error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketWebsite", func() error {
		_, err := client.PutBucketWebsite(ctx, &s3.PutBucketWebsiteInput{
			Bucket: aws.String(bucketName),
			WebsiteConfiguration: &types.WebsiteConfiguration{
				IndexDocument: &types.IndexDocument{
					Suffix: aws.String("index.html"),
				},
				ErrorDocument: &types.ErrorDocument{
					Key: aws.String("error.html"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketWebsite failed: %w", err)
		}
		getResp, err := client.GetBucketWebsite(ctx, &s3.GetBucketWebsiteInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketWebsite failed: %w", err)
		}
		if getResp.IndexDocument == nil || getResp.IndexDocument.Suffix == nil || *getResp.IndexDocument.Suffix != "index.html" {
			return fmt.Errorf("IndexDocument.Suffix not index.html")
		}
		if getResp.ErrorDocument == nil || getResp.ErrorDocument.Key == nil || *getResp.ErrorDocument.Key != "error.html" {
			return fmt.Errorf("ErrorDocument.Key not error.html")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketWebsite", func() error {
		_, err := client.DeleteBucketWebsite(ctx, &s3.DeleteBucketWebsiteInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("DeleteBucketWebsite failed: %w", err)
		}
		_, err = client.GetBucketWebsite(ctx, &s3.GetBucketWebsiteInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting website, got nil")
		}
		if !strings.Contains(err.Error(), "NoSuchWebsiteConfiguration") {
			return fmt.Errorf("expected NoSuchWebsiteConfiguration error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketNotificationConfiguration", func() error {
		snsCfg, snsErr := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if snsErr != nil {
			return fmt.Errorf("failed to load SNS config: %w", snsErr)
		}
		snsClient := sns.NewFromConfig(snsCfg)
		topicName := fmt.Sprintf("test-topic-%s", ts)
		createTopicResp, err := snsClient.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String(topicName),
		})
		if err != nil {
			return fmt.Errorf("CreateTopic failed: %w", err)
		}
		topicArn := aws.ToString(createTopicResp.TopicArn)
		defer snsClient.DeleteTopic(ctx, &sns.DeleteTopicInput{
			TopicArn: createTopicResp.TopicArn,
		})
		_, err = client.PutBucketNotificationConfiguration(ctx, &s3.PutBucketNotificationConfigurationInput{
			Bucket: aws.String(bucketName),
			NotificationConfiguration: &types.NotificationConfiguration{
				TopicConfigurations: []types.TopicConfiguration{
					{
						TopicArn: aws.String(topicArn),
						Events: []types.Event{
							types.EventS3ObjectCreatedPut,
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketNotificationConfiguration failed: %w", err)
		}
		getResp, err := client.GetBucketNotificationConfiguration(ctx, &s3.GetBucketNotificationConfigurationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketNotificationConfiguration failed: %w", err)
		}
		if len(getResp.TopicConfigurations) == 0 {
			return fmt.Errorf("TopicConfigurations is empty")
		}
		tc := getResp.TopicConfigurations[0]
		if tc.TopicArn == nil || *tc.TopicArn != topicArn {
			return fmt.Errorf("expected TopicArn %s, got %v", topicArn, tc.TopicArn)
		}
		foundEvent := false
		for _, e := range tc.Events {
			if string(e) == "s3:ObjectCreated:Put" {
				foundEvent = true
				break
			}
		}
		if !foundEvent {
			return fmt.Errorf("Events do not contain s3:ObjectCreated:Put: %v", tc.Events)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketLogging", func() error {
		_, err := client.PutBucketLogging(ctx, &s3.PutBucketLoggingInput{
			Bucket: aws.String(bucketName),
			BucketLoggingStatus: &types.BucketLoggingStatus{
				LoggingEnabled: &types.LoggingEnabled{
					TargetBucket: aws.String(bucketName),
					TargetPrefix: aws.String("logs/"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketLogging failed: %w", err)
		}
		getResp, err := client.GetBucketLogging(ctx, &s3.GetBucketLoggingInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketLogging failed: %w", err)
		}
		if getResp.LoggingEnabled == nil {
			return fmt.Errorf("LoggingEnabled is nil")
		}
		if getResp.LoggingEnabled.TargetPrefix == nil || *getResp.LoggingEnabled.TargetPrefix != "logs/" {
			return fmt.Errorf("expected TargetPrefix logs/, got %v", getResp.LoggingEnabled.TargetPrefix)
		}
		if getResp.LoggingEnabled.TargetBucket == nil || *getResp.LoggingEnabled.TargetBucket != bucketName {
			return fmt.Errorf("expected TargetBucket %s, got %v", bucketName, getResp.LoggingEnabled.TargetBucket)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutPublicAccessBlock", func() error {
		_, err := client.PutPublicAccessBlock(ctx, &s3.PutPublicAccessBlockInput{
			Bucket: aws.String(bucketName),
			PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
				BlockPublicAcls:       aws.Bool(true),
				BlockPublicPolicy:     aws.Bool(true),
				IgnorePublicAcls:      aws.Bool(true),
				RestrictPublicBuckets: aws.Bool(true),
			},
		})
		if err != nil {
			return fmt.Errorf("PutPublicAccessBlock failed: %w", err)
		}
		getResp, err := client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetPublicAccessBlock failed: %w", err)
		}
		cfg := getResp.PublicAccessBlockConfiguration
		if cfg == nil {
			return fmt.Errorf("PublicAccessBlockConfiguration is nil")
		}
		if cfg.BlockPublicAcls == nil || !*cfg.BlockPublicAcls {
			return fmt.Errorf("BlockPublicAcls not true")
		}
		if cfg.BlockPublicPolicy == nil || !*cfg.BlockPublicPolicy {
			return fmt.Errorf("BlockPublicPolicy not true")
		}
		if cfg.IgnorePublicAcls == nil || !*cfg.IgnorePublicAcls {
			return fmt.Errorf("IgnorePublicAcls not true")
		}
		if cfg.RestrictPublicBuckets == nil || !*cfg.RestrictPublicBuckets {
			return fmt.Errorf("RestrictPublicBuckets not true")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeletePublicAccessBlock", func() error {
		_, err := client.DeletePublicAccessBlock(ctx, &s3.DeletePublicAccessBlockInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("DeletePublicAccessBlock failed: %w", err)
		}
		_, err = client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting public access block, got nil")
		}
		if !strings.Contains(err.Error(), "NoSuchPublicAccessBlockConfiguration") {
			return fmt.Errorf("expected NoSuchPublicAccessBlockConfiguration error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketOwnershipControls", func() error {
		_, err := client.PutBucketOwnershipControls(ctx, &s3.PutBucketOwnershipControlsInput{
			Bucket: aws.String(bucketName),
			OwnershipControls: &types.OwnershipControls{
				Rules: []types.OwnershipControlsRule{
					{
						ObjectOwnership: types.ObjectOwnershipBucketOwnerPreferred,
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketOwnershipControls failed: %w", err)
		}
		getResp, err := client.GetBucketOwnershipControls(ctx, &s3.GetBucketOwnershipControlsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketOwnershipControls failed: %w", err)
		}
		if len(getResp.OwnershipControls.Rules) == 0 {
			return fmt.Errorf("Rules is empty")
		}
		if getResp.OwnershipControls.Rules[0].ObjectOwnership != types.ObjectOwnershipBucketOwnerPreferred {
			return fmt.Errorf("expected ObjectOwnership BucketOwnerPreferred, got %s", getResp.OwnershipControls.Rules[0].ObjectOwnership)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketOwnershipControls", func() error {
		_, err := client.DeleteBucketOwnershipControls(ctx, &s3.DeleteBucketOwnershipControlsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("DeleteBucketOwnershipControls failed: %w", err)
		}
		_, err = client.GetBucketOwnershipControls(ctx, &s3.GetBucketOwnershipControlsInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting ownership controls, got nil")
		}
		if !strings.Contains(err.Error(), "OwnershipControlsNotFoundError") {
			return fmt.Errorf("expected OwnershipControlsNotFoundError, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketRequestPayment", func() error {
		_, err := client.PutBucketRequestPayment(ctx, &s3.PutBucketRequestPaymentInput{
			Bucket: aws.String(bucketName),
			RequestPaymentConfiguration: &types.RequestPaymentConfiguration{
				Payer: types.PayerRequester,
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketRequestPayment failed: %w", err)
		}
		getResp, err := client.GetBucketRequestPayment(ctx, &s3.GetBucketRequestPaymentInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketRequestPayment failed: %w", err)
		}
		if getResp.Payer != types.PayerRequester {
			return fmt.Errorf("expected Payer Requester, got %s", getResp.Payer)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketAccelerateConfiguration", func() error {
		_, err := client.PutBucketAccelerateConfiguration(ctx, &s3.PutBucketAccelerateConfigurationInput{
			Bucket: aws.String(bucketName),
			AccelerateConfiguration: &types.AccelerateConfiguration{
				Status: types.BucketAccelerateStatusSuspended,
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketAccelerateConfiguration failed: %w", err)
		}
		getResp, err := client.GetBucketAccelerateConfiguration(ctx, &s3.GetBucketAccelerateConfigurationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketAccelerateConfiguration failed: %w", err)
		}
		if getResp.Status != types.BucketAccelerateStatusSuspended {
			return fmt.Errorf("expected Status Suspended, got %s", getResp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketReplication_GetVerify", func() error {
		repArn, repCleanup, err := s3CreateReplicationDest(ctx, client, s3Bucket(ts, "repl-dest"))
		if err != nil {
			return err
		}
		defer repCleanup()

		// Replication requires versioning on the source bucket.
		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		if err := r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("rule-1", repArn, "")); err != nil {
			return err
		}

		getResp, err := client.GetBucketReplication(ctx, &s3.GetBucketReplicationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketReplication failed: %w", err)
		}
		if getResp.ReplicationConfiguration == nil {
			return fmt.Errorf("ReplicationConfiguration is nil")
		}
		if len(getResp.ReplicationConfiguration.Rules) != 1 {
			return fmt.Errorf("expected 1 rule, got %d", len(getResp.ReplicationConfiguration.Rules))
		}
		rule := getResp.ReplicationConfiguration.Rules[0]
		if rule.Status != types.ReplicationRuleStatusEnabled {
			return fmt.Errorf("expected Status Enabled, got %s", rule.Status)
		}
		if rule.Destination == nil || rule.Destination.Bucket == nil {
			return fmt.Errorf("Destination or Destination.Bucket is nil")
		}
		if *rule.Destination.Bucket != repArn {
			return fmt.Errorf("expected Destination.Bucket %s, got %s", repArn, *rule.Destination.Bucket)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "Replication_ObjectCopiedToDest", func() error {
		repBucket := s3Bucket(ts, "repl-copy-dest")
		repArn, repCleanup, err := s3CreateReplicationDest(ctx, client, repBucket)
		if err != nil {
			return err
		}
		defer repCleanup()

		// Replication requires versioning on the source bucket.
		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		if err := r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("copy-rule", repArn, "")); err != nil {
			return err
		}

		repKey := "repl-test-object.txt"
		if _, err := s3PutObject(ctx, client, bucketName, repKey, "replicated content"); err != nil {
			return err
		}

		if !waitForObject(ctx, client, repBucket, repKey, 5*time.Second) {
			return fmt.Errorf("GetObject from destination bucket failed (replication did not occur)")
		}

		getResp, _, err := s3GetRead(ctx, client, repBucket, repKey)
		if err != nil {
			return fmt.Errorf("GetObject from destination bucket failed (replication did not occur): %w", err)
		}
		if getResp.ContentLength == nil || *getResp.ContentLength != int64(len("replicated content")) {
			return fmt.Errorf("replicated object size mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "Replication_TagFilter", func() error {
		tagBucket := s3Bucket(ts, "repl-tag-dest")
		tagArn, tagCleanup, err := s3CreateReplicationDest(ctx, client, tagBucket)
		if err != nil {
			return err
		}
		defer tagCleanup()

		// Replication requires versioning on the source bucket.
		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		// The tag-filter rule keeps its literal rule: s3ReplicationRule
		// covers only the prefix-filter form.
		if err := r.s3PutReplication(ctx, client, bucketName, types.ReplicationRule{
			ID:       aws.String("tag-rule"),
			Status:   types.ReplicationRuleStatusEnabled,
			Priority: aws.Int32(1),
			Filter: &types.ReplicationRuleFilter{
				Tag: &types.Tag{Key: aws.String("replicate"), Value: aws.String("true")},
			},
			Destination: &types.Destination{
				Bucket: aws.String(tagArn),
			},
			DeleteMarkerReplication: &types.DeleteMarkerReplication{
				Status: types.DeleteMarkerReplicationStatusDisabled,
			},
		}); err != nil {
			return err
		}

		// Object WITH the matching tag should replicate.
		matchedKey := "tag-matched.txt"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:  aws.String(bucketName),
			Key:     aws.String(matchedKey),
			Body:    strings.NewReader("tagged content"),
			Tagging: aws.String("replicate=true"),
		}); err != nil {
			return fmt.Errorf("PutObject (tagged) failed: %w", err)
		}

		// Object WITHOUT the matching tag should NOT replicate.
		unmatchedKey := "tag-unmatched.txt"
		if _, err := s3PutObject(ctx, client, bucketName, unmatchedKey, "untagged content"); err != nil {
			return err
		}

		if !waitForObject(ctx, client, tagBucket, matchedKey, 5*time.Second) {
			return fmt.Errorf("matched object was not replicated to destination")
		}

		// Unmatched object should NOT exist in destination.
		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(tagBucket),
			Key:    aws.String(unmatchedKey),
		})
		if err == nil {
			return fmt.Errorf("unmatched object should not have been replicated")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "Replication_DeleteMarkerPropagated", func() error {
		// Re-enable versioning on the source bucket (a prior test suspends it).
		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		dmBucket := s3Bucket(ts, "repl-dm-dest")
		dmArn, dmCleanup, err := s3CreateReplicationDest(ctx, client, dmBucket)
		if err != nil {
			return err
		}
		defer dmCleanup()

		// Enable versioning on destination so delete markers are supported.
		if err := s3EnableVersioning(ctx, client, dmBucket); err != nil {
			return err
		}

		// The delete-marker-enabled rule keeps its literal rule:
		// s3ReplicationRule covers only the delete-marker-disabled form.
		if err := r.s3PutReplication(ctx, client, bucketName, types.ReplicationRule{
			ID:       aws.String("dm-rule"),
			Status:   types.ReplicationRuleStatusEnabled,
			Priority: aws.Int32(1),
			Filter:   &types.ReplicationRuleFilter{Prefix: aws.String("")},
			Destination: &types.Destination{
				Bucket: aws.String(dmArn),
			},
			DeleteMarkerReplication: &types.DeleteMarkerReplication{
				Status: types.DeleteMarkerReplicationStatusEnabled,
			},
		}); err != nil {
			return err
		}

		dmKey := "dm-test-object.txt"
		if _, err := s3PutObject(ctx, client, bucketName, dmKey, "dm content"); err != nil {
			return err
		}

		if !waitForObject(ctx, client, dmBucket, dmKey, 5*time.Second) {
			return fmt.Errorf("PutObject did not replicate to destination before delete marker test")
		}

		_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(dmKey),
		})
		if err != nil {
			return fmt.Errorf("DeleteObject failed: %w", err)
		}

		if !waitForObjectDeleted(ctx, client, dmBucket, dmKey, 5*time.Second) {
			return fmt.Errorf("delete marker was not propagated to destination")
		}

		s3CleanupBucket(ctx, client, dmBucket)

		// Versioning stays Enabled: suspending it on a bucket that still
		// carries the replication configuration is rejected, and the
		// configuration must remain for the deletion test that follows.
		// Every later replication test re-enables versioning explicitly.
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketReplication_VerifyGone", func() error {
		_, err := client.DeleteBucketReplication(ctx, &s3.DeleteBucketReplicationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("DeleteBucketReplication failed: %w", err)
		}
		_, err = client.GetBucketReplication(ctx, &s3.GetBucketReplicationInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error after DeleteBucketReplication, got nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketReplication_SourceVersioningRequired", func() error {
		srcBucket := s3Bucket(ts, "repl-nover-src")
		if err := s3CreateBucket(ctx, client, srcBucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, srcBucket)

		repArn, repCleanup, err := s3CreateReplicationDest(ctx, client, s3Bucket(ts, "repl-nover-dest"))
		if err != nil {
			return err
		}
		defer repCleanup()

		err = r.s3PutReplication(ctx, client, srcBucket, s3ReplicationRule("ver-rule", repArn, ""))
		return expectAWSErrorCode(err, "InvalidRequest")
	}))

	results = append(results, r.RunTest("s3", "PutBucketReplication_DestVersioningRequired", func() error {
		// The shared source bucket has versioning enabled by the earlier
		// replication tests in this suite.
		plainDest := s3Bucket(ts, "repl-plain-dest")
		if err := s3CreateBucket(ctx, client, plainDest); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, plainDest)

		err := r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("dest-ver-rule", "arn:aws:s3:::"+plainDest, ""))
		if err := expectAWSErrorCode(err, "InvalidRequest"); err != nil {
			return fmt.Errorf("unversioned destination should be rejected with InvalidRequest: %w", err)
		}

		// A destination bucket that does not exist is rejected the same way.
		err = r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("dest-ver-rule", "arn:aws:s3:::"+s3Bucket(ts, "repl-nonexistent"), ""))
		return expectAWSErrorCode(err, "InvalidRequest")
	}))

	results = append(results, r.RunTest("s3", "Replication_CompleteMultipartUploadReplicated", func() error {
		repBucket := s3Bucket(ts, "repl-mpu-dest")
		repArn, repCleanup, err := s3CreateReplicationDest(ctx, client, repBucket)
		if err != nil {
			return err
		}
		defer repCleanup()

		// Replication requires versioning on the source bucket.
		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		if err := r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("mpu-rule", repArn, "")); err != nil {
			return err
		}

		mpuKey := "repl-mpu-object.bin"
		createResp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(mpuKey),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload failed: %w", err)
		}
		uploadId := aws.ToString(createResp.UploadId)

		partSize := 5 * 1024 * 1024
		partData := bytes.Repeat([]byte("m"), partSize)
		var completedParts []types.CompletedPart
		for partNumber := int32(1); partNumber <= 2; partNumber++ {
			uploadResp, err := client.UploadPart(ctx, &s3.UploadPartInput{
				Bucket:     aws.String(bucketName),
				Key:        aws.String(mpuKey),
				UploadId:   aws.String(uploadId),
				PartNumber: aws.Int32(partNumber),
				Body:       bytes.NewReader(partData),
			})
			if err != nil {
				return fmt.Errorf("UploadPart %d failed: %w", partNumber, err)
			}
			completedParts = append(completedParts, types.CompletedPart{
				PartNumber: aws.Int32(partNumber),
				ETag:       uploadResp.ETag,
			})
		}
		if _, err := client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:          aws.String(bucketName),
			Key:             aws.String(mpuKey),
			UploadId:        aws.String(uploadId),
			MultipartUpload: &types.CompletedMultipartUpload{Parts: completedParts},
		}); err != nil {
			return fmt.Errorf("CompleteMultipartUpload failed: %w", err)
		}

		if !waitForReplicationStatus(ctx, client, repBucket, mpuKey, types.ReplicationStatusReplica, 5*time.Second) {
			return fmt.Errorf("replicated multipart object not marked REPLICA in destination bucket")
		}
		headResp, err := s3HeadObject(ctx, client, repBucket, mpuKey)
		if err != nil {
			return fmt.Errorf("HeadObject on replicated multipart object failed: %w", err)
		}
		if headResp.ContentLength == nil || *headResp.ContentLength != int64(2*partSize) {
			return fmt.Errorf("replicated multipart object size mismatch: got %v, want %d", headResp.ContentLength, 2*partSize)
		}
		if headResp.ReplicationStatus != types.ReplicationStatusReplica {
			return fmt.Errorf("expected replica ReplicationStatus REPLICA, got %s", headResp.ReplicationStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "Replication_CopyObjectReplicated", func() error {
		repBucket := s3Bucket(ts, "repl-cpo-dest")
		repArn, repCleanup, err := s3CreateReplicationDest(ctx, client, repBucket)
		if err != nil {
			return err
		}
		defer repCleanup()

		// Replication requires versioning on the source bucket.
		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		if err := r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("cpo-rule", repArn, "")); err != nil {
			return err
		}

		srcKey := "copy-source.txt"
		copyKey := "copied-via-copy.txt"
		body := "copy replication content"
		if _, err := s3PutObject(ctx, client, bucketName, srcKey, body); err != nil {
			return err
		}
		if _, err := client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(copyKey),
			CopySource: aws.String(bucketName + "/" + srcKey),
		}); err != nil {
			return fmt.Errorf("CopyObject failed: %w", err)
		}

		if !waitForObject(ctx, client, repBucket, copyKey, 5*time.Second) {
			return fmt.Errorf("copied object not found in destination bucket")
		}
		_, gotBody, err := s3GetRead(ctx, client, repBucket, copyKey)
		if err != nil {
			return fmt.Errorf("GetObject on copied replica failed: %w", err)
		}
		if gotBody != body {
			return fmt.Errorf("copied replica content mismatch: got %q, want %q", gotBody, body)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "Replication_ReplicaMetadata", func() error {
		repBucket := s3Bucket(ts, "repl-meta-dest")
		repArn, repCleanup, err := s3CreateReplicationDest(ctx, client, repBucket)
		if err != nil {
			return err
		}
		defer repCleanup()

		// Replication requires versioning on the source bucket.
		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		// The storage-class-override rule keeps its literal rule:
		// s3ReplicationRule covers only the plain destination form.
		if err := r.s3PutReplication(ctx, client, bucketName, types.ReplicationRule{
			ID:       aws.String("meta-rule"),
			Status:   types.ReplicationRuleStatusEnabled,
			Priority: aws.Int32(1),
			Filter:   &types.ReplicationRuleFilter{Prefix: aws.String("meta/")},
			Destination: &types.Destination{
				Bucket:       aws.String(repArn),
				StorageClass: types.StorageClassReducedRedundancy,
			},
			DeleteMarkerReplication: &types.DeleteMarkerReplication{
				Status: types.DeleteMarkerReplicationStatusDisabled,
			},
		}); err != nil {
			return err
		}

		// Override path: the rule's StorageClass wins and the source tags
		// are carried to the replica.
		overrideKey := "meta/override.txt"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:  aws.String(bucketName),
			Key:     aws.String(overrideKey),
			Body:    strings.NewReader("override content"),
			Tagging: aws.String("proj=replica&env=test"),
		}); err != nil {
			return fmt.Errorf("PutObject (override) failed: %w", err)
		}
		if !waitForReplicationStatus(ctx, client, repBucket, overrideKey, types.ReplicationStatusReplica, 5*time.Second) {
			return fmt.Errorf("override replica not marked REPLICA in destination bucket")
		}
		overrideHead, err := s3HeadObject(ctx, client, repBucket, overrideKey)
		if err != nil {
			return fmt.Errorf("HeadObject on override replica failed: %w", err)
		}
		if overrideHead.ReplicationStatus != types.ReplicationStatusReplica {
			return fmt.Errorf("expected replica ReplicationStatus REPLICA, got %s", overrideHead.ReplicationStatus)
		}
		if overrideHead.StorageClass != types.StorageClassReducedRedundancy {
			return fmt.Errorf("expected rule StorageClass REDUCED_REDUNDANCY on replica, got %s", overrideHead.StorageClass)
		}
		overrideTags, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(repBucket),
			Key:    aws.String(overrideKey),
		})
		if err != nil {
			return fmt.Errorf("GetObjectTagging on override replica failed: %w", err)
		}
		tagMap := map[string]string{}
		for _, t := range overrideTags.TagSet {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["proj"] != "replica" || tagMap["env"] != "test" {
			return fmt.Errorf("source tags not carried to replica: got %v", tagMap)
		}

		// Inherit path: without a rule StorageClass the replica keeps the
		// source object's own storage class. The rule is replaced with a
		// StorageClass-free one first — a rule's StorageClass applies to
		// every replica under it, so the override would win otherwise.
		if err := r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("meta-inherit-rule", repArn, "meta/")); err != nil {
			return err
		}

		inheritKey := "meta/inherit.txt"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String(inheritKey),
			Body:         strings.NewReader("inherit content"),
			StorageClass: types.StorageClassStandardIa,
		}); err != nil {
			return fmt.Errorf("PutObject (inherit) failed: %w", err)
		}
		if !waitForReplicationStatus(ctx, client, repBucket, inheritKey, types.ReplicationStatusReplica, 5*time.Second) {
			return fmt.Errorf("inherit replica not marked REPLICA in destination bucket")
		}
		inheritHead, err := s3HeadObject(ctx, client, repBucket, inheritKey)
		if err != nil {
			return fmt.Errorf("HeadObject on inherit replica failed: %w", err)
		}
		if inheritHead.StorageClass != types.StorageClassStandardIa {
			return fmt.Errorf("expected inherited StorageClass STANDARD_IA on replica, got %s", inheritHead.StorageClass)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "Replication_ReplicaNotReReplicated", func() error {
		// Chain A -> B -> C: the replica in B must not replicate onwards
		// to C just because B itself has a replication configuration.
		midBucket := s3Bucket(ts, "repl-chain-mid")
		midArn, midCleanup, err := s3CreateReplicationDest(ctx, client, midBucket)
		if err != nil {
			return err
		}
		defer midCleanup()

		finalBucket := s3Bucket(ts, "repl-chain-final")
		finalArn, finalCleanup, err := s3CreateReplicationDest(ctx, client, finalBucket)
		if err != nil {
			return err
		}
		defer finalCleanup()

		// Replication requires versioning on the source bucket.
		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		if err := r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("chain-rule", midArn, "chain/")); err != nil {
			return err
		}
		if err := r.s3PutReplication(ctx, client, midBucket, s3ReplicationRule("onwards-rule", finalArn, "")); err != nil {
			return err
		}

		chainKey := "chain/object.txt"
		if _, err := s3PutObject(ctx, client, bucketName, chainKey, "chain content"); err != nil {
			return err
		}
		if !waitForReplicationStatus(ctx, client, midBucket, chainKey, types.ReplicationStatusReplica, 5*time.Second) {
			return fmt.Errorf("replica not marked REPLICA in mid bucket")
		}

		midHead, err := s3HeadObject(ctx, client, midBucket, chainKey)
		if err != nil {
			return fmt.Errorf("HeadObject on mid replica failed: %w", err)
		}
		if midHead.ReplicationStatus != types.ReplicationStatusReplica {
			return fmt.Errorf("expected mid replica ReplicationStatus REPLICA, got %s", midHead.ReplicationStatus)
		}

		// The mid replica is marked REPLICA, so the onwards rule on the mid
		// bucket must leave it alone: after the replication pipeline has
		// settled (mid replica visible and marked), the object stays absent
		// from the final bucket.
		time.Sleep(2 * time.Second)
		_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(finalBucket),
			Key:    aws.String(chainKey),
		})
		if err == nil {
			return fmt.Errorf("replica must not be re-replicated to the final bucket")
		}
		return nil
	}))

	// An object eligible for replication carries the x-amz-replication-status
	// header from the moment the upload returns: PENDING while the async copy
	// is in flight, settling to COMPLETED afterwards. The transitional status
	// is transient by contract, so the immediate read accepts PENDING or
	// COMPLETED; the defect it pins is the header never existing between
	// upload and completion.
	results = append(results, r.RunTest("s3", "Replication_PendingStatusVisible", func() error {
		repBucket := s3Bucket(ts, "repl-pending-dest")
		repArn, repCleanup, err := s3CreateReplicationDest(ctx, client, repBucket)
		if err != nil {
			return err
		}
		defer repCleanup()

		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		if err := r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("pending-rule", repArn, "pending/")); err != nil {
			return err
		}

		// A large body keeps the async copy in flight long enough that the
		// transitional status is observable before it settles.
		pendingKey := "pending/large-object.bin"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(pendingKey),
			Body:   bytes.NewReader(bytes.Repeat([]byte("p"), 8*1024*1024)),
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		head, err := s3HeadObject(ctx, client, bucketName, pendingKey)
		if err != nil {
			return fmt.Errorf("HeadObject immediately after upload failed: %w", err)
		}
		if head.ReplicationStatus != types.ReplicationStatusPending && head.ReplicationStatus != types.ReplicationStatusCompleted {
			return fmt.Errorf("expected replication status PENDING or COMPLETED right after upload, got %q", head.ReplicationStatus)
		}

		if !waitForReplicationStatus(ctx, client, bucketName, pendingKey, types.ReplicationStatusCompleted, 10*time.Second) {
			return fmt.Errorf("source object never reached replication status COMPLETED")
		}
		return nil
	}))

	// With several destinations, the source status aggregates over every
	// matched rule: COMPLETED only when all destinations received the copy,
	// FAILED as soon as one destination is unreachable.
	results = append(results, r.RunTest("s3", "Replication_MultiDestinationStatusAggregation", func() error {
		dest1 := s3Bucket(ts, "repl-multi-dest1")
		dest1Arn, dest1Cleanup, err := s3CreateReplicationDest(ctx, client, dest1)
		if err != nil {
			return err
		}
		defer dest1Cleanup()

		// The second destination is managed manually: it is deleted
		// mid-test so one rule fails while the other still copies.
		dest2 := s3Bucket(ts, "repl-multi-dest2")
		if err := s3CreateBucket(ctx, client, dest2); err != nil {
			return err
		}
		if err := s3EnableVersioning(ctx, client, dest2); err != nil {
			return err
		}
		dest2Arn := "arn:aws:s3:::" + dest2

		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		if err := r.s3PutReplication(ctx, client, bucketName,
			s3ReplicationRule("multi-rule-1", dest1Arn, "multi/"),
			// The second rule keeps its literal form: rule priorities are
			// unique per configuration, so it cannot reuse the
			// priority-1 helper.
			types.ReplicationRule{
				ID:       aws.String("multi-rule-2"),
				Status:   types.ReplicationRuleStatusEnabled,
				Priority: aws.Int32(2),
				Filter:   &types.ReplicationRuleFilter{Prefix: aws.String("multi/")},
				Destination: &types.Destination{
					Bucket: aws.String(dest2Arn),
				},
				DeleteMarkerReplication: &types.DeleteMarkerReplication{
					Status: types.DeleteMarkerReplicationStatusDisabled,
				},
			},
		); err != nil {
			return err
		}

		// Remove the second destination so rule 2 fails while rule 1 still
		// receives the copy.
		if _, err := client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(dest2),
		}); err != nil {
			return fmt.Errorf("DeleteBucket (dest2) failed: %w", err)
		}

		aggKey := "multi/object.txt"
		if _, err := s3PutObject(ctx, client, bucketName, aggKey, "aggregated content"); err != nil {
			return err
		}

		if !waitForReplicationStatus(ctx, client, bucketName, aggKey, types.ReplicationStatusFailed, 5*time.Second) {
			return fmt.Errorf("expected aggregated replication status FAILED with one unreachable destination")
		}
		if !waitForObject(ctx, client, dest1, aggKey, 5*time.Second) {
			return fmt.Errorf("the reachable destination should still receive the replica")
		}

		// With both destinations reachable again the same configuration must
		// aggregate to COMPLETED.
		if err := s3CreateBucket(ctx, client, dest2); err != nil {
			return err
		}
		if err := s3EnableVersioning(ctx, client, dest2); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, dest2)

		aggKey2 := "multi/object2.txt"
		if _, err := s3PutObject(ctx, client, bucketName, aggKey2, "aggregated content 2"); err != nil {
			return err
		}
		if !waitForReplicationStatus(ctx, client, bucketName, aggKey2, types.ReplicationStatusCompleted, 5*time.Second) {
			return fmt.Errorf("expected aggregated replication status COMPLETED with both destinations reachable")
		}
		return nil
	}))

	// Suspending versioning on a destination bucket is accepted by the API,
	// but replication to it must then fail and leave the source FAILED
	// rather than writing a degraded replica.
	results = append(results, r.RunTest("s3", "Replication_DestVersioningSuspendedFails", func() error {
		repBucket := s3Bucket(ts, "repl-susp-dest")
		repArn, repCleanup, err := s3CreateReplicationDest(ctx, client, repBucket)
		if err != nil {
			return err
		}
		defer repCleanup()

		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		if err := r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("susp-rule", repArn, "susp/")); err != nil {
			return err
		}

		if _, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(repBucket),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusSuspended,
			},
		}); err != nil {
			return fmt.Errorf("PutBucketVersioning (dest suspend) failed: %w", err)
		}

		suspKey := "susp/object.txt"
		if _, err := s3PutObject(ctx, client, bucketName, suspKey, "suspended destination content"); err != nil {
			return err
		}

		if !waitForReplicationStatus(ctx, client, bucketName, suspKey, types.ReplicationStatusFailed, 5*time.Second) {
			return fmt.Errorf("expected replication status FAILED with versioning suspended on the destination")
		}
		if _, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(repBucket),
			Key:    aws.String(suspKey),
		}); err == nil {
			return fmt.Errorf("no replica may be written to a versioning-suspended destination")
		}
		return nil
	}))

	// Suspending versioning on a replication source is rejected until the
	// replication configuration is removed.
	results = append(results, r.RunTest("s3", "PutBucketVersioning_SuspendedWithReplicationRejected", func() error {
		repBucket := s3Bucket(ts, "repl-vers-dest")
		repArn, repCleanup, err := s3CreateReplicationDest(ctx, client, repBucket)
		if err != nil {
			return err
		}
		defer repCleanup()

		if err := s3EnableVersioning(ctx, client, bucketName); err != nil {
			return err
		}

		if err := r.s3PutReplication(ctx, client, bucketName, s3ReplicationRule("vers-rule", repArn, "vers/")); err != nil {
			return err
		}

		_, err = client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusSuspended,
			},
		})
		if err := expectAWSErrorCode(err, "InvalidRequest"); err != nil {
			return fmt.Errorf("suspending versioning with a replication configuration must be rejected: %w", err)
		}

		// Removing the configuration lifts the rejection.
		if _, err := client.DeleteBucketReplication(ctx, &s3.DeleteBucketReplicationInput{
			Bucket: aws.String(bucketName),
		}); err != nil {
			return fmt.Errorf("DeleteBucketReplication failed: %w", err)
		}
		if _, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusSuspended,
			},
		}); err != nil {
			return fmt.Errorf("suspending versioning after removing the replication configuration failed: %w", err)
		}

		// Restore versioning for the tests that follow.
		if _, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusEnabled,
			},
		}); err != nil {
			return fmt.Errorf("PutBucketVersioning (restore) failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketCORS_PreflightVerify", func() error {
		_, err := client.PutBucketCors(ctx, &s3.PutBucketCorsInput{
			Bucket: aws.String(bucketName),
			CORSConfiguration: &types.CORSConfiguration{
				CORSRules: []types.CORSRule{
					{
						AllowedOrigins: []string{"https://example.com"},
						AllowedMethods: []string{"GET", "PUT", "DELETE"},
						AllowedHeaders: []string{"x-amz-*", "Content-Type"},
						ExposeHeaders:  []string{"ETag", "x-amz-version-id"},
						MaxAgeSeconds:  aws.Int32(3600),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketCors failed: %w", err)
		}

		endpoint := r.endpoint
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = "http://" + endpoint
		}

		preflightReq, err := http.NewRequest("OPTIONS", endpoint+"/"+bucketName+"/cors-test.txt", nil)
		if err != nil {
			return fmt.Errorf("failed to create preflight request: %w", err)
		}
		preflightReq.Header.Set("Origin", "https://example.com")
		preflightReq.Header.Set("Access-Control-Request-Method", "PUT")
		preflightReq.Header.Set("Access-Control-Request-Headers", "x-amz-meta-author")

		resp, err := testHTTPClient.Do(preflightReq)
		if err != nil {
			return fmt.Errorf("preflight request failed: %w", err)
		}
		defer resp.Body.Close()

		allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
		if allowOrigin != "https://example.com" {
			return fmt.Errorf("expected Access-Control-Allow-Origin 'https://example.com', got '%s'", allowOrigin)
		}
		allowMethods := resp.Header.Get("Access-Control-Allow-Methods")
		if !strings.Contains(allowMethods, "PUT") {
			return fmt.Errorf("expected Access-Control-Allow-Methods to contain PUT, got '%s'", allowMethods)
		}
		allowHeaders := resp.Header.Get("Access-Control-Allow-Headers")
		if !strings.Contains(allowHeaders, "x-amz-*") {
			return fmt.Errorf("expected Access-Control-Allow-Headers to contain x-amz-*, got '%s'", allowHeaders)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CORS_ActualRequestHeaders", func() error {
		endpoint := r.endpoint
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = "http://" + endpoint
		}

		req, err := http.NewRequest("GET", endpoint+"/"+bucketName+"/cors-test.txt", nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Origin", "https://example.com")

		resp, err := testHTTPClient.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
		if allowOrigin != "https://example.com" {
			return fmt.Errorf("expected Access-Control-Allow-Origin 'https://example.com', got '%s'", allowOrigin)
		}
		exposeHeaders := resp.Header.Get("Access-Control-Expose-Headers")
		if !strings.Contains(exposeHeaders, "ETag") {
			return fmt.Errorf("expected Access-Control-Expose-Headers to contain ETag, got '%s'", exposeHeaders)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketLifecycle_VerifyConfig", func() error {
		_, err := client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucketName),
			LifecycleConfiguration: &types.BucketLifecycleConfiguration{
				Rules: []types.LifecycleRule{
					{
						ID:     aws.String("expire-old"),
						Status: types.ExpirationStatusEnabled,
						Filter: &types.LifecycleRuleFilter{
							Prefix: aws.String("expire-me/"),
						},
						Expiration: &types.LifecycleExpiration{
							Days: aws.Int32(1),
						},
					},
					{
						ID:     aws.String("abort-uploads"),
						Status: types.ExpirationStatusEnabled,
						Filter: &types.LifecycleRuleFilter{
							Prefix: aws.String(""),
						},
						AbortIncompleteMultipartUpload: &types.AbortIncompleteMultipartUpload{
							DaysAfterInitiation: aws.Int32(1),
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketLifecycleConfiguration failed: %w", err)
		}

		getResp, err := client.GetBucketLifecycleConfiguration(ctx, &s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketLifecycleConfiguration failed: %w", err)
		}
		if getResp.Rules == nil || len(getResp.Rules) != 2 {
			return fmt.Errorf("expected 2 lifecycle rules, got %d", len(getResp.Rules))
		}
		found := false
		for _, rule := range getResp.Rules {
			if rule.Expiration != nil && rule.Expiration.Days != nil && *rule.Expiration.Days == 1 {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected a rule with Expiration.Days=1")
		}
		return nil
	}))

	// GetBucketPolicyStatus classifies the bucket's effective policy: a
	// bucket without a policy is not public, a Principal-"*" policy is
	// public, and a policy bound to the account root principal is not.
	results = append(results, r.RunTest("s3", "GetBucketPolicyStatus", func() error {
		for _, c := range []struct {
			name       string
			suffix     string
			policy     func(bucket string) string // nil places no policy
			wantPublic bool
		}{
			{
				name:       "bucket without a policy",
				suffix:     "polstat-none",
				policy:     nil,
				wantPublic: false,
			},
			{
				name:   "public policy with wildcard principal",
				suffix: "polstat-public",
				policy: func(bucket string) string {
					return `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:GetObject","Resource":"arn:aws:s3:::` + bucket + `/*"}]}`
				},
				wantPublic: true,
			},
			{
				name:   "non-public policy with fixed principal",
				suffix: "polstat-fixed",
				policy: func(bucket string) string {
					principalARN := fmt.Sprintf("arn:aws:iam::%s:root", r.accountID)
					return `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"` + principalARN + `"},"Action":"s3:GetObject","Resource":"arn:aws:s3:::` + bucket + `/*"}]}`
				},
				wantPublic: false,
			},
		} {
			bucket := s3Bucket(ts, c.suffix)
			if err := s3CreateBucket(ctx, client, bucket); err != nil {
				return fmt.Errorf("%s: %w", c.name, err)
			}
			defer s3CleanupBucket(ctx, client, bucket)

			if c.policy != nil {
				if _, err := client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
					Bucket: aws.String(bucket),
					Policy: aws.String(c.policy(bucket)),
				}); err != nil {
					return fmt.Errorf("%s: PutBucketPolicy failed: %w", c.name, err)
				}
			}

			resp, err := client.GetBucketPolicyStatus(ctx, &s3.GetBucketPolicyStatusInput{
				Bucket: aws.String(bucket),
			})
			if err != nil {
				return fmt.Errorf("%s: GetBucketPolicyStatus failed: %w", c.name, err)
			}
			if resp.PolicyStatus == nil || resp.PolicyStatus.IsPublic == nil {
				return fmt.Errorf("%s: PolicyStatus or IsPublic is nil", c.name)
			}
			if *resp.PolicyStatus.IsPublic != c.wantPublic {
				return fmt.Errorf("%s: expected IsPublic %v, got %v", c.name, c.wantPublic, *resp.PolicyStatus.IsPublic)
			}
		}
		return nil
	}))

	// CreateBucket accepts the x-amz-acl and x-amz-object-ownership request
	// headers, which seed the bucket's ACL and ownership controls.
	results = append(results, r.RunTest("s3", "CreateBucket_ACLHeader", func() error {
		aclBucket := s3Bucket(ts, "create-acl")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(aclBucket),
			ACL:    types.BucketCannedACLPublicRead,
		})
		if err != nil {
			return fmt.Errorf("CreateBucket with ACL header failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, aclBucket)

		aclResp, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
			Bucket: aws.String(aclBucket),
		})
		if err != nil {
			return fmt.Errorf("GetBucketAcl failed: %w", err)
		}
		if !s3HasPublicReadGrant(aclResp.Grants) {
			return fmt.Errorf("expected an AllUsers READ grant from the public-read ACL header")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CreateBucket_ObjectOwnershipHeader", func() error {
		ownBucket := s3Bucket(ts, "create-own")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket:          aws.String(ownBucket),
			ObjectOwnership: types.ObjectOwnershipBucketOwnerPreferred,
		})
		if err != nil {
			return fmt.Errorf("CreateBucket with ObjectOwnership header failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, ownBucket)

		ownResp, err := client.GetBucketOwnershipControls(ctx, &s3.GetBucketOwnershipControlsInput{
			Bucket: aws.String(ownBucket),
		})
		if err != nil {
			return fmt.Errorf("GetBucketOwnershipControls failed: %w", err)
		}
		if ownResp.OwnershipControls == nil || len(ownResp.OwnershipControls.Rules) != 1 {
			return fmt.Errorf("expected one ownership rule, got %v", ownResp.OwnershipControls)
		}
		if ownResp.OwnershipControls.Rules[0].ObjectOwnership != types.ObjectOwnershipBucketOwnerPreferred {
			return fmt.Errorf("expected BucketOwnerPreferred ownership, got %s", ownResp.OwnershipControls.Rules[0].ObjectOwnership)
		}
		return nil
	}))

	// With Object Ownership set to BucketOwnerEnforced the bucket "accepts
	// only PUT requests that do not specify an ACL or PUT requests with
	// bucket owner full control ACLs"; other ACLs fail with 400
	// AccessControlListNotSupported, and later set-ACL requests fail too.
	results = append(results, r.RunTest("s3", "ObjectOwnershipEnforced_RejectsNonOwnerACLs", func() error {
		enforcedBucket := s3Bucket(ts, "own-enf")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket:          aws.String(enforcedBucket),
			ACL:             types.BucketCannedACLPublicRead,
			ObjectOwnership: types.ObjectOwnershipBucketOwnerEnforced,
		})
		if err == nil {
			return fmt.Errorf("expected AccessControlListNotSupported for public-read ACL with enforced ownership, got nil")
		}
		if err := expectS3Error(err, "AccessControlListNotSupported", http.StatusBadRequest); err != nil {
			return err
		}

		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket:          aws.String(s3Bucket(ts, "own-enf-ok")),
			ACL:             types.BucketCannedACL("bucket-owner-full-control"),
			ObjectOwnership: types.ObjectOwnershipBucketOwnerEnforced,
		})
		if err != nil {
			return fmt.Errorf("CreateBucket with bucket-owner-full-control under enforced ownership failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, s3Bucket(ts, "own-enf-ok"))

		// After the ownership setting is enforced, set/update ACL requests
		// fail regardless of the ACL they carry.
		_, err = client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
			Bucket: aws.String(s3Bucket(ts, "own-enf-ok")),
			ACL:    types.BucketCannedACLPublicRead,
		})
		if err == nil {
			return fmt.Errorf("expected AccessControlListNotSupported for PutBucketAcl on enforced bucket, got nil")
		}
		if err := expectAWSErrorCode(err, "AccessControlListNotSupported"); err != nil {
			return fmt.Errorf("expected AccessControlListNotSupported from PutBucketAcl: %v", err)
		}

		if _, err := s3PutObject(ctx, client, s3Bucket(ts, "own-enf-ok"), "acl-key.txt", "data"); err != nil {
			return err
		}
		_, err = client.PutObjectAcl(ctx, &s3.PutObjectAclInput{
			Bucket: aws.String(s3Bucket(ts, "own-enf-ok")),
			Key:    aws.String("acl-key.txt"),
			ACL:    types.ObjectCannedACLPublicRead,
		})
		if err == nil {
			return fmt.Errorf("expected AccessControlListNotSupported for PutObjectAcl on enforced bucket, got nil")
		}
		if err := expectAWSErrorCode(err, "AccessControlListNotSupported"); err != nil {
			return fmt.Errorf("expected AccessControlListNotSupported from PutObjectAcl: %v", err)
		}
		return nil
	}))

	// With ACLs disabled (BucketOwnerEnforced), uploads "with bucket owner
	// full control ACLs or uploads that don't specify an ACL" are accepted;
	// uploads specifying any other ACL fail with 400
	// AccessControlListNotSupported.
	results = append(results, r.RunTest("s3", "PutObject_ACLHeader_EnforcedBucket", func() error {
		enforcedBucket := s3Bucket(ts, "up-enf")
		if _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket:          aws.String(enforcedBucket),
			ObjectOwnership: types.ObjectOwnershipBucketOwnerEnforced,
		}); err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, enforcedBucket)

		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(enforcedBucket),
			Key:    aws.String("rejected.txt"),
			Body:   strings.NewReader("data"),
			ACL:    types.ObjectCannedACLPublicRead,
		})
		if err == nil {
			return fmt.Errorf("expected AccessControlListNotSupported for public-read upload ACL, got nil")
		}
		if err := expectS3Error(err, "AccessControlListNotSupported", http.StatusBadRequest); err != nil {
			return err
		}

		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(enforcedBucket),
			Key:    aws.String("accepted.txt"),
			Body:   strings.NewReader("data"),
			ACL:    types.ObjectCannedACLBucketOwnerFullControl,
		}); err != nil {
			return fmt.Errorf("PutObject with bucket-owner-full-control under enforced ownership failed: %w", err)
		}
		return nil
	}))

	// With ACLs enabled (ObjectWriter), the x-amz-acl upload header builds
	// the object's ACL.
	results = append(results, r.RunTest("s3", "PutObject_ACLHeader_AppliedOnACLsEnabledBucket", func() error {
		writerBucket := s3Bucket(ts, "up-ow")
		if _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket:          aws.String(writerBucket),
			ObjectOwnership: types.ObjectOwnershipObjectWriter,
		}); err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, writerBucket)

		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(writerBucket),
			Key:    aws.String("acl-object.txt"),
			Body:   strings.NewReader("data"),
			ACL:    types.ObjectCannedACLPublicRead,
		}); err != nil {
			return fmt.Errorf("PutObject with ACL header failed: %w", err)
		}

		aclResp, err := client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
			Bucket: aws.String(writerBucket),
			Key:    aws.String("acl-object.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObjectAcl failed: %w", err)
		}
		if !s3HasPublicReadGrant(aclResp.Grants) {
			return fmt.Errorf("expected an AllUsers READ grant from the public-read upload ACL")
		}
		return nil
	}))

	// CreateMultipartUpload validates the ACL headers of the initiating
	// request under the same object-ownership rules as PutObject.
	results = append(results, r.RunTest("s3", "CreateMultipartUpload_ACLHeader_EnforcedBucket", func() error {
		enforcedBucket := s3Bucket(ts, "mp-enf")
		if _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket:          aws.String(enforcedBucket),
			ObjectOwnership: types.ObjectOwnershipBucketOwnerEnforced,
		}); err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, enforcedBucket)

		_, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(enforcedBucket),
			Key:    aws.String("rejected.bin"),
			ACL:    types.ObjectCannedACLPublicRead,
		})
		if err == nil {
			return fmt.Errorf("expected AccessControlListNotSupported for public-read multipart ACL, got nil")
		}
		return expectS3Error(err, "AccessControlListNotSupported", http.StatusBadRequest)
	}))

	// CopyObject carries the ACL headers of a direct upload to the created
	// copy.
	results = append(results, r.RunTest("s3", "CopyObject_ACLHeader_AppliedOnACLsEnabledBucket", func() error {
		writerBucket := s3Bucket(ts, "cp-ow")
		if _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket:          aws.String(writerBucket),
			ObjectOwnership: types.ObjectOwnershipObjectWriter,
		}); err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, writerBucket)

		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(writerBucket),
			Key:    aws.String("src.txt"),
			Body:   strings.NewReader("data"),
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		if _, err := client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(writerBucket),
			Key:        aws.String("dst.txt"),
			CopySource: aws.String(writerBucket + "/src.txt"),
			ACL:        types.ObjectCannedACLPublicRead,
		}); err != nil {
			return fmt.Errorf("CopyObject with ACL header failed: %w", err)
		}

		aclResp, err := client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
			Bucket: aws.String(writerBucket),
			Key:    aws.String("dst.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObjectAcl failed: %w", err)
		}
		if !s3HasPublicReadGrant(aclResp.Grants) {
			return fmt.Errorf("expected an AllUsers READ grant on the copy from the public-read ACL")
		}
		return nil
	}))

	// --- Bucket inventory configurations ---

	destARN, cleanupDest, destErr := s3CreateReplicationDest(ctx, client, s3Bucket(ts, "inv-dest"))
	if destErr != nil {
		return append(results, TestResult{Service: "s3", TestName: "PutBucketInventoryConfiguration", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create inventory destination bucket: %v", destErr)})
	}
	defer cleanupDest()

	results = append(results, r.RunTest("s3", "PutBucketInventoryConfiguration", func() error {
		put := &types.InventoryConfiguration{
			Id:                     aws.String("report-main"),
			IsEnabled:              aws.Bool(false),
			IncludedObjectVersions: types.InventoryIncludedObjectVersionsAll,
			Filter:                 &types.InventoryFilter{Prefix: aws.String("docs/")},
			OptionalFields: []types.InventoryOptionalField{
				types.InventoryOptionalFieldSize,
				types.InventoryOptionalFieldLastModifiedDate,
				types.InventoryOptionalFieldETag,
				types.InventoryOptionalFieldStorageClass,
			},
			Schedule: &types.InventorySchedule{Frequency: types.InventoryFrequencyDaily},
			Destination: &types.InventoryDestination{
				S3BucketDestination: &types.InventoryS3BucketDestination{
					AccountId: aws.String("123456789012"),
					Bucket:    aws.String(destARN),
					Format:    types.InventoryFormatCsv,
					Prefix:    aws.String("inv"),
					Encryption: &types.InventoryEncryption{
						SSES3: &types.SSES3{},
					},
				},
			},
		}
		if _, err := client.PutBucketInventoryConfiguration(ctx, &s3.PutBucketInventoryConfigurationInput{
			Bucket:                 aws.String(bucketName),
			Id:                     aws.String("report-main"),
			InventoryConfiguration: put,
		}); err != nil {
			return fmt.Errorf("PutBucketInventoryConfiguration failed: %w", err)
		}

		got, err := client.GetBucketInventoryConfiguration(ctx, &s3.GetBucketInventoryConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("report-main"),
		})
		if err != nil {
			return fmt.Errorf("GetBucketInventoryConfiguration failed: %w", err)
		}
		cfg := got.InventoryConfiguration
		if cfg == nil || aws.ToString(cfg.Id) != "report-main" {
			return fmt.Errorf("round-trip lost the configuration id: %+v", cfg)
		}
		if cfg.IsEnabled == nil || aws.ToBool(cfg.IsEnabled) {
			return fmt.Errorf("round-trip lost the disabled IsEnabled flag: %+v", cfg)
		}
		if cfg.IncludedObjectVersions != types.InventoryIncludedObjectVersionsAll {
			return fmt.Errorf("round-trip lost IncludedObjectVersions: %+v", cfg)
		}
		if cfg.Filter == nil || aws.ToString(cfg.Filter.Prefix) != "docs/" {
			return fmt.Errorf("round-trip lost the filter prefix: %+v", cfg.Filter)
		}
		if cfg.Schedule == nil || cfg.Schedule.Frequency != types.InventoryFrequencyDaily {
			return fmt.Errorf("round-trip lost the schedule: %+v", cfg.Schedule)
		}
		dest := cfg.Destination
		if dest == nil || dest.S3BucketDestination == nil {
			return fmt.Errorf("round-trip lost the destination: %+v", dest)
		}
		gotDest := dest.S3BucketDestination
		if aws.ToString(gotDest.Bucket) != destARN || gotDest.Format != types.InventoryFormatCsv ||
			aws.ToString(gotDest.Prefix) != "inv" || aws.ToString(gotDest.AccountId) != "123456789012" {
			return fmt.Errorf("round-trip lost the destination details: %+v", gotDest)
		}
		if gotDest.Encryption == nil || gotDest.Encryption.SSES3 == nil {
			return fmt.Errorf("round-trip lost the SSE-S3 report encryption: %+v", gotDest.Encryption)
		}
		if len(cfg.OptionalFields) != 4 {
			return fmt.Errorf("round-trip lost optional fields: %+v", cfg.OptionalFields)
		}

		// A second Put under the same id is a full replacement.
		put.Schedule.Frequency = types.InventoryFrequencyWeekly
		put.OptionalFields = []types.InventoryOptionalField{types.InventoryOptionalFieldSize}
		if _, err := client.PutBucketInventoryConfiguration(ctx, &s3.PutBucketInventoryConfigurationInput{
			Bucket:                 aws.String(bucketName),
			Id:                     aws.String("report-main"),
			InventoryConfiguration: put,
		}); err != nil {
			return fmt.Errorf("replacement PutBucketInventoryConfiguration failed: %w", err)
		}
		replaced, err := client.GetBucketInventoryConfiguration(ctx, &s3.GetBucketInventoryConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("report-main"),
		})
		if err != nil {
			return fmt.Errorf("GetBucketInventoryConfiguration after replacement failed: %w", err)
		}
		if replaced.InventoryConfiguration.Schedule.Frequency != types.InventoryFrequencyWeekly ||
			len(replaced.InventoryConfiguration.OptionalFields) != 1 {
			return fmt.Errorf("replacement was not a full replacement: %+v", replaced.InventoryConfiguration)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetBucketInventoryConfiguration_NotFound", func() error {
		_, err := client.GetBucketInventoryConfiguration(ctx, &s3.GetBucketInventoryConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("no-such-inventory-id"),
		})
		if err == nil {
			return fmt.Errorf("expected NoSuchConfiguration, got nil")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) || apiErr.ErrorCode() != "NoSuchConfiguration" {
			return fmt.Errorf("expected NoSuchConfiguration, got: %v", err)
		}
		return nil
	}))

	// The destination bucket field is documented as a bucket ARN; a value
	// that does not parse as one, or whose resource part is not a valid
	// bucket name ("bucket/path" parses as an ARN but names no bucket), can
	// never receive a delivery, so the Put rejects both.
	results = append(results, r.RunTest("s3", "PutBucketInventoryConfiguration_DestinationArnValidated", func() error {
		for _, badDest := range []string{"plain-name", "arn:aws:s3:::bucket/path"} {
			_, err := client.PutBucketInventoryConfiguration(ctx, &s3.PutBucketInventoryConfigurationInput{
				Bucket: aws.String(bucketName),
				Id:     aws.String("bad-dest"),
				InventoryConfiguration: &types.InventoryConfiguration{
					Id:                     aws.String("bad-dest"),
					IsEnabled:              aws.Bool(true),
					IncludedObjectVersions: types.InventoryIncludedObjectVersionsCurrent,
					Schedule:               &types.InventorySchedule{Frequency: types.InventoryFrequencyDaily},
					Destination: &types.InventoryDestination{
						S3BucketDestination: &types.InventoryS3BucketDestination{
							Bucket: aws.String(badDest),
							Format: types.InventoryFormatCsv,
						},
					},
				},
			})
			if err == nil {
				return fmt.Errorf("expected InvalidArgument for destination %q, got nil", badDest)
			}
			var apiErr smithy.APIError
			if !errors.As(err, &apiErr) || apiErr.ErrorCode() != "InvalidArgument" {
				return fmt.Errorf("expected InvalidArgument for destination %q, got: %v", badDest, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListBucketInventoryConfigurations_Pagination", func() error {
		// AWS pages at 100 configurations; 102 forces a two-page traversal.
		const total = 102
		for i := 0; i < total; i++ {
			id := fmt.Sprintf("pg-i-%03d", i)
			if _, err := client.PutBucketInventoryConfiguration(ctx, &s3.PutBucketInventoryConfigurationInput{
				Bucket: aws.String(bucketName),
				Id:     aws.String(id),
				InventoryConfiguration: &types.InventoryConfiguration{
					Id:                     aws.String(id),
					IsEnabled:              aws.Bool(false),
					IncludedObjectVersions: types.InventoryIncludedObjectVersionsCurrent,
					Schedule:               &types.InventorySchedule{Frequency: types.InventoryFrequencyWeekly},
					Destination: &types.InventoryDestination{
						S3BucketDestination: &types.InventoryS3BucketDestination{
							Bucket: aws.String(destARN),
							Format: types.InventoryFormatCsv,
						},
					},
				},
			}); err != nil {
				return fmt.Errorf("filler Put %d failed: %w", i, err)
			}
		}
		defer func() {
			for i := 0; i < total; i++ {
				client.DeleteBucketInventoryConfiguration(ctx, &s3.DeleteBucketInventoryConfigurationInput{
					Bucket: aws.String(bucketName),
					Id:     aws.String(fmt.Sprintf("pg-i-%03d", i)),
				})
			}
		}()

		seen := make(map[string]bool)
		pages := 0
		var token *string
		for {
			resp, err := client.ListBucketInventoryConfigurations(ctx, &s3.ListBucketInventoryConfigurationsInput{
				Bucket:            aws.String(bucketName),
				ContinuationToken: token,
			})
			if err != nil {
				return fmt.Errorf("ListBucketInventoryConfigurations failed: %w", err)
			}
			pages++
			if len(resp.InventoryConfigurationList) > 100 {
				return fmt.Errorf("page %d returned %d configurations, over the documented 100 cap", pages, len(resp.InventoryConfigurationList))
			}
			for _, cfg := range resp.InventoryConfigurationList {
				seen[aws.ToString(cfg.Id)] = true
			}
			if !aws.ToBool(resp.IsTruncated) {
				if aws.ToString(token) != "" && aws.ToString(resp.ContinuationToken) != aws.ToString(token) {
					return fmt.Errorf("ContinuationToken echo mismatch: sent %q, got %q", aws.ToString(token), aws.ToString(resp.ContinuationToken))
				}
				break
			}
			token = resp.NextContinuationToken
			if token == nil || *token == "" {
				return fmt.Errorf("IsTruncated without NextContinuationToken")
			}
			if pages > 5 {
				return fmt.Errorf("pagination did not terminate after %d pages", pages)
			}
		}
		if pages < 2 {
			return fmt.Errorf("expected at least two pages for %d configurations, got %d", total, pages)
		}
		if len(seen) < total {
			return fmt.Errorf("traversal collected %d configurations, want at least %d (plus the round-trip one)", len(seen), total)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketInventoryConfiguration", func() error {
		if _, err := client.DeleteBucketInventoryConfiguration(ctx, &s3.DeleteBucketInventoryConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("report-main"),
		}); err != nil {
			return fmt.Errorf("DeleteBucketInventoryConfiguration failed: %w", err)
		}
		_, err := client.GetBucketInventoryConfiguration(ctx, &s3.GetBucketInventoryConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("report-main"),
		})
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) || apiErr.ErrorCode() != "NoSuchConfiguration" {
			return fmt.Errorf("expected NoSuchConfiguration after delete, got: %v", err)
		}
		// Deleting an absent id is idempotent.
		if _, err := client.DeleteBucketInventoryConfiguration(ctx, &s3.DeleteBucketInventoryConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("report-main"),
		}); err != nil {
			return fmt.Errorf("idempotent re-delete failed: %w", err)
		}
		return nil
	}))

	// --- Bucket metrics configurations ---

	results = append(results, r.RunTest("s3", "PutBucketMetricsConfiguration", func() error {
		cases := []struct {
			id     string
			filter types.MetricsFilter
		}{
			{id: "entire-bucket", filter: nil},
			{id: "docs-prefix", filter: &types.MetricsFilterMemberPrefix{Value: "docs/"}},
			{id: "classified-docs", filter: &types.MetricsFilterMemberTag{Value: types.Tag{Key: aws.String("class"), Value: aws.String("blue")}}},
			{id: "and-docs", filter: &types.MetricsFilterMemberAnd{Value: types.MetricsAndOperator{
				Prefix: aws.String("docs/"),
				Tags: []types.Tag{
					{Key: aws.String("priority"), Value: aws.String("high")},
					{Key: aws.String("class"), Value: aws.String("blue")},
				},
			}}},
		}
		for _, tc := range cases {
			put := &types.MetricsConfiguration{Id: aws.String(tc.id), Filter: tc.filter}
			if _, err := client.PutBucketMetricsConfiguration(ctx, &s3.PutBucketMetricsConfigurationInput{
				Bucket:               aws.String(bucketName),
				Id:                   aws.String(tc.id),
				MetricsConfiguration: put,
			}); err != nil {
				return fmt.Errorf("PutBucketMetricsConfiguration(%s) failed: %w", tc.id, err)
			}
			got, err := client.GetBucketMetricsConfiguration(ctx, &s3.GetBucketMetricsConfigurationInput{
				Bucket: aws.String(bucketName),
				Id:     aws.String(tc.id),
			})
			if err != nil {
				return fmt.Errorf("GetBucketMetricsConfiguration(%s) failed: %w", tc.id, err)
			}
			if aws.ToString(got.MetricsConfiguration.Id) != tc.id {
				return fmt.Errorf("round-trip lost id %s: %+v", tc.id, got.MetricsConfiguration)
			}
			if tc.filter == nil {
				if got.MetricsConfiguration.Filter != nil {
					return fmt.Errorf("nil filter round-tripped non-nil: %+v", got.MetricsConfiguration.Filter)
				}
				continue
			}
			f := got.MetricsConfiguration.Filter
			switch want := tc.filter.(type) {
			case *types.MetricsFilterMemberPrefix:
				gotPrefix, ok := f.(*types.MetricsFilterMemberPrefix)
				if !ok || gotPrefix.Value != want.Value {
					return fmt.Errorf("prefix filter round-trip mismatch: %+v", f)
				}
			case *types.MetricsFilterMemberTag:
				gotTag, ok := f.(*types.MetricsFilterMemberTag)
				if !ok || aws.ToString(gotTag.Value.Key) != "class" || aws.ToString(gotTag.Value.Value) != "blue" {
					return fmt.Errorf("tag filter round-trip mismatch: %+v", f)
				}
			case *types.MetricsFilterMemberAnd:
				gotAnd, ok := f.(*types.MetricsFilterMemberAnd)
				if !ok || aws.ToString(gotAnd.Value.Prefix) != "docs/" || len(gotAnd.Value.Tags) != 2 {
					return fmt.Errorf("and filter round-trip mismatch: %+v", f)
				}
			}
		}
		return nil
	}))

	// The model pins PutBucketMetricsConfiguration's success status at 200.
	// The SDK cannot observe it because the operation has no output shape,
	// so the pin goes through a raw request and then round-trips the stored
	// configuration through the SDK.
	results = append(results, r.RunTest("s3", "PutBucketMetricsConfiguration_HttpStatus200", func() error {
		endpoint := r.endpoint
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = "http://" + endpoint
		}
		body := "<MetricsConfiguration><Id>status-pin</Id></MetricsConfiguration>"
		req, err := http.NewRequest("PUT", endpoint+"/"+bucketName+"?metrics&id=status-pin", strings.NewReader(body))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		resp, err := testHTTPClient.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("expected HTTP 200 for PutBucketMetricsConfiguration, got %d: %s", resp.StatusCode, string(respBody))
		}
		got, err := client.GetBucketMetricsConfiguration(ctx, &s3.GetBucketMetricsConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("status-pin"),
		})
		if err != nil {
			return fmt.Errorf("GetBucketMetricsConfiguration(status-pin) failed: %w", err)
		}
		if aws.ToString(got.MetricsConfiguration.Id) != "status-pin" {
			return fmt.Errorf("round-trip lost id status-pin: %+v", got.MetricsConfiguration)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetBucketMetricsConfiguration_NotFound", func() error {
		_, err := client.GetBucketMetricsConfiguration(ctx, &s3.GetBucketMetricsConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("no-such-metrics-id"),
		})
		if err == nil {
			return fmt.Errorf("expected NoSuchConfiguration, got nil")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) || apiErr.ErrorCode() != "NoSuchConfiguration" {
			return fmt.Errorf("expected NoSuchConfiguration, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListBucketMetricsConfigurations", func() error {
		seen := make(map[string]bool)
		var token *string
		pages := 0
		for {
			resp, err := client.ListBucketMetricsConfigurations(ctx, &s3.ListBucketMetricsConfigurationsInput{
				Bucket:            aws.String(bucketName),
				ContinuationToken: token,
			})
			if err != nil {
				return fmt.Errorf("ListBucketMetricsConfigurations failed: %w", err)
			}
			pages++
			for _, cfg := range resp.MetricsConfigurationList {
				seen[aws.ToString(cfg.Id)] = true
			}
			if !aws.ToBool(resp.IsTruncated) {
				break
			}
			token = resp.NextContinuationToken
			if pages > 5 {
				return fmt.Errorf("pagination did not terminate after %d pages", pages)
			}
		}
		for _, id := range []string{"entire-bucket", "docs-prefix", "classified-docs", "and-docs"} {
			if !seen[id] {
				return fmt.Errorf("metrics list is missing configuration %q (seen: %v)", id, seen)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketMetricsConfiguration", func() error {
		if _, err := client.DeleteBucketMetricsConfiguration(ctx, &s3.DeleteBucketMetricsConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("docs-prefix"),
		}); err != nil {
			return fmt.Errorf("DeleteBucketMetricsConfiguration failed: %w", err)
		}
		_, err := client.GetBucketMetricsConfiguration(ctx, &s3.GetBucketMetricsConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("docs-prefix"),
		})
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) || apiErr.ErrorCode() != "NoSuchConfiguration" {
			return fmt.Errorf("expected NoSuchConfiguration after delete, got: %v", err)
		}
		if _, err := client.DeleteBucketMetricsConfiguration(ctx, &s3.DeleteBucketMetricsConfigurationInput{
			Bucket: aws.String(bucketName),
			Id:     aws.String("docs-prefix"),
		}); err != nil {
			return fmt.Errorf("idempotent re-delete failed: %w", err)
		}
		return nil
	}))

	// Request metrics: a metrics configuration with a prefix filter drives
	// per-filter CloudWatch metrics in the AWS/S3 namespace under the
	// BucketName and FilterId dimensions.
	results = append(results, r.RunTest("s3", "BucketMetricsConfiguration_RequestMetrics", func() error {
		metricsBucket := s3Bucket(ts, "req-metrics")
		if err := s3CreateBucket(ctx, client, metricsBucket); err != nil {
			return fmt.Errorf("failed to create the request-metrics bucket: %w", err)
		}
		defer s3CleanupBucket(ctx, client, metricsBucket)

		filterID := "data-prefix"
		if _, err := client.PutBucketMetricsConfiguration(ctx, &s3.PutBucketMetricsConfigurationInput{
			Bucket: aws.String(metricsBucket),
			Id:     aws.String(filterID),
			MetricsConfiguration: &types.MetricsConfiguration{
				Id:     aws.String(filterID),
				Filter: &types.MetricsFilterMemberPrefix{Value: "data/"},
			},
		}); err != nil {
			return fmt.Errorf("PutBucketMetricsConfiguration failed: %w", err)
		}

		body := "request-metrics-body"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(metricsBucket),
			Key:    aws.String("data/one.txt"),
			Body:   strings.NewReader(body),
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		if _, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(metricsBucket),
			Key:    aws.String("data/one.txt"),
		}); err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		// An object outside the filter must not feed the filter's metrics.
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(metricsBucket),
			Key:    aws.String("other/two.txt"),
			Body:   strings.NewReader(body),
		}); err != nil {
			return fmt.Errorf("outside-filter PutObject failed: %w", err)
		}

		cwCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load cloudwatch config: %w", err)
		}
		cwClient := cloudwatch.NewFromConfig(cwCfg)
		dimensions := []cwtypes.Dimension{
			{Name: aws.String("BucketName"), Value: aws.String(metricsBucket)},
			{Name: aws.String("FilterId"), Value: aws.String(filterID)},
		}

		// TEST_MODE compresses the aggregation window; poll until the
		// datapoints land.
		sums := map[string]float64{}
		deadline := time.Now().Add(60 * time.Second)
		for time.Now().Before(deadline) {
			sums = map[string]float64{}
			for _, metric := range []string{"AllRequests", "PutRequests", "GetRequests"} {
				resp, err := cwClient.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
					StartTime: aws.Time(time.Now().Add(-15 * time.Minute)),
					EndTime:   aws.Time(time.Now().Add(time.Minute)),
					ScanBy:    cwtypes.ScanByTimestampAscending,
					MetricDataQueries: []cwtypes.MetricDataQuery{{
						Id: aws.String("sum"),
						MetricStat: &cwtypes.MetricStat{
							Metric: &cwtypes.Metric{
								Namespace:  aws.String("AWS/S3"),
								MetricName: aws.String(metric),
								Dimensions: dimensions,
							},
							Period: aws.Int32(60),
							Stat:   aws.String("Sum"),
						},
					}},
				})
				if err != nil {
					return fmt.Errorf("GetMetricData(%s) failed: %w", metric, err)
				}
				for _, result := range resp.MetricDataResults {
					for _, value := range result.Values {
						sums[metric] += value
					}
				}
			}
			if sums["AllRequests"] >= 2 && sums["PutRequests"] >= 1 && sums["GetRequests"] >= 1 {
				return nil
			}
			time.Sleep(2 * time.Second)
		}
		return fmt.Errorf("request metrics not populated in time: %v", sums)
	}))

	// Bucket-plane requests: AllRequests counts every HTTP request made to
	// the bucket regardless of type, HeadRequests covers bucket heads, and
	// DeleteObjects requests classify as DeleteRequests.
	results = append(results, r.RunTest("s3", "BucketMetricsConfiguration_BucketPlaneRequests", func() error {
		planeBucket := s3Bucket(ts, "req-metrics-plane")
		if err := s3CreateBucket(ctx, client, planeBucket); err != nil {
			return fmt.Errorf("failed to create the bucket-plane metrics bucket: %w", err)
		}
		defer s3CleanupBucket(ctx, client, planeBucket)

		filterID := "whole-bucket"
		if _, err := client.PutBucketMetricsConfiguration(ctx, &s3.PutBucketMetricsConfigurationInput{
			Bucket: aws.String(planeBucket),
			Id:     aws.String(filterID),
			MetricsConfiguration: &types.MetricsConfiguration{
				Id: aws.String(filterID),
			},
		}); err != nil {
			return fmt.Errorf("PutBucketMetricsConfiguration failed: %w", err)
		}

		if _, err := client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(planeBucket)}); err != nil {
			return fmt.Errorf("HeadBucket failed: %w", err)
		}
		body := "bucket-plane-body"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(planeBucket),
			Key:    aws.String("gone/one.txt"),
			Body:   strings.NewReader(body),
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		if _, err := client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(planeBucket),
			Delete: &types.Delete{Objects: []types.ObjectIdentifier{{Key: aws.String("gone/one.txt")}}},
		}); err != nil {
			return fmt.Errorf("DeleteObjects failed: %w", err)
		}

		cwCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load cloudwatch config: %w", err)
		}
		cwClient := cloudwatch.NewFromConfig(cwCfg)
		dimensions := []cwtypes.Dimension{
			{Name: aws.String("BucketName"), Value: aws.String(planeBucket)},
			{Name: aws.String("FilterId"), Value: aws.String(filterID)},
		}

		sums := map[string]float64{}
		deadline := time.Now().Add(60 * time.Second)
		for time.Now().Before(deadline) {
			sums = map[string]float64{}
			for _, metric := range []string{"AllRequests", "HeadRequests", "DeleteRequests"} {
				resp, err := cwClient.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
					StartTime: aws.Time(time.Now().Add(-15 * time.Minute)),
					EndTime:   aws.Time(time.Now().Add(time.Minute)),
					ScanBy:    cwtypes.ScanByTimestampAscending,
					MetricDataQueries: []cwtypes.MetricDataQuery{{
						Id: aws.String("sum"),
						MetricStat: &cwtypes.MetricStat{
							Metric: &cwtypes.Metric{
								Namespace:  aws.String("AWS/S3"),
								MetricName: aws.String(metric),
								Dimensions: dimensions,
							},
							Period: aws.Int32(60),
							Stat:   aws.String("Sum"),
						},
					}},
				})
				if err != nil {
					return fmt.Errorf("GetMetricData(%s) failed: %w", metric, err)
				}
				for _, result := range resp.MetricDataResults {
					for _, value := range result.Values {
						sums[metric] += value
					}
				}
			}
			if sums["AllRequests"] >= 4 && sums["HeadRequests"] >= 1 && sums["DeleteRequests"] >= 1 {
				return nil
			}
			time.Sleep(2 * time.Second)
		}
		return fmt.Errorf("bucket-plane request metrics not populated in time: %v", sums)
	}))

	// CopyObject: GetRequests is incremented for the source of each copy,
	// so a copy out of a configured bucket must show up in that bucket's
	// GetRequests even though the request itself was a PUT to another
	// bucket.
	results = append(results, r.RunTest("s3", "CopyObject_SourceGetRequests", func() error {
		srcBucket := s3Bucket(ts, "copy-metrics-src")
		dstBucket := s3Bucket(ts, "copy-metrics-dst")
		for _, b := range []string{srcBucket, dstBucket} {
			if err := s3CreateBucket(ctx, client, b); err != nil {
				return fmt.Errorf("failed to create %s: %w", b, err)
			}
			defer s3CleanupBucket(ctx, client, b)
		}
		filterID := "whole-bucket"
		if _, err := client.PutBucketMetricsConfiguration(ctx, &s3.PutBucketMetricsConfigurationInput{
			Bucket:               aws.String(srcBucket),
			Id:                   aws.String(filterID),
			MetricsConfiguration: &types.MetricsConfiguration{Id: aws.String(filterID)},
		}); err != nil {
			return fmt.Errorf("PutBucketMetricsConfiguration failed: %w", err)
		}
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(srcBucket),
			Key:    aws.String("docs/orig.txt"),
			Body:   strings.NewReader("copy-source-body"),
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		if _, err := client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(dstBucket),
			Key:        aws.String("docs/copy.txt"),
			CopySource: aws.String(srcBucket + "/docs/orig.txt"),
		}); err != nil {
			return fmt.Errorf("CopyObject failed: %w", err)
		}

		cwCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load cloudwatch config: %w", err)
		}
		cwClient := cloudwatch.NewFromConfig(cwCfg)
		dimensions := []cwtypes.Dimension{
			{Name: aws.String("BucketName"), Value: aws.String(srcBucket)},
			{Name: aws.String("FilterId"), Value: aws.String(filterID)},
		}
		deadline := time.Now().Add(60 * time.Second)
		for time.Now().Before(deadline) {
			resp, err := cwClient.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
				StartTime: aws.Time(time.Now().Add(-15 * time.Minute)),
				EndTime:   aws.Time(time.Now().Add(time.Minute)),
				ScanBy:    cwtypes.ScanByTimestampAscending,
				MetricDataQueries: []cwtypes.MetricDataQuery{{
					Id: aws.String("sum"),
					MetricStat: &cwtypes.MetricStat{
						Metric: &cwtypes.Metric{
							Namespace:  aws.String("AWS/S3"),
							MetricName: aws.String("GetRequests"),
							Dimensions: dimensions,
						},
						Period: aws.Int32(60),
						Stat:   aws.String("Sum"),
					},
				}},
			})
			if err != nil {
				return fmt.Errorf("GetMetricData failed: %w", err)
			}
			sum := 0.0
			for _, result := range resp.MetricDataResults {
				for _, value := range result.Values {
					sum += value
				}
			}
			if sum >= 1 {
				return nil
			}
			time.Sleep(2 * time.Second)
		}
		return fmt.Errorf("copy-source GetRequests not populated in time")
	}))

	// --- Inventory report delivery (CSV, Parquet, and ORC end to end) ---

	results = append(results, r.RunTest("s3", "InventoryReportDelivery_CsvParquetOrc", func() error {
		srcBucket := s3Bucket(ts, "inv-src")
		if err := s3CreateBucket(ctx, client, srcBucket); err != nil {
			return fmt.Errorf("failed to create the report source bucket: %w", err)
		}
		defer s3CleanupBucket(ctx, client, srcBucket)

		for _, key := range []string{"kept/alpha.txt", "kept/beta.bin", "skipped/gamma.txt"} {
			if _, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(srcBucket),
				Key:    aws.String(key),
				Body:   strings.NewReader("inventory-report-body"),
			}); err != nil {
				return fmt.Errorf("PutObject %s failed: %w", key, err)
			}
		}

		csvID := "csv-delivery"
		pqID := "pq-delivery"
		orcID := "orc-delivery"
		inventoryCfg := func(id string, prefix string, format types.InventoryFormat) *types.InventoryConfiguration {
			return &types.InventoryConfiguration{
				Id:                     aws.String(id),
				IsEnabled:              aws.Bool(true),
				IncludedObjectVersions: types.InventoryIncludedObjectVersionsCurrent,
				Schedule:               &types.InventorySchedule{Frequency: types.InventoryFrequencyDaily},
				OptionalFields: []types.InventoryOptionalField{
					types.InventoryOptionalFieldSize,
					types.InventoryOptionalFieldLastModifiedDate,
					types.InventoryOptionalFieldStorageClass,
				},
				Destination: &types.InventoryDestination{
					S3BucketDestination: &types.InventoryS3BucketDestination{
						Bucket: aws.String(destARN),
						Format: format,
						Prefix: aws.String(prefix),
					},
				},
			}
		}
		// The CSV configuration filters to kept/, the Parquet one covers the
		// whole bucket through a destination prefix.
		if _, err := client.PutBucketInventoryConfiguration(ctx, &s3.PutBucketInventoryConfigurationInput{
			Bucket: aws.String(srcBucket), Id: aws.String(csvID),
			InventoryConfiguration: &types.InventoryConfiguration{
				Id:                     aws.String(csvID),
				IsEnabled:              aws.Bool(true),
				IncludedObjectVersions: types.InventoryIncludedObjectVersionsCurrent,
				Filter:                 &types.InventoryFilter{Prefix: aws.String("kept/")},
				Schedule:               &types.InventorySchedule{Frequency: types.InventoryFrequencyDaily},
				OptionalFields:         []types.InventoryOptionalField{types.InventoryOptionalFieldSize, types.InventoryOptionalFieldStorageClass},
				Destination: &types.InventoryDestination{
					S3BucketDestination: &types.InventoryS3BucketDestination{
						Bucket: aws.String(destARN),
						Format: types.InventoryFormatCsv,
					},
				},
			},
		}); err != nil {
			return fmt.Errorf("PutBucketInventoryConfiguration(csv) failed: %w", err)
		}
		if _, err := client.PutBucketInventoryConfiguration(ctx, &s3.PutBucketInventoryConfigurationInput{
			Bucket: aws.String(srcBucket), Id: aws.String(pqID),
			InventoryConfiguration: inventoryCfg(pqID, "reports", types.InventoryFormatParquet),
		}); err != nil {
			return fmt.Errorf("PutBucketInventoryConfiguration(parquet) failed: %w", err)
		}
		// The ORC configuration covers the whole bucket through a destination
		// prefix.
		if _, err := client.PutBucketInventoryConfiguration(ctx, &s3.PutBucketInventoryConfigurationInput{
			Bucket: aws.String(srcBucket), Id: aws.String(orcID),
			InventoryConfiguration: inventoryCfg(orcID, "orc", types.InventoryFormatOrc),
		}); err != nil {
			return fmt.Errorf("PutBucketInventoryConfiguration(orc) failed: %w", err)
		}
		defer client.DeleteBucketInventoryConfiguration(ctx, &s3.DeleteBucketInventoryConfigurationInput{Bucket: aws.String(srcBucket), Id: aws.String(csvID)})
		defer client.DeleteBucketInventoryConfiguration(ctx, &s3.DeleteBucketInventoryConfigurationInput{Bucket: aws.String(srcBucket), Id: aws.String(pqID)})
		defer client.DeleteBucketInventoryConfiguration(ctx, &s3.DeleteBucketInventoryConfigurationInput{Bucket: aws.String(srcBucket), Id: aws.String(orcID)})

		destBucket := s3Bucket(ts, "inv-dest")
		findManifest := func(prefix string) (string, bool) {
			resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket: aws.String(destBucket),
				Prefix: aws.String(prefix),
			})
			if err != nil {
				return "", false
			}
			for _, obj := range resp.Contents {
				if strings.HasSuffix(aws.ToString(obj.Key), "/manifest.json") {
					return aws.ToString(obj.Key), true
				}
			}
			return "", false
		}

		// TEST_MODE compresses the delivery cadence; a minute covers the
		// compressed period plus the worker tick.
		deadline := time.Now().Add(120 * time.Second)
		var csvManifestKey, pqManifestKey, orcManifestKey string
		for time.Now().Before(deadline) {
			if csvManifestKey == "" {
				if key, ok := findManifest(srcBucket + "/" + csvID + "/"); ok {
					csvManifestKey = key
				}
			}
			if pqManifestKey == "" {
				if key, ok := findManifest("reports/" + srcBucket + "/" + pqID + "/"); ok {
					pqManifestKey = key
				}
			}
			if orcManifestKey == "" {
				if key, ok := findManifest("orc/" + srcBucket + "/" + orcID + "/"); ok {
					orcManifestKey = key
				}
			}
			if csvManifestKey != "" && pqManifestKey != "" && orcManifestKey != "" {
				break
			}
			time.Sleep(2 * time.Second)
		}
		if csvManifestKey == "" || pqManifestKey == "" || orcManifestKey == "" {
			return fmt.Errorf("inventory reports not delivered in time (csv=%q parquet=%q orc=%q)", csvManifestKey, pqManifestKey, orcManifestKey)
		}

		verifyManifest := func(manifestKey string) (manifest struct {
			SourceBucket      string `json:"sourceBucket"`
			DestinationBucket string `json:"destinationBucket"`
			Version           string `json:"version"`
			FileFormat        string `json:"fileFormat"`
			FileSchema        string `json:"fileSchema"`
			Files             []struct {
				Key         string `json:"key"`
				Size        int64  `json:"size"`
				MD5Checksum string `json:"MD5checksum"`
			} `json:"files"`
		}, dataKey string, dataBytes []byte, err error) {
			mResp, err := client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(destBucket), Key: aws.String(manifestKey)})
			if err != nil {
				return manifest, "", nil, fmt.Errorf("GetObject manifest failed: %w", err)
			}
			manifestBytes, err := io.ReadAll(mResp.Body)
			mResp.Body.Close()
			if err != nil {
				return manifest, "", nil, fmt.Errorf("read manifest failed: %w", err)
			}
			if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
				return manifest, "", nil, fmt.Errorf("manifest json malformed: %w", err)
			}

			checksumResp, err := client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(destBucket), Key: aws.String(manifestKey + ".checksum")})
			if err != nil {
				return manifest, "", nil, fmt.Errorf("GetObject manifest.checksum failed: %w", err)
			}
			checksumBytes, err := io.ReadAll(checksumResp.Body)
			checksumResp.Body.Close()
			if err != nil {
				return manifest, "", nil, fmt.Errorf("read checksum failed: %w", err)
			}
			if strings.TrimSpace(string(checksumBytes)) != fmt.Sprintf("%x", md5.Sum(manifestBytes)) {
				return manifest, "", nil, fmt.Errorf("manifest.checksum %q does not match the manifest md5", strings.TrimSpace(string(checksumBytes)))
			}

			if len(manifest.Files) == 0 {
				return manifest, "", nil, fmt.Errorf("manifest lists no data files")
			}
			dataKey = manifest.Files[0].Key
			dResp, err := client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(destBucket), Key: aws.String(dataKey)})
			if err != nil {
				return manifest, "", nil, fmt.Errorf("GetObject data file failed: %w", err)
			}
			dataBytes, err = io.ReadAll(dResp.Body)
			dResp.Body.Close()
			if err != nil {
				return manifest, "", nil, fmt.Errorf("read data file failed: %w", err)
			}
			if manifest.Files[0].Size != int64(len(dataBytes)) {
				return manifest, "", nil, fmt.Errorf("manifest size %d != actual %d", manifest.Files[0].Size, len(dataBytes))
			}
			if manifest.Files[0].MD5Checksum != fmt.Sprintf("%x", md5.Sum(dataBytes)) {
				return manifest, "", nil, fmt.Errorf("manifest MD5checksum does not match the data file md5")
			}
			return manifest, dataKey, dataBytes, nil
		}

		csvManifest, csvDataKey, csvData, err := verifyManifest(csvManifestKey)
		if err != nil {
			return err
		}
		if csvManifest.SourceBucket != srcBucket || csvManifest.DestinationBucket != destARN {
			return fmt.Errorf("csv manifest bucket fields wrong: %+v", csvManifest)
		}
		if csvManifest.Version != "2016-11-30" || csvManifest.FileFormat != "CSV" {
			return fmt.Errorf("csv manifest version/format wrong: %+v", csvManifest)
		}
		if csvManifest.FileSchema != "Bucket, Key, Size, StorageClass" {
			return fmt.Errorf("csv fileSchema = %q, want the selected-column order", csvManifest.FileSchema)
		}
		if !strings.HasPrefix(csvDataKey, srcBucket+"/"+csvID+"/data/") || !strings.HasSuffix(csvDataKey, ".csv.gz") {
			return fmt.Errorf("csv data key = %q, want the documented layout", csvDataKey)
		}
		gzipReader, err := gzip.NewReader(bytes.NewReader(csvData))
		if err != nil {
			return fmt.Errorf("csv data is not gzip: %w", err)
		}
		csvText, err := io.ReadAll(gzipReader)
		if err != nil {
			return fmt.Errorf("csv gunzip failed: %w", err)
		}
		if !strings.Contains(string(csvText), srcBucket+",kept/alpha.txt,") || !strings.Contains(string(csvText), srcBucket+",kept/beta.bin,") {
			return fmt.Errorf("csv report lost the kept objects: %s", csvText)
		}
		if strings.Contains(string(csvText), "skipped/gamma.txt") {
			return fmt.Errorf("csv report included an object outside the filter: %s", csvText)
		}

		pqManifest, pqDataKey, pqData, err := verifyManifest(pqManifestKey)
		if err != nil {
			return err
		}
		if pqManifest.FileFormat != "Parquet" || !strings.Contains(pqManifest.FileSchema, "required binary bucket (UTF8)") {
			return fmt.Errorf("parquet manifest wrong: %+v", pqManifest)
		}
		if !strings.HasPrefix(pqDataKey, "reports/"+srcBucket+"/"+pqID+"/data/") || !strings.HasSuffix(pqDataKey, ".parquet") {
			return fmt.Errorf("parquet data key = %q, want the prefixed layout", pqDataKey)
		}
		if len(pqData) < 8 || string(pqData[:4]) != "PAR1" || string(pqData[len(pqData)-4:]) != "PAR1" {
			return fmt.Errorf("parquet data file does not carry the parquet magic bytes")
		}

		orcManifest, orcDataKey, orcData, err := verifyManifest(orcManifestKey)
		if err != nil {
			return err
		}
		if orcManifest.FileFormat != "ORC" || !strings.HasPrefix(orcManifest.FileSchema, "struct<bucket:string,key:string") {
			return fmt.Errorf("orc manifest wrong: %+v", orcManifest)
		}
		if !strings.HasPrefix(orcDataKey, "orc/"+srcBucket+"/"+orcID+"/data/") || !strings.HasSuffix(orcDataKey, ".orc") {
			return fmt.Errorf("orc data key = %q, want the prefixed layout", orcDataKey)
		}
		orcReader, err := orc.NewReader(bytes.NewReader(orcData))
		if err != nil {
			return fmt.Errorf("orc reader: %w", err)
		}
		defer orcReader.Close()
		cursor := orcReader.Select("bucket", "key", "size")
		orcRows := map[string]int64{}
		orcBuckets := map[string]bool{}
		for cursor.Stripes() {
			for cursor.Next() {
				row := cursor.Row()
				orcRows[row[1].(string)] = row[2].(int64)
				orcBuckets[row[0].(string)] = true
			}
		}
		if err := cursor.Err(); err != nil {
			return fmt.Errorf("orc cursor: %w", err)
		}
		if len(orcRows) != 3 {
			return fmt.Errorf("orc report rows = %d (%v), want 3", len(orcRows), orcRows)
		}
		if len(orcBuckets) != 1 || !orcBuckets[srcBucket] {
			return fmt.Errorf("orc report bucket column = %v, want only %q", orcBuckets, srcBucket)
		}
		for _, key := range []string{"kept/alpha.txt", "kept/beta.bin", "skipped/gamma.txt"} {
			if orcRows[key] != 21 {
				return fmt.Errorf("orc row %s = %d, want 21", key, orcRows[key])
			}
		}
		return nil
	}))

	// Delivered reports honour the configuration's InventoryEncryption
	// choice: the SSE-S3 delivery carries AES256 and the SSE-KMS one
	// aws:kms plus the configured key.
	results = append(results, r.RunTest("s3", "InventoryReportDelivery_Encrypted", func() error {
		_, keyID, keyCleanup, keyErr := r.s3CreateTestKMSKey(ctx, "inventory encrypted delivery")
		if keyErr != nil {
			return keyErr
		}
		defer keyCleanup()

		srcBucket := s3Bucket(ts, "inv-enc-src")
		if err := s3CreateBucket(ctx, client, srcBucket); err != nil {
			return fmt.Errorf("failed to create the encrypted-delivery source bucket: %w", err)
		}
		defer s3CleanupBucket(ctx, client, srcBucket)
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:               aws.String(srcBucket),
			Key:                  aws.String("kept/one.txt"),
			Body:                 strings.NewReader("encrypted-delivery-body"),
			ServerSideEncryption: types.ServerSideEncryptionAwsKms,
			SSEKMSKeyId:          aws.String(keyID),
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(srcBucket),
			Key:    aws.String("plain-two.txt"),
			Body:   strings.NewReader("plain body"),
		}); err != nil {
			return fmt.Errorf("plain PutObject failed: %w", err)
		}

		encDestBucket := s3Bucket(ts, "inv-enc-dest")
		if err := s3CreateBucket(ctx, client, encDestBucket); err != nil {
			return fmt.Errorf("failed to create the encrypted-delivery destination bucket: %w", err)
		}
		defer s3CleanupBucket(ctx, client, encDestBucket)
		encDestARN := fmt.Sprintf("arn:aws:s3:::%s", encDestBucket)

		encryptedCfg := func(id string, encryption *types.InventoryEncryption) *types.InventoryConfiguration {
			return &types.InventoryConfiguration{
				Id:                     aws.String(id),
				IsEnabled:              aws.Bool(true),
				IncludedObjectVersions: types.InventoryIncludedObjectVersionsCurrent,
				OptionalFields:         []types.InventoryOptionalField{types.InventoryOptionalFieldSize, types.InventoryOptionalFieldBucketKeyStatus},
				Schedule:               &types.InventorySchedule{Frequency: types.InventoryFrequencyDaily},
				Destination: &types.InventoryDestination{
					S3BucketDestination: &types.InventoryS3BucketDestination{
						Bucket:     aws.String(encDestARN),
						Format:     types.InventoryFormatCsv,
						Encryption: encryption,
					},
				},
			}
		}
		s3ID, kmsID := "enc-s3", "enc-kms"
		for id, cfg := range map[string]*types.InventoryConfiguration{
			s3ID:  encryptedCfg(s3ID, &types.InventoryEncryption{SSES3: &types.SSES3{}}),
			kmsID: encryptedCfg(kmsID, &types.InventoryEncryption{SSEKMS: &types.SSEKMS{KeyId: aws.String(keyID)}}),
		} {
			if _, err := client.PutBucketInventoryConfiguration(ctx, &s3.PutBucketInventoryConfigurationInput{
				Bucket: aws.String(srcBucket), Id: aws.String(id), InventoryConfiguration: cfg,
			}); err != nil {
				return fmt.Errorf("PutBucketInventoryConfiguration(%s) failed: %w", id, err)
			}
			defer client.DeleteBucketInventoryConfiguration(ctx, &s3.DeleteBucketInventoryConfigurationInput{Bucket: aws.String(srcBucket), Id: aws.String(id)})
		}

		findManifest := func(prefix string) (string, bool) {
			resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket: aws.String(encDestBucket),
				Prefix: aws.String(prefix),
			})
			if err != nil {
				return "", false
			}
			for _, obj := range resp.Contents {
				if strings.HasSuffix(aws.ToString(obj.Key), "/manifest.json") {
					return aws.ToString(obj.Key), true
				}
			}
			return "", false
		}

		// TEST_MODE compresses the delivery cadence; a minute covers the
		// compressed period plus the worker tick.
		deadline := time.Now().Add(120 * time.Second)
		manifests := map[string]string{}
		for time.Now().Before(deadline) && len(manifests) < 2 {
			for _, id := range []string{s3ID, kmsID} {
				if manifests[id] != "" {
					continue
				}
				if key, ok := findManifest(srcBucket + "/" + id + "/"); ok {
					manifests[id] = key
				}
			}
			time.Sleep(2 * time.Second)
		}
		if len(manifests) < 2 {
			return fmt.Errorf("encrypted inventory reports not delivered in time (seen: %v)", manifests)
		}

		wantEncryption := map[string]types.ServerSideEncryption{
			s3ID:  types.ServerSideEncryptionAes256,
			kmsID: types.ServerSideEncryptionAwsKms,
		}
		for id, manifestKey := range manifests {
			mResp, err := client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(encDestBucket), Key: aws.String(manifestKey)})
			if err != nil {
				return fmt.Errorf("GetObject %s manifest failed: %w", id, err)
			}
			manifestBytes, err := io.ReadAll(mResp.Body)
			mResp.Body.Close()
			if err != nil {
				return fmt.Errorf("read %s manifest failed: %w", id, err)
			}
			if mResp.ServerSideEncryption != wantEncryption[id] {
				return fmt.Errorf("%s manifest encryption = %s, want %s", id, mResp.ServerSideEncryption, wantEncryption[id])
			}
			var manifest struct {
				Files []struct {
					Key string `json:"key"`
				} `json:"files"`
			}
			if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
				return fmt.Errorf("%s manifest json malformed: %w", id, err)
			}
			if len(manifest.Files) == 0 {
				return fmt.Errorf("%s manifest lists no data files", id)
			}
			dResp, err := client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(encDestBucket), Key: aws.String(manifest.Files[0].Key)})
			if err != nil {
				return fmt.Errorf("GetObject %s data file failed: %w", id, err)
			}
			dataBytes, err := io.ReadAll(dResp.Body)
			dResp.Body.Close()
			if err != nil {
				return fmt.Errorf("read %s data file failed: %w", id, err)
			}
			if dResp.ServerSideEncryption != wantEncryption[id] {
				return fmt.Errorf("%s data file encryption = %s, want %s", id, dResp.ServerSideEncryption, wantEncryption[id])
			}
			if id == kmsID && aws.ToString(dResp.SSEKMSKeyId) != keyID {
				return fmt.Errorf("%s data file kms key = %q, want %q", id, aws.ToString(dResp.SSEKMSKeyId), keyID)
			}
			gz, err := gzip.NewReader(bytes.NewReader(dataBytes))
			if err != nil {
				return fmt.Errorf("%s data file is not gzip: %w", id, err)
			}
			csvText, err := io.ReadAll(gz)
			if err != nil {
				return fmt.Errorf("%s csv gunzip failed: %w", id, err)
			}
			if !strings.Contains(string(csvText), srcBucket+",kept/one.txt,") {
				return fmt.Errorf("%s report lost the object row: %s", id, csvText)
			}
			// The optional BucketKeyStatus column describes each listed
			// object: the SSE-KMS object reports DISABLED (the platform has
			// no bucket keys) and the unencrypted one carries no status.
			checkRow := func(fragment, wantSuffix string) error {
				for _, line := range strings.Split(string(csvText), "\n") {
					if strings.Contains(line, fragment) {
						if !strings.HasSuffix(line, wantSuffix) {
							return fmt.Errorf("%s report row %q does not end with %q", id, line, wantSuffix)
						}
						return nil
					}
				}
				return fmt.Errorf("%s report lost the %q row: %s", id, fragment, csvText)
			}
			if err := checkRow(",kept/one.txt,", ",DISABLED"); err != nil {
				return err
			}
			if err := checkRow(",plain-two.txt,", ","); err != nil {
				return err
			}
		}
		return nil
	}))

	return results
}
