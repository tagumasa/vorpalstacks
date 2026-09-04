package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func (r *TestRunner) runRoute53VPCTests(tc *route53TestContext) []TestResult {
	var results []TestResult

	var testVPC1, testVPC2 string
	results = append(results, r.RunTest("route53", "VPCSetup", func() error {
		var err error
		testVPC1, err = tc.createTestVPC("10.201.0.0/16")
		if err != nil {
			return err
		}
		testVPC2, err = tc.createTestVPC("10.202.0.0/16")
		if err != nil {
			return err
		}
		return nil
	}))
	defer func() {
		if testVPC1 != "" {
			tc.deleteTestVPC(testVPC1)
		}
		if testVPC2 != "" {
			tc.deleteTestVPC(testVPC2)
		}
	}()

	results = append(results, r.RunTest("route53", "AssociateVPCWithHostedZone", func() error {
		if testVPC1 == "" || testVPC2 == "" {
			return fmt.Errorf("VPC setup failed")
		}
		privateDomain := tc.domain("private")
		createResp, err := tc.createPrivateZone(privateDomain, tc.callerRef("privref"), testVPC1, "private zone for VPC test")
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		privateZoneID := aws.ToString(createResp.HostedZone.Id)

		defer func() {
			tc.client.DisassociateVPCFromHostedZone(tc.ctx, &route53.DisassociateVPCFromHostedZoneInput{
				HostedZoneId: aws.String(privateZoneID),
				VPC: &types.VPC{
					VPCId:     aws.String(testVPC2),
					VPCRegion: types.VPCRegion(r.region),
				},
			})
			tc.client.DisassociateVPCFromHostedZone(tc.ctx, &route53.DisassociateVPCFromHostedZoneInput{
				HostedZoneId: aws.String(privateZoneID),
				VPC: &types.VPC{
					VPCId:     aws.String(testVPC1),
					VPCRegion: types.VPCRegion(r.region),
				},
			})
			tc.deleteZone(privateZoneID)
		}()

		assocResp, err := tc.client.AssociateVPCWithHostedZone(tc.ctx, &route53.AssociateVPCWithHostedZoneInput{
			HostedZoneId: aws.String(privateZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String(testVPC2),
				VPCRegion: types.VPCRegion(r.region),
			},
		})
		if err != nil {
			return fmt.Errorf("associate: %v", err)
		}
		if assocResp == nil || assocResp.ChangeInfo == nil {
			return fmt.Errorf("associate response or change info is nil")
		}

		getResp, err := tc.getZone(privateZoneID)
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
		if !vpcIDs[testVPC1] || !vpcIDs[testVPC2] {
			return fmt.Errorf("expected VPCs %s and %s, got %v", testVPC1, testVPC2, vpcIDs)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "DisassociateVPCFromHostedZone", func() error {
		if testVPC1 == "" || testVPC2 == "" {
			return fmt.Errorf("VPC setup failed")
		}
		dsDomain := tc.domain("disassoc")
		createResp, err := tc.createPrivateZone(dsDomain, tc.callerRef("dsref"), testVPC1, "")
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		dsZoneID := aws.ToString(createResp.HostedZone.Id)

		defer tc.deleteZone(dsZoneID)

		_, err = tc.client.AssociateVPCWithHostedZone(tc.ctx, &route53.AssociateVPCWithHostedZoneInput{
			HostedZoneId: aws.String(dsZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String(testVPC2),
				VPCRegion: types.VPCRegion(r.region),
			},
		})
		if err != nil {
			return fmt.Errorf("associate: %v", err)
		}

		disassocResp, err := tc.client.DisassociateVPCFromHostedZone(tc.ctx, &route53.DisassociateVPCFromHostedZoneInput{
			HostedZoneId: aws.String(dsZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String(testVPC2),
				VPCRegion: types.VPCRegion(r.region),
			},
		})
		if err != nil {
			return fmt.Errorf("disassociate: %v", err)
		}
		if disassocResp == nil || disassocResp.ChangeInfo == nil {
			return fmt.Errorf("disassociate response or change info is nil")
		}

		getResp, err := tc.getZone(dsZoneID)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(getResp.VPCs) != 1 {
			return fmt.Errorf("expected 1 VPC after disassociation, got %d", len(getResp.VPCs))
		}
		if aws.ToString(getResp.VPCs[0].VPCId) != testVPC1 {
			return fmt.Errorf("remaining VPC mismatch: got %q", aws.ToString(getResp.VPCs[0].VPCId))
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "DisassociateVPCFromHostedZone_NotAssociated", func() error {
		if testVPC1 == "" || testVPC2 == "" {
			return fmt.Errorf("VPC setup failed")
		}
		naDomain := tc.domain("notassoc")
		createResp, err := tc.createPrivateZone(naDomain, tc.callerRef("naref"), testVPC1, "")
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		naZoneID := aws.ToString(createResp.HostedZone.Id)

		defer tc.deleteZone(naZoneID)

		// testVPC2 was never associated with this zone; the removal is
		// rejected with a 404 VPCAssociationNotFound.
		_, err = tc.client.DisassociateVPCFromHostedZone(tc.ctx, &route53.DisassociateVPCFromHostedZoneInput{
			HostedZoneId: aws.String(naZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String(testVPC2),
				VPCRegion: types.VPCRegion(r.region),
			},
		})
		return expectRoute53Error(err, "VPCAssociationNotFound", 404)
	}))

	results = append(results, r.RunTest("route53", "AssociateVPCWithHostedZone_PublicZone", func() error {
		pubDomain := tc.domain("pub-vpc-test")
		createResp, err := tc.createZone(pubDomain, tc.callerRef("pubvpc"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		pubZoneID := aws.ToString(createResp.HostedZone.Id)

		defer tc.deleteZone(pubZoneID)

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
