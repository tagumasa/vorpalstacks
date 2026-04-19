package iam

import (
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
)

const groupBucketName = "iam_groups"

// GroupStore manages IAM group data in persistent storage.
type GroupStore struct {
	entityStore[Group]
	arnBuilder *ARNBuilder
}

// NewGroupStore creates a new store for IAM groups.
func NewGroupStore(store storage.BasicStorage, accountId string) *GroupStore {
	return &GroupStore{
		entityStore: newEntityStore[Group](store, groupBucketName),
		arnBuilder:  NewARNBuilder(accountId),
	}
}

// GetByArn retrieves a group by its ARN.
func (s *GroupStore) GetByArn(arn string) (*Group, error) {
	result, err := s.List("", "", 1000)
	if err != nil {
		return nil, err
	}
	for _, g := range result.Groups {
		if g.Arn == arn {
			return g, nil
		}
	}
	return nil, NewStoreError("get_group_by_arn", ErrGroupNotFound)
}

// GetByPath retrieves all groups with a given path prefix.
func (s *GroupStore) GetByPath(pathPrefix string) ([]*Group, error) {
	result, err := s.List("", "", 1000)
	if err != nil {
		return nil, err
	}
	if pathPrefix == "" {
		return result.Groups, nil
	}
	var filtered []*Group
	for _, g := range result.Groups {
		if strings.HasPrefix(g.Path, pathPrefix) {
			filtered = append(filtered, g)
		}
	}
	return filtered, nil
}

// Put stores a group.
func (s *GroupStore) Put(group *Group) error {
	return s.BaseStore.Put(group.GroupName, group)
}

// List returns a paginated list of groups.
func (s *GroupStore) List(pathPrefix, marker string, maxItems int) (*GroupListResult, error) {
	items, truncated, nextMarker, err := listEntitiesWithPathPrefix(s.BaseStore, pathPrefix, marker, maxItems, func(g *Group) string { return g.Path })
	if err != nil {
		return nil, err
	}
	return &GroupListResult{
		Groups:      items,
		IsTruncated: truncated,
		Marker:      nextMarker,
	}, nil
}

// Create creates a new IAM group.
func (s *GroupStore) Create(groupName, path, accountId string) (*Group, error) {
	if s.Exists(groupName) {
		return nil, NewStoreError("create_group", ErrGroupAlreadyExists)
	}

	groupID, err := GenerateGroupID()
	if err != nil {
		return nil, NewStoreError("generate_group_id", err)
	}

	group := &Group{
		ID:         groupID,
		Path:       path,
		GroupName:  groupName,
		AccountId:  accountId,
		CreateDate: time.Now().UTC(),
	}
	group.Arn = s.arnBuilder.GroupARN(path, groupName)

	if err := s.Put(group); err != nil {
		return nil, err
	}
	return group, nil
}

// Update updates an existing group.
func (s *GroupStore) Update(group *Group) error {
	return s.Put(group)
}
