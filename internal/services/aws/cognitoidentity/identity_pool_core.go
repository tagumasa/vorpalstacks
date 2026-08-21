package cognitoidentity

import (
	"errors"

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

	// Enforce per-account identity pool limit (AWS default: 50).
	existing, err := store.ListIdentityPools(storecommon.ListOptions{MaxItems: 100})
	if err != nil {
		return nil, ErrInternalError
	}
	if len(existing.Items) >= 50 {
		return nil, ErrLimitExceeded
	}

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
