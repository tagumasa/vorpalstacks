package neptune

import (
	"vorpalstacks/internal/common/pagination"
)

func paginateItems(items []interface{}, marker string, maxRecords int, keyFn func(interface{}) string) ([]interface{}, string, bool) {
	if maxRecords <= 0 {
		maxRecords = 100
	}
	if maxRecords > 100 {
		maxRecords = 100
	}
	result := pagination.PaginateSlice(items, marker, maxRecords, keyFn)
	return result.Items, result.NextMarker, result.IsTruncated
}
