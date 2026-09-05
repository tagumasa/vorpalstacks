package testutil

import (
	"fmt"
	"strings"

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
		invName := tc.uniqueName("InvExprSched")
		_, err := tc.createSchedule(invName, "not-a-valid-expression", &types.Target{
			Arn:     aws.String(tc.lambdaARN()),
			RoleArn: aws.String(rARN),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
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
