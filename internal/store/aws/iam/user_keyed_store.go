package iam

import (
	"encoding/json"

	"vorpalstacks/internal/store/aws/common"
)

// userKeyed provides the shared bodies for user-scoped credential records:
// records keyed by their own ID but queried by the owning user's name —
// signing certificates, SSH public keys and service-specific credentials.
// The per-type specifics (key field, user-name field, operation names and
// any rename-time derived fields) arrive as closures, following the
// principalKeyed precedent.
type userKeyed[T any] struct {
	*common.BaseStore
	kl           common.KeyLocker
	idFunc       func(*T) string
	userNameFunc func(*T) string
}

func newUserKeyed[T any](bs *common.BaseStore, idFunc, userNameFunc func(*T) string) userKeyed[T] {
	return userKeyed[T]{BaseStore: bs, idFunc: idFunc, userNameFunc: userNameFunc}
}

// listByUserName returns all records belonging to the given user.
func (s *userKeyed[T]) listByUserName(userName, opName string) ([]*T, error) {
	var records []*T
	err := s.ForEach(func(k string, v []byte) error {
		var record T
		if err := json.Unmarshal(v, &record); err != nil {
			return err
		}
		if s.userNameFunc(&record) == userName {
			records = append(records, &record)
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError(opName, err)
	}
	return records, nil
}

// deleteAllForUser removes all records belonging to the given user.
func (s *userKeyed[T]) deleteAllForUser(userName, opName string) error {
	var ids []string
	err := s.ForEach(func(k string, v []byte) error {
		var record T
		if err := json.Unmarshal(v, &record); err != nil {
			return err
		}
		if s.userNameFunc(&record) == userName {
			ids = append(ids, s.idFunc(&record))
		}
		return nil
	})
	if err != nil {
		return NewStoreError(opName, err)
	}
	for _, id := range ids {
		if err := s.BaseStore.Delete(id); err != nil {
			return err
		}
	}
	return nil
}

// migrateUser moves all records from oldUserName to newUserName; rename
// applies the new user name and any rename-time derived fields.
func (s *userKeyed[T]) migrateUser(oldUserName, newUserName, opName string, rename func(*T, string)) error {
	var toUpdate []*T
	err := s.ForEach(func(k string, v []byte) error {
		var record T
		if err := json.Unmarshal(v, &record); err != nil {
			return err
		}
		if s.userNameFunc(&record) == oldUserName {
			toUpdate = append(toUpdate, &record)
		}
		return nil
	})
	if err != nil {
		return NewStoreError(opName, err)
	}
	for _, record := range toUpdate {
		rename(record, newUserName)
		if err := s.BaseStore.Put(s.idFunc(record), record); err != nil {
			return err
		}
	}
	return nil
}

// updateStatus changes a record's status under its per-ID lock.  get
// resolves the record and maps its not-found error; setStatus applies the
// field.
func (s *userKeyed[T]) updateStatus(id, status string, get func(string) (*T, error), setStatus func(*T, string)) error {
	return s.kl.WithLock(id, func() error {
		record, err := get(id)
		if err != nil {
			return err
		}
		setStatus(record, status)
		return s.BaseStore.Put(s.idFunc(record), record)
	})
}
