package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

func pgTimestamp() int64 {
	return time.Now().UnixNano()
}

func (r *TestRunner) runSecretsManagerEdgeTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "Secret_NonExistent", func() error {
		// Every secret-scoped operation rejects an unknown SecretId with
		// ResourceNotFoundException.
		probes := []struct {
			name  string
			probe func() error
		}{
			{"GetSecretValue", func() error {
				_, err := tc.client.GetSecretValue(tc.ctx, &secretsmanager.GetSecretValueInput{
					SecretId: aws.String("nonexistent-secret-xyz"),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"DescribeSecret", func() error {
				_, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
					SecretId: aws.String("nonexistent-secret-xyz"),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"DeleteSecret", func() error {
				_, err := tc.client.DeleteSecret(tc.ctx, &secretsmanager.DeleteSecretInput{
					SecretId:                   aws.String("nonexistent-delete-xyz"),
					ForceDeleteWithoutRecovery: aws.Bool(true),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"RestoreSecret", func() error {
				_, err := tc.client.RestoreSecret(tc.ctx, &secretsmanager.RestoreSecretInput{
					SecretId: aws.String("nonexistent-restore-xyz"),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"GetResourcePolicy", func() error {
				_, err := tc.client.GetResourcePolicy(tc.ctx, &secretsmanager.GetResourcePolicyInput{
					SecretId: aws.String("nonexistent-policy-secret"),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"PutResourcePolicy", func() error {
				_, err := tc.client.PutResourcePolicy(tc.ctx, &secretsmanager.PutResourcePolicyInput{
					SecretId:       aws.String("nonexistent-policy-secret"),
					ResourcePolicy: aws.String(`{"Version":"2012-10-17"}`),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
		}
		for _, p := range probes {
			if err := p.probe(); err != nil {
				return fmt.Errorf("%s: %w", p.name, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_Filters", func() error {
		prefix := tc.uniqueName("FilterTest")
		for _, suffix := range []string{"alpha", "beta"} {
			_, err := tc.createSecret(prefix+"-"+suffix, suffix)
			if err != nil {
				return err
			}
			defer tc.forceDeleteSecret(prefix + "-" + suffix)
		}

		resp, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{
				{Key: types.FilterNameStringTypeName, Values: []string{prefix + "-alpha"}},
			},
		})
		if err != nil {
			return fmt.Errorf("list with filter: %v", err)
		}
		if len(resp.SecretList) != 1 {
			return fmt.Errorf("expected 1 secret, got %d", len(resp.SecretList))
		}
		if resp.SecretList[0].Name == nil || *resp.SecretList[0].Name != prefix+"-alpha" {
			return fmt.Errorf("wrong secret returned: %q", aws.ToString(resp.SecretList[0].Name))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_DefaultSortCreatedDate", func() error {
		base := tc.uniqueName("SortDefault")
		older := base + "-zzz" // created first, sorts last by name
		newer := base + "-aaa"
		for _, n := range []string{older, newer} {
			_, err := tc.createSecret(n, "sort-default")
			if err != nil {
				return err
			}
			defer tc.forceDeleteSecret(n)
		}

		pos := map[string]int{}
		idx := 0
		var nextToken *string
		for {
			resp, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
				Filters:   []types.Filter{{Key: types.FilterNameStringTypeName, Values: []string{base}}},
				NextToken: nextToken,
			})
			if err != nil {
				return fmt.Errorf("list with filter: %v", err)
			}
			for _, s := range resp.SecretList {
				name := aws.ToString(s.Name)
				if _, seen := pos[name]; !seen {
					pos[name] = idx
				}
				idx++
			}
			if resp.NextToken == nil {
				break
			}
			nextToken = resp.NextToken
		}

		posOlder, okOlder := pos[older]
		posNewer, okNewer := pos[newer]
		if !okOlder || !okNewer {
			return fmt.Errorf("filter did not return both sort fixtures: %v", pos)
		}
		if posOlder > posNewer {
			return fmt.Errorf("default sort is not CreatedDate: %s (created first) listed after %s", older, newer)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "List_MaxResultsOutOfRange", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			{"ListSecrets", func() error {
				for _, bad := range []int32{0, 101} {
					_, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
						MaxResults: aws.Int32(bad),
					})
					if err == nil {
						return fmt.Errorf("MaxResults=%d should be rejected", bad)
					}
					if e := expectAWSErrorCode(err, "InvalidParameterException"); e != nil {
						return fmt.Errorf("MaxResults=%d: %w", bad, e)
					}
				}
				return nil
			}},
			{"ListSecretVersionIds", func() error {
				name := tc.uniqueName("VerListRange")
				if _, err := tc.createSecret(name, "range"); err != nil {
					return err
				}
				defer tc.forceDeleteSecret(name)
				for _, bad := range []int32{0, 101} {
					_, err := tc.client.ListSecretVersionIds(tc.ctx, &secretsmanager.ListSecretVersionIdsInput{
						SecretId:   aws.String(name),
						MaxResults: aws.Int32(bad),
					})
					if err == nil {
						return fmt.Errorf("MaxResults=%d should be rejected", bad)
					}
					if e := expectAWSErrorCode(err, "InvalidParameterException"); e != nil {
						return fmt.Errorf("MaxResults=%d: %w", bad, e)
					}
				}
				return nil
			}},
			{"BatchGetSecretValue", func() error {
				prefix := tc.uniqueName("BatchRange")
				if _, err := tc.createSecret(prefix, "range"); err != nil {
					return err
				}
				defer tc.forceDeleteSecret(prefix)
				for _, bad := range []int32{0, 21} {
					_, err := tc.client.BatchGetSecretValue(tc.ctx, &secretsmanager.BatchGetSecretValueInput{
						Filters:    []types.Filter{{Key: types.FilterNameStringTypeName, Values: []string{prefix}}},
						MaxResults: aws.Int32(bad),
					})
					if err == nil {
						return fmt.Errorf("MaxResults=%d should be rejected", bad)
					}
					if e := expectAWSErrorCode(err, "InvalidParameterException"); e != nil {
						return fmt.Errorf("MaxResults=%d: %w", bad, e)
					}
				}
				// MaxResults is documented as requiring Filters; pairing it with
				// SecretIdList is rejected rather than silently ignored.
				_, err := tc.client.BatchGetSecretValue(tc.ctx, &secretsmanager.BatchGetSecretValueInput{
					SecretIdList: []string{prefix},
					MaxResults:   aws.Int32(5),
				})
				if err == nil {
					return fmt.Errorf("MaxResults with SecretIdList should be rejected")
				}
				if e := expectAWSErrorCode(err, "InvalidParameterException"); e != nil {
					return fmt.Errorf("MaxResults with SecretIdList: %w", e)
				}
				return nil
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %w", row.name, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "List_InvalidNextTokenRejected", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			{"ListSecrets", func() error {
				name := tc.uniqueName("BadToken")
				if _, err := tc.createSecret(name, "bad-token"); err != nil {
					return err
				}
				defer tc.forceDeleteSecret(name)
				for _, bad := range []string{"9999", "-1", "abc", "1x"} {
					_, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
						NextToken: aws.String(bad),
					})
					if err == nil {
						return fmt.Errorf("NextToken=%q should be rejected", bad)
					}
					if e := expectAWSErrorCode(err, "InvalidNextTokenException"); e != nil {
						return fmt.Errorf("NextToken=%q: %w", bad, e)
					}
				}
				return nil
			}},
			{"ListSecretVersionIds", func() error {
				name := tc.uniqueName("VerListBadToken")
				if _, err := tc.createSecret(name, "bad-token"); err != nil {
					return err
				}
				defer tc.forceDeleteSecret(name)
				for _, bad := range []string{"9999", "-1", "abc", "1x"} {
					_, err := tc.client.ListSecretVersionIds(tc.ctx, &secretsmanager.ListSecretVersionIdsInput{
						SecretId:  aws.String(name),
						NextToken: aws.String(bad),
					})
					if err == nil {
						return fmt.Errorf("NextToken=%q should be rejected", bad)
					}
					if e := expectAWSErrorCode(err, "InvalidNextTokenException"); e != nil {
						return fmt.Errorf("NextToken=%q: %w", bad, e)
					}
				}
				return nil
			}},
			{"BatchGetSecretValue", func() error {
				prefix := tc.uniqueName("BatchBadToken")
				if _, err := tc.createSecret(prefix, "bad-token"); err != nil {
					return err
				}
				defer tc.forceDeleteSecret(prefix)
				for _, bad := range []string{"9999", "-1", "abc", "1x"} {
					_, err := tc.client.BatchGetSecretValue(tc.ctx, &secretsmanager.BatchGetSecretValueInput{
						Filters:   []types.Filter{{Key: types.FilterNameStringTypeName, Values: []string{prefix}}},
						NextToken: aws.String(bad),
					})
					if err == nil {
						return fmt.Errorf("NextToken=%q should be rejected", bad)
					}
					if e := expectAWSErrorCode(err, "InvalidNextTokenException"); e != nil {
						return fmt.Errorf("NextToken=%q: %w", bad, e)
					}
				}
				return nil
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %w", row.name, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_FilterSemantics", func() error {
		base := tc.uniqueName("Pfx")
		name := base + "-DataSecret"
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("filter-semantics"),
			Description:  aws.String("My Filter Description"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		count := func(filters []types.Filter) (int, error) {
			resp, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
				Filters: filters,
			})
			if err != nil {
				return 0, err
			}
			n := 0
			for _, s := range resp.SecretList {
				if strings.HasPrefix(aws.ToString(s.Name), base) {
					n++
				}
			}
			return n, nil
		}

		// name: prefix, case-sensitive.
		if n, err := count([]types.Filter{{Key: types.FilterNameStringTypeName, Values: []string{name}}}); err != nil {
			return fmt.Errorf("name exact: %v", err)
		} else if n != 1 {
			return fmt.Errorf("name exact prefix: got %d secrets, want 1", n)
		}
		if n, err := count([]types.Filter{{Key: types.FilterNameStringTypeName, Values: []string{strings.ToLower(name)}}}); err != nil {
			return fmt.Errorf("name lower: %v", err)
		} else if n != 0 {
			return fmt.Errorf("name filter must be case-sensitive: got %d secrets, want 0", n)
		}
		if n, err := count([]types.Filter{{Key: types.FilterNameStringTypeName, Values: []string{"-DataSecret"}}}); err != nil {
			return fmt.Errorf("name suffix: %v", err)
		} else if n != 0 {
			return fmt.Errorf("name filter must be a prefix match, not substring: got %d secrets, want 0", n)
		}

		// description: prefix, not case-sensitive.
		if n, err := count([]types.Filter{{Key: types.FilterNameStringTypeDescription, Values: []string{"my filter"}}}); err != nil {
			return fmt.Errorf("description lower: %v", err)
		} else if n != 1 {
			return fmt.Errorf("description filter must be case-insensitive prefix: got %d secrets, want 1", n)
		}
		if n, err := count([]types.Filter{{Key: types.FilterNameStringTypeDescription, Values: []string{"filter description"}}}); err != nil {
			return fmt.Errorf("description infix: %v", err)
		} else if n != 0 {
			return fmt.Errorf("description filter must be a prefix match: got %d secrets, want 0", n)
		}

		// all: word-split, not case-sensitive.
		if n, err := count([]types.Filter{{Key: types.FilterNameStringTypeAll, Values: []string{"DATASECRET"}}}); err != nil {
			return fmt.Errorf("all upper: %v", err)
		} else if n != 1 {
			return fmt.Errorf("all filter must be case-insensitive word match: got %d secrets, want 1", n)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_FilterNegationAndUnknownKeyRejected", func() error {
		base := tc.uniqueName("Neg")
		for _, suffix := range []string{"keep", "drop"} {
			_, err := tc.createSecret(base+"-"+suffix, suffix)
			if err != nil {
				return err
			}
			defer tc.forceDeleteSecret(base + "-" + suffix)
		}

		// Negation: positive base prefix plus a negated value excludes the
		// matching secret.
		resp, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{
				{Key: types.FilterNameStringTypeName, Values: []string{base, "!" + base + "-drop"}},
			},
		})
		if err != nil {
			return fmt.Errorf("list with negation: %v", err)
		}
		names := map[string]bool{}
		for _, s := range resp.SecretList {
			if strings.HasPrefix(aws.ToString(s.Name), base) {
				names[aws.ToString(s.Name)] = true
			}
		}
		if names[base+"-drop"] {
			return fmt.Errorf("negated secret %s must be excluded", base+"-drop")
		}
		if !names[base+"-keep"] {
			return fmt.Errorf("non-negated secret %s must be included", base+"-keep")
		}

		// Unknown filter keys are rejected, not ignored.
		_, err = tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{{Key: types.FilterNameStringType("bogus"), Values: []string{"x"}}},
		})
		if err == nil {
			return fmt.Errorf("unknown filter key must be rejected")
		}
		if e := expectAWSErrorCode(err, "InvalidParameterException"); e != nil {
			return fmt.Errorf("unknown filter key: %v", e)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", pgTimestamp())
		var pgSecrets []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagSecret-%s-%d", pgTs, i)
			_, err := tc.createSecret(name, "pagval")
			if err != nil {
				for _, sn := range pgSecrets {
					tc.forceDeleteSecret(sn)
				}
				return err
			}
			pgSecrets = append(pgSecrets, name)
		}

		var allSecrets []string
		var nextToken *string
		for {
			resp, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, sn := range pgSecrets {
					tc.forceDeleteSecret(sn)
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
			tc.forceDeleteSecret(sn)
		}
		if len(allSecrets) != 5 {
			return fmt.Errorf("expected 5 paginated secrets, got %d", len(allSecrets))
		}
		return nil
	}))

	return results
}
