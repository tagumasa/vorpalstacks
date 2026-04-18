// Package acm provides ACM (AWS Certificate Manager) storage functionality for vorpalstacks.
package acm

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// CertificateStore provides storage operations for ACM certificates.
type CertificateStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	mu         sync.RWMutex
}

func certificateBucketName(region string) string {
	return "acm_certificates-" + region
}

// NewCertificateStore creates a new CertificateStore instance with the specified storage, account ID, and region.
func NewCertificateStore(store storage.BasicStorage, accountId, region string) *CertificateStore {
	return &CertificateStore{
		BaseStore:  common.NewBaseStore(store.Bucket(certificateBucketName(region)), "acm"),
		arnBuilder: NewARNBuilder(accountId, region),
	}
}

// Get retrieves an ACM certificate by its ARN from the store.
// Returns the certificate or an error if not found.
func (s *CertificateStore) Get(arn string) (*Certificate, error) {
	certId := s.extractCertificateId(arn)
	var cert Certificate
	if err := s.BaseStore.Get(certId, &cert); err != nil {
		return nil, NewStoreError("get_certificate", err)
	}
	return &cert, nil
}

// GetByDomain retrieves an ACM certificate by its domain name from the store.
// Returns the certificate or an error if not found.
func (s *CertificateStore) GetByDomain(domain string) (*Certificate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return common.FindFirst[Certificate](s.BaseStore, func(c *Certificate) bool { return strings.EqualFold(c.DomainName, domain) })
}

// List returns a list of ACM certificates from the store with pagination support.
func (s *CertificateStore) List(marker string, maxItems int) (*CertificateListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var certs []*CertificateSummary
	count := 0
	started := marker == ""
	hasMore := false
	var lastIncludedArn string

	err := s.ForEach(func(key string, value []byte) error {
		var cert Certificate
		if err := json.Unmarshal(value, &cert); err != nil {
			return err
		}

		if !started {
			if cert.CertificateArn > marker {
				started = true
			}
			if !started {
				return nil
			}
		}

		if count < maxItems {
			certs = append(certs, certificateToSummary(&cert))
			count++
			lastIncludedArn = cert.CertificateArn
		} else {
			hasMore = true
			return nil
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_certificates", err)
	}

	nextToken := ""
	if hasMore && lastIncludedArn != "" {
		nextToken = lastIncludedArn
	}
	return &CertificateListResult{
		Certificates: certs,
		IsTruncated:  hasMore,
		NextToken:    nextToken,
	}, nil
}

// ListAll returns all ACM certificates from the store.
func (s *CertificateStore) ListAll() ([]*Certificate, error) {
	return common.ListAll[Certificate](s.BaseStore)
}

// Create creates a new ACM certificate in the store.
// Returns an error if the certificate already exists.
func (s *CertificateStore) Create(cert *Certificate) error {
	certId := s.extractCertificateId(cert.CertificateArn)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists(certId) {
		return NewStoreError("create_certificate", ErrCertificateExists)
	}
	cert.CreatedAt = time.Now()
	if err := s.BaseStore.Put(certId, cert); err != nil {
		return NewStoreError("create_certificate", err)
	}
	return nil
}

// Update updates an existing ACM certificate in the store.
// Returns an error if the certificate does not exist.
func (s *CertificateStore) Update(cert *Certificate) error {
	certId := s.extractCertificateId(cert.CertificateArn)

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(certId) {
		return NewStoreError("update_certificate", ErrCertificateNotFound)
	}
	if err := s.BaseStore.Put(certId, cert); err != nil {
		return NewStoreError("update_certificate", err)
	}
	return nil
}

// Delete deletes an ACM certificate by its ARN from the store.
// Returns an error if the certificate does not exist.
func (s *CertificateStore) Delete(arn string) error {
	certId := s.extractCertificateId(arn)

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(certId) {
		return NewStoreError("delete_certificate", ErrCertificateNotFound)
	}
	if err := s.BaseStore.Delete(certId); err != nil {
		return NewStoreError("delete_certificate", err)
	}
	return nil
}

// Exists checks whether an ACM certificate exists in the store by its ARN.
func (s *CertificateStore) Exists(arn string) bool {
	certId := s.extractCertificateId(arn)
	return s.BaseStore.Exists(certId)
}

func (s *CertificateStore) extractCertificateId(arn string) string {
	parts := strings.Split(arn, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return arn
}

func formatEpochSeconds(t time.Time) float64 {
	return float64(t.Unix()) + float64(t.Nanosecond())/1e9
}

// CertificateToSummary converts a Certificate to a CertificateSummary.
func CertificateToSummary(cert *Certificate) *CertificateSummary {
	summary := &CertificateSummary{
		CertificateArn:                       cert.CertificateArn,
		DomainName:                           cert.DomainName,
		SubjectAlternativeNameSummaries:      cert.SubjectAlternativeNames,
		HasAdditionalSubjectAlternativeNames: len(cert.SubjectAlternativeNames) > 0,
		Status:                               cert.Status,
		Type:                                 cert.Type,
		RenewalEligibility:                   cert.RenewalEligibility,
		KeyAlgorithm:                         cert.KeyAlgorithm,
		InUse:                                len(cert.InUseBy) > 0,
		Exported:                             cert.PrivateKey != "",
	}

	if !cert.NotBefore.IsZero() {
		summary.NotBefore = formatEpochSeconds(cert.NotBefore)
	}
	if !cert.NotAfter.IsZero() {
		summary.NotAfter = formatEpochSeconds(cert.NotAfter)
	}
	if !cert.CreatedAt.IsZero() {
		summary.CreatedAt = formatEpochSeconds(cert.CreatedAt)
	}
	if !cert.IssuedAt.IsZero() {
		summary.IssuedAt = formatEpochSeconds(cert.IssuedAt)
	}
	if !cert.ImportedAt.IsZero() {
		summary.ImportedAt = formatEpochSeconds(cert.ImportedAt)
	}

	if cert.Options != nil {
		summary.ExportOption = cert.Options.Export
	}

	return summary
}

func certificateToSummary(cert *Certificate) *CertificateSummary {
	return CertificateToSummary(cert)
}

// GetTags retrieves the tags associated with an ACM certificate.
func (s *CertificateStore) GetTags(arn string) ([]common.Tag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cert, err := s.Get(arn)
	if err != nil {
		return nil, err
	}
	return cert.Tags, nil
}

// AddTags adds tags to an ACM certificate.
func (s *CertificateStore) AddTags(arn string, tags []common.Tag) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cert, err := s.Get(arn)
	if err != nil {
		return err
	}
	cert.Tags = append(cert.Tags, tags...)
	return s.BaseStore.Put(s.extractCertificateId(arn), cert)
}

// RemoveTags removes tags from an ACM certificate.
func (s *CertificateStore) RemoveTags(arn string, tagKeys []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cert, err := s.Get(arn)
	if err != nil {
		return err
	}
	keySet := make(map[string]bool)
	for _, k := range tagKeys {
		keySet[k] = true
	}
	var remaining []common.Tag
	for _, t := range cert.Tags {
		if !keySet[t.Key] {
			remaining = append(remaining, t)
		}
	}
	cert.Tags = remaining
	return s.BaseStore.Put(s.extractCertificateId(arn), cert)
}

var accountConfigurations = make(map[string]*AccountConfiguration)
var accountConfigMutex sync.RWMutex

func accountConfigKey(accountID, region string) string {
	return accountID + "/" + region
}

// GetAccountConfiguration retrieves the account configuration for ACM certificates.
func (s *CertificateStore) GetAccountConfiguration(accountID, region string) (*AccountConfiguration, error) {
	accountConfigMutex.RLock()
	defer accountConfigMutex.RUnlock()

	key := accountConfigKey(accountID, region)
	if config, ok := accountConfigurations[key]; ok {
		return config, nil
	}
	return &AccountConfiguration{
		ExpiryEvents: ExpiryEventsConfiguration{
			DaysBeforeExpiry: 45,
		},
	}, nil
}

// PutAccountConfiguration stores the account configuration for ACM certificates.
func (s *CertificateStore) PutAccountConfiguration(accountID, region string, config *AccountConfiguration) error {
	accountConfigMutex.Lock()
	defer accountConfigMutex.Unlock()

	key := accountConfigKey(accountID, region)
	accountConfigurations[key] = config
	return nil
}

// DeleteAccountConfiguration removes the account configuration for ACM certificates.
func (s *CertificateStore) DeleteAccountConfiguration(accountID, region string) {
	accountConfigMutex.Lock()
	defer accountConfigMutex.Unlock()

	delete(accountConfigurations, accountConfigKey(accountID, region))
}
