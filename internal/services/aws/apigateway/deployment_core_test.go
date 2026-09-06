package apigateway

import (
	"testing"

	"vorpalstacks/internal/store/aws/apigateway"
)

// TestApplyMapPatchPointerUnescaping pins the RFC 6901 unescaping of map keys
// extracted from patch paths: a slash travels as ~1 and a tilde as ~0, and
// applyMapPatch must store the unescaped key on every operation kind.
func TestApplyMapPatchPointerUnescaping(t *testing.T) {
	vars := map[string]string{"existing": "old"}

	if err := applyMapPatch(vars, PatchOperation{Op: "replace", Path: "/variables/a~1b", Value: "v1"}, "/variables/", nil, nil); err != nil {
		t.Fatalf("replace failed: %v", err)
	}
	if vars["a/b"] != "v1" {
		t.Fatalf("escaped slash key not unescaped: %v", vars)
	}
	if _, ok := vars["a~1b"]; ok {
		t.Fatal("escaped form stored as a literal key")
	}

	if err := applyMapPatch(vars, PatchOperation{Op: "add", Path: "/variables/tilde~0key", Value: "v2"}, "/variables/", nil, nil); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if vars["tilde~key"] != "v2" {
		t.Fatalf("escaped tilde key not unescaped: %v", vars)
	}

	if err := applyMapPatch(vars, PatchOperation{Op: "remove", Path: "/variables/a~1b"}, "/variables/", nil, nil); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if _, ok := vars["a/b"]; ok {
		t.Fatal("remove did not delete the unescaped key")
	}

	if err := applyMapPatch(vars, PatchOperation{Op: "replace", Path: "/variables/plain", Value: "v3"}, "/variables/", nil, nil); err != nil {
		t.Fatalf("plain replace failed: %v", err)
	}
	if vars["plain"] != "v3" {
		t.Fatalf("plain key mishandled: %v", vars)
	}
}

// TestValidateStageVariables pins the ingress validation of a whole stage
// variables map: legal names and values pass; charset violations and empty
// values reject.
func TestValidateStageVariables(t *testing.T) {
	if err := validateStageVariables(map[string]string{"env": "prod", "tier_2": "a-z.~:/?#&=,"}); err != nil {
		t.Fatalf("valid variables rejected: %v", err)
	}
	for name, vars := range map[string]map[string]string{
		"invalid name":  {"has space": "v"},
		"invalid value": {"env": "has space"},
		"empty value":   {"env": ""},
	} {
		if err := validateStageVariables(vars); err == nil {
			t.Fatalf("%s accepted: %v", name, vars)
		}
	}
}

// TestSplitApiStageValue pins the apiId:stageName identifier split the
// documented /apiStages addressing uses.
func TestSplitApiStageValue(t *testing.T) {
	if api, stage, ok := splitApiStageValue("abc123:Prod"); !ok || api != "abc123" || stage != "Prod" {
		t.Fatalf("split failed: %q %q %v", api, stage, ok)
	}
	for _, bad := range []string{"", "abc123", ":Prod", "abc123:", "a:b:c"} {
		if _, _, ok := splitApiStageValue(bad); ok {
			t.Fatalf("%q accepted", bad)
		}
	}
}

// TestApplyApiStageThrottlePatch pins the documented per-stage throttle
// patch family: the whole-throttle JSON form keyed by
// {resourcePath}/{httpMethod}, the single-method rateLimit/burstLimit
// members, and the remove semantics of each row.
func TestApplyApiStageThrottlePatch(t *testing.T) {
	stage := &apigateway.ApiStage{}

	if err := applyApiStageThrottlePatch(stage, PatchOperation{
		Op:    "replace",
		Path:  "/apiStages/abc123:Prod/throttle",
		Value: `{"//GET":{"rateLimit":1,"burstLimit":2},"/pets/GET":{"rateLimit":0.5}}`,
	}); err != nil {
		t.Fatalf("whole replace failed: %v", err)
	}
	if len(stage.Throttle) != 2 || stage.Throttle["//GET"].RateLimit != 1 || stage.Throttle["//GET"].BurstLimit != 2 {
		t.Fatalf("whole replace not applied: %+v", stage.Throttle)
	}

	if err := applyApiStageThrottlePatch(stage, PatchOperation{
		Op:    "add",
		Path:  "/apiStages/abc123:Prod/throttle//GET/rateLimit",
		Value: "0.25",
	}); err != nil {
		t.Fatalf("single rateLimit failed: %v", err)
	}
	if stage.Throttle["//GET"].RateLimit != 0.25 {
		t.Fatalf("single rateLimit not applied: %+v", stage.Throttle["//GET"])
	}

	if err := applyApiStageThrottlePatch(stage, PatchOperation{
		Op:   "remove",
		Path: "/apiStages/abc123:Prod/throttle//GET",
	}); err != nil {
		t.Fatalf("method remove failed: %v", err)
	}
	if _, ok := stage.Throttle["//GET"]; ok {
		t.Fatalf("method remove did not delete the key: %+v", stage.Throttle)
	}

	// The escaped resource path token stays as addressed: the key is
	// "~1items~1{id}/GET", the convention the official CLI update-stage
	// example output shows for the method-keyed maps.
	if err := applyApiStageThrottlePatch(stage, PatchOperation{
		Op:    "add",
		Path:  "/apiStages/abc123:Prod/throttle/~1items~1{id}/GET/rateLimit",
		Value: "3",
	}); err != nil {
		t.Fatalf("escaped-key rateLimit failed: %v", err)
	}
	if stage.Throttle["~1items~1{id}/GET"] == nil || stage.Throttle["~1items~1{id}/GET"].RateLimit != 3 {
		t.Fatalf("escaped-key entry not stored as addressed: %+v", stage.Throttle)
	}
	if err := applyApiStageThrottlePatch(stage, PatchOperation{
		Op:   "remove",
		Path: "/apiStages/abc123:Prod/throttle/~1items~1{id}/GET",
	}); err != nil {
		t.Fatalf("escaped-key remove failed: %v", err)
	}
	if _, ok := stage.Throttle["~1items~1{id}/GET"]; ok {
		t.Fatalf("escaped-key remove did not delete the key: %+v", stage.Throttle)
	}

	if err := applyApiStageThrottlePatch(stage, PatchOperation{
		Op:   "remove",
		Path: "/apiStages/abc123:Prod/throttle",
	}); err != nil {
		t.Fatalf("whole remove failed: %v", err)
	}
	if len(stage.Throttle) != 0 {
		t.Fatalf("whole remove did not clear the map: %+v", stage.Throttle)
	}

	// The method-only row documents remove only.
	if err := applyApiStageThrottlePatch(stage, PatchOperation{
		Op:    "add",
		Path:  "/apiStages/abc123:Prod/throttle//GET",
		Value: "x",
	}); err == nil {
		t.Fatal("add on the method-only form accepted")
	}

	// Element addressing outside the throttle family is not documented.
	if err := applyApiStageThrottlePatch(stage, PatchOperation{
		Op:    "replace",
		Path:  "/apiStages/abc123:Prod/apiId",
		Value: "x",
	}); err == nil {
		t.Fatal("non-throttle element path accepted")
	}
}

// TestRequirePatchOp pins the shared (path, op) gate: operations outside the
// set transcribed from the official patch table reject with the
// unknown-patch-path wording, operations inside it pass.
func TestRequirePatchOp(t *testing.T) {
	if err := requirePatchOp(PatchOperation{Op: "add", Path: "/description"}, opReplace); err == nil {
		t.Fatal("add on a replace-only row was accepted")
	}
	if err := requirePatchOp(PatchOperation{Op: "add", Path: "/description"}, opAdd|opReplace); err != nil {
		t.Fatalf("add on an add+replace row rejected: %v", err)
	}
	if err := requirePatchOp(PatchOperation{Op: "copy", Path: "/deploymentId"}, opReplace|opCopy); err != nil {
		t.Fatalf("copy on a copy-supporting row rejected: %v", err)
	}
	if err := requirePatchOp(PatchOperation{Op: "copy", Path: "/name"}, opReplace); err == nil {
		t.Fatal("copy on a replace-only row was accepted")
	}
}

// TestParsePatchOperationsCopyAndFrom pins the wire parsing of the copy
// operation: copy is admitted only with a from path, and its source travels
// on the from member while move and test remain rejected outright.
func TestParsePatchOperationsCopyAndFrom(t *testing.T) {
	ops, err := parsePatchOperations(map[string]interface{}{
		"patchOperations": []interface{}{
			map[string]interface{}{"op": "copy", "path": "/deploymentId", "from": "/canarySettings/deploymentId"},
		},
	})
	if err != nil {
		t.Fatalf("copy with from rejected: %v", err)
	}
	if len(ops) != 1 || ops[0].From != "/canarySettings/deploymentId" {
		t.Fatalf("from member not parsed: %+v", ops)
	}

	if _, err := parsePatchOperations(map[string]interface{}{
		"patchOperations": []interface{}{
			map[string]interface{}{"op": "copy", "path": "/deploymentId"},
		},
	}); err == nil {
		t.Fatal("copy without from was accepted")
	}

	if _, err := parsePatchOperations(map[string]interface{}{
		"patchOperations": []interface{}{
			map[string]interface{}{"op": "move", "path": "/name", "from": "/description"},
		},
	}); err == nil {
		t.Fatal("move was accepted")
	}
}

// TestApplyMapPatchValidators pins the per-key validator hook: names and
// values run through the validators on add and replace, remove bypasses
// them (it carries no value), and nil validators keep the prior behaviour.
func TestApplyMapPatchValidators(t *testing.T) {
	vars := map[string]string{"existing": "old"}
	badName := func(string) bool { return false }
	badValue := func(string) bool { return false }

	if err := applyMapPatch(vars, PatchOperation{Op: "replace", Path: "/variables/env", Value: "prod"}, "/variables/", badName, nil); err == nil {
		t.Fatal("invalid name accepted")
	}
	if err := applyMapPatch(vars, PatchOperation{Op: "replace", Path: "/variables/env", Value: "prod"}, "/variables/", nil, badValue); err == nil {
		t.Fatal("invalid value accepted")
	}
	if err := applyMapPatch(vars, PatchOperation{Op: "remove", Path: "/variables/existing"}, "/variables/", badName, badValue); err != nil {
		t.Fatalf("remove ran through validators: %v", err)
	}
	if _, ok := vars["existing"]; ok {
		t.Fatal("remove did not delete the entry")
	}
	if err := applyMapPatch(vars, PatchOperation{Op: "replace", Path: "/variables/env", Value: "prod"}, "/variables/", nil, nil); err != nil {
		t.Fatalf("nil validators rejected a plain replace: %v", err)
	}
}

// TestApplyMapPatchEmptyTokenRejected pins the shared guard: a
// trailing-slash path carries an empty key token, which is not a
// documented patch form — every operation kind is rejected and the map
// is left unchanged instead of gaining an empty-string key.
func TestApplyMapPatchEmptyTokenRejected(t *testing.T) {
	vars := map[string]string{"existing": "old"}

	for _, op := range []string{"add", "replace", "remove"} {
		if err := applyMapPatch(vars, PatchOperation{Op: op, Path: "/variables/"}, "/variables/", nil, nil); err == nil {
			t.Fatalf("op %q on an empty key token was accepted", op)
		}
	}
	if len(vars) != 1 || vars["existing"] != "old" {
		t.Fatalf("map mutated by rejected patches: %v", vars)
	}
}

// TestApplyBoolMapPatchEmptyTokenRejected is the bool-map twin of the
// empty-token guard pin.
func TestApplyBoolMapPatchEmptyTokenRejected(t *testing.T) {
	params := map[string]bool{"method.request.header.Authorization": true}

	for _, op := range []string{"add", "replace", "remove"} {
		if err := applyBoolMapPatch(params, PatchOperation{Op: op, Path: "/requestParameters/"}, "/requestParameters/", nil, nil); err == nil {
			t.Fatalf("op %q on an empty key token was accepted", op)
		}
	}
	if len(params) != 1 {
		t.Fatalf("map mutated by rejected patches: %v", params)
	}
}

// TestApplyCanarySettingsPatchVariableOverrideUnescaping pins the canary
// stageVariableOverrides branch: override names carried with the documented
// JSON Pointer escaping arrive unescaped in the stored map.
func TestApplyCanarySettingsPatchOverrideConstraints(t *testing.T) {
	stage := &apigateway.Stage{}

	if err := applyCanarySettingsPatch(stage, PatchOperation{
		Op:    "replace",
		Path:  "/canarySettings/stageVariableOverrides/feature_x",
		Value: "on",
	}); err != nil {
		t.Fatalf("valid override rejected: %v", err)
	}
	if got := stage.CanarySettings.StageVariableOverrides["feature_x"]; got != "on" {
		t.Fatalf("override not stored: %v", stage.CanarySettings.StageVariableOverrides)
	}

	// a~1b unescapes to a/b, which the documented stage variable name
	// charset rejects.
	if err := applyCanarySettingsPatch(stage, PatchOperation{
		Op:    "replace",
		Path:  "/canarySettings/stageVariableOverrides/a~1b",
		Value: "v",
	}); err == nil {
		t.Fatal("charset-invalid override name accepted")
	}
	if _, ok := stage.CanarySettings.StageVariableOverrides["a/b"]; ok {
		t.Fatal("charset-invalid override name stored")
	}

	// The per-key override form admits replace only, mirroring the
	// /variables/* row that governs stage-variable per-key patches.
	for _, op := range []string{"add", "remove"} {
		if err := applyCanarySettingsPatch(stage, PatchOperation{
			Op:   op,
			Path: "/canarySettings/stageVariableOverrides/extra",
		}); err == nil {
			t.Fatalf("op %q on a per-key override accepted", op)
		}
		if _, ok := stage.CanarySettings.StageVariableOverrides["extra"]; ok {
			t.Fatalf("op %q mutated the overrides", op)
		}
	}
	if len(stage.CanarySettings.StageVariableOverrides) != 1 || stage.CanarySettings.StageVariableOverrides["feature_x"] != "on" {
		t.Fatalf("rejected operations mutated the overrides: %v", stage.CanarySettings.StageVariableOverrides)
	}
}

// TestApplyCanarySettingsPatchEmptyOverrideTokenRejected pins the
// trailing-slash canary override path: the empty key token must surface
// as an error through the canary applier, not as an empty-string key.
func TestApplyCanarySettingsPatchEmptyOverrideTokenRejected(t *testing.T) {
	stage := &apigateway.Stage{}
	if err := applyCanarySettingsPatch(stage, PatchOperation{Op: "add", Path: "/canarySettings/stageVariableOverrides/"}); err == nil {
		t.Fatal("empty override key token was accepted")
	}
	if stage.CanarySettings != nil && len(stage.CanarySettings.StageVariableOverrides) != 0 {
		t.Fatalf("overrides mutated by the rejected patch: %v", stage.CanarySettings.StageVariableOverrides)
	}
}

// TestParseWholeStringMapValue pins the whole-member map value decoding:
// the JSON object string form, the non-empty name rule, and the stage
// variable name and value constraints where they apply.
func TestParseWholeStringMapValue(t *testing.T) {
	parsed, err := parseWholeStringMapValue(
		PatchOperation{Op: "replace", Path: "/variables", Value: `{"env":"prod","tier_2":"gold"}`},
		validateStageVariableName, validateStageVariableValue)
	if err != nil {
		t.Fatalf("valid value rejected: %v", err)
	}
	if parsed["env"] != "prod" || parsed["tier_2"] != "gold" || len(parsed) != 2 {
		t.Fatalf("map not decoded: %v", parsed)
	}

	for _, tc := range []struct {
		name  string
		value string
	}{
		{"malformed json", `{"a":`},
		{"non-string value", `{"a":1}`},
		{"empty name", `{"":"v"}`},
		{"invalid name charset", `{"a/b":"v"}`},
		{"invalid value charset", `{"a":"has space"}`},
	} {
		if _, err := parseWholeStringMapValue(
			PatchOperation{Op: "replace", Path: "/variables", Value: tc.value},
			validateStageVariableName, validateStageVariableValue); err == nil {
			t.Fatalf("%s was accepted", tc.name)
		}
	}

	// Without member constraints only the empty name is rejected.
	if _, err := parseWholeStringMapValue(
		PatchOperation{Op: "replace", Path: "/requestModels", Value: `{"a/b":"Empty"}`}, nil, nil); err != nil {
		t.Fatalf("constrained-free member rejected a legal key: %v", err)
	}
}

// TestApplyCanarySettingsPatchWholeMemberOverride pins the whole-member
// canary override row: replace sets the map from the JSON object value
// (with the stage variable constraints), and the operations the patch
// table does not support return an error instead of a silent success.
func TestApplyCanarySettingsPatchWholeMemberOverride(t *testing.T) {
	stage := &apigateway.Stage{
		CanarySettings: &apigateway.CanarySettings{
			StageVariableOverrides: map[string]string{"old": "v"},
		},
	}

	if err := applyCanarySettingsPatch(stage, PatchOperation{
		Op:    "replace",
		Path:  "/canarySettings/stageVariableOverrides",
		Value: `{"env":"prod"}`,
	}); err != nil {
		t.Fatalf("whole-member replace failed: %v", err)
	}
	if got := stage.CanarySettings.StageVariableOverrides["env"]; got != "prod" {
		t.Fatalf("whole-member replace not applied: %v", stage.CanarySettings.StageVariableOverrides)
	}
	if len(stage.CanarySettings.StageVariableOverrides) != 1 {
		t.Fatalf("whole-member replace did not replace the map: %v", stage.CanarySettings.StageVariableOverrides)
	}

	for _, op := range []string{"add", "remove"} {
		if err := applyCanarySettingsPatch(stage, PatchOperation{Op: op, Path: "/canarySettings/stageVariableOverrides"}); err == nil {
			t.Fatalf("unsupported op %q on the whole member was accepted", op)
		}
	}
	if stage.CanarySettings.StageVariableOverrides["env"] != "prod" {
		t.Fatalf("map mutated by rejected patches: %v", stage.CanarySettings.StageVariableOverrides)
	}

	if err := applyCanarySettingsPatch(stage, PatchOperation{
		Op:    "replace",
		Path:  "/canarySettings/stageVariableOverrides",
		Value: `{"bad name":"v"}`,
	}); err == nil {
		t.Fatal("override name violating the stage variable charset was accepted")
	}
}

// TestParseWholeMethodSettingsValue pins the whole-member /methodSettings
// replace decoding: raw-slash keys, AWS member names, and the same
// validators as the per-setting patch form.
func TestParseWholeMethodSettingsValue(t *testing.T) {
	settings, err := parseWholeMethodSettingsValue(PatchOperation{
		Op:    "replace",
		Path:  "/methodSettings",
		Value: `{"*/GET":{"metricsEnabled":true,"loggingLevel":"INFO","throttlingRateLimit":12.5},"/pets/{petId}/GET":{"cachingEnabled":true,"cacheTtlInSeconds":300}}`,
	})
	if err != nil {
		t.Fatalf("valid value rejected: %v", err)
	}
	wildcard := settings["*/GET"]
	if wildcard == nil || !wildcard.MetricsEnabled || wildcard.LoggingLevel != "INFO" || wildcard.ThrottlingRateLimit != 12.5 {
		t.Fatalf("wildcard entry not decoded: %+v", wildcard)
	}
	pets := settings["/pets/{petId}/GET"]
	if pets == nil || !pets.CachingEnabled || pets.CacheTtlInSeconds != 300 {
		t.Fatalf("raw-slash key entry not decoded: %+v", pets)
	}

	for _, tc := range []struct {
		name  string
		value string
	}{
		{"malformed json", `{"*/GET":`},
		{"empty key", `{"":{}}`},
		{"bad loggingLevel", `{"*/GET":{"loggingLevel":"NOISY"}}`},
		{"out-of-range burst", `{"*/GET":{"throttlingBurstLimit":200000}}`},
		{"out-of-range ttl", `{"*/GET":{"cacheTtlInSeconds":86401}}`},
	} {
		if _, err := parseWholeMethodSettingsValue(PatchOperation{
			Op: "replace", Path: "/methodSettings", Value: tc.value,
		}); err == nil {
			t.Fatalf("%s was accepted", tc.name)
		}
	}
}

// TestApplyWholeMapPatches pins the shared whole-member semantics for the
// map members whose official patch table rows allow add, replace and
// remove: the JSON object value replaces the map, remove clears it, and
// any other operation is rejected.
func TestApplyWholeMapPatches(t *testing.T) {
	strMap := map[string]string{"old": "v"}
	if err := applyWholeStringMapPatch(&strMap, PatchOperation{
		Op: "replace", Path: "/requestModels", Value: `{"application/json":"Empty"}`,
	}, nil, nil); err != nil {
		t.Fatalf("replace failed: %v", err)
	}
	if strMap["application/json"] != "Empty" || len(strMap) != 1 {
		t.Fatalf("whole replace not applied: %v", strMap)
	}
	if err := applyWholeStringMapPatch(&strMap, PatchOperation{Op: "remove", Path: "/requestModels"}, nil, nil); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if strMap != nil {
		t.Fatalf("whole remove did not clear the map: %v", strMap)
	}
	if err := applyWholeStringMapPatch(&strMap, PatchOperation{Op: "copy", Path: "/requestModels"}, nil, nil); err == nil {
		t.Fatal("unsupported op on a whole-member path was accepted")
	}

	boolMap := map[string]bool{"method.request.header.Authorization": false}
	if err := applyWholeBoolMapPatch(&boolMap, PatchOperation{
		Op: "add", Path: "/requestParameters", Value: `{"method.request.header.XApiKey":true}`,
	}); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if !boolMap["method.request.header.XApiKey"] || len(boolMap) != 1 {
		t.Fatalf("whole add not applied: %v", boolMap)
	}
	if err := applyWholeBoolMapPatch(&boolMap, PatchOperation{
		Op: "add", Path: "/requestParameters", Value: `{"":true}`,
	}); err == nil {
		t.Fatal("empty entry name in a whole-member value was accepted")
	}
}

// TestApplyMethodSettingsSettingPaths pins the documented
// /{resourcePath}/{httpMethod}/{group}/{member} setting family: the map key
// is the as-addressed pointer token joined to the method — the official CLI
// update-stage example output returns "~1resourceName/GET" and stores the
// wildcard as the literal "*/*" — and every row of the family documents
// replace only.
func TestApplyMethodSettingsSettingPaths(t *testing.T) {
	stage := &apigateway.Stage{}

	if err := applyMethodSettingsSetting(stage, PatchOperation{
		Op:    "replace",
		Path:  "/~1pets~1{petId}/GET/logging/loglevel",
		Value: "INFO",
	}, mustParseMethodSettingsSettingPath(t, "/~1pets~1{petId}/GET/logging/loglevel")); err != nil {
		t.Fatalf("escaped-form patch failed: %v", err)
	}
	if setting, ok := stage.MethodSettings["~1pets~1{petId}/GET"]; !ok || setting.LoggingLevel != "INFO" {
		t.Fatalf("as-addressed key missing: %v", stage.MethodSettings)
	}

	if err := applyMethodSettingsSetting(stage, PatchOperation{
		Op:    "replace",
		Path:  "/*/*/metrics/enabled",
		Value: "true",
	}, mustParseMethodSettingsSettingPath(t, "/*/*/metrics/enabled")); err != nil {
		t.Fatalf("wildcard patch failed: %v", err)
	}
	if setting, ok := stage.MethodSettings["*/*"]; !ok || !setting.MetricsEnabled {
		t.Fatalf("wildcard key missing: %v", stage.MethodSettings)
	}

	if err := applyMethodSettingsSetting(stage, PatchOperation{
		Op:    "replace",
		Path:  "//GET/caching/ttlInSeconds",
		Value: "300",
	}, mustParseMethodSettingsSettingPath(t, "//GET/caching/ttlInSeconds")); err != nil {
		t.Fatalf("root patch failed: %v", err)
	}
	if setting, ok := stage.MethodSettings["//GET"]; !ok || setting.CacheTtlInSeconds != 300 {
		t.Fatalf("root key missing: %v", stage.MethodSettings)
	}

	if err := applyMethodSettingsSetting(stage, PatchOperation{
		Op:    "replace",
		Path:  "/~1pets/GET/caching/unauthorizedCacheControlHeaderStrategy",
		Value: "SUCCEED_WITH_RESPONSE_HEADER",
	}, mustParseMethodSettingsSettingPath(t, "/~1pets/GET/caching/unauthorizedCacheControlHeaderStrategy")); err != nil {
		t.Fatalf("strategy patch failed: %v", err)
	}
	if setting, ok := stage.MethodSettings["~1pets/GET"]; !ok ||
		setting.UnauthorizedCacheControlHeaderStrategy != "SUCCEED_WITH_RESPONSE_HEADER" {
		t.Fatalf("strategy not stored: %v", stage.MethodSettings)
	}
	if err := applyMethodSettingsSetting(stage, PatchOperation{
		Op:    "replace",
		Path:  "/~1pets/GET/caching/unauthorizedCacheControlHeaderStrategy",
		Value: "FAIL_WITH_40X",
	}, mustParseMethodSettingsSettingPath(t, "/~1pets/GET/caching/unauthorizedCacheControlHeaderStrategy")); err == nil {
		t.Fatal("invalid strategy value accepted")
	}

	before := len(stage.MethodSettings)
	for _, op := range []string{"add", "remove"} {
		po := PatchOperation{Op: op, Path: "/~1pets/GET/logging/dataTrace"}
		setting, ok := parseMethodSettingsSettingPath(po.Path)
		if !ok {
			t.Fatalf("path %q did not parse", po.Path)
		}
		if err := applyMethodSettingsSetting(stage, po, setting); err == nil {
			t.Fatalf("op %q accepted", op)
		}
	}
	if len(stage.MethodSettings) != before {
		t.Fatalf("rejected ops mutated the settings: %v", stage.MethodSettings)
	}

	// The undocumented per-key form /methodSettings/{key}/{member} must not
	// parse as a row of the family.
	for _, path := range []string{
		"/methodSettings/~1pets/GET/loggingLevel",
		"/methodSettings/~1pets~1GET",
		"/pets/GET/logging",
		"/pets/GET/logging/loglevel/extra",
	} {
		if _, ok := parseMethodSettingsSettingPath(path); ok {
			t.Fatalf("path %q parsed as a documented setting row", path)
		}
	}
}

// mustParseMethodSettingsSettingPath asserts the documented family accepts
// the given path and returns the parse for apply-side tests.
func mustParseMethodSettingsSettingPath(t *testing.T, path string) methodSettingsSettingPath {
	t.Helper()
	p, ok := parseMethodSettingsSettingPath(path)
	if !ok {
		t.Fatalf("path %q did not parse as a documented setting row", path)
	}
	return p
}

// TestMethodMapKeyMatchesStoreDerivation pins the cross-plane key invariant:
// a method-settings patch addressed with the escaped form of a resource path
// stores (methodMapKey) exactly the key the execution plane derives for the
// same route (apigateway.MethodSettingsKey). The root route spells
// differently per the model documentation ("//GET", the empty resource
// token) — the execution plane consults that form through its raw-slash
// candidate and the fully-escaped "~1/GET" address through the store
// derivation, so both are pinned here.
func TestMethodMapKeyMatchesStoreDerivation(t *testing.T) {
	for _, c := range []struct{ resourcePath, method string }{
		{"/pets", "GET"},
		{"/pets/{petId}", "GET"},
		{"/a/b~c/d", "POST"},
	} {
		address := "/" + apigateway.EscapeResourcePath(c.resourcePath) + "/" + c.method + "/throttling/rateLimit"
		p := mustParseMethodSettingsSettingPath(t, address)
		got := methodMapKey(p.resourceTokens, p.httpMethod)
		want := apigateway.MethodSettingsKey(c.resourcePath, c.method)
		if got != want {
			t.Fatalf("key mismatch for %s %s: control plane %q vs execution plane %q",
				c.resourcePath, c.method, got, want)
		}
	}

	if got := methodMapKey([]string{""}, "GET"); got != "//GET" {
		t.Fatalf("root model-documented key = %q, want //GET", got)
	}
	if got := apigateway.MethodSettingsKey("/", "GET"); got != "~1/GET" {
		t.Fatalf("root fully-escaped key = %q, want ~1/GET", got)
	}
}
