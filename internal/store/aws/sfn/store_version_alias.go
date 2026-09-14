// This file holds the state-machine version and alias surfaces: version
// publication with the recovered monotonic counter, and alias CRUD with
// routing configuration.

// Package stepfunction provides Step Functions storage functionality for vorpalstacks.
package sfn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"vorpalstacks/internal/store/aws/common"
)

func (s *StepFunctionStore) buildVersionARN(smArn string, version int64) string {
	return smArn + fmt.Sprintf(":%d", version)
}

func (s *StepFunctionStore) buildAliasARN(smArn, aliasName string) string {
	// The alias ARN extends the state machine ARN with the alias name, so
	// alias names are namespaced per state machine rather than per account.
	return smArn + ":" + aliasName
}

func (s *StepFunctionStore) nextVersionNumber(smArn string) int64 {
	s.versionCountersMu.Lock()
	defer s.versionCountersMu.Unlock()
	s.versionCounters[smArn]++
	return s.versionCounters[smArn]
}

func (s *StepFunctionStore) recoverVersionCounter(smArn string) {
	versions, err := common.ListMatching[StateMachineVersion](s.versionsStore, smArn+":", nil)
	if err != nil {
		return
	}
	var maxVersion int64
	for _, v := range versions {
		if v.Version > maxVersion {
			maxVersion = v.Version
		}
	}
	s.versionCountersMu.Lock()
	if maxVersion > s.versionCounters[smArn] {
		s.versionCounters[smArn] = maxVersion
	}
	s.versionCountersMu.Unlock()
}

// PublishStateMachineVersion publishes the state machine's current
// revision as a version. Publishing is idempotent per revision: when a
// version of the current revision already exists it is returned instead of
// publishing a duplicate. The idempotency check and the write run under
// createMu — the same check-then-write serialisation CreateExecution and
// CreateStateMachineAlias use — so concurrent publishes of one revision
// cannot both miss the find and mint duplicate version records.
func (s *StepFunctionStore) PublishStateMachineVersion(ctx context.Context, smArn string, description string) (*StateMachineVersion, error) {
	sm, err := s.GetStateMachine(ctx, smArn)
	if err != nil {
		return nil, err
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	if existing, err := s.FindVersionByRevision(smArn, sm.RevisionId); err == nil {
		return existing, nil
	} else if !errors.Is(err, ErrStateMachineVersionNotFound) {
		return nil, err
	}

	s.versionCountersMu.Lock()
	if _, exists := s.versionCounters[smArn]; !exists {
		s.versionCountersMu.Unlock()
		s.recoverVersionCounter(smArn)
	} else {
		s.versionCountersMu.Unlock()
	}

	version := s.nextVersionNumber(smArn)
	versionArn := s.buildVersionARN(smArn, version)

	v := &StateMachineVersion{
		StateMachineVersionArn: versionArn,
		StateMachineArn:        smArn,
		Version:                version,
		Description:            description,
		CreationDate:           time.Now().UTC(),
		Definition:             sm.Definition,
		RevisionId:             sm.RevisionId,
	}

	if err := s.versionsStore.Put(versionArn, v); err != nil {
		return nil, err
	}

	return v, nil
}

// FindVersionByRevision returns the version published for the given state
// machine revision, or ErrStateMachineVersionNotFound when the revision has
// never been published. When several records claim the same revision the
// highest version number wins, keeping the result deterministic.
func (s *StepFunctionStore) FindVersionByRevision(smArn, revisionId string) (*StateMachineVersion, error) {
	versions, err := common.ListMatching[StateMachineVersion](s.versionsStore, smArn+":", func(v *StateMachineVersion) bool {
		return v.RevisionId == revisionId
	})
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, ErrStateMachineVersionNotFound
	}
	newest := versions[0]
	for _, v := range versions[1:] {
		if v.Version > newest.Version {
			newest = v
		}
	}
	return newest, nil
}

// CountStateMachineVersions returns the number of published versions of a
// state machine; the quota is one thousand versions per state machine.
func (s *StepFunctionStore) CountStateMachineVersions(smArn string) (int, error) {
	versions, err := common.ListMatching[StateMachineVersion](s.versionsStore, smArn+":", nil)
	if err != nil {
		return 0, err
	}
	return len(versions), nil
}

// GetStateMachineVersion retrieves a state machine version by its ARN. A
// missing record is the not-found sentinel; any other storage fault surfaces
// as itself.
func (s *StepFunctionStore) GetStateMachineVersion(ctx context.Context, arn string) (*StateMachineVersion, error) {
	var v StateMachineVersion
	if err := s.versionsStore.Get(arn, &v); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrStateMachineVersionNotFound
		}
		return nil, err
	}
	return &v, nil
}

// GetStateMachineVersionByNumber retrieves a version of a state machine by
// its sequential version number.
func (s *StepFunctionStore) GetStateMachineVersionByNumber(ctx context.Context, smArn string, version int64) (*StateMachineVersion, error) {
	return s.GetStateMachineVersion(ctx, s.buildVersionARN(smArn, version))
}

// DeleteStateMachineVersion removes a state machine version from the store.
// A version an alias routing configuration still points at cannot be
// deleted; the caller maps the sentinel to ConflictException.
func (s *StepFunctionStore) DeleteStateMachineVersion(ctx context.Context, arn string) error {
	if !s.versionsStore.Exists(arn) {
		return ErrStateMachineVersionNotFound
	}
	aliases, err := common.ListMatching[StateMachineAlias](s.aliasesStore, "", func(a *StateMachineAlias) bool {
		for _, rc := range a.RoutingConfiguration {
			if rc.StateMachineVersionArn == arn {
				return true
			}
		}
		return false
	})
	if err != nil {
		return err
	}
	if len(aliases) > 0 {
		return ErrStateMachineVersionInUse
	}
	return s.versionsStore.Delete(arn)
}

// ListAllStateMachineVersions returns every published version of a state
// machine without pagination. The ListStateMachineVersions contract orders
// results by creation time ("The results are sorted in descending order
// of the version creation time"), so the service layer fetches the full
// set before sorting and paging it.
func (s *StepFunctionStore) ListAllStateMachineVersions(smArn string) ([]*StateMachineVersion, error) {
	return common.ListMatching[StateMachineVersion](s.versionsStore, smArn+":", nil)
}

// CreateStateMachineAlias creates a new alias for a state machine.
func (s *StepFunctionStore) CreateStateMachineAlias(ctx context.Context, alias *StateMachineAlias) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	if alias.Name == "" {
		return ErrInvalidARN
	}

	aliasArn := s.buildAliasARN(alias.StateMachineArn, alias.Name)
	storageKey := aliasArn
	if s.aliasesStore.Exists(storageKey) {
		return ErrStateMachineAliasAlreadyExists
	}

	now := time.Now().UTC()
	alias.StateMachineAliasArn = aliasArn
	alias.CreationDate = now
	alias.UpdateDate = now

	return s.aliasesStore.Put(storageKey, alias)
}

// GetStateMachineAlias retrieves a state machine alias by its ARN. A missing
// record is the not-found sentinel; any other storage fault surfaces as
// itself.
func (s *StepFunctionStore) GetStateMachineAlias(ctx context.Context, arn string) (*StateMachineAlias, error) {
	var alias StateMachineAlias
	if err := s.aliasesStore.Get(arn, &alias); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrStateMachineAliasNotFound
		}
		return nil, err
	}
	return &alias, nil
}

// GetStateMachineAliasByName retrieves a state machine alias by its state
// machine ARN and alias name.
func (s *StepFunctionStore) GetStateMachineAliasByName(ctx context.Context, smArn, name string) (*StateMachineAlias, error) {
	return s.GetStateMachineAlias(ctx, s.buildAliasARN(smArn, name))
}

// UpdateStateMachineAlias updates an existing state machine alias.
func (s *StepFunctionStore) UpdateStateMachineAlias(ctx context.Context, alias *StateMachineAlias) error {
	if !s.aliasesStore.Exists(alias.StateMachineAliasArn) {
		return ErrStateMachineAliasNotFound
	}
	alias.UpdateDate = time.Now().UTC()
	return s.aliasesStore.Put(alias.StateMachineAliasArn, alias)
}

// DeleteStateMachineAlias removes a state machine alias from the store.
func (s *StepFunctionStore) DeleteStateMachineAlias(ctx context.Context, arn string) error {
	if !s.aliasesStore.Exists(arn) {
		return ErrStateMachineAliasNotFound
	}
	return s.aliasesStore.Delete(arn)
}

// ListAllStateMachineAliases returns every alias of a state machine
// without pagination. The ListAliases contract orders results by creation
// time ("Results are sorted by time, with the most recently created
// aliases listed first"), so the service layer fetches the full set
// before sorting and paging it.
func (s *StepFunctionStore) ListAllStateMachineAliases(smArn string) ([]*StateMachineAlias, error) {
	return common.ListMatching[StateMachineAlias](s.aliasesStore, "", func(a *StateMachineAlias) bool {
		return a.StateMachineArn == smArn
	})
}

// CountStateMachineAliases returns the number of aliases defined for a
// state machine; the quota is one hundred aliases per state machine.
func (s *StepFunctionStore) CountStateMachineAliases(smArn string) (int, error) {
	aliases, err := common.ListMatching[StateMachineAlias](s.aliasesStore, "", func(a *StateMachineAlias) bool {
		return a.StateMachineArn == smArn
	})
	if err != nil {
		return 0, err
	}
	return len(aliases), nil
}
