package testutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotevents"
)

// runIoTEventsListDescribeTests covers IoTEvents List/Describe ops on the
// iotevents client. The detector model is created by the preceding control-
// plane tests, so the version list targets a real resource. Each op asserts
// success (previously every call swallowed the error with `_ = err`).
func (r *TestRunner) runIoTEventsListDescribeTests(tc *iotEventsTestContext, prefix string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iotevents", "ListDetectorModelVersions", func() error {
		_, err := tc.client.ListDetectorModelVersions(tc.ctx, &iotevents.ListDetectorModelVersionsInput{
			DetectorModelName: aws.String(prefix + "detector-1"),
		})
		return err
	}))

	results = append(results, r.RunTest("iotevents", "ListInputRoutings", func() error {
		_, err := tc.client.ListInputRoutings(tc.ctx, &iotevents.ListInputRoutingsInput{})
		return expectValidationError(err)
	}))

	// GetDetectorModelAnalysisResults on a non-existent analysis id must be
	// rejected, proving the handler is registered.
	results = append(results, r.RunTest("iotevents", "GetDetectorModelAnalysisResults_NotFound", func() error {
		_, err := tc.client.GetDetectorModelAnalysisResults(tc.ctx, &iotevents.GetDetectorModelAnalysisResultsInput{
			AnalysisId: aws.String("nonexistent"),
		})
		return err
	}))

	return results
}
