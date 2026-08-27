package iot

import (
	"context"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
)

// parseTimestampParam normalises caller-supplied timestamp strings to Unix
// epoch seconds (int64). Both epoch digit strings and ISO-8601 are
// accepted; an empty or unparseable input returns 0.
func parseTimestampParam(v string) int64 {
	if v == "" {
		return 0
	}
	if n, err := strconv.ParseInt(v, 10, 64); err == nil {
		return n
	}
	if t, err := time.Parse(time.RFC3339, v); err == nil {
		return t.Unix()
	}
	return 0
}

// timestampMemberParam normalises a Timestamp-typed request member to Unix
// epoch seconds. restJson1 serialises timestamps as JSON numbers (epoch
// seconds), but some clients send ISO-8601 strings; both forms are
// accepted. Absent or unparseable values return 0. Reading through
// string-only accessors would silently drop the numeric form.
func timestampMemberParam(params map[string]interface{}, key string) int64 {
	for _, k := range []string{key, request.LowerFirst(key), strings.ToLower(key)} {
		raw, ok := params[k]
		if !ok {
			continue
		}
		switch v := raw.(type) {
		case float64:
			return int64(v)
		case int64:
			return v
		case int:
			return int64(v)
		case string:
			return parseTimestampParam(v)
		}
	}
	return 0
}

// resourceIdentifierFromParams extracts the AWS ResourceIdentifier shape from
// request parameters, preserving only the non-empty members.
func resourceIdentifierFromParams(params map[string]interface{}) map[string]interface{} {
	raw := request.GetMapParamCaseInsensitive(params, "resourceIdentifier")
	if raw == nil {
		raw = map[string]interface{}{}
	}
	out := make(map[string]interface{})
	for k, v := range raw {
		if s, ok := v.(string); ok && s != "" {
			out[k] = s
			continue
		}
		if m, ok := v.(map[string]interface{}); ok && len(m) > 0 {
			out[k] = m
		}
	}
	return out
}

func auditSuppressionResponse(rec *AuditSuppressionRecord) map[string]interface{} {
	return map[string]interface{}{
		"checkName":            rec.CheckName,
		"resourceIdentifier":   rec.ResourceIdentifier,
		"expirationDate":       rec.ExpirationDate,
		"suppressIndefinitely": rec.SuppressIndefinitely,
		"description":          rec.Description,
	}
}

func (s *IoTService) CreateAuditSuppression(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := CreateAuditSuppressionInput{
		CheckName:                    request.GetParamCaseInsensitive(req.Parameters, "checkName"),
		ResourceIdentifier:           resourceIdentifierFromParams(req.Parameters),
		ExpirationDate:               timestampMemberParam(req.Parameters, "expirationDate"),
		ExpirationProvided:           request.HasParam(req.Parameters, "expirationDate"),
		SuppressIndefinitely:         request.GetBoolParam(req.Parameters, "suppressIndefinitely"),
		SuppressIndefinitelyProvided: request.HasParam(req.Parameters, "suppressIndefinitely"),
		Description:                  request.GetParamCaseInsensitive(req.Parameters, "description"),
		ClientRequestToken:           request.GetParamCaseInsensitive(req.Parameters, "clientRequestToken"),
	}
	if err := s.createAuditSuppressionCore(store, in); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DeleteAuditSuppression(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteAuditSuppressionCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "checkName"),
		resourceIdentifierFromParams(req.Parameters)); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeAuditSuppression(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.describeAuditSuppressionCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "checkName"),
		resourceIdentifierFromParams(req.Parameters))
	if err != nil {
		return nil, err
	}
	return auditSuppressionResponse(rec), nil
}
func (s *IoTService) ListAuditSuppressions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	suppressions, err := s.listAuditSuppressionsCore(store)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(suppressions))
	for _, rec := range suppressions {
		items = append(items, auditSuppressionResponse(rec))
	}
	return paginatedMaps("suppressions", items, req.Parameters)
}
func (s *IoTService) UpdateAuditSuppression(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := UpdateAuditSuppressionInput{
		CheckName:                    request.GetParamCaseInsensitive(req.Parameters, "checkName"),
		ResourceIdentifier:           resourceIdentifierFromParams(req.Parameters),
		Description:                  request.GetParamCaseInsensitive(req.Parameters, "description"),
		ExpirationDate:               timestampMemberParam(req.Parameters, "expirationDate"),
		ExpirationProvided:           request.HasParam(req.Parameters, "expirationDate"),
		SuppressIndefinitely:         request.GetBoolParam(req.Parameters, "suppressIndefinitely"),
		SuppressIndefinitelyProvided: request.HasParam(req.Parameters, "suppressIndefinitely"),
	}
	if err := s.updateAuditSuppressionCore(store, in); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
