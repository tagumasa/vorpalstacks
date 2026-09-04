package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func (r *TestRunner) runRoute53EdgeTests(tc *route53TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("route53", "NonExistentResources", func() error {
		for _, c := range []struct {
			name string
			call func() error
			code string
		}{
			{
				name: "GetHostedZone on a non-existent zone ID",
				call: func() error {
					_, err := tc.getZone("Z00000000000000000000")
					return err
				},
				code: "NoSuchHostedZone",
			},
			{
				name: "DeleteHostedZone on a non-existent zone ID",
				call: func() error {
					_, err := tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{
						Id: aws.String("Z00000000000000000000"),
					})
					return err
				},
				code: "NoSuchHostedZone",
			},
			{
				name: "GetChange on a non-existent change ID",
				call: func() error {
					_, err := tc.client.GetChange(tc.ctx, &route53.GetChangeInput{
						Id: aws.String("C0000000000000000000000000"),
					})
					return err
				},
				code: "NoSuchChange",
			},
		} {
			if err := AssertErrorContains(c.call(), c.code); err != nil {
				return fmt.Errorf("%s: %w", c.name, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "CreateHostedZone_InvalidName", func() error {
		_, err := tc.createZone("invalid name with spaces", tc.callerRef("badref"))
		if err := AssertErrorContains(err, "InvalidDomainName"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListHostedZones_Pagination", func() error {
		pgTs := fmt.Sprintf("pzpg%d", tc.uniq)
		var pgZoneIDs []string
		for i := 0; i < 5; i++ {
			resp, err := tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
				Name:            aws.String(fmt.Sprintf("%s-%d.example.com.", pgTs, i)),
				CallerReference: aws.String(fmt.Sprintf("%s-ref-%d", pgTs, i)),
			})
			if err != nil {
				return fmt.Errorf("create hosted zone: %v", err)
			}
			pgZoneIDs = append(pgZoneIDs, aws.ToString(resp.HostedZone.Id))
		}

		pageCount := 0
		totalCount := 0
		var marker *string
		for {
			resp, err := tc.client.ListHostedZones(tc.ctx, &route53.ListHostedZonesInput{
				Marker:   marker,
				MaxItems: aws.Int32(2),
			})
			if err != nil {
				for _, zid := range pgZoneIDs {
					tc.deleteZone(zid)
				}
				return fmt.Errorf("list hosted zones page: %v", err)
			}
			pageCount++
			totalCount += len(resp.HostedZones)
			if resp.IsTruncated && resp.NextMarker != nil {
				marker = resp.NextMarker
			} else {
				break
			}
		}

		for _, zid := range pgZoneIDs {
			tc.deleteZone(zid)
		}

		if pageCount < 2 {
			return fmt.Errorf("expected at least 2 pages, got %d (total zones: %d)", pageCount, totalCount)
		}
		if totalCount < 5 {
			return fmt.Errorf("expected at least 5 zones total across pages, got %d", totalCount)
		}
		return nil
	}))

	return results
}
