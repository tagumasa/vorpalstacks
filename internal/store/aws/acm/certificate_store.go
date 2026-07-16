// Package acm provides ACM (AWS Certificate Manager) storage functionality for vorpalstacks.
package acm

import (
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// CertificateStore provides storage operations for ACM certificates.
type CertificateStore struct {
	*common.BaseStore
	arnBuilder  *ARNBuilder
	mu          sync.RWMutex
	configStore *common.BaseStore
}

func certificateBucketName(region string) string {
	return "acm_certificates-" + region
}

func configBucketName(region string) string {
	return "acm_config-" + region
}

// NewCertificateStore creates a new CertificateStore instance with the specified storage, account ID, and region.
func NewCertificateStore(store storage.BasicStorage, accountId, region string) *CertificateStore {
	return &CertificateStore{
		BaseStore:   common.NewBaseStore(store.Bucket(certificateBucketName(region)), "acm"),
		arnBuilder:  NewARNBuilder(accountId, region),
		configStore: common.NewBaseStore(store.Bucket(configBucketName(region)), "acm-config"),
	}
}

func (s *CertificateStore) extractCertificateId(arn string) string {
	parts := strings.Split(arn, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return arn
}

// Get retrieves an ACM certificate by its ARN from the store.
func (s *CertificateStore) Get(arn string) (*Certificate, error) {
	certId := s.extractCertificateId(arn)
	var cert Certificate
	if err := s.BaseStore.Get(certId, &cert); err != nil {
		return nil, NewStoreError("get_certificate", err)
	}
	return &cert, nil
}

// List returns a paginated list of ACM certificates from the store.
func (s *CertificateStore) List(marker string, maxItems int) (*CertificateListResult, error) {
	opts := common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	}
	result, err := common.List[Certificate](s.BaseStore, opts, nil)
	if err != nil {
		return nil, err
	}
	summaries := make([]*CertificateSummary, len(result.Items))
	for i, cert := range result.Items {
		summaries[i] = CertificateToSummary(cert)
	}
	return &CertificateListResult{
		Certificates: summaries,
		IsTruncated:  result.IsTruncated,
		NextToken:    result.NextMarker,
	}, nil
}

// ListByStatus returns a paginated list of ACM certificates filtered by status.
func (s *CertificateStore) ListByStatus(statuses []string, marker string, maxItems int) (*CertificateListResult, error) {
	statusSet := make(map[string]bool, len(statuses))
	for _, s := range statuses {
		statusSet[s] = true
	}
	opts := common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	}
	result, err := common.List[Certificate](s.BaseStore, opts, func(cert *Certificate) bool {
		return statusSet[cert.Status]
	})
	if err != nil {
		return nil, err
	}
	summaries := make([]*CertificateSummary, len(result.Items))
	for i, cert := range result.Items {
		summaries[i] = CertificateToSummary(cert)
	}
	return &CertificateListResult{
		Certificates: summaries,
		IsTruncated:  result.IsTruncated,
		NextToken:    result.NextMarker,
	}, nil
}

// ListAll returns all ACM certificates from the store.
func (s *CertificateStore) ListAll() ([]*Certificate, error) {
	return common.ListAll[Certificate](s.BaseStore)
}

// Create creates a new ACM certificate in the store.
func (s *CertificateStore) Create(cert *Certificate) error {
	certId := s.extractCertificateId(cert.CertificateArn)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists(certId) {
		return NewStoreError("create_certificate", ErrCertificateExists)
	}
	if err := s.BaseStore.Put(certId, cert); err != nil {
		return NewStoreError("create_certificate", err)
	}
	return nil
}

// Update updates an existing ACM certificate in the store.
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

// CertificateToSummary converts a Certificate to a CertificateSummary.
func CertificateToSummary(cert *Certificate) *CertificateSummary {
	summary := &CertificateSummary{
		CertificateArn:                       cert.CertificateArn,
		DomainName:                           cert.DomainName,
		SubjectAlternativeNameSummaries:      cert.SubjectAlternativeNames,
		HasAdditionalSubjectAlternativeNames: len(cert.SubjectAlternativeNames) > 100,
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

func accountConfigKey(accountID, region string) string {
	return accountID + "/" + region
}

// GetAccountConfiguration retrieves the account configuration for ACM certificates.
func (s *CertificateStore) GetAccountConfiguration(accountID, region string) (*AccountConfiguration, error) {
	key := accountConfigKey(accountID, region)
	var config AccountConfiguration
	if err := s.configStore.Get(key, &config); err != nil {
		return &AccountConfiguration{
			ExpiryEvents: ExpiryEventsConfiguration{
				DaysBeforeExpiry: 45,
			},
		}, nil
	}
	return &config, nil
}

// PutAccountConfiguration stores the account configuration for ACM certificates.
func (s *CertificateStore) PutAccountConfiguration(accountID, region string, config *AccountConfiguration) error {
	key := accountConfigKey(accountID, region)
	return s.configStore.Put(key, config)
}

func formatEpochSeconds(t time.Time) float64 {
	return float64(t.Unix()) + float64(t.Nanosecond())/1e9
}
