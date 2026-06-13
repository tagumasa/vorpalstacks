package iot

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

type IotStoreInterface interface {
	// Things
	CreateThing(thing *Thing) (*Thing, error)
	GetThing(thingName string) (*Thing, error)
	UpdateThing(thing *Thing) error
	DeleteThing(thingName string) error
	ListThings(opts common.ListOptions) (*common.ListResult[Thing], error)
	ListThingsForThingType(thingTypeName string, opts common.ListOptions) (*common.ListResult[Thing], error)

	// ThingTypes
	CreateThingType(tt *ThingType) (*ThingType, error)
	GetThingType(name string) (*ThingType, error)
	UpdateThingType(tt *ThingType) error
	DeleteThingType(name string) error
	ListThingTypes(opts common.ListOptions) (*common.ListResult[ThingType], error)

	// ThingGroups
	CreateThingGroup(group *ThingGroup) (*ThingGroup, error)
	GetThingGroup(name string) (*ThingGroup, error)
	UpdateThingGroup(group *ThingGroup) error
	DeleteThingGroup(name string) error
	ListThingGroups(opts common.ListOptions) (*common.ListResult[ThingGroup], error)

	// BillingGroups
	CreateBillingGroup(bg *BillingGroup) (*BillingGroup, error)
	GetBillingGroup(name string) (*BillingGroup, error)
	UpdateBillingGroup(bg *BillingGroup) error
	DeleteBillingGroup(name string) error
	ListBillingGroups(opts common.ListOptions) (*common.ListResult[BillingGroup], error)

	// Certificates
	CreateCertificate(cert *Certificate) (*Certificate, error)
	GetCertificate(certificateID string) (*Certificate, error)
	UpdateCertificate(cert *Certificate) error
	DeleteCertificate(certificateID string) error
	ListCertificates(opts common.ListOptions) (*common.ListResult[Certificate], error)

	// Policies
	CreatePolicy(policy *Policy) (*Policy, error)
	GetPolicy(policyName string) (*Policy, error)
	DeletePolicy(policyName string) error
	ListPolicies(opts common.ListOptions) (*common.ListResult[Policy], error)

	// TopicRules
	CreateRule(rule *TopicRule) (*TopicRule, error)
	GetRule(ruleName string) (*TopicRule, error)
	UpdateRule(rule *TopicRule) error
	DeleteRule(ruleName string) error
	ListRules(opts common.ListOptions) (*common.ListResult[TopicRule], error)

	// RoleAliases
	CreateRoleAlias(ra *RoleAlias) (*RoleAlias, error)
	GetRoleAlias(alias string) (*RoleAlias, error)
	UpdateRoleAlias(ra *RoleAlias) error
	DeleteRoleAlias(alias string) error
	ListRoleAliases(opts common.ListOptions) (*common.ListResult[RoleAlias], error)

	// Jobs
	CreateJob(job *Job) (*Job, error)
	GetJob(jobID string) (*Job, error)
	UpdateJob(job *Job) error
	DeleteJob(jobID string) error
	ListJobs(opts common.ListOptions) (*common.ListResult[Job], error)

	// Shadows
	GetShadow(thingName, shadowName string) (*ShadowDocument, error)
	PutShadow(thingName, shadowName string, doc *ShadowDocument) error
	PutShadowWithVersion(thingName, shadowName string, doc *ShadowDocument, clientVersion int64) error
	DeleteShadow(thingName, shadowName string) error
	ListShadowNames(thingName string) ([]string, error)

	// Authorizers
	CreateAuthorizer(a *Authorizer) (*Authorizer, error)
	GetAuthorizer(name string) (*Authorizer, error)
	UpdateAuthorizer(a *Authorizer) error
	DeleteAuthorizer(name string) error
	ListAuthorizers(opts common.ListOptions) (*common.ListResult[Authorizer], error)

	// ProvisioningTemplates
	CreateProvisioningTemplate(t *ProvisioningTemplate) (*ProvisioningTemplate, error)
	GetProvisioningTemplate(name string) (*ProvisioningTemplate, error)
	UpdateProvisioningTemplate(t *ProvisioningTemplate) error
	DeleteProvisioningTemplate(name string) error
	ListProvisioningTemplates(opts common.ListOptions) (*common.ListResult[ProvisioningTemplate], error)

	// DetectorModels
	CreateDetectorModel(d *DetectorModel) (*DetectorModel, error)
	GetDetectorModel(name string) (*DetectorModel, error)
	UpdateDetectorModel(d *DetectorModel) error
	DeleteDetectorModel(name string) error
	ListDetectorModels(opts common.ListOptions) (*common.ListResult[DetectorModel], error)

	// Inputs
	CreateInput(i *Input) (*Input, error)
	GetInput(name string) (*Input, error)
	UpdateInput(i *Input) error
	DeleteInput(name string) error
	ListInputs(opts common.ListOptions) (*common.ListResult[Input], error)

	// Tags (via embedded TagStore)
	ListTags(resourceARN string) (map[string]string, error)
	TagResource(resourceARN string, tags map[string]string) error
	UntagResource(resourceARN string, tagKeys []string) error

	// Policy Attachments
	AttachPolicyToPrincipal(policyName, principal string) error
	DetachPolicyFromPrincipal(policyName, principal string) error
	ListPrincipalsForPolicy(policyName string) ([]string, error)
	ListPoliciesForPrincipal(principal string) ([]string, error)

	// Thing Principal Attachments
	AttachThingPrincipal(thingName, principal string) error
	DetachThingPrincipal(thingName, principal string) error
	ListPrincipalsForThing(thingName string) ([]string, error)
	AddThingToThingGroup(thingName, groupName string) error
	RemoveThingFromThingGroup(thingName, groupName string) error
	ListThingsInGroup(groupName string) ([]string, error)
	ListGroupsForThing(thingName string) ([]string, error)
	ListThingsForPrincipal(principal string) ([]string, error)

	Storage() storage.BasicStorage
	GetAccountID() string
	GetRegion() string
}
