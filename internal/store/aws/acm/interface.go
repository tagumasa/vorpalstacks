package acm

import (
	"vorpalstacks/internal/store/aws/common"
)

// CertificateStoreInterface defines operations for managing ACM certificates.
type CertificateStoreInterface interface {
	Get(arn string) (*Certificate, error)
	GetByDomain(domain string) (*Certificate, error)
	List(marker string, maxItems int) (*CertificateListResult, error)
	ListAll() ([]*Certificate, error)
	Create(cert *Certificate) error
	Update(cert *Certificate) error
	Delete(arn string) error
	Exists(arn string) bool
	GetTags(arn string) ([]common.Tag, error)
	AddTags(arn string, tags []common.Tag) error
	RemoveTags(arn string, tagKeys []string) error
	GetAccountConfiguration(accountID, region string) (*AccountConfiguration, error)
	PutAccountConfiguration(accountID, region string, config *AccountConfiguration) error
	DeleteAccountConfiguration(accountID, region string)
}

var _ CertificateStoreInterface = (*CertificateStore)(nil)
