package kinesis

import (
	"encoding/base64"
	"regexp"

	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// streamNamePattern matches valid Kinesis stream names per the Smithy model:
// pattern ^[a-zA-Z0-9_.-]+, length 1-128.
var streamNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{1,128}$`)

// consumerNamePattern matches valid Kinesis consumer names per the Smithy
// model: ConsumerName has the same pattern and length traits as StreamName.
var consumerNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{1,128}$`)

// shardIdPattern matches valid Kinesis shard IDs per the Smithy model:
// pattern ^[a-zA-Z0-9_.-]+, length 1-128.
var shardIdPattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{1,128}$`)

// resourceARNPattern validates Kinesis resource ARNs per the Smithy model:
// ^arn:aws.*:kinesis:.*:\d{12}:.*stream/\S+$
var resourceARNPattern = regexp.MustCompile(`^arn:aws.*:kinesis:.*:\d{12}:.*stream/\S+$`)

// validShardLevelMetrics is the set of valid ShardLevelMetrics enum values.
// Values match the enumValue traits in the Smithy model (camelCase, not
// the SCREAMING_SNAKE_CASE member names).
var validShardLevelMetrics = map[string]bool{
	"IncomingBytes":                      true,
	"IncomingRecords":                    true,
	"OutgoingBytes":                      true,
	"OutgoingRecords":                    true,
	"WriteProvisionedThroughputExceeded": true,
	"ReadProvisionedThroughputExceeded":  true,
	"IteratorAgeMilliseconds":            true,
	"ALL":                                true,
}

// validStreamModes is the set of valid StreamMode enum values from the
// Smithy model.
var validStreamModes = map[string]bool{
	"PROVISIONED": true,
	"ON_DEMAND":   true,
}

// validIteratorTypes is the set of valid ShardIteratorType enum values from
// the Smithy model.
var validIteratorTypes = map[string]bool{
	"TRIM_HORIZON":          true,
	"LATEST":                true,
	"AT_SEQUENCE_NUMBER":    true,
	"AFTER_SEQUENCE_NUMBER": true,
	"AT_TIMESTAMP":          true,
}

// validShardFilterTypes is the set of valid ShardFilterType enum values from
// the Smithy model.
var validShardFilterTypes = map[string]bool{
	"AT_TIMESTAMP":          true,
	"FROM_TRIM_HORIZON":     true,
	"FROM_LATEST":           true,
	"AT_SEQUENCE_NUMBER":    true,
	"AFTER_SEQUENCE_NUMBER": true,
}

// validShardOrders is the set of valid ShardOrder enum values from the
// Smithy model.
var validShardOrders = map[string]bool{
	"ASCENDING":  true,
	"DESCENDING": true,
}

// validateStreamName checks that a stream name matches the AWS naming rules.
// Smithy: length 1-128, pattern ^[a-zA-Z0-9_.-]+
func validateStreamName(name string) bool {
	return streamNamePattern.MatchString(name)
}

// validateConsumerName checks that a consumer name matches the AWS naming
// rules. Smithy: length 1-128, pattern ^[a-zA-Z0-9_.-]+
func validateConsumerName(name string) bool {
	return consumerNamePattern.MatchString(name)
}

// validateShardId checks that a shard ID matches the AWS naming rules.
// Smithy: length 1-128, pattern ^[a-zA-Z0-9_.-]+
func validateShardId(shardID string) bool {
	return shardIdPattern.MatchString(shardID)
}

// validatePartitionKey checks that a partition key is within the allowed
// length. Smithy: length 1-256.
func validatePartitionKey(key string) bool {
	return len(key) >= 1 && len(key) <= 256
}

// validateShardCount checks that a shard count is positive.
// Smithy ShardCountObject range is 0-1000000, but for CreateStream the
// practical minimum is 1.
func validateShardCount(count int32) bool {
	return count >= 1
}

// validateRetentionPeriod checks that a retention period is within the
// AWS-allowed range. AWS docs: minimum 24, maximum 8760 hours (365 days).
func validateRetentionPeriod(hours int32) bool {
	return hours >= 24 && hours <= 8760
}

// validateMaxRecordSizeInKiB checks that a max record size is within the
// Smithy range trait. Smithy MaxRecordSizeInKiB: range 1024-10240.
func validateMaxRecordSizeInKiB(kib int32) bool {
	return kib >= 1024 && kib <= 10240
}

// validateWarmThroughputMiBps checks that a warm throughput value is positive.
// AWS docs do not specify an upper bound for WarmThroughputMiBps.
func validateWarmThroughputMiBps(mibps int32) bool {
	return mibps >= 1
}

// validateGetRecordsLimit checks that a GetRecords Limit is within the
// Smithy range trait. GetRecordsInputLimit: range 1-10000.
func validateGetRecordsLimit(limit int32) bool {
	return limit >= 1 && limit <= 10000
}

// validateListStreamsLimit checks that a ListStreams Limit is within the
// Smithy range trait. ListStreamsInputLimit: range 1-10000.
func validateListStreamsLimit(limit int) bool {
	return limit >= 1 && limit <= 10000
}

// validateListShardsLimit checks that a ListShards MaxResults is within the
// Smithy range trait. ListShardsInputLimit: range 1-10000.
func validateListShardsLimit(limit int) bool {
	return limit >= 1 && limit <= 10000
}

// validateListStreamConsumersLimit checks that a ListStreamConsumers
// MaxResults is within the Smithy range trait. ListStreamConsumersInputLimit:
// range 1-10000.
func validateListStreamConsumersLimit(limit int) bool {
	return limit >= 1 && limit <= 10000
}

// validateResourceARN checks that a resource ARN matches the Kinesis ARN
// format. Smithy: pattern ^arn:aws.*:kinesis:.*:\d{12}:.*stream/\S+$,
// length 1-2048.
func validateResourceARN(arn string) bool {
	if len(arn) < 1 || len(arn) > 2048 {
		return false
	}
	return resourceARNPattern.MatchString(arn)
}

// validateShardLevelMetric checks that a metric name is a valid
// ShardLevelMetrics enum value.
func validateShardLevelMetric(metric string) bool {
	return validShardLevelMetrics[metric]
}

// validateStreamMode checks that a stream mode is a valid StreamMode enum value.
func validateStreamMode(mode string) bool {
	return validStreamModes[mode]
}

// validateIteratorType checks that an iterator type is a valid
// ShardIteratorType enum value.
func validateIteratorType(t string) bool {
	return validIteratorTypes[t]
}

// validateShardFilterType checks that a filter type is a valid
// ShardFilterType enum value.
func validateShardFilterType(t string) bool {
	return validShardFilterTypes[t]
}

// validateShardOrder checks that an order is a valid ShardOrder enum value.
func validateShardOrder(order string) bool {
	return validShardOrders[order]
}

// validateKeyId checks that a KMS key identifier is within the allowed length.
// Smithy: KeyId has length trait 1-2048 but no pattern trait. AWS accepts
// UUID, key ARN, alias name (alias/my-key), and alias ARN.
func validateKeyId(keyID string) bool {
	return len(keyID) >= 1 && len(keyID) <= 2048
}

// validateRecordDataSize decodes the base64-encoded Data and checks the
// decoded byte length against the stream's max record size. AWS measures
// record size on the raw payload, not the base64-encoded representation.
// Smithy Data shape: blob length 0-10485760 (10 MiB).
func validateRecordDataSize(b64Data string, maxKiB int32) bool {
	decoded, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		return false
	}
	maxBytes := int(maxKiB) * 1024
	if maxBytes <= 0 {
		maxBytes = 1048576 // 1 MiB default
	}
	return len(decoded) <= maxBytes
}

// validateStreamModeValue validates a kinesisstore.StreamMode value.
func validateStreamModeValue(mode kinesisstore.StreamMode) bool {
	return validStreamModes[string(mode)]
}
