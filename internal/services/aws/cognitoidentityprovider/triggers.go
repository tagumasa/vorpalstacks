package cognitoidentityprovider

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// TriggerVersion is the Cognito Lambda trigger payload version, matching
// the AWS Cognito User Pools trigger event format (currently version 1).
const TriggerVersion = 1

// PreSignUpSignUp is the trigger source constant for the PreSignUp Lambda
// trigger fired during user self-registration.
const PreSignUpSignUp = "PreSignUp_SignUp"

// PreSignUpAdminCreateUser is the trigger source constant for the PreSignUp
// Lambda trigger fired during administrative user creation.
const PreSignUpAdminCreateUser = "PreSignUp_AdminCreateUser"

// PostConfirmationConfirmSignUp is the trigger source constant for the
// PostConfirmation Lambda trigger fired after user self-confirmation.
const PostConfirmationConfirmSignUp = "PostConfirmation_ConfirmSignUp"

// PostConfirmationAdminCreateUser is the trigger source constant for the
// PostConfirmation Lambda trigger fired after administrative user creation.
const PostConfirmationAdminCreateUser = "PostConfirmation_AdminCreateUser"

// PreAuthentication is the trigger source constant for the Lambda trigger
// fired before user authentication.
const PreAuthentication = "PreAuthentication_Authentication"

// PostAuthentication is the trigger source constant for the Lambda trigger
// fired after successful user authentication.
const PostAuthentication = "PostAuthentication_Authentication"

// TokenGenerationAuthentication is the trigger source constant for the
// PreTokenGeneration Lambda trigger fired during user authentication.
const TokenGenerationAuthentication = "TokenGeneration_Authentication"

// TokenGenerationRefreshTokens is the trigger source constant for the
// PreTokenGeneration Lambda trigger fired during token refresh.
const TokenGenerationRefreshTokens = "TokenGeneration_RefreshTokens"

// TokenGenerationHostedAuth is the trigger source constant for the
// PreTokenGeneration Lambda trigger fired during hosted UI authentication.
const TokenGenerationHostedAuth = "TokenGeneration_HostedAuth"

// CustomMessageSignUp is the trigger source constant for the CustomMessage
// Lambda trigger fired during user self-registration.
const CustomMessageSignUp = "CustomMessage_SignUp"

// CustomMessageForgotPassword is the trigger source constant for the
// CustomMessage Lambda trigger fired during forgot password flow.
const CustomMessageForgotPassword = "CustomMessage_ForgotPassword"

// CustomMessageResendCode is the trigger source constant for the
// CustomMessage Lambda trigger fired when resending a confirmation code.
const CustomMessageResendCode = "CustomMessage_ResendCode"

// CustomMessageAdminCreateUser is the trigger source constant for the
// CustomMessage Lambda trigger fired during administrative user creation.
const CustomMessageAdminCreateUser = "CustomMessage_AdminCreateUser"

// CustomMessageVerifyUserAttribute is the trigger source constant for the
// CustomMessage Lambda trigger fired during user attribute verification.
const CustomMessageVerifyUserAttribute = "CustomMessage_VerifyUserAttribute"

// CustomMessageUpdateUserAttribute is the trigger source constant for the
// CustomMessage Lambda trigger fired during user attribute update.
const CustomMessageUpdateUserAttribute = "CustomMessage_UpdateUserAttribute"

// UserMigrationAuthentication is the trigger source constant for the
// UserMigration Lambda trigger fired during authentication when a user
// is not found in the user pool.
const UserMigrationAuthentication = "UserMigration_Authentication"

// DefineAuthChallenge is the trigger source constant for the Lambda trigger
// that defines which authentication challenge to present.
const DefineAuthChallenge = "DefineAuthChallenge_Authentication"

// CreateAuthChallenge is the trigger source constant for the Lambda trigger
// that creates a custom authentication challenge.
const CreateAuthChallenge = "CreateAuthChallenge_Authentication"

// VerifyAuthChallengeResponse is the trigger source constant for the Lambda
// trigger that verifies a user's response to an authentication challenge.
const VerifyAuthChallengeResponse = "VerifyAuthChallengeResponse_Authentication"

func defaultTriggerARNResolver(lambdaARN string) string {
	_, _, region, _, _ := svcarn.SplitARN(lambdaARN)
	return region
}

// handleCognitoTrigger is the bus handler for CognitoTriggerEvent. It
// extracts the function name from the Lambda ARN and invokes the Lambda
// function with the full AWS Cognito trigger payload.
func (s *CognitoService) handleCognitoTrigger(ctx context.Context, event *eventbus.CognitoTriggerEvent) eventbus.HandlerResult {
	if s.lambdaInvoker == nil || event.LambdaARN == "" {
		return eventbus.HandlerResult{StatusCode: 200}
	}

	_, payload, err := s.lambdaInvoker.InvokeForGateway(ctx, event.LambdaARN, event.Payload)
	if err != nil {
		logs.Error("cognito trigger Lambda invocation failed",
			logs.String("trigger_source", event.TriggerSource),
			logs.String("user_pool_id", event.UserPoolID),
			logs.String("username", event.Username),
			logs.String("lambda_arn", event.LambdaARN),
			logs.Err(err),
		)
		return eventbus.HandlerResult{StatusCode: 500, Error: err}
	}

	return eventbus.HandlerResult{StatusCode: 200, Payload: payload}
}

// invokeTrigger is the core trigger invocation function. It builds the
// AWS Cognito trigger payload, publishes a CognitoTriggerEvent via
// bus.PublishSync, and returns the Lambda response. When blocking is
// true, errors are returned to the caller; otherwise errors are logged
// and nil is returned.
func (s *CognitoService) invokeTrigger(
	ctx context.Context,
	triggerSource, userPoolID, username, clientID, lambdaARN string,
	requestPayload map[string]interface{},
	responseDefaults map[string]interface{},
	blocking bool,
) (map[string]interface{}, error) {
	if s.bus == nil || lambdaARN == "" {
		return nil, nil
	}

	region := defaultTriggerARNResolver(lambdaARN)

	triggerPayload := map[string]interface{}{
		"version":       TriggerVersion,
		"triggerSource": triggerSource,
		"region":        region,
		"userPoolId":    userPoolID,
		"userName":      username,
		"callerContext": map[string]interface{}{
			"awsSdkVersion": "",
			"clientId":      clientID,
		},
		"request":  requestPayload,
		"response": responseDefaults,
	}

	payloadBytes, err := json.Marshal(triggerPayload)
	if err != nil {
		if blocking {
			return nil, err
		}
		logs.Error("cognito trigger: failed to marshal trigger payload",
			logs.String("trigger_source", triggerSource),
			logs.Err(err),
		)
		return nil, nil
	}

	cte := &eventbus.CognitoTriggerEvent{
		EventBase: eventbus.EventBase{
			Region: region,
		},
		TriggerSource: triggerSource,
		UserPoolID:    userPoolID,
		Username:      username,
		ClientID:      clientID,
		LambdaARN:     lambdaARN,
		Version:       TriggerVersion,
		Payload:       payloadBytes,
	}

	result, busErr := s.bus.PublishSync(ctx, cte)
	if busErr != nil {
		if blocking {
			return nil, busErr
		}
		logs.Error("cognito trigger: bus PublishSync failed",
			logs.String("trigger_source", triggerSource),
			logs.Err(busErr),
		)
		return nil, nil
	}

	if result.Error != nil {
		if blocking {
			return nil, result.Error
		}
		logs.Error("cognito trigger: handler returned error",
			logs.String("trigger_source", triggerSource),
			logs.Err(result.Error),
		)
		return nil, nil
	}

	if len(result.Payload) == 0 {
		return responseDefaults, nil
	}

	var response map[string]interface{}
	if err := json.Unmarshal(result.Payload, &response); err != nil {
		if blocking {
			return nil, err
		}
		logs.Error("cognito trigger: failed to unmarshal trigger response",
			logs.String("trigger_source", triggerSource),
			logs.Err(err),
		)
		return responseDefaults, nil
	}

	return response, nil
}

// resolveTriggerARN maps a trigger source constant to the corresponding
// Lambda ARN stored in the user pool's LambdaConfig. Returns an empty
// string if the config is nil or the trigger source is unrecognised.
func resolveTriggerARN(config *cognitostore.LambdaConfig, triggerSource string) string {
	if config == nil {
		return ""
	}

	switch triggerSource {
	case PreSignUpSignUp, PreSignUpAdminCreateUser:
		return config.PreSignUp
	case PostConfirmationConfirmSignUp, PostConfirmationAdminCreateUser:
		return config.PostConfirmation
	case PreAuthentication:
		return config.PreAuthentication
	case PostAuthentication:
		return config.PostAuthentication
	case TokenGenerationAuthentication, TokenGenerationRefreshTokens, TokenGenerationHostedAuth:
		return config.PreTokenGeneration
	case CustomMessageSignUp, CustomMessageForgotPassword, CustomMessageResendCode,
		CustomMessageAdminCreateUser, CustomMessageVerifyUserAttribute, CustomMessageUpdateUserAttribute:
		return config.CustomMessage
	case UserMigrationAuthentication:
		return config.UserMigration
	case DefineAuthChallenge:
		return config.DefineAuthChallenge
	case CreateAuthChallenge:
		return config.CreateAuthChallenge
	case VerifyAuthChallengeResponse:
		return config.VerifyAuthChallengeResponse
	default:
		return ""
	}
}

// userAttributesMap builds a map of user attributes from a Cognito user,
// including the synthetic "sub" claim set to the user's internal ID.
func userAttributesMap(user *cognitostore.User) map[string]string {
	attrs := make(map[string]string)
	if user.Attributes != nil {
		for k, v := range user.Attributes {
			attrs[k] = v
		}
	}
	attrs["sub"] = user.ID
	return attrs
}

// preSignUpResult carries the parsed response from a PreSignUp Lambda trigger.
type preSignUpResult struct {
	AutoConfirmUser bool
	AutoVerifyEmail bool
	AutoVerifyPhone bool
	UserAttributes  map[string]string
}

// invokePreSignUp invokes the PreSignUp Lambda trigger synchronously. The
// trigger can auto-confirm the user, verify email/phone, and modify user
// attributes before creation. Returns an error if the trigger fails
// (blocking).
func invokePreSignUp(
	ctx context.Context,
	s *CognitoService,
	triggerSource, userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs map[string]string,
) (*preSignUpResult, error) {
	lambdaARN := resolveTriggerARN(config, triggerSource)
	if lambdaARN == "" {
		return &preSignUpResult{}, nil
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"validationData": nil,
		"clientMetadata": nil,
	}

	responseDefaults := map[string]interface{}{
		"autoConfirmUser": false,
		"autoVerifyEmail": false,
		"autoVerifyPhone": false,
	}

	response, err := s.invokeTrigger(ctx, triggerSource, userPoolID, username, clientID, lambdaARN, request, responseDefaults, true)
	if err != nil {
		return nil, err
	}

	result := &preSignUpResult{
		AutoConfirmUser: boolFromMap(response, "autoConfirmUser"),
		AutoVerifyEmail: boolFromMap(response, "autoVerifyEmail"),
		AutoVerifyPhone: boolFromMap(response, "autoVerifyPhone"),
	}

	if attrs, ok := response["userAttributes"].(map[string]interface{}); ok {
		merged := make(map[string]string)
		for k, v := range userAttrs {
			merged[k] = v
		}
		for k, v := range attrs {
			if vs, ok := v.(string); ok {
				merged[k] = vs
			}
		}
		result.UserAttributes = merged
	} else {
		result.UserAttributes = userAttrs
	}

	return result, nil
}

// invokePostConfirmation invokes the PostConfirmation Lambda trigger. This
// trigger is non-blocking: errors are logged but do not prevent the
// operation from succeeding.
func invokePostConfirmation(
	ctx context.Context,
	s *CognitoService,
	triggerSource, userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs map[string]string,
) error {
	lambdaARN := resolveTriggerARN(config, triggerSource)
	if lambdaARN == "" {
		return nil
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"clientMetadata": nil,
	}

	_, err := s.invokeTrigger(ctx, triggerSource, userPoolID, username, clientID, lambdaARN, request, map[string]interface{}{}, false)
	return err
}

// invokePreAuthentication invokes the PreAuthentication Lambda trigger
// before user authentication. This trigger is non-blocking: errors are
// logged but do not prevent authentication from proceeding.
func invokePreAuthentication(
	ctx context.Context,
	s *CognitoService,
	userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs map[string]string,
) error {
	lambdaARN := resolveTriggerARN(config, PreAuthentication)
	if lambdaARN == "" {
		return nil
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"clientMetadata": nil,
	}

	_, err := s.invokeTrigger(ctx, PreAuthentication, userPoolID, username, clientID, lambdaARN, request, map[string]interface{}{}, false)
	return err
}

// invokePostAuthentication invokes the PostAuthentication Lambda trigger
// after successful user authentication. This trigger is non-blocking.
func invokePostAuthentication(
	ctx context.Context,
	s *CognitoService,
	userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs map[string]string,
) error {
	lambdaARN := resolveTriggerARN(config, PostAuthentication)
	if lambdaARN == "" {
		return nil
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"newDeviceUsed":  false,
		"clientMetadata": nil,
	}

	_, err := s.invokeTrigger(ctx, PostAuthentication, userPoolID, username, clientID, lambdaARN, request, map[string]interface{}{}, false)
	return err
}

// preTokenGenerationResult carries the parsed response from a
// PreTokenGeneration Lambda trigger, including claim overrides and group
// overrides to apply to the generated tokens.
type preTokenGenerationResult struct {
	ClaimsToAddOrOverride map[string]string
	ClaimsToSuppress      []string
	GroupsToOverride      []string
	IAMRolesToOverride    []string
	PreferredRole         string
}

// invokePreTokenGeneration invokes the PreTokenGeneration Lambda trigger
// synchronously before token generation. The trigger can add/override/suppress
// claims and override groups in the generated tokens. Returns an error if
// the trigger fails (blocking).
func invokePreTokenGeneration(
	ctx context.Context,
	s *CognitoService,
	triggerSource, userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs map[string]string,
	userGroups []string,
) (*preTokenGenerationResult, error) {
	lambdaARN := resolveTriggerARN(config, triggerSource)
	if lambdaARN == "" {
		return &preTokenGenerationResult{}, nil
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"groupConfiguration": map[string]interface{}{
			"groupsToOverride":   nil,
			"iamRolesToOverride": nil,
			"preferredRole":      nil,
		},
		"clientMetadata": nil,
	}

	responseDefaults := map[string]interface{}{
		"claimsOverrideDetails": map[string]interface{}{
			"claimsToAddOrOverride": nil,
			"claimsToSuppress":      nil,
			"groupOverrideDetails": map[string]interface{}{
				"groupsToOverride":   nil,
				"iamRolesToOverride": nil,
				"preferredRole":      nil,
			},
		},
	}

	response, err := s.invokeTrigger(ctx, triggerSource, userPoolID, username, clientID, lambdaARN, request, responseDefaults, true)
	if err != nil {
		return nil, err
	}

	result := &preTokenGenerationResult{}

	if cod, ok := response["claimsOverrideDetails"].(map[string]interface{}); ok {
		if cta, ok := cod["claimsToAddOrOverride"].(map[string]interface{}); ok {
			m := make(map[string]string, len(cta))
			for k, v := range cta {
				if vs, ok := v.(string); ok {
					m[k] = vs
				}
			}
			result.ClaimsToAddOrOverride = m
		}

		if cts, ok := cod["claimsToSuppress"].([]interface{}); ok {
			slice := make([]string, 0, len(cts))
			for _, v := range cts {
				if vs, ok := v.(string); ok {
					slice = append(slice, vs)
				}
			}
			result.ClaimsToSuppress = slice
		}

		if god, ok := cod["groupOverrideDetails"].(map[string]interface{}); ok {
			if gto, ok := god["groupsToOverride"].([]interface{}); ok {
				result.GroupsToOverride = interfaceSliceToStrings(gto)
			}
			if iro, ok := god["iamRolesToOverride"].([]interface{}); ok {
				result.IAMRolesToOverride = interfaceSliceToStrings(iro)
			}
			if pr, ok := god["preferredRole"].(string); ok {
				result.PreferredRole = pr
			}
		}
	}

	return result, nil
}

// customMessageResult carries the parsed response from a CustomMessage
// Lambda trigger, containing the customised SMS message, email message,
// and email subject.
type customMessageResult struct {
	SMSMessage   string
	EmailMessage string
	EmailSubject string
}

// invokeCustomMessage invokes the CustomMessage Lambda trigger to allow
// customisation of the confirmation/forgot-password message content. This
// trigger is non-blocking: if it fails, the default message is used.
func invokeCustomMessage(
	ctx context.Context,
	s *CognitoService,
	triggerSource, userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	codeParameter string,
	userAttrs map[string]string,
) (*customMessageResult, error) {
	lambdaARN := resolveTriggerARN(config, triggerSource)
	if lambdaARN == "" {
		return nil, nil
	}

	request := map[string]interface{}{
		"codeParameter":     codeParameter,
		"usernameParameter": nil,
		"userAttributes":    userAttrs,
		"clientMetadata":    nil,
	}

	responseDefaults := map[string]interface{}{
		"smsMessage":   nil,
		"emailMessage": nil,
		"emailSubject": nil,
	}

	response, err := s.invokeTrigger(ctx, triggerSource, userPoolID, username, clientID, lambdaARN, request, responseDefaults, false)
	if err != nil {
		return nil, err
	}

	result := &customMessageResult{
		SMSMessage:   stringFromMap(response, "smsMessage"),
		EmailMessage: stringFromMap(response, "emailMessage"),
		EmailSubject: stringFromMap(response, "emailSubject"),
	}

	return result, nil
}

// userMigrationResult carries the parsed response from a UserMigration
// Lambda trigger, containing the migrated user attributes and desired
// initial state.
type userMigrationResult struct {
	UserAttributes         map[string]string
	FinalUserStatus        string
	MessageAction          string
	DesiredDeliveryMediums []string
	ForceAliasCreation     bool
}

// invokeUserMigration invokes the UserMigration Lambda trigger
// synchronously when a user is not found during authentication. The trigger
// can create the user from an external identity store. Returns an error if
// the trigger fails (blocking).
func invokeUserMigration(
	ctx context.Context,
	s *CognitoService,
	userPoolID, username, clientID, password string,
	config *cognitostore.LambdaConfig,
) (*userMigrationResult, error) {
	lambdaARN := resolveTriggerARN(config, UserMigrationAuthentication)
	if lambdaARN == "" {
		return nil, nil
	}

	request := map[string]interface{}{
		"password":       password,
		"validationData": nil,
		"clientMetadata": nil,
	}

	responseDefaults := map[string]interface{}{
		"userAttributes":         nil,
		"finalUserStatus":        "CONFIRMED",
		"messageAction":          "SUPPRESS",
		"desiredDeliveryMediums": nil,
		"forceAliasCreation":     false,
	}

	response, err := s.invokeTrigger(ctx, UserMigrationAuthentication, userPoolID, username, clientID, lambdaARN, request, responseDefaults, true)
	if err != nil {
		return nil, err
	}

	result := &userMigrationResult{
		FinalUserStatus:    stringFromMap(response, "finalUserStatus"),
		MessageAction:      stringFromMap(response, "messageAction"),
		ForceAliasCreation: boolFromMap(response, "forceAliasCreation"),
	}

	if result.FinalUserStatus == "" {
		result.FinalUserStatus = "CONFIRMED"
	}

	if ua, ok := response["userAttributes"].(map[string]interface{}); ok {
		m := make(map[string]string, len(ua))
		for k, v := range ua {
			if vs, ok := v.(string); ok {
				m[k] = vs
			}
		}
		result.UserAttributes = m
	}

	if ddm, ok := response["desiredDeliveryMediums"].([]interface{}); ok {
		result.DesiredDeliveryMediums = interfaceSliceToStrings(ddm)
	}

	return result, nil
}

// defineAuthChallengeResult carries the parsed response from a
// DefineAuthChallenge Lambda trigger, specifying which challenge to
// present or whether to issue tokens directly.
type defineAuthChallengeResult struct {
	ChallengeName      string
	IssueTokens        bool
	FailAuthentication bool
}

// invokeDefineAuthChallenge invokes the DefineAuthChallenge Lambda trigger
// synchronously. The trigger determines which authentication challenge to
// present, or whether to issue tokens immediately. Blocking.
func invokeDefineAuthChallenge(
	ctx context.Context,
	s *CognitoService,
	userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs map[string]string,
	session []map[string]interface{},
) (*defineAuthChallengeResult, error) {
	lambdaARN := resolveTriggerARN(config, DefineAuthChallenge)
	if lambdaARN == "" {
		return nil, nil
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"session":        session,
		"clientMetadata": nil,
	}

	responseDefaults := map[string]interface{}{
		"challengeName":      nil,
		"issueTokens":        false,
		"failAuthentication": false,
	}

	response, err := s.invokeTrigger(ctx, DefineAuthChallenge, userPoolID, username, clientID, lambdaARN, request, responseDefaults, true)
	if err != nil {
		return nil, err
	}

	result := &defineAuthChallengeResult{
		ChallengeName:      stringFromMap(response, "challengeName"),
		IssueTokens:        boolFromMap(response, "issueTokens"),
		FailAuthentication: boolFromMap(response, "failAuthentication"),
	}

	return result, nil
}

// createAuthChallengeResult carries the parsed response from a
// CreateAuthChallenge Lambda trigger, including public and private
// challenge parameters.
type createAuthChallengeResult struct {
	PublicChallengeParameters  map[string]string
	PrivateChallengeParameters map[string]string
	ChallengeMetadata          string
}

// invokeCreateAuthChallenge invokes the CreateAuthChallenge Lambda trigger
// synchronously. The trigger generates the challenge parameters for a
// custom authentication flow. Blocking.
func invokeCreateAuthChallenge(
	ctx context.Context,
	s *CognitoService,
	userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs map[string]string,
	challengeName string,
	session []map[string]interface{},
) (*createAuthChallengeResult, error) {
	lambdaARN := resolveTriggerARN(config, CreateAuthChallenge)
	if lambdaARN == "" {
		return nil, nil
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"challengeName":  challengeName,
		"session":        session,
		"clientMetadata": nil,
	}

	responseDefaults := map[string]interface{}{
		"publicChallengeParameters":  nil,
		"privateChallengeParameters": nil,
		"challengeMetadata":          nil,
	}

	response, err := s.invokeTrigger(ctx, CreateAuthChallenge, userPoolID, username, clientID, lambdaARN, request, responseDefaults, true)
	if err != nil {
		return nil, err
	}

	result := &createAuthChallengeResult{
		ChallengeMetadata: stringFromMap(response, "challengeMetadata"),
	}

	if pcp, ok := response["publicChallengeParameters"].(map[string]interface{}); ok {
		m := make(map[string]string, len(pcp))
		for k, v := range pcp {
			if vs, ok := v.(string); ok {
				m[k] = vs
			}
		}
		result.PublicChallengeParameters = m
	}

	if prcp, ok := response["privateChallengeParameters"].(map[string]interface{}); ok {
		m := make(map[string]string, len(prcp))
		for k, v := range prcp {
			if vs, ok := v.(string); ok {
				m[k] = vs
			}
		}
		result.PrivateChallengeParameters = m
	}

	return result, nil
}

// verifyAuthChallengeResponseResult carries the parsed response from a
// VerifyAuthChallengeResponse Lambda trigger, indicating whether the
// user's challenge answer is correct.
type verifyAuthChallengeResponseResult struct {
	AnswerCorrect bool
}

// invokeVerifyAuthChallengeResponse invokes the VerifyAuthChallengeResponse
// Lambda trigger synchronously. The trigger validates the user's answer to
// a custom authentication challenge. Blocking.
func invokeVerifyAuthChallengeResponse(
	ctx context.Context,
	s *CognitoService,
	userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs map[string]string,
	challengeName string,
	challengeResponse map[string]interface{},
) (*verifyAuthChallengeResponseResult, error) {
	lambdaARN := resolveTriggerARN(config, VerifyAuthChallengeResponse)
	if lambdaARN == "" {
		return nil, nil
	}

	request := map[string]interface{}{
		"userAttributes":    userAttrs,
		"challengeName":     challengeName,
		"challengeResponse": challengeResponse,
		"clientMetadata":    nil,
	}

	responseDefaults := map[string]interface{}{
		"answerCorrect": false,
	}

	response, err := s.invokeTrigger(ctx, VerifyAuthChallengeResponse, userPoolID, username, clientID, lambdaARN, request, responseDefaults, true)
	if err != nil {
		return nil, err
	}

	result := &verifyAuthChallengeResponseResult{
		AnswerCorrect: boolFromMap(response, "answerCorrect"),
	}

	return result, nil
}

func boolFromMap(m map[string]interface{}, key string) bool {
	if m == nil {
		return false
	}
	v, ok := m[key]
	if !ok {
		return false
	}
	switch b := v.(type) {
	case bool:
		return b
	case string:
		return b == "true"
	}
	return false
}

func stringFromMap(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func interfaceSliceToStrings(slice []interface{}) []string {
	if slice == nil {
		return nil
	}
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
