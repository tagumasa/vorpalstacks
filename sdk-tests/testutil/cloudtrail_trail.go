package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runCloudTrailTrailTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	var trailName string
	results = append(results, r.RunTest("cloudtrail", "CreateTrail", func() error {
		trailName = tc.uniqueName("test-trail")
		if err := tc.ensureTrailBucket("test-bucket"); err != nil {
			return err
		}
		resp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:                       aws.String(trailName),
			S3BucketName:               aws.String("test-bucket"),
			IncludeGlobalServiceEvents: aws.Bool(true),
			IsMultiRegionTrail:         aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != trailName {
			return fmt.Errorf("trail name mismatch, got %v", resp.Name)
		}
		if resp.S3BucketName == nil || *resp.S3BucketName != "test-bucket" {
			return fmt.Errorf("S3 bucket name mismatch, got %v", resp.S3BucketName)
		}
		if resp.TrailARN == nil || *resp.TrailARN == "" {
			return fmt.Errorf("TrailARN should be set")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "ListTrails", func() error {
		resp, err := tc.client.ListTrails(tc.ctx, &cloudtrail.ListTrailsInput{})
		if err != nil {
			return err
		}
		if resp.Trails == nil {
			return fmt.Errorf("trails list is nil")
		}
		found := false
		for _, t := range resp.Trails {
			if t.Name != nil && *t.Name == trailName {
				found = true
				if t.TrailARN == nil || *t.TrailARN == "" {
					return fmt.Errorf("TrailARN is empty for listed trail")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created trail %s not found in ListTrails", trailName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrail", func() error {
		resp, err := tc.client.GetTrail(tc.ctx, &cloudtrail.GetTrailInput{
			Name: aws.String(trailName),
		})
		if err != nil {
			return err
		}
		if resp.Trail == nil {
			return fmt.Errorf("trail is nil")
		}
		if resp.Trail.Name == nil || *resp.Trail.Name != trailName {
			return fmt.Errorf("trail name mismatch, got %v", resp.Trail.Name)
		}
		if resp.Trail.S3BucketName == nil || *resp.Trail.S3BucketName != "test-bucket" {
			return fmt.Errorf("S3 bucket name mismatch, got %v", resp.Trail.S3BucketName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DescribeTrails", func() error {
		resp, err := tc.client.DescribeTrails(tc.ctx, &cloudtrail.DescribeTrailsInput{
			TrailNameList: []string{trailName},
		})
		if err != nil {
			return err
		}
		if len(resp.TrailList) != 1 {
			return fmt.Errorf("expected 1 trail, got %d", len(resp.TrailList))
		}
		if resp.TrailList[0].Name == nil || *resp.TrailList[0].Name != trailName {
			return fmt.Errorf("trail name mismatch, got %v", resp.TrailList[0].Name)
		}
		if resp.TrailList[0].S3BucketName == nil || *resp.TrailList[0].S3BucketName != "test-bucket" {
			return fmt.Errorf("S3 bucket name mismatch, got %v", resp.TrailList[0].S3BucketName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "UpdateTrail", func() error {
		if err := tc.ensureTrailBucket("updated-bucket"); err != nil {
			return err
		}
		resp, err := tc.client.UpdateTrail(tc.ctx, &cloudtrail.UpdateTrailInput{
			Name:         aws.String(trailName),
			S3BucketName: aws.String("updated-bucket"),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != trailName {
			return fmt.Errorf("trail name mismatch, got %v", resp.Name)
		}
		if resp.S3BucketName == nil || *resp.S3BucketName != "updated-bucket" {
			return fmt.Errorf("S3 bucket name not updated, got %v", resp.S3BucketName)
		}
		getResp, err := tc.client.GetTrail(tc.ctx, &cloudtrail.GetTrailInput{Name: aws.String(trailName)})
		if err != nil {
			return fmt.Errorf("get trail after update: %v", err)
		}
		if getResp.Trail.S3BucketName == nil || *getResp.Trail.S3BucketName != "updated-bucket" {
			return fmt.Errorf("S3 bucket not persisted after update, got %v", getResp.Trail.S3BucketName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DeleteTrail", func() error {
		resp, err := tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{
			Name: aws.String(trailName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		_, err = tc.client.GetTrail(tc.ctx, &cloudtrail.GetTrailInput{Name: aws.String(trailName)})
		if err == nil {
			return fmt.Errorf("expected TrailNotFoundException after delete")
		}
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// Deleting a trail that is currently logging succeeds: the DeleteTrail
	// reference documents no StopLogging precondition ("Deleting a
	// multi-Region trail will stop logging of events in all AWS Regions
	// enabled in your AWS account").
	results = append(results, r.RunTest("cloudtrail", "DeleteTrail_WhileLogging", func() error {
		name := tc.uniqueName("del-logging")
		if _, err := tc.createTrail(name, "del-logging-bucket"); err != nil {
			return err
		}
		if _, err := tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)}); err != nil {
			return fmt.Errorf("StartLogging failed: %w", err)
		}
		if _, err := tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)}); err != nil {
			return fmt.Errorf("DeleteTrail on a logging trail failed: %w", err)
		}
		_, err := tc.client.GetTrail(tc.ctx, &cloudtrail.GetTrailInput{Name: aws.String(name)})
		return AssertErrorContains(err, "TrailNotFoundException")
	}))

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_DefaultFields", func() error {
		name := tc.uniqueName("defaults")
		defer tc.deleteTrail(name)

		resp, err := tc.createTrail(name, "defaults-bucket")
		if err != nil {
			return err
		}
		if resp.IncludeGlobalServiceEvents == nil || !*resp.IncludeGlobalServiceEvents {
			return fmt.Errorf("IncludeGlobalServiceEvents should default to true")
		}
		if resp.IsMultiRegionTrail != nil && *resp.IsMultiRegionTrail {
			return fmt.Errorf("IsMultiRegionTrail should default to false")
		}
		if resp.LogFileValidationEnabled != nil && *resp.LogFileValidationEnabled {
			return fmt.Errorf("LogFileValidationEnabled should default to false")
		}
		if resp.IsOrganizationTrail != nil && *resp.IsOrganizationTrail {
			return fmt.Errorf("IsOrganizationTrail should default to false")
		}
		if resp.TrailARN == nil || *resp.TrailARN == "" {
			return fmt.Errorf("TrailARN should be set")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_Duplicate", func() error {
		name := tc.uniqueName("dup-trail")
		defer tc.deleteTrail(name)

		_, err := tc.createTrail(name, "dup-bucket")
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}

		_, err = tc.createTrail(name, "dup-bucket")
		if err := AssertErrorContains(err, "TrailAlreadyExists"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrail_ByARN", func() error {
		name := tc.uniqueName("arn-trail")
		defer tc.deleteTrail(name)

		createResp, err := tc.createTrail(name, "arn-bucket")
		if err != nil {
			return err
		}
		if createResp.TrailARN == nil {
			return fmt.Errorf("trail ARN is nil")
		}

		getResp, err := tc.client.GetTrail(tc.ctx, &cloudtrail.GetTrailInput{
			Name: createResp.TrailARN,
		})
		if err != nil {
			return fmt.Errorf("get by ARN: %v", err)
		}
		if getResp.Trail == nil || getResp.Trail.Name == nil || *getResp.Trail.Name != name {
			return fmt.Errorf("trail name mismatch after get by ARN")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DescribeTrails_ByARN", func() error {
		name := tc.uniqueName("desc-arn")
		defer tc.deleteTrail(name)

		createResp, err := tc.createTrail(name, "desc-arn-bucket")
		if err != nil {
			return err
		}

		if createResp.TrailARN == nil {
			return fmt.Errorf("CreateTrail did not return TrailARN")
		}

		resp, err := tc.client.DescribeTrails(tc.ctx, &cloudtrail.DescribeTrailsInput{
			TrailNameList: []string{*createResp.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("describe by ARN: %v", err)
		}
		if len(resp.TrailList) != 1 {
			return fmt.Errorf("expected 1 trail, got %d", len(resp.TrailList))
		}
		if resp.TrailList[0].Name == nil || *resp.TrailList[0].Name != name {
			return fmt.Errorf("trail name mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DescribeTrails_ListAll", func() error {
		name := tc.uniqueName("listall")
		defer tc.deleteTrail(name)

		_, err := tc.createTrail(name, "listall-bucket")
		if err != nil {
			return err
		}

		resp, err := tc.client.DescribeTrails(tc.ctx, &cloudtrail.DescribeTrailsInput{})
		if err != nil {
			return err
		}
		if resp.TrailList == nil {
			return fmt.Errorf("trail list is nil")
		}
		if len(resp.TrailList) == 0 {
			return fmt.Errorf("expected at least 1 trail")
		}
		found := false
		for _, t := range resp.TrailList {
			if t.Name != nil && *t.Name == name {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("trail %s not found in ListAll result", name)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "UpdateTrail_EnableLogFileValidation", func() error {
		name := tc.uniqueName("lfv")
		defer tc.deleteTrail(name)

		_, err := tc.createTrail(name, "lfv-bucket")
		if err != nil {
			return err
		}

		resp, err := tc.client.UpdateTrail(tc.ctx, &cloudtrail.UpdateTrailInput{
			Name:                    aws.String(name),
			EnableLogFileValidation: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp.LogFileValidationEnabled == nil || !*resp.LogFileValidationEnabled {
			return fmt.Errorf("expected LogFileValidationEnabled=true")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_WithLogFileValidation", func() error {
		name := tc.uniqueName("lfv-create")
		defer tc.deleteTrail(name)

		if err := tc.ensureTrailBucket("lfv-create-bucket"); err != nil {
			return err
		}
		resp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:                       aws.String(name),
			S3BucketName:               aws.String("lfv-create-bucket"),
			EnableLogFileValidation:    aws.Bool(true),
			IncludeGlobalServiceEvents: aws.Bool(true),
			IsMultiRegionTrail:         aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.LogFileValidationEnabled == nil || !*resp.LogFileValidationEnabled {
			return fmt.Errorf("expected LogFileValidationEnabled=true")
		}
		if resp.IncludeGlobalServiceEvents == nil || !*resp.IncludeGlobalServiceEvents {
			return fmt.Errorf("expected IncludeGlobalServiceEvents=true")
		}
		if resp.IsMultiRegionTrail == nil || !*resp.IsMultiRegionTrail {
			return fmt.Errorf("expected IsMultiRegionTrail=true")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_ContentVerify", func() error {
		name := tc.uniqueName("verify-trail")
		defer tc.deleteTrail(name)

		if err := tc.ensureTrailBucket("verify-bucket"); err != nil {
			return err
		}
		resp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:                       aws.String(name),
			S3BucketName:               aws.String("verify-bucket"),
			IncludeGlobalServiceEvents: aws.Bool(true),
			IsMultiRegionTrail:         aws.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		if resp.Name == nil || *resp.Name != name {
			return fmt.Errorf("trail name mismatch")
		}
		if resp.S3BucketName == nil || *resp.S3BucketName != "verify-bucket" {
			return fmt.Errorf("S3 bucket name mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "UpdateTrail_VerifyChange", func() error {
		name := tc.uniqueName("upd-verify")
		defer tc.deleteTrail(name)

		_, err := tc.createTrail(name, "upd-verify-bucket")
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}

		if err := tc.ensureTrailBucket("updated-verify-bucket"); err != nil {
			return err
		}
		_, err = tc.client.UpdateTrail(tc.ctx, &cloudtrail.UpdateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("updated-verify-bucket"),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		resp, err := tc.client.GetTrail(tc.ctx, &cloudtrail.GetTrailInput{
			Name: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Trail == nil {
			return fmt.Errorf("GetTrail returned nil Trail")
		}
		if resp.Trail.S3BucketName == nil || *resp.Trail.S3BucketName != "updated-verify-bucket" {
			return fmt.Errorf("S3 bucket name not updated, got %v", resp.Trail.S3BucketName)
		}
		return nil
	}))

	// CreateTrail verifies the delivery destinations before accepting the
	// configuration: the bucket must exist, its policy must grant CloudTrail
	// the documented write access, and a configured SNS topic must resolve
	// to an existing, publish-permitted topic.
	results = append(results, r.RunTest("cloudtrail", "CreateTrail_MissingBucket", func() error {
		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(tc.uniqueName("ct-missing-bucket")),
			S3BucketName: aws.String(tc.uniqueName("ct-ghost-bucket")),
		})
		return AssertErrorContains(err, "S3BucketDoesNotExistException")
	}))

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_InsufficientBucketPolicy", func() error {
		// The bucket exists but carries no CloudTrail policy.
		bucket := tc.uniqueName("ct-unprivileged-bucket")
		if _, err := tc.s3Client.CreateBucket(tc.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer tc.s3Client.DeleteBucket(tc.ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(tc.uniqueName("ct-bad-policy")),
			S3BucketName: aws.String(bucket),
		})
		return AssertErrorContains(err, "InsufficientS3BucketPolicyException")
	}))

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_MissingSnsTopic", func() error {
		if err := tc.ensureTrailBucket("topic-validation-bucket"); err != nil {
			return err
		}
		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(tc.uniqueName("ct-missing-topic")),
			S3BucketName: aws.String("topic-validation-bucket"),
			SnsTopicName: aws.String(tc.uniqueName("ct-ghost-topic")),
		})
		return AssertErrorContains(err, "InsufficientSnsTopicPolicyException")
	}))

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_SnsTopicResolvesARN", func() error {
		if err := tc.ensureTrailBucket("topic-validation-bucket"); err != nil {
			return err
		}
		topic, err := tc.snsClient.CreateTopic(tc.ctx, &sns.CreateTopicInput{Name: aws.String(tc.uniqueName("ct-notify"))})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.snsClient.DeleteTopic(tc.ctx, &sns.DeleteTopicInput{TopicArn: topic.TopicArn})
		topicPolicy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Sid":"AWSCloudTrailSNSPolicy20131101","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"SNS:Publish","Resource":"%s"}]}`,
			aws.ToString(topic.TopicArn))
		if _, err := tc.snsClient.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn: topic.TopicArn, AttributeName: aws.String("Policy"), AttributeValue: aws.String(topicPolicy),
		}); err != nil {
			return fmt.Errorf("set topic policy: %v", err)
		}

		name := tc.uniqueName("ct-topic-trail")
		defer tc.deleteTrail(name)
		resp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("topic-validation-bucket"),
			SnsTopicName: topic.TopicArn,
		})
		if err != nil {
			return err
		}
		// The member carries a name or an ARN; the ARN form echoes verbatim
		// as the resolved destination.
		if aws.ToString(resp.SnsTopicARN) != aws.ToString(topic.TopicArn) {
			return fmt.Errorf("SnsTopicARN = %q, want %q", aws.ToString(resp.SnsTopicARN), aws.ToString(topic.TopicArn))
		}
		return nil
	}))

	return results
}
