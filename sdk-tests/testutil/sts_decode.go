package testutil

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func (r *TestRunner) runSTSDecodeTests(tc *stsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sts", "DecodeAuthorizationMessage_JSON", func() error {
		originalMsg := `{"ErrorCode":"AccessDenied","Message":"Not authorized"}`
		encoded := base64.StdEncoding.EncodeToString([]byte(originalMsg))
		resp, err := tc.client.DecodeAuthorizationMessage(tc.ctx, &sts.DecodeAuthorizationMessageInput{
			EncodedMessage: aws.String(encoded),
		})
		if err != nil {
			return err
		}
		if resp.DecodedMessage == nil || *resp.DecodedMessage != originalMsg {
			return fmt.Errorf("decoded message mismatch, got: %v", resp.DecodedMessage)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "DecodeAuthorizationMessage_PlainText", func() error {
		originalMsg := "Plain text error message"
		encoded := base64.StdEncoding.EncodeToString([]byte(originalMsg))
		resp, err := tc.client.DecodeAuthorizationMessage(tc.ctx, &sts.DecodeAuthorizationMessageInput{
			EncodedMessage: aws.String(encoded),
		})
		if err != nil {
			return err
		}
		if resp.DecodedMessage == nil || *resp.DecodedMessage != originalMsg {
			return fmt.Errorf("decoded message mismatch, got: %v", resp.DecodedMessage)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "DecodeAuthorizationMessage_InvalidBase64", func() error {
		_, err := tc.client.DecodeAuthorizationMessage(tc.ctx, &sts.DecodeAuthorizationMessageInput{
			EncodedMessage: aws.String("not-valid-base64!!!"),
		})
		if err := AssertErrorContains(err, "InvalidAuthorizationMessageException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "DecodeAuthorizationMessage_Empty", func() error {
		_, err := tc.client.DecodeAuthorizationMessage(tc.ctx, &sts.DecodeAuthorizationMessageInput{
			EncodedMessage: aws.String(""),
		})
		if err := AssertErrorContains(err, "InvalidAuthorizationMessageException"); err != nil {
			return err
		}
		return nil
	}))

	// GetAccessKeyInfo: verify against a real temporary access key
	// obtained from GetSessionToken. The server validates key
	// existence; fake key IDs are rejected with
	// InvalidClientTokenId.
	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_RealTempKey", func() error {
		// Obtain a temporary session to get a real ASIA-prefixed key.
		stsResp, err := tc.client.GetSessionToken(tc.ctx, &sts.GetSessionTokenInput{
			DurationSeconds: aws.Int32(900),
		})
		if err != nil {
			return fmt.Errorf("failed to get session token: %w", err)
		}
		resp, err := tc.client.GetAccessKeyInfo(tc.ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: stsResp.Credentials.AccessKeyId,
		})
		if err != nil {
			return err
		}
		if resp.Account == nil || *resp.Account == "" {
			return fmt.Errorf("account is nil or empty for real temporary key")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_NonExistentKey", func() error {
		_, err := tc.client.GetAccessKeyInfo(tc.ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: aws.String("AKIAIOSFODNN7EXAMPLE"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent access key, got nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_Invalid", func() error {
		_, err := tc.client.GetAccessKeyInfo(tc.ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: aws.String(""),
		})
		if err := AssertErrorContains(err, "InvalidAccessKeyId"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
