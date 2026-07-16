package testutil

func (r *TestRunner) RunSESv2Tests() []TestResult {
	tc, err := newSESTestContext(r.endpoint, r.region)
	if err != nil {
		return []TestResult{{
			Service:  "sesv2",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		}}
	}

	var results []TestResult
	results = append(results, r.runSESv2AccountTests(tc)...)
	results = append(results, r.runSESv2ConfigSetTests(tc)...)
	results = append(results, r.runSESv2IdentityTests(tc)...)
	results = append(results, r.runSESv2EmailTests(tc)...)
	results = append(results, r.runSESv2TemplateTests(tc)...)
	results = append(results, r.runSESv2PoolTests(tc)...)
	results = append(results, r.runSESv2SuppressionTests(tc)...)
	results = append(results, r.runSESv2ContactTests(tc)...)
	results = append(results, r.runSESv2EdgeTests(tc)...)
	return results
}
