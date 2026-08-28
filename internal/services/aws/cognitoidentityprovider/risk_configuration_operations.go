package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// SetRiskConfiguration sets the risk configuration for a user pool or client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SetRiskConfiguration.html
func (s *CognitoService) SetRiskConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := SetRiskConfigurationInput{
		Region:     reqCtx.GetRegion(),
		UserPoolID: req.GetParam("UserPoolId"),
		ClientID:   req.GetParam("ClientId"),
	}
	if m, ok := req.Parameters["CompromisedCredentialsRiskConfiguration"].(map[string]interface{}); ok {
		in.CompromisedCredentialsRiskConfiguration = m
	}
	if m, ok := req.Parameters["AccountTakeoverRiskConfiguration"].(map[string]interface{}); ok {
		in.AccountTakeoverRiskConfiguration = m
	}
	if m, ok := req.Parameters["RiskExceptionConfiguration"].(map[string]interface{}); ok {
		in.RiskExceptionConfiguration = m
	}

	cfg, err := s.setRiskConfigurationCore(in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"RiskConfiguration": formatRiskConfiguration(cfg),
	}, nil
}

// formatRiskConfiguration converts a stored RiskConfiguration to the API response format.
func formatRiskConfiguration(cfg *cognitostore.RiskConfiguration) map[string]interface{} {
	result := map[string]interface{}{
		"UserPoolId": cfg.UserPoolID,
	}

	if cfg.CompromisedCredentialsEventAction != "" || len(cfg.CompromisedCredentialsEventFilter) > 0 {
		cc := map[string]interface{}{}
		if len(cfg.CompromisedCredentialsEventFilter) > 0 {
			cc["EventFilter"] = cfg.CompromisedCredentialsEventFilter
		} else {
			cc["EventFilter"] = []string{}
		}
		cc["Actions"] = map[string]interface{}{
			"EventAction": defaultIfEmpty(cfg.CompromisedCredentialsEventAction, "NO_ACTION"),
		}
		result["CompromisedCredentialsRiskConfiguration"] = cc
	}

	if cfg.AccountTakeoverLowAction != "" || cfg.AccountTakeoverMediumAction != "" ||
		cfg.AccountTakeoverHighAction != "" || cfg.NotifyFrom != "" || cfg.NotifySourceArn != "" {
		at := map[string]interface{}{}
		if cfg.NotifyFrom != "" || cfg.NotifySourceArn != "" {
			notify := map[string]interface{}{}
			if cfg.NotifyFrom != "" {
				notify["From"] = cfg.NotifyFrom
			}
			if cfg.NotifyReplyTo != "" {
				notify["ReplyTo"] = cfg.NotifyReplyTo
			}
			if cfg.NotifySourceArn != "" {
				notify["SourceArn"] = cfg.NotifySourceArn
			}
			at["NotifyConfiguration"] = notify
		}
		at["Actions"] = map[string]interface{}{
			"HighAction":   map[string]interface{}{"EventAction": defaultIfEmpty(cfg.AccountTakeoverHighAction, "NO_ACTION"), "Notify": cfg.AccountTakeoverHighNotify},
			"MediumAction": map[string]interface{}{"EventAction": defaultIfEmpty(cfg.AccountTakeoverMediumAction, "NO_ACTION"), "Notify": cfg.AccountTakeoverMediumNotify},
			"LowAction":    map[string]interface{}{"EventAction": defaultIfEmpty(cfg.AccountTakeoverLowAction, "NO_ACTION"), "Notify": cfg.AccountTakeoverLowNotify},
		}
		result["AccountTakeoverRiskConfiguration"] = at
	}

	if len(cfg.BlockedIPRangeList) > 0 || len(cfg.SkippedIPRangeList) > 0 {
		result["RiskExceptionConfiguration"] = map[string]interface{}{
			"BlockedIPRangeList": cfg.BlockedIPRangeList,
			"SkippedIPRangeList": cfg.SkippedIPRangeList,
		}
	}

	if !cfg.LastModifiedDate.IsZero() {
		result["LastModifiedDate"] = cfg.LastModifiedDate.Unix()
	}

	if cfg.ClientID != "" {
		result["ClientId"] = cfg.ClientID
	}

	return result
}

func defaultIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

// isValidAccountTakeoverAction validates against the Smithy enum
// AccountTakeoverEventActionType: BLOCK, MFA_IF_CONFIGURED, MFA_REQUIRED, NO_ACTION.
func isValidAccountTakeoverAction(action string) bool {
	switch action {
	case "BLOCK", "MFA_IF_CONFIGURED", "MFA_REQUIRED", "NO_ACTION":
		return true
	}
	return false
}
