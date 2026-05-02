package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func (r *TestRunner) runRoute53ZoneTests(tc *route53TestContext) []TestResult {
	var results []TestResult

	domainName := tc.domain("example")

	var hostedZoneID string
	results = append(results, r.RunTest("route53", "CreateHostedZone", func() error {
		resp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(domainName),
			CallerReference: aws.String(tc.callerRef("ref")),
		})
		if err != nil {
			return err
		}
		if resp.HostedZone == nil {
			return fmt.Errorf("hosted zone is nil")
		}
		if resp.HostedZone.Id == nil {
			return fmt.Errorf("hosted zone ID is nil")
		}
		if aws.ToString(resp.HostedZone.Name) != domainName {
			return fmt.Errorf("domain name mismatch: got %q, want %q", aws.ToString(resp.HostedZone.Name), domainName)
		}
		if resp.DelegationSet == nil || len(resp.DelegationSet.NameServers) == 0 {
			return fmt.Errorf("delegation set or name servers missing")
		}
		hostedZoneID = aws.ToString(resp.HostedZone.Id)
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListHostedZones", func() error {
		var marker *string
		found := false
		for {
			resp, err := tc.client.ListHostedZones(tc.ctx, &route53.ListHostedZonesInput{
				Marker:   marker,
				MaxItems: aws.Int32(10),
			})
			if err != nil {
				return err
			}
			if resp.HostedZones == nil {
				return fmt.Errorf("hosted zones list is nil")
			}
			for _, hz := range resp.HostedZones {
				if aws.ToString(hz.Id) == hostedZoneID {
					found = true
					if aws.ToString(hz.Name) != domainName {
						return fmt.Errorf("domain name mismatch in list: got %q", aws.ToString(hz.Name))
					}
					break
				}
			}
			if found {
				break
			}
			if !resp.IsTruncated || resp.NextMarker == nil {
				break
			}
			marker = resp.NextMarker
		}
		if !found {
			return fmt.Errorf("created hosted zone %s not found in ListHostedZones", hostedZoneID)
		}
		return nil
	}))

	if hostedZoneID != "" {
		results = append(results, r.RunTest("route53", "GetHostedZone", func() error {
			resp, err := tc.client.GetHostedZone(tc.ctx, &route53.GetHostedZoneInput{
				Id: aws.String(hostedZoneID),
			})
			if err != nil {
				return err
			}
			if resp.HostedZone == nil {
				return fmt.Errorf("hosted zone is nil")
			}
			if aws.ToString(resp.HostedZone.Id) != hostedZoneID {
				return fmt.Errorf("ID mismatch: got %q, want %q", aws.ToString(resp.HostedZone.Id), hostedZoneID)
			}
			if aws.ToString(resp.HostedZone.Name) != domainName {
				return fmt.Errorf("name mismatch: got %q, want %q", aws.ToString(resp.HostedZone.Name), domainName)
			}
			if resp.DelegationSet == nil || len(resp.DelegationSet.NameServers) == 0 {
				return fmt.Errorf("delegation set missing in GetHostedZone response")
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "DeleteHostedZone", func() error {
			resp, err := tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{
				Id: aws.String(hostedZoneID),
			})
			if err != nil {
				return err
			}
			if resp == nil || resp.ChangeInfo == nil {
				return fmt.Errorf("response or change info is nil")
			}
			if resp.ChangeInfo.Id == nil {
				return fmt.Errorf("change info ID is nil")
			}
			return nil
		}))
	}

	results = append(results, r.RunTest("route53", "CreateHostedZone_ContentVerify", func() error {
		verifyDomain := tc.domain("verify")
		verifyRef := fmt.Sprintf("ref-%d", tc.uniq)
		resp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(verifyDomain),
			CallerReference: aws.String(verifyRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		hzID := aws.ToString(resp.HostedZone.Id)

		defer tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(hzID)})

		if resp.HostedZone.Name == nil || *resp.HostedZone.Name != verifyDomain {
			return fmt.Errorf("domain name mismatch: got %q, want %q", aws.ToString(resp.HostedZone.Name), verifyDomain)
		}
		if resp.HostedZone.Id == nil || *resp.HostedZone.Id == "" {
			return fmt.Errorf("hosted zone ID is empty")
		}
		if resp.ChangeInfo == nil || resp.ChangeInfo.Id == nil {
			return fmt.Errorf("change info missing in create response")
		}

		getResp, err := tc.client.GetHostedZone(tc.ctx, &route53.GetHostedZoneInput{Id: aws.String(hzID)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.HostedZone.Name == nil || *getResp.HostedZone.Name != verifyDomain {
			return fmt.Errorf("get domain name mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListHostedZonesByName", func() error {
		resp, err := tc.client.ListHostedZonesByName(tc.ctx, &route53.ListHostedZonesByNameInput{
			MaxItems: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.HostedZones == nil {
			return fmt.Errorf("hosted zones list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListHostedZonesByName_WithDNSName", func() error {
		testDomain := tc.domain("sorttest")
		hzResp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(testDomain),
			CallerReference: aws.String(tc.callerRef("sortref")),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}

		resp, err := tc.client.ListHostedZonesByName(tc.ctx, &route53.ListHostedZonesByNameInput{
			DNSName:  aws.String(testDomain),
			MaxItems: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.HostedZones == nil {
			return fmt.Errorf("hosted zones list is nil")
		}
		found := false
		for _, hz := range resp.HostedZones {
			if aws.ToString(hz.Name) == testDomain {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created zone %q not found in ListHostedZonesByName", testDomain)
		}
		tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{
			Id: hzResp.HostedZone.Id,
		})
		return nil
	}))

	results = append(results, r.RunTest("route53", "UpdateHostedZoneComment", func() error {
		ucDomain := tc.domain("updatecomment")
		ucRef := fmt.Sprintf("ucref-%d", tc.uniq)
		createResp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(ucDomain),
			CallerReference: aws.String(ucRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		ucID := aws.ToString(createResp.HostedZone.Id)

		defer tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(ucID)})

		comment := "test comment for zone"
		updateResp, err := tc.client.UpdateHostedZoneComment(tc.ctx, &route53.UpdateHostedZoneCommentInput{
			Id:      aws.String(ucID),
			Comment: aws.String(comment),
		})
		if err != nil {
			return fmt.Errorf("update comment: %v", err)
		}
		if updateResp.HostedZone == nil {
			return fmt.Errorf("update response hosted zone is nil")
		}

		getResp, err := tc.client.GetHostedZone(tc.ctx, &route53.GetHostedZoneInput{Id: aws.String(ucID)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.HostedZone.Config == nil || aws.ToString(getResp.HostedZone.Config.Comment) != comment {
			return fmt.Errorf("comment mismatch: got %q, want %q", aws.ToString(getResp.HostedZone.Config.Comment), comment)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "CreateHostedZone_PrivateWithComment", func() error {
		pvtDomain := tc.domain("private-comment")
		pvtRef := fmt.Sprintf("pvtref-%d", tc.uniq)
		resp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(pvtDomain),
			CallerReference: aws.String(pvtRef),
			HostedZoneConfig: &types.HostedZoneConfig{
				PrivateZone: true,
				Comment:     aws.String("private zone with comment"),
			},
			VPC: &types.VPC{
				VPCId:     aws.String("vpc-pvttest"),
				VPCRegion: types.VPCRegionEuWest1,
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		pvtID := aws.ToString(resp.HostedZone.Id)

		defer tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(pvtID)})

		getResp, err := tc.client.GetHostedZone(tc.ctx, &route53.GetHostedZoneInput{Id: aws.String(pvtID)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.HostedZone.Config == nil {
			return fmt.Errorf("config is nil")
		}
		if !getResp.HostedZone.Config.PrivateZone {
			return fmt.Errorf("expected PrivateZone=true")
		}
		if aws.ToString(getResp.HostedZone.Config.Comment) != "private zone with comment" {
			return fmt.Errorf("comment mismatch: got %q", aws.ToString(getResp.HostedZone.Config.Comment))
		}
		if len(getResp.VPCs) != 1 {
			return fmt.Errorf("expected 1 VPC, got %d", len(getResp.VPCs))
		}
		if aws.ToString(getResp.VPCs[0].VPCId) != "vpc-pvttest" {
			return fmt.Errorf("VPC ID mismatch: got %q", aws.ToString(getResp.VPCs[0].VPCId))
		}
		if getResp.VPCs[0].VPCRegion != types.VPCRegionEuWest1 {
			return fmt.Errorf("VPC region mismatch: got %v", getResp.VPCs[0].VPCRegion)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "DelegationSet_Persisted", func() error {
		dsDomain := tc.domain("ds-persist")
		dsRef := fmt.Sprintf("dspersist-%d", tc.uniq)
		createResp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(dsDomain),
			CallerReference: aws.String(dsRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		dsZoneID := aws.ToString(createResp.HostedZone.Id)

		defer tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(dsZoneID)})

		createNS := createResp.DelegationSet.NameServers
		if len(createNS) == 0 {
			return fmt.Errorf("name servers empty in create response")
		}

		getResp, err := tc.client.GetHostedZone(tc.ctx, &route53.GetHostedZoneInput{Id: aws.String(dsZoneID)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		getNS := getResp.DelegationSet.NameServers
		if len(getNS) == 0 {
			return fmt.Errorf("name servers empty in get response")
		}
		if len(createNS) != len(getNS) {
			return fmt.Errorf("name server count mismatch: create=%d, get=%d", len(createNS), len(getNS))
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListReusableDelegationSets", func() error {
		resp, err := tc.client.ListReusableDelegationSets(tc.ctx, &route53.ListReusableDelegationSetsInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.DelegationSets == nil {
			return fmt.Errorf("delegation sets is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "GetDNSSEC", func() error {
		domain := tc.domain("dnssectest")
		cr, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(domain),
			CallerReference: aws.String(tc.callerRef("dnssec")),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		hzID := aws.ToString(cr.HostedZone.Id)
		defer tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(hzID)})

		resp, err := tc.client.GetDNSSEC(tc.ctx, &route53.GetDNSSECInput{
			HostedZoneId: aws.String(hzID),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "CreateHostedZone_DuplicateCallerRef", func() error {
		dupDomain := tc.domain("dupref")
		dupRef := fmt.Sprintf("dupref-%d", tc.uniq)
		resp1, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(dupDomain),
			CallerReference: aws.String(dupRef),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		hzID1 := aws.ToString(resp1.HostedZone.Id)
		defer tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(hzID1)})

		resp2, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(dupDomain),
			CallerReference: aws.String(dupRef),
		})
		if err != nil {
			return fmt.Errorf("duplicate caller ref: %v", err)
		}
		hzID2 := aws.ToString(resp2.HostedZone.Id)
		if hzID1 != hzID2 {
			return fmt.Errorf("idempotent create returned different ID: %q vs %q", hzID1, hzID2)
		}
		return nil
	}))

	return results
}
