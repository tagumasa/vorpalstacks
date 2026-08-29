package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Kinesis Destination Core — single validation + persistence path for the
// Kinesis streaming destination operations.
//
// Both the HTTP API handlers (kinesis_operations.go) and any future admin
// handler delegate to these methods to ensure identical behaviour.
// ---------------------------------------------------------------------------

// Kinesis streaming destination states and the
// ApproximateCreationDateTimePrecision enum values from the model.
const (
	kinesisDestinationEnabling  = "ENABLING"
	kinesisDestinationActive    = "ACTIVE"
	kinesisDestinationDisabling = "DISABLING"
	kinesisDestinationUpdating  = "UPDATING"

	kinesisPrecisionMillisecond = "MILLISECOND"
	kinesisPrecisionMicrosecond = "MICROSECOND"

	// kinesisDestinationTransitionDelay is the delay before a destination
	// reaches the state its triggering request announced (ENABLING/DISABLING/
	// UPDATING are observable intermediate states, mirroring the
	// asynchronous service-side workflow).
	kinesisDestinationTransitionDelay = 1 * time.Second
)

// kinesisPrecisionFromConfig reads the ApproximateCreationDateTimePrecision
// enum out of the named configuration member. An empty string means the
// member was absent; a value outside the model enum is a client error.
func kinesisPrecisionFromConfig(params map[string]interface{}, member string) (string, error) {
	configMap, ok := params[member].(map[string]interface{})
	if !ok {
		return "", nil
	}
	value, _ := configMap["ApproximateCreationDateTimePrecision"].(string)
	if value == "" {
		return "", nil
	}
	if value != kinesisPrecisionMillisecond && value != kinesisPrecisionMicrosecond {
		return "", ErrInvalidParameter
	}
	return value, nil
}

// kinesisStreamExists reports whether the destination stream exists in the
// table's region. A DynamoDB table may only stream to a Kinesis data stream
// in the same account and region, so anything else is rejected as a
// nonexistent resource.
func (s *DynamoDBService) kinesisStreamExists(reqCtx *request.RequestContext, streamArn string) bool {
	if s.bus == nil {
		return true
	}
	invoker := s.bus.KinesisInvoker()
	if invoker == nil {
		return true
	}
	exists, err := invoker.StreamExists(reqCtx.Context, reqCtx.GetRegion(), streamArn)
	if err != nil {
		logs.Warn("Kinesis destination existence check failed; allowing request",
			logs.Err(err),
			logs.String("streamArn", streamArn),
		)
		return true
	}
	return exists
}

// copyKinesisDestinations deep-copies a destination slice so handlers and
// background transitions never mutate a slice another goroutine holds.
func copyKinesisDestinations(destinations []*dbstore.KinesisDataStreamDestination) []*dbstore.KinesisDataStreamDestination {
	copied := make([]*dbstore.KinesisDataStreamDestination, len(destinations))
	for i, d := range destinations {
		entry := *d
		copied[i] = &entry
	}
	return copied
}

// kinesisDestinationApply decides what a delayed transition does with the
// destination entry it finds. ok=false leaves the entry untouched (a newer
// request changed its state in the meantime); remove=true drops the entry;
// otherwise the function has mutated the entry in place and the result must
// be persisted.
type kinesisDestinationApply func(d *dbstore.KinesisDataStreamDestination) (remove bool, ok bool)

// scheduleKinesisDestinationTransition performs a destination transition
// after the observable intermediate state has been visible for a moment. The
// transition is conditional: the table is re-read under kinesisDestMu after
// the delay, and the change applies only when the entry still carries the
// state the triggering request left it in. A stale transition can therefore
// never resurrect a destination that a newer Enable/Disable/Update has since
// changed or removed.
func (s *DynamoDBService) scheduleKinesisDestinationTransition(store dbstore.DynamoDBStoreInterface, tableName, streamArn string, apply kinesisDestinationApply) {
	s.bgWg.Add(1)
	go func() {
		defer func() { resilience.RecoverPanic("dynamodb Kinesis destination transition") }()
		defer s.bgWg.Done()
		select {
		case <-time.After(kinesisDestinationTransitionDelay):
		case <-s.bgCtx.Done():
			return
		}

		s.kinesisDestMu.Lock()
		defer s.kinesisDestMu.Unlock()

		table, err := store.Tables().Get(tableName)
		if err != nil || table == nil {
			return
		}
		for i, d := range table.KinesisDataStreamDestinations {
			if d.StreamArn != streamArn {
				continue
			}
			remove, ok := apply(d)
			if !ok {
				return
			}
			next := table.KinesisDataStreamDestinations
			if remove {
				next = make([]*dbstore.KinesisDataStreamDestination, 0, len(table.KinesisDataStreamDestinations)-1)
				next = append(next, table.KinesisDataStreamDestinations[:i]...)
				next = append(next, table.KinesisDataStreamDestinations[i+1:]...)
			}
			if err := store.Tables().SetKinesisStreamingDestination(tableName, next); err != nil {
				logs.Error("Failed to complete Kinesis destination transition",
					logs.Err(err),
					logs.String("tableName", tableName),
					logs.String("streamArn", streamArn),
				)
			}
			return
		}
	}()
}

// enableKinesisStreamingDestinationInput carries the raw wire parameters
// for EnableKinesisStreamingDestination.
type enableKinesisStreamingDestinationInput struct {
	Parameters map[string]interface{}
}

// enableKinesisStreamingDestinationCore validates the request, then enables
// Kinesis streaming from the named table to the requested stream.
func (s *DynamoDBService) enableKinesisStreamingDestinationCore(ctx context.Context, reqCtx *request.RequestContext, in enableKinesisStreamingDestinationInput) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, in.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	streamArn := request.GetStringParam(in.Parameters, "StreamArn")
	if !validateStreamArn(streamArn) {
		return nil, ErrInvalidParameter
	}

	precision, err := kinesisPrecisionFromConfig(in.Parameters, "EnableKinesisStreamingConfiguration")
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if !s.kinesisStreamExists(reqCtx, streamArn) {
		return nil, ErrResourceNotFound
	}

	s.kinesisDestMu.Lock()
	defer s.kinesisDestMu.Unlock()

	// Re-read under the lock: another request's transition may have changed
	// the destination list since the initial fetch.
	current, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrTableNotFound
	}
	// A table streams to at most one Kinesis data stream; an entry in any
	// state (including a transitional one) occupies that slot.
	if len(current.KinesisDataStreamDestinations) > 0 {
		return nil, ErrResourceAlreadyExists
	}

	destinations := []*dbstore.KinesisDataStreamDestination{{
		StreamArn:                            streamArn,
		DestinationStatus:                    kinesisDestinationEnabling,
		ApproximateCreationDateTimePrecision: precision,
	}}
	if err := store.Tables().SetKinesisStreamingDestination(tableName, destinations); err != nil {
		return nil, err
	}

	s.scheduleKinesisDestinationTransition(store, tableName, streamArn, func(d *dbstore.KinesisDataStreamDestination) (bool, bool) {
		if d.DestinationStatus != kinesisDestinationEnabling {
			return false, false
		}
		d.DestinationStatus = kinesisDestinationActive
		return false, true
	})

	response := map[string]interface{}{
		"TableName":         tableName,
		"StreamArn":         streamArn,
		"DestinationStatus": kinesisDestinationEnabling,
	}
	if precision != "" {
		response["EnableKinesisStreamingConfiguration"] = map[string]interface{}{
			"ApproximateCreationDateTimePrecision": precision,
		}
	}
	return response, nil
}

// disableKinesisStreamingDestinationInput carries the raw wire parameters
// for DisableKinesisStreamingDestination.
type disableKinesisStreamingDestinationInput struct {
	Parameters map[string]interface{}
}

// disableKinesisStreamingDestinationCore validates the request, then
// disables Kinesis streaming from the named table.
func (s *DynamoDBService) disableKinesisStreamingDestinationCore(ctx context.Context, reqCtx *request.RequestContext, in disableKinesisStreamingDestinationInput) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, in.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	streamArn := request.GetStringParam(in.Parameters, "StreamArn")
	if !validateStreamArn(streamArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	s.kinesisDestMu.Lock()
	defer s.kinesisDestMu.Unlock()

	// Re-read under the lock so the removal decision is made against the
	// current destination list.
	current, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrTableNotFound
	}
	idx := -1
	for i, d := range current.KinesisDataStreamDestinations {
		if d.StreamArn == streamArn {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil, ErrResourceNotFound
	}

	destinations := copyKinesisDestinations(current.KinesisDataStreamDestinations)
	destinations[idx].DestinationStatus = kinesisDestinationDisabling
	if err := store.Tables().SetKinesisStreamingDestination(tableName, destinations); err != nil {
		return nil, err
	}

	s.scheduleKinesisDestinationTransition(store, tableName, streamArn, func(d *dbstore.KinesisDataStreamDestination) (bool, bool) {
		if d.DestinationStatus != kinesisDestinationDisabling {
			return false, false
		}
		return true, true
	})

	return map[string]interface{}{
		"TableName":         tableName,
		"StreamArn":         streamArn,
		"DestinationStatus": kinesisDestinationDisabling,
	}, nil
}

// updateKinesisStreamingDestinationInput carries the raw wire parameters
// for UpdateKinesisStreamingDestination.
type updateKinesisStreamingDestinationInput struct {
	Parameters map[string]interface{}
}

// updateKinesisStreamingDestinationCore validates the request, then updates
// the Kinesis streaming destination of the named table.
func (s *DynamoDBService) updateKinesisStreamingDestinationCore(ctx context.Context, reqCtx *request.RequestContext, in updateKinesisStreamingDestinationInput) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, in.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	streamArn := request.GetStringParam(in.Parameters, "StreamArn")
	if !validateStreamArn(streamArn) {
		return nil, ErrInvalidParameter
	}

	precision, err := kinesisPrecisionFromConfig(in.Parameters, "UpdateKinesisStreamingConfiguration")
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	s.kinesisDestMu.Lock()
	defer s.kinesisDestMu.Unlock()

	// Re-read under the lock: precision applies to the destination as it
	// exists now, and the state transition must not race another request.
	current, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrTableNotFound
	}
	idx := -1
	for i, d := range current.KinesisDataStreamDestinations {
		if d.StreamArn == streamArn {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil, ErrResourceNotFound
	}

	destinations := copyKinesisDestinations(current.KinesisDataStreamDestinations)
	if precision != "" {
		destinations[idx].ApproximateCreationDateTimePrecision = precision
	}
	destinations[idx].DestinationStatus = kinesisDestinationUpdating
	if err := store.Tables().SetKinesisStreamingDestination(tableName, destinations); err != nil {
		return nil, err
	}

	s.scheduleKinesisDestinationTransition(store, tableName, streamArn, func(d *dbstore.KinesisDataStreamDestination) (bool, bool) {
		if d.DestinationStatus != kinesisDestinationUpdating {
			return false, false
		}
		d.DestinationStatus = kinesisDestinationActive
		return false, true
	})

	response := map[string]interface{}{
		"TableName":         tableName,
		"StreamArn":         streamArn,
		"DestinationStatus": kinesisDestinationUpdating,
	}
	if precision != "" {
		response["UpdateKinesisStreamingConfiguration"] = map[string]interface{}{
			"ApproximateCreationDateTimePrecision": precision,
		}
	}
	return response, nil
}
