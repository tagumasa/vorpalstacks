package sfn

import (
	sfnstore "vorpalstacks/internal/store/aws/sfn"

	"fmt"
	"math"
	"strings"
)

// validateMapState validates the Map contract: exactly one of Iterator or
// ItemProcessor, the ItemProcessor processing mode (Distributed mode is
// Standard-only), the ItemReader and ResultWriter structures, the
// tolerated-failure thresholds and the Label naming rules.
func (v *aslValidatorContext) validateMapState(name string, stateMap map[string]interface{}, stateLoc string, jsonata bool) {
	_, hasIterator := stateMap["Iterator"]
	_, hasItemProcessor := stateMap["ItemProcessor"]
	if hasIterator == hasItemProcessor {
		v.schemaf(stateLoc, "Map state '%s' must specify exactly one of Iterator or ItemProcessor", name)
	}
	if raw, ok := stateMap["Iterator"].(map[string]interface{}); ok {
		v.validateMapSubMachine(name, "Iterator", raw, stateLoc)
	}
	if raw, ok := stateMap["ItemProcessor"].(map[string]interface{}); ok {
		v.validateMapSubMachine(name, "ItemProcessor", raw, stateLoc)
	}

	if _, ok := stateMap["MaxConcurrency"]; ok {
		// "In JSONata states, you can specify a JSONata expression that
		// evaluates to an integer."
		v.validateCountField(stateMap, "MaxConcurrency", stateLoc, jsonata, "Map state MaxConcurrency must be a non-negative integer", nil)
	}

	// Label naming: at most 40 characters, unique within the definition,
	// and free of the documented forbidden characters.
	if raw, ok := stateMap["Label"]; ok {
		label, isString := raw.(string)
		if !isString {
			v.schemaf(stateLoc+"/Label", "Map state Label must be a string")
		} else {
			if len([]rune(label)) > sfnstore.MaxMapLabelLength {
				v.add("ERROR", "INVALID_LABEL_NAME",
					fmt.Sprintf("The label '%s' of state '%s' exceeds the allowed length of %d characters", label, name, sfnstore.MaxMapLabelLength),
					stateLoc+"/Label")
			}
			if strings.ContainsFunc(label, isForbiddenLabelRune) {
				v.add("ERROR", "INVALID_LABEL_NAME",
					fmt.Sprintf("The label '%s' of state '%s' contains forbidden characters", label, name),
					stateLoc+"/Label")
			}
			if prior, dup := v.labels[label]; dup {
				v.add("ERROR", "DUPLICATE_LABEL_NAME",
					fmt.Sprintf("The label name '%s' appears more than once (also used by state %s)", label, prior),
					stateLoc+"/Label")
			} else {
				v.labels[label] = name
			}
		}
	}

	v.validateItemReader(name, stateMap, stateLoc)
	v.validateItemBatcher(name, stateMap, stateLoc, jsonata)
	v.validateResultWriter(name, stateMap, stateLoc)
	v.validateToleratedFailure(name, stateMap, stateLoc, jsonata)
}

// validateMapSubMachine validates the Iterator or ItemProcessor
// sub-machine and its Distributed-mode constraints.
func (v *aslValidatorContext) validateMapSubMachine(name, field string, sub map[string]interface{}, stateLoc string) {
	subLoc := stateLoc + "/" + field
	startAt, hasStartAt := sub["StartAt"].(string)
	subStates, hasStates := sub["States"].(map[string]interface{})
	if !hasStartAt || startAt == "" {
		v.schemaf(subLoc+"/StartAt", "The %s of Map state '%s' must specify a string 'StartAt'", field, name)
	}
	if !hasStates {
		v.schemaf(subLoc+"/States", "The %s of Map state '%s' must specify a 'States' object", field, name)
		return
	}
	if hasStartAt {
		if _, exists := subStates[startAt]; !exists {
			v.schemaf(subLoc+"/StartAt", "StartAt '%s' of the %s does not reference a state in it", startAt, field)
		} else {
			v.checkTerminalReachability(startAt, subStates)
		}
	}

	if pc, ok := sub["ProcessorConfig"].(map[string]interface{}); ok {
		mode, _ := pc["Mode"].(string)
		if mode != "" && mode != "INLINE" && mode != "DISTRIBUTED" {
			v.schemaf(subLoc+"/ProcessorConfig/Mode", "ProcessorConfig Mode must be INLINE or DISTRIBUTED, got %q", mode)
		}
		if mode == "DISTRIBUTED" {
			executionType, _ := pc["ExecutionType"].(string)
			if executionType != "STANDARD" && executionType != "EXPRESS" {
				v.schemaf(subLoc+"/ProcessorConfig/ExecutionType", "DISTRIBUTED processing requires ExecutionType STANDARD or EXPRESS")
			}
			if v.smType == "EXPRESS" {
				v.schemaf(subLoc+"/ProcessorConfig/Mode", "Distributed mode is supported in Standard workflows but not in Express workflows")
			}
		}
	}
	v.validateStatesScope(sub, subStates, subLoc)
}

// validateItemReader validates the ItemReader structure: resource, reader
// configuration enums, the MaxItems exclusivity and range, CSV header
// requirements and the manifest constraints.
func (v *aslValidatorContext) validateItemReader(name string, stateMap map[string]interface{}, stateLoc string) {
	raw, ok := stateMap["ItemReader"].(map[string]interface{})
	if !ok {
		return
	}
	readerLoc := stateLoc + "/ItemReader"

	resource, _ := raw["Resource"].(string)
	if resource == "" {
		v.schemaf(readerLoc, "The ItemReader of Map state '%s' must specify a Resource", name)
	} else {
		known := strings.HasSuffix(resource, ":s3:getObject") ||
			strings.HasSuffix(resource, ":aws-sdk:s3:getObject") ||
			strings.HasSuffix(resource, ":s3:listObjectsV2")
		if !known {
			v.schemaf(readerLoc+"/Resource", "The ItemReader Resource '%s' of Map state '%s' is not a supported reader", resource, name)
		}
	}
	_, hasParams := raw["Parameters"]
	_, hasArgs := raw["Arguments"]
	if hasParams && hasArgs {
		v.schemaf(readerLoc, "The ItemReader of Map state '%s' cannot specify both Parameters and Arguments", name)
	}

	rc, hasRC := raw["ReaderConfig"].(map[string]interface{})
	if !hasRC {
		return
	}
	rcLoc := readerLoc + "/ReaderConfig"

	if inputType, ok := rc["InputType"].(string); ok {
		switch inputType {
		case "CSV", "JSON", "JSONL", "PARQUET", "MANIFEST":
		default:
			v.schemaf(rcLoc+"/InputType", "ReaderConfig InputType '%s' is not one of CSV, JSON, JSONL, PARQUET or MANIFEST", inputType)
		}
	}
	if csvHeaderLocation, ok := rc["CSVHeaderLocation"].(string); ok {
		if csvHeaderLocation != "FIRST_ROW" && csvHeaderLocation != "GIVEN" {
			v.schemaf(rcLoc+"/CSVHeaderLocation", "ReaderConfig CSVHeaderLocation must be FIRST_ROW or GIVEN")
		}
		if csvHeaderLocation == "GIVEN" {
			headers, hasHeaders := rc["CSVHeaders"].([]interface{})
			if !hasHeaders || len(headers) == 0 {
				v.schemaf(rcLoc+"/CSVHeaders", "ReaderConfig CSVHeaders is required when CSVHeaderLocation is GIVEN")
			} else {
				total := 0
				for _, h := range headers {
					if s, ok := h.(string); ok {
						total += len(s)
					}
				}
				if total > sfnstore.MaxCSVHeaderBytes {
					v.schemaf(rcLoc+"/CSVHeaders", "ReaderConfig CSVHeaders exceed the %d byte header ceiling", sfnstore.MaxCSVHeaderBytes)
				}
			}
		}
	}
	if delimiter, ok := rc["CSVDelimiter"].(string); ok {
		switch delimiter {
		case "COMMA", "PIPE", "SEMICOLON", "SPACE", "TAB":
		default:
			v.schemaf(rcLoc+"/CSVDelimiter", "ReaderConfig CSVDelimiter '%s' is not one of COMMA, PIPE, SEMICOLON, SPACE or TAB", delimiter)
		}
	}
	_, hasMaxItems := rc["MaxItems"]
	maxItemsPath, hasMaxItemsPath := rc["MaxItemsPath"].(string)
	if hasMaxItems && hasMaxItemsPath {
		v.schemaf(rcLoc+"/MaxItems", "ReaderConfig cannot specify both MaxItems and MaxItemsPath")
	}
	if hasMaxItems {
		if value, ok := rc["MaxItems"].(float64); ok {
			if value != math.Trunc(value) || value < 0 || value > float64(sfnstore.MaxItemReaderItems) {
				v.schemaf(rcLoc+"/MaxItems", "ReaderConfig MaxItems must be an integer from 0 to %d", sfnstore.MaxItemReaderItems)
			}
		} else {
			v.schemaf(rcLoc+"/MaxItems", "ReaderConfig MaxItems must be an integer")
		}
	}
	if hasMaxItemsPath && maxItemsPath == "" {
		v.schemaf(rcLoc+"/MaxItemsPath", "ReaderConfig MaxItemsPath must be a non-empty path string")
	}
	if pointer, ok := rc["ItemsPointer"].(string); ok && pointer != "" {
		if inputType, _ := rc["InputType"].(string); inputType != "JSON" {
			v.schemaf(rcLoc+"/ItemsPointer", "ReaderConfig ItemsPointer can only be specified when InputType is JSON")
		}
	}
	if transformation, ok := rc["Transformation"].(string); ok {
		if transformation != "NONE" && transformation != "LOAD_AND_FLATTEN" {
			v.schemaf(rcLoc+"/Transformation", "ReaderConfig Transformation must be NONE or LOAD_AND_FLATTEN")
		}
		if transformation == "LOAD_AND_FLATTEN" {
			if inputType, _ := rc["InputType"].(string); inputType == "" {
				v.schemaf(rcLoc+"/InputType", "InputType is required when Transformation is LOAD_AND_FLATTEN")
			}
		}
	}
	if manifestType, ok := rc["ManifestType"].(string); ok {
		if manifestType != "ATHENA_DATA" && manifestType != "S3_INVENTORY" {
			v.schemaf(rcLoc+"/ManifestType", "ReaderConfig ManifestType must be ATHENA_DATA or S3_INVENTORY")
		}
		if manifestType == "S3_INVENTORY" {
			if _, hasInputType := rc["InputType"]; hasInputType {
				v.schemaf(rcLoc+"/InputType", "InputType cannot be specified when ManifestType is S3_INVENTORY")
			}
		}
		if manifestType == "ATHENA_DATA" {
			if inputType, _ := rc["InputType"].(string); inputType == "" {
				v.schemaf(rcLoc+"/InputType", "InputType is required when ManifestType is ATHENA_DATA")
			}
		}
	}
}

// validateItemBatcher validates the ItemBatcher structure: the literal and
// reference-path forms are mutually exclusive within each pair, at least
// one sizing value is required, item counts are positive integers, the
// batch byte cap stays within the 256 KiB child-execution input bound and
// the fixed BatchInput is an object (ItemBatcher documentation).
func (v *aslValidatorContext) validateItemBatcher(name string, stateMap map[string]interface{}, stateLoc string, jsonata bool) {
	raw, ok := stateMap["ItemBatcher"].(map[string]interface{})
	if !ok {
		return
	}
	batcherLoc := stateLoc + "/ItemBatcher"

	_, hasMaxItems := raw["MaxItemsPerBatch"]
	maxItemsPath, hasMaxItemsPath := raw["MaxItemsPerBatchPath"].(string)
	if hasMaxItems && hasMaxItemsPath {
		v.schemaf(batcherLoc+"/MaxItemsPerBatch", "ItemBatcher cannot specify both MaxItemsPerBatch and MaxItemsPerBatchPath")
	}
	if hasMaxItems {
		if value, ok := raw["MaxItemsPerBatch"].(float64); ok {
			if value != math.Trunc(value) || value < 1 {
				v.schemaf(batcherLoc+"/MaxItemsPerBatch", "ItemBatcher MaxItemsPerBatch must be a positive integer")
			}
		} else if expr, isString := raw["MaxItemsPerBatch"].(string); !isString || !jsonata || !IsExpression(expr) {
			// "For JSONata-based states, you can also provide a JSONata
			// expression that evaluates to a positive integer."
			v.schemaf(batcherLoc+"/MaxItemsPerBatch", "ItemBatcher MaxItemsPerBatch must be a positive integer")
		}
	}
	if hasMaxItemsPath && maxItemsPath == "" {
		v.schemaf(batcherLoc+"/MaxItemsPerBatchPath", "ItemBatcher MaxItemsPerBatchPath must be a non-empty reference path")
	}

	_, hasMaxBytes := raw["MaxInputBytesPerBatch"]
	maxBytesPath, hasMaxBytesPath := raw["MaxInputBytesPerBatchPath"].(string)
	if hasMaxBytes && hasMaxBytesPath {
		v.schemaf(batcherLoc+"/MaxInputBytesPerBatch", "ItemBatcher cannot specify both MaxInputBytesPerBatch and MaxInputBytesPerBatchPath")
	}
	if hasMaxBytes {
		if value, ok := raw["MaxInputBytesPerBatch"].(float64); ok {
			if value != math.Trunc(value) || value < 1 || value > float64(sfnstore.MaxExecutionDataBytes) {
				v.schemaf(batcherLoc+"/MaxInputBytesPerBatch", "ItemBatcher MaxInputBytesPerBatch must be an integer from 1 to %d", sfnstore.MaxExecutionDataBytes)
			}
		} else if expr, isString := raw["MaxInputBytesPerBatch"].(string); !isString || !jsonata || !IsExpression(expr) {
			v.schemaf(batcherLoc+"/MaxInputBytesPerBatch", "ItemBatcher MaxInputBytesPerBatch must be a positive integer")
		}
	}
	if hasMaxBytesPath && maxBytesPath == "" {
		v.schemaf(batcherLoc+"/MaxInputBytesPerBatchPath", "ItemBatcher MaxInputBytesPerBatchPath must be a non-empty reference path")
	}

	if !hasMaxItems && !hasMaxItemsPath && !hasMaxBytes && !hasMaxBytesPath {
		v.schemaf(batcherLoc, "ItemBatcher must specify MaxItemsPerBatch, MaxInputBytesPerBatch or both to batch items")
	}

	_, hasBatchInput := raw["BatchInput"]
	_, hasBatchInputPath := raw["BatchInputPath"]
	if hasBatchInput && hasBatchInputPath {
		v.schemaf(batcherLoc+"/BatchInput", "ItemBatcher cannot specify both BatchInput and BatchInputPath")
	}
	if hasBatchInput {
		// "For JSONata-based states, you can provide JSONata expressions
		// directly to BatchInput, or use JSONata expressions inside JSON
		// objects or arrays."
		if _, isObject := raw["BatchInput"].(map[string]interface{}); !isObject {
			expr, isString := raw["BatchInput"].(string)
			if !isString || !jsonata || !IsExpression(expr) {
				v.schemaf(batcherLoc+"/BatchInput", "ItemBatcher BatchInput must be a JSON object")
			}
		}
	}
}

// validateResultWriter validates the documented required field
// combinations: WriterConfig alone, Resource with Parameters, or all
// three.
func (v *aslValidatorContext) validateResultWriter(name string, stateMap map[string]interface{}, stateLoc string) {
	raw, ok := stateMap["ResultWriter"].(map[string]interface{})
	if !ok {
		return
	}
	writerLoc := stateLoc + "/ResultWriter"

	_, hasResource := raw["Resource"]
	_, hasParams := raw["Parameters"]
	_, hasArgs := raw["Arguments"]
	_, hasWriterConfig := raw["WriterConfig"]

	if hasResource && !hasParams && !hasArgs && !hasWriterConfig {
		v.schemaf(writerLoc, "The ResultWriter of Map state '%s' must specify Parameters with Resource", name)
	}
	if (hasParams || hasArgs) && !hasResource && !hasWriterConfig {
		v.schemaf(writerLoc, "The ResultWriter of Map state '%s' must specify a Resource with Parameters", name)
	}
	if hasParams && hasArgs {
		v.schemaf(writerLoc, "The ResultWriter of Map state '%s' cannot specify both Parameters and Arguments", name)
	}
	if resource, ok := raw["Resource"].(string); ok && resource != "" &&
		!strings.HasSuffix(resource, ":s3:putObject") {
		v.schemaf(writerLoc+"/Resource", "The ResultWriter Resource '%s' is not the s3:putObject writer", resource)
	}
	if wc, ok := raw["WriterConfig"].(map[string]interface{}); ok {
		if transformation, ok := wc["Transformation"].(string); ok {
			if transformation != "NONE" && transformation != "COMPACT" && transformation != "FLATTEN" {
				v.schemaf(writerLoc+"/WriterConfig/Transformation", "WriterConfig Transformation must be NONE, COMPACT or FLATTEN")
			}
		}
		if outputType, ok := wc["OutputType"].(string); ok {
			if outputType != "JSON" && outputType != "JSONL" {
				v.schemaf(writerLoc+"/WriterConfig/OutputType", "WriterConfig OutputType must be JSON or JSONL")
			}
		}
	}
}

// validateToleratedFailure validates the threshold fields: non-negative
// counts, percentages within zero to one hundred, and the value/path
// exclusivity.
func (v *aslValidatorContext) validateToleratedFailure(name string, stateMap map[string]interface{}, stateLoc string, jsonata bool) {
	if _, ok := stateMap["ToleratedFailureCount"]; ok {
		v.validateCountField(stateMap, "ToleratedFailureCount", stateLoc, jsonata,
			"Map state ToleratedFailureCount must be a non-negative integer", nil)
	}
	if _, ok := stateMap["ToleratedFailurePercentage"]; ok {
		v.validateCountField(stateMap, "ToleratedFailurePercentage", stateLoc, jsonata,
			"Map state ToleratedFailurePercentage must be between zero and 100",
			func(value float64) bool { return value < 0 || value > 100 })
	}
}

// validateCountField applies the shape every Map count field shares: a
// non-negative whole number, or in a JSONata state a {% %} expression
// evaluating to one; the Field/FieldPath twins are mutually exclusive. The
// optional rangeCheck narrows the numeric bound (the percentage fields cap
// at one hundred); the field/path exclusivity message names the field.
func (v *aslValidatorContext) validateCountField(stateMap map[string]interface{}, field, stateLoc string, jsonata bool, message string, rangeCheck func(float64) bool) {
	raw := stateMap[field]
	value, isNumber := raw.(float64)
	outOfRange := !isNumber || value != math.Trunc(value) || value < 0
	if rangeCheck != nil && isNumber {
		outOfRange = rangeCheck(value)
	}
	if expr, isString := raw.(string); isString {
		// "In JSONata states, you can specify a JSONata expression that
		// evaluates to an integer."
		if !jsonata || !IsExpression(expr) {
			v.schemaf(stateLoc+"/"+field, "%s", message)
		}
	} else if outOfRange {
		v.schemaf(stateLoc+"/"+field, "%s", message)
	}
	if _, also := stateMap[field+"Path"]; also {
		v.schemaf(stateLoc+"/"+field, "Map state cannot specify both %s and %sPath", field, field)
	}
}
