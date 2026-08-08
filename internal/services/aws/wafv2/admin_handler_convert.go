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
