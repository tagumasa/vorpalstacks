package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

func parseTarget(params map[string]interface{}) (*schedulerstore.Target, error) {
	targetData, ok := params["Target"]
	if !ok || targetData == nil {
		return nil, nil
	}

	rawMap, err := coerceToMap(targetData)
	if err != nil {
		return nil, ErrInvalidTarget
	}

	// Arn and RoleArn are required members of the Target shape; their
	// absence is rejected by validateTarget in the Core validation path.
	target, err := parseTargetFromMap(rawMap)
	if err != nil {
		return nil, err
	}
	return &target, nil
}

// parseMaxResultsParam reads an optional integer MaxResults member. It only
// converts the wire text; the range check lives in the Core validation path
// (resolveListMaxResults).
func parseMaxResultsParam(params map[string]interface{}) (*int32, error) {
	raw := request.GetStringParam(params, "MaxResults")
	if raw == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return nil, ErrValidation
	}
	v := int32(parsed)
	return &v, nil
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

// parseTargetFromMap binds the Target body members to the member names the
// model declares (Pascal at the Target level; the nested shapes use the
// names the model gives them — see each nested parser). A present member
// whose JSON type does not match the shape is a wire-format violation
// reported to the caller, never a silently dropped value the Core validator
// cannot see; an absent member flows through as the zero value and is
// rejected by the Core validation path when the shape requires it.
func parseTargetFromMap(m map[string]interface{}) (schedulerstore.Target, error) {
	var target schedulerstore.Target
	var err error

	if target.Arn, err = stringMember(m, "Arn", "Target.Arn"); err != nil {
		return target, err
	}
	if target.RoleArn, err = stringMember(m, "RoleArn", "Target.RoleArn"); err != nil {
		return target, err
	}
	if target.Input, err = stringMember(m, "Input", "Target.Input"); err != nil {
		return target, err
	}

	dl, err := mapMember(m, "DeadLetterConfig", "Target.DeadLetterConfig")
	if err != nil {
		return target, err
	}
	if dl != nil {
		arn, err := stringMember(dl, "Arn", "Target.DeadLetterConfig.Arn")
		if err != nil {
			return target, err
		}
		target.DeadLetterConfig = &schedulerstore.DeadLetterConfig{Arn: arn}
	}
	rp, err := mapMember(m, "RetryPolicy", "Target.RetryPolicy")
	if err != nil {
		return target, err
	}
	if rp != nil {
		if target.RetryPolicy, err = parseRetryPolicyFromMap(rp); err != nil {
			return target, err
		}
	}
	sqs, err := mapMember(m, "SqsParameters", "Target.SqsParameters")
	if err != nil {
		return target, err
	}
	if sqs != nil {
		if target.SqsParameters, err = parseSqsParameters(sqs); err != nil {
			return target, err
		}
	}
	ecs, err := mapMember(m, "EcsParameters", "Target.EcsParameters")
	if err != nil {
		return target, err
	}
	if ecs != nil {
		if target.EcsParameters, err = parseEcsParameters(ecs); err != nil {
			return target, err
		}
	}
	eb, err := mapMember(m, "EventBridgeParameters", "Target.EventBridgeParameters")
	if err != nil {
		return target, err
	}
	if eb != nil {
		if target.EventBridgeParameters, err = parseEventBridgeParameters(eb); err != nil {
			return target, err
		}
	}
	kinesis, err := mapMember(m, "KinesisParameters", "Target.KinesisParameters")
	if err != nil {
		return target, err
	}
	if kinesis != nil {
		if target.KinesisParameters, err = parseKinesisParameters(kinesis); err != nil {
			return target, err
		}
	}
	// SageMaker templated targets are permanently out of scope on this
	// platform, so a present SageMakerPipelineParameters member is never a
	// carryable value on an accepted target — it is reported like every
	// other sub-parameter the target's service cannot use, instead of being
	// silently dropped where the Core validator cannot see it.
	if _, ok := m["SageMakerPipelineParameters"]; ok {
		return target, awserrors.NewValidationException(
			"SageMakerPipelineParameters can only be specified for SageMaker targets; SageMaker targets are not supported on this platform")
	}
	return target, nil
}

// stringMember reads one modelled string member: an absent or null member
// yields "" (required members are rejected by the Core validation path),
// while a present value of another JSON type is a wire-format violation.
// The path names the member for the error message.
func stringMember(m map[string]interface{}, key, path string) (string, error) {
	v, ok := m[key]
	if !ok || v == nil {
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", awserrors.NewValidationException(fmt.Sprintf("%s must be a string", path))
	}
	return s, nil
}

// mapMember reads one modelled structure member under the same absent/null
// and type-mismatch rules as stringMember.
func mapMember(m map[string]interface{}, key, path string) (map[string]interface{}, error) {
	v, ok := m[key]
	if !ok || v == nil {
		return nil, nil
	}
	mp, ok := v.(map[string]interface{})
	if !ok {
		return nil, awserrors.NewValidationException(fmt.Sprintf("%s must be a structure", path))
	}
	return mp, nil
}

// sliceMember reads one modelled list member under the same absent/null and
// type-mismatch rules as stringMember.
func sliceMember(m map[string]interface{}, key, path string) ([]interface{}, error) {
	v, ok := m[key]
	if !ok || v == nil {
		return nil, nil
	}
	sl, ok := v.([]interface{})
	if !ok {
		return nil, awserrors.NewValidationException(fmt.Sprintf("%s must be a list", path))
	}
	return sl, nil
}

// numberMember reads one modelled numeric member; the boolean reports
// presence so a zero stays distinguishable from an absent member.
// encoding/json decodes every JSON number to float64 — the int/int32/int64
// arms serve callers that build the map with Go-native values. A
// fractional JSON number is a wire violation for the modelled integer
// shapes (restJson1), reported rather than truncated into a reshaped
// value.
func numberMember(m map[string]interface{}, key, path string) (int, bool, error) {
	v, ok := m[key]
	if !ok || v == nil {
		return 0, false, nil
	}
	switch n := v.(type) {
	case float64:
		if n != math.Trunc(n) {
			return 0, false, awserrors.NewValidationException(fmt.Sprintf("%s must be an integer", path))
		}
		return int(n), true, nil
	case int:
		return n, true, nil
	case int32:
		return int(n), true, nil
	case int64:
		return int(n), true, nil
	default:
		return 0, false, awserrors.NewValidationException(fmt.Sprintf("%s must be a number", path))
	}
}

// boolMember reads one modelled boolean member; nil means absent.
func boolMember(m map[string]interface{}, key, path string) (*bool, error) {
	v, ok := m[key]
	if !ok || v == nil {
		return nil, nil
	}
	b, ok := v.(bool)
	if !ok {
		return nil, awserrors.NewValidationException(fmt.Sprintf("%s must be a boolean", path))
	}
	return &b, nil
}

// stringSlice binds a modelled list of strings; a non-string entry is a
// wire-format violation rather than a silently dropped value.
func stringSlice(items []interface{}, path string) ([]string, error) {
	out := make([]string, 0, len(items))
	for i, item := range items {
		s, ok := item.(string)
		if !ok {
			return nil, awserrors.NewValidationException(fmt.Sprintf("%s[%d] must be a string", path, i))
		}
		out = append(out, s)
	}
	return out, nil
}

func parseRetryPolicyFromMap(retryPolicy map[string]interface{}) (*schedulerstore.RetryPolicy, error) {
	rp := &schedulerstore.RetryPolicy{}
	// Accept the raw values without range filtering. Range validation is
	// performed by validateTarget in validators.go.
	if val, ok, err := numberMember(retryPolicy, "MaximumEventAgeInSeconds", "RetryPolicy.MaximumEventAgeInSeconds"); err != nil {
		return nil, err
	} else if ok {
		rp.MaximumEventAgeInSeconds = &val
	}
	if val, ok, err := numberMember(retryPolicy, "MaximumRetryAttempts", "RetryPolicy.MaximumRetryAttempts"); err != nil {
		return nil, err
	} else if ok {
		rp.MaximumRetryAttempts = &val
	}
	return rp, nil
}

// parseSqsParameters binds SqsParameters to its single modelled member.
// MessageGroupId carries @length(1, 128): an explicitly present empty
// string is out of bounds on the wire, where presence is still visible —
// distinct from an absent member, which leaves the DTO empty and lets
// delivery fall back to the schedule name for FIFO queues.
func parseSqsParameters(m map[string]interface{}) (*schedulerstore.SqsParameters, error) {
	sqs := &schedulerstore.SqsParameters{}
	v, ok := m["MessageGroupId"]
	if !ok || v == nil {
		return sqs, nil
	}
	s, ok := v.(string)
	if !ok {
		return nil, awserrors.NewValidationException("SqsParameters.MessageGroupId must be a string")
	}
	if s == "" {
		return nil, awserrors.NewValidationException(fmt.Sprintf(
			"SqsParameters.MessageGroupId must be 1-%d characters", MaxMessageGroupIdLength))
	}
	sqs.MessageGroupId = s
	return sqs, nil
}

func parseEcsParameters(data map[string]interface{}) (*schedulerstore.EcsParameters, error) {
	params := &schedulerstore.EcsParameters{}
	var err error
	if params.TaskDefinitionArn, err = stringMember(data, "TaskDefinitionArn", "EcsParameters.TaskDefinitionArn"); err != nil {
		return nil, err
	}
	if params.LaunchType, err = stringMember(data, "LaunchType", "EcsParameters.LaunchType"); err != nil {
		return nil, err
	}
	if params.PlatformVersion, err = stringMember(data, "PlatformVersion", "EcsParameters.PlatformVersion"); err != nil {
		return nil, err
	}
	if params.Group, err = stringMember(data, "Group", "EcsParameters.Group"); err != nil {
		return nil, err
	}
	if params.PropagateTags, err = stringMember(data, "PropagateTags", "EcsParameters.PropagateTags"); err != nil {
		return nil, err
	}
	if params.ReferenceId, err = stringMember(data, "ReferenceId", "EcsParameters.ReferenceId"); err != nil {
		return nil, err
	}
	if val, ok, err := numberMember(data, "TaskCount", "EcsParameters.TaskCount"); err != nil {
		return nil, err
	} else if ok {
		params.TaskCount = &val
	}
	if val, err := boolMember(data, "EnableECSManagedTags", "EcsParameters.EnableECSManagedTags"); err != nil {
		return nil, err
	} else if val != nil {
		params.EnableECSManagedTags = val
	}
	if val, err := boolMember(data, "EnableExecuteCommand", "EcsParameters.EnableExecuteCommand"); err != nil {
		return nil, err
	} else if val != nil {
		params.EnableExecuteCommand = val
	}
	nc, err := mapMember(data, "NetworkConfiguration", "EcsParameters.NetworkConfiguration")
	if err != nil {
		return nil, err
	}
	if nc != nil {
		if params.NetworkConfiguration, err = parseNetworkConfiguration(nc); err != nil {
			return nil, err
		}
	}
	cps, err := sliceMember(data, "CapacityProviderStrategy", "EcsParameters.CapacityProviderStrategy")
	if err != nil {
		return nil, err
	}
	if cps != nil {
		if params.CapacityProviderStrategy, err = parseCapacityProviderStrategy(cps); err != nil {
			return nil, err
		}
	}
	pc, err := sliceMember(data, "PlacementConstraints", "EcsParameters.PlacementConstraints")
	if err != nil {
		return nil, err
	}
	if pc != nil {
		if params.PlacementConstraints, err = parsePlacementConstraints(pc); err != nil {
			return nil, err
		}
	}
	ps, err := sliceMember(data, "PlacementStrategy", "EcsParameters.PlacementStrategy")
	if err != nil {
		return nil, err
	}
	if ps != nil {
		if params.PlacementStrategy, err = parsePlacementStrategy(ps); err != nil {
			return nil, err
		}
	}
	tags, err := sliceMember(data, "Tags", "EcsParameters.Tags")
	if err != nil {
		return nil, err
	}
	if tags != nil {
		if params.Tags, err = parseEcsTags(tags); err != nil {
			return nil, err
		}
	}
	return params, nil
}

// parseNetworkConfiguration binds the NetworkConfiguration member. The
// modelled shape carries exactly one member, awsvpcConfiguration (the
// model's own lowerCamel name), so an absent sub-member means the structure
// itself is absent — the empty wrapper is never stored, because restJson1
// omits an absent structure rather than round-tripping it as null.
func parseNetworkConfiguration(data map[string]interface{}) (*schedulerstore.NetworkConfiguration, error) {
	// Absent and null both mean the structure is omitted (restJson1 never
	// round-trips an absent structure as null); any other value whose JSON
	// type contradicts the shape is a wire-format violation reported like
	// every other member, never a silently dropped configuration the Core
	// validator cannot see.
	v, present := data["awsvpcConfiguration"]
	if !present || v == nil {
		return nil, nil
	}
	vpc, ok := v.(map[string]interface{})
	if !ok {
		return nil, awserrors.NewValidationException("NetworkConfiguration.awsvpcConfiguration must be a structure")
	}
	parsed, err := parseAwsVpcConfiguration(vpc)
	if err != nil {
		return nil, err
	}
	return &schedulerstore.NetworkConfiguration{
		AwsVpcConfiguration: parsed,
	}, nil
}

func parseAwsVpcConfiguration(data map[string]interface{}) (*schedulerstore.AwsVpcConfiguration, error) {
	vpc := &schedulerstore.AwsVpcConfiguration{}
	var err error
	if vpc.AssignPublicIp, err = stringMember(data, "AssignPublicIp", "NetworkConfiguration.awsvpcConfiguration.AssignPublicIp"); err != nil {
		return nil, err
	}
	subnets, err := sliceMember(data, "Subnets", "NetworkConfiguration.awsvpcConfiguration.Subnets")
	if err != nil {
		return nil, err
	}
	if subnets != nil {
		if vpc.Subnets, err = stringSlice(subnets, "NetworkConfiguration.awsvpcConfiguration.Subnets"); err != nil {
			return nil, err
		}
	}
	sgs, err := sliceMember(data, "SecurityGroups", "NetworkConfiguration.awsvpcConfiguration.SecurityGroups")
	if err != nil {
		return nil, err
	}
	if sgs != nil {
		// stringSlice allocates before filling so an explicitly empty list
		// stays distinguishable from an omitted member: the model constrains
		// a provided SecurityGroups list to at least one entry, while an
		// omitted member selects the VPC default security group.
		if vpc.SecurityGroups, err = stringSlice(sgs, "NetworkConfiguration.awsvpcConfiguration.SecurityGroups"); err != nil {
			return nil, err
		}
	}
	return vpc, nil
}

func parseCapacityProviderStrategy(data []interface{}) ([]schedulerstore.CapacityProviderStrategyItem, error) {
	var result []schedulerstore.CapacityProviderStrategyItem
	for i, item := range data {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.CapacityProviderStrategy[%d] must be a structure", i))
		}
		cps := schedulerstore.CapacityProviderStrategyItem{}
		var err error
		if cps.CapacityProvider, err = stringMember(m, "capacityProvider", "CapacityProviderStrategy.capacityProvider"); err != nil {
			return nil, err
		}
		if w, ok, err := numberMember(m, "weight", "CapacityProviderStrategy.weight"); err != nil {
			return nil, err
		} else if ok {
			cps.Weight = &w
		}
		if b, ok, err := numberMember(m, "base", "CapacityProviderStrategy.base"); err != nil {
			return nil, err
		} else if ok {
			cps.Base = &b
		}
		result = append(result, cps)
	}
	return result, nil
}

func parsePlacementConstraints(data []interface{}) ([]schedulerstore.PlacementConstraint, error) {
	var result []schedulerstore.PlacementConstraint
	for i, item := range data {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.PlacementConstraints[%d] must be a structure", i))
		}
		var pc schedulerstore.PlacementConstraint
		var err error
		if pc.Type, err = stringMember(m, "type", "PlacementConstraints.type"); err != nil {
			return nil, err
		}
		if pc.Expression, err = stringMember(m, "expression", "PlacementConstraints.expression"); err != nil {
			return nil, err
		}
		result = append(result, pc)
	}
	return result, nil
}

func parsePlacementStrategy(data []interface{}) ([]schedulerstore.PlacementStrategy, error) {
	var result []schedulerstore.PlacementStrategy
	for i, item := range data {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.PlacementStrategy[%d] must be a structure", i))
		}
		var ps schedulerstore.PlacementStrategy
		var err error
		if ps.Type, err = stringMember(m, "type", "PlacementStrategy.type"); err != nil {
			return nil, err
		}
		if ps.Field, err = stringMember(m, "field", "PlacementStrategy.field"); err != nil {
			return nil, err
		}
		result = append(result, ps)
	}
	return result, nil
}

func parseEventBridgeParameters(data map[string]interface{}) (*schedulerstore.EventBridgeParameters, error) {
	eb := &schedulerstore.EventBridgeParameters{}
	var err error
	if eb.DetailType, err = stringMember(data, "DetailType", "EventBridgeParameters.DetailType"); err != nil {
		return nil, err
	}
	if eb.Source, err = stringMember(data, "Source", "EventBridgeParameters.Source"); err != nil {
		return nil, err
	}
	return eb, nil
}

func parseKinesisParameters(data map[string]interface{}) (*schedulerstore.KinesisParameters, error) {
	kin := &schedulerstore.KinesisParameters{}
	var err error
	if kin.PartitionKey, err = stringMember(data, "PartitionKey", "KinesisParameters.PartitionKey"); err != nil {
		return nil, err
	}
	return kin, nil
}

// parseFlexibleTimeWindow binds the FlexibleTimeWindow body member: the
// wire key and the nested Mode/MaximumWindowInMinutes members are the
// Pascal names the model declares. The value may arrive as a structure or
// as a JSON string carrying one; both bind through the same members.
func parseFlexibleTimeWindow(params map[string]interface{}) (*schedulerstore.FlexibleTimeWindow, error) {
	ftwData, ok := params["FlexibleTimeWindow"]
	if !ok || ftwData == nil {
		return nil, nil
	}

	var m map[string]interface{}
	switch v := ftwData.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, ErrInvalidFlexibleTimeWindow
		}
	case map[string]interface{}:
		m = v
	default:
		return nil, ErrInvalidFlexibleTimeWindow
	}

	ftw := &schedulerstore.FlexibleTimeWindow{}
	mode, err := stringMember(m, "Mode", "FlexibleTimeWindow.Mode")
	if err != nil {
		return nil, err
	}
	if mode != "" {
		ftw.Mode = schedulerstore.FlexibleTimeWindowMode(mode)
	}
	if maxWindow, ok, err := numberMember(m, "MaximumWindowInMinutes", "FlexibleTimeWindow.MaximumWindowInMinutes"); err != nil {
		return nil, err
	} else if ok {
		ftw.MaximumWindowInMinutes = &maxWindow
	}

	// Mode is a required member of the FlexibleTimeWindow shape; an empty
	// Mode is rejected by validateFlexibleTimeWindow in validators.go.

	return ftw, nil
}

// validateVpcConfig validates the AwsVpcConfiguration subnets and security
// groups against the EC2 service via the event bus. All resources must exist
// and belong to the same VPC. Accepts region directly so both the HTTP API
// and admin console paths can call it.
func (s *SchedulerService) validateVpcConfig(ctx context.Context, region string, target *schedulerstore.Target) error {
	// The engine is attached by BuildEngine, which runs before the service
	// takes traffic in the running server. A service without one (unit
	// harnesses, the construction window) skips the cross-service existence
	// check — the modelled trait validation in validators.go still applies.
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

	// A built engine whose bus carries no EC2 invoker is a deployment whose
	// EC2 service is unavailable: the creation request fails closed instead
	// of persisting a VPC configuration nobody verified.
	ec2 := s.engine.bus.EC2Invoker()
	if ec2 == nil {
		return awserrors.NewValidationException("scheduler: EC2 service not available for VPC configuration validation")
	}

	for _, subnetId := range vpc.Subnets {
		if _, _, err := ec2.LookupSubnet(ctx, region, subnetId); err != nil {
			return awserrors.NewValidationException(fmt.Sprintf("scheduler: subnet %s not found: %v", subnetId, err))
		}
	}

	for _, sgId := range vpc.SecurityGroups {
		if _, err := ec2.LookupSecurityGroup(ctx, region, sgId); err != nil {
			return awserrors.NewValidationException(fmt.Sprintf("scheduler: security group %s not found: %v", sgId, err))
		}
	}

	return nil
}

// validateKmsKey checks that a customer managed KMS key referenced by
// KmsKeyArn exists. AWS validates the key when the schedule is created:
// the Encryption at rest documentation requires kms:DescribeKey on the
// principal "that calls the EventBridge Scheduler API when creating a
// schedule" — "Required in order to validate that the key you provide is
// a symmetric encryption KMS key."
func (s *SchedulerService) validateKmsKey(ctx context.Context, region, kmsKeyArn string) error {
	if kmsKeyArn == "" {
		return nil
	}
	// The engine is attached by BuildEngine, which runs before the service
	// takes traffic in the running server. A service without one (unit
	// harnesses, the construction window) skips the cross-service key
	// check — the KmsKeyArn @pattern/@length validation in validators.go
	// still applies.
	if s.engine == nil || s.engine.bus == nil {
		return nil
	}
	// A built engine whose bus carries no KMS invoker is a deployment whose
	// KMS service is unavailable: the creation request fails closed (the
	// posture validateVpcConfig already applied for EC2) instead of
	// accepting a key ARN nobody can verify.
	kms := s.engine.bus.KMSInvoker()
	if kms == nil {
		return awserrors.NewValidationException("scheduler: KMS service not available for KMS key validation")
	}
	if !kms.KeyExists(ctx, kmsKeyArn) {
		return awserrors.NewValidationException(fmt.Sprintf("scheduler: KMS key %s does not exist", kmsKeyArn))
	}
	// AWS validates the key class at schedule creation through
	// kms:DescribeKey: only a symmetric encryption key may protect a
	// schedule's data.
	if !kms.SymmetricEncryptionKeyExists(ctx, kmsKeyArn) {
		return awserrors.NewValidationException(fmt.Sprintf("scheduler: KMS key %s is not a symmetric encryption KMS key", kmsKeyArn))
	}
	return nil
}

// claimClientToken owns the ClientToken idempotency protocol shared by the
// create, update and delete cores: an empty token claims nothing, an invalid
// token is a validation error, and a valid token is claimed for the
// (scope, resourceArn) pair. A replay returns the first application's
// resource ARN with a no-op release; a fresh claim returns a release closure
// that rolls the claim back — it must be called on every error path after
// the claim, and is idempotent-safe to skip on success (the entry expires
// with the token TTL).
func claimClientToken(store *schedulerstore.SchedulerStore, token, resourceArn, scope string) (replayedArn string, release func(), err error) {
	if token == "" {
		return "", func() {}, nil
	}
	if err := validateClientToken(token); err != nil {
		return "", nil, err
	}
	if entry, created := store.ClientTokens().LookupOrClaim(token, resourceArn, scope); !created {
		return entry.ResourceArn, func() {}, nil
	}
	return "", func() { store.ClientTokens().Release(token, resourceArn, scope) }, nil
}
