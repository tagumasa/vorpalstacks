package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func (r *TestRunner) runSESv2ConfigSetTests(tc *sesv2TestContext) []TestResult {
	var results []TestResult

	configSetName := fmt.Sprintf("test-cs-%d", tc.uid)

	results = append(results, r.RunTest("sesv2", "CreateConfigurationSet", func() error {
		_, err := tc.client.CreateConfigurationSet(tc.ctx, &sesv2.CreateConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetConfigurationSet(tc.ctx, &sesv2.GetConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return fmt.Errorf("get after create: %v", err)
		}
		if resp.ConfigurationSetName == nil || *resp.ConfigurationSetName != configSetName {
			return fmt.Errorf("expected name %s, got %v", configSetName, resp.ConfigurationSetName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetConfigurationSet", func() error {
		resp, err := tc.client.GetConfigurationSet(tc.ctx, &sesv2.GetConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return err
		}
		if resp.ConfigurationSetName == nil || *resp.ConfigurationSetName != configSetName {
			return fmt.Errorf("expected %s, got %v", configSetName, resp.ConfigurationSetName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListConfigurationSets", func() error {
		all, err := tc.listAllConfigSets()
		if err != nil {
			return err
		}
		if !containsString(all, configSetName) {
			return fmt.Errorf("created config set %s not found in list", configSetName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetSendingOptions", func() error {
		_, err := tc.client.PutConfigurationSetSendingOptions(tc.ctx, &sesv2.PutConfigurationSetSendingOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			SendingEnabled:       false,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetConfigurationSet(tc.ctx, &sesv2.GetConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return fmt.Errorf("get after put sending options: %v", err)
		}
		if resp.SendingOptions == nil {
			return fmt.Errorf("SendingOptions is nil after put")
		}
		if resp.SendingOptions.SendingEnabled {
			return fmt.Errorf("expected SendingEnabled=false after setting false")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetReputationOptions", func() error {
		_, err := tc.client.PutConfigurationSetReputationOptions(tc.ctx, &sesv2.PutConfigurationSetReputationOptionsInput{
			ConfigurationSetName:     aws.String(configSetName),
			ReputationMetricsEnabled: false,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetConfigurationSet(tc.ctx, &sesv2.GetConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return fmt.Errorf("get after put reputation options: %v", err)
		}
		if resp.ReputationOptions == nil {
			return fmt.Errorf("ReputationOptions is nil after put")
		}
		if resp.ReputationOptions.ReputationMetricsEnabled {
			return fmt.Errorf("expected ReputationMetricsEnabled=false")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetDeliveryOptions", func() error {
		_, err := tc.client.PutConfigurationSetDeliveryOptions(tc.ctx, &sesv2.PutConfigurationSetDeliveryOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			MaxDeliverySeconds:   aws.Int64(30),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetConfigurationSet(tc.ctx, &sesv2.GetConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return fmt.Errorf("get after put delivery options: %v", err)
		}
		if resp.DeliveryOptions == nil {
			return fmt.Errorf("DeliveryOptions is nil after put")
		}
		if resp.DeliveryOptions.MaxDeliverySeconds == nil || *resp.DeliveryOptions.MaxDeliverySeconds != 30 {
			return fmt.Errorf("expected MaxDeliverySeconds=30, got %v", resp.DeliveryOptions.MaxDeliverySeconds)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetSuppressionOptions", func() error {
		reasons := []types.SuppressionListReason{
			types.SuppressionListReasonBounce,
			types.SuppressionListReasonComplaint,
		}
		_, err := tc.client.PutConfigurationSetSuppressionOptions(tc.ctx, &sesv2.PutConfigurationSetSuppressionOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			SuppressedReasons:    reasons,
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetTrackingOptions", func() error {
		_, err := tc.client.PutConfigurationSetTrackingOptions(tc.ctx, &sesv2.PutConfigurationSetTrackingOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			HttpsPolicy:          types.HttpsPolicyRequire,
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetVdmOptions", func() error {
		_, err := tc.client.PutConfigurationSetVdmOptions(tc.ctx, &sesv2.PutConfigurationSetVdmOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			VdmOptions: &types.VdmOptions{
				DashboardOptions: &types.DashboardOptions{
					EngagementMetrics: types.FeatureStatusEnabled,
				},
			},
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetArchivingOptions", func() error {
		_, err := tc.client.PutConfigurationSetArchivingOptions(tc.ctx, &sesv2.PutConfigurationSetArchivingOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			ArchiveArn:           aws.String("arn:aws:mailmanager:us-east-1:000000000000:archive/test"),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	eventDestName := fmt.Sprintf("test-dest-%d", tc.uid)

	results = append(results, r.RunTest("sesv2", "CreateConfigurationSetEventDestination", func() error {
		_, err := tc.client.CreateConfigurationSetEventDestination(tc.ctx, &sesv2.CreateConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(configSetName),
			EventDestinationName: aws.String(eventDestName),
			EventDestination: &types.EventDestinationDefinition{
				Enabled:            true,
				MatchingEventTypes: []types.EventType{types.EventTypeSend},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetConfigurationSetEventDestinations(tc.ctx, &sesv2.GetConfigurationSetEventDestinationsInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return fmt.Errorf("get event destinations after create: %v", err)
		}
		found := false
		for _, ed := range resp.EventDestinations {
			if ed.Name != nil && *ed.Name == eventDestName {
				found = true
				if !ed.Enabled {
					return fmt.Errorf("expected Enabled=true")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("event destination %s not found after creation", eventDestName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetConfigurationSetEventDestinations", func() error {
		resp, err := tc.client.GetConfigurationSetEventDestinations(tc.ctx, &sesv2.GetConfigurationSetEventDestinationsInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return err
		}
		if len(resp.EventDestinations) == 0 {
			return fmt.Errorf("expected at least 1 event destination")
		}
		found := false
		for _, ed := range resp.EventDestinations {
			if ed.Name != nil && *ed.Name == eventDestName {
				found = true
				hasSend := false
				for _, et := range ed.MatchingEventTypes {
					if et == types.EventTypeSend {
						hasSend = true
					}
				}
				if !hasSend {
					return fmt.Errorf("expected MatchingEventTypes to contain SEND")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("event destination %s not found", eventDestName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UpdateConfigurationSetEventDestination", func() error {
		_, err := tc.client.UpdateConfigurationSetEventDestination(tc.ctx, &sesv2.UpdateConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(configSetName),
			EventDestinationName: aws.String(eventDestName),
			EventDestination: &types.EventDestinationDefinition{
				Enabled:            true,
				MatchingEventTypes: []types.EventType{types.EventTypeSend, types.EventTypeDelivery},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetConfigurationSetEventDestinations(tc.ctx, &sesv2.GetConfigurationSetEventDestinationsInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return fmt.Errorf("get after update: %v", err)
		}
		for _, ed := range resp.EventDestinations {
			if ed.Name != nil && *ed.Name == eventDestName {
				hasDelivery := false
				for _, et := range ed.MatchingEventTypes {
					if et == types.EventTypeDelivery {
						hasDelivery = true
					}
				}
				if !hasDelivery {
					return fmt.Errorf("DELIVERY not in MatchingEventTypes after update")
				}
				return nil
			}
		}
		return fmt.Errorf("event destination %s not found after update", eventDestName)
	}))

	results = append(results, r.RunTest("sesv2", "DeleteConfigurationSetEventDestination", func() error {
		_, err := tc.client.DeleteConfigurationSetEventDestination(tc.ctx, &sesv2.DeleteConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(configSetName),
			EventDestinationName: aws.String(eventDestName),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetConfigurationSetEventDestinations(tc.ctx, &sesv2.GetConfigurationSetEventDestinationsInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return fmt.Errorf("get after delete: %v", err)
		}
		for _, ed := range resp.EventDestinations {
			if ed.Name != nil && *ed.Name == eventDestName {
				return fmt.Errorf("event destination should have been deleted")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "TagResource_ConfigSet", func() error {
		_, err := tc.client.TagResource(tc.ctx, &sesv2.TagResourceInput{
			ResourceArn: aws.String(tc.configSetARN(configSetName)),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListTagsForResource_ConfigSet", func() error {
		resp, err := tc.client.ListTagsForResource(tc.ctx, &sesv2.ListTagsForResourceInput{
			ResourceArn: aws.String(tc.configSetARN(configSetName)),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) == 0 {
			return fmt.Errorf("expected at least 1 tag")
		}
		found := false
		for _, t := range resp.Tags {
			if t.Key != nil && *t.Key == "Environment" && t.Value != nil && *t.Value == "test" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("tag Environment=test not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UntagResource_ConfigSet", func() error {
		_, err := tc.client.UntagResource(tc.ctx, &sesv2.UntagResourceInput{
			ResourceArn: aws.String(tc.configSetARN(configSetName)),
			TagKeys:     []string{"Environment"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListTagsForResource(tc.ctx, &sesv2.ListTagsForResourceInput{
			ResourceArn: aws.String(tc.configSetARN(configSetName)),
		})
		if err != nil {
			return fmt.Errorf("list tags after untag: %v", err)
		}
		for _, t := range resp.Tags {
			if t.Key != nil && *t.Key == "Environment" {
				return fmt.Errorf("tag Environment should have been removed")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteConfigurationSet", func() error {
		_, err := tc.client.DeleteConfigurationSet(tc.ctx, &sesv2.DeleteConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetConfigurationSet(tc.ctx, &sesv2.GetConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting config set")
		}
		return nil
	}))

	return results
}
