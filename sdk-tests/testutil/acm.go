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

func (r *TestRunner) RunACMTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "acm",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := acm.NewFromConfig(cfg)
	ctx := context.Background()

	testCertPEM := []byte(`-----BEGIN CERTIFICATE-----
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

	testKeyPEM := []byte(`-----BEGIN RSA PRIVATE KEY-----
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

	testChainPEM := []byte(`-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHHCgVZU1JUMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3RjYTAeFw0yNDAxMDEwMDAwMDBaFw0yNTAxMDEwMDAwMDBaMBExDzANBgNVBAMM
BnRlc3RjYTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAwK0j6f8C6hJ7u8P
-----END CERTIFICATE-----`)

	// --- RequestCertificate ---

	results = append(results, r.RunTest("acm", "RequestCertificate_WithSubjectAlternativeNames", func() error {
		domain := fmt.Sprintf("san-test-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:              aws.String(domain),
			ValidationMethod:        types.ValidationMethodDns,
			SubjectAlternativeNames: []string{fmt.Sprintf("www.%s", domain), fmt.Sprintf("api.%s", domain)},
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})
		desc, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(desc.Certificate.SubjectAlternativeNames) != 2 {
			return fmt.Errorf("expected 2 SANs, got %d", len(desc.Certificate.SubjectAlternativeNames))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RequestCertificate_WithOptions", func() error {
		domain := fmt.Sprintf("opts-test-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
			Options: &types.CertificateOptions{
				CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled,
			},
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})
		desc, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if desc.Certificate.Options == nil {
			return fmt.Errorf("Options is nil")
		}
		if desc.Certificate.Options.CertificateTransparencyLoggingPreference != types.CertificateTransparencyLoggingPreferenceDisabled {
			return fmt.Errorf("expected DISABLED, got %s", desc.Certificate.Options.CertificateTransparencyLoggingPreference)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RequestCertificate_WithTags", func() error {
		domain := fmt.Sprintf("tag-test-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
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
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})
		tagResp, err := client.ListTagsForCertificate(ctx, &acm.ListTagsForCertificateInput{CertificateArn: resp.CertificateArn})
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

	// --- DescribeCertificate ---

	results = append(results, r.RunTest("acm", "DescribeCertificate_AMAZON_ISSUED_Fields", func() error {
		domain := fmt.Sprintf("desc-ai-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})
		desc, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		c := desc.Certificate
		if c.Status != types.CertificateStatusIssued {
			return fmt.Errorf("expected ISSUED, got %s", c.Status)
		}
		if c.Type != types.CertificateTypeAmazonIssued {
			return fmt.Errorf("expected AMAZON_ISSUED, got %s", c.Type)
		}
		if c.RenewalEligibility != types.RenewalEligibilityEligible {
			return fmt.Errorf("expected ELIGIBLE, got %s", c.RenewalEligibility)
		}
		if aws.ToString(c.DomainName) != domain {
			return fmt.Errorf("expected domain %s, got %s", domain, aws.ToString(c.DomainName))
		}
		if c.KeyAlgorithm != types.KeyAlgorithmRsa2048 {
			return fmt.Errorf("expected RSA_2048, got %s", c.KeyAlgorithm)
		}
		if c.CertificateArn == nil || !strings.Contains(aws.ToString(c.CertificateArn), "acm") {
			return fmt.Errorf("CertificateArn missing or malformed: %s", aws.ToString(c.CertificateArn))
		}
		if c.CreatedAt == nil {
			return fmt.Errorf("CreatedAt is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "DescribeCertificate_DomainValidationOptions_DNS", func() error {
		domain := fmt.Sprintf("dv-dns-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})
		desc, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(desc.Certificate.DomainValidationOptions) != 1 {
			return fmt.Errorf("expected 1 DVO, got %d", len(desc.Certificate.DomainValidationOptions))
		}
		dvo := desc.Certificate.DomainValidationOptions[0]
		if dvo.ValidationMethod != types.ValidationMethodDns {
			return fmt.Errorf("expected DNS, got %s", dvo.ValidationMethod)
		}
		if dvo.ResourceRecord == nil {
			return fmt.Errorf("ResourceRecord is nil")
		}
		if dvo.ResourceRecord.Type != "CNAME" {
			return fmt.Errorf("expected CNAME, got %s", dvo.ResourceRecord.Type)
		}
		if !strings.Contains(aws.ToString(dvo.ResourceRecord.Name), domain) {
			return fmt.Errorf("ResourceRecord.Name should contain domain, got %s", aws.ToString(dvo.ResourceRecord.Name))
		}
		if aws.ToString(dvo.ResourceRecord.Value) == "" {
			return fmt.Errorf("ResourceRecord.Value is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "DescribeCertificate_IMPORTED_Fields", func() error {
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  testKeyPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: importResp.CertificateArn})
		desc, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: importResp.CertificateArn})
		if err != nil {
			return err
		}
		c := desc.Certificate
		if c.Status != types.CertificateStatusIssued {
			return fmt.Errorf("expected ISSUED, got %s", c.Status)
		}
		if c.Type != types.CertificateTypeImported {
			return fmt.Errorf("expected IMPORTED, got %s", c.Type)
		}
		if c.RenewalEligibility != types.RenewalEligibilityIneligible {
			return fmt.Errorf("expected INELIGIBLE, got %s", c.RenewalEligibility)
		}
		if c.ImportedAt == nil {
			return fmt.Errorf("ImportedAt is nil")
		}
		if c.NotBefore == nil || c.NotAfter == nil {
			return fmt.Errorf("NotBefore or NotAfter is nil")
		}
		if c.KeyAlgorithm != types.KeyAlgorithmRsa2048 {
			return fmt.Errorf("expected RSA_2048, got %s", c.KeyAlgorithm)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "DescribeCertificate_NonExistent", func() error {
		_, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// --- GetCertificate ---

	results = append(results, r.RunTest("acm", "GetCertificate_ImportedCert", func() error {
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate:      testCertPEM,
			PrivateKey:       testKeyPEM,
			CertificateChain: testChainPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: importResp.CertificateArn})
		getResp, err := client.GetCertificate(ctx, &acm.GetCertificateInput{CertificateArn: importResp.CertificateArn})
		if err != nil {
			return err
		}
		if aws.ToString(getResp.Certificate) == "" {
			return fmt.Errorf("Certificate is empty")
		}
		if aws.ToString(getResp.CertificateChain) == "" {
			return fmt.Errorf("CertificateChain is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "GetCertificate_NonExistent", func() error {
		_, err := client.GetCertificate(ctx, &acm.GetCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// --- ListCertificates ---

	results = append(results, r.RunTest("acm", "ListCertificates", func() error {
		resp, err := client.ListCertificates(ctx, &acm.ListCertificatesInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.CertificateSummaryList == nil {
			return fmt.Errorf("certificate summary list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ListCertificates_Pagination", func() error {
		var arns []string
		for i := 0; i < 3; i++ {
			resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
				DomainName:       aws.String(fmt.Sprintf("page-%d-%d.com", time.Now().UnixNano(), i)),
				ValidationMethod: types.ValidationMethodDns,
			})
			if err != nil {
				return err
			}
			arns = append(arns, aws.ToString(resp.CertificateArn))
		}
		defer func() {
			for _, arn := range arns {
				client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: aws.String(arn)})
			}
		}()

		resp, err := client.ListCertificates(ctx, &acm.ListCertificatesInput{
			MaxItems: aws.Int32(2),
		})
		if err != nil {
			return err
		}
		if len(resp.CertificateSummaryList) != 2 {
			return fmt.Errorf("expected 2, got %d", len(resp.CertificateSummaryList))
		}
		if resp.NextToken == nil || aws.ToString(resp.NextToken) == "" {
			return fmt.Errorf("expected non-empty NextToken")
		}

		resp2, err := client.ListCertificates(ctx, &acm.ListCertificatesInput{
			NextToken: resp.NextToken,
			MaxItems:  aws.Int32(2),
		})
		if err != nil {
			return err
		}
		if len(resp2.CertificateSummaryList) < 1 {
			return fmt.Errorf("expected at least 1 on page 2, got %d", len(resp2.CertificateSummaryList))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ListCertificates_CertificateStatusesFilter", func() error {
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  testKeyPEM,
		})
		if err != nil {
			return err
		}
		importArn := aws.ToString(importResp.CertificateArn)
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: importResp.CertificateArn})

		resp, err := client.ListCertificates(ctx, &acm.ListCertificatesInput{
			CertificateStatuses: []types.CertificateStatus{types.CertificateStatusIssued},
		})
		if err != nil {
			return err
		}
		found := false
		for _, s := range resp.CertificateSummaryList {
			if aws.ToString(s.CertificateArn) == importArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("imported ISSUED cert not found in filtered list")
		}
		for _, s := range resp.CertificateSummaryList {
			if s.Status != types.CertificateStatusIssued {
				return fmt.Errorf("found non-ISSUED cert in filtered list: %s", s.Status)
			}
		}
		return nil
	}))

	// --- DeleteCertificate ---

	results = append(results, r.RunTest("acm", "DeleteCertificate_NonExistent", func() error {
		_, err := client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// --- Tags ---

	results = append(results, r.RunTest("acm", "AddTagsToCertificate_UpdateExistingTag", func() error {
		domain := fmt.Sprintf("tagupd-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})

		client.AddTagsToCertificate(ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: resp.CertificateArn,
			Tags:           []types.Tag{{Key: aws.String("Env"), Value: aws.String("dev")}},
		})
		client.AddTagsToCertificate(ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: resp.CertificateArn,
			Tags:           []types.Tag{{Key: aws.String("Env"), Value: aws.String("prod")}},
		})
		tagResp, err := client.ListTagsForCertificate(ctx, &acm.ListTagsForCertificateInput{CertificateArn: resp.CertificateArn})
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
		return fmt.Errorf("Env tag not found")
	}))

	results = append(results, r.RunTest("acm", "AddTagsToCertificate_VerifyContent", func() error {
		domain := fmt.Sprintf("tagver-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})

		client.AddTagsToCertificate(ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: resp.CertificateArn,
			Tags: []types.Tag{
				{Key: aws.String("Key1"), Value: aws.String("Val1")},
				{Key: aws.String("Key2"), Value: aws.String("Val2")},
			},
		})
		tagResp, err := client.ListTagsForCertificate(ctx, &acm.ListTagsForCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(tagResp.Tags) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(tagResp.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RemoveTagsFromCertificate_VerifyEmpty", func() error {
		domain := fmt.Sprintf("tagrm-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})

		client.AddTagsToCertificate(ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: resp.CertificateArn,
			Tags:           []types.Tag{{Key: aws.String("X"), Value: aws.String("Y")}},
		})
		client.RemoveTagsFromCertificate(ctx, &acm.RemoveTagsFromCertificateInput{
			CertificateArn: resp.CertificateArn,
			Tags:           []types.Tag{{Key: aws.String("X")}},
		})
		tagResp, err := client.ListTagsForCertificate(ctx, &acm.ListTagsForCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(tagResp.Tags) != 0 {
			return fmt.Errorf("expected 0 tags after removal, got %d", len(tagResp.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "AddTagsToCertificate_NonExistent", func() error {
		_, err := client.AddTagsToCertificate(ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
			Tags:           []types.Tag{{Key: aws.String("X"), Value: aws.String("Y")}},
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RemoveTagsFromCertificate_NonExistent", func() error {
		_, err := client.RemoveTagsFromCertificate(ctx, &acm.RemoveTagsFromCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
			Tags:           []types.Tag{{Key: aws.String("X")}},
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ListTagsForCertificate_NonExistent", func() error {
		_, err := client.ListTagsForCertificate(ctx, &acm.ListTagsForCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// --- ResendValidationEmail ---

	results = append(results, r.RunTest("acm", "ResendValidationEmail", func() error {
		domain := fmt.Sprintf("resend-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodEmail,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})
		_, err = client.ResendValidationEmail(ctx, &acm.ResendValidationEmailInput{
			CertificateArn:   resp.CertificateArn,
			Domain:           aws.String(domain),
			ValidationDomain: aws.String(domain),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ResendValidationEmail_IssuedStatus", func() error {
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  testKeyPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: importResp.CertificateArn})
		_, err = client.ResendValidationEmail(ctx, &acm.ResendValidationEmailInput{
			CertificateArn:   importResp.CertificateArn,
			Domain:           aws.String("example.com"),
			ValidationDomain: aws.String("example.com"),
		})
		if err := AssertErrorContains(err, "InvalidStateException"); err != nil {
			return err
		}
		return nil
	}))

	// --- UpdateCertificateOptions ---

	results = append(results, r.RunTest("acm", "UpdateCertificateOptions_VerifyInDescribe", func() error {
		domain := fmt.Sprintf("updopt-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})

		_, err = client.UpdateCertificateOptions(ctx, &acm.UpdateCertificateOptionsInput{
			CertificateArn: resp.CertificateArn,
			Options: &types.CertificateOptions{
				CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled,
			},
		})
		if err != nil {
			return err
		}
		desc, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if desc.Certificate.Options == nil {
			return fmt.Errorf("Options is nil after update")
		}
		if desc.Certificate.Options.CertificateTransparencyLoggingPreference != types.CertificateTransparencyLoggingPreferenceDisabled {
			return fmt.Errorf("expected DISABLED, got %s", desc.Certificate.Options.CertificateTransparencyLoggingPreference)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "UpdateCertificateOptions_NonExistent", func() error {
		_, err := client.UpdateCertificateOptions(ctx, &acm.UpdateCertificateOptionsInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
			Options: &types.CertificateOptions{
				CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceEnabled,
			},
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// --- RenewCertificate ---

	results = append(results, r.RunTest("acm", "RenewCertificate_ImportedCert_Error", func() error {
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  testKeyPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: importResp.CertificateArn})
		_, err = client.RenewCertificate(ctx, &acm.RenewCertificateInput{
			CertificateArn: importResp.CertificateArn,
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RenewCertificate_NonExistent", func() error {
		_, err := client.RenewCertificate(ctx, &acm.RenewCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// --- RevokeCertificate ---

	results = append(results, r.RunTest("acm", "RevokeCertificate_ImportedCert", func() error {
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  testKeyPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: importResp.CertificateArn})
		_, err = client.RevokeCertificate(ctx, &acm.RevokeCertificateInput{
			CertificateArn:   importResp.CertificateArn,
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return err
		}
		desc, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: importResp.CertificateArn})
		if err != nil {
			return err
		}
		if desc.Certificate.Status != types.CertificateStatusRevoked {
			return fmt.Errorf("expected REVOKED, got %s", desc.Certificate.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_VerifyRevokedAt", func() error {
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  testKeyPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: importResp.CertificateArn})
		_, err = client.RevokeCertificate(ctx, &acm.RevokeCertificateInput{
			CertificateArn:   importResp.CertificateArn,
			RevocationReason: types.RevocationReasonSuperseded,
		})
		if err != nil {
			return err
		}
		desc, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: importResp.CertificateArn})
		if err != nil {
			return err
		}
		if desc.Certificate.RevokedAt == nil {
			return fmt.Errorf("RevokedAt is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_VerifyRevocationReason", func() error {
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  testKeyPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: importResp.CertificateArn})
		_, err = client.RevokeCertificate(ctx, &acm.RevokeCertificateInput{
			CertificateArn:   importResp.CertificateArn,
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return err
		}
		desc, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: importResp.CertificateArn})
		if err != nil {
			return err
		}
		if desc.Certificate.RevocationReason != types.RevocationReasonKeyCompromise {
			return fmt.Errorf("expected KEY_COMPROMISE, got %s", desc.Certificate.RevocationReason)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_AlreadyRevoked", func() error {
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  testKeyPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: importResp.CertificateArn})
		_, err = client.RevokeCertificate(ctx, &acm.RevokeCertificateInput{
			CertificateArn:   importResp.CertificateArn,
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return fmt.Errorf("first revoke failed: %v", err)
		}
		_, err = client.RevokeCertificate(ctx, &acm.RevokeCertificateInput{
			CertificateArn:   importResp.CertificateArn,
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err == nil {
			return fmt.Errorf("expected error for already revoked cert")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_PendingValidation", func() error {
		domain := fmt.Sprintf("revoke-pv-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})
		_, err = client.RevokeCertificate(ctx, &acm.RevokeCertificateInput{
			CertificateArn:   resp.CertificateArn,
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return fmt.Errorf("RevokeCertificate failed: %v", err)
		}
		desc, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
			CertificateArn: resp.CertificateArn,
		})
		if err != nil {
			return err
		}
		if desc.Certificate.Status != types.CertificateStatusRevoked {
			return fmt.Errorf("expected REVOKED, got %s", desc.Certificate.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_NonExistent", func() error {
		_, err := client.RevokeCertificate(ctx, &acm.RevokeCertificateInput{
			CertificateArn:   aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// --- ImportCertificate ---

	results = append(results, r.RunTest("acm", "ImportCertificate", func() error {
		resp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  testKeyPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})
		if resp.CertificateArn == nil {
			return fmt.Errorf("CertificateArn is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ImportCertificate_WithCertificateChain", func() error {
		resp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate:      testCertPEM,
			PrivateKey:       testKeyPEM,
			CertificateChain: testChainPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})
		getResp, err := client.GetCertificate(ctx, &acm.GetCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if aws.ToString(getResp.CertificateChain) == "" {
			return fmt.Errorf("CertificateChain is empty")
		}
		return nil
	}))

	// --- ExportCertificate ---

	results = append(results, r.RunTest("acm", "ExportCertificate", func() error {
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  testKeyPEM,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: importResp.CertificateArn})
		exportResp, err := client.ExportCertificate(ctx, &acm.ExportCertificateInput{
			CertificateArn: importResp.CertificateArn,
			Passphrase:     []byte("test-passphrase"),
		})
		if err != nil {
			return err
		}
		if aws.ToString(exportResp.Certificate) == "" {
			return fmt.Errorf("Certificate is empty")
		}
		if exportResp.PrivateKey == nil || aws.ToString(exportResp.PrivateKey) == "" {
			return fmt.Errorf("PrivateKey is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ExportCertificate_AMAZON_ISSUED_Error", func() error {
		domain := fmt.Sprintf("export-ai-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: resp.CertificateArn})
		_, err = client.ExportCertificate(ctx, &acm.ExportCertificateInput{
			CertificateArn: resp.CertificateArn,
			Passphrase:     []byte("test-passphrase"),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	// --- AccountConfiguration ---

	results = append(results, r.RunTest("acm", "GetAccountConfiguration_DefaultValues", func() error {
		_, _ = client.PutAccountConfiguration(ctx, &acm.PutAccountConfigurationInput{
			IdempotencyToken: aws.String(fmt.Sprintf("reset-%d", time.Now().UnixNano())),
			ExpiryEvents:     &types.ExpiryEventsConfiguration{DaysBeforeExpiry: aws.Int32(45)},
		})
		resp, err := client.GetAccountConfiguration(ctx, &acm.GetAccountConfigurationInput{})
		if err != nil {
			return err
		}
		if resp.ExpiryEvents == nil {
			return fmt.Errorf("ExpiryEvents is nil")
		}
		if resp.ExpiryEvents.DaysBeforeExpiry == nil {
			return fmt.Errorf("DaysBeforeExpiry is nil")
		}
		if aws.ToInt32(resp.ExpiryEvents.DaysBeforeExpiry) != 45 {
			return fmt.Errorf("expected default 45, got %d", aws.ToInt32(resp.ExpiryEvents.DaysBeforeExpiry))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "GetAccountConfiguration_RoundTrip", func() error {
		_, err := client.PutAccountConfiguration(ctx, &acm.PutAccountConfigurationInput{
			IdempotencyToken: aws.String(fmt.Sprintf("rt-%d", time.Now().UnixNano())),
			ExpiryEvents: &types.ExpiryEventsConfiguration{
				DaysBeforeExpiry: aws.Int32(30),
			},
		})
		if err != nil {
			return err
		}
		resp, err := client.GetAccountConfiguration(ctx, &acm.GetAccountConfigurationInput{})
		if err != nil {
			return err
		}
		if resp.ExpiryEvents == nil {
			return fmt.Errorf("ExpiryEvents is nil")
		}
		if aws.ToInt32(resp.ExpiryEvents.DaysBeforeExpiry) != 30 {
			return fmt.Errorf("expected 30, got %d", aws.ToInt32(resp.ExpiryEvents.DaysBeforeExpiry))
		}
		return nil
	}))

	return results
}
