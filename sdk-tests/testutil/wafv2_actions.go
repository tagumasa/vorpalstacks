package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

// actionURIStatement builds the plain statement the action tests hang their
// rules on; the action, not the statement, is what these tests exercise.
func actionURIStatement(fragment string) *types.Statement {
	return &types.Statement{
		ByteMatchStatement: &types.ByteMatchStatement{
			FieldToMatch:         &types.FieldToMatch{UriPath: &types.UriPath{}},
			PositionalConstraint: types.PositionalConstraintContains,
			SearchString:         []byte(fragment),
			TextTransformations: []types.TextTransformation{
				{Priority: 0, Type: types.TextTransformationTypeNone},
			},
		},
	}
}

// testMonetizationConfig is a complete, valid MonetizationConfig on the
// Base Sepolia test network; the wallet is a checksummed EIP-55 vector.
func testMonetizationConfig() *types.MonetizationConfig {
	return testMonetizationConfigWithAmount("0.010")
}

// testMonetizationConfigWithAmount builds the reference monetization
// configuration around a specific price amount.
func testMonetizationConfigWithAmount(amount string) *types.MonetizationConfig {
	return &types.MonetizationConfig{
		CurrencyMode: types.CurrencyModeTest,
		CryptoConfig: &types.CryptoConfig{
			PaymentNetworks: []types.PaymentNetwork{{
				Chain:         types.BlockchainChainBaseSepolia,
				WalletAddress: aws.String("0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"),
				Prices: []types.Price{{
					Amount:   aws.String(amount),
					Currency: types.CryptoCurrencyUsdc,
				}},
			}},
		},
	}
}

func (r *TestRunner) runWAFv2ActionTests(tc *wafv2TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("wafv2", "RuleAction_CaptchaChallenge_RoundTrip", func() error {
		name := tc.uniqueName("test-captcha-acl")
		create, err := tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
			Name:  aws.String(name),
			Scope: tc.scope,
			DefaultAction: &types.DefaultAction{
				Allow: &types.AllowAction{},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("captcha-roundtrip-metric"),
			},
			Rules: []types.Rule{
				{
					Name:     aws.String("CaptchaRule"),
					Priority: 1,
					Action:   &types.RuleAction{Captcha: &types.CaptchaAction{}},
					CaptchaConfig: &types.CaptchaConfig{
						ImmunityTimeProperty: &types.ImmunityTimeProperty{ImmunityTime: aws.Int64(600)},
					},
					Statement:        actionURIStatement("/captcha-path"),
					VisibilityConfig: wafActionRuleVisibility("CaptchaRule"),
				},
				{
					Name:     aws.String("ChallengeRule"),
					Priority: 2,
					Action:   &types.RuleAction{Challenge: &types.ChallengeAction{}},
					ChallengeConfig: &types.ChallengeConfig{
						ImmunityTimeProperty: &types.ImmunityTimeProperty{ImmunityTime: aws.Int64(300)},
					},
					Statement:        actionURIStatement("/challenge-path"),
					VisibilityConfig: wafActionRuleVisibility("ChallengeRule"),
				},
			},
		})
		if err != nil {
			return err
		}
		aclID := aws.ToString(create.Summary.Id)
		lockToken := aws.ToString(create.Summary.LockToken)
		defer tc.deleteWebACL(name, aclID, lockToken)

		resp, err := tc.client.GetWebACL(tc.ctx, &wafv2.GetWebACLInput{
			Name: aws.String(name), Scope: tc.scope, Id: aws.String(aclID),
		})
		if err != nil {
			return err
		}
		if len(resp.WebACL.Rules) != 2 {
			return fmt.Errorf("expected 2 rules, got %d", len(resp.WebACL.Rules))
		}
		captchaRule, challengeRule := resp.WebACL.Rules[0], resp.WebACL.Rules[1]
		if captchaRule.Action == nil || captchaRule.Action.Captcha == nil {
			return fmt.Errorf("first rule lost its Captcha action: %+v", captchaRule.Action)
		}
		if captchaRule.CaptchaConfig == nil ||
			captchaRule.CaptchaConfig.ImmunityTimeProperty == nil ||
			aws.ToInt64(captchaRule.CaptchaConfig.ImmunityTimeProperty.ImmunityTime) != 600 {
			return fmt.Errorf("CaptchaConfig immunity time did not round-trip: %+v", captchaRule.CaptchaConfig)
		}
		if challengeRule.Action == nil || challengeRule.Action.Challenge == nil {
			return fmt.Errorf("second rule lost its Challenge action: %+v", challengeRule.Action)
		}
		if challengeRule.ChallengeConfig == nil ||
			challengeRule.ChallengeConfig.ImmunityTimeProperty == nil ||
			aws.ToInt64(challengeRule.ChallengeConfig.ImmunityTimeProperty.ImmunityTime) != 300 {
			return fmt.Errorf("ChallengeConfig immunity time did not round-trip: %+v", challengeRule.ChallengeConfig)
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "Monetize_WithoutConfig_Rejected", func() error {
		_, err := tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
			Name:  aws.String(tc.uniqueName("test-monetize-noconf")),
			Scope: tc.scope,
			DefaultAction: &types.DefaultAction{
				Allow: &types.AllowAction{},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   false,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("monetize-noconf-metric"),
			},
			Rules: []types.Rule{
				{
					Name:     aws.String("MonetizeRule"),
					Priority: 1,
					Action: &types.RuleAction{
						Monetize: &types.MonetizeAction{PriceMultiplier: aws.String("2")},
					},
					Statement:        actionURIStatement("/premium"),
					VisibilityConfig: wafActionRuleVisibility("MonetizeRule"),
				},
			},
		})
		return AssertErrorContains(err, "WAFInvalidParameterException")
	}))

	results = append(results, r.RunTest("wafv2", "Monetize_RegionalScope_Rejected", func() error {
		_, err := tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
			Name:  aws.String(tc.uniqueName("test-monetize-regional")),
			Scope: types.ScopeRegional,
			DefaultAction: &types.DefaultAction{
				Allow: &types.AllowAction{},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   false,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("monetize-regional-metric"),
			},
			MonetizationConfig: testMonetizationConfig(),
			Rules: []types.Rule{
				{
					Name:     aws.String("MonetizeRule"),
					Priority: 1,
					Action: &types.RuleAction{
						Monetize: &types.MonetizeAction{PriceMultiplier: aws.String("2")},
					},
					Statement:        actionURIStatement("/premium"),
					VisibilityConfig: wafActionRuleVisibility("MonetizeRule"),
				},
			},
		})
		return AssertErrorContains(err, "WAFInvalidParameterException")
	}))

	results = append(results, r.RunTest("wafv2", "Monetize_InvalidAmountFormat_Rejected", func() error {
		_, err := tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
			Name:  aws.String(tc.uniqueName("test-monetize-amount")),
			Scope: types.ScopeCloudfront,
			DefaultAction: &types.DefaultAction{
				Allow: &types.AllowAction{},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   false,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("monetize-amount-metric"),
			},
			MonetizationConfig: testMonetizationConfigWithAmount("1."),
			Rules: []types.Rule{
				{
					Name:     aws.String("MonetizeRule"),
					Priority: 1,
					Action: &types.RuleAction{
						Monetize: &types.MonetizeAction{PriceMultiplier: aws.String("2")},
					},
					Statement:        actionURIStatement("/premium"),
					VisibilityConfig: wafActionRuleVisibility("MonetizeRule"),
				},
			},
		})
		return AssertErrorContains(err, "WAFInvalidParameterException")
	}))

	results = append(results, r.RunTest("wafv2", "Monetize_WithCloudFrontConfig_Accepted", func() error {
		name := tc.uniqueName("test-monetize-acl")
		create, err := tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
			Name:  aws.String(name),
			Scope: types.ScopeCloudfront,
			DefaultAction: &types.DefaultAction{
				Allow: &types.AllowAction{},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   false,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("monetize-accepted-metric"),
			},
			MonetizationConfig: testMonetizationConfig(),
			Rules: []types.Rule{
				{
					Name:     aws.String("MonetizeRule"),
					Priority: 1,
					Action: &types.RuleAction{
						Monetize: &types.MonetizeAction{PriceMultiplier: aws.String("2")},
					},
					Statement:        actionURIStatement("/premium"),
					VisibilityConfig: wafActionRuleVisibility("MonetizeRule"),
				},
			},
		})
		if err != nil {
			return err
		}
		aclID := aws.ToString(create.Summary.Id)
		lockToken := aws.ToString(create.Summary.LockToken)
		defer tc.deleteWebACL(name, aclID, lockToken)

		resp, err := tc.client.GetWebACL(tc.ctx, &wafv2.GetWebACLInput{
			Name: aws.String(name), Scope: types.ScopeCloudfront, Id: aws.String(aclID),
		})
		if err != nil {
			return err
		}
		config := resp.WebACL.MonetizationConfig
		if config == nil || config.CryptoConfig == nil || len(config.CryptoConfig.PaymentNetworks) != 1 {
			return fmt.Errorf("MonetizationConfig did not round-trip: %+v", config)
		}
		network := config.CryptoConfig.PaymentNetworks[0]
		if network.Chain != types.BlockchainChainBaseSepolia ||
			aws.ToString(network.WalletAddress) != "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed" ||
			len(network.Prices) != 1 || aws.ToString(network.Prices[0].Amount) != "0.010" ||
			network.Prices[0].Currency != types.CryptoCurrencyUsdc {
			return fmt.Errorf("payment network did not round-trip: %+v", network)
		}
		if len(resp.WebACL.Rules) != 1 ||
			resp.WebACL.Rules[0].Action == nil || resp.WebACL.Rules[0].Action.Monetize == nil ||
			aws.ToString(resp.WebACL.Rules[0].Action.Monetize.PriceMultiplier) != "2" {
			return fmt.Errorf("Monetize rule action did not round-trip: %+v", resp.WebACL.Rules)
		}
		return nil
	}))

	return results
}

// wafActionRuleVisibility is the per-rule visibility configuration shared by
// the action tests.
func wafActionRuleVisibility(rule string) *types.VisibilityConfig {
	return &types.VisibilityConfig{
		SampledRequestsEnabled:   true,
		CloudWatchMetricsEnabled: false,
		MetricName:               aws.String(rule + "-metric"),
	}
}
