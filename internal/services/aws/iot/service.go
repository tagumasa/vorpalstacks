// Package iot provides AWS IoT Core service operations for vorpalstacks.
package iot

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/services/aws/iot/broker"
	"vorpalstacks/internal/services/aws/iot/ca"
	"vorpalstacks/internal/services/aws/iot/rules"
	"vorpalstacks/internal/services/aws/iot/rules/actions"
	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ErrNotInitialised is returned when the service is accessed before Init.
var ErrNotInitialised = errors.New("iot service not initialised")

// IoTServiceDeps holds all dependencies required for IoTService initialisation.
type IoTServiceDeps struct {
	StorageManager   *storage.RegionStorageManager
	CA               *ca.CertificateAuthority
	Broker           *broker.Broker
	EventBus         eventbus.Bus
	ActionDispatcher *actions.Dispatcher
}

// IoTService provides AWS IoT Core operations for managing things,
// certificates, policies, shadows, rules, and jobs.
type IoTService struct {
	accountID    string
	stores       sync.Map
	deps         IoTServiceDeps
	executor     *rules.Executor
	stateMachine *detectorStateMachine
	once         sync.Once
	initialised  bool
}

// NewIoTService creates a new IoT Core service instance.
// Call Init before using the service.
func NewIoTService(accountID string) *IoTService {
	return &IoTService{
		accountID: accountID,
	}
}

// Init sets all dependencies atomically, guarded by sync.Once.
// Must be called before RegisterHandlers or any operation.
// Calling Init multiple times is safe and idempotent.
func (s *IoTService) Init(deps IoTServiceDeps) {
	s.once.Do(func() {
		s.deps = deps

		if deps.ActionDispatcher != nil && deps.StorageManager != nil {
			dispatchFn := s.makeActionDispatcher(deps.ActionDispatcher)
			s.executor = rules.NewExecutor(dispatchFn, slog.Default())
			s.executor.Start()
			s.hydrateRules()

			s.stateMachine = newDetectorStateMachine(func(modelName, key, actionType string, payload map[string]interface{}) {
				slog.Debug("detector action", "model", modelName, "key", key, "action", actionType)
			})
		}

		s.initialised = true
	})
}

func (s *IoTService) hydrateRules() {
	if s.executor == nil || s.deps.StorageManager == nil {
		return
	}
	stores := s.deps.StorageManager.ListRegions()
	if len(stores) == 0 {
		return
	}
	defaultRegion := stores[0]
	st, err := s.deps.StorageManager.GetStorage(defaultRegion)
	if err != nil {
		slog.Warn("failed to get storage for IoT rules hydrate", "region", defaultRegion, "error", err)
		return
	}
	store := iotstore.NewIotStore(st, s.accountID, defaultRegion)

	var allRules []iotstore.TopicRule
	var opts storecommon.ListOptions
	for {
		result, err := store.ListRules(opts)
		if err != nil {
			slog.Warn("failed to list rules for hydrate", "error", err)
			break
		}
		for _, r := range result.Items {
			allRules = append(allRules, *r)
		}
		if result.NextMarker == "" {
			break
		}
		opts.Marker = result.NextMarker
	}

	for _, r := range allRules {
		if !r.RuleDisabled && r.SQL != "" {
			if err := s.executor.AddRule(r.RuleName, r.TopicPattern, r.SQL, actionsToList(r.Actions)); err != nil {
				slog.Warn("failed to hydrate rule", "rule", r.RuleName, "error", err)
			}
		}
	}
	if len(allRules) > 0 {
		slog.Info("hydrated IoT rules executor", "total", len(allRules), "active", s.executor.RulesCount())
	}
}

func (s *IoTService) Shutdown() {
	if s.executor != nil {
		s.executor.Stop()
	}
}

func (s *IoTService) Executor() *rules.Executor {
	return s.executor
}

func (s *IoTService) makeActionDispatcher(d *actions.Dispatcher) rules.ActionDispatcher {
	return func(ruleName, topic string, actionList []map[string]interface{}, payload map[string]interface{}) error {
		ctx := context.Background()
		for _, action := range actionList {
			for actionType, cfg := range action {
				ac := &actions.ActionConfig{Type: actionType}
				if m, ok := cfg.(map[string]interface{}); ok {
					ac.Extra = m
				}
				if err := d.Dispatch(ctx, ac, topic, payload); err != nil {
					slog.Warn("rule action dispatch failed", "rule", ruleName, "type", actionType, "error", err)
				}
			}
		}
		return nil
	}
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
	d.RegisterHandlerForService("iot", "CreateSecurityProfile", s.CreateSecurityProfile)
	d.RegisterHandlerForService("iot", "DescribeSecurityProfile", s.DescribeSecurityProfile)
	d.RegisterHandlerForService("iot", "UpdateSecurityProfile", s.UpdateSecurityProfile)
	d.RegisterHandlerForService("iot", "DeleteSecurityProfile", s.DeleteSecurityProfile)
	d.RegisterHandlerForService("iot", "ListSecurityProfiles", s.ListSecurityProfiles)

	// DetectorModel operations (IoT Events integration)
	d.RegisterHandlerForService("iot", "CreateDetectorModel", s.CreateDetectorModel)
	d.RegisterHandlerForService("iot", "DescribeDetectorModel", s.DescribeDetectorModel)
	d.RegisterHandlerForService("iot", "UpdateDetectorModel", s.UpdateDetectorModel)
	d.RegisterHandlerForService("iot", "DeleteDetectorModel", s.DeleteDetectorModel)
	d.RegisterHandlerForService("iot", "ListDetectorModels", s.ListDetectorModels)

	// Input operations (IoT Events integration)
	d.RegisterHandlerForService("iot", "CreateInput", s.CreateInput)
	d.RegisterHandlerForService("iot", "DescribeInput", s.DescribeInput)
	d.RegisterHandlerForService("iot", "UpdateInput", s.UpdateInput)
	d.RegisterHandlerForService("iot", "DeleteInput", s.DeleteInput)
	d.RegisterHandlerForService("iot", "ListInputs", s.ListInputs)

	d.RegisterHandlerForService("iot", "BatchPutMessage", s.BatchPutMessage)
}
