package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

// cfContinuousDeploymentTests covers the continuous deployment policy
// family: create, get, get-config, list, update, attach to a primary
// distribution, the in-use delete guard, detach, delete, and the
// documented configuration validation.
func cfContinuousDeploymentTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	distID, _, err := tc.createDistribution(tc.uniquePrefix("cdp-primary"), "CDP test distribution", "127.0.0.1:50080")
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudfront",
			TestName: "CDP_Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		})
	}

	// The staging distribution comes from CopyDistribution, the only
	// documented way to obtain one.
	copyResp, err := client.CopyDistribution(ctx, &cloudfront.CopyDistributionInput{
		PrimaryDistributionId: aws.String(distID),
		CallerReference:       aws.String(tc.uniquePrefix("cdp-staging")),
		Staging:               aws.Bool(true),
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudfront",
			TestName: "CDP_Lifecycle",
			Status:   "FAIL",
			Error:    fmt.Sprintf("copy distribution: %v", err),
		})
	}
	stagingID := aws.ToString(copyResp.Distribution.Id)
	stagingDomain := aws.ToString(copyResp.Distribution.DomainName)
	if stagingDomain == "" {
		stagingDomain = stagingID + ".cloudfront.net"
	}

	defer func() { _ = tc.disableAndDeleteDistributionByID(stagingID) }()
	defer func() { _ = tc.disableAndDeleteDistributionByID(distID) }()

	policyConfig := func(weight float64) *cftypes.ContinuousDeploymentPolicyConfig {
		return &cftypes.ContinuousDeploymentPolicyConfig{
			StagingDistributionDnsNames: &cftypes.StagingDistributionDnsNames{
				Quantity: aws.Int32(1),
				Items:    []string{stagingDomain},
			},
			Enabled: aws.Bool(true),
			TrafficConfig: &cftypes.TrafficConfig{
				Type: cftypes.ContinuousDeploymentPolicyTypeSingleWeight,
				SingleWeightConfig: &cftypes.ContinuousDeploymentSingleWeightConfig{
					Weight: aws.Float32(float32(weight)),
					SessionStickinessConfig: &cftypes.SessionStickinessConfig{
						IdleTTL:    aws.Int32(300),
						MaximumTTL: aws.Int32(600),
					},
				},
			},
		}
	}

	var policyID, policyETag string
	results = append(results, tc.runner.RunTest("cloudfront", "CDP_Create_Verify", func() error {
		resp, err := client.CreateContinuousDeploymentPolicy(ctx, &cloudfront.CreateContinuousDeploymentPolicyInput{
			ContinuousDeploymentPolicyConfig: policyConfig(0.10),
		})
		if err != nil {
			return err
		}
		if resp.ContinuousDeploymentPolicy == nil || aws.ToString(resp.ContinuousDeploymentPolicy.Id) == "" {
			return fmt.Errorf("created policy carries no id")
		}
		if aws.ToString(resp.ETag) == "" {
			return fmt.Errorf("created policy carries no ETag")
		}
		policyID = aws.ToString(resp.ContinuousDeploymentPolicy.Id)
		policyETag = aws.ToString(resp.ETag)
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CDP_Get_Verify", func() error {
		resp, err := client.GetContinuousDeploymentPolicy(ctx, &cloudfront.GetContinuousDeploymentPolicyInput{Id: aws.String(policyID)})
		if err != nil {
			return err
		}
		cfg := resp.ContinuousDeploymentPolicy.ContinuousDeploymentPolicyConfig
		if cfg == nil || !aws.ToBool(cfg.Enabled) {
			return fmt.Errorf("policy config missing or disabled")
		}
		// The wire type is a 32-bit float; compare with tolerance.
		weight := float64(aws.ToFloat32(cfg.TrafficConfig.SingleWeightConfig.Weight))
		if weight < 0.099 || weight > 0.101 {
			return fmt.Errorf("weight mismatch: got %v", weight)
		}
		if len(cfg.StagingDistributionDnsNames.Items) != 1 || cfg.StagingDistributionDnsNames.Items[0] != stagingDomain {
			return fmt.Errorf("staging DNS mismatch: %v", cfg.StagingDistributionDnsNames.Items)
		}
		if aws.ToString(resp.ETag) != policyETag {
			return fmt.Errorf("etag drift between create and get")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CDP_GetConfig_Verify", func() error {
		resp, err := client.GetContinuousDeploymentPolicyConfig(ctx, &cloudfront.GetContinuousDeploymentPolicyConfigInput{Id: aws.String(policyID)})
		if err != nil {
			return err
		}
		stickiness := resp.ContinuousDeploymentPolicyConfig.TrafficConfig.SingleWeightConfig.SessionStickinessConfig
		if stickiness == nil || aws.ToInt32(stickiness.IdleTTL) != 300 || aws.ToInt32(stickiness.MaximumTTL) != 600 {
			return fmt.Errorf("stickiness mismatch: %+v", stickiness)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CDP_List_ContainsPolicy", func() error {
		policies, err := paginate(func(next *string) ([]cftypes.ContinuousDeploymentPolicySummary, *string, error) {
			resp, lerr := client.ListContinuousDeploymentPolicies(ctx, &cloudfront.ListContinuousDeploymentPoliciesInput{Marker: next})
			if lerr != nil {
				return nil, nil, lerr
			}
			if resp.ContinuousDeploymentPolicyList == nil {
				return nil, nil, nil
			}
			return resp.ContinuousDeploymentPolicyList.Items, resp.ContinuousDeploymentPolicyList.NextMarker, nil
		})
		if err != nil {
			return err
		}
		for _, summary := range policies {
			if aws.ToString(summary.ContinuousDeploymentPolicy.Id) == policyID {
				return nil
			}
		}
		return fmt.Errorf("policy %q not found in the list", policyID)
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CDP_Update_Verify", func() error {
		resp, err := client.UpdateContinuousDeploymentPolicy(ctx, &cloudfront.UpdateContinuousDeploymentPolicyInput{
			Id:                               aws.String(policyID),
			IfMatch:                          aws.String(policyETag),
			ContinuousDeploymentPolicyConfig: policyConfig(0.05),
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.ETag) == policyETag {
			return fmt.Errorf("update must rotate the ETag")
		}
		policyETag = aws.ToString(resp.ETag)
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CDP_AttachToDistribution", func() error {
		cfgResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)})
		if err != nil {
			return err
		}
		cfgResp.DistributionConfig.ContinuousDeploymentPolicyId = aws.String(policyID)
		updResp, err := client.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
			Id: aws.String(distID), IfMatch: cfgResp.ETag, DistributionConfig: cfgResp.DistributionConfig,
		})
		if err != nil {
			return err
		}
		policyETagDist := aws.ToString(updResp.ETag)
		afterResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)})
		if err != nil {
			return err
		}
		if aws.ToString(afterResp.DistributionConfig.ContinuousDeploymentPolicyId) != policyID {
			return fmt.Errorf("distribution policy id mismatch: got %q", aws.ToString(afterResp.DistributionConfig.ContinuousDeploymentPolicyId))
		}
		_ = policyETagDist

		// A policy still attached to a distribution cannot be deleted.
		_, err = client.DeleteContinuousDeploymentPolicy(ctx, &cloudfront.DeleteContinuousDeploymentPolicyInput{
			Id: aws.String(policyID), IfMatch: aws.String(policyETag),
		})
		return AssertErrorContains(err, "ContinuousDeploymentPolicyInUse")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CDP_DetachAndDelete", func() error {
		cfgResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)})
		if err != nil {
			return err
		}
		cfgResp.DistributionConfig.ContinuousDeploymentPolicyId = aws.String("")
		if _, err := client.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
			Id: aws.String(distID), IfMatch: cfgResp.ETag, DistributionConfig: cfgResp.DistributionConfig,
		}); err != nil {
			return err
		}

		if _, err := client.DeleteContinuousDeploymentPolicy(ctx, &cloudfront.DeleteContinuousDeploymentPolicyInput{
			Id: aws.String(policyID), IfMatch: aws.String(policyETag),
		}); err != nil {
			return err
		}
		_, err = client.GetContinuousDeploymentPolicy(ctx, &cloudfront.GetContinuousDeploymentPolicyInput{Id: aws.String(policyID)})
		return AssertErrorContains(err, "NoSuchContinuousDeploymentPolicy")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CDP_Create_WeightAboveQuota_Rejected", func() error {
		_, err := client.CreateContinuousDeploymentPolicy(ctx, &cloudfront.CreateContinuousDeploymentPolicyInput{
			ContinuousDeploymentPolicyConfig: policyConfig(0.2),
		})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CDP_Create_UnknownStagingDistribution_Rejected", func() error {
		cfg := policyConfig(0.1)
		cfg.StagingDistributionDnsNames = &cftypes.StagingDistributionDnsNames{
			Quantity: aws.Int32(1),
			Items:    []string{"d99999999zzzzzz.cloudfront.net"},
		}
		_, err := client.CreateContinuousDeploymentPolicy(ctx, &cloudfront.CreateContinuousDeploymentPolicyInput{
			ContinuousDeploymentPolicyConfig: cfg,
		})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CDP_Create_HeaderWithoutPrefix_Rejected", func() error {
		_, err := client.CreateContinuousDeploymentPolicy(ctx, &cloudfront.CreateContinuousDeploymentPolicyInput{
			ContinuousDeploymentPolicyConfig: &cftypes.ContinuousDeploymentPolicyConfig{
				StagingDistributionDnsNames: &cftypes.StagingDistributionDnsNames{
					Quantity: aws.Int32(1),
					Items:    []string{stagingDomain},
				},
				Enabled: aws.Bool(true),
				TrafficConfig: &cftypes.TrafficConfig{
					Type: cftypes.ContinuousDeploymentPolicyTypeSingleHeader,
					SingleHeaderConfig: &cftypes.ContinuousDeploymentSingleHeaderConfig{
						Header: aws.String("x-missing-prefix"),
						Value:  aws.String("staging"),
					},
				},
			},
		})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CDP_Create_StickinessBounds_Rejected", func() error {
		cfg := policyConfig(0.1)
		cfg.TrafficConfig.SingleWeightConfig.SessionStickinessConfig = &cftypes.SessionStickinessConfig{
			IdleTTL:    aws.Int32(60),
			MaximumTTL: aws.Int32(600),
		}
		_, err := client.CreateContinuousDeploymentPolicy(ctx, &cloudfront.CreateContinuousDeploymentPolicyInput{
			ContinuousDeploymentPolicyConfig: cfg,
		})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	return results
}
