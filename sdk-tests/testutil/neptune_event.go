package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

func (r *TestRunner) runNeptuneEventSubscriptionTests(tc *neptuneContext) []TestResult {
	var results []TestResult
	subName := fmt.Sprintf("test-sub-%d", tc.ts)

	results = append(results, r.RunTest("neptune", "CreateEventSubscription", func() error {
		resp, err := tc.client.CreateEventSubscription(tc.ctx, &neptune.CreateEventSubscriptionInput{
			SubscriptionName: aws.String(subName),
			SnsTopicArn:      aws.String("arn:aws:sns:us-east-1:000000000000:test-topic"),
			SourceType:       aws.String("db-cluster"),
			SourceIds:        []string{tc.clusterID},
			EventCategories:  []string{"creation", "deletion", "failover"},
		})
		if err != nil {
			return err
		}
		if resp.EventSubscription == nil {
			return fmt.Errorf("expected EventSubscription in response")
		}
		es := resp.EventSubscription
		if es.CustSubscriptionId == nil || *es.CustSubscriptionId != subName {
			return fmt.Errorf("expected CustSubscriptionId=%s, got %v", subName, es.CustSubscriptionId)
		}
		if es.SnsTopicArn == nil || *es.SnsTopicArn != "arn:aws:sns:us-east-1:000000000000:test-topic" {
			return fmt.Errorf("expected SnsTopicArn=arn:aws:sns:us-east-1:000000000000:test-topic, got %v", es.SnsTopicArn)
		}
		if es.SourceType == nil || *es.SourceType != "db-cluster" {
			return fmt.Errorf("expected SourceType=db-cluster, got %v", es.SourceType)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeEventSubscriptions", func() error {
		resp, err := tc.client.DescribeEventSubscriptions(tc.ctx, &neptune.DescribeEventSubscriptionsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, sub := range resp.EventSubscriptionsList {
			if sub.CustSubscriptionId != nil && *sub.CustSubscriptionId == subName {
				found = true
				if sub.SnsTopicArn == nil || *sub.SnsTopicArn != "arn:aws:sns:us-east-1:000000000000:test-topic" {
					return fmt.Errorf("expected SnsTopicArn in subscription, got %v", sub.SnsTopicArn)
				}
				if sub.Status == nil || *sub.Status == "" {
					return fmt.Errorf("expected non-empty Status in subscription")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created event subscription not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "AddSourceIdentifierToSubscription", func() error {
		resp, err := tc.client.AddSourceIdentifierToSubscription(tc.ctx, &neptune.AddSourceIdentifierToSubscriptionInput{
			SubscriptionName: aws.String(subName),
			SourceIdentifier: aws.String("cluster-extra"),
		})
		if err != nil {
			return err
		}
		if resp.EventSubscription == nil {
			return fmt.Errorf("expected EventSubscription in AddSourceIdentifierToSubscription response")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "AddSourceIdentifierToSubscription_Verify", func() error {
		resp, err := tc.client.DescribeEventSubscriptions(tc.ctx, &neptune.DescribeEventSubscriptionsInput{
			SubscriptionName: aws.String(subName),
		})
		if err != nil {
			return err
		}
		if len(resp.EventSubscriptionsList) != 1 {
			return fmt.Errorf("expected 1 subscription, got %d", len(resp.EventSubscriptionsList))
		}
		found := false
		for _, id := range resp.EventSubscriptionsList[0].SourceIdsList {
			if id == "cluster-extra" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected 'cluster-extra' in SourceIdsList after AddSourceIdentifierToSubscription")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "RemoveSourceIdentifierFromSubscription", func() error {
		resp, err := tc.client.RemoveSourceIdentifierFromSubscription(tc.ctx, &neptune.RemoveSourceIdentifierFromSubscriptionInput{
			SubscriptionName: aws.String(subName),
			SourceIdentifier: aws.String("cluster-extra"),
		})
		if err != nil {
			return err
		}
		if resp.EventSubscription == nil {
			return fmt.Errorf("expected EventSubscription in RemoveSourceIdentifierFromSubscription response")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "RemoveSourceIdentifierFromSubscription_Verify", func() error {
		resp, err := tc.client.DescribeEventSubscriptions(tc.ctx, &neptune.DescribeEventSubscriptionsInput{
			SubscriptionName: aws.String(subName),
		})
		if err != nil {
			return err
		}
		if len(resp.EventSubscriptionsList) != 1 {
			return fmt.Errorf("expected 1 subscription, got %d", len(resp.EventSubscriptionsList))
		}
		for _, id := range resp.EventSubscriptionsList[0].SourceIdsList {
			if id == "cluster-extra" {
				return fmt.Errorf("expected 'cluster-extra' to be removed from SourceIdsList")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyEventSubscription", func() error {
		_, err := tc.client.ModifyEventSubscription(tc.ctx, &neptune.ModifyEventSubscriptionInput{
			SubscriptionName: aws.String(subName),
			SourceType:       aws.String("db-instance"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ModifyEventSubscription_Verify", func() error {
		resp, err := tc.client.DescribeEventSubscriptions(tc.ctx, &neptune.DescribeEventSubscriptionsInput{
			SubscriptionName: aws.String(subName),
		})
		if err != nil {
			return err
		}
		if len(resp.EventSubscriptionsList) != 1 {
			return fmt.Errorf("expected 1 subscription, got %d", len(resp.EventSubscriptionsList))
		}
		sub := resp.EventSubscriptionsList[0]
		if sub.SourceType == nil || *sub.SourceType != "db-instance" {
			return fmt.Errorf("expected SourceType=db-instance after modify, got %v", sub.SourceType)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DeleteEventSubscription", func() error {
		_, err := tc.client.DeleteEventSubscription(tc.ctx, &neptune.DeleteEventSubscriptionInput{
			SubscriptionName: aws.String(subName),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteEventSubscription_VerifyDeleted", func() error {
		resp, err := tc.client.DescribeEventSubscriptions(tc.ctx, &neptune.DescribeEventSubscriptionsInput{})
		if err != nil {
			return err
		}
		for _, sub := range resp.EventSubscriptionsList {
			if sub.CustSubscriptionId != nil && *sub.CustSubscriptionId == subName {
				return fmt.Errorf("expected subscription to be deleted but still found")
			}
		}
		return nil
	}))

	return results
}
