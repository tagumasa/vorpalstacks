package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTStreamRegistrationTests covers the Stream CRUD lifecycle and
// Registration Code lifecycle. These operations were previously unregistered
// stubs and are now real handlers backed by GenericKV persistence.
func (r *TestRunner) runIoTStreamRegistrationTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	streamID := uniqueName("test-stream")

	// --- Stream lifecycle ---
	results = append(results, r.RunTest("iot", "Stream_CreateStream", func() error {
		_, err := tc.client.CreateStream(tc.ctx, &iot.CreateStreamInput{
			StreamId:    aws.String(streamID),
			Description: aws.String("test stream"),
			RoleArn:     aws.String(fmt.Sprintf("arn:aws:iam::%s:role/iot-stream-role", tc.accountID)),
			Files: []iottypes.StreamFile{{
				FileId:     aws.Int32(1),
				S3Location: &iottypes.S3Location{Bucket: aws.String("test-bucket"), Key: aws.String("test-key")},
			}},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Stream_DescribeStream", func() error {
		out, err := tc.client.DescribeStream(tc.ctx, &iot.DescribeStreamInput{
			StreamId: aws.String(streamID),
		})
		if err != nil {
			return err
		}
		if out.StreamInfo == nil || out.StreamInfo.StreamId == nil {
			return fmt.Errorf("streamInfo missing in response")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Stream_DescribeStream_NotFound", func() error {
		_, err := tc.client.DescribeStream(tc.ctx, &iot.DescribeStreamInput{
			StreamId: aws.String(uniqueName("nope-stream")),
		})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "Stream_UpdateStream", func() error {
		_, err := tc.client.UpdateStream(tc.ctx, &iot.UpdateStreamInput{
			StreamId:    aws.String(streamID),
			Description: aws.String("updated description"),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Stream_ListStreams_IncludesCreated", func() error {
		out, err := tc.client.ListStreams(tc.ctx, &iot.ListStreamsInput{})
		if err != nil {
			return err
		}
		for _, s := range out.Streams {
			if s.StreamId != nil && *s.StreamId == streamID {
				return nil
			}
		}
		return fmt.Errorf("created stream not found in ListStreams")
	}))

	results = append(results, r.RunTest("iot", "Stream_DeleteStream", func() error {
		_, err := tc.client.DeleteStream(tc.ctx, &iot.DeleteStreamInput{
			StreamId: aws.String(streamID),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Stream_DeleteStream_NotFound", func() error {
		_, err := tc.client.DeleteStream(tc.ctx, &iot.DeleteStreamInput{
			StreamId: aws.String(uniqueName("nope-stream")),
		})
		return expectNotFound(err)
	}))

	// --- Registration code lifecycle ---
	results = append(results, r.RunTest("iot", "RegistrationCode_Get", func() error {
		out, err := tc.client.GetRegistrationCode(tc.ctx, &iot.GetRegistrationCodeInput{})
		if err != nil {
			return err
		}
		if out.RegistrationCode == nil || *out.RegistrationCode == "" {
			return fmt.Errorf("registrationCode missing in response")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "RegistrationCode_Get_Idempotent", func() error {
		out1, err := tc.client.GetRegistrationCode(tc.ctx, &iot.GetRegistrationCodeInput{})
		if err != nil {
			return err
		}
		out2, err := tc.client.GetRegistrationCode(tc.ctx, &iot.GetRegistrationCodeInput{})
		if err != nil {
			return err
		}
		if out1.RegistrationCode == nil || out2.RegistrationCode == nil {
			return fmt.Errorf("registrationCode missing in response")
		}
		if *out1.RegistrationCode != *out2.RegistrationCode {
			return fmt.Errorf("registration code not idempotent across calls")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "RegistrationCode_Delete", func() error {
		_, err := tc.client.DeleteRegistrationCode(tc.ctx, &iot.DeleteRegistrationCodeInput{})
		return err
	}))

	results = append(results, r.RunTest("iot", "RegistrationCode_Get_AfterDelete_NewCode", func() error {
		out, err := tc.client.GetRegistrationCode(tc.ctx, &iot.GetRegistrationCodeInput{})
		if err != nil {
			return err
		}
		if out.RegistrationCode == nil || *out.RegistrationCode == "" {
			return fmt.Errorf("registrationCode missing after delete")
		}
		return nil
	}))

	return results
}
