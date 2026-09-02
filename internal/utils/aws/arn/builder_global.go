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
	return b.BuildNoRegion("cloudfront", "distribution/"+id)
}

// CachePolicy constructs an ARN for a CloudFront cache policy.
func (b *CloudFrontBuilder) CachePolicy(id string) string {
	return b.BuildNoRegion("cloudfront", "cache-policy/"+id)
}

// OriginRequestPolicy constructs an ARN for a CloudFront origin request policy.
func (b *CloudFrontBuilder) OriginRequestPolicy(id string) string {
	return b.BuildNoRegion("cloudfront", "origin-request-policy/"+id)
}

// OriginAccessControl constructs an ARN for a CloudFront origin access control.
func (b *CloudFrontBuilder) OriginAccessControl(id string) string {
	return b.BuildNoRegion("cloudfront", "origin-access-control/"+id)
}

// ResponseHeadersPolicy constructs an ARN for a CloudFront response headers policy.
func (b *CloudFrontBuilder) ResponseHeadersPolicy(id string) string {
	return b.BuildNoRegion("cloudfront", "response-headers-policy/"+id)
}

// PublicKey constructs an ARN for a CloudFront public key.
func (b *CloudFrontBuilder) PublicKey(id string) string {
	return b.BuildNoRegion("cloudfront", "public-key/"+id)
}

// KeyGroup constructs an ARN for a CloudFront key group.
func (b *CloudFrontBuilder) KeyGroup(id string) string {
	return b.BuildNoRegion("cloudfront", "key-group/"+id)
}

// WAFBuilder provides methods for constructing WAF ARNs.
type WAFBuilder struct{ *ARNBuilder }

// WAF returns a WAFBuilder for constructing WAF ARNs.
func (b *ARNBuilder) WAF() *WAFBuilder { return &WAFBuilder{b} }

// wafScopeSegment maps a WAF scope to the ARN resource-path segment:
// CLOUDFRONT-scoped resources live under "global", everything else
// under "regional" (the IAM documentation formats wafv2 ARNs as
// arn:<partition>:wafv2:<region>:<account>:<scope>/<type>/<name>/<id>).
func wafScopeSegment(scope string) string {
	if strings.EqualFold(scope, "CLOUDFRONT") || strings.EqualFold(scope, "global") {
		return "global"
	}
	return "regional"
}

// WAFCloudFrontRegion is the region every CLOUDFRONT-scoped wafv2
// resource belongs to: the IAM documentation's guidance for resources
// that protect Amazon CloudFront distributions is to set the region to
// us-east-1, independent of the caller's configured region, and the
// WAFv2 API reference requires the us-east-1 endpoint for the CloudFront
// scope. Every layer that must locate a global-scope resource — ARN
// construction, storage routing, lookups — references this constant.
const WAFCloudFrontRegion = "us-east-1"

// wafARN builds a wafv2 resource ARN under the given scope; global
// (CloudFront) scope fixes the region field to us-east-1.
func (b *WAFBuilder) wafARN(scope, resource string) string {
	segment := wafScopeSegment(scope)
	if segment == "global" {
		return b.BuildInRegion("wafv2", WAFCloudFrontRegion, segment+"/"+resource)
	}
	return b.Build("wafv2", segment+"/"+resource)
}

// WebACL constructs an ARN for a WAF Web ACL.
func (b *WAFBuilder) WebACL(name, id, scope string) string {
	return b.wafARN(scope, "webacl/"+name+"/"+id)
}

// RuleGroup constructs an ARN for a WAF rule group.
func (b *WAFBuilder) RuleGroup(name, id, scope string) string {
	return b.wafARN(scope, "rulegroup/"+name+"/"+id)
}

// IPSet constructs an ARN for a WAF IP set.
func (b *WAFBuilder) IPSet(name, id, scope string) string {
	return b.wafARN(scope, "ipset/"+name+"/"+id)
}

// RegexPatternSet constructs a WAF regex pattern set ARN.
func (b *WAFBuilder) RegexPatternSet(name, id, scope string) string {
	return b.wafARN(scope, "regexpatternset/"+name+"/"+id)
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

// CapacityReservation constructs an ARN for an Athena capacity reservation.
func (b *AthenaBuilder) CapacityReservation(name string) string {
	return b.Build("athena", "capacityreservation/"+name)
}

// ParseWorkGroupName extracts the work group name from an Athena work group ARN.
func (b *AthenaBuilder) ParseWorkGroupName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "workgroup/") {
		return strings.TrimPrefix(resource, "workgroup/")
	}
	return ""
}
