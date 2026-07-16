package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (r *TestRunner) runKMSGrantTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kms", "CreateGrant", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		granteePrincipal := "arn:aws:iam::000000000000:user/TestUser"
		resp, err := tc.client.CreateGrant(tc.ctx, &kms.CreateGrantInput{
			KeyId:            aws.String(tc.keyID),
			GranteePrincipal: aws.String(granteePrincipal),
			Operations:       []types.GrantOperation{types.GrantOperationEncrypt},
		})
		if err != nil {
			return err
		}
		if resp.GrantToken == nil || *resp.GrantToken == "" {
			return fmt.Errorf("grant token is nil or empty")
		}
		if resp.GrantId == nil || *resp.GrantId == "" {
			return fmt.Errorf("grant ID is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListGrants", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := tc.client.ListGrants(tc.ctx, &kms.ListGrantsInput{
			KeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return err
		}
		if resp.Grants == nil {
			return fmt.Errorf("grants list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListRetirableGrants", func() error {
		resp, err := tc.client.ListRetirableGrants(tc.ctx, &kms.ListRetirableGrantsInput{
			RetiringPrincipal: aws.String("arn:aws:iam::000000000000:user/TestUser"),
		})
		if err != nil {
			return err
		}
		if resp.Grants == nil {
			return fmt.Errorf("grants list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "RetireGrant", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		granteePrincipal := "arn:aws:iam::000000000000:user/RetireUser"
		createResp, err := tc.client.CreateGrant(tc.ctx, &kms.CreateGrantInput{
			KeyId:            aws.String(tc.keyID),
			GranteePrincipal: aws.String(granteePrincipal),
			Operations:       []types.GrantOperation{types.GrantOperationDecrypt},
		})
		if err != nil {
			return fmt.Errorf("create grant: %v", err)
		}

		_, err = tc.client.RetireGrant(tc.ctx, &kms.RetireGrantInput{
			GrantId: createResp.GrantId,
			KeyId:   aws.String(tc.keyID),
		})
		if err != nil {
			return fmt.Errorf("retire grant: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "RevokeGrant", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		granteePrincipal := "arn:aws:iam::000000000000:user/RevokeUser"
		createResp, err := tc.client.CreateGrant(tc.ctx, &kms.CreateGrantInput{
			KeyId:            aws.String(tc.keyID),
			GranteePrincipal: aws.String(granteePrincipal),
			Operations:       []types.GrantOperation{types.GrantOperationEncrypt},
		})
		if err != nil {
			return fmt.Errorf("create grant: %v", err)
		}

		_, err = tc.client.RevokeGrant(tc.ctx, &kms.RevokeGrantInput{
			KeyId:   aws.String(tc.keyID),
			GrantId: createResp.GrantId,
		})
		if err != nil {
			return fmt.Errorf("revoke grant: %v", err)
		}

		listResp, err := tc.client.ListGrants(tc.ctx, &kms.ListGrantsInput{KeyId: aws.String(tc.keyID)})
		if err != nil {
			return fmt.Errorf("list grants: %v", err)
		}
		for _, g := range listResp.Grants {
			if aws.ToString(g.GrantId) == *createResp.GrantId {
				return fmt.Errorf("revoked grant still in list")
			}
		}
		return nil
	}))

	return results
}
