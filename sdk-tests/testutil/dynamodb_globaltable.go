package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"vorpalstacks-sdk-tests/config"
)

// createGlobalTableTestTable creates an empty table streaming both item
// images — the shape a table needs to join a global table.
func createGlobalTableTestTable(ctx context.Context, client *dynamodb.Client, name string, withSortKey bool) error {
	attrs := []types.AttributeDefinition{
		{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
	}
	schema := []types.KeySchemaElement{
		{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
	}
	if withSortKey {
		attrs = append(attrs, types.AttributeDefinition{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS})
		schema = append(schema, types.KeySchemaElement{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange})
	}
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName:            aws.String(name),
		AttributeDefinitions: attrs,
		KeySchema:            schema,
		BillingMode:          types.BillingModePayPerRequest,
		StreamSpecification: &types.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: types.StreamViewTypeNewAndOldImages,
		},
	})
	return err
}

// replicaRegionClient builds a client for a second region against the
// same endpoint — the second member region the global-table tests form
// their two-or-more replication groups over.
func (r *TestRunner) replicaRegionClient() (*dynamodb.Client, string, error) {
	replicaRegion := "us-west-2"
	if r.region == replicaRegion {
		replicaRegion = "eu-west-1"
	}
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   replicaRegion,
	})
	if err != nil {
		return nil, "", err
	}
	return dynamodb.NewFromConfig(cfg), replicaRegion, nil
}

// dynamoDBGlobalTableValidationTests pins the documented CreateGlobalTable
// preconditions: each replica region must hold an empty table with the same
// name and key/index schemas streaming both item images, regions may appear
// once, and a duplicate global table name is rejected.
func (r *TestRunner) dynamoDBGlobalTableValidationTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult
	suffix := time.Now().UnixNano()

	replicaClient, replicaRegion, cfgErr := r.replicaRegionClient()
	if cfgErr != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "GlobalTableValidation_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("load replica config: %v", cfgErr),
		})
	}

	results = append(results, r.RunTest("dynamodb", "CreateGlobalTable_RequiresBackingTable", func() error {
		name := fmt.Sprintf("GtVal-NoTable-%d", suffix)
		_, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String(r.region)},
				{RegionName: aws.String(replicaRegion)},
			},
		})
		if err == nil {
			return errors.New("expected an error without a backing table")
		}
		var nf *types.TableNotFoundException
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected TableNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "CreateGlobalTable_RequiresStreams", func() error {
		name := fmt.Sprintf("GtVal-NoStreams-%d", suffix)
		for _, c := range []*dynamodb.Client{client, replicaClient} {
			if _, err := c.CreateTable(ctx, &dynamodb.CreateTableInput{
				TableName: aws.String(name),
				AttributeDefinitions: []types.AttributeDefinition{
					{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
				},
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
				},
				BillingMode: types.BillingModePayPerRequest,
			}); err != nil {
				return err
			}
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})
		defer replicaClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})

		_, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String(r.region)},
				{RegionName: aws.String(replicaRegion)},
			},
		})
		if err == nil {
			return errors.New("expected an error without streaming enabled")
		}
		return expectAWSErrorCode(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "CreateGlobalTable_RequiresEmptyTable", func() error {
		name := fmt.Sprintf("GtVal-NonEmpty-%d", suffix)
		if err := createGlobalTableTestTable(ctx, client, name, false); err != nil {
			return err
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})
		if err := createGlobalTableTestTable(ctx, replicaClient, name, false); err != nil {
			return err
		}
		defer replicaClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})

		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(name),
			Item: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "seed"},
			},
		}); err != nil {
			return err
		}
		_, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String(r.region)},
				{RegionName: aws.String(replicaRegion)},
			},
		})
		if err == nil {
			return errors.New("expected an error for a non-empty replica table")
		}
		return expectAWSErrorCode(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "CreateGlobalTable_DuplicateRegion_Rejected", func() error {
		name := fmt.Sprintf("GtVal-DupRegion-%d", suffix)
		if err := createGlobalTableTestTable(ctx, client, name, false); err != nil {
			return err
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})

		_, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String(r.region)},
				{RegionName: aws.String(r.region)},
			},
		})
		if err == nil {
			return errors.New("expected an error for a duplicated replica region")
		}
		return expectAWSErrorCode(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "CreateGlobalTable_SchemaMustMatch", func() error {
		name := fmt.Sprintf("GtVal-Schema-%d", suffix)
		if err := createGlobalTableTestTable(ctx, client, name, false); err != nil {
			return err
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})
		if err := createGlobalTableTestTable(ctx, replicaClient, name, true); err != nil {
			return err
		}
		defer replicaClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})

		_, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String(r.region)},
				{RegionName: aws.String(replicaRegion)},
			},
		})
		if err == nil {
			return errors.New("expected an error for mismatched key schemas")
		}
		return expectAWSErrorCode(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "CreateGlobalTable_DuplicateName_Rejected", func() error {
		name := fmt.Sprintf("GtVal-DupName-%d", suffix)
		if err := createGlobalTableTestTable(ctx, client, name, false); err != nil {
			return err
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})
		if err := createGlobalTableTestTable(ctx, replicaClient, name, false); err != nil {
			return err
		}
		defer replicaClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})

		if _, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String(r.region)},
				{RegionName: aws.String(replicaRegion)},
			},
		}); err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		_, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String(r.region)},
				{RegionName: aws.String(replicaRegion)},
			},
		})
		if err == nil {
			return errors.New("expected GlobalTableAlreadyExistsException")
		}
		var ge *types.GlobalTableAlreadyExistsException
		if !errors.As(err, &ge) {
			return fmt.Errorf("expected GlobalTableAlreadyExistsException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeGlobalTable_ReturnsReplicationGroup", func() error {
		name := fmt.Sprintf("GtVal-Describe-%d", suffix)
		if err := createGlobalTableTestTable(ctx, client, name, false); err != nil {
			return err
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})
		if err := createGlobalTableTestTable(ctx, replicaClient, name, false); err != nil {
			return err
		}
		defer replicaClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})

		if _, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String(r.region)},
				{RegionName: aws.String(replicaRegion)},
			},
		}); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		out, err := client.DescribeGlobalTable(ctx, &dynamodb.DescribeGlobalTableInput{
			GlobalTableName: aws.String(name),
		})
		if err != nil {
			return err
		}
		desc := out.GlobalTableDescription
		if aws.ToString(desc.GlobalTableName) != name {
			return fmt.Errorf("GlobalTableName = %s, want %s", aws.ToString(desc.GlobalTableName), name)
		}
		if len(desc.ReplicationGroup) != 2 {
			return fmt.Errorf("ReplicationGroup = %+v, want 2 entries", desc.ReplicationGroup)
		}
		for _, replica := range desc.ReplicationGroup {
			region := aws.ToString(replica.RegionName)
			if region != r.region && region != replicaRegion {
				return fmt.Errorf("replica region = %s, want %s or %s", region, r.region, replicaRegion)
			}
			if replica.ReplicaStatus != types.ReplicaStatusActive {
				return fmt.Errorf("replica %s status = %s, want ACTIVE", region, replica.ReplicaStatus)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateGlobalTable_AddAndRemoveReplica", func() error {
		name := fmt.Sprintf("GtVal-Update-%d", suffix)
		if err := createGlobalTableTestTable(ctx, client, name, false); err != nil {
			return err
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})
		if err := createGlobalTableTestTable(ctx, replicaClient, name, false); err != nil {
			return err
		}
		defer replicaClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})

		if _, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String(r.region)},
				{RegionName: aws.String(replicaRegion)},
			},
		}); err != nil {
			return fmt.Errorf("create: %v", err)
		}

		// Replicas leave individually: removing one member leaves the
		// surviving member's record intact.
		oneLeft, err := client.UpdateGlobalTable(ctx, &dynamodb.UpdateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicaUpdates: []types.ReplicaUpdate{{
				Delete: &types.DeleteReplicaAction{RegionName: aws.String(r.region)},
			}},
		})
		if err != nil {
			return fmt.Errorf("remove first replica: %v", err)
		}
		if len(oneLeft.GlobalTableDescription.ReplicationGroup) != 1 ||
			aws.ToString(oneLeft.GlobalTableDescription.ReplicationGroup[0].RegionName) != replicaRegion {
			return fmt.Errorf("replication group after the first removal = %+v, want %s alone",
				oneLeft.GlobalTableDescription.ReplicationGroup, replicaRegion)
		}

		// The final member's departure empties the group and takes the
		// record with it, deleting the last member's table — the terminal
		// state teardown tooling polls DescribeGlobalTable for, and
		// re-creation goes through CreateGlobalTable.
		removed, err := client.UpdateGlobalTable(ctx, &dynamodb.UpdateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicaUpdates: []types.ReplicaUpdate{{
				Delete: &types.DeleteReplicaAction{RegionName: aws.String(replicaRegion)},
			}},
		})
		if err != nil {
			return fmt.Errorf("remove final replica: %v", err)
		}
		if len(removed.GlobalTableDescription.ReplicationGroup) != 0 {
			return fmt.Errorf("replication group after removal = %+v, want empty", removed.GlobalTableDescription.ReplicationGroup)
		}
		_, err = client.UpdateGlobalTable(ctx, &dynamodb.UpdateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicaUpdates: []types.ReplicaUpdate{{
				Create: &types.CreateReplicaAction{RegionName: aws.String(replicaRegion)},
			}},
		})
		var gnf *types.GlobalTableNotFoundException
		if !errors.As(err, &gnf) {
			return fmt.Errorf("re-add after emptying: expected GlobalTableNotFoundException, got: %T: %v", err, err)
		}
		if _, err := client.DescribeGlobalTable(ctx, &dynamodb.DescribeGlobalTableInput{
			GlobalTableName: aws.String(name),
		}); !errors.As(err, &gnf) {
			return fmt.Errorf("describe after emptying: expected GlobalTableNotFoundException, got: %T: %v", err, err)
		}

		// Re-creation over both regions again: the emptied group's
		// teardown deleted both member tables, so the backing tables are
		// recreated first.
		if err := createGlobalTableTestTable(ctx, client, name, false); err != nil {
			return fmt.Errorf("re-create backing table: %v", err)
		}
		if err := createGlobalTableTestTable(ctx, replicaClient, name, false); err != nil {
			return fmt.Errorf("re-create replica backing table: %v", err)
		}
		added, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String(r.region)},
				{RegionName: aws.String(replicaRegion)},
			},
		})
		if err != nil {
			return fmt.Errorf("re-create global table: %v", err)
		}
		if len(added.GlobalTableDescription.ReplicationGroup) != 2 {
			return fmt.Errorf("replication group after re-create = %+v, want 2 entries", added.GlobalTableDescription.ReplicationGroup)
		}
		for _, replica := range added.GlobalTableDescription.ReplicationGroup {
			region := aws.ToString(replica.RegionName)
			if region != r.region && region != replicaRegion {
				return fmt.Errorf("re-created replica region = %s, want %s or %s", region, r.region, replicaRegion)
			}
		}

		// Adding a replica whose region holds no table must fail.
		_, err = client.UpdateGlobalTable(ctx, &dynamodb.UpdateGlobalTableInput{
			GlobalTableName: aws.String(name),
			ReplicaUpdates: []types.ReplicaUpdate{{
				Create: &types.CreateReplicaAction{RegionName: aws.String("ap-south-1")},
			}},
		})
		if err == nil {
			return errors.New("expected an error adding a replica without a backing table")
		}
		var nf *types.TableNotFoundException
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected TableNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	return results
}
