package neptunedata

import (
	"fmt"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage/graphengine"
)

// FastResetInput carries the parsed ExecuteFastReset wire payload plus the
// request-context state the reset Core needs. Graph holds the cluster's graph
// store (extracted from the request context graph writer by the handler); it
// is nil when the request carries no graph context.
type FastResetInput struct {
	Action string
	Token  string
	Region string
	Graph  graphengine.GraphStore
}

// executeFastResetCore implements the two-phase database reset protocol. The
// initiateDatabaseReset action issues a time-limited token; performDatabaseReset
// validates the token and clears all graph data.
func (s *NeptuneDataService) executeFastResetCore(in *FastResetInput) (map[string]interface{}, error) {
	switch in.Action {
	case "initiateDatabaseReset":
		token := generateFastResetToken()
		now := time.Now()
		s.fastTokens.Store(token, now.Add(30*time.Second))
		s.fastTokens.Range(func(key, value any) bool {
			if t, ok := value.(time.Time); ok {
				if now.After(t) {
					s.fastTokens.Delete(key)
				}
			}
			return true
		})
		return map[string]interface{}{
			"payload": map[string]interface{}{
				"token": token,
			},
		}, nil

	case "performDatabaseReset":
		if in.Token == "" {
			return nil, missingParameter("token")
		}
		val, ok := s.fastTokens.Load(in.Token)
		if !ok {
			return nil, preconditionFailed("invalid or expired token")
		}
		expiry, typeOk := val.(time.Time)
		if !typeOk || time.Now().After(expiry) {
			s.fastTokens.Delete(in.Token)
			return nil, preconditionFailed("invalid or expired token")
		}
		s.fastTokens.Delete(in.Token)

		if in.Graph != nil {
			if err := in.Graph.Clear(); err != nil {
				logs.Error("failed to clear graph store during fast reset", logs.Err(err))
				return nil, internalFailure(fmt.Sprintf("fast reset failed: graph clear error: %v", err))
			}
		}
		region := in.Region
		s.mu.Lock()
		s.statsDisabled = false
		s.autoComputeEnabled = true
		s.mu.Unlock()
		s.statsMap.Store(region, &GraphStatistics{
			LabelCounts: make(map[string]int64),
			RelCounts:   make(map[string]int64),
		})
		if resetStore, err := s.GetStoreForRegion(region); err == nil {
			queries, _ := resetStore.ListQueries()
			for _, q := range queries {
				if err := resetStore.DeleteQuery(q.GetQueryId()); err != nil {
					logs.Warn("failed to delete query during reset", logs.String("queryId", q.GetQueryId()), logs.Err(err))
				}
			}
			jobs, _ := resetStore.ListLoaderJobs()
			for _, j := range jobs {
				if err := resetStore.DeleteLoaderJob(j.GetLoadId()); err != nil {
					logs.Warn("failed to delete loader job during reset", logs.String("loadId", j.GetLoadId()), logs.Err(err))
				}
			}
		}
		return map[string]interface{}{
			"status": "200 OK",
		}, nil

	default:
		return nil, invalidParameter(fmt.Sprintf("unknown action: %s", in.Action))
	}
}
