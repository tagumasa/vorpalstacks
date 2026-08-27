package iot

import (
	"errors"
	"strings"
	"testing"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

func strPtr(s string) *string { return &s }
func i64Ptr(v int64) *int64   { return &v }

// validateFleetMetricCreate runs entirely before store access, so the
// rejection paths are safe to exercise with a nil store.
func TestValidateFleetMetricCreate(t *testing.T) {
	valid := FleetMetricInput{
		MetricName:       "m",
		QueryString:      strPtr("*"),
		AggregationField: strPtr("thingName"),
		AggregationType:  map[string]interface{}{"name": "Statistics"},
		Period:           i64Ptr(60),
	}
	if err := validateFleetMetricCreate(valid); err != nil {
		t.Fatalf("valid input rejected: %v", err)
	}

	missing := []struct {
		name string
		mut  func(in *FleetMetricInput)
	}{
		{"queryString", func(in *FleetMetricInput) { in.QueryString = strPtr("") }},
		{"queryString omitted", func(in *FleetMetricInput) { in.QueryString = nil }},
		{"aggregationField", func(in *FleetMetricInput) { in.AggregationField = strPtr("") }},
		{"aggregationType name", func(in *FleetMetricInput) { in.AggregationType = map[string]interface{}{} }},
		{"period omitted", func(in *FleetMetricInput) { in.Period = nil }},
	}
	for _, tt := range missing {
		t.Run("missing "+tt.name, func(t *testing.T) {
			in := valid
			tt.mut(&in)
			if err := validateFleetMetricCreate(in); !errors.Is(err, iotstore.ErrMissingParam) {
				t.Fatalf("expected ErrMissingParam, got %v", err)
			}
		})
	}

	badPeriods := []struct {
		name   string
		period int64
	}{
		{"below minimum", 30},
		{"not a multiple of 60", 61},
		{"above maximum", MaxFleetMetricPeriod + 60},
	}
	for _, tt := range badPeriods {
		t.Run("period "+tt.name, func(t *testing.T) {
			in := valid
			in.Period = i64Ptr(tt.period)
			if err := validateFleetMetricCreate(in); !errors.Is(err, iotstore.ErrInvalidRequest) {
				t.Fatalf("expected ErrInvalidRequest, got %v", err)
			}
		})
	}
}

// updateFleetMetricCore's indexName requirement and period validation run
// before store access.
func TestUpdateFleetMetricRequiresIndexName(t *testing.T) {
	svc := &IoTService{}
	in := FleetMetricInput{MetricName: "m", QueryString: strPtr("*"), Period: i64Ptr(60), AggregationType: map[string]interface{}{"name": "Statistics"}, AggregationField: strPtr("thingName")}
	if _, err := svc.updateFleetMetricCore(nil, in); !errors.Is(err, iotstore.ErrMissingParam) {
		t.Fatalf("expected ErrMissingParam for missing indexName, got %v", err)
	}
	in.IndexName = "AWS_Things"
	in.Period = i64Ptr(61)
	if _, err := svc.updateFleetMetricCore(nil, in); !errors.Is(err, iotstore.ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest for period 61, got %v", err)
	}
}

// fleetMetricFields must carry only the provided members so the update
// merge preserves omitted optional members; an explicitly provided empty
// string, by contrast, is a deliberate overwrite and must survive.
func TestFleetMetricFieldsOmitAbsentMembers(t *testing.T) {
	in := FleetMetricInput{
		MetricName:  "m",
		IndexName:   "AWS_Things",
		QueryString: strPtr("thingName:*"),
		Period:      i64Ptr(120),
	}
	fields := fleetMetricFields(in)
	if _, ok := fields["aggregationType"]; ok {
		t.Fatal("omitted aggregationType must not be written")
	}
	if _, ok := fields["aggregationField"]; ok {
		t.Fatal("omitted aggregationField must not be written")
	}
	if _, ok := fields["unit"]; ok {
		t.Fatal("omitted unit must not be written")
	}
	if fields["period"] != int64(120) {
		t.Fatalf("provided period must be carried, got %v", fields["period"])
	}

	in.QueryString = strPtr("")
	fields = fleetMetricFields(in)
	if v, ok := fields["queryString"].(string); !ok || v != "" {
		t.Fatal("an explicitly provided empty queryString is a deliberate overwrite")
	}
}

// validateFleetMetricMembers applies the Smithy shape constraints to every
// provided member; aggregation violations map to InvalidAggregationException
// and everything else to InvalidRequestException. It runs before store
// access, so a nil store is safe here.
func TestValidateFleetMetricMembers(t *testing.T) {
	valid := FleetMetricInput{
		MetricName:      "m",
		IndexName:       "AWS_Things",
		QueryString:     strPtr("thingName:*"),
		AggregationType: map[string]interface{}{"name": "Statistics", "values": []interface{}{"avg", "count"}},
		Unit:            strPtr("Seconds"),
		Description:     strPtr("fleet metric"),
	}
	if err := validateFleetMetricMembers(valid); err != nil {
		t.Fatalf("valid input rejected: %v", err)
	}

	invalidRequest := []struct {
		name string
		mut  func(in *FleetMetricInput)
	}{
		{"name pattern", func(in *FleetMetricInput) { in.MetricName = "bad name!" }},
		{"name over length", func(in *FleetMetricInput) { in.MetricName = strings.Repeat("a", MaxFleetMetricNameLength+1) }},
		{"empty queryString", func(in *FleetMetricInput) { in.QueryString = strPtr("") }},
		{"empty aggregationField", func(in *FleetMetricInput) { in.AggregationField = strPtr("") }},
		{"description over length", func(in *FleetMetricInput) {
			in.Description = strPtr(strings.Repeat("a", MaxFleetMetricDescriptionLength+1))
		}},
		{"description control character", func(in *FleetMetricInput) { in.Description = strPtr("line\nbreak") }},
		{"indexName pattern", func(in *FleetMetricInput) { in.IndexName = "AWS Things!" }},
		{"unit enum", func(in *FleetMetricInput) { in.Unit = strPtr("SolarFlux") }},
	}
	for _, tt := range invalidRequest {
		t.Run(tt.name, func(t *testing.T) {
			in := valid
			tt.mut(&in)
			if err := validateFleetMetricMembers(in); !errors.Is(err, iotstore.ErrInvalidRequest) {
				t.Fatalf("expected ErrInvalidRequest, got %v", err)
			}
		})
	}

	invalidAggregation := []struct {
		name string
		mut  func(in *FleetMetricInput)
	}{
		{"name enum", func(in *FleetMetricInput) {
			in.AggregationType = map[string]interface{}{"name": "Bogus"}
		}},
		{"value pattern", func(in *FleetMetricInput) {
			in.AggregationType = map[string]interface{}{"name": "Statistics", "values": []interface{}{"p50!"}}
		}},
		{"value over length", func(in *FleetMetricInput) {
			in.AggregationType = map[string]interface{}{"name": "Statistics", "values": []interface{}{strings.Repeat("a", MaxAggregationValueLength+1)}}
		}},
	}
	for _, tt := range invalidAggregation {
		t.Run("aggregation "+tt.name, func(t *testing.T) {
			in := valid
			tt.mut(&in)
			if err := validateFleetMetricMembers(in); !errors.Is(err, iotstore.ErrInvalidAggregation) {
				t.Fatalf("expected ErrInvalidAggregation, got %v", err)
			}
		})
	}
}
