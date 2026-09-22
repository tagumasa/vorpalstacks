package cloudwatchlogs

import (
	"fmt"
	"strings"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/utils/aws/arn"
)

// maxDestinationHops bounds CloudWatch Logs destination indirection during
// delivery. A destination's target is a Kinesis stream ARN (PutDestination's
// targetArn member documentation), so one hop suffices; the bound exists so
// a legacy or mis-stored chain cannot cycle the dispatch forever.
const maxDestinationHops = 3

// dispatchSubscriptionDelivery delivers one matched, gzip-compressed
// subscription payload to its destination. Every delivery leg — the event
// bus handler and the bus-less fallback — routes through this one
// dispatch, so the accepted-destination vocabulary and the delivered
// vocabulary cannot drift apart again. A nil return means the payload
// reached a destination leg; a non-nil return is a delivery failure the
// caller carries into the documented retry window ("Throttled
// deliverables are retried for up to 24 hours. After 24 hours, the failed
// deliverables are dropped").
//
// The destination forms are the ones the PutSubscriptionFilterRequest
// destinationArn member documentation lists: a Kinesis stream ARN, a Lambda
// function ARN, and a CloudWatch Logs destination ARN (the "logical
// destination" form the platform's own PutDestination mints). Kinesis Data
// Firehose is a documented form, but the platform Firehose service does
// not exist yet, so PutSubscriptionFilter and PutDestination reject
// Firehose ARNs outright; the dispatch still recognises one defensively —
// for filters stored before that rejection — reporting it as a delivery
// failure instead of discarding the batch silently.
func (s *LogsService) dispatchSubscriptionDelivery(region, destArn, logGroup, logStream, distribution string, compressed []byte) error {
	current := destArn
	currentRegion := region
	for hop := 0; hop < maxDestinationHops; hop++ {
		switch {
		case arn.IsLambdaARN(current):
			return s.invokeLambda(current, compressed)
		case arn.IsKinesisARN(current):
			return s.putToKinesis(current, logGroup, logStream, distribution, compressed)
		case isFirehoseARN(current):
			logs.Warn("Subscription filter Firehose destination is not deliverable until the platform Firehose service exists",
				logs.String("destinationArn", current))
			return fmt.Errorf("firehose destination %s is not deliverable", current)
		case isCloudWatchLogsDestinationARN(current):
			target, targetRegion := s.resolveDestinationTarget(currentRegion, current)
			if target == "" {
				return fmt.Errorf("destination %s is unresolvable", current)
			}
			current, currentRegion = target, targetRegion
		default:
			logs.Warn("Subscription filter destination is not a deliverable ARN form",
				logs.String("destinationArn", current))
			return fmt.Errorf("destination %s is not a deliverable ARN form", current)
		}
	}
	logs.Warn("Subscription filter destination chain exceeded the destination indirection bound",
		logs.String("destinationArn", destArn),
		logs.Int("bound", maxDestinationHops))
	return fmt.Errorf("destination chain from %s exceeded the indirection bound", destArn)
}

// resolveDestinationTarget resolves the CloudWatch Logs destination ARN to
// the target ARN its stored record carries, in the destination ARN's own
// region. An empty return means unresolvable; the caller has already
// warned.
func (s *LogsService) resolveDestinationTarget(region, destArn string) (string, string) {
	_, _, destRegion, _, resource := arn.SplitARN(destArn)
	if destRegion == "" {
		destRegion = region
	}
	name := destinationNameFromResource(resource)
	if name == "" {
		logs.Warn("Subscription filter destination ARN carries no destination name",
			logs.String("destinationArn", destArn))
		return "", ""
	}
	store, err := s.getLogsStoreByRegion(destRegion)
	if err != nil {
		logs.Warn("Failed to resolve logs store for subscription destination",
			logs.String("region", destRegion), logs.Err(err))
		return "", ""
	}
	dest, err := store.GetDestination(name)
	if err != nil {
		logs.Warn("Subscription filter destination not found",
			logs.String("destinationArn", destArn), logs.Err(err))
		return "", ""
	}
	return dest.TargetArn, destRegion
}

// isCloudWatchLogsDestinationARN reports whether the ARN addresses a
// CloudWatch Logs destination — service "logs" with a destination/ or
// destination: resource prefix (the platform mints the colon form; AWS
// documentation shows the slash form — both address the same record).
func isCloudWatchLogsDestinationARN(ar string) bool {
	if arn.GetServiceFromARN(ar) != "logs" {
		return false
	}
	_, _, _, _, resource := arn.SplitARN(ar)
	return destinationNameFromResource(resource) != ""
}

// destinationNameFromResource extracts the destination name from a logs
// ARN's resource part, or "" when the resource is not a destination.
func destinationNameFromResource(resource string) string {
	if strings.HasPrefix(resource, "destination:") {
		return strings.TrimPrefix(resource, "destination:")
	}
	if strings.HasPrefix(resource, "destination/") {
		return strings.TrimPrefix(resource, "destination/")
	}
	return ""
}
