package testutil

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// sshRsaPublicKeyBody renders an RSA key as an authorized-keys style SSH
// public key line, built from the SSH wire format without external
// dependencies.
func sshRsaPublicKeyBody(key *rsa.PrivateKey, comment string) string {
	buf := new(bytes.Buffer)
	writeSSHString := func(b []byte) {
		var lenBytes [4]byte
		lenBytes[0] = byte(len(b) >> 24)
		lenBytes[1] = byte(len(b) >> 16)
		lenBytes[2] = byte(len(b) >> 8)
		lenBytes[3] = byte(len(b))
		buf.Write(lenBytes[:])
		buf.Write(b)
	}
	writeSSHString([]byte("ssh-rsa"))
	writeSSHString(big.NewInt(int64(key.E)).Bytes())
	modulus := key.PublicKey.N.Bytes()
	if modulus[0]&0x80 != 0 {
		modulus = append([]byte{0}, modulus...)
	}
	writeSSHString(modulus)

	body := "ssh-rsa " + base64.StdEncoding.EncodeToString(buf.Bytes())
	if comment != "" {
		body += " " + comment
	}
	return body
}

func (r *TestRunner) iamSSHPublicKeyTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	newSSHUser := func(suffix string) (string, func()) {
		user := fmt.Sprintf("SSHKey-%s-%s", suffix, tc.ts)
		if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(user)}); err != nil {
			return "", func() {}
		}
		return user, func() {
			_, _ = tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(user)})
		}
	}

	results = append(results, r.RunTest("iam", "UploadSSHPublicKey_Invalid", func() error {
		user, cleanup := newSSHUser("invalid")
		defer cleanup()
		if user == "" {
			return fmt.Errorf("user setup failed")
		}

		_, err := tc.client.UploadSSHPublicKey(tc.ctx, &iam.UploadSSHPublicKeyInput{
			UserName:         aws.String(user),
			SSHPublicKeyBody: aws.String("not an ssh public key at all"),
		})
		if err == nil {
			return fmt.Errorf("an invalid key body must be rejected")
		}
		if !containsErrorCode(err, "InvalidPublicKey") {
			return fmt.Errorf("invalid key: got %v, want InvalidPublicKey", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UploadSSHPublicKey_Duplicate", func() error {
		user, cleanup := newSSHUser("dup")
		defer cleanup()
		if user == "" {
			return fmt.Errorf("user setup failed")
		}

		key := generateTestRSAKey()
		body := sshRsaPublicKeyBody(key, "dup@example.com")
		if _, err := tc.client.UploadSSHPublicKey(tc.ctx, &iam.UploadSSHPublicKeyInput{
			UserName:         aws.String(user),
			SSHPublicKeyBody: aws.String(body),
		}); err != nil {
			return err
		}

		// Re-uploading the same key material under a different comment
		// still duplicates the key.
		_, err := tc.client.UploadSSHPublicKey(tc.ctx, &iam.UploadSSHPublicKeyInput{
			UserName:         aws.String(user),
			SSHPublicKeyBody: aws.String(sshRsaPublicKeyBody(key, "other@example.com")),
		})
		if err == nil {
			return fmt.Errorf("uploading the same key material twice must be rejected")
		}
		if !containsErrorCode(err, "DuplicateSSHPublicKey") {
			return fmt.Errorf("duplicate key: got %v, want DuplicateSSHPublicKey", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UploadSSHPublicKey_Quota", func() error {
		user, cleanup := newSSHUser("quota")
		defer cleanup()
		if user == "" {
			return fmt.Errorf("user setup failed")
		}

		for i := 0; i < 5; i++ {
			body := sshRsaPublicKeyBody(generateTestRSAKey(), fmt.Sprintf("q%d@example.com", i))
			if _, err := tc.client.UploadSSHPublicKey(tc.ctx, &iam.UploadSSHPublicKeyInput{
				UserName:         aws.String(user),
				SSHPublicKeyBody: aws.String(body),
			}); err != nil {
				return fmt.Errorf("upload %d within the quota failed: %w", i+1, err)
			}
		}

		_, err := tc.client.UploadSSHPublicKey(tc.ctx, &iam.UploadSSHPublicKeyInput{
			UserName:         aws.String(user),
			SSHPublicKeyBody: aws.String(sshRsaPublicKeyBody(generateTestRSAKey(), "sixth@example.com")),
		})
		if err == nil {
			return fmt.Errorf("the sixth key must exceed the per-user quota")
		}
		if !containsErrorCode(err, "LimitExceeded") {
			return fmt.Errorf("quota: got %v, want LimitExceeded", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetSSHPublicKey_Encodings", func() error {
		user, cleanup := newSSHUser("enc")
		defer cleanup()
		if user == "" {
			return fmt.Errorf("user setup failed")
		}

		upload, err := tc.client.UploadSSHPublicKey(tc.ctx, &iam.UploadSSHPublicKeyInput{
			UserName:         aws.String(user),
			SSHPublicKeyBody: aws.String(sshRsaPublicKeyBody(generateTestRSAKey(), "enc@example.com")),
		})
		if err != nil {
			return err
		}
		keyID := upload.SSHPublicKey.SSHPublicKeyId

		// The upload response itself carries the key body in SSH form.
		uploadBody := aws.ToString(upload.SSHPublicKey.SSHPublicKeyBody)
		if !strings.HasPrefix(uploadBody, "ssh-rsa") {
			return fmt.Errorf("upload response must include the SSH-form key body, got %q", uploadBody)
		}

		sshResp, err := tc.client.GetSSHPublicKey(tc.ctx, &iam.GetSSHPublicKeyInput{
			UserName:       aws.String(user),
			SSHPublicKeyId: keyID,
			Encoding:       "SSH",
		})
		if err != nil {
			return err
		}
		if body := aws.ToString(sshResp.SSHPublicKey.SSHPublicKeyBody); len(body) < 7 || body[:7] != "ssh-rsa" {
			return fmt.Errorf("SSH encoding must return the ssh-rsa single-line form, got %q", body)
		}
		if body := aws.ToString(sshResp.SSHPublicKey.SSHPublicKeyBody); body != uploadBody {
			return fmt.Errorf("SSH encoding must return the same body as the upload response")
		}

		pemResp, err := tc.client.GetSSHPublicKey(tc.ctx, &iam.GetSSHPublicKeyInput{
			UserName:       aws.String(user),
			SSHPublicKeyId: keyID,
			Encoding:       "PEM",
		})
		if err != nil {
			return err
		}
		if body := aws.ToString(pemResp.SSHPublicKey.SSHPublicKeyBody); !strings.HasPrefix(body, "-----BEGIN PUBLIC KEY-----") {
			return fmt.Errorf("PEM encoding must return a PUBLIC KEY PEM block, got %q", body)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "SSHPublicKey_Ownership", func() error {
		owner, cleanup := newSSHUser("owner")
		defer cleanup()
		if owner == "" {
			return fmt.Errorf("user setup failed")
		}
		other := fmt.Sprintf("SSHKey-other-%s", tc.ts)
		if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(other)}); err != nil {
			return err
		}
		defer tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(other)})

		upload, err := tc.client.UploadSSHPublicKey(tc.ctx, &iam.UploadSSHPublicKeyInput{
			UserName:         aws.String(owner),
			SSHPublicKeyBody: aws.String(sshRsaPublicKeyBody(generateTestRSAKey(), "own@example.com")),
		})
		if err != nil {
			return err
		}
		keyID := upload.SSHPublicKey.SSHPublicKeyId

		if _, err := tc.client.UpdateSSHPublicKey(tc.ctx, &iam.UpdateSSHPublicKeyInput{
			UserName:       aws.String(other),
			SSHPublicKeyId: keyID,
			Status:         "Inactive",
		}); err == nil {
			return fmt.Errorf("updating through a non-owner UserName must be rejected")
		} else if !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("non-owner update: got %v, want NoSuchEntity", err)
		}

		if _, err := tc.client.DeleteSSHPublicKey(tc.ctx, &iam.DeleteSSHPublicKeyInput{
			UserName:       aws.String(other),
			SSHPublicKeyId: keyID,
		}); err == nil {
			return fmt.Errorf("deleting through a non-owner UserName must be rejected")
		} else if !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("non-owner delete: got %v, want NoSuchEntity", err)
		}

		// The owner can still delete the key.
		if _, err := tc.client.DeleteSSHPublicKey(tc.ctx, &iam.DeleteSSHPublicKeyInput{
			UserName:       aws.String(owner),
			SSHPublicKeyId: keyID,
		}); err != nil {
			return err
		}
		return nil
	}))

	return results
}
