package cloudwatchlogs

import (
	"context"
	"fmt"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/crypto"
)

// lookupTableEncryptionContextKey is the encryption context key carrying the
// lookup table ARN, mirroring the aws:s3:arn context used by S3 SSE-KMS.
const lookupTableEncryptionContextKey = "aws:cloudwatchlogs:arn"

// validateLookupTableKmsKey checks the documented length ceiling and, when a
// KMS invoker is available, that the key exists.
func (s *LogsService) validateLookupTableKmsKey(kmsKeyID string) error {
	if kmsKeyID == "" {
		return nil
	}
	if len(kmsKeyID) > logsstore.MaxKmsKeyIdLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("kmsKeyId exceeds %d characters", logsstore.MaxKmsKeyIdLength), 400)
	}
	if s.kms == nil {
		return nil
	}
	if !s.kms.KeyExists(context.Background(), kmsKeyID) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("KMS key %s not found", kmsKeyID), 400)
	}
	return nil
}

// encryptLookupTableBody envelope-encrypts the CSV body with the customer
// managed KMS key: a data key from KMS protects the body with AES-GCM and the
// wrapped data key is stored alongside the ciphertext.
func (s *LogsService) encryptLookupTableBody(body []byte, kmsKeyID, tableArn string) (encrypted, dataKey, nonce []byte, err error) {
	if s.kms == nil {
		return nil, nil, nil, NewLogsError("InvalidParameterException",
			"KMS encryption is not available", 400)
	}
	encryptionContext := map[string]string{lookupTableEncryptionContextKey: tableArn}
	dk, err := s.kms.GenerateDataKey(context.Background(), kmsKeyID, "AES_256", encryptionContext, tableArn)
	if err != nil {
		return nil, nil, nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("failed to generate data key for lookup table: %v", err), 400)
	}
	nonce, err = crypto.RandomNonce()
	if err != nil {
		return nil, nil, nil, err
	}
	encrypted, err = crypto.AESGCMEncryptWithNonce(dk.Plaintext, body, nonce)
	if err != nil {
		return nil, nil, nil, err
	}
	return encrypted, dk.CiphertextBlob, nonce, nil
}

// lookupTablePlainBody returns the CSV content of a stored lookup table,
// decrypting the envelope-encrypted body when a customer-managed KMS key is
// configured.
func (s *LogsService) lookupTablePlainBody(lt *logsstore.LookupTable, region string) (string, error) {
	if len(lt.EncryptedBody) == 0 {
		return lt.TableBody, nil
	}
	if s.kms == nil {
		return "", NewLogsError("InvalidParameterException",
			"KMS encryption is not available", 400)
	}
	tableArn := lookupTableArn(region, s.accountID, lt.Name)
	encryptionContext := map[string]string{lookupTableEncryptionContextKey: tableArn}
	plaintextKey, err := s.kms.Decrypt(context.Background(), lt.KmsKeyId, lt.EncryptedDataKey, encryptionContext, tableArn)
	if err != nil {
		return "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("failed to decrypt lookup table %s: %v", lt.Name, err), 400)
	}
	plaintext, err := crypto.AESGCMDecryptWithNonce(plaintextKey, lt.EncryptedBody, lt.ContentNonce)
	if err != nil {
		return "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("failed to decrypt lookup table %s: %v", lt.Name, err), 400)
	}
	return string(plaintext), nil
}
