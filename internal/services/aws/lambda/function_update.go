// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"os"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// UpdateFunctionCode updates the code of the specified Lambda function.
func (s *LambdaService) UpdateFunctionCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	functionName = extractFunctionName(functionName)
	if err := validateFunctionName(functionName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	codeMap := request.GetMapParam(req.Parameters, "Code")
	zipFileStr := ""
	imageUri := ""
	s3Bucket := ""
	s3Key := ""
	s3Version := ""

	if codeMap != nil {
		if z, ok := codeMap["ZipFile"].(string); ok {
			zipFileStr = z
		}
		if i, ok := codeMap["ImageUri"].(string); ok {
			imageUri = i
		}
		if b, ok := codeMap["S3Bucket"].(string); ok {
			s3Bucket = b
		}
		if k, ok := codeMap["S3Key"].(string); ok {
			s3Key = k
		}
		if v, ok := codeMap["S3ObjectVersion"].(string); ok {
			s3Version = v
		}
	}
	if zipFileStr == "" {
		if z, ok := req.Parameters["ZipFile"].(string); ok {
			zipFileStr = z
		}
	}
	if imageUri == "" {
		if i, ok := req.Parameters["ImageUri"].(string); ok {
			imageUri = i
		}
	}
	// UpdateFunctionCode carries the code members at the top level of the
	// request (only CreateFunction nests them under Code), so the S3
	// placement has a flat fallback like ZipFile and ImageUri.
	if s3Bucket == "" {
		if b, ok := req.Parameters["S3Bucket"].(string); ok {
			s3Bucket = b
		}
	}
	if s3Key == "" {
		if k, ok := req.Parameters["S3Key"].(string); ok {
			s3Key = k
		}
	}
	if s3Version == "" {
		if v, ok := req.Parameters["S3ObjectVersion"].(string); ok {
			s3Version = v
		}
	}

	codeMeta, err := s.prepareFunctionCodeUpdateCore(ctx, reqCtx.GetRegion(), functionName, zipFileStr, imageUri, s3Bucket, s3Key, s3Version)
	if err != nil {
		return nil, err
	}

	// Parse architectures into string slice.
	var architectures []string
	if archs, ok := req.Parameters["Architectures"].([]interface{}); ok {
		architectures = make([]string, 0, len(archs))
		for _, a := range archs {
			if as, ok := a.(string); ok {
				architectures = append(architectures, as)
			}
		}
	}

	// DryRun validates the request without modifying the code.
	if request.GetBoolParam(req.Parameters, "DryRun") {
		current, err := s.getFunctionForDryRunCore(store, functionName)
		if err != nil {
			return nil, err
		}
		return s.toFunctionConfiguration(current), nil
	}

	function, published, err := s.updateFunctionCodeCore(store, &UpdateFunctionCodeInput{
		FunctionName:  functionName,
		CodeLocation:  codeMeta.CodeLocation,
		CodeSize:      codeMeta.CodeSize,
		CodeSha256:    codeMeta.CodeSha256,
		ImageUri:      imageUri,
		Architectures: architectures,
		Publish:       request.GetBoolParam(req.Parameters, "Publish"),
		Region:        reqCtx.GetRegion(),
		RevisionId:    request.GetStringParam(req.Parameters, "RevisionId"),
	})
	if err != nil {
		return nil, err
	}

	// When Publish was requested, the response describes the published
	// version rather than $LATEST.
	if published != nil {
		return s.toVersionConfiguration(published), nil
	}
	return s.toFunctionConfiguration(function), nil
}

// UpdateFunctionConfiguration updates the configuration of the specified Lambda function.
func (s *LambdaService) UpdateFunctionConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	functionName = extractFunctionName(functionName)
	if err := validateFunctionName(functionName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	runtime := request.GetStringParam(req.Parameters, "Runtime")
	role := request.GetStringParam(req.Parameters, "Role")
	if role != "" {
		validator := reqCtx.GetIAMValidator()
		if os.Getenv("TEST_MODE") != "true" {
			if err := validator.ValidateRoleForService(ctx, role, iam.ServicePrincipalLambda); err != nil {
				return nil, err
			}
		}
	}

	// Parse and resolve VpcConfig before the Core call so that
	// EC2 subnet validation (I/O) happens outside the store lock.
	var newVpcConfig *lambdastore.VpcConfig
	if request.GetMapParam(req.Parameters, "VpcConfig") != nil {
		newVpcConfig = parseVpcConfig(req.Parameters)
		if newVpcConfig != nil && len(newVpcConfig.SubnetIds) > 0 {
			if err := s.resolveVpcConfig(ctx, reqCtx.GetRegion(), newVpcConfig); err != nil {
				return nil, err
			}
		}
	}

	var newEnvironment *lambdastore.Environment
	if request.GetMapParam(req.Parameters, "Environment") != nil {
		newEnvironment = parseEnvironment(req.Parameters)
	}

	var newDeadLetterConfig *lambdastore.DeadLetterConfig
	if request.GetMapParam(req.Parameters, "DeadLetterConfig") != nil {
		newDeadLetterConfig, err = parseDeadLetterConfig(req.Parameters)
		if err != nil {
			return nil, err
		}
	}

	var newTracingConfig *lambdastore.TracingConfig
	if request.GetMapParam(req.Parameters, "TracingConfig") != nil {
		newTracingConfig, err = parseTracingConfig(req.Parameters)
		if err != nil {
			return nil, err
		}
	}

	var newLoggingConfig *lambdastore.LoggingConfig
	if logMap := request.GetMapParam(req.Parameters, "LoggingConfig"); logMap != nil {
		newLoggingConfig = parseLoggingConfig(logMap)
	}

	var newImageConfig *lambdastore.ImageConfig
	if imgMap := request.GetMapParam(req.Parameters, "ImageConfig"); imgMap != nil {
		newImageConfig = parseImageConfig(imgMap)
	}

	var newEphemeralStorage *lambdastore.EphemeralStorage
	if esMap := request.GetMapParam(req.Parameters, "EphemeralStorage"); esMap != nil {
		newEphemeralStorage = &lambdastore.EphemeralStorage{
			Size: int32(request.GetIntParam(esMap, "Size")),
		}
	}

	var newSnapStart *lambdastore.SnapStart
	if ssMap := request.GetMapParam(req.Parameters, "SnapStart"); ssMap != nil {
		newSnapStart = &lambdastore.SnapStart{
			ApplyOn: request.GetStringParam(ssMap, "ApplyOn"),
		}
	}

	// A present member is validated as-is: negative or zero values are
	// rejected instead of being silently ignored as "not provided".
	if _, ok := req.Parameters["Timeout"]; ok {
		if err := validateTimeout(int32(request.GetIntParam(req.Parameters, "Timeout"))); err != nil {
			return nil, err
		}
	}
	if _, ok := req.Parameters["MemorySize"]; ok {
		if err := validateMemorySize(int32(request.GetIntParam(req.Parameters, "MemorySize"))); err != nil {
			return nil, err
		}
	}

	var newFileSystemConfigs []lambdastore.FileSystemConfig
	if fscs, ok := req.Parameters["FileSystemConfigs"].([]interface{}); ok {
		for _, fsc := range fscs {
			if m, ok := fsc.(map[string]interface{}); ok {
				newFileSystemConfigs = append(newFileSystemConfigs, lambdastore.FileSystemConfig{
					Arn:            request.GetStringParam(m, "Arn"),
					LocalMountPath: request.GetStringParam(m, "LocalMountPath"),
				})
			}
		}
	}

	var newLayers []lambdastore.LayerReference
	if layers, ok := req.Parameters["Layers"].([]interface{}); ok {
		newLayers = make([]lambdastore.LayerReference, 0, len(layers))
		for _, l := range layers {
			if ls, ok := l.(string); ok {
				newLayers = append(newLayers, lambdastore.LayerReference{Arn: ls})
			}
		}
	}

	function, err := s.updateFunctionConfigurationCore(ctx, store, &UpdateFunctionConfigurationInput{
		FunctionName:         functionName,
		Runtime:              runtime,
		Role:                 role,
		Handler:              request.GetStringParam(req.Parameters, "Handler"),
		Description:          request.GetStringParam(req.Parameters, "Description"),
		Timeout:              int32(request.GetIntParam(req.Parameters, "Timeout")),
		MemorySize:           int32(request.GetIntParam(req.Parameters, "MemorySize")),
		KMSKeyArn:            request.GetStringParam(req.Parameters, "KMSKeyArn"),
		CodeSigningConfigArn: request.GetStringParam(req.Parameters, "CodeSigningConfigArn"),
		VpcConfig:            newVpcConfig,
		Environment:          newEnvironment,
		DeadLetterConfig:     newDeadLetterConfig,
		TracingConfig:        newTracingConfig,
		LoggingConfig:        newLoggingConfig,
		ImageConfig:          newImageConfig,
		EphemeralStorage:     newEphemeralStorage,
		SnapStart:            newSnapStart,
		FileSystemConfigs:    newFileSystemConfigs,
		Layers:               newLayers,
	})
	if err != nil {
		return nil, err
	}

	return s.toFunctionConfiguration(function), nil
}
