package testutil

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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

func IAMCreateRole(endpoint, roleName, trustPolicy string) error {
	form := url.Values{}
	form.Set("Action", "CreateRole")
	form.Set("Version", "2010-05-08")
	form.Set("RoleName", roleName)
	form.Set("AssumeRolePolicyDocument", trustPolicy)
	req, err := http.NewRequestWithContext(context.Background(), "POST", endpoint, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func IAMDeleteRole(endpoint, roleName string) {
	form := url.Values{}
	form.Set("Action", "DeleteRole")
	form.Set("Version", "2010-05-08")
	form.Set("RoleName", roleName)
	req, _ := http.NewRequestWithContext(context.Background(), "POST", endpoint, bytes.NewBufferString(form.Encode()))
	if req != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
	}
}

func IAMDeleteRoleCtx(ctx context.Context, endpoint, roleName string) {
	form := url.Values{}
	form.Set("Action", "DeleteRole")
	form.Set("Version", "2010-05-08")
	form.Set("RoleName", roleName)
	req, _ := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBufferString(form.Encode()))
	if req != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if resp, _ := http.DefaultClient.Do(req); resp != nil {
			resp.Body.Close()
		}
	}
}
