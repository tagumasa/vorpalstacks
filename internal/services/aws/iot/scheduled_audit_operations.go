package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---- Scheduled Audits ----------------------------------------------

func (s *IoTService) CreateScheduledAudit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, err := s.bulkCreate(reqCtx, "scheduledAudit", req, "scheduledAuditName", map[string]interface{}{
		"frequency":        request.GetParamCaseInsensitive(req.Parameters, "frequency"),
		"dayOfMonth":       request.GetParamCaseInsensitive(req.Parameters, "dayOfMonth"),
		"dayOfWeek":        request.GetParamCaseInsensitive(req.Parameters, "dayOfWeek"),
		"targetCheckNames": request.GetStringList(req.Parameters, "targetCheckNames"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"scheduledAuditArn": iotstore.BuildScheduledAuditARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
	}, nil
}
func (s *IoTService) DeleteScheduledAudit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.bulkDelete(reqCtx, "scheduledAudit", req, "scheduledAuditName"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeScheduledAudit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, _, exists, err := s.bulkGet(reqCtx, "scheduledAudit", req, "scheduledAuditName")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrScheduledAuditNotFound
	}
	return map[string]interface{}{
		"scheduledAuditName": rec["name"],
		"scheduledAuditArn":  iotstore.BuildScheduledAuditARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"frequency":          rec["frequency"],
		"targetCheckNames":   rec["targetCheckNames"],
		"lastModifiedDate":   rec["lastModifiedDate"],
	}, nil
}
func (s *IoTService) ListScheduledAudits(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	items, err := s.bulkList(reqCtx, "scheduledAudit")
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		name, _ := item["name"].(string)
		out = append(out, map[string]interface{}{
			"scheduledAuditName": name,
			"scheduledAuditArn":  iotstore.BuildScheduledAuditARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
			"frequency":          item["frequency"],
		})
	}
	return paginatedMaps("scheduledAudits", out, req.Parameters), nil
}
func (s *IoTService) UpdateScheduledAudit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, exists, err := s.bulkUpdate(reqCtx, "scheduledAudit", req, "scheduledAuditName", map[string]interface{}{
		"frequency":        request.GetParamCaseInsensitive(req.Parameters, "frequency"),
		"dayOfMonth":       request.GetParamCaseInsensitive(req.Parameters, "dayOfMonth"),
		"dayOfWeek":        request.GetParamCaseInsensitive(req.Parameters, "dayOfWeek"),
		"targetCheckNames": request.GetStringList(req.Parameters, "targetCheckNames"),
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrScheduledAuditNotFound
	}
	return map[string]interface{}{
		"scheduledAuditArn": iotstore.BuildScheduledAuditARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
	}, nil
}
