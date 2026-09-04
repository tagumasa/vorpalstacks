package testutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"vorpalstacks-sdk-tests/config"
)

type kmsTestContext struct {
	client *kms.Client
	ctx    context.Context

	region    string
	accountID string

	keyID     string
	keyAlias  string
	keyARN    string
	rsaKeyID  string
	hmacKeyID string

	ciphertextBlob []byte
	signature      []byte
	macValue       []byte

	cleanupKeyIDs  []string
	cleanupAliases []string
}

func (tc *kmsTestContext) addCleanupKey(keyID string) {
	tc.cleanupKeyIDs = append(tc.cleanupKeyIDs, keyID)
}

func (tc *kmsTestContext) addCleanupAlias(alias string) {
	tc.cleanupAliases = append(tc.cleanupAliases, alias)
}

func (tc *kmsTestContext) scheduleDeletion(keyID string) {
	tc.client.ScheduleKeyDeletion(tc.ctx, &kms.ScheduleKeyDeletionInput{
		KeyId:               aws.String(keyID),
		PendingWindowInDays: aws.Int32(7),
	})
}

// requireKeyID reports whether the shared symmetric-key fixture is available.
func (tc *kmsTestContext) requireKeyID() error {
	if tc.keyID == "" {
		return fmt.Errorf("key ID not available")
	}
	return nil
}

// requireRSAKeyID reports whether the shared RSA-key fixture is available.
func (tc *kmsTestContext) requireRSAKeyID() error {
	if tc.rsaKeyID == "" {
		return fmt.Errorf("RSA key ID not available")
	}
	return nil
}

// requireHMACKeyID reports whether the shared HMAC-key fixture is available.
func (tc *kmsTestContext) requireHMACKeyID() error {
	if tc.hmacKeyID == "" {
		return fmt.Errorf("HMAC key ID not available")
	}
	return nil
}

// createKey creates a key and registers it for end-of-run deletion.
func (tc *kmsTestContext) createKey(input *kms.CreateKeyInput) (*types.KeyMetadata, error) {
	resp, err := tc.client.CreateKey(tc.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create key failed: %w", err)
	}
	tc.addCleanupKey(*resp.KeyMetadata.KeyId)
	return resp.KeyMetadata, nil
}

// createAlias creates an alias and registers it for end-of-run deletion.
func (tc *kmsTestContext) createAlias(aliasName, targetKeyID string) error {
	if _, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
		AliasName:   aws.String(aliasName),
		TargetKeyId: aws.String(targetKeyID),
	}); err != nil {
		return fmt.Errorf("create alias failed: %w", err)
	}
	tc.addCleanupAlias(aliasName)
	return nil
}

// kmsAssertErrorAs asserts the error chain carries the given AWS error type.
func kmsAssertErrorAs[T any, PT interface {
	*T
	error
}](err error) error {
	var target PT
	if !errors.As(err, &target) {
		return fmt.Errorf("expected %T, got: %T: %v", target, err, err)
	}
	return nil
}

// nonexistentKeyARN builds the ARN of a key UUID that never exists.
func (tc *kmsTestContext) nonexistentKeyARN() string {
	return fmt.Sprintf("arn:aws:kms:%s:%s:key/ffffffff-ffff-ffff-ffff-ffffffffffff", tc.region, tc.accountID)
}

// iamUserARN builds an IAM user principal ARN under the test account.
func (tc *kmsTestContext) iamUserARN(userName string) string {
	return fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, userName)
}

func (r *TestRunner) RunKMSTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "kms",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	tc := &kmsTestContext{
		client:    kms.NewFromConfig(cfg),
		ctx:       context.Background(),
		region:    r.region,
		accountID: r.accountID,
	}

	results = append(results, r.runKMSKeyTests(tc)...)
	results = append(results, r.runKMSCryptoTests(tc)...)
	results = append(results, r.runKMSSignTests(tc)...)
	results = append(results, r.runKMSAliasTests(tc)...)
	results = append(results, r.runKMSGrantTests(tc)...)
	results = append(results, r.runKMSPolicyTests(tc)...)
	results = append(results, r.runKMSTagTests(tc)...)
	results = append(results, r.runKMSRotationTests(tc)...)
	results = append(results, r.runKMSImportTests(tc)...)
	results = append(results, r.runKMSMultiRegionTests(tc)...)
	results = append(results, r.runKMSEdgeTests(tc)...)

	for _, alias := range tc.cleanupAliases {
		tc.client.DeleteAlias(tc.ctx, &kms.DeleteAliasInput{AliasName: aws.String(alias)})
	}
	for _, kid := range tc.cleanupKeyIDs {
		tc.scheduleDeletion(kid)
	}
	if tc.rsaKeyID != "" {
		tc.scheduleDeletion(tc.rsaKeyID)
	}
	if tc.hmacKeyID != "" {
		tc.scheduleDeletion(tc.hmacKeyID)
	}

	return results
}
