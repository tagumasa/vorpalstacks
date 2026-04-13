package errors

import (
	"fmt"
	"net/http"
)

// NewValidationException creates a new ValidationException error.
func NewValidationException(message string) *AWSError {
	return NewAWSError("ValidationException", message, http.StatusBadRequest)
}

// NewResourceNotFoundException creates a new ResourceNotFoundException error.
func NewResourceNotFoundException(resourceType, identifier string) *AWSError {
	return NewAWSError("ResourceNotFoundException",
		fmt.Sprintf("%s %s not found", resourceType, identifier), http.StatusNotFound)
}

// NewNotFoundException creates a new NotFoundException error.
func NewNotFoundException(resource string) *AWSError {
	return NewAWSError("NotFoundException",
		fmt.Sprintf("%s not found", resource), http.StatusNotFound)
}

// NewInvalidParameterException creates a new InvalidParameterException error.
func NewInvalidParameterException(message string) *AWSError {
	return NewAWSError("InvalidParameterException", message, http.StatusBadRequest)
}

// NewAccessDeniedException creates a new AccessDeniedException error.
func NewAccessDeniedException(message string) *AWSError {
	return NewAWSError("AccessDeniedException", message, http.StatusForbidden)
}

// NewThrottlingException creates a new ThrottlingException error.
func NewThrottlingException(message string) *AWSError {
	return NewAWSError("ThrottlingException", message, http.StatusTooManyRequests)
}

// NewServiceUnavailableException creates a new ServiceUnavailableException error.
func NewServiceUnavailableException(message string) *AWSError {
	return NewAWSError("ServiceUnavailableException", message, http.StatusServiceUnavailable)
}

// NewConflictException creates a new ConflictException error.
func NewConflictException(message string) *AWSError {
	return NewAWSError("ConflictException", message, http.StatusConflict)
}

// NewResourceAlreadyExistsException creates a new ResourceAlreadyExistsException error.
func NewResourceAlreadyExistsException(resource string) *AWSError {
	return NewAWSError("ResourceAlreadyExistsException",
		fmt.Sprintf("%s already exists", resource), http.StatusConflict)
}

// NewLimitExceededException creates a new LimitExceededException error.
func NewLimitExceededException(message string) *AWSError {
	return NewAWSError("LimitExceededException", message, http.StatusBadRequest)
}

// NewBadRequestException creates a new BadRequestException error.
func NewBadRequestException(message string) *AWSError {
	return NewAWSError("BadRequestException", message, http.StatusBadRequest)
}

// NewInternalErrorException creates a new InternalError exception error.
func NewInternalErrorException(message string) *AWSError {
	return NewAWSError("InternalError", message, http.StatusInternalServerError)
}

// NewInternalFailureException creates a new InternalFailure error.
func NewInternalFailureException(message string) *AWSError {
	return NewAWSError("InternalFailure", message, http.StatusInternalServerError)
}

// NewNoSuchEntityException creates a new NoSuchEntity error.
func NewNoSuchEntityException(resourceType, identifier string) *AWSError {
	return NewAWSError("NoSuchEntity",
		fmt.Sprintf("The %s with name %s cannot be found.", resourceType, identifier), http.StatusNotFound)
}

// NewEntityAlreadyExistsException creates a new EntityAlreadyExists error.
func NewEntityAlreadyExistsException(resource string) *AWSError {
	return NewAWSError("EntityAlreadyExists",
		fmt.Sprintf("%s already exists", resource), http.StatusConflict)
}

// NewDeleteConflictException creates a new DeleteConflict error.
func NewDeleteConflictException(message string) *AWSError {
	return NewAWSError("DeleteConflict", message, http.StatusConflict)
}

// NewInvalidInputException creates a new InvalidInput error.
func NewInvalidInputException(message string) *AWSError {
	return NewAWSError("InvalidInput", message, http.StatusBadRequest)
}

// NewInvalidParameterValueException creates a new InvalidParameterValue error.
func NewInvalidParameterValueException(message string) *AWSError {
	return NewAWSError("InvalidParameterValue", message, http.StatusBadRequest)
}

// NewResourceInUseException creates a new ResourceInUseException error.
func NewResourceInUseException(message string) *AWSError {
	return NewAWSError("ResourceInUseException", message, http.StatusBadRequest)
}
