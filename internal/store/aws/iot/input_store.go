package iot

import (
	"fmt"
	"vorpalstacks/internal/store/aws/common"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
)
func (s *IotStore) CreateInput(i *Input) (*Input, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i.InputName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Input{}
	if err := s.inputsBase.GetProto(i.InputName, existing); err == nil {
		return nil, ErrInputAlreadyExists
	}
	i.InputARN = BuildInputARN(s.accountID, s.region, i.InputName)
	pb, err := InputToProto(i)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise input: %w", err)
	}
	return i, s.inputsBase.PutProto(i.InputName, pb)
}

func (s *IotStore) GetInput(name string) (*Input, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := &pb.Input{}
	if err := s.inputsBase.GetProto(name, p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrInputNotFound
		}
		return nil, err
	}
	return ProtoToInput(p), nil
}

func (s *IotStore) UpdateInput(i *Input) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing := &pb.Input{}
	if err := s.inputsBase.GetProto(i.InputName, existing); err != nil {
		if common.IsNotFound(err) {
			return ErrInputNotFound
		}
		return err
	}
	pb, err := InputToProto(i)
	if err != nil {
		return fmt.Errorf("failed to serialise input: %w", err)
	}
	return s.inputsBase.PutProto(i.InputName, pb)
}

func (s *IotStore) DeleteInput(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.inputsBase.Delete(name)
}

func (s *IotStore) ListInputs(opts common.ListOptions) (*common.ListResult[Input], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.inputsBase, opts, func() *pb.Input { return &pb.Input{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*Input, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToInput(p))
	}
	return &common.ListResult[Input]{Items: items, NextMarker: result.NextMarker}, nil
}

// UpdateJob persists changes to an existing job.
