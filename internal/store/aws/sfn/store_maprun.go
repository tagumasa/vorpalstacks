// This file holds the Distributed Map run records: the run sequence
// counter and the run CRUD the Map Run APIs and the engine reclaim read.

// Package stepfunction provides Step Functions storage functionality for vorpalstacks.
package sfn

import (
	"context"
	"strconv"
	"strings"
	"sync/atomic"

	"vorpalstacks/internal/store/aws/common"
)

func (s *StepFunctionStore) nextMapRunSeq() int64 {
	if atomic.LoadInt64(&s.mapRunSeq) == 0 {
		s.recoverMapRunSeq()
	}
	return atomic.AddInt64(&s.mapRunSeq, 1)
}

// NextMapRunSeq returns the next sequential identifier for a map run.
func (s *StepFunctionStore) NextMapRunSeq() int64 {
	return s.nextMapRunSeq()
}

// MapRunSeqFromARN extracts the monotonic sequence number a map run ARN's
// run identifier carries (mapRun:<stateMachine>/<label>/mapRun-<seq>-<timestamp>).
// The second return is false when the ARN carries no parsable sequence —
// callers fall back to other ordering evidence rather than treating a
// zero as a sequence.
func MapRunSeqFromARN(arn string) (int64, bool) {
	idx := strings.LastIndex(arn, "/mapRun-")
	if idx < 0 {
		return 0, false
	}
	rest := arn[idx+len("/mapRun-"):]
	end := 0
	for end < len(rest) && rest[end] >= '0' && rest[end] <= '9' {
		end++
	}
	if end == 0 {
		return 0, false
	}
	n, err := strconv.ParseInt(rest[:end], 10, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (s *StepFunctionStore) recoverMapRunSeq() {
	mapRuns, err := common.ListMatching[MapRun](s.mapRunsStore, "", nil)
	if err != nil {
		return
	}
	var maxSeq int64
	for _, mr := range mapRuns {
		if n, ok := MapRunSeqFromARN(mr.MapRunArn); ok && n > maxSeq {
			maxSeq = n
		}
	}
	if maxSeq > s.mapRunSeq {
		atomic.StoreInt64(&s.mapRunSeq, maxSeq)
	}
}

// CreateMapRun persists a new map run in Pebble-backed storage.
func (s *StepFunctionStore) CreateMapRun(ctx context.Context, mr *MapRun) error {
	return s.mapRunsStore.Put(mr.MapRunArn, mr)
}

// UpdateMapRun persists an updated map run record.
func (s *StepFunctionStore) UpdateMapRun(ctx context.Context, mr *MapRun) error {
	return s.mapRunsStore.Put(mr.MapRunArn, mr)
}

// GetMapRun retrieves a map run by its ARN. A missing record is the
// not-found sentinel; any other storage fault surfaces as itself.
func (s *StepFunctionStore) GetMapRun(ctx context.Context, mapRunArn string) (*MapRun, error) {
	var mr MapRun
	if err := s.mapRunsStore.Get(mapRunArn, &mr); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrMapRunNotFound
		}
		return nil, err
	}
	return &mr, nil
}

// ListMapRunsByExecution returns all map runs for a given execution ARN.
func (s *StepFunctionStore) ListMapRunsByExecution(ctx context.Context, executionArn string) ([]*MapRun, error) {
	return common.ListMatching[MapRun](s.mapRunsStore, "", func(mr *MapRun) bool {
		return mr.ExecutionArn == executionArn
	})
}

// ListAllMapRuns returns all map runs, optionally filtered by execution
// ARN. Inline-map checkpoint records are excluded here — before
// pagination — because they are not Map Run resources on the API surface:
// a post-listing filter would let a page of inline records satisfy the
// page size while returning no visible items.
func (s *StepFunctionStore) ListAllMapRuns(ctx context.Context, executionArn string, limit int32, nextToken string) (*MapRunListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[MapRun](s.mapRunsStore, opts, func(mr *MapRun) bool {
		if mr.Inline {
			return false
		}
		if executionArn != "" && mr.ExecutionArn != executionArn {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return &MapRunListResult{
		MapRuns:   result.Items,
		NextToken: result.NextMarker,
	}, nil
}
