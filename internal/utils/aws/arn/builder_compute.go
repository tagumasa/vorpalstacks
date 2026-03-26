// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

import (
	"fmt"
	"strconv"
	"strings"
)

// LambdaBuilder provides methods for constructing Lambda ARNs.
type LambdaBuilder struct{ *ARNBuilder }

// Lambda returns a LambdaBuilder for constructing Lambda ARNs.
func (b *ARNBuilder) Lambda() *LambdaBuilder { return &LambdaBuilder{b} }

// Function constructs an ARN for a Lambda function.
func (b *LambdaBuilder) Function(name string) string { return b.Build("lambda", "function:"+name) }

// FunctionVersion constructs an ARN for a specific version of a Lambda function.
func (b *LambdaBuilder) FunctionVersion(name, version string) string {
	return b.Build("lambda", "function:"+name+":"+version)
}

// FunctionAlias constructs an ARN for a Lambda function alias.
func (b *LambdaBuilder) FunctionAlias(name, alias string) string {
	return b.Build("lambda", "function:"+name+":"+alias)
}

// Layer constructs an ARN for a Lambda layer.
func (b *LambdaBuilder) Layer(name string) string { return b.Build("lambda", "layer:"+name) }

// LayerVersion constructs an ARN for a specific version of a Lambda layer.
func (b *LambdaBuilder) LayerVersion(name string, version int64) string {
	return b.Build("lambda", fmt.Sprintf("layer:%s:%d", name, version))
}

// EventSourceMapping constructs an ARN for a Lambda event source mapping.
func (b *LambdaBuilder) EventSourceMapping(uuid string) string {
	return b.Build("lambda", "event-source-mapping:"+uuid)
}

// CodeSigningConfig constructs an ARN for a Lambda code signing configuration.
func (b *LambdaBuilder) CodeSigningConfig(id string) string {
	return b.Build("lambda", "code-signing-config:"+id)
}

// ParseFunctionName extracts the function name from a Lambda function ARN.
func (b *LambdaBuilder) ParseFunctionName(arn string) string { return ExtractFunctionNameFromARN(arn) }

// ParseLayerName extracts the layer name from a Lambda layer ARN.
func (b *LambdaBuilder) ParseLayerName(arn string) string {
	parts := strings.Split(arn, ":")
	for i, p := range parts {
		if p == "layer" && i+1 < len(parts) {
			return strings.Split(parts[i+1], ":")[0]
		}
	}
	return ""
}

// ParseLayerVersion extracts the version number from a Lambda layer ARN.
func (b *LambdaBuilder) ParseLayerVersion(arn string) int64 {
	parts := strings.Split(arn, ":")
	for i, p := range parts {
		if p == "layer" && i+1 < len(parts) {
			subParts := strings.Split(parts[i+1], ":")
			if len(subParts) > 1 {
				if v, err := strconv.ParseInt(subParts[1], 10, 64); err == nil {
					return v
				}
			}
		}
	}
	return 0
}

// StepFunctionsBuilder provides methods for constructing Step Functions ARNs.
type StepFunctionsBuilder struct{ *ARNBuilder }

// StepFunctions returns a StepFunctionsBuilder for constructing Step Functions ARNs.
func (b *ARNBuilder) StepFunctions() *StepFunctionsBuilder { return &StepFunctionsBuilder{b} }

// StateMachine constructs an ARN for a Step Functions state machine.
func (b *StepFunctionsBuilder) StateMachine(name string) string {
	return b.Build("states", "stateMachine:"+name)
}

// Execution constructs an ARN for a Step Functions execution.
func (b *StepFunctionsBuilder) Execution(smName, execName string) string {
	return b.Build("states", fmt.Sprintf("execution:%s:%s", smName, execName))
}

// Activity constructs an ARN for a Step Functions activity.
func (b *StepFunctionsBuilder) Activity(name string) string {
	return b.Build("states", "activity:"+name)
}

// ParseStateMachineName extracts the state machine name from a Step Functions ARN.
func (b *StepFunctionsBuilder) ParseStateMachineName(arn string) string {
	return ExtractStateMachineNameFromARN(arn)
}

// ParseExecutionName extracts the execution name from a Step Functions execution ARN.
func (b *StepFunctionsBuilder) ParseExecutionName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "execution:") {
		parts := strings.Split(strings.TrimPrefix(resource, "execution:"), ":")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return ""
}
