package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	appsyncTypes "github.com/aws/aws-sdk-go-v2/service/appsync/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgeTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

// runEventBridgeToAppSync verifies the AppSync target delivery path end to
// end. Success is observed through the dead-letter channel rather than the
// mutation's return value (which EventBridge discards): a rule pointing at a
// healthy AWS_IAM GraphQL API with a schema must leave its dead-letter queue
// empty, while a control rule pointing at a non-existent API on the same bus
// must dead-letter its event — the control proves the observation channel
// itself works, so silence on the valid target means the mutation executed.
func (r *TestRunner) runEventBridgeToAppSync(ic *integClients, ts string) TestResult {
	const testName = "EventBridge_AppSync"

	apiResp, err := ic.appsync.CreateGraphqlApi(ic.ctx, &appsync.CreateGraphqlApiInput{
		Name:               aws.String(fmt.Sprintf("integ-eb-gql-%s", ts)),
		AuthenticationType: appsyncTypes.AuthenticationTypeAwsIam,
	})
	if err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create graphql api: %w", err) })
	}
	apiID := aws.ToString(apiResp.GraphqlApi.ApiId)
	defer ic.appsync.DeleteGraphqlApi(ic.ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(apiID)})

	schemaResp, err := ic.appsync.StartSchemaCreation(ic.ctx, &appsync.StartSchemaCreationInput{
		ApiId:      aws.String(apiID),
		Definition: []byte("type Mutation { pushEvent: String }"),
	})
	if err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("apply schema: %w", err) })
	}
	if schemaResp.Status != "PROCESSING" && schemaResp.Status != "SUCCESS" {
		return r.RunTest(integSvc, testName, func() error {
			return fmt.Errorf("schema creation status %s", schemaResp.Status)
		})
	}

	okQueueURL, err := ic.createQueue(fmt.Sprintf("integ-eb-appsync-ok-dlq-%s", ts))
	if err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create ok dlq: %w", err) })
	}
	defer ic.deleteQueue(okQueueURL)
	badQueueURL, err := ic.createQueue(fmt.Sprintf("integ-eb-appsync-bad-dlq-%s", ts))
	if err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create control dlq: %w", err) })
	}
	defer ic.deleteQueue(badQueueURL)
	okDLQArn := fmt.Sprintf("arn:aws:sqs:%s:%s:integ-eb-appsync-ok-dlq-%s", ic.region, ic.accountID, ts)
	badDLQArn := fmt.Sprintf("arn:aws:sqs:%s:%s:integ-eb-appsync-bad-dlq-%s", ic.region, ic.accountID, ts)

	busName := fmt.Sprintf("integ-eb-appsync-bus-%s", ts)
	if _, err := ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)}); err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create bus: %w", err) })
	}
	defer ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})

	operation := "mutation { pushEvent }"
	pattern := aws.String(`{"source":["com.integration.appsync"]}`)
	okRule := fmt.Sprintf("integ-eb-appsync-ok-%s", ts)
	badRule := fmt.Sprintf("integ-eb-appsync-bad-%s", ts)
	for _, rule := range []struct{ name, apiID, dlqArn string }{
		{okRule, apiID, okDLQArn},
		{badRule, "nonexistent", badDLQArn},
	} {
		if _, err := ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{
			Name:         aws.String(rule.name),
			EventBusName: aws.String(busName),
			EventPattern: pattern,
		}); err != nil {
			return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("put rule %s: %w", rule.name, err) })
		}
		defer ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{
			Name: aws.String(rule.name), EventBusName: aws.String(busName), Force: true,
		})
		if _, err := ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(rule.name),
			EventBusName: aws.String(busName),
			Targets: []eventbridgeTypes.Target{{
				Id:                aws.String("t1"),
				Arn:               aws.String(fmt.Sprintf("arn:aws:appsync:%s:%s:apis/%s/endpoints/GRAPHQL", ic.region, ic.accountID, rule.apiID)),
				AppSyncParameters: &eventbridgeTypes.AppSyncParameters{GraphQLOperation: aws.String(operation)},
				DeadLetterConfig:  &eventbridgeTypes.DeadLetterConfig{Arn: aws.String(rule.dlqArn)},
			}},
		}); err != nil {
			return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("put targets %s: %w", rule.name, err) })
		}
		defer ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(rule.name), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}

	if _, err := ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []eventbridgeTypes.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.appsync"),
			DetailType:   aws.String("AppSyncDeliveryTest"),
			Detail:       aws.String(`{"kind":"appsync-target"}`),
		}},
	}); err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("put events: %w", err) })
	}

	// The control target dead-letters after its retry budget; both targets
	// share the fire time and the retry schedule, so by the time the control
	// copy lands a failing valid delivery would have landed too.
	control := r.pollVerify(testName, defaultPollTimeout, func() error {
		return ic.verifyMessageContains(badQueueURL, "AppSyncDeliveryTest")
	})
	if control.Status != "PASS" {
		return control
	}

	return r.RunTest(integSvc, testName, func() error {
		msgs, err := ic.receiveMessages(okQueueURL, 5, 3)
		if err != nil {
			return fmt.Errorf("receive from the valid target's dlq: %w", err)
		}
		if len(msgs) != 0 {
			return fmt.Errorf("the valid AppSync target dead-lettered %d message(s); the mutation did not execute", len(msgs))
		}
		return nil
	})
}
