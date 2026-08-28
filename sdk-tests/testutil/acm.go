package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"vorpalstacks-sdk-tests/config"
)

type acmTestContext struct {
	client    *acm.Client
	ctx       context.Context
	region    string
	accountID string
}

func (r *TestRunner) newACMTestContext() (*acmTestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &acmTestContext{
		client:    acm.NewFromConfig(cfg),
		ctx:       context.Background(),
		region:    r.region,
		accountID: r.accountID,
	}, nil
}

func acmUniqueDomain(prefix string) string {
	return fmt.Sprintf("%s-%d.com", prefix, time.Now().UnixNano())
}

// requestCert submits a RequestCertificate call and returns the new
// certificate's ARN.
func (tc *acmTestContext) requestCert(input *acm.RequestCertificateInput) (string, error) {
	resp, err := tc.client.RequestCertificate(tc.ctx, input)
	if err != nil {
		return "", err
	}
	arn := aws.ToString(resp.CertificateArn)
	if arn == "" {
		return "", fmt.Errorf("expected CertificateArn in RequestCertificate response")
	}
	return arn, nil
}

// requestDNSCert requests a DNS-validated certificate and returns its ARN.
func (tc *acmTestContext) requestDNSCert(domain string) (string, error) {
	return tc.requestCert(&acm.RequestCertificateInput{
		DomainName:       aws.String(domain),
		ValidationMethod: types.ValidationMethodDns,
	})
}

// requestEmailCert requests an email-validated certificate and returns its ARN.
func (tc *acmTestContext) requestEmailCert(domain string) (string, error) {
	return tc.requestCert(&acm.RequestCertificateInput{
		DomainName:       aws.String(domain),
		ValidationMethod: types.ValidationMethodEmail,
	})
}

// requestCertSANs requests a DNS-validated certificate with subject
// alternative names and returns its ARN.
func (tc *acmTestContext) requestCertSANs(domain string, sans ...string) (string, error) {
	return tc.requestCert(&acm.RequestCertificateInput{
		DomainName:              aws.String(domain),
		ValidationMethod:        types.ValidationMethodDns,
		SubjectAlternativeNames: sans,
	})
}

func (tc *acmTestContext) importDefaultCert() (string, error) {
	return tc.importCertificate(testCertPEM, testKeyPEM, nil)
}

func (tc *acmTestContext) importCertWithChain() (string, error) {
	return tc.importCertificate(testCertPEM, testKeyPEM, testChainPEM)
}

// importCertificate imports a certificate body and returns the new
// certificate's ARN.
func (tc *acmTestContext) importCertificate(certificate, privateKey, chain []byte) (string, error) {
	input := &acm.ImportCertificateInput{
		Certificate: certificate,
		PrivateKey:  privateKey,
	}
	if chain != nil {
		input.CertificateChain = chain
	}
	resp, err := tc.client.ImportCertificate(tc.ctx, input)
	if err != nil {
		return "", err
	}
	arn := aws.ToString(resp.CertificateArn)
	if arn == "" {
		return "", fmt.Errorf("expected CertificateArn in ImportCertificate response")
	}
	return arn, nil
}

// allCertificates walks ListCertificates to completion collecting every
// certificate summary across all pages. Optional statuses filter the list.
func (tc *acmTestContext) allCertificates(statuses []types.CertificateStatus) ([]types.CertificateSummary, error) {
	return paginate(func(next *string) ([]types.CertificateSummary, *string, error) {
		resp, err := tc.client.ListCertificates(tc.ctx, &acm.ListCertificatesInput{
			MaxItems:            aws.Int32(100),
			CertificateStatuses: statuses,
			NextToken:           next,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("list certificates page: %v", err)
		}
		return resp.CertificateSummaryList, resp.NextToken, nil
	})
}

// cleanupStaleCertificates removes certificates left over from previous test
// runs. These accumulate across sessions and eventually fill the first page
// of ListCertificates, causing pagination-dependent tests to fail.
func (tc *acmTestContext) cleanupStaleCertificates() {
	certs, err := tc.allCertificates(nil)
	if err != nil {
		return
	}
	for _, s := range certs {
		domain := aws.ToString(s.DomainName)
		if strings.HasPrefix(domain, "page-") || strings.HasPrefix(domain, "summary-") || strings.HasPrefix(domain, "test-") {
			tc.deleteCert(aws.ToString(s.CertificateArn))
		}
	}
}

// nonExistentArn builds the ARN of a certificate that never exists, for
// negative-path tests pinning ResourceNotFoundException.
func (tc *acmTestContext) nonExistentArn() string {
	return fmt.Sprintf("arn:aws:acm:%s:%s:certificate/nonexistent", tc.region, tc.accountID)
}

func (tc *acmTestContext) deleteCert(arn string) {
	tc.client.DeleteCertificate(tc.ctx, &acm.DeleteCertificateInput{
		CertificateArn: aws.String(arn),
	})
}

func (r *TestRunner) RunACMTests() []TestResult {
	tc, err := r.newACMTestContext()
	if err != nil {
		return []TestResult{{
			Service:  "acm",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		}}
	}

	tc.cleanupStaleCertificates()

	var results []TestResult
	results = append(results, r.runACMCertificateTests(tc)...)
	results = append(results, r.runACMLifecycleTests(tc)...)
	results = append(results, r.runACMListTests(tc)...)
	results = append(results, r.runACMTagTests(tc)...)
	results = append(results, r.runACMAccountTests(tc)...)
	results = append(results, r.runACMEdgeTests(tc)...)
	return results
}
