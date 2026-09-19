package kinesis

import (
	"math"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// parseStreamModeDetails reads the StreamModeDetails member. The second
// return reports whether the member travelled at all: CreateStream's member
// is optional (an absent one defaults to PROVISIONED in the Core) while
// UpdateStreamMode's is required, so the two cores treat absence
// differently. A present details object without its own StreamMode member
// reads as an empty mode — the inner member is required, and the cores
// reject it.
func parseStreamModeDetails(params map[string]interface{}) (kinesisstore.StreamMode, bool) {
	streamModeDetails := request.GetMapParam(params, "StreamModeDetails")
	if streamModeDetails == nil {
		streamModeDetails = request.GetMapParam(params, "streamModeDetails")
	}
	if streamModeDetails == nil {
		return "", false
	}
	if v, ok := streamModeDetails["StreamMode"].(string); ok {
		return kinesisstore.StreamMode(v), true
	}
	if v, ok := streamModeDetails["streamMode"].(string); ok {
		return kinesisstore.StreamMode(v), true
	}
	return "", true
}

// formatEpochSeconds renders a timestamp as the model's epoch-seconds
// double with fractional precision — sub-second arrival order survives.
func formatEpochSeconds(t time.Time) float64 {
	return float64(t.UnixNano()) / 1e9
}

// formatConsumer converts a Consumer to its API response map. The Consumer
// shape carries the consumer's own members only — the stream reference
// lives on ConsumerDescription (DescribeStreamConsumer), a different model
// shape with its own formatter.
func formatConsumer(c *kinesisstore.Consumer) map[string]interface{} {
	return map[string]interface{}{
		"ConsumerName":              c.ConsumerName,
		"ConsumerARN":               c.ConsumerARN,
		"ConsumerStatus":            c.ConsumerStatus,
		"ConsumerCreationTimestamp": formatEpochSeconds(c.ConsumerCreationTimestamp),
	}
}

// formatConsumerDescription converts a Consumer to the
// DescribeStreamConsumer response map — the ConsumerDescription shape,
// which adds the owning stream's ARN.
func formatConsumerDescription(c *kinesisstore.Consumer) map[string]interface{} {
	m := formatConsumer(c)
	m["StreamARN"] = c.StreamARN
	return m
}

// strictIntParam reads a typed Integer member. An omitted member keeps the
// operation's documented default (present=false); a member that is present
// but cannot be read as an integer — a fractional or int32-overflowing JSON
// number, or a non-numeric string — is a wire-type violation rejected as
// InvalidArgumentException, the error shape every Kinesis operation
// declares. The lenient readers truncate float64 values, so typed members
// must come through here.
func strictIntParam(params map[string]interface{}, key string) (int, bool, error) {
	value, present, err := request.GetIntParamStrictCaseInsensitive(params, key)
	if err != nil {
		return 0, present, ErrInvalidArgument
	}
	return value, present, nil
}

// strictStringParam reads a typed String member. An omitted member stays
// empty (present=false); a member that is present but is not a string — the
// JSON wire form of a String or blob member is a string, so a number or
// object is a wire-type violation — is rejected as InvalidArgumentException,
// never coerced to the empty string the lenient readers return. A present
// null reads as absent: the awsJson1_1 conformance model drops null
// structure values, the reading the other strict getters and HasParam take.
func strictStringParam(params map[string]interface{}, key string) (string, bool, error) {
	for _, k := range []string{key, request.LowerFirst(key), strings.ToLower(key)} {
		if v, ok := params[k]; ok {
			if v == nil {
				return "", false, nil
			}
			if s, ok := v.(string); ok {
				return s, true, nil
			}
			return "", true, ErrInvalidArgument
		}
	}
	return "", false, nil
}

// strictTimestampParam reads a typed Timestamp member. An omitted member
// stays empty (present=false); a member that is present but cannot be read
// as a timestamp — the SDKs serialise these members as epoch-second JSON
// numbers, and both documented string notations are accepted — is a
// wire-type violation rejected as InvalidArgumentException, never a
// silently dropped position that degrades to horizon behaviour.
func strictTimestampParam(params map[string]interface{}, key string) (string, bool, error) {
	value, present, err := request.GetTimestampParamStrict(params, key)
	if err != nil {
		return "", present, ErrInvalidArgument
	}
	return value, present, nil
}

// maxConvertibleEpochSeconds is 2^63 — the first epoch-second value the
// int64 narrowing in parseTimestampMember cannot carry (float64 rounds
// math.MaxInt64 up to exactly this bound).
const maxConvertibleEpochSeconds = 9223372036854775808.0

// maxUsableEpochSeconds is the highest epoch-second value a time.Time can
// carry: the internal seconds field stores sec + unixToInternal
// (62135596800, the seconds from year 1 to 1970), so a value beyond
// math.MaxInt64 minus that offset overflows the field and the instant's
// calendar arithmetic breaks. The low end has no such bound — adding the
// offset to math.MinInt64 moves the sum toward zero.
const maxUsableEpochSeconds = int64(math.MaxInt64 - 62135596800)

// parseTimestampMember reads a Timestamp member from its canonical wire
// string: epoch seconds with fractional seconds honoured, or — the second
// value form the member's documentation shows — RFC 3339 notation. A member
// that is present but unparseable, non-finite, or beyond the seconds range
// time.Time's internal field carries is a wire-type violation answered
// with InvalidArgumentException — the int64 narrowing outside its range is
// implementation-dependent in Go, and a value past the representation's
// usable ceiling is never a working instant.
func parseTimestampMember(value string) (time.Time, error) {
	seconds, err := strconv.ParseFloat(value, 64)
	if err != nil {
		if t, terr := time.Parse(time.RFC3339, value); terr == nil {
			return t.UTC(), nil
		}
		return time.Time{}, ErrInvalidArgument
	}
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) ||
		seconds >= maxConvertibleEpochSeconds || seconds < -maxConvertibleEpochSeconds {
		return time.Time{}, ErrInvalidArgument
	}
	sec := int64(seconds)
	if sec > maxUsableEpochSeconds {
		return time.Time{}, ErrInvalidArgument
	}
	nsec := int64((seconds - float64(sec)) * 1e9)
	return time.Unix(sec, nsec).UTC(), nil
}

// resolveEncryptionType returns the encryption type string, defaulting to "NONE".
func resolveEncryptionType(stream *kinesisstore.Stream) string {
	if stream.EncryptionType != "" {
		return stream.EncryptionType
	}
	return "NONE"
}

// verifyStreamGeneration applies the StreamCreationTimestamp identity
// check the list operations document for the member: it must parse (a
// wire violation answers InvalidArgumentException) and it must name the
// resolved stream's own generation — a deleted+recreated stream name
// answers ResourceNotFoundException. Callers that validated the member
// early for the strict-read precedence re-derive it here, where the
// resolved stream is available to compare against.
func (s *KinesisService) verifyStreamGeneration(member string, stream *kinesisstore.Stream) error {
	if member == "" {
		return nil
	}
	ts, err := parseTimestampMember(member)
	if err != nil {
		return err
	}
	if stream.CreatedAt.Unix() != ts.Unix() {
		return s.mapStoreError(kinesisstore.ErrStreamNotFound)
	}
	return nil
}
