package testutil

import (
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
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
		if *resp.SubscriptionArn != "pending confirmation" {
			tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: resp.SubscriptionArn})
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
		if getResp.Attributes["Protocol"] != "sqs" {
			return fmt.Errorf("Protocol mismatch: got %q, want %q", getResp.Attributes["Protocol"], "sqs")
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

	results = append(results, r.RunTest("sns", "Subscribe_ApplicationPendingConfirmation", func() error {
		appTopicArn, err := tc.createTopic(tc.uniqueName("AppTopic"))
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(appTopicArn)

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(appTopicArn),
			Protocol: aws.String("application"),
			Endpoint: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:app/FAKE/fake-endpoint", reg, acct)),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		if sResp.SubscriptionArn == nil || *sResp.SubscriptionArn == "" {
			return fmt.Errorf("application subscription should return ARN")
		}
		if *sResp.SubscriptionArn == "pending confirmation" {
			return fmt.Errorf("application protocol should be auto-confirmed, got pending confirmation")
		}

		getResp, err := tc.getSubscriptionAttributes(*sResp.SubscriptionArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes["PendingConfirmation"] != "false" {
			return fmt.Errorf("application subscription should be auto-confirmed, got PendingConfirmation=%s", getResp.Attributes["PendingConfirmation"])
		}
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

	return results
}
