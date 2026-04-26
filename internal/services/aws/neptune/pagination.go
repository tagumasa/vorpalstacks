package neptune

func paginateItems(items []interface{}, marker string, maxRecords int, keyFn func(interface{}) string) (resultItems []interface{}, nextMarker string, isTruncated bool) {
	if maxRecords <= 0 {
		maxRecords = 100
	}
	if maxRecords > 100 {
		maxRecords = 100
	}

	if len(items) == 0 {
		return items, "", false
	}

	startIdx := 0
	if marker != "" {
		found := false
		for i, item := range items {
			if keyFn(item) == marker {
				startIdx = i + 1
				found = true
				break
			}
		}
		if !found {
			return []interface{}{}, "", false
		}
	}

	endIdx := startIdx + maxRecords
	if endIdx > len(items) {
		endIdx = len(items)
	}

	if startIdx >= len(items) {
		return []interface{}{}, "", false
	}

	resultItems = items[startIdx:endIdx]
	isTruncated = endIdx < len(items)
	if isTruncated && len(resultItems) > 0 {
		nextMarker = keyFn(resultItems[len(resultItems)-1])
	}

	return resultItems, nextMarker, isTruncated
}
