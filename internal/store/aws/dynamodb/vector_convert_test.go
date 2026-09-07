package dynamodb

import (
	"testing"
)

func TestVectorIndexProtoRoundTrip(t *testing.T) {
	table := &Table{
		Name: "vec-table",
		ARN:  "arn:aws:dynamodb:us-east-1:123456789012:table/vec-table",
		VectorIndexes: []*VectorIndex{
			{
				IndexName:           "vec-idx",
				IndexArn:            "arn:aws:dynamodb:us-east-1:123456789012:table/vec-table/index/vec-idx",
				VectorAttributeName: "embedding",
				Dimensions:          3,
				DistanceFunction:    "EUCLIDEAN",
				Projection:          &Projection{ProjectionType: "INCLUDE", NonKeyAttributes: []string{"title"}},
				SearchSchema: []*SearchSchemaElement{
					{AttributeName: "category", SearchSchemaElementType: "HASH"},
					{AttributeName: "year", SearchSchemaElementType: "INLINE_FILTER"},
				},
				IndexStatus:    IndexStatusActive,
				Backfilling:    true,
				IndexSizeBytes: 1234,
				ItemCount:      56,
			},
		},
	}

	rt := ProtoToTable(TableToProto(table))
	if len(rt.VectorIndexes) != 1 {
		t.Fatalf("round trip lost vector indexes: %d", len(rt.VectorIndexes))
	}
	got := rt.VectorIndexes[0]
	want := table.VectorIndexes[0]
	if got.IndexName != want.IndexName ||
		got.IndexArn != want.IndexArn ||
		got.VectorAttributeName != want.VectorAttributeName ||
		got.Dimensions != want.Dimensions ||
		got.DistanceFunction != want.DistanceFunction ||
		got.IndexStatus != want.IndexStatus ||
		got.Backfilling != want.Backfilling ||
		got.IndexSizeBytes != want.IndexSizeBytes ||
		got.ItemCount != want.ItemCount {
		t.Errorf("scalar fields mismatch: got %+v", got)
	}
	if got.Projection == nil || got.Projection.ProjectionType != "INCLUDE" || len(got.Projection.NonKeyAttributes) != 1 {
		t.Errorf("projection round trip: %+v", got.Projection)
	}
	if len(got.SearchSchema) != 2 ||
		got.SearchSchema[0].AttributeName != "category" || got.SearchSchema[0].SearchSchemaElementType != "HASH" ||
		got.SearchSchema[1].AttributeName != "year" || got.SearchSchema[1].SearchSchemaElementType != "INLINE_FILTER" {
		t.Errorf("search schema round trip: %+v", got.SearchSchema)
	}
}

func TestBackupVectorIndexProtoRoundTrip(t *testing.T) {
	backup := &Backup{
		BackupName: "vb",
		VectorIndexes: []*VectorIndex{
			{IndexName: "vec-idx", VectorAttributeName: "embedding", Dimensions: 8, DistanceFunction: "COSINE"},
		},
	}
	rt := ProtoToBackup(BackupToProto(backup))
	if len(rt.VectorIndexes) != 1 || rt.VectorIndexes[0].IndexName != "vec-idx" || rt.VectorIndexes[0].Dimensions != 8 {
		t.Errorf("backup vector round trip: %+v", rt.VectorIndexes)
	}
}
