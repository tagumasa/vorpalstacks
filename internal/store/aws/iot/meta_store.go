package iot

import (
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func (s *IotStore) GetIndexingConfiguration() (*IndexingConfiguration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	pbIC := &pb.IndexingConfiguration{}
	if err := s.indexingConfigBase.GetProto("default", pbIC); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrIndexingConfigurationNotFound
		}
		return nil, err
	}
	return ProtoToIndexingConfiguration(pbIC), nil
}

func (s *IotStore) UpdateIndexingConfiguration(ic *IndexingConfiguration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	pbIC := IndexingConfigurationToProto(ic)
	return s.indexingConfigBase.PutProto("default", pbIC)
}
