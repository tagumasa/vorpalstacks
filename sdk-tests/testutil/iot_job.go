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

	results = append(results, r.RunTest("iot", "Job_ListJobExecutions_MemberShapes", func() error {
		forThing, err := tc.client.ListJobExecutionsForThing(tc.ctx, &iot.ListJobExecutionsForThingInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("ListJobExecutionsForThing failed: %w", err)
		}
		found := false
		for _, sum := range forThing.ExecutionSummaries {
			if sum.JobId != nil && *sum.JobId == jobID {
				found = true
				if sum.JobExecutionSummary == nil || sum.JobExecutionSummary.Status != iottypes.JobExecutionStatusQueued {
					return fmt.Errorf("expected QUEUED summary status for %s", jobID)
				}
			}
		}
		if !found {
			return fmt.Errorf("expected execution summary with jobId=%s for thing %s", jobID, thingName)
		}

		forJob, err := tc.client.ListJobExecutionsForJob(tc.ctx, &iot.ListJobExecutionsForJobInput{JobId: aws.String(jobID)})
		if err != nil {
			return fmt.Errorf("ListJobExecutionsForJob failed: %w", err)
		}
		found = false
		for _, sum := range forJob.ExecutionSummaries {
			if sum.ThingArn != nil && *sum.ThingArn == target {
				found = true
				if sum.JobExecutionSummary == nil || sum.JobExecutionSummary.Status != iottypes.JobExecutionStatusQueued {
					return fmt.Errorf("expected QUEUED summary status for %s", target)
				}
			}
		}
		if !found {
			return fmt.Errorf("expected execution summary with thingArn=%s", target)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Job_CancelJob_RecordsReasonCodeAndComment", func() error {
		out, err := tc.client.CancelJob(tc.ctx, &iot.CancelJobInput{
			JobId:      aws.String(jobID),
			ReasonCode: aws.String("DEV_TEST"),
			Comment:    aws.String("cancelled by sdk test"),
		})
		if err != nil {
			return fmt.Errorf("CancelJob failed: %w", err)
		}
		if aws.ToString(out.JobId) != jobID {
			return fmt.Errorf("expected jobId=%s, got %v", jobID, out.JobId)
		}
		if aws.ToString(out.JobArn) == "" {
			return fmt.Errorf("expected non-empty jobArn in cancel response")
		}
		desc, err := tc.client.DescribeJob(tc.ctx, &iot.DescribeJobInput{JobId: aws.String(jobID)})
		if err != nil {
			return fmt.Errorf("DescribeJob after cancel failed: %w", err)
		}
		if desc.Job == nil {
			return fmt.Errorf("expected job in describe response")
		}
		if aws.ToString(desc.Job.ReasonCode) != "DEV_TEST" {
			return fmt.Errorf("expected reasonCode=DEV_TEST, got %v", desc.Job.ReasonCode)
		}
		if aws.ToString(desc.Job.Comment) != "cancelled by sdk test" {
			return fmt.Errorf("expected comment persisted, got %v", desc.Job.Comment)
		}
		if desc.Job.Status != iottypes.JobStatusCanceled {
			return fmt.Errorf("expected status=CANCELED, got %v", desc.Job.Status)
		}
		return nil
	}))

	// ── List filters and per-execution member contracts ──
	contJobID := uniqueName("job-cont")
	groupName := uniqueName("job-group")
	groupJobID := uniqueName("job-grp")

	defer func() {
		_, _ = tc.client.DeleteJob(tc.ctx, &iot.DeleteJobInput{JobId: aws.String(contJobID)})
		_, _ = tc.client.DeleteJob(tc.ctx, &iot.DeleteJobInput{JobId: aws.String(groupJobID)})
	}()

	results = append(results, r.RunTest("iot", "Job_ListJobs_TargetSelectionFilter", func() error {
		if _, err := tc.client.CreateJob(tc.ctx, &iot.CreateJobInput{
			JobId:           aws.String(contJobID),
			Document:        aws.String(doc),
			Targets:         []string{target},
			TargetSelection: iottypes.TargetSelectionContinuous,
		}); err != nil {
			return fmt.Errorf("CreateJob CONTINUOUS prerequisite failed: %w", err)
		}
		continuous, err := tc.client.ListJobs(tc.ctx, &iot.ListJobsInput{
			TargetSelection: iottypes.TargetSelectionContinuous,
		})
		if err != nil {
			return fmt.Errorf("ListJobs CONTINUOUS failed: %w", err)
		}
		hasCont, leakedSnap := false, false
		for _, j := range continuous.Jobs {
			if aws.ToString(j.JobId) == contJobID {
				hasCont = true
			}
			if aws.ToString(j.JobId) == jobID {
				leakedSnap = true
			}
		}
		if !hasCont || leakedSnap {
			return fmt.Errorf("CONTINUOUS filter: hasCont=%v leakedSnapshot=%v", hasCont, leakedSnap)
		}
		// An omitted targetSelection documents as SNAPSHOT, so the plain
		// job must appear under the SNAPSHOT filter and the CONTINUOUS one
		// must not.
		snapshot, err := tc.client.ListJobs(tc.ctx, &iot.ListJobsInput{
			TargetSelection: iottypes.TargetSelectionSnapshot,
		})
		if err != nil {
			return fmt.Errorf("ListJobs SNAPSHOT failed: %w", err)
		}
		hasSnap, leakedCont := false, false
		for _, j := range snapshot.Jobs {
			if aws.ToString(j.JobId) == jobID {
				hasSnap = true
			}
			if aws.ToString(j.JobId) == contJobID {
				leakedCont = true
			}
		}
		if !hasSnap || leakedCont {
			return fmt.Errorf("SNAPSHOT filter: hasSnapshot=%v leakedContinuous=%v", hasSnap, leakedCont)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Job_ListJobs_ThingGroupFilter", func() error {
		cleanupGroup, err := tc.createThingGroup(groupName)
		if err != nil {
			return fmt.Errorf("CreateThingGroup prerequisite failed: %w", err)
		}
		defer cleanupGroup()
		if _, err := tc.client.CreateJob(tc.ctx, &iot.CreateJobInput{
			JobId:    aws.String(groupJobID),
			Document: aws.String(doc),
			Targets:  []string{tc.arn("iot", "thinggroup", groupName)},
		}); err != nil {
			return fmt.Errorf("CreateJob group-target prerequisite failed: %w", err)
		}
		out, err := tc.client.ListJobs(tc.ctx, &iot.ListJobsInput{
			ThingGroupName: aws.String(groupName),
		})
		if err != nil {
			return fmt.Errorf("ListJobs thingGroupName failed: %w", err)
		}
		found, leaked := false, false
		for _, j := range out.Jobs {
			if aws.ToString(j.JobId) == groupJobID {
				found = true
			}
			if aws.ToString(j.JobId) == jobID {
				leaked = true
			}
		}
		if !found || leaked {
			return fmt.Errorf("thing-group filter: found=%v leakedThingTarget=%v", found, leaked)
		}
		unknown, err := tc.client.ListJobs(tc.ctx, &iot.ListJobsInput{
			ThingGroupName: aws.String(uniqueName("no-such-group")),
		})
		if err != nil {
			return fmt.Errorf("ListJobs with an unknown group failed: %w", err)
		}
		if len(unknown.Jobs) != 0 {
			return fmt.Errorf("unknown group returned %d jobs", len(unknown.Jobs))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Job_CancelJobExecution_ExpectedVersionAndStatusDetails", func() error {
		// An explicitly supplied zero never matches the stored version
		// (versions start at one), so it is a version conflict rather
		// than an alias for an omitted member.
		_, err := tc.client.CancelJobExecution(tc.ctx, &iot.CancelJobExecutionInput{
			JobId:           aws.String(jobID),
			ThingName:       aws.String(thingName),
			ExpectedVersion: aws.Int64(0),
		})
		if codeErr := expectAWSErrorCode(err, "VersionConflictException"); codeErr != nil {
			return fmt.Errorf("explicit zero expectedVersion: %v", codeErr)
		}
		_, err = tc.client.CancelJobExecution(tc.ctx, &iot.CancelJobExecutionInput{
			JobId:           aws.String(jobID),
			ThingName:       aws.String(thingName),
			ExpectedVersion: aws.Int64(999),
		})
		if err == nil {
			return fmt.Errorf("expected VersionConflictException for a mismatched expectedVersion")
		}
		if codeErr := expectAWSErrorCode(err, "VersionConflictException"); codeErr != nil {
			return codeErr
		}
		if _, err := tc.client.CancelJobExecution(tc.ctx, &iot.CancelJobExecutionInput{
			JobId:           aws.String(jobID),
			ThingName:       aws.String(thingName),
			ExpectedVersion: aws.Int64(1),
			StatusDetails:   map[string]string{"progress": "halted"},
		}); err != nil {
			return fmt.Errorf("CancelJobExecution failed: %w", err)
		}
		out, err := tc.client.DescribeJobExecution(tc.ctx, &iot.DescribeJobExecutionInput{
			JobId:     aws.String(jobID),
			ThingName: aws.String(thingName),
		})
		if err != nil {
			return fmt.Errorf("DescribeJobExecution failed: %w", err)
		}
		if out.Execution == nil || out.Execution.Status != iottypes.JobExecutionStatusCanceled {
			return fmt.Errorf("expected CANCELED execution, got %+v", out.Execution)
		}
		if out.Execution.StatusDetails == nil || out.Execution.StatusDetails.DetailsMap["progress"] != "halted" {
			return fmt.Errorf("expected statusDetails round-trip, got %+v", out.Execution.StatusDetails)
		}
		return nil
	}))

	// The execution above is now CANCELED; only QUEUED (or IN_PROGRESS
	// with force) executions can be canceled, so a second cancel is an
	// invalid state transition.
	results = append(results, r.RunTest("iot", "Job_CancelJobExecution_TerminalStateRejected", func() error {
		_, err := tc.client.CancelJobExecution(tc.ctx, &iot.CancelJobExecutionInput{
			JobId:     aws.String(jobID),
			ThingName: aws.String(thingName),
		})
		if err == nil {
			return fmt.Errorf("expected canceling a terminal execution to be rejected")
		}
		return expectAWSErrorCode(err, "InvalidStateTransitionException")
	}))

	results = append(results, r.RunTest("iot", "Job_DeleteJobExecution_WrongNumberRejected", func() error {
		_, err := tc.client.DeleteJobExecution(tc.ctx, &iot.DeleteJobExecutionInput{
			JobId:           aws.String(jobID),
			ThingName:       aws.String(thingName),
			ExecutionNumber: aws.Int64(42),
		})
		if err == nil {
			return fmt.Errorf("expected ResourceNotFoundException for a wrong executionNumber")
		}
		if codeErr := expectAWSErrorCode(err, "ResourceNotFoundException"); codeErr != nil {
			return codeErr
		}
		if _, err := tc.client.DeleteJobExecution(tc.ctx, &iot.DeleteJobExecutionInput{
			JobId:           aws.String(jobID),
			ThingName:       aws.String(thingName),
			ExecutionNumber: aws.Int64(1),
		}); err != nil {
			return fmt.Errorf("DeleteJobExecution with the stored number failed: %w", err)
		}
		return nil
	}))

	// A QUEUED execution is not in a terminal state, so deleting it
	// without force is an invalid state transition; force deletes it.
	results = append(results, r.RunTest("iot", "Job_DeleteJobExecution_ForceGate", func() error {
		forceThing := uniqueName("job-force-thing")
		forceJob := uniqueName("job-force")
		forceTarget := tc.arn("iot", "thing", forceThing)
		defer tc.client.DeleteJob(tc.ctx, &iot.DeleteJobInput{JobId: aws.String(forceJob)})
		defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(forceThing)})
		if _, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String(forceThing)}); err != nil {
			return fmt.Errorf("CreateThing prerequisite failed: %w", err)
		}
		if _, err := tc.client.CreateJob(tc.ctx, &iot.CreateJobInput{
			JobId:    aws.String(forceJob),
			Document: aws.String(doc),
			Targets:  []string{forceTarget},
		}); err != nil {
			return fmt.Errorf("CreateJob prerequisite failed: %w", err)
		}
		_, err := tc.client.DeleteJobExecution(tc.ctx, &iot.DeleteJobExecutionInput{
			JobId:           aws.String(forceJob),
			ThingName:       aws.String(forceThing),
			ExecutionNumber: aws.Int64(1),
		})
		if err == nil {
			return fmt.Errorf("expected deleting a QUEUED execution without force to be rejected")
		}
		if codeErr := expectAWSErrorCode(err, "InvalidStateTransitionException"); codeErr != nil {
			return codeErr
		}
		if _, err := tc.client.DeleteJobExecution(tc.ctx, &iot.DeleteJobExecutionInput{
			JobId:           aws.String(forceJob),
			ThingName:       aws.String(forceThing),
			ExecutionNumber: aws.Int64(1),
			Force:           true,
		}); err != nil {
			return fmt.Errorf("DeleteJobExecution with force failed: %w", err)
		}
		return nil
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
			PresignedUrlConfig: &iottypes.PresignedUrlConfig{
				RoleArn:      aws.String(tc.iamRoleARN("tmpl")),
				ExpiresInSec: aws.Int64(3600),
			},
			JobExecutionsRolloutConfig: &iottypes.JobExecutionsRolloutConfig{
				MaximumPerMinute: aws.Int32(100),
			},
			TimeoutConfig: &iottypes.TimeoutConfig{
				InProgressTimeoutInMinutes: aws.Int64(30),
			},
			JobExecutionsRetryConfig: &iottypes.JobExecutionsRetryConfig{
				CriteriaList: []iottypes.RetryCriteria{{
					FailureType:     iottypes.RetryableFailureTypeFailed,
					NumberOfRetries: aws.Int32(3),
				}},
			},
			MaintenanceWindows: []iottypes.MaintenanceWindow{{
				StartTime:         aws.String("Mon 00:00"),
				DurationInMinutes: aws.Int32(60),
			}},
			DestinationPackageVersions: []string{"test-pkg@1.0"},
		})
		if err != nil {
			return fmt.Errorf("CreateJobTemplate failed: %w", err)
		}
		if out.JobTemplateId == nil || *out.JobTemplateId != tmplID {
			return fmt.Errorf("expected jobTemplateId=%s, got %v", tmplID, out.JobTemplateId)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "JobTemplate_MissingDocumentRejected", func() error {
		bareID := uniqueName("tmpl-bare")
		defer tc.client.DeleteJobTemplate(tc.ctx, &iot.DeleteJobTemplateInput{JobTemplateId: aws.String(bareID)})
		// The job document is required unless documentSource or jobArn is
		// supplied, so a request carrying none of the three is rejected.
		_, err := tc.client.CreateJobTemplate(tc.ctx, &iot.CreateJobTemplateInput{
			JobTemplateId: aws.String(bareID),
			Description:   aws.String("template without any document carrier"),
		})
		if err == nil {
			return fmt.Errorf("expected CreateJobTemplate without document, documentSource and jobArn to be rejected")
		}
		return expectAWSErrorCode(err, "InvalidRequestException")
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

	results = append(results, r.RunTest("iot", "JobTemplate_DescribeJobTemplate_MembersRoundTrip", func() error {
		out, err := tc.client.DescribeJobTemplate(tc.ctx, &iot.DescribeJobTemplateInput{JobTemplateId: aws.String(tmplID)})
		if err != nil {
			return fmt.Errorf("DescribeJobTemplate failed: %w", err)
		}
		if aws.ToString(out.JobTemplateArn) == "" {
			return fmt.Errorf("expected non-empty jobTemplateArn")
		}
		if out.PresignedUrlConfig == nil ||
			aws.ToString(out.PresignedUrlConfig.RoleArn) == "" ||
			aws.ToInt64(out.PresignedUrlConfig.ExpiresInSec) != 3600 {
			return fmt.Errorf("expected presignedUrlConfig round-trip, got %+v", out.PresignedUrlConfig)
		}
		if out.JobExecutionsRolloutConfig == nil ||
			aws.ToInt32(out.JobExecutionsRolloutConfig.MaximumPerMinute) != 100 {
			return fmt.Errorf("expected jobExecutionsRolloutConfig round-trip, got %+v", out.JobExecutionsRolloutConfig)
		}
		if out.TimeoutConfig == nil ||
			aws.ToInt64(out.TimeoutConfig.InProgressTimeoutInMinutes) != 30 {
			return fmt.Errorf("expected timeoutConfig round-trip, got %+v", out.TimeoutConfig)
		}
		if out.JobExecutionsRetryConfig == nil ||
			len(out.JobExecutionsRetryConfig.CriteriaList) != 1 ||
			aws.ToInt32(out.JobExecutionsRetryConfig.CriteriaList[0].NumberOfRetries) != 3 {
			return fmt.Errorf("expected jobExecutionsRetryConfig round-trip, got %+v", out.JobExecutionsRetryConfig)
		}
		if len(out.MaintenanceWindows) != 1 ||
			aws.ToInt32(out.MaintenanceWindows[0].DurationInMinutes) != 60 {
			return fmt.Errorf("expected maintenanceWindows round-trip, got %+v", out.MaintenanceWindows)
		}
		if len(out.DestinationPackageVersions) != 1 ||
			out.DestinationPackageVersions[0] != "test-pkg@1.0" {
			return fmt.Errorf("expected destinationPackageVersions round-trip, got %v", out.DestinationPackageVersions)
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
			AwsJobExecutionsRolloutConfig: &iottypes.AwsJobExecutionsRolloutConfig{
				MaximumPerMinute: aws.Int32(200),
			},
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
		if out.OtaUpdateStatus != iottypes.OTAUpdateStatusCreateComplete {
			return fmt.Errorf("expected otaUpdateStatus=CREATE_COMPLETE, got %v", out.OtaUpdateStatus)
		}
		return nil
	}))

	defer tc.client.DeleteOTAUpdate(tc.ctx, &iot.DeleteOTAUpdateInput{OtaUpdateId: aws.String(otaID), ForceDeleteAWSJob: true})

	// The awsIotJobId the create mints identifies a real IoT job; the OTA
	// job is not in a terminal state, so deleting the update without
	// forceDeleteAWSJob is rejected.
	results = append(results, r.RunTest("iot", "OTAUpdate_AwsJobMaterialised", func() error {
		out, err := tc.client.GetOTAUpdate(tc.ctx, &iot.GetOTAUpdateInput{OtaUpdateId: aws.String(otaID)})
		if err != nil {
			return fmt.Errorf("GetOTAUpdate failed: %w", err)
		}
		if out.OtaUpdateInfo == nil {
			return fmt.Errorf("expected otaUpdateInfo for %s", otaID)
		}
		awsJobID := aws.ToString(out.OtaUpdateInfo.AwsIotJobId)
		if awsJobID == "" {
			return fmt.Errorf("expected awsIotJobId in otaUpdateInfo")
		}
		job, err := tc.client.DescribeJob(tc.ctx, &iot.DescribeJobInput{JobId: aws.String(awsJobID)})
		if err != nil {
			return fmt.Errorf("DescribeJob of the OTA job failed: %w", err)
		}
		if job.Job == nil || len(job.Job.Targets) == 0 {
			return fmt.Errorf("expected the OTA job to carry targets, got %+v", job.Job)
		}
		_, err = tc.client.DeleteOTAUpdate(tc.ctx, &iot.DeleteOTAUpdateInput{OtaUpdateId: aws.String(otaID)})
		if codeErr := expectAWSErrorCode(err, "InvalidRequestException"); codeErr != nil {
			return fmt.Errorf("delete without forceDeleteAWSJob: %v", codeErr)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "OTAUpdate_GetOTAUpdate", func() error {
		out, err := tc.client.GetOTAUpdate(tc.ctx, &iot.GetOTAUpdateInput{OtaUpdateId: aws.String(otaID)})
		if err != nil {
			return fmt.Errorf("GetOTAUpdate failed: %w", err)
		}
		if out.OtaUpdateInfo == nil || aws.ToString(out.OtaUpdateInfo.OtaUpdateId) != otaID {
			return fmt.Errorf("expected otaUpdateInfo for %s", otaID)
		}
		info := out.OtaUpdateInfo
		if info.OtaUpdateStatus != iottypes.OTAUpdateStatusCreateComplete {
			return fmt.Errorf("expected otaUpdateStatus=CREATE_COMPLETE, got %v", info.OtaUpdateStatus)
		}
		if info.CreationDate == nil {
			return fmt.Errorf("expected non-nil creationDate in otaUpdateInfo")
		}
		if info.LastModifiedDate == nil {
			return fmt.Errorf("expected non-nil lastModifiedDate in otaUpdateInfo")
		}
		if aws.ToString(info.OtaUpdateArn) == "" {
			return fmt.Errorf("expected non-empty otaUpdateArn in otaUpdateInfo")
		}
		if len(info.OtaUpdateFiles) != 1 || aws.ToString(info.OtaUpdateFiles[0].FileName) != "firmware.bin" {
			return fmt.Errorf("expected otaUpdateFiles round-trip, got %+v", info.OtaUpdateFiles)
		}
		if info.AwsJobExecutionsRolloutConfig == nil ||
			aws.ToInt32(info.AwsJobExecutionsRolloutConfig.MaximumPerMinute) != 200 {
			return fmt.Errorf("expected awsJobExecutionsRolloutConfig round-trip, got %+v", info.AwsJobExecutionsRolloutConfig)
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
				if o.CreationDate == nil {
					return fmt.Errorf("expected non-nil creationDate in OTA update summary")
				}
				return nil
			}
		}
		return fmt.Errorf("%s not found in list of %d ota updates", otaID, len(out.OtaUpdates))
	}))

	results = append(results, r.RunTest("iot", "OTAUpdate_ListOTAUpdates_StatusFilter", func() error {
		complete, err := tc.client.ListOTAUpdates(tc.ctx, &iot.ListOTAUpdatesInput{
			OtaUpdateStatus: iottypes.OTAUpdateStatusCreateComplete,
		})
		if err != nil {
			return fmt.Errorf("ListOTAUpdates CREATE_COMPLETE failed: %w", err)
		}
		found := false
		for _, o := range complete.OtaUpdates {
			if o.OtaUpdateId != nil && *o.OtaUpdateId == otaID {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("%s missing from the CREATE_COMPLETE filter", otaID)
		}
		failed, err := tc.client.ListOTAUpdates(tc.ctx, &iot.ListOTAUpdatesInput{
			OtaUpdateStatus: iottypes.OTAUpdateStatusCreateFailed,
		})
		if err != nil {
			return fmt.Errorf("ListOTAUpdates CREATE_FAILED failed: %w", err)
		}
		for _, o := range failed.OtaUpdates {
			if o.OtaUpdateId != nil && *o.OtaUpdateId == otaID {
				return fmt.Errorf("%s leaked into the CREATE_FAILED filter", otaID)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "OTAUpdate_DeleteOTAUpdate_DeleteStreamFlag", func() error {
		roleArn := tc.iamRoleARN("ota-stream")
		streamKeep := uniqueName("ota-stream-keep")
		streamDrop := uniqueName("ota-stream-drop")
		otaKeep := uniqueName("ota-keep")
		otaDrop := uniqueName("ota-drop")
		defer func() {
			_, _ = tc.client.DeleteOTAUpdate(tc.ctx, &iot.DeleteOTAUpdateInput{OtaUpdateId: aws.String(otaKeep), ForceDeleteAWSJob: true})
			_, _ = tc.client.DeleteOTAUpdate(tc.ctx, &iot.DeleteOTAUpdateInput{OtaUpdateId: aws.String(otaDrop), ForceDeleteAWSJob: true})
			_, _ = tc.client.DeleteStream(tc.ctx, &iot.DeleteStreamInput{StreamId: aws.String(streamKeep)})
			_, _ = tc.client.DeleteStream(tc.ctx, &iot.DeleteStreamInput{StreamId: aws.String(streamDrop)})
		}()
		mkStream := func(id string) error {
			_, err := tc.client.CreateStream(tc.ctx, &iot.CreateStreamInput{
				StreamId: aws.String(id),
				RoleArn:  aws.String(roleArn),
				Files: []iottypes.StreamFile{{
					FileId:     aws.Int32(1),
					S3Location: &iottypes.S3Location{Bucket: aws.String("ota-bucket"), Key: aws.String("firmware.bin")},
				}},
			})
			return err
		}
		mkOTA := func(id, streamID string) error {
			_, err := tc.client.CreateOTAUpdate(tc.ctx, &iot.CreateOTAUpdateInput{
				OtaUpdateId: aws.String(id),
				RoleArn:     aws.String(roleArn),
				Targets:     []string{target},
				Files: []iottypes.OTAUpdateFile{{
					FileName:     aws.String("firmware.bin"),
					FileLocation: &iottypes.FileLocation{Stream: &iottypes.Stream{StreamId: aws.String(streamID)}},
				}},
			})
			return err
		}
		if err := mkStream(streamKeep); err != nil {
			return fmt.Errorf("CreateStream keep prerequisite failed: %w", err)
		}
		if err := mkStream(streamDrop); err != nil {
			return fmt.Errorf("CreateStream drop prerequisite failed: %w", err)
		}
		if err := mkOTA(otaKeep, streamKeep); err != nil {
			return fmt.Errorf("CreateOTAUpdate keep prerequisite failed: %w", err)
		}
		if err := mkOTA(otaDrop, streamDrop); err != nil {
			return fmt.Errorf("CreateOTAUpdate drop prerequisite failed: %w", err)
		}

		// Without the flag the referenced stream survives the OTA deletion
		// (the in-progress OTA job is force-deleted alongside the update).
		if _, err := tc.client.DeleteOTAUpdate(tc.ctx, &iot.DeleteOTAUpdateInput{
			OtaUpdateId:       aws.String(otaKeep),
			ForceDeleteAWSJob: true,
		}); err != nil {
			return fmt.Errorf("DeleteOTAUpdate without the flag failed: %w", err)
		}
		if _, err := tc.client.DescribeStream(tc.ctx, &iot.DescribeStreamInput{
			StreamId: aws.String(streamKeep),
		}); err != nil {
			return fmt.Errorf("stream must survive without deleteStream: %v", err)
		}
		// deleteStream only removes streams the OTAUpdate process itself
		// created; a user-supplied stream referenced through the files'
		// stream location is ignored by the flag, so the stream survives
		// the OTA deletion either way.
		if _, err := tc.client.DeleteOTAUpdate(tc.ctx, &iot.DeleteOTAUpdateInput{
			OtaUpdateId:       aws.String(otaDrop),
			DeleteStream:      true,
			ForceDeleteAWSJob: true,
		}); err != nil {
			return fmt.Errorf("DeleteOTAUpdate with deleteStream failed: %w", err)
		}
		if _, err := tc.client.DescribeStream(tc.ctx, &iot.DescribeStreamInput{
			StreamId: aws.String(streamDrop),
		}); err != nil {
			return fmt.Errorf("user-supplied stream must survive deleteStream: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ManagedJobTemplate_NotFoundAndEmptyList", func() error {
		// The platform ships no AWS-copyrighted managed-template
		// catalogue, so the list is legitimately empty and any describe
		// resolves to the documented not-found error.
		list, err := tc.client.ListManagedJobTemplates(tc.ctx, &iot.ListManagedJobTemplatesInput{})
		if err != nil {
			return fmt.Errorf("ListManagedJobTemplates failed: %w", err)
		}
		if len(list.ManagedJobTemplates) != 0 {
			return fmt.Errorf("expected an empty managed-template catalogue, got %d entries", len(list.ManagedJobTemplates))
		}
		_, err = tc.client.DescribeManagedJobTemplate(tc.ctx, &iot.DescribeManagedJobTemplateInput{
			TemplateName: aws.String("PackageManager"),
		})
		if err == nil {
			return fmt.Errorf("expected ResourceNotFoundException for an unknown managed template")
		}
		return expectAWSErrorCode(err, "ResourceNotFoundException")
	}))

	return results
}
