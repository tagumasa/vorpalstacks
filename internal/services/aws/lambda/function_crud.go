// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"os"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// CreateFunction creates a new Lambda function.
func (s *LambdaService) CreateFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	functionName = extractFunctionName(functionName)

	runtime := request.GetStringParam(req.Parameters, "Runtime")
	role := request.GetStringParam(req.Parameters, "Role")
	handler := request.GetStringParam(req.Parameters, "Handler")
	packageType := request.GetStringParam(req.Parameters, "PackageType")

	// IAM role validation (requires reqCtx — not transport-agnostic).
	if role != "" {
		validator := reqCtx.GetIAMValidator()
		if os.Getenv("TEST_MODE") != "true" {
			if err := validator.ValidateRoleForService(ctx, role, iam.ServicePrincipalLambda); err != nil {
				return nil, err
			}
		}
	}

	// Code handling — fetch from S3, decode ZipFile, or accept ImageUri.
	// Results are passed to createFunctionCore as pre-processed metadata.
	codeMeta, imageUri, packageType, err := s.prepareCreateFunctionCodeCore(ctx, reqCtx.GetRegion(), functionName,
		request.GetMapParam(req.Parameters, "Code"), packageType)
	if err != nil {
		return nil, err
	}

	// Parse optional configurations.
	var ephemeralStorage *lambdastore.EphemeralStorage
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
			ephemeralStorage = &lambdastore.EphemeralStorage{Size: sizeValue}
		}
	}

	var vpcConfig *lambdastore.VpcConfig
	if request.GetMapParam(req.Parameters, "VpcConfig") != nil {
		vpcConfig = parseVpcConfig(req.Parameters)
		if vpcConfig != nil && len(vpcConfig.SubnetIds) > 0 {
			if err := s.resolveVpcConfig(ctx, reqCtx.GetRegion(), vpcConfig); err != nil {
				return nil, err
			}
		}
	}

	var environment *lambdastore.Environment
	if request.GetMapParam(req.Parameters, "Environment") != nil {
		environment = parseEnvironment(req.Parameters)
	}

	var deadLetterConfig *lambdastore.DeadLetterConfig
	if request.GetMapParam(req.Parameters, "DeadLetterConfig") != nil {
		dl, err := parseDeadLetterConfig(req.Parameters)
		if err != nil {
			return nil, err
		}
		deadLetterConfig = dl
	}

	var tracingConfig *lambdastore.TracingConfig
	if request.GetMapParam(req.Parameters, "TracingConfig") != nil {
		trace, err := parseTracingConfig(req.Parameters)
		if err != nil {
			return nil, err
		}
		tracingConfig = trace
	}

	var snapStart *lambdastore.SnapStart
	if snapMap := request.GetMapParam(req.Parameters, "SnapStart"); snapMap != nil {
		applyOn, _ := snapMap["ApplyOn"].(string)
		snapStart = &lambdastore.SnapStart{ApplyOn: applyOn}
	}

	var loggingConfig *lambdastore.LoggingConfig
	if logMap := request.GetMapParam(req.Parameters, "LoggingConfig"); logMap != nil {
		loggingConfig = parseLoggingConfig(logMap)
	}

	var imageConfig *lambdastore.ImageConfig
	if imgMap := request.GetMapParam(req.Parameters, "ImageConfig"); imgMap != nil {
		imageConfig = parseImageConfig(imgMap)
	}

	var fileSystemConfigs []lambdastore.FileSystemConfig
	if fscs, ok := req.Parameters["FileSystemConfigs"].([]interface{}); ok {
		for _, fsc := range fscs {
			if m, ok := fsc.(map[string]interface{}); ok {
				fileSystemConfigs = append(fileSystemConfigs, lambdastore.FileSystemConfig{
					Arn:            request.GetStringParam(m, "Arn"),
					LocalMountPath: request.GetStringParam(m, "LocalMountPath"),
				})
			}
		}
	}

	var layers []lambdastore.LayerReference
	if layerList, ok := req.Parameters["Layers"].([]interface{}); ok {
		for _, l := range layerList {
			if ls, ok := l.(string); ok {
				if !isValidLayerARN(ls) {
					return nil, NewInvalidParameter("Layers", "Invalid layer ARN format: "+ls)
				}
				layers = append(layers, lambdastore.LayerReference{Arn: ls})
			}
		}
	}

	var architectures []string
	if archs, ok := req.Parameters["Architectures"].([]interface{}); ok {
		architectures = make([]string, 0, len(archs))
		for _, a := range archs {
			if as, ok := a.(string); ok {
				architectures = append(architectures, as)
			}
		}
	}

	var tags map[string]string
	if tm, ok := req.Parameters["Tags"].(map[string]interface{}); ok {
		tags = tagutil.ToMap(tagutil.MapInterfaceToTags(tm))
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	created, published, err := s.createFunctionCore(store, &CreateFunctionInput{
		FunctionName:         functionName,
		Runtime:              runtime,
		Role:                 role,
		Handler:              handler,
		Description:          request.GetStringParam(req.Parameters, "Description"),
		PackageType:          packageType,
		KMSKeyArn:            request.GetStringParam(req.Parameters, "KMSKeyArn"),
		Publish:              request.GetBoolParam(req.Parameters, "Publish"),
		CodeSigningConfigArn: request.GetStringParam(req.Parameters, "CodeSigningConfigArn"),
		Region:               reqCtx.GetRegion(),
		CodeLocation:         codeMeta.CodeLocation,
		CodeSize:             codeMeta.CodeSize,
		CodeSha256:           codeMeta.CodeSha256,
		ImageUri:             imageUri,
		Timeout:              int32(request.GetIntParam(req.Parameters, "Timeout")),
		MemorySize:           int32(request.GetIntParam(req.Parameters, "MemorySize")),
		VpcConfig:            vpcConfig,
		Environment:          environment,
		DeadLetterConfig:     deadLetterConfig,
		TracingConfig:        tracingConfig,
		SnapStart:            snapStart,
		LoggingConfig:        loggingConfig,
		ImageConfig:          imageConfig,
		EphemeralStorage:     ephemeralStorage,
		FileSystemConfigs:    fileSystemConfigs,
		Layers:               layers,
		Architectures:        architectures,
		Tags:                 tags,
	})
	if err != nil {
		return nil, err
	}

	// When Publish was requested, the response describes the published
	// version rather than $LATEST.
	if published != nil {
		return s.toVersionConfiguration(published), nil
	}
	return s.toFunctionConfiguration(created), nil
}

// DeleteFunction deletes the specified Lambda function.
func (s *LambdaService) DeleteFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionNameRaw := request.GetStringParam(req.Parameters, "FunctionName")
	functionName, embeddedQualifier := resolveFunctionRef(functionNameRaw)
	if err := validateFunctionName(functionName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteFunctionCore(ctx, store, &DeleteFunctionInput{
		FunctionName: functionName,
		Qualifier:    mergeQualifier(request.GetStringParam(req.Parameters, "Qualifier"), embeddedQualifier),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
