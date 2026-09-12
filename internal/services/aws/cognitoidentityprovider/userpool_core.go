package cognitoidentityprovider

import (
	"errors"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
	"vorpalstacks/pkg/vsjwt"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// ListUserPoolsInput carries the pagination parameters for ListUserPools in a
// format independent of the wire protocol (HTTP Query vs gRPC-Web).
type ListUserPoolsInput struct {
	MaxResults int
	NextToken  string
}

// UserPoolSummary is the transport-agnostic representation of a single
// UserPoolDescriptionType entry returned by ListUserPools.
type UserPoolSummary struct {
	ID               string
	Name             string
	Arn              string
	Status           string
	CreationDate     time.Time
	LastModifiedDate time.Time
}

// ListUserPoolsResult is the paginated result of ListUserPools.
type ListUserPoolsResult struct {
	UserPools []UserPoolSummary
	NextToken string
}

// CreateUserPoolInput carries the pre-built store-level UserPool, the region
// for store resolution, and optional tags. Both the HTTP handler (via
// applyUserPoolUpdates) and the admin handler (via createUserPoolFromAdmin)
// construct the UserPool and delegate to createUserPoolCore so that
// validation, persistence, and tagging share a single code path.
type CreateUserPoolInput struct {
	Pool   *cognitostore.UserPool
	Region string
	Tags   map[string]string
}

// AdminCreateUserPoolInput carries the simplified parameters that the admin
// console provides when creating a user pool. Unlike the HTTP API which
// constructs a full store UserPool from query parameters, the admin console
// sends proto-derived primitives and relies on the Core to build the store
// type.
type AdminCreateUserPoolInput struct {
	PoolName          string
	Region            string
	AutoVerifiedAttrs []string
	PasswordPolicy    *AdminPasswordPolicy
	Tags              map[string]string
}

// AdminPasswordPolicy is the service-layer representation of the Cognito
// password policy, decoupled from the store type so that admin handler DTOs
// never reference store packages.
type AdminPasswordPolicy struct {
	MinimumLength                 int
	RequireUppercase              bool
	RequireLowercase              bool
	RequireNumbers                bool
	RequireSymbols                bool
	TemporaryPasswordValidityDays int
	PasswordHistorySize           int
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// listUserPoolsCore lists Cognito user pools with pagination. It resolves the
// store for the given region internally so that callers (HTTP handler, admin
// handler) never need to reference store types directly.
func (s *CognitoService) listUserPoolsCore(region string, in ListUserPoolsInput) (*ListUserPoolsResult, error) {
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	maxResults := applyListLimitDefaults(in.MaxResults)

	result, err := store.ListUserPoolsPaginated(storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   in.NextToken,
	})
	if err != nil {
		return nil, err
	}

	pools := make([]UserPoolSummary, 0, len(result.Items))
	for _, pool := range result.Items {
		pools = append(pools, UserPoolSummary{
			ID:               pool.ID,
			Name:             pool.Name,
			Arn:              pool.Arn,
			Status:           pool.Status,
			CreationDate:     pool.CreationDate,
			LastModifiedDate: pool.LastModifiedDate,
		})
	}

	return &ListUserPoolsResult{
		UserPools: pools,
		NextToken: result.NextMarker,
	}, nil
}

// getJWKSCore returns the JSON Web Key Set for a user pool. The JWKS route
// serves external verifiers, so the pool's key material is decoded and the
// key set built here, in Core; the transport serves the result as-is.
func (s *CognitoService) getJWKSCore(region, userPoolID string) (interface{}, error) {
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	privateKey, err := vsjwt.DecodePrivateKeyFromPEM(userPool.JwtPrivateKey)
	if err != nil {
		return nil, ErrInternalError
	}

	issuer := fmt.Sprintf("https://%s/%s", cognitoIdpHost(region), userPoolID)
	jwtManager, err := vsjwt.NewManager(privateKey, userPool.JwtKeyID, issuer)
	if err != nil {
		return nil, ErrInternalError
	}
	return jwtManager.GetJWKS(), nil
}

// createUserPoolCore creates a new Cognito user pool. Validation, password
// policy enforcement, persistence, and tagging all happen here so that both
// the HTTP handler and the admin handler share a single code path. Returns
// the created store-level UserPool so each transport layer can convert to its
// own wire format (formatUserPool for HTTP, userPoolToProto for gRPC-Web).
func (s *CognitoService) createUserPoolCore(in CreateUserPoolInput) (*cognitostore.UserPool, error) {
	if in.Pool == nil || in.Pool.Name == "" {
		return nil, ErrInvalidParameter
	}

	// Single model-derived validation entry point shared by every transport.
	if err := validateUserPoolConfig(in.Pool); err != nil {
		return nil, err
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, err
	}

	created, err := store.CreateUserPool(in.Pool)
	if err != nil {
		return nil, err
	}

	if len(in.Tags) > 0 {
		if err := store.Tag(created.Arn, in.Tags); err != nil {
			return nil, err
		}
	}

	return created, nil
}

// createUserPoolFromAdmin is the Core entry point for the admin console.
// It builds a store-level UserPool from the simplified admin parameters
// and delegates to createUserPoolCore, ensuring that validation, password
// policy enforcement, and persistence follow the single code path.
func (s *CognitoService) createUserPoolFromAdmin(in AdminCreateUserPoolInput) (*cognitostore.UserPool, error) {
	pool := cognitostore.NewUserPool(in.PoolName, in.Region)
	pool.AutoVerifiedAttributes = in.AutoVerifiedAttrs
	if in.PasswordPolicy != nil {
		pool.PasswordPolicy = &cognitostore.PasswordPolicy{
			MinimumLength:                 in.PasswordPolicy.MinimumLength,
			RequireUppercase:              in.PasswordPolicy.RequireUppercase,
			RequireLowercase:              in.PasswordPolicy.RequireLowercase,
			RequireNumbers:                in.PasswordPolicy.RequireNumbers,
			RequireSymbols:                in.PasswordPolicy.RequireSymbols,
			TemporaryPasswordValidityDays: in.PasswordPolicy.TemporaryPasswordValidityDays,
			PasswordHistorySize:           in.PasswordPolicy.PasswordHistorySize,
		}
	}
	return s.createUserPoolCore(CreateUserPoolInput{
		Pool:   pool,
		Region: in.Region,
		Tags:   in.Tags,
	})
}

// createUserPoolFromRequestCore is the Core entry point for the HTTP
// CreateUserPool operation: it builds the store-level UserPool from the raw
// wire request (identity, create-only members, configuration members and
// tags) and delegates to createUserPoolCore, so validation and persistence
// follow the single code path shared with the admin console.
func (s *CognitoService) createUserPoolFromRequestCore(region, poolName string, req *request.ParsedRequest) (*cognitostore.UserPool, error) {
	pool, err := s.newUserPoolCore(poolName, region)
	if err != nil {
		return nil, err
	}
	// CreateUserPool is the only operation that carries the Schema member;
	// apply it before the shared member application so the whole-pool
	// validation still sees the schema definitions.
	schema, err := parseSchemaAttributes(req)
	if err != nil {
		return nil, err
	}
	pool.SchemaAttributes = schema
	if err := applyCreateUserPoolRequest(pool, req, region); err != nil {
		return nil, err
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "UserPoolTags"))
	return s.createUserPoolCore(CreateUserPoolInput{
		Pool:   pool,
		Region: region,
		Tags:   tags,
	})
}

// deleteUserPoolCore deletes a Cognito user pool by ID.
func (s *CognitoService) deleteUserPoolCore(region, userPoolID string) error {
	if userPoolID == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteUserPool(userPoolID); err != nil {
		return ErrResourceNotFound
	}
	return nil
}

// describeUserPoolCore retrieves a single Cognito user pool by ID, including
// its tags. Returns the store-level UserPool type so that both the HTTP
// handler (via formatUserPool) and the admin handler (via userPoolToProto in
// admin_handler_convert.go) can convert to their respective wire formats.
func (s *CognitoService) describeUserPoolCore(region, userPoolID string) (*cognitostore.UserPool, error) {
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	tags, _ := store.ListAsSlice(userPool.Arn)
	userPool.Tags = tags

	return userPool, nil
}

// getUserPoolCore loads a user pool without the tag side effect of
// describeUserPoolCore; mutation flows (UpdateUserPool, MFA configuration)
// must not bake tags into the persisted pool record.
func (s *CognitoService) getUserPoolCore(region, userPoolID string) (*cognitostore.UserPool, error) {
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return userPool, nil
}

// newUserPoolCore validates the create-path required members and constructs
// the store object. Store constructors are store-package calls, so
// construction lives behind Core.
func (s *CognitoService) newUserPoolCore(poolName, region string) (*cognitostore.UserPool, error) {
	if poolName == "" {
		return nil, ErrInvalidParameter
	}
	if !validateUserPoolNamePattern(poolName) {
		return nil, ErrInvalidParameter
	}
	return cognitostore.NewUserPool(poolName, region), nil
}

// UpdateUserPoolInput carries an UpdateUserPool request in the format the
// Core needs: the region for store resolution, the target pool and the
// parsed wire request whose present members replace the pool configuration.
type UpdateUserPoolInput struct {
	Region     string
	UserPoolID string
	Req        *request.ParsedRequest
}

// updateUserPoolCore applies the model's UpdateUserPool reset contract as
// one serialised pool mutation: the request rebuilds the pool configuration
// with omitted members at their defaults, UserPoolTags replaces the tag set
// (its default is an untagged pool), and a pool deleted before or during the
// mutation surfaces ResourceNotFoundException instead of an internal error.
func (s *CognitoService) updateUserPoolCore(in UpdateUserPoolInput) error {
	if in.UserPoolID == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return err
	}

	err = store.UpdateUserPoolFunc(in.UserPoolID, func(pool *cognitostore.UserPool) error {
		rebuilt, rerr := rebuildUserPoolFromUpdate(pool, in.Req, in.Region)
		if rerr != nil {
			return rerr
		}
		*pool = *rebuilt
		// The reset contract governs tags like every other member: a
		// request that omits UserPoolTags leaves the pool untagged (the
		// default), a request that carries it replaces the whole set.
		tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(in.Req.Parameters, "UserPoolTags"))
		return store.Replace(pool.Arn, tags)
	})
	if err != nil {
		if errors.Is(err, cognitostore.ErrUserPoolNotFound) {
			return ErrResourceNotFound
		}
		// Validation rejections from the rebuild are already wire-shaped
		// service errors; everything else is a storage failure.
		var awsErr *awserrors.AWSError
		if errors.As(err, &awsErr) {
			return err
		}
		return ErrInternalError
	}
	return nil
}

// SetUserPoolMfaConfigInput carries the raw nested wire members of
// SetUserPoolMfaConfig; absent members arrive as nil maps.
type SetUserPoolMfaConfigInput struct {
	Region                        string
	UserPoolID                    string
	MfaConfiguration              string
	SmsMfaConfiguration           map[string]interface{}
	SoftwareTokenMfaConfiguration map[string]interface{}
	EmailMfaConfiguration         map[string]interface{}
	WebAuthnConfiguration         map[string]interface{}
}

// setUserPoolMfaConfigCore applies the MFA configuration members onto the
// stored pool as one serialised pool mutation and returns the updated record
// for response serialisation.
func (s *CognitoService) setUserPoolMfaConfigCore(in SetUserPoolMfaConfigInput) (*cognitostore.UserPool, error) {
	if in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, err
	}

	var updated *cognitostore.UserPool
	err = store.UpdateUserPoolFunc(in.UserPoolID, func(userPool *cognitostore.UserPool) error {
		updated = userPool
		if in.MfaConfiguration != "" {
			if !validateUserPoolMfaConfig(in.MfaConfiguration) {
				return ErrInvalidParameter
			}
			userPool.MfaConfiguration = in.MfaConfiguration
		}

		if m := in.SmsMfaConfiguration; m != nil {
			smsMfa := &cognitostore.SmsMfaConfig{}
			if v, ok := m["SmsAuthenticationMessage"].(string); ok {
				smsMfa.SmsAuthenticationMessage = v
			}
			if smsConfig, ok := m["SmsConfiguration"].(map[string]interface{}); ok {
				poolSmsConfig := &cognitostore.SmsConfiguration{}
				if v, ok := smsConfig["SnsCallerArn"].(string); ok {
					poolSmsConfig.SnsCallerArn = v
				}
				if v, ok := smsConfig["ExternalId"].(string); ok {
					poolSmsConfig.ExternalId = v
				}
				if v, ok := smsConfig["SnsRegion"].(string); ok {
					poolSmsConfig.SnsRegion = v
				}
				smsMfa.SmsConfiguration = poolSmsConfig
			}
			userPool.MfaConfigurationSms = smsMfa
		}
		if m := in.SoftwareTokenMfaConfiguration; m != nil {
			swMfa := &cognitostore.MfaConfigurationType{}
			if enabled, ok := m["Enabled"].(bool); ok {
				swMfa.Enabled = enabled
			}
			userPool.MfaConfigurationSoftwareToken = swMfa
		}
		if m := in.EmailMfaConfiguration; m != nil {
			emailMfa := &cognitostore.EmailMfaConfig{}
			if v, ok := m["Message"].(string); ok {
				emailMfa.Message = v
			}
			if v, ok := m["Subject"].(string); ok {
				emailMfa.Subject = v
			}
			userPool.EmailMfaConfig = emailMfa
		}
		if m := in.WebAuthnConfiguration; m != nil {
			wa := &cognitostore.WebAuthnConfiguration{}
			if v, ok := m["RelyingPartyId"].(string); ok {
				wa.RelyingPartyId = v
			}
			if v, ok := m["UserVerification"].(string); ok {
				wa.UserVerification = v
			}
			userPool.WebAuthnConfiguration = wa
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, cognitostore.ErrUserPoolNotFound) {
			return nil, ErrResourceNotFound
		}
		var awsErr *awserrors.AWSError
		if errors.As(err, &awsErr) {
			return nil, err
		}
		return nil, ErrInternalError
	}

	return updated, nil
}

// addCustomAttributesCore appends the requested custom attribute schemas to
// the pool as one serialised pool mutation, rejecting a list outside the
// CustomAttributesListType bounds, invalid names, non-member data types,
// non-map entries and duplicates of existing schema attributes.
func (s *CognitoService) addCustomAttributesCore(region, userPoolID string, customAttributes []interface{}) error {
	if userPoolID == "" {
		return ErrInvalidParameter
	}
	if len(customAttributes) < cognitostore.MinCustomAttributesPerAdd ||
		len(customAttributes) > cognitostore.MaxCustomAttributesPerAdd {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	err = store.UpdateUserPoolFunc(userPoolID, func(pool *cognitostore.UserPool) error {
		for _, a := range customAttributes {
			m, ok := a.(map[string]interface{})
			if !ok {
				return ErrInvalidParameter
			}
			attrName := getStringParam(m, "Name")
			attrType := getStringParam(m, "AttributeDataType")
			if attrName == "" || attrType == "" {
				return ErrInvalidParameter
			}
			if err := validateCustomAttributeName(attrName); err != nil {
				return err
			}
			if !validateAttributeDataType(attrType) {
				return ErrInvalidParameter
			}
			for _, existing := range pool.SchemaAttributes {
				if existing.Name == attrName {
					return ErrInvalidParameter
				}
			}
			newAttr := cognitostore.SchemaAttributeType{
				Name:              attrName,
				AttributeDataType: attrType,
			}
			if dev, ok := m["DeveloperOnlyAttribute"].(bool); ok {
				newAttr.DeveloperOnlyAttribute = dev
			}
			if mut, ok := m["Mutable"].(bool); ok {
				newAttr.Mutable = mut
			}
			if reqVal, ok := m["Required"].(bool); ok {
				newAttr.Required = reqVal
			}
			if nac, ok := m["NumberAttributeConstraints"].(map[string]interface{}); ok {
				nc := &cognitostore.NumberAttributeConstraints{}
				if v, ok := nac["MinValue"].(string); ok {
					nc.MinValue = v
				}
				if v, ok := nac["MaxValue"].(string); ok {
					nc.MaxValue = v
				}
				newAttr.NumberAttributeConstraints = nc
			}
			if sac, ok := m["StringAttributeConstraints"].(map[string]interface{}); ok {
				sc := &cognitostore.StringAttributeConstraints{}
				if v, ok := sac["MinLength"].(string); ok {
					sc.MinLength = v
				}
				if v, ok := sac["MaxLength"].(string); ok {
					sc.MaxLength = v
				}
				newAttr.StringAttributeConstraints = sc
			}
			pool.SchemaAttributes = append(pool.SchemaAttributes, newAttr)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, cognitostore.ErrUserPoolNotFound) {
			return ErrResourceNotFound
		}
		var awsErr *awserrors.AWSError
		if errors.As(err, &awsErr) {
			return err
		}
		return ErrInternalError
	}

	return nil
}
