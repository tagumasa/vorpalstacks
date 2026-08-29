package iot

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	behaviors, behaviorsErr := parseDualForm(req.Parameters, "behaviors", parseBehaviors, parseBehaviorsParam)
	alertTargets, alertTargetsErr := parseDualForm(req.Parameters, "alertTargets", parseAlertTargets, parseAlertTargetsParam)
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	spTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		spTags[t.Key] = t.Value
	}
	in := CreateSecurityProfileInput{
		Name:                        request.GetParamCaseInsensitive(req.Parameters, "securityProfileName"),
		Description:                 request.GetParamCaseInsensitive(req.Parameters, "securityProfileDescription"),
		Behaviors:                   behaviors,
		BehaviorsMalformed:          behaviorsErr != nil,
		AlertTargets:                alertTargets,
		AlertTargetsMalformed:       alertTargetsErr != nil,
		AdditionalMetricsToRetainV2: parseDualFormNoError(req.Parameters, "additionalMetricsToRetainV2", parseMetricsToRetainJSON, parseMetricsToRetainParam),
		AdditionalMetricsToRetain:   request.GetStringList(req.Parameters, "additionalMetricsToRetain"),
		MetricsExportConfig:         request.GetParamCaseInsensitive(req.Parameters, "metricsExportConfig"),
		Tags:                        spTags,
	}
	created, err := s.createSecurityProfileCore(store, in)
	if err != nil {
		return nil, err
	}
	// CreateSecurityProfileResponse carries only the name and ARN.
	return map[string]interface{}{
		"securityProfileName": created.SecurityProfileName,
		"securityProfileArn":  created.SecurityProfileARN,
	}, nil
}

func (s *IoTService) DescribeSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	sp, err := s.describeSecurityProfileCore(store, request.GetParamCaseInsensitive(req.Parameters, "securityProfileName"))
	if err != nil {
		return nil, err
	}
	return securityProfileToResponse(sp), nil
}

func (s *IoTService) UpdateSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	behaviors, behaviorsErr := parseDualForm(req.Parameters, "behaviors", parseBehaviors, parseBehaviorsParam)
	alertTargets, alertTargetsErr := parseDualForm(req.Parameters, "alertTargets", parseAlertTargets, parseAlertTargetsParam)
	in := UpdateSecurityProfileInput{
		Name:                            request.GetParamCaseInsensitive(req.Parameters, "securityProfileName"),
		ExpectedVersion:                 int64(request.GetIntParam(req.Parameters, "expectedVersion")),
		Description:                     request.GetParamCaseInsensitive(req.Parameters, "securityProfileDescription"),
		Behaviors:                       behaviors,
		BehaviorsMalformed:              behaviorsErr != nil,
		DeleteBehaviors:                 request.GetBoolParam(req.Parameters, "deleteBehaviors"),
		AlertTargets:                    alertTargets,
		AlertTargetsMalformed:           alertTargetsErr != nil,
		DeleteAlertTargets:              request.GetBoolParam(req.Parameters, "deleteAlertTargets"),
		AdditionalMetricsToRetainV2:     parseDualFormNoError(req.Parameters, "additionalMetricsToRetainV2", parseMetricsToRetainJSON, parseMetricsToRetainParam),
		AdditionalMetricsToRetain:       request.GetStringList(req.Parameters, "additionalMetricsToRetain"),
		DeleteAdditionalMetricsToRetain: request.GetBoolParam(req.Parameters, "deleteAdditionalMetricsToRetain"),
		MetricsExportConfig:             request.GetParamCaseInsensitive(req.Parameters, "metricsExportConfig"),
		DeleteMetricsExportConfig:       request.GetBoolParam(req.Parameters, "deleteMetricsExportConfig"),
	}
	updated, err := s.updateSecurityProfileCore(store, in)
	if err != nil {
		return nil, err
	}
	return securityProfileToResponse(updated), nil
}

func (s *IoTService) DeleteSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteSecurityProfileCore(store, request.GetParamCaseInsensitive(req.Parameters, "securityProfileName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListSecurityProfiles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, nextMarker, err := s.listSecurityProfilesCore(store, parseListOptions(req.Parameters), ListSecurityProfilesInput{
		DimensionName: request.GetParamCaseInsensitive(req.Parameters, "dimensionName"),
		MetricName:    request.GetParamCaseInsensitive(req.Parameters, "metricName"),
	})
	if err != nil {
		return nil, err
	}
	identifiers := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		identifiers = append(identifiers, securityProfileIdentifierResponse(item.Name, item.Arn))
	}
	return listResponse("securityProfileIdentifiers", identifiers, nextMarker), nil
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
			"validationErrors": []map[string]interface{}{
				{"errorMessage": "failed to parse behaviors JSON: " + err.Error()},
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
				"errorMessage": "behavior name must not be empty",
			})
			continue
		}
		if b.Criteria == nil {
			valid = false
			msgs = append(msgs, map[string]interface{}{
				"errorMessage": "behavior '" + b.Name + "' must have criteria",
			})
			continue
		}
		if !validOperators[b.Criteria.ComparisonOperator] {
			valid = false
			msgs = append(msgs, map[string]interface{}{
				"errorMessage": "behavior '" + b.Name + "' has invalid comparisonOperator '" + b.Criteria.ComparisonOperator + "'",
			})
		}
		if b.Criteria.DurationSeconds < 0 {
			valid = false
			msgs = append(msgs, map[string]interface{}{
				"errorMessage": "behavior '" + b.Name + "' has negative durationSeconds",
			})
		}
	}

	// ValidateSecurityProfileBehaviorsResponse carries valid plus
	// validationErrors, each a ValidationError whose only member is
	// errorMessage.
	result := map[string]interface{}{"valid": valid}
	if len(msgs) > 0 {
		result["validationErrors"] = msgs
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
			Name:           strVal(r["name"]),
			Metric:         strVal(r["metric"]),
			SuppressAlerts: boolVal(r["suppressAlerts"]),
			ExportMetric:   boolVal(r["exportMetric"]),
		}
		// The wire form of metricDimension is the MetricDimension structure
		// {dimensionName, operator}; the operator is the optional IN/NOT_IN
		// enum.
		if dim, ok := r["metricDimension"].(map[string]interface{}); ok {
			b.MetricDimension = strVal(dim["dimensionName"])
			switch operator := strVal(dim["operator"]); operator {
			case "":
			case "IN", "NOT_IN":
				b.MetricDimensionOperator = operator
			default:
				return nil, iotstore.ErrInvalidRequest
			}
		}
		if c, ok := r["criteria"].(map[string]interface{}); ok {
			criteria := &iotstore.BehaviorCriteria{
				ComparisonOperator:           strVal(c["comparisonOperator"]),
				Value:                        float64Val(c["value"]),
				DurationSeconds:              int64Val(c["durationSeconds"]),
				ConsecutiveDatapointsToAlarm: int64Val(c["consecutiveDatapointsToAlarm"]),
				ConsecutiveDatapointsToClear: int64Val(c["consecutiveDatapointsToClear"]),
			}
			if st, ok := c["statisticalThreshold"].(map[string]interface{}); ok {
				criteria.StatisticalThreshold = &iotstore.StatisticalThreshold{
					Statistic: strVal(st["statistic"]),
				}
			}
			if ml, ok := c["mlDetectionConfig"].(map[string]interface{}); ok {
				criteria.MLDetectionConfig = &iotstore.MachineLearningDetectionConfig{
					ConfidenceLevel: strVal(ml["confidenceLevel"]),
				}
			}
			b.Criteria = criteria
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

// parseMetricsToRetainJSON parses the JSON form of the
// additionalMetricsToRetainV2 member: a list of MetricToRetain objects.
func parseMetricsToRetainJSON(jsonStr string) []*iotstore.MetricToRetain {
	if jsonStr == "" {
		return nil
	}
	var list []interface{}
	if json.Unmarshal([]byte(jsonStr), &list) != nil {
		return nil
	}
	return metricsToRetainFromList(list)
}

// parseMetricsToRetainParam parses the structured form of the
// additionalMetricsToRetainV2 member.
func parseMetricsToRetainParam(v interface{}) []*iotstore.MetricToRetain {
	switch val := v.(type) {
	case string:
		return parseMetricsToRetainJSON(val)
	case []interface{}:
		return metricsToRetainFromList(val)
	default:
		return nil
	}
}

func metricsToRetainFromList(list []interface{}) []*iotstore.MetricToRetain {
	entries := make([]*iotstore.MetricToRetain, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		entry := &iotstore.MetricToRetain{
			Metric:       strVal(m["metric"]),
			ExportMetric: boolVal(m["exportMetric"]),
		}
		if dim, ok := m["metricDimension"].(map[string]interface{}); ok {
			entry.MetricDimension = strVal(dim["dimensionName"])
			entry.Operator = strVal(dim["operator"])
		}
		entries = append(entries, entry)
	}
	return entries
}

func securityProfileToResponse(sp *iotstore.SecurityProfile) map[string]interface{} {
	// Behaviours serialise through the shared lowerCamel projection so the
	// typed SDK can parse the members.
	behaviors := make([]interface{}, 0, len(sp.Behaviors))
	for _, b := range sp.Behaviors {
		behaviors = append(behaviors, behaviorResponse(b))
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

	metricsV2 := make([]interface{}, 0, len(sp.AdditionalMetricsToRetainV2))
	for _, m := range sp.AdditionalMetricsToRetainV2 {
		entry := map[string]interface{}{
			"metric":       m.Metric,
			"exportMetric": m.ExportMetric,
		}
		if m.MetricDimension != "" {
			dimension := map[string]interface{}{
				"dimensionName": m.MetricDimension,
			}
			if m.Operator != "" {
				dimension["operator"] = m.Operator
			}
			entry["metricDimension"] = dimension
		}
		metricsV2 = append(metricsV2, entry)
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
		"additionalMetricsToRetainV2": metricsV2,
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.attachSecurityProfileCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "securityProfileName"),
		request.GetParamCaseInsensitive(req.Parameters, "securityProfileTargetArn")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DetachSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.detachSecurityProfileCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "securityProfileName"),
		request.GetParamCaseInsensitive(req.Parameters, "securityProfileTargetArn")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) ListSecurityProfilesForTarget(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	targetArn := request.GetParamCaseInsensitive(req.Parameters, "securityProfileTargetArn")
	names, err := s.listSecurityProfilesForTargetCore(store, targetArn)
	if err != nil {
		return nil, err
	}
	mappings := make([]map[string]interface{}, 0, len(names))
	for _, profileName := range names {
		mappings = append(mappings, map[string]interface{}{
			"securityProfileIdentifier": map[string]interface{}{
				"name": profileName,
			},
			"target": map[string]interface{}{
				"arn": targetArn,
			},
		})
	}
	return paginatedMaps("securityProfileTargetMappings", mappings, req.Parameters)
}
func (s *IoTService) ListTargetsForSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	targets, err := s.listTargetsForSecurityProfileCore(store, request.GetParamCaseInsensitive(req.Parameters, "securityProfileName"))
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(targets))
	for _, targetArn := range targets {
		items = append(items, map[string]interface{}{
			"arn": targetArn,
		})
	}
	return paginatedMaps("securityProfileTargets", items, req.Parameters)
}
func (s *IoTService) PutVerificationStateOnViolation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putVerificationStateOnViolationCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "violationId"),
		request.GetParamCaseInsensitive(req.Parameters, "verificationState"),
		request.GetParamCaseInsensitive(req.Parameters, "verificationStateDescription")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
