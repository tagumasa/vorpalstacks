package iot

import (
	"fmt"
	"vorpalstacks/internal/store/aws/common"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
)
func (s *IotStore) CreateDetectorModel(d *DetectorModel) (*DetectorModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.DetectorModelName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.DetectorModel{}
	if err := s.detectorModelsBase.GetProto(d.DetectorModelName, existing); err == nil {
		return nil, ErrDetectorModelAlreadyExists
	}
	d.DetectorModelARN = BuildDetectorModelARN(s.accountID, s.region, d.DetectorModelName)
	pb, err := DetectorModelToProto(d)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise detector model: %w", err)
	}
	return d, s.detectorModelsBase.PutProto(d.DetectorModelName, pb)
}

func (s *IotStore) GetDetectorModel(name string) (*DetectorModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := &pb.DetectorModel{}
	if err := s.detectorModelsBase.GetProto(name, p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrDetectorModelNotFound
		}
		return nil, err
	}
	return ProtoToDetectorModel(p), nil
}

func (s *IotStore) UpdateDetectorModel(d *DetectorModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Existence check — cannot use ProtoStore because DetectorModelToProto returns an error.
	existing := &pb.DetectorModel{}
	if err := s.detectorModelsBase.GetProto(d.DetectorModelName, existing); err != nil {
		if common.IsNotFound(err) {
			return ErrDetectorModelNotFound
		}
		return err
	}
	pb, err := DetectorModelToProto(d)
	if err != nil {
		return fmt.Errorf("failed to serialise detector model: %w", err)
	}
	return s.detectorModelsBase.PutProto(d.DetectorModelName, pb)
}

func (s *IotStore) DeleteDetectorModel(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.detectorModelsBase.Delete(name)
}

func (s *IotStore) ListDetectorModels(opts common.ListOptions) (*common.ListResult[DetectorModel], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.detectorModelsBase, opts, func() *pb.DetectorModel { return &pb.DetectorModel{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*DetectorModel, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToDetectorModel(p))
	}
	return &common.ListResult[DetectorModel]{Items: items, NextMarker: result.NextMarker}, nil
}
