package testutil

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func (r *TestRunner) runSNSTopicTests(tc *snsTestContext) []TestResult {
	var results []TestResult

	topicName := tc.uniqueName("TestTopic")
	var topicArn string

	results = append(results, r.RunTest("sns", "CreateTopic", func() error {
		resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(topicName),
		})
		if err != nil {
			return err
		}
		if resp.TopicArn == nil {
			return fmt.Errorf("TopicArn is nil")
		}
		topicArn = *resp.TopicArn
		if !strings.Contains(topicArn, topicName) {
			return fmt.Errorf("TopicArn should contain topic name %q, got %q", topicName, topicArn)
		}

		getResp, err := tc.getTopicAttributes(topicArn)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes == nil {
			return fmt.Errorf("Attributes is nil")
		}
		if getResp.Attributes["TopicArn"] != topicArn {
			return fmt.Errorf("TopicArn mismatch: got %q, want %q", getResp.Attributes["TopicArn"], topicArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetTopicAttributes_GetVerify", func() error {
		attrTopicArn, err := tc.createTopic(tc.uniqueName("AttrTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(attrTopicArn)

		_, err = tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn:       aws.String(attrTopicArn),
			AttributeName:  aws.String("DisplayName"),
			AttributeValue: aws.String("MyDisplayName"),
		})
		if err != nil {
			return fmt.Errorf("set: %v", err)
		}

		getResp, err := tc.getTopicAttributes(attrTopicArn)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes == nil {
			return fmt.Errorf("Attributes is nil")
		}
		if getResp.Attributes["DisplayName"] != "MyDisplayName" {
			return fmt.Errorf("DisplayName mismatch: got %q, want %q", getResp.Attributes["DisplayName"], "MyDisplayName")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "DeleteTopic", func() error {
		delTopicArn, err := tc.createTopic(tc.uniqueName("DelTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		_, err = tc.client.DeleteTopic(tc.ctx, &sns.DeleteTopicInput{
			TopicArn: aws.String(delTopicArn),
		})
		if err != nil {
			return err
		}
		_, err = tc.getTopicAttributes(delTopicArn)
		return tc.expectNotFound("GetTopicAttributes", err)
	}))

	results = append(results, r.RunTest("sns", "CreateTopic_DuplicateIdempotent", func() error {
		dupTopicName := tc.uniqueName("DupTopic")
		arn1, err := tc.createTopic(dupTopicName)
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.deleteTopic(arn1)

		arn2, err := tc.createTopic(dupTopicName)
		if err != nil {
			return fmt.Errorf("duplicate create should be idempotent, got: %v", err)
		}
		if arn1 != arn2 {
			return fmt.Errorf("ARN mismatch: %q vs %q", arn1, arn2)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "CreateTopic_FIFO", func() error {
		fifoTopicName := tc.uniqueName("TestFifoTopic") + ".fifo"
		resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoTopicName),
			Attributes: map[string]string{
				"ContentBasedDeduplication": "true",
				"FifoTopic":                 "true",
			},
		})
		if err != nil {
			return err
		}
		if resp.TopicArn == nil {
			return fmt.Errorf("TopicArn is nil")
		}
		// Deferred before the assertions: a failed assertion must not
		// orphan the FIFO topic.
		defer tc.deleteTopic(*resp.TopicArn)
		if !strings.HasSuffix(*resp.TopicArn, ".fifo") {
			return fmt.Errorf("FIFO topic ARN should end in .fifo, got %q", *resp.TopicArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetTopicAttributes_FifoAttributes", func() error {
		fifoAttrName := tc.uniqueName("FifoAttrTopic") + ".fifo"
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(fifoAttrName),
			Attributes: map[string]string{
				"FifoTopic":                 "true",
				"ContentBasedDeduplication": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		getResp, err := tc.getTopicAttributes(*tResp.TopicArn)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes["FifoTopic"] != "true" {
			return fmt.Errorf("FifoTopic should be true, got %q", getResp.Attributes["FifoTopic"])
		}
		if getResp.Attributes["ContentBasedDeduplication"] != "true" {
			return fmt.Errorf("ContentBasedDeduplication should be true, got %q", getResp.Attributes["ContentBasedDeduplication"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetTopicAttributes_PolicyDefault", func() error {
		policyTopicArn, err := tc.createTopic(tc.uniqueName("PolicyTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(policyTopicArn)

		getResp, err := tc.getTopicAttributes(policyTopicArn)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		policy, ok := getResp.Attributes["Policy"]
		if !ok || policy == "" {
			return fmt.Errorf("default Policy should be set")
		}
		if !strings.Contains(policy, "Version") {
			return fmt.Errorf("policy should contain Version, got: %s", policy)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListTopics_ContainsCreated", func() error {
		ltTopicArn, err := tc.createTopic(tc.uniqueName("LTTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(ltTopicArn)

		arns, err := tc.listAllTopics()
		if err != nil {
			return err
		}
		if len(arns) == 0 {
			return fmt.Errorf("ListTopics returned no topics")
		}
		if !slices.Contains(arns, ltTopicArn) {
			return fmt.Errorf("created topic not found in ListTopics")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListTopics_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", nanoTime())
		var pgTopicARNs []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagTopic-%s-%d", pgTs, i)
			arn, err := tc.createTopic(name)
			if err != nil {
				for _, created := range pgTopicARNs {
					tc.deleteTopic(created)
				}
				return fmt.Errorf("create topic %s: %v", name, err)
			}
			pgTopicARNs = append(pgTopicARNs, arn)
		}

		allARNs, err := tc.listAllTopics()
		if err != nil {
			for _, created := range pgTopicARNs {
				tc.deleteTopic(created)
			}
			return fmt.Errorf("list topics: %v", err)
		}
		var allTopics []string
		for _, arn := range allARNs {
			if strings.Contains(arn, "PagTopic-"+pgTs) {
				allTopics = append(allTopics, arn)
			}
		}

		for _, created := range pgTopicARNs {
			tc.deleteTopic(created)
		}
		if len(allTopics) != 5 {
			return fmt.Errorf("expected 5 paginated topics, got %d", len(allTopics))
		}
		return nil
	}))

	// GetTopicAttributes follows the documented enumeration: bookkeeping
	// timestamps are not attributes, the optional attributes appear only
	// when explicitly set, and EffectiveDeliveryPolicy is synthesised from
	// the stored policy and the system defaults.
	results = append(results, r.RunTest("sns", "TopicAttributes_DocumentedSet", func() error {
		bareArn, err := tc.createTopic(tc.uniqueName("DocSetTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(bareArn)

		getResp, err := tc.getTopicAttributes(bareArn)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		for _, absent := range []string{"CreatedDate", "LastModifiedTime", "DisplayName", "SignatureVersion", "MaximumMessageSize"} {
			if _, present := getResp.Attributes[absent]; present {
				return fmt.Errorf("bare topic returns %s — the documented set returns it only when explicitly set", absent)
			}
		}
		for _, present := range []string{"TopicArn", "Owner", "Policy", "EffectiveDeliveryPolicy", "SubscriptionsConfirmed", "SubscriptionsDeleted", "SubscriptionsPending"} {
			if _, ok := getResp.Attributes[present]; !ok {
				return fmt.Errorf("bare topic omits the documented key %s", present)
			}
		}
		if !strings.Contains(getResp.Attributes["EffectiveDeliveryPolicy"], `"numRetries":3`) {
			return fmt.Errorf("default EffectiveDeliveryPolicy = %q, want the system default of three retries", getResp.Attributes["EffectiveDeliveryPolicy"])
		}
		return nil
	}))

	// SignatureVersion accepts 1 and 2 only, echoes back once set, and —
	// absence means 1 — is not emitted for a topic that never set it.
	results = append(results, r.RunTest("sns", "SetTopicAttributes_SignatureVersion", func() error {
		sigArn, err := tc.createTopic(tc.uniqueName("SigVerTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(sigArn)

		if _, err := tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn: aws.String(sigArn), AttributeName: aws.String("SignatureVersion"), AttributeValue: aws.String("3"),
		}); expectAWSErrorCode(err, "InvalidParameter") != nil {
			return fmt.Errorf("SignatureVersion 3 accepted: %v", err)
		}

		if _, err := tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn: aws.String(sigArn), AttributeName: aws.String("SignatureVersion"), AttributeValue: aws.String("2"),
		}); err != nil {
			return fmt.Errorf("set SignatureVersion 2: %v", err)
		}
		getResp, err := tc.getTopicAttributes(sigArn)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes["SignatureVersion"] != "2" {
			return fmt.Errorf("SignatureVersion = %q, want 2", getResp.Attributes["SignatureVersion"])
		}
		return nil
	}))

	// MaximumMessageSize is honoured: a topic set to 1024 bytes rejects a
	// larger publish, and out-of-range values are rejected at set time (the
	// documented range floor is 1024; this platform's ceiling is its flat
	// 256 KiB transport cap).
	results = append(results, r.RunTest("sns", "SetTopicAttributes_MaximumMessageSize", func() error {
		capArn, err := tc.createTopic(tc.uniqueName("MaxSizeTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(capArn)

		if _, err := tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn: aws.String(capArn), AttributeName: aws.String("MaximumMessageSize"), AttributeValue: aws.String("1023"),
		}); expectAWSErrorCode(err, "InvalidParameter") != nil {
			return fmt.Errorf("MaximumMessageSize below the 1024 floor accepted: %v", err)
		}

		if _, err := tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn: aws.String(capArn), AttributeName: aws.String("MaximumMessageSize"), AttributeValue: aws.String("1024"),
		}); err != nil {
			return fmt.Errorf("set MaximumMessageSize 1024: %v", err)
		}
		getResp, err := tc.getTopicAttributes(capArn)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes["MaximumMessageSize"] != "1024" {
			return fmt.Errorf("MaximumMessageSize = %q, want 1024", getResp.Attributes["MaximumMessageSize"])
		}

		if _, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
			TopicArn: aws.String(capArn), Message: aws.String(strings.Repeat("a", 1025)),
		}); expectAWSErrorCode(err, "InvalidParameter") != nil {
			return fmt.Errorf("publish above the topic's MaximumMessageSize accepted: %v", err)
		}
		return nil
	}))

	// The write plane refuses the keys it must: FifoTopic is create-only,
	// the FIFO-only attributes are rejected on standard topics, and the
	// synthesised and derived keys are read-only.
	results = append(results, r.RunTest("sns", "SetTopicAttributes_ImmutableAndReadOnlyRejected", func() error {
		stdArn, err := tc.createTopic(tc.uniqueName("RejectStdTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(stdArn)
		fifoResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name:       aws.String(tc.uniqueName("RejectFifoTopic") + ".fifo"),
			Attributes: map[string]string{"FifoTopic": "true"},
		})
		if err != nil {
			return fmt.Errorf("create FIFO topic: %v", err)
		}
		defer tc.deleteTopic(*fifoResp.TopicArn)

		for name, value := range map[string]string{
			"FifoTopic":                 "true",
			"EffectiveDeliveryPolicy":   "{}",
			"TopicArn":                  stdArn,
			"SubscriptionsConfirmed":    "0",
			"ContentBasedDeduplication": "true",
			"FifoThroughputScope":       "MessageGroup",
		} {
			if _, err := tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
				TopicArn: aws.String(stdArn), AttributeName: aws.String(name), AttributeValue: aws.String(value),
			}); expectAWSErrorCode(err, "InvalidParameter") != nil {
				return fmt.Errorf("standard topic accepted %s: %v", name, err)
			}
		}
		if _, err := tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn: fifoResp.TopicArn, AttributeName: aws.String("FifoTopic"), AttributeValue: aws.String("true"),
		}); expectAWSErrorCode(err, "InvalidParameter") != nil {
			return fmt.Errorf("FIFO topic accepted the create-only FifoTopic via SetTopicAttributes: %v", err)
		}
		return nil
	}))

	// FifoThroughputScope=MessageGroup scopes the topic's deduplication
	// window to each message group: the same deduplication ID in two
	// message groups delivers both copies, where the default Topic scope
	// admits only the first. The subscribed queue is a STANDARD queue: a
	// FIFO queue would deduplicate the second copy again by the forwarded
	// deduplication ID, hiding the topic-side scope behind the queue's own
	// documented window.
	results = append(results, r.RunTest("sns", "FifoTopic_PerGroupDedupScope", func() error {
		topicResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(tc.uniqueName("ScopeFifoTopic") + ".fifo"),
			Attributes: map[string]string{
				"FifoTopic":           "true",
				"FifoThroughputScope": "MessageGroup",
			},
		})
		if err != nil {
			return fmt.Errorf("create scoped FIFO topic: %v", err)
		}
		defer tc.deleteTopic(*topicResp.TopicArn)

		sqsClient, err := tc.sqsClient()
		if err != nil {
			return err
		}
		queueURL, cleanupQueue, err := createTestQueue(tc.ctx, sqsClient, tc.uniqueName("ScopeStdQueue"), nil)
		if err != nil {
			return fmt.Errorf("create delivery queue: %w", err)
		}
		defer cleanupQueue()

		queueARN, err := queueArn(tc.ctx, sqsClient, queueURL)
		if err != nil {
			return fmt.Errorf("queue ARN: %w", err)
		}
		subResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: topicResp.TopicArn,
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(queueARN),
		})
		if err != nil {
			return fmt.Errorf("subscribe: %w", err)
		}
		defer tc.client.Unsubscribe(tc.ctx, &sns.UnsubscribeInput{SubscriptionArn: subResp.SubscriptionArn})

		for i, group := range []string{"scope-group-a", "scope-group-b"} {
			if _, err := tc.client.Publish(tc.ctx, &sns.PublishInput{
				TopicArn:               topicResp.TopicArn,
				Message:                aws.String(fmt.Sprintf("scoped-%d", i)),
				MessageGroupId:         aws.String(group),
				MessageDeduplicationId: aws.String("shared-scope-dedup"),
			}); err != nil {
				return fmt.Errorf("publish to %s: %w", group, err)
			}
		}

		var bodies []string
		deadline := time.Now().Add(20 * time.Second)
		for len(bodies) < 2 && time.Now().Before(deadline) {
			resp, err := sqsClient.ReceiveMessage(tc.ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            queueURL,
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     1,
			})
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			for _, m := range resp.Messages {
				var envelope struct {
					Message string `json:"Message"`
				}
				if err := json.Unmarshal([]byte(aws.ToString(m.Body)), &envelope); err != nil {
					return fmt.Errorf("delivered body is not the notification envelope: %v", err)
				}
				bodies = append(bodies, envelope.Message)
			}
		}
		if !slices.Contains(bodies, "scoped-0") || !slices.Contains(bodies, "scoped-1") {
			return fmt.Errorf("per-group dedup scope delivered %v, want both message groups' copies", bodies)
		}
		return nil
	}))

	if topicArn != "" {
		tc.deleteTopic(topicArn)
	}

	return results
}
