package iot

import (
		"github.com/google/uuid"
	"time"
	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/store/aws/common"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
)
func (s *IotStore) CreateProvisioningTemplate(t *ProvisioningTemplate) (*ProvisioningTemplate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t.TemplateName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.ProvisioningTemplate{}
	if err := s.templatesBase.GetProto(t.TemplateName, existing); err == nil {
		return nil, ErrTemplateAlreadyExists
	}
	t.TemplateARN = BuildProvisioningTemplateARN(s.accountID, s.region, t.TemplateName)
	return t, s.provisioningTplPS.Create(t)
}

func (s *IotStore) GetProvisioningTemplate(name string) (*ProvisioningTemplate, error) {
	return s.provisioningTplPS.Get(name)
}

func (s *IotStore) UpdateProvisioningTemplate(name string, opts ProvisioningTemplateUpdateOpts) (*ProvisioningTemplate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.provisioningTplPS.Get(name)
	if err != nil {
		return nil, ErrTemplateNotFound
	}
	if opts.Description != "" {
		existing.Description = opts.Description
	}
	if opts.RoleARN != "" {
		existing.ProvisioningRoleARN = opts.RoleARN
	}
	if opts.Enabled != nil {
		existing.Enabled = *opts.Enabled
	}
	existing.LastModifiedDate = time.Now().UTC()
	return existing, s.provisioningTplPS.Update(existing)
}

func (s *IotStore) DeleteProvisioningTemplate(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.provisioningTplPS.DeleteIfExists(name)
}

// ListProvisioningTemplates returns all provisioning templates.
func (s *IotStore) ListProvisioningTemplates(opts common.ListOptions) (*common.ListResult[ProvisioningTemplate], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.templatesBase, opts, func() *pb.ProvisioningTemplate { return &pb.ProvisioningTemplate{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*ProvisioningTemplate, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToProvisioningTemplate(p))
	}
	return &common.ListResult[ProvisioningTemplate]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) CreateProvisioningTemplateVersion(name string, v *ProvisioningTemplateVersion) (*ProvisioningTemplateVersion, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v.VersionID == "" {
		v.VersionID = uuid.New().String()
	}
	if v.CreationDate.IsZero() {
		v.CreationDate = time.Now().UTC()
	}
	key := name + "\x00" + v.VersionID
	pbV, err := ProvisioningTemplateVersionToProto(v)
	if err != nil {
		return nil, err
	}
	return v, s.templatesBase.PutProto(key, pbV)
}

func (s *IotStore) GetProvisioningTemplateVersion(name, versionID string) (*ProvisioningTemplateVersion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := name + "\x00" + versionID
	pbV := &pb.ProvisioningTemplateVersion{}
	if err := s.templatesBase.GetProto(key, pbV); err != nil {
		return nil, ErrTemplateVersionNotFound
	}
	return ProtoToProvisioningTemplateVersion(pbV), nil
}

func (s *IotStore) DeleteProvisioningTemplateVersion(name, versionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := name + "\x00" + versionID
	return s.templatesBase.Delete(key)
}

func (s *IotStore) ListProvisioningTemplateVersions(name string, opts common.ListOptions) ([]*ProvisioningTemplateVersion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var versions []*ProvisioningTemplateVersion
	prefix := name + "\x00"
	err := s.templatesBase.ScanPrefix(prefix, func(key string, value []byte) error {
		pbV := &pb.ProvisioningTemplateVersion{}
		if proto.Unmarshal(value, pbV) == nil {
			versions = append(versions, ProtoToProvisioningTemplateVersion(pbV))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return versions, nil
}

var _ IotStoreInterface = (*IotStore)(nil)
