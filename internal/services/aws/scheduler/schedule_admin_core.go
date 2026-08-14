package scheduler

import (
	"context"

	"vorpalstacks/internal/common/iam"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"

	pb "vorpalstacks/internal/pb/aws/scheduler"
)

// AdminCreateScheduleInput carries the proto-derived parameters that the
// admin console sends when creating a schedule. Complex nested types
// (Target, FlexibleTimeWindow) are passed as proto messages and converted
// to store types inside createScheduleFromAdmin, so that admin handler
// code never constructs store structs directly (store-import prohibition).
type AdminCreateScheduleInput struct {
	Name                       string
	GroupName                  string
	ScheduleExpression         string
	ScheduleExpressionTimezone string
	Description                string
	State                      string
	KmsKeyArn                  string
	StartDate                  string
	EndDate                    string
	ActionAfterCompletion      string
	Target                     *pb.Target
	FlexibleTimeWindow         *pb.FlexibleTimeWindow
	ClientToken                string
	Region                     string
	IAMValidator               *iam.IAMValidator
}

// createScheduleFromAdmin is the Core entry point for the admin console.
// It converts the proto-derived input to the transport-agnostic
// ScheduleSpec and delegates to createScheduleCore.
func (s *SchedulerService) createScheduleFromAdmin(ctx context.Context, store *schedulerstore.SchedulerStore, in AdminCreateScheduleInput) (*CreateScheduleResult, error) {
	spec := &ScheduleSpec{
		Name:                       in.Name,
		GroupName:                  in.GroupName,
		ScheduleExpression:         in.ScheduleExpression,
		ScheduleExpressionTimezone: in.ScheduleExpressionTimezone,
		Description:                in.Description,
		State:                      in.State,
		KmsKeyArn:                  in.KmsKeyArn,
		StartDate:                  in.StartDate,
		EndDate:                    in.EndDate,
		ActionAfterCompletion:      in.ActionAfterCompletion,
		Target:                     protoTargetToStore(in.Target),
		FlexibleTimeWindow:         protoFTWToStore(in.FlexibleTimeWindow),
	}

	return s.createScheduleCore(ctx, store, &CreateScheduleInput{
		Spec:         spec,
		ClientToken:  in.ClientToken,
		Region:       in.Region,
		IAMValidator: in.IAMValidator,
	})
}

// protoFTWToStore converts a protobuf FlexibleTimeWindow to the store type.
func protoFTWToStore(pbFTW *pb.FlexibleTimeWindow) *schedulerstore.FlexibleTimeWindow {
	if pbFTW == nil {
		return &schedulerstore.FlexibleTimeWindow{Mode: schedulerstore.FlexibleTimeWindowModeOff}
	}
	ftw := &schedulerstore.FlexibleTimeWindow{
		Mode: schedulerstore.FlexibleTimeWindowMode(pbFTW.Mode),
	}
	if pbFTW.Maximumwindowinminutes != nil {
		v := int(*pbFTW.Maximumwindowinminutes)
		ftw.MaximumWindowInMinutes = &v
	}
	return ftw
}

// protoTargetToStore converts a protobuf Target message to the store Target
// type, mapping all sub-fields faithfully.
func protoTargetToStore(pbTarget *pb.Target) *schedulerstore.Target {
	if pbTarget == nil {
		return nil
	}
	t := &schedulerstore.Target{
		Arn:     pbTarget.Arn,
		Input:   pbTarget.Input,
		RoleArn: pbTarget.Rolearn,
	}

	if pbTarget.Deadletterconfig != nil {
		t.DeadLetterConfig = &schedulerstore.DeadLetterConfig{
			Arn: pbTarget.Deadletterconfig.Arn,
		}
	}

	if pbTarget.Retrypolicy != nil {
		rp := &schedulerstore.RetryPolicy{}
		if pbTarget.Retrypolicy.Maximumeventageinseconds != nil {
			v := int(*pbTarget.Retrypolicy.Maximumeventageinseconds)
			rp.MaximumEventAgeInSeconds = &v
		}
		if pbTarget.Retrypolicy.Maximumretryattempts != nil {
			v := int(*pbTarget.Retrypolicy.Maximumretryattempts)
			rp.MaximumRetryAttempts = &v
		}
		t.RetryPolicy = rp
	}

	if pbTarget.Sqsparameters != nil {
		t.SqsParameters = &schedulerstore.SqsParameters{
			MessageGroupId: pbTarget.Sqsparameters.Messagegroupid,
		}
	}

	if pbTarget.Ecsparameters != nil {
		t.EcsParameters = protoEcsParametersToStore(pbTarget.Ecsparameters)
	}

	if pbTarget.Eventbridgeparameters != nil {
		t.EventBridgeParameters = &schedulerstore.EventBridgeParameters{
			DetailType: pbTarget.Eventbridgeparameters.Detailtype,
			Source:     pbTarget.Eventbridgeparameters.Source,
		}
	}

	if pbTarget.Kinesisparameters != nil {
		t.KinesisParameters = &schedulerstore.KinesisParameters{
			PartitionKey: pbTarget.Kinesisparameters.Partitionkey,
		}
	}

	return t
}

// protoEcsParametersToStore converts a protobuf EcsParameters to the store
// type, including all nested sub-structures.
func protoEcsParametersToStore(pbEcs *pb.EcsParameters) *schedulerstore.EcsParameters {
	ecs := &schedulerstore.EcsParameters{
		TaskDefinitionArn: pbEcs.Taskdefinitionarn,
		LaunchType:        pbEcs.Launchtype,
		PlatformVersion:   pbEcs.Platformversion,
		Group:             pbEcs.Group,
		PropagateTags:     pbEcs.Propagatetags,
		ReferenceId:       pbEcs.Referenceid,
	}

	if pbEcs.Taskcount != nil {
		v := int(*pbEcs.Taskcount)
		ecs.TaskCount = &v
	}
	if pbEcs.Enableecsmanagedtags != nil {
		b := *pbEcs.Enableecsmanagedtags
		ecs.EnableECSManagedTags = &b
	}
	if pbEcs.Enableexecutecommand != nil {
		b := *pbEcs.Enableexecutecommand
		ecs.EnableExecuteCommand = &b
	}

	if pbEcs.Networkconfiguration != nil {
		ecs.NetworkConfiguration = protoNetworkConfigToStore(pbEcs.Networkconfiguration)
	}

	for _, cps := range pbEcs.Capacityproviderstrategy {
		item := schedulerstore.CapacityProviderStrategyItem{
			CapacityProvider: cps.Capacityprovider,
		}
		if cps.Weight != nil {
			w := int(*cps.Weight)
			item.Weight = &w
		}
		if cps.Base != nil {
			b := int(*cps.Base)
			item.Base = &b
		}
		ecs.CapacityProviderStrategy = append(ecs.CapacityProviderStrategy, item)
	}

	for _, pc := range pbEcs.Placementconstraints {
		ecs.PlacementConstraints = append(ecs.PlacementConstraints, schedulerstore.PlacementConstraint{
			Type:       pc.Type,
			Expression: pc.Expression,
		})
	}

	for _, ps := range pbEcs.Placementstrategy {
		ecs.PlacementStrategy = append(ecs.PlacementStrategy, schedulerstore.PlacementStrategy{
			Type:  ps.Type,
			Field: ps.Field,
		})
	}

	return ecs
}

// protoNetworkConfigToStore converts a protobuf NetworkConfiguration to
// the store type.
func protoNetworkConfigToStore(pbNet *pb.NetworkConfiguration) *schedulerstore.NetworkConfiguration {
	if pbNet.Awsvpcconfiguration == nil {
		return &schedulerstore.NetworkConfiguration{}
	}
	vpc := pbNet.Awsvpcconfiguration
	nc := &schedulerstore.NetworkConfiguration{
		AwsVpcConfiguration: &schedulerstore.AwsVpcConfiguration{
			Subnets:        vpc.Subnets,
			SecurityGroups: vpc.Securitygroups,
		},
	}
	if vpc.Assignpublicip != "" {
		nc.AwsVpcConfiguration.AssignPublicIp = vpc.Assignpublicip
	}
	return nc
}
