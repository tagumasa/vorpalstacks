package lambda

// Payload size limits per the API model's Invoke documentation: "The
// maximum payload size is 6 MB for synchronous invocations and 1 MB for
// asynchronous invocations."
const (
	syncMaxPayloadSize  = 6 * 1024 * 1024 // 6 MB for synchronous invocation
	asyncMaxPayloadSize = 1024 * 1024     // 1 MB for asynchronous invocation
)

// InvokeRequest carries the wire members shared by the invocation
// handlers: the enum members from the parameters and the raw httpPayload
// body.
type InvokeRequest struct {
	InvocationType string
	LogType        string
	Payload        []byte
}

// InvocationMode names the transport execution mode of a validated
// invocation request.
type InvocationMode string

const (
	// InvocationModeSync executes the function synchronously and answers
	// with its result.
	InvocationModeSync InvocationMode = "sync"
	// InvocationModeEvent dispatches the function asynchronously and
	// answers 202 immediately.
	InvocationModeEvent InvocationMode = "event"
	// InvocationModeDryRun validates the request and answers 204 without
	// invoking the function.
	InvocationModeDryRun InvocationMode = "dry-run"
)

// prepareInvocationCore validates an Invoke request against the modelled
// contract — the InvocationType and LogType enums and the payload cap of
// the selected mode — and classifies the mode the transport executes.
func prepareInvocationCore(in *InvokeRequest) (InvocationMode, error) {
	if err := validateInvocationType(in.InvocationType); err != nil {
		return "", err
	}
	switch in.InvocationType {
	case "DryRun":
		return InvocationModeDryRun, nil
	case "Event":
		if err := validateInvokeAsyncCore(in); err != nil {
			return "", err
		}
		return InvocationModeEvent, nil
	default:
		return validateInvokeSyncCore(in)
	}
}

// validateInvokeSyncCore validates a synchronous invocation request: the
// 6 MB payload cap and the LogType enum.
func validateInvokeSyncCore(in *InvokeRequest) (InvocationMode, error) {
	if len(in.Payload) > syncMaxPayloadSize {
		return "", ErrRequestTooLarge
	}
	if err := validateLogType(in.LogType); err != nil {
		return "", err
	}
	return InvocationModeSync, nil
}

// validateInvokeAsyncCore validates an asynchronous invocation request:
// the 1 MB payload cap.
func validateInvokeAsyncCore(in *InvokeRequest) error {
	if len(in.Payload) > asyncMaxPayloadSize {
		return ErrRequestTooLarge
	}
	return nil
}
