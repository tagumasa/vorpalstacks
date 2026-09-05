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
		_, err := tc.createSecret(name, "initial")
		if err != nil {
			return err
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

		getResp, err := tc.getSecretValue(name)
		if err != nil {
			return err
		}
		if getResp.SecretString == nil || *getResp.SecretString != newValue {
			return fmt.Errorf("value mismatch: got %q, want %q", aws.ToString(getResp.SecretString), newValue)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "PutSecretValue_MultipleVersions", func() error {
		name := tc.uniqueName("PutVerify")
		_, err := tc.createSecret(name, "v1")
		if err != nil {
			return err
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
		createResp, err := tc.createSecret(name, "v1")
		if err != nil {
			return err
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

	results = append(results, r.RunTest("secretsmanager", "GetSecretValue_UpdatesLastAccessedDate", func() error {
		name := tc.uniqueName("LastAccess")
		_, err := tc.createSecret(name, "accessed")
		if err != nil {
			return err
		}
		defer tc.forceDeleteSecret(name)

		before, err := tc.describeSecret(name)
		if err != nil {
			return fmt.Errorf("describe before: %w", err)
		}
		if before.LastAccessedDate != nil {
			return fmt.Errorf("LastAccessedDate set before any retrieval")
		}

		if _, err := tc.client.GetSecretValue(tc.ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(name),
		}); err != nil {
			return fmt.Errorf("get: %v", err)
		}

		after, err := tc.describeSecret(name)
		if err != nil {
			return fmt.Errorf("describe after: %w", err)
		}
		if after.LastAccessedDate == nil {
			return fmt.Errorf("LastAccessedDate not set after GetSecretValue")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "UpdateSecretVersionStage_Basic", func() error {
		name := tc.uniqueName("VersionStage")
		_, err := tc.createSecret(name, "v1")
		if err != nil {
			return err
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

		descResp, err := tc.describeSecret(name)
		if err != nil {
			return err
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

		descResp2, err := tc.describeSecret(name)
		if err != nil {
			return fmt.Errorf("describe after: %w", err)
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

	results = append(results, r.RunTest("secretsmanager", "UpdateSecretVersionStage_AddNewLabelToVersion", func() error {
		name := tc.uniqueName("StageAdd")
		_, err := tc.createSecret(name, "v1")
		if err != nil {
			return err
		}
		defer tc.forceDeleteSecret(name)

		putResp, err := tc.client.PutSecretValue(tc.ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(name),
			SecretString: aws.String("v2"),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}

		// A staging label not attached to any version is added by
		// MoveToVersionId alone (API_UpdateSecretVersionStage Example 1).
		if _, err := tc.client.UpdateSecretVersionStage(tc.ctx, &secretsmanager.UpdateSecretVersionStageInput{
			SecretId:        aws.String(name),
			VersionStage:    aws.String("EXAMPLESTAGING"),
			MoveToVersionId: putResp.VersionId,
		}); err != nil {
			return fmt.Errorf("add new label: %v", err)
		}

		descResp, err := tc.describeSecret(name)
		if err != nil {
			return err
		}
		stages, ok := descResp.VersionIdsToStages[*putResp.VersionId]
		if !ok {
			return fmt.Errorf("v2 missing from VersionIdsToStages")
		}
		hasLabel := false
		for _, s := range stages {
			if s == "EXAMPLESTAGING" {
				hasLabel = true
			}
		}
		if !hasLabel {
			return fmt.Errorf("new label not attached to v2")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "UpdateSecretVersionStage_MoveWithoutRemoveFromRejected", func() error {
		createResp, err := tc.createSecret(tc.uniqueName("StageOmit"), "v1")
		if err != nil {
			return err
		}
		name := *createResp.Name
		defer tc.forceDeleteSecret(name)

		putResp, err := tc.client.PutSecretValue(tc.ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(name),
			SecretString: aws.String("v2"),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}

		// Attach the label to v1 first (a pure add).
		if _, err := tc.client.UpdateSecretVersionStage(tc.ctx, &secretsmanager.UpdateSecretVersionStageInput{
			SecretId:        aws.String(name),
			VersionStage:    aws.String("MOVINGLABEL"),
			MoveToVersionId: createResp.VersionId,
		}); err != nil {
			return fmt.Errorf("attach label to v1: %v", err)
		}

		// Moving a label already attached to another version without
		// RemoveFromVersionId must fail: "If the label is attached and
		// you either do not specify this parameter ... then the operation
		// fails."
		_, err = tc.client.UpdateSecretVersionStage(tc.ctx, &secretsmanager.UpdateSecretVersionStageInput{
			SecretId:        aws.String(name),
			VersionStage:    aws.String("MOVINGLABEL"),
			MoveToVersionId: putResp.VersionId,
		})
		if err == nil {
			return fmt.Errorf("move without RemoveFromVersionId should be rejected")
		}
		return expectAWSErrorCode(err, "InvalidParameterException")
	}))

	results = append(results, r.RunTest("secretsmanager", "UpdateSecretVersionStage_MoveRemoveFromMismatchRejected", func() error {
		createResp, err := tc.createSecret(tc.uniqueName("StageMis"), "v1")
		if err != nil {
			return err
		}
		name := *createResp.Name
		defer tc.forceDeleteSecret(name)

		putResp, err := tc.client.PutSecretValue(tc.ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(name),
			SecretString: aws.String("v2"),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}
		thirdResp, err := tc.client.PutSecretValue(tc.ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(name),
			SecretString: aws.String("v3"),
		})
		if err != nil {
			return fmt.Errorf("put v3: %v", err)
		}

		// Attach the label to v1 (a pure add).
		if _, err := tc.client.UpdateSecretVersionStage(tc.ctx, &secretsmanager.UpdateSecretVersionStageInput{
			SecretId:        aws.String(name),
			VersionStage:    aws.String("MISMATCHLABEL"),
			MoveToVersionId: createResp.VersionId,
		}); err != nil {
			return fmt.Errorf("attach label to v1: %v", err)
		}

		// A RemoveFromVersionId that does not match the version actually
		// holding the label must fail: "... or the version ID does not
		// match, then the operation fails."
		_, err = tc.client.UpdateSecretVersionStage(tc.ctx, &secretsmanager.UpdateSecretVersionStageInput{
			SecretId:            aws.String(name),
			VersionStage:        aws.String("MISMATCHLABEL"),
			MoveToVersionId:     putResp.VersionId,
			RemoveFromVersionId: thirdResp.VersionId,
		})
		if err == nil {
			return fmt.Errorf("move with mismatched RemoveFromVersionId should be rejected")
		}
		if e := expectAWSErrorCode(err, "InvalidRequestException"); e != nil {
			return e
		}

		// The matching move succeeds and relocates the label.
		if _, err := tc.client.UpdateSecretVersionStage(tc.ctx, &secretsmanager.UpdateSecretVersionStageInput{
			SecretId:            aws.String(name),
			VersionStage:        aws.String("MISMATCHLABEL"),
			MoveToVersionId:     putResp.VersionId,
			RemoveFromVersionId: createResp.VersionId,
		}); err != nil {
			return fmt.Errorf("move with matching RemoveFromVersionId: %v", err)
		}

		descResp, err := tc.describeSecret(name)
		if err != nil {
			return err
		}
		labelOn := func(versionId *string) bool {
			stages, ok := descResp.VersionIdsToStages[*versionId]
			if !ok {
				return false
			}
			for _, s := range stages {
				if s == "MISMATCHLABEL" {
					return true
				}
			}
			return false
		}
		if !labelOn(putResp.VersionId) {
			return fmt.Errorf("label not moved to v2")
		}
		if labelOn(createResp.VersionId) {
			return fmt.Errorf("label still attached to v1 after move")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "PutSecretValue_AWSPENDING_WithoutRotationToken", func() error {
		name := tc.uniqueName("Pending")
		_, err := tc.createSecret(name, "initial")
		if err != nil {
			return err
		}
		defer tc.forceDeleteSecret(name)

		// The standard same-account rotation Lambda template calls
		// PutSecretValue with VersionStages=["AWSPENDING"] and a
		// ClientRequestToken (used as the version ID) but WITHOUT a
		// RotationToken parameter. AWS does not gate AWSPENDING on
		// RotationToken presence — it is only required for cross-account
		// rotation. This test verifies that the server accepts the call.
		token := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
		putResp, err := tc.client.PutSecretValue(tc.ctx, &secretsmanager.PutSecretValueInput{
			SecretId:           aws.String(name),
			SecretString:       aws.String("pending-value"),
			ClientRequestToken: aws.String(token),
			VersionStages:      []string{"AWSPENDING"},
		})
		if err != nil {
			return fmt.Errorf("put with AWSPENDING (no RotationToken): %v", err)
		}
		if putResp.VersionId == nil || *putResp.VersionId != token {
			return fmt.Errorf("VersionId mismatch: expected %s", token)
		}

		verResp, err := tc.client.ListSecretVersionIds(tc.ctx, &secretsmanager.ListSecretVersionIdsInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("list versions: %v", err)
		}
		foundPending := false
		for _, v := range verResp.Versions {
			for _, s := range v.VersionStages {
				if s == "AWSPENDING" {
					foundPending = true
				}
			}
		}
		if !foundPending {
			return fmt.Errorf("AWSPENDING staging label not found in any version")
		}
		return nil
	}))

	return results
}
