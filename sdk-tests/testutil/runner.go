package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
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

type ServiceFactory func(*TestRunner) []TestResult

var serviceRegistry = map[string]ServiceFactory{}

func RegisterService(name string, factory ServiceFactory) {
	serviceRegistry[name] = factory
}

func init() {
	RegisterService("dynamodb", (*TestRunner).RunDynamoDBTests)
	RegisterService("sqs", (*TestRunner).RunSQSTests)
	RegisterService("sns", (*TestRunner).RunSNSTests)
	RegisterService("s3", (*TestRunner).RunS3Tests)
	RegisterService("lambda", (*TestRunner).RunLambdaTests)
	RegisterService("iam", (*TestRunner).RunIAMTests)
	RegisterService("kms", (*TestRunner).RunKMSTests)
	RegisterService("events", (*TestRunner).RunEventBridgeTests)
	RegisterService("stepfunctions", (*TestRunner).RunStepFunctionsTests)
	RegisterService("apigateway", (*TestRunner).RunAPIGatewayTests)
	RegisterService("logs", (*TestRunner).RunCloudWatchLogsTests)
	RegisterService("cloudwatch", (*TestRunner).RunCloudWatchTests)
	RegisterService("ssm", (*TestRunner).RunSSMTests)
	RegisterService("cloudtrail", (*TestRunner).RunCloudTrailTests)
	RegisterService("acm", (*TestRunner).RunACMTests)
	RegisterService("cognito", (*TestRunner).RunCognitoTests)
	RegisterService("cognito-identity", (*TestRunner).RunCognitoIdentityTests)
	RegisterService("secretsmanager", (*TestRunner).RunSecretsManagerTests)
	RegisterService("kinesis", (*TestRunner).RunKinesisTests)
	RegisterService("sts", (*TestRunner).RunSTSTests)
	RegisterService("scheduler", (*TestRunner).RunSchedulerTests)
	RegisterService("athena", (*TestRunner).RunAthenaTests)
	RegisterService("timestream", (*TestRunner).RunTimestreamTests)
	RegisterService("sesv2", (*TestRunner).RunSESv2Tests)
	RegisterService("route53", (*TestRunner).RunRoute53Tests)
	RegisterService("cloudfront", (*TestRunner).RunCloudFrontTests)
	RegisterService("wafv2", (*TestRunner).RunWAFv2Tests)
	RegisterService("neptune", (*TestRunner).RunNeptuneTests)
	RegisterService("neptunedata", (*TestRunner).RunNeptunedataTests)
	RegisterService("appsync", (*TestRunner).RunAppSyncTests)
	RegisterService("appsync-ws", (*TestRunner).RunAppSyncWSTests)
}

func NewTestRunner(endpoint, region string, verbose bool) *TestRunner {
	return &TestRunner{
		endpoint: endpoint,
		region:   region,
		client:   &http.Client{Timeout: 30 * time.Second},
		verbose:  verbose,
	}
}

func (r *TestRunner) Endpoint() string { return r.endpoint }
func (r *TestRunner) Region() string   { return r.region }

func (r *TestRunner) GetAllServices() []string {
	services := make([]string, 0, len(serviceRegistry))
	for name := range serviceRegistry {
		services = append(services, name)
	}
	sort.Strings(services)
	return services
}

func (r *TestRunner) RunServiceTests(service string) []TestResult {
	factory, ok := serviceRegistry[service]
	if !ok {
		return []TestResult{{
			Service:  service,
			TestName: "Unknown",
			Status:   "SKIP",
			Error:    "Unknown service",
		}}
	}
	return factory(r)
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

func (r *TestRunner) SkipTest(service, testName, reason string) TestResult {
	if r.verbose {
		fmt.Printf("  - %s (skipped: %s)\n", testName, reason)
	}
	return TestResult{
		Service:  service,
		TestName: testName,
		Status:   "SKIP",
		Error:    reason,
	}
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
	resp.Body.Close()

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
