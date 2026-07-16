package iotevents

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/iotutil"
	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTEventsService) CreateAlarmModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "alarmModelName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	def, err := parseStructParam(req.Parameters, "alarmModelDefinition")
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	am := &iotstore.AlarmModel{
		AlarmModelName:        name,
		AlarmModelDescription: request.GetParamCaseInsensitive(req.Parameters, "alarmModelDescription"),
		RoleARN:               request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		AlarmModelDefinition:  def,
		Severity:              request.GetParamCaseInsensitive(req.Parameters, "severity"),
		Tags:                  parseAlarmTagsParam(req.Parameters),
	}

	created, err := store.CreateAlarmModel(am)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"alarmModelName":        created.AlarmModelName,
		"alarmModelArn":         created.AlarmModelARN,
		"alarmModelDescription": created.AlarmModelDescription,
		"roleArn":               created.RoleARN,
		"alarmModelVersion":     created.AlarmModelVersion,
		"status":                created.Status,
		"creationTime":          created.CreationDate.Unix(),
		"lastUpdateTime":        created.LastModifiedDate.Unix(),
	}, nil
}

func (s *IoTEventsService) DescribeAlarmModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "alarmModelName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	am, err := store.GetAlarmModel(name)
	if err != nil {
		return nil, err
	}

	return alarmModelDescribeResponse(am), nil
}

func (s *IoTEventsService) UpdateAlarmModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "alarmModelName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	def, err := parseStructParam(req.Parameters, "alarmModelDefinition")
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	update := &iotstore.AlarmModel{
		AlarmModelName:        name,
		AlarmModelDescription: request.GetParamCaseInsensitive(req.Parameters, "alarmModelDescription"),
		RoleARN:               request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		AlarmModelDefinition:  def,
		Severity:              request.GetParamCaseInsensitive(req.Parameters, "severity"),
	}

	updated, err := store.UpdateAlarmModel(update)
	if err != nil {
		return nil, err
	}

	return alarmModelDescribeResponse(updated), nil
}

func (s *IoTEventsService) DeleteAlarmModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "alarmModelName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteAlarmModel(name); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTEventsService) ListAlarmModels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListAlarmModels(parseAlarmListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(result.Items))
	for _, am := range result.Items {
		items = append(items, alarmModelSummary(am))
	}

	return listResponse("alarmModelSummaries", items, result.NextMarker), nil
}

func (s *IoTEventsService) ListAlarmModelVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "alarmModelName")
	if name == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	versions, err := store.ListAlarmModelVersions(name)
	if err != nil {
		return nil, err
	}

	return paginatedMaps("alarmModelVersionSummaries", versions, req.Parameters), nil
}

func alarmModelDescribeResponse(am *iotstore.AlarmModel) map[string]interface{} {
	resp := map[string]interface{}{
		"alarmModelName":        am.AlarmModelName,
		"alarmModelArn":         am.AlarmModelARN,
		"alarmModelDescription": am.AlarmModelDescription,
		"roleArn":               am.RoleARN,
		"status":                am.Status,
		"alarmModelVersion":     am.AlarmModelVersion,
		"creationTime":          am.CreationDate.Unix(),
		"lastUpdateTime":        am.LastModifiedDate.Unix(),
	}
	// Severity may be stored as a numeric string or a textual label.
	// The AWS API returns it as an integer when possible.
	if am.Severity != "" {
		if sev, err := strconv.Atoi(am.Severity); err == nil {
			resp["severity"] = int32(sev)
		} else {
			resp["severity"] = am.Severity
		}
	}
	if am.AlarmModelDefinition != nil {
		resp["alarmModelDefinition"] = am.AlarmModelDefinition
	}
	return resp
}

func alarmModelSummary(am *iotstore.AlarmModel) map[string]interface{} {
	return map[string]interface{}{
		"alarmModelName":    am.AlarmModelName,
		"alarmModelArn":     am.AlarmModelARN,
		"alarmModelVersion": am.AlarmModelVersion,
		"creationTime":      am.CreationDate.Unix(),
		"lastUpdateTime":    am.LastModifiedDate.Unix(),
	}
}

func parseAlarmTagsParam(params map[string]interface{}) map[string]string {
	raw, ok := params["tags"]
	if !ok {
		return nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]string, len(list))
	for _, item := range list {
		if m, ok := item.(map[string]interface{}); ok {
			k, _ := m["key"].(string)
			if k == "" {
				k, _ = m["Key"].(string)
			}
			v, _ := m["value"].(string)
			if v == "" {
				v, _ = m["Value"].(string)
			}
			if k != "" {
				result[k] = v
			}
		}
	}
	return result
}

func parseAlarmListOptions(params map[string]interface{}) storecommon.ListOptions {
	opts := storecommon.ListOptions{}
	if v := request.GetParamCaseInsensitive(params, "nextToken"); v != "" {
		opts.Marker = v
	}
	if maxStr := request.GetParamCaseInsensitive(params, "maxResults"); maxStr != "" {
		if n, err := strconv.Atoi(maxStr); err == nil && n > 0 {
			opts.MaxItems = n
		}
	}
	return opts
}

func listResponse(key string, items []map[string]interface{}, nextMarker string) map[string]interface{} {
	return iotutil.ListResponse(key, items, nextMarker)
}
