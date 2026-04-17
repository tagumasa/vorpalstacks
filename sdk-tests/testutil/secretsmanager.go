package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunSecretsManagerTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "secretsmanager",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := secretsmanager.NewFromConfig(cfg)
	ctx := context.Background()

	secretName := fmt.Sprintf("TestSecret-%d", time.Now().UnixNano())
	secretValue := "my-secret-value"

	results = append(results, r.RunTest("secretsmanager", "CreateSecret", func() error {
		resp, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secretName),
			SecretString: aws.String(secretValue),
		})
		if err != nil {
			return err
		}
		if resp.ARN == nil {
			return fmt.Errorf("secret ARN is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DescribeSecret", func() error {
		resp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secretName),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil {
			return fmt.Errorf("secret name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetSecretValue", func() error {
		resp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		})
		if err != nil {
			return err
		}
		if resp.SecretString == nil && resp.SecretBinary == nil {
			return fmt.Errorf("secret value is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets", func() error {
		resp, err := client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{})
		if err != nil {
			return err
		}
		if resp.SecretList == nil {
			return fmt.Errorf("secrets list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "UpdateSecret", func() error {
		resp, err := client.UpdateSecret(ctx, &secretsmanager.UpdateSecretInput{
			SecretId:     aws.String(secretName),
			SecretString: aws.String("updated-secret-value"),
		})
		if err != nil {
			return err
		}
		if resp.ARN == nil {
			return fmt.Errorf("secret ARN is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "TagResource", func() error {
		resp, err := client.TagResource(ctx, &secretsmanager.TagResourceInput{
			SecretId: aws.String(secretName),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Project"), Value: aws.String("sdk-tests")},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecretVersionIds", func() error {
		resp, err := client.ListSecretVersionIds(ctx, &secretsmanager.ListSecretVersionIdsInput{
			SecretId: aws.String(secretName),
		})
		if err != nil {
			return err
		}
		if resp.Versions == nil {
			return fmt.Errorf("versions list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DeleteSecret", func() error {
		resp, err := client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secretName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("secretsmanager", "GetSecretValue_NonExistent", func() error {
		_, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String("arn:aws:secretsmanager:us-east-1:000000000000:secret:nonexistent-secret-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DescribeSecret_NonExistent", func() error {
		_, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String("arn:aws:secretsmanager:us-east-1:000000000000:secret:nonexistent-secret-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "CreateSecret_Duplicate", func() error {
		dupName := fmt.Sprintf("DupSecret-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(dupName),
			SecretString: aws.String("initial-value"),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(dupName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		_, err = client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(dupName),
			SecretString: aws.String("duplicate-value"),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate secret name")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetSecretValue_ContentVerify", func() error {
		secName := fmt.Sprintf("VerifySecret-%d", time.Now().UnixNano())
		secValue := "my-verified-secret-123"
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String(secValue),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.SecretString == nil || *resp.SecretString != secValue {
			return fmt.Errorf("secret value mismatch: got %q, want %q", aws.ToString(resp.SecretString), secValue)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "UpdateSecret_ContentVerify", func() error {
		secName := fmt.Sprintf("UpdateVerify-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("original-value"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		updatedValue := "updated-secret-value-456"
		_, err = client.UpdateSecret(ctx, &secretsmanager.UpdateSecretInput{
			SecretId:     aws.String(secName),
			SecretString: aws.String(updatedValue),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		resp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.SecretString == nil || *resp.SecretString != updatedValue {
			return fmt.Errorf("secret value not updated: got %q, want %q", aws.ToString(resp.SecretString), updatedValue)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_ContainsCreated", func() error {
		secName := fmt.Sprintf("ListVerify-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("list-test"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, s := range resp.SecretList {
			if s.Name != nil && *s.Name == secName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created secret not found in ListSecrets")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DeleteSecret_NonExistent", func() error {
		_, err := client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String("arn:aws:secretsmanager:us-east-1:000000000000:secret:nonexistent-xyz"),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// === PUT SECRET VALUE TESTS ===

	results = append(results, r.RunTest("secretsmanager", "PutSecretValue_Basic", func() error {
		secName := fmt.Sprintf("PutValue-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("initial"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(secName),
			SecretString: aws.String("new-value"),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		if resp.VersionId == nil {
			return fmt.Errorf("version ID is nil")
		}

		getResp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.SecretString == nil || *getResp.SecretString != "new-value" {
			return fmt.Errorf("value mismatch: got %q", aws.ToString(getResp.SecretString))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "PutSecretValue_ContentVerify", func() error {
		secName := fmt.Sprintf("PutVerify-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("v1"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		_, err = client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(secName),
			SecretString: aws.String("v2"),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}
		_, err = client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(secName),
			SecretString: aws.String("v3"),
		})
		if err != nil {
			return fmt.Errorf("put v3: %v", err)
		}

		verResp, err := client.ListSecretVersionIds(ctx, &secretsmanager.ListSecretVersionIdsInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("list versions: %v", err)
		}
		if len(verResp.Versions) != 3 {
			return fmt.Errorf("expected 3 versions, got %d", len(verResp.Versions))
		}
		return nil
	}))

	// === RESTORE SECRET TESTS ===

	results = append(results, r.RunTest("secretsmanager", "RestoreSecret_Basic", func() error {
		secName := fmt.Sprintf("RestoreSec-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("restore-test"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}

		_, err = client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		listResp, err := client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		for _, s := range listResp.SecretList {
			if s.Name != nil && *s.Name == secName {
				return fmt.Errorf("soft-deleted secret should not appear in ListSecrets")
			}
		}

		_, err = client.RestoreSecret(ctx, &secretsmanager.RestoreSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("restore: %v", err)
		}

		getResp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("get after restore: %v", err)
		}
		if getResp.SecretString == nil || *getResp.SecretString != "restore-test" {
			return fmt.Errorf("value mismatch after restore")
		}

		client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "RestoreSecret_NonExistent", func() error {
		_, err := client.RestoreSecret(ctx, &secretsmanager.RestoreSecretInput{
			SecretId: aws.String("nonexistent-restore-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// === GET RANDOM PASSWORD TESTS ===

	results = append(results, r.RunTest("secretsmanager", "GetRandomPassword_Basic", func() error {
		resp, err := client.GetRandomPassword(ctx, &secretsmanager.GetRandomPasswordInput{})
		if err != nil {
			return err
		}
		if resp.RandomPassword == nil || len(*resp.RandomPassword) != 32 {
			return fmt.Errorf("expected default password length 32, got %s", aws.ToString(resp.RandomPassword))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetRandomPassword_CustomLength", func() error {
		resp, err := client.GetRandomPassword(ctx, &secretsmanager.GetRandomPasswordInput{
			PasswordLength: aws.Int64(16),
		})
		if err != nil {
			return err
		}
		if resp.RandomPassword == nil || len(*resp.RandomPassword) != 16 {
			return fmt.Errorf("expected password length 16, got %d", len(*resp.RandomPassword))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetRandomPassword_ExcludeCharacters", func() error {
		resp, err := client.GetRandomPassword(ctx, &secretsmanager.GetRandomPasswordInput{
			PasswordLength:     aws.Int64(50),
			ExcludeCharacters:  aws.String("abcdefABCDEF0123456789"),
			ExcludePunctuation: aws.Bool(true),
			IncludeSpace:       aws.Bool(false),
		})
		if err != nil {
			return err
		}
		for _, c := range *resp.RandomPassword {
			if c >= 'a' && c <= 'f' {
				return fmt.Errorf("found excluded lowercase char: %c", c)
			}
			if c >= 'A' && c <= 'F' {
				return fmt.Errorf("found excluded uppercase char: %c", c)
			}
			if c >= '0' && c <= '5' {
				return fmt.Errorf("found excluded digit: %c", c)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetRandomPassword_RequireEachIncludedType", func() error {
		resp, err := client.GetRandomPassword(ctx, &secretsmanager.GetRandomPasswordInput{
			PasswordLength:          aws.Int64(20),
			RequireEachIncludedType: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		hasLower, hasUpper, hasDigit, hasPunct := false, false, false, false
		for _, c := range *resp.RandomPassword {
			if c >= 'a' && c <= 'z' {
				hasLower = true
			}
			if c >= 'A' && c <= 'Z' {
				hasUpper = true
			}
			if c >= '0' && c <= '9' {
				hasDigit = true
			}
			if c >= '!' && c <= '/' || c >= ':' && c <= '@' || c >= '[' && c <= '`' || c >= '{' && c <= '~' {
				hasPunct = true
			}
		}
		if !hasLower || !hasUpper || !hasDigit || !hasPunct {
			return fmt.Errorf("missing required types: lower=%v upper=%v digit=%v punct=%v", hasLower, hasUpper, hasDigit, hasPunct)
		}
		return nil
	}))

	// === UPDATE SECRET VERSION STAGE TESTS ===

	results = append(results, r.RunTest("secretsmanager", "UpdateSecretVersionStage_Basic", func() error {
		secName := fmt.Sprintf("VersionStage-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("v1"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		putResp, err := client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(secName),
			SecretString: aws.String("v2"),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}
		v2VersionId := putResp.VersionId

		descResp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secName),
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

		_, err = client.UpdateSecretVersionStage(ctx, &secretsmanager.UpdateSecretVersionStageInput{
			SecretId:        aws.String(secName),
			VersionStage:    aws.String("AWSCURRENT"),
			MoveToVersionId: v2VersionId,
		})
		if err != nil {
			return fmt.Errorf("update version stage: %v", err)
		}

		descResp2, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secName),
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

	// === TAG OPERATIONS TESTS ===

	results = append(results, r.RunTest("secretsmanager", "UntagResource_Basic", func() error {
		secName := fmt.Sprintf("UntagTest-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name: aws.String(secName),
			Tags: []types.Tag{
				{Key: aws.String("env"), Value: aws.String("test")},
				{Key: aws.String("team"), Value: aws.String("dev")},
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		_, err = client.UntagResource(ctx, &secretsmanager.UntagResourceInput{
			SecretId: aws.String(secName),
			TagKeys:  []string{"env"},
		})
		if err != nil {
			return fmt.Errorf("untag: %v", err)
		}

		descResp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		for _, t := range descResp.Tags {
			if t.Key != nil && *t.Key == "env" {
				return fmt.Errorf("env tag should have been removed")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DescribeSecret_Tags", func() error {
		secName := fmt.Sprintf("ListTags-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name: aws.String(secName),
			Tags: []types.Tag{
				{Key: aws.String("key1"), Value: aws.String("val1")},
				{Key: aws.String("key2"), Value: aws.String("val2")},
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Tags) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(resp.Tags))
		}
		tagMap := make(map[string]string)
		for _, t := range resp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["key1"] != "val1" || tagMap["key2"] != "val2" {
			return fmt.Errorf("tag content mismatch: %v", tagMap)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "CreateSecret_WithTags", func() error {
		secName := fmt.Sprintf("TagCreate-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("tagged"),
			Tags: []types.Tag{
				{Key: aws.String("Owner"), Value: aws.String("test-suite")},
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		found := false
		for _, t := range resp.Tags {
			if t.Key != nil && *t.Key == "Owner" && t.Value != nil && *t.Value == "test-suite" {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("owner tag not found in DescribeSecret")
		}
		return nil
	}))

	// === RESOURCE POLICY TESTS ===

	results = append(results, r.RunTest("secretsmanager", "PutResourcePolicy_Basic", func() error {
		secName := fmt.Sprintf("PolicyTest-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("policy-test"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"secretsmanager:GetSecretValue","Resource":"*"}]}`
		_, err = client.PutResourcePolicy(ctx, &secretsmanager.PutResourcePolicyInput{
			SecretId:       aws.String(secName),
			ResourcePolicy: aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("put policy: %v", err)
		}

		getResp, err := client.GetResourcePolicy(ctx, &secretsmanager.GetResourcePolicyInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("get policy: %v", err)
		}
		if getResp.ResourcePolicy == nil || *getResp.ResourcePolicy != policy {
			return fmt.Errorf("policy mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DeleteResourcePolicy_Basic", func() error {
		secName := fmt.Sprintf("DelPolicy-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("del-policy"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		policy := `{"Version":"2012-10-17","Statement":[]}`
		_, err = client.PutResourcePolicy(ctx, &secretsmanager.PutResourcePolicyInput{
			SecretId:       aws.String(secName),
			ResourcePolicy: aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("put policy: %v", err)
		}

		_, err = client.DeleteResourcePolicy(ctx, &secretsmanager.DeleteResourcePolicyInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("delete policy: %v", err)
		}

		getResp, err := client.GetResourcePolicy(ctx, &secretsmanager.GetResourcePolicyInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("get policy: %v", err)
		}
		if getResp.ResourcePolicy != nil {
			return fmt.Errorf("policy should be nil after deletion")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ValidateResourcePolicy_Valid", func() error {
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"*","Resource":"*"}]}`
		resp, err := client.ValidateResourcePolicy(ctx, &secretsmanager.ValidateResourcePolicyInput{
			ResourcePolicy: aws.String(policy),
		})
		if err != nil {
			return err
		}
		if !resp.PolicyValidationPassed {
			return fmt.Errorf("expected policy validation to pass")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ValidateResourcePolicy_Invalid", func() error {
		resp, err := client.ValidateResourcePolicy(ctx, &secretsmanager.ValidateResourcePolicyInput{
			ResourcePolicy: aws.String("not valid json {"),
		})
		if err != nil {
			return err
		}
		if resp.PolicyValidationPassed {
			return fmt.Errorf("expected policy validation to fail for invalid JSON")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetResourcePolicy_NonExistent", func() error {
		_, err := client.GetResourcePolicy(ctx, &secretsmanager.GetResourcePolicyInput{
			SecretId: aws.String("nonexistent-policy-secret"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "PutResourcePolicy_NonExistent", func() error {
		_, err := client.PutResourcePolicy(ctx, &secretsmanager.PutResourcePolicyInput{
			SecretId:       aws.String("nonexistent-policy-secret"),
			ResourcePolicy: aws.String(`{"Version":"2012-10-17"}`),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// === BATCH GET SECRET VALUE TESTS ===

	results = append(results, r.RunTest("secretsmanager", "BatchGetSecretValue_Basic", func() error {
		sec1 := fmt.Sprintf("Batch1-%d", time.Now().UnixNano())
		sec2 := fmt.Sprintf("Batch2-%d", time.Now().UnixNano()+1)
		for _, name := range []string{sec1, sec2} {
			_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
				Name:         aws.String(name),
				SecretString: aws.String("batch-value-" + name),
			})
			if err != nil {
				return fmt.Errorf("create %s: %v", name, err)
			}
			defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
				SecretId:                   aws.String(name),
				ForceDeleteWithoutRecovery: aws.Bool(true),
			})
		}

		resp, err := client.BatchGetSecretValue(ctx, &secretsmanager.BatchGetSecretValueInput{
			SecretIdList: []string{sec1, sec2},
		})
		if err != nil {
			return fmt.Errorf("batch get: %v", err)
		}
		if len(resp.SecretValues) != 2 {
			return fmt.Errorf("expected 2 secret values, got %d", len(resp.SecretValues))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "BatchGetSecretValue_NonExistent", func() error {
		secName := fmt.Sprintf("BatchNE-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("exists"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.BatchGetSecretValue(ctx, &secretsmanager.BatchGetSecretValueInput{
			SecretIdList: []string{secName, "nonexistent-batch-secret"},
		})
		if err != nil {
			return fmt.Errorf("batch get: %v", err)
		}
		if len(resp.SecretValues) != 1 {
			return fmt.Errorf("expected 1 secret value, got %d", len(resp.SecretValues))
		}
		if len(resp.Errors) != 1 {
			return fmt.Errorf("expected 1 error, got %d", len(resp.Errors))
		}
		return nil
	}))

	// === ROTATION TESTS ===

	results = append(results, r.RunTest("secretsmanager", "RotateSecret_Basic", func() error {
		secName := fmt.Sprintf("RotateTest-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("rotate-me"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.RotateSecret(ctx, &secretsmanager.RotateSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("rotate: %v", err)
		}
		if resp.VersionId == nil {
			return fmt.Errorf("version ID is nil after rotation")
		}

		descResp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if descResp.LastRotatedDate == nil {
			return fmt.Errorf("LastRotatedDate should be set after rotation")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "CancelRotateSecret_Basic", func() error {
		secName := fmt.Sprintf("CancelRot-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("cancel-rotate"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		_, err = client.RotateSecret(ctx, &secretsmanager.RotateSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("rotate: %v", err)
		}

		_, err = client.CancelRotateSecret(ctx, &secretsmanager.CancelRotateSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("cancel rotate: %v", err)
		}

		descResp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if descResp.RotationEnabled != nil && *descResp.RotationEnabled {
			return fmt.Errorf("rotation should be disabled after cancel")
		}
		return nil
	}))

	// === CREATE SECRET WITH DESCRIPTION ===

	results = append(results, r.RunTest("secretsmanager", "CreateSecret_WithDescription", func() error {
		secName := fmt.Sprintf("DescTest-%d", time.Now().UnixNano())
		desc := "My test secret description"
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("desc-value"),
			Description:  aws.String(desc),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if resp.Description == nil || *resp.Description != desc {
			return fmt.Errorf("description mismatch: got %q, want %q", aws.ToString(resp.Description), desc)
		}
		return nil
	}))

	// === LIST SECRETS FILTER TESTS ===

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_Filters", func() error {
		prefix := fmt.Sprintf("FilterTest-%d-", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(prefix + "alpha"),
			SecretString: aws.String("a"),
		})
		if err != nil {
			return fmt.Errorf("create alpha: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(prefix + "alpha"),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
		_, err = client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(prefix + "beta"),
			SecretString: aws.String("b"),
		})
		if err != nil {
			return fmt.Errorf("create beta: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(prefix + "beta"),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{
				{Key: types.FilterNameStringTypeName, Values: []string{prefix + "alpha"}},
			},
		})
		if err != nil {
			return fmt.Errorf("list with filter: %v", err)
		}
		if len(resp.SecretList) != 1 {
			return fmt.Errorf("expected 1 secret, got %d", len(resp.SecretList))
		}
		if resp.SecretList[0].Name == nil || *resp.SecretList[0].Name != prefix+"alpha" {
			return fmt.Errorf("wrong secret returned")
		}
		return nil
	}))

	// === DESCRIBE SECRET CONTENT VERIFY ===

	results = append(results, r.RunTest("secretsmanager", "DescribeSecret_ContentVerify", func() error {
		secName := fmt.Sprintf("DescVerify-%d", time.Now().UnixNano())
		createResp, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("desc-content"),
			Description:  aws.String("test description"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if resp.Name == nil || *resp.Name != secName {
			return fmt.Errorf("name mismatch")
		}
		if resp.ARN == nil || *resp.ARN != *createResp.ARN {
			return fmt.Errorf("ARN mismatch")
		}
		if resp.Description == nil || *resp.Description != "test description" {
			return fmt.Errorf("description mismatch")
		}
		if resp.CreatedDate == nil {
			return fmt.Errorf("CreatedDate is nil")
		}
		return nil
	}))

	// === SECRET BINARY TESTS ===

	results = append(results, r.RunTest("secretsmanager", "CreateSecret_Binary", func() error {
		secName := fmt.Sprintf("BinaryTest-%d", time.Now().UnixNano())
		binaryData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretBinary: binaryData,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.SecretBinary == nil {
			return fmt.Errorf("SecretBinary is nil")
		}
		if string(resp.SecretBinary) != string(binaryData) {
			return fmt.Errorf("binary data mismatch")
		}
		return nil
	}))

	// === UPDATE SECRET DESCRIPTION CLEAR TEST ===

	results = append(results, r.RunTest("secretsmanager", "UpdateSecret_ClearDescription", func() error {
		secName := fmt.Sprintf("ClearDesc-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("clear-desc"),
			Description:  aws.String("initial description"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		_, err = client.UpdateSecret(ctx, &secretsmanager.UpdateSecretInput{
			SecretId:    aws.String(secName),
			Description: aws.String(""),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		resp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if resp.Description != nil && *resp.Description != "" {
			return fmt.Errorf("description should be cleared, got %q", aws.ToString(resp.Description))
		}
		return nil
	}))

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgSecrets []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagSecret-%s-%d", pgTs, i)
			_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
				Name:         aws.String(name),
				SecretString: aws.String("pagval"),
			})
			if err != nil {
				for _, sn := range pgSecrets {
					client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
						SecretId:                   aws.String(sn),
						ForceDeleteWithoutRecovery: aws.Bool(true),
					})
				}
				return fmt.Errorf("create secret %s: %v", name, err)
			}
			pgSecrets = append(pgSecrets, name)
		}

		var allSecrets []string
		var nextToken *string
		for {
			resp, err := client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, sn := range pgSecrets {
					client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
						SecretId:                   aws.String(sn),
						ForceDeleteWithoutRecovery: aws.Bool(true),
					})
				}
				return fmt.Errorf("list secrets page: %v", err)
			}
			for _, s := range resp.SecretList {
				if s.Name != nil && strings.Contains(*s.Name, "PagSecret-"+pgTs) {
					allSecrets = append(allSecrets, *s.Name)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, sn := range pgSecrets {
			client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
				SecretId:                   aws.String(sn),
				ForceDeleteWithoutRecovery: aws.Bool(true),
			})
		}
		if len(allSecrets) != 5 {
			return fmt.Errorf("expected 5 paginated secrets, got %d", len(allSecrets))
		}
		return nil
	}))

	return results
}
