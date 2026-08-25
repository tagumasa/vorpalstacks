package testutil

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awssfn "github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"

	"vorpalstacks-sdk-tests/config"
)

// runSFNItemReaderTests pins the Distributed Map ItemReader and ResultWriter
// data plane over real S3 objects, the ValidateStateMachineDefinition
// documented diagnostic codes, and the TestState stateConfiguration
// contracts.
func (r *TestRunner) runSFNItemReaderTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return []TestResult{{Service: "stepfunctions", TestName: "ItemReaderSetup", Status: "FAIL", Error: fmt.Sprintf("load config: %v", err)}}
	}
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	bucket := "sfn-itemreader-" + ts
	if _, err := s3Client.CreateBucket(tc.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
		return []TestResult{{Service: "stepfunctions", TestName: "ItemReaderSetup", Status: "FAIL", Error: fmt.Sprintf("create bucket: %v", err)}}
	}
	defer s3CleanupBucket(tc.ctx, s3Client, bucket)

	put := func(key string, body string) error {
		_, err := s3Client.PutObject(tc.ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket), Key: aws.String(key), Body: strings.NewReader(body),
		})
		return err
	}

	runAndCollect := func(smArn string) (string, []interface{}, error) {
		_, output, err := tc.runWithInput(smArn, "", `{}`)
		if err != nil {
			return "", nil, err
		}
		var items []interface{}
		_ = json.Unmarshal([]byte(output), &items)
		return output, items, nil
	}

	processor := func() map[string]interface{} {
		return map[string]interface{}{
			"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
			"StartAt":         "W",
			"States": map[string]interface{}{
				"W": map[string]interface{}{"Type": "Pass", "ResultPath": "$", "End": true},
			},
		}
	}

	results = append(results, r.RunTest("stepfunctions", "ItemReader_CSV_FirstRow", func() error {
		if err := put("items.csv", "n\n1\n2\n3\n"); err != nil {
			return err
		}
		smARN, cerr := tc.createSingleStateSM(fmt.Sprintf("IrCsv-%d", time.Now().UnixNano()), map[string]interface{}{
			"Type": "Map",
			"ItemReader": map[string]interface{}{
				"Resource":     "arn:aws:states:::s3:getObject",
				"Parameters":   map[string]interface{}{"Bucket": bucket, "Key": "items.csv"},
				"ReaderConfig": map[string]interface{}{"InputType": "CSV", "CSVHeaderLocation": "FIRST_ROW"},
			},
			"ItemProcessor": processor(),
			"End":           true,
		})
		if cerr != nil {
			return cerr
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		_, items, err := runAndCollect(smARN)
		if err != nil {
			return err
		}
		if len(items) != 3 {
			return fmt.Errorf("csv items = %v, want 3", items)
		}
		if items[0].(map[string]interface{})["n"] != "1" {
			return fmt.Errorf("first csv item = %v (fields must stay strings)", items[0])
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ItemReader_CSV_GivenHeaders_Pipe", func() error {
		if err := put("pipe.csv", "1|307\n1|481\n"); err != nil {
			return err
		}
		smARN, cerr := tc.createSingleStateSM(fmt.Sprintf("IrPipe-%d", time.Now().UnixNano()), map[string]interface{}{
			"Type": "Map",
			"ItemReader": map[string]interface{}{
				"Resource":   "arn:aws:states:::s3:getObject",
				"Parameters": map[string]interface{}{"Bucket": bucket, "Key": "pipe.csv"},
				"ReaderConfig": map[string]interface{}{
					"InputType": "CSV", "CSVHeaderLocation": "GIVEN",
					"CSVHeaders": []string{"userId", "movieId"}, "CSVDelimiter": "PIPE",
				},
			},
			"ItemProcessor": processor(),
			"End":           true,
		})
		if cerr != nil {
			return cerr
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		_, items, err := runAndCollect(smARN)
		if err != nil {
			return err
		}
		if len(items) != 2 {
			return fmt.Errorf("pipe items = %v, want 2", items)
		}
		if items[0].(map[string]interface{})["movieId"] != "307" {
			return fmt.Errorf("pipe item = %v", items[0])
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ItemReader_JSON_ItemsPointer_MaxItems", func() error {
		if err := put("nested.json", `{"data":{"items":[{"id":1},{"id":2},{"id":3}]}}`); err != nil {
			return err
		}
		smARN, cerr := tc.createSingleStateSM(fmt.Sprintf("IrJson-%d", time.Now().UnixNano()), map[string]interface{}{
			"Type": "Map",
			"ItemReader": map[string]interface{}{
				"Resource":   "arn:aws:states:::s3:getObject",
				"Parameters": map[string]interface{}{"Bucket": bucket, "Key": "nested.json"},
				"ReaderConfig": map[string]interface{}{
					"InputType": "JSON", "ItemsPointer": "/data/items", "MaxItems": 2,
				},
			},
			"ItemProcessor": processor(),
			"End":           true,
		})
		if cerr != nil {
			return cerr
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		_, items, err := runAndCollect(smARN)
		if err != nil {
			return err
		}
		if len(items) != 2 {
			return fmt.Errorf("items = %v, want the MaxItems cap of 2", items)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ItemReader_JSONL_Gzip", func() error {
		if _, err := s3Client.PutObject(tc.ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket), Key: aws.String("lines.jsonl.gz"),
			Body: strings.NewReader(gzipData("{\"n\":1}\n{\"n\":2}\n")),
		}); err != nil {
			return err
		}
		smARN, cerr := tc.createSingleStateSM(fmt.Sprintf("IrJsonl-%d", time.Now().UnixNano()), map[string]interface{}{
			"Type": "Map",
			"ItemReader": map[string]interface{}{
				"Resource":     "arn:aws:states:::s3:getObject",
				"Parameters":   map[string]interface{}{"Bucket": bucket, "Key": "lines.jsonl.gz"},
				"ReaderConfig": map[string]interface{}{"InputType": "JSONL"},
			},
			"ItemProcessor": processor(),
			"End":           true,
		})
		if cerr != nil {
			return cerr
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		_, items, err := runAndCollect(smARN)
		if err != nil {
			return err
		}
		if len(items) != 2 {
			return fmt.Errorf("jsonl items = %v, want 2", items)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ItemReader_ListObjectsV2_Prefix", func() error {
		if err := put("data/f1.json", `{"k":"f1"}`); err != nil {
			return err
		}
		if err := put("data/f2.json", `{"k":"f2"}`); err != nil {
			return err
		}
		smARN, cerr := tc.createSingleStateSM(fmt.Sprintf("IrList-%d", time.Now().UnixNano()), map[string]interface{}{
			"Type": "Map",
			"ItemReader": map[string]interface{}{
				"Resource":   "arn:aws:states:::s3:listObjectsV2",
				"Parameters": map[string]interface{}{"Bucket": bucket, "Prefix": "data/"},
			},
			"ItemProcessor": processor(),
			"End":           true,
		})
		if cerr != nil {
			return cerr
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		_, items, err := runAndCollect(smARN)
		if err != nil {
			return err
		}
		if len(items) != 2 {
			return fmt.Errorf("listed items = %v, want 2", items)
		}
		entry := items[0].(map[string]interface{})
		if _, ok := entry["Key"]; !ok {
			return fmt.Errorf("list entry missing Key member: %v", entry)
		}
		if _, ok := entry["Etag"]; !ok {
			return fmt.Errorf("list entry missing Etag member: %v", entry)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ItemReader_ToleratedFailure_ExceedThreshold", func() error {
		// Every iteration fails (the processor divides by a missing
		// field) with a threshold of zero: the documented default fails
		// the Map Run with States.ExceedToleratedFailureThreshold.
		if err := put("failing.csv", "v\n1\n2\n"); err != nil {
			return err
		}
		def := map[string]interface{}{
			"StartAt": "M",
			"States": map[string]interface{}{
				"M": map[string]interface{}{
					"Type": "Map",
					"ItemReader": map[string]interface{}{
						"Resource":     "arn:aws:states:::s3:getObject",
						"Parameters":   map[string]interface{}{"Bucket": bucket, "Key": "failing.csv"},
						"ReaderConfig": map[string]interface{}{"InputType": "CSV", "CSVHeaderLocation": "FIRST_ROW"},
					},
					"ItemProcessor": map[string]interface{}{
						"ProcessorConfig": map[string]interface{}{"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
						"StartAt":         "B",
						"States": map[string]interface{}{
							"B": map[string]interface{}{"Type": "Fail", "Error": "ItemFailed", "Cause": "deliberate"},
						},
					},
					"End": true,
				},
			},
		}
		defJSON, _ := json.Marshal(def)
		resp, cerr := tc.client.CreateStateMachine(tc.ctx, &awssfn.CreateStateMachineInput{
			Name:       aws.String(fmt.Sprintf("IrFail-%d", time.Now().UnixNano())),
			Definition: aws.String(string(defJSON)),
			RoleArn:    aws.String(tc.roleARN),
		})
		if cerr != nil {
			return cerr
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: resp.StateMachineArn})

		execArn, eerr := tc.startExecution(aws.ToString(resp.StateMachineArn), "", `{}`)
		if eerr != nil {
			return eerr
		}
		desc, werr := tc.awaitTerminal(execArn, 200*time.Millisecond, 60)
		if werr != nil {
			return werr
		}
		if desc.Status == types.ExecutionStatusFailed {
			if aws.ToString(desc.Error) != "States.ExceedToleratedFailureThreshold" {
				return fmt.Errorf("error = %s, want States.ExceedToleratedFailureThreshold", aws.ToString(desc.Error))
			}
			return nil
		}
		if desc.Status == types.ExecutionStatusSucceeded {
			return fmt.Errorf("execution unexpectedly succeeded")
		}
		return fmt.Errorf("execution ended %s", desc.Status)
	}))

	results = append(results, r.RunTest("stepfunctions", "ResultWriter_ExportsManifestAndResultFiles", func() error {
		if err := put("rw.csv", "n\n1\n2\n"); err != nil {
			return err
		}
		smARN, cerr := tc.createSingleStateSM(fmt.Sprintf("IrRw-%d", time.Now().UnixNano()), map[string]interface{}{
			"Type": "Map",
			"ItemReader": map[string]interface{}{
				"Resource":     "arn:aws:states:::s3:getObject",
				"Parameters":   map[string]interface{}{"Bucket": bucket, "Key": "rw.csv"},
				"ReaderConfig": map[string]interface{}{"InputType": "CSV", "CSVHeaderLocation": "FIRST_ROW"},
			},
			"ItemProcessor": processor(),
			"ResultWriter": map[string]interface{}{
				"Resource":   "arn:aws:states:::s3:putObject",
				"Parameters": map[string]interface{}{"Bucket": bucket, "Prefix": "exports"},
			},
			"End": true,
		})
		if cerr != nil {
			return cerr
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})
		output, _, err := runAndCollect(smARN)
		if err != nil {
			return err
		}
		var summary map[string]interface{}
		if uerr := json.Unmarshal([]byte(output), &summary); uerr != nil {
			return fmt.Errorf("output is not the export summary: %s", output)
		}
		details, ok := summary["ResultWriterDetails"].(map[string]interface{})
		if !ok || details["Bucket"] != bucket {
			return fmt.Errorf("summary = %v", summary)
		}
		manifestKey, _ := details["Key"].(string)
		if !strings.HasPrefix(manifestKey, "exports/") || !strings.HasSuffix(manifestKey, "/manifest.json") {
			return fmt.Errorf("manifest key = %q", manifestKey)
		}
		get := func(key string) (string, error) {
			obj, gerr := s3Client.GetObject(tc.ctx, &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
			if gerr != nil {
				return "", gerr
			}
			defer obj.Body.Close()
			var sb strings.Builder
			buf := make([]byte, 4096)
			for {
				n, rerr := obj.Body.Read(buf)
				sb.Write(buf[:n])
				if rerr != nil {
					break
				}
			}
			return sb.String(), nil
		}
		manifest, gerr := get(manifestKey)
		if gerr != nil {
			return fmt.Errorf("manifest missing: %v", gerr)
		}
		var m map[string]interface{}
		if uerr := json.Unmarshal([]byte(manifest), &m); uerr != nil {
			return fmt.Errorf("manifest not JSON: %v", uerr)
		}
		if m["MapRunArn"] == nil || m["ResultLocation"] != fmt.Sprintf("s3://%s/exports", bucket) {
			return fmt.Errorf("manifest = %v", m)
		}
		succeededKey := strings.TrimSuffix(manifestKey, "manifest.json") + "SUCCEEDED_0.json"
		records, gerr := get(succeededKey)
		if gerr != nil {
			return fmt.Errorf("result file missing: %v", gerr)
		}
		var arr []map[string]interface{}
		if uerr := json.Unmarshal([]byte(records), &arr); uerr != nil {
			return fmt.Errorf("result file not a JSON array: %v", uerr)
		}
		if len(arr) != 2 || arr[0]["Status"] != "SUCCEEDED" || arr[0]["ExecutionArn"] == nil {
			return fmt.Errorf("records = %v", arr)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ValidateStateMachineDefinition_DocumentedCodes", func() error {
		resp, verr := tc.client.ValidateStateMachineDefinition(tc.ctx, &awssfn.ValidateStateMachineDefinitionInput{
			Definition: aws.String(`{"StartAt":"A","States":{"A":{"Type":"Pass","Next":"Missing"}}}`),
		})
		if verr != nil {
			return verr
		}
		if resp.Result != types.ValidateStateMachineDefinitionResultCodeFail {
			return fmt.Errorf("result = %s, want FAIL", resp.Result)
		}
		found := false
		for _, d := range resp.Diagnostics {
			if d.Severity == types.ValidateStateMachineDefinitionSeverityError &&
				aws.ToString(d.Code) == "MISSING_TRANSITION_TARGET" {
				found = true
				if d.Location == nil || *d.Location != "/States/A/Next" {
					return fmt.Errorf("location = %v, want /States/A/Next", d.Location)
				}
			}
		}
		if !found {
			return fmt.Errorf("MISSING_TRANSITION_TARGET missing: %v", resp.Diagnostics)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ValidateStateMachineDefinition_MissingEndState", func() error {
		resp, verr := tc.client.ValidateStateMachineDefinition(tc.ctx, &awssfn.ValidateStateMachineDefinitionInput{
			Definition: aws.String(`{"StartAt":"A","States":{"A":{"Type":"Pass","Next":"B"},"B":{"Type":"Pass","Next":"A"}}}`),
		})
		if verr != nil {
			return verr
		}
		if resp.Result != types.ValidateStateMachineDefinitionResultCodeFail {
			return fmt.Errorf("result = %s, want FAIL", resp.Result)
		}
		for _, d := range resp.Diagnostics {
			if aws.ToString(d.Code) == "MISSING_END_STATE" {
				return nil
			}
		}
		return fmt.Errorf("MISSING_END_STATE missing: %v", resp.Diagnostics)
	}))

	results = append(results, r.RunTest("stepfunctions", "TestState_RetrierRetryCount_Retriable", func() error {
		def := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:lambda:` + r.region + `:` + r.accountID + `:function:none",`
		def += `"Retry":[{"ErrorEquals":["States.TaskFailed"],"IntervalSeconds":2,"MaxAttempts":3,"BackoffRate":2.0}],`
		def += `"Next":"Done"},"Done":{"Type":"Succeed"}}}`
		result, terr := tc.rawTestState(map[string]interface{}{
			"definition": def,
			"stateName":  "T",
			"input":      `{}`,
			"mock": map[string]interface{}{
				"errorOutput": map[string]interface{}{"error": "States.TaskFailed", "cause": "task failed"},
			},
			"stateConfiguration": map[string]interface{}{"retrierRetryCount": 1},
			"inspectionLevel":    "DEBUG",
		})
		if terr != nil {
			return terr
		}
		if result["status"] != "RETRIABLE" {
			return fmt.Errorf("status = %v, want RETRIABLE: %v", result["status"], result)
		}
		inspection, _ := result["inspectionData"].(map[string]interface{})
		if inspection == nil {
			return fmt.Errorf("inspectionData missing")
		}
		details, _ := inspection["errorDetails"].(map[string]interface{})
		if details == nil {
			return fmt.Errorf("errorDetails missing from inspection data: %v", inspection)
		}
		if n, ok := details["retryIndex"].(float64); !ok || n != 0 {
			return fmt.Errorf("retryIndex = %v, want 0", details["retryIndex"])
		}
		if n, ok := details["retryBackoffIntervalSeconds"].(float64); !ok || n != 4 {
			return fmt.Errorf("retryBackoffIntervalSeconds = %v, want 4", details["retryBackoffIntervalSeconds"])
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "TestState_RetrierExhausted_CaughtError", func() error {
		def := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:lambda:us-east-1:000000000000:function:none",`
		def += `"Retry":[{"ErrorEquals":["States.TaskFailed"],"MaxAttempts":2}],`
		def += `"Catch":[{"ErrorEquals":["States.ALL"],"ResultPath":"$.error","Next":"Handle"}],"Next":"Done"},`
		def += `"Handle":{"Type":"Pass","End":true},"Done":{"Type":"Succeed"}}}`
		result, terr := tc.rawTestState(map[string]interface{}{
			"definition": def,
			"stateName":  "T",
			"input":      `{"data":"value"}`,
			"mock": map[string]interface{}{
				"errorOutput": map[string]interface{}{"error": "States.TaskFailed", "cause": "task failed"},
			},
			"stateConfiguration": map[string]interface{}{"retrierRetryCount": 2},
			"inspectionLevel":    "DEBUG",
		})
		if terr != nil {
			return terr
		}
		if result["status"] != "CAUGHT_ERROR" {
			return fmt.Errorf("status = %v, want CAUGHT_ERROR: %v", result["status"], result)
		}
		if result["nextState"] != "Handle" {
			return fmt.Errorf("nextState = %v, want Handle", result["nextState"])
		}
		inspection, _ := result["inspectionData"].(map[string]interface{})
		details, _ := inspection["errorDetails"].(map[string]interface{})
		if n, ok := details["catchIndex"].(float64); !ok || n != 0 {
			return fmt.Errorf("catchIndex = %v, want 0", details["catchIndex"])
		}
		var output map[string]interface{}
		if uerr := json.Unmarshal([]byte(result["output"].(string)), &output); uerr != nil {
			return fmt.Errorf("caught output not JSON: %v", uerr)
		}
		if output["error"] == nil || output["data"] == nil {
			return fmt.Errorf("ResultPath must add the error to the input: %v", output)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "TestState_MapIterationFailureCount_ExceedsThreshold", func() error {
		def := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemsPath":"$.items",`
		def += `"ItemProcessor":{"ProcessorConfig":{"Mode":"DISTRIBUTED","ExecutionType":"STANDARD"},"StartAt":"W","States":{"W":{"Type":"Pass","End":true}}},`
		def += `"ToleratedFailureCount":1,"End":true}}}`
		result, terr := tc.rawTestState(map[string]interface{}{
			"definition": def,
			"stateName":  "M",
			"input":      `{"items":[1,2,3]}`,
			"mock":       map[string]interface{}{"result": "[1,2,3]"},
			"stateConfiguration": map[string]interface{}{
				"mapIterationFailureCount": 2,
			},
			"inspectionLevel": "DEBUG",
		})
		if terr != nil {
			return terr
		}
		if result["status"] != "FAILED" {
			return fmt.Errorf("status = %v, want FAILED: %v", result["status"], result)
		}
		if result["error"] != "States.ExceedToleratedFailureThreshold" {
			return fmt.Errorf("error = %v, want States.ExceedToleratedFailureThreshold", result["error"])
		}
		inspection, _ := result["inspectionData"].(map[string]interface{})
		if inspection == nil || inspection["afterItemsPath"] == nil {
			return fmt.Errorf("afterItemsPath missing from inspection data: %v", inspection)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "TestState_MapItemReaderData_SubstitutesReaderSource", func() error {
		def := `{"StartAt":"M","States":{"M":{"Type":"Map",`
		def += `"ItemReader":{"Resource":"arn:aws:states:::s3:getObject","Parameters":{"Bucket":"` + bucket + `","Key":"absent.csv"},`
		def += `"ReaderConfig":{"InputType":"CSV","CSVHeaderLocation":"FIRST_ROW"}},`
		def += `"ItemProcessor":{"ProcessorConfig":{"Mode":"DISTRIBUTED","ExecutionType":"STANDARD"},"StartAt":"W","States":{"W":{"Type":"Pass","End":true}}},`
		def += `"End":true}}}`
		result, terr := tc.rawTestState(map[string]interface{}{
			"definition": def,
			"stateName":  "M",
			"input":      `{}`,
			"mock":       map[string]interface{}{"result": `[{"n":"1"},{"n":"2"}]`},
			"stateConfiguration": map[string]interface{}{
				"mapItemReaderData": "n\n1\n2\n",
			},
			"inspectionLevel": "DEBUG",
		})
		if terr != nil {
			return terr
		}
		if result["status"] != "SUCCEEDED" {
			return fmt.Errorf("status = %v: %v %v", result["status"], result["error"], result["cause"])
		}
		inspection, _ := result["inspectionData"].(map[string]interface{})
		if inspection == nil || inspection["afterItemsPath"] == nil {
			return fmt.Errorf("afterItemsPath missing")
		}
		var items []interface{}
		if uerr := json.Unmarshal([]byte(inspection["afterItemsPath"].(string)), &items); uerr != nil || len(items) != 2 {
			return fmt.Errorf("afterItemsPath = %v, want 2 items from the substituted source", inspection["afterItemsPath"])
		}
		return nil
	}))

	return results
}

// gzipData compresses a string payload for the compressed-dataset tests.
func gzipData(data string) string {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write([]byte(data))
	w.Close()
	return buf.String()
}
