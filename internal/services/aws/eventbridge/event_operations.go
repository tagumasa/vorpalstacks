package eventbridge

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// PutEvents delivers one or more events to EventBridge.
// Validates required fields (Source, DetailType, Detail) and delivers
// events to matching rules on the specified event bus.
func (s *EventsService) PutEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	entries, ok := req.Parameters["Entries"].([]interface{})
	if !ok {
		entries, ok = req.Parameters["entries"].([]interface{})
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

// buildTargetPayload renders the payload one target receives: the
// constant Input verbatim, the InputPath extraction of the matched event,
// the transformer output, or — with no input configuration — the whole
// event envelope. ruleARN and ruleName feed the transformer's reserved
// variables.
func (s *EventsService) buildTargetPayload(ruleARN, ruleName string, event *eventsstore.Event, target eventsstore.Target) []byte {
	switch {
	case target.Input != "":
		// Input is valid JSON text validated at PutTargets time and
		// overrides the event entirely; it is passed through verbatim so
		// non-object JSON constants survive.
		return []byte(target.Input)
	case target.InputPath != "":
		if value, ok := resolveEventPathString(eventEnvelope(event), target.InputPath); ok {
			if b, err := json.Marshal(value); err == nil {
				return b
			}
		}
		// An unresolvable InputPath passes the whole event: AWS validates
		// no path against the event shape ("There is no validation when
		// creating JSON path for your template"), so misses are tolerated
		// at delivery and the no-input default — the entire event — holds.
		return marshalEventEnvelope(event)
	case target.InputTransformer != nil:
		return s.applyInputTransform(ruleARN, ruleName, event, target.InputTransformer)
	default:
		return marshalEventEnvelope(event)
	}
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
