package iam

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

// UploadServerCertificate uploads a server certificate entity for the account.
func (s *IAMService) UploadServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "ServerCertificateName")
	if name == "" {
		return nil, NewValidationError("ServerCertificateName")
	}

	path := request.GetStringParam(req.Parameters, "Path")
	if path == "" {
		path = "/"
	}
	certificateBody := request.GetStringParam(req.Parameters, "CertificateBody")
	certificateChain := request.GetStringParam(req.Parameters, "CertificateChain")

	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := store.ServerCertificates().Create(name, path, certificateBody, certificateChain, newTags)
	if err != nil {
		if errors.Is(err, iamstore.ErrServerCertificateAlreadyExists) {
			return nil, NewEntityAlreadyExistsError("Server Certificate " + name)
		}
		return nil, err
	}

	return map[string]interface{}{
		"ServerCertificateMetadata": s.serverCertificateMetadataToResponse(cert),
	}, nil
}

// GetServerCertificate retrieves information about the specified server certificate.
func (s *IAMService) GetServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "ServerCertificateName")
	if name == "" {
		return nil, NewValidationError("ServerCertificateName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := store.ServerCertificates().Get(name)
	if err != nil {
		return nil, NewNoSuchEntityError("server certificate", name)
	}

	serverCert := map[string]interface{}{
		"ServerCertificateMetadata": s.serverCertificateMetadataToResponse(cert),
		"CertificateBody":           cert.CertificateBody,
	}
	if cert.CertificateChain != "" {
		serverCert["CertificateChain"] = cert.CertificateChain
	}

	return map[string]interface{}{
		"ServerCertificate": serverCert,
	}, nil
}

// UpdateServerCertificate updates the name or path of a server certificate.
func (s *IAMService) UpdateServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "ServerCertificateName")
	if name == "" {
		return nil, NewValidationError("ServerCertificateName")
	}

	newPath := request.GetStringParam(req.Parameters, "NewPath")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.ServerCertificates().Exists(name) {
		return nil, NewNoSuchEntityError("server certificate", name)
	}
	if err := store.ServerCertificates().Update(name, newPath, "", ""); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteServerCertificate deletes a server certificate.
func (s *IAMService) DeleteServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "ServerCertificateName")
	if name == "" {
		return nil, NewValidationError("ServerCertificateName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.ServerCertificates().Exists(name) {
		return nil, NewNoSuchEntityError("server certificate", name)
	}
	if err := store.ServerCertificates().Delete(name); err != nil {
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
	result, err := store.ServerCertificates().List(pathPrefix, marker, maxItems)
	if err != nil {
		return nil, err
	}

	metadataList := make([]interface{}, len(result.ServerCertificateMetadataList))
	for i, cert := range result.ServerCertificateMetadataList {
		metadataList[i] = s.serverCertificateMetadataToResponse(cert)
	}

	response := map[string]interface{}{
		"ServerCertificateMetadataList": metadataList,
		"IsTruncated":                   result.IsTruncated,
	}

	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

var serverCertificateTagOps = tagOps[*iamstore.ServerCertificate]{
	paramName:  "ServerCertificateName",
	emptyErr:   NewValidationError("ServerCertificateName"),
	notFoundFn: func(n string) error { return NewNoSuchEntityError("server certificate", n) },
	getFn: func(s *iamstore.IAMStore, n string) (*iamstore.ServerCertificate, error) {
		return s.ServerCertificates().Get(n)
	},
	putFn:  func(s *iamstore.IAMStore, r *iamstore.ServerCertificate) error { return s.ServerCertificates().Put(r) },
	tagsFn: func(r *iamstore.ServerCertificate) *[]types.Tag { return &r.Tags },
}

// TagServerCertificate adds tags to a server certificate.
func (s *IAMService) TagServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagResource(ctx, s, reqCtx, req, serverCertificateTagOps)
}

// UntagServerCertificate removes tags from a server certificate.
func (s *IAMService) UntagServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return untagResource(ctx, s, reqCtx, req, serverCertificateTagOps)
}

// ListServerCertificateTags lists the tags attached to a server certificate.
func (s *IAMService) ListServerCertificateTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listResourceTags(ctx, s, reqCtx, req, serverCertificateTagOps)
}

func (s *IAMService) serverCertificateMetadataToResponse(cert *iamstore.ServerCertificate) map[string]interface{} {
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
