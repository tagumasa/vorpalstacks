package iam

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const userGroupBucketName = "iam_user_groups"

// UserGroupStore manages user-group membership.
type UserGroupStore struct {
	*common.BaseStore
	kl common.KeyLocker
}

// NewUserGroupStore creates a new UserGroupStore.
func NewUserGroupStore(store storage.BasicStorage) *UserGroupStore {
	return &UserGroupStore{
		BaseStore: common.NewBaseStore(store.Bucket(userGroupBucketName), "iam"),
	}
}

func (s *UserGroupStore) membershipKey(userName, groupName string) string {
	return userName + ":" + groupName
}

// AddUserToGroup adds a user to a group.
func (s *UserGroupStore) AddUserToGroup(userName, groupName string) error {
	lockKey := userName + ":" + groupName
	return s.kl.WithLock(lockKey, func() error {
		if s.IsUserInGroup(userName, groupName) {
			return NewStoreError("add_user_to_group", ErrUserAlreadyInGroup)
		}

		membership := &UserGroupMembership{
			UserName:  userName,
			GroupName: groupName,
			JoinDate:  time.Now().UTC(),
		}

		return s.BaseStore.Put(s.membershipKey(userName, groupName), membership)
	})
}

// RemoveUserFromGroup removes a user from a group.
func (s *UserGroupStore) RemoveUserFromGroup(userName, groupName string) error {
	lockKey := userName + ":" + groupName
	return s.kl.WithLock(lockKey, func() error {
		if !s.IsUserInGroup(userName, groupName) {
			return NewStoreError("remove_user_from_group", ErrUserNotInGroup)
		}

		return s.BaseStore.Delete(s.membershipKey(userName, groupName))
	})
}

// IsUserInGroup checks whether a user is in a group.
func (s *UserGroupStore) IsUserInGroup(userName, groupName string) bool {
	return s.BaseStore.Exists(s.membershipKey(userName, groupName))
}

// ListGroupsForUser returns the groups a user belongs to.
func (s *UserGroupStore) ListGroupsForUser(userName string) ([]string, error) {
	var groups []string
	prefix := userName + ":"

	err := s.ForEach(func(k string, v []byte) error {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			groups = append(groups, k[len(prefix):])
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_groups_for_user", err)
	}
	return groups, nil
}

// ListUsersInGroup returns the users in a group.
func (s *UserGroupStore) ListUsersInGroup(groupName string) ([]string, error) {
	var users []string
	suffix := ":" + groupName

	err := s.ForEach(func(k string, v []byte) error {
		if len(k) > len(suffix) && k[len(k)-len(suffix):] == suffix {
			users = append(users, k[:len(k)-len(suffix)])
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_users_in_group", err)
	}
	return users, nil
}

// RemoveAllGroupsForUser removes a user from all groups.
func (s *UserGroupStore) RemoveAllGroupsForUser(userName string) error {
	lockKey := "removeall:" + userName
	return s.kl.WithLock(lockKey, func() error {
		prefix := userName + ":"
		var keysToDelete []string

		err := s.ForEach(func(k string, v []byte) error {
			if len(k) > len(prefix) && k[:len(prefix)] == prefix {
				keysToDelete = append(keysToDelete, k)
			}
			return nil
		})

		if err != nil {
			return NewStoreError("remove_all_groups_for_user", err)
		}

		for _, key := range keysToDelete {
			if err := s.BaseStore.Delete(key); err != nil {
				return NewStoreError("remove_all_groups_for_user", err)
			}
		}
		return nil
	})
}

// RemoveAllUsersFromGroup removes all users from a group.
func (s *UserGroupStore) RemoveAllUsersFromGroup(groupName string) error {
	lockKey := "removeall:" + groupName
	return s.kl.WithLock(lockKey, func() error {
		suffix := ":" + groupName
		var keysToDelete []string

		err := s.ForEach(func(k string, v []byte) error {
			if len(k) > len(suffix) && k[len(k)-len(suffix):] == suffix {
				keysToDelete = append(keysToDelete, k)
			}
			return nil
		})

		if err != nil {
			return NewStoreError("remove_all_users_from_group", err)
		}

		for _, key := range keysToDelete {
			if err := s.BaseStore.Delete(key); err != nil {
				return NewStoreError("remove_all_users_from_group", err)
			}
		}
		return nil
	})
}

// CountUsersInGroup returns the number of users in a group.
func (s *UserGroupStore) CountUsersInGroup(groupName string) int {
	users, _ := s.ListUsersInGroup(groupName)
	return len(users)
}

// MigrateUser moves user memberships to a new user name.
func (s *UserGroupStore) MigrateUser(oldUserName, newUserName string) error {
	lockKey := "migrate:" + oldUserName
	return s.kl.WithLock(lockKey, func() error {
		prefix := oldUserName + ":"
		var memberships []UserGroupMembership

		err := s.ForEach(func(k string, v []byte) error {
			if len(k) > len(prefix) && k[:len(prefix)] == prefix {
				var membership UserGroupMembership
				if err := json.Unmarshal(v, &membership); err != nil {
					return err
				}
				memberships = append(memberships, membership)
			}
			return nil
		})

		if err != nil {
			return NewStoreError("migrate_user", err)
		}

		for _, membership := range memberships {
			if err := s.BaseStore.Delete(s.membershipKey(oldUserName, membership.GroupName)); err != nil {
				return NewStoreError("migrate_user", err)
			}
			membership.UserName = newUserName
			if err := s.BaseStore.Put(s.membershipKey(newUserName, membership.GroupName), &membership); err != nil {
				return NewStoreError("migrate_user", err)
			}
		}

		return nil
	})
}
