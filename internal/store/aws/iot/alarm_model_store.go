package iot

import (
	"fmt"
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func (s *IotStore) CreateAlarmModel(a *AlarmModel) (*AlarmModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a.AlarmModelName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.AlarmModel{}
	if err := s.alarmModelsBase.GetProto(a.AlarmModelName, existing); err == nil {
		return nil, ErrAlarmModelAlreadyExists
	}
	a.AlarmModelARN = BuildAlarmModelARN(s.accountID, s.region, a.AlarmModelName)
	if a.CreationDate.IsZero() {
		a.CreationDate = time.Now().UTC()
	}
	a.LastModifiedDate = a.CreationDate
	if a.Status == "" {
		a.Status = "ACTIVE"
	}
	if a.AlarmModelVersion == "" {
		a.AlarmModelVersion = "1"
	}
	p, err := AlarmModelToProto(a)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise alarm model: %w", err)
	}
	return a, s.alarmModelsBase.PutProto(a.AlarmModelName, p)
}

func (s *IotStore) GetAlarmModel(name string) (*AlarmModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := &pb.AlarmModel{}
	if err := s.alarmModelsBase.GetProto(name, p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrAlarmModelNotFound
		}
		return nil, err
	}
	return ProtoToAlarmModel(p), nil
}

func (s *IotStore) UpdateAlarmModel(a *AlarmModel) (*AlarmModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing := &pb.AlarmModel{}
	if err := s.alarmModelsBase.GetProto(a.AlarmModelName, existing); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrAlarmModelNotFound
		}
		return nil, err
	}
	prev := ProtoToAlarmModel(existing)
	if a.AlarmModelDescription != "" {
		prev.AlarmModelDescription = a.AlarmModelDescription
	}
	if a.RoleARN != "" {
		prev.RoleARN = a.RoleARN
	}
	if a.AlarmModelDefinition != nil {
		prev.AlarmModelDefinition = a.AlarmModelDefinition
	}
	if a.Severity != "" {
		prev.Severity = a.Severity
	}
	prev.LastModifiedDate = time.Now().UTC()
	p, err := AlarmModelToProto(prev)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise alarm model: %w", err)
	}
	if err := s.alarmModelsBase.PutProto(prev.AlarmModelName, p); err != nil {
		return nil, err
	}
	return prev, nil
}

func (s *IotStore) DeleteAlarmModel(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.alarmModelsBase.Delete(name)
}

func (s *IotStore) ListAlarmModels(opts common.ListOptions) (*common.ListResult[AlarmModel], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.alarmModelsBase, opts, func() *pb.AlarmModel { return &pb.AlarmModel{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*AlarmModel, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToAlarmModel(p))
	}
	return &common.ListResult[AlarmModel]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) ListAlarmModelVersions(name string) ([]map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := &pb.AlarmModel{}
	if err := s.alarmModelsBase.GetProto(name, p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrAlarmModelNotFound
		}
		return nil, err
	}
	return []map[string]interface{}{
		{
			"alarmModelVersion": p.AlarmModelVersion,
			"creationDate":      protoToTime(p.CreationDate).Unix(),
			"lastModifiedDate":  protoToTime(p.LastModifiedDate).Unix(),
		},
	}, nil
}
