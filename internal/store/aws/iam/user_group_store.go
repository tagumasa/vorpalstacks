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

// userMembershipLockKey is the single lock key under which every mutation
// of one user's membership set runs — single additions, quota-checked
// additions, removals, sweeps and migrations. A membership key belongs to
// exactly one user, so the per-user lock fully serialises it; group-wide
// sweeps acquire the per-user locks in the sorted order their listing
// already yields, which keeps concurrent sweeps deadlock-free.
func userMembershipLockKey(userName string) string {
	return "user:" + userName
}

// AddUserToGroup adds a user to a group.
func (s *UserGroupStore) AddUserToGroup(userName, groupName string) error {
	return s.kl.WithLock(userMembershipLockKey(userName), func() error {
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

// AddUserToGroupWithLimit atomically checks the per-user group quota and
// writes the membership inside a single user-scoped lock, preventing the
// race condition where concurrent additions to different groups could both
// observe a count below the limit and both succeed.
func (s *UserGroupStore) AddUserToGroupWithLimit(userName, groupName string, maxGroups int) error {
	return s.kl.WithLock(userMembershipLockKey(userName), func() error {
		if s.IsUserInGroup(userName, groupName) {
			return NewStoreError("add_user_to_group", ErrUserAlreadyInGroup)
		}

		groups, err := s.ListGroupsForUser(userName)
		if err != nil {
			return err
		}
		if len(groups) >= maxGroups {
			return NewStoreError("add_user_to_group", ErrUserGroupLimitExceeded)
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
	return s.kl.WithLock(userMembershipLockKey(userName), func() error {
		if !s.IsUserInGroup(userName, groupName) {
			return NewStoreError("remove_user_from_group", ErrUserNotInGroup)
		}

		return s.BaseStore.Delete(s.membershipKey(userName, groupName))
	})
}

// MigrateUserGroup moves one user's membership between groups under the
// user's membership lock, add-first: the pair of writes is atomic against
// the quota check and every other membership mutation for that user, and a
// failure between the two writes leaves the member in both groups — which
// the caller's reverse step resolves — rather than in neither.
func (s *UserGroupStore) MigrateUserGroup(userName, fromGroup, toGroup string) error {
	return s.kl.WithLock(userMembershipLockKey(userName), func() error {
		if !s.IsUserInGroup(userName, fromGroup) {
			return NewStoreError("migrate_user_group", ErrUserNotInGroup)
		}
		membership := &UserGroupMembership{
			UserName:  userName,
			GroupName: toGroup,
			JoinDate:  time.Now().UTC(),
		}
		if err := s.BaseStore.Put(s.membershipKey(userName, toGroup), membership); err != nil {
			return NewStoreError("migrate_user_group", err)
		}
		if err := s.BaseStore.Delete(s.membershipKey(userName, fromGroup)); err != nil {
			return NewStoreError("migrate_user_group", err)
		}
		return nil
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
	return s.kl.WithLock(userMembershipLockKey(userName), func() error {
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

// RemoveAllUsersFromGroup removes all users from a group. The users are
// listed first, then each membership is deleted under its user's
// membership lock — the per-user locks are acquired in the listing's
// (sorted) order, so concurrent group-wide sweeps cannot deadlock.
func (s *UserGroupStore) RemoveAllUsersFromGroup(groupName string) error {
	users, err := s.ListUsersInGroup(groupName)
	if err != nil {
		return NewStoreError("remove_all_users_from_group", err)
	}
	for _, userName := range users {
		if err := s.kl.WithLock(userMembershipLockKey(userName), func() error {
			// A membership that vanished between the listing and this
			// lock is already gone; deleting an absent key is a no-op.
			return s.BaseStore.Delete(s.membershipKey(userName, groupName))
		}); err != nil {
			return NewStoreError("remove_all_users_from_group", err)
		}
	}
	return nil
}

// CountUsersInGroup returns the number of users in a group.  A listing
// failure propagates: callers checking a group's emptiness before deletion
// must not mistake a storage error for an empty group.
func (s *UserGroupStore) CountUsersInGroup(groupName string) (int, error) {
	users, err := s.ListUsersInGroup(groupName)
	if err != nil {
		return 0, err
	}
	return len(users), nil
}

// MigrateUser moves user memberships to a new user name. Both the old and
// the new name's locks are held — acquired in sorted order so a concurrent
// reverse migration cannot deadlock — because the move reads the old set
// and writes both sets.
func (s *UserGroupStore) MigrateUser(oldUserName, newUserName string) error {
	first, second := userMembershipLockKey(oldUserName), userMembershipLockKey(newUserName)
	if first > second {
		first, second = second, first
	}
	return s.kl.WithLock(first, func() error {
		return s.kl.WithLock(second, func() error {
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
	})
}
