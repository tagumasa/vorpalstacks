// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/utils/timeutils"
)

// domainNameRequestParams extracts the domain identity members from the wire
// request: the domainName parameter (query, then path label) and the
// domainNameId parameter. At least one must be present; the Core enforces it.
func domainNameRequestParams(req *request.ParsedRequest) (string, string) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		domainName = getPathParam(req, "domainName")
	}
	return domainName, request.GetStringParam(req.Parameters, "domainNameId")
}

// basePathRequestParam extracts the basePath member from the wire request
// (query parameter, then path label).
func basePathRequestParam(req *request.ParsedRequest) string {
	basePath := request.GetStringParam(req.Parameters, "basePath")
	if basePath == "" {
		basePath = getPathParam(req, "basePath")
	}
	return basePath
}

// CreateDomainName creates a new domain name for API Gateway.
func (s *APIGatewayService) CreateDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := &DomainNameCreateInput{
		DomainName:                          request.GetStringParam(req.Parameters, "domainName"),
		CertificateArn:                      request.GetStringParam(req.Parameters, "certificateArn"),
		CertificateName:                     request.GetStringParam(req.Parameters, "certificateName"),
		CertificateBody:                     request.GetStringParam(req.Parameters, "certificateBody"),
		CertificatePrivateKey:               request.GetStringParam(req.Parameters, "certificatePrivateKey"),
		CertificateChain:                    request.GetStringParam(req.Parameters, "certificateChain"),
		RegionalCertificateArn:              request.GetStringParam(req.Parameters, "regionalCertificateArn"),
		RegionalCertificateName:             request.GetStringParam(req.Parameters, "regionalCertificateName"),
		SecurityPolicy:                      request.GetStringParam(req.Parameters, "securityPolicy"),
		OwnershipVerificationCertificateArn: request.GetStringParam(req.Parameters, "ownershipVerificationCertificateArn"),
		EndpointAccessMode:                  request.GetStringParam(req.Parameters, "endpointAccessMode"),
		Policy:                              request.GetStringParam(req.Parameters, "policy"),
		RoutingMode:                         request.GetStringParam(req.Parameters, "routingMode"),
	}
	if mutualTls, ok := req.Parameters["mutualTlsAuthentication"].(map[string]interface{}); ok {
		in.HasMutualTlsAuthentication = true
		if v, ok := mutualTls["truststoreUri"].(string); ok {
			in.MutualTlsTruststoreUri = v
		}
		if v, ok := mutualTls["truststoreVersion"].(string); ok {
			in.MutualTlsTruststoreVersion = v
		}
	}
	if endpointConfig, ok := req.Parameters["endpointConfiguration"].(map[string]interface{}); ok {
		in.HasEndpointConfiguration = true
		if types, ok := endpointConfig["types"].([]interface{}); ok {
			for _, t := range types {
				if ts, ok := t.(string); ok {
					in.EndpointTypes = append(in.EndpointTypes, ts)
				}
			}
		}
	}
	if tags, ok := req.Parameters["tags"].(map[string]interface{}); ok {
		in.Tags = tagutil.MapInterfaceToTags(tags)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createDomainNameCore(ctx, stores, reqCtx.GetRegion(), in)
	if err != nil {
		return nil, err
	}
	return s.toDomainNameResponse(created), nil
}

// GetDomainName retrieves a domain name by its name.
func (s *APIGatewayService) GetDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, domainNameId := domainNameRequestParams(req)
	domain, err := s.getDomainNameCore(stores, domainName, domainNameId)
	if err != nil {
		return nil, err
	}
	return s.toDomainNameResponse(domain), nil
}

// DeleteDomainName deletes a domain name.
func (s *APIGatewayService) DeleteDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, domainNameId := domainNameRequestParams(req)
	if err := s.deleteDomainNameCore(ctx, stores, reqCtx.GetRegion(), domainName, domainNameId); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// UpdateDomainName updates an existing domain name.
func (s *APIGatewayService) UpdateDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	domainName, domainNameId := domainNameRequestParams(req)
	domain, err := s.updateDomainNameCore(ctx, stores, reqCtx.GetRegion(), domainName, domainNameId, ops)
	if err != nil {
		return nil, err
	}
	return s.toDomainNameResponse(domain), nil
}

// GetDomainNames lists all domain names.
func (s *APIGatewayService) GetDomainNames(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxItems, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	marker := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listDomainNamesCore(stores, marker, maxItems)
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

func (s *APIGatewayService) toDomainNameResponse(d *apigateway.DomainName) map[string]interface{} {
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
	if d.ManagementPolicy != "" {
		response["managementPolicy"] = d.ManagementPolicy
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
	domainName, domainNameId := domainNameRequestParams(req)
	in := &BasePathMappingInput{
		BasePath:  request.GetStringParam(req.Parameters, "basePath"),
		RestApiId: request.GetStringParam(req.Parameters, "restApiId"),
		Stage:     request.GetStringParam(req.Parameters, "stage"),
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createBasePathMappingCore(stores, domainName, domainNameId, in)
	if err != nil {
		return nil, err
	}
	return s.toBasePathMappingResponse(created), nil
}

// GetBasePathMapping retrieves a base path mapping.
func (s *APIGatewayService) GetBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, domainNameId := domainNameRequestParams(req)
	mapping, err := s.getBasePathMappingCore(stores, domainName, domainNameId, basePathRequestParam(req))
	if err != nil {
		return nil, err
	}
	return s.toBasePathMappingResponse(mapping), nil
}

// DeleteBasePathMapping deletes a base path mapping.
func (s *APIGatewayService) DeleteBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, domainNameId := domainNameRequestParams(req)
	if err := s.deleteBasePathMappingCore(stores, domainName, domainNameId, basePathRequestParam(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// UpdateBasePathMapping updates an existing base path mapping.
func (s *APIGatewayService) UpdateBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	domainName, domainNameId := domainNameRequestParams(req)
	mapping, err := s.updateBasePathMappingCore(stores, domainName, domainNameId, basePathRequestParam(req), ops)
	if err != nil {
		return nil, err
	}
	return s.toBasePathMappingResponse(mapping), nil
}

// GetBasePathMappings lists all base path mappings for a domain name.
func (s *APIGatewayService) GetBasePathMappings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, domainNameId := domainNameRequestParams(req)
	domainName, err = s.resolveDomainNameCore(stores, domainName, domainNameId)
	if err != nil {
		return nil, err
	}

	maxItems, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	marker := request.GetStringParam(req.Parameters, "position")

	result, err := s.listBasePathMappingsCore(stores, domainName, marker, maxItems)
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

func (s *APIGatewayService) toBasePathMappingResponse(m *apigateway.BasePathMapping) map[string]interface{} {
	return map[string]interface{}{
		"basePath":  m.BasePath,
		"restApiId": m.RestApiId,
		"stage":     m.Stage,
	}
}
