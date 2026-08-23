package cognitoidentityprovider

import (
	"encoding/json"
	"strconv"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
)

// ListGroupsPaginated lists groups in a Cognito user pool with server-side pagination.
func (s *CognitoStore) ListGroupsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[Group], error) {
	opts.Prefix = userPoolID + "#"
	return common.List[Group](s.groupsStore, opts, nil)
}

// CreateGroup creates a new Cognito group.
func (s *CognitoStore) CreateGroup(group *Group) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if group.Name == "" {
		return ErrInvalidGroupName
	}

	key := userPoolGroupKey(group.UserPoolID, group.Name)
	if s.groupsStore.Exists(key) {
		return ErrGroupAlreadyExists
	}

	now := time.Now().UTC()
	group.CreationDate = now
	group.LastModifiedDate = now

	return s.groupsStore.Put(key, group)
}

// GetGroup retrieves a Cognito group by user pool ID and group name.
func (s *CognitoStore) GetGroup(userPoolID, groupName string) (*Group, error) {
	key := userPoolGroupKey(userPoolID, groupName)
	var group Group
	if err := s.groupsStore.Get(key, &group); err != nil {
		return nil, ErrGroupNotFound
	}
	return &group, nil
}

// UpdateGroup updates an existing Cognito group.
func (s *CognitoStore) UpdateGroup(group *Group) error {
	key := userPoolGroupKey(group.UserPoolID, group.Name)
	if !s.groupsStore.Exists(key) {
		return ErrGroupNotFound
	}
	group.LastModifiedDate = time.Now().UTC()
	return s.groupsStore.Put(key, group)
}

// DeleteGroup deletes a Cognito group.
func (s *CognitoStore) DeleteGroup(userPoolID, groupName string) error {
	s.groupMu.Lock()
	defer s.groupMu.Unlock()
	key := userPoolGroupKey(userPoolID, groupName)
	if !s.groupsStore.Exists(key) {
		return ErrGroupNotFound
	}
	group, err := s.GetGroup(userPoolID, groupName)
	if err != nil {
		return err
	}
	for _, member := range group.Members {
		user, err := s.GetUser(userPoolID, member)
		if err != nil {
			continue
		}
		var newGroups []string
		for _, g := range user.Groups {
			if g != groupName {
				newGroups = append(newGroups, g)
			}
		}
		user.Groups = newGroups
		// Key the write-back by the user's own canonical username rather
		// than the member string.
		if err := s.usersStore.Put(userPoolUserKey(userPoolID, user.Username), user); err != nil {
			logs.Warn("failed to update user after group deletion", logs.String("user", user.Username), logs.Err(err))
		}
	}
	return s.groupsStore.Delete(key)
}

// ListGroups lists all groups in a Cognito user pool.
func (s *CognitoStore) ListGroups(userPoolID string) ([]*Group, error) {
	var groups []*Group
	prefix := userPoolID + "#"
	err := s.groupsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var group Group
		if err := json.Unmarshal(value, &group); err != nil {
			return err
		}
		groups = append(groups, &group)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// AddUserToGroup adds a user to a Cognito group. The membership is keyed
// on the canonical stored username so that differently-cased spellings of
// the same user share one entry, matching the case-insensitive username
// resolution used by every other user lookup.
func (s *CognitoStore) AddUserToGroup(userPoolID, groupName, username string) error {
	s.groupMu.Lock()
	defer s.groupMu.Unlock()
	group, err := s.GetGroup(userPoolID, groupName)
	if err != nil {
		return err
	}

	canonical := s.resolveUsername(userPoolID, username)
	user, err := s.GetUser(userPoolID, canonical)
	if err != nil {
		return err
	}

	for _, member := range group.Members {
		if member == canonical {
			return ErrUserAlreadyInGroup
		}
	}

	group.Members = append(group.Members, canonical)
	if err := s.UpdateGroup(group); err != nil {
		return err
	}

	for _, g := range user.Groups {
		if g == groupName {
			return nil
		}
	}

	user.Groups = append(user.Groups, groupName)
	return s.UpdateUser(user)
}

// RemoveUserFromGroup removes a user from a Cognito group. Like the add
// path, the username is resolved to its canonical form first so a
// differently-cased spelling still matches the stored membership.
func (s *CognitoStore) RemoveUserFromGroup(userPoolID, groupName, username string) error {
	s.groupMu.Lock()
	defer s.groupMu.Unlock()
	group, err := s.GetGroup(userPoolID, groupName)
	if err != nil {
		return err
	}

	canonical := s.resolveUsername(userPoolID, username)
	found := false
	var newMembers []string
	for _, member := range group.Members {
		if member == canonical {
			found = true
		} else {
			newMembers = append(newMembers, member)
		}
	}

	if !found {
		return ErrUserNotInGroup
	}

	group.Members = newMembers
	if err := s.UpdateGroup(group); err != nil {
		return err
	}

	user, err := s.GetUser(userPoolID, canonical)
	if err != nil {
		return err
	}

	var newGroups []string
	for _, g := range user.Groups {
		if g != groupName {
			newGroups = append(newGroups, g)
		}
	}

	user.Groups = newGroups
	return s.UpdateUser(user)
}

// ListGroupsForUser lists all groups for a Cognito user.
func (s *CognitoStore) ListGroupsForUser(userPoolID, username string) ([]*Group, error) {
	user, err := s.GetUser(userPoolID, username)
	if err != nil {
		return nil, err
	}

	var groups []*Group
	for _, groupName := range user.Groups {
		group, err := s.GetGroup(userPoolID, groupName)
		if err == nil {
			groups = append(groups, group)
		}
	}
	return groups, nil
}

// ListUsersInGroup lists all users in a Cognito group.
func (s *CognitoStore) ListUsersInGroup(userPoolID, groupName string) ([]*User, error) {
	group, err := s.GetGroup(userPoolID, groupName)
	if err != nil {
		return nil, err
	}

	var users []*User
	for _, username := range group.Members {
		user, err := s.GetUser(userPoolID, username)
		if err == nil {
			users = append(users, user)
		}
	}
	return users, nil
}

// ListUsersInGroupPaginated returns a page of users belonging to the
// specified group. The Marker is a zero-based index into group.Members.
func (s *CognitoStore) ListUsersInGroupPaginated(userPoolID, groupName string, opts common.ListOptions) (*common.ListResult[User], error) {
	group, err := s.GetGroup(userPoolID, groupName)
	if err != nil {
		return nil, err
	}

	members := group.Members
	start := 0
	if opts.Marker != "" {
		if idx, perr := strconv.Atoi(opts.Marker); perr == nil && idx >= 0 && idx < len(members) {
			start = idx
		}
	}

	limit := opts.MaxItems
	if limit <= 0 {
		limit = 60
	}

	end := start + limit
	if end > len(members) {
		end = len(members)
	}

	var users []*User
	for i := start; i < end; i++ {
		user, err := s.GetUser(userPoolID, members[i])
		if err == nil {
			users = append(users, user)
		}
	}

	result := &common.ListResult[User]{
		Items:       users,
		IsTruncated: end < len(members),
	}
	if end < len(members) {
		result.NextMarker = strconv.Itoa(end)
	}
	return result, nil
}

// ListGroupsForUserPaginated returns a page of groups the specified user
// belongs to. The Marker is a zero-based index into user.Groups.
func (s *CognitoStore) ListGroupsForUserPaginated(userPoolID, username string, opts common.ListOptions) (*common.ListResult[Group], error) {
	user, err := s.GetUser(userPoolID, username)
	if err != nil {
		return nil, err
	}

	groupNames := user.Groups
	start := 0
	if opts.Marker != "" {
		if idx, perr := strconv.Atoi(opts.Marker); perr == nil && idx >= 0 && idx < len(groupNames) {
			start = idx
		}
	}

	limit := opts.MaxItems
	if limit <= 0 {
		limit = 60
	}

	end := start + limit
	if end > len(groupNames) {
		end = len(groupNames)
	}

	var groups []*Group
	for i := start; i < end; i++ {
		group, err := s.GetGroup(userPoolID, groupNames[i])
		if err == nil {
			groups = append(groups, group)
		}
	}

	result := &common.ListResult[Group]{
		Items:       groups,
		IsTruncated: end < len(groupNames),
	}
	if end < len(groupNames) {
		result.NextMarker = strconv.Itoa(end)
	}
	return result, nil
}
