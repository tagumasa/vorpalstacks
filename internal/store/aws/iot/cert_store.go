package iot

import (
	"time"
	"vorpalstacks/internal/store/aws/common"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
)
func (s *IotStore) CreateCertificate(cert *Certificate) (*Certificate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cert.CertificateID == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Certificate{}
	if err := s.certsBase.GetProto(cert.CertificateID, existing); err == nil {
		return nil, ErrCertificateAlreadyExists
	}
	cert.CertificateARN = BuildCertificateARN(s.accountID, s.region, cert.CertificateID)
	return cert, s.certificatePS.Create(cert)
}

func (s *IotStore) GetCertificate(certificateID string) (*Certificate, error) {
	return s.certificatePS.Get(certificateID)
}

func (s *IotStore) UpdateCertificate(certID string, opts CertificateUpdateOpts) (*Certificate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.certificatePS.Get(certID)
	if err != nil {
		return nil, ErrCertificateNotFound
	}
	if opts.NewStatus != "" {
		existing.Status = opts.NewStatus
	}
	existing.LastModifiedDate = time.Now().UTC()
	return existing, s.certificatePS.Update(existing)
}

func (s *IotStore) DeleteCertificate(certificateID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cert, err := s.certificatePS.Get(certificateID)
	if err != nil {
		return err
	}
	if cert.Status != "INACTIVE" {
		return ErrInvalidCertStatus
	}
	certARN := cert.CertificateARN
	hasAttach := false
	s.policyAttachBase.ScanPrefix(certificateID+"\x00", func(key string, _ []byte) error {
		hasAttach = true
		return errFound
	})
	if !hasAttach {
		s.policyAttachBase.ScanPrefix(certARN+"\x00", func(key string, _ []byte) error {
			hasAttach = true
			return errFound
		})
	}
	if hasAttach {
		return ErrCertHasAttachments
	}
	hasThingAttach := false
	s.principalThingBase.ScanPrefix(certificateID+"\x00", func(key string, _ []byte) error {
		hasThingAttach = true
		return errFound
	})
	if !hasThingAttach {
		s.principalThingBase.ScanPrefix(certARN+"\x00", func(key string, _ []byte) error {
			hasThingAttach = true
			return errFound
		})
	}
	if hasThingAttach {
		return ErrCertHasAttachments
	}
	return s.certificatePS.DeleteIfExists(certificateID)
}

func (s *IotStore) ListCertificates(opts common.ListOptions) (*common.ListResult[Certificate], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.certsBase, opts, func() *pb.Certificate { return &pb.Certificate{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*Certificate, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToCertificate(p))
	}
	return &common.ListResult[Certificate]{Items: items, NextMarker: result.NextMarker}, nil
}
