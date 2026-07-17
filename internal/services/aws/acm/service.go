// Package acm provides ACM (AWS Certificate Manager) service operations for vorpalstacks.
package acm

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	acmstore "vorpalstacks/internal/store/aws/acm"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// acmStores holds the various ACM stores.
type acmStores struct {
	certificates acmstore.CertificateStoreInterface
	arnBuilder   *acmstore.ARNBuilder
}

// ACMService provides ACM (AWS Certificate Manager) operations for managing certificates.
type ACMService struct {
	accountID      string
	region         string
	stores         sync.Map // region → *acmStores
	storageManager *storage.RegionStorageManager
}

// NewACMService creates a new ACM service instance with the given storage, account ID, and region.
func NewACMService(accountID, region string) *ACMService {
	return &ACMService{
		accountID: accountID,
		region:    region,
	}
}

// SetStorageManager injects the region storage manager for lazy store creation.
func (s *ACMService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

func (s *ACMService) store(reqCtx *request.RequestContext) (*acmStores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (*acmStores, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return &acmStores{
			certificates: acmstore.NewCertificateStore(storage, s.accountID, reqCtx.GetRegion()),
			arnBuilder:   acmstore.NewARNBuilder(s.accountID, reqCtx.GetRegion()),
		}, nil
	})
}

// GetStoreForRegion returns the cached ACM stores for the given region,
// creating a new store instance if not already cached.
func (s *ACMService) GetStoreForRegion(region string) (*acmStores, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(*acmStores), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("acm storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	stores := &acmStores{
		certificates: acmstore.NewCertificateStore(st, s.accountID, region),
		arnBuilder:   acmstore.NewARNBuilder(s.accountID, region),
	}
	actual, _ := s.stores.LoadOrStore(region, stores)
	return actual.(*acmStores), nil
}

// RegisterHandlers registers all ACM operation handlers with the dispatcher.
func (s *ACMService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("acm", "RequestCertificate", s.RequestCertificate)
	d.RegisterHandlerForService("acm", "GetCertificate", s.GetCertificate)
	d.RegisterHandlerForService("acm", "ListCertificates", s.ListCertificates)
	d.RegisterHandlerForService("acm", "DeleteCertificate", s.DeleteCertificate)
	d.RegisterHandlerForService("acm", "DescribeCertificate", s.DescribeCertificate)
	d.RegisterHandlerForService("acm", "ResendValidationEmail", s.ResendValidationEmail)
	d.RegisterHandlerForService("acm", "AddTagsToCertificate", s.AddTagsToCertificate)
	d.RegisterHandlerForService("acm", "RemoveTagsFromCertificate", s.RemoveTagsFromCertificate)
	d.RegisterHandlerForService("acm", "ListTagsForCertificate", s.ListTagsForCertificate)
	d.RegisterHandlerForService("acm", "ImportCertificate", s.ImportCertificate)
	d.RegisterHandlerForService("acm", "ExportCertificate", s.ExportCertificate)
	d.RegisterHandlerForService("acm", "GetAccountConfiguration", s.GetAccountConfiguration)
	d.RegisterHandlerForService("acm", "PutAccountConfiguration", s.PutAccountConfiguration)
	d.RegisterHandlerForService("acm", "UpdateCertificateOptions", s.UpdateCertificateOptions)
	d.RegisterHandlerForService("acm", "RenewCertificate", s.RenewCertificate)
	d.RegisterHandlerForService("acm", "RevokeCertificate", s.RevokeCertificate)
	d.RegisterHandlerForService("acm", "TagResource", s.TagResource)
	d.RegisterHandlerForService("acm", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("acm", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("acm", "SearchCertificates", s.SearchCertificates)
}
