// Package smithy provides Smithy model handling for vorpalstacks.
package smithy

import "encoding/json"

// ExtractStringTrait extracts a string value from traits.
func ExtractStringTrait(traits map[string]interface{}, key string) string {
	if trait, ok := traits[key]; ok {
		if str, ok := trait.(string); ok {
			return str
		}
		if m, ok := trait.(map[string]interface{}); ok {
			if v, ok := m["value"].(string); ok {
				return v
			}
		}
	}
	return ""
}

// ExtractIntTrait extracts an integer value from traits.
func ExtractIntTrait(traits map[string]interface{}, key string) int {
	if trait, ok := traits[key]; ok {
		if i, ok := trait.(int); ok {
			return i
		}
		if f, ok := trait.(float64); ok {
			return int(f)
		}
		if m, ok := trait.(map[string]interface{}); ok {
			if v, ok := m["value"].(float64); ok {
				return int(v)
			}
			if v, ok := m["value"].(int); ok {
				return v
			}
		}
	}
	return 0
}

// ExtractBoolTrait extracts a boolean value from traits.
func ExtractBoolTrait(traits map[string]interface{}, key string) bool {
	if trait, ok := traits[key]; ok {
		if b, ok := trait.(bool); ok {
			return b
		}
		if _, ok := trait.(map[string]interface{}); ok {
			return true
		}
	}
	return false
}

// ExtractDocumentation extracts documentation from traits.
func ExtractDocumentation(traits map[string]interface{}) string {
	if docTrait, ok := traits[TraitDocs]; ok {
		if str, ok := docTrait.(string); ok {
			return str
		}
		if m, ok := docTrait.(map[string]interface{}); ok {
			if value, ok := m["value"].(string); ok {
				return value
			}
		}
	}
	if docTrait, ok := traits[TraitDocumentation]; ok {
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

// ExtractErrorFault extracts the error fault type from traits.
func ExtractErrorFault(traits map[string]interface{}) string {
	return ExtractStringTrait(traits, TraitSmithyAPIError)
}

// ExtractHTTPStatusCode extracts the HTTP status code from traits.
func ExtractHTTPStatusCode(traits map[string]interface{}) int {
	if code := ExtractIntTrait(traits, TraitSmithyAPIHttpError); code != 0 {
		return code
	}
	return ExtractIntTrait(traits, TraitSmithyAPIHttpCode)
}

// ExtractErrorCodeOverride extracts the error code override from traits.
func ExtractErrorCodeOverride(traits map[string]interface{}) string {
	return ExtractStringTrait(traits, TraitAwsAPIErrorCode)
}

// ExtractServiceInfo extracts service information from traits.
func ExtractServiceInfo(traits map[string]interface{}) (sdkID, arnNamespace, cloudFormationName, cloudTrailEventSource, endpointPrefix string) {
	if serviceTrait, ok := traits[TraitAwsAPIService]; ok {
		if traitMap, ok := serviceTrait.(map[string]interface{}); ok {
			sdkID = ExtractStringTrait(traitMap, "sdkId")
			arnNamespace = ExtractStringTrait(traitMap, "arnNamespace")
			cloudFormationName = ExtractStringTrait(traitMap, "cloudFormationName")
			cloudTrailEventSource = ExtractStringTrait(traitMap, "cloudTrailEventSource")
			endpointPrefix = ExtractStringTrait(traitMap, "endpointPrefix")
		}
	}
	return
}

// ExtractServiceVersion extracts the service version from traits.
func ExtractServiceVersion(traits map[string]interface{}) string {
	if versionTrait, ok := traits[TraitAwsAPIVersion]; ok {
		if v, ok := versionTrait.(string); ok {
			return v
		}
		if m, ok := versionTrait.(map[string]interface{}); ok {
			return ExtractStringTrait(m, "value")
		}
	}
	return ""
}

// ExtractEnumValue extracts the enum value from traits.
func ExtractEnumValue(traits map[string]interface{}) string {
	return ExtractStringTrait(traits, TraitSmithyAPIEnumValue)
}

// FormatTraitValue formats a trait value to a string.
func FormatTraitValue(traitValue interface{}) string {
	if str, ok := traitValue.(string); ok {
		return str
	}
	if bytes, err := json.Marshal(traitValue); err == nil {
		return string(bytes)
	}
	return ""
}
