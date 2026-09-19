package kinesis

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// The MaxResults input window and the effective page follow two different
// documented rules: the ListShardsInputLimit range trait accepts 1-10000 on
// the wire, while the MaxResults member documentation fixes the default
// page at 1000 and returns at most 1000 results per call.
const (
	DefaultListShardsResults = 1000
	MaxListShardsResults     = 1000
)

// ListShardsResult carries a page of shards and its continuation state:
// NextToken is the opaque envelope carrying the stream identification and
// the page's last shard for the follow-up call; a following page exists
// exactly when it is non-empty.
type ListShardsResult struct {
	Shards    []*kinesisstore.Shard
	NextToken string
}

// ListShardsInput is the transport-agnostic input for ListShards. NextToken
// is the raw wire member — the envelope parses in the Core so both planes
// share one token definition; ExclusiveStartShardId is the model's explicit
// resumption member.
type ListShardsInput struct {
	StreamName              string
	StreamARN               string
	StreamCreationTimestamp string
	ShardFilter             *kinesisstore.ShardFilter
	MaxResults              int
	HasMaxResults           bool
	NextToken               string
	ExclusiveStartShardId   string
}

// SplitShardInput is the transport-agnostic input for SplitShard.
type SplitShardInput struct {
	StreamName         string
	StreamARN          string
	ShardToSplit       string
	NewStartingHashKey string
}

// MergeShardsInput is the transport-agnostic input for MergeShards.
type MergeShardsInput struct {
	StreamName           string
	StreamARN            string
	ShardToMerge         string
	AdjacentShardToMerge string
}

// UpdateShardCountInput is the transport-agnostic input for UpdateShardCount.
type UpdateShardCountInput struct {
	StreamName       string
	StreamARN        string
	TargetShardCount int32
	ScalingType      string
}

// UpdateShardCountResult carries the reshaping receipt.
type UpdateShardCountResult struct {
	StreamName        string
	CurrentShardCount int32
	TargetShardCount  int32
	StreamARN         string
}

// ensureReshapableStream pins the documented preconditions the reshaping
// trio shares: the operations are supported on provisioned-capacity streams
// alone (an on-demand stream rejects), and the model declares
// ResourceInUseException for any stream status other than ACTIVE.
func ensureReshapableStream(stream *kinesisstore.Stream) error {
	if stream.StreamModeDetails != nil && stream.StreamModeDetails.StreamMode == kinesisstore.StreamModeOnDemand {
		return ErrInvalidArgument
	}
	if stream.StreamStatus != kinesisstore.StreamStatusActive {
		return ErrResourceInUse
	}
	return nil
}

// encodeListShardsToken packs the continuation state into the opaque
// envelope the model documents: the token "unambiguously identifies" the
// stream — its name and creation timestamp disambiguate a deleted and
// recreated stream — and carries the resumption position together with its
// own issuance time, the field the documented 300-second validity window
// is measured against.
func encodeListShardsToken(streamName string, createdAt, issuedAtNanos int64, lastShardID string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s#%d#%s#%d", streamName, createdAt, lastShardID, issuedAtNanos)))
}

// decodeListShardsToken unpacks the continuation envelope. Anything that
// does not carry the four fields is not one of this service's tokens.
func decodeListShardsToken(token string) (streamName string, createdAt, issuedAtNanos int64, lastShardID string, ok bool) {
	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", 0, 0, "", false
	}
	parts := strings.Split(string(raw), "#")
	if len(parts) != 4 {
		return "", 0, 0, "", false
	}
	created, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, 0, "", false
	}
	issued, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return "", 0, 0, "", false
	}
	return parts[0], created, issued, parts[2], true
}

// listShardsCore validates the MaxResults window, resolves the stream — by
// name or ARN, or through the continuation envelope that carries the
// identification itself — verifies the creation timestamp when provided,
// validates the shard filter against the model's enum and its member
// pairings, and pages with one over-fetch so the exact final page carries
// no continuation token.
func (s *KinesisService) listShardsCore(reqCtx *request.RequestContext, input ListShardsInput) (ListShardsResult, error) {
	// The range trait window is enforced on the provided value, then the
	// effective page applies the documented cap.
	limit := input.MaxResults
	if input.HasMaxResults {
		if !validateListShardsLimit(limit) {
			return ListShardsResult{}, ErrInvalidArgument
		}
		if limit > MaxListShardsResults {
			limit = MaxListShardsResults
		}
	} else {
		limit = DefaultListShardsResults
	}

	// The creation timestamp parses before storage is acquired; the value
	// is compared against the resolved stream below.
	if input.StreamCreationTimestamp != "" {
		if _, err := parseTimestampMember(input.StreamCreationTimestamp); err != nil {
			return ListShardsResult{}, err
		}
	}

	// The resumption member is mutually exclusive with the request's own
	// identification and positioning members: the model's NextToken
	// documentation states StreamName and StreamCreationTimestamp cannot
	// travel with it ("the latter unambiguously identifies the stream"),
	// and ExclusiveStartShardId's own documentation states the same — an
	// agreeing value is still a member the contract forbids.
	if input.NextToken != "" && (input.StreamName != "" || input.StreamCreationTimestamp != "" || input.ExclusiveStartShardId != "") {
		return ListShardsResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return ListShardsResult{}, err
	}

	// The token carries the stream identification; a StreamARN travelling
	// with it is addressing, not identification, so it must agree with the
	// stream the token names.
	resumeShardID := ""
	var stream *kinesisstore.Stream
	if input.NextToken != "" {
		tokenStream, tokenCreated, tokenIssued, tokenShardID, ok := decodeListShardsToken(input.NextToken)
		if !ok {
			return ListShardsResult{}, ErrExpiredNextToken
		}
		// The documented validity window: a token older than it answers
		// ExpiredNextTokenException whatever it names.
		if time.Since(time.Unix(0, tokenIssued)) > kinesisstore.NextTokenValidity {
			return ListShardsResult{}, ErrExpiredNextToken
		}
		if input.StreamARN != "" {
			arnStream, err := store.GetStreamByARN(input.StreamARN)
			if err != nil {
				return ListShardsResult{}, s.mapStoreError(err)
			}
			if arnStream.StreamName != tokenStream {
				return ListShardsResult{}, ErrInvalidArgument
			}
		}
		resumeShardID = tokenShardID
		stream, err = store.GetStream(tokenStream)
		if err != nil {
			return ListShardsResult{}, s.mapStoreError(err)
		}
		// A stream deleted and recreated under the same name is a
		// different stream: the token's creation timestamp identifies its
		// own, now-expired, generation.
		if stream.CreatedAt.Unix() != tokenCreated {
			return ListShardsResult{}, ErrExpiredNextToken
		}
	} else {
		streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
		if err != nil {
			return ListShardsResult{}, err
		}
		stream, err = store.GetStream(streamName)
		if err != nil {
			return ListShardsResult{}, s.mapStoreError(err)
		}
		resumeShardID = input.ExclusiveStartShardId
	}

	// Verify StreamCreationTimestamp matches when provided (used to
	// disambiguate deleted+recreated streams); the parse already happened
	// before storage was acquired, and the identity check itself is the
	// shared generation rule.
	if err := s.verifyStreamGeneration(input.StreamCreationTimestamp, stream); err != nil {
		return ListShardsResult{}, err
	}

	// The filter's Type is a required property of the ShardFilter and takes
	// the model's enum; the paired members travel only with their types and
	// their types require them. An absent ShardFilter member is the
	// documented default type.
	filter := input.ShardFilter
	if filter == nil {
		filter = &kinesisstore.ShardFilter{Type: "FROM_TRIM_HORIZON"}
	} else if err := validateShardFilter(filter); err != nil {
		return ListShardsResult{}, err
	}

	// One over-fetch past the effective page: a full result set answers
	// hasMore exactly, so the final page — even one that fills the page
	// exactly — carries no continuation token.
	shards, err := store.ListShards(stream.StreamName, filter, resumeShardID, limit+1)
	if err != nil {
		return ListShardsResult{}, s.mapStoreError(err)
	}
	hasMore := len(shards) > limit
	if hasMore {
		shards = shards[:limit]
	}
	nextToken := ""
	if hasMore && len(shards) > 0 {
		nextToken = encodeListShardsToken(stream.StreamName, stream.CreatedAt.Unix(), time.Now().UnixNano(), shards[len(shards)-1].ShardID)
	}

	return ListShardsResult{Shards: shards, NextToken: nextToken}, nil
}

// splitShardCore splits a shard at the given starting hash key.
func (s *KinesisService) splitShardCore(reqCtx *request.RequestContext, input SplitShardInput) (string, error) {
	// NewStartingHashKey is a required member; the shared validator treats
	// an empty member as unset, so the presence check comes first.
	if input.ShardToSplit == "" || !validateShardId(input.ShardToSplit) || input.NewStartingHashKey == "" {
		return "", ErrInvalidArgument
	}
	if !validateExplicitHashKey(input.NewStartingHashKey) {
		return "", ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return "", err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return "", s.mapStoreError(err)
	}
	if err := ensureReshapableStream(stream); err != nil {
		return "", err
	}

	// The new key must fall inside the parent shard's hash key range: the
	// member documentation requires it, and the store lays the children out
	// as [start, key-1] and [key, end] — a key at the starting boundary
	// inverts the lower child. A shard's hash key range is fixed at
	// creation, so this check cannot race a concurrent reshape.
	parent, err := store.GetShard(streamName, input.ShardToSplit)
	if err != nil {
		return "", s.mapStoreError(err)
	}
	if parent.HashKeyRange == nil || !hashKeyWithinRange(input.NewStartingHashKey, parent.HashKeyRange) {
		return "", ErrInvalidArgument
	}

	if err := store.SplitShard(streamName, input.ShardToSplit, input.NewStartingHashKey); err != nil {
		return "", s.mapStoreError(err)
	}

	return store.BuildStreamARN(streamName), nil
}

// mergeShardsCore merges two adjacent shards.
func (s *KinesisService) mergeShardsCore(reqCtx *request.RequestContext, input MergeShardsInput) (string, error) {
	if input.ShardToMerge == "" || input.AdjacentShardToMerge == "" || !validateShardId(input.ShardToMerge) || !validateShardId(input.AdjacentShardToMerge) {
		return "", ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return "", err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return "", s.mapStoreError(err)
	}
	if err := ensureReshapableStream(stream); err != nil {
		return "", err
	}

	if err := store.MergeShards(streamName, input.ShardToMerge, input.AdjacentShardToMerge); err != nil {
		return "", s.mapStoreError(err)
	}

	return store.BuildStreamARN(streamName), nil
}

// updateShardCountCore reshapes the stream to the target shard count.
func (s *KinesisService) updateShardCountCore(reqCtx *request.RequestContext, input UpdateShardCountInput) (UpdateShardCountResult, error) {
	if !validateShardCount(input.TargetShardCount) {
		return UpdateShardCountResult{}, ErrInvalidArgument
	}

	// ScalingType is a required member whose enum carries the single value
	// UNIFORM_SCALING — the SDK's own required-member validation masks this
	// for typed clients, so the empty member reaches the Core only through
	// untyped transports.
	if input.ScalingType != "UNIFORM_SCALING" {
		return UpdateShardCountResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return UpdateShardCountResult{}, err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return UpdateShardCountResult{}, err
	}

	// CurrentShardCount reports the count at request time, before the
	// reshaping runs — the documented example answers the pre-update count
	// alongside the target (scaling a three-shard stream to six returns
	// CurrentShardCount 3, TargetShardCount 6).
	pre, err := store.GetStream(streamName)
	if err != nil {
		return UpdateShardCountResult{}, s.mapStoreError(err)
	}
	if err := ensureReshapableStream(pre); err != nil {
		return UpdateShardCountResult{}, err
	}

	if err := store.UpdateShardCount(streamName, input.TargetShardCount); err != nil {
		return UpdateShardCountResult{}, s.mapStoreError(err)
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return UpdateShardCountResult{}, s.mapStoreError(err)
	}

	return UpdateShardCountResult{
		StreamName:        streamName,
		CurrentShardCount: pre.ShardCount,
		TargetShardCount:  input.TargetShardCount,
		StreamARN:         stream.StreamARN,
	}, nil
}
