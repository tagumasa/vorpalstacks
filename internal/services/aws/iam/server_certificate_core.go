// Transport-agnostic Core functions for IAM server certificates:
// validation and store operations shared by the AWS-compatible HTTP API
// handlers and any admin surface (the xxxCore pattern).
package iam

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"unicode/utf8"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// UploadServerCertificateInput holds the parameters for uploading a server
// certificate.
type UploadServerCertificateInput struct {
	ServerCertificateName string
	Path                  string
	CertificateBody       string
	PrivateKey            string
	CertificateChain      string
	Tags                  []tags.Tag
}

// UpdateServerCertificateInput holds the parameters for renaming or
// re-pathing a server certificate.
type UpdateServerCertificateInput struct {
	ServerCertificateName    string
	NewPath                  string
	NewServerCertificateName string
}

// uploadServerCertificateCore validates input and uploads a server
// certificate, returning the created record.
func (s *IAMService) uploadServerCertificateCore(store *iamstore.IAMStore, input *UploadServerCertificateInput) (*iamstore.ServerCertificate, error) {
	if input.ServerCertificateName == "" {
		return nil, NewValidationError("ServerCertificateName")
	}
	if err := validateEntityName128(input.ServerCertificateName, "ServerCertificateName"); err != nil {
		return nil, err
	}

	path := input.Path
	if path == "" {
		path = "/"
	}
	if !validatePath(path) {
		return nil, NewInvalidInputError("Path", "must be a valid path starting and ending with /")
	}
	if input.CertificateBody == "" {
		return nil, NewValidationError("CertificateBody")
	}
	// certificateBodyType / certificateChainType / privateKeyType carry
	// Latin-1 patterns, so lengths count Unicode characters.
	if utf8.RuneCountInString(input.CertificateBody) > maxCertificateBodyLength {
		return nil, NewInvalidInputError("CertificateBody", fmt.Sprintf("must be 1 to %d characters", maxCertificateBodyLength))
	}
	if input.CertificateChain != "" && utf8.RuneCountInString(input.CertificateChain) > maxCertificateChainLength {
		return nil, NewInvalidInputError("CertificateChain", fmt.Sprintf("must be 1 to %d characters", maxCertificateChainLength))
	}
	if input.PrivateKey == "" {
		return nil, NewValidationError("PrivateKey")
	}
	if len(input.PrivateKey) > maxPrivateKeyLength {
		return nil, NewInvalidInputError("PrivateKey", fmt.Sprintf("must be 1 to %d characters", maxPrivateKeyLength))
	}

	parsedCert, err := parseCertificate(input.CertificateBody)
	if err != nil {
		return nil, ErrMalformedCertificate
	}
	privKey, err := parsePrivateKey(input.PrivateKey)
	if err != nil {
		return nil, ErrMalformedCertificate
	}
	if !keyPairMatches(parsedCert, privKey) {
		return nil, ErrKeyPairMismatch
	}
	if input.CertificateChain != "" {
		rest := []byte(input.CertificateChain)
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

	if err := validateNewTags(input.Tags); err != nil {
		return nil, err
	}

	cert, err := store.ServerCertificates().Create(input.ServerCertificateName, path, input.CertificateBody, input.PrivateKey, input.CertificateChain, expiration, input.Tags)
	if err != nil {
		if errors.Is(err, iamstore.ErrServerCertificateAlreadyExists) {
			return nil, NewEntityAlreadyExistsError("Server Certificate " + input.ServerCertificateName)
		}
		return nil, err
	}
	return cert, nil
}

// getServerCertificateCore returns the server certificate with the given
// name.
func (s *IAMService) getServerCertificateCore(store *iamstore.IAMStore, name string) (*iamstore.ServerCertificate, error) {
	if name == "" {
		return nil, NewValidationError("ServerCertificateName")
	}
	cert, err := store.ServerCertificates().Get(name)
	if err != nil {
		return nil, NewNoSuchEntityError("server certificate", name)
	}
	return cert, nil
}

// ServerCertificateMaterial implements invokers.IAMServerCertificateProvider.
// It resolves a server certificate by its unique certificate ID (the ASCA…
// identifier cross-service consumers such as CloudFront reference) and
// returns the PEM material so the consumer's listener can terminate TLS.
func (s *IAMService) ServerCertificateMaterial(ctx context.Context, serverCertificateId string) (invokers.TLSCertificateMaterial, error) {
	if serverCertificateId == "" {
		return invokers.TLSCertificateMaterial{}, NewValidationError("ServerCertificateId")
	}
	store, err := s.GetStoreForRegion("")
	if err != nil {
		return invokers.TLSCertificateMaterial{}, err
	}
	cert, err := store.ServerCertificates().GetByID(serverCertificateId)
	if err != nil {
		return invokers.TLSCertificateMaterial{}, NewNoSuchEntityError("server certificate", serverCertificateId)
	}
	if cert.CertificateBody == "" || cert.PrivateKey == "" {
		return invokers.TLSCertificateMaterial{}, fmt.Errorf("server certificate %s does not have serving material available", serverCertificateId)
	}
	return invokers.TLSCertificateMaterial{
		Certificate:      cert.CertificateBody,
		PrivateKey:       cert.PrivateKey,
		CertificateChain: cert.CertificateChain,
	}, nil
}

// updateServerCertificateCore validates input and updates the name or path
// of a server certificate.
func (s *IAMService) updateServerCertificateCore(store *iamstore.IAMStore, input *UpdateServerCertificateInput) error {
	if input.ServerCertificateName == "" {
		return NewValidationError("ServerCertificateName")
	}

	if input.NewPath != "" && !validatePath(input.NewPath) {
		return NewInvalidInputError("NewPath", "must be a valid path starting and ending with /")
	}
	if input.NewServerCertificateName != "" {
		if err := validateEntityName128(input.NewServerCertificateName, "NewServerCertificateName"); err != nil {
			return err
		}
	}

	if !store.ServerCertificates().Exists(input.ServerCertificateName) {
		return NewNoSuchEntityError("server certificate", input.ServerCertificateName)
	}
	if input.NewServerCertificateName != "" && input.NewServerCertificateName != input.ServerCertificateName && store.ServerCertificates().Exists(input.NewServerCertificateName) {
		return NewEntityAlreadyExistsError("Server Certificate " + input.NewServerCertificateName)
	}
	return store.ServerCertificates().Update(input.ServerCertificateName, input.NewPath, input.NewServerCertificateName, "", "")
}

// deleteServerCertificateCore deletes the server certificate with the given
// name.
func (s *IAMService) deleteServerCertificateCore(store *iamstore.IAMStore, name string) error {
	if name == "" {
		return NewValidationError("ServerCertificateName")
	}
	if !store.ServerCertificates().Exists(name) {
		return NewNoSuchEntityError("server certificate", name)
	}
	return store.ServerCertificates().Delete(name)
}

// listServerCertificatesCore returns a paginated list of server
// certificates matching the path prefix.
func (s *IAMService) listServerCertificatesCore(store *iamstore.IAMStore, pathPrefix, marker string, maxItems int) (*iamstore.ServerCertificateListResult, error) {
	return store.ServerCertificates().List(pathPrefix, marker, maxItems)
}
