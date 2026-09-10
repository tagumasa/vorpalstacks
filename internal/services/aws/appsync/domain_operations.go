package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

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

	config, tags, err := s.createDomainNameCore(store, createDomainNameInput{
		DomainName:     request.GetStringParam(req.Parameters, "domainName"),
		CertificateArn: request.GetStringParam(req.Parameters, "certificateArn"),
		Description:    request.GetStringParam(req.Parameters, "description"),
		Tags:           tagMap,
	})
	if err != nil {
		return nil, err
	}

	m := domainNameConfigToMap(config)
	if len(tags) > 0 {
		m["tags"] = tags
	}
	return map[string]interface{}{"domainNameConfig": m}, nil
}

// ListDomainNames lists all custom domain names.
func (s *AppSyncService) ListDomainNames(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	entries, nextToken, err := s.listDomainNamesCore(store, request.GetIntParam(req.Parameters, "maxResults"), request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		m := domainNameConfigToMap(entry.Config)
		if len(entry.Tags) > 0 {
			m["tags"] = entry.Tags
		}
		items = append(items, m)
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

	config, tags, err := s.getDomainNameCore(store, request.GetStringParam(req.Parameters, "domainName"))
	if err != nil {
		return nil, err
	}

	m := domainNameConfigToMap(config)
	if len(tags) > 0 {
		m["tags"] = tags
	}
	return map[string]interface{}{"domainNameConfig": m}, nil
}

// UpdateDomainName updates a custom domain name description. Tag changes go
// through the tag operations, whose store is the single tag source.
func (s *AppSyncService) UpdateDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	config, tags, err := s.updateDomainNameCore(store, updateDomainNameInput{
		DomainName:  request.GetStringParam(req.Parameters, "domainName"),
		Description: request.GetStringParam(req.Parameters, "description"),
	})
	if err != nil {
		return nil, err
	}

	m := domainNameConfigToMap(config)
	if len(tags) > 0 {
		m["tags"] = tags
	}
	return map[string]interface{}{"domainNameConfig": m}, nil
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
