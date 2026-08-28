package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamAdvancedTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	// ========== OpenID Connect Provider ==========
	oidcUrl := "https://oidc.example.com/" + tc.ts
	oidcThumbprint := "9e99a48a9960b14926bb7f3b02e22da0b5ec98f2"

	results = append(results, r.RunTest("iam", "CreateOpenIDConnectProvider", func() error {
		resp, err := tc.client.CreateOpenIDConnectProvider(tc.ctx, &iam.CreateOpenIDConnectProviderInput{
			Url:            aws.String(oidcUrl),
			ThumbprintList: []string{oidcThumbprint},
			ClientIDList:   []string{"my-client-id"},
			Tags: []types.Tag{
				{Key: aws.String("Source"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp.OpenIDConnectProviderArn == nil {
			return fmt.Errorf("oidc provider arn is nil")
		}
		tc.oidcProviderArn = *resp.OpenIDConnectProviderArn
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetOpenIDConnectProvider", func() error {
		// An ARN shorter than the documented 20-character minimum is
		// malformed input, not a missing provider.
		if _, err := tc.client.GetOpenIDConnectProvider(tc.ctx, &iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String("x"),
		}); err == nil {
			return fmt.Errorf("a malformed provider ARN must be rejected")
		} else if !containsErrorCode(err, "InvalidInput") {
			return fmt.Errorf("malformed provider ARN: got %v, want InvalidInput", err)
		}

		resp, err := tc.client.GetOpenIDConnectProvider(tc.ctx, &iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.Url) != oidcUrl {
			return fmt.Errorf("url mismatch: got %s, want %s", aws.ToString(resp.Url), oidcUrl)
		}
		clientFound := false
		for _, id := range resp.ClientIDList {
			if id == "my-client-id" {
				clientFound = true
				break
			}
		}
		if !clientFound {
			return fmt.Errorf("client ID not found in GetOpenIDConnectProvider")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListOpenIDConnectProviders", func() error {
		resp, err := tc.client.ListOpenIDConnectProviders(tc.ctx, &iam.ListOpenIDConnectProvidersInput{})
		if err != nil {
			return err
		}
		if resp.OpenIDConnectProviderList == nil {
			return fmt.Errorf("oidc provider list is nil")
		}
		found := false
		for _, p := range resp.OpenIDConnectProviderList {
			if aws.ToString(p.Arn) == tc.oidcProviderArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("oidc provider %s not found in list", tc.oidcProviderArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "AddClientIDToOpenIDConnectProvider", func() error {
		_, err := tc.client.AddClientIDToOpenIDConnectProvider(tc.ctx, &iam.AddClientIDToOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
			ClientID:                 aws.String("second-client-id"),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetOpenIDConnectProvider(tc.ctx, &iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
		})
		if err != nil {
			return err
		}
		found := false
		for _, id := range resp.ClientIDList {
			if id == "second-client-id" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("second-client-id not found after AddClientID")
		}
		return nil
	}))

	// Concurrent AddClientID race regression guard.  Five goroutines
	// each add a distinct client ID to the same provider.  Under the old
	// unlocked Get-modify-Put path, at least one update was lost.  The
	// atomic store-side AddClientID holds the per-ARN lock across the
	// cycle, so all five IDs must be present afterwards.
	results = append(results, r.RunTest("iam", "AddClientIDToOpenIDConnectProvider_Concurrent", func() error {
		const concurrency = 5
		ids := make([]string, concurrency)
		for i := range ids {
			ids[i] = fmt.Sprintf("concurrent-client-%d-%s", i, tc.ts)
		}
		errs := make(chan error, concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id string) {
				_, err := tc.client.AddClientIDToOpenIDConnectProvider(tc.ctx, &iam.AddClientIDToOpenIDConnectProviderInput{
					OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
					ClientID:                 aws.String(id),
				})
				errs <- err
			}(ids[i])
		}
		for i := 0; i < concurrency; i++ {
			if err := <-errs; err != nil {
				return fmt.Errorf("concurrent AddClientID failed: %w", err)
			}
		}
		resp, err := tc.client.GetOpenIDConnectProvider(tc.ctx, &iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
		})
		if err != nil {
			return err
		}
		present := make(map[string]bool, len(resp.ClientIDList))
		for _, id := range resp.ClientIDList {
			present[id] = true
		}
		for _, id := range ids {
			if !present[id] {
				return fmt.Errorf("client id %s missing after concurrent AddClientID (race regression)", id)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateOpenIDConnectProviderThumbprint", func() error {
		newThumbprint := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		_, err := tc.client.UpdateOpenIDConnectProviderThumbprint(tc.ctx, &iam.UpdateOpenIDConnectProviderThumbprintInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
			ThumbprintList:           []string{newThumbprint},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetOpenIDConnectProvider(tc.ctx, &iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
		})
		if err != nil {
			return err
		}
		thumbFound := false
		for _, tp := range resp.ThumbprintList {
			if tp == newThumbprint {
				thumbFound = true
				break
			}
		}
		if !thumbFound {
			return fmt.Errorf("updated thumbprint not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "RemoveClientIDFromOpenIDConnectProvider", func() error {
		_, err := tc.client.RemoveClientIDFromOpenIDConnectProvider(tc.ctx, &iam.RemoveClientIDFromOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
			ClientID:                 aws.String("second-client-id"),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetOpenIDConnectProvider(tc.ctx, &iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
		})
		if err != nil {
			return err
		}
		for _, id := range resp.ClientIDList {
			if id == "second-client-id" {
				return fmt.Errorf("second-client-id should be removed")
			}
		}
		return nil
	}))

	// OIDC tags
	results = append(results, r.RunTest("iam", "TagOpenIDConnectProvider", func() error {
		_, err := tc.client.TagOpenIDConnectProvider(tc.ctx, &iam.TagOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListOpenIDConnectProviderTags(tc.ctx, &iam.ListOpenIDConnectProviderTagsInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
		})
		if err != nil {
			return fmt.Errorf("ListOpenIDConnectProviderTags after tag: %w", err)
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found after TagOpenIDConnectProvider")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListOpenIDConnectProviderTags", func() error {
		resp, err := tc.client.ListOpenIDConnectProviderTags(tc.ctx, &iam.ListOpenIDConnectProviderTagsInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
		})
		if err != nil {
			return err
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagOpenIDConnectProvider", func() error {
		_, err := tc.client.UntagOpenIDConnectProvider(tc.ctx, &iam.UntagOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
			TagKeys:                  []string{"Environment"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListOpenIDConnectProviderTags(tc.ctx, &iam.ListOpenIDConnectProviderTagsInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
		})
		if err != nil {
			return fmt.Errorf("ListOpenIDConnectProviderTags after untag: %w", err)
		}
		if iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag still present after UntagOpenIDConnectProvider")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteOpenIDConnectProvider", func() error {
		_, err := tc.client.DeleteOpenIDConnectProvider(tc.ctx, &iam.DeleteOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetOpenIDConnectProvider(tc.ctx, &iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: aws.String(tc.oidcProviderArn),
		})
		if err == nil {
			return fmt.Errorf("GetOpenIDConnectProvider should fail after DeleteOpenIDConnectProvider")
		}
		return nil
	}))

	// ========== SSH Public Key ==========
	sshUserName := fmt.Sprintf("SSHUser-%s", tc.ts)
	testSSHPublicKey := sshRsaPublicKeyBody(generateTestRSAKey(), "advanced@example.com")

	results = append(results, r.RunTest("iam", "_Advanced_CreateUserForSSH", func() error {
		_, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{
			UserName: aws.String(sshUserName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "UploadSSHPublicKey", func() error {
		resp, err := tc.client.UploadSSHPublicKey(tc.ctx, &iam.UploadSSHPublicKeyInput{
			UserName:         aws.String(sshUserName),
			SSHPublicKeyBody: aws.String(testSSHPublicKey),
		})
		if err != nil {
			return err
		}
		if resp.SSHPublicKey == nil {
			return fmt.Errorf("ssh public key is nil")
		}
		if aws.ToString(resp.SSHPublicKey.SSHPublicKeyId) == "" {
			return fmt.Errorf("ssh public key id is empty")
		}
		tc.sshPublicKeyId = *resp.SSHPublicKey.SSHPublicKeyId
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListSSHPublicKeys", func() error {
		resp, err := tc.client.ListSSHPublicKeys(tc.ctx, &iam.ListSSHPublicKeysInput{
			UserName: aws.String(sshUserName),
		})
		if err != nil {
			return err
		}
		found := false
		for _, k := range resp.SSHPublicKeys {
			if aws.ToString(k.SSHPublicKeyId) == tc.sshPublicKeyId {
				found = true
				if aws.ToString(k.UserName) != sshUserName {
					return fmt.Errorf("username mismatch in SSH key")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("ssh public key %s not found", tc.sshPublicKeyId)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetSSHPublicKey", func() error {
		resp, err := tc.client.GetSSHPublicKey(tc.ctx, &iam.GetSSHPublicKeyInput{
			UserName:       aws.String(sshUserName),
			SSHPublicKeyId: aws.String(tc.sshPublicKeyId),
			Encoding:       types.EncodingTypeSsh,
		})
		if err != nil {
			return err
		}
		if resp.SSHPublicKey == nil {
			return fmt.Errorf("ssh public key is nil")
		}
		if aws.ToString(resp.SSHPublicKey.SSHPublicKeyBody) == "" {
			return fmt.Errorf("ssh public key body is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateSSHPublicKey", func() error {
		_, err := tc.client.UpdateSSHPublicKey(tc.ctx, &iam.UpdateSSHPublicKeyInput{
			UserName:       aws.String(sshUserName),
			SSHPublicKeyId: aws.String(tc.sshPublicKeyId),
			Status:         types.StatusTypeInactive,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetSSHPublicKey(tc.ctx, &iam.GetSSHPublicKeyInput{
			UserName:       aws.String(sshUserName),
			SSHPublicKeyId: aws.String(tc.sshPublicKeyId),
			Encoding:       types.EncodingTypeSsh,
		})
		if err != nil {
			return err
		}
		if resp.SSHPublicKey.Status != types.StatusTypeInactive {
			return fmt.Errorf("expected Inactive status, got %s", resp.SSHPublicKey.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteSSHPublicKey", func() error {
		_, err := tc.client.DeleteSSHPublicKey(tc.ctx, &iam.DeleteSSHPublicKeyInput{
			UserName:       aws.String(sshUserName),
			SSHPublicKeyId: aws.String(tc.sshPublicKeyId),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListSSHPublicKeys(tc.ctx, &iam.ListSSHPublicKeysInput{
			UserName: aws.String(sshUserName),
		})
		if err != nil {
			return fmt.Errorf("ListSSHPublicKeys after delete: %w", err)
		}
		for _, k := range resp.SSHPublicKeys {
			if aws.ToString(k.SSHPublicKeyId) == tc.sshPublicKeyId {
				return fmt.Errorf("ssh key %s still present after DeleteSSHPublicKey", tc.sshPublicKeyId)
			}
		}
		return nil
	}))

	// ========== Server Certificate ==========
	serverCertName := fmt.Sprintf("TestCert-%s", tc.ts)
	testCertBody, testPrivateKey, err := generateSelfSignedCertificate(serverCertName)
	if err != nil {
		return []TestResult{{Service: "iam", TestName: "Setup", Status: "FAIL", Error: err.Error()}}
	}

	results = append(results, r.RunTest("iam", "UploadServerCertificate", func() error {
		resp, err := tc.client.UploadServerCertificate(tc.ctx, &iam.UploadServerCertificateInput{
			ServerCertificateName: aws.String(serverCertName),
			CertificateBody:       aws.String(testCertBody),
			PrivateKey:            aws.String(testPrivateKey),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp.ServerCertificateMetadata == nil {
			return fmt.Errorf("server certificate metadata is nil")
		}
		if aws.ToString(resp.ServerCertificateMetadata.ServerCertificateName) != serverCertName {
			return fmt.Errorf("server certificate name mismatch")
		}
		if len(resp.Tags) == 0 {
			return fmt.Errorf("UploadServerCertificate response missing Tags")
		}
		if resp.ServerCertificateMetadata.Expiration == nil {
			return fmt.Errorf("UploadServerCertificate response missing Expiration")
		}
		tc.serverCertArn = aws.ToString(resp.ServerCertificateMetadata.Arn)
		return nil
	}))

	results = append(results, r.RunTest("iam", "UploadServerCertificate_Malformed", func() error {
		_, err := tc.client.UploadServerCertificate(tc.ctx, &iam.UploadServerCertificateInput{
			ServerCertificateName: aws.String(serverCertName + "-malformed"),
			CertificateBody:       aws.String("garbage certificate body"),
			PrivateKey:            aws.String(testPrivateKey),
		})
		if err == nil {
			return fmt.Errorf("a malformed certificate body must be rejected")
		}
		if !containsErrorCode(err, "MalformedCertificate") {
			return fmt.Errorf("malformed certificate: got %v, want MalformedCertificate", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UploadServerCertificate_KeyPairMismatch", func() error {
		_, otherKey, err := generateSelfSignedCertificate(serverCertName + "-other")
		if err != nil {
			return err
		}
		_, err = tc.client.UploadServerCertificate(tc.ctx, &iam.UploadServerCertificateInput{
			ServerCertificateName: aws.String(serverCertName + "-mismatch"),
			CertificateBody:       aws.String(testCertBody),
			PrivateKey:            aws.String(otherKey),
		})
		if err == nil {
			return fmt.Errorf("a private key from another certificate must be rejected")
		}
		if !containsErrorCode(err, "KeyPairMismatch") {
			return fmt.Errorf("key pair mismatch: got %v, want KeyPairMismatch", err)
		}
		return nil
	}))

	// UploadServerCertificate without PrivateKey (a required field) must fail
	results = append(results, r.RunTest("iam", "Error_UploadServerCertificate_MissingPrivateKey", func() error {
		_, err := tc.client.UploadServerCertificate(tc.ctx, &iam.UploadServerCertificateInput{
			ServerCertificateName: aws.String(serverCertName + "-nokey"),
			CertificateBody:       aws.String(testCertBody),
		})
		if err == nil {
			return fmt.Errorf("expected validation error for missing PrivateKey")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetServerCertificate", func() error {
		resp, err := tc.client.GetServerCertificate(tc.ctx, &iam.GetServerCertificateInput{
			ServerCertificateName: aws.String(serverCertName),
		})
		if err != nil {
			return err
		}
		if resp.ServerCertificate == nil {
			return fmt.Errorf("server certificate is nil")
		}
		if aws.ToString(resp.ServerCertificate.CertificateBody) == "" {
			return fmt.Errorf("certificate body is empty")
		}
		if len(resp.ServerCertificate.Tags) == 0 {
			return fmt.Errorf("GetServerCertificate response missing Tags")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListServerCertificates", func() error {
		resp, err := tc.client.ListServerCertificates(tc.ctx, &iam.ListServerCertificatesInput{})
		if err != nil {
			return err
		}
		found := false
		for _, cert := range resp.ServerCertificateMetadataList {
			if aws.ToString(cert.ServerCertificateName) == serverCertName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("server certificate %s not found in list", serverCertName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateServerCertificate", func() error {
		newName := serverCertName + "-renamed"
		_, err := tc.client.UpdateServerCertificate(tc.ctx, &iam.UpdateServerCertificateInput{
			ServerCertificateName:    aws.String(serverCertName),
			NewServerCertificateName: aws.String(newName),
		})
		if err != nil {
			return err
		}
		serverCertName = newName
		resp, err := tc.client.GetServerCertificate(tc.ctx, &iam.GetServerCertificateInput{
			ServerCertificateName: aws.String(newName),
		})
		if err != nil {
			return fmt.Errorf("GetServerCertificate with new name: %w", err)
		}
		if aws.ToString(resp.ServerCertificate.ServerCertificateMetadata.ServerCertificateName) != newName {
			return fmt.Errorf("renamed cert name mismatch: got %s", aws.ToString(resp.ServerCertificate.ServerCertificateMetadata.ServerCertificateName))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateServerCertificate_InvalidNewName", func() error {
		_, err := tc.client.UpdateServerCertificate(tc.ctx, &iam.UpdateServerCertificateInput{
			ServerCertificateName:    aws.String(serverCertName),
			NewServerCertificateName: aws.String("invalid name with spaces!"),
		})
		if err == nil {
			return fmt.Errorf("an invalid new certificate name must be rejected")
		}
		if !isInvalidInputError(err) {
			return fmt.Errorf("invalid new name: got %v, want InvalidInput", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateServerCertificate_DuplicateName", func() error {
		otherName := serverCertName + "-occupied"
		otherCert, otherKey, err := generateSelfSignedCertificate(otherName)
		if err != nil {
			return err
		}
		if _, err := tc.client.UploadServerCertificate(tc.ctx, &iam.UploadServerCertificateInput{
			ServerCertificateName: aws.String(otherName),
			CertificateBody:       aws.String(otherCert),
			PrivateKey:            aws.String(otherKey),
		}); err != nil {
			return err
		}
		defer tc.client.DeleteServerCertificate(tc.ctx, &iam.DeleteServerCertificateInput{
			ServerCertificateName: aws.String(otherName),
		})

		_, err = tc.client.UpdateServerCertificate(tc.ctx, &iam.UpdateServerCertificateInput{
			ServerCertificateName:    aws.String(serverCertName),
			NewServerCertificateName: aws.String(otherName),
		})
		if err == nil {
			return fmt.Errorf("renaming onto an existing certificate name must be rejected")
		}
		if !containsErrorCode(err, "EntityAlreadyExists") {
			return fmt.Errorf("duplicate new name: got %v, want EntityAlreadyExists", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "TagServerCertificate", func() error {
		_, err := tc.client.TagServerCertificate(tc.ctx, &iam.TagServerCertificateInput{
			ServerCertificateName: aws.String(serverCertName),
			Tags: []types.Tag{
				{Key: aws.String("Env"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListServerCertificateTags(tc.ctx, &iam.ListServerCertificateTagsInput{
			ServerCertificateName: aws.String(serverCertName),
		})
		if err != nil {
			return fmt.Errorf("ListServerCertificateTags after tag: %w", err)
		}
		if !iamTagPresent(resp.Tags, "Env", "test") {
			return fmt.Errorf("Env=test tag not found after TagServerCertificate")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListServerCertificateTags", func() error {
		resp, err := tc.client.ListServerCertificateTags(tc.ctx, &iam.ListServerCertificateTagsInput{
			ServerCertificateName: aws.String(serverCertName),
		})
		if err != nil {
			return err
		}
		if !iamTagPresent(resp.Tags, "Env", "test") {
			return fmt.Errorf("Env=test tag not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagServerCertificate", func() error {
		_, err := tc.client.UntagServerCertificate(tc.ctx, &iam.UntagServerCertificateInput{
			ServerCertificateName: aws.String(serverCertName),
			TagKeys:               []string{"Env"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListServerCertificateTags(tc.ctx, &iam.ListServerCertificateTagsInput{
			ServerCertificateName: aws.String(serverCertName),
		})
		if err != nil {
			return fmt.Errorf("ListServerCertificateTags after untag: %w", err)
		}
		if iamTagPresent(resp.Tags, "Env", "test") {
			return fmt.Errorf("Env=test tag still present after UntagServerCertificate")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteServerCertificate", func() error {
		_, err := tc.client.DeleteServerCertificate(tc.ctx, &iam.DeleteServerCertificateInput{
			ServerCertificateName: aws.String(serverCertName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetServerCertificate(tc.ctx, &iam.GetServerCertificateInput{
			ServerCertificateName: aws.String(serverCertName),
		})
		if err == nil {
			return fmt.Errorf("GetServerCertificate should fail after DeleteServerCertificate")
		}
		return nil
	}))

	// ========== Service-Specific Credential ==========
	results = append(results, r.RunTest("iam", "CreateServiceSpecificCredential", func() error {
		resp, err := tc.client.CreateServiceSpecificCredential(tc.ctx, &iam.CreateServiceSpecificCredentialInput{
			UserName:          aws.String(sshUserName),
			ServiceName:       aws.String("codecommit.amazonaws.com"),
			CredentialAgeDays: aws.Int32(30),
		})
		if err != nil {
			return err
		}
		if resp.ServiceSpecificCredential == nil {
			return fmt.Errorf("service specific credential is nil")
		}
		if aws.ToString(resp.ServiceSpecificCredential.ServiceUserName) == "" {
			return fmt.Errorf("service user name is empty")
		}
		// ExpirationDate must be set when CredentialAgeDays is specified
		if resp.ServiceSpecificCredential.ExpirationDate == nil {
			return fmt.Errorf("ExpirationDate is nil despite CredentialAgeDays=30")
		}
		tc.serviceCredId = aws.ToString(resp.ServiceSpecificCredential.ServiceSpecificCredentialId)
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListServiceSpecificCredentials", func() error {
		resp, err := tc.client.ListServiceSpecificCredentials(tc.ctx, &iam.ListServiceSpecificCredentialsInput{
			UserName: aws.String(sshUserName),
		})
		if err != nil {
			return err
		}
		found := false
		for _, c := range resp.ServiceSpecificCredentials {
			if aws.ToString(c.ServiceSpecificCredentialId) == tc.serviceCredId {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("service specific credential %s not found", tc.serviceCredId)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateServiceSpecificCredential", func() error {
		_, err := tc.client.UpdateServiceSpecificCredential(tc.ctx, &iam.UpdateServiceSpecificCredentialInput{
			UserName:                    aws.String(sshUserName),
			ServiceSpecificCredentialId: aws.String(tc.serviceCredId),
			Status:                      types.StatusTypeInactive,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListServiceSpecificCredentials(tc.ctx, &iam.ListServiceSpecificCredentialsInput{
			UserName: aws.String(sshUserName),
		})
		if err != nil {
			return err
		}
		for _, c := range resp.ServiceSpecificCredentials {
			if aws.ToString(c.ServiceSpecificCredentialId) == tc.serviceCredId {
				if c.Status != types.StatusTypeInactive {
					return fmt.Errorf("expected Inactive status, got %s", c.Status)
				}
				return nil
			}
		}
		return fmt.Errorf("service credential not found after update")
	}))

	results = append(results, r.RunTest("iam", "ServiceSpecificCredential_Ownership", func() error {
		other := fmt.Sprintf("SSCOwner-%s", tc.ts)
		cleanupOther, err := tc.createUser(other)
		if err != nil {
			return err
		}
		defer cleanupOther()

		// Mutating the credential through a user that does not own it
		// fails with NoSuchEntity.
		_, err = tc.client.UpdateServiceSpecificCredential(tc.ctx, &iam.UpdateServiceSpecificCredentialInput{
			UserName:                    aws.String(other),
			ServiceSpecificCredentialId: aws.String(tc.serviceCredId),
			Status:                      types.StatusTypeActive,
		})
		if err == nil {
			return fmt.Errorf("update through a non-owner must be rejected")
		}
		if !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("non-owner update: got %v, want NoSuchEntity", err)
		}

		_, err = tc.client.DeleteServiceSpecificCredential(tc.ctx, &iam.DeleteServiceSpecificCredentialInput{
			UserName:                    aws.String(other),
			ServiceSpecificCredentialId: aws.String(tc.serviceCredId),
		})
		if err == nil {
			return fmt.Errorf("delete through a non-owner must be rejected")
		}
		if !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("non-owner delete: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "CreateServiceSpecificCredential_UnsupportedService", func() error {
		_, err := tc.client.CreateServiceSpecificCredential(tc.ctx, &iam.CreateServiceSpecificCredentialInput{
			UserName:    aws.String(sshUserName),
			ServiceName: aws.String("s3.amazonaws.com"),
		})
		if err == nil {
			return fmt.Errorf("a service without service-specific credential support must be rejected")
		}
		if !containsErrorCode(err, "NotSupportedService") {
			return fmt.Errorf("unsupported service: got %v, want NotSupportedService", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ResetServiceSpecificCredential", func() error {
		resp, err := tc.client.ResetServiceSpecificCredential(tc.ctx, &iam.ResetServiceSpecificCredentialInput{
			UserName:                    aws.String(sshUserName),
			ServiceSpecificCredentialId: aws.String(tc.serviceCredId),
		})
		if err != nil {
			return err
		}
		if resp.ServiceSpecificCredential == nil {
			return fmt.Errorf("service specific credential is nil in reset response")
		}
		if aws.ToString(resp.ServiceSpecificCredential.ServicePassword) == "" {
			return fmt.Errorf("new service password is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteServiceSpecificCredential", func() error {
		_, err := tc.client.DeleteServiceSpecificCredential(tc.ctx, &iam.DeleteServiceSpecificCredentialInput{
			UserName:                    aws.String(sshUserName),
			ServiceSpecificCredentialId: aws.String(tc.serviceCredId),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListServiceSpecificCredentials(tc.ctx, &iam.ListServiceSpecificCredentialsInput{
			UserName: aws.String(sshUserName),
		})
		if err != nil {
			return fmt.Errorf("ListServiceSpecificCredentials after delete: %w", err)
		}
		for _, c := range resp.ServiceSpecificCredentials {
			if aws.ToString(c.ServiceSpecificCredentialId) == tc.serviceCredId {
				return fmt.Errorf("credential %s still present after delete", tc.serviceCredId)
			}
		}
		return nil
	}))

	// Cleanup temp user
	results = append(results, r.RunTest("iam", "_Advanced_DeleteTempUser", func() error {
		_, err := tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{
			UserName: aws.String(sshUserName),
		})
		return err
	}))

	// ========== Credential Report ==========
	results = append(results, r.RunTest("iam", "GenerateCredentialReport", func() error {
		resp, err := tc.client.GenerateCredentialReport(tc.ctx, &iam.GenerateCredentialReportInput{})
		if err != nil {
			return err
		}
		if resp.State != types.ReportStateTypeComplete && resp.State != types.ReportStateTypeStarted {
			return fmt.Errorf("unexpected report state: %v", resp.State)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetCredentialReport", func() error {
		var lastErr error
		for i := 0; i < 10; i++ {
			resp, err := tc.client.GetCredentialReport(tc.ctx, &iam.GetCredentialReportInput{})
			if err == nil {
				if resp.Content == nil || len(resp.Content) == 0 {
					return fmt.Errorf("credential report content is empty")
				}
				if resp.ReportFormat != "text/csv" {
					return fmt.Errorf("report format mismatch: %v", resp.ReportFormat)
				}
				return nil
			}
			if !strings.Contains(err.Error(), "ReportInProgress") {
				return err
			}
			lastErr = err
			time.Sleep(500 * time.Millisecond)
		}
		return lastErr
	}))

	// ========== Service Last Accessed ==========
	results = append(results, r.RunTest("iam", "GenerateServiceLastAccessedDetails", func() error {
		resp, err := tc.client.GenerateServiceLastAccessedDetails(tc.ctx, &iam.GenerateServiceLastAccessedDetailsInput{
			Arn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		if resp.JobId == nil {
			return fmt.Errorf("job id is nil")
		}
		return nil
	}))

	return results
}
