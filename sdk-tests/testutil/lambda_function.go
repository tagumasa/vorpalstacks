package testutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"vorpalstacks-sdk-tests/config"
)

func runLambdaFunctionTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	functionName := tc.unique("TestFunction")
	roleARN, cleanupRole, err := tc.createRole(tc.unique("TestRole"))
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "CreateFunction", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer cleanupRole()
	defer deleteLambdaLogGroup(tc.cwl, tc.ctx, functionName)
	var functionARN string

	results = append(results, tc.r.RunTest("lambda", "CreateFunction", func() error {
		zipCode, err := zipLambdaCode(lambdaFunctionCode)
		if err != nil {
			return fmt.Errorf("zip lambda code: %w", err)
		}
		resp, err := tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(functionName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: zipCode},
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
		functionARN = aws.ToString(resp.FunctionArn)
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetFunction", func() error {
		resp, err := tc.client.GetFunction(tc.ctx, &lambda.GetFunctionInput{
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

	// FunctionName accepts the bare name, the full ARN, and the partial
	// ARN (account-id:function:name) per the API reference; all three must
	// resolve to the same function.
	results = append(results, tc.r.RunTest("lambda", "GetFunction_ArnForms", func() error {
		if functionARN == "" {
			return fmt.Errorf("function ARN not captured from CreateFunction")
		}
		for _, ref := range []string{functionARN, fmt.Sprintf("%s:function:%s", tc.r.accountID, functionName)} {
			resp, err := tc.client.GetFunction(tc.ctx, &lambda.GetFunctionInput{
				FunctionName: aws.String(ref),
			})
			if err != nil {
				return fmt.Errorf("GetFunction with %q: %v", ref, err)
			}
			if aws.ToString(resp.Configuration.FunctionName) != functionName {
				return fmt.Errorf("FunctionName mismatch for %q: got %v", ref, resp.Configuration.FunctionName)
			}
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetFunctionConfiguration", func() error {
		resp, err := tc.client.GetFunctionConfiguration(tc.ctx, &lambda.GetFunctionConfigurationInput{
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

	results = append(results, tc.r.RunTest("lambda", "ListFunctions", func() error {
		lfFunc, cleanupLfFn, err := tc.setupFunction("LfFunc", "exports.handler = async () => { return 1; };")
		if err != nil {
			return err
		}
		defer cleanupLfFn()

		all, err := paginate[types.FunctionConfiguration](func(next *string) ([]types.FunctionConfiguration, *string, error) {
			resp, err := tc.client.ListFunctions(tc.ctx, &lambda.ListFunctionsInput{Marker: next})
			if err != nil {
				return nil, nil, err
			}
			return resp.Functions, resp.NextMarker, nil
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		for _, f := range all {
			if aws.ToString(f.FunctionName) == lfFunc {
				if f.Runtime != types.RuntimeNodejs22x {
					return fmt.Errorf("runtime mismatch in list, got %v", f.Runtime)
				}
				if aws.ToString(f.Handler) != "index.handler" {
					return fmt.Errorf("handler mismatch in list, got %v", f.Handler)
				}
				return nil
			}
		}
		return fmt.Errorf("created function %s not found in ListFunctions", lfFunc)
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateFunctionCode", func() error {
		newCode, err := zipLambdaCode("exports.handler = async (event) => { return { statusCode: 200, body: 'Updated' }; };")
		if err != nil {
			return fmt.Errorf("zip lambda code: %w", err)
		}
		resp, err := tc.client.UpdateFunctionCode(tc.ctx, &lambda.UpdateFunctionCodeInput{
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

	results = append(results, tc.r.RunTest("lambda", "FunctionCode_S3ObjectVersion", func() error {
		cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: tc.r.endpoint,
			Region:   tc.r.region,
		})
		if err != nil {
			return fmt.Errorf("load config: %v", err)
		}
		s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })

		suffix := time.Now().UnixNano()
		bucket := fmt.Sprintf("lambda-versioned-code-%d", suffix)
		key := "code.zip"
		if _, err := s3Client.CreateBucket(tc.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer func() {
			s3Client.DeleteObject(tc.ctx, &s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
			s3Client.DeleteBucket(tc.ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})
		}()
		if _, err := s3Client.PutBucketVersioning(tc.ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucket),
			VersioningConfiguration: &s3types.VersioningConfiguration{
				Status: s3types.BucketVersioningStatusEnabled,
			},
		}); err != nil {
			return fmt.Errorf("enable versioning: %v", err)
		}

		zip1, err := zipLambdaCode("exports.handler = async () => 'v1';")
		if err != nil {
			return fmt.Errorf("zip v1: %v", err)
		}
		zip2, err := zipLambdaCode("exports.handler = async () => 'v2';")
		if err != nil {
			return fmt.Errorf("zip v2: %v", err)
		}
		put1, err := s3Client.PutObject(tc.ctx, &s3.PutObjectInput{Bucket: aws.String(bucket), Key: aws.String(key), Body: bytes.NewReader(zip1)})
		if err != nil {
			return fmt.Errorf("put v1: %v", err)
		}
		put2, err := s3Client.PutObject(tc.ctx, &s3.PutObjectInput{Bucket: aws.String(bucket), Key: aws.String(key), Body: bytes.NewReader(zip2)})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}
		version1 := aws.ToString(put1.VersionId)
		version2 := aws.ToString(put2.VersionId)
		if version1 == "" || version2 == "" || version1 == version2 {
			return fmt.Errorf("expected two distinct version ids, got %q and %q", version1, version2)
		}
		sum1 := sha256.Sum256(zip1)
		sum2 := sha256.Sum256(zip2)
		hash1 := base64.StdEncoding.EncodeToString(sum1[:])
		hash2 := base64.StdEncoding.EncodeToString(sum2[:])

		fnName := fmt.Sprintf("VersionedCodeFn-%d", suffix)
		created, err := tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(fnName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(roleARN),
			Handler:      aws.String("index.handler"),
			Code: &types.FunctionCode{
				S3Bucket:        aws.String(bucket),
				S3Key:           aws.String(key),
				S3ObjectVersion: aws.String(version1),
			},
		})
		if err != nil {
			return fmt.Errorf("create function pinned to version 1: %v", err)
		}
		defer tc.deleteFunctionAndLogs(fnName)
		if aws.ToString(created.CodeSha256) != hash1 {
			return fmt.Errorf("CodeSHA256 = %s, want version-1 hash %s (latest is version 2)", aws.ToString(created.CodeSha256), hash1)
		}

		updated, err := tc.client.UpdateFunctionCode(tc.ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName:    aws.String(fnName),
			S3Bucket:        aws.String(bucket),
			S3Key:           aws.String(key),
			S3ObjectVersion: aws.String(version2),
		})
		if err != nil {
			return fmt.Errorf("update function code pinned to version 2: %v", err)
		}
		if aws.ToString(updated.CodeSha256) != hash2 {
			return fmt.Errorf("updated CodeSHA256 = %s, want version-2 hash %s", aws.ToString(updated.CodeSha256), hash2)
		}

		layerName := fmt.Sprintf("versioned-code-layer-%d", suffix)
		layer, err := tc.client.PublishLayerVersion(tc.ctx, &lambda.PublishLayerVersionInput{
			LayerName: aws.String(layerName),
			Content: &types.LayerVersionContentInput{
				S3Bucket:        aws.String(bucket),
				S3Key:           aws.String(key),
				S3ObjectVersion: aws.String(version1),
			},
		})
		if err != nil {
			return fmt.Errorf("publish layer pinned to version 1: %v", err)
		}
		defer tc.client.DeleteLayerVersion(tc.ctx, &lambda.DeleteLayerVersionInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(layer.Version),
		})
		if layer.Content == nil || aws.ToString(layer.Content.CodeSha256) != hash1 {
			return fmt.Errorf("layer CodeSHA256 = %v, want version-1 hash %s", layer.Content, hash1)
		}

		if _, err := tc.client.UpdateFunctionCode(tc.ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName:    aws.String(fnName),
			S3Bucket:        aws.String(bucket),
			S3Key:           aws.String(key),
			S3ObjectVersion: aws.String("0000000000000000000000000000000000000000"),
		}); err == nil {
			return fmt.Errorf("a nonexistent object version must be rejected")
		} else if err := expectAWSErrorCode(err, "InvalidParameterValueException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateFunctionConfiguration", func() error {
		description := "Updated function"
		resp, err := tc.client.UpdateFunctionConfiguration(tc.ctx, &lambda.UpdateFunctionConfigurationInput{
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

	results = append(results, tc.r.RunTest("lambda", "Environment_InvalidKeyRejected", func() error {
		envFnName := tc.unique("EnvVarFn")
		code, err := zipLambdaCode("exports.handler = async () => 'ok';")
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}

		// "Keys start with a letter and are at least two characters. Keys
		// only contain letters, numbers, and the underscore character (_)."
		for _, key := range []string{"1BAD", "a-b", "a", "_x"} {
			_, err := tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
				FunctionName: aws.String(envFnName),
				Runtime:      types.RuntimeNodejs22x,
				Role:         aws.String(roleARN),
				Handler:      aws.String("index.handler"),
				Code:         &types.FunctionCode{ZipFile: code},
				Environment: &types.Environment{
					Variables: map[string]string{key: "v"},
				},
			})
			if err == nil {
				return fmt.Errorf("environment key %q must be rejected", key)
			}
			if err := expectAWSErrorCode(err, "InvalidParameterValueException"); err != nil {
				return err
			}
		}

		validName := envFnName
		_, cleanupValidFn, err := tc.createFunction(validName, roleARN, "exports.handler = async () => 'ok';",
			func(input *lambda.CreateFunctionInput) {
				input.Environment = &types.Environment{
					Variables: map[string]string{"GOOD_KEY": "v"},
				}
			})
		if err != nil {
			return fmt.Errorf("valid environment key rejected: %v", err)
		}
		defer cleanupValidFn()
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke", func() error {
		resp, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
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

	results = append(results, tc.r.RunTest("lambda", "PutFunctionConcurrency", func() error {
		_, err := tc.client.PutFunctionConcurrency(tc.ctx, &lambda.PutFunctionConcurrencyInput{
			FunctionName:                 aws.String(functionName),
			ReservedConcurrentExecutions: aws.Int32(10),
		})
		return err
	}))

	results = append(results, tc.r.RunTest("lambda", "GetFunctionConcurrency", func() error {
		resp, err := tc.client.GetFunctionConcurrency(tc.ctx, &lambda.GetFunctionConcurrencyInput{
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

	results = append(results, tc.r.RunTest("lambda", "DeleteFunctionConcurrency", func() error {
		_, err := tc.client.DeleteFunctionConcurrency(tc.ctx, &lambda.DeleteFunctionConcurrencyInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		// Once the limit is removed, the concurrency sub-resource no longer
		// exists and AWS answers GetFunctionConcurrency with 404
		// ResourceNotFoundException.
		_, err = tc.client.GetFunctionConcurrency(tc.ctx, &lambda.GetFunctionConcurrencyInput{
			FunctionName: aws.String(functionName),
		})
		if err := expectAWSErrorCode(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("GetFunctionConcurrency after delete: %v", err)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateFunctionCode_RevisionIdPrecondition", func() error {
		revCode, err := zipLambdaCode("exports.handler = async () => { return 2; };")
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		_, err = tc.client.UpdateFunctionCode(tc.ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(functionName),
			ZipFile:      revCode,
			RevisionId:   aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err := expectAWSErrorCode(err, "ResourceConflictException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetFunctionConfiguration_UnqualifiedVersionIsLatest", func() error {
		if _, err := tc.client.PublishVersion(tc.ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(functionName),
		}); err != nil {
			return fmt.Errorf("publish: %v", err)
		}
		// A request without a qualifier addresses the unpublished version,
		// which AWS always reports as $LATEST even after publishing.
		resp, err := tc.client.GetFunctionConfiguration(tc.ctx, &lambda.GetFunctionConfigurationInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp.Version == nil || *resp.Version != "$LATEST" {
			return fmt.Errorf("unqualified Version should be $LATEST, got %v", resp.Version)
		}
		// The published version itself reports its own number.
		qualified, err := tc.client.GetFunctionConfiguration(tc.ctx, &lambda.GetFunctionConfigurationInput{
			FunctionName: aws.String(functionName),
			Qualifier:    aws.String("1"),
		})
		if err != nil {
			return err
		}
		if qualified.Version == nil || *qualified.Version != "1" {
			return fmt.Errorf("qualified Version should be 1, got %v", qualified.Version)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetFunction_IncludesConcurrency", func() error {
		if _, err := tc.client.PutFunctionConcurrency(tc.ctx, &lambda.PutFunctionConcurrencyInput{
			FunctionName:                 aws.String(functionName),
			ReservedConcurrentExecutions: aws.Int32(5),
		}); err != nil {
			return fmt.Errorf("put concurrency: %v", err)
		}
		resp, err := tc.client.GetFunction(tc.ctx, &lambda.GetFunctionInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		if resp.Concurrency == nil {
			return fmt.Errorf("Concurrency is nil")
		}
		if resp.Concurrency.ReservedConcurrentExecutions == nil || *resp.Concurrency.ReservedConcurrentExecutions != 5 {
			return fmt.Errorf("ReservedConcurrentExecutions mismatch, got %v", resp.Concurrency.ReservedConcurrentExecutions)
		}

		// Unreserved concurrency reflects the reservation just configured.
		settings, err := tc.client.GetAccountSettings(tc.ctx, &lambda.GetAccountSettingsInput{})
		if err != nil {
			return err
		}
		if settings.AccountLimit == nil || settings.AccountLimit.UnreservedConcurrentExecutions == nil {
			return fmt.Errorf("UnreservedConcurrentExecutions is nil")
		}
		if *settings.AccountLimit.UnreservedConcurrentExecutions != 995 {
			return fmt.Errorf("UnreservedConcurrentExecutions should be 995 after reserving 5, got %d", *settings.AccountLimit.UnreservedConcurrentExecutions)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetFunctionConfiguration_ImageConfigResponse", func() error {
		imgFunc := tc.unique("ImageFunc")
		created, err := tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(imgFunc),
			Role:         aws.String(roleARN),
			PackageType:  types.PackageTypeImage,
			Code:         &types.FunctionCode{ImageUri: aws.String("123456789012.dkr.ecr.us-east-1.amazonaws.com/test/image:1")},
			ImageConfig: &types.ImageConfig{
				Command: []string{"/bin/app"},
			},
		})
		if err != nil {
			return fmt.Errorf("create image function: %v", err)
		}
		defer tc.deleteFunctionAndLogs(imgFunc)
		if created.ImageConfigResponse == nil {
			return fmt.Errorf("create response ImageConfigResponse is nil")
		}

		resp, err := tc.client.GetFunctionConfiguration(tc.ctx, &lambda.GetFunctionConfigurationInput{
			FunctionName: aws.String(imgFunc),
		})
		if err != nil {
			return err
		}
		// The wire member is ImageConfigResponse nesting ImageConfig; a
		// flat ImageConfig member would leave this nil.
		if resp.ImageConfigResponse == nil || resp.ImageConfigResponse.ImageConfig == nil {
			return fmt.Errorf("ImageConfigResponse.ImageConfig is nil")
		}
		if len(resp.ImageConfigResponse.ImageConfig.Command) != 1 || resp.ImageConfigResponse.ImageConfig.Command[0] != "/bin/app" {
			return fmt.Errorf("ImageConfig.Command mismatch, got %v", resp.ImageConfigResponse.ImageConfig.Command)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "PublishedVersion_ExecutesSnapshotCode", func() error {
		snapFunc, cleanupSnapFn, err := tc.setupFunction("SnapFunc", "exports.handler = async () => 'v1-output';")
		if err != nil {
			return err
		}
		defer cleanupSnapFn()

		if _, err := tc.client.PublishVersion(tc.ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(snapFunc),
		}); err != nil {
			return fmt.Errorf("publish: %v", err)
		}

		// Replace the $LATEST code after publishing; the published version
		// must keep executing its own snapshot.
		v2Code, err := zipLambdaCode("exports.handler = async () => 'v2-output';")
		if err != nil {
			return fmt.Errorf("zip v2 code: %v", err)
		}
		if _, err := tc.client.UpdateFunctionCode(tc.ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(snapFunc),
			ZipFile:      v2Code,
		}); err != nil {
			return fmt.Errorf("update code: %v", err)
		}

		versioned, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(snapFunc),
			Qualifier:    aws.String("1"),
		})
		if err != nil {
			return err
		}
		if versioned.ExecutedVersion == nil || *versioned.ExecutedVersion != "1" {
			return fmt.Errorf("ExecutedVersion mismatch, got %v", versioned.ExecutedVersion)
		}
		if !strings.Contains(string(versioned.Payload), "v1-output") {
			return fmt.Errorf("published version should execute the v1 snapshot, got payload %q", string(versioned.Payload))
		}

		latest, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(snapFunc),
		})
		if err != nil {
			return err
		}
		if !strings.Contains(string(latest.Payload), "v2-output") {
			return fmt.Errorf("unqualified invoke should execute the updated $LATEST code, got payload %q", string(latest.Payload))
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_FunctionError_Classification", func() error {
		feFunc, cleanupFeFn, err := tc.setupFunction("FeFunc", "exports.handler = async () => { throw new Error('boom'); };")
		if err != nil {
			return err
		}
		defer cleanupFeFn()

		// A thrown error is intercepted by the runtime: Unhandled.
		unhandled, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(feFunc),
		})
		if err != nil {
			return err
		}
		if unhandled.FunctionError == nil || *unhandled.FunctionError != "Unhandled" {
			return fmt.Errorf("thrown error should classify as Unhandled, got %v", unhandled.FunctionError)
		}
		// The intercepted failure reaches the caller as the runtime error
		// document, not an empty payload.
		var errDoc map[string]interface{}
		if err := json.Unmarshal(unhandled.Payload, &errDoc); err != nil {
			return fmt.Errorf("unhandled payload must be the error document JSON, got %q: %v", string(unhandled.Payload), err)
		}
		if errDoc["errorMessage"] != "boom" {
			return fmt.Errorf("error document errorMessage mismatch: got %v", errDoc["errorMessage"])
		}
		if errDoc["errorType"] != "Error" {
			return fmt.Errorf("error document errorType mismatch: got %v", errDoc["errorType"])
		}
		trace, ok := errDoc["stackTrace"].([]interface{})
		if !ok || len(trace) == 0 {
			return fmt.Errorf("error document stackTrace must be non-empty, got %v", errDoc["stackTrace"])
		}

		// A returned error-shaped document is a plain success on AWS: the
		// FunctionError header may only be present when an error actually
		// occurred, so a handler returning {errorMessage: ...} carries no
		// header and the document round-trips as the payload.
		returnedCode, err := zipLambdaCode("exports.handler = async () => ({ errorMessage: 'handled-boom', errorType: 'Error' });")
		if err != nil {
			return fmt.Errorf("zip returned-envelope code: %v", err)
		}
		if _, err := tc.client.UpdateFunctionCode(tc.ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(feFunc),
			ZipFile:      returnedCode,
		}); err != nil {
			return fmt.Errorf("update code: %v", err)
		}
		returned, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(feFunc),
		})
		if err != nil {
			return err
		}
		if returned.FunctionError != nil {
			return fmt.Errorf("returned error document must not set FunctionError, got %v", *returned.FunctionError)
		}
		if !strings.Contains(string(returned.Payload), "handled-boom") {
			return fmt.Errorf("returned error document must round-trip as the payload, got %q", string(returned.Payload))
		}

		// A failure the handler signalled through the callback is a
		// Handled error: the runtime reports callback-signalled errors as
		// handled, unlike uncaught failures.
		cbCode, err := zipLambdaCode("exports.handler = (event, context, callback) => { callback(new Error('cb-boom')); };")
		if err != nil {
			return fmt.Errorf("zip callback-error code: %v", err)
		}
		if _, err := tc.client.UpdateFunctionCode(tc.ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(feFunc),
			ZipFile:      cbCode,
		}); err != nil {
			return fmt.Errorf("update code: %v", err)
		}
		cbHandled, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(feFunc),
		})
		if err != nil {
			return err
		}
		if cbHandled.FunctionError == nil || *cbHandled.FunctionError != "Handled" {
			return fmt.Errorf("callback error should classify as Handled, got %v", cbHandled.FunctionError)
		}
		var cbDoc map[string]interface{}
		if err := json.Unmarshal(cbHandled.Payload, &cbDoc); err != nil {
			return fmt.Errorf("callback error payload must be the error document JSON, got %q: %v", string(cbHandled.Payload), err)
		}
		if cbDoc["errorMessage"] != "cb-boom" {
			return fmt.Errorf("callback error document errorMessage mismatch: got %v", cbDoc["errorMessage"])
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_HandledErrorWithLogs", func() error {
		logFunc, cleanupLogFn, err := tc.setupFunction("LogFeFunc",
			"exports.handler = (event, context, callback) => { console.log('about to fail'); callback(new Error('logged-boom')); };")
		if err != nil {
			return err
		}
		defer cleanupLogFn()

		// Handler console output reaches stdout before the wrapper appends
		// the returned payload, so logs and payload stay separate channels.
		// A callback-signalled failure after logging is a Handled error and
		// the error document still reaches the caller as the payload.
		result, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(logFunc),
		})
		if err != nil {
			return err
		}
		if result.FunctionError == nil || *result.FunctionError != "Handled" {
			return fmt.Errorf("callback error after log lines should classify as Handled, got %v", result.FunctionError)
		}
		if !strings.Contains(string(result.Payload), "logged-boom") {
			return fmt.Errorf("payload should carry the error document after the log lines, got %q", string(result.Payload))
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_PayloadExcludesConsoleOutput", func() error {
		logFunc, cleanupLogFn, err := tc.setupFunction("PayloadCleanFn",
			"exports.handler = async () => { console.log('log-line-should-not-leak'); return { ok: true, value: 7 }; };")
		if err != nil {
			return err
		}
		defer cleanupLogFn()

		// On AWS the response payload is the handler's return value only;
		// console output belongs to the logs. A handler that logs before
		// returning must not leak its log lines into the payload.
		result, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(logFunc),
		})
		if err != nil {
			return err
		}
		if result.FunctionError != nil {
			return fmt.Errorf("a returning handler is a success, got FunctionError=%v", *result.FunctionError)
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(result.Payload, &payload); err != nil {
			return fmt.Errorf("payload must be exactly the returned object, got %q: %v", string(result.Payload), err)
		}
		if payload["ok"] != true || payload["value"] != float64(7) {
			return fmt.Errorf("payload must carry the returned members, got %q", string(result.Payload))
		}
		if strings.Contains(string(result.Payload), "log-line-should-not-leak") {
			return fmt.Errorf("console output must not leak into the response payload, got %q", string(result.Payload))
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ListFunctions_AllPaginatesEntries", func() error {
		listFunc, cleanupListFn, err := tc.setupFunction("ListAllFn", "exports.handler = async () => ({});")
		if err != nil {
			return err
		}
		defer cleanupListFn()

		for i := 0; i < 3; i++ {
			if _, err := tc.client.PublishVersion(tc.ctx, &lambda.PublishVersionInput{
				FunctionName: aws.String(listFunc),
			}); err != nil {
				return fmt.Errorf("publish version: %v", err)
			}
		}

		// The function contributes four entries ($LATEST plus three
		// published versions); with MaxItems=2 every page must carry at
		// most two entries and the walk must see each entry exactly once.
		seen := map[string]int{}
		marker := ""
		for page := 0; ; page++ {
			input := &lambda.ListFunctionsInput{
				FunctionVersion: types.FunctionVersionAll,
				MaxItems:        aws.Int32(2),
			}
			if marker != "" {
				input.Marker = aws.String(marker)
			}
			out, err := tc.client.ListFunctions(tc.ctx, input)
			if err != nil {
				return fmt.Errorf("list functions: %v", err)
			}
			if len(out.Functions) > 2 {
				return fmt.Errorf("page %d carries %d entries, the documented cap is 2", page, len(out.Functions))
			}
			for _, f := range out.Functions {
				if aws.ToString(f.FunctionName) == listFunc {
					seen[aws.ToString(f.Version)]++
				}
			}
			if out.NextMarker == nil || *out.NextMarker == "" {
				break
			}
			marker = *out.NextMarker
			if page > 50 {
				return fmt.Errorf("pagination did not terminate")
			}
		}

		wantVersions := map[string]int{"$LATEST": 1, "1": 1, "2": 1, "3": 1}
		if len(seen) != len(wantVersions) {
			return fmt.Errorf("expected entries %v for the function, got %v", wantVersions, seen)
		}
		for v, count := range wantVersions {
			if seen[v] != count {
				return fmt.Errorf("version %q seen %d times, want %d", v, seen[v], count)
			}
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetAccountSettings", func() error {
		resp, err := tc.client.GetAccountSettings(tc.ctx, &lambda.GetAccountSettingsInput{})
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

	results = append(results, tc.r.RunTest("lambda", "DeleteFunction", func() error {
		_, err := tc.client.DeleteFunction(tc.ctx, &lambda.DeleteFunctionInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetFunction(tc.ctx, &lambda.GetFunctionInput{
			FunctionName: aws.String(functionName),
		})
		if err := expectAWSErrorCode(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// === ERROR CASES ===

	// Every operation that addresses a function by name must answer
	// ResourceNotFoundException for a name that does not exist.
	results = append(results, tc.r.RunTest("lambda", "NonExistentFunctionReference", func() error {
		const missing = "NoSuchFunction_xyz_12345"
		for _, c := range []struct {
			op   string
			call func() error
		}{
			{"GetFunction", func() error {
				_, err := tc.client.GetFunction(tc.ctx, &lambda.GetFunctionInput{
					FunctionName: aws.String(missing),
				})
				return err
			}},
			{"Invoke", func() error {
				_, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
					FunctionName: aws.String(missing),
				})
				return err
			}},
			{"UpdateFunctionCode", func() error {
				_, err := tc.client.UpdateFunctionCode(tc.ctx, &lambda.UpdateFunctionCodeInput{
					FunctionName: aws.String(missing),
					ZipFile:      []byte("code"),
				})
				return err
			}},
			{"DeleteFunction", func() error {
				_, err := tc.client.DeleteFunction(tc.ctx, &lambda.DeleteFunctionInput{
					FunctionName: aws.String(missing),
				})
				return err
			}},
		} {
			if err := expectAWSErrorCode(c.call(), "ResourceNotFoundException"); err != nil {
				return fmt.Errorf("%s: %w", c.op, err)
			}
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "CreateFunction_DuplicateName", func() error {
		dupName := tc.unique("DupFunc")
		dupRole, cleanupDupRole, err := tc.createRole(tc.unique("DupRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupDupRole()
		dupCode, err := zipLambdaCode("exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("zip lambda code: %w", err)
		}
		_, err = tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(dupName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(dupRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: dupCode},
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.deleteFunctionAndLogs(dupName)

		_, err = tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(dupName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(dupRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: dupCode},
		})
		if err := expectAWSErrorCode(err, "ResourceConflictException"); err != nil {
			return err
		}
		return nil
	}))

	// CreateFunction rejects configurations the requested runtime does not
	// support: an unknown runtime value, and SnapStart on a runtime that
	// cannot snapshot.
	results = append(results, tc.r.RunTest("lambda", "CreateFunction_RejectsUnsupportedConfiguration", func() error {
		rejRole, cleanupRejRole, err := tc.createRole(tc.unique("RejRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupRejRole()
		for _, c := range []struct {
			name   string
			mutate func(*lambda.CreateFunctionInput)
		}{
			{"InvalidRuntime", func(in *lambda.CreateFunctionInput) {
				in.Runtime = types.Runtime("invalid_runtime_99")
			}},
			{"SnapStartUnsupportedRuntime", func(in *lambda.CreateFunctionInput) {
				in.SnapStart = &types.SnapStart{ApplyOn: types.SnapStartApplyOnPublishedVersions}
			}},
		} {
			input := &lambda.CreateFunctionInput{
				FunctionName: aws.String(tc.unique("RejFunc")),
				Runtime:      types.RuntimeNodejs22x,
				Role:         aws.String(rejRole),
				Handler:      aws.String("index.handler"),
				Code:         &types.FunctionCode{ZipFile: []byte("code")},
			}
			c.mutate(input)
			if _, err := tc.client.CreateFunction(tc.ctx, input); err != nil {
				if err := expectAWSErrorCode(err, "InvalidParameterValueException"); err != nil {
					return fmt.Errorf("%s: %w", c.name, err)
				}
				continue
			}
			return fmt.Errorf("%s must be rejected with InvalidParameterValueException", c.name)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "CreateFunction_JavaAl2023Runtime", func() error {
		alRtFuncName := tc.unique("AlRtFunc")
		alRtRole, cleanupAlRtRole, err := tc.createRole(tc.unique("AlRtRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupAlRtRole()
		alRtCode, err := zipLambdaCode("exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		resp, err := tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(alRtFuncName),
			Role:         aws.String(alRtRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: alRtCode},
			Runtime:      types.RuntimeJava17al2023,
		})
		if err != nil {
			return err
		}
		defer tc.deleteFunctionAndLogs(alRtFuncName)
		if resp.Runtime != types.RuntimeJava17al2023 {
			return fmt.Errorf("Runtime mismatch, got %v", resp.Runtime)
		}
		return nil
	}))

	// === VERIFICATION TESTS ===

	results = append(results, tc.r.RunTest("lambda", "Invoke_VerifyResponsePayload", func() error {
		invFunc, cleanupInvFn, err := tc.setupFunction("InvFunc",
			"exports.handler = async (event) => { return { statusCode: 200, body: JSON.stringify({result: 'ok'}) }; };")
		if err != nil {
			return err
		}
		defer cleanupInvFn()

		resp, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
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

	results = append(results, tc.r.RunTest("lambda", "Invoke_AsyncPayloadLimit", func() error {
		apFunc, cleanupApFn, err := tc.setupFunction("AsyncLimitFn", "exports.handler = async () => ({ ok: 1 });")
		if err != nil {
			return err
		}
		defer cleanupApFn()

		// 300 KB stays within the 1 MB asynchronous invocation limit and
		// must be accepted with 202.
		medium := []byte(`{"pad":"` + strings.Repeat("a", 300*1024) + `"}`)
		acc, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName:   aws.String(apFunc),
			InvocationType: types.InvocationTypeEvent,
			Payload:        medium,
		})
		if err != nil {
			return fmt.Errorf("300 KB event invoke: %v", err)
		}
		if acc.StatusCode != 202 {
			return fmt.Errorf("300 KB event invoke must answer 202, got %d", acc.StatusCode)
		}

		// 1 MB + 1 KB exceeds the asynchronous limit and must be rejected.
		oversized := []byte(`{"pad":"` + strings.Repeat("a", 1024*1024+1024) + `"}`)
		_, err = tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName:   aws.String(apFunc),
			InvocationType: types.InvocationTypeEvent,
			Payload:        oversized,
		})
		if err == nil {
			return fmt.Errorf("1 MB+1 KB event invoke must be rejected")
		}
		if err := expectAWSErrorCode(err, "RequestTooLargeException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_HandlerContext", func() error {
		ctxFunc, cleanupCtxFn, err := tc.setupFunction("CtxFn",
			"exports.handler = async (event, context) => { return { rid: context.awsRequestId, fn: context.functionName, ver: context.functionVersion, arn: context.invokedFunctionArn, mem: context.memoryLimitInMB, lg: context.logGroupName, ls: context.logStreamName, rem: context.getRemainingTimeInMillis() }; }",
			func(input *lambda.CreateFunctionInput) {
				input.Timeout = aws.Int32(10)
				input.MemorySize = aws.Int32(256)
			})
		if err != nil {
			return err
		}
		defer cleanupCtxFn()

		_, out, err := tc.invokeAndDecode(ctxFunc, nil)
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		rid, _ := out["rid"].(string)
		if len(rid) != 36 || strings.Count(rid, "-") != 4 {
			return fmt.Errorf("awsRequestId must be a UUID, got %q", rid)
		}
		if out["fn"] != ctxFunc {
			return fmt.Errorf("functionName mismatch: got %v", out["fn"])
		}
		if out["ver"] != "$LATEST" {
			return fmt.Errorf("functionVersion mismatch: got %v", out["ver"])
		}
		arn, _ := out["arn"].(string)
		if !strings.HasSuffix(arn, "function:"+ctxFunc) {
			return fmt.Errorf("invokedFunctionArn must end with the function name, got %q", arn)
		}
		if out["mem"] != float64(256) {
			return fmt.Errorf("memoryLimitInMB mismatch: got %v", out["mem"])
		}
		if out["lg"] != "/aws/lambda/"+ctxFunc {
			return fmt.Errorf("logGroupName mismatch: got %v", out["lg"])
		}
		ls, _ := out["ls"].(string)
		if !strings.Contains(ls, "[$LATEST]") {
			return fmt.Errorf("logStreamName must carry the version bracket, got %q", ls)
		}
		rem, ok := out["rem"].(float64)
		if !ok || rem <= 0 || rem > 10000 {
			return fmt.Errorf("getRemainingTimeInMillis must be within the 10s timeout, got %v", out["rem"])
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_HandlerContext_Qualifier", func() error {
		qualFunc := tc.unique("QualCtxFn")
		qualRole, cleanupQualRole, err := tc.createRole(tc.unique("QualCtxRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupQualRole()
		qualARN, cleanupQualFn, err := tc.createFunction(qualFunc, qualRole,
			"exports.handler = async (event, context) => { return { ver: context.functionVersion, arn: context.invokedFunctionArn }; }")
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer cleanupQualFn()

		pv, err := tc.client.PublishVersion(tc.ctx, &lambda.PublishVersionInput{FunctionName: aws.String(qualFunc)})
		if err != nil {
			return fmt.Errorf("publish version: %v", err)
		}
		_, err = tc.client.CreateAlias(tc.ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(qualFunc),
			Name:            aws.String("ctxqa"),
			FunctionVersion: pv.Version,
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		defer tc.client.DeleteAlias(tc.ctx, &lambda.DeleteAliasInput{FunctionName: aws.String(qualFunc), Name: aws.String("ctxqa")})

		_, out, err := tc.invokeAndDecode(qualFunc, &lambda.InvokeInput{
			FunctionName: aws.String(qualARN),
			Qualifier:    aws.String("ctxqa"),
		})
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		if out["ver"] != aws.ToString(pv.Version) {
			return fmt.Errorf("functionVersion must be the alias target, got %v want %v", out["ver"], aws.ToString(pv.Version))
		}
		arn, _ := out["arn"].(string)
		if !strings.HasSuffix(arn, ":"+qualFunc+":ctxqa") {
			return fmt.Errorf("invokedFunctionArn must carry the alias qualifier, got %q", arn)
		}

		// A qualifier embedded in the FunctionName ARN reaches the context the
		// same way an explicit Qualifier parameter does.
		_, embOut, err := tc.invokeAndDecode(qualFunc, &lambda.InvokeInput{
			FunctionName: aws.String(qualARN + ":" + aws.ToString(pv.Version)),
		})
		if err != nil {
			return fmt.Errorf("invoke with embedded qualifier: %v", err)
		}
		if embOut["ver"] != aws.ToString(pv.Version) {
			return fmt.Errorf("embedded qualifier must execute the published version, got %v want %v", embOut["ver"], aws.ToString(pv.Version))
		}
		embArn, _ := embOut["arn"].(string)
		if !strings.HasSuffix(embArn, ":"+qualFunc+":"+aws.ToString(pv.Version)) {
			return fmt.Errorf("invokedFunctionArn must carry the embedded version qualifier, got %q", embArn)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_ClientContext", func() error {
		ccFunc, cleanupCcFn, err := tc.setupFunction("CcFn",
			"exports.handler = async (event, context) => { return { k: context.clientContext && context.clientContext.custom && context.clientContext.custom.k }; }")
		if err != nil {
			return err
		}
		defer cleanupCcFn()

		_, out, err := tc.invokeAndDecode(ccFunc, &lambda.InvokeInput{
			ClientContext: aws.String(base64.StdEncoding.EncodeToString([]byte(`{"custom":{"k":"v"}}`))),
		})
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		if out["k"] != "v" {
			return fmt.Errorf("clientContext.custom.k must round-trip, got %v", out["k"])
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_ContextRequestIdInLogs", func() error {
		ridFunc, cleanupRidFn, err := tc.setupFunction("RidFn",
			"exports.handler = async (event, context) => { return { rid: context.awsRequestId }; }",
			func(input *lambda.CreateFunctionInput) {
				input.MemorySize = aws.Int32(256)
			})
		if err != nil {
			return err
		}
		defer cleanupRidFn()

		resp, out, err := tc.invokeAndDecode(ridFunc, &lambda.InvokeInput{
			LogType: types.LogTypeTail,
		})
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		rid, _ := out["rid"].(string)
		if rid == "" {
			return fmt.Errorf("handler must return its awsRequestId")
		}
		logBytes, err := base64.StdEncoding.DecodeString(aws.ToString(resp.LogResult))
		if err != nil {
			return fmt.Errorf("decode LogResult: %v", err)
		}
		logTail := string(logBytes)
		for _, needle := range []string{
			"START RequestId: " + rid + " Version: $LATEST",
			"END RequestId: " + rid,
			"REPORT RequestId: " + rid,
			"Memory Size: 256 MB",
		} {
			if !strings.Contains(logTail, needle) {
				return fmt.Errorf("tailed logs must contain %q, got:\n%s", needle, logTail)
			}
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_PythonHandlerContext", func() error {
		pyFunc := tc.unique("PyCtxFn")
		pyRole, cleanupPyRole, err := tc.createRole(tc.unique("PyCtxRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupPyRole()
		pyZip, err := zipLambdaPythonCode("def handler(event, context):\n    return {'rid': context.aws_request_id, 'fn': context.function_name, 'arn': context.invoked_function_arn, 'mem': context.memory_limit_in_mb, 'lg': context.log_group_name, 'rem': context.get_remaining_time_in_millis(), 's': 'x', 'b': True, 'n': None}\n")
		if err != nil {
			return fmt.Errorf("zip python code: %v", err)
		}
		_, err = tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(pyFunc),
			Runtime:      types.RuntimePython313,
			Role:         aws.String(pyRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: pyZip},
			Timeout:      aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteFunctionAndLogs(pyFunc)

		_, out, err := tc.invokeAndDecode(pyFunc, nil)
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		rid, _ := out["rid"].(string)
		if len(rid) != 36 || strings.Count(rid, "-") != 4 {
			return fmt.Errorf("aws_request_id must be a UUID, got %q", rid)
		}
		if out["fn"] != pyFunc {
			return fmt.Errorf("function_name mismatch: got %v", out["fn"])
		}
		arn, _ := out["arn"].(string)
		if !strings.HasSuffix(arn, "function:"+pyFunc) {
			return fmt.Errorf("invoked_function_arn must end with the function name, got %q", arn)
		}
		if out["mem"] != float64(128) {
			return fmt.Errorf("memory_limit_in_mb must be the default 128, got %v", out["mem"])
		}
		if out["lg"] != "/aws/lambda/"+pyFunc {
			return fmt.Errorf("log_group_name mismatch: got %v", out["lg"])
		}
		rem, ok := out["rem"].(float64)
		if !ok || rem <= 0 || rem > 10000 {
			return fmt.Errorf("get_remaining_time_in_millis must be within the 10s timeout, got %v", out["rem"])
		}
		// The mixed-type members must serialise strictly as JSON: strings
		// quoted, booleans lower-case, None as null. The whole-payload
		// Unmarshal above already proves strict JSON (a bare True/None/x
		// would not parse); these assertions pin the decoded types.
		if out["s"] != "x" {
			return fmt.Errorf("python string return must serialise as a JSON string, got %v", out["s"])
		}
		if out["b"] != true {
			return fmt.Errorf("python True must serialise as JSON true, got %v", out["b"])
		}
		if out["n"] != nil {
			return fmt.Errorf("python None must serialise as JSON null, got %v", out["n"])
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_Timeout_Envelope", func() error {
		toFunc, cleanupToFn, err := tc.setupFunction("ToFn",
			"exports.handler = async () => { await new Promise(r => setTimeout(r, 4000)); return { ok: 1 }; }",
			func(input *lambda.CreateFunctionInput) {
				input.Timeout = aws.Int32(1)
			})
		if err != nil {
			return err
		}
		defer cleanupToFn()

		resp, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{FunctionName: aws.String(toFunc)})
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("a timeout still answers 200, got %d", resp.StatusCode)
		}
		if aws.ToString(resp.FunctionError) != "Unhandled" {
			return fmt.Errorf("FunctionError must be Unhandled, got %q", aws.ToString(resp.FunctionError))
		}
		if !strings.Contains(string(resp.Payload), "Task timed out after 1") {
			return fmt.Errorf("payload must carry the timeout message, got %s", string(resp.Payload))
		}

		// A callback response is held until the event loop drains; an
		// interval behind the callback means it never drains, so the
		// timeout error replaces the recorded response.
		cbToFunc, cleanupCbToFn, err := tc.setupFunction("CbToFn",
			"exports.handler = (event, context, callback) => { setTimeout(() => callback(null, { ok: 1 }), 20); setInterval(() => {}, 500); }",
			func(input *lambda.CreateFunctionInput) {
				input.Timeout = aws.Int32(1)
			})
		if err != nil {
			return err
		}
		defer cleanupCbToFn()

		cbResp, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{FunctionName: aws.String(cbToFunc)})
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		if aws.ToString(cbResp.FunctionError) != "Unhandled" {
			return fmt.Errorf("a callback handler that never drains must time out as Unhandled, got %q", aws.ToString(cbResp.FunctionError))
		}
		if !strings.Contains(string(cbResp.Payload), "Task timed out after 1") {
			return fmt.Errorf("payload must carry the timeout message, got %s", string(cbResp.Payload))
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_AsyncHandlerLingeringTimer", func() error {
		// An async handler settles its response when the promise resolves;
		// a timer left on the event loop must not hold the invocation open
		// until the timeout overrides the delivered result.
		ltFunc, cleanupLtFn, err := tc.setupFunction("LtFn",
			"exports.handler = async () => { setInterval(() => {}, 1000); return { ok: 1 }; }",
			func(input *lambda.CreateFunctionInput) {
				input.Timeout = aws.Int32(3)
			})
		if err != nil {
			return err
		}
		defer cleanupLtFn()

		invokeStart := time.Now()
		resp, out, err := tc.invokeAndDecode(ltFunc, nil)
		elapsed := time.Since(invokeStart)
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		if elapsed >= 3*time.Second {
			return fmt.Errorf("the response must arrive at the promise resolution, not the timeout, took %v", elapsed)
		}
		if out["ok"] != float64(1) {
			return fmt.Errorf("the delivered result must survive the lingering timer, got %v (%s)", out["ok"], string(resp.Payload))
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_CallbackHandler", func() error {
		cbFunc, cleanupCbFn, err := tc.setupFunction("CbFn",
			"exports.handler = (event, context, callback) => { setTimeout(() => callback(null, { ok: 1 }), 20); }")
		if err != nil {
			return err
		}
		defer cleanupCbFn()

		resp, out, err := tc.invokeAndDecode(cbFunc, nil)
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		if out["ok"] != float64(1) {
			return fmt.Errorf("callback result must be the payload, got %v (%s)", out["ok"], string(resp.Payload))
		}

		// A scalar callback result must serialise as JSON: a string
		// arrives quoted, not as a bare token.
		cbStrFunc, cleanupCbStrFn, err := tc.setupFunction("CbStrFn",
			"exports.handler = (event, context, callback) => { callback(null, 'ok'); }")
		if err != nil {
			return err
		}
		defer cleanupCbStrFn()
		strResp, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{FunctionName: aws.String(cbStrFunc)})
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		if string(strResp.Payload) != `"ok"` {
			return fmt.Errorf("a string callback result must serialise as a quoted JSON string, got %s", string(strResp.Payload))
		}

		// The callback response is not sent until all event loop tasks
		// have finished: work scheduled behind the callback delays the
		// response, so the answer arrives after the trailing timer runs.
		drainFunc, cleanupDrainFn, err := tc.setupFunction("CbDrainFn",
			"exports.handler = (event, context, callback) => { setTimeout(() => callback(null, { ok: 1 }), 20); setTimeout(() => console.log('late work'), 400); }")
		if err != nil {
			return err
		}
		defer cleanupDrainFn()
		invokeStart := time.Now()
		_, drainOut, err := tc.invokeAndDecode(drainFunc, nil)
		elapsed := time.Since(invokeStart)
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		if elapsed < 380*time.Millisecond {
			return fmt.Errorf("the callback response must wait for the event loop to drain, took %v", elapsed)
		}
		if drainOut["ok"] != float64(1) {
			return fmt.Errorf("the drained callback result must still be the payload, got %v", drainOut["ok"])
		}
		return nil
	}))

	// A custom runtime exchanges the event and its answer over the Runtime
	// API. The test bootstrap follows the AWS runtime pattern — it answers,
	// then loops back to /invocation/next for the next event — so the test
	// also pins that a looping runtime does not hold the synchronous invoke
	// open until the function timeout: the response must come back as soon
	// as the answer POST lands.
	results = append(results, tc.r.RunTest("lambda", "Invoke_ProvidedRuntime_Context", func() error {
		provFunc := tc.unique("ProvFn")
		provRole, cleanupProvRole, err := tc.createRole(tc.unique("ProvRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupProvRole()
		bootstrap := `#!/bin/sh
set -eu
API="http://$AWS_LAMBDA_RUNTIME_API"
while :; do
	hdr="$(curl -sS -f -D - -o /tmp/event.json "$API/2018-06-01/runtime/invocation/next")"
	rid="$(printf '%s\n' "$hdr" | tr -d '\r' | sed -n 's/^Lambda-Runtime-Aws-Request-Id: //p')"
	resp="$(printf '{"rid":"%s","event":' "$rid"; cat /tmp/event.json; printf '}')"
	curl -sS -f -X POST --data-binary "$resp" "$API/2018-06-01/runtime/invocation/$rid/response"
done
`
		provZip, err := zipLambdaBootstrap(bootstrap)
		if err != nil {
			return fmt.Errorf("zip bootstrap: %v", err)
		}
		_, err = tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(provFunc),
			Runtime:      types.RuntimeProvidedal2023,
			Role:         aws.String(provRole),
			Handler:      aws.String("bootstrap"),
			Code:         &types.FunctionCode{ZipFile: provZip},
			Timeout:      aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteFunctionAndLogs(provFunc)

		invokeStart := time.Now()
		resp, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(provFunc),
			Payload:      []byte(`{"hello":"provided"}`),
		})
		elapsed := time.Since(invokeStart)
		if err != nil {
			return fmt.Errorf("invoke: %v", err)
		}
		if elapsed >= 8*time.Second {
			return fmt.Errorf("a looping runtime must not hold the invoke until the timeout, took %v", elapsed)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if aws.ToString(resp.FunctionError) != "" {
			return fmt.Errorf("unexpected function error %q (payload %s)", aws.ToString(resp.FunctionError), string(resp.Payload))
		}
		var out struct {
			RID   string                 `json:"rid"`
			Event map[string]interface{} `json:"event"`
		}
		if err := json.Unmarshal(resp.Payload, &out); err != nil {
			return fmt.Errorf("parse payload: %v (%s)", err, string(resp.Payload))
		}
		if len(out.RID) != 36 || strings.Count(out.RID, "-") != 4 {
			return fmt.Errorf("request id from the Runtime API header must be a UUID, got %q", out.RID)
		}
		if out.Event["hello"] != "provided" {
			return fmt.Errorf("event from /invocation/next must round-trip, got %v (%s)", out.Event["hello"], string(resp.Payload))
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetFunction_ContainsCodeConfig", func() error {
		gfcDesc := "Test description for verification"
		gfcFunc, cleanupGfcFn, err := tc.setupFunction("GfcFunc", "exports.handler = async () => { return 1; };",
			func(input *lambda.CreateFunctionInput) {
				input.Description = aws.String(gfcDesc)
				input.Timeout = aws.Int32(15)
				input.MemorySize = aws.Int32(256)
			})
		if err != nil {
			return err
		}
		defer cleanupGfcFn()

		resp, err := tc.client.GetFunction(tc.ctx, &lambda.GetFunctionInput{
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

	results = append(results, tc.r.RunTest("lambda", "UpdateFunctionConfiguration_VerifyUpdate", func() error {
		ucFunc, cleanupUcFn, err := tc.setupFunction("UcFunc", "exports.handler = async () => { return 1; };",
			func(input *lambda.CreateFunctionInput) {
				input.Description = aws.String("original")
			})
		if err != nil {
			return err
		}
		defer cleanupUcFn()

		newDesc := "updated description"
		newTimeout := int32(30)
		newMemory := int32(512)
		_, err = tc.client.UpdateFunctionConfiguration(tc.ctx, &lambda.UpdateFunctionConfigurationInput{
			FunctionName: aws.String(ucFunc),
			Description:  aws.String(newDesc),
			Timeout:      aws.Int32(newTimeout),
			MemorySize:   aws.Int32(newMemory),
		})
		if err != nil {
			return fmt.Errorf("update config: %v", err)
		}

		resp, err := tc.client.GetFunctionConfiguration(tc.ctx, &lambda.GetFunctionConfigurationInput{
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

	results = append(results, tc.r.RunTest("lambda", "UpdateFunctionConfiguration_EphemeralStorageAndSnapStart", func() error {
		esFunc, cleanupEsFn, err := tc.setupFunction("EsFunc", "exports.handler = async () => { return 1; };",
			func(input *lambda.CreateFunctionInput) {
				input.Runtime = types.RuntimeJava21
				input.Handler = aws.String("org.example.App::handleRequest")
			})
		if err != nil {
			return err
		}
		defer cleanupEsFn()

		if _, err := tc.client.UpdateFunctionConfiguration(tc.ctx, &lambda.UpdateFunctionConfigurationInput{
			FunctionName:     aws.String(esFunc),
			EphemeralStorage: &types.EphemeralStorage{Size: aws.Int32(1024)},
			SnapStart:        &types.SnapStart{ApplyOn: types.SnapStartApplyOnPublishedVersions},
		}); err != nil {
			return fmt.Errorf("update config: %v", err)
		}

		resp, err := tc.client.GetFunctionConfiguration(tc.ctx, &lambda.GetFunctionConfigurationInput{
			FunctionName: aws.String(esFunc),
		})
		if err != nil {
			return err
		}
		if resp.EphemeralStorage == nil || resp.EphemeralStorage.Size == nil || *resp.EphemeralStorage.Size != 1024 {
			return fmt.Errorf("EphemeralStorage not applied, got %v", resp.EphemeralStorage)
		}
		if resp.SnapStart == nil || resp.SnapStart.ApplyOn != types.SnapStartApplyOnPublishedVersions {
			return fmt.Errorf("SnapStart not applied, got %v", resp.SnapStart)
		}

		// A present-but-negative member is rejected rather than silently
		// ignored.
		_, err = tc.client.UpdateFunctionConfiguration(tc.ctx, &lambda.UpdateFunctionConfigurationInput{
			FunctionName: aws.String(esFunc),
			Timeout:      aws.Int32(-1),
		})
		if err := expectAWSErrorCode(err, "InvalidParameterValueException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ListFunctions_Pagination", func() error {
		pgTs := tc.ts
		var pgFuncs []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagFunc-%s-%d", pgTs, i)
			_, cleanupPgFn, err := tc.createFunction(name, roleARN, lambdaFunctionCode)
			if err != nil {
				return fmt.Errorf("create function %s: %v", name, err)
			}
			defer cleanupPgFn()
			pgFuncs = append(pgFuncs, name)
		}

		all, err := paginate[types.FunctionConfiguration](func(next *string) ([]types.FunctionConfiguration, *string, error) {
			resp, err := tc.client.ListFunctions(tc.ctx, &lambda.ListFunctionsInput{
				Marker:   next,
				MaxItems: aws.Int32(2),
			})
			if err != nil {
				return nil, nil, err
			}
			return resp.Functions, resp.NextMarker, nil
		})
		if err != nil {
			return fmt.Errorf("list functions page: %v", err)
		}
		var allFuncs []string
		for _, f := range all {
			if name := aws.ToString(f.FunctionName); strings.HasPrefix(name, "PagFunc-"+pgTs) {
				allFuncs = append(allFuncs, name)
			}
		}
		if len(allFuncs) != 5 {
			return fmt.Errorf("expected 5 paginated functions, got %d", len(allFuncs))
		}
		return nil
	}))

	return results
}
