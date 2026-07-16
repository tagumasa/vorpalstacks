package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

func (r *TestRunner) runSecretsManagerReplicationTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "ReplicateSecretToRegions_Basic", func() error {
		name := tc.uniqueName("ReplicateTest")
		secretValue := "replicate-test-value"
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String(secretValue),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		resp, err := tc.client.ReplicateSecretToRegions(tc.ctx, &secretsmanager.ReplicateSecretToRegionsInput{
			SecretId: aws.String(name),
			AddReplicaRegions: []types.ReplicaRegionType{
				{Region: aws.String("us-west-2")},
			},
		})
		if err != nil {
			return fmt.Errorf("replicate: %v", err)
		}
		if resp.ARN == nil {
			return fmt.Errorf("ARN is nil")
		}
		if len(resp.ReplicationStatus) == 0 {
			return fmt.Errorf("ReplicationStatus is empty")
		}

		found := false
		for _, rs := range resp.ReplicationStatus {
			if rs.Region != nil && *rs.Region == "us-west-2" {
				found = true
				if rs.Status != types.StatusTypeInSync {
					return fmt.Errorf("replica status mismatch: got %q, want %q", rs.Status, types.StatusTypeInSync)
				}
			}
		}
		if !found {
			return fmt.Errorf("us-west-2 not found in ReplicationStatus")
		}

		descResp, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.ReplicationStatus) == 0 {
			return fmt.Errorf("ReplicationStatus missing from DescribeSecret")
		}

		rmResp, err := tc.client.RemoveRegionsFromReplication(tc.ctx, &secretsmanager.RemoveRegionsFromReplicationInput{
			SecretId:             aws.String(name),
			RemoveReplicaRegions: []string{"us-west-2"},
		})
		if err != nil {
			return fmt.Errorf("remove regions: %v", err)
		}
		if rmResp.ARN == nil {
			return fmt.Errorf("ARN is nil after remove")
		}
		for _, rs := range rmResp.ReplicationStatus {
			if rs.Region != nil && *rs.Region == "us-west-2" {
				return fmt.Errorf("us-west-2 should have been removed from ReplicationStatus")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ReplicateSecretToRegions_DuplicateRegion", func() error {
		name := tc.uniqueName("ReplicateDup")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("dup-repl"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		_, err = tc.client.ReplicateSecretToRegions(tc.ctx, &secretsmanager.ReplicateSecretToRegionsInput{
			SecretId: aws.String(name),
			AddReplicaRegions: []types.ReplicaRegionType{
				{Region: aws.String("eu-west-1")},
			},
		})
		if err != nil {
			return fmt.Errorf("first replicate: %v", err)
		}

		_, err = tc.client.ReplicateSecretToRegions(tc.ctx, &secretsmanager.ReplicateSecretToRegionsInput{
			SecretId: aws.String(name),
			AddReplicaRegions: []types.ReplicaRegionType{
				{Region: aws.String("eu-west-1")},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate replication region")
		}

		_, _ = tc.client.RemoveRegionsFromReplication(tc.ctx, &secretsmanager.RemoveRegionsFromReplicationInput{
			SecretId:             aws.String(name),
			RemoveReplicaRegions: []string{"eu-west-1"},
		})
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "StopReplicationToReplica", func() error {
		name := tc.uniqueName("StopRepl")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("stop-repl"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		_, err = tc.client.ReplicateSecretToRegions(tc.ctx, &secretsmanager.ReplicateSecretToRegionsInput{
			SecretId: aws.String(name),
			AddReplicaRegions: []types.ReplicaRegionType{
				{Region: aws.String("ap-northeast-1")},
			},
		})
		if err != nil {
			return fmt.Errorf("replicate: %v", err)
		}

		stopResp, err := tc.client.StopReplicationToReplica(tc.ctx, &secretsmanager.StopReplicationToReplicaInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("stop replication: %v", err)
		}
		if stopResp.ARN == nil {
			return fmt.Errorf("ARN is nil")
		}

		descResp, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.ReplicationStatus) > 0 {
			return fmt.Errorf("ReplicationStatus should be empty after stop")
		}
		return nil
	}))

	return results
}
