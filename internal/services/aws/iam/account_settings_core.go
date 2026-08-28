package iam

import (
	"fmt"
	"net/http"
	"net/url"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/config"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// enableOutboundWebIdentityFederationCore enables outbound web identity
// federation for the account and returns the derived issuer identifier.
func (s *IAMService) enableOutboundWebIdentityFederationCore(store *iamstore.IAMStore) (string, error) {
	settings, err := store.AccountSettings().Get()
	if err != nil {
		return "", err
	}

	// Enabling an already enabled feature fails with FeatureEnabled.
	if settings.OutboundWebIdentityFederationEnabled {
		return "", ErrFeatureEnabled
	}

	issuerIdentifier := deriveOIDCIssuer(s.accountID)

	settings.OutboundWebIdentityFederationEnabled = true
	settings.IssuerIdentifier = issuerIdentifier

	if err := store.AccountSettings().Put(settings); err != nil {
		return "", err
	}
	return issuerIdentifier, nil
}

// deriveOIDCIssuer builds the OIDC issuer URL from the configured
// endpoints.base_url, falling back to the amazonaws.com format when
// the base URL cannot be parsed.
//
// The host segment preserves non-default ports so that edge deployments
// exposed on a non-443 port produce a usable issuer URL (e.g.
// https://oidc.<account>.edge.example:8443). Default ports (80 for http,
// 443 for https) are stripped because including them in the issuer would
// diverge from the canonical URL form expected by JWT consumers.
func deriveOIDCIssuer(accountID string) string {
	baseURL := config.BaseURL()
	if baseURL == "" {
		return fmt.Sprintf("https://oidc.%s.amazonaws.com", accountID)
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Hostname() == "" {
		return fmt.Sprintf("https://oidc.%s.amazonaws.com", accountID)
	}
	host := parsed.Hostname()
	if port := parsed.Port(); port != "" && port != "80" && port != "443" {
		host = host + ":" + port
	}
	return fmt.Sprintf("https://oidc.%s.%s", accountID, host)
}

// disableOutboundWebIdentityFederationCore disables outbound web identity
// federation for the account.
func (s *IAMService) disableOutboundWebIdentityFederationCore(store *iamstore.IAMStore) error {
	settings, err := store.AccountSettings().Get()
	if err != nil {
		return err
	}

	// Disabling a feature that is not enabled fails with FeatureDisabled.
	if !settings.OutboundWebIdentityFederationEnabled {
		return ErrFeatureDisabled
	}

	settings.OutboundWebIdentityFederationEnabled = false
	settings.IssuerIdentifier = ""

	return store.AccountSettings().Put(settings)
}

// getOutboundWebIdentityFederationCore retrieves the outbound web identity
// federation configuration for the account.  Reading the configuration of
// a disabled feature fails with FeatureDisabled.
func (s *IAMService) getOutboundWebIdentityFederationCore(store *iamstore.IAMStore) (*iamstore.AccountSettings, error) {
	settings, err := store.AccountSettings().Get()
	if err != nil {
		return nil, err
	}
	if !settings.OutboundWebIdentityFederationEnabled {
		return nil, ErrFeatureDisabled
	}
	return settings, nil
}

// setSecurityTokenServicePreferencesCore validates input and sets the
// global endpoint token version preference for the account.
func (s *IAMService) setSecurityTokenServicePreferencesCore(store *iamstore.IAMStore, tokenVersion string) error {
	if tokenVersion == "" {
		return NewValidationError("GlobalEndpointTokenVersion")
	}

	if tokenVersion != "v1Token" && tokenVersion != "v2Token" {
		return errors.NewAWSError("InvalidInput", tokenVersion+" is not a valid value for GlobalEndpointTokenVersion. Must be v1Token or v2Token.", http.StatusBadRequest)
	}

	settings, err := store.AccountSettings().Get()
	if err != nil {
		return err
	}

	settings.GlobalEndpointTokenVersion = tokenVersion

	return store.AccountSettings().Put(settings)
}
