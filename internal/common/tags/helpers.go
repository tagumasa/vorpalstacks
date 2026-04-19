// Package tags provides AWS tag handling utilities for vorpalstacks.
package tags

import (
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/utils/aws/types"
)

// ParseTags parses tags from request parameters.
func ParseTags(params map[string]interface{}, key string) []types.Tag {
	var result []types.Tag
	if tagsParam, ok := params[key]; ok {
		if tagList, ok := tagsParam.([]interface{}); ok {
			for _, t := range tagList {
				if tagMap, ok := t.(map[string]interface{}); ok {
					result = append(result, types.Tag{
						Key:   request.GetStringParam(tagMap, "Key"),
						Value: request.GetStringParam(tagMap, "Value"),
					})
				}
			}
		} else if tagMap, ok := tagsParam.(map[string]interface{}); ok {
			result = parseTagsFromMap(tagMap)
		}
	}
	if len(result) == 0 {
		if itemsParam, ok := params["Items"]; ok {
			if itemsMap, ok := itemsParam.(map[string]interface{}); ok {
				result = parseTagsFromMap(itemsMap)
			}
		}
	}
	return result
}

func parseTagsFromMap(tagMap map[string]interface{}) []types.Tag {
	var result []types.Tag
	if items, ok := tagMap["Items"].([]interface{}); ok {
		for _, item := range items {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if tagInner, ok := itemMap["Tag"].(map[string]interface{}); ok {
					result = append(result, types.Tag{
						Key:   request.GetStringParam(tagInner, "Key"),
						Value: request.GetStringParam(tagInner, "Value"),
					})
				} else {
					result = append(result, types.Tag{
						Key:   request.GetStringParam(itemMap, "Key"),
						Value: request.GetStringParam(itemMap, "Value"),
					})
				}
			}
		}
	} else if tagInner, ok := tagMap["Tag"].(map[string]interface{}); ok {
		result = append(result, types.Tag{
			Key:   request.GetStringParam(tagInner, "Key"),
			Value: request.GetStringParam(tagInner, "Value"),
		})
	} else {
		for k, v := range tagMap {
			if val, ok := v.(string); ok {
				result = append(result, types.Tag{Key: k, Value: val})
			}
		}
	}
	return result
}

// ParseTagsWithQueryFallback parses tags from request parameters with fallback to query-style format.
func ParseTagsWithQueryFallback(params map[string]interface{}, jsonKey string) []types.Tag {
	tags := ParseTags(params, jsonKey)
	if len(tags) == 0 {
		tags = ParseTags(params, request.LowerFirst(jsonKey))
	}
	if len(tags) == 0 {
		tags = ParseTags(params, strings.ToLower(jsonKey))
	}
	if len(tags) > 0 {
		return tags
	}

	var result []types.Tag
	for i := 1; ; i++ {
		keyKey := jsonKey + ".member." + strconv.Itoa(i) + ".Key"
		valueKey := jsonKey + ".member." + strconv.Itoa(i) + ".Value"

		key := request.GetStringParam(params, keyKey)
		if key == "" {
			keyKeyLower := jsonKey + ".member." + strconv.Itoa(i) + ".key"
			key = request.GetStringParam(params, keyKeyLower)
		}
		if key == "" {
			break
		}

		value := request.GetStringParam(params, valueKey)
		if value == "" {
			valueKeyLower := jsonKey + ".member." + strconv.Itoa(i) + ".value"
			value = request.GetStringParam(params, valueKeyLower)
		}

		result = append(result, types.Tag{Key: key, Value: value})
	}
	return result
}

// ParseTagKeys parses tag keys from request parameters.
func ParseTagKeys(params map[string]interface{}, key string) map[string]bool {
	result := make(map[string]bool)
	if tagKeysParam, ok := params[key]; ok {
		if tagKeys, ok := tagKeysParam.([]interface{}); ok {
			for _, k := range tagKeys {
				if keyStr, ok := k.(string); ok {
					result[keyStr] = true
				}
			}
		}
	}
	return result
}

// ParseTagKeysWithQueryFallback parses tag keys from request parameters with fallback to query-style format.
func ParseTagKeysWithQueryFallback(params map[string]interface{}, jsonKey string) []string {
	keys := ParseTagKeysAsSlice(params, jsonKey)
	if len(keys) > 0 {
		return keys
	}

	var result []string
	for i := 1; ; i++ {
		keyKey := jsonKey + ".member." + strconv.Itoa(i)
		key := request.GetStringParam(params, keyKey)
		if key == "" {
			break
		}
		result = append(result, key)
	}

	if len(result) == 0 {
		if singleKey := request.GetStringParam(params, jsonKey); singleKey != "" {
			for _, k := range strings.Split(singleKey, ",") {
				if trimmed := strings.TrimSpace(k); trimmed != "" {
					result = append(result, trimmed)
				}
			}
		}
	}

	return result
}

// ParseTagKeysAsSlice parses tag keys from request parameters as a slice.
func ParseTagKeysAsSlice(params map[string]interface{}, key string) []string {
	var result []string
	if tagKeysParam, ok := params[key]; ok {
		switch v := tagKeysParam.(type) {
		case []interface{}:
			for _, k := range v {
				if keyStr, ok := k.(string); ok {
					result = append(result, keyStr)
				}
			}
		case []string:
			result = append(result, v...)
		case string:
			if v != "" {
				result = append(result, v)
			}
		}
	}
	return result
}

// ToResponse converts tags to response format.
func ToResponse(tags []types.Tag) []map[string]interface{} {
	if len(tags) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		result[i] = map[string]interface{}{
			"Key":   tag.Key,
			"Value": tag.Value,
		}
	}
	return result
}

// ToResponseWithKeyNames converts tags to response format with custom key names.
func ToResponseWithKeyNames(tags []types.Tag, keyName, valueName string) []map[string]interface{} {
	if len(tags) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		result[i] = map[string]interface{}{
			keyName:   tag.Key,
			valueName: tag.Value,
		}
	}
	return result
}

// Apply merges new tags into existing tags.
func Apply(existing []types.Tag, newTags []types.Tag) []types.Tag {
	tagMap := make(map[string]string)
	for _, tag := range existing {
		tagMap[tag.Key] = tag.Value
	}
	for _, tag := range newTags {
		tagMap[tag.Key] = tag.Value
	}
	result := make([]types.Tag, 0, len(tagMap))
	for key, value := range tagMap {
		result = append(result, types.Tag{Key: key, Value: value})
	}
	return result
}

// Remove removes tags with specified keys from a tag slice.
func Remove(tags []types.Tag, keysToRemove map[string]bool) []types.Tag {
	result := make([]types.Tag, 0)
	for _, tag := range tags {
		if !keysToRemove[tag.Key] {
			result = append(result, tag)
		}
	}
	return result
}

// RemoveByTagKeys removes tags whose keys appear in the provided slice.
func RemoveByTagKeys(t []types.Tag, keys []string) []types.Tag {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return Remove(t, m)
}

// ToMap converts a slice of tags to a map.
func ToMap(tags []types.Tag) map[string]string {
	result := make(map[string]string)
	for _, tag := range tags {
		result[tag.Key] = tag.Value
	}
	return result
}

// MapToTags converts a map to a slice of tags.
func MapToTags(m map[string]string) []types.Tag {
	if m == nil {
		return nil
	}
	result := make([]types.Tag, 0, len(m))
	for k, v := range m {
		result = append(result, types.Tag{Key: k, Value: v})
	}
	return result
}

// MapInterfaceToTags converts a map with interface values to a slice of tags.
func MapInterfaceToTags(m map[string]interface{}) []types.Tag {
	if m == nil {
		return nil
	}
	result := make([]types.Tag, 0, len(m))
	for k, v := range m {
		if vs, ok := v.(string); ok {
			result = append(result, types.Tag{Key: k, Value: vs})
		}
	}
	return result
}

// MapToResponse converts a tag map to response format.
func MapToResponse(m map[string]string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(m))
	for k, v := range m {
		result = append(result, map[string]interface{}{
			"Key":   k,
			"Value": v,
		})
	}
	return result
}

// ParseTagsAsMap parses tags from request parameters and returns them as a map.
func ParseTagsAsMap(params map[string]interface{}, key string) map[string]string {
	tagsIf, ok := params[key]
	if !ok {
		tagsIf, ok = params[request.LowerFirst(key)]
	}
	if !ok {
		tagsIf, ok = params[strings.ToLower(key)]
	}
	if !ok {
		return make(map[string]string)
	}

	if tagsMap, ok := tagsIf.(map[string]interface{}); ok {
		result := make(map[string]string)
		for k, v := range tagsMap {
			if vs, ok := v.(string); ok {
				result[k] = vs
			}
		}
		return result
	}

	return ToMap(ParseTags(params, key))
}

// TagConverter converts a type T to a types.Tag.
type TagConverter[T any] func(T) types.Tag

// ConvertTags converts a slice of type T to a slice of types.Tag using a converter function.
func ConvertTags[T any](tags []T, converter TagConverter[T]) []types.Tag {
	if tags == nil {
		return nil
	}
	result := make([]types.Tag, len(tags))
	for i, tag := range tags {
		result[i] = converter(tag)
	}
	return result
}

// ParseTagsWithKeyNames parses tags from params with custom key and value names.
func ParseTagsWithKeyNames(params map[string]interface{}, listKey, keyName, valueName string) []types.Tag {
	var result []types.Tag
	if tagsParam, ok := params[listKey]; ok {
		if tagList, ok := tagsParam.([]interface{}); ok {
			for _, t := range tagList {
				if tagMap, ok := t.(map[string]interface{}); ok {
					result = append(result, types.Tag{
						Key:   request.GetStringParam(tagMap, keyName),
						Value: request.GetStringParam(tagMap, valueName),
					})
				}
			}
		}
	}
	return result
}

// ParseTagsWithQueryFallbackAndKeyNames parses tags from request parameters with query fallback and custom key names.
func ParseTagsWithQueryFallbackAndKeyNames(params map[string]interface{}, jsonKey, keyName, valueName string) []types.Tag {
	tags := ParseTagsWithKeyNames(params, jsonKey, keyName, valueName)
	if len(tags) > 0 {
		return tags
	}

	var result []types.Tag
	for i := 1; ; i++ {
		keyKey := jsonKey + ".member." + strconv.Itoa(i) + "." + keyName
		valueKey := jsonKey + ".member." + strconv.Itoa(i) + "." + valueName

		key := request.GetStringParam(params, keyKey)
		if key == "" {
			keyKeyLower := jsonKey + ".member." + strconv.Itoa(i) + "." + strings.ToLower(keyName)
			key = request.GetStringParam(params, keyKeyLower)
		}
		if key == "" {
			break
		}

		value := request.GetStringParam(params, valueKey)
		if value == "" {
			valueKeyLower := jsonKey + ".member." + strconv.Itoa(i) + "." + strings.ToLower(valueName)
			value = request.GetStringParam(params, valueKeyLower)
		}

		result = append(result, types.Tag{Key: key, Value: value})
	}
	return result
}

// ParseTagsWithPrefix parses tags from request parameters using a prefix pattern.
func ParseTagsWithPrefix(params map[string]interface{}, prefix string) []types.Tag {
	var result []types.Tag
	for i := 1; ; i++ {
		keyKey := prefix + "." + strconv.Itoa(i) + ".Key"
		valueKey := prefix + "." + strconv.Itoa(i) + ".Value"

		key := request.GetStringParam(params, keyKey)
		if key == "" {
			key = request.GetStringParam(params, strings.ToLower(keyKey))
		}
		if key == "" {
			break
		}

		value := request.GetStringParam(params, valueKey)
		if value == "" {
			value = request.GetStringParam(params, strings.ToLower(valueKey))
		}

		result = append(result, types.Tag{Key: key, Value: value})
	}
	return result
}

// ParseTagKeysWithKeyName parses tag keys from request parameters using a custom key name.
func ParseTagKeysWithKeyName(params map[string]interface{}, listKey, keyName string) []string {
	var result []string
	if tagKeysParam, ok := params[listKey]; ok {
		if tagKeys, ok := tagKeysParam.([]interface{}); ok {
			for _, k := range tagKeys {
				if keyStr, ok := k.(string); ok {
					result = append(result, keyStr)
				} else if tagMap, ok := k.(map[string]interface{}); ok {
					if keyStr := request.GetStringParam(tagMap, keyName); keyStr != "" {
						result = append(result, keyStr)
					}
				}
			}
		}
	}
	return result
}

// ConvertToMapSlice converts a slice of tags to a slice of string maps.
func ConvertToMapSlice(tags []types.Tag) []map[string]string {
	result := make([]map[string]string, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]string{
			"Key":   tag.Key,
			"Value": tag.Value,
		})
	}
	return result
}

// MessageTag represents a tag with a name and value for message routing.
type MessageTag struct {
	Name  string
	Value string
}

// ParseMessageTags parses message tags from request parameters.
func ParseMessageTags(params map[string]interface{}, listKey string) []MessageTag {
	tagsIf, ok := params[listKey]
	if !ok {
		return nil
	}
	tagList, ok := tagsIf.([]interface{})
	if !ok {
		return nil
	}
	return ParseMessageTagsFromList(tagList)
}

// ParseMessageTagsFromList parses message tags from a list of interface values.
func ParseMessageTagsFromList(tagList []interface{}) []MessageTag {
	var result []MessageTag
	for _, t := range tagList {
		tagMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, MessageTag{
			Name:  request.GetStringParam(tagMap, "Name"),
			Value: request.GetStringParam(tagMap, "Value"),
		})
	}
	return result
}

// ParseEcsTags parses ECS-style tags from a list of interface values.
func ParseEcsTags(data []interface{}) []map[string]string {
	if len(data) == 0 {
		return nil
	}
	var result []map[string]string
	for _, item := range data {
		if m, ok := item.(map[string]interface{}); ok {
			tag := make(map[string]string)
			for k, v := range m {
				if str, ok := v.(string); ok {
					tag[k] = str
				}
			}
			result = append(result, tag)
		}
	}
	return result
}

// ConvertMessageTagsToSESv2 converts message tags to SESv2 format.
func ConvertMessageTagsToSESv2(tags []MessageTag) []struct {
	Name  string
	Value string
} {
	if len(tags) == 0 {
		return nil
	}
	result := make([]struct {
		Name  string
		Value string
	}, len(tags))
	for i, t := range tags {
		result[i].Name = t.Name
		result[i].Value = t.Value
	}
	return result
}
