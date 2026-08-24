package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// CreateStateMachineInput carries every field that CreateStateMachine needs,
// in a format independent of the wire protocol (HTTP Query/JSON vs gRPC-Web).
// Both the HTTP API handler and the admin gRPC handler build this struct from
// their respective request formats and delegate to createStateMachineCore,
// ensuring that validation, role checking, and persistence follow a single
// code path.
//
// Note: the Smithy CreateStateMachineInput shape does NOT include a
// "description" field. Only PublishStateMachineVersion has a description
// (mapped to the VersionDescription shape).
type CreateStateMachineInput struct {
	Name                 string
	Definition           string
	RoleArn              string
	Type                 string
	LoggingConfiguration *sfnstore.LoggingConfiguration
	EncryptionConfig     *sfnstore.EncryptionConfiguration
	TracingConfig        *sfnstore.TracingConfiguration
	Tags                 map[string]string
	Publish              bool
	VersionDescription   string
}

// CreateStateMachineResult is the transport-agnostic result of
// createStateMachineCore.
type CreateStateMachineResult struct {
	StateMachineArn        string
	CreationDate           time.Time
	StateMachineVersionArn string
}

// UpdateStateMachineInput carries every field that UpdateStateMachine needs.
type UpdateStateMachineInput struct {
	StateMachineArn      string
	Definition           string
	DefinitionProvided   bool
	RoleArn              string
	RoleArnProvided      bool
	LoggingConfiguration *sfnstore.LoggingConfiguration
	EncryptionConfig     *sfnstore.EncryptionConfiguration
	TracingConfig        *sfnstore.TracingConfiguration
	Publish              bool
	VersionDescription   string
}

// UpdateStateMachineResult is the transport-agnostic result of
// updateStateMachineCore.
type UpdateStateMachineResult struct {
	StateMachineArn        string
	UpdateDate             time.Time
	RevisionId             string
	StateMachineVersionArn string
}

// ListStateMachinesInput carries the parameters for ListStateMachines.
type ListStateMachinesInput struct {
	MaxResults int32
	NextToken  string
}

// ListStateMachinesResult is the transport-agnostic result of
// listStateMachinesCore.
type ListStateMachinesResult struct {
	StateMachines []*sfnstore.StateMachine
	NextToken     string
}

// DeleteStateMachineInput carries the parameters for DeleteStateMachine.
type DeleteStateMachineInput struct {
	StateMachineArn string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createStateMachineCore is the single entry point for state machine creation
// logic shared by the HTTP API and the admin gRPC handler. It performs all
// Smithy-conformant validation, persists the state machine, and applies tags
// atomically — rolling back the creation if tagging fails.
func (s *StepFunctionService) createStateMachineCore(ctx context.Context, store *sfnstore.StepFunctionStore, in CreateStateMachineInput) (*CreateStateMachineResult, error) {
	// 1. Name (Smithy Name: @length(1,80), required via empty check).
	if err := validateResourceName(in.Name); err != nil {
		return nil, err
	}

	// 2. Definition (required, must be valid JSON).
	if err := validateDefinitionJSON(in.Definition); err != nil {
		return nil, err
	}

	// 3. Type (Smithy StateMachineType enum, default STANDARD).
	smType, err := validateStateMachineType(in.Type)
	if err != nil {
		return nil, err
	}

	// 4. ASL structural validation (Wait contract, JSONata field
	// separation and the full structural checks). Any ERROR-severity
	// diagnostic rejects the definition at creation time.
	if err := validateDefinitionStructure(in.Definition, smType); err != nil {
		return nil, err
	}

	// 5. RoleArn (Smithy @required on CreateStateMachineInput member).
	if err := validateRoleArnRequired(in.RoleArn); err != nil {
		return nil, err
	}

	// 6. LoggingConfiguration (Smithy LogLevel enum + AWS docs size-1 limit).
	if err := validateLoggingConfiguration(in.LoggingConfiguration); err != nil {
		return nil, err
	}

	// 7. EncryptionConfiguration (Smithy type enum + @range on kmsDataKeyReusePeriodSeconds).
	if in.EncryptionConfig != nil {
		if err := validateEncryptionConfiguration(in.EncryptionConfig); err != nil {
			return nil, err
		}
	}

	// 8. Tags (hard quota: fifty tags per resource).
	if len(in.Tags) > sfnstore.MaxTagsPerResource {
		return nil, NewTooManyTags(fmt.Sprintf("Too many tags: %d, maximum allowed %d", len(in.Tags), sfnstore.MaxTagsPerResource))
	}

	sm := &sfnstore.StateMachine{
		Name:                    in.Name,
		Definition:              in.Definition,
		RoleArn:                 in.RoleArn,
		Type:                    smType,
		Tags:                    in.Tags,
		VariableReferences:      extractVariableReferences(in.Definition),
		RevisionId:              generateRevisionId(),
		LoggingConfiguration:    in.LoggingConfiguration,
		EncryptionConfiguration: in.EncryptionConfig,
		TracingConfiguration:    in.TracingConfig,
	}

	if err := store.CreateStateMachine(ctx, sm); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAlreadyExists) {
			return s.idempotentCreateStateMachine(ctx, store, in, smType)
		}
		return nil, err
	}

	// Apply tags via TagStore so that ListTagsForResource works correctly.
	// If tagging fails, roll back the creation to avoid an untagged orphan.
	if len(in.Tags) > 0 {
		if err := store.Tag(sm.StateMachineArn, in.Tags); err != nil {
			_ = store.DeleteStateMachine(ctx, sm.StateMachineArn)
			return nil, err
		}
	}

	result := &CreateStateMachineResult{
		StateMachineArn: sm.StateMachineArn,
		CreationDate:    sm.CreationDate,
	}

	if in.Publish {
		version, err := store.PublishStateMachineVersion(ctx, sm.StateMachineArn, in.VersionDescription)
		if err != nil {
			return nil, err
		}
		result.StateMachineVersionArn = version.StateMachineVersionArn
	}

	return result, nil
}

// idempotentCreateStateMachine serves the documented CreateStateMachine
// idempotency when the name is already taken: a retry whose name,
// definition, type, logging, tracing and encryption configurations match
// the existing state machine — and whose publish and versionDescription
// parameters match the first version it published — returns the original
// resource. Differences in roleArn or tags are ignored. Any other
// same-name request stays a StateMachineAlreadyExists failure.
func (s *StepFunctionService) idempotentCreateStateMachine(ctx context.Context, store *sfnstore.StepFunctionStore, in CreateStateMachineInput, smType string) (*CreateStateMachineResult, error) {
	conflict := func() (*CreateStateMachineResult, error) {
		return nil, NewStateMachineAlreadyExists("A state machine with the same name already exists: " + in.Name)
	}

	existing, err := store.GetStateMachineByName(ctx, in.Name)
	if err != nil {
		return conflict()
	}

	var firstVersion *sfnstore.StateMachineVersion
	if v1, err := store.GetStateMachineVersionByNumber(ctx, existing.StateMachineArn, 1); err == nil {
		firstVersion = v1
	}

	publishMatches := in.Publish == (firstVersion != nil)
	descriptionMatches := true
	if in.Publish {
		if firstVersion == nil || firstVersion.Description != in.VersionDescription {
			descriptionMatches = false
		}
	} else if in.VersionDescription != "" {
		descriptionMatches = false
	}

	if existing.Definition != in.Definition ||
		existing.Type != smType ||
		!reflect.DeepEqual(existing.LoggingConfiguration, in.LoggingConfiguration) ||
		!reflect.DeepEqual(existing.EncryptionConfiguration, in.EncryptionConfig) ||
		!reflect.DeepEqual(existing.TracingConfiguration, in.TracingConfig) ||
		!publishMatches ||
		!descriptionMatches {
		return conflict()
	}

	result := &CreateStateMachineResult{
		StateMachineArn: existing.StateMachineArn,
		CreationDate:    existing.CreationDate,
	}
	if firstVersion != nil {
		result.StateMachineVersionArn = firstVersion.StateMachineVersionArn
	}
	return result, nil
}

// updateStateMachineCore is the single entry point for state machine updates
// shared by the HTTP API and the admin gRPC handler.
func (s *StepFunctionService) updateStateMachineCore(ctx context.Context, store *sfnstore.StepFunctionStore, in UpdateStateMachineInput) (*UpdateStateMachineResult, error) {
	if err := validateArnRequired(in.StateMachineArn, "stateMachineArn"); err != nil {
		return nil, err
	}

	// The update must carry at least one of definition or roleArn; the
	// API reference returns MissingRequiredParameter otherwise, even when
	// configuration-only fields are present.
	if !in.DefinitionProvided && !in.RoleArnProvided {
		return nil, NewMissingRequiredParameter("Request is missing a required parameter: update must include definition or roleArn")
	}
	// versionDescription is only valid together with publish=true.
	if in.VersionDescription != "" && !in.Publish {
		return nil, NewValidationException("versionDescription can only be specified when publish is true")
	}

	sm, err := store.GetStateMachine(ctx, in.StateMachineArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + in.StateMachineArn)
		}
		return nil, err
	}

	if in.DefinitionProvided {
		if err := validateDefinitionJSON(in.Definition); err != nil {
			return nil, err
		}
		if err := validateDefinitionStructure(in.Definition, sm.Type); err != nil {
			return nil, err
		}
		sm.Definition = in.Definition
		sm.VariableReferences = extractVariableReferences(in.Definition)
	}

	if in.RoleArnProvided {
		if err := validateRoleArnOptional(in.RoleArn); err != nil {
			return nil, err
		}
		sm.RoleArn = in.RoleArn
	}

	if in.LoggingConfiguration != nil {
		if err := validateLoggingConfiguration(in.LoggingConfiguration); err != nil {
			return nil, err
		}
		sm.LoggingConfiguration = in.LoggingConfiguration
	}

	if in.EncryptionConfig != nil {
		if err := validateEncryptionConfiguration(in.EncryptionConfig); err != nil {
			return nil, err
		}
		sm.EncryptionConfiguration = in.EncryptionConfig
	}

	if in.TracingConfig != nil {
		sm.TracingConfiguration = in.TracingConfig
	}

	sm.UpdateDate = time.Now().UTC()
	sm.RevisionId = generateRevisionId()

	if err := store.UpdateStateMachine(ctx, sm); err != nil {
		return nil, err
	}

	result := &UpdateStateMachineResult{
		StateMachineArn: sm.StateMachineArn,
		UpdateDate:      sm.UpdateDate,
		RevisionId:      sm.RevisionId,
	}

	if in.Publish {
		if err := enforceVersionQuota(ctx, store, sm.StateMachineArn); err != nil {
			return nil, err
		}
		version, err := store.PublishStateMachineVersion(ctx, sm.StateMachineArn, in.VersionDescription)
		if err != nil {
			return nil, err
		}
		result.StateMachineVersionArn = version.StateMachineVersionArn
	}

	return result, nil
}

// deleteStateMachineCore is the single entry point for state machine deletion.
func (s *StepFunctionService) deleteStateMachineCore(ctx context.Context, store *sfnstore.StepFunctionStore, in DeleteStateMachineInput) error {
	if err := validateArnRequired(in.StateMachineArn, "stateMachineArn"); err != nil {
		return err
	}
	if err := store.DeleteStateMachine(ctx, in.StateMachineArn); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return NewStateMachineDoesNotExist("State Machine Does not exist: " + in.StateMachineArn)
		}
		return err
	}
	return nil
}

// listStateMachinesCore is the single entry point for listing state machines.
func (s *StepFunctionService) listStateMachinesCore(ctx context.Context, store *sfnstore.StepFunctionStore, in ListStateMachinesInput) (*ListStateMachinesResult, error) {
	result, err := store.ListStateMachines(ctx, in.MaxResults, in.NextToken)
	if err != nil {
		return nil, err
	}
	return &ListStateMachinesResult{
		StateMachines: result.StateMachines,
		NextToken:     result.NextToken,
	}, nil
}

// ---------------------------------------------------------------------------
// Describe / validation Core — single validation + persistence path
// ---------------------------------------------------------------------------

// DescribeStateMachineInput carries the parameters for DescribeStateMachine.
type DescribeStateMachineInput struct {
	StateMachineArn string
	IncludedData    string
}

// DescribeStateMachineForExecutionInput carries the parameters for
// DescribeStateMachineForExecution.
type DescribeStateMachineForExecutionInput struct {
	ExecutionArn string
	IncludedData string
}

// ValidateStateMachineDefinitionInput carries the parameters for
// ValidateStateMachineDefinition.
type ValidateStateMachineDefinitionInput struct {
	Definition string
	SMType     string
	Severity   string
	MaxResults int32
}

// validateIncludedData validates the IncludedData enum (ALL_DATA |
// METADATA_ONLY); the empty value means ALL_DATA.
func validateIncludedData(includedData string) error {
	switch includedData {
	case "", "ALL_DATA", "METADATA_ONLY":
		return nil
	default:
		return NewValidationException("includedData must be ALL_DATA or METADATA_ONLY, got " + includedData)
	}
}

// describeStateMachineCore is the single entry point for
// DescribeStateMachine. A version-qualified ARN describes that version —
// the response carries the version ARN, the version creation date, the
// version description and the version's definition snapshot (API
// reference); an alias-qualified ARN describes the underlying state
// machine. includedData=METADATA_ONLY returns the definition as "{}".
// The response member set follows the DescribeStateMachineOutput shape:
// the operation does not return tags.
func (s *StepFunctionService) describeStateMachineCore(ctx context.Context, store *sfnstore.StepFunctionStore, in DescribeStateMachineInput) (map[string]interface{}, error) {
	if err := validateIncludedData(in.IncludedData); err != nil {
		return nil, err
	}

	ref, err := resolveStateMachineReference(ctx, store, in.StateMachineArn)
	if err != nil {
		return nil, err
	}

	sm := ref.StateMachine
	response := map[string]interface{}{
		"stateMachineArn": sm.StateMachineArn,
		"name":            sm.Name,
		"type":            sm.Type,
		"creationDate":    sm.CreationDate.Unix(),
	}
	if sm.Status != "" {
		response["status"] = sm.Status
	}
	if sm.RoleArn != "" {
		response["roleArn"] = sm.RoleArn
	}
	if sm.RevisionId != "" {
		response["revisionId"] = sm.RevisionId
	}
	if sm.LoggingConfiguration != nil {
		response["loggingConfiguration"] = sm.LoggingConfiguration
	}
	if sm.TracingConfiguration != nil {
		response["tracingConfiguration"] = sm.TracingConfiguration
	}
	if sm.EncryptionConfiguration != nil {
		response["encryptionConfiguration"] = sm.EncryptionConfiguration
	}

	definition := ref.definition()
	if ref.Version != nil {
		response["stateMachineArn"] = ref.Version.StateMachineVersionArn
		response["creationDate"] = ref.Version.CreationDate.Unix()
		if ref.Version.Description != "" {
			response["description"] = ref.Version.Description
		}
	}

	if in.IncludedData == "METADATA_ONLY" {
		response["definition"] = "{}"
	} else if definition != "" {
		response["definition"] = definition
	}
	if refs := extractVariableReferences(definition); len(refs) > 0 {
		response["variableReferences"] = refs
	}

	return response, nil
}

// describeStateMachineForExecutionCore is the single entry point for
// DescribeStateMachineForExecution. The response member set follows the
// DescribeStateMachineForExecutionOutput shape (updateDate, no type or
// status) and never includes tags. When the execution was started with a
// version or alias ARN, the described definition is the version snapshot
// the execution ran.
func (s *StepFunctionService) describeStateMachineForExecutionCore(ctx context.Context, store *sfnstore.StepFunctionStore, in DescribeStateMachineForExecutionInput) (map[string]interface{}, error) {
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

	sm, err := store.GetStateMachine(ctx, exec.StateMachineArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + exec.StateMachineArn)
		}
		return nil, err
	}

	definition := sm.Definition
	if exec.StateMachineVersionArn != "" {
		if version, verr := store.GetStateMachineVersion(ctx, exec.StateMachineVersionArn); verr == nil {
			definition = version.Definition
		}
	}

	response := map[string]interface{}{
		"stateMachineArn": sm.StateMachineArn,
		"name":            sm.Name,
	}
	if sm.RoleArn != "" {
		response["roleArn"] = sm.RoleArn
	}
	if !sm.UpdateDate.IsZero() {
		response["updateDate"] = sm.UpdateDate.Unix()
	}
	if sm.RevisionId != "" {
		response["revisionId"] = sm.RevisionId
	}
	if sm.LoggingConfiguration != nil {
		response["loggingConfiguration"] = sm.LoggingConfiguration
	}
	if sm.TracingConfiguration != nil {
		response["tracingConfiguration"] = sm.TracingConfiguration
	}
	if sm.EncryptionConfiguration != nil {
		response["encryptionConfiguration"] = sm.EncryptionConfiguration
	}
	if exec.MapRunArn != "" {
		response["mapRunArn"] = exec.MapRunArn
	}

	if in.IncludedData == "METADATA_ONLY" {
		response["definition"] = "{}"
	} else if definition != "" {
		response["definition"] = definition
	}
	if refs := extractVariableReferences(definition); len(refs) > 0 {
		response["variableReferences"] = refs
	}

	return response, nil
}

// validateStateMachineDefinitionCore is the single entry point for
// ValidateStateMachineDefinition: it validates the request parameters and
// computes the diagnostics over the definition without persisting
// anything, so it needs no store handle. The diagnostics come from the
// shared ASL structural validator and follow the documented
// ValidateStateMachineDefinitionDiagnostic code set; AWS guarantees only
// the stability of the result value, not the exact code, order or count.
func (s *StepFunctionService) validateStateMachineDefinitionCore(in ValidateStateMachineDefinitionInput) (map[string]interface{}, error) {
	// The length bound is a request-shape constraint; emptiness and JSON
	// syntax problems surface as diagnostics instead.
	if len(in.Definition) > sfnstore.MaxDefinitionLength {
		return nil, NewInvalidDefinitionException(fmt.Sprintf("State Machine definition must be at most %d bytes, got %d", sfnstore.MaxDefinitionLength, len(in.Definition)))
	}

	severity, err := validateSeverity(in.Severity)
	if err != nil {
		return nil, err
	}

	smType, err := validateStateMachineType(in.SMType)
	if err != nil {
		return nil, err
	}

	maxResults := in.MaxResults
	if err := validateMaxResults(maxResults, 0, sfnstore.MaxValidateDefinitionResults, "maxResults"); err != nil {
		return nil, err
	}
	if maxResults == 0 {
		maxResults = sfnstore.MaxValidateDefinitionResults
	}

	structural := validateASLStructure(in.Definition, smType)

	diagnostics := []map[string]string{}
	for _, d := range structural {
		entry := map[string]string{
			"severity": d.Severity,
			"code":     d.Code,
			"message":  d.Message,
		}
		if d.Location != "" {
			entry["location"] = d.Location
		}
		diagnostics = append(diagnostics, entry)
	}

	if severity == "ERROR" {
		filtered := []map[string]string{}
		for _, d := range diagnostics {
			if d["severity"] == "ERROR" {
				filtered = append(filtered, d)
			}
		}
		diagnostics = filtered
	}

	truncated := false
	if maxResults > 0 && len(diagnostics) > int(maxResults) {
		diagnostics = diagnostics[:maxResults]
		truncated = true
	}

	// Warnings do not prevent deploying a workflow definition, so the
	// result only fails when at least one ERROR-severity diagnostic is
	// present; a warning-only definition stays OK.
	result := "OK"
	for _, d := range diagnostics {
		if d["severity"] == "ERROR" {
			result = "FAIL"
			break
		}
	}

	return map[string]interface{}{
		"result":      result,
		"diagnostics": diagnostics,
		"truncated":   truncated,
	}, nil
}

// validateDefinitionStructure runs the shared ASL structural validator
// for the creation and update paths: any ERROR-severity diagnostic
// rejects the definition with the creation-time InvalidDefinition shape.
func validateDefinitionStructure(definition, smType string) error {
	for _, d := range validateASLStructure(definition, smType) {
		if d.Severity != "ERROR" {
			continue
		}
		location := ""
		if d.Location != "" {
			location = " at " + d.Location
		}
		return NewInvalidDefinitionException(d.Message + location)
	}
	return nil
}

// ---------------------------------------------------------------------------
// JSON parsing helpers (used by both HTTP handler and core)
// ---------------------------------------------------------------------------

// parseLoggingConfigurationFromJSON deserialises a LoggingConfiguration
// from a raw interface{} value (typically from request parameters).
func parseLoggingConfigurationFromJSON(raw interface{}) (*sfnstore.LoggingConfiguration, error) {
	if raw == nil {
		return nil, nil
	}
	// The configuration must be a JSON object; reject strings, arrays and
	// scalars with the specific configuration error instead of an opaque
	// marshal/unmarshal failure.
	if _, ok := raw.(map[string]interface{}); !ok {
		return nil, NewInvalidLoggingConfiguration("loggingConfiguration must be a JSON object")
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, NewInvalidLoggingConfiguration("loggingConfiguration is not serialisable: " + err.Error())
	}
	var lc sfnstore.LoggingConfiguration
	if err := json.Unmarshal(bytes, &lc); err != nil {
		return nil, NewInvalidLoggingConfiguration("loggingConfiguration is not valid JSON: " + err.Error())
	}
	return &lc, nil
}

// parseEncryptionConfigurationFromJSON deserialises an EncryptionConfiguration
// from a raw interface{} value.
func parseEncryptionConfigurationFromJSON(raw interface{}) (*sfnstore.EncryptionConfiguration, error) {
	if raw == nil {
		return nil, nil
	}
	if _, ok := raw.(map[string]interface{}); !ok {
		return nil, NewInvalidEncryptionConfiguration("encryptionConfiguration must be a JSON object")
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, NewInvalidEncryptionConfiguration("encryptionConfiguration is not serialisable: " + err.Error())
	}
	var ec sfnstore.EncryptionConfiguration
	if err := json.Unmarshal(bytes, &ec); err != nil {
		return nil, NewInvalidEncryptionConfiguration("encryptionConfiguration is not valid JSON: " + err.Error())
	}
	return &ec, nil
}

// parseTracingConfigurationFromJSON deserialises a TracingConfiguration
// from a raw interface{} value.
func parseTracingConfigurationFromJSON(raw interface{}) (*sfnstore.TracingConfiguration, error) {
	if raw == nil {
		return nil, nil
	}
	if _, ok := raw.(map[string]interface{}); !ok {
		return nil, NewInvalidTracingConfiguration("tracingConfiguration must be a JSON object")
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, NewInvalidTracingConfiguration("tracingConfiguration is not serialisable: " + err.Error())
	}
	var tc sfnstore.TracingConfiguration
	if err := json.Unmarshal(bytes, &tc); err != nil {
		return nil, NewInvalidTracingConfiguration("tracingConfiguration is not valid JSON: " + err.Error())
	}
	return &tc, nil
}
