package sts

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	stsstore "vorpalstacks/internal/store/aws/sts"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

func validateDurationSecondsExtended(durationSeconds int) (int, error) {
	if durationSeconds == 0 {
		return DefaultDurationSeconds, nil
	}
	if durationSeconds < MinDurationSeconds || durationSeconds > MaxDurationSecondsExtended {
		return 0, ErrInvalidDurationExtended
	}
	return durationSeconds, nil
}

// validateFederationDurationSeconds validates the DurationSeconds parameter
// for GetFederationToken. The Smithy durationSecondsType trait allows
// 900-129600 with a default of 43200 (12 hours). When the caller is the root
// user, AWS caps the session at 3600 seconds (1 hour); the isRoot parameter
// enforces this cap.
func validateFederationDurationSeconds(durationSeconds int, isRoot bool) (int, error) {
	if durationSeconds == 0 {
		if isRoot {
			return 3600, nil
		}
		return DefaultFederationDurationSeconds, nil
	}
	if durationSeconds < MinDurationSeconds || durationSeconds > MaxDurationSecondsExtended {
		return 0, ErrInvalidDurationExtended
	}
	if isRoot && durationSeconds > 3600 {
		return 0, ErrInvalidFederationRootDuration
	}
	return durationSeconds, nil
}

// federationNamePattern mirrors the Smithy userNameType trait used by
// GetFederationToken's Name parameter: [\w+=,.@-]* with length 2-32.
var federationNamePattern = regexp.MustCompile(`^[\w+=,.@-]+$`)

// GetSessionToken returns a set of temporary credentials for an AWS account or IAM user.
func (s *STSService) GetSessionToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")
	tokenCode := request.GetStringParam(req.Parameters, "TokenCode")

	validDuration, err := validateDurationSecondsExtended(durationSeconds)
	if err != nil {
		return nil, err
	}

	callerArn, callerName, err := s.resolveCallerArnOrReject(reqCtx, req)
	if err != nil {
		return nil, err
	}

	// AWS caps GetSessionToken duration at 3600 seconds for root users.
	isRoot := strings.HasSuffix(callerArn, ":root")
	if isRoot && validDuration > 3600 {
		return nil, ErrInvalidRootSessionDuration
	}

	if err := validateMFACredentials(serialNumber, tokenCode); err != nil {
		return nil, err
	}

	mfaPresent := false
	if serialNumber != "" && tokenCode != "" {
		iamStore, err := s.iamStore(reqCtx)
		if err != nil {
			return nil, err
		}
		actx := &assumeContext{
			SerialNumber: serialNumber,
			TokenCode:    tokenCode,
			CallerName:   callerName,
		}
		if _, err := s.verifyCallerMFA(reqCtx, iamStore, actx); err != nil {
			return nil, err
		}
		mfaPresent = true
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:          "User",
		PrincipalName:          callerName,
		PrincipalArn:           callerArn,
		DurationSeconds:        validDuration,
		MultiFactorAuthPresent: mfaPresent,
	})
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

// GetFederationToken returns a set of temporary security credentials for a federated user.
func (s *STSService) GetFederationToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	policy := request.GetStringParam(req.Parameters, "Policy")
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")

	if name == "" {
		return nil, ErrInvalidFederationName
	}
	// userNameType Smithy trait: length 2-32, pattern [\w+=,.@-]*.
	if len(name) < 2 || len(name) > 32 || !federationNamePattern.MatchString(name) {
		return nil, ErrInvalidFederationName
	}

	callerArn, _, err := s.resolveCallerArnOrReject(reqCtx, req)
	if err != nil {
		return nil, err
	}
	isRoot := strings.HasSuffix(callerArn, ":root")

	validDuration, err := validateFederationDurationSeconds(durationSeconds, isRoot)
	if err != nil {
		return nil, err
	}

	if err := validateSessionPolicy(policy); err != nil {
		return nil, err
	}

	fedTags, err := extractSessionTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	policyArns, err := extractPolicyArns(req.Parameters)
	if err != nil {
		return nil, err
	}

	packedPolicySize := computePackedPolicySize(policy, policyArns, nil, fedTags)
	if packedPolicySize > 100 {
		return nil, ErrPackedPolicyTooLarge
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "FederatedUser",
		PrincipalName:   name,
		PrincipalArn:    callerArn,
		RoleSessionName: name,
		DurationSeconds: validDuration,
		Tags:            fedTags,
		Policy:          policy,
		PolicyArns:      policyArns,
	})
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
		"FederatedUser": map[string]interface{}{
			"FederatedUserId": reqCtx.GetAccountID() + ":" + name,
			"Arn":             arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().FederatedUser(name),
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	principalArn, err := store.RedeemDelegationToken(tradeInToken)
	if err != nil {
		if errors.Is(err, stsstore.ErrDelegationTokenExpired) {
			return nil, ErrExpiredTradeInToken
		}
		if errors.Is(err, stsstore.ErrDelegationTokenNotFound) {
			return nil, ErrInvalidTradeInToken
		}
		// Unexpected errors (I/O, JSON unmarshal) indicate infrastructure
		// failure, not an invalid token. Masking them as InvalidToken
		// makes debugging impossible for the caller.
		return nil, ErrInternalError
	}

	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "DelegatedAccess",
		PrincipalName:   principalArn,
		PrincipalArn:    principalArn,
		DurationSeconds: 3600,
	})
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

// GetWebIdentityToken returns a signed JSON Web Token (JWT) representing the
// calling AWS identity. The returned JWT can be used to authenticate with
// external services that support OIDC discovery. The token is signed using
// the caller-specified algorithm (RS256 or ES384).
//
// AWS spec: https://docs.aws.amazon.com/STS/latest/APIReference/API_GetWebIdentityToken.html
func (s *STSService) GetWebIdentityToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")
	signingAlgorithm := request.GetStringParam(req.Parameters, "SigningAlgorithm")

	validDuration, err := validateWebIdentityDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if signingAlgorithm != "RS256" && signingAlgorithm != "ES384" {
		return nil, ErrInvalidSigningAlgorithm
	}

	var audiences []string
	for i := 1; ; i++ {
		key := fmt.Sprintf("Audience.member.%d", i)
		val := request.GetStringParam(req.Parameters, key)
		if val == "" {
			break
		}
		audiences = append(audiences, val)
	}
	if len(audiences) == 0 {
		return nil, ErrAudienceRequired
	}
	if len(audiences) > 10 {
		return nil, ErrTooManyAudiences
	}

	_, callerName, err := s.resolveCallerArnOrReject(reqCtx, req)
	if err != nil {
		return nil, err
	}

	// When the caller is using temporary credentials, the generated
	// web identity token must not outlive the caller's own session.
	// AWS rejects this with SessionDurationEscalationException (403).
	callerSession := s.resolveCallerSession(reqCtx, req)
	if callerSession != nil {
		remaining := int(time.Until(callerSession.Expiration).Seconds())
		if validDuration > remaining {
			return nil, ErrSessionDurationEscalation
		}
	}

	// Parse caller-supplied Tags and forward them to the JWT signer so
	// external services consuming the token can read the caller-attached
	// session tags from the standard "tags" JWT claim.
	tags, err := extractSessionTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	token, expiration, err := webIdentityTokenManagerInstance.generateToken(
		callerName, reqCtx.GetAccountID(), audiences, signingAlgorithm, validDuration, tags,
	)
	if err != nil {
		return nil, err
	}

	// Smithy GetWebIdentityToken can return JWTPayloadSizeExceededException
	// when the serialised JWT exceeds the platform limit. AWS enforces a
	// practical ceiling; we check the final token length as a guard.
	if len(token) > 32768 {
		return nil, ErrJWTPayloadSizeExceeded
	}

	return map[string]interface{}{
		"WebIdentityToken": token,
		"Expiration":       expiration.Format(timeutils.ISO8601SimpleFormat),
	}, nil
}
