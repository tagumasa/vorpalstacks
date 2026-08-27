package iot

import (
	"context"
	"strings"

	"vorpalstacks/internal/core/storage"
)

// principalExclusiveMarker is the stored value of a thing↔principal pair
// declared EXCLUSIVE_THING by a provisioning template; every other value
// (including the plain "1" written by the untyped API path) denotes
// NON_EXCLUSIVE_THING.
const principalExclusiveMarker = "EXCLUSIVE_THING"

// AttachThingPrincipal records a NON_EXCLUSIVE_THING attachment between the
// thing and the principal.
func (s *IotStore) AttachThingPrincipal(thingName, principal string) error {
	return s.attachThingPrincipal(thingName, principal, false)
}

// AttachThingPrincipalExclusive records an EXCLUSIVE_THING attachment: the
// principal may attach to this thing only, and a principal already marked
// exclusive rejects every further attachment regardless of the requested
// type.
func (s *IotStore) AttachThingPrincipalExclusive(thingName, principal string) error {
	return s.attachThingPrincipal(thingName, principal, true)
}

func (s *IotStore) attachThingPrincipal(thingName, principal string, exclusive bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.thingsBase.Exists(thingName) {
		return ErrThingNotFound
	}
	// Exclusivity: a principal marked exclusive cannot attach to a second
	// thing, and an exclusive request fails when the principal is already
	// attached anywhere else.
	var conflict bool
	err := s.principalThingBase.ScanPrefix(principal+"\x00", func(key string, val []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 && parts[1] != thingName {
			if exclusive || string(val) == principalExclusiveMarker {
				conflict = true
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if conflict {
		return ErrInvalidRequest
	}
	value := "1"
	if exclusive {
		value = principalExclusiveMarker
	}
	tk := thingPrincipalKey(thingName, principal)
	pk := principalThingKey(principal, thingName)
	tpBucket := bucketThingPrincipal + s.rs
	ptBucket := bucketPrincipalThing + s.rs
	return s.ts.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(tpBucket).Put([]byte(tk), []byte(value)); err != nil {
			return err
		}
		return txn.Bucket(ptBucket).Put([]byte(pk), []byte(value))
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
