package cognitoidentityprovider

// Core functions for the federated provider linking family. The handlers
// extract the wire members; validation and store access live here.

// adminDisableProviderForUserCore detaches the federated provider identity
// from the referenced user.
func (s *CognitoService) adminDisableProviderForUserCore(region, userPoolID string, user map[string]interface{}) error {
	if userPoolID == "" {
		return ErrInvalidParameter
	}
	if user == nil {
		return ErrInvalidParameter
	}

	providerName := getStringParam(user, "ProviderName")
	providerAttrValue := getStringParam(user, "ProviderAttributeValue")
	if providerName == "" || providerAttrValue == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	u, err := store.GetUserByProvider(userPoolID, providerName, providerAttrValue)
	if err != nil {
		return ErrUserNotFound
	}

	u.ProviderName = ""
	u.ProviderAttributeName = ""
	u.ProviderAttributeValue = ""
	if err := store.UpdateUser(u); err != nil {
		return ErrInternalError
	}

	return nil
}

// adminLinkProviderForUserCore links a federated source identity onto the
// local destination user.
func (s *CognitoService) adminLinkProviderForUserCore(region, userPoolID string, destinationUser, sourceUser map[string]interface{}) error {
	if userPoolID == "" {
		return ErrInvalidParameter
	}
	if destinationUser == nil {
		return ErrInvalidParameter
	}
	if sourceUser == nil {
		return ErrInvalidParameter
	}

	destUsername := getStringParam(destinationUser, "ProviderAttributeValue")
	destProviderName := getStringParam(destinationUser, "ProviderName")
	if destProviderName == "" {
		return ErrInvalidParameter
	}
	srcProviderName := getStringParam(sourceUser, "ProviderName")
	srcProviderAttrName := getStringParam(sourceUser, "ProviderAttributeName")
	srcProviderAttrValue := getStringParam(sourceUser, "ProviderAttributeValue")

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	u, err := store.GetUser(userPoolID, destUsername)
	if err != nil {
		return ErrUserNotFound
	}

	u.ProviderName = srcProviderName
	u.ProviderAttributeName = srcProviderAttrName
	u.ProviderAttributeValue = srcProviderAttrValue
	if err := store.UpdateUser(u); err != nil {
		return ErrInternalError
	}

	return nil
}
