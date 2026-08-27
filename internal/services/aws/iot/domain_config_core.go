package iot

import (
	"encoding/json"
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Domain configuration Core
// ---------------------------------------------------------------------------

// DomainConfigAuthorizerConfig is the wire authorizer-config structure
// (the Smithy AuthorizerConfig shape).
type DomainConfigAuthorizerConfig struct {
	DefaultAuthorizerName   string
	AllowAuthorizerOverride bool
	OverrideProvided        bool
}

// CreateDomainConfigurationInput carries the fields for
// CreateDomainConfiguration.
type CreateDomainConfigurationInput struct {
	DomainConfigurationName  string
	DomainName               string
	ServerCertificateARNs    []string
	ValidationCertificateARN string
	AuthorizerConfig         DomainConfigAuthorizerConfig
	AuthorizerConfigProvided bool
	ServiceType              string
	AuthenticationType       string
	ApplicationProtocol      string
	Tags                     map[string]string
}

// CreateDomainConfigurationResult is the transport-agnostic result of
// CreateDomainConfiguration.
type CreateDomainConfigurationResult struct {
	DomainConfigurationName string
	DomainConfigurationARN  string
}

// DescribeDomainConfigurationInput carries the fields for
// DescribeDomainConfiguration. An empty name resolves to the reserved
// "default" configuration.
type DescribeDomainConfigurationInput struct {
	DomainConfigurationName string
}

// UpdateDomainConfigurationInput carries the fields for
// UpdateDomainConfiguration. Only the model's update members travel here;
// AuthorizerConfig applies when provided, RemoveAuthorizerConfig clears it
// (and wins when both are supplied).
type UpdateDomainConfigurationInput struct {
	DomainConfigurationName   string
	AuthorizerConfig          DomainConfigAuthorizerConfig
	AuthorizerConfigProvided  bool
	DomainConfigurationStatus string
	RemoveAuthorizerConfig    bool
	AuthenticationType        string
	ApplicationProtocol       string
}

// ListDomainConfigurationsResult is the transport-agnostic result of
// ListDomainConfigurations.
type ListDomainConfigurationsResult struct {
	DomainConfigurations []*iotstore.DomainConfiguration
	NextMarker           string
}

// The Smithy enum value sets for the domain-configuration members.
var (
	domainConfigServiceTypes = map[string]struct{}{
		"DATA": {}, "CREDENTIAL_PROVIDER": {}, "JOBS": {},
	}
	domainConfigStatuses = map[string]struct{}{
		"ENABLED": {}, "DISABLED": {},
	}
	domainConfigAuthTypes = map[string]struct{}{
		"CUSTOM_AUTH_X509": {}, "CUSTOM_AUTH": {}, "AWS_X509": {}, "AWS_SIGV4": {}, "DEFAULT": {},
	}
	domainConfigAppProtocols = map[string]struct{}{
		"SECURE_MQTT": {}, "MQTT_WSS": {}, "HTTPS": {}, "DEFAULT": {},
	}
)

// validateDomainConfigEnum rejects a non-empty value that is outside its
// enum value set.
func validateDomainConfigEnum(value string, allowed map[string]struct{}) error {
	if value == "" {
		return nil
	}
	if _, ok := allowed[value]; !ok {
		return iotstore.ErrInvalidRequest
	}
	return nil
}

// domainConfigAuthorizerConfigJSON serialises the authorizer-config
// structure for the store's string field.
func domainConfigAuthorizerConfigJSON(cfg DomainConfigAuthorizerConfig) string {
	payload := map[string]interface{}{
		"defaultAuthorizerName": cfg.DefaultAuthorizerName,
	}
	if cfg.OverrideProvided {
		payload["allowAuthorizerOverride"] = cfg.AllowAuthorizerOverride
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(b)
}

// domainConfigAuthorizerConfigValue decodes the stored authorizer-config
// JSON into the response structure; nil when unset or unparseable.
func domainConfigAuthorizerConfigValue(stored string) map[string]interface{} {
	if stored == "" {
		return nil
	}
	var cfg struct {
		DefaultAuthorizerName   string `json:"defaultAuthorizerName"`
		AllowAuthorizerOverride *bool  `json:"allowAuthorizerOverride"`
	}
	if err := json.Unmarshal([]byte(stored), &cfg); err != nil {
		return nil
	}
	out := map[string]interface{}{
		"defaultAuthorizerName": cfg.DefaultAuthorizerName,
	}
	if cfg.AllowAuthorizerOverride != nil {
		out["allowAuthorizerOverride"] = *cfg.AllowAuthorizerOverride
	}
	return out
}

// createDomainConfigurationCore validates and persists a domain
// configuration; duplicate names are rejected.
func (s *IoTService) createDomainConfigurationCore(store iotstore.IotStoreInterface, in CreateDomainConfigurationInput) (*CreateDomainConfigurationResult, error) {
	if in.DomainConfigurationName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if err := validateDomainConfigEnum(in.ServiceType, domainConfigServiceTypes); err != nil {
		return nil, err
	}
	if err := validateDomainConfigEnum(in.AuthenticationType, domainConfigAuthTypes); err != nil {
		return nil, err
	}
	if err := validateDomainConfigEnum(in.ApplicationProtocol, domainConfigAppProtocols); err != nil {
		return nil, err
	}
	// The model caps the server-certificate list at one ARN.
	if len(in.ServerCertificateARNs) > 1 {
		return nil, iotstore.ErrInvalidRequest
	}

	if _, err := store.GetDomainConfiguration(in.DomainConfigurationName); err == nil {
		return nil, iotstore.ErrDomainConfigurationAlreadyExists
	}

	now := time.Now().UTC()
	dc := &iotstore.DomainConfiguration{
		DomainConfigurationName:   in.DomainConfigurationName,
		DomainName:                in.DomainName,
		ServerCertificateARNs:     in.ServerCertificateARNs,
		ValidationCertificateARN:  in.ValidationCertificateARN,
		ServiceType:               in.ServiceType,
		AuthenticationType:        in.AuthenticationType,
		ApplicationProtocol:       in.ApplicationProtocol,
		Tags:                      in.Tags,
		DomainConfigurationStatus: "ENABLED",
		CreationDate:              now,
		LastModifiedDate:          now,
	}
	if in.AuthorizerConfigProvided {
		dc.AuthorizerConfig = domainConfigAuthorizerConfigJSON(in.AuthorizerConfig)
	}

	created, err := store.CreateDomainConfiguration(dc)
	if err != nil {
		return nil, err
	}

	// Create-time tags live on the ARN-keyed tag store so
	// ListTagsForResource can see them.
	arn := iotstore.BuildDomainConfigurationARN(store.GetAccountID(), store.GetRegion(), created.DomainConfigurationName)
	if len(in.Tags) > 0 {
		if err := store.TagResource(arn, in.Tags); err != nil {
			return nil, err
		}
	}

	return &CreateDomainConfigurationResult{
		DomainConfigurationName: created.DomainConfigurationName,
		DomainConfigurationARN:  created.DomainConfigurationARN,
	}, nil
}

// describeDomainConfigurationCore retrieves a domain configuration by name,
// resolving an omitted name to the reserved "default" configuration.
func (s *IoTService) describeDomainConfigurationCore(store iotstore.IotStoreInterface, in DescribeDomainConfigurationInput) (*iotstore.DomainConfiguration, error) {
	name := in.DomainConfigurationName
	if name == "" {
		name = "default"
	}
	return store.GetDomainConfiguration(name)
}

// updateDomainConfigurationCore applies the model's update members to an
// existing domain configuration: the authorizer config (or its removal) and
// the status.
func (s *IoTService) updateDomainConfigurationCore(store iotstore.IotStoreInterface, in UpdateDomainConfigurationInput) (*iotstore.DomainConfiguration, error) {
	if in.DomainConfigurationName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if err := validateDomainConfigEnum(in.DomainConfigurationStatus, domainConfigStatuses); err != nil {
		return nil, err
	}
	if err := validateDomainConfigEnum(in.AuthenticationType, domainConfigAuthTypes); err != nil {
		return nil, err
	}
	if err := validateDomainConfigEnum(in.ApplicationProtocol, domainConfigAppProtocols); err != nil {
		return nil, err
	}

	dc, err := store.GetDomainConfiguration(in.DomainConfigurationName)
	if err != nil {
		return nil, err
	}

	if in.AuthorizerConfigProvided {
		dc.AuthorizerConfig = domainConfigAuthorizerConfigJSON(in.AuthorizerConfig)
	}
	if in.RemoveAuthorizerConfig {
		dc.AuthorizerConfig = ""
	}
	if in.DomainConfigurationStatus != "" {
		dc.DomainConfigurationStatus = in.DomainConfigurationStatus
	}
	if in.AuthenticationType != "" {
		dc.AuthenticationType = in.AuthenticationType
	}
	if in.ApplicationProtocol != "" {
		dc.ApplicationProtocol = in.ApplicationProtocol
	}
	dc.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateDomainConfiguration(in.DomainConfigurationName, dc); err != nil {
		return nil, err
	}
	return dc, nil
}

// deleteDomainConfigurationCore removes a domain configuration and its tags.
func (s *IoTService) deleteDomainConfigurationCore(store iotstore.IotStoreInterface, name string) error {
	if name == "" {
		return iotstore.ErrMissingParam
	}

	arn := iotstore.BuildDomainConfigurationARN(store.GetAccountID(), store.GetRegion(), name)
	_ = store.DeleteAllTags(arn)

	return store.DeleteDomainConfiguration(name)
}

// listDomainConfigurationsCore lists domain configurations with pagination.
func (s *IoTService) listDomainConfigurationsCore(store iotstore.IotStoreInterface, marker string, maxItems int) (*ListDomainConfigurationsResult, error) {
	maxItems = listMaxItems(maxItems)
	dcs, err := store.ListDomainConfigurations(iotstoreListOpts(maxItems, marker))
	if err != nil {
		return nil, err
	}
	return &ListDomainConfigurationsResult{
		DomainConfigurations: dcs.Items,
		NextMarker:           dcs.NextMarker,
	}, nil
}
