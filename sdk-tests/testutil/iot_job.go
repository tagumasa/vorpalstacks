package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

func (r *TestRunner) runIoTJobTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "Job_CreateJob_WithDocument", func() error {
		_, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{
			ThingName: aws.String("test-job-thing"),
		})
		if err != nil {
			return fmt.Errorf("CreateThing prerequisite failed: %w", err)
		}
		_, err = tc.client.CreateJob(tc.ctx, &iot.CreateJobInput{
			JobId:     aws.String("test-job-doc-001"),
			Document:  aws.String(`{"version":1,"cmd":"reboot"}`),
			Targets:   []string{"arn:aws:iot:us-east-1:000000000000:thing/test-job-thing"},
		})
		if err != nil {
			return fmt.Errorf("CreateJob failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Job_GetJobDocument_RoundTrip", func() error {
		out, err := tc.client.GetJobDocument(tc.ctx, &iot.GetJobDocumentInput{
			JobId: aws.String("test-job-doc-001"),
		})
		if err != nil {
			return fmt.Errorf("GetJobDocument failed: %w", err)
		}
		if out.Document == nil {
			return fmt.Errorf("expected non-nil document")
		}
		if *out.Document != `{"version":1,"cmd":"reboot"}` {
			return fmt.Errorf("document mismatch: got %q", *out.Document)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Job_DescribeJob_Targets", func() error {
		out, err := tc.client.DescribeJob(tc.ctx, &iot.DescribeJobInput{
			JobId: aws.String("test-job-doc-001"),
		})
		if err != nil {
			return fmt.Errorf("DescribeJob failed: %w", err)
		}
		if out.Job == nil {
			return fmt.Errorf("expected non-nil job")
		}
		if len(out.Job.Targets) == 0 {
			return fmt.Errorf("expected non-empty targets")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Job_ListJobs", func() error {
		out, err := tc.client.ListJobs(tc.ctx, &iot.ListJobsInput{})
		if err != nil {
			return fmt.Errorf("ListJobs failed: %w", err)
		}
		found := false
		for _, j := range out.Jobs {
			if j.JobId != nil && *j.JobId == "test-job-doc-001" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("test-job-doc-001 not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Job_DeleteJob", func() error {
		_, err := tc.client.DeleteJob(tc.ctx, &iot.DeleteJobInput{
			JobId: aws.String("test-job-doc-001"),
		})
		if err != nil {
			return fmt.Errorf("DeleteJob failed: %w", err)
		}
		_, err = tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{
			ThingName: aws.String("test-job-thing"),
		})
		return nil
	}))

	return results
}
