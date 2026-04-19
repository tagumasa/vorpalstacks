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
	"sync"
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

type TestCategory string

const (
	CategorySDK         TestCategory = "sdk"
	CategoryWS          TestCategory = "ws"
	CategoryIntegration TestCategory = "integration"
)

type serviceEntry struct {
	factory  ServiceFactory
	category TestCategory
}

var serviceRegistry = map[string]serviceEntry{}

func RegisterService(name string, category TestCategory, factory ServiceFactory) {
	serviceRegistry[name] = serviceEntry{factory: factory, category: category}
}

func init() {
	RegisterService("dynamodb", CategorySDK, (*TestRunner).RunDynamoDBTests)
	RegisterService("sqs", CategorySDK, (*TestRunner).RunSQSTests)
	RegisterService("sns", CategorySDK, (*TestRunner).RunSNSTests)
	RegisterService("s3", CategorySDK, (*TestRunner).RunS3Tests)
	RegisterService("lambda", CategorySDK, (*TestRunner).RunLambdaTests)
	RegisterService("iam", CategorySDK, (*TestRunner).RunIAMTests)
	RegisterService("kms", CategorySDK, (*TestRunner).RunKMSTests)
	RegisterService("events", CategorySDK, (*TestRunner).RunEventBridgeTests)
	RegisterService("stepfunctions", CategorySDK, (*TestRunner).RunStepFunctionsTests)
	RegisterService("apigateway", CategorySDK, (*TestRunner).RunAPIGatewayTests)
	RegisterService("logs", CategorySDK, (*TestRunner).RunCloudWatchLogsTests)
	RegisterService("cloudwatch", CategorySDK, (*TestRunner).RunCloudWatchTests)
	RegisterService("ssm", CategorySDK, (*TestRunner).RunSSMTests)
	RegisterService("cloudtrail", CategorySDK, (*TestRunner).RunCloudTrailTests)
	RegisterService("acm", CategorySDK, (*TestRunner).RunACMTests)
	RegisterService("cognito", CategorySDK, (*TestRunner).RunCognitoTests)
	RegisterService("cognito-identity", CategorySDK, (*TestRunner).RunCognitoIdentityTests)
	RegisterService("secretsmanager", CategorySDK, (*TestRunner).RunSecretsManagerTests)
	RegisterService("kinesis", CategorySDK, (*TestRunner).RunKinesisTests)
	RegisterService("sts", CategorySDK, (*TestRunner).RunSTSTests)
	RegisterService("scheduler", CategorySDK, (*TestRunner).RunSchedulerTests)
	RegisterService("athena", CategorySDK, (*TestRunner).RunAthenaTests)
	RegisterService("timestream", CategorySDK, (*TestRunner).RunTimestreamTests)
	RegisterService("sesv2", CategorySDK, (*TestRunner).RunSESv2Tests)
	RegisterService("route53", CategorySDK, (*TestRunner).RunRoute53Tests)
	RegisterService("cloudfront", CategorySDK, (*TestRunner).RunCloudFrontTests)
	RegisterService("wafv2", CategorySDK, (*TestRunner).RunWAFv2Tests)
	RegisterService("neptune", CategorySDK, (*TestRunner).RunNeptuneTests)
	RegisterService("neptunedata", CategorySDK, (*TestRunner).RunNeptunedataTests)
	RegisterService("neptunegraph", CategorySDK, (*TestRunner).RunNeptunegraphTests)
	RegisterService("appsync", CategorySDK, (*TestRunner).RunAppSyncTests)
	RegisterService("appsync-ws", CategoryWS, (*TestRunner).RunAppSyncWSTests)
	RegisterService("integration", CategoryIntegration, (*TestRunner).RunIntegrationTests)
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
	return r.GetServicesByCategory("")
}

func (r *TestRunner) GetServicesByCategory(category TestCategory) []string {
	services := make([]string, 0, len(serviceRegistry))
	for name, entry := range serviceRegistry {
		if category == "" || entry.category == category {
			services = append(services, name)
		}
	}
	sort.Strings(services)
	return services
}

func (r *TestRunner) GetServiceCategory(name string) TestCategory {
	if entry, ok := serviceRegistry[name]; ok {
		return entry.category
	}
	return CategorySDK
}

func (r *TestRunner) RunServiceTests(service string) []TestResult {
	entry, ok := serviceRegistry[service]
	if !ok {
		return []TestResult{{
			Service:  service,
			TestName: "Unknown",
			Status:   "SKIP",
			Error:    "Unknown service",
		}}
	}
	return entry.factory(r)
}

var iamDependentServices = map[string]bool{
	"lambda":         true,
	"stepfunctions":  true,
	"scheduler":      true,
	"timestream":     true,
	"logs":           true,
	"sts":            true,
}

func (r *TestRunner) RunServicesParallel(services []string, parallelism int) map[string][]TestResult {
	if parallelism <= 1 {
		results := make(map[string][]TestResult)
		for _, svc := range services {
			results[svc] = r.RunServiceTests(svc)
		}
		return results
	}

	var phase1, phase2, phase3 []string
	for _, svc := range services {
		switch {
		case svc == "iam":
			phase1 = append(phase1, svc)
		case svc == "integration":
			phase3 = append(phase3, svc)
		default:
			phase2 = append(phase2, svc)
		}
	}

	results := make(map[string][]TestResult)
	var mu sync.Mutex

	runPhase := func(svcs []string) {
		if len(svcs) == 0 {
			return
		}
		sem := make(chan struct{}, parallelism)
		var wg sync.WaitGroup
		for _, svc := range svcs {
			wg.Add(1)
			go func(s string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				res := r.RunServiceTests(s)
				mu.Lock()
				results[s] = res
				mu.Unlock()
			}(svc)
		}
		wg.Wait()
	}

	for _, phase := range [][]string{phase1, phase2, phase3} {
		runPhase(phase)
	}

	return results
}

const perTestTimeout = 60 * time.Second

func (r *TestRunner) RunTest(service, testName string, testFunc func() error) TestResult {
	if r.verbose {
		fmt.Printf("  Running: %s...\n", testName)
	}

	start := time.Now()
	errCh := make(chan error, 1)
	go func() {
		errCh <- testFunc()
	}()

	var err error
	select {
	case err = <-errCh:
	case <-time.After(perTestTimeout):
		err = fmt.Errorf("test timed out after %v", perTestTimeout)
	}
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

type CategorySummary struct {
	Passed  int
	Failed  int
	Skipped int
}

func (r *TestRunner) SummariseByCategory(results map[string][]TestResult) map[TestCategory]*CategorySummary {
	summary := map[TestCategory]*CategorySummary{
		CategorySDK:         {},
		CategoryWS:          {},
		CategoryIntegration: {},
	}
	for svc, svcResults := range results {
		cat := r.GetServiceCategory(svc)
		for _, result := range svcResults {
			switch result.Status {
			case "PASS":
				summary[cat].Passed++
			case "FAIL":
				summary[cat].Failed++
			case "SKIP":
				summary[cat].Skipped++
			}
		}
	}
	return summary
}

func (r *TestRunner) printTableReport(results map[string][]TestResult) {
	categoryOrder := []TestCategory{CategorySDK, CategoryIntegration, CategoryWS}
	categoryHeaders := map[TestCategory]string{
		CategorySDK:         "SDK TESTS",
		CategoryIntegration: "INTEGRATION TESTS",
		CategoryWS:          "WEBSOCKET TESTS",
	}

	for _, cat := range categoryOrder {
		catResults := make(map[string][]TestResult)
		for svc, svcResults := range results {
			if r.GetServiceCategory(svc) == cat {
				catResults[svc] = svcResults
			}
		}
		if len(catResults) == 0 {
			continue
		}

		fmt.Printf("\n========== %s ==========\n", categoryHeaders[cat])
		for _, service := range r.GetServicesByCategory(cat) {
			serviceResults, exists := catResults[service]
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
}

func (r *TestRunner) printJSONReport(results map[string][]TestResult) {
	type categoryBreakdown struct {
		Passed  int `json:"passed"`
		Failed  int `json:"failed"`
		Skipped int `json:"skipped"`
	}
	report := struct {
		Results map[string][]TestResult `json:"results"`
		Summary struct {
			Passed  int `json:"passed"`
			Failed  int `json:"failed"`
			Skipped int `json:"skipped"`
		} `json:"summary"`
		ByCategory map[TestCategory]*categoryBreakdown `json:"byCategory"`
	}{
		Results:    results,
		ByCategory: map[TestCategory]*categoryBreakdown{},
	}

	catSummary := r.SummariseByCategory(results)
	for cat, s := range catSummary {
		report.ByCategory[cat] = &categoryBreakdown{
			Passed:  s.Passed,
			Failed:  s.Failed,
			Skipped: s.Skipped,
		}
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
