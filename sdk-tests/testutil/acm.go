package testutil

import (
	"context"
	"fmt"
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

	// Test 1: List Certificates
	results = append(results, r.RunTest("acm", "ListCertificates", func() error {
		_, err := client.ListCertificates(ctx, &acm.ListCertificatesInput{
			MaxItems: aws.Int32(10),
		})
		return err
	}))

	// Test 2: Describe Certificate - Create a cert first, then describe
	var testCertARN string
	results = append(results, r.RunTest("acm", "DescribeCertificate", func() error {
		testDomain := fmt.Sprintf("describe-test-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(testDomain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		testCertARN = aws.ToString(resp.CertificateArn)
		_, err = client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
			CertificateArn: resp.CertificateArn,
		})
		return err
	}))

	// Test 3: Get Certificate
	results = append(results, r.RunTest("acm", "GetCertificate", func() error {
		if testCertARN == "" {
			return fmt.Errorf("no certificate ARN available from DescribeCertificate test")
		}
		_, err := client.GetCertificate(ctx, &acm.GetCertificateInput{
			CertificateArn: aws.String(testCertARN),
		})
		return err
	}))

	// Test 4: Request Certificate
	domainName := fmt.Sprintf("example-%d.com", time.Now().UnixNano())
	var certARN string
	results = append(results, r.RunTest("acm", "RequestCertificate", func() error {
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domainName),
			ValidationMethod: types.ValidationMethodDns,
			IdempotencyToken: aws.String("test-token"),
		})
		if err == nil {
			certARN = aws.ToString(resp.CertificateArn)
		}
		return err
	}))

	// Test 5: Describe Certificate (with created cert)
	if certARN != "" {
		results = append(results, r.RunTest("acm", "DescribeCertificateCreated", func() error {
			_, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
				CertificateArn: aws.String(certARN),
			})
			return err
		}))

		// Test 6: Delete Certificate
		results = append(results, r.RunTest("acm", "DeleteCertificate", func() error {
			_, err := client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{
				CertificateArn: aws.String(certARN),
			})
			return err
		}))
	}

	// Cleanup testCertARN
	if testCertARN != "" {
		client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{
			CertificateArn: aws.String(testCertARN),
		})
	}

	// Test 7: Add Tags to Certificate
	results = append(results, r.RunTest("acm", "AddTagsToCertificate", func() error {
		domainName2 := fmt.Sprintf("example2-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domainName2),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		certARN2 := aws.ToString(resp.CertificateArn)
		_, err = client.AddTagsToCertificate(ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: aws.String(certARN2),
			Tags: []types.Tag{
				{
					Key:   aws.String("Environment"),
					Value: aws.String("test"),
				},
				{
					Key:   aws.String("Owner"),
					Value: aws.String("test-user"),
				},
			},
		})
		// Cleanup
		client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{
			CertificateArn: aws.String(certARN2),
		})
		return err
	}))

	// Test 8: List Tags for Certificate
	results = append(results, r.RunTest("acm", "ListTagsForCertificate", func() error {
		domainName3 := fmt.Sprintf("example3-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domainName3),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		certARN3 := aws.ToString(resp.CertificateArn)
		client.AddTagsToCertificate(ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: aws.String(certARN3),
			Tags: []types.Tag{
				{
					Key:   aws.String("Test"),
					Value: aws.String("value"),
				},
			},
		})
		_, err = client.ListTagsForCertificate(ctx, &acm.ListTagsForCertificateInput{
			CertificateArn: aws.String(certARN3),
		})
		// Cleanup
		client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{
			CertificateArn: aws.String(certARN3),
		})
		return err
	}))

	// Test 9: Remove Tags from Certificate
	results = append(results, r.RunTest("acm", "RemoveTagsFromCertificate", func() error {
		domainName4 := fmt.Sprintf("example4-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domainName4),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		certARN4 := aws.ToString(resp.CertificateArn)
		client.AddTagsToCertificate(ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: aws.String(certARN4),
			Tags: []types.Tag{
				{
					Key:   aws.String("Test"),
					Value: aws.String("value"),
				},
			},
		})
		_, err = client.RemoveTagsFromCertificate(ctx, &acm.RemoveTagsFromCertificateInput{
			CertificateArn: aws.String(certARN4),
			Tags: []types.Tag{
				{
					Key: aws.String("Test"),
				},
			},
		})
		// Cleanup
		client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{
			CertificateArn: aws.String(certARN4),
		})
		return err
	}))

	// Test 10: Resend Validation Email
	results = append(results, r.RunTest("acm", "ResendValidationEmail", func() error {
		domainName5 := fmt.Sprintf("example5-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domainName5),
			ValidationMethod: types.ValidationMethodEmail,
		})
		if err != nil {
			return err
		}
		certARN5 := aws.ToString(resp.CertificateArn)
		_, err = client.ResendValidationEmail(ctx, &acm.ResendValidationEmailInput{
			CertificateArn:   aws.String(certARN5),
			Domain:           aws.String(domainName5),
			ValidationDomain: aws.String(domainName5),
		})
		// Cleanup
		client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{
			CertificateArn: aws.String(certARN5),
		})
		return err
	}))

	// Test 11: Update Certificate Options
	results = append(results, r.RunTest("acm", "UpdateCertificateOptions", func() error {
		domainName6 := fmt.Sprintf("example6-%d.com", time.Now().UnixNano())
		resp, err := client.RequestCertificate(ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domainName6),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		certARN6 := aws.ToString(resp.CertificateArn)
		_, err = client.UpdateCertificateOptions(ctx, &acm.UpdateCertificateOptionsInput{
			CertificateArn: aws.String(certARN6),
			Options: &types.CertificateOptions{
				CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceEnabled,
			},
		})
		// Cleanup
		client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{
			CertificateArn: aws.String(certARN6),
		})
		return err
	}))

	// Test 12: Import Certificate
	results = append(results, r.RunTest("acm", "ImportCertificate", func() error {
		// Use a dummy certificate for testing (will fail validation but API should work)
		certPem := []byte(`-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHHCgVZU1JUMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3RjYTAeFw0yNDAxMDEwMDAwMDBaFw0yNTAxMDEwMDAwMDBaMBExDzANBgNVBAMM
BnRlc3RjYTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAwK0j6f8C6hJ7u8P
-----END CERTIFICATE-----`)
		privateKeyPem := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBOQIBAAJBAKjHCBmV1SlQwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMCr
-----END RSA PRIVATE KEY-----`)
		_, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: certPem,
			PrivateKey:  privateKeyPem,
		})
		return err
	}))

	// Test 13: Get Account Configuration
	results = append(results, r.RunTest("acm", "GetAccountConfiguration", func() error {
		_, err := client.GetAccountConfiguration(ctx, &acm.GetAccountConfigurationInput{})
		return err
	}))

	// Test 14: Put Account Configuration
	results = append(results, r.RunTest("acm", "PutAccountConfiguration", func() error {
		_, err := client.PutAccountConfiguration(ctx, &acm.PutAccountConfigurationInput{
			IdempotencyToken: aws.String("test-token"),
			ExpiryEvents: &types.ExpiryEventsConfiguration{
				DaysBeforeExpiry: aws.Int32(30),
			},
		})
		return err
	}))

	// Test 15: Export Certificate (import a cert, then export it)
	results = append(results, r.RunTest("acm", "ExportCertificate", func() error {
		certPem := []byte(`-----BEGIN CERTIFICATE-----
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
		privateKeyPem := []byte(`-----BEGIN RSA PRIVATE KEY-----
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
		importResp, err := client.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate: certPem,
			PrivateKey:  privateKeyPem,
		})
		if err != nil {
			return err
		}
		defer client.DeleteCertificate(ctx, &acm.DeleteCertificateInput{
			CertificateArn: importResp.CertificateArn,
		})
		_, err = client.ExportCertificate(ctx, &acm.ExportCertificateInput{
			CertificateArn: importResp.CertificateArn,
			Passphrase:     []byte("test-passphrase"),
		})
		return err
	}))

	return results
}
