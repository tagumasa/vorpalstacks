package waf

// Package waf provides WAF (Web Application Firewall) data store implementations
// for vorpalstacks.

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ARNBuilder wraps svcarn.WAFBuilder for backward compatibility
type ARNBuilder struct {
	builder *svcarn.WAFBuilder
}

// NewARNBuilder creates a new WAF ARN builder.
func NewARNBuilder(accountId, region string) *ARNBuilder {
	return &ARNBuilder{
		builder: svcarn.NewARNBuilder(accountId, region).WAF(),
	}
}

// BuildWebACLARN builds an ARN for a WAF Web ACL.
func (b *ARNBuilder) BuildWebACLARN(id, scope string) string {
	return b.builder.WebACL(id, scope)
}

// BuildRuleGroupARN builds an ARN for a WAF Rule Group.
func (b *ARNBuilder) BuildRuleGroupARN(id string) string {
	return b.builder.RuleGroup(id, "REGIONAL")
}

// BuildIPSetARN builds an ARN for a WAF IP Set.
func (b *ARNBuilder) BuildIPSetARN(id string) string {
	return b.builder.IPSet(id, "REGIONAL")
}

// BuildRegexPatternSetARN builds an ARN for a WAF Regex Pattern Set.
func (b *ARNBuilder) BuildRegexPatternSetARN(id string) string {
	return b.builder.RegexPatternSet(id, "REGIONAL")
}

// BuildRuleARN builds an ARN for a WAF Classic Rule.
func (b *ARNBuilder) BuildRuleARN(id string) string {
	return b.builder.Build("waf", "rule/"+id)
}
