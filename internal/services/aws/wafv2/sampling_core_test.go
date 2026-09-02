package wafv2

import (
	"strings"
	"testing"
	"time"
)

// TestGetSampledRequestsCoreMemberErrors pins the member-specific
// required-member and range errors, all before any store access.
func TestGetSampledRequestsCoreMemberErrors(t *testing.T) {
	svc := &WAFv2Service{}
	base := GetSampledRequestsInput{
		WebACLArn:         "arn",
		RuleMetricName:    "metric",
		Scope:             "REGIONAL",
		StartTime:         time.Now().Add(-time.Minute),
		EndTime:           time.Now(),
		MaxItems:          10,
		TimeWindowPresent: true,
		MaxItemsPresent:   true,
	}

	cases := map[string]func(*GetSampledRequestsInput){
		"missing WebAclArn":      func(in *GetSampledRequestsInput) { in.WebACLArn = "" },
		"missing RuleMetricName": func(in *GetSampledRequestsInput) { in.RuleMetricName = "" },
		"missing Scope":          func(in *GetSampledRequestsInput) { in.Scope = "" },
		"missing TimeWindow":     func(in *GetSampledRequestsInput) { in.TimeWindowPresent = false },
		"missing StartTime":      func(in *GetSampledRequestsInput) { in.StartTime = time.Time{} },
		"missing EndTime":        func(in *GetSampledRequestsInput) { in.EndTime = time.Time{} },
		"end before start":       func(in *GetSampledRequestsInput) { in.EndTime = in.StartTime },
		"missing MaxItems":       func(in *GetSampledRequestsInput) { in.MaxItemsPresent = false },
		"MaxItems below range":   func(in *GetSampledRequestsInput) { in.MaxItems = 0 },
		"MaxItems above range":   func(in *GetSampledRequestsInput) { in.MaxItems = 501 },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			input := base
			mutate(&input)
			_, err := svc.getSampledRequestsCore(nil, input)
			if err == nil {
				t.Fatal("expected a validation error, got nil")
			}
			if !strings.Contains(err.Error(), "WAFInvalidParameterException") {
				t.Fatalf("expected WAFInvalidParameterException, got %v", err)
			}
		})
	}
}

// TestGetRateBasedManagedKeysCoreMemberErrors pins the required-member
// errors of the managed-keys operation.
func TestGetRateBasedManagedKeysCoreMemberErrors(t *testing.T) {
	svc := &WAFv2Service{}
	base := GetRateBasedManagedKeysInput{
		Scope:      "REGIONAL",
		WebACLName: "acl",
		WebACLId:   "id",
		RuleName:   "rate",
	}

	cases := map[string]func(*GetRateBasedManagedKeysInput){
		"missing Scope":      func(in *GetRateBasedManagedKeysInput) { in.Scope = "" },
		"missing WebACLName": func(in *GetRateBasedManagedKeysInput) { in.WebACLName = "" },
		"missing WebACLId":   func(in *GetRateBasedManagedKeysInput) { in.WebACLId = "" },
		"missing RuleName":   func(in *GetRateBasedManagedKeysInput) { in.RuleName = "" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			input := base
			mutate(&input)
			_, err := svc.getRateBasedManagedKeysCore(nil, input)
			if err == nil {
				t.Fatal("expected a validation error, got nil")
			}
			if !strings.Contains(err.Error(), "WAFInvalidParameterException") {
				t.Fatalf("expected WAFInvalidParameterException, got %v", err)
			}
		})
	}
}
