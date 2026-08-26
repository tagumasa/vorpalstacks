package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"vorpalstacks-sdk-tests/config"
)

type disableHostPrefixFinalize struct{}

func (*disableHostPrefixFinalize) ID() string { return "DisableHostPrefixFinalize" }

func (*disableHostPrefixFinalize) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
	ctx = smithyhttp.SetHostnameImmutable(ctx, true)
	return next.HandleFinalize(ctx, in)
}

func newNeptuneGraphClient(r *TestRunner) (*neptunegraph.Client, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, err
	}
	return neptunegraph.NewFromConfig(cfg, func(o *neptunegraph.Options) {
		o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
			return stack.Finalize.Add(&disableHostPrefixFinalize{}, middleware.Before)
		})
	}), nil
}

type neptunegraphContext struct {
	client     *neptunegraph.Client
	ctx        context.Context
	graphID    string
	graphARN   string
	snapshotID string
	tsNano     string
	region     string
	accountID  string
}

// unique returns a per-suite unique resource name.
func (tc *neptunegraphContext) unique(prefix string) string {
	return prefix + "-" + tc.tsNano[len(tc.tsNano)-8:]
}

// requireID fails a test whose prerequisite resource was never created.
func requireID(id, what string) error {
	if id == "" {
		return fmt.Errorf("no %s", what)
	}
	return nil
}

func (tc *neptunegraphContext) requireGraph() error {
	return requireID(tc.graphID, "graph ID")
}

func (tc *neptunegraphContext) requireGraphARN() error {
	return requireID(tc.graphARN, "graph ARN")
}

func (tc *neptunegraphContext) requireSnapshot() error {
	return requireID(tc.snapshotID, "snapshot ID")
}

// getGraph fetches the current graph summary for id.
func (tc *neptunegraphContext) getGraph(id string) (*neptunegraph.GetGraphOutput, error) {
	return tc.client.GetGraph(tc.ctx, &neptunegraph.GetGraphInput{
		GraphIdentifier: aws.String(id),
	})
}

func (r *TestRunner) RunNeptunegraphTests() []TestResult {
	var results []TestResult

	client, err := newNeptuneGraphClient(r)
	if err != nil {
		return append(results, SetupFailResult("neptunegraph", "Failed to create client: %v", err))
	}

	tc := &neptunegraphContext{
		client:    client,
		ctx:       context.Background(),
		region:    r.region,
		accountID: r.accountID,
	}
	tc.tsNano = fmt.Sprintf("%d", time.Now().UnixNano())

	results = append(results, r.runNeptunegraphGraphTests(tc)...)
	results = append(results, r.runNeptunegraphSnapshotTests(tc)...)
	results = append(results, r.runNeptunegraphEndpointTests(tc)...)
	results = append(results, r.runNeptunegraphImportTaskTests(tc)...)
	results = append(results, r.runNeptunegraphExportTaskTests(tc)...)
	results = append(results, r.runNeptunegraphTagTests(tc)...)
	results = append(results, r.runNeptunegraphQueryTests(tc)...)
	results = append(results, r.runNeptunegraphEdgeTests(tc)...)

	return results
}
