package testutil

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/smithy-go/middleware"
	"vorpalstacks-sdk-tests/config"
)

// The Live Tail pins (E5): the streaming session's wire contract through
// the AWS SDK's event-stream consumer and the session-request rejection
// rows that reach the server.

// newLiveTailClient builds a CloudWatch Logs client without host-prefix
// injection: the modelled "stream-" host prefix would rewrite the
// endpoint host away from the platform.
func (tc *cwlogsTestCtx) newLiveTailClient() (*cloudwatchlogs.Client, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: tc.runner.endpoint,
		Region:   tc.runner.region,
	})
	if err != nil {
		return nil, err
	}
	return cloudwatchlogs.NewFromConfig(cfg, func(o *cloudwatchlogs.Options) {
		o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
			return stack.Finalize.Add(&disableHostPrefixFinalize{}, middleware.Before)
		})
	}), nil
}

func (tc *cwlogsTestCtx) liveTailTests() []TestResult {
	var results []TestResult

	// One full session: the sessionStart event carries the session ids
	// and the resolved ARN identifiers; an event ingested after the
	// session opened streams in a per-second sessionUpdate while the
	// pre-session event never does; the filter pattern holds.
	results = append(results, tc.runner.RunTest("logs", "LiveTail_SessionStream", func() error {
		client, err := tc.newLiveTailClient()
		if err != nil {
			return err
		}
		group, cleanupGroup, err := tc.newGroupStreamFixture("livetail", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		preMark := "livetail-before-session"
		if err := tc.putLogEvent(group, "s1", preMark, time.Now().UnixMilli()); err != nil {
			return err
		}

		resp, err := client.StartLiveTail(tc.ctx, &cloudwatchlogs.StartLiveTailInput{
			LogGroupIdentifiers:   []string{tc.logGroupARN(group)},
			LogEventFilterPattern: aws.String("livetail-marker"),
		})
		if err != nil {
			return fmt.Errorf("start live tail: %v", err)
		}
		stream := resp.GetStream()
		defer stream.Close()

		const postMark = "livetail-marker streamed"
		started := false
		deadline := time.Now().Add(12 * time.Second)
		for time.Now().Before(deadline) {
			event, ok := <-stream.Events()
			if !ok {
				return fmt.Errorf("event stream closed early: %v", stream.Err())
			}
			switch ev := event.(type) {
			case *types.StartLiveTailResponseStreamMemberSessionStart:
				started = true
				if aws.ToString(ev.Value.RequestId) == "" || aws.ToString(ev.Value.SessionId) == "" {
					return fmt.Errorf("sessionStart ids: %q/%q", aws.ToString(ev.Value.RequestId), aws.ToString(ev.Value.SessionId))
				}
				if len(ev.Value.LogGroupIdentifiers) != 1 || !strings.HasSuffix(ev.Value.LogGroupIdentifiers[0], group) {
					return fmt.Errorf("sessionStart identifiers = %v, want the group ARN", ev.Value.LogGroupIdentifiers)
				}
				// The session is live: put the marked event and a
				// non-matching sibling now; the per-second updates carry
				// them within the window.
				if err := tc.putLogEvent(group, "s1", postMark, time.Now().UnixMilli()); err != nil {
					return err
				}
				if err := tc.putLogEvent(group, "s1", "unmarked stays out", time.Now().UnixMilli()); err != nil {
					return err
				}
			case *types.StartLiveTailResponseStreamMemberSessionUpdate:
				if !started {
					return fmt.Errorf("sessionUpdate arrived before sessionStart")
				}
				if ev.Value.SessionMetadata != nil && ev.Value.SessionMetadata.Sampled {
					return fmt.Errorf("sampled = true under the 500-event ceiling")
				}
				for _, result := range ev.Value.SessionResults {
					switch aws.ToString(result.Message) {
					case preMark:
						return fmt.Errorf("pre-session event streamed: the tail carries no replay")
					case postMark:
						if aws.ToString(result.LogStreamName) != "s1" || !strings.HasSuffix(aws.ToString(result.LogGroupIdentifier), group) {
							return fmt.Errorf("streamed event identifiers: %v", result)
						}
						return nil
					}
				}
			default:
				return fmt.Errorf("unexpected stream event %T", event)
			}
		}
		return fmt.Errorf("the marked event never streamed within the window")
	}))

	// The session-request rows the SDK passes through: the unknown group
	// is ResourceNotFound, and the stream-filter pairing rules reject.
	results = append(results, tc.runner.RunTest("logs", "LiveTail_RejectionRows", func() error {
		client, err := tc.newLiveTailClient()
		if err != nil {
			return err
		}
		group, cleanupGroup, err := tc.newLogGroupFixture("livetail-reject")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		other, cleanupOther, err := tc.newLogGroupFixture("livetail-reject-2")
		if err != nil {
			return err
		}
		defer cleanupOther()

		_, err = client.StartLiveTail(tc.ctx, &cloudwatchlogs.StartLiveTailInput{
			LogGroupIdentifiers: []string{tc.logGroupARN(tc.uniquePrefix("livetail-no-such"))},
		})
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("unknown group = %v, want ResourceNotFoundException", err)
		}

		// A bare name is not the member's documented addressing form
		// ("Specify each log group by its ARN").
		_, err = client.StartLiveTail(tc.ctx, &cloudwatchlogs.StartLiveTailInput{
			LogGroupIdentifiers: []string{group},
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("bare group name: %v", err)
		}

		_, err = client.StartLiveTail(tc.ctx, &cloudwatchlogs.StartLiveTailInput{
			LogGroupIdentifiers:   []string{tc.logGroupARN(group)},
			LogStreamNames:        []string{"s1"},
			LogStreamNamePrefixes: []string{"s"},
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("both stream filters: %v", err)
		}

		_, err = client.StartLiveTail(tc.ctx, &cloudwatchlogs.StartLiveTailInput{
			LogGroupIdentifiers: []string{tc.logGroupARN(group), tc.logGroupARN(other)},
			LogStreamNames:      []string{"s1"},
		})
		return AssertErrorContains(err, "InvalidParameterException")
	}))

	return results
}
