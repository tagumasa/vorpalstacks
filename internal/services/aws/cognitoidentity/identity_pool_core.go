package cognitoidentity

import (
	"errors"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// CreateIdentityPoolInput carries every field that CreateIdentityPool needs,
// in a format independent of the wire protocol (HTTP Query/JSON vs gRPC-Web).
type CreateIdentityPoolInput struct {
	IdentityPoolName               string
	AllowUnauthenticatedIdentities bool
	AllowClassicFlow               bool
	AllowClassicFlowProvided       bool
	CognitoIdentityProviders       []ProviderOut
	DeveloperProviderName          string
	SupportedLoginProviders        map[string]string
	OpenIdConnectProviderARNs      []string
	SamlProviderARNs               []string
	Tags                           map[string]string
	TagsProvided                   bool
	Region                         string
}

// ListIdentityPoolsInput carries every field that ListIdentityPools needs.
type ListIdentityPoolsInput struct {
	MaxResults         int
	MaxResultsProvided bool
	NextToken          string
}

// IdentityPoolOut is the transport-agnostic representation of an IdentityPool,
// used by both the HTTP API and admin gRPC handler to format responses.
type IdentityPoolOut struct {
	ID                             string
	Name                           string
	AllowUnauthenticatedIdentities bool
	AllowClassicFlow               bool
	CognitoIdentityProviders       []ProviderOut
	DeveloperProviderName          string
	SupportedLoginProviders        map[string]string
	OpenIdConnectProviderARNs      []string
	SamlProviderARNs               []string
	Tags                           map[string]string
	Arn                            string
}

// ProviderOut is the transport-agnostic representation of a Cognito identity
// provider configuration entry.
type ProviderOut struct {
	ProviderName         string
	ClientID             string
	ServerSideTokenCheck bool
}

// IdentityPoolShortOut is the short description returned by ListIdentityPools.
type IdentityPoolShortOut struct {
	ID   string
	Name string
}

// RoleMappingInput is the transport-agnostic representation of a RoleMapping,
// used by validators and service-layer Core functions without depending on
// the store package.
type RoleMappingInput struct {
	Type                    string
	AmbiguousRoleResolution string
	RulesConfiguration      *RulesConfigInput
}

// RulesConfigInput is the transport-agnostic representation of
// RulesConfiguration.
type RulesConfigInput struct {
	Rules []MappingRuleInput
}

// MappingRuleInput is the transport-agnostic representation of a MappingRule.
type MappingRuleInput struct {
	Claim     string
	MatchType string
	Value     string
	RoleARN   string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createIdentityPoolCore is the single entry point for identity pool creation
// shared by the HTTP API and the admin gRPC handler.
func (s *CognitoIdentityService) createIdentityPoolCore(store cognitoidentitystore.CognitoIdentityStoreInterface, in CreateIdentityPoolInput) (*IdentityPoolOut, error) {
	if !validateIdentityPoolName(in.IdentityPoolName) {
		return nil, ErrInvalidParameter
	}

	if in.DeveloperProviderName != "" {
		if !validateDeveloperProviderName(in.DeveloperProviderName) {
			return nil, ErrInvalidParameter
		}
	}

	for _, p := range in.CognitoIdentityProviders {
		if p.ProviderName != "" {
			if !validateProviderName(p.ProviderName) {
				return nil, ErrInvalidParameter
			}
		}
		if p.ClientID != "" {
			if !validateProviderClientId(p.ClientID) {
				return nil, ErrInvalidParameter
			}
		}
	}

	// SupportedLoginProviders carries the model's IdentityProviders map
	// with its @length max of 10.
	if !validateMapSize(len(in.SupportedLoginProviders), maxLoginProviders) {
		return nil, ErrInvalidParameter
	}
	if in.TagsProvided && len(in.Tags) > 0 {
		if !validateTagKeys(in.Tags) || !validateTagValues(in.Tags) ||
			!validateMapSize(len(in.Tags), tagutil.MaxTagsPerResource) {
			return nil, ErrInvalidParameter
		}
	}
	for _, arn := range in.OpenIdConnectProviderARNs {
		if !validateRoleARN(arn) {
			return nil, ErrInvalidParameter
		}
	}
	for _, arn := range in.SamlProviderARNs {
		if !validateRoleARN(arn) {
			return nil, ErrInvalidParameter
		}
	}

	// The per-account quota is enforced by the store atomically with the
	// creation itself, so no separate pre-check is needed here.

	pool := cognitoidentitystore.NewIdentityPool(in.IdentityPoolName, in.AllowUnauthenticatedIdentities, in.Region)

	if len(in.CognitoIdentityProviders) > 0 {
		pool.CognitoIdentityProviders = providerOutsToStore(in.CognitoIdentityProviders)
	}
	if in.DeveloperProviderName != "" {
		pool.DeveloperProviderName = in.DeveloperProviderName
	}
	if len(in.SupportedLoginProviders) > 0 {
		pool.SupportedLoginProviders = in.SupportedLoginProviders
	}
	if len(in.OpenIdConnectProviderARNs) > 0 {
		pool.OpenIdConnectProviderARNs = in.OpenIdConnectProviderARNs
	}
	if len(in.SamlProviderARNs) > 0 {
		pool.SamlProviderARNs = in.SamlProviderARNs
	}
	if in.AllowClassicFlowProvided {
		pool.AllowClassicFlow = in.AllowClassicFlow
	}

	created, err := store.CreateIdentityPool(pool)
	if err != nil {
		if errors.Is(err, cognitoidentitystore.ErrIdentityPoolAlreadyExists) {
			return nil, ErrResourceInUse
		}
		if errors.Is(err, cognitoidentitystore.ErrTooManyIdentityPools) {
			return nil, ErrLimitExceeded
		}
		return nil, ErrInternalError
	}

	if in.TagsProvided && len(in.Tags) > 0 {
		if err := store.Tag(created.Arn, in.Tags); err != nil {
			logs.Error("Failed to tag identity pool, attempting cleanup", logs.String("poolId", created.ID), logs.Err(err))
			if delErr := store.DeleteIdentityPool(created.ID); delErr != nil {
				logs.Error("Failed to cleanup identity pool after tag failure", logs.String("poolId", created.ID), logs.Err(delErr))
				// The compensating delete failed: an untagged pool remains,
				// and the caller must see it rather than retry the create
				// into ResourceInUse.
				return nil, fmt.Errorf("%w: identity pool %s was created but the tag write and its rollback both failed; the pool remains", ErrInternalError, created.ID)
			}
			return nil, ErrInternalError
		}
		created.Tags = in.Tags
	}

	return poolToOut(created), nil
}

// listIdentityPoolsShortCore returns only the short description (ID + Name),
// matching the Smithy IdentityPoolShortDescription shape.
func (s *CognitoIdentityService) listIdentityPoolsShortCore(store cognitoidentitystore.CognitoIdentityStoreInterface, in ListIdentityPoolsInput) ([]IdentityPoolShortOut, string, error) {
	if !in.MaxResultsProvided {
		return nil, "", ErrInvalidParameter
	}
	if !validateQueryLimit(in.MaxResults) {
		return nil, "", ErrInvalidParameter
	}
	if !validatePaginationKey(in.NextToken) {
		return nil, "", ErrInvalidParameter
	}

	opts := storecommon.ListOptions{
		MaxItems: in.MaxResults,
		Marker:   in.NextToken,
	}

	result, err := store.ListIdentityPools(opts)
	if err != nil {
		return nil, "", ErrInternalError
	}

	items := make([]IdentityPoolShortOut, 0, len(result.Items))
	for _, pool := range result.Items {
		items = append(items, IdentityPoolShortOut{
			ID:   pool.ID,
			Name: pool.Name,
		})
	}

	return items, result.NextMarker, nil
}

// deleteIdentityPoolCore is the single entry point for identity pool deletion.
func (s *CognitoIdentityService) deleteIdentityPoolCore(store cognitoidentitystore.CognitoIdentityStoreInterface, poolID string) error {
	if !validateIdentityPoolId(poolID) {
		return ErrInvalidParameter
	}
	if err := store.DeleteIdentityPool(poolID); err != nil {
		if errors.Is(err, cognitoidentitystore.ErrIdentityPoolNotFound) {
			return ErrResourceNotFound
		}
		return ErrInternalError
	}
	return nil
}

// describeIdentityPoolCore is the single entry point for identity pool retrieval.
func (s *CognitoIdentityService) describeIdentityPoolCore(store cognitoidentitystore.CognitoIdentityStoreInterface, poolID string) (*IdentityPoolOut, error) {
	if !validateIdentityPoolId(poolID) {
		return nil, ErrInvalidParameter
	}
	pool, err := store.GetIdentityPool(poolID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}
	tags, err := store.List(pool.Arn)
	if err != nil {
		logs.Error("Failed to list identity pool tags", logs.String("poolId", pool.ID), logs.Err(err))
		return nil, ErrInternalError
	}
	if len(tags) > 0 {
		pool.Tags = tags
	}
	return poolToOut(pool), nil
}

// UpdateIdentityPoolInput carries every field that UpdateIdentityPool needs.
// The Raw members carry the untyped wire values with their presence flags so
// the Core can type-check and validate each member at the exact position the
// wire contract requires.
type UpdateIdentityPoolInput struct {
	IdentityPoolID string
	PoolName       string

	AllowUnauthProvided  bool
	AllowUnauthRaw       interface{}
	AllowClassicProvided bool
	AllowClassicRaw      interface{}

	DeveloperProviderName     string
	ProvidersRaw              interface{}
	SupportedLoginProviders   map[string]string
	OpenIdConnectProviderARNs []string
	SamlProviderARNs          []string

	TagsProvided bool
	Tags         map[string]string
}

// updateIdentityPoolCore is the single entry point for identity pool updates.
// The whole read-validate-merge-write cycle runs inside the store's pool-lock
// callback, so two concurrent updates serialise instead of losing the earlier
// caller's field changes, and a pool deleted mid-flight maps to
// ResourceNotFoundException instead of InternalErrorException.
func (s *CognitoIdentityService) updateIdentityPoolCore(reqCtx *request.RequestContext, in UpdateIdentityPoolInput) (*IdentityPoolOut, error) {
	if !validateIdentityPoolId(in.IdentityPoolID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var updated *cognitoidentitystore.IdentityPool
	err = store.UpdateIdentityPoolFunc(in.IdentityPoolID, func(pool *cognitoidentitystore.IdentityPool) error {
		// IdentityPoolName is @required in the Smithy IdentityPool shape.
		if in.PoolName == "" {
			return ErrInvalidParameter
		}
		if !validateIdentityPoolName(in.PoolName) {
			return ErrInvalidParameter
		}
		pool.Name = in.PoolName

		// AllowUnauthenticatedIdentities is @required in the Smithy shape.
		if !in.AllowUnauthProvided {
			return ErrInvalidParameter
		}
		if b, ok := in.AllowUnauthRaw.(bool); ok {
			pool.AllowUnauthenticatedIdentities = b
		} else {
			return ErrInvalidParameter
		}

		if in.AllowClassicProvided {
			b, ok := in.AllowClassicRaw.(bool)
			if !ok {
				return ErrInvalidParameter
			}
			pool.AllowClassicFlow = b
		} else {
			pool.AllowClassicFlow = false
		}
		// The model's UpdateIdentityPool contract resets every omitted member
		// to its default value ("If you don't provide a value for a
		// parameter, Amazon Cognito sets it to its default value."), so an
		// absent member and an explicitly empty member both clear the stored
		// value; only the @required members above must be present.
		if in.DeveloperProviderName != "" {
			if !validateDeveloperProviderName(in.DeveloperProviderName) {
				return ErrInvalidParameter
			}
		}
		pool.DeveloperProviderName = in.DeveloperProviderName
		providers, err := parseCognitoIdentityProviders(in.ProvidersRaw)
		if err != nil {
			return err
		}
		pool.CognitoIdentityProviders = providerOutsToStore(providers)
		// SupportedLoginProviders carries the model's IdentityProviders
		// map with its @length max of 10.
		if !validateMapSize(len(in.SupportedLoginProviders), maxLoginProviders) {
			return ErrInvalidParameter
		}
		if in.SupportedLoginProviders == nil {
			pool.SupportedLoginProviders = make(map[string]string)
		} else {
			pool.SupportedLoginProviders = in.SupportedLoginProviders
		}
		for _, arn := range in.OpenIdConnectProviderARNs {
			if !validateRoleARN(arn) {
				return ErrInvalidParameter
			}
		}
		pool.OpenIdConnectProviderARNs = in.OpenIdConnectProviderARNs
		for _, arn := range in.SamlProviderARNs {
			if !validateRoleARN(arn) {
				return ErrInvalidParameter
			}
		}
		pool.SamlProviderARNs = in.SamlProviderARNs

		if in.TagsProvided {
			if !validateTagKeys(in.Tags) || !validateTagValues(in.Tags) ||
				!validateMapSize(len(in.Tags), tagutil.MaxTagsPerResource) {
				return ErrInvalidParameter
			}
			pool.Tags = in.Tags
		} else {
			tags, err := store.List(pool.Arn)
			if err != nil {
				return err
			}
			pool.Tags = tags
		}

		updated = pool
		return nil
	})
	if err != nil {
		if errors.Is(err, cognitoidentitystore.ErrIdentityPoolNotFound) {
			return nil, ErrResourceNotFound
		}
		// Validation failures raised inside the callback are service errors
		// and pass through unchanged; anything else is a storage failure.
		var awsErr *awserrors.AWSError
		if errors.As(err, &awsErr) {
			return nil, err
		}
		return nil, ErrInternalError
	}

	// The tag store follows the pool record, not the other way round: with
	// the replace inside the mutation a pool write that later failed would
	// have left the new tag set applied to an update that never landed. A
	// single replace write swaps the whole set under the tag store's lock,
	// so a failure here leaves no partially-untagged resource behind.
	if in.TagsProvided {
		if err := store.Replace(updated.Arn, in.Tags); err != nil {
			return nil, ErrInternalError
		}
	}

	return poolToOut(updated), nil
}

// GetIdentityPoolRolesResult is the transport-agnostic role configuration of
// an identity pool. RoleMappings keeps the store representation; only the
// serialising layers consume it.
type GetIdentityPoolRolesResult struct {
	IdentityPoolID string
	AuthRole       string
	UnauthRole     string
	RoleMappings   map[string]cognitoidentitystore.RoleMapping
}

// getIdentityPoolRolesCore is the single entry point for GetIdentityPoolRoles.
func (s *CognitoIdentityService) getIdentityPoolRolesCore(reqCtx *request.RequestContext, poolID string) (*GetIdentityPoolRolesResult, error) {
	if !validateIdentityPoolId(poolID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	authRole, unauthRole, mappings, err := store.GetIdentityPoolRoles(poolID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	return &GetIdentityPoolRolesResult{
		IdentityPoolID: poolID,
		AuthRole:       authRole,
		UnauthRole:     unauthRole,
		RoleMappings:   mappings,
	}, nil
}

// SetIdentityPoolRolesInput carries every field that SetIdentityPoolRoles
// needs. RolesRaw and RoleMappingsRaw carry the untyped wire values.
type SetIdentityPoolRolesInput struct {
	IdentityPoolID  string
	RolesProvided   bool
	RolesRaw        interface{}
	RoleMappingsRaw interface{}
}

// setIdentityPoolRolesCore is the single entry point for SetIdentityPoolRoles.
func (s *CognitoIdentityService) setIdentityPoolRolesCore(reqCtx *request.RequestContext, in SetIdentityPoolRolesInput) error {
	if !validateIdentityPoolId(in.IdentityPoolID) {
		return ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	authRole, unauthRole := "", ""
	if in.RolesProvided {
		rolesMap, ok := in.RolesRaw.(map[string]interface{})
		if !ok {
			return ErrInvalidParameter
		}
		for k := range rolesMap {
			if !validRoleTypes[k] {
				return ErrInvalidParameter
			}
		}
		if !validateMapSize(len(rolesMap), 2) {
			return ErrInvalidParameter
		}
		if v, ok := rolesMap["authenticated"].(string); ok {
			authRole = v
		}
		if v, ok := rolesMap["unauthenticated"].(string); ok {
			unauthRole = v
		}
		if authRole != "" && !validateRoleARN(authRole) {
			return ErrInvalidParameter
		}
		if unauthRole != "" && !validateRoleARN(unauthRole) {
			return ErrInvalidParameter
		}
	} else {
		// Roles is semantically required by AWS. Absent Roles would silently
		// clear all existing roles — a destructive operation that AWS rejects
		// with InvalidParameterException.
		return ErrInvalidParameter
	}
	if !validateRoleKeys(authRole, unauthRole) {
		return ErrInvalidParameter
	}
	mappingDTOs, err := parseRoleMappings(in.RoleMappingsRaw)
	if err != nil {
		return err
	}
	if !validateMapSize(len(mappingDTOs), maxRoleMappingsPerPool) {
		return ErrInvalidParameter
	}

	if err := store.SetIdentityPoolRoles(in.IdentityPoolID, authRole, unauthRole, roleMappingMapToStore(mappingDTOs)); err != nil {
		return mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	return nil
}

// parseCognitoIdentityProviders parses the raw CognitoIdentityProviders wire
// value (nil when the member is absent) and validates every entry against the
// Smithy provider shapes.
func parseCognitoIdentityProviders(raw interface{}) ([]ProviderOut, error) {
	if raw == nil {
		return nil, nil
	}
	slice, ok := raw.([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	providers := make([]ProviderOut, 0)
	for _, v := range slice {
		m, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		provider := ProviderOut{}
		if name, ok := m["ProviderName"].(string); ok {
			if !validateProviderName(name) {
				return nil, ErrInvalidParameter
			}
			provider.ProviderName = name
		}
		if clientID, ok := m["ClientId"].(string); ok {
			if !validateProviderClientId(clientID) {
				return nil, ErrInvalidParameter
			}
			provider.ClientID = clientID
		}
		if check, ok := m["ServerSideTokenCheck"].(bool); ok {
			provider.ServerSideTokenCheck = check
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

// parseRoleMappings parses the raw RoleMappings wire value (nil when the
// member is absent) and validates the complete mapping set.
func parseRoleMappings(raw interface{}) (map[string]RoleMappingInput, error) {
	if raw == nil {
		return nil, nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	result := make(map[string]RoleMappingInput)
	for k, v := range m {
		mapping, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		rm := RoleMappingInput{}
		if t, ok := mapping["Type"].(string); ok {
			rm.Type = t
		}
		if arr, ok := mapping["AmbiguousRoleResolution"].(string); ok {
			rm.AmbiguousRoleResolution = arr
		}
		if rules, ok := mapping["RulesConfiguration"].(map[string]interface{}); ok {
			rc, err := parseRulesConfiguration(rules)
			if err != nil {
				return nil, err
			}
			rm.RulesConfiguration = rc
		}
		result[k] = rm
	}
	if !validateRoleMappings(result) {
		return nil, ErrInvalidParameter
	}
	return result, nil
}

func parseRulesConfiguration(m map[string]interface{}) (*RulesConfigInput, error) {
	rules, ok := m["Rules"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	config := &RulesConfigInput{
		Rules: make([]MappingRuleInput, 0),
	}
	for _, r := range rules {
		rule, ok := r.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		mr := MappingRuleInput{}
		if claim, ok := rule["Claim"].(string); ok {
			mr.Claim = claim
		}
		if matchType, ok := rule["MatchType"].(string); ok {
			mr.MatchType = matchType
		}
		if value, ok := rule["Value"].(string); ok {
			mr.Value = value
		}
		if roleArn, ok := rule["RoleARN"].(string); ok {
			mr.RoleARN = roleArn
		}
		config.Rules = append(config.Rules, mr)
	}
	return config, nil
}

// ---------------------------------------------------------------------------
// Conversion helpers
// ---------------------------------------------------------------------------

func poolToOut(p *cognitoidentitystore.IdentityPool) *IdentityPoolOut {
	out := &IdentityPoolOut{
		ID:                             p.ID,
		Name:                           p.Name,
		AllowUnauthenticatedIdentities: p.AllowUnauthenticatedIdentities,
		AllowClassicFlow:               p.AllowClassicFlow,
		DeveloperProviderName:          p.DeveloperProviderName,
		SupportedLoginProviders:        p.SupportedLoginProviders,
		OpenIdConnectProviderARNs:      p.OpenIdConnectProviderARNs,
		SamlProviderARNs:               p.SamlProviderARNs,
		Tags:                           p.Tags,
		Arn:                            p.Arn,
	}
	for _, cp := range p.CognitoIdentityProviders {
		out.CognitoIdentityProviders = append(out.CognitoIdentityProviders, ProviderOut{
			ProviderName:         cp.ProviderName,
			ClientID:             cp.ClientID,
			ServerSideTokenCheck: cp.ServerSideTokenCheck,
		})
	}
	return out
}

// poolOutToHTTP converts the transport-agnostic IdentityPoolOut to the HTTP
// API response format (map[string]interface{}).
func poolOutToHTTP(p *IdentityPoolOut) map[string]interface{} {
	result := map[string]interface{}{
		"IdentityPoolId":                 p.ID,
		"IdentityPoolName":               p.Name,
		"AllowUnauthenticatedIdentities": p.AllowUnauthenticatedIdentities,
		"AllowClassicFlow":               p.AllowClassicFlow,
	}
	if len(p.CognitoIdentityProviders) > 0 {
		providers := make([]map[string]interface{}, 0, len(p.CognitoIdentityProviders))
		for _, cp := range p.CognitoIdentityProviders {
			providers = append(providers, map[string]interface{}{
				"ProviderName":         cp.ProviderName,
				"ClientId":             cp.ClientID,
				"ServerSideTokenCheck": cp.ServerSideTokenCheck,
			})
		}
		result["CognitoIdentityProviders"] = providers
	}
	if p.DeveloperProviderName != "" {
		result["DeveloperProviderName"] = p.DeveloperProviderName
	}
	if len(p.SupportedLoginProviders) > 0 {
		result["SupportedLoginProviders"] = p.SupportedLoginProviders
	}
	if len(p.OpenIdConnectProviderARNs) > 0 {
		result["OpenIdConnectProviderARNs"] = p.OpenIdConnectProviderARNs
	}
	if len(p.SamlProviderARNs) > 0 {
		result["SamlProviderARNs"] = p.SamlProviderARNs
	}
	if len(p.Tags) > 0 {
		result["IdentityPoolTags"] = p.Tags
	}
	return result
}
