package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func (r *TestRunner) runWAFv2RegexPatternSetTests(tc *wafv2TestContext) []TestResult {
	var results []TestResult

	regexName := tc.uniqueName("test-regex")
	var regexID, regexLockToken string

	results = append(results, r.RunTest("wafv2", "CreateRegexPatternSet", func() error {
		resp, err := tc.client.CreateRegexPatternSet(tc.ctx, &wafv2.CreateRegexPatternSetInput{
			Name:  aws.String(regexName),
			Scope: tc.scope,
			RegularExpressionList: []types.Regex{
				{RegexString: aws.String(`[a-z]+@[a-z]+\.[a-z]+`)},
			},
		})
		if err != nil {
			return err
		}
		if resp.Summary == nil {
			return fmt.Errorf("Summary is nil")
		}
		if aws.ToString(resp.Summary.Id) == "" {
			return fmt.Errorf("Summary.Id is empty")
		}
		if aws.ToString(resp.Summary.LockToken) == "" {
			return fmt.Errorf("Summary.LockToken is empty")
		}
		if aws.ToString(resp.Summary.Name) != regexName {
			return fmt.Errorf("Summary.Name mismatch: expected %s, got %s", regexName, aws.ToString(resp.Summary.Name))
		}
		regexID = aws.ToString(resp.Summary.Id)
		regexLockToken = aws.ToString(resp.Summary.LockToken)
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetRegexPatternSet", func() error {
		resp, err := tc.client.GetRegexPatternSet(tc.ctx, &wafv2.GetRegexPatternSetInput{
			Name: aws.String(regexName), Scope: tc.scope, Id: aws.String(regexID),
		})
		if err != nil {
			return err
		}
		if resp.RegexPatternSet == nil {
			return fmt.Errorf("RegexPatternSet is nil")
		}
		if aws.ToString(resp.RegexPatternSet.Name) != regexName {
			return fmt.Errorf("Name mismatch: expected %s, got %s", regexName, aws.ToString(resp.RegexPatternSet.Name))
		}
		if aws.ToString(resp.RegexPatternSet.Id) != regexID {
			return fmt.Errorf("Id mismatch")
		}
		if len(resp.RegexPatternSet.RegularExpressionList) != 1 {
			return fmt.Errorf("expected 1 pattern, got %d", len(resp.RegexPatternSet.RegularExpressionList))
		}
		if resp.LockToken == nil {
			return fmt.Errorf("LockToken is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListRegexPatternSets", func() error {
		resp, err := tc.client.ListRegexPatternSets(tc.ctx, &wafv2.ListRegexPatternSetsInput{
			Scope: tc.scope,
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.RegexPatternSets == nil {
			return fmt.Errorf("RegexPatternSets list is nil")
		}
		found := false
		for _, s := range resp.RegexPatternSets {
			if s.Id != nil && *s.Id == regexID {
				found = true
				if aws.ToString(s.Name) != regexName {
					return fmt.Errorf("RegexPatternSet name mismatch in list")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("RegexPatternSet not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateRegexPatternSet", func() error {
		resp, err := tc.client.UpdateRegexPatternSet(tc.ctx, &wafv2.UpdateRegexPatternSetInput{
			Name: aws.String(regexName), Scope: tc.scope,
			Id: aws.String(regexID), LockToken: aws.String(regexLockToken),
			RegularExpressionList: []types.Regex{
				{RegexString: aws.String(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)},
				{RegexString: aws.String(`[0-9a-f]{2}:[0-9a-f]{2}:[0-9a-f]{2}`)},
			},
		})
		if err != nil {
			return err
		}
		if resp.NextLockToken == nil || aws.ToString(resp.NextLockToken) == "" {
			return fmt.Errorf("NextLockToken is nil or empty")
		}
		regexLockToken = aws.ToString(resp.NextLockToken)
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateRegexPatternSet_ContentVerify", func() error {
		resp, err := tc.client.GetRegexPatternSet(tc.ctx, &wafv2.GetRegexPatternSetInput{
			Name: aws.String(regexName), Scope: tc.scope, Id: aws.String(regexID),
		})
		if err != nil {
			return err
		}
		if len(resp.RegexPatternSet.RegularExpressionList) != 2 {
			return fmt.Errorf("expected 2 patterns, got %d", len(resp.RegexPatternSet.RegularExpressionList))
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DeleteRegexPatternSet", func() error {
		resp, err := tc.client.DeleteRegexPatternSet(tc.ctx, &wafv2.DeleteRegexPatternSetInput{
			Name: aws.String(regexName), Scope: tc.scope,
			Id: aws.String(regexID), LockToken: aws.String(regexLockToken),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetRegexPatternSet_NonExistent", func() error {
		_, err := tc.client.GetRegexPatternSet(tc.ctx, &wafv2.GetRegexPatternSetInput{
			Name: aws.String(regexName), Scope: tc.scope, Id: aws.String(regexID),
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	return results
}
