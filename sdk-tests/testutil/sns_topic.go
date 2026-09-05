package testutil

import (
	"fmt"
	"slices"
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

	if topicArn != "" {
		tc.deleteTopic(topicArn)
	}

	return results
}
