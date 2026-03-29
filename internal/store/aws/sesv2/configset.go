package sesv2

import (
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
)

// SendingOptions specifies the sending options for a configuration set.
type SendingOptions struct {
	SendingEnabled bool `json:"sendingEnabled"`
}

// ReputationOptions specifies the reputation options for a configuration set.
type ReputationOptions struct {
	ReputationMetricsEnabled bool      `json:"reputationMetricsEnabled"`
	LastFreshStart           time.Time `json:"lastFreshStart,omitempty"`
}

// DeliveryOptions specifies the delivery options for a configuration set.
type DeliveryOptions struct {
	SendingPoolName    string `json:"sendingPoolName,omitempty"`
	MaxDeliverySeconds int32  `json:"maxDeliverySeconds,omitempty"`
}

// TrackingOptions specifies the tracking options for a configuration set.
type TrackingOptions struct {
	CustomRedirectDomain string `json:"customRedirectDomain,omitempty"`
	HttpsPolicy          string `json:"httpsPolicy,omitempty"`
}

// SuppressionOptions specifies the suppression options for a configuration set.
type SuppressionOptions struct {
	SuppressedReasons []string `json:"suppressedReasons,omitempty"`
}

// VdmOptions specifies the Verifiable Deliverability Metrics options for a configuration set.
type VdmOptions struct {
	DashboardOptions *VDMDashboardOptions `json:"dashboardOptions,omitempty"`
	GuardianOptions  *VDMGuardianOptions  `json:"guardianOptions,omitempty"`
}

// VDMDashboardOptions specifies the dashboard options for VDM.
type VDMDashboardOptions struct {
	EngagementMetrics string `json:"engagementMetrics,omitempty"`
}

// VDMGuardianOptions specifies the guardian options for VDM.
type VDMGuardianOptions struct {
	OptimizedSharedDelivery string `json:"optimizedSharedDelivery,omitempty"`
}

// ConfigurationSet represents an SES v2 configuration set.
type ConfigurationSet struct {
	Name                 string              `json:"name"`
	SendingOptions       *SendingOptions     `json:"sendingOptions,omitempty"`
	ReputationOptions    *ReputationOptions  `json:"reputationOptions,omitempty"`
	DeliveryOptions      *DeliveryOptions    `json:"deliveryOptions,omitempty"`
	TrackingOptions      *TrackingOptions    `json:"trackingOptions,omitempty"`
	SuppressionOptions   *SuppressionOptions `json:"suppressionOptions,omitempty"`
	VdmOptions           *VdmOptions         `json:"vdmOptions,omitempty"`
	CreatedTimestamp     time.Time           `json:"createdTimestamp"`
	LastUpdatedTimestamp time.Time           `json:"lastUpdatedTimestamp"`
}

// EventDestinationDefinition defines an event destination.
type EventDestinationDefinition struct {
	Enabled                    bool                        `json:"enabled"`
	MatchingEventTypes         []string                    `json:"matchingEventTypes,omitempty"`
	CloudWatchDestination      *CloudWatchDestination      `json:"cloudWatchDestination,omitempty"`
	KinesisFirehoseDestination *KinesisFirehoseDestination `json:"kinesisFirehoseDestination,omitempty"`
	PinpointDestination        *PinpointDestination        `json:"pinpointDestination,omitempty"`
	SnsDestination             *SnsDestination             `json:"snsDestination,omitempty"`
	EventBridgeDestination     *EventBridgeDestination     `json:"eventBridgeDestination,omitempty"`
}

// CloudWatchDestination specifies a CloudWatch event destination.
type CloudWatchDestination struct {
	DimensionConfigurations []CloudWatchDimensionConfiguration `json:"dimensionConfigurations,omitempty"`
}

// CloudWatchDimensionConfiguration specifies dimension configuration for CloudWatch.
type CloudWatchDimensionConfiguration struct {
	DimensionName         string `json:"dimensionName"`
	DimensionValueSource  string `json:"dimensionValueSource"`
	DefaultDimensionValue string `json:"defaultDimensionValue"`
}

// KinesisFirehoseDestination specifies a Kinesis Firehose event destination.
type KinesisFirehoseDestination struct {
	DeliveryStreamARN string `json:"deliveryStreamArn"`
	IAMRoleARN        string `json:"iamRoleArn"`
}

// PinpointDestination specifies a Pinpoint event destination.
type PinpointDestination struct {
	ApplicationARN string `json:"applicationArn"`
}

// SnsDestination specifies an SNS event destination.
type SnsDestination struct {
	TopicARN string `json:"topicArn"`
}

// EventBridgeDestination specifies an EventBridge event destination.
type EventBridgeDestination struct {
	EventBusARN string `json:"eventBusArn"`
}

// EventDestination represents an SES event destination.
type EventDestination struct {
	Name                       string                      `json:"name"`
	Enabled                    bool                        `json:"enabled"`
	MatchingEventTypes         []string                    `json:"matchingEventTypes,omitempty"`
	EventDestinationDefinition *EventDestinationDefinition `json:"eventDestinationDefinition,omitempty"`
}

// NewConfigurationSet creates a new configuration set with default values.
func NewConfigurationSet(name string) *ConfigurationSet {
	now := time.Now().UTC()
	return &ConfigurationSet{
		Name:                 name,
		SendingOptions:       &SendingOptions{SendingEnabled: true},
		ReputationOptions:    &ReputationOptions{ReputationMetricsEnabled: true},
		CreatedTimestamp:     now,
		LastUpdatedTimestamp: now,
	}
}

// CreateConfigurationSet creates a new configuration set.
func (s *SESv2Store) CreateConfigurationSet(configSet *ConfigurationSet) (*ConfigurationSet, error) {
	if s.configSetStore.Exists(configSet.Name) {
		return nil, ErrConfigSetAlreadyExists
	}

	configSet.CreatedTimestamp = time.Now().UTC()
	configSet.LastUpdatedTimestamp = configSet.CreatedTimestamp

	if err := s.configSetStore.Put(configSet.Name, configSet); err != nil {
		return nil, err
	}

	return configSet, nil
}

// GetConfigurationSet retrieves a configuration set by name.
func (s *SESv2Store) GetConfigurationSet(name string) (*ConfigurationSet, error) {
	var configSet ConfigurationSet
	if err := s.configSetStore.Get(name, &configSet); err != nil {
		return nil, ErrConfigSetNotFound
	}
	return &configSet, nil
}

// DeleteConfigurationSet deletes a configuration set by name.
func (s *SESv2Store) DeleteConfigurationSet(name string) error {
	if !s.configSetStore.Exists(name) {
		return ErrConfigSetNotFound
	}
	arn := s.BuildConfigSetArn(name)
	if err := s.TagStore.Delete(arn); err != nil {
		logs.Error("Failed to delete tags for config set", logs.String("name", name), logs.Err(err))
	}
	return s.configSetStore.Delete(name)
}

// UpdateConfigurationSet updates an existing configuration set.
func (s *SESv2Store) UpdateConfigurationSet(configSet *ConfigurationSet) error {
	if !s.configSetStore.Exists(configSet.Name) {
		return ErrConfigSetNotFound
	}
	configSet.LastUpdatedTimestamp = time.Now().UTC()
	return s.configSetStore.Put(configSet.Name, configSet)
}

// ListConfigurationSets lists configuration sets.
func (s *SESv2Store) ListConfigurationSets(opts common.ListOptions) (*common.ListResult[ConfigurationSet], error) {
	return common.List[ConfigurationSet](s.configSetStore, opts, nil)
}

// ConfigurationSetExists checks if a configuration set exists.
func (s *SESv2Store) ConfigurationSetExists(name string) bool {
	return s.configSetStore.Exists(name)
}
