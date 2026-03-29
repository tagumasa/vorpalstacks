// Package apigateway provides API Gateway storage functionality for vorpalstacks.
package apigateway

import (
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// DomainStore provides storage operations for API Gateway domain names.
type DomainStore struct {
	*common.BaseStore
	arnBuilder *ArnBuilder
	accountId  string
	region     string
	mu         sync.Mutex
}

func domainBucketName(region string) string {
	return "apigateway-domains-" + region
}

// NewDomainStore creates a new DomainStore instance.
func NewDomainStore(store storage.BasicStorage, accountId, region string) *DomainStore {
	bucket := store.Bucket(domainBucketName(region))
	return &DomainStore{
		BaseStore:  common.NewBaseStore(bucket, "apigateway-domains"),
		arnBuilder: NewArnBuilder(accountId, region),
		accountId:  accountId,
		region:     region,
	}
}

// CreateDomainName creates a new domain name for API Gateway.
func (s *DomainStore) CreateDomainName(domain *DomainName) (*DomainName, error) {
	if domain.DomainName == "" {
		return nil, ErrInvalidParameter
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists("domain#" + domain.DomainName) {
		return nil, ErrDomainNameAlreadyExists
	}

	domain.DomainNameStatus = "AVAILABLE"
	domain.CertificateUploadDate = time.Now().UTC()

	domainId := generateId("", 12)
	domain.DomainNameId = domainId
	domain.DomainNameArn = s.arnBuilder.DomainNameArn(domain.DomainName)

	domain.DistributionDomainName = fmt.Sprintf("d%s.cloudfront.net", generateId("", 22))
	domain.DistributionHostedZoneId = "Z2FDTNDATAQYW2"

	domain.RegionalDomainName = fmt.Sprintf("d-%s.execute-api.%s.amazonaws.com", domainId, s.region)
	domain.RegionalHostedZoneId = "Z2OJLY3DKBEYEU"

	if err := s.Put("domain#"+domain.DomainName, domain); err != nil {
		return nil, err
	}

	return domain, nil
}

// GetDomainName retrieves a domain name by its name.
func (s *DomainStore) GetDomainName(domainName string) (*DomainName, error) {
	var domain DomainName
	if err := s.BaseStore.Get("domain#"+domainName, &domain); err != nil {
		return nil, ErrDomainNameNotFound
	}
	return &domain, nil
}

// UpdateDomainName updates an existing domain name.
func (s *DomainStore) UpdateDomainName(domain *DomainName) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists("domain#" + domain.DomainName) {
		return ErrDomainNameNotFound
	}
	return s.Put("domain#"+domain.DomainName, domain)
}

// DeleteDomainName deletes a domain name.
func (s *DomainStore) DeleteDomainName(domainName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists("domain#" + domainName) {
		return ErrDomainNameNotFound
	}

	mappings, err := s.ListBasePathMappings(domainName, common.ListOptions{MaxItems: 1000})
	if err == nil {
		for _, mapping := range mappings.Items {
			if err := s.DeleteBasePathMapping(domainName, mapping.BasePath); err != nil {
				logs.Error("Failed to delete base path mapping", logs.String("domain", domainName), logs.String("basePath", mapping.BasePath), logs.Err(err))
			}
		}
	}

	return s.BaseStore.Delete("domain#" + domainName)
}

// ListDomainNames returns all domain names.
func (s *DomainStore) ListDomainNames(opts common.ListOptions) (*common.ListResult[DomainName], error) {
	return common.List[DomainName](s.BaseStore, common.ListOptions{
		Prefix:   "domain#",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}, nil)
}

// CreateBasePathMapping creates a new base path mapping for a domain name.
func (s *DomainStore) CreateBasePathMapping(domainName string, mapping *BasePathMapping) (*BasePathMapping, error) {
	if domainName == "" {
		return nil, ErrInvalidParameter
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.GetDomainName(domainName); err != nil {
		return nil, err
	}

	key := "mapping#" + domainName + "#" + mapping.BasePath
	if s.Exists(key) {
		return nil, ErrBasePathMappingAlreadyExists
	}

	mapping.DomainName = domainName

	if err := s.Put(key, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

// GetBasePathMapping retrieves a base path mapping for a domain name.
func (s *DomainStore) GetBasePathMapping(domainName, basePath string) (*BasePathMapping, error) {
	var mapping BasePathMapping
	if err := s.BaseStore.Get("mapping#"+domainName+"#"+basePath, &mapping); err != nil {
		return nil, ErrBasePathMappingNotFound
	}
	return &mapping, nil
}

// UpdateBasePathMapping updates an existing base path mapping.
func (s *DomainStore) UpdateBasePathMapping(domainName, basePath string, mapping *BasePathMapping) error {
	key := "mapping#" + domainName + "#" + basePath
	if !s.Exists(key) {
		return ErrBasePathMappingNotFound
	}
	mapping.DomainName = domainName
	return s.Put(key, mapping)
}

// DeleteBasePathMapping deletes a base path mapping.
func (s *DomainStore) DeleteBasePathMapping(domainName, basePath string) error {
	key := "mapping#" + domainName + "#" + basePath
	if !s.Exists(key) {
		return ErrBasePathMappingNotFound
	}
	return s.BaseStore.Delete(key)
}

// ListBasePathMappings returns all base path mappings for a domain name.
func (s *DomainStore) ListBasePathMappings(domainName string, opts common.ListOptions) (*common.ListResult[BasePathMapping], error) {
	return common.List[BasePathMapping](s.BaseStore, common.ListOptions{
		Prefix:   "mapping#" + domainName + "#",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}, nil)
}
