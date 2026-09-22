// Transport-agnostic Core functions for IAM service-specific credentials:
// validation and store operations shared by the AWS-compatible HTTP API
// handlers and any admin plane paths (the xxxCore pattern).
package iam

import (
	"errors"
	"fmt"
	"time"

	"vorpalstacks/internal/common/pagination"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreateServiceSpecificCredentialInput holds the parameters for creating a
// service-specific credential.  A nil CredentialAgeDays leaves the
// credential non-expiring.
type CreateServiceSpecificCredentialInput struct {
	UserName          string
	ServiceName       string
	CredentialAgeDays *int
}

// UpdateServiceSpecificCredentialInput holds the parameters for updating a
// service-specific credential status.  UserName is optional: when set it
// must own the credential.
type UpdateServiceSpecificCredentialInput struct {
	ServiceSpecificCredentialId string
	UserName                    string
	Status                      string
}

// ServiceSpecificCredentialListResult holds a paginated service-specific
// credential listing.
type ServiceSpecificCredentialListResult struct {
	Credentials []*iamstore.ServiceSpecificCredential
	IsTruncated bool
	NextMarker  string
}

// supportedServiceSpecificCredentialServices lists the services that
// support service-specific credentials: CodeCommit (Git HTTPS credentials),
// Amazon Keyspaces (Cassandra credentials) and CloudWatch Logs (ACWL
// bearer tokens for the HTTP ingestion endpoints — the developer guide's
// "Setting up bearer token authentication" issues them with
// --service-name logs.amazonaws.com).
var supportedServiceSpecificCredentialServices = map[string]bool{
	"codecommit.amazonaws.com": true,
	"cassandra.amazonaws.com":  true,
	"logs.amazonaws.com":       true,
}

func supportedServiceSpecificCredentialService(serviceName string) bool {
	return supportedServiceSpecificCredentialServices[serviceName]
}

// createServiceSpecificCredentialCore validates input and generates a new
// service-specific credential for the specified user and service.
func (s *IAMService) createServiceSpecificCredentialCore(store *iamstore.IAMStore, input *CreateServiceSpecificCredentialInput) (*iamstore.ServiceSpecificCredential, error) {
	if input.UserName == "" {
		return nil, NewValidationError("UserName")
	}
	serviceName := input.ServiceName
	if serviceName == "" {
		return nil, NewValidationError("ServiceName")
	}
	if !validateServiceNamespace(serviceName) {
		return nil, NewInvalidInputError("ServiceName", fmt.Sprintf("must be 1 to %d characters", MaxServiceNamespaceLength))
	}
	// Only services that support service-specific credentials accept them;
	// any other service name fails with NotSupportedService.
	if !supportedServiceSpecificCredentialService(serviceName) {
		return nil, ErrNotSupportedService
	}

	// CredentialAgeDays (documented range [1, 36600]). When not specified the
	// credential does not expire.
	credentialAgeDays := 0
	if input.CredentialAgeDays != nil {
		credentialAgeDays = *input.CredentialAgeDays
		if credentialAgeDays < MinCredentialAgeDays || credentialAgeDays > MaxCredentialAgeDays {
			return nil, NewInvalidInputError("CredentialAgeDays", fmt.Sprintf("must be between %d and %d", MinCredentialAgeDays, MaxCredentialAgeDays))
		}
	}

	if !store.Users().Exists(input.UserName) {
		return nil, NewNoSuchUserError(input.UserName)
	}

	return store.ServiceSpecificCredentials().Create(input.UserName, serviceName, credentialAgeDays)
}

// deleteServiceSpecificCredentialCore validates input and deletes the
// specified service-specific credential.  A non-empty userName that does
// not own the credential yields NoSuchEntity.
func (s *IAMService) deleteServiceSpecificCredentialCore(store *iamstore.IAMStore, credentialId, userName string) error {
	if credentialId == "" {
		return NewValidationError("ServiceSpecificCredentialId")
	}

	if userName != "" && !store.Users().Exists(userName) {
		return NewNoSuchUserError(userName)
	}
	cred, err := store.ServiceSpecificCredentials().Get(credentialId)
	if err != nil {
		return storeReadError(err, iamstore.ErrServiceSpecificCredentialNotFound, NewNoSuchEntityError("service-specific credential", credentialId))
	}
	// A named user that does not own the credential yields NoSuchEntity.
	if userName != "" && cred.UserName != userName {
		return NewNoSuchEntityError("service-specific credential", credentialId)
	}

	return store.ServiceSpecificCredentials().Delete(credentialId)
}

// listServiceSpecificCredentialsCore returns a paginated list of the
// service-specific credentials for the specified user, optionally filtered
// by service name.
func (s *IAMService) listServiceSpecificCredentialsCore(store *iamstore.IAMStore, userName, serviceName, marker string, maxItems int) (*ServiceSpecificCredentialListResult, error) {
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	creds, err := store.ServiceSpecificCredentials().ListByUserName(userName)
	if err != nil {
		return nil, err
	}

	filtered := make([]*iamstore.ServiceSpecificCredential, 0, len(creds))
	for _, cred := range creds {
		if serviceName != "" && cred.ServiceName != serviceName {
			continue
		}
		filtered = append(filtered, cred)
	}

	paged := pagination.PaginateSlice(filtered, marker, maxItems, func(c *iamstore.ServiceSpecificCredential) string {
		return c.ServiceSpecificCredentialId
	})

	return &ServiceSpecificCredentialListResult{
		Credentials: paged.Items,
		IsTruncated: paged.IsTruncated,
		NextMarker:  paged.NextMarker,
	}, nil
}

// resetServiceSpecificCredentialCore validates input and resets the
// password for a service-specific credential.  A non-empty userName that
// does not own the credential yields NoSuchEntity.
func (s *IAMService) resetServiceSpecificCredentialCore(store *iamstore.IAMStore, credentialId, userName string) (*iamstore.ServiceSpecificCredential, error) {
	if credentialId == "" {
		return nil, NewValidationError("ServiceSpecificCredentialId")
	}

	if userName != "" && !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}
	cred, err := store.ServiceSpecificCredentials().Get(credentialId)
	if err != nil {
		return nil, storeReadError(err, iamstore.ErrServiceSpecificCredentialNotFound, NewNoSuchEntityError("service-specific credential", credentialId))
	}
	// A named user that does not own the credential yields NoSuchEntity.
	if userName != "" && cred.UserName != userName {
		return nil, NewNoSuchEntityError("service-specific credential", credentialId)
	}

	cred, err = store.ServiceSpecificCredentials().ResetPassword(credentialId)
	if err != nil {
		return nil, storeReadError(err, iamstore.ErrServiceSpecificCredentialNotFound, NewNoSuchEntityError("service-specific credential", credentialId))
	}
	return cred, nil
}

// updateServiceSpecificCredentialCore validates input and sets the status
// of a service-specific credential to Active or Inactive.  A non-empty
// userName that does not own the credential yields NoSuchEntity.
func (s *IAMService) updateServiceSpecificCredentialCore(store *iamstore.IAMStore, input *UpdateServiceSpecificCredentialInput) error {
	credentialId := input.ServiceSpecificCredentialId
	if credentialId == "" {
		return NewValidationError("ServiceSpecificCredentialId")
	}
	status := input.Status
	if status == "" {
		return NewValidationError("Status")
	}
	if status != "Active" && status != "Inactive" {
		return NewInvalidInputError("Status", "must be Active or Inactive")
	}

	userName := input.UserName
	if userName != "" && !store.Users().Exists(userName) {
		return NewNoSuchUserError(userName)
	}
	cred, err := store.ServiceSpecificCredentials().Get(credentialId)
	if err != nil {
		return storeReadError(err, iamstore.ErrServiceSpecificCredentialNotFound, NewNoSuchEntityError("service-specific credential", credentialId))
	}
	// A named user that does not own the credential yields NoSuchEntity.
	if userName != "" && cred.UserName != userName {
		return NewNoSuchEntityError("service-specific credential", credentialId)
	}

	return store.ServiceSpecificCredentials().UpdateStatus(credentialId, status)
}

// LogIngestionServiceName is the service name of the CloudWatch Logs ACWL
// bearer tokens: a service-specific credential of this service
// authenticates the CloudWatch Logs HTTP ingestion endpoints.
const LogIngestionServiceName = "logs.amazonaws.com"

// AuthenticateLogIngestionToken validates an ACWL bearer token presented
// to the CloudWatch Logs HTTP ingestion endpoints: the token must be the
// ServiceCredentialSecret of an Active, unexpired service-specific
// credential of service logs.amazonaws.com. A nil error authenticates;
// every failure carries a short reason for the 401 response.
func (s *IAMService) AuthenticateLogIngestionToken(token string) error {
	store, err := s.GetStoreForRegion("")
	if err != nil {
		return fmt.Errorf("IAM store unavailable: %w", err)
	}
	cred, err := store.ServiceSpecificCredentials().FindByServiceAndSecret(LogIngestionServiceName, token)
	if err != nil {
		return errors.New("invalid bearer token")
	}
	if cred.Status != "Active" {
		return errors.New("bearer token is inactive")
	}
	if cred.ExpirationDate != nil && time.Now().After(*cred.ExpirationDate) {
		return errors.New("bearer token has expired")
	}
	return nil
}
