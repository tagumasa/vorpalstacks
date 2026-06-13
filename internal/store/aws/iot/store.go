package iot

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

type IotStore struct {
	thingsBase         *common.BaseStore
	thingTypesBase     *common.BaseStore
	thingGroupsBase    *common.BaseStore
	certsBase          *common.BaseStore
	policiesBase       *common.BaseStore
	rulesBase          *common.BaseStore
	jobsBase           *common.BaseStore
	shadowsBase        *common.BaseStore
	roleAliasBase      *common.BaseStore
	policyAttachBase   *common.BaseStore
	thingPrincipalBase *common.BaseStore
	principalThingBase *common.BaseStore
	policyPrincipalBase *common.BaseStore
	thingGroupMemberBase *common.BaseStore
	groupThingMemberBase *common.BaseStore
	billingGroupsBase  *common.BaseStore
	authorizersBase    *common.BaseStore
	templatesBase      *common.BaseStore
	detectorModelsBase *common.BaseStore
	inputsBase         *common.BaseStore
	*common.TagStore
	storage    storage.BasicStorage
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	mu         sync.RWMutex
}

func NewIotStore(store storage.BasicStorage, accountID, region string) *IotStore {
	regionSuffix := "-" + region
	return &IotStore{
		thingsBase:         common.NewBaseStore(store.Bucket("iot-things-"+regionSuffix), "iot-things"),
		thingTypesBase:     common.NewBaseStore(store.Bucket("iot-thing-types-"+regionSuffix), "iot-thing-types"),
		thingGroupsBase:    common.NewBaseStore(store.Bucket("iot-thing-groups-"+regionSuffix), "iot-thing-groups"),
		certsBase:          common.NewBaseStore(store.Bucket("iot-certificates-"+regionSuffix), "iot-certificates"),
		policiesBase:       common.NewBaseStore(store.Bucket("iot-policies-"+regionSuffix), "iot-policies"),
		rulesBase:          common.NewBaseStore(store.Bucket("iot-rules-"+regionSuffix), "iot-rules"),
		jobsBase:           common.NewBaseStore(store.Bucket("iot-jobs-"+regionSuffix), "iot-jobs"),
		shadowsBase:        common.NewBaseStore(store.Bucket("iot-shadows-"+regionSuffix), "iot-shadows"),
		roleAliasBase:      common.NewBaseStore(store.Bucket("iot-role-aliases-"+regionSuffix), "iot-role-aliases"),
		policyAttachBase:   common.NewBaseStore(store.Bucket("iot-policy-attach-"+regionSuffix), "iot-policy-attach"),
		thingPrincipalBase: common.NewBaseStore(store.Bucket("iot-thing-principal-"+regionSuffix), "iot-thing-principal"),
		principalThingBase: common.NewBaseStore(store.Bucket("iot-principal-thing-"+regionSuffix), "iot-principal-thing"),
		policyPrincipalBase: common.NewBaseStore(store.Bucket("iot-policy-principal-"+regionSuffix), "iot-policy-principal"),
		thingGroupMemberBase: common.NewBaseStore(store.Bucket("iot-thing-group-member-"+regionSuffix), "iot-thing-group-member"),
		groupThingMemberBase: common.NewBaseStore(store.Bucket("iot-group-thing-member-"+regionSuffix), "iot-group-thing-member"),
		billingGroupsBase:  common.NewBaseStore(store.Bucket("iot-billing-groups-"+regionSuffix), "iot-billing-groups"),
		authorizersBase:    common.NewBaseStore(store.Bucket("iot-authorizers-"+regionSuffix), "iot-authorizers"),
		templatesBase:      common.NewBaseStore(store.Bucket("iot-provisioning-templates-"+regionSuffix), "iot-provisioning-templates"),
		detectorModelsBase: common.NewBaseStore(store.Bucket("iot-detector-models-"+regionSuffix), "iot-detector-models"),
		inputsBase:         common.NewBaseStore(store.Bucket("iot-inputs-"+regionSuffix), "iot-inputs"),
		TagStore:           common.NewTagStoreWithRegion(store, "iot", region),
		storage:            store,
		arnBuilder:         svcarn.NewARNBuilder(accountID, region),
		accountID:          accountID,
		region:             region,
	}
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

func (s *IotStore) CreateThing(thing *Thing) (*Thing, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if thing.ThingName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Thing{}
	if err := s.thingsBase.GetProto(thing.ThingName, existing); err == nil {
		return nil, ErrThingAlreadyExists
	}
	if thing.ThingID == "" {
		thing.ThingID = uuid.New().String()
	}
	thing.ThingARN = BuildThingARN(s.accountID, s.region, thing.ThingName)
	return thing, s.thingsBase.PutProto(thing.ThingName, ThingToProto(thing))
}

func (s *IotStore) GetThing(thingName string) (*Thing, error) {
	pb := &pb.Thing{}
	if err := s.thingsBase.GetProto(thingName, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrThingNotFound
		}
		return nil, err
	}
	return ProtoToThing(pb), nil
}

func (s *IotStore) UpdateThing(thing *Thing) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.thingsBase.PutProto(thing.ThingName, ThingToProto(thing))
}

// errFound is a sentinel error used to short-circuit ScanPrefix
// once the first matching entry has been found.
var errFound = errors.New("found")

func (s *IotStore) DeleteThing(thingName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hasAttachments := false
	s.thingPrincipalBase.ScanPrefix(thingName+"\x00", func(key string, _ []byte) error {
		hasAttachments = true
		return errFound
	})
	if hasAttachments {
		return ErrDeleteConflict
	}
	// Remove ThingGroup memberships to prevent dangling references in
	// thingGroupMemberBase and groupThingMemberBase.
	var thingGroupKeys, groupThingKeys []string
	s.thingGroupMemberBase.ScanPrefix(thingName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) >= 2 {
			groupName := parts[1]
			groupThingKeys = append(groupThingKeys, groupName+"\x00"+thingName)
		}
		thingGroupKeys = append(thingGroupKeys, key)
		return nil
	})
	for _, k := range groupThingKeys {
		_ = s.groupThingMemberBase.Delete(k)
	}
	for _, k := range thingGroupKeys {
		_ = s.thingGroupMemberBase.Delete(k)
	}
	if err := s.shadowsBase.DeleteByPrefix(thingName + "/"); err != nil {
		return err
	}
	return s.thingsBase.Delete(thingName)
}

func (s *IotStore) ListThings(opts common.ListOptions) (*common.ListResult[Thing], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.thingsBase, opts, func() *pb.Thing { return &pb.Thing{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*Thing, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToThing(p))
	}
	return &common.ListResult[Thing]{
		Items:      items,
		NextMarker: result.NextMarker,
	}, nil
}

func (s *IotStore) ListThingsForThingType(thingTypeName string, opts common.ListOptions) (*common.ListResult[Thing], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	filter := func(p *pb.Thing) bool { return p.ThingTypeName == thingTypeName }
	result, err := common.ListProto(s.thingsBase, opts, func() *pb.Thing { return &pb.Thing{} }, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*Thing, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToThing(p))
	}
	return &common.ListResult[Thing]{
		Items:      items,
		NextMarker: result.NextMarker,
	}, nil
}

func (s *IotStore) CreateThingType(tt *ThingType) (*ThingType, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if tt.ThingTypeName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.ThingType{}
	if err := s.thingTypesBase.GetProto(tt.ThingTypeName, existing); err == nil {
		return nil, ErrThingTypeAlreadyExists
	}
	if tt.ThingTypeID == "" {
		tt.ThingTypeID = uuid.New().String()
	}
	tt.ThingTypeARN = BuildThingTypeARN(s.accountID, s.region, tt.ThingTypeName)
	return tt, s.thingTypesBase.PutProto(tt.ThingTypeName, ThingTypeToProto(tt))
}

func (s *IotStore) GetThingType(name string) (*ThingType, error) {
	pb := &pb.ThingType{}
	if err := s.thingTypesBase.GetProto(name, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrThingTypeNotFound
		}
		return nil, err
	}
	return ProtoToThingType(pb), nil
}

func (s *IotStore) DeleteThingType(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.thingTypesBase.Delete(name)
}

func (s *IotStore) ListThingTypes(opts common.ListOptions) (*common.ListResult[ThingType], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.thingTypesBase, opts, func() *pb.ThingType { return &pb.ThingType{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*ThingType, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToThingType(p))
	}
	return &common.ListResult[ThingType]{
		Items:      items,
		NextMarker: result.NextMarker,
	}, nil
}

func (s *IotStore) CreateThingGroup(group *ThingGroup) (*ThingGroup, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if group.GroupName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.ThingGroup{}
	if err := s.thingGroupsBase.GetProto(group.GroupName, existing); err == nil {
		return nil, ErrThingGroupAlreadyExists
	}
	if group.GroupID == "" {
		group.GroupID = uuid.New().String()
	}
	group.GroupARN = BuildThingGroupARN(s.accountID, s.region, group.GroupName)
	if group.Version == 0 {
		group.Version = 1
	}
	return group, s.thingGroupsBase.PutProto(group.GroupName, ThingGroupToProto(group))
}

func (s *IotStore) GetThingGroup(name string) (*ThingGroup, error) {
	pb := &pb.ThingGroup{}
	if err := s.thingGroupsBase.GetProto(name, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrThingGroupNotFound
		}
		return nil, err
	}
	return ProtoToThingGroup(pb), nil
}

func (s *IotStore) DeleteThingGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.thingGroupsBase.Delete(name)
}

func (s *IotStore) ListThingGroups(opts common.ListOptions) (*common.ListResult[ThingGroup], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.thingGroupsBase, opts, func() *pb.ThingGroup { return &pb.ThingGroup{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*ThingGroup, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToThingGroup(p))
	}
	return &common.ListResult[ThingGroup]{
		Items:      items,
		NextMarker: result.NextMarker,
	}, nil
}

func (s *IotStore) CreateCertificate(cert *Certificate) (*Certificate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cert.CertificateID == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Certificate{}
	if err := s.certsBase.GetProto(cert.CertificateID, existing); err == nil {
		return nil, ErrCertificateAlreadyExists
	}
	cert.CertificateARN = BuildCertificateARN(s.accountID, s.region, cert.CertificateID)
	return cert, s.certsBase.PutProto(cert.CertificateID, CertificateToProto(cert))
}

func (s *IotStore) GetCertificate(certificateID string) (*Certificate, error) {
	pb := &pb.Certificate{}
	if err := s.certsBase.GetProto(certificateID, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrCertificateNotFound
		}
		return nil, err
	}
	return ProtoToCertificate(pb), nil
}

func (s *IotStore) UpdateCertificate(cert *Certificate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.certsBase.PutProto(cert.CertificateID, CertificateToProto(cert))
}

func (s *IotStore) DeleteCertificate(certificateID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.certsBase.Delete(certificateID)
}

func (s *IotStore) ListCertificates(opts common.ListOptions) (*common.ListResult[Certificate], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.certsBase, opts, func() *pb.Certificate { return &pb.Certificate{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*Certificate, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToCertificate(p))
	}
	return &common.ListResult[Certificate]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) CreatePolicy(policy *Policy) (*Policy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if policy.PolicyName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Policy{}
	if err := s.policiesBase.GetProto(policy.PolicyName, existing); err == nil {
		return nil, ErrPolicyAlreadyExists
	}
	policy.PolicyARN = BuildPolicyARN(s.accountID, s.region, policy.PolicyName)
	return policy, s.policiesBase.PutProto(policy.PolicyName, PolicyToProto(policy))
}

func (s *IotStore) GetPolicy(policyName string) (*Policy, error) {
	pb := &pb.Policy{}
	if err := s.policiesBase.GetProto(policyName, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrPolicyNotFound
		}
		return nil, err
	}
	return ProtoToPolicy(pb), nil
}

func (s *IotStore) DeletePolicy(policyName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	principals, err := s.listPrincipalsForPolicyLocked(policyName)
	if err != nil {
		return err
	}
	if len(principals) > 0 {
		return ErrDeleteConflict
	}
	return s.policiesBase.Delete(policyName)
}

func (s *IotStore) ListPolicies(opts common.ListOptions) (*common.ListResult[Policy], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.policiesBase, opts, func() *pb.Policy { return &pb.Policy{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*Policy, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToPolicy(p))
	}
	return &common.ListResult[Policy]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) CreateRule(rule *TopicRule) (*TopicRule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if rule.RuleName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.TopicRule{}
	if err := s.rulesBase.GetProto(rule.RuleName, existing); err == nil {
		return nil, ErrRuleAlreadyExists
	}
	rule.ARN = BuildRuleARN(s.accountID, s.region, rule.RuleName)
	return rule, s.rulesBase.PutProto(rule.RuleName, RuleToProto(rule))
}

func (s *IotStore) GetRule(ruleName string) (*TopicRule, error) {
	pb := &pb.TopicRule{}
	if err := s.rulesBase.GetProto(ruleName, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrRuleNotFound
		}
		return nil, err
	}
	return ProtoToRule(pb), nil
}

func (s *IotStore) DeleteRule(ruleName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rulesBase.Delete(ruleName)
}

func (s *IotStore) ListRules(opts common.ListOptions) (*common.ListResult[TopicRule], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.rulesBase, opts, func() *pb.TopicRule { return &pb.TopicRule{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*TopicRule, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToRule(p))
	}
	return &common.ListResult[TopicRule]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) CreateJob(job *Job) (*Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job.JobID == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Job{}
	if err := s.jobsBase.GetProto(job.JobID, existing); err == nil {
		return nil, ErrJobAlreadyExists
	}
	job.JobARN = BuildJobARN(s.accountID, s.region, job.JobID)
	return job, s.jobsBase.PutProto(job.JobID, JobToProto(job))
}

func (s *IotStore) GetJob(jobID string) (*Job, error) {
	pb := &pb.Job{}
	if err := s.jobsBase.GetProto(jobID, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}
	return ProtoToJob(pb), nil
}

func (s *IotStore) DeleteJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.jobsBase.Delete(jobID)
}

func (s *IotStore) ListJobs(opts common.ListOptions) (*common.ListResult[Job], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.jobsBase, opts, func() *pb.Job { return &pb.Job{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*Job, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToJob(p))
	}
	return &common.ListResult[Job]{Items: items, NextMarker: result.NextMarker}, nil
}

func shadowKey(thingName, shadowName string) string {
	if shadowName == "" || shadowName == "classic" {
		return thingName + "/$current"
	}
	return thingName + "/" + shadowName
}

func (s *IotStore) GetShadow(thingName, shadowName string) (*ShadowDocument, error) {
	pb := &pb.ShadowDocument{}
	if err := s.shadowsBase.GetProto(shadowKey(thingName, shadowName), pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrShadowNotFound
		}
		return nil, err
	}
	return ProtoToShadow(pb), nil
}

func (s *IotStore) PutShadow(thingName, shadowName string, doc *ShadowDocument) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.shadowsBase.PutProto(shadowKey(thingName, shadowName), ShadowToProto(doc))
}

// PutShadowWithVersion atomically reads the current shadow, checks the version,
// and writes the new document if the version matches.  Returns ErrVersionConflict
// if clientVersion > 0 and does not match the stored version.
func (s *IotStore) PutShadowWithVersion(thingName, shadowName string, doc *ShadowDocument, clientVersion int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pb := &pb.ShadowDocument{}
	err := s.shadowsBase.GetProto(shadowKey(thingName, shadowName), pb)
	if err == nil {
		stored := ProtoToShadow(pb)
		if clientVersion > 0 && stored.VersionNumber != clientVersion {
			return ErrVersionConflict
		}
	} else if !common.IsNotFound(err) {
		return err
	}
	return s.shadowsBase.PutProto(shadowKey(thingName, shadowName), ShadowToProto(doc))
}

func (s *IotStore) DeleteShadow(thingName, shadowName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.shadowsBase.Delete(shadowKey(thingName, shadowName))
}

// ListShadowNames returns the distinct shadow names for a given thing by
// scanning the shadows bucket and extracting the shadow name component.
func (s *IotStore) ListShadowNames(thingName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	prefix := thingName + "/"
	var names []string
	seen := make(map[string]bool)
	err := s.shadowsBase.ScanPrefix(prefix, func(key string, _ []byte) error {
		parts := strings.SplitN(key[len(prefix):], "/", 2)
		name := parts[0]
		if name == "$current" {
			return nil
		}
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
		return nil
	})
	return names, err
}

func (s *IotStore) ListTags(resourceARN string) (map[string]string, error) {
	return s.TagStore.List(resourceARN)
}

func (s *IotStore) TagResource(resourceARN string, tags map[string]string) error {
	return s.TagStore.Tag(resourceARN, tags)
}

func (s *IotStore) UntagResource(resourceARN string, tagKeys []string) error {
	return s.TagStore.Untag(resourceARN, tagKeys)
}

func (s *IotStore) UpdateRule(rule *TopicRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rulesBase.PutProto(rule.RuleName, RuleToProto(rule))
}

// UpdateThingType persists changes to an existing thing type.
func (s *IotStore) UpdateThingType(tt *ThingType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.thingTypesBase.PutProto(tt.ThingTypeName, ThingTypeToProto(tt))
}

// UpdateThingGroup persists changes to an existing thing group.
func (s *IotStore) UpdateThingGroup(group *ThingGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	group.Version++
	return s.thingGroupsBase.PutProto(group.GroupName, ThingGroupToProto(group))
}

// CreateBillingGroup persists a new billing group.
func (s *IotStore) CreateBillingGroup(bg *BillingGroup) (*BillingGroup, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if bg.GroupName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.BillingGroup{}
	if err := s.billingGroupsBase.GetProto(bg.GroupName, existing); err == nil {
		return nil, ErrBillingGroupAlreadyExists
	}
	if bg.GroupID == "" {
		bg.GroupID = uuid.New().String()
	}
	bg.GroupARN = BuildBillingGroupARN(s.accountID, s.region, bg.GroupName)
	return bg, s.billingGroupsBase.PutProto(bg.GroupName, BillingGroupToProto(bg))
}

// GetBillingGroup retrieves a billing group by name.
func (s *IotStore) GetBillingGroup(name string) (*BillingGroup, error) {
	pb := &pb.BillingGroup{}
	if err := s.billingGroupsBase.GetProto(name, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrBillingGroupNotFound
		}
		return nil, err
	}
	return ProtoToBillingGroup(pb), nil
}

// UpdateBillingGroup persists changes to an existing billing group.
func (s *IotStore) UpdateBillingGroup(bg *BillingGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.billingGroupsBase.PutProto(bg.GroupName, BillingGroupToProto(bg))
}

// DeleteBillingGroup removes a billing group by name.
func (s *IotStore) DeleteBillingGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.billingGroupsBase.Delete(name)
}

// ListBillingGroups returns all billing groups.
func (s *IotStore) ListBillingGroups(opts common.ListOptions) (*common.ListResult[BillingGroup], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.billingGroupsBase, opts, func() *pb.BillingGroup { return &pb.BillingGroup{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*BillingGroup, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToBillingGroup(p))
	}
	return &common.ListResult[BillingGroup]{
		Items:      items,
		NextMarker: result.NextMarker,
	}, nil
}

// CreateAuthorizer persists a new custom authorizer.
func (s *IotStore) CreateAuthorizer(a *Authorizer) (*Authorizer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a.AuthorizerName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Authorizer{}
	if err := s.authorizersBase.GetProto(a.AuthorizerName, existing); err == nil {
		return nil, ErrAuthorizerAlreadyExists
	}
	a.AuthorizerARN = BuildAuthorizerARN(s.accountID, s.region, a.AuthorizerName)
	return a, s.authorizersBase.PutProto(a.AuthorizerName, AuthorizerToProto(a))
}

// GetAuthorizer retrieves an authorizer by name.
func (s *IotStore) GetAuthorizer(name string) (*Authorizer, error) {
	pb := &pb.Authorizer{}
	if err := s.authorizersBase.GetProto(name, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrAuthorizerNotFound
		}
		return nil, err
	}
	return ProtoToAuthorizer(pb), nil
}

// UpdateAuthorizer persists changes to an existing authorizer.
func (s *IotStore) UpdateAuthorizer(a *Authorizer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.authorizersBase.PutProto(a.AuthorizerName, AuthorizerToProto(a))
}

// DeleteAuthorizer removes an authorizer by name.
func (s *IotStore) DeleteAuthorizer(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.authorizersBase.Delete(name)
}

// ListAuthorizers returns all custom authorizers.
func (s *IotStore) ListAuthorizers(opts common.ListOptions) (*common.ListResult[Authorizer], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.authorizersBase, opts, func() *pb.Authorizer { return &pb.Authorizer{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*Authorizer, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToAuthorizer(p))
	}
	return &common.ListResult[Authorizer]{Items: items, NextMarker: result.NextMarker}, nil
}

// CreateProvisioningTemplate persists a new provisioning template.
func (s *IotStore) CreateProvisioningTemplate(t *ProvisioningTemplate) (*ProvisioningTemplate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t.TemplateName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.ProvisioningTemplate{}
	if err := s.templatesBase.GetProto(t.TemplateName, existing); err == nil {
		return nil, ErrTemplateAlreadyExists
	}
	t.TemplateARN = BuildProvisioningTemplateARN(s.accountID, s.region, t.TemplateName)
	return t, s.templatesBase.PutProto(t.TemplateName, ProvisioningTemplateToProto(t))
}

// GetProvisioningTemplate retrieves a provisioning template by name.
func (s *IotStore) GetProvisioningTemplate(name string) (*ProvisioningTemplate, error) {
	pb := &pb.ProvisioningTemplate{}
	if err := s.templatesBase.GetProto(name, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	return ProtoToProvisioningTemplate(pb), nil
}

// UpdateProvisioningTemplate persists changes to an existing provisioning template.
func (s *IotStore) UpdateProvisioningTemplate(t *ProvisioningTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.templatesBase.PutProto(t.TemplateName, ProvisioningTemplateToProto(t))
}

// DeleteProvisioningTemplate removes a provisioning template by name.
func (s *IotStore) DeleteProvisioningTemplate(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.templatesBase.Delete(name)
}

// ListProvisioningTemplates returns all provisioning templates.
func (s *IotStore) ListProvisioningTemplates(opts common.ListOptions) (*common.ListResult[ProvisioningTemplate], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.templatesBase, opts, func() *pb.ProvisioningTemplate { return &pb.ProvisioningTemplate{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*ProvisioningTemplate, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToProvisioningTemplate(p))
	}
	return &common.ListResult[ProvisioningTemplate]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) CreateDetectorModel(d *DetectorModel) (*DetectorModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.DetectorModelName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.DetectorModel{}
	if err := s.detectorModelsBase.GetProto(d.DetectorModelName, existing); err == nil {
		return nil, ErrDetectorModelAlreadyExists
	}
	d.DetectorModelARN = BuildDetectorModelARN(s.accountID, s.region, d.DetectorModelName)
	pb, err := DetectorModelToProto(d)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise detector model: %w", err)
	}
	return d, s.detectorModelsBase.PutProto(d.DetectorModelName, pb)
}

func (s *IotStore) GetDetectorModel(name string) (*DetectorModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := &pb.DetectorModel{}
	if err := s.detectorModelsBase.GetProto(name, p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrDetectorModelNotFound
		}
		return nil, err
	}
	return ProtoToDetectorModel(p), nil
}

func (s *IotStore) UpdateDetectorModel(d *DetectorModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	pb, err := DetectorModelToProto(d)
	if err != nil {
		return fmt.Errorf("failed to serialise detector model: %w", err)
	}
	return s.detectorModelsBase.PutProto(d.DetectorModelName, pb)
}

func (s *IotStore) DeleteDetectorModel(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.detectorModelsBase.Delete(name)
}

func (s *IotStore) ListDetectorModels(opts common.ListOptions) (*common.ListResult[DetectorModel], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.detectorModelsBase, opts, func() *pb.DetectorModel { return &pb.DetectorModel{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*DetectorModel, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToDetectorModel(p))
	}
	return &common.ListResult[DetectorModel]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) CreateInput(i *Input) (*Input, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i.InputName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Input{}
	if err := s.inputsBase.GetProto(i.InputName, existing); err == nil {
		return nil, ErrInputAlreadyExists
	}
	i.InputARN = BuildInputARN(s.accountID, s.region, i.InputName)
	pb, err := InputToProto(i)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise input: %w", err)
	}
	return i, s.inputsBase.PutProto(i.InputName, pb)
}

func (s *IotStore) GetInput(name string) (*Input, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := &pb.Input{}
	if err := s.inputsBase.GetProto(name, p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrInputNotFound
		}
		return nil, err
	}
	return ProtoToInput(p), nil
}

func (s *IotStore) UpdateInput(i *Input) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	pb, err := InputToProto(i)
	if err != nil {
		return fmt.Errorf("failed to serialise input: %w", err)
	}
	return s.inputsBase.PutProto(i.InputName, pb)
}

func (s *IotStore) DeleteInput(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.inputsBase.Delete(name)
}

func (s *IotStore) ListInputs(opts common.ListOptions) (*common.ListResult[Input], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.inputsBase, opts, func() *pb.Input { return &pb.Input{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*Input, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToInput(p))
	}
	return &common.ListResult[Input]{Items: items, NextMarker: result.NextMarker}, nil
}

// UpdateJob persists changes to an existing job.
func (s *IotStore) UpdateJob(job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.jobsBase.PutProto(job.JobID, JobToProto(job))
}

func (s *IotStore) CreateRoleAlias(ra *RoleAlias) (*RoleAlias, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ra.RoleAlias == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.RoleAlias{}
	if err := s.roleAliasBase.GetProto(ra.RoleAlias, existing); err == nil {
		return nil, ErrRoleAliasAlreadyExists
	}
	ra.RoleAliasARN = BuildRoleAliasARN(s.accountID, s.region, ra.RoleAlias)
	return ra, s.roleAliasBase.PutProto(ra.RoleAlias, RoleAliasToProto(ra))
}

func (s *IotStore) GetRoleAlias(alias string) (*RoleAlias, error) {
	pb := &pb.RoleAlias{}
	if err := s.roleAliasBase.GetProto(alias, pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrRoleAliasNotFound
		}
		return nil, err
	}
	return ProtoToRoleAlias(pb), nil
}

// UpdateRoleAlias persists changes to an existing role alias.
func (s *IotStore) UpdateRoleAlias(ra *RoleAlias) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.roleAliasBase.PutProto(ra.RoleAlias, RoleAliasToProto(ra))
}

func (s *IotStore) DeleteRoleAlias(alias string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.roleAliasBase.Delete(alias)
}

func (s *IotStore) ListRoleAliases(opts common.ListOptions) (*common.ListResult[RoleAlias], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.roleAliasBase, opts, func() *pb.RoleAlias { return &pb.RoleAlias{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*RoleAlias, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToRoleAlias(p))
	}
	return &common.ListResult[RoleAlias]{Items: items, NextMarker: result.NextMarker}, nil
}

func policyAttachKey(policyName, principal string) string {
	return principal + "\x00" + policyName
}

func policyPrincipalKey(policyName, principal string) string {
	return policyName + "\x00" + principal
}

func (s *IotStore) AttachPolicyToPrincipal(policyName, principal string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ak := policyAttachKey(policyName, principal)
	pk := policyPrincipalKey(policyName, principal)
	if err := s.policyAttachBase.Put(ak, []byte("1")); err != nil {
		return err
	}
	if err := s.policyPrincipalBase.Put(pk, []byte("1")); err != nil {
		// Compensate: remove the already-written attach entry to avoid dangling reference.
		_ = s.policyAttachBase.Delete(ak)
		return err
	}
	return nil
}

func (s *IotStore) DetachPolicyFromPrincipal(policyName, principal string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ak := policyAttachKey(policyName, principal)
	pk := policyPrincipalKey(policyName, principal)
	if err := s.policyAttachBase.Delete(ak); err != nil {
		return err
	}
	if err := s.policyPrincipalBase.Delete(pk); err != nil {
		// Compensate: re-create the already-deleted attach entry.
		_ = s.policyAttachBase.Put(ak, []byte("1"))
		return err
	}
	return nil
}

// listPrincipalsForPolicyLocked returns principals attached to a policy.
// Caller must hold s.mu.
func (s *IotStore) listPrincipalsForPolicyLocked(policyName string) ([]string, error) {
	var principals []string
	err := s.policyPrincipalBase.ScanPrefix(policyName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			principals = append(principals, parts[1])
		}
		return nil
	})
	return principals, err
}

func (s *IotStore) ListPrincipalsForPolicy(policyName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.listPrincipalsForPolicyLocked(policyName)
}

func (s *IotStore) ListPoliciesForPrincipal(principal string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var policies []string
	err := s.policyAttachBase.ScanPrefix(principal+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			policies = append(policies, parts[1])
		}
		return nil
	})
	return policies, err
}

func thingPrincipalKey(thingName, principal string) string {
	return thingName + "\x00" + principal
}

func principalThingKey(principal, thingName string) string {
	return principal + "\x00" + thingName
}

func (s *IotStore) AttachThingPrincipal(thingName, principal string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tk := thingPrincipalKey(thingName, principal)
	pk := principalThingKey(principal, thingName)
	if err := s.thingPrincipalBase.Put(tk, []byte("1")); err != nil {
		return err
	}
	if err := s.principalThingBase.Put(pk, []byte("1")); err != nil {
		// Compensate: remove the already-written thing entry to avoid dangling reference.
		_ = s.thingPrincipalBase.Delete(tk)
		return err
	}
	return nil
}

func (s *IotStore) DetachThingPrincipal(thingName, principal string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tk := thingPrincipalKey(thingName, principal)
	pk := principalThingKey(principal, thingName)
	if err := s.thingPrincipalBase.Delete(tk); err != nil {
		return err
	}
	if err := s.principalThingBase.Delete(pk); err != nil {
		// Compensate: re-create the already-deleted thing entry.
		_ = s.thingPrincipalBase.Put(tk, []byte("1"))
		return err
	}
	return nil
}

func (s *IotStore) ListPrincipalsForThing(thingName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var principals []string
	err := s.thingPrincipalBase.ScanPrefix(thingName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			principals = append(principals, parts[1])
		}
		return nil
	})
	return principals, err
}

func (s *IotStore) ListThingsForPrincipal(principal string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var things []string
	err := s.principalThingBase.ScanPrefix(principal+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			things = append(things, parts[1])
		}
		return nil
	})
	return things, err
}

func (s *IotStore) AddThingToThingGroup(thingName, groupName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tgk := thingName + "\x00" + groupName
	gtk := groupName + "\x00" + thingName
	if err := s.thingGroupMemberBase.Put(tgk, []byte("1")); err != nil {
		return err
	}
	if err := s.groupThingMemberBase.Put(gtk, []byte("1")); err != nil {
		_ = s.thingGroupMemberBase.Delete(tgk)
		return err
	}
	return nil
}

func (s *IotStore) RemoveThingFromThingGroup(thingName, groupName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tgk := thingName + "\x00" + groupName
	gtk := groupName + "\x00" + thingName
	if err := s.thingGroupMemberBase.Delete(tgk); err != nil {
		return err
	}
	if err := s.groupThingMemberBase.Delete(gtk); err != nil {
		_ = s.thingGroupMemberBase.Put(tgk, []byte("1"))
		return err
	}
	return nil
}

func (s *IotStore) ListThingsInGroup(groupName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var things []string
	err := s.groupThingMemberBase.ScanPrefix(groupName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			things = append(things, parts[1])
		}
		return nil
	})
	return things, err
}

func (s *IotStore) ListGroupsForThing(thingName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var groups []string
	err := s.thingGroupMemberBase.ScanPrefix(thingName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			groups = append(groups, parts[1])
		}
		return nil
	})
	return groups, err
}

var _ IotStoreInterface = (*IotStore)(nil)
