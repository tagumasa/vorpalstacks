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

var testCertPEM = []byte(`-----BEGIN CERTIFICATE-----
MIICnzCCAYegAwIBAgIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDEwh0ZXN0
Y2VydDAeFw0yNjAzMjUwNzE1MTVaFw0yNzAzMjUwNzE1MTVaMBMxETAPBgNVBAMT
CHRlc3RjZXJ0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuQ1frJ4y
NJURfuVZ+ZoXaJnzH9Aca7cAl4kUlZauacQe9GeBiK9MH/gZahS5Nk7uYB3SEFf2
hRFy5O0FOhk89rztdB/iWZn346+RqRHAxBEl1LGRX0HTCaaf/uxl8uj6qraDJrOm
rCaBAU3zBQ+x7xJO0GmYT4y2rsnDdJwnElVIcNW6EcF/e7mN5F8qItLuNvLeZcgI
CEifF1Jxhj6/0LnOB2ywsrvs974lIDfvOs8wbkQJZIOZX7TOkwtUNo9FaBua5a8s
Q03SXxas6nXXBHE7yl/BlJZfneAO8KT1w067ohWpuPjCGfJN6LgXg347nE5IgyFM
gksV2rXM9SdkowIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBCHYc/ZkJBo6m8G4I8
3/u2joYJAgo0MpsQiKre1lRuEgvsWHFbyPMBWXQkGdTydV8AIz23YV+rpPDt3/s/
BliGOu4L4o2bCjiPO5V2cv36id6e7FRfJyAmRe/S3M06jJh9HB3/uUTABITkGgee
Sa35wq1cRp86PGHhCGkEg79J8WRQmNrelttmCz/Fs4N5leuwnOlTlgCoEaLt+QSY
1DR2aPlMB0iC7yQ2UMSwdLvdWQ7ted02yYV0Hqgq/QT3wA7vfjI0SG0OUqfaJ5d2
QOl0rfDrYF2ZQNqiUX827TRg9kYRJveMjGxLhFMNVxyZJkQsbGoxJPIMikWULfk2
Xwdo
-----END CERTIFICATE-----`)

var testKeyPEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAuQ1frJ4yNJURfuVZ+ZoXaJnzH9Aca7cAl4kUlZauacQe9GeB
iK9MH/gZahS5Nk7uYB3SEFf2hRFy5O0FOhk89rztdB/iWZn346+RqRHAxBEl1LGR
X0HTCaaf/uxl8uj6qraDJrOmrCaBAU3zBQ+x7xJO0GmYT4y2rsnDdJwnElVIcNW6
EcF/e7mN5F8qItLuNvLeZcgICEifF1Jxhj6/0LnOB2ywsrvs974lIDfvOs8wbkQJ
ZIOZX7TOkwtUNo9FaBua5a8sQ03SXxas6nXXBHE7yl/BlJZfneAO8KT1w067ohWp
uPjCGfJN6LgXg347nE5IgyFMgksV2rXM9SdkowIDAQABAoIBAEseUi2kxBWTQ5hi
6szHT+ROxiIuXTMehPd+lmQI2EEn8zbcQ3lkS38Yu9xTkEGq9dn/kPPAeVpYFG84
hewpLZWtaKjAfqZHuZhr/zGF+t28ZkJ6WFw2QMBEquMVPGdISuT8lK2jtK9iK/EH
HvT5g43cPTEeBE2afdfjIFwYPUYTto2bC1dIsPJ66IH3AUN6uwnYLfIlyomvIxGJ
iwsNZloOEMtjpvf8Q/5JbfioTYwBMGS4SZetPl4CSnASLI44jZPU7hWHDRhlM0OS
U3TzqacbAxNm1tzkzJARxyCd9GatyuLqNSgph0QqW/VXkOH10kCMoGRUmmMc6gWe
40yaH4kCgYEAziBz4/RnEs+WMqs0IwoKQtgU9blXTXNgIs9WS0Eyn3XgigZB9IxB
BIJHbltDSeX5/TO3iyhE0hIeEDukSsuDzMt+O2N2ZOac+4UnRcqczD12XiLLysru
mfep9MUNrIUj+UaMB1ZPyVfxGyfIc8An9RlayBN7jsDZi+Pj3dTSfpsCgYEA5dOP
JVTGgC0ZcK7w+xev2iCDixHMvX4ofm4iKd9eYM3RXKnUulDbTL4GcACjbW82IG5z
0TfEdAF4lNW3C1bCSDhIWM1P3Nc1zPnH65RZju3oSYvToDk/1PXcWTmCcWXA2twq
JE8NRBaHtFjBMqu/5KddcIoohIlTRiC/V/d7zpkCgYEAg88UzIwY7Vp5PWVlLZLa
BOyQWqFuRkSlER1snSrP6FBEiX5+5pZZbTyx2MvbN4IsXdGYaRATEhIrz02UPY/u
dCMcUXXE27jsYZpABs0Nfz0+V+wATWl/Mk3BDJiFqfBplJmcKYTz+FiYATlrYTlb
U8wm1RJATITdmCreJ5hUEkkCgYEAhqlvNnB13qSOQ3g9uuImJ6jlapcDYASLtYjS
e7ZlllMCWUkpXAIEfPLa0sWM/JItJNOTCQOkGFTEUnDmz74GGEriGSYzpTJ0U6YH
fgFueFDtyioj1b21qRJmCeGojMkSNyrJhnzLSRnqacGXchkwVsm59jb9hqrwICcP
9nsMEAECgYB498ktMUMajMgNyKc4bIL92EzScPcTIfn+1a22wd0ZJkiEtTotMwPh
Nw1sf/uZ5JyJwTEr6FU4qBk+zc/M3+4f8VG5ChVMt6mPEVwHAlgscwODj1pxO7nz
Vzw7YxT498cnLJsBFDy+kk9uKMf7cpLCdRF1gRpeIP3K6sFLNF96Gw==
-----END RSA PRIVATE KEY-----`)

var testChainPEM = []byte(`-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHHCgVZU1JUMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3RjYTAeFw0yNDAxMDEwMDAwMDBaFw0yNTAxMDEwMDAwMDBaMBExDzANBgNVBAMM
BnRlc3RjYTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAwK0j6f8C6hJ7u8P
-----END CERTIFICATE-----`)

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
