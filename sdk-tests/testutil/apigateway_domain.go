package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"vorpalstacks-sdk-tests/config"
)

// createTestCertArn imports a self-signed certificate into ACM and returns
// its ARN.  API Gateway custom domains require a valid ACM certificate ARN;
// CertificateName alone is not sufficient.
func (r *TestRunner) createTestCertArn(ctx context.Context) (string, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return "", fmt.Errorf("failed to load ACM config: %w", err)
	}
	acmClient := acm.NewFromConfig(cfg)
	resp, err := acmClient.ImportCertificate(ctx, &acm.ImportCertificateInput{
		Certificate: testCertPEM,
		PrivateKey:  testKeyPEM,
	})
	if err != nil {
		return "", fmt.Errorf("failed to import test certificate: %w", err)
	}
	return *resp.CertificateArn, nil
}

// cleanupStaleDomainNames deletes accumulated test-residue domain names from
// previous test runs. When the test runner is interrupted (SIGKILL, timeout),
// deferred cleanup may not execute, leaving domains behind. Over many sessions
// these accumulate and push newly created domains beyond the first page of
// GetDomainNames, causing false failures.
func (r *TestRunner) cleanupStaleDomainNames(ctx context.Context, client *apigateway.Client) {
	prefixes := []string{"test-", "none-", "full-lifecycle-"}
	var position *string
	for {
		resp, err := client.GetDomainNames(ctx, &apigateway.GetDomainNamesInput{
			Limit:    aws.Int32(500),
			Position: position,
		})
		if err != nil {
			return
		}
		for _, item := range resp.Items {
			if item.DomainName == nil {
				continue
			}
			dn := *item.DomainName
			for _, p := range prefixes {
				if strings.HasPrefix(dn, p) {
					client.DeleteDomainName(ctx, &apigateway.DeleteDomainNameInput{
						DomainName: aws.String(dn),
					})
					break
				}
			}
		}
		if resp.Position == nil {
			break
		}
		position = resp.Position
	}
}

func (r *TestRunner) runAPIGatewayDomainTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	ctx, client := tc.ctx, tc.client

	// Clean up stale domains from previous interrupted test runs.
	r.cleanupStaleDomainNames(ctx, client)

	// Provision an ACM certificate for all domain tests.
	certArn, err := r.createTestCertArn(ctx)
	if err != nil {
		return []TestResult{{
			Service:  "apigateway",
			TestName: "DomainSetup",
			Status:   "FAIL",
			Error:    err.Error(),
		}}
	}

	var domainName string
	// Ensure the created domain is deleted even if a subsequent test fails
	// and aborts the test runner, preventing test-residue accumulation.
	defer func() {
		if domainName != "" {
			client.DeleteDomainName(ctx, &apigateway.DeleteDomainNameInput{
				DomainName: aws.String(domainName),
			})
		}
	}()
	results = append(results, r.RunTest("apigateway", "CreateDomainName", func() error {
		domain := fmt.Sprintf("test-%d.example.com", time.Now().UnixNano())
		resp, err := client.CreateDomainName(ctx, &apigateway.CreateDomainNameInput{
			DomainName:     aws.String(domain),
			CertificateArn: aws.String(certArn),
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
		domainName = domain
		return nil
	}))
	results = append(results, r.RunTest("apigateway", "GetDomainNames", func() error {
		// Paginate through all pages — accumulated domains from previous
		// test runs can push the newly created domain beyond the first page.
		var position *string
		for {
			resp, err := client.GetDomainNames(ctx, &apigateway.GetDomainNamesInput{
				Limit:    aws.Int32(100),
				Position: position,
			})
			if err != nil {
				return err
			}
			for _, item := range resp.Items {
				if item.DomainName != nil && *item.DomainName == domainName {
					return nil
				}
			}
			if resp.Position == nil {
				break
			}
			position = resp.Position
		}
		return fmt.Errorf("created domain %q not found in list", domainName)
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
		if tc.apiID == "" || domainName == "" {
			return fmt.Errorf("API ID or domain name not available")
		}
		resp, err := client.CreateBasePathMapping(ctx, &apigateway.CreateBasePathMappingInput{
			DomainName: aws.String(domainName),
			RestApiId:  aws.String(tc.apiID),
			BasePath:   aws.String("v1"),
			Stage:      aws.String("prod"),
		})
		if err != nil {
			return err
		}
		if resp.BasePath == nil || *resp.BasePath != "v1" {
			return fmt.Errorf("basePath mismatch, got %v", resp.BasePath)
		}
		if resp.RestApiId == nil || *resp.RestApiId != tc.apiID {
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

	results = append(results, r.RunTest("apigateway", "CreateBasePathMapping_DefaultNone", func() error {
		apiID, _, err := tc.createAPI(tc.uniqueName("BnAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(apiID)

		domain := fmt.Sprintf("none-%d.example.com", time.Now().UnixNano())
		_, err = client.CreateDomainName(ctx, &apigateway.CreateDomainNameInput{
			DomainName:     aws.String(domain),
			CertificateArn: aws.String(certArn),
		})
		if err != nil {
			return fmt.Errorf("create domain: %v", err)
		}
		defer client.DeleteDomainName(ctx, &apigateway.DeleteDomainNameInput{DomainName: aws.String(domain)})

		mappingResp, err := client.CreateBasePathMapping(ctx, &apigateway.CreateBasePathMappingInput{
			DomainName: aws.String(domain),
			RestApiId:  aws.String(apiID),
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
		return nil
	}))

	return results
}
