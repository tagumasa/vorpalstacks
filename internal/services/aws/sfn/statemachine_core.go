package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	Type                 string
	TypeProvided         bool
	LoggingConfiguration *sfnstore.LoggingConfiguration
	EncryptionConfig     *sfnstore.EncryptionConfiguration
	TracingConfig        *sfnstore.TracingConfiguration
	Publish              bool
	VersionDescription   string
	RevisionId           string
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
	if err := validateWaitStates(in.Definition); err != nil {
		return nil, err
	}

	// 3. Type (Smithy StateMachineType enum, default STANDARD).
	smType, err := validateStateMachineType(in.Type)
	if err != nil {
		return nil, err
	}

	// 4. RoleArn (Smithy @required on CreateStateMachineInput member).
	if err := validateRoleArnRequired(in.RoleArn); err != nil {
		return nil, err
	}

	// 5. LoggingConfiguration (Smithy LogLevel enum + AWS docs size-1 limit).
	if err := validateLoggingConfiguration(in.LoggingConfiguration); err != nil {
		return nil, err
	}

	// 6. EncryptionConfiguration (Smithy type enum + @range on kmsDataKeyReusePeriodSeconds).
	if in.EncryptionConfig != nil {
		if err := validateEncryptionConfiguration(in.EncryptionConfig); err != nil {
			return nil, err
		}
	}

	// 7. JSONata field validation.
	if err := validateDefinitionJSONataFields(in.Definition); err != nil {
		return nil, err
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
			return nil, NewStateMachineAlreadyExists("A state machine with the same name already exists: " + in.Name)
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

// updateStateMachineCore is the single entry point for state machine updates
// shared by the HTTP API and the admin gRPC handler.
func (s *StepFunctionService) updateStateMachineCore(ctx context.Context, store *sfnstore.StepFunctionStore, in UpdateStateMachineInput) (*UpdateStateMachineResult, error) {
	if err := validateArnRequired(in.StateMachineArn, "stateMachineArn"); err != nil {
		return nil, err
	}

	sm, err := store.GetStateMachine(ctx, in.StateMachineArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + in.StateMachineArn)
		}
		return nil, err
	}

	// Optimistic concurrency check.
	if in.RevisionId != "" && sm.RevisionId != in.RevisionId {
		return nil, NewValidationException("revisionId mismatch: expected " + sm.RevisionId + ", got " + in.RevisionId)
	}

	if in.DefinitionProvided {
		if err := validateDefinitionJSON(in.Definition); err != nil {
			return nil, err
		}
		if err := validateWaitStates(in.Definition); err != nil {
			return nil, err
		}
		if err := validateDefinitionJSONataFields(in.Definition); err != nil {
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

	if in.TypeProvided {
		smType, err := validateStateMachineType(in.Type)
		if err != nil {
			return nil, err
		}
		sm.Type = smType
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
