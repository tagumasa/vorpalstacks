package iam

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const groupBucketName = "iam_groups"

// GroupStore manages IAM group data in persistent storage.
type GroupStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
}

// NewGroupStore creates a new store for IAM groups.
func NewGroupStore(store storage.BasicStorage, accountId string) *GroupStore {
	return &GroupStore{
		BaseStore:  common.NewBaseStore(store.Bucket(groupBucketName), "iam"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves a group by its name.
func (s *GroupStore) Get(groupName string) (*Group, error) {
	var group Group
	if err := s.BaseStore.Get(groupName, &group); err != nil {
		return nil, err
	}
	return &group, nil
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
	return nil, ErrGroupNotFound
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

// Delete removes a group by its name.
func (s *GroupStore) Delete(groupName string) error {
	return s.BaseStore.Delete(groupName)
}

// Exists checks whether a group exists.
func (s *GroupStore) Exists(groupName string) bool {
	return s.BaseStore.Exists(groupName)
}

// List returns a paginated list of groups.
func (s *GroupStore) List(pathPrefix, marker string, maxItems int) (*GroupListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var groups []*Group
	count := 0
	started := marker == ""
	hasMore := false

	err := s.ForEach(func(key string, value []byte) error {
		var group Group
		if err := json.Unmarshal(value, &group); err != nil {
			return err
		}

		if !started {
			if group.GroupName == marker {
				started = true
			}
			return nil
		}

		if pathPrefix != "" && !strings.HasPrefix(group.Path, pathPrefix) {
			return nil
		}

		if count < maxItems {
			groups = append(groups, &group)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_groups", err)
	}

	result := &GroupListResult{
		Groups:      groups,
		IsTruncated: hasMore,
	}
	if len(groups) > 0 {
		result.Marker = groups[len(groups)-1].GroupName
	}

	return result, nil
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

// Count returns the total number of groups.
func (s *GroupStore) Count() int {
	return s.BaseStore.Count()
}
