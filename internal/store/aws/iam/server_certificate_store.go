package iam

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/utils/aws/types"
)

const serverCertificateBucketName = "iam_server_certificates"

// ServerCertificateStore provides storage operations for IAM server certificates.
type ServerCertificateStore struct {
	entityStore[ServerCertificate]
	arnBuilder *ARNBuilder
}

// NewServerCertificateStore creates a new ServerCertificateStore instance.
func NewServerCertificateStore(store storage.BasicStorage, accountID string) *ServerCertificateStore {
	return &ServerCertificateStore{
		entityStore: newEntityStore[ServerCertificate](store, serverCertificateBucketName),
		arnBuilder:  NewARNBuilder(accountID),
	}
}

// Put stores a server certificate, keyed by its name.
func (s *ServerCertificateStore) Put(cert *ServerCertificate) error {
	if cert.Tags == nil {
		cert.Tags = []types.Tag{}
	}
	return s.BaseStore.Put(cert.ServerCertificateName, cert)
}

// Create creates a new server certificate with the given name, path, certificate body, and chain.
func (s *ServerCertificateStore) Create(name, path, certificateBody, certificateChain string, tags []types.Tag) (*ServerCertificate, error) {
	if s.Exists(name) {
		return nil, NewStoreError("create_server_certificate", ErrServerCertificateAlreadyExists)
	}
	id, err := GenerateServerCertificateID()
	if err != nil {
		return nil, err
	}
	cert := &ServerCertificate{
		ID:                    id,
		Path:                  path,
		ServerCertificateName: name,
		Arn:                   s.arnBuilder.ServerCertificateARN(name),
		AccountId:             s.arnBuilder.AccountID(),
		CreateDate:            time.Now().UTC(),
		CertificateBody:       certificateBody,
		CertificateChain:      certificateChain,
		Tags:                  tags,
	}
	if err := s.Put(cert); err != nil {
		return nil, err
	}
	return cert, nil
}

// Update modifies the path, certificate body, or chain of an existing server certificate.
func (s *ServerCertificateStore) Update(name, newPath, newCertificateBody, newCertificateChain string) error {
	return s.kl.WithLock(name, func() error {
		cert, err := s.Get(name)
		if err != nil {
			return err
		}
		if newPath != "" {
			cert.Path = newPath
		}
		if newCertificateBody != "" {
			cert.CertificateBody = newCertificateBody
		}
		if newCertificateChain != "" {
			cert.CertificateChain = newCertificateChain
		}
		return s.Put(cert)
	})
}

// List returns server certificates matching the given path prefix, with pagination support.
func (s *ServerCertificateStore) List(pathPrefix, marker string, maxItems int) (*ServerCertificateListResult, error) {
	items, truncated, nextMarker, err := listEntitiesWithPathPrefix(s.BaseStore, pathPrefix, marker, maxItems, func(c *ServerCertificate) string { return c.Path })
	if err != nil {
		return nil, err
	}
	return &ServerCertificateListResult{
		ServerCertificateMetadataList: items,
		IsTruncated:                   truncated,
		Marker:                        nextMarker,
	}, nil
}
