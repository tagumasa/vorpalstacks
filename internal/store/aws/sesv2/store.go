// Package sesv2 provides AWS SES v2 storage functionality for vorpalstacks.
package sesv2

import (
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// SESv2Store provides SES v2 email storage functionality.
type SESv2Store struct {
	*common.BaseStore
	*common.TagStore
	arnBuilder       *svcarn.ARNBuilder
	accountID        string
	region           string
	accountMu        sync.Mutex
	configSetStore   *common.BaseStore
	templateStore    *common.BaseStore
	emailRecordStore *common.BaseStore
	ipPoolStore      *common.BaseStore
	suppressionStore *common.BaseStore
	accountStore     *common.BaseStore
	contactListStore *common.BaseStore
	contactStore     *common.BaseStore
	policyStore      *common.BaseStore
}

// NewSESv2Store creates a new SESv2 store instance.
func NewSESv2Store(store storage.BasicStorage, accountID, region string) *SESv2Store {
	bucketName := "sesv2-identities-" + region
	return &SESv2Store{
		BaseStore:        common.NewBaseStore(store.Bucket(bucketName), "sesv2-identities"),
		TagStore:         common.NewTagStoreWithRegion(store, "sesv2", region),
		arnBuilder:       svcarn.NewARNBuilder(accountID, region),
		accountID:        accountID,
		region:           region,
		configSetStore:   common.NewBaseStore(store.Bucket("sesv2-configsets-"+region), "sesv2-configsets"),
		templateStore:    common.NewBaseStore(store.Bucket("sesv2-templates-"+region), "sesv2-templates"),
		emailRecordStore: common.NewBaseStore(store.Bucket("sesv2-emails-"+region), "sesv2-emails"),
		ipPoolStore:      common.NewBaseStore(store.Bucket("sesv2-ippools-"+region), "sesv2-ippools"),
		suppressionStore: common.NewBaseStore(store.Bucket("sesv2-suppression-"+region), "sesv2-suppression"),
		accountStore:     common.NewBaseStore(store.Bucket("sesv2-account-"+region), "sesv2-account"),
		contactListStore: common.NewBaseStore(store.Bucket("sesv2-contactlists-"+region), "sesv2-contactlists"),
		contactStore:     common.NewBaseStore(store.Bucket("sesv2-contacts-"+region), "sesv2-contacts"),
		policyStore:      common.NewBaseStore(store.Bucket("sesv2-policies-"+region), "sesv2-policies"),
	}
}

// GetAccountID returns the AWS account ID.
func (s *SESv2Store) GetAccountID() string { return s.accountID }

// GetRegion returns the AWS region.
func (s *SESv2Store) GetRegion() string { return s.region }

// BuildIdentityArn builds an ARN for an SES identity.
func (s *SESv2Store) BuildIdentityArn(identity string) string {
	return s.arnBuilder.Build("ses", "identity/"+identity)
}

// BuildConfigSetArn builds an ARN for an SES configuration set.
func (s *SESv2Store) BuildConfigSetArn(name string) string {
	return s.arnBuilder.Build("ses", "configuration-set/"+name)
}

// BuildTemplateArn builds an ARN for an SES template.
func (s *SESv2Store) BuildTemplateArn(name string) string {
	return s.arnBuilder.Build("ses", "template/"+name)
}

func (s *SESv2Store) buildEmailArn(messageId string) string {
	return s.arnBuilder.Build("ses", "email/"+messageId)
}

// DkimAttributes represents DKIM attributes for an email identity.
type DkimAttributes struct {
	Status                  string   `json:"status,omitempty"`
	SigningEnabled          bool     `json:"signingEnabled"`
	Tokens                  []string `json:"tokens,omitempty"`
	CurrentSigningKeyLength string   `json:"currentSigningKeyLength,omitempty"`
	NextSigningKeyLength    string   `json:"nextSigningKeyLength,omitempty"`
	SigningAttributesOrigin string   `json:"signingAttributesOrigin,omitempty"`
}

// EmailIdentity represents an email identity in SES.
type EmailIdentity struct {
	Identity             string              `json:"identity"`
	IdentityType         string              `json:"identityType"`
	VerifiedForSending   bool                `json:"verifiedForSending"`
	DkimAttributes       *DkimAttributes     `json:"dkimAttributes,omitempty"`
	ConfigurationSetName string              `json:"configurationSetName,omitempty"`
	MailFromAttributes   *MailFromAttributes `json:"mailFromAttributes,omitempty"`
	FeedbackForwarding   bool                `json:"feedbackForwarding"`
	CreatedTimestamp     time.Time           `json:"createdTimestamp"`
	LastUpdatedTimestamp time.Time           `json:"lastUpdatedTimestamp"`
}

// MailFromAttributes represents the MailFrom attributes for an email identity.
type MailFromAttributes struct {
	MailFromDomain       string `json:"mailFromDomain,omitempty"`
	BehaviorOnMxFailure  string `json:"behaviorOnMxFailure,omitempty"`
	MailFromDomainStatus string `json:"mailFromDomainStatus,omitempty"`
}

// NewEmailIdentity creates a new email identity.
func NewEmailIdentity(identity string) *EmailIdentity {
	now := time.Now().UTC()
	identityType := "EMAIL_ADDRESS"
	if isDomain(identity) {
		identityType = "DOMAIN"
	}

	return &EmailIdentity{
		Identity:             identity,
		IdentityType:         identityType,
		VerifiedForSending:   true,
		DkimAttributes:       GenerateDkimAttributes(identityType),
		FeedbackForwarding:   true,
		CreatedTimestamp:     now,
		LastUpdatedTimestamp: now,
	}
}

func isDomain(identity string) bool {
	return len(identity) > 0 && !strings.Contains(identity, "@")
}

// GenerateDkimAttributes generates DKIM attributes for an email identity.
func GenerateDkimAttributes(identityType string) *DkimAttributes {
	if identityType == "EMAIL_ADDRESS" {
		return nil
	}
	return &DkimAttributes{
		Status:         "SUCCESS",
		SigningEnabled: true,
		Tokens: []string{
			generateDkimToken(),
			generateDkimToken(),
			generateDkimToken(),
		},
		CurrentSigningKeyLength: "RSA_1024_BIT",
		SigningAttributesOrigin: "AWS_SES",
	}
}

func generateDkimToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[(time.Now().UnixNano()+int64(i))%int64(len(charset))]
	}
	return string(b)
}

// CreateEmailIdentity creates a new email identity.
func (s *SESv2Store) CreateEmailIdentity(identity *EmailIdentity) (*EmailIdentity, error) {
	arn := s.BuildIdentityArn(identity.Identity)
	if s.Exists(arn) {
		return nil, ErrIdentityAlreadyExists
	}

	identity.CreatedTimestamp = time.Now().UTC()
	identity.LastUpdatedTimestamp = identity.CreatedTimestamp

	if err := s.Put(arn, identity); err != nil {
		return nil, err
	}

	return identity, nil
}

// GetEmailIdentity retrieves an email identity.
func (s *SESv2Store) GetEmailIdentity(identity string) (*EmailIdentity, error) {
	arn := s.BuildIdentityArn(identity)
	var emailIdentity EmailIdentity
	if err := s.BaseStore.Get(arn, &emailIdentity); err != nil {
		return nil, ErrIdentityNotFound
	}
	return &emailIdentity, nil
}

// DeleteEmailIdentity deletes an email identity.
func (s *SESv2Store) DeleteEmailIdentity(identity string) error {
	arn := s.BuildIdentityArn(identity)
	if !s.Exists(arn) {
		return ErrIdentityNotFound
	}
	if err := s.TagStore.Delete(arn); err != nil {
		logs.Error("Failed to delete tags for identity", logs.String("identity", identity), logs.Err(err))
	}
	return s.BaseStore.Delete(arn)
}

// UpdateEmailIdentity updates an email identity.
func (s *SESv2Store) UpdateEmailIdentity(identity *EmailIdentity) error {
	arn := s.BuildIdentityArn(identity.Identity)
	if !s.Exists(arn) {
		return ErrIdentityNotFound
	}
	identity.LastUpdatedTimestamp = time.Now().UTC()
	return s.Put(arn, identity)
}

// ListEmailIdentities lists email identities.
func (s *SESv2Store) ListEmailIdentities(opts common.ListOptions) (*common.ListResult[EmailIdentity], error) {
	return common.List[EmailIdentity](s.BaseStore, opts, nil)
}

// IdentityExists checks if an email identity exists.
func (s *SESv2Store) IdentityExists(identity string) bool {
	arn := s.BuildIdentityArn(identity)
	return s.Exists(arn)
}

// DedicatedIpPool represents a dedicated IP pool.
type DedicatedIpPool struct {
	PoolName    string            `json:"poolName"`
	ScalingMode string            `json:"scalingMode,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// CreateDedicatedIpPool creates a dedicated IP pool.
func (s *SESv2Store) CreateDedicatedIpPool(pool *DedicatedIpPool) error {
	if s.ipPoolStore.Exists(pool.PoolName) {
		return ErrDedicatedIpPoolAlreadyExists
	}
	return s.ipPoolStore.Put(pool.PoolName, pool)
}

// GetDedicatedIpPool retrieves a dedicated IP pool.
func (s *SESv2Store) GetDedicatedIpPool(poolName string) (*DedicatedIpPool, error) {
	var pool DedicatedIpPool
	if err := s.ipPoolStore.Get(poolName, &pool); err != nil {
		return nil, err
	}
	return &pool, nil
}

// DeleteDedicatedIpPool deletes a dedicated IP pool.
func (s *SESv2Store) DeleteDedicatedIpPool(poolName string) error {
	if !s.ipPoolStore.Exists(poolName) {
		return ErrDedicatedIpPoolNotFound
	}
	return s.ipPoolStore.Delete(poolName)
}

// ListDedicatedIpPools lists dedicated IP pools.
func (s *SESv2Store) ListDedicatedIpPools(opts common.ListOptions) (*common.ListResult[DedicatedIpPool], error) {
	return common.List[DedicatedIpPool](s.ipPoolStore, opts, nil)
}

// SuppressedDestination represents a suppressed email destination.
type SuppressedDestination struct {
	EmailAddress      string            `json:"emailAddress"`
	Reason            string            `json:"reason"`
	SuppressionReason string            `json:"suppressionReason"`
	LastUpdateTime    string            `json:"lastUpdateTime"`
	Attributes        map[string]string `json:"attributes,omitempty"`
}

// PutSuppressedDestination adds or updates a suppressed destination.
func (s *SESv2Store) PutSuppressedDestination(dest *SuppressedDestination) error {
	return s.suppressionStore.Put(dest.EmailAddress, dest)
}

// GetSuppressedDestination retrieves a suppressed destination.
func (s *SESv2Store) GetSuppressedDestination(emailAddress string) (*SuppressedDestination, error) {
	var dest SuppressedDestination
	if err := s.suppressionStore.Get(emailAddress, &dest); err != nil {
		return nil, err
	}
	return &dest, nil
}

// DeleteSuppressedDestination removes a suppressed destination.
func (s *SESv2Store) DeleteSuppressedDestination(emailAddress string) error {
	if !s.suppressionStore.Exists(emailAddress) {
		return &SESv2StoreError{common.NewAWSError("NotFoundException", "Suppressed destination not found", 404)}
	}
	return s.suppressionStore.Delete(emailAddress)
}

// ListSuppressedDestinations lists suppressed destinations.
func (s *SESv2Store) ListSuppressedDestinations(opts common.ListOptions) (*common.ListResult[SuppressedDestination], error) {
	return common.List[SuppressedDestination](s.suppressionStore, opts, nil)
}
