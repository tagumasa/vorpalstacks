package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func (r *TestRunner) runSecretsManagerValueTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "PutSecretValue_Basic", func() error {
		name := tc.uniqueName("PutValue")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("initial"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		newValue := "new-value"
		resp, err := tc.client.PutSecretValue(tc.ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(name),
			SecretString: aws.String(newValue),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		if resp.VersionId == nil {
			return fmt.Errorf("VersionId is nil")
		}
		if resp.Name == nil || *resp.Name != name {
			return fmt.Errorf("name mismatch")
		}

		getResp, err := tc.client.GetSecretValue(tc.ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.SecretString == nil || *getResp.SecretString != newValue {
			return fmt.Errorf("value mismatch: got %q, want %q", aws.ToString(getResp.SecretString), newValue)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "PutSecretValue_MultipleVersions", func() error {
		name := tc.uniqueName("PutVerify")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("v1"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		putResp2, err := tc.client.PutSecretValue(tc.ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(name),
			SecretString: aws.String("v2"),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}
		putResp3, err := tc.client.PutSecretValue(tc.ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(name),
			SecretString: aws.String("v3"),
		})
		if err != nil {
			return fmt.Errorf("put v3: %v", err)
		}
		if putResp2.VersionId == nil || putResp3.VersionId == nil {
			return fmt.Errorf("VersionId is nil")
		}
		if *putResp2.VersionId == *putResp3.VersionId {
			return fmt.Errorf("version IDs should be unique")
		}

		verResp, err := tc.client.ListSecretVersionIds(tc.ctx, &secretsmanager.ListSecretVersionIdsInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("list versions: %v", err)
		}
		if len(verResp.Versions) != 3 {
			return fmt.Errorf("expected 3 versions, got %d", len(verResp.Versions))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecretVersionIds", func() error {
		name := tc.uniqueName("VerList")
		createResp, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("v1"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		resp, err := tc.client.ListSecretVersionIds(tc.ctx, &secretsmanager.ListSecretVersionIdsInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return err
		}
		if resp.Versions == nil {
			return fmt.Errorf("Versions is nil")
		}
		if len(resp.Versions) < 1 {
			return fmt.Errorf("expected at least 1 version, got %d", len(resp.Versions))
		}
		if resp.Name == nil || *resp.Name != name {
			return fmt.Errorf("name mismatch")
		}
		if resp.ARN == nil {
			return fmt.Errorf("ARN is nil")
		}

		initialVersion := resp.Versions[0]
		foundInitial := false
		for _, stage := range initialVersion.VersionStages {
			if stage == "AWSCURRENT" {
				foundInitial = true
			}
		}
		if !foundInitial {
			return fmt.Errorf("initial version should have AWSCURRENT stage")
		}
		if createResp.VersionId != nil && initialVersion.VersionId != nil && *initialVersion.VersionId != *createResp.VersionId {
			return fmt.Errorf("version ID mismatch with CreateSecret response")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "UpdateSecretVersionStage_Basic", func() error {
		name := tc.uniqueName("VersionStage")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("v1"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		putResp, err := tc.client.PutSecretValue(tc.ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(name),
			SecretString: aws.String("v2"),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}
		v2VersionId := putResp.VersionId

		descResp, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if descResp.VersionIdsToStages == nil {
			return fmt.Errorf("VersionIdsToStages is nil")
		}
		if stages, ok := descResp.VersionIdsToStages[*v2VersionId]; !ok {
			return fmt.Errorf("v2 not in VersionIdsToStages")
		} else {
			hasCurrent := false
			for _, s := range stages {
				if s == "AWSCURRENT" {
					hasCurrent = true
				}
			}
			if !hasCurrent {
				return fmt.Errorf("v2 should have AWSCURRENT stage")
			}
		}

		_, err = tc.client.UpdateSecretVersionStage(tc.ctx, &secretsmanager.UpdateSecretVersionStageInput{
			SecretId:        aws.String(name),
			VersionStage:    aws.String("AWSCURRENT"),
			MoveToVersionId: v2VersionId,
		})
		if err != nil {
			return fmt.Errorf("update version stage: %v", err)
		}

		descResp2, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe after: %v", err)
		}
		if descResp2.VersionIdsToStages == nil {
			return fmt.Errorf("VersionIdsToStages is nil after update")
		}
		stages, ok := descResp2.VersionIdsToStages[*v2VersionId]
		if !ok {
			return fmt.Errorf("v2 not in VersionIdsToStages after update")
		}
		hasCurrent := false
		for _, s := range stages {
			if s == "AWSCURRENT" {
				hasCurrent = true
			}
		}
		if !hasCurrent {
			return fmt.Errorf("v2 should still have AWSCURRENT")
		}
		return nil
	}))

	return results
}
