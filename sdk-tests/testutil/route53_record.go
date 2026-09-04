package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func (r *TestRunner) runRoute53RecordTests(tc *route53TestContext) []TestResult {
	var results []TestResult

	domainName := tc.domain("record")
	zoneResp, err := tc.createZone(domainName, tc.callerRef("recref"))
	if err != nil {
		return []TestResult{{
			Service:  "route53",
			TestName: "RecordSetup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("CreateHostedZone for record tests: %v", err),
		}}
	}
	zoneID := aws.ToString(zoneResp.HostedZone.Id)
	defer tc.deleteZone(zoneID)

	var changeID string
	recordName := fmt.Sprintf("test.%s", domainName)
	results = append(results, r.RunTest("route53", "ChangeResourceRecordSets_Create", func() error {
		resp, err := tc.changeRecords(zoneID, rrChange(types.ChangeActionCreate, recordName, types.RRTypeA, 300, "192.0.2.1"))
		if err != nil {
			return err
		}
		if resp == nil || resp.ChangeInfo == nil {
			return fmt.Errorf("response or change info is nil")
		}
		if resp.ChangeInfo.Id == nil {
			return fmt.Errorf("change info ID is nil")
		}
		if resp.ChangeInfo.Status != types.ChangeStatusInsync {
			return fmt.Errorf("expected INSYNC status, got %v", resp.ChangeInfo.Status)
		}
		changeID = aws.ToString(resp.ChangeInfo.Id)
		return nil
	}))

	if changeID != "" {
		results = append(results, r.RunTest("route53", "GetChange", func() error {
			resp, err := tc.client.GetChange(tc.ctx, &route53.GetChangeInput{
				Id: aws.String(changeID),
			})
			if err != nil {
				return err
			}
			if resp.ChangeInfo == nil {
				return fmt.Errorf("change info is nil")
			}
			if aws.ToString(resp.ChangeInfo.Id) != changeID {
				return fmt.Errorf("change ID mismatch: got %q, want %q", aws.ToString(resp.ChangeInfo.Id), changeID)
			}
			return nil
		}))
	}

	results = append(results, r.RunTest("route53", "ListResourceRecordSets", func() error {
		resp, err := tc.listRecords(zoneID)
		if err != nil {
			return err
		}
		if resp.ResourceRecordSets == nil {
			return fmt.Errorf("resource record sets list is nil")
		}
		found := false
		for _, rs := range resp.ResourceRecordSets {
			if aws.ToString(rs.Name) == recordName && rs.Type == types.RRTypeA {
				found = true
				if aws.ToInt64(rs.TTL) != 300 {
					return fmt.Errorf("TTL mismatch: got %d, want 300", aws.ToInt64(rs.TTL))
				}
				if len(rs.ResourceRecords) != 1 || aws.ToString(rs.ResourceRecords[0].Value) != "192.0.2.1" {
					return fmt.Errorf("resource record value mismatch")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created A record not found in ListResourceRecordSets")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "DeleteResourceRecord", func() error {
		resp, err := tc.changeRecords(zoneID, rrChange(types.ChangeActionDelete, recordName, types.RRTypeA, 300, "192.0.2.1"))
		if err != nil {
			return err
		}
		if resp == nil || resp.ChangeInfo == nil {
			return fmt.Errorf("response or change info is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ChangeResourceRecordSets_Upsert", func() error {
		upsertDomain := tc.domain("upsert")
		createResp, err := tc.createZone(upsertDomain, tc.callerRef("upsertref"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		upsertZoneID := aws.ToString(createResp.HostedZone.Id)

		defer tc.deleteZone(upsertZoneID)

		upRecordName := fmt.Sprintf("upsert.%s", upsertDomain)
		_, err = tc.changeRecords(upsertZoneID, rrChange(types.ChangeActionUpsert, upRecordName, types.RRTypeA, 300, "10.0.0.1"))
		if err != nil {
			return fmt.Errorf("upsert create: %v", err)
		}

		_, err = tc.changeRecords(upsertZoneID, rrChange(types.ChangeActionUpsert, upRecordName, types.RRTypeA, 600, "10.0.0.2"))
		if err != nil {
			return fmt.Errorf("upsert update: %v", err)
		}

		listResp, err := tc.listRecords(upsertZoneID)
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		found := false
		for _, rs := range listResp.ResourceRecordSets {
			if aws.ToString(rs.Name) == upRecordName && rs.Type == types.RRTypeA {
				if aws.ToInt64(rs.TTL) != 600 {
					return fmt.Errorf("TTL mismatch after upsert: got %d, want 600", aws.ToInt64(rs.TTL))
				}
				if len(rs.ResourceRecords) != 1 {
					return fmt.Errorf("expected 1 resource record, got %d", len(rs.ResourceRecords))
				}
				if aws.ToString(rs.ResourceRecords[0].Value) != "10.0.0.2" {
					return fmt.Errorf("value mismatch after upsert: got %q", aws.ToString(rs.ResourceRecords[0].Value))
				}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("upserted record not found")
		}

		tc.changeRecords(upsertZoneID, rrChange(types.ChangeActionDelete, upRecordName, types.RRTypeA, 600, "10.0.0.2"))
		return nil
	}))

	results = append(results, r.RunTest("route53", "ChangeResourceRecordSets_MultipleTypes", func() error {
		mtDomain := tc.domain("multitype")
		createResp, err := tc.createZone(mtDomain, tc.callerRef("mtref"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		mtZoneID := aws.ToString(createResp.HostedZone.Id)

		defer tc.changeRecords(mtZoneID,
			rrChange(types.ChangeActionDelete, fmt.Sprintf("www.%s", mtDomain), types.RRTypeCname, 300, "target.example.com"),
			rrChange(types.ChangeActionDelete, fmt.Sprintf("txt.%s", mtDomain), types.RRTypeTxt, 300, "v=spf1 include:example.com ~all"),
		)
		defer tc.deleteZone(mtZoneID)

		_, err = tc.changeRecords(mtZoneID,
			rrChange(types.ChangeActionCreate, fmt.Sprintf("www.%s", mtDomain), types.RRTypeCname, 300, "target.example.com"),
			rrChange(types.ChangeActionCreate, fmt.Sprintf("txt.%s", mtDomain), types.RRTypeTxt, 300, "v=spf1 include:example.com ~all"),
		)
		if err != nil {
			return fmt.Errorf("create records: %v", err)
		}

		listResp, err := tc.listRecords(mtZoneID)
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}

		foundCNAME := false
		foundTXT := false
		for _, rs := range listResp.ResourceRecordSets {
			if rs.Type == types.RRTypeCname && strings.HasSuffix(aws.ToString(rs.Name), mtDomain) {
				foundCNAME = true
				if len(rs.ResourceRecords) != 1 || aws.ToString(rs.ResourceRecords[0].Value) != "target.example.com" {
					return fmt.Errorf("CNAME value mismatch")
				}
			}
			if rs.Type == types.RRTypeTxt && strings.HasSuffix(aws.ToString(rs.Name), mtDomain) {
				foundTXT = true
			}
		}
		if !foundCNAME {
			return fmt.Errorf("CNAME record not found")
		}
		if !foundTXT {
			return fmt.Errorf("TXT record not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "DeleteHostedZone_NotEmpty", func() error {
		neDomain := tc.domain("notempty")
		createResp, err := tc.createZone(neDomain, tc.callerRef("neref"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		neZoneID := aws.ToString(createResp.HostedZone.Id)

		neRecord := fmt.Sprintf("keep.%s", neDomain)
		_, err = tc.changeRecords(neZoneID, rrChange(types.ChangeActionCreate, neRecord, types.RRTypeA, 300, "10.0.0.1"))
		defer tc.deleteZone(neZoneID)
		defer tc.changeRecords(neZoneID, rrChange(types.ChangeActionDelete, neRecord, types.RRTypeA, 300, "10.0.0.1"))

		_, err = tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{
			Id: aws.String(neZoneID),
		})
		if err := AssertErrorContains(err, "HostedZoneNotEmpty"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListResourceRecordSets_Pagination", func() error {
		pgDomain := tc.domain("rrsetpg")
		cr, err := tc.createZone(pgDomain, tc.callerRef("rrsetpg"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		pgZoneID := aws.ToString(cr.HostedZone.Id)

		var pgRecordNames []string
		for i := 0; i < 5; i++ {
			pgRecordNames = append(pgRecordNames, fmt.Sprintf("rec%d.%s", i, pgDomain))
		}
		defer tc.deleteZone(pgZoneID)
		defer func() {
			for _, rn := range pgRecordNames {
				tc.changeRecords(pgZoneID, rrChange(types.ChangeActionDelete, rn, types.RRTypeA, 300, "10.0.0.0"))
			}
		}()

		var recordNames []string
		for i := 0; i < 5; i++ {
			rn := fmt.Sprintf("rec%d.%s", i, pgDomain)
			recordNames = append(recordNames, rn)
			_, err := tc.changeRecords(pgZoneID, rrChange(types.ChangeActionCreate, rn, types.RRTypeA, 300, fmt.Sprintf("10.0.0.%d", i)))
			if err != nil {
				return fmt.Errorf("create record %d: %v", i, err)
			}
		}

		var allNames []string
		var startName *string
		var startType types.RRType
		var hasStartType bool
		for {
			input := &route53.ListResourceRecordSetsInput{
				HostedZoneId: aws.String(pgZoneID),
				MaxItems:     aws.Int32(2),
			}
			if startName != nil {
				input.StartRecordName = startName
			}
			if hasStartType {
				input.StartRecordType = startType
			}
			resp, err := tc.client.ListResourceRecordSets(tc.ctx, input)
			if err != nil {
				return fmt.Errorf("list page: %v", err)
			}
			for _, rs := range resp.ResourceRecordSets {
				for _, rn := range recordNames {
					if aws.ToString(rs.Name) == rn {
						allNames = append(allNames, rn)
					}
				}
			}
			if !resp.IsTruncated || resp.NextRecordName == nil {
				break
			}
			startName = resp.NextRecordName
			startType = resp.NextRecordType
			hasStartType = true
		}
		if len(allNames) != 5 {
			return fmt.Errorf("expected 5 records across pages, got %d", len(allNames))
		}
		for _, rn := range recordNames {
			found := false
			for _, got := range allNames {
				if got == rn {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("record %q not found in paginated results", rn)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ChangeResourceRecordSets_DuplicateCreate", func() error {
		ddDomain := tc.domain("duprec")
		cr, err := tc.createZone(ddDomain, tc.callerRef("ddref"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		ddZoneID := aws.ToString(cr.HostedZone.Id)

		ddRecord := fmt.Sprintf("dup.%s", ddDomain)
		defer tc.deleteZone(ddZoneID)
		defer tc.changeRecords(ddZoneID, rrChange(types.ChangeActionDelete, ddRecord, types.RRTypeA, 300, "10.0.0.1"))

		_, err = tc.changeRecords(ddZoneID, rrChange(types.ChangeActionCreate, ddRecord, types.RRTypeA, 300, "10.0.0.1"))
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}

		_, err = tc.changeRecords(ddZoneID, rrChange(types.ChangeActionCreate, ddRecord, types.RRTypeA, 300, "10.0.0.2"))
		if err := AssertErrorContains(err, "InvalidChangeBatch"); err != nil {
			return fmt.Errorf("duplicate CREATE should fail with InvalidChangeBatch: %v", err)
		}
		return nil
	}))

	// Record-set changes are validated against the API model bounds before
	// being applied; both rejections below surface InvalidChangeBatch with
	// HTTP 400.
	results = append(results, r.RunTest("route53", "Record_InvalidBatchRejected", func() error {
		for _, c := range []struct {
			name string
			call func() error
		}{
			{
				// ResourceRecords carries @length(min 1) in the API model: a
				// non-alias record set with a present-but-empty value list
				// must be rejected instead of creating a record set with
				// zero values.
				name: "empty ResourceRecords list",
				call: func() error {
					_, err := tc.changeRecords(zoneID, rrChange(types.ChangeActionCreate, fmt.Sprintf("empty-rr.%s", domainName), types.RRTypeA, 300))
					return err
				},
			},
			{
				// ResourceRecord.Value is bounded by the RData shape
				// (@length 0..4000, counted in characters): an over-long
				// value invalidates the change batch instead of being
				// persisted.
				name: "over-long record value",
				call: func() error {
					_, err := tc.changeRecords(zoneID, rrChange(types.ChangeActionCreate, fmt.Sprintf("longvalue.%s", domainName), types.RRTypeTxt, 300, strings.Repeat("a", 4001)))
					return err
				},
			},
		} {
			if err := expectRoute53Error(c.call(), "InvalidChangeBatch", 400); err != nil {
				return fmt.Errorf("%s: %w", c.name, err)
			}
		}
		return nil
	}))

	return results
}
