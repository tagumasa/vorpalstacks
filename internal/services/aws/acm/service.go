// Package acm provides ACM (AWS Certificate Manager) service operations for vorpalstacks.
package acm

import (
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
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
	accountID string
	region    string
	stores    sync.Map // region → *acmStores
}

// NewACMService creates a new ACM service instance with the given storage, account ID, and region.
func NewACMService(accountID, region string) *ACMService {
	return &ACMService{
		accountID: accountID,
		region:    region,
	}
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
}
