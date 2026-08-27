package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTRemainingSmokeTests covers certificate-transfer, audit/detect task and
// job-execution handlers on non-existent resources. Many of these handlers are
// stubs that accept any id and no-op; the meaningful assertion is that the
// handler is registered and responds without a server error (previously every
// call swallowed the error with `_ = err`, so a 500 passed silently).
func (r *TestRunner) runIoTRemainingSmokeTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	fakeCert := "0000000000000000000000000000000000000000000000000000000000000000"
	fakeTask := uniqueName("nope-task")
	fakeFinding := uniqueName("nope-finding")
	fakeDestARN := tc.arn("iot", "ruledestination", "nonexistent")

	// Certificate transfer handlers (the platform records/updates a transfer
	// state for any id; assert the handler responds).
	transfer := func(name string, fn func() error) {
		results = append(results, r.RunTest("iot", name, fn))
	}
	transfer("TransferCertificate", func() error {
		_, err := tc.client.TransferCertificate(tc.ctx, &iot.TransferCertificateInput{
			CertificateId:    aws.String(fakeCert),
			TargetAwsAccount: aws.String("000000000001"),
		})
		return expectNotFound(err)
	})
	transfer("AcceptCertificateTransfer", func() error {
		_, err := tc.client.AcceptCertificateTransfer(tc.ctx, &iot.AcceptCertificateTransferInput{CertificateId: aws.String(fakeCert)})
		// Completing a transfer that does not exist returns NotFound.
		return expectNotFound(err)
	})
	transfer("CancelCertificateTransfer", func() error {
		_, err := tc.client.CancelCertificateTransfer(tc.ctx, &iot.CancelCertificateTransferInput{CertificateId: aws.String(fakeCert)})
		return expectNotFound(err)
	})
	transfer("RejectCertificateTransfer", func() error {
		_, err := tc.client.RejectCertificateTransfer(tc.ctx, &iot.RejectCertificateTransferInput{CertificateId: aws.String(fakeCert)})
		return expectNotFound(err)
	})

	// RegisterCertificateWithoutCA with an invalid PEM.
	results = append(results, r.RunTest("iot", "RegisterCertificateWithoutCA", func() error {
		out, err := tc.client.RegisterCertificateWithoutCA(tc.ctx, &iot.RegisterCertificateWithoutCAInput{
			CertificatePem: aws.String(uniqueName("pem-")),
		})
		if err != nil {
			return err
		}
		if out.CertificateId != nil && *out.CertificateId != "" {
			defer tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{
				CertificateId: out.CertificateId,
				ForceDelete:   true,
			})
		}
		return nil
	}))

	// Audit / detect task handlers persist Start* output and resolve task ids
	// via ResourceNotFoundException for unknown ids, matching the Smithy error
	// traits on Cancel/Describe for each operation.
	task := func(name string, fn func() error) {
		results = append(results, r.RunTest("iot", name, fn))
	}
	task("CancelAuditTask_NotFound", func() error {
		_, err := tc.client.CancelAuditTask(tc.ctx, &iot.CancelAuditTaskInput{TaskId: aws.String(fakeTask)})
		return expectNotFound(err)
	})
	task("DescribeAuditTask_NotFound", func() error {
		_, err := tc.client.DescribeAuditTask(tc.ctx, &iot.DescribeAuditTaskInput{TaskId: aws.String(fakeTask)})
		return expectNotFound(err)
	})
	task("DescribeAuditFinding_NotFound", func() error {
		_, err := tc.client.DescribeAuditFinding(tc.ctx, &iot.DescribeAuditFindingInput{FindingId: aws.String(fakeFinding)})
		return expectNotFound(err)
	})
	task("CancelAuditMitigationActionsTask_NotFound", func() error {
		_, err := tc.client.CancelAuditMitigationActionsTask(tc.ctx, &iot.CancelAuditMitigationActionsTaskInput{TaskId: aws.String(fakeTask)})
		return expectNotFound(err)
	})
	task("DescribeAuditMitigationActionsTask_NotFound", func() error {
		_, err := tc.client.DescribeAuditMitigationActionsTask(tc.ctx, &iot.DescribeAuditMitigationActionsTaskInput{TaskId: aws.String(fakeTask)})
		return expectNotFound(err)
	})
	task("CancelDetectMitigationActionsTask_NotFound", func() error {
		_, err := tc.client.CancelDetectMitigationActionsTask(tc.ctx, &iot.CancelDetectMitigationActionsTaskInput{TaskId: aws.String(fakeTask)})
		return expectNotFound(err)
	})
	task("DescribeDetectMitigationActionsTask_NotFound", func() error {
		_, err := tc.client.DescribeDetectMitigationActionsTask(tc.ctx, &iot.DescribeDetectMitigationActionsTaskInput{TaskId: aws.String(fakeTask)})
		return expectNotFound(err)
	})

	// Job execution / registration task reads: DescribeJobExecution now
	// returns ResourceNotFoundException for unknown job/thing; the
	// ThingRegistrationTask variants remain validation-error stubs.
	results = append(results, r.RunTest("iot", "DescribeJobExecution_Validation", func() error {
		_, err := tc.client.DescribeJobExecution(tc.ctx, &iot.DescribeJobExecutionInput{
			JobId: aws.String("nonexistent"), ThingName: aws.String("nonexistent"),
		})
		return expectNotFound(err)
	}))
	results = append(results, r.RunTest("iot", "DescribeThingRegistrationTask_Validation", func() error {
		// An unknown, well-formed task ID yields the documented
		// ResourceNotFoundException (404).
		_, err := tc.client.DescribeThingRegistrationTask(tc.ctx, &iot.DescribeThingRegistrationTaskInput{TaskId: aws.String(fakeTask)})
		return expectNotFound(err)
	}))
	task("StopThingRegistrationTask", func() error {
		_, err := tc.client.StopThingRegistrationTask(tc.ctx, &iot.StopThingRegistrationTaskInput{TaskId: aws.String(fakeTask)})
		return expectNotFound(err)
	})

	// Rule destination / connectivity / managed-template reads.
	results = append(results, r.RunTest("iot", "GetTopicRuleDestination_NotFound", func() error {
		_, err := tc.client.GetTopicRuleDestination(tc.ctx, &iot.GetTopicRuleDestinationInput{Arn: aws.String(fakeDestARN)})
		return expectNotFound(err)
	}))
	results = append(results, r.RunTest("iot", "GetThingConnectivityData", func() error {
		_, err := tc.client.GetThingConnectivityData(tc.ctx, &iot.GetThingConnectivityDataInput{ThingName: aws.String("nonexistent")})
		return expectNotFound(err)
	}))
	results = append(results, r.RunTest("iot", "DescribeManagedJobTemplate", func() error {
		_, err := tc.client.DescribeManagedJobTemplate(tc.ctx, &iot.DescribeManagedJobTemplateInput{TemplateName: aws.String("AWS-Reset-Accesskey")})
		return err
	}))

	// Associate targets / update custom metric / start detect mitigation task
	// (stubs that accept any id).
	task("AssociateTargetsWithJob_NotFound", func() error {
		_, err := tc.client.AssociateTargetsWithJob(tc.ctx, &iot.AssociateTargetsWithJobInput{
			JobId:   aws.String(uniqueName("nope-job")),
			Targets: []string{tc.arn("iot", "thing", "nonexistent")},
		})
		return expectNotFound(err)
	})
	task("UpdateCustomMetric", func() error {
		_, err := tc.client.UpdateCustomMetric(tc.ctx, &iot.UpdateCustomMetricInput{
			MetricName:  aws.String(uniqueName("nope-metric")),
			DisplayName: aws.String("updated"),
		})
		// Updating a non-existent custom metric now correctly returns
		// ResourceNotFoundException (404); assert that rather than swallowing.
		return expectNotFound(err)
	})
	task("StartDetectMitigationActionsTask", func() error {
		_, err := tc.client.StartDetectMitigationActionsTask(tc.ctx, &iot.StartDetectMitigationActionsTaskInput{
			TaskId:                        aws.String(uniqueName("detect-mit")),
			Target:                        &iottypes.DetectMitigationActionsTaskTarget{},
			Actions:                       []string{"deviceCertMitigation"},
			ViolationEventOccurrenceRange: &iottypes.ViolationEventOccurrenceRange{},
		})
		return expectValidationError(err)
	})

	results = append(results, r.RunTest("iot", "DetectMitigation_TaskSummaryShape", func() error {
		taskId := uniqueName("detect-mit-shape")
		defer tc.client.CancelDetectMitigationActionsTask(tc.ctx, &iot.CancelDetectMitigationActionsTaskInput{TaskId: aws.String(taskId)})
		if _, err := tc.client.StartDetectMitigationActionsTask(tc.ctx, &iot.StartDetectMitigationActionsTaskInput{
			TaskId:             aws.String(taskId),
			Target:             &iottypes.DetectMitigationActionsTaskTarget{ViolationIds: []string{"violation-1"}},
			Actions:            []string{"deviceCertMitigation"},
			ClientRequestToken: aws.String("detect-mit-shape-token"),
		}); err != nil {
			return fmt.Errorf("StartDetectMitigationActionsTask failed: %w", err)
		}
		out, err := tc.client.DescribeDetectMitigationActionsTask(tc.ctx, &iot.DescribeDetectMitigationActionsTaskInput{TaskId: aws.String(taskId)})
		if err != nil {
			return err
		}
		if out.TaskSummary == nil {
			return fmt.Errorf("expected taskSummary in response")
		}
		if aws.ToString(out.TaskSummary.TaskId) != taskId {
			return fmt.Errorf("expected taskSummary.taskId=%s, got %v", taskId, out.TaskSummary.TaskId)
		}
		if out.TaskSummary.TaskStatus != iottypes.DetectMitigationActionsTaskStatusInProgress {
			return fmt.Errorf("expected taskStatus=IN_PROGRESS, got %s", out.TaskSummary.TaskStatus)
		}
		if out.TaskSummary.TaskStartTime == nil {
			return fmt.Errorf("expected non-nil taskStartTime")
		}
		if out.TaskSummary.Target == nil || len(out.TaskSummary.Target.ViolationIds) != 1 {
			return fmt.Errorf("expected target.violationIds echoed")
		}
		return nil
	}))

	return results
}
