package testutil

import (
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runSNSPlatformTests(tc *snsTestContext) []TestResult {
	var results []TestResult

	var platformAppArn string
	results = append(results, r.RunTest("sns", "CreatePlatformApplication", func() error {
		appName := tc.uniqueName("TestApp")
		resp, err := tc.client.CreatePlatformApplication(tc.ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(appName),
			Platform: aws.String("GCM"),
			Attributes: map[string]string{
				"PlatformCredential": "fake-credential",
			},
		})
		if err != nil {
			return err
		}
		if resp.PlatformApplicationArn == nil {
			return fmt.Errorf("PlatformApplicationArn is nil")
		}
		platformAppArn = *resp.PlatformApplicationArn

		getResp, err := tc.getPlatformApplicationAttributes(platformAppArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes == nil {
			return fmt.Errorf("Attributes is nil")
		}
		if getResp.Attributes["PlatformCredential"] != "fake-credential" {
			return fmt.Errorf("PlatformCredential mismatch: got %q", getResp.Attributes["PlatformCredential"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "CreatePlatformApplication_Duplicate", func() error {
		if platformAppArn == "" {
			return fmt.Errorf("platformAppArn not set from previous test")
		}
		appName := tc.uniqueName("TestApp")
		dupResp, err := tc.client.CreatePlatformApplication(tc.ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(appName),
			Platform: aws.String("GCM"),
			Attributes: map[string]string{
				"PlatformCredential": "fake-credential",
			},
		})
		if err != nil {
			return fmt.Errorf("create with unique name should succeed: %v", err)
		}
		tc.client.DeletePlatformApplication(tc.ctx, &sns.DeletePlatformApplicationInput{
			PlatformApplicationArn: dupResp.PlatformApplicationArn,
		})
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListPlatformApplications", func() error {
		arns, err := tc.listAllPlatformApplications()
		if err != nil {
			return err
		}
		if len(arns) == 0 {
			return fmt.Errorf("PlatformApplications is empty")
		}
		if !slices.Contains(arns, platformAppArn) {
			return fmt.Errorf("created platform application not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetPlatformApplicationAttributes", func() error {
		_, err := tc.client.SetPlatformApplicationAttributes(tc.ctx, &sns.SetPlatformApplicationAttributesInput{
			PlatformApplicationArn: aws.String(platformAppArn),
			Attributes: map[string]string{
				"EventEndpointCreated": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("set: %v", err)
		}

		getResp, err := tc.getPlatformApplicationAttributes(platformAppArn)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes["EventEndpointCreated"] != "true" {
			return fmt.Errorf("EventEndpointCreated not set: got %q", getResp.Attributes["EventEndpointCreated"])
		}
		return nil
	}))

	var endpointArn string
	results = append(results, r.RunTest("sns", "CreatePlatformEndpoint", func() error {
		resp, err := tc.client.CreatePlatformEndpoint(tc.ctx, &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: aws.String(platformAppArn),
			Token:                  aws.String("fake-device-token-12345"),
			CustomUserData:         aws.String("user-data-123"),
		})
		if err != nil {
			return err
		}
		if resp.EndpointArn == nil {
			return fmt.Errorf("EndpointArn is nil")
		}
		endpointArn = *resp.EndpointArn

		getResp, err := tc.getEndpointAttributes(endpointArn)
		if err != nil {
			return fmt.Errorf("get attrs: %v", err)
		}
		if getResp.Attributes == nil {
			return fmt.Errorf("Attributes is nil")
		}
		if getResp.Attributes["Token"] != "fake-device-token-12345" {
			return fmt.Errorf("token mismatch: got %q", getResp.Attributes["Token"])
		}
		if getResp.Attributes["Enabled"] != "true" {
			return fmt.Errorf("enabled should be true by default, got %q", getResp.Attributes["Enabled"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "ListEndpointsByPlatformApplication", func() error {
		arns, err := tc.listAllEndpointsByPlatformApplication(platformAppArn)
		if err != nil {
			return err
		}
		if len(arns) == 0 {
			return fmt.Errorf("expected at least one endpoint")
		}
		if !slices.Contains(arns, endpointArn) {
			return fmt.Errorf("created endpoint not found in ListEndpointsByPlatformApplication")
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "SetEndpointAttributes", func() error {
		_, err := tc.client.SetEndpointAttributes(tc.ctx, &sns.SetEndpointAttributesInput{
			EndpointArn: aws.String(endpointArn),
			Attributes: map[string]string{
				"CustomUserData": "updated-user-data",
				"Enabled":        "false",
			},
		})
		if err != nil {
			return fmt.Errorf("set: %v", err)
		}

		getResp, err := tc.getEndpointAttributes(endpointArn)
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Attributes["CustomUserData"] != "updated-user-data" {
			return fmt.Errorf("CustomUserData mismatch: got %q", getResp.Attributes["CustomUserData"])
		}
		if getResp.Attributes["Enabled"] != "false" {
			return fmt.Errorf("enabled should be false, got %q", getResp.Attributes["Enabled"])
		}
		return nil
	}))

	results = append(results, r.RunTest("sns", "DeletePlatformApplication_Cascade", func() error {
		cascadeAppName := tc.uniqueName("CascadeApp")
		appResp, err := tc.client.CreatePlatformApplication(tc.ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(cascadeAppName),
			Platform: aws.String("APNS"),
			Attributes: map[string]string{
				"PlatformCredential": "cascade-cred",
			},
		})
		if err != nil {
			return fmt.Errorf("create app: %v", err)
		}
		cascadeAppArn := *appResp.PlatformApplicationArn

		epResp, err := tc.client.CreatePlatformEndpoint(tc.ctx, &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: aws.String(cascadeAppArn),
			Token:                  aws.String("cascade-token-123"),
		})
		if err != nil {
			return fmt.Errorf("create endpoint: %v", err)
		}
		cascadeEndpointArn := *epResp.EndpointArn

		_, err = tc.client.DeletePlatformApplication(tc.ctx, &sns.DeletePlatformApplicationInput{
			PlatformApplicationArn: aws.String(cascadeAppArn),
		})
		if err != nil {
			return fmt.Errorf("delete app: %v", err)
		}

		_, err = tc.getEndpointAttributes(cascadeEndpointArn)
		return tc.expectNotFound("GetEndpointAttributes", err)
	}))

	results = append(results, r.RunTest("sns", "DeleteEndpoint", func() error {
		delEpAppName := tc.uniqueName("DelEpApp")
		appResp, err := tc.client.CreatePlatformApplication(tc.ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(delEpAppName),
			Platform: aws.String("GCM"),
			Attributes: map[string]string{
				"PlatformCredential": "del-ep-cred",
			},
		})
		if err != nil {
			return fmt.Errorf("create app: %v", err)
		}
		defer tc.client.DeletePlatformApplication(tc.ctx, &sns.DeletePlatformApplicationInput{
			PlatformApplicationArn: appResp.PlatformApplicationArn,
		})

		epResp, err := tc.client.CreatePlatformEndpoint(tc.ctx, &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: appResp.PlatformApplicationArn,
			Token:                  aws.String("del-ep-token"),
		})
		if err != nil {
			return fmt.Errorf("create endpoint: %v", err)
		}

		_, err = tc.client.DeleteEndpoint(tc.ctx, &sns.DeleteEndpointInput{
			EndpointArn: epResp.EndpointArn,
		})
		if err != nil {
			return fmt.Errorf("delete endpoint: %v", err)
		}

		_, err = tc.getEndpointAttributes(*epResp.EndpointArn)
		return tc.expectNotFound("GetEndpointAttributes", err)
	}))

	results = append(results, r.RunTest("sns", "DeletePlatformApplication_Standalone", func() error {
		appName := tc.uniqueName("DelApp")
		appResp, err := tc.client.CreatePlatformApplication(tc.ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(appName),
			Platform: aws.String("GCM"),
			Attributes: map[string]string{
				"PlatformCredential": "standalone-cred",
			},
		})
		if err != nil {
			return fmt.Errorf("create app: %v", err)
		}
		appArn := *appResp.PlatformApplicationArn

		_, err = tc.client.DeletePlatformApplication(tc.ctx, &sns.DeletePlatformApplicationInput{
			PlatformApplicationArn: appResp.PlatformApplicationArn,
		})
		if err != nil {
			return fmt.Errorf("delete app: %v", err)
		}

		_, err = tc.getPlatformApplicationAttributes(appArn)
		return tc.expectNotFound("GetPlatformApplicationAttributes", err)
	}))

	results = append(results, r.RunTest("sns", "CreatePlatformEndpoint_DuplicateToken", func() error {
		appName := tc.uniqueName("DupEpApp")
		appResp, err := tc.client.CreatePlatformApplication(tc.ctx, &sns.CreatePlatformApplicationInput{
			Name:     aws.String(appName),
			Platform: aws.String("GCM"),
			Attributes: map[string]string{
				"PlatformCredential": "dup-ep-cred",
			},
		})
		if err != nil {
			return fmt.Errorf("create app: %v", err)
		}
		defer tc.client.DeletePlatformApplication(tc.ctx, &sns.DeletePlatformApplicationInput{
			PlatformApplicationArn: appResp.PlatformApplicationArn,
		})

		ep1, err := tc.client.CreatePlatformEndpoint(tc.ctx, &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: appResp.PlatformApplicationArn,
			Token:                  aws.String("dup-token-999"),
		})
		if err != nil {
			return fmt.Errorf("create endpoint 1: %v", err)
		}

		ep2, err := tc.client.CreatePlatformEndpoint(tc.ctx, &sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: appResp.PlatformApplicationArn,
			Token:                  aws.String("dup-token-999"),
		})
		if err != nil {
			return fmt.Errorf("create endpoint 2 with same token: %v", err)
		}
		if *ep1.EndpointArn != *ep2.EndpointArn {
			return fmt.Errorf("duplicate token should return same endpoint ARN: %q vs %q", *ep1.EndpointArn, *ep2.EndpointArn)
		}
		return nil
	}))

	if endpointArn != "" {
		tc.client.DeleteEndpoint(tc.ctx, &sns.DeleteEndpointInput{EndpointArn: aws.String(endpointArn)})
	}
	if platformAppArn != "" {
		tc.client.DeletePlatformApplication(tc.ctx, &sns.DeletePlatformApplicationInput{
			PlatformApplicationArn: aws.String(platformAppArn),
		})
	}

	return results
}
