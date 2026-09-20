package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runSNSPolicyTests(tc *snsTestContext) []TestResult {
	var results []TestResult
	reg := tc.region
	acct := tc.accountID

	results = append(results, r.RunTest("sns", "AddPermission", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("PermTopic"))
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		_, err = tc.client.AddPermission(tc.ctx, &sns.AddPermissionInput{
			TopicArn:     aws.String(topicArn),
			Label:        aws.String("TestPermission"),
			AWSAccountId: []string{acct},
			ActionName:   []string{"Publish"},
		})
		if err != nil {
			return err
		}

		getResp, err := tc.getTopicAttributes(topicArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		policy := getResp.Attributes["Policy"]
		if policy == "" {
			return fmt.Errorf("Policy should be set after AddPermission")
		}
		return nil
	}))

	// Label is "A unique identifier for the new policy statement" — a
	// duplicate label is refused rather than overwriting the live
	// statement.
	results = append(results, r.RunTest("sns", "AddPermission_DuplicateLabelRejected", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("DupLabelTopic"))
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		addInput := &sns.AddPermissionInput{
			TopicArn:     aws.String(topicArn),
			Label:        aws.String("UniqueLabel"),
			AWSAccountId: []string{acct},
			ActionName:   []string{"Publish"},
		}
		if _, err := tc.client.AddPermission(tc.ctx, addInput); err != nil {
			return fmt.Errorf("first AddPermission: %v", err)
		}
		_, err = tc.client.AddPermission(tc.ctx, addInput)
		if codeErr := expectAWSErrorCode(err, "InvalidParameter"); codeErr != nil {
			return fmt.Errorf("duplicate label: %v", codeErr)
		}

		// The refused duplicate left the original statement intact.
		getResp, err := tc.getTopicAttributes(topicArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		policy := getResp.Attributes["Policy"]
		if len(policy) == 0 {
			return fmt.Errorf("Policy missing after the refused duplicate")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "RemovePermission", func() error {
		topicArn, err := tc.createTopic(tc.uniqueName("RemPermTopic"))
		if err != nil {
			return err
		}
		defer tc.deleteTopic(topicArn)

		_, err = tc.client.AddPermission(tc.ctx, &sns.AddPermissionInput{
			TopicArn:     aws.String(topicArn),
			Label:        aws.String("TestPermission"),
			AWSAccountId: []string{acct},
			ActionName:   []string{"Publish"},
		})
		if err != nil {
			return fmt.Errorf("add permission: %v", err)
		}

		_, err = tc.client.RemovePermission(tc.ctx, &sns.RemovePermissionInput{
			TopicArn: aws.String(topicArn),
			Label:    aws.String("TestPermission"),
		})
		return err
	}))

	results = append(results, r.RunTest("sns", "PutDataProtectionPolicy", func() error {
		dppTopicArn, err := tc.createTopic(tc.uniqueName("DppTopic"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(dppTopicArn)

		policy := `{"Name":"data-protection-policy","Description":"Test policy","Version":"2021-06-01"}`
		_, err = tc.client.PutDataProtectionPolicy(tc.ctx, &sns.PutDataProtectionPolicyInput{
			ResourceArn:          aws.String(dppTopicArn),
			DataProtectionPolicy: aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		getResp, err := tc.client.GetDataProtectionPolicy(tc.ctx, &sns.GetDataProtectionPolicyInput{
			ResourceArn: aws.String(dppTopicArn),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.DataProtectionPolicy == nil {
			return fmt.Errorf("DataProtectionPolicy is nil")
		}
		if *getResp.DataProtectionPolicy != policy {
			return fmt.Errorf("policy mismatch: got %q, want %q", *getResp.DataProtectionPolicy, policy)
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "GetDataProtectionPolicy_NonExistent", func() error {
		_, err := tc.client.GetDataProtectionPolicy(tc.ctx, &sns.GetDataProtectionPolicyInput{
			ResourceArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:nonexistent-dpp-topic", reg, acct)),
		})
		return tc.expectNotFound("GetDataProtectionPolicy", err)
	}))

	// CreateTopic accepts the inline DataProtectionPolicy parameter; the
	// policy is retrievable via GetDataProtectionPolicy and excluded from
	// GetTopicAttributes output.
	results = append(results, r.RunTest("sns", "CreateTopic_InlineDataProtectionPolicy", func() error {
		policy := `{"Name":"inline-policy","Version":"2021-06-01"}`
		tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name:                 aws.String(tc.uniqueName("InlineDppTopic")),
			DataProtectionPolicy: aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteTopic(*tResp.TopicArn)

		getResp, err := tc.client.GetDataProtectionPolicy(tc.ctx, &sns.GetDataProtectionPolicyInput{
			ResourceArn: tResp.TopicArn,
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.DataProtectionPolicy == nil || *getResp.DataProtectionPolicy != policy {
			return fmt.Errorf("policy mismatch: got %v, want %q", getResp.DataProtectionPolicy, policy)
		}

		attrsResp, err := tc.getTopicAttributes(*tResp.TopicArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if _, ok := attrsResp.Attributes["DataProtectionPolicy"]; ok {
			return fmt.Errorf("DataProtectionPolicy must not appear in GetTopicAttributes")
		}
		return nil
	}))

	// CreateTopic with an invalid JSON DataProtectionPolicy must be rejected.
	results = append(results, r.RunTest("sns", "CreateTopic_InlineDataProtectionPolicy_InvalidJSON_Rejected", func() error {
		_, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
			Name:                 aws.String(tc.uniqueName("BadDppTopic")),
			DataProtectionPolicy: aws.String("not-json"),
		})
		return expectAWSErrorCode(err, "InvalidParameter")
	}))

	return results
}
