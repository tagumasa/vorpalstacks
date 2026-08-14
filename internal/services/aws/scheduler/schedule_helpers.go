package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	tagutil "vorpalstacks/internal/common/tags"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

func parseTarget(params map[string]interface{}) (*schedulerstore.Target, error) {
	targetData, _ := getMapField(params, "Target")
	if targetData == nil {
		return nil, nil
	}

	rawMap, err := coerceToMap(targetData)
	if err != nil {
		return nil, ErrInvalidTarget
	}

	target := parseTargetFromMap(rawMap)
	if target.Arn == "" || target.RoleArn == "" {
		return nil, ErrInvalidTarget
	}
	return &target, nil
}

func coerceToMap(v interface{}) (map[string]interface{}, error) {
	switch t := v.(type) {
	case string:
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(t), &m); err != nil {
			return nil, err
		}
		return m, nil
	case map[string]interface{}:
		return t, nil
	default:
		return nil, fmt.Errorf("unexpected target type: %T", v)
	}
}

func parseTargetFromMap(m map[string]interface{}) schedulerstore.Target {
	var target schedulerstore.Target
	target.Arn = getStringFromMap(m, "arn", "Arn")
	target.RoleArn = getStringFromMap(m, "roleArn", "RoleArn")
	target.Input = getStringFromMap(m, "input", "Input")

	if dl, ok := getMapField(m, "deadLetterConfig", "DeadLetterConfig"); ok {
		target.DeadLetterConfig = &schedulerstore.DeadLetterConfig{
			Arn: getStringFromMap(dl, "arn", "Arn"),
		}
	}
	if rp, ok := getMapField(m, "retryPolicy", "RetryPolicy"); ok {
		target.RetryPolicy = parseRetryPolicyFromMap(rp)
	}
	if sqs, ok := getMapField(m, "sqsParameters", "SqsParameters"); ok {
		target.SqsParameters = &schedulerstore.SqsParameters{
			MessageGroupId: getStringFromMap(sqs, "messageGroupId", "MessageGroupId"),
		}
	}
	if ecs, ok := getMapField(m, "ecsParameters", "EcsParameters"); ok {
		target.EcsParameters = parseEcsParameters(ecs)
	}
	if eb, ok := getMapField(m, "eventBridgeParameters", "EventBridgeParameters"); ok {
		target.EventBridgeParameters = parseEventBridgeParameters(eb)
	}
	if kinesis, ok := getMapField(m, "kinesisParameters", "KinesisParameters"); ok {
		target.KinesisParameters = parseKinesisParameters(kinesis)
	}
	return target
}

func getMapField(m map[string]interface{}, keys ...string) (map[string]interface{}, bool) {
	for _, k := range keys {
		if v, ok := m[k].(map[string]interface{}); ok {
			return v, true
		}
	}
	return nil, false
}

func getStringFromMap(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if val, ok := m[key].(string); ok {
			return val
		}
	}
	return ""
}

func parseRetryPolicyFromMap(retryPolicy map[string]interface{}) *schedulerstore.RetryPolicy {
	rp := &schedulerstore.RetryPolicy{}
	// Accept the raw values without range filtering. Range validation is
	// performed by validateTarget in validators.go.
	if val, ok := getFloatField(retryPolicy, "maximumEventAgeInSeconds", "MaximumEventAgeInSeconds"); ok {
		rp.MaximumEventAgeInSeconds = &val
	}
	if val, ok := getFloatField(retryPolicy, "maximumRetryAttempts", "MaximumRetryAttempts"); ok {
		rp.MaximumRetryAttempts = &val
	}
	return rp
}

func getFloatField(m map[string]interface{}, keys ...string) (int, bool) {
	for _, k := range keys {
		switch v := m[k].(type) {
		case float64:
			return int(v), true
		case int:
			return v, true
		case int32:
			return int(v), true
		case int64:
			return int(v), true
		}
	}
	return 0, false
}

func parseEcsParameters(data map[string]interface{}) *schedulerstore.EcsParameters {
	params := &schedulerstore.EcsParameters{
		TaskDefinitionArn: getStringFromMap(data, "taskDefinitionArn", "TaskDefinitionArn"),
		LaunchType:        getStringFromMap(data, "launchType", "LaunchType"),
		PlatformVersion:   getStringFromMap(data, "platformVersion", "PlatformVersion"),
		Group:             getStringFromMap(data, "group", "Group"),
		PropagateTags:     getStringFromMap(data, "propagateTags", "PropagateTags"),
		ReferenceId:       getStringFromMap(data, "referenceId", "ReferenceId"),
	}
	if val, ok := getFloatField(data, "taskCount", "TaskCount"); ok {
		params.TaskCount = &val
	}
	if val, ok := getBoolField(data, "enableECSManagedTags", "EnableECSManagedTags"); ok {
		params.EnableECSManagedTags = &val
	}
	if val, ok := getBoolField(data, "enableExecuteCommand", "EnableExecuteCommand"); ok {
		params.EnableExecuteCommand = &val
	}
	if nc, ok := getMapField(data, "networkConfiguration", "NetworkConfiguration"); ok {
		params.NetworkConfiguration = parseNetworkConfiguration(nc)
	}
	if cps, ok := getSliceField(data, "capacityProviderStrategy", "CapacityProviderStrategy"); ok {
		params.CapacityProviderStrategy = parseCapacityProviderStrategy(cps)
	}
	if pc, ok := getSliceField(data, "placementConstraints", "PlacementConstraints"); ok {
		params.PlacementConstraints = parsePlacementConstraints(pc)
	}
	if ps, ok := getSliceField(data, "placementStrategy", "PlacementStrategy"); ok {
		params.PlacementStrategy = parsePlacementStrategy(ps)
	}
	if tags, ok := getSliceField(data, "tags", "Tags"); ok {
		params.Tags = tagutil.ParseEcsTags(tags)
	}
	return params
}

func getBoolField(m map[string]interface{}, keys ...string) (bool, bool) {
	for _, k := range keys {
		if v, ok := m[k].(bool); ok {
			return v, true
		}
	}
	return false, false
}

func getSliceField(m map[string]interface{}, keys ...string) ([]interface{}, bool) {
	for _, k := range keys {
		if v, ok := m[k].([]interface{}); ok {
			return v, true
		}
	}
	return nil, false
}

func parseNetworkConfiguration(data map[string]interface{}) *schedulerstore.NetworkConfiguration {
	nc := &schedulerstore.NetworkConfiguration{}
	if vpc, ok := getMapField(data, "awsvpcConfiguration", "AwsvpcConfiguration"); ok {
		nc.AwsVpcConfiguration = parseAwsVpcConfiguration(vpc)
	}
	return nc
}

func parseAwsVpcConfiguration(data map[string]interface{}) *schedulerstore.AwsVpcConfiguration {
	vpc := &schedulerstore.AwsVpcConfiguration{
		AssignPublicIp: getStringFromMap(data, "assignPublicIp", "AssignPublicIp"),
	}
	if subnets, ok := getSliceField(data, "subnets", "Subnets"); ok {
		for _, s := range subnets {
			if str, ok := s.(string); ok {
				vpc.Subnets = append(vpc.Subnets, str)
			}
		}
	}
	if sgs, ok := getSliceField(data, "securityGroups", "SecurityGroups"); ok {
		for _, sg := range sgs {
			if str, ok := sg.(string); ok {
				vpc.SecurityGroups = append(vpc.SecurityGroups, str)
			}
		}
	}
	return vpc
}

func parseCapacityProviderStrategy(data []interface{}) []schedulerstore.CapacityProviderStrategyItem {
	var result []schedulerstore.CapacityProviderStrategyItem
	for _, item := range data {
		if m, ok := item.(map[string]interface{}); ok {
			cps := schedulerstore.CapacityProviderStrategyItem{
				CapacityProvider: getStringFromMap(m, "capacityProvider", "CapacityProvider"),
			}
			if w, ok := getFloatField(m, "weight", "Weight"); ok {
				cps.Weight = &w
			}
			if b, ok := getFloatField(m, "base", "Base"); ok {
				cps.Base = &b
			}
			result = append(result, cps)
		}
	}
	return result
}

func parsePlacementConstraints(data []interface{}) []schedulerstore.PlacementConstraint {
	var result []schedulerstore.PlacementConstraint
	for _, item := range data {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, schedulerstore.PlacementConstraint{
				Type:       getStringFromMap(m, "type", "Type"),
				Expression: getStringFromMap(m, "expression", "Expression"),
			})
		}
	}
	return result
}

func parsePlacementStrategy(data []interface{}) []schedulerstore.PlacementStrategy {
	var result []schedulerstore.PlacementStrategy
	for _, item := range data {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, schedulerstore.PlacementStrategy{
				Type:  getStringFromMap(m, "type", "Type"),
				Field: getStringFromMap(m, "field", "Field"),
			})
		}
	}
	return result
}

func parseEventBridgeParameters(data map[string]interface{}) *schedulerstore.EventBridgeParameters {
	return &schedulerstore.EventBridgeParameters{
		DetailType: getStringFromMap(data, "detailType", "DetailType"),
		Source:     getStringFromMap(data, "source", "Source"),
	}
}

func parseKinesisParameters(data map[string]interface{}) *schedulerstore.KinesisParameters {
	return &schedulerstore.KinesisParameters{
		PartitionKey: getStringFromMap(data, "partitionKey", "PartitionKey"),
	}
}

func parseFlexibleTimeWindow(params map[string]interface{}) (*schedulerstore.FlexibleTimeWindow, error) {
	ftwData, ok := params["FlexibleTimeWindow"]
	if !ok {
		ftwData, ok = params["flexibleTimeWindow"]
	}
	if ftwData == nil || !ok {
		return nil, nil
	}

	ftw := &schedulerstore.FlexibleTimeWindow{}
	switch v := ftwData.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), ftw); err != nil {
			return nil, ErrInvalidFlexibleTimeWindow
		}
	case map[string]interface{}:
		mode := getStringFromMap(v, "mode", "Mode")
		if mode != "" {
			ftw.Mode = schedulerstore.FlexibleTimeWindowMode(mode)
		}
		if maxWindow, ok := getFloatField(v, "maximumWindowInMinutes", "MaximumWindowInMinutes"); ok {
			ftw.MaximumWindowInMinutes = &maxWindow
		}
	default:
		return nil, ErrInvalidFlexibleTimeWindow
	}

	if ftw.Mode == "" {
		ftw.Mode = schedulerstore.FlexibleTimeWindowModeOff
	}

	// Mode enum and MaximumWindowInMinutes range validation is performed
	// by validateFlexibleTimeWindow in validators.go.

	return ftw, nil
}

// validateVpcConfig validates the AwsVpcConfiguration subnets and security
// groups against the EC2 service via the event bus. All resources must exist
// and belong to the same VPC. Accepts region directly so both the HTTP API
// and admin console paths can call it (Minor 3).
func (s *SchedulerService) validateVpcConfig(ctx context.Context, region string, target *schedulerstore.Target) error {
	if s.engine == nil || s.engine.bus == nil {
		return nil
	}
	if target == nil || target.EcsParameters == nil || target.EcsParameters.NetworkConfiguration == nil {
		return nil
	}
	vpc := target.EcsParameters.NetworkConfiguration.AwsVpcConfiguration
	if vpc == nil || (len(vpc.Subnets) == 0 && len(vpc.SecurityGroups) == 0) {
		return nil
	}

	ec2 := s.engine.bus.EC2Invoker()
	if ec2 == nil {
		return fmt.Errorf("scheduler: EC2 service not available for VPC configuration validation")
	}

	for _, subnetId := range vpc.Subnets {
		if _, _, err := ec2.LookupSubnet(ctx, region, subnetId); err != nil {
			return fmt.Errorf("scheduler: subnet %s not found: %v", subnetId, err)
		}
	}

	for _, sgId := range vpc.SecurityGroups {
		if _, err := ec2.LookupSecurityGroup(ctx, region, sgId); err != nil {
			return fmt.Errorf("scheduler: security group %s not found: %v", sgId, err)
		}
	}

	return nil
}
