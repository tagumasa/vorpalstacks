package iot

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
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
		return nil, err
	}
	// A defaultVersionId that matches no stored version must not reach the
	// template record: it would point the template at a version nothing
	// can resolve while the previously flagged version keeps claiming to
	// be the default. A valid repoint also moves the versions'
	// IsDefaultVersion flags so exactly one version claims the default.
	if opts.DefaultVersionID != 0 {
		versionIDStr := fmt.Sprintf("%d", opts.DefaultVersionID)
		versions, listErr := s.listProvisioningTemplateVersionsLocked(name)
		if listErr != nil {
			return nil, listErr
		}
		if !provisioningVersionExists(versions, versionIDStr) {
			return nil, ErrTemplateVersionNotFound
		}
		if err := s.repointDefaultVersionFlagsLocked(name, versionIDStr, versions); err != nil {
			return nil, err
		}
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
	if opts.TemplateBody != nil {
		existing.TemplateBody = *opts.TemplateBody
	}
	if opts.DefaultVersionID != 0 {
		existing.DefaultVersionID = opts.DefaultVersionID
	}
	if opts.PreProvisioningHook != "" {
		existing.PreProvisioningHook = opts.PreProvisioningHook
	}
	if opts.RemovePreProvisioningHook {
		existing.PreProvisioningHook = ""
	}
	existing.LastModifiedDate = time.Now().UTC()
	return existing, s.provisioningTplPS.Update(existing)
}

// SetDefaultProvisioningTemplateVersion repoints a template's default
// version atomically: the template record, the IsDefaultVersion flag of
// the chosen version, and the flags of every sibling are updated under one
// lock so exactly one version reports itself as the default — the state
// ListProvisioningTemplateVersions exposes to clients. A versionID that
// matches no stored version is rejected before any mutation; clearing the
// sibling flags first would leave the template with no default at all.
func (s *IotStore) SetDefaultProvisioningTemplateVersion(name string, versionID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	template, err := s.provisioningTplPS.Get(name)
	if err != nil {
		return err
	}

	versions, err := s.listProvisioningTemplateVersionsLocked(name)
	if err != nil {
		return err
	}

	versionIDStr := fmt.Sprintf("%d", versionID)
	if !provisioningVersionExists(versions, versionIDStr) {
		return ErrTemplateVersionNotFound
	}
	if err := s.repointDefaultVersionFlagsLocked(name, versionIDStr, versions); err != nil {
		return err
	}

	template.DefaultVersionID = versionID
	template.LastModifiedDate = time.Now().UTC()
	return s.provisioningTplPS.Update(template)
}

// provisioningVersionExists reports whether the version list contains the
// given version ID.
func provisioningVersionExists(versions []*ProvisioningTemplateVersion, versionIDStr string) bool {
	for _, v := range versions {
		if v.VersionID == versionIDStr {
			return true
		}
	}
	return false
}

// repointDefaultVersionFlagsLocked clears the IsDefaultVersion flag of
// every sibling and sets it on the chosen version, so exactly one stored
// version claims to be the template's default. Callers hold s.mu.
func (s *IotStore) repointDefaultVersionFlagsLocked(name, versionIDStr string, versions []*ProvisioningTemplateVersion) error {
	for _, v := range versions {
		isDefault := v.VersionID == versionIDStr
		if v.IsDefaultVersion == isDefault {
			continue
		}
		v.IsDefaultVersion = isDefault
		if _, err := s.putProvisioningTemplateVersionLocked(name, v); err != nil {
			return err
		}
	}
	return nil
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
		t, err := ProtoToProvisioningTemplate(p)
		if err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	return &common.ListResult[ProvisioningTemplate]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) CreateProvisioningTemplateVersion(name string, v *ProvisioningTemplateVersion) (*ProvisioningTemplateVersion, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.putProvisioningTemplateVersionLocked(name, v)
}

// putProvisioningTemplateVersionLocked persists one version record.
// Callers must hold s.mu (the mutex is not reentrant, so compound
// operations reuse this form instead of the public method).
func (s *IotStore) putProvisioningTemplateVersionLocked(name string, v *ProvisioningTemplateVersion) (*ProvisioningTemplateVersion, error) {
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
		if common.IsNotFound(err) {
			return nil, ErrTemplateVersionNotFound
		}
		return nil, err
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
	return s.listProvisioningTemplateVersionsLocked(name)
}

// listProvisioningTemplateVersionsLocked scans a template's version
// records. Callers must hold s.mu.
func (s *IotStore) listProvisioningTemplateVersionsLocked(name string) ([]*ProvisioningTemplateVersion, error) {
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
