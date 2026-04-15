package testutil

import (
	"context"
	"fmt"
	"strings"
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
	uniq := time.Now().UnixNano()

	domainName := fmt.Sprintf("example-%d.com.", uniq)

	var hostedZoneID string
	results = append(results, r.RunTest("route53", "CreateHostedZone", func() error {
		resp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(domainName),
			CallerReference: aws.String(fmt.Sprintf("ref-%d", uniq)),
		})
		if err != nil {
			return err
		}
		if resp.HostedZone == nil {
			return fmt.Errorf("hosted zone is nil")
		}
		hostedZoneID = aws.ToString(resp.HostedZone.Id)
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListHostedZones", func() error {
		resp, err := client.ListHostedZones(ctx, &route53.ListHostedZonesInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.HostedZones == nil {
			return fmt.Errorf("hosted zones list is nil")
		}
		return nil
	}))

	if hostedZoneID != "" {
		results = append(results, r.RunTest("route53", "GetHostedZone", func() error {
			resp, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{
				Id: aws.String(hostedZoneID),
			})
			if err != nil {
				return err
			}
			if resp.HostedZone == nil {
				return fmt.Errorf("hosted zone is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "ListResourceRecordSets", func() error {
			resp, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
				HostedZoneId: aws.String(hostedZoneID),
				MaxItems:     aws.Int32(10),
			})
			if err != nil {
				return err
			}
			if resp.ResourceRecordSets == nil {
				return fmt.Errorf("resource record sets list is nil")
			}
			return nil
		}))

		var changeID string
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
			if err != nil {
				return err
			}
			if resp == nil || resp.ChangeInfo == nil {
				return fmt.Errorf("response or change info is nil")
			}
			changeID = aws.ToString(resp.ChangeInfo.Id)
			return nil
		}))

		if changeID != "" {
			results = append(results, r.RunTest("route53", "GetChange", func() error {
				resp, err := client.GetChange(ctx, &route53.GetChangeInput{
					Id: aws.String(changeID),
				})
				if err != nil {
					return err
				}
				if resp.ChangeInfo == nil {
					return fmt.Errorf("change info is nil")
				}
				return nil
			}))
		}

		results = append(results, r.RunTest("route53", "DeleteResourceRecord", func() error {
			resp, err := client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
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
			if err != nil {
				return err
			}
			if resp == nil {
				return fmt.Errorf("response is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "GetDNSSEC", func() error {
			resp, err := client.GetDNSSEC(ctx, &route53.GetDNSSECInput{
				HostedZoneId: aws.String(hostedZoneID),
			})
			if err != nil {
				return err
			}
			if resp == nil {
				return fmt.Errorf("response is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "DeleteHostedZone", func() error {
			resp, err := client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{
				Id: aws.String(hostedZoneID),
			})
			if err != nil {
				return err
			}
			if resp == nil {
				return fmt.Errorf("response is nil")
			}
			return nil
		}))
	}

	results = append(results, r.RunTest("route53", "ListReusableDelegationSets", func() error {
		resp, err := client.ListReusableDelegationSets(ctx, &route53.ListReusableDelegationSetsInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.DelegationSets != nil && len(resp.DelegationSets) != 0 {
			return fmt.Errorf("expected no delegation sets, got %d", len(resp.DelegationSets))
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "GetHostedZone_NonExistent", func() error {
		_, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{
			Id: aws.String("Z00000000000000000000"),
		})
		if err := AssertErrorContains(err, "NoSuchHostedZone"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "DeleteHostedZone_NonExistent", func() error {
		_, err := client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{
			Id: aws.String("Z00000000000000000000"),
		})
		if err := AssertErrorContains(err, "NoSuchHostedZone"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "GetChange_NonExistent", func() error {
		_, err := client.GetChange(ctx, &route53.GetChangeInput{
			Id: aws.String("C0000000000000000000000000"),
		})
		if err := AssertErrorContains(err, "NoSuchChange"); err != nil {
			return err
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

	results = append(results, r.RunTest("route53", "ListHostedZonesByName", func() error {
		resp, err := client.ListHostedZonesByName(ctx, &route53.ListHostedZonesByNameInput{
			MaxItems: aws.Int32(10),
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
		testDomain := fmt.Sprintf("sorttest-%d.com.", time.Now().UnixNano())
		hzResp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(testDomain),
			CallerReference: aws.String(fmt.Sprintf("sortref-%d", time.Now().UnixNano())),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}

		resp, err := client.ListHostedZonesByName(ctx, &route53.ListHostedZonesByNameInput{
			DNSName:  aws.String(testDomain),
			MaxItems: aws.Int32(10),
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
		client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{
			Id: hzResp.HostedZone.Id,
		})
		return nil
	}))

	results = append(results, r.RunTest("route53", "UpdateHostedZoneComment", func() error {
		ucDomain := fmt.Sprintf("updatecomment-%d.com.", time.Now().UnixNano())
		ucRef := fmt.Sprintf("ucref-%d", time.Now().UnixNano())
		createResp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(ucDomain),
			CallerReference: aws.String(ucRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		ucID := aws.ToString(createResp.HostedZone.Id)

		defer client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(ucID)})

		comment := "test comment for zone"
		_, err = client.UpdateHostedZoneComment(ctx, &route53.UpdateHostedZoneCommentInput{
			Id:      aws.String(ucID),
			Comment: aws.String(comment),
		})
		if err != nil {
			return fmt.Errorf("update comment: %v", err)
		}

		getResp, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{Id: aws.String(ucID)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.HostedZone.Config == nil || aws.ToString(getResp.HostedZone.Config.Comment) != comment {
			return fmt.Errorf("comment mismatch: got %q, want %q", aws.ToString(getResp.HostedZone.Config.Comment), comment)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "CreateHealthCheck", func() error {
		resp, err := client.CreateHealthCheck(ctx, &route53.CreateHealthCheckInput{
			CallerReference: aws.String(fmt.Sprintf("hcref-%d", time.Now().UnixNano())),
			HealthCheckConfig: &types.HealthCheckConfig{
				Type:                         types.HealthCheckTypeHttp,
				ResourcePath:                 aws.String("/health"),
				FullyQualifiedDomainName:     aws.String("example.com"),
				RequestInterval:              aws.Int32(30),
				FailureThreshold:             aws.Int32(3),
				MeasureLatency:               aws.Bool(true),
				Disabled:                     aws.Bool(false),
				EnableSNI:                    aws.Bool(true),
				IPAddress:                    aws.String("192.0.2.1"),
				Port:                         aws.Int32(443),
				Inverted:                     aws.Bool(false),
				InsufficientDataHealthStatus: types.InsufficientDataHealthStatusLastKnownStatus,
			},
		})
		if err != nil {
			return err
		}
		client.DeleteHealthCheck(ctx, &route53.DeleteHealthCheckInput{
			HealthCheckId: resp.HealthCheck.Id,
		})
		return nil
	}))

	var healthCheckID string
	results = append(results, r.RunTest("route53", "CreateHealthCheck_GetID", func() error {
		hcRef := fmt.Sprintf("hcref2-%d", time.Now().UnixNano())
		resp, err := client.CreateHealthCheck(ctx, &route53.CreateHealthCheckInput{
			CallerReference: aws.String(hcRef),
			HealthCheckConfig: &types.HealthCheckConfig{
				Type:                     types.HealthCheckTypeTcp,
				FullyQualifiedDomainName: aws.String("hc.example.com"),
				Port:                     aws.Int32(8080),
			},
		})
		if err != nil {
			return err
		}
		if resp.HealthCheck == nil || resp.HealthCheck.Id == nil {
			return fmt.Errorf("health check ID is nil")
		}
		healthCheckID = aws.ToString(resp.HealthCheck.Id)
		return nil
	}))

	if healthCheckID != "" {
		results = append(results, r.RunTest("route53", "GetHealthCheck", func() error {
			resp, err := client.GetHealthCheck(ctx, &route53.GetHealthCheckInput{
				HealthCheckId: aws.String(healthCheckID),
			})
			if err != nil {
				return err
			}
			if resp.HealthCheck == nil {
				return fmt.Errorf("health check is nil")
			}
			if resp.HealthCheck.HealthCheckConfig == nil {
				return fmt.Errorf("health check config is nil")
			}
			if resp.HealthCheck.HealthCheckConfig.Type != types.HealthCheckTypeTcp {
				return fmt.Errorf("health check type mismatch: got %v", resp.HealthCheck.HealthCheckConfig.Type)
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "UpdateHealthCheck", func() error {
			resp, err := client.UpdateHealthCheck(ctx, &route53.UpdateHealthCheckInput{
				HealthCheckId:            aws.String(healthCheckID),
				ResourcePath:             aws.String("/updated"),
				FailureThreshold:         aws.Int32(5),
				Disabled:                 aws.Bool(true),
				Inverted:                 aws.Bool(true),
				EnableSNI:                aws.Bool(false),
				FullyQualifiedDomainName: aws.String("updated.example.com"),
			})
			if err != nil {
				return err
			}
			if resp.HealthCheck == nil {
				return fmt.Errorf("health check is nil after update")
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "UpdateHealthCheck_VerifyContent", func() error {
			resp, err := client.GetHealthCheck(ctx, &route53.GetHealthCheckInput{
				HealthCheckId: aws.String(healthCheckID),
			})
			if err != nil {
				return err
			}
			if aws.ToInt32(resp.HealthCheck.HealthCheckConfig.FailureThreshold) != 5 {
				return fmt.Errorf("failure threshold mismatch: got %d", aws.ToInt32(resp.HealthCheck.HealthCheckConfig.FailureThreshold))
			}
			if aws.ToString(resp.HealthCheck.HealthCheckConfig.ResourcePath) != "/updated" {
				return fmt.Errorf("resource path mismatch: got %q", aws.ToString(resp.HealthCheck.HealthCheckConfig.ResourcePath))
			}
			if !aws.ToBool(resp.HealthCheck.HealthCheckConfig.Disabled) {
				return fmt.Errorf("expected disabled=true")
			}
			if !aws.ToBool(resp.HealthCheck.HealthCheckConfig.Inverted) {
				return fmt.Errorf("expected inverted=true")
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "DeleteHealthCheck", func() error {
			_, err := client.DeleteHealthCheck(ctx, &route53.DeleteHealthCheckInput{
				HealthCheckId: aws.String(healthCheckID),
			})
			return err
		}))

		results = append(results, r.RunTest("route53", "GetHealthCheck_NonExistent", func() error {
			_, err := client.GetHealthCheck(ctx, &route53.GetHealthCheckInput{
				HealthCheckId: aws.String("00000000-0000-0000-0000-000000000000"),
			})
			if err := AssertErrorContains(err, "NoSuchHealthCheck"); err != nil {
				return err
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "DeleteHealthCheck_NonExistent", func() error {
			_, err := client.DeleteHealthCheck(ctx, &route53.DeleteHealthCheckInput{
				HealthCheckId: aws.String("00000000-0000-0000-0000-000000000000"),
			})
			if err := AssertErrorContains(err, "NoSuchHealthCheck"); err != nil {
				return err
			}
			return nil
		}))
	}

	results = append(results, r.RunTest("route53", "ListHealthChecks", func() error {
		resp, err := client.ListHealthChecks(ctx, &route53.ListHealthChecksInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.HealthChecks == nil {
			return fmt.Errorf("health checks list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "HealthCheckConfig_DefaultPort", func() error {
		resp, err := client.CreateHealthCheck(ctx, &route53.CreateHealthCheckInput{
			CallerReference: aws.String(fmt.Sprintf("hcref-port-%d", time.Now().UnixNano())),
			HealthCheckConfig: &types.HealthCheckConfig{
				Type:                     types.HealthCheckTypeHttp,
				FullyQualifiedDomainName: aws.String("porttest.example.com"),
			},
		})
		if err != nil {
			return err
		}
		hcID := aws.ToString(resp.HealthCheck.Id)

		defer client.DeleteHealthCheck(ctx, &route53.DeleteHealthCheckInput{HealthCheckId: aws.String(hcID)})

		getResp, err := client.GetHealthCheck(ctx, &route53.GetHealthCheckInput{
			HealthCheckId: aws.String(hcID),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		port := aws.ToInt32(getResp.HealthCheck.HealthCheckConfig.Port)
		if port != 80 {
			return fmt.Errorf("expected default port 80, got %d", port)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "AssociateVPCWithHostedZone", func() error {
		privateDomain := fmt.Sprintf("private-%d.com.", time.Now().UnixNano())
		privateRef := fmt.Sprintf("privref-%d", time.Now().UnixNano())
		createResp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
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
			client.DisassociateVPCFromHostedZone(ctx, &route53.DisassociateVPCFromHostedZoneInput{
				HostedZoneId: aws.String(privateZoneID),
				VPC: &types.VPC{
					VPCId:     aws.String("vpc-xyz12345"),
					VPCRegion: types.VPCRegionUsEast1,
				},
			})
			client.DisassociateVPCFromHostedZone(ctx, &route53.DisassociateVPCFromHostedZoneInput{
				HostedZoneId: aws.String(privateZoneID),
				VPC: &types.VPC{
					VPCId:     aws.String("vpc-abcdef01"),
					VPCRegion: types.VPCRegionUsEast1,
				},
			})
			client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(privateZoneID)})
		}()

		_, err = client.AssociateVPCWithHostedZone(ctx, &route53.AssociateVPCWithHostedZoneInput{
			HostedZoneId: aws.String(privateZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String("vpc-xyz12345"),
				VPCRegion: types.VPCRegionUsEast1,
			},
		})
		if err != nil {
			return fmt.Errorf("associate: %v", err)
		}

		getResp, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{
			Id: aws.String(privateZoneID),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(getResp.VPCs) < 2 {
			return fmt.Errorf("expected at least 2 VPCs, got %d", len(getResp.VPCs))
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "DisassociateVPCFromHostedZone", func() error {
		dsDomain := fmt.Sprintf("disassoc-%d.com.", time.Now().UnixNano())
		dsRef := fmt.Sprintf("dsref-%d", time.Now().UnixNano())
		createResp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
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

		defer client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(dsZoneID)})

		_, err = client.AssociateVPCWithHostedZone(ctx, &route53.AssociateVPCWithHostedZoneInput{
			HostedZoneId: aws.String(dsZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String("vpc-disassoc2"),
				VPCRegion: types.VPCRegionUsEast1,
			},
		})
		if err != nil {
			return fmt.Errorf("associate: %v", err)
		}

		_, err = client.DisassociateVPCFromHostedZone(ctx, &route53.DisassociateVPCFromHostedZoneInput{
			HostedZoneId: aws.String(dsZoneID),
			VPC: &types.VPC{
				VPCId:     aws.String("vpc-disassoc2"),
				VPCRegion: types.VPCRegionUsEast1,
			},
		})
		if err != nil {
			return fmt.Errorf("disassociate: %v", err)
		}

		getResp, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{
			Id: aws.String(dsZoneID),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(getResp.VPCs) != 1 {
			return fmt.Errorf("expected 1 VPC after disassociation, got %d", len(getResp.VPCs))
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ChangeTagsForResource", func() error {
		tagDomain := fmt.Sprintf("tags-%d.com.", time.Now().UnixNano())
		tagRef := fmt.Sprintf("tagref-%d", time.Now().UnixNano())
		createResp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(tagDomain),
			CallerReference: aws.String(tagRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		tagZoneID := aws.ToString(createResp.HostedZone.Id)

		defer client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(tagZoneID)})

		_, err = client.ChangeTagsForResource(ctx, &route53.ChangeTagsForResourceInput{
			ResourceType: types.TagResourceTypeHostedzone,
			ResourceId:   aws.String(tagZoneID),
			AddTags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Owner"), Value: aws.String("team-a")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags: %v", err)
		}

		listResp, err := client.ListTagsForResource(ctx, &route53.ListTagsForResourceInput{
			ResourceType: types.TagResourceTypeHostedzone,
			ResourceId:   aws.String(tagZoneID),
		})
		if err != nil {
			return fmt.Errorf("list tags after add: %v", err)
		}
		if listResp.ResourceTagSet == nil || len(listResp.ResourceTagSet.Tags) != 2 {
			return fmt.Errorf("expected 2 tags after add, got %d", len(listResp.ResourceTagSet.Tags))
		}

		_, err = client.ChangeTagsForResource(ctx, &route53.ChangeTagsForResourceInput{
			ResourceType:  types.TagResourceTypeHostedzone,
			ResourceId:    aws.String(tagZoneID),
			RemoveTagKeys: []string{"Owner"},
		})
		if err != nil {
			return fmt.Errorf("remove tags: %v", err)
		}

		listResp2, err := client.ListTagsForResource(ctx, &route53.ListTagsForResourceInput{
			ResourceType: types.TagResourceTypeHostedzone,
			ResourceId:   aws.String(tagZoneID),
		})
		if err != nil {
			return fmt.Errorf("list tags after remove: %v", err)
		}
		if len(listResp2.ResourceTagSet.Tags) != 1 {
			return fmt.Errorf("expected 1 tag after removal, got %d", len(listResp2.ResourceTagSet.Tags))
		}
		if aws.ToString(listResp2.ResourceTagSet.Tags[0].Key) != "Environment" {
			return fmt.Errorf("expected Environment tag, got %q", aws.ToString(listResp2.ResourceTagSet.Tags[0].Key))
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListTagsForResource_HealthCheck", func() error {
		hcRef := fmt.Sprintf("hctagref-%d", time.Now().UnixNano())
		hcResp, err := client.CreateHealthCheck(ctx, &route53.CreateHealthCheckInput{
			CallerReference: aws.String(hcRef),
			HealthCheckConfig: &types.HealthCheckConfig{
				Type:                     types.HealthCheckTypeTcp,
				FullyQualifiedDomainName: aws.String("hctag.example.com"),
			},
		})
		if err != nil {
			return fmt.Errorf("create hc: %v", err)
		}
		hcID := aws.ToString(hcResp.HealthCheck.Id)

		defer client.DeleteHealthCheck(ctx, &route53.DeleteHealthCheckInput{HealthCheckId: aws.String(hcID)})

		_, err = client.ChangeTagsForResource(ctx, &route53.ChangeTagsForResourceInput{
			ResourceType: types.TagResourceTypeHealthcheck,
			ResourceId:   aws.String(hcID),
			AddTags: []types.Tag{
				{Key: aws.String("Monitor"), Value: aws.String("enabled")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags: %v", err)
		}

		listResp, err := client.ListTagsForResource(ctx, &route53.ListTagsForResourceInput{
			ResourceType: types.TagResourceTypeHealthcheck,
			ResourceId:   aws.String(hcID),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(listResp.ResourceTagSet.Tags) != 1 {
			return fmt.Errorf("expected 1 tag, got %d", len(listResp.ResourceTagSet.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ChangeResourceRecordSets_Upsert", func() error {
		upsertDomain := fmt.Sprintf("upsert-%d.com.", time.Now().UnixNano())
		upsertRef := fmt.Sprintf("upsertref-%d", time.Now().UnixNano())
		createResp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(upsertDomain),
			CallerReference: aws.String(upsertRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		upsertZoneID := aws.ToString(createResp.HostedZone.Id)

		defer client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(upsertZoneID)})

		recordName := fmt.Sprintf("upsert.%s", upsertDomain)
		_, err = client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: aws.String(upsertZoneID),
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: types.ChangeActionUpsert,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: aws.String(recordName),
							Type: types.RRTypeA,
							TTL:  aws.Int64(300),
							ResourceRecords: []types.ResourceRecord{
								{Value: aws.String("10.0.0.1")},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("upsert create: %v", err)
		}

		_, err = client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: aws.String(upsertZoneID),
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: types.ChangeActionUpsert,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: aws.String(recordName),
							Type: types.RRTypeA,
							TTL:  aws.Int64(600),
							ResourceRecords: []types.ResourceRecord{
								{Value: aws.String("10.0.0.2")},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("upsert update: %v", err)
		}

		listResp, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(upsertZoneID),
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		found := false
		for _, rs := range listResp.ResourceRecordSets {
			if aws.ToString(rs.Name) == recordName && rs.Type == types.RRTypeA {
				if aws.ToInt64(rs.TTL) != 600 {
					return fmt.Errorf("TTL mismatch after upsert: got %d, want 600", aws.ToInt64(rs.TTL))
				}
				if len(rs.ResourceRecords) == 1 && aws.ToString(rs.ResourceRecords[0].Value) != "10.0.0.2" {
					return fmt.Errorf("value mismatch after upsert: got %q", aws.ToString(rs.ResourceRecords[0].Value))
				}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("upserted record not found")
		}

		client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: aws.String(upsertZoneID),
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: types.ChangeActionDelete,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name:            aws.String(recordName),
							Type:            types.RRTypeA,
							TTL:             aws.Int64(600),
							ResourceRecords: []types.ResourceRecord{{Value: aws.String("10.0.0.2")}},
						},
					},
				},
			},
		})
		return nil
	}))

	results = append(results, r.RunTest("route53", "CreateHostedZone_PrivateWithComment", func() error {
		pvtDomain := fmt.Sprintf("private-comment-%d.com.", time.Now().UnixNano())
		pvtRef := fmt.Sprintf("pvtref-%d", time.Now().UnixNano())
		resp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
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

		defer client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(pvtID)})

		getResp, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{Id: aws.String(pvtID)})
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
		return nil
	}))

	results = append(results, r.RunTest("route53", "DelegationSet_Persisted", func() error {
		dsDomain := fmt.Sprintf("ds-persist-%d.com.", time.Now().UnixNano())
		dsRef := fmt.Sprintf("dspersist-%d", time.Now().UnixNano())
		createResp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(dsDomain),
			CallerReference: aws.String(dsRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		dsZoneID := aws.ToString(createResp.HostedZone.Id)

		defer client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(dsZoneID)})

		createDSID := aws.ToString(createResp.DelegationSet.Id)
		if createDSID == "" {
			return fmt.Errorf("delegation set ID is empty in create response")
		}

		getResp, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{Id: aws.String(dsZoneID)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		getDSID := aws.ToString(getResp.DelegationSet.Id)
		if getDSID == "" {
			return fmt.Errorf("delegation set ID is empty in get response")
		}
		if createDSID != getDSID {
			return fmt.Errorf("delegation set ID mismatch: create=%q, get=%q", createDSID, getDSID)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "CreateHostedZone_InvalidName", func() error {
		_, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String("invalid name with spaces"),
			CallerReference: aws.String(fmt.Sprintf("badref-%d", time.Now().UnixNano())),
		})
		if err := AssertErrorContains(err, "InvalidInput"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "AssociateVPCWithHostedZone_PublicZone", func() error {
		pubDomain := fmt.Sprintf("pub-vpc-test-%d.com.", time.Now().UnixNano())
		pubRef := fmt.Sprintf("pubvpc-%d", time.Now().UnixNano())
		createResp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(pubDomain),
			CallerReference: aws.String(pubRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		pubZoneID := aws.ToString(createResp.HostedZone.Id)

		defer client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(pubZoneID)})

		_, err = client.AssociateVPCWithHostedZone(ctx, &route53.AssociateVPCWithHostedZoneInput{
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

	results = append(results, r.RunTest("route53", "DeleteHostedZone_NotEmpty", func() error {
		neDomain := fmt.Sprintf("notempty-%d.com.", time.Now().UnixNano())
		neRef := fmt.Sprintf("neref-%d", time.Now().UnixNano())
		createResp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(neDomain),
			CallerReference: aws.String(neRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		neZoneID := aws.ToString(createResp.HostedZone.Id)

		defer client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(neZoneID)})

		client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: aws.String(neZoneID),
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: types.ChangeActionCreate,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: aws.String(fmt.Sprintf("keep.%s", neDomain)),
							Type: types.RRTypeA,
							TTL:  aws.Int64(300),
							ResourceRecords: []types.ResourceRecord{
								{Value: aws.String("10.0.0.1")},
							},
						},
					},
				},
			},
		})

		_, err = client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{
			Id: aws.String(neZoneID),
		})
		if err := AssertErrorContains(err, "HostedZoneNotEmpty"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ChangeResourceRecordSets_MultipleTypes", func() error {
		mtDomain := fmt.Sprintf("multitype-%d.com.", time.Now().UnixNano())
		mtRef := fmt.Sprintf("mtref-%d", time.Now().UnixNano())
		createResp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
			Name:            aws.String(mtDomain),
			CallerReference: aws.String(mtRef),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		mtZoneID := aws.ToString(createResp.HostedZone.Id)

		defer client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: aws.String(mtZoneID),
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{Action: types.ChangeActionDelete, ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(fmt.Sprintf("www.%s", mtDomain)), Type: types.RRTypeCname, TTL: aws.Int64(300),
						ResourceRecords: []types.ResourceRecord{{Value: aws.String("target.example.com")}},
					}},
					{Action: types.ChangeActionDelete, ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(fmt.Sprintf("txt.%s", mtDomain)), Type: types.RRTypeTxt, TTL: aws.Int64(300),
						ResourceRecords: []types.ResourceRecord{{Value: aws.String("v=spf1 include:example.com ~all")}},
					}},
				},
			},
		})
		defer client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(mtZoneID)})

		_, err = client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: aws.String(mtZoneID),
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: types.ChangeActionCreate,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: aws.String(fmt.Sprintf("www.%s", mtDomain)),
							Type: types.RRTypeCname,
							TTL:  aws.Int64(300),
							ResourceRecords: []types.ResourceRecord{
								{Value: aws.String("target.example.com")},
							},
						},
					},
					{
						Action: types.ChangeActionCreate,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: aws.String(fmt.Sprintf("txt.%s", mtDomain)),
							Type: types.RRTypeTxt,
							TTL:  aws.Int64(300),
							ResourceRecords: []types.ResourceRecord{
								{Value: aws.String("v=spf1 include:example.com ~all")},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create records: %v", err)
		}

		listResp, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(mtZoneID),
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}

		foundCNAME := false
		foundTXT := false
		for _, rs := range listResp.ResourceRecordSets {
			if rs.Type == types.RRTypeCname && strings.HasSuffix(aws.ToString(rs.Name), mtDomain) {
				foundCNAME = true
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

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("route53", "ListHostedZones_Pagination", func() error {
		pgTs := fmt.Sprintf("pzpg%d", time.Now().UnixNano())
		var pgZoneIDs []string
		for i := 0; i < 5; i++ {
			resp, err := client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
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
			resp, err := client.ListHostedZones(ctx, &route53.ListHostedZonesInput{
				Marker:   marker,
				MaxItems: aws.Int32(2),
			})
			if err != nil {
				for _, zid := range pgZoneIDs {
					client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(zid)})
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
			client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{Id: aws.String(zid)})
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
