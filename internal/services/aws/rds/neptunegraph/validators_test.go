package neptunegraph

import (
	"strings"
	"testing"
	"time"
)

func TestValidateGraphName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "mygraph", false},
		{"valid with hyphens", "my-graph", false},
		{"valid alphanumeric", "graph123", false},
		{"valid max length", "aabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghij", false},
		{"empty", "", true},
		{"starts with g-", "g-mygraph", true},
		{"starts with digit", "1graph", true},
		{"contains uppercase", "MyGraph", true},
		{"too long", "aabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklkl", true},
		{"ends with hyphen", "mygraph-", true},
		{"consecutive hyphens", "my--graph", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateGraphName(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validateGraphName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateProvisionedMemory(t *testing.T) {
	tests := []struct {
		name     string
		mem      int
		required bool
		wantErr  bool
	}{
		{"min valid not required", 16, false, false},
		{"max valid not required", 24576, false, false},
		{"typical not required", 128, false, false},
		{"unset not required", 0, false, false},
		{"min valid required", 16, true, false},
		{"max valid required", 24576, true, false},
		{"unset required", 0, true, true},
		{"below min", 15, false, true},
		{"above max", 24577, false, true},
		{"negative", -1, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateProvisionedMemory(tt.mem, tt.required); (err != nil) != tt.wantErr {
				t.Errorf("validateProvisionedMemory(%d, %v) error = %v, wantErr %v", tt.mem, tt.required, err, tt.wantErr)
			}
		})
	}
}

func TestValidateQueryLanguage(t *testing.T) {
	tests := []struct {
		name    string
		lang    string
		wantErr bool
	}{
		{"OPEN_CYPHER", "OPEN_CYPHER", false},
		{"empty (required)", "", true},
		{"CYPHER (invalid)", "CYPHER", true},
		{"OPENCYPHER (invalid)", "OPENCYPHER", true},
		{"GREMLIN (invalid)", "GREMLIN", true},
		{"lowercase", "open_cypher", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateQueryLanguage(tt.lang); (err != nil) != tt.wantErr {
				t.Errorf("validateQueryLanguage(%q) error = %v, wantErr %v", tt.lang, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePlanCache(t *testing.T) {
	tests := []struct {
		name    string
		pc      string
		wantErr bool
	}{
		{"empty (optional)", "", false},
		{"ENABLED", "ENABLED", false},
		{"DISABLED", "DISABLED", false},
		{"AUTO", "AUTO", false},
		{"invalid", "FOREVER", true},
		{"lowercase", "enabled", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePlanCache(tt.pc); (err != nil) != tt.wantErr {
				t.Errorf("validatePlanCache(%q) error = %v, wantErr %v", tt.pc, err, tt.wantErr)
			}
		})
	}
}

func TestValidateExportFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"CSV", "CSV", false},
		{"PARQUET", "PARQUET", false},
		{"empty (required)", "", true},
		{"CSV+BINARY (invalid)", "CSV+BINARY", true},
		{"JSON (invalid)", "JSON", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateExportFormat(tt.format); (err != nil) != tt.wantErr {
				t.Errorf("validateExportFormat(%q) error = %v, wantErr %v", tt.format, err, tt.wantErr)
			}
		})
	}
}

func TestValidateImportFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"empty (optional)", "", false},
		{"CSV", "CSV", false},
		{"OPEN_CYPHER", "OPEN_CYPHER", false},
		{"PARQUET", "PARQUET", false},
		{"NTRIPLES", "NTRIPLES", false},
		{"json (invalid)", "json", true},
		{"nquad (invalid)", "nquad", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateImportFormat(tt.format); (err != nil) != tt.wantErr {
				t.Errorf("validateImportFormat(%q) error = %v, wantErr %v", tt.format, err, tt.wantErr)
			}
		})
	}
}

func TestValidateVectorSearchDimension(t *testing.T) {
	tests := []struct {
		name    string
		dim     int
		wantErr bool
	}{
		{"min valid", 1, false},
		{"max valid 65536", 65536, false},
		{"typical", 128, false},
		{"zero", 0, true},
		{"above max 65537", 65537, true},
		{"negative", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateVectorSearchDimension(tt.dim); (err != nil) != tt.wantErr {
				t.Errorf("validateVectorSearchDimension(%d) error = %v, wantErr %v", tt.dim, err, tt.wantErr)
			}
		})
	}
}

func TestValidateReplicaCount(t *testing.T) {
	tests := []struct {
		name    string
		rc      int
		wantErr bool
	}{
		{"zero", 0, false},
		{"one", 1, false},
		{"two (max)", 2, false},
		{"three (over max)", 3, true},
		{"negative", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateReplicaCount(tt.rc); (err != nil) != tt.wantErr {
				t.Errorf("validateReplicaCount(%d) error = %v, wantErr %v", tt.rc, err, tt.wantErr)
			}
		})
	}
}

func TestValidateKmsKeyArn(t *testing.T) {
	validArn := "arn:aws:kms:us-east-1:123456789012:key/abcd1234-5678-90ef-abcd-1234567890ab"
	tests := []struct {
		name     string
		arn      string
		required bool
		wantErr  bool
	}{
		{"valid not required", validArn, false, false},
		{"valid required", validArn, true, false},
		{"empty not required", "", false, false},
		{"empty required", "", true, true},
		{"invalid format", "not-an-arn", false, true},
		{"too long", "x" + validArn[:1] + validArn[1:], false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateKmsKeyArn(tt.arn, tt.required); (err != nil) != tt.wantErr {
				t.Errorf("validateKmsKeyArn(%q, %v) error = %v, wantErr %v", tt.arn, tt.required, err, tt.wantErr)
			}
		})
	}
}

func TestValidateRoleArn(t *testing.T) {
	validArn := "arn:aws:iam::123456789012:role/my-neptune-role"
	tests := []struct {
		name    string
		arn     string
		wantErr bool
	}{
		{"valid", validArn, false},
		{"valid service-role", "arn:aws:iam::123456789012:role/service-role/my-role", false},
		{"empty (required)", "", true},
		{"invalid format", "not-an-arn", true},
		{"wrong service", "arn:aws:s3::123456789012:role/my-role", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRoleArn(tt.arn); (err != nil) != tt.wantErr {
				t.Errorf("validateRoleArn(%q) error = %v, wantErr %v", tt.arn, err, tt.wantErr)
			}
		})
	}
}

func TestValidateExplainMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		wantErr bool
	}{
		{"empty (optional)", "", false},
		{"STATIC", "STATIC", false},
		{"DETAILS", "DETAILS", false},
		{"invalid", "FULL", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateExplainMode(tt.mode); (err != nil) != tt.wantErr {
				t.Errorf("validateExplainMode(%q) error = %v, wantErr %v", tt.mode, err, tt.wantErr)
			}
		})
	}
}

func TestValidateParquetType(t *testing.T) {
	tests := []struct {
		name    string
		pt      string
		wantErr bool
	}{
		{"empty (optional)", "", false},
		{"COLUMNAR", "COLUMNAR", false},
		{"invalid", "ROW", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateParquetType(tt.pt); (err != nil) != tt.wantErr {
				t.Errorf("validateParquetType(%q) error = %v, wantErr %v", tt.pt, err, tt.wantErr)
			}
		})
	}
}

func TestValidateBlankNodeHandling(t *testing.T) {
	tests := []struct {
		name    string
		bnh     string
		wantErr bool
	}{
		{"empty (optional)", "", false},
		{"convertToIri", "convertToIri", false},
		{"invalid", "discard", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBlankNodeHandling(tt.bnh); (err != nil) != tt.wantErr {
				t.Errorf("validateBlankNodeHandling(%q) error = %v, wantErr %v", tt.bnh, err, tt.wantErr)
			}
		})
	}
}

func TestValidateGraphSummaryMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		wantErr bool
	}{
		{"empty (optional)", "", false},
		{"BASIC", "BASIC", false},
		{"DETAILED", "DETAILED", false},
		{"invalid", "FULL", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateGraphSummaryMode(tt.mode); (err != nil) != tt.wantErr {
				t.Errorf("validateGraphSummaryMode(%q) error = %v, wantErr %v", tt.mode, err, tt.wantErr)
			}
		})
	}
}

func TestValidateQueryStateInput(t *testing.T) {
	tests := []struct {
		name    string
		state   string
		wantErr bool
	}{
		{"empty (optional)", "", false},
		{"ALL", "ALL", false},
		{"RUNNING", "RUNNING", false},
		{"WAITING", "WAITING", false},
		{"CANCELLING", "CANCELLING", false},
		{"invalid", "COMPLETED", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateQueryStateInput(tt.state); (err != nil) != tt.wantErr {
				t.Errorf("validateQueryStateInput(%q) error = %v, wantErr %v", tt.state, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDestination(t *testing.T) {
	tests := []struct {
		name    string
		dest    string
		wantErr bool
	}{
		{"valid", "file:///tmp/export", false},
		{"empty (required)", "", true},
		{"too long", string(make([]byte, 1025)), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDestination(tt.dest); (err != nil) != tt.wantErr {
				t.Errorf("validateDestination(%q) error = %v, wantErr %v", tt.dest, err, tt.wantErr)
			}
		})
	}
}

func TestValidateSnapshotName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "mysnap", false},
		{"valid with hyphens", "my-snapshot", false},
		{"empty", "", true},
		{"starts with gs-", "gs-mysnap", true},
		{"too long", "aabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklkl", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateSnapshotName(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validateSnapshotName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestClampMaxResults(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{"zero defaults to 100", 0, 100},
		{"negative defaults to 100", -1, 100},
		{"one", 1, 1},
		{"fifty", 50, 50},
		{"hundred (max)", 100, 100},
		{"over max clamped", 200, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clampMaxResults(tt.in); got != tt.want {
				t.Errorf("clampMaxResults(%d) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestPlanCacheBasic(t *testing.T) {
	cache := newQueryPlanCache(3, 100*1000*1000)
	cache.put("key1", "plan1")
	cache.put("key2", "plan2")

	if v, ok := cache.get("key1"); !ok || v != "plan1" {
		t.Error("expected key1 to return plan1")
	}

	cache.put("key3", "plan3")
	cache.put("key4", "plan4")

	if _, ok := cache.get("key2"); ok {
		t.Error("expected key2 to be evicted (LRU)")
	}
	if v, ok := cache.get("key4"); !ok || v != "plan4" {
		t.Error("expected key4 to return plan4")
	}
}

func TestPlanCacheTTL(t *testing.T) {
	cache := newQueryPlanCache(10, 1)
	cache.put("key1", "plan1")
	time.Sleep(3 * time.Millisecond)
	if _, ok := cache.get("key1"); ok {
		t.Error("expected key1 to be expired")
	}
}

func TestPlanCacheKeyDeterministic(t *testing.T) {
	params := map[string]any{
		"alpha":   "first",
		"beta":    "second",
		"gamma":   "third",
		"delta":   42,
		"epsilon": true,
	}
	key1 := planCacheKey("graph-1", "MATCH (n) RETURN n", params)
	for i := 0; i < 100; i++ {
		keyN := planCacheKey("graph-1", "MATCH (n) RETURN n", params)
		if keyN != key1 {
			t.Fatalf("iteration %d: key mismatch\n  got:  %s\n  want: %s", i, keyN, key1)
		}
	}
}

func TestPlanCacheKeyGraphIDPrefix(t *testing.T) {
	key := planCacheKey("mygraph", "MATCH (n) RETURN n", nil)
	if !strings.HasPrefix(key, "mygraph:") {
		t.Fatalf("cache key %q does not have expected 'mygraph:' prefix", key)
	}
}

func TestPlanCachePurgeByGraph(t *testing.T) {
	cache := newQueryPlanCache(100, 5*time.Minute)

	k1 := planCacheKey("graph-A", "MATCH (n) RETURN n", map[string]any{"x": 1})
	k2 := planCacheKey("graph-A", "MATCH (m) RETURN m", map[string]any{"y": 2})
	k3 := planCacheKey("graph-B", "MATCH (n) RETURN n", nil)

	cache.put(k1, "plan1")
	cache.put(k2, "plan2")
	cache.put(k3, "plan3")

	cache.purgeByGraph("graph-A")

	if _, ok := cache.get(k1); ok {
		t.Error("graph-A key1 should have been purged")
	}
	if _, ok := cache.get(k2); ok {
		t.Error("graph-A key2 should have been purged")
	}
	if v, ok := cache.get(k3); !ok || v != "plan3" {
		t.Error("graph-B key3 should still exist after purge of graph-A")
	}
}
