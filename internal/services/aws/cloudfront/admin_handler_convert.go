package cloudfront

import (
	"google.golang.org/protobuf/proto"

	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/cloudfront"
)

// This file is the sole exception to the store-import prohibition: it is the only admin
// handler file that imports the store package. It contains only pure proto
// conversion helpers (toPb* functions) that translate store types to proto
// types for response marshalling.

// toPbDistribution converts a store Distribution to a proto Distribution.
func toPbDistribution(d *cloudfrontstore.Distribution) *pb.Distribution {
	return &pb.Distribution{
		Id:               d.ID,
		Arn:              d.ARN,
		Status:           d.Status,
		Domainname:       d.DomainName,
		Lastmodifiedtime: d.LastModifiedAt.Format(timeutils.ISO8601UTCFormat),
	}
}

// toPbDistributionSummary converts a store Distribution to a proto
// DistributionSummary for list responses.
func toPbDistributionSummary(d *cloudfrontstore.Distribution) *pb.DistributionSummary {
	summary := &pb.DistributionSummary{
		Id:         d.ID,
		Arn:        d.ARN,
		Status:     d.Status,
		Enabled:    proto.Bool(d.Enabled),
		Staging:    proto.Bool(d.Staging),
		Etag:       d.ETag,
		Domainname: d.DomainName,
	}
	if d.DistributionConfig != nil {
		summary.Comment = d.DistributionConfig.Comment
	}
	if !d.LastModifiedAt.IsZero() {
		summary.Lastmodifiedtime = d.LastModifiedAt.Format(timeutils.ISO8601UTCFormat)
	}
	return summary
}
