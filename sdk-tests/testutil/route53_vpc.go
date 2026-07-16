package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func (r *TestRunner) runRoute53VPCTests(tc *route53TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("route53", "AssociateVPCWithHostedZone", func() error {
		privateDomain := tc.domain("private")
		privateRef := fmt.Sprintf("privref-%d", tc.uniq)
		createResp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(privateDomain),
			CallerReference: aws.String(privateRef),
			HostedZoneConfig: &types.HostedZoneConfig{
				PrivateZone: true,
				Comment:     aws.String("private zone for VPC test"),
			},
			VPC: &types.VPC{
				VPCId:     aws.String("vpc-abcdef01"),
				VPCRegion: types.VPCRegionUsEast1,
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		privateZoneID := aws.ToString(createResp.HostedZone.Id)

		defer func() {
			tc.client.DisassociateVPCFromHostedZone(tc.ctx, &route53.DisassociateVPCFromHostedZoneInput{
				HostedZoneId: aws.String(privateZoneID),
				VPC: &types.VPC{
					VPCId:     aws.String("vpc-xyz12345"),
					VPCRegion: types.VPCRegionUsEast1,
				},
			})
			tc.client.DisassociateVPCFromHostedZone(tc.ctx, &route53.DisassociateVPCFromHostedZoneInput{
				HostedZoneId: aws.String(privateZoneID),
				VPC: &types.VPC{
					VPCId:     aws.String("vpc-abcdef01"),
					VPCRegion: types.VPCRegionUsEast1,
				},
			})
			tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(privateZoneID)})
		}()

		assocResp, err := tc.client.AssociateVPCWithHostedZone(tc.ctx, &route53.AssociateVPCWithHostedZoneInput{
			HostedZoneId: aws.String(privateZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String("vpc-xyz12345"),
				VPCRegion: types.VPCRegionUsEast1,
			},
		})
		if err != nil {
			return fmt.Errorf("associate: %v", err)
		}
		if assocResp == nil || assocResp.ChangeInfo == nil {
			return fmt.Errorf("associate response or change info is nil")
		}

		getResp, err := tc.client.GetHostedZone(tc.ctx, &route53.GetHostedZoneInput{
			Id: aws.String(privateZoneID),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(getResp.VPCs) < 2 {
			return fmt.Errorf("expected at least 2 VPCs, got %d", len(getResp.VPCs))
		}
		vpcIDs := make(map[string]bool)
		for _, v := range getResp.VPCs {
			vpcIDs[aws.ToString(v.VPCId)] = true
		}
		if !vpcIDs["vpc-abcdef01"] || !vpcIDs["vpc-xyz12345"] {
			return fmt.Errorf("expected VPCs vpc-abcdef01 and vpc-xyz12345, got %v", vpcIDs)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "DisassociateVPCFromHostedZone", func() error {
		dsDomain := tc.domain("disassoc")
		dsRef := fmt.Sprintf("dsref-%d", tc.uniq)
		createResp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(dsDomain),
			CallerReference: aws.String(dsRef),
			HostedZoneConfig: &types.HostedZoneConfig{
				PrivateZone: true,
			},
			VPC: &types.VPC{
				VPCId:     aws.String("vpc-disassoc1"),
				VPCRegion: types.VPCRegionUsEast1,
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		dsZoneID := aws.ToString(createResp.HostedZone.Id)

		defer tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(dsZoneID)})

		_, err = tc.client.AssociateVPCWithHostedZone(tc.ctx, &route53.AssociateVPCWithHostedZoneInput{
			HostedZoneId: aws.String(dsZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String("vpc-disassoc2"),
				VPCRegion: types.VPCRegionUsEast1,
			},
		})
		if err != nil {
			return fmt.Errorf("associate: %v", err)
		}

		disassocResp, err := tc.client.DisassociateVPCFromHostedZone(tc.ctx, &route53.DisassociateVPCFromHostedZoneInput{
			HostedZoneId: aws.String(dsZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String("vpc-disassoc2"),
				VPCRegion: types.VPCRegionUsEast1,
			},
		})
		if err != nil {
			return fmt.Errorf("disassociate: %v", err)
		}
		if disassocResp == nil || disassocResp.ChangeInfo == nil {
			return fmt.Errorf("disassociate response or change info is nil")
		}

		getResp, err := tc.client.GetHostedZone(tc.ctx, &route53.GetHostedZoneInput{
			Id: aws.String(dsZoneID),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(getResp.VPCs) != 1 {
			return fmt.Errorf("expected 1 VPC after disassociation, got %d", len(getResp.VPCs))
		}
		if aws.ToString(getResp.VPCs[0].VPCId) != "vpc-disassoc1" {
			return fmt.Errorf("remaining VPC mismatch: got %q", aws.ToString(getResp.VPCs[0].VPCId))
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "AssociateVPCWithHostedZone_PublicZone", func() error {
		pubDomain := tc.domain("pub-vpc-test")
		pubRef := fmt.Sprintf("pubvpc-%d", tc.uniq)
		createResp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(pubDomain),
			CallerReference: aws.String(pubRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		pubZoneID := aws.ToString(createResp.HostedZone.Id)

		defer tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(pubZoneID)})

		_, err = tc.client.AssociateVPCWithHostedZone(tc.ctx, &route53.AssociateVPCWithHostedZoneInput{
			HostedZoneId: aws.String(pubZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String("vpc-test123"),
				VPCRegion: types.VPCRegionUsEast1,
			},
		})
		if err := AssertErrorContains(err, "InvalidInput"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
