// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"encoding/base64"
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

	if zipFileStr == "" && imageUri == "" && s3Bucket == "" {
		return nil, NewInvalidParameter("Code", "Either ZipFile, ImageUri, or S3Bucket/S3Key must be provided")
	}

	var codeLocation string
	var codeSize int64
	var codeSha256 string
	if zipFileStr != "" {
		zipFile, err := base64.StdEncoding.DecodeString(zipFileStr)
		if err != nil {
			return nil, NewInvalidParameter("ZipFile", "Invalid base64 encoding")
		}
		codeLocation, codeSize, err = s.storeCode(functionName, "$LATEST", zipFile, reqCtx.GetRegion())
		if err != nil {
			return nil, err
		}
		codeSha256 = lambdastore.GenerateCodeHash(zipFile)
	} else if s3Bucket != "" {
		if s3Key == "" {
			return nil, NewInvalidParameter("Code.S3Key", "S3Key is required when S3Bucket is specified")
		}
		zipFile, err := s.fetchCodeFromS3(ctx, s3Bucket, s3Key, reqCtx.GetRegion())
		if err != nil {
			return nil, NewInvalidParameter("Code", err.Error())
		}
		codeLocation, codeSize, err = s.storeCode(functionName, "$LATEST", zipFile, reqCtx.GetRegion())
		if err != nil {
			return nil, err
		}
		codeSha256 = lambdastore.GenerateCodeHash(zipFile)
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

	function, err := s.updateFunctionCodeCore(store, &UpdateFunctionCodeInput{
		FunctionName:  functionName,
		CodeLocation:  codeLocation,
		CodeSize:      codeSize,
		CodeSha256:    codeSha256,
		ImageUri:      imageUri,
		Architectures: architectures,
		Publish:       request.GetBoolParam(req.Parameters, "Publish"),
	})
	if err != nil {
		return nil, err
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
		FileSystemConfigs:    newFileSystemConfigs,
		Layers:               newLayers,
	})
	if err != nil {
		return nil, err
	}

	return s.toFunctionConfiguration(function), nil
}
