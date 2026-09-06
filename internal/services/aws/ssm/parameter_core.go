package ssm

import (
	"context"
	"errors"
	"net/http"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/logs"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// DescribeParametersInput is the transport-agnostic input for listing SSM
// parameters. Both the admin handler and the AWS API HTTP handler convert
// their request-specific types into this struct before calling the Core.
type DescribeParametersInput struct {
	Filters    []ssmstore.ParameterFilter
	MaxResults int32
	NextToken  string
}

// DescribeParametersResult holds the raw store types returned by the Core.
// Each caller converts to its own response format (protobuf or map).
type DescribeParametersResult struct {
	Parameters []*ssmstore.ParameterMetadata
	NextToken  string
}

// describeParametersCore is the single entry point for SSM parameter
// listing. It applies default pagination and delegates to the store.
func (s *SSMService) describeParametersCore(store ssmstore.SSMStoreInterface, in DescribeParametersInput) (*DescribeParametersResult, error) {
	maxResults := in.MaxResults
	if maxResults <= 0 {
		maxResults = 50
	}

	params, nextToken, err := store.DescribeParameters(in.Filters, maxResults, in.NextToken)
	if err != nil {
		return nil, err
	}

	return &DescribeParametersResult{
		Parameters: params,
		NextToken:  nextToken,
	}, nil
}

// DescribeParametersWireInput carries the raw AWS API wire parameters for
// DescribeParameters so that MaxResults range enforcement and filter parsing
// happen in the Core rather than in a handler.
type DescribeParametersWireInput struct {
	Parameters map[string]interface{}
	MaxResults int32
	NextToken  string
}

// describeParametersWireCore is the HTTP-plane entry point for
// DescribeParameters: it validates MaxResults, parses the raw
// ParameterFilters/Filters wire values, and delegates to the shared
// describeParametersCore used by both protocol planes.
func (s *SSMService) describeParametersWireCore(store ssmstore.SSMStoreInterface, in DescribeParametersWireInput) (*DescribeParametersResult, error) {
	maxResults, err := validateMaxResultsForPage(in.MaxResults)
	if err != nil {
		return nil, err
	}

	filters, err := parseParameterFilters(in.Parameters)
	if err != nil {
		return nil, err
	}

	return s.describeParametersCore(store, DescribeParametersInput{
		Filters:    filters,
		MaxResults: maxResults,
		NextToken:  in.NextToken,
	})
}

// ParameterPutFields holds the raw string inputs of a PutParameter call.
// Both the HTTP query-protocol path and the admin gRPC path populate this
// from their respective wire formats and then call normalisePutParameter so
// validation, defaulting and Tier auto-promotion happen once in a single place.
type ParameterPutFields struct {
	Name           string
	Value          string
	Type           string
	Description    string
	KeyID          string
	AllowedPattern string
	DataType       string
	Tier           string
	Policies       string
	Tags           map[string]string
}

// PutParameterInput carries the raw PutParameter fields plus the overwrite
// flag and the modifying principal from either protocol plane.
type PutParameterInput struct {
	Fields     ParameterPutFields
	Overwrite  bool
	ModifiedBy string
}

// PutParameterResult reports the stored version and the effective tier
// (after auto-promotion) of a PutParameter call.
type PutParameterResult struct {
	Version int64
	Tier    ssmstore.ParameterTier
}

// putParameterCore is the single entry point for parameter creation and
// update on both protocol planes. It normalises and validates the raw
// fields, applies SecureString encryption via the shared
// putParameterWithEncryption helper, and maps store sentinel errors to
// their AWS API shapes.
func (s *SSMService) putParameterCore(ctx context.Context, store ssmstore.SSMStoreInterface, in PutParameterInput) (*PutParameterResult, error) {
	param, err := normalisePutParameter(in.Fields)
	if err != nil {
		return nil, err
	}

	// Tags travel inside the parameter: store.PutParameter persists them
	// alongside the new version and returns an error on failure, so a tag
	// problem fails the call instead of being silently dropped.
	version, err := s.putParameterWithEncryption(ctx, store, param, in.Overwrite, in.ModifiedBy)
	if err != nil {
		if errors.Is(err, ssmstore.ErrParameterAlreadyExists) {
			return nil, ErrParameterAlreadyExists
		}
		if errors.Is(err, ssmstore.ErrReservedParameterName) {
			return nil, ErrParameterPatternMismatch
		}
		if errors.Is(err, ssmstore.ErrInvalidAllowedPattern) {
			return nil, ErrInvalidAllowedPattern
		}
		if errors.Is(err, ssmstore.ErrParameterPatternMismatch) {
			return nil, ErrParameterPatternMismatch
		}
		if errors.Is(err, ssmstore.ErrHierarchyLevelLimitExceeded) {
			return nil, ErrHierarchyLevelLimitExceeded
		}
		return nil, err
	}

	return &PutParameterResult{Version: version, Tier: param.Tier}, nil
}

// normalisePutParameter validates every PutParameter input field and returns
// a fully-populated Parameter. Defaults (DataType="text", Tier="Standard") and
// Tier auto-promotion (Standard -> Advanced on >4KB or Policies) are applied
// here so callers cannot diverge. The returned Parameter is ready for the
// store; callers only need to supply the LastModifiedBy before storing.
func normalisePutParameter(in ParameterPutFields) (*ssmstore.Parameter, error) {
	if in.Name == "" {
		return nil, ErrInvalidParameterName
	}
	// Smithy marks Value as required; AWS rejects an empty value with
	// ValidationException.
	if in.Value == "" {
		return nil, ErrInvalidParameterValue
	}
	if utf8.RuneCountInString(in.Description) > ssmstore.MaxParameterDescriptionLength {
		return nil, ErrInvalidParameterValue
	}

	paramType := ssmstore.ParameterType(in.Type)
	if paramType == "" {
		paramType = ssmstore.ParameterTypeString
	}
	switch paramType {
	case ssmstore.ParameterTypeString, ssmstore.ParameterTypeStringList, ssmstore.ParameterTypeSecureString:
	default:
		return nil, ErrInvalidParameterType
	}

	if err := validateKeyID(in.KeyID); err != nil {
		return nil, err
	}
	if err := validateAllowedPattern(in.AllowedPattern); err != nil {
		return nil, err
	}

	dataType := in.DataType
	if dataType == "" {
		dataType = "text"
	}
	if err := validateDataType(dataType); err != nil {
		return nil, err
	}

	tier := ssmstore.ParameterTier(in.Tier)
	if in.Tier != "" {
		if err := validateTier(tier); err != nil {
			return nil, err
		}
	}

	if err := validatePolicies(in.Policies); err != nil {
		return nil, err
	}

	param := ssmstore.NewParameter(in.Name, in.Value, paramType)
	param.Description = in.Description
	param.KeyID = in.KeyID
	param.AllowedPattern = in.AllowedPattern
	param.DataType = dataType
	param.Tags = in.Tags
	param.Tier = tier
	param.Policies = in.Policies

	// AWS auto-promotes Standard-tier parameters to Advanced when the value
	// exceeds the 4KB Standard-tier limit or when any Policies are attached.
	// The value limit is documented in bytes ("Standard parameters have a
	// value limit of 4 KB", PutParameter API reference), so byte length is
	// the correct metric here — unlike @length traits, which count Unicode
	// characters. An empty Tier means the caller omitted it — AWS treats
	// that as Standard-equivalent, so the same promotion rule must apply.
	if param.Tier == "" || param.Tier == ssmstore.ParameterTierStandard {
		if len(param.Value) > 4096 || param.Policies != "" {
			param.Tier = ssmstore.ParameterTierAdvanced
		}
	}

	return param, nil
}

// deleteParameterCore is the single entry point for SSM parameter deletion.
// It validates the name and delegates to the store. Callers map the returned
// error to their protocol-specific error type.
func (s *SSMService) deleteParameterCore(store ssmstore.SSMStoreInterface, name string) error {
	if name == "" {
		return ErrInvalidParameterName
	}
	if err := store.DeleteParameter(name); err != nil {
		if errors.Is(err, ssmstore.ErrParameterNotFound) {
			return ErrParameterNotFound
		}
		logs.Error("Failed to delete parameter from store",
			logs.String("name", name), logs.Err(err))
		// InternalServerError is the error DeleteParameter declares in the
		// Smithy model (awsJson protocol: the wire code is the shape name).
		return awserrors.NewAWSError("InternalServerError", "failed to delete parameter", http.StatusInternalServerError)
	}
	return nil
}

// DeleteParametersInput carries the names for a bulk delete.
type DeleteParametersInput struct {
	Names []string
}

// DeleteParametersResult reports the per-name deletion outcomes.
type DeleteParametersResult struct {
	DeletedParameters []string
	InvalidParameters []string
}

// deleteParametersCore is the single entry point for bulk parameter
// deletion. Names that do not exist are reported as invalid, matching the
// AWS partial-success contract.
func (s *SSMService) deleteParametersCore(store ssmstore.SSMStoreInterface, in DeleteParametersInput) (*DeleteParametersResult, error) {
	if err := validateParameterNameList(in.Names); err != nil {
		return nil, err
	}

	deleted, invalid := store.DeleteParameters(in.Names)
	return &DeleteParametersResult{
		DeletedParameters: deleted,
		InvalidParameters: invalid,
	}, nil
}
