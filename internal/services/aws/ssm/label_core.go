package ssm

import (
	"errors"

	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// LabelParameterVersionInput carries the fields for LabelParameterVersion.
type LabelParameterVersionInput struct {
	Name             string
	ParameterVersion int64
	Labels           []string
}

// LabelParameterVersionResult carries the per-label outcomes.
type LabelParameterVersionResult struct {
	InvalidLabels    []string
	ParameterVersion int64
}

// UnlabelParameterVersionInput carries the fields for UnlabelParameterVersion.
type UnlabelParameterVersionInput struct {
	Name             string
	ParameterVersion int64
	Labels           []string
}

// UnlabelParameterVersionResult carries the removed labels.
type UnlabelParameterVersionResult struct {
	RemovedLabels []string
}

// labelParameterVersionCore is the single entry point for attaching labels
// to a parameter version. A zero ParameterVersion resolves to the latest
// version first, matching the AWS defaulting contract.
func (s *SSMService) labelParameterVersionCore(store ssmstore.SSMStoreInterface, in LabelParameterVersionInput) (*LabelParameterVersionResult, error) {
	if in.Name == "" {
		return nil, ErrInvalidParameterName
	}

	if err := validateLabels(in.Labels); err != nil {
		return nil, err
	}

	parameterVersion := in.ParameterVersion
	if parameterVersion == 0 {
		// AWS spec: if ParameterVersion is omitted, default to the latest version.
		param, err := store.GetParameter(in.Name, false)
		if err != nil {
			return nil, ErrParameterNotFound
		}
		parameterVersion = param.Version
	}

	invalidLabels, err := store.LabelParameterVersion(in.Name, parameterVersion, in.Labels)
	if err != nil {
		if errors.Is(err, ssmstore.ErrParameterNotFound) {
			return nil, ErrParameterNotFound
		}
		if errors.Is(err, ssmstore.ErrParameterVersionNotFound) {
			return nil, ErrParameterVersionNotFound
		}
		return nil, err
	}

	return &LabelParameterVersionResult{
		InvalidLabels:    invalidLabels,
		ParameterVersion: parameterVersion,
	}, nil
}

// unlabelParameterVersionCore is the single entry point for removing labels
// from a parameter version. ParameterVersion is required by the Smithy
// contract — an omitted version fails the call instead of defaulting to the
// latest version (unlike LabelParameterVersion, whose version is optional).
func (s *SSMService) unlabelParameterVersionCore(store ssmstore.SSMStoreInterface, in UnlabelParameterVersionInput) (*UnlabelParameterVersionResult, error) {
	if in.Name == "" {
		return nil, ErrInvalidParameterName
	}

	if err := validateLabels(in.Labels); err != nil {
		return nil, err
	}

	// Smithy marks ParameterVersion as required for UnlabelParameterVersion:
	// "If it isn't present, the call will fail." The wire layer cannot
	// distinguish omission from an explicit zero, so both are rejected here.
	if in.ParameterVersion == 0 {
		return nil, ErrValidationException
	}

	removedLabels, err := store.UnlabelParameterVersion(in.Name, in.ParameterVersion, in.Labels)
	if err != nil {
		if errors.Is(err, ssmstore.ErrParameterNotFound) {
			return nil, ErrParameterNotFound
		}
		if errors.Is(err, ssmstore.ErrParameterVersionNotFound) {
			return nil, ErrParameterVersionNotFound
		}
		return nil, err
	}

	return &UnlabelParameterVersionResult{
		RemovedLabels: removedLabels,
	}, nil
}
