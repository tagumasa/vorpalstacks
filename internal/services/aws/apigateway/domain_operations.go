// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	store "vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateDomainName creates a new domain name for API Gateway.
func (s *APIGatewayService) CreateDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	domain := &store.DomainName{
		DomainName:                          domainName,
		CertificateArn:                      request.GetStringParam(req.Parameters, "certificateArn"),
		CertificateName:                     request.GetStringParam(req.Parameters, "certificateName"),
		RegionalCertificateArn:              request.GetStringParam(req.Parameters, "regionalCertificateArn"),
		RegionalCertificateName:             request.GetStringParam(req.Parameters, "regionalCertificateName"),
		SecurityPolicy:                      request.GetStringParam(req.Parameters, "securityPolicy"),
		OwnershipVerificationCertificateArn: request.GetStringParam(req.Parameters, "ownershipVerificationCertificateArn"),
		EndpointAccessMode:                  request.GetStringParam(req.Parameters, "endpointAccessMode"),
		Policy:                              request.GetStringParam(req.Parameters, "policy"),
		RoutingMode:                         request.GetStringParam(req.Parameters, "routingMode"),
	}

	if mutualTls, ok := req.Parameters["mutualTlsAuthentication"].(map[string]interface{}); ok {
		domain.MutualTlsAuthentication = &store.MutualTlsAuthentication{}
		if v, ok := mutualTls["truststoreUri"].(string); ok {
			domain.MutualTlsAuthentication.TruststoreUri = v
		}
		if v, ok := mutualTls["truststoreVersion"].(string); ok {
			domain.MutualTlsAuthentication.TruststoreVersion = v
		}
	}

	if endpointConfig, ok := req.Parameters["endpointConfiguration"].(map[string]interface{}); ok {
		domain.EndpointConfiguration = &store.EndpointConfiguration{}
		if types, ok := endpointConfig["types"].([]interface{}); ok {
			for _, t := range types {
				if ts, ok := t.(string); ok {
					domain.EndpointConfiguration.Types = append(domain.EndpointConfiguration.Types, ts)
				}
			}
		}
	}

	if tags, ok := req.Parameters["tags"].(map[string]interface{}); ok {
		domain.Tags = tagutil.MapInterfaceToTags(tags)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.domains.CreateDomainName(domain)
	if err != nil {
		return nil, err
	}

	return s.toDomainNameResponse(created), nil
}

// GetDomainName retrieves a domain name by its name.
func (s *APIGatewayService) GetDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		domainName = getPathParam(req, "domainName")
	}
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domain, err := stores.domains.GetDomainName(domainName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toDomainNameResponse(domain), nil
}

// DeleteDomainName deletes a domain name.
func (s *APIGatewayService) DeleteDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		domainName = getPathParam(req, "domainName")
	}
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.domains.DeleteDomainName(domainName); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// UpdateDomainName updates an existing domain name.
func (s *APIGatewayService) UpdateDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		domainName = getPathParam(req, "domainName")
	}
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domain, err := stores.domains.GetDomainName(domainName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch po.Path {
		case "/certificateArn":
			domain.CertificateArn = po.Value
		case "/regionalCertificateArn":
			domain.RegionalCertificateArn = po.Value
		case "/certificateName":
			domain.CertificateName = po.Value
		}
	}

	if err := stores.domains.UpdateDomainName(domain); err != nil {
		return nil, err
	}

	return s.toDomainNameResponse(domain), nil
}

// GetDomainNames lists all domain names.
func (s *APIGatewayService) GetDomainNames(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxItems := request.GetIntParam(req.Parameters, "limit")
	if maxItems <= 0 {
		maxItems = 25
	}
	marker := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.domains.ListDomainNames(common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, d := range result.Items {
		items = append(items, s.toDomainNameResponse(d))
	}

	response := map[string]interface{}{
		"item": items,
	}
	if result.IsTruncated {
		response["position"] = result.NextMarker
	}

	return response, nil
}

func (s *APIGatewayService) toDomainNameResponse(d *store.DomainName) map[string]interface{} {
	response := map[string]interface{}{
		"domainName": d.DomainName,
	}

	if d.DomainNameId != "" {
		response["domainNameId"] = d.DomainNameId
	}
	if d.DomainNameArn != "" {
		response["domainNameArn"] = d.DomainNameArn
	}
	if d.CertificateArn != "" {
		response["certificateArn"] = d.CertificateArn
	}
	if d.CertificateName != "" {
		response["certificateName"] = d.CertificateName
	}
	if !d.CertificateUploadDate.IsZero() {
		response["certificateUploadDate"] = timeutils.FormatEpochSeconds(d.CertificateUploadDate)
	}
	if d.DistributionDomainName != "" {
		response["distributionDomainName"] = d.DistributionDomainName
	}
	if d.DistributionHostedZoneId != "" {
		response["distributionHostedZoneId"] = d.DistributionHostedZoneId
	}
	if d.RegionalDomainName != "" {
		response["regionalDomainName"] = d.RegionalDomainName
	}
	if d.RegionalHostedZoneId != "" {
		response["regionalHostedZoneId"] = d.RegionalHostedZoneId
	}
	if d.RegionalCertificateArn != "" {
		response["regionalCertificateArn"] = d.RegionalCertificateArn
	}
	if d.RegionalCertificateName != "" {
		response["regionalCertificateName"] = d.RegionalCertificateName
	}
	if d.DomainNameStatus != "" {
		response["domainNameStatus"] = d.DomainNameStatus
	}
	if d.DomainNameStatusMessage != "" {
		response["domainNameStatusMessage"] = d.DomainNameStatusMessage
	}
	if d.SecurityPolicy != "" {
		response["securityPolicy"] = d.SecurityPolicy
	}
	if d.EndpointAccessMode != "" {
		response["endpointAccessMode"] = d.EndpointAccessMode
	}
	if d.MutualTlsAuthentication != nil {
		mtls := map[string]interface{}{}
		if d.MutualTlsAuthentication.TruststoreUri != "" {
			mtls["truststoreUri"] = d.MutualTlsAuthentication.TruststoreUri
		}
		if d.MutualTlsAuthentication.TruststoreVersion != "" {
			mtls["truststoreVersion"] = d.MutualTlsAuthentication.TruststoreVersion
		}
		if len(d.MutualTlsAuthentication.TruststoreWarnings) > 0 {
			mtls["truststoreWarnings"] = d.MutualTlsAuthentication.TruststoreWarnings
		}
		response["mutualTlsAuthentication"] = mtls
	}
	if d.OwnershipVerificationCertificateArn != "" {
		response["ownershipVerificationCertificateArn"] = d.OwnershipVerificationCertificateArn
	}
	if d.Policy != "" {
		response["policy"] = d.Policy
	}
	if d.RoutingMode != "" {
		response["routingMode"] = d.RoutingMode
	}
	if d.EndpointConfiguration != nil {
		response["endpointConfiguration"] = map[string]interface{}{
			"types": d.EndpointConfiguration.Types,
		}
	}
	if len(d.Tags) > 0 {
		response["tags"] = tagutil.ToMap(d.Tags)
	}

	return response
}

// CreateBasePathMapping creates a new base path mapping.
func (s *APIGatewayService) CreateBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		domainName = getPathParam(req, "domainName")
	}
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	restApiId := request.GetStringParam(req.Parameters, "restApiId")
	if restApiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	mapping := &store.BasePathMapping{
		BasePath:  request.GetStringParam(req.Parameters, "basePath"),
		RestApiId: restApiId,
		Stage:     request.GetStringParam(req.Parameters, "stage"),
	}

	if mapping.BasePath == "" {
		mapping.BasePath = "(none)"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.domains.CreateBasePathMapping(domainName, mapping)
	if err != nil {
		return nil, err
	}

	return s.toBasePathMappingResponse(created), nil
}

// GetBasePathMapping retrieves a base path mapping.
func (s *APIGatewayService) GetBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		domainName = getPathParam(req, "domainName")
	}
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	basePath := request.GetStringParam(req.Parameters, "basePath")
	if basePath == "" {
		basePath = getPathParam(req, "basePath")
	}
	if basePath == "" {
		return nil, NewBadRequestException("basePath is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	mapping, err := stores.domains.GetBasePathMapping(domainName, basePath)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toBasePathMappingResponse(mapping), nil
}

// DeleteBasePathMapping deletes a base path mapping.
func (s *APIGatewayService) DeleteBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		domainName = getPathParam(req, "domainName")
	}
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	basePath := request.GetStringParam(req.Parameters, "basePath")
	if basePath == "" {
		basePath = getPathParam(req, "basePath")
	}
	if basePath == "" {
		return nil, NewBadRequestException("basePath is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.domains.DeleteBasePathMapping(domainName, basePath); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// UpdateBasePathMapping updates an existing base path mapping.
func (s *APIGatewayService) UpdateBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		domainName = getPathParam(req, "domainName")
	}
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	basePath := request.GetStringParam(req.Parameters, "basePath")
	if basePath == "" {
		basePath = getPathParam(req, "basePath")
	}
	if basePath == "" {
		return nil, NewBadRequestException("basePath is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stores.keyLocker.Lock(domainName + ":" + basePath)
	defer stores.keyLocker.Unlock(domainName + ":" + basePath)

	mapping, err := stores.domains.GetBasePathMapping(domainName, basePath)
	if err != nil {
		return nil, ErrNotFoundException
	}

	renamed := false
	for _, po := range parsePatchOperations(req.Parameters) {
		switch po.Path {
		case "/restApiId":
			mapping.RestApiId = po.Value
		case "/stage":
			mapping.Stage = po.Value
		case "/basePath":
			oldPath := basePath
			basePath = po.Value
			mapping.BasePath = po.Value
			renamed = true
			if err := stores.domains.DeleteBasePathMapping(domainName, oldPath); err != nil {
				return nil, err
			}
		}
	}

	if renamed {
		if _, err := stores.domains.CreateBasePathMapping(domainName, mapping); err != nil {
			return nil, err
		}
	} else {
		if err := stores.domains.UpdateBasePathMapping(domainName, basePath, mapping); err != nil {
			return nil, err
		}
	}

	return s.toBasePathMappingResponse(mapping), nil
}

// GetBasePathMappings lists all base path mappings for a domain name.
func (s *APIGatewayService) GetBasePathMappings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		domainName = getPathParam(req, "domainName")
	}
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}

	maxItems := request.GetIntParam(req.Parameters, "limit")
	if maxItems <= 0 {
		maxItems = 25
	}
	marker := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.domains.ListBasePathMappings(domainName, common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, m := range result.Items {
		items = append(items, s.toBasePathMappingResponse(m))
	}

	response := map[string]interface{}{
		"item": items,
	}
	if result.IsTruncated {
		response["position"] = result.NextMarker
	}

	return response, nil
}

func (s *APIGatewayService) toBasePathMappingResponse(m *store.BasePathMapping) map[string]interface{} {
	return map[string]interface{}{
		"basePath":  m.BasePath,
		"restApiId": m.RestApiId,
		"stage":     m.Stage,
	}
}
