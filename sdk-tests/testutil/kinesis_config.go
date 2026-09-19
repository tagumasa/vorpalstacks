package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

func (r *TestRunner) kinesisConfigTests(ctx context.Context, client *kinesis.Client, kmsClient *kms.Client, ts string) []TestResult {
	var results []TestResult

	// Stream encryption resolves its key through the platform's KMS, so
	// the encryption tests adopt a real created key.
	keyResp, err := kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
		Description: aws.String("kinesis stream-encryption test key"),
	})
	if err != nil {
		return []TestResult{SetupFailResult("kinesis", "CreateKey setup failed: %v", err)}
	}
	keyARN := aws.ToString(keyResp.KeyMetadata.Arn)
	defer func() {
		_, _ = kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId: aws.String(keyARN), PendingWindowInDays: aws.Int32(7),
		})
	}()

	streamName := kinesisStream(ts, "cfg")
	if _, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
		StreamName: aws.String(streamName),
		ShardCount: aws.Int32(1),
	}); err != nil {
		return []TestResult{SetupFailResult("kinesis", "CreateStream setup failed: %v", err)}
	}
	if _, err := kinesisDescribeWhenReady(ctx, client, streamName, 10*time.Second); err != nil {
		return []TestResult{SetupFailResult("kinesis", "CreateStream setup failed: stream not ready: %v", err)}
	}

	results = append(results, r.RunTest("kinesis", "EnableEnhancedMonitoring", func() error {
		resp, err := client.EnableEnhancedMonitoring(ctx, &kinesis.EnableEnhancedMonitoringInput{
			StreamName: aws.String(streamName),
			ShardLevelMetrics: []types.MetricsName{
				types.MetricsNameIncomingBytes,
				types.MetricsNameOutgoingBytes,
			},
		})
		if err != nil {
			return err
		}
		if resp.CurrentShardLevelMetrics == nil {
			return fmt.Errorf("CurrentShardLevelMetrics is nil")
		}
		if len(resp.DesiredShardLevelMetrics) != 2 {
			return fmt.Errorf("DesiredShardLevelMetrics: expected 2, got %d", len(resp.DesiredShardLevelMetrics))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DisableEnhancedMonitoring", func() error {
		resp, err := client.DisableEnhancedMonitoring(ctx, &kinesis.DisableEnhancedMonitoringInput{
			StreamName: aws.String(streamName),
			ShardLevelMetrics: []types.MetricsName{
				types.MetricsNameIncomingBytes,
				types.MetricsNameOutgoingBytes,
			},
		})
		if err != nil {
			return err
		}
		if resp.CurrentShardLevelMetrics == nil {
			return fmt.Errorf("CurrentShardLevelMetrics is nil")
		}
		if len(resp.DesiredShardLevelMetrics) != 0 {
			return fmt.Errorf("DesiredShardLevelMetrics: expected 0 after disabling all, got %d", len(resp.DesiredShardLevelMetrics))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "IncreaseStreamRetentionPeriod", func() error {
		_, err := client.IncreaseStreamRetentionPeriod(ctx, &kinesis.IncreaseStreamRetentionPeriodInput{
			StreamName:           aws.String(streamName),
			RetentionPeriodHours: aws.Int32(48),
		})
		if err != nil {
			return err
		}
		descResp, err := client.DescribeStreamSummary(ctx, &kinesis.DescribeStreamSummaryInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return fmt.Errorf("describe summary after retention increase: %v", err)
		}
		if aws.ToInt32(descResp.StreamDescriptionSummary.RetentionPeriodHours) != 48 {
			return fmt.Errorf("RetentionPeriodHours: got %d, want 48", aws.ToInt32(descResp.StreamDescriptionSummary.RetentionPeriodHours))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DecreaseStreamRetentionPeriod", func() error {
		_, err := client.DecreaseStreamRetentionPeriod(ctx, &kinesis.DecreaseStreamRetentionPeriodInput{
			StreamName:           aws.String(streamName),
			RetentionPeriodHours: aws.Int32(24),
		})
		if err != nil {
			return err
		}
		descResp, err := client.DescribeStreamSummary(ctx, &kinesis.DescribeStreamSummaryInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return fmt.Errorf("describe summary after retention decrease: %v", err)
		}
		if aws.ToInt32(descResp.StreamDescriptionSummary.RetentionPeriodHours) != 24 {
			return fmt.Errorf("RetentionPeriodHours: got %d, want 24", aws.ToInt32(descResp.StreamDescriptionSummary.RetentionPeriodHours))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "StartStreamEncryption", func() error {
		_, err := client.StartStreamEncryption(ctx, &kinesis.StartStreamEncryptionInput{
			StreamName:     aws.String(streamName),
			EncryptionType: types.EncryptionTypeKms,
			KeyId:          aws.String(keyARN),
		})
		if err != nil {
			return err
		}
		descResp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return fmt.Errorf("describe after encryption: %v", err)
		}
		if descResp.StreamDescription.EncryptionType != types.EncryptionTypeKms {
			return fmt.Errorf("EncryptionType: expected KMS, got %s", descResp.StreamDescription.EncryptionType)
		}

		// The key must exist: a key ARN that resolves to nothing answers
		// the modelled KMS identity.
		_, err = client.StartStreamEncryption(ctx, &kinesis.StartStreamEncryptionInput{
			StreamName:     aws.String(streamName),
			EncryptionType: types.EncryptionTypeKms,
			KeyId:          aws.String(r.kinesisKMSKeyARN("00000000-0000-0000-0000-000000000000")),
		})
		return expectAWSErrorCode(err, "KMSNotFoundException")
	}))

	results = append(results, r.RunTest("kinesis", "StopStreamEncryption", func() error {
		_, err := client.StopStreamEncryption(ctx, &kinesis.StopStreamEncryptionInput{
			StreamName:     aws.String(streamName),
			EncryptionType: types.EncryptionTypeKms,
			KeyId:          aws.String(keyARN),
		})
		if err != nil {
			return err
		}
		descResp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return fmt.Errorf("describe after stop encryption: %v", err)
		}
		if descResp.StreamDescription.EncryptionType != types.EncryptionTypeNone {
			return fmt.Errorf("EncryptionType: expected NONE, got %s", descResp.StreamDescription.EncryptionType)
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DescribeStream_VerifyEncryption", func() error {
		sn := kinesisStream(ts, "enc")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		_, err = client.StartStreamEncryption(ctx, &kinesis.StartStreamEncryptionInput{
			StreamName:     aws.String(sn),
			EncryptionType: types.EncryptionTypeKms,
			KeyId:          aws.String(keyARN),
		})
		if err != nil {
			return fmt.Errorf("start encryption: %v", err)
		}

		resp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(sn)})
		if err != nil {
			return err
		}
		if resp.StreamDescription.EncryptionType != types.EncryptionTypeKms {
			return fmt.Errorf("expected KMS encryption, got %s", resp.StreamDescription.EncryptionType)
		}
		if resp.StreamDescription.KeyId == nil {
			return fmt.Errorf("KeyId is nil after encryption")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "UpdateStreamMode", func() error {
		sn := kinesisStream(ts, "mode")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		descResp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(sn)})
		if err != nil {
			return err
		}
		_, err = client.UpdateStreamMode(ctx, &kinesis.UpdateStreamModeInput{
			StreamARN: descResp.StreamDescription.StreamARN,
			StreamModeDetails: &types.StreamModeDetails{
				StreamMode: types.StreamModeOnDemand,
			},
		})
		if err != nil {
			return err
		}

		postResp, err := client.DescribeStreamSummary(ctx, &kinesis.DescribeStreamSummaryInput{
			StreamName: aws.String(sn),
		})
		if err != nil {
			return fmt.Errorf("describe summary after mode change: %v", err)
		}
		if postResp.StreamDescriptionSummary.StreamModeDetails == nil {
			return fmt.Errorf("StreamModeDetails is nil")
		}
		if postResp.StreamDescriptionSummary.StreamModeDetails.StreamMode != types.StreamModeOnDemand {
			return fmt.Errorf("StreamMode: expected ON_DEMAND, got %s", postResp.StreamDescriptionSummary.StreamModeDetails.StreamMode)
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DescribeAccountSettings", func() error {
		resp, err := client.DescribeAccountSettings(ctx, &kinesis.DescribeAccountSettingsInput{})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "UpdateAccountSettings", func() error {
		resp, err := client.UpdateAccountSettings(ctx, &kinesis.UpdateAccountSettingsInput{
			MinimumThroughputBillingCommitment: &types.MinimumThroughputBillingCommitmentInput{
				Status: types.MinimumThroughputBillingCommitmentInputStatusEnabled,
			},
		})
		if err != nil {
			return err
		}
		if resp.MinimumThroughputBillingCommitment == nil {
			return fmt.Errorf("MinimumThroughputBillingCommitment is nil in response")
		}
		if resp.MinimumThroughputBillingCommitment.Status != types.MinimumThroughputBillingCommitmentOutputStatusEnabled {
			return fmt.Errorf("response status: got %s, want ENABLED", resp.MinimumThroughputBillingCommitment.Status)
		}
		return nil
	}))

	// The persistence contract: Describe reports what Update answered —
	// the write round-trips instead of vanishing into an echo stub. The
	// setting resets to DISABLED so later runs describe a fresh account.
	results = append(results, r.RunTest("kinesis", "DescribeAccountSettings_RoundTrip", func() error {
		resp, err := client.DescribeAccountSettings(ctx, &kinesis.DescribeAccountSettingsInput{})
		if err != nil {
			return err
		}
		if resp.MinimumThroughputBillingCommitment == nil || resp.MinimumThroughputBillingCommitment.Status != types.MinimumThroughputBillingCommitmentOutputStatusEnabled {
			return fmt.Errorf("describe after update: got %v, want ENABLED", resp.MinimumThroughputBillingCommitment)
		}
		_, err = client.UpdateAccountSettings(ctx, &kinesis.UpdateAccountSettingsInput{
			MinimumThroughputBillingCommitment: &types.MinimumThroughputBillingCommitmentInput{
				Status: types.MinimumThroughputBillingCommitmentInputStatusDisabled,
			},
		})
		if err != nil {
			return fmt.Errorf("reset to DISABLED: %v", err)
		}
		return nil
	}))

	results = append(results, r.kinesisResourcePolicyTests(ctx, client, ts)...)
	results = append(results, r.kinesisAdvancedConfigTests(ctx, client, ts)...)

	_, _ = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})

	return results
}

func (r *TestRunner) kinesisResourcePolicyTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	policyStreamName := kinesisStream(ts, "policy")
	if _, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
		StreamName: aws.String(policyStreamName),
		ShardCount: aws.Int32(1),
	}); err != nil {
		return []TestResult{SetupFailResult("kinesis", "CreateStream (policy) setup failed: %v", err)}
	}
	policyStreamDesc, err := kinesisDescribeWhenReady(ctx, client, policyStreamName, 10*time.Second)
	if err != nil {
		return []TestResult{SetupFailResult("kinesis", "DescribeStream (policy) setup failed: %v", err)}
	}
	var policyStreamARN string
	if policyStreamDesc != nil && policyStreamDesc.StreamDescription != nil {
		policyStreamARN = aws.ToString(policyStreamDesc.StreamDescription.StreamARN)
	}

	if policyStreamARN != "" {
		testPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"kinesis:*","Resource":"*"}]}`

		results = append(results, r.RunTest("kinesis", "PutResourcePolicy", func() error {
			_, err := client.PutResourcePolicy(ctx, &kinesis.PutResourcePolicyInput{
				ResourceARN: aws.String(policyStreamARN),
				Policy:      aws.String(testPolicy),
			})
			if err != nil {
				return err
			}
			return nil
		}))

		results = append(results, r.RunTest("kinesis", "GetResourcePolicy", func() error {
			resp, err := client.GetResourcePolicy(ctx, &kinesis.GetResourcePolicyInput{
				ResourceARN: aws.String(policyStreamARN),
			})
			if err != nil {
				return err
			}
			policy := aws.ToString(resp.Policy)
			if policy == "" {
				return fmt.Errorf("expected non-empty policy")
			}
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(policy), &parsed); err != nil {
				return fmt.Errorf("policy is not valid JSON: %v", err)
			}
			return nil
		}))

		results = append(results, r.RunTest("kinesis", "DeleteResourcePolicy", func() error {
			_, err := client.DeleteResourcePolicy(ctx, &kinesis.DeleteResourcePolicyInput{
				ResourceARN: aws.String(policyStreamARN),
			})
			if err != nil {
				return err
			}
			// A resource whose policy was deleted answers the
			// operation's declared ResourceNotFoundException — the
			// missing-policy case, never an empty-string success.
			resp, err := client.GetResourcePolicy(ctx, &kinesis.GetResourcePolicyInput{
				ResourceARN: aws.String(policyStreamARN),
			})
			if err != nil {
				return expectAWSErrorCode(err, "ResourceNotFoundException")
			}
			return fmt.Errorf("expected ResourceNotFoundException after delete, got policy %q", aws.ToString(resp.Policy))
		}))
	} else {
		results = append(results, TestResult{Service: "kinesis", TestName: "PutResourcePolicy", Status: "SKIP", Error: "policyStreamARN not available"})
		results = append(results, TestResult{Service: "kinesis", TestName: "GetResourcePolicy", Status: "SKIP", Error: "policyStreamARN not available"})
		results = append(results, TestResult{Service: "kinesis", TestName: "DeleteResourcePolicy", Status: "SKIP", Error: "policyStreamARN not available"})
	}

	_, _ = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(policyStreamName)})
	return results
}

func (r *TestRunner) kinesisAdvancedConfigTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	maxRecordStreamName := kinesisStream(ts, "maxrec")
	if _, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
		StreamName: aws.String(maxRecordStreamName),
		ShardCount: aws.Int32(1),
	}); err != nil {
		return []TestResult{SetupFailResult("kinesis", "CreateStream (maxrec) setup failed: %v", err)}
	}
	maxRecordDesc, err := kinesisDescribeWhenReady(ctx, client, maxRecordStreamName, 10*time.Second)
	if err != nil {
		return []TestResult{SetupFailResult("kinesis", "DescribeStream (maxrec) setup failed: %v", err)}
	}

	if maxRecordDesc != nil && maxRecordDesc.StreamDescription != nil {
		results = append(results, r.RunTest("kinesis", "UpdateMaxRecordSize", func() error {
			_, err := client.UpdateMaxRecordSize(ctx, &kinesis.UpdateMaxRecordSizeInput{
				StreamARN:          maxRecordDesc.StreamDescription.StreamARN,
				MaxRecordSizeInKiB: aws.Int32(1024),
			})
			if err != nil {
				return err
			}
			return nil
		}))
	} else {
		results = append(results, TestResult{Service: "kinesis", TestName: "UpdateMaxRecordSize", Status: "SKIP", Error: "stream not available"})
	}
	_, _ = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(maxRecordStreamName)})

	warmStreamName := kinesisStream(ts, "warm")
	if _, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
		StreamName: aws.String(warmStreamName),
		ShardCount: aws.Int32(1),
	}); err != nil {
		return []TestResult{SetupFailResult("kinesis", "CreateStream (warm) setup failed: %v", err)}
	}
	warmDesc, err := kinesisDescribeWhenReady(ctx, client, warmStreamName, 10*time.Second)
	if err != nil {
		return []TestResult{SetupFailResult("kinesis", "DescribeStream (warm) setup failed: %v", err)}
	}

	// UpdateStreamWarmThroughput is documented for on-demand data streams:
	// switch the fixture's mode before the call under test.
	if _, err := client.UpdateStreamMode(ctx, &kinesis.UpdateStreamModeInput{
		StreamARN:         warmDesc.StreamDescription.StreamARN,
		StreamModeDetails: &types.StreamModeDetails{StreamMode: types.StreamModeOnDemand},
	}); err != nil {
		return []TestResult{SetupFailResult("kinesis", "UpdateStreamMode (warm) setup failed: %v", err)}
	}

	if warmDesc != nil && warmDesc.StreamDescription != nil {
		results = append(results, r.RunTest("kinesis", "UpdateStreamWarmThroughput", func() error {
			resp, err := client.UpdateStreamWarmThroughput(ctx, &kinesis.UpdateStreamWarmThroughputInput{
				StreamARN:           warmDesc.StreamDescription.StreamARN,
				WarmThroughputMiBps: aws.Int32(256),
			})
			if err != nil {
				return err
			}
			if resp.WarmThroughput == nil {
				return fmt.Errorf("WarmThroughput is nil")
			}
			if resp.WarmThroughput.TargetMiBps == nil {
				return fmt.Errorf("WarmThroughput.TargetMiBps is nil")
			}
			if got := *resp.WarmThroughput.TargetMiBps; got != 256 {
				return fmt.Errorf("WarmThroughput.TargetMiBps: expected 256, got %d", got)
			}
			return nil
		}))
	} else {
		results = append(results, TestResult{Service: "kinesis", TestName: "UpdateStreamWarmThroughput", Status: "SKIP", Error: "stream not available"})
	}
	_, _ = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(warmStreamName)})

	return results
}
