package iot

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/services/aws/iot/policy"
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

// parseTimestampParam normalises caller-supplied timestamp values to Unix
// epoch seconds (int64). restJson1 Timestamp traits serialise as JSON numbers,
// but some clients send ISO-8601 strings; both are accepted here. An empty
// or unparseable input returns 0.
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

// Handlers for Security, Detect, Audit, Fleet Indexing, Logging,
// Encryption and Topic Rule Destinations. State is persisted via the generic-KV
// store (Pebble-backed), so it survives restarts. The CRUD entities
// (CustomMetric, Dimension, MitigationAction, FleetMetric, ScheduledAudit,
// TopicRuleDestination) implement full create/read/update/delete with field
// echo, NotFound semantics and error propagation, matching the AWS IoT
// Control Plane. The asynchronous task families (Detect tasks, on-demand
// Audit tasks) and the fleet index aggregator require dedicated engines that
// are out of scope; those handlers remain structural stubs.

// ---- bulk generic-KV helpers ------------------------------------------------

// bulkCreate persists a new record under "<category>/<name>", capturing the
// supplied extra fields alongside the name and timestamps. Returns the stored
// record so callers can shape the AWS response. Store errors are propagated.
func (s *IoTService) bulkCreate(reqCtx *request.RequestContext, category string, req *request.ParsedRequest, nameKey string, extra map[string]interface{}) (map[string]interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetParamCaseInsensitive(req.Parameters, nameKey)
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	// Reject duplicates: AWS returns ResourceAlreadyExistsException (409)
	// when a named resource already exists. Only check when the caller
	// supplied a name (auto-generated UUIDs are inherently unique).
	if userProvided := request.GetParamCaseInsensitive(req.Parameters, nameKey); userProvided != "" {
		exists, err := store.GetGenericExists(category+"/"+name, &map[string]interface{}{})
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, iotstore.ErrResourceAlreadyExists
		}
	}
	now := time.Now().Unix()
	rec := map[string]interface{}{
		"name":             name,
		"creationDate":     now,
		"lastModifiedDate": now,
	}
	for k, v := range extra {
		if v != nil {
			rec[k] = v
		}
	}
	if err := store.PutGeneric(category+"/"+name, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// bulkGet loads a single record. exists is false (with nil error) when the
// record is absent, fixing the previous dead logic where a missing key left a
// pre-allocated map in place and was reported as found.
func (s *IoTService) bulkGet(reqCtx *request.RequestContext, category string, req *request.ParsedRequest, nameKey string) (rec map[string]interface{}, name string, exists bool, err error) {
	name = request.GetParamCaseInsensitive(req.Parameters, nameKey)
	store, sErr := s.store(reqCtx)
	if sErr != nil {
		return nil, name, false, sErr
	}
	rec = map[string]interface{}{}
	exists, err = store.GetGenericExists(category+"/"+name, &rec)
	if err != nil {
		return nil, name, false, err
	}
	if !exists {
		return nil, name, false, nil
	}
	return rec, name, true, nil
}

// bulkUpdate merges the supplied fields into an existing record and refreshes
// lastModifiedDate. exists is false (nil error) when absent, so callers return
// ResourceNotFoundException.
func (s *IoTService) bulkUpdate(reqCtx *request.RequestContext, category string, req *request.ParsedRequest, nameKey string, merge map[string]interface{}) (map[string]interface{}, bool, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, nameKey)
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, false, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(category+"/"+name, &rec)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	for k, v := range merge {
		if v != nil {
			rec[k] = v
		}
	}
	rec["lastModifiedDate"] = time.Now().Unix()
	if err := store.PutGeneric(category+"/"+name, rec); err != nil {
		return nil, false, err
	}
	return rec, true, nil
}

// bulkDelete removes a record. BaseStore.Delete is idempotent (no error on a
// missing key), matching AWS IoT delete semantics; only genuine store errors
// propagate.
func (s *IoTService) bulkDelete(reqCtx *request.RequestContext, category string, req *request.ParsedRequest, nameKey string) error {
	name := request.GetParamCaseInsensitive(req.Parameters, nameKey)
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	return store.DeleteGeneric(category + "/" + name)
}

// bulkList lists all records under a category prefix.
func (s *IoTService) bulkList(reqCtx *request.RequestContext, category string) ([]map[string]interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return store.ListGeneric(category + "/")
}

func bulkName(rec map[string]interface{}) string {
	if name, ok := rec["name"].(string); ok {
		return name
	}
	return ""
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

// ---- Custom Metrics -----------------------------------------------

func (s *IoTService) CreateCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}
	rec, err := s.bulkCreate(reqCtx, "customMetric", req, "metricName", map[string]interface{}{
		"metricType":         request.GetParamCaseInsensitive(req.Parameters, "metricType"),
		"displayName":        request.GetParamCaseInsensitive(req.Parameters, "displayName"),
		"clientRequestToken": request.GetParamCaseInsensitive(req.Parameters, "clientRequestToken"),
		"tags":               recTags,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"metricName": rec["name"],
		"metricArn":  iotstore.BuildCustomMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
	}, nil
}
func (s *IoTService) DeleteCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.bulkDelete(reqCtx, "customMetric", req, "metricName"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, _, exists, err := s.bulkGet(reqCtx, "customMetric", req, "metricName")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCustomMetricNotFound
	}
	return map[string]interface{}{
		"metricName":       rec["name"],
		"metricArn":        iotstore.BuildCustomMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"metricType":       rec["metricType"],
		"displayName":      rec["displayName"],
		"creationDate":     rec["creationDate"],
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}
func (s *IoTService) ListCustomMetrics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy: ListCustomMetricsResponse.metricNames is a list of MetricName (string).
	items, err := s.bulkList(reqCtx, "customMetric")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(items))
	for _, item := range items {
		if name, ok := item["name"].(string); ok {
			names = append(names, name)
		}
	}
	return paginatedStrings("metricNames", names, req.Parameters), nil
}
func (s *IoTService) UpdateCustomMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, exists, err := s.bulkUpdate(reqCtx, "customMetric", req, "metricName", map[string]interface{}{
		"displayName": request.GetParamCaseInsensitive(req.Parameters, "displayName"),
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCustomMetricNotFound
	}
	return map[string]interface{}{
		"metricName":       rec["name"],
		"metricArn":        iotstore.BuildCustomMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"metricType":       rec["metricType"],
		"displayName":      rec["displayName"],
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}

// ---- Dimensions ----------------------------------------------------

func (s *IoTService) CreateDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, err := s.bulkCreate(reqCtx, "dimension", req, "name", map[string]interface{}{
		"type":         request.GetParamCaseInsensitive(req.Parameters, "type"),
		"stringValues": request.GetStringList(req.Parameters, "stringValues"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name": rec["name"],
		"arn":  iotstore.BuildDimensionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
	}, nil
}
func (s *IoTService) DeleteDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.bulkDelete(reqCtx, "dimension", req, "name"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, _, exists, err := s.bulkGet(reqCtx, "dimension", req, "name")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDimensionNotFound
	}
	return map[string]interface{}{
		"name":             rec["name"],
		"arn":              iotstore.BuildDimensionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"type":             rec["type"],
		"stringValues":     rec["stringValues"],
		"creationDate":     rec["creationDate"],
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}
func (s *IoTService) ListDimensions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy: ListDimensionsResponse.dimensionNames is list<DimensionName> (string).
	items, err := s.bulkList(reqCtx, "dimension")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(items))
	for _, item := range items {
		if name, ok := item["name"].(string); ok {
			names = append(names, name)
		}
	}
	return paginatedStrings("dimensionNames", names, req.Parameters), nil
}
func (s *IoTService) UpdateDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, exists, err := s.bulkUpdate(reqCtx, "dimension", req, "name", map[string]interface{}{
		"stringValues": request.GetStringList(req.Parameters, "stringValues"),
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDimensionNotFound
	}
	return map[string]interface{}{
		"name":             rec["name"],
		"arn":              iotstore.BuildDimensionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"stringValues":     rec["stringValues"],
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}

// ---- Mitigation Actions --------------------------------------------

func (s *IoTService) CreateMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := request.GetMapParamCaseInsensitive(req.Parameters, "actionParams")
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}
	rec, err := s.bulkCreate(reqCtx, "mitigationAction", req, "actionName", map[string]interface{}{
		"actionType":   deriveMitigationActionType(params),
		"actionParams": params,
		"roleArn":      request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"tags":         recTags,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"actionArn": iotstore.BuildMitigationActionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"actionId":  uuid.New().String(),
	}, nil
}

// deriveMitigationActionType infers the action type from the actionParams keys.
// AWS derives the type automatically based on which params member is set.
func deriveMitigationActionType(params map[string]interface{}) string {
	if _, ok := params["addThingsToThingGroupParams"]; ok {
		return "ADD_THINGS_TO_THING_GROUP"
	}
	if _, ok := params["enableIoTLoggingParams"]; ok {
		return "ENABLE_IOT_LOGGING"
	}
	if _, ok := params["publishFindingToSnsParams"]; ok {
		return "PUBLISH_FINDING_TO_SNS"
	}
	if _, ok := params["addThingsToCertPoolParams"]; ok {
		return "ADD_THINGS_TO_CERTIFICATE_POOL"
	}
	if _, ok := params["replaceCACertificateParams"]; ok {
		return "REPLACE_CA_CERTIFICATE"
	}
	if _, ok := params["updateDeviceCertificateParams"]; ok {
		return "UPDATE_DEVICE_CERTIFICATE"
	}
	return ""
}
func (s *IoTService) DeleteMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.bulkDelete(reqCtx, "mitigationAction", req, "actionName"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, _, exists, err := s.bulkGet(reqCtx, "mitigationAction", req, "actionName")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrMitigationActionNotFound
	}
	return map[string]interface{}{
		"actionName":       rec["name"],
		"actionType":       rec["actionType"],
		"actionParams":     rec["actionParams"],
		"roleArn":          rec["roleArn"],
		"actionArn":        iotstore.BuildMitigationActionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"creationDate":     rec["creationDate"],
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}
func (s *IoTService) ListMitigationActions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	items, err := s.bulkList(reqCtx, "mitigationAction")
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		name, _ := item["name"].(string)
		out = append(out, map[string]interface{}{
			"actionName":   name,
			"actionArn":    iotstore.BuildMitigationActionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
			"creationDate": item["creationDate"],
		})
	}
	return paginatedMaps("actionIdentifiers", out, req.Parameters), nil
}
func (s *IoTService) UpdateMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, exists, err := s.bulkUpdate(reqCtx, "mitigationAction", req, "actionName", map[string]interface{}{
		"roleArn":      request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"actionParams": request.GetMapParamCaseInsensitive(req.Parameters, "actionParams"),
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrMitigationActionNotFound
	}
	return map[string]interface{}{
		"actionName":       rec["name"],
		"actionType":       rec["actionType"],
		"actionParams":     rec["actionParams"],
		"roleArn":          rec["roleArn"],
		"actionArn":        iotstore.BuildMitigationActionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}

// ---- Detect Mitigation Actions Tasks --------------------------------
// These handlers persist task records so that Cancel/Describe can resolve the
// identifier and return ResourceNotFoundException for unknown task ids,
// matching the Smithy error trait set on each operation.

func (s *IoTService) StartDetectMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		taskId = uuid.New().String()
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"taskId":         taskId,
		"status":         "IN_PROGRESS",
		"startTime":      time.Now().UTC().Unix(),
		"target":         request.GetParamCaseInsensitive(req.Parameters, "target"),
		"actions":        request.GetParamCaseInsensitive(req.Parameters, "actions"),
		"violationEvent": request.GetParamCaseInsensitive(req.Parameters, "violationEvent"),
	}
	if err := store.PutGeneric("detectMitigationTask/"+taskId, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{"taskId": taskId}, nil
}
func (s *IoTService) CancelDetectMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "detectMitigationTask/" + taskId
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDetectMitigationTaskNotFound
	}
	rec["status"] = "CANCELED"
	rec["endTime"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeDetectMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("detectMitigationTask/"+taskId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDetectMitigationTaskNotFound
	}
	return map[string]interface{}{
		"taskId":         rec["taskId"],
		"status":         rec["status"],
		"startTime":      rec["startTime"],
		"endTime":        rec["endTime"],
		"target":         rec["target"],
		"actions":        rec["actions"],
		"violationEvent": rec["violationEvent"],
	}, nil
}
func (s *IoTService) ListDetectMitigationActionsExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return paginatedMaps("taskExecutions", []map[string]interface{}{}, req.Parameters), nil
}
func (s *IoTService) ListDetectMitigationActionsTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("detectMitigationTask/")
	if err != nil {
		return nil, err
	}
	tasks := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		tasks = append(tasks, map[string]interface{}{
			"taskId":        rec["taskId"],
			"taskStatus":    rec["status"],
			"taskStartTime": rec["startTime"],
		})
	}
	return paginatedMaps("tasks", tasks, req.Parameters), nil
}

// ---- Audit (task/findings) ------------------------------------------
// Audit task lifecycle mirrors the Detect Mitigation task pattern: Start
// persists the task id, Cancel/Describe enforce ResourceNotFoundException for
// unknown ids. Audit findings are not generated without a Defender engine,
// so DescribeAuditFinding always returns NotFound for an arbitrary id.

func (s *IoTService) DescribeAccountAuditConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/accountAudit", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]interface{}{
			"auditCheckConfigurations":              map[string]interface{}{},
			"auditNotificationTargetConfigurations": map[string]interface{}{},
		}, nil
	}
	return rec, nil
}
func (s *IoTService) UpdateAccountAuditConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"roleArn":                               request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"auditCheckConfigurations":              request.GetMapParamCaseInsensitive(req.Parameters, "auditCheckConfigurations"),
		"auditNotificationTargetConfigurations": request.GetMapParamCaseInsensitive(req.Parameters, "auditNotificationTargetConfigurations"),
	}
	if err := store.PutGeneric("config/accountAudit", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DeleteAccountAuditConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteGeneric("config/accountAudit"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) StartOnDemandAuditTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := uuid.New().String()
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"taskId":         taskId,
		"status":         "IN_PROGRESS",
		"startTime":      time.Now().UTC().Unix(),
		"targetAccounts": request.GetParamCaseInsensitive(req.Parameters, "targetAccounts"),
		"auditChecks":    request.GetParamCaseInsensitive(req.Parameters, "auditChecks"),
	}
	if err := store.PutGeneric("auditTask/"+taskId, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{"taskId": taskId}, nil
}
func (s *IoTService) CancelAuditTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "auditTask/" + taskId
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditTaskNotFound
	}
	rec["status"] = "CANCELED"
	rec["endTime"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeAuditTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditTask/"+taskId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditTaskNotFound
	}
	return map[string]interface{}{
		"taskId":           rec["taskId"],
		"taskStatus":       rec["status"],
		"auditTaskDetails": rec,
	}, nil
}
func (s *IoTService) ListAuditTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("auditTask/")
	if err != nil {
		return nil, err
	}
	tasks := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		tasks = append(tasks, map[string]interface{}{
			"taskId":     rec["taskId"],
			"taskStatus": rec["status"],
			"taskType":   "ON_DEMAND_AUDIT_TASK",
		})
	}
	return paginatedMaps("tasks", tasks, req.Parameters), nil
}
func (s *IoTService) DescribeAuditFinding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	findingId := request.GetParamCaseInsensitive(req.Parameters, "findingId")
	if findingId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditFinding/"+findingId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		// No Defender engine generates findings, so any caller-supplied id
		// is unknown to the platform. AWS returns ResourceNotFoundException.
		return nil, iotstore.ErrAuditFindingNotFound
	}
	return map[string]interface{}{"finding": rec}, nil
}
func (s *IoTService) ListAuditFindings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("auditFinding/")
	if err != nil {
		return nil, err
	}
	findings := make([]map[string]interface{}, 0, len(items))
	findings = append(findings, items...)
	return paginatedMaps("findings", findings, req.Parameters), nil
}
func (s *IoTService) ListRelatedResourcesForAuditFinding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	findingID := request.GetParamCaseInsensitive(req.Parameters, "findingId")
	if findingID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditFinding/"+findingID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditFindingNotFound
	}
	resources := []map[string]interface{}{}
	if raw, ok := rec["relatedResources"].([]interface{}); ok {
		for _, r := range raw {
			if m, ok := r.(map[string]interface{}); ok {
				resources = append(resources, m)
			}
		}
	}
	return paginatedMaps("relatedResources", resources, req.Parameters), nil
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
		"expirationDate":       parseTimestampParam(request.GetParamCaseInsensitive(req.Parameters, "expirationDate")),
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
	for _, field := range []string{"expirationDate", "description"} {
		if v := request.GetParamCaseInsensitive(req.Parameters, field); v != "" {
			rec[field] = v
		}
	}
	if request.HasParam(req.Parameters, "suppressIndefinitely") {
		rec["suppressIndefinitely"] = request.GetBoolParam(req.Parameters, "suppressIndefinitely")
	}
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) StartAuditMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"taskId":    taskId,
		"status":    "IN_PROGRESS",
		"startTime": time.Now().UTC().Unix(),
	}
	if err := store.PutGeneric("auditMitigationTask/"+taskId, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) CancelAuditMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "auditMitigationTask/" + taskId
	exists, err := store.GetGenericExists(key, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditMitigationTaskNotFound
	}
	if err := store.DeleteGeneric(key); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeAuditMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditMitigationTask/"+taskId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditMitigationTaskNotFound
	}
	return rec, nil
}
func (s *IoTService) ListAuditMitigationActionsExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return paginatedMaps("taskExecutions", []map[string]interface{}{}, req.Parameters), nil
}
func (s *IoTService) ListAuditMitigationActionsTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("auditMitigationTask/")
	if err != nil {
		return nil, err
	}
	tasks := make([]map[string]interface{}, 0, len(items))
	tasks = append(tasks, items...)
	return paginatedMaps("tasks", tasks, req.Parameters), nil
}

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
		"auditChecks":        rec["auditChecksToInclude"],
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
func (s *IoTService) AssociateTargetsWithJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobId := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	if jobId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	job, err := store.GetJob(jobId)
	if err != nil {
		return nil, err
	}
	if job == nil || job.JobID == "" {
		return nil, iotstore.ErrJobNotFound
	}
	// Merge new targets into the existing target list (de-duplicated).
	newTargets := request.GetStringList(req.Parameters, "targets")
	seen := make(map[string]bool, len(job.Targets))
	for _, t := range job.Targets {
		seen[t] = true
	}
	for _, t := range newTargets {
		if !seen[t] {
			job.Targets = append(job.Targets, t)
			seen[t] = true
		}
	}
	if comment := request.GetParamCaseInsensitive(req.Parameters, "comment"); comment != "" && job.Description == "" {
		job.Description = comment
	}
	if _, err := store.UpdateJob(jobId, iotstore.JobUpdateOpts{
		Description: job.Description,
		Targets:     job.Targets,
	}); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"jobArn":      job.JobARN,
		"jobId":       job.JobID,
		"description": job.Description,
	}, nil
}
func (s *IoTService) TestAuthorization(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	principal := request.GetParamCaseInsensitive(req.Parameters, "principal")
	clientID := request.GetParamCaseInsensitive(req.Parameters, "clientId")
	cognitoIdentityPoolId := request.GetParamCaseInsensitive(req.Parameters, "cognitoIdentityPoolId")
	policyNamesToAdd := request.GetStringList(req.Parameters, "policyNamesToAdd")
	policyNamesToSkip := request.GetStringList(req.Parameters, "policyNamesToSkip")
	if principal == "" {
		return map[string]interface{}{"authResults": []map[string]interface{}{}}, nil
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Resolve the policies attached to the principal and parse them for
	// evaluation, mirroring the MQTT broker auth hook (auth_hooks.go).
	policyNames, err := store.ListPoliciesForPrincipal(principal)
	if err != nil {
		return nil, err
	}
	// Merge additional policies supplied by the caller (policyNamesToAdd).
	skipSet := make(map[string]bool, len(policyNamesToSkip))
	for _, n := range policyNamesToSkip {
		skipSet[n] = true
	}
	for _, n := range policyNamesToAdd {
		if !skipSet[n] {
			policyNames = append(policyNames, n)
		}
	}
	// When cognitoIdentityPoolId is supplied, also include policies attached
	// to the cognito identity pool principal.
	if cognitoIdentityPoolId != "" {
		cognitoPolicies, err := store.ListPoliciesForPrincipal(cognitoIdentityPoolId)
		if err == nil {
			for _, n := range cognitoPolicies {
				if !skipSet[n] {
					policyNames = append(policyNames, n)
				}
			}
		}
	}
	var versions []*policy.PolicyVersion
	matchedNames := make([]string, 0, len(policyNames))
	for _, name := range policyNames {
		if skipSet[name] {
			continue
		}
		p, gErr := store.GetPolicy(name)
		if gErr != nil {
			continue
		}
		pv, pErr := policy.ParsePolicyVersion([]byte(p.PolicyDocument))
		if pErr != nil {
			continue
		}
		versions = append(versions, pv)
		matchedNames = append(matchedNames, name)
	}
	// Evaluate each requested action/resource combination (authInfos). When no
	// authInfos are supplied, return the matched policies with an empty result.
	var authInfos []map[string]interface{}
	if raw, ok := req.Parameters["authInfos"].([]interface{}); ok {
		for _, item := range raw {
			if m, ok := item.(map[string]interface{}); ok {
				authInfos = append(authInfos, m)
			}
		}
	}
	policyEntries := make([]map[string]interface{}, 0, len(matchedNames))
	for _, n := range matchedNames {
		policyEntries = append(policyEntries, map[string]interface{}{
			"policyName": n,
			"policyArn":  iotstore.BuildPolicyARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), n),
		})
	}
	results := make([]map[string]interface{}, 0, len(authInfos))
	for _, info := range authInfos {
		action, _ := info["actionType"].(string)
		resources := request.GetStringList(info, "resources")
		evalResource := "*"
		if len(resources) > 0 {
			evalResource = resources[0]
		}
		// Normalise the action to "iot:TitleCase" for policy evaluation.
		// The CLI sends lowercase (e.g. "connect") but policies use "iot:Connect".
		iotAction := action
		if action != "" && !strings.HasPrefix(action, "iot:") {
			iotAction = "iot:" + strings.ToUpper(action[:1]) + strings.ToLower(action[1:])
		}
		allowed, _ := policy.Evaluate(&policy.EvaluateRequest{
			Policies: versions,
			Action:   iotAction,
			Resource: evalResource,
			ClientID: clientID,
		})
		entry := map[string]interface{}{
			"authInfo": map[string]interface{}{
				"actionType": action,
				"resources":  resources,
			},
			"matchedPolicies": matchedNames,
		}
		if allowed {
			entry["allowed"] = map[string]interface{}{"policies": policyEntries}
			entry["decision"] = "ALLOWED"
		} else {
			entry["denied"] = map[string]interface{}{
				"implicitDeny": map[string]interface{}{"policies": policyEntries},
			}
			entry["decision"] = "IMPLICIT_DENY"
		}
		results = append(results, entry)
	}
	return map[string]interface{}{"authResults": results}, nil
}

// ---- Fleet Indexing / Metrics --------------------------------------

func (s *IoTService) DescribeIndex(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "indexName")
	return map[string]interface{}{"indexName": name, "indexStatus": "ACTIVE"}, nil
}
func (s *IoTService) ListIndices(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy: ListIndicesResponse.indexNames is a list of IndexName (string).
	return map[string]interface{}{"indexNames": []string{"AWS_Things"}}, nil
}
func (s *IoTService) SearchIndex(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.searchIndexImpl(ctx, reqCtx, req)
}
func (s *IoTService) GetBucketsAggregation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getBucketsAggregationImpl(ctx, reqCtx, req)
}
func (s *IoTService) GetCardinality(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getCardinalityImpl(ctx, reqCtx, req)
}
func (s *IoTService) GetPercentiles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getPercentilesImpl(ctx, reqCtx, req)
}
func (s *IoTService) GetStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getStatisticsImpl(ctx, reqCtx, req)
}
func (s *IoTService) ListMetricValues(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	metricName := request.GetParamCaseInsensitive(req.Parameters, "metricName")
	if metricName == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec, _, exists, err := s.bulkGet(reqCtx, "fleetMetric", req, "metricName")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrFleetMetricNotFound
	}

	// Execute the fleet metric's query to compute the current value.
	queryString, _ := rec["queryString"].(string)
	aggField, _ := rec["aggregationField"].(string)
	aggType, _ := rec["aggregationType"].(string)

	matched, err := s.searchThings(reqCtx, queryString)
	if err != nil {
		return nil, err
	}

	var value float64
	var count int64
	if aggField != "" {
		var values []float64
		for _, t := range matched {
			if v := getNumericAttribute(t, aggField); !math.IsNaN(v) {
				values = append(values, v)
			}
		}
		count = int64(len(values))
		if count > 0 {
			switch strings.ToUpper(aggType) {
			case "AVERAGE":
				sum := 0.0
				for _, v := range values {
					sum += v
				}
				value = sum / float64(count)
			case "SUM":
				for _, v := range values {
					value += v
				}
			case "MINIMUM":
				value = values[0]
				for _, v := range values {
					if v < value {
						value = v
					}
				}
			case "MAXIMUM":
				value = values[0]
				for _, v := range values {
					if v > value {
						value = v
					}
				}
			default: // COUNT or unspecified
				value = float64(count)
			}
		}
	} else {
		value = float64(len(matched))
		count = int64(len(matched))
	}

	now := time.Now().UTC().Unix()
	values := []map[string]interface{}{
		{
			"timestamp": now,
			"value":     value,
		},
	}
	return paginatedMaps("metricValues", values, req.Parameters), nil
}

func (s *IoTService) CreateFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, err := s.bulkCreate(reqCtx, "fleetMetric", req, "metricName", map[string]interface{}{
		"queryString":      request.GetParamCaseInsensitive(req.Parameters, "queryString"),
		"aggregationType":  request.GetMapParamCaseInsensitive(req.Parameters, "aggregationType"),
		"period":           int64(request.GetIntParam(req.Parameters, "period")),
		"aggregationField": request.GetParamCaseInsensitive(req.Parameters, "aggregationField"),
		"unit":             request.GetParamCaseInsensitive(req.Parameters, "unit"),
		"indexName":        request.GetParamCaseInsensitive(req.Parameters, "indexName"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"metricName": rec["name"],
		"metricArn":  iotstore.BuildFleetMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
	}, nil
}
func (s *IoTService) DeleteFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.bulkDelete(reqCtx, "fleetMetric", req, "metricName"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, _, exists, err := s.bulkGet(reqCtx, "fleetMetric", req, "metricName")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrFleetMetricNotFound
	}
	return map[string]interface{}{
		"metricName":       rec["name"],
		"metricArn":        iotstore.BuildFleetMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"queryString":      rec["queryString"],
		"aggregationType":  rec["aggregationType"],
		"period":           rec["period"],
		"aggregationField": rec["aggregationField"],
		"unit":             rec["unit"],
		"indexName":        rec["indexName"],
		"creationDate":     rec["creationDate"],
		"lastModifiedDate": rec["lastModifiedDate"],
		"version":          int64(1),
	}, nil
}
func (s *IoTService) ListFleetMetrics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	items, err := s.bulkList(reqCtx, "fleetMetric")
	if err != nil {
		return nil, err
	}
	// Transform internal records to the AWS FleetMetricNameAndArn summary
	// shape. Without this the response items are empty objects because the
	// internal field names ("name") do not match the expected output members
	// ("metricName", "metricArn").
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		name := bulkName(item)
		result = append(result, map[string]interface{}{
			"metricName": name,
			"metricArn":  iotstore.BuildFleetMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
		})
	}
	return paginatedMaps("fleetMetrics", result, req.Parameters), nil
}
func (s *IoTService) UpdateFleetMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, exists, err := s.bulkUpdate(reqCtx, "fleetMetric", req, "metricName", map[string]interface{}{
		"queryString":      request.GetParamCaseInsensitive(req.Parameters, "queryString"),
		"aggregationType":  request.GetMapParamCaseInsensitive(req.Parameters, "aggregationType"),
		"period":           int64(request.GetIntParam(req.Parameters, "period")),
		"aggregationField": request.GetParamCaseInsensitive(req.Parameters, "aggregationField"),
		"unit":             request.GetParamCaseInsensitive(req.Parameters, "unit"),
		"indexName":        request.GetParamCaseInsensitive(req.Parameters, "indexName"),
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrFleetMetricNotFound
	}
	return map[string]interface{}{
		"metricName":       rec["name"],
		"metricArn":        iotstore.BuildFleetMetricARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}

// ---- Logging / Event / Encryption config --------------------------
// Persisted via GenericKV under "config/<name>". A missing key means "not yet
// configured"; the handlers return a default/empty shape for that case and
// propagate genuine store errors.

func (s *IoTService) GetV2LoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/v2Logging", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]interface{}{
			"defaultLogLevel": "DISABLED",
			"disableAllLogs":  true,
		}, nil
	}
	if rec["defaultLogLevel"] == nil || rec["defaultLogLevel"] == "" {
		rec["defaultLogLevel"] = "INFO"
	}
	// GetV2LoggingOptions output shape is flat (roleArn, defaultLogLevel,
	// disableAllLogs at the top level). Wrapping in "loggingOptions" causes
	// the AWS SDK parser to discard all fields.
	return rec, nil
}
func (s *IoTService) SetV2LoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	defaultLogLevel := request.GetParamCaseInsensitive(req.Parameters, "defaultLogLevel")
	if defaultLogLevel == "" {
		defaultLogLevel = "INFO"
	}
	rec := map[string]interface{}{
		"roleArn":         request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"defaultLogLevel": defaultLogLevel,
		"disableAllLogs":  request.GetBoolParam(req.Parameters, "disableAllLogs"),
	}
	if err := store.PutGeneric("config/v2Logging", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DeleteV2LoggingLevel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	targetType := request.GetParamCaseInsensitive(req.Parameters, "targetType")
	targetName := request.GetParamCaseInsensitive(req.Parameters, "targetName")
	if targetType == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteGeneric("v2LoggingLevel/" + targetType + "/" + targetName); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) ListV2LoggingLevels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	prefix := "v2LoggingLevel/"
	if tt := request.GetParamCaseInsensitive(req.Parameters, "targetType"); tt != "" {
		prefix = "v2LoggingLevel/" + tt + "/"
	}
	items, err := store.ListGeneric(prefix)
	if err != nil {
		return nil, err
	}
	configs := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		logTarget, _ := rec["logTarget"].(map[string]interface{})
		if logTarget == nil {
			logTarget = map[string]interface{}{}
		}
		configs = append(configs, map[string]interface{}{
			"logTarget": logTarget,
			"logLevel":  rec["logLevel"],
		})
	}
	return paginatedMaps("logTargetConfigurations", configs, req.Parameters), nil
}
func (s *IoTService) SetV2LoggingLevel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logTarget := request.GetMapParamCaseInsensitive(req.Parameters, "logTarget")
	logLevel := request.GetParamCaseInsensitive(req.Parameters, "logLevel")
	targetType, _ := logTarget["targetType"].(string)
	targetName, _ := logTarget["targetName"].(string)
	if targetType == "" || logLevel == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.PutGeneric("v2LoggingLevel/"+targetType+"/"+targetName, map[string]interface{}{
		"logTarget": logTarget,
		"logLevel":  logLevel,
	}); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) GetLoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/legacyLogging", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		// AWS returns an empty response when logging has never been configured.
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{
		"roleArn":  rec["roleArn"],
		"logLevel": rec["logLevel"],
	}, nil
}
func (s *IoTService) SetLoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// AWS uses httpPayload on LoggingOptionsPayload so the body is wrapped.
	props := request.GetMapParamCaseInsensitive(req.Parameters, "loggingOptionsPayload")
	if props == nil {
		// Some SDKs send the members flat (no wrapper); accept either form.
		props = req.Parameters
	}
	roleArn := request.GetParamCaseInsensitive(props, "roleArn")
	if roleArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{
		"roleArn":  roleArn,
		"logLevel": request.GetParamCaseInsensitive(props, "logLevel"),
	}
	if err := store.PutGeneric("config/legacyLogging", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeEventConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	if _, err := store.GetGenericExists("config/eventConfigurations", &rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{"eventConfigurations": rec}, nil
}
func (s *IoTService) UpdateEventConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	if _, err := store.GetGenericExists("config/eventConfigurations", &rec); err != nil {
		return nil, err
	}
	// Merge incoming configuration attributes into the persisted map so that
	// partial updates behave like AWS IoT (per-event-type toggles).
	if incoming, ok := req.Parameters["eventConfigurations"].(map[string]interface{}); ok {
		for k, v := range incoming {
			rec[k] = v
		}
	}
	if err := store.PutGeneric("config/eventConfigurations", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeEncryptionConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/encryptionConfiguration", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]interface{}{
			"encryptionType": "TLS",
		}, nil
	}
	// DescribeEncryptionConfiguration output shape is flat (encryptionType,
	// kmsKeyArn at the top level). Wrapping in "encryptionConfiguration"
	// causes the AWS SDK parser to discard all fields.
	return rec, nil
}
func (s *IoTService) UpdateEncryptionConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"kmsKeyArn":        request.GetParamCaseInsensitive(req.Parameters, "kmsKeyArn"),
		"kmsAccessRoleArn": request.GetParamCaseInsensitive(req.Parameters, "kmsAccessRoleArn"),
	}
	et := request.GetParamCaseInsensitive(req.Parameters, "encryptionType")
	if et == "" {
		et = "TLS"
	}
	rec["encryptionType"] = et
	if err := store.PutGeneric("config/encryptionConfiguration", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

// ---- Topic Rule Destinations --------------------------------------
// AWS identifies destinations by ARN (auto-generated at create time, derived
// from a UUID) and resolves Confirm by confirmationToken. The earlier
// "destinationName" keying was a misreading of the Smithy model: neither
// CreateTopicRuleDestinationRequest nor any other request shape carries a
// destinationName member. The handlers below persist by ARN and produce the
// canonical TopicRuleDestination response shape.

// topicRuleDestinationResponse shapes a stored record into the AWS
// TopicRuleDestination output structure. The configuration sub-map is echoed
// verbatim so callers see the httpUrlProperties / vpcProperties they supplied.
func topicRuleDestinationResponse(rec map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{
		"arn":          rec["arn"],
		"status":       rec["status"],
		"statusReason": rec["statusReason"],
	}
	if cfg, ok := rec["destinationConfiguration"].(map[string]interface{}); ok {
		if http, ok := cfg["httpUrlConfiguration"].(map[string]interface{}); ok && len(http) > 0 {
			out["httpUrlProperties"] = http
		}
		if vpc, ok := cfg["vpcConfiguration"].(map[string]interface{}); ok && len(vpc) > 0 {
			out["vpcProperties"] = vpc
		}
	}
	if v, ok := rec["createdAt"]; ok {
		out["createdAt"] = v
	}
	if v, ok := rec["lastUpdatedAt"]; ok {
		out["lastUpdatedAt"] = v
	}
	return out
}

func (s *IoTService) CreateTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	cfg := request.GetMapParamCaseInsensitive(req.Parameters, "destinationConfiguration")
	if len(cfg) == 0 {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// AWS assigns a UUID-derived identifier and builds the ARN from it.
	// destinationName is not part of the AWS API.
	uid := uuid.New().String()
	arn := iotstore.BuildTopicRuleDestinationARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), uid)
	confirmationToken := uuid.New().String()
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"arn":                      arn,
		"status":                   "IN_PROGRESS",
		"confirmationToken":        confirmationToken,
		"destinationConfiguration": cfg,
		"createdAt":                now,
		"lastUpdatedAt":            now,
	}
	if err := store.PutGeneric("topicRuleDestination/"+arn, rec); err != nil {
		return nil, err
	}
	// Also index by confirmationToken so ConfirmTopicRuleDestination can look
	// the record up by token alone.
	if err := store.PutGeneric("topicRuleDestinationToken/"+confirmationToken, map[string]interface{}{"arn": arn}); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"topicRuleDestination": topicRuleDestinationResponse(rec),
	}, nil
}
func (s *IoTService) DeleteTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "arn")
	if arn == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "topicRuleDestination/" + arn
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	if err := store.DeleteGeneric(key); err != nil {
		return nil, err
	}
	if token, ok := rec["confirmationToken"].(string); ok && token != "" {
		if err := store.DeleteGeneric("topicRuleDestinationToken/" + token); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) GetTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "arn")
	if arn == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("topicRuleDestination/"+arn, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	return map[string]interface{}{
		"topicRuleDestination": topicRuleDestinationResponse(rec),
	}, nil
}
func (s *IoTService) ListTopicRuleDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("topicRuleDestination/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		summaries = append(summaries, topicRuleDestinationResponse(rec))
	}
	return paginatedMaps("destinationSummaries", summaries, req.Parameters), nil
}
func (s *IoTService) UpdateTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "arn")
	status := request.GetParamCaseInsensitive(req.Parameters, "status")
	if arn == "" || status == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "topicRuleDestination/" + arn
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	rec["status"] = status
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) ConfirmTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	token := request.GetParamCaseInsensitive(req.Parameters, "confirmationToken")
	if token == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Resolve token -> arn, then flip the destination status to ENABLED.
	tokenRec := map[string]interface{}{}
	tokenExists, err := store.GetGenericExists("topicRuleDestinationToken/"+token, &tokenRec)
	if err != nil {
		return nil, err
	}
	if !tokenExists {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	arn, _ := tokenRec["arn"].(string)
	if arn == "" {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	destKey := "topicRuleDestination/" + arn
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(destKey, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	rec["status"] = "ENABLED"
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(destKey, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

// CreateProvisioningClaim issues a temporary provisioning claim bound to an
// existing provisioning template. AWS uses the claim for just-in-time
// provisioning during device manufacturing. The claim consists of a
// short-lived self-signed X.509 certificate and its private key.
func (s *IoTService) CreateProvisioningClaim(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	if templateName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tmpl, err := store.GetProvisioningTemplate(templateName)
	if err != nil {
		return nil, err
	}
	if tmpl == nil || tmpl.TemplateName == "" {
		return nil, iotstore.ErrTemplateNotFound
	}

	// Generate an ECDSA P-256 private key for the claim certificate.
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}

	// Build a self-signed X.509 certificate valid for 1 hour.
	certID := uuid.New().String()
	now := time.Now().UTC()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName:   certID,
			Organization: []string{"vorpalstacks-iot-provisioning-claim"},
		},
		NotBefore:             now.Add(-1 * time.Minute),
		NotAfter:              now.Add(1 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}

	certPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}))
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}
	pubKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyBytes}))
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}
	privKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privKeyBytes}))

	return map[string]interface{}{
		"certificateId":  certID,
		"certificatePem": certPEM,
		"keyPair": map[string]interface{}{
			"PublicKey":  pubKeyPEM,
			"PrivateKey": privKeyPEM,
		},
		"expiration": now.Add(time.Hour).UTC().Unix(),
	}, nil
}

// ---------------------------------------------------------------------------
// Stream operations (MQTT-based file delivery).
// Streams are lightweight metadata records keyed by streamId. The actual
// file payload delivery happens over MQTT; the control-plane API manages
// the stream catalog and versioning only.
// ---------------------------------------------------------------------------

func (s *IoTService) CreateStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamID := request.GetParamCaseInsensitive(req.Parameters, "streamId")
	if streamID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"streamId":      streamID,
		"streamArn":     iotstore.BuildStreamARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), streamID),
		"streamVersion": int64(1),
		"description":   request.GetParamCaseInsensitive(req.Parameters, "description"),
		"files":         req.Parameters["files"],
		"roleArn":       request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"tags":          recTags,
		"createdAt":     now,
		"lastUpdatedAt": now,
	}
	if err := store.PutGeneric("stream/"+streamID, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"streamId":      streamID,
		"streamArn":     rec["streamArn"],
		"description":   rec["description"],
		"streamVersion": rec["streamVersion"],
	}, nil
}

func (s *IoTService) DeleteStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamID := request.GetParamCaseInsensitive(req.Parameters, "streamId")
	if streamID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exists, err := store.GetGenericExists("stream/"+streamID, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrStreamNotFound
	}
	if err := store.DeleteGeneric("stream/" + streamID); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamID := request.GetParamCaseInsensitive(req.Parameters, "streamId")
	if streamID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("stream/"+streamID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrStreamNotFound
	}
	return map[string]interface{}{"streamInfo": rec}, nil
}

func (s *IoTService) ListStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("stream/")
	if err != nil {
		return nil, err
	}
	return paginatedMaps("streams", items, req.Parameters), nil
}

func (s *IoTService) UpdateStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamID := request.GetParamCaseInsensitive(req.Parameters, "streamId")
	if streamID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("stream/"+streamID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrStreamNotFound
	}
	if desc := request.GetParamCaseInsensitive(req.Parameters, "description"); desc != "" {
		rec["description"] = desc
	}
	if files, ok := req.Parameters["files"]; ok {
		rec["files"] = files
	}
	if roleArn := request.GetParamCaseInsensitive(req.Parameters, "roleArn"); roleArn != "" {
		rec["roleArn"] = roleArn
	}
	if v, ok := rec["streamVersion"].(int64); ok {
		rec["streamVersion"] = v + 1
	} else if v, ok := rec["streamVersion"].(float64); ok {
		rec["streamVersion"] = int64(v) + 1
	} else {
		rec["streamVersion"] = int64(2)
	}
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric("stream/"+streamID, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"streamId":      streamID,
		"streamArn":     rec["streamArn"],
		"description":   rec["description"],
		"streamVersion": rec["streamVersion"],
	}, nil
}

// ---------------------------------------------------------------------------
// Registration code operations.
// AWS IoT generates a per-account registration code used to register CA
// certificates. The code is lazily created on first GetRegistrationCode call
// and persists until explicitly deleted.
// ---------------------------------------------------------------------------

func (s *IoTService) GetRegistrationCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/registrationCode", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		code := uuid.New().String()
		rec["registrationCode"] = code
		if err := store.PutGeneric("config/registrationCode", rec); err != nil {
			return nil, err
		}
	}
	return rec, nil
}

func (s *IoTService) DeleteRegistrationCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteGeneric("config/registrationCode"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

// ---------------------------------------------------------------------------
// ProvisioningTemplateVersion operations.
// ---------------------------------------------------------------------------

func (s *IoTService) CreateProvisioningTemplateVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	if templateName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if tmpl, err := store.GetProvisioningTemplate(templateName); err != nil {
		return nil, err
	} else if tmpl == nil || tmpl.TemplateName == "" {
		return nil, iotstore.ErrTemplateNotFound
	}
	// Determine the next version ID by scanning existing versions.
	existing, err := store.ListProvisioningTemplateVersions(templateName, parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}
	maxVersion := 0
	for _, v := range existing {
		var n int
		fmt.Sscanf(v.VersionID, "%d", &n)
		if n > maxVersion {
			maxVersion = n
		}
	}
	versionID := fmt.Sprintf("%d", maxVersion+1)
	v := &iotstore.ProvisioningTemplateVersion{
		VersionID:        versionID,
		TemplateBody:     request.GetParamCaseInsensitive(req.Parameters, "templateBody"),
		IsDefaultVersion: request.GetBoolParam(req.Parameters, "setAsDefault"),
		CreationDate:     time.Now().UTC(),
	}
	if _, err := store.CreateProvisioningTemplateVersion(templateName, v); err != nil {
		return nil, err
	}
	versionIDInt := int32(maxVersion + 1)
	return map[string]interface{}{
		"templateArn":      iotstore.BuildProvisioningTemplateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), templateName),
		"templateName":     templateName,
		"versionId":        versionIDInt,
		"isDefaultVersion": v.IsDefaultVersion,
	}, nil
}

func (s *IoTService) DeleteProvisioningTemplateVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	versionID := request.GetParamCaseInsensitive(req.Parameters, "versionId")
	if templateName == "" || versionID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetProvisioningTemplateVersion(templateName, versionID); err != nil {
		return nil, err
	}
	if err := store.DeleteProvisioningTemplateVersion(templateName, versionID); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeProvisioningTemplateVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	versionID := request.GetParamCaseInsensitive(req.Parameters, "versionId")
	if templateName == "" || versionID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	v, err := store.GetProvisioningTemplateVersion(templateName, versionID)
	if err != nil {
		return nil, err
	}
	var versionIDInt int32
	fmt.Sscanf(v.VersionID, "%d", &versionIDInt)
	return map[string]interface{}{
		"templateName":     templateName,
		"versionId":        versionIDInt,
		"templateBody":     v.TemplateBody,
		"isDefaultVersion": v.IsDefaultVersion,
		"creationDate":     v.CreationDate.UTC().Unix(),
	}, nil
}

// ---------------------------------------------------------------------------
// UpdateThingGroupsForThing: atomically add/remove thing group memberships.
// ---------------------------------------------------------------------------

func (s *IoTService) UpdateThingGroupsForThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	add := request.GetStringList(req.Parameters, "thingGroupsToAdd")
	remove := request.GetStringList(req.Parameters, "thingGroupsToRemove")
	for _, g := range add {
		if err := store.AddThingToThingGroup(thingName, g); err != nil {
			return nil, err
		}
	}
	for _, g := range remove {
		if err := store.RemoveThingFromThingGroup(thingName, g); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{}, nil
}

// ---------------------------------------------------------------------------
// GetThingConnectivityData: returns MQTT connection status. Without a real
// MQTT broker feeding connection state, we return connected=false which is
// honest for a platform that has no active device connections.
// ---------------------------------------------------------------------------

func (s *IoTService) GetThingConnectivityData(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetThing(thingName); err != nil {
		return nil, err
	}

	// Check if any certificate principal attached to this thing is currently
	// connected to any MQTT broker (not just the request-region broker).
	connected := false
	connectedAt := int64(0)
	principals, _ := store.ListPrincipalsForThing(thingName)
	for _, principal := range principals {
		certID := extractCertIDFromPrincipal(principal)
		if certID == "" {
			continue
		}
		for _, brk := range s.brokers {
			if c, ts := brk.IsCertConnected(certID); c {
				connected = true
				connectedAt = ts
				break
			}
		}
		if connected {
			break
		}
	}

	return map[string]interface{}{
		"thingName":        thingName,
		"connected":        connected,
		"timestamp":        time.Now().UTC().Unix(),
		"connectTime":      connectedAt,
		"disconnectReason": "",
	}, nil
}

// ---------------------------------------------------------------------------
// V2 principal/thing listing with richer output.
// ---------------------------------------------------------------------------

func (s *IoTService) ListPrincipalThingsV2(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	principal := request.GetParamCaseInsensitive(req.Parameters, "principal")
	if principal == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	things, err := store.ListThingsForPrincipal(principal)
	if err != nil {
		return nil, err
	}
	objects := make([]map[string]interface{}, 0, len(things))
	for _, t := range things {
		objects = append(objects, map[string]interface{}{
			"thingName":          t,
			"thingPrincipalType": "EXCLUSIVE_THING",
		})
	}
	return paginatedMaps("principalThingObjects", objects, req.Parameters), nil
}

// extractCertIDFromPrincipal extracts the certificate ID from an IoT
// principal ARN (e.g. arn:aws:iot:us-east-1:123:cert/abcdef).
func extractCertIDFromPrincipal(principal string) string {
	idx := strings.LastIndex(principal, "cert/")
	if idx < 0 {
		return ""
	}
	return principal[idx+len("cert/"):]
}

func (s *IoTService) ListThingPrincipalsV2(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	principals, err := store.ListPrincipalsForThing(thingName)
	if err != nil {
		return nil, err
	}
	objects := make([]map[string]interface{}, 0, len(principals))
	for _, p := range principals {
		objects = append(objects, map[string]interface{}{
			"principal":          p,
			"thingPrincipalType": "EXCLUSIVE_THING",
		})
	}
	return paginatedMaps("thingPrincipalObjects", objects, req.Parameters), nil
}
