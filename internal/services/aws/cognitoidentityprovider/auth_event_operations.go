package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// AdminListUserAuthEvents lists authentication events for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminListUserAuthEvents.html
func (s *CognitoService) AdminListUserAuthEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminListUserAuthEventsCore(reqCtx, AdminListUserAuthEventsInput{
		UserPoolID: req.GetParam("UserPoolId"),
		Username:   getUsername(req),
		NextToken:  req.GetParam("NextToken"),
		Params:     req.Parameters,
	})
}

// AdminUpdateAuthEventFeedback updates feedback for an auth event (admin).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminUpdateAuthEventFeedback.html
func (s *CognitoService) AdminUpdateAuthEventFeedback(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminUpdateAuthEventFeedbackCore(reqCtx, AdminUpdateAuthEventFeedbackInput{
		UserPoolID:    req.GetParam("UserPoolId"),
		Username:      getUsername(req),
		EventID:       req.GetParam("EventId"),
		FeedbackValue: req.GetParam("FeedbackValue"),
	})
}

// UpdateAuthEventFeedback updates feedback for an auth event (user via access token).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateAuthEventFeedback.html
func (s *CognitoService) UpdateAuthEventFeedback(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateAuthEventFeedbackCore(reqCtx, UpdateAuthEventFeedbackInput{
		AccessToken:   getAccessToken(req),
		EventID:       req.GetParam("EventId"),
		FeedbackValue: req.GetParam("FeedbackValue"),
	})
}
