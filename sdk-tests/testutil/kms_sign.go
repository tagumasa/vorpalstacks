package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (r *TestRunner) runKMSSignTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kms", "GetPublicKey_RSA", func() error {
		if tc.rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		resp, err := tc.client.GetPublicKey(tc.ctx, &kms.GetPublicKeyInput{
			KeyId: aws.String(tc.rsaKeyID),
		})
		if err != nil {
			return err
		}
		if len(resp.PublicKey) == 0 {
			return fmt.Errorf("public key is nil or empty")
		}
		if resp.KeySpec != types.KeySpecRsa2048 {
			return fmt.Errorf("expected KeySpec=RSA_2048, got %s", resp.KeySpec)
		}
		if resp.KeyUsage != types.KeyUsageTypeSignVerify {
			return fmt.Errorf("expected KeyUsage=SIGN_VERIFY, got %s", resp.KeyUsage)
		}
		if len(resp.SigningAlgorithms) == 0 {
			return fmt.Errorf("SigningAlgorithms is nil or empty")
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Sign_RSA", func() error {
		if tc.rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		message := []byte("message to sign")
		resp, err := tc.client.Sign(tc.ctx, &kms.SignInput{
			KeyId:            aws.String(tc.rsaKeyID),
			Message:          message,
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: types.SigningAlgorithmSpecRsassaPkcs1V15Sha256,
		})
		if err != nil {
			return err
		}
		if len(resp.Signature) == 0 {
			return fmt.Errorf("signature is nil or empty")
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		if resp.SigningAlgorithm != types.SigningAlgorithmSpecRsassaPkcs1V15Sha256 {
			return fmt.Errorf("expected SigningAlgorithm=RSASSA_PKCS1_V1_5_SHA_256, got %s", resp.SigningAlgorithm)
		}
		tc.signature = resp.Signature
		return nil
	}))

	results = append(results, r.RunTest("kms", "Verify_RSA", func() error {
		if tc.rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		if tc.signature == nil {
			return fmt.Errorf("signature not available")
		}
		message := []byte("message to sign")
		resp, err := tc.client.Verify(tc.ctx, &kms.VerifyInput{
			KeyId:            aws.String(tc.rsaKeyID),
			Message:          message,
			Signature:        tc.signature,
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: types.SigningAlgorithmSpecRsassaPkcs1V15Sha256,
		})
		if err != nil {
			return err
		}
		if !resp.SignatureValid {
			return fmt.Errorf("expected signature to be valid")
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Verify_RSA_InvalidSignature", func() error {
		if tc.rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		message := []byte("different message")
		badSig := make([]byte, 256)
		for i := range badSig {
			badSig[i] = 0xFF
		}
		resp, err := tc.client.Verify(tc.ctx, &kms.VerifyInput{
			KeyId:            aws.String(tc.rsaKeyID),
			Message:          message,
			Signature:        badSig,
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: types.SigningAlgorithmSpecRsassaPkcs1V15Sha256,
		})
		if err != nil {
			return err
		}
		if resp.SignatureValid {
			return fmt.Errorf("expected invalid signature")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateMac", func() error {
		if tc.hmacKeyID == "" {
			return fmt.Errorf("HMAC key ID not available")
		}
		message := []byte("message to mac")
		resp, err := tc.client.GenerateMac(tc.ctx, &kms.GenerateMacInput{
			KeyId:        aws.String(tc.hmacKeyID),
			Message:      message,
			MacAlgorithm: types.MacAlgorithmSpecHmacSha256,
		})
		if err != nil {
			return err
		}
		if len(resp.Mac) == 0 {
			return fmt.Errorf("MAC is nil or empty")
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		if resp.MacAlgorithm != types.MacAlgorithmSpecHmacSha256 {
			return fmt.Errorf("expected MacAlgorithm=HMAC_SHA_256, got %s", resp.MacAlgorithm)
		}
		tc.macValue = resp.Mac
		return nil
	}))

	results = append(results, r.RunTest("kms", "VerifyMac", func() error {
		if tc.hmacKeyID == "" {
			return fmt.Errorf("HMAC key ID not available")
		}
		if tc.macValue == nil {
			return fmt.Errorf("MAC value not available")
		}
		message := []byte("message to mac")
		resp, err := tc.client.VerifyMac(tc.ctx, &kms.VerifyMacInput{
			KeyId:        aws.String(tc.hmacKeyID),
			Message:      message,
			Mac:          tc.macValue,
			MacAlgorithm: types.MacAlgorithmSpecHmacSha256,
		})
		if err != nil {
			return err
		}
		if !resp.MacValid {
			return fmt.Errorf("expected MAC to be valid")
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "VerifyMac_InvalidMac", func() error {
		if tc.hmacKeyID == "" {
			return fmt.Errorf("HMAC key ID not available")
		}
		message := []byte("message to mac")
		badMac := make([]byte, 32)
		for i := range badMac {
			badMac[i] = 0xFF
		}
		resp, err := tc.client.VerifyMac(tc.ctx, &kms.VerifyMacInput{
			KeyId:        aws.String(tc.hmacKeyID),
			Message:      message,
			Mac:          badMac,
			MacAlgorithm: types.MacAlgorithmSpecHmacSha256,
		})
		if err != nil {
			return err
		}
		if resp.MacValid {
			return fmt.Errorf("expected invalid MAC")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKeyPair", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.keyID)})
		resp, err := tc.client.GenerateDataKeyPair(tc.ctx, &kms.GenerateDataKeyPairInput{
			KeyId:       aws.String(tc.keyID),
			KeyPairSpec: types.DataKeyPairSpecRsa2048,
		})
		if err != nil {
			return err
		}
		if len(resp.PrivateKeyCiphertextBlob) == 0 {
			return fmt.Errorf("private key ciphertext is nil or empty")
		}
		if len(resp.PrivateKeyPlaintext) == 0 {
			return fmt.Errorf("private key plaintext is nil or empty")
		}
		if len(resp.PublicKey) == 0 {
			return fmt.Errorf("public key is nil or empty")
		}
		if resp.KeyPairSpec != types.DataKeyPairSpecRsa2048 {
			return fmt.Errorf("expected KeyPairSpec=RSA_2048, got %s", resp.KeyPairSpec)
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKeyPairWithoutPlaintext", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.keyID)})
		resp, err := tc.client.GenerateDataKeyPairWithoutPlaintext(tc.ctx, &kms.GenerateDataKeyPairWithoutPlaintextInput{
			KeyId:       aws.String(tc.keyID),
			KeyPairSpec: types.DataKeyPairSpecRsa2048,
		})
		if err != nil {
			return err
		}
		if len(resp.PrivateKeyCiphertextBlob) == 0 {
			return fmt.Errorf("private key ciphertext is nil or empty")
		}
		if len(resp.PublicKey) == 0 {
			return fmt.Errorf("public key is nil or empty")
		}
		return nil
	}))

	return results
}
