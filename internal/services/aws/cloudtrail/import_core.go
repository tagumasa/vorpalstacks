package cloudtrail

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"sort"
	"strings"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/resilience"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// StartImportInput carries the members for StartImport. A non-empty ImportID
// selects the retry path; a new import requires Destinations and
// ImportSource. The time bounds keep both wire forms (RFC3339 string and
// Unix epoch number) because the JSON protocol serialises timestamps as
// epochs while query strings arrive as RFC3339 text; an unparseable present
// value is rejected, never silently dropped.
type StartImportInput struct {
	ImportID             string
	Destinations         []string
	ImportSourceRaw      interface{}
	ImportSourceProvided bool
	StartEventTimeStr    string
	StartEventTimeRaw    interface{}
	EndEventTimeStr      string
	EndEventTimeRaw      interface{}
}

// ImportIDInput carries the import identifier.
type ImportIDInput struct {
	ImportID string
}

// ListImportsInput carries the filter and pagination members for ListImports.
type ListImportsInput struct {
	NextToken    string
	MaxResults   int
	Destination  string
	ImportStatus string
}

// ListImportFailuresInput carries the pagination members for
// ListImportFailures.
type ListImportFailuresInput struct {
	ImportID   string
	NextToken  string
	MaxResults int
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// startImportCore is the single entry point for StartImport.
func (s *CloudTrailService) startImportCore(ctx context.Context, store cloudtrailstore.CloudTrailStoreInterface, in StartImportInput) (map[string]interface{}, error) {
	// Retry case: existing ImportId provided. The admission check and the
	// INITIALIZING write run under the store mutex as one step. The attempt
	// counter increments with the reset: an executor from the previous
	// attempt that is still draining detects the change at its next
	// checkpoint and abandons instead of racing the new attempt's walk.
	if in.ImportID != "" {
		retried, err := store.MutateImport(in.ImportID, func(imp *cloudtrailstore.Import) error {
			if imp.ImportStatus != "STOPPED" && imp.ImportStatus != "FAILED" {
				return newOperationNotPermittedException(
					"Cannot retry an import that is still in progress")
			}
			imp.ImportStatus = "INITIALIZING"
			imp.Attempt++
			imp.ImportStatistics = cloudtrailstore.ImportStatistics{}
			// The retry is a fresh attempt: the previous run's failure
			// entries describe that attempt, not the one now starting.
			imp.Failures = nil
			return nil
		})
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		go s.runImport(store, retried.ImportID)
		return formatStartImportResponse(retried), nil
	}

	// New import: Destinations and ImportSource required.
	destinations := in.Destinations
	if len(destinations) == 0 {
		return nil, newInvalidParameterException(
			"Destinations is required for a new import")
	}
	// The model bounds ImportDestinations to a single entry.
	if len(destinations) > 1 {
		return nil, newInvalidParameterException(
			"Destinations must contain exactly one event data store")
	}

	if !in.ImportSourceProvided {
		return nil, newInvalidParameterException(
			"ImportSource is required for a new import")
	}

	// The time bounds parse strictly in both wire forms: a present value
	// that neither parses as RFC3339 text nor as an epoch number is an
	// invalid parameter, and EndEventTime must not precede StartEventTime.
	startEventTime, err := parseWireTime(in.StartEventTimeStr, in.StartEventTimeRaw)
	if err != nil {
		return nil, err
	}
	endEventTime, err := parseWireTime(in.EndEventTimeStr, in.EndEventTimeRaw)
	if err != nil {
		return nil, err
	}
	if startEventTime != nil && endEventTime != nil && endEventTime.Before(*startEventTime) {
		return nil, newInvalidParameterException(
			"EndEventTime must not precede StartEventTime")
	}

	source := parseImportSource(in.ImportSourceRaw)
	if source.S3LocationURI == "" {
		return nil, newInvalidImportSourceException(
			"S3LocationUri is required")
	}
	if source.S3BucketRegion == "" {
		return nil, newInvalidImportSourceException(
			"S3BucketRegion is required")
	}
	// The model marks the access role required alongside the other two
	// source members.
	if source.S3BucketAccessRoleARN == "" {
		return nil, newInvalidImportSourceException(
			"S3BucketAccessRoleArn is required")
	}
	bucket, _, err := parseS3ImportURI(source.S3LocationURI)
	if err != nil {
		return nil, err
	}

	// Validate destination EDS exists and is active.
	eds, err := s.resolveEventDataStore(store, destinations[0])
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	if eds.Status == "PENDING_DELETION" {
		return nil, newInactiveEventDataStoreException(
			"Cannot import into a PENDING_DELETION event data store")
	}

	// The ongoing-import admission is atomic with the creation write
	// (CreateImportIfNoOngoing): "This exception is thrown when you start a
	// new import and a previous import is still in progress"
	// (AccountHasOngoingImportException).

	// The source bucket must exist in the region the import names: "This
	// exception is thrown when the provided source S3 bucket is not valid
	// for import" (InvalidImportSourceException, StartImport). A source the
	// platform cannot verify — no S3 service — is rejected rather than
	// accepted unverified.
	invoker := s.s3Invoker()
	if invoker == nil {
		return nil, newInvalidImportSourceException(
			"Cannot verify the source S3 bucket: the S3 service is unavailable")
	}
	if exists, err := invoker.BucketExists(ctx, source.S3BucketRegion, bucket); err != nil {
		return nil, newInvalidImportSourceException(
			fmt.Sprintf("Cannot verify the source S3 bucket %s: %s", bucket, err.Error()))
	} else if !exists {
		return nil, newInvalidImportSourceException(
			fmt.Sprintf("The source S3 bucket does not exist: %s", bucket))
	}

	imp := cloudtrailstore.NewImport(destinations, source)

	if startEventTime != nil {
		imp.StartEventTime = startEventTime
	}
	if endEventTime != nil {
		imp.EndEventTime = endEventTime
	}

	created, err := store.CreateImportIfNoOngoing(imp)
	if err != nil {
		if errors.Is(err, cloudtrailstore.ErrImportOngoing) {
			return nil, newAccountHasOngoingImportException(
				"Cannot start a new import while a previous import is still in progress")
		}
		return nil, s.mapStoreError(err)
	}

	go s.runImport(store, created.ImportID)

	return formatStartImportResponse(created), nil
}

// importTerminalStatus reports whether an import status admits no further
// transitions.
func importTerminalStatus(status string) bool {
	return status == "COMPLETED" || status == "FAILED" || status == "STOPPED"
}

// stopImportCore is the single entry point for StopImport. The terminal
// check and the STOPPED write run under the store mutex as one step, so a
// completion that lands first makes the stop refuse, and a stop that lands
// first holds against the executor's remaining transitions.
func (s *CloudTrailService) stopImportCore(store cloudtrailstore.CloudTrailStoreInterface, in ImportIDInput) (map[string]interface{}, error) {
	if in.ImportID == "" {
		return nil, newInvalidParameterException(
			"ImportId is required")
	}

	imp, err := store.MutateImport(in.ImportID, func(imp *cloudtrailstore.Import) error {
		if importTerminalStatus(imp.ImportStatus) {
			return newOperationNotPermittedException(
				"Cannot stop an import that has already finished")
		}
		imp.ImportStatus = "STOPPED"
		return nil
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return formatImportDetail(imp), nil
}

// getImportCore is the single entry point for GetImport.
func (s *CloudTrailService) getImportCore(store cloudtrailstore.CloudTrailStoreInterface, in ImportIDInput) (map[string]interface{}, error) {
	if in.ImportID == "" {
		return nil, newInvalidParameterException(
			"ImportId is required")
	}

	imp, err := store.GetImport(in.ImportID)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return formatImportDetail(imp), nil
}

// listImportsCore is the single entry point for ListImports.
func (s *CloudTrailService) listImportsCore(store cloudtrailstore.CloudTrailStoreInterface, in ListImportsInput) (map[string]interface{}, error) {
	if in.ImportStatus != "" {
		if err := validateImportStatusFilter(in.ImportStatus); err != nil {
			return nil, err
		}
	}
	opts := storecommon.ListOptions{MaxItems: cloudtrailstore.DefaultListImportsResults}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}
	if in.MaxResults > 0 {
		// ListImports declares no InvalidMaxResultsException; its generic
		// declared invalid-input error carries the bound rejection.
		if in.MaxResults > cloudtrailstore.MaxListImportsResults {
			return nil, newInvalidParameterException(
				fmt.Sprintf("MaxResults exceeds the maximum of %d", cloudtrailstore.MaxListImportsResults))
		}
		opts.MaxItems = in.MaxResults
	}

	result, err := store.ListImports(opts, in.Destination, in.ImportStatus)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	imports := make([]map[string]interface{}, 0, len(result.Items))
	for _, imp := range result.Items {
		imports = append(imports, formatImportBase(imp))
	}

	resp := map[string]interface{}{
		"Imports": imports,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}

	return resp, nil
}

// listImportFailuresCore is the single entry point for ListImportFailures.
func (s *CloudTrailService) listImportFailuresCore(store cloudtrailstore.CloudTrailStoreInterface, in ListImportFailuresInput) (map[string]interface{}, error) {
	if in.ImportID == "" {
		return nil, newInvalidParameterException(
			"ImportId is required")
	}

	opts := storecommon.ListOptions{MaxItems: cloudtrailstore.DefaultListImportsResults}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}
	if in.MaxResults > 0 {
		opts.MaxItems = in.MaxResults
	}

	result, err := store.ListImportFailures(in.ImportID, opts)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	failures := make([]map[string]interface{}, 0, len(result.Items))
	for _, f := range result.Items {
		fm := map[string]interface{}{}
		if f.Location != "" {
			fm["Location"] = f.Location
		}
		if f.Status != "" {
			fm["Status"] = f.Status
		}
		if f.ErrorType != "" {
			fm["ErrorType"] = f.ErrorType
		}
		if f.ErrorMessage != "" {
			fm["ErrorMessage"] = f.ErrorMessage
		}
		if f.LastUpdatedTime != nil {
			fm["LastUpdatedTime"] = f.LastUpdatedTime.Unix()
		}
		failures = append(failures, fm)
	}

	resp := map[string]interface{}{
		"Failures": failures,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}

	return resp, nil
}

// runImport executes the import asynchronously: it walks the source S3
// location, downloads each gzip log file, and copies every trail record
// inside the time bounds and the destination's retention into the
// destination event data store — "CloudTrail copies all trail events
// regardless of the configuration of the destination event data store's
// event types, advanced event selectors, or AWS Region" (user guide,
// considerations for copying trail events), so the writes go directly
// into the store. "An import finishes with a status of COMPLETED if there
// were no failures, or FAILED if there were failures" (StartImport
// response, ImportStatus). Every transition runs under the store mutex
// and only from its expected prior status: a StopImport that lands
// mid-execution has already written STOPPED, and no later transition
// overwrites that verdict.
func (s *CloudTrailService) runImport(store cloudtrailstore.CloudTrailStoreInterface, importID string) {
	defer func() {
		if r := recover(); r != nil {
			resilience.LogPanic("cloudtrail.runImport", r)
			_, _ = store.MutateImport(importID, func(imp *cloudtrailstore.Import) error {
				if importTerminalStatus(imp.ImportStatus) {
					return cloudtrailstore.ErrUnchanged
				}
				imp.ImportStatus = "FAILED"
				return nil
			})
			slog.Error("Panic recovered in runImport",
				"importId", importID, "panic", r)
		}
	}()

	rec, err := store.MutateImport(importID, func(imp *cloudtrailstore.Import) error {
		if imp.ImportStatus != "INITIALIZING" {
			return cloudtrailstore.ErrUnchanged
		}
		imp.ImportStatus = "IN_PROGRESS"
		return nil
	})
	if err != nil {
		slog.Error("Failed to update import status to IN_PROGRESS",
			"importId", importID, "error", err)
		return
	}
	if rec.ImportStatus != "IN_PROGRESS" {
		// Someone else moved the record out of INITIALIZING first; this
		// executor stands down.
		return
	}
	// The attempt number this executor runs for: every checkpoint and the
	// finalise compare it against the record, so a superseded executor (a
	// retry reset the record while this goroutine was still draining) can
	// neither continue the walk nor settle the verdict.
	attempt := rec.Attempt

	imp, err := store.GetImport(importID)
	if err != nil {
		slog.Error("Failed to read the import record",
			"importId", importID, "error", err)
		return
	}

	stats, failures := s.executeImport(store, importID, imp, attempt)

	status := "COMPLETED"
	if stats.FailedEntries > 0 {
		status = "FAILED"
	}
	if _, err := store.MutateImport(importID, func(imp *cloudtrailstore.Import) error {
		if imp.ImportStatus != "IN_PROGRESS" || imp.Attempt != attempt {
			return cloudtrailstore.ErrUnchanged
		}
		imp.ImportStatus = status
		imp.ImportStatistics = stats
		imp.Failures = append(imp.Failures, failures...)
		return nil
	}); err != nil {
		slog.Error("Failed to finalise the import",
			"importId", importID, "status", status, "error", err)
	}
}

// executeImport performs the import walk and returns the measured
// statistics and failure entries. The progress statistics are written
// after every file — also the point a StopImport between files, or a retry
// that superseded this attempt, is honoured: a write that finds the import
// no longer IN_PROGRESS on this attempt abandons the walk with what it
// has, leaving the verdict to whoever flipped the status or restarted.
func (s *CloudTrailService) executeImport(store cloudtrailstore.CloudTrailStoreInterface, importID string, imp *cloudtrailstore.Import, attempt int64) (cloudtrailstore.ImportStatistics, []cloudtrailstore.ImportFailure) {
	var stats cloudtrailstore.ImportStatistics
	var failures []cloudtrailstore.ImportFailure

	ctx := s.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	invoker := s.s3Invoker()
	bucket, prefix, err := parseS3ImportURI(imp.ImportSource.S3LocationURI)
	if invoker == nil || err != nil {
		message := "the S3 service is unavailable"
		if err != nil {
			message = err.Error()
		}
		failures = append(failures, importFailure(imp.ImportSource.S3LocationURI, "InvalidImportSource", message))
		stats.FailedEntries++
		return stats, failures
	}
	// The S3LocationUri names one location prefix; the walk covers it, and
	// the statistics count it as found once and completed once the walk
	// ends. AWS documents neither a finer PrefixesFound granularity nor
	// this one — the single location prefix is the reading the platform
	// pins.
	stats.PrefixesFound = 1

	keys, err := invoker.ListObjects(ctx, imp.ImportSource.S3BucketRegion, bucket, prefix, 0)
	if err != nil {
		failures = append(failures, importFailure(
			fmt.Sprintf("s3://%s/%s", bucket, prefix), "ListObjectsError", err.Error()))
		stats.FailedEntries++
		return stats, failures
	}
	sort.Strings(keys)

	eds, err := s.resolveEventDataStore(store, imp.Destinations[0])
	if err != nil {
		failures = append(failures, importFailure(imp.ImportSource.S3LocationURI, "DestinationUnavailable", err.Error()))
		stats.FailedEntries++
		return stats, failures
	}
	edsID := eds.EventDataStoreID
	// "CloudTrail only copies trail events that have an eventTime within
	// the event data store's retention period" (user guide).
	retentionCutoff := time.Now().UTC().AddDate(0, 0, -int(eds.RetentionPeriod))

	for _, key := range keys {
		if !importFileCandidate(key, prefix) {
			continue
		}
		if !importFileWithinBounds(key, imp.StartEventTime, imp.EndEventTime) {
			continue
		}
		location := fmt.Sprintf("s3://%s/%s", bucket, key)
		events, fileErr, recordFailures := readImportFile(ctx, invoker, imp.ImportSource.S3BucketRegion, bucket, key, retentionCutoff)
		if fileErr != nil {
			failures = append(failures, *fileErr)
			stats.FailedEntries++
			continue
		}
		for _, f := range recordFailures {
			failures = append(failures, f)
			stats.FailedEntries++
		}
		imported := 0
		for _, event := range events {
			if err := store.PutEventIntoEDS(edsID, event); err != nil {
				code := "WriteError"
				if errors.Is(err, cloudtrailstore.ErrEventDataStoreNotFound) {
					// The destination record vanished mid-walk (the restore
					// window's hard delete): the destination is gone, not
					// merely unwritable.
					code = "DestinationUnavailable"
				}
				failures = append(failures, importFailure(location, code, err.Error()))
				stats.FailedEntries++
				continue
			}
			imported++
		}
		stats.EventsCompleted += int64(imported)
		stats.FilesCompleted++

		// MutateImport answers a refused guard with the untouched record
		// and no error, so the refusal is detected on the returned record:
		// a writer that flipped the import out of IN_PROGRESS, or a retry
		// that bumped the attempt number, owns the verdict and the walk
		// abandons with what it has.
		after, err := store.MutateImport(importID, func(imp *cloudtrailstore.Import) error {
			if imp.ImportStatus != "IN_PROGRESS" || imp.Attempt != attempt {
				return cloudtrailstore.ErrUnchanged
			}
			imp.ImportStatistics = stats
			return nil
		})
		if err != nil {
			slog.Error("Failed to record import progress",
				"importId", importID, "error", err)
			return stats, failures
		}
		if after.ImportStatus != "IN_PROGRESS" || after.Attempt != attempt {
			slog.Warn("Import stopped, finished, or superseded underneath the walk; abandoning",
				"importId", importID, "status", after.ImportStatus)
			return stats, failures
		}
	}
	stats.PrefixesCompleted = 1
	return stats, failures
}

// readImportFile downloads and decodes one gzip trail log file into
// events. Its failures split by level: fileErr is set when the file itself
// could not be read, decompressed, or parsed as a {"Records":[...]} log
// file — such a file did not complete; recordFailures carries one entry
// per record that does not parse, while the file's other records still
// import and the file still counts as completed.
func readImportFile(ctx context.Context, invoker invokers.S3Invoker, region, bucket, key string, retentionCutoff time.Time) ([]*cloudtrailstore.Event, *cloudtrailstore.ImportFailure, []cloudtrailstore.ImportFailure) {
	location := fmt.Sprintf("s3://%s/%s", bucket, key)

	data, err := invoker.GetObject(ctx, region, bucket, key, 0)
	if err != nil {
		f := importFailure(location, "ObjectReadError", err.Error())
		return nil, &f, nil
	}
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		f := importFailure(location, "InvalidFileFormat",
			"The object is not a gzip-compressed CloudTrail log file")
		return nil, &f, nil
	}
	raw, err := io.ReadAll(zr)
	zr.Close()
	if err != nil {
		f := importFailure(location, "InvalidFileFormat",
			fmt.Sprintf("The gzip stream of the log file is corrupt: %s", err.Error()))
		return nil, &f, nil
	}
	var envelope struct {
		Records json.RawMessage `json:"Records"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil || len(envelope.Records) == 0 {
		f := importFailure(location, "InvalidFileFormat",
			`The log file is not a {"Records":[...]} JSON document`)
		return nil, &f, nil
	}
	var records []json.RawMessage
	if err := json.Unmarshal(envelope.Records, &records); err != nil {
		f := importFailure(location, "InvalidFileFormat",
			`The log file is not a {"Records":[...]} JSON document`)
		return nil, &f, nil
	}

	var events []*cloudtrailstore.Event
	var failures []cloudtrailstore.ImportFailure
	for _, rec := range records {
		event, err := cloudtrailstore.RecordEventFromPayload(string(rec), "", "")
		if err != nil {
			failures = append(failures, importFailure(location, "InvalidEventRecord", err.Error()))
			continue
		}
		// The retention bound drops the record silently — the copy contract
		// keeps only events inside the store's retention, it is not a
		// failure of the import.
		if event.EventTime.Before(retentionCutoff) {
			continue
		}
		events = append(events, event)
	}
	return events, nil, failures
}

// importFileCandidate reports whether one listed key is an importable
// trail log file. "CloudTrail only copies trail events from Gzip
// compressed log files" (user guide), so only .json.gz keys are
// candidates, and a walk of the whole bucket — the default when the
// S3LocationUri carries no prefix — admits only keys under a CloudTrail
// path segment: "By default, CloudTrail only imports events contained in
// the S3 bucket's CloudTrail prefix and the prefixes inside the
// CloudTrail prefix, and does not check prefixes for other AWS services"
// (StartImport). The segment match is exact, so CloudTrail-Digest files
// stay out.
func importFileCandidate(key, prefix string) bool {
	if !strings.HasSuffix(key, ".json.gz") {
		return false
	}
	if prefix != "" {
		return true
	}
	for _, segment := range strings.Split(key, "/") {
		if segment == "CloudTrail" {
			return true
		}
	}
	return false
}

// importFileNameDateRe matches the delivery-time stamp in a trail log
// file's name — <YYYYMMDDTHHmmZ> in
// <account>_CloudTrail_<region>_<stamp>_<unique>.json.gz. Digest names
// carry seconds (<YYYYMMDDTHHMMSSZ>) and do not match.
var importFileNameDateRe = regexp.MustCompile(`(\d{8}T\d{4}Z)`)

// importFileWithinBounds applies the import's time bounds at file
// granularity: "CloudTrail checks the prefix and log file names to verify
// the names contain a date between the specified StartEventTime and
// EndEventTime before attempting to import events" (StartImport). A name
// that carries no parseable date cannot be verified, so with bounds set
// the file is not attempted.
func importFileWithinBounds(key string, start, end *time.Time) bool {
	if start == nil && end == nil {
		return true
	}
	match := importFileNameDateRe.FindStringSubmatch(key)
	if match == nil {
		return false
	}
	at, err := time.Parse("20060102T1504Z", match[1])
	if err != nil {
		return false
	}
	if start != nil && at.Before(*start) {
		return false
	}
	if end != nil && at.After(*end) {
		return false
	}
	return true
}

// parseS3ImportURI splits the import source location into its bucket and
// prefix. The location must be an s3:// URI naming a bucket; anything
// after the bucket is the prefix the walk is scoped to.
func parseS3ImportURI(uri string) (string, string, error) {
	rest, ok := strings.CutPrefix(uri, "s3://")
	if !ok || rest == "" {
		return "", "", newInvalidImportSourceException(
			fmt.Sprintf("S3LocationUri must be an s3:// URI: %s", uri))
	}
	bucket, prefix, _ := strings.Cut(strings.TrimPrefix(rest, "/"), "/")
	if bucket == "" {
		return "", "", newInvalidImportSourceException(
			fmt.Sprintf("S3LocationUri must name a bucket: %s", uri))
	}
	return bucket, prefix, nil
}

// importFailure builds one ListImportFailures entry.
func importFailure(location, errorType, message string) cloudtrailstore.ImportFailure {
	now := time.Now().UTC()
	return cloudtrailstore.ImportFailure{
		Location:        location,
		Status:          "FAILED",
		ErrorType:       errorType,
		ErrorMessage:    message,
		LastUpdatedTime: &now,
	}
}

// storeHasOngoingImport reports whether the store holds any import that
// has not reached a terminal status — "still in progress" spans both
// INITIALIZING and IN_PROGRESS. An empty destinationARN scans every
// import in the store (the account-scoped StartImport guard); a non-empty
// one scopes the scan to one destination event data store (the update and
// delete guards).
func storeHasOngoingImport(store cloudtrailstore.CloudTrailStoreInterface, destinationARN string) (bool, error) {
	for _, status := range []string{"INITIALIZING", "IN_PROGRESS"} {
		result, err := store.ListImports(storecommon.ListOptions{MaxItems: 1}, destinationARN, status)
		if err != nil {
			return false, err
		}
		if len(result.Items) > 0 {
			return true, nil
		}
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// formatStartImportResponse formats an import for the StartImport response.
// Per Smithy StartImportResponse: includes ImportSource, StartEventTime,
// EndEventTime but NOT ImportStatistics.
func formatStartImportResponse(imp *cloudtrailstore.Import) map[string]interface{} {
	resp := formatImportBase(imp)
	if imp.ImportSource.S3LocationURI != "" {
		s3Map := map[string]interface{}{
			"S3LocationUri":  imp.ImportSource.S3LocationURI,
			"S3BucketRegion": imp.ImportSource.S3BucketRegion,
		}
		if imp.ImportSource.S3BucketAccessRoleARN != "" {
			s3Map["S3BucketAccessRoleArn"] = imp.ImportSource.S3BucketAccessRoleARN
		}
		resp["ImportSource"] = map[string]interface{}{
			"S3": s3Map,
		}
	}
	if imp.StartEventTime != nil {
		resp["StartEventTime"] = imp.StartEventTime.Unix()
	}
	if imp.EndEventTime != nil {
		resp["EndEventTime"] = imp.EndEventTime.Unix()
	}
	return resp
}

// formatImportBase returns the common fields shared by all import responses.
// Its field set matches the Smithy ImportsListItem shape: only ImportId,
// ImportStatus, Destinations, CreatedTimestamp, UpdatedTimestamp.
func formatImportBase(imp *cloudtrailstore.Import) map[string]interface{} {
	resp := map[string]interface{}{
		"ImportId":         imp.ImportID,
		"ImportStatus":     imp.ImportStatus,
		"CreatedTimestamp": imp.CreatedTimestamp.Unix(),
		"UpdatedTimestamp": imp.UpdatedTimestamp.Unix(),
	}
	if len(imp.Destinations) > 0 {
		resp["Destinations"] = imp.Destinations
	}
	return resp
}

// formatImportDetail formats an import for GetImport and StopImport responses.
// Includes all fields: ImportSource, StartEventTime, EndEventTime,
// ImportStatistics.
func formatImportDetail(imp *cloudtrailstore.Import) map[string]interface{} {
	resp := formatStartImportResponse(imp)
	if imp.ImportStatistics.EventsCompleted > 0 || imp.ImportStatistics.FilesCompleted > 0 ||
		imp.ImportStatistics.PrefixesFound > 0 || imp.ImportStatistics.PrefixesCompleted > 0 ||
		imp.ImportStatistics.FailedEntries > 0 {
		stats := map[string]interface{}{}
		if imp.ImportStatistics.PrefixesFound > 0 {
			stats["PrefixesFound"] = imp.ImportStatistics.PrefixesFound
		}
		if imp.ImportStatistics.PrefixesCompleted > 0 {
			stats["PrefixesCompleted"] = imp.ImportStatistics.PrefixesCompleted
		}
		if imp.ImportStatistics.FilesCompleted > 0 {
			stats["FilesCompleted"] = imp.ImportStatistics.FilesCompleted
		}
		if imp.ImportStatistics.EventsCompleted > 0 {
			stats["EventsCompleted"] = imp.ImportStatistics.EventsCompleted
		}
		if imp.ImportStatistics.FailedEntries > 0 {
			stats["FailedEntries"] = imp.ImportStatistics.FailedEntries
		}
		resp["ImportStatistics"] = stats
	}
	return resp
}

// parseImportSource parses the ImportSource member from its raw wire value.
func parseImportSource(raw interface{}) cloudtrailstore.ImportSource {
	var source cloudtrailstore.ImportSource
	m, ok := raw.(map[string]interface{})
	if !ok {
		return source
	}
	if s3Raw, ok := m["S3"].(map[string]interface{}); ok {
		if uri, ok := s3Raw["S3LocationUri"].(string); ok {
			source.S3LocationURI = uri
		}
		if region, ok := s3Raw["S3BucketRegion"].(string); ok {
			source.S3BucketRegion = region
		}
		if roleARN, ok := s3Raw["S3BucketAccessRoleArn"].(string); ok {
			source.S3BucketAccessRoleARN = roleARN
		}
	}
	return source
}

// parseDestinationsList parses a list of destination ARNs from the raw wire
// value.
func parseDestinationsList(raw interface{}) []string {
	arr, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
