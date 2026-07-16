package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
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
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListTopics", func() error {
		resp, err := tc.client.ListTopics(tc.ctx, &sns.ListTopicsInput{})
		if err != nil {
			return err
		}
		if resp.Topics == nil {
			return fmt.Errorf("Topics is nil")
		}
		found := false
		for _, t := range resp.Topics {
			if t.TopicArn != nil && *t.TopicArn == topicArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created topic %q not found in ListTopics", topicArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetTopicAttributes", func() error {
		resp, err := tc.client.GetTopicAttributes(tc.ctx, &sns.GetTopicAttributesInput{
			TopicArn: aws.String(topicArn),
		})
		if err != nil {
			return err
		}
		if resp.Attributes == nil {
			return fmt.Errorf("Attributes is nil")
		}
		if resp.Attributes["TopicArn"] != topicArn {
			return fmt.Errorf("TopicArn mismatch: got %q, want %q", resp.Attributes["TopicArn"], topicArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetTopicAttributes", func() error {
		_, err := tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn:       aws.String(topicArn),
			AttributeName:  aws.String("DisplayName"),
			AttributeValue: aws.String("Test Topic"),
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "SetTopicAttributes_GetVerify", func() error {
		attrTopicName := tc.uniqueName("AttrTopic")
		resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(attrTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*resp.TopicArn)

		_, err = tc.client.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn:       resp.TopicArn,
			AttributeName:  aws.String("DisplayName"),
			AttributeValue: aws.String("MyDisplayName"),
		})
		if err != nil {
			return fmt.Errorf("set: %v", err)
		}

		getResp, err := tc.client.GetTopicAttributes(tc.ctx, &sns.GetTopicAttributesInput{
			TopicArn: resp.TopicArn,
		})
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

	var delTopicArn string
	results = append(results, r.RunTest("sns", "DeleteTopic", func() error {
		delTopicName := tc.uniqueName("DelTopic")
		resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(delTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		delTopicArn = *resp.TopicArn
		_, err = tc.client.DeleteTopic(tc.ctx, &sns.DeleteTopicInput{
			TopicArn: resp.TopicArn,
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "DeleteTopic_VerifyGone", func() error {
		_, err := tc.client.GetTopicAttributes(tc.ctx, &sns.GetTopicAttributesInput{
			TopicArn: aws.String(delTopicArn),
		})
		return AssertErrorContains(err, "NotFound")
	}))

	results = append(results, r.RunTest("sns", "CreateTopic_DuplicateIdempotent", func() error {
		dupTopicName := tc.uniqueName("DupTopic")
		resp1, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(dupTopicName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.deleteTopic(*resp1.TopicArn)

		resp2, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(dupTopicName),
		})
		if err != nil {
			return fmt.Errorf("duplicate create should be idempotent, got: %v", err)
		}
		if *resp1.TopicArn != *resp2.TopicArn {
			return fmt.Errorf("ARN mismatch: %q vs %q", *resp1.TopicArn, *resp2.TopicArn)
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
		if !strings.HasSuffix(*resp.TopicArn, ".fifo") {
			return fmt.Errorf("FIFO topic ARN should end with .fifo, got %q", *resp.TopicArn)
		}
		tc.deleteTopic(*resp.TopicArn)
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

		getResp, err := tc.client.GetTopicAttributes(tc.ctx, &sns.GetTopicAttributesInput{
			TopicArn: tResp.TopicArn,
		})
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
		policyTopicName := tc.uniqueName("PolicyTopic")
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(policyTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		getResp, err := tc.client.GetTopicAttributes(tc.ctx, &sns.GetTopicAttributesInput{
			TopicArn: tResp.TopicArn,
		})
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
		ltTopicName := tc.uniqueName("LTTopic")
		resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name: aws.String(ltTopicName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*resp.TopicArn)

		listResp, err := tc.client.ListTopics(tc.ctx, &sns.ListTopicsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, t := range listResp.Topics {
			if t.TopicArn != nil && *t.TopicArn == *resp.TopicArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created topic not found in ListTopics")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListTopics_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", nanoTime())
		var pgTopicARNs []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagTopic-%s-%d", pgTs, i)
			resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{Name: aws.String(name)})
			if err != nil {
				for _, arn := range pgTopicARNs {
					tc.deleteTopic(arn)
				}
				return fmt.Errorf("create topic %s: %v", name, err)
			}
			pgTopicARNs = append(pgTopicARNs, *resp.TopicArn)
		}

		var allTopics []string
		var nextToken *string
		for {
			resp, err := tc.client.ListTopics(tc.ctx, &sns.ListTopicsInput{NextToken: nextToken})
			if err != nil {
				for _, arn := range pgTopicARNs {
					tc.deleteTopic(arn)
				}
				return fmt.Errorf("list topics page: %v", err)
			}
			for _, t := range resp.Topics {
				if strings.Contains(aws.ToString(t.TopicArn), "PagTopic-"+pgTs) {
					allTopics = append(allTopics, aws.ToString(t.TopicArn))
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, arn := range pgTopicARNs {
			tc.deleteTopic(arn)
		}
		if len(allTopics) != 5 {
			return fmt.Errorf("expected 5 paginated topics, got %d", len(allTopics))
		}
		return nil
	}))

	if topicArn != "" {
		tc.deleteTopic(topicArn)
	}

	return results
}
