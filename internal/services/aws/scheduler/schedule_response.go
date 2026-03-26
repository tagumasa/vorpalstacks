package scheduler

import (
	"time"

	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

func formatEpochSeconds(t time.Time) float64 {
	return float64(t.Unix()) + float64(t.Nanosecond())/1e9
}

func scheduleToResponse(schedule *schedulerstore.Schedule) map[string]interface{} {
	resp := map[string]interface{}{
		"Arn":                  schedule.ARN,
		"Name":                 schedule.Name,
		"GroupName":            schedule.GroupName,
		"ScheduleExpression":   schedule.ScheduleExpression,
		"State":                string(schedule.State),
		"CreationDate":         formatEpochSeconds(schedule.CreationDate),
		"LastModificationDate": formatEpochSeconds(schedule.LastModificationDate),
	}

	if schedule.Description != "" {
		resp["Description"] = schedule.Description
	}
	if schedule.ScheduleExpressionTimezone != "" {
		resp["ScheduleExpressionTimezone"] = schedule.ScheduleExpressionTimezone
	}
	if schedule.KmsKeyArn != "" {
		resp["KmsKeyArn"] = schedule.KmsKeyArn
	}
	if schedule.StartDate != nil {
		resp["StartDate"] = formatEpochSeconds(*schedule.StartDate)
	}
	if schedule.EndDate != nil {
		resp["EndDate"] = formatEpochSeconds(*schedule.EndDate)
	}
	if schedule.ActionAfterCompletion != "" {
		resp["ActionAfterCompletion"] = string(schedule.ActionAfterCompletion)
	}
	if schedule.Target != nil {
		resp["Target"] = targetToResponse(schedule.Target)
	}
	if schedule.FlexibleTimeWindow != nil {
		ftw := map[string]interface{}{
			"Mode": string(schedule.FlexibleTimeWindow.Mode),
		}
		if schedule.FlexibleTimeWindow.MaximumWindowInMinutes != nil {
			ftw["MaximumWindowInMinutes"] = *schedule.FlexibleTimeWindow.MaximumWindowInMinutes
		}
		resp["FlexibleTimeWindow"] = ftw
	}

	return resp
}

func targetToResponse(target *schedulerstore.Target) map[string]interface{} {
	resp := map[string]interface{}{
		"Arn":     target.Arn,
		"RoleArn": target.RoleArn,
	}
	if target.Input != "" {
		resp["Input"] = target.Input
	}
	if target.DeadLetterConfig != nil {
		resp["DeadLetterConfig"] = map[string]interface{}{
			"Arn": target.DeadLetterConfig.Arn,
		}
	}
	if target.RetryPolicy != nil {
		rp := map[string]interface{}{}
		if target.RetryPolicy.MaximumEventAgeInSeconds != nil {
			rp["MaximumEventAgeInSeconds"] = *target.RetryPolicy.MaximumEventAgeInSeconds
		}
		if target.RetryPolicy.MaximumRetryAttempts != nil {
			rp["MaximumRetryAttempts"] = *target.RetryPolicy.MaximumRetryAttempts
		}
		resp["RetryPolicy"] = rp
	}
	if target.SqsParameters != nil {
		resp["SqsParameters"] = map[string]interface{}{
			"MessageGroupId": target.SqsParameters.MessageGroupId,
		}
	}
	if target.EcsParameters != nil {
		ecs := map[string]interface{}{
			"TaskDefinitionArn": target.EcsParameters.TaskDefinitionArn,
		}
		if target.EcsParameters.LaunchType != "" {
			ecs["LaunchType"] = target.EcsParameters.LaunchType
		}
		if target.EcsParameters.PlatformVersion != "" {
			ecs["PlatformVersion"] = target.EcsParameters.PlatformVersion
		}
		if target.EcsParameters.Group != "" {
			ecs["Group"] = target.EcsParameters.Group
		}
		if target.EcsParameters.TaskCount != nil {
			ecs["TaskCount"] = *target.EcsParameters.TaskCount
		}
		if target.EcsParameters.PropagateTags != "" {
			ecs["PropagateTags"] = target.EcsParameters.PropagateTags
		}
		if target.EcsParameters.ReferenceId != "" {
			ecs["ReferenceId"] = target.EcsParameters.ReferenceId
		}
		if target.EcsParameters.EnableECSManagedTags != nil {
			ecs["EnableECSManagedTags"] = *target.EcsParameters.EnableECSManagedTags
		}
		if target.EcsParameters.EnableExecuteCommand != nil {
			ecs["EnableExecuteCommand"] = *target.EcsParameters.EnableExecuteCommand
		}
		if target.EcsParameters.NetworkConfiguration != nil {
			ecs["NetworkConfiguration"] = networkConfigurationToResponse(target.EcsParameters.NetworkConfiguration)
		}
		if len(target.EcsParameters.CapacityProviderStrategy) > 0 {
			ecs["CapacityProviderStrategy"] = capacityProviderStrategyToResponse(target.EcsParameters.CapacityProviderStrategy)
		}
		if len(target.EcsParameters.PlacementConstraints) > 0 {
			ecs["PlacementConstraints"] = placementConstraintsToResponse(target.EcsParameters.PlacementConstraints)
		}
		if len(target.EcsParameters.PlacementStrategy) > 0 {
			ecs["PlacementStrategy"] = placementStrategyToResponse(target.EcsParameters.PlacementStrategy)
		}
		if len(target.EcsParameters.Tags) > 0 {
			ecs["Tags"] = target.EcsParameters.Tags
		}
		resp["EcsParameters"] = ecs
	}
	if target.EventBridgeParameters != nil {
		resp["EventBridgeParameters"] = map[string]interface{}{
			"DetailType": target.EventBridgeParameters.DetailType,
			"Source":     target.EventBridgeParameters.Source,
		}
	}
	if target.KinesisParameters != nil {
		resp["KinesisParameters"] = map[string]interface{}{
			"PartitionKey": target.KinesisParameters.PartitionKey,
		}
	}
	return resp
}

func networkConfigurationToResponse(nc *schedulerstore.NetworkConfiguration) map[string]interface{} {
	if nc == nil || nc.AwsVpcConfiguration == nil {
		return nil
	}
	vpc := map[string]interface{}{}
	if len(nc.AwsVpcConfiguration.Subnets) > 0 {
		vpc["Subnets"] = nc.AwsVpcConfiguration.Subnets
	}
	if len(nc.AwsVpcConfiguration.SecurityGroups) > 0 {
		vpc["SecurityGroups"] = nc.AwsVpcConfiguration.SecurityGroups
	}
	if nc.AwsVpcConfiguration.AssignPublicIp != "" {
		vpc["AssignPublicIp"] = nc.AwsVpcConfiguration.AssignPublicIp
	}
	return map[string]interface{}{"AwsVpcConfiguration": vpc}
}

func capacityProviderStrategyToResponse(cps []schedulerstore.CapacityProviderStrategyItem) []map[string]interface{} {
	result := make([]map[string]interface{}, len(cps))
	for i, cp := range cps {
		item := map[string]interface{}{"CapacityProvider": cp.CapacityProvider}
		if cp.Weight != nil {
			item["Weight"] = *cp.Weight
		}
		if cp.Base != nil {
			item["Base"] = *cp.Base
		}
		result[i] = item
	}
	return result
}

func placementConstraintsToResponse(pc []schedulerstore.PlacementConstraint) []map[string]interface{} {
	result := make([]map[string]interface{}, len(pc))
	for i, c := range pc {
		item := map[string]interface{}{}
		if c.Type != "" {
			item["Type"] = c.Type
		}
		if c.Expression != "" {
			item["Expression"] = c.Expression
		}
		result[i] = item
	}
	return result
}

func placementStrategyToResponse(ps []schedulerstore.PlacementStrategy) []map[string]interface{} {
	result := make([]map[string]interface{}, len(ps))
	for i, s := range ps {
		item := map[string]interface{}{}
		if s.Type != "" {
			item["Type"] = s.Type
		}
		if s.Field != "" {
			item["Field"] = s.Field
		}
		result[i] = item
	}
	return result
}
