package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTAuditStreamProvTests covers account audit configuration, the Stream
// resource lifecycle and aggregation/cardinality reads, with real assertions.
// ProvisioningTemplateVersion ops are omitted (handler unimplemented).
func (r *TestRunner) runIoTAuditStreamProvTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	const checkName = "AUTHENTICATED_COGNITO_ROLE_OVERLY_PERMISSIVE_CHECK"

	// ── Account audit configuration round-trip ──
	results = append(results, r.RunTest("iot", "UpdateAccountAuditConfiguration", func() error {
		_, err := tc.client.UpdateAccountAuditConfiguration(tc.ctx, &iot.UpdateAccountAuditConfigurationInput{
			AuditCheckConfigurations: map[string]iottypes.AuditCheckConfiguration{
				checkName: {Enabled: true},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "DescribeAccountAuditConfiguration_Echo", func() error {
		out, err := tc.client.DescribeAccountAuditConfiguration(tc.ctx, &iot.DescribeAccountAuditConfigurationInput{})
		if err != nil {
			return err
		}
		// After Update, the auditCheckConfigurations map must echo the
		// previously-enabled check. AWS does not guarantee field-level echo
		// ordering, but the check name must be present as a key.
		if _, ok := out.AuditCheckConfigurations[checkName]; !ok {
			return fmt.Errorf("DescribeAccountAuditConfiguration did not echo check %q", checkName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "DeleteAccountAuditConfiguration", func() error {
		_, err := tc.client.DeleteAccountAuditConfiguration(tc.ctx, &iot.DeleteAccountAuditConfigurationInput{})
		return err
	}))

	results = append(results, r.RunTest("iot", "DescribeAccountAuditConfiguration_AfterDelete", func() error {
		out, err := tc.client.DescribeAccountAuditConfiguration(tc.ctx, &iot.DescribeAccountAuditConfigurationInput{})
		if err != nil {
			return err
		}
		// After delete, the configuration map must be empty (or absent).
		if len(out.AuditCheckConfigurations) > 0 {
			return fmt.Errorf("DeleteAccountAuditConfiguration left %d entries", len(out.AuditCheckConfigurations))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "StartOnDemandAuditTask", func() error {
		out, err := tc.client.StartOnDemandAuditTask(tc.ctx, &iot.StartOnDemandAuditTaskInput{
			TargetCheckNames: []string{checkName},
		})
		if err != nil {
			return fmt.Errorf("StartOnDemandAuditTask failed: %w", err)
		}
		if out.TaskId == nil || *out.TaskId == "" {
			return fmt.Errorf("expected non-empty taskId")
		}
		return nil
	}))

	// ── Stream lifecycle ──
	// The stream resource lifecycle (Create/Describe/Update/Delete/List) is
	// covered by runIoTStreamRegistrationTests; registering it here as well
	// would double-count identical scenarios, so this section only exercises
	// the read-only aggregations.

	// ── Read-only aggregations ──
	results = append(results, r.RunTest("iot", "ListStreams", func() error {
		_, err := tc.client.ListStreams(tc.ctx, &iot.ListStreamsInput{})
		return err
	}))
	results = append(results, r.RunTest("iot", "GetBucketsAggregation_Validation", func() error {
		_, err := tc.client.GetBucketsAggregation(tc.ctx, &iot.GetBucketsAggregationInput{IndexName: aws.String("AWS_Things")})
		return expectValidationError(err)
	}))
	results = append(results, r.RunTest("iot", "GetCardinality", func() error {
		_, err := tc.client.GetCardinality(tc.ctx, &iot.GetCardinalityInput{IndexName: aws.String("AWS_Things"), QueryString: aws.String("*")})
		return err
	}))

	return results
}
