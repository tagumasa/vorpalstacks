package cloudwatchlogs

import (
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The event-plane validators: the scheduled query description, the
// StartFromHead date constraint, the PutLogEvents event rows, and the
// scheduled query destination configuration.

// validateScheduledQueryDescription enforces the
// ScheduledQueryDescription length trait (0-1024).
func validateScheduledQueryDescription(d string) error {
	if utf8.RuneCountInString(d) > logsstore.MaxScheduledQueryDescriptionLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Scheduled query description must not exceed %d characters", logsstore.MaxScheduledQueryDescriptionLength), 400)
	}
	return nil
}

// --- StartFromHead date constraint ---

// jan012024Millis is Jan 1, 2024 00:00:00 UTC in epoch milliseconds.
// AWS requires startTime on or after this date when startFromHead=false.
const jan012024Millis = 1704067200000

// validateStartFromHeadDate enforces the AWS requirement that when
// startFromHead is false, startTime must be on or after Jan 1, 2024.
func validateStartFromHeadDate(startFromHead bool, startTime int64) error {
	if !startFromHead && startTime > 0 && startTime < jan012024Millis {
		return NewLogsError("InvalidParameterException",
			"Setting startFromHead to false is supported only when startTime is on or after Jan 1, 2024 00:00:00 UTC", 400)
	}
	return nil
}

// --- Log event validation (PutLogEvents) ---

const (
	// maxEventsTimeSpan is the maximum allowed time span (in milliseconds)
	// for a single PutLogEvents batch. AWS rejects the entire batch if the
	// span between the earliest and latest event exceeds 24 hours.
	maxEventsTimeSpan int64 = 24 * 60 * 60 * 1000

	// tooNewThreshold is the maximum future offset (in milliseconds) for
	// an event timestamp. Events more than 2 hours in the future are
	// rejected individually.
	tooNewThreshold int64 = 2 * 60 * 60 * 1000

	// tooOldThreshold is the maximum age (in milliseconds) for an event
	// timestamp. Events older than 14 days are rejected individually.
	tooOldThreshold int64 = 14 * 24 * 60 * 60 * 1000
)

// retentionCutoffMillis returns the epoch-millisecond horizon of a
// retention setting at the given now: events stamped before it precede
// the group's retention period. A non-positive setting never expires.
func retentionCutoffMillis(now int64, retentionInDays int32) int64 {
	if retentionInDays <= 0 {
		return 0
	}
	return now - int64(retentionInDays)*24*60*60*1000
}

// validateLogEvents checks that log events satisfy the PutLogEvents
// constraints required by AWS CloudWatch Logs:
//   - Events must be in chronological order (by timestamp).
//   - No event may be more than 2 hours in the future.
//   - No event may be older than 14 days.
//   - No event may precede the log group's retention period ("Events
//     older than 14 days or preceding the log group's retention period
//     are rejected while processing remaining valid events", PutLogEvents
//     operation documentation) — retentionCutoff is the group's horizon
//     in epoch milliseconds, 0 for a never-expiring group.
//   - The timespan between the earliest and latest valid event must not
//     exceed 24 hours.
//
// Events that fall outside the age thresholds are silently excluded from
// the returned valid slice; the caller receives information about the
// rejected indices via the returned map, which is suitable for inclusion
// in the response as rejectedLogEventsInfo.
func validateLogEvents(events []logsstore.LogEntry, retentionCutoff int64) ([]logsstore.LogEntry, map[string]interface{}, error) {
	now := time.Now().UnixMilli()

	// Chronological order check must be performed on ALL events in the
	// batch, not just the age-valid subset. AWS rejects the entire batch
	// if any event is out of order, regardless of whether some events are
	// later individually rejected for being too old or too new.
	for i := 1; i < len(events); i++ {
		if events[i].Timestamp < events[i-1].Timestamp {
			return nil, nil, NewLogsError("InvalidParameterException",
				"log events in the batch must be in chronological order", 400)
		}
	}

	var valid []logsstore.LogEntry
	var tooOldEndIndex int
	expiredEndIndex := 0
	tooNewStartIndex := -1

	for i, e := range events {
		if e.Timestamp > now+tooNewThreshold {
			if tooNewStartIndex == -1 || i < tooNewStartIndex {
				tooNewStartIndex = i
			}
			continue
		}
		if e.Timestamp < now-tooOldThreshold {
			tooOldEndIndex = i + 1
			continue
		}
		// The retention horizon bites only when it is tighter than the
		// 14-day rule; a chronological batch keeps the expired run right
		// after the too-old run, so the exclusive end index accumulates.
		if retentionCutoff > now-tooOldThreshold && e.Timestamp < retentionCutoff {
			expiredEndIndex = i + 1
			continue
		}
		valid = append(valid, e)
	}

	if len(valid) == 0 {
		rejected := formatRejectedInfo(tooOldEndIndex, expiredEndIndex, tooNewStartIndex, len(events))
		return nil, rejected, nil
	}

	span := valid[len(valid)-1].Timestamp - valid[0].Timestamp
	if span > maxEventsTimeSpan {
		return nil, nil, NewLogsError("InvalidParameterException",
			"Events span must not exceed 24 hours", 400)
	}

	rejected := formatRejectedInfo(tooOldEndIndex, expiredEndIndex, tooNewStartIndex, len(events))
	return valid, rejected, nil
}

// formatRejectedInfo constructs the rejectedLogEventsInfo response member
// from the computed too-old, expired and too-new indices. If no events
// were rejected an empty map is returned.
func formatRejectedInfo(tooOldEndIndex, expiredEndIndex, tooNewStartIndex, totalEvents int) map[string]interface{} {
	if tooOldEndIndex == 0 && expiredEndIndex == 0 && tooNewStartIndex == -1 {
		return nil
	}
	info := make(map[string]interface{})
	if tooOldEndIndex > 0 {
		info["tooOldLogEventEndIndex"] = tooOldEndIndex
	}
	if expiredEndIndex > 0 {
		info["expiredLogEventEndIndex"] = expiredEndIndex
	}
	if tooNewStartIndex >= 0 {
		info["tooNewLogEventStartIndex"] = tooNewStartIndex
	}
	return info
}

// --- Scheduled query destination configuration ---

var (
	// s3DestinationURIPattern is the documented S3Uri pattern of the
	// scheduled-query S3 destination.
	s3DestinationURIPattern = regexp.MustCompile(`^s3://[a-z0-9][.\-a-z0-9]{1,61}[a-z0-9](/.*)?$`)

	// ownerAccountIDPattern is the documented 12-digit account pattern of
	// the scheduled-query S3 destination.
	ownerAccountIDPattern = regexp.MustCompile(`^\d{12}$`)
)

// destinationParam reads a string member of a destination configuration map
// across its case variants.
func destinationParam(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	if v, ok := m[request.LowerFirst(key)]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

// destinationTags reads the tags member of a destination configuration map.
func destinationTags(m map[string]interface{}) map[string]string {
	raw, ok := m["tags"]
	if !ok {
		return nil
	}
	asMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	tags := make(map[string]string, len(asMap))
	for k, v := range asMap {
		if str, ok := v.(string); ok {
			tags[k] = str
		}
	}
	return tags
}

// validateDestinationConfiguration checks the documented member constraints
// of a scheduled-query destination configuration. Delivery to a lookup table
// requires a table name and an IAM role ARN; delivery to S3 requires an S3
// URI and an IAM role ARN.
func validateDestinationConfiguration(dc map[string]interface{}) error {
	if lt, ok := dc["lookupTableConfiguration"].(map[string]interface{}); ok {
		name := destinationParam(lt, "tableName")
		if name == "" {
			return NewLogsError("InvalidParameterException",
				"lookupTableConfiguration requires tableName", 400)
		}
		if len(name) > logsstore.MaxLookupTableNameLength || !lookupTableNameRe.MatchString(name) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Invalid lookupTableName %q: 1-%d characters, alphanumeric and underscores only",
					name, logsstore.MaxLookupTableNameLength), 400)
		}
		roleArn := destinationParam(lt, "roleArn")
		if roleArn == "" {
			return NewLogsError("InvalidParameterException",
				"lookupTableConfiguration requires roleArn", 400)
		}
		if err := validateIAMRoleArn(roleArn); err != nil {
			return err
		}
		if len(destinationParam(lt, "description")) > logsstore.MaxLookupTableDescriptionLength {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("description exceeds %d characters", logsstore.MaxLookupTableDescriptionLength), 400)
		}
		if len(destinationParam(lt, "kmsKeyId")) > logsstore.MaxKmsKeyIdLength {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("kmsKeyId exceeds %d characters", logsstore.MaxKmsKeyIdLength), 400)
		}
		if len(destinationTags(lt)) > logsstore.MaxLookupTableTags {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("a maximum of %d tags can be attached to a lookup table", logsstore.MaxLookupTableTags), 400)
		}
	}
	if s3, ok := dc["s3Configuration"].(map[string]interface{}); ok {
		uri := destinationParam(s3, "destinationIdentifier")
		if uri == "" {
			return NewLogsError("InvalidParameterException",
				"s3Configuration requires destinationIdentifier", 400)
		}
		if len(uri) > logsstore.MaxLookupTableDestinationIdentifierLength || !s3DestinationURIPattern.MatchString(uri) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Invalid S3 destination identifier %q: must be a valid s3:// URI", uri), 400)
		}
		roleArn := destinationParam(s3, "roleArn")
		if roleArn == "" {
			return NewLogsError("InvalidParameterException",
				"s3Configuration requires roleArn", 400)
		}
		if err := validateIAMRoleArn(roleArn); err != nil {
			return err
		}
		if len(destinationParam(s3, "kmsKeyId")) > logsstore.MaxKmsKeyIdLength {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("kmsKeyId exceeds %d characters", logsstore.MaxKmsKeyIdLength), 400)
		}
		if owner := destinationParam(s3, "ownerAccountId"); owner != "" && !ownerAccountIDPattern.MatchString(owner) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Invalid ownerAccountId %q: must be 12 digits", owner), 400)
		}
	}
	return nil
}
