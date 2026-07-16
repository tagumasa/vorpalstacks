package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func (r *TestRunner) runRoute53TagTests(tc *route53TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("route53", "ChangeTagsForResource", func() error {
		tagDomain := tc.domain("tags")
		tagRef := fmt.Sprintf("tagref-%d", tc.uniq)
		createResp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(tagDomain),
			CallerReference: aws.String(tagRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		tagZoneID := aws.ToString(createResp.HostedZone.Id)

		defer tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(tagZoneID)})

		_, err = tc.client.ChangeTagsForResource(tc.ctx, &route53.ChangeTagsForResourceInput{
			ResourceType: types.TagResourceTypeHostedzone,
			ResourceId:   aws.String(tagZoneID),
			AddTags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Owner"), Value: aws.String("team-a")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags: %v", err)
		}

		listResp, err := tc.client.ListTagsForResource(tc.ctx, &route53.ListTagsForResourceInput{
			ResourceType: types.TagResourceTypeHostedzone,
			ResourceId:   aws.String(tagZoneID),
		})
		if err != nil {
			return fmt.Errorf("list tags after add: %v", err)
		}
		if listResp.ResourceTagSet == nil {
			return fmt.Errorf("resource tag set is nil")
		}
		if len(listResp.ResourceTagSet.Tags) != 2 {
			return fmt.Errorf("expected 2 tags after add, got %d", len(listResp.ResourceTagSet.Tags))
		}
		tagMap := make(map[string]string)
		for _, t := range listResp.ResourceTagSet.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["Environment"] != "test" {
			return fmt.Errorf("Environment tag mismatch: got %q", tagMap["Environment"])
		}
		if tagMap["Owner"] != "team-a" {
			return fmt.Errorf("Owner tag mismatch: got %q", tagMap["Owner"])
		}

		_, err = tc.client.ChangeTagsForResource(tc.ctx, &route53.ChangeTagsForResourceInput{
			ResourceType:  types.TagResourceTypeHostedzone,
			ResourceId:    aws.String(tagZoneID),
			RemoveTagKeys: []string{"Owner"},
		})
		if err != nil {
			return fmt.Errorf("remove tags: %v", err)
		}

		listResp2, err := tc.client.ListTagsForResource(tc.ctx, &route53.ListTagsForResourceInput{
			ResourceType: types.TagResourceTypeHostedzone,
			ResourceId:   aws.String(tagZoneID),
		})
		if err != nil {
			return fmt.Errorf("list tags after remove: %v", err)
		}
		if len(listResp2.ResourceTagSet.Tags) != 1 {
			return fmt.Errorf("expected 1 tag after removal, got %d", len(listResp2.ResourceTagSet.Tags))
		}
		if aws.ToString(listResp2.ResourceTagSet.Tags[0].Key) != "Environment" {
			return fmt.Errorf("expected Environment tag, got %q", aws.ToString(listResp2.ResourceTagSet.Tags[0].Key))
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListTagsForResource_HealthCheck", func() error {
		hcRef := fmt.Sprintf("hctagref-%d", tc.uniq)
		hcResp, err := tc.client.CreateHealthCheck(tc.ctx, &route53.CreateHealthCheckInput{
			CallerReference: aws.String(hcRef),
			HealthCheckConfig: &types.HealthCheckConfig{
				Type:                     types.HealthCheckTypeTcp,
				FullyQualifiedDomainName: aws.String("hctag.example.com"),
			},
		})
		if err != nil {
			return fmt.Errorf("create hc: %v", err)
		}
		hcID := aws.ToString(hcResp.HealthCheck.Id)

		defer tc.client.DeleteHealthCheck(tc.ctx, &route53.DeleteHealthCheckInput{HealthCheckId: aws.String(hcID)})

		_, err = tc.client.ChangeTagsForResource(tc.ctx, &route53.ChangeTagsForResourceInput{
			ResourceType: types.TagResourceTypeHealthcheck,
			ResourceId:   aws.String(hcID),
			AddTags: []types.Tag{
				{Key: aws.String("Monitor"), Value: aws.String("enabled")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags: %v", err)
		}

		listResp, err := tc.client.ListTagsForResource(tc.ctx, &route53.ListTagsForResourceInput{
			ResourceType: types.TagResourceTypeHealthcheck,
			ResourceId:   aws.String(hcID),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if listResp.ResourceTagSet == nil {
			return fmt.Errorf("resource tag set is nil")
		}
		if len(listResp.ResourceTagSet.Tags) != 1 {
			return fmt.Errorf("expected 1 tag, got %d", len(listResp.ResourceTagSet.Tags))
		}
		if aws.ToString(listResp.ResourceTagSet.Tags[0].Key) != "Monitor" {
			return fmt.Errorf("expected Monitor tag, got %q", aws.ToString(listResp.ResourceTagSet.Tags[0].Key))
		}
		if aws.ToString(listResp.ResourceTagSet.Tags[0].Value) != "enabled" {
			return fmt.Errorf("expected enabled value, got %q", aws.ToString(listResp.ResourceTagSet.Tags[0].Value))
		}
		return nil
	}))

	return results
}
