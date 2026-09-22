package cloudwatchlogs

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The vended-logs delivery family: delivery sources, delivery
// destinations, their policies, deliveries, and the configuration
// templates that publish the platform's valid/default values. The
// vocabulary the cores accept is scoped to the substrate the platform
// carries — a delivery source is a CloudWatch Logs log group (the one
// resource whose log events the platform's delivery engine can read), a
// delivery destination is a CloudWatch Logs log group (type CWL,
// delivered through the shared ingestion seam) or an S3 bucket (type S3,
// delivered through the S3 invoker) — and every other source or
// destination form rejects at its Put with the operation's declared
// ValidationException rather than accepting a delivery that could never
// flow.

// errDeliveryValidation is the delivery family's parameter rejection
// identity: every operation of the family declares ValidationException
// (their error lists carry ValidationException, ConflictException,
// ResourceNotFoundException and quota/throttle identities — not the
// InvalidParameterException vocabulary of the log-plane operations).
func errDeliveryValidation(format string, args ...interface{}) error {
	return NewLogsError("ValidationException", fmt.Sprintf(format, args...), 400)
}

// errDeliveryConflict is the family's conflicting-state identity
// (declared ConflictException: "This operation attempted to create a
// resource that already exists" — and the documented you-can't-delete
// guards of the delete operations).
func errDeliveryConflict(message string) error {
	return NewLogsError("ConflictException", message, 400)
}

var (
	deliveryNameRe    = regexp.MustCompile(`^[\w-]*$`)
	deliveryLogTypeRe = regexp.MustCompile(`^[\w]*$`)
)

// The platform's delivery-source record vocabulary. The model leaves
// record fields to the source service's own vocabulary (the
// ConfigurationTemplate's allowedFields); the platform's source form is
// a log group, whose event records carry the event time, the message and
// the originating stream.
var deliveryRecordFieldVocabulary = map[string]bool{
	"time":    true,
	"message": true,
	"stream":  true,
}

// deliveryDefaultRecordFields is the template's defaultDeliveryConfigValues
// recordFields set: every vocabulary field, in the template's order.
var deliveryDefaultRecordFields = []string{"time", "stream", "message"}

// deliveryAllowedFieldDelimiters is the field-delimiter vocabulary the
// platform honours for plain/w3c/raw output (the template's
// allowedFieldDelimiters).
var deliveryAllowedFieldDelimiters = []string{"\t", " ", ",", ";", "|"}

// deliveryDefaultFieldDelimiter is the template's default
// fieldDelimiter.
const deliveryDefaultFieldDelimiter = "\t"

// deliverySuffixPathVariables are the variable sections the suffixPath
// member may carry (the template's allowedSuffixPathFields); the engine
// substitutes {name} occurrences with the batch's values.
var deliverySuffixPathVariables = []string{"accountId", "yyyy", "MM", "dd", "HH"}

// --- HTTP handlers ---

// deliverySourceRequest reads the delivery source members of one request.
func deliverySourceRequest(req *request.ParsedRequest) *PutDeliverySourceInput {
	config := map[string]string{}
	if raw, ok := req.Parameters["deliverySourceConfiguration"].(map[string]interface{}); ok {
		for k, v := range raw {
			if str, ok := v.(string); ok {
				config[k] = str
			}
		}
	}
	return &PutDeliverySourceInput{
		Name:                        request.GetParamLowerFirst(req.Parameters, "Name"),
		ResourceArn:                 request.GetParamLowerFirst(req.Parameters, "ResourceArn"),
		LogType:                     request.GetParamLowerFirst(req.Parameters, "LogType"),
		DeliverySourceConfiguration: config,
		Tags:                        deliveryTagsFromParams(req.Parameters),
	}
}

// deliveryTagsFromParams reads the tags map member the delivery family's
// requests carry. A present-but-empty map stays present: the Tags map's
// own length trait is 1-50, so the Core — not the coercion — decides
// whether the empty set is a member violation.
func deliveryTagsFromParams(params map[string]interface{}) map[string]string {
	raw := request.GetMapParamLowerFirst(params, "Tags")
	if raw == nil {
		return nil
	}
	tags := make(map[string]string, len(raw))
	for k, v := range raw {
		if str, ok := v.(string); ok {
			tags[k] = str
		}
	}
	return tags
}

// validateDeliveryTags enforces the Tags map traits on the delivery
// family's Put operations: 1-50 entries ("Map Entries: Maximum number of
// 50 items", PutDeliveryDestination API reference; the shared Tags
// shape's @length min 1) whose keys and values satisfy the TagKey and
// TagValue traits. The family's declared error list carries
// ValidationException — neither the log-plane InvalidParameterException
// the shared tag validator raises nor TooManyTagsException — so the
// violations reject with the family's own identity.
func validateDeliveryTags(tags map[string]string) error {
	if tags == nil {
		return nil
	}
	if len(tags) == 0 || len(tags) > tagutil.MaxTagsPerResource {
		return errDeliveryValidation("tags must contain between 1 and %d entries", tagutil.MaxTagsPerResource)
	}
	if violation := tagEntryViolation(tags); violation != "" {
		return errDeliveryValidation("%s", violation)
	}
	return nil
}

func (s *LogsService) PutDeliverySource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := deliverySourceRequest(req)
	input.Region = reqCtx.GetRegion()
	source, err := s.putDeliverySourceCore(input)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"deliverySource": formatDeliverySource(source),
	}, nil
}

func (s *LogsService) DescribeDeliverySources(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return describeDeliveryFamilyHandler(reqCtx, req, "deliverySources",
		s.describeDeliverySourcesCore, formatDeliverySource)
}

func (s *LogsService) GetDeliverySource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return getDeliveryFamilyHandler(reqCtx, req, "Name", "deliverySource",
		s.getDeliverySourceCore, formatDeliverySource)
}

func (s *LogsService) DeleteDeliverySource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return deleteDeliveryFamilyHandler(reqCtx, req, "Name", "deliverySourceName",
		s.deleteDeliverySourceCore)
}

// describeDeliveryFamilyHandler serves one Describe envelope of the
// delivery family: parse the paging members, page the core, format each
// item and wrap the response member. Handler plumbing only — parsing,
// the core call and serialisation.
func describeDeliveryFamilyHandler[T any](reqCtx *request.RequestContext, req *request.ParsedRequest,
	member string, page func(region, token string, limit int32) ([]T, string, error),
	format func(T) map[string]interface{}) (interface{}, error) {
	items, nextMarker, err := page(reqCtx.GetRegion(),
		request.GetParamLowerFirst(req.Parameters, "NextToken"), int32(request.GetIntParam(req.Parameters, "Limit")))
	if err != nil {
		return nil, err
	}
	formatted := make([]map[string]interface{}, len(items))
	for i, item := range items {
		formatted[i] = format(item)
	}
	resp := map[string]interface{}{member: formatted}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}
	return resp, nil
}

// describeStoreFamilyHandler serves one paginated Describe envelope
// over a store-resolved core — the envelope shape the export and
// subscription listings share with the delivery family (parse the paging
// members inside the closure, page the core, format, wrap the member).
func describeStoreFamilyHandler[T any](s *LogsService, reqCtx *request.RequestContext, req *request.ParsedRequest,
	member string, page func(store *logsstore.Store, nextToken string, limit int32) ([]T, string, error),
	format func(T) map[string]interface{}) (interface{}, error) {
	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	items, nextMarker, err := page(store,
		request.GetParamLowerFirst(req.Parameters, "NextToken"), int32(request.GetIntParam(req.Parameters, "Limit")))
	if err != nil {
		return nil, err
	}
	formatted := make([]map[string]interface{}, len(items))
	for i, item := range items {
		formatted[i] = format(item)
	}
	resp := map[string]interface{}{member: formatted}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}
	return resp, nil
}

// getDeliveryFamilyHandler serves one Get envelope of the delivery
// family: read and require the identifier member, resolve through the
// core, wrap the single formatted record.
func getDeliveryFamilyHandler[T any](reqCtx *request.RequestContext, req *request.ParsedRequest,
	param, member string, resolve func(region, id string) (T, error),
	format func(T) map[string]interface{}) (interface{}, error) {
	id := request.GetParamLowerFirst(req.Parameters, param)
	if id == "" {
		return nil, errValidationMember(strings.ToLower(param[:1]) + param[1:])
	}
	one, err := resolve(reqCtx.GetRegion(), id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{member: format(one)}, nil
}

// deleteDeliveryFamilyHandler serves one Delete envelope of the delivery
// family: read and require the identifier member, run the core, answer
// the empty response.
func deleteDeliveryFamilyHandler(reqCtx *request.RequestContext, req *request.ParsedRequest,
	param, memberLabel string, del func(region, id string) error) (interface{}, error) {
	id := request.GetParamLowerFirst(req.Parameters, param)
	if id == "" {
		return nil, errValidationMember(strings.ToLower(param[:1]) + param[1:])
	}
	if err := del(reqCtx.GetRegion(), id); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func (s *LogsService) PutDeliveryDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &PutDeliveryDestinationInput{
		Name:                    request.GetParamLowerFirst(req.Parameters, "Name"),
		OutputFormat:            request.GetParamLowerFirst(req.Parameters, "OutputFormat"),
		DestinationResourceArn:  "",
		ExplicitDestinationType: request.GetParamLowerFirst(req.Parameters, "DeliveryDestinationType"),
		Tags:                    deliveryTagsFromParams(req.Parameters),
		Region:                  reqCtx.GetRegion(),
	}
	if config := request.GetMapParamLowerFirst(req.Parameters, "DeliveryDestinationConfiguration"); config != nil {
		input.DestinationResourceArn = destinationParam(config, "DestinationResourceArn")
	}
	destination, err := s.putDeliveryDestinationCore(input)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"deliveryDestination": formatDeliveryDestination(destination),
	}, nil
}

func (s *LogsService) DescribeDeliveryDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return describeDeliveryFamilyHandler(reqCtx, req, "deliveryDestinations",
		s.describeDeliveryDestinationsCore, formatDeliveryDestination)
}

func (s *LogsService) GetDeliveryDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return getDeliveryFamilyHandler(reqCtx, req, "Name", "deliveryDestination",
		s.getDeliveryDestinationCore, formatDeliveryDestination)
}

func (s *LogsService) DeleteDeliveryDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return deleteDeliveryFamilyHandler(reqCtx, req, "Name", "deliveryDestinationName",
		s.deleteDeliveryDestinationCore)
}

func (s *LogsService) PutDeliveryDestinationPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "DeliveryDestinationName")
	policyDocument := request.GetParamLowerFirst(req.Parameters, "DeliveryDestinationPolicy")
	policy, err := s.putDeliveryDestinationPolicyCore(reqCtx.GetRegion(), name, policyDocument)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"policy": map[string]interface{}{"deliveryDestinationPolicy": policy.PolicyDocument},
	}, nil
}

func (s *LogsService) GetDeliveryDestinationPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "DeliveryDestinationName")
	policy, err := s.getDeliveryDestinationPolicyCore(reqCtx.GetRegion(), name)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"policy": map[string]interface{}{"deliveryDestinationPolicy": policy.PolicyDocument},
	}, nil
}

func (s *LogsService) DeleteDeliveryDestinationPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "DeliveryDestinationName")
	if name == "" {
		return nil, errValidationMember("deliveryDestinationName")
	}
	if err := s.deleteDeliveryDestinationPolicyCore(reqCtx.GetRegion(), name); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func (s *LogsService) CreateDelivery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &CreateDeliveryInput{
		DeliverySourceName:     request.GetParamLowerFirst(req.Parameters, "DeliverySourceName"),
		DeliveryDestinationArn: request.GetParamLowerFirst(req.Parameters, "DeliveryDestinationArn"),
		RecordFields:           request.GetStringList(req.Parameters, "RecordFields"),
		FieldDelimiter:         request.GetParamLowerFirst(req.Parameters, "FieldDelimiter"),
		Tags:                   deliveryTagsFromParams(req.Parameters),
		Region:                 reqCtx.GetRegion(),
	}
	input.S3SuffixPath, input.S3HiveCompatiblePath, input.S3ConfigPresent = s3DeliveryConfigurationFromParams(req.Parameters)
	delivery, err := s.createDeliveryCore(input)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"delivery": formatDelivery(delivery),
	}, nil
}

// s3DeliveryConfigurationFromParams reads the suffixPath and
// enableHiveCompatiblePath members out of the request's
// s3DeliveryConfiguration object. The third return reports whether the
// member was present at all: the structure is "valid only when the
// delivery's delivery destination is an S3 bucket", so its presence on a
// non-S3 delivery rejects even when it carries no parameters.
func s3DeliveryConfigurationFromParams(params map[string]interface{}) (suffixPath string, hive, present bool) {
	config := request.GetMapParamLowerFirst(params, "S3DeliveryConfiguration")
	if config == nil {
		return "", false, false
	}
	suffixPath = destinationParam(config, "SuffixPath")
	if raw, ok := config["enableHiveCompatiblePath"].(bool); ok {
		hive = raw
	}
	return suffixPath, hive, true
}

func (s *LogsService) DescribeDeliveries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return describeDeliveryFamilyHandler(reqCtx, req, "deliveries",
		s.describeDeliveriesCore, formatDelivery)
}

func (s *LogsService) GetDelivery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return getDeliveryFamilyHandler(reqCtx, req, "Id", "delivery",
		s.getDeliveryCore, formatDelivery)
}

func (s *LogsService) DeleteDelivery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return deleteDeliveryFamilyHandler(reqCtx, req, "Id", "id", s.deleteDeliveryCore)
}

func (s *LogsService) UpdateDeliveryConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &UpdateDeliveryConfigurationInput{
		Id:             request.GetParamLowerFirst(req.Parameters, "Id"),
		RecordFields:   request.GetStringList(req.Parameters, "RecordFields"),
		FieldDelimiter: request.GetParamLowerFirst(req.Parameters, "FieldDelimiter"),
		Region:         reqCtx.GetRegion(),
	}
	input.S3SuffixPath, input.S3HiveCompatiblePath, input.S3ConfigPresent = s3DeliveryConfigurationFromParams(req.Parameters)
	if err := s.updateDeliveryConfigurationCore(input); err != nil {
		return nil, err
	}
	// The operation's response shape carries no members.
	return response.EmptyResponse(), nil
}
