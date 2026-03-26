// Package pagination provides AWS pagination utilities for vorpalstacks.
package pagination

import (
	"vorpalstacks/internal/services/aws/common/request"
)

// DefaultMaxItems is the default maximum number of items to return in a list operation.
const DefaultMaxItems = 100

// GetMaxItems extracts the MaxItems parameter from the given params map.
func GetMaxItems(params map[string]interface{}, defaultVal int) int {
	if defaultVal <= 0 {
		defaultVal = DefaultMaxItems
	}
	maxItems := request.GetIntParam(params, "MaxItems")
	if maxItems <= 0 {
		return defaultVal
	}
	return maxItems
}

// GetMarker extracts the Marker parameter from the given params map.
func GetMarker(params map[string]interface{}) string {
	return request.GetStringParam(params, "Marker")
}

// ListResult represents the result of a list operation with pagination support.
type ListResult struct {
	Items       interface{}
	IsTruncated bool
	NextMarker  string
}

// BuildNextMarker builds the next marker for pagination based on the count and last key.
func BuildNextMarker(count, maxItems int, lastKey string) (string, bool) {
	if count >= maxItems && lastKey != "" {
		return lastKey, true
	}
	return "", false
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
