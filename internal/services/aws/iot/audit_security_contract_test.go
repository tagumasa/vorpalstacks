package iot

import (
	"testing"

	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// The required-member negatives below are unreachable through the typed AWS
// SDK (it validates the required members client-side), so the server-side
// rejections are pinned here. Every guard runs before any store access, so
// a nil store is sufficient.

func TestListAuditTasksCoreRequiredTimeRange(t *testing.T) {
	svc := &IoTService{}
	if _, err := svc.listAuditTasksCore(nil, ListAuditTasksInput{}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("no time range: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listAuditTasksCore(nil, ListAuditTasksInput{StartTimeProvided: true}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("startTime only: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listAuditTasksCore(nil, ListAuditTasksInput{EndTimeProvided: true}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("endTime only: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listAuditTasksCore(nil, ListAuditTasksInput{
		StartTimeProvided: true, EndTimeProvided: true, TaskType: "HOURLY",
	}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("off-enum taskType: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listAuditTasksCore(nil, ListAuditTasksInput{
		StartTimeProvided: true, EndTimeProvided: true, TaskStatus: "PAUSED",
	}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("off-enum taskStatus: got %v, want ErrInvalidRequest", err)
	}
}

func TestAuditTaskMatchesFilters(t *testing.T) {
	base := ListAuditTasksInput{StartTime: 100, EndTime: 200, StartTimeProvided: true, EndTimeProvided: true}
	rec := map[string]interface{}{"taskId": "t1", "status": "COMPLETED", "startTime": int64(150)}
	if !auditTaskMatchesFilters(rec, base) {
		t.Fatal("task inside the window must match")
	}
	narrow := base
	narrow.StartTime = 151
	if auditTaskMatchesFilters(rec, narrow) {
		t.Fatal("task before the window start must not match")
	}
	narrow = base
	narrow.EndTime = 149
	if auditTaskMatchesFilters(rec, narrow) {
		t.Fatal("task after the window end must not match")
	}
	if !auditTaskMatchesFilters(rec, withType(base, "ON_DEMAND_AUDIT_TASK")) {
		t.Fatal("record without a type member is an on-demand task")
	}
	if auditTaskMatchesFilters(rec, withType(base, "SCHEDULED_AUDIT_TASK")) {
		t.Fatal("on-demand record must not match the SCHEDULED filter")
	}
	statusFilter := base
	statusFilter.TaskStatus = "COMPLETED"
	if !auditTaskMatchesFilters(rec, statusFilter) {
		t.Fatal("status filter must match the stored status")
	}
	statusFilter.TaskStatus = "FAILED"
	if auditTaskMatchesFilters(rec, statusFilter) {
		t.Fatal("status filter must reject a mismatching status")
	}
}

func withType(in ListAuditTasksInput, taskType string) ListAuditTasksInput {
	in.TaskType = taskType
	return in
}

func TestListAuditFindingsCoreTaskIdTimeRangeExclusive(t *testing.T) {
	svc := &IoTService{}
	if _, err := svc.listAuditFindingsCore(nil, ListAuditFindingsInput{
		TaskID: "task-1", StartTimeProvided: true,
	}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("taskId plus time range: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listAuditFindingsCore(nil, ListAuditFindingsInput{
		TaskID: "task-1", EndTimeProvided: true,
	}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("taskId plus endTime: got %v, want ErrInvalidRequest", err)
	}
}

func TestAuditFindingMatchesFilters(t *testing.T) {
	rec := map[string]interface{}{
		"taskId":             "task-1",
		"checkName":          "CHECK_A",
		"findingTime":        int64(150),
		"resourceIdentifier": map[string]interface{}{"deviceCertificateId": "cert-1"},
	}
	if !auditFindingMatchesFilters(rec, ListAuditFindingsInput{TaskID: "task-1"}) {
		t.Fatal("taskId filter must match")
	}
	if auditFindingMatchesFilters(rec, ListAuditFindingsInput{TaskID: "task-2"}) {
		t.Fatal("taskId filter must reject a mismatch")
	}
	if !auditFindingMatchesFilters(rec, ListAuditFindingsInput{CheckName: "CHECK_A"}) {
		t.Fatal("checkName filter must match")
	}
	byResource := ListAuditFindingsInput{ResourceIdentifier: map[string]interface{}{"deviceCertificateId": "cert-1"}}
	if !auditFindingMatchesFilters(rec, byResource) {
		t.Fatal("resourceIdentifier member filter must match")
	}
	byResource.ResourceIdentifier = map[string]interface{}{"deviceCertificateId": "cert-2"}
	if auditFindingMatchesFilters(rec, byResource) {
		t.Fatal("resourceIdentifier member filter must reject a mismatch")
	}
	window := ListAuditFindingsInput{StartTime: 100, EndTime: 200, StartTimeProvided: true, EndTimeProvided: true}
	if !auditFindingMatchesFilters(rec, window) {
		t.Fatal("finding inside the time range must match")
	}
	window.StartTime = 151
	if auditFindingMatchesFilters(rec, window) {
		t.Fatal("finding before the range start must not match")
	}
	suppressed := ListAuditFindingsInput{ListSuppressed: boolPtr(false)}
	if !auditFindingMatchesFilters(rec, suppressed) {
		t.Fatal("a record without isSuppressed counts as not suppressed")
	}
	rec["isSuppressed"] = true
	if auditFindingMatchesFilters(rec, suppressed) {
		t.Fatal("explicit false selector must exclude suppressed findings")
	}
	if !auditFindingMatchesFilters(rec, ListAuditFindingsInput{ListSuppressed: boolPtr(true)}) {
		t.Fatal("explicit true selector must include suppressed findings")
	}
	if !auditFindingMatchesFilters(rec, ListAuditFindingsInput{}) {
		t.Fatal("omitted selector must return both kinds")
	}
}

func TestSuppressionFilterAndOrder(t *testing.T) {
	early := &AuditSuppressionRecord{CheckName: "early", ExpirationDate: int64(100)}
	late := &AuditSuppressionRecord{CheckName: "late", ExpirationDate: int64(200)}
	indefinite := &AuditSuppressionRecord{CheckName: "indefinite", ExpirationDate: nil, SuppressIndefinitely: true}

	filter := ListAuditSuppressionsInput{CheckName: "early"}
	if !suppressionMatchesFilter(map[string]interface{}{"checkName": "early"}, filter) {
		t.Fatal("checkName filter must match")
	}
	if suppressionMatchesFilter(map[string]interface{}{"checkName": "other"}, filter) {
		t.Fatal("checkName filter must reject a mismatch")
	}
	byResource := ListAuditSuppressionsInput{ResourceIdentifier: map[string]interface{}{"caCertificateId": "ca-1"}}
	stored := map[string]interface{}{
		"checkName":          "x",
		"resourceIdentifier": map[string]interface{}{"caCertificateId": "ca-1", "deviceCertificateId": "d-1"},
	}
	if !suppressionMatchesFilter(stored, byResource) {
		t.Fatal("resourceIdentifier filter must match the named members")
	}
	byResource.ResourceIdentifier["caCertificateId"] = "ca-2"
	if suppressionMatchesFilter(stored, byResource) {
		t.Fatal("resourceIdentifier filter must reject a member mismatch")
	}

	entries := []*AuditSuppressionRecord{late, indefinite, early}
	sortAuditSuppressions(entries, nil)
	if entries[0] != early || entries[1] != late || entries[2] != indefinite {
		t.Fatal("default order must be ascending by expiration date with indefinite last")
	}
	entries = []*AuditSuppressionRecord{early, late}
	sortAuditSuppressions(entries, boolPtr(false))
	if entries[0] != late || entries[1] != early {
		t.Fatal("explicit false must sort descending by expiration date")
	}
}

func TestUpdateAuditSuppressionCoreExclusivity(t *testing.T) {
	svc := &IoTService{}
	// An expiration date together with an indefinite suppression is the
	// contradictory pairing; an explicit false alongside a date is the
	// documented CLI pairing and accepted.
	err := svc.updateAuditSuppressionCore(nil, UpdateAuditSuppressionInput{
		CheckName:                    "c",
		ResourceIdentifier:           map[string]interface{}{"deviceCertificateId": "d"},
		ExpirationProvided:           true,
		SuppressIndefinitelyProvided: true,
		SuppressIndefinitely:         true,
	})
	if err != iotstore.ErrInvalidRequest {
		t.Fatalf("date + indefinite: got %v, want ErrInvalidRequest", err)
	}
}

func TestListSecurityProfilesCoreExclusiveFilters(t *testing.T) {
	svc := &IoTService{}
	_, _, err := svc.listSecurityProfilesCore(nil, storecommon.ListOptions{}, ListSecurityProfilesInput{
		DimensionName: "dim", MetricName: "metric",
	})
	if err != iotstore.ErrInvalidRequest {
		t.Fatalf("both filters: got %v, want ErrInvalidRequest", err)
	}
}

func TestSecurityProfileMatchesFilter(t *testing.T) {
	byMetric := &iotstore.SecurityProfile{Behaviors: []*iotstore.Behavior{
		{Name: "b", Metric: "aws:num-connected-devices"},
	}}
	byDimension := &iotstore.SecurityProfile{Behaviors: []*iotstore.Behavior{
		{Name: "b", MetricDimension: "fleet-dim"},
	}}
	byRetained := &iotstore.SecurityProfile{AdditionalMetricsToRetainV2: []*iotstore.MetricToRetain{
		{Metric: "custom-metric", MetricDimension: "retained-dim"},
	}}
	if !securityProfileMatchesFilter(byMetric, ListSecurityProfilesInput{MetricName: "aws:num-connected-devices"}) {
		t.Fatal("behaviour metric filter must match")
	}
	if securityProfileMatchesFilter(byMetric, ListSecurityProfilesInput{DimensionName: "fleet-dim"}) {
		t.Fatal("metric-only profile must not match a dimension filter")
	}
	if !securityProfileMatchesFilter(byDimension, ListSecurityProfilesInput{DimensionName: "fleet-dim"}) {
		t.Fatal("behaviour dimension filter must match")
	}
	if securityProfileMatchesFilter(byDimension, ListSecurityProfilesInput{MetricName: "custom-metric"}) {
		t.Fatal("dimension-only profile must not match a metric filter")
	}
	if !securityProfileMatchesFilter(byRetained, ListSecurityProfilesInput{MetricName: "custom-metric"}) {
		t.Fatal("retained metric filter must match")
	}
	if !securityProfileMatchesFilter(byRetained, ListSecurityProfilesInput{DimensionName: "retained-dim"}) {
		t.Fatal("retained dimension filter must match")
	}
}

func TestListViolationEventsCoreRequiredTimeRange(t *testing.T) {
	svc := &IoTService{}
	if _, err := svc.listViolationEventsCore(nil, storecommon.ListOptions{}, ListViolationEventsInput{}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("no time range: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listViolationEventsCore(nil, storecommon.ListOptions{}, ListViolationEventsInput{
		StartTimeProvided: true,
	}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("startTime only: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listActiveViolationsCore(nil, ListActiveViolationsInput{
		BehaviorCriteriaType: "NEURAL",
	}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("off-enum behaviorCriteriaType: got %v, want ErrInvalidRequest", err)
	}
}

func TestBehaviorCriteriaTypeOf(t *testing.T) {
	if behaviorCriteriaTypeOf(nil) != "" || behaviorCriteriaTypeOf(&iotstore.Behavior{}) != "" {
		t.Fatal("behaviour without criteria has no type")
	}
	static := &iotstore.Behavior{Criteria: &iotstore.BehaviorCriteria{ComparisonOperator: "less-than"}}
	if behaviorCriteriaTypeOf(static) != "STATIC" {
		t.Fatalf("plain criteria: got %q, want STATIC", behaviorCriteriaTypeOf(static))
	}
	statistical := &iotstore.Behavior{Criteria: &iotstore.BehaviorCriteria{
		StatisticalThreshold: &iotstore.StatisticalThreshold{Statistic: "p99"},
	}}
	if behaviorCriteriaTypeOf(statistical) != "STATISTICAL" {
		t.Fatalf("statistical criteria: got %q, want STATISTICAL", behaviorCriteriaTypeOf(statistical))
	}
	ml := &iotstore.Behavior{Criteria: &iotstore.BehaviorCriteria{
		MLDetectionConfig: &iotstore.MachineLearningDetectionConfig{ConfidenceLevel: "HIGH"},
	}}
	if behaviorCriteriaTypeOf(ml) != "MACHINE_LEARNING" {
		t.Fatalf("ML criteria: got %q, want MACHINE_LEARNING", behaviorCriteriaTypeOf(ml))
	}
}

func TestViolationMatchesFilters(t *testing.T) {
	violation := &iotstore.ViolationEvent{
		ViolationID:         "v1",
		SecurityProfileName: "profile-a",
		VerificationState:   "TRUE_POSITIVE",
		Behavior:            &iotstore.Behavior{Name: "b", Criteria: &iotstore.BehaviorCriteria{}},
	}
	if !violationMatchesFilters(violation, "profile-a", "STATIC", "TRUE_POSITIVE", nil) {
		t.Fatal("matching filters must accept the violation")
	}
	if violationMatchesFilters(violation, "profile-b", "", "", nil) {
		t.Fatal("profile mismatch must reject")
	}
	if violationMatchesFilters(violation, "", "STATISTICAL", "", nil) {
		t.Fatal("criteria type mismatch must reject")
	}
	if violationMatchesFilters(violation, "", "", "FALSE_POSITIVE", nil) {
		t.Fatal("verification state mismatch must reject")
	}
	suppressed := &iotstore.ViolationEvent{
		Behavior: &iotstore.Behavior{Name: "b", SuppressAlerts: true},
	}
	if violationMatchesFilters(suppressed, "", "", "", nil) {
		t.Fatal("suppressed alert must be excluded by default")
	}
	if violationMatchesFilters(suppressed, "", "", "", boolPtr(false)) {
		t.Fatal("explicit false selector must exclude suppressed alerts")
	}
	if !violationMatchesFilters(suppressed, "", "", "", boolPtr(true)) {
		t.Fatal("explicit true selector must include suppressed alerts")
	}
}

func TestRawToBehaviorsStructuredMembers(t *testing.T) {
	behaviors, err := rawToBehaviors([]map[string]interface{}{
		{
			"name":            "b1",
			"metricDimension": map[string]interface{}{"dimensionName": "fleet-dim", "operator": "IN"},
			"criteria": map[string]interface{}{
				"comparisonOperator":   "less-than",
				"statisticalThreshold": map[string]interface{}{"statistic": "p90"},
				"mlDetectionConfig":    map[string]interface{}{"confidenceLevel": "HIGH"},
			},
		},
	})
	if err != nil {
		t.Fatalf("rawToBehaviors failed: %v", err)
	}
	if len(behaviors) != 1 {
		t.Fatalf("got %d behaviours, want 1", len(behaviors))
	}
	b := behaviors[0]
	if b.MetricDimension != "fleet-dim" {
		t.Fatalf("metricDimension: got %q, want fleet-dim", b.MetricDimension)
	}
	if b.Criteria == nil || b.Criteria.StatisticalThreshold == nil || b.Criteria.StatisticalThreshold.Statistic != "p90" {
		t.Fatal("statisticalThreshold member was dropped")
	}
	if b.Criteria.MLDetectionConfig == nil || b.Criteria.MLDetectionConfig.ConfidenceLevel != "HIGH" {
		t.Fatal("mlDetectionConfig member was dropped")
	}
}

func TestSetV2LoggingOptionsCoreValidation(t *testing.T) {
	svc := &IoTService{}
	err := svc.setV2LoggingOptionsCore(nil, SetV2LoggingOptionsInput{
		EventConfigurations: []interface{}{
			map[string]interface{}{"eventType": "Connect", "logLevel": "LOUD"},
		},
	})
	if err != iotstore.ErrInvalidRequest {
		t.Fatalf("off-enum logLevel: got %v, want ErrInvalidRequest", err)
	}
	err = svc.setV2LoggingOptionsCore(nil, SetV2LoggingOptionsInput{
		EventConfigurations: []interface{}{
			map[string]interface{}{"logLevel": "WARN"},
		},
	})
	if err != iotstore.ErrMissingParam {
		t.Fatalf("missing eventType: got %v, want ErrMissingParam", err)
	}
	err = svc.setV2LoggingOptionsCore(nil, SetV2LoggingOptionsInput{DefaultLogLevel: "LOUD"})
	if err != iotstore.ErrInvalidRequest {
		t.Fatalf("off-enum defaultLogLevel: got %v, want ErrInvalidRequest", err)
	}
}

func TestCreateRequiresClientRequestToken(t *testing.T) {
	svc := &IoTService{}
	if _, err := svc.createCustomMetricCore(nil, CreateCustomMetricInput{
		MetricName: "m", MetricType: "string-list",
	}); err != iotstore.ErrMissingParam {
		t.Fatalf("custom metric without token: got %v, want ErrMissingParam", err)
	}
	if _, err := svc.createDimensionCore(nil, CreateDimensionInput{
		Name: "d", Type: "TOPIC_FILTER", StringValues: []string{"t/#"},
	}); err != iotstore.ErrMissingParam {
		t.Fatalf("dimension without token: got %v, want ErrMissingParam", err)
	}
}

func boolPtr(v bool) *bool { return &v }

// TestListAuditFindingsRequiredFilterSelection pins the documented filter
// requirement: either the taskId or the full startTime/endTime pair must
// be specified, but never both. The guards run before any store access,
// so a nil store is sufficient.
func TestListAuditFindingsRequiredFilterSelection(t *testing.T) {
	svc := &IoTService{}
	if _, err := svc.listAuditFindingsCore(nil, ListAuditFindingsInput{}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("no filter: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listAuditFindingsCore(nil, ListAuditFindingsInput{CheckName: "DEVICE_CERTIFICATE_EXPIRING_CHECK"}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("checkName only: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listAuditFindingsCore(nil, ListAuditFindingsInput{StartTimeProvided: true}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("startTime only: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listAuditFindingsCore(nil, ListAuditFindingsInput{EndTimeProvided: true}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("endTime only: got %v, want ErrInvalidRequest", err)
	}
	if _, err := svc.listAuditFindingsCore(nil, ListAuditFindingsInput{
		TaskID: "audit-task-1", StartTimeProvided: true, EndTimeProvided: true,
	}); err != iotstore.ErrInvalidRequest {
		t.Fatalf("taskId plus time range: got %v, want ErrInvalidRequest", err)
	}
}
