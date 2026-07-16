package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// SetRiskConfiguration sets the risk configuration for a user pool or client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SetRiskConfiguration.html
func (s *CognitoService) SetRiskConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	clientID := req.GetParam("ClientId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	cfg := &cognitostore.RiskConfiguration{
		UserPoolID: userPoolID,
		ClientID:   clientID,
	}

	if ccRaw, ok := req.Parameters["CompromisedCredentialsRiskConfiguration"]; ok {
		if m, ok := ccRaw.(map[string]interface{}); ok {
			if actions, ok := m["Actions"].(map[string]interface{}); ok {
				action := getStringParam(actions, "EventAction")
				if action != "BLOCK" && action != "NO_ACTION" {
					return nil, ErrInvalidParameter
				}
				cfg.CompromisedCredentialsEventAction = action
			}
			if ef, ok := m["EventFilter"].([]interface{}); ok {
				for _, v := range ef {
					if s, ok := v.(string); ok {
						cfg.CompromisedCredentialsEventFilter = append(cfg.CompromisedCredentialsEventFilter, s)
					}
				}
			}
		}
	}

	if atRaw, ok := req.Parameters["AccountTakeoverRiskConfiguration"]; ok {
		if m, ok := atRaw.(map[string]interface{}); ok {
			if actions, ok := m["Actions"].(map[string]interface{}); ok {
				if low, ok := actions["LowAction"].(map[string]interface{}); ok {
					action := getStringParam(low, "EventAction")
					if !isValidAccountTakeoverAction(action) {
						return nil, ErrInvalidParameter
					}
					cfg.AccountTakeoverLowAction = action
					if notify, ok := low["Notify"].(bool); ok {
						cfg.AccountTakeoverLowNotify = notify
					}
				}
				if med, ok := actions["MediumAction"].(map[string]interface{}); ok {
					action := getStringParam(med, "EventAction")
					if !isValidAccountTakeoverAction(action) {
						return nil, ErrInvalidParameter
					}
					cfg.AccountTakeoverMediumAction = action
					if notify, ok := med["Notify"].(bool); ok {
						cfg.AccountTakeoverMediumNotify = notify
					}
				}
				if high, ok := actions["HighAction"].(map[string]interface{}); ok {
					action := getStringParam(high, "EventAction")
					if !isValidAccountTakeoverAction(action) {
						return nil, ErrInvalidParameter
					}
					cfg.AccountTakeoverHighAction = action
					if notify, ok := high["Notify"].(bool); ok {
						cfg.AccountTakeoverHighNotify = notify
					}
				}
			}
			if notify, ok := m["NotifyConfiguration"].(map[string]interface{}); ok {
				cfg.NotifyFrom = getStringParam(notify, "From")
				cfg.NotifyReplyTo = getStringParam(notify, "ReplyTo")
				cfg.NotifySourceArn = getStringParam(notify, "SourceArn")
				if blockEmail, ok := notify["BlockEmail"].(map[string]interface{}); ok {
					cfg.NotifyBlockEmailSubject = getStringParam(blockEmail, "Subject")
					cfg.NotifyBlockEmailHtml = getStringParam(blockEmail, "HtmlBody")
				}
				if noAction, ok := notify["NoActionEmail"].(map[string]interface{}); ok {
					cfg.NotifyNoActionEmailSubject = getStringParam(noAction, "Subject")
					cfg.NotifyNoActionEmailHtml = getStringParam(noAction, "HtmlBody")
				}
				if mfa, ok := notify["MfaEmail"].(map[string]interface{}); ok {
					cfg.NotifyMfaEmailSubject = getStringParam(mfa, "Subject")
					cfg.NotifyMfaEmailHtml = getStringParam(mfa, "HtmlBody")
				}
			}
		}
	}

	if reRaw, ok := req.Parameters["RiskExceptionConfiguration"]; ok {
		if m, ok := reRaw.(map[string]interface{}); ok {
			if blocked, ok := m["BlockedIPRangeList"].([]interface{}); ok {
				for _, v := range blocked {
					if s, ok := v.(string); ok {
						cfg.BlockedIPRangeList = append(cfg.BlockedIPRangeList, s)
					}
				}
			}
			if skipped, ok := m["SkippedIPRangeList"].([]interface{}); ok {
				for _, v := range skipped {
					if s, ok := v.(string); ok {
						cfg.SkippedIPRangeList = append(cfg.SkippedIPRangeList, s)
					}
				}
			}
		}
	}

	if err := store.SaveRiskConfiguration(cfg); err != nil {
		return nil, ErrInternalError
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
