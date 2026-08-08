// Package iot provides AWS IoT Core service operations for vorpalstacks.
package iot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/services/aws/iot/broker"
	"vorpalstacks/internal/services/aws/iot/ca"
	"vorpalstacks/internal/services/aws/iot/rules"
	"vorpalstacks/internal/services/aws/iot/rules/actions"
	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// IoTServiceDeps holds all dependencies required for IoTService
// initialisation. Brokers, CAs and Dispatchers are keyed by region so
// each region gets an independent data plane (MQTT broker, signing CA
// and rule-action dispatcher), mirroring AWS IoT Core's fully regional
// design.
type IoTServiceDeps struct {
	StorageManager *storage.RegionStorageManager
	CAs            map[string]*ca.CertificateAuthority
	Brokers        map[string]*broker.Broker
	EventBus       eventbus.Bus
	Dispatchers    map[string]*actions.Dispatcher
}

// IoTService provides AWS IoT Core operations for managing things,
// certificates, policies, shadows, rules, and jobs. Each region owns an
// independent broker, rule executor and action dispatcher.
type IoTService struct {
	accountID   string
	deps        IoTServiceDeps
	executors   map[string]*rules.Executor
	brokers     map[string]*broker.Broker
	dispatchers map[string]*actions.Dispatcher
	taskEngine  *taskEngine
	throttle    *throttleLimiter
	once        sync.Once
	initialised bool
}

// throttleLimiter provides a simple sliding-window rate limiter for
// data-plane Publish operations. AWS IoT throttles these operations per
// account; the limiter returns ErrThrottling when the configured threshold
// is exceeded within the window.
type throttleLimiter struct {
	mu     sync.Mutex
	window time.Duration
	limit  int
	hits   []time.Time
}

func newThrottleLimiter(limit int, window time.Duration) *throttleLimiter {
	return &throttleLimiter{window: window, limit: limit}
}

// allow returns true if the request is within the rate limit, false if
// throttled. Thread-safe.
func (t *throttleLimiter) allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-t.window)
	idx := 0
	for ; idx < len(t.hits); idx++ {
		if t.hits[idx].After(cutoff) {
			break
		}
	}
	t.hits = t.hits[idx:]
	if len(t.hits) >= t.limit {
		return false
	}
	t.hits = append(t.hits, now)
	return true
}

// NewIoTService creates a new IoT Core service instance.
// Call Init before using the service.
func NewIoTService(accountID string) *IoTService {
	return &IoTService{
		accountID: accountID,
		throttle:  newThrottleLimiter(10000, time.Second),
	}
}

// Init sets all dependencies atomically, guarded by sync.Once.
// Must be called before RegisterHandlers or any operation.
// Calling Init multiple times is safe and idempotent.
//
// A per-region rule executor and broker message handler are wired for
// every region in deps.StorageManager.ListRegions(), so each region's
// MQTT messages are evaluated against that region's topic rules.
func (s *IoTService) Init(deps IoTServiceDeps) {
	s.once.Do(func() {
		s.deps = deps
		s.executors = make(map[string]*rules.Executor)
		s.brokers = deps.Brokers
		if s.brokers == nil {
			s.brokers = map[string]*broker.Broker{}
		}
		s.dispatchers = deps.Dispatchers
		if s.dispatchers == nil {
			s.dispatchers = map[string]*actions.Dispatcher{}
		}

		if deps.StorageManager != nil {
			for _, region := range deps.StorageManager.ListRegions() {
				dispatcher := s.dispatchers[region]
				var dispatchFn rules.ActionDispatcher
				if dispatcher != nil {
					dispatchFn = s.makeActionDispatcher(region, dispatcher)
				}
				executor := rules.NewExecutor(dispatchFn, slog.Default())
				executor.Start()
				s.executors[region] = executor

				if brk := s.brokers[region]; brk != nil {
					brk.SetMessageHandler(func(_ string, topic string, payload []byte) {
						executor.OnMessage(topic, payload)
					})
				}

				s.hydrateRulesForRegion(region, executor)
			}

			// Start the background task engine to transition IN_PROGRESS
			// tasks (Detect Mitigation, On-Demand Audit) to COMPLETED.
			s.taskEngine = newTaskEngine(deps.StorageManager, s.accountID)
			s.taskEngine.Start(context.Background())
		}

		s.initialised = true
	})
}

// hydrateRulesForRegion loads the topic rules for a single region into the
// given executor. Each region's rules stay isolated in their own executor,
// matching AWS IoT Core's per-region rule engine.
func (s *IoTService) hydrateRulesForRegion(region string, executor *rules.Executor) {
	if executor == nil || s.deps.StorageManager == nil {
		return
	}
	st, err := s.deps.StorageManager.GetStorage(region)
	if err != nil {
		slog.Warn("failed to get storage for IoT rules hydrate", "region", region, "error", err)
		return
	}
	store := iotstore.GetOrCreateStore(st, s.accountID, region)

	var opts storecommon.ListOptions
	var count int
	for {
		result, err := store.ListRules(opts)
		if err != nil {
			slog.Warn("failed to list rules for hydrate", "region", region, "error", err)
			break
		}
		for _, r := range result.Items {
			if !r.RuleDisabled && r.SQL != "" {
				if err := executor.AddRule(r.RuleName, r.TopicPattern, r.SQL, actionsToList(r.Actions), r.ErrorAction); err != nil {
					slog.Warn("failed to hydrate rule", "rule", r.RuleName, "region", region, "error", err)
					continue
				}
				count++
			}
		}
		if result.NextMarker == "" {
			break
		}
		opts.Marker = result.NextMarker
	}
	if count > 0 {
		slog.Info("hydrated IoT rules executor", "region", region, "total", count, "active", executor.RulesCount())
	}
}

// Shutdown stops every regional rule executor and MQTT broker. Safe to
// call multiple times; executor.Stop and broker.Stop are themselves idempotent.
func (s *IoTService) Shutdown() {
	if s.taskEngine != nil {
		s.taskEngine.Stop()
	}
	for _, executor := range s.executors {
		executor.Stop()
	}
	for _, brk := range s.brokers {
		_ = brk.Stop()
	}
}

// ExecutorForRegion returns the rule executor for the given region, or nil.
func (s *IoTService) ExecutorForRegion(region string) *rules.Executor {
	return s.executors[region]
}

// BrokerForRegion returns the MQTT broker for the given region, or nil.
func (s *IoTService) BrokerForRegion(region string) *broker.Broker {
	return s.brokers[region]
}

// executorForReq returns the rule executor for the request's region, or
// nil when the region has no executor (e.g. tests that never called Init).
func (s *IoTService) executorForReq(reqCtx *request.RequestContext) *rules.Executor {
	if reqCtx == nil {
		return nil
	}
	return s.executors[reqCtx.GetRegion()]
}

// caForReq returns the signing CA for the request's region, or nil.
func (s *IoTService) caForReq(reqCtx *request.RequestContext) *ca.CertificateAuthority {
	if reqCtx == nil {
		return nil
	}
	return s.deps.CAs[reqCtx.GetRegion()]
}

// brokerForReq returns the MQTT broker for the request's region, or nil.
func (s *IoTService) brokerForReq(reqCtx *request.RequestContext) *broker.Broker {
	if reqCtx == nil {
		return nil
	}
	return s.brokers[reqCtx.GetRegion()]
}

func (s *IoTService) makeActionDispatcher(region string, d *actions.Dispatcher) rules.ActionDispatcher {
	return func(ctx context.Context, ruleName, topic string, actionList []map[string]interface{}, errorAction map[string]interface{}, payload map[string]interface{}, iteration int) error {
		// Embed iteration count for republish chain tracking. dispatchRepublish
		// reads this field and increments it before publishing.
		if iteration > 0 {
			payload["_iotRuleIteration"] = iteration
		}
		hadError := false
		for _, action := range actionList {
			for actionType, cfg := range action {
				var ac *actions.ActionConfig
				if m, ok := cfg.(map[string]interface{}); ok {
					ac = actions.NewActionConfigFromMap(actionType, m)
				} else {
					ac = &actions.ActionConfig{Type: actionType}
				}
				ac.Region = region
				if err := d.Dispatch(ctx, ac, topic, payload); err != nil {
					slog.Warn("rule action dispatch failed", "rule", ruleName, "region", region, "type", actionType, "error", err)
					hadError = true
				}
			}
		}
		// When any primary action fails, dispatch the error action (AWS IoT
		// dead-letter behaviour) so the user can route failures to SQS/SNS.
		if hadError && len(errorAction) > 0 {
			for eaType, eaCfg := range errorAction {
				if m, ok := eaCfg.(map[string]interface{}); ok {
					eaConfig := actions.NewActionConfigFromMap(eaType, m)
					eaConfig.Region = region
					if err := d.Dispatch(ctx, eaConfig, topic, payload); err != nil {
						slog.Warn("rule error action dispatch failed", "rule", ruleName, "region", region, "type", eaType, "error", err)
					}
				}
			}
		}
		if hadError {
			return fmt.Errorf("rule %s: one or more actions failed", ruleName)
		}
		return nil
	}
}

// WireActionCallbacksForRegion connects the rule action dispatcher
// callbacks for the given region to their concrete implementations.
// Republish routes to that region's broker so messages stay within the
// region, matching AWS IoT Core's per-region rule engine.
func (s *IoTService) WireActionCallbacksForRegion(region string) {
	d := s.dispatchers[region]
	if d == nil {
		return
	}

	if brk := s.brokers[region]; brk != nil {
		d.SetRepublishFn(func(ctx context.Context, topic string, payload map[string]interface{}) error {
			data, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			return brk.Publish(topic, data)
		})
	}

	d.SetHTTPPostFn(func(ctx context.Context, url string, payload []byte) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode >= 400 {
			return fmt.Errorf("http action: %s returned status %d", url, resp.StatusCode)
		}
		return nil
	})
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
	d.RegisterHandlerForService("iot", "Publish", s.Publish)

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
	d.RegisterHandlerForService("iot", "CreateProvisioningClaim", s.CreateProvisioningClaim)

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

	// Device Management Extensions
	d.RegisterHandlerForService("iot", "CreateDynamicThingGroup", s.CreateDynamicThingGroup)
	d.RegisterHandlerForService("iot", "DeleteDynamicThingGroup", s.DeleteDynamicThingGroup)
	d.RegisterHandlerForService("iot", "UpdateDynamicThingGroup", s.UpdateDynamicThingGroup)
	d.RegisterHandlerForService("iot", "AddThingToBillingGroup", s.AddThingToBillingGroup)
	d.RegisterHandlerForService("iot", "RemoveThingFromBillingGroup", s.RemoveThingFromBillingGroup)
	d.RegisterHandlerForService("iot", "ListThingsInBillingGroup", s.ListThingsInBillingGroup)
	d.RegisterHandlerForService("iot", "RegisterThing", s.RegisterThing)
	d.RegisterHandlerForService("iot", "StartThingRegistrationTask", s.StartThingRegistrationTask)
	d.RegisterHandlerForService("iot", "StopThingRegistrationTask", s.StopThingRegistrationTask)
	d.RegisterHandlerForService("iot", "DescribeThingRegistrationTask", s.DescribeThingRegistrationTask)
	d.RegisterHandlerForService("iot", "ListThingRegistrationTasks", s.ListThingRegistrationTasks)
	d.RegisterHandlerForService("iot", "ListThingRegistrationTaskReports", s.ListThingRegistrationTaskReports)

	// Certificates & CA
	d.RegisterHandlerForService("iot", "RegisterCACertificate", s.RegisterCACertificate)
	d.RegisterHandlerForService("iot", "DescribeCACertificate", s.DescribeCACertificate)
	d.RegisterHandlerForService("iot", "ListCACertificates", s.ListCACertificates)
	d.RegisterHandlerForService("iot", "UpdateCACertificate", s.UpdateCACertificate)
	d.RegisterHandlerForService("iot", "DeleteCACertificate", s.DeleteCACertificate)
	d.RegisterHandlerForService("iot", "ListCertificatesByCA", s.ListCertificatesByCA)
	d.RegisterHandlerForService("iot", "RegisterCertificateWithoutCA", s.RegisterCertificateWithoutCA)
	d.RegisterHandlerForService("iot", "TransferCertificate", s.TransferCertificate)
	d.RegisterHandlerForService("iot", "AcceptCertificateTransfer", s.AcceptCertificateTransfer)
	d.RegisterHandlerForService("iot", "CancelCertificateTransfer", s.CancelCertificateTransfer)
	d.RegisterHandlerForService("iot", "RejectCertificateTransfer", s.RejectCertificateTransfer)
	d.RegisterHandlerForService("iot", "ListOutgoingCertificates", s.ListOutgoingCertificates)

	// Policy/Authorizer add-ons
	d.RegisterHandlerForService("iot", "CreatePolicyVersion", s.CreatePolicyVersion)
	d.RegisterHandlerForService("iot", "DeletePolicyVersion", s.DeletePolicyVersion)
	d.RegisterHandlerForService("iot", "SetDefaultPolicyVersion", s.SetDefaultPolicyVersion)
	d.RegisterHandlerForService("iot", "ListTargetsForPolicy", s.ListTargetsForPolicy)
	d.RegisterHandlerForService("iot", "AttachPrincipalPolicy", s.AttachPrincipalPolicy)
	d.RegisterHandlerForService("iot", "DetachPrincipalPolicy", s.DetachPrincipalPolicy)
	d.RegisterHandlerForService("iot", "SetDefaultAuthorizer", s.SetDefaultAuthorizer)
	d.RegisterHandlerForService("iot", "ClearDefaultAuthorizer", s.ClearDefaultAuthorizer)
	d.RegisterHandlerForService("iot", "DescribeDefaultAuthorizer", s.DescribeDefaultAuthorizer)
	d.RegisterHandlerForService("iot", "TestInvokeAuthorizer", s.TestInvokeAuthorizer)

	// Jobs ecosystem
	d.RegisterHandlerForService("iot", "CreateJobTemplate", s.CreateJobTemplate)
	d.RegisterHandlerForService("iot", "DeleteJobTemplate", s.DeleteJobTemplate)
	d.RegisterHandlerForService("iot", "DescribeJobTemplate", s.DescribeJobTemplate)
	d.RegisterHandlerForService("iot", "ListJobTemplates", s.ListJobTemplates)
	d.RegisterHandlerForService("iot", "DescribeManagedJobTemplate", s.DescribeManagedJobTemplate)
	d.RegisterHandlerForService("iot", "ListManagedJobTemplates", s.ListManagedJobTemplates)
	d.RegisterHandlerForService("iot", "CancelJobExecution", s.CancelJobExecution)
	d.RegisterHandlerForService("iot", "DeleteJobExecution", s.DeleteJobExecution)
	d.RegisterHandlerForService("iot", "DescribeJobExecution", s.DescribeJobExecution)
	d.RegisterHandlerForService("iot", "ListJobExecutionsForJob", s.ListJobExecutionsForJob)
	d.RegisterHandlerForService("iot", "ListJobExecutionsForThing", s.ListJobExecutionsForThing)
	d.RegisterHandlerForService("iot", "CreateOTAUpdate", s.CreateOTAUpdate)
	d.RegisterHandlerForService("iot", "DeleteOTAUpdate", s.DeleteOTAUpdate)
	d.RegisterHandlerForService("iot", "GetOTAUpdate", s.GetOTAUpdate)
	d.RegisterHandlerForService("iot", "ListOTAUpdates", s.ListOTAUpdates)

	// Stream (MQTT file delivery) + Registration code
	d.RegisterHandlerForService("iot", "CreateStream", s.CreateStream)
	d.RegisterHandlerForService("iot", "DeleteStream", s.DeleteStream)
	d.RegisterHandlerForService("iot", "DescribeStream", s.DescribeStream)
	d.RegisterHandlerForService("iot", "ListStreams", s.ListStreams)
	d.RegisterHandlerForService("iot", "UpdateStream", s.UpdateStream)
	d.RegisterHandlerForService("iot", "GetRegistrationCode", s.GetRegistrationCode)
	d.RegisterHandlerForService("iot", "DeleteRegistrationCode", s.DeleteRegistrationCode)
	d.RegisterHandlerForService("iot", "CreateProvisioningTemplateVersion", s.CreateProvisioningTemplateVersion)
	d.RegisterHandlerForService("iot", "DeleteProvisioningTemplateVersion", s.DeleteProvisioningTemplateVersion)
	d.RegisterHandlerForService("iot", "DescribeProvisioningTemplateVersion", s.DescribeProvisioningTemplateVersion)
	d.RegisterHandlerForService("iot", "UpdateThingGroupsForThing", s.UpdateThingGroupsForThing)
	d.RegisterHandlerForService("iot", "GetThingConnectivityData", s.GetThingConnectivityData)
	d.RegisterHandlerForService("iot", "ListPrincipalThingsV2", s.ListPrincipalThingsV2)
	d.RegisterHandlerForService("iot", "ListThingPrincipalsV2", s.ListThingPrincipalsV2)

	// Package/SBOM (IoT Device Management — software package distribution)
	d.RegisterHandlerForService("iot", "CreatePackage", s.CreatePackage)
	d.RegisterHandlerForService("iot", "GetPackage", s.GetPackage)
	d.RegisterHandlerForService("iot", "UpdatePackage", s.UpdatePackage)
	d.RegisterHandlerForService("iot", "DeletePackage", s.DeletePackage)
	d.RegisterHandlerForService("iot", "ListPackages", s.ListPackages)
	d.RegisterHandlerForService("iot", "CreatePackageVersion", s.CreatePackageVersion)
	d.RegisterHandlerForService("iot", "GetPackageVersion", s.GetPackageVersion)
	d.RegisterHandlerForService("iot", "UpdatePackageVersion", s.UpdatePackageVersion)
	d.RegisterHandlerForService("iot", "DeletePackageVersion", s.DeletePackageVersion)
	d.RegisterHandlerForService("iot", "ListPackageVersions", s.ListPackageVersions)
	d.RegisterHandlerForService("iot", "GetPackageConfiguration", s.GetPackageConfiguration)
	d.RegisterHandlerForService("iot", "UpdatePackageConfiguration", s.UpdatePackageConfiguration)
	d.RegisterHandlerForService("iot", "AssociateSbomWithPackageVersion", s.AssociateSbomWithPackageVersion)
	d.RegisterHandlerForService("iot", "DisassociateSbomFromPackageVersion", s.DisassociateSbomFromPackageVersion)
	d.RegisterHandlerForService("iot", "ListSbomValidationResults", s.ListSbomValidationResults)

	// CertificateProvider (self-managed certificate signing via Lambda)
	d.RegisterHandlerForService("iot", "CreateCertificateProvider", s.CreateCertificateProvider)
	d.RegisterHandlerForService("iot", "DescribeCertificateProvider", s.DescribeCertificateProvider)
	d.RegisterHandlerForService("iot", "UpdateCertificateProvider", s.UpdateCertificateProvider)
	d.RegisterHandlerForService("iot", "DeleteCertificateProvider", s.DeleteCertificateProvider)
	d.RegisterHandlerForService("iot", "ListCertificateProviders", s.ListCertificateProviders)

	// Command (remote command execution metadata management)
	d.RegisterHandlerForService("iot", "CreateCommand", s.CreateCommand)
	d.RegisterHandlerForService("iot", "GetCommand", s.GetCommand)
	d.RegisterHandlerForService("iot", "UpdateCommand", s.UpdateCommand)
	d.RegisterHandlerForService("iot", "DeleteCommand", s.DeleteCommand)
	d.RegisterHandlerForService("iot", "ListCommands", s.ListCommands)
	d.RegisterHandlerForService("iot", "GetCommandExecution", s.GetCommandExecution)
	d.RegisterHandlerForService("iot", "DeleteCommandExecution", s.DeleteCommandExecution)
	d.RegisterHandlerForService("iot", "ListCommandExecutions", s.ListCommandExecutions)

	// Security/Audit/Fleet/Logging/TopicRuleDest bulk handlers
	d.RegisterHandlerForService("iot", "AttachSecurityProfile", s.AttachSecurityProfile)
	d.RegisterHandlerForService("iot", "DetachSecurityProfile", s.DetachSecurityProfile)
	d.RegisterHandlerForService("iot", "ListSecurityProfilesForTarget", s.ListSecurityProfilesForTarget)
	d.RegisterHandlerForService("iot", "ListTargetsForSecurityProfile", s.ListTargetsForSecurityProfile)
	d.RegisterHandlerForService("iot", "PutVerificationStateOnViolation", s.PutVerificationStateOnViolation)
	d.RegisterHandlerForService("iot", "CreateCustomMetric", s.CreateCustomMetric)
	d.RegisterHandlerForService("iot", "DeleteCustomMetric", s.DeleteCustomMetric)
	d.RegisterHandlerForService("iot", "DescribeCustomMetric", s.DescribeCustomMetric)
	d.RegisterHandlerForService("iot", "ListCustomMetrics", s.ListCustomMetrics)
	d.RegisterHandlerForService("iot", "UpdateCustomMetric", s.UpdateCustomMetric)
	d.RegisterHandlerForService("iot", "CreateDimension", s.CreateDimension)
	d.RegisterHandlerForService("iot", "DeleteDimension", s.DeleteDimension)
	d.RegisterHandlerForService("iot", "DescribeDimension", s.DescribeDimension)
	d.RegisterHandlerForService("iot", "ListDimensions", s.ListDimensions)
	d.RegisterHandlerForService("iot", "UpdateDimension", s.UpdateDimension)
	d.RegisterHandlerForService("iot", "CreateMitigationAction", s.CreateMitigationAction)
	d.RegisterHandlerForService("iot", "DeleteMitigationAction", s.DeleteMitigationAction)
	d.RegisterHandlerForService("iot", "DescribeMitigationAction", s.DescribeMitigationAction)
	d.RegisterHandlerForService("iot", "ListMitigationActions", s.ListMitigationActions)
	d.RegisterHandlerForService("iot", "UpdateMitigationAction", s.UpdateMitigationAction)
	d.RegisterHandlerForService("iot", "StartDetectMitigationActionsTask", s.StartDetectMitigationActionsTask)
	d.RegisterHandlerForService("iot", "CancelDetectMitigationActionsTask", s.CancelDetectMitigationActionsTask)
	d.RegisterHandlerForService("iot", "DescribeDetectMitigationActionsTask", s.DescribeDetectMitigationActionsTask)
	d.RegisterHandlerForService("iot", "ListDetectMitigationActionsExecutions", s.ListDetectMitigationActionsExecutions)
	d.RegisterHandlerForService("iot", "ListDetectMitigationActionsTasks", s.ListDetectMitigationActionsTasks)
	d.RegisterHandlerForService("iot", "DescribeAccountAuditConfiguration", s.DescribeAccountAuditConfiguration)
	d.RegisterHandlerForService("iot", "UpdateAccountAuditConfiguration", s.UpdateAccountAuditConfiguration)
	d.RegisterHandlerForService("iot", "DeleteAccountAuditConfiguration", s.DeleteAccountAuditConfiguration)
	d.RegisterHandlerForService("iot", "StartOnDemandAuditTask", s.StartOnDemandAuditTask)
	d.RegisterHandlerForService("iot", "CancelAuditTask", s.CancelAuditTask)
	d.RegisterHandlerForService("iot", "DescribeAuditTask", s.DescribeAuditTask)
	d.RegisterHandlerForService("iot", "ListAuditTasks", s.ListAuditTasks)
	d.RegisterHandlerForService("iot", "DescribeAuditFinding", s.DescribeAuditFinding)
	d.RegisterHandlerForService("iot", "ListAuditFindings", s.ListAuditFindings)
	d.RegisterHandlerForService("iot", "ListRelatedResourcesForAuditFinding", s.ListRelatedResourcesForAuditFinding)
	d.RegisterHandlerForService("iot", "CreateAuditSuppression", s.CreateAuditSuppression)
	d.RegisterHandlerForService("iot", "DeleteAuditSuppression", s.DeleteAuditSuppression)
	d.RegisterHandlerForService("iot", "DescribeAuditSuppression", s.DescribeAuditSuppression)
	d.RegisterHandlerForService("iot", "ListAuditSuppressions", s.ListAuditSuppressions)
	d.RegisterHandlerForService("iot", "UpdateAuditSuppression", s.UpdateAuditSuppression)
	d.RegisterHandlerForService("iot", "StartAuditMitigationActionsTask", s.StartAuditMitigationActionsTask)
	d.RegisterHandlerForService("iot", "CancelAuditMitigationActionsTask", s.CancelAuditMitigationActionsTask)
	d.RegisterHandlerForService("iot", "DescribeAuditMitigationActionsTask", s.DescribeAuditMitigationActionsTask)
	d.RegisterHandlerForService("iot", "ListAuditMitigationActionsExecutions", s.ListAuditMitigationActionsExecutions)
	d.RegisterHandlerForService("iot", "ListAuditMitigationActionsTasks", s.ListAuditMitigationActionsTasks)
	d.RegisterHandlerForService("iot", "CreateScheduledAudit", s.CreateScheduledAudit)
	d.RegisterHandlerForService("iot", "DeleteScheduledAudit", s.DeleteScheduledAudit)
	d.RegisterHandlerForService("iot", "DescribeScheduledAudit", s.DescribeScheduledAudit)
	d.RegisterHandlerForService("iot", "ListScheduledAudits", s.ListScheduledAudits)
	d.RegisterHandlerForService("iot", "UpdateScheduledAudit", s.UpdateScheduledAudit)
	d.RegisterHandlerForService("iot", "AssociateTargetsWithJob", s.AssociateTargetsWithJob)
	d.RegisterHandlerForService("iot", "TestAuthorization", s.TestAuthorization)
	d.RegisterHandlerForService("iot", "DescribeIndex", s.DescribeIndex)
	d.RegisterHandlerForService("iot", "ListIndices", s.ListIndices)
	d.RegisterHandlerForService("iot", "SearchIndex", s.SearchIndex)
	d.RegisterHandlerForService("iot", "GetBucketsAggregation", s.GetBucketsAggregation)
	d.RegisterHandlerForService("iot", "GetCardinality", s.GetCardinality)
	d.RegisterHandlerForService("iot", "GetPercentiles", s.GetPercentiles)
	d.RegisterHandlerForService("iot", "GetStatistics", s.GetStatistics)
	d.RegisterHandlerForService("iot", "ListMetricValues", s.ListMetricValues)
	d.RegisterHandlerForService("iot", "CreateFleetMetric", s.CreateFleetMetric)
	d.RegisterHandlerForService("iot", "DeleteFleetMetric", s.DeleteFleetMetric)
	d.RegisterHandlerForService("iot", "DescribeFleetMetric", s.DescribeFleetMetric)
	d.RegisterHandlerForService("iot", "ListFleetMetrics", s.ListFleetMetrics)
	d.RegisterHandlerForService("iot", "UpdateFleetMetric", s.UpdateFleetMetric)
	d.RegisterHandlerForService("iot", "GetV2LoggingOptions", s.GetV2LoggingOptions)
	d.RegisterHandlerForService("iot", "SetV2LoggingOptions", s.SetV2LoggingOptions)
	d.RegisterHandlerForService("iot", "DeleteV2LoggingLevel", s.DeleteV2LoggingLevel)
	d.RegisterHandlerForService("iot", "ListV2LoggingLevels", s.ListV2LoggingLevels)
	d.RegisterHandlerForService("iot", "SetV2LoggingLevel", s.SetV2LoggingLevel)
	d.RegisterHandlerForService("iot", "GetLoggingOptions", s.GetLoggingOptions)
	d.RegisterHandlerForService("iot", "SetLoggingOptions", s.SetLoggingOptions)
	d.RegisterHandlerForService("iot", "DescribeEventConfigurations", s.DescribeEventConfigurations)
	d.RegisterHandlerForService("iot", "UpdateEventConfigurations", s.UpdateEventConfigurations)
	d.RegisterHandlerForService("iot", "DescribeEncryptionConfiguration", s.DescribeEncryptionConfiguration)
	d.RegisterHandlerForService("iot", "UpdateEncryptionConfiguration", s.UpdateEncryptionConfiguration)
	d.RegisterHandlerForService("iot", "CreateTopicRuleDestination", s.CreateTopicRuleDestination)
	d.RegisterHandlerForService("iot", "DeleteTopicRuleDestination", s.DeleteTopicRuleDestination)
	d.RegisterHandlerForService("iot", "GetTopicRuleDestination", s.GetTopicRuleDestination)
	d.RegisterHandlerForService("iot", "ListTopicRuleDestinations", s.ListTopicRuleDestinations)
	d.RegisterHandlerForService("iot", "UpdateTopicRuleDestination", s.UpdateTopicRuleDestination)
	d.RegisterHandlerForService("iot", "ConfirmTopicRuleDestination", s.ConfirmTopicRuleDestination)
}

// IoTAuthBrokerProvider implements broker.AuthProvider by delegating to the
// IoTService for the region-scoped store and CA. Each regional broker is
// given a provider with its own Region so certificate and policy lookups
// hit that region's singleton IotStore, never stores[0].
type IoTAuthBrokerProvider struct {
	Service *IoTService
	Region  string
}

func (p *IoTAuthBrokerProvider) GetCA() *ca.CertificateAuthority {
	return p.Service.deps.CAs[p.Region]
}

func (p *IoTAuthBrokerProvider) GetStore() iotstore.IotStoreInterface {
	if p.Service.deps.StorageManager == nil {
		slog.Error("iot broker auth: storage manager not available, denying all MQTT connections", "region", p.Region)
		return nil
	}
	st, err := p.Service.deps.StorageManager.GetStorage(p.Region)
	if err != nil {
		slog.Error("iot broker auth: storage unavailable, denying all MQTT connections", "region", p.Region, "error", err)
		return nil
	}
	return iotstore.GetOrCreateStore(st, p.Service.accountID, p.Region)
}

// CertificatePrincipal builds the ARN for a certificate ID so the auth
// hook can look up policy attachments (which are keyed by ARN).
func (p *IoTAuthBrokerProvider) CertificatePrincipal(certID string) string {
	return iotstore.BuildCertificateARN(p.Service.accountID, p.Region, certID)
}
