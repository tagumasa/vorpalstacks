package cloudtrail

import (
	"log/slog"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/resilience"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// StartImportInput carries the members for StartImport. A non-empty ImportID
// selects the retry path; a new import requires Destinations and
// ImportSource.
type StartImportInput struct {
	ImportID             string
	Destinations         []string
	ImportSourceRaw      interface{}
	ImportSourceProvided bool
	StartEventTime       *time.Time
	EndEventTime         *time.Time
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
func (s *CloudTrailService) startImportCore(store cloudtrailstore.CloudTrailStoreInterface, in StartImportInput) (map[string]interface{}, error) {
	// Retry case: existing ImportId provided.
	if in.ImportID != "" {
		existing, err := store.GetImport(in.ImportID)
		if err != nil {
			return nil, awserrors.NewAWSError("ImportNotFoundException",
				"Import not found", 404)
		}
		if existing.ImportStatus != "STOPPED" && existing.ImportStatus != "FAILED" {
			return nil, awserrors.NewAWSError("OperationNotPermittedException",
				"Cannot retry an import that is still in progress", 400)
		}
		existing.ImportStatus = "INITIALIZING"
		existing.UpdatedTimestamp = time.Now().UTC()
		existing.ImportStatistics = cloudtrailstore.ImportStatistics{}
		if err := store.UpdateImport(existing); err != nil {
			return nil, s.mapStoreError(err)
		}
		go s.runImport(store, existing.ImportID)
		return formatStartImportResponse(existing), nil
	}

	// New import: Destinations and ImportSource required.
	destinations := in.Destinations
	if len(destinations) == 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"Destinations is required for a new import", 400)
	}

	if !in.ImportSourceProvided {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"ImportSource is required for a new import", 400)
	}
	source := parseImportSource(in.ImportSourceRaw)
	if source.S3LocationURI == "" {
		return nil, awserrors.NewAWSError("InvalidImportSourceException",
			"S3LocationUri is required", 400)
	}
	if source.S3BucketRegion == "" {
		return nil, awserrors.NewAWSError("InvalidImportSourceException",
			"S3BucketRegion is required", 400)
	}

	// Validate destination EDS exists and is active.
	edsID := destinations[0]
	if idx := strings.LastIndex(edsID, "/"); idx >= 0 {
		edsID = edsID[idx+1:]
	}
	eds, err := store.GetEventDataStore(edsID)
	if err != nil {
		return nil, awserrors.NewAWSError("EventDataStoreNotFoundException",
			"Event data store not found", 404)
	}
	if eds.Status == "PENDING_DELETION" {
		return nil, awserrors.NewAWSError("InactiveEventDataStoreException",
			"Cannot import into a PENDING_DELETION event data store", 400)
	}

	imp := cloudtrailstore.NewImport(destinations, source)

	if in.StartEventTime != nil {
		imp.StartEventTime = in.StartEventTime
	}
	if in.EndEventTime != nil {
		imp.EndEventTime = in.EndEventTime
	}

	created, err := store.CreateImport(imp)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	go s.runImport(store, created.ImportID)

	return formatStartImportResponse(created), nil
}

// stopImportCore is the single entry point for StopImport.
func (s *CloudTrailService) stopImportCore(store cloudtrailstore.CloudTrailStoreInterface, in ImportIDInput) (map[string]interface{}, error) {
	if in.ImportID == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"ImportId is required", 400)
	}

	imp, err := store.GetImport(in.ImportID)
	if err != nil {
		return nil, awserrors.NewAWSError("ImportNotFoundException",
			"Import not found", 404)
	}

	if imp.ImportStatus == "COMPLETED" || imp.ImportStatus == "FAILED" || imp.ImportStatus == "STOPPED" {
		return nil, awserrors.NewAWSError("OperationNotPermittedException",
			"Cannot stop an import that has already finished", 400)
	}

	imp.ImportStatus = "STOPPED"
	if err := store.UpdateImport(imp); err != nil {
		return nil, s.mapStoreError(err)
	}

	return formatImportDetail(imp), nil
}

// getImportCore is the single entry point for GetImport.
func (s *CloudTrailService) getImportCore(store cloudtrailstore.CloudTrailStoreInterface, in ImportIDInput) (map[string]interface{}, error) {
	if in.ImportID == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"ImportId is required", 400)
	}

	imp, err := store.GetImport(in.ImportID)
	if err != nil {
		return nil, awserrors.NewAWSError("ImportNotFoundException",
			"Import not found", 404)
	}

	return formatImportDetail(imp), nil
}

// listImportsCore is the single entry point for ListImports.
func (s *CloudTrailService) listImportsCore(store cloudtrailstore.CloudTrailStoreInterface, in ListImportsInput) (map[string]interface{}, error) {
	opts := storecommon.ListOptions{MaxItems: 50}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}
	if in.MaxResults > 0 {
		opts.MaxItems = in.MaxResults
	}

	result, err := store.ListImports(opts, in.Destination, in.ImportStatus)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	imports := make([]map[string]interface{}, 0, len(result.Items))
	for _, imp := range result.Items {
		imports = append(imports, formatImportsListItem(imp))
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
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"ImportId is required", 400)
	}

	opts := storecommon.ListOptions{MaxItems: 50}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}
	if in.MaxResults > 0 {
		opts.MaxItems = in.MaxResults
	}

	result, err := store.ListImportFailures(in.ImportID, opts)
	if err != nil {
		if err == cloudtrailstore.ErrImportNotFound {
			return nil, awserrors.NewAWSError("ImportNotFoundException",
				"Import not found", 404)
		}
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

// runImport simulates the import process asynchronously. Since there is no
// actual S3 bucket to read from, it transitions the import through
// IN_PROGRESS to COMPLETED with simulated statistics.
func (s *CloudTrailService) runImport(store cloudtrailstore.CloudTrailStoreInterface, importID string) {
	defer func() {
		if r := resilience.RecoverPanic("cloudtrail.runImport"); r != nil {
			imp, err := store.GetImport(importID)
			if err == nil {
				imp.ImportStatus = "FAILED"
				imp.UpdatedTimestamp = time.Now().UTC()
				_ = store.UpdateImport(imp)
			}
			slog.Error("Panic recovered in runImport",
				"importId", importID, "panic", r)
		}
	}()

	imp, err := store.GetImport(importID)
	if err != nil {
		slog.Error("Failed to get import for processing",
			"importId", importID, "error", err)
		return
	}

	imp.ImportStatus = "IN_PROGRESS"
	imp.UpdatedTimestamp = time.Now().UTC()
	if err := store.UpdateImport(imp); err != nil {
		slog.Error("Failed to update import status to IN_PROGRESS",
			"importId", importID, "error", err)
		return
	}

	// Simulate processing completion.
	imp.ImportStatus = "COMPLETED"
	imp.UpdatedTimestamp = time.Now().UTC()
	imp.ImportStatistics = cloudtrailstore.ImportStatistics{
		PrefixesFound:     1,
		PrefixesCompleted: 1,
		FilesCompleted:    1,
		EventsCompleted:   0,
	}
	if err := store.UpdateImport(imp); err != nil {
		slog.Error("Failed to update import status to COMPLETED",
			"importId", importID, "error", err)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// formatImportsListItem formats an import for the ListImports response.
// Per Smithy ImportsListItem: only ImportId, ImportStatus, Destinations,
// CreatedTimestamp, UpdatedTimestamp.
func formatImportsListItem(imp *cloudtrailstore.Import) map[string]interface{} {
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
