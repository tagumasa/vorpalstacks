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
		targetServices = tester.GetAllServices()
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

	passed := 0
	failed := 0
	for _, results := range allResults {
		for _, r := range results {
			if r.Status == "PASS" {
				passed++
			} else if r.Status == "FAIL" {
				failed++
			}
		}
	}

	fmt.Printf("\nSummary: %d passed, %d failed\n", passed, failed)
	if failed > 0 {
		os.Exit(1)
	}
}
