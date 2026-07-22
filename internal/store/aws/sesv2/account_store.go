package sesv2

import (
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
)

// AccountDetails represents the details associated with an SESv2 account.
type AccountDetails struct {
	MailType                        string   `json:"mailType,omitempty"`
	AdditionalContactEmailAddresses []string `json:"additionalContactEmailAddresses,omitempty"`
	UseCaseDescription              string   `json:"useCaseDescription,omitempty"`
	WebsiteURL                      string   `json:"websiteUrl,omitempty"`
	ContactLanguage                 string   `json:"contactLanguage,omitempty"`
	ProductionAccessEnabled         bool     `json:"productionAccessEnabled"`
	// ProductionAccessProvided is set true when the caller explicitly
	// supplied ProductionAccessEnabled in a PutAccountDetails call. Per
	// Smithy the wire type is EnabledWrapper, which distinguishes unset
	// from explicit false; without this sentinel we could not tell the
	// two apart because Go bool's zero value is false.
	ProductionAccessProvided bool `json:"productionAccessProvided,omitempty"`
}

// SendingAttributes represents the sending quota and rate attributes for an SESv2 account.
type SendingAttributes struct {
	MaxSendRate           float64 `json:"maxSendRate,omitempty"`
	Max24HourSend         float64 `json:"max24HourSend,omitempty"`
	SendQuota             float64 `json:"sendQuota,omitempty"`
	SentLast24Hours       float64 `json:"sentLast24Hours,omitempty"`
	DedicatedIpAutoWarmup bool    `json:"dedicatedIpAutoWarmupEnabled"`
}

// SuppressionAttributes represents the suppression reasons configured for an SESv2 account.
type SuppressionAttributes struct {
	SuppressedReasons []string `json:"suppressedReasons,omitempty"`
}

// VdmAttributes represents the VDM (Visibility in Delivery Metrics) attributes for an account.
type VdmAttributes struct {
	VdmEnabled                      bool     `json:"vdmEnabled"`
	DashboardAttributes             string   `json:"dashboardAttributes,omitempty"`
	GuardianAttributes              string   `json:"guardianAttributes,omitempty"`
	AdditionalContactEmailAddresses []string `json:"additionalContactEmailAddresses,omitempty"`
}

// Account represents the overall SESv2 account configuration.
//
// Note on SendingEnabled: this top-level field mirrors the wire-level
// response member on com.amazonaws.sesv2#GetAccountResponse. It is
// distinct from SendingAttributes (which carries the SendQuota fields).
// The PutAccountSendingAttributes operation flips this top-level flag.
// The two structures are intentionally separate; the apparent redundancy
// matches the AWS wire shape.
type Account struct {
	Details                 *AccountDetails        `json:"details,omitempty"`
	SendingAttributes       *SendingAttributes     `json:"sendingAttributes"`
	SuppressionAttributes   *SuppressionAttributes `json:"suppressionAttributes"`
	VdmAttributes           *VdmAttributes         `json:"vdmAttributes,omitempty"`
	EnforcementStatus       string                 `json:"enforcementStatus"`
	ProductionAccessEnabled bool                   `json:"productionAccessEnabled"`
	SendingEnabled          bool                   `json:"sendingEnabled"`
}

// DefaultAccount returns an Account populated with default values.
func DefaultAccount() *Account {
	return &Account{
		SendingAttributes: &SendingAttributes{
			MaxSendRate:     14.0,
			Max24HourSend:   50.0,
			SendQuota:       50.0,
			SentLast24Hours: 0.0,
		},
		SuppressionAttributes: &SuppressionAttributes{
			SuppressedReasons: []string{"BOUNCE", "COMPLAINT"},
		},
		EnforcementStatus:       "HEALTHY",
		ProductionAccessEnabled: true,
		SendingEnabled:          true,
	}
}

// GetAccount retrieves the SESv2 account configuration, initialising defaults if none exist.
func (s *SESv2Store) GetAccount() (*Account, error) {
	s.accountMu.Lock()
	defer s.accountMu.Unlock()
	var account Account
	if err := s.accountStore.Get("account", &account); err != nil {
		account = *DefaultAccount()
		if putErr := s.accountStore.Put("account", &account); putErr != nil {
			logs.Warn("failed to persist default account configuration", logs.Err(putErr))
		}
		return &account, nil
	}
	return &account, nil
}

// PutAccountDetails updates the account details for the SESv2 account.
// Per Smithy, ProductionAccessEnabled is an EnabledWrapper — when the
// caller omits the parameter, the existing stored value is preserved
// rather than being implicitly reset to false.
func (s *SESv2Store) PutAccountDetails(details *AccountDetails) error {
	s.accountMu.Lock()
	defer s.accountMu.Unlock()
	account, err := s.getAccountLocked()
	if err != nil {
		return err
	}
	if account.Details == nil {
		account.Details = &AccountDetails{}
	}
	existing := account.Details
	if details.ProductionAccessProvided {
		existing.ProductionAccessEnabled = details.ProductionAccessEnabled
		existing.ProductionAccessProvided = true
	}
	existing.MailType = details.MailType
	existing.UseCaseDescription = details.UseCaseDescription
	existing.WebsiteURL = details.WebsiteURL
	existing.ContactLanguage = details.ContactLanguage
	if len(details.AdditionalContactEmailAddresses) > 0 {
		existing.AdditionalContactEmailAddresses = details.AdditionalContactEmailAddresses
	}
	account.Details = existing
	return s.accountStore.Put("account", account)
}

// PutSendingAttributes updates the sending enabled flag for the SESv2 account.
func (s *SESv2Store) PutSendingAttributes(sendingEnabled bool) error {
	s.accountMu.Lock()
	defer s.accountMu.Unlock()
	account, err := s.getAccountLocked()
	if err != nil {
		return err
	}
	account.SendingEnabled = sendingEnabled
	return s.accountStore.Put("account", account)
}

// PutSuppressionAttributes updates the suppression reasons for the SESv2 account.
func (s *SESv2Store) PutSuppressionAttributes(reasons []string) error {
	s.accountMu.Lock()
	defer s.accountMu.Unlock()
	account, err := s.getAccountLocked()
	if err != nil {
		return err
	}
	account.SuppressionAttributes.SuppressedReasons = reasons
	return s.accountStore.Put("account", account)
}

// PutVdmAttributes updates the VDM attributes for the SESv2 account.
func (s *SESv2Store) PutVdmAttributes(vdm *VdmAttributes) error {
	s.accountMu.Lock()
	defer s.accountMu.Unlock()
	account, err := s.getAccountLocked()
	if err != nil {
		return err
	}
	account.VdmAttributes = vdm
	return s.accountStore.Put("account", account)
}

// PutDedicatedIpAutoWarmupEnabled updates the dedicated IP auto warmup flag.
func (s *SESv2Store) PutDedicatedIpAutoWarmupEnabled(enabled bool) error {
	s.accountMu.Lock()
	defer s.accountMu.Unlock()
	account, err := s.getAccountLocked()
	if err != nil {
		return err
	}
	if account.SendingAttributes == nil {
		account.SendingAttributes = &SendingAttributes{}
	}
	account.SendingAttributes.DedicatedIpAutoWarmup = enabled
	return s.accountStore.Put("account", account)
}

func (s *SESv2Store) getAccountLocked() (*Account, error) {
	var account Account
	if err := s.accountStore.Get("account", &account); err != nil {
		account = *DefaultAccount()
		if putErr := s.accountStore.Put("account", &account); putErr != nil {
			logs.Warn("failed to persist default account configuration", logs.Err(putErr))
		}
		return &account, nil
	}
	return &account, nil
}

// IdentityPolicy represents a policy attached to an SESv2 email identity.
type IdentityPolicy struct {
	PolicyName string `json:"policyName"`
	Policy     string `json:"policy"`
}

// PutEmailIdentityPolicy attaches a policy to an SESv2 email identity.
func (s *SESv2Store) PutEmailIdentityPolicy(identity, policyName, policy string) error {
	arn := s.BuildIdentityArn(identity)
	key := "policy#" + arn + "#" + policyName
	p := &IdentityPolicy{
		PolicyName: policyName,
		Policy:     policy,
	}
	return s.policyStore.Put(key, p)
}

// GetEmailIdentityPolicy retrieves a specific policy attached to an SESv2 email identity.
func (s *SESv2Store) GetEmailIdentityPolicy(identity, policyName string) (*IdentityPolicy, error) {
	arn := s.BuildIdentityArn(identity)
	key := "policy#" + arn + "#" + policyName
	var p IdentityPolicy
	if err := s.policyStore.Get(key, &p); err != nil {
		return nil, ErrPolicyNotFound
	}
	return &p, nil
}

// DeleteEmailIdentityPolicy removes a policy from an SESv2 email identity.
func (s *SESv2Store) DeleteEmailIdentityPolicy(identity, policyName string) error {
	arn := s.BuildIdentityArn(identity)
	key := "policy#" + arn + "#" + policyName
	if !s.policyStore.Exists(key) {
		return ErrPolicyNotFound
	}
	return s.policyStore.Delete(key)
}

// ListEmailIdentityPolicies returns all policies attached to an SESv2 email identity.
func (s *SESv2Store) ListEmailIdentityPolicies(identity string) ([]*IdentityPolicy, error) {
	arn := s.BuildIdentityArn(identity)
	prefix := "policy#" + arn + "#"
	policies, err := common.ListMatching[IdentityPolicy](s.policyStore, prefix, nil)
	if err != nil {
		return nil, err
	}
	result := make([]*IdentityPolicy, len(policies))
	copy(result, policies)
	return result, nil
}
