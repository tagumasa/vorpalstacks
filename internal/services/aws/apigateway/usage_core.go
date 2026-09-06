package apigateway

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/store/aws/common"
)

// ApiKeyInput is the transport-agnostic input for creating an API key.
// GenerateDistinctId mirrors the AWS member semantics: nil means true. When
// false the caller-supplied Value becomes the key id and must be 20-128
// characters; when true the store mints the id and the value is discarded.
type ApiKeyInput struct {
	Name               string
	Description        string
	Enabled            bool
	CustomerId         string
	Value              string
	Id                 string
	GenerateDistinctId *bool
	StageKeys          []string
	Tags               []types.Tag
}

// UsagePlanInput is the transport-agnostic input for creating a usage plan.
type UsagePlanInput struct {
	Name        string
	Description string
	ApiStages   []ApiStageInput
	Quota       *QuotaInput
	Throttle    *ThrottleInput
	Tags        []types.Tag
}

// ApiStageInput describes a single api-stage attachment on a usage plan.
type ApiStageInput struct {
	ApiId    string
	Stage    string
	Throttle map[string]*apigateway.Throttle
}

// QuotaInput is the transport-agnostic quota settings for a usage plan.
type QuotaInput struct {
	Limit  int64
	Offset int64
	Period string
}

// ThrottleInput is the transport-agnostic throttle settings for a usage
// plan.
type ThrottleInput struct {
	BurstLimit int64
	RateLimit  float64
}

// UsagePlanKeyInput is the transport-agnostic input for attaching an API
// key to a usage plan.
type UsagePlanKeyInput struct {
	KeyId   string
	KeyType string
}

// createApiKeyCore persists an API key. The name is optional in the API
// Gateway model, so both planes accept an empty name; the core validates
// the stageKey entries first and then the generateDistinctId/value
// pairing, preserving the data-plane failure precedence.
func (s *APIGatewayService) createApiKeyCore(stores *apiGatewayStores, in *ApiKeyInput) (*apigateway.ApiKey, error) {
	for _, sk := range in.StageKeys {
		if !validateStageKey(sk) {
			return nil, NewBadRequestException("invalid stageKey format, expected restApiId/stageName: " + sk)
		}
	}
	generateDistinctId := true
	if in.GenerateDistinctId != nil {
		generateDistinctId = *in.GenerateDistinctId
	}
	if !generateDistinctId {
		if in.Value == "" {
			return nil, NewBadRequestException("value is required when generateDistinctId is false")
		}
		if len(in.Value) < 20 || len(in.Value) > 128 {
			return nil, NewBadRequestException("value must be between 20 and 128 characters")
		}
		in.Id = in.Value
	} else {
		in.Value = ""
	}

	apiKey := &apigateway.ApiKey{
		Name:        in.Name,
		Description: in.Description,
		Enabled:     in.Enabled,
		CustomerId:  in.CustomerId,
		Value:       in.Value,
		Id:          in.Id,
		StageKeys:   in.StageKeys,
		Tags:        in.Tags,
	}

	created, err := stores.usage.CreateApiKey(apiKey)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getApiKeyCore retrieves a single API key.
func (s *APIGatewayService) getApiKeyCore(stores *apiGatewayStores, apiKeyId string) (*apigateway.ApiKey, error) {
	if apiKeyId == "" {
		return nil, NewBadRequestException("apiKey is required")
	}
	return stores.usage.GetApiKey(apiKeyId)
}

// deleteApiKeyCore removes an API key and asks the runtime server to drop
// any cached state keyed by the api key id.
func (s *APIGatewayService) deleteApiKeyCore(stores *apiGatewayStores, apiKeyId string) error {
	if apiKeyId == "" {
		return NewBadRequestException("apiKey is required")
	}
	if err := stores.usage.DeleteApiKey(apiKeyId); err != nil {
		return toApiGatewayError(err)
	}
	if s.runtimeServer != nil {
		s.runtimeServer.RemoveApiKey(apiKeyId)
	}
	return nil
}

// listApiKeysCore returns a page of API keys.
func (s *APIGatewayService) listApiKeysCore(stores *apiGatewayStores, limit int, marker string) (*common.ListResult[apigateway.ApiKey], error) {
	resolved, err := resolvePageLimit(limit)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return stores.usage.ListApiKeys(common.ListOptions{
		Marker:   marker,
		MaxItems: resolved,
	})
}

// updateApiKeyCore applies patch operations to an API key under the
// per-key key locker and asks the runtime server to refresh its cache.
func (s *APIGatewayService) updateApiKeyCore(
	stores *apiGatewayStores,
	apiKeyId string,
	patches []PatchOperation,
) (*apigateway.ApiKey, error) {
	if apiKeyId == "" {
		return nil, NewBadRequestException("apiKey is required")
	}

	stores.keyLocker.Lock(apiKeyId)
	defer stores.keyLocker.Unlock(apiKeyId)

	apiKey, err := stores.usage.GetApiKey(apiKeyId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range patches {
		handled := false
		switch {
		case po.Path == "/name":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			apiKey.Name = po.Value
		case po.Path == "/description":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			apiKey.Description = po.Value
		case po.Path == "/customerId":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			apiKey.CustomerId = po.Value
		case po.Path == "/enabled":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			apiKey.Enabled = po.Value == "true"
		case po.Path == "/stages":
			handled = true
			// The /stages row of the official UpdateApiKey patch table
			// documents add and remove; the value uses the stageKeys
			// member format, restApiId/stageName, that the official CLI
			// reference output shows.
			if err := requirePatchOp(po, opAdd|opRemove); err != nil {
				return nil, err
			}
			if !validateStageKey(po.Value) {
				return nil, NewBadRequestException("invalid stage format, expected restApiId/stageName: " + po.Value)
			}
			if po.Op == "add" {
				if apiKey.StageKeys == nil {
					apiKey.StageKeys = []string{}
				}
				if !slices.Contains(apiKey.StageKeys, po.Value) {
					apiKey.StageKeys = append(apiKey.StageKeys, po.Value)
				}
			} else {
				apiKey.StageKeys = removeString(apiKey.StageKeys, po.Value)
			}
		case strings.HasPrefix(po.Path, "/stageKeys/"):
			handled = true
			// Not a documented patch form: the table addresses the stages
			// through the /stages row only.
			return nil, unknownPatchPathError(po)
		case po.Path == "/labels":
			handled = true
			// The /labels row of the official UpdateApiKey patch table
			// documents add and remove. ApiKey carries no labels member —
			// tags is the only string-to-string map the model and the API
			// reference define, so the row addresses the tags: add sets
			// them from the JSON object value, remove clears them.
			if err := requirePatchOp(po, opAdd|opRemove); err != nil {
				return nil, err
			}
			if po.Op == "remove" {
				apiKey.Tags = nil
			} else {
				parsed, err := parseWholeStringMapValue(po, nil, nil)
				if err != nil {
					return nil, err
				}
				apiKey.Tags = tagsFromMap(parsed)
			}
		}
		if !handled {
			return nil, unknownPatchPathError(po)
		}
	}
	if err := stores.usage.UpdateApiKey(apiKey); err != nil {
		return nil, toApiGatewayError(err)
	}

	if s.runtimeServer != nil {
		s.runtimeServer.RemoveApiKey(apiKeyId)
	}
	return apiKey, nil
}

// createUsagePlanCore persists a usage plan.
func (s *APIGatewayService) createUsagePlanCore(stores *apiGatewayStores, in *UsagePlanInput) (*apigateway.UsagePlan, error) {
	if in.Quota != nil && in.Quota.Period == "" {
		return nil, NewBadRequestException("quota period is required when quota is set")
	}
	if in.Name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if !validateUsagePlanNameLen(in.Name) {
		return nil, NewBadRequestException("name must be between 1 and 255 characters")
	}

	plan := &apigateway.UsagePlan{
		Name:        in.Name,
		Description: in.Description,
		Tags:        in.Tags,
	}

	for _, as := range in.ApiStages {
		for _, t := range as.Throttle {
			if !validateThrottleBurstLimit(t.BurstLimit) {
				return nil, NewBadRequestException("per-stage throttle burstLimit must be between 0 and 10000")
			}
			if !validateThrottleRateLimit(t.RateLimit) {
				return nil, NewBadRequestException("per-stage throttle rateLimit must be between 0 and 10000")
			}
		}
		stage := apigateway.ApiStage{
			ApiId:    as.ApiId,
			Stage:    as.Stage,
			Throttle: as.Throttle,
		}
		plan.ApiStages = append(plan.ApiStages, stage)
	}

	if in.Quota != nil {
		if !validateQuotaPeriod(in.Quota.Period) {
			return nil, NewBadRequestException("Invalid quota period: must be DAY, WEEK, or MONTH")
		}
		plan.Quota = &apigateway.Quota{
			Limit:  in.Quota.Limit,
			Offset: in.Quota.Offset,
			Period: in.Quota.Period,
		}
	}

	if in.Throttle != nil {
		if !validateThrottleBurstLimit(in.Throttle.BurstLimit) {
			return nil, NewBadRequestException("throttle burstLimit must be between 0 and 10000")
		}
		if !validateThrottleRateLimit(in.Throttle.RateLimit) {
			return nil, NewBadRequestException("throttle rateLimit must be between 0 and 10000")
		}
		plan.Throttle = &apigateway.Throttle{
			BurstLimit: in.Throttle.BurstLimit,
			RateLimit:  in.Throttle.RateLimit,
		}
	}

	created, err := stores.usage.CreateUsagePlan(plan)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getUsagePlanCore retrieves a single usage plan.
func (s *APIGatewayService) getUsagePlanCore(stores *apiGatewayStores, usagePlanId string) (*apigateway.UsagePlan, error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}
	plan, err := stores.usage.GetUsagePlan(usagePlanId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return plan, nil
}

// deleteUsagePlanCore removes a usage plan.
func (s *APIGatewayService) deleteUsagePlanCore(stores *apiGatewayStores, usagePlanId string) error {
	if usagePlanId == "" {
		return NewBadRequestException("usagePlanId is required")
	}
	if err := stores.usage.DeleteUsagePlan(usagePlanId); err != nil {
		return toApiGatewayError(err)
	}
	return nil
}

// listUsagePlansCore returns a page of usage plans.
func (s *APIGatewayService) listUsagePlansCore(stores *apiGatewayStores, limit int, marker string) (*common.ListResult[apigateway.UsagePlan], error) {
	resolved, err := resolvePageLimit(limit)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return stores.usage.ListUsagePlans(common.ListOptions{
		Marker:   marker,
		MaxItems: resolved,
	})
}

// updateUsagePlanCore applies patch operations to a usage plan under the
// per-plan key locker.
func (s *APIGatewayService) updateUsagePlanCore(
	stores *apiGatewayStores,
	usagePlanId string,
	patches []PatchOperation,
) (*apigateway.UsagePlan, error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}

	stores.keyLocker.Lock(usagePlanId)
	defer stores.keyLocker.Unlock(usagePlanId)

	usagePlan, err := stores.usage.GetUsagePlan(usagePlanId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range patches {
		handled := false
		switch {
		case po.Path == "/name":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			usagePlan.Name = po.Value
		case po.Path == "/description":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			usagePlan.Description = po.Value
		case po.Path == "/productCode":
			handled = true
			// The row documents add, replace and remove.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			usagePlan.ProductCode = po.Value
		case po.Path == "/quota":
			handled = true
			// The whole-member row documents remove only.
			if err := requirePatchOp(po, opRemove); err != nil {
				return nil, err
			}
			usagePlan.Quota = nil
		case po.Path == "/throttle":
			handled = true
			// The whole-member row documents remove only.
			if err := requirePatchOp(po, opRemove); err != nil {
				return nil, err
			}
			usagePlan.Throttle = nil
		case po.Path == "/quota/limit":
			handled = true
			// The row documents add and replace.
			if err := requirePatchOp(po, opAdd|opReplace); err != nil {
				return nil, err
			}
			if usagePlan.Quota == nil {
				usagePlan.Quota = &apigateway.Quota{}
			}
			if v, err := parseInt64(po.Value); err != nil {
				return nil, NewBadRequestException("invalid quota limit: not a number")
			} else {
				usagePlan.Quota.Limit = v
			}
		case po.Path == "/quota/period":
			handled = true
			// The row documents add and replace.
			if err := requirePatchOp(po, opAdd|opReplace); err != nil {
				return nil, err
			}
			if usagePlan.Quota == nil {
				usagePlan.Quota = &apigateway.Quota{}
			}
			if !validateQuotaPeriod(po.Value) {
				return nil, NewBadRequestException("Invalid quota period: must be DAY, WEEK, or MONTH")
			}
			usagePlan.Quota.Period = po.Value
		case po.Path == "/quota/offset":
			handled = true
			// The row documents add and replace.
			if err := requirePatchOp(po, opAdd|opReplace); err != nil {
				return nil, err
			}
			if usagePlan.Quota == nil {
				usagePlan.Quota = &apigateway.Quota{}
			}
			if v, err := parseInt64(po.Value); err != nil {
				return nil, NewBadRequestException("invalid quota offset: not a number")
			} else {
				usagePlan.Quota.Offset = v
			}
		case po.Path == "/throttle/burstLimit":
			handled = true
			// The row documents add and replace.
			if err := requirePatchOp(po, opAdd|opReplace); err != nil {
				return nil, err
			}
			if usagePlan.Throttle == nil {
				usagePlan.Throttle = &apigateway.Throttle{}
			}
			if v, err := parseInt64(po.Value); err != nil {
				return nil, NewBadRequestException("invalid throttle burstLimit: not a number")
			} else if !validateThrottleBurstLimit(v) {
				return nil, NewBadRequestException("throttle burstLimit must be between 0 and 10000")
			} else {
				usagePlan.Throttle.BurstLimit = v
			}
		case po.Path == "/throttle/rateLimit":
			handled = true
			// The row documents add and replace.
			if err := requirePatchOp(po, opAdd|opReplace); err != nil {
				return nil, err
			}
			if usagePlan.Throttle == nil {
				usagePlan.Throttle = &apigateway.Throttle{}
			}
			if v, err := parseFloat64(po.Value); err != nil {
				return nil, NewBadRequestException("invalid throttle rateLimit: not a number")
			} else if !validateThrottleRateLimit(v) {
				return nil, NewBadRequestException("throttle rateLimit must be between 0 and 10000")
			} else {
				usagePlan.Throttle.RateLimit = v
			}
		case po.Path == "/apiStages":
			handled = true
			// The /apiStages row documents add and remove; the developer
			// guide example carries the value as apiId:stageName.
			if err := requirePatchOp(po, opAdd|opRemove); err != nil {
				return nil, err
			}
			apiId, stageName, ok := splitApiStageValue(po.Value)
			if !ok {
				return nil, NewBadRequestException("invalid apiStages value, expected apiId:stageName: " + po.Value)
			}
			if po.Op == "add" {
				if findApiStage(usagePlan.ApiStages, apiId, stageName) < 0 {
					usagePlan.ApiStages = append(usagePlan.ApiStages, apigateway.ApiStage{
						ApiId: apiId,
						Stage: stageName,
					})
				}
			} else {
				usagePlan.ApiStages = removeApiStage(usagePlan.ApiStages, apiId, stageName)
			}
		case strings.HasPrefix(po.Path, "/apiStages/"):
			handled = true
			// The documented element addressing is
			// /apiStages/{apiId}:{stageName}/throttle/... — the numeric
			// index forms appear nowhere in the official tables.
			rest := strings.TrimPrefix(po.Path, "/apiStages/")
			addr := rest
			if idx := strings.Index(rest, "/"); idx >= 0 {
				addr = rest[:idx]
			}
			apiId, stageName, ok := splitApiStageValue(addr)
			if !ok {
				return nil, NewBadRequestException("invalid apiStages address, expected apiId:stageName: " + addr)
			}
			stageIdx := findApiStage(usagePlan.ApiStages, apiId, stageName)
			if stageIdx < 0 {
				return nil, NewBadRequestException("api stage not found in the usage plan: " + addr)
			}
			if err := applyApiStageThrottlePatch(&usagePlan.ApiStages[stageIdx], po); err != nil {
				return nil, err
			}
		}
		if !handled {
			return nil, unknownPatchPathError(po)
		}
	}

	if usagePlan.Quota != nil && !validateQuotaPeriod(usagePlan.Quota.Period) {
		return nil, NewBadRequestException("quota period is required when quota is set")
	}

	if err := stores.usage.UpdateUsagePlan(usagePlan); err != nil {
		return nil, toApiGatewayError(err)
	}
	return usagePlan, nil
}

// splitApiStageValue splits the apiId:stageName identifier the documented
// /apiStages patch addressing uses. Both segments must be non-empty and the
// stage name must be a legal stage name.
func splitApiStageValue(v string) (string, string, bool) {
	apiId, stage, ok := strings.Cut(v, ":")
	if !ok || apiId == "" || !validateStageName(stage) {
		return "", "", false
	}
	return apiId, stage, true
}

// findApiStage returns the index of the apiStages entry for the api and
// stage, or -1 when the plan does not carry it.
func findApiStage(stages []apigateway.ApiStage, apiId, stage string) int {
	return slices.IndexFunc(stages, func(s apigateway.ApiStage) bool {
		return s.ApiId == apiId && s.Stage == stage
	})
}

// removeApiStage drops the apiStages entry for the api and stage.
func removeApiStage(stages []apigateway.ApiStage, apiId, stage string) []apigateway.ApiStage {
	idx := findApiStage(stages, apiId, stage)
	if idx < 0 {
		return stages
	}
	return slices.Delete(stages, idx, idx+1)
}

// applyApiStageThrottlePatch applies the documented per-stage throttle patch
// family to one apiStages entry:
//
//	/apiStages/{apiId}:{stage}/throttle                                  add/replace set the whole throttle map from a JSON object; remove clears it
//	/apiStages/{apiId}:{stage}/throttle/{resourcePath}/{httpMethod}      remove deletes that throttle key
//	/apiStages/{apiId}:{stage}/throttle/{resourcePath}/{httpMethod}/rateLimit|burstLimit  add/replace set the limit
//
// The method key is the as-addressed pointer derivation of methodMapKey: the
// resource path token arrives JSON Pointer escaped and stays as addressed
// ("/throttle/~1pets/GET" keys as "~1pets/GET"; the root's empty token keys
// as "//GET"), matching the key convention the official CLI update-stage
// example output shows for the stage methodSettings map that documents the
// same {resource_path}/{http_method} idiom. The whole-throttle JSON value's
// entry keys are data in the request and pass through unchanged.
func applyApiStageThrottlePatch(stage *apigateway.ApiStage, po PatchOperation) error {
	rest := strings.TrimPrefix(po.Path, "/apiStages/")
	if idx := strings.Index(rest, "/"); idx >= 0 {
		rest = rest[idx+1:]
	}
	if rest != "throttle" && !strings.HasPrefix(rest, "throttle/") {
		return unknownPatchPathError(po)
	}
	body := strings.TrimPrefix(rest, "throttle")

	if body == "" || body == "/" {
		// The whole-throttle row documents add, replace and remove.
		if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
			return err
		}
		if po.Op == "remove" {
			stage.Throttle = nil
			return nil
		}
		parsed, err := parseApiStageThrottleValue(po)
		if err != nil {
			return err
		}
		stage.Throttle = parsed
		return nil
	}

	// body addresses one method: "/{resourcePath}/{httpMethod}" with an
	// optional trailing rateLimit or burstLimit member.
	field := ""
	if suffix, ok := strings.CutSuffix(body, "/rateLimit"); ok {
		field, body = "rateLimit", suffix
	} else if suffix, ok := strings.CutSuffix(body, "/burstLimit"); ok {
		field, body = "burstLimit", suffix
	}
	tokens := splitPatchTokens(body)
	if len(tokens) < 2 {
		return unknownPatchPathError(po)
	}
	key := methodMapKey(tokens[:len(tokens)-1], tokens[len(tokens)-1])
	if field == "" {
		// The method row without a trailing member documents remove only.
		if err := requirePatchOp(po, opRemove); err != nil {
			return err
		}
		delete(stage.Throttle, key)
		return nil
	}
	if err := requirePatchOp(po, opAdd|opReplace); err != nil {
		return err
	}
	if stage.Throttle == nil {
		stage.Throttle = make(map[string]*apigateway.Throttle)
	}
	t := stage.Throttle[key]
	if t == nil {
		t = &apigateway.Throttle{}
		stage.Throttle[key] = t
	}
	if field == "rateLimit" {
		v, err := parseFloat64(po.Value)
		if err != nil {
			return NewBadRequestException("invalid throttle rateLimit: not a number")
		}
		if !validateThrottleRateLimit(v) {
			return NewBadRequestException("throttle rateLimit must be between 0 and 10000")
		}
		t.RateLimit = v
	} else {
		v, err := parseInt64(po.Value)
		if err != nil {
			return NewBadRequestException("invalid throttle burstLimit: not a number")
		}
		if !validateThrottleBurstLimit(v) {
			return NewBadRequestException("throttle burstLimit must be between 0 and 10000")
		}
		t.BurstLimit = v
	}
	return nil
}

// throttleSettingsPatchValue mirrors the JSON form of one throttle entry
// inside a whole-throttle replace value, with the AWS member names.
type throttleSettingsPatchValue struct {
	RateLimit  *float64 `json:"rateLimit"`
	BurstLimit *int64   `json:"burstLimit"`
}

// parseApiStageThrottleValue decodes the whole-throttle patch value: a JSON
// object keyed by "{resourcePath}/{httpMethod}" whose entries carry the AWS
// member names. The value shape follows the PatchOperation value
// documentation's JSON-object rule; no official example pins the entry-key
// form, so the keys pass through as data. Every entry runs through the same
// validators as the plan-level throttle.
func parseApiStageThrottleValue(po PatchOperation) (map[string]*apigateway.Throttle, error) {
	var raw map[string]throttleSettingsPatchValue
	if err := json.Unmarshal([]byte(po.Value), &raw); err != nil {
		return nil, NewBadRequestException(fmt.Sprintf(
			"Invalid patch value for '%s': expected a JSON object of method keys to throttle settings", po.Path))
	}
	parsed := make(map[string]*apigateway.Throttle, len(raw))
	for key, entry := range raw {
		if key == "" {
			return nil, NewBadRequestException(fmt.Sprintf(
				"Invalid patch value for '%s': entry keys must not be empty", po.Path))
		}
		t := &apigateway.Throttle{}
		if entry.RateLimit != nil {
			if !validateThrottleRateLimit(*entry.RateLimit) {
				return nil, NewBadRequestException("throttle rateLimit must be between 0 and 10000")
			}
			t.RateLimit = *entry.RateLimit
		}
		if entry.BurstLimit != nil {
			if !validateThrottleBurstLimit(*entry.BurstLimit) {
				return nil, NewBadRequestException("throttle burstLimit must be between 0 and 10000")
			}
			t.BurstLimit = *entry.BurstLimit
		}
		parsed[key] = t
	}
	return parsed, nil
}

// createUsagePlanKeyCore attaches an API key to a usage plan. The API key
// must already exist; the key's value and name are copied into the new
// usage plan key entry.
func (s *APIGatewayService) createUsagePlanKeyCore(
	stores *apiGatewayStores,
	usagePlanId string,
	in *UsagePlanKeyInput,
) (*apigateway.UsagePlanKey, error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}
	if in.KeyId == "" {
		return nil, NewBadRequestException("keyId is required")
	}
	if in.KeyType == "" {
		return nil, NewBadRequestException("keyType is required")
	}
	if in.KeyType != "API_KEY" {
		return nil, NewBadRequestException("keyType must be API_KEY")
	}

	apiKey, err := stores.usage.GetApiKey(in.KeyId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	key := &apigateway.UsagePlanKey{
		Id:    in.KeyId,
		Type:  in.KeyType,
		Value: apiKey.Value,
		Name:  apiKey.Name,
	}

	created, err := stores.usage.CreateUsagePlanKey(usagePlanId, key)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getUsagePlanKeyCore retrieves a single usage plan key.
func (s *APIGatewayService) getUsagePlanKeyCore(stores *apiGatewayStores, usagePlanId, keyId string) (*apigateway.UsagePlanKey, error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}
	if keyId == "" {
		return nil, NewBadRequestException("keyId is required")
	}
	key, err := stores.usage.GetUsagePlanKey(usagePlanId, keyId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return key, nil
}

// deleteUsagePlanKeyCore removes an API key association from a usage plan.
func (s *APIGatewayService) deleteUsagePlanKeyCore(stores *apiGatewayStores, usagePlanId, keyId string) error {
	if usagePlanId == "" {
		return NewBadRequestException("usagePlanId is required")
	}
	if keyId == "" {
		return NewBadRequestException("keyId is required")
	}
	if err := stores.usage.DeleteUsagePlanKey(usagePlanId, keyId); err != nil {
		return toApiGatewayError(err)
	}
	return nil
}

// listUsagePlanKeysCore returns a page of usage plan keys.
func (s *APIGatewayService) listUsagePlanKeysCore(stores *apiGatewayStores, usagePlanId string, limit int, marker string) (*common.ListResult[apigateway.UsagePlanKey], error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}
	resolved, err := resolvePageLimit(limit)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return stores.usage.ListUsagePlanKeys(usagePlanId, common.ListOptions{
		Marker:   marker,
		MaxItems: resolved,
	})
}

// getUsageCore returns aggregated usage data for a usage plan over a date
// range. The caller pre-validates the date format and the 90-day range
// because that is transport-specific.
func (s *APIGatewayService) getUsageCore(
	stores *apiGatewayStores,
	usagePlanId, keyId, startDate, endDate string,
) (map[string]interface{}, error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}
	if startDate == "" || endDate == "" {
		return nil, NewBadRequestException("startDate and endDate are required")
	}
	if !validateUsageDateFormat(startDate) || !validateUsageDateFormat(endDate) {
		return nil, NewBadRequestException("startDate and endDate must be in YYYY-MM-DD format")
	}

	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)
	if startTime.After(endTime) {
		return nil, NewBadRequestException("startDate must not be after endDate")
	}
	if endTime.Sub(startTime).Hours()/24 > maxUsageDateRangeDays {
		return nil, NewBadRequestException("The date range must not exceed 90 days")
	}

	if _, err := stores.usage.GetUsagePlan(usagePlanId); err != nil {
		return nil, toApiGatewayError(err)
	}

	var apiKeys []string
	if keyId != "" {
		apiKeys = []string{keyId}
	} else {
		// Page through the store's own listing API rather than reaching into
		// its key space: the key schema stays an implementation detail of
		// the UsageStore.
		marker := ""
		for {
			page, err := stores.usage.ListUsagePlanKeys(usagePlanId, common.ListOptions{Marker: marker})
			if err != nil {
				return nil, toApiGatewayError(err)
			}
			for _, key := range page.Items {
				apiKeys = append(apiKeys, key.Id)
			}
			if !page.IsTruncated {
				break
			}
			marker = page.NextMarker
		}
	}

	usageCounts := make(map[string]int64)
	for _, keyId := range apiKeys {
		records, err := stores.usage.ListUsageRecordsForAPIKey(usagePlanId, keyId, startDate, endDate)
		if err != nil {
			logs.Warn("GetUsage: failed to list usage records for key", logs.String("keyId", keyId), logs.Err(err))
			continue
		}
		for _, record := range records {
			usageCounts[record.Date] += record.RequestCount
		}
	}

	items := make([]interface{}, 0, len(usageCounts))
	for date, count := range usageCounts {
		items = append(items, map[string]interface{}{
			"date":         date,
			"requestCount": count,
		})
	}

	return map[string]interface{}{
		"usagePlanId": usagePlanId,
		"startDate":   startDate,
		"endDate":     endDate,
		"items":       items,
	}, nil
}
