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
// The wire request carries only domainName and description; tag changes go
// through the tag operations.
type updateDomainNameInput struct {
	DomainName  string
	Description string
}

// cloudFrontHostedZoneID is the fixed Route 53 hosted zone ID for all
// CloudFront distributions. AppSync custom domains are backed by CloudFront,
// so this value is always returned in DomainNameConfig.hostedZoneId.
const cloudFrontHostedZoneID = "Z2FDTNDATAQYW2"

// createDomainNameCore validates the request, registers a custom domain name
// for AppSync, and returns the configuration with its tag-store view.
func (s *AppSyncService) createDomainNameCore(store *appsyncstore.AppSyncStore, in createDomainNameInput) (*appsyncstore.DomainNameConfig, map[string]string, error) {
	if in.DomainName == "" {
		return nil, nil, NewBadRequestException("domainName is required")
	}
	if err := validateDomainName(in.DomainName); err != nil {
		return nil, nil, err
	}
	if in.CertificateArn == "" {
		return nil, nil, NewBadRequestException("certificateArn is required")
	}
	if err := validateCertificateArn(in.CertificateArn); err != nil {
		return nil, nil, err
	}

	if err := validateDescription(in.Description); err != nil {
		return nil, nil, err
	}

	config := &appsyncstore.DomainNameConfig{
		DomainName:        in.DomainName,
		CertificateArn:    in.CertificateArn,
		Description:       in.Description,
		AppsyncDomainName: in.DomainName + ".appsync-api." + store.GetRegion() + ".amazonaws.com",
		DomainNameArn:     store.BuildDomainNameARN(in.DomainName),
		HostedZoneId:      cloudFrontHostedZoneID,
	}

	if err := store.CreateDomainName(config); err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	if err := applyCreateTags(store, config.DomainNameArn, in.Tags); err != nil {
		return nil, nil, err
	}

	return config, listTagsIfAny(store, config.DomainNameArn), nil
}

// domainNameWithTags pairs a domain configuration with its tag-store view.
type domainNameWithTags struct {
	Config *appsyncstore.DomainNameConfig
	Tags   map[string]string
}

// listDomainNamesCore lists custom domain names with pagination, pairing each
// configuration with its tag-store view.
func (s *AppSyncService) listDomainNamesCore(store *appsyncstore.AppSyncStore, maxResults int, nextToken string) ([]*domainNameWithTags, string, error) {
	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	configs, nextToken, err := store.ListDomainNames(opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	entries := make([]*domainNameWithTags, len(configs))
	for i, c := range configs {
		entries[i] = &domainNameWithTags{Config: c, Tags: listTagsIfAny(store, c.DomainNameArn)}
	}
	return entries, nextToken, nil
}

// getDomainNameCore fetches a custom domain name configuration together with
// its tag-store view.
func (s *AppSyncService) getDomainNameCore(store *appsyncstore.AppSyncStore, domainName string) (*appsyncstore.DomainNameConfig, map[string]string, error) {
	if domainName == "" {
		return nil, nil, NewBadRequestException("domainName is required")
	}

	config, err := store.GetDomainName(domainName)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	return config, listTagsIfAny(store, config.DomainNameArn), nil
}

// updateDomainNameCore applies a description update to an existing custom
// domain name and returns the configuration with its tag-store view.
func (s *AppSyncService) updateDomainNameCore(store *appsyncstore.AppSyncStore, in updateDomainNameInput) (*appsyncstore.DomainNameConfig, map[string]string, error) {
	if in.DomainName == "" {
		return nil, nil, NewBadRequestException("domainName is required")
	}

	config, err := store.GetDomainName(in.DomainName)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	if in.Description != "" {
		if err := validateDescription(in.Description); err != nil {
			return nil, nil, err
		}
		config.Description = in.Description
	}

	if err := store.UpdateDomainName(config); err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	return config, listTagsIfAny(store, config.DomainNameArn), nil
}

// deleteDomainNameCore removes a custom domain name, disassociating any API
// first to prevent dangling references.
func (s *AppSyncService) deleteDomainNameCore(store *appsyncstore.AppSyncStore, domainName string) error {
	if domainName == "" {
		return NewBadRequestException("domainName is required")
	}

	// Disassociate API before deleting the domain to prevent dangling
	// references; a failed disassociation aborts the delete so the
	// association cannot outlive the domain silently.
	if assoc, err := store.GetApiAssociation(domainName); err == nil && assoc != nil {
		if err := store.DisassociateApi(domainName); err != nil {
			return mapStoreErrorE(err)
		}
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
