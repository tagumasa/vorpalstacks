// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"time"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func tableBucketName(region string) string {
	return "dynamodb_tables-" + region
}

// TableStore manages DynamoDB table metadata in persistent storage.
type TableStore struct {
	*common.BaseStore
	TagStore         *common.TagStore
	arnBuilder       *svcarn.DynamoDBBuilder
	keyLocker        common.KeyLocker
	region           string
	autoScalingStore *common.BaseStore
}

// NewTableStore creates a new store for DynamoDB tables.
func NewTableStore(store storage.BasicStorage, accountId, region string) *TableStore {
	return &TableStore{
		BaseStore:        common.NewBaseStore(store.Bucket(tableBucketName(region)), "dynamodb"),
		TagStore:         common.NewTagStoreWithRegion(store, "dynamodb", region),
		arnBuilder:       svcarn.NewARNBuilder(accountId, region).DynamoDB(),
		region:           region,
		autoScalingStore: common.NewBaseStore(store.Bucket("dynamodb_autoscaling-"+region), "dynamodb"),
	}
}

// Get retrieves a table by its name.
func (s *TableStore) Get(name string) (*Table, error) {
	var pbTable pb.Table
	if err := s.BaseStore.GetProto(name, &pbTable); err != nil {
		return nil, err
	}
	return ProtoToTable(&pbTable), nil
}

// CreateTableParams carries every field TableStore.Create persists. Name,
// KeySchema and BillingMode are required; the rest are optional.
type CreateTableParams struct {
	Name                      string
	KeySchema                 []*KeySchemaElement
	AttributeDefinitions      []*AttributeDefinition
	BillingMode               BillingMode
	ProvisionedThroughput     *ProvisionedThroughput
	GlobalSecondaryIndexes    []*GlobalSecondaryIndex
	LocalSecondaryIndexes     []*LocalSecondaryIndex
	StreamSpecification       *StreamSpecification
	Tags                      []types.Tag
	DeletionProtectionEnabled bool
}

// Create creates a new DynamoDB table.
func (s *TableStore) Create(params CreateTableParams) (*Table, error) {
	s.keyLocker.Lock(params.Name)
	defer s.keyLocker.Unlock(params.Name)
	if s.Exists(params.Name) {
		return nil, ErrTableAlreadyExists
	}

	now := time.Now().UTC()
	table := &Table{
		Name:                      params.Name,
		ARN:                       s.arnBuilder.Table(params.Name),
		Status:                    TableStatusActive,
		CreationDateTime:          now,
		LastUpdatedDateTime:       now,
		KeySchema:                 params.KeySchema,
		AttributeDefinitions:      params.AttributeDefinitions,
		BillingMode:               params.BillingMode,
		ProvisionedThroughput:     params.ProvisionedThroughput,
		GlobalSecondaryIndexes:    params.GlobalSecondaryIndexes,
		LocalSecondaryIndexes:     params.LocalSecondaryIndexes,
		StreamSpecification:       params.StreamSpecification,
		Tags:                      params.Tags,
		TableSizeBytes:            0,
		ItemCount:                 0,
		DeletionProtectionEnabled: params.DeletionProtectionEnabled,
	}

	for _, g := range table.GlobalSecondaryIndexes {
		g.IndexArn = table.ARN + "/index/" + g.IndexName
	}

	if params.StreamSpecification != nil && params.StreamSpecification.StreamEnabled {
		table.StreamArn = table.ARN + "/stream/" + now.Format("2006-01-02T15:04:05.000")
		table.LatestStreamLabel = now.Format("2006-01-02T15:04:05.000")
	}

	if err := s.BaseStore.PutProto(params.Name, TableToProto(table)); err != nil {
		return nil, err
	}

	return table, nil
}

// Put stores or updates a DynamoDB table under its record lock. Callers
// that read the table before writing it must use Update instead, which
// holds the lock across both halves of the read-modify-write.
func (s *TableStore) Put(table *Table) error {
	return s.keyLocker.WithLock(table.Name, func() error {
		return s.put(table)
	})
}

// put writes the table record without locking; the caller must already hold
// the table's record lock.
func (s *TableStore) put(table *Table) error {
	return s.BaseStore.PutProto(table.Name, TableToProto(table))
}

// Update loads the table under its record lock, applies mutate to it, and
// persists the result. A mutate error aborts without writing. This is the
// single locked read-modify-write path for whole table records: the
// storage layer is read-committed, so an unlocked Get→Put sequence can
// lose concurrent changes (metric deltas, settings writes).
func (s *TableStore) Update(name string, mutate func(*Table) error) (*Table, error) {
	var updated *Table
	err := s.keyLocker.WithLock(name, func() error {
		table, err := s.Get(name)
		if err != nil {
			return err
		}
		if err := mutate(table); err != nil {
			return err
		}
		if err := s.put(table); err != nil {
			return err
		}
		updated = table
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// WithTableLock runs fn while holding the table's record lock, serialising
// it against Update, Put, applyMetricDeltas, and the post-commit metric
// flush. DeleteTable uses it so a cascade (including its commit) cannot
// interleave with a concurrent flush that would otherwise write the record
// back after the cascade deleted it. fn must not queue table metric deltas
// for this table: the post-commit flush inside store.Update takes the same
// lock.
func (s *TableStore) WithTableLock(name string, fn func() error) error {
	return s.keyLocker.WithLock(name, fn)
}

// applyMetricDeltas adds item-count and size deltas to a table record under
// the table's record lock. The read and the write both happen inside the
// lock with immediate (non-batch) writes, so concurrent metric flushes and
// locked table updates cannot lose increments. A table that no longer
// exists (deleted after the carrying write committed) is not an error: its
// counters no longer exist.
func (s *TableStore) applyMetricDeltas(name string, countDelta, sizeDelta int64) error {
	if countDelta == 0 && sizeDelta == 0 {
		return nil
	}
	_, err := s.Update(name, func(table *Table) error {
		table.ItemCount += countDelta
		if table.ItemCount < 0 {
			table.ItemCount = 0
		}
		table.TableSizeBytes += sizeDelta
		if table.TableSizeBytes < 0 {
			table.TableSizeBytes = 0
		}
		return nil
	})
	if err != nil && common.IsNotFound(err) {
		return nil
	}
	return err
}

// Delete removes a DynamoDB table by name.
func (s *TableStore) Delete(name string) error {
	return s.keyLocker.WithLock(name, func() error {
		return s.BaseStore.Delete(name)
	})
}

// Exists checks if a DynamoDB table exists.
func (s *TableStore) Exists(name string) bool {
	return s.BaseStore.Exists(name)
}

// List returns a list of DynamoDB tables with pagination.
func (s *TableStore) List(marker string, limit int) ([]*Table, string, error) {
	opts := common.ListOptions{
		Marker:   marker,
		MaxItems: limit,
	}
	result, err := common.ListProto[*pb.Table](s.BaseStore, opts, func() *pb.Table { return &pb.Table{} }, nil)
	if err != nil {
		return nil, "", err
	}
	tables := make([]*Table, len(result.Items))
	for i, pbTable := range result.Items {
		tables[i] = ProtoToTable(pbTable)
	}
	if !result.IsTruncated {
		return tables, "", nil
	}
	return tables, result.NextMarker, nil
}

// UpdateItemCount updates the item count for a DynamoDB table through the
// locked metric-delta path shared with the post-commit flush.
func (s *TableStore) UpdateItemCount(name string, delta int64) error {
	return s.applyMetricDeltas(name, delta, 0)
}

// UpdateTableSize updates the table size for a DynamoDB table through the
// locked metric-delta path shared with the post-commit flush.
func (s *TableStore) UpdateTableSize(name string, delta int64) error {
	return s.applyMetricDeltas(name, 0, delta)
}

// ARNBuilder returns the ARN builder for DynamoDB tables.
func (s *TableStore) ARNBuilder() *svcarn.DynamoDBBuilder {
	return s.arnBuilder
}

// Tags returns the tag store for DynamoDB tables.
func (s *TableStore) Tags() *common.TagStore {
	return s.TagStore
}

// SetTimeToLive sets the time-to-live specification for a DynamoDB table.
func (s *TableStore) SetTimeToLive(name string, ttl *TimeToLiveSpecification) error {
	_, err := s.Update(name, func(table *Table) error {
		table.TimeToLive = ttl
		return nil
	})
	return err
}

// GetTimeToLive returns the time-to-live specification for a DynamoDB table.
func (s *TableStore) GetTimeToLive(name string) (*TimeToLiveSpecification, error) {
	table, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	return table.TimeToLive, nil
}

// SetPointInTimeRecovery sets the point-in-time recovery description for a DynamoDB table.
func (s *TableStore) SetPointInTimeRecovery(name string, pitr *PointInTimeRecoveryDescription) error {
	_, err := s.Update(name, func(table *Table) error {
		table.PointInTimeRecovery = pitr
		return nil
	})
	return err
}

// GetPointInTimeRecovery returns the point-in-time recovery description for a DynamoDB table.
func (s *TableStore) GetPointInTimeRecovery(name string) (*PointInTimeRecoveryDescription, error) {
	table, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	return table.PointInTimeRecovery, nil
}

// SetResourcePolicy sets the resource policy for a DynamoDB table.
func (s *TableStore) SetResourcePolicy(name string, policy string) error {
	_, err := s.Update(name, func(table *Table) error {
		table.ResourcePolicy = policy
		table.ResourcePolicyRevisionId++
		return nil
	})
	return err
}

// GetResourcePolicyRevisionId returns the current resource policy revision
// number for a table. Returns 0 if no policy has been set.
func (s *TableStore) GetResourcePolicyRevisionId(name string) (int, error) {
	table, err := s.Get(name)
	if err != nil {
		return 0, err
	}
	return table.ResourcePolicyRevisionId, nil
}

// GetResourcePolicy returns the resource policy for a DynamoDB table.
func (s *TableStore) GetResourcePolicy(name string) (string, error) {
	table, err := s.Get(name)
	if err != nil {
		return "", err
	}
	return table.ResourcePolicy, nil
}

// DeleteResourcePolicy removes the resource policy from a DynamoDB table.
func (s *TableStore) DeleteResourcePolicy(name string) error {
	_, err := s.Update(name, func(table *Table) error {
		table.ResourcePolicy = ""
		return nil
	})
	return err
}

// SetKinesisStreamingDestination sets the Kinesis streaming destination for a DynamoDB table.
func (s *TableStore) SetKinesisStreamingDestination(name string, destinations []*KinesisDataStreamDestination) error {
	_, err := s.Update(name, func(table *Table) error {
		table.KinesisDataStreamDestinations = destinations
		return nil
	})
	return err
}

// SetContributorInsights enables or disables contributor insights for a DynamoDB table.
// When mode is empty, the existing ContributorInsightsMode is preserved.
func (s *TableStore) SetContributorInsights(name string, enabled bool, mode string) error {
	_, err := s.Update(name, func(table *Table) error {
		table.ContributorInsightsEnabled = enabled
		if mode != "" {
			table.ContributorInsightsMode = mode
		}
		table.ContributorInsightsUpdatedAt = time.Now().UTC()
		return nil
	})
	return err
}

// SetAutoScalingSettings stores the auto-scaling settings for a table as a
// typed protobuf record. Auto-scaling policy execution (Application Auto
// Scaling) is not implemented; these settings round-trip for API
// compatibility.
func (s *TableStore) SetAutoScalingSettings(name string, settings *TableReplicaAutoScalingSettings) error {
	return s.autoScalingStore.PutProto(name, tableReplicaAutoScalingToProto(settings))
}

// GetAutoScalingSettings returns the stored auto-scaling settings for a
// table. Returns nil if no settings have been stored.
func (s *TableStore) GetAutoScalingSettings(name string) (*TableReplicaAutoScalingSettings, error) {
	var pbSettings pb.TableReplicaAutoScalingSettings
	if err := s.autoScalingStore.GetProto(name, &pbSettings); err != nil {
		if common.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return protoToTableReplicaAutoScaling(&pbSettings), nil
}
