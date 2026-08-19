package testutil

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// dynamoDBUtilityTests pins the DescribeLimits and DescribeEndpoints
// response shapes: the four modelled limit members carry the documented
// default quotas, and the reported endpoint is the client-facing host, not
// an unreachable AWS hostname.
func (r *TestRunner) dynamoDBUtilityTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "DescribeLimits_ReturnsDefaultQuotas", func() error {
		out, err := client.DescribeLimits(ctx, &dynamodb.DescribeLimitsInput{})
		if err != nil {
			return err
		}
		const accountMax = 80000
		const tableMax = 40000
		if out.AccountMaxReadCapacityUnits == nil || *out.AccountMaxReadCapacityUnits != accountMax {
			return fmt.Errorf("AccountMaxReadCapacityUnits = %v, want %d", out.AccountMaxReadCapacityUnits, accountMax)
		}
		if out.AccountMaxWriteCapacityUnits == nil || *out.AccountMaxWriteCapacityUnits != accountMax {
			return fmt.Errorf("AccountMaxWriteCapacityUnits = %v, want %d", out.AccountMaxWriteCapacityUnits, accountMax)
		}
		if out.TableMaxReadCapacityUnits == nil || *out.TableMaxReadCapacityUnits != tableMax {
			return fmt.Errorf("TableMaxReadCapacityUnits = %v, want %d", out.TableMaxReadCapacityUnits, tableMax)
		}
		if out.TableMaxWriteCapacityUnits == nil || *out.TableMaxWriteCapacityUnits != tableMax {
			return fmt.Errorf("TableMaxWriteCapacityUnits = %v, want %d", out.TableMaxWriteCapacityUnits, tableMax)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeEndpoints_ReportsReachableAddress", func() error {
		out, err := client.DescribeEndpoints(ctx, &dynamodb.DescribeEndpointsInput{})
		if err != nil {
			return err
		}
		if len(out.Endpoints) == 0 {
			return errors.New("no endpoints returned")
		}
		ep := out.Endpoints[0]
		if ep.Address == nil || *ep.Address == "" {
			return errors.New("endpoint address is empty")
		}
		if strings.HasSuffix(*ep.Address, ".amazonaws.com") {
			return fmt.Errorf("endpoint address %q is not reachable on this deployment", *ep.Address)
		}
		if ep.CachePeriodInMinutes <= 0 {
			return fmt.Errorf("CachePeriodInMinutes = %d, want positive", ep.CachePeriodInMinutes)
		}
		return nil
	}))

	return results
}
