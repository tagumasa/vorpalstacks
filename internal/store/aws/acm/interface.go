package acm

// ListFilters defines filter options for listing ACM certificates.
type ListFilters struct {
	Statuses         []string
	KeyTypes         []string
	KeyUsage         []string
	ExtendedKeyUsage []string
	ExportOption     string
	ManagedBy        string
	Origins          []string
	SortBy           string
	SortOrder        string
}

// CertificateStoreInterface defines operations for managing ACM certificates.
type CertificateStoreInterface interface {
	Get(arn string) (*Certificate, error)
	ListWithFilters(filters ListFilters, nextToken string, maxItems int) (*CertificateListResult, error)
	ListAll() ([]*Certificate, error)
	Create(cert *Certificate) error
	Update(cert *Certificate) error
	Delete(arn string) error
	Exists(arn string) bool
	GetAccountConfiguration(accountID, region string) (*AccountConfiguration, error)
	PutAccountConfiguration(accountID, region string, config *AccountConfiguration) error
	AddInUseBy(certArn, resourceArn string) error
	RemoveInUseBy(certArn, resourceArn string) error
}

var _ CertificateStoreInterface = (*CertificateStore)(nil)
