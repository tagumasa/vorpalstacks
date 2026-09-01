package eventbridge

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

const maxPutEventsEntries = 10

// PutEvents delivers one or more events to EventBridge.
// Validates required fields (Source, DetailType, Detail) and delivers
// events to matching rules on the specified event bus.
func (s *EventsService) PutEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	entries, ok := req.Parameters["Entries"].([]interface{})
	if !ok {
		entries, ok = req.Parameters["entries"].([]interface{})
	}
	if !ok || len(entries) == 0 {
		return nil, awserrors.NewValidationException("Entries are required")
	}
	if len(entries) > maxPutEventsEntries {
		return nil, awserrors.NewValidationException("Maximum 10 entries allowed per request")
	}

	// EndpointId routes the request through a global endpoint. Global
	// endpoints are out of scope for this edge platform, but we still
	// accept the parameter so SDK clients do not receive an unexpected
	// ValidationException when populating it.
	if endpointID, _ := req.Parameters["EndpointId"].(string); endpointID != "" {
		logs.Debug("PutEvents EndpointId ignored (global endpoints out of scope)",
			logs.String("endpointId", endpointID))
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.putEventsCore(ctx, store, PutEventsInput{
		Entries: entries,
		Region:  reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"FailedEntryCount": result.FailedEntryCount,
		"Entries":          result.Entries,
	}, nil
}

func (s *EventsService) buildTargetPayload(event *eventsstore.Event, target eventsstore.Target) map[string]interface{} {
	payload := map[string]interface{}{
		"version":     event.Version,
		"id":          event.ID,
		"detail-type": event.DetailType,
		"source":      event.Source,
		"account":     event.Account,
		"time":        event.Time,
		"region":      event.Region,
		"resources":   event.Resources,
		"detail":      event.Detail,
	}

	if target.Input != "" {
		var inputPayload map[string]interface{}
		if err := json.Unmarshal([]byte(target.Input), &inputPayload); err == nil {
			payload = inputPayload
		}
	} else if target.InputPath != "" {
		extracted := s.extractInputPath(payload, target.InputPath)
		if extracted != nil {
			payload = extracted
		}
	} else if target.InputTransformer != nil {
		transformed := s.applyInputTransformer(payload, target.InputTransformer)
		if transformed != nil {
			payload = transformed
		}
	}

	return payload
}

func (s *EventsService) TestEventPattern(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	patternStr := request.GetStringParam(req.Parameters, "EventPattern")
	eventStr := request.GetStringParam(req.Parameters, "Event")

	result, err := s.testEventPatternCore(patternStr, eventStr)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Result": result,
	}, nil
}

func generateEventID() string {
	return uuid.New().String()
}
