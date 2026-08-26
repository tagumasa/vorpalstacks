package testutil

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/smithy-go"
)

// testHTTPClient is a shared *http.Client with a sensible timeout so that a
// hung server connection cannot block the test goroutine indefinitely.
var testHTTPClient = &http.Client{Timeout: 30 * time.Second}

func SetupFailResult(service, msg string, args ...interface{}) TestResult {
	return TestResult{
		Service:  service,
		TestName: "Setup",
		Status:   "FAIL",
		Error:    fmt.Sprintf(msg, args...),
	}
}

func AssertErrorContains(err error, expectedType string) error {
	if err == nil {
		return fmt.Errorf("expected error containing %q, got nil", expectedType)
	}
	if !strings.Contains(err.Error(), expectedType) {
		return fmt.Errorf("expected error containing %q, got: %v", expectedType, err)
	}
	return nil
}

func AssertNoError(err error, context string) error {
	if err != nil {
		return fmt.Errorf("%s: %v", context, err)
	}
	return nil
}

func AssertNotNil(v interface{}, name string) error {
	if v == nil {
		return fmt.Errorf("%s is nil", name)
	}
	return nil
}

// expectAWSErrorCode asserts that err is a smithy API error carrying the
// expected AWS error code string (for example "ValidationException" or
// "InvalidObjectState"), regardless of which service client produced it.
func expectAWSErrorCode(err error, code string) error {
	if err == nil {
		return fmt.Errorf("expected %s error, got nil", code)
	}
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		return fmt.Errorf("expected %s error, got non-API error: %T: %v", code, err, err)
	}
	if apiErr.ErrorCode() != code {
		return fmt.Errorf("expected %s, got %s: %v", code, apiErr.ErrorCode(), err)
	}
	return nil
}

// paginate walks a token-based list API to completion, collecting every
// item across all pages. The fetch closure performs one page request and
// returns the page items plus the continuation token; pagination stops
// when the token is nil or empty.
func paginate[T any](fetch func(next *string) ([]T, *string, error)) ([]T, error) {
	var all []T
	var next *string
	for {
		items, token, err := fetch(next)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if token == nil || *token == "" {
			return all, nil
		}
		next = token
	}
}

// containsID returns a pointer to the first item for which match reports
// true, or nil when no item matches.
func containsID[T any](items []T, match func(*T) bool) *T {
	for i := range items {
		if match(&items[i]) {
			return &items[i]
		}
	}
	return nil
}

func IAMCreateRole(client *iam.Client, roleName, trustPolicy string) error {
	_, err := client.CreateRole(context.Background(), &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(trustPolicy),
	})
	return err
}

func IAMDeleteRole(client *iam.Client, roleName string) {
	client.DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})
}

func IAMDeleteRoleCtx(ctx context.Context, client *iam.Client, roleName string) {
	client.DeleteRole(ctx, &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})
}
