package naming

// Package naming provides naming utility functions for vorpalstacks.

import (
	"encoding/json"
)

// ExtractTraitString extracts a string value from a trait map
func ExtractTraitString(traits map[string]interface{}, traitName string) string {
	if traits == nil {
		return ""
	}
	if val, ok := traits[traitName]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// ExtractTraitBool extracts a boolean value from a trait map
func ExtractTraitBool(traits map[string]interface{}, traitName string) bool {
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

// ExtractTraitInt extracts an integer value from a trait map
func ExtractTraitInt(traits map[string]interface{}, traitName string) int {
	if traits == nil {
		return 0
	}
	if val, ok := traits[traitName]; ok {
		if f, ok := val.(float64); ok {
			return int(f)
		}
	}
	return 0
}

// ExtractTraitJSON extracts a value from a trait map and returns it as JSON string
func ExtractTraitJSON(traits map[string]interface{}, traitName string) string {
	if traits == nil {
		return ""
	}
	if val, ok := traits[traitName]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		// Convert complex types to JSON
		if bytes, err := json.Marshal(val); err == nil {
			return string(bytes)
		}
	}
	return ""
}

// ExtractTraitMap extracts a map value from a trait map
func ExtractTraitMap(traits map[string]interface{}, traitName string) map[string]interface{} {
	if traits == nil {
		return nil
	}
	if val, ok := traits[traitName]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			return m
		}
	}
	return nil
}

// ExtractTraitArray extracts an array value from a trait map
func ExtractTraitArray(traits map[string]interface{}, traitName string) []interface{} {
	if traits == nil {
		return nil
	}
	if val, ok := traits[traitName]; ok {
		if arr, ok := val.([]interface{}); ok {
			return arr
		}
	}
	return nil
}
