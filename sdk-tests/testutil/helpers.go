package testutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
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
