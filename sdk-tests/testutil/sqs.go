package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"vorpalstacks-sdk-tests/config"
)

// createTestQueue creates a throwaway queue for a single test case and
// returns its URL together with a cleanup function that deletes it. The URL
// is taken from the CreateQueue response itself; no extra GetQueueUrl
// round-trip is needed.
func createTestQueue(ctx context.Context, client *sqs.Client, name string, attributes map[string]string) (*string, func(), error) {
	resp, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName:  aws.String(name),
		Attributes: attributes,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create %s: %v", name, err)
	}
	if resp.QueueUrl == nil {
		return nil, nil, fmt.Errorf("CreateQueue returned nil QueueUrl for %s", name)
	}
	cleanup := func() {
		client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: resp.QueueUrl})
	}
	return resp.QueueUrl, cleanup, nil
}

// queueArn fetches the QueueArn attribute of the given queue.
func queueArn(ctx context.Context, client *sqs.Client, queueURL *string) (string, error) {
	resp, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       queueURL,
		AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
	})
	if err != nil {
		return "", fmt.Errorf("get arn: %v", err)
	}
	return resp.Attributes[string(types.QueueAttributeNameQueueArn)], nil
}

func (r *TestRunner) RunSQSTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "sqs",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := sqs.NewFromConfig(cfg)
	ctx := context.Background()

	queueName := fmt.Sprintf("TestQueue-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("sqs", "CreateQueue", func() error {
		resp, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		if resp.QueueUrl == nil {
			return fmt.Errorf("CreateQueue returned nil QueueUrl")
		}
		if !strings.Contains(*resp.QueueUrl, queueName) {
			return fmt.Errorf("QueueUrl does not contain queue name %q: %s", queueName, *resp.QueueUrl)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "GetQueueUrl", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		if resp.QueueUrl == nil {
			return fmt.Errorf("GetQueueUrl returned nil QueueUrl")
		}
		if !strings.Contains(*resp.QueueUrl, queueName) {
			return fmt.Errorf("QueueUrl does not contain queue name %q: %s", queueName, *resp.QueueUrl)
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "GetQueueAttributes", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		attrResp, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return err
		}
		if attrResp.Attributes == nil {
			return fmt.Errorf("GetQueueAttributes returned nil Attributes")
		}
		if attrResp.Attributes[string(types.QueueAttributeNameVisibilityTimeout)] == "" {
			return fmt.Errorf("VisibilityTimeout is empty")
		}
		if _, ok := attrResp.Attributes[string(types.QueueAttributeNameQueueArn)]; !ok {
			return fmt.Errorf("GetQueueAttributes missing QueueArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("sqs", "GetQueueAttributes_SpecificAttributes", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		attrResp, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl: resp.QueueUrl,
			AttributeNames: []types.QueueAttributeName{
				types.QueueAttributeNameQueueArn,
				types.QueueAttributeNameVisibilityTimeout,
			},
		})
		if err != nil {
			return err
		}
		if attrResp.Attributes == nil {
			return fmt.Errorf("GetQueueAttributes returned nil Attributes")
		}
		if _, ok := attrResp.Attributes[string(types.QueueAttributeNameQueueArn)]; !ok {
			return fmt.Errorf("GetQueueAttributes missing QueueArn")
		}
		if !strings.HasPrefix(attrResp.Attributes[string(types.QueueAttributeNameQueueArn)], "arn:") {
			return fmt.Errorf("QueueArn does not start with 'arn:': %s", attrResp.Attributes[string(types.QueueAttributeNameQueueArn)])
		}
		if _, ok := attrResp.Attributes[string(types.QueueAttributeNameVisibilityTimeout)]; !ok {
			return fmt.Errorf("GetQueueAttributes missing VisibilityTimeout")
		}
		return nil
	}))

	results = append(results, r.runSQSQueueTests(ctx, client, queueName)...)
	results = append(results, r.runSQSMessageTests(ctx, client, queueName)...)
	results = append(results, r.runSQSPermissionTests(ctx, client, queueName)...)
	results = append(results, r.runSQSDLQTests(ctx, client)...)
	results = append(results, r.runSQSEdgeTests(ctx, client, queueName)...)

	results = append(results, r.RunTest("sqs", "DeleteQueue", func() error {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		_, err = client.DeleteQueue(ctx, &sqs.DeleteQueueInput{
			QueueUrl: resp.QueueUrl,
		})
		if err != nil {
			return err
		}
		_, err = client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
		if err == nil {
			return fmt.Errorf("GetQueueUrl should fail after DeleteQueue")
		}
		return AssertErrorContains(err, "QueueDoesNotExist")
	}))

	return results
}
