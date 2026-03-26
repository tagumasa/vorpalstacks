package iam

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
)

const serverCertificateBucketName = "iam_server_certificates"

// ServerCertificateStore provides storage operations for IAM server certificates.
type ServerCertificateStore struct {
	bucket     storage.Bucket
	arnBuilder *ARNBuilder
}

// NewServerCertificateStore creates a new ServerCertificateStore instance.
func NewServerCertificateStore(store storage.BasicStorage, accountID string) *ServerCertificateStore {
	return &ServerCertificateStore{
		bucket:     store.Bucket(serverCertificateBucketName),
		arnBuilder: NewARNBuilder(accountID),
	}
}

// Get retrieves a server certificate by name.
func (s *ServerCertificateStore) Get(name string) (*ServerCertificate, error) {
	data, err := s.bucket.Get([]byte(name))
	if err != nil {
		return nil, NewStoreError("get_server_certificate", err)
	}
	if data == nil {
		return nil, NewStoreError("get_server_certificate", ErrServerCertificateNotFound)
	}
	var cert ServerCertificate
	if err := json.Unmarshal(data, &cert); err != nil {
		return nil, NewStoreError("get_server_certificate", err)
	}
	return &cert, nil
}

// Put stores a server certificate, keyed by its name.
func (s *ServerCertificateStore) Put(cert *ServerCertificate) error {
	if cert.Tags == nil {
		cert.Tags = []Tag{}
	}
	data, err := json.Marshal(cert)
	if err != nil {
		return NewStoreError("put_server_certificate", err)
	}
	if err := s.bucket.Put([]byte(cert.ServerCertificateName), data); err != nil {
		return NewStoreError("put_server_certificate", err)
	}
	return nil
}

// Delete removes a server certificate by name.
func (s *ServerCertificateStore) Delete(name string) error {
	if err := s.bucket.Delete([]byte(name)); err != nil {
		return NewStoreError("delete_server_certificate", err)
	}
	return nil
}

// Exists reports whether a server certificate exists with the given name.
func (s *ServerCertificateStore) Exists(name string) bool {
	return s.bucket.Has([]byte(name))
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
}

// List returns server certificates matching the given path prefix, with pagination support.
func (s *ServerCertificateStore) List(pathPrefix, marker string, maxItems int) (*ServerCertificateListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}
	var certs []*ServerCertificate
	count := 0
	started := marker == ""
	hasMore := false

	err := s.bucket.ForEach(func(k, v []byte) error {
		var cert ServerCertificate
		if err := json.Unmarshal(v, &cert); err != nil {
			return err
		}
		if pathPrefix != "" && !strings.HasPrefix(cert.Path, pathPrefix) {
			return nil
		}
		if !started {
			if cert.ServerCertificateName >= marker {
				started = true
			}
			return nil
		}
		if count < maxItems {
			certs = append(certs, &cert)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_server_certificates", err)
	}

	result := &ServerCertificateListResult{
		ServerCertificateMetadataList: certs,
		IsTruncated:                   hasMore,
	}
	if len(certs) > 0 {
		result.Marker = certs[len(certs)-1].ServerCertificateName
	}
	return result, nil
}

// Count returns the total number of server certificates.
func (s *ServerCertificateStore) Count() int {
	return s.bucket.Count()
}
