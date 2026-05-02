package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

func (r *TestRunner) runNeptuneClusterEndpointTests(tc *neptuneContext) []TestResult {
	var results []TestResult
	endpointID := fmt.Sprintf("test-ep-%d", tc.ts)

	results = append(results, r.RunTest("neptune", "CreateDBClusterEndpoint", func() error {
		resp, err := tc.client.CreateDBClusterEndpoint(tc.ctx, &neptune.CreateDBClusterEndpointInput{
			DBClusterEndpointIdentifier: aws.String(endpointID),
			DBClusterIdentifier:         aws.String(tc.clusterID),
			EndpointType:                aws.String("READER"),
		})
		if err != nil {
			return err
		}
		if resp.DBClusterEndpointIdentifier == nil || *resp.DBClusterEndpointIdentifier != endpointID {
			return fmt.Errorf("expected endpoint ID=%s, got %v", endpointID, resp.DBClusterEndpointIdentifier)
		}
		if resp.DBClusterIdentifier == nil || *resp.DBClusterIdentifier != tc.clusterID {
			return fmt.Errorf("expected DBClusterIdentifier=%s, got %v", tc.clusterID, resp.DBClusterIdentifier)
		}
		if resp.EndpointType == nil || *resp.EndpointType != "READER" {
			return fmt.Errorf("expected EndpointType=READER, got %v", resp.EndpointType)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterEndpoints", func() error {
		resp, err := tc.client.DescribeDBClusterEndpoints(tc.ctx, &neptune.DescribeDBClusterEndpointsInput{
			DBClusterIdentifier: aws.String(tc.clusterID),
		})
		if err != nil {
			return err
		}
		found := false
		for _, ep := range resp.DBClusterEndpoints {
			if ep.DBClusterEndpointIdentifier != nil && *ep.DBClusterEndpointIdentifier == endpointID {
				found = true
				if ep.EndpointType == nil || *ep.EndpointType != "READER" {
					return fmt.Errorf("expected EndpointType READER, got %v", ep.EndpointType)
				}
				if ep.Status == nil || *ep.Status == "" {
					return fmt.Errorf("expected non-empty Status")
				}
				if ep.DBClusterIdentifier == nil || *ep.DBClusterIdentifier != tc.clusterID {
					return fmt.Errorf("expected DBClusterIdentifier=%s, got %v", tc.clusterID, ep.DBClusterIdentifier)
				}
			}
		}
		if !found {
			return fmt.Errorf("expected endpoint %s in DescribeDBClusterEndpoints response", endpointID)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBClusterEndpoint", func() error {
		_, err := tc.client.ModifyDBClusterEndpoint(tc.ctx, &neptune.ModifyDBClusterEndpointInput{
			DBClusterEndpointIdentifier: aws.String(endpointID),
			EndpointType:                aws.String("ANY"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBClusterEndpoint_Verify", func() error {
		resp, err := tc.client.DescribeDBClusterEndpoints(tc.ctx, &neptune.DescribeDBClusterEndpointsInput{
			DBClusterEndpointIdentifier: aws.String(endpointID),
		})
		if err != nil {
			return err
		}
		if len(resp.DBClusterEndpoints) != 1 {
			return fmt.Errorf("expected 1 endpoint, got %d", len(resp.DBClusterEndpoints))
		}
		ep := resp.DBClusterEndpoints[0]
		if ep.EndpointType == nil || *ep.EndpointType != "ANY" {
			return fmt.Errorf("expected EndpointType=ANY after modify, got %v", ep.EndpointType)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBClusterEndpoint", func() error {
		resp, err := tc.client.DeleteDBClusterEndpoint(tc.ctx, &neptune.DeleteDBClusterEndpointInput{
			DBClusterEndpointIdentifier: aws.String(endpointID),
		})
		if err != nil {
			return err
		}
		if resp.DBClusterEndpointIdentifier == nil || *resp.DBClusterEndpointIdentifier != endpointID {
			return fmt.Errorf("expected DBClusterEndpointIdentifier=%s in delete response, got %v", endpointID, resp.DBClusterEndpointIdentifier)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBClusterEndpoint_VerifyDeleted", func() error {
		resp, err := tc.client.DescribeDBClusterEndpoints(tc.ctx, &neptune.DescribeDBClusterEndpointsInput{
			DBClusterIdentifier: aws.String(tc.clusterID),
		})
		if err != nil {
			return err
		}
		for _, ep := range resp.DBClusterEndpoints {
			if ep.DBClusterEndpointIdentifier != nil && *ep.DBClusterEndpointIdentifier == endpointID {
				return fmt.Errorf("expected endpoint %s to be deleted but still found", endpointID)
			}
		}
		return nil
	}))

	endpointID2 := fmt.Sprintf("test-ep2-%d", tc.ts)
	results = append(results, r.RunTest("neptune", "DescribeDBClusterEndpoints_FilterByEndpointID", func() error {
		resp, err := tc.client.CreateDBClusterEndpoint(tc.ctx, &neptune.CreateDBClusterEndpointInput{
			DBClusterEndpointIdentifier: aws.String(endpointID2),
			DBClusterIdentifier:         aws.String(tc.clusterID),
			EndpointType:                aws.String("ANY"),
		})
		if err != nil {
			return err
		}
		if resp.DBClusterEndpointIdentifier == nil || *resp.DBClusterEndpointIdentifier != endpointID2 {
			return fmt.Errorf("endpoint ID mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterEndpoints_ByEndpointID", func() error {
		resp, err := tc.client.DescribeDBClusterEndpoints(tc.ctx, &neptune.DescribeDBClusterEndpointsInput{
			DBClusterEndpointIdentifier: aws.String(endpointID2),
		})
		if err != nil {
			return err
		}
		if len(resp.DBClusterEndpoints) != 1 {
			return fmt.Errorf("expected 1 endpoint when filtering by ID, got %d", len(resp.DBClusterEndpoints))
		}
		if resp.DBClusterEndpoints[0].DBClusterEndpointIdentifier == nil || *resp.DBClusterEndpoints[0].DBClusterEndpointIdentifier != endpointID2 {
			return fmt.Errorf("endpoint ID mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBClusterEndpoint2", func() error {
		resp, err := tc.client.DeleteDBClusterEndpoint(tc.ctx, &neptune.DeleteDBClusterEndpointInput{
			DBClusterEndpointIdentifier: aws.String(endpointID2),
		})
		if err != nil {
			return err
		}
		if resp.DBClusterEndpointIdentifier == nil || *resp.DBClusterEndpointIdentifier != endpointID2 {
			return fmt.Errorf("expected DBClusterEndpointIdentifier=%s in delete response, got %v", endpointID2, resp.DBClusterEndpointIdentifier)
		}
		return nil
	}))

	return results
}
