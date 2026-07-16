package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

func (r *TestRunner) runCloudTrailSelectorTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "GetEventSelectors", func() error {
		name := fmt.Sprintf("ges-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("ges-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.GetEventSelectors(tc.ctx, &cloudtrail.GetEventSelectorsInput{
			TrailName: aws.String(name),
		})
		if err != nil {
			return err
		}
		if len(resp.EventSelectors) < 1 {
			return fmt.Errorf("expected at least 1 default event selector, got %d", len(resp.EventSelectors))
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

	results = append(results, r.RunTest("cloudtrail", "PutEventSelectors", func() error {
		name := fmt.Sprintf("pes-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("pes-bucket"),
		})
		if err != nil {
			return err
		}

		putResp, err := tc.client.PutEventSelectors(tc.ctx, &cloudtrail.PutEventSelectorsInput{
			TrailName: aws.String(name),
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
		if len(putResp.EventSelectors) != 1 {
			return fmt.Errorf("expected 1 event selector in response, got %d", len(putResp.EventSelectors))
		}
		if putResp.EventSelectors[0].ReadWriteType != types.ReadWriteTypeAll {
			return fmt.Errorf("ReadWriteType mismatch, got %v", putResp.EventSelectors[0].ReadWriteType)
		}

		getResp, err := tc.client.GetEventSelectors(tc.ctx, &cloudtrail.GetEventSelectorsInput{TrailName: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get event selectors: %v", err)
		}
		if len(getResp.EventSelectors) != 1 || getResp.EventSelectors[0].ReadWriteType != types.ReadWriteTypeAll {
			return fmt.Errorf("event selectors not persisted after PutEventSelectors")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutEventSelectors_VerifyContent", func() error {
		name := fmt.Sprintf("vc-es-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("vc-es-bucket"),
		})
		if err != nil {
			return err
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
		_, err = tc.client.PutEventSelectors(tc.ctx, &cloudtrail.PutEventSelectorsInput{
			TrailName:      aws.String(name),
			EventSelectors: selectors,
		})
		if err != nil {
			return fmt.Errorf("put event selectors: %v", err)
		}
		resp, err := tc.client.GetEventSelectors(tc.ctx, &cloudtrail.GetEventSelectorsInput{
			TrailName: aws.String(name),
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

	results = append(results, r.RunTest("cloudtrail", "PutEventSelectors_ExcludeManagementEventSources", func() error {
		name := fmt.Sprintf("emes-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
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
		_, err = tc.client.PutEventSelectors(tc.ctx, &cloudtrail.PutEventSelectorsInput{
			TrailName:      aws.String(name),
			EventSelectors: selectors,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := tc.client.GetEventSelectors(tc.ctx, &cloudtrail.GetEventSelectorsInput{
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

	results = append(results, r.RunTest("cloudtrail", "GetEventSelectors_DefaultValues", func() error {
		name := fmt.Sprintf("def-es-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("def-es-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.GetEventSelectors(tc.ctx, &cloudtrail.GetEventSelectorsInput{
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

	results = append(results, r.RunTest("cloudtrail", "PutInsightSelectors", func() error {
		name := fmt.Sprintf("insight-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("insight-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.PutInsightSelectors(tc.ctx, &cloudtrail.PutInsightSelectorsInput{
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
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("get-insight-bucket"),
		})
		if err != nil {
			return err
		}

		_, err = tc.client.PutInsightSelectors(tc.ctx, &cloudtrail.PutInsightSelectorsInput{
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

		resp, err := tc.client.GetInsightSelectors(tc.ctx, &cloudtrail.GetInsightSelectorsInput{
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
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("empty-insight-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.GetInsightSelectors(tc.ctx, &cloudtrail.GetInsightSelectorsInput{
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

	return results
}
