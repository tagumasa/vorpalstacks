package iot

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// IotStoreInterface is the composition of all per-entity Ops interfaces.
// Callers should depend on the narrowest interface needed (ISP/SOLID).
type IotStoreInterface interface {
	ThingOps
	ThingTypeOps
	ThingGroupOps
	BillingGroupOps
	CertificateOps
	PolicyOps
	TopicRuleOps
	RoleAliasOps
	JobOps
	ShadowOps
	AuthorizerOps
	ProvisioningTemplateOps
	GenericKVOps
	SecurityProfileOps
	ViolationOps
	DomainConfigOps
	IndexingConfigOps
	ProvisioningTemplateVersionOps
	TagOps
	PolicyAttachmentOps
	ThingPrincipalOps
	MetaOps
}

type ThingOps interface {
	CreateThing(thing *Thing) (*Thing, error)
	GetThing(thingName string) (*Thing, error)
	UpdateThing(thingName string, opts ThingUpdateOpts) (*Thing, error)
	DeleteThing(thingName string) error
	ListThings(opts common.ListOptions, attributeName, attributeValue string) (*common.ListResult[Thing], error)
	ListThingsForThingType(thingTypeName string, opts common.ListOptions) (*common.ListResult[Thing], error)
}

type ThingTypeOps interface {
	CreateThingType(tt *ThingType) (*ThingType, error)
	GetThingType(name string) (*ThingType, error)
	UpdateThingType(name string, opts ThingTypeUpdateOpts) (*ThingType, error)
	SetThingTypeDeprecation(name string, deprecated bool) (*ThingType, error)
	DeleteThingType(name string) error
	ListThingTypes(opts common.ListOptions) (*common.ListResult[ThingType], error)
}

type ThingGroupOps interface {
	CreateThingGroup(group *ThingGroup) (*ThingGroup, error)
	GetThingGroup(name string) (*ThingGroup, error)
	UpdateThingGroup(groupName string, opts ThingGroupUpdateOpts) (*ThingGroup, error)
	DeleteThingGroup(name string) error
	ListThingGroups(opts common.ListOptions, parentGroupName string) (*common.ListResult[ThingGroup], error)
}

type BillingGroupOps interface {
	CreateBillingGroup(bg *BillingGroup) (*BillingGroup, error)
	GetBillingGroup(name string) (*BillingGroup, error)
	UpdateBillingGroup(name string, opts BillingGroupUpdateOpts) (*BillingGroup, error)
	DeleteBillingGroup(name string) error
	ListBillingGroups(opts common.ListOptions) (*common.ListResult[BillingGroup], error)
	AddThingToBillingGroup(thingName, billingGroup string) error
	RemoveThingFromBillingGroup(thingName, billingGroup string) error
	ListThingsInBillingGroup(billingGroup string) ([]string, error)
	ListBillingGroupsForThing(thingName string) ([]string, error)
}

type CertificateOps interface {
	CreateCertificate(cert *Certificate) (*Certificate, error)
	GetCertificate(certificateID string) (*Certificate, error)
	UpdateCertificate(certID string, opts CertificateUpdateOpts) (*Certificate, error)
	DeleteCertificate(certificateID string) error
	ListCertificates(opts common.ListOptions) (*common.ListResult[Certificate], error)
}

type PolicyOps interface {
	CreatePolicy(policy *Policy) (*Policy, error)
	GetPolicy(policyName string) (*Policy, error)
	UpdatePolicy(policy *Policy) error
	DeletePolicy(policyName string) error
	ListPolicies(opts common.ListOptions) (*common.ListResult[Policy], error)
}

type TopicRuleOps interface {
	CreateRule(rule *TopicRule) (*TopicRule, error)
	GetRule(ruleName string) (*TopicRule, error)
	UpdateRule(ruleName string, opts RuleUpdateOpts) (*TopicRule, error)
	DeleteRule(ruleName string) error
	ListRules(opts common.ListOptions) (*common.ListResult[TopicRule], error)
}

type RoleAliasOps interface {
	CreateRoleAlias(ra *RoleAlias) (*RoleAlias, error)
	GetRoleAlias(alias string) (*RoleAlias, error)
	UpdateRoleAlias(alias string, opts RoleAliasUpdateOpts) (*RoleAlias, error)
	DeleteRoleAlias(alias string) error
	ListRoleAliases(opts common.ListOptions) (*common.ListResult[RoleAlias], error)
}

type JobOps interface {
	CreateJob(job *Job) (*Job, error)
	AssociateJobTargets(jobID string, newTargets []string, comment string) (*Job, error)
	GetJob(jobID string) (*Job, error)
	UpdateJob(jobID string, opts JobUpdateOpts) (*Job, error)
	DeleteJob(jobID string) error
	ListJobs(opts common.ListOptions, statusFilter string) (*common.ListResult[Job], error)
}

type ShadowOps interface {
	GetShadow(thingName, shadowName string) (*ShadowDocument, error)
	UpdateShadow(thingName, shadowName string, incoming ShadowState, clientVersion int64) (*ShadowUpdateResult, error)
	PutShadow(thingName, shadowName string, doc *ShadowDocument) error
	PutShadowWithVersion(thingName, shadowName string, doc *ShadowDocument, clientVersion int64) error
	DeleteShadow(thingName, shadowName string) error
	ListShadowNames(thingName string, opts common.ListOptions) ([]string, string, error)
}

type AuthorizerOps interface {
	CreateAuthorizer(a *Authorizer) (*Authorizer, error)
	GetAuthorizer(name string) (*Authorizer, error)
	UpdateAuthorizer(name string, opts AuthorizerUpdateOpts) (*Authorizer, error)
	DeleteAuthorizer(name string) error
	ListAuthorizers(opts common.ListOptions) (*common.ListResult[Authorizer], error)
}

type ProvisioningTemplateOps interface {
	CreateProvisioningTemplate(t *ProvisioningTemplate) (*ProvisioningTemplate, error)
	GetProvisioningTemplate(name string) (*ProvisioningTemplate, error)
	UpdateProvisioningTemplate(name string, opts ProvisioningTemplateUpdateOpts) (*ProvisioningTemplate, error)
	DeleteProvisioningTemplate(name string) error
	ListProvisioningTemplates(opts common.ListOptions) (*common.ListResult[ProvisioningTemplate], error)
}

type GenericKVOps interface {
	PutGeneric(key string, value interface{}) error
	GetGeneric(key string, dest interface{}) error
	GetGenericExists(key string, dest interface{}) (bool, error)
	DeleteGeneric(key string) error
	ListGeneric(prefix string) ([]map[string]interface{}, error)
}

type SecurityProfileOps interface {
	CreateSecurityProfile(sp *SecurityProfile) (*SecurityProfile, error)
	GetSecurityProfile(name string) (*SecurityProfile, error)
	UpdateSecurityProfile(name string, sp *SecurityProfile) error
	DeleteSecurityProfile(name string) error
	ListSecurityProfiles(opts common.ListOptions) (*common.ListResult[SecurityProfile], error)
}

type ViolationOps interface {
	ListActiveViolations(thingName string) ([]*ViolationEvent, error)
	ListViolationEvents(opts common.ListOptions, securityProfileName, thingName string) ([]*ViolationEvent, error)
}

type DomainConfigOps interface {
	CreateDomainConfiguration(dc *DomainConfiguration) (*DomainConfiguration, error)
	GetDomainConfiguration(name string) (*DomainConfiguration, error)
	UpdateDomainConfiguration(name string, dc *DomainConfiguration) error
	DeleteDomainConfiguration(name string) error
	ListDomainConfigurations(opts common.ListOptions) (*common.ListResult[DomainConfiguration], error)
}

type IndexingConfigOps interface {
	GetIndexingConfiguration() (*IndexingConfiguration, error)
	UpdateIndexingConfiguration(ic *IndexingConfiguration) error
}

type ProvisioningTemplateVersionOps interface {
	CreateProvisioningTemplateVersion(name string, v *ProvisioningTemplateVersion) (*ProvisioningTemplateVersion, error)
	SetDefaultProvisioningTemplateVersion(name string, versionID int64) error
	GetProvisioningTemplateVersion(name, versionID string) (*ProvisioningTemplateVersion, error)
	DeleteProvisioningTemplateVersion(name, versionID string) error
	ListProvisioningTemplateVersions(name string, opts common.ListOptions) ([]*ProvisioningTemplateVersion, error)
}

type TagOps interface {
	ListTags(resourceARN string) (map[string]string, error)
	TagResource(resourceARN string, tags map[string]string) error
	UntagResource(resourceARN string, tagKeys []string) error
	DeleteAllTags(resourceARN string) error
}

type PolicyAttachmentOps interface {
	AttachPolicyToPrincipal(policyName, principal string) error
	DetachPolicyFromPrincipal(policyName, principal string) error
	ListPrincipalsForPolicy(policyName string) ([]string, error)
	ListPoliciesForPrincipal(principal string) ([]string, error)
}

type ThingPrincipalOps interface {
	AttachThingPrincipal(thingName, principal string) error
	// AttachThingPrincipalExclusive records an EXCLUSIVE_THING attachment
	// (provisioning-template ThingPrincipalType): the principal may attach
	// to this thing only.
	AttachThingPrincipalExclusive(thingName, principal string) error
	DetachThingPrincipal(thingName, principal string) error
	ListPrincipalsForThing(thingName string) ([]string, error)
	AddThingToThingGroup(thingName, groupName string) error
	RemoveThingFromThingGroup(thingName, groupName string) error
	ListThingsInGroup(groupName string) ([]string, error)
	ListGroupsForThing(thingName string) ([]string, error)
	ListThingsForPrincipal(principal string) ([]string, error)
}

type MetaOps interface {
	Storage() storage.BasicStorage
	GetAccountID() string
	GetRegion() string
}
