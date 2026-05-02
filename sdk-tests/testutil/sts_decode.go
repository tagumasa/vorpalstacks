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
		if err := AssertErrorContains(err, "InvalidEncodedMessage"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "DecodeAuthorizationMessage_Empty", func() error {
		_, err := tc.client.DecodeAuthorizationMessage(tc.ctx, &sts.DecodeAuthorizationMessageInput{
			EncodedMessage: aws.String(""),
		})
		if err := AssertErrorContains(err, "InvalidEncodedMessage"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_AKIAPrefix", func() error {
		resp, err := tc.client.GetAccessKeyInfo(tc.ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: aws.String("AKIAIOSFODNN7EXAMPLE"),
		})
		if err != nil {
			return err
		}
		if resp.Account == nil || *resp.Account == "" {
			return fmt.Errorf("account is nil or empty for AKIA prefix")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_ASIAPrefix", func() error {
		resp, err := tc.client.GetAccessKeyInfo(tc.ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: aws.String("ASIAIOSFODNN7EXAMPLE"),
		})
		if err != nil {
			return err
		}
		if resp.Account == nil || *resp.Account == "" {
			return fmt.Errorf("account is nil or empty for ASIA prefix")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_ABIAPrefix", func() error {
		resp, err := tc.client.GetAccessKeyInfo(tc.ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: aws.String("ABIAIOSFODNN7EXAMPLE"),
		})
		if err != nil {
			return err
		}
		if resp.Account == nil || *resp.Account == "" {
			return fmt.Errorf("account is nil or empty for ABIA prefix")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetAccessKeyInfo_ACCAPrefix", func() error {
		resp, err := tc.client.GetAccessKeyInfo(tc.ctx, &sts.GetAccessKeyInfoInput{
			AccessKeyId: aws.String("ACCAIOSFODNN7EXAMPLE"),
		})
		if err != nil {
			return err
		}
		if resp.Account == nil || *resp.Account == "" {
			return fmt.Errorf("account is nil or empty for ACCA prefix")
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
