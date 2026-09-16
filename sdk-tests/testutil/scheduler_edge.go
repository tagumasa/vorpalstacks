package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmtypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

func (tc *schedTestContext) runEdgeTests() []TestResult {
	var results []TestResult

	rn, rARN := tc.createIAMRole()
	defer tc.deleteIAMRole(rn)

	// Addressing a schedule that does not exist is
	// ResourceNotFoundException on every mutating operation.
	results = append(results, tc.runner.RunTest("scheduler", "Schedule_NonExistent", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			{"GetSchedule", func() error {
				_, err := tc.getSchedule("nonexistent-schedule-xyz")
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"DeleteSchedule", func() error {
				_, err := tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{
					Name: aws.String("nonexistent-schedule-xyz"),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"UpdateSchedule", func() error {
				_, err := tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
					Name:               aws.String("nonexistent-schedule-xyz"),
					ScheduleExpression: aws.String("rate(30 minutes)"),
					Target: &types.Target{
						Arn:     aws.String(tc.lambdaARN()),
						RoleArn: aws.String(rARN),
					},
					FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %w", row.name, err)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_DuplicateName", func() error {
		dupName := tc.uniqueName("DupSched")
		dupRN, dupRARN := tc.createIAMRole()
		defer tc.deleteIAMRole(dupRN)

		_, err := tc.createSchedule(dupName, "rate(30 minutes)", tc.defaultTarget(dupRARN))
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.cleanupSchedule(dupName)

		_, err = tc.createSchedule(dupName, "rate(60 minutes)", tc.defaultTarget(dupRARN))
		if err := AssertErrorContains(err, "ConflictException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_InvalidExpression", func() error {
		// A syntactically impossible expression and an out-of-range cron
		// field are both rejected at creation: AWS validates every cron
		// field's range, and an accepted-but-never-matching expression
		// would silently never fire.
		invalidExpressions := []string{
			"not-a-valid-expression",
			"cron(99 12 * * ? *)",
		}
		for _, expr := range invalidExpressions {
			invName := tc.uniqueName("InvExprSched")
			defer tc.cleanupSchedule(invName)
			_, err := tc.createSchedule(invName, expr, &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
			})
			if err := AssertErrorContains(err, "ValidationException"); err != nil {
				return fmt.Errorf("expression %q: %v", expr, err)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_UnsupportedTarget_Rejected", func() error {
		schedName := tc.uniqueName("UnsupportedTarget")
		defer tc.cleanupSchedule(schedName)
		_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
			Arn:     aws.String(fmt.Sprintf("arn:aws:codebuild:%s:%s:project/FakeProject", tc.region, tc.accountID)),
			RoleArn: aws.String(rARN),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	// The universal-target ARN form (arn:aws:scheduler:::aws-sdk:{service}:{action})
	// is accepted at creation for any SDK service — delivery decides per
	// target, mirroring AWS — while the read-only action prefixes, non-JSON
	// Input and templated sub-parameters on a universal target are
	// creation-time rejections.
	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_UniversalTarget_Validation", func() error {
		schedName := tc.uniqueName("UniversalTarget")
		defer tc.cleanupSchedule(schedName)

		universalARN := "arn:aws:scheduler:::aws-sdk:batch:submitJob"
		_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
			Arn:     aws.String(universalARN),
			RoleArn: aws.String(rARN),
			Input:   aws.String(`{"JobName":"validation-pin"}`),
		})
		if err != nil {
			return fmt.Errorf("an unimplemented-service universal target must be creatable: %w", err)
		}
		got, err := tc.getSchedule(schedName)
		if err != nil {
			return fmt.Errorf("get schedule: %w", err)
		}
		if got.Target == nil || aws.ToString(got.Target.Arn) != universalARN {
			return fmt.Errorf("universal ARN must round-trip, got %+v", got.Target)
		}
		if _, err := tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)}); err != nil {
			return fmt.Errorf("delete schedule: %w", err)
		}

		rows := []struct {
			name   string
			target *types.Target
		}{
			{"read-only action prefix", &types.Target{
				Arn:     aws.String("arn:aws:scheduler:::aws-sdk:sqs:getQueueUrl"),
				RoleArn: aws.String(rARN),
			}},
			{"non-JSON Input", &types.Target{
				Arn:     aws.String("arn:aws:scheduler:::aws-sdk:lambda:invoke"),
				RoleArn: aws.String(rARN),
				Input:   aws.String(`not-json`),
			}},
			{"templated sub-parameters on a universal target", &types.Target{
				Arn:           aws.String("arn:aws:scheduler:::aws-sdk:sqs:sendMessage"),
				RoleArn:       aws.String(rARN),
				Input:         aws.String(`{"QueueUrl":"http://q","MessageBody":"m"}`),
				SqsParameters: &types.SqsParameters{MessageGroupId: aws.String("g")},
			}},
			{"malformed scheduler-service ARN", &types.Target{
				Arn:     aws.String("arn:aws:scheduler:::aws-sdk:sqs"),
				RoleArn: aws.String(rARN),
			}},
		}
		for _, row := range rows {
			if _, err := tc.createSchedule(tc.uniqueName("UniversalReject"), "rate(30 minutes)", row.target); err == nil ||
				!strings.Contains(err.Error(), "ValidationException") {
				return fmt.Errorf("%s: expected ValidationException, got %v", row.name, err)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_DeadLetterConfig_SnsRejected", func() error {
		schedName := tc.uniqueName("SnsDLQ")
		defer tc.cleanupSchedule(schedName)
		_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
			Arn:     aws.String(tc.lambdaARN()),
			RoleArn: aws.String(rARN),
			DeadLetterConfig: &types.DeadLetterConfig{
				Arn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:MyTopic", tc.region, tc.accountID)),
			},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	// The EcsParameters members carry one shared validation substrate:
	// every out-of-range, malformed or misplaced member is a
	// ValidationException at schedule creation.
	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EcsParameters_Validation", func() error {
		taskDef := fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/FakeTask:1", tc.region, tc.accountID)
		clusterArn := fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/ValidationCluster", tc.region, tc.accountID)
		rows := []struct {
			name   string
			target string
			ecs    func() *types.EcsParameters
		}{
			{"TaskCountOutOfRange", tc.lambdaARN(), func() *types.EcsParameters {
				taskCount := int32(100)
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef), TaskCount: &taskCount}
			}},
			{"SubnetNotFound", clusterArn, func() *types.EcsParameters {
				taskCount := int32(1)
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef), TaskCount: &taskCount,
					NetworkConfiguration: &types.NetworkConfiguration{
						AwsvpcConfiguration: &types.AwsVpcConfiguration{Subnets: []string{"subnet-nonexistent-xyz"}},
					}}
			}},
			{"EmptySecurityGroupsRejected", clusterArn, func() *types.EcsParameters {
				taskCount := int32(1)
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef), TaskCount: &taskCount,
					NetworkConfiguration: &types.NetworkConfiguration{
						AwsvpcConfiguration: &types.AwsVpcConfiguration{
							Subnets:        []string{"subnet-0123456789abcdef0"},
							SecurityGroups: []string{},
						},
					}}
			}},
			{"CapacityProviderWeightOutOfRange", clusterArn, func() *types.EcsParameters {
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef),
					CapacityProviderStrategy: []types.CapacityProviderStrategyItem{
						{CapacityProvider: aws.String("FARGATE"), Weight: 5000},
					}}
			}},
			{"PlacementStrategyInvalidType", clusterArn, func() *types.EcsParameters {
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef),
					PlacementStrategy: []types.PlacementStrategy{
						{Type: types.PlacementStrategyType("bogus")},
					}}
			}},
			{"PlacementExpressionTooLong", clusterArn, func() *types.EcsParameters {
				// PlacementConstraintExpression @length(max 2000).
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef),
					PlacementConstraints: []types.PlacementConstraint{
						{Type: types.PlacementConstraintType("memberOf"), Expression: aws.String(strings.Repeat("a", 2001))},
					}}
			}},
			{"PlacementFieldTooLong", clusterArn, func() *types.EcsParameters {
				// PlacementStrategyField @length(max 255).
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef),
					PlacementStrategy: []types.PlacementStrategy{
						{Type: types.PlacementStrategyType("spread"), Field: aws.String(strings.Repeat("a", 256))},
					}}
			}},
			{"TagsTooManyItems", clusterArn, func() *types.EcsParameters {
				// Tags @length(0, 50): the list of TagMap entries.
				tags := make([]map[string]string, 51)
				for i := range tags {
					tags[i] = map[string]string{"env": "prod"}
				}
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef), Tags: tags}
			}},
			{"TagsKeyTooLong", clusterArn, func() *types.EcsParameters {
				// TagKey @length(1, 128).
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef),
					Tags: []map[string]string{{strings.Repeat("k", 129): "v"}}}
			}},
			{"TagsValueTooLong", clusterArn, func() *types.EcsParameters {
				// TagValue @length(1, 256).
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef),
					Tags: []map[string]string{{"k": strings.Repeat("v", 257)}}}
			}},
			{"EmptySubnetsRejected", clusterArn, func() *types.EcsParameters {
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef),
					NetworkConfiguration: &types.NetworkConfiguration{
						AwsvpcConfiguration: &types.AwsVpcConfiguration{Subnets: []string{}},
					}}
			}},
			{"NonEcsTarget_Rejected", tc.lambdaARN(), func() *types.EcsParameters {
				return &types.EcsParameters{TaskDefinitionArn: aws.String(taskDef)}
			}},
		}
		for _, row := range rows {
			schedName := tc.uniqueName("EcsValidation")
			defer tc.cleanupSchedule(schedName)
			_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
				Arn:           aws.String(row.target),
				RoleArn:       aws.String(rARN),
				EcsParameters: row.ecs(),
			})
			if err := AssertErrorContains(err, "ValidationException"); err != nil {
				return fmt.Errorf("%s: %w", row.name, err)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_AsymmetricKmsKeyRejected", func() error {
		// AWS validates the key class at schedule creation via
		// kms:DescribeKey; an asymmetric key cannot protect a schedule.
		key, err := tc.kmsClient.CreateKey(tc.ctx, &kms.CreateKeyInput{
			KeySpec:  kmtypes.KeySpecRsa2048,
			KeyUsage: kmtypes.KeyUsageTypeEncryptDecrypt,
		})
		if err != nil {
			return fmt.Errorf("create asymmetric KMS key: %v", err)
		}
		defer func() {
			// ScheduleKeyDeletion requires a disabled key.
			tc.kmsClient.DisableKey(tc.ctx, &kms.DisableKeyInput{KeyId: key.KeyMetadata.KeyId})
			tc.kmsClient.ScheduleKeyDeletion(tc.ctx, &kms.ScheduleKeyDeletionInput{KeyId: key.KeyMetadata.KeyId})
		}()

		schedName := tc.uniqueName("AsymKms")
		defer tc.cleanupSchedule(schedName)
		_, err = tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			KmsKeyArn:          key.KeyMetadata.Arn,
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_KmsKeyNotFound", func() error {
		schedName := tc.uniqueName("BadKms")
		defer tc.cleanupSchedule(schedName)
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			// AWS validates the customer managed key at schedule
			// creation (kms:DescribeKey requirement in the encryption
			// documentation); a non-existent key must be rejected.
			KmsKeyArn: aws.String(fmt.Sprintf("arn:aws:kms:%s:%s:key/11111111-2222-3333-4444-555555555555", tc.region, tc.accountID)),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	// The schedule-level modelled traits are enforced at creation: KmsKeyArn
	// @pattern + @length(1,2048), the DeadLetterConfig.Arn member @pattern,
	// the FlexibleTimeWindow.MaximumWindowInMinutes @range(1,1440)
	// (unconditional on the shape, so a bound outside the range is rejected
	// in OFF mode too), and ScheduleExpressionTimezone taking IANA zone
	// names only.
	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_ModelTraitValidation", func() error {
		rows := []struct {
			name      string
			kmsKeyArn string
			timezone  string
		}{
			{"KmsKeyArnWrongResourceType", fmt.Sprintf("arn:aws:kms:%s:%s:secret/not-a-key", tc.region, tc.accountID), ""},
			{"KmsKeyArnOverLength", "arn:aws:kms:" + tc.region + ":" + tc.accountID + ":key/" + strings.Repeat("k", 2048), ""},
			{"TimezoneLocalRejected", "", "Local"},
		}
		for _, row := range rows {
			schedName := tc.uniqueName("ModelTrait")
			defer tc.cleanupSchedule(schedName)
			input := &scheduler.CreateScheduleInput{
				Name:               aws.String(schedName),
				ScheduleExpression: aws.String("rate(30 minutes)"),
				Target: &types.Target{
					Arn:     aws.String(tc.lambdaARN()),
					RoleArn: aws.String(rARN),
				},
				FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			}
			if row.kmsKeyArn != "" {
				input.KmsKeyArn = aws.String(row.kmsKeyArn)
			}
			if row.timezone != "" {
				input.ScheduleExpressionTimezone = aws.String(row.timezone)
			} else {
				input.ScheduleExpressionTimezone = aws.String("UTC")
			}
			_, err := tc.client.CreateSchedule(tc.ctx, input)
			if err := AssertErrorContains(err, "ValidationException"); err != nil {
				return fmt.Errorf("%s: %w", row.name, err)
			}
		}

		// DeadLetterConfig.Arn outside the member pattern (queue-name
		// charset is [a-zA-Z0-9-_]).
		schedName := tc.uniqueName("ModelTrait")
		defer tc.cleanupSchedule(schedName)
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
				DeadLetterConfig: &types.DeadLetterConfig{
					Arn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:bad name", tc.region, tc.accountID)),
				},
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return fmt.Errorf("DeadLetterConfigBadQueueCharset: %w", err)
		}

		// MaximumWindowInMinutes above the @range maximum in OFF mode.
		schedName = tc.uniqueName("ModelTrait")
		defer tc.cleanupSchedule(schedName)
		_, err = tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode:                   types.FlexibleTimeWindowModeOff,
				MaximumWindowInMinutes: aws.Int32(1441),
			},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return fmt.Errorf("OffModeWindowOutOfRange: %w", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_LogsTargetRejected", func() error {
		schedName := tc.uniqueName("LogsTarget")
		defer tc.cleanupSchedule(schedName)
		_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
			// CloudWatch Logs is not an EventBridge Scheduler
			// templated target in AWS.
			Arn:     aws.String(fmt.Sprintf("arn:aws:logs:%s:%s:log-group:/scheduler/delivery", tc.region, tc.accountID)),
			RoleArn: aws.String(rARN),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EventBridgeParameters_Source_InvalidPattern", func() error {
		schedName := tc.uniqueName("BadSource")
		defer tc.cleanupSchedule(schedName)
		_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
			Arn:     aws.String(tc.lambdaARN()),
			RoleArn: aws.String(rARN),
			EventBridgeParameters: &types.EventBridgeParameters{
				Source:     aws.String("!@#invalid"),
				DetailType: aws.String("MyDetailType"),
			},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	// SqsParameters without a MessageGroupId is model-valid input: the
	// member carries no @required in the model, and the FIFO fallback to
	// the schedule name is applied at delivery.
	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_SqsParametersEmptyAccepted", func() error {
		schedName := tc.uniqueName("EmptySqs")
		defer tc.cleanupSchedule(schedName)
		_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
			Arn:           aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:sqs-parameters-empty-pin", tc.region, tc.accountID)),
			RoleArn:       aws.String(rARN),
			SqsParameters: &types.SqsParameters{},
		})
		if err != nil {
			return err
		}
		out, err := tc.getSchedule(schedName)
		if err != nil {
			return err
		}
		if out.Target == nil || out.Target.SqsParameters == nil {
			return fmt.Errorf("GetSchedule did not return the SqsParameters of the created schedule")
		}
		// The unset MessageGroupId member is omitted from the read-back
		// rather than emitted as an empty string.
		if out.Target.SqsParameters.MessageGroupId != nil {
			return fmt.Errorf("unset MessageGroupId emitted as %q, want omission", *out.Target.SqsParameters.MessageGroupId)
		}

		// An explicitly present empty MessageGroupId is out of bounds on
		// the wire (@length min 1) — distinct from the omitted member
		// above, which stays accepted.
		presentEmpty := tc.uniqueName("EmptySqs2")
		defer tc.cleanupSchedule(presentEmpty)
		_, err = tc.createSchedule(presentEmpty, "rate(30 minutes)", &types.Target{
			Arn:     aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:sqs-parameters-present-empty-pin", tc.region, tc.accountID)),
			RoleArn: aws.String(rARN),
			SqsParameters: &types.SqsParameters{
				MessageGroupId: aws.String(""),
			},
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	// The model's Target.Input contract for the JSON-bound templated
	// families: "If you are configuring a templated Lambda, AWS Step
	// Functions, or Amazon EventBridge target, the input must be a
	// well-formed JSON" — any JSON value qualifies, and EventBridgeParameters
	// is optional on the wire. The object form is a delivery-time concern of
	// the events family, not a creation-time restriction.
	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EventBridgeTarget_Validation", func() error {
		busArn := fmt.Sprintf("arn:aws:events:%s:%s:event-bus/default", tc.region, tc.accountID)

		// Model-legal requests are created: a missing EventBridgeParameters
		// and non-object JSON Input values alike (delivery, not creation,
		// decides their fate).
		accepted := []struct {
			name string
			in   string
			eb   *types.EventBridgeParameters
		}{
			{"MissingEventBridgeParameters", `{"k":"v"}`, nil},
			{"ArrayInput", `[1,2]`, &types.EventBridgeParameters{
				Source: aws.String("vorpal.test"), DetailType: aws.String("validation-pin"),
			}},
		}
		for _, row := range accepted {
			schedName := tc.uniqueName("EbValidation")
			defer tc.cleanupSchedule(schedName)
			_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
				Arn:                   aws.String(busArn),
				RoleArn:               aws.String(rARN),
				Input:                 aws.String(row.in),
				EventBridgeParameters: row.eb,
			})
			if err != nil {
				return fmt.Errorf("%s: created-schedule rejected: %w", row.name, err)
			}
		}

		// A non-JSON Input on a JSON-bound templated family (events and
		// lambda) is rejected at creation.
		rejected := []struct {
			name string
			arn  string
			in   string
		}{
			{"MalformedEventsInput", busArn, `{k:v}`},
			{"NonJSONLambdaInput", tc.lambdaARN(), "hello"},
		}
		for _, row := range rejected {
			schedName := tc.uniqueName("EbValidation")
			defer tc.cleanupSchedule(schedName)
			_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
				Arn:     aws.String(row.arn),
				RoleArn: aws.String(rARN),
				Input:   aws.String(row.in),
			})
			if err := AssertErrorContains(err, "ValidationException"); err != nil {
				return fmt.Errorf("%s: %w", row.name, err)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_InvalidGroupName", func() error {
		schedName := tc.uniqueName("BadGroup")
		defer tc.cleanupSchedule(schedName)
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			GroupName:          aws.String("invalid/group name"),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	// The ClientToken pattern is enforced on both mutating operations.
	results = append(results, tc.runner.RunTest("scheduler", "Schedule_InvalidClientTokenRejected", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			{"CreateSchedule", func() error {
				schedName := tc.uniqueName("BadToken")
				defer tc.cleanupSchedule(schedName)
				_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
					Name:               aws.String(schedName),
					ScheduleExpression: aws.String("rate(30 minutes)"),
					ClientToken:        aws.String("bad token!"),
					Target: &types.Target{
						Arn:     aws.String(tc.lambdaARN()),
						RoleArn: aws.String(rARN),
					},
					FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
				})
				return AssertErrorContains(err, "ValidationException")
			}},
			{"UpdateSchedule", func() error {
				schedName := tc.uniqueName("UpdBadToken")
				_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
					Arn:     aws.String(tc.lambdaARN()),
					RoleArn: aws.String(rARN),
				})
				if err != nil {
					return fmt.Errorf("create: %v", err)
				}
				defer tc.cleanupSchedule(schedName)

				_, err = tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
					Name:               aws.String(schedName),
					ScheduleExpression: aws.String("rate(30 minutes)"),
					ClientToken:        aws.String("bad token!"),
					Target: &types.Target{
						Arn:     aws.String(tc.lambdaARN()),
						RoleArn: aws.String(rARN),
					},
					FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
				})
				return AssertErrorContains(err, "ValidationException")
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %w", row.name, err)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_TargetRoleArnInvalidRejected", func() error {
		invalidRoleArns := []string{
			// The Target shape requires RoleArn to reference an IAM role; a
			// queue ARN must be rejected as invalid input.
			fmt.Sprintf("arn:aws:sqs:%s:%s:not-a-role", tc.region, tc.accountID),
			// The RoleArn pattern also demands a 12-digit account and a
			// non-empty role path/name after role/.
			"arn:aws:iam::123456789012:role/",
			"arn:aws:iam::abc:role/x",
		}
		for _, roleArn := range invalidRoleArns {
			schedName := tc.uniqueName("BadRoleArn")
			defer tc.cleanupSchedule(schedName)
			_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(roleArn),
			})
			if err := AssertErrorContains(err, "ValidationException"); err != nil {
				return fmt.Errorf("RoleArn %q: %v", roleArn, err)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_UntrustedRoleRejected", func() error {
		// A role that exists but whose trust policy names a principal other
		// than scheduler.amazonaws.com is refused at creation, and the
		// refusal carries the modelled ValidationException identity — the
		// scheduler model defines no Lambda-style
		// InvalidParameterValueException or InvalidArn error shape.
		roleName := fmt.Sprintf("SchedUntrustedRole-%d", time.Now().UnixNano())
		lambdaOnlyTrust := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
		IAMCreateRole(tc.iamClient, roleName, lambdaOnlyTrust)
		defer IAMDeleteRole(tc.iamClient, roleName)
		untrustedARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.accountID, roleName)

		schedName := tc.uniqueName("UntrustedRole")
		defer tc.cleanupSchedule(schedName)
		_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
			Arn:     aws.String(tc.lambdaARN()),
			RoleArn: aws.String(untrustedARN),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		if strings.Contains(err.Error(), "InvalidParameterValueException") || strings.Contains(err.Error(), "InvalidArn") {
			return fmt.Errorf("role refusal surfaced an unmodelled error identity: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_TargetArnLengthRejected", func() error {
		// Target.Arn and RoleArn carry a 1600-character maximum in the
		// model; both over-length variants must be rejected.
		sqsPrefix := fmt.Sprintf("arn:aws:sqs:%s:%s:queue/", tc.region, tc.accountID)
		tooLongArn := sqsPrefix + strings.Repeat("q", 1601-len(sqsPrefix))
		rolePrefix := "arn:aws:iam::123456789012:role/"
		tooLongRole := rolePrefix + strings.Repeat("r", 1601-len(rolePrefix))

		schedName := tc.uniqueName("LongArn")
		defer tc.cleanupSchedule(schedName)
		_, err := tc.createSchedule(schedName, "rate(30 minutes)", &types.Target{
			Arn:     aws.String(tooLongArn),
			RoleArn: aws.String(rARN),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return fmt.Errorf("over-length Target.Arn: %v", err)
		}

		schedName2 := tc.uniqueName("LongRoleArn")
		defer tc.cleanupSchedule(schedName2)
		_, err = tc.createSchedule(schedName2, "rate(30 minutes)", &types.Target{
			Arn:     aws.String(tc.lambdaARN()),
			RoleArn: aws.String(tooLongRole),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return fmt.Errorf("over-length RoleArn: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "DeleteSchedule_ClientTokenNotSharedWithUpdate", func() error {
		schedName := tc.uniqueName("TokUpdDel")
		token := tc.uniqueName("tok")
		defer tc.cleanupSchedule(schedName)

		_, err := tc.createSchedule(schedName, "rate(30 minutes)", tc.defaultTarget(rARN))

		_, err = tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(45 minutes)"),
			ClientToken:        aws.String(token),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		// The same token on the delete belongs to a different request: the
		// deletion must be applied, not replayed as the update's outcome.
		if _, err := tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{
			Name:        aws.String(schedName),
			ClientToken: aws.String(token),
		}); err != nil {
			return fmt.Errorf("delete with reused token: %v", err)
		}

		_, err = tc.getSchedule(schedName)
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("schedule must be deleted after the token-reused delete: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "UpdateSchedule_ClientTokenNotSharedWithCreate", func() error {
		nameA := tc.uniqueName("TokCrateA")
		nameB := tc.uniqueName("TokCrateB")
		token := tc.uniqueName("tok")
		defer tc.cleanupSchedule(nameA)
		defer tc.cleanupSchedule(nameB)

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(nameA),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			ClientToken:        aws.String(token),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("create A: %v", err)
		}

		_, err = tc.createSchedule(nameB, "rate(30 minutes)", tc.defaultTarget(rARN))
		if err != nil {
			return fmt.Errorf("create B: %v", err)
		}

		out, err := tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(nameB),
			ScheduleExpression: aws.String("rate(45 minutes)"),
			ClientToken:        aws.String(token),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("update B: %v", err)
		}
		if aws.ToString(out.ScheduleArn) != tc.scheduleARN(nameB) {
			return fmt.Errorf("update must report B's ARN %q, got %q", tc.scheduleARN(nameB), aws.ToString(out.ScheduleArn))
		}

		got, err := tc.getSchedule(nameB)
		if err != nil {
			return fmt.Errorf("get B: %v", err)
		}
		if aws.ToString(got.ScheduleExpression) != "rate(45 minutes)" {
			return fmt.Errorf("B must carry the updated expression, got %q", aws.ToString(got.ScheduleExpression))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ScheduleClientToken_ReplayReusesOutcome", func() error {
		schedName := tc.uniqueName("TokReplay")
		createToken := tc.uniqueName("tok")
		updateToken := tc.uniqueName("tok")
		deleteToken := tc.uniqueName("tok")

		first, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			ClientToken:        aws.String(createToken),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.cleanupSchedule(schedName)

		// A replayed create reports the first application's ARN instead of a
		// name conflict.
		second, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			ClientToken:        aws.String(createToken),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("create replay: %v", err)
		}
		if aws.ToString(second.ScheduleArn) != aws.ToString(first.ScheduleArn) {
			return fmt.Errorf("create replay must return the first ARN %q, got %q",
				aws.ToString(first.ScheduleArn), aws.ToString(second.ScheduleArn))
		}

		updFirst, err := tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(45 minutes)"),
			ClientToken:        aws.String(updateToken),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		updReplay, err := tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(45 minutes)"),
			ClientToken:        aws.String(updateToken),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("update replay: %v", err)
		}
		if aws.ToString(updReplay.ScheduleArn) != aws.ToString(updFirst.ScheduleArn) {
			return fmt.Errorf("update replay must return the first ARN %q, got %q",
				aws.ToString(updFirst.ScheduleArn), aws.ToString(updReplay.ScheduleArn))
		}

		if _, err := tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{
			Name:        aws.String(schedName),
			ClientToken: aws.String(deleteToken),
		}); err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		// A replayed delete reports the first deletion's outcome: success,
		// not a second ResourceNotFoundException.
		if _, err := tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{
			Name:        aws.String(schedName),
			ClientToken: aws.String(deleteToken),
		}); err != nil {
			return fmt.Errorf("delete replay: %v", err)
		}
		return nil
	}))

	return results
}
