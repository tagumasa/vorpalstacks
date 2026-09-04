package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func (r *TestRunner) runRoute53ZoneTests(tc *route53TestContext) []TestResult {
	var results []TestResult

	domainName := tc.domain("example")

	var hostedZoneID string
	results = append(results, r.RunTest("route53", "CreateHostedZone", func() error {
		resp, err := tc.createZone(domainName, tc.callerRef("ref"))
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
		if resp.ChangeInfo == nil || resp.ChangeInfo.Id == nil {
			return fmt.Errorf("change info missing in create response")
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
			resp, err := tc.getZone(hostedZoneID)
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
		hzResp, err := tc.createZone(testDomain, tc.callerRef("sortref"))
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
		tc.deleteZone(aws.ToString(hzResp.HostedZone.Id))
		return nil
	}))

	results = append(results, r.RunTest("route53", "UpdateHostedZoneComment", func() error {
		ucDomain := tc.domain("updatecomment")
		createResp, err := tc.createZone(ucDomain, tc.callerRef("ucref"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		ucID := aws.ToString(createResp.HostedZone.Id)

		defer tc.deleteZone(ucID)

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

		getResp, err := tc.getZone(ucID)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.HostedZone.Config == nil || aws.ToString(getResp.HostedZone.Config.Comment) != comment {
			return fmt.Errorf("comment mismatch: got %q, want %q", aws.ToString(getResp.HostedZone.Config.Comment), comment)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "CreateHostedZone_PrivateWithComment", func() error {
		pvtVPCID, err := tc.createTestVPC("10.210.0.0/16")
		if err != nil {
			return fmt.Errorf("vpc setup: %v", err)
		}
		defer tc.deleteTestVPC(pvtVPCID)

		pvtDomain := tc.domain("private-comment")
		resp, err := tc.createPrivateZone(pvtDomain, tc.callerRef("pvtref"), pvtVPCID, "private zone with comment")
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		pvtID := aws.ToString(resp.HostedZone.Id)

		defer tc.deleteZone(pvtID)

		getResp, err := tc.getZone(pvtID)
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
		if aws.ToString(getResp.VPCs[0].VPCId) != pvtVPCID {
			return fmt.Errorf("VPC ID mismatch: got %q", aws.ToString(getResp.VPCs[0].VPCId))
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "DelegationSet_Persisted", func() error {
		dsDomain := tc.domain("ds-persist")
		createResp, err := tc.createZone(dsDomain, tc.callerRef("dspersist"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		dsZoneID := aws.ToString(createResp.HostedZone.Id)

		defer tc.deleteZone(dsZoneID)

		createNS := createResp.DelegationSet.NameServers
		if len(createNS) == 0 {
			return fmt.Errorf("name servers empty in create response")
		}

		getResp, err := tc.getZone(dsZoneID)
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
		cr, err := tc.createZone(domain, tc.callerRef("dnssec"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		hzID := aws.ToString(cr.HostedZone.Id)
		defer tc.deleteZone(hzID)

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
		dupRef := tc.callerRef("dupref")
		resp1, err := tc.createZone(dupDomain, dupRef)
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		hzID1 := aws.ToString(resp1.HostedZone.Id)
		defer tc.deleteZone(hzID1)

		resp2, err := tc.createZone(dupDomain, dupRef)
		if err != nil {
			return fmt.Errorf("duplicate caller ref: %v", err)
		}
		hzID2 := aws.ToString(resp2.HostedZone.Id)
		if hzID1 != hzID2 {
			return fmt.Errorf("idempotent create returned different ID: %q vs %q", hzID1, hzID2)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "CreateHostedZone_CallerReferenceRejected", func() error {
		for _, c := range []struct {
			name   string
			domain string
			ref    string
		}{
			{name: "empty caller reference", domain: tc.domain("emptyref"), ref: ""},
			{name: "caller reference over the 128-character bound", domain: tc.domain("longref"), ref: strings.Repeat("a", 129)},
		} {
			_, err := tc.createZone(c.domain, c.ref)
			if err := AssertErrorContains(err, "InvalidInput"); err != nil {
				return fmt.Errorf("%s: %w", c.name, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListHostedZones_PrivateHostedZoneFilter", func() error {
		pvtVPCID, err := tc.createTestVPC("10.220.0.0/16")
		if err != nil {
			return fmt.Errorf("vpc setup: %v", err)
		}
		defer tc.deleteTestVPC(pvtVPCID)

		pvtDomain := tc.domain("pvt-filter")
		pvtResp, err := tc.createPrivateZone(pvtDomain, tc.callerRef("pvtfilter"), pvtVPCID, "")
		if err != nil {
			return fmt.Errorf("create private zone: %v", err)
		}
		pvtZoneID := aws.ToString(pvtResp.HostedZone.Id)
		defer tc.deleteZone(pvtZoneID)

		pubDomain := tc.domain("pub-filter")
		pubResp, err := tc.createZone(pubDomain, tc.callerRef("pubfilter"))
		if err != nil {
			return fmt.Errorf("create public zone: %v", err)
		}
		pubZoneID := aws.ToString(pubResp.HostedZone.Id)
		defer tc.deleteZone(pubZoneID)

		var marker *string
		pvtFound := false
		pubFound := false
		for {
			resp, err := tc.client.ListHostedZones(tc.ctx, &route53.ListHostedZonesInput{
				HostedZoneType: types.HostedZoneTypePrivateHostedZone,
				Marker:         marker,
				MaxItems:       aws.Int32(100),
			})
			if err != nil {
				return fmt.Errorf("list with HostedZoneType filter: %v", err)
			}
			for _, hz := range resp.HostedZones {
				if aws.ToString(hz.Id) == pvtZoneID {
					pvtFound = true
				}
				if aws.ToString(hz.Id) == pubZoneID {
					pubFound = true
				}
			}
			if !resp.IsTruncated || resp.NextMarker == nil {
				break
			}
			marker = resp.NextMarker
		}
		if !pvtFound {
			return fmt.Errorf("private zone %s not found in PrivateHostedZone-filtered list", pvtZoneID)
		}
		if pubFound {
			return fmt.Errorf("public zone %s must not appear in PrivateHostedZone-filtered list", pubZoneID)
		}
		return nil
	}))

	// The ResourceDescription comment bound is a Smithy @length trait
	// counted in Unicode characters; a 128-character CJK comment is 384
	// bytes but must be accepted.
	results = append(results, r.RunTest("route53", "CreateHostedZone_CommentMultibyteAccepted", func() error {
		mbDomain := tc.domain("mb-comment")
		mbComment := strings.Repeat("\u65e5", 128)
		resp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
			Name:             aws.String(mbDomain),
			CallerReference:  aws.String(tc.callerRef("mb-comment")),
			HostedZoneConfig: &types.HostedZoneConfig{Comment: aws.String(mbComment)},
		})
		if err != nil {
			return fmt.Errorf("create with multibyte comment: %v", err)
		}
		hzID := aws.ToString(resp.HostedZone.Id)
		defer tc.deleteZone(hzID)

		getResp, err := tc.getZone(hzID)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.HostedZone.Config == nil || aws.ToString(getResp.HostedZone.Config.Comment) != mbComment {
			return fmt.Errorf("multibyte comment not persisted faithfully")
		}
		return nil
	}))

	return results
}
