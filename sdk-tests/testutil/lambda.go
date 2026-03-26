package testutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	ctx := context.Background()

	functionName := fmt.Sprintf("TestFunction-%d", time.Now().UnixNano())
	roleName := fmt.Sprintf("TestRole-%d", time.Now().UnixNano())
	roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)

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
		_, err := client.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(functionName),
			ZipFile:      newCode,
		})
		return err
	}))

	results = append(results, r.RunTest("lambda", "UpdateFunctionConfiguration", func() error {
		description := "Updated function"
		_, err := client.UpdateFunctionConfiguration(ctx, &lambda.UpdateFunctionConfigurationInput{
			FunctionName: aws.String(functionName),
			Description:  aws.String(description),
		})
		return err
	}))

	results = append(results, r.RunTest("lambda", "PublishVersion", func() error {
		_, err := client.PublishVersion(ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(functionName),
		})
		return err
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
		_, err := client.CreateAlias(ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(functionName),
			Name:            aws.String("live"),
			FunctionVersion: aws.String("$LATEST"),
		})
		return err
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
		_, err := client.UpdateAlias(ctx, &lambda.UpdateAliasInput{
			FunctionName: aws.String(functionName),
			Name:         aws.String("live"),
			Description:  aws.String("Production alias"),
		})
		return err
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
		_, err := client.PutFunctionConcurrency(ctx, &lambda.PutFunctionConcurrencyInput{
			FunctionName:                 aws.String(functionName),
			ReservedConcurrentExecutions: aws.Int32(10),
		})
		return err
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
		_, err := client.DeleteFunctionConcurrency(ctx, &lambda.DeleteFunctionConcurrencyInput{
			FunctionName: aws.String(functionName),
		})
		return err
	}))

	results = append(results, r.RunTest("lambda", "AddPermission", func() error {
		statementID := fmt.Sprintf("stmt-%d", time.Now().UnixNano())
		_, err := client.AddPermission(ctx, &lambda.AddPermissionInput{
			FunctionName: aws.String(functionName),
			StatementId:  aws.String(statementID),
			Action:       aws.String("lambda:InvokeFunction"),
			Principal:    aws.String("apigateway.amazonaws.com"),
		})
		return err
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
		_, err := client.UntagResource(ctx, &lambda.UntagResourceInput{
			Resource: aws.String(functionARN),
			TagKeys:  []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("lambda", "DeleteAlias", func() error {
		_, err := client.DeleteAlias(ctx, &lambda.DeleteAliasInput{
			FunctionName: aws.String(functionName),
			Name:         aws.String("live"),
		})
		return err
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
		_, err := client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
			FunctionName: aws.String(functionName),
		})
		return err
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
		dupRole := fmt.Sprintf("arn:aws:iam::000000000000:role/DupRole-%d", time.Now().UnixNano())
		dupCode := []byte("exports.handler = async () => { return 1; };")
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
		invRole := fmt.Sprintf("arn:aws:iam::000000000000:role/InvRole-%d", time.Now().UnixNano())
		invCode := []byte("exports.handler = async (event) => { return { statusCode: 200, body: JSON.stringify({result: 'ok'}) }; };")
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
		gfcRole := fmt.Sprintf("arn:aws:iam::000000000000:role/GfcRole-%d", time.Now().UnixNano())
		gfcCode := []byte("exports.handler = async () => { return 1; };")
		gfcDesc := "Test description for verification"
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
		pvRole := fmt.Sprintf("arn:aws:iam::000000000000:role/PvRole-%d", time.Now().UnixNano())
		pvCode := []byte("exports.handler = async () => { return 1; };")
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
		lfRole := fmt.Sprintf("arn:aws:iam::000000000000:role/LfRole-%d", time.Now().UnixNano())
		lfCode := []byte("exports.handler = async () => { return 1; };")
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
		caRole := fmt.Sprintf("arn:aws:iam::000000000000:role/CaRole-%d", time.Now().UnixNano())
		caCode := []byte("exports.handler = async () => { return 1; };")
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
		ucRole := fmt.Sprintf("arn:aws:iam::000000000000:role/UcRole-%d", time.Now().UnixNano())
		ucCode := []byte("exports.handler = async () => { return 1; };")
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

	return results
}
