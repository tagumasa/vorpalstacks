package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// Neutral values from the Smithy EventRiskType enums, used for every
// recorded authentication event in the absence of a risk engine.
const (
	riskDecisionNoRisk = "NoRisk"
	riskLevelLow       = "Low"
)

// AdminListUserAuthEvents lists authentication events for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminListUserAuthEvents.html
func (s *CognitoService) AdminListUserAuthEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Smithy QueryLimitType: range {min: 0, max: 60}
	maxResults, err := parseListLimit(req.Parameters, "MaxResults", 60)
	if err != nil {
		return nil, err
	}

	result, err := store.ListAuthEventsPaginated(userPoolID, user.ID, storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   req.GetParam("NextToken"),
	})
	if err != nil {
		return nil, ErrInternalError
	}

	formatted := make([]map[string]interface{}, 0, len(result.Items))
	for _, e := range result.Items {
		formatted = append(formatted, formatAuthEvent(e))
	}

	resp := map[string]interface{}{
		"AuthEvents": formatted,
	}
	if result.IsTruncated && result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// AdminUpdateAuthEventFeedback updates feedback for an auth event (admin).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminUpdateAuthEventFeedback.html
func (s *CognitoService) AdminUpdateAuthEventFeedback(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	username := getUsername(req)
	eventID := req.GetParam("EventId")
	feedbackValue := req.GetParam("FeedbackValue")
	if userPoolID == "" || username == "" || eventID == "" || feedbackValue == "" {
		return nil, ErrInvalidParameter
	}
	if feedbackValue != "Valid" && feedbackValue != "Invalid" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	event, err := store.GetAuthEvent(userPoolID, user.ID, eventID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	event.FeedbackValue = feedbackValue
	event.FeedbackProvider = "admin"
	event.FeedbackDate = time.Now().UTC()
	if err := store.UpdateAuthEvent(event); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// UpdateAuthEventFeedback updates feedback for an auth event (user via access token).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateAuthEventFeedback.html
func (s *CognitoService) UpdateAuthEventFeedback(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	eventID := req.GetParam("EventId")
	feedbackValue := req.GetParam("FeedbackValue")
	if accessToken == "" || eventID == "" || feedbackValue == "" {
		return nil, ErrInvalidParameter
	}
	if feedbackValue != "Valid" && feedbackValue != "Invalid" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	event, err := store.GetAuthEvent(user.UserPoolID, user.ID, eventID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	event.FeedbackValue = feedbackValue
	event.FeedbackProvider = "user"
	event.FeedbackDate = time.Now().UTC()
	if err := store.UpdateAuthEvent(event); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// recordAuthEvent creates and stores an authentication event. Called from
// authenticateUser and other auth flows.
//
// Every event is recorded with the neutral risk assessment values from the
// Smithy EventRiskType enums (RiskDecision NoRisk, RiskLevel Low): this
// platform does not run a threat-protection risk engine, so no event is ever
// escalated to a riskier decision or level.
func (s *CognitoService) recordAuthEvent(reqCtx *request.RequestContext, userPoolID, userID, username, clientID, eventType, eventResponse string) {
	store, err := s.store(reqCtx)
	if err != nil {
		return
	}

	event := &cognitostore.AuthEvent{
		EventID:       generateEventID(),
		UserName:      username,
		ClientID:      clientID,
		UserPoolID:    userPoolID,
		UserID:        userID,
		EventType:     eventType,
		CreationDate:  time.Now().UTC(),
		EventResponse: eventResponse,
		RiskDecision:  riskDecisionNoRisk,
		RiskLevel:     riskLevelLow,
	}

	if err := store.CreateAuthEvent(event); err != nil {
		return
	}

	s.publishAuthEventLog(reqCtx, userPoolID, event)
}

func generateEventID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("evt-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", b)
}

func formatAuthEvent(e *cognitostore.AuthEvent) map[string]interface{} {
	result := map[string]interface{}{
		"EventId":       e.EventID,
		"EventType":     e.EventType,
		"CreationDate":  e.CreationDate.Unix(),
		"EventResponse": e.EventResponse,
		"EventRisk": map[string]interface{}{
			"RiskDecision": e.RiskDecision,
			"RiskLevel":    e.RiskLevel,
		},
	}
	if e.CompromisedFlag {
		result["EventRisk"].(map[string]interface{})["CompromisedCredentialsDetected"] = true
	}
	if len(e.ChallengeResponses) > 0 {
		challenges := make([]map[string]string, 0, len(e.ChallengeResponses))
		for _, c := range e.ChallengeResponses {
			challenges = append(challenges, map[string]string{
				"ChallengeName":     c.ChallengeName,
				"ChallengeResponse": c.ChallengeResponse,
			})
		}
		result["ChallengeResponses"] = challenges
	}
	ctxData := map[string]interface{}{}
	if e.ContextIPAddress != "" {
		ctxData["IpAddress"] = e.ContextIPAddress
	}
	if e.ContextDeviceName != "" {
		ctxData["DeviceName"] = e.ContextDeviceName
	}
	if e.ContextCity != "" {
		ctxData["City"] = e.ContextCity
	}
	if e.ContextCountry != "" {
		ctxData["Country"] = e.ContextCountry
	}
	if e.ContextTimezone != "" {
		ctxData["Timezone"] = e.ContextTimezone
	}
	if len(ctxData) > 0 {
		result["EventContextData"] = ctxData
	}
	if e.FeedbackValue != "" {
		result["EventFeedback"] = map[string]interface{}{
			"FeedbackValue": e.FeedbackValue,
			"Provider":      e.FeedbackProvider,
			"FeedbackDate":  e.FeedbackDate.Unix(),
		}
	}
	return result
}
