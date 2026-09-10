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

// deleteTestCertArn removes an ACM certificate imported for one test run,
// so repeated runs do not accumulate imported certificates.
func (r *TestRunner) deleteTestCertArn(ctx context.Context, arn string) {
	if arn == "" {
		return
	}
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return
	}
	_, _ = acm.NewFromConfig(cfg).DeleteCertificate(ctx, &acm.DeleteCertificateInput{
		CertificateArn: aws.String(arn),
	})
}

// cleanupStaleDomainNames deletes accumulated test-residue domain names from
// previous test runs. When the test runner is interrupted (SIGKILL, timeout),
// deferred cleanup may not execute, leaving domains behind. Over many sessions
// these accumulate and push newly created domains beyond the first page of
// GetDomainNames, causing false failures.
func (r *TestRunner) cleanupStaleDomainNames(tc *apigwTestContext) {
	prefixes := []string{"test-", "none-", "full-lifecycle-", "dup-"}
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
	// The imported certificate is per-run state: remove it once the domain
	// release defer above has run (defers run last-registered-first).
	defer r.deleteTestCertArn(ctx, certArn)
	results = append(results, r.RunTest("apigateway", "CreateDomainName", func() error {
		domain := fmt.Sprintf("test-%d.example.com", time.Now().UnixNano())
		resp, err := client.CreateDomainName(ctx, &apigateway.CreateDomainNameInput{
			DomainName:     aws.String(domain),
			CertificateArn: aws.String(certArn),
			EndpointConfiguration: &types.EndpointConfiguration{
				Types: []types.EndpointType{"EDGE"},
			},
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

	results = append(results, r.RunTest("apigateway", "GetDomainNames_ResourceOwner", func() error {
		if err := tc.require(domainName); err != nil {
			return err
		}
		// Every domain on this single-account platform is SELF-owned;
		// OTHER_ACCOUNTS matches none.
		selfOwned, err := tc.client.GetDomainNames(tc.ctx, &apigateway.GetDomainNamesInput{
			ResourceOwner: types.ResourceOwnerSelf,
			Limit:         aws.Int32(500),
		})
		if err != nil {
			return err
		}
		found := false
		for _, item := range selfOwned.Items {
			if aws.ToString(item.DomainName) == domainName {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("SELF listing missing the created domain, got %d items", len(selfOwned.Items))
		}
		otherAccounts, err := tc.client.GetDomainNames(tc.ctx, &apigateway.GetDomainNamesInput{
			ResourceOwner: types.ResourceOwnerOtherAccounts,
			Limit:         aws.Int32(500),
		})
		if err != nil {
			return err
		}
		if len(otherAccounts.Items) != 0 || otherAccounts.Position != nil {
			return fmt.Errorf("OTHER_ACCOUNTS matched %d domains on a single-account platform", len(otherAccounts.Items))
		}
		return nil
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
		if aws.ToString(resp.CertificateArn) != certArn {
			return fmt.Errorf("certificateArn mismatch, got %v", resp.CertificateArn)
		}
		if resp.Tags["domain"] != "test" {
			return fmt.Errorf("tags mismatch, got %v", resp.Tags)
		}
		if resp.EndpointConfiguration == nil || len(resp.EndpointConfiguration.Types) == 0 ||
			resp.EndpointConfiguration.Types[0] != types.EndpointTypeEdge {
			return fmt.Errorf("edge endpointConfiguration mismatch, got %+v", resp.EndpointConfiguration)
		}
		if p := string(resp.SecurityPolicy); p != "TLS_1_2" {
			return fmt.Errorf("securityPolicy must default to TLS_1_2 on creation, got %q", p)
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

		// The securityPolicy row: replace between the documented policies.
		spResp, err := client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/securityPolicy"), Value: aws.String("TLS_1_2")},
			},
		})
		if err != nil {
			return fmt.Errorf("securityPolicy replace: %v", err)
		}
		if string(spResp.SecurityPolicy) != "TLS_1_2" {
			return fmt.Errorf("securityPolicy not applied, got %v", spResp.SecurityPolicy)
		}

		// The mutual TLS rows: truststoreUri and truststoreVersion accept
		// add, replace and remove.
		mtlsResp, err := client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/mutualTlsAuthentication/truststoreUri"), Value: aws.String("s3://my-bucket/truststore.pem")},
				{Op: types.OpAdd, Path: aws.String("/mutualTlsAuthentication/truststoreVersion"), Value: aws.String("2")},
			},
		})
		if err != nil {
			return fmt.Errorf("mutualTls add: %v", err)
		}
		if mtlsResp.MutualTlsAuthentication == nil ||
			aws.ToString(mtlsResp.MutualTlsAuthentication.TruststoreUri) != "s3://my-bucket/truststore.pem" ||
			aws.ToString(mtlsResp.MutualTlsAuthentication.TruststoreVersion) != "2" {
			return fmt.Errorf("mutualTls rows not applied, got %+v", mtlsResp.MutualTlsAuthentication)
		}
		_, err = client.UpdateDomainName(ctx, &apigateway.UpdateDomainNameInput{
			DomainName: aws.String(domainName),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/mutualTlsAuthentication/truststoreUri")},
				{Op: types.OpRemove, Path: aws.String("/mutualTlsAuthentication/truststoreVersion")},
			},
		})
		if err != nil {
			return fmt.Errorf("mutualTls remove: %v", err)
		}
		mtlsGet, err := client.GetDomainName(ctx, &apigateway.GetDomainNameInput{DomainName: aws.String(domainName)})
		if err != nil {
			return err
		}
		if mtlsGet.MutualTlsAuthentication != nil &&
			(aws.ToString(mtlsGet.MutualTlsAuthentication.TruststoreUri) != "" ||
				aws.ToString(mtlsGet.MutualTlsAuthentication.TruststoreVersion) != "") {
			return fmt.Errorf("mutualTls remove not applied, got %+v", mtlsGet.MutualTlsAuthentication)
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
		// The domain was created EDGE, so adding REGIONAL follows the
		// documented migration semantics: both types coexist until the
		// cutover remove completes the transition.
		if getResp.EndpointConfiguration == nil ||
			len(getResp.EndpointConfiguration.Types) != 2 {
			return fmt.Errorf("endpoint types add did not append REGIONAL to EDGE, got %+v", getResp.EndpointConfiguration)
		}
		sawEdge, sawRegional := false, false
		for _, t := range getResp.EndpointConfiguration.Types {
			if t == types.EndpointTypeEdge {
				sawEdge = true
			}
			if t == types.EndpointTypeRegional {
				sawRegional = true
			}
		}
		if !sawEdge || !sawRegional {
			return fmt.Errorf("expected EDGE and REGIONAL to coexist, got %+v", getResp.EndpointConfiguration)
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
		defer r.deleteTestCertArn(ctx, regionalArn)
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

		// A forced multi-page walk: two more mappings bring the domain past
		// a two-mapping page limit, with every page's position fed back.
		for _, bp := range []string{"pg1", "pg2"} {
			if _, err := client.CreateBasePathMapping(ctx, &apigateway.CreateBasePathMappingInput{
				DomainName: aws.String(domainName),
				RestApiId:  aws.String(tc.apiID),
				BasePath:   aws.String(bp),
				Stage:      aws.String("prod"),
			}); err != nil {
				return fmt.Errorf("create mapping %s: %v", bp, err)
			}
			defer client.DeleteBasePathMapping(ctx, &apigateway.DeleteBasePathMappingInput{
				DomainName: aws.String(domainName), BasePath: aws.String(bp),
			})
		}
		want := map[string]bool{"v1": true, "pg1": true, "pg2": true}
		got := map[string]bool{}
		var position *string
		pages := 0
		for {
			page, err := client.GetBasePathMappings(ctx, &apigateway.GetBasePathMappingsInput{
				DomainName: aws.String(domainName),
				Limit:      aws.Int32(2),
				Position:   position,
			})
			if err != nil {
				return fmt.Errorf("page %d: %v", pages, err)
			}
			pages++
			for _, m := range page.Items {
				got[aws.ToString(m.BasePath)] = true
			}
			position = page.Position
			if position == nil {
				break
			}
		}
		if pages < 2 {
			return fmt.Errorf("expected the page limit to force multiple pages, got %d", pages)
		}
		for bp := range want {
			if !got[bp] {
				return fmt.Errorf("created mapping %q missing from the walk, got %v", bp, got)
			}
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

		// A real retarget: the documented example flow moves the mapping to
		// another API; moving it back keeps the later tests addressing this
		// API through "v1".
		otherAPI, _, err := tc.createOwnAPI("BpmRetgt")
		if err != nil {
			return fmt.Errorf("create retarget api: %v", err)
		}
		defer tc.deleteAPI(otherAPI)
		resp, err = client.UpdateBasePathMapping(ctx, &apigateway.UpdateBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/restapiId"), Value: aws.String(otherAPI)},
			},
		})
		if err != nil {
			return fmt.Errorf("retarget to another api: %v", err)
		}
		if aws.ToString(resp.RestApiId) != otherAPI {
			return fmt.Errorf("restapiId retarget not applied, got %v", resp.RestApiId)
		}
		resp, err = client.UpdateBasePathMapping(ctx, &apigateway.UpdateBasePathMappingInput{
			DomainName: aws.String(domainName),
			BasePath:   aws.String("v1"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/restapiId"), Value: aws.String(currentAPI)},
			},
		})
		if err != nil {
			return fmt.Errorf("retarget back: %v", err)
		}
		if aws.ToString(resp.RestApiId) != currentAPI {
			return fmt.Errorf("restapiId restore not applied, got %v", resp.RestApiId)
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
