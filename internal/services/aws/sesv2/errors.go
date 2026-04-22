package sesv2

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// Error variables for common SESv2 error conditions.
var (
	// ErrAlreadyExists is returned when a resource already exists.
	ErrAlreadyExists = awserrors.NewAWSError("AlreadyExistsException", "Resource already exists", http.StatusConflict)
	// ErrBadRequest is returned when a request is invalid.
	ErrBadRequest = awserrors.NewBadRequestException("Invalid request")
	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = awserrors.NewNotFoundException("Resource")
	// ErrLimitExceeded is returned when a limit is exceeded.
	ErrLimitExceeded = awserrors.NewLimitExceededException("Limit exceeded")
	// ErrTooManyRequests is returned when too many requests are sent.
	ErrTooManyRequests = awserrors.NewThrottlingException("Too many requests")
	// ErrMessageRejected is returned when a message is rejected.
	ErrMessageRejected = awserrors.NewAWSError("MessageRejected", "Message rejected", http.StatusBadRequest)
	// ErrAccountSuspended is returned when the account is suspended.
	ErrAccountSuspended = awserrors.NewAWSError("AccountSuspendedException", "Account suspended", http.StatusBadRequest)
	// ErrSendingPaused is returned when sending is paused.
	ErrSendingPaused = awserrors.NewAWSError("SendingPausedException", "Sending paused", http.StatusBadRequest)
	// ErrMailFromDomainNotVerified is returned when the MAIL FROM domain is not verified.
	ErrMailFromDomainNotVerified = awserrors.NewAWSError("MailFromDomainNotVerifiedException", "MAIL FROM domain not verified", http.StatusBadRequest)
	// ErrConcurrentModification is returned when a concurrent modification conflict occurs.
	ErrConcurrentModification = awserrors.NewAWSError("ConcurrentModificationException", "Concurrent modification", http.StatusInternalServerError)
	// ErrMissingParameter is returned when a required parameter is missing.
	ErrMissingParameter = awserrors.NewBadRequestException("Missing required parameter")
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = awserrors.NewBadRequestException("Invalid parameter")
	// ErrIdentityNotFound is returned when an email identity is not found.
	ErrIdentityNotFound = awserrors.NewNotFoundException("Email identity")
	// ErrConfigurationSetNotFound is returned when a configuration set is not found.
	ErrConfigurationSetNotFound = awserrors.NewNotFoundException("Configuration set")
	// ErrTemplateNotFound is returned when an email template is not found.
	ErrTemplateNotFound = awserrors.NewNotFoundException("Email template")
	// ErrContactListNotFound is returned when a contact list is not found.
	ErrContactListNotFound = awserrors.NewNotFoundException("Contact list")
	// ErrContactNotFound is returned when a contact is not found.
	ErrContactNotFound = awserrors.NewNotFoundException("Contact")
)

// NewAlreadyExistsException creates a new AlreadyExistsException error.
func NewAlreadyExistsException(resource string) *awserrors.AWSError {
	return awserrors.NewResourceAlreadyExistsException(resource)
}
