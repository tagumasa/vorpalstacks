package appsync

import (
	"context"

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

	tagMap, err := parseTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	config, err := s.createDomainNameCore(store, createDomainNameInput{
		DomainName:     request.GetStringParam(req.Parameters, "domainName"),
		CertificateArn: request.GetStringParam(req.Parameters, "certificateArn"),
		Description:    request.GetStringParam(req.Parameters, "description"),
		Tags:           tagMap,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"domainNameConfig": domainNameConfigToMap(config)}, nil
}

// ListDomainNames lists all custom domain names.
func (s *AppSyncService) ListDomainNames(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	configs, nextToken, err := s.listDomainNamesCore(store, request.GetIntParam(req.Parameters, "maxResults"), request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
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

	config, err := s.getDomainNameCore(store, request.GetStringParam(req.Parameters, "domainName"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"domainNameConfig": domainNameConfigToMap(config)}, nil
}

// UpdateDomainName updates a custom domain name description.
func (s *AppSyncService) UpdateDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	tagMap, err := parseTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	config, err := s.updateDomainNameCore(store, updateDomainNameInput{
		DomainName:  request.GetStringParam(req.Parameters, "domainName"),
		Description: request.GetStringParam(req.Parameters, "description"),
		Tags:        tagMap,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"domainNameConfig": domainNameConfigToMap(config)}, nil
}

// DeleteDomainName deletes a custom domain name.
func (s *AppSyncService) DeleteDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	if err := s.deleteDomainNameCore(store, request.GetStringParam(req.Parameters, "domainName")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// AssociateApi associates a GraphQL API with a custom domain name.
func (s *AppSyncService) AssociateApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	assoc, err := s.associateApiCore(store, request.GetStringParam(req.Parameters, "domainName"), request.GetStringParam(req.Parameters, "apiId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"apiAssociation": apiAssociationToMap(assoc)}, nil
}

// DisassociateApi disassociates a GraphQL API from a custom domain name.
func (s *AppSyncService) DisassociateApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	if err := s.disassociateApiCore(store, request.GetStringParam(req.Parameters, "domainName")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// GetApiAssociation retrieves the API association for a domain name.
func (s *AppSyncService) GetApiAssociation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	assoc, err := s.getApiAssociationCore(store, request.GetStringParam(req.Parameters, "domainName"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"apiAssociation": apiAssociationToMap(assoc)}, nil
}
