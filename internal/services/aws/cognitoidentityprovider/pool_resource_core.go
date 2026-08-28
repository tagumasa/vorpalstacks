package cognitoidentityprovider

import (
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/config"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	"vorpalstacks/internal/store/aws/common"
)

// Core functions for the user-pool domain, resource-server and risk
// configuration families. The HTTP handlers are thin transport adapters:
// they extract the wire members into the Input structs below and serialise
// the Core results; every validation and store call lives here.

// CreateUserPoolDomainInput carries the raw wire members of
// CreateUserPoolDomain. The nested maps are the unmodified request payloads;
// the Core performs the member extraction and validation.
type CreateUserPoolDomainInput struct {
	Region              string
	Domain              string
	UserPoolID          string
	ManagedLoginVersion interface{}
	CustomDomainConfig  map[string]interface{}
	Routing             map[string]interface{}
}

// UpdateUserPoolDomainInput carries the raw wire members of
// UpdateUserPoolDomain. Absent optional members preserve the stored values.
type UpdateUserPoolDomainInput struct {
	Region              string
	Domain              string
	UserPoolID          string
	ManagedLoginVersion interface{}
	CustomDomainConfig  map[string]interface{}
	Routing             map[string]interface{}
}

// UpdateResourceServerInput carries the wire members of UpdateResourceServer.
// ScopesPresent distinguishes an omitted Scopes member from an empty list.
type UpdateResourceServerInput struct {
	Region        string
	UserPoolID    string
	Identifier    string
	Name          string
	Scopes        []interface{}
	ScopesPresent bool
}

// SetRiskConfigurationInput carries the raw nested wire members of
// SetRiskConfiguration; absent members arrive as nil maps.
type SetRiskConfigurationInput struct {
	Region                                  string
	UserPoolID                              string
	ClientID                                string
	CompromisedCredentialsRiskConfiguration map[string]interface{}
	AccountTakeoverRiskConfiguration        map[string]interface{}
	RiskExceptionConfiguration              map[string]interface{}
}

// createUserPoolDomainCore validates the request, resolves the CloudFront
// domain name from the platform endpoint suffix and persists the domain
// entry. The returned entry feeds both the response serialisation and the
// stored state.
func (s *CognitoService) createUserPoolDomainCore(in CreateUserPoolDomainInput) (*cognitostore.UserPoolDomain, error) {
	if in.Domain == "" || in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(in.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	cognitoSuffix := config.GetString("endpoints.cognito_suffix")
	cfDomain := fmt.Sprintf("%s.auth.%s", in.Domain, strings.Replace(cognitoSuffix, "{region}", in.Region, 1))
	domainEntry := &cognitostore.UserPoolDomain{
		Domain:           in.Domain,
		UserPoolID:       in.UserPoolID,
		CloudFrontDomain: cfDomain,
		CreatedDate:      time.Now().UTC(),
		Status:           "ACTIVE",
	}
	if in.ManagedLoginVersion != nil {
		switch n := in.ManagedLoginVersion.(type) {
		case int:
			mlv := n
			domainEntry.ManagedLoginVersion = &mlv
		case float64:
			mlv := int(n)
			domainEntry.ManagedLoginVersion = &mlv
		}
	}
	if cdc := in.CustomDomainConfig; cdc != nil {
		cfg := &cognitostore.CustomDomainConfig{}
		if v, ok := cdc["CertificateArn"].(string); ok {
			cfg.CertificateArn = v
		}
		if v, ok := cdc["SecurityPolicy"].(string); ok {
			if !validateSecurityPolicy(v) {
				return nil, ErrInvalidParameter
			}
			cfg.SecurityPolicy = v
		}
		domainEntry.CustomDomainConfig = cfg
	}
	if rt := in.Routing; rt != nil {
		routing := &cognitostore.Routing{}
		if fo, ok := rt["Failover"].(map[string]interface{}); ok {
			failover := &cognitostore.FailoverType{}
			if v, ok := fo["SecondaryRegion"].(string); ok {
				failover.SecondaryRegion = v
			}
			if v, ok := fo["PrimaryRoute53HealthCheckId"].(string); ok {
				failover.PrimaryRoute53HealthCheckId = v
			}
			routing.Failover = failover
		}
		domainEntry.Routing = routing
	}
	if err := store.SetUserPoolDomain(in.Domain, domainEntry); err != nil {
		return nil, err
	}

	return domainEntry, nil
}

// describeUserPoolDomainCore loads a user-pool domain entry by domain name.
func (s *CognitoService) describeUserPoolDomainCore(region, domain string) (*cognitostore.UserPoolDomain, error) {
	if domain == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	domainEntry, err := store.GetUserPoolDomain(domain)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return domainEntry, nil
}

// deleteUserPoolDomainCore removes a user-pool domain entry.
func (s *CognitoService) deleteUserPoolDomainCore(region, domain string) error {
	if domain == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	return store.DeleteUserPoolDomain(domain)
}

// updateUserPoolDomainCore rebuilds the stored domain entry for a domain,
// preserving the original creation date, and returns the CloudFront domain.
func (s *CognitoService) updateUserPoolDomainCore(in UpdateUserPoolDomainInput) (*cognitostore.UserPoolDomain, error) {
	if in.Domain == "" || in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(in.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	existing, err := store.GetUserPoolDomain(in.Domain)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	cognitoSuffix := config.GetString("endpoints.cognito_suffix")
	cfDomain := fmt.Sprintf("%s.auth.%s", in.Domain, strings.Replace(cognitoSuffix, "{region}", in.Region, 1))
	domainEntry := &cognitostore.UserPoolDomain{
		Domain:           in.Domain,
		UserPoolID:       in.UserPoolID,
		CloudFrontDomain: cfDomain,
		CreatedDate:      existing.CreatedDate,
		Status:           existing.Status,
	}

	// The three optional members follow the update contract: a provided
	// member replaces the stored value, an absent member preserves it.
	if in.ManagedLoginVersion != nil {
		switch n := in.ManagedLoginVersion.(type) {
		case int:
			mlv := n
			domainEntry.ManagedLoginVersion = &mlv
		case float64:
			mlv := int(n)
			domainEntry.ManagedLoginVersion = &mlv
		default:
			domainEntry.ManagedLoginVersion = existing.ManagedLoginVersion
		}
	} else {
		domainEntry.ManagedLoginVersion = existing.ManagedLoginVersion
	}

	if in.CustomDomainConfig != nil {
		cfg := &cognitostore.CustomDomainConfig{}
		if v, ok := in.CustomDomainConfig["CertificateArn"].(string); ok {
			cfg.CertificateArn = v
		}
		if v, ok := in.CustomDomainConfig["SecurityPolicy"].(string); ok {
			if !validateSecurityPolicy(v) {
				return nil, ErrInvalidParameter
			}
			cfg.SecurityPolicy = v
		}
		domainEntry.CustomDomainConfig = cfg
	} else {
		domainEntry.CustomDomainConfig = existing.CustomDomainConfig
	}

	if in.Routing != nil {
		routing := &cognitostore.Routing{}
		if fo, ok := in.Routing["Failover"].(map[string]interface{}); ok {
			failover := &cognitostore.FailoverType{}
			if v, ok := fo["SecondaryRegion"].(string); ok {
				failover.SecondaryRegion = v
			}
			if v, ok := fo["PrimaryRoute53HealthCheckId"].(string); ok {
				failover.PrimaryRoute53HealthCheckId = v
			}
			routing.Failover = failover
		}
		domainEntry.Routing = routing
	} else {
		domainEntry.Routing = existing.Routing
	}

	if err := store.SetUserPoolDomain(in.Domain, domainEntry); err != nil {
		return nil, err
	}

	return domainEntry, nil
}

// createResourceServerCore validates the identifier, name and scope list,
// then persists the resource server.
func (s *CognitoService) createResourceServerCore(region, userPoolID, identifier, name string, scopes []interface{}) (*cognitostore.ResourceServer, error) {
	if userPoolID == "" || identifier == "" || name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	rs := &cognitostore.ResourceServer{
		UserPoolID: userPoolID,
		Identifier: identifier,
		Name:       name,
		Scopes:     []cognitostore.ResourceServerScope{},
	}

	if scopes != nil {
		parsed, err := parseResourceServerScopes(scopes)
		if err != nil {
			return nil, err
		}
		rs.Scopes = parsed
	}

	if err := store.CreateResourceServer(rs); err != nil {
		return nil, err
	}

	return rs, nil
}

// describeResourceServerCore loads a resource server by pool and identifier.
func (s *CognitoService) describeResourceServerCore(region, userPoolID, identifier string) (*cognitostore.ResourceServer, error) {
	if userPoolID == "" || identifier == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	rs, err := store.GetResourceServer(userPoolID, identifier)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return rs, nil
}

// updateResourceServerCore applies the mutable members (name, full scope
// list replacement) onto the stored resource server.
func (s *CognitoService) updateResourceServerCore(in UpdateResourceServerInput) (*cognitostore.ResourceServer, error) {
	if in.UserPoolID == "" || in.Identifier == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, err
	}

	rs, err := store.GetResourceServer(in.UserPoolID, in.Identifier)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if in.Name != "" {
		rs.Name = in.Name
	}

	if in.ScopesPresent {
		newScopes, err := parseResourceServerScopes(in.Scopes)
		if err != nil {
			return nil, err
		}
		rs.Scopes = newScopes
	}

	if err := store.UpdateResourceServer(rs); err != nil {
		return nil, ErrInternalError
	}

	return rs, nil
}

// deleteResourceServerCore removes a resource server.
func (s *CognitoService) deleteResourceServerCore(region, userPoolID, identifier string) error {
	if userPoolID == "" || identifier == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteResourceServer(userPoolID, identifier); err != nil {
		return ErrResourceNotFound
	}
	return nil
}

// listResourceServersCore pages through the pool's resource servers with the
// documented default and maximum page size of 50.
func (s *CognitoService) listResourceServersCore(region, userPoolID string, maxResults int, nextToken string) (*common.ListResult[cognitostore.ResourceServer], error) {
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	if maxResults <= 0 || maxResults > 50 {
		maxResults = 50
	}

	opts := common.ListOptions{
		MaxItems: maxResults,
		Marker:   nextToken,
	}

	return store.ListResourceServersPaginated(userPoolID, opts)
}

// getCSVHeaderCore resolves the pool whose schema drives the CSV header.
func (s *CognitoService) getCSVHeaderCore(region, userPoolID string) (*cognitostore.UserPool, error) {
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return pool, nil
}

// describeRiskConfigurationCore resolves the effective risk configuration:
// the client-specific entry when present, otherwise the pool-level entry,
// otherwise a fresh default bound to the pool.
func (s *CognitoService) describeRiskConfigurationCore(region, userPoolID, clientID string) (*cognitostore.RiskConfiguration, error) {
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	cfg, err := store.GetRiskConfiguration(userPoolID, clientID)
	if err == nil && cfg != nil {
		return cfg, nil
	}

	if clientID != "" {
		cfg, err = store.GetRiskConfiguration(userPoolID, "")
		if err == nil && cfg != nil {
			return cfg, nil
		}
	}

	return &cognitostore.RiskConfiguration{UserPoolID: userPoolID}, nil
}

// setRiskConfigurationCore parses and validates the risk-configuration
// members and persists the resulting configuration.
func (s *CognitoService) setRiskConfigurationCore(in SetRiskConfigurationInput) (*cognitostore.RiskConfiguration, error) {
	if in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(in.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	cfg := &cognitostore.RiskConfiguration{
		UserPoolID: in.UserPoolID,
		ClientID:   in.ClientID,
	}

	if m := in.CompromisedCredentialsRiskConfiguration; m != nil {
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
					if !validateEventFilter(s) {
						return nil, ErrInvalidParameter
					}
					cfg.CompromisedCredentialsEventFilter = append(cfg.CompromisedCredentialsEventFilter, s)
				}
			}
		}
	}

	if m := in.AccountTakeoverRiskConfiguration; m != nil {
		if err := applyAccountTakeoverRiskConfiguration(cfg, m); err != nil {
			return nil, err
		}
	}

	if m := in.RiskExceptionConfiguration; m != nil {
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

	if err := store.SaveRiskConfiguration(cfg); err != nil {
		return nil, ErrInternalError
	}

	return cfg, nil
}

// applyAccountTakeoverRiskConfiguration parses the account-takeover actions
// and notify configuration onto cfg, rejecting non-member enum values.
func applyAccountTakeoverRiskConfiguration(cfg *cognitostore.RiskConfiguration, m map[string]interface{}) error {
	if actions, ok := m["Actions"].(map[string]interface{}); ok {
		if low, ok := actions["LowAction"].(map[string]interface{}); ok {
			action := getStringParam(low, "EventAction")
			if !isValidAccountTakeoverAction(action) {
				return ErrInvalidParameter
			}
			cfg.AccountTakeoverLowAction = action
			if notify, ok := low["Notify"].(bool); ok {
				cfg.AccountTakeoverLowNotify = notify
			}
		}
		if med, ok := actions["MediumAction"].(map[string]interface{}); ok {
			action := getStringParam(med, "EventAction")
			if !isValidAccountTakeoverAction(action) {
				return ErrInvalidParameter
			}
			cfg.AccountTakeoverMediumAction = action
			if notify, ok := med["Notify"].(bool); ok {
				cfg.AccountTakeoverMediumNotify = notify
			}
		}
		if high, ok := actions["HighAction"].(map[string]interface{}); ok {
			action := getStringParam(high, "EventAction")
			if !isValidAccountTakeoverAction(action) {
				return ErrInvalidParameter
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
	return nil
}

// parseResourceServerScopes converts the raw wire scope list into store
// scope records, rejecting entries without a scope name.
func parseResourceServerScopes(scopes []interface{}) ([]cognitostore.ResourceServerScope, error) {
	parsed := make([]cognitostore.ResourceServerScope, 0, len(scopes))
	for _, sc := range scopes {
		m, ok := sc.(map[string]interface{})
		if !ok {
			continue
		}
		scopeName, _ := m["ScopeName"].(string)
		scopeDesc, _ := m["ScopeDescription"].(string)
		if scopeName == "" {
			return nil, ErrInvalidParameter
		}
		parsed = append(parsed, cognitostore.ResourceServerScope{
			ScopeName:        scopeName,
			ScopeDescription: scopeDesc,
		})
	}
	return parsed, nil
}
