package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func (r *TestRunner) runKMSAliasTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kms", "CreateAlias", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(tc.keyAlias),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return err
		}
		tc.addCleanupAlias(tc.keyAlias)
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListAliases", func() error {
		resp, err := tc.client.ListAliases(tc.ctx, &kms.ListAliasesInput{})
		if err != nil {
			return err
		}
		if resp.Aliases == nil {
			return fmt.Errorf("aliases list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "UpdateAlias", func() error {
		_, err := tc.client.UpdateAlias(tc.ctx, &kms.UpdateAliasInput{
			AliasName:   aws.String(tc.keyAlias),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return err
		}
		listResp, err := tc.client.ListAliases(tc.ctx, &kms.ListAliasesInput{})
		if err != nil {
			return fmt.Errorf("list aliases after update: %v", err)
		}
		found := false
		for _, a := range listResp.Aliases {
			if aws.ToString(a.AliasName) == tc.keyAlias {
				found = true
				if a.TargetKeyId == nil || *a.TargetKeyId != tc.keyID {
					return fmt.Errorf("alias target key mismatch")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("alias %q not found after update", tc.keyAlias)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListAliases_ContainsCreated", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		testAlias := fmt.Sprintf("alias/list-test-%d", time.Now().UnixNano())
		_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(testAlias),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		tc.addCleanupAlias(testAlias)

		resp, err := tc.client.ListAliases(tc.ctx, &kms.ListAliasesInput{})
		if err != nil {
			return err
		}
		found := false
		for _, a := range resp.Aliases {
			if a.AliasName != nil && *a.AliasName == testAlias {
				found = true
				if a.TargetKeyId == nil || *a.TargetKeyId != tc.keyID {
					return fmt.Errorf("alias target key mismatch")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created alias %q not found in ListAliases", testAlias)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListAliases_FilterByKeyID", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		filterAlias := fmt.Sprintf("alias/filter-test-%d", time.Now().UnixNano())
		_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(filterAlias),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		tc.addCleanupAlias(filterAlias)

		resp, err := tc.client.ListAliases(tc.ctx, &kms.ListAliasesInput{
			KeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return err
		}
		for _, a := range resp.Aliases {
			if a.AliasName != nil && strings.HasPrefix(*a.AliasName, "alias/aws/") {
				continue
			}
			if a.TargetKeyId != nil && *a.TargetKeyId != tc.keyID {
				return fmt.Errorf("alias %q has wrong target key %q", *a.AliasName, *a.TargetKeyId)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateAlias_Duplicate", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		dupAlias := fmt.Sprintf("alias/dup-test-%d", time.Now().UnixNano())
		_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(dupAlias),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		tc.addCleanupAlias(dupAlias)
		_, err = tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(dupAlias),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate alias")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DeleteAlias", func() error {
		delAlias := fmt.Sprintf("alias/del-test-%d", time.Now().UnixNano())
		_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(delAlias),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias for delete: %v", err)
		}
		_, err = tc.client.DeleteAlias(tc.ctx, &kms.DeleteAliasInput{
			AliasName: aws.String(delAlias),
		})
		return err
	}))

	return results
}
