// Package athena provides AWS Athena storage functionality for vorpalstacks.
package athena

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_athena"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func workGroupBucketName(region string) string         { return "athena_workgroups-" + region }
func namedQueryBucketName(region string) string        { return "athena_named_queries-" + region }
func preparedStatementBucketName(region string) string { return "athena_prepared_statements-" + region }
func queryExecutionBucketName(region string) string    { return "athena_query_executions-" + region }
func resultBucketName(region string) string            { return "athena_results-" + region }
func dataCatalogBucketName(region string) string       { return "athena_data_catalogs-" + region }
func databaseBucketName(region string) string          { return "athena_databases-" + region }
func tableBucketName(region string) string             { return "athena_tables-" + region }
func tableDataBucketName(region string) string         { return "athena_table_data-" + region }

// Store provides Athena storage operations.
type Store struct {
	*common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
}

// NewStore creates a new Athena store.
func NewStore(store storage.BasicStorage, accountID, region string) *Store {
	return &Store{
		BaseStore:  common.NewBaseStore(store.Bucket(workGroupBucketName(region)), "athena"),
		TagStore:   common.NewTagStoreWithRegion(store, "athena", region),
		arnBuilder: svcarn.NewARNBuilder(accountID, region),
		accountID:  accountID,
		region:     region,
	}
}

// GetAccountID returns the AWS account ID associated with the store.
func (s *Store) GetAccountID() string { return s.accountID }

// GetRegion returns the AWS region associated with the store.
func (s *Store) GetRegion() string { return s.region }

// WorkGroupStore provides Athena work group storage operations.
type WorkGroupStore struct {
	*common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	kl         common.KeyLocker
}

// NewWorkGroupStore creates a new Athena work group store.
func NewWorkGroupStore(store storage.BasicStorage, accountID, region string) *WorkGroupStore {
	s := &WorkGroupStore{
		BaseStore:  common.NewBaseStore(store.Bucket(workGroupBucketName(region)), "athena-workgroup"),
		TagStore:   common.NewTagStoreWithRegion(store, "athena-workgroup", region),
		arnBuilder: svcarn.NewARNBuilder(accountID, region),
		accountID:  accountID,
		region:     region,
	}
	s.ensurePrimaryWorkGroup()
	return s
}

func (s *WorkGroupStore) ensurePrimaryWorkGroup() {
	s.kl.WithLock("primary", func() error {
		if !s.Exists("primary") {
			primaryWg := &WorkGroup{
				Name:        "primary",
				Description: "The primary workgroup",
				State:       WorkGroupStateEnabled,
				CreatedTime: time.Now().UTC(),
			}
			if err := s.PutProto("primary", WorkGroupToProto(primaryWg)); err != nil {
				logs.Error("Failed to create primary workgroup", logs.Err(err))
			}
		}
		return nil
	})
}

// CreateWorkGroup creates a new Athena work group.
func (s *WorkGroupStore) CreateWorkGroup(wg *WorkGroup) error {
	return s.kl.WithLock(wg.Name, func() error {
		if s.Exists(wg.Name) {
			return ErrWorkGroupAlreadyExists
		}

		wg.CreatedTime = time.Now().UTC()
		if wg.State == "" {
			wg.State = WorkGroupStateEnabled
		}

		return s.PutProto(wg.Name, WorkGroupToProto(wg))
	})
}

// GetWorkGroup retrieves an Athena work group by name.
func (s *WorkGroupStore) GetWorkGroup(name string) (*WorkGroup, error) {
	var p pb.WorkGroup
	if err := s.GetProto(name, &p); err != nil {
		return nil, ErrWorkGroupNotFound
	}
	return ProtoToWorkGroup(&p), nil
}

// UpdateWorkGroup updates an existing Athena work group.
func (s *WorkGroupStore) UpdateWorkGroup(wg *WorkGroup) error {
	return s.kl.WithLock(wg.Name, func() error {
		if !s.Exists(wg.Name) {
			return ErrWorkGroupNotFound
		}
		return s.PutProto(wg.Name, WorkGroupToProto(wg))
	})
}

// DeleteWorkGroup deletes an Athena work group by name.
func (s *WorkGroupStore) DeleteWorkGroup(name string) error {
	return s.kl.WithLock(name, func() error {
		if !s.Exists(name) {
			return ErrWorkGroupNotFound
		}
		arn := s.GetARN(name)
		_ = s.TagStore.Delete(arn)
		return s.BaseStore.Delete(name)
	})
}

// ListWorkGroups returns Athena work groups with pagination support.
func (s *WorkGroupStore) ListWorkGroups(opts common.ListOptions) (*common.ListResult[WorkGroup], error) {
	result, err := common.ListProto[*pb.WorkGroup](s.BaseStore, opts, func() *pb.WorkGroup { return &pb.WorkGroup{} }, nil)
	if err != nil {
		return nil, err
	}
	workGroups := make([]*WorkGroup, len(result.Items))
	for i, p := range result.Items {
		workGroups[i] = ProtoToWorkGroup(p)
	}
	return &common.ListResult[WorkGroup]{
		Items:       workGroups,
		NextMarker:  result.NextMarker,
		IsTruncated: result.IsTruncated,
	}, nil
}

// GetARN returns the ARN for an Athena work group.
func (s *WorkGroupStore) GetARN(name string) string {
	return s.arnBuilder.Athena().WorkGroup(name)
}

// NamedQueryStore provides Athena named query storage operations.
type NamedQueryStore struct {
	*common.BaseStore
	kl common.KeyLocker
}

// NewNamedQueryStore creates a new Athena named query store.
func NewNamedQueryStore(store storage.BasicStorage, region string) *NamedQueryStore {
	return &NamedQueryStore{
		BaseStore: common.NewBaseStore(store.Bucket(namedQueryBucketName(region)), "athena-named-query"),
	}
}

func (s *NamedQueryStore) namedQueryByNameKey(workGroup, name string) string {
	return "#name:" + workGroup + ":" + name
}

// CreateNamedQuery creates a new Athena named query.
func (s *NamedQueryStore) CreateNamedQuery(nq *NamedQuery) error {
	if nq.NamedQueryId == "" {
		nq.NamedQueryId = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return s.kl.WithLock(nq.NamedQueryId, func() error {
		nq.CreatedTime = time.Now().UTC()

		nameKey := s.namedQueryByNameKey(nq.WorkGroup, nq.Name)

		if s.Exists(nameKey) {
			var existingIdBytes []byte
			err := s.Get(nameKey, &existingIdBytes)
			existingId := string(existingIdBytes)
			if err != nil || existingId == "" {
				// Stale name index from a crashed update — overwrite it
			} else if s.Exists(existingId) {
				return ErrNamedQueryAlreadyExists
			}
			// Stale index pointing to deleted query — overwrite
		}

		if err := s.Put(nameKey, []byte(nq.NamedQueryId)); err != nil {
			return err
		}

		if err := s.PutProto(nq.NamedQueryId, NamedQueryToProto(nq)); err != nil {
			if delErr := s.BaseStore.Delete(nameKey); delErr != nil {
				logs.Warn("NamedQuery: failed to rollback name index on PutProto error",
					logs.String("nameKey", nameKey), logs.Err(delErr))
			}
			return err
		}

		return nil
	})
}

// GetNamedQuery retrieves an Athena named query by ID.
func (s *NamedQueryStore) GetNamedQuery(id string) (*NamedQuery, error) {
	var p pb.NamedQuery
	if err := s.GetProto(id, &p); err != nil {
		return nil, ErrNamedQueryNotFound
	}
	return ProtoToNamedQuery(&p), nil
}

// UpdateNamedQuery updates an existing named query.
func (s *NamedQueryStore) UpdateNamedQuery(id string, nq *NamedQuery) error {
	return s.kl.WithLock(id, func() error {
		if !s.Exists(id) {
			return ErrNamedQueryNotFound
		}

		oldNq, err := s.GetNamedQuery(id)
		if err != nil {
			return err
		}

		newNameKey := s.namedQueryByNameKey(nq.WorkGroup, nq.Name)
		if err := s.Put(newNameKey, []byte(nq.NamedQueryId)); err != nil {
			return err
		}

		if err := s.PutProto(id, NamedQueryToProto(nq)); err != nil {
			if delErr := s.BaseStore.Delete(newNameKey); delErr != nil {
				logs.Warn("NamedQuery: failed to rollback new name index on PutProto error",
					logs.String("newNameKey", newNameKey), logs.Err(delErr))
			}
			return err
		}

		if oldNq.Name != nq.Name || oldNq.WorkGroup != nq.WorkGroup {
			oldNameKey := s.namedQueryByNameKey(oldNq.WorkGroup, oldNq.Name)
			if err := s.BaseStore.Delete(oldNameKey); err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteNamedQuery deletes an Athena named query by ID.
func (s *NamedQueryStore) DeleteNamedQuery(id string) error {
	return s.kl.WithLock(id, func() error {
		nq, err := s.GetNamedQuery(id)
		if err != nil {
			return err
		}

		nameKey := s.namedQueryByNameKey(nq.WorkGroup, nq.Name)
		if err := s.BaseStore.Delete(nameKey); err != nil {
			return err
		}

		if err := s.BaseStore.Delete(id); err != nil {
			if putErr := s.Put(nameKey, []byte(nq.NamedQueryId)); putErr != nil {
				logs.Warn("NamedQuery: failed to restore name index on Delete error",
					logs.String("nameKey", nameKey), logs.Err(putErr))
			}
			return err
		}

		return nil
	})
}

// ListNamedQueries returns all Athena named queries for a work group.
func (s *NamedQueryStore) ListNamedQueries(workGroup string) ([]*NamedQuery, error) {
	var namedQueries []*NamedQuery

	err := s.ForEach(func(key string, value []byte) error {
		if len(key) > 0 && key[0] == '#' {
			return nil
		}
		if len(key) > 0 && key[0] == '_' {
			return nil
		}
		var p pb.NamedQuery
		if err := proto.Unmarshal(value, &p); err != nil {
			return nil
		}
		nq := ProtoToNamedQuery(&p)
		if workGroup == "" || nq.WorkGroup == workGroup {
			namedQueries = append(namedQueries, nq)
		}
		return nil
	})

	return namedQueries, err
}

// DeleteNamedQueriesByWorkGroup removes all named queries belonging to the given workgroup.
func (s *NamedQueryStore) DeleteNamedQueriesByWorkGroup(workGroup string) error {
	queries, err := s.ListNamedQueries(workGroup)
	if err != nil {
		return err
	}
	for _, nq := range queries {
		nameKey := s.namedQueryByNameKey(nq.WorkGroup, nq.Name)
		_ = s.BaseStore.Delete(nameKey)
		_ = s.BaseStore.Delete(nq.NamedQueryId)
	}
	return nil
}

// PreparedStatementStore provides Athena prepared statement storage operations.
type PreparedStatementStore struct {
	*common.BaseStore
	kl common.KeyLocker
}

// NewPreparedStatementStore creates a new Athena prepared statement store.
func NewPreparedStatementStore(store storage.BasicStorage, region string) *PreparedStatementStore {
	return &PreparedStatementStore{
		BaseStore: common.NewBaseStore(store.Bucket(preparedStatementBucketName(region)), "athena-prepared-statement"),
	}
}

func (s *PreparedStatementStore) preparedStatementKey(workGroup, name string) string {
	return workGroup + ":" + name
}

// CreatePreparedStatement creates a new Athena prepared statement.
func (s *PreparedStatementStore) CreatePreparedStatement(ps *PreparedStatement) error {
	key := s.preparedStatementKey(ps.WorkGroupName, ps.StatementName)

	if s.Exists(key) {
		return ErrPreparedStatementAlreadyExists
	}

	ps.CreatedTime = time.Now().UTC()
	ps.LastModifiedTime = ps.CreatedTime

	return s.PutProto(key, PreparedStatementToProto(ps))
}

// GetPreparedStatement retrieves an Athena prepared statement.
func (s *PreparedStatementStore) GetPreparedStatement(workGroup, name string) (*PreparedStatement, error) {
	key := s.preparedStatementKey(workGroup, name)

	var p pb.PreparedStatement
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrPreparedStatementNotFound
	}
	return ProtoToPreparedStatement(&p), nil
}

// UpdatePreparedStatement updates an existing Athena prepared statement.
func (s *PreparedStatementStore) UpdatePreparedStatement(ps *PreparedStatement) error {
	key := s.preparedStatementKey(ps.WorkGroupName, ps.StatementName)

	return s.kl.WithLock(key, func() error {
		existing, err := s.GetPreparedStatement(ps.WorkGroupName, ps.StatementName)
		if err != nil {
			return err
		}

		ps.CreatedTime = existing.CreatedTime
		ps.LastModifiedTime = time.Now().UTC()

		return s.PutProto(key, PreparedStatementToProto(ps))
	})
}

// DeletePreparedStatement deletes an Athena prepared statement.
func (s *PreparedStatementStore) DeletePreparedStatement(workGroup, name string) error {
	key := s.preparedStatementKey(workGroup, name)
	return s.kl.WithLock(key, func() error {
		if !s.Exists(key) {
			return ErrPreparedStatementNotFound
		}
		return s.BaseStore.Delete(key)
	})
}

// ListPreparedStatements returns all Athena prepared statements for a work group.
func (s *PreparedStatementStore) ListPreparedStatements(workGroup string) ([]*PreparedStatement, error) {
	prefix := ""
	if workGroup != "" {
		prefix = workGroup + ":"
	}

	result, err := common.ListProto[*pb.PreparedStatement](s.BaseStore, common.ListOptions{Prefix: prefix}, func() *pb.PreparedStatement { return &pb.PreparedStatement{} }, nil)
	if err != nil {
		return nil, err
	}
	statements := make([]*PreparedStatement, len(result.Items))
	for i, p := range result.Items {
		statements[i] = ProtoToPreparedStatement(p)
	}
	return statements, nil
}

// DeletePreparedStatementsByWorkGroup removes all prepared statements belonging to the given workgroup.
func (s *PreparedStatementStore) DeletePreparedStatementsByWorkGroup(workGroup string) error {
	prefix := workGroup + ":"
	result, err := common.ListProto[*pb.PreparedStatement](s.BaseStore, common.ListOptions{Prefix: prefix}, func() *pb.PreparedStatement { return &pb.PreparedStatement{} }, nil)
	if err != nil {
		return err
	}
	for _, p := range result.Items {
		key := s.preparedStatementKey(p.WorkGroupName, p.StatementName)
		_ = s.BaseStore.Delete(key)
	}
	return nil
}

// QueryExecutionStore provides Athena query execution storage operations.
type QueryExecutionStore struct {
	*common.BaseStore
}

// NewQueryExecutionStore creates a new Athena query execution store.
func NewQueryExecutionStore(store storage.BasicStorage, region string) *QueryExecutionStore {
	return &QueryExecutionStore{
		BaseStore: common.NewBaseStore(store.Bucket(queryExecutionBucketName(region)), "athena-query-execution"),
	}
}

func (s *QueryExecutionStore) workGroupIndexKey(workGroup, queryExecutionId string) string {
	return "#wg:" + workGroup + ":" + queryExecutionId
}

// extractIDFromIndexKey parses "#wg:<workGroup>:<id>" and returns the id portion.
func extractIDFromIndexKey(key string) string {
	parts := strings.SplitN(key, ":", 3)
	if len(parts) == 3 {
		return parts[2]
	}
	return ""
}

// CreateQueryExecution creates a new Athena query execution.
func (s *QueryExecutionStore) CreateQueryExecution(qe *QueryExecution) error {
	if qe.QueryExecutionId == "" {
		qe.QueryExecutionId = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	if qe.Status == nil {
		qe.Status = &QueryExecutionStatus{
			State:              QueryExecutionStateQueued,
			SubmissionDateTime: time.Now().UTC(),
		}
	}

	if err := s.PutProto(qe.QueryExecutionId, QueryExecutionToProto(qe)); err != nil {
		return err
	}

	if qe.WorkGroup != "" {
		if err := s.Put(s.workGroupIndexKey(qe.WorkGroup, qe.QueryExecutionId), nil); err != nil {
			_ = s.BaseStore.Delete(qe.QueryExecutionId)
			return fmt.Errorf("failed to write workgroup index: %w", err)
		}
	}

	return nil
}

// GetQueryExecution retrieves an Athena query execution by ID.
func (s *QueryExecutionStore) GetQueryExecution(id string) (*QueryExecution, error) {
	var p pb.QueryExecution
	if err := s.GetProto(id, &p); err != nil {
		return nil, ErrQueryExecutionNotFound
	}
	return ProtoToQueryExecution(&p), nil
}

// UpdateQueryExecution updates an existing Athena query execution.
func (s *QueryExecutionStore) UpdateQueryExecution(qe *QueryExecution) error {
	return s.PutProto(qe.QueryExecutionId, QueryExecutionToProto(qe))
}

// DeleteExpiredQueryExecutions removes query executions older than the given threshold.
// Returns the number of deleted executions and their IDs for cascade cleanup.
func (s *QueryExecutionStore) DeleteExpiredQueryExecutions(olderThan time.Time) (int, []string, error) {
	var ids []string
	err := s.ForEach(func(key string, value []byte) error {
		if len(key) > 0 && key[0] == '#' {
			return nil
		}
		var p pb.QueryExecution
		if err := proto.Unmarshal(value, &p); err != nil {
			return nil
		}
		if p.Status == nil || p.Status.SubmissionDateTime == nil {
			return nil
		}
		if !p.Status.SubmissionDateTime.AsTime().Before(olderThan) {
			return nil
		}
		if p.WorkGroup != "" {
			_ = s.BaseStore.Delete(s.workGroupIndexKey(p.WorkGroup, p.QueryExecutionId))
		}
		ids = append(ids, p.QueryExecutionId)
		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	for _, id := range ids {
		_ = s.BaseStore.Delete(id)
	}
	return len(ids), ids, nil
}

// DeleteQueryExecution deletes an Athena query execution by ID.
func (s *QueryExecutionStore) DeleteQueryExecution(id string) error {
	qe, err := s.GetQueryExecution(id)
	if err != nil {
		return err
	}
	if qe.WorkGroup != "" {
		_ = s.BaseStore.Delete(s.workGroupIndexKey(qe.WorkGroup, id))
	}
	return s.BaseStore.Delete(id)
}

// ListQueryExecutionIDs returns query execution IDs for a work group.
func (s *QueryExecutionStore) ListQueryExecutionIDs(workGroup string, maxResults int) ([]string, error) {
	if workGroup != "" {
		return s.listQueryExecutionIDsByIndex(workGroup, maxResults)
	}

	var ids []string
	err := s.ForEach(func(key string, value []byte) error {
		if len(key) > 0 && key[0] == '#' {
			return nil
		}
		var p pb.QueryExecution
		if err := proto.Unmarshal(value, &p); err != nil {
			return nil
		}
		ids = append(ids, p.QueryExecutionId)
		return nil
	})
	return ids, err
}

var errStopScan = fmt.Errorf("stop scan")

func (s *QueryExecutionStore) listQueryExecutionIDsByIndex(workGroup string, maxResults int) ([]string, error) {
	prefix := "#wg:" + workGroup + ":"
	var ids []string
	err := s.ScanPrefix(prefix, func(key string, _ []byte) error {
		id := extractIDFromIndexKey(key)
		if id != "" {
			ids = append(ids, id)
		}
		if maxResults > 0 && len(ids) >= maxResults {
			return errStopScan
		}
		return nil
	})
	if err != nil && err != errStopScan {
		return nil, err
	}
	return ids, nil
}

// DeleteQueryExecutionsByWorkGroup removes all query executions for a workgroup and returns their IDs.
func (s *QueryExecutionStore) DeleteQueryExecutionsByWorkGroup(workGroup string) ([]string, error) {
	prefix := "#wg:" + workGroup + ":"
	var ids []string
	err := s.ScanPrefix(prefix, func(key string, _ []byte) error {
		id := extractIDFromIndexKey(key)
		if id != "" {
			ids = append(ids, id)
			_ = s.BaseStore.Delete(key)
			_ = s.BaseStore.Delete(id)
		}
		return nil
	})
	return ids, err
}

// ResultStore provides Athena query result storage operations.
type ResultStore struct {
	*common.BaseStore
}

// NewResultStore creates a new Athena result store.
func NewResultStore(store storage.BasicStorage, region string) *ResultStore {
	return &ResultStore{
		BaseStore: common.NewBaseStore(store.Bucket(resultBucketName(region)), "athena-result"),
	}
}

// StoreResult stores an Athena query result.
func (s *ResultStore) StoreResult(queryExecutionId string, result *QueryResult) error {
	return s.PutProto(queryExecutionId, QueryResultToProto(result))
}

// GetResult retrieves an Athena query result by execution ID.
func (s *ResultStore) GetResult(queryExecutionId string) (*QueryResult, error) {
	var p pb.QueryResult
	if err := s.GetProto(queryExecutionId, &p); err != nil {
		return nil, fmt.Errorf("result not found")
	}
	return ProtoToQueryResult(&p), nil
}

// DeleteResult deletes an Athena query result by execution ID.
func (s *ResultStore) DeleteResult(queryExecutionId string) error {
	if _, err := s.GetResult(queryExecutionId); err != nil {
		return err
	}
	return s.BaseStore.Delete(queryExecutionId)
}

// DeleteResultsByIDs removes query results for the given execution IDs.
func (s *ResultStore) DeleteResultsByIDs(ids []string) {
	for _, id := range ids {
		_ = s.BaseStore.Delete(id)
	}
}

// DataCatalogStore provides Athena data catalog storage operations.
type DataCatalogStore struct {
	*common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	kl         common.KeyLocker
}

// NewDataCatalogStore creates a new Athena data catalog store.
func NewDataCatalogStore(store storage.BasicStorage, accountID, region string) *DataCatalogStore {
	return &DataCatalogStore{
		BaseStore:  common.NewBaseStore(store.Bucket(dataCatalogBucketName(region)), "athena-data-catalog"),
		TagStore:   common.NewTagStoreWithRegion(store, "athena-data-catalog", region),
		arnBuilder: svcarn.NewARNBuilder(accountID, region),
	}
}

// CreateDataCatalog creates a new Athena data catalog.
func (s *DataCatalogStore) CreateDataCatalog(dc *DataCatalog) error {
	return s.kl.WithLock(dc.Name, func() error {
		if s.Exists(dc.Name) {
			return ErrDataCatalogAlreadyExists
		}
		return s.PutProto(dc.Name, DataCatalogToProto(dc))
	})
}

// GetDataCatalog retrieves an Athena data catalog by name.
func (s *DataCatalogStore) GetDataCatalog(name string) (*DataCatalog, error) {
	var p pb.DataCatalog
	if err := s.GetProto(name, &p); err != nil {
		return nil, ErrDataCatalogNotFound
	}
	return ProtoToDataCatalog(&p), nil
}

// UpdateDataCatalog updates an existing Athena data catalog.
func (s *DataCatalogStore) UpdateDataCatalog(dc *DataCatalog) error {
	return s.kl.WithLock(dc.Name, func() error {
		if !s.Exists(dc.Name) {
			return ErrDataCatalogNotFound
		}
		return s.PutProto(dc.Name, DataCatalogToProto(dc))
	})
}

// DeleteDataCatalog deletes an Athena data catalog by name.
func (s *DataCatalogStore) DeleteDataCatalog(name string) error {
	return s.kl.WithLock(name, func() error {
		if !s.Exists(name) {
			return ErrDataCatalogNotFound
		}
		arn := s.arnBuilder.Athena().DataCatalog(name)
		_ = s.TagStore.Delete(arn)
		return s.BaseStore.Delete(name)
	})
}

// ListDataCatalogs returns all Athena data catalogs.
func (s *DataCatalogStore) ListDataCatalogs() ([]*DataCatalog, error) {
	result, err := common.ListProto[*pb.DataCatalog](s.BaseStore, common.ListOptions{}, func() *pb.DataCatalog { return &pb.DataCatalog{} }, nil)
	if err != nil {
		return nil, err
	}
	catalogs := make([]*DataCatalog, len(result.Items))
	for i, p := range result.Items {
		catalogs[i] = ProtoToDataCatalog(p)
	}
	return catalogs, nil
}

// DatabaseStore provides Athena database storage operations.
type DatabaseStore struct {
	*common.BaseStore
}

// NewDatabaseStore creates a new Athena database store.
func NewDatabaseStore(store storage.BasicStorage, region string) *DatabaseStore {
	return &DatabaseStore{
		BaseStore: common.NewBaseStore(store.Bucket(databaseBucketName(region)), "athena-database"),
	}
}

func (s *DatabaseStore) databaseKey(catalog, name string) string {
	return catalog + ":" + name
}

// CreateDatabase creates a new Athena database.
func (s *DatabaseStore) CreateDatabase(catalog string, db *Database) error {
	key := s.databaseKey(catalog, db.Name)

	if s.Exists(key) {
		return fmt.Errorf("database already exists")
	}

	return s.PutProto(key, DatabaseToProto(db))
}

// GetDatabase retrieves an Athena database by catalog and name.
func (s *DatabaseStore) GetDatabase(catalog, name string) (*Database, error) {
	key := s.databaseKey(catalog, name)

	var p pb.Database
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrDatabaseNotFound
	}
	return ProtoToDatabase(&p), nil
}

// ListDatabases returns all databases in an Athena catalog.
func (s *DatabaseStore) ListDatabases(catalog string) ([]*Database, error) {
	prefix := catalog + ":"

	result, err := common.ListProto[*pb.Database](s.BaseStore, common.ListOptions{Prefix: prefix}, func() *pb.Database { return &pb.Database{} }, nil)
	if err != nil {
		return nil, err
	}
	databases := make([]*Database, len(result.Items))
	for i, p := range result.Items {
		databases[i] = ProtoToDatabase(p)
	}
	return databases, nil
}

// DeleteDatabase deletes an Athena database by catalog and name.
func (s *DatabaseStore) DeleteDatabase(catalog, name string) error {
	key := s.databaseKey(catalog, name)
	if !s.Exists(key) {
		return ErrDatabaseNotFound
	}
	return s.BaseStore.Delete(key)
}

// TableStore provides Athena table metadata storage operations.
type TableStore struct {
	*common.BaseStore
}

// NewTableStore creates a new Athena table store.
func NewTableStore(store storage.BasicStorage, region string) *TableStore {
	return &TableStore{
		BaseStore: common.NewBaseStore(store.Bucket(tableBucketName(region)), "athena-table"),
	}
}

func (s *TableStore) tableKey(catalog, database, name string) string {
	return catalog + ":" + database + ":" + name
}

// CreateTable creates a new Athena table.
func (s *TableStore) CreateTable(catalog, database string, table *TableMetadata) error {
	key := s.tableKey(catalog, database, table.Name)

	if s.Exists(key) {
		return fmt.Errorf("table already exists")
	}

	return s.PutProto(key, TableMetadataToProto(table))
}

// GetTable retrieves an Athena table by catalog, database, and name.
func (s *TableStore) GetTable(catalog, database, name string) (*TableMetadata, error) {
	key := s.tableKey(catalog, database, name)

	var p pb.TableMetadata
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrTableNotFound
	}
	return ProtoToTableMetadata(&p), nil
}

// ListTables returns all tables in an Athena database.
func (s *TableStore) ListTables(catalog, database string) ([]*TableMetadata, error) {
	prefix := catalog + ":" + database + ":"

	result, err := common.ListProto[*pb.TableMetadata](s.BaseStore, common.ListOptions{Prefix: prefix}, func() *pb.TableMetadata { return &pb.TableMetadata{} }, nil)
	if err != nil {
		return nil, err
	}
	tables := make([]*TableMetadata, len(result.Items))
	for i, p := range result.Items {
		tables[i] = ProtoToTableMetadata(p)
	}
	return tables, nil
}

// DeleteTable deletes an Athena table by catalog, database, and name.
func (s *TableStore) DeleteTable(catalog, database, name string) error {
	key := s.tableKey(catalog, database, name)
	if !s.Exists(key) {
		return ErrTableNotFound
	}
	return s.BaseStore.Delete(key)
}

// UpdateTable updates an existing Athena table.
func (s *TableStore) UpdateTable(catalog, database string, table *TableMetadata) error {
	key := s.tableKey(catalog, database, table.Name)
	if !s.Exists(key) {
		return ErrTableNotFound
	}
	return s.PutProto(key, TableMetadataToProto(table))
}

// DeleteTablesByDatabase removes all table metadata entries for the given catalog and database.
func (s *TableStore) DeleteTablesByDatabase(catalog, database string) error {
	tables, err := s.ListTables(catalog, database)
	if err != nil {
		return err
	}
	for _, t := range tables {
		key := s.tableKey(catalog, database, t.Name)
		_ = s.BaseStore.Delete(key)
	}
	return nil
}

// TableDataStore provides Athena table data storage operations.
type TableDataStore struct {
	*common.BaseStore
}

// NewTableDataStore creates a new Athena table data store.
func NewTableDataStore(store storage.BasicStorage, region string) *TableDataStore {
	return &TableDataStore{
		BaseStore: common.NewBaseStore(store.Bucket(tableDataBucketName(region)), "athena-table-data"),
	}
}

func (s *TableDataStore) tableDataKey(catalog, database, table string) string {
	return catalog + ":" + database + ":" + table
}

// StoreTableData stores data for an Athena table.
func (s *TableDataStore) StoreTableData(catalog, database, table string, storedTable *StoredTable) error {
	key := s.tableDataKey(catalog, database, table)
	proto, err := StoredTableToProto(storedTable)
	if err != nil {
		return fmt.Errorf("store table data: %w", err)
	}
	return s.PutProto(key, proto)
}

// GetTableData retrieves data for an Athena table.
func (s *TableDataStore) GetTableData(catalog, database, table string) (*StoredTable, error) {
	key := s.tableDataKey(catalog, database, table)

	var p pb.StoredTable
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrTableNotFound
	}
	return ProtoToStoredTable(&p), nil
}

// DeleteTableData deletes data for an Athena table.
func (s *TableDataStore) DeleteTableData(catalog, database, table string) error {
	key := s.tableDataKey(catalog, database, table)
	if !s.Exists(key) {
		return nil
	}
	return s.BaseStore.Delete(key)
}

// DeleteTableDataByDatabase removes all table data (rows) for the given catalog and database.
func (s *TableDataStore) DeleteTableDataByDatabase(catalog, database string) error {
	prefix := catalog + ":" + database + ":"
	result, err := common.ListProto[*pb.StoredTable](s.BaseStore, common.ListOptions{Prefix: prefix}, func() *pb.StoredTable { return &pb.StoredTable{} }, nil)
	if err != nil {
		return err
	}
	for _, p := range result.Items {
		key := s.tableDataKey(catalog, database, p.TableName)
		_ = s.BaseStore.Delete(key)
	}
	return nil
}

// --- Capacity Reservation ---

func capacityReservationBucketName(region string) string {
	return "athena_capacity_reservations-" + region
}

// CapacityReservationStore manages Athena capacity reservations.
type CapacityReservationStore struct {
	*common.BaseStore
	accountID string
	region    string
}

// NewCapacityReservationStore creates a new CapacityReservationStore.
func NewCapacityReservationStore(store storage.BasicStorage, accountID, region string) *CapacityReservationStore {
	return &CapacityReservationStore{
		BaseStore: common.NewBaseStore(store.Bucket(capacityReservationBucketName(region)), "athena-capacity-reservation"),
		accountID: accountID,
		region:    region,
	}
}

// CreateCapacityReservation persists a new capacity reservation.
func (s *CapacityReservationStore) CreateCapacityReservation(cr *CapacityReservation) (*CapacityReservation, error) {
	existing, err := s.GetCapacityReservation(cr.Name)
	if err == nil && existing != nil {
		return nil, ErrCapacityReservationAlreadyExists
	}
	cr.CreationTime = time.Now().UTC()
	cr.Status = CapacityReservationStatusActive
	return cr, s.Put(cr.Name, cr)
}

// GetCapacityReservation retrieves a capacity reservation by name.
func (s *CapacityReservationStore) GetCapacityReservation(name string) (*CapacityReservation, error) {
	var cr CapacityReservation
	if err := s.Get(name, &cr); err != nil {
		return nil, ErrCapacityReservationNotFound
	}
	return &cr, nil
}

// UpdateCapacityReservation updates a capacity reservation.
func (s *CapacityReservationStore) UpdateCapacityReservation(cr *CapacityReservation) error {
	return s.Put(cr.Name, cr)
}

// ListCapacityReservations returns all capacity reservations, optionally
// filtered by work group.
func (s *CapacityReservationStore) ListCapacityReservations(workGroup string) ([]*CapacityReservation, error) {
	var result []*CapacityReservation
	err := s.ForEach(func(_ string, value []byte) error {
		var cr CapacityReservation
		if err := json.Unmarshal(value, &cr); err != nil {
			return nil
		}
		if workGroup != "" && cr.TargetWorkGroup != workGroup {
			return nil
		}
		result = append(result, &cr)
		return nil
	})
	return result, err
}

// DeleteCapacityReservation deletes a capacity reservation by name.
func (s *CapacityReservationStore) DeleteCapacityReservation(name string) error {
	return s.BaseStore.Delete(name)
}
