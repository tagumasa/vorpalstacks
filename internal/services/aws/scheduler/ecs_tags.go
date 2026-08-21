package scheduler

// parseEcsTags parses ECS-style tags from a list of interface values
// (EcsParameters.Tags).
func parseEcsTags(data []interface{}) []map[string]string {
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
