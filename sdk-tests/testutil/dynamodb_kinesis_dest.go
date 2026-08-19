package testutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	kinesistypes "github.com/aws/aws-sdk-go-v2/service/kinesis/types"

	"vorpalstacks-sdk-tests/config"
)

// kinesisDestinationPollDeadline bounds the wait for the asynchronous
// destination transitions (ENABLING/DISABLING/UPDATING settle after a
// service-side delay).
const kinesisDestinationPollDeadline = 6 * time.Second

// kinesisDestFixture bundles the resources the Kinesis destination tests
// need: two streams and a table in the default region.
type kinesisDestFixture struct {
	streamA   string
	streamArn string
	streamB   string
	tableName string
}

// setupKinesisDestFixture creates two Kinesis streams and a DynamoDB table.
func setupKinesisDestFixture(ctx context.Context, kc *kinesis.Client, dc *dynamodb.Client, suffix int64) (*kinesisDestFixture, error) {
	f := &kinesisDestFixture{
		streamA:   fmt.Sprintf("ddb-kds-a-%d", suffix),
		streamB:   fmt.Sprintf("ddb-kds-b-%d", suffix),
		tableName: fmt.Sprintf("ddb-kds-tbl-%d", suffix),
	}
	for _, sn := range []string{f.streamA, f.streamB} {
		if _, err := kc.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(sn),
			ShardCount: aws.Int32(1),
		}); err != nil {
			return nil, fmt.Errorf("create stream %s: %v", sn, err)
		}
	}
	deadline := time.Now().Add(5 * time.Second)
	for _, sn := range []string{f.streamA, f.streamB} {
		for {
			out, err := kc.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(sn)})
			if err == nil && out.StreamDescription.StreamStatus == kinesistypes.StreamStatusActive {
				if sn == f.streamA {
					f.streamArn = aws.ToString(out.StreamDescription.StreamARN)
				}
				break
			}
			if time.Now().After(deadline) {
				return nil, fmt.Errorf("stream %s did not become ACTIVE", sn)
			}
			time.Sleep(150 * time.Millisecond)
		}
	}
	if _, err := dc.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(f.tableName),
		AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamodbtypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
		},
		BillingMode: dynamodbtypes.BillingModePayPerRequest,
	}); err != nil {
		return nil, fmt.Errorf("create table: %v", err)
	}
	if err := waitKinesisDestTableActive(ctx, dc, f.tableName); err != nil {
		return nil, err
	}
	return f, nil
}

// waitKinesisDestTableActive polls DescribeTable until the table is ACTIVE.
func waitKinesisDestTableActive(ctx context.Context, dc *dynamodb.Client, tableName string) error {
	deadline := time.Now().Add(5 * time.Second)
	for {
		out, err := dc.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
		if err == nil && out.Table.TableStatus == dynamodbtypes.TableStatusActive {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("table %s did not become ACTIVE", tableName)
		}
		time.Sleep(150 * time.Millisecond)
	}
}

// teardown releases the fixture resources; a destination left behind is
// disabled first so the streams can be deleted.
func (f *kinesisDestFixture) teardown(ctx context.Context, kc *kinesis.Client, dc *dynamodb.Client) {
	if f.streamArn != "" {
		_, _ = dc.DisableKinesisStreamingDestination(ctx, &dynamodb.DisableKinesisStreamingDestinationInput{
			TableName: aws.String(f.tableName),
			StreamArn: aws.String(f.streamArn),
		})
	}
	time.Sleep(1500 * time.Millisecond)
	_, _ = dc.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(f.tableName)})
	for _, sn := range []string{f.streamA, f.streamB} {
		_, _ = kc.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(sn)})
	}
}

// describeDestination fetches the destinations list for the table.
func describeKinesisDestinations(ctx context.Context, dc *dynamodb.Client, tableName string) ([]dynamodbtypes.KinesisDataStreamDestination, error) {
	out, err := dc.DescribeKinesisStreamingDestination(ctx, &dynamodb.DescribeKinesisStreamingDestinationInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}
	return out.KinesisDataStreamDestinations, nil
}

// waitKinesisDestinationStatus polls until the destination for streamArn
// reaches the wanted status, or the deadline lapses.
func waitKinesisDestinationStatus(ctx context.Context, dc *dynamodb.Client, tableName, streamArn, want string) error {
	deadline := time.Now().Add(kinesisDestinationPollDeadline)
	for {
		dests, err := describeKinesisDestinations(ctx, dc, tableName)
		if err != nil {
			return err
		}
		for _, d := range dests {
			if aws.ToString(d.StreamArn) == streamArn && string(d.DestinationStatus) == want {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("destination did not reach %s", want)
		}
		time.Sleep(150 * time.Millisecond)
	}
}

// kinesisDestRecordPayload mirrors the JSON DynamoDB writes into a Kinesis
// destination stream for one item change.
type kinesisDestRecordPayload struct {
	Keys                        map[string]map[string]string `json:"Keys"`
	EventName                   string                       `json:"eventName"`
	ApproximateCreationDateTime float64                      `json:"ApproximateCreationDateTime"`
}

// readFirstKinesisDestRecord reads the first record from the stream's first
// shard and decodes the DynamoDB destination payload.
func readFirstKinesisDestRecord(ctx context.Context, kc *kinesis.Client, streamName string) (*kinesisDestRecordPayload, error) {
	desc, err := kc.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(streamName)})
	if err != nil {
		return nil, fmt.Errorf("describe stream: %v", err)
	}
	if len(desc.StreamDescription.Shards) == 0 {
		return nil, errors.New("stream has no shards")
	}
	shardID := aws.ToString(desc.StreamDescription.Shards[0].ShardId)
	iter, err := kc.GetShardIterator(ctx, &kinesis.GetShardIteratorInput{
		StreamName:        aws.String(streamName),
		ShardId:           aws.String(shardID),
		ShardIteratorType: kinesistypes.ShardIteratorTypeTrimHorizon,
	})
	if err != nil {
		return nil, fmt.Errorf("get shard iterator: %v", err)
	}
	recs, err := kc.GetRecords(ctx, &kinesis.GetRecordsInput{
		ShardIterator: iter.ShardIterator,
	})
	if err != nil {
		return nil, fmt.Errorf("get records: %v", err)
	}
	for _, r := range recs.Records {
		var payload kinesisDestRecordPayload
		if err := json.Unmarshal(r.Data, &payload); err != nil {
			return nil, fmt.Errorf("decode record: %v", err)
		}
		return &payload, nil
	}
	return nil, errors.New("no records in stream")
}

// pollFirstKinesisDestRecord waits for the asynchronous record emit.
func pollFirstKinesisDestRecord(ctx context.Context, kc *kinesis.Client, streamName string) (*kinesisDestRecordPayload, error) {
	deadline := time.Now().Add(6 * time.Second)
	for {
		payload, err := readFirstKinesisDestRecord(ctx, kc, streamName)
		if err == nil {
			return payload, nil
		}
		if time.Now().After(deadline) {
			return nil, err
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// dynamoDBKinesisDestinationTests pins the Kinesis streaming destination
// contract: one stream per table, destinations must exist in the table's
// region, the ENABLING/DISABLING/UPDATING intermediate states settle
// asynchronously, records carry millisecond timestamps by default, and a
// disable issued during the enable transition must win.
func (r *TestRunner) dynamoDBKinesisDestinationTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "KinesisDestination_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("load config: %v", err),
		})
	}
	kc := kinesis.NewFromConfig(cfg)

	results = append(results, r.RunTest("dynamodb", "EnableKinesisStreamingDestination_UnknownStream_Rejected", func() error {
		f, err := setupKinesisDestFixture(ctx, kc, client, time.Now().UnixNano())
		if err != nil {
			return err
		}
		defer f.teardown(ctx, kc, client)

		missing := fmt.Sprintf("arn:aws:kinesis:%s:111122223333:stream/never-created-%d", r.region, time.Now().UnixNano())
		_, err = client.EnableKinesisStreamingDestination(ctx, &dynamodb.EnableKinesisStreamingDestinationInput{
			TableName: aws.String(f.tableName),
			StreamArn: aws.String(missing),
		})
		if err == nil {
			return errors.New("expected an error for a nonexistent destination stream")
		}
		var nf *dynamodbtypes.ResourceNotFoundException
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		dests, err := describeKinesisDestinations(ctx, client, f.tableName)
		if err != nil {
			return err
		}
		if len(dests) != 0 {
			return fmt.Errorf("nonexistent stream must not be recorded, got %d destinations", len(dests))
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "KinesisDestination_LifecycleAndRecordDelivery", func() error {
		f, err := setupKinesisDestFixture(ctx, kc, client, time.Now().UnixNano())
		if err != nil {
			return err
		}
		defer f.teardown(ctx, kc, client)

		enable, err := client.EnableKinesisStreamingDestination(ctx, &dynamodb.EnableKinesisStreamingDestinationInput{
			TableName: aws.String(f.tableName),
			StreamArn: aws.String(f.streamArn),
			EnableKinesisStreamingConfiguration: &dynamodbtypes.EnableKinesisStreamingConfiguration{
				ApproximateCreationDateTimePrecision: dynamodbtypes.ApproximateCreationDateTimePrecisionMillisecond,
			},
		})
		if err != nil {
			return fmt.Errorf("enable: %v", err)
		}
		if enable.DestinationStatus != dynamodbtypes.DestinationStatusEnabling {
			return fmt.Errorf("enable returned %s, want ENABLING", enable.DestinationStatus)
		}
		if enable.EnableKinesisStreamingConfiguration == nil ||
			enable.EnableKinesisStreamingConfiguration.ApproximateCreationDateTimePrecision != dynamodbtypes.ApproximateCreationDateTimePrecisionMillisecond {
			return errors.New("enable must echo the configuration")
		}
		if err := waitKinesisDestinationStatus(ctx, client, f.tableName, f.streamArn, "ACTIVE"); err != nil {
			return err
		}

		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(f.tableName),
			Item: map[string]dynamodbtypes.AttributeValue{
				"pk": &dynamodbtypes.AttributeValueMemberS{Value: "item-1"},
			},
		}); err != nil {
			return fmt.Errorf("put item: %v", err)
		}
		payload, err := pollFirstKinesisDestRecord(ctx, kc, f.streamA)
		if err != nil {
			return fmt.Errorf("record delivery: %v", err)
		}
		if payload.EventName != "INSERT" {
			return fmt.Errorf("record eventName = %s, want INSERT", payload.EventName)
		}
		if payload.Keys["pk"]["S"] != "item-1" {
			return fmt.Errorf("record keys = %v, want pk item-1", payload.Keys)
		}
		// Millisecond precision: epoch-milliseconds (~1.7e12), not seconds.
		if payload.ApproximateCreationDateTime < 1e12 {
			return fmt.Errorf("record ApproximateCreationDateTime = %f, want millisecond precision", payload.ApproximateCreationDateTime)
		}

		update, err := client.UpdateKinesisStreamingDestination(ctx, &dynamodb.UpdateKinesisStreamingDestinationInput{
			TableName: aws.String(f.tableName),
			StreamArn: aws.String(f.streamArn),
			UpdateKinesisStreamingConfiguration: &dynamodbtypes.UpdateKinesisStreamingConfiguration{
				ApproximateCreationDateTimePrecision: dynamodbtypes.ApproximateCreationDateTimePrecisionMicrosecond,
			},
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if update.DestinationStatus != dynamodbtypes.DestinationStatusUpdating {
			return fmt.Errorf("update returned %s, want UPDATING", update.DestinationStatus)
		}
		if update.UpdateKinesisStreamingConfiguration == nil ||
			update.UpdateKinesisStreamingConfiguration.ApproximateCreationDateTimePrecision != dynamodbtypes.ApproximateCreationDateTimePrecisionMicrosecond {
			return errors.New("update must echo the configuration")
		}
		if err := waitKinesisDestinationStatus(ctx, client, f.tableName, f.streamArn, "ACTIVE"); err != nil {
			return err
		}
		dests, err := describeKinesisDestinations(ctx, client, f.tableName)
		if err != nil {
			return err
		}
		if len(dests) != 1 || dests[0].ApproximateCreationDateTimePrecision != dynamodbtypes.ApproximateCreationDateTimePrecisionMicrosecond {
			return fmt.Errorf("describe after update = %+v, want MICROSECOND precision", dests)
		}

		disable, err := client.DisableKinesisStreamingDestination(ctx, &dynamodb.DisableKinesisStreamingDestinationInput{
			TableName: aws.String(f.tableName),
			StreamArn: aws.String(f.streamArn),
		})
		if err != nil {
			return fmt.Errorf("disable: %v", err)
		}
		if disable.DestinationStatus != dynamodbtypes.DestinationStatusDisabling {
			return fmt.Errorf("disable returned %s, want DISABLING", disable.DestinationStatus)
		}
		deadline := time.Now().Add(kinesisDestinationPollDeadline)
		for {
			dests, err := describeKinesisDestinations(ctx, client, f.tableName)
			if err != nil {
				return err
			}
			if len(dests) == 0 {
				break
			}
			if time.Now().After(deadline) {
				return errors.New("destination was not removed after disable")
			}
			time.Sleep(150 * time.Millisecond)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "EnableKinesisStreamingDestination_SecondStream_Rejected", func() error {
		f, err := setupKinesisDestFixture(ctx, kc, client, time.Now().UnixNano())
		if err != nil {
			return err
		}
		defer f.teardown(ctx, kc, client)

		descB, err := kc.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(f.streamB)})
		if err != nil {
			return err
		}
		if _, err := client.EnableKinesisStreamingDestination(ctx, &dynamodb.EnableKinesisStreamingDestinationInput{
			TableName: aws.String(f.tableName),
			StreamArn: aws.String(f.streamArn),
		}); err != nil {
			return fmt.Errorf("first enable: %v", err)
		}
		if err := waitKinesisDestinationStatus(ctx, client, f.tableName, f.streamArn, "ACTIVE"); err != nil {
			return err
		}

		_, err = client.EnableKinesisStreamingDestination(ctx, &dynamodb.EnableKinesisStreamingDestinationInput{
			TableName: aws.String(f.tableName),
			StreamArn: descB.StreamDescription.StreamARN,
		})
		if err == nil {
			return errors.New("expected an error enabling a second stream")
		}
		var ri *dynamodbtypes.ResourceInUseException
		if !errors.As(err, &ri) {
			return fmt.Errorf("expected ResourceInUseException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DisableKinesisStreamingDestination_ImmediateAfterEnable_Wins", func() error {
		f, err := setupKinesisDestFixture(ctx, kc, client, time.Now().UnixNano())
		if err != nil {
			return err
		}
		defer f.teardown(ctx, kc, client)

		if _, err := client.EnableKinesisStreamingDestination(ctx, &dynamodb.EnableKinesisStreamingDestinationInput{
			TableName: aws.String(f.tableName),
			StreamArn: aws.String(f.streamArn),
		}); err != nil {
			return fmt.Errorf("enable: %v", err)
		}
		// Disable while the enable transition is still in flight. The
		// enable's delayed write-back must not resurrect the destination:
		// from the moment Disable returns, the observed status only moves
		// forward (DISABLING, then gone) and never back to ACTIVE.
		time.Sleep(500 * time.Millisecond)
		if _, err := client.DisableKinesisStreamingDestination(ctx, &dynamodb.DisableKinesisStreamingDestinationInput{
			TableName: aws.String(f.tableName),
			StreamArn: aws.String(f.streamArn),
		}); err != nil {
			return fmt.Errorf("disable during enable transition: %v", err)
		}
		deadline := time.Now().Add(2500 * time.Millisecond)
		for {
			dests, err := describeKinesisDestinations(ctx, client, f.tableName)
			if err != nil {
				return err
			}
			for _, d := range dests {
				if aws.ToString(d.StreamArn) != f.streamArn {
					continue
				}
				if d.DestinationStatus != dynamodbtypes.DestinationStatusDisabling {
					return fmt.Errorf("disable during enable transition was resurrected: status %s", d.DestinationStatus)
				}
			}
			if len(dests) == 0 {
				return nil
			}
			if time.Now().After(deadline) {
				return errors.New("destination was not removed after disable")
			}
			time.Sleep(100 * time.Millisecond)
		}
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeKinesisStreamingDestination_NoDestinations", func() error {
		f, err := setupKinesisDestFixture(ctx, kc, client, time.Now().UnixNano())
		if err != nil {
			return err
		}
		defer f.teardown(ctx, kc, client)

		out, err := client.DescribeKinesisStreamingDestination(ctx, &dynamodb.DescribeKinesisStreamingDestinationInput{
			TableName: aws.String(f.tableName),
		})
		if err != nil {
			return err
		}
		if aws.ToString(out.TableName) != f.tableName {
			return fmt.Errorf("TableName = %s, want %s", aws.ToString(out.TableName), f.tableName)
		}
		if len(out.KinesisDataStreamDestinations) != 0 {
			return fmt.Errorf("expected no destinations, got %+v", out.KinesisDataStreamDestinations)
		}
		return nil
	}))

	return results
}
