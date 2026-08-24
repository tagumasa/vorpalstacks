package sfn

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// This file holds the version- and alias-operation Core path shared by
// the HTTP API and the admin console handler.

// PublishStateMachineVersionInput carries the parameters for
// PublishStateMachineVersion.
type PublishStateMachineVersionInput struct {
	StateMachineArn string
	Description     string
	RevisionId      string
}

// CreateStateMachineAliasInput carries every field that
// CreateStateMachineAlias needs.
type CreateStateMachineAliasInput struct {
	Name          string
	Description   string
	RoutingConfig []sfnstore.RoutingConfiguration
}

// UpdateStateMachineAliasInput carries every field that
// UpdateStateMachineAlias needs.
type UpdateStateMachineAliasInput struct {
	StateMachineAliasArn string
	Description          string
	DescriptionProvided  bool
	RoutingConfig        []sfnstore.RoutingConfiguration
	RoutingProvided      bool
}

// aliasNamePattern pins the CharacterRestrictedName shape's character
// class for alias names: letters, digits, hyphen, underscore and dot. The
// shape's Smithy pattern also requires at least one letter, hyphen,
// underscore or dot; Go's RE2 has no lookahead, so that half is a separate
// scan.
var aliasNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

// validateAliasName validates an alias name against the
// CharacterRestrictedName shape: 1-80 characters from the class, with at
// least one letter, hyphen, underscore or dot.
func validateAliasName(name string) error {
	if name == "" {
		return NewInvalidName("name is required")
	}
	if len(name) > sfnstore.MaxResourceNameLength {
		return NewInvalidName(fmt.Sprintf("name must be 1-80 characters, got %d", len(name)))
	}
	if !aliasNamePattern.MatchString(name) {
		return NewInvalidName("name must contain only letters, digits, hyphens, underscores and dots")
	}
	hasLetterOrSeparator := false
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == '-' || r == '.' {
			hasLetterOrSeparator = true
			break
		}
	}
	if !hasLetterOrSeparator {
		return NewInvalidName("name must contain at least one letter, underscore, hyphen or dot")
	}
	return nil
}

// baseStateMachineArn strips the numeric version qualifier from a state
// machine version ARN and returns the unqualified state machine ARN.
func baseStateMachineArn(versionArn string) (string, error) {
	partition, service, region, account, resource := svcarn.SplitARN(versionArn)
	if service != "states" {
		return "", fmt.Errorf("not a States ARN: %s", versionArn)
	}
	rest, ok := strings.CutPrefix(resource, "stateMachine:")
	if !ok {
		return "", fmt.Errorf("not a state machine ARN: %s", versionArn)
	}
	name, qualifier, qualified := strings.Cut(rest, ":")
	if name == "" || !qualified {
		return "", fmt.Errorf("not a version ARN: %s", versionArn)
	}
	if _, err := strconv.Atoi(qualifier); err != nil {
		return "", fmt.Errorf("not a version ARN: %s", versionArn)
	}
	return fmt.Sprintf("arn:%s:%s:%s:%s:stateMachine:%s", partition, service, region, account, name), nil
}

// publishStateMachineVersionCore is the single entry point for
// PublishStateMachineVersion, enforcing the VersionDescription length
// bound, the optimistic revisionId check (a mismatch is a conflict), the
// per-revision idempotency and the version quota.
func (s *StepFunctionService) publishStateMachineVersionCore(ctx context.Context, store *sfnstore.StepFunctionStore, in PublishStateMachineVersionInput) (map[string]interface{}, error) {
	if err := validateArnRequired(in.StateMachineArn, "stateMachineArn"); err != nil {
		return nil, err
	}
	if len(in.Description) > sfnstore.MaxVersionDescriptionLength {
		return nil, NewValidationException(fmt.Sprintf("description must be at most %d characters, got %d", sfnstore.MaxVersionDescriptionLength, len(in.Description)))
	}

	sm, err := store.GetStateMachine(ctx, in.StateMachineArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + in.StateMachineArn)
		}
		return nil, err
	}

	if in.RevisionId != "" && sm.RevisionId != in.RevisionId {
		return nil, NewConflictException("revisionId mismatch: expected " + sm.RevisionId + ", got " + in.RevisionId)
	}

	// Publishing is idempotent per revision: a retry that finds a version
	// of the current revision returns it, even when the quota is already
	// exhausted.
	if existing, err := store.FindVersionByRevision(in.StateMachineArn, sm.RevisionId); err == nil {
		return map[string]interface{}{
			"stateMachineVersionArn": existing.StateMachineVersionArn,
			"creationDate":           existing.CreationDate.Unix(),
		}, nil
	} else if !errors.Is(err, sfnstore.ErrStateMachineVersionNotFound) {
		return nil, err
	}

	if err := enforceVersionQuota(ctx, store, in.StateMachineArn); err != nil {
		return nil, err
	}

	version, err := store.PublishStateMachineVersion(ctx, in.StateMachineArn, in.Description)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"stateMachineVersionArn": version.StateMachineVersionArn,
		"creationDate":           version.CreationDate.Unix(),
	}, nil
}

// enforceVersionQuota rejects publishing past the documented version quota
// of one thousand versions per state machine.
func enforceVersionQuota(ctx context.Context, store *sfnstore.StepFunctionStore, smArn string) error {
	count, err := store.CountStateMachineVersions(smArn)
	if err != nil {
		return err
	}
	if count >= sfnstore.MaxVersionsPerStateMachine {
		return NewServiceQuotaExceededException(fmt.Sprintf("Maximum number of versions per state machine (%d) reached", sfnstore.MaxVersionsPerStateMachine))
	}
	return nil
}

// deleteStateMachineVersionCore is the single entry point for
// DeleteStateMachineVersion. The documented error set carries no
// not-found code, so an unknown well-formed version ARN maps to
// ValidationException; a version an alias still routes to maps to
// ConflictException.
func (s *StepFunctionService) deleteStateMachineVersionCore(ctx context.Context, store *sfnstore.StepFunctionStore, versionArn string) error {
	if err := validateArnRequired(versionArn, "stateMachineVersionArn"); err != nil {
		return err
	}
	if err := store.DeleteStateMachineVersion(ctx, versionArn); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineVersionNotFound) {
			return NewValidationException("State Machine Version Does not exist: " + versionArn)
		}
		if errors.Is(err, sfnstore.ErrStateMachineVersionInUse) {
			return NewConflictException("State Machine Version is referenced by an alias routing configuration: " + versionArn)
		}
		return err
	}
	return nil
}

// listStateMachineVersionsCore is the single entry point for
// ListStateMachineVersions.
func (s *StepFunctionService) listStateMachineVersionsCore(ctx context.Context, store *sfnstore.StepFunctionStore, smArn string, maxResults int32, nextToken string) (map[string]interface{}, error) {
	if err := validateArnRequired(smArn, "stateMachineArn"); err != nil {
		return nil, err
	}
	if err := validateMaxResults(maxResults, 0, sfnstore.MaxPageSize, "maxResults"); err != nil {
		return nil, err
	}

	result, err := store.ListStateMachineVersions(ctx, smArn, maxResults, nextToken)
	if err != nil {
		return nil, err
	}

	versions := make([]map[string]interface{}, len(result.Versions))
	for i, v := range result.Versions {
		versions[i] = map[string]interface{}{
			"stateMachineVersionArn": v.StateMachineVersionArn,
			"creationDate":           v.CreationDate.Unix(),
		}
	}

	response := map[string]interface{}{"stateMachineVersions": versions}
	if result.NextToken != "" {
		response["nextToken"] = result.NextToken
	}
	return response, nil
}

// validateAliasRoutingConfig validates a routing configuration against the
// RoutingConfigurationList shape: one or two entries whose weights lie in
// [0, 100] and sum to 100, naming existing versions of the given state
// machine.
func validateAliasRoutingConfig(ctx context.Context, store *sfnstore.StepFunctionStore, smArn string, rc []sfnstore.RoutingConfiguration) error {
	if len(rc) < 1 || len(rc) > sfnstore.MaxRoutingConfigEntries {
		return NewValidationException(fmt.Sprintf("routingConfiguration must contain 1-%d entries, got %d", sfnstore.MaxRoutingConfigEntries, len(rc)))
	}
	if err := validateRoutingConfiguration(rc); err != nil {
		return err
	}
	for _, entry := range rc {
		if base, err := baseStateMachineArn(entry.StateMachineVersionArn); err != nil || base != smArn {
			return NewValidationException("routingConfiguration entries must reference versions of the same state machine: " + entry.StateMachineVersionArn)
		}
		if _, err := store.GetStateMachineVersion(ctx, entry.StateMachineVersionArn); err != nil {
			return NewResourceNotFound("State Machine Version Does not exist: " + entry.StateMachineVersionArn)
		}
	}
	return nil
}

// routingConfigsEqual compares two alias routing configurations as sets:
// the one-or-two entry list has no meaningful order, so the idempotency
// key matches entries regardless of their position.
func routingConfigsEqual(a, b []sfnstore.RoutingConfiguration) bool {
	if len(a) != len(b) {
		return false
	}
	matched := make([]bool, len(b))
	for _, ea := range a {
		found := false
		for i, eb := range b {
			if !matched[i] && ea.StateMachineVersionArn == eb.StateMachineVersionArn && ea.Weight == eb.Weight {
				matched[i] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// createStateMachineAliasCore is the single entry point for
// CreateStateMachineAlias. The state machine is derived from the routing
// configuration's version ARNs (the operation carries no stateMachineArn
// member); every routed version must belong to that state machine. The
// operation is idempotent: a retry with the same description, name and
// routing configuration returns the existing alias, while the same name
// with a different configuration is a conflict. The alias quota is one
// hundred aliases per state machine.
func (s *StepFunctionService) createStateMachineAliasCore(ctx context.Context, store *sfnstore.StepFunctionStore, in CreateStateMachineAliasInput) (map[string]interface{}, error) {
	if err := validateAliasName(in.Name); err != nil {
		return nil, err
	}
	if len(in.Description) > sfnstore.MaxVersionDescriptionLength {
		return nil, NewValidationException(fmt.Sprintf("description must be at most %d characters, got %d", sfnstore.MaxVersionDescriptionLength, len(in.Description)))
	}
	if len(in.RoutingConfig) == 0 {
		return nil, NewValidationException("routingConfiguration is required")
	}

	smArn, err := baseStateMachineArn(in.RoutingConfig[0].StateMachineVersionArn)
	if err != nil {
		return nil, NewInvalidArnException("routingConfiguration entries must be state machine version ARNs")
	}

	sm, err := store.GetStateMachine(ctx, smArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + smArn)
		}
		return nil, err
	}
	if sm.Status == "DELETING" {
		return nil, NewStateMachineDeleting("State Machine is being deleted: " + smArn)
	}

	if err := validateAliasRoutingConfig(ctx, store, smArn, in.RoutingConfig); err != nil {
		return nil, err
	}

	if existing, err := store.GetStateMachineAliasByName(ctx, smArn, in.Name); err == nil {
		if existing.Description == in.Description && routingConfigsEqual(existing.RoutingConfiguration, in.RoutingConfig) {
			return map[string]interface{}{
				"stateMachineAliasArn": existing.StateMachineAliasArn,
				"creationDate":         existing.CreationDate.Unix(),
			}, nil
		}
		return nil, NewConflictException(fmt.Sprintf("State Machine Alias already exists: %s", in.Name))
	} else if !errors.Is(err, sfnstore.ErrStateMachineAliasNotFound) {
		return nil, err
	}

	aliasCount, err := store.CountStateMachineAliases(smArn)
	if err != nil {
		return nil, err
	}
	if aliasCount >= sfnstore.MaxAliasesPerStateMachine {
		return nil, NewServiceQuotaExceededException(fmt.Sprintf("Maximum number of aliases per state machine (%d) reached", sfnstore.MaxAliasesPerStateMachine))
	}

	alias := &sfnstore.StateMachineAlias{
		StateMachineArn:      smArn,
		Name:                 in.Name,
		Description:          in.Description,
		RoutingConfiguration: in.RoutingConfig,
	}
	if err := store.CreateStateMachineAlias(ctx, alias); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAliasAlreadyExists) {
			return nil, NewConflictException(fmt.Sprintf("State Machine Alias already exists: %s", in.Name))
		}
		return nil, err
	}

	return map[string]interface{}{
		"stateMachineAliasArn": alias.StateMachineAliasArn,
		"creationDate":         alias.CreationDate.Unix(),
	}, nil
}

// describeStateMachineAliasCore is the single entry point for
// DescribeStateMachineAlias.
func (s *StepFunctionService) describeStateMachineAliasCore(ctx context.Context, store *sfnstore.StepFunctionStore, aliasArn string) (map[string]interface{}, error) {
	if err := validateArnRequired(aliasArn, "stateMachineAliasArn"); err != nil {
		return nil, err
	}

	alias, err := store.GetStateMachineAlias(ctx, aliasArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAliasNotFound) {
			return nil, NewResourceNotFound("State Machine Alias Does not exist: " + aliasArn)
		}
		return nil, err
	}

	resp := map[string]interface{}{
		"stateMachineAliasArn": alias.StateMachineAliasArn,
		"name":                 alias.Name,
		"creationDate":         alias.CreationDate.Unix(),
		"updateDate":           alias.UpdateDate.Unix(),
	}
	if alias.Description != "" {
		resp["description"] = alias.Description
	}
	if len(alias.RoutingConfiguration) > 0 {
		resp["routingConfiguration"] = formatRoutingConfiguration(alias.RoutingConfiguration)
	}
	return resp, nil
}

// updateStateMachineAliasCore is the single entry point for
// UpdateStateMachineAlias. An update must carry at least one of the
// description or routingConfiguration parameters.
func (s *StepFunctionService) updateStateMachineAliasCore(ctx context.Context, store *sfnstore.StepFunctionStore, in UpdateStateMachineAliasInput) (map[string]interface{}, error) {
	if err := validateArnRequired(in.StateMachineAliasArn, "stateMachineAliasArn"); err != nil {
		return nil, err
	}
	if !in.DescriptionProvided && !in.RoutingProvided {
		return nil, NewValidationException("Request must include at least one of description or routingConfiguration")
	}
	if len(in.Description) > sfnstore.MaxVersionDescriptionLength {
		return nil, NewValidationException(fmt.Sprintf("description must be at most %d characters, got %d", sfnstore.MaxVersionDescriptionLength, len(in.Description)))
	}

	alias, err := store.GetStateMachineAlias(ctx, in.StateMachineAliasArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAliasNotFound) {
			return nil, NewResourceNotFound("State Machine Alias Does not exist: " + in.StateMachineAliasArn)
		}
		return nil, err
	}

	if in.RoutingProvided {
		if err := validateAliasRoutingConfig(ctx, store, alias.StateMachineArn, in.RoutingConfig); err != nil {
			return nil, err
		}
	}

	if in.DescriptionProvided && in.Description != "" {
		alias.Description = in.Description
	}
	if in.RoutingProvided && len(in.RoutingConfig) > 0 {
		alias.RoutingConfiguration = in.RoutingConfig
	}

	if err := store.UpdateStateMachineAlias(ctx, alias); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAliasNotFound) {
			return nil, NewResourceNotFound("State Machine Alias Does not exist: " + in.StateMachineAliasArn)
		}
		return nil, err
	}

	return map[string]interface{}{"updateDate": alias.UpdateDate.Unix()}, nil
}

// deleteStateMachineAliasCore is the single entry point for
// DeleteStateMachineAlias.
func (s *StepFunctionService) deleteStateMachineAliasCore(ctx context.Context, store *sfnstore.StepFunctionStore, aliasArn string) error {
	if err := validateArnRequired(aliasArn, "stateMachineAliasArn"); err != nil {
		return err
	}
	if err := store.DeleteStateMachineAlias(ctx, aliasArn); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAliasNotFound) {
			return NewResourceNotFound("State Machine Alias Does not exist: " + aliasArn)
		}
		return err
	}
	return nil
}

// listStateMachineAliasesCore is the single entry point for
// ListStateMachineAliases; the state machine must exist.
func (s *StepFunctionService) listStateMachineAliasesCore(ctx context.Context, store *sfnstore.StepFunctionStore, smArn string, maxResults int32, nextToken string) (map[string]interface{}, error) {
	if err := validateArnRequired(smArn, "stateMachineArn"); err != nil {
		return nil, err
	}
	if err := validateMaxResults(maxResults, 0, sfnstore.MaxPageSize, "maxResults"); err != nil {
		return nil, err
	}
	if _, err := store.GetStateMachine(ctx, smArn); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + smArn)
		}
		return nil, err
	}

	result, err := store.ListStateMachineAliases(ctx, smArn, maxResults, nextToken)
	if err != nil {
		return nil, err
	}

	aliases := make([]map[string]interface{}, len(result.Aliases))
	for i, a := range result.Aliases {
		aliases[i] = map[string]interface{}{
			"stateMachineAliasArn": a.StateMachineAliasArn,
			"creationDate":         a.CreationDate.Unix(),
		}
	}

	response := map[string]interface{}{"stateMachineAliases": aliases}
	if result.NextToken != "" {
		response["nextToken"] = result.NextToken
	}
	return response, nil
}
