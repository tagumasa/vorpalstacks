package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"vorpalstacks-sdk-tests/testutil"
)

var (
	endpoint = flag.String("endpoint", "http://localhost:8080", "VorpalStacks endpoint")
	region   = flag.String("region", "us-east-1", "AWS region")
	services = flag.String("service", "", "Comma-separated list of services to test (or 'all')")
	testType = flag.String("type", "all", "Test type to run: all, sdk, ws, integration")
	format   = flag.String("format", "table", "Output format: table, json")
	verbose  = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "Usage: go run main.go [options]\n")
		os.Exit(1)
	}

	tester := testutil.NewTestRunner(*endpoint, *region, *verbose)

	var targetServices []string
	if *services == "" || *services == "all" {
		cat := testutil.TestCategory(*testType)
		if *testType == "all" {
			cat = ""
		}
		targetServices = tester.GetServicesByCategory(cat)
	} else {
		targetServices = strings.Split(*services, ",")
	}

	allResults := make(map[string][]testutil.TestResult)
	for _, svc := range targetServices {
		svc = strings.TrimSpace(svc)
		if svc == "" {
			continue
		}
		results := tester.RunServiceTests(svc)
		allResults[svc] = results
	}

	tester.PrintReport(allResults, *format)

	catSummary := tester.SummariseByCategory(allResults)
	totalPassed, totalFailed, totalSkipped := 0, 0, 0

	fmt.Println()
	for _, cat := range []testutil.TestCategory{testutil.CategorySDK, testutil.CategoryIntegration, testutil.CategoryWS} {
		s := catSummary[cat]
		totalPassed += s.Passed
		totalFailed += s.Failed
		totalSkipped += s.Skipped
		fmt.Printf("  %-15s %d passed, %d failed, %d skipped\n", cat+":", s.Passed, s.Failed, s.Skipped)
	}
	fmt.Printf("\n  %-15s %d passed, %d failed, %d skipped\n", "TOTAL:", totalPassed, totalFailed, totalSkipped)

	if totalFailed > 0 {
		os.Exit(1)
	}
}
