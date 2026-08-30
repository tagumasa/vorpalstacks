package athena

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// ErrInvalidRequestException is returned when the input fails constraints.
	ErrInvalidRequestException = awserrors.NewAWSError("InvalidRequestException", "The input failed to satisfy the constraints specified by an AWS service.", http.StatusBadRequest)
	// ErrMetadataException is returned when an error occurs while accessing metadata.
	ErrMetadataException = awserrors.NewAWSError("MetadataException", "An error occurred while accessing metadata.", http.StatusBadRequest)
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

// capacityReservationNotFound returns an InvalidRequestException for the
// specified capacity reservation. Per the Smithy model, CapacityReservation
// operations only define InternalServerException and InvalidRequestException —
// ResourceNotFoundException is not among the declared errors.
func capacityReservationNotFound(name string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidRequestException",
		fmt.Sprintf("CapacityReservation %s not found", name), http.StatusBadRequest)
}

// alreadyExistsInvalidRequest returns the InvalidRequestException every
// Athena create operation declares for a duplicate resource name — the
// model defines no ResourceAlreadyExistsException shape and the create
// operations list only InternalServerException and InvalidRequestException.
func alreadyExistsInvalidRequest(resourceType, name string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidRequestException",
		fmt.Sprintf("%s %s already exists", resourceType, name), http.StatusBadRequest)
}

// invalidRequestParameter returns the InvalidRequestException every Athena
// operation declares for a constraint-violating parameter — the model
// defines no InvalidParameterException shape, so parameter-constraint
// violations ride the operation-declared InvalidRequestException.
func invalidRequestParameter(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidRequestException", message, http.StatusBadRequest)
}
