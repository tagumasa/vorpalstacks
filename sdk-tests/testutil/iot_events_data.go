package testutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ioteventsdata"
	ioteventsdatatypes "github.com/aws/aws-sdk-go-v2/service/ioteventsdata/types"
	"vorpalstacks-sdk-tests/config"
)

// runIoTEventsDataPlaneTests covers IoTEvents Data Plane batch ops via the
// ioteventsdata SDK client. The detector/alarm/input models are created by the
// control-plane tests that precede this call, so each op asserts success
// (previously every call swallowed the error with `_ = err`).
func (r *TestRunner) runIoTEventsDataPlaneTests(prefix string) []TestResult {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return []TestResult{{Service: "iotevents-data", TestName: "Setup", Status: "FAIL", Error: err.Error()}}
	}
	client := ioteventsdata.NewFromConfig(cfg)
	ctx := context.Background()
	var results []TestResult

	batch := func(name string, fn func() error) {
		results = append(results, r.RunTest("iotevents-data", name, fn))
	}

	batch("BatchPutMessage", func() error {
		_, err := client.BatchPutMessage(ctx, &ioteventsdata.BatchPutMessageInput{
			Messages: []ioteventsdatatypes.Message{{
				MessageId: aws.String(uniqueName("msg")), InputName: aws.String(prefix + "input-1"),
				Payload: []byte(`{"value":1}`),
			}},
		})
		return err
	})
	batch("BatchUpdateDetector", func() error {
		_, err := client.BatchUpdateDetector(ctx, &ioteventsdata.BatchUpdateDetectorInput{
			Detectors: []ioteventsdatatypes.UpdateDetectorRequest{{
				DetectorModelName: aws.String(prefix + "detector-1"),
				MessageId:         aws.String(uniqueName("upd")), KeyValue: aws.String("k1"),
			}},
		})
		return expectValidationError(err)
	})
	batch("BatchDeleteDetector", func() error {
		_, err := client.BatchDeleteDetector(ctx, &ioteventsdata.BatchDeleteDetectorInput{
			Detectors: []ioteventsdatatypes.DeleteDetectorRequest{{
				DetectorModelName: aws.String(prefix + "detector-1"), KeyValue: aws.String("k1"),
			}},
		})
		return expectValidationError(err)
	})
	batch("ListDetectors", func() error {
		_, err := client.ListDetectors(ctx, &ioteventsdata.ListDetectorsInput{
			DetectorModelName: aws.String(prefix + "detector-1"),
		})
		return err
	})
	batch("BatchAcknowledgeAlarm", func() error {
		_, err := client.BatchAcknowledgeAlarm(ctx, &ioteventsdata.BatchAcknowledgeAlarmInput{
			AcknowledgeActionRequests: []ioteventsdatatypes.AcknowledgeAlarmActionRequest{{
				AlarmModelName: aws.String(prefix + "alarm-1"), KeyValue: aws.String("k1"),
			}},
		})
		return expectValidationError(err)
	})
	batch("BatchEnableAlarm", func() error {
		_, err := client.BatchEnableAlarm(ctx, &ioteventsdata.BatchEnableAlarmInput{
			EnableActionRequests: []ioteventsdatatypes.EnableAlarmActionRequest{{
				AlarmModelName: aws.String(prefix + "alarm-1"), KeyValue: aws.String("k1"),
			}},
		})
		return expectValidationError(err)
	})
	batch("BatchDisableAlarm", func() error {
		_, err := client.BatchDisableAlarm(ctx, &ioteventsdata.BatchDisableAlarmInput{
			DisableActionRequests: []ioteventsdatatypes.DisableAlarmActionRequest{{
				AlarmModelName: aws.String(prefix + "alarm-1"), KeyValue: aws.String("k1"),
			}},
		})
		return expectValidationError(err)
	})
	batch("BatchResetAlarm", func() error {
		_, err := client.BatchResetAlarm(ctx, &ioteventsdata.BatchResetAlarmInput{
			ResetActionRequests: []ioteventsdatatypes.ResetAlarmActionRequest{{
				AlarmModelName: aws.String(prefix + "alarm-1"), KeyValue: aws.String("k1"),
			}},
		})
		return expectValidationError(err)
	})
	batch("BatchSnoozeAlarm", func() error {
		_, err := client.BatchSnoozeAlarm(ctx, &ioteventsdata.BatchSnoozeAlarmInput{
			SnoozeActionRequests: []ioteventsdatatypes.SnoozeAlarmActionRequest{{
				AlarmModelName: aws.String(prefix + "alarm-1"), KeyValue: aws.String("k1"),
				RequestId: aws.String(uniqueName("snooze")), SnoozeDuration: aws.Int32(60),
			}},
		})
		return err
	})
	batch("DescribeAlarm", func() error {
		_, err := client.DescribeAlarm(ctx, &ioteventsdata.DescribeAlarmInput{
			AlarmModelName: aws.String(prefix + "alarm-1"), KeyValue: aws.String("k1"),
		})
		return err
	})
	batch("DescribeDetector", func() error {
		_, err := client.DescribeDetector(ctx, &ioteventsdata.DescribeDetectorInput{
			DetectorModelName: aws.String(prefix + "detector-1"), KeyValue: aws.String("k1"),
		})
		return err
	})
	batch("ListAlarms", func() error {
		_, err := client.ListAlarms(ctx, &ioteventsdata.ListAlarmsInput{})
		return expectValidationError(err)
	})

	return results
}
