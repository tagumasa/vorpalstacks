package apigateway

import (
	"encoding/json"
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
			return nil, NewBadRequestException(validCacheClusterSizesMessage())
		}
		if err := validateStageVariables(in.Variables); err != nil {
			return nil, err
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
			canary, err := canarySettingsFromInput(in.CanarySettings, created.Id)
			if err != nil {
				// Compensating delete, as below: a validation failure must
				// not leave an orphaned deployment.
				_ = stores.restApis.DeleteDeployment(apiId, created.Id)
				return nil, err
			}
			stage.CanarySettings = canary
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

// updateDeploymentCore applies JSON Patch operations to a deployment under
// the api key lock.
func (s *APIGatewayService) updateDeploymentCore(
	stores *apiGatewayStores,
	apiId, deploymentId string,
	ops []PatchOperation,
) (*apigateway.Deployment, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if deploymentId == "" {
		return nil, NewBadRequestException("deploymentId is required")
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	deployment, err := stores.restApis.GetDeployment(apiId, deploymentId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range ops {
		handled := false
		switch po.Path {
		case "/description":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			deployment.Description = po.Value
		}
		if !handled {
			return nil, unknownPatchPathError(po)
		}
	}

	if err := stores.restApis.UpdateDeployment(apiId, deployment); err != nil {
		return nil, toApiGatewayError(err)
	}

	return deployment, nil
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
		return nil, NewBadRequestException(validCacheClusterSizesMessage())
	}
	if in.DeploymentId == "" {
		return nil, NewBadRequestException("deploymentId is required")
	}
	if err := validateStageVariables(in.Variables); err != nil {
		return nil, err
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
		canary, err := canarySettingsFromInput(in.CanarySettings, "")
		if err != nil {
			return nil, err
		}
		stage.CanarySettings = canary
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
		return nil, toApiGatewayError(err)
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
		return toApiGatewayError(err)
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
		handled := false
		switch {
		case po.Path == "/description":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			stage.Description = po.Value
		case po.Path == "/deploymentId":
			handled = true
			// The official row documents replace and copy; the copy form
			// promotes the canary deployment — the documented example in
			// the PatchOperation from member description.
			if po.Op == "copy" {
				if po.From != "/canarySettings/deploymentId" {
					return nil, unknownPatchPathError(po)
				}
				if stage.CanarySettings == nil || stage.CanarySettings.DeploymentId == "" {
					return nil, NewBadRequestException("cannot promote the canary deployment: the stage has no canary deploymentId")
				}
				stage.DeploymentId = stage.CanarySettings.DeploymentId
			} else {
				if err := requirePatchOp(po, opReplace); err != nil {
					return nil, err
				}
				if po.Value != "" {
					if _, err := stores.restApis.GetDeployment(apiId, po.Value); err != nil {
						return nil, NewBadRequestException("deployment not found")
					}
				}
				stage.DeploymentId = po.Value
			}
		case po.Path == "/cacheClusterEnabled":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			stage.CacheClusterEnabled = po.Value == "true"
		case po.Path == "/cacheClusterSize":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if !validateCacheClusterSize(po.Value) {
				return nil, NewBadRequestException(validCacheClusterSizesMessage())
			}
			stage.CacheClusterSize = po.Value
		case po.Path == "/tracingEnabled":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			stage.TracingEnabled = po.Value == "true"
		case po.Path == "/clientCertificateId":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			stage.ClientCertificateId = po.Value
		case po.Path == "/documentationVersion":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			stage.DocumentationVersion = po.Value
		case po.Path == "/variables":
			handled = true
			// Whole-member row of the official UpdateStage patch table:
			// replace sets the map from the JSON object value, remove
			// clears it; every other operation is unsupported there.
			if err := requirePatchOp(po, opReplace|opRemove); err != nil {
				return nil, err
			}
			if po.Op == "remove" {
				stage.Variables = nil
			} else {
				parsed, err := parseWholeStringMapValue(po, validateStageVariableName, validateStageVariableValue)
				if err != nil {
					return nil, err
				}
				stage.Variables = parsed
			}
		case strings.HasPrefix(po.Path, "/variables/"):
			handled = true
			// The /variables/* row documents replace only, and the
			// per-key form enforces the documented stage variable name
			// and value constraints on the unescaped key and the value.
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if stage.Variables == nil {
				stage.Variables = make(map[string]string)
			}
			if err := applyMapPatch(stage.Variables, po, "/variables/", validateStageVariableName, validateStageVariableValue); err != nil {
				return nil, err
			}
		case po.Path == "/accessLogSettings":
			handled = true
			// The whole-member row documents remove only.
			if err := requirePatchOp(po, opRemove); err != nil {
				return nil, err
			}
			stage.AccessLogSettings = nil
		case strings.HasPrefix(po.Path, "/accessLogSettings/"):
			handled = true
			// The destinationArn and format rows document add, replace
			// and remove.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
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
			default:
				return nil, unknownPatchPathError(po)
			}
		case po.Path == "/methodSettings":
			handled = true
			// The whole-member row documents replace only.
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			settings, err := parseWholeMethodSettingsValue(po)
			if err != nil {
				return nil, err
			}
			stage.MethodSettings = settings
		case po.Path == "/canarySettings":
			handled = true
			// The whole-member row documents remove only.
			if err := requirePatchOp(po, opRemove); err != nil {
				return nil, err
			}
			stage.CanarySettings = nil
		case strings.HasPrefix(po.Path, "/canarySettings/"):
			handled = true
			if err := applyCanarySettingsPatch(stage, po); err != nil {
				return nil, toApiGatewayError(err)
			}
		}
		if !handled {
			// The documented setting family addresses members as
			// /{resourcePath}/{httpMethod}/{group}/{member}, which no fixed
			// prefix can single out — parse it before rejecting the path.
			setting, ok := parseMethodSettingsSettingPath(po.Path)
			if !ok {
				return nil, unknownPatchPathError(po)
			}
			if err := applyMethodSettingsSetting(stage, po, setting); err != nil {
				return nil, err
			}
		}
	}

	if err := stores.restApis.UpdateStage(apiId, stage); err != nil {
		return nil, toApiGatewayError(err)
	}
	return stage, nil
}

// methodSettingPatchValue mirrors the JSON form of one MethodSetting entry
// inside a whole-member /methodSettings replace value: AWS member names,
// pointer members so an omitted field keeps the zero value distinct from
// an explicitly set one.
type methodSettingPatchValue struct {
	MetricsEnabled                      *bool             `json:"metricsEnabled"`
	LoggingLevel                        *string           `json:"loggingLevel"`
	DataTraceEnabled                    *bool             `json:"dataTraceEnabled"`
	ThrottlingBurstLimit                *int64            `json:"throttlingBurstLimit"`
	ThrottlingRateLimit                 *float64          `json:"throttlingRateLimit"`
	CachingEnabled                      *bool             `json:"cachingEnabled"`
	CacheTtlInSeconds                   *int64            `json:"cacheTtlInSeconds"`
	CacheDataEncrypted                  *bool             `json:"cacheDataEncrypted"`
	RequireAuthorizationForCacheControl *bool             `json:"requireAuthorizationForCacheControl"`
	UnreservedCacheParameters           map[string]string `json:"unreservedCacheParameters"`
}

// parseWholeMethodSettingsValue decodes the whole-member /methodSettings
// replace value: a JSON object whose keys are {resource_path}/{http_method}
// method-setting keys with real slashes (they are data, not JSON Pointer
// tokens, so no unescaping applies) and whose entries use the AWS member
// names. Every entry runs through the same validators as the per-setting
// patch form.
func parseWholeMethodSettingsValue(po PatchOperation) (map[string]*apigateway.MethodSetting, error) {
	var raw map[string]methodSettingPatchValue
	if err := json.Unmarshal([]byte(po.Value), &raw); err != nil {
		return nil, NewBadRequestException(
			"Invalid patch value for '/methodSettings': expected a JSON object of method-setting keys to settings objects")
	}
	parsed := make(map[string]*apigateway.MethodSetting, len(raw))
	for key, entry := range raw {
		if key == "" {
			return nil, NewBadRequestException(
				"Invalid patch value for '/methodSettings': entry keys must not be empty")
		}
		ms := &apigateway.MethodSetting{
			UnreservedCacheParameters: entry.UnreservedCacheParameters,
		}
		if entry.MetricsEnabled != nil {
			ms.MetricsEnabled = *entry.MetricsEnabled
		}
		if entry.LoggingLevel != nil {
			if !validateLoggingLevel(*entry.LoggingLevel) {
				return nil, NewBadRequestException("Invalid loggingLevel: must be OFF, ERROR, or INFO")
			}
			ms.LoggingLevel = *entry.LoggingLevel
		}
		if entry.DataTraceEnabled != nil {
			ms.DataTraceEnabled = *entry.DataTraceEnabled
		}
		if entry.ThrottlingBurstLimit != nil {
			if !validateMethodSettingThrottleBurstLimit(*entry.ThrottlingBurstLimit) {
				return nil, NewBadRequestException("throttlingBurstLimit must be between 0 and 100000")
			}
			ms.ThrottlingBurstLimit = int32(*entry.ThrottlingBurstLimit)
		}
		if entry.ThrottlingRateLimit != nil {
			if !validateMethodSettingThrottleRateLimit(*entry.ThrottlingRateLimit) {
				return nil, NewBadRequestException("throttlingRateLimit must be between 0 and 100000")
			}
			ms.ThrottlingRateLimit = *entry.ThrottlingRateLimit
		}
		if entry.CachingEnabled != nil {
			ms.CachingEnabled = *entry.CachingEnabled
		}
		if entry.CacheTtlInSeconds != nil {
			if !validateCacheTtlInSeconds(int32(*entry.CacheTtlInSeconds)) {
				return nil, NewBadRequestException("cacheTtlInSeconds must be between 0 and 86400")
			}
			ms.CacheTtlInSeconds = int32(*entry.CacheTtlInSeconds)
		}
		if entry.CacheDataEncrypted != nil {
			ms.CacheDataEncrypted = *entry.CacheDataEncrypted
		}
		if entry.RequireAuthorizationForCacheControl != nil {
			ms.RequireAuthorizationForCacheControl = *entry.RequireAuthorizationForCacheControl
		}
		parsed[key] = ms
	}
	return parsed, nil
}

// methodSettingsSettingPath is one parsed setting address of the documented
// UpdateStage patch family /{resourcePath}/{httpMethod}/{group}/{member} —
// e.g. /pets/GET/logging/loglevel, the root as //GET/metrics/enabled, an
// escaped resource path as /~1pets~1{petId}/GET/throttling/burstLimit, and
// the all-methods wildcard /*/*/logging/dataTrace of the CLI reference
// example. Every row of the family documents replace only.
type methodSettingsSettingPath struct {
	resourceTokens []string
	httpMethod     string
	group          string
	member         string
}

// methodSettingsGroups are the setting groups the documented family
// addresses; the member token is validated against the row set instead.
var methodSettingsGroups = map[string]bool{
	"logging":    true,
	"metrics":    true,
	"throttling": true,
	"caching":    true,
}

// isMethodSettingsMethodToken reports whether a path token addresses the
// httpMethod position: an HTTP verb (ANY included) or the documented
// all-methods wildcard "*".
func isMethodSettingsMethodToken(token string) bool {
	return token == "*" || validHTTPMethods[token]
}

// parseMethodSettingsSettingPath matches a path against the documented
// setting family. One or more resource tokens precede the method token,
// which a group token and a member token follow exactly; anything else is
// not a row of the family and rejects as an unknown patch path.
func parseMethodSettingsSettingPath(path string) (methodSettingsSettingPath, bool) {
	tokens := splitPatchTokens(path)
	for i := 1; i+2 < len(tokens); i++ {
		if !isMethodSettingsMethodToken(tokens[i]) || !methodSettingsGroups[tokens[i+1]] {
			continue
		}
		if len(tokens) != i+3 {
			return methodSettingsSettingPath{}, false
		}
		return methodSettingsSettingPath{
			resourceTokens: tokens[:i],
			httpMethod:     tokens[i],
			group:          tokens[i+1],
			member:         tokens[i+2],
		}, true
	}
	return methodSettingsSettingPath{}, false
}

// applyMethodSettingsSetting applies one documented setting replace. The
// settings map is keyed by the as-addressed method key (methodMapKey): the
// official CLI update-stage example output shows the keys as the addressed
// pointer token ("~1resourceName/GET") and the wildcard ("*/*").
func applyMethodSettingsSetting(stage *apigateway.Stage, po PatchOperation, p methodSettingsSettingPath) *ApiGatewayError {
	if err := requirePatchOp(po, opReplace); err != nil {
		return err
	}
	if stage.MethodSettings == nil {
		stage.MethodSettings = make(map[string]*apigateway.MethodSetting)
	}
	key := methodMapKey(p.resourceTokens, p.httpMethod)
	ms, ok := stage.MethodSettings[key]
	if !ok {
		ms = &apigateway.MethodSetting{}
		stage.MethodSettings[key] = ms
	}
	switch p.group + "/" + p.member {
	case "logging/loglevel":
		if !validateLoggingLevel(po.Value) {
			return NewBadRequestException("Invalid loggingLevel: must be OFF, ERROR, or INFO")
		}
		ms.LoggingLevel = po.Value
	case "logging/dataTrace":
		ms.DataTraceEnabled = po.Value == "true"
	case "metrics/enabled":
		ms.MetricsEnabled = po.Value == "true"
	case "throttling/burstLimit":
		v, err := strconv.ParseInt(po.Value, 10, 32)
		if err != nil {
			return NewBadRequestException("Invalid throttlingBurstLimit: not a number")
		}
		if !validateMethodSettingThrottleBurstLimit(v) {
			return NewBadRequestException("throttlingBurstLimit must be between 0 and 100000")
		}
		ms.ThrottlingBurstLimit = int32(v)
	case "throttling/rateLimit":
		v, err := strconv.ParseFloat(po.Value, 64)
		if err != nil {
			return NewBadRequestException("Invalid throttlingRateLimit: not a number")
		}
		if !validateMethodSettingThrottleRateLimit(v) {
			return NewBadRequestException("throttlingRateLimit must be between 0 and 100000")
		}
		ms.ThrottlingRateLimit = v
	case "caching/enabled":
		ms.CachingEnabled = po.Value == "true"
	case "caching/dataEncrypted":
		ms.CacheDataEncrypted = po.Value == "true"
	case "caching/ttlInSeconds":
		v, err := strconv.ParseInt(po.Value, 10, 32)
		if err != nil {
			return NewBadRequestException("Invalid cacheTtlInSeconds: not a number")
		}
		if !validateCacheTtlInSeconds(int32(v)) {
			return NewBadRequestException("cacheTtlInSeconds must be between 0 and 86400")
		}
		ms.CacheTtlInSeconds = int32(v)
	case "caching/requireAuthorizationForCacheControl":
		ms.RequireAuthorizationForCacheControl = po.Value == "true"
	case "caching/unauthorizedCacheControlHeaderStrategy":
		if !validateUnauthorizedCacheControlHeaderStrategy(po.Value) {
			return NewBadRequestException("Invalid unauthorizedCacheControlHeaderStrategy: must be FAIL_WITH_403, SUCCEED_WITH_RESPONSE_HEADER, or SUCCEED_WITHOUT_RESPONSE_HEADER")
		}
		ms.UnauthorizedCacheControlHeaderStrategy = po.Value
	default:
		return unknownPatchPathError(po)
	}
	return nil
}

func applyCanarySettingsPatch(stage *apigateway.Stage, po PatchOperation) error {
	if stage.CanarySettings == nil {
		stage.CanarySettings = &apigateway.CanarySettings{}
	}

	rest := strings.TrimPrefix(po.Path, "/canarySettings/")
	parts := strings.SplitN(rest, "/", 2)

	switch parts[0] {
	case "percentTraffic":
		// The row documents replace only.
		if err := requirePatchOp(po, opReplace); err != nil {
			return err
		}
		v, err := strconv.ParseFloat(po.Value, 64)
		if err != nil {
			return NewBadRequestException("Invalid percentTraffic: not a number")
		}
		if !validatePercentTraffic(v) {
			return NewBadRequestException("percentTraffic must be between 0 and 100")
		}
		stage.CanarySettings.PercentTraffic = v
	case "deploymentId":
		// The row documents replace only.
		if err := requirePatchOp(po, opReplace); err != nil {
			return err
		}
		stage.CanarySettings.DeploymentId = po.Value
	case "useStageCache":
		// The row documents replace only.
		if err := requirePatchOp(po, opReplace); err != nil {
			return err
		}
		stage.CanarySettings.UseStageCache = po.Value == "true"
	case "stageVariableOverrides":
		if len(parts) < 2 {
			// Whole-member form: the official patch table documents
			// replace only. The overrides are stage variables, so the
			// documented name and value constraints apply.
			if err := requirePatchOp(po, opReplace); err != nil {
				return err
			}
			parsed, err := parseWholeStringMapValue(po, validateStageVariableName, validateStageVariableValue)
			if err != nil {
				return err
			}
			stage.CanarySettings.StageVariableOverrides = parsed
			return nil
		}
		if stage.CanarySettings.StageVariableOverrides == nil {
			stage.CanarySettings.StageVariableOverrides = make(map[string]string)
		}
		// The per-key override form admits the same operations the
		// /variables/* row admits for stage variables: replace only —
		// the overrides are stage variables and no other row governs
		// their per-key form.
		if err := requirePatchOp(po, opReplace); err != nil {
			return err
		}
		return applyMapPatch(stage.CanarySettings.StageVariableOverrides, po, "/canarySettings/stageVariableOverrides/", validateStageVariableName, validateStageVariableValue)
	default:
		return unknownPatchPathError(po)
	}
	return nil
}
