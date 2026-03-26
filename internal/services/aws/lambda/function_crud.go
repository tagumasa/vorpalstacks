// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/common/iam"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// CreateFunction creates a new Lambda function.
func (s *LambdaService) CreateFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	if err := validateFunctionName(functionName); err != nil {
		return nil, err
	}

	runtime := request.GetStringParam(req.Parameters, "Runtime")
	role := request.GetStringParam(req.Parameters, "Role")
	handler := request.GetStringParam(req.Parameters, "Handler")
	packageType := request.GetStringParam(req.Parameters, "PackageType")

	if runtime == "" && packageType != "Image" {
		return nil, NewInvalidParameter("Runtime", "Runtime is required for Zip package type")
	}

	if runtime != "" && !ValidateRuntime(runtime) {
		return nil, NewInvalidParameter("Runtime", "Runtime '"+runtime+"' is not supported")
	}

	if handler == "" && packageType != "Image" {
		return nil, NewInvalidParameter("Handler", "Handler is required for Zip package type")
	}

	if handler != "" && runtime != "" {
		if err := ValidateHandler(runtime, handler); err != nil {
			return nil, err
		}
	}

	if role == "" {
		return nil, NewInvalidParameter("Role", "Role ARN is required")
	}

	validator := reqCtx.GetIAMValidator()
	if err := validator.ValidateRoleForService(ctx, role, iam.ServicePrincipalLambda); err != nil {
		return nil, err
	}

	codeMap := request.GetMapParam(req.Parameters, "Code")
	if codeMap == nil {
		return nil, NewInvalidParameter("Code", "Code is required")
	}

	function := &lambdastore.Function{
		FunctionName: functionName,
		Runtime:      lambdastore.Runtime(runtime),
		Role:         role,
		Handler:      handler,
		Description:  request.GetStringParam(req.Parameters, "Description"),
		Publish:      request.GetBoolParam(req.Parameters, "Publish"),
		KMSKeyArn:    request.GetStringParam(req.Parameters, "KMSKeyArn"),
		PackageType:  packageType,
	}

	if timeout := request.GetIntParam(req.Parameters, "Timeout"); timeout > 0 {
		if err := validateTimeout(int32(timeout)); err != nil {
			return nil, err
		}
		function.Timeout = int32(timeout)
	} else {
		function.Timeout = 3
	}
	if memorySize := request.GetIntParam(req.Parameters, "MemorySize"); memorySize > 0 {
		if err := validateMemorySize(int32(memorySize)); err != nil {
			return nil, err
		}
		function.MemorySize = int32(memorySize)
	} else {
		function.MemorySize = 128
	}

	if ephemeralMap := request.GetMapParam(req.Parameters, "EphemeralStorage"); ephemeralMap != nil {
		if size, ok := ephemeralMap["Size"]; ok {
			var sizeValue int32
			switch v := size.(type) {
			case int:
				sizeValue = int32(v)
			case float64:
				sizeValue = int32(v)
			default:
				return nil, NewInvalidParameter("EphemeralStorage.Size", "must be an integer")
			}
			if sizeValue < 512 || sizeValue > 10240 {
				return nil, NewInvalidParameter("EphemeralStorage.Size", "must be between 512 and 10240 MB")
			}
			function.EphemeralStorage = &lambdastore.EphemeralStorage{Size: sizeValue}
		}
	}

	if vpcMap := request.GetMapParam(req.Parameters, "VpcConfig"); vpcMap != nil {
		function.VpcConfig = parseVpcConfig(req.Parameters)
	}

	if envMap := request.GetMapParam(req.Parameters, "Environment"); envMap != nil {
		function.Environment = parseEnvironment(req.Parameters)
	}

	if dlMap := request.GetMapParam(req.Parameters, "DeadLetterConfig"); dlMap != nil {
		dl, err := parseDeadLetterConfig(req.Parameters)
		if err != nil {
			return nil, err
		}
		function.DeadLetterConfig = dl
	}

	if traceMap := request.GetMapParam(req.Parameters, "TracingConfig"); traceMap != nil {
		trace, err := parseTracingConfig(req.Parameters)
		if err != nil {
			return nil, err
		}
		function.TracingConfig = trace
	}

	if snapMap := request.GetMapParam(req.Parameters, "SnapStart"); snapMap != nil {
		function.SnapStart = &lambdastore.SnapStart{}
		if applyOn, ok := snapMap["ApplyOn"].(string); ok {
			function.SnapStart.ApplyOn = applyOn
		}
	}

	if imageUri, ok := codeMap["ImageUri"].(string); ok && imageUri != "" {
		function.ImageUri = imageUri
		function.PackageType = "Image"
	}

	if s3Bucket, ok := codeMap["S3Bucket"].(string); ok && s3Bucket != "" {
		s3Key, _ := codeMap["S3Key"].(string)
		if s3Key == "" {
			return nil, NewInvalidParameter("Code.S3Key", "S3Key is required when S3Bucket is specified")
		}
		zipFile, err := s.fetchCodeFromS3(ctx, s3Bucket, s3Key, reqCtx.GetRegion())
		if err != nil {
			return nil, NewInvalidParameter("Code", err.Error())
		}
		codePath, codeSize, err := s.storeCode(functionName, "$LATEST", zipFile, reqCtx.GetRegion())
		if err != nil {
			return nil, err
		}
		function.CodeLocation = codePath
		function.CodeSize = codeSize
		function.CodeSha256 = lambdastore.GenerateCodeHash(zipFile)
	}

	if zipFileStr, ok := codeMap["ZipFile"].(string); ok && zipFileStr != "" {
		zipFile, err := base64.StdEncoding.DecodeString(zipFileStr)
		if err != nil {
			return nil, NewInvalidParameter("Code.ZipFile", "Invalid base64 encoding: "+err.Error())
		}
		codePath, codeSize, err := s.storeCode(functionName, "$LATEST", zipFile, reqCtx.GetRegion())
		if err != nil {
			return nil, err
		}
		function.CodeLocation = codePath
		function.CodeSize = codeSize
		function.CodeSha256 = lambdastore.GenerateCodeHash(zipFile)
	}

	if tagsMap, ok := req.Parameters["Tags"].(map[string]interface{}); ok {
		function.Tags = tagutil.MapInterfaceToTags(tagsMap)
	}

	if layers, ok := req.Parameters["Layers"].([]interface{}); ok {
		for _, l := range layers {
			if ls, ok := l.(string); ok {
				if !isValidLayerARN(ls) {
					return nil, NewInvalidParameter("Layers", "Invalid layer ARN format: "+ls)
				}
				function.Layers = append(function.Layers, lambdastore.LayerReference{Arn: ls})
			}
		}
	}

	if archs, ok := req.Parameters["Architectures"].([]interface{}); ok {
		function.Architectures = make([]string, 0, len(archs))
		for _, a := range archs {
			if as, ok := a.(string); ok {
				if as != "x86_64" && as != "arm64" {
					return nil, NewInvalidParameter("Architectures", "must be x86_64 or arm64")
				}
				function.Architectures = append(function.Architectures, as)
			}
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := store.Functions.Create(function)
	if err != nil {
		if errors.Is(err, lambdastore.ErrFunctionAlreadyExists) {
			return nil, NewResourceConflict(fmt.Sprintf("Function already exist: %s", function.FunctionName))
		}
		return nil, err
	}

	if function.Publish {
		_, err = store.Functions.PublishVersion(created, "")
		if err != nil {
			return nil, err
		}
	}

	return s.toFunctionConfiguration(created), nil
}

// DeleteFunction deletes the specified Lambda function.
func (s *LambdaService) DeleteFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	qualifier := request.GetStringParam(req.Parameters, "Qualifier")
	if qualifier == "$LATEST" {
		return nil, NewInvalidParameter("Qualifier", "Cannot delete $LATEST version of a function")
	}

	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if qualifier != "" {
		if err := store.Functions.DeleteVersion(function.FunctionName, qualifier); err != nil {
			return nil, err
		}
		return response.EmptyResponse(), nil
	}

	mappings, err := store.EventSources.ListByFunction(function.FunctionArn)
	if err == nil && len(mappings) > 0 {
		return nil, ErrResourceInUse
	}

	if function.ContainerID != "" {
		if err := s.dockerClient.RemoveContainer(ctx, function.ContainerID, true); err != nil {
			logs.Warn("Failed to remove container for function", logs.String("containerID", function.ContainerID), logs.String("function", function.FunctionName), logs.Err(err))
		}
	}

	if err := store.Functions.Delete(function.FunctionName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
