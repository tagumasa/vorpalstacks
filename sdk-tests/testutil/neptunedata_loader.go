package testutil

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata/types"
)

func (r *TestRunner) runNeptunedataLoaderTests(tc *neptunedataContext) []TestResult {
	var results []TestResult
	var loaderJobID string

	results = append(results, r.RunTest("neptunedata", "StartLoaderJob", func() error {
		resp, err := tc.client.StartLoaderJob(tc.ctx, &neptunedata.StartLoaderJobInput{
			Source:         aws.String("s3://test-bucket/data"),
			Format:         types.FormatCsv,
			IamRoleArn:     aws.String("arn:aws:iam::000000000000:role/NeptuneLoadRole"),
			S3BucketRegion: types.S3BucketRegionUsEast1,
		})
		if err != nil {
			return err
		}
		if resp.Payload == nil {
			return fmt.Errorf("expected non-nil loader job payload")
		}
		data := marshalDoc(resp.Payload)
		var payloadMap map[string]interface{}
		if err := json.Unmarshal([]byte(data), &payloadMap); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		if lid, ok := payloadMap["loadId"].(string); ok && lid != "" {
			loaderJobID = lid
		} else {
			return fmt.Errorf("expected loadId in payload, got %s", data)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetLoaderJobStatus", func() error {
		if loaderJobID == "" {
			return fmt.Errorf("no loader job ID from StartLoaderJob")
		}
		resp, err := tc.client.GetLoaderJobStatus(tc.ctx, &neptunedata.GetLoaderJobStatusInput{
			LoadId: aws.String(loaderJobID),
		})
		if err != nil {
			return err
		}
		if resp.Payload == nil {
			return fmt.Errorf("expected non-nil loader job status payload")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ListLoaderJobs", func() error {
		resp, err := tc.client.ListLoaderJobs(tc.ctx, &neptunedata.ListLoaderJobsInput{})
		if err != nil {
			return err
		}
		if resp.Payload == nil {
			return fmt.Errorf("expected non-nil list loader jobs payload")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "CancelLoaderJob", func() error {
		if loaderJobID == "" {
			return fmt.Errorf("no loader job ID from StartLoaderJob")
		}
		resp, err := tc.client.CancelLoaderJob(tc.ctx, &neptunedata.CancelLoaderJobInput{
			LoadId: aws.String(loaderJobID),
		})
		if err != nil {
			return err
		}
		if resp.Status == nil || *resp.Status == "" {
			return fmt.Errorf("expected non-empty status from cancel")
		}
		return nil
	}))

	return results
}
