// Package converter provides types and utilities for parsing Smithy model files.
package converter

import (
	"encoding/json"
	"strings"
)

// ShapeSkipTraits defines trait names that should be skipped when storing shape traits.
// These traits are used for processing but not stored as separate records.
var ShapeSkipTraits = map[string]bool{
	"docs":                 true,
	"documentation":        true,
	"smithy.api#input":     true,
	"smithy.api#output":    true,
	"smithy.api#errors":    true,
	"smithy.api#http":      true,
	"smithy.api#readonly":  true,
	"smithy.api#error":     true,
	"smithy.api#httpError": true,
	"smithy.api#httpCode":  true,
	"aws.api#errorCode":    true,
	"aws.api#service":      true,
	"aws.api#version":      true,
	"smithy.api#enumValue": true,
	"smithy.api#enum":      true,
}

// MemberSkipTraits defines trait names that should be skipped when storing member traits.
// These traits are used for processing but not stored as separate records.
var MemberSkipTraits = map[string]bool{
	"docs":                                true,
	"documentation":                       true,
	"smithy.api#httpLabel":                true,
	"smithy.api#httpQuery":                true,
	"smithy.api#httpHeader":               true,
	"smithy.api#httpPayload":              true,
	"smithy.api#jsonName":                 true,
	"smithy.api#required":                 true,
	"smithy.api#httpPrefix":               true,
	"smithy.api#httpLocationName":         true,
	"smithy.http#defaultPayloadMediaType": true,
	"smithy.api#enumValue":                true,
}

// extractStringTrait extracts a string value from a trait map.
// Traits can be either direct string values or map structures with a "value" key.
func extractStringTrait(traits map[string]interface{}, traitName string) string {
	if traits == nil {
		return ""
	}

	if val, ok := traits[traitName]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		if m, ok := val.(map[string]interface{}); ok {
			if value, ok := m["value"].(string); ok {
				return value
			}
		}
	}

	return ""
}

// extractIntTrait extracts an integer value from a trait map.
// Traits can be either direct numeric values or map structures with a "value" key.
func extractIntTrait(traits map[string]interface{}, traitName string) int {
	if traits == nil {
		return 0
	}

	if val, ok := traits[traitName]; ok {
		if f, ok := val.(float64); ok {
			return int(f)
		}
		if i, ok := val.(int); ok {
			return i
		}
	}

	return 0
}

// extractBoolTrait extracts a boolean value from a trait map.
// Traits can be either direct boolean values or map structures.
func extractBoolTrait(traits map[string]interface{}, traitName string) bool {
	if traits == nil {
		return false
	}

	if val, ok := traits[traitName]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}

	return false
}

// extractHTTPBinding extracts HTTP binding information from a trait map.
// It checks multiple trait keys in order and returns the first valid HTTP binding found.
func extractHTTPBinding(traits map[string]interface{}, traitKeys []string) *HTTPBinding {
	if traits == nil {
		return nil
	}

	for _, key := range traitKeys {
		if httpTrait, ok := traits[key]; ok {
			binding := &HTTPBinding{}

			httpMap, ok := httpTrait.(map[string]interface{})
			if !ok {
				continue
			}

			if method, ok := httpMap["method"].(string); ok {
				binding.Method = method
			}
			if uri, ok := httpMap["uri"].(string); ok {
				binding.URI = uri
			}
			if code, ok := httpMap["code"].(float64); ok {
				binding.Code = int(code)
			} else if code, ok := httpMap["code"].(int); ok {
				binding.Code = code
			}

			if binding.Method != "" || binding.URI != "" {
				return binding
			}
		}
	}

	return nil
}

// extractShapeName extracts the shape name from a shape ID.
// Shape IDs can be in formats like "namespace#ShapeName" or "com.example#ShapeName".
// This function returns just the name portion without the namespace.
func extractShapeName(shapeID string) string {
	if shapeID == "" {
		return ""
	}

	parts := strings.Split(shapeID, "#")
	if len(parts) > 1 {
		return parts[1]
	}

	namespaceParts := strings.Split(parts[0], ".")
	if len(namespaceParts) > 0 {
		return namespaceParts[len(namespaceParts)-1]
	}

	return ""
}

// formatTraitValue converts a trait value to a string representation for storage.
// Complex trait values are JSON-encoded, while simple string values are used directly.
func formatTraitValue(traitValue interface{}) string {
	if traitValue == nil {
		return ""
	}
	if str, ok := traitValue.(string); ok {
		return str
	}
	if bytes, err := json.Marshal(traitValue); err == nil {
		return string(bytes)
	}
	return ""
}

// extractEnumValue extracts the enum value trait from a trait map.
// This trait specifies the actual string value for enum members.
func extractEnumValue(traits map[string]interface{}) string {
	return extractStringTrait(traits, "smithy.api#enumValue")
}

// extractErrorFault extracts the error fault type from a trait map.
// Error faults indicate whether an error is a client fault or server fault.
func extractErrorFault(traits map[string]interface{}) string {
	return extractStringTrait(traits, "smithy.api#error")
}

// extractHTTPStatusCode extracts the HTTP status code for error responses.
// Status codes can be specified via httpError or httpCode traits.
func extractHTTPStatusCode(traits map[string]interface{}) int {
	if code := extractIntTrait(traits, "smithy.api#httpError"); code != 0 {
		return code
	}
	return extractIntTrait(traits, "smithy.api#httpCode")
}

// extractErrorCodeOverride extracts the AWS error code override from a trait map.
// This allows specifying a custom error code different from the default.
func extractErrorCodeOverride(traits map[string]interface{}) string {
	return extractStringTrait(traits, "aws.api#errorCode")
}

// ShouldStoreShapeTrait returns true if a shape trait should be stored.
// Some traits are skipped because they are used for processing only.
func ShouldStoreShapeTrait(traitName string) bool {
	return !ShapeSkipTraits[traitName]
}

// ShouldStoreMemberTrait returns true if a member trait should be stored.
// Some traits are skipped because they are used for processing only.
func ShouldStoreMemberTrait(traitName string) bool {
	return !MemberSkipTraits[traitName]
}

// ParseErrorsTrait parses the errors trait and returns a list of target shape IDs.
// The errors trait can be in various formats depending on the Smithy version.
func ParseErrorsTrait(errorsTrait interface{}) []string {
	if errorsTrait == nil {
		return nil
	}

	var targets []string

	switch v := errorsTrait.(type) {
	case []interface{}:
		for _, e := range v {
			if m, ok := e.(map[string]interface{}); ok {
				if t, ok := m["target"].(string); ok {
					targets = append(targets, t)
				}
			} else if s, ok := e.(string); ok {
				targets = append(targets, s)
			}
		}
	case []map[string]interface{}:
		for _, e := range v {
			if t, ok := e["target"].(string); ok {
				targets = append(targets, t)
			}
		}
	case []string:
		targets = v
	case map[string]interface{}:
		if t, ok := v["target"].(string); ok {
			targets = append(targets, t)
		}
	case string:
		targets = append(targets, v)
	}

	return targets
}

// extractServiceInfo extracts AWS service metadata from a trait map.
// This includes the SDK ID, ARN namespace, CloudFormation name, CloudTrail event source, and endpoint prefix.
func extractServiceInfo(traits map[string]interface{}) (sdkID, arnNamespace, cloudFormationName, cloudTrailEventSource, endpointPrefix string) {
	if traits == nil {
		return
	}

	serviceTrait, ok := traits["aws.api#service"]
	if !ok {
		return
	}

	traitMap, ok := serviceTrait.(map[string]interface{})
	if !ok {
		return
	}

	sdkID = extractStringTrait(traitMap, "sdkId")
	arnNamespace = extractStringTrait(traitMap, "arnNamespace")
	cloudFormationName = extractStringTrait(traitMap, "cloudFormationName")
	cloudTrailEventSource = extractStringTrait(traitMap, "cloudTrailEventSource")
	endpointPrefix = extractStringTrait(traitMap, "endpointPrefix")

	return
}

// extractServiceVersion extracts the service version from a trait map.
// The version is specified in the aws.api#version trait.
func extractServiceVersion(traits map[string]interface{}) string {
	if traits == nil {
		return ""
	}

	versionTrait, ok := traits["aws.api#version"]
	if !ok {
		return ""
	}

	if v, ok := versionTrait.(string); ok {
		return v
	}

	if m, ok := versionTrait.(map[string]interface{}); ok {
		return extractStringTrait(m, "value")
	}

	return ""
}

// extractDocumentation extracts documentation text from a trait map.
// Documentation can be provided via either the "docs" or "documentation" trait.
func extractDocumentation(traits map[string]interface{}) string {
	if traits == nil {
		return ""
	}

	if docTrait, ok := traits["docs"]; ok {
		if str, ok := docTrait.(string); ok {
			return str
		}
		if m, ok := docTrait.(map[string]interface{}); ok {
			if value, ok := m["value"].(string); ok {
				return value
			}
		}
	}

	if docTrait, ok := traits["documentation"]; ok {
		if str, ok := docTrait.(string); ok {
			return str
		}
		if m, ok := docTrait.(map[string]interface{}); ok {
			if value, ok := m["value"].(string); ok {
				return value
			}
		}
	}

	return ""
}
