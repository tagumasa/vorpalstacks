package iot

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

type Thing struct {
	ThingName        string
	ThingARN         string
	ThingID          string
	ThingTypeName    string
	AttributeNames   []string
	Attributes       map[string]string
	Version          int64
	CreationDate     time.Time
	LastModifiedDate time.Time
}

func ThingToProto(t *Thing) *pb.Thing {
	return &pb.Thing{
		ThingName:        t.ThingName,
		ThingArn:         t.ThingARN,
		ThingId:          t.ThingID,
		ThingTypeName:    t.ThingTypeName,
		AttributeNames:   t.AttributeNames,
		Attributes:       t.Attributes,
		Version:          t.Version,
		CreationDate:     timeToProto(t.CreationDate),
		LastModifiedDate: timeToProto(t.LastModifiedDate),
	}
}

func ProtoToThing(p *pb.Thing) *Thing {
	return &Thing{
		ThingName:        p.ThingName,
		ThingARN:         p.ThingArn,
		ThingID:          p.ThingId,
		ThingTypeName:    p.ThingTypeName,
		AttributeNames:   p.AttributeNames,
		Attributes:       p.Attributes,
		Version:          p.Version,
		CreationDate:     protoToTime(p.CreationDate),
		LastModifiedDate: protoToTime(p.LastModifiedDate),
	}
}

type ThingType struct {
	ThingTypeName       string
	ThingTypeARN        string
	ThingTypeID         string
	Description         string
	ThingTypeProperties []ThingTypeProperty
	Tags                map[string]string
	Version             int64
	CreationDate        time.Time
	LastModifiedDate    time.Time
}

type ThingTypeProperty struct {
	Name        string
	DataType    string
	Description string
	Required    bool
	FieldName   string
}

func ThingTypeToProto(t *ThingType) *pb.ThingType {
	props := make([]*pb.ThingTypeProperty, len(t.ThingTypeProperties))
	for i, p := range t.ThingTypeProperties {
		props[i] = &pb.ThingTypeProperty{
			Name: p.Name, DataType: p.DataType, Description: p.Description,
			Required: p.Required, FieldName: p.FieldName,
		}
	}
	return &pb.ThingType{
		ThingTypeName:       t.ThingTypeName,
		ThingTypeArn:        t.ThingTypeARN,
		ThingTypeId:         t.ThingTypeID,
		Description:         t.Description,
		ThingTypeProperties: props,
		Tags:                t.Tags,
		Version:             t.Version,
		CreationDate:        timeToProto(t.CreationDate),
		LastModifiedDate:    timeToProto(t.LastModifiedDate),
	}
}

func ProtoToThingType(p *pb.ThingType) *ThingType {
	props := make([]ThingTypeProperty, len(p.ThingTypeProperties))
	for i, pp := range p.ThingTypeProperties {
		props[i] = ThingTypeProperty{
			Name: pp.Name, DataType: pp.DataType, Description: pp.Description,
			Required: pp.Required, FieldName: pp.FieldName,
		}
	}
	return &ThingType{
		ThingTypeName:       p.ThingTypeName,
		ThingTypeARN:        p.ThingTypeArn,
		ThingTypeID:         p.ThingTypeId,
		Description:         p.Description,
		ThingTypeProperties: props,
		Tags:                p.Tags,
		Version:             p.Version,
		CreationDate:        protoToTime(p.CreationDate),
		LastModifiedDate:    protoToTime(p.LastModifiedDate),
	}
}

type ThingGroup struct {
	GroupName        string
	GroupARN         string
	GroupID          string
	ParentGroupName  string
	Description      string
	Attributes       map[string]string
	Tags             map[string]string
	CreationDate     time.Time
	LastModifiedDate time.Time
	Version          int64
}

func ThingGroupToProto(g *ThingGroup) *pb.ThingGroup {
	return &pb.ThingGroup{
		GroupName: g.GroupName, GroupArn: g.GroupARN, GroupId: g.GroupID,
		ParentGroupName: g.ParentGroupName, Description: g.Description,
		Attributes: g.Attributes, Tags: g.Tags,
		CreationDate: timeToProto(g.CreationDate), LastModifiedDate: timeToProto(g.LastModifiedDate),
		Version: g.Version,
	}
}

func ProtoToThingGroup(p *pb.ThingGroup) *ThingGroup {
	return &ThingGroup{
		GroupName: p.GroupName, GroupARN: p.GroupArn, GroupID: p.GroupId,
		ParentGroupName: p.ParentGroupName, Description: p.Description,
		Attributes: p.Attributes, Tags: p.Tags,
		CreationDate: protoToTime(p.CreationDate), LastModifiedDate: protoToTime(p.LastModifiedDate),
		Version: p.Version,
	}
}

type BillingGroup struct {
	GroupName        string
	GroupARN         string
	GroupID          string
	Description      string
	Attributes       map[string]string
	Tags             map[string]string
	CreationDate     time.Time
	LastModifiedDate time.Time
}

func BillingGroupToProto(g *BillingGroup) *pb.BillingGroup {
	return &pb.BillingGroup{
		GroupName: g.GroupName, GroupArn: g.GroupARN, GroupId: g.GroupID,
		Description: g.Description, Attributes: g.Attributes, Tags: g.Tags,
		CreationDate: timeToProto(g.CreationDate), LastModifiedDate: timeToProto(g.LastModifiedDate),
	}
}

func ProtoToBillingGroup(p *pb.BillingGroup) *BillingGroup {
	return &BillingGroup{
		GroupName: p.GroupName, GroupARN: p.GroupArn, GroupID: p.GroupId,
		Description: p.Description, Attributes: p.Attributes, Tags: p.Tags,
		CreationDate: protoToTime(p.CreationDate), LastModifiedDate: protoToTime(p.LastModifiedDate),
	}
}

type Certificate struct {
	CertificateID       string
	CertificateARN      string
	CertificatePEM      string
	CertificateMode     string
	CreationDate        time.Time
	LastModifiedDate    time.Time
	Status              string
	CaCertificateID     string
	PublicKey           []byte
	TransferDate        time.Time
	TransferToAccountID string
	AutoActive          bool
}

func CertificateToProto(c *Certificate) *pb.Certificate {
	return &pb.Certificate{
		CertificateId: c.CertificateID, CertificateArn: c.CertificateARN,
		CertificatePem: c.CertificatePEM, CertificateMode: c.CertificateMode,
		CreationDate: timeToProto(c.CreationDate), LastModifiedDate: timeToProto(c.LastModifiedDate),
		Status: c.Status, CaCertificateId: c.CaCertificateID, PublicKey: c.PublicKey,
		TransferDate: timeToProto(c.TransferDate), TransferToAccountId: c.TransferToAccountID,
		AutoActive: c.AutoActive,
	}
}

func ProtoToCertificate(p *pb.Certificate) *Certificate {
	return &Certificate{
		CertificateID: p.CertificateId, CertificateARN: p.CertificateArn,
		CertificatePEM: p.CertificatePem, CertificateMode: p.CertificateMode,
		CreationDate: protoToTime(p.CreationDate), LastModifiedDate: protoToTime(p.LastModifiedDate),
		Status: p.Status, CaCertificateID: p.CaCertificateId, PublicKey: p.PublicKey,
		TransferDate: protoToTime(p.TransferDate), TransferToAccountID: p.TransferToAccountId,
		AutoActive: p.AutoActive,
	}
}

type Policy struct {
	PolicyName       string
	PolicyARN        string
	PolicyDocument   string
	CreationDate     time.Time
	LastModifiedDate time.Time
	Version          int64
}

func PolicyToProto(p *Policy) *pb.Policy {
	return &pb.Policy{
		PolicyName: p.PolicyName, PolicyArn: p.PolicyARN,
		PolicyDocument: p.PolicyDocument,
		CreationDate:   timeToProto(p.CreationDate), LastModifiedDate: timeToProto(p.LastModifiedDate),
		Version: p.Version,
	}
}

func ProtoToPolicy(p *pb.Policy) *Policy {
	return &Policy{
		PolicyName: p.PolicyName, PolicyARN: p.PolicyArn,
		PolicyDocument: p.PolicyDocument,
		CreationDate:   protoToTime(p.CreationDate), LastModifiedDate: protoToTime(p.LastModifiedDate),
		Version: p.Version,
	}
}

type TopicRule struct {
	RuleName         string
	ARN              string
	TopicPattern     string
	Description      string
	RuleDisabled     bool
	SQL              string
	CreatedAt        string
	AwsIotSqlVersion string
}

func RuleToProto(r *TopicRule) *pb.TopicRule {
	return &pb.TopicRule{
		RuleName: r.RuleName, Arn: r.ARN, TopicPattern: r.TopicPattern,
		Description: r.Description, RuleDisabled: r.RuleDisabled,
		Sql: r.SQL, CreatedAt: r.CreatedAt, AwsIotSqlVersion: r.AwsIotSqlVersion,
	}
}

func ProtoToRule(p *pb.TopicRule) *TopicRule {
	return &TopicRule{
		RuleName: p.RuleName, ARN: p.Arn, TopicPattern: p.TopicPattern,
		Description: p.Description, RuleDisabled: p.RuleDisabled,
		SQL: p.Sql, CreatedAt: p.CreatedAt, AwsIotSqlVersion: p.AwsIotSqlVersion,
	}
}

type Job struct {
	JobARN          string
	JobID           string
	Description     string
	Version         int64
	Force           bool
	CreatedAt       time.Time
	LastUpdatedAt   time.Time
	CompletedAt     string
	JobTemplateARN  string
	Status          string
	TargetSelection string
	Tags            map[string]string
}

func JobToProto(j *Job) *pb.Job {
	return &pb.Job{
		JobArn: j.JobARN, JobId: j.JobID, Description: j.Description,
		Version: j.Version, Force: j.Force,
		CreatedAt: timeToProto(j.CreatedAt), LastUpdatedAt: timeToProto(j.LastUpdatedAt),
		CompletedAt: j.CompletedAt, JobTemplateArn: j.JobTemplateARN,
		Status: j.Status, TargetSelection: j.TargetSelection, Tags: j.Tags,
	}
}

func ProtoToJob(p *pb.Job) *Job {
	return &Job{
		JobARN: p.JobArn, JobID: p.JobId, Description: p.Description,
		Version: p.Version, Force: p.Force,
		CreatedAt: protoToTime(p.CreatedAt), LastUpdatedAt: protoToTime(p.LastUpdatedAt),
		CompletedAt: p.CompletedAt, JobTemplateARN: p.JobTemplateArn,
		Status: p.Status, TargetSelection: p.TargetSelection, Tags: p.Tags,
	}
}

type ShadowDocument struct {
	ThingName     string
	Version       time.Time
	VersionNumber int64
	Timestamp     time.Time
	State         string
	Metadata      string
}

func ShadowToProto(s *ShadowDocument) *pb.ShadowDocument {
	return &pb.ShadowDocument{
		ThingName: s.ThingName, Version: timeToProto(s.Version),
		VersionNumber: s.VersionNumber, Timestamp: timeToProto(s.Timestamp),
		State: s.State, Metadata: s.Metadata,
	}
}

func ProtoToShadow(p *pb.ShadowDocument) *ShadowDocument {
	return &ShadowDocument{
		ThingName: p.ThingName, Version: protoToTime(p.Version),
		VersionNumber: p.VersionNumber, Timestamp: protoToTime(p.Timestamp),
		State: p.State, Metadata: p.Metadata,
	}
}

type Authorizer struct {
	AuthorizerName        string
	AuthorizerARN         string
	AuthorizerFunctionARN string
	TokenName             string
	TokenSignature        string
	Status                bool
	CreationDate          time.Time
	LastModifiedDate      time.Time
	EnableCachingForHTTP  bool
	CachingDisabled       int64
}

func AuthorizerToProto(a *Authorizer) *pb.Authorizer {
	return &pb.Authorizer{
		AuthorizerName: a.AuthorizerName, AuthorizerArn: a.AuthorizerARN,
		AuthorizerFunctionArn: a.AuthorizerFunctionARN,
		TokenKeyName:          a.TokenName, TokenSignature: a.TokenSignature, Status: a.Status,
		CreationDate: timeToProto(a.CreationDate), LastModifiedDate: timeToProto(a.LastModifiedDate),
		EnableCachingForHttp: a.EnableCachingForHTTP, CachingDisabled: a.CachingDisabled,
	}
}

func ProtoToAuthorizer(p *pb.Authorizer) *Authorizer {
	return &Authorizer{
		AuthorizerName: p.AuthorizerName, AuthorizerARN: p.AuthorizerArn,
		AuthorizerFunctionARN: p.AuthorizerFunctionArn,
		TokenName:             p.TokenKeyName, TokenSignature: p.TokenSignature, Status: p.Status,
		CreationDate: protoToTime(p.CreationDate), LastModifiedDate: protoToTime(p.LastModifiedDate),
		EnableCachingForHTTP: p.EnableCachingForHttp, CachingDisabled: p.CachingDisabled,
	}
}

type RoleAlias struct {
	RoleAlias                 string
	RoleAliasARN              string
	RoleARN                   string
	CredentialDurationSeconds string
	Owner                     string
	CreationDate              time.Time
	LastModifiedDate          time.Time
}

func RoleAliasToProto(r *RoleAlias) *pb.RoleAlias {
	return &pb.RoleAlias{
		RoleAlias: r.RoleAlias, RoleAliasArn: r.RoleAliasARN,
		RoleArn: r.RoleARN, CredentialDurationSeconds: r.CredentialDurationSeconds,
		Owner:        r.Owner,
		CreationDate: timeToProto(r.CreationDate), LastModifiedDate: timeToProto(r.LastModifiedDate),
	}
}

func ProtoToRoleAlias(p *pb.RoleAlias) *RoleAlias {
	return &RoleAlias{
		RoleAlias: p.RoleAlias, RoleAliasARN: p.RoleAliasArn,
		RoleARN: p.RoleArn, CredentialDurationSeconds: p.CredentialDurationSeconds,
		Owner:        p.Owner,
		CreationDate: protoToTime(p.CreationDate), LastModifiedDate: protoToTime(p.LastModifiedDate),
	}
}

type ProvisioningTemplate struct {
	TemplateName        string
	TemplateARN         string
	Description         string
	Enabled             bool
	ProvisioningRoleARN string
	Tags                map[string]string
	CreationDate        time.Time
	LastModifiedDate    time.Time
}

func ProvisioningTemplateToProto(t *ProvisioningTemplate) *pb.ProvisioningTemplate {
	return &pb.ProvisioningTemplate{
		TemplateName: t.TemplateName, TemplateArn: t.TemplateARN,
		Description: t.Description, Enabled: t.Enabled,
		ProvisioningRoleArn: t.ProvisioningRoleARN, Tags: t.Tags,
		CreationDate: timeToProto(t.CreationDate), LastModifiedDate: timeToProto(t.LastModifiedDate),
	}
}

func ProtoToProvisioningTemplate(p *pb.ProvisioningTemplate) *ProvisioningTemplate {
	return &ProvisioningTemplate{
		TemplateName: p.TemplateName, TemplateARN: p.TemplateArn,
		Description: p.Description, Enabled: p.Enabled,
		ProvisioningRoleARN: p.ProvisioningRoleArn, Tags: p.Tags,
		CreationDate: protoToTime(p.CreationDate), LastModifiedDate: protoToTime(p.LastModifiedDate),
	}
}

func DetectorModelToProto(d *DetectorModel) (*pb.DetectorModel, error) {
	def, err := mapToStruct(d.DetectorModelDefinition)
	if err != nil {
		return nil, fmt.Errorf("detector model definition: %w", err)
	}
	return &pb.DetectorModel{
		DetectorModelName:        d.DetectorModelName,
		DetectorModelArn:         d.DetectorModelARN,
		DetectorModelDescription: d.DetectorModelDescription,
		RoleArn:                  d.RoleARN,
		DetectorModelDefinition:  def,
		Tags:                     d.Tags,
		CreationDate:             timeToProto(d.CreationDate),
		LastModifiedDate:         timeToProto(d.LastModifiedDate),
		Status:                   d.Status,
		Key:                      d.Key,
		EvaluationMethod:         d.EvaluationMethod,
		DetectorModelVersion:     d.DetectorModelVersion,
	}, nil
}

func ProtoToDetectorModel(p *pb.DetectorModel) *DetectorModel {
	return &DetectorModel{
		DetectorModelName:        p.DetectorModelName,
		DetectorModelARN:         p.DetectorModelArn,
		DetectorModelDescription: p.DetectorModelDescription,
		RoleARN:                  p.RoleArn,
		DetectorModelDefinition:  structToMap(p.DetectorModelDefinition),
		Tags:                     p.Tags,
		CreationDate:             protoToTime(p.CreationDate),
		LastModifiedDate:         protoToTime(p.LastModifiedDate),
		Status:                   p.Status,
		Key:                      p.Key,
		EvaluationMethod:         p.EvaluationMethod,
		DetectorModelVersion:     p.DetectorModelVersion,
	}
}

func InputToProto(i *Input) (*pb.Input, error) {
	def, err := mapToStruct(i.InputDefinition)
	if err != nil {
		return nil, fmt.Errorf("input definition: %w", err)
	}
	return &pb.Input{
		InputName:        i.InputName,
		InputArn:         i.InputARN,
		InputDescription: i.InputDescription,
		InputDefinition:  def,
		Tags:             i.Tags,
		CreationDate:     timeToProto(i.CreationDate),
		LastModifiedDate: timeToProto(i.LastModifiedDate),
		Status:           i.Status,
	}, nil
}

func ProtoToInput(p *pb.Input) *Input {
	return &Input{
		InputName:        p.InputName,
		InputARN:         p.InputArn,
		InputDescription: p.InputDescription,
		InputDefinition:  structToMap(p.InputDefinition),
		Tags:             p.Tags,
		CreationDate:     protoToTime(p.CreationDate),
		LastModifiedDate: protoToTime(p.LastModifiedDate),
		Status:           p.Status,
	}
}

func timeToProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func protoToTime(p *timestamppb.Timestamp) time.Time {
	if p == nil {
		return time.Time{}
	}
	return p.AsTime()
}

// BuildThingARN constructs an ARN for an IoT thing.
func BuildThingARN(accountID, region, thingName string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("thing/%s", thingName))
}

// BuildCertificateARN constructs an ARN for an IoT certificate.
func BuildCertificateARN(accountID, region, certID string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("cert/%s", certID))
}

// BuildPolicyARN constructs an ARN for an IoT policy.
func BuildPolicyARN(accountID, region, policyName string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("policy/%s", policyName))
}

// BuildRuleARN constructs an ARN for an IoT topic rule.
func BuildRuleARN(accountID, region, ruleName string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("rule/%s", ruleName))
}

// BuildJobARN constructs an ARN for an IoT job.
func BuildJobARN(accountID, region, jobID string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("job/%s", jobID))
}

// BuildThingTypeARN constructs an ARN for an IoT thing type.
func BuildThingTypeARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("thingtype/%s", name))
}

// BuildThingGroupARN constructs an ARN for an IoT thing group.
func BuildThingGroupARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("thinggroup/%s", name))
}

// BuildRoleAliasARN constructs an ARN for an IoT role alias.
func BuildRoleAliasARN(accountID, region, alias string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("rolealias/%s", alias))
}

// BuildBillingGroupARN constructs an ARN for an IoT billing group.
func BuildBillingGroupARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("billinggroup/%s", name))
}

// BuildAuthorizerARN constructs an ARN for an IoT custom authorizer.
func BuildAuthorizerARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("authorizer/%s", name))
}

// BuildProvisioningTemplateARN constructs an ARN for a fleet provisioning template.
func BuildProvisioningTemplateARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("provisioningtemplate/%s", name))
}

// BuildDetectorModelARN constructs an ARN for an IoT Events detector model.
func BuildDetectorModelARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iotevents", fmt.Sprintf("detectorModel/%s", name))
}

// BuildInputARN constructs an ARN for an IoT Events input.
func BuildInputARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iotevents", fmt.Sprintf("input/%s", name))
}

// DetectorModel represents an IoT Events detector model.
type DetectorModel struct {
	DetectorModelName        string
	DetectorModelARN         string
	DetectorModelDescription string
	RoleARN                  string
	DetectorModelDefinition  map[string]interface{}
	Tags                     map[string]string
	CreationDate             time.Time
	LastModifiedDate         time.Time
	Status                   string
	Key                      string
	EvaluationMethod         string
	DetectorModelVersion     string
}

// Input represents an IoT Events input definition.
type Input struct {
	InputName        string
	InputARN         string
	InputDescription string
	InputDefinition  map[string]interface{}
	Tags             map[string]string
	CreationDate     time.Time
	LastModifiedDate time.Time
	Status           string
}

func mapToStruct(m map[string]interface{}) (*structpb.Struct, error) {
	if m == nil {
		return nil, nil
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil, fmt.Errorf("failed to convert map to protobuf struct: %w", err)
	}
	return s, nil
}

func structToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	return s.AsMap()
}
