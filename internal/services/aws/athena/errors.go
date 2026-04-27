package athena

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// ErrInvalidRequestException is returned when the input fails constraints.
	ErrInvalidRequestException = awserrors.NewAWSError("InvalidRequestException", "The input failed to satisfy the constraints specified by an AWS service.", http.StatusBadRequest)
	// ErrInternalServerException is returned when an internal server error occurs.
	ErrInternalServerException = awserrors.NewInternalErrorException("An internal server error occurred.")
	// ErrInvalidParameterException is returned when a parameter value is invalid.
	ErrInvalidParameterException = awserrors.NewInvalidParameterException("Invalid parameter value.")
	// ErrResourceAlreadyExistsException is returned when a resource already exists.
	ErrResourceAlreadyExistsException = awserrors.NewResourceAlreadyExistsException("Resource")
	// ErrTooManyRequestsException is returned when too many requests are received.
	ErrTooManyRequestsException = awserrors.NewThrottlingException("Too many requests have been received.")
	// ErrMetadataException is returned when an error occurs while accessing metadata.
	ErrMetadataException = awserrors.NewAWSError("MetadataException", "An error occurred while accessing metadata.", http.StatusBadRequest)
	// ErrInvalidConfigurationException is returned when the configuration is invalid.
	ErrInvalidConfigurationException = awserrors.NewBadRequestException("The configuration is invalid.")
)

// workGroupNotFound returns a ResourceNotFoundException for the specified work group.
func workGroupNotFound(name string) *awserrors.AWSError {
	return awserrors.NewResourceNotFoundException("WorkGroup", name)
}

// namedQueryNotFound returns a ResourceNotFoundException for the specified named query.
func namedQueryNotFound(id string) *awserrors.AWSError {
	return awserrors.NewResourceNotFoundException("NamedQuery", id)
}

// dataCatalogNotFound returns a ResourceNotFoundException for the specified data catalog.
func dataCatalogNotFound(name string) *awserrors.AWSError {
	return awserrors.NewResourceNotFoundException("DataCatalog", name)
}

// queryExecutionNotFound returns a ResourceNotFoundException for the specified query execution.
func queryExecutionNotFound(id string) *awserrors.AWSError {
	return awserrors.NewResourceNotFoundException("QueryExecution", id)
}

// preparedStatementNotFound returns a ResourceNotFoundException for the specified prepared statement.
func preparedStatementNotFound(name string) *awserrors.AWSError {
	return awserrors.NewResourceNotFoundException("PreparedStatement", name)
}
