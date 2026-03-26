package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// AliasStore manages KMS key aliases.
type AliasStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
}

func aliasBucketName(region string) string {
	return "kms_aliases-" + region
}

// NewAliasStore creates a new AliasStore.
func NewAliasStore(store storage.BasicStorage, accountId, region string) *AliasStore {
	return &AliasStore{
		BaseStore:  common.NewBaseStore(store.Bucket(aliasBucketName(region)), "kms"),
		arnBuilder: NewARNBuilder(accountId, region),
	}
}

func (s *AliasStore) normalizeAliasName(aliasName string) string {
	aliasName = s.arnBuilder.ParseAliasName(aliasName)
	if !strings.HasPrefix(aliasName, "alias/") {
		aliasName = "alias/" + aliasName
	}
	return aliasName
}

// Create creates a new alias for a KMS key.
func (s *AliasStore) Create(aliasName, targetKeyID string) (*Alias, error) {
	aliasName = s.normalizeAliasName(aliasName)

	if strings.HasPrefix(aliasName, "alias/aws/") {
		return nil, ErrInvalidAliasName
	}

	if s.Exists(aliasName) {
		return nil, ErrAliasAlreadyExists
	}

	now := time.Now()
	alias := &Alias{
		AliasName:       aliasName,
		AliasArn:        s.arnBuilder.AliasArn(aliasName),
		TargetKeyID:     targetKeyID,
		TargetKeyArn:    s.arnBuilder.KeyArn(targetKeyID),
		CreationDate:    now,
		LastUpdatedDate: now,
	}

	if err := s.save(alias); err != nil {
		return nil, err
	}

	return alias, nil
}

// Get retrieves an alias by name.
func (s *AliasStore) Get(aliasName string) (*Alias, error) {
	aliasName = s.normalizeAliasName(aliasName)

	var alias Alias
	if err := s.BaseStore.Get(aliasName, &alias); err != nil {
		return nil, err
	}
	return &alias, nil
}

func (s *AliasStore) save(alias *Alias) error {
	return s.BaseStore.Put(alias.AliasName, alias)
}

// Delete removes an alias.
func (s *AliasStore) Delete(aliasName string) error {
	aliasName = s.normalizeAliasName(aliasName)
	return s.BaseStore.Delete(aliasName)
}

// Exists checks whether an alias exists.
func (s *AliasStore) Exists(aliasName string) bool {
	aliasName = s.normalizeAliasName(aliasName)
	return s.BaseStore.Exists(aliasName)
}

// UpdateTarget updates the target key for an alias.
func (s *AliasStore) UpdateTarget(aliasName, targetKeyID string) error {
	alias, err := s.Get(aliasName)
	if err != nil {
		return err
	}

	alias.TargetKeyID = targetKeyID
	alias.TargetKeyArn = s.arnBuilder.KeyArn(targetKeyID)
	alias.LastUpdatedDate = time.Now()

	return s.save(alias)
}

// List returns a list of aliases with optional filtering by key ID.
func (s *AliasStore) List(marker string, maxItems int, keyID string) (*AliasListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var aliases []*Alias
	count := 0
	started := marker == ""
	var lastAliasName string
	hasMore := false

	err := s.ForEach(func(key string, value []byte) error {
		var alias Alias
		if err := json.Unmarshal(value, &alias); err != nil {
			return err
		}

		if !started {
			if alias.AliasName == marker {
				started = true
			}
			return nil
		}

		if keyID == "" || alias.TargetKeyID == keyID {
			if count < maxItems {
				aliases = append(aliases, &alias)
				count++
				lastAliasName = alias.AliasName
			} else {
				hasMore = true
			}
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_aliases", err)
	}

	result := &AliasListResult{
		Aliases:     aliases,
		IsTruncated: hasMore,
	}
	if result.IsTruncated {
		result.NextMarker = lastAliasName
	}

	return result, nil
}

// ListForKeyID returns aliases that target a specific key.
func (s *AliasStore) ListForKeyID(keyID string) ([]*Alias, error) {
	var aliases []*Alias

	err := s.ForEach(func(key string, value []byte) error {
		var alias Alias
		if err := json.Unmarshal(value, &alias); err != nil {
			return err
		}

		if alias.TargetKeyID == keyID {
			aliases = append(aliases, &alias)
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_aliases_for_key", err)
	}

	return aliases, nil
}

// ARNBuilder returns the ARN builder for this store.
func (s *AliasStore) ARNBuilder() *ARNBuilder {
	return s.arnBuilder
}
