package scheduler

import (
	"encoding/json"
	"regexp"

	tagutil "vorpalstacks/internal/common/tags"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

var (
	atExpressionRegex   = regexp.MustCompile(`^at\((\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})\)$`)
	rateExpressionRegex = regexp.MustCompile(`^rate\((\d+)\s+(minute|minutes|hour|hours|day|days)\)$`)
	cronExpressionRegex = regexp.MustCompile(`^cron\((.+)\)$`)
)

func isValidScheduleExpression(expr string) bool {
	return atExpressionRegex.MatchString(expr) ||
		rateExpressionRegex.MatchString(expr) ||
		cronExpressionRegex.MatchString(expr)
}

func parseTarget(params map[string]interface{}) (*schedulerstore.Target, error) {
	targetData, ok := params["Target"]
	if !ok {
		targetData, ok = params["target"]
	}
	if targetData == nil || !ok {
		return nil, nil
	}

	var target schedulerstore.Target
	switch v := targetData.(type) {
	case string:
		var rawMap map[string]interface{}
		if err := json.Unmarshal([]byte(v), &rawMap); err != nil {
			return nil, ErrInvalidTarget
		}
		target.Arn = getStringFromMap(rawMap, "arn", "Arn")
		target.RoleArn = getStringFromMap(rawMap, "roleArn", "RoleArn")
		target.Input = getStringFromMap(rawMap, "input", "Input")
		if dlConfig, ok := rawMap["deadLetterConfig"].(map[string]interface{}); ok {
			target.DeadLetterConfig = &schedulerstore.DeadLetterConfig{
				Arn: getStringFromMap(dlConfig, "arn", "Arn"),
			}
		} else if dlConfig, ok := rawMap["DeadLetterConfig"].(map[string]interface{}); ok {
			target.DeadLetterConfig = &schedulerstore.DeadLetterConfig{
				Arn: getStringFromMap(dlConfig, "arn", "Arn"),
			}
		}
		if retryPolicy, ok := rawMap["retryPolicy"].(map[string]interface{}); ok {
			target.RetryPolicy = parseRetryPolicyFromMap(retryPolicy)
		} else if retryPolicy, ok := rawMap["RetryPolicy"].(map[string]interface{}); ok {
			target.RetryPolicy = parseRetryPolicyFromMap(retryPolicy)
		}
		if sqsParams, ok := rawMap["sqsParameters"].(map[string]interface{}); ok {
			target.SqsParameters = &schedulerstore.SqsParameters{
				MessageGroupId: getStringFromMap(sqsParams, "messageGroupId", "MessageGroupId"),
			}
		} else if sqsParams, ok := rawMap["SqsParameters"].(map[string]interface{}); ok {
			target.SqsParameters = &schedulerstore.SqsParameters{
				MessageGroupId: getStringFromMap(sqsParams, "messageGroupId", "MessageGroupId"),
			}
		}
		if ecsParams, ok := rawMap["ecsParameters"].(map[string]interface{}); ok {
			target.EcsParameters = parseEcsParameters(ecsParams)
		} else if ecsParams, ok := rawMap["EcsParameters"].(map[string]interface{}); ok {
			target.EcsParameters = parseEcsParameters(ecsParams)
		}
		if ebParams, ok := rawMap["eventBridgeParameters"].(map[string]interface{}); ok {
			target.EventBridgeParameters = parseEventBridgeParameters(ebParams)
		} else if ebParams, ok := rawMap["EventBridgeParameters"].(map[string]interface{}); ok {
			target.EventBridgeParameters = parseEventBridgeParameters(ebParams)
		}
		if kinesisParams, ok := rawMap["kinesisParameters"].(map[string]interface{}); ok {
			target.KinesisParameters = parseKinesisParameters(kinesisParams)
		} else if kinesisParams, ok := rawMap["KinesisParameters"].(map[string]interface{}); ok {
			target.KinesisParameters = parseKinesisParameters(kinesisParams)
		}
	case map[string]interface{}:
		target.Arn = getStringFromMap(v, "arn", "Arn")
		target.RoleArn = getStringFromMap(v, "roleArn", "RoleArn")
		target.Input = getStringFromMap(v, "input", "Input")

		if dlConfig, ok := v["deadLetterConfig"].(map[string]interface{}); ok {
			target.DeadLetterConfig = &schedulerstore.DeadLetterConfig{
				Arn: getStringFromMap(dlConfig, "arn", "Arn"),
			}
		} else if dlConfig, ok := v["DeadLetterConfig"].(map[string]interface{}); ok {
			target.DeadLetterConfig = &schedulerstore.DeadLetterConfig{
				Arn: getStringFromMap(dlConfig, "arn", "Arn"),
			}
		}
		if retryPolicy, ok := v["retryPolicy"].(map[string]interface{}); ok {
			target.RetryPolicy = parseRetryPolicyFromMap(retryPolicy)
		} else if retryPolicy, ok := v["RetryPolicy"].(map[string]interface{}); ok {
			target.RetryPolicy = parseRetryPolicyFromMap(retryPolicy)
		}
		if sqsParams, ok := v["sqsParameters"].(map[string]interface{}); ok {
			target.SqsParameters = &schedulerstore.SqsParameters{
				MessageGroupId: getStringFromMap(sqsParams, "messageGroupId", "MessageGroupId"),
			}
		} else if sqsParams, ok := v["SqsParameters"].(map[string]interface{}); ok {
			target.SqsParameters = &schedulerstore.SqsParameters{
				MessageGroupId: getStringFromMap(sqsParams, "messageGroupId", "MessageGroupId"),
			}
		}
		if ecsParams, ok := v["ecsParameters"].(map[string]interface{}); ok {
			target.EcsParameters = parseEcsParameters(ecsParams)
		} else if ecsParams, ok := v["EcsParameters"].(map[string]interface{}); ok {
			target.EcsParameters = parseEcsParameters(ecsParams)
		}
		if ebParams, ok := v["eventBridgeParameters"].(map[string]interface{}); ok {
			target.EventBridgeParameters = parseEventBridgeParameters(ebParams)
		} else if ebParams, ok := v["EventBridgeParameters"].(map[string]interface{}); ok {
			target.EventBridgeParameters = parseEventBridgeParameters(ebParams)
		}
		if kinesisParams, ok := v["kinesisParameters"].(map[string]interface{}); ok {
			target.KinesisParameters = parseKinesisParameters(kinesisParams)
		} else if kinesisParams, ok := v["KinesisParameters"].(map[string]interface{}); ok {
			target.KinesisParameters = parseKinesisParameters(kinesisParams)
		}
	default:
		return nil, ErrInvalidTarget
	}

	if target.Arn == "" || target.RoleArn == "" {
		return nil, ErrInvalidTarget
	}

	return &target, nil
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
	if maxAge, ok := retryPolicy["maximumEventAgeInSeconds"].(float64); ok {
		val := int(maxAge)
		if val >= 60 && val <= 86400 {
			rp.MaximumEventAgeInSeconds = &val
		}
	} else if maxAge, ok := retryPolicy["MaximumEventAgeInSeconds"].(float64); ok {
		val := int(maxAge)
		if val >= 60 && val <= 86400 {
			rp.MaximumEventAgeInSeconds = &val
		}
	}
	if maxRetry, ok := retryPolicy["maximumRetryAttempts"].(float64); ok {
		val := int(maxRetry)
		if val >= 0 && val <= 185 {
			rp.MaximumRetryAttempts = &val
		}
	} else if maxRetry, ok := retryPolicy["MaximumRetryAttempts"].(float64); ok {
		val := int(maxRetry)
		if val >= 0 && val <= 185 {
			rp.MaximumRetryAttempts = &val
		}
	}
	return rp
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
	if taskCount, ok := data["taskCount"].(float64); ok {
		val := int(taskCount)
		params.TaskCount = &val
	}
	if taskCount, ok := data["TaskCount"].(float64); ok {
		val := int(taskCount)
		params.TaskCount = &val
	}
	if enabled, ok := data["enableECSManagedTags"].(bool); ok {
		params.EnableECSManagedTags = &enabled
	}
	if enabled, ok := data["EnableECSManagedTags"].(bool); ok {
		params.EnableECSManagedTags = &enabled
	}
	if enabled, ok := data["enableExecuteCommand"].(bool); ok {
		params.EnableExecuteCommand = &enabled
	}
	if enabled, ok := data["EnableExecuteCommand"].(bool); ok {
		params.EnableExecuteCommand = &enabled
	}
	if nc, ok := data["networkConfiguration"].(map[string]interface{}); ok {
		params.NetworkConfiguration = parseNetworkConfiguration(nc)
	}
	if nc, ok := data["NetworkConfiguration"].(map[string]interface{}); ok {
		params.NetworkConfiguration = parseNetworkConfiguration(nc)
	}
	if cps, ok := data["capacityProviderStrategy"].([]interface{}); ok {
		params.CapacityProviderStrategy = parseCapacityProviderStrategy(cps)
	}
	if cps, ok := data["CapacityProviderStrategy"].([]interface{}); ok {
		params.CapacityProviderStrategy = parseCapacityProviderStrategy(cps)
	}
	if pc, ok := data["placementConstraints"].([]interface{}); ok {
		params.PlacementConstraints = parsePlacementConstraints(pc)
	}
	if pc, ok := data["PlacementConstraints"].([]interface{}); ok {
		params.PlacementConstraints = parsePlacementConstraints(pc)
	}
	if ps, ok := data["placementStrategy"].([]interface{}); ok {
		params.PlacementStrategy = parsePlacementStrategy(ps)
	}
	if ps, ok := data["PlacementStrategy"].([]interface{}); ok {
		params.PlacementStrategy = parsePlacementStrategy(ps)
	}
	if tags, ok := data["tags"].([]interface{}); ok {
		params.Tags = tagutil.ParseEcsTags(tags)
	}
	if tags, ok := data["Tags"].([]interface{}); ok {
		params.Tags = tagutil.ParseEcsTags(tags)
	}
	return params
}

func parseNetworkConfiguration(data map[string]interface{}) *schedulerstore.NetworkConfiguration {
	nc := &schedulerstore.NetworkConfiguration{}
	if vpc, ok := data["awsvpcConfiguration"].(map[string]interface{}); ok {
		nc.AwsVpcConfiguration = parseAwsVpcConfiguration(vpc)
	} else if vpc, ok := data["AwsVpcConfiguration"].(map[string]interface{}); ok {
		nc.AwsVpcConfiguration = parseAwsVpcConfiguration(vpc)
	}
	return nc
}

func parseAwsVpcConfiguration(data map[string]interface{}) *schedulerstore.AwsVpcConfiguration {
	vpc := &schedulerstore.AwsVpcConfiguration{
		AssignPublicIp: getStringFromMap(data, "assignPublicIp", "AssignPublicIp"),
	}
	if subnets, ok := data["subnets"].([]interface{}); ok {
		for _, s := range subnets {
			if str, ok := s.(string); ok {
				vpc.Subnets = append(vpc.Subnets, str)
			}
		}
	} else if subnets, ok := data["Subnets"].([]interface{}); ok {
		for _, s := range subnets {
			if str, ok := s.(string); ok {
				vpc.Subnets = append(vpc.Subnets, str)
			}
		}
	}
	if sgs, ok := data["securityGroups"].([]interface{}); ok {
		for _, sg := range sgs {
			if str, ok := sg.(string); ok {
				vpc.SecurityGroups = append(vpc.SecurityGroups, str)
			}
		}
	} else if sgs, ok := data["SecurityGroups"].([]interface{}); ok {
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
			if w, ok := m["weight"].(float64); ok {
				val := int(w)
				cps.Weight = &val
			}
			if w, ok := m["Weight"].(float64); ok {
				val := int(w)
				cps.Weight = &val
			}
			if b, ok := m["base"].(float64); ok {
				val := int(b)
				cps.Base = &val
			}
			if b, ok := m["Base"].(float64); ok {
				val := int(b)
				cps.Base = &val
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
		if maxWindow, ok := v["maximumWindowInMinutes"].(float64); ok {
			val := int(maxWindow)
			ftw.MaximumWindowInMinutes = &val
		} else if maxWindow, ok := v["MaximumWindowInMinutes"].(float64); ok {
			val := int(maxWindow)
			ftw.MaximumWindowInMinutes = &val
		}
	default:
		return nil, ErrInvalidFlexibleTimeWindow
	}

	if ftw.Mode == "" {
		ftw.Mode = schedulerstore.FlexibleTimeWindowModeOff
	}

	return ftw, nil
}
