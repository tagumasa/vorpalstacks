package cloudwatchlogs

import (
	"context"
	"fmt"
	"net/http"
	"unicode/utf8"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/crypto"
)

// lookupTableEncryptionContextKey is the encryption context key carrying the
// lookup table ARN, mirroring the aws:s3:arn context used by S3 SSE-KMS.
const lookupTableEncryptionContextKey = "aws:cloudwatchlogs:arn"

// validateLookupTableKmsKey checks the member's length ceiling and, when
// a KMS invoker is available, that the key exists. The lookup-table
// family resolves bare key ids and aliases through the invoker (the ARN
// form is not its gate), so only the shared length bound applies here;
// the ARN form is the group-surface validator's contract.
func (s *LogsService) validateLookupTableKmsKey(kmsKeyID string) error {
	if kmsKeyID == "" {
		return nil
	}
	if utf8.RuneCountInString(kmsKeyID) > logsstore.MaxKmsKeyIdLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("kmsKeyId exceeds %d characters", logsstore.MaxKmsKeyIdLength), 400)
	}
	if s.kmsInvoker() == nil {
		return nil
	}
	if !s.kmsInvoker().KeyExists(context.Background(), kmsKeyID) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("KMS key %s not found", kmsKeyID), 400)
	}
	return nil
}

// encryptLookupTableBody envelope-encrypts the CSV body with the customer
// managed KMS key: a data key from KMS protects the body with AES-GCM and the
// wrapped data key is stored alongside the ciphertext.
func (s *LogsService) encryptLookupTableBody(body []byte, kmsKeyID, tableArn string) (encrypted, dataKey, nonce []byte, err error) {
	if s.kmsInvoker() == nil {
		return nil, nil, nil, NewLogsError("InvalidParameterException",
			"KMS encryption is not available", 400)
	}
	encryptionContext := map[string]string{lookupTableEncryptionContextKey: tableArn}
	dk, err := s.kmsInvoker().GenerateDataKey(context.Background(), kmsKeyID, "AES_256", encryptionContext, tableArn)
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
	// The read's parameters were valid: a missing invoker, a key the KMS
	// substrate refuses to decrypt (disabled, rotated, granted to another
	// context) or an undecryptable stored body are server-side conditions
	// of the table's storage, not the caller's request — they surface as
	// the operation's declared ServiceUnavailableException instead of
	// blaming the parameters. (AWS's own surface for an undecryptable
	// table is undocumented on every reachable page; this mapping follows
	// the platform's established substrate-failure identity.)
	if s.kmsInvoker() == nil {
		return "", NewLogsError("ServiceUnavailableException",
			"KMS encryption is not available to read this lookup table", http.StatusServiceUnavailable)
	}
	tableArn := lookupTableArn(region, s.accountID, lt.Name)
	encryptionContext := map[string]string{lookupTableEncryptionContextKey: tableArn}
	plaintextKey, err := s.kmsInvoker().Decrypt(context.Background(), lt.KmsKeyId, lt.EncryptedDataKey, encryptionContext, tableArn)
	if err != nil {
		return "", NewLogsError("ServiceUnavailableException",
			fmt.Sprintf("failed to decrypt lookup table %s: %v", lt.Name, err), http.StatusServiceUnavailable)
	}
	plaintext, err := crypto.AESGCMDecryptWithNonce(plaintextKey, lt.EncryptedBody, lt.ContentNonce)
	if err != nil {
		return "", NewLogsError("ServiceUnavailableException",
			fmt.Sprintf("failed to decrypt lookup table %s: %v", lt.Name, err), http.StatusServiceUnavailable)
	}
	return string(plaintext), nil
}
