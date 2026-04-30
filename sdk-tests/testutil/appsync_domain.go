package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

func (r *TestRunner) runAppSyncDomainTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client
	uid := res.uid

	res.domainName = fmt.Sprintf("test-domain-%d.example.com", uid)

	results = append(results, r.RunTest("appsync", "CreateDomainName", func() error {
		resp, err := client.CreateDomainName(ctx, &appsync.CreateDomainNameInput{
			DomainName:     aws.String(res.domainName),
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/test-cert"),
			Description:    aws.String("test domain"),
		})
		if err != nil {
			return err
		}
		if resp.DomainNameConfig == nil {
			return fmt.Errorf("DomainNameConfig is nil")
		}
		if *resp.DomainNameConfig.DomainName != res.domainName {
			return fmt.Errorf("expected domain %s, got %s", res.domainName, *resp.DomainNameConfig.DomainName)
		}
		if resp.DomainNameConfig.DomainNameArn == nil || *resp.DomainNameConfig.DomainNameArn == "" {
			return fmt.Errorf("DomainNameArn is nil")
		}
		if resp.DomainNameConfig.AppsyncDomainName == nil || *resp.DomainNameConfig.AppsyncDomainName == "" {
			return fmt.Errorf("AppsyncDomainName is nil")
		}
		if resp.DomainNameConfig.HostedZoneId == nil || *resp.DomainNameConfig.HostedZoneId == "" {
			return fmt.Errorf("HostedZoneId is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateDomainName_WithTags", func() error {
		res.tagDomainName = fmt.Sprintf("tag-domain-%d.example.com", uid)
		resp, err := client.CreateDomainName(ctx, &appsync.CreateDomainNameInput{
			DomainName:     aws.String(res.tagDomainName),
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/tag-cert"),
			Tags:           map[string]string{"env": "prod"},
		})
		if err != nil {
			return err
		}
		if resp.DomainNameConfig == nil || resp.DomainNameConfig.Tags == nil {
			return fmt.Errorf("tags not returned")
		}
		if resp.DomainNameConfig.Tags["env"] != "prod" {
			return fmt.Errorf("tag not persisted: %v", resp.DomainNameConfig.Tags)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetDomainName", func() error {
		resp, err := client.GetDomainName(ctx, &appsync.GetDomainNameInput{
			DomainName: aws.String(res.domainName),
		})
		if err != nil {
			return err
		}
		if resp.DomainNameConfig == nil || *resp.DomainNameConfig.DomainName != res.domainName {
			return fmt.Errorf("expected domain %s, got %v", res.domainName, resp.DomainNameConfig)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetDomainName_NonExistent", func() error {
		_, err := client.GetDomainName(ctx, &appsync.GetDomainNameInput{
			DomainName: aws.String("does-not-exist.example.com"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "ListDomainNames", func() error {
		resp, err := client.ListDomainNames(ctx, &appsync.ListDomainNamesInput{})
		if err != nil {
			return err
		}
		if len(resp.DomainNameConfigs) < 2 {
			return fmt.Errorf("expected at least 2 domain names, got %d", len(resp.DomainNameConfigs))
		}
		found := false
		for _, dnc := range resp.DomainNameConfigs {
			if dnc.DomainName != nil && *dnc.DomainName == res.domainName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created domain %s not found in list", res.domainName)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListDomainNames_WithPagination", func() error {
		resp, err := client.ListDomainNames(ctx, &appsync.ListDomainNamesInput{MaxResults: 1})
		if err != nil {
			return err
		}
		if len(resp.DomainNameConfigs) != 1 {
			return fmt.Errorf("expected 1 domain name with maxResults=1, got %d", len(resp.DomainNameConfigs))
		}
		if resp.NextToken == nil {
			return fmt.Errorf("expected NextToken when more results exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateDomainName", func() error {
		resp, err := client.UpdateDomainName(ctx, &appsync.UpdateDomainNameInput{
			DomainName:  aws.String(res.domainName),
			Description: aws.String("updated description"),
		})
		if err != nil {
			return err
		}
		if resp.DomainNameConfig == nil {
			return fmt.Errorf("DomainNameConfig is nil")
		}
		if resp.DomainNameConfig.Description == nil || *resp.DomainNameConfig.Description != "updated description" {
			return fmt.Errorf("description not updated: %v", resp.DomainNameConfig.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateDomainName_NonExistent", func() error {
		_, err := client.UpdateDomainName(ctx, &appsync.UpdateDomainNameInput{
			DomainName: aws.String("does-not-exist.example.com"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.runAppSyncDomainAssociationTests(res)...)
	results = append(results, r.runAppSyncMergedApiTests(res)...)

	return results
}

func (r *TestRunner) runAppSyncDomainAssociationTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "AssociateApi", func() error {
		resp, err := client.AssociateApi(ctx, &appsync.AssociateApiInput{
			DomainName: aws.String(res.domainName),
			ApiId:      aws.String(res.gqlApiId),
		})
		if err != nil {
			return err
		}
		if resp.ApiAssociation == nil {
			return fmt.Errorf("ApiAssociation is nil")
		}
		if resp.ApiAssociation.ApiId == nil || *resp.ApiAssociation.ApiId != res.gqlApiId {
			return fmt.Errorf("expected apiId %s, got %v", res.gqlApiId, resp.ApiAssociation.ApiId)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetApiAssociation", func() error {
		resp, err := client.GetApiAssociation(ctx, &appsync.GetApiAssociationInput{
			DomainName: aws.String(res.domainName),
		})
		if err != nil {
			return err
		}
		if resp.ApiAssociation == nil {
			return fmt.Errorf("ApiAssociation is nil")
		}
		if resp.ApiAssociation.ApiId == nil || *resp.ApiAssociation.ApiId != res.gqlApiId {
			return fmt.Errorf("expected apiId %s, got %v", res.gqlApiId, resp.ApiAssociation.ApiId)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateApi", func() error {
		_, err := client.DisassociateApi(ctx, &appsync.DisassociateApiInput{
			DomainName: aws.String(res.domainName),
		})
		if err != nil {
			return err
		}
		_, err = client.GetApiAssociation(ctx, &appsync.GetApiAssociationInput{
			DomainName: aws.String(res.domainName),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DisassociateApi_NotAssociated", func() error {
		_, err := client.DisassociateApi(ctx, &appsync.DisassociateApiInput{
			DomainName: aws.String(res.domainName),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	return results
}

func (r *TestRunner) runAppSyncMergedApiTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client
	uid := res.uid

	results = append(results, r.RunTest("appsync", "CreateGraphqlApi_ForMerged", func() error {
		resp, err := client.CreateGraphqlApi(ctx, &appsync.CreateGraphqlApiInput{
			Name:               aws.String(fmt.Sprintf("merged-api-%d", uid)),
			AuthenticationType: types.AuthenticationTypeApiKey,
			ApiType:            types.GraphQLApiTypeMerged,
		})
		if err != nil {
			return err
		}
		if resp.GraphqlApi == nil || resp.GraphqlApi.ApiId == nil {
			return fmt.Errorf("invalid response")
		}
		if resp.GraphqlApi.ApiType != types.GraphQLApiTypeMerged {
			return fmt.Errorf("expected MERGED api type, got %s", resp.GraphqlApi.ApiType)
		}
		res.mergedApiId = *resp.GraphqlApi.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateGraphqlApi_ForSource", func() error {
		resp, err := client.CreateGraphqlApi(ctx, &appsync.CreateGraphqlApiInput{
			Name:               aws.String(fmt.Sprintf("source-api-%d", uid)),
			AuthenticationType: types.AuthenticationTypeApiKey,
		})
		if err != nil {
			return err
		}
		if resp.GraphqlApi == nil || resp.GraphqlApi.ApiId == nil {
			return fmt.Errorf("invalid response")
		}
		if resp.GraphqlApi.AuthenticationType != types.AuthenticationTypeApiKey {
			return fmt.Errorf("expected API_KEY auth type, got %s", resp.GraphqlApi.AuthenticationType)
		}
		res.sourceApiId2 = *resp.GraphqlApi.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "AssociateSourceGraphqlApi", func() error {
		resp, err := client.AssociateSourceGraphqlApi(ctx, &appsync.AssociateSourceGraphqlApiInput{
			MergedApiIdentifier: aws.String(res.mergedApiId),
			SourceApiIdentifier: aws.String(res.sourceApiId2),
			Description:         aws.String("test association"),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociation == nil {
			return fmt.Errorf("SourceApiAssociation is nil")
		}
		if resp.SourceApiAssociation.AssociationId == nil || *resp.SourceApiAssociation.AssociationId == "" {
			return fmt.Errorf("AssociationId is empty")
		}
		if resp.SourceApiAssociation.SourceApiAssociationStatus != types.SourceApiAssociationStatusMergeScheduled {
			return fmt.Errorf("expected MERGE_SCHEDULED, got %s", resp.SourceApiAssociation.SourceApiAssociationStatus)
		}
		if resp.SourceApiAssociation.Description == nil || *resp.SourceApiAssociation.Description != "test association" {
			return fmt.Errorf("description not set: %v", resp.SourceApiAssociation.Description)
		}
		res.associationId = *resp.SourceApiAssociation.AssociationId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetSourceApiAssociation", func() error {
		resp, err := client.GetSourceApiAssociation(ctx, &appsync.GetSourceApiAssociationInput{
			MergedApiIdentifier: aws.String(res.mergedApiId),
			AssociationId:       aws.String(res.associationId),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociation == nil || *resp.SourceApiAssociation.AssociationId != res.associationId {
			return fmt.Errorf("expected association %s, got %v", res.associationId, resp.SourceApiAssociation)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetSourceApiAssociation_NonExistent", func() error {
		_, err := client.GetSourceApiAssociation(ctx, &appsync.GetSourceApiAssociationInput{
			MergedApiIdentifier: aws.String(res.mergedApiId),
			AssociationId:       aws.String("does-not-exist"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "UpdateSourceApiAssociation", func() error {
		resp, err := client.UpdateSourceApiAssociation(ctx, &appsync.UpdateSourceApiAssociationInput{
			MergedApiIdentifier: aws.String(res.mergedApiId),
			AssociationId:       aws.String(res.associationId),
			Description:         aws.String("updated association"),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociation == nil {
			return fmt.Errorf("SourceApiAssociation is nil")
		}
		if resp.SourceApiAssociation.Description == nil || *resp.SourceApiAssociation.Description != "updated association" {
			return fmt.Errorf("description not updated: %v", resp.SourceApiAssociation.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "StartSchemaMerge", func() error {
		resp, err := client.StartSchemaMerge(ctx, &appsync.StartSchemaMergeInput{
			MergedApiIdentifier: aws.String(res.mergedApiId),
			AssociationId:       aws.String(res.associationId),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociationStatus != types.SourceApiAssociationStatusMergeSuccess {
			return fmt.Errorf("expected MERGE_SUCCESS, got %s", resp.SourceApiAssociationStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListSourceApiAssociations", func() error {
		resp, err := client.ListSourceApiAssociations(ctx, &appsync.ListSourceApiAssociationsInput{
			ApiId: aws.String(res.mergedApiId),
		})
		if err != nil {
			return err
		}
		if len(resp.SourceApiAssociationSummaries) < 1 {
			return fmt.Errorf("expected at least 1 association, got %d", len(resp.SourceApiAssociationSummaries))
		}
		found := false
		for _, a := range resp.SourceApiAssociationSummaries {
			if a.AssociationId != nil && *a.AssociationId == res.associationId {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created association not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateSourceGraphqlApi", func() error {
		resp, err := client.DisassociateSourceGraphqlApi(ctx, &appsync.DisassociateSourceGraphqlApiInput{
			MergedApiIdentifier: aws.String(res.mergedApiId),
			AssociationId:       aws.String(res.associationId),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociationStatus != types.SourceApiAssociationStatusDeletionScheduled {
			return fmt.Errorf("expected DELETION_SCHEDULED, got %s", resp.SourceApiAssociationStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateSourceGraphqlApi_VerifyDeletionScheduled", func() error {
		resp, err := client.GetSourceApiAssociation(ctx, &appsync.GetSourceApiAssociationInput{
			MergedApiIdentifier: aws.String(res.mergedApiId),
			AssociationId:       aws.String(res.associationId),
		})
		if err != nil {
			return fmt.Errorf("association should still exist in DELETION_SCHEDULED state: %v", err)
		}
		if resp.SourceApiAssociation == nil {
			return fmt.Errorf("SourceApiAssociation is nil")
		}
		if resp.SourceApiAssociation.SourceApiAssociationStatus != types.SourceApiAssociationStatusDeletionScheduled {
			return fmt.Errorf("expected DELETION_SCHEDULED, got %s", resp.SourceApiAssociation.SourceApiAssociationStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateSourceGraphqlApi_NonExistent", func() error {
		_, err := client.DisassociateSourceGraphqlApi(ctx, &appsync.DisassociateSourceGraphqlApiInput{
			MergedApiIdentifier: aws.String(res.mergedApiId),
			AssociationId:       aws.String("already-deleted"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.runAppSyncMergedGraphqlApiAssociationTests(res)...)

	return results
}

func (r *TestRunner) runAppSyncMergedGraphqlApiAssociationTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "AssociateMergedGraphqlApi", func() error {
		resp, err := client.AssociateMergedGraphqlApi(ctx, &appsync.AssociateMergedGraphqlApiInput{
			SourceApiIdentifier: aws.String(res.sourceApiId2),
			MergedApiIdentifier: aws.String(res.mergedApiId),
			Description:         aws.String("merged from source side"),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociation == nil {
			return fmt.Errorf("SourceApiAssociation is nil")
		}
		if resp.SourceApiAssociation.AssociationId == nil || *resp.SourceApiAssociation.AssociationId == "" {
			return fmt.Errorf("AssociationId is empty")
		}
		if resp.SourceApiAssociation.Description == nil || *resp.SourceApiAssociation.Description != "merged from source side" {
			return fmt.Errorf("description not set: %v", resp.SourceApiAssociation.Description)
		}
		res.mergedAssocId = *resp.SourceApiAssociation.AssociationId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateMergedGraphqlApi", func() error {
		resp, err := client.DisassociateMergedGraphqlApi(ctx, &appsync.DisassociateMergedGraphqlApiInput{
			SourceApiIdentifier: aws.String(res.sourceApiId2),
			AssociationId:       aws.String(res.mergedAssocId),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociationStatus != types.SourceApiAssociationStatusDeletionScheduled {
			return fmt.Errorf("expected DELETION_SCHEDULED, got %s", resp.SourceApiAssociationStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateMergedGraphqlApi_VerifyDeletionScheduled", func() error {
		resp, err := client.GetSourceApiAssociation(ctx, &appsync.GetSourceApiAssociationInput{
			MergedApiIdentifier: aws.String(res.mergedApiId),
			AssociationId:       aws.String(res.mergedAssocId),
		})
		if err != nil {
			return fmt.Errorf("association should still exist in DELETION_SCHEDULED state: %v", err)
		}
		if resp.SourceApiAssociation == nil {
			return fmt.Errorf("SourceApiAssociation is nil")
		}
		if resp.SourceApiAssociation.SourceApiAssociationStatus != types.SourceApiAssociationStatusDeletionScheduled {
			return fmt.Errorf("expected DELETION_SCHEDULED, got %s", resp.SourceApiAssociation.SourceApiAssociationStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateMergedGraphqlApi_NonExistent", func() error {
		_, err := client.DisassociateMergedGraphqlApi(ctx, &appsync.DisassociateMergedGraphqlApiInput{
			SourceApiIdentifier: aws.String(res.sourceApiId2),
			AssociationId:       aws.String("already-deleted"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	return results
}
