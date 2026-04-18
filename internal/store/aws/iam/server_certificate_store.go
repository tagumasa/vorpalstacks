package iam

import (
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const serverCertificateBucketName = "iam_server_certificates"

// ServerCertificateStore provides storage operations for IAM server certificates.
type ServerCertificateStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	kl         common.KeyLocker
}

// NewServerCertificateStore creates a new ServerCertificateStore instance.
func NewServerCertificateStore(store storage.BasicStorage, accountID string) *ServerCertificateStore {
	return &ServerCertificateStore{
		BaseStore:  common.NewBaseStore(store.Bucket(serverCertificateBucketName), "iam"),
		arnBuilder: NewARNBuilder(accountID),
	}
}

// Get retrieves a server certificate by name.
func (s *ServerCertificateStore) Get(name string) (*ServerCertificate, error) {
	var cert ServerCertificate
	if err := s.BaseStore.Get(name, &cert); err != nil {
		return nil, err
	}
	return &cert, nil
}

// Put stores a server certificate, keyed by its name.
func (s *ServerCertificateStore) Put(cert *ServerCertificate) error {
	if cert.Tags == nil {
		cert.Tags = []Tag{}
	}
	return s.BaseStore.Put(cert.ServerCertificateName, cert)
}

// Delete removes a server certificate by name.
func (s *ServerCertificateStore) Delete(name string) error {
	return s.BaseStore.Delete(name)
}

// Exists reports whether a server certificate exists with the given name.
func (s *ServerCertificateStore) Exists(name string) bool {
	return s.BaseStore.Exists(name)
}

// Create creates a new server certificate with the given name, path, certificate body, and chain.
func (s *ServerCertificateStore) Create(name, path, certificateBody, certificateChain string, tags []Tag) (*ServerCertificate, error) {
	if s.Exists(name) {
		return nil, ErrServerCertificateAlreadyExists
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
	var filter common.FilterFunc[ServerCertificate]
	if pathPrefix != "" {
		filter = func(c *ServerCertificate) bool { return strings.HasPrefix(c.Path, pathPrefix) }
	}
	result, err := common.List[ServerCertificate](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, filter)
	if err != nil {
		return nil, err
	}
	return &ServerCertificateListResult{
		ServerCertificateMetadataList: result.Items,
		IsTruncated:                   result.IsTruncated,
		Marker:                        result.NextMarker,
	}, nil
}

// Count returns the total number of server certificates.
func (s *ServerCertificateStore) Count() int {
	return s.BaseStore.Count()
}
