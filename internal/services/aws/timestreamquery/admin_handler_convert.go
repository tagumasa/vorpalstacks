package timestreamquery

import (
	"net/http"
	"vorpalstacks/internal/common/defaults"

	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/timestreamquery"
)

// This file is the sole exception to the #29 rule: it imports store packages
// and provides conversion functions so that admin_handler.go stays free of
// store imports.

// getStores retrieves the store group for the request region from the shared
// service cache.
func (h *AdminHandler) getStores(region string) (*tsQueryStores, error) {
	return h.service.GetScheduledQueryStoreForRegionStores(region)
}

// resolveRegion extracts the region from the request header.
func resolveRegion(header http.Header) string {
	return defaults.GetRegionFromHeader(header)
}

// toPbScheduledQuery converts a service-layer ScheduledQuerySummary DTO into
// the protobuf representation used by the admin console.
func toPbScheduledQuery(summary *ScheduledQuerySummary) *pb.ScheduledQuery {
	sq := &pb.ScheduledQuery{
		Arn:          summary.ARN,
		Name:         summary.Name,
		Creationtime: summary.CreationTime.Format(timeutils.ISO8601UTCFormat),
		State:        pb.ScheduledQueryState_SCHEDULED_QUERY_STATE_ENABLED,
	}

	if summary.State == "DISABLED" {
		sq.State = pb.ScheduledQueryState_SCHEDULED_QUERY_STATE_DISABLED
	}

	if summary.LastRunStatus != "" {
		sq.Lastrunstatus = mapLastRunStatusToProto(summary.LastRunStatus)
	}

	if summary.ErrorReportConfiguration != nil && summary.ErrorReportConfiguration.S3Configuration != nil {
		sq.Errorreportconfiguration = &pb.ErrorReportConfiguration{
			S3Configuration: &pb.S3Configuration{
				Bucketname: summary.ErrorReportConfiguration.S3Configuration.BucketName,
			},
		}
		if summary.ErrorReportConfiguration.S3Configuration.ObjectKeyPrefix != "" {
			sq.Errorreportconfiguration.S3Configuration.Objectkeyprefix = summary.ErrorReportConfiguration.S3Configuration.ObjectKeyPrefix
		}
	}

	if summary.TargetConfiguration != nil && summary.TargetConfiguration.TimestreamConfiguration != nil {
		sq.Targetdestination = &pb.TargetDestination{
			Timestreamdestination: &pb.TimestreamDestination{
				Databasename: summary.TargetConfiguration.TimestreamConfiguration.DatabaseName,
				Tablename:    summary.TargetConfiguration.TimestreamConfiguration.TableName,
			},
		}
	}

	if !summary.PreviousRunTime.IsZero() {
		sq.Previousinvocationtime = summary.PreviousRunTime.Format(timeutils.ISO8601UTCFormat)
	}

	if !summary.NextRunTime.IsZero() {
		sq.Nextinvocationtime = summary.NextRunTime.Format(timeutils.ISO8601UTCFormat)
	}

	return sq
}

// mapLastRunStatusToProto maps the store-level ScheduledQueryRunStatus string
// to the proto enum used by the admin console.
func mapLastRunStatusToProto(status string) pb.ScheduledQueryRunStatus {
	switch status {
	case "MANUAL_TRIGGER_SUCCESS":
		return pb.ScheduledQueryRunStatus_SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_SUCCESS
	case "MANUAL_TRIGGER_FAILURE":
		return pb.ScheduledQueryRunStatus_SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_FAILURE
	case "AUTO_TRIGGER_FAILURE":
		return pb.ScheduledQueryRunStatus_SCHEDULED_QUERY_RUN_STATUS_AUTO_TRIGGER_FAILURE
	case "AUTO_TRIGGER_SUCCESS":
		return pb.ScheduledQueryRunStatus_SCHEDULED_QUERY_RUN_STATUS_AUTO_TRIGGER_SUCCESS
	default:
		return pb.ScheduledQueryRunStatus_SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_SUCCESS
	}
}
