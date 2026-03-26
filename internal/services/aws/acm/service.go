// Package acm provides ACM (AWS Certificate Manager) service operations for vorpalstacks.
package acm

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	acmstore "vorpalstacks/internal/store/aws/acm"
)

// acmStores holds the various ACM stores.
type acmStores struct {
	certificates acmstore.CertificateStoreInterface
	arnBuilder   *acmstore.ARNBuilder
}

// ACMService provides ACM (AWS Certificate Manager) operations for managing certificates.
type ACMService struct{}

// NewACMService creates a new ACM service instance with the given storage, account ID, and region.
func NewACMService(store storage.BasicStorage, accountID, region string) *ACMService {
	return &ACMService{}
}

func (s *ACMService) store(reqCtx *request.RequestContext) (*acmStores, error) {
	if certStore := reqCtx.GetACMStore(); certStore != nil {
		return &acmStores{
			certificates: certStore,
			arnBuilder:   acmstore.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()),
		}, nil
	}
	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	return &acmStores{
		certificates: acmstore.NewCertificateStore(storage, reqCtx.GetAccountID(), reqCtx.GetRegion()),
		arnBuilder:   acmstore.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()),
	}, nil
}

// RegisterHandlers registers all ACM operation handlers with the dispatcher.
func (s *ACMService) RegisterHandlers(d *dispatcher.Dispatcher) {
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
