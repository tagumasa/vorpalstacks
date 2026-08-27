package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTAuditSuppressionTests covers the AuditSuppression CRUD lifecycle on a
// real (non-existent-cert) resource identifier plus a mitigation-actions task
// start. The lifecycle exercises Create -> Describe (echo) -> Update (partial
// merge) -> Describe (verify) -> Delete -> Describe (NotFound). A mitigation
// task is also started and cancelled to exercise the persistence path.
func (r *TestRunner) runIoTAuditSuppressionTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	checkName := "DEVICE_CERTIFICATE_EXPIRING_CHECK"
	resourceID := &iottypes.ResourceIdentifier{
		DeviceCertificateArn: aws.String(tc.arn("iot", "cert", "0000000000000000000000000000000000000000000000000000000000000000")),
	}

	results = append(results, r.RunTest("iot", "CreateAuditSuppression", func() error {
		_, err := tc.client.CreateAuditSuppression(tc.ctx, &iot.CreateAuditSuppressionInput{
			CheckName:            aws.String(checkName),
			ResourceIdentifier:   resourceID,
			SuppressIndefinitely: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		// Recreating the same (checkName, resourceIdentifier) tuple must
		// raise ResourceAlreadyExistsException, not overwrite the record.
		_, err = tc.client.CreateAuditSuppression(tc.ctx, &iot.CreateAuditSuppressionInput{
			CheckName:            aws.String(checkName),
			ResourceIdentifier:   resourceID,
			SuppressIndefinitely: aws.Bool(true),
		})
		return expectAWSErrorCode(err, "ResourceAlreadyExistsException")
	}))

	results = append(results, r.RunTest("iot", "Suppression_ExpirationExclusivity", func() error {
		// The expirationDate and suppressIndefinitely members are mutually
		// exclusive: neither or both must be rejected as InvalidRequest.
		_, err := tc.client.CreateAuditSuppression(tc.ctx, &iot.CreateAuditSuppressionInput{
			CheckName:          aws.String(uniqueName("xor-neither")),
			ResourceIdentifier: resourceID,
		})
		if err := expectAWSErrorCode(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("neither member: %w", err)
		}
		_, err = tc.client.CreateAuditSuppression(tc.ctx, &iot.CreateAuditSuppressionInput{
			CheckName:            aws.String(uniqueName("xor-both")),
			ResourceIdentifier:   resourceID,
			SuppressIndefinitely: aws.Bool(true),
			ExpirationDate:       aws.Time(time.Now().Add(24 * time.Hour)),
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	results = append(results, r.RunTest("iot", "DescribeAuditSuppression_Echo", func() error {
		out, err := tc.client.DescribeAuditSuppression(tc.ctx, &iot.DescribeAuditSuppressionInput{
			CheckName:          aws.String(checkName),
			ResourceIdentifier: resourceID,
		})
		if err != nil {
			return err
		}
		if out.CheckName == nil || *out.CheckName != checkName {
			return fmt.Errorf("DescribeAuditSuppression checkName mismatch: got %v", out.CheckName)
		}
		if !aws.ToBool(out.SuppressIndefinitely) {
			return fmt.Errorf("expected suppressIndefinitely echoed back as true")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "UpdateAuditSuppression", func() error {
		_, err := tc.client.UpdateAuditSuppression(tc.ctx, &iot.UpdateAuditSuppressionInput{
			CheckName:          aws.String(checkName),
			ResourceIdentifier: resourceID,
			Description:        aws.String("updated by test"),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "DescribeAuditSuppression_AfterUpdate", func() error {
		out, err := tc.client.DescribeAuditSuppression(tc.ctx, &iot.DescribeAuditSuppressionInput{
			CheckName:          aws.String(checkName),
			ResourceIdentifier: resourceID,
		})
		if err != nil {
			return err
		}
		if out.Description == nil || *out.Description != "updated by test" {
			return fmt.Errorf("UpdateAuditSuppression did not persist description: got %v", out.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "DeleteAuditSuppression", func() error {
		_, err := tc.client.DeleteAuditSuppression(tc.ctx, &iot.DeleteAuditSuppressionInput{
			CheckName:          aws.String(checkName),
			ResourceIdentifier: resourceID,
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "DescribeAuditSuppression_NotFound", func() error {
		_, err := tc.client.DescribeAuditSuppression(tc.ctx, &iot.DescribeAuditSuppressionInput{
			CheckName:          aws.String(checkName),
			ResourceIdentifier: resourceID,
		})
		return expectNotFound(err)
	}))

	mitTaskId := uniqueName("audit-mit")
	results = append(results, r.RunTest("iot", "StartAuditMitigationActionsTask", func() error {
		out, err := tc.client.StartAuditMitigationActionsTask(tc.ctx, &iot.StartAuditMitigationActionsTaskInput{
			TaskId: aws.String(mitTaskId),
			Target: &iottypes.AuditMitigationActionsTaskTarget{AuditTaskId: aws.String("nonexistent")},
			AuditCheckToActionsMapping: map[string][]string{
				checkName: {"deviceCertMitigation"},
			},
		})
		if err != nil {
			return err
		}
		if aws.ToString(out.TaskId) != mitTaskId {
			return fmt.Errorf("expected taskId %q echoed, got %v", mitTaskId, out.TaskId)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "AuditMitigation_TaskDescribeShape", func() error {
		out, err := tc.client.DescribeAuditMitigationActionsTask(tc.ctx, &iot.DescribeAuditMitigationActionsTaskInput{TaskId: aws.String(mitTaskId)})
		if err != nil {
			return err
		}
		if out.TaskStatus != iottypes.AuditMitigationActionsTaskStatusInProgress {
			return fmt.Errorf("expected taskStatus=IN_PROGRESS, got %s", out.TaskStatus)
		}
		if out.StartTime == nil {
			return fmt.Errorf("expected non-nil startTime")
		}
		if out.Target == nil || aws.ToString(out.Target.AuditTaskId) != "nonexistent" {
			return fmt.Errorf("expected target.auditTaskId echoed, got %+v", out.Target)
		}
		if _, ok := out.AuditCheckToActionsMapping[checkName]; !ok {
			return fmt.Errorf("expected auditCheckToActionsMapping to echo check %q", checkName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "AuditMitigation_ListTasks_Members", func() error {
		out, err := tc.client.ListAuditMitigationActionsTasks(tc.ctx, &iot.ListAuditMitigationActionsTasksInput{
			StartTime: aws.Time(time.Now().Add(-time.Hour)),
			EndTime:   aws.Time(time.Now().Add(time.Hour)),
		})
		if err != nil {
			return err
		}
		for _, task := range out.Tasks {
			if aws.ToString(task.TaskId) == mitTaskId {
				if task.StartTime == nil {
					return fmt.Errorf("expected non-nil startTime on task metadata")
				}
				return nil
			}
		}
		return fmt.Errorf("started task %q not found in ListAuditMitigationActionsTasks", mitTaskId)
	}))

	results = append(results, r.RunTest("iot", "CancelAuditMitigationActionsTask_NotFound", func() error {
		_, err := tc.client.CancelAuditMitigationActionsTask(tc.ctx, &iot.CancelAuditMitigationActionsTaskInput{
			TaskId: aws.String(uniqueName("nope-mit-task")),
		})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "CancelAuditMitigationActionsTask_Existing", func() error {
		_, err := tc.client.CancelAuditMitigationActionsTask(tc.ctx, &iot.CancelAuditMitigationActionsTaskInput{
			TaskId: aws.String(mitTaskId),
		})
		return err
	}))

	return results
}
