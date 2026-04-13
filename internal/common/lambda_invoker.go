// Package common provides common AWS service utilities for vorpalstacks.
package common

import "context"

// LambdaInvoker defines the interface for invoking Lambda functions
type LambdaInvoker interface {
	InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error)
}
