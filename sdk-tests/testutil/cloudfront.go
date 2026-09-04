package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"vorpalstacks-sdk-tests/config"
)

type cfTestContext struct {
	runner *TestRunner
	client *cloudfront.Client
	ctx    context.Context
}

func (tc *cfTestContext) uniquePrefix(tag string) string {
	return fmt.Sprintf("%s-%d", tag, time.Now().UnixNano())
}

// baseDistributionConfig returns the minimal valid distribution
// configuration shared by every test that creates a distribution: one
// custom origin, a default cache behaviour forwarding nothing, the
// default viewer certificate, and no geo restriction. Callers mutate
// the returned value to add the behaviour under test.
func (tc *cfTestContext) baseDistributionConfig(callerRef, comment, originDomain string) *types.DistributionConfig {
	originID := tc.uniquePrefix("origin")
	return &types.DistributionConfig{
		CallerReference: aws.String(callerRef),
		Enabled:         aws.Bool(true),
		Comment:         aws.String(comment),
		Origins: &types.Origins{
			Quantity: aws.Int32(1),
			Items: []types.Origin{
				{
					Id:         aws.String(originID),
					DomainName: aws.String(originDomain),
					CustomOriginConfig: &types.CustomOriginConfig{
						HTTPPort:             aws.Int32(80),
						HTTPSPort:            aws.Int32(443),
						OriginProtocolPolicy: types.OriginProtocolPolicyHttpOnly,
					},
				},
			},
		},
		DefaultCacheBehavior: &types.DefaultCacheBehavior{
			TargetOriginId:       aws.String(originID),
			ViewerProtocolPolicy: types.ViewerProtocolPolicyAllowAll,
			ForwardedValues: &types.ForwardedValues{
				QueryString: aws.Bool(false),
				Cookies: &types.CookiePreference{
					Forward: types.ItemSelectionNone,
				},
			},
		},
		ViewerCertificate: &types.ViewerCertificate{
			CloudFrontDefaultCertificate: aws.Bool(true),
		},
		Restrictions: &types.Restrictions{
			GeoRestriction: &types.GeoRestriction{
				RestrictionType: types.GeoRestrictionTypeNone,
				Quantity:        aws.Int32(0),
			},
		},
	}
}

// disableAndDeleteDistributionByID tears down a distribution without
// requiring callers to track its ETag: the fetched config supplies the
// If-Match for the disable update and the update response supplies the
// one for the delete.
func (tc *cfTestContext) disableAndDeleteDistributionByID(id string) error {
	cfgResp, err := tc.client.GetDistributionConfig(tc.ctx, &cloudfront.GetDistributionConfigInput{
		Id: aws.String(id),
	})
	if err != nil {
		return err
	}
	cfgResp.DistributionConfig.Enabled = aws.Bool(false)
	updResp, err := tc.client.UpdateDistribution(tc.ctx, &cloudfront.UpdateDistributionInput{
		Id:                 aws.String(id),
		IfMatch:            cfgResp.ETag,
		DistributionConfig: cfgResp.DistributionConfig,
	})
	if err != nil {
		return err
	}
	_, err = tc.client.DeleteDistribution(tc.ctx, &cloudfront.DeleteDistributionInput{
		Id:      aws.String(id),
		IfMatch: updResp.ETag,
	})
	return err
}

// createPublicKey creates a throwaway public key for key-group tests;
// callers delete it with the returned ETag once the key group is gone.
func (tc *cfTestContext) createPublicKey(tag string) (id, etag string, err error) {
	resp, err := tc.client.CreatePublicKey(tc.ctx, &cloudfront.CreatePublicKeyInput{
		PublicKeyConfig: &types.PublicKeyConfig{
			Name:            aws.String(tc.uniquePrefix(tag)),
			EncodedKey:      aws.String(cloudfrontEncodedKeyB64),
			CallerReference: aws.String(tc.uniquePrefix(tag + "-ref")),
		},
	})
	if err != nil {
		return "", "", err
	}
	return aws.ToString(resp.PublicKey.Id), aws.ToString(resp.ETag), nil
}

func (r *TestRunner) RunCloudFrontTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudfront",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	tc := &cfTestContext{
		runner: r,
		client: cloudfront.NewFromConfig(cfg),
		ctx:    context.Background(),
	}

	results = append(results, cfDistributionTests(tc)...)
	results = append(results, cfTagTests(tc)...)
	results = append(results, cfInvalidationTests(tc)...)
	results = append(results, cfOACTests(tc)...)
	results = append(results, cfCachePolicyTests(tc)...)
	results = append(results, cfOriginRequestPolicyTests(tc)...)
	results = append(results, cfResponseHeadersPolicyTests(tc)...)
	results = append(results, cfWebACLCopyTests(tc)...)
	results = append(results, cfKeyGroupTests(tc)...)
	results = append(results, cfContinuousDeploymentTests(tc)...)
	results = append(results, cfEdgeTests(tc)...)

	return results
}
