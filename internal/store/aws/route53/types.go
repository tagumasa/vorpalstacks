package route53

// Package route53 provides Route 53 data store implementations for vorpalstacks.

import "time"

// HostedZone represents a Route 53 hosted zone.
type HostedZone struct {
	ID                     string            `json:"id"`
	Name                   string            `json:"name"`
	CallerReference        string            `json:"callerReference"`
	Config                 *HostedZoneConfig `json:"config,omitempty"`
	ResourceRecordSetCount int               `json:"resourceRecordSetCount"`
	Private                bool              `json:"private"`
	VPCs                   []*VPC            `json:"vpcs,omitempty"`
	DelegationSetID        string            `json:"delegationSetId,omitempty"`
	NameServers            []string          `json:"nameServers,omitempty"`
	Region                 string            `json:"region"`
	AccountID              string            `json:"accountId"`
	CreatedAt              time.Time         `json:"createdAt"`
}

// HostedZoneConfig holds the configuration for a hosted zone.
type HostedZoneConfig struct {
	Comment     string `json:"comment,omitempty"`
	PrivateZone bool   `json:"privateZone"`
}

// VPC represents an Amazon VPC associated with a hosted zone.
type VPC struct {
	VPCRegion string `json:"vpcRegion"`
	VPCID     string `json:"vpcId"`
}

// ResourceRecordSet represents a Route 53 resource record set.
type ResourceRecordSet struct {
	Name                    string            `json:"name"`
	Type                    string            `json:"type"`
	SetIdentifier           string            `json:"setIdentifier,omitempty"`
	Weight                  int64             `json:"weight,omitempty"`
	Region                  string            `json:"region,omitempty"`
	GeoLocation             *GeoLocation      `json:"geoLocation,omitempty"`
	Failover                string            `json:"failover,omitempty"`
	MultiValueAnswer        bool              `json:"multiValueAnswer,omitempty"`
	TTL                     int64             `json:"ttl,omitempty"`
	ResourceRecords         []*ResourceRecord `json:"resourceRecords,omitempty"`
	AliasTarget             *AliasTarget      `json:"aliasTarget,omitempty"`
	HealthCheckID           string            `json:"healthCheckId,omitempty"`
	TrafficPolicyInstanceID string            `json:"trafficPolicyInstanceId,omitempty"`
}

// ResourceRecord represents a DNS resource record.
type ResourceRecord struct {
	Value string `json:"value"`
}

// AliasTarget specifies the alias target for a resource record set.
type AliasTarget struct {
	HostedZoneID         string `json:"hostedZoneId"`
	DNSName              string `json:"dNSName"`
	EvaluateTargetHealth bool   `json:"evaluateTargetHealth"`
}

// GeoLocation specifies the geographic location for a geolocation routing policy.
type GeoLocation struct {
	ContinentCode   string `json:"continentCode,omitempty"`
	CountryCode     string `json:"countryCode,omitempty"`
	SubdivisionCode string `json:"subdivisionCode,omitempty"`
}

// HealthCheck represents a Route 53 health check.
type HealthCheck struct {
	ID                           string                        `json:"id"`
	CallerReference              string                        `json:"callerReference"`
	HealthCheckConfig            *HealthCheckConfig            `json:"healthCheckConfig"`
	HealthCheckVersion           string                        `json:"healthCheckVersion"`
	CloudWatchAlarmConfiguration *CloudWatchAlarmConfiguration `json:"cloudWatchAlarmConfiguration,omitempty"`
	Region                       string                        `json:"region"`
	AccountID                    string                        `json:"accountId"`
	CreatedAt                    time.Time                     `json:"createdAt"`
}

// HealthCheckConfig holds the configuration for a health check.
type HealthCheckConfig struct {
	IPAddress                    string           `json:"iPAddress,omitempty"`
	Port                         int64            `json:"port"`
	Type                         string           `json:"type"`
	ResourcePath                 string           `json:"resourcePath,omitempty"`
	FullyQualifiedDomainName     string           `json:"fullyQualifiedDomainName,omitempty"`
	SearchString                 string           `json:"searchString,omitempty"`
	RequestInterval              int64            `json:"requestInterval"`
	FailureThreshold             int64            `json:"failureThreshold"`
	MeasureLatency               bool             `json:"measureLatency"`
	Inverted                     bool             `json:"inverted"`
	Disabled                     bool             `json:"disabled"`
	EnableSNI                    bool             `json:"enableSNI"`
	Regions                      []string         `json:"regions,omitempty"`
	AlarmIdentifier              *AlarmIdentifier `json:"alarmIdentifier,omitempty"`
	InsufficientDataHealthStatus string           `json:"insufficientDataHealthStatus,omitempty"`
}

// CloudWatchAlarmConfiguration holds the CloudWatch alarm configuration for a health check.
type CloudWatchAlarmConfiguration struct {
	AlarmName   string       `json:"alarmName"`
	AlarmRegion string       `json:"alarmRegion,omitempty"`
	Dimensions  []*Dimension `json:"dimensions,omitempty"`
}

// Dimension represents a CloudWatch dimension for an alarm.
type Dimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AlarmIdentifier identifies a CloudWatch alarm.
type AlarmIdentifier struct {
	Region string `json:"region"`
	Name   string `json:"name"`
}

// ChangeInfo represents information about a change to a resource record set.
type ChangeInfo struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	SubmittedAt time.Time `json:"submittedAt"`
	Comment     string    `json:"comment,omitempty"`
}

// Change represents a change to a resource record set.
type Change struct {
	Action            string             `json:"action"`
	ResourceRecordSet *ResourceRecordSet `json:"resourceRecordSet"`
}

// ChangeBatch represents a batch of changes to resource record sets.
type ChangeBatch struct {
	Comment string    `json:"comment,omitempty"`
	Changes []*Change `json:"changes"`
}

// HostedZoneListResult represents the result of listing hosted zones.
type HostedZoneListResult struct {
	HostedZones []*HostedZone
	IsTruncated bool
	Marker      string
}

// HealthCheckListResult represents the result of listing health checks.
type HealthCheckListResult struct {
	HealthChecks []*HealthCheck
	IsTruncated  bool
	Marker       string
}

// ResourceRecordSetListResult represents the result of listing resource record sets.
type ResourceRecordSetListResult struct {
	ResourceRecordSets []*ResourceRecordSet
	IsTruncated        bool
	Marker             string
}
