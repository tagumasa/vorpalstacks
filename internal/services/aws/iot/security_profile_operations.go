package iot

import (
	"context"
	"encoding/json"
	"time"

	"vorpalstacks/internal/common/request"
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

	bStr := request.GetParamCaseInsensitive(req.Parameters, "behaviors")
	var behaviors []*iotstore.Behavior
	if bStr != "" {
		behaviors, err = parseBehaviors(bStr)
	} else if bRaw := req.Parameters["behaviors"]; bRaw != nil {
		behaviors, err = parseBehaviorsParam(bRaw)
	}
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	aStr := request.GetParamCaseInsensitive(req.Parameters, "alertTargets")
	var alertTargets map[string]*iotstore.AlertTarget
	if aStr != "" {
		alertTargets, err = parseAlertTargets(aStr)
	} else if aRaw := req.Parameters["alertTargets"]; aRaw != nil {
		alertTargets, err = parseAlertTargetsParam(aRaw)
	}
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	mStr := request.GetParamCaseInsensitive(req.Parameters, "additionalMetricsToRetainVersionTwo")
	var metrics []string
	if mStr != "" {
		metrics = parseStringList(mStr)
	} else if mRaw := req.Parameters["additionalMetricsToRetainVersionTwo"]; mRaw != nil {
		metrics = parseMetricsParam(mRaw)
	}

	sp := &iotstore.SecurityProfile{
		SecurityProfileName:         name,
		SecurityProfileDescription:  request.GetParamCaseInsensitive(req.Parameters, "securityProfileDescription"),
		Behaviors:                   behaviors,
		AlertTargets:                alertTargets,
		AdditionalMetricsToRetainV2: metrics,
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
		return nil, iotstore.ErrSecurityProfileNotFound
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
		return nil, iotstore.ErrSecurityProfileNotFound
	}

	if expVer := request.GetIntParam(req.Parameters, "expectedVersion"); expVer > 0 && int64(expVer) != existing.Version {
		return nil, iotstore.ErrVersionConflict
	}

	desc := request.GetParamCaseInsensitive(req.Parameters, "securityProfileDescription")
	if desc != "" {
		existing.SecurityProfileDescription = desc
	}

	if bStr := request.GetParamCaseInsensitive(req.Parameters, "behaviors"); bStr != "" {
		behaviors, err := parseBehaviors(bStr)
		if err != nil {
			return nil, iotstore.ErrInvalidRequest
		}
		existing.Behaviors = behaviors
	} else if bRaw := req.Parameters["behaviors"]; bRaw != nil {
		behaviors, err := parseBehaviorsParam(bRaw)
		if err != nil {
			return nil, iotstore.ErrInvalidRequest
		}
		existing.Behaviors = behaviors
	}

	if aStr := request.GetParamCaseInsensitive(req.Parameters, "alertTargets"); aStr != "" {
		alertTargets, err := parseAlertTargets(aStr)
		if err != nil {
			return nil, iotstore.ErrInvalidRequest
		}
		existing.AlertTargets = alertTargets
	} else if aRaw := req.Parameters["alertTargets"]; aRaw != nil {
		alertTargets, err := parseAlertTargetsParam(aRaw)
		if err != nil {
			return nil, iotstore.ErrInvalidRequest
		}
		existing.AlertTargets = alertTargets
	}

	if mStr := request.GetParamCaseInsensitive(req.Parameters, "additionalMetricsToRetain"); mStr != "" {
		existing.AdditionalMetricsToRetainV2 = parseStringList(mStr)
	} else if mRaw := req.Parameters["additionalMetricsToRetain"]; mRaw != nil {
		existing.AdditionalMetricsToRetainV2 = parseMetricsParam(mRaw)
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
		items = append(items, securityProfileResponse(sp))
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

	return map[string]interface{}{
		"securityProfileName":                 sp.SecurityProfileName,
		"securityProfileArn":                  sp.SecurityProfileARN,
		"securityProfileDescription":          sp.SecurityProfileDescription,
		"behaviors":                           behaviors,
		"alertTargets":                        alertTargets,
		"additionalMetricsToRetainVersionTwo": metrics,
		"version":                             sp.Version,
		"creationDate":                        sp.CreationDate.Unix(),
		"lastModifiedDate":                    sp.LastModifiedDate.Unix(),
	}
}


