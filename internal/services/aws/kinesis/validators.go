package kinesis

import (
	"encoding/base64"
	"math/big"
	"regexp"
	"unicode/utf8"

	"vorpalstacks/internal/common/tags"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// streamNamePattern matches valid Kinesis stream names per the Smithy model:
// pattern ^[a-zA-Z0-9_.-]+, length 1-128.
var streamNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{1,128}$`)

// consumerNamePattern matches valid Kinesis consumer names per the Smithy
// model: ConsumerName has the same pattern and length traits as StreamName.
var consumerNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{1,128}$`)

// shardIdPattern matches valid Kinesis shard IDs per the Smithy model's
// ShardId shape: pattern ^[a-zA-Z0-9_.-]+, length 1-128.
var shardIdPattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{1,128}$`)

// resourceARNPattern validates Kinesis resource ARNs per the Smithy model:
// ^arn:aws.*:kinesis:.*:\d{12}:.*(stream|channel)/\S+$ — the shape accepts
// both stream and channel ARNs.
var resourceARNPattern = regexp.MustCompile(`^arn:aws.*:kinesis:.*:\d{12}:.*(stream|channel)/\S+$`)

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

// validShardFilterTypes is the set of valid ShardFilterType enum values
// from the Smithy model.
var validShardFilterTypes = map[string]bool{
	"AFTER_SHARD_ID":    true,
	"AT_TRIM_HORIZON":   true,
	"FROM_TRIM_HORIZON": true,
	"AT_LATEST":         true,
	"AT_TIMESTAMP":      true,
	"FROM_TIMESTAMP":    true,
}

// validMTBCStatuses is the set of valid
// MinimumThroughputBillingCommitmentInputStatus enum values from the
// Smithy model — the store package's two constants are the vocabulary.
var validMTBCStatuses = map[string]bool{
	kinesisstore.MTBCStatusEnabled:  true,
	kinesisstore.MTBCStatusDisabled: true,
}

// sequenceNumberPattern matches the model's SequenceNumber shape: an
// optional "0" or a non-zero-leading decimal of up to 129 digits.
var sequenceNumberPattern = regexp.MustCompile(`^(0|[1-9][0-9]{0,128})$`)

// The MetricsNameList length trait bounds one enhanced-monitoring request's
// ShardLevelMetrics member: at least one metric, at most the seven
// individual metrics the MetricsName enum defines (the ALL wildcard counts
// as one).
const (
	minShardLevelMetricsPerRequest = 1
	maxShardLevelMetricsPerRequest = 7
)

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

// validateShardId checks a shard-id member (GetShardIterator ShardId,
// SplitShard ShardToSplit, MergeShards' pair, SubscribeToShard ShardId)
// against the ShardId shape's traits. Smithy: length 1-128,
// pattern ^[a-zA-Z0-9_.-]+ — a malformed ID answers invalid input instead
// of falling into a not-found lookup.
func validateShardId(id string) bool {
	return shardIdPattern.MatchString(id)
}

// validatePartitionKey checks that a partition key is within the allowed
// length. Smithy: length 1-256 counted in Unicode characters (the shape
// carries no pattern, so multibyte partition keys are valid input).
func validatePartitionKey(key string) bool {
	n := utf8.RuneCountInString(key)
	return n >= 1 && n <= 256
}

// explicitHashKeyPattern matches the model's HashKey shape: optional "0" or
// a non-zero-leading decimal.
var explicitHashKeyPattern = regexp.MustCompile(`^(0|[1-9][0-9]{0,38})$`)

// validateExplicitHashKey checks a hash-key member (PutRecord and
// PutRecords ExplicitHashKey, SplitShard NewStartingHashKey) against the
// model pattern and the shard hash-key space: shard ranges span the
// model's 128-bit domain, so a value above 2^128-1 can address no shard —
// it is rejected up front as invalid input instead of failing to route at
// write time. An empty member is unset and valid.
func validateExplicitHashKey(key string) bool {
	if key == "" {
		return true
	}
	if !explicitHashKeyPattern.MatchString(key) {
		return false
	}
	value, ok := new(big.Int).SetString(key, 10)
	if !ok {
		return false
	}
	return value.Cmp(kinesisstore.MaxShardHashKeyInt()) <= 0
}

// validateShardCount checks that a shard count is within the allowed
// range: the CreateStream member targets the PositiveIntegerObject shape
// (range trait min 1) and the platform caps counts at the store's
// documented maximum.
func validateShardCount(count int32) bool {
	return count >= 1 && count <= kinesisstore.MaxShardCount
}

// validateRetentionPeriod checks that a retention period is within the
// AWS-allowed range. AWS docs: minimum 24, maximum 8760 hours (365 days).
func validateRetentionPeriod(hours int32) bool {
	return hours >= kinesisstore.MinRetentionPeriodHours && hours <= kinesisstore.MaxRetentionPeriodHours
}

// validateMaxRecordSizeInKiB checks that a max record size is within the
// Smithy range trait. Smithy MaxRecordSizeInKiB: range 1024-10240.
func validateMaxRecordSizeInKiB(kib int32) bool {
	return kib >= kinesisstore.MinMaxRecordSizeInKiB && kib <= kinesisstore.MaxMaxRecordSizeInKiB
}

// validateWarmThroughputMiBps checks that a warm throughput value is within
// the model's range. Every warm-throughput member targets the
// NaturalIntegerObject shape whose range trait carries min 0 alone — zero
// is the documented release floor (the API reference's valid range is a
// minimum of 0, and releasing excess capacity means setting the same or a
// lower value) — so negatives alone reject.
func validateWarmThroughputMiBps(mibps int32) bool {
	return mibps >= 0
}

// validateGetRecordsLimit checks that a GetRecords Limit is within the
// Smithy range trait. GetRecordsInputLimit: range 1-10000.
func validateGetRecordsLimit(limit int32) bool {
	return limit >= kinesisstore.MinListResultsLimit && limit <= kinesisstore.MaxListResultsLimit
}

// validateListStreamsLimit checks that a ListStreams Limit is within the
// Smithy range trait. ListStreamsInputLimit: range 1-10000.
func validateListStreamsLimit(limit int) bool {
	return limit >= kinesisstore.MinListResultsLimit && limit <= kinesisstore.MaxListResultsLimit
}

// validateListShardsLimit checks that a ListShards MaxResults is within the
// Smithy range trait. ListShardsInputLimit: range 1-10000.
func validateListShardsLimit(limit int) bool {
	return limit >= kinesisstore.MinListResultsLimit && limit <= kinesisstore.MaxListResultsLimit
}

// validateListStreamConsumersLimit checks that a ListStreamConsumers
// MaxResults is within the Smithy range trait. ListStreamConsumersInputLimit:
// range 1-10000.
func validateListStreamConsumersLimit(limit int) bool {
	return limit >= kinesisstore.MinListResultsLimit && limit <= kinesisstore.MaxListResultsLimit
}

// validateResourceARN checks that a resource ARN matches the Kinesis ARN
// format. Smithy: pattern ^arn:aws.*:kinesis:.*:\d{12}:.*(stream|channel)/\S+$,
// length 1-2048 — the shape accepts both stream and channel ARNs.
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

// validateShardFilter checks a ShardFilter against the model: the Type is
// a required property taking the enum, and the paired members travel only
// with their types — the ShardId member's documentation limits it to
// AFTER_SHARD_ID and the Timestamp member's to AT_TIMESTAMP/FROM_TIMESTAMP,
// while those types' own semantics presuppose their member.
func validateShardFilter(filter *kinesisstore.ShardFilter) error {
	if !validShardFilterTypes[filter.Type] {
		return ErrInvalidArgument
	}
	if filter.ShardID != "" && filter.Type != "AFTER_SHARD_ID" {
		return ErrInvalidArgument
	}
	if filter.Type == "AFTER_SHARD_ID" && filter.ShardID == "" {
		return ErrInvalidArgument
	}
	if filter.Timestamp != nil && filter.Type != "AT_TIMESTAMP" && filter.Type != "FROM_TIMESTAMP" {
		return ErrInvalidArgument
	}
	if (filter.Type == "AT_TIMESTAMP" || filter.Type == "FROM_TIMESTAMP") && filter.Timestamp == nil {
		return ErrInvalidArgument
	}
	return nil
}

// validateKeyId checks that a KMS key identifier is within the allowed length.
// Smithy: KeyId has length trait 1-2048 but no pattern trait. AWS accepts
// UUID, key ARN, alias name (alias/my-key), and alias ARN. Lengths count
// Unicode characters.
func validateKeyId(keyID string) bool {
	n := utf8.RuneCountInString(keyID)
	return n >= 1 && n <= 2048
}

// validateRecordSize measures one record against the stream's maximum
// record size the way the Data member's documentation defines the limit:
// the decoded payload ("the data blob, the payload before base64-encoding")
// added to the partition key size must not exceed the maximum — a payload
// that fits alone can still fail through its key. The limit counts bytes;
// AWS measures record size on the raw payload, not the base64-encoded
// representation.
func validateRecordSize(b64Data, partitionKey string, maxKiB int32) bool {
	decoded, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		return false
	}
	maxBytes := int(maxKiB) * 1024
	if maxBytes <= 0 {
		maxBytes = kinesisstore.DefaultMaxRecordSizeBytes
	}
	return len(decoded)+len(partitionKey) <= maxBytes
}

// validateSequenceNumber checks a sequence-number member (GetShardIterator
// StartingSequenceNumber, PutRecord SequenceNumberForOrdering) against the
// model's SequenceNumber shape pattern. An empty member is unset and valid —
// presence requirements are cross-member rules the callers enforce.
func validateSequenceNumber(sn string) bool {
	if sn == "" {
		return true
	}
	return sequenceNumberPattern.MatchString(sn)
}

// validateMTBCStatus checks a commitment status against the
// MinimumThroughputBillingCommitmentInputStatus enum.
func validateMTBCStatus(status string) bool {
	return validMTBCStatuses[status]
}

// validateShardLevelMetricsList checks the MetricsNameList length trait on
// one enhanced-monitoring request: the member is required and carries
// between one and seven metric names.
func validateShardLevelMetricsList(metrics []string) bool {
	return len(metrics) >= minShardLevelMetricsPerRequest && len(metrics) <= maxShardLevelMetricsPerRequest
}

// hashKeyWithinRange reports whether a hash key lies inside a shard's hash
// key range so that both children of a split are well-formed: the split
// member's documentation requires the NewStartingHashKey to be in the range
// being mapped into the shard, and the store lays the children out as
// [start, NSK-1] and [NSK, end] — a key at the range's starting boundary
// inverts the lower child, and a key beyond the ending boundary leaves the
// upper child outside the parent.
func hashKeyWithinRange(key string, r *kinesisstore.HashKeyRange) bool {
	low, okLow := new(big.Int).SetString(r.StartingHashKey, 10)
	high, okHigh := new(big.Int).SetString(r.EndingHashKey, 10)
	value, okValue := new(big.Int).SetString(key, 10)
	if !okLow || !okHigh || !okValue {
		return false
	}
	return value.Cmp(low) > 0 && value.Cmp(high) <= 0
}

// validateTagKeysList checks a TagKeyList against the model's traits: the
// list carries at most MaxTagKeysPerUntagRequest keys and every key is a
// well-formed TagKey (1-128 characters, the shared tag-key bound).
func validateTagKeysList(tagKeys []string) bool {
	if len(tagKeys) == 0 || len(tagKeys) > kinesisstore.MaxTagKeysPerUntagRequest {
		return false
	}
	for _, key := range tagKeys {
		if utf8.RuneCountInString(key) < 1 || utf8.RuneCountInString(key) > tags.MaxTagKeyLength {
			return false
		}
	}
	return true
}

// validateStartingPosition applies the position-member pairing the model
// documents identically on GetShardIterator's ShardIteratorType and
// SubscribeToShard's StartingPosition.Type: the type must be one of the
// five modelled values (checked by validateIteratorType), a sequence
// number that travels must satisfy the SequenceNumber shape's pattern and
// is required by AT_SEQUENCE_NUMBER and AFTER_SEQUENCE_NUMBER, and a
// timestamp is required by AT_TIMESTAMP — an unpaired or absent position
// has no documented reading and must not silently degrade to the horizon.
func validateStartingPosition(posType, sequenceNumber string, hasTimestamp bool) bool {
	if !validateIteratorType(posType) {
		return false
	}
	if !validateSequenceNumber(sequenceNumber) {
		return false
	}
	switch posType {
	case "AT_SEQUENCE_NUMBER", "AFTER_SEQUENCE_NUMBER":
		return sequenceNumber != ""
	case "AT_TIMESTAMP":
		return hasTimestamp
	}
	return true
}
