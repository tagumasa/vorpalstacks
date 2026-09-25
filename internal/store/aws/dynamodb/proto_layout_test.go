package dynamodb

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

// The storage proto's header documents the numbering convention every
// message family follows: a shared semantic layout — identity first, then
// schema and configuration, then runtime state, then timestamps — with
// same-meaning members across a family on the same number and new members
// inserted at their semantic position, never appended past it. These
// probes pin that convention on the generated descriptors: a member added
// on a trailing free number, or a family member renumbered alone, fails
// here instead of silently drifting the persisted layout.

func descriptor(t *testing.T, m protoreflect.ProtoMessage) protoreflect.MessageDescriptor {
	t.Helper()
	return m.ProtoReflect().Descriptor()
}

func fieldNumber(t *testing.T, m protoreflect.MessageDescriptor, name protoreflect.Name) protoreflect.FieldNumber {
	t.Helper()
	fd := m.Fields().ByName(name)
	if fd == nil {
		t.Fatalf("%v has no field %v", m.FullName(), name)
	}
	return fd.Number()
}

// increasing asserts the listed members appear in strictly rising number
// order, the form a block takes when it is authored in layout order.
func increasing(t *testing.T, m protoreflect.MessageDescriptor, names ...protoreflect.Name) {
	t.Helper()
	prev := protoreflect.FieldNumber(0)
	for _, n := range names {
		num := fieldNumber(t, m, n)
		if num <= prev {
			t.Fatalf("%v: %v = %d sits at or before its preceding listed member (%d); the documented layout orders members identity → schema and configuration → runtime state → timestamps, and members are numbered in that order", m.FullName(), n, num, prev)
		}
		prev = num
	}
}

// sameNumber asserts two same-meaning members of one family carry one
// number, the rule that keeps the family layout shared.
func sameNumber(t *testing.T, a protoreflect.MessageDescriptor, aName protoreflect.Name, b protoreflect.MessageDescriptor, bName protoreflect.Name) {
	t.Helper()
	an, bn := fieldNumber(t, a, aName), fieldNumber(t, b, bName)
	if an != bn {
		t.Fatalf("family members %v.%v and %v.%v carry different numbers (%d vs %d); same-meaning members across a family share one number", a.FullName(), aName, b.FullName(), bName, an, bn)
	}
}

func TestStorageProtoIndexFamilyLayout(t *testing.T) {
	gsi := descriptor(t, &pb.GlobalSecondaryIndex{})
	lsi := descriptor(t, &pb.LocalSecondaryIndex{})
	vec := descriptor(t, &pb.VectorIndex{})

	// The shared slots: name, arn, schema, projection, capacity, status,
	// size, count — one number each across the family.
	sameNumber(t, gsi, "index_name", lsi, "index_name")
	sameNumber(t, gsi, "index_name", vec, "index_name")
	sameNumber(t, gsi, "key_schema", lsi, "key_schema")
	sameNumber(t, gsi, "projection", lsi, "projection")
	sameNumber(t, gsi, "projection", vec, "projection")
	sameNumber(t, gsi, "index_arn", vec, "index_arn")
	sameNumber(t, gsi, "index_status", vec, "index_status")
	sameNumber(t, gsi, "index_size_bytes", lsi, "index_size_bytes")
	sameNumber(t, gsi, "index_size_bytes", vec, "index_size_bytes")
	sameNumber(t, gsi, "item_count", lsi, "item_count")
	sameNumber(t, gsi, "item_count", vec, "item_count")

	// Capacity is configuration: it sits between the schema block and the
	// runtime-state members, not appended after them.
	increasing(t, gsi,
		"index_name", "index_arn", "key_schema", "projection",
		"provisioned_throughput", "on_demand_throughput", "warm_throughput",
		"index_status", "index_size_bytes", "item_count")

	// Family-specific members follow the shared block.
	increasing(t, vec, "item_count", "vector_attribute_name", "dimensions", "distance_function", "backfilling")
}

func TestStorageProtoTableLayout(t *testing.T) {
	table := descriptor(t, &pb.Table{})

	// Identity, then schema and configuration, then runtime state, then
	// timestamps — each member numbered inside its section.
	increasing(t, table,
		"name", "arn", "table_id", "key_schema", "attribute_definitions",
		"billing_mode", "provisioned_throughput", "on_demand_throughput", "warm_throughput",
		"global_secondary_indexes", "local_secondary_indexes", "vector_indexes",
		"stream_specification", "sse_description", "table_class",
		"deletion_protection_enabled", "time_to_live", "point_in_time_recovery",
		"resource_policy", "resource_policy_revision_id",
		"kinesis_data_stream_destinations",
		"contributor_insights_enabled", "contributor_insights_mode", "contributor_insights_updated_at",
		"global_table_source_arn",
		"status", "table_size_bytes", "item_count", "stream_arn", "latest_stream_label",
		"restore_summary", "billing_mode_switches",
		"creation_date_time", "last_updated_date_time")
}

func TestStorageProtoBackupLayout(t *testing.T) {
	backup := descriptor(t, &pb.Backup{})
	table := descriptor(t, &pb.Table{})

	// Backup identity, then the source-table summary with the source
	// table's identity inside it, then the snapshot, state and times.
	increasing(t, backup,
		"backup_name", "backup_arn",
		"source_table_name", "source_table_arn", "source_table_id",
		"source_table_creation_time", "source_table_size_bytes", "source_table_item_count",
		"key_schema", "attribute_definitions", "billing_mode", "provisioned_throughput",
		"global_secondary_indexes", "local_secondary_indexes", "vector_indexes",
		"backup_status", "backup_type", "backup_size_bytes",
		"backup_creation_date_time", "backup_expiry_date_time")

	// The snapshot block mirrors the Table's schema-and-configuration
	// block member-for-member: the same members in the same order.
	increasing(t, table,
		"key_schema", "attribute_definitions", "billing_mode", "provisioned_throughput",
		"global_secondary_indexes", "local_secondary_indexes", "vector_indexes")
}

func TestStorageProtoImportExportFamilyLayout(t *testing.T) {
	imp := descriptor(t, &pb.ImportTableDescription{})
	exp := descriptor(t, &pb.ExportDescription{})

	sameNumber(t, imp, "import_arn", exp, "export_arn")
	sameNumber(t, imp, "import_status", exp, "export_status")
	for _, name := range []protoreflect.Name{
		"table_arn", "table_id", "start_time", "end_time",
		"failure_code", "failure_message", "client_token", "s3_bucket_source",
	} {
		sameNumber(t, imp, name, exp, name)
	}

	// Family-specific members follow the shared block.
	increasing(t, imp, "s3_bucket_source", "processed_item_count")
	increasing(t, exp, "s3_bucket_source", "item_count")
}

func TestStorageProtoClosedDomainMembersAreEnums(t *testing.T) {
	imp := descriptor(t, &pb.ImportTableDescription{})
	exp := descriptor(t, &pb.ExportDescription{})
	elem := descriptor(t, &pb.SearchSchemaElement{})
	vec := descriptor(t, &pb.VectorIndex{})

	for _, tc := range []struct {
		m    protoreflect.MessageDescriptor
		name protoreflect.Name
	}{
		{imp, "input_format"},
		{imp, "input_compression_type"},
		{exp, "export_format"},
		{exp, "export_type"},
		{elem, "search_schema_element_type"},
		{vec, "distance_function"},
	} {
		fd := tc.m.Fields().ByName(tc.name)
		if fd == nil {
			t.Fatalf("%v has no field %v", tc.m.FullName(), tc.name)
		}
		if fd.Kind() != protoreflect.EnumKind {
			t.Fatalf("%v.%v is %v, not an enum; members whose AWS counterpart is a closed enum are stored as enums", tc.m.FullName(), tc.name, fd.Kind())
		}
		zero := fd.Enum().Values().Get(0).Name()
		if !strings.HasSuffix(string(zero), "UNSPECIFIED") {
			t.Fatalf("%v.%v zero value %v is not the UNSPECIFIED absence value", tc.m.FullName(), tc.name, zero)
		}
	}
}
