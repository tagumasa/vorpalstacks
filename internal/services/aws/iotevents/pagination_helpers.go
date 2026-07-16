package iotevents

import (
	pagination "vorpalstacks/internal/common/pagination"
)

// paginatedMaps applies offset-based pagination to a slice of maps. The
// nextToken parameter is treated as an integer offset; maxResults limits the
// page size. When maxResults is zero the entire slice is returned (no
// pagination), preserving backward compatibility for callers that do not
// request pagination.
func paginatedMaps(key string, items []map[string]interface{}, params map[string]interface{}) map[string]interface{} {
	return pagination.PaginateOffsetMaps(key, items, params)
}
