package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func (r *TestRunner) runACMTagTests(tc *acmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("acm", "RequestCertificate_WithTags", func() error {
		domain := acmUniqueDomain("tag-test")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
			Tags: []types.Tag{
				{Key: aws.String("Env"), Value: aws.String("prod")},
				{Key: aws.String("Team"), Value: aws.String("platform")},
			},
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		tagResp, err := tc.client.ListTagsForCertificate(tc.ctx, &acm.ListTagsForCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(tagResp.Tags) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(tagResp.Tags))
		}
		tagMap := make(map[string]string)
		for _, t := range tagResp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["Env"] != "prod" {
			return fmt.Errorf("expected Env=prod, got %s", tagMap["Env"])
		}
		if tagMap["Team"] != "platform" {
			return fmt.Errorf("expected Team=platform, got %s", tagMap["Team"])
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "AddTagsToCertificate_UpdateExistingTag", func() error {
		domain := acmUniqueDomain("tagupd")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		tc.client.AddTagsToCertificate(tc.ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: resp.CertificateArn,
			Tags:           []types.Tag{{Key: aws.String("Env"), Value: aws.String("dev")}},
		})
		tc.client.AddTagsToCertificate(tc.ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: resp.CertificateArn,
			Tags:           []types.Tag{{Key: aws.String("Env"), Value: aws.String("prod")}},
		})
		tagResp, err := tc.client.ListTagsForCertificate(tc.ctx, &acm.ListTagsForCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		for _, t := range tagResp.Tags {
			if aws.ToString(t.Key) == "Env" {
				if aws.ToString(t.Value) != "prod" {
					return fmt.Errorf("expected Env=prod, got %s", aws.ToString(t.Value))
				}
				return nil
			}
		}
		return fmt.Errorf("env tag not found")
	}))

	results = append(results, r.RunTest("acm", "AddTagsToCertificate_VerifyContent", func() error {
		domain := acmUniqueDomain("tagver")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		tc.client.AddTagsToCertificate(tc.ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: resp.CertificateArn,
			Tags: []types.Tag{
				{Key: aws.String("Key1"), Value: aws.String("Val1")},
				{Key: aws.String("Key2"), Value: aws.String("Val2")},
			},
		})
		tagResp, err := tc.client.ListTagsForCertificate(tc.ctx, &acm.ListTagsForCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(tagResp.Tags) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(tagResp.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RemoveTagsFromCertificate_VerifyEmpty", func() error {
		domain := acmUniqueDomain("tagrm")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		tc.client.AddTagsToCertificate(tc.ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: resp.CertificateArn,
			Tags:           []types.Tag{{Key: aws.String("X"), Value: aws.String("Y")}},
		})
		tc.client.RemoveTagsFromCertificate(tc.ctx, &acm.RemoveTagsFromCertificateInput{
			CertificateArn: resp.CertificateArn,
			Tags:           []types.Tag{{Key: aws.String("X")}},
		})
		tagResp, err := tc.client.ListTagsForCertificate(tc.ctx, &acm.ListTagsForCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(tagResp.Tags) != 0 {
			return fmt.Errorf("expected 0 tags after removal, got %d", len(tagResp.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ListTagsForCertificate_VerifyMultipleTags", func() error {
		domain := acmUniqueDomain("listtag")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
			Tags: []types.Tag{
				{Key: aws.String("A"), Value: aws.String("1")},
				{Key: aws.String("B"), Value: aws.String("2")},
				{Key: aws.String("C"), Value: aws.String("3")},
			},
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		tagResp, err := tc.client.ListTagsForCertificate(tc.ctx, &acm.ListTagsForCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(tagResp.Tags) != 3 {
			return fmt.Errorf("expected 3 tags, got %d", len(tagResp.Tags))
		}
		tagMap := make(map[string]string)
		for _, t := range tagResp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["A"] != "1" || tagMap["B"] != "2" || tagMap["C"] != "3" {
			return fmt.Errorf("tag values mismatch: %+v", tagMap)
		}
		return nil
	}))

	// Generic tag API: TagResource / UntagResource / ListTagsForResource
	results = append(results, r.RunTest("acm", "TagResource_AddAndList", func() error {
		domain := acmUniqueDomain("generic-tag")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		_, err = tc.client.TagResource(tc.ctx, &acm.TagResourceInput{
			ResourceArn: resp.CertificateArn,
			Tags: []types.Tag{
				{Key: aws.String("Project"), Value: aws.String("acme")},
				{Key: aws.String("Owner"), Value: aws.String("devops")},
			},
		})
		if err != nil {
			return err
		}

		listResp, err := tc.client.ListTagsForResource(tc.ctx, &acm.ListTagsForResourceInput{
			ResourceArn: resp.CertificateArn,
		})
		if err != nil {
			return err
		}
		if len(listResp.Tags) != 2 {
			return fmt.Errorf("expected 2 tags via ListTagsForResource, got %d", len(listResp.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "UntagResource_RemoveByKey", func() error {
		domain := acmUniqueDomain("generic-untag")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
			Tags: []types.Tag{
				{Key: aws.String("Keep"), Value: aws.String("yes")},
				{Key: aws.String("Remove"), Value: aws.String("no")},
			},
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		_, err = tc.client.UntagResource(tc.ctx, &acm.UntagResourceInput{
			ResourceArn: resp.CertificateArn,
			TagKeys:     []string{"Remove"},
		})
		if err != nil {
			return err
		}

		listResp, err := tc.client.ListTagsForResource(tc.ctx, &acm.ListTagsForResourceInput{
			ResourceArn: resp.CertificateArn,
		})
		if err != nil {
			return err
		}
		if len(listResp.Tags) != 1 {
			return fmt.Errorf("expected 1 tag after untag, got %d", len(listResp.Tags))
		}
		if aws.ToString(listResp.Tags[0].Key) != "Keep" {
			return fmt.Errorf("expected Keep tag, got %s", aws.ToString(listResp.Tags[0].Key))
		}
		return nil
	}))

	// SearchCertificates
	results = append(results, r.RunTest("acm", "SearchCertificates_ReturnsResults", func() error {
		domain := acmUniqueDomain("search-test")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		searchResp, err := tc.client.SearchCertificates(tc.ctx, &acm.SearchCertificatesInput{})
		if err != nil {
			return err
		}
		if len(searchResp.Results) == 0 {
			return fmt.Errorf("expected at least 1 search result")
		}
		// Verify each result has required fields.
		for _, r := range searchResp.Results {
			if aws.ToString(r.CertificateArn) == "" {
				return fmt.Errorf("CertificateArn empty in search result")
			}
		}
		return nil
	}))

	return results
}
