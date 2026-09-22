package cloudwatchlogs

import (
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The field-index document and scan machinery: the policy document's
// vocabulary, its parser (shared by the group-level PutIndexPolicy and
// the account-level FIELD_INDEX_POLICY) and the ingestion-time scan that
// maintains each field's matching-event bounds.

// --- Policy document vocabulary ---

// indexGeneratedFields is the generated-field set the syntax page allows
// as single-@ field indexes: "Of these generated fields, the following
// are supported for use as field indexes."
var indexGeneratedFields = map[string]bool{
	"@logStream": true, "@ingestionTime": true, "@requestId": true,
	"@type": true, "@initDuration": true, "@duration": true,
	"@billedDuration": true, "@memorySize": true, "@maxMemoryUsed": true,
	"@xrayTraceId": true, "@xraySegmentId": true,
}

// defaultFieldIndexes is the DEFAULT index category's fixed listing:
// "CloudWatch Logs provides default field indexes for all log groups in
// the Standard log class", in the documented order.
var defaultFieldIndexes = []string{
	"@logStream", "@aws.region", "@aws.account", "@source.log",
	"@data_source_name", "@data_source_type", "@data_format",
	"traceId", "severityText", "attributes.session.id",
}

// parsedIndexPolicyDocument is the typed form of the policy document:
// "Fields": ["TransactionId"], "FieldsV2": {"RequestId": {"type":
// "FIELD_INDEX"}, "APIName": {"type": "FACET"}}.
type parsedIndexPolicyDocument struct {
	Fields   []string `json:"Fields"`
	FieldsV2 map[string]struct {
		Type string `json:"type"`
	} `json:"FieldsV2"`
}

// parseIndexPolicyDocument validates a policy document and returns its
// field specs in a deterministic order (the Fields array's order, then
// FieldsV2 by name — the document is a JSON object whose member order the
// wire does not preserve). PutIndexPolicy and the account-level
// FIELD_INDEX_POLICY share the document syntax and its restrictions.
func parseIndexPolicyDocument(document string) ([]logsstore.IndexFieldSpec, error) {
	var doc parsedIndexPolicyDocument
	if err := json.Unmarshal([]byte(document), &doc); err != nil {
		return nil, NewLogsError("InvalidParameterException",
			"The policyDocument must be a JSON object with Fields and/or FieldsV2 members", 400)
	}
	inV2 := make(map[string]bool, len(doc.FieldsV2))
	v2Names := make([]string, 0, len(doc.FieldsV2))
	for name, cfg := range doc.FieldsV2 {
		if cfg.Type != "FIELD_INDEX" && cfg.Type != "FACET" {
			return nil, NewLogsError("InvalidParameterException",
				fmt.Sprintf("Unsupported index type %q for field %s. Supported types are FIELD_INDEX and FACET", cfg.Type, name), 400)
		}
		inV2[name] = true
		v2Names = append(v2Names, name)
	}
	sort.Strings(v2Names)
	seen := make(map[string]bool, len(doc.Fields)+len(doc.FieldsV2))
	total := 0
	for _, name := range doc.Fields {
		if inV2[name] {
			return nil, NewLogsError("InvalidParameterException",
				fmt.Sprintf("Field %s appears in both Fields and FieldsV2; field names within Fields and FieldsV2 must be mutually exclusive", name), 400)
		}
		seen[name] = true
		total++
	}
	for _, name := range v2Names {
		seen[name] = true
		total++
	}
	if total == 0 {
		return nil, NewLogsError("InvalidParameterException",
			"The policy document must include at least one field index", 400)
	}
	if total > logsstore.IndexPolicyFieldsMax {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("An index policy can include as many as %d fields; the document carries %d", logsstore.IndexPolicyFieldsMax, total), 400)
	}
	for name := range seen {
		if name == "" || len(name) > logsstore.IndexPolicyFieldNameMax {
			return nil, NewLogsError("InvalidParameterException",
				fmt.Sprintf("Field name %q must be 1 to %d characters", name, logsstore.IndexPolicyFieldNameMax), 400)
		}
		// "To create an index of a custom field in your log groups that
		// starts with @, you must specify the field with an extra @ at
		// the beginning"; the generated fields are the single-@ names.
		if strings.HasPrefix(name, "@") && !strings.HasPrefix(name, "@@") && !indexGeneratedFields[name] {
			return nil, NewLogsError("InvalidParameterException",
				fmt.Sprintf("Field %s is not a supported generated field; specify custom fields that start with @ with an extra @ (for example @@%s)", name, strings.TrimPrefix(name, "@")), 400)
		}
	}
	specs := make([]logsstore.IndexFieldSpec, 0, total)
	for _, name := range doc.Fields {
		specs = append(specs, logsstore.IndexFieldSpec{Name: name, Type: "FIELD_INDEX"})
	}
	for _, name := range v2Names {
		specs = append(specs, logsstore.IndexFieldSpec{Name: name, Type: doc.FieldsV2[name].Type})
	}
	return specs, nil
}

// indexedFieldName resolves the field index name a spec denotes: a
// doubled-@ spec indexes the single-@ custom field; every other spec is
// the name itself ("@logStream" and friends included).
func indexedFieldName(spec string) string {
	if strings.HasPrefix(spec, "@@") {
		return spec[1:]
	}
	return spec
}

// --- Ingestion-time field-index scan ---

// scanIngestedFieldIndexes maintains the group's field-index bounds for
// one ingested batch: the scan runs at ingestion ("The most recent time
// that CloudWatch Logs scanned ingested log events to search for this
// field index" — FieldIndex.lastScanTime), accumulating each watched
// field's matching-event bounds so the DescribeFieldIndexes listing
// reads the record and never walks the events. The message form is the
// transformed copy where one exists — the query plane's view.
func (s *LogsService) scanIngestedFieldIndexes(store *logsstore.Store, groupName, logStream string, events []logsstore.LogEntry, transformed map[string]string) {
	watched := s.watchedFieldIndexFields(store, groupName)
	messages := make([]logsstore.FieldScanMessage, 0, len(events))
	for _, e := range events {
		message := e.Message
		if out, ok := transformed[logsstore.TransformedMessageDigest(e.Timestamp, logStream, e.Message)]; ok {
			message = out
		}
		messages = append(messages, logsstore.FieldScanMessage{Timestamp: e.Timestamp, Message: message})
	}
	store.MergeFieldIndexBounds(groupName, watched, messages)
}

// watchedFieldIndexFields names the fields the ingestion scan matches:
// the DEFAULT category's non-identity members (an identity field matches
// every event by construction, so its DEFAULT bounds derive from the
// group's stream records) plus the effective policy's fields — a policy
// may index an identity field, and that CUSTOM row then carries
// policy-scoped accumulated bounds.
func (s *LogsService) watchedFieldIndexFields(store *logsstore.Store, groupName string) []string {
	watched := make([]string, 0, len(defaultFieldIndexes)+logsstore.IndexPolicyFieldsMax)
	for _, name := range defaultFieldIndexes {
		if !logsstore.IndexIdentityFields[name] {
			watched = append(watched, name)
		}
	}
	if eff := s.effectiveIndexPolicy(store, groupName); eff != nil {
		for _, spec := range eff.fields {
			name := indexedFieldName(spec.Name)
			if !slices.Contains(watched, name) {
				watched = append(watched, name)
			}
		}
	}
	return watched
}
