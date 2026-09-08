// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// UpdateFunctionCode updates the code of the specified Lambda function.
func (s *LambdaService) UpdateFunctionCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	functionName = extractFunctionName(functionName)
	if err := validateFunctionName(functionName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	codeMap := request.GetMapParam(req.Parameters, "Code")
	if codeMap == nil {
		codeMap = map[string]interface{}{}
	}
	// UpdateFunctionCode carries the code members at the top level of the
	// request (only CreateFunction nests them under Code), so each member
	// has a flat fallback like ZipFile and ImageUri; a member already read
	// from the Code map wins over the flat one.
	for _, member := range []string{"ZipFile", "ImageUri", "S3Bucket", "S3Key", "S3ObjectVersion"} {
		if _, ok := codeMap[member]; !ok {
			if v, ok := req.Parameters[member].(string); ok {
				codeMap[member] = v
			}
		}
	}

	codeMeta, err := s.prepareFunctionCodeUpdateCore(ctx, reqCtx.GetRegion(), functionName, codeMap)
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

	imageUri, _ := codeMap["ImageUri"].(string)
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

	var newFileSystemConfigs []lambdastore.FileSystemConfig
	if fscs, ok := req.Parameters["FileSystemConfigs"].([]interface{}); ok {
		newFileSystemConfigs = parseFileSystemConfigs(fscs)
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

	update := &UpdateFunctionConfigurationInput{
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
		IAMValidator:         reqCtx.GetIAMValidator(),
		RevisionId:           request.GetStringParam(req.Parameters, "RevisionId"),
	}
	// Presence detection for the members whose zero value carries request
	// meaning (an empty Description clears, a zero Timeout or MemorySize is
	// range-rejected); the Core owns the semantics for both planes.
	if _, ok := req.Parameters["Description"]; ok {
		update.HasDescription = true
	}
	if _, ok := req.Parameters["Timeout"]; ok {
		update.HasTimeout = true
	}
	if _, ok := req.Parameters["MemorySize"]; ok {
		update.HasMemorySize = true
	}
	if request.GetMapParam(req.Parameters, "DeadLetterConfig") != nil {
		update.HasDeadLetterConfig = true
	}

	function, err := s.updateFunctionConfigurationCore(ctx, store, update)
	if err != nil {
		return nil, err
	}

	return s.toFunctionConfiguration(function), nil
}
