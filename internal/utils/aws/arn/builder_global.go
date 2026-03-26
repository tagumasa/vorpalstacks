// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

import "strings"

// Route53Builder provides methods for constructing Route 53 ARNs.
type Route53Builder struct{ *ARNBuilder }

// Route53 returns a Route53Builder for constructing Route 53 ARNs.
func (b *ARNBuilder) Route53() *Route53Builder { return &Route53Builder{b} }

// HostedZone constructs an ARN for a Route 53 hosted zone.
func (b *Route53Builder) HostedZone(id string) string {
	return b.BuildGlobal("route53", "hostedzone/"+id)
}

// HealthCheck constructs an ARN for a Route 53 health check.
func (b *Route53Builder) HealthCheck(id string) string {
	return b.BuildGlobal("route53", "healthcheck/"+id)
}

// CloudFrontBuilder provides methods for constructing CloudFront ARNs.
type CloudFrontBuilder struct{ *ARNBuilder }

// CloudFront returns a CloudFrontBuilder for constructing CloudFront ARNs.
func (b *ARNBuilder) CloudFront() *CloudFrontBuilder { return &CloudFrontBuilder{b} }

// Distribution constructs an ARN for a CloudFront distribution.
func (b *CloudFrontBuilder) Distribution(id string) string {
	return b.BuildGlobal("cloudfront", "distribution/"+id)
}

// CachePolicy constructs an ARN for a CloudFront cache policy.
func (b *CloudFrontBuilder) CachePolicy(id string) string {
	return b.BuildGlobal("cloudfront", "cache-policy/"+id)
}

// OriginRequestPolicy constructs an ARN for a CloudFront origin request policy.
func (b *CloudFrontBuilder) OriginRequestPolicy(id string) string {
	return b.BuildGlobal("cloudfront", "origin-request-policy/"+id)
}

// OriginAccessControl constructs an ARN for a CloudFront origin access control.
func (b *CloudFrontBuilder) OriginAccessControl(id string) string {
	return b.BuildGlobal("cloudfront", "origin-access-control/"+id)
}

// ResponseHeadersPolicy constructs an ARN for a CloudFront response headers policy.
func (b *CloudFrontBuilder) ResponseHeadersPolicy(id string) string {
	return b.BuildGlobal("cloudfront", "response-headers-policy/"+id)
}

// WAFBuilder provides methods for constructing WAF ARNs.
type WAFBuilder struct{ *ARNBuilder }

// WAF returns a WAFBuilder for constructing WAF ARNs.
func (b *ARNBuilder) WAF() *WAFBuilder { return &WAFBuilder{b} }

// WebACL constructs an ARN for a WAF Web ACL.
func (b *WAFBuilder) WebACL(id, scope string) string {
	return b.Build("waf", "webacl/"+scope+"/"+id)
}

// RuleGroup constructs an ARN for a WAF rule group.
func (b *WAFBuilder) RuleGroup(id, scope string) string {
	return b.Build("waf", "rulegroup/"+scope+"/"+id)
}

// IPSet constructs an ARN for a WAF IP set.
func (b *WAFBuilder) IPSet(id, scope string) string {
	return b.Build("waf", "ipset/"+scope+"/"+id)
}

// RegexPatternSet constructs a WAF regex pattern set ARN.
func (b *WAFBuilder) RegexPatternSet(id, scope string) string {
	return b.Build("waf", "regexpatternset/"+scope+"/"+id)
}

// TimestreamBuilder provides methods for constructing Timestream ARNs.
type TimestreamBuilder struct{ *ARNBuilder }

// Timestream returns a TimestreamBuilder for constructing Timestream ARNs.
func (b *ARNBuilder) Timestream() *TimestreamBuilder { return &TimestreamBuilder{b} }

// Database constructs an ARN for a Timestream database.
func (b *TimestreamBuilder) Database(name string) string {
	return b.Build("timestream", "database/"+name)
}

// Table constructs an ARN for a Timestream table.
func (b *TimestreamBuilder) Table(db, table string) string {
	return b.Build("timestream", "database/"+db+"/table/"+table)
}

// ScheduledQuery constructs an ARN for a Timestream scheduled query.
func (b *TimestreamBuilder) ScheduledQuery(name string) string {
	return b.Build("timestream", "scheduled-query/"+name)
}

// ParseDatabaseName extracts the database name from a Timestream database ARN.
func (b *TimestreamBuilder) ParseDatabaseName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "database/") {
		parts := strings.Split(strings.TrimPrefix(resource, "database/"), "/")
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

// ParseTableName extracts the table name from a Timestream table ARN.
func (b *TimestreamBuilder) ParseTableName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "database/") {
		parts := strings.Split(strings.TrimPrefix(resource, "database/"), "/")
		if len(parts) >= 3 && parts[1] == "table" {
			return parts[2]
		}
	}
	return ""
}

// ParseScheduledQueryName extracts the scheduled query name from a Timestream scheduled query ARN.
func (b *TimestreamBuilder) ParseScheduledQueryName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "scheduled-query/") {
		return strings.TrimPrefix(resource, "scheduled-query/")
	}
	return ""
}

// AthenaBuilder provides methods for constructing Athena ARNs.
type AthenaBuilder struct{ *ARNBuilder }

// Athena returns an AthenaBuilder for constructing Athena ARNs.
func (b *ARNBuilder) Athena() *AthenaBuilder { return &AthenaBuilder{b} }

// WorkGroup constructs an ARN for an Athena work group.
func (b *AthenaBuilder) WorkGroup(name string) string { return b.Build("athena", "workgroup/"+name) }

// DataCatalog constructs an ARN for an Athena data catalog.
func (b *AthenaBuilder) DataCatalog(name string) string {
	return b.Build("athena", "datacatalog/"+name)
}

// ParseWorkGroupName extracts the work group name from an Athena work group ARN.
func (b *AthenaBuilder) ParseWorkGroupName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "workgroup/") {
		return strings.TrimPrefix(resource, "workgroup/")
	}
	return ""
}
