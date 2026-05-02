package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runSNSSubscriptionTests(tc *snsTestContext) []TestResult {
	var results []TestResult

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
		sqsTopicName := tc.uniqueName("SqsTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(sqsTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		sqsSubTopicArn = *tResp.TopicArn

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("arn:aws:sqs:us-east-1:000000000000:auto-confirm-queue"),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		sqsSubArn = *sResp.SubscriptionArn

		getResp, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
		})
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes["PendingConfirmation"] != "false" {
			return fmt.Errorf("SQS subscription should be auto-confirmed, got PendingConfirmation=%s", getResp.Attributes["PendingConfirmation"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetSubscriptionAttributes", func() error {
		resp, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(sqsSubArn),
		})
		if err != nil {
			return err
		}
		if resp.Attributes == nil {
			return fmt.Errorf("Attributes is nil")
		}
		if resp.Attributes["SubscriptionArn"] != sqsSubArn {
			return fmt.Errorf("SubscriptionArn mismatch: got %q", resp.Attributes["SubscriptionArn"])
		}
		if resp.Attributes["Protocol"] != "sqs" {
			return fmt.Errorf("Protocol mismatch: got %q, want %q", resp.Attributes["Protocol"], "sqs")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetSubscriptionAttributes", func() error {
		_, err := tc.client.SetSubscriptionAttributes(tc.ctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(sqsSubArn),
			AttributeName:   aws.String("RawMessageDelivery"),
			AttributeValue:  aws.String("true"),
		})
		if err != nil {
			return fmt.Errorf("set: %v", err)
		}

		getResp, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(sqsSubArn),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes["RawMessageDelivery"] != "true" {
			return fmt.Errorf("RawMessageDelivery not set: got %q", getResp.Attributes["RawMessageDelivery"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptions", func() error {
		resp, err := tc.client.ListSubscriptions(tc.ctx, &sns.ListSubscriptionsInput{})
		if err != nil {
			return err
		}
		if resp.Subscriptions == nil {
			return fmt.Errorf("Subscriptions is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptions_ContainsCreated", func() error {
		lsTopicName := tc.uniqueName("LSTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(lsTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("arn:aws:sqs:us-east-1:000000000000:list-all-sub-queue"),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		listResp, err := tc.client.ListSubscriptions(tc.ctx, &sns.ListSubscriptionsInput{})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		found := false
		for _, sub := range listResp.Subscriptions {
			if sub.SubscriptionArn != nil && *sub.SubscriptionArn == *sResp.SubscriptionArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created subscription not found in ListSubscriptions")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptionsByTopic", func() error {
		lstTopicName := tc.uniqueName("LstSubTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(lstTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		_, err = tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("arn:aws:sqs:us-east-1:000000000000:list-sub-queue"),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		listResp, err := tc.client.ListSubscriptionsByTopic(tc.ctx, &sns.ListSubscriptionsByTopicInput{
			TopicArn: tResp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("list by topic: %v", err)
		}
		if len(listResp.Subscriptions) == 0 {
			return fmt.Errorf("expected at least one subscription")
		}
		found := false
		for _, sub := range listResp.Subscriptions {
			if sub.TopicArn != nil && *sub.TopicArn == *tResp.TopicArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("subscription for topic not found in ListSubscriptionsByTopic")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListSubscriptionsByTopic_Empty", func() error {
		emptySubTopicName := tc.uniqueName("EmptySubTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(emptySubTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		listResp, err := tc.client.ListSubscriptionsByTopic(tc.ctx, &sns.ListSubscriptionsByTopicInput{
			TopicArn: tResp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		if len(listResp.Subscriptions) != 0 {
			return fmt.Errorf("expected 0 subscriptions for new topic, got %d", len(listResp.Subscriptions))
		}
		return nil
	}))

	var unsubSubArn string
	results = append(results, r.RunTest("sns", "Unsubscribe", func() error {
		unsubTopicName := tc.uniqueName("UnsubTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(unsubTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("sqs"),
			Endpoint: aws.String("arn:aws:sqs:us-east-1:000000000000:fake-queue"),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}
		unsubSubArn = *sResp.SubscriptionArn

		_, err = tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{
			SubscriptionArn: sResp.SubscriptionArn,
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "Unsubscribe_VerifyGone", func() error {
		_, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(unsubSubArn),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "Subscribe_EmailPendingConfirmation", func() error {
		emailTopicName := tc.uniqueName("EmailTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(emailTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(*tResp.TopicArn),
			Protocol:              aws.String("email"),
			Endpoint:              aws.String("pending@example.com"),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		getResp, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
		})
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes["PendingConfirmation"] != "true" {
			return fmt.Errorf("email subscription should be pending, got PendingConfirmation=%s", getResp.Attributes["PendingConfirmation"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ConfirmSubscription", func() error {
		confTopicName := tc.uniqueName("ConfTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(confTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn:              aws.String(*tResp.TopicArn),
			Protocol:              aws.String("email"),
			Endpoint:              aws.String("confirm-test@example.com"),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			return fmt.Errorf("subscribe: %v", err)
		}

		getResp, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
		})
		if err != nil {
			return fmt.Errorf("get attrs before confirm: %v", err)
		}
		if getResp.Attributes["PendingConfirmation"] != "true" {
			return fmt.Errorf("should be pending before confirm, got %s", getResp.Attributes["PendingConfirmation"])
		}
		token := getResp.Attributes["Token"]
		if token == "" {
			return fmt.Errorf("Token attribute should be present for pending subscription")
		}

		confResp, err := tc.client.ConfirmSubscription(tc.ctx, &sns.ConfirmSubscriptionInput{
			TopicArn: tResp.TopicArn,
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

		afterResp, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
		})
		if err != nil {
			return fmt.Errorf("get attrs after confirm: %v", err)
		}
		if afterResp.Attributes["PendingConfirmation"] != "false" {
			return fmt.Errorf("should be confirmed, got PendingConfirmation=%s", afterResp.Attributes["PendingConfirmation"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "Subscribe_ApplicationPendingConfirmation", func() error {
		appTopicName := tc.uniqueName("AppTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(appTopicName),
		})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: aws.String(*tResp.TopicArn),
			Protocol: aws.String("application"),
			Endpoint: aws.String("arn:aws:sns:us-east-1:000000000000:app/FAKE/fake-endpoint"),
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

		getResp, err := tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
			SubscriptionArn: sResp.SubscriptionArn,
		})
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

	return results
}
