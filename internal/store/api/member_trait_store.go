package api

import (
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/core/storage"
)

// MemberTraitStore stores member trait metadata.
type MemberTraitStore struct {
	bucket storage.Bucket
}

// NewMemberTraitStore creates a new member trait store.
func NewMemberTraitStore(store storage.BasicStorage) *MemberTraitStore {
	return &MemberTraitStore{
		bucket: store.Bucket("api_member_traits"),
	}
}

func (s *MemberTraitStore) makeKey(memberID int64, traitName string) []byte {
	return []byte(fmt.Sprintf("%d:%s", memberID, traitName))
}

// Put saves a member trait to the store.
func (s *MemberTraitStore) Put(trait *MemberTrait) error {
	data, err := json.Marshal(trait)
	if err != nil {
		return err
	}
	return s.bucket.Put(s.makeKey(trait.MemberID, trait.TraitName), data)
}

// GetByMemberID retrieves all traits for a member.
func (s *MemberTraitStore) GetByMemberID(memberID int64) ([]*MemberTrait, error) {
	prefix := []byte(fmt.Sprintf("%d:", memberID))
	var traits []*MemberTrait
	iter := s.bucket.ScanPrefix(prefix)
	defer iter.Close()

	for iter.Next() {
		var t MemberTrait
		if err := json.Unmarshal(iter.Value(), &t); err != nil {
			return nil, err
		}
		traits = append(traits, &t)
	}

	return traits, iter.Error()
}

// GetTraitValue retrieves a specific trait value for a member.
func (s *MemberTraitStore) GetTraitValue(memberID int64, traitName string) (string, error) {
	data, err := s.bucket.Get(s.makeKey(memberID, traitName))
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	var t MemberTrait
	if err := json.Unmarshal(data, &t); err != nil {
		return "", err
	}
	return t.TraitValue, nil
}

// Delete removes a member trait from the store.
func (s *MemberTraitStore) Delete(memberID int64, traitName string) error {
	return s.bucket.Delete(s.makeKey(memberID, traitName))
}

// DeleteByMemberID removes all traits for a member.
func (s *MemberTraitStore) DeleteByMemberID(memberID int64) error {
	prefix := []byte(fmt.Sprintf("%d:", memberID))
	var keys [][]byte
	iter := s.bucket.ScanPrefix(prefix)
	defer iter.Close()

	for iter.Next() {
		key := make([]byte, len(iter.Key()))
		copy(key, iter.Key())
		keys = append(keys, key)
	}

	if err := iter.Error(); err != nil {
		return err
	}

	for _, key := range keys {
		if err := s.bucket.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

// ForEach iterates over all member traits.
func (s *MemberTraitStore) ForEach(fn func(*MemberTrait) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		var t MemberTrait
		if err := json.Unmarshal(v, &t); err != nil {
			return err
		}
		return fn(&t)
	})
}

// Count returns the number of member traits in the store.
func (s *MemberTraitStore) Count() int {
	return s.bucket.Count()
}
