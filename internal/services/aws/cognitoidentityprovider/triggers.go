package cognitoidentityprovider

import (
	"context"
	"encoding/json"
	"errors"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// TriggerVersion is the Cognito Lambda trigger payload version, matching
// the AWS Cognito User Pools trigger event format (currently version 1).
const TriggerVersion = 1

// triggerCallerSdkVersion identifies the platform component that invokes the
// trigger, occupying the awsSdkVersion field of the callerContext that AWS
// fills with its own internal SDK version.
const triggerCallerSdkVersion = "vorpalstacks/1.0"

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

// TokenGenerationClientCredentials is the trigger source constant for the
// PreTokenGeneration Lambda trigger fired after a machine-to-machine client
// credentials grant. The user pool sends this event only when the trigger is
// configured at event version V3_0.
const TokenGenerationClientCredentials = "TokenGeneration_ClientCredentials"

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
	if s.bus == nil || event.LambdaARN == "" {
		return eventbus.HandlerResult{StatusCode: 200}
	}

	invocation, err := s.bus.LambdaInvoker().InvokeForTrigger(ctx, event.LambdaARN, event.Payload)
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

	// A raised function error still means the invoke transport succeeded;
	// surface it distinctly so trigger consumers can classify it as a
	// function-level outcome rather than an infrastructure failure.
	if invocation.FunctionError != "" {
		return eventbus.HandlerResult{StatusCode: 200, Error: &lambdaFunctionError{msg: invocation.FunctionError}}
	}

	return eventbus.HandlerResult{StatusCode: 200, Payload: invocation.Payload}
}

// lambdaFunctionError marks a trigger invocation where the Lambda function
// itself raised an error, as opposed to an invocation-transport failure.
type lambdaFunctionError struct {
	msg string
}

func (e *lambdaFunctionError) Error() string { return e.msg }

// classifyMigrationFailure maps a UserMigration trigger failure onto the
// AWS contract: a Lambda function that raises an error fails the sign-in
// with NotAuthorizedException (the documented way to reject a bad password
// during migration), while an invocation-transport failure is an
// infrastructure error.
func classifyMigrationFailure(err error) error {
	var fnErr *lambdaFunctionError
	if errors.As(err, &fnErr) {
		return ErrIncorrectPassword
	}
	return ErrInternalError
}

// classifyTriggerFailure maps an authentication-path trigger failure onto
// the AWS contract: a Lambda function that raises an error fails the
// authentication with UserLambdaValidationException, while an
// invocation-transport failure is an infrastructure error.
func classifyTriggerFailure(err error) error {
	var fnErr *lambdaFunctionError
	if errors.As(err, &fnErr) {
		return ErrUserLambdaValidation
	}
	return ErrInternalError
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
			"awsSdkVersion": triggerCallerSdkVersion,
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
			// A response that is not a JSON object is malformed trigger
			// output: the documented InvalidLambdaResponseException, not a
			// raw encoding error that maps onto an internal failure.
			return nil, ErrInvalidLambdaResponse
		}
		logs.Error("cognito trigger: failed to unmarshal trigger response",
			logs.String("trigger_source", triggerSource),
			logs.Err(err),
		)
		return nil, nil
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
	case TokenGenerationClientCredentials:
		// The M2M source is configured through the versioned trigger form:
		// resolve its ARN first, falling back to the plain-ARN member.
		if config.PreTokenGenerationConfig != nil && config.PreTokenGenerationConfig.LambdaArn != "" {
			return config.PreTokenGenerationConfig.LambdaArn
		}
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
	validationData map[string]string,
	clientMetadata map[string]string,
) (*preSignUpResult, error) {
	lambdaARN := resolveTriggerARN(config, triggerSource)
	if lambdaARN == "" {
		return &preSignUpResult{}, nil
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"validationData": validationData,
		"clientMetadata": clientMetadata,
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
		// Only the trigger's own overrides travel here; the caller merges
		// them onto the validated client map so the two attribute sources
		// stay distinguishable.
		merged := make(map[string]string, len(attrs))
		for k, v := range attrs {
			if vs, ok := v.(string); ok {
				merged[k] = vs
			}
		}
		result.UserAttributes = merged
	}

	return result, nil
}

// invokePostConfirmation invokes the PostConfirmation Lambda trigger. This
// trigger is non-blocking by contract: every failure is logged inside
// invokeTrigger and never prevents the confirmed operation from succeeding,
// so the invoker returns nothing — there is no error to guard on at the
// call sites.
func invokePostConfirmation(
	ctx context.Context,
	s *CognitoService,
	triggerSource, userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs map[string]string,
) {
	lambdaARN := resolveTriggerARN(config, triggerSource)
	if lambdaARN == "" {
		return
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"clientMetadata": nil,
	}

	s.invokeTrigger(ctx, triggerSource, userPoolID, username, clientID, lambdaARN, request, map[string]interface{}{}, false)
}

// invokePreAuthentication invokes the PreAuthentication Lambda trigger
// before user authentication. The trigger is blocking: raising an error is
// the documented way to deny the sign-in, so a function-level failure fails
// the authentication with UserLambdaValidationException and a transport
// failure with InternalErrorException.
func invokePreAuthentication(
	ctx context.Context,
	s *CognitoService,
	userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs, clientMetadata map[string]string,
) error {
	lambdaARN := resolveTriggerARN(config, PreAuthentication)
	if lambdaARN == "" {
		return nil
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"clientMetadata": clientMetadata,
	}

	_, err := s.invokeTrigger(ctx, PreAuthentication, userPoolID, username, clientID, lambdaARN, request, map[string]interface{}{}, true)
	if err != nil {
		return classifyTriggerFailure(err)
	}
	return nil
}

// invokePostAuthentication invokes the PostAuthentication Lambda trigger
// after successful user authentication. Per AWS docs, a trigger failure
// causes the authentication to fail and returns an error to the client.
func invokePostAuthentication(
	ctx context.Context,
	s *CognitoService,
	userPoolID, username, clientID string,
	config *cognitostore.LambdaConfig,
	userAttrs, clientMetadata map[string]string,
) error {
	lambdaARN := resolveTriggerARN(config, PostAuthentication)
	if lambdaARN == "" {
		return nil
	}

	request := map[string]interface{}{
		"userAttributes": userAttrs,
		"newDeviceUsed":  false,
		"clientMetadata": clientMetadata,
	}

	_, err := s.invokeTrigger(ctx, PostAuthentication, userPoolID, username, clientID, lambdaARN, request, map[string]interface{}{}, true)
	if err != nil {
		return classifyTriggerFailure(err)
	}
	return nil
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
	clientMetadata map[string]string,
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
		"clientMetadata": clientMetadata,
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
	clientMetadata map[string]string,
) (*customMessageResult, error) {
	lambdaARN := resolveTriggerARN(config, triggerSource)
	if lambdaARN == "" {
		return nil, nil
	}

	request := map[string]interface{}{
		"codeParameter":     codeParameter,
		"usernameParameter": nil,
		"userAttributes":    userAttrs,
		"clientMetadata":    clientMetadata,
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
// Lambda trigger. FinalUserStatus and DesiredDeliveryMediums arrive already
// resolved to their documented defaults (RESET_REQUIRED and SMS); an empty
// MessageAction means the welcome message is sent. Only members the
// migration path consumes are parsed: delivery mediums name the welcome
// message's channel and alias-creation semantics belong to the alias
// machinery, carried for it.
type userMigrationResult struct {
	UserAttributes         map[string]string
	FinalUserStatus        string
	MessageAction          string
	ForceAliasCreation     bool
	DesiredDeliveryMediums []string
	EnableSMSMFA           bool
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
	validationData, clientMetadata map[string]string,
) (*userMigrationResult, error) {
	lambdaARN := resolveTriggerARN(config, UserMigrationAuthentication)
	if lambdaARN == "" {
		return nil, nil
	}

	request := map[string]interface{}{
		"password":       password,
		"validationData": validationData,
		"clientMetadata": clientMetadata,
	}

	responseDefaults := map[string]interface{}{
		"userAttributes": nil,
	}

	response, err := s.invokeTrigger(ctx, UserMigrationAuthentication, userPoolID, username, clientID, lambdaARN, request, responseDefaults, true)
	if err != nil {
		return nil, err
	}

	result := &userMigrationResult{
		FinalUserStatus:    stringFromMap(response, "finalUserStatus"),
		MessageAction:      stringFromMap(response, "messageAction"),
		ForceAliasCreation: boolFromMap(response, "forceAliasCreation"),
		EnableSMSMFA:       boolFromMap(response, "enableSMSMFA"),
	}

	// finalUserStatus and messageAction carry fixed vocabularies: the
	// documented migration statuses are CONFIRMED, RESET_REQUIRED and
	// FORCE_CHANGE_PASSWORD, and the message actions are SUPPRESS and
	// RESEND. Anything else is a malformed trigger response — a stray
	// value must never become a stored user status.
	switch result.FinalUserStatus {
	case "", "CONFIRMED", "RESET_REQUIRED", "FORCE_CHANGE_PASSWORD":
	default:
		return nil, ErrInvalidLambdaResponse
	}
	switch result.MessageAction {
	case "", "SUPPRESS", "RESEND":
	default:
		return nil, ErrInvalidLambdaResponse
	}

	// The documented defaults: an omitted finalUserStatus resolves to
	// RESET_REQUIRED ("If you don't set this attribute to CONFIRMED, it's
	// set to RESET_REQUIRED") and an omitted desiredDeliveryMediums sends
	// the welcome message by SMS.
	if result.FinalUserStatus == "" {
		result.FinalUserStatus = "RESET_REQUIRED"
	}
	if raw, ok := response["desiredDeliveryMediums"].([]interface{}); ok {
		if len(raw) == 0 {
			return nil, ErrInvalidLambdaResponse
		}
		mediums := make([]string, 0, len(raw))
		for _, v := range raw {
			s, isString := v.(string)
			if !isString || (s != "EMAIL" && s != "SMS") {
				return nil, ErrInvalidLambdaResponse
			}
			mediums = append(mediums, s)
		}
		result.DesiredDeliveryMediums = mediums
	} else {
		result.DesiredDeliveryMediums = []string{"SMS"}
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
