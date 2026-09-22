package cloudwatchlogs

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
)

// --- CreateLogStream ---

type CreateLogStreamInput struct {
	LogGroupName  string
	LogStreamName string
	Region        string
}

func (s *LogsService) createLogStreamCore(input CreateLogStreamInput) error {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return err
	}
	if err := validateLogStreamName(input.LogStreamName); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	ls := logsstore.NewLogStream(input.LogStreamName, input.LogGroupName)
	if err := store.CreateLogStream(ls); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// --- DeleteLogStream ---

type DeleteLogStreamInput struct {
	LogGroupName  string
	LogStreamName string
	Region        string
}

func (s *LogsService) deleteLogStreamCore(input DeleteLogStreamInput) error {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return err
	}
	if err := validateLogStreamName(input.LogStreamName); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	if err := store.DeleteLogStream(input.LogGroupName, input.LogStreamName); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// --- PutLogEvents ---

// PutLogEvent is one parsed input log event: the store entry plus the
// wire-presence fact the required-member row needs. The timestamp's valid
// range starts at zero, so a present zero is a value the age rules reject
// per event, while an absent member must reject the request — the two are
// indistinguishable on the entry alone.
type PutLogEvent struct {
	logsstore.LogEntry
	TimestampSet bool
}

type PutLogEventsInput struct {
	LogGroupName  string
	LogStreamName string
	Events        []PutLogEvent
	Region        string
}

type PutLogEventsResult struct {
	NextSequenceToken string
	RejectedLogEvents map[string]interface{}
}

// putLogEventsCore is the single ingestion seam (see ingestion.go): every
// write path — the HTTP handler, the bus handlers and the cross-service
// invoker — validates, writes and fans out to metric and subscription
// filters through this function.
func (s *LogsService) putLogEventsCore(input PutLogEventsInput) (*PutLogEventsResult, error) {
	// The modelled name patterns reject here, ahead of the store's
	// existence discrimination: a malformed name is a parameter error, not
	// a ResourceNotFound.
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return nil, err
	}
	if err := validateLogStreamName(input.LogStreamName); err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	if len(input.Events) == 0 {
		return nil, errRequiredMember("logEvents")
	}

	if len(input.Events) > logsstore.MaxBatchLogEvents {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Maximum number of log events in a single batch is %d", logsstore.MaxBatchLogEvents), 400)
	}

	// The InputLogEvent required members: a timestamp the request must
	// carry and a message of at least one character ("Minimum length of
	// 1", live reference). The Timestamp shape's range starts at zero,
	// so a negative timestamp is a range violation, not an age
	// rejection. All three reject the whole request.
	batchBytes := 0
	entries := make([]logsstore.LogEntry, 0, len(input.Events))
	for _, e := range input.Events {
		if !e.TimestampSet {
			return nil, NewLogsError("InvalidParameterException",
				"The timestamp member of every log event is required", 400)
		}
		if e.Timestamp < 0 {
			return nil, NewLogsError("InvalidParameterException",
				"The timestamp member of every log event must not be negative", 400)
		}
		if e.Message == "" {
			return nil, NewLogsError("InvalidParameterException",
				"The message member of every log event must be at least 1 character", 400)
		}
		if len(e.Message) > logsstore.MaxLogEventMessageBytes {
			return nil, NewLogsError("InvalidParameterException",
				fmt.Sprintf("Each log event can be no larger than %d bytes", logsstore.MaxLogEventMessageBytes), 400)
		}
		batchBytes += len(e.Message) + logsstore.PutLogEventOverheadBytes
		entries = append(entries, e.LogEntry)
	}
	// "The maximum batch size is 1,048,576 bytes. This size is calculated
	// as the sum of all event messages in UTF-8, plus 26 bytes for each
	// log event" (PutLogEvents API reference).
	if batchBytes > logsstore.MaxPutLogEventsBatchBytes {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("The maximum batch size is %d bytes, calculated as the sum of all event messages in UTF-8 plus %d bytes for each log event",
				logsstore.MaxPutLogEventsBatchBytes, logsstore.PutLogEventOverheadBytes), 400)
	}

	// The group's retention horizon participates in the batch validation
	// ("Events older than 14 days or preceding the log group's retention
	// period are rejected while processing remaining valid events",
	// PutLogEvents operation documentation); a group without a retention
	// setting never expires.
	retentionCutoff := int64(0)
	if lg, err := store.GetLogGroup(input.LogGroupName); err == nil {
		retentionCutoff = retentionCutoffMillis(time.Now().UnixMilli(), lg.RetentionInDays)
	}
	validEvents, rejectedInfo, valErr := validateLogEvents(entries, retentionCutoff)
	if valErr != nil {
		return nil, valErr
	}
	if len(validEvents) == 0 {
		return &PutLogEventsResult{
			NextSequenceToken: "",
			RejectedLogEvents: rejectedInfo,
		}, nil
	}

	nextToken, err := store.PutLogEvents(input.LogGroupName, input.LogStreamName, validEvents)
	if err != nil {
		return nil, mapStoreError(err)
	}

	// The effective transformer runs at ingestion (the documented
	// contract: transformations happen only during ingestion, the
	// original events stay the GetLogEvents/FilterLogEvents surface and
	// the transformed copies feed the query plane and the
	// applyOnTransformedLogs filters).
	transformed := s.transformIngestedBatch(store, input.Region, input.LogGroupName, input.LogStreamName, validEvents)

	// The field-index bounds scan rides the same ingestion pass, over
	// the transformed form where one exists — the same view the query
	// plane serves.
	s.scanIngestedFieldIndexes(store, input.LogGroupName, input.LogStreamName, validEvents, transformed)

	s.dispatchFilterFanOut(store, input.Region, input.LogGroupName, input.LogStreamName, validEvents, transformed)

	return &PutLogEventsResult{
		NextSequenceToken: nextToken,
		RejectedLogEvents: rejectedInfo,
	}, nil
}

// --- GetLogEvents ---

type GetLogEventsInput struct {
	LogGroupName  string
	LogStreamName string
	StartTime     int64
	EndTime       int64
	Limit         int32
	StartFromHead bool
	NextToken     string
	Region        string
}

type GetLogEventsResult struct {
	Events            []*logsstore.OutputLogEvent
	NextForwardToken  string
	NextBackwardToken string
}

func (s *LogsService) getLogEventsCore(input GetLogEventsInput) (*GetLogEventsResult, error) {
	// The modelled name patterns reject as parameter errors rather than
	// surfacing as ResourceNotFound from the store walk.
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return nil, err
	}
	if err := validateLogStreamName(input.LogStreamName); err != nil {
		return nil, err
	}
	// The window members target Timestamp, whose @range declares a minimum
	// of 0 ("Valid Range: Minimum value of 0").
	if input.StartTime < 0 || input.EndTime < 0 {
		return nil, NewLogsError("InvalidParameterException",
			"startTime and endTime must be non-negative epoch milliseconds", 400)
	}

	limit, err := validateListLimit(input.Limit, logsstore.DefaultEventsLimit, logsstore.MaxEventsLimit)
	if err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	events, nextForwardToken, nextBackwardToken, err := store.GetLogEvents(
		input.LogGroupName, input.LogStreamName, input.StartTime, input.EndTime,
		int(limit), input.StartFromHead, input.NextToken,
	)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &GetLogEventsResult{
		Events:            events,
		NextForwardToken:  nextForwardToken,
		NextBackwardToken: nextBackwardToken,
	}, nil
}

// --- FilterLogEvents ---

type FilterLogEventsInput struct {
	LogGroupName      string
	LogStreamNames    []string
	LogStreamNamePref string
	StartTime         int64
	EndTime           int64
	FilterPattern     string
	Limit             int32
	StartFromHead     bool
	NextToken         string
	Region            string
}

type FilterLogEventsResult struct {
	Events    []*logsstore.OutputLogEvent
	NextToken string
}

func (s *LogsService) filterLogEventsCore(input FilterLogEventsInput) (*FilterLogEventsResult, error) {
	// The group name's modelled pattern rejects as a parameter error
	// rather than surfacing as ResourceNotFound from the store walk.
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return nil, err
	}

	if len(input.LogStreamNames) > 0 && input.LogStreamNamePref != "" {
		return nil, NewLogsError("InvalidParameterException",
			"Cannot specify both logStreamNames and logStreamNamePrefix", 400)
	}

	// The stream array carries the InputLogStreamNames @length max ("Array
	// Members: Minimum number of 1 item. Maximum number of 100 items"); an
	// absent array means every stream, so the maximum alone rejects.
	if len(input.LogStreamNames) > logsstore.MaxInputLogStreamNames {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("logStreamNames must contain at most %d items", logsstore.MaxInputLogStreamNames), 400)
	}
	// The same trait's minimum: a present-but-empty parse (marked by the
	// non-nil empty slice) rejects instead of silently meaning "absent".
	if input.LogStreamNames != nil && len(input.LogStreamNames) == 0 {
		return nil, NewLogsError("InvalidParameterException",
			"logStreamNames must contain at least 1 entry", 400)
	}

	// The prefix targets LogStreamName and so carries its pattern and
	// length ceiling.
	if err := validateLogStreamNamePrefix(input.LogStreamNamePref); err != nil {
		return nil, err
	}

	// The window members target Timestamp, whose @range declares a minimum
	// of 0 ("Valid Range: Minimum value of 0").
	if input.StartTime < 0 || input.EndTime < 0 {
		return nil, NewLogsError("InvalidParameterException",
			"startTime and endTime must be non-negative epoch milliseconds", 400)
	}

	// The pattern syntax validates on the read plane too: an invalid
	// pattern rejects with InvalidParameterException instead of matching
	// nothing, and the two-regex ceiling of delimited and JSON patterns
	// applies "when filtering log events or Live Tail".
	if err := validateFilterPattern(input.FilterPattern); err != nil {
		return nil, err
	}

	limit, err := validateListLimit(input.Limit, logsstore.DefaultEventsLimit, logsstore.MaxEventsLimit)
	if err != nil {
		return nil, err
	}

	if err := validateStartFromHeadDate(input.StartFromHead, input.StartTime); err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	logStreamNames := input.LogStreamNames
	if input.LogStreamNamePref != "" {
		prefixStreams, err := fetchAllLogStreams(store, input.LogGroupName, input.LogStreamNamePref)
		if err != nil {
			return nil, mapStoreError(err)
		}
		for _, ls := range prefixStreams {
			logStreamNames = append(logStreamNames, ls.Name)
		}
	}

	events, nextMarker, err := store.FilterLogEvents(
		input.LogGroupName, logStreamNames, input.StartTime, input.EndTime,
		input.FilterPattern, int(limit), input.StartFromHead, input.NextToken,
	)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &FilterLogEventsResult{
		Events:    events,
		NextToken: nextMarker,
	}, nil
}

// --- Retention policy ---

type PutRetentionPolicyInput struct {
	LogGroupName    string
	RetentionInDays int32
	Region          string
}

func (s *LogsService) putRetentionPolicyCore(input PutRetentionPolicyInput) error {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return err
	}
	if !logsstore.IsValidRetentionDays(input.RetentionInDays) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("%d is not a valid retention value. Allowed values: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653", input.RetentionInDays),
			400)
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	if err := store.MutateLogGroup(input.LogGroupName, func(lg *logsstore.LogGroup) error {
		lg.SetRetention(input.RetentionInDays)
		return nil
	}); err != nil {
		return mapStoreError(err)
	}
	return nil
}

type DeleteRetentionPolicyInput struct {
	LogGroupName string
	Region       string
}

func (s *LogsService) deleteRetentionPolicyCore(input DeleteRetentionPolicyInput) error {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	if err := store.MutateLogGroup(input.LogGroupName, func(lg *logsstore.LogGroup) error {
		lg.SetRetention(0)
		return nil
	}); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// --- KMS key association ---

type AssociateKmsKeyInput struct {
	LogGroupName       string
	ResourceIdentifier string
	KmsKeyId           string
	Region             string
}

func (s *LogsService) associateKmsKeyCore(input AssociateKmsKeyInput) error {
	if input.LogGroupName == "" && input.ResourceIdentifier == "" {
		return errRequiredMember("logGroupName or resourceIdentifier")
	}
	if input.KmsKeyId == "" {
		return errRequiredMember("kmsKeyId")
	}
	if err := validateKmsTargetExclusivity(input.LogGroupName, input.ResourceIdentifier); err != nil {
		return err
	}
	if err := validateKmsKeyId(input.KmsKeyId); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	target := logGroupTarget(input.LogGroupName, input.ResourceIdentifier)

	// The account-level query-result ARN form: "Specify the following ARN
	// to have future GetQueryResults operations in this account encrypt
	// the results with the specified KMS key" (resourceIdentifier member
	// documentation). The platform records the association at account
	// scope the same way it records a log group's key (ingested data is
	// stored under the platform's own at-rest protection either way;
	// GetQueryResults serves plaintext per the operation documentation).
	if isQueryResultResourceIdentifier(input.ResourceIdentifier) {
		if err := store.PutQueryResultKmsKey(input.KmsKeyId); err != nil {
			return mapStoreError(err)
		}
		return nil
	}

	// The log-group association documents the key's class and state as a
	// parameter error — the query-result form's documentation does not,
	// so the check binds the log-group path alone.
	if err := s.validateLogGroupKmsKey(input.KmsKeyId); err != nil {
		return err
	}

	if err := store.MutateLogGroup(target, func(lg *logsstore.LogGroup) error {
		lg.KmsKeyId = input.KmsKeyId
		return nil
	}); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// validateLogGroupKmsKey rejects the key classes and states the log-group
// association documents as InvalidParameterException: "If you attempt to
// associate a KMS key with the log group but the KMS key does not exist or
// the KMS key is disabled, you receive an InvalidParameterException error"
// and "CloudWatch Logs supports only symmetric KMS keys. Do not associate
// an asymmetric KMS key with your log group" (CreateLogGroup and
// AssociateKmsKey, identical sentences on both pages). The check needs the
// platform KMS invoker; without one configured the ARN form alone governs.
func (s *LogsService) validateLogGroupKmsKey(kmsKeyId string) error {
	if s.kmsInvoker() == nil {
		return nil
	}
	if !s.kmsInvoker().SymmetricEncryptionKeyUsable(context.Background(), kmsKeyId) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("KMS key %s does not exist, is disabled, or is not a symmetric encryption key; CloudWatch Logs supports only symmetric KMS keys", kmsKeyId), 400)
	}
	return nil
}

// validateKmsTargetExclusivity enforces the members' shared rule: "you
// must specify either the resourceIdentifier parameter or the logGroup
// parameter, but you can't specify both" (both member documentations).
func validateKmsTargetExclusivity(logGroupName, resourceIdentifier string) error {
	if logGroupName != "" && resourceIdentifier != "" {
		return NewLogsError("InvalidParameterException",
			"You must specify either the resourceIdentifier parameter or the logGroupName parameter, but you can't specify both", 400)
	}
	return nil
}

// isQueryResultResourceIdentifier reports whether the identifier is the
// account-level query-result ARN form the resourceIdentifier member
// documents: arn:aws:logs:REGION:ACCOUNT_ID:query-result:*. The resource
// segment is matched exactly; a malformed or foreign ARN falls through to
// the log-group resolution, which reports it as not found.
func isQueryResultResourceIdentifier(identifier string) bool {
	if !strings.HasPrefix(identifier, "arn:") {
		return false
	}
	parsed, err := arn.ParseARN(identifier)
	if err != nil {
		return false
	}
	return parsed.Service == "logs" && parsed.Resource == "query-result:*"
}

type DisassociateKmsKeyInput struct {
	LogGroupName       string
	ResourceIdentifier string
	Region             string
}

func (s *LogsService) disassociateKmsKeyCore(input DisassociateKmsKeyInput) error {
	if input.LogGroupName == "" && input.ResourceIdentifier == "" {
		return errRequiredMember("logGroupName or resourceIdentifier")
	}
	if err := validateKmsTargetExclusivity(input.LogGroupName, input.ResourceIdentifier); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	if isQueryResultResourceIdentifier(input.ResourceIdentifier) {
		if err := store.DeleteQueryResultKmsKey(); err != nil {
			return mapStoreError(err)
		}
		return nil
	}

	target := logGroupTarget(input.LogGroupName, input.ResourceIdentifier)

	if err := store.MutateLogGroup(target, func(lg *logsstore.LogGroup) error {
		lg.KmsKeyId = ""
		return nil
	}); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// --- Deletion protection ---

type PutLogGroupDeletionProtectionInput struct {
	LogGroupIdentifier string
	// DeletionProtectionEnabled is the modelled setting; the wire presence
	// flag carries the member's required trait (a bool's zero value is a
	// legitimate explicit false, so absence must travel separately).
	DeletionProtectionEnabled bool
	DeletionProtectionSet     bool
	Region                    string
}

func (s *LogsService) putLogGroupDeletionProtectionCore(input PutLogGroupDeletionProtectionInput) error {
	// deletionProtectionEnabled is required: an omitted member rejects
	// rather than silently writing false.
	if !input.DeletionProtectionSet {
		return NewLogsError("InvalidParameterException",
			"The deletionProtectionEnabled member is required", 400)
	}
	target := logGroupTarget("", input.LogGroupIdentifier)
	if target == "" {
		return errRequiredMember("logGroupIdentifier")
	}
	// The identifier member targets the shared LogGroupIdentifier shape:
	// an over-trait value is a parameter error, not a not-found.
	if err := validateLogGroupIdentifierElements([]string{input.LogGroupIdentifier}); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	if err := store.MutateLogGroup(target, func(lg *logsstore.LogGroup) error {
		lg.DeletionProtectionEnabled = input.DeletionProtectionEnabled
		return nil
	}); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// DescribeLogStreamsInput is the transport-agnostic input for listing log streams.
type DescribeLogStreamsInput struct {
	LogGroupName        string
	LogStreamNamePrefix string
	OrderBy             string
	Descending          bool
	NextToken           string
	Limit               int32
	Region              string
}

// DescribeLogStreamsResult holds the outcome of a successful log stream listing.
type DescribeLogStreamsResult struct {
	LogStreams []*logsstore.LogStream
	NextToken  string
}

// describeLogStreamsCore lists log streams with ordering and pagination.
// Every ordering path mints and consumes the one streamPageToken
// vocabulary: the token carries the request identity (group, prefix,
// ordering) and a value cursor, never a raw store key or a bare offset.
func (s *LogsService) describeLogStreamsCore(input DescribeLogStreamsInput) (*DescribeLogStreamsResult, error) {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return nil, err
	}
	// The prefix targets LogStreamName, so its pattern and length ceiling
	// apply; a prefix violating them is a parameter error, not a
	// ResourceNotFound.
	if err := validateLogStreamNamePrefix(input.LogStreamNamePrefix); err != nil {
		return nil, err
	}

	limit, err := validateListLimit(input.Limit, logsstore.DefaultDescribeLimit, logsstore.MaxDescribeLimit)
	if err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	// The operation declares ResourceNotFoundException: a listing against a
	// group that does not exist fails rather than serving an empty page
	// (the sibling event reads carry the same check).
	if _, err := store.GetLogGroup(input.LogGroupName); err != nil {
		return nil, mapStoreError(err)
	}

	// The model's default ordering is LogStreamName; a supplied value
	// outside the OrderBy enum rejects (the operation declares
	// InvalidParameterException).
	orderBy := input.OrderBy
	if orderBy == "" {
		orderBy = "LogStreamName"
	}
	if orderBy != "LogStreamName" && orderBy != "LastEventTime" {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid orderBy: %s. Valid values: LogStreamName, LastEventTime", input.OrderBy), 400)
	}

	var cursor streamPageToken
	if input.NextToken != "" {
		cursor, err = decodeStreamPageToken(input.NextToken)
		if err != nil {
			return nil, err
		}
		if cursor.Group != input.LogGroupName || cursor.Prefix != input.LogStreamNamePrefix ||
			cursor.OrderBy != orderBy || cursor.Descending != input.Descending {
			return nil, errInvalidPageToken
		}
	}

	if orderBy == "LastEventTime" {
		if input.LogStreamNamePrefix != "" {
			return nil, NewLogsError("InvalidParameterException",
				"Cannot specify logStreamNamePrefix when orderBy is LastEventTime", 400)
		}

		allStreams, err := fetchAllLogStreams(store, input.LogGroupName, "")
		if err != nil {
			return nil, mapStoreError(err)
		}

		sort.Slice(allStreams, func(i, j int) bool {
			if allStreams[i].LastEventTs != allStreams[j].LastEventTs {
				if input.Descending {
					return allStreams[i].LastEventTs > allStreams[j].LastEventTs
				}
				return allStreams[i].LastEventTs < allStreams[j].LastEventTs
			}
			// The name tiebreak keeps the order total (names are unique
			// within a group) and the cursor exact.
			if input.Descending {
				return allStreams[i].Name > allStreams[j].Name
			}
			return allStreams[i].Name < allStreams[j].Name
		})

		startIdx := 0
		if cursor.LastName != "" || cursor.LastEventTs != 0 {
			startIdx = len(allStreams)
			for i, ls := range allStreams {
				if streamSortsAfterCursor(ls, cursor, input.Descending) {
					startIdx = i
					break
				}
			}
		}

		endIdx := startIdx + int(limit)
		if endIdx > len(allStreams) {
			endIdx = len(allStreams)
		}
		if startIdx > len(allStreams) {
			startIdx = len(allStreams)
		}

		result := &DescribeLogStreamsResult{}
		if startIdx < endIdx {
			result.LogStreams = allStreams[startIdx:endIdx]
			last := allStreams[endIdx-1]
			if endIdx < len(allStreams) {
				result.NextToken, err = encodeStreamPageToken(streamPageToken{
					Group: input.LogGroupName, OrderBy: orderBy, Descending: input.Descending,
					LastName: last.Name, LastEventTs: last.LastEventTs,
				})
				if err != nil {
					return nil, err
				}
			}
		}
		return result, nil
	}

	if input.Descending {
		allStreams, err := fetchAllLogStreams(store, input.LogGroupName, input.LogStreamNamePrefix)
		if err != nil {
			return nil, mapStoreError(err)
		}

		sort.Slice(allStreams, func(i, j int) bool {
			return allStreams[i].Name > allStreams[j].Name
		})

		// Descending name order: resume at the first stream whose name
		// sorts strictly below the cursor's.
		startIdx := 0
		if cursor.LastName != "" {
			startIdx = len(allStreams)
			for i, ls := range allStreams {
				if ls.Name < cursor.LastName {
					startIdx = i
					break
				}
			}
		}

		endIdx := startIdx + int(limit)
		if endIdx > len(allStreams) {
			endIdx = len(allStreams)
		}

		result := &DescribeLogStreamsResult{}
		if startIdx < endIdx {
			result.LogStreams = allStreams[startIdx:endIdx]
			if endIdx < len(allStreams) {
				result.NextToken, err = encodeStreamPageToken(streamPageToken{
					Group: input.LogGroupName, Prefix: input.LogStreamNamePrefix,
					OrderBy: orderBy, Descending: true,
					LastName: allStreams[endIdx-1].Name,
				})
				if err != nil {
					return nil, err
				}
			}
		}
		return result, nil
	}

	// Ascending name order is the store's own scan order, so the cursor
	// stream name doubles as the scan marker and the page stays a single
	// store read.
	streams, nextMarker, err := store.ListLogStreams(input.LogGroupName, input.LogStreamNamePrefix, cursor.LastName, int(limit))
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := &DescribeLogStreamsResult{
		LogStreams: streams,
	}
	if nextMarker != "" && len(streams) > 0 {
		result.NextToken, err = encodeStreamPageToken(streamPageToken{
			Group: input.LogGroupName, Prefix: input.LogStreamNamePrefix, OrderBy: orderBy,
			LastName: streams[len(streams)-1].Name,
		})
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// streamSortsAfterCursor reports whether a stream's (LastEventTime, name)
// key sorts strictly after the cursor's key in the requested direction —
// the resume position of an event-time-ordered page.
func streamSortsAfterCursor(ls *logsstore.LogStream, cursor streamPageToken, descending bool) bool {
	if ls.LastEventTs != cursor.LastEventTs {
		if descending {
			return ls.LastEventTs < cursor.LastEventTs
		}
		return ls.LastEventTs > cursor.LastEventTs
	}
	if ls.Name != cursor.LastName {
		if descending {
			return ls.Name < cursor.LastName
		}
		return ls.Name > cursor.LastName
	}
	return false
}

// --- Legacy tag operations (TagLogGroup / UntagLogGroup / ListTagsLogGroup) ---
