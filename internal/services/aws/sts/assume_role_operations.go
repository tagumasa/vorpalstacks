package sts

import (
	"context"

	"vorpalstacks/internal/common/request"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// AssumeRole returns a set of temporary security credentials that you can use to access AWS resources.
func (s *STSService) AssumeRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.assumeRoleCore(reqCtx, WireInput{
		Parameters:    req.Parameters,
		AccessKeyID:   callerAccessKeyID(req),
		SecurityToken: req.Headers.Get("X-Amz-Security-Token"),
	})
	if err != nil {
		return nil, err
	}

	return withSourceIdentity(map[string]interface{}{
		"Credentials": credentialsMap(result.Credentials),
		"AssumedRoleUser": map[string]interface{}{
			"AssumedRoleId": result.RoleID + ":" + result.RoleSessionName,
			"Arn":           arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(result.RoleName, result.RoleSessionName),
		},
		"PackedPolicySize": result.PackedPolicySize,
	}, result.SourceIdentity), nil
}

// AssumeRoleWithSAML returns a set of temporary security credentials for users who have been authenticated via a SAML authentication response.
// VorpalStacks cannot validate SAML assertions against external IdPs, so this is only available in TEST_MODE.
func (s *STSService) AssumeRoleWithSAML(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.assumeRoleWithSAMLCore(reqCtx, WireInput{Parameters: req.Parameters})
	if err != nil {
		return nil, err
	}

	return withSourceIdentity(map[string]interface{}{
		"Credentials": credentialsMap(result.Credentials),
		"AssumedRoleUser": map[string]interface{}{
			"AssumedRoleId": result.RoleID + ":" + result.RoleSessionName,
			"Arn":           arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(result.RoleName, result.RoleSessionName),
		},
		"Subject":          result.Subject,
		"SubjectType":      "persistent",
		"Issuer":           result.Issuer,
		"NameQualifier":    "SAML",
		"Audience":         result.Audience,
		"PackedPolicySize": result.PackedPolicySize,
		// SourceIdentity is optional in the Smithy
		// AssumeRoleWithSAMLResponse shape. AWS derives the value from the
		// SAML assertion's saml:AttributeStatement, which VorpalStacks
		// does not parse in TEST_MODE; withSourceIdentity omits the field
		// until a real SAML parser is introduced.
	}, ""), nil
}

// AssumeRoleWithWebIdentity returns a set of temporary security credentials for users who have been authenticated in a mobile or web application with a web identity provider.
// VorpalStacks cannot validate web identity tokens against external IdPs, so this is only available in TEST_MODE.
func (s *STSService) AssumeRoleWithWebIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.assumeRoleWithWebIdentityCore(reqCtx, WireInput{Parameters: req.Parameters})
	if err != nil {
		return nil, err
	}

	return withSourceIdentity(map[string]interface{}{
		"Credentials": credentialsMap(result.Credentials),
		"AssumedRoleUser": map[string]interface{}{
			"AssumedRoleId": result.RoleID + ":" + result.RoleSessionName,
			"Arn":           arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(result.RoleName, result.RoleSessionName),
		},
		"Provider":                    result.Provider,
		"SubjectFromWebIdentityToken": result.Subject,
		"Audience":                    result.Audience,
		"PackedPolicySize":            result.PackedPolicySize,
	}, result.SourceIdentity), nil
}

// AssumeRoot returns a set of temporary security credentials for performing
// privileged tasks on a member account. AWS requires the caller to be an
// IAM user or role in the Organizations management account (or an IAM
// delegated administrator) with an explicit sts:AssumeRoot grant, and
// explicitly forbids calling AssumeRoot with root user credentials. VorpalStacks
// does not implement Organizations (see docs/services.md "No organisations
// integration"), so the caller-side check reduces to the standard IAM policy
// evaluation for sts:AssumeRoot plus the root-caller rejection below. The
// session itself is scoped by the task policy: per the AWS AssumeRoot
// contract, TaskPolicyArn restricts the temporary credentials to the
// privileged task's action set instead of granting unrestricted root.
func (s *STSService) AssumeRoot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.assumeRootCore(reqCtx, WireInput{Parameters: req.Parameters})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": credentialsMap(*result),
	}, nil
}
