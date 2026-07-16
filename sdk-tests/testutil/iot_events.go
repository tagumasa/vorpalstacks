package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotevents"
	ioteventstypes "github.com/aws/aws-sdk-go-v2/service/iotevents/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"vorpalstacks-sdk-tests/config"
)

type iotEventsTestContext struct {
	client *iotevents.Client
	ctx    context.Context
}

func (r *TestRunner) newIoTEventsTestContext() (*iotEventsTestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &iotEventsTestContext{
		client: iotevents.NewFromConfig(cfg),
		ctx:    context.Background(),
	}, nil
}

func (r *TestRunner) RunIoTEventsTests() []TestResult {
	tc, err := r.newIoTEventsTestContext()
	if err != nil {
		return []TestResult{{
			Service:  "iotevents",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		}}
	}

	// Use a unique prefix so parallel runs and re-runs never collide on
	// fixed names like "vstest-input-1".
	prefix := uniqueName("iotevents") + "-"

	// Best-effort cleanup: defer at function level guarantees resource removal
	// even when individual test cases fail or panic.
	defer tc.client.DeleteInput(tc.ctx, &iotevents.DeleteInputInput{
		InputName: ptrString(prefix + "input-1"),
	})
	defer tc.client.DeleteDetectorModel(tc.ctx, &iotevents.DeleteDetectorModelInput{
		DetectorModelName: ptrString(prefix + "detector-1"),
	})
	defer tc.client.DeleteAlarmModel(tc.ctx, &iotevents.DeleteAlarmModelInput{
		AlarmModelName: ptrString(prefix + "alarm-1"),
	})

	var results []TestResult

	results = append(results, r.RunTest("iotevents", "Input_CreateInput", func() error {
		out, err := tc.client.CreateInput(tc.ctx, &iotevents.CreateInputInput{
			InputName: ptrString(prefix + "input-1"),
			InputDefinition: &ioteventstypes.InputDefinition{
				Attributes: []ioteventstypes.Attribute{
					{JsonPath: ptrString("$.value")},
				},
			},
		})
		if err != nil {
			return err
		}
		if out.InputConfiguration == nil {
			return fmt.Errorf("expected non-nil InputConfiguration in CreateInput response")
		}
		return nil
	}))

	results = append(results, r.RunTest("iotevents", "Input_DescribeInput", func() error {
		out, err := tc.client.DescribeInput(tc.ctx, &iotevents.DescribeInputInput{
			InputName: ptrString(prefix + "input-1"),
		})
		if err != nil {
			return err
		}
		if out.Input == nil || out.Input.InputConfiguration == nil {
			return fmt.Errorf("DescribeInput returned nil Input or InputConfiguration")
		}
		if out.Input.InputConfiguration.InputName == nil || *out.Input.InputConfiguration.InputName != prefix+"input-1" {
			return fmt.Errorf("input name mismatch: got %v", out.Input.InputConfiguration.InputName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iotevents", "Input_ListInputs", func() error {
		var token *string
		for {
			out, err := tc.client.ListInputs(tc.ctx, &iotevents.ListInputsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, in := range out.InputSummaries {
				if in.InputName != nil && *in.InputName == prefix+"input-1" {
					return nil
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return fmt.Errorf("ListInputs did not contain test input")
	}))

	results = append(results, r.RunTest("iotevents", "Input_UpdateInput", func() error {
		_, err := tc.client.UpdateInput(tc.ctx, &iotevents.UpdateInputInput{
			InputName: ptrString(prefix + "input-1"),
			InputDefinition: &ioteventstypes.InputDefinition{
				Attributes: []ioteventstypes.Attribute{
					{JsonPath: ptrString("$.value")},
					{JsonPath: ptrString("$.timestamp")},
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iotevents", "Input_DeleteInput", func() error {
		_, err := tc.client.DeleteInput(tc.ctx, &iotevents.DeleteInputInput{
			InputName: ptrString(prefix + "input-1"),
		})
		return err
	}))

	results = append(results, r.RunTest("iotevents", "DetectorModel_CreateDetectorModel", func() error {
		out, err := tc.client.CreateDetectorModel(tc.ctx, &iotevents.CreateDetectorModelInput{
			DetectorModelName: ptrString(prefix + "detector-1"),
			RoleArn:           ptrString("arn:aws:iam::000000000000:role/test"),
			DetectorModelDefinition: &ioteventstypes.DetectorModelDefinition{
				InitialStateName: ptrString("Idle"),
				States: []ioteventstypes.State{{
					StateName: ptrString("Idle"),
					OnEnter: &ioteventstypes.OnEnterLifecycle{
						Events: []ioteventstypes.Event{{
							EventName: ptrString("Init"),
							Condition: ptrString("true"),
							Actions:   []ioteventstypes.Action{},
						}},
					},
				}},
			},
		})
		if err != nil {
			return err
		}
		if out.DetectorModelConfiguration == nil {
			return fmt.Errorf("expected non-nil DetectorModelConfiguration in response")
		}
		return nil
	}))

	results = append(results, r.RunTest("iotevents", "DetectorModel_DescribeDetectorModel", func() error {
		out, err := tc.client.DescribeDetectorModel(tc.ctx, &iotevents.DescribeDetectorModelInput{
			DetectorModelName: ptrString(prefix + "detector-1"),
		})
		if err != nil {
			return err
		}
		if out.DetectorModel == nil || out.DetectorModel.DetectorModelConfiguration == nil {
			return fmt.Errorf("DescribeDetectorModel returned nil DetectorModel or Configuration")
		}
		if out.DetectorModel.DetectorModelConfiguration.DetectorModelName == nil ||
			*out.DetectorModel.DetectorModelConfiguration.DetectorModelName != prefix+"detector-1" {
			return fmt.Errorf("detector model name mismatch: got %v", out.DetectorModel.DetectorModelConfiguration.DetectorModelName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iotevents", "DetectorModel_ListDetectorModels", func() error {
		var token *string
		for {
			out, err := tc.client.ListDetectorModels(tc.ctx, &iotevents.ListDetectorModelsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, dm := range out.DetectorModelSummaries {
				if dm.DetectorModelName != nil && *dm.DetectorModelName == prefix+"detector-1" {
					return nil
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return fmt.Errorf("ListDetectorModels did not contain test detector")
	}))

	results = append(results, r.RunTest("iotevents", "DetectorModel_UpdateDetectorModel", func() error {
		_, err := tc.client.UpdateDetectorModel(tc.ctx, &iotevents.UpdateDetectorModelInput{
			DetectorModelName: ptrString(prefix + "detector-1"),
			RoleArn:           ptrString("arn:aws:iam::000000000000:role/test"),
			DetectorModelDefinition: &ioteventstypes.DetectorModelDefinition{
				InitialStateName: ptrString("Idle"),
				States: []ioteventstypes.State{{
					StateName: ptrString("Idle"),
					OnEnter: &ioteventstypes.OnEnterLifecycle{
						Events: []ioteventstypes.Event{{
							EventName: ptrString("Init"),
							Condition: ptrString("true"),
							Actions:   []ioteventstypes.Action{},
						}},
					},
				}},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iotevents", "DetectorModel_DeleteDetectorModel", func() error {
		_, err := tc.client.DeleteDetectorModel(tc.ctx, &iotevents.DeleteDetectorModelInput{
			DetectorModelName: ptrString(prefix + "detector-1"),
		})
		return err
	}))

	results = append(results, r.RunTest("iotevents", "BatchPutMessage_DetectorTransition", func() error {
		return r.testBatchPutMessageDetectorTransition(tc, prefix)
	}))

	results = append(results, r.RunTest("iotevents", "BatchPutMessage_IotEventsAction_NoDeadlock", func() error {
		return r.testIotEventsActionNoDeadlock(tc, prefix)
	}))

	results = append(results, r.runIoTEventsAlarmModelTests(tc, prefix)...)
	results = append(results, r.runIoTEventsDataPlaneTests(prefix)...)
	results = append(results, r.runIoTEventsListDescribeTests(tc, prefix)...)

	return results
}

func (r *TestRunner) runIoTEventsAlarmModelTests(tc *iotEventsTestContext, prefix string) []TestResult {
	var results []TestResult
	alarmName := prefix + "alarm-1"

	results = append(results, r.RunTest("iotevents", "AlarmModel_CreateAlarmModel", func() error {
		out, err := tc.client.CreateAlarmModel(tc.ctx, &iotevents.CreateAlarmModelInput{
			AlarmModelName: ptrString(alarmName),
			RoleArn:        ptrString("arn:aws:iam::000000000000:role/test"),
			AlarmRule: &ioteventstypes.AlarmRule{
				SimpleRule: &ioteventstypes.SimpleRule{
					InputProperty:      ptrString("$.value"),
					ComparisonOperator: ioteventstypes.ComparisonOperatorGreater,
					Threshold:          ptrString("0"),
				},
			},
		})
		if err != nil {
			return err
		}
		if out.AlarmModelArn == nil || *out.AlarmModelArn == "" {
			return fmt.Errorf("expected non-empty AlarmModelArn in CreateAlarmModel response")
		}
		return nil
	}))

	results = append(results, r.RunTest("iotevents", "AlarmModel_DescribeAlarmModel", func() error {
		out, err := tc.client.DescribeAlarmModel(tc.ctx, &iotevents.DescribeAlarmModelInput{
			AlarmModelName: ptrString(alarmName),
		})
		if err != nil {
			return err
		}
		if out.AlarmModelName == nil || *out.AlarmModelName != alarmName {
			return fmt.Errorf("expected alarmModelName=%s, got %v", alarmName, out.AlarmModelName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iotevents", "AlarmModel_UpdateAlarmModel", func() error {
		_, err := tc.client.UpdateAlarmModel(tc.ctx, &iotevents.UpdateAlarmModelInput{
			AlarmModelName:        ptrString(alarmName),
			AlarmModelDescription: ptrString("updated description"),
			RoleArn:               ptrString("arn:aws:iam::000000000000:role/test"),
			AlarmRule: &ioteventstypes.AlarmRule{
				SimpleRule: &ioteventstypes.SimpleRule{
					InputProperty:      ptrString("$.value"),
					ComparisonOperator: ioteventstypes.ComparisonOperatorGreater,
					Threshold:          ptrString("1"),
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iotevents", "AlarmModel_ListAlarmModels", func() error {
		var token *string
		for {
			out, err := tc.client.ListAlarmModels(tc.ctx, &iotevents.ListAlarmModelsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, am := range out.AlarmModelSummaries {
				if am.AlarmModelName != nil && *am.AlarmModelName == alarmName {
					return nil
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return fmt.Errorf("ListAlarmModels did not contain test alarm %s", alarmName)
	}))

	results = append(results, r.RunTest("iotevents", "AlarmModel_ListAlarmModelVersions", func() error {
		out, err := tc.client.ListAlarmModelVersions(tc.ctx, &iotevents.ListAlarmModelVersionsInput{
			AlarmModelName: ptrString(alarmName),
		})
		if err != nil {
			return err
		}
		if len(out.AlarmModelVersionSummaries) == 0 {
			return fmt.Errorf("expected at least one alarm model version")
		}
		return nil
	}))

	results = append(results, r.RunTest("iotevents", "AlarmModel_DeleteAlarmModel", func() error {
		_, err := tc.client.DeleteAlarmModel(tc.ctx, &iotevents.DeleteAlarmModelInput{
			AlarmModelName: ptrString(alarmName),
		})
		return err
	}))

	results = append(results, r.RunTest("iotevents", "Analysis_StartDetectorModelAnalysis", func() error {
		_, err := tc.client.StartDetectorModelAnalysis(tc.ctx, &iotevents.StartDetectorModelAnalysisInput{
			DetectorModelDefinition: &ioteventstypes.DetectorModelDefinition{
				InitialStateName: ptrString("Idle"),
				States: []ioteventstypes.State{{
					StateName: ptrString("Idle"),
					OnEnter: &ioteventstypes.OnEnterLifecycle{
						Events: []ioteventstypes.Event{{
							EventName: ptrString("Init"),
							Condition: ptrString("true"),
							Actions:   []ioteventstypes.Action{},
						}},
					},
				}},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iotevents", "Analysis_DescribeDetectorModelAnalysis", func() error {
		startOut, err := tc.client.StartDetectorModelAnalysis(tc.ctx, &iotevents.StartDetectorModelAnalysisInput{
			DetectorModelDefinition: &ioteventstypes.DetectorModelDefinition{
				InitialStateName: ptrString("Idle"),
				States: []ioteventstypes.State{{
					StateName: ptrString("Idle"),
					OnEnter: &ioteventstypes.OnEnterLifecycle{
						Events: []ioteventstypes.Event{{
							EventName: ptrString("Init"),
							Condition: ptrString("true"),
							Actions:   []ioteventstypes.Action{},
						}},
					},
				}},
			},
		})
		if err != nil {
			return fmt.Errorf("StartDetectorModelAnalysis setup failed: %w", err)
		}
		if startOut.AnalysisId == nil {
			return nil
		}
		_, err = tc.client.DescribeDetectorModelAnalysis(tc.ctx, &iotevents.DescribeDetectorModelAnalysisInput{
			AnalysisId: startOut.AnalysisId,
		})
		return err
	}))

	results = append(results, r.RunTest("iotevents", "Logging_PutLoggingOptions", func() error {
		_, err := tc.client.PutLoggingOptions(tc.ctx, &iotevents.PutLoggingOptionsInput{
			LoggingOptions: &ioteventstypes.LoggingOptions{
				Enabled: true,
				Level:   "INFO",
				RoleArn: ptrString("arn:aws:iam::000000000000:role/test"),
			},
		})
		return err
	}))

	return results
}

func (r *TestRunner) testBatchPutMessageDetectorTransition(tc *iotEventsTestContext, prefix string) error {
	sqsCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return fmt.Errorf("failed to create SQS config: %w", err)
	}
	sqsClient := sqs.NewFromConfig(sqsCfg)
	sqsCtx := context.Background()

	queueName := fmt.Sprintf("detector-test-%d", time.Now().UnixNano())

	createResp, err := sqsClient.CreateQueue(sqsCtx, &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return fmt.Errorf("create SQS queue: %w", err)
	}
	queueURL := aws.ToString(createResp.QueueUrl)

	detectorName := prefix + "pipeline-detector"

	defer func() {
		sqsClient.DeleteQueue(sqsCtx, &sqs.DeleteQueueInput{
			QueueUrl: aws.String(queueURL),
		})
		tc.client.DeleteDetectorModel(tc.ctx, &iotevents.DeleteDetectorModelInput{
			DetectorModelName: ptrString(detectorName),
		})
	}()

	tc.client.DeleteDetectorModel(tc.ctx, &iotevents.DeleteDetectorModelInput{
		DetectorModelName: ptrString(detectorName),
	})

	_, err = tc.client.CreateDetectorModel(tc.ctx, &iotevents.CreateDetectorModelInput{
		DetectorModelName: ptrString(detectorName),
		RoleArn:           ptrString("arn:aws:iam::000000000000:role/test"),
		DetectorModelDefinition: &ioteventstypes.DetectorModelDefinition{
			InitialStateName: ptrString("Idle"),
			States: []ioteventstypes.State{
				{
					StateName: ptrString("Idle"),
					OnInput: &ioteventstypes.OnInputLifecycle{
						TransitionEvents: []ioteventstypes.TransitionEvent{
							{
								EventName: ptrString("HighTemp"),
								Condition: ptrString("true"),
								Actions: []ioteventstypes.Action{
									{
										Sqs: &ioteventstypes.SqsAction{
											QueueUrl: aws.String(queueURL),
										},
									},
								},
								NextState: ptrString("Alarm"),
							},
						},
					},
				},
				{
					StateName: ptrString("Alarm"),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("create detector model: %w", err)
	}

	payload := `{"temperature":75}`

	batchBody, _ := json.Marshal(map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"inputName":      "TempInput",
				"messagePayload": payload,
			},
		},
	})

	iotEndpoint := fmt.Sprintf("%s/messages", r.endpoint)
	req, err := http.NewRequest("POST", iotEndpoint, bytes.NewReader(batchBody))
	if err != nil {
		return fmt.Errorf("create BatchPutMessage request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("BatchPutMessage HTTP: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("BatchPutMessage HTTP %d: %s", resp.StatusCode, string(body))
	}

	receiveResp, err := sqsClient.ReceiveMessage(sqsCtx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     2,
	})
	if err != nil {
		return fmt.Errorf("ReceiveMessage from SQS: %w", err)
	}

	for _, msg := range receiveResp.Messages {
		if msg.Body != nil && len(*msg.Body) > 0 {
			return nil
		}
	}

	return fmt.Errorf("no message received from SQS queue within timeout")
}

func (r *TestRunner) testIotEventsActionNoDeadlock(tc *iotEventsTestContext, prefix string) error {
	// Create a detector model whose action is iotEvents (recursive dispatch).
	// Before the async fix, this would deadlock because EvaluateEvent holds
	// sm.mu.Lock() and the iotEvents action calls BatchEvaluate which needs
	// sm.mu.RLock() on the same goroutine. Verify the BatchPutMessage returns
	// within a few seconds instead of hanging indefinitely.
	detectorName := prefix + "iot-events-deadlock-test"

	defer func() {
		tc.client.DeleteDetectorModel(tc.ctx, &iotevents.DeleteDetectorModelInput{
			DetectorModelName: ptrString(detectorName),
		})
	}()

	tc.client.DeleteDetectorModel(tc.ctx, &iotevents.DeleteDetectorModelInput{
		DetectorModelName: ptrString(detectorName),
	})

	_, err := tc.client.CreateDetectorModel(tc.ctx, &iotevents.CreateDetectorModelInput{
		DetectorModelName: ptrString(detectorName),
		RoleArn:           ptrString("arn:aws:iam::000000000000:role/test"),
		DetectorModelDefinition: &ioteventstypes.DetectorModelDefinition{
			InitialStateName: ptrString("Idle"),
			States: []ioteventstypes.State{
				{
					StateName: ptrString("Idle"),
					OnInput: &ioteventstypes.OnInputLifecycle{
						TransitionEvents: []ioteventstypes.TransitionEvent{
							{
								EventName: ptrString("Trigger"),
								Condition: ptrString("true"),
								Actions: []ioteventstypes.Action{
									{
										IotEvents: &ioteventstypes.IotEventsAction{
											InputName: ptrString("RecursiveInput"),
										},
									},
								},
								NextState: ptrString("Triggered"),
							},
						},
					},
				},
				{
					StateName: ptrString("Triggered"),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("create detector model with iotEvents action: %w", err)
	}

	payload := `{"sensor":1}`
	batchBody, _ := json.Marshal(map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"inputName":      "TempInput",
				"messagePayload": payload,
			},
		},
	})

	iotEndpoint := fmt.Sprintf("%s/messages", r.endpoint)

	// Use a short client timeout so that a deadlock would surface as a timeout
	// rather than hanging the entire test suite.
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", iotEndpoint, bytes.NewReader(batchBody))
	if err != nil {
		return fmt.Errorf("create BatchPutMessage request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("BatchPutMessage HTTP (deadlocked?): %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("BatchPutMessage HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
