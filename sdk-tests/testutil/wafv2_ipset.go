package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func (r *TestRunner) runWAFv2IPSetTests(tc *wafv2TestContext) []TestResult {
	var results []TestResult

	ipSetName := tc.uniqueName("test-ipset")
	var ipSetID, ipSetARN, ipSetLockToken string

	results = append(results, r.RunTest("wafv2", "ListIPSets_Empty", func() error {
		resp, err := tc.client.ListIPSets(tc.ctx, &wafv2.ListIPSetsInput{
			Scope: tc.scope,
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.IPSets == nil {
			return fmt.Errorf("IPSets is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "CreateIPSet", func() error {
		resp, err := tc.client.CreateIPSet(tc.ctx, &wafv2.CreateIPSetInput{
			Name:             aws.String(ipSetName),
			Scope:            tc.scope,
			IPAddressVersion: types.IPAddressVersionIpv4,
			Addresses:        []string{},
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Owner"), Value: aws.String("sdk-tests")},
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
		if aws.ToString(resp.Summary.ARN) == "" {
			return fmt.Errorf("Summary.ARN is empty")
		}
		if aws.ToString(resp.Summary.LockToken) == "" {
			return fmt.Errorf("Summary.LockToken is empty")
		}
		if aws.ToString(resp.Summary.Name) != ipSetName {
			return fmt.Errorf("Summary.Name mismatch: expected %s, got %s", ipSetName, aws.ToString(resp.Summary.Name))
		}
		ipSetID = aws.ToString(resp.Summary.Id)
		ipSetARN = aws.ToString(resp.Summary.ARN)
		ipSetLockToken = aws.ToString(resp.Summary.LockToken)
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetIPSet", func() error {
		resp, err := tc.client.GetIPSet(tc.ctx, &wafv2.GetIPSetInput{
			Name: aws.String(ipSetName), Scope: tc.scope, Id: aws.String(ipSetID),
		})
		if err != nil {
			return err
		}
		if resp.IPSet == nil {
			return fmt.Errorf("IPSet is nil")
		}
		if resp.IPSet.IPAddressVersion != types.IPAddressVersionIpv4 {
			return fmt.Errorf("expected IPV4, got %s", resp.IPSet.IPAddressVersion)
		}
		if aws.ToString(resp.IPSet.Name) != ipSetName {
			return fmt.Errorf("Name mismatch: expected %s, got %s", ipSetName, aws.ToString(resp.IPSet.Name))
		}
		if aws.ToString(resp.IPSet.Id) != ipSetID {
			return fmt.Errorf("Id mismatch")
		}
		if resp.LockToken == nil {
			return fmt.Errorf("LockToken is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListIPSets_ContainsCreated", func() error {
		resp, err := tc.client.ListIPSets(tc.ctx, &wafv2.ListIPSetsInput{
			Scope: tc.scope,
		})
		if err != nil {
			return err
		}
		found := false
		for _, s := range resp.IPSets {
			if s.Id != nil && *s.Id == ipSetID {
				found = true
				if aws.ToString(s.Name) != ipSetName {
					return fmt.Errorf("IPSet name mismatch")
				}
				if aws.ToString(s.ARN) == "" {
					return fmt.Errorf("IPSet ARN is empty")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("IPSet not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListTagsForResource_IPSet", func() error {
		resp, err := tc.client.ListTagsForResource(tc.ctx, &wafv2.ListTagsForResourceInput{
			ResourceARN: aws.String(ipSetARN),
		})
		if err != nil {
			return err
		}
		if resp.TagInfoForResource == nil || resp.TagInfoForResource.TagList == nil {
			return fmt.Errorf("expected tags in response")
		}
		if len(resp.TagInfoForResource.TagList) < 2 {
			return fmt.Errorf("expected at least 2 tags, got %d", len(resp.TagInfoForResource.TagList))
		}
		if aws.ToString(resp.TagInfoForResource.ResourceARN) != ipSetARN {
			return fmt.Errorf("ResourceARN mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateIPSet", func() error {
		resp, err := tc.client.UpdateIPSet(tc.ctx, &wafv2.UpdateIPSetInput{
			Name: aws.String(ipSetName), Scope: tc.scope,
			Id: aws.String(ipSetID), LockToken: aws.String(ipSetLockToken),
			Addresses: []string{"192.0.2.0/24", "203.0.113.0/24"},
		})
		if err != nil {
			return err
		}
		if resp.NextLockToken == nil || aws.ToString(resp.NextLockToken) == "" {
			return fmt.Errorf("NextLockToken is nil or empty")
		}
		ipSetLockToken = aws.ToString(resp.NextLockToken)
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateIPSet_ContentVerify", func() error {
		resp, err := tc.client.GetIPSet(tc.ctx, &wafv2.GetIPSetInput{
			Name: aws.String(ipSetName), Scope: tc.scope, Id: aws.String(ipSetID),
		})
		if err != nil {
			return err
		}
		if len(resp.IPSet.Addresses) != 2 {
			return fmt.Errorf("expected 2 addresses, got %d", len(resp.IPSet.Addresses))
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DeleteIPSet", func() error {
		resp, err := tc.client.DeleteIPSet(tc.ctx, &wafv2.DeleteIPSetInput{
			Name: aws.String(ipSetName), Scope: tc.scope,
			Id: aws.String(ipSetID), LockToken: aws.String(ipSetLockToken),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetIPSet_NonExistent", func() error {
		_, err := tc.client.GetIPSet(tc.ctx, &wafv2.GetIPSetInput{
			Name: aws.String(ipSetName), Scope: tc.scope, Id: aws.String(ipSetID),
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	results = append(results, r.RunTest("wafv2", "DeleteIPSet_NonExistent", func() error {
		_, err := tc.client.DeleteIPSet(tc.ctx, &wafv2.DeleteIPSetInput{
			Name: aws.String(ipSetName), Scope: tc.scope,
			Id: aws.String(ipSetID), LockToken: aws.String("fake-lock-token"),
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	results = append(results, r.RunTest("wafv2", "UpdateIPSet_StaleLockToken", func() error {
		vName := tc.uniqueName("verify-ipset")
		vID, _, vLock, err := tc.createIPSet(vName, nil)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteIPSet(vName, vID, vLock)

		_, err = tc.client.UpdateIPSet(tc.ctx, &wafv2.UpdateIPSetInput{
			Name: aws.String(vName), Scope: tc.scope,
			Id: aws.String(vID), LockToken: aws.String("stale-token-should-fail"),
			Addresses: []string{"192.168.0.0/16"},
		})
		return AssertErrorContains(err, "WAFInvalidLockTokenException")
	}))

	results = append(results, r.RunTest("wafv2", "ListIPSets_Pagination", func() error {
		return r.wafv2ListIPSetsPagination(tc)
	}))

	return results
}

func (r *TestRunner) wafv2ListIPSetsPagination(tc *wafv2TestContext) error {
	pgTs := tc.uniqueName("pgwaf")
	type pgIPSet struct {
		id, name, lockToken string
	}
	var pgSets []pgIPSet
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("%s-%d", pgTs, i)
		id, _, lock, err := tc.createIPSet(name, nil)
		if err != nil {
			return fmt.Errorf("create ip set %s: %v", name, err)
		}
		pgSets = append(pgSets, pgIPSet{id: id, name: name, lockToken: lock})
	}

	pgIDSet := make(map[string]bool)
	for _, s := range pgSets {
		pgIDSet[s.id] = true
	}
	var foundCount int
	var nextMarker *string
	for {
		resp, err := tc.client.ListIPSets(tc.ctx, &wafv2.ListIPSetsInput{
			Scope:      tc.scope,
			Limit:      aws.Int32(2),
			NextMarker: nextMarker,
		})
		if err != nil {
			for _, s := range pgSets {
				tc.deleteIPSet(s.name, s.id, s.lockToken)
			}
			return fmt.Errorf("list ip sets page: %v", err)
		}
		for _, s := range resp.IPSets {
			if pgIDSet[aws.ToString(s.Id)] {
				foundCount++
			}
		}
		if resp.NextMarker != nil && *resp.NextMarker != "" {
			nextMarker = resp.NextMarker
		} else {
			break
		}
	}

	for _, s := range pgSets {
		tc.deleteIPSet(s.name, s.id, s.lockToken)
	}
	if foundCount != 5 {
		return fmt.Errorf("expected 5 paginated IP sets, got %d", foundCount)
	}
	return nil
}
