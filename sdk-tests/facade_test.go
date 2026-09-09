package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"vorpalstacks-sdk-tests/testutil"
)

// The go test facade exercises the same registered suites as the binary in
// main.go: each service builder first runs in register-only mode so its
// test list is captured without execution, then every captured entry runs
// as a t.Run subtest. That gives the standard runner native discovery,
// -run filtering, and per-test failure attribution. Chained scenarios
// depend on registration order, so subtests run sequentially.

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// facadeServices honours SDK_TEST_SERVICES (comma-separated); the default
// mirrors the binary's parallel phases: iam first, then the remaining SDK
// services and the WebSocket suites, then integration.
func facadeServices(r *testutil.TestRunner) []string {
	if v := os.Getenv("SDK_TEST_SERVICES"); v != "" {
		var out []string
		for _, s := range strings.Split(v, ",") {
			if s = strings.TrimSpace(s); s != "" {
				out = append(out, s)
			}
		}
		return out
	}
	ordered := []string{"iam"}
	for _, s := range r.GetServicesByCategory(testutil.CategorySDK) {
		if s != "iam" {
			ordered = append(ordered, s)
		}
	}
	ordered = append(ordered, r.GetServicesByCategory(testutil.CategoryWS)...)
	ordered = append(ordered, r.GetServicesByCategory(testutil.CategoryIntegration)...)
	return ordered
}

func TestMain(m *testing.M) {
	endpoint := envOrDefault("SDK_TEST_ENDPOINT", "http://localhost:50080")
	region := envOrDefault("SDK_TEST_REGION", "us-east-1")

	// The health gate keeps `go test ./...` green on a machine without a
	// server: every test is skipped, with the reason printed once.
	probe := testutil.NewTestRunner(endpoint, region, false)
	if err := probe.CheckServerHealth(); err != nil {
		fmt.Printf("sdk-tests facade: server at %s not reachable (%v); skipping all tests\n", endpoint, err)
		os.Exit(0)
	}
	probe.CleanupStaleContainers()
	probe.CleanupResidualLambdaFixtures()
	os.Exit(m.Run())
}

func TestService(t *testing.T) {
	endpoint := envOrDefault("SDK_TEST_ENDPOINT", "http://localhost:50080")
	region := envOrDefault("SDK_TEST_REGION", "us-east-1")

	probe := testutil.NewTestRunner(endpoint, region, false)
	for _, svc := range facadeServices(probe) {
		t.Run(svc, func(t *testing.T) {
			r := testutil.NewTestRunner(endpoint, region, false)
			r.SetRegisterOnly(true)
			r.RunServiceTests(svc)
			r.SetRegisterOnly(false)

			pending := r.PendingTests()
			if len(pending) == 0 {
				// Builders report setup failures by appending a FAIL
				// TestResult directly instead of registering; without this
				// guard such a failure would masquerade as an empty pass.
				t.Errorf("service %s registered no tests: builder setup failed outside RunTest", svc)
				return
			}
			for _, p := range pending {
				t.Run(p.TestName, func(t *testing.T) {
					res := r.ExecuteRegistered(p)
					switch res.Status {
					case "FAIL":
						t.Errorf("%s", res.Error)
					case "SKIP":
						t.Skipf("%s", res.Error)
					}
				})
			}
			// Service cleanups registered during the registration pass run
			// after the service's subtests, mirroring the binary's
			// builder-return tail.
			r.RunServiceCleanups(svc)
		})
	}
}
