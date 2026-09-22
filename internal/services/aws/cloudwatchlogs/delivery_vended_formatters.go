package cloudwatchlogs

import (
	"strings"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// The wire formatters of the delivery family and the S3 bucket ARN
// helper they share.

// --- formatters ---// formatDeliverySource renders the model's DeliverySource shape from the
// core-computed view (the status members come from the view; the
// formatter stays a pure renderer).
func formatDeliverySource(view *DeliverySourceView) map[string]interface{} {
	source := view.Source
	result := map[string]interface{}{
		"name":         source.Name,
		"arn":          source.Arn,
		"resourceArns": []string{source.ResourceArn},
		"service":      source.Service,
		"logType":      source.LogType,
		"status":       view.Status,
	}
	if view.StatusReason != "" {
		result["statusReason"] = view.StatusReason
	}
	if len(source.DeliverySourceConfiguration) > 0 {
		result["deliverySourceConfiguration"] = source.DeliverySourceConfiguration
	}
	if len(source.Tags) > 0 {
		result["tags"] = source.Tags
	}
	return result
}

// formatDeliveryDestination renders the model's DeliveryDestination
// shape.
func formatDeliveryDestination(dest *logsstore.DeliveryDestination) map[string]interface{} {
	result := map[string]interface{}{
		"name":                    dest.Name,
		"arn":                     dest.Arn,
		"deliveryDestinationType": dest.DeliveryDestinationType,
		"deliveryDestinationConfiguration": map[string]interface{}{
			"destinationResourceArn": dest.DestinationResourceArn,
		},
	}
	if dest.OutputFormat != "" {
		result["outputFormat"] = dest.OutputFormat
	}
	if len(dest.Tags) > 0 {
		result["tags"] = dest.Tags
	}
	return result
}

// formatDelivery renders the model's Delivery shape.
func formatDelivery(d *logsstore.Delivery) map[string]interface{} {
	result := map[string]interface{}{
		"id":                      d.Id,
		"arn":                     d.Arn,
		"deliverySourceName":      d.DeliverySourceName,
		"deliveryDestinationArn":  d.DeliveryDestinationArn,
		"deliveryDestinationType": d.DeliveryDestinationType,
	}
	if len(d.RecordFields) > 0 {
		result["recordFields"] = d.RecordFields
	}
	if d.FieldDelimiter != "" {
		result["fieldDelimiter"] = d.FieldDelimiter
	}
	if d.S3SuffixPath != "" || d.S3HiveCompatiblePath {
		s3 := map[string]interface{}{}
		if d.S3SuffixPath != "" {
			s3["suffixPath"] = d.S3SuffixPath
		}
		if d.S3HiveCompatiblePath {
			s3["enableHiveCompatiblePath"] = true
		}
		result["s3DeliveryConfiguration"] = s3
	}
	if len(d.Tags) > 0 {
		result["tags"] = d.Tags
	}
	return result
}

// s3BucketFromArn extracts the bucket name from an S3 bucket ARN. The
// bucket sits in the ARN's resource field (arn:aws:s3:::bucket), with a
// key suffix when the ARN addresses an object.
func s3BucketFromArn(ar string) (string, bool) {
	bucket, _, ok := s3BucketAndPrefixFromArn(ar)
	return bucket, ok
}

// s3BucketAndPrefixFromArn splits an S3 destination ARN into its bucket
// and destination prefix: "An optional path that you append to the
// bucket ARN when you call PutDeliveryDestination. For example,
// arn:aws:s3:::{{bucket-name}}/{{MyLogPrefix}}. Delivered objects begin
// with this prefix." — the resource's post-slash remainder, with its
// boundary slashes trimmed.
func s3BucketAndPrefixFromArn(ar string) (string, string, bool) {
	parsed, err := svcarn.ParseARN(ar)
	if err != nil || parsed.Service != "s3" {
		return "", "", false
	}
	bucket := parsed.Resource
	prefix := ""
	if idx := strings.Index(bucket, "/"); idx >= 0 {
		prefix = strings.Trim(bucket[idx+1:], "/")
		bucket = bucket[:idx]
	}
	if bucket == "" {
		return "", "", false
	}
	return bucket, prefix, true
}
