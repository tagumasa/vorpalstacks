package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func (r *TestRunner) runRoute53HealthCheckTests(tc *route53TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("route53", "CreateHealthCheck", func() error {
		resp, err := tc.client.CreateHealthCheck(tc.ctx, &route53.CreateHealthCheckInput{
			CallerReference: aws.String(fmt.Sprintf("hcref-%d", tc.uniq)),
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
		if resp.HealthCheck == nil || resp.HealthCheck.Id == nil {
			return fmt.Errorf("health check or ID is nil")
		}
		if resp.HealthCheck.HealthCheckConfig == nil {
			return fmt.Errorf("health check config is nil")
		}
		if resp.HealthCheck.HealthCheckConfig.Type != types.HealthCheckTypeHttp {
			return fmt.Errorf("type mismatch: got %v", resp.HealthCheck.HealthCheckConfig.Type)
		}
		tc.client.DeleteHealthCheck(tc.ctx, &route53.DeleteHealthCheckInput{
			HealthCheckId: resp.HealthCheck.Id,
		})
		return nil
	}))

	var healthCheckID string
	results = append(results, r.RunTest("route53", "CreateHealthCheck_GetID", func() error {
		hcRef := fmt.Sprintf("hcref2-%d", tc.uniq)
		resp, err := tc.client.CreateHealthCheck(tc.ctx, &route53.CreateHealthCheckInput{
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
		if resp.HealthCheck.CallerReference == nil || *resp.HealthCheck.CallerReference != hcRef {
			return fmt.Errorf("caller reference mismatch: got %q, want %q", aws.ToString(resp.HealthCheck.CallerReference), hcRef)
		}
		return nil
	}))

	if healthCheckID != "" {
		results = append(results, r.RunTest("route53", "GetHealthCheck", func() error {
			resp, err := tc.client.GetHealthCheck(tc.ctx, &route53.GetHealthCheckInput{
				HealthCheckId: aws.String(healthCheckID),
			})
			if err != nil {
				return err
			}
			if resp.HealthCheck == nil {
				return fmt.Errorf("health check is nil")
			}
			if aws.ToString(resp.HealthCheck.Id) != healthCheckID {
				return fmt.Errorf("ID mismatch: got %q, want %q", aws.ToString(resp.HealthCheck.Id), healthCheckID)
			}
			if resp.HealthCheck.HealthCheckConfig == nil {
				return fmt.Errorf("health check config is nil")
			}
			if resp.HealthCheck.HealthCheckConfig.Type != types.HealthCheckTypeTcp {
				return fmt.Errorf("health check type mismatch: got %v", resp.HealthCheck.HealthCheckConfig.Type)
			}
			if aws.ToInt32(resp.HealthCheck.HealthCheckConfig.Port) != 8080 {
				return fmt.Errorf("port mismatch: got %d, want 8080", aws.ToInt32(resp.HealthCheck.HealthCheckConfig.Port))
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "UpdateHealthCheck", func() error {
			resp, err := tc.client.UpdateHealthCheck(tc.ctx, &route53.UpdateHealthCheckInput{
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
			resp, err := tc.client.GetHealthCheck(tc.ctx, &route53.GetHealthCheckInput{
				HealthCheckId: aws.String(healthCheckID),
			})
			if err != nil {
				return err
			}
			cfg := resp.HealthCheck.HealthCheckConfig
			if aws.ToInt32(cfg.FailureThreshold) != 5 {
				return fmt.Errorf("failure threshold mismatch: got %d", aws.ToInt32(cfg.FailureThreshold))
			}
			if aws.ToString(cfg.ResourcePath) != "/updated" {
				return fmt.Errorf("resource path mismatch: got %q", aws.ToString(cfg.ResourcePath))
			}
			if !aws.ToBool(cfg.Disabled) {
				return fmt.Errorf("expected disabled=true")
			}
			if !aws.ToBool(cfg.Inverted) {
				return fmt.Errorf("expected inverted=true")
			}
			if aws.ToString(cfg.FullyQualifiedDomainName) != "updated.example.com" {
				return fmt.Errorf("domain name mismatch: got %q", aws.ToString(cfg.FullyQualifiedDomainName))
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "DeleteHealthCheck", func() error {
			_, err := tc.client.DeleteHealthCheck(tc.ctx, &route53.DeleteHealthCheckInput{
				HealthCheckId: aws.String(healthCheckID),
			})
			return err
		}))

		results = append(results, r.RunTest("route53", "GetHealthCheck_AfterDelete", func() error {
			_, err := tc.client.GetHealthCheck(tc.ctx, &route53.GetHealthCheckInput{
				HealthCheckId: aws.String(healthCheckID),
			})
			if err := AssertErrorContains(err, "NoSuchHealthCheck"); err != nil {
				return err
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "GetHealthCheck_NonExistent", func() error {
			_, err := tc.client.GetHealthCheck(tc.ctx, &route53.GetHealthCheckInput{
				HealthCheckId: aws.String("00000000-0000-0000-0000-000000000000"),
			})
			if err := AssertErrorContains(err, "NoSuchHealthCheck"); err != nil {
				return err
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "DeleteHealthCheck_NonExistent", func() error {
			_, err := tc.client.DeleteHealthCheck(tc.ctx, &route53.DeleteHealthCheckInput{
				HealthCheckId: aws.String("00000000-0000-0000-0000-000000000000"),
			})
			if err := AssertErrorContains(err, "NoSuchHealthCheck"); err != nil {
				return err
			}
			return nil
		}))
	}

	results = append(results, r.RunTest("route53", "ListHealthChecks", func() error {
		resp, err := tc.client.ListHealthChecks(tc.ctx, &route53.ListHealthChecksInput{
			MaxItems: aws.Int32(100),
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
		resp, err := tc.client.CreateHealthCheck(tc.ctx, &route53.CreateHealthCheckInput{
			CallerReference: aws.String(fmt.Sprintf("hcref-port-%d", tc.uniq)),
			HealthCheckConfig: &types.HealthCheckConfig{
				Type:                     types.HealthCheckTypeHttp,
				FullyQualifiedDomainName: aws.String("porttest.example.com"),
			},
		})
		if err != nil {
			return err
		}
		hcID := aws.ToString(resp.HealthCheck.Id)

		defer tc.client.DeleteHealthCheck(tc.ctx, &route53.DeleteHealthCheckInput{HealthCheckId: aws.String(hcID)})

		getResp, err := tc.client.GetHealthCheck(tc.ctx, &route53.GetHealthCheckInput{
			HealthCheckId: aws.String(hcID),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		port := aws.ToInt32(getResp.HealthCheck.HealthCheckConfig.Port)
		if port != 8089 {
			return fmt.Errorf("expected default port 8089, got %d", port)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListHealthChecks_Pagination", func() error {
		var hcIDs []string
		for i := 0; i < 5; i++ {
			resp, err := tc.client.CreateHealthCheck(tc.ctx, &route53.CreateHealthCheckInput{
				CallerReference: aws.String(fmt.Sprintf("hcpagref-%d-%d", tc.uniq, i)),
				HealthCheckConfig: &types.HealthCheckConfig{
					Type:                     types.HealthCheckTypeTcp,
					FullyQualifiedDomainName: aws.String(fmt.Sprintf("hcpag%d.example.com", i)),
					Port:                     aws.Int32(80),
				},
			})
			if err != nil {
				return fmt.Errorf("create health check %d: %v", i, err)
			}
			hcIDs = append(hcIDs, aws.ToString(resp.HealthCheck.Id))
		}

		var marker *string
		totalCount := 0
		pageCount := 0
		for {
			resp, err := tc.client.ListHealthChecks(tc.ctx, &route53.ListHealthChecksInput{
				Marker:   marker,
				MaxItems: aws.Int32(2),
			})
			if err != nil {
				return fmt.Errorf("list page: %v", err)
			}
			pageCount++
			totalCount += len(resp.HealthChecks)
			if !resp.IsTruncated || resp.NextMarker == nil {
				break
			}
			marker = resp.NextMarker
		}

		for _, id := range hcIDs {
			tc.client.DeleteHealthCheck(tc.ctx, &route53.DeleteHealthCheckInput{HealthCheckId: aws.String(id)})
		}

		if pageCount < 2 {
			return fmt.Errorf("expected at least 2 pages, got %d (total: %d)", pageCount, totalCount)
		}
		if totalCount < 5 {
			return fmt.Errorf("expected at least 5 health checks, got %d", totalCount)
		}
		return nil
	}))

	return results
}
