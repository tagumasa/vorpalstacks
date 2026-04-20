// Package pagination provides AWS pagination utilities for vorpalstacks.
package pagination

import (
	"vorpalstacks/internal/common/request"
)

// DefaultMaxItems is the default maximum number of items to return in a list operation.
const DefaultMaxItems = 100

// AbsoluteMaxItems is the hard upper limit for any pagination parameter.
const AbsoluteMaxItems = 1000

// GetMaxItems extracts a pagination limit parameter from the given params map.
// It checks paramName first, then falls back to "MaxItems" for backward compatibility.
func GetMaxItems(params map[string]interface{}, defaultVal int, paramName ...string) int {
	if defaultVal <= 0 {
		defaultVal = DefaultMaxItems
	}
	if len(paramName) > 0 {
		maxItems := request.GetIntParam(params, paramName[0])
		if maxItems <= 0 {
			maxItems = request.GetIntParam(params, "MaxItems")
		}
		if maxItems > 0 {
			if maxItems > AbsoluteMaxItems {
				return AbsoluteMaxItems
			}
			return maxItems
		}
	}
	maxItems := request.GetIntParam(params, "MaxItems")
	if maxItems <= 0 {
		return defaultVal
	}
	if maxItems > AbsoluteMaxItems {
		return AbsoluteMaxItems
	}
	return maxItems
}

// GetMarker extracts a pagination marker/nextToken parameter from the given params map.
// It checks paramName first, then falls back to "Marker" for backward compatibility.
func GetMarker(params map[string]interface{}, paramName ...string) string {
	if len(paramName) > 0 {
		if m := request.GetStringParam(params, paramName[0]); m != "" {
			return m
		}
	}
	return request.GetStringParam(params, "Marker")
}

// SetNextToken sets a pagination token in the response map under the given key.
func SetNextToken(response map[string]interface{}, key string, value string) {
	if value != "" {
		response[key] = value
	}
}

// BuildListResponse builds a list response map with the given items and next token.
func BuildListResponse(itemsKey string, items interface{}, nextToken string) map[string]interface{} {
	response := map[string]interface{}{
		itemsKey: items,
	}
	if nextToken != "" {
		response["NextToken"] = nextToken
	}
	return response
}

// SliceResult represents the result of paginating a slice.
type SliceResult[T any] struct {
	Items       []T
	NextMarker  string
	IsTruncated bool
}

// KeyExtractor is a function that extracts a marker key from an item.
type KeyExtractor[T any] func(item T) string

// PaginateSlice paginates a slice based on marker and maxItems.
// The keyExtractor function is used to find the starting position and generate the next marker.
func PaginateSlice[T any](items []T, marker string, maxItems int, keyExtractor KeyExtractor[T]) SliceResult[T] {
	if len(items) == 0 {
		return SliceResult[T]{
			Items:       []T{},
			NextMarker:  "",
			IsTruncated: false,
		}
	}

	startIdx := 0
	if marker != "" {
		for i, item := range items {
			if keyExtractor(item) == marker {
				startIdx = i + 1
				break
			}
		}
	}

	endIdx := startIdx + maxItems
	if endIdx > len(items) {
		endIdx = len(items)
	}

	var resultItems []T
	if startIdx < len(items) {
		resultItems = items[startIdx:endIdx]
	} else {
		resultItems = []T{}
	}

	isTruncated := endIdx < len(items)
	var nextMarker string
	if isTruncated && len(resultItems) > 0 {
		nextMarker = keyExtractor(resultItems[len(resultItems)-1])
	}

	return SliceResult[T]{
		Items:       resultItems,
		NextMarker:  nextMarker,
		IsTruncated: isTruncated,
	}
}
