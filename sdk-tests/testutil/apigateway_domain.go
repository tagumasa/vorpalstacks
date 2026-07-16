package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayDomainTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "CreateDomainName", func() error {
		domain := fmt.Sprintf("test-%d.example.com", time.Now().UnixNano())
		resp, err := client.CreateDomainName(ctx, &apigateway.CreateDomainNameInput{
			DomainName:      aws.String(domain),
			CertificateName: aws.String("test-cert"),
			Tags: map[string]string{
				"domain": "test",
			},
		})
		if err != nil {
			return err
		}
		if resp.DomainName == nil || *resp.DomainName != domain {
			return fmt.Errorf("domain name mismatch, got %v", resp.DomainName)
		}
		if resp.DomainNameId == nil {
			return fmt.Errorf("domain name ID is nil")
		}
		return nil
	}))

	var domainName string
	results = append(results, r.RunTest("apigateway", "GetDomainNames", func() error {
		resp, err := client.GetDomainNames(ctx, &apigateway.GetDomainNamesInput{
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 domain name")
		}
		domainName = *resp.Items[0].DomainName
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetDomainName", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		resp, err := client.GetDomainName(ctx, &apigateway.GetDomainNameInput{
			DomainName: aws.String(domainName),
		})
		if err != nil {
			return err
		}
		if resp.DomainName == nil || *resp.DomainName != domainName {
			return fmt.Errorf("domain name mismatch, got %v", resp.DomainName)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateDomainName", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		resp, err := client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/certificateName"),
					Value: aws.String("updated-cert"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.CertificateName == nil || *resp.CertificateName != "updated-cert" {
			return fmt.Errorf("certificateName not updated, got %v", resp.CertificateName)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateBasePathMapping", func() error {
		if apiID == "" || domainName == "" {
			return fmt.Errorf("API ID or domain name not available")
		}
		resp, err := client.CreateBasePathMapping(ctx, &apigateway.CreateBasePathMappingInput{
			DomainName: aws.String(domainName),
			RestApiId:  aws.String(apiID),
			BasePath:   aws.String("v1"),
			Stage:      aws.String("prod"),
		})
		if err != nil {
			return err
		}
		if resp.BasePath == nil || *resp.BasePath != "v1" {
			return fmt.Errorf("basePath mismatch, got %v", resp.BasePath)
		}
		if resp.RestApiId == nil || *resp.RestApiId != apiID {
			return fmt.Errorf("restApiId mismatch, got %v", resp.RestApiId)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetBasePathMappings", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		resp, err := client.GetBasePathMappings(ctx, &apigateway.GetBasePathMappingsInput{
			DomainName: aws.String(domainName),
			Limit:      aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 base path mapping")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetBasePathMapping", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		resp, err := client.GetBasePathMapping(ctx, &apigateway.GetBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
		})
		if err != nil {
			return err
		}
		if resp.BasePath == nil || *resp.BasePath != "v1" {
			return fmt.Errorf("basePath mismatch, got %v", resp.BasePath)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateBasePathMapping", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		resp, err := client.UpdateBasePathMapping(ctx, &apigateway.UpdateBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/stage"),
					Value: aws.String("staging"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Stage == nil || *resp.Stage != "staging" {
			return fmt.Errorf("stage not updated, got %v", resp.Stage)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteBasePathMapping", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		_, err := client.DeleteBasePathMapping(ctx, &apigateway.DeleteBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetBasePathMapping(ctx, &apigateway.GetBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
		})
		if err == nil {
			return fmt.Errorf("GetBasePathMapping should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteDomainName", func() error {
		if domainName == "" {
			return fmt.Errorf("domain name not available")
		}
		_, err := client.DeleteDomainName(ctx, &apigateway.DeleteDomainNameInput{
			DomainName: aws.String(domainName),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetDomainName(ctx, &apigateway.GetDomainNameInput{
			DomainName: aws.String(domainName),
		})
		if err == nil {
			return fmt.Errorf("GetDomainName should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DomainName_BasePathMapping_FullLifecycle", func() error {
		dbAPI := fmt.Sprintf("DbAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(dbAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		domain := fmt.Sprintf("lc-%d.example.com", time.Now().UnixNano())
		dnResp, err := client.CreateDomainName(ctx, &apigateway.CreateDomainNameInput{
			DomainName:      aws.String(domain),
			CertificateName: aws.String("lc-cert"),
		})
		if err != nil {
			return fmt.Errorf("create domain: %v", err)
		}

		_, err = client.CreateBasePathMapping(ctx, &apigateway.CreateBasePathMappingInput{
			DomainName: aws.String(domain),
			RestApiId:  createResp.Id,
			BasePath:   aws.String("(none)"),
			Stage:      aws.String("prod"),
		})
		if err != nil {
			return fmt.Errorf("create base path mapping: %v", err)
		}

		_, err = client.GetBasePathMappings(ctx, &apigateway.GetBasePathMappingsInput{
			DomainName: aws.String(domain),
		})
		if err != nil {
			return fmt.Errorf("get base path mappings: %v", err)
		}

		_, err = client.DeleteBasePathMapping(ctx, &apigateway.DeleteBasePathMappingInput{
			DomainName: aws.String(domain),
			BasePath:   aws.String("(none)"),
		})
		if err != nil {
			return fmt.Errorf("delete base path mapping: %v", err)
		}

		_ = dnResp
		_, err = client.DeleteDomainName(ctx, &apigateway.DeleteDomainNameInput{
			DomainName: aws.String(domain),
		})
		return err
	}))

	results = append(results, r.RunTest("apigateway", "CreateBasePathMapping_DefaultNone", func() error {
		bnAPI := fmt.Sprintf("BnAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(bnAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		domain := fmt.Sprintf("none-%d.example.com", time.Now().UnixNano())
		_, err = client.CreateDomainName(ctx, &apigateway.CreateDomainNameInput{
			DomainName:      aws.String(domain),
			CertificateName: aws.String("none-cert"),
		})
		if err != nil {
			return fmt.Errorf("create domain: %v", err)
		}
		defer client.DeleteDomainName(ctx, &apigateway.DeleteDomainNameInput{DomainName: aws.String(domain)})

		mappingResp, err := client.CreateBasePathMapping(ctx, &apigateway.CreateBasePathMappingInput{
			DomainName: aws.String(domain),
			RestApiId:  createResp.Id,
			Stage:      aws.String("prod"),
		})
		if err != nil {
			return fmt.Errorf("create base path mapping: %v", err)
		}
		if mappingResp.BasePath == nil || *mappingResp.BasePath != "(none)" {
			return fmt.Errorf("expected basePath (none) when omitted, got %v", mappingResp.BasePath)
		}

		getResp, err := client.GetBasePathMapping(ctx, &apigateway.GetBasePathMappingInput{
			DomainName: aws.String(domain),
			BasePath:   aws.String("(none)"),
		})
		if err != nil {
			return fmt.Errorf("get base path mapping with (none): %v", err)
		}
		if getResp.BasePath == nil || *getResp.BasePath != "(none)" {
			return fmt.Errorf("get basePath mismatch, got %v", getResp.BasePath)
		}
		return nil
	}))

	return results
}
