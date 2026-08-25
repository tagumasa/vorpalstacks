package testutil

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssfn "github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

// runSFNIntrinsicTests pins the ASL intrinsic function family end to end:
// payload templates with intrinsics, the Fail ErrorPath/CausePath intrinsic
// form, the States.Runtime failure classification, and the Parallel/Map
// Parameters and ResultSelector payload templates.
func (r *TestRunner) runSFNIntrinsicTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	ts := fmt.Sprintf("%d", time.Now().UnixNano())

	results = append(results, r.RunTest("stepfunctions", "Intrinsics_PassParametersApplied", func() error {
		smARN, err := tc.createSingleStateSM("Intrin-"+ts, map[string]interface{}{
			"Type": "Pass",
			"Parameters": map[string]interface{}{
				"greeting.$":   "States.Format('Hello, my name is {}.', $.name)",
				"buildId.$":    "States.Array($.Id, $.Id)",
				"chunks.$":     "States.ArrayPartition($.nine, 4)",
				"contains.$":   "States.ArrayContains($.nine, 5)",
				"range.$":      "States.ArrayRange(1, 9, 2)",
				"item.$":       "States.ArrayGetItem($.nine, 5)",
				"length.$":     "States.ArrayLength($.nine)",
				"unique.$":     "States.ArrayUnique($.dups)",
				"merged.$":     "States.JsonMerge($.json1, $.json2, false)",
				"parsed.$":     "States.StringToJson($.escaped)",
				"serialised.$": "States.JsonToString($.json1)",
				"sum.$":        "States.MathAdd(111, -1)",
				"split.$":      "States.StringSplit('This.is+a,test=string', '.+,=')",
				"seededA.$":    "States.MathRandom(1, 999, 42)",
				"seededB.$":    "States.MathRandom(1, 999, 42)",
				"literal.$":    "States.Format('escaped \\{\\} and {}', $.name)",
			},
			"ResultPath": "$",
			"End":        true,
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		_, output, rerr := tc.runWithInput(smARN, "intrin-"+ts,
			`{"name":"Arnav","Id":123456,"nine":[1,2,3,4,5,6,7,8,9],"dups":[1,2,3,3,3,4],"json1":{"a":1,"keep":true},"json2":{"b":2},"escaped":"{\"foo\": \"bar\"}"}`)
		if rerr != nil {
			return rerr
		}
		var got map[string]interface{}
		if err := json.Unmarshal([]byte(output), &got); err != nil {
			return fmt.Errorf("output not JSON: %v (%s)", err, output)
		}
		want := map[string]interface{}{
			"greeting":   "Hello, my name is Arnav.",
			"contains":   true,
			"item":       float64(6),
			"length":     float64(9),
			"sum":        float64(110),
			"seededA":    got["seededB"], // same seed draws the same number
			"literal":    "escaped {} and Arnav",
			"serialised": `{"a":1,"keep":true}`,
		}
		for key, expected := range want {
			if got[key] != expected {
				return fmt.Errorf("%s = %v, want %v", key, got[key], expected)
			}
		}
		if buildID, ok := got["buildId"].([]interface{}); !ok || len(buildID) != 2 {
			return fmt.Errorf("buildId = %v", got["buildId"])
		}
		if chunks, ok := got["chunks"].([]interface{}); !ok || len(chunks) != 3 {
			return fmt.Errorf("chunks = %v", got["chunks"])
		}
		if rng, ok := got["range"].([]interface{}); !ok || len(rng) != 5 || rng[4] != float64(9) {
			return fmt.Errorf("range = %v", got["range"])
		}
		if uniq, ok := got["unique"].([]interface{}); !ok || len(uniq) != 4 {
			return fmt.Errorf("unique = %v", got["unique"])
		}
		merged, ok := got["merged"].(map[string]interface{})
		if !ok || merged["a"] != float64(1) || merged["b"] != float64(2) || merged["keep"] != true {
			return fmt.Errorf("merged = %v", got["merged"])
		}
		parsed, ok := got["parsed"].(map[string]interface{})
		if !ok || parsed["foo"] != "bar" {
			return fmt.Errorf("parsed = %v", got["parsed"])
		}
		split, ok := got["split"].([]interface{})
		if !ok || len(split) != 5 || split[0] != "This" || split[4] != "string" {
			return fmt.Errorf("split = %v", got["split"])
		}
		if n, ok := got["seededA"].(float64); !ok || n < 1 || n >= 999 {
			return fmt.Errorf("seededA = %v, want 1 <= n < 999", got["seededA"])
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Intrinsics_EncodingHashUUID", func() error {
		digest := sha256.Sum256([]byte("input data"))
		wantHash := hex.EncodeToString(digest[:])
		smARN, err := tc.createSingleStateSM("IntrinEnc-"+ts, map[string]interface{}{
			"Type": "Pass",
			"Parameters": map[string]interface{}{
				"encoded.$":    "States.Base64Encode($.plain)",
				"decoded.$":    "States.Base64Decode($.b64)",
				"hashed.$":     "States.Hash($.plain, 'SHA-256')",
				"hashedPath.$": "States.Hash($.plain, $.algo)",
				"uuid.$":       "States.UUID()",
			},
			"ResultPath": "$",
			"End":        true,
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		_, output, rerr := tc.runWithInput(smARN, "intrinenc-"+ts,
			`{"plain":"input data","b64":"RGF0YSB0byBlbmNvZGU=","algo":"SHA-256"}`)
		if rerr != nil {
			return rerr
		}
		var got map[string]interface{}
		if err := json.Unmarshal([]byte(output), &got); err != nil {
			return fmt.Errorf("output not JSON: %v (%s)", err, output)
		}
		if got["decoded"] != "Data to encode" {
			return fmt.Errorf("decoded = %v", got["decoded"])
		}
		if got["hashed"] != wantHash || got["hashedPath"] != wantHash {
			return fmt.Errorf("hash = %v / %v, want %s", got["hashed"], got["hashedPath"], wantHash)
		}
		if enc, ok := got["encoded"].(string); !ok || enc == "" {
			return fmt.Errorf("encoded = %v", got["encoded"])
		}
		uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
		if id, ok := got["uuid"].(string); !ok || !uuidPattern.MatchString(id) {
			return fmt.Errorf("uuid = %v, not a v4 UUID", got["uuid"])
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Intrinsics_FailErrorCausePathIntrinsic", func() error {
		smARN, err := tc.createSingleStateSM("IntrinFail-"+ts, map[string]interface{}{
			"Type":      "Fail",
			"ErrorPath": "States.Format('Code-{}', $.n)",
			"CausePath": "States.Format('failed at {}', $.n)",
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		execArn, serr := tc.startExecution(smARN, "intrinfail-"+ts, `{"n":7}`)
		if serr != nil {
			return serr
		}
		desc, ferr := tc.awaitTerminal(execArn, 200*time.Millisecond, 60)
		if ferr != nil {
			return ferr
		}
		if desc.Status != types.ExecutionStatusFailed {
			return fmt.Errorf("status = %s, want FAILED", desc.Status)
		}
		if aws.ToString(desc.Error) != "Code-7" {
			return fmt.Errorf("error = %q, want Code-7", aws.ToString(desc.Error))
		}
		if aws.ToString(desc.Cause) != "failed at 7" {
			return fmt.Errorf("cause = %q, want 'failed at 7'", aws.ToString(desc.Cause))
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Intrinsics_InvalidIntrinsicFailsAsRuntime", func() error {
		smARN, err := tc.createSingleStateSM("IntrinBad-"+ts, map[string]interface{}{
			"Type": "Pass",
			"Parameters": map[string]interface{}{
				"x.$": "States.ArrayRange(1, 9, 0)",
			},
			"End": true,
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		execArn, serr := tc.startExecution(smARN, "intrinbad-"+ts, `{}`)
		if serr != nil {
			return serr
		}
		desc, ferr := tc.awaitTerminal(execArn, 200*time.Millisecond, 60)
		if ferr != nil {
			return ferr
		}
		if desc.Status != types.ExecutionStatusFailed {
			return fmt.Errorf("status = %s, want FAILED", desc.Status)
		}
		if aws.ToString(desc.Error) != "States.Runtime" {
			return fmt.Errorf("error = %q, want States.Runtime", aws.ToString(desc.Error))
		}
		if aws.ToString(desc.Cause) == "" {
			return fmt.Errorf("cause is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Intrinsics_ParallelParametersResultSelector", func() error {
		smARN, err := tc.createSingleStateSM("IntrinPar-"+ts, map[string]interface{}{
			"Type": "Parallel",
			"Parameters": map[string]interface{}{
				"word.$": "States.Format('{}-par', $.w)",
			},
			"Branches": []map[string]interface{}{
				{
					"StartAt": "B1",
					"States": map[string]interface{}{
						"B1": map[string]interface{}{"Type": "Pass", "ResultPath": "$.b1", "End": true},
					},
				},
				{
					"StartAt": "B2",
					"States": map[string]interface{}{
						"B2": map[string]interface{}{"Type": "Pass", "ResultPath": "$.b2", "End": true},
					},
				},
			},
			"ResultSelector": map[string]interface{}{
				"firstWord.$": "$.[0].word",
				"count.$":     "States.ArrayLength($)",
			},
			"ResultPath": "$.combined",
			"End":        true,
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		_, output, rerr := tc.runWithInput(smARN, "intrinpar-"+ts, `{"w":"go"}`)
		if rerr != nil {
			return rerr
		}
		var got map[string]interface{}
		if err := json.Unmarshal([]byte(output), &got); err != nil {
			return fmt.Errorf("output not JSON: %v (%s)", err, output)
		}
		combined, ok := got["combined"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("combined = %v", got["combined"])
		}
		if combined["firstWord"] != "go-par" {
			return fmt.Errorf("firstWord = %v", combined["firstWord"])
		}
		if combined["count"] != float64(2) {
			return fmt.Errorf("count = %v", combined["count"])
		}
		if got["word"] != "go-par" {
			return fmt.Errorf("post-Parameters input missing: %v", got["word"])
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Intrinsics_MapLegacyParametersResultSelector", func() error {
		smARN, err := tc.createSingleStateSM("IntrinMap-"+ts, map[string]interface{}{
			"Type":      "Map",
			"ItemsPath": "$.items",
			"Parameters": map[string]interface{}{
				"n.$":    "$.value.n",
				"size.$": "States.ArrayLength($.value.tags)",
			},
			"Iterator": map[string]interface{}{
				"StartAt": "W",
				"States": map[string]interface{}{
					"W": map[string]interface{}{"Type": "Pass", "ResultPath": "$.out", "End": true},
				},
			},
			"ResultSelector": map[string]interface{}{
				"firstSize.$": "$.[0].size",
				"count.$":     "States.ArrayLength($)",
			},
			"ResultPath": "$.combined",
			"End":        true,
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		_, output, rerr := tc.runWithInput(smARN, "intrinmap-"+ts,
			`{"items":[{"value":{"n":1,"tags":["a","b"]}},{"value":{"n":2,"tags":["c"]}}]}`)
		if rerr != nil {
			return rerr
		}
		var got map[string]interface{}
		if err := json.Unmarshal([]byte(output), &got); err != nil {
			return fmt.Errorf("output not JSON: %v (%s)", err, output)
		}
		combined, ok := got["combined"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("combined = %v", got["combined"])
		}
		if combined["firstSize"] != float64(2) {
			return fmt.Errorf("firstSize = %v", combined["firstSize"])
		}
		if combined["count"] != float64(2) {
			return fmt.Errorf("count = %v", combined["count"])
		}
		return nil
	}))

	return results
}
