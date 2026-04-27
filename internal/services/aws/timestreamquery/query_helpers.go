package timestreamquery

import (
	"fmt"
)

func (s *TimestreamQueryService) formatRowsForResponse(rows []map[string]interface{}, columnInfo []ColumnInfo) []map[string]interface{} {
	var result []map[string]interface{}
	for _, row := range rows {
		var data []map[string]interface{}
		for _, col := range columnInfo {
			val, exists := row[col.Name]
			if !exists || val == nil {
				data = append(data, map[string]interface{}{
					"NullValue": true,
				})
			} else {
				data = append(data, map[string]interface{}{
					"ScalarValue": fmt.Sprintf("%v", val),
				})
			}
		}
		result = append(result, map[string]interface{}{
			"Data": data,
		})
	}
	return result
}

func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
