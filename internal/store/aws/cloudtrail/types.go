// Package cloudtrail provides CloudTrail data store functionality for vorpalstacks.
package cloudtrail

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// Trail represents a CloudTrail trail.
type Trail struct {
	Name                       string            `json:"name"`
	TrailARN                   string            `json:"trailArn"`
	S3BucketName               string            `json:"s3BucketName,omitempty"`
	S3KeyPrefix                string            `json:"s3KeyPrefix,omitempty"`
	SnsTopicName               string            `json:"snsTopicName,omitempty"`
	SnsTopicARN                string            `json:"snsTopicArn,omitempty"`
	IncludeGlobalServiceEvents bool              `json:"includeGlobalServiceEvents"`
	IsMultiRegionTrail         bool              `json:"isMultiRegionTrail"`
	HomeRegion                 string            `json:"homeRegion"`
	IsOrganizationTrail        bool              `json:"isOrganizationTrail"`
	IsLogging                  bool              `json:"isLogging"`
	LogFileValidationEnabled   bool              `json:"logFileValidationEnabled"`
	CloudWatchLogsLogGroupARN  string            `json:"cloudWatchLogsLogGroupArn,omitempty"`
	CloudWatchLogsRoleARN      string            `json:"cloudWatchLogsRoleArn,omitempty"`
	KMSKeyID                   string            `json:"kmsKeyId,omitempty"`
	HasCustomEventSelectors    bool              `json:"hasCustomEventSelectors"`
	HasInsightSelectors        bool              `json:"hasInsightSelectors"`
	EventSelectors             []EventSelector   `json:"eventSelectors,omitempty"`
	InsightSelectors           []InsightSelector `json:"insightSelectors,omitempty"`
	CreatedAt                  time.Time         `json:"createdAt"`
	LastUpdated                time.Time         `json:"lastUpdated"`
	StartedLoggingAt           *time.Time        `json:"startedLoggingAt,omitempty"`
	StoppedLoggingAt           *time.Time        `json:"stoppedLoggingAt,omitempty"`
	Tags                       map[string]string `json:"tags,omitempty"`
}

// EventSelector represents event selector settings for a CloudTrail trail.
type EventSelector struct {
	ReadWriteType                 string         `json:"readWriteType,omitempty"`
	IncludeManagementEvents       bool           `json:"includeManagementEvents,omitempty"`
	DataResources                 []DataResource `json:"dataResources,omitempty"`
	ExcludeManagementEventSources []string       `json:"excludeManagementEventSources,omitempty"`
}

// DataResource represents a data resource for CloudTrail event selection.
type DataResource struct {
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

// InsightSelector represents insight selector settings for a CloudTrail trail.
type InsightSelector struct {
	InsightType string `json:"insightType"`
}

// Event represents a CloudTrail event.
type Event struct {
	EventID           string                 `json:"eventId"`
	EventName         string                 `json:"eventName"`
	ReadOnly          string                 `json:"readOnly"`
	AccessKeyId       string                 `json:"accessKeyId,omitempty"`
	EventSource       string                 `json:"eventSource"`
	EventTime         time.Time              `json:"eventTime"`
	EventType         string                 `json:"eventType"`
	EventVersion      string                 `json:"eventVersion"`
	UserIdentity      *UserIdentity          `json:"userIdentity"`
	Resources         []Resource             `json:"resources,omitempty"`
	CloudTrailEvent   string                 `json:"cloudTrailEvent,omitempty"`
	RequestParameters map[string]interface{} `json:"requestParameters,omitempty"`
	ResponseElements  map[string]interface{} `json:"responseElements,omitempty"`
	RequestID         string                 `json:"requestId,omitempty"`
	SourceIPAddress   string                 `json:"sourceIpAddress,omitempty"`
	UserAgent         string                 `json:"userAgent,omitempty"`
	ErrorCode         string                 `json:"errorCode,omitempty"`
	ErrorMessage      string                 `json:"errorMessage,omitempty"`
	Tags              map[string]string      `json:"tags,omitempty"`
}

// UserIdentity represents the identity information for a CloudTrail event.
type UserIdentity struct {
	Type           string          `json:"type"`
	PrincipalID    string          `json:"principalId,omitempty"`
	ARN            string          `json:"arn,omitempty"`
	AccountID      string          `json:"accountId,omitempty"`
	AccessKeyID    string          `json:"accessKeyId,omitempty"`
	UserName       string          `json:"userName,omitempty"`
	SessionContext *SessionContext `json:"sessionContext,omitempty"`
}

// SessionContext represents session context information for a CloudTrail event.
type SessionContext struct {
	SessionIssuer *SessionIssuer     `json:"sessionIssuer,omitempty"`
	Attributes    *SessionAttributes `json:"attributes,omitempty"`
}

// SessionIssuer represents session issuer information for a CloudTrail event.
type SessionIssuer struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
	ARN         string `json:"arn"`
	UserName    string `json:"userName,omitempty"`
}

// SessionAttributes represents session attributes for a CloudTrail event.
type SessionAttributes struct {
	CreationDate     time.Time `json:"creationDate"`
	MFAAuthenticated string    `json:"mfaAuthenticated"`
}

// Resource represents a resource associated with a CloudTrail event.
type Resource struct {
	ResourceType string `json:"resourceType,omitempty"`
	ResourceName string `json:"resourceName,omitempty"`
	ARN          string `json:"ARN,omitempty"`
}

// ResourcePolicy represents a resource policy for CloudTrail.
type ResourcePolicy struct {
	ResourceARN string `json:"resourceArn"`
	Policy      string `json:"policy"`
}

// PublicKey represents a CloudTrail public key used for log file validation.
type PublicKey struct {
	PublicKeyID       string    `json:"publicKeyId"`
	Value             []byte    `json:"value"`
	ValidityStartTime time.Time `json:"validityStartTime"`
	ValidityEndTime   time.Time `json:"validityEndTime"`
	TrailName         string    `json:"trailName,omitempty"`
}

// GenerateKeyPair generates an RSA key pair using the common crypto utility.
// Returns the DER-encoded public key bytes suitable for CloudTrail log file validation.
func GenerateKeyPair() (*PublicKey, error) {
	privKey, err := vcrypto.GenerateRSAKey(vcrypto.DefaultRSAKeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	derBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key to DER: %w", err)
	}

	now := time.Now().UTC()
	return &PublicKey{
		PublicKeyID:       "K" + uuid.NewString(),
		Value:             derBytes,
		ValidityStartTime: now,
		ValidityEndTime:   now.Add(365 * 24 * time.Hour),
	}, nil
}

// NewTrail creates a new CloudTrail trail with default settings.
func NewTrail(name, s3BucketName, region string) *Trail {
	now := time.Now().UTC()
	return &Trail{
		Name:                       name,
		S3BucketName:               s3BucketName,
		HomeRegion:                 region,
		IncludeGlobalServiceEvents: true,
		IsLogging:                  false,
		LogFileValidationEnabled:   false,
		HasCustomEventSelectors:    false,
		HasInsightSelectors:        false,
		EventSelectors: []EventSelector{
			{
				ReadWriteType:           "All",
				IncludeManagementEvents: true,
			},
		},
		CreatedAt:   now,
		LastUpdated: now,
		Tags:        make(map[string]string),
	}
}

// NewEvent creates a new CloudTrail event.
func NewEvent(eventName, eventSource string, userIdentity *UserIdentity) *Event {
	now := time.Now().UTC()
	return &Event{
		EventID:      generateEventID(),
		EventName:    eventName,
		EventSource:  eventSource,
		EventTime:    now,
		EventType:    "AwsApiCall",
		EventVersion: "1.08",
		UserIdentity: userIdentity,
		Resources:    []Resource{},
		ReadOnly:     "false",
	}
}

// NewReadOnlyEvent creates a new read-only CloudTrail event.
func NewReadOnlyEvent(eventName, eventSource string, userIdentity *UserIdentity) *Event {
	return &Event{
		EventID:      generateEventID(),
		EventName:    eventName,
		EventSource:  eventSource,
		EventTime:    time.Now().UTC(),
		EventType:    "AwsApiCall",
		EventVersion: "1.08",
		UserIdentity: userIdentity,
		Resources:    []Resource{},
		ReadOnly:     "true",
	}
}

func generateEventID() string {
	return uuid.NewString()
}

// generateCloudTrailEvent populates the CloudTrailEvent field with a JSON
// representation of the event, matching the AWS API response format.
func (e *Event) generateCloudTrailEvent() {
	if e.CloudTrailEvent != "" {
		return
	}
	clone := *e
	clone.CloudTrailEvent = ""
	b, err := json.Marshal(clone)
	if err != nil {
		return
	}
	e.CloudTrailEvent = string(b)
}
