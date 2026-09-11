package iam

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// UploadServerCertificate uploads a server certificate entity for the account.
func (s *IAMService) UploadServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UploadServerCertificateInput{
		ServerCertificateName: request.GetStringParam(req.Parameters, "ServerCertificateName"),
		Path:                  request.GetStringParam(req.Parameters, "Path"),
		CertificateBody:       request.GetStringParam(req.Parameters, "CertificateBody"),
		PrivateKey:            request.GetStringParam(req.Parameters, "PrivateKey"),
		CertificateChain:      request.GetStringParam(req.Parameters, "CertificateChain"),
		Tags:                  tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	cert, err := s.uploadServerCertificateCore(store, input)
	if err != nil {
		return nil, err
	}

	uploadResp := map[string]interface{}{
		"ServerCertificateMetadata": serverCertificateMetadataToResponse(cert),
	}
	if certTags := tags.ToResponse(cert.Tags); certTags != nil {
		uploadResp["Tags"] = certTags
	}
	return uploadResp, nil
}

// GetServerCertificate retrieves information about the specified server certificate.
func (s *IAMService) GetServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := s.getServerCertificateCore(store, request.GetStringParam(req.Parameters, "ServerCertificateName"))
	if err != nil {
		return nil, err
	}

	serverCert := map[string]interface{}{
		"ServerCertificateMetadata": serverCertificateMetadataToResponse(cert),
		"CertificateBody":           cert.CertificateBody,
	}
	if cert.CertificateChain != "" {
		serverCert["CertificateChain"] = cert.CertificateChain
	}
	if certTags := tags.ToResponse(cert.Tags); certTags != nil {
		serverCert["Tags"] = certTags
	}

	return map[string]interface{}{
		"ServerCertificate": serverCert,
	}, nil
}

// UpdateServerCertificate updates the name or path of a server certificate.
func (s *IAMService) UpdateServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UpdateServerCertificateInput{
		ServerCertificateName:    request.GetStringParam(req.Parameters, "ServerCertificateName"),
		NewPath:                  request.GetStringParam(req.Parameters, "NewPath"),
		NewServerCertificateName: request.GetStringParam(req.Parameters, "NewServerCertificateName"),
	}
	if err := s.updateServerCertificateCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteServerCertificate deletes a server certificate.
func (s *IAMService) DeleteServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteServerCertificateCore(store, request.GetStringParam(req.Parameters, "ServerCertificateName")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListServerCertificates lists the server certificates in the account.
func (s *IAMService) ListServerCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listServerCertificatesCore(store, pathPrefix, marker, maxItems)
	if err != nil {
		return nil, err
	}

	metadataList := make([]interface{}, len(result.ServerCertificateMetadataList))
	for i, cert := range result.ServerCertificateMetadataList {
		metadataList[i] = serverCertificateMetadataToResponse(cert)
	}

	resp := map[string]interface{}{
		"ServerCertificateMetadataList": metadataList,
		"IsTruncated":                   result.IsTruncated,
	}

	if result.Marker != "" {
		resp["Marker"] = result.Marker
	}

	return resp, nil
}

// TagServerCertificate adds tags to a server certificate.
func (s *IAMService) TagServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &TagResourceInput{
		ResourceName: request.GetStringParam(req.Parameters, "ServerCertificateName"),
		Tags:         tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	if err := tagResourceCore(store, serverCertificateTagOps, input); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// UntagServerCertificate removes tags from a server certificate.
func (s *IAMService) UntagServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UntagResourceInput{
		ResourceName: request.GetStringParam(req.Parameters, "ServerCertificateName"),
		TagKeys:      tags.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys"),
	}
	if err := untagResourceCore(store, serverCertificateTagOps, input); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListServerCertificateTags lists the tags attached to a server certificate.
func (s *IAMService) ListServerCertificateTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &ListResourceTagsInput{
		ResourceName: request.GetStringParam(req.Parameters, "ServerCertificateName"),
		Marker:       request.GetStringParam(req.Parameters, "Marker"),
		MaxItems:     pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems),
	}
	result, err := listResourceTagsCore(store, serverCertificateTagOps, input)
	if err != nil {
		return nil, err
	}
	resp := map[string]interface{}{
		"Tags":        tags.ToResponse(result.Tags),
		"IsTruncated": result.IsTruncated,
	}
	if result.Marker != "" {
		resp["Marker"] = result.Marker
	}
	return resp, nil
}

func serverCertificateMetadataToResponse(cert *iamstore.ServerCertificate) map[string]interface{} {
	resp := map[string]interface{}{
		"ServerCertificateId":   cert.ID,
		"ServerCertificateName": cert.ServerCertificateName,
		"Arn":                   cert.Arn,
		"Path":                  cert.Path,
		"UploadDate":            cert.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}

	if cert.Expiration != nil {
		resp["Expiration"] = cert.Expiration.Format(timeutils.ISO8601SimpleFormat)
	}

	return resp
}
