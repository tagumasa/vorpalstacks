package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

func (r *TestRunner) runCloudTrailTrailTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	var trailName string
	results = append(results, r.RunTest("cloudtrail", "CreateTrail", func() error {
		trailName = fmt.Sprintf("test-trail-%d", time.Now().UnixNano())
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

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_DefaultFields", func() error {
		name := fmt.Sprintf("defaults-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		resp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("defaults-bucket"),
		})
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
		name := fmt.Sprintf("dup-trail-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("dup-bucket"),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}

		_, err = tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("dup-bucket"),
		})
		if err := AssertErrorContains(err, "TrailAlreadyExists"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrail_ByARN", func() error {
		name := fmt.Sprintf("arn-trail-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("arn-bucket"),
		})
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
		name := fmt.Sprintf("desc-arn-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("desc-arn-bucket"),
		})
		if err != nil {
			return err
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
		name := fmt.Sprintf("listall-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("listall-bucket"),
		})
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
		name := fmt.Sprintf("lfv-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("lfv-bucket"),
		})
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
		name := fmt.Sprintf("lfv-create-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

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
		name := fmt.Sprintf("verify-trail-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

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
		name := fmt.Sprintf("upd-verify-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("upd-verify-bucket"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
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
		if resp.Trail == nil || resp.Trail.S3BucketName == nil || *resp.Trail.S3BucketName != "updated-verify-bucket" {
			return fmt.Errorf("S3 bucket name not updated, got %v", resp.Trail.S3BucketName)
		}
		return nil
	}))

	return results
}
