package testutil

import (
	"fmt"

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

	results = append(results, tc.runner.RunTest("scheduler", "GetSchedule_NonExistent", func() error {
		_, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{
			Name: aws.String("nonexistent-schedule-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "DeleteSchedule_NonExistent", func() error {
		_, err := tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String("nonexistent-schedule-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_DuplicateName", func() error {
		dupName := tc.uniqueName("DupSched")
		dupRN, dupRARN := tc.createIAMRole()
		defer tc.deleteIAMRole(dupRN)

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(dupName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             tc.defaultTarget(dupRARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(dupName)})

		_, err = tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(dupName),
			ScheduleExpression: aws.String("rate(60 minutes)"),
			Target:             tc.defaultTarget(dupRARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ConflictException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "UpdateSchedule_NonExistent", func() error {
		_, err := tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String("nonexistent-schedule-xyz"),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_InvalidExpression", func() error {
		invName := tc.uniqueName("InvExprSched")
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(invName),
			ScheduleExpression: aws.String("not-a-valid-expression"),
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

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_UnsupportedTarget_Rejected", func() error {
		schedName := tc.uniqueName("UnsupportedTarget")
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:codebuild:%s:%s:project/FakeProject", tc.region, tc.accountID)),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_DeadLetterConfig_SnsRejected", func() error {
		schedName := tc.uniqueName("SnsDLQ")
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
				DeadLetterConfig: &types.DeadLetterConfig{
					Arn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:MyTopic", tc.region, tc.accountID)),
				},
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EcsParameters_TaskCountOutOfRange", func() error {
		schedName := tc.uniqueName("EcsBadCount")
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		taskCount := int32(100)
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
				EcsParameters: &types.EcsParameters{
					TaskDefinitionArn: aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/FakeTask:1", tc.region, tc.accountID)),
					TaskCount:         &taskCount,
				},
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EcsParameters_SubnetNotFound", func() error {
		schedName := tc.uniqueName("BadSubnet")
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		taskCount := int32(1)
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/ValidationCluster", tc.region, tc.accountID)),
				RoleArn: aws.String(rARN),
				EcsParameters: &types.EcsParameters{
					TaskDefinitionArn: aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/FakeTask:1", tc.region, tc.accountID)),
					TaskCount:         &taskCount,
					NetworkConfiguration: &types.NetworkConfiguration{
						AwsvpcConfiguration: &types.AwsVpcConfiguration{
							Subnets: []string{"subnet-nonexistent-xyz"},
						},
					},
				},
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EcsParameters_EmptySecurityGroupsRejected", func() error {
		schedName := tc.uniqueName("EmptySG")
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		taskCount := int32(1)
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/ValidationCluster", tc.region, tc.accountID)),
				RoleArn: aws.String(rARN),
				EcsParameters: &types.EcsParameters{
					TaskDefinitionArn: aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/FakeTask:1", tc.region, tc.accountID)),
					TaskCount:         &taskCount,
					NetworkConfiguration: &types.NetworkConfiguration{
						AwsvpcConfiguration: &types.AwsVpcConfiguration{
							Subnets:        []string{"subnet-0123456789abcdef0"},
							SecurityGroups: []string{},
						},
					},
				},
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EcsParameters_CapacityProviderWeightOutOfRange", func() error {
		schedName := tc.uniqueName("EcsBadWeight")
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		weight := int32(5000)
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/ValidationCluster", tc.region, tc.accountID)),
				RoleArn: aws.String(rARN),
				EcsParameters: &types.EcsParameters{
					TaskDefinitionArn: aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/FakeTask:1", tc.region, tc.accountID)),
					CapacityProviderStrategy: []types.CapacityProviderStrategyItem{
						{CapacityProvider: aws.String("FARGATE"), Weight: weight},
					},
				},
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EcsParameters_PlacementStrategyInvalidType", func() error {
		schedName := tc.uniqueName("EcsBadStrategy")
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/ValidationCluster", tc.region, tc.accountID)),
				RoleArn: aws.String(rARN),
				EcsParameters: &types.EcsParameters{
					TaskDefinitionArn: aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/FakeTask:1", tc.region, tc.accountID)),
					PlacementStrategy: []types.PlacementStrategy{
						{Type: types.PlacementStrategyType("bogus")},
					},
				},
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EcsParameters_EmptySubnetsRejected", func() error {
		schedName := tc.uniqueName("EcsNoSubnets")
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/ValidationCluster", tc.region, tc.accountID)),
				RoleArn: aws.String(rARN),
				EcsParameters: &types.EcsParameters{
					TaskDefinitionArn: aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/FakeTask:1", tc.region, tc.accountID)),
					NetworkConfiguration: &types.NetworkConfiguration{
						AwsvpcConfiguration: &types.AwsVpcConfiguration{
							Subnets: []string{},
						},
					},
				},
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EcsParameters_NonEcsTarget_Rejected", func() error {
		schedName := tc.uniqueName("EcsOnLambda")
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
				EcsParameters: &types.EcsParameters{
					TaskDefinitionArn: aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/FakeTask:1", tc.region, tc.accountID)),
				},
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
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
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
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
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
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
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				// CloudWatch Logs is not an EventBridge Scheduler
				// templated target in AWS.
				Arn:     aws.String(fmt.Sprintf("arn:aws:logs:%s:%s:log-group:/scheduler/delivery", tc.region, tc.accountID)),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_EventBridgeParameters_Source_InvalidPattern", func() error {
		schedName := tc.uniqueName("BadSource")
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
				EventBridgeParameters: &types.EventBridgeParameters{
					Source:     aws.String("!@#invalid"),
					DetailType: aws.String("MyDetailType"),
				},
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
