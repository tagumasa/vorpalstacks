package cloudwatchlogs

import (
	"encoding/json"
	"sort"
	"time"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The field-discovery cores: GetLogGroupFields' sampled field inventory
// over a group's recent events, and GetLogFields' data-source field
// listing.

func (s *LogsService) getLogGroupFieldsCore(store *logsstore.Store, logGroupName string, timeSeconds int64) ([]map[string]interface{}, error) {
	if _, err := store.GetLogGroup(logGroupName); err != nil {
		return nil, mapStoreError(err)
	}

	standardFields := []string{"@timestamp", "@message", "@logStream", "@log", "@ingestionTime"}
	// jsonFieldCounts carries each discovered field's occurrence count
	// over the sampled events — the member's documented meaning is
	// "The percentage of log events queried that contained the field",
	// a sampled ratio rather than a constant.
	jsonFieldCounts := make(map[string]int)
	sampled := 0

	var startTime, endTime int64
	if timeSeconds > 0 {
		center := timeSeconds * 1000
		startTime = center - 8*60*1000
		endTime = center + 8*60*1000
	} else {
		now := time.Now().UnixMilli()
		startTime = now - 15*60*1000
		endTime = now
	}
	// The scan reads through the fetch-all layer; the sample sizes stay
	// explicit bounds over the complete stream list rather than an
	// inherited single-page truncation.
	streams, err := fetchAllLogStreams(store, logGroupName, "")
	if err != nil {
		return nil, mapStoreError(err)
	}
	if len(streams) > logsstore.GetLogGroupFieldsSampleSize {
		streams = streams[:logsstore.GetLogGroupFieldsSampleSize]
	}
	for _, ls := range streams {
		events, err := fetchEventsUpTo(store, logGroupName, ls.Name, startTime, endTime, logsstore.GetLogGroupFieldsSampleSize)
		if err != nil {
			continue
		}
		for _, evt := range events {
			sampled++
			var data map[string]interface{}
			if json.Unmarshal([]byte(evt.Message), &data) == nil {
				for k := range data {
					jsonFieldCounts[k]++
				}
			}
		}
	}

	fields := make([]map[string]interface{}, 0, len(standardFields)+len(jsonFieldCounts))
	for _, f := range standardFields {
		fields = append(fields, map[string]interface{}{
			"name":    f,
			"percent": 100,
		})
	}

	jsonFieldNames := make([]string, 0, len(jsonFieldCounts))
	for f := range jsonFieldCounts {
		jsonFieldNames = append(jsonFieldNames, f)
	}
	sort.Strings(jsonFieldNames)
	for _, f := range jsonFieldNames {
		percent := 0
		if sampled > 0 {
			percent = jsonFieldCounts[f] * 100 / sampled
		}
		fields = append(fields, map[string]interface{}{
			"name":    f,
			"percent": percent,
		})
	}

	return fields, nil
}

// getLogFieldsCore validates input and retrieves field names from log events.
func (s *LogsService) getLogFieldsCore(store *logsstore.Store, dataSourceName, dataSourceType string) ([]string, error) {
	jsonFields := make(map[string]bool)
	// Same seam and sample-bound discipline as getLogGroupFieldsCore.
	streams, err := fetchAllLogStreams(store, dataSourceName, "")
	if err != nil {
		return nil, mapStoreError(err)
	}
	if len(streams) > logsstore.GetLogFieldsSampleSize {
		streams = streams[:logsstore.GetLogFieldsSampleSize]
	}
	for _, ls := range streams {
		events, err := fetchEventsUpTo(store, dataSourceName, ls.Name, 0, 0, logsstore.GetLogFieldsSampleSize)
		if err != nil {
			continue
		}
		for _, evt := range events {
			var data map[string]interface{}
			if json.Unmarshal([]byte(evt.Message), &data) == nil {
				for k := range data {
					jsonFields[k] = true
				}
			}
		}
	}

	fieldNames := []string{"@timestamp", "@message", "@logStream"}
	for f := range jsonFields {
		fieldNames = append(fieldNames, f)
	}
	return fieldNames, nil
}
