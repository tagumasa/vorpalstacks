package cloudwatchlogs

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/scheduleexpr"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The scheduled-query family's Core layer: creation and update
// validation, the CRUD cores, the execution-history core and the
// listing core the operations delegate to.

// --- Scheduled Query Core ---

// CreateScheduledQueryInput holds the parsed parameters for
// CreateScheduledQuery.
type CreateScheduledQueryInput struct {
	Name                     string
	Description              string
	QueryString              string
	QueryLanguage            string
	LogGroupIdentifiers      []string
	ScheduleExpression       string
	State                    string
	ExecutionRoleArn         string
	Timezone                 string
	StartTimeOffset          int64
	EndTimeOffset            int64
	ScheduleStartTime        int64
	ScheduleEndTime          int64
	DestinationConfiguration map[string]interface{}
	Tags                     map[string]string
}

// validateScheduledQuerySpec applies the query-spec validation shared by
// CreateScheduledQuery and UpdateScheduledQuery so that an updated query can
// never hold a weaker specification than a newly created one.
func validateScheduledQuerySpec(queryLanguage, queryString, scheduleExpression, executionRoleArn, timezone string) error {
	if err := validateQueryString(queryString, 1); err != nil {
		return asScheduledQueryValidation(err)
	}
	if err := validateQueryPipeline(queryString); err != nil {
		return asScheduledQueryValidation(err)
	}
	// The member documents the cron grammar ("A cron expression that
	// defines when the scheduled query runs. The expression uses standard
	// cron syntax") and the family shares the cron/rate profile of the
	// expression package — the at() one-shot form belongs to the
	// EventBridge Scheduler profile alone.
	if !scheduleexpr.ValidateCronOrRateExpression(scheduleExpression) {
		return NewLogsError("ValidationException",
			fmt.Sprintf("Invalid schedule expression: %s. Must be a valid rate() or cron() expression", scheduleExpression), 400)
	}
	// "The timezone for evaluating the schedule expression" — the value
	// names an IANA time zone, the LoadLocation vocabulary; a value that
	// resolves to no zone cannot evaluate anything, and the Go runtime
	// alias "Local" names the host clock rather than a database entry.
	if timezone != "" {
		if timezone == "Local" {
			return NewLogsError("ValidationException",
				"timezone must name an IANA time zone; \"Local\" is not an IANA zone name", 400)
		}
		if _, err := time.LoadLocation(timezone); err != nil {
			return NewLogsError("ValidationException",
				fmt.Sprintf("Invalid timezone: %s. Must be an IANA time zone name", timezone), 400)
		}
	}
	// The scheduled-query operations always carry a queryLanguage (both
	// cores reject an absent member before reaching here), so the spec
	// validator judges the language on its own enum without the
	// optional-member tolerance StartQuery's validator applies.
	if queryLanguage == "" || !validQueryLanguages[queryLanguage] {
		return NewLogsError("ValidationException",
			fmt.Sprintf("Invalid queryLanguage: %s. Allowed values: CWLI, SQL, PPL", queryLanguage), 400)
	}
	return asScheduledQueryValidation(validateIAMRoleArn(executionRoleArn))
}

// asScheduledQueryValidation normalises a rejection from the shared
// validators into the identity the scheduled-query operations declare.
// CreateScheduledQuery and UpdateScheduledQuery carry ValidationException
// (with AccessDenied, Conflict, InternalServer, ResourceNotFound and
// Throttling) — neither InvalidParameterException nor
// MalformedQueryException, the latter being StartQuery's compile-error
// identity that the shared query pipeline speaks.
func asScheduledQueryValidation(err error) error {
	var ae *awserrors.AWSError
	if errors.As(err, &ae) {
		switch ae.GetCode() {
		case "InvalidParameterException", "MalformedQueryException":
			return NewLogsError("ValidationException", ae.GetMessage(), 400)
		}
	}
	return err
}

// createScheduledQueryCore validates input and persists a new scheduled query.
func (s *LogsService) createScheduledQueryCore(store *logsstore.Store, input *CreateScheduledQueryInput) (*logsstore.ScheduledQuery, error) {
	if err := asScheduledQueryValidation(validateScheduledQueryName(input.Name)); err != nil {
		return nil, err
	}
	// The model marks queryString, queryLanguage, scheduleExpression and
	// executionRoleArn required on the create request — the same members
	// the full-specification update path enforces.
	if input.QueryString == "" || input.ScheduleExpression == "" || input.QueryLanguage == "" || input.ExecutionRoleArn == "" {
		return nil, errValidationMember("queryString, scheduleExpression, queryLanguage and executionRoleArn")
	}
	if err := validateScheduledQuerySpec(input.QueryLanguage, input.QueryString, input.ScheduleExpression, input.ExecutionRoleArn, input.Timezone); err != nil {
		return nil, err
	}
	// ScheduledQueryLogGroupIdentifiers carries @length {1,50} ("Array
	// Members: Minimum number of 1 item."): the optional member cannot
	// arrive as the empty list — a present-but-empty parse (marked by the
	// non-nil empty slice) rejects instead of silently meaning "absent".
	if input.LogGroupIdentifiers != nil && len(input.LogGroupIdentifiers) == 0 {
		return nil, NewLogsError("ValidationException",
			"logGroupIdentifiers must contain between 1 and 50 entries", 400)
	}
	if err := asScheduledQueryValidation(validateLogGroupIdentifierCount(input.LogGroupIdentifiers)); err != nil {
		return nil, err
	}
	if err := asScheduledQueryValidation(validateLogGroupIdentifierElements(input.LogGroupIdentifiers)); err != nil {
		return nil, err
	}
	if err := asScheduledQueryValidation(validateScheduledQueryDescription(input.Description)); err != nil {
		return nil, err
	}
	// The shared Tags shape's length trait is 1-50 and every entry carries
	// the TagKey/TagValue traits; both reject through the family's
	// ValidationException identity.
	if len(input.Tags) > logsstore.MaxScheduledQueryTags {
		return nil, NewLogsError("ValidationException",
			fmt.Sprintf("A scheduled query accepts at most %d tags", logsstore.MaxScheduledQueryTags), 400)
	}
	if err := asScheduledQueryValidation(validateTagEntries(input.Tags)); err != nil {
		return nil, err
	}
	// The schedule window members carry "Valid Range: Minimum value of 0";
	// the zero value is the members' absent form, so only negatives reject.
	if input.ScheduleStartTime < 0 {
		return nil, NewLogsError("ValidationException", "scheduleStartTime must not be negative", 400)
	}
	if input.ScheduleEndTime < 0 {
		return nil, NewLogsError("ValidationException", "scheduleEndTime must not be negative", 400)
	}
	if input.DestinationConfiguration != nil {
		if err := validateDestinationConfiguration(input.DestinationConfiguration); err != nil {
			return nil, err
		}
	}

	state := input.State
	if state == "" {
		state = logsstore.ScheduledQueryStateEnabled
	}
	if !validateScheduledQueryState(state) {
		return nil, NewLogsError("ValidationException",
			fmt.Sprintf("Invalid state: %s. Allowed values: ENABLED, DISABLED", state), 400)
	}

	// "The name of the scheduled query. The name must be unique within
	// your account and region" — a create naming an existing scheduled
	// query rejects with the operation's declared ConflictException
	// ("This operation attempted to create a resource that already
	// exists").
	existing, err := store.ListScheduledQueries("")
	if err != nil {
		return nil, mapStoreError(err)
	}
	for _, sq := range existing {
		if sq.Name == input.Name {
			return nil, NewLogsError("ConflictException",
				fmt.Sprintf("A scheduled query named %s already exists", input.Name), 400)
		}
	}

	id := fmt.Sprintf("sq-%d", time.Now().UnixNano())

	sq := &logsstore.ScheduledQuery{
		Id:                       id,
		Name:                     input.Name,
		Description:              input.Description,
		QueryString:              input.QueryString,
		QueryLanguage:            input.QueryLanguage,
		LogGroupIdentifiers:      input.LogGroupIdentifiers,
		ScheduleExpression:       input.ScheduleExpression,
		ScheduleType:             "CUSTOMER_MANAGED",
		State:                    state,
		ExecutionRoleArn:         input.ExecutionRoleArn,
		Timezone:                 input.Timezone,
		StartTimeOffset:          input.StartTimeOffset,
		EndTimeOffset:            input.EndTimeOffset,
		ScheduleStartTime:        input.ScheduleStartTime,
		ScheduleEndTime:          input.ScheduleEndTime,
		DestinationConfiguration: input.DestinationConfiguration,
		CreationTime:             time.Now().UTC().UnixMilli(),
		Tags:                     input.Tags,
	}

	if err := store.PutScheduledQuery(sq); err != nil {
		return nil, mapStoreError(err)
	}
	return sq, nil
}

func (s *LogsService) deleteScheduledQueryCore(store *logsstore.Store, identifier string) error {
	if identifier == "" {
		return errValidationMember("identifier")
	}
	id, err := resolveScheduledQueryRef(store, identifier)
	if err != nil {
		return err
	}
	if err := store.DeleteScheduledQuery(id); err != nil {
		return mapStoreError(err)
	}
	return nil
}

func (s *LogsService) getScheduledQueryCore(store *logsstore.Store, identifier string) (*logsstore.ScheduledQuery, error) {
	if identifier == "" {
		return nil, errValidationMember("identifier")
	}
	id, err := resolveScheduledQueryRef(store, identifier)
	if err != nil {
		return nil, err
	}
	sq, err := store.GetScheduledQuery(id)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return sq, nil
}

// UpdateScheduledQueryInput holds the parsed parameters for
// UpdateScheduledQuery. The Smithy model marks identifier, queryLanguage,
// queryString, scheduleExpression and executionRoleArn as required because
// the update replaces the full specification; the remaining members are
// replaced only when provided.
type UpdateScheduledQueryInput struct {
	Identifier               string
	Description              string
	QueryString              string
	QueryLanguage            string
	LogGroupIdentifiers      []string
	ScheduleExpression       string
	State                    string
	ExecutionRoleArn         string
	Timezone                 string
	StartTimeOffset          *int64
	EndTimeOffset            *int64
	ScheduleStartTime        *int64
	ScheduleEndTime          *int64
	DestinationConfiguration map[string]interface{}
}

// updateScheduledQueryCore validates input and applies the full-specification
// update semantics to the stored scheduled query.
func (s *LogsService) updateScheduledQueryCore(store *logsstore.Store, input *UpdateScheduledQueryInput) (*logsstore.ScheduledQuery, error) {
	if input.Identifier == "" {
		return nil, errValidationMember("identifier")
	}
	if input.QueryLanguage == "" || input.QueryString == "" || input.ScheduleExpression == "" || input.ExecutionRoleArn == "" {
		return nil, errValidationMember("queryLanguage, queryString, scheduleExpression and executionRoleArn")
	}
	if err := validateScheduledQuerySpec(input.QueryLanguage, input.QueryString, input.ScheduleExpression, input.ExecutionRoleArn, input.Timezone); err != nil {
		return nil, err
	}
	if input.State != "" && !validateScheduledQueryState(input.State) {
		return nil, NewLogsError("ValidationException",
			fmt.Sprintf("Invalid state: %s. Allowed values: ENABLED, DISABLED", input.State), 400)
	}
	if input.LogGroupIdentifiers != nil {
		// The same @length {1,50} trait the create path enforces: an
		// explicitly empty update list rejects instead of clearing the
		// query's group set.
		if len(input.LogGroupIdentifiers) == 0 {
			return nil, NewLogsError("ValidationException",
				"logGroupIdentifiers must contain between 1 and 50 entries", 400)
		}
		if err := asScheduledQueryValidation(validateLogGroupIdentifierCount(input.LogGroupIdentifiers)); err != nil {
			return nil, err
		}
	}
	if input.DestinationConfiguration != nil {
		if err := validateDestinationConfiguration(input.DestinationConfiguration); err != nil {
			return nil, err
		}
	}

	id, err := resolveScheduledQueryRef(store, input.Identifier)
	if err != nil {
		return nil, err
	}

	// The whole read-modify-write runs under the store's record lock so
	// a concurrent delivery touch or delete cannot interleave.
	err = store.MutateScheduledQuery(id, func(sq *logsstore.ScheduledQuery) error {
		sq.QueryString = input.QueryString
		sq.QueryLanguage = input.QueryLanguage
		sq.ScheduleExpression = input.ScheduleExpression
		sq.ExecutionRoleArn = input.ExecutionRoleArn
		if input.Description != "" {
			sq.Description = input.Description
		}
		if input.State != "" {
			sq.State = input.State
		}
		if input.LogGroupIdentifiers != nil {
			sq.LogGroupIdentifiers = input.LogGroupIdentifiers
		}
		if input.Timezone != "" {
			sq.Timezone = input.Timezone
		}
		if input.StartTimeOffset != nil {
			sq.StartTimeOffset = *input.StartTimeOffset
		}
		if input.EndTimeOffset != nil {
			sq.EndTimeOffset = *input.EndTimeOffset
		}
		if input.ScheduleStartTime != nil {
			sq.ScheduleStartTime = *input.ScheduleStartTime
		}
		if input.ScheduleEndTime != nil {
			sq.ScheduleEndTime = *input.ScheduleEndTime
		}
		if input.DestinationConfiguration != nil {
			sq.DestinationConfiguration = input.DestinationConfiguration
		}
		return nil
	})
	if err != nil {
		return nil, mapStoreError(err)
	}

	updated, err := store.GetScheduledQuery(id)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return updated, nil
}

// ScheduledQueryHistoryResult carries one page of scheduled query executions
// plus the scheduled query itself for response formatting.
type ScheduledQueryHistoryResult struct {
	// Id is the resolved id-keyed identifier of the queried scheduled
	// query: the Identifier member arrives ARN-or-name and the handler
	// mints the response ARN from the resolved id.
	Id         string
	Executions []*logsstore.ScheduledQueryExecution
	NextMarker string
	Query      *logsstore.ScheduledQuery
}

// GetScheduledQueryHistoryInput holds the parsed parameters for
// GetScheduledQueryHistory. StartTimeSet and EndTimeSet carry the wire
// presence of the two required members — an int64 zero value cannot
// distinguish an explicit zero (legal, the members' minimum) from an
// omitted member.
type GetScheduledQueryHistoryInput struct {
	Identifier        string
	StartTime         int64
	EndTime           int64
	StartTimeSet      bool
	EndTimeSet        bool
	ExecutionStatuses []string
	NextToken         string
	MaxResults        int32
}

// getScheduledQueryHistoryCore validates input and returns a page of the
// scheduled query's execution history.
func (s *LogsService) getScheduledQueryHistoryCore(store *logsstore.Store, input *GetScheduledQueryHistoryInput) (*ScheduledQueryHistoryResult, error) {
	if input.Identifier == "" {
		return nil, errValidationMember("identifier")
	}
	maxResults, err := validateListLimitValidation(input.MaxResults, logsstore.DefaultListMaxResults, logsstore.MaxListMaxResults)
	if err != nil {
		return nil, err
	}
	input.MaxResults = maxResults

	// executionStatuses filters by the ExecutionStatus enumeration; a value
	// outside the vocabulary rejects with ValidationException — the
	// validation error this operation declares (its error list carries
	// ValidationException where the service's other operations declare
	// InvalidParameterException).
	for _, st := range input.ExecutionStatuses {
		if !validExecutionStatus[st] {
			return nil, NewLogsError("ValidationException",
				fmt.Sprintf("Invalid executionStatus: %s. Valid values: Running, InvalidQuery, Complete, Failed, Timeout", st), 400)
		}
	}

	// Both window members are required on the request ("Required: Yes")
	// and carry "Valid Range: Minimum value of 0"; an omitting or negative
	// request rejects rather than silently serving the unbounded history.
	if !input.StartTimeSet || !input.EndTimeSet {
		return nil, errValidationMember("startTime and endTime")
	}
	if input.StartTime < 0 || input.EndTime < 0 {
		return nil, NewLogsError("ValidationException", "startTime and endTime must not be negative", 400)
	}

	id, err := resolveScheduledQueryRef(store, input.Identifier)
	if err != nil {
		return nil, err
	}
	// The history serves an existing scheduled query: the operation
	// declares ResourceNotFoundException and an unknown identifier reaches
	// it here rather than synthesising an empty history.
	sq, err := store.GetScheduledQuery(id)
	if err != nil {
		if errors.Is(err, logsstore.ErrResourceNotFound) {
			return nil, NewLogsError("ResourceNotFoundException",
				fmt.Sprintf("The specified scheduled query does not exist: %s", input.Identifier), 400)
		}
		return nil, mapStoreError(err)
	}

	allExecs, err := store.ListScheduledQueryExecutions(id, input.StartTime, input.EndTime)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if len(input.ExecutionStatuses) > 0 {
		statusSet := make(map[string]bool)
		for _, st := range input.ExecutionStatuses {
			statusSet[st] = true
		}
		var filtered []*logsstore.ScheduledQueryExecution
		for _, exec := range allExecs {
			mappedStatus := mapExecutionStatus(exec.Status)
			if statusSet[mappedStatus] {
				filtered = append(filtered, exec)
			}
		}
		allExecs = filtered
	}

	// The page marker rides the scoped listing vocabulary: the token
	// carries the request identity and its documented expiry rather than
	// a bare trigger-time key.
	statusKey := ""
	if len(input.ExecutionStatuses) > 0 {
		sorted := make([]string, len(input.ExecutionStatuses))
		copy(sorted, input.ExecutionStatuses)
		sort.Strings(sorted)
		statusKey = strings.Join(sorted, "\x1f")
	}
	scope := listingScope("scheduledqueryexecutions", id, statusKey,
		strconv.FormatInt(input.StartTime, 10), strconv.FormatInt(input.EndTime, 10))
	result, err := paginateScopedListing(scope, input.NextToken, allExecs, int(input.MaxResults), func(e *logsstore.ScheduledQueryExecution) string {
		return strconv.FormatInt(e.TriggerTime, 10)
	})
	if err != nil {
		return nil, err
	}

	history := &ScheduledQueryHistoryResult{
		Id:         id,
		Executions: result.Items,
		NextMarker: result.NextMarker,
		Query:      sq,
	}
	return history, nil
}

// listScheduledQueriesCore validates input and returns a page of scheduled
// queries filtered by state and schedule type.
func (s *LogsService) listScheduledQueriesCore(store *logsstore.Store, stateFilter, scheduleTypeFilter string, maxResults int32, nextToken string) ([]*logsstore.ScheduledQuery, string, error) {
	maxResults, err := validateListLimitValidation(maxResults, logsstore.DefaultListMaxResults, logsstore.MaxListMaxResults)
	if err != nil {
		return nil, "", err
	}
	// Both filter members are enumerations (state "ENABLED | DISABLED",
	// scheduleType "CUSTOMER_MANAGED | AWS_MANAGED"); a value outside the
	// vocabulary rejects with the family's ValidationException instead of
	// silently matching nothing.
	if stateFilter != "" && stateFilter != logsstore.ScheduledQueryStateEnabled && stateFilter != logsstore.ScheduledQueryStateDisabled {
		return nil, "", NewLogsError("ValidationException",
			fmt.Sprintf("Invalid state: %s. Valid values: ENABLED, DISABLED", stateFilter), 400)
	}
	if scheduleTypeFilter != "" && scheduleTypeFilter != "CUSTOMER_MANAGED" && scheduleTypeFilter != "AWS_MANAGED" {
		return nil, "", NewLogsError("ValidationException",
			fmt.Sprintf("Invalid scheduleType: %s. Valid values: CUSTOMER_MANAGED, AWS_MANAGED", scheduleTypeFilter), 400)
	}
	all, err := store.ListScheduledQueries(stateFilter)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	if scheduleTypeFilter != "" {
		var filtered []*logsstore.ScheduledQuery
		for _, sq := range all {
			if sq.ScheduleType == scheduleTypeFilter {
				filtered = append(filtered, sq)
			}
		}
		all = filtered
	}

	// The page marker rides the scoped listing vocabulary: the token
	// carries the request identity and its documented expiry rather than
	// a bare query id.
	scope := listingScope("scheduledqueries", stateFilter, scheduleTypeFilter)
	result, err := paginateScopedListing(scope, nextToken, all, int(maxResults), func(sq *logsstore.ScheduledQuery) string {
		return sq.Id
	})
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// getLogGroupFieldsCore validates input and retrieves field names from
// recent log events in the specified log group. The request's time member
// is epoch SECONDS (the model documentation's unit) and centres a
// ±8-minute window; when it is omitted the most recent 15 minutes up to
// the current time are searched.
