package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The S3-typed destination family of the vended-logs delivery engine:
// the pass writes one gzipped object per delivery under the documented
// AWSLogs/<account-id>/ prefix, with the record shaping (fields,
// delimiter, output format) and the object key layout (suffix path
// variables, Hive-compatible form, destination prefix) the delivery and
// destination configurations select.

// deliverEventsToS3 writes an S3-typed delivery as one gzipped object
// per pass, under the documented AWSLogs/<account-id>/ prefix. The
// object name carries the delivery source and a unique batch stamp; the
// path section honours the delivery's suffixPath variables, the default
// date layout, and the Hive-compatible form. The boolean reports
// success; the cursor advances only then.
func (s *LogsService) deliverEventsToS3(ctx context.Context, region string, dest *logsstore.DeliveryDestination, source *logsstore.DeliverySource, delivery *logsstore.Delivery, events []deliveryEvent) bool {
	if s.eventBus() == nil || s.eventBus().S3Invoker() == nil {
		logs.Warn("Delivery to S3 destination has no S3 invoker configured",
			logs.String("deliveryId", delivery.Id))
		return false
	}
	bucket, destPrefix, ok := s3BucketAndPrefixFromArn(dest.DestinationResourceArn)
	if !ok {
		logs.Warn("Delivery destination ARN does not address an S3 bucket",
			logs.String("deliveryId", delivery.Id),
			logs.String("destinationResourceArn", dest.DestinationResourceArn))
		return false
	}

	fields := deliveryRecordFields(delivery)
	body, err := renderDeliveryBody(dest.OutputFormat, delivery.FieldDelimiter, fields, events)
	if err != nil {
		logs.Error("Failed to render delivery body",
			logs.String("deliveryId", delivery.Id), logs.Err(err))
		return false
	}
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if _, err := gw.Write(body); err != nil {
		return false
	}
	if err := gw.Close(); err != nil {
		return false
	}

	now := time.Now().UTC()
	key := deliveryS3ObjectKey(s.accountID, source.Name, delivery, dest, now, destPrefix)
	contentType := "application/x-gzip"
	if err := s.eventBus().S3Invoker().PutObject(ctx, region, bucket, key, buf.Bytes(), contentType); err != nil {
		logs.Warn("Delivery to S3 destination failed",
			logs.String("deliveryId", delivery.Id),
			logs.String("bucket", bucket), logs.String("key", key), logs.Err(err))
		return false
	}
	logs.Info("Delivered logs to S3 destination",
		logs.String("deliveryId", delivery.Id),
		logs.String("bucket", bucket),
		logs.String("key", key),
		logs.Int("events", len(events)))
	return true
}

// deliveryRecordFields resolves the delivery's effective record fields:
// the configured set, or the template default (every vocabulary field)
// when the delivery carries none.
func deliveryRecordFields(delivery *logsstore.Delivery) []string {
	if len(delivery.RecordFields) > 0 {
		return delivery.RecordFields
	}
	return deliveryDefaultRecordFields
}

// renderDeliveryBody shapes the delivery's events per the destination's
// output format. json renders one object per line carrying the delivery's
// record fields; plain and w3c render the fields joined by the delimiter
// (w3c under the W3C "#Fields:" header line); raw renders the event
// message alone.
func renderDeliveryBody(outputFormat, fieldDelimiter string, fields []string, events []deliveryEvent) ([]byte, error) {
	var b strings.Builder
	if outputFormat == "w3c" {
		b.WriteString("#Fields: " + strings.Join(fields, " ") + "\n")
	}
	for _, e := range events {
		line, err := deliveryEventLine(outputFormat, fieldDelimiter, fields, e)
		if err != nil {
			return nil, err
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return []byte(b.String()), nil
}

// deliveryEventLine renders one event as the single shaped line the
// delivered form carries for it — the unit both the S3 body and a CWL
// destination's re-ingested message are built from.
func deliveryEventLine(outputFormat, fieldDelimiter string, fields []string, e deliveryEvent) (string, error) {
	if fieldDelimiter == "" {
		fieldDelimiter = deliveryDefaultFieldDelimiter
	}
	switch outputFormat {
	case "raw":
		return e.Message, nil
	case "plain", "w3c", "":
		return strings.Join(deliveryFieldValues(e, fields), fieldDelimiter), nil
	default: // json
		values := deliveryFieldValues(e, fields)
		record := make(map[string]interface{}, len(fields))
		for i, f := range fields {
			record[f] = values[i]
		}
		line, err := json.Marshal(record)
		if err != nil {
			return "", err
		}
		return string(line), nil
	}
}

// deliveryFieldValues renders the vocabulary fields of one event in the
// configured order. The time field renders as epoch milliseconds (the
// event's own unit); stream is the originating log stream; message is
// the event message.
func deliveryFieldValues(e deliveryEvent, fields []string) []string {
	values := make([]string, len(fields))
	for i, f := range fields {
		switch f {
		case "time":
			values[i] = fmt.Sprintf("%d", e.Timestamp)
		case "stream":
			values[i] = e.LogStreamName
		default:
			values[i] = e.Message
		}
	}
	return values
}

// deliveryS3ObjectKey builds the delivered object's key. The head is the
// documented granted prefix AWSLogs/<account-id>/<source> — rendered
// AWSLogs/aws-account-id=<account-id>/<source> under Hive-compatible
// paths ("the default account segment AWSLogs/{{source-account-id}}/
// becomes AWSLogs/aws-account-id={{source-account-id}}/") — unless the
// destination ARN carries a destination prefix, which "replaces that
// default path" ("For log types that otherwise use a default
// AWSLogs/source-account-id/service-name/ path, the destination prefix
// replaces that default path"). The path section honours the delivery's
// suffixPath variables (rendered key=value under Hive-compatible paths —
// "variables in the effective path are rendered as key=value", the
// documented suffix example myFolder/{yyyy}/{MM}/{dd} yielding
// myFolder/year=2026/month=09/day=10/) or the default date layout, and
// the object name carries the delivery source and a unique batch stamp.
func deliveryS3ObjectKey(accountID, sourceName string, delivery *logsstore.Delivery, dest *logsstore.DeliveryDestination, now time.Time, destPrefix string) string {
	head := fmt.Sprintf("AWSLogs/%s/%s", accountID, sourceName)
	if delivery.S3HiveCompatiblePath {
		head = fmt.Sprintf("AWSLogs/aws-account-id=%s/%s", accountID, sourceName)
	}
	if destPrefix != "" {
		head = destPrefix
	}
	pathSection := now.Format("2006/01/02")
	if delivery.S3HiveCompatiblePath {
		pathSection = fmt.Sprintf("year=%04d/month=%02d/day=%02d",
			now.Year(), now.Month(), now.Day())
	}
	if delivery.S3SuffixPath != "" {
		pathSection = substituteSuffixPath(delivery.S3SuffixPath, accountID, now, delivery.S3HiveCompatiblePath)
	}
	ext := "json"
	switch dest.OutputFormat {
	case "plain", "w3c", "raw":
		ext = "txt"
	}
	batch := fmt.Sprintf("%s_%s_%s.%s.gz",
		sourceName, now.Format("20060102T150405Z"), randomBatchToken(), ext)
	return fmt.Sprintf("%s/%s/%s", head, pathSection, batch)
}

// substituteSuffixPath expands the suffixPath member's variable sections
// into the batch's values. The variable form is {name} with the
// template's allowedSuffixPathFields vocabulary; under Hive-compatible
// paths each variable renders as its key=value form instead.
func substituteSuffixPath(suffixPath, accountID string, now time.Time, hive bool) string {
	replacements := map[string]string{
		"accountId": accountID,
		"yyyy":      now.Format("2006"),
		"MM":        now.Format("01"),
		"dd":        now.Format("02"),
		"HH":        now.Format("15"),
	}
	hiveKeys := map[string]string{
		"accountId": "aws-account-id",
		"yyyy":      "year",
		"MM":        "month",
		"dd":        "day",
		"HH":        "hour",
	}
	result := suffixPath
	for name, value := range replacements {
		if hive {
			value = hiveKeys[name] + "=" + value
		}
		result = strings.ReplaceAll(result, "{"+name+"}", value)
	}
	return result
}

// randomBatchToken is the uniqueness component of the delivered object's
// name.
func randomBatchToken() string {
	return randomHexWithFallback(4, func() string { return fmt.Sprintf("%d", nowNanos()) })
}

func nowNanos() int64 { return time.Now().UnixNano() }
