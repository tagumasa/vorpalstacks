package iam

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// UploadServerCertificate uploads a server certificate entity for the account.
func (s *IAMService) UploadServerCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "ServerCertificateName")
	if name == "" {
		return nil, NewValidationError("ServerCertificateName")
	}
	if err := validateEntityName128(name, "ServerCertificateName"); err != nil {
		return nil, err
	}

	path := request.GetStringParam(req.Parameters, "Path")
	if path == "" {
		path = "/"
	}
	if !validatePath(path) {
		return nil, NewInvalidInputError("Path", "must be a valid path starting and ending with /")
	}
	certificateBody := request.GetStringParam(req.Parameters, "CertificateBody")
	if certificateBody == "" {
		return nil, NewValidationError("CertificateBody")
	}
	if len(certificateBody) > maxCertificateBodyLength {
		return nil, NewInvalidInputError("CertificateBody", fmt.Sprintf("must be 1 to %d characters", maxCertificateBodyLength))
	}
	certificateChain := request.GetStringParam(req.Parameters, "CertificateChain")
	if certificateChain != "" && len(certificateChain) > 2097152 {
		return nil, NewInvalidInputError("CertificateChain", "must be 1 to 2097152 characters")
	}
	privateKey := request.GetStringParam(req.Parameters, "PrivateKey")
	if privateKey == "" {
		return nil, NewValidationError("PrivateKey")
	}
	if len(privateKey) > maxPrivateKeyLength {
		return nil, NewInvalidInputError("PrivateKey", fmt.Sprintf("must be 1 to %d characters", maxPrivateKeyLength))
	}

	parsedCert, err := parseCertificate(certificateBody)
	if err != nil {
		return nil, ErrMalformedCertificate
	}
	privKey, err := parsePrivateKey(privateKey)
	if err != nil {
		return nil, ErrMalformedCertificate
	}
	if !keyPairMatches(parsedCert, privKey) {
		return nil, ErrKeyPairMismatch
	}
	if certificateChain != "" {
		rest := []byte(certificateChain)
		for {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				break
			}
			if _, err := x509.ParseCertificate(block.Bytes); err != nil {
				return nil, ErrMalformedCertificate
			}
		}
	}
	expiration := &parsedCert.NotAfter

	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")
	if err := validateNewTags(newTags); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := store.ServerCertificates().Create(name, path, certificateBody, privateKey, certificateChain, expiration, newTags)
	if err != nil {
		if errors.Is(err, iamstore.ErrServerCertificateAlreadyExists) {
			return nil, NewEntityAlreadyExistsError("Server Certificate " + name)
		}
		return nil, err
	}

	uploadResp := map[string]interface{}{
		"ServerCertificateMetadata": s.serverCertificateMetadataToResponse(cert),
	}
	if certTags := tags.ToResponse(cert.Tags); certTags != nil {
		uploadResp["Tags"] = certTags
	}
	return uploadResp, nil
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
	if certTags := tags.ToResponse(cert.Tags); certTags != nil {
		serverCert["Tags"] = certTags
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
	newName := request.GetStringParam(req.Parameters, "NewServerCertificateName")
	if newPath != "" && !validatePath(newPath) {
		return nil, NewInvalidInputError("NewPath", "must be a valid path starting and ending with /")
	}
	if newName != "" {
		if err := validateEntityName128(newName, "NewServerCertificateName"); err != nil {
			return nil, err
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.ServerCertificates().Exists(name) {
		return nil, NewNoSuchEntityError("server certificate", name)
	}
	if newName != "" && newName != name && store.ServerCertificates().Exists(newName) {
		return nil, NewEntityAlreadyExistsError("Server Certificate " + newName)
	}
	if err := store.ServerCertificates().Update(name, newPath, newName, "", ""); err != nil {
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
	tagsFn: func(r *iamstore.ServerCertificate) *[]tags.Tag { return &r.Tags },
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
