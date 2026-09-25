package dynamodb

import (
	"testing"

	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

// TestEnumWireTablesRoundTrip pins the single-table contract of every closed
// domain: each value in a domain's wire table converts to its protobuf enum
// and back to itself. Because both directions read the same table, a value
// missing from the table would surface here as a round-trip to the zero
// value, keeping a future value addition from silently dropping on persist.
func TestEnumWireTablesRoundTrip(t *testing.T) {
	t.Run("table status", func(t *testing.T) {
		for _, entry := range tableStatusWire {
			if got := protoToTableStatus(tableStatusToProto(entry.value)); got != entry.value {
				t.Errorf("table status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("billing mode", func(t *testing.T) {
		for _, entry := range billingModeWire {
			if got := protoToBillingMode(billingModeToProto(entry.value)); got != entry.value {
				t.Errorf("billing mode %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("table class", func(t *testing.T) {
		for _, entry := range tableClassWire {
			if got := protoToTableClass(tableClassToProto(entry.value)); got != entry.value {
				t.Errorf("table class %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("sse status", func(t *testing.T) {
		for _, entry := range sseStatusWire {
			if got := protoToSSEStatus(sseStatusToProto(entry.value)); got != entry.value {
				t.Errorf("sse status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("sse type", func(t *testing.T) {
		for _, entry := range sseTypeWire {
			if got := protoToSSEType(sseTypeToProto(entry.value)); got != entry.value {
				t.Errorf("sse type %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("destination status", func(t *testing.T) {
		for _, entry := range destinationStatusWire {
			if got := protoToDestinationStatus(destinationStatusToProto(entry.value)); got != entry.value {
				t.Errorf("destination status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("contributor insights mode", func(t *testing.T) {
		for _, entry := range contributorInsightsModeWire {
			if got := protoToContributorInsightsMode(contributorInsightsModeToProto(entry.value)); got != entry.value {
				t.Errorf("contributor insights mode %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("acdt precision", func(t *testing.T) {
		for _, entry := range acdtPrecisionWire {
			if got := protoToApproximateCreationDateTimePrecision(approximateCreationDateTimePrecisionToProto(entry.value)); got != entry.value {
				t.Errorf("acdt precision %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("stream view type", func(t *testing.T) {
		for _, entry := range streamViewTypeWire {
			if got := protoToStreamViewType(streamViewTypeToProto(entry.value)); got != entry.value {
				t.Errorf("stream view type %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("ttl status", func(t *testing.T) {
		for _, entry := range ttlStatusWire {
			if got := protoToTTLStatus(ttlStatusToProto(entry.value)); got != entry.value {
				t.Errorf("ttl status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("pitr status", func(t *testing.T) {
		for _, entry := range pitrStatusWire {
			if got := protoToPointInTimeRecoveryStatus(pointInTimeRecoveryStatusToProto(entry.value)); got != entry.value {
				t.Errorf("pitr status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("backup status", func(t *testing.T) {
		for _, entry := range backupStatusWire {
			if got := protoToBackupStatus(backupStatusToProto(entry.value)); got != entry.value {
				t.Errorf("backup status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("backup type", func(t *testing.T) {
		for _, entry := range backupTypeWire {
			if got := protoToBackupType(backupTypeToProto(entry.value)); got != entry.value {
				t.Errorf("backup type %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("global table status", func(t *testing.T) {
		for _, entry := range globalTableStatusWire {
			if got := protoToGlobalTableStatus(globalTableStatusToProto(entry.value)); got != entry.value {
				t.Errorf("global table status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("replica status", func(t *testing.T) {
		for _, entry := range replicaStatusWire {
			if got := protoToReplicaStatus(replicaStatusToProto(entry.value)); got != entry.value {
				t.Errorf("replica status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("import status", func(t *testing.T) {
		for _, entry := range importStatusWire {
			if got := protoToImportStatus(importStatusToProto(entry.value)); got != entry.value {
				t.Errorf("import status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("export status", func(t *testing.T) {
		for _, entry := range exportStatusWire {
			if got := protoToExportStatus(exportStatusToProto(entry.value)); got != entry.value {
				t.Errorf("export status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("input format", func(t *testing.T) {
		for _, entry := range inputFormatWire {
			if got := protoToInputFormat(inputFormatToProto(entry.value)); got != entry.value {
				t.Errorf("input format %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("input compression", func(t *testing.T) {
		for _, entry := range inputCompressionTypeWire {
			if got := protoToInputCompressionType(inputCompressionTypeToProto(entry.value)); got != entry.value {
				t.Errorf("input compression %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("export format", func(t *testing.T) {
		for _, entry := range exportFormatWire {
			if got := protoToExportFormat(exportFormatToProto(entry.value)); got != entry.value {
				t.Errorf("export format %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("export type", func(t *testing.T) {
		for _, entry := range exportTypeWire {
			if got := protoToExportType(exportTypeToProto(entry.value)); got != entry.value {
				t.Errorf("export type %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("export view type", func(t *testing.T) {
		for _, entry := range exportViewTypeWire {
			if got := protoToExportViewType(exportViewTypeToProto(entry.value)); got != entry.value {
				t.Errorf("export view type %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("s3 sse algorithm", func(t *testing.T) {
		for _, entry := range s3SseAlgorithmWire {
			if got := protoToS3SseAlgorithm(s3SseAlgorithmToProto(entry.value)); got != entry.value {
				t.Errorf("s3 sse algorithm %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("projection type", func(t *testing.T) {
		for _, entry := range projectionTypeWire {
			if got := protoToProjectionType(projectionTypeToProto(entry.value)); got != entry.value {
				t.Errorf("projection type %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("key type", func(t *testing.T) {
		for _, entry := range keyTypeWire {
			if got := protoToKeyType(keyTypeToProto(entry.value)); got != entry.value {
				t.Errorf("key type %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("scalar attribute type", func(t *testing.T) {
		for _, entry := range scalarAttributeTypeWire {
			if got := protoToScalarAttributeType(scalarAttributeTypeToProto(entry.value)); got != entry.value {
				t.Errorf("scalar attribute type %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("index status", func(t *testing.T) {
		for _, entry := range indexStatusWire {
			if got := protoToIndexStatus(indexStatusToProto(entry.value)); got != entry.value {
				t.Errorf("index status %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("vector distance function", func(t *testing.T) {
		for _, entry := range vectorDistanceFunctionWire {
			if got := protoToVectorDistanceFunction(vectorDistanceFunctionToProto(entry.value)); got != entry.value {
				t.Errorf("distance function %q round-tripped as %q", entry.value, got)
			}
		}
	})
	t.Run("search schema element type", func(t *testing.T) {
		for _, entry := range searchSchemaElementTypeWire {
			if got := protoToSearchSchemaElementType(searchSchemaElementTypeToProto(entry.value)); got != entry.value {
				t.Errorf("search schema element type %q round-tripped as %q", entry.value, got)
			}
		}
	})
}

// TestEnumWireUnknownValuesFallToZero pins the fallback contract: a value
// outside the table persists as UNSPECIFIED and reads back as the empty
// value, on every domain alike.
func TestEnumWireUnknownValuesFallToZero(t *testing.T) {
	if got := tableStatusToProto(TableStatus("nonsense")); got != pb.TableStatus_TABLE_STATUS_UNSPECIFIED {
		t.Errorf("unknown table status mapped to %v", got)
	}
	if got := protoToTableStatus(pb.TableStatus_TABLE_STATUS_UNSPECIFIED); got != "" {
		t.Errorf("unspecified table status read back as %q", got)
	}
	if got := exportStatusToProto(ExportStatus("nonsense")); got != pb.ExportStatus_EXPORT_STATUS_UNSPECIFIED {
		t.Errorf("unknown export status mapped to %v", got)
	}
	if got := protoToExportStatus(pb.ExportStatus_EXPORT_STATUS_UNSPECIFIED); got != "" {
		t.Errorf("unspecified export status read back as %q", got)
	}
	if got := replicaStatusToProto(ReplicaStatus("nonsense")); got != pb.ReplicaStatus_REPLICA_STATUS_UNSPECIFIED {
		t.Errorf("unknown replica status mapped to %v", got)
	}
	if got := protoToReplicaStatus(pb.ReplicaStatus_REPLICA_STATUS_UNSPECIFIED); got != "" {
		t.Errorf("unspecified replica status read back as %q", got)
	}
	if got := vectorDistanceFunctionToProto(VectorDistanceFunction("nonsense")); got != pb.VectorDistanceFunction_VECTOR_DISTANCE_FUNCTION_UNSPECIFIED {
		t.Errorf("unknown distance function mapped to %v", got)
	}
	if got := protoToVectorDistanceFunction(pb.VectorDistanceFunction_VECTOR_DISTANCE_FUNCTION_UNSPECIFIED); got != "" {
		t.Errorf("unspecified distance function read back as %q", got)
	}
}
