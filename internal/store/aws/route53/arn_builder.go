package route53

// Package route53 provides Route 53 data store implementations for vorpalstacks.

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ARNBuilder wraps svcarn.Route53Builder for backward compatibility
type ARNBuilder struct {
	builder *svcarn.Route53Builder
}

// NewARNBuilder creates a new ARNBuilder instance with the specified account ID.
func NewARNBuilder(accountId string) *ARNBuilder {
	return &ARNBuilder{
		builder: svcarn.NewARNBuilder(accountId, "").Route53(),
	}
}

// BuildHostedZoneARN builds an ARN for a Route53 hosted zone.
func (b *ARNBuilder) BuildHostedZoneARN(zoneId string) string {
	return b.builder.HostedZone(zoneId)
}

// BuildHealthCheckARN builds an ARN for a Route53 health check.
func (b *ARNBuilder) BuildHealthCheckARN(healthCheckId string) string {
	return b.builder.HealthCheck(healthCheckId)
}
