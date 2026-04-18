// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"time"

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
	arnBuilder *svcarn.DynamoDBBuilder
	keyLocker  common.KeyLocker
	region     string
}

// NewTableStore creates a new store for DynamoDB tables.
func NewTableStore(store storage.BasicStorage, accountId, region string) *TableStore {
	return &TableStore{
		BaseStore:  common.NewBaseStore(store.Bucket(tableBucketName(region)), "dynamodb"),
		arnBuilder: svcarn.NewARNBuilder(accountId, region).DynamoDB(),
		region:     region,
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

// Create creates a new DynamoDB table.
func (s *TableStore) Create(
	name string,
	keySchema []*KeySchemaElement,
	attributeDefinitions []*AttributeDefinition,
	billingMode BillingMode,
	provisionedThroughput *ProvisionedThroughput,
	gsi []*GlobalSecondaryIndex,
	lsi []*LocalSecondaryIndex,
	streamSpec *StreamSpecification,
	tags []common.Tag,
	deletionProtectionEnabled bool,
) (*Table, error) {
	s.keyLocker.Lock(name)
	defer s.keyLocker.Unlock(name)
	if s.Exists(name) {
		return nil, ErrTableAlreadyExists
	}

	now := time.Now().UTC()
	table := &Table{
		Name:                      name,
		ARN:                       s.arnBuilder.Table(name),
		Status:                    TableStatusActive,
		CreationDateTime:          now,
		LastUpdatedDateTime:       now,
		KeySchema:                 keySchema,
		AttributeDefinitions:      attributeDefinitions,
		BillingMode:               billingMode,
		ProvisionedThroughput:     provisionedThroughput,
		GlobalSecondaryIndexes:    gsi,
		LocalSecondaryIndexes:     lsi,
		StreamSpecification:       streamSpec,
		Tags:                      tags,
		TableSizeBytes:            0,
		ItemCount:                 0,
		DeletionProtectionEnabled: deletionProtectionEnabled,
	}

	for _, g := range table.GlobalSecondaryIndexes {
		g.IndexArn = table.ARN + "/index/" + g.IndexName
	}

	if streamSpec != nil && streamSpec.StreamEnabled {
		table.StreamArn = table.ARN + "/stream/" + now.Format("20060102150405")
		table.LatestStreamLabel = now.Format("2006-01-02T15:04:05.000")
	}

	if err := s.BaseStore.PutProto(name, TableToProto(table)); err != nil {
		return nil, err
	}

	return table, nil
}

// Put stores or updates a DynamoDB table.
func (s *TableStore) Put(table *Table) error {
	return s.BaseStore.PutProto(table.Name, TableToProto(table))
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

// UpdateItemCount updates the item count for a DynamoDB table.
func (s *TableStore) UpdateItemCount(name string, delta int64) error {
	return s.keyLocker.WithLock(name, func() error {
		table, err := s.Get(name)
		if err != nil {
			return err
		}
		table.ItemCount += delta
		if table.ItemCount < 0 {
			table.ItemCount = 0
		}
		return s.Put(table)
	})
}

// UpdateTableSize updates the table size for a DynamoDB table.
func (s *TableStore) UpdateTableSize(name string, delta int64) error {
	return s.keyLocker.WithLock(name, func() error {
		table, err := s.Get(name)
		if err != nil {
			return err
		}
		table.TableSizeBytes += delta
		if table.TableSizeBytes < 0 {
			table.TableSizeBytes = 0
		}
		return s.Put(table)
	})
}

// SetTags sets the tags for a DynamoDB table.
func (s *TableStore) SetTags(name string, tags []common.Tag) error {
	return s.keyLocker.WithLock(name, func() error {
		table, err := s.Get(name)
		if err != nil {
			return err
		}
		table.Tags = tags
		return s.Put(table)
	})
}

// ARNBuilder returns the ARN builder for DynamoDB tables.
func (s *TableStore) ARNBuilder() *svcarn.DynamoDBBuilder {
	return s.arnBuilder
}

// GetPartitionKey returns the partition key attribute name for a table.
func (s *TableStore) GetPartitionKey(table *Table) string {
	for _, ks := range table.KeySchema {
		if ks.KeyType == KeyTypeHash {
			return ks.AttributeName
		}
	}
	return ""
}

// GetSortKey returns the sort key attribute name for a table.
func (s *TableStore) GetSortKey(table *Table) string {
	for _, ks := range table.KeySchema {
		if ks.KeyType == KeyTypeRange {
			return ks.AttributeName
		}
	}
	return ""
}

// SetTimeToLive sets the time-to-live specification for a DynamoDB table.
func (s *TableStore) SetTimeToLive(name string, ttl *TimeToLiveSpecification) error {
	return s.keyLocker.WithLock(name, func() error {
		table, err := s.Get(name)
		if err != nil {
			return err
		}
		table.TimeToLive = ttl
		return s.Put(table)
	})
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
	return s.keyLocker.WithLock(name, func() error {
		table, err := s.Get(name)
		if err != nil {
			return err
		}
		table.PointInTimeRecovery = pitr
		return s.Put(table)
	})
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
	return s.keyLocker.WithLock(name, func() error {
		table, err := s.Get(name)
		if err != nil {
			return err
		}
		table.ResourcePolicy = policy
		return s.Put(table)
	})
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
	return s.keyLocker.WithLock(name, func() error {
		table, err := s.Get(name)
		if err != nil {
			return err
		}
		table.ResourcePolicy = ""
		return s.Put(table)
	})
}

// SetKinesisStreamingDestination sets the Kinesis streaming destination for a DynamoDB table.
func (s *TableStore) SetKinesisStreamingDestination(name string, destinations []*KinesisDataStreamDestination) error {
	return s.keyLocker.WithLock(name, func() error {
		table, err := s.Get(name)
		if err != nil {
			return err
		}
		table.KinesisDataStreamDestinations = destinations
		return s.Put(table)
	})
}

// SetContributorInsights enables or disables contributor insights for a DynamoDB table.
func (s *TableStore) SetContributorInsights(name string, enabled bool) error {
	return s.keyLocker.WithLock(name, func() error {
		table, err := s.Get(name)
		if err != nil {
			return err
		}
		table.ContributorInsightsEnabled = enabled
		return s.Put(table)
	})
}
