package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
)

func (r *TestRunner) runAppSyncChannelTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client
	uid := res.uid

	res.nsName = fmt.Sprintf("testns-%d", uid)

	results = append(results, r.RunTest("appsync", "CreateChannelNamespace", func() error {
		resp, err := client.CreateChannelNamespace(ctx, &appsync.CreateChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String(res.nsName),
			Tags: map[string]string{
				"type": "test",
			},
		})
		if err != nil {
			return err
		}
		if resp.ChannelNamespace == nil {
			return fmt.Errorf("ChannelNamespace is nil")
		}
		if *resp.ChannelNamespace.Name != res.nsName {
			return fmt.Errorf("expected name %s, got %s", res.nsName, *resp.ChannelNamespace.Name)
		}
		if resp.ChannelNamespace.ApiId == nil || *resp.ChannelNamespace.ApiId != res.apiId {
			return fmt.Errorf("ApiId mismatch: %v", resp.ChannelNamespace.ApiId)
		}
		if resp.ChannelNamespace.ChannelNamespaceArn == nil {
			return fmt.Errorf("ChannelNamespaceArn is nil")
		}
		if resp.ChannelNamespace.Created == nil {
			return fmt.Errorf("created timestamp is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetChannelNamespace", func() error {
		resp, err := client.GetChannelNamespace(ctx, &appsync.GetChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String(res.nsName),
		})
		if err != nil {
			return err
		}
		if resp.ChannelNamespace == nil || *resp.ChannelNamespace.Name != res.nsName {
			return fmt.Errorf("expected namespace %s, got %v", res.nsName, resp.ChannelNamespace)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetChannelNamespace_NonExistent", func() error {
		_, err := client.GetChannelNamespace(ctx, &appsync.GetChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String("does-not-exist"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "ListChannelNamespaces", func() error {
		resp, err := client.ListChannelNamespaces(ctx, &appsync.ListChannelNamespacesInput{
			ApiId: aws.String(res.apiId),
		})
		if err != nil {
			return err
		}
		if len(resp.ChannelNamespaces) < 1 {
			return fmt.Errorf("expected at least 1 namespace, got %d", len(resp.ChannelNamespaces))
		}
		found := false
		for _, ns := range resp.ChannelNamespaces {
			if ns.Name != nil && *ns.Name == res.nsName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created namespace %s not found in list", res.nsName)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateChannelNamespace", func() error {
		resp, err := client.UpdateChannelNamespace(ctx, &appsync.UpdateChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String(res.nsName),
		})
		if err != nil {
			return err
		}
		if resp.ChannelNamespace == nil || *resp.ChannelNamespace.Name != res.nsName {
			return fmt.Errorf("expected namespace %s, got %v", res.nsName, resp.ChannelNamespace)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateChannelNamespace_NonExistent", func() error {
		_, err := client.UpdateChannelNamespace(ctx, &appsync.UpdateChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String("does-not-exist"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	return results
}
