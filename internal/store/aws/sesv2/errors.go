package sesv2

import (
	"fmt"

	storecommon "vorpalstacks/internal/store/aws/common"
)

// SESv2StoreError represents an error that occurs during SESv2 store operations.
type SESv2StoreError struct {
	*storecommon.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *SESv2StoreError) Unwrap() error {
	return e.AWSError
}

var (
	// ErrIdentityNotFound is returned when the specified email identity
	// does not exist.
	ErrIdentityNotFound = &SESv2StoreError{storecommon.NewAWSError("NotFoundException", "Email identity not found", 404)}

	// ErrIdentityAlreadyExists is returned when attempting to create an identity
	// that already exists.
	ErrIdentityAlreadyExists = &SESv2StoreError{storecommon.NewAWSError("AlreadyExistsException", "Email identity already exists", 400)}

	// ErrConfigSetNotFound is returned when the specified configuration set
	// does not exist.
	ErrConfigSetNotFound = &SESv2StoreError{storecommon.NewAWSError("NotFoundException", "Configuration set not found", 404)}

	// ErrConfigSetAlreadyExists is returned when attempting to create a
	// configuration set that already exists.
	ErrConfigSetAlreadyExists = &SESv2StoreError{storecommon.NewAWSError("AlreadyExistsException", "Configuration set already exists", 400)}

	// ErrTemplateNotFound is returned when the specified email template
	// does not exist.
	ErrTemplateNotFound = &SESv2StoreError{storecommon.NewAWSError("NotFoundException", "Email template not found", 404)}

	// ErrTemplateAlreadyExists is returned when attempting to create an email
	// template that already exists.
	ErrTemplateAlreadyExists = &SESv2StoreError{storecommon.NewAWSError("AlreadyExistsException", "Email template already exists", 400)}

	// ErrEmailNotFound is returned when the specified email does not exist.
	ErrEmailNotFound = &SESv2StoreError{storecommon.NewAWSError("NotFoundException", "Email not found", 404)}

	// ErrInvalidParameter is returned when a parameter is not valid.
	ErrInvalidParameter = &SESv2StoreError{storecommon.NewAWSError("BadRequestException", "Invalid parameter", 400)}

	// ErrLimitExceeded is returned when a service limit has been exceeded.
	ErrLimitExceeded = &SESv2StoreError{storecommon.NewAWSError("LimitExceededException", "Limit exceeded", 400)}

	// ErrContactListNotFound is returned when the specified contact list does not exist.
	ErrContactListNotFound = &SESv2StoreError{storecommon.NewAWSError("NotFoundException", "Contact list not found", 404)}

	// ErrContactListAlreadyExists is returned when attempting to create a contact list that already exists.
	ErrContactListAlreadyExists = &SESv2StoreError{storecommon.NewAWSError("AlreadyExistsException", "Contact list already exists", 400)}

	// ErrContactNotFound is returned when the specified contact does not exist.
	ErrContactNotFound = &SESv2StoreError{storecommon.NewAWSError("NotFoundException", "Contact not found", 404)}

	// ErrContactAlreadyExists is returned when attempting to create a contact that already exists.
	ErrContactAlreadyExists = &SESv2StoreError{storecommon.NewAWSError("AlreadyExistsException", "Contact already exists", 400)}

	// ErrPolicyNotFound is returned when the specified policy does not exist.
	ErrPolicyNotFound = &SESv2StoreError{storecommon.NewAWSError("NotFoundException", "Policy not found", 404)}

	// ErrDedicatedIpPoolNotFound is returned when the specified IP pool does not exist.
	ErrDedicatedIpPoolNotFound = &SESv2StoreError{storecommon.NewAWSError("NotFoundException", "Dedicated IP pool not found", 404)}

	// ErrDedicatedIpPoolAlreadyExists is returned when attempting to create a pool that already exists.
	ErrDedicatedIpPoolAlreadyExists = &SESv2StoreError{storecommon.NewAWSError("AlreadyExistsException", "Dedicated IP pool already exists", 400)}
)

// NewIdentityNotFoundException creates a new error indicating that the specified
// email identity was not found.
func NewIdentityNotFoundException(identity string) *SESv2StoreError {
	return &SESv2StoreError{storecommon.NewAWSError("NotFoundException", fmt.Sprintf("Email identity '%s' not found", identity), 404)}
}

// NewConfigSetNotFoundException creates a new error indicating that the specified
// configuration set was not found.
func NewConfigSetNotFoundException(name string) *SESv2StoreError {
	return &SESv2StoreError{storecommon.NewAWSError("NotFoundException", fmt.Sprintf("Configuration set '%s' not found", name), 404)}
}

// NewTemplateNotFoundException creates a new error indicating that the specified
// email template was not found.
func NewTemplateNotFoundException(name string) *SESv2StoreError {
	return &SESv2StoreError{storecommon.NewAWSError("NotFoundException", fmt.Sprintf("Email template '%s' not found", name), 404)}
}
