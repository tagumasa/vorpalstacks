package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaFunctionTests(
	r *TestRunner,
	ctx context.Context,
	client *lambda.Client,
	createIAMRole func(string) error,
	deleteIAMRole func(string),
) []TestResult {
	var results []TestResult

	functionName := fmt.Sprintf("TestFunction-%d", time.Now().UnixNano())
	roleName := fmt.Sprintf("TestRole-%d", time.Now().UnixNano())
	roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)

	if err := createIAMRole(roleName); err != nil {
		return []TestResult{{Service: "lambda", TestName: "CreateFunction", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer deleteIAMRole(roleName)

	results = append(results, r.RunTest("lambda", "CreateFunction", func() error {
		resp, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(functionName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: []byte(lambdaFunctionCode)},
		})
		if err != nil {
			return err
		}
		if resp.FunctionName == nil || *resp.FunctionName != functionName {
			return fmt.Errorf("FunctionName mismatch, got %v", resp.FunctionName)
		}
		if resp.Runtime != types.RuntimeNodejs22x {
			return fmt.Errorf("Runtime mismatch, got %v", resp.Runtime)
		}
		if resp.Handler == nil || *resp.Handler != "index.handler" {
			return fmt.Errorf("Handler mismatch, got %v", resp.Handler)
		}
		if resp.Role == nil || *resp.Role != roleARN {
			return fmt.Errorf("Role mismatch, got %v", resp.Role)
		}
		if resp.CodeSize == 0 {
			return fmt.Errorf("CodeSize should be > 0")
		}
		if resp.CodeSha256 == nil || *resp.CodeSha256 == "" {
			return fmt.Errorf("CodeSha256 is nil or empty")
		}
		return nil
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
		if resp.Configuration.FunctionName == nil || *resp.Configuration.FunctionName != functionName {
			return fmt.Errorf("FunctionName mismatch, got %v", resp.Configuration.FunctionName)
		}
		if resp.Configuration.Runtime != types.RuntimeNodejs22x {
			return fmt.Errorf("Runtime mismatch, got %v", resp.Configuration.Runtime)
		}
		if resp.Configuration.Handler == nil || *resp.Configuration.Handler != "index.handler" {
			return fmt.Errorf("Handler mismatch, got %v", resp.Configuration.Handler)
		}
		if resp.Configuration.FunctionArn == nil || *resp.Configuration.FunctionArn == "" {
			return fmt.Errorf("FunctionArn is nil or empty")
		}
		if resp.Code == nil || resp.Code.Location == nil {
			return fmt.Errorf("Code.Location is nil")
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
		if resp.FunctionName == nil || *resp.FunctionName != functionName {
			return fmt.Errorf("FunctionName mismatch, got %v", resp.FunctionName)
		}
		if resp.Runtime != types.RuntimeNodejs22x {
			return fmt.Errorf("Runtime mismatch, got %v", resp.Runtime)
		}
		if resp.Handler == nil || *resp.Handler != "index.handler" {
			return fmt.Errorf("Handler mismatch, got %v", resp.Handler)
		}
		if resp.Role == nil || *resp.Role != roleARN {
			return fmt.Errorf("Role mismatch, got %v", resp.Role)
		}
		if resp.CodeSha256 == nil || *resp.CodeSha256 == "" {
			return fmt.Errorf("CodeSha256 is nil or empty")
		}
		if resp.State != types.StateActive && resp.State != types.StatePending {
			return fmt.Errorf("State should be Active or Pending, got %v", resp.State)
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
		found := false
		for _, f := range resp.Functions {
			if f.FunctionName != nil && *f.FunctionName == functionName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created function %s not found in ListFunctions", functionName)
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
		if resp.FunctionName == nil || *resp.FunctionName != functionName {
			return fmt.Errorf("FunctionName mismatch, got %v", resp.FunctionName)
		}
		if resp.CodeSha256 == nil || *resp.CodeSha256 == "" {
			return fmt.Errorf("CodeSha256 is nil or empty")
		}
		if resp.LastModified == nil || *resp.LastModified == "" {
			return fmt.Errorf("LastModified is nil or empty")
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
		if resp.Description == nil || *resp.Description != description {
			return fmt.Errorf("Description mismatch, got %v", resp.Description)
		}
		if resp.FunctionName == nil || *resp.FunctionName != functionName {
			return fmt.Errorf("FunctionName mismatch, got %v", resp.FunctionName)
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
		if resp.StatusCode != 200 {
			return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if resp.ExecutedVersion == nil || *resp.ExecutedVersion == "" {
			return fmt.Errorf("ExecutedVersion is nil or empty")
		}
		if len(resp.Payload) == 0 {
			return fmt.Errorf("Payload should not be empty")
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
			return fmt.Errorf("ReservedConcurrentExecutions is nil")
		}
		if *resp.ReservedConcurrentExecutions != 10 {
			return fmt.Errorf("ReservedConcurrentExecutions mismatch, expected 10, got %d", *resp.ReservedConcurrentExecutions)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "DeleteFunctionConcurrency", func() error {
		_, err := client.DeleteFunctionConcurrency(ctx, &lambda.DeleteFunctionConcurrencyInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		concurrencyResp, err := client.GetFunctionConcurrency(ctx, &lambda.GetFunctionConcurrencyInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return fmt.Errorf("GetFunctionConcurrency after delete: %v", err)
		}
		if concurrencyResp.ReservedConcurrentExecutions != nil {
			return fmt.Errorf("ReservedConcurrentExecutions should be nil after delete, got %d", *concurrencyResp.ReservedConcurrentExecutions)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetAccountSettings", func() error {
		resp, err := client.GetAccountSettings(ctx, &lambda.GetAccountSettingsInput{})
		if err != nil {
			return err
		}
		if resp.AccountLimit == nil {
			return fmt.Errorf("AccountLimit is nil")
		}
		if resp.AccountUsage == nil {
			return fmt.Errorf("AccountUsage is nil")
		}
		if resp.AccountLimit.ConcurrentExecutions == 0 {
			return fmt.Errorf("ConcurrentExecutions limit should be > 0")
		}
		if resp.AccountLimit.TotalCodeSize == 0 {
			return fmt.Errorf("TotalCodeSize limit should be > 0")
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "DeleteFunction", func() error {
		_, err := client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		_, err = client.GetFunction(ctx, &lambda.GetFunctionInput{
			FunctionName: aws.String(functionName),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// === ERROR CASES ===

	results = append(results, r.RunTest("lambda", "GetFunction_NonExistent", func() error {
		_, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "Invoke_NonExistent", func() error {
		_, err := client.Invoke(ctx, &lambda.InvokeInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "UpdateFunctionCode_NonExistent", func() error {
		_, err := client.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
			ZipFile:      []byte("code"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "DeleteFunction_NonExistent", func() error {
		_, err := client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
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
		if err := AssertErrorContains(err, "ResourceConflictException"); err != nil {
			return err
		}
		return nil
	}))

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
		if err := AssertErrorContains(err, "InvalidParameterValueException"); err != nil {
			return err
		}
		return nil
	}))

	// === VERIFICATION TESTS ===

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
		var lambdaResult map[string]interface{}
		if err := json.Unmarshal(payload, &lambdaResult); err != nil {
			return fmt.Errorf("parse payload JSON: %v (payload: %s)", err, string(payload))
		}
		if lambdaResult["statusCode"] != float64(200) {
			return fmt.Errorf("expected statusCode 200, got %v", lambdaResult["statusCode"])
		}
		bodyStr, ok := lambdaResult["body"].(string)
		if !ok || bodyStr == "" {
			return fmt.Errorf("expected non-empty body string in payload, got %v", lambdaResult["body"])
		}
		var body map[string]interface{}
		if err := json.Unmarshal([]byte(bodyStr), &body); err != nil {
			return fmt.Errorf("parse body JSON: %v", err)
		}
		if body["result"] != "ok" {
			return fmt.Errorf("body result mismatch: got %v", body["result"])
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
				Code:         &types.FunctionCode{ZipFile: []byte(lambdaFunctionCode)},
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
