package testutil

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssfn "github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

// runSFNContractTests pins the wire contracts re-established by the SFN
// Core remediation: execution input validation, EXPRESS-only synchronous
// starts, the alias ARN format with per-state-machine namespaces,
// qualified-ARN starts, redrive eligibility, tag targets and the
// includedData behaviour.
func (r *TestRunner) runSFNContractTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	smARN, err := tc.createPassSM(fmt.Sprintf("Contract-%d", time.Now().UnixNano()), "contract")
	if err != nil {
		return []TestResult{{Service: "stepfunctions", TestName: "ContractSetup", Status: "FAIL", Error: fmt.Sprintf("create SM: %v", err)}}
	}
	defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(smARN)})

	results = append(results, r.RunTest("stepfunctions", "StartExecution_InvalidJSONInput", func() error {
		_, err := tc.client.StartExecution(tc.ctx, &awssfn.StartExecutionInput{
			StateMachineArn: aws.String(smARN),
			Input:           aws.String("{not json"),
		})
		return expectAWSErrorCode(err, "InvalidExecutionInput")
	}))

	results = append(results, r.RunTest("stepfunctions", "StartExecution_NameForbiddenCharacters", func() error {
		_, err := tc.client.StartExecution(tc.ctx, &awssfn.StartExecutionInput{
			StateMachineArn: aws.String(smARN),
			Name:            aws.String("bad name!"),
		})
		return expectAWSErrorCode(err, "InvalidName")
	}))

	results = append(results, r.RunTest("stepfunctions", "StartSyncExecution_StandardRejected", func() error {
		// The SDK resolves the synchronous operation onto a sync-prefixed
		// endpoint, so the rejection is asserted over the raw JSON-1.0
		// protocol the router serves.
		status, errBody, err := tc.rawJSONCall("AWSStepFunctions.StartSyncExecution", map[string]interface{}{"stateMachineArn": smARN})
		if err != nil {
			return err
		}
		if status != http.StatusBadRequest {
			return fmt.Errorf("status %d for STANDARD sync start, want 400: %v", status, errBody)
		}
		if errType, _ := errBody["__type"].(string); !strings.HasSuffix(errType, "StateMachineTypeNotSupported") {
			return fmt.Errorf("error type %q, want StateMachineTypeNotSupported", errType)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "UpdateStateMachine_MissingRequiredParameter", func() error {
		_, err := tc.client.UpdateStateMachine(tc.ctx, &awssfn.UpdateStateMachineInput{
			StateMachineArn: aws.String(smARN),
		})
		return expectAWSErrorCode(err, "MissingRequiredParameter")
	}))

	results = append(results, r.RunTest("stepfunctions", "ListMapRuns_UnknownExecution", func() error {
		_, err := tc.client.ListMapRuns(tc.ctx, &awssfn.ListMapRunsInput{
			ExecutionArn: aws.String("arn:aws:states:" + r.region + ":" + r.accountID + ":execution:no-such-sm:no-such-exec"),
		})
		return expectAWSErrorCode(err, "ExecutionDoesNotExist")
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeStateMachine_MetadataOnly", func() error {
		resp, err := tc.client.DescribeStateMachine(tc.ctx, &awssfn.DescribeStateMachineInput{
			StateMachineArn: aws.String(smARN),
			IncludedData:    types.IncludedDataMetadataOnly,
		})
		if err != nil {
			return err
		}
		if resp.Definition == nil || *resp.Definition != "{}" {
			return fmt.Errorf("METADATA_ONLY definition = %v, want {}", resp.Definition)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "SendTaskSuccess_UnknownToken", func() error {
		_, err := tc.client.SendTaskSuccess(tc.ctx, &awssfn.SendTaskSuccessInput{
			TaskToken: aws.String("no-such-token"),
			Output:    aws.String("{}"),
		})
		return expectAWSErrorCode(err, "TaskDoesNotExist")
	}))

	results = append(results, r.RunTest("stepfunctions", "SendTaskSuccess_InvalidOutput", func() error {
		_, err := tc.client.SendTaskSuccess(tc.ctx, &awssfn.SendTaskSuccessInput{
			TaskToken: aws.String("no-such-token"),
			Output:    aws.String("{not json"),
		})
		return expectAWSErrorCode(err, "InvalidOutput")
	}))

	results = append(results, r.RunTest("stepfunctions", "CreateActivity_DuplicateName", func() error {
		name := fmt.Sprintf("ContractAct-%d", time.Now().UnixNano())
		if _, err := tc.client.CreateActivity(tc.ctx, &awssfn.CreateActivityInput{Name: aws.String(name)}); err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.client.DeleteActivity(tc.ctx, &awssfn.DeleteActivityInput{ActivityArn: aws.String("arn:aws:states:" + r.region + ":" + r.accountID + ":activity:" + name)})
		_, err := tc.client.CreateActivity(tc.ctx, &awssfn.CreateActivityInput{Name: aws.String(name)})
		return expectAWSErrorCode(err, "ActivityAlreadyExists")
	}))

	results = append(results, r.RunTest("stepfunctions", "CreateActivity_NameForbiddenCharacters", func() error {
		_, err := tc.client.CreateActivity(tc.ctx, &awssfn.CreateActivityInput{Name: aws.String("bad name!")})
		return expectAWSErrorCode(err, "InvalidName")
	}))

	results = append(results, r.RunTest("stepfunctions", "TagResource_VersionArnRejected", func() error {
		// Versions are not taggable resources; a qualified ARN is rejected.
		versionArn := smARN + ":1"
		_, err := tc.client.TagResource(tc.ctx, &awssfn.TagResourceInput{
			ResourceArn: aws.String(versionArn),
			Tags:        []types.Tag{{Key: aws.String("k"), Value: aws.String("v")}},
		})
		return expectAWSErrorCode(err, "ResourceNotFound")
	}))

	results = append(results, r.RunTest("stepfunctions", "AliasARNFormatAndNamespace", func() error {
		// Second state machine so the same alias name can exist on both,
		// proving the per-state-machine namespace of the AWS alias ARN.
		otherARN, err := tc.createPassSM(fmt.Sprintf("ContractNS-%d", time.Now().UnixNano()), "namespace")
		if err != nil {
			return fmt.Errorf("create second SM: %v", err)
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(otherARN)})

		mkVersion := func(arn string) (string, error) {
			resp, err := tc.client.PublishStateMachineVersion(tc.ctx, &awssfn.PublishStateMachineVersionInput{
				StateMachineArn: aws.String(arn),
			})
			if err != nil {
				return "", err
			}
			return *resp.StateMachineVersionArn, nil
		}

		v1, err := mkVersion(smARN)
		if err != nil {
			return fmt.Errorf("publish v1: %v", err)
		}
		v2, err := mkVersion(otherARN)
		if err != nil {
			return fmt.Errorf("publish v2: %v", err)
		}

		aliasName := fmt.Sprintf("Shared-%d", time.Now().UnixNano())
		first, err := tc.client.CreateStateMachineAlias(tc.ctx, &awssfn.CreateStateMachineAliasInput{
			Name: aws.String(aliasName),
			RoutingConfiguration: []types.RoutingConfigurationListItem{
				{StateMachineVersionArn: aws.String(v1), Weight: 100},
			},
		})
		if err != nil {
			return fmt.Errorf("first alias: %v", err)
		}
		defer tc.client.DeleteStateMachineAlias(tc.ctx, &awssfn.DeleteStateMachineAliasInput{StateMachineAliasArn: first.StateMachineAliasArn})

		wantPrefix := smARN + ":"
		if !strings.HasPrefix(*first.StateMachineAliasArn, wantPrefix) {
			return fmt.Errorf("alias ARN %q must extend the state machine ARN with a colon (AWS format stateMachine:<name>:<alias>)", *first.StateMachineAliasArn)
		}

		second, err := tc.client.CreateStateMachineAlias(tc.ctx, &awssfn.CreateStateMachineAliasInput{
			Name: aws.String(aliasName),
			RoutingConfiguration: []types.RoutingConfigurationListItem{
				{StateMachineVersionArn: aws.String(v2), Weight: 100},
			},
		})
		if err != nil {
			return fmt.Errorf("same alias name on another state machine must be accepted (per-SM namespace): %v", err)
		}
		defer tc.client.DeleteStateMachineAlias(tc.ctx, &awssfn.DeleteStateMachineAliasInput{StateMachineAliasArn: second.StateMachineAliasArn})
		if *second.StateMachineAliasArn == *first.StateMachineAliasArn {
			return fmt.Errorf("alias ARNs of different state machines must differ")
		}

		// A version an alias routes to cannot be deleted.
		_, delErr := tc.client.DeleteStateMachineVersion(tc.ctx, &awssfn.DeleteStateMachineVersionInput{
			StateMachineVersionArn: aws.String(v1),
		})
		if codeErr := expectAWSErrorCode(delErr, "ConflictException"); codeErr != nil {
			return codeErr
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "StartExecution_WithVersionAndAliasARNs", func() error {
		versionResp, err := tc.client.PublishStateMachineVersion(tc.ctx, &awssfn.PublishStateMachineVersionInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return fmt.Errorf("publish version: %v", err)
		}
		versionArn := *versionResp.StateMachineVersionArn
		defer tc.client.DeleteStateMachineVersion(tc.ctx, &awssfn.DeleteStateMachineVersionInput{StateMachineVersionArn: aws.String(versionArn)})

		versionExec, err := tc.client.StartExecution(tc.ctx, &awssfn.StartExecutionInput{
			StateMachineArn: aws.String(versionArn),
		})
		if err != nil {
			return fmt.Errorf("start with version ARN: %v", err)
		}
		vd, err := tc.client.DescribeExecution(tc.ctx, &awssfn.DescribeExecutionInput{
			ExecutionArn: versionExec.ExecutionArn,
		})
		if err != nil {
			return fmt.Errorf("describe version-pinned execution: %v", err)
		}
		if vd.StateMachineVersionArn == nil || *vd.StateMachineVersionArn != versionArn {
			return fmt.Errorf("version-pinned execution must carry stateMachineVersionArn = %s, got %v", versionArn, vd.StateMachineVersionArn)
		}
		if vd.StateMachineAliasArn != nil {
			return fmt.Errorf("version-pinned execution must not carry stateMachineAliasArn")
		}

		aliasName := fmt.Sprintf("Exec-%d", time.Now().UnixNano())
		aliasResp, err := tc.client.CreateStateMachineAlias(tc.ctx, &awssfn.CreateStateMachineAliasInput{
			Name: aws.String(aliasName),
			RoutingConfiguration: []types.RoutingConfigurationListItem{
				{StateMachineVersionArn: aws.String(versionArn), Weight: 100},
			},
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		defer tc.client.DeleteStateMachineAlias(tc.ctx, &awssfn.DeleteStateMachineAliasInput{StateMachineAliasArn: aliasResp.StateMachineAliasArn})

		aliasExec, err := tc.client.StartExecution(tc.ctx, &awssfn.StartExecutionInput{
			StateMachineArn: aliasResp.StateMachineAliasArn,
		})
		if err != nil {
			return fmt.Errorf("start with alias ARN: %v", err)
		}
		ad, err := tc.client.DescribeExecution(tc.ctx, &awssfn.DescribeExecutionInput{
			ExecutionArn: aliasExec.ExecutionArn,
		})
		if err != nil {
			return fmt.Errorf("describe alias-pinned execution: %v", err)
		}
		if ad.StateMachineAliasArn == nil || *ad.StateMachineAliasArn != *aliasResp.StateMachineAliasArn {
			return fmt.Errorf("alias-pinned execution must carry stateMachineAliasArn, got %v", ad.StateMachineAliasArn)
		}
		if ad.StateMachineVersionArn == nil {
			return fmt.Errorf("alias-pinned execution must carry the routed stateMachineVersionArn")
		}

		// Plain starts carry neither association.
		plainExec, err := tc.client.StartExecution(tc.ctx, &awssfn.StartExecutionInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return fmt.Errorf("plain start: %v", err)
		}
		pd, err := tc.client.DescribeExecution(tc.ctx, &awssfn.DescribeExecutionInput{
			ExecutionArn: plainExec.ExecutionArn,
		})
		if err != nil {
			return fmt.Errorf("describe plain execution: %v", err)
		}
		if pd.StateMachineVersionArn != nil || pd.StateMachineAliasArn != nil {
			return fmt.Errorf("unqualified start must carry no version or alias association")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "CreateStateMachine_IdempotentRetry", func() error {
		name := fmt.Sprintf("IdemSM-%d", time.Now().UnixNano())
		def := `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`
		in := &awssfn.CreateStateMachineInput{
			Name:       aws.String(name),
			Definition: aws.String(def),
			RoleArn:    aws.String(tc.roleARN),
		}
		first, err := tc.client.CreateStateMachine(tc.ctx, in)
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: first.StateMachineArn})

		// A retry with the same parameters returns the original resource;
		// roleArn and tags differences are ignored. The differing role is
		// a real one — the role ARN is still validated on the wire.
		_, otherRoleARN, otherRoleCleanup := tc.createRoleForSM("IdemRole")
		defer otherRoleCleanup()
		retry := *in
		retry.RoleArn = aws.String(otherRoleARN)
		second, err := tc.client.CreateStateMachine(tc.ctx, &retry)
		if err != nil {
			return fmt.Errorf("identical retry must succeed idempotently: %v", err)
		}
		if *second.StateMachineArn != *first.StateMachineArn {
			return fmt.Errorf("retry ARN %s, want the original %s", *second.StateMachineArn, *first.StateMachineArn)
		}

		// The same name with a different definition stays a conflict.
		_, err = tc.client.CreateStateMachine(tc.ctx, &awssfn.CreateStateMachineInput{
			Name:       aws.String(name),
			Definition: aws.String(`{"StartAt":"B","States":{"B":{"Type":"Succeed"}}}`),
			RoleArn:    aws.String(tc.roleARN),
		})
		return expectAWSErrorCode(err, "StateMachineAlreadyExists")
	}))

	results = append(results, r.RunTest("stepfunctions", "CreateStateMachine_NameU10FFFFRejected", func() error {
		_, err := tc.client.CreateStateMachine(tc.ctx, &awssfn.CreateStateMachineInput{
			Name:       aws.String("bad\U0010FFFFname"),
			Definition: aws.String(`{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`),
			RoleArn:    aws.String(tc.roleARN),
		})
		return expectAWSErrorCode(err, "InvalidName")
	}))

	results = append(results, r.RunTest("stepfunctions", "UpdateStateMachine_IgnoresUnknownTypeAndRevisionId", func() error {
		// Neither member exists on UpdateStateMachineInput, so both are
		// ignored like any other unknown member; the request is asserted
		// over the raw JSON-1.0 protocol because the typed SDK input has
		// no fields to carry them.
		status, _, err := tc.rawJSONCall("AWSStepFunctions.UpdateStateMachine", map[string]interface{}{
			"stateMachineArn": smARN,
			"roleArn":         tc.roleARN,
			"type":            "EXPRESS",
			"revisionId":      "stale-revision",
		})
		if err != nil {
			return err
		}
		if status != http.StatusOK {
			return fmt.Errorf("status %d for update with unknown members, want 200", status)
		}

		desc, err := tc.client.DescribeStateMachine(tc.ctx, &awssfn.DescribeStateMachineInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return err
		}
		if desc.Type != types.StateMachineTypeStandard {
			return fmt.Errorf("state machine type mutated to %q, want STANDARD", desc.Type)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ValidateStateMachineDefinition_WarningOnlyOK", func() error {
		warningOnly := `{"StartAt":"P","States":{"P":{"Type":"Pass","Result":"$.payload","End":true}}}`
		resp, err := tc.client.ValidateStateMachineDefinition(tc.ctx, &awssfn.ValidateStateMachineDefinitionInput{
			Definition: aws.String(warningOnly),
			Severity:   types.ValidateStateMachineDefinitionSeverityWarning,
		})
		if err != nil {
			return err
		}
		if resp.Result != types.ValidateStateMachineDefinitionResultCodeOk {
			return fmt.Errorf("warning-only definition result = %s, want OK", resp.Result)
		}
		if len(resp.Diagnostics) == 0 {
			return fmt.Errorf("expected the warning diagnostic to be reported")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "PublishStateMachineVersion_IdempotentSameRevision", func() error {
		name := fmt.Sprintf("PubSM-%d", time.Now().UnixNano())
		sm, err := tc.createPassSM(name, "publish idempotency")
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(sm)})
		first, err := tc.client.PublishStateMachineVersion(tc.ctx, &awssfn.PublishStateMachineVersionInput{
			StateMachineArn: aws.String(sm),
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachineVersion(tc.ctx, &awssfn.DeleteStateMachineVersionInput{
			StateMachineVersionArn: first.StateMachineVersionArn,
		})
		second, err := tc.client.PublishStateMachineVersion(tc.ctx, &awssfn.PublishStateMachineVersionInput{
			StateMachineArn: aws.String(sm),
		})
		if err != nil {
			return fmt.Errorf("same-revision republish must be idempotent: %v", err)
		}
		if *second.StateMachineVersionArn != *first.StateMachineVersionArn {
			return fmt.Errorf("republish created %s, want the existing %s", *second.StateMachineVersionArn, *first.StateMachineVersionArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "PublishStateMachineVersion_RevisionMismatchConflict", func() error {
		name := fmt.Sprintf("PubConflictSM-%d", time.Now().UnixNano())
		sm, err := tc.createPassSM(name, "publish conflict")
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(sm)})
		_, err = tc.client.PublishStateMachineVersion(tc.ctx, &awssfn.PublishStateMachineVersionInput{
			StateMachineArn: aws.String(sm),
			RevisionId:      aws.String("stale-revision-id"),
		})
		return expectAWSErrorCode(err, "ConflictException")
	}))

	results = append(results, r.RunTest("stepfunctions", "UpdateStateMachineAlias_RequiresDescriptionOrRouting", func() error {
		name := fmt.Sprintf("AliasReqSM-%d", time.Now().UnixNano())
		sm, err := tc.createPassSM(name, "alias required parameters")
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &awssfn.DeleteStateMachineInput{StateMachineArn: aws.String(sm)})
		version, err := tc.client.PublishStateMachineVersion(tc.ctx, &awssfn.PublishStateMachineVersionInput{
			StateMachineArn: aws.String(sm),
		})
		if err != nil {
			return err
		}
		aliasName := fmt.Sprintf("req-alias-%d", time.Now().UnixNano())
		alias, err := tc.client.CreateStateMachineAlias(tc.ctx, &awssfn.CreateStateMachineAliasInput{
			Name: aws.String(aliasName),
			RoutingConfiguration: []types.RoutingConfigurationListItem{{
				StateMachineVersionArn: version.StateMachineVersionArn,
				Weight:                 100,
			}},
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteStateMachineAlias(tc.ctx, &awssfn.DeleteStateMachineAliasInput{
			StateMachineAliasArn: alias.StateMachineAliasArn,
		})

		_, err = tc.client.UpdateStateMachineAlias(tc.ctx, &awssfn.UpdateStateMachineAliasInput{
			StateMachineAliasArn: alias.StateMachineAliasArn,
		})
		return expectAWSErrorCode(err, "ValidationException")
	}))

	results = append(results, r.RunTest("stepfunctions", "ListExecutions_RedriveFilterWithStateMachineArnRejected", func() error {
		_, err := tc.client.ListExecutions(tc.ctx, &awssfn.ListExecutionsInput{
			StateMachineArn: aws.String(smARN),
			RedriveFilter:   types.ExecutionRedriveFilterRedriven,
		})
		return expectAWSErrorCode(err, "ValidationException")
	}))

	results = append(results, r.RunTest("stepfunctions", "TestState_InspectionDataCarriesResultMember", func() error {
		def := `{"StartAt":"F","States":{"F":{"Type":"Fail","Error":"CustomError","Cause":"deliberate"}}}`
		result, err := tc.rawTestState(map[string]interface{}{
			"definition":      def,
			"stateName":       "F",
			"input":           `{}`,
			"inspectionLevel": "DEBUG",
		})
		if err != nil {
			return err
		}
		inspection, ok := result["inspectionData"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("inspectionData missing or not an object: %v", result["inspectionData"])
		}
		if _, ok := inspection["result"]; !ok {
			return fmt.Errorf("inspectionData.result missing: %v", inspection)
		}
		if _, ok := inspection["output"]; ok {
			return fmt.Errorf("inspectionData carries the non-existent output member: %v", inspection)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "TestState_MockResult", func() error {
		def := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","Next":"Done"},"Done":{"Type":"Succeed"}}}`
		result, err := tc.rawTestState(map[string]interface{}{
			"definition": def,
			"stateName":  "T",
			"mock":       map[string]interface{}{"result": `{"value":42}`},
		})
		if err != nil {
			return err
		}
		if result["status"] != "SUCCEEDED" {
			return fmt.Errorf("status %v, want SUCCEEDED", result["status"])
		}
		if result["output"] != `{"value":42}` {
			return fmt.Errorf("output %v, want the mocked result", result["output"])
		}
		if result["nextState"] != "Done" {
			return fmt.Errorf("nextState %v, want Done", result["nextState"])
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "TestState_MockOnPassRejected", func() error {
		def := `{"StartAt":"P","States":{"P":{"Type":"Pass","Result":{"a":1},"End":true}}}`
		status, errBody, err := tc.rawJSONCall("AWSStepFunctions.TestState", map[string]interface{}{
			"definition": def,
			"stateName":  "P",
			"mock":       map[string]interface{}{"result": `{"value":42}`},
		})
		if err != nil {
			return err
		}
		if status != http.StatusBadRequest {
			return fmt.Errorf("status %d for a mock on a Pass state, want 400: %v", status, errBody)
		}
		if errType, _ := errBody["__type"].(string); !strings.HasSuffix(errType, "ValidationException") {
			return fmt.Errorf("error type %q, want ValidationException", errType)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "TestState_ContextRequiresMock", func() error {
		def := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","End":true}}}`
		status, errBody, err := tc.rawJSONCall("AWSStepFunctions.TestState", map[string]interface{}{
			"definition": def,
			"stateName":  "T",
			"context":    `{"Execution":{"Name":"x"}}`,
		})
		if err != nil {
			return err
		}
		if status != http.StatusBadRequest {
			return fmt.Errorf("status %d for context without a mock, want 400: %v", status, errBody)
		}
		if errType, _ := errBody["__type"].(string); !strings.HasSuffix(errType, "ValidationException") {
			return fmt.Errorf("error type %q, want ValidationException", errType)
		}
		return nil
	}))

	return results
}
