package scheduler

import (
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
)

// parseEcsTags parses EcsParameters.Tags — a list of TagMap entries, each a
// map of string to string. The parse is strict: a non-map entry or a
// non-string pair value is a wire-format violation reported to the caller,
// never a value silently dropped where the Core validator cannot see it.
func parseEcsTags(data []interface{}) ([]map[string]string, error) {
	if len(data) == 0 {
		return nil, nil
	}
	result := make([]map[string]string, 0, len(data))
	for i, item := range data {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.Tags[%d] must be a map of string to string", i))
		}
		tag := make(map[string]string, len(m))
		for k, v := range m {
			str, ok := v.(string)
			if !ok {
				return nil, awserrors.NewValidationException(fmt.Sprintf(
					"EcsParameters.Tags[%d].%s must be a string", i, k))
			}
			tag[k] = str
		}
		result = append(result, tag)
	}
	return result, nil
}
