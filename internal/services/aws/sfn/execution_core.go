package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"vorpalstacks/internal/core/logs"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// This file holds the execution-operation Core path shared by the HTTP API,
// the event-bus start path and the admin console, plus the qualified-ARN
// resolution both start and describe rely on.

// stateMachineReference is the resolution of a possibly qualified state
// machine ARN. A qualified ARN names a version
// (stateMachine:<name>:<number>) or an alias
// (stateMachine:<name>:<aliasName>); an unqualified ARN names the state
// machine's latest revision.
type stateMachineReference struct {
	// StateMachine is the base state machine record.
	StateMachine *sfnstore.StateMachine
	// Version is the version an execution must run: set when the ARN named
	// a version directly or an alias routed to one. AWS associates
	// executions started with a version or alias ARN with that version.
	Version *sfnstore.StateMachineVersion
	// Alias is set when the ARN named an alias.
	Alias *sfnstore.StateMachineAlias
}

// definition returns the definition an execution must run: the pinned
// version snapshot when the reference is version-qualified, the live
// state machine definition otherwise.
func (r *stateMachineReference) definition() string {
	if r.Version != nil {
		return r.Version.Definition
	}
	return r.StateMachine.Definition
}

// resolveStateMachineReference resolves a state machine ARN that may carry
// a version or alias qualifier. The Distributed Map label form
// (stateMachine:<name>/<label>) is rejected with ValidationException per
// the StartExecution documentation. Malformed ARNs return InvalidArn;
// well-formed ARNs whose state machine, version or alias does not exist
// return StateMachineDoesNotExist.
func resolveStateMachineReference(ctx context.Context, store *sfnstore.StepFunctionStore, arn string) (*stateMachineReference, error) {
	if err := validateArnRequired(arn, "stateMachineArn"); err != nil {
		return nil, err
	}

	_, service, _, _, resource := svcarn.SplitARN(arn)
	if service != "states" {
		return nil, NewInvalidArnException("stateMachineArn is not a States ARN: " + arn)
	}
	rest, ok := strings.CutPrefix(resource, "stateMachine:")
	if !ok {
		return nil, NewInvalidArnException("stateMachineArn is not a state machine ARN: " + arn)
	}
	if strings.Contains(rest, "/") {
		// A slash-qualified ARN names a Distributed Map state within a
		// state machine, which is not a startable or describable target.
		return nil, NewValidationException("stateMachineArn must not refer to a Distributed Map state: " + arn)
	}

	name, qualifier, qualified := strings.Cut(rest, ":")
	if name == "" {
		return nil, NewInvalidArnException("stateMachineArn is not a state machine ARN: " + arn)
	}

	if !qualified {
		sm, err := store.GetStateMachine(ctx, arn)
		if err != nil {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
		}
		return &stateMachineReference{StateMachine: sm}, nil
	}

	if _, err := strconv.Atoi(qualifier); err == nil {
		// Numeric qualifier: a version ARN.
		version, err := store.GetStateMachineVersion(ctx, arn)
		if err != nil {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
		}
		sm, err := store.GetStateMachine(ctx, version.StateMachineArn)
		if err != nil {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + version.StateMachineArn)
		}
		return &stateMachineReference{StateMachine: sm, Version: version}, nil
	}

	// Non-numeric qualifier: an alias ARN.
	alias, err := store.GetStateMachineAlias(ctx, arn)
	if err != nil {
		return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
	}
	sm, err := store.GetStateMachine(ctx, alias.StateMachineArn)
	if err != nil {
		return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + alias.StateMachineArn)
	}

	versionArn, err := selectVersionByWeight(alias.RoutingConfiguration)
	if err != nil {
		return nil, err
	}
	version, err := store.GetStateMachineVersion(ctx, versionArn)
	if err != nil {
		return nil, NewStateMachineDoesNotExist("State Machine Version Does not exist: " + versionArn)
	}
	return &stateMachineReference{StateMachine: sm, Version: version, Alias: alias}, nil
}

// selectVersionByWeight picks the version an alias-routed execution runs.
// Step Functions randomly chooses among the routing configuration entries
// based on the traffic percentage assigned to each version (state machine
// alias documentation); a single entry is deterministic.
func selectVersionByWeight(rc []sfnstore.RoutingConfiguration) (string, error) {
	if len(rc) == 0 {
		return "", NewValidationException("alias routing configuration is empty")
	}
	if len(rc) == 1 {
		return rc[0].StateMachineVersionArn, nil
	}

	pick := rand.IntN(100) + 1
	cumulative := 0
	for _, entry := range rc {
		cumulative += int(entry.Weight)
		if pick <= cumulative {
			return entry.StateMachineVersionArn, nil
		}
	}
	return rc[len(rc)-1].StateMachineVersionArn, nil
}

// executionAlreadyExistsIdempotent implements the documented StartExecution
// idempotency for STANDARD workflows: the same name with the same input on
// a running execution returns that execution instead of an error. The
// returned flag reports whether the idempotent path applied.
func executionAlreadyExistsIdempotent(store *sfnstore.StepFunctionStore, executionArn, input string) (*sfnstore.Execution, bool, error) {
	existing, err := store.GetExecution(context.Background(), executionArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrExecutionNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if existing.Status == "RUNNING" && existing.Input == input {
		return existing, true, nil
	}
	return nil, false, fmt.Errorf("execution %s already exists with different input or a closed status", executionArn)
}

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// StartExecutionInput carries every field that StartExecution needs.
type StartExecutionInput struct {
	StateMachineArn string
	Name            string
	Input           string
	TraceHeader     string
}

// StartExecutionResult carries the StartExecutionOutput members: the
// operation returns the execution ARN and start date only.
type StartExecutionResult struct {
	ExecutionArn string
	StartDate    time.Time
}

// StartSyncExecutionInput carries every field that StartSyncExecution needs.
type StartSyncExecutionInput struct {
	StateMachineArn string
	Name            string
	Input           string
	TraceHeader     string
	IncludedData    string
}

// StopExecutionInput carries every field that StopExecution needs.
type StopExecutionInput struct {
	ExecutionArn string
	Error        string
	Cause        string
}

// DescribeExecutionInput carries the parameters for DescribeExecution.
type DescribeExecutionInput struct {
	ExecutionArn string
	IncludedData string
}

// ListExecutionsInput carries the parameters for ListExecutions.
type ListExecutionsInput struct {
	StateMachineArn string
	StatusFilter    string
	MapRunArn       string
	RedriveFilter   string
	MaxResults      int32
	NextToken       string
}

// ListExecutionsResult carries the paginated ListExecutions output.
type ListExecutionsResult struct {
	Executions []*sfnstore.Execution
	NextToken  string
}

// GetExecutionHistoryInput carries the parameters for GetExecutionHistory.
type GetExecutionHistoryInput struct {
	ExecutionArn         string
	MaxResults           int32
	NextToken            string
	IncludeExecutionData bool
	ReverseOrder         bool
}

// RedriveExecutionInput carries the parameters for RedriveExecution.
type RedriveExecutionInput struct {
	ExecutionArn string
	ClientToken  string
}

// RedriveExecutionResult carries the RedriveExecutionOutput member.
type RedriveExecutionResult struct {
	RedriveDate time.Time
}

// ---------------------------------------------------------------------------
// Shared validation helpers
// ---------------------------------------------------------------------------

// validateExecutionInputData enforces the SensitiveData contract on
// execution input: non-empty input must be valid JSON within the
// 262,144-byte UTF-8 bound. AWS rejects malformed input with
// InvalidExecutionInput.
func validateExecutionInputData(input string) error {
	if input == "" {
		return nil
	}
	if len(input) > sfnstore.MaxExecutionDataBytes {
		return NewInvalidExecutionInput(fmt.Sprintf("Invalid State Machine Execution Input: input must be at most %d bytes, got %d", sfnstore.MaxExecutionDataBytes, len(input)))
	}
	if !json.Valid([]byte(input)) {
		return NewInvalidExecutionInput("Invalid State Machine Execution Input: input must be valid JSON")
	}
	return nil
}

// validateTraceHeader enforces the TraceHeader contract: at most 256
// ASCII characters.
func validateTraceHeader(traceHeader string) error {
	if len(traceHeader) > sfnstore.MaxTraceHeaderLength {
		return NewValidationException(fmt.Sprintf("traceHeader must be at most %d characters, got %d", sfnstore.MaxTraceHeaderLength, len(traceHeader)))
	}
	for i := 0; i < len(traceHeader); i++ {
		if traceHeader[i] > unicode.MaxASCII {
			return NewValidationException("traceHeader must contain only ASCII characters")
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// StartExecution / StartSyncExecution Core
// ---------------------------------------------------------------------------

// startExecutionCore is the single entry point for starting an execution:
// parameter validation, qualified-ARN resolution, persistence and the
// asynchronous launch all run here, so the HTTP handler and the event-bus
// start path behave identically. The response carries the execution ARN
// and start date only, per the StartExecutionOutput shape.
func (s *StepFunctionService) startExecutionCore(ctx context.Context, store *sfnstore.StepFunctionStore, in StartExecutionInput) (*StartExecutionResult, error) {
	if err := validateArnRequired(in.StateMachineArn, "stateMachineArn"); err != nil {
		return nil, err
	}
	if err := validateExecutionInputData(in.Input); err != nil {
		return nil, err
	}
	if err := validateTraceHeader(in.TraceHeader); err != nil {
		return nil, err
	}

	name := in.Name
	if name == "" {
		name = generateExecutionName()
	} else if err := validateExecutionName(name); err != nil {
		return nil, err
	}

	ref, err := resolveStateMachineReference(ctx, store, in.StateMachineArn)
	if err != nil {
		return nil, err
	}
	sm := ref.StateMachine

	executionArn := svcarn.NewARNBuilder(s.accountID, store.GetRegion()).StepFunctions().
		Execution(svcarn.ExtractStateMachineNameFromARN(sm.StateMachineArn), name)

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, in.Input, in.TraceHeader)
	exec.ExecutionArn = executionArn
	exec.StateMachineVersionArn = ref.versionArn()
	exec.StateMachineAliasArn = ref.aliasArn()

	if err := store.CreateExecution(ctx, exec); err != nil {
		if errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			// StartExecution is idempotent for STANDARD workflows: the
			// same name and input on a running execution returns that
			// execution; EXPRESS workflows are not idempotent.
			if sm.Type != "EXPRESS" {
				if existing, idempotent, ierr := executionAlreadyExistsIdempotent(store, executionArn, in.Input); ierr == nil && idempotent {
					return &StartExecutionResult{ExecutionArn: existing.ExecutionArn, StartDate: existing.StartDate}, nil
				}
			}
			return nil, NewExecutionAlreadyExists("An execution with the same name already exists: " + executionArn)
		}
		return nil, err
	}

	s.launchExecution(store, exec)

	return &StartExecutionResult{ExecutionArn: exec.ExecutionArn, StartDate: exec.StartDate}, nil
}

// launchExecution runs an execution asynchronously with panic isolation,
// registration for StopExecution cancellation and terminal-state
// persistence on panic.
func (s *StepFunctionService) launchExecution(store *sfnstore.StepFunctionStore, exec *sfnstore.Execution) {
	executor := NewExecutorWithStores(store, s.bus, s.accountID, store.GetRegion())
	execCtx, cancel := context.WithCancel(context.Background())
	store.RegisterExecution(exec.ExecutionArn, cancel)
	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer store.UnregisterExecution(exec.ExecutionArn)
		defer func() {
			if r := recover(); r != nil {
				logs.Error("sfn: panic in execution", logs.String("arn", exec.ExecutionArn), logs.Any("panic", r))
				exec.Status = "FAILED"
				exec.Error = "States.InternalError"
				exec.Cause = fmt.Sprintf("internal panic: %v", r)
				exec.StopDate = time.Now().UTC()
				_ = store.UpdateExecution(context.Background(), exec)
			}
		}()
		if err := executor.ExecuteStateMachine(execCtx, exec); err != nil {
			logs.Error("sfn: execution error", logs.String("arn", exec.ExecutionArn), logs.Err(err))
		}
	}()
}

// versionArn returns the version ARN an execution must be associated
// with, or the empty string for unqualified starts.
func (r *stateMachineReference) versionArn() string {
	if r.Version != nil {
		return r.Version.StateMachineVersionArn
	}
	return ""
}

// aliasArn returns the alias ARN an execution must be associated with, or
// the empty string when the start ARN was not an alias.
func (r *stateMachineReference) aliasArn() string {
	if r.Alias != nil {
		return r.Alias.StateMachineAliasArn
	}
	return ""
}

// startSyncExecutionCore is the single entry point for the synchronous
// start: it enforces the EXPRESS-only contract, validates like
// StartExecution, runs the execution to completion and returns the
// StartSyncExecutionOutput members.
func (s *StepFunctionService) startSyncExecutionCore(ctx context.Context, store *sfnstore.StepFunctionStore, in StartSyncExecutionInput) (map[string]interface{}, error) {
	if err := validateArnRequired(in.StateMachineArn, "stateMachineArn"); err != nil {
		return nil, err
	}
	if err := validateIncludedData(in.IncludedData); err != nil {
		return nil, err
	}
	if err := validateExecutionInputData(in.Input); err != nil {
		return nil, err
	}
	if err := validateTraceHeader(in.TraceHeader); err != nil {
		return nil, err
	}

	name := in.Name
	if name == "" {
		name = generateExecutionName()
	} else if err := validateExecutionName(name); err != nil {
		return nil, err
	}

	ref, err := resolveStateMachineReference(ctx, store, in.StateMachineArn)
	if err != nil {
		return nil, err
	}
	sm := ref.StateMachine
	if sm.Type != "EXPRESS" {
		return nil, NewStateMachineTypeNotSupported("StartSyncExecution is not available for STANDARD workflows")
	}

	executionArn := svcarn.NewARNBuilder(s.accountID, store.GetRegion()).StepFunctions().
		Execution(svcarn.ExtractStateMachineNameFromARN(sm.StateMachineArn), name)

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, in.Input, in.TraceHeader)
	exec.ExecutionArn = executionArn
	exec.StateMachineVersionArn = ref.versionArn()
	exec.StateMachineAliasArn = ref.aliasArn()

	if err := store.CreateExecution(ctx, exec); err != nil {
		if errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			return nil, NewExecutionAlreadyExists("An execution with the same name already exists: " + executionArn)
		}
		return nil, err
	}

	executor := NewExecutorWithStores(store, s.bus, s.accountID, store.GetRegion())
	_ = executor.ExecuteStateMachine(ctx, exec)

	updated, err := store.GetExecution(ctx, executionArn)
	if err != nil {
		updated = exec
	}

	metadataOnly := in.IncludedData == "METADATA_ONLY"
	result := map[string]interface{}{
		"executionArn":    updated.ExecutionArn,
		"stateMachineArn": updated.StateMachineArn,
		"name":            updated.Name,
		"startDate":       updated.StartDate.Unix(),
		"status":          updated.Status,
		"inputDetails":    map[string]interface{}{"included": !metadataOnly},
		"outputDetails":   map[string]interface{}{"included": !metadataOnly},
	}
	if !updated.StopDate.IsZero() {
		result["stopDate"] = updated.StopDate.Unix()
	}
	if !metadataOnly {
		if updated.Input != "" {
			result["input"] = updated.Input
		}
		if updated.Output != "" {
			result["output"] = updated.Output
		}
	}
	if updated.Error != "" {
		result["error"] = updated.Error
	}
	if updated.Cause != "" {
		result["cause"] = updated.Cause
	}
	if updated.TraceHeader != "" {
		result["traceHeader"] = updated.TraceHeader
	}
	return result, nil
}

// startExecutionForBusCore serves event-bus start requests (EventBridge,
// Scheduler, CloudWatch Alarms): it resolves the state machine by name
// when the event carries no ARN and runs the same validation and launch
// path as the HTTP StartExecution handler.
func (s *StepFunctionService) startExecutionForBusCore(ctx context.Context, store *sfnstore.StepFunctionStore, stateMachineArn, stateMachineName, input string) error {
	if stateMachineArn == "" && stateMachineName != "" {
		sm, err := store.GetStateMachineByName(ctx, stateMachineName)
		if err != nil {
			return err
		}
		stateMachineArn = sm.StateMachineArn
	}
	_, err := s.startExecutionCore(ctx, store, StartExecutionInput{
		StateMachineArn: stateMachineArn,
		Name:            "bus-" + generateExecutionName(),
		Input:           input,
	})
	return err
}

// ---------------------------------------------------------------------------
// Stop / Describe / List / History Core
// ---------------------------------------------------------------------------

// stopExecutionCore is the single entry point for StopExecution. The
// error and cause strings obey the SensitiveError (256) and SensitiveCause
// (32768) bounds; stopping an already-terminal execution returns its stop
// date without error.
func (s *StepFunctionService) stopExecutionCore(ctx context.Context, store *sfnstore.StepFunctionStore, in StopExecutionInput) (map[string]interface{}, error) {
	if err := validateArnRequired(in.ExecutionArn, "executionArn"); err != nil {
		return nil, err
	}
	if len(in.Error) > sfnstore.MaxErrorLength {
		return nil, NewValidationException(fmt.Sprintf("error must be at most %d characters, got %d", sfnstore.MaxErrorLength, len(in.Error)))
	}
	if len(in.Cause) > sfnstore.MaxCauseLength {
		return nil, NewValidationException(fmt.Sprintf("cause must be at most %d characters, got %d", sfnstore.MaxCauseLength, len(in.Cause)))
	}

	exec, err := store.GetExecution(ctx, in.ExecutionArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrExecutionNotFound) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + in.ExecutionArn)
		}
		return nil, err
	}

	if isTerminalStatus(exec.Status) {
		return map[string]interface{}{"stopDate": exec.StopDate.Unix()}, nil
	}

	store.CancelExecution(in.ExecutionArn)

	exec.Status = "ABORTED"
	exec.StopDate = time.Now().UTC()
	exec.Error = in.Error
	exec.Cause = in.Cause

	if err := store.UpdateExecution(ctx, exec); err != nil {
		return nil, err
	}

	return map[string]interface{}{"stopDate": exec.StopDate.Unix()}, nil
}

// describeExecutionCore is the single entry point for DescribeExecution.
// includedData=METADATA_ONLY omits the input and output payloads and
// reports them as not included.
func (s *StepFunctionService) describeExecutionCore(ctx context.Context, store *sfnstore.StepFunctionStore, in DescribeExecutionInput) (map[string]interface{}, error) {
	if err := validateArnRequired(in.ExecutionArn, "executionArn"); err != nil {
		return nil, err
	}
	if err := validateIncludedData(in.IncludedData); err != nil {
		return nil, err
	}

	exec, err := store.GetExecution(ctx, in.ExecutionArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrExecutionNotFound) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + in.ExecutionArn)
		}
		return nil, err
	}

	response := executionToResponse(exec)
	if in.IncludedData == "METADATA_ONLY" {
		delete(response, "input")
		delete(response, "output")
		response["inputDetails"] = map[string]interface{}{"included": false}
		response["outputDetails"] = map[string]interface{}{"included": false}
	} else {
		response["inputDetails"] = map[string]interface{}{"included": true}
		response["outputDetails"] = map[string]interface{}{"included": true}
	}
	return response, nil
}

// listExecutionsCore is the single entry point for ListExecutions. Results
// are sorted by time with the most recent execution first — running
// executions by startDate or redriveDate, the rest by stopDate — and the
// state machine ARN may be qualified with a version or alias to list the
// executions associated with it.
func (s *StepFunctionService) listExecutionsCore(ctx context.Context, store *sfnstore.StepFunctionStore, in ListExecutionsInput) (*ListExecutionsResult, error) {
	if !isValidExecutionStatus(in.StatusFilter) {
		return nil, NewValidationException("statusFilter must be one of RUNNING, SUCCEEDED, FAILED, TIMED_OUT, ABORTED, PENDING_REDRIVE, got " + in.StatusFilter)
	}
	if in.RedriveFilter != "" && in.RedriveFilter != "REDRIVEN" && in.RedriveFilter != "NOT_REDRIVEN" {
		return nil, NewValidationException("redriveFilter must be REDRIVEN or NOT_REDRIVEN, got " + in.RedriveFilter)
	}
	// The redrive filter applies to Distributed Map child executions listed
	// through mapRunArn; pairing it with a state machine ARN is documented
	// to fail with a validation exception.
	if in.RedriveFilter != "" && in.StateMachineArn != "" {
		return nil, NewValidationException("redriveFilter cannot be combined with stateMachineArn")
	}
	// PENDING_REDRIVE lists child workflow executions awaiting redrive;
	// those only exist in the scope of a Map Run, so the documented
	// contract requires mapRunArn and rejects a stateMachineArn pairing
	// with a validation exception — unconditionally, even when a
	// mapRunArn is also present.
	if in.StatusFilter == "PENDING_REDRIVE" && in.StateMachineArn != "" {
		return nil, NewValidationException("statusFilter PENDING_REDRIVE requires mapRunArn; providing stateMachineArn with PENDING_REDRIVE is not supported")
	}
	if in.MapRunArn != "" && in.StateMachineArn != "" {
		return nil, NewValidationException("mapRunArn and stateMachineArn are mutually exclusive")
	}
	if err := validateMaxResults(in.MaxResults, 0, sfnstore.MaxPageSize, "maxResults"); err != nil {
		return nil, err
	}
	maxResults := in.MaxResults
	if maxResults == 0 {
		maxResults = sfnstore.DefaultPageSize
	}

	filterArn := in.StateMachineArn
	association := ""
	if filterArn != "" {
		if err := validateArnRequired(filterArn, "stateMachineArn"); err != nil {
			return nil, err
		}
		// A qualified ARN filters by association: executions started with
		// the version ARN (or an alias routing to it) carry it.
		ref, err := resolveStateMachineReference(ctx, store, filterArn)
		if err != nil {
			return nil, err
		}
		filterArn = ref.StateMachine.StateMachineArn
		if ref.Alias != nil {
			association = ref.Alias.StateMachineAliasArn
		} else if ref.Version != nil {
			association = ref.Version.StateMachineVersionArn
		}
	} else if in.MapRunArn == "" {
		return nil, NewValidationException("stateMachineArn or mapRunArn is required")
	}

	all, err := store.ListAllExecutions(ctx, filterArn, in.StatusFilter, in.MapRunArn, in.RedriveFilter)
	if err != nil {
		return nil, err
	}

	executions := make([]*sfnstore.Execution, 0, len(all))
	for _, exec := range all {
		if association == "" ||
			exec.StateMachineVersionArn == association ||
			exec.StateMachineAliasArn == association {
			executions = append(executions, exec)
		}
	}

	sort.SliceStable(executions, func(i, j int) bool {
		ti, tj := listSortTime(executions[i]), listSortTime(executions[j])
		if !ti.Equal(tj) {
			return ti.After(tj)
		}
		return executions[i].ExecutionArn > executions[j].ExecutionArn
	})

	offset := 0
	if in.NextToken != "" {
		parsed, parseErr := strconv.Atoi(in.NextToken)
		if parseErr != nil || parsed < 0 || parsed > len(executions) {
			return nil, NewInvalidToken("Invalid nextToken: " + in.NextToken)
		}
		offset = parsed
	}

	end := offset + int(maxResults)
	if end > len(executions) {
		end = len(executions)
	}
	page := executions[offset:end]

	nextToken := ""
	if end < len(executions) {
		nextToken = strconv.Itoa(end)
	}

	return &ListExecutionsResult{Executions: page, NextToken: nextToken}, nil
}

// listSortTime returns the timestamp ListExecutions orders a running
// execution by (redriveDate when redriven, startDate otherwise) and a
// closed execution by (stopDate).
func listSortTime(exec *sfnstore.Execution) time.Time {
	if exec.Status == "RUNNING" {
		if !exec.RedriveDate.IsZero() {
			return exec.RedriveDate
		}
		return exec.StartDate
	}
	if !exec.StopDate.IsZero() {
		return exec.StopDate
	}
	return exec.StartDate
}

// getExecutionHistoryCore is the single entry point for GetExecutionHistory.
func (s *StepFunctionService) getExecutionHistoryCore(ctx context.Context, store *sfnstore.StepFunctionStore, in GetExecutionHistoryInput) (map[string]interface{}, error) {
	if err := validateArnRequired(in.ExecutionArn, "executionArn"); err != nil {
		return nil, err
	}
	if err := validateMaxResults(in.MaxResults, 0, sfnstore.MaxPageSize, "maxResults"); err != nil {
		return nil, err
	}
	limit := in.MaxResults
	if limit == 0 {
		limit = sfnstore.DefaultPageSize
	}

	if _, err := store.GetExecution(ctx, in.ExecutionArn); err != nil {
		return nil, NewExecutionDoesNotExist("Execution Does not exist: " + in.ExecutionArn)
	}

	// Reverse order must paginate in reverse as a whole: the store serves
	// newest-first pages with a direction-consistent marker, so reversing
	// an ascending page in place would scramble the global order.
	events, nextToken, err := store.GetExecutionHistory(ctx, in.ExecutionArn, limit, in.NextToken, in.ReverseOrder)
	if err != nil {
		if errors.Is(err, sfnstore.ErrInvalidToken) {
			return nil, NewInvalidToken("Invalid nextToken: " + in.NextToken)
		}
		return nil, err
	}

	history := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		history = append(history, historyEventToResponse(event, in.IncludeExecutionData))
	}

	response := map[string]interface{}{"events": history}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

// ---------------------------------------------------------------------------
// RedriveExecution Core
// ---------------------------------------------------------------------------

// redriveTokenEntry records one client token of a successful redrive for
// the documented idempotency window.
type redriveTokenEntry struct {
	token       string
	redriveDate time.Time
	expiresAt   time.Time
}

// redriveTokenCache holds the last client tokens per execution (at most
// ten, valid fifteen minutes) so retried RedriveExecution calls with the
// same token return the original redriveDate instead of redriving again.
type redriveTokenCache struct {
	mu      sync.Mutex
	entries map[string][]redriveTokenEntry
}

var redriveTokens = &redriveTokenCache{entries: map[string][]redriveTokenEntry{}}

// lookupToken returns the cached redriveDate for a token that is still
// within its validity window.
func (c *redriveTokenCache) lookupToken(executionArn, token string) (time.Time, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now().UTC()
	for _, e := range c.entries[executionArn] {
		if e.token == token && now.Before(e.expiresAt) {
			return e.redriveDate, true
		}
	}
	return time.Time{}, false
}

// record stores a token after a successful redrive, dropping expired
// entries and keeping at most ten per execution. Keys whose entries have
// all expired are removed as well, so the cache does not retain one entry
// per redriven execution for the process lifetime.
func (c *redriveTokenCache) record(executionArn, token string, redriveDate time.Time) {
	if token == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now().UTC()
	live := c.entries[executionArn][:0]
	for _, e := range c.entries[executionArn] {
		if now.Before(e.expiresAt) {
			live = append(live, e)
		}
	}
	live = append(live, redriveTokenEntry{
		token:       token,
		redriveDate: redriveDate,
		expiresAt:   now.Add(15 * time.Minute),
	})
	if len(live) > 10 {
		live = live[len(live)-10:]
	}
	c.entries[executionArn] = live

	for arn, list := range c.entries {
		anyLive := false
		for _, e := range list {
			if now.Before(e.expiresAt) {
				anyLive = true
				break
			}
		}
		if !anyLive {
			delete(c.entries, arn)
		}
	}
}

// redriveExecutionCore is the single entry point for RedriveExecution: it
// enforces the documented eligibility contract — STANDARD workflows only,
// unsuccessful terminal status, within fourteen days of completion and
// below the 24,999-event history ceiling — honours the clientToken
// idempotency window, and resumes the execution from its failed state
// with its ARN, input and history preserved.
func (s *StepFunctionService) redriveExecutionCore(ctx context.Context, store *sfnstore.StepFunctionStore, in RedriveExecutionInput) (*RedriveExecutionResult, error) {
	if err := validateArnRequired(in.ExecutionArn, "executionArn"); err != nil {
		return nil, err
	}

	if in.ClientToken != "" {
		if cached, ok := redriveTokens.lookupToken(in.ExecutionArn, in.ClientToken); ok {
			return &RedriveExecutionResult{RedriveDate: cached}, nil
		}
	}

	exec, err := store.GetExecution(ctx, in.ExecutionArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrExecutionNotFound) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + in.ExecutionArn)
		}
		return nil, err
	}

	sm, err := store.GetStateMachine(ctx, exec.StateMachineArn)
	if err != nil {
		return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + exec.StateMachineArn)
	}

	if sm.Type == "EXPRESS" {
		return nil, NewExecutionNotRedrivable(fmt.Sprintf("Execution %s is an EXPRESS workflow execution and cannot be redriven", in.ExecutionArn))
	}
	if !isRedrivableStatus(exec.Status) {
		return nil, NewExecutionNotRedrivable(fmt.Sprintf("Execution %s is in %s status and cannot be redriven", in.ExecutionArn, exec.Status))
	}
	if !exec.StopDate.IsZero() && time.Since(exec.StopDate) > sfnstore.RedriveWindowDays*24*time.Hour {
		return nil, NewExecutionNotRedrivable(fmt.Sprintf("Execution %s has exceeded the redrivable period of %d days", in.ExecutionArn, sfnstore.RedriveWindowDays))
	}
	eventCount, err := store.CountExecutionHistory(ctx, in.ExecutionArn)
	if err != nil {
		return nil, err
	}
	if eventCount >= sfnstore.MaxRedriveEventHistory {
		return nil, NewExecutionNotRedrivable(fmt.Sprintf("Execution %s has exceeded the execution history event limit", in.ExecutionArn))
	}

	definition, err := parseStateMachineDefinition(sm.Definition)
	if err != nil {
		return nil, NewInvalidDefinitionException("Invalid state machine definition: " + err.Error())
	}

	rp, err := determineResumePoint(ctx, store, in.ExecutionArn, definition)
	if err != nil {
		return nil, NewValidationException(fmt.Sprintf("failed to determine resume point: %v", err))
	}

	redriveDate := time.Now().UTC()
	exec.Status = "RUNNING"
	exec.Error = ""
	exec.Cause = ""
	exec.StopDate = time.Time{}
	exec.RedriveCount++
	exec.RedriveDate = redriveDate
	exec.Output = ""

	if err := store.UpdateExecution(ctx, exec); err != nil {
		return nil, NewConflictException(fmt.Sprintf("failed to update execution for redrive: %v", err))
	}

	executor := NewExecutorWithStores(store, s.bus, s.accountID, store.GetRegion())

	// Record the redrive in the history: the resumed state's next event
	// follows this one, so the resume point passes lastEventId+1.
	redriveEventId := rp.LastEventId + 1
	if err := executor.addExecutionHistoryEvent(ctx, exec, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn: in.ExecutionArn,
		EventId:      redriveEventId,
		Type:         "ExecutionRedriven",
		Timestamp:    redriveDate,
		ExecutionRedrivedEventDetails: &sfnstore.ExecutionRedrivedEventDetails{
			RedriveDate:     redriveDate,
			StateMachineArn: exec.StateMachineArn,
			ExecutionArn:    in.ExecutionArn,
		},
	}); err != nil {
		logs.Error("Failed to add ExecutionRedriven event", logs.Err(err))
	}

	resumeCtx, cancel := context.WithCancel(context.Background())
	store.RegisterExecution(in.ExecutionArn, cancel)
	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer store.UnregisterExecution(in.ExecutionArn)
		defer func() {
			if r := recover(); r != nil {
				logs.Error("sfn: panic in redrive execution", logs.String("arn", in.ExecutionArn), logs.Any("panic", r))
				exec.Status = "FAILED"
				exec.Error = "States.InternalError"
				exec.Cause = fmt.Sprintf("internal panic: %v", r)
				exec.StopDate = time.Now().UTC()
				_ = store.UpdateExecution(context.Background(), exec)
			}
		}()
		if err := executor.ExecuteStateMachineFromState(resumeCtx, exec, rp.StateName, rp.Input, redriveEventId); err != nil {
			logs.Error("sfn: redrive execution failed", logs.String("arn", in.ExecutionArn), logs.Err(err))
		}
	}()

	redriveTokens.record(in.ExecutionArn, in.ClientToken, redriveDate)

	return &RedriveExecutionResult{RedriveDate: redriveDate}, nil
}
