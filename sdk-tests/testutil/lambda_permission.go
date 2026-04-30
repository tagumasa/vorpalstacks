package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaPermissionTests(
	r *TestRunner,
	ctx context.Context,
	client *lambda.Client,
	createIAMRole func(string) error,
	deleteIAMRole func(string),
	region string,
) []TestResult {
	var results []TestResult

	funcName := fmt.Sprintf("PermFunc-%d", time.Now().UnixNano())
	roleName := fmt.Sprintf("PermRole-%d", time.Now().UnixNano())
	roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)

	if err := createIAMRole(roleName); err != nil {
		return []TestResult{{Service: "lambda", TestName: "Permission_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer deleteIAMRole(roleName)

	_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
		FunctionName: aws.String(funcName),
		Runtime:      types.RuntimeNodejs22x,
		Role:         aws.String(roleARN),
		Handler:      aws.String("index.handler"),
		Code:         &types.FunctionCode{ZipFile: []byte(lambdaFunctionCode)},
	})
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Permission_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(funcName)})

	var addPermStatementID string

	results = append(results, r.RunTest("lambda", "AddPermission", func() error {
		addPermStatementID = fmt.Sprintf("stmt-%d", time.Now().UnixNano())
		resp, err := client.AddPermission(ctx, &lambda.AddPermissionInput{
			FunctionName: aws.String(funcName),
			StatementId:  aws.String(addPermStatementID),
			Action:       aws.String("lambda:InvokeFunction"),
			Principal:    aws.String("apigateway.amazonaws.com"),
		})
		if err != nil {
			return err
		}
		if resp.Statement == nil || *resp.Statement == "" {
			return fmt.Errorf("Statement is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetPolicy", func() error {
		resp, err := client.GetPolicy(ctx, &lambda.GetPolicyInput{
			FunctionName: aws.String(funcName),
		})
		if err != nil {
			return err
		}
		if resp.Policy == nil || *resp.Policy == "" {
			return fmt.Errorf("policy is empty")
		}
		var policy map[string]interface{}
		if err := json.Unmarshal([]byte(*resp.Policy), &policy); err != nil {
			return fmt.Errorf("policy is not valid JSON: %v", err)
		}
		statements, ok := policy["Statement"].([]interface{})
		if !ok || len(statements) == 0 {
			return fmt.Errorf("policy has no Statement array")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "RemovePermission", func() error {
		statementID := fmt.Sprintf("stmt-%d", time.Now().UnixNano())
		_, err := client.AddPermission(ctx, &lambda.AddPermissionInput{
			FunctionName: aws.String(funcName),
			StatementId:  aws.String(statementID),
			Action:       aws.String("lambda:InvokeFunction"),
			Principal:    aws.String("apigateway.amazonaws.com"),
		})
		if err != nil {
			return err
		}
		_, err = client.RemovePermission(ctx, &lambda.RemovePermissionInput{
			FunctionName: aws.String(funcName),
			StatementId:  aws.String(statementID),
		})
		if err != nil {
			return err
		}
		policyResp, err := client.GetPolicy(ctx, &lambda.GetPolicyInput{
			FunctionName: aws.String(funcName),
		})
		if err != nil {
			return fmt.Errorf("GetPolicy after remove: %v", err)
		}
		var policy map[string]interface{}
		if err := json.Unmarshal([]byte(*policyResp.Policy), &policy); err != nil {
			return fmt.Errorf("policy is not valid JSON: %v", err)
		}
		statements, _ := policy["Statement"].([]interface{})
		for _, s := range statements {
			stmt, _ := s.(map[string]interface{})
			if sid, _ := stmt["Sid"].(string); sid == statementID {
				return fmt.Errorf("removed statement %s still present in policy", statementID)
			}
		}
		return nil
	}))

	functionARN := fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:%s", region, funcName)

	results = append(results, r.RunTest("lambda", "TagResource", func() error {
		_, err := client.TagResource(ctx, &lambda.TagResourceInput{
			Resource: aws.String(functionARN),
			Tags: map[string]string{
				"Environment": "test",
				"Project":     "sdk-tests",
			},
		})
		return err
	}))

	results = append(results, r.RunTest("lambda", "ListTags", func() error {
		resp, err := client.ListTags(ctx, &lambda.ListTagsInput{
			Resource: aws.String(functionARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		if resp.Tags["Project"] != "sdk-tests" {
			return fmt.Errorf("expected tag Project=sdk-tests, got %v", resp.Tags["Project"])
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "UntagResource", func() error {
		_, err := client.UntagResource(ctx, &lambda.UntagResourceInput{
			Resource: aws.String(functionARN),
			TagKeys:  []string{"Environment"},
		})
		if err != nil {
			return err
		}
		tagResp, err := client.ListTags(ctx, &lambda.ListTagsInput{
			Resource: aws.String(functionARN),
		})
		if err != nil {
			return fmt.Errorf("ListTags after untag: %v", err)
		}
		if _, ok := tagResp.Tags["Environment"]; ok {
			return fmt.Errorf("tag 'Environment' should have been removed")
		}
		if tagResp.Tags["Project"] != "sdk-tests" {
			return fmt.Errorf("tag 'Project' should still exist, got %v", tagResp.Tags["Project"])
		}
		return nil
	}))

	return results
}
