// Package conditions provides AWS-specific utility functions for vorpalstacks.
package conditions

import (
	"time"
)

// BuildEvaluationContextFromMap creates an evaluation context from a map of values.
// This is useful for testing or when you have context information already extracted.
//
// Parameters:
//   - values: The map of values to create context from
//
// Returns:
//   - *EvaluationContext: The evaluation context
func BuildEvaluationContextFromMap(values map[string]string) *EvaluationContext {
	ctx := &EvaluationContext{
		CurrentTime:    time.Now(),
		RequestHeaders: make(map[string]string),
		Environment:    make(map[string]string),
	}

	for k, v := range values {
		switch k {
		case "SourceIP":
			ctx.SourceIP = v
		case "UserAgent":
			ctx.UserAgent = v
		case "Referer":
			ctx.Referer = v
		case "PrincipalARN":
			ctx.PrincipalARN = v
		case "PrincipalID":
			ctx.PrincipalID = v
		case "ResourceARN":
			ctx.ResourceARN = v
		case "Action":
			ctx.Action = v
		}
		ctx.RequestHeaders[k] = v
	}

	return ctx
}
