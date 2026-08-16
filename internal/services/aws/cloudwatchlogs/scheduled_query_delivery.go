package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"strings"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// Destination delivery statuses of the ActionStatus enumeration reported
// through GetScheduledQueryHistory.
const (
	destinationStatusComplete    = "COMPLETE"
	destinationStatusClientError = "CLIENT_ERROR"
	destinationStatusFailed      = "FAILED"
)

// deliverScheduledQueryResults renders the result rows and delivers them to
// every configured destination, returning one delivery record per
// destination.
func (s *LogsService) deliverScheduledQueryResults(region string, store *logsstore.Store, sq *logsstore.ScheduledQuery, queryId string, rows []queryResultRow) []*logsstore.ScheduledQueryDestination {
	var dests []*logsstore.ScheduledQueryDestination
	if sq.DestinationConfiguration == nil {
		return dests
	}
	if lt, ok := sq.DestinationConfiguration["lookupTableConfiguration"].(map[string]interface{}); ok {
		dests = append(dests, s.deliverToLookupTable(region, store, lt, rows))
	}
	if s3c, ok := sq.DestinationConfiguration["s3Configuration"].(map[string]interface{}); ok {
		dests = append(dests, s.deliverToS3(region, s3c, queryId, rows))
	}
	return dests
}

// deliverToLookupTable populates or refreshes the configured lookup table
// with the query results, creating the table on first delivery.
func (s *LogsService) deliverToLookupTable(region string, store *logsstore.Store, cfg map[string]interface{}, rows []queryResultRow) *logsstore.ScheduledQueryDestination {
	name := destinationParam(cfg, "tableName")
	dest := &logsstore.ScheduledQueryDestination{
		DestinationType:       "LOOKUP_TABLE",
		DestinationIdentifier: lookupTableArn(region, s.accountID, name),
		ProcessedIdentifier:   name,
		Status:                destinationStatusComplete,
	}
	fail := func(status, message string) *logsstore.ScheduledQueryDestination {
		dest.Status = status
		dest.ErrorMessage = message
		return dest
	}

	body := rowsToCSV(rows)
	if body == "" {
		return fail(destinationStatusClientError, "query returned no result rows to populate the lookup table")
	}
	existing, err := store.GetLookupTable(name)
	if err != nil {
		// The table does not exist yet: create it from the results. The
		// configuration's description, KMS key and tags apply only during
		// initial table creation.
		in := &LookupTableInput{
			Name:        name,
			Description: destinationParam(cfg, "description"),
			TableBody:   body,
			KmsKeyId:    destinationParam(cfg, "kmsKeyId"),
			Tags:        destinationTags(cfg),
		}
		if _, _, err := s.createLookupTableCore(store, in, region); err != nil {
			return fail(destinationStatusClientError, fmt.Sprintf("failed to create lookup table: %v", err))
		}
		return dest
	}
	if err := s.applyLookupTableBody(existing, body, region); err != nil {
		return fail(destinationStatusClientError, fmt.Sprintf("failed to refresh lookup table: %v", err))
	}
	if err := store.PutLookupTable(existing); err != nil {
		return fail(destinationStatusFailed, fmt.Sprintf("failed to store lookup table: %v", err))
	}
	return dest
}

// deliverToS3 writes the query results as gzipped CSV under the configured
// S3 URI prefix.
func (s *LogsService) deliverToS3(region string, cfg map[string]interface{}, queryId string, rows []queryResultRow) *logsstore.ScheduledQueryDestination {
	uri := destinationParam(cfg, "destinationIdentifier")
	dest := &logsstore.ScheduledQueryDestination{
		DestinationType:       "S3",
		DestinationIdentifier: uri,
		Status:                destinationStatusComplete,
	}
	fail := func(status, message string) *logsstore.ScheduledQueryDestination {
		dest.Status = status
		dest.ErrorMessage = message
		return dest
	}

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if _, err := gw.Write([]byte(rowsToCSV(rows))); err != nil {
		return fail(destinationStatusFailed, fmt.Sprintf("gzip error: %v", err))
	}
	if err := gw.Close(); err != nil {
		return fail(destinationStatusFailed, fmt.Sprintf("gzip error: %v", err))
	}
	bucket, prefix := splitS3DestinationURI(uri)
	key := queryId + "/results.csv.gz"
	if prefix != "" {
		key = prefix + "/" + key
	}
	if s.bus == nil || s.bus.S3Invoker() == nil {
		return fail(destinationStatusFailed, "S3 delivery is not available")
	}
	if err := s.bus.S3Invoker().PutObject(context.Background(), region, bucket, key, buf.Bytes(), "application/x-gzip"); err != nil {
		return fail(destinationStatusFailed, fmt.Sprintf("S3 upload error: %v", err))
	}
	dest.ProcessedIdentifier = key
	return dest
}

// splitS3DestinationURI splits an s3:// URI into its bucket and key prefix.
func splitS3DestinationURI(uri string) (bucket, prefix string) {
	rest := strings.TrimPrefix(uri, "s3://")
	if idx := strings.Index(rest, "/"); idx >= 0 {
		return rest[:idx], strings.Trim(rest[idx+1:], "/")
	}
	return rest, ""
}
