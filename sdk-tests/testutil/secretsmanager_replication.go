package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) runSecretsManagerReplicationTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "ReplicateSecretToRegions_Basic", func() error {
		name := tc.uniqueName("ReplicateTest")
		secretValue := "replicate-test-value"
		_, err := tc.createSecret(name, secretValue)
		if err != nil {
			return err
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

		descResp, err := tc.describeSecret(name)
		if err != nil {
			return err
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

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_PrimaryRegionOnReplica", func() error {
		name := tc.uniqueName("PrimRegion")
		_, err := tc.createSecret(name, "prim-region")
		if err != nil {
			return err
		}
		defer tc.forceDeleteSecret(name)

		if _, err := tc.client.ReplicateSecretToRegions(tc.ctx, &secretsmanager.ReplicateSecretToRegionsInput{
			SecretId: aws.String(name),
			AddReplicaRegions: []types.ReplicaRegionType{
				{Region: aws.String("us-west-2")},
			},
		}); err != nil {
			return fmt.Errorf("replicate: %v", err)
		}
		defer tc.client.RemoveRegionsFromReplication(tc.ctx, &secretsmanager.RemoveRegionsFromReplicationInput{
			SecretId:             aws.String(name),
			RemoveReplicaRegions: []string{"us-west-2"},
		})

		// The replica lives in the target region and its SecretListEntry
		// must carry PrimaryRegion pointing at the source region.
		cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   "us-west-2",
		})
		if err != nil {
			return fmt.Errorf("replica config: %v", err)
		}
		replicaClient := secretsmanager.NewFromConfig(cfg)

		resp, err := replicaClient.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{{Key: types.FilterNameStringTypeName, Values: []string{name}}},
		})
		if err != nil {
			return fmt.Errorf("replica list: %v", err)
		}
		found := false
		for _, entry := range resp.SecretList {
			if aws.ToString(entry.Name) == name {
				found = true
				if aws.ToString(entry.PrimaryRegion) != tc.region {
					return fmt.Errorf("replica entry PrimaryRegion = %q, want %q", aws.ToString(entry.PrimaryRegion), tc.region)
				}
			}
		}
		if !found {
			return fmt.Errorf("replica not listed in target region")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "RemoveRegionsFromReplication_ErrReplicaNotFound", func() error {
		name := tc.uniqueName("RmNonExist")
		_, err := tc.createSecret(name, "rm-nonexist")
		if err != nil {
			return err
		}
		defer tc.forceDeleteSecret(name)

		_, err = tc.client.RemoveRegionsFromReplication(tc.ctx, &secretsmanager.RemoveRegionsFromReplicationInput{
			SecretId:             aws.String(name),
			RemoveReplicaRegions: []string{"ap-south-1"},
		})
		if err == nil {
			return fmt.Errorf("removing a region with no replica should fail")
		}
		return expectAWSErrorCode(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("secretsmanager", "ReplicateSecretToRegions_DuplicateRegion", func() error {
		name := tc.uniqueName("ReplicateDup")
		_, err := tc.createSecret(name, "dup-repl")
		if err != nil {
			return err
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
		_, err := tc.createSecret(name, "stop-repl")
		if err != nil {
			return err
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

		descResp, err := tc.describeSecret(name)
		if err != nil {
			return err
		}
		if len(descResp.ReplicationStatus) > 0 {
			return fmt.Errorf("ReplicationStatus should be empty after stop")
		}
		return nil
	}))

	return results
}
