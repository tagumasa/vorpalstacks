// Package lambdautil holds the cross-service Lambda invocation contract.
// Gateways that trigger functions without importing the Lambda service
// package depend on this interface; the eventbus invoker provides the
// implementation.
package lambdautil

import "context"

// Invoker defines the interface for invoking Lambda functions.
type Invoker interface {
	InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error)
}
