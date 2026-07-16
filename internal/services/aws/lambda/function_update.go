// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// UpdateFunctionCode updates the code of the specified Lambda function.
func (s *LambdaService) UpdateFunctionCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	functionName = extractFunctionName(functionName)

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

	var zipFile []byte
	var codePath string
	var codeSize int64
	if zipFileStr != "" {
		zipFile, err = base64.StdEncoding.DecodeString(zipFileStr)
		if err != nil {
			return nil, NewInvalidParameter("ZipFile", "Invalid base64 encoding")
		}
		codePath, codeSize, err = s.storeCode(functionName, "$LATEST", zipFile, reqCtx.GetRegion())
		if err != nil {
			return nil, err
		}
	} else if s3Bucket != "" {
		if s3Key == "" {
			return nil, NewInvalidParameter("Code.S3Key", "S3Key is required when S3Bucket is specified")
		}
		zipFile, err = s.fetchCodeFromS3(ctx, s3Bucket, s3Key, reqCtx.GetRegion())
		if err != nil {
			return nil, NewInvalidParameter("Code", err.Error())
		}
		codePath, codeSize, err = s.storeCode(functionName, "$LATEST", zipFile, reqCtx.GetRegion())
		if err != nil {
			return nil, err
		}
	}

	archs := request.GetSliceParam(req.Parameters, "Architectures")

	function, err := store.Functions.UpdateAtomically(functionName, func(fn *lambdastore.Function) error {
		if zipFileStr != "" {
			fn.CodeLocation = codePath
			fn.CodeSize = codeSize
			fn.CodeSha256 = lambdastore.GenerateCodeHash(zipFile)
		}

		if imageUri != "" {
			fn.ImageUri = imageUri
			fn.PackageType = "Image"
		}

		if len(archs) > 0 {
			fn.Architectures = make([]string, 0, len(archs))
			for _, a := range archs {
				if as, ok := a.(string); ok {
					if as != "x86_64" && as != "arm64" {
						return NewInvalidParameter("Architectures", "must be x86_64 or arm64")
					}
					fn.Architectures = append(fn.Architectures, as)
				}
			}
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return nil, NewResourceNotFound("Function", functionName)
		}
		return nil, err
	}

	publish := request.GetBoolParam(req.Parameters, "Publish")
	if publish {
		_, err = store.Functions.PublishVersion(function, "")
		if err != nil {
			return nil, err
		}
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	runtime := request.GetStringParam(req.Parameters, "Runtime")
	if runtime != "" && !ValidateRuntime(runtime) {
		return nil, NewInvalidParameter("Runtime", "Runtime '"+runtime+"' is not supported")
	}

	role := request.GetStringParam(req.Parameters, "Role")
	if role != "" {
		validator := reqCtx.GetIAMValidator()
		if os.Getenv("TEST_MODE") != "true" {
			if err := validator.ValidateRoleForService(ctx, role, iam.ServicePrincipalLambda); err != nil {
				return nil, err
			}
		}
	}

	timeout := request.GetIntParam(req.Parameters, "Timeout")
	if timeout > 0 {
		if err := validateTimeout(int32(timeout)); err != nil {
			return nil, err
		}
	}
	memorySize := request.GetIntParam(req.Parameters, "MemorySize")
	if memorySize > 0 {
		if err := validateMemorySize(int32(memorySize)); err != nil {
			return nil, err
		}
	}

	var oldContainerID string
	function, err := store.Functions.UpdateAtomically(functionName, func(fn *lambdastore.Function) error {
		if runtime != "" {
			fn.Runtime = lambdastore.Runtime(runtime)
		}
		if role != "" {
			fn.Role = role
		}
		if handler := request.GetStringParam(req.Parameters, "Handler"); handler != "" {
			fn.Handler = handler
		}
		if desc := request.GetStringParam(req.Parameters, "Description"); desc != "" {
			fn.Description = desc
		}
		if timeout > 0 {
			fn.Timeout = int32(timeout)
		}
		if memorySize > 0 {
			fn.MemorySize = int32(memorySize)
		}
		if kmsKeyArn := request.GetStringParam(req.Parameters, "KMSKeyArn"); kmsKeyArn != "" {
			fn.KMSKeyArn = kmsKeyArn
		}

		if request.GetMapParam(req.Parameters, "VpcConfig") != nil {
			fn.VpcConfig = parseVpcConfig(req.Parameters)
		}

		if request.GetMapParam(req.Parameters, "Environment") != nil {
			fn.Environment = parseEnvironment(req.Parameters)
		}

		if request.GetMapParam(req.Parameters, "DeadLetterConfig") != nil {
			dl, err := parseDeadLetterConfig(req.Parameters)
			if err != nil {
				return err
			}
			fn.DeadLetterConfig = dl
		}

		if request.GetMapParam(req.Parameters, "TracingConfig") != nil {
			trace, err := parseTracingConfig(req.Parameters)
			if err != nil {
				return err
			}
			fn.TracingConfig = trace
		}

		if layers, ok := req.Parameters["Layers"].([]interface{}); ok {
			fn.Layers = make([]lambdastore.LayerReference, 0, len(layers))
			for _, l := range layers {
				if ls, ok := l.(string); ok {
					if !isValidLayerARN(ls) {
						return NewInvalidParameter("Layers", "Invalid layer ARN format: "+ls)
					}
					fn.Layers = append(fn.Layers, lambdastore.LayerReference{Arn: ls})
				}
			}
		}

		// Invalidate any running container so the next invoke creates a fresh
		// one with the updated runtime, handler, memory, or environment.
		oldContainerID = fn.ContainerID
		fn.ContainerID = ""
		fn.ContainerImageID = ""
		return nil
	})
	if err != nil {
		if IsLambdaError(err) {
			return nil, err
		}
		return nil, ErrResourceNotFound
	}

	// Remove the previous container to prevent orphaned Docker containers.
	// The ID was captured inside the callback before it was cleared.
	if oldContainerID != "" {
		if rmErr := s.dockerClient.RemoveContainer(ctx, oldContainerID, true); rmErr != nil {
			logs.Warn("Failed to remove container after configuration update",
				logs.String("containerID", oldContainerID),
				logs.String("function", functionName),
				logs.Err(rmErr))
		}
	}

	return s.toFunctionConfiguration(function), nil
}
