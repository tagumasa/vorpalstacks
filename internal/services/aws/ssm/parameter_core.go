package ssm

import (
	"context"

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

// deleteParameterCore is the single entry point for SSM parameter deletion.
// It validates the name and delegates to the store. Callers map the returned
// error to their protocol-specific error type.
func (s *SSMService) deleteParameterCore(store ssmstore.SSMStoreInterface, name string) error {
	if name == "" {
		return ErrInvalidParameterName
	}
	return store.DeleteParameter(name)
}

// putParameterCore is the single entry point for SSM parameter creation
// and update. It delegates to putParameterWithEncryption which handles
// normalisation, KMS encryption for SecureString, and version management.
func (s *SSMService) putParameterCore(ctx context.Context, store ssmstore.SSMStoreInterface, param *ssmstore.Parameter, overwrite bool, modifiedBy string) (int64, error) {
	return s.putParameterWithEncryption(ctx, store, param, overwrite, modifiedBy)
}
