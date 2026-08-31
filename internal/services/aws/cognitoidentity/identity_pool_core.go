package cognitoidentity

import (
	"errors"

	"vorpalstacks/internal/common/request"
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

	if !validateMapSize(len(in.SupportedLoginProviders), 10) {
		return nil, ErrInvalidParameter
	}
	if in.TagsProvided && len(in.Tags) > 0 {
		if !validateTagKeys(in.Tags) || !validateTagValues(in.Tags) {
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
	tags, _ := store.List(pool.Arn)
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
func (s *CognitoIdentityService) updateIdentityPoolCore(reqCtx *request.RequestContext, in UpdateIdentityPoolInput) (*IdentityPoolOut, error) {
	if !validateIdentityPoolId(in.IdentityPoolID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pool, err := store.GetIdentityPool(in.IdentityPoolID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	// IdentityPoolName is @required in the Smithy IdentityPool shape.
	if in.PoolName == "" {
		return nil, ErrInvalidParameter
	}
	if !validateIdentityPoolName(in.PoolName) {
		return nil, ErrInvalidParameter
	}
	pool.Name = in.PoolName

	// AllowUnauthenticatedIdentities is @required in the Smithy shape.
	if !in.AllowUnauthProvided {
		return nil, ErrInvalidParameter
	}
	if b, ok := in.AllowUnauthRaw.(bool); ok {
		pool.AllowUnauthenticatedIdentities = b
	} else {
		return nil, ErrInvalidParameter
	}

	if in.AllowClassicProvided {
		b, ok := in.AllowClassicRaw.(bool)
		if !ok {
			return nil, ErrInvalidParameter
		}
		pool.AllowClassicFlow = b
	}
	if in.DeveloperProviderName != "" {
		if !validateDeveloperProviderName(in.DeveloperProviderName) {
			return nil, ErrInvalidParameter
		}
		pool.DeveloperProviderName = in.DeveloperProviderName
	}
	if providers, err := parseCognitoIdentityProviders(in.ProvidersRaw); err != nil {
		return nil, err
	} else if len(providers) > 0 {
		pool.CognitoIdentityProviders = providerOutsToStore(providers)
	}
	if len(in.SupportedLoginProviders) > 0 {
		if !validateMapSize(len(in.SupportedLoginProviders), 10) {
			return nil, ErrInvalidParameter
		}
		pool.SupportedLoginProviders = in.SupportedLoginProviders
	}
	if len(in.OpenIdConnectProviderARNs) > 0 {
		for _, arn := range in.OpenIdConnectProviderARNs {
			if !validateRoleARN(arn) {
				return nil, ErrInvalidParameter
			}
		}
		pool.OpenIdConnectProviderARNs = in.OpenIdConnectProviderARNs
	}
	if len(in.SamlProviderARNs) > 0 {
		for _, arn := range in.SamlProviderARNs {
			if !validateRoleARN(arn) {
				return nil, ErrInvalidParameter
			}
		}
		pool.SamlProviderARNs = in.SamlProviderARNs
	}

	var updatedTags map[string]string
	if in.TagsProvided {
		if !validateTagKeys(in.Tags) || !validateTagValues(in.Tags) {
			return nil, ErrInvalidParameter
		}
		updatedTags = in.Tags
		// A single replace write swaps the whole tag set under the tag
		// store's lock, so a failure cannot leave a partially-untagged
		// resource and no rollback path is needed.
		if err := store.Replace(pool.Arn, updatedTags); err != nil {
			return nil, ErrInternalError
		}
	} else {
		updatedTags, _ = store.List(pool.Arn)
	}

	if err := store.UpdateIdentityPool(pool); err != nil {
		return nil, ErrInternalError
	}

	pool.Tags = updatedTags
	return poolToOut(pool), nil
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
	if !validateMapSize(len(mappingDTOs), 10) {
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
