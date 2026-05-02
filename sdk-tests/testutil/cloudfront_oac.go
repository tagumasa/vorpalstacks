package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func cfOACTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("cloudfront", "ListOriginAccessControls_VerifyFields", func() error {
		resp, err := client.ListOriginAccessControls(ctx, &cloudfront.ListOriginAccessControlsInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.OriginAccessControlList == nil {
			return fmt.Errorf("OAC list is nil")
		}
		if resp.OriginAccessControlList.MaxItems == nil {
			return fmt.Errorf("maxItems is nil")
		}
		return nil
	}))

	oacName := tc.uniquePrefix("test-oac")
	var oacID string
	var oacETag string
	results = append(results, tc.runner.RunTest("cloudfront", "CreateOriginAccessControl_VerifyFields", func() error {
		resp, err := client.CreateOriginAccessControl(ctx, &cloudfront.CreateOriginAccessControlInput{
			OriginAccessControlConfig: &types.OriginAccessControlConfig{
				Name:                          aws.String(oacName),
				OriginAccessControlOriginType: types.OriginAccessControlOriginTypesS3,
				SigningBehavior:               types.OriginAccessControlSigningBehaviorsNever,
				SigningProtocol:               types.OriginAccessControlSigningProtocolsSigv4,
			},
		})
		if err != nil {
			return err
		}
		if resp.OriginAccessControl == nil {
			return fmt.Errorf("OAC is nil")
		}
		if resp.OriginAccessControl.Id == nil || *resp.OriginAccessControl.Id == "" {
			return fmt.Errorf("OAC ID is nil or empty")
		}
		if resp.OriginAccessControl.OriginAccessControlConfig == nil {
			return fmt.Errorf("OAC config is nil")
		}
		cfg := resp.OriginAccessControl.OriginAccessControlConfig
		if cfg.Name == nil || *cfg.Name != oacName {
			return fmt.Errorf("OAC name mismatch: got %q, want %q", aws.ToString(cfg.Name), oacName)
		}
		if cfg.SigningBehavior != types.OriginAccessControlSigningBehaviorsNever {
			return fmt.Errorf("signing behavior mismatch: got %s", cfg.SigningBehavior)
		}
		if cfg.SigningProtocol != types.OriginAccessControlSigningProtocolsSigv4 {
			return fmt.Errorf("signing protocol mismatch: got %s", cfg.SigningProtocol)
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("ETag is nil or empty")
		}
		if resp.Location == nil || *resp.Location == "" {
			return fmt.Errorf("location is nil or empty")
		}
		oacID = *resp.OriginAccessControl.Id
		oacETag = *resp.ETag
		return nil
	}))

	if oacID != "" {
		results = append(results, tc.runner.RunTest("cloudfront", "GetOriginAccessControl_VerifyFields", func() error {
			resp, err := client.GetOriginAccessControl(ctx, &cloudfront.GetOriginAccessControlInput{
				Id: aws.String(oacID),
			})
			if err != nil {
				return err
			}
			if resp.OriginAccessControl == nil {
				return fmt.Errorf("OAC is nil")
			}
			if resp.OriginAccessControl.Id == nil || *resp.OriginAccessControl.Id != oacID {
				return fmt.Errorf("OAC ID mismatch: got %q, want %q", aws.ToString(resp.OriginAccessControl.Id), oacID)
			}
			if resp.ETag == nil || *resp.ETag == "" {
				return fmt.Errorf("ETag is nil or empty")
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "GetOriginAccessControlConfig_VerifyFields", func() error {
			resp, err := client.GetOriginAccessControlConfig(ctx, &cloudfront.GetOriginAccessControlConfigInput{
				Id: aws.String(oacID),
			})
			if err != nil {
				return err
			}
			if resp.OriginAccessControlConfig == nil {
				return fmt.Errorf("OAC config is nil")
			}
			if resp.OriginAccessControlConfig.Name == nil || *resp.OriginAccessControlConfig.Name != oacName {
				return fmt.Errorf("OAC config name mismatch: got %q, want %q", aws.ToString(resp.OriginAccessControlConfig.Name), oacName)
			}
			if resp.ETag == nil || *resp.ETag == "" {
				return fmt.Errorf("ETag is nil or empty")
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "UpdateOriginAccessControl_VerifyFields", func() error {
			updatedName := oacName + "-updated"
			resp, err := client.UpdateOriginAccessControl(ctx, &cloudfront.UpdateOriginAccessControlInput{
				Id:      aws.String(oacID),
				IfMatch: aws.String(oacETag),
				OriginAccessControlConfig: &types.OriginAccessControlConfig{
					Name:                          aws.String(updatedName),
					OriginAccessControlOriginType: types.OriginAccessControlOriginTypesS3,
					SigningBehavior:               types.OriginAccessControlSigningBehaviorsAlways,
					SigningProtocol:               types.OriginAccessControlSigningProtocolsSigv4,
					Description:                   aws.String("updated description"),
				},
			})
			if err != nil {
				return err
			}
			if resp.OriginAccessControl == nil {
				return fmt.Errorf("updated OAC is nil")
			}
			if resp.OriginAccessControl.OriginAccessControlConfig == nil {
				return fmt.Errorf("OAC config is nil")
			}
			cfg := resp.OriginAccessControl.OriginAccessControlConfig
			if cfg.Name == nil || *cfg.Name != updatedName {
				return fmt.Errorf("OAC name not updated: got %q, want %q", aws.ToString(cfg.Name), updatedName)
			}
			if cfg.SigningBehavior != types.OriginAccessControlSigningBehaviorsAlways {
				return fmt.Errorf("signing behavior not updated: got %s", cfg.SigningBehavior)
			}
			if cfg.Description == nil || *cfg.Description != "updated description" {
				return fmt.Errorf("description not updated: got %q", aws.ToString(cfg.Description))
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "DeleteOriginAccessControl", func() error {
			_, err := client.DeleteOriginAccessControl(ctx, &cloudfront.DeleteOriginAccessControlInput{
				Id: aws.String(oacID),
			})
			return err
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "GetOriginAccessControl_NonExistent", func() error {
			_, err := client.GetOriginAccessControl(ctx, &cloudfront.GetOriginAccessControlInput{
				Id: aws.String(oacID),
			})
			return AssertErrorContains(err, "NoSuchOriginAccessControl")
		}))
	}

	return results
}
