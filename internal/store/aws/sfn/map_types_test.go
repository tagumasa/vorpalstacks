package sfn

import (
	"encoding/json"
	"testing"
)

// TestMapStateDistributedFieldsParse pins the JSON deserialisation of the
// Distributed Map surface: ItemReader, ResultWriter, the ItemProcessor
// ProcessorConfig, the tolerated-failure thresholds and the Label field.
func TestMapStateDistributedFieldsParse(t *testing.T) {
	definition := `{
	  "StartAt": "M",
	  "States": {
	    "M": {
	      "Type": "Map",
	      "Label": "csvMap",
	      "ItemReader": {
	        "Resource": "arn:aws:states:::s3:getObject",
	        "Parameters": {"Bucket": "b", "Key": "d.csv", "VersionId": "v1"},
	        "ReaderConfig": {
	          "InputType": "CSV",
	          "CSVHeaderLocation": "GIVEN",
	          "CSVHeaders": ["a", "b"],
	          "CSVDelimiter": "PIPE",
	          "MaxItems": 100
	        }
	      },
	      "ItemProcessor": {
	        "ProcessorConfig": {"Mode": "DISTRIBUTED", "ExecutionType": "EXPRESS"},
	        "StartAt": "W",
	        "States": {"W": {"Type": "Wait", "Seconds": 1, "End": true}}
	      },
	      "ToleratedFailureCount": 2,
	      "ToleratedFailurePercentage": 15.5,
	      "ResultWriter": {
	        "Resource": "arn:aws:states:::s3:putObject",
	        "Parameters": {"Bucket": "out", "Prefix": "jobs"},
	        "WriterConfig": {"Transformation": "FLATTEN", "OutputType": "JSONL"}
	      },
	      "End": true
	    }
	  }
	}`

	var def StateMachineDefinition
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	raw := def.States["M"]
	var ms MapState
	if err := jsonRemarshal(raw, &ms); err != nil {
		t.Fatalf("remarshal Map state failed: %v", err)
	}

	if ms.Label != "csvMap" {
		t.Errorf("Label = %q, want csvMap", ms.Label)
	}
	ir := ms.ItemReader
	if ir == nil {
		t.Fatalf("ItemReader not parsed")
	}
	if ir.Resource != "arn:aws:states:::s3:getObject" {
		t.Errorf("ItemReader.Resource = %q", ir.Resource)
	}
	var params map[string]string
	if err := json.Unmarshal(ir.Parameters, &params); err != nil {
		t.Fatalf("ItemReader.Parameters not valid JSON: %v", err)
	}
	if params["Bucket"] != "b" || params["Key"] != "d.csv" || params["VersionId"] != "v1" {
		t.Errorf("ItemReader.Parameters = %v", params)
	}
	rc := ir.ReaderConfig
	if rc == nil {
		t.Fatalf("ReaderConfig not parsed")
	}
	if rc.InputType != "CSV" || rc.CSVHeaderLocation != "GIVEN" || rc.CSVDelimiter != "PIPE" {
		t.Errorf("ReaderConfig scalars = %+v", rc)
	}
	if len(rc.CSVHeaders) != 2 || rc.CSVHeaders[0] != "a" {
		t.Errorf("CSVHeaders = %v", rc.CSVHeaders)
	}
	if rc.MaxItems == nil || *rc.MaxItems != 100 {
		t.Errorf("MaxItems = %v", rc.MaxItems)
	}

	ip := ms.ItemProcessor
	if ip == nil || ip.ProcessorConfig == nil {
		t.Fatalf("ItemProcessor.ProcessorConfig not parsed")
	}
	if ip.ProcessorConfig.Mode != "DISTRIBUTED" || ip.ProcessorConfig.ExecutionType != "EXPRESS" {
		t.Errorf("ProcessorConfig = %+v", ip.ProcessorConfig)
	}

	if ms.ToleratedFailureCount == nil || *ms.ToleratedFailureCount != 2 {
		t.Errorf("ToleratedFailureCount = %v", ms.ToleratedFailureCount)
	}
	if ms.ToleratedFailurePercentage == nil || *ms.ToleratedFailurePercentage != 15.5 {
		t.Errorf("ToleratedFailurePercentage = %v", ms.ToleratedFailurePercentage)
	}

	rw := ms.ResultWriter
	if rw == nil {
		t.Fatalf("ResultWriter not parsed")
	}
	if rw.Resource != "arn:aws:states:::s3:putObject" {
		t.Errorf("ResultWriter.Resource = %q", rw.Resource)
	}
	var rwParams map[string]string
	if err := json.Unmarshal(rw.Parameters, &rwParams); err != nil {
		t.Fatalf("ResultWriter.Parameters not valid JSON: %v", err)
	}
	if rwParams["Bucket"] != "out" || rwParams["Prefix"] != "jobs" {
		t.Errorf("ResultWriter.Parameters = %v", rwParams)
	}
	if rw.WriterConfig == nil || rw.WriterConfig.Transformation != "FLATTEN" || rw.WriterConfig.OutputType != "JSONL" {
		t.Errorf("WriterConfig = %+v", rw.WriterConfig)
	}
}

// TestMapStatePathVariantsParse pins the *Path variants of the tolerated
// failure thresholds, which must survive a round trip unchanged.
func TestMapStatePathVariantsParse(t *testing.T) {
	definition := `{
	  "StartAt": "M",
	  "States": {
	    "M": {
	      "Type": "Map",
	      "Iterator": {"StartAt": "S", "States": {"S": {"Type": "Succeed"}}},
	      "ToleratedFailureCountPath": "$.count",
	      "ToleratedFailurePercentagePath": "$.pct",
	      "End": true
	    }
	  }
	}`

	var def StateMachineDefinition
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	var ms MapState
	if err := jsonRemarshal(def.States["M"], &ms); err != nil {
		t.Fatalf("remarshal Map state failed: %v", err)
	}
	if ms.ToleratedFailureCountPath != "$.count" {
		t.Errorf("ToleratedFailureCountPath = %q", ms.ToleratedFailureCountPath)
	}
	if ms.ToleratedFailurePercentagePath != "$.pct" {
		t.Errorf("ToleratedFailurePercentagePath = %q", ms.ToleratedFailurePercentagePath)
	}
	if ms.ItemReader != nil || ms.ResultWriter != nil || ms.Label != "" {
		t.Errorf("absent distributed fields should stay nil: %+v", ms)
	}
}

// TestMapStateItemBatcherParse pins the JSON deserialisation of the
// ItemBatcher: both sizing pairs, their reference-path variants and the
// fixed BatchInput object.
func TestMapStateItemBatcherParse(t *testing.T) {
	definition := `{
	  "StartAt": "M",
	  "States": {
	    "M": {
	      "Type": "Map",
	      "ItemProcessor": {
	        "ProcessorConfig": {"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
	        "StartAt": "P",
	        "States": {"P": {"Type": "Pass", "End": true}}
	      },
	      "ItemBatcher": {
	        "MaxItemsPerBatch": 500,
	        "MaxInputBytesPerBatchPath": "$.batchBytes",
	        "BatchInput": {"key": "value"}
	      },
	      "End": true
	    }
	  }
	}`

	var def StateMachineDefinition
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	var ms MapState
	if err := jsonRemarshal(def.States["M"], &ms); err != nil {
		t.Fatalf("remarshal Map state failed: %v", err)
	}
	ib := ms.ItemBatcher
	if ib == nil {
		t.Fatalf("ItemBatcher not parsed")
	}
	if ib.MaxItemsPerBatch != float64(500) {
		t.Errorf("MaxItemsPerBatch = %v (%T)", ib.MaxItemsPerBatch, ib.MaxItemsPerBatch)
	}
	if ib.MaxItemsPerBatchPath != "" {
		t.Errorf("MaxItemsPerBatchPath = %q, want empty", ib.MaxItemsPerBatchPath)
	}
	if ib.MaxInputBytesPerBatch != nil {
		t.Errorf("MaxInputBytesPerBatch = %v, want nil", ib.MaxInputBytesPerBatch)
	}
	if ib.MaxInputBytesPerBatchPath != "$.batchBytes" {
		t.Errorf("MaxInputBytesPerBatchPath = %q", ib.MaxInputBytesPerBatchPath)
	}
	bi, ok := ib.BatchInput.(map[string]interface{})
	if !ok || bi["key"] != "value" {
		t.Errorf("BatchInput = %#v", ib.BatchInput)
	}
	if ib.BatchInputPath != "" {
		t.Errorf("BatchInputPath = %q, want empty", ib.BatchInputPath)
	}
}

// TestMapStateItemBatcherPathVariantsParse pins the reference-path forms
// of the ItemBatcher sizing and batch input fields.
func TestMapStateItemBatcherPathVariantsParse(t *testing.T) {
	definition := `{
	  "StartAt": "M",
	  "States": {
	    "M": {
	      "Type": "Map",
	      "ItemProcessor": {"StartAt": "S", "States": {"S": {"Type": "Succeed"}}},
	      "ItemBatcher": {
	        "MaxItemsPerBatchPath": "$.max",
	        "MaxInputBytesPerBatch": 1024,
	        "BatchInputPath": "$.fixed"
	      },
	      "End": true
	    }
	  }
	}`

	var def StateMachineDefinition
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	var ms MapState
	if err := jsonRemarshal(def.States["M"], &ms); err != nil {
		t.Fatalf("remarshal Map state failed: %v", err)
	}
	ib := ms.ItemBatcher
	if ib == nil {
		t.Fatalf("ItemBatcher not parsed")
	}
	if ib.MaxItemsPerBatchPath != "$.max" {
		t.Errorf("MaxItemsPerBatchPath = %q", ib.MaxItemsPerBatchPath)
	}
	if ib.MaxInputBytesPerBatch == nil || *ib.MaxInputBytesPerBatch != 1024 {
		t.Errorf("MaxInputBytesPerBatch = %v", ib.MaxInputBytesPerBatch)
	}
	if ib.BatchInputPath != "$.fixed" {
		t.Errorf("BatchInputPath = %q", ib.BatchInputPath)
	}
}

// jsonRemarshal converts a decoded generic value back through JSON so a
// concrete struct can be deserialised from it.
func jsonRemarshal(v interface{}, target interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}
