package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"vorpalstacks-sdk-tests/config"
)

type kmsTestContext struct {
	client *kms.Client
	ctx    context.Context

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
		client: kms.NewFromConfig(cfg),
		ctx:    context.Background(),
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
