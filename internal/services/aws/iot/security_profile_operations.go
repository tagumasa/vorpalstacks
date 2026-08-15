package iot

import (
	"context"
	"encoding/json"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "securityProfileName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	behaviors, err := parseDualForm(req.Parameters, "behaviors", parseBehaviors, parseBehaviorsParam)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	alertTargets, err := parseDualForm(req.Parameters, "alertTargets", parseAlertTargets, parseAlertTargetsParam)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	metrics := parseDualFormNoError(req.Parameters, "additionalMetricsToRetainV2", parseStringList, parseMetricsParam)
	additionalMetricsV1 := request.GetStringList(req.Parameters, "additionalMetricsToRetain")
	metricsExportConfig := request.GetParamCaseInsensitive(req.Parameters, "metricsExportConfig")
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	spTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		spTags[t.Key] = t.Value
	}

	sp := &iotstore.SecurityProfile{
		SecurityProfileName:         name,
		SecurityProfileDescription:  request.GetParamCaseInsensitive(req.Parameters, "securityProfileDescription"),
		Behaviors:                   behaviors,
		AlertTargets:                alertTargets,
		AdditionalMetricsToRetainV2: metrics,
		AdditionalMetricsToRetain:   additionalMetricsV1,
		MetricsExportConfig:         metricsExportConfig,
		Tags:                        spTags,
		Version:                     1,
		CreationDate:                time.Now().UTC(),
		LastModifiedDate:            time.Now().UTC(),
	}

	created, err := store.CreateSecurityProfile(sp)
	if err != nil {
		return nil, err
	}

	return securityProfileToResponse(created), nil
}

func (s *IoTService) DescribeSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "securityProfileName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sp, err := store.GetSecurityProfile(name)
	if err != nil {
		return nil, err
	}

	return securityProfileToResponse(sp), nil
}

func (s *IoTService) UpdateSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "securityProfileName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.GetSecurityProfile(name)
	if err != nil {
		return nil, err
	}

	if expVer := request.GetIntParam(req.Parameters, "expectedVersion"); expVer > 0 && int64(expVer) != existing.Version {
		return nil, iotstore.ErrVersionConflict
	}

	desc := request.GetParamCaseInsensitive(req.Parameters, "securityProfileDescription")
	if desc != "" {
		existing.SecurityProfileDescription = desc
	}

	if behaviors, err := parseDualForm(req.Parameters, "behaviors", parseBehaviors, parseBehaviorsParam); err != nil {
		return nil, iotstore.ErrInvalidRequest
	} else if behaviors != nil {
		existing.Behaviors = behaviors
	}

	if alertTargets, err := parseDualForm(req.Parameters, "alertTargets", parseAlertTargets, parseAlertTargetsParam); err != nil {
		return nil, iotstore.ErrInvalidRequest
	} else if alertTargets != nil {
		existing.AlertTargets = alertTargets
	}

	if metrics := parseDualFormNoError(req.Parameters, "additionalMetricsToRetainV2", parseStringList, parseMetricsParam); metrics != nil {
		existing.AdditionalMetricsToRetainV2 = metrics
	}
	if metricsV1 := request.GetStringList(req.Parameters, "additionalMetricsToRetain"); metricsV1 != nil {
		existing.AdditionalMetricsToRetain = metricsV1
	}
	if mec := request.GetParamCaseInsensitive(req.Parameters, "metricsExportConfig"); mec != "" {
		existing.MetricsExportConfig = mec
	}

	existing.Version++
	existing.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateSecurityProfile(name, existing); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DeleteSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "securityProfileName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	arn := iotstore.BuildSecurityProfileARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteSecurityProfile(name); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListSecurityProfiles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	profiles, err := store.ListSecurityProfiles(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(profiles.Items))
	for _, sp := range profiles.Items {
		items = append(items, securityProfileIdentifierResponse(sp.SecurityProfileName, sp.SecurityProfileARN))
	}

	return listResponse("securityProfileIdentifiers", items, profiles.NextMarker), nil
}

func (s *IoTService) ValidateSecurityProfileBehaviors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	bRaw := req.Parameters["behaviors"]

	var behaviors []*iotstore.Behavior
	var err error

	switch v := bRaw.(type) {
	case string:
		if v == "" {
			return map[string]interface{}{"valid": true}, nil
		}
		behaviors, err = parseBehaviors(v)
	case []interface{}:
		behaviors, err = parseBehaviorsFromList(v)
	default:
		return map[string]interface{}{"valid": true}, nil
	}

	if err != nil {
		return map[string]interface{}{
			"valid": false,
			"validationMessages": []map[string]interface{}{
				{"message": "failed to parse behaviors JSON: " + err.Error()},
			},
		}, nil
	}

	valid := true
	var msgs []map[string]interface{}

	validOperators := map[string]bool{
		"less-than": true, "less-than-equals": true,
		"greater-than": true, "greater-than-equals": true,
		"not-equal": true, "equal": true,
	}

	for _, b := range behaviors {
		if b.Name == "" {
			valid = false
			msgs = append(msgs, map[string]interface{}{
				"message": "behavior name must not be empty",
			})
			continue
		}
		if b.Criteria == nil {
			valid = false
			msgs = append(msgs, map[string]interface{}{
				"message": "behavior '" + b.Name + "' must have criteria",
			})
			continue
		}
		if !validOperators[b.Criteria.ComparisonOperator] {
			valid = false
			msgs = append(msgs, map[string]interface{}{
				"message": "behavior '" + b.Name + "' has invalid comparisonOperator '" + b.Criteria.ComparisonOperator + "'",
			})
		}
		if b.Criteria.DurationSeconds < 0 {
			valid = false
			msgs = append(msgs, map[string]interface{}{
				"message": "behavior '" + b.Name + "' has negative durationSeconds",
			})
		}
	}

	result := map[string]interface{}{"valid": valid}
	if len(msgs) > 0 {
		result["validationMessages"] = msgs
	}
	return result, nil
}

func parseDualForm[T any](params map[string]interface{}, key string, fromJSON func(string) (T, error), fromParam func(interface{}) (T, error)) (T, error) {
	var zero T
	if strVal := request.GetParamCaseInsensitive(params, key); strVal != "" {
		return fromJSON(strVal)
	}
	if rawVal := params[key]; rawVal != nil {
		return fromParam(rawVal)
	}
	return zero, nil
}

func parseDualFormNoError[T any](params map[string]interface{}, key string, fromJSON func(string) T, fromParam func(interface{}) T) T {
	if strVal := request.GetParamCaseInsensitive(params, key); strVal != "" {
		return fromJSON(strVal)
	}
	if rawVal := params[key]; rawVal != nil {
		return fromParam(rawVal)
	}
	var zero T
	return zero
}

func parseBehaviors(jsonStr string) ([]*iotstore.Behavior, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var raw []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return nil, err
	}
	return rawToBehaviors(raw)
}

func parseBehaviorsFromList(list []interface{}) ([]*iotstore.Behavior, error) {
	if len(list) == 0 {
		return nil, nil
	}
	var raw []map[string]interface{}
	for _, item := range list {
		if m, ok := item.(map[string]interface{}); ok {
			raw = append(raw, m)
		}
	}
	return rawToBehaviors(raw)
}

func rawToBehaviors(raw []map[string]interface{}) ([]*iotstore.Behavior, error) {
	behaviors := make([]*iotstore.Behavior, 0, len(raw))
	for _, r := range raw {
		b := &iotstore.Behavior{
			Name:            strVal(r["name"]),
			Metric:          strVal(r["metric"]),
			MetricDimension: strVal(r["metricDimension"]),
			SuppressAlerts:  boolVal(r["suppressAlerts"]),
			ExportMetric:    boolVal(r["exportMetric"]),
		}
		if c, ok := r["criteria"].(map[string]interface{}); ok {
			b.Criteria = &iotstore.BehaviorCriteria{
				ComparisonOperator:           strVal(c["comparisonOperator"]),
				Value:                        float64Val(c["value"]),
				DurationSeconds:              int64Val(c["durationSeconds"]),
				ConsecutiveDatapointsToAlarm: int64Val(c["consecutiveDatapointsToAlarm"]),
				ConsecutiveDatapointsToClear: int64Val(c["consecutiveDatapointsToClear"]),
			}
		}
		behaviors = append(behaviors, b)
	}
	return behaviors, nil
}

func parseBehaviorsParam(v interface{}) ([]*iotstore.Behavior, error) {
	switch val := v.(type) {
	case string:
		return parseBehaviors(val)
	case []interface{}:
		return parseBehaviorsFromList(val)
	default:
		return nil, nil
	}
}

func parseAlertTargets(jsonStr string) (map[string]*iotstore.AlertTarget, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return nil, err
	}
	return rawToAlertTargets(raw), nil
}

func rawToAlertTargets(raw map[string]interface{}) map[string]*iotstore.AlertTarget {
	targets := make(map[string]*iotstore.AlertTarget)
	for k, v := range raw {
		if m, ok := v.(map[string]interface{}); ok {
			targets[k] = &iotstore.AlertTarget{
				AlertTargetARN: strVal(m["alertTargetArn"]),
				RoleARN:        strVal(m["roleArn"]),
			}
		}
	}
	return targets
}

func parseAlertTargetsParam(v interface{}) (map[string]*iotstore.AlertTarget, error) {
	switch val := v.(type) {
	case string:
		return parseAlertTargets(val)
	case map[string]interface{}:
		return rawToAlertTargets(val), nil
	default:
		return nil, nil
	}
}

func parseStringList(jsonStr string) []string {
	if jsonStr == "" {
		return nil
	}
	var list []string
	if json.Unmarshal([]byte(jsonStr), &list) == nil {
		return list
	}
	return nil
}

func parseMetricsParam(v interface{}) []string {
	switch val := v.(type) {
	case string:
		return parseStringList(val)
	case []interface{}:
		var list []string
		for _, item := range val {
			if s, ok := item.(string); ok {
				list = append(list, s)
			}
		}
		return list
	default:
		return nil
	}
}

func securityProfileToResponse(sp *iotstore.SecurityProfile) map[string]interface{} {
	var behaviors interface{} = []interface{}{}
	if sp.Behaviors != nil {
		behaviors = sp.Behaviors
	}

	alertTargets := map[string]interface{}{}
	if sp.AlertTargets != nil {
		for k, v := range sp.AlertTargets {
			alertTargets[k] = map[string]interface{}{
				"alertTargetArn": v.AlertTargetARN,
				"roleArn":        v.RoleARN,
			}
		}
	}

	metrics := sp.AdditionalMetricsToRetainV2
	if metrics == nil {
		metrics = []string{}
	}

	metricsV1 := sp.AdditionalMetricsToRetain
	if metricsV1 == nil {
		metricsV1 = []string{}
	}

	resp := map[string]interface{}{
		"securityProfileName":         sp.SecurityProfileName,
		"securityProfileArn":          sp.SecurityProfileARN,
		"securityProfileDescription":  sp.SecurityProfileDescription,
		"behaviors":                   behaviors,
		"alertTargets":                alertTargets,
		"additionalMetricsToRetain":   metricsV1,
		"additionalMetricsToRetainV2": metrics,
		"version":                     sp.Version,
		"creationDate":                sp.CreationDate.Unix(),
		"lastModifiedDate":            sp.LastModifiedDate.Unix(),
	}
	if sp.MetricsExportConfig != "" {
		var mec interface{}
		if err := json.Unmarshal([]byte(sp.MetricsExportConfig), &mec); err == nil {
			resp["metricsExportConfig"] = mec
		} else {
			resp["metricsExportConfig"] = sp.MetricsExportConfig
		}
	}
	return resp
}

// ---- Security Profile attach ---------------------------------------
// AWS persists profile<->target associations so that ListSecurityProfilesForTarget
// and ListTargetsForSecurityProfile return real data. Attach/Detach enforce
// ResourceNotFoundException when the association does not exist (Detach) and
// return empty responses per the Smithy output shapes.

func (s *IoTService) AttachSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	profileName := request.GetParamCaseInsensitive(req.Parameters, "securityProfileName")
	targetArn := request.GetParamCaseInsensitive(req.Parameters, "securityProfileTargetArn")
	if profileName == "" || targetArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Forward mapping: profile -> target. Reverse mapping: target -> profile.
	// Both are stored so that ListSecurityProfilesForTarget (target->profile)
	// and ListTargetsForSecurityProfile (profile->target) can scan a single
	// prefix.
	forwardKey := "secProfileTarget/" + profileName + "/" + targetArn
	reverseKey := "secTargetProfile/" + targetArn + "/" + profileName
	assocValue := map[string]interface{}{
		"securityProfileName":      profileName,
		"securityProfileTargetArn": targetArn,
	}
	if err := store.PutGeneric(forwardKey, assocValue); err != nil {
		return nil, err
	}
	if err := store.PutGeneric(reverseKey, assocValue); err != nil {
		// Rollback forward write to maintain bidirectional consistency.
		_ = store.DeleteGeneric(forwardKey)
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DetachSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	profileName := request.GetParamCaseInsensitive(req.Parameters, "securityProfileName")
	targetArn := request.GetParamCaseInsensitive(req.Parameters, "securityProfileTargetArn")
	if profileName == "" || targetArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	forwardKey := "secProfileTarget/" + profileName + "/" + targetArn
	reverseKey := "secTargetProfile/" + targetArn + "/" + profileName
	exists, err := store.GetGenericExists(forwardKey, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrSecurityProfileAttachmentNotFound
	}
	// Attempt both deletes so a partial failure does not leave stale mappings
	// that block subsequent retries (the existence check above would reject
	// a retry after a partial delete).
	errForward := store.DeleteGeneric(forwardKey)
	errReverse := store.DeleteGeneric(reverseKey)
	if errForward != nil {
		return nil, errForward
	}
	if errReverse != nil {
		return nil, errReverse
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) ListSecurityProfilesForTarget(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	targetArn := request.GetParamCaseInsensitive(req.Parameters, "securityProfileTargetArn")
	if targetArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("secTargetProfile/" + targetArn + "/")
	if err != nil {
		return nil, err
	}
	mappings := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		profileName, _ := rec["securityProfileName"].(string)
		mappings = append(mappings, map[string]interface{}{
			"securityProfileIdentifier": map[string]interface{}{
				"name": profileName,
			},
			"target": map[string]interface{}{
				"arn": targetArn,
			},
		})
	}
	return paginatedMaps("securityProfileTargetMappings", mappings, req.Parameters), nil
}
func (s *IoTService) ListTargetsForSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	profileName := request.GetParamCaseInsensitive(req.Parameters, "securityProfileName")
	if profileName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("secProfileTarget/" + profileName + "/")
	if err != nil {
		return nil, err
	}
	targets := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		targetArn, _ := rec["securityProfileTargetArn"].(string)
		targets = append(targets, map[string]interface{}{
			"arn": targetArn,
		})
	}
	return paginatedMaps("securityProfileTargets", targets, req.Parameters), nil
}
func (s *IoTService) PutVerificationStateOnViolation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	violationId := request.GetParamCaseInsensitive(req.Parameters, "violationId")
	verificationState := request.GetParamCaseInsensitive(req.Parameters, "verificationState")
	if violationId == "" || verificationState == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "violation/" + violationId
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	// No Device Defender engine generates violations, so the record usually
	// does not exist. AWS lists only InvalidRequestException in the Smithy
	// errors trait (not ResourceNotFoundException), so return InvalidRequest
	// for an unknown violation id rather than 404.
	if !exists {
		return nil, iotstore.ErrInvalidRequest
	}
	rec["verificationState"] = verificationState
	if desc := request.GetParamCaseInsensitive(req.Parameters, "verificationStateDescription"); desc != "" {
		rec["verificationStateDescription"] = desc
	}
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
