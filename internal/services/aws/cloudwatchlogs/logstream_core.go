package cloudwatchlogs

import (
	"fmt"

	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
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

type PutLogEventsInput struct {
	LogGroupName  string
	LogStreamName string
	Events        []logsstore.LogEntry
	Region        string
}

type PutLogEventsResult struct {
	NextSequenceToken string
	RejectedLogEvents map[string]interface{}
}

func (s *LogsService) putLogEventsCore(input PutLogEventsInput) (*PutLogEventsResult, error) {
	if input.LogGroupName == "" || input.LogStreamName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	if len(input.Events) == 0 {
		return nil, ErrMissingParameter
	}

	if len(input.Events) > logsstore.MaxChunkSize {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Maximum number of log events in a single batch is %d", logsstore.MaxChunkSize), 400)
	}

	validEvents, rejectedInfo, valErr := validateLogEvents(input.Events)
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

	eventsCopy := make([]logsstore.LogEntry, len(validEvents))
	copy(eventsCopy, validEvents)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("Panic in async metric/subscription processing",
					logs.String("logGroup", input.LogGroupName),
					logs.Any("panic", r))
			}
		}()
		s.evaluateMetricFilters(store, input.Region, input.LogGroupName, eventsCopy)
		s.deliverSubscriptionEvents(store, input.Region, input.LogGroupName, input.LogStreamName, eventsCopy)
	}()

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
	if input.LogGroupName == "" || input.LogStreamName == "" {
		return nil, ErrMissingParameter
	}

	limit, err := validateListLimit(input.Limit, 10000, 10000)
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
	Events          []*logsstore.OutputLogEvent
	SearchedStreams map[string]bool
	NextToken       string
}

func (s *LogsService) filterLogEventsCore(input FilterLogEventsInput) (*FilterLogEventsResult, error) {
	logGroupName := input.LogGroupName
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	if len(input.LogStreamNames) > 0 && input.LogStreamNamePref != "" {
		return nil, NewLogsError("InvalidParameterException",
			"Cannot specify both logStreamNames and logStreamNamePrefix", 400)
	}

	limit, err := validateListLimit(input.Limit, 10000, 10000)
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
		prefixStreams, err := fetchAllLogStreams(store, logGroupName, input.LogStreamNamePref)
		if err != nil {
			return nil, mapStoreError(err)
		}
		for _, ls := range prefixStreams {
			logStreamNames = append(logStreamNames, ls.Name)
		}
	}

	events, searchedStreams, nextMarker, err := store.FilterLogEvents(
		logGroupName, logStreamNames, input.StartTime, input.EndTime,
		input.FilterPattern, int(limit), input.StartFromHead, input.NextToken,
	)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &FilterLogEventsResult{
		Events:          events,
		SearchedStreams: searchedStreams,
		NextToken:       nextMarker,
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

	lg, err := store.GetLogGroup(input.LogGroupName)
	if err != nil {
		return mapStoreError(err)
	}

	lg.SetRetention(input.RetentionInDays)
	if err := store.PutLogGroup(lg); err != nil {
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

	lg, err := store.GetLogGroup(input.LogGroupName)
	if err != nil {
		return mapStoreError(err)
	}

	lg.SetRetention(0)
	if err := store.PutLogGroup(lg); err != nil {
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
		return ErrMissingParameter
	}
	if input.KmsKeyId == "" {
		return ErrMissingParameter
	}
	if err := validateKmsKeyId(input.KmsKeyId); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	target := input.LogGroupName
	if target == "" {
		target = input.ResourceIdentifier
	}

	lg, err := store.GetLogGroup(target)
	if err != nil {
		return mapStoreError(err)
	}

	lg.KmsKeyId = input.KmsKeyId
	if err := store.PutLogGroup(lg); err != nil {
		return mapStoreError(err)
	}
	return nil
}

type DisassociateKmsKeyInput struct {
	LogGroupName       string
	ResourceIdentifier string
	Region             string
}

func (s *LogsService) disassociateKmsKeyCore(input DisassociateKmsKeyInput) error {
	if input.LogGroupName == "" && input.ResourceIdentifier == "" {
		return ErrMissingParameter
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	target := input.LogGroupName
	if target == "" {
		target = input.ResourceIdentifier
	}

	lg, err := store.GetLogGroup(target)
	if err != nil {
		return mapStoreError(err)
	}

	lg.KmsKeyId = ""
	if err := store.PutLogGroup(lg); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// --- Deletion protection ---

type PutLogGroupDeletionProtectionInput struct {
	LogGroupIdentifier        string
	LogGroupName              string
	DeletionProtectionEnabled bool
	Region                    string
}

func (s *LogsService) putLogGroupDeletionProtectionCore(input PutLogGroupDeletionProtectionInput) error {
	target := input.LogGroupName
	if target == "" {
		target = input.LogGroupIdentifier
	}
	if target == "" {
		return ErrMissingParameter
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	lg, err := store.GetLogGroup(target)
	if err != nil {
		return mapStoreError(err)
	}

	lg.DeletionProtectionEnabled = input.DeletionProtectionEnabled
	if err := store.PutLogGroup(lg); err != nil {
		return mapStoreError(err)
	}
	return nil
}
