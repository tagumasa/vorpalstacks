// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Default throughput quotas per the DynamoDB quotas documentation, per
// Region: a table may provision at most 40,000 read and 40,000 write
// capacity units, and the account-wide total across all provisioned tables
// and GSIs is 80,000 read and 80,000 write capacity units.
const (
	AccountMaxReadCapacityUnits  = 80000
	AccountMaxWriteCapacityUnits = 80000
	TableMaxReadCapacityUnits    = 40000
	TableMaxWriteCapacityUnits   = 40000
)

func tableBucketName(region string) string {
	return "dynamodb_tables-" + region
}

func autoScalingBucketName(region string) string {
	return "dynamodb_autoscaling-" + region
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
		TagStore:         common.NewTagStoreWithRegion(store, "dynamodb", region, common.StandardTagBudget("com.amazon.coral.validate#ValidationException")),
		arnBuilder:       svcarn.NewARNBuilder(accountId, region).DynamoDB(),
		region:           region,
		autoScalingStore: common.NewBaseStore(store.Bucket(autoScalingBucketName(region)), "dynamodb"),
	}
}

// Get retrieves a table by its name. Absence answers with the package
// sentinel (wrapped in the store-error shape with the key), so both read
// paths of a table report one not-found contract.
func (s *TableStore) Get(name string) (*Table, error) {
	var pbTable pb.Table
	if err := s.BaseStore.GetProto(name, &pbTable); err != nil {
		if common.IsNotFound(err) {
			return nil, common.NewStoreErrorWithKey("dynamodb", "get", name, ErrTableNotFound)
		}
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
	// The table UUID is minted here, at the single creation path every
	// surface (create, restore, import, replica join) funnels through: the
	// identity is stable for the table's lifetime and distinct across
	// same-name generations, matching the model's TableId pattern
	// (^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$).
	table := &Table{
		Name:                      params.Name,
		ARN:                       s.arnBuilder.Table(params.Name),
		TableId:                   uuid.NewString(),
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
		TableSizeBytes:            0,
		ItemCount:                 0,
		DeletionProtectionEnabled: params.DeletionProtectionEnabled,
	}

	for _, g := range table.GlobalSecondaryIndexes {
		g.IndexArn = s.arnBuilder.Index(params.Name, g.IndexName)
	}

	if params.StreamSpecification != nil && params.StreamSpecification.StreamEnabled {
		streamLabel := now.Format("2006-01-02T15:04:05.000")
		table.StreamArn = s.arnBuilder.TableStream(params.Name, streamLabel)
		table.LatestStreamLabel = streamLabel
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
// counters no longer exist — and neither is a record of a DIFFERENT
// generation under the same name (a successor created between the write's
// commit and this flush): the queued deltas belong to a generation whose
// counters are gone, not to whatever carries the name now.
func (s *TableStore) applyMetricDeltas(name, tableId string, countDelta, sizeDelta int64) error {
	if countDelta == 0 && sizeDelta == 0 {
		return nil
	}
	_, err := s.Update(name, func(table *Table) error {
		if tableId != "" && table.TableId != tableId {
			return errTableGenerationGone
		}
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
	if err != nil && (common.IsNotFound(err) || errors.Is(err, errTableGenerationGone)) {
		return nil
	}
	return err
}

// Exists checks if a DynamoDB table exists.
func (s *TableStore) Exists(name string) bool {
	return s.BaseStore.Exists(name)
}

// List returns a list of DynamoDB tables with pagination.
func (s *TableStore) List(marker string, limit int) ([]*Table, string, error) {
	return listProtoConverted(s.BaseStore, marker, limit, func() *pb.Table { return &pb.Table{} }, ProtoToTable, nil)
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

// PolicyRevisionUnchecked instructs the revision-checked resource policy
// writes to skip the optimistic-lock comparison — the store-side value of
// the wire contract's absent ExpectedRevisionId.
const PolicyRevisionUnchecked = -1

// SetResourcePolicyExpected sets the table's resource policy and advances
// the revision inside one locked read-modify-write of the table record:
// the write applies only when the current revision equals expectedRev, so
// two concurrent writers that both observed the same revision cannot both
// apply — the loser answers ErrPolicyRevisionMismatch. Returns the new
// revision number.
func (s *TableStore) SetResourcePolicyExpected(name string, policy string, expectedRev int) (int, error) {
	var newRev int
	_, err := s.Update(name, func(table *Table) error {
		if expectedRev != PolicyRevisionUnchecked && table.ResourcePolicyRevisionId != expectedRev {
			return ErrPolicyRevisionMismatch
		}
		table.ResourcePolicy = policy
		table.ResourcePolicyRevisionId++
		newRev = table.ResourcePolicyRevisionId
		return nil
	})
	if err != nil {
		return 0, err
	}
	return newRev, nil
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

// DeleteResourcePolicyExpected removes the resource policy from a table
// under the same revision-checked, locked write SetResourcePolicyExpected
// applies. A successful delete advances the revision too: the delete
// output's documented RevisionId member is a fresh id, and without the
// advance the token would sit unmoved and a concurrent put holding the
// same expected revision could still apply past the delete. Returns the
// new revision number.
func (s *TableStore) DeleteResourcePolicyExpected(name string, expectedRev int) (int, error) {
	var newRev int
	_, err := s.Update(name, func(table *Table) error {
		if expectedRev != PolicyRevisionUnchecked && table.ResourcePolicyRevisionId != expectedRev {
			return ErrPolicyRevisionMismatch
		}
		table.ResourcePolicy = ""
		table.ResourcePolicyRevisionId++
		newRev = table.ResourcePolicyRevisionId
		return nil
	})
	if err != nil {
		return 0, err
	}
	return newRev, nil
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
func (s *TableStore) SetContributorInsights(name string, enabled bool, mode ContributorInsightsMode) error {
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

// UpdateAutoScalingSettings applies merge to the stored auto-scaling
// settings under the table's record lock and persists the result — the
// locked read-modify-write the whole-record Update provides for table
// records, for the settings record that lives in its own bucket keyed by
// the same table name. An unlocked read-merge-write here lets two
// concurrent updates start from the same stored state, and the second
// write silently discards the first's replicas. Auto-scaling policy
// execution (Application Auto Scaling) is not implemented; the settings
// round-trip for API compatibility.
func (s *TableStore) UpdateAutoScalingSettings(name string, merge func(existing *TableReplicaAutoScalingSettings) *TableReplicaAutoScalingSettings) (*TableReplicaAutoScalingSettings, error) {
	var merged *TableReplicaAutoScalingSettings
	err := s.keyLocker.WithLock(name, func() error {
		existing, err := s.GetAutoScalingSettings(name)
		if err != nil {
			return err
		}
		merged = merge(existing)
		return s.autoScalingStore.PutProto(name, tableReplicaAutoScalingToProto(merged))
	})
	if err != nil {
		return nil, err
	}
	return merged, nil
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
