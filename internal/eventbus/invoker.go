package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
)

// ServiceInvoker defines the contract for invoking a target service as
// part of event dispatch.
type ServiceInvoker interface {
	Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult
	ServiceType() string
}

// EC2SubnetLookup provides subnet and security group metadata resolution for
// cross-service consumers that need to map resource IDs to their parent VPC.
type EC2SubnetLookup interface {
	LookupSubnet(ctx context.Context, region string, subnetId string) (vpcId string, availabilityZone string, err error)
	LookupSecurityGroup(ctx context.Context, region string, groupId string) (vpcId string, err error)
}

// EC2InvokerAdapter adapts an EC2SubnetLookup to the ServiceInvoker
// interface, exposing subnet and security group queries through the event bus.
type EC2InvokerAdapter struct {
	Lookup EC2SubnetLookup
}

// Invoke dispatches an EC2 lookup action. The only supported action types are
// "LookupSubnet" and "LookupSecurityGroup", which resolve resource IDs to
// their VPC metadata, returning a JSON payload.
func (a *EC2InvokerAdapter) Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult {
	switch action.Type {
	case "LookupSubnet":
		vpcId, az, err := a.Lookup.LookupSubnet(ctx, action.Region, action.Identifier)
		if err != nil {
			return HandlerResult{Error: err}
		}
		result := map[string]string{"VpcId": vpcId, "AvailabilityZone": az}
		data, err := json.Marshal(result)
		if err != nil {
			return HandlerResult{Error: fmt.Errorf("ec2: failed to marshal subnet lookup result: %w", err)}
		}
		return HandlerResult{StatusCode: 200, Payload: data}
	case "LookupSecurityGroup":
		vpcId, err := a.Lookup.LookupSecurityGroup(ctx, action.Region, action.Identifier)
		if err != nil {
			return HandlerResult{Error: err}
		}
		result := map[string]string{"VpcId": vpcId}
		data, err := json.Marshal(result)
		if err != nil {
			return HandlerResult{Error: fmt.Errorf("ec2: failed to marshal security group lookup result: %w", err)}
		}
		return HandlerResult{StatusCode: 200, Payload: data}
	default:
		return HandlerResult{Error: fmt.Errorf("ec2: unsupported invoker action type: %s", action.Type)}
	}
}

// ServiceType returns "ec2" for this adapter.
func (a *EC2InvokerAdapter) ServiceType() string { return "ec2" }
