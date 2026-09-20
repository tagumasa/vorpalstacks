package testutil

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func (r *TestRunner) runSNSSubscriptionTests(tc *snsTestContext) []TestResult {
	var results []TestResult
	reg := tc.region
	acct := tc.accountID

	results = append(results, r.RunTest("sns", "Subscribe", func() error {
		topicName := tc.uniqueName("SubTopic")
		topicArn, err := tc.createTopic(topicName)
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		resp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("email"),
			Endpoint: aws.String("test@example.com"),
		})
		if err != nil {
			return err
		}
		if resp.SubscriptionArn == nil {
			return fmt.Errorf("SubscriptionArn is nil")
		}
		// Email subscriptions never auto-confirm on this platform (email
		// delivery is a recorded exclusion), so the response must carry the
		// pending marker — a confirmable ARN here would mean the protocol
		// auto-confirmed.
		if *resp.SubscriptionArn != "pending confirmation" {
			tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: resp.SubscriptionArn})
			return fmt.Errorf("email Subscribe returned %q, want the pending-confirmation marker", *resp.SubscriptionArn)
		}
		return nil
	}))

	var sqsSubArn string
	var sqsSubTopicArn string
	results = append(results, r.RunTest("sns", "Subscribe_SQS_AutoConfirmed", func() error {
		var err error
		sqsSubTopicArn, err = tc.createTopic(tc.uniqueName("SqsTopic"))
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(sqsSubTopicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:auto-confirm-queue", reg, acct)),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		sqsSubArn = *sResp.SubscriptionArn

		getResp, err := tc.getSubscriptionAttributes(sqsSubArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes["PendingConfirmation"] != "false" {
			return fmt.Errorf("SQS subscription should be auto-confirmed, got PendingConfirmation=%s", getResp.Attributes["PendingConfirmation"])
		}
		if getResp.Attributes["SubscriptionArn"] != sqsSubArn {
			return fmt.Errorf("SubscriptionArn mismatch: got %q", getResp.Attributes["SubscriptionArn"])
		}
		// Protocol and Endpoint are not attribute-plane keys — the documented
		// GetSubscriptionAttributes enumeration carries neither; the
		// subscription's protocol reaches clients through ListSubscriptions.
		if _, present := getResp.Attributes["Protocol"]; present {
			return fmt.Errorf("subscription attributes carry the invented key Protocol")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptions_ContainsCreated", func() error {
		lsTopicArn, err := tc.createTopic(tc.uniqueName("LSTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(lsTopicArn)

		subArn, err := tc.subscribeSQS(lsTopicArn, "list-all-sub-queue")
		if err != nil {
			return err
		}

		subs, err := tc.listAllSubscriptions()
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		if len(subs) == 0 {
			return fmt.Errorf("Subscriptions is empty")
		}
		found := slices.ContainsFunc(subs, func(s types.Subscription) bool {
			return aws.ToString(s.SubscriptionArn) == subArn
		})
		if !found {
			return fmt.Errorf("created subscription not found in ListSubscriptions")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptionsByTopic", func() error {
		lstTopicArn, err := tc.createTopic(tc.uniqueName("LstSubTopic"))
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(lstTopicArn)

		if _, err := tc.subscribeSQS(lstTopicArn, "list-sub-queue"); err != nil {
			return err
		}

		subs, err := tc.listAllSubscriptionsByTopic(lstTopicArn)
		if err != nil {
			return fmt.Errorf("list by topic: %v", err)
		}
		if len(subs) == 0 {
			return fmt.Errorf("expected at least one subscription")
		}
		found := slices.ContainsFunc(subs, func(s types.Subscription) bool {
			return aws.ToString(s.TopicArn) == lstTopicArn
		})
		if !found {
			return fmt.Errorf("subscription for topic not found in ListSubscriptionsByTopic")
		}

		// A topic ARN addressing no existing topic is rejected with the
		// model's awsQueryError code for NotFoundException ("NotFound"), not
		// the shape name.
		missingArn := fmt.Sprintf("arn:aws:sns:%s:%s:%s", tc.region, tc.accountID, tc.uniqueName("NoSuchTopic"))
		_, err = tc.listAllSubscriptionsByTopic(missingArn)
		if err == nil {
			return fmt.Errorf("ListSubscriptionsByTopic on a missing topic: expected an error, got none")
		}
		if nfe := tc.expectNotFound("ListSubscriptionsByTopic", err); nfe != nil {
			return nfe
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptionsByTopic_Empty", func() error {
		emptySubTopicArn, err := tc.createTopic(tc.uniqueName("EmptySubTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(emptySubTopicArn)

		subs, err := tc.listAllSubscriptionsByTopic(emptySubTopicArn)
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		if len(subs) != 0 {
			return fmt.Errorf("expected 0 subscriptions for new topic, got %d", len(subs))
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Unsubscribe", func() error {
		unsubTopicArn, err := tc.createTopic(tc.uniqueName("UnsubTopic"))
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(unsubTopicArn)

		subArn, err := tc.subscribeSQS(unsubTopicArn, "fake-queue")
		if err != nil {
			return err
		}

		_, err = tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{
			SubscriptionArn: aws.String(subArn),
		})
		if err != nil {
			return err
		}
		_, err = tc.getSubscriptionAttributes(subArn)
		return tc.expectNotFound("GetSubscriptionAttributes", err)
	}))

	results = append(results, r.RunTest("sns", "Subscribe_EmailPendingConfirmation", func() error {
		emailTopicArn, err := tc.createTopic(tc.uniqueName("EmailTopic"))
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(emailTopicArn)

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(emailTopicArn),
			Protocol:              aws.String("email"),
			Endpoint:              aws.String("pending@example.com"),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		getResp, err := tc.getSubscriptionAttributes(*sResp.SubscriptionArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes["PendingConfirmation"] != "true" {
			return fmt.Errorf("email subscription should be pending, got PendingConfirmation=%s", getResp.Attributes["PendingConfirmation"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ConfirmSubscription", func() error {
		confTopicArn, err := tc.createTopic(tc.uniqueName("ConfTopic"))
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(confTopicArn)

		// Set up a local HTTP server to capture the confirmation token.
		// AWS sends the token to the endpoint out-of-band; it is not
		// exposed through GetSubscriptionAttributes.
		srv, tokenCh := tc.startTokenCaptureServer()
		defer srv.Close()

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(confTopicArn),
			Protocol:              aws.String("http"),
			Endpoint:              aws.String(srv.URL),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		// Wait for the async confirmation message to arrive.
		var token string
		select {
		case token = <-tokenCh:
		case <-time.After(10 * time.Second):
			return fmt.Errorf("timed out waiting for subscription confirmation")
		}
		if token == "" {
			return fmt.Errorf("Token should not be empty from confirmation message")
		}

		confResp, err := tc.client.ConfirmSubscription(tc.ctx, &sns.ConfirmSubscriptionInput{
			TopicArn: aws.String(confTopicArn),
			Token:    aws.String(token),
		})
		if err != nil {
			return fmt.Errorf("confirm: %v", err)
		}
		if confResp.SubscriptionArn == nil || *confResp.SubscriptionArn == "" {
			return fmt.Errorf("confirmed SubscriptionArn should be non-empty")
		}
		if *confResp.SubscriptionArn != *sResp.SubscriptionArn {
			return fmt.Errorf("confirmed ARN mismatch: got %q, want %q", *confResp.SubscriptionArn, *sResp.SubscriptionArn)
		}

		afterResp, err := tc.getSubscriptionAttributes(*sResp.SubscriptionArn)
		if err != nil {
			return fmt.Errorf("get attrs after confirm: %v", err)
		}
		if afterResp.Attributes["PendingConfirmation"] != "false" {
			return fmt.Errorf("should be confirmed, got PendingConfirmation=%s", afterResp.Attributes["PendingConfirmation"])
		}
		if afterResp.Attributes["ConfirmationWasAuthenticated"] != "true" {
			return fmt.Errorf("ConfirmationWasAuthenticated = %q, want true (signed ConfirmSubscription call)",
				afterResp.Attributes["ConfirmationWasAuthenticated"])
		}
		if _, ok := afterResp.Attributes["AuthenticateOnUnsubscribe"]; ok {
			return fmt.Errorf("AuthenticateOnUnsubscribe must not be exposed in GetSubscriptionAttributes")
		}
		return nil
	}))

	// The application protocol cannot be delivered on this platform (its
	// delivery needs external push services), so Subscribe refuses it
	// instead of auto-confirming a subscription that would silently drop
	// every publish.
	results = append(results, r.RunTest("sns", "Subscribe_ApplicationProtocolRejected", func() error {
		appTopicArn, err := tc.createTopic(tc.uniqueName("AppTopic"))
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(appTopicArn)

		_, err = tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(appTopicArn),
			Protocol: aws.String("application"),
			Endpoint: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:app/FAKE/fake-endpoint", reg, acct)),
		})
		return expectAWSErrorCode(err, "InvalidParameter")
	}))

	// A repeated Subscribe of the same endpoint returns the live
	// subscription's ARN — duplicates never accumulate.
	results = append(results, r.RunTest("sns", "Subscribe_IdempotentSameEndpoint", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("IdemTopic"))
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		sqsClient, err := tc.sqsClient()
		if err != nil {
			return err
		}
		queueURL, cleanupQueue, err := createTestQueue(tc.ctx, sqsClient, tc.uniqueName("IdemQueue"), nil)
		if err != nil {
			return fmt.Errorf("create queue: %w", err)
		}
		defer cleanupQueue()
		queueARN, err := queueArn(tc.ctx, sqsClient, queueURL)
		if err != nil {
			return fmt.Errorf("queue ARN: %w", err)
		}

		first, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(queueARN),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: first.SubscriptionArn})

		second, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(queueARN),
		})
		if err != nil {
			return fmt.Errorf("re-subscribe: %v", err)
		}
		if aws.ToString(second.SubscriptionArn) != aws.ToString(first.SubscriptionArn) {
			return fmt.Errorf("re-subscribe returned %q after %q — duplicates accumulate",
				aws.ToString(second.SubscriptionArn), aws.ToString(first.SubscriptionArn))
		}

		subs, err := tc.listAllSubscriptionsByTopic(topicArn)
		if err != nil {
			return fmt.Errorf("list subscriptions: %v", err)
		}
		if len(subs) != 1 {
			return fmt.Errorf("topic carries %d subscriptions for one endpoint, want 1", len(subs))
		}
		return nil
	}))

	// FIFO topics refuse every customer-managed endpoint protocol — email,
	// mobile apps, SMS and HTTP(S) — while the AWS-managed protocols stay
	// subscribable.
	results = append(results, r.RunTest("sns", "FifoTopic_CustomerManagedEndpointRejected", func() error {
		fifoArn, err := tc.createTopic(tc.uniqueName("FifoProtoTopic") + ".fifo")
		if err != nil {
			return err
		}
		defer tc.deleteTopic(fifoArn)

		for _, tc2 := range []struct{ protocol, endpoint string }{
			{"http", "http://example.invalid/hook"},
			{"email", "subscriber@example.com"},
		} {
			_, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
				TopicArn: aws.String(fifoArn),
				Protocol: aws.String(tc2.protocol),
				Endpoint: aws.String(tc2.endpoint),
			})
			if codeErr := expectAWSErrorCode(err, "InvalidParameter"); codeErr != nil {
				return fmt.Errorf("protocol %s on a FIFO topic: %v", tc2.protocol, codeErr)
			}
		}

		sqsClient, err := tc.sqsClient()
		if err != nil {
			return err
		}
		queueURL, cleanupQueue, err := createTestQueue(tc.ctx, sqsClient, tc.uniqueName("FifoProtoQueue")+".fifo", map[string]string{
			"FifoQueue": "true",
		})
		if err != nil {
			return fmt.Errorf("create fifo queue: %w", err)
		}
		defer cleanupQueue()
		queueARN, err := queueArn(tc.ctx, sqsClient, queueURL)
		if err != nil {
			return fmt.Errorf("queue ARN: %w", err)
		}
		sub, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(fifoArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(queueARN),
		})
		if err != nil {
			return fmt.Errorf("sqs stays subscribable on a FIFO topic: %v", err)
		}
		tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: sub.SubscriptionArn})
		return nil
	}))

	if sqsSubTopicArn != "" {
		tc.deleteTopic(sqsSubTopicArn)
	}
	if sqsSubArn != "" {
		tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: aws.String(sqsSubArn)})
	}

	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes_FilterPolicy", func() error {
		topicName := tc.uniqueName("FilterTopic")
		topicArn, err := tc.createTopic(topicName)
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		// NOTE: The endpoint ARN points to a non-existent SQS queue. This test
		// verifies attribute round-trip (Set → Get) only, not actual delivery.
		subResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(topicArn),
			Protocol:              aws.String("sqs"),
			Endpoint:              aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:dummy-queue", reg, acct)),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		subArn := *subResp.SubscriptionArn
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: aws.String(subArn)})

		filterPolicy := `{"event": ["order_created", "order_updated"]}`
		_, err = tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(subArn),
			AttributeName:   aws.String("FilterPolicy"),
			AttributeValue:  aws.String(filterPolicy),
		})
		if err != nil {
			return fmt.Errorf("set FilterPolicy: %v", err)
		}

		getResp, err := tc.getSubscriptionAttributes(subArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes["FilterPolicy"] != filterPolicy {
			return fmt.Errorf("FilterPolicy mismatch: got %q, want %q", getResp.Attributes["FilterPolicy"], filterPolicy)
		}
		return nil
	}))

	// The documented operator grammar is accepted exactly: every operator
	// of the string/numeric matching pages sets cleanly, and operand
	// shapes the matcher cannot evaluate are rejected at set time.
	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes_FilterPolicyOperators", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("OperatorTopic"))
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		subArn, err := tc.subscribeSQS(topicArn, "operator-queue")
		if err != nil {
			return err
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: aws.String(subArn)})

		accepted := []string{
			`{"k": [{"prefix": "bas"}, {"suffix": "ball"}, {"wildcard": "*ball"}]}`,
			`{"k": [{"equals-ignore-case": "tennis"}]}`,
			`{"ip": [{"cidr": "10.0.0.0/24"}]}`,
			`{"n": [{"numeric": [">", 0, "<=", 150]}, {"numeric": ["=", 3.015e2]}]}`,
			`{"k": [{"anything-but": ["rugby", 100]}, {"anything-but": {"prefix": "order-"}}]}`,
			`{"k": ["v", 10, true, null, {"exists": false}]}`,
		}
		for _, policy := range accepted {
			if _, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
				SubscriptionArn: aws.String(subArn),
				AttributeName:   aws.String("FilterPolicy"),
				AttributeValue:  aws.String(policy),
			}); err != nil {
				return fmt.Errorf("documented policy %s rejected: %v", policy, err)
			}
		}

		rejected := []string{
			`{"k": [{"prefix": 123}]}`,
			`{"k": [{"exists": "yes"}]}`,
			`{"k": [{"numeric": [">", 0, "<"]}]}`,
			`{"k": [{"wildcard": "*a*b*c*d"}]}`,
			`{"k": [{"anything-but": {"numeric": [">", 1]}}]}`,
			`{"a": [1], "b": [1], "c": [1], "d": [1], "e": [1], "f": [1]}`,
		}
		for _, policy := range rejected {
			_, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
				SubscriptionArn: aws.String(subArn),
				AttributeName:   aws.String("FilterPolicy"),
				AttributeValue:  aws.String(policy),
			})
			if codeErr := expectAWSErrorCode(err, "InvalidParameter"); codeErr != nil {
				return fmt.Errorf("invalid policy %s: %v", policy, codeErr)
			}
		}
		return nil
	}))

	// Numeric matching normalises across text forms: a policy value of 10
	// delivers a message whose Number attribute carries "10.0", while a
	// non-equal value stays filtered.
	results = append(results, r.RunTest("sns", "FilterPolicy_NumericNormalisation", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("NumericTopic"))
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		sqsClient, err := tc.sqsClient()
		if err != nil {
			return err
		}
		queueURL, cleanupQueue, err := createTestQueue(tc.ctx, sqsClient, tc.uniqueName("NumericQueue"), nil)
		if err != nil {
			return fmt.Errorf("create queue: %w", err)
		}
		defer cleanupQueue()
		queueARN, err := queueArn(tc.ctx, sqsClient, queueURL)
		if err != nil {
			return fmt.Errorf("queue ARN: %w", err)
		}
		subResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(queueARN),
			Attributes: map[string]string{
				"FilterPolicy": `{"price": [{"numeric": ["=", 10]}]}`,
			},
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: subResp.SubscriptionArn})

		publish := func(value string) error {
			_, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
				TopicArn: aws.String(topicArn),
				Message:  aws.String("numeric-filter-probe"),
				MessageAttributes: map[string]types.MessageAttributeValue{
					"price": {DataType: aws.String("Number"), StringValue: aws.String(value)},
				},
			})
			return err
		}
		if err := publish("10.0"); err != nil {
			return fmt.Errorf("publish 10.0: %v", err)
		}
		if err := publish("10.5"); err != nil {
			return fmt.Errorf("publish 10.5: %v", err)
		}

		var bodies []string
		deadline := time.Now().Add(10 * time.Second)
		for time.Now().Before(deadline) {
			resp, err := sqsClient.ReceiveMessage(tc.ctx, &sqs.ReceiveMessageInput{
				QueueUrl:        queueURL,
				WaitTimeSeconds: 1,
			})
			if err != nil {
				return fmt.Errorf("receive: %v", err)
			}
			for _, m := range resp.Messages {
				bodies = append(bodies, aws.ToString(m.Body))
				if _, err := sqsClient.DeleteMessage(tc.ctx, &sqs.DeleteMessageInput{
					QueueUrl:      queueURL,
					ReceiptHandle: m.ReceiptHandle,
				}); err != nil {
					return fmt.Errorf("delete received copy: %v", err)
				}
			}
			// Both publishes share one body, so the count alone separates
			// the matching copy from a broken filter delivering both — the
			// negative half is only settled by an empty receive AFTER the
			// first copy, never by breaking on the first message.
			if len(bodies) > 0 && len(resp.Messages) == 0 {
				break
			}
		}
		if len(bodies) != 1 {
			return fmt.Errorf("numeric filter delivered %d messages, want exactly the 10.0 copy", len(bodies))
		}
		return nil
	}))

	// Nested filter policies are payload-based filtering: a nested policy is
	// refused without FilterPolicyScope MessageBody, delivers on the nested
	// body paths once scoped, and the scope cannot be flipped back while the
	// nested policy is attached.
	results = append(results, r.RunTest("sns", "FilterPolicy_NestedMessageBody", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("NestedTopic"))
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		sqsClient, err := tc.sqsClient()
		if err != nil {
			return err
		}
		queueURL, cleanupQueue, err := createTestQueue(tc.ctx, sqsClient, tc.uniqueName("NestedQueue"), nil)
		if err != nil {
			return fmt.Errorf("create queue: %w", err)
		}
		defer cleanupQueue()
		queueARN, err := queueArn(tc.ctx, sqsClient, queueURL)
		if err != nil {
			return fmt.Errorf("queue ARN: %w", err)
		}

		nested := `{"key_a": {"key_b": ["value_one"]}}`
		_, err = tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(queueARN),
			Attributes: map[string]string{
				"FilterPolicy": nested,
			},
		})
		if codeErr := expectAWSErrorCode(err, "InvalidParameter"); codeErr != nil {
			return fmt.Errorf("nested policy without FilterPolicyScope MessageBody: %v", codeErr)
		}

		subResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(queueARN),
			Attributes: map[string]string{
				"FilterPolicy":      nested,
				"FilterPolicyScope": "MessageBody",
			},
		})
		if err != nil {
			return fmt.Errorf("subscribe with nested policy and body scope: %v", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: subResp.SubscriptionArn})

		publish := func(body string) error {
			_, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
				TopicArn: aws.String(topicArn),
				Message:  aws.String(body),
			})
			return err
		}
		if err := publish(`{"key_a": {"key_b": "value_one"}, "probe": "nested-match"}`); err != nil {
			return fmt.Errorf("publish matching body: %v", err)
		}
		if err := publish(`{"key_a": {"key_b": "value_two"}, "probe": "nested-miss"}`); err != nil {
			return fmt.Errorf("publish non-matching body: %v", err)
		}

		var bodies []string
		deadline := time.Now().Add(10 * time.Second)
		for time.Now().Before(deadline) {
			resp, err := sqsClient.ReceiveMessage(tc.ctx, &sqs.ReceiveMessageInput{
				QueueUrl:        queueURL,
				WaitTimeSeconds: 1,
			})
			if err != nil {
				return fmt.Errorf("receive: %v", err)
			}
			for _, m := range resp.Messages {
				bodies = append(bodies, aws.ToString(m.Body))
				if _, err := sqsClient.DeleteMessage(tc.ctx, &sqs.DeleteMessageInput{
					QueueUrl:      queueURL,
					ReceiptHandle: m.ReceiptHandle,
				}); err != nil {
					return fmt.Errorf("delete received copy: %v", err)
				}
			}
			if len(bodies) >= 1 {
				break
			}
		}
		if len(bodies) != 1 {
			return fmt.Errorf("nested body filter delivered %d messages, want exactly the matching copy", len(bodies))
		}
		if !strings.Contains(bodies[0], "nested-match") {
			return fmt.Errorf("delivered copy is not the matching publish: %s", bodies[0])
		}

		_, err = tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: subResp.SubscriptionArn,
			AttributeName:   aws.String("FilterPolicyScope"),
			AttributeValue:  aws.String("MessageAttributes"),
		})
		if codeErr := expectAWSErrorCode(err, "InvalidParameter"); codeErr != nil {
			return fmt.Errorf("scope flip away from a nested policy: %v", codeErr)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes_RawMessageDelivery", func() error {
		topicName := tc.uniqueName("RawTopic")
		topicArn, err := tc.createTopic(topicName)
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		// NOTE: The endpoint ARN points to a non-existent SQS queue. This test
		// verifies attribute round-trip (Set → Get) only, not actual delivery.
		subResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(topicArn),
			Protocol:              aws.String("sqs"),
			Endpoint:              aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:dummy-queue", reg, acct)),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		subArn := *subResp.SubscriptionArn
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: aws.String(subArn)})

		_, err = tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(subArn),
			AttributeName:   aws.String("RawMessageDelivery"),
			AttributeValue:  aws.String("true"),
		})
		if err != nil {
			return fmt.Errorf("set RawMessageDelivery: %v", err)
		}

		getResp, err := tc.getSubscriptionAttributes(subArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes["RawMessageDelivery"] != "true" {
			return fmt.Errorf("RawMessageDelivery mismatch: got %q, want \"true\"", getResp.Attributes["RawMessageDelivery"])
		}
		return nil
	}))

	// SetSubscriptionAttributes must reject the reserved
	// AuthenticateOnUnsubscribe key (it is set via the ConfirmSubscription
	// input parameter only).
	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes_AuthenticateOnUnsubscribe_Rejected", func() error {
		sub, err := tc.createHTTPSubscription("AuthOnUnsub")
		if err != nil {
			return err
		}
		defer tc.deleteTopic(sub.TopicArn)

		_, err = tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(sub.SubscriptionArn),
			AttributeName:   aws.String("AuthenticateOnUnsubscribe"),
			AttributeValue:  aws.String("true"),
		})
		return expectAWSErrorCode(err, "InvalidParameter")
	}))

	// Unsubscribe of a subscription confirmed without
	// AuthenticateOnUnsubscribe must succeed normally (the flag is absent).
	results = append(results, r.RunTest("sns", "Unsubscribe_AfterPlainConfirm_Succeeds", func() error {
		sub, err := tc.createHTTPSubscription("PlainConfirm")
		if err != nil {
			return err
		}
		defer tc.deleteTopic(sub.TopicArn)

		_, err = tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{
			SubscriptionArn: aws.String(sub.SubscriptionArn),
		})
		return err
	}))

	// ConfirmSubscription must reject a non-boolean
	// AuthenticateOnUnsubscribe value instead of treating it as false.
	results = append(results, r.RunTest("sns", "ConfirmSubscription_AuthenticateOnUnsubscribe_InvalidValue_Rejected", func() error {
		badAuthTopicArn, err := tc.createTopic(tc.uniqueName("BadAuthOnUnsub"))
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(badAuthTopicArn)

		srv, tokenCh := tc.startTokenCaptureServer()
		defer srv.Close()

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(badAuthTopicArn),
			Protocol:              aws.String("http"),
			Endpoint:              aws.String(srv.URL),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		var token string
		select {
		case token = <-tokenCh:
		case <-time.After(10 * time.Second):
			return fmt.Errorf("timed out waiting for subscription confirmation")
		}

		_, err = tc.client.ConfirmSubscription(tc.ctx, &sns.ConfirmSubscriptionInput{
			TopicArn:                  aws.String(badAuthTopicArn),
			Token:                     aws.String(token),
			AuthenticateOnUnsubscribe: aws.String("xyz"),
		})
		if err := expectAWSErrorCode(err, "InvalidParameter"); err != nil {
			return err
		}

		// The subscription must remain unconfirmed and carry no flag.
		afterResp, err := tc.getSubscriptionAttributes(*sResp.SubscriptionArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if afterResp.Attributes["PendingConfirmation"] != "true" {
			return fmt.Errorf("subscription should remain pending after rejected confirm, got PendingConfirmation=%s",
				afterResp.Attributes["PendingConfirmation"])
		}
		return nil
	}))

	// GetSubscriptionAttributes follows the documented enumeration: the
	// invented Protocol/Endpoint keys are gone and the synthesised
	// EffectiveDeliveryPolicy is present.
	results = append(results, r.RunTest("sns", "GetSubscriptionAttributes_DocumentedSet", func() error {
		dsTopicArn, err := tc.createTopic(tc.uniqueName("DocSetSubTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(dsTopicArn)

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(dsTopicArn),
			Protocol:              aws.String("sqs"),
			Endpoint:              aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:doc-set-queue", reg, acct)),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: sResp.SubscriptionArn})

		getResp, err := tc.getSubscriptionAttributes(*sResp.SubscriptionArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		for _, invented := range []string{"Protocol", "Endpoint"} {
			if _, present := getResp.Attributes[invented]; present {
				return fmt.Errorf("subscription attributes carry the invented key %s", invented)
			}
		}
		for _, required := range []string{"SubscriptionArn", "TopicArn", "Owner", "PendingConfirmation", "ConfirmationWasAuthenticated", "EffectiveDeliveryPolicy"} {
			if _, ok := getResp.Attributes[required]; !ok {
				return fmt.Errorf("subscription attributes omit the documented key %s", required)
			}
		}
		return nil
	}))

	// A RedrivePolicy must name an existing dead-letter queue in the same
	// account and Region — a typo'd ARN is rejected at set time, not
	// discovered when the message is already lost.
	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes_RedrivePolicyQueueMustExist", func() error {
		rdTopicArn, err := tc.createTopic(tc.uniqueName("RedriveTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(rdTopicArn)

		sqsClient, err := tc.sqsClient()
		if err != nil {
			return err
		}
		dlqURL, cleanupQueue, err := createTestQueue(tc.ctx, sqsClient, tc.uniqueName("RedriveDLQ"), nil)
		if err != nil {
			return fmt.Errorf("create dead-letter queue: %w", err)
		}
		defer cleanupQueue()
		dlqARN, err := queueArn(tc.ctx, sqsClient, dlqURL)
		if err != nil {
			return fmt.Errorf("queue ARN: %w", err)
		}

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(rdTopicArn),
			Protocol:              aws.String("sqs"),
			Endpoint:              aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:redrive-delivery-queue", reg, acct)),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: sResp.SubscriptionArn})

		ghostPolicy := fmt.Sprintf(`{"deadLetterTargetArn":"arn:aws:sqs:%s:%s:no-such-dead-letter-queue"}`, reg, acct)
		if _, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
			AttributeName:   aws.String("RedrivePolicy"),
			AttributeValue:  aws.String(ghostPolicy),
		}); expectAWSErrorCode(err, "InvalidParameter") != nil {
			return fmt.Errorf("RedrivePolicy naming a missing queue accepted: %v", err)
		}

		livePolicy := fmt.Sprintf(`{"deadLetterTargetArn":"%s"}`, dlqARN)
		if _, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
			AttributeName:   aws.String("RedrivePolicy"),
			AttributeValue:  aws.String(livePolicy),
		}); err != nil {
			return fmt.Errorf("RedrivePolicy naming the live queue rejected: %v", err)
		}
		return nil
	}))

	// RawMessageDelivery takes boolean literals on the documented protocol
	// set (Amazon SQS and HTTP/S), and the read-only lifecycle keys are
	// refused on the write plane.
	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes_RawMessageDeliveryValidation", func() error {
		rawTopicArn, err := tc.createTopic(tc.uniqueName("RawValidTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(rawTopicArn)

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(rawTopicArn),
			Protocol:              aws.String("sqs"),
			Endpoint:              aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:raw-validation-queue", reg, acct)),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: sResp.SubscriptionArn})

		if _, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
			AttributeName:   aws.String("RawMessageDelivery"),
			AttributeValue:  aws.String("yes"),
		}); expectAWSErrorCode(err, "InvalidParameter") != nil {
			return fmt.Errorf("RawMessageDelivery \"yes\" accepted: %v", err)
		}

		if _, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
			AttributeName:   aws.String("PendingConfirmation"),
			AttributeValue:  aws.String("false"),
		}); expectAWSErrorCode(err, "InvalidParameter") != nil {
			return fmt.Errorf("the read-only PendingConfirmation accepted: %v", err)
		}

		// The documented raw-delivery protocol set is Amazon SQS and
		// HTTP/S — an email subscription must refuse the attribute.
		emailResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(rawTopicArn),
			Protocol:              aws.String("email"),
			Endpoint:              aws.String("raw-validation@example.com"),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("email subscribe: %v", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: emailResp.SubscriptionArn})
		if _, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: emailResp.SubscriptionArn,
			AttributeName:   aws.String("RawMessageDelivery"),
			AttributeValue:  aws.String("true"),
		}); expectAWSErrorCode(err, "InvalidParameter") != nil {
			return fmt.Errorf("RawMessageDelivery accepted on an email subscription: %v", err)
		}
		return nil
	}))

	return results
}
