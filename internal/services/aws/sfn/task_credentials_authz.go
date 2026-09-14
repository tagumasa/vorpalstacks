package sfn

// This file implements the invocation half of a Task state's Credentials:
// deriving the IAM action and resource ARN each integration dispatch
// targets, and running the assume-role authorisation chain before the
// invocation leaves the executor.

import (
	"context"
	"encoding/json"
	"strings"

	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// TaskCredentialsAuthoriser is the authorisation seam for a Task state's
// Credentials role: it validates that the States service can assume the
// resolved role and that the assumed role's identity policies allow the
// integration action on the target resource. The server wiring supplies
// the implementation; a nil seam leaves task invocations on the
// machine-role behaviour.
type TaskCredentialsAuthoriser interface {
	AuthoriseTaskCredentials(ctx context.Context, taskRoleArn, action, resourceArn string) error
}

// integrationTarget is one action/resource pair an integration dispatch
// must be authorised for.
type integrationTarget struct {
	action      string
	resourceArn string
}

// dynamodbIntegrationActions maps the integration action segment to the
// IAM action the table operation authorises as.
var dynamodbIntegrationActions = map[string]string{
	"getItem":    "dynamodb:GetItem",
	"putItem":    "dynamodb:PutItem",
	"deleteItem": "dynamodb:DeleteItem",
	"updateItem": "dynamodb:UpdateItem",
}

// authoriseTaskCredentialsInvocation runs the per-task authorisation for
// one attempt's dispatch. Credentials "specifies a target role the state
// machine's execution role must assume before invoking the specified
// Resource", and a task without the privileges fails with
// States.Permissions ("A Task state failed because it had insufficient
// privileges to run the specified code") before any integration call.
// Activity resources carry no Credentials contract ("This field is
// supported by the Task types that use Lambda functions and a supported
// service") and no SFN-side invocation exists to authorise, so their
// target list is empty and the authorisation passes through.
func (e *Executor) authoriseTaskCredentialsInvocation(ctx context.Context, resource, input, roleArn string, resolved *resolvedInvocation) error {
	if e.taskAuthz == nil {
		return nil
	}
	for _, target := range e.integrationAuthorisationTargets(ctx, resource, input, resolved) {
		if err := e.taskAuthz.AuthoriseTaskCredentials(ctx, roleArn, target.action, target.resourceArn); err != nil {
			return &taskFailure{code: "States.Permissions", cause: err.Error(), failedToStart: true}
		}
	}
	return nil
}

// integrationAuthorisationTargets derives the action/resource pairs the
// integration invocation is authorised against, mirroring the dispatch
// switch's resource families. A family whose target the parameters cannot
// name yields no target: the invocation itself produces the failure for
// the missing parameter, so the authorisation has nothing to add. What the
// derivation resolves is recorded on the attempt's carrier so the dispatch
// reuses it.
func (e *Executor) integrationAuthorisationTargets(ctx context.Context, resource, input string, resolved *resolvedInvocation) []integrationTarget {
	var params map[string]interface{}
	if input != "" {
		_ = json.Unmarshal([]byte(input), &params)
	}
	switch {
	case arnutil.IsLambdaARN(resource):
		return []integrationTarget{{action: "lambda:InvokeFunction", resourceArn: resource}}
	case resource == "arn:aws:states:::lambda:invoke" || resource == "arn:aws:states:::aws-sdk:lambda:invoke":
		functionName := getStr(params, "FunctionName")
		if functionName == "" {
			return nil
		}
		if !strings.HasPrefix(functionName, "arn:") {
			functionName = arnutil.NewARNBuilder(e.accountID, e.region).Lambda().Function(functionName)
		}
		return []integrationTarget{{action: "lambda:InvokeFunction", resourceArn: functionName}}
	case strings.HasPrefix(resource, "arn:aws:states:::sqs:") || strings.HasPrefix(resource, "arn:aws:states:::aws-sdk:sqs:"):
		return e.sqsSendMessageTarget(ctx, params, resolved)
	case strings.HasPrefix(resource, "arn:aws:states:::sns:") || strings.HasPrefix(resource, "arn:aws:states:::aws-sdk:sns:"):
		topicArn, err := e.resolveSnsTopicArn(params)
		if err != nil {
			return nil
		}
		return []integrationTarget{{action: "sns:Publish", resourceArn: topicArn}}
	case strings.HasPrefix(resource, "arn:aws:states:::events:") || strings.HasPrefix(resource, "arn:aws:states:::aws-sdk:eventbridge:"):
		return e.putEventsTargets(params, resource)
	case strings.HasPrefix(resource, "arn:aws:states:::dynamodb:") || strings.HasPrefix(resource, "arn:aws:states:::aws-sdk:dynamodb:"):
		action, ok := dynamodbIntegrationActions[integrationAction(resource)]
		if !ok {
			return nil
		}
		tableName := getStr(params, "TableName")
		if tableName == "" {
			return nil
		}
		return []integrationTarget{{action: action, resourceArn: arnutil.NewARNBuilder(e.accountID, e.region).DynamoDB().Table(tableName)}}
	case e.isStartExecutionResource(resource):
		smArn := getStr(params, "StateMachineArn")
		if smArn == "" {
			return nil
		}
		return []integrationTarget{{action: "states:StartExecution", resourceArn: smArn}}
	}
	return nil
}

// sqsSendMessageTarget resolves the queue the SendMessage authorises
// against, through the send path's own queue resolution, then the queue's
// ARN. The resolved queue rides the attempt's carrier for the dispatch to
// reuse.
func (e *Executor) sqsSendMessageTarget(ctx context.Context, params map[string]interface{}, resolved *resolvedInvocation) []integrationTarget {
	if e.bus == nil || e.bus.SQSInvoker() == nil {
		return nil
	}
	queueURL, region, err := e.resolveSqsQueueURL(ctx, params)
	if err != nil {
		return nil
	}
	if resolved != nil {
		resolved.sqsQueueURL, resolved.sqsRegion = queueURL, region
	}
	queueARN, err := e.bus.SQSInvoker().GetQueueARN(ctx, region, queueURL)
	if err != nil || queueARN == "" {
		return nil
	}
	return []integrationTarget{{action: "sqs:SendMessage", resourceArn: queueARN}}
}

// putEventsTargets derives one target per distinct event bus the entries
// address — every bus the invocation writes to must be allowed.
func (e *Executor) putEventsTargets(params map[string]interface{}, resource string) []integrationTarget {
	eventsRegion := e.regionOfARN(resource)
	var targets []integrationTarget
	for _, busName := range eventBusNames(params) {
		targets = append(targets, integrationTarget{
			action:      "events:PutEvents",
			resourceArn: arnutil.NewARNBuilder(e.accountID, eventsRegion).Events().EventBus(busName),
		})
	}
	return targets
}
