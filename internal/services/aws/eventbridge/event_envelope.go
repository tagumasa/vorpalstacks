package eventbridge

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/common/request"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// marshalEventEnvelope renders the envelope as the payload bytes a target
// receives when no input configuration narrows the event.
func marshalEventEnvelope(event *eventsstore.Event) []byte {
	b, err := json.Marshal(eventEnvelope(event))
	if err != nil {
		return []byte("{}")
	}
	return b
}

// eventEnvelope renders the AWS event wire form: the nine envelope members
// under one time rendering, plus the replay-name metadata field on replayed
// events. The rendering is shared by target payloads, archives and pattern
// matching so the three cannot drift, and RFC3339Nano keeps sub-second
// timestamps through the archive and replay round-trip. The X-Ray trace
// header is deliberately not an envelope member: AWS propagates it on the
// transport (the SQS message attribute for dead-letter queues, the
// X-Amzn-Trace-Id header for invocations), never on the event itself — the
// EventBridge X-Ray integration states the trace header is not available on
// the event delivered to the target, nor on archived and replayed events.
// formatEventTime renders an event timestamp for the wire envelope — the
// single definition of the event-time format, shared by the envelope codec
// and the transformer's ingestion-time reserved variable.
func formatEventTime(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func eventEnvelope(event *eventsstore.Event) map[string]interface{} {
	envelope := map[string]interface{}{
		"version":     event.Version,
		"id":          event.ID,
		"detail-type": event.DetailType,
		"source":      event.Source,
		"account":     event.Account,
		"time":        formatEventTime(event.Time),
		"region":      event.Region,
		"resources":   event.Resources,
		"detail":      event.Detail,
	}
	// EventBridge adds the replay-name metadata field to events it replays
	// (the archives user guide); live events never carry it, so the member
	// is rendered only when the replay path stamped it — which is also what
	// makes the documented {"replay-name":[{"exists":false}]} pattern form
	// discriminate replays from live events under pattern matching.
	if event.ReplayName != "" {
		envelope["replay-name"] = event.ReplayName
	}
	return envelope
}

// eventFromEnvelope is eventEnvelope's inverse: it rebuilds the Event from a
// serialised envelope, wherever the envelope was produced (bus-to-bus
// delivery input, an archived event). eventBusName names the bus the event is
// being delivered onto — a property of the delivery, not of the envelope —
// so the caller supplies it.
func eventFromEnvelope(envelope map[string]interface{}, eventBusName string) *eventsstore.Event {
	event := &eventsstore.Event{
		ID:           request.GetStringParam(envelope, "id"),
		Version:      request.GetStringParam(envelope, "version"),
		DetailType:   request.GetStringParam(envelope, "detail-type"),
		Source:       request.GetStringParam(envelope, "source"),
		Account:      request.GetStringParam(envelope, "account"),
		Region:       request.GetStringParam(envelope, "region"),
		EventBusName: eventBusName,
		// The replay-name stamp survives the envelope hop so a replayed
		// event forwarded through a bus target stays recognisable — and
		// stays excluded from archives — on the downstream bus.
		ReplayName: request.GetStringParam(envelope, "replay-name"),
	}

	if timeStr := request.GetStringParam(envelope, "time"); timeStr != "" {
		if parsed, err := time.Parse(time.RFC3339, timeStr); err == nil {
			event.Time = parsed
		}
	}

	if detail, ok := envelope["detail"].(map[string]interface{}); ok {
		event.Detail = detail
	}

	// Resources arrive as []interface{} across a JSON hop and stay []string
	// on an in-process round-trip; the codec accepts both.
	switch resources := envelope["resources"].(type) {
	case []interface{}:
		for _, r := range resources {
			if rStr, ok := r.(string); ok {
				event.Resources = append(event.Resources, rStr)
			}
		}
	case []string:
		event.Resources = append(event.Resources, resources...)
	}

	return event
}
