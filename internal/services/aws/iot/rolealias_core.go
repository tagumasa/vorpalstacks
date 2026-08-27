package iot

import (
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Role alias Core
// ---------------------------------------------------------------------------

// CreateRoleAliasInput carries the fields for CreateRoleAlias. The credential
// duration defaults to one hour when DurationProvided is false.
type CreateRoleAliasInput struct {
	RoleAlias                 string
	RoleARN                   string
	CredentialDurationSeconds int64
	DurationProvided          bool
}

// UpdateRoleAliasInput carries the fields for UpdateRoleAlias. The role ARN
// and credential duration are applied only when provided.
type UpdateRoleAliasInput struct {
	RoleAlias                 string
	RoleARN                   string
	CredentialDurationSeconds int64
	DurationProvided          bool
}

// ListRoleAliasesResult is the transport-agnostic result of ListRoleAliases.
type ListRoleAliasesResult struct {
	RoleAliases []string
	NextMarker  string
}

// validateRoleAliasDuration enforces the CredentialDurationSeconds range
// documented for role aliases.
func validateRoleAliasDuration(seconds int64) error {
	if seconds < MinRoleAliasCredentialDuration || seconds > MaxRoleAliasCredentialDuration {
		return iotstore.ErrInvalidRequest
	}
	return nil
}

// createRoleAliasCore validates and persists a role alias.
func (s *IoTService) createRoleAliasCore(store iotstore.IotStoreInterface, in CreateRoleAliasInput) (*iotstore.RoleAlias, error) {
	if in.RoleAlias == "" || in.RoleARN == "" {
		return nil, iotstore.ErrMissingParam
	}

	credentialDuration := int64(3600)
	if in.DurationProvided {
		if err := validateRoleAliasDuration(in.CredentialDurationSeconds); err != nil {
			return nil, err
		}
		credentialDuration = in.CredentialDurationSeconds
	}

	ra := &iotstore.RoleAlias{
		RoleAlias:                 in.RoleAlias,
		RoleARN:                   in.RoleARN,
		CredentialDurationSeconds: credentialDuration,
		Owner:                     store.GetAccountID(),
		CreationDate:              time.Now().UTC(),
		LastModifiedDate:          time.Now().UTC(),
	}

	return store.CreateRoleAlias(ra)
}

// describeRoleAliasCore retrieves a role alias by name.
func (s *IoTService) describeRoleAliasCore(store iotstore.IotStoreInterface, roleAlias string) (*iotstore.RoleAlias, error) {
	if roleAlias == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.GetRoleAlias(roleAlias)
}

// updateRoleAliasCore applies the supplied fields to an existing role alias.
func (s *IoTService) updateRoleAliasCore(store iotstore.IotStoreInterface, in UpdateRoleAliasInput) (*iotstore.RoleAlias, error) {
	if in.RoleAlias == "" {
		return nil, iotstore.ErrMissingParam
	}

	var opts iotstore.RoleAliasUpdateOpts
	opts.RoleARN = in.RoleARN
	if in.DurationProvided {
		if err := validateRoleAliasDuration(in.CredentialDurationSeconds); err != nil {
			return nil, err
		}
		opts.DurationSeconds = in.CredentialDurationSeconds
	}

	return store.UpdateRoleAlias(in.RoleAlias, opts)
}

// deleteRoleAliasCore removes a role alias and its tags.
func (s *IoTService) deleteRoleAliasCore(store iotstore.IotStoreInterface, roleAlias string) error {
	if roleAlias == "" {
		return iotstore.ErrMissingParam
	}

	arn := iotstore.BuildRoleAliasARN(store.GetAccountID(), store.GetRegion(), roleAlias)
	_ = store.DeleteAllTags(arn)

	return store.DeleteRoleAlias(roleAlias)
}

// listRoleAliasesCore lists role alias names with pagination.
func (s *IoTService) listRoleAliasesCore(store iotstore.IotStoreInterface, marker string, maxItems int) (*ListRoleAliasesResult, error) {
	maxItems = listMaxItems(maxItems)
	aliases, err := store.ListRoleAliases(iotstoreListOpts(maxItems, marker))
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(aliases.Items))
	for _, a := range aliases.Items {
		names = append(names, a.RoleAlias)
	}
	return &ListRoleAliasesResult{
		RoleAliases: names,
		NextMarker:  aliases.NextMarker,
	}, nil
}
