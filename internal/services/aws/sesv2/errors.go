package sesv2

import (
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// SESv2Error represents an error returned by the SESv2 service.
type SESv2Error struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error for error chaining.
func (e *SESv2Error) Unwrap() error {
	return e.AWSError
}

// Error variables for common SESv2 error conditions.
var (
	ErrAlreadyExists             = &SESv2Error{awserrors.NewAWSError("AlreadyExistsException", "Resource already exists", http.StatusConflict)}
	ErrBadRequest                = &SESv2Error{awserrors.NewBadRequestException("Invalid request")}
	ErrNotFound                  = &SESv2Error{awserrors.NewNotFoundException("Resource")}
	ErrLimitExceeded             = &SESv2Error{awserrors.NewLimitExceededException("Limit exceeded")}
	ErrTooManyRequests           = &SESv2Error{awserrors.NewThrottlingException("Too many requests")}
	ErrMessageRejected           = &SESv2Error{awserrors.NewAWSError("MessageRejected", "Message rejected", http.StatusBadRequest)}
	ErrAccountSuspended          = &SESv2Error{awserrors.NewAWSError("AccountSuspendedException", "Account suspended", http.StatusBadRequest)}
	ErrSendingPaused             = &SESv2Error{awserrors.NewAWSError("SendingPausedException", "Sending paused", http.StatusBadRequest)}
	ErrMailFromDomainNotVerified = &SESv2Error{awserrors.NewAWSError("MailFromDomainNotVerifiedException", "MAIL FROM domain not verified", http.StatusBadRequest)}
	ErrConcurrentModification    = &SESv2Error{awserrors.NewAWSError("ConcurrentModificationException", "Concurrent modification", http.StatusInternalServerError)}
	ErrMissingParameter          = &SESv2Error{awserrors.NewBadRequestException("Missing required parameter")}
	ErrInvalidParameter          = &SESv2Error{awserrors.NewBadRequestException("Invalid parameter")}
	ErrIdentityNotFound          = &SESv2Error{awserrors.NewNotFoundException("Email identity")}
	ErrConfigurationSetNotFound  = &SESv2Error{awserrors.NewNotFoundException("Configuration set")}
	ErrTemplateNotFound          = &SESv2Error{awserrors.NewNotFoundException("Email template")}
	ErrContactListNotFound       = &SESv2Error{awserrors.NewNotFoundException("Contact list")}
	ErrContactNotFound           = &SESv2Error{awserrors.NewNotFoundException("Contact")}
)

// NewAlreadyExistsException creates a new AlreadyExistsException error.
func NewAlreadyExistsException(resource string) *SESv2Error {
	return &SESv2Error{awserrors.NewResourceAlreadyExistsException(resource)}
}

// NewNotFoundException creates a new NotFoundException error.
func NewNotFoundException(resource string) *SESv2Error {
	return &SESv2Error{awserrors.NewNotFoundException(resource)}
}

// NewBadRequestException creates a new BadRequestException error.
func NewBadRequestException(message string) *SESv2Error {
	return &SESv2Error{awserrors.NewBadRequestException(message)}
}

// NewMessageRejectedException creates a new MessageRejectedException error.
func NewMessageRejectedException(message string) *SESv2Error {
	return &SESv2Error{awserrors.NewAWSError("MessageRejected", message, http.StatusBadRequest)}
}

// NewLimitExceededException creates a new LimitExceededException error.
func NewLimitExceededException(message string) *SESv2Error {
	return &SESv2Error{awserrors.NewLimitExceededException(message)}
}

// NewInvalidParameterException creates a new InvalidParameterException error.
func NewInvalidParameterException(message string) *SESv2Error {
	return &SESv2Error{awserrors.NewBadRequestException(message)}
}
