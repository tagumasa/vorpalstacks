package appsync

import (
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// createDomainNameInput carries the parsed CreateDomainName request payload.
type createDomainNameInput struct {
	DomainName     string
	CertificateArn string
	Description    string
	Tags           map[string]string
}

// updateDomainNameInput carries the parsed UpdateDomainName request payload.
type updateDomainNameInput struct {
	DomainName  string
	Description string
	Tags        map[string]string
}

// createDomainNameCore validates the request and registers a custom domain
// name for AppSync.
func (s *AppSyncService) createDomainNameCore(store *appsyncstore.AppSyncStore, in createDomainNameInput) (*appsyncstore.DomainNameConfig, error) {
	if in.DomainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}
	if err := validateDomainName(in.DomainName); err != nil {
		return nil, err
	}
	if in.CertificateArn == "" {
		return nil, NewBadRequestException("certificateArn is required")
	}
	if err := validateCertificateArn(in.CertificateArn); err != nil {
		return nil, err
	}

	if err := validateDescription(in.Description); err != nil {
		return nil, err
	}

	config := &appsyncstore.DomainNameConfig{
		DomainName:        in.DomainName,
		CertificateArn:    in.CertificateArn,
		Description:       in.Description,
		AppsyncDomainName: in.DomainName + ".appsync-api." + store.GetRegion() + ".amazonaws.com",
		DomainNameArn:     store.BuildDomainNameARN(in.DomainName),
		HostedZoneId:      cloudFrontHostedZoneID,
		// Tags must be parsed from the request and persisted at creation time.
		Tags: in.Tags,
	}

	if err := store.CreateDomainName(config); err != nil {
		return nil, mapStoreErrorE(err)
	}

	return config, nil
}

// listDomainNamesCore lists custom domain names with pagination.
func (s *AppSyncService) listDomainNamesCore(store *appsyncstore.AppSyncStore, maxResults int, nextToken string) ([]*appsyncstore.DomainNameConfig, string, error) {
	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	configs, nextToken, err := store.ListDomainNames(opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	return configs, nextToken, nil
}

// getDomainNameCore fetches a custom domain name configuration.
func (s *AppSyncService) getDomainNameCore(store *appsyncstore.AppSyncStore, domainName string) (*appsyncstore.DomainNameConfig, error) {
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	config, err := store.GetDomainName(domainName)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return config, nil
}

// updateDomainNameCore applies description/tag updates to an existing custom
// domain name.
func (s *AppSyncService) updateDomainNameCore(store *appsyncstore.AppSyncStore, in updateDomainNameInput) (*appsyncstore.DomainNameConfig, error) {
	if in.DomainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	config, err := store.GetDomainName(in.DomainName)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	if in.Description != "" {
		if err := validateDescription(in.Description); err != nil {
			return nil, err
		}
		config.Description = in.Description
	}

	// Update tags if provided in the request.
	if len(in.Tags) > 0 {
		config.Tags = in.Tags
	}

	if err := store.UpdateDomainName(config); err != nil {
		return nil, mapStoreErrorE(err)
	}

	return config, nil
}

// deleteDomainNameCore removes a custom domain name, disassociating any API
// first to prevent dangling references.
func (s *AppSyncService) deleteDomainNameCore(store *appsyncstore.AppSyncStore, domainName string) error {
	if domainName == "" {
		return NewBadRequestException("domainName is required")
	}

	// Disassociate API before deleting the domain to prevent dangling references.
	if assoc, err := store.GetApiAssociation(domainName); err == nil && assoc != nil {
		_ = store.DisassociateApi(domainName)
	}

	if err := store.DeleteDomainName(domainName); err != nil {
		return mapStoreErrorE(err)
	}

	return nil
}

// associateApiCore associates a GraphQL API with a custom domain name.
func (s *AppSyncService) associateApiCore(store *appsyncstore.AppSyncStore, domainName, apiId string) (*appsyncstore.ApiAssociation, error) {
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	if _, err := store.GetDomainName(domainName); err != nil {
		return nil, mapStoreErrorE(err)
	}

	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if _, err := store.GetGraphqlApiById(apiId); err != nil {
		return nil, mapStoreErrorE(err)
	}

	assoc := &appsyncstore.ApiAssociation{
		ApiId:             apiId,
		DomainName:        domainName,
		AssociationStatus: "SUCCESS",
	}

	if err := store.AssociateApi(assoc); err != nil {
		return nil, mapStoreErrorE(err)
	}

	return assoc, nil
}

// disassociateApiCore disassociates a GraphQL API from a custom domain name.
func (s *AppSyncService) disassociateApiCore(store *appsyncstore.AppSyncStore, domainName string) error {
	if domainName == "" {
		return NewBadRequestException("domainName is required")
	}

	if err := store.DisassociateApi(domainName); err != nil {
		return mapStoreErrorE(err)
	}

	return nil
}

// getApiAssociationCore fetches the API association for a domain name.
func (s *AppSyncService) getApiAssociationCore(store *appsyncstore.AppSyncStore, domainName string) (*appsyncstore.ApiAssociation, error) {
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	assoc, err := store.GetApiAssociation(domainName)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return assoc, nil
}
