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
		lfFunc := tc.unique("LfFunc")
		lfRole, cleanupLfRole, err := tc.createRole(tc.unique("LfRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupLfRole()
		_, cleanupLfFn, err := tc.createFunction(lfFunc, lfRole, "exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("create: %v", err)
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
		defer tc.client.DeleteFunction(tc.ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(fnName)})
		defer deleteLambdaLogGroup(tc.cwl, tc.ctx, fnName)
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
		} else if err := AssertErrorContains(err, "InvalidParameterValueException"); err != nil {
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
			if err := AssertErrorContains(err, "InvalidParameterValueException"); err != nil {
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
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
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
		if err := AssertErrorContains(err, "ResourceConflictException"); err != nil {
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
		defer tc.client.DeleteFunction(tc.ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(imgFunc)})
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
		snapFunc := tc.unique("SnapFunc")
		snapRole, cleanupSnapRole, err := tc.createRole(tc.unique("SnapRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupSnapRole()
		_, cleanupSnapFn, err := tc.createFunction(snapFunc, snapRole, "exports.handler = async () => 'v1-output';")
		if err != nil {
			return fmt.Errorf("create function: %v", err)
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
		feFunc := tc.unique("FeFunc")
		feRole, cleanupFeRole, err := tc.createRole(tc.unique("FeRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupFeRole()

		// A thrown error is intercepted by the runtime: Unhandled.
		_, cleanupFeFn, err := tc.createFunction(feFunc, feRole, "exports.handler = async () => { throw new Error('boom'); };")
		if err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer cleanupFeFn()

		unhandled, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(feFunc),
		})
		if err != nil {
			return err
		}
		if unhandled.FunctionError == nil || *unhandled.FunctionError != "Unhandled" {
			return fmt.Errorf("thrown error should classify as Unhandled, got %v", unhandled.FunctionError)
		}

		// A returned errorMessage envelope is a handled error: HTTP 200
		// with the error document as the payload.
		handledCode, err := zipLambdaCode("exports.handler = async () => ({ errorMessage: 'handled-boom', errorType: 'Error' });")
		if err != nil {
			return fmt.Errorf("zip handled code: %v", err)
		}
		if _, err := tc.client.UpdateFunctionCode(tc.ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(feFunc),
			ZipFile:      handledCode,
		}); err != nil {
			return fmt.Errorf("update code: %v", err)
		}
		handled, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(feFunc),
		})
		if err != nil {
			return err
		}
		if handled.FunctionError == nil || *handled.FunctionError != "Handled" {
			return fmt.Errorf("returned error document should classify as Handled, got %v", handled.FunctionError)
		}
		if !strings.Contains(string(handled.Payload), "handled-boom") {
			return fmt.Errorf("Handled payload should carry the error document, got %q", string(handled.Payload))
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_HandledErrorWithLogs", func() error {
		logFunc := tc.unique("LogFeFunc")
		logRole, cleanupLogRole, err := tc.createRole(tc.unique("LogFeRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupLogRole()

		// Handler console output reaches stdout before the runtime appends
		// the returned payload, so the error envelope must be read from the
		// final JSON document rather than the whole output.
		_, cleanupLogFn, err := tc.createFunction(logFunc, logRole, "exports.handler = async () => { console.log('about to fail'); return { errorMessage: 'logged-boom', errorType: 'Error' }; };")
		if err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer cleanupLogFn()

		result, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(logFunc),
		})
		if err != nil {
			return err
		}
		if result.FunctionError == nil || *result.FunctionError != "Handled" {
			return fmt.Errorf("error envelope after log lines should classify as Handled, got %v", result.FunctionError)
		}
		if !strings.Contains(string(result.Payload), "logged-boom") {
			return fmt.Errorf("payload should carry the error document after the log lines, got %q", string(result.Payload))
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_PayloadExcludesConsoleOutput", func() error {
		logFunc := tc.unique("PayloadCleanFn")
		logRole, cleanupLogRole, err := tc.createRole(tc.unique("PayloadCleanRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupLogRole()

		// On AWS the response payload is the handler's return value only;
		// console output belongs to the logs. A handler that logs before
		// returning must not leak its log lines into the payload.
		_, cleanupLogFn, err := tc.createFunction(logFunc, logRole, "exports.handler = async () => { console.log('log-line-should-not-leak'); return { ok: true, value: 7 }; };")
		if err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer cleanupLogFn()

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
		listFunc := tc.unique("ListAllFn")
		listRole, cleanupListRole, err := tc.createRole(tc.unique("ListAllRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupListRole()

		_, cleanupListFn, err := tc.createFunction(listFunc, listRole, "exports.handler = async () => ({});")
		if err != nil {
			return fmt.Errorf("create function: %v", err)
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
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// === ERROR CASES ===

	results = append(results, tc.r.RunTest("lambda", "GetFunction_NonExistent", func() error {
		_, err := tc.client.GetFunction(tc.ctx, &lambda.GetFunctionInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "Invoke_NonExistent", func() error {
		_, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateFunctionCode_NonExistent", func() error {
		_, err := tc.client.UpdateFunctionCode(tc.ctx, &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
			ZipFile:      []byte("code"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "DeleteFunction_NonExistent", func() error {
		_, err := tc.client.DeleteFunction(tc.ctx, &lambda.DeleteFunctionInput{
			FunctionName: aws.String("NoSuchFunction_xyz_12345"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
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
		defer tc.client.DeleteFunction(tc.ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(dupName)})
		defer deleteLambdaLogGroup(tc.cwl, tc.ctx, dupName)

		_, err = tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
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

	results = append(results, tc.r.RunTest("lambda", "CreateFunction_InvalidRuntime", func() error {
		invRtFuncName := tc.unique("InvRtFunc")
		invRtRole, cleanupInvRtRole, err := tc.createRole(tc.unique("InvRtRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupInvRtRole()
		_, err = tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
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
		defer tc.client.DeleteFunction(tc.ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(alRtFuncName)})
		if resp.Runtime != types.RuntimeJava17al2023 {
			return fmt.Errorf("Runtime mismatch, got %v", resp.Runtime)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "CreateFunction_SnapStartUnsupportedRuntime", func() error {
		ssFuncName := tc.unique("SsFunc")
		ssRole, cleanupSsRole, err := tc.createRole(tc.unique("SsRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupSsRole()
		ssCode, err := zipLambdaCode("exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("zip code: %v", err)
		}
		_, err = tc.client.CreateFunction(tc.ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(ssFuncName),
			Role:         aws.String(ssRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: ssCode},
			Runtime:      types.RuntimeNodejs22x,
			SnapStart:    &types.SnapStart{ApplyOn: types.SnapStartApplyOnPublishedVersions},
		})
		if err := AssertErrorContains(err, "InvalidParameterValueException"); err != nil {
			return err
		}
		return nil
	}))

	// === VERIFICATION TESTS ===

	results = append(results, tc.r.RunTest("lambda", "Invoke_VerifyResponsePayload", func() error {
		invFunc := tc.unique("InvFunc")
		invRole, cleanupInvRole, err := tc.createRole(tc.unique("InvRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupInvRole()
		_, cleanupInvFn, err := tc.createFunction(invFunc, invRole, "exports.handler = async (event) => { return { statusCode: 200, body: JSON.stringify({result: 'ok'}) }; };")
		if err != nil {
			return fmt.Errorf("create: %v", err)
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

	results = append(results, tc.r.RunTest("lambda", "GetFunction_ContainsCodeConfig", func() error {
		gfcFunc := tc.unique("GfcFunc")
		gfcDesc := "Test description for verification"
		gfcRole, cleanupGfcRole, err := tc.createRole(tc.unique("GfcRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupGfcRole()
		_, cleanupGfcFn, err := tc.createFunction(gfcFunc, gfcRole, "exports.handler = async () => { return 1; };",
			func(input *lambda.CreateFunctionInput) {
				input.Description = aws.String(gfcDesc)
				input.Timeout = aws.Int32(15)
				input.MemorySize = aws.Int32(256)
			})
		if err != nil {
			return fmt.Errorf("create: %v", err)
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
		ucFunc := tc.unique("UcFunc")
		ucRole, cleanupUcRole, err := tc.createRole(tc.unique("UcRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupUcRole()
		_, cleanupUcFn, err := tc.createFunction(ucFunc, ucRole, "exports.handler = async () => { return 1; };",
			func(input *lambda.CreateFunctionInput) {
				input.Description = aws.String("original")
			})
		if err != nil {
			return fmt.Errorf("create: %v", err)
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
		esFunc := tc.unique("EsFunc")
		esRole, cleanupEsRole, err := tc.createRole(tc.unique("EsRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupEsRole()
		_, cleanupEsFn, err := tc.createFunction(esFunc, esRole, "exports.handler = async () => { return 1; };",
			func(input *lambda.CreateFunctionInput) {
				input.Runtime = types.RuntimeJava21
				input.Handler = aws.String("org.example.App::handleRequest")
			})
		if err != nil {
			return fmt.Errorf("create: %v", err)
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
		if err := AssertErrorContains(err, "InvalidParameterValueException"); err != nil {
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
