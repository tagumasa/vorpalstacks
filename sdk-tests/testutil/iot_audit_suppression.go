package testutil

import (
	"fmt"

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
			CheckName:          aws.String(checkName),
			ResourceIdentifier: resourceID,
		})
		return err
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
		_, err := tc.client.StartAuditMitigationActionsTask(tc.ctx, &iot.StartAuditMitigationActionsTaskInput{
			TaskId: aws.String(mitTaskId),
			Target: &iottypes.AuditMitigationActionsTaskTarget{AuditTaskId: aws.String("nonexistent")},
			AuditCheckToActionsMapping: map[string][]string{
				checkName: {"deviceCertMitigation"},
			},
		})
		return err
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
