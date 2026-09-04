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
		resp, err := tc.createHealthCheck(tc.callerRef("hcref"), &types.HealthCheckConfig{
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
		tc.deleteHealthCheck(aws.ToString(resp.HealthCheck.Id))
		return nil
	}))

	var healthCheckID string
	results = append(results, r.RunTest("route53", "CreateHealthCheck_GetID", func() error {
		hcRef := tc.callerRef("hcref2")
		resp, err := tc.createHealthCheck(hcRef, tcpHealthCheck("hc.example.com", 8080))
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

	// CreateHealthCheck documents retry semantics for the CallerReference: a
	// resend with the same reference and the same settings returns the
	// existing health check instead of creating a duplicate.
	results = append(results, r.RunTest("route53", "CreateHealthCheck_IdempotentRetry", func() error {
		ref := tc.callerRef("hcref-idem")
		first, err := tc.createHealthCheck(ref, tcpHealthCheck("idem.example.com", 8081))
		if err != nil {
			return err
		}
		hcID := aws.ToString(first.HealthCheck.Id)
		defer tc.deleteHealthCheck(hcID)

		second, err := tc.createHealthCheck(ref, tcpHealthCheck("idem.example.com", 8081))
		if err != nil {
			return fmt.Errorf("idempotent retry: %v", err)
		}
		if aws.ToString(second.HealthCheck.Id) != hcID {
			return fmt.Errorf("idempotent retry returned different ID: %q vs %q", aws.ToString(second.HealthCheck.Id), hcID)
		}
		if aws.ToInt64(second.HealthCheck.HealthCheckVersion) != aws.ToInt64(first.HealthCheck.HealthCheckVersion) {
			return fmt.Errorf("idempotent retry changed version: %d vs %d",
				aws.ToInt64(second.HealthCheck.HealthCheckVersion), aws.ToInt64(first.HealthCheck.HealthCheckVersion))
		}
		return nil
	}))

	// The same CallerReference with different settings is rejected with
	// HealthCheckAlreadyExists.
	results = append(results, r.RunTest("route53", "CreateHealthCheck_SameCallerRefDifferentSettings", func() error {
		ref := tc.callerRef("hcref-diff")
		first, err := tc.createHealthCheck(ref, tcpHealthCheck("diff.example.com", 8082))
		if err != nil {
			return err
		}
		defer tc.deleteHealthCheck(aws.ToString(first.HealthCheck.Id))

		_, err = tc.createHealthCheck(ref, tcpHealthCheck("diff.example.com", 8083))
		if err := AssertErrorContains(err, "HealthCheckAlreadyExists"); err != nil {
			return err
		}
		return nil
	}))

	// A retry whose CallerReference matches a recently deleted health check
	// fails with HealthCheckAlreadyExists; the reference is retained for a
	// limited period after deletion.
	results = append(results, r.RunTest("route53", "CreateHealthCheck_RetryAfterDelete", func() error {
		ref := tc.callerRef("hcref-del")
		config := tcpHealthCheck("delref.example.com", 8084)
		created, err := tc.createHealthCheck(ref, config)
		if err != nil {
			return err
		}
		if _, err := tc.client.DeleteHealthCheck(tc.ctx, &route53.DeleteHealthCheckInput{
			HealthCheckId: created.HealthCheck.Id,
		}); err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		_, err = tc.createHealthCheck(ref, config)
		if err := AssertErrorContains(err, "HealthCheckAlreadyExists"); err != nil {
			return err
		}
		return nil
	}))

	if healthCheckID != "" {
		results = append(results, r.RunTest("route53", "GetHealthCheck", func() error {
			resp, err := tc.getHealthCheck(healthCheckID)
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
			resp, err := tc.getHealthCheck(healthCheckID)
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

		results = append(results, r.RunTest("route53", "UpdateHealthCheck_VersionMismatch", func() error {
			createResp, err := tc.createHealthCheck(tc.callerRef("hcref-ver"), tcpHealthCheck("ver.example.com", 8080))
			if err != nil {
				return fmt.Errorf("create: %v", err)
			}
			hcID := aws.ToString(createResp.HealthCheck.Id)

			defer tc.deleteHealthCheck(hcID)

			// A freshly created health check is at version 1; a stale
			// version must be rejected with a 409 conflict.
			_, err = tc.client.UpdateHealthCheck(tc.ctx, &route53.UpdateHealthCheckInput{
				HealthCheckId:      aws.String(hcID),
				HealthCheckVersion: aws.Int64(2),
				Port:               aws.Int32(9090),
			})
			return expectRoute53Error(err, "HealthCheckVersionMismatch", 409)
		}))

		results = append(results, r.RunTest("route53", "DeleteHealthCheck", func() error {
			_, err := tc.client.DeleteHealthCheck(tc.ctx, &route53.DeleteHealthCheckInput{
				HealthCheckId: aws.String(healthCheckID),
			})
			return err
		}))

		results = append(results, r.RunTest("route53", "GetHealthCheck_AfterDelete", func() error {
			_, err := tc.getHealthCheck(healthCheckID)
			if err := AssertErrorContains(err, "NoSuchHealthCheck"); err != nil {
				return err
			}
			return nil
		}))

		results = append(results, r.RunTest("route53", "HealthCheck_NonExistent", func() error {
			for _, c := range []struct {
				name string
				call func() error
			}{
				{
					name: "get a non-existent health check ID",
					call: func() error {
						_, err := tc.getHealthCheck("00000000-0000-0000-0000-000000000000")
						return err
					},
				},
				{
					name: "delete a non-existent health check ID",
					call: func() error {
						_, err := tc.client.DeleteHealthCheck(tc.ctx, &route53.DeleteHealthCheckInput{
							HealthCheckId: aws.String("00000000-0000-0000-0000-000000000000"),
						})
						return err
					},
				},
			} {
				if err := AssertErrorContains(c.call(), "NoSuchHealthCheck"); err != nil {
					return fmt.Errorf("%s: %w", c.name, err)
				}
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
		resp, err := tc.createHealthCheck(tc.callerRef("hcref-port"), &types.HealthCheckConfig{
			Type:                     types.HealthCheckTypeHttp,
			FullyQualifiedDomainName: aws.String("porttest.example.com"),
		})
		if err != nil {
			return err
		}
		hcID := aws.ToString(resp.HealthCheck.Id)

		defer tc.deleteHealthCheck(hcID)

		getResp, err := tc.getHealthCheck(hcID)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		// route53HCPort mirrors internal/common/serviceports.Route53HC.
		const route53HCPort = 50089
		port := aws.ToInt32(getResp.HealthCheck.HealthCheckConfig.Port)
		if port != route53HCPort {
			return fmt.Errorf("expected default port %d, got %d", route53HCPort, port)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ListHealthChecks_Pagination", func() error {
		var hcIDs []string
		for i := 0; i < 5; i++ {
			resp, err := tc.createHealthCheck(tc.callerRef(fmt.Sprintf("hcpagref-%d", i)), tcpHealthCheck(fmt.Sprintf("hcpag%d.example.com", i), 80))
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
			tc.deleteHealthCheck(id)
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
