// Package ssm provides Systems Manager Parameter Store storage functionality for vorpalstacks.
package ssm

import (
	"time"
)

// ParameterType defines the type of a parameter in Parameter Store.
type ParameterType string

// ParameterType constants define the supported parameter types.
const (
	ParameterTypeString       ParameterType = "String"
	ParameterTypeStringList   ParameterType = "StringList"
	ParameterTypeSecureString ParameterType = "SecureString"
)

// ParameterTier defines the tier level of a parameter.
type ParameterTier string

// ParameterTier constants define the supported parameter tiers.
const (
	ParameterTierStandard           ParameterTier = "Standard"
	ParameterTierAdvanced           ParameterTier = "Advanced"
	ParameterTierIntelligentTiering ParameterTier = "Intelligent-Tiering"
)

// Parameter represents a Systems Manager Parameter Store parameter.
type Parameter struct {
	Name             string             `json:"name"`
	Value            string             `json:"value"`
	Type             ParameterType      `json:"type"`
	Tier             ParameterTier      `json:"tier"`
	Description      string             `json:"description,omitempty"`
	KeyID            string             `json:"keyId,omitempty"`
	Version          int64              `json:"version"`
	LastModifiedDate time.Time          `json:"lastModifiedDate"`
	DataType         string             `json:"dataType,omitempty"`
	AllowedPattern   string             `json:"allowedPattern,omitempty"`
	Policies         string             `json:"policies,omitempty"`
	ARN              string             `json:"arn"`
	Tags             map[string]string  `json:"tags,omitempty"`
	VersionLabels    map[int64][]string `json:"versionLabels,omitempty"`
}

// ParameterVersion represents a specific version of a parameter.
type ParameterVersion struct {
	ParameterName    string        `json:"parameterName"`
	Version          int64         `json:"version"`
	Value            string        `json:"value"`
	Type             ParameterType `json:"type"`
	LastModifiedDate time.Time     `json:"lastModifiedDate"`
	DataType         string        `json:"dataType,omitempty"`
	Labels           []string      `json:"labels,omitempty"`
}

// ParameterMetadata represents the metadata of a parameter without its value.
type ParameterMetadata struct {
	Name             string        `json:"name"`
	Type             ParameterType `json:"type"`
	Tier             ParameterTier `json:"tier"`
	Description      string        `json:"description,omitempty"`
	KeyID            string        `json:"keyId,omitempty"`
	Version          int64         `json:"version"`
	LastModifiedDate time.Time     `json:"lastModifiedDate"`
	DataType         string        `json:"dataType,omitempty"`
	AllowedPattern   string        `json:"allowedPattern,omitempty"`
	Policies         string        `json:"policies,omitempty"`
	ARN              string        `json:"arn"`
}

// Tag represents a tag for an SSM parameter.
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewParameter creates a new Parameter with the specified name, value, and type.
func NewParameter(name, value string, paramType ParameterType) *Parameter {
	return &Parameter{
		Name:             name,
		Value:            value,
		Type:             paramType,
		Tier:             ParameterTierStandard,
		Version:          1,
		LastModifiedDate: time.Now().UTC(),
		DataType:         "text",
		Tags:             make(map[string]string),
	}
}

// NewParameterMetadata creates ParameterMetadata from a Parameter.
func NewParameterMetadata(param *Parameter) *ParameterMetadata {
	return &ParameterMetadata{
		Name:             param.Name,
		Type:             param.Type,
		Tier:             param.Tier,
		Description:      param.Description,
		KeyID:            param.KeyID,
		Version:          param.Version,
		LastModifiedDate: param.LastModifiedDate,
		DataType:         param.DataType,
		AllowedPattern:   param.AllowedPattern,
		Policies:         param.Policies,
		ARN:              param.ARN,
	}
}
