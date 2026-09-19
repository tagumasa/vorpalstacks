package kinesis

import (
	"math/big"
	"time"
)

// AWS Kinesis service and model limits — the single definition site for
// every documented bound; the raw numbers appear nowhere else. Sources are
// cited per constant: the vendored Smithy model (models/kinesis
// /2013-12-02) range and length traits, and the AWS Kinesis API reference
// quotas.
const (
	// MaxShardHashKey is the inclusive upper bound of the platform's shard
	// hash-key space: 2^128-1, the model's HashKey domain (the shape's
	// pattern admits 39 decimal digits, and AWS's describe-stream output
	// shows a single-shard stream's EndingHashKey at exactly this value).
	// A string constant because the value exceeds uint64;
	// MaxShardHashKeyInt returns it as a big.Int for range and validation
	// arithmetic.
	MaxShardHashKey = "340282366920938463463374607431768211455"

	// Retention period, hours. API reference: the minimum is 24 hours,
	// which is also the create default, up to 8760 hours.
	MinRetentionPeriodHours = 24
	MaxRetentionPeriodHours = 8760

	// MaxRecordSizeInKiB window (model range trait on the member) and the
	// create default of 1 MiB.
	MinMaxRecordSizeInKiB     = 1024
	MaxMaxRecordSizeInKiB     = 10240
	DefaultMaxRecordSizeInKiB = 1024

	// DefaultMaxRecordSizeBytes is the record-size bound applied when a
	// stream carries no max record size: 1 MiB.
	DefaultMaxRecordSizeBytes = 1048576

	// Wire window shared by the paginated inputs' Limit and MaxResults
	// members (model range traits on ListStreamsInputLimit,
	// ListShardsInputLimit, ListStreamConsumersInputLimit and
	// GetRecordsInputLimit).
	MinListResultsLimit = 1
	MaxListResultsLimit = 10000

	// DefaultGetRecordsLimit is the GetRecords documented default page.
	DefaultGetRecordsLimit = 10000

	// MaxPutRecordsEntries is the PutRecords batch size cap (model @length
	// max on Records).
	MaxPutRecordsEntries = 500

	// MaxPutRecordsRequestBytes is the PutRecords whole-request cap (API
	// reference, PutRecords doc): each record can be as large as the
	// per-record bound, "up to a limit of 10 MiB for the entire request,
	// including partition keys" — the sum of every entry's decoded payload
	// and partition key.
	MaxPutRecordsRequestBytes = 10 * 1024 * 1024

	// MaxShardCount is the shard-count upper bound (model ShardCountObject
	// range 0-1000000).
	MaxShardCount = 1000000

	// Account quotas reported by DescribeLimits (API reference quotas).
	ShardLimit               = 500
	OnDemandStreamCountLimit = 50

	// DescribeLimitsStreamsWalkBound is the per-request ceiling of the
	// stream walk the DescribeLimits aggregation performs; the walk pages
	// through the shared paginator's hard item cap until every stream is
	// counted, so the cap bounds one page, never the total.
	DescribeLimitsStreamsWalkBound = 10000

	// DescribeStreamShardPageCeiling is the effective shard page of
	// DescribeStream: the Limit member's documentation fixes the default
	// at one hundred and serves at most one hundred shards per call —
	// "If you specify a value greater than 100, at most 100 results are
	// returned" — even though the member's range trait accepts 1-10000.
	DescribeStreamShardPageCeiling = 100

	// DefaultOnDemandShardCount: the on-demand creation default. No current
	// AWS source states an initial shard count — the model and API
	// reference are silent, and the quotas page expresses the on-demand
	// default as 4 MB/s write throughput, which four shards deliver under
	// the same page's 1 MB/s-per-shard write rate.
	DefaultOnDemandShardCount = 4

	// SubscribePollPageSize is the record read bound per SubscribeToShard
	// pump cycle.
	SubscribePollPageSize = 1000

	// NextTokenValidity is the documented lifetime of a list-operation
	// continuation token: "Tokens expire after 300 seconds. When you obtain
	// a value for NextToken... you have 300 seconds to use that value. If
	// you specify an expired token... you get ExpiredNextTokenException."
	// (ListShards and ListStreamConsumers member documentation.)
	NextTokenValidity = 300 * time.Second

	// ATTimestampScanBound caps the record scan an AT_TIMESTAMP iterator
	// positioning walk examines: a target position deeper than this many
	// records past the shard's starting sequence number clamps to the last
	// record the walk reached. A platform bound — AWS leaves the walk's
	// depth undocumented and indexes the position directly.
	ATTimestampScanBound = 10000

	// MaxTagPageResults is the ListTagsForStreamInputLimit range-trait
	// ceiling (the model bounds the member 1-50); an absent Limit serves
	// one full page of this size.
	MaxTagPageResults = 50

	// MaxTagsPerResource is the documented per-resource tag total: "You
	// can assign up to 50 tags to a data stream" (AddTagsToStream), "You
	// can add up to 50 tags per resource" (the Tags members) — the bound
	// applies to the merged key set, not to each write.
	MaxTagsPerResource = 50

	// MaxTagKeysPerUntagRequest is the TagKeyList length-trait ceiling: at
	// most fifty tag keys per RemoveTagsFromStream or UntagResource call.
	MaxTagKeysPerUntagRequest = 50

	// SubscribeTTL is the documented subscription lifetime:
	// SubscribeToShard serves events "for up to 5 minutes, after which
	// time you need to call SubscribeToShard again to renew the
	// subscription if you want to continue to receive records".
	SubscribeTTL = 5 * time.Minute

	// SubscribeTakeoverWindow is the documented re-subscription window:
	// a repeat call with the same ConsumerARN and ShardId within five
	// seconds of a successful call answers ResourceInUseException; at
	// five seconds or more the new call takes the subscription over and
	// the previous connection expires.
	SubscribeTakeoverWindow = 5 * time.Second

	// MaxConsumersPerStream is the documented per-stream registered-consumer
	// quota for enhanced fan-out: twenty for On-demand Standard and
	// Provisioned streams, fifty for On-demand Advantage streams. The
	// platform carries no Advantage tier, so twenty applies to every stream
	// mode it implements.
	MaxConsumersPerStream = 20
)

// MaxShardHashKeyInt returns the shard hash-key space bound as a fresh
// big.Int; callers may mutate the result.
func MaxShardHashKeyInt() *big.Int {
	value, ok := new(big.Int).SetString(MaxShardHashKey, 10)
	if !ok {
		panic("kinesis: MaxShardHashKey is not a decimal constant")
	}
	return value
}

// RetentionCutoff is the arrival timestamp beyond which a record has aged
// out of the stream's retention window — the single derivation every read
// path's visibility filter and the physical trim share.
func RetentionCutoff(retentionPeriodHours int32) time.Time {
	return time.Now().UTC().Add(-time.Duration(retentionPeriodHours) * time.Hour)
}
