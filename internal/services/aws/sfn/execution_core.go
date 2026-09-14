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

	"github.com/google/uuid"

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
			if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
				return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
			}
			return nil, err
		}
		return &stateMachineReference{StateMachine: sm}, nil
	}

	if _, err := strconv.Atoi(qualifier); err == nil {
		// Numeric qualifier: a version ARN.
		version, err := store.GetStateMachineVersion(ctx, arn)
		if err != nil {
			if errors.Is(err, sfnstore.ErrStateMachineVersionNotFound) {
				return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
			}
			return nil, err
		}
		sm, err := store.GetStateMachine(ctx, version.StateMachineArn)
		if err != nil {
			if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
				return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + version.StateMachineArn)
			}
			return nil, err
		}
		return &stateMachineReference{StateMachine: sm, Version: version}, nil
	}

	// Non-numeric qualifier: an alias ARN.
	alias, err := store.GetStateMachineAlias(ctx, arn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAliasNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
		}
		return nil, err
	}
	sm, err := store.GetStateMachine(ctx, alias.StateMachineArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + alias.StateMachineArn)
		}
		return nil, err
	}

	versionArn, err := selectVersionByWeight(alias.RoutingConfiguration)
	if err != nil {
		return nil, err
	}
	version, err := store.GetStateMachineVersion(ctx, versionArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineVersionNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Version Does not exist: " + versionArn)
		}
		return nil, err
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

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, in.Input, in.TraceHeader)
	exec.StateMachineVersionArn = ref.versionArn()
	exec.StateMachineAliasArn = ref.aliasArn()

	if err := persistExecutionForStart(ctx, store, exec, s.accountID, store.GetRegion(), sm.StateMachineArn, sm.Type == "EXPRESS"); err != nil {
		if errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			// StartExecution is idempotent for STANDARD workflows: the
			// same name and input on a running execution returns that
			// execution. EXPRESS reuse never reaches this branch —
			// persistExecutionForStart disambiguates its ARN.
			if existing, idempotent, ierr := executionAlreadyExistsIdempotent(store, exec.ExecutionArn, in.Input); ierr == nil && idempotent {
				return &StartExecutionResult{ExecutionArn: existing.ExecutionArn, StartDate: existing.StartDate}, nil
			}
			return nil, NewExecutionAlreadyExists("An execution with the same name already exists: " + exec.ExecutionArn)
		}
		return nil, err
	}

	s.launchExecution(store, exec)

	return &StartExecutionResult{ExecutionArn: exec.ExecutionArn, StartDate: exec.StartDate}, nil
}

// launchExecution runs an execution asynchronously with panic isolation,
// registration for StopExecution cancellation and terminal-state
// persistence on panic.
// launchExecutionGoroutine registers the execution for StopExecution
// cancellation, runs run under the service's async wait group, and
// persists a FAILED terminal state if the goroutine panics — a panic must
// never leave a RUNNING zombie behind. The runLabel names the launch path
// in logs (execution, recovered execution, redrive execution).
// executionDrainWait bounds how long a redrive waits for the superseded
// run's goroutine to finish its terminal history and record writes. The
// drain is a handful of store writes in practice; the bound only guards
// against a wedged goroutine holding the redrive hostage — on timeout the
// transition's fresh-status guard still decides eligibility.
const executionDrainWait = 5 * time.Second

func (s *StepFunctionService) launchExecutionGoroutine(store *sfnstore.StepFunctionStore, exec *sfnstore.Execution, runLabel string, run func(ctx context.Context) error) {
	execCtx, cancel := context.WithCancel(context.Background())
	handle := store.RegisterExecution(exec.ExecutionArn, cancel)
	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer store.UnregisterExecution(exec.ExecutionArn, handle)
		defer func() {
			if r := recover(); r != nil {
				logs.Error("sfn: panic in "+runLabel, logs.String("arn", exec.ExecutionArn), logs.Any("panic", r))
				exec.Status = "FAILED"
				exec.Error = "States.Runtime"
				exec.Cause = fmt.Sprintf("internal panic: %v", r)
				exec.StopDate = time.Now().UTC()
				// A terminal status without its terminal event would leave
				// the history ending at the last state event; the event
				// continues the stored sequence.
				appendTerminalFailureEvent(context.Background(), store, exec, exec.Error, exec.Cause)
				_ = store.UpdateExecution(context.Background(), exec)
			}
		}()
		if err := run(execCtx); err != nil {
			logs.Error("sfn: "+runLabel+" failed", logs.String("arn", exec.ExecutionArn), logs.Err(err))
		}
	}()
}

func (s *StepFunctionService) launchExecution(store *sfnstore.StepFunctionStore, exec *sfnstore.Execution) {
	executor := NewExecutorWithStores(store, s.bus, s.accountID, store.GetRegion(), s.taskCredentialsAuthz)
	s.launchExecutionGoroutine(store, exec, "execution", func(ctx context.Context) error {
		return executor.ExecuteStateMachine(ctx, exec)
	})
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

// persistExecutionForStart creates the execution record under the two
// dialects' name contracts: STANDARD keeps the name-derived ARN and
// surfaces ErrExecutionAlreadyExists for the caller's idempotence check,
// while "For EXPRESS workflows, execution names can be reused" — the
// store keys executions by ARN, so a reused name persists under a freshly
// suffixed ARN while the record keeps the caller's name.
func persistExecutionForStart(ctx context.Context, store *sfnstore.StepFunctionStore, exec *sfnstore.Execution, accountID, region, smArn string, express bool) error {
	exec.ExecutionArn = svcarn.NewARNBuilder(accountID, region).StepFunctions().
		Execution(svcarn.ExtractStateMachineNameFromARN(smArn), exec.Name)
	for {
		err := store.CreateExecution(ctx, exec)
		if err == nil {
			return nil
		}
		if !express || !errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			return err
		}
		exec.ExecutionArn = svcarn.NewARNBuilder(accountID, region).StepFunctions().
			Execution(svcarn.ExtractStateMachineNameFromARN(smArn), exec.Name+"-"+uuid.New().String()[:8])
	}
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

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, in.Input, in.TraceHeader)
	exec.StateMachineVersionArn = ref.versionArn()
	exec.StateMachineAliasArn = ref.aliasArn()

	// StartSyncExecution is EXPRESS-only, so a name reuse persists under a
	// freshly suffixed ARN — the operation's response carries it.
	if err := persistExecutionForStart(ctx, store, exec, s.accountID, store.GetRegion(), sm.StateMachineArn, true); err != nil {
		if errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			return nil, NewExecutionAlreadyExists("An execution with the same name already exists: " + exec.ExecutionArn)
		}
		return nil, err
	}
	executionArn := exec.ExecutionArn

	executor := NewExecutorWithStores(store, s.bus, s.accountID, store.GetRegion(), s.taskCredentialsAuthz)
	// The inline run registers for StopExecution like the async launch:
	// without registration a stop only flips the stored record while the
	// executor keeps running and overwrites the ABORTED status with its
	// own terminal write.
	execCtx, cancel := context.WithCancel(ctx)
	handle := store.RegisterExecution(executionArn, cancel)
	_ = executor.ExecuteStateMachine(execCtx, exec)
	store.UnregisterExecution(executionArn, handle)
	cancel()

	updated, err := store.GetExecution(ctx, executionArn)
	if err != nil {
		updated = exec
	}

	metadataOnly := in.IncludedData == "METADATA_ONLY"
	result := map[string]interface{}{
		"executionArn":    updated.ExecutionArn,
		"stateMachineArn": updated.StateMachineArn,
		"name":            updated.Name,
		"startDate":       awsEpochSeconds(updated.StartDate),
		"status":          updated.Status,
		"inputDetails":    map[string]interface{}{"included": !metadataOnly},
	}
	// The details members follow the serialiser's presence contract:
	// outputDetails appears only when an output exists — a failed sync
	// execution has none — or when METADATA_ONLY explicitly reports the
	// withheld payload.
	if metadataOnly || updated.Output != "" {
		result["outputDetails"] = map[string]interface{}{"included": !metadataOnly}
	}
	if !updated.StopDate.IsZero() {
		result["stopDate"] = awsEpochSeconds(updated.StopDate)
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
		return map[string]interface{}{"stopDate": awsEpochSeconds(exec.StopDate)}, nil
	}

	// Persist the terminal record before cancelling: the executor's
	// terminal write must not overwrite the caller's error/cause pair in a
	// last-writer race.
	exec.Status = "ABORTED"
	exec.StopDate = time.Now().UTC()
	exec.Error = in.Error
	exec.Cause = in.Cause

	if err := store.UpdateExecution(ctx, exec); err != nil {
		return nil, err
	}

	store.CancelExecution(in.ExecutionArn)

	return map[string]interface{}{"stopDate": awsEpochSeconds(exec.StopDate)}, nil
}

// appendTerminalFailureEvent writes the ExecutionFailed terminal event for
// a path that fails an execution outside the executor (a goroutine panic
// or restart recovery): it continues the stored history's id sequence.
func appendTerminalFailureEvent(ctx context.Context, store *sfnstore.StepFunctionStore, exec *sfnstore.Execution, errorCode, cause string) {
	lastId, err := store.LastExecutionHistoryId(ctx, exec.ExecutionArn)
	if err != nil {
		logs.Error("sfn: failed to read the history tail for a terminal event", logs.String("arn", exec.ExecutionArn), logs.Err(err))
		return
	}
	event := &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    exec.ExecutionArn,
		EventId:         lastId + 1,
		PreviousEventId: lastId,
		Type:            "ExecutionFailed",
		Timestamp:       time.Now().UTC(),
		ExecutionFailedEventDetails: &sfnstore.ExecutionFailedEventDetails{
			Error: errorCode,
			Cause: cause,
		},
	}
	if err := store.AddExecutionHistoryEvent(ctx, event); err != nil {
		logs.Error("sfn: failed to append the ExecutionFailed event", logs.String("arn", exec.ExecutionArn), logs.Err(err))
	}
}

// describeExecutionCore is the single entry point for DescribeExecution.
// includedData=METADATA_ONLY omits the input and output payloads and
// reports them as not included; in the default mode the serialiser owns
// the presence contract (outputDetails appears only when an output
// exists).
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

	response := executionToResponse(ctx, store, exec)
	if in.IncludedData == "METADATA_ONLY" {
		delete(response, "input")
		delete(response, "output")
		response["inputDetails"] = map[string]interface{}{"included": false}
		response["outputDetails"] = map[string]interface{}{"included": false}
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
	maxResults := normaliseListLimit(in.MaxResults)

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
	limit := normaliseListLimit(in.MaxResults)

	if _, err := store.GetExecution(ctx, in.ExecutionArn); err != nil {
		if errors.Is(err, sfnstore.ErrExecutionNotFound) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + in.ExecutionArn)
		}
		return nil, err
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

// evaluateRedriveEligibility applies the documented RedriveExecution
// eligibility contract — STANDARD workflow, not a Distributed Map child
// (those "can only be redriven by their Map Run"), an unsuccessful terminal
// status, within the fourteen-day redrive window, below the history-event
// ceiling — and is the single verdict both RedriveExecution and the
// DescribeExecution redriveStatus derivation use, so the two can never
// disagree. The reason continues the sentence "Execution <arn> …" and is
// set exactly when eligible is false.
func evaluateRedriveEligibility(ctx context.Context, store *sfnstore.StepFunctionStore, sm *sfnstore.StateMachine, exec *sfnstore.Execution) (eligible bool, reason string, err error) {
	if sm.Type == "EXPRESS" {
		return false, "is an EXPRESS workflow execution and cannot be redriven", nil
	}
	if exec.MapRunArn != "" {
		return false, "is a Distributed Map child execution and can only be redriven by its Map Run", nil
	}
	if !isRedrivableStatus(exec.Status) {
		return false, "is in " + exec.Status + " status and cannot be redriven", nil
	}
	return redrivePeriodChecks(ctx, store, exec)
}

// redrivePeriodChecks applies the time window and history-event ceiling
// shared by the direct-redrive verdict and the child redriveStatus
// derivation. The reason continues the sentence "Execution <arn> …".
func redrivePeriodChecks(ctx context.Context, store *sfnstore.StepFunctionStore, exec *sfnstore.Execution) (eligible bool, reason string, err error) {
	if !exec.StopDate.IsZero() && time.Since(exec.StopDate) > sfnstore.RedriveWindowDays*24*time.Hour {
		return false, fmt.Sprintf("has exceeded the redrivable period of %d days", sfnstore.RedriveWindowDays), nil
	}
	eventCount, err := store.CountExecutionHistory(ctx, exec.ExecutionArn)
	if err != nil {
		return false, "", err
	}
	if eventCount >= sfnstore.MaxRedriveEventHistory {
		return false, "has exceeded the execution history event limit", nil
	}
	return true, "", nil
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
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + exec.StateMachineArn)
		}
		return nil, err
	}

	if eligible, why, elErr := evaluateRedriveEligibility(ctx, store, sm, exec); elErr != nil {
		return nil, elErr
	} else if !eligible {
		return nil, NewExecutionNotRedrivable(fmt.Sprintf("Execution %s %s", in.ExecutionArn, why))
	}

	definition, err := parseStateMachineDefinition(sm.Definition)
	if err != nil {
		return nil, NewInvalidDefinitionException("Invalid state machine definition: " + err.Error())
	}

	rp, err := determineResumePoint(ctx, store, in.ExecutionArn, definition)
	if err != nil {
		return nil, NewValidationException(fmt.Sprintf("failed to determine resume point: %v", err))
	}

	// The superseded run's goroutine flushes its terminal events and its
	// final record write after the terminal status became visible (the
	// stop persists the status first), and it numbers events from an
	// in-memory counter this path cannot see. Transitioning or numbering
	// before it drains lets its writes clobber the fresh record and the
	// ExecutionRedriven event, and a stale ABORTED write can re-arm a
	// second concurrent redrive. The registration's lifetime spans every
	// write the goroutine makes, so waiting for it to disappear orders
	// all of them before anything below; on timeout the transition guard
	// still re-checks the fresh status.
	store.WaitExecutionInactive(in.ExecutionArn, executionDrainWait)

	redriveDate := time.Now().UTC()
	// The status transition is serialised in the store: the eligibility
	// verdict above ran against the pre-transition record, so a concurrent
	// redrive of the same execution is rejected here against the fresh
	// status instead of double-transitioning it.
	exec, terr := store.TransitionExecutionForRedrive(ctx, in.ExecutionArn,
		func(fresh *sfnstore.Execution) bool {
			return fresh.MapRunArn == "" && isRedrivableStatus(fresh.Status)
		},
		func(fresh *sfnstore.Execution) {
			fresh.Status = "RUNNING"
			fresh.Error = ""
			fresh.Cause = ""
			fresh.StopDate = time.Time{}
			fresh.RedriveCount++
			fresh.RedriveDate = redriveDate
			fresh.Output = ""
		})
	if terr != nil {
		if errors.Is(terr, sfnstore.ErrExecutionNotRedrivable) {
			return nil, NewExecutionNotRedrivable(fmt.Sprintf("Execution %s %s", in.ExecutionArn, "is no longer in a redrivable state"))
		}
		if errors.Is(terr, sfnstore.ErrExecutionNotFound) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + in.ExecutionArn)
		}
		return nil, NewConflictException(fmt.Sprintf("failed to update execution for redrive: %v", terr))
	}

	executor := NewExecutorWithStores(store, s.bus, s.accountID, store.GetRegion(), s.taskCredentialsAuthz)

	// The event id follows the store's current maximum, not the resume
	// point's derivation: the drained goroutine's terminal events sit
	// above the failed state's last event, and the redriven history must
	// continue after all of them.
	lastId, lidErr := store.LastExecutionHistoryId(ctx, in.ExecutionArn)
	if lidErr != nil {
		logs.Error("Failed to read the last history id for redrive numbering", logs.Err(lidErr))
		lastId = rp.LastEventId
	}
	redriveEventId := lastId + 1
	if err := executor.addExecutionHistoryEvent(ctx, exec, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    in.ExecutionArn,
		EventId:         redriveEventId,
		PreviousEventId: lastId,
		Type:            "ExecutionRedriven",
		Timestamp:       redriveDate,
		ExecutionRedrivenEventDetails: &sfnstore.ExecutionRedrivenEventDetails{
			RedriveCount: exec.RedriveCount,
		},
	}); err != nil {
		logs.Error("Failed to add ExecutionRedriven event", logs.Err(err))
	}

	s.launchExecutionGoroutine(store, exec, "redrive execution", func(ctx context.Context) error {
		// A redrive does not re-arm the state-machine-level timeout.
		return executor.ExecuteStateMachineFromState(ctx, exec, rp.StateName, rp.Input, redriveEventId, false)
	})

	redriveTokens.record(in.ExecutionArn, in.ClientToken, redriveDate)

	return &RedriveExecutionResult{RedriveDate: redriveDate}, nil
}
