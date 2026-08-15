package iot

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// auditSuppressionKey builds the generic-KV key for an AuditSuppression
// record. AWS identifies suppressions by the (checkName, resourceIdentifier)
// tuple; resourceIdentifier is a structure with up to ten optional member
// fields, so a canonical-JSON SHA-256 digest gives a stable, collision-free
// suffix without coupling the key layout to AWS-internal representation.
func auditSuppressionKey(checkName string, resourceIdentifier map[string]interface{}) string {
	canonical, _ := json.Marshal(resourceIdentifier)
	digest := sha256.Sum256(canonical)
	return fmt.Sprintf("auditSuppression/%s/%s", checkName, hex.EncodeToString(digest[:])[:16])
}

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

func (s *IoTService) CreateAuditSuppression(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	checkName := request.GetParamCaseInsensitive(req.Parameters, "checkName")
	if checkName == "" {
		return nil, iotstore.ErrMissingParam
	}
	resourceIdentifier := resourceIdentifierFromParams(req.Parameters)
	if len(resourceIdentifier) == 0 {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"checkName":            checkName,
		"resourceIdentifier":   resourceIdentifier,
		"expirationDate":       timestampMemberParam(req.Parameters, "expirationDate"),
		"suppressIndefinitely": request.GetBoolParam(req.Parameters, "suppressIndefinitely"),
		"description":          request.GetParamCaseInsensitive(req.Parameters, "description"),
		"clientRequestToken":   request.GetParamCaseInsensitive(req.Parameters, "clientRequestToken"),
	}
	if err := store.PutGeneric(auditSuppressionKey(checkName, resourceIdentifier), rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DeleteAuditSuppression(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	checkName := request.GetParamCaseInsensitive(req.Parameters, "checkName")
	if checkName == "" {
		return nil, iotstore.ErrMissingParam
	}
	resourceIdentifier := resourceIdentifierFromParams(req.Parameters)
	if len(resourceIdentifier) == 0 {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := auditSuppressionKey(checkName, resourceIdentifier)
	exists, err := store.GetGenericExists(key, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditSuppressionNotFound
	}
	if err := store.DeleteGeneric(key); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeAuditSuppression(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	checkName := request.GetParamCaseInsensitive(req.Parameters, "checkName")
	if checkName == "" {
		return nil, iotstore.ErrMissingParam
	}
	resourceIdentifier := resourceIdentifierFromParams(req.Parameters)
	if len(resourceIdentifier) == 0 {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(auditSuppressionKey(checkName, resourceIdentifier), &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditSuppressionNotFound
	}
	return map[string]interface{}{
		"checkName":            rec["checkName"],
		"resourceIdentifier":   rec["resourceIdentifier"],
		"expirationDate":       rec["expirationDate"],
		"suppressIndefinitely": rec["suppressIndefinitely"],
		"description":          rec["description"],
	}, nil
}
func (s *IoTService) ListAuditSuppressions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("auditSuppression/")
	if err != nil {
		return nil, err
	}
	suppressions := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		suppressions = append(suppressions, map[string]interface{}{
			"checkName":            rec["checkName"],
			"resourceIdentifier":   rec["resourceIdentifier"],
			"expirationDate":       rec["expirationDate"],
			"suppressIndefinitely": rec["suppressIndefinitely"],
			"description":          rec["description"],
		})
	}
	return paginatedMaps("suppressions", suppressions, req.Parameters), nil
}
func (s *IoTService) UpdateAuditSuppression(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	checkName := request.GetParamCaseInsensitive(req.Parameters, "checkName")
	if checkName == "" {
		return nil, iotstore.ErrMissingParam
	}
	resourceIdentifier := resourceIdentifierFromParams(req.Parameters)
	if len(resourceIdentifier) == 0 {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := auditSuppressionKey(checkName, resourceIdentifier)
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditSuppressionNotFound
	}
	// Partial update: only overwrite fields that are explicitly supplied.
	if v := request.GetParamCaseInsensitive(req.Parameters, "description"); v != "" {
		rec["description"] = v
	}
	if request.HasParam(req.Parameters, "expirationDate") {
		// Normalise like the create path so the stored value keeps its
		// numeric epoch form regardless of how the client serialised it.
		rec["expirationDate"] = timestampMemberParam(req.Parameters, "expirationDate")
	}
	if request.HasParam(req.Parameters, "suppressIndefinitely") {
		rec["suppressIndefinitely"] = request.GetBoolParam(req.Parameters, "suppressIndefinitely")
	}
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
