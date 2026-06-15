package iot

import (
	"sync"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)
type IotStore struct {
	thingsBase           *common.BaseStore
	thingTypesBase       *common.BaseStore
	thingGroupsBase      *common.BaseStore
	certsBase            *common.BaseStore
	policiesBase         *common.BaseStore
	rulesBase            *common.BaseStore
	jobsBase             *common.BaseStore
	shadowsBase          *common.BaseStore
	roleAliasBase        *common.BaseStore
	policyAttachBase     *common.BaseStore
	thingPrincipalBase   *common.BaseStore
	principalThingBase   *common.BaseStore
	policyPrincipalBase  *common.BaseStore
	thingGroupMemberBase *common.BaseStore
	groupThingMemberBase *common.BaseStore
	billingGroupsBase    *common.BaseStore
	authorizersBase      *common.BaseStore
	templatesBase        *common.BaseStore
	detectorModelsBase   *common.BaseStore
	inputsBase           *common.BaseStore
	securityProfilesBase *common.BaseStore
	domainConfigsBase    *common.BaseStore
	violationEventsBase  *common.BaseStore
	indexingConfigBase   *common.BaseStore
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
	storage    storage.BasicStorage
	ts         storage.TransactionalStorage
	rs         string
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	mu         sync.RWMutex
}


func NewIotStore(store storage.BasicStorage, accountID, region string) *IotStore {
	bp := "-" + region
	s := &IotStore{
		thingsBase:           common.NewBaseStore(store.Bucket("iot-things"+bp), "iot-things"),
		thingTypesBase:       common.NewBaseStore(store.Bucket("iot-thing-types"+bp), "iot-thing-types"),
		thingGroupsBase:      common.NewBaseStore(store.Bucket("iot-thing-groups"+bp), "iot-thing-groups"),
		certsBase:            common.NewBaseStore(store.Bucket("iot-certificates"+bp), "iot-certificates"),
		policiesBase:         common.NewBaseStore(store.Bucket("iot-policies"+bp), "iot-policies"),
		rulesBase:            common.NewBaseStore(store.Bucket("iot-rules"+bp), "iot-rules"),
		jobsBase:             common.NewBaseStore(store.Bucket("iot-jobs"+bp), "iot-jobs"),
		shadowsBase:          common.NewBaseStore(store.Bucket("iot-shadows"+bp), "iot-shadows"),
		roleAliasBase:        common.NewBaseStore(store.Bucket("iot-role-aliases"+bp), "iot-role-aliases"),
		policyAttachBase:     common.NewBaseStore(store.Bucket("iot-policy-attach"+bp), "iot-policy-attach"),
		thingPrincipalBase:   common.NewBaseStore(store.Bucket("iot-thing-principal"+bp), "iot-thing-principal"),
		principalThingBase:   common.NewBaseStore(store.Bucket("iot-principal-thing"+bp), "iot-principal-thing"),
		policyPrincipalBase:  common.NewBaseStore(store.Bucket("iot-policy-principal"+bp), "iot-policy-principal"),
		thingGroupMemberBase: common.NewBaseStore(store.Bucket("iot-thing-group-member"+bp), "iot-thing-group-member"),
		groupThingMemberBase: common.NewBaseStore(store.Bucket("iot-group-thing-member"+bp), "iot-group-thing-member"),
		billingGroupsBase:    common.NewBaseStore(store.Bucket("iot-billing-groups"+bp), "iot-billing-groups"),
		authorizersBase:      common.NewBaseStore(store.Bucket("iot-authorizers"+bp), "iot-authorizers"),
		templatesBase:        common.NewBaseStore(store.Bucket("iot-provisioning-templates"+bp), "iot-provisioning-templates"),
		detectorModelsBase:   common.NewBaseStore(store.Bucket("iot-detector-models"+bp), "iot-detector-models"),
		inputsBase:           common.NewBaseStore(store.Bucket("iot-inputs"+bp), "iot-inputs"),
		securityProfilesBase: common.NewBaseStore(store.Bucket("iot-security-profiles"+bp), "iot-security-profiles"),
		domainConfigsBase:    common.NewBaseStore(store.Bucket("iot-domain-configs"+bp), "iot-domain-configs"),
		violationEventsBase:  common.NewBaseStore(store.Bucket("iot-violation-events"+bp), "iot-violation-events"),
		indexingConfigBase:   common.NewBaseStore(store.Bucket("iot-indexing-config"+bp), "iot-indexing-config"),
		TagStore:             common.NewTagStoreWithRegion(store, "iot", region),
		storage:              store,
		arnBuilder:           svcarn.NewARNBuilder(accountID, region),
		accountID:            accountID,
		region:               region,
		rs:                   bp,
	}
	s.ts, _ = store.(storage.TransactionalStorage)
	initProtoStores(s)
	return s
}

func initProtoStores(s *IotStore) {
	s.thingPS = common.NewProtoStore(common.ProtoStoreConfig[Thing]{
		Store: s.thingsBase, NewProto: func() proto.Message { return &pb.Thing{} },
		ToDomain:    func(m proto.Message) *Thing { return ProtoToThing(m.(*pb.Thing)) },
		ToProto:     func(t *Thing) proto.Message { return ThingToProto(t) },
		IDFunc:      func(t *Thing) string { return t.ThingName },
		NotFoundErr: ErrThingNotFound, AlreadyExist: ErrThingAlreadyExists,
	})
	s.thingTypePS = common.NewProtoStore(common.ProtoStoreConfig[ThingType]{
		Store: s.thingTypesBase, NewProto: func() proto.Message { return &pb.ThingType{} },
		ToDomain:    func(m proto.Message) *ThingType { return ProtoToThingType(m.(*pb.ThingType)) },
		ToProto:     func(t *ThingType) proto.Message { return ThingTypeToProto(t) },
		IDFunc:      func(t *ThingType) string { return t.ThingTypeName },
		NotFoundErr: ErrThingTypeNotFound, AlreadyExist: ErrThingTypeAlreadyExists,
	})
	s.thingGroupPS = common.NewProtoStore(common.ProtoStoreConfig[ThingGroup]{
		Store: s.thingGroupsBase, NewProto: func() proto.Message { return &pb.ThingGroup{} },
		ToDomain:    func(m proto.Message) *ThingGroup { return ProtoToThingGroup(m.(*pb.ThingGroup)) },
		ToProto:     func(g *ThingGroup) proto.Message { return ThingGroupToProto(g) },
		IDFunc:      func(g *ThingGroup) string { return g.GroupName },
		NotFoundErr: ErrThingGroupNotFound, AlreadyExist: ErrThingGroupAlreadyExists,
	})
	s.billingGroupPS = common.NewProtoStore(common.ProtoStoreConfig[BillingGroup]{
		Store: s.billingGroupsBase, NewProto: func() proto.Message { return &pb.BillingGroup{} },
		ToDomain:    func(m proto.Message) *BillingGroup { return ProtoToBillingGroup(m.(*pb.BillingGroup)) },
		ToProto:     func(b *BillingGroup) proto.Message { return BillingGroupToProto(b) },
		IDFunc:      func(b *BillingGroup) string { return b.GroupName },
		NotFoundErr: ErrBillingGroupNotFound, AlreadyExist: ErrBillingGroupAlreadyExists,
	})
	s.certificatePS = common.NewProtoStore(common.ProtoStoreConfig[Certificate]{
		Store: s.certsBase, NewProto: func() proto.Message { return &pb.Certificate{} },
		ToDomain:    func(m proto.Message) *Certificate { return ProtoToCertificate(m.(*pb.Certificate)) },
		ToProto:     func(c *Certificate) proto.Message { return CertificateToProto(c) },
		IDFunc:      func(c *Certificate) string { return c.CertificateID },
		NotFoundErr: ErrCertificateNotFound, AlreadyExist: ErrCertificateAlreadyExists,
	})
	s.policyPS = common.NewProtoStore(common.ProtoStoreConfig[Policy]{
		Store: s.policiesBase, NewProto: func() proto.Message { return &pb.Policy{} },
		ToDomain:    func(m proto.Message) *Policy { return ProtoToPolicy(m.(*pb.Policy)) },
		ToProto:     func(p *Policy) proto.Message { return PolicyToProto(p) },
		IDFunc:      func(p *Policy) string { return p.PolicyName },
		NotFoundErr: ErrPolicyNotFound, AlreadyExist: ErrPolicyAlreadyExists,
	})
	s.topicRulePS = common.NewProtoStore(common.ProtoStoreConfig[TopicRule]{
		Store: s.rulesBase, NewProto: func() proto.Message { return &pb.TopicRule{} },
		ToDomain:    func(m proto.Message) *TopicRule { return ProtoToRule(m.(*pb.TopicRule)) },
		ToProto:     func(r *TopicRule) proto.Message { return RuleToProto(r) },
		IDFunc:      func(r *TopicRule) string { return r.RuleName },
		NotFoundErr: ErrRuleNotFound, AlreadyExist: ErrRuleAlreadyExists,
	})
	s.jobPS = common.NewProtoStore(common.ProtoStoreConfig[Job]{
		Store: s.jobsBase, NewProto: func() proto.Message { return &pb.Job{} },
		ToDomain:    func(m proto.Message) *Job { return ProtoToJob(m.(*pb.Job)) },
		ToProto:     func(j *Job) proto.Message { return JobToProto(j) },
		IDFunc:      func(j *Job) string { return j.JobID },
		NotFoundErr: ErrJobNotFound, AlreadyExist: ErrJobAlreadyExists,
	})
	s.authorizerPS = common.NewProtoStore(common.ProtoStoreConfig[Authorizer]{
		Store: s.authorizersBase, NewProto: func() proto.Message { return &pb.Authorizer{} },
		ToDomain:    func(m proto.Message) *Authorizer { return ProtoToAuthorizer(m.(*pb.Authorizer)) },
		ToProto:     func(a *Authorizer) proto.Message { return AuthorizerToProto(a) },
		IDFunc:      func(a *Authorizer) string { return a.AuthorizerName },
		NotFoundErr: ErrAuthorizerNotFound, AlreadyExist: ErrAuthorizerAlreadyExists,
	})
	s.roleAliasPS = common.NewProtoStore(common.ProtoStoreConfig[RoleAlias]{
		Store: s.roleAliasBase, NewProto: func() proto.Message { return &pb.RoleAlias{} },
		ToDomain:    func(m proto.Message) *RoleAlias { return ProtoToRoleAlias(m.(*pb.RoleAlias)) },
		ToProto:     func(r *RoleAlias) proto.Message { return RoleAliasToProto(r) },
		IDFunc:      func(r *RoleAlias) string { return r.RoleAlias },
		NotFoundErr: ErrRoleAliasNotFound, AlreadyExist: ErrRoleAliasAlreadyExists,
	})
	s.provisioningTplPS = common.NewProtoStore(common.ProtoStoreConfig[ProvisioningTemplate]{
		Store: s.templatesBase, NewProto: func() proto.Message { return &pb.ProvisioningTemplate{} },
		ToDomain: func(m proto.Message) *ProvisioningTemplate {
			return ProtoToProvisioningTemplate(m.(*pb.ProvisioningTemplate))
		},
		ToProto:     func(t *ProvisioningTemplate) proto.Message { return ProvisioningTemplateToProto(t) },
		IDFunc:      func(t *ProvisioningTemplate) string { return t.TemplateName },
		NotFoundErr: ErrTemplateNotFound, AlreadyExist: ErrTemplateAlreadyExists,
	})
	s.securityProfilePS = common.NewProtoStore(common.ProtoStoreConfig[SecurityProfile]{
		Store: s.securityProfilesBase, NewProto: func() proto.Message { return &pb.SecurityProfile{} },
		ToDomain: func(m proto.Message) *SecurityProfile {
			return ProtoToSecurityProfile(m.(*pb.SecurityProfile))
		},
		ToProto: func(sp *SecurityProfile) proto.Message {
			p, _ := SecurityProfileToProto(sp)
			return p
		},
		IDFunc:      func(sp *SecurityProfile) string { return sp.SecurityProfileName },
		NotFoundErr: ErrSecurityProfileNotFound, AlreadyExist: ErrSecurityProfileAlreadyExists,
	})
	s.domainConfigPS = common.NewProtoStore(common.ProtoStoreConfig[DomainConfiguration]{
		Store: s.domainConfigsBase, NewProto: func() proto.Message { return &pb.DomainConfiguration{} },
		ToDomain: func(m proto.Message) *DomainConfiguration {
			return ProtoToDomainConfiguration(m.(*pb.DomainConfiguration))
		},
		ToProto: func(dc *DomainConfiguration) proto.Message {
			p, _ := DomainConfigurationToProto(dc)
			return p
		},
		IDFunc:      func(dc *DomainConfiguration) string { return dc.DomainConfigurationName },
		NotFoundErr: ErrDomainConfigurationNotFound, AlreadyExist: ErrDomainConfigurationAlreadyExists,
	})
	s.violationEventPS = common.NewProtoStore(common.ProtoStoreConfig[ViolationEvent]{
		Store: s.violationEventsBase, NewProto: func() proto.Message { return &pb.ViolationEvent{} },
		ToDomain: func(m proto.Message) *ViolationEvent {
			return ProtoToViolationEvent(m.(*pb.ViolationEvent))
		},
		ToProto: func(v *ViolationEvent) proto.Message {
			p, _ := ViolationEventToProto(v)
			return p
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

