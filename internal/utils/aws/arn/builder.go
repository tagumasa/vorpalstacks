// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

import "strings"

// ARNBuilder is a builder for constructing Amazon Resource Names.
type ARNBuilder struct {
	accountID string
	region    string
	partition string
}

// NewARNBuilder creates a new ARNBuilder with the specified account ID and region.
func NewARNBuilder(accountID, region string) *ARNBuilder {
	partition := PartitionAWS
	if strings.HasPrefix(region, "cn-") {
		partition = PartitionAWSCN
	}
	return &ARNBuilder{
		accountID: accountID,
		region:    region,
		partition: partition,
	}
}

// AccountId returns the account ID associated with the builder.
func (b *ARNBuilder) AccountId() string { return b.accountID }

// Region returns the region associated with the builder.
func (b *ARNBuilder) Region() string { return b.region }

// Partition returns the partition associated with the builder.
func (b *ARNBuilder) Partition() string { return b.partition }

// Build constructs an ARN for the given service and resource.
func (b *ARNBuilder) Build(service, resource string) string {
	return "arn:" + b.partition + ":" + service + ":" + b.region + ":" + b.accountID + ":" + resource
}

// BuildGlobal constructs a global ARN (without region) for the given service and resource.
func (b *ARNBuilder) BuildGlobal(service, resource string) string {
	return "arn:" + b.partition + ":" + service + ":::" + resource
}
