// Package acm provides ACM (AWS Certificate Manager) service operations for vorpalstacks.
package acm

import (
	"context"
	"fmt"
	"sync"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	acmstore "vorpalstacks/internal/store/aws/acm"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// acmStores holds the various ACM stores together with the account/region
// binding they were constructed for, so Core functions operating on a
// resolved stores value never need the transport request context.
type acmStores struct {
	certificates acmstore.CertificateStoreInterface
	arnBuilder   *acmstore.ARNBuilder
	accountID    string
	region       string
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

// createStores builds the per-region ACM store bundle. Called by both
// store() and GetStoreForRegion() so a single construction path defines
// what an acmStores carries.
func (s *ACMService) createStores(st storage.BasicStorage, region string) *acmStores {
	return &acmStores{
		certificates: acmstore.NewCertificateStore(st, s.accountID, region),
		arnBuilder:   acmstore.NewARNBuilder(s.accountID, region),
		accountID:    s.accountID,
		region:       region,
	}
}

func (s *ACMService) store(reqCtx *request.RequestContext) (*acmStores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (*acmStores, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return s.createStores(storage, reqCtx.GetRegion()), nil
	})
}

// GetStoreForRegion returns the cached ACM stores for the given region,
// creating a new store instance if not already cached. Invoker-facing entry
// (cross-service certificate registration, admin console); shares the same
// per-region cache and the same construction path as store().
func (s *ACMService) GetStoreForRegion(region string) (*acmStores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (*acmStores, error) {
		if s.storageManager == nil {
			return nil, fmt.Errorf("acm storage manager not initialised")
		}
		st, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}
		return s.createStores(st, region), nil
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
	d.RegisterHandlerForService("acm", "TagResource", s.TagResource)
	d.RegisterHandlerForService("acm", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("acm", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("acm", "SearchCertificates", s.SearchCertificates)
}

// RegisterCertificateUsage implements invokers.ACMInvoker. It records that
// the resource identified by resourceArn is now using the ACM certificate
// identified by certArn, in the specified region.
func (s *ACMService) RegisterCertificateUsage(ctx context.Context, region, certArn, resourceArn string) error {
	stores, err := s.GetStoreForRegion(region)
	if err != nil {
		return fmt.Errorf("acm: failed to get store for region %s: %w", region, err)
	}
	return stores.certificates.AddInUseBy(certArn, resourceArn)
}

// UnregisterCertificateUsage implements invokers.ACMInvoker. It removes the
// resource identified by resourceArn from the certificate's InUseBy list,
// indicating the resource no longer references the certificate.
func (s *ACMService) UnregisterCertificateUsage(ctx context.Context, region, certArn, resourceArn string) error {
	stores, err := s.GetStoreForRegion(region)
	if err != nil {
		return fmt.Errorf("acm: failed to get store for region %s: %w", region, err)
	}
	return stores.certificates.RemoveInUseBy(certArn, resourceArn)
}

// CertificateExists implements invokers.ACMInvoker. It performs a
// pre-validation check before a cross-service consumer saves a resource that
// references an ACM certificate. This eliminates the most common failure mode
// (invalid cert ARN) before the resource is created, reducing the compensating
// transaction path to Pebble I/O errors only.
func (s *ACMService) CertificateExists(ctx context.Context, region, certArn string) bool {
	stores, err := s.GetStoreForRegion(region)
	if err != nil {
		return false
	}
	return stores.certificates.Exists(certArn)
}

// CertificateMaterial implements invokers.ACMCertificateProvider. It returns
// the PEM material of an issued certificate so a cross-service listener
// (e.g. the CloudFront distribution plane) can terminate TLS with it. Only
// issued certificates with a retained key pair resolve; anything else is an
// error so the caller fails the handshake instead of serving a mismatched
// certificate.
func (s *ACMService) CertificateMaterial(ctx context.Context, region, certArn string) (invokers.TLSCertificateMaterial, error) {
	stores, err := s.GetStoreForRegion(region)
	if err != nil {
		return invokers.TLSCertificateMaterial{}, fmt.Errorf("acm: failed to get store for region %s: %w", region, err)
	}
	cert, err := stores.certificates.Get(certArn)
	if err != nil {
		return invokers.TLSCertificateMaterial{}, fmt.Errorf("acm: certificate %s not found: %w", certArn, err)
	}
	if cert.Status != "ISSUED" {
		return invokers.TLSCertificateMaterial{}, awserrors.NewValidationException("Certificate is not in the ISSUED state: " + certArn)
	}
	if cert.Certificate == "" || cert.PrivateKey == "" {
		return invokers.TLSCertificateMaterial{}, awserrors.NewValidationException("Certificate does not have serving material available: " + certArn)
	}
	return invokers.TLSCertificateMaterial{
		Certificate:      cert.Certificate,
		PrivateKey:       cert.PrivateKey,
		CertificateChain: cert.CertificateChain,
	}, nil
}
