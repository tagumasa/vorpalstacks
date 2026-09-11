package iam

import (
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
	var group *Group
	err := s.kl.WithLock(groupName, func() error {
		if s.Exists(groupName) {
			return NewStoreError("create_group", ErrGroupAlreadyExists)
		}

		groupID, err := GenerateGroupID()
		if err != nil {
			return NewStoreError("generate_group_id", err)
		}

		group = &Group{
			ID:         groupID,
			Path:       path,
			GroupName:  groupName,
			AccountId:  accountId,
			CreateDate: time.Now().UTC(),
		}
		group.Arn = s.arnBuilder.GroupARN(path, groupName)

		return s.Put(group)
	})
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Update updates an existing group.
func (s *GroupStore) Update(group *Group) error {
	return s.Put(group)
}
