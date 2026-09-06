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
func (r *TestRunner) cleanupStaleDomainNames(tc *apigwTestContext) {
	prefixes := []string{"test-", "none-", "full-lifecycle-"}
	items, err := tc.allDomainNames()
	if err != nil {
		return
	}
	for _, item := range items {
		if item.DomainName == nil {
			continue
		}
		dn := *item.DomainName
		for _, p := range prefixes {
			if strings.HasPrefix(dn, p) {
				tc.client.DeleteDomainName(tc.ctx, &apigateway.DeleteDomainNameInput{
					DomainName: aws.String(dn),
				})
				break
			}
		}
	}
}

func (r *TestRunner) runAPIGatewayDomainTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	ctx, client := tc.ctx, tc.client

	// Clean up stale domains from previous interrupted test runs.
	r.cleanupStaleDomainNames(tc)

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
		items, err := tc.allDomainNames()
		if err != nil {
			return err
		}
		for _, item := range items {
			if item.DomainName != nil && *item.DomainName == domainName {
				return nil
			}
		}
		return fmt.Errorf("created domain %q not found in list", domainName)
	}))

	results = append(results, r.RunTest("apigateway", "GetDomainName", func() error {
		if err := tc.require(domainName); err != nil {
			return err
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
		if err := tc.require(domainName); err != nil {
			return err
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

	results = append(results, r.RunTest("apigateway", "UpdateDomainName_CertificateAndEndpointPatches", func() error {
		if err := tc.require(domainName); err != nil {
			return err
		}
		// The regionalCertificateName row: replace sets the friendly name.
		_, err := client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/regionalCertificateName"), Value: aws.String("regional-cert")},
			},
		})
		if err != nil {
			return err
		}
		getResp, err := client.GetDomainName(ctx, &apigateway.GetDomainNameInput{DomainName: aws.String(domainName)})
		if err != nil {
			return err
		}
		if aws.ToString(getResp.RegionalCertificateName) != "regional-cert" {
			return fmt.Errorf("regionalCertificateName not set, got %v", getResp.RegionalCertificateName)
		}

		// Removing the edge certificate without a regional one rejects: the
		// documented remove serves the edge-to-regional transition only.
		_, err = client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/certificateArn")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for certificateArn remove without a regional certificate, got: %v", err)
		}

		// The certificate rows exclude add and remove on the same path
		// within one request.
		_, err = client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/certificateName"), Value: aws.String("x")},
				{Op: types.OpRemove, Path: aws.String("/certificateName")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for same-request add+remove, got: %v", err)
		}

		// The endpointConfiguration rows: types replace is Not supported,
		// add serves the edge/regional update; ipAddressType is replace-only
		// with ipv4|dualstack values.
		_, err = client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/endpointConfiguration/types"), Value: aws.String("REGIONAL")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for endpointConfiguration/types replace, got: %v", err)
		}
		_, err = client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/endpointConfiguration/types"), Value: aws.String("REGIONAL")},
				{Op: types.OpReplace, Path: aws.String("/endpointConfiguration/ipAddressType"), Value: aws.String("dualstack")},
			},
		})
		if err != nil {
			return err
		}
		getResp, err = client.GetDomainName(ctx, &apigateway.GetDomainNameInput{DomainName: aws.String(domainName)})
		if err != nil {
			return err
		}
		if getResp.EndpointConfiguration == nil ||
			len(getResp.EndpointConfiguration.Types) != 1 || getResp.EndpointConfiguration.Types[0] != "REGIONAL" {
			return fmt.Errorf("endpoint types not applied, got %+v", getResp.EndpointConfiguration)
		}
		if string(getResp.EndpointConfiguration.IpAddressType) != "dualstack" {
			return fmt.Errorf("ipAddressType not applied, got %+v", getResp.EndpointConfiguration)
		}

		// The developer guide's migration semantics: the new type joins the
		// existing list (its output example shows "types": ["EDGE",
		// "REGIONAL"] with both coexisting until the DNS cutover) and
		// removing the obsolete type later completes the transition.
		_, err = client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/endpointConfiguration/types"), Value: aws.String("EDGE")},
			},
		})
		if err != nil {
			return err
		}
		getResp, err = client.GetDomainName(ctx, &apigateway.GetDomainNameInput{DomainName: aws.String(domainName)})
		if err != nil {
			return err
		}
		if getResp.EndpointConfiguration == nil || len(getResp.EndpointConfiguration.Types) != 2 {
			return fmt.Errorf("endpoint types add did not append, got %+v", getResp.EndpointConfiguration)
		}
		_, err = client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/endpointConfiguration/types"), Value: aws.String("REGIONAL")},
			},
		})
		if err != nil {
			return err
		}
		getResp, err = client.GetDomainName(ctx, &apigateway.GetDomainNameInput{DomainName: aws.String(domainName)})
		if err != nil {
			return err
		}
		if getResp.EndpointConfiguration == nil ||
			len(getResp.EndpointConfiguration.Types) != 1 || getResp.EndpointConfiguration.Types[0] != "EDGE" {
			return fmt.Errorf("endpoint types remove did not drop the addressed value, got %+v", getResp.EndpointConfiguration)
		}
		_, err = client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/endpointConfiguration/ipAddressType"), Value: aws.String("ipv6")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for ipAddressType ipv6, got: %v", err)
		}

		// The documented edge-to-regional transition: with a regional
		// certificate present, the edge certificate remove succeeds.
		regionalArn, err := r.createTestCertArn(ctx)
		if err != nil {
			return fmt.Errorf("import regional certificate: %v", err)
		}
		_, err = client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/regionalCertificateArn"), Value: aws.String(regionalArn)},
			},
		})
		if err != nil {
			return err
		}
		_, err = client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/certificateArn")},
			},
		})
		if err != nil {
			return err
		}
		getResp, err = client.GetDomainName(ctx, &apigateway.GetDomainNameInput{DomainName: aws.String(domainName)})
		if err != nil {
			return err
		}
		if getResp.CertificateArn != nil {
			return fmt.Errorf("edge certificate not cleared, got %v", *getResp.CertificateArn)
		}
		if aws.ToString(getResp.RegionalCertificateArn) != regionalArn {
			return fmt.Errorf("regional certificate mismatch, got %v", getResp.RegionalCertificateArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateBasePathMapping", func() error {
		if err := tc.require(tc.apiID, domainName); err != nil {
			return err
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
		if err := tc.require(domainName); err != nil {
			return err
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
		if err := tc.require(domainName); err != nil {
			return err
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
		if err := tc.require(domainName); err != nil {
			return err
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

		// The restApiId row accepts the table spelling ("/restapiId") and
		// the member casing alike; the basePath row accepts both spellings
		// too ("path='/basePath'" is the official CLI example form).
		currentAPI := aws.ToString(resp.RestApiId)
		resp, err = client.UpdateBasePathMapping(ctx, &apigateway.UpdateBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/restapiId"), Value: aws.String(currentAPI)},
			},
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.RestApiId) != currentAPI {
			return fmt.Errorf("restapiId replace not applied, got %v", resp.RestApiId)
		}

		resp, err = client.UpdateBasePathMapping(ctx, &apigateway.UpdateBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/basePath"), Value: aws.String("v2")},
			},
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.BasePath) != "v2" {
			return fmt.Errorf("basePath rename not applied, got %v", resp.BasePath)
		}
		if _, err = client.GetBasePathMapping(ctx, &apigateway.GetBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v2"),
		}); err != nil {
			return fmt.Errorf("renamed mapping not readable: %v", err)
		}
		// Rename back so the delete test keeps addressing "v1".
		_, err = client.UpdateBasePathMapping(ctx, &apigateway.UpdateBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v2"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/basepath"), Value: aws.String("v1")},
			},
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteBasePathMapping", func() error {
		if err := tc.require(domainName); err != nil {
			return err
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
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetBasePathMapping should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteDomainName", func() error {
		if err := tc.require(domainName); err != nil {
			return err
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
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetDomainName should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateBasePathMapping_DefaultNone", func() error {
		apiID, _, err := tc.createOwnAPI("BnAPI")
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

	results = append(results, r.RunTest("apigateway", "CreateDomainName_DuplicateConflict", func() error {
		domain := fmt.Sprintf("dup-%d.example.com", time.Now().UnixNano())
		_, err := client.CreateDomainName(ctx, &apigateway.CreateDomainNameInput{
			DomainName:     aws.String(domain),
			CertificateArn: aws.String(certArn),
		})
		if err != nil {
			return fmt.Errorf("create domain: %v", err)
		}
		defer client.DeleteDomainName(ctx, &apigateway.DeleteDomainNameInput{DomainName: aws.String(domain)})

		_, err = client.CreateDomainName(ctx, &apigateway.CreateDomainNameInput{
			DomainName:     aws.String(domain),
			CertificateArn: aws.String(certArn),
		})
		if aerr := AssertErrorContains(err, "ConflictException"); aerr != nil {
			return fmt.Errorf("duplicate CreateDomainName should fail with ConflictException: %v", aerr)
		}
		return nil
	}))

	return results
}
