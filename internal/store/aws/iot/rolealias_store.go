package iot

import (
	"time"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func (s *IotStore) CreateRoleAlias(ra *RoleAlias) (*RoleAlias, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ra.RoleAlias == "" {
		return nil, ErrInvalidRequest
	}
	ra.RoleAliasARN = BuildRoleAliasARN(s.accountID, s.region, ra.RoleAlias)
	return ra, s.roleAliasPS.Create(ra)
}

func (s *IotStore) GetRoleAlias(alias string) (*RoleAlias, error) {
	return s.roleAliasPS.Get(alias)
}

// UpdateRoleAlias persists changes to an existing role alias.
func (s *IotStore) UpdateRoleAlias(alias string, opts RoleAliasUpdateOpts) (*RoleAlias, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.roleAliasPS.Get(alias)
	if err != nil {
		return nil, err
	}
	if opts.RoleARN != "" {
		existing.RoleARN = opts.RoleARN
	}
	if opts.DurationSeconds != 0 {
		existing.CredentialDurationSeconds = opts.DurationSeconds
	}
	existing.LastModifiedDate = time.Now().UTC()
	return existing, s.roleAliasPS.Update(existing)
}

func (s *IotStore) DeleteRoleAlias(alias string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.roleAliasPS.DeleteIfExists(alias)
}

func (s *IotStore) ListRoleAliases(opts common.ListOptions) (*common.ListResult[RoleAlias], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.roleAliasBase, opts, func() *pb.RoleAlias { return &pb.RoleAlias{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*RoleAlias, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToRoleAlias(p))
	}
	return &common.ListResult[RoleAlias]{Items: items, NextMarker: result.NextMarker}, nil
}
