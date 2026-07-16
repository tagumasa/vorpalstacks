package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTJobTests covers the Job lifecycle (Create with inline document/
// GetJobDocument round-trip/Describe/List/Cancel/Delete) plus the JobTemplate
// and OTAUpdate ecosystems. All ids/names use uniqueName.
func (r *TestRunner) runIoTJobTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	thingName := uniqueName("job-thing")
	jobID := uniqueName("job")
	doc := `{"version":1,"cmd":"reboot"}`
	target := tc.arn("iot", "thing", thingName)

	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingName)})

	results = append(results, r.RunTest("iot", "Job_CreateJob_WithDocument", func() error {
		if _, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String(thingName)}); err != nil {
			return fmt.Errorf("CreateThing prerequisite failed: %w", err)
		}
		if _, err := tc.client.CreateJob(tc.ctx, &iot.CreateJobInput{
			JobId:    aws.String(jobID),
			Document: aws.String(doc),
			Targets:  []string{target},
		}); err != nil {
			return fmt.Errorf("CreateJob failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Job_GetJobDocument_RoundTrip", func() error {
		out, err := tc.client.GetJobDocument(tc.ctx, &iot.GetJobDocumentInput{JobId: aws.String(jobID)})
		if err != nil {
			return fmt.Errorf("GetJobDocument failed: %w", err)
		}
		if out.Document == nil || *out.Document != doc {
			return fmt.Errorf("document mismatch: got %v", out.Document)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Job_DescribeJob_Targets", func() error {
		out, err := tc.client.DescribeJob(tc.ctx, &iot.DescribeJobInput{JobId: aws.String(jobID)})
		if err != nil {
			return fmt.Errorf("DescribeJob failed: %w", err)
		}
		if out.Job == nil {
			return fmt.Errorf("expected non-nil job")
		}
		if aws.ToString(out.Job.JobId) != jobID {
			return fmt.Errorf("expected jobId=%s, got %v", jobID, out.Job.JobId)
		}
		if len(out.Job.Targets) == 0 {
			return fmt.Errorf("expected non-empty targets")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Job_ListJobs_IncludesCreated", func() error {
		out, err := tc.client.ListJobs(tc.ctx, &iot.ListJobsInput{})
		if err != nil {
			return fmt.Errorf("ListJobs failed: %w", err)
		}
		for _, j := range out.Jobs {
			if j.JobId != nil && *j.JobId == jobID {
				return nil
			}
		}
		return fmt.Errorf("%s not found in list of %d jobs", jobID, len(out.Jobs))
	}))

	results = append(results, r.RunTest("iot", "Job_DescribeJob_NotFound", func() error {
		_, err := tc.client.DescribeJob(tc.ctx, &iot.DescribeJobInput{JobId: aws.String(uniqueName("nope-job"))})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "Job_DeleteJob", func() error {
		if _, err := tc.client.DeleteJob(tc.ctx, &iot.DeleteJobInput{JobId: aws.String(jobID)}); err != nil {
			return fmt.Errorf("DeleteJob failed: %w", err)
		}
		return nil
	}))

	// ── JobTemplate ecosystem ──
	tmplID := uniqueName("job-tmpl")
	defer tc.client.DeleteJobTemplate(tc.ctx, &iot.DeleteJobTemplateInput{JobTemplateId: aws.String(tmplID)})

	results = append(results, r.RunTest("iot", "JobTemplate_CreateJobTemplate", func() error {
		out, err := tc.client.CreateJobTemplate(tc.ctx, &iot.CreateJobTemplateInput{
			JobTemplateId: aws.String(tmplID),
			Document:      aws.String(`{"version":1}`),
			Description:   aws.String("test template"),
		})
		if err != nil {
			return fmt.Errorf("CreateJobTemplate failed: %w", err)
		}
		if out.JobTemplateId == nil || *out.JobTemplateId != tmplID {
			return fmt.Errorf("expected jobTemplateId=%s, got %v", tmplID, out.JobTemplateId)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "JobTemplate_DescribeJobTemplate", func() error {
		out, err := tc.client.DescribeJobTemplate(tc.ctx, &iot.DescribeJobTemplateInput{JobTemplateId: aws.String(tmplID)})
		if err != nil {
			return fmt.Errorf("DescribeJobTemplate failed: %w", err)
		}
		if aws.ToString(out.JobTemplateId) != tmplID {
			return fmt.Errorf("expected jobTemplateId=%s", tmplID)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "JobTemplate_ListJobTemplates_IncludesCreated", func() error {
		out, err := tc.client.ListJobTemplates(tc.ctx, &iot.ListJobTemplatesInput{})
		if err != nil {
			return fmt.Errorf("ListJobTemplates failed: %w", err)
		}
		for _, t := range out.JobTemplates {
			if t.JobTemplateId != nil && *t.JobTemplateId == tmplID {
				return nil
			}
		}
		return fmt.Errorf("%s not found in list of %d templates", tmplID, len(out.JobTemplates))
	}))

	// ── OTAUpdate ecosystem ──
	otaID := uniqueName("ota")

	results = append(results, r.RunTest("iot", "OTAUpdate_CreateOTAUpdate", func() error {
		out, err := tc.client.CreateOTAUpdate(tc.ctx, &iot.CreateOTAUpdateInput{
			OtaUpdateId: aws.String(otaID),
			RoleArn:     aws.String(tc.iamRoleARN("ota")),
			Targets:     []string{target},
			Files: []iottypes.OTAUpdateFile{{
				FileName:     aws.String("firmware.bin"),
				FileVersion:  aws.String("1.0"),
				FileLocation: &iottypes.FileLocation{},
			}},
		})
		if err != nil {
			return fmt.Errorf("CreateOTAUpdate failed: %w", err)
		}
		if out.OtaUpdateId == nil || *out.OtaUpdateId != otaID {
			return fmt.Errorf("expected otaUpdateId=%s, got %v", otaID, out.OtaUpdateId)
		}
		return nil
	}))

	defer tc.client.DeleteOTAUpdate(tc.ctx, &iot.DeleteOTAUpdateInput{OtaUpdateId: aws.String(otaID)})

	results = append(results, r.RunTest("iot", "OTAUpdate_GetOTAUpdate", func() error {
		out, err := tc.client.GetOTAUpdate(tc.ctx, &iot.GetOTAUpdateInput{OtaUpdateId: aws.String(otaID)})
		if err != nil {
			return fmt.Errorf("GetOTAUpdate failed: %w", err)
		}
		if out.OtaUpdateInfo == nil || aws.ToString(out.OtaUpdateInfo.OtaUpdateId) != otaID {
			return fmt.Errorf("expected otaUpdateInfo for %s", otaID)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "OTAUpdate_ListOTAUpdates_IncludesCreated", func() error {
		out, err := tc.client.ListOTAUpdates(tc.ctx, &iot.ListOTAUpdatesInput{})
		if err != nil {
			return fmt.Errorf("ListOTAUpdates failed: %w", err)
		}
		for _, o := range out.OtaUpdates {
			if o.OtaUpdateId != nil && *o.OtaUpdateId == otaID {
				return nil
			}
		}
		return fmt.Errorf("%s not found in list of %d ota updates", otaID, len(out.OtaUpdates))
	}))

	return results
}
