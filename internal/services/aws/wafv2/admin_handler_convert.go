package wafv2

import (
	pb "vorpalstacks/internal/pb/aws/wafv2"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// toPbWebACLSummary converts a store-layer WebACL to the protobuf
// WebACLSummary used in admin console responses. This file is the sole
// location in the admin handler layer that imports the store package.
func toPbWebACLSummary(wa *wafstore.WebACL) *pb.WebACLSummary {
	return &pb.WebACLSummary{
		Id:          wa.ID,
		Name:        wa.Name,
		Arn:         wa.ARN,
		Description: wa.Description,
		Locktoken:   wa.LockToken,
	}
}

// defaultCreateWebACLInput builds the transport-agnostic create input
// for the admin console. CreateWebACLRequest marks DefaultAction and
// VisibilityConfig as required, but the console UI does not collect
// them, so this adapter synthesises safe defaults (allow-all with
// metrics enabled; the metric name derives from the ACL name, which
// satisfies the MetricName pattern because EntityName is a subset of
// it). Keeping the synthesis here preserves the rule that only this
// file references store types.
func defaultCreateWebACLInput(name, description, scope string) CreateWebACLInput {
	return CreateWebACLInput{
		Name:        name,
		Description: description,
		Scope:       scope,
		DefaultAction: &wafstore.Action{
			Allow: &wafstore.AllowAction{},
		},
		VisibilityConfig: &wafstore.VisibilityConfig{
			SampledRequestsEnabled:   true,
			CloudWatchMetricsEnabled: true,
			MetricName:               name,
		},
	}
}
