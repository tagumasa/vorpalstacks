package apigateway

import (
	"errors"
	"strconv"
	"strings"

	"vorpalstacks/internal/store/aws/apigateway"
)

// DeploymentInput is the transport-agnostic input for creating a
// deployment. An optional embedded stage is created atomically with the
// deployment; the core rolls the deployment back if stage creation fails.
type DeploymentInput struct {
	Description         string
	StageName           string
	StageDescription    string
	CacheClusterSize    string
	CacheClusterEnabled bool
	TracingEnabled      bool
	Variables           map[string]string
	CanarySettings      *CanarySettingsInput
	Tags                map[string]string
}

// StageInput is the transport-agnostic input for creating or replacing a
// stage.
type StageInput struct {
	StageName            string
	DeploymentId         string
	Description          string
	CacheClusterEnabled  bool
	CacheClusterSize     string
	DocumentationVersion string
	TracingEnabled       bool
	Variables            map[string]string
	Tags                 map[string]string
	CanarySettings       *CanarySettingsInput
}

// CanarySettingsInput is the transport-agnostic input for canary settings.
type CanarySettingsInput struct {
	PercentTraffic         float64
	DeploymentId           string
	StageVariableOverrides map[string]string
	UseStageCache          bool
}

// createDeploymentCore persists a deployment and, when StageName is set,
// also creates a matching stage. It uses the per-API key locker to
// serialise concurrent deployment attempts and rolls back the deployment
// if the embedded stage creation fails so that a validation failure does
// not leave an orphaned resource.
func (s *APIGatewayService) createDeploymentCore(
	stores *apiGatewayStores,
	apiId string,
	in *DeploymentInput,
) (*apigateway.Deployment, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if in.StageName != "" {
		if !validateStageName(in.StageName) {
			return nil, NewBadRequestException("Invalid stage name: must be alphanumeric, underscore, or hyphen, max 128 characters")
		}
		if !validateCacheClusterSize(in.CacheClusterSize) {
			return nil, NewBadRequestException("Invalid cacheClusterSize: must be one of 0.5, 1.6, 6.1, 13.5, 28.4, 58.2, 118, 237")
		}
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	deployment := &apigateway.Deployment{Description: in.Description}
	created, err := stores.restApis.CreateDeployment(apiId, deployment)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	if in.StageName != "" {
		stageDescription := in.StageDescription
		if stageDescription == "" {
			stageDescription = "Auto-created stage"
		}
		stage := &apigateway.Stage{
			StageName:           in.StageName,
			DeploymentId:        created.Id,
			Description:         stageDescription,
			CacheClusterEnabled: in.CacheClusterEnabled,
			CacheClusterSize:    in.CacheClusterSize,
			TracingEnabled:      in.TracingEnabled,
			Variables:           in.Variables,
			Tags:                tagsFromMap(in.Tags),
		}
		if in.CanarySettings != nil {
			stage.CanarySettings = canarySettingsFromInput(in.CanarySettings, created.Id)
		}
		if _, err := stores.restApis.CreateStage(apiId, stage); err != nil {
			// Compensating delete: roll back the deployment so that a
			// failed stage creation does not leave an orphaned resource.
			_ = stores.restApis.DeleteDeployment(apiId, created.Id)
			return nil, toApiGatewayError(err)
		}
	}

	return created, nil
}

// getDeploymentCore retrieves a deployment by id.
func (s *APIGatewayService) getDeploymentCore(stores *apiGatewayStores, apiId, deploymentId string) (*apigateway.Deployment, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if deploymentId == "" {
		return nil, NewBadRequestException("deploymentId is required")
	}
	deployment, err := stores.restApis.GetDeployment(apiId, deploymentId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return deployment, nil
}

// deleteDeploymentCore removes a deployment. Maps ErrDeploymentInUse to a
// ConflictException so callers can surface a clear 409 response.
func (s *APIGatewayService) deleteDeploymentCore(stores *apiGatewayStores, apiId, deploymentId string) error {
	if apiId == "" {
		return NewBadRequestException("restApiId is required")
	}
	if deploymentId == "" {
		return NewBadRequestException("deploymentId is required")
	}
	if err := stores.restApis.DeleteDeployment(apiId, deploymentId); err != nil {
		if errors.Is(err, apigateway.ErrDeploymentInUse) {
			return NewConflictException("Deployment is in use by a stage")
		}
		return toApiGatewayError(err)
	}
	return nil
}

// listDeploymentsCore returns all deployments for an api id.
func (s *APIGatewayService) listDeploymentsCore(stores *apiGatewayStores, apiId string) ([]*apigateway.Deployment, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	deployments, err := stores.restApis.ListDeployments(apiId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return deployments, nil
}

// createStageCore persists a stage.
func (s *APIGatewayService) createStageCore(stores *apiGatewayStores, apiId string, in *StageInput) (*apigateway.Stage, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if in.StageName == "" {
		return nil, NewBadRequestException("stageName is required")
	}
	if !validateStageName(in.StageName) {
		return nil, NewBadRequestException("Invalid stage name: must be alphanumeric, underscore, or hyphen, max 128 characters")
	}
	if !validateCacheClusterSize(in.CacheClusterSize) {
		return nil, NewBadRequestException("Invalid cacheClusterSize: must be one of 0.5, 1.6, 6.1, 13.5, 28.4, 58.2, 118, 237")
	}
	if in.DeploymentId == "" {
		return nil, NewBadRequestException("deploymentId is required")
	}

	stage := &apigateway.Stage{
		StageName:            in.StageName,
		DeploymentId:         in.DeploymentId,
		Description:          in.Description,
		CacheClusterEnabled:  in.CacheClusterEnabled,
		CacheClusterSize:     in.CacheClusterSize,
		DocumentationVersion: in.DocumentationVersion,
		TracingEnabled:       in.TracingEnabled,
		Variables:            in.Variables,
		Tags:                 tagsFromMap(in.Tags),
	}
	if in.CanarySettings != nil {
		stage.CanarySettings = canarySettingsFromInput(in.CanarySettings, "")
	}

	created, err := stores.restApis.CreateStage(apiId, stage)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getStageCore retrieves a stage by name.
func (s *APIGatewayService) getStageCore(stores *apiGatewayStores, apiId, stageName string) (*apigateway.Stage, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if stageName == "" {
		return nil, NewBadRequestException("stageName is required")
	}
	stage, err := stores.restApis.GetStage(apiId, stageName)
	if err != nil {
		return nil, ErrNotFoundException
	}
	return stage, nil
}

// deleteStageCore removes a stage and asks the runtime server to clean
// up any throttlers keyed by the stage name.
func (s *APIGatewayService) deleteStageCore(stores *apiGatewayStores, apiId, stageName string) error {
	if apiId == "" {
		return NewBadRequestException("restApiId is required")
	}
	if stageName == "" {
		return NewBadRequestException("stageName is required")
	}
	if err := stores.restApis.DeleteStage(apiId, stageName); err != nil {
		return ErrNotFoundException
	}
	if s.runtimeServer != nil {
		s.runtimeServer.CleanupStageThrottlers(stageName)
	}
	return nil
}

// listStagesCore returns all stages for an api id.
func (s *APIGatewayService) listStagesCore(stores *apiGatewayStores, apiId string) ([]*apigateway.Stage, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	return stores.restApis.ListStages(apiId)
}

// updateStageCore applies the patch operations to a stage under the
// per-API key locker and persists the result.
func (s *APIGatewayService) updateStageCore(
	stores *apiGatewayStores,
	apiId, stageName string,
	patches []PatchOperation,
) (*apigateway.Stage, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if stageName == "" {
		return nil, NewBadRequestException("stageName is required")
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	stage, err := stores.restApis.GetStage(apiId, stageName)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range patches {
		switch {
		case po.Path == "/description":
			stage.Description = po.Value
		case po.Path == "/deploymentId":
			if po.Value != "" {
				if _, err := stores.restApis.GetDeployment(apiId, po.Value); err != nil {
					return nil, NewBadRequestException("deployment not found")
				}
			}
			stage.DeploymentId = po.Value
		case po.Path == "/cacheClusterEnabled":
			stage.CacheClusterEnabled = po.Value == "true"
		case po.Path == "/cacheClusterSize":
			if !validateCacheClusterSize(po.Value) {
				return nil, NewBadRequestException("Invalid cacheClusterSize: must be one of 0.5, 1.6, 6.1, 13.5, 28.4, 58.2, 118, 237")
			}
			stage.CacheClusterSize = po.Value
		case po.Path == "/tracingEnabled":
			stage.TracingEnabled = po.Value == "true"
		case po.Path == "/clientCertificateId":
			stage.ClientCertificateId = po.Value
		case po.Path == "/documentationVersion":
			stage.DocumentationVersion = po.Value
		case strings.HasPrefix(po.Path, "/variables/"):
			if stage.Variables == nil {
				stage.Variables = make(map[string]string)
			}
			varName := strings.TrimPrefix(po.Path, "/variables/")
			if po.Op == "remove" {
				delete(stage.Variables, varName)
			} else {
				stage.Variables[varName] = po.Value
			}
		case strings.HasPrefix(po.Path, "/accessLogSettings/"):
			if stage.AccessLogSettings == nil {
				stage.AccessLogSettings = &apigateway.AccessLogSettings{}
			}
			switch strings.TrimPrefix(po.Path, "/accessLogSettings/") {
			case "destinationArn":
				if po.Value != "" && !validateAccessLogDestinationArn(po.Value) {
					return nil, NewBadRequestException("Invalid accessLogSettings destinationArn: must be a CloudWatch Logs log group ARN or a Kinesis Firehose delivery stream ARN whose name begins with amazon-apigateway-")
				}
				stage.AccessLogSettings.DestinationArn = po.Value
			case "format":
				stage.AccessLogSettings.Format = po.Value
			}
		case strings.HasPrefix(po.Path, "/methodSettings/"):
			if err := applyMethodSettingsPatch(stage, po); err != nil {
				return nil, toApiGatewayError(err)
			}
		case strings.HasPrefix(po.Path, "/canarySettings/"):
			if err := applyCanarySettingsPatch(stage, po); err != nil {
				return nil, toApiGatewayError(err)
			}
		}
	}

	if err := stores.restApis.UpdateStage(apiId, stage); err != nil {
		return nil, toApiGatewayError(err)
	}
	return stage, nil
}

func applyMethodSettingsPatch(stage *apigateway.Stage, po PatchOperation) error {
	if stage.MethodSettings == nil {
		stage.MethodSettings = make(map[string]*apigateway.MethodSetting)
	}

	rest := strings.TrimPrefix(po.Path, "/methodSettings/")
	parts := strings.SplitN(rest, "/", 3)

	resourcePath := parts[0]
	httpMethod := ""
	settingName := ""
	if len(parts) >= 2 {
		httpMethod = parts[1]
	}
	if len(parts) >= 3 {
		settingName = parts[2]
	}

	key := resourcePath + "/" + httpMethod

	// "remove" on the entire entry (no settingName) deletes the entry.
	if po.Op == "remove" && settingName == "" {
		delete(stage.MethodSettings, key)
		return nil
	}

	ms, ok := stage.MethodSettings[key]
	if !ok {
		ms = &apigateway.MethodSetting{}
		stage.MethodSettings[key] = ms
	}

	if po.Op == "remove" {
		switch settingName {
		case "metricsEnabled":
			ms.MetricsEnabled = false
		case "loggingLevel":
			ms.LoggingLevel = ""
		case "dataTraceEnabled":
			ms.DataTraceEnabled = false
		case "throttlingBurstLimit":
			ms.ThrottlingBurstLimit = 0
		case "throttlingRateLimit":
			ms.ThrottlingRateLimit = 0
		case "cachingEnabled":
			ms.CachingEnabled = false
		case "cacheTtlInSeconds":
			ms.CacheTtlInSeconds = 0
		case "cacheDataEncrypted":
			ms.CacheDataEncrypted = false
		case "requireAuthorizationForCacheControl":
			ms.RequireAuthorizationForCacheControl = false
		}
		if isMethodSettingEmpty(ms) {
			delete(stage.MethodSettings, key)
		}
		return nil
	}

	switch settingName {
	case "metricsEnabled":
		ms.MetricsEnabled = po.Value == "true"
	case "loggingLevel":
		if !validateLoggingLevel(po.Value) {
			return NewBadRequestException("Invalid loggingLevel: must be OFF, ERROR, or INFO")
		}
		ms.LoggingLevel = po.Value
	case "dataTraceEnabled":
		ms.DataTraceEnabled = po.Value == "true"
	case "throttlingBurstLimit":
		v, err := strconv.ParseInt(po.Value, 10, 32)
		if err != nil {
			return NewBadRequestException("Invalid throttlingBurstLimit: not a number")
		}
		if !validateMethodSettingThrottleBurstLimit(v) {
			return NewBadRequestException("throttlingBurstLimit must be between 0 and 100000")
		}
		ms.ThrottlingBurstLimit = int32(v)
	case "throttlingRateLimit":
		v, err := strconv.ParseFloat(po.Value, 64)
		if err != nil {
			return NewBadRequestException("Invalid throttlingRateLimit: not a number")
		}
		if !validateMethodSettingThrottleRateLimit(v) {
			return NewBadRequestException("throttlingRateLimit must be between 0 and 100000")
		}
		ms.ThrottlingRateLimit = v
	case "cachingEnabled":
		ms.CachingEnabled = po.Value == "true"
	case "cacheTtlInSeconds":
		v, err := strconv.ParseInt(po.Value, 10, 32)
		if err != nil {
			return NewBadRequestException("Invalid cacheTtlInSeconds: not a number")
		}
		if !validateCacheTtlInSeconds(int32(v)) {
			return NewBadRequestException("cacheTtlInSeconds must be between 0 and 86400")
		}
		ms.CacheTtlInSeconds = int32(v)
	case "cacheDataEncrypted":
		ms.CacheDataEncrypted = po.Value == "true"
	case "requireAuthorizationForCacheControl":
		ms.RequireAuthorizationForCacheControl = po.Value == "true"
	}
	return nil
}

func isMethodSettingEmpty(ms *apigateway.MethodSetting) bool {
	return !ms.MetricsEnabled &&
		ms.LoggingLevel == "" &&
		!ms.DataTraceEnabled &&
		ms.ThrottlingBurstLimit == 0 &&
		ms.ThrottlingRateLimit == 0 &&
		!ms.CachingEnabled &&
		ms.CacheTtlInSeconds == 0 &&
		!ms.CacheDataEncrypted &&
		!ms.RequireAuthorizationForCacheControl &&
		len(ms.UnreservedCacheParameters) == 0
}

func applyCanarySettingsPatch(stage *apigateway.Stage, po PatchOperation) error {
	if stage.CanarySettings == nil {
		stage.CanarySettings = &apigateway.CanarySettings{}
	}

	rest := strings.TrimPrefix(po.Path, "/canarySettings/")
	parts := strings.SplitN(rest, "/", 2)

	switch parts[0] {
	case "percentTraffic":
		if po.Op == "remove" {
			stage.CanarySettings.PercentTraffic = 0
		} else {
			v, err := strconv.ParseFloat(po.Value, 64)
			if err != nil {
				return NewBadRequestException("Invalid percentTraffic: not a number")
			}
			if !validatePercentTraffic(v) {
				return NewBadRequestException("percentTraffic must be between 0 and 100")
			}
			stage.CanarySettings.PercentTraffic = v
		}
	case "deploymentId":
		if po.Op == "remove" {
			stage.CanarySettings.DeploymentId = ""
		} else {
			stage.CanarySettings.DeploymentId = po.Value
		}
	case "useStageCache":
		if po.Op == "remove" {
			stage.CanarySettings.UseStageCache = false
		} else {
			stage.CanarySettings.UseStageCache = po.Value == "true"
		}
	case "stageVariableOverrides":
		if len(parts) < 2 {
			return nil
		}
		if stage.CanarySettings.StageVariableOverrides == nil {
			stage.CanarySettings.StageVariableOverrides = make(map[string]string)
		}
		varName := parts[1]
		if po.Op == "remove" {
			delete(stage.CanarySettings.StageVariableOverrides, varName)
		} else {
			stage.CanarySettings.StageVariableOverrides[varName] = po.Value
		}
	}
	return nil
}
