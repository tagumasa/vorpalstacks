package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunRoute53Tests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "route53",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := route53.NewFromConfig(cfg)
	ctx := context.Background()

	// Test 1: List Hosted Zones
	results = append(results, r.RunTest("route53", "ListHostedZones", func() error {
		_, err := client.ListHostedZones(ctx, &route53.ListHostedZonesInput{
			MaxItems: aws.Int32(10),
		})
		return err
	}))

	// Test 2: Create Hosted Zone
	domainName := fmt.Sprintf("example-%d.com.", time.Now().UnixNano())
	var hostedZoneID string
	results = append(results, r.RunTest("route53", "CreateHostedZone", func() error {
		resp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(domainName),
			CallerReference: aws.String(fmt.Sprintf("ref-%d", time.Now().UnixNano())),
		})
		if err == nil {
			hostedZoneID = aws.ToString(resp.HostedZone.Id)
		}
		return err
	}))

	// Test 3: Get Hosted Zone
	if hostedZoneID != "" {
		results = append(results, r.RunTest("route53", "GetHostedZone", func() error {
			_, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{
				Id: aws.String(hostedZoneID),
			})
			return err
		}))

		// Test 4: List Resource Record Sets
		results = append(results, r.RunTest("route53", "ListResourceRecordSets", func() error {
			_, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
				HostedZoneId: aws.String(hostedZoneID),
				MaxItems:     aws.Int32(10),
			})
			return err
		}))

		var changeID string
		// Test 5: Change Resource Record Sets (Create A Record)
		results = append(results, r.RunTest("route53", "ChangeResourceRecordSets", func() error {
			resp, err := client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
				HostedZoneId: aws.String(hostedZoneID),
				ChangeBatch: &types.ChangeBatch{
					Changes: []types.Change{
						{
							Action: types.ChangeActionCreate,
							ResourceRecordSet: &types.ResourceRecordSet{
								Name: aws.String(fmt.Sprintf("test.%s", domainName)),
								Type: types.RRTypeA,
								TTL:  aws.Int64(300),
								ResourceRecords: []types.ResourceRecord{
									{
										Value: aws.String("192.0.2.1"),
									},
								},
							},
						},
					},
				},
			})
			if err == nil {
				changeID = aws.ToString(resp.ChangeInfo.Id)
			}
			return err
		}))

		// Test 6: Get Change (only if changeID was created)
		if changeID != "" {
			results = append(results, r.RunTest("route53", "GetChange", func() error {
				_, err := client.GetChange(ctx, &route53.GetChangeInput{
					Id: aws.String(changeID),
				})
				return err
			}))
		}

		// Test 7: Delete Resource Record (cleanup before deleting hosted zone)
		results = append(results, r.RunTest("route53", "DeleteResourceRecord", func() error {
			_, err := client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
				HostedZoneId: aws.String(hostedZoneID),
				ChangeBatch: &types.ChangeBatch{
					Changes: []types.Change{
						{
							Action: types.ChangeActionDelete,
							ResourceRecordSet: &types.ResourceRecordSet{
								Name:            aws.String(fmt.Sprintf("test.%s", domainName)),
								Type:            types.RRTypeA,
								TTL:             aws.Int64(300),
								ResourceRecords: []types.ResourceRecord{{Value: aws.String("192.0.2.1")}},
							},
						},
					},
				},
			})
			return err
		}))

		// Test 8: GetDNSSEC (before deleting hosted zone)
		results = append(results, r.RunTest("route53", "GetDNSSEC", func() error {
			_, err := client.GetDNSSEC(ctx, &route53.GetDNSSECInput{
				HostedZoneId: aws.String(hostedZoneID),
			})
			return err
		}))

		// Test 9: Delete Hosted Zone
		results = append(results, r.RunTest("route53", "DeleteHostedZone", func() error {
			_, err := client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{
				Id: aws.String(hostedZoneID),
			})
			return err
		}))
	}

	// Test 10: ListReusableDelegationSets
	results = append(results, r.RunTest("route53", "ListReusableDelegationSets", func() error {
		_, err := client.ListReusableDelegationSets(ctx, &route53.ListReusableDelegationSetsInput{
			MaxItems: aws.Int32(10),
		})
		return err
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("route53", "GetHostedZone_NonExistent", func() error {
		_, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{
			Id: aws.String("Z00000000000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent hosted zone")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "DeleteHostedZone_NonExistent", func() error {
		_, err := client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{
			Id: aws.String("Z00000000000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent hosted zone")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "GetChange_NonExistent", func() error {
		_, err := client.GetChange(ctx, &route53.GetChangeInput{
			Id: aws.String("C0000000000000000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent change")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "CreateHostedZone_ContentVerify", func() error {
		verifyDomain := fmt.Sprintf("verify-%d.com.", time.Now().UnixNano())
		verifyRef := fmt.Sprintf("ref-%d", time.Now().UnixNano())
		resp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(verifyDomain),
			CallerReference: aws.String(verifyRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		hzID := aws.ToString(resp.HostedZone.Id)

		defer client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(hzID)})

		if resp.HostedZone.Name == nil || *resp.HostedZone.Name != verifyDomain {
			return fmt.Errorf("domain name mismatch: got %q, want %q", aws.ToString(resp.HostedZone.Name), verifyDomain)
		}

		getResp, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{Id: aws.String(hzID)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.HostedZone.Name == nil || *getResp.HostedZone.Name != verifyDomain {
			return fmt.Errorf("get domain name mismatch")
		}
		return nil
	}))

	return results
}
