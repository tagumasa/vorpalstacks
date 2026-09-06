// Package apigateway provides API Gateway storage functionality for vorpalstacks.
package apigateway

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

var errDomainStopIteration = errors.New("stop domain iteration")

// DomainStore provides storage operations for API Gateway domain names.
type DomainStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
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
		arnBuilder: NewARNBuilder(accountId, region),
		accountId:  accountId,
		region:     region,
	}
}

// deleteBasePathMappingLocked removes a base path mapping without acquiring
// the store lock. The caller must hold s.mu.
func (s *DomainStore) deleteBasePathMappingLocked(domainName, basePath string) error {
	key := "mapping#" + domainName + "#" + basePath
	if !s.Exists(key) {
		return ErrBasePathMappingNotFound
	}
	return s.BaseStore.Delete(key)
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
		if common.IsNotFound(err) {
			return nil, ErrDomainNameNotFound
		}
		return nil, err
	}
	return &domain, nil
}

// GetDomainNameById retrieves a domain name by its domainNameId.
// This scans all domain entries to find the one with matching DomainNameId.
func (s *DomainStore) GetDomainNameById(domainNameId string) (*DomainName, error) {
	var found *DomainName
	err := common.ForEachAll[DomainName](s.BaseStore, "domain#", nil, func(d *DomainName) error {
		if d.DomainNameId == domainNameId {
			found = d
			return errDomainStopIteration
		}
		return nil
	})
	if err != nil && err != errDomainStopIteration {
		return nil, err
	}
	if found == nil {
		return nil, ErrDomainNameNotFound
	}
	return found, nil
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

// DeleteDomainName deletes a domain name and its base path mappings.
func (s *DomainStore) DeleteDomainName(domainName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists("domain#" + domainName) {
		return ErrDomainNameNotFound
	}

	if err := common.ForEachAll[BasePathMapping](s.BaseStore, "mapping#"+domainName+"#", nil, func(m *BasePathMapping) error {
		if delErr := s.deleteBasePathMappingLocked(domainName, m.BasePath); delErr != nil {
			logs.Error("Failed to delete base path mapping", logs.String("domain", domainName), logs.String("basePath", m.BasePath), logs.Err(delErr))
		}
		return nil
	}); err != nil {
		logs.Error("Failed to clean up base path mappings", logs.String("domain", domainName), logs.Err(err))
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
		if common.IsNotFound(err) {
			return nil, ErrBasePathMappingNotFound
		}
		return nil, err
	}
	return &mapping, nil
}

// UpdateBasePathMapping updates an existing base path mapping.
func (s *DomainStore) UpdateBasePathMapping(domainName, basePath string, mapping *BasePathMapping) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := "mapping#" + domainName + "#" + basePath
	if !s.Exists(key) {
		return ErrBasePathMappingNotFound
	}
	mapping.DomainName = domainName
	return s.Put(key, mapping)
}

// DeleteBasePathMapping deletes a base path mapping.
func (s *DomainStore) DeleteBasePathMapping(domainName, basePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.deleteBasePathMappingLocked(domainName, basePath)
}

// ListBasePathMappings returns all base path mappings for a domain name.
func (s *DomainStore) ListBasePathMappings(domainName string, opts common.ListOptions) (*common.ListResult[BasePathMapping], error) {
	return common.List[BasePathMapping](s.BaseStore, common.ListOptions{
		Prefix:   "mapping#" + domainName + "#",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}, nil)
}

// TagDomainName adds or updates tags on a domain name.
func (s *DomainStore) TagDomainName(domainName string, inputTags map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	domain, err := s.GetDomainName(domainName)
	if err != nil {
		return err
	}
	if domain.Tags == nil {
		domain.Tags = []tags.Tag{}
	}
	domain.Tags = tags.Apply(domain.Tags, tags.MapToTags(inputTags))
	return s.Put("domain#"+domainName, domain)
}

// UntagDomainName removes tags from a domain name.
func (s *DomainStore) UntagDomainName(domainName string, tagKeys []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	domain, err := s.GetDomainName(domainName)
	if err != nil {
		return err
	}
	domain.Tags = tags.RemoveByTagKeys(domain.Tags, tagKeys)
	return s.Put("domain#"+domainName, domain)
}

// GetDomainNameTags returns tags for a domain name.
func (s *DomainStore) GetDomainNameTags(domainName string) ([]tags.Tag, error) {
	domain, err := s.GetDomainName(domainName)
	if err != nil {
		return nil, err
	}
	return domain.Tags, nil
}

// RemoveBasePathMappingsForApi deletes all base path mappings that reference
// the given restApiId across all domain names. This is called when an API is
// deleted to prevent dangling references. Both the domain walk and the
// per-domain mapping walk use ForEachAll so every page is visited — a
// fixed-page listing leaves mappings beyond the first page as dangling
// references.
func (s *DomainStore) RemoveBasePathMappingsForApi(restApiId string) error {
	type mappingRef struct {
		domain   string
		basePath string
	}
	var targets []mappingRef
	err := common.ForEachAll[DomainName](s.BaseStore, "domain#", nil, func(d *DomainName) error {
		// Best-effort per domain: a listing failure for one domain must not
		// abort the sweep of the remaining domains.
		_ = common.ForEachAll[BasePathMapping](s.BaseStore, "mapping#"+d.DomainName+"#", nil, func(m *BasePathMapping) error {
			if m.RestApiId == restApiId {
				targets = append(targets, mappingRef{domain: d.DomainName, basePath: m.BasePath})
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return err
	}
	// Deletions run after the walk so the iterator never observes its own
	// writes.
	for _, t := range targets {
		_ = s.DeleteBasePathMapping(t.domain, t.basePath)
	}
	return nil
}
