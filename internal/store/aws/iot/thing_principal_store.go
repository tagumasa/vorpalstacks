package iot

import (
	"context"
	"strings"
	"vorpalstacks/internal/core/storage"
)

func (s *IotStore) AttachThingPrincipal(thingName, principal string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.thingsBase.Exists(thingName) {
		return ErrThingNotFound
	}
	tk := thingPrincipalKey(thingName, principal)
	pk := principalThingKey(principal, thingName)
	tpBucket := bucketThingPrincipal + s.rs
	ptBucket := bucketPrincipalThing + s.rs
	return s.ts.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(tpBucket).Put([]byte(tk), []byte("1")); err != nil {
			return err
		}
		return txn.Bucket(ptBucket).Put([]byte(pk), []byte("1"))
	})
}

func (s *IotStore) DetachThingPrincipal(thingName, principal string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.thingsBase.Exists(thingName) {
		return ErrThingNotFound
	}
	tk := thingPrincipalKey(thingName, principal)
	pk := principalThingKey(principal, thingName)
	tpBucket := bucketThingPrincipal + s.rs
	ptBucket := bucketPrincipalThing + s.rs
	return s.ts.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(tpBucket).Delete([]byte(tk)); err != nil {
			return err
		}
		return txn.Bucket(ptBucket).Delete([]byte(pk))
	})
}

func (s *IotStore) ListPrincipalsForThing(thingName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var principals []string
	err := s.thingPrincipalBase.ScanPrefix(thingName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			principals = append(principals, parts[1])
		}
		return nil
	})
	return principals, err
}

func (s *IotStore) ListThingsForPrincipal(principal string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var things []string
	err := s.principalThingBase.ScanPrefix(principal+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			things = append(things, parts[1])
		}
		return nil
	})
	return things, err
}

func thingPrincipalKey(thingName, principal string) string {
	return thingName + "\x00" + principal
}

func principalThingKey(principal, thingName string) string {
	return principal + "\x00" + thingName
}
