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

// UpdateIdentityProviderInput carries the update members of
// UpdateIdentityProvider. The model's UpdateIdentityProviderRequest members
// are AttributeMapping, IdpIdentifiers, ProviderDetails, ProviderName and
// UserPoolId — ProviderType is not an update member, so it has no field
// here. nil maps keep the stored value; non-nil maps replace it. Presence,
// not length, decides for IdpIdentifiers: IdpIdentifiersProvided=true with
// an empty list clears the stored identifiers.
type UpdateIdentityProviderInput struct {
	UserPoolID             string
	ProviderName           string
	ProviderDetails        map[string]string
	AttributeMapping       map[string]string
	IdpIdentifiers         []string
	IdpIdentifiersProvided bool
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
// transport-agnostic input. It is the shared create path of the admin
// console handler and the HTTP CreateIdentityProvider operation.
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

	maxResults := applyListLimitDefaults(in.MaxResults)

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

// updateIdentityProviderCore applies the update members onto the stored
// identity provider and persists it, returning the updated record for
// response serialisation. Members absent from the input keep their stored
// value; present members replace it.
func (s *CognitoService) updateIdentityProviderCore(region string, in UpdateIdentityProviderInput) (*cognitostore.IdentityProvider, error) {
	if in.UserPoolID == "" || in.ProviderName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	ip, err := store.GetIdentityProvider(in.UserPoolID, in.ProviderName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if in.ProviderDetails != nil {
		ip.ProviderDetails = in.ProviderDetails
	}
	if in.AttributeMapping != nil {
		ip.AttributeMapping = in.AttributeMapping
	}
	if in.IdpIdentifiersProvided {
		ip.IdpIdentifiers = in.IdpIdentifiers
	}

	if err := store.UpdateIdentityProvider(ip); err != nil {
		return nil, ErrInternalError
	}
	return ip, nil
}

// getIdentityProviderByIdentifierCore resolves the identity provider that
// declares the given IdpIdentifier.
func (s *CognitoService) getIdentityProviderByIdentifierCore(region, userPoolID, idpIdentifier string) (*cognitostore.IdentityProvider, error) {
	if userPoolID == "" || idpIdentifier == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	providers, err := store.ListIdentityProviders(userPoolID)
	if err != nil {
		return nil, ErrInternalError
	}

	for _, idp := range providers {
		for _, ident := range idp.IdpIdentifiers {
			if ident == idpIdentifier {
				return idp, nil
			}
		}
	}

	return nil, ErrResourceNotFound
}
