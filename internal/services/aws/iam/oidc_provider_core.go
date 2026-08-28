// Transport-agnostic Core functions for IAM OpenID Connect providers:
// validation and store operations shared by the AWS-compatible HTTP API
// handlers and any admin surface (the xxxCore pattern).
package iam

import (
	"errors"
	"unicode/utf8"

	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreateOpenIDConnectProviderInput holds the parameters for creating an
// OpenID Connect provider.
type CreateOpenIDConnectProviderInput struct {
	Url            string
	ThumbprintList []string
	ClientIDList   []string
	Tags           []tags.Tag
}

// UpdateOpenIDConnectProviderThumbprintInput holds the parameters for
// replacing an OpenID Connect provider's thumbprint list.
type UpdateOpenIDConnectProviderThumbprintInput struct {
	ProviderArn    string
	ThumbprintList []string
}

// OpenIDConnectProviderClientIDInput holds the parameters for the
// add/remove client ID operations.
type OpenIDConnectProviderClientIDInput struct {
	ProviderArn string
	ClientID    string
}

// createOpenIDConnectProviderCore validates input and creates an OpenID
// Connect provider, returning its ARN.
func (s *IAMService) createOpenIDConnectProviderCore(store *iamstore.IAMStore, input *CreateOpenIDConnectProviderInput) (string, error) {
	if input.Url == "" {
		return "", NewValidationError("Url")
	}
	// OpenIDConnectProviderUrlType @length(1,255) counts Unicode characters
	// (the shape carries no pattern).
	if utf8.RuneCountInString(input.Url) > 255 {
		return "", NewInvalidInputError("Url", "must be 1 to 255 characters")
	}

	for _, tp := range input.ThumbprintList {
		if !validateThumbprint(tp) {
			return "", NewInvalidInputError("ThumbprintList", "each thumbprint must be a 40-character hex-encoded SHA-1 hash")
		}
	}
	for _, cid := range input.ClientIDList {
		if !validateClientID(cid) {
			return "", NewInvalidInputError("ClientIDList", "each client ID must be 1 to 255 characters")
		}
	}
	if err := validateNewTags(input.Tags); err != nil {
		return "", err
	}

	provider, err := store.OpenIDConnectProviders().Create(input.Url, input.ThumbprintList, input.ClientIDList, input.Tags)
	if err != nil {
		if errors.Is(err, iamstore.ErrOpenIDConnectProviderAlreadyExists) {
			return "", NewEntityAlreadyExistsError("OpenID Connect Provider " + input.Url)
		}
		return "", err
	}
	return provider.Arn, nil
}

// getOpenIDConnectProviderCore returns the OpenID Connect provider with the
// given ARN.
func (s *IAMService) getOpenIDConnectProviderCore(store *iamstore.IAMStore, providerArn string) (*iamstore.OpenIDConnectProvider, error) {
	if providerArn == "" {
		return nil, NewValidationError("OpenIDConnectProviderArn")
	}
	if err := validateARNParameter("OpenIDConnectProviderArn", providerArn); err != nil {
		return nil, err
	}
	provider, err := store.OpenIDConnectProviders().Get(providerArn)
	if err != nil {
		return nil, NewNoSuchEntityError("OpenID Connect provider", providerArn)
	}
	return provider, nil
}

// listOpenIDConnectProvidersCore lists the OpenID Connect providers in the
// account.
func (s *IAMService) listOpenIDConnectProvidersCore(store *iamstore.IAMStore) ([]*iamstore.OpenIDConnectProvider, error) {
	result, err := store.OpenIDConnectProviders().List()
	if err != nil {
		return nil, err
	}
	return result.OpenIDConnectProviderList, nil
}

// updateOpenIDConnectProviderThumbprintCore validates input and replaces the
// thumbprint list of the given OpenID Connect provider.
func (s *IAMService) updateOpenIDConnectProviderThumbprintCore(store *iamstore.IAMStore, input *UpdateOpenIDConnectProviderThumbprintInput) error {
	if input.ProviderArn == "" {
		return NewValidationError("OpenIDConnectProviderArn")
	}
	if err := validateARNParameter("OpenIDConnectProviderArn", input.ProviderArn); err != nil {
		return err
	}
	for _, tp := range input.ThumbprintList {
		if !validateThumbprint(tp) {
			return NewInvalidInputError("ThumbprintList", "each thumbprint must be a 40-character hex-encoded SHA-1 hash")
		}
	}

	if !store.OpenIDConnectProviders().Exists(input.ProviderArn) {
		return NewNoSuchEntityError("OpenID Connect provider", input.ProviderArn)
	}
	return store.OpenIDConnectProviders().Update(input.ProviderArn, input.ThumbprintList, nil)
}

// addClientIDToOpenIDConnectProviderCore adds a client ID to the given
// OpenID Connect provider.
func (s *IAMService) addClientIDToOpenIDConnectProviderCore(store *iamstore.IAMStore, input *OpenIDConnectProviderClientIDInput) error {
	if input.ProviderArn == "" {
		return NewValidationError("OpenIDConnectProviderArn")
	}
	if err := validateARNParameter("OpenIDConnectProviderArn", input.ProviderArn); err != nil {
		return err
	}
	if input.ClientID == "" {
		return NewValidationError("ClientID")
	}

	// Atomic AddClientID holds the per-ARN lock across the
	// read-modify-write cycle, preventing lost updates from concurrent
	// callers.
	if err := store.OpenIDConnectProviders().AddClientID(input.ProviderArn, input.ClientID); err != nil {
		if errors.Is(err, iamstore.ErrOpenIDConnectProviderNotFound) {
			return NewNoSuchEntityError("OpenID Connect provider", input.ProviderArn)
		}
		if errors.Is(err, iamstore.ErrOpenIDConnectProviderClientIDExists) {
			return NewInvalidInputError("ClientID", "already exists")
		}
		return err
	}
	return nil
}

// removeClientIDFromOpenIDConnectProviderCore removes a client ID from the
// given OpenID Connect provider.
func (s *IAMService) removeClientIDFromOpenIDConnectProviderCore(store *iamstore.IAMStore, input *OpenIDConnectProviderClientIDInput) error {
	if input.ProviderArn == "" {
		return NewValidationError("OpenIDConnectProviderArn")
	}
	if err := validateARNParameter("OpenIDConnectProviderArn", input.ProviderArn); err != nil {
		return err
	}
	if input.ClientID == "" {
		return NewValidationError("ClientID")
	}

	// Atomic RemoveClientID mirrors AddClientID's locking guarantee.
	if err := store.OpenIDConnectProviders().RemoveClientID(input.ProviderArn, input.ClientID); err != nil {
		if errors.Is(err, iamstore.ErrOpenIDConnectProviderNotFound) {
			return NewNoSuchEntityError("OpenID Connect provider", input.ProviderArn)
		}
		if errors.Is(err, iamstore.ErrOpenIDConnectProviderClientIDNotFound) {
			return NewNoSuchEntityError("ClientID", input.ClientID)
		}
		return err
	}
	return nil
}

// deleteOpenIDConnectProviderCore deletes the OpenID Connect provider with
// the given ARN.
func (s *IAMService) deleteOpenIDConnectProviderCore(store *iamstore.IAMStore, providerArn string) error {
	if providerArn == "" {
		return NewValidationError("OpenIDConnectProviderArn")
	}
	if err := validateARNParameter("OpenIDConnectProviderArn", providerArn); err != nil {
		return err
	}
	if !store.OpenIDConnectProviders().Exists(providerArn) {
		return NewNoSuchEntityError("OpenID Connect provider", providerArn)
	}
	return store.OpenIDConnectProviders().Delete(providerArn)
}
