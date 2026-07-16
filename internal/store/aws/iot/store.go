package iot

import (
	"context"
	"sync"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Bucket base names. These are the single source of truth for every bucket
// name in this package; both NewIotStore (which constructs the BaseStore
// wrappers) and the transactional helpers in thing_store.go /
// policy_store.go / thing_principal_store.go / thing_group_store.go append
// the region suffix ("-{region}") to the SAME constant, so a rename can
// never diverge between the read path (BaseStore) and the write path
// (txn.Bucket). See AGENTS.md事故12 (DRY store instance discipline).
const (
	bucketThings             = "iot-things"
	bucketThingTypes         = "iot-thing-types"
	bucketThingGroups        = "iot-thing-groups"
	bucketCertificates       = "iot-certificates"
	bucketPolicies           = "iot-policies"
	bucketRules              = "iot-rules"
	bucketJobs               = "iot-jobs"
	bucketShadows            = "iot-shadows"
	bucketRoleAliases        = "iot-role-aliases"
	bucketPolicyAttach       = "iot-policy-attach"
	bucketThingPrincipal     = "iot-thing-principal"
	bucketPrincipalThing     = "iot-principal-thing"
	bucketPolicyPrincipal    = "iot-policy-principal"
	bucketThingGroupMember   = "iot-thing-group-member"
	bucketGroupThingMember   = "iot-group-thing-member"
	bucketThingBillingMember = "iot-thing-billing-member"
	bucketBillingThingMember = "iot-billing-thing-member"
	bucketBillingGroups      = "iot-billing-groups"
	bucketAuthorizers        = "iot-authorizers"
	bucketProvisioningTpls   = "iot-provisioning-templates"
	bucketDetectorModels     = "iot-detector-models"
	bucketInputs             = "iot-inputs"
	bucketAlarmModels        = "iot-alarm-models"
	bucketGenericKV          = "iot-generic-kv"
	bucketSecurityProfiles   = "iot-security-profiles"
	bucketDomainConfigs      = "iot-domain-configs"
	bucketViolationEvents    = "iot-violation-events"
	bucketIndexingConfig     = "iot-indexing-config"
)

type IotStore struct {
	thingsBase             *common.BaseStore
	thingTypesBase         *common.BaseStore
	thingGroupsBase        *common.BaseStore
	certsBase              *common.BaseStore
	policiesBase           *common.BaseStore
	rulesBase              *common.BaseStore
	jobsBase               *common.BaseStore
	shadowsBase            *common.BaseStore
	roleAliasBase          *common.BaseStore
	policyAttachBase       *common.BaseStore
	thingPrincipalBase     *common.BaseStore
	principalThingBase     *common.BaseStore
	policyPrincipalBase    *common.BaseStore
	thingGroupMemberBase   *common.BaseStore
	groupThingMemberBase   *common.BaseStore
	thingBillingMemberBase *common.BaseStore
	billingThingMemberBase *common.BaseStore
	billingGroupsBase      *common.BaseStore
	authorizersBase        *common.BaseStore
	templatesBase          *common.BaseStore
	detectorModelsBase     *common.BaseStore
	inputsBase             *common.BaseStore
	alarmModelsBase        *common.BaseStore
	genericKVBase          *common.BaseStore
	securityProfilesBase   *common.BaseStore
	domainConfigsBase      *common.BaseStore
	violationEventsBase    *common.BaseStore
	indexingConfigBase     *common.BaseStore
	// ProtoStore instances provide built-in existence checks for
	// Create (AlreadyExist), Update (NotFound), and DeleteIfExists
	// (NotFound).  Each wraps the corresponding BaseStore above.
	thingPS           *common.ProtoStore[Thing]
	thingTypePS       *common.ProtoStore[ThingType]
	thingGroupPS      *common.ProtoStore[ThingGroup]
	billingGroupPS    *common.ProtoStore[BillingGroup]
	certificatePS     *common.ProtoStore[Certificate]
	policyPS          *common.ProtoStore[Policy]
	topicRulePS       *common.ProtoStore[TopicRule]
	jobPS             *common.ProtoStore[Job]
	authorizerPS      *common.ProtoStore[Authorizer]
	roleAliasPS       *common.ProtoStore[RoleAlias]
	provisioningTplPS *common.ProtoStore[ProvisioningTemplate]
	securityProfilePS *common.ProtoStore[SecurityProfile]
	domainConfigPS    *common.ProtoStore[DomainConfiguration]
	violationEventPS  *common.ProtoStore[ViolationEvent]
	*common.TagStore
	storage      storage.BasicStorage
	ts           storage.TransactionalStorage
	rs           string
	arnBuilder   *svcarn.ARNBuilder
	accountID    string
	region       string
	mu           sync.RWMutex
	stateMachine *DetectorStateMachine
	// shadowLocker provides per-shadow-key mutual exclusion for atomic
	// read-merge-write sequences, replacing the store-wide s.mu. This
	// allows concurrent shadow updates for different things to proceed
	// without blocking each other.
	shadowLocker *common.KeyLocker
}

func NewIotStore(store storage.BasicStorage, accountID, region string, onAction func(string, string, string, map[string]interface{})) *IotStore {
	bp := "-" + region
	s := &IotStore{
		thingsBase:             common.NewBaseStore(store.Bucket(bucketThings+bp), bucketThings),
		thingTypesBase:         common.NewBaseStore(store.Bucket(bucketThingTypes+bp), bucketThingTypes),
		thingGroupsBase:        common.NewBaseStore(store.Bucket(bucketThingGroups+bp), bucketThingGroups),
		certsBase:              common.NewBaseStore(store.Bucket(bucketCertificates+bp), bucketCertificates),
		policiesBase:           common.NewBaseStore(store.Bucket(bucketPolicies+bp), bucketPolicies),
		rulesBase:              common.NewBaseStore(store.Bucket(bucketRules+bp), bucketRules),
		jobsBase:               common.NewBaseStore(store.Bucket(bucketJobs+bp), bucketJobs),
		shadowsBase:            common.NewBaseStore(store.Bucket(bucketShadows+bp), bucketShadows),
		roleAliasBase:          common.NewBaseStore(store.Bucket(bucketRoleAliases+bp), bucketRoleAliases),
		policyAttachBase:       common.NewBaseStore(store.Bucket(bucketPolicyAttach+bp), bucketPolicyAttach),
		thingPrincipalBase:     common.NewBaseStore(store.Bucket(bucketThingPrincipal+bp), bucketThingPrincipal),
		principalThingBase:     common.NewBaseStore(store.Bucket(bucketPrincipalThing+bp), bucketPrincipalThing),
		policyPrincipalBase:    common.NewBaseStore(store.Bucket(bucketPolicyPrincipal+bp), bucketPolicyPrincipal),
		thingGroupMemberBase:   common.NewBaseStore(store.Bucket(bucketThingGroupMember+bp), bucketThingGroupMember),
		groupThingMemberBase:   common.NewBaseStore(store.Bucket(bucketGroupThingMember+bp), bucketGroupThingMember),
		thingBillingMemberBase: common.NewBaseStore(store.Bucket(bucketThingBillingMember+bp), bucketThingBillingMember),
		billingThingMemberBase: common.NewBaseStore(store.Bucket(bucketBillingThingMember+bp), bucketBillingThingMember),
		billingGroupsBase:      common.NewBaseStore(store.Bucket(bucketBillingGroups+bp), bucketBillingGroups),
		authorizersBase:        common.NewBaseStore(store.Bucket(bucketAuthorizers+bp), bucketAuthorizers),
		templatesBase:          common.NewBaseStore(store.Bucket(bucketProvisioningTpls+bp), bucketProvisioningTpls),
		detectorModelsBase:     common.NewBaseStore(store.Bucket(bucketDetectorModels+bp), bucketDetectorModels),
		inputsBase:             common.NewBaseStore(store.Bucket(bucketInputs+bp), bucketInputs),
		alarmModelsBase:        common.NewBaseStore(store.Bucket(bucketAlarmModels+bp), bucketAlarmModels),
		genericKVBase:          common.NewBaseStore(store.Bucket(bucketGenericKV+bp), bucketGenericKV),
		securityProfilesBase:   common.NewBaseStore(store.Bucket(bucketSecurityProfiles+bp), bucketSecurityProfiles),
		domainConfigsBase:      common.NewBaseStore(store.Bucket(bucketDomainConfigs+bp), bucketDomainConfigs),
		violationEventsBase:    common.NewBaseStore(store.Bucket(bucketViolationEvents+bp), bucketViolationEvents),
		indexingConfigBase:     common.NewBaseStore(store.Bucket(bucketIndexingConfig+bp), bucketIndexingConfig),
		TagStore:               common.NewTagStoreWithRegion(store, "iot", region),
		storage:                store,
		arnBuilder:             svcarn.NewARNBuilder(accountID, region),
		accountID:              accountID,
		region:                 region,
		rs:                     bp,
		shadowLocker:           &common.KeyLocker{},
	}
	s.ts, _ = store.(storage.TransactionalStorage)
	initProtoStores(s)
	s.stateMachine = NewDetectorStateMachine(onAction)
	return s
}

// regionalStores ensures a single IotStore instance per (accountID, region)
// pair. All services (iot, iotevents, MQTT broker auth) must use
// GetOrCreateStore to avoid duplicate store instances and state machine
// divergence, mirroring the IAM/STS/admin_auth pattern (see
// GetOrCreateGlobalStore in store/aws/iam/store.go).
var regionalStores sync.Map

// GetOrCreateStore returns the cached IotStore for the given account and
// region, creating it on first access. The detector action callback is
// left nil; callers that need detector action dispatch must invoke
// SetActionCallback on the returned store exactly once after creation.
func GetOrCreateStore(store storage.BasicStorage, accountID, region string) *IotStore {
	key := accountID + "/" + region
	if cached, ok := regionalStores.Load(key); ok {
		if typed, ok := cached.(*IotStore); ok {
			return typed
		}
	}
	newStore := NewIotStore(store, accountID, region, nil)
	actual, _ := regionalStores.LoadOrStore(key, newStore)
	return actual.(*IotStore)
}

// SetActionCallback installs the detector action dispatch callback on the
// state machine. Intended to be called at most once during process
// initialisation; subsequent calls overwrite the previous callback. Safe
// for concurrent EvaluateEvent readers via the state machine's internal
// RWMutex.
func (s *IotStore) SetActionCallback(fn func(modelName, key, actionType string, payload map[string]interface{})) {
	if s.stateMachine != nil {
		s.stateMachine.SetActionCallback(fn)
	}
}

func initProtoStores(s *IotStore) {
	s.thingPS = common.NewProtoStore(common.ProtoStoreConfig[Thing]{
		Store: s.thingsBase, NewProto: func() proto.Message { return &pb.Thing{} },
		ToDomain:    func(m proto.Message) (*Thing, error) { return ProtoToThing(m.(*pb.Thing)), nil },
		ToProto:     func(t *Thing) (proto.Message, error) { return ThingToProto(t), nil },
		IDFunc:      func(t *Thing) string { return t.ThingName },
		NotFoundErr: ErrThingNotFound, AlreadyExist: ErrThingAlreadyExists,
	})
	s.thingTypePS = common.NewProtoStore(common.ProtoStoreConfig[ThingType]{
		Store: s.thingTypesBase, NewProto: func() proto.Message { return &pb.ThingType{} },
		ToDomain:    func(m proto.Message) (*ThingType, error) { return ProtoToThingType(m.(*pb.ThingType)), nil },
		ToProto:     func(t *ThingType) (proto.Message, error) { return ThingTypeToProto(t), nil },
		IDFunc:      func(t *ThingType) string { return t.ThingTypeName },
		NotFoundErr: ErrThingTypeNotFound, AlreadyExist: ErrThingTypeAlreadyExists,
	})
	s.thingGroupPS = common.NewProtoStore(common.ProtoStoreConfig[ThingGroup]{
		Store: s.thingGroupsBase, NewProto: func() proto.Message { return &pb.ThingGroup{} },
		ToDomain:    func(m proto.Message) (*ThingGroup, error) { return ProtoToThingGroup(m.(*pb.ThingGroup)), nil },
		ToProto:     func(g *ThingGroup) (proto.Message, error) { return ThingGroupToProto(g), nil },
		IDFunc:      func(g *ThingGroup) string { return g.GroupName },
		NotFoundErr: ErrThingGroupNotFound, AlreadyExist: ErrThingGroupAlreadyExists,
	})
	s.billingGroupPS = common.NewProtoStore(common.ProtoStoreConfig[BillingGroup]{
		Store: s.billingGroupsBase, NewProto: func() proto.Message { return &pb.BillingGroup{} },
		ToDomain:    func(m proto.Message) (*BillingGroup, error) { return ProtoToBillingGroup(m.(*pb.BillingGroup)), nil },
		ToProto:     func(b *BillingGroup) (proto.Message, error) { return BillingGroupToProto(b), nil },
		IDFunc:      func(b *BillingGroup) string { return b.GroupName },
		NotFoundErr: ErrBillingGroupNotFound, AlreadyExist: ErrBillingGroupAlreadyExists,
	})
	s.certificatePS = common.NewProtoStore(common.ProtoStoreConfig[Certificate]{
		Store: s.certsBase, NewProto: func() proto.Message { return &pb.Certificate{} },
		ToDomain:    func(m proto.Message) (*Certificate, error) { return ProtoToCertificate(m.(*pb.Certificate)), nil },
		ToProto:     func(c *Certificate) (proto.Message, error) { return CertificateToProto(c), nil },
		IDFunc:      func(c *Certificate) string { return c.CertificateID },
		NotFoundErr: ErrCertificateNotFound, AlreadyExist: ErrCertificateAlreadyExists,
	})
	s.policyPS = common.NewProtoStore(common.ProtoStoreConfig[Policy]{
		Store: s.policiesBase, NewProto: func() proto.Message { return &pb.Policy{} },
		ToDomain:    func(m proto.Message) (*Policy, error) { return ProtoToPolicy(m.(*pb.Policy)), nil },
		ToProto:     func(p *Policy) (proto.Message, error) { return PolicyToProto(p), nil },
		IDFunc:      func(p *Policy) string { return p.PolicyName },
		NotFoundErr: ErrPolicyNotFound, AlreadyExist: ErrPolicyAlreadyExists,
	})
	s.topicRulePS = common.NewProtoStore(common.ProtoStoreConfig[TopicRule]{
		Store: s.rulesBase, NewProto: func() proto.Message { return &pb.TopicRule{} },
		ToDomain:    func(m proto.Message) (*TopicRule, error) { return ProtoToRule(m.(*pb.TopicRule)), nil },
		ToProto:     func(r *TopicRule) (proto.Message, error) { return RuleToProto(r), nil },
		IDFunc:      func(r *TopicRule) string { return r.RuleName },
		NotFoundErr: ErrRuleNotFound, AlreadyExist: ErrRuleAlreadyExists,
	})
	s.jobPS = common.NewProtoStore(common.ProtoStoreConfig[Job]{
		Store: s.jobsBase, NewProto: func() proto.Message { return &pb.Job{} },
		ToDomain:    func(m proto.Message) (*Job, error) { return ProtoToJob(m.(*pb.Job)), nil },
		ToProto:     func(j *Job) (proto.Message, error) { return JobToProto(j), nil },
		IDFunc:      func(j *Job) string { return j.JobID },
		NotFoundErr: ErrJobNotFound, AlreadyExist: ErrJobAlreadyExists,
	})
	s.authorizerPS = common.NewProtoStore(common.ProtoStoreConfig[Authorizer]{
		Store: s.authorizersBase, NewProto: func() proto.Message { return &pb.Authorizer{} },
		ToDomain:    func(m proto.Message) (*Authorizer, error) { return ProtoToAuthorizer(m.(*pb.Authorizer)), nil },
		ToProto:     func(a *Authorizer) (proto.Message, error) { return AuthorizerToProto(a), nil },
		IDFunc:      func(a *Authorizer) string { return a.AuthorizerName },
		NotFoundErr: ErrAuthorizerNotFound, AlreadyExist: ErrAuthorizerAlreadyExists,
	})
	s.roleAliasPS = common.NewProtoStore(common.ProtoStoreConfig[RoleAlias]{
		Store: s.roleAliasBase, NewProto: func() proto.Message { return &pb.RoleAlias{} },
		ToDomain:    func(m proto.Message) (*RoleAlias, error) { return ProtoToRoleAlias(m.(*pb.RoleAlias)), nil },
		ToProto:     func(r *RoleAlias) (proto.Message, error) { return RoleAliasToProto(r), nil },
		IDFunc:      func(r *RoleAlias) string { return r.RoleAlias },
		NotFoundErr: ErrRoleAliasNotFound, AlreadyExist: ErrRoleAliasAlreadyExists,
	})
	s.provisioningTplPS = common.NewProtoStore(common.ProtoStoreConfig[ProvisioningTemplate]{
		Store: s.templatesBase, NewProto: func() proto.Message { return &pb.ProvisioningTemplate{} },
		ToDomain: func(m proto.Message) (*ProvisioningTemplate, error) {
			return ProtoToProvisioningTemplate(m.(*pb.ProvisioningTemplate))
		},
		ToProto:     func(t *ProvisioningTemplate) (proto.Message, error) { return ProvisioningTemplateToProto(t) },
		IDFunc:      func(t *ProvisioningTemplate) string { return t.TemplateName },
		NotFoundErr: ErrTemplateNotFound, AlreadyExist: ErrTemplateAlreadyExists,
	})
	s.securityProfilePS = common.NewProtoStore(common.ProtoStoreConfig[SecurityProfile]{
		Store: s.securityProfilesBase, NewProto: func() proto.Message { return &pb.SecurityProfile{} },
		ToDomain: func(m proto.Message) (*SecurityProfile, error) {
			return ProtoToSecurityProfile(m.(*pb.SecurityProfile)), nil
		},
		ToProto: func(sp *SecurityProfile) (proto.Message, error) {
			return SecurityProfileToProto(sp)
		},
		IDFunc:      func(sp *SecurityProfile) string { return sp.SecurityProfileName },
		NotFoundErr: ErrSecurityProfileNotFound, AlreadyExist: ErrSecurityProfileAlreadyExists,
	})
	s.domainConfigPS = common.NewProtoStore(common.ProtoStoreConfig[DomainConfiguration]{
		Store: s.domainConfigsBase, NewProto: func() proto.Message { return &pb.DomainConfiguration{} },
		ToDomain: func(m proto.Message) (*DomainConfiguration, error) {
			return ProtoToDomainConfiguration(m.(*pb.DomainConfiguration)), nil
		},
		ToProto: func(dc *DomainConfiguration) (proto.Message, error) {
			return DomainConfigurationToProto(dc)
		},
		IDFunc:      func(dc *DomainConfiguration) string { return dc.DomainConfigurationName },
		NotFoundErr: ErrDomainConfigurationNotFound, AlreadyExist: ErrDomainConfigurationAlreadyExists,
	})
	s.violationEventPS = common.NewProtoStore(common.ProtoStoreConfig[ViolationEvent]{
		Store: s.violationEventsBase, NewProto: func() proto.Message { return &pb.ViolationEvent{} },
		ToDomain: func(m proto.Message) (*ViolationEvent, error) {
			return ProtoToViolationEvent(m.(*pb.ViolationEvent)), nil
		},
		ToProto: func(v *ViolationEvent) (proto.Message, error) {
			return ViolationEventToProto(v)
		},
		IDFunc: func(v *ViolationEvent) string { return v.ViolationID },
	})
}

func (s *IotStore) Storage() storage.BasicStorage {
	return s.storage
}

func (s *IotStore) GetAccountID() string {
	return s.accountID
}
func (s *IotStore) GetRegion() string {
	return s.region
}

// LoadDetectorModel registers a detector model in the state machine for
// event evaluation. No-op if no state machine is initialised.
func (s *IotStore) LoadDetectorModel(dm *DetectorModel) {
	if s.stateMachine != nil {
		s.stateMachine.LoadModel(dm)
	}
}

// UnloadModel removes a detector model and its instances from the state machine.
func (s *IotStore) UnloadModel(name string) {
	if s.stateMachine != nil {
		s.stateMachine.UnloadModel(name)
	}
}

// BatchEvaluate processes a batch of input messages against all loaded detector
// models via the state machine. Returns error entries for failed evaluations.
func (s *IotStore) BatchEvaluate(ctx context.Context, messages []InputMessage) []map[string]interface{} {
	if s.stateMachine == nil {
		return nil
	}
	return s.stateMachine.BatchEvaluate(ctx, messages)
}

// StateMachine returns the underlying DetectorStateMachine, or nil if not initialised.
func (s *IotStore) StateMachine() *DetectorStateMachine {
	return s.stateMachine
}
