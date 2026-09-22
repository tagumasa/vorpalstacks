package cloudwatchlogs

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The configuration-template listing: the platform's published valid
// and default delivery values.

// --- Configuration templates ---

// deliveryConfigurationTemplates is the platform's template catalog: one
// row per (source form, destination type) the platform carries. The
// model leaves the catalog to the service; the platform's source form is
// a CloudWatch Logs log group, and the logType cell stays empty because
// the platform does not constrain the source's log type vocabulary.
var deliveryConfigurationTemplates = []map[string]interface{}{
	{
		"service":                 "logs",
		"resourceType":            "logGroup",
		"deliveryDestinationType": "CWL",
		"defaultDeliveryConfigValues": map[string]interface{}{
			"recordFields": deliveryDefaultRecordFields,
		},
		"allowedFields":        deliveryAllowedFields(),
		"allowedOutputFormats": []string{},
	},
	{
		"service":                 "logs",
		"resourceType":            "logGroup",
		"deliveryDestinationType": "S3",
		"defaultDeliveryConfigValues": map[string]interface{}{
			"recordFields":   deliveryDefaultRecordFields,
			"fieldDelimiter": deliveryDefaultFieldDelimiter,
			"s3DeliveryConfiguration": map[string]interface{}{
				"enableHiveCompatiblePath": false,
			},
		},
		"allowedFields":           deliveryAllowedFields(),
		"allowedOutputFormats":    []string{"json", "plain", "w3c", "raw"},
		"allowedFieldDelimiters":  deliveryAllowedFieldDelimiters,
		"allowedSuffixPathFields": deliverySuffixPathVariables,
	},
}

// deliveryAllowedFields renders the template's AllowedFields list: the
// model types each element as the RecordField structure (a field header
// plus whether it is mandatory), and the platform's source vocabulary
// carries no mandatory field — a delivery may narrow recordFields to any
// subset, including the message alone.
func deliveryAllowedFields() []map[string]interface{} {
	fields := make([]map[string]interface{}, 0, len(deliveryDefaultRecordFields))
	for _, f := range deliveryDefaultRecordFields {
		fields = append(fields, map[string]interface{}{"name": f, "mandatory": false})
	}
	return fields
}

func (s *LogsService) DescribeConfigurationTemplates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	items, nextMarker, err := s.describeConfigurationTemplatesCore(&DescribeConfigurationTemplatesInput{
		Service:                  request.GetParamLowerFirst(req.Parameters, "Service"),
		LogTypes:                 request.GetStringList(req.Parameters, "LogTypes"),
		ResourceTypes:            request.GetStringList(req.Parameters, "ResourceTypes"),
		DeliveryDestinationTypes: request.GetStringList(req.Parameters, "DeliveryDestinationTypes"),
		NextToken:                request.GetParamLowerFirst(req.Parameters, "NextToken"),
		Limit:                    int32(request.GetIntParam(req.Parameters, "Limit")),
	})
	if err != nil {
		return nil, err
	}
	resp := map[string]interface{}{"configurationTemplates": items}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}
	return resp, nil
}

// DescribeConfigurationTemplatesInput is the service-layer DTO for the
// configuration-template listing (both planes build it; the core owns
// the filter and bound validation).
type DescribeConfigurationTemplatesInput struct {
	Service                  string
	LogTypes                 []string
	ResourceTypes            []string
	DeliveryDestinationTypes []string
	NextToken                string
	Limit                    int32
}

// describeConfigurationTemplatesCore validates the filter members and
// returns one page of the matching template catalog rows.
func (s *LogsService) describeConfigurationTemplatesCore(input *DescribeConfigurationTemplatesInput) ([]map[string]interface{}, string, error) {
	if len(input.LogTypes) > logsstore.ConfigurationTemplatesLogTypesMax || len(input.ResourceTypes) > logsstore.ConfigurationTemplatesLogTypesMax {
		return nil, "", errDeliveryValidation("logTypes and resourceTypes accept at most %d entries", logsstore.ConfigurationTemplatesLogTypesMax)
	}
	if len(input.DeliveryDestinationTypes) > logsstore.ConfigurationTemplatesDestinationTypesMax {
		return nil, "", errDeliveryValidation("deliveryDestinationTypes accepts at most %d entries", logsstore.ConfigurationTemplatesDestinationTypesMax)
	}
	limit, err := validateListLimitValidation(input.Limit, logsstore.DefaultDescribeLimit, logsstore.MaxDescribeLimit)
	if err != nil {
		return nil, "", err
	}

	matchesAny := func(filter []string, rowValue string) bool {
		if len(filter) == 0 {
			return true
		}
		for _, f := range filter {
			if f == rowValue {
				return true
			}
		}
		return false
	}

	var matched []map[string]interface{}
	for _, tmpl := range deliveryConfigurationTemplates {
		if input.Service != "" && tmpl["service"] != input.Service {
			continue
		}
		// The platform's templates carry no logType cell: they apply to
		// every log type, so a logTypes filter never excludes a row.
		if !matchesAny(input.ResourceTypes, tmpl["resourceType"].(string)) {
			continue
		}
		if !matchesAny(input.DeliveryDestinationTypes, tmpl["deliveryDestinationType"].(string)) {
			continue
		}
		matched = append(matched, tmpl)
	}

	// The page marker rides the scoped listing vocabulary: the token
	// carries the filter identity and its documented expiry rather than
	// a bare catalog key.
	scope := listingScope("configurationtemplates", input.Service,
		strings.Join(input.LogTypes, "\x1f"), strings.Join(input.ResourceTypes, "\x1f"),
		strings.Join(input.DeliveryDestinationTypes, "\x1f"))
	result, err := paginateScopedListing(scope, input.NextToken, matched, int(limit), func(t map[string]interface{}) string {
		return fmt.Sprintf("%v", t["deliveryDestinationType"])
	})
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}
