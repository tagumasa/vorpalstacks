package cloudwatchlogs

import (
	"context"
	"sort"

	"vorpalstacks/internal/common/invokers"
	corelogs "vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The fetch-all read layer: the single seam every internal bulk consumer
// (Insights query gather, export execution, import object listing, field
// sampling) reads through. Each helper pages its store or invoker call to
// exhaustion so a bulk consumer never inherits a single-page truncation;
// callers that want a sample cap pass it as an explicit bound instead of
// relying on a page size.

// fetchAllLogStreams lists every log stream of a group under an optional
// name prefix, in the store's ascending name order, paging to exhaustion.
func fetchAllLogStreams(store *logsstore.Store, logGroupName, prefix string) ([]*logsstore.LogStream, error) {
	var all []*logsstore.LogStream
	marker := ""
	for {
		streams, nextMarker, err := store.ListLogStreams(logGroupName, prefix, marker, logsstore.ListingPageSize)
		if err != nil {
			return nil, err
		}
		all = append(all, streams...)
		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}
	return all, nil
}

// fetchAllLogGroups lists every log group of the region's store in the
// store's ascending name order, paging to exhaustion — the base the
// pattern-filtered DescribeLogGroups listing selects from.
func fetchAllLogGroups(store *logsstore.Store) ([]*logsstore.LogGroup, error) {
	var all []*logsstore.LogGroup
	marker := ""
	for {
		groups, nextMarker, err := store.ListLogGroups("", marker, logsstore.ListingPageSize)
		if err != nil {
			return nil, err
		}
		all = append(all, groups...)
		if nextMarker == "" || len(groups) == 0 {
			break
		}
		marker = nextMarker
	}
	return all, nil
}

// fetchAllMetricFilters lists every metric filter of a group under an
// optional name prefix, paging to exhaustion — the ingestion fan-out
// evaluates the group's complete filter set, never one listing page.
func fetchAllMetricFilters(store *logsstore.Store, logGroupName, prefix string) ([]*logsstore.MetricFilter, error) {
	var all []*logsstore.MetricFilter
	marker := ""
	for {
		filters, nextMarker, err := store.ListMetricFilters(logGroupName, prefix, marker, logsstore.ListingPageSize)
		if err != nil {
			return nil, err
		}
		all = append(all, filters...)
		if nextMarker == "" || len(filters) == 0 {
			break
		}
		marker = nextMarker
	}
	return all, nil
}

// fetchAllLogEvents reads every event of one stream inside the window, in
// the read engine's deterministic order, paging to exhaustion. The endTime
// boundary is exclusive — the read engine's default, for internal windows
// whose boundary no API document defines.
func fetchAllLogEvents(store *logsstore.Store, logGroupName, logStreamName string, startTime, endTime int64) ([]*logsstore.OutputLogEvent, error) {
	return fetchEventsWindow(store, logGroupName, logStreamName, startTime, endTime, 0, false)
}

// fetchAllLogEventsEndInclusive reads the same full window with an
// inclusive endTime — the boundary the APIs that define these windows
// document ("The range is inclusive, so the specified end time is included
// in the query", StartQuery; "Events with a timestamp later than this time
// are not exported", CreateExportTask). The Insights query gather and the
// export executor read through this form.
func fetchAllLogEventsEndInclusive(store *logsstore.Store, logGroupName, logStreamName string, startTime, endTime int64) ([]*logsstore.OutputLogEvent, error) {
	return fetchEventsWindow(store, logGroupName, logStreamName, startTime, endTime, 0, true)
}

// fetchEventsUpTo reads a stream's windowed events, paging until the
// stream is exhausted or max events have been collected. max <= 0 means
// no bound. The endTime boundary stays the exclusive default; every direct
// caller defines an internal sampling or cursor window.
func fetchEventsUpTo(store *logsstore.Store, logGroupName, logStreamName string, startTime, endTime int64, max int) ([]*logsstore.OutputLogEvent, error) {
	return fetchEventsWindow(store, logGroupName, logStreamName, startTime, endTime, max, false)
}

// fetchEventsWindow is the single fetch-all paging loop over one stream's
// window; endInclusive selects the page read whose endTime boundary the
// consuming API documents.
func fetchEventsWindow(store *logsstore.Store, logGroupName, logStreamName string, startTime, endTime int64, max int, endInclusive bool) ([]*logsstore.OutputLogEvent, error) {
	var all []*logsstore.OutputLogEvent
	token := ""
	for {
		pageLimit := logsstore.MaxEventsLimit
		if max > 0 {
			remaining := max - len(all)
			if remaining <= 0 {
				break
			}
			if remaining < pageLimit {
				pageLimit = remaining
			}
		}
		var events []*logsstore.OutputLogEvent
		var nextForwardToken string
		var err error
		if endInclusive {
			events, nextForwardToken, _, err = store.GetLogEventsEndInclusive(logGroupName, logStreamName, startTime, endTime, pageLimit, true, token)
		} else {
			events, nextForwardToken, _, err = store.GetLogEvents(logGroupName, logStreamName, startTime, endTime, pageLimit, true, token)
		}
		if err != nil {
			return nil, err
		}
		all = append(all, events...)
		// GetLogEvents re-offers the presented token once the stream is
		// exhausted (its documented end-of-stream rule), so pagination
		// ends on the empty page rather than on an absent token.
		if len(events) == 0 || nextForwardToken == token {
			break
		}
		token = nextForwardToken
	}
	if max > 0 && len(all) > max {
		all = all[:max]
	}
	return all, nil
}

// fetchGroupEventsForQuery gathers the query input set: every stream of
// every named group inside the window, in one deterministic global order
// (timestamp, then ingestion time, then stream name — the read engine's
// documented key with the stream carrying the put-request component's
// place). Stream listing and event read errors surface through the server
// log; unresolvable groups yield no events, matching the documented
// behaviour of querying a group without matching events.
func fetchGroupEventsForQuery(store *logsstore.Store, groups []string, startTime, endTime int64) []logEventWithContext {
	var allEvents []logEventWithContext
	for _, lgName := range groups {
		streams, err := fetchAllLogStreams(store, lgName, "")
		if err != nil {
			corelogs.Error("Failed to list log streams for query",
				corelogs.String("logGroup", lgName), corelogs.Err(err))
			continue
		}
		// The query plane reads the transformed form where one exists —
		// "After log events have been transformed, you must use
		// CloudWatch Logs Insights queries to view the transformed
		// versions"; GetLogEvents and FilterLogEvents keep serving
		// originals.
		transformed, _ := store.TransformedMessageMap(lgName)
		for _, ls := range streams {
			events, err := fetchAllLogEventsEndInclusive(store, lgName, ls.Name, startTime, endTime)
			if err != nil {
				corelogs.Error("Failed to read log events for query",
					corelogs.String("logGroup", lgName), corelogs.String("logStream", ls.Name), corelogs.Err(err))
				continue
			}
			for _, evt := range events {
				message := evt.Message
				isTransformed := false
				if out, ok := transformed[logsstore.TransformedMessageDigest(evt.Timestamp, ls.Name, evt.Message)]; ok {
					message = out
					isTransformed = true
				}
				allEvents = append(allEvents, logEventWithContext{
					timestamp:     evt.Timestamp,
					message:       message,
					ingestionTime: evt.IngestionTime,
					logGroup:      lgName,
					logStream:     ls.Name,
					transformed:   isTransformed,
				})
			}
		}
	}
	sort.SliceStable(allEvents, func(i, j int) bool {
		if allEvents[i].timestamp != allEvents[j].timestamp {
			return allEvents[i].timestamp < allEvents[j].timestamp
		}
		if allEvents[i].ingestionTime != allEvents[j].ingestionTime {
			return allEvents[i].ingestionTime < allEvents[j].ingestionTime
		}
		return allEvents[i].logStream < allEvents[j].logStream
	})
	return allEvents
}

// fetchAllS3ObjectKeys lists every object key under a bucket prefix
// through the S3 invoker's fetch-all convention (maxKeys 0 pages the
// invoker internally, with its own safety cap).
func fetchAllS3ObjectKeys(ctx context.Context, invoker invokers.S3Invoker, region, bucket, prefix string) ([]string, error) {
	return invoker.ListObjects(ctx, region, bucket, prefix, 0)
}
