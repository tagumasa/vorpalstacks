package kms

// alias_core.go carries the Core functions of the KMS alias family. Each
// Core takes the per-region store bundle plus the transport-agnostic
// operation members, performs validation, authorisation and persistence in
// the original failure-precedence order, and leaves wire parsing and
// response serialisation to the handlers.

import (
	"errors"
	"strings"

	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// mapAliasStoreError translates store-layer alias sentinel errors into the
// service-layer AWS error shapes.
func (s *KMSService) mapAliasStoreError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, kmsstore.ErrAliasAlreadyExists):
		return ErrAliasAlreadyExists
	case errors.Is(err, kmsstore.ErrAliasNotFound):
		return ErrAliasNotFound
	case errors.Is(err, kmsstore.ErrAliasNameReserved):
		return ErrInvalidAliasName
	case errors.Is(err, kmsstore.ErrInvalidKeyState):
		return ErrKeyPendingDeletion
	default:
		return err
	}
}

// createAliasCore is the single entry point for creating an alias. The
// validation ladder (alias name format, target key resolution,
// authorisation, target key state) runs before the store write, matching
// the original failure precedence.
func (s *KMSService) createAliasCore(stores *kmsStores, principal, aliasName, targetKeyID string) error {
	if err := validateAliasName(aliasName); err != nil {
		return err
	}

	if targetKeyID == "" {
		return ErrKeyNotFound
	}

	key, err := s.resolveKeyByKeyID(stores, targetKeyID)
	if err != nil {
		return err
	}
	if err := s.authorizeOperation(stores, principal, "CreateAlias", key.KeyID, nil); err != nil {
		return err
	}
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return ErrKeyPendingDeletion
	}
	if _, err := stores.aliases.Create(aliasName, key.KeyID); err != nil {
		return s.mapAliasStoreError(err)
	}
	return nil
}

// deleteAliasCore is the single entry point for deleting an alias.
func (s *KMSService) deleteAliasCore(stores *kmsStores, principal, aliasName string) error {
	if aliasName == "" {
		return ErrInvalidAliasName
	}

	if strings.HasPrefix(aliasName, "alias/aws/") {
		return ErrInvalidAliasName
	}

	alias, err := stores.aliases.Get(aliasName)
	if err != nil {
		return ErrAliasNotFound
	}

	if err := s.authorizeOperation(stores, principal, "DeleteAlias", alias.TargetKeyID, nil); err != nil {
		return err
	}

	if err := stores.aliases.Delete(aliasName); err != nil {
		return s.mapAliasStoreError(err)
	}
	return nil
}

// listAliasesCore is the single entry point for listing aliases,
// optionally narrowed to a target key.
func (s *KMSService) listAliasesCore(stores *kmsStores, marker string, maxItems int, keyID string) (*kmsstore.AliasListResult, error) {
	if err := validateMarkerLength(marker); err != nil {
		return nil, err
	}
	return stores.aliases.List(marker, maxItems, keyID)
}

// updateAliasCore is the single entry point for retargeting an alias.
func (s *KMSService) updateAliasCore(stores *kmsStores, principal, aliasName, targetKeyID string) error {
	if err := validateAliasName(aliasName); err != nil {
		return err
	}

	if targetKeyID == "" {
		return ErrKeyNotFound
	}

	key, err := s.resolveKeyByKeyID(stores, targetKeyID)
	if err != nil {
		return err
	}
	if err := s.authorizeOperation(stores, principal, "UpdateAlias", key.KeyID, nil); err != nil {
		return err
	}
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return ErrKeyPendingDeletion
	}
	if err := stores.aliases.UpdateTarget(aliasName, key.KeyID); err != nil {
		return s.mapAliasStoreError(err)
	}
	return nil
}
