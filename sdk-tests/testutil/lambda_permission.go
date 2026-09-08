package testutil

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaPermissionTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	funcName := tc.unique("PermFunc")
	roleARN, cleanupRole, err := tc.createRole(tc.unique("PermRole"))
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Permission_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer cleanupRole()

	functionARN, cleanupFn, err := tc.createFunction(funcName, roleARN, lambdaFunctionCode)
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Permission_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer cleanupFn()

	var addPermStatementID string

	results = append(results, tc.r.RunTest("lambda", "AddPermission", func() error {
		addPermStatementID = fmt.Sprintf("stmt-%d", time.Now().UnixNano())
		resp, err := tc.client.AddPermission(tc.ctx, &lambda.AddPermissionInput{
			FunctionName: aws.String(funcName),
			StatementId:  aws.String(addPermStatementID),
			Action:       aws.String("lambda:InvokeFunction"),
			Principal:    aws.String("apigateway.amazonaws.com"),
		})
		if err != nil {
			return err
		}
		if resp.Statement == nil || *resp.Statement == "" {
			return fmt.Errorf("Statement is nil or empty")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetPolicy", func() error {
		resp, err := tc.client.GetPolicy(tc.ctx, &lambda.GetPolicyInput{
			FunctionName: aws.String(funcName),
		})
		if err != nil {
			return err
		}
		if resp.Policy == nil || *resp.Policy == "" {
			return fmt.Errorf("policy is empty")
		}
		var policy map[string]interface{}
		if err := json.Unmarshal([]byte(*resp.Policy), &policy); err != nil {
			return fmt.Errorf("policy is not valid JSON: %v", err)
		}
		statements, ok := policy["Statement"].([]interface{})
		if !ok || len(statements) == 0 {
			return fmt.Errorf("policy has no Statement array")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "RemovePermission", func() error {
		statementID := fmt.Sprintf("stmt-%d", time.Now().UnixNano())
		_, err := tc.client.AddPermission(tc.ctx, &lambda.AddPermissionInput{
			FunctionName: aws.String(funcName),
			StatementId:  aws.String(statementID),
			Action:       aws.String("lambda:InvokeFunction"),
			Principal:    aws.String("apigateway.amazonaws.com"),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.RemovePermission(tc.ctx, &lambda.RemovePermissionInput{
			FunctionName: aws.String(funcName),
			StatementId:  aws.String(statementID),
		})
		if err != nil {
			return err
		}
		policyResp, err := tc.client.GetPolicy(tc.ctx, &lambda.GetPolicyInput{
			FunctionName: aws.String(funcName),
		})
		if err != nil {
			return fmt.Errorf("GetPolicy after remove: %v", err)
		}
		var policy map[string]interface{}
		if err := json.Unmarshal([]byte(*policyResp.Policy), &policy); err != nil {
			return fmt.Errorf("policy is not valid JSON: %v", err)
		}
		statements, _ := policy["Statement"].([]interface{})
		for _, s := range statements {
			stmt, _ := s.(map[string]interface{})
			if sid, _ := stmt["Sid"].(string); sid == statementID {
				return fmt.Errorf("removed statement %s still present in policy", statementID)
			}
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "AddPermission_ConditionAndQualifierScope", func() error {
		if _, err := tc.client.PublishVersion(tc.ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(funcName),
		}); err != nil {
			return fmt.Errorf("publish: %v", err)
		}
		sid := fmt.Sprintf("scoped-%d", time.Now().UnixNano())
		sourceArn := fmt.Sprintf("arn:aws:s3:::source-bucket-%d", time.Now().UnixNano())
		_, err := tc.client.AddPermission(tc.ctx, &lambda.AddPermissionInput{
			FunctionName:        aws.String(funcName),
			StatementId:         aws.String(sid),
			Action:              aws.String("lambda:InvokeFunction"),
			Principal:           aws.String("s3.amazonaws.com"),
			SourceArn:           aws.String(sourceArn),
			SourceAccount:       aws.String(tc.r.accountID),
			Qualifier:           aws.String("1"),
			FunctionUrlAuthType: types.FunctionUrlAuthTypeNone,
		})
		if err != nil {
			return err
		}

		policyResp, err := tc.client.GetPolicy(tc.ctx, &lambda.GetPolicyInput{
			FunctionName: aws.String(funcName),
		})
		if err != nil {
			return err
		}
		var policy struct {
			Statement []struct {
				Sid       string                       `json:"Sid"`
				Resource  string                       `json:"Resource"`
				Condition map[string]map[string]string `json:"Condition"`
			} `json:"Statement"`
		}
		if err := json.Unmarshal([]byte(*policyResp.Policy), &policy); err != nil {
			return fmt.Errorf("policy is not valid JSON: %v", err)
		}
		for _, stmt := range policy.Statement {
			if stmt.Sid != sid {
				continue
			}
			if !strings.HasSuffix(stmt.Resource, funcName+":1") {
				return fmt.Errorf("Resource should be scoped with the :1 qualifier, got %s", stmt.Resource)
			}
			if stmt.Condition["StringLike"]["aws:SourceArn"] != sourceArn {
				return fmt.Errorf("Condition should carry aws:SourceArn via StringLike, got %v", stmt.Condition)
			}
			if stmt.Condition["StringEquals"]["aws:SourceAccount"] != tc.r.accountID {
				return fmt.Errorf("Condition should carry aws:SourceAccount via StringEquals, got %v", stmt.Condition)
			}
			if stmt.Condition["StringEquals"]["lambda:FunctionUrlAuthType"] != "NONE" {
				return fmt.Errorf("Condition should carry lambda:FunctionUrlAuthType, got %v", stmt.Condition)
			}
			return nil
		}
		return fmt.Errorf("statement %s not found in policy", sid)
	}))

	results = append(results, tc.r.RunTest("lambda", "ResourcePolicy_PutGetDelete", func() error {
		sid := fmt.Sprintf("rbp-%d", time.Now().UnixNano())
		doc := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Sid":"%s","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::%s:root"},"Action":"lambda:InvokeFunction","Resource":"%s"}]}`,
			sid, tc.r.accountID, functionARN)

		put, err := tc.client.PutResourcePolicy(tc.ctx, &lambda.PutResourcePolicyInput{
			ResourceArn: aws.String(functionARN),
			Policy:      aws.String(doc),
		})
		if err != nil {
			return err
		}
		if put.RevisionId == nil || *put.RevisionId == "" {
			return fmt.Errorf("PutResourcePolicy returned no RevisionId")
		}
		if put.Policy == nil || *put.Policy == "" {
			return fmt.Errorf("PutResourcePolicy returned no Policy")
		}

		got, err := tc.client.GetResourcePolicy(tc.ctx, &lambda.GetResourcePolicyInput{
			ResourceArn: aws.String(functionARN),
		})
		if err != nil {
			return err
		}
		if got.RevisionId == nil || *got.RevisionId != *put.RevisionId {
			return fmt.Errorf("GetResourcePolicy RevisionId %v, want the Put revision %v", got.RevisionId, put.RevisionId)
		}
		var gotDoc struct {
			Statement []struct {
				Sid       string `json:"Sid"`
				Effect    string `json:"Effect"`
				Principal struct {
					AWS string `json:"AWS"`
				} `json:"Principal"`
			} `json:"Statement"`
		}
		if err := json.Unmarshal([]byte(*got.Policy), &gotDoc); err != nil {
			return fmt.Errorf("GetResourcePolicy policy is not valid JSON: %v", err)
		}
		if len(gotDoc.Statement) != 1 || gotDoc.Statement[0].Sid != sid {
			return fmt.Errorf("GetResourcePolicy statements %v, want exactly %s", gotDoc.Statement, sid)
		}
		if gotDoc.Statement[0].Effect != "Allow" {
			return fmt.Errorf("statement Effect %q, want Allow", gotDoc.Statement[0].Effect)
		}
		if gotDoc.Statement[0].Principal.AWS != "arn:aws:iam::"+tc.r.accountID+":root" {
			return fmt.Errorf("statement Principal %v, want the account root ARN", gotDoc.Statement[0].Principal)
		}

		// The whole-document replacement also replaces permissions added
		// through AddPermission, so GetPolicy serves the same statement.
		oldView, err := tc.client.GetPolicy(tc.ctx, &lambda.GetPolicyInput{
			FunctionName: aws.String(funcName),
		})
		if err != nil {
			return fmt.Errorf("GetPolicy after PutResourcePolicy: %v", err)
		}
		if !strings.Contains(*oldView.Policy, sid) {
			return fmt.Errorf("GetPolicy does not serve the PutResourcePolicy statement: %s", *oldView.Policy)
		}

		// A stale revision cannot replace or delete the policy.
		_, err = tc.client.PutResourcePolicy(tc.ctx, &lambda.PutResourcePolicyInput{
			ResourceArn: aws.String(functionARN),
			Policy:      aws.String(doc),
			RevisionId:  aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err := expectAWSErrorCode(err, "PreconditionFailedException"); err != nil {
			return err
		}
		_, err = tc.client.DeleteResourcePolicy(tc.ctx, &lambda.DeleteResourcePolicyInput{
			ResourceArn: aws.String(functionARN),
			RevisionId:  aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err := expectAWSErrorCode(err, "PreconditionFailedException"); err != nil {
			return err
		}

		if _, err := tc.client.DeleteResourcePolicy(tc.ctx, &lambda.DeleteResourcePolicyInput{
			ResourceArn: aws.String(functionARN),
		}); err != nil {
			return fmt.Errorf("DeleteResourcePolicy: %v", err)
		}
		_, err = tc.client.GetResourcePolicy(tc.ctx, &lambda.GetResourcePolicyInput{
			ResourceArn: aws.String(functionARN),
		})
		return expectAWSErrorCode(err, "ResourceNotFoundException")
	}))

	results = append(results, tc.r.RunTest("lambda", "ResourcePolicy_Validation", func() error {
		_, err := tc.client.PutResourcePolicy(tc.ctx, &lambda.PutResourcePolicyInput{
			ResourceArn: aws.String(functionARN),
			Policy:      aws.String(`not json`),
		})
		if err := expectAWSErrorCode(err, "InvalidParameterValueException"); err != nil {
			return fmt.Errorf("malformed document: %v", err)
		}

		_, err = tc.client.PutResourcePolicy(tc.ctx, &lambda.PutResourcePolicyInput{
			ResourceArn: aws.String(functionARN),
			Policy:      aws.String(`{"Version":"2012-10-17"}`),
		})
		if err := expectAWSErrorCode(err, "InvalidParameterValueException"); err != nil {
			return fmt.Errorf("statement-less document: %v", err)
		}

		_, err = tc.client.PutResourcePolicy(tc.ctx, &lambda.PutResourcePolicyInput{
			ResourceArn: aws.String(functionARN),
			Policy:      aws.String(`{"Version":"2012-10-17","Statement":[{"Sid":"s","Principal":"*","Action":"lambda:InvokeFunction","Resource":"arn"}]}`),
		})
		if err := expectAWSErrorCode(err, "InvalidParameterValueException"); err != nil {
			return fmt.Errorf("statement missing Effect: %v", err)
		}

		big := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Sid":"big","Effect":"Allow","Principal":"*","Action":"lambda:InvokeFunction","Resource":"%s"}]}`,
			strings.Repeat("a", 21000))
		_, err = tc.client.PutResourcePolicy(tc.ctx, &lambda.PutResourcePolicyInput{
			ResourceArn: aws.String(functionARN),
			Policy:      aws.String(big),
		})
		if err := expectAWSErrorCode(err, "PolicyLengthExceededException"); err != nil {
			return fmt.Errorf("oversized document: %v", err)
		}

		missingArn := fmt.Sprintf("arn:aws:lambda:us-east-1:%s:function:no-such-function-%d", tc.r.accountID, time.Now().UnixNano())
		_, err = tc.client.GetResourcePolicy(tc.ctx, &lambda.GetResourcePolicyInput{
			ResourceArn: aws.String(missingArn),
		})
		if err := expectAWSErrorCode(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("get on unknown function: %v", err)
		}
		_, err = tc.client.PutResourcePolicy(tc.ctx, &lambda.PutResourcePolicyInput{
			ResourceArn: aws.String(missingArn),
			Policy:      aws.String(`{"Version":"2012-10-17","Statement":[{"Sid":"s","Effect":"Allow","Principal":"*","Action":"lambda:InvokeFunction","Resource":"arn"}]}`),
		})
		if err := expectAWSErrorCode(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("put on unknown function: %v", err)
		}

		_, err = tc.client.PutResourcePolicy(tc.ctx, &lambda.PutResourcePolicyInput{
			ResourceArn: aws.String("not-an-arn"),
			Policy:      aws.String(`{"Version":"2012-10-17","Statement":[{"Sid":"s","Effect":"Allow","Principal":"*","Action":"lambda:InvokeFunction","Resource":"arn"}]}`),
		})
		return expectAWSErrorCode(err, "InvalidParameterValueException")
	}))

	results = append(results, tc.r.RunTest("lambda", "TagResource", func() error {
		_, err := tc.client.TagResource(tc.ctx, &lambda.TagResourceInput{
			Resource: aws.String(functionARN),
			Tags: map[string]string{
				"Environment": "test",
				"Project":     "sdk-tests",
			},
		})
		return err
	}))

	results = append(results, tc.r.RunTest("lambda", "ListTags", func() error {
		resp, err := tc.client.ListTags(tc.ctx, &lambda.ListTagsInput{
			Resource: aws.String(functionARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		if resp.Tags["Project"] != "sdk-tests" {
			return fmt.Errorf("expected tag Project=sdk-tests, got %v", resp.Tags["Project"])
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UntagResource", func() error {
		_, err := tc.client.UntagResource(tc.ctx, &lambda.UntagResourceInput{
			Resource: aws.String(functionARN),
			TagKeys:  []string{"Environment"},
		})
		if err != nil {
			return err
		}
		tagResp, err := tc.client.ListTags(tc.ctx, &lambda.ListTagsInput{
			Resource: aws.String(functionARN),
		})
		if err != nil {
			return fmt.Errorf("ListTags after untag: %v", err)
		}
		if _, ok := tagResp.Tags["Environment"]; ok {
			return fmt.Errorf("tag 'Environment' should have been removed")
		}
		if tagResp.Tags["Project"] != "sdk-tests" {
			return fmt.Errorf("tag 'Project' should still exist, got %v", tagResp.Tags["Project"])
		}
		return nil
	}))

	return results
}
