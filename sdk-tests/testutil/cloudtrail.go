package testutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunCloudTrailTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudtrail",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := cloudtrail.NewFromConfig(cfg)
	ctx := context.Background()

	trailName := fmt.Sprintf("test-trail-%d", time.Now().UnixNano())

	// === BASIC CRUD ===

	results = append(results, r.RunTest("cloudtrail", "ListTrails", func() error {
		resp, err := client.ListTrails(ctx, &cloudtrail.ListTrailsInput{})
		if err != nil {
			return err
		}
		if resp.Trails == nil {
			return fmt.Errorf("trails list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "CreateTrail", func() error {
		resp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:                       aws.String(trailName),
			S3BucketName:               aws.String("test-bucket"),
			IncludeGlobalServiceEvents: aws.Bool(true),
			IsMultiRegionTrail:         aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil {
			return fmt.Errorf("trail name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrail", func() error {
		resp, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{
			Name: aws.String(trailName),
		})
		if err != nil {
			return err
		}
		if resp.Trail == nil {
			return fmt.Errorf("trail is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DescribeTrails", func() error {
		resp, err := client.DescribeTrails(ctx, &cloudtrail.DescribeTrailsInput{
			TrailNameList: []string{trailName},
		})
		if err != nil {
			return err
		}
		if resp.TrailList == nil {
			return fmt.Errorf("trail list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "StartLogging", func() error {
		resp, err := client.StartLogging(ctx, &cloudtrail.StartLoggingInput{
			Name: aws.String(trailName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "StopLogging", func() error {
		resp, err := client.StopLogging(ctx, &cloudtrail.StopLoggingInput{
			Name: aws.String(trailName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus", func() error {
		resp, err := client.GetTrailStatus(ctx, &cloudtrail.GetTrailStatusInput{
			Name: aws.String(trailName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "UpdateTrail", func() error {
		resp, err := client.UpdateTrail(ctx, &cloudtrail.UpdateTrailInput{
			Name:         aws.String(trailName),
			S3BucketName: aws.String("updated-bucket"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetEventSelectors", func() error {
		resp, err := client.GetEventSelectors(ctx, &cloudtrail.GetEventSelectorsInput{
			TrailName: aws.String(trailName),
		})
		if err != nil {
			return err
		}
		if resp.EventSelectors == nil {
			return fmt.Errorf("event selectors list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutEventSelectors", func() error {
		resp, err := client.PutEventSelectors(ctx, &cloudtrail.PutEventSelectorsInput{
			TrailName: aws.String(trailName),
			EventSelectors: []types.EventSelector{
				{
					ReadWriteType:           types.ReadWriteTypeAll,
					IncludeManagementEvents: aws.Bool(true),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	var trailARN string
	getTrailResp, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{
		Name: aws.String(trailName),
	})
	if err == nil && getTrailResp.Trail != nil && getTrailResp.Trail.TrailARN != nil {
		trailARN = *getTrailResp.Trail.TrailARN
	}

	results = append(results, r.RunTest("cloudtrail", "AddTags", func() error {
		resourceID := trailARN
		if resourceID == "" {
			resourceID = trailName
		}
		resp, err := client.AddTags(ctx, &cloudtrail.AddTagsInput{
			ResourceId: aws.String(resourceID),
			TagsList: []types.Tag{
				{
					Key:   aws.String("Environment"),
					Value: aws.String("test"),
				},
				{
					Key:   aws.String("Owner"),
					Value: aws.String("test-user"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "ListTags", func() error {
		resourceID := trailARN
		if resourceID == "" {
			resourceID = trailName
		}
		resp, err := client.ListTags(ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{resourceID},
		})
		if err != nil {
			return err
		}
		if resp.ResourceTagList == nil {
			return fmt.Errorf("resource tag list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "RemoveTags", func() error {
		resourceID := trailARN
		if resourceID == "" {
			resourceID = trailName
		}
		resp, err := client.RemoveTags(ctx, &cloudtrail.RemoveTagsInput{
			ResourceId: aws.String(resourceID),
			TagsList: []types.Tag{
				{
					Key: aws.String("Environment"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents", func() error {
		resp, err := client.LookupEvents(ctx, &cloudtrail.LookupEventsInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.Events == nil {
			return fmt.Errorf("events list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DeleteTrail", func() error {
		resp, err := client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{
			Name: aws.String(trailName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("cloudtrail", "GetTrail_NonExistent", func() error {
		_, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{
			Name: aws.String("nonexistent-trail-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DeleteTrail_NonExistent", func() error {
		_, err := client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{
			Name: aws.String("nonexistent-trail-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "StartLogging_NonExistent", func() error {
		_, err := client.StartLogging(ctx, &cloudtrail.StartLoggingInput{
			Name: aws.String("nonexistent-trail-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DescribeTrails_NonExistent", func() error {
		resp, err := client.DescribeTrails(ctx, &cloudtrail.DescribeTrailsInput{
			TrailNameList: []string{"nonexistent-trail-xyz"},
		})
		if err != nil {
			return fmt.Errorf("DescribeTrails should not error for non-existent trail, got: %v", err)
		}
		if len(resp.TrailList) != 0 {
			return fmt.Errorf("expected empty trail list for non-existent trail, got %d", len(resp.TrailList))
		}
		return nil
	}))

	// === CONTENT VERIFICATION TESTS ===

	var verifyTrailName string
	results = append(results, r.RunTest("cloudtrail", "CreateTrail_ContentVerify", func() error {
		verifyTrailName = fmt.Sprintf("verify-trail-%d", time.Now().UnixNano())
		resp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:                       aws.String(verifyTrailName),
			S3BucketName:               aws.String("verify-bucket"),
			IncludeGlobalServiceEvents: aws.Bool(true),
			IsMultiRegionTrail:         aws.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		if resp.Name == nil || *resp.Name != verifyTrailName {
			return fmt.Errorf("trail name mismatch")
		}
		if resp.S3BucketName == nil || *resp.S3BucketName != "verify-bucket" {
			return fmt.Errorf("S3 bucket name mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "UpdateTrail_VerifyChange", func() error {
		if verifyTrailName == "" {
			return fmt.Errorf("trail name not available")
		}
		_, err := client.UpdateTrail(ctx, &cloudtrail.UpdateTrailInput{
			Name:         aws.String(verifyTrailName),
			S3BucketName: aws.String("updated-verify-bucket"),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		resp, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{
			Name: aws.String(verifyTrailName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Trail == nil || resp.Trail.S3BucketName == nil || *resp.Trail.S3BucketName != "updated-verify-bucket" {
			return fmt.Errorf("S3 bucket name not updated, got %v", resp.Trail.S3BucketName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutEventSelectors_VerifyContent", func() error {
		if verifyTrailName == "" {
			return fmt.Errorf("trail name not available")
		}
		selectors := []types.EventSelector{
			{
				ReadWriteType:           types.ReadWriteTypeReadOnly,
				IncludeManagementEvents: aws.Bool(false),
				DataResources: []types.DataResource{
					{
						Type:   aws.String("AWS::S3::Object"),
						Values: []string{"arn:aws:s3:::"},
					},
				},
			},
		}
		_, err := client.PutEventSelectors(ctx, &cloudtrail.PutEventSelectorsInput{
			TrailName:      aws.String(verifyTrailName),
			EventSelectors: selectors,
		})
		if err != nil {
			return fmt.Errorf("put event selectors: %v", err)
		}
		resp, err := client.GetEventSelectors(ctx, &cloudtrail.GetEventSelectorsInput{
			TrailName: aws.String(verifyTrailName),
		})
		if err != nil {
			return fmt.Errorf("get event selectors: %v", err)
		}
		if len(resp.EventSelectors) != 1 {
			return fmt.Errorf("expected 1 event selector, got %d", len(resp.EventSelectors))
		}
		if resp.EventSelectors[0].ReadWriteType != types.ReadWriteTypeReadOnly {
			return fmt.Errorf("ReadWriteType mismatch, got %v", resp.EventSelectors[0].ReadWriteType)
		}
		return nil
	}))

	if verifyTrailName != "" {
		client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(verifyTrailName)})
	}

	// === EXPANDED TESTS ===

	// --- Trail default fields verification ---

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_DefaultFields", func() error {
		name := fmt.Sprintf("defaults-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		resp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
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

	// --- Duplicate trail creation ---

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_Duplicate", func() error {
		name := fmt.Sprintf("dup-trail-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("dup-bucket"),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}

		_, err = client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("dup-bucket"),
		})
		if err := AssertErrorContains(err, "TrailAlreadyExists"); err != nil {
			return err
		}
		return nil
	}))

	// --- GetTrail by ARN ---

	results = append(results, r.RunTest("cloudtrail", "GetTrail_ByARN", func() error {
		name := fmt.Sprintf("arn-trail-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("arn-bucket"),
		})
		if err != nil {
			return err
		}
		if createResp.TrailARN == nil {
			return fmt.Errorf("trail ARN is nil")
		}

		getResp, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{
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

	// --- DescribeTrails by ARN ---

	results = append(results, r.RunTest("cloudtrail", "DescribeTrails_ByARN", func() error {
		name := fmt.Sprintf("desc-arn-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("desc-arn-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := client.DescribeTrails(ctx, &cloudtrail.DescribeTrailsInput{
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

	// --- DescribeTrails without filter (list all) ---

	results = append(results, r.RunTest("cloudtrail", "DescribeTrails_ListAll", func() error {
		resp, err := client.DescribeTrails(ctx, &cloudtrail.DescribeTrailsInput{})
		if err != nil {
			return err
		}
		if resp.TrailList == nil {
			return fmt.Errorf("trail list is nil")
		}
		return nil
	}))

	// --- Trail status after start/stop logging ---

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus_AfterStart", func() error {
		name := fmt.Sprintf("status-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("status-bucket"),
		})
		if err != nil {
			return err
		}

		_, err = client.StartLogging(ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("start logging: %v", err)
		}

		status, err := client.GetTrailStatus(ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get status: %v", err)
		}
		if status.IsLogging == nil || !*status.IsLogging {
			return fmt.Errorf("expected IsLogging=true after StartLogging")
		}
		if status.StartLoggingTime == nil {
			return fmt.Errorf("expected StartLoggingTime to be set")
		}
		if status.LatestDeliveryTime == nil {
			return fmt.Errorf("expected LatestDeliveryTime when logging is active")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus_AfterStop", func() error {
		name := fmt.Sprintf("stopstat-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("stopstat-bucket"),
		})
		if err != nil {
			return err
		}

		_, err = client.StartLogging(ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("start logging: %v", err)
		}
		_, err = client.StopLogging(ctx, &cloudtrail.StopLoggingInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("stop logging: %v", err)
		}

		status, err := client.GetTrailStatus(ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get status: %v", err)
		}
		if status.IsLogging != nil && *status.IsLogging {
			return fmt.Errorf("expected IsLogging=false after StopLogging")
		}
		if status.StopLoggingTime == nil {
			return fmt.Errorf("expected StopLoggingTime to be set")
		}
		return nil
	}))

	// --- UpdateTrail with LogFileValidationEnabled ---

	results = append(results, r.RunTest("cloudtrail", "UpdateTrail_EnableLogFileValidation", func() error {
		name := fmt.Sprintf("lfv-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("lfv-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := client.UpdateTrail(ctx, &cloudtrail.UpdateTrailInput{
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

	// --- CreateTrail with LogFileValidationEnabled (auto-generates public key) ---

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_WithLogFileValidation", func() error {
		name := fmt.Sprintf("lfv-create-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		resp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
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

	// --- ListPublicKeys (after creating trail with LogFileValidation) ---

	results = append(results, r.RunTest("cloudtrail", "ListPublicKeys", func() error {
		resp, err := client.ListPublicKeys(ctx, &cloudtrail.ListPublicKeysInput{})
		if err != nil {
			return fmt.Errorf("list public keys: %v", err)
		}
		if resp.PublicKeyList == nil {
			return fmt.Errorf("PublicKeyList is nil")
		}
		if len(resp.PublicKeyList) == 0 {
			return fmt.Errorf("expected at least 1 public key (from previous LogFileValidation trail)")
		}
		pk := resp.PublicKeyList[0]
		if pk.Fingerprint == nil || *pk.Fingerprint == "" {
			return fmt.Errorf("expected non-empty Fingerprint")
		}
		if len(pk.Value) == 0 {
			return fmt.Errorf("expected non-empty Value (DER bytes)")
		}
		_, err = base64.StdEncoding.DecodeString(base64.StdEncoding.EncodeToString(pk.Value))
		if err != nil {
			return fmt.Errorf("value should be base64-decodable DER bytes")
		}
		if pk.ValidityStartTime == nil || pk.ValidityStartTime.IsZero() {
			return fmt.Errorf("expected ValidityStartTime to be set")
		}
		if pk.ValidityEndTime == nil || pk.ValidityEndTime.IsZero() {
			return fmt.Errorf("expected ValidityEndTime to be set")
		}
		if pk.ValidityEndTime.Before(*pk.ValidityStartTime) {
			return fmt.Errorf("ValidityEndTime should be after ValidityStartTime")
		}
		return nil
	}))

	// --- ListPublicKeys with time filter ---

	results = append(results, r.RunTest("cloudtrail", "ListPublicKeys_TimeFilter", func() error {
		now := time.Now().UTC()
		past := now.Add(-1 * time.Hour)
		future := now.Add(1 * time.Hour)

		resp, err := client.ListPublicKeys(ctx, &cloudtrail.ListPublicKeysInput{
			StartTime: &past,
			EndTime:   &future,
		})
		if err != nil {
			return fmt.Errorf("list public keys with time filter: %v", err)
		}
		if resp.PublicKeyList == nil {
			return fmt.Errorf("PublicKeyList is nil")
		}
		if len(resp.PublicKeyList) == 0 {
			return fmt.Errorf("expected at least 1 public key in time range")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "ListPublicKeys_OutsideTimeRange", func() error {
		farFuture := time.Now().UTC().Add(10 * 24 * 365 * time.Hour)
		beyond := farFuture.Add(1 * time.Hour)

		resp, err := client.ListPublicKeys(ctx, &cloudtrail.ListPublicKeysInput{
			StartTime: &farFuture,
			EndTime:   &beyond,
		})
		if err != nil {
			return fmt.Errorf("list public keys outside range: %v", err)
		}
		if resp.PublicKeyList == nil {
			return fmt.Errorf("PublicKeyList is nil")
		}
		if len(resp.PublicKeyList) != 0 {
			return fmt.Errorf("expected 0 public keys outside validity range, got %d", len(resp.PublicKeyList))
		}
		return nil
	}))

	// --- EventSelectors with ExcludeManagementEventSources ---

	results = append(results, r.RunTest("cloudtrail", "PutEventSelectors_ExcludeManagementEventSources", func() error {
		name := fmt.Sprintf("emes-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("emes-bucket"),
		})
		if err != nil {
			return err
		}

		selectors := []types.EventSelector{
			{
				ReadWriteType:                 types.ReadWriteTypeAll,
				IncludeManagementEvents:       aws.Bool(true),
				ExcludeManagementEventSources: []string{"kms.amazonaws.com"},
			},
		}
		_, err = client.PutEventSelectors(ctx, &cloudtrail.PutEventSelectorsInput{
			TrailName:      aws.String(name),
			EventSelectors: selectors,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.GetEventSelectors(ctx, &cloudtrail.GetEventSelectorsInput{
			TrailName: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(resp.EventSelectors) != 1 {
			return fmt.Errorf("expected 1 selector, got %d", len(resp.EventSelectors))
		}
		if len(resp.EventSelectors[0].ExcludeManagementEventSources) != 1 {
			return fmt.Errorf("expected 1 ExcludeManagementEventSource, got %d", len(resp.EventSelectors[0].ExcludeManagementEventSources))
		}
		if resp.EventSelectors[0].ExcludeManagementEventSources[0] != "kms.amazonaws.com" {
			return fmt.Errorf("expected kms.amazonaws.com, got %s", resp.EventSelectors[0].ExcludeManagementEventSources[0])
		}
		return nil
	}))

	// --- GetEventSelectors default values ---

	results = append(results, r.RunTest("cloudtrail", "GetEventSelectors_DefaultValues", func() error {
		name := fmt.Sprintf("def-es-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("def-es-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := client.GetEventSelectors(ctx, &cloudtrail.GetEventSelectorsInput{
			TrailName: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(resp.EventSelectors) != 1 {
			return fmt.Errorf("expected 1 default event selector, got %d", len(resp.EventSelectors))
		}
		es := resp.EventSelectors[0]
		if es.ReadWriteType != types.ReadWriteTypeAll {
			return fmt.Errorf("expected default ReadWriteType=All, got %s", es.ReadWriteType)
		}
		if es.IncludeManagementEvents == nil || !*es.IncludeManagementEvents {
			return fmt.Errorf("expected default IncludeManagementEvents=true")
		}
		return nil
	}))

	// --- InsightSelectors ---

	results = append(results, r.RunTest("cloudtrail", "PutInsightSelectors", func() error {
		name := fmt.Sprintf("insight-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("insight-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := client.PutInsightSelectors(ctx, &cloudtrail.PutInsightSelectorsInput{
			TrailName: aws.String(name),
			InsightSelectors: []types.InsightSelector{
				{
					InsightType: types.InsightTypeApiCallRateInsight,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put insight selectors: %v", err)
		}
		if resp.InsightSelectors == nil || len(resp.InsightSelectors) != 1 {
			return fmt.Errorf("expected 1 insight selector in response")
		}
		if resp.InsightSelectors[0].InsightType != types.InsightTypeApiCallRateInsight {
			return fmt.Errorf("insight type mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetInsightSelectors", func() error {
		name := fmt.Sprintf("get-insight-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("get-insight-bucket"),
		})
		if err != nil {
			return err
		}

		_, err = client.PutInsightSelectors(ctx, &cloudtrail.PutInsightSelectorsInput{
			TrailName: aws.String(name),
			InsightSelectors: []types.InsightSelector{
				{
					InsightType: types.InsightTypeApiErrorRateInsight,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.GetInsightSelectors(ctx, &cloudtrail.GetInsightSelectorsInput{
			TrailName: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get insight selectors: %v", err)
		}
		if resp.InsightSelectors == nil || len(resp.InsightSelectors) != 1 {
			return fmt.Errorf("expected 1 insight selector, got %d", len(resp.InsightSelectors))
		}
		if resp.InsightSelectors[0].InsightType != types.InsightTypeApiErrorRateInsight {
			return fmt.Errorf("expected ApiErrorRateInsight, got %s", resp.InsightSelectors[0].InsightType)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetInsightSelectors_Empty", func() error {
		name := fmt.Sprintf("empty-insight-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("empty-insight-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := client.GetInsightSelectors(ctx, &cloudtrail.GetInsightSelectorsInput{
			TrailName: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(resp.InsightSelectors) != 0 {
			return fmt.Errorf("expected 0 insight selectors for new trail, got %d", len(resp.InsightSelectors))
		}
		return nil
	}))

	// --- Resource Policy ---

	results = append(results, r.RunTest("cloudtrail", "PutResourcePolicy_GetResourcePolicy", func() error {
		name := fmt.Sprintf("policy-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("policy-bucket"),
		})
		if err != nil {
			return err
		}

		policyDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"cloudtrail:GetTrail","Resource":"*"}]}`
		putResp, err := client.PutResourcePolicy(ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    createResp.TrailARN,
			ResourcePolicy: &policyDoc,
		})
		if err != nil {
			return fmt.Errorf("put resource policy: %v", err)
		}
		if putResp.ResourceArn == nil || *putResp.ResourceArn != *createResp.TrailARN {
			return fmt.Errorf("resource ARN mismatch in put response")
		}

		getResp, err := client.GetResourcePolicy(ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		if err != nil {
			return fmt.Errorf("get resource policy: %v", err)
		}
		if getResp.ResourcePolicy == nil || *getResp.ResourcePolicy != policyDoc {
			return fmt.Errorf("policy content mismatch, got: %v", getResp.ResourcePolicy)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetResourcePolicy_NotFound", func() error {
		name := fmt.Sprintf("nopolicy-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("nopolicy-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := client.GetResourcePolicy(ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		if err != nil {
			return fmt.Errorf("get resource policy (no policy): %v", err)
		}
		if resp.ResourcePolicy != nil && *resp.ResourcePolicy != "" {
			return fmt.Errorf("expected empty policy, got: %s", *resp.ResourcePolicy)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DeleteResourcePolicy", func() error {
		name := fmt.Sprintf("delpolicy-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("delpolicy-bucket"),
		})
		if err != nil {
			return err
		}

		policyDoc := `{"Version":"2012-10-17","Statement":[]}`
		_, err = client.PutResourcePolicy(ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    createResp.TrailARN,
			ResourcePolicy: &policyDoc,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		_, err = client.DeleteResourcePolicy(ctx, &cloudtrail.DeleteResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		if err != nil {
			return fmt.Errorf("delete resource policy: %v", err)
		}

		resp, err := client.GetResourcePolicy(ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		if err != nil {
			return fmt.Errorf("get after delete: %v", err)
		}
		if resp.ResourcePolicy != nil && *resp.ResourcePolicy != "" {
			return fmt.Errorf("expected empty policy after delete")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutResourcePolicy_NonExistentTrail", func() error {
		fakeARN := "arn:aws:cloudtrail:us-east-1:123456789012:trail/nonexistent-policy-trail"
		_, err := client.PutResourcePolicy(ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    aws.String(fakeARN),
			ResourcePolicy: aws.String(`{"Version":"2012-10-17"}`),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// --- Tag lifecycle ---

	results = append(results, r.RunTest("cloudtrail", "Tags_Lifecycle", func() error {
		name := fmt.Sprintf("tagcycle-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("tagcycle-bucket"),
		})
		if err != nil {
			return err
		}

		_, err = client.AddTags(ctx, &cloudtrail.AddTagsInput{
			ResourceId: createResp.TrailARN,
			TagsList: []types.Tag{
				{Key: aws.String("Key1"), Value: aws.String("Value1")},
				{Key: aws.String("Key2"), Value: aws.String("Value2")},
				{Key: aws.String("Key3"), Value: aws.String("Value3")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags: %v", err)
		}

		listResp, err := client.ListTags(ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*createResp.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(listResp.ResourceTagList) == 0 || len(listResp.ResourceTagList[0].TagsList) != 3 {
			return fmt.Errorf("expected 3 tags, got %d", len(listResp.ResourceTagList[0].TagsList))
		}

		_, err = client.RemoveTags(ctx, &cloudtrail.RemoveTagsInput{
			ResourceId: createResp.TrailARN,
			TagsList: []types.Tag{
				{Key: aws.String("Key2")},
			},
		})
		if err != nil {
			return fmt.Errorf("remove tags: %v", err)
		}

		listResp2, err := client.ListTags(ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*createResp.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("list tags after remove: %v", err)
		}
		if len(listResp2.ResourceTagList[0].TagsList) != 2 {
			return fmt.Errorf("expected 2 tags after remove, got %d", len(listResp2.ResourceTagList[0].TagsList))
		}

		_, err = client.AddTags(ctx, &cloudtrail.AddTagsInput{
			ResourceId: createResp.TrailARN,
			TagsList: []types.Tag{
				{Key: aws.String("Key1"), Value: aws.String("UpdatedValue1")},
			},
		})
		if err != nil {
			return fmt.Errorf("update tag: %v", err)
		}

		listResp3, err := client.ListTags(ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*createResp.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("list tags after update: %v", err)
		}
		found := false
		for _, t := range listResp3.ResourceTagList[0].TagsList {
			if t.Key != nil && *t.Key == "Key1" && t.Value != nil && *t.Value == "UpdatedValue1" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected Key1=UpdatedValue1 after update")
		}
		return nil
	}))

	// --- CreateTrail with tags ---

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_WithTags", func() error {
		name := fmt.Sprintf("tagtrail-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("tagtrail-bucket"),
			TagsList: []types.Tag{
				{Key: aws.String("CreatedBy"), Value: aws.String("sdk-test")},
				{Key: aws.String("Env"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return fmt.Errorf("create with tags: %v", err)
		}

		resp, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get trail: %v", err)
		}

		listResp, err := client.ListTags(ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*resp.Trail.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(listResp.ResourceTagList[0].TagsList) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(listResp.ResourceTagList[0].TagsList))
		}
		return nil
	}))

	// --- LookupEvents with time range ---

	results = append(results, r.RunTest("cloudtrail", "LookupEvents_WithTimeRange", func() error {
		now := time.Now().UTC()
		pastHour := now.Add(-1 * time.Hour)

		resp, err := client.LookupEvents(ctx, &cloudtrail.LookupEventsInput{
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
		return nil
	}))

	// --- Additional error cases ---

	results = append(results, r.RunTest("cloudtrail", "StopLogging_NonExistent", func() error {
		_, err := client.StopLogging(ctx, &cloudtrail.StopLoggingInput{
			Name: aws.String("nonexistent-stop-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus_NonExistent", func() error {
		_, err := client.GetTrailStatus(ctx, &cloudtrail.GetTrailStatusInput{
			Name: aws.String("nonexistent-status-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "UpdateTrail_NonExistent", func() error {
		_, err := client.UpdateTrail(ctx, &cloudtrail.UpdateTrailInput{
			Name:         aws.String("nonexistent-update-xyz"),
			S3BucketName: aws.String("some-bucket"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetEventSelectors_NonExistent", func() error {
		_, err := client.GetEventSelectors(ctx, &cloudtrail.GetEventSelectorsInput{
			TrailName: aws.String("nonexistent-es-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutEventSelectors_NonExistent", func() error {
		_, err := client.PutEventSelectors(ctx, &cloudtrail.PutEventSelectorsInput{
			TrailName: aws.String("nonexistent-es-xyz"),
			EventSelectors: []types.EventSelector{
				{ReadWriteType: types.ReadWriteTypeAll},
			},
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutInsightSelectors_NonExistent", func() error {
		_, err := client.PutInsightSelectors(ctx, &cloudtrail.PutInsightSelectorsInput{
			TrailName: aws.String("nonexistent-is-xyz"),
			InsightSelectors: []types.InsightSelector{
				{InsightType: types.InsightTypeApiCallRateInsight},
			},
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetInsightSelectors_NonExistent", func() error {
		_, err := client.GetInsightSelectors(ctx, &cloudtrail.GetInsightSelectorsInput{
			TrailName: aws.String("nonexistent-is-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "AddTags_NonExistent", func() error {
		_, err := client.AddTags(ctx, &cloudtrail.AddTagsInput{
			ResourceId: aws.String("arn:aws:cloudtrail:us-east-1:123456789012:trail/nonexistent-tag-xyz"),
			TagsList:   []types.Tag{{Key: aws.String("K"), Value: aws.String("V")}},
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "RemoveTags_NonExistent", func() error {
		_, err := client.RemoveTags(ctx, &cloudtrail.RemoveTagsInput{
			ResourceId: aws.String("arn:aws:cloudtrail:us-east-1:123456789012:trail/nonexistent-rm-xyz"),
			TagsList:   []types.Tag{{Key: aws.String("K")}},
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "ListTags_NonExistent", func() error {
		_, err := client.ListTags(ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{"arn:aws:cloudtrail:us-east-1:123456789012:trail/nonexistent-lt-xyz"},
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetResourcePolicy_NonExistentTrail", func() error {
		fakeARN := "arn:aws:cloudtrail:us-east-1:123456789012:trail/nonexistent-grp-xyz"
		_, err := client.GetResourcePolicy(ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: aws.String(fakeARN),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// --- formatTrail IsLogging field verification (via GetTrailStatus) ---

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus_DefaultIsNotLogging", func() error {
		name := fmt.Sprintf("islog-%d", time.Now().UnixNano())
		defer client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("islog-bucket"),
		})
		if err != nil {
			return err
		}

		status, err := client.GetTrailStatus(ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get status: %v", err)
		}
		if status.IsLogging != nil && *status.IsLogging {
			return fmt.Errorf("expected IsLogging=false for newly created trail")
		}
		if status.LatestDeliveryTime != nil {
			return fmt.Errorf("expected nil LatestDeliveryTime when not logging")
		}
		return nil
	}))

	return results
}
