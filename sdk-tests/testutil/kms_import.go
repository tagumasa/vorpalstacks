package testutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (r *TestRunner) runKMSImportTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	var extKeyID string
	var importToken []byte
	var publicKeyBytes []byte

	results = append(results, r.RunTest("kms", "GetParametersForImport", func() error {
		resp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
			Description: aws.String("Import origin key"),
			Origin:      types.OriginTypeExternal,
		})
		if err != nil {
			return fmt.Errorf("create external key: %v", err)
		}
		extKeyID = *resp.KeyMetadata.KeyId
		tc.addCleanupKey(extKeyID)

		if resp.KeyMetadata.KeyState != types.KeyStatePendingImport {
			return fmt.Errorf("expected KeyState=PendingImport for EXTERNAL origin, got %s", resp.KeyMetadata.KeyState)
		}
		if resp.KeyMetadata.Enabled {
			return fmt.Errorf("expected Enabled=false for EXTERNAL origin")
		}

		paramsResp, err := tc.client.GetParametersForImport(tc.ctx, &kms.GetParametersForImportInput{
			KeyId:             aws.String(extKeyID),
			WrappingAlgorithm: types.AlgorithmSpecRsaesOaepSha256,
			WrappingKeySpec:   types.WrappingKeySpecRsa2048,
		})
		if err != nil {
			return fmt.Errorf("get parameters: %v", err)
		}
		if len(paramsResp.ImportToken) == 0 {
			return fmt.Errorf("ImportToken is nil or empty")
		}
		if len(paramsResp.PublicKey) == 0 {
			return fmt.Errorf("PublicKey is nil or empty")
		}
		if paramsResp.KeyId == nil || *paramsResp.KeyId == "" {
			return fmt.Errorf("KeyId is nil or empty")
		}
		if paramsResp.ParametersValidTo == nil {
			return fmt.Errorf("ParametersValidTo is nil")
		}

		pubKey, err := x509.ParsePKIXPublicKey(paramsResp.PublicKey)
		if err != nil {
			return fmt.Errorf("public key is not valid DER PKIX: %v", err)
		}
		if _, ok := pubKey.(*rsa.PublicKey); !ok {
			return fmt.Errorf("public key is not RSA")
		}

		importToken = paramsResp.ImportToken
		publicKeyBytes = paramsResp.PublicKey
		return nil
	}))

	results = append(results, r.RunTest("kms", "ImportKeyMaterial", func() error {
		if extKeyID == "" || len(importToken) == 0 || len(publicKeyBytes) == 0 {
			return fmt.Errorf("import prerequisites not available")
		}

		pubKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
		if err != nil {
			return fmt.Errorf("parse public key: %v", err)
		}
		rsaPubKey, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("public key is not RSA")
		}

		keyMaterial := make([]byte, 32)
		if _, err := rand.Read(keyMaterial); err != nil {
			return fmt.Errorf("generate key material: %v", err)
		}

		encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, keyMaterial, nil)
		if err != nil {
			return fmt.Errorf("encrypt key material: %v", err)
		}

		_, err = tc.client.ImportKeyMaterial(tc.ctx, &kms.ImportKeyMaterialInput{
			KeyId:                aws.String(extKeyID),
			ImportToken:          importToken,
			EncryptedKeyMaterial: encryptedKey,
		})
		if err != nil {
			return fmt.Errorf("import key material: %v", err)
		}

		descResp, err := tc.client.DescribeKey(tc.ctx, &kms.DescribeKeyInput{KeyId: aws.String(extKeyID)})
		if err != nil {
			return fmt.Errorf("describe after import: %v", err)
		}
		if descResp.KeyMetadata.KeyState != types.KeyStateEnabled {
			return fmt.Errorf("expected KeyState=Enabled after import, got %s", descResp.KeyMetadata.KeyState)
		}
		if !descResp.KeyMetadata.Enabled {
			return fmt.Errorf("expected Enabled=true after import")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ImportKeyMaterial_EncryptDecryptRoundtrip", func() error {
		if extKeyID == "" {
			return fmt.Errorf("external key ID not available")
		}
		plaintext := []byte("test-imported-key-material-works")
		encResp, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(extKeyID),
			Plaintext: plaintext,
		})
		if err != nil {
			return fmt.Errorf("encrypt with imported key: %v", err)
		}
		if len(encResp.CiphertextBlob) == 0 {
			return fmt.Errorf("ciphertext is empty")
		}
		decResp, err := tc.client.Decrypt(tc.ctx, &kms.DecryptInput{
			CiphertextBlob: encResp.CiphertextBlob,
		})
		if err != nil {
			return fmt.Errorf("decrypt with imported key: %v", err)
		}
		if string(decResp.Plaintext) != string(plaintext) {
			return fmt.Errorf("plaintext mismatch: got %q, want %q", string(decResp.Plaintext), string(plaintext))
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ImportKeyMaterial_InvalidToken", func() error {
		if extKeyID == "" {
			return fmt.Errorf("external key ID not available")
		}
		badToken := make([]byte, 32)
		rand.Read(badToken)
		_, err := tc.client.ImportKeyMaterial(tc.ctx, &kms.ImportKeyMaterialInput{
			KeyId:                aws.String(extKeyID),
			ImportToken:          badToken,
			EncryptedKeyMaterial: []byte("garbage"),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid import token")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DeleteImportedKeyMaterial", func() error {
		if extKeyID == "" {
			return fmt.Errorf("external key ID not available")
		}
		_, err := tc.client.DeleteImportedKeyMaterial(tc.ctx, &kms.DeleteImportedKeyMaterialInput{
			KeyId: aws.String(extKeyID),
		})
		if err != nil {
			return err
		}

		descResp, err := tc.client.DescribeKey(tc.ctx, &kms.DescribeKeyInput{KeyId: aws.String(extKeyID)})
		if err != nil {
			return fmt.Errorf("describe after delete: %v", err)
		}
		if descResp.KeyMetadata.KeyState != types.KeyStatePendingImport {
			return fmt.Errorf("expected KeyState=PendingImport, got %s", descResp.KeyMetadata.KeyState)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DeleteImportedKeyMaterial_EncryptFails", func() error {
		if extKeyID == "" {
			return fmt.Errorf("external key ID not available")
		}
		_, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(extKeyID),
			Plaintext: []byte("should fail without key material"),
		})
		if err == nil {
			return fmt.Errorf("expected error when encrypting with key in PendingImport state")
		}
		return nil
	}))

	return results
}
