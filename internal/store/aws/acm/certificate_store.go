// Package acm provides ACM (AWS Certificate Manager) storage functionality for vorpalstacks.
package acm

import (
	"sort"
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
	if arn == "" {
		return ""
	}
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

// ListWithFilters returns a paginated list of ACM certificates with advanced filtering.
func (s *CertificateStore) ListWithFilters(filters ListFilters, marker string, maxItems int) (*CertificateListResult, error) {
	allCerts, err := common.ListAll[Certificate](s.BaseStore)
	if err != nil {
		return nil, err
	}

	// Apply filters.
	filtered := make([]*Certificate, 0, len(allCerts))
	for _, cert := range allCerts {
		if !matchStatuses(cert, filters.Statuses) {
			continue
		}
		if !matchKeyTypes(cert, filters.KeyTypes) {
			continue
		}
		if !matchKeyUsage(cert, filters.KeyUsage) {
			continue
		}
		if !matchExtendedKeyUsage(cert, filters.ExtendedKeyUsage) {
			continue
		}
		if !matchExportOption(cert, filters.ExportOption) {
			continue
		}
		if !matchManagedBy(cert, filters.ManagedBy) {
			continue
		}
		if !matchOrigins(cert, filters.Origins) {
			continue
		}
		filtered = append(filtered, cert)
	}

	// Sort. AWS default sort order is ASCENDING when unspecified.
	if filters.SortBy == "CREATED_AT" {
		sort.Slice(filtered, func(i, j int) bool {
			if filters.SortOrder == "DESCENDING" {
				return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
			}
			return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
		})
	}

	// Paginate.
	start := 0
	if marker != "" {
		for i, cert := range filtered {
			if cert.CertificateArn == marker {
				start = i + 1
				break
			}
		}
	}

	end := len(filtered)
	if maxItems > 0 && start+maxItems < end {
		end = start + maxItems
	}
	page := filtered[start:end]

	summaries := make([]*CertificateSummary, len(page))
	for i, cert := range page {
		summaries[i] = CertificateToSummary(cert)
	}

	nextToken := ""
	if end < len(filtered) && len(page) > 0 {
		nextToken = page[len(page)-1].CertificateArn
	}

	return &CertificateListResult{
		Certificates: summaries,
		IsTruncated:  end < len(filtered),
		NextToken:    nextToken,
	}, nil
}

func matchStatuses(cert *Certificate, statuses []string) bool {
	if len(statuses) == 0 {
		return true
	}
	for _, s := range statuses {
		if cert.Status == s {
			return true
		}
	}
	return false
}

func matchKeyTypes(cert *Certificate, keyTypes []string) bool {
	if len(keyTypes) == 0 {
		return true
	}
	for _, kt := range keyTypes {
		if cert.KeyAlgorithm == kt {
			return true
		}
	}
	return false
}

func matchKeyUsage(cert *Certificate, keyUsages []string) bool {
	if len(keyUsages) == 0 {
		return true
	}
	for _, filterKU := range keyUsages {
		for _, certKU := range cert.KeyUsages {
			if certKU.Name == filterKU {
				return true
			}
		}
	}
	return false
}

func matchExtendedKeyUsage(cert *Certificate, ekus []string) bool {
	if len(ekus) == 0 {
		return true
	}
	for _, filterEKU := range ekus {
		for _, certEKU := range cert.ExtendedKeyUsages {
			if certEKU.Name == filterEKU || certEKU.OID == filterEKU {
				return true
			}
		}
	}
	return false
}

func matchExportOption(cert *Certificate, exportOption string) bool {
	if exportOption == "" {
		return true
	}
	if cert.Options == nil {
		return exportOption == "DISABLED"
	}
	return cert.Options.Export == exportOption
}

func matchManagedBy(cert *Certificate, managedBy string) bool {
	if managedBy == "" {
		return true
	}
	return cert.ManagedBy == managedBy
}

func matchOrigins(cert *Certificate, origins []string) bool {
	if len(origins) == 0 {
		return true
	}
	for _, origin := range origins {
		if cert.CertificateKeyPairOrigin == origin {
			return true
		}
	}
	return false
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

	summary.ManagedBy = cert.ManagedBy
	summary.CertificateKeyPairOrigin = cert.CertificateKeyPairOrigin

	if !cert.RevokedAt.IsZero() {
		summary.RevokedAt = formatEpochSeconds(cert.RevokedAt)
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
