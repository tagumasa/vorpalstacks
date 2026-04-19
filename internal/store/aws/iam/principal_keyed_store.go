package iam

import (
	"encoding/json"

	"vorpalstacks/internal/store/aws/common"
)

// principalKeyed provides shared prefix-based CRUD operations for IAM stores
// that use a three-part key scheme: principalType:principalName:thirdKey.
// Both InlinePolicyStore and AttachedPolicyStore follow this pattern.
type principalKeyed[T any] struct {
	*common.BaseStore
	kl          common.KeyLocker
	keyFunc     func(*T) string
	migrateFunc func(*T, string)
}

func newPrincipalKeyed[T any](bs *common.BaseStore, keyFunc func(*T) string, migrateFunc func(*T, string)) principalKeyed[T] {
	return principalKeyed[T]{BaseStore: bs, keyFunc: keyFunc, migrateFunc: migrateFunc}
}

func (s *principalKeyed[T]) buildPrefix(principalType, principalName string) string {
	return principalType + ":" + principalName + ":"
}

func (s *principalKeyed[T]) listByPrincipal(principalType, principalName string) ([]string, error) {
	var names []string
	prefix := s.buildPrefix(principalType, principalName)

	err := s.ForEach(func(k string, _ []byte) error {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			names = append(names, k[len(prefix):])
		}
		return nil
	})

	return names, err
}

func (s *principalKeyed[T]) deleteAllForPrincipal(principalType, principalName, opName string) error {
	lockKey := "removeall:" + principalType + ":" + principalName
	return s.kl.WithLock(lockKey, func() error {
		prefix := s.buildPrefix(principalType, principalName)
		var keysToDelete []string

		err := s.ForEach(func(k string, _ []byte) error {
			if len(k) > len(prefix) && k[:len(prefix)] == prefix {
				keysToDelete = append(keysToDelete, k)
			}
			return nil
		})

		if err != nil {
			return NewStoreError(opName, err)
		}

		for _, key := range keysToDelete {
			if err := s.BaseStore.Delete(key); err != nil {
				return NewStoreError(opName, err)
			}
		}
		return nil
	})
}

func (s *principalKeyed[T]) migratePrincipal(oldName, newName, principalType, opName string) error {
	lockKey := "migrate:" + principalType + ":" + oldName
	return s.kl.WithLock(lockKey, func() error {
		prefix := s.buildPrefix(principalType, oldName)
		var items []T

		err := s.ForEach(func(k string, v []byte) error {
			if len(k) > len(prefix) && k[:len(prefix)] == prefix {
				var item T
				if err := json.Unmarshal(v, &item); err != nil {
					return err
				}
				items = append(items, item)
			}
			return nil
		})

		if err != nil {
			return NewStoreError(opName, err)
		}

		oldPrefix := s.buildPrefix(principalType, oldName)
		newPrefix := s.buildPrefix(principalType, newName)
		for i := range items {
			item := &items[i]
			if err := s.BaseStore.Delete(oldPrefix + s.keyFunc(item)); err != nil {
				return NewStoreError(opName, err)
			}
			s.migrateFunc(item, newName)
			if err := s.BaseStore.Put(newPrefix+s.keyFunc(item), item); err != nil {
				return NewStoreError(opName, err)
			}
		}

		return nil
	})
}
