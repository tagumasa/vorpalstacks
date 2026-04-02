package sesv2

import "vorpalstacks/internal/store/aws/common"

// SESv2StoreInterface defines operations for managing SESv2 resources.
type SESv2StoreInterface interface {
	CreateConfigurationSet(configSet *ConfigurationSet) (*ConfigurationSet, error)
	GetConfigurationSet(name string) (*ConfigurationSet, error)
	DeleteConfigurationSet(name string) error
	UpdateConfigurationSet(configSet *ConfigurationSet) error
	ListConfigurationSets(opts common.ListOptions) (*common.ListResult[ConfigurationSet], error)
	ConfigurationSetExists(name string) bool

	SaveEmailRecord(email *EmailRecord) error
	GetEmailRecord(messageId string) (*EmailRecord, error)
	ListEmailRecords(opts common.ListOptions) (*common.ListResult[EmailRecord], error)

	GetAccountID() string
	GetRegion() string
	BuildIdentityArn(identity string) string
	BuildConfigSetArn(name string) string
	BuildTemplateArn(name string) string

	CreateEmailIdentity(identity *EmailIdentity) (*EmailIdentity, error)
	GetEmailIdentity(identity string) (*EmailIdentity, error)
	DeleteEmailIdentity(identity string) error
	UpdateEmailIdentity(identity *EmailIdentity) error
	ListEmailIdentities(opts common.ListOptions) (*common.ListResult[EmailIdentity], error)
	IdentityExists(identity string) bool

	CreateDedicatedIpPool(pool *DedicatedIpPool) error
	GetDedicatedIpPool(poolName string) (*DedicatedIpPool, error)
	DeleteDedicatedIpPool(poolName string) error
	ListDedicatedIpPools(opts common.ListOptions) (*common.ListResult[DedicatedIpPool], error)

	PutSuppressedDestination(dest *SuppressedDestination) error
	GetSuppressedDestination(emailAddress string) (*SuppressedDestination, error)
	DeleteSuppressedDestination(emailAddress string) error
	ListSuppressedDestinations(opts common.ListOptions) (*common.ListResult[SuppressedDestination], error)

	CreateEmailTemplate(template *EmailTemplate) (*EmailTemplate, error)
	GetEmailTemplate(name string) (*EmailTemplate, error)
	DeleteEmailTemplate(name string) error
	UpdateEmailTemplate(template *EmailTemplate) error
	ListEmailTemplates(opts common.ListOptions) (*common.ListResult[EmailTemplate], error)

	GetAccount() (*Account, error)
	PutAccountDetails(details *AccountDetails) error
	PutSendingAttributes(sendingEnabled bool) error
	PutSuppressionAttributes(reasons []string) error
	PutVdmAttributes(vdm *VdmAttributes) error
	PutDedicatedIpAutoWarmupEnabled(enabled bool) error

	PutEmailIdentityPolicy(identity, policyName, policy string) error
	DeleteEmailIdentityPolicy(identity, policyName string) error
	ListEmailIdentityPolicies(identity string) ([]*IdentityPolicy, error)

	CreateContactList(cl *ContactList) (*ContactList, error)
	GetContactList(name string) (*ContactList, error)
	DeleteContactList(name string) error
	UpdateContactList(cl *ContactList) error
	ListContactLists(opts common.ListOptions) (*common.ListResult[ContactList], error)

	CreateContact(c *Contact) (*Contact, error)
	GetContact(contactListName, emailAddress string) (*Contact, error)
	DeleteContact(contactListName, emailAddress string) error
	UpdateContact(c *Contact) error
	ListContacts(contactListName string, opts common.ListOptions) (*common.ListResult[Contact], error)

	TagResourceFromSlice(resourceKey string, tags []common.Tag) error
	UntagResource(resourceKey string, tagKeys []string) error
	ListTagsAsSlice(resourceKey string) ([]common.Tag, error)

	Raw() *SESv2Store
}

// Raw returns the underlying SESv2 store.
func (s *SESv2Store) Raw() *SESv2Store {
	return s
}
