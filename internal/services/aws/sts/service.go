// Package sts provides STS (Security Token Service) operations for vorpalstacks.
package sts

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/iam"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/iam/policy"
	iamstore "vorpalstacks/internal/store/aws/iam"
	stsstore "vorpalstacks/internal/store/aws/sts"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// STSService provides AWS Security Token Service operations.
type STSService struct{}

// NewSTSService creates a new STS service instance.
func NewSTSService() *STSService {
	return &STSService{}
}

func (s *STSService) store(reqCtx *request.RequestContext) (stsstore.SessionStoreInterface, error) {
	if stsStore := reqCtx.GetSTSStore(); stsStore != nil {
		return stsStore, nil
	}
	storage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get global storage: %w", err)
	}
	region := reqCtx.GetRegion()
	return stsstore.NewSessionStore(storage, region), nil
}

func (s *STSService) iamStore(reqCtx *request.RequestContext) (iamstore.IAMStoreInterface, error) {
	storage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get global storage: %w", err)
	}
	return iamstore.NewIAMStore(storage, reqCtx.GetAccountID()), nil
}

// RegisterHandlers registers all STS operation handlers with the dispatcher.
func (s *STSService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("sts", "GetCallerIdentity", s.GetCallerIdentity)
	d.RegisterHandlerForService("sts", "AssumeRole", s.AssumeRole)
	d.RegisterHandlerForService("sts", "GetSessionToken", s.GetSessionToken)
	d.RegisterHandlerForService("sts", "AssumeRoleWithSAML", s.AssumeRoleWithSAML)
	d.RegisterHandlerForService("sts", "AssumeRoleWithWebIdentity", s.AssumeRoleWithWebIdentity)
	d.RegisterHandlerForService("sts", "AssumeRoot", s.AssumeRoot)
	d.RegisterHandlerForService("sts", "DecodeAuthorizationMessage", s.DecodeAuthorizationMessage)
	d.RegisterHandlerForService("sts", "GetAccessKeyInfo", s.GetAccessKeyInfo)
	d.RegisterHandlerForService("sts", "GetFederationToken", s.GetFederationToken)
	d.RegisterHandlerForService("sts", "GetDelegatedAccessToken", s.GetDelegatedAccessToken)
	d.RegisterHandlerForService("sts", "GetWebIdentityToken", s.GetWebIdentityToken)
}

// GetCallerIdentity returns details about the IAM user or role whose credentials are used to call the operation.
func (s *STSService) GetCallerIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	securityToken := req.Headers.Get("X-Amz-Security-Token")

	var userId, arn string

	if securityToken != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		session, err := store.Get(securityToken)
		if err == nil && session != nil {
			userId = session.AccessKeyId
			switch session.PrincipalType {
			case "AssumedRole":
				roleName := arnutil.ExtractRoleNameFromARN(session.RoleArn)
				arn = "arn:aws:sts::" + reqCtx.GetAccountID() + ":assumed-role/" + roleName + "/" + session.RoleSessionName
			case "SAML":
				arn = session.RoleArn + "/saml/" + session.PrincipalName
			case "WebIdentity":
				arn = session.RoleArn + "/assumed-role/" + session.RoleSessionName
			case "FederatedUser":
				arn = session.PrincipalArn + ":" + session.PrincipalName
			case "Root":
				arn = "arn:aws:iam::" + reqCtx.GetAccountID() + ":root"
			default:
				arn = session.PrincipalArn
			}
		}
	}

	if userId == "" && arn == "" {
		accessKeyId := req.Headers.Get("X-Amz-Access-Key")
		if accessKeyId != "" {
			iamStore, err := s.iamStore(reqCtx)
			if err == nil {
				accessKey, err := iamStore.AccessKeys().Get(accessKeyId)
				if err == nil && accessKey != nil {
					userId = accessKey.UserName
					arn = "arn:aws:iam::" + reqCtx.GetAccountID() + ":user/" + accessKey.UserName
				}
			}
		}
	}

	if userId == "" {
		userId = reqCtx.GetAccountID()
	}

	if arn == "" {
		arn = "arn:aws:iam::" + reqCtx.GetAccountID() + ":root"
	}

	return map[string]interface{}{
		"UserId":  userId,
		"Account": reqCtx.GetAccountID(),
		"Arn":     arn,
	}, nil
}

// AssumeRole returns a set of temporary security credentials that you can use to access AWS resources.
func (s *STSService) AssumeRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleArn := request.GetStringParam(req.Parameters, "RoleArn")
	roleSessionName := request.GetStringParam(req.Parameters, "RoleSessionName")
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")

	validDuration, err := validateDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if roleArn == "" {
		return nil, ErrInvalidRoleArn
	}

	if roleSessionName == "" {
		return nil, ErrInvalidParameter
	}

	roleName := arnutil.ExtractRoleNameFromARN(roleArn)
	if roleName == "" {
		return nil, ErrInvalidRoleArn
	}

	iamStore, err := s.iamStore(reqCtx)
	if err != nil {
		return nil, err
	}
	if !iamStore.Roles().Exists(roleName) {
		return nil, ErrNoSuchRole
	}

	trustPolicyDoc, err := iamStore.Roles().GetAssumeRolePolicyDocument(roleName)
	if err != nil {
		return nil, ErrNoSuchRole
	}

	parsedPolicy, err := policy.ParseDocument(trustPolicyDoc)
	if err != nil {
		return nil, ErrInvalidRoleArn
	}

	callerArn, _ := s.resolveCallerIdentity(reqCtx, req)
	if callerArn == "" {
		callerArn = "arn:aws:iam::" + reqCtx.GetAccountID() + ":root"
	}

	evalCtx := iam.BuildEvaluationContext(reqCtx.GetAccountID(), callerArn)
	if err := iam.EvaluateTrustPolicy(parsedPolicy, callerArn, evalCtx); err != nil {
		return nil, ErrAccessDenied
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create("AssumedRole", roleSessionName, roleArn, roleArn, roleSessionName, validDuration)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
		"AssumedRoleUser": map[string]interface{}{
			"AssumedRoleId": session.AccessKeyId + ":" + roleSessionName,
			"Arn":           "arn:aws:sts::" + reqCtx.GetAccountID() + ":assumed-role/" + roleName + "/" + roleSessionName,
		},
	}, nil
}

// GetSessionToken returns a set of temporary credentials for an AWS account or IAM user.
func (s *STSService) GetSessionToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")

	validDuration, err := validateDurationSecondsExtended(durationSeconds)
	if err != nil {
		return nil, err
	}

	accessKeyId := req.Headers.Get("X-Amz-Access-Key")
	if accessKeyId == "" {
		accessKeyId = request.GetStringParam(req.Parameters, "AccessKeyId")
	}

	var callerArn, callerName string

	if accessKeyId != "" {
		iamStore, err := s.iamStore(reqCtx)
		if err != nil {
			return nil, err
		}
		accessKey, err := iamStore.AccessKeys().Get(accessKeyId)
		if err == nil && accessKey != nil {
			callerName = accessKey.UserName
			callerArn = "arn:aws:iam::" + reqCtx.GetAccountID() + ":user/" + callerName
		}
	}

	if callerArn == "" {
		callerName = reqCtx.GetAccountID()
		callerArn = "arn:aws:iam::" + reqCtx.GetAccountID() + ":root"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create("User", callerName, callerArn, "", "", validDuration)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
	}, nil
}

func validateDurationSeconds(durationSeconds int) (int, error) {
	if durationSeconds == 0 {
		return DefaultDurationSeconds, nil
	}
	if durationSeconds < MinDurationSeconds || durationSeconds > MaxDurationSeconds {
		return 0, ErrInvalidDuration
	}
	return durationSeconds, nil
}

func (s *STSService) resolveCallerIdentity(reqCtx *request.RequestContext, req *request.ParsedRequest) (arn, name string) {
	accessKeyId := req.Headers.Get("X-Amz-Access-Key")
	if accessKeyId == "" {
		authHeader := req.Headers.Get("Authorization")
		if authHeader != "" {
			accessKeyId = request.ExtractAccessKeyIDFromAuth(authHeader)
		}
	}
	if accessKeyId == "" {
		return "", ""
	}
	iamStore, err := s.iamStore(reqCtx)
	if err != nil {
		return "", ""
	}
	accessKey, err := iamStore.AccessKeys().Get(accessKeyId)
	if err != nil || accessKey == nil {
		return "", ""
	}
	return "arn:aws:iam::" + reqCtx.GetAccountID() + ":user/" + accessKey.UserName, accessKey.UserName
}

func validateDurationSecondsExtended(durationSeconds int) (int, error) {
	if durationSeconds == 0 {
		return DefaultDurationSeconds, nil
	}
	if durationSeconds < MinDurationSeconds || durationSeconds > MaxDurationSecondsExtended {
		return 0, ErrInvalidDurationExtended
	}
	return durationSeconds, nil
}

// AssumeRoleWithSAML returns a set of temporary security credentials for users who have been authenticated via a SAML authentication response.
func (s *STSService) AssumeRoleWithSAML(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleArn := request.GetStringParam(req.Parameters, "RoleArn")
	principalArn := request.GetStringParam(req.Parameters, "PrincipalArn")
	samlAssertion := request.GetStringParam(req.Parameters, "SAMLAssertion")
	roleSessionName := request.GetStringParam(req.Parameters, "RoleSessionName")
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")

	if roleSessionName == "" {
		roleSessionName = "SAML"
	}

	validDuration, err := validateDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if roleArn == "" {
		return nil, ErrInvalidRoleArn
	}

	if principalArn == "" {
		return nil, ErrInvalidPrincipalArn
	}

	if samlAssertion == "" {
		return nil, ErrInvalidSAMLAssertion
	}

	roleName := arnutil.ExtractRoleNameFromARN(roleArn)
	if roleName == "" {
		return nil, ErrInvalidRoleArn
	}

	iamStore, err := s.iamStore(reqCtx)
	if err != nil {
		return nil, err
	}
	if !iamStore.Roles().Exists(roleName) {
		return nil, ErrNoSuchRole
	}

	trustPolicyDoc, err := iamStore.Roles().GetAssumeRolePolicyDocument(roleName)
	if err != nil {
		return nil, ErrNoSuchRole
	}

	parsedPolicy, err := policy.ParseDocument(trustPolicyDoc)
	if err != nil {
		return nil, ErrInvalidRoleArn
	}

	evalCtx := iam.BuildEvaluationContext(reqCtx.GetAccountID(), principalArn)
	if err := iam.EvaluateTrustPolicyForAction(parsedPolicy, principalArn, "sts:AssumeRoleWithSAML", evalCtx); err != nil {
		return nil, ErrAccessDenied
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create("SAML", principalArn, roleArn, roleArn, roleSessionName, validDuration)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
		"AssumedRoleUser": map[string]interface{}{
			"AssumedRoleId": session.AccessKeyId + ":" + roleSessionName,
			"Arn":           "arn:aws:sts::" + reqCtx.GetAccountID() + ":assumed-role/" + roleName + "/" + roleSessionName,
		},
		"Subject":       principalArn,
		"SubjectType":   "persistent",
		"Issuer":        "VorpalStacks",
		"NameQualifier": "SAML",
		"Audience":      "STS",
	}, nil
}

// AssumeRoleWithWebIdentity returns a set of temporary security credentials for users who have been authenticated in a mobile or web application with a web identity provider.
func (s *STSService) AssumeRoleWithWebIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleArn := request.GetStringParam(req.Parameters, "RoleArn")
	roleSessionName := request.GetStringParam(req.Parameters, "RoleSessionName")
	webIdentityToken := request.GetStringParam(req.Parameters, "WebIdentityToken")
	providerId := request.GetStringParam(req.Parameters, "ProviderId")
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")

	validDuration, err := validateDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if roleArn == "" {
		return nil, ErrInvalidRoleArn
	}

	if roleSessionName == "" {
		roleSessionName = "web-identity-session"
	}

	if webIdentityToken == "" {
		return nil, ErrInvalidWebIdentityToken
	}

	roleName := arnutil.ExtractRoleNameFromARN(roleArn)
	if roleName == "" {
		return nil, ErrInvalidRoleArn
	}

	iamStore, err := s.iamStore(reqCtx)
	if err != nil {
		return nil, err
	}
	if !iamStore.Roles().Exists(roleName) {
		return nil, ErrNoSuchRole
	}

	trustPolicyDoc, err := iamStore.Roles().GetAssumeRolePolicyDocument(roleName)
	if err != nil {
		return nil, ErrNoSuchRole
	}

	parsedPolicy, err := policy.ParseDocument(trustPolicyDoc)
	if err != nil {
		return nil, ErrInvalidRoleArn
	}

	federatedPrincipal := ""
	if providerId != "" {
		federatedPrincipal = "arn:aws:iam::" + reqCtx.GetAccountID() + ":oidc-provider/" + providerId
	}

	evalCtx := iam.BuildEvaluationContext(reqCtx.GetAccountID(), federatedPrincipal)
	if err := iam.EvaluateTrustPolicyForAction(parsedPolicy, federatedPrincipal, "sts:AssumeRoleWithWebIdentity", evalCtx); err != nil {
		return nil, ErrAccessDenied
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create("WebIdentity", roleSessionName, roleArn, roleArn, roleSessionName, validDuration)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
		"AssumedRoleUser": map[string]interface{}{
			"AssumedRoleId": session.AccessKeyId + ":" + roleSessionName,
			"Arn":           "arn:aws:sts::" + reqCtx.GetAccountID() + ":assumed-role/" + roleName + "/" + roleSessionName,
		},
		"Provider":                    providerId,
		"SubjectFromWebIdentityToken": roleSessionName,
		"Audience":                    "sts.amazonaws.com",
	}, nil
}

// AssumeRoot returns a set of temporary security credentials for the root user of an account.
func (s *STSService) AssumeRoot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")
	targetPrincipal := request.GetStringParam(req.Parameters, "TargetPrincipal")

	validDuration, err := validateDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create("Root", "root", "arn:aws:iam::"+reqCtx.GetAccountID()+":root", "", targetPrincipal, validDuration)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
	}, nil
}

// DecodeAuthorizationMessage decodes additional information about the authorization status of a request from an encoded message.
func (s *STSService) DecodeAuthorizationMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	encodedMessage := request.GetStringParam(req.Parameters, "EncodedMessage")

	if encodedMessage == "" {
		return nil, ErrInvalidEncodedMessage
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(encodedMessage)
	if err != nil {
		decodedBytes, err = base64.URLEncoding.DecodeString(encodedMessage)
		if err != nil {
			return nil, ErrInvalidEncodedMessage
		}
	}

	var decodedData map[string]interface{}
	if err := json.Unmarshal(decodedBytes, &decodedData); err != nil {
		return map[string]interface{}{
			"DecodedMessage": string(decodedBytes),
		}, nil
	}

	decodedJSON, err := json.Marshal(decodedData)
	if err != nil {
		return map[string]interface{}{
			"DecodedMessage": string(decodedBytes),
		}, nil
	}

	return map[string]interface{}{
		"DecodedMessage": string(decodedJSON),
	}, nil
}

// GetAccessKeyInfo returns information about the access key in the request.
func (s *STSService) GetAccessKeyInfo(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessKeyId := request.GetStringParam(req.Parameters, "AccessKeyId")

	if accessKeyId == "" {
		return nil, ErrInvalidAccessKeyId
	}

	if len(accessKeyId) >= 4 {
		prefix := accessKeyId[:4]
		if prefix == "ASIA" || prefix == "AKIA" || prefix == "ABIA" || prefix == "ACCA" {
			return map[string]interface{}{
				"Account": reqCtx.GetAccountID(),
			}, nil
		}
	}

	return map[string]interface{}{
		"Account": reqCtx.GetAccountID(),
	}, nil
}

// GetFederationToken returns a set of temporary security credentials for a federated user.
func (s *STSService) GetFederationToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	policy := request.GetStringParam(req.Parameters, "Policy")
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")

	validDuration, err := validateDurationSecondsExtended(durationSeconds)
	if err != nil {
		return nil, err
	}

	if name == "" {
		return nil, ErrInvalidFederationName
	}

	if policy != "" {
		var js interface{}
		if err := json.Unmarshal([]byte(policy), &js); err != nil {
			return nil, ErrMalformedPolicyDocument
		}
	}

	callerArn, callerName := s.resolveCallerIdentity(reqCtx, req)
	if callerArn == "" {
		callerArn = "arn:aws:iam::" + reqCtx.GetAccountID() + ":root"
	}
	if callerName == "" {
		callerName = reqCtx.GetAccountID()
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create("FederatedUser", name, callerArn, "", name, validDuration)
	if err != nil {
		return nil, err
	}

	const maxPolicySize = 2048

	packedPolicySize := 0
	if policy != "" {
		packedPolicySize = (len(policy) * 100) / maxPolicySize
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
		"FederatedUser": map[string]interface{}{
			"FederatedUserId": reqCtx.GetAccountID() + ":" + name,
			"Arn":             callerArn + ":" + name,
		},
		"PackedPolicySize": packedPolicySize,
	}, nil
}

// GetDelegatedAccessToken returns a set of temporary security credentials that represent an IAM identity centre user.
func (s *STSService) GetDelegatedAccessToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tradeInToken := request.GetStringParam(req.Parameters, "TradeInToken")

	if tradeInToken == "" {
		return nil, ErrInvalidTradeInToken
	}

	principalArn, principalName := s.resolveCallerIdentity(reqCtx, req)
	if principalArn == "" {
		principalArn = "arn:aws:iam::" + reqCtx.GetAccountID() + ":root"
	}
	if principalName == "" {
		principalName = reqCtx.GetAccountID()
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create("DelegatedAccess", principalName, principalArn, "", "", 3600)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
		"AssumedPrincipal": principalArn,
		"PackedPolicySize": 0,
	}, nil
}

// GetWebIdentityToken returns a web identity token for the caller.
func (s *STSService) GetWebIdentityToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")

	validDuration, err := validateDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	webIdentityToken := base64.StdEncoding.EncodeToString(tokenBytes)

	expiration := time.Now().UTC().Add(time.Duration(validDuration) * time.Second)

	return map[string]interface{}{
		"WebIdentityToken": webIdentityToken,
		"Expiration":       expiration.Format(timeutils.ISO8601SimpleFormat),
	}, nil
}
