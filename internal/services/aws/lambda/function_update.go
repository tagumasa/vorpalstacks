// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"encoding/base64"
	"errors"

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
		if err := validator.ValidateRoleForService(ctx, role, iam.ServicePrincipalLambda); err != nil {
			return nil, err
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

		if vpcMap := request.GetMapParam(req.Parameters, "VpcConfig"); vpcMap != nil {
			fn.VpcConfig = &lambdastore.VpcConfig{}
			if subnets, ok := vpcMap["SubnetIds"].([]interface{}); ok {
				for _, sub := range subnets {
					if str, ok := sub.(string); ok {
						fn.VpcConfig.SubnetIds = append(fn.VpcConfig.SubnetIds, str)
					}
				}
			}
			if sgs, ok := vpcMap["SecurityGroupIds"].([]interface{}); ok {
				for _, sg := range sgs {
					if str, ok := sg.(string); ok {
						fn.VpcConfig.SecurityGroupIds = append(fn.VpcConfig.SecurityGroupIds, str)
					}
				}
			}
		}

		if envMap := request.GetMapParam(req.Parameters, "Environment"); envMap != nil {
			fn.Environment = &lambdastore.Environment{}
			if vars, ok := envMap["Variables"].(map[string]interface{}); ok {
				fn.Environment.Variables = make(map[string]string)
				for k, v := range vars {
					if str, ok := v.(string); ok {
						fn.Environment.Variables[k] = str
					}
				}
			}
		}

		if dlMap := request.GetMapParam(req.Parameters, "DeadLetterConfig"); dlMap != nil {
			fn.DeadLetterConfig = &lambdastore.DeadLetterConfig{}
			if targetArn, ok := dlMap["TargetArn"].(string); ok {
				fn.DeadLetterConfig.TargetArn = targetArn
			}
		}

		if traceMap := request.GetMapParam(req.Parameters, "TracingConfig"); traceMap != nil {
			fn.TracingConfig = &lambdastore.TracingConfig{}
			if mode, ok := traceMap["Mode"].(string); ok {
				fn.TracingConfig.Mode = mode
			}
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
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.toFunctionConfiguration(function), nil
}
