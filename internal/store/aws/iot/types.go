package iot

import (
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "vorpalstacks/internal/pb/storage/storage_iot"
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
	DefaultClientId  string
	BillingGroupName string
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
		DefaultClientId:  t.DefaultClientId,
		BillingGroupName: t.BillingGroupName,
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
		DefaultClientId:  p.DefaultClientId,
		BillingGroupName: p.BillingGroupName,
	}
}

type ThingType struct {
	ThingTypeName        string
	ThingTypeARN         string
	ThingTypeID          string
	Description          string
	ThingTypeProperties  []ThingTypeProperty
	SearchableAttributes []string
	Tags                 map[string]string
	Version              int64
	CreationDate         time.Time
	LastModifiedDate     time.Time
	Deprecated           bool
	DeprecationDate      time.Time
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
		ThingTypeName:        t.ThingTypeName,
		ThingTypeArn:         t.ThingTypeARN,
		ThingTypeId:          t.ThingTypeID,
		Description:          t.Description,
		ThingTypeProperties:  props,
		SearchableAttributes: t.SearchableAttributes,
		Tags:                 t.Tags,
		Version:              t.Version,
		CreationDate:         timeToProto(t.CreationDate),
		LastModifiedDate:     timeToProto(t.LastModifiedDate),
		Deprecated:           t.Deprecated,
		DeprecationDate:      timeToProto(t.DeprecationDate),
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
		ThingTypeName:        p.ThingTypeName,
		ThingTypeARN:         p.ThingTypeArn,
		ThingTypeID:          p.ThingTypeId,
		Description:          p.Description,
		ThingTypeProperties:  props,
		SearchableAttributes: p.SearchableAttributes,
		Tags:                 p.Tags,
		Version:              p.Version,
		CreationDate:         protoToTime(p.CreationDate),
		LastModifiedDate:     protoToTime(p.LastModifiedDate),
		Deprecated:           p.Deprecated,
		DeprecationDate:      protoToTime(p.DeprecationDate),
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
	Version          int64
}

func BillingGroupToProto(g *BillingGroup) *pb.BillingGroup {
	return &pb.BillingGroup{
		GroupName: g.GroupName, GroupArn: g.GroupARN, GroupId: g.GroupID,
		Description: g.Description, Attributes: g.Attributes, Tags: g.Tags,
		CreationDate: timeToProto(g.CreationDate), LastModifiedDate: timeToProto(g.LastModifiedDate),
		Version: g.Version,
	}
}

func ProtoToBillingGroup(p *pb.BillingGroup) *BillingGroup {
	return &BillingGroup{
		GroupName: p.GroupName, GroupARN: p.GroupArn, GroupID: p.GroupId,
		Description: p.Description, Attributes: p.Attributes, Tags: p.Tags,
		CreationDate: protoToTime(p.CreationDate), LastModifiedDate: protoToTime(p.LastModifiedDate),
		Version: p.Version,
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
	Actions          map[string]interface{}
	ErrorAction      map[string]interface{}
}

func RuleToProto(r *TopicRule) *pb.TopicRule {
	p := &pb.TopicRule{
		RuleName: r.RuleName, Arn: r.ARN, TopicPattern: r.TopicPattern,
		Description: r.Description, RuleDisabled: r.RuleDisabled,
		Sql: r.SQL, CreatedAt: r.CreatedAt, AwsIotSqlVersion: r.AwsIotSqlVersion,
	}
	if r.Actions != nil {
		if s, err := mapToStruct(r.Actions); err == nil {
			p.Actions = s
		}
	}
	if r.ErrorAction != nil {
		if s, err := mapToStruct(r.ErrorAction); err == nil {
			p.ErrorAction = s
		}
	}
	return p
}

func ProtoToRule(p *pb.TopicRule) *TopicRule {
	return &TopicRule{
		RuleName: p.RuleName, ARN: p.Arn, TopicPattern: p.TopicPattern,
		Description: p.Description, RuleDisabled: p.RuleDisabled,
		SQL: p.Sql, CreatedAt: p.CreatedAt, AwsIotSqlVersion: p.AwsIotSqlVersion,
		Actions:     structToMap(p.Actions),
		ErrorAction: structToMap(p.ErrorAction),
	}
}

type Job struct {
	JobARN                     string
	JobID                      string
	Description                string
	Version                    int64
	Force                      bool
	CreatedAt                  time.Time
	LastUpdatedAt              time.Time
	CompletedAt                string
	JobTemplateARN             string
	Status                     string
	TargetSelection            string
	Tags                       map[string]string
	Document                   string
	Targets                    []string
	ReasonCode                 string
	Comment                    string
	NamespaceID                string
	ForceCanceledFlag          bool
	IsConcurrent               bool
	PresignedUrlConfig         string
	JobExecutionsRolloutConfig string
	AbortConfig                string
	TimeoutConfig              string
	JobExecutionsRetryConfig   string
	DocumentParameters         string
	SchedulingConfig           string
	ScheduledJobRollouts       string
	DestinationPackageVersions string
}

func JobToProto(j *Job) *pb.Job {
	return &pb.Job{
		JobArn: j.JobARN, JobId: j.JobID, Description: j.Description,
		Version: j.Version, Force: j.Force,
		CreatedAt: timeToProto(j.CreatedAt), LastUpdatedAt: timeToProto(j.LastUpdatedAt),
		CompletedAt: j.CompletedAt, JobTemplateArn: j.JobTemplateARN,
		Status: j.Status, TargetSelection: j.TargetSelection, Tags: j.Tags,
		Document: j.Document, Targets: j.Targets,
		ReasonCode: j.ReasonCode, Comment: j.Comment, NamespaceId: j.NamespaceID,
		ForceCanceledFlag: j.ForceCanceledFlag, IsConcurrent: j.IsConcurrent,
		PresignedUrlConfig: j.PresignedUrlConfig, JobExecutionsRolloutConfig: j.JobExecutionsRolloutConfig,
		AbortConfig: j.AbortConfig, TimeoutConfig: j.TimeoutConfig,
		JobExecutionsRetryConfig: j.JobExecutionsRetryConfig, DocumentParameters: j.DocumentParameters,
		SchedulingConfig: j.SchedulingConfig, ScheduledJobRollouts: j.ScheduledJobRollouts,
		DestinationPackageVersions: j.DestinationPackageVersions,
	}
}

func ProtoToJob(p *pb.Job) *Job {
	return &Job{
		JobARN: p.JobArn, JobID: p.JobId, Description: p.Description,
		Version: p.Version, Force: p.Force,
		CreatedAt: protoToTime(p.CreatedAt), LastUpdatedAt: protoToTime(p.LastUpdatedAt),
		CompletedAt: p.CompletedAt, JobTemplateARN: p.JobTemplateArn,
		Status: p.Status, TargetSelection: p.TargetSelection, Tags: p.Tags,
		Document: p.Document, Targets: p.Targets,
		ReasonCode: p.ReasonCode, Comment: p.Comment, NamespaceID: p.NamespaceId,
		ForceCanceledFlag: p.ForceCanceledFlag, IsConcurrent: p.IsConcurrent,
		PresignedUrlConfig: p.PresignedUrlConfig, JobExecutionsRolloutConfig: p.JobExecutionsRolloutConfig,
		AbortConfig: p.AbortConfig, TimeoutConfig: p.TimeoutConfig,
		JobExecutionsRetryConfig: p.JobExecutionsRetryConfig, DocumentParameters: p.DocumentParameters,
		SchedulingConfig: p.SchedulingConfig, ScheduledJobRollouts: p.ScheduledJobRollouts,
		DestinationPackageVersions: p.DestinationPackageVersions,
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
	AuthorizerName         string
	AuthorizerARN          string
	AuthorizerFunctionARN  string
	TokenName              string
	TokenSignature         string
	Status                 bool
	CreationDate           time.Time
	LastModifiedDate       time.Time
	EnableCachingForHTTP   bool
	CachingDisabled        int64
	TokenSigningPublicKeys map[string]string
	SigningDisabled        bool
}

func AuthorizerToProto(a *Authorizer) *pb.Authorizer {
	return &pb.Authorizer{
		AuthorizerName: a.AuthorizerName, AuthorizerArn: a.AuthorizerARN,
		AuthorizerFunctionArn: a.AuthorizerFunctionARN,
		TokenKeyName:          a.TokenName, TokenSignature: a.TokenSignature, Status: a.Status,
		CreationDate: timeToProto(a.CreationDate), LastModifiedDate: timeToProto(a.LastModifiedDate),
		EnableCachingForHttp: a.EnableCachingForHTTP, CachingDisabled: a.CachingDisabled,
		TokenSigningPublicKeys: a.TokenSigningPublicKeys, SigningDisabled: a.SigningDisabled,
	}
}

func ProtoToAuthorizer(p *pb.Authorizer) *Authorizer {
	return &Authorizer{
		AuthorizerName: p.AuthorizerName, AuthorizerARN: p.AuthorizerArn,
		AuthorizerFunctionARN: p.AuthorizerFunctionArn,
		TokenName:             p.TokenKeyName, TokenSignature: p.TokenSignature, Status: p.Status,
		CreationDate: protoToTime(p.CreationDate), LastModifiedDate: protoToTime(p.LastModifiedDate),
		EnableCachingForHTTP: p.EnableCachingForHttp, CachingDisabled: p.CachingDisabled,
		TokenSigningPublicKeys: p.TokenSigningPublicKeys, SigningDisabled: p.SigningDisabled,
	}
}

type RoleAlias struct {
	RoleAlias                 string
	RoleAliasARN              string
	RoleARN                   string
	CredentialDurationSeconds int64
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
	TemplateBody        string
	Tags                map[string]string
	CreationDate        time.Time
	LastModifiedDate    time.Time
	Type                string
	PreProvisioningHook string
	DefaultVersionID    int64
}

func ProvisioningTemplateToProto(t *ProvisioningTemplate) (*pb.ProvisioningTemplate, error) {
	out := &pb.ProvisioningTemplate{
		TemplateName: t.TemplateName, TemplateArn: t.TemplateARN,
		Description: t.Description, Enabled: t.Enabled,
		ProvisioningRoleArn: t.ProvisioningRoleARN, Tags: t.Tags,
		CreationDate: timeToProto(t.CreationDate), LastModifiedDate: timeToProto(t.LastModifiedDate),
		Type: t.Type, DefaultVersionId: t.DefaultVersionID,
	}
	if t.PreProvisioningHook != "" {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(t.PreProvisioningHook), &m); err != nil {
			return nil, fmt.Errorf("provisioning template %q hook unmarshal: %w", t.TemplateName, err)
		}
		s, err := mapToStruct(m)
		if err != nil {
			return nil, fmt.Errorf("provisioning template %q hook to struct: %w", t.TemplateName, err)
		}
		out.PreProvisioningHook = s
	}
	if t.TemplateBody != "" {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(t.TemplateBody), &m); err != nil {
			return nil, fmt.Errorf("provisioning template %q body unmarshal: %w", t.TemplateName, err)
		}
		s, err := mapToStruct(m)
		if err != nil {
			return nil, fmt.Errorf("provisioning template %q body to struct: %w", t.TemplateName, err)
		}
		out.TemplateBody = s
	}
	return out, nil
}

func ProtoToProvisioningTemplate(p *pb.ProvisioningTemplate) (*ProvisioningTemplate, error) {
	t := &ProvisioningTemplate{
		TemplateName: p.TemplateName, TemplateARN: p.TemplateArn,
		Description: p.Description, Enabled: p.Enabled,
		ProvisioningRoleARN: p.ProvisioningRoleArn, Tags: p.Tags,
		CreationDate: protoToTime(p.CreationDate), LastModifiedDate: protoToTime(p.LastModifiedDate),
		Type: p.Type, DefaultVersionID: p.DefaultVersionId,
	}
	if p.PreProvisioningHook != nil {
		data, err := json.Marshal(structToMap(p.PreProvisioningHook))
		if err != nil {
			return nil, fmt.Errorf("provisioning template %q hook marshal: %w", p.TemplateName, err)
		}
		t.PreProvisioningHook = string(data)
	}
	if p.TemplateBody != nil {
		data, err := json.Marshal(structToMap(p.TemplateBody))
		if err != nil {
			return nil, fmt.Errorf("provisioning template %q body marshal: %w", p.TemplateName, err)
		}
		t.TemplateBody = string(data)
	}
	return t, nil
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

// SecurityProfile represents a Device Defender security profile.
type SecurityProfile struct {
	SecurityProfileName         string
	SecurityProfileARN          string
	SecurityProfileDescription  string
	Behaviors                   []*Behavior
	AlertTargets                map[string]*AlertTarget
	AdditionalMetricsToRetainV2 []*MetricToRetain
	Version                     int64
	CreationDate                time.Time
	LastModifiedDate            time.Time
	Tags                        map[string]string
	MetricsExportConfig         string
	AdditionalMetricsToRetain   []string
}

// AlertTarget represents an alert destination for security profile violations.
type AlertTarget struct {
	AlertTargetARN string
	RoleARN        string
}

// MetricToRetain is one retained-metric entry of a security profile: the
// metric name, an optional scoping dimension (name plus the optional
// IN/NOT_IN operator) and the export flag.
type MetricToRetain struct {
	Metric          string
	MetricDimension string
	Operator        string
	ExportMetric    bool
}

// Behavior represents a security profile behaviour definition.
type Behavior struct {
	Name            string
	Metric          string
	MetricDimension string
	Criteria        *BehaviorCriteria
	SuppressAlerts  bool
	ExportMetric    bool
}

// BehaviorCriteria defines the conditions that trigger a behaviour alert.
type BehaviorCriteria struct {
	ComparisonOperator           string
	Value                        float64
	DurationSeconds              int64
	ConsecutiveDatapointsToAlarm int64
	ConsecutiveDatapointsToClear int64
	StatisticalThreshold         *StatisticalThreshold
	MLDetectionConfig            *MachineLearningDetectionConfig
}

// StatisticalThreshold configures statistical anomaly detection.
type StatisticalThreshold struct {
	Statistic string
}

// MachineLearningDetectionConfig configures ML-based anomaly detection.
type MachineLearningDetectionConfig struct {
	ConfidenceLevel string
}

// MetricValue represents a metric measurement in a violation event.
type MetricValue struct {
	Count   int64
	Cidrs   []string
	Ports   []int64
	Number  float64
	Numbers []float64
	Strings []string
}

// ViolationEvent represents a Device Defender violation event.
type ViolationEvent struct {
	ViolationID                  string
	ThingName                    string
	SecurityProfileName          string
	Behavior                     *Behavior
	MetricValue                  *MetricValue
	ViolationEventType           string
	VerificationState            string
	VerificationStateDescription string
	ViolationEventTime           time.Time
}

// DomainConfiguration represents an IoT domain configuration.
type DomainConfiguration struct {
	DomainConfigurationName   string
	DomainConfigurationARN    string
	DomainName                string
	ServerCertificateARNs     []string
	ValidationCertificateARN  string
	AuthorizerConfig          string
	ServiceType               string
	DomainConfigurationStatus string
	AuthenticationType        string
	ApplicationProtocol       string
	CreationDate              time.Time
	LastModifiedDate          time.Time
	Tags                      map[string]string
}

// IndexingConfiguration represents the IoT thing indexing configuration.
type IndexingConfiguration struct {
	ThingIndexingMode                 string
	ThingGroupIndexingMode            string
	ThingConnectivityIndexingMode     string
	DeviceDefenderIndexingMode        string
	NamedShadowIndexingMode           string
	ManagedFields                     []string
	CustomFields                      []string
	ThingIndexingConfigurationVersion string
}

// ProvisioningTemplateVersion represents a version of a provisioning template.
type ProvisioningTemplateVersion struct {
	VersionID        string
	CreationDate     time.Time
	IsDefaultVersion bool
	TemplateBody     string
}

func SecurityProfileToProto(s *SecurityProfile) (*pb.SecurityProfile, error) {
	pbBehaviors := make([]*pb.Behavior, 0, len(s.Behaviors))
	for _, b := range s.Behaviors {
		pbBehavior, err := BehaviorToProto(b)
		if err != nil {
			return nil, fmt.Errorf("behavior %s: %w", b.Name, err)
		}
		pbBehaviors = append(pbBehaviors, pbBehavior)
	}
	pbAlertTargets := make(map[string]*pb.AlertTarget, len(s.AlertTargets))
	for k, v := range s.AlertTargets {
		pbAlertTargets[k] = &pb.AlertTarget{
			AlertTargetArn: v.AlertTargetARN,
			RoleArn:        v.RoleARN,
		}
	}
	// The retained-metric entries are structures on the wire; the storage
	// proto carries them as one JSON-encoded string per entry.
	pbMetricsV2 := make([]string, 0, len(s.AdditionalMetricsToRetainV2))
	for _, m := range s.AdditionalMetricsToRetainV2 {
		encoded, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		pbMetricsV2 = append(pbMetricsV2, string(encoded))
	}
	return &pb.SecurityProfile{
		SecurityProfileName:         s.SecurityProfileName,
		SecurityProfileArn:          s.SecurityProfileARN,
		SecurityProfileDescription:  s.SecurityProfileDescription,
		Behaviors:                   pbBehaviors,
		AlertTargets:                pbAlertTargets,
		AdditionalMetricsToRetainV2: pbMetricsV2,
		Version:                     s.Version,
		CreationDate:                timeToProto(s.CreationDate),
		LastModifiedDate:            timeToProto(s.LastModifiedDate),
		Tags:                        s.Tags,
		MetricsExportConfig:         s.MetricsExportConfig,
		AdditionalMetricsToRetain:   s.AdditionalMetricsToRetain,
	}, nil
}

func ProtoToSecurityProfile(p *pb.SecurityProfile) *SecurityProfile {
	domainBehaviors := make([]*Behavior, 0, len(p.Behaviors))
	for _, b := range p.Behaviors {
		domainBehaviors = append(domainBehaviors, ProtoToBehavior(b))
	}
	domainAlertTargets := make(map[string]*AlertTarget, len(p.AlertTargets))
	for k, v := range p.AlertTargets {
		domainAlertTargets[k] = &AlertTarget{
			AlertTargetARN: v.AlertTargetArn,
			RoleARN:        v.RoleArn,
		}
	}
	domainMetricsV2 := make([]*MetricToRetain, 0, len(p.AdditionalMetricsToRetainV2))
	for _, entry := range p.AdditionalMetricsToRetainV2 {
		var m MetricToRetain
		if err := json.Unmarshal([]byte(entry), &m); err != nil {
			continue
		}
		domainMetricsV2 = append(domainMetricsV2, &m)
	}
	return &SecurityProfile{
		SecurityProfileName:         p.SecurityProfileName,
		SecurityProfileARN:          p.SecurityProfileArn,
		SecurityProfileDescription:  p.SecurityProfileDescription,
		Behaviors:                   domainBehaviors,
		AlertTargets:                domainAlertTargets,
		AdditionalMetricsToRetainV2: domainMetricsV2,
		Version:                     p.Version,
		CreationDate:                protoToTime(p.CreationDate),
		LastModifiedDate:            protoToTime(p.LastModifiedDate),
		Tags:                        p.Tags,
		MetricsExportConfig:         p.MetricsExportConfig,
		AdditionalMetricsToRetain:   p.AdditionalMetricsToRetain,
	}
}

func BehaviorToProto(b *Behavior) (*pb.Behavior, error) {
	var pbCriteria *pb.BehaviorCriteria
	if b.Criteria != nil {
		pbCriteria = &pb.BehaviorCriteria{
			ComparisonOperator:           b.Criteria.ComparisonOperator,
			Value:                        b.Criteria.Value,
			DurationSeconds:              b.Criteria.DurationSeconds,
			ConsecutiveDatapointsToAlarm: b.Criteria.ConsecutiveDatapointsToAlarm,
			ConsecutiveDatapointsToClear: b.Criteria.ConsecutiveDatapointsToClear,
		}
		if b.Criteria.StatisticalThreshold != nil {
			pbCriteria.StatisticalThreshold = &pb.StatisticalThreshold{
				Statistic: b.Criteria.StatisticalThreshold.Statistic,
			}
		}
		if b.Criteria.MLDetectionConfig != nil {
			pbCriteria.MlDetectionConfig = &pb.MachineLearningDetectionConfig{
				ConfidenceLevel: b.Criteria.MLDetectionConfig.ConfidenceLevel,
			}
		}
	}
	return &pb.Behavior{
		Name:            b.Name,
		Metric:          b.Metric,
		MetricDimension: b.MetricDimension,
		Criteria:        pbCriteria,
		SuppressAlerts:  b.SuppressAlerts,
		ExportMetric:    b.ExportMetric,
	}, nil
}

func ProtoToBehavior(p *pb.Behavior) *Behavior {
	var domainCriteria *BehaviorCriteria
	if p.Criteria != nil {
		domainCriteria = &BehaviorCriteria{
			ComparisonOperator:           p.Criteria.ComparisonOperator,
			Value:                        p.Criteria.Value,
			DurationSeconds:              p.Criteria.DurationSeconds,
			ConsecutiveDatapointsToAlarm: p.Criteria.ConsecutiveDatapointsToAlarm,
			ConsecutiveDatapointsToClear: p.Criteria.ConsecutiveDatapointsToClear,
		}
		if p.Criteria.StatisticalThreshold != nil {
			domainCriteria.StatisticalThreshold = &StatisticalThreshold{
				Statistic: p.Criteria.StatisticalThreshold.Statistic,
			}
		}
		if p.Criteria.MlDetectionConfig != nil {
			domainCriteria.MLDetectionConfig = &MachineLearningDetectionConfig{
				ConfidenceLevel: p.Criteria.MlDetectionConfig.ConfidenceLevel,
			}
		}
	}
	return &Behavior{
		Name:            p.Name,
		Metric:          p.Metric,
		MetricDimension: p.MetricDimension,
		Criteria:        domainCriteria,
		SuppressAlerts:  p.SuppressAlerts,
		ExportMetric:    p.ExportMetric,
	}
}

func MetricValueToProto(m *MetricValue) *pb.MetricValue {
	return &pb.MetricValue{
		Count:   m.Count,
		Cidrs:   m.Cidrs,
		Ports:   m.Ports,
		Number:  m.Number,
		Numbers: m.Numbers,
		Strings: m.Strings,
	}
}

func ProtoToMetricValue(p *pb.MetricValue) *MetricValue {
	return &MetricValue{
		Count:   p.Count,
		Cidrs:   p.Cidrs,
		Ports:   p.Ports,
		Number:  p.Number,
		Numbers: p.Numbers,
		Strings: p.Strings,
	}
}

func ViolationEventToProto(v *ViolationEvent) (*pb.ViolationEvent, error) {
	var pbBehavior *pb.Behavior
	if v.Behavior != nil {
		b, err := BehaviorToProto(v.Behavior)
		if err != nil {
			return nil, fmt.Errorf("behavior: %w", err)
		}
		pbBehavior = b
	}
	return &pb.ViolationEvent{
		ViolationId:                  v.ViolationID,
		ThingName:                    v.ThingName,
		SecurityProfileName:          v.SecurityProfileName,
		Behavior:                     pbBehavior,
		MetricValue:                  MetricValueToProto(v.MetricValue),
		ViolationEventType:           v.ViolationEventType,
		VerificationState:            v.VerificationState,
		VerificationStateDescription: v.VerificationStateDescription,
		ViolationEventTime:           timeToProto(v.ViolationEventTime),
	}, nil
}

func ProtoToViolationEvent(p *pb.ViolationEvent) *ViolationEvent {
	return &ViolationEvent{
		ViolationID:                  p.ViolationId,
		ThingName:                    p.ThingName,
		SecurityProfileName:          p.SecurityProfileName,
		Behavior:                     ProtoToBehavior(p.Behavior),
		MetricValue:                  ProtoToMetricValue(p.MetricValue),
		ViolationEventType:           p.ViolationEventType,
		VerificationState:            p.VerificationState,
		VerificationStateDescription: p.VerificationStateDescription,
		ViolationEventTime:           protoToTime(p.ViolationEventTime),
	}
}

func DomainConfigurationToProto(d *DomainConfiguration) (*pb.DomainConfiguration, error) {
	return &pb.DomainConfiguration{
		DomainConfigurationName:   d.DomainConfigurationName,
		DomainConfigurationArn:    d.DomainConfigurationARN,
		DomainName:                d.DomainName,
		ServerCertificateArns:     d.ServerCertificateARNs,
		ValidationCertificateArn:  d.ValidationCertificateARN,
		AuthorizerConfig:          d.AuthorizerConfig,
		ServiceType:               d.ServiceType,
		DomainConfigurationStatus: d.DomainConfigurationStatus,
		AuthenticationType:        d.AuthenticationType,
		ApplicationProtocol:       d.ApplicationProtocol,
		CreationDate:              timeToProto(d.CreationDate),
		LastModifiedDate:          timeToProto(d.LastModifiedDate),
		Tags:                      d.Tags,
	}, nil
}

func ProtoToDomainConfiguration(p *pb.DomainConfiguration) *DomainConfiguration {
	return &DomainConfiguration{
		DomainConfigurationName:   p.DomainConfigurationName,
		DomainConfigurationARN:    p.DomainConfigurationArn,
		DomainName:                p.DomainName,
		ServerCertificateARNs:     p.ServerCertificateArns,
		ValidationCertificateARN:  p.ValidationCertificateArn,
		AuthorizerConfig:          p.AuthorizerConfig,
		ServiceType:               p.ServiceType,
		DomainConfigurationStatus: p.DomainConfigurationStatus,
		AuthenticationType:        p.AuthenticationType,
		ApplicationProtocol:       p.ApplicationProtocol,
		CreationDate:              protoToTime(p.CreationDate),
		LastModifiedDate:          protoToTime(p.LastModifiedDate),
		Tags:                      p.Tags,
	}
}

func IndexingConfigurationToProto(i *IndexingConfiguration) *pb.IndexingConfiguration {
	return &pb.IndexingConfiguration{
		ThingIndexingMode:                 i.ThingIndexingMode,
		ThingGroupIndexingMode:            i.ThingGroupIndexingMode,
		ThingConnectivityIndexingMode:     i.ThingConnectivityIndexingMode,
		DeviceDefenderIndexingMode:        i.DeviceDefenderIndexingMode,
		NamedShadowIndexingMode:           i.NamedShadowIndexingMode,
		ManagedFields:                     i.ManagedFields,
		CustomFields:                      i.CustomFields,
		ThingIndexingConfigurationVersion: i.ThingIndexingConfigurationVersion,
	}
}

func ProtoToIndexingConfiguration(p *pb.IndexingConfiguration) *IndexingConfiguration {
	return &IndexingConfiguration{
		ThingIndexingMode:                 p.ThingIndexingMode,
		ThingGroupIndexingMode:            p.ThingGroupIndexingMode,
		ThingConnectivityIndexingMode:     p.ThingConnectivityIndexingMode,
		DeviceDefenderIndexingMode:        p.DeviceDefenderIndexingMode,
		NamedShadowIndexingMode:           p.NamedShadowIndexingMode,
		ManagedFields:                     p.ManagedFields,
		CustomFields:                      p.CustomFields,
		ThingIndexingConfigurationVersion: p.ThingIndexingConfigurationVersion,
	}
}

func ProvisioningTemplateVersionToProto(v *ProvisioningTemplateVersion) (*pb.ProvisioningTemplateVersion, error) {
	return &pb.ProvisioningTemplateVersion{
		VersionId:        v.VersionID,
		CreationDate:     timeToProto(v.CreationDate),
		IsDefaultVersion: v.IsDefaultVersion,
		TemplateBody:     v.TemplateBody,
	}, nil
}

func ProtoToProvisioningTemplateVersion(p *pb.ProvisioningTemplateVersion) *ProvisioningTemplateVersion {
	return &ProvisioningTemplateVersion{
		VersionID:        p.VersionId,
		CreationDate:     protoToTime(p.CreationDate),
		IsDefaultVersion: p.IsDefaultVersion,
		TemplateBody:     p.TemplateBody,
	}
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

func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
