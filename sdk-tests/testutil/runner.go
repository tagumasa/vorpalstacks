package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type TestResult struct {
	Service  string
	TestName string
	Status   string
	Error    string
	Duration time.Duration
}

type TestRunner struct {
	endpoint string
	region   string
	client   *http.Client
	verbose  bool
}

func NewTestRunner(endpoint, region string, verbose bool) *TestRunner {
	return &TestRunner{
		endpoint: endpoint,
		region:   region,
		client:   &http.Client{Timeout: 30 * time.Second},
		verbose:  verbose,
	}
}

func (r *TestRunner) GetAllServices() []string {
	return []string{
		"dynamodb", "sqs", "sns", "s3", "lambda", "iam", "kms", "events",
		"stepfunctions", "apigateway", "logs", "cloudwatch", "ssm", "cloudtrail",
		"acm", "cognito", "cognito-identity", "secretsmanager", "kinesis", "sts", "scheduler",
		"athena", "timestream", "sesv2", "route53", "cloudfront", "waf",
		"neptune", "neptunedata",
	}
}

func (r *TestRunner) RunServiceTests(service string) []TestResult {
	var results []TestResult

	switch service {
	case "dynamodb":
		results = r.RunDynamoDBTests()
	case "sqs":
		results = r.RunSQSTests()
	case "sns":
		results = r.RunSNSTests()
	case "s3":
		results = r.RunS3Tests()
	case "lambda":
		results = r.RunLambdaTests()
	case "iam":
		results = r.RunIAMTests()
	case "kms":
		results = r.RunKMSTests()
	case "events":
		results = r.RunEventBridgeTests()
	case "stepfunctions":
		results = r.RunStepFunctionsTests()
	case "apigateway":
		results = r.RunAPIGatewayTests()
	case "logs":
		results = r.RunCloudWatchLogsTests()
	case "cloudwatch":
		results = r.RunCloudWatchTests()
	case "ssm":
		results = r.RunSSMTests()
	case "cloudtrail":
		results = r.RunCloudTrailTests()
	case "acm":
		results = r.RunACMTests()
	case "cognito":
		results = r.RunCognitoTests()
	case "cognito-identity":
		results = r.RunCognitoIdentityTests()
	case "secretsmanager":
		results = r.RunSecretsManagerTests()
	case "kinesis":
		results = r.RunKinesisTests()
	case "sts":
		results = r.RunSTSTests()
	case "scheduler":
		results = r.RunSchedulerTests()
	case "athena":
		results = r.RunAthenaTests()
	case "timestream":
		results = r.RunTimestreamTests()
	case "sesv2":
		results = r.RunSESv2Tests()
	case "route53":
		results = r.RunRoute53Tests()
	case "cloudfront":
		results = r.RunCloudFrontTests()
	case "waf":
		results = r.RunWAFTests()
	case "neptune":
		results = r.RunNeptuneTests()
	case "neptunedata":
		results = r.RunNeptunedataTests()
	default:
		results = []TestResult{{
			Service:  service,
			TestName: "Unknown",
			Status:   "SKIP",
			Error:    "Unknown service",
		}}
	}

	return results
}

func (r *TestRunner) RunTest(service, testName string, testFunc func() error) TestResult {
	if r.verbose {
		fmt.Printf("  Running: %s...\n", testName)
	}

	start := time.Now()
	err := testFunc()
	duration := time.Since(start)

	result := TestResult{
		Service:  service,
		TestName: testName,
		Status:   "PASS",
		Duration: duration,
	}

	if err != nil {
		result.Status = "FAIL"
		result.Error = err.Error()
	}

	if r.verbose {
		if result.Status == "PASS" {
			fmt.Printf("  ✓ %s (%.2fs)\n", testName, duration.Seconds())
		} else {
			fmt.Printf("  ✗ %s: %s\n", testName, err.Error())
		}
	}

	return result
}

func (r *TestRunner) PrintReport(results map[string][]TestResult, format string) {
	if format == "json" {
		r.printJSONReport(results)
	} else {
		r.printTableReport(results)
	}
}

func (r *TestRunner) printTableReport(results map[string][]TestResult) {
	for _, service := range r.GetAllServices() {
		serviceResults, exists := results[service]
		if !exists || len(serviceResults) == 0 {
			continue
		}

		fmt.Printf("\n=== %s ===\n", strings.ToUpper(service))
		for _, result := range serviceResults {
			statusSymbol := "✓"
			if result.Status == "FAIL" {
				statusSymbol = "✗"
			} else if result.Status == "SKIP" {
				statusSymbol = "-"
			}
			fmt.Printf("  %s %s (%.2fs)", statusSymbol, result.TestName, result.Duration.Seconds())
			if result.Error != "" {
				fmt.Printf(" - %s", result.Error)
			}
			fmt.Println()
		}
	}
}

func (r *TestRunner) printJSONReport(results map[string][]TestResult) {
	report := struct {
		Results map[string][]TestResult
		Summary struct {
			Passed  int
			Failed  int
			Skipped int
		}
	}{
		Results: results,
	}

	for _, serviceResults := range results {
		for _, result := range serviceResults {
			switch result.Status {
			case "PASS":
				report.Summary.Passed++
			case "FAIL":
				report.Summary.Failed++
			case "SKIP":
				report.Summary.Skipped++
			}
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(report)
}

func (r *TestRunner) CheckServerHealth() error {
	resp, err := r.client.Get(r.endpoint + "/health")
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *TestRunner) SendHTTPTestRequest(method, path string, body []byte) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, r.endpoint+path, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Amz-Target", "test")

	return r.client.Do(req)
}
