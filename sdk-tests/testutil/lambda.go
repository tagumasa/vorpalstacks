package testutil

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunLambdaTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "lambda",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := lambda.NewFromConfig(cfg)
	iamClient := iam.NewFromConfig(cfg)
	ctx := context.Background()

	createIAMRole := func(roleName string) error {
		return IAMCreateRole(iamClient, roleName, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}`)
	}

	deleteIAMRole := func(roleName string) {
		IAMDeleteRole(iamClient, roleName)
	}

	functionName := fmt.Sprintf("TestFunction-%d", time.Now().UnixNano())
	roleName := fmt.Sprintf("TestRole-%d", time.Now().UnixNano())
	roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)

	if err := createIAMRole(roleName); err != nil {
		results = append(results, TestResult{Service: "lambda", TestName: "Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
		return results
	}
	defer deleteIAMRole(roleName)

	functionCode := []byte("exports.handler = async (event) => { return { statusCode: 200, body: 'Hello' }; };")

	results = append(results, r.RunTest("lambda", "CreateFunction", func() error {
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(functionName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code: &types.FunctionCode{
				ZipFile: functionCode,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("lambda", "GetFunction", func() error {
		resp, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp.Configuration == nil {
			return fmt.Errorf("configuration is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetFunctionConfiguration", func() error {
		resp, err := client.GetFunctionConfiguration(ctx, &lambda.GetFunctionConfigurationInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp.FunctionName == nil {
			return fmt.Errorf("function name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ListFunctions", func() error {
		resp, err := client.ListFunctions(ctx, &lambda.ListFunctionsInput{})
		if err != nil {
			return err
		}
		if resp.Functions == nil {
			return fmt.Errorf("functions list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "UpdateFunctionCode", func() error {
		newCode := []byte("exports.handler = async (event) => { return { statusCode: 200, body: 'Updated' }; };")
		resp, err := client.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(functionName),
			ZipFile:      newCode,
		})
		if err != nil {
			return err
		}
		if resp.LastModified == nil {
			return fmt.Errorf("LastModified is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "UpdateFunctionConfiguration", func() error {
		description := "Updated function"
		resp, err := client.UpdateFunctionConfiguration(ctx, &lambda.UpdateFunctionConfigurationInput{
			FunctionName: aws.String(functionName),
			Description:  aws.String(description),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "PublishVersion", func() error {
		resp, err := client.PublishVersion(ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp.Version == nil {
			return fmt.Errorf("Version is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ListVersionsByFunction", func() error {
		resp, err := client.ListVersionsByFunction(ctx, &lambda.ListVersionsByFunctionInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp.Versions == nil {
			return fmt.Errorf("versions list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "CreateAlias", func() error {
		resp, err := client.CreateAlias(ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(functionName),
			Name:            aws.String("live"),
			FunctionVersion: aws.String("$LATEST"),
		})
		if err != nil {
			return err
		}
		if resp.AliasArn == nil {
			return fmt.Errorf("AliasArn is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetAlias", func() error {
		resp, err := client.GetAlias(ctx, &lambda.GetAliasInput{
			FunctionName: aws.String(functionName),
			Name:         aws.String("live"),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil {
			return fmt.Errorf("alias name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "UpdateAlias", func() error {
		resp, err := client.UpdateAlias(ctx, &lambda.UpdateAliasInput{
			FunctionName: aws.String(functionName),
			Name:         aws.String("live"),
			Description:  aws.String("Production alias"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ListAliases", func() error {
		resp, err := client.ListAliases(ctx, &lambda.ListAliasesInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp.Aliases == nil {
			return fmt.Errorf("aliases list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "Invoke", func() error {
		resp, err := client.Invoke(ctx, &lambda.InvokeInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == 0 {
			return fmt.Errorf("status code is zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "PutFunctionConcurrency", func() error {
		resp, err := client.PutFunctionConcurrency(ctx, &lambda.PutFunctionConcurrencyInput{
			FunctionName:                 aws.String(functionName),
			ReservedConcurrentExecutions: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetFunctionConcurrency", func() error {
		resp, err := client.GetFunctionConcurrency(ctx, &lambda.GetFunctionConcurrencyInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp.ReservedConcurrentExecutions == nil {
			return fmt.Errorf("concurrency is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "DeleteFunctionConcurrency", func() error {
		resp, err := client.DeleteFunctionConcurrency(ctx, &lambda.DeleteFunctionConcurrencyInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	var addPermStatementID string

	results = append(results, r.RunTest("lambda", "AddPermission", func() error {
		addPermStatementID = fmt.Sprintf("stmt-%d", time.Now().UnixNano())
		resp, err := client.AddPermission(ctx, &lambda.AddPermissionInput{
			FunctionName: aws.String(functionName),
			StatementId:  aws.String(addPermStatementID),
			Action:       aws.String("lambda:InvokeFunction"),
			Principal:    aws.String("apigateway.amazonaws.com"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("AddPermission response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetPolicy", func() error {
		resp, err := client.GetPolicy(ctx, &lambda.GetPolicyInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp.Policy == nil || *resp.Policy == "" {
			return fmt.Errorf("policy is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "RemovePermission", func() error {
		statementID := fmt.Sprintf("stmt-%d", time.Now().UnixNano())
		_, err := client.AddPermission(ctx, &lambda.AddPermissionInput{
			FunctionName: aws.String(functionName),
			StatementId:  aws.String(statementID),
			Action:       aws.String("lambda:InvokeFunction"),
			Principal:    aws.String("apigateway.amazonaws.com"),
		})
		if err != nil {
			return err
		}
		_, err = client.RemovePermission(ctx, &lambda.RemovePermissionInput{
			FunctionName: aws.String(functionName),
			StatementId:  aws.String(statementID),
		})
		return err
	}))

	results = append(results, r.RunTest("lambda", "TagResource", func() error {
		functionARN := fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:%s", r.region, functionName)
		resp, err := client.TagResource(ctx, &lambda.TagResourceInput{
			Resource: aws.String(functionARN),
			Tags: map[string]string{
				"Environment": "test",
				"Project":     "sdk-tests",
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ListTags", func() error {
		functionARN := fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:%s", r.region, functionName)
		resp, err := client.ListTags(ctx, &lambda.ListTagsInput{
			Resource: aws.String(functionARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "UntagResource", func() error {
		functionARN := fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:%s", r.region, functionName)
		resp, err := client.UntagResource(ctx, &lambda.UntagResourceInput{
			Resource: aws.String(functionARN),
			TagKeys:  []string{"Environment"},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "DeleteAlias", func() error {
		resp, err := client.DeleteAlias(ctx, &lambda.DeleteAliasInput{
			FunctionName: aws.String(functionName),
			Name:         aws.String("live"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetAccountSettings", func() error {
		resp, err := client.GetAccountSettings(ctx, &lambda.GetAccountSettingsInput{})
		if err != nil {
			return err
		}
		if resp.AccountLimit == nil {
			return fmt.Errorf("account limit is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "DeleteFunction", func() error {
		if addPermStatementID != "" {
			client.RemovePermission(ctx, &lambda.RemovePermissionInput{
				FunctionName: aws.String(functionName),
				StatementId:  aws.String(addPermStatementID),
			})
		}
		resp, err := client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("lambda", "GetFunction_NonExistent", func() error {
		_, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent function")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "Invoke_NonExistent", func() error {
		_, err := client.Invoke(ctx, &lambda.InvokeInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent function")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "UpdateFunctionCode_NonExistent", func() error {
		_, err := client.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
			ZipFile:      []byte("code"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent function")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "DeleteFunction_NonExistent", func() error {
		_, err := client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent function")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "CreateFunction_DuplicateName", func() error {
		dupName := fmt.Sprintf("DupFunc-%d", time.Now().UnixNano())
		dupRoleName := fmt.Sprintf("DupRole-%d", time.Now().UnixNano())
		dupRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", dupRoleName)
		dupCode := []byte("exports.handler = async () => { return 1; };")
		if err := createIAMRole(dupRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(dupRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(dupName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(dupRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: dupCode},
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(dupName)})

		_, err = client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(dupName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(dupRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: dupCode},
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate function name")
		}
		var riu *types.ResourceConflictException
		if !errors.As(err, &riu) {
			return fmt.Errorf("expected ResourceConflictException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "Invoke_VerifyResponsePayload", func() error {
		invFunc := fmt.Sprintf("InvFunc-%d", time.Now().UnixNano())
		invRoleName := fmt.Sprintf("InvRole-%d", time.Now().UnixNano())
		invRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", invRoleName)
		invCode := []byte("exports.handler = async (event) => { return { statusCode: 200, body: JSON.stringify({result: 'ok'}) }; };")
		if err := createIAMRole(invRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(invRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(invFunc),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(invRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: invCode},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(invFunc)})

		resp, err := client.Invoke(ctx, &lambda.InvokeInput{
			FunctionName: aws.String(invFunc),
		})
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if len(resp.Payload) == 0 {
			return fmt.Errorf("expected non-empty payload")
		}
		payload, err := io.ReadAll(bytes.NewReader(resp.Payload))
		if err != nil {
			return fmt.Errorf("read payload: %v", err)
		}
		if string(payload) == "" {
			return fmt.Errorf("payload should not be empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetFunction_ContainsCodeConfig", func() error {
		gfcFunc := fmt.Sprintf("GfcFunc-%d", time.Now().UnixNano())
		gfcRoleName := fmt.Sprintf("GfcRole-%d", time.Now().UnixNano())
		gfcRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", gfcRoleName)
		gfcCode := []byte("exports.handler = async () => { return 1; };")
		gfcDesc := "Test description for verification"
		if err := createIAMRole(gfcRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(gfcRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(gfcFunc),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(gfcRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: gfcCode},
			Description:  aws.String(gfcDesc),
			Timeout:      aws.Int32(15),
			MemorySize:   aws.Int32(256),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(gfcFunc)})

		resp, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
			FunctionName: aws.String(gfcFunc),
		})
		if err != nil {
			return fmt.Errorf("get function: %v", err)
		}
		if resp.Configuration == nil {
			return fmt.Errorf("configuration is nil")
		}
		if resp.Configuration.Description == nil || *resp.Configuration.Description != gfcDesc {
			return fmt.Errorf("description mismatch, got %v", resp.Configuration.Description)
		}
		if resp.Configuration.Timeout == nil || *resp.Configuration.Timeout != 15 {
			return fmt.Errorf("timeout mismatch, got %v", resp.Configuration.Timeout)
		}
		if resp.Configuration.MemorySize == nil || *resp.Configuration.MemorySize != 256 {
			return fmt.Errorf("memory size mismatch, got %v", resp.Configuration.MemorySize)
		}
		if resp.Configuration.Runtime != types.RuntimeNodejs22x {
			return fmt.Errorf("runtime mismatch, got %v", resp.Configuration.Runtime)
		}
		if resp.Code == nil || resp.Code.Location == nil {
			return fmt.Errorf("code location should not be nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "PublishVersion_VerifyVersion", func() error {
		pvFunc := fmt.Sprintf("PvFunc-%d", time.Now().UnixNano())
		pvRoleName := fmt.Sprintf("PvRole-%d", time.Now().UnixNano())
		pvRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", pvRoleName)
		pvCode := []byte("exports.handler = async () => { return 1; };")
		if err := createIAMRole(pvRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(pvRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(pvFunc),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(pvRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: pvCode},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(pvFunc)})

		resp, err := client.PublishVersion(ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(pvFunc),
		})
		if err != nil {
			return fmt.Errorf("publish: %v", err)
		}
		if resp.Version == nil || *resp.Version == "$LATEST" {
			return fmt.Errorf("published version should not be $LATEST, got %v", resp.Version)
		}
		if resp.Version != nil && *resp.Version != "1" {
			return fmt.Errorf("first published version should be 1, got %v", resp.Version)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ListFunctions_ReturnsCreated", func() error {
		lfFunc := fmt.Sprintf("LfFunc-%d", time.Now().UnixNano())
		lfRoleName := fmt.Sprintf("LfRole-%d", time.Now().UnixNano())
		lfRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", lfRoleName)
		lfCode := []byte("exports.handler = async () => { return 1; };")
		if err := createIAMRole(lfRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(lfRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(lfFunc),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(lfRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: lfCode},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(lfFunc)})

		resp, err := client.ListFunctions(ctx, &lambda.ListFunctionsInput{})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		found := false
		for _, f := range resp.Functions {
			if f.FunctionName != nil && *f.FunctionName == lfFunc {
				found = true
				if f.Runtime != types.RuntimeNodejs22x {
					return fmt.Errorf("runtime mismatch in list, got %v", f.Runtime)
				}
				if f.Handler == nil || *f.Handler != "index.handler" {
					return fmt.Errorf("handler mismatch in list, got %v", f.Handler)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created function %s not found in ListFunctions", lfFunc)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "CreateAlias_DuplicateName", func() error {
		caFunc := fmt.Sprintf("CaFunc-%d", time.Now().UnixNano())
		caRoleName := fmt.Sprintf("CaRole-%d", time.Now().UnixNano())
		caRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", caRoleName)
		caCode := []byte("exports.handler = async () => { return 1; };")
		if err := createIAMRole(caRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(caRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(caFunc),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(caRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: caCode},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(caFunc)})

		_, err = client.CreateAlias(ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(caFunc),
			Name:            aws.String("prod"),
			FunctionVersion: aws.String("$LATEST"),
		})
		if err != nil {
			return fmt.Errorf("first alias: %v", err)
		}

		_, err = client.CreateAlias(ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(caFunc),
			Name:            aws.String("prod"),
			FunctionVersion: aws.String("$LATEST"),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate alias name")
		}
		var riu *types.ResourceConflictException
		if !errors.As(err, &riu) {
			return fmt.Errorf("expected ResourceConflictException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "UpdateFunctionConfiguration_VerifyUpdate", func() error {
		ucFunc := fmt.Sprintf("UcFunc-%d", time.Now().UnixNano())
		ucRoleName := fmt.Sprintf("UcRole-%d", time.Now().UnixNano())
		ucRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", ucRoleName)
		ucCode := []byte("exports.handler = async () => { return 1; };")
		if err := createIAMRole(ucRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(ucRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(ucFunc),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(ucRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: ucCode},
			Description:  aws.String("original"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(ucFunc)})

		newDesc := "updated description"
		newTimeout := int32(30)
		newMemory := int32(512)
		_, err = client.UpdateFunctionConfiguration(ctx, &lambda.UpdateFunctionConfigurationInput{
			FunctionName: aws.String(ucFunc),
			Description:  aws.String(newDesc),
			Timeout:      aws.Int32(newTimeout),
			MemorySize:   aws.Int32(newMemory),
		})
		if err != nil {
			return fmt.Errorf("update config: %v", err)
		}

		resp, err := client.GetFunctionConfiguration(ctx, &lambda.GetFunctionConfigurationInput{
			FunctionName: aws.String(ucFunc),
		})
		if err != nil {
			return fmt.Errorf("get config: %v", err)
		}
		if resp.Description == nil || *resp.Description != newDesc {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}
		if resp.Timeout == nil || *resp.Timeout != newTimeout {
			return fmt.Errorf("timeout not updated, got %v", resp.Timeout)
		}
		if resp.MemorySize == nil || *resp.MemorySize != newMemory {
			return fmt.Errorf("memory size not updated, got %v", resp.MemorySize)
		}
		return nil
	}))

	// === GROUP A: LAYER OPERATIONS ===

	layerName := fmt.Sprintf("TestLayer-%d", time.Now().UnixNano())
	layerZipContent := base64.StdEncoding.EncodeToString([]byte("exports.handler = async (event) => { return 1; };"))

	results = append(results, r.RunTest("lambda", "PublishLayerVersion", func() error {
		resp, err := client.PublishLayerVersion(ctx, &lambda.PublishLayerVersionInput{
			LayerName: aws.String(layerName),
			Content: &types.LayerVersionContentInput{
				ZipFile: []byte(layerZipContent),
			},
			Description:        aws.String("Test layer version"),
			CompatibleRuntimes: []types.Runtime{types.RuntimeNodejs22x},
		})
		if err != nil {
			return err
		}
		if resp.LayerArn == nil {
			return fmt.Errorf("LayerArn is nil")
		}
		if resp.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Version)
		}
		if resp.Content == nil || resp.Content.CodeSha256 == nil {
			return fmt.Errorf("CodeSha256 is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetLayerVersion", func() error {
		resp, err := client.GetLayerVersion(ctx, &lambda.GetLayerVersionInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(1),
		})
		if err != nil {
			return err
		}
		if resp.Content == nil || resp.Content.CodeSha256 == nil {
			return fmt.Errorf("CodeSha256 is nil")
		}
		if resp.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Version)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ListLayers", func() error {
		resp, err := client.ListLayers(ctx, &lambda.ListLayersInput{})
		if err != nil {
			return err
		}
		if resp.Layers == nil {
			return fmt.Errorf("layers list is nil")
		}
		found := false
		for _, l := range resp.Layers {
			if l.LayerName != nil && *l.LayerName == layerName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("layer %s not found in ListLayers", layerName)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ListLayerVersions", func() error {
		resp, err := client.ListLayerVersions(ctx, &lambda.ListLayerVersionsInput{
			LayerName: aws.String(layerName),
		})
		if err != nil {
			return err
		}
		if resp.LayerVersions == nil {
			return fmt.Errorf("layer versions list is nil")
		}
		if len(resp.LayerVersions) == 0 {
			return fmt.Errorf("expected at least 1 layer version")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "DeleteLayerVersion", func() error {
		_, err := client.DeleteLayerVersion(ctx, &lambda.DeleteLayerVersionInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(1),
		})
		return err
	}))

	// === GROUP B: EVENT SOURCE MAPPING ===

	esmFuncName := fmt.Sprintf("EsmFunc-%d", time.Now().UnixNano())
	esmRoleName := fmt.Sprintf("EsmRole-%d", time.Now().UnixNano())
	esmRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", esmRoleName)
	esmCode := []byte("exports.handler = async () => { return 1; };")
	esmEventSourceArn := "arn:aws:sqs:us-east-1:000000000000:test-queue"

	if err := createIAMRole(esmRoleName); err != nil {
		results = append(results, TestResult{Service: "lambda", TestName: "CreateEventSourceMapping_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
	} else {
		defer deleteIAMRole(esmRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(esmFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(esmRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: esmCode},
		})
		if err != nil {
			results = append(results, TestResult{Service: "lambda", TestName: "CreateEventSourceMapping_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create function: %v", err)})
		} else {
			defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(esmFuncName)})

			var esmUUID string

			results = append(results, r.RunTest("lambda", "CreateEventSourceMapping", func() error {
				resp, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
					FunctionName:   aws.String(esmFuncName),
					EventSourceArn: aws.String(esmEventSourceArn),
					BatchSize:      aws.Int32(10),
					Enabled:        aws.Bool(true),
				})
				if err != nil {
					return err
				}
				if resp.UUID == nil || *resp.UUID == "" {
					return fmt.Errorf("UUID is nil or empty")
				}
				esmUUID = *resp.UUID
				return nil
			}))

			results = append(results, r.RunTest("lambda", "GetEventSourceMapping", func() error {
				if esmUUID == "" {
					return fmt.Errorf("no UUID from CreateEventSourceMapping")
				}
				resp, err := client.GetEventSourceMapping(ctx, &lambda.GetEventSourceMappingInput{
					UUID: aws.String(esmUUID),
				})
				if err != nil {
					return err
				}
				if resp.FunctionArn == nil {
					return fmt.Errorf("FunctionArn is nil")
				}
				if resp.EventSourceArn == nil || *resp.EventSourceArn != esmEventSourceArn {
					return fmt.Errorf("EventSourceArn mismatch, got %v", resp.EventSourceArn)
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "UpdateEventSourceMapping", func() error {
				if esmUUID == "" {
					return fmt.Errorf("no UUID from CreateEventSourceMapping")
				}
				resp, err := client.UpdateEventSourceMapping(ctx, &lambda.UpdateEventSourceMappingInput{
					UUID:      aws.String(esmUUID),
					BatchSize: aws.Int32(50),
					Enabled:   aws.Bool(false),
				})
				if err != nil {
					return err
				}
				if resp.BatchSize == nil || *resp.BatchSize != 50 {
					return fmt.Errorf("BatchSize not updated, got %v", resp.BatchSize)
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "ListEventSourceMappings", func() error {
				resp, err := client.ListEventSourceMappings(ctx, &lambda.ListEventSourceMappingsInput{
					FunctionName: aws.String(esmFuncName),
				})
				if err != nil {
					return err
				}
				if resp.EventSourceMappings == nil {
					return fmt.Errorf("event source mappings list is nil")
				}
				if len(resp.EventSourceMappings) == 0 {
					return fmt.Errorf("expected at least 1 event source mapping")
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "DeleteEventSourceMapping", func() error {
				if esmUUID == "" {
					return fmt.Errorf("no UUID from CreateEventSourceMapping")
				}
				_, err := client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{
					UUID: aws.String(esmUUID),
				})
				return err
			}))

			results = append(results, r.RunTest("lambda", "GetEventSourceMapping_NonExistent", func() error {
				_, err := client.GetEventSourceMapping(ctx, &lambda.GetEventSourceMappingInput{
					UUID: aws.String("00000000-0000-0000-0000-000000000000"),
				})
				if err == nil {
					return fmt.Errorf("expected error for non-existent event source mapping")
				}
				var rnf *types.ResourceNotFoundException
				if !errors.As(err, &rnf) {
					return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
				}
				return nil
			}))
		}
	}

	// === GROUP C: PROVISIONED CONCURRENCY ===

	pcFuncName := fmt.Sprintf("PcFunc-%d", time.Now().UnixNano())
	pcRoleName := fmt.Sprintf("PcRole-%d", time.Now().UnixNano())
	pcRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", pcRoleName)
	pcCode := []byte("exports.handler = async () => { return 1; };")

	if err := createIAMRole(pcRoleName); err != nil {
		results = append(results, TestResult{Service: "lambda", TestName: "PutProvisionedConcurrencyConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
	} else {
		defer deleteIAMRole(pcRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(pcFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(pcRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: pcCode},
		})
		if err != nil {
			results = append(results, TestResult{Service: "lambda", TestName: "PutProvisionedConcurrencyConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create function: %v", err)})
		} else {
			defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(pcFuncName)})

			publishResp, err := client.PublishVersion(ctx, &lambda.PublishVersionInput{
				FunctionName: aws.String(pcFuncName),
			})
			if err != nil {
				results = append(results, TestResult{Service: "lambda", TestName: "PutProvisionedConcurrencyConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to publish version: %v", err)})
			} else {
				pcVersion := *publishResp.Version

				results = append(results, r.RunTest("lambda", "PutProvisionedConcurrencyConfig", func() error {
					resp, err := client.PutProvisionedConcurrencyConfig(ctx, &lambda.PutProvisionedConcurrencyConfigInput{
						FunctionName:                    aws.String(pcFuncName),
						Qualifier:                       aws.String(pcVersion),
						ProvisionedConcurrentExecutions: aws.Int32(5),
					})
					if err != nil {
						return err
					}
					if resp.AllocatedProvisionedConcurrentExecutions == nil {
						return fmt.Errorf("AllocatedProvisionedConcurrentExecutions is nil")
					}
					return nil
				}))

				results = append(results, r.RunTest("lambda", "GetProvisionedConcurrencyConfig", func() error {
					resp, err := client.GetProvisionedConcurrencyConfig(ctx, &lambda.GetProvisionedConcurrencyConfigInput{
						FunctionName: aws.String(pcFuncName),
						Qualifier:    aws.String(pcVersion),
					})
					if err != nil {
						return err
					}
					if resp.Status == "" {
						return fmt.Errorf("Status is empty")
					}
					return nil
				}))

				results = append(results, r.RunTest("lambda", "ListProvisionedConcurrencyConfigs", func() error {
					resp, err := client.ListProvisionedConcurrencyConfigs(ctx, &lambda.ListProvisionedConcurrencyConfigsInput{
						FunctionName: aws.String(pcFuncName),
					})
					if err != nil {
						return err
					}
					if resp.ProvisionedConcurrencyConfigs == nil {
						return fmt.Errorf("configs list is nil")
					}
					if len(resp.ProvisionedConcurrencyConfigs) == 0 {
						return fmt.Errorf("expected at least 1 config")
					}
					return nil
				}))

				results = append(results, r.RunTest("lambda", "DeleteProvisionedConcurrencyConfig", func() error {
					_, err := client.DeleteProvisionedConcurrencyConfig(ctx, &lambda.DeleteProvisionedConcurrencyConfigInput{
						FunctionName: aws.String(pcFuncName),
						Qualifier:    aws.String(pcVersion),
					})
					return err
				}))

				results = append(results, r.RunTest("lambda", "GetProvisionedConcurrencyConfig_NonExistent", func() error {
					_, err := client.GetProvisionedConcurrencyConfig(ctx, &lambda.GetProvisionedConcurrencyConfigInput{
						FunctionName: aws.String(pcFuncName),
						Qualifier:    aws.String(pcVersion),
					})
					if err == nil {
						return fmt.Errorf("expected error for deleted provisioned concurrency config")
					}
					var rnf *types.ResourceNotFoundException
					if !errors.As(err, &rnf) {
						return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
					}
					return nil
				}))
			}
		}
	}

	// === GROUP D: FUNCTION EVENT INVOKE CONFIG ===

	eicFuncName := fmt.Sprintf("EicFunc-%d", time.Now().UnixNano())
	eicRoleName := fmt.Sprintf("EicRole-%d", time.Now().UnixNano())
	eicRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", eicRoleName)
	eicCode := []byte("exports.handler = async () => { return 1; };")

	if err := createIAMRole(eicRoleName); err != nil {
		results = append(results, TestResult{Service: "lambda", TestName: "PutFunctionEventInvokeConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
	} else {
		defer deleteIAMRole(eicRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(eicFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(eicRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: eicCode},
		})
		if err != nil {
			results = append(results, TestResult{Service: "lambda", TestName: "PutFunctionEventInvokeConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create function: %v", err)})
		} else {
			defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(eicFuncName)})

			results = append(results, r.RunTest("lambda", "PutFunctionEventInvokeConfig", func() error {
				maxAge := int32(3600)
				maxRetries := int32(2)
				resp, err := client.PutFunctionEventInvokeConfig(ctx, &lambda.PutFunctionEventInvokeConfigInput{
					FunctionName:             aws.String(eicFuncName),
					MaximumEventAgeInSeconds: aws.Int32(maxAge),
					MaximumRetryAttempts:     aws.Int32(maxRetries),
				})
				if err != nil {
					return err
				}
				if resp.LastModified == nil {
					return fmt.Errorf("LastModified is nil")
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "GetFunctionEventInvokeConfig", func() error {
				resp, err := client.GetFunctionEventInvokeConfig(ctx, &lambda.GetFunctionEventInvokeConfigInput{
					FunctionName: aws.String(eicFuncName),
				})
				if err != nil {
					return err
				}
				if resp.MaximumEventAgeInSeconds == nil || *resp.MaximumEventAgeInSeconds != 3600 {
					return fmt.Errorf("MaximumEventAgeInSeconds mismatch, got %v", resp.MaximumEventAgeInSeconds)
				}
				if resp.MaximumRetryAttempts == nil || *resp.MaximumRetryAttempts != 2 {
					return fmt.Errorf("MaximumRetryAttempts mismatch, got %v", resp.MaximumRetryAttempts)
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "ListFunctionEventInvokeConfigs", func() error {
				resp, err := client.ListFunctionEventInvokeConfigs(ctx, &lambda.ListFunctionEventInvokeConfigsInput{
					FunctionName: aws.String(eicFuncName),
				})
				if err != nil {
					return err
				}
				if resp.FunctionEventInvokeConfigs == nil {
					return fmt.Errorf("configs list is nil")
				}
				if len(resp.FunctionEventInvokeConfigs) == 0 {
					return fmt.Errorf("expected at least 1 config")
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "DeleteFunctionEventInvokeConfig", func() error {
				_, err := client.DeleteFunctionEventInvokeConfig(ctx, &lambda.DeleteFunctionEventInvokeConfigInput{
					FunctionName: aws.String(eicFuncName),
				})
				return err
			}))
		}
	}

	// === GROUP E: FUNCTION URL CONFIG ===

	furlFuncName := fmt.Sprintf("FurlFunc-%d", time.Now().UnixNano())
	furlRoleName := fmt.Sprintf("FurlRole-%d", time.Now().UnixNano())
	furlRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", furlRoleName)
	furlCode := []byte("exports.handler = async () => { return 1; };")

	if err := createIAMRole(furlRoleName); err != nil {
		results = append(results, TestResult{Service: "lambda", TestName: "CreateFunctionUrlConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
	} else {
		defer deleteIAMRole(furlRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(furlFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(furlRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: furlCode},
		})
		if err != nil {
			results = append(results, TestResult{Service: "lambda", TestName: "CreateFunctionUrlConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create function: %v", err)})
		} else {
			defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(furlFuncName)})

			results = append(results, r.RunTest("lambda", "CreateFunctionUrlConfig", func() error {
				resp, err := client.CreateFunctionUrlConfig(ctx, &lambda.CreateFunctionUrlConfigInput{
					FunctionName: aws.String(furlFuncName),
					AuthType:     types.FunctionUrlAuthTypeNone,
				})
				if err != nil {
					return err
				}
				if resp.FunctionUrl == nil || *resp.FunctionUrl == "" {
					return fmt.Errorf("FunctionUrl is nil or empty")
				}
				if resp.AuthType != types.FunctionUrlAuthTypeNone {
					return fmt.Errorf("AuthType mismatch, got %v", resp.AuthType)
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "GetFunctionUrlConfig", func() error {
				resp, err := client.GetFunctionUrlConfig(ctx, &lambda.GetFunctionUrlConfigInput{
					FunctionName: aws.String(furlFuncName),
				})
				if err != nil {
					return err
				}
				if resp.FunctionUrl == nil || *resp.FunctionUrl == "" {
					return fmt.Errorf("FunctionUrl is nil or empty")
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "UpdateFunctionUrlConfig", func() error {
				resp, err := client.UpdateFunctionUrlConfig(ctx, &lambda.UpdateFunctionUrlConfigInput{
					FunctionName: aws.String(furlFuncName),
					AuthType:     types.FunctionUrlAuthTypeAwsIam,
				})
				if err != nil {
					return err
				}
				if resp.AuthType != types.FunctionUrlAuthTypeAwsIam {
					return fmt.Errorf("AuthType not updated, got %v", resp.AuthType)
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "ListFunctionUrlConfigs", func() error {
				resp, err := client.ListFunctionUrlConfigs(ctx, &lambda.ListFunctionUrlConfigsInput{
					FunctionName: aws.String(furlFuncName),
				})
				if err != nil {
					return err
				}
				if resp.FunctionUrlConfigs == nil {
					return fmt.Errorf("url configs list is nil")
				}
				if len(resp.FunctionUrlConfigs) == 0 {
					return fmt.Errorf("expected at least 1 url config")
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "DeleteFunctionUrlConfig", func() error {
				_, err := client.DeleteFunctionUrlConfig(ctx, &lambda.DeleteFunctionUrlConfigInput{
					FunctionName: aws.String(furlFuncName),
				})
				return err
			}))
		}
	}

	// === GROUP F: INVOKE ASYNC & RESPONSE STREAM ===

	iaFuncName := fmt.Sprintf("IaFunc-%d", time.Now().UnixNano())
	iaRoleName := fmt.Sprintf("IaRole-%d", time.Now().UnixNano())
	iaRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", iaRoleName)
	iaCode := []byte("exports.handler = async () => { return { statusCode: 200 }; };")

	if err := createIAMRole(iaRoleName); err != nil {
		results = append(results, TestResult{Service: "lambda", TestName: "InvokeAsync_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
	} else {
		defer deleteIAMRole(iaRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(iaFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(iaRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: iaCode},
		})
		if err != nil {
			results = append(results, TestResult{Service: "lambda", TestName: "InvokeAsync_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create function: %v", err)})
		} else {
			defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(iaFuncName)})

			results = append(results, r.RunTest("lambda", "InvokeAsync", func() error {
				resp, err := client.InvokeAsync(ctx, &lambda.InvokeAsyncInput{
					FunctionName: aws.String(iaFuncName),
					InvokeArgs:   strings.NewReader(`{"test": true}`),
				})
				if err != nil {
					return err
				}
				if resp.Status != 202 {
					return fmt.Errorf("expected status 202, got %d", resp.Status)
				}
				return nil
			}))

			results = append(results, r.RunTest("lambda", "InvokeWithResponseStream", func() error {
				resp, err := client.InvokeWithResponseStream(ctx, &lambda.InvokeWithResponseStreamInput{
					FunctionName: aws.String(iaFuncName),
				})
				if err != nil {
					return err
				}
				if resp.StatusCode != 200 {
					return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
				}
				if resp.ResponseStreamContentType == nil {
					return fmt.Errorf("ResponseStreamContentType is nil")
				}
				return nil
			}))
		}
	}

	// === GROUP G: ERROR CASES ===

	results = append(results, r.RunTest("lambda", "CreateFunction_InvalidRuntime", func() error {
		invRtFuncName := fmt.Sprintf("InvRtFunc-%d", time.Now().UnixNano())
		invRtRoleName := fmt.Sprintf("InvRtRole-%d", time.Now().UnixNano())
		invRtRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", invRtRoleName)
		if err := createIAMRole(invRtRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(invRtRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(invRtFuncName),
			Runtime:      types.Runtime("invalid_runtime_99"),
			Role:         aws.String(invRtRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: []byte("code")},
		})
		if err == nil {
			client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(invRtFuncName)})
			return fmt.Errorf("expected error for invalid runtime")
		}
		var ipve *types.InvalidParameterValueException
		if !errors.As(err, &ipve) {
			return fmt.Errorf("expected InvalidParameterValueException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetAlias_NonExistent", func() error {
		_, err := client.GetAlias(ctx, &lambda.GetAliasInput{
			FunctionName: aws.String(functionName),
			Name:         aws.String("nonexistent-alias-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent alias")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetLayerVersion_NonExistent", func() error {
		_, err := client.GetLayerVersion(ctx, &lambda.GetLayerVersionInput{
			LayerName:     aws.String("nonexistent-layer-xyz"),
			VersionNumber: aws.Int64(999),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent layer version")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetFunctionUrlConfig_NoConfig", func() error {
		nofcFuncName := fmt.Sprintf("NofcFunc-%d", time.Now().UnixNano())
		nofcRoleName := fmt.Sprintf("NofcRole-%d", time.Now().UnixNano())
		nofcRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", nofcRoleName)
		if err := createIAMRole(nofcRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(nofcRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(nofcFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(nofcRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: []byte("code")},
		})
		if err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(nofcFuncName)})

		_, err = client.GetFunctionUrlConfig(ctx, &lambda.GetFunctionUrlConfigInput{
			FunctionName: aws.String(nofcFuncName),
		})
		if err == nil {
			return fmt.Errorf("expected error when no URL config set")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "PutFunctionEventInvokeConfig_NonExistent", func() error {
		_, err := client.PutFunctionEventInvokeConfig(ctx, &lambda.PutFunctionEventInvokeConfigInput{
			FunctionName:             aws.String("nonexistent-func-xyz-123"),
			MaximumEventAgeInSeconds: aws.Int32(3600),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent function")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("lambda", "ListFunctions_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgFuncs []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagFunc-%s-%d", pgTs, i)
			_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
				FunctionName: aws.String(name),
				Runtime:      types.RuntimeNodejs22x,
				Role:         aws.String(roleARN),
				Handler:      aws.String("index.handler"),
				Code:         &types.FunctionCode{ZipFile: functionCode},
			})
			if err != nil {
				return fmt.Errorf("create function %s: %v", name, err)
			}
			pgFuncs = append(pgFuncs, name)
		}

		var allFuncs []string
		var marker *string
		for {
			resp, err := client.ListFunctions(ctx, &lambda.ListFunctionsInput{
				Marker:   marker,
				MaxItems: aws.Int32(2),
			})
			if err != nil {
				for _, name := range pgFuncs {
					client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(name)})
				}
				return fmt.Errorf("list functions page: %v", err)
			}
			for _, f := range resp.Functions {
				if strings.HasPrefix(aws.ToString(f.FunctionName), "PagFunc-"+pgTs) {
					allFuncs = append(allFuncs, aws.ToString(f.FunctionName))
				}
			}
			if resp.NextMarker != nil && *resp.NextMarker != "" {
				marker = resp.NextMarker
			} else {
				break
			}
		}

		for _, name := range pgFuncs {
			client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(name)})
		}
		if len(allFuncs) != 5 {
			return fmt.Errorf("expected 5 paginated functions, got %d", len(allFuncs))
		}
		return nil
	}))

	return results
}
