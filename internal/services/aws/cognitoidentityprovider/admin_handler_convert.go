package cognitoidentityprovider

import (
	"google.golang.org/protobuf/proto"

	pb "vorpalstacks/internal/pb/aws/cognitoidentityprovider"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	"vorpalstacks/internal/utils/timeutils"
)

// This file is the sole location where admin handler proto-conversion helpers
// may import store packages (the sole exception to the store-import prohibition).

// ---------------------------------------------------------------------------
// Enum / scalar helpers
// ---------------------------------------------------------------------------

func aliasAttributeToProto(s string) pb.AliasAttributeType {
	switch s {
	case "phone_number":
		return pb.AliasAttributeType_ALIAS_ATTRIBUTE_TYPE_PHONE_NUMBER
	case "email":
		return pb.AliasAttributeType_ALIAS_ATTRIBUTE_TYPE_EMAIL
	case "preferred_username":
		return pb.AliasAttributeType_ALIAS_ATTRIBUTE_TYPE_PREFERRED_USERNAME
	default:
		return 0
	}
}

func usernameAttributeToProto(s string) pb.UsernameAttributeType {
	switch s {
	case "phone_number":
		return pb.UsernameAttributeType_USERNAME_ATTRIBUTE_TYPE_PHONE_NUMBER
	case "email":
		return pb.UsernameAttributeType_USERNAME_ATTRIBUTE_TYPE_EMAIL
	default:
		return 0
	}
}

func verifiedAttributeToProto(s string) pb.VerifiedAttributeType {
	switch s {
	case "phone_number":
		return pb.VerifiedAttributeType_VERIFIED_ATTRIBUTE_TYPE_PHONE_NUMBER
	case "email":
		return pb.VerifiedAttributeType_VERIFIED_ATTRIBUTE_TYPE_EMAIL
	default:
		return 0
	}
}

func statusToProto(s string) pb.StatusType {
	if s == "ENABLED" {
		return pb.StatusType_STATUS_TYPE_ENABLED
	}
	return pb.StatusType_STATUS_TYPE_DISABLED
}

func mfaConfigurationToProto(s string) pb.UserPoolMfaType {
	switch s {
	case "ON":
		return pb.UserPoolMfaType_USER_POOL_MFA_TYPE_ON
	case "OPTIONAL":
		return pb.UserPoolMfaType_USER_POOL_MFA_TYPE_OPTIONAL
	default:
		return pb.UserPoolMfaType_USER_POOL_MFA_TYPE_OFF
	}
}

func deletionProtectionToProto(s string) pb.DeletionProtectionType {
	if s == "ACTIVE" {
		return pb.DeletionProtectionType_DELETION_PROTECTION_TYPE_ACTIVE
	}
	return pb.DeletionProtectionType_DELETION_PROTECTION_TYPE_INACTIVE
}

// ---------------------------------------------------------------------------
// UserPool → proto
// ---------------------------------------------------------------------------

// userPoolToProto converts a store-level UserPool to the proto UserPoolType.
// Complex nested configuration types (EmailConfiguration, SmsConfiguration,
// LambdaConfig, etc.) are omitted — the admin console shows these as raw
// JSON when needed. All scalar fields, tags, password policy, and enum
// arrays are fully converted.
func userPoolToProto(pool *cognitostore.UserPool) *pb.UserPoolType {
	if pool == nil {
		return nil
	}

	result := &pb.UserPoolType{
		Id:                 pool.ID,
		Name:               pool.Name,
		Arn:                pool.Arn,
		Status:             statusToProto(pool.Status),
		Mfaconfiguration:   mfaConfigurationToProto(pool.MfaConfiguration),
		Deletionprotection: deletionProtectionToProto(pool.DeletionProtection),
	}

	if !pool.CreationDate.IsZero() {
		result.Creationdate = pool.CreationDate.Format(timeutils.ISO8601UTCFormat)
	}
	if !pool.LastModifiedDate.IsZero() {
		result.Lastmodifieddate = pool.LastModifiedDate.Format(timeutils.ISO8601UTCFormat)
	}

	if len(pool.AliasAttributes) > 0 {
		result.Aliasattributes = make([]pb.AliasAttributeType, len(pool.AliasAttributes))
		for i, a := range pool.AliasAttributes {
			result.Aliasattributes[i] = aliasAttributeToProto(a)
		}
	}
	if len(pool.UsernameAttributes) > 0 {
		result.Usernameattributes = make([]pb.UsernameAttributeType, len(pool.UsernameAttributes))
		for i, u := range pool.UsernameAttributes {
			result.Usernameattributes[i] = usernameAttributeToProto(u)
		}
	}
	if len(pool.AutoVerifiedAttributes) > 0 {
		result.Autoverifiedattributes = make([]pb.VerifiedAttributeType, len(pool.AutoVerifiedAttributes))
		for i, a := range pool.AutoVerifiedAttributes {
			result.Autoverifiedattributes[i] = verifiedAttributeToProto(a)
		}
	}

	if pool.PasswordPolicy != nil {
		result.Policies = &pb.UserPoolPolicyType{
			Passwordpolicy: &pb.PasswordPolicyType{
				Minimumlength:                 proto.Int32(int32(pool.PasswordPolicy.MinimumLength)),
				Requireuppercase:              proto.Bool(pool.PasswordPolicy.RequireUppercase),
				Requirelowercase:              proto.Bool(pool.PasswordPolicy.RequireLowercase),
				Requirenumbers:                proto.Bool(pool.PasswordPolicy.RequireNumbers),
				Requiresymbols:                proto.Bool(pool.PasswordPolicy.RequireSymbols),
				Temporarypasswordvaliditydays: proto.Int32(int32(pool.PasswordPolicy.TemporaryPasswordValidityDays)),
				Passwordhistorysize:           proto.Int32(int32(pool.PasswordPolicy.PasswordHistorySize)),
			},
		}
	}

	if pool.EstimatedNumberOfUsers > 0 {
		result.Estimatednumberofusers = proto.Int32(int32(pool.EstimatedNumberOfUsers))
	}

	if len(pool.Tags) > 0 {
		tags := make(map[string]string, len(pool.Tags))
		for _, t := range pool.Tags {
			tags[t.Key] = t.Value
		}
		result.Userpooltags = tags
	}

	return result
}

// ---------------------------------------------------------------------------
// IdentityProvider → proto
// ---------------------------------------------------------------------------

func identityProviderTypeToProto(s string) pb.IdentityProviderTypeType {
	switch s {
	case "LoginWithAmazon", "LOGIN_WITH_AMAZON":
		return pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_LOGINWITHAMAZON
	case "Google", "GOOGLE":
		return pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_GOOGLE
	case "Facebook", "FACEBOOK":
		return pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_FACEBOOK
	case "SAML":
		return pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_SAML
	case "SignInWithApple", "SIGN_IN_WITH_APPLE":
		return pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_SIGNINWITHAPPLE
	case "OIDC":
		return pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_OIDC
	default:
		return 0
	}
}

func identityProviderTypeFromProto(t pb.IdentityProviderTypeType) string {
	switch t {
	case pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_LOGINWITHAMAZON:
		return "LoginWithAmazon"
	case pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_GOOGLE:
		return "Google"
	case pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_FACEBOOK:
		return "Facebook"
	case pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_SAML:
		return "SAML"
	case pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_SIGNINWITHAPPLE:
		return "SignInWithApple"
	case pb.IdentityProviderTypeType_IDENTITY_PROVIDER_TYPE_TYPE_OIDC:
		return "OIDC"
	default:
		return ""
	}
}

// identityProviderToProto converts a store-level IdentityProvider to proto.
func identityProviderToProto(ip *cognitostore.IdentityProvider) *pb.IdentityProviderType {
	if ip == nil {
		return nil
	}

	result := &pb.IdentityProviderType{
		Providername:     ip.ProviderName,
		Providertype:     identityProviderTypeToProto(ip.ProviderType),
		Providerdetails:  ip.ProviderDetails,
		Attributemapping: ip.AttributeMapping,
		Idpidentifiers:   ip.IdpIdentifiers,
	}
	if !ip.CreationDate.IsZero() {
		result.Creationdate = ip.CreationDate.Format(timeutils.ISO8601UTCFormat)
	}
	if !ip.LastModifiedDate.IsZero() {
		result.Lastmodifieddate = ip.LastModifiedDate.Format(timeutils.ISO8601UTCFormat)
	}
	return result
}

// providerDescriptionToProto converts a store-level IdentityProvider to the
// list-entry proto type used by ListIdentityProviders.
func providerDescriptionToProto(ip *cognitostore.IdentityProvider) *pb.ProviderDescription {
	if ip == nil {
		return nil
	}
	result := &pb.ProviderDescription{
		Providername: ip.ProviderName,
		Providertype: identityProviderTypeToProto(ip.ProviderType),
	}
	if !ip.CreationDate.IsZero() {
		result.Creationdate = ip.CreationDate.Format(timeutils.ISO8601UTCFormat)
	}
	if !ip.LastModifiedDate.IsZero() {
		result.Lastmodifieddate = ip.LastModifiedDate.Format(timeutils.ISO8601UTCFormat)
	}
	return result
}

// ---------------------------------------------------------------------------
// User → proto
// ---------------------------------------------------------------------------

func userStatusToProto(s string) pb.UserStatusType {
	switch s {
	case "CONFIRMED":
		return pb.UserStatusType_USER_STATUS_TYPE_CONFIRMED
	case "UNCONFIRMED":
		return pb.UserStatusType_USER_STATUS_TYPE_UNCONFIRMED
	case "FORCE_CHANGE_PASSWORD":
		return pb.UserStatusType_USER_STATUS_TYPE_FORCE_CHANGE_PASSWORD
	case "RESET_REQUIRED":
		return pb.UserStatusType_USER_STATUS_TYPE_RESET_REQUIRED
	case "COMPROMISED":
		return pb.UserStatusType_USER_STATUS_TYPE_COMPROMISED
	case "ARCHIVED":
		return pb.UserStatusType_USER_STATUS_TYPE_ARCHIVED
	case "EXTERNAL_PROVIDER":
		return pb.UserStatusType_USER_STATUS_TYPE_EXTERNAL_PROVIDER
	default:
		return pb.UserStatusType_USER_STATUS_TYPE_UNKNOWN
	}
}

// userToProto converts a store-level User to the proto UserType.
func userToProto(user *cognitostore.User) *pb.UserType {
	if user == nil {
		return nil
	}

	u := &pb.UserType{
		Username:   user.Username,
		Userstatus: userStatusToProto(user.UserStatus),
		Enabled:    proto.Bool(user.Enabled),
	}
	if !user.CreatedDate.IsZero() {
		u.Usercreatedate = user.CreatedDate.Format(timeutils.ISO8601UTCFormat)
	}
	if !user.LastModifiedDate.IsZero() {
		u.Userlastmodifieddate = user.LastModifiedDate.Format(timeutils.ISO8601UTCFormat)
	}
	if len(user.Attributes) > 0 {
		attrs := make([]*pb.AttributeType, 0, len(user.Attributes))
		for k, v := range user.Attributes {
			attrs = append(attrs, &pb.AttributeType{Name: k, Value: v})
		}
		u.Attributes = attrs
	}
	return u
}

// ---------------------------------------------------------------------------
// Group → proto
// ---------------------------------------------------------------------------

// groupToProto converts a store-level Group to the proto GroupType.
func groupToProto(group *cognitostore.Group) *pb.GroupType {
	if group == nil {
		return nil
	}

	result := &pb.GroupType{
		Groupname: group.Name,
		Rolearn:   group.RoleArn,
	}
	if group.Description != "" {
		result.Description = group.Description
	}
	if !group.CreationDate.IsZero() {
		result.Creationdate = group.CreationDate.Format(timeutils.ISO8601UTCFormat)
	}
	if !group.LastModifiedDate.IsZero() {
		result.Lastmodifieddate = group.LastModifiedDate.Format(timeutils.ISO8601UTCFormat)
	}
	if group.Precedence != nil {
		result.Precedence = proto.Int32(int32(*group.Precedence))
	}
	return result
}

// ---------------------------------------------------------------------------
// UserPoolClient → proto
// ---------------------------------------------------------------------------

func explicitAuthFlowToProto(s string) pb.ExplicitAuthFlowsType {
	switch s {
	case "ADMIN_NO_SRP_AUTH":
		return pb.ExplicitAuthFlowsType_EXPLICIT_AUTH_FLOWS_TYPE_ADMIN_NO_SRP_AUTH
	case "ADMIN_USER_PASSWORD_AUTH", "ALLOW_ADMIN_USER_PASSWORD_AUTH":
		return pb.ExplicitAuthFlowsType_EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_ADMIN_USER_PASSWORD_AUTH
	case "USER_SRP_AUTH", "ALLOW_USER_SRP_AUTH":
		return pb.ExplicitAuthFlowsType_EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_SRP_AUTH
	case "USER_PASSWORD_AUTH":
		return pb.ExplicitAuthFlowsType_EXPLICIT_AUTH_FLOWS_TYPE_USER_PASSWORD_AUTH
	case "ALLOW_USER_PASSWORD_AUTH":
		return pb.ExplicitAuthFlowsType_EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_PASSWORD_AUTH
	case "ALLOW_CUSTOM_AUTH":
		return pb.ExplicitAuthFlowsType_EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_CUSTOM_AUTH
	case "ALLOW_REFRESH_TOKEN_AUTH":
		return pb.ExplicitAuthFlowsType_EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_REFRESH_TOKEN_AUTH
	default:
		return 0
	}
}

func oauthFlowToProto(s string) pb.OAuthFlowType {
	switch s {
	case "code":
		return pb.OAuthFlowType_O_AUTH_FLOW_TYPE_CODE
	case "implicit":
		return pb.OAuthFlowType_O_AUTH_FLOW_TYPE_IMPLICIT
	case "client_credentials":
		return pb.OAuthFlowType_O_AUTH_FLOW_TYPE_CLIENT_CREDENTIALS
	default:
		return 0
	}
}

func preventUserExistenceErrorsToProto(s string) pb.PreventUserExistenceErrorTypes {
	if s == "ENABLED" {
		return pb.PreventUserExistenceErrorTypes_PREVENT_USER_EXISTENCE_ERROR_TYPES_ENABLED
	}
	return pb.PreventUserExistenceErrorTypes_PREVENT_USER_EXISTENCE_ERROR_TYPES_LEGACY
}

// userPoolClientToProto converts a store-level UserPoolClient to the proto
// UserPoolClientType. All scalar fields and enum arrays are converted.
func userPoolClientToProto(client *cognitostore.UserPoolClient) *pb.UserPoolClientType {
	if client == nil {
		return nil
	}

	result := &pb.UserPoolClientType{
		Clientid:                   client.ClientID,
		Clientname:                 client.ClientName,
		Userpoolid:                 client.UserPoolID,
		Clientsecret:               client.ClientSecret,
		Defaultredirecturi:         client.DefaultRedirectURI,
		Preventuserexistenceerrors: preventUserExistenceErrorsToProto(client.PreventUserExistenceErrors),
	}

	if !client.CreationDate.IsZero() {
		result.Creationdate = client.CreationDate.Format(timeutils.ISO8601UTCFormat)
	}
	if !client.LastModifiedDate.IsZero() {
		result.Lastmodifieddate = client.LastModifiedDate.Format(timeutils.ISO8601UTCFormat)
	}

	if client.RefreshTokenValidity > 0 {
		result.Refreshtokenvalidity = proto.Int32(int32(client.RefreshTokenValidity))
	}
	if client.AccessTokenValidity > 0 {
		result.Accesstokenvalidity = proto.Int32(int32(client.AccessTokenValidity))
	}
	if client.IDTokenValidity > 0 {
		result.Idtokenvalidity = proto.Int32(int32(client.IDTokenValidity))
	}
	if client.AuthSessionValidity > 0 {
		result.Authsessionvalidity = proto.Int32(int32(client.AuthSessionValidity))
	}

	if len(client.ExplicitAuthFlows) > 0 {
		result.Explicitauthflows = make([]pb.ExplicitAuthFlowsType, len(client.ExplicitAuthFlows))
		for i, f := range client.ExplicitAuthFlows {
			result.Explicitauthflows[i] = explicitAuthFlowToProto(f)
		}
	}
	if len(client.AllowedOAuthFlows) > 0 {
		result.Allowedoauthflows = make([]pb.OAuthFlowType, len(client.AllowedOAuthFlows))
		for i, f := range client.AllowedOAuthFlows {
			result.Allowedoauthflows[i] = oauthFlowToProto(f)
		}
	}

	result.Callbackurls = client.CallbackURLs
	result.Logouturls = client.LogoutURLs
	result.Allowedoauthscopes = client.AllowedOAuthScopes
	result.Supportedidentityproviders = client.SupportedIdentityProviders
	result.Readattributes = client.ReadAttributes
	result.Writeattributes = client.WriteAttributes

	result.Allowedoauthflowsuserpoolclient = proto.Bool(client.AllowedOAuthFlowsUserPoolClient)
	result.Enablepropagateadditionalusercontextdata = proto.Bool(client.EnablePropagateAdditionalUserContextData)
	result.Enabletokenrevocation = proto.Bool(client.EnableTokenRevocation)

	return result
}
