package appsync

import (
	"context"

	"github.com/google/uuid"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

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
	certificateArn := request.GetStringParam(req.Parameters, "certificateArn")
	if certificateArn == "" {
		return nil, NewBadRequestException("certificateArn is required")
	}

	description := request.GetStringParam(req.Parameters, "description")

	config := &appsyncstore.DomainNameConfig{
		DomainName:        domainName,
		CertificateArn:    certificateArn,
		Description:       description,
		AppsyncDomainName: domainName + ".appsync-api." + store.GetRegion() + ".amazonaws.com",
		DomainNameArn:     store.BuildDomainNameARN(domainName),
		HostedZoneId:      uuid.New().String(),
		// Tags must be parsed from the request and persisted at creation time.
		Tags: parseTags(req.Parameters),
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

	opts := parsePaginationOptions(req)
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
		config.Description = description
	}

	// Update tags if provided in the request.
	tags := parseTags(req.Parameters)
	if len(tags) > 0 {
		config.Tags = tags
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

func domainNameConfigToMap(c *appsyncstore.DomainNameConfig) map[string]interface{} {
	result := map[string]interface{}{
		"domainName":        c.DomainName,
		"appsyncDomainName": c.AppsyncDomainName,
	}
	if c.CertificateArn != "" {
		result["certificateArn"] = c.CertificateArn
	}
	if c.Description != "" {
		result["description"] = c.Description
	}
	if c.DomainNameArn != "" {
		result["domainNameArn"] = c.DomainNameArn
	}
	if c.HostedZoneId != "" {
		result["hostedZoneId"] = c.HostedZoneId
	}
	if len(c.Tags) > 0 {
		result["tags"] = c.Tags
	}
	return result
}

func apiAssociationToMap(a *appsyncstore.ApiAssociation) map[string]interface{} {
	result := map[string]interface{}{
		"domainName":        a.DomainName,
		"associationStatus": a.AssociationStatus,
	}
	if a.ApiId != "" {
		result["apiId"] = a.ApiId
	}
	if a.DeploymentDetail != "" {
		result["deploymentDetail"] = a.DeploymentDetail
	}
	return result
}
