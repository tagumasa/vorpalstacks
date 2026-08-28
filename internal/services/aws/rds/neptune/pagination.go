package neptune

import (
	"vorpalstacks/internal/common/pagination"
	rdssvc "vorpalstacks/internal/services/aws/rds"
)

func paginateItems(items []interface{}, marker string, maxRecords int, keyFn func(interface{}) string) ([]interface{}, string, bool) {
	resolved := rdssvc.ResolveDescribeMaxRecords(maxRecords)
	result := pagination.PaginateSlice(items, marker, resolved, keyFn)
	return result.Items, result.NextMarker, result.IsTruncated
}
