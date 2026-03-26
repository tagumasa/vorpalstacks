package s3

import (
	pb "vorpalstacks/internal/pb/storage/storage_s3"
)

func websiteConfigurationToProto(w *WebsiteConfiguration) *pb.WebsiteConfiguration {
	if w == nil {
		return nil
	}
	return &pb.WebsiteConfiguration{
		IndexDocument:         w.IndexDocument,
		ErrorDocument:         w.ErrorDocument,
		RedirectAllRequestsTo: redirectAllRequestsToToProto(w.RedirectAllRequestsTo),
		RoutingRules:          routingRulesToProto(w.RoutingRules),
	}
}

func protoToWebsiteConfiguration(p *pb.WebsiteConfiguration) *WebsiteConfiguration {
	if p == nil {
		return nil
	}
	return &WebsiteConfiguration{
		IndexDocument:         p.IndexDocument,
		ErrorDocument:         p.ErrorDocument,
		RedirectAllRequestsTo: protoToRedirectAllRequestsTo(p.RedirectAllRequestsTo),
		RoutingRules:          protoToRoutingRules(p.RoutingRules),
	}
}

func redirectAllRequestsToToProto(r *RedirectAllRequestsTo) *pb.RedirectAllRequestsTo {
	if r == nil {
		return nil
	}
	return &pb.RedirectAllRequestsTo{
		HostName: r.HostName,
		Protocol: r.Protocol,
	}
}

func protoToRedirectAllRequestsTo(p *pb.RedirectAllRequestsTo) *RedirectAllRequestsTo {
	if p == nil {
		return nil
	}
	return &RedirectAllRequestsTo{
		HostName: p.HostName,
		Protocol: p.Protocol,
	}
}

func routingRulesToProto(rules []RoutingRule) []*pb.RoutingRule {
	if rules == nil {
		return nil
	}
	result := make([]*pb.RoutingRule, len(rules))
	for i, r := range rules {
		result[i] = routingRuleToProto(&r)
	}
	return result
}

func protoToRoutingRules(p []*pb.RoutingRule) []RoutingRule {
	if p == nil {
		return nil
	}
	result := make([]RoutingRule, len(p))
	for i, r := range p {
		result[i] = *protoToRoutingRule(r)
	}
	return result
}

func routingRuleToProto(r *RoutingRule) *pb.RoutingRule {
	if r == nil {
		return nil
	}
	return &pb.RoutingRule{
		Condition: routingRuleConditionToProto(r.Condition),
		Redirect:  routingRuleRedirectToProto(r.Redirect),
	}
}

func protoToRoutingRule(p *pb.RoutingRule) *RoutingRule {
	if p == nil {
		return nil
	}
	return &RoutingRule{
		Condition: protoToRoutingRuleCondition(p.Condition),
		Redirect:  protoToRoutingRuleRedirect(p.Redirect),
	}
}

func routingRuleConditionToProto(c *RoutingRuleCondition) *pb.RoutingRuleCondition {
	if c == nil {
		return nil
	}
	return &pb.RoutingRuleCondition{
		HttpErrorCodeReturnedEquals: stringPtrToString(c.HTTPErrorCodeReturnedEquals),
		KeyPrefixEquals:             stringPtrToString(c.KeyPrefixEquals),
	}
}

func protoToRoutingRuleCondition(p *pb.RoutingRuleCondition) *RoutingRuleCondition {
	if p == nil {
		return nil
	}
	return &RoutingRuleCondition{
		HTTPErrorCodeReturnedEquals: stringToStringPtr(p.HttpErrorCodeReturnedEquals),
		KeyPrefixEquals:             stringToStringPtr(p.KeyPrefixEquals),
	}
}

func routingRuleRedirectToProto(r *RoutingRuleRedirect) *pb.RoutingRuleRedirect {
	if r == nil {
		return nil
	}
	return &pb.RoutingRuleRedirect{
		HostName:             stringPtrToString(r.HostName),
		HttpRedirectCode:     stringPtrToString(r.HTTPRedirectCode),
		Protocol:             stringPtrToString(r.Protocol),
		ReplaceKeyPrefixWith: stringPtrToString(r.ReplaceKeyPrefixWith),
		ReplaceKeyWith:       stringPtrToString(r.ReplaceKeyWith),
	}
}

func protoToRoutingRuleRedirect(p *pb.RoutingRuleRedirect) *RoutingRuleRedirect {
	if p == nil {
		return nil
	}
	return &RoutingRuleRedirect{
		HostName:             stringToStringPtr(p.HostName),
		HTTPRedirectCode:     stringToStringPtr(p.HttpRedirectCode),
		Protocol:             stringToStringPtr(p.Protocol),
		ReplaceKeyPrefixWith: stringToStringPtr(p.ReplaceKeyPrefixWith),
		ReplaceKeyWith:       stringToStringPtr(p.ReplaceKeyWith),
	}
}

func corsConfigurationToProto(c *CORSConfiguration) *pb.CORSConfiguration {
	if c == nil {
		return nil
	}
	return &pb.CORSConfiguration{
		CorsRules: corsRulesToProto(c.CORSRules),
	}
}

func protoToCORSConfiguration(p *pb.CORSConfiguration) *CORSConfiguration {
	if p == nil {
		return nil
	}
	return &CORSConfiguration{
		CORSRules: protoToCORSRules(p.CorsRules),
	}
}

func corsRulesToProto(rules []CORSRule) []*pb.CORSRule {
	if rules == nil {
		return nil
	}
	result := make([]*pb.CORSRule, len(rules))
	for i, r := range rules {
		result[i] = corsRuleToProto(&r)
	}
	return result
}

func protoToCORSRules(p []*pb.CORSRule) []CORSRule {
	if p == nil {
		return nil
	}
	result := make([]CORSRule, len(p))
	for i, r := range p {
		result[i] = *protoToCORSRule(r)
	}
	return result
}

func corsRuleToProto(r *CORSRule) *pb.CORSRule {
	if r == nil {
		return nil
	}
	return &pb.CORSRule{
		AllowedHeaders: r.AllowedHeaders,
		AllowedMethods: r.AllowedMethods,
		AllowedOrigins: r.AllowedOrigins,
		ExposeHeaders:  r.ExposeHeaders,
		MaxAgeSeconds:  intPtrToInt32(r.MaxAgeSeconds),
		Id:             r.ID,
	}
}

func protoToCORSRule(p *pb.CORSRule) *CORSRule {
	if p == nil {
		return nil
	}
	return &CORSRule{
		AllowedHeaders: p.AllowedHeaders,
		AllowedMethods: p.AllowedMethods,
		AllowedOrigins: p.AllowedOrigins,
		ExposeHeaders:  p.ExposeHeaders,
		MaxAgeSeconds:  int32ToIntPtr(p.MaxAgeSeconds),
		ID:             p.Id,
	}
}
