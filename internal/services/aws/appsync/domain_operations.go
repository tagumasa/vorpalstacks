package appsync

import (
	"context"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// cloudFrontHostedZoneID is the fixed Route 53 hosted zone ID for all
// CloudFront distributions. AppSync custom domains are backed by CloudFront,
// so this value is always returned in DomainNameConfig.hostedZoneId.
const cloudFrontHostedZoneID = "Z2FDTNDATAQYW2"

// CreateDomainName creates a custom domain name for AppSync.
func (s *AppSyncService) CreateDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}
	if err := validateDomainName(domainName); err != nil {
		return nil, err
	}
	certificateArn := request.GetStringParam(req.Parameters, "certificateArn")
	if certificateArn == "" {
		return nil, NewBadRequestException("certificateArn is required")
	}
	if err := validateCertificateArn(certificateArn); err != nil {
		return nil, err
	}

	description := request.GetStringParam(req.Parameters, "description")
	if err := validateDescription(description); err != nil {
		return nil, err
	}

	tagMap, err := parseTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	config := &appsyncstore.DomainNameConfig{
		DomainName:        domainName,
		CertificateArn:    certificateArn,
		Description:       description,
		AppsyncDomainName: domainName + ".appsync-api." + store.GetRegion() + ".amazonaws.com",
		DomainNameArn:     store.BuildDomainNameARN(domainName),
		HostedZoneId:      cloudFrontHostedZoneID,
		// Tags must be parsed from the request and persisted at creation time.
		Tags: tagMap,
	}

	if err := store.CreateDomainName(config); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"domainNameConfig": domainNameConfigToMap(config)}, nil
}

// ListDomainNames lists all custom domain names.
func (s *AppSyncService) ListDomainNames(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	opts, err := parsePaginationOptions(req)
	if err != nil {
		return nil, err
	}
	configs, nextToken, err := store.ListDomainNames(opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]map[string]interface{}, 0, len(configs))
	for _, c := range configs {
		items = append(items, domainNameConfigToMap(c))
	}

	response := map[string]interface{}{"domainNameConfigs": items}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

// GetDomainName retrieves a custom domain name configuration.
func (s *AppSyncService) GetDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	config, err := store.GetDomainName(domainName)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"domainNameConfig": domainNameConfigToMap(config)}, nil
}

// UpdateDomainName updates a custom domain name description.
func (s *AppSyncService) UpdateDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	config, err := store.GetDomainName(domainName)
	if err != nil {
		return mapStoreError(err)
	}

	description := request.GetStringParam(req.Parameters, "description")
	if description != "" {
		if err := validateDescription(description); err != nil {
			return nil, err
		}
		config.Description = description
	}

	// Update tags if provided in the request.
	tagMap, err := parseTags(req.Parameters)
	if err != nil {
		return nil, err
	}
	if len(tagMap) > 0 {
		config.Tags = tagMap
	}

	if err := store.UpdateDomainName(config); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"domainNameConfig": domainNameConfigToMap(config)}, nil
}

// DeleteDomainName deletes a custom domain name.
func (s *AppSyncService) DeleteDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	// Disassociate API before deleting the domain to prevent dangling references.
	if assoc, err := store.GetApiAssociation(domainName); err == nil && assoc != nil {
		_ = store.DisassociateApi(domainName)
	}

	if err := store.DeleteDomainName(domainName); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// AssociateApi associates a GraphQL API with a custom domain name.
func (s *AppSyncService) AssociateApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	if _, err := store.GetDomainName(domainName); err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if _, err := store.GetGraphqlApiById(apiId); err != nil {
		return mapStoreError(err)
	}

	assoc := &appsyncstore.ApiAssociation{
		ApiId:             apiId,
		DomainName:        domainName,
		AssociationStatus: "SUCCESS",
	}

	if err := store.AssociateApi(assoc); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"apiAssociation": apiAssociationToMap(assoc)}, nil
}

// DisassociateApi disassociates a GraphQL API from a custom domain name.
func (s *AppSyncService) DisassociateApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	if err := store.DisassociateApi(domainName); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// GetApiAssociation retrieves the API association for a domain name.
func (s *AppSyncService) GetApiAssociation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	assoc, err := store.GetApiAssociation(domainName)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"apiAssociation": apiAssociationToMap(assoc)}, nil
}
