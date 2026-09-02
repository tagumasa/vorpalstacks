package waf

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ARNBuilder wraps svcarn.WAFBuilder for ARN construction
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
func (b *ARNBuilder) BuildWebACLARN(name, id, scope string) string {
	return b.builder.WebACL(name, id, scope)
}

// BuildRuleGroupARN builds an ARN for a WAF Rule Group.
func (b *ARNBuilder) BuildRuleGroupARN(name, id, scope string) string {
	if scope == "" {
		scope = "REGIONAL"
	}
	return b.builder.RuleGroup(name, id, scope)
}

// BuildIPSetARN builds an ARN for a WAF IP Set.
func (b *ARNBuilder) BuildIPSetARN(name, id, scope string) string {
	if scope == "" {
		scope = "REGIONAL"
	}
	return b.builder.IPSet(name, id, scope)
}

// BuildRegexPatternSetARN builds an ARN for a WAF Regex Pattern Set.
func (b *ARNBuilder) BuildRegexPatternSetARN(name, id, scope string) string {
	if scope == "" {
		scope = "REGIONAL"
	}
	return b.builder.RegexPatternSet(name, id, scope)
}

// BuildRuleARN builds an ARN for a WAF Classic Rule.
func (b *ARNBuilder) BuildRuleARN(id string) string {
	return b.builder.Build("waf", "rule/"+id)
}
