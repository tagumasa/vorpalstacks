package testutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
)

func (r *TestRunner) runAppSyncCleanupTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "DeleteDomainName", func() error {
		_, err := client.DeleteDomainName(ctx, &appsync.DeleteDomainNameInput{
			DomainName: aws.String(res.domainName),
		})
		if err != nil {
			return err
		}
		_, err = client.GetDomainName(ctx, &appsync.GetDomainNameInput{
			DomainName: aws.String(res.domainName),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteDomainName_NonExistent", func() error {
		_, err := client.DeleteDomainName(ctx, &appsync.DeleteDomainNameInput{
			DomainName: aws.String("already-deleted.example.com"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteDomainName_WithTags", func() error {
		_, err := client.DeleteDomainName(ctx, &appsync.DeleteDomainNameInput{
			DomainName: aws.String(res.tagDomainName),
		})
		if err != nil {
			return err
		}
		_, err = client.GetDomainName(ctx, &appsync.GetDomainNameInput{
			DomainName: aws.String(res.tagDomainName),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteGraphqlApi", func() error {
		if res.descApiKeyId != "" {
			client.DeleteApiKey(ctx, &appsync.DeleteApiKeyInput{
				ApiId: aws.String(res.gqlApiId),
				Id:    aws.String(res.descApiKeyId),
			})
		}
		_, err := client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(res.gqlApiId)})
		if err != nil {
			return err
		}
		_, err = client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String(res.gqlApiId)})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteGraphqlApi_NonExistent", func() error {
		_, err := client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String("already-deleted")})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteGraphqlApi_WithTags", func() error {
		_, err := client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(res.gqlTagsApiId)})
		if err != nil {
			return err
		}
		_, err = client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String(res.gqlTagsApiId)})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteGraphqlApi_Merged", func() error {
		_, err := client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(res.mergedApiId)})
		if err != nil {
			return err
		}
		_, err = client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String(res.mergedApiId)})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteGraphqlApi_Source", func() error {
		_, err := client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(res.sourceApiId2)})
		if err != nil {
			return err
		}
		_, err = client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String(res.sourceApiId2)})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteChannelNamespace", func() error {
		_, err := client.DeleteChannelNamespace(ctx, &appsync.DeleteChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String(res.nsName),
		})
		if err != nil {
			return err
		}
		_, err = client.GetChannelNamespace(ctx, &appsync.GetChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String(res.nsName),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteChannelNamespace_NonExistent", func() error {
		_, err := client.DeleteChannelNamespace(ctx, &appsync.DeleteChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String("already-deleted"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteChannelNamespace_ForTagging", func() error {
		_, err := client.DeleteChannelNamespace(ctx, &appsync.DeleteChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String(res.tagNsName),
		})
		if err != nil {
			return err
		}
		_, err = client.GetChannelNamespace(ctx, &appsync.GetChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String(res.tagNsName),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.runAppSyncEventApiPaginationTests(res)...)

	results = append(results, r.RunTest("appsync", "DeleteApi", func() error {
		_, err := client.DeleteApi(ctx, &appsync.DeleteApiInput{
			ApiId: aws.String(res.apiId),
		})
		if err != nil {
			return err
		}
		_, err = client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String(res.apiId)})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteApi_NonExistent", func() error {
		_, err := client.DeleteApi(ctx, &appsync.DeleteApiInput{
			ApiId: aws.String("already-deleted"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteApi_WithTags", func() error {
		_, err := client.DeleteApi(ctx, &appsync.DeleteApiInput{
			ApiId: aws.String(res.tagsApiId),
		})
		if err != nil {
			return err
		}
		_, err = client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String(res.tagsApiId)})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteApi_ForTagging", func() error {
		_, err := client.DeleteApi(ctx, &appsync.DeleteApiInput{
			ApiId: aws.String(res.taggedApiId),
		})
		if err != nil {
			return err
		}
		_, err = client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String(res.taggedApiId)})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteApi_WithOwnerContact", func() error {
		_, err := client.DeleteApi(ctx, &appsync.DeleteApiInput{
			ApiId: aws.String(res.ownerApiId),
		})
		if err != nil {
			return err
		}
		_, err = client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String(res.ownerApiId)})
		return AssertErrorContains(err, "NotFoundException")
	}))

	return results
}
