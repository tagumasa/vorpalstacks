package api

// Package api provides API model data store implementations for vorpalstacks.

import (
	"encoding/json"
	"strconv"

	"vorpalstacks/internal/core/storage"
)

// MemberStore stores API member definitions.
type MemberStore struct {
	bucket storage.Bucket
}

// NewMemberStore creates a new member store.
func NewMemberStore(store storage.BasicStorage) *MemberStore {
	return &MemberStore{
		bucket: store.Bucket("api_members"),
	}
}

func (s *MemberStore) makeKey(shapeID int64, memberName string) []byte {
	return []byte(int64ToString(shapeID) + ":" + memberName)
}

// Get retrieves a member by shape ID and name.
func (s *MemberStore) Get(shapeID int64, memberName string) (*Member, error) {
	data, err := s.bucket.Get(s.makeKey(shapeID, memberName))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var member Member
	if err := json.Unmarshal(data, &member); err != nil {
		return nil, err
	}
	return &member, nil
}

// Put saves a member to the store.
func (s *MemberStore) Put(member *Member) error {
	data, err := json.Marshal(member)
	if err != nil {
		return err
	}
	return s.bucket.Put(s.makeKey(member.ShapeID, member.Name), data)
}

// Delete removes a member from the store.
func (s *MemberStore) Delete(shapeID int64, memberName string) error {
	return s.bucket.Delete(s.makeKey(shapeID, memberName))
}

// ForEach iterates over all members.
func (s *MemberStore) ForEach(fn func(*Member) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		var member Member
		if err := json.Unmarshal(v, &member); err != nil {
			return err
		}
		return fn(&member)
	})
}

// ForEachByShape iterates over members for a specific shape.
func (s *MemberStore) ForEachByShape(shapeID int64, fn func(*Member) error) error {
	prefix := []byte(int64ToString(shapeID) + ":")
	iter := s.bucket.ScanPrefix(prefix)
	defer iter.Close()

	for iter.Next() {
		var member Member
		if err := json.Unmarshal(iter.Value(), &member); err != nil {
			return err
		}
		if err := fn(&member); err != nil {
			return err
		}
	}
	return iter.Error()
}

// Count returns the number of members in the store.
func (s *MemberStore) Count() int {
	return s.bucket.Count()
}

func int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}
