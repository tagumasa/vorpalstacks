package cloudwatchlogs

import (
	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/common/pbutil"
	pb "vorpalstacks/internal/pb/aws/cloudwatchlogs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// toPbLogGroupSummary converts a store-level LogGroup to the proto
// representation used by the admin console.
func toPbLogGroupSummary(lg *logsstore.LogGroup) *pb.LogGroupSummary {
	summary := &pb.LogGroupSummary{
		Loggroupname: proto.String(lg.Name),
		Loggrouparn:  proto.String(lg.ARN),
	}
	switch lg.LogGroupClass {
	case "DELIVERY":
		summary.Loggroupclass = pbutil.Enum(pb.LogGroupClass_LOG_GROUP_CLASS_DELIVERY)
	case "INFREQUENT_ACCESS":
		summary.Loggroupclass = pbutil.Enum(pb.LogGroupClass_LOG_GROUP_CLASS_INFREQUENT_ACCESS)
	default:
		summary.Loggroupclass = pbutil.Enum(pb.LogGroupClass_LOG_GROUP_CLASS_STANDARD)
	}
	return summary
}

// toPbLogStream converts a store-level LogStream to the proto
// representation used by the admin console.
func toPbLogStream(ls *logsstore.LogStream) *pb.LogStream {
	return &pb.LogStream{
		Logstreamname: proto.String(ls.Name),
		Arn:           proto.String(ls.ARN),
		Creationtime:  proto.Int64(ls.CreatedAt.UnixMilli()),
		// storedBytes is always zero for log streams (the LogStream
		// shape's deprecated member documentation).
		Storedbytes:         proto.Int64(0),
		Firsteventtimestamp: proto.Int64(ls.FirstEventTs),
		Lasteventtimestamp:  proto.Int64(ls.LastEventTs),
		Lastingestiontime:   proto.Int64(ls.LastIngestionTs),
		Uploadsequencetoken: proto.String(ls.UploadSequenceToken),
	}
}

// pbLogGroupClassToString converts a proto LogGroupClass enum to the
// string representation used by the service layer.
func pbLogGroupClassToString(c pb.LogGroupClass) string {
	switch c {
	case pb.LogGroupClass_LOG_GROUP_CLASS_DELIVERY:
		return "DELIVERY"
	case pb.LogGroupClass_LOG_GROUP_CLASS_INFREQUENT_ACCESS:
		return "INFREQUENT_ACCESS"
	default:
		return "STANDARD"
	}
}
