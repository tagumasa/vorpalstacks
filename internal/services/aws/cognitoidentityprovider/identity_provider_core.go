package cognitoidentityprovider

import (
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// CreateIdentityProviderInput carries the fields needed to create a Cognito
// identity provider in a wire-protocol-independent format.
type CreateIdentityProviderInput struct {
	UserPoolID       string
	ProviderName     string
	ProviderType     string
	ProviderDetails  map[string]string
	AttributeMapping map[string]string
	IdpIdentifiers   []string
}

// ListIdentityProvidersInput carries pagination parameters.
type ListIdentityProvidersInput struct {
	UserPoolID string
	MaxResults int
	NextToken  string
}

// ListIdentityProvidersResult is the paginated result.
type ListIdentityProvidersResult struct {
	Providers []*cognitostore.IdentityProvider
	NextToken string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// createIdentityProviderFromInputCore creates an identity provider from
// transport-agnostic input. Used by the admin handler.
func (s *CognitoService) createIdentityProviderFromInputCore(region string, in CreateIdentityProviderInput) (*cognitostore.IdentityProvider, error) {
	if in.UserPoolID == "" || in.ProviderName == "" || in.ProviderType == "" {
		return nil, ErrInvalidParameter
	}
	if !validateProviderType(in.ProviderType) {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetUserPool(in.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	ip := &cognitostore.IdentityProvider{
		UserPoolID:       in.UserPoolID,
		ProviderName:     in.ProviderName,
		ProviderType:     in.ProviderType,
		ProviderDetails:  in.ProviderDetails,
		AttributeMapping: in.AttributeMapping,
		IdpIdentifiers:   in.IdpIdentifiers,
	}

	if err := store.CreateIdentityProvider(ip); err != nil {
		return nil, err
	}
	return ip, nil
}

// describeIdentityProviderCore retrieves an identity provider by name.
func (s *CognitoService) describeIdentityProviderCore(region, userPoolID, providerName string) (*cognitostore.IdentityProvider, error) {
	if userPoolID == "" || providerName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	ip, err := store.GetIdentityProvider(userPoolID, providerName)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return ip, nil
}

// deleteIdentityProviderCore deletes an identity provider by name.
func (s *CognitoService) deleteIdentityProviderCore(region, userPoolID, providerName string) error {
	if userPoolID == "" || providerName == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteIdentityProvider(userPoolID, providerName); err != nil {
		return ErrResourceNotFound
	}
	return nil
}

// listIdentityProvidersCore lists identity providers in a user pool.
func (s *CognitoService) listIdentityProvidersCore(region string, in ListIdentityProvidersInput) (*ListIdentityProvidersResult, error) {
	if in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	maxResults := in.MaxResults
	if maxResults <= 0 || maxResults > listLimitMax {
		maxResults = listLimitMax
	}

	result, err := store.ListIdentityProvidersPaginated(in.UserPoolID, storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   in.NextToken,
	})
	if err != nil {
		return nil, err
	}

	return &ListIdentityProvidersResult{
		Providers: result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// updateIdentityProviderCore persists updates to an identity provider.
func (s *CognitoService) updateIdentityProviderCore(region string, ip *cognitostore.IdentityProvider) error {
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	if err := store.UpdateIdentityProvider(ip); err != nil {
		return ErrInternalError
	}
	return nil
}
