package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
)

// ---- Scheduled Audits ----------------------------------------------

func (s *IoTService) CreateScheduledAudit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	saTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		saTags[t.Key] = t.Value
	}
	arn, err := s.createScheduledAuditCore(store, ScheduledAuditInput{
		Name:                     request.GetParamCaseInsensitive(req.Parameters, "scheduledAuditName"),
		Frequency:                request.GetParamCaseInsensitive(req.Parameters, "frequency"),
		FrequencyProvided:        true,
		DayOfMonth:               request.GetParamCaseInsensitive(req.Parameters, "dayOfMonth"),
		DayOfMonthProvided:       hasParam(req.Parameters, "dayOfMonth"),
		DayOfWeek:                request.GetParamCaseInsensitive(req.Parameters, "dayOfWeek"),
		DayOfWeekProvided:        hasParam(req.Parameters, "dayOfWeek"),
		TargetCheckNames:         request.GetStringList(req.Parameters, "targetCheckNames"),
		TargetCheckNamesProvided: true,
		Tags:                     saTags,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"scheduledAuditArn": arn,
	}, nil
}
func (s *IoTService) DeleteScheduledAudit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteScheduledAuditCore(store, request.GetParamCaseInsensitive(req.Parameters, "scheduledAuditName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeScheduledAudit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.describeScheduledAuditCore(store, request.GetParamCaseInsensitive(req.Parameters, "scheduledAuditName"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"scheduledAuditName": rec.Rec["name"],
		"scheduledAuditArn":  rec.Arn,
		"frequency":          rec.Rec["frequency"],
		"dayOfMonth":         rec.Rec["dayOfMonth"],
		"dayOfWeek":          rec.Rec["dayOfWeek"],
		"targetCheckNames":   rec.Rec["targetCheckNames"],
	}, nil
}
func (s *IoTService) ListScheduledAudits(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := s.listScheduledAuditsCore(store)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		out = append(out, map[string]interface{}{
			"scheduledAuditName": item.Name,
			"scheduledAuditArn":  item.Arn,
			"frequency":          item.Frequency,
		})
	}
	return paginatedMaps("scheduledAudits", out, req.Parameters)
}
func (s *IoTService) UpdateScheduledAudit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn, err := s.updateScheduledAuditCore(store, ScheduledAuditInput{
		Name:                     request.GetParamCaseInsensitive(req.Parameters, "scheduledAuditName"),
		Frequency:                request.GetParamCaseInsensitive(req.Parameters, "frequency"),
		FrequencyProvided:        hasParam(req.Parameters, "frequency"),
		DayOfMonth:               request.GetParamCaseInsensitive(req.Parameters, "dayOfMonth"),
		DayOfMonthProvided:       hasParam(req.Parameters, "dayOfMonth"),
		DayOfWeek:                request.GetParamCaseInsensitive(req.Parameters, "dayOfWeek"),
		DayOfWeekProvided:        hasParam(req.Parameters, "dayOfWeek"),
		TargetCheckNames:         request.GetStringList(req.Parameters, "targetCheckNames"),
		TargetCheckNamesProvided: hasParam(req.Parameters, "targetCheckNames"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"scheduledAuditArn": arn,
	}, nil
}
