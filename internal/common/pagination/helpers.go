// Package pagination provides AWS pagination utilities for vorpalstacks.
package pagination

import (
	"errors"
	"strconv"

	"vorpalstacks/internal/common/request"
)

// DefaultMaxItems is the default maximum number of items to return in a list operation.
const DefaultMaxItems = 100

// AbsoluteMaxItems is the hard upper limit for any pagination parameter.
const AbsoluteMaxItems = 1000

// ErrInvalidOffsetToken reports a nextToken that is not a usable offset
// into the result set: non-numeric, negative, or at/beyond the set's
// length. Consumers surface it as their operation's documented
// invalid-token error (for the iot list operations,
// InvalidRequestException).
var ErrInvalidOffsetToken = errors.New("invalid offset pagination token")

// ParseOffsetToken parses an integer-offset pagination token against the
// result-set length. Tokens are opaque to clients but internally a plain
// decimal offset, so anything non-numeric, negative, or beyond the result
// set is invalid — silently restarting at the first page or clamping
// would desynchronise a client's walk over the list.
func ParseOffsetToken(token string, total int) (int, error) {
	if token == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(token)
	if err != nil || n < 0 || n >= total {
		return 0, ErrInvalidOffsetToken
	}
	return n, nil
}

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

// ResolveMaxItems applies an operation's Smithy @range bounds to a
// caller-supplied page size with default-on-zero semantics: an unset or
// zero value maps to defaultVal, values within [min,max] pass through,
// and anything else is rejected through errFactory so each service keeps
// its own AWS error shape.
func ResolveMaxItems(value, defaultVal, min, max int, errFactory func(got int) error) (int, error) {
	if value == 0 {
		return defaultVal, nil
	}
	if value < min || value > max {
		return 0, errFactory(value)
	}
	return value, nil
}

// ClampMaxItems normalises a page size without erroring: non-positive
// values map to defaultVal and values above max are clamped to max. Use
// for operations whose AWS behaviour is clamping rather than rejection
// (for example Lambda ListFunctions caps MaxItems instead of rejecting).
func ClampMaxItems(value, defaultVal, max int) int {
	if value <= 0 {
		return defaultVal
	}
	if value > max {
		return max
	}
	return value
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
		found := false
		for i, item := range items {
			if keyExtractor(item) == marker {
				startIdx = i + 1
				found = true
				break
			}
		}
		if !found {
			return SliceResult[T]{
				Items:       []T{},
				NextMarker:  "",
				IsTruncated: false,
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

// PaginateOffsetMaps paginates a slice of maps using integer-offset tokens.
// The offset is read from the "nextToken" parameter (case-insensitive)
// and the page size from "maxResults". The response map uses the given key
// for the page items and "nextToken" for the continuation token. An
// unusable token is reported through ErrInvalidOffsetToken rather than
// being ignored.
func PaginateOffsetMaps(key string, items []map[string]interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	offset, err := ParseOffsetToken(request.GetParamCaseInsensitive(params, "nextToken"), len(items))
	if err != nil {
		return nil, err
	}
	max := len(items)
	if m := request.GetIntParam(params, "maxResults"); m > 0 {
		max = m
	}
	end := offset + max
	if end > len(items) {
		end = len(items)
	}
	page := items[offset:end]
	resp := map[string]interface{}{key: page}
	if end < len(items) {
		resp["nextToken"] = strconv.Itoa(end)
	}
	return resp, nil
}

// PaginateOffsetStrings paginates a string slice using the same offset-based
// scheme as PaginateOffsetMaps.
func PaginateOffsetStrings(key string, items []string, params map[string]interface{}) (map[string]interface{}, error) {
	offset, err := ParseOffsetToken(request.GetParamCaseInsensitive(params, "nextToken"), len(items))
	if err != nil {
		return nil, err
	}
	max := len(items)
	if m := request.GetIntParam(params, "maxResults"); m > 0 {
		max = m
	}
	end := offset + max
	if end > len(items) {
		end = len(items)
	}
	page := items[offset:end]
	resp := map[string]interface{}{key: page}
	if end < len(items) {
		resp["nextToken"] = strconv.Itoa(end)
	}
	return resp, nil
}
