package timestreamwrite

import (
	"net/http"
	"vorpalstacks/internal/common/defaults"

	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/timestreamwrite"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// This file is the sole location permitted to import the tsstore package
// directly: it contains the getStore helpers and pure proto conversion
// helpers (toPb* functions) that translate store types to proto types for
// response marshalling.

// getStoreFromHeader resolves the per-region tsWriteStores for the given
// gRPC-Web request headers.
func (h *AdminHandler) getStoreFromHeader(header http.Header) (*tsWriteStores, error) {
	region := defaults.GetRegionFromHeader(header)
	stores, err := h.service.GetStoresForRegion(region)
	if err != nil {
		return nil, err
	}
	return stores, nil
}

// toPbDatabase converts a transport-agnostic DatabaseResult to a proto
// Database for gRPC-Web response marshalling.
func toPbDatabase(db *DatabaseResult) *pb.Database {
	pbDb := &pb.Database{
		Arn:             db.ARN,
		Databasename:    db.DatabaseName,
		Tablecount:      proto.Int64(db.TableCount),
		Kmskeyid:        db.KmsKeyId,
		Creationtime:    db.CreationTime.Format(timeutils.ISO8601UTCFormat),
		Lastupdatedtime: db.LastUpdatedTime.Format(timeutils.ISO8601UTCFormat),
	}
	return pbDb
}

// toPbTable converts a transport-agnostic TableResult to a proto Table for
// gRPC-Web response marshalling.
func toPbTable(t *TableResult) *pb.Table {
	table := &pb.Table{
		Arn:             t.ARN,
		Tablename:       t.TableName,
		Databasename:    t.DatabaseName,
		Creationtime:    t.CreationTime.Format(timeutils.ISO8601UTCFormat),
		Lastupdatedtime: t.LastUpdatedTime.Format(timeutils.ISO8601UTCFormat),
	}

	switch t.TableStatus {
	case tsstore.TableStatusActive:
		table.Tablestatus = pb.TableStatus_TABLE_STATUS_ACTIVE
	case tsstore.TableStatusDeleting:
		table.Tablestatus = pb.TableStatus_TABLE_STATUS_DELETING
	case tsstore.TableStatusRestoring:
		table.Tablestatus = pb.TableStatus_TABLE_STATUS_RESTORING
	}

	if t.RetentionProperties != nil {
		table.Retentionproperties = &pb.RetentionProperties{
			Memorystoreretentionperiodinhours:  t.RetentionProperties.MemoryStoreRetentionPeriodInHours,
			Magneticstoreretentionperiodindays: t.RetentionProperties.MagneticStoreRetentionPeriodInDays,
		}
	}

	if t.Schema != nil && len(t.Schema.CompositePartitionKey) > 0 {
		schema := &pb.Schema{}
		for _, pk := range t.Schema.CompositePartitionKey {
			cpk := &pb.PartitionKey{}
			switch pk.Type {
			case tsstore.PartitionKeyTypeMeasure:
				cpk.Type = pb.PartitionKeyType_PARTITION_KEY_TYPE_MEASURE
			default:
				cpk.Type = pb.PartitionKeyType_PARTITION_KEY_TYPE_DIMENSION
			}
			if pk.Name != "" {
				cpk.Name = pk.Name
			}
			switch pk.EnforcementInRecord {
			case tsstore.EnforcementInRecordRequired:
				cpk.Enforcementinrecord = pb.PartitionKeyEnforcementLevel_PARTITION_KEY_ENFORCEMENT_LEVEL_REQUIRED
			default:
				cpk.Enforcementinrecord = pb.PartitionKeyEnforcementLevel_PARTITION_KEY_ENFORCEMENT_LEVEL_OPTIONAL
			}
			schema.Compositepartitionkey = append(schema.Compositepartitionkey, cpk)
		}
		table.Schema = schema
	}

	return table
}
