package cloudwatchlogs

import (
	"context"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// newQueryDefinitionEnv builds the service and request context the saved
// query definition tests call handlers through.
func newQueryDefinitionEnv(t *testing.T) (*LogsService, *request.RequestContext) {
	t.Helper()
	svc, _ := newTestService(t)
	reqCtx := request.NewRequestContext(context.Background(), svc.storageManager, "000000000000", "us-east-1")
	return svc, reqCtx
}

// Saved query parameters ride the wire as a list of
// name/defaultValue/description members; the bound members reject and the
// round trip keeps every present member while absent members stay absent.
func TestQueryDefinitionParametersListShape(t *testing.T) {
	svc, reqCtx := newQueryDefinitionEnv(t)

	qdRequest := func(params interface{}) map[string]interface{} {
		return map[string]interface{}{
			"name":        "params-qd",
			"queryString": "fields @message | filter level = {{logLevel}}",
			"parameters":  params,
		}
	}

	if _, err := svc.PutQueryDefinition(context.Background(), reqCtx, vocabRequest(qdRequest([]interface{}{
		map[string]interface{}{"name": "logLevel", "defaultValue": "ERROR", "description": "Level to filter on"},
		map[string]interface{}{"name": "region"},
	}))); err != nil {
		t.Fatalf("put with parameters: %v", err)
	}

	defs, _, err := svc.describeQueryDefinitionsCore("", "", "", "us-east-1", 0)
	if err != nil {
		t.Fatalf("describe: %v", err)
	}
	if len(defs) != 1 || len(defs[0].Parameters) != 2 {
		t.Fatalf("parameters round trip = %+v", defs)
	}
	if defs[0].Parameters[0].Name != "logLevel" || defs[0].Parameters[0].DefaultValue != "ERROR" || defs[0].Parameters[0].Description != "Level to filter on" {
		t.Fatalf("first parameter = %+v", defs[0].Parameters[0])
	}
	if defs[0].Parameters[1].Name != "region" || defs[0].Parameters[1].DefaultValue != "" {
		t.Fatalf("bare parameter = %+v", defs[0].Parameters[1])
	}

	resp, err := svc.DescribeQueryDefinitions(context.Background(), reqCtx, vocabRequest(nil))
	if err != nil {
		t.Fatalf("describe handler: %v", err)
	}
	entries := resp.(map[string]interface{})["queryDefinitions"].([]map[string]interface{})
	if len(entries) != 1 {
		t.Fatalf("handler entries = %d", len(entries))
	}
	wire, ok := entries[0]["parameters"].([]map[string]interface{})
	if !ok || len(wire) != 2 {
		t.Fatalf("wire parameters = %#v", entries[0]["parameters"])
	}
	if wire[0]["name"] != "logLevel" || wire[0]["defaultValue"] != "ERROR" || wire[0]["description"] != "Level to filter on" {
		t.Fatalf("wire first parameter = %+v", wire[0])
	}
	if _, present := wire[1]["defaultValue"]; present {
		t.Fatalf("absent default value emitted: %+v", wire[1])
	}

	rows := []struct {
		name   string
		params []logsstore.QueryParameter
	}{
		{"twenty-one parameters", make([]logsstore.QueryParameter, logsstore.MaxQueryParameters+1)},
		{"missing name", []logsstore.QueryParameter{{DefaultValue: "x"}}},
		{"leading digit name", []logsstore.QueryParameter{{Name: "1level"}}},
		{"oversized default value", []logsstore.QueryParameter{{Name: "level", DefaultValue: string(make([]byte, logsstore.MaxQueryParameterDefaultValueLen+1))}}},
		{"oversized description", []logsstore.QueryParameter{{Name: "level", Description: string(make([]byte, logsstore.MaxQueryParameterDescriptionLen+1))}}},
	}
	for _, row := range rows {
		if _, err := svc.putQueryDefinitionCore("params-qd", "fields @message", "CWLI", "", nil, row.params, "us-east-1"); err == nil {
			t.Fatalf("%s: accepted", row.name)
		}
	}
}

// lastModified is epoch seconds — the documented response examples carry
// ten-digit values — and the queryLanguage filter limits the listing to
// the definitions written in that language, with an absent stored
// language counting as the CWLI default.
func TestQueryDefinitionLastModifiedSecondsAndLanguageFilter(t *testing.T) {
	svc, reqCtx := newQueryDefinitionEnv(t)

	if _, err := svc.PutQueryDefinition(context.Background(), reqCtx, vocabRequest(map[string]interface{}{
		"name": "seconds-cwli", "queryString": "fields @message", "queryLanguage": "CWLI",
	})); err != nil {
		t.Fatalf("put CWLI: %v", err)
	}
	if _, err := svc.PutQueryDefinition(context.Background(), reqCtx, vocabRequest(map[string]interface{}{
		"name": "seconds-sql", "queryString": "SELECT * FROM log", "queryLanguage": "SQL",
	})); err != nil {
		t.Fatalf("put SQL: %v", err)
	}
	if _, err := svc.PutQueryDefinition(context.Background(), reqCtx, vocabRequest(map[string]interface{}{
		"name": "seconds-default", "queryString": "fields @message",
	})); err != nil {
		t.Fatalf("put default-language: %v", err)
	}

	all, _, err := svc.describeQueryDefinitionsCore("", "", "", "us-east-1", 0)
	if err != nil {
		t.Fatalf("describe all: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("unfiltered listing = %d definitions", len(all))
	}
	now := time.Now().Unix()
	for _, qd := range all {
		if qd.LastModified < now-60 || qd.LastModified > now+60 {
			t.Fatalf("lastModified %d is not epoch seconds (now=%d)", qd.LastModified, now)
		}
	}

	sqlOnly, _, err := svc.describeQueryDefinitionsCore("", "SQL", "", "us-east-1", 0)
	if err != nil {
		t.Fatalf("describe SQL: %v", err)
	}
	if len(sqlOnly) != 1 || sqlOnly[0].Name != "seconds-sql" {
		t.Fatalf("SQL filter = %+v", sqlOnly)
	}

	cwliOnly, _, err := svc.describeQueryDefinitionsCore("", "CWLI", "", "us-east-1", 0)
	if err != nil {
		t.Fatalf("describe CWLI: %v", err)
	}
	if len(cwliOnly) != 2 {
		t.Fatalf("CWLI filter = %d definitions, want the explicit and the defaulted one", len(cwliOnly))
	}

	resp, err := svc.DescribeQueryDefinitions(context.Background(), reqCtx, vocabRequest(map[string]interface{}{"queryLanguage": "SQL"}))
	if err != nil {
		t.Fatalf("describe handler filtered: %v", err)
	}
	entries := resp.(map[string]interface{})["queryDefinitions"].([]map[string]interface{})
	if len(entries) != 1 || entries[0]["name"] != "seconds-sql" {
		t.Fatalf("handler language filter = %+v", entries)
	}
}

// The update form of PutQueryDefinition names the definition it updates
// ("To update a query definition, specify its queryDefinitionId in your
// request"): an unknown queryDefinitionId is the operation's declared
// ResourceNotFoundException, never a silent create; the member's length
// trait (1-256) bounds the id; and a known id updates in place,
// returning the same id with the request's members alone retained.
func TestPutQueryDefinitionUpdateForm(t *testing.T) {
	svc, reqCtx := newQueryDefinitionEnv(t)

	resp, err := svc.PutQueryDefinition(context.Background(), reqCtx, vocabRequest(map[string]interface{}{
		"name": "update-form", "queryString": "fields @message",
	}))
	if err != nil {
		t.Fatal(err)
	}
	id := resp.(map[string]interface{})["queryDefinitionId"].(string)

	updated, err := svc.PutQueryDefinition(context.Background(), reqCtx, vocabRequest(map[string]interface{}{
		"name": "update-form-b", "queryString": "fields @timestamp", "queryDefinitionId": id,
	}))
	if err != nil {
		t.Fatalf("update known id: %v", err)
	}
	if got := updated.(map[string]interface{})["queryDefinitionId"].(string); got != id {
		t.Fatalf("update returned %q, want the named id %q", got, id)
	}
	defs, _, err := svc.describeQueryDefinitionsCore("", "", "", "us-east-1", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(defs) != 1 || defs[0].Name != "update-form-b" || defs[0].QueryString != "fields @timestamp" {
		t.Fatalf("listing after update = %+v, want the one definition updated in place", defs)
	}

	if _, err := svc.PutQueryDefinition(context.Background(), reqCtx, vocabRequest(map[string]interface{}{
		"name": "update-form-c", "queryString": "fields @message", "queryDefinitionId": "no-such-definition",
	})); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("unknown queryDefinitionId: error = %v, want ResourceNotFoundException", err)
	}
	defs, _, err = svc.describeQueryDefinitionsCore("", "", "", "us-east-1", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(defs) != 1 {
		t.Fatalf("listing after rejected update = %d definitions, want the rejected form to have created nothing", len(defs))
	}

	oversized := strings.Repeat("q", logsstore.MaxQueryDefinitionIdLength+1)
	if _, err := svc.putQueryDefinitionCore("update-form-d", "fields @message", "", oversized, nil, nil, "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("oversized queryDefinitionId: error = %v, want InvalidParameterException", err)
	}
}
