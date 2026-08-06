package testutil

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"vorpalstacks-sdk-tests/config"
)

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

	results = append(results, r.RunTest("s3", "PutBucketVersioning_Enabled", func() error {
		_, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusEnabled,
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketVersioning failed: %w", err)
		}
		getResp, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketVersioning failed: %w", err)
		}
		if getResp.Status != types.BucketVersioningStatusEnabled {
			return fmt.Errorf("expected Status Enabled, got %s", getResp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketVersioning_Suspended", func() error {
		_, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusSuspended,
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketVersioning failed: %w", err)
		}
		getResp, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketVersioning failed: %w", err)
		}
		if getResp.Status != types.BucketVersioningStatusSuspended {
			return fmt.Errorf("expected Status Suspended, got %s", getResp.Status)
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
		repBucket := s3Bucket(ts, "repl-dest")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(repBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket (replication-dest) failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, repBucket)

		repArn := "arn:aws:s3:::" + repBucket
		_, err = client.PutBucketReplication(ctx, &s3.PutBucketReplicationInput{
			Bucket: aws.String(bucketName),
			ReplicationConfiguration: &types.ReplicationConfiguration{
				Role: aws.String(fmt.Sprintf("arn:aws:iam::%s:role/s3-replication", r.accountID)),
				Rules: []types.ReplicationRule{
					{
						ID:       aws.String("rule-1"),
						Status:   types.ReplicationRuleStatusEnabled,
						Priority: aws.Int32(1),
						Filter: &types.ReplicationRuleFilter{
							Prefix: aws.String(""),
						},
						Destination: &types.Destination{
							Bucket: aws.String(repArn),
						},
						DeleteMarkerReplication: &types.DeleteMarkerReplication{
							Status: types.DeleteMarkerReplicationStatusDisabled,
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketReplication failed: %w", err)
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
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(repBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket (replication-copy-dest) failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, repBucket)

		repArn := "arn:aws:s3:::" + repBucket
		_, err = client.PutBucketReplication(ctx, &s3.PutBucketReplicationInput{
			Bucket: aws.String(bucketName),
			ReplicationConfiguration: &types.ReplicationConfiguration{
				Role: aws.String(fmt.Sprintf("arn:aws:iam::%s:role/s3-replication", r.accountID)),
				Rules: []types.ReplicationRule{{
					ID:       aws.String("copy-rule"),
					Status:   types.ReplicationRuleStatusEnabled,
					Priority: aws.Int32(1),
					Filter:   &types.ReplicationRuleFilter{Prefix: aws.String("")},
					Destination: &types.Destination{
						Bucket: aws.String(repArn),
					},
					DeleteMarkerReplication: &types.DeleteMarkerReplication{
						Status: types.DeleteMarkerReplicationStatusDisabled,
					},
				}},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketReplication failed: %w", err)
		}

		repKey := "repl-test-object.txt"
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(repKey),
			Body:   strings.NewReader("replicated content"),
		})
		if err != nil {
			return fmt.Errorf("PutObject to source bucket failed: %w", err)
		}

		if !waitForObject(ctx, client, repBucket, repKey, 5*time.Second) {
			return fmt.Errorf("GetObject from destination bucket failed (replication did not occur)")
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(repBucket),
			Key:    aws.String(repKey),
		})
		if err != nil {
			return fmt.Errorf("GetObject from destination bucket failed (replication did not occur): %w", err)
		}
		defer getResp.Body.Close()
		if getResp.ContentLength == nil || *getResp.ContentLength != int64(len("replicated content")) {
			return fmt.Errorf("replicated object size mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "Replication_TagFilter", func() error {
		tagBucket := s3Bucket(ts, "repl-tag-dest")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(tagBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket (repl-tag-dest) failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, tagBucket)

		tagArn := "arn:aws:s3:::" + tagBucket
		_, err = client.PutBucketReplication(ctx, &s3.PutBucketReplicationInput{
			Bucket: aws.String(bucketName),
			ReplicationConfiguration: &types.ReplicationConfiguration{
				Role: aws.String(fmt.Sprintf("arn:aws:iam::%s:role/s3-replication", r.accountID)),
				Rules: []types.ReplicationRule{{
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
				}},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketReplication (tag-filter) failed: %w", err)
		}

		// Object WITH the matching tag should replicate.
		matchedKey := "tag-matched.txt"
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:  aws.String(bucketName),
			Key:     aws.String(matchedKey),
			Body:    strings.NewReader("tagged content"),
			Tagging: aws.String("replicate=true"),
		})
		if err != nil {
			return fmt.Errorf("PutObject (tagged) failed: %w", err)
		}

		// Object WITHOUT the matching tag should NOT replicate.
		unmatchedKey := "tag-unmatched.txt"
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(unmatchedKey),
			Body:   strings.NewReader("untagged content"),
		})
		if err != nil {
			return fmt.Errorf("PutObject (untagged) failed: %w", err)
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
		_, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusEnabled,
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketVersioning (source) failed: %w", err)
		}

		dmBucket := s3Bucket(ts, "repl-dm-dest")
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(dmBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket (repl-dm-dest) failed: %w", err)
		}

		// Enable versioning on destination so delete markers are supported.
		_, err = client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(dmBucket),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusEnabled,
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketVersioning (dest) failed: %w", err)
		}

		dmArn := "arn:aws:s3:::" + dmBucket
		_, err = client.PutBucketReplication(ctx, &s3.PutBucketReplicationInput{
			Bucket: aws.String(bucketName),
			ReplicationConfiguration: &types.ReplicationConfiguration{
				Role: aws.String(fmt.Sprintf("arn:aws:iam::%s:role/s3-replication", r.accountID)),
				Rules: []types.ReplicationRule{{
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
				}},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketReplication (dm) failed: %w", err)
		}

		dmKey := "dm-test-object.txt"
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(dmKey),
			Body:   strings.NewReader("dm content"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
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

		// Restore versioning to Suspended (the state before this test changed it)
		// so that subsequent tests see the same bucket state.
		_, err = client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusSuspended,
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketVersioning (restore suspended) failed: %w", err)
		}

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

	return results
}
