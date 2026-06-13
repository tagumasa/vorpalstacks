// Package iot provides AWS IoT Core service operations for vorpalstacks.
package iot

import (
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/services/aws/iot/broker"
	"vorpalstacks/internal/services/aws/iot/ca"
	"vorpalstacks/internal/services/aws/iot/rules/actions"
)

// IoTService provides AWS IoT Core operations for managing things,
// certificates, policies, shadows, rules, and jobs.
type IoTService struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	stores         sync.Map
	ca             *ca.CertificateAuthority
	broker         *broker.Broker
	bus            eventbus.Bus
	actionDisp     *actions.Dispatcher
}

// NewIoTService creates a new IoT Core service instance.
func NewIoTService(accountID string) *IoTService {
	return &IoTService{
		accountID: accountID,
	}
}

// SetStorageManager injects the region storage manager for admin console access.
func (s *IoTService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// SetCA injects the built-in certificate authority for certificate operations.
func (s *IoTService) SetCA(certAuth *ca.CertificateAuthority) {
	s.ca = certAuth
}

// SetBroker injects the MQTT broker for shadow delta publishing.
func (s *IoTService) SetBroker(b *broker.Broker) {
	s.broker = b
}

// SetEventBus injects the event bus for cross-service action dispatching.
func (s *IoTService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
}

// SetActionDispatcher injects the rules action dispatcher.
func (s *IoTService) SetActionDispatcher(d *actions.Dispatcher) {
	s.actionDisp = d
}

// RegisterHandlers registers all IoT operation handlers with the dispatcher.
func (s *IoTService) RegisterHandlers(d handler.Registrar) {
	// Thing Registry
	d.RegisterHandlerForService("iot", "CreateThing", s.CreateThing)
	d.RegisterHandlerForService("iot", "DescribeThing", s.DescribeThing)
	d.RegisterHandlerForService("iot", "UpdateThing", s.UpdateThing)
	d.RegisterHandlerForService("iot", "DeleteThing", s.DeleteThing)
	d.RegisterHandlerForService("iot", "ListThings", s.ListThings)
	d.RegisterHandlerForService("iot", "ListThingsForThingType", s.ListThingsForThingType)

	// Certificate
	d.RegisterHandlerForService("iot", "CreateKeysAndCertificate", s.CreateKeysAndCertificate)
	d.RegisterHandlerForService("iot", "DescribeCertificate", s.DescribeCertificate)
	d.RegisterHandlerForService("iot", "UpdateCertificate", s.UpdateCertificate)
	d.RegisterHandlerForService("iot", "DeleteCertificate", s.DeleteCertificate)
	d.RegisterHandlerForService("iot", "ListCertificates", s.ListCertificates)
	d.RegisterHandlerForService("iot", "RegisterCertificate", s.RegisterCertificate)
	d.RegisterHandlerForService("iot", "CreateCertificateFromCsr", s.CreateCertificateFromCsr)

	// Policy
	d.RegisterHandlerForService("iot", "CreatePolicy", s.CreatePolicy)
	d.RegisterHandlerForService("iot", "GetPolicy", s.GetPolicy)
	d.RegisterHandlerForService("iot", "DeletePolicy", s.DeletePolicy)
	d.RegisterHandlerForService("iot", "ListPolicies", s.ListPolicies)
	d.RegisterHandlerForService("iot", "AttachPolicy", s.AttachPolicy)
	d.RegisterHandlerForService("iot", "DetachPolicy", s.DetachPolicy)
	d.RegisterHandlerForService("iot", "AttachThingPrincipal", s.AttachThingPrincipal)
	d.RegisterHandlerForService("iot", "DetachThingPrincipal", s.DetachThingPrincipal)
	d.RegisterHandlerForService("iot", "ListPolicyPrincipals", s.ListPolicyPrincipals)
	d.RegisterHandlerForService("iot", "ListPrincipalPolicies", s.ListPrincipalPolicies)
	d.RegisterHandlerForService("iot", "ListAttachedPolicies", s.ListPrincipalPolicies)
	d.RegisterHandlerForService("iot", "ListThingPrincipals", s.ListThingPrincipals)
	d.RegisterHandlerForService("iot", "ListPrincipalThings", s.ListPrincipalThings)
	d.RegisterHandlerForService("iot", "GetEffectivePolicies", s.GetEffectivePolicies)
	d.RegisterHandlerForService("iot", "ListPolicyVersions", s.ListPolicyVersions)
	d.RegisterHandlerForService("iot", "GetPolicyVersion", s.GetPolicyVersion)

	// Shadow
	d.RegisterHandlerForService("iot", "GetThingShadow", s.GetThingShadow)
	d.RegisterHandlerForService("iot", "UpdateThingShadow", s.UpdateThingShadow)
	d.RegisterHandlerForService("iot", "DeleteThingShadow", s.DeleteThingShadow)
	d.RegisterHandlerForService("iot", "ListNamedShadowsForThing", s.ListNamedShadowsForThing)

	// RoleAlias
	d.RegisterHandlerForService("iot", "CreateRoleAlias", s.CreateRoleAlias)
	d.RegisterHandlerForService("iot", "DescribeRoleAlias", s.DescribeRoleAlias)
	d.RegisterHandlerForService("iot", "UpdateRoleAlias", s.UpdateRoleAlias)
	d.RegisterHandlerForService("iot", "DeleteRoleAlias", s.DeleteRoleAlias)
	d.RegisterHandlerForService("iot", "ListRoleAliases", s.ListRoleAliases)

	// Topic Rules
	d.RegisterHandlerForService("iot", "CreateTopicRule", s.CreateTopicRule)
	d.RegisterHandlerForService("iot", "DescribeTopicRule", s.DescribeTopicRule)
	d.RegisterHandlerForService("iot", "ReplaceTopicRule", s.ReplaceTopicRule)
	d.RegisterHandlerForService("iot", "DeleteTopicRule", s.DeleteTopicRule)
	d.RegisterHandlerForService("iot", "ListTopicRules", s.ListTopicRules)
	d.RegisterHandlerForService("iot", "EnableTopicRule", s.EnableTopicRule)
	d.RegisterHandlerForService("iot", "DisableTopicRule", s.DisableTopicRule)
	d.RegisterHandlerForService("iot", "GetTopicRule", s.GetTopicRule)

	// Jobs
	d.RegisterHandlerForService("iot", "CreateJob", s.CreateJob)
	d.RegisterHandlerForService("iot", "DescribeJob", s.DescribeJob)
	d.RegisterHandlerForService("iot", "DeleteJob", s.DeleteJob)
	d.RegisterHandlerForService("iot", "ListJobs", s.ListJobs)
	d.RegisterHandlerForService("iot", "CancelJob", s.CancelJob)
	d.RegisterHandlerForService("iot", "GetJobDocument", s.GetJobDocument)
	d.RegisterHandlerForService("iot", "UpdateJob", s.UpdateJob)

	// ThingGroup / ThingType / BillingGroup
	d.RegisterHandlerForService("iot", "CreateThingType", s.CreateThingType)
	d.RegisterHandlerForService("iot", "DescribeThingType", s.DescribeThingType)
	d.RegisterHandlerForService("iot", "UpdateThingType", s.UpdateThingType)
	d.RegisterHandlerForService("iot", "DeleteThingType", s.DeleteThingType)
	d.RegisterHandlerForService("iot", "ListThingTypes", s.ListThingTypes)
	d.RegisterHandlerForService("iot", "DeprecateThingType", s.DeprecateThingType)
	d.RegisterHandlerForService("iot", "CreateThingGroup", s.CreateThingGroup)
	d.RegisterHandlerForService("iot", "DescribeThingGroup", s.DescribeThingGroup)
	d.RegisterHandlerForService("iot", "UpdateThingGroup", s.UpdateThingGroup)
	d.RegisterHandlerForService("iot", "DeleteThingGroup", s.DeleteThingGroup)
	d.RegisterHandlerForService("iot", "ListThingGroups", s.ListThingGroups)
	d.RegisterHandlerForService("iot", "AddThingToThingGroup", s.AddThingToThingGroup)
	d.RegisterHandlerForService("iot", "RemoveThingFromThingGroup", s.RemoveThingFromThingGroup)
	d.RegisterHandlerForService("iot", "ListThingsInThingGroup", s.ListThingsInThingGroup)
	d.RegisterHandlerForService("iot", "ListThingGroupsForThing", s.ListThingGroupsForThing)
	d.RegisterHandlerForService("iot", "CreateBillingGroup", s.CreateBillingGroup)
	d.RegisterHandlerForService("iot", "DescribeBillingGroup", s.DescribeBillingGroup)
	d.RegisterHandlerForService("iot", "UpdateBillingGroup", s.UpdateBillingGroup)
	d.RegisterHandlerForService("iot", "DeleteBillingGroup", s.DeleteBillingGroup)
	d.RegisterHandlerForService("iot", "ListBillingGroups", s.ListBillingGroups)

	// Authorizer / ProvisioningTemplate
	d.RegisterHandlerForService("iot", "CreateAuthorizer", s.CreateAuthorizer)
	d.RegisterHandlerForService("iot", "DescribeAuthorizer", s.DescribeAuthorizer)
	d.RegisterHandlerForService("iot", "UpdateAuthorizer", s.UpdateAuthorizer)
	d.RegisterHandlerForService("iot", "DeleteAuthorizer", s.DeleteAuthorizer)
	d.RegisterHandlerForService("iot", "ListAuthorizers", s.ListAuthorizers)
	d.RegisterHandlerForService("iot", "CreateProvisioningTemplate", s.CreateProvisioningTemplate)
	d.RegisterHandlerForService("iot", "DescribeProvisioningTemplate", s.DescribeProvisioningTemplate)
	d.RegisterHandlerForService("iot", "UpdateProvisioningTemplate", s.UpdateProvisioningTemplate)
	d.RegisterHandlerForService("iot", "DeleteProvisioningTemplate", s.DeleteProvisioningTemplate)
	d.RegisterHandlerForService("iot", "ListProvisioningTemplates", s.ListProvisioningTemplates)
	d.RegisterHandlerForService("iot", "ListProvisioningTemplateVersions", s.ListProvisioningTemplateVersions)

	// Endpoint / DomainConfiguration
	d.RegisterHandlerForService("iot", "DescribeEndpoint", s.DescribeEndpoint)
	d.RegisterHandlerForService("iot", "GetIndexingConfiguration", s.GetIndexingConfiguration)
	d.RegisterHandlerForService("iot", "UpdateIndexingConfiguration", s.UpdateIndexingConfiguration)
	d.RegisterHandlerForService("iot", "DescribeDomainConfiguration", s.DescribeDomainConfiguration)
	d.RegisterHandlerForService("iot", "CreateDomainConfiguration", s.CreateDomainConfiguration)
	d.RegisterHandlerForService("iot", "UpdateDomainConfiguration", s.UpdateDomainConfiguration)
	d.RegisterHandlerForService("iot", "DeleteDomainConfiguration", s.DeleteDomainConfiguration)
	d.RegisterHandlerForService("iot", "ListDomainConfigurations", s.ListDomainConfigurations)

	// Tags / Security Profile
	d.RegisterHandlerForService("iot", "TagResource", s.TagResource)
	d.RegisterHandlerForService("iot", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("iot", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("iot", "ListActiveViolations", s.ListActiveViolations)
	d.RegisterHandlerForService("iot", "ListViolationEvents", s.ListViolationEvents)
	d.RegisterHandlerForService("iot", "GetBehaviorModelTrainingSummaries", s.GetBehaviorModelTrainingSummaries)
	d.RegisterHandlerForService("iot", "ValidateSecurityProfileBehaviors", s.ValidateSecurityProfileBehaviors)
	d.RegisterHandlerForService("iot", "DescribeSecurityProfile", s.DescribeSecurityProfile)
	d.RegisterHandlerForService("iot", "UpdateSecurityProfile", s.UpdateSecurityProfile)
}
