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
	testSSHPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0g+Z+s+Z+s+Z test@example.com"

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
	testCertBody := `-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHHCgVZU65BMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3QgY2EwHhcNMjQwMTAxMDAwMDAwWhcNMjUwMTAxMDAwMDAwWjARMQ8wDQYDVQQD
DAZ0ZXN0IGNhMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0g+Z+s+Z+s+Z+
s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s
IDAQABMA0GCSqGSIb3DQEBCwUAA4GBAKHHCgVZU65BMA0GCSqGSIb3DQEBCwUAMB
ExDzANBgNVBAMMBnRlc3QgY2EwHhcNMjQwMTAxMDAwMDAwWhcNMjUwMTAxMDAwM
DAwWjARMQ8wDQYDVQQDDAZ0ZXN0IGNh
-----END CERTIFICATE-----`
	testPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIBkTCB+wIJAKHHCgVZU65BMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3QgY2EwHhcNMjQwMTAxMDAwMDAwWhcNMjUwMTAxMDAwMDAwWjARMQ8wDQYDVQQD
DAZ0ZXN0IGNhMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0g+Z+s+Z+s+Z+
s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s+Z+s
IDAQABMA0GCSqGSIb3DQEBCwUAA4GBAKHHCgVZU65BMA0GCSqGSIb3DQEBCwUAMB
-----END RSA PRIVATE KEY-----`

	results = append(results, r.RunTest("iam", "UploadServerCertificate", func() error {
		resp, err := tc.client.UploadServerCertificate(tc.ctx, &iam.UploadServerCertificateInput{
			ServerCertificateName: aws.String(serverCertName),
			CertificateBody:       aws.String(testCertBody),
			PrivateKey:            aws.String(testPrivateKey),
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
		tc.serverCertArn = aws.ToString(resp.ServerCertificateMetadata.Arn)
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
			UserName:    aws.String(sshUserName),
			ServiceName: aws.String("codecommit.amazonaws.com"),
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
