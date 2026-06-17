package iot

import (
	"time"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func (s *IotStore) CreateDomainConfiguration(dc *DomainConfiguration) (*DomainConfiguration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if dc.DomainConfigurationName == "" {
		return nil, ErrInvalidRequest
	}
	dc.DomainConfigurationARN = BuildDomainConfigurationARN(s.accountID, s.region, dc.DomainConfigurationName)
	if dc.CreationDate.IsZero() {
		dc.CreationDate = time.Now().UTC()
	}
	if dc.LastModifiedDate.IsZero() {
		dc.LastModifiedDate = dc.CreationDate
	}
	if dc.DomainConfigurationStatus == "" {
		dc.DomainConfigurationStatus = "ENABLED"
	}
	return dc, s.domainConfigPS.Create(dc)
}

func (s *IotStore) GetDomainConfiguration(name string) (*DomainConfiguration, error) {
	return s.domainConfigPS.Get(name)
}

func (s *IotStore) UpdateDomainConfiguration(name string, dc *DomainConfiguration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.domainConfigPS.Get(name)
	if err != nil {
		return ErrDomainConfigurationNotFound
	}
	if dc.AuthorizerConfig != "" {
		existing.AuthorizerConfig = dc.AuthorizerConfig
	}
	if dc.DomainConfigurationStatus != "" {
		existing.DomainConfigurationStatus = dc.DomainConfigurationStatus
	}
	if dc.AuthenticationType != "" {
		existing.AuthenticationType = dc.AuthenticationType
	}
	if dc.ApplicationProtocol != "" {
		existing.ApplicationProtocol = dc.ApplicationProtocol
	}
	if dc.ServerCertificateARNs != nil {
		existing.ServerCertificateARNs = dc.ServerCertificateARNs
	}
	if dc.ValidationCertificateARN != "" {
		existing.ValidationCertificateARN = dc.ValidationCertificateARN
	}
	existing.LastModifiedDate = time.Now().UTC()
	return s.domainConfigPS.Update(existing)
}

func (s *IotStore) DeleteDomainConfiguration(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.domainConfigPS.DeleteIfExists(name)
}

func (s *IotStore) ListDomainConfigurations(opts common.ListOptions) (*common.ListResult[DomainConfiguration], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.domainConfigsBase, opts, func() *pb.DomainConfiguration { return &pb.DomainConfiguration{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*DomainConfiguration, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToDomainConfiguration(p))
	}
	return &common.ListResult[DomainConfiguration]{Items: items, NextMarker: result.NextMarker}, nil
}
