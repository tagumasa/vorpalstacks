package testutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"vorpalstacks-sdk-tests/config"
)

// generateTestRSAKey creates a fresh RSA key for credential material tests.
func generateTestRSAKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("rsa key generation failed: %v", err))
	}
	return key
}

// generateSelfSignedCertificate creates a self-signed X.509 certificate and
// returns its PEM body plus the PEM private key.
func generateSelfSignedCertificate(commonName string) (string, string, error) {
	key := generateTestRSAKey()
	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:   time.Now().Add(-time.Hour),
		NotAfter:    time.Now().Add(24 * 365 * time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return "", "", err
	}
	certPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
	keyDER := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDER}))
	return certPEM, keyPEM, nil
}

// newIAMClientForUser builds an IAM client authenticated as the given user
// via a freshly created access key. The returned key must be deleted by the
// caller.
func newIAMClientForUser(r *TestRunner, tc *iamTestContext, userName string) (*iam.Client, *types.AccessKey, error) {
	key, err := tc.client.CreateAccessKey(tc.ctx, &iam.CreateAccessKeyInput{UserName: aws.String(userName)})
	if err != nil {
		return nil, nil, err
	}
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		_, _ = tc.client.DeleteAccessKey(tc.ctx, &iam.DeleteAccessKeyInput{
			UserName:    aws.String(userName),
			AccessKeyId: key.AccessKey.AccessKeyId,
		})
		return nil, nil, err
	}
	cfg.Credentials = credentials.NewStaticCredentialsProvider(
		*key.AccessKey.AccessKeyId,
		*key.AccessKey.SecretAccessKey,
		"",
	)
	return iam.NewFromConfig(cfg), key.AccessKey, nil
}

func (r *TestRunner) iamSigningCertificateTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	newCertUser := func(suffix string) (string, *iam.Client, *types.AccessKey, func()) {
		user := fmt.Sprintf("SignCert-%s-%s", suffix, tc.ts)
		if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(user)}); err != nil {
			return "", nil, nil, func() {}
		}
		client, key, err := newIAMClientForUser(r, tc, user)
		if err != nil {
			return "", nil, nil, func() {}
		}
		cleanup := func() {
			_, _ = tc.client.DeleteAccessKey(tc.ctx, &iam.DeleteAccessKeyInput{
				UserName:    aws.String(user),
				AccessKeyId: key.AccessKeyId,
			})
			_, _ = tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(user)})
		}
		return user, client, key, cleanup
	}

	results = append(results, r.RunTest("iam", "UploadSigningCertificate", func() error {
		user, client, key, cleanup := newCertUser("upload")
		defer cleanup()
		if client == nil {
			return fmt.Errorf("user client setup failed")
		}
		_ = key

		certPEM, _, err := generateSelfSignedCertificate(user)
		if err != nil {
			return err
		}

		// UserName is optional on the wire; omitting it addresses the
		// certificate to the authenticated caller.
		resp, err := client.UploadSigningCertificate(tc.ctx, &iam.UploadSigningCertificateInput{
			CertificateBody: aws.String(certPEM),
		})
		if err != nil {
			return err
		}
		if resp.Certificate == nil {
			return fmt.Errorf("certificate is nil")
		}
		if aws.ToString(resp.Certificate.UserName) != user {
			return fmt.Errorf("certificate user: got %s, want %s", aws.ToString(resp.Certificate.UserName), user)
		}
		if string(resp.Certificate.Status) != "Active" {
			return fmt.Errorf("certificate status: got %s, want Active", string(resp.Certificate.Status))
		}
		if aws.ToString(resp.Certificate.CertificateId) == "" {
			return fmt.Errorf("certificate id is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UploadSigningCertificate_Malformed", func() error {
		_, client, _, cleanup := newCertUser("malformed")
		defer cleanup()
		if client == nil {
			return fmt.Errorf("user client setup failed")
		}

		_, err := client.UploadSigningCertificate(tc.ctx, &iam.UploadSigningCertificateInput{
			CertificateBody: aws.String("this is not a certificate"),
		})
		if err == nil {
			return fmt.Errorf("malformed certificate body must be rejected")
		}
		if !containsErrorCode(err, "MalformedCertificate") {
			return fmt.Errorf("malformed certificate: got %v, want MalformedCertificate", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UploadSigningCertificate_Duplicate", func() error {
		user, client, _, cleanup := newCertUser("duplicate")
		defer cleanup()
		if client == nil {
			return fmt.Errorf("user client setup failed")
		}

		certPEM, _, err := generateSelfSignedCertificate(user)
		if err != nil {
			return err
		}
		if _, err := client.UploadSigningCertificate(tc.ctx, &iam.UploadSigningCertificateInput{
			CertificateBody: aws.String(certPEM),
		}); err != nil {
			return err
		}

		_, err = client.UploadSigningCertificate(tc.ctx, &iam.UploadSigningCertificateInput{
			CertificateBody: aws.String(certPEM),
		})
		if err == nil {
			return fmt.Errorf("uploading the same certificate twice must be rejected")
		}
		if !containsErrorCode(err, "DuplicateCertificate") {
			return fmt.Errorf("duplicate certificate: got %v, want DuplicateCertificate", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UploadSigningCertificate_Quota", func() error {
		user, client, _, cleanup := newCertUser("quota")
		defer cleanup()
		if client == nil {
			return fmt.Errorf("user client setup failed")
		}

		for i := 0; i < 2; i++ {
			certPEM, _, err := generateSelfSignedCertificate(fmt.Sprintf("%s-%d", user, i))
			if err != nil {
				return err
			}
			if _, err := client.UploadSigningCertificate(tc.ctx, &iam.UploadSigningCertificateInput{
				CertificateBody: aws.String(certPEM),
			}); err != nil {
				return fmt.Errorf("upload %d within the quota failed: %w", i+1, err)
			}
		}

		certPEM, _, err := generateSelfSignedCertificate(user + "-third")
		if err != nil {
			return err
		}
		_, err = client.UploadSigningCertificate(tc.ctx, &iam.UploadSigningCertificateInput{
			CertificateBody: aws.String(certPEM),
		})
		if err == nil {
			return fmt.Errorf("the third certificate must exceed the per-user quota")
		}
		if !containsErrorCode(err, "LimitExceeded") {
			return fmt.Errorf("quota: got %v, want LimitExceeded", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListSigningCertificates", func() error {
		user, client, _, cleanup := newCertUser("list")
		defer cleanup()
		if client == nil {
			return fmt.Errorf("user client setup failed")
		}

		certPEM, _, err := generateSelfSignedCertificate(user)
		if err != nil {
			return err
		}
		if _, err := client.UploadSigningCertificate(tc.ctx, &iam.UploadSigningCertificateInput{
			CertificateBody: aws.String(certPEM),
		}); err != nil {
			return err
		}

		// Omitting UserName lists the caller's own certificates.
		resp, err := client.ListSigningCertificates(tc.ctx, &iam.ListSigningCertificatesInput{})
		if err != nil {
			return err
		}
		if len(resp.Certificates) != 1 {
			return fmt.Errorf("certificate count: got %d, want 1", len(resp.Certificates))
		}
		if aws.ToString(resp.Certificates[0].UserName) != user {
			return fmt.Errorf("certificate owner: got %s, want %s", aws.ToString(resp.Certificates[0].UserName), user)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateSigningCertificate", func() error {
		user, client, _, cleanup := newCertUser("update")
		defer cleanup()
		if client == nil {
			return fmt.Errorf("user client setup failed")
		}

		certPEM, _, err := generateSelfSignedCertificate(user)
		if err != nil {
			return err
		}
		upload, err := client.UploadSigningCertificate(tc.ctx, &iam.UploadSigningCertificateInput{
			CertificateBody: aws.String(certPEM),
		})
		if err != nil {
			return err
		}
		certID := upload.Certificate.CertificateId

		if _, err := tc.client.UpdateSigningCertificate(tc.ctx, &iam.UpdateSigningCertificateInput{
			UserName:      aws.String(user),
			CertificateId: certID,
			Status:        types.StatusTypeInactive,
		}); err != nil {
			return err
		}

		list, err := tc.client.ListSigningCertificates(tc.ctx, &iam.ListSigningCertificatesInput{
			UserName: aws.String(user),
		})
		if err != nil {
			return err
		}
		foundInactive := false
		for _, c := range list.Certificates {
			if aws.ToString(c.CertificateId) == aws.ToString(certID) && string(c.Status) == "Inactive" {
				foundInactive = true
			}
		}
		if !foundInactive {
			return fmt.Errorf("certificate %s not Inactive after update", aws.ToString(certID))
		}

		// Naming a UserName that does not own the certificate must fail
		// with NoSuchEntity instead of mutating the certificate.
		other := fmt.Sprintf("SignCert-other-%s", tc.ts)
		cleanupOther, err := tc.createUser(other)
		if err != nil {
			return err
		}
		defer cleanupOther()

		_, err = tc.client.UpdateSigningCertificate(tc.ctx, &iam.UpdateSigningCertificateInput{
			UserName:      aws.String(other),
			CertificateId: certID,
			Status:        types.StatusTypeActive,
		})
		if err == nil {
			return fmt.Errorf("updating through a non-owner UserName must be rejected")
		}
		if !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("non-owner update: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteSigningCertificate", func() error {
		user, client, _, cleanup := newCertUser("delete")
		defer cleanup()
		if client == nil {
			return fmt.Errorf("user client setup failed")
		}

		certPEM, _, err := generateSelfSignedCertificate(user)
		if err != nil {
			return err
		}
		upload, err := client.UploadSigningCertificate(tc.ctx, &iam.UploadSigningCertificateInput{
			CertificateBody: aws.String(certPEM),
		})
		if err != nil {
			return err
		}
		certID := upload.Certificate.CertificateId

		// Deleting through a UserName that does not own the certificate
		// must fail with NoSuchEntity.
		other := fmt.Sprintf("SignCert-delother-%s", tc.ts)
		cleanupDelOther, err := tc.createUser(other)
		if err != nil {
			return err
		}
		defer cleanupDelOther()

		_, err = tc.client.DeleteSigningCertificate(tc.ctx, &iam.DeleteSigningCertificateInput{
			UserName:      aws.String(other),
			CertificateId: certID,
		})
		if err == nil {
			return fmt.Errorf("deleting through a non-owner UserName must be rejected")
		}
		if !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("non-owner delete: got %v, want NoSuchEntity", err)
		}

		if _, err := tc.client.DeleteSigningCertificate(tc.ctx, &iam.DeleteSigningCertificateInput{
			UserName:      aws.String(user),
			CertificateId: certID,
		}); err != nil {
			return err
		}

		list, err := tc.client.ListSigningCertificates(tc.ctx, &iam.ListSigningCertificatesInput{
			UserName: aws.String(user),
		})
		if err != nil {
			return err
		}
		if len(list.Certificates) != 0 {
			return fmt.Errorf("certificate still listed after delete: %d left", len(list.Certificates))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "SigningCertificate_ImplicitCallerOwnership", func() error {
		owner, ownerClient, _, cleanup := newCertUser("implicit")
		defer cleanup()
		if ownerClient == nil {
			return fmt.Errorf("owner client setup failed")
		}
		_, otherClient, _, otherCleanup := newCertUser("implicit-other")
		defer otherCleanup()
		if otherClient == nil {
			return fmt.Errorf("other client setup failed")
		}

		certPEM, _, err := generateSelfSignedCertificate(owner)
		if err != nil {
			return err
		}
		upload, err := ownerClient.UploadSigningCertificate(tc.ctx, &iam.UploadSigningCertificateInput{
			CertificateBody: aws.String(certPEM),
		})
		if err != nil {
			return err
		}
		certID := upload.Certificate.CertificateId

		// The owner may act on their own certificate without naming a user.
		if _, err := ownerClient.UpdateSigningCertificate(tc.ctx, &iam.UpdateSigningCertificateInput{
			CertificateId: certID,
			Status:        types.StatusTypeInactive,
		}); err != nil {
			return fmt.Errorf("implicit-caller update by the owner: %w", err)
		}

		// An omitted UserName resolves to the caller; another user must
		// not be able to update or delete a certificate they do not own.
		if _, err := otherClient.UpdateSigningCertificate(tc.ctx, &iam.UpdateSigningCertificateInput{
			CertificateId: certID,
			Status:        types.StatusTypeActive,
		}); err == nil {
			return fmt.Errorf("implicit-caller update by a non-owner must be rejected")
		} else if !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("non-owner implicit update: got %v, want NoSuchEntity", err)
		}

		if _, err := otherClient.DeleteSigningCertificate(tc.ctx, &iam.DeleteSigningCertificateInput{
			CertificateId: certID,
		}); err == nil {
			return fmt.Errorf("implicit-caller delete by a non-owner must be rejected")
		} else if !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("non-owner implicit delete: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	return results
}
