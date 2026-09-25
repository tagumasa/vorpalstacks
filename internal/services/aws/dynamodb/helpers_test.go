package dynamodb

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// serviceRegionFixture builds the single-region service harness: a service
// over a temp-dir storage manager (closed with the test), its request
// context, and its cached region store. Tests that create their own tables
// — several, or one whose record they keep — start here.
func serviceRegionFixture(t *testing.T) (*DynamoDBService, *request.RequestContext, dbstore.DynamoDBStoreInterface) {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)
	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	return svc, request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1"), store
}

// serviceTableFixture builds the region harness with one created table.
func serviceTableFixture(t *testing.T, params dbstore.CreateTableParams) (*DynamoDBService, *request.RequestContext, dbstore.DynamoDBStoreInterface) {
	t.Helper()
	svc, reqCtx, store := serviceRegionFixture(t)
	if _, err := store.Tables().Create(params); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return svc, reqCtx, store
}

// ttlFaultTables makes the TTL precondition read fail while every other
// table operation runs through the real store — the pin's way of holding
// the read fault and the write apart.
type ttlFaultTables struct {
	dbstore.TableStoreInterface
}

func (w ttlFaultTables) GetTimeToLive(name string) (*dbstore.TimeToLiveSpecification, error) {
	return nil, errors.New("ttl precondition read fault")
}

type ttlFaultStore struct {
	dbstore.DynamoDBStoreInterface
}

func (w ttlFaultStore) Tables() dbstore.TableStoreInterface {
	return ttlFaultTables{w.DynamoDBStoreInterface.Tables()}
}

// snapshotFaultBucket fails deletes of snapshot keys while every other
// operation runs through the real bucket — the pin's way of holding the
// snapshot-delete fault and the record-delete apart.
type snapshotFaultBucket struct {
	storage.Bucket
}

func (b snapshotFaultBucket) Delete(key []byte) error {
	if len(key) > 0 && key[0] == '#' {
		return errors.New("snapshot delete fault")
	}
	return b.Bucket.Delete(key)
}

// TestStorageFaultsNeverMasqueradeAsAbsence pins the documented
// fault-versus-absence split on the paths outside the describeByArn family:
// the backup delete read, the three stream cores' table read, the TTL
// duplicate-enable precondition read, and the backup snapshot delete. A
// storage fault is reported as the storage error it is (or, for the
// snapshot, fails the whole deletion with the record kept for retry) —
// never flattened into a not-found answer, a silently passed rejection, or
// an orphaned snapshot.
func TestStorageFaultsNeverMasqueradeAsAbsence(t *testing.T) {
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)
	reqCtx := request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")
	ctx := context.Background()
	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	regionStorage, err := sm.GetStorage("us-east-1")
	if err != nil {
		t.Fatalf("region storage: %v", err)
	}
	isStorageFault := func(err error) bool {
		return err != nil && !errors.Is(err, ErrResourceNotFound) && !errors.Is(err, ErrBackupNotFound)
	}

	// A streamed table and a backup of it, both healthy to start with.
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "FaultTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
		StreamSpecification:  &dbstore.StreamSpecification{StreamEnabled: true, StreamViewType: dbstore.StreamViewTypeNewAndOldImages},
	}); err != nil {
		t.Fatalf("create streamed table: %v", err)
	}
	backupResp, err := svc.CreateBackup(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "FaultTable", "BackupName": "fault-bk",
	}})
	if err != nil {
		t.Fatalf("create backup: %v", err)
	}
	backupArn := backupResp.(map[string]interface{})["BackupDetails"].(map[string]interface{})["BackupArn"].(string)

	// The backup delete: an unreadable record is the storage error, not
	// the family's not-found answer; an absent record still is.
	if err := regionStorage.Bucket("dynamodb_backups-us-east-1").Put([]byte(svcarn.ExtractBackupIdFromARN(backupArn)), []byte("not-a-protobuf-record")); err != nil {
		t.Fatalf("corrupt backup record: %v", err)
	}
	if _, err := svc.DeleteBackup(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"BackupArn": backupArn,
	}}); !isStorageFault(err) {
		t.Fatalf("backup delete on corrupt record: err = %v, want the storage fault", err)
	}
	if _, err := svc.DeleteBackup(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"BackupArn": "arn:aws:dynamodb:us-east-1:123456789012:table/FaultTable/backup/01234567890123-aaaaaaaa",
	}}); !errors.Is(err, ErrBackupNotFound) {
		t.Fatalf("backup delete on absent record: err = %v, want ErrBackupNotFound", err)
	}

	// The stream cores: an iterator minted while the record was healthy,
	// then an unreadable table record — describe, iterator issue, and
	// read all answer the storage fault, not ResourceNotFound.
	table, err := store.Tables().Get("FaultTable")
	if err != nil {
		t.Fatalf("read table: %v", err)
	}
	if _, err := svc.GetShardIterator(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": table.StreamArn, "ShardId": dbstore.ShardIDForStream(table.StreamArn), "ShardIteratorType": "LATEST",
	}}); err != nil {
		t.Fatalf("mint iterator: %v", err)
	}
	if err := regionStorage.Bucket("dynamodb_tables-us-east-1").Put([]byte("FaultTable"), []byte("not-a-protobuf-record")); err != nil {
		t.Fatalf("corrupt table record: %v", err)
	}
	if _, err := svc.DescribeStream(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": table.StreamArn,
	}}); !isStorageFault(err) {
		t.Fatalf("describe stream on corrupt record: err = %v, want the storage fault", err)
	}
	if _, err := svc.GetShardIterator(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": table.StreamArn, "ShardId": dbstore.ShardIDForStream(table.StreamArn), "ShardIteratorType": "LATEST",
	}}); !isStorageFault(err) {
		t.Fatalf("iterator issue on corrupt record: err = %v, want the storage fault", err)
	}
	signingKey, err := store.Streams().IteratorSigningKey()
	if err != nil {
		t.Fatalf("signing key: %v", err)
	}
	iterator := encodeShardIterator(signingKey, table.StreamArn, "FaultTable", 0, "LATEST")
	if _, err := svc.GetRecords(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ShardIterator": iterator,
	}}); !isStorageFault(err) {
		t.Fatalf("records read on corrupt record: err = %v, want the storage fault", err)
	}

	// The TTL precondition: a faulted read must not silently pass the
	// duplicate-enable rejection and overwrite an unreadable state.
	if _, err := svc.updateTimeToLiveCore(ctx, ttlFaultStore{store}, UpdateTimeToLiveInput{
		TableName: "FaultTable", AttributeName: "ttl", Enabled: true, SpecPresent: true,
	}); err == nil || err.Error() != "ttl precondition read fault" {
		t.Fatalf("ttl enable under a faulted precondition read: err = %v, want the fault", err)
	}

	// The snapshot delete: its failure fails the whole deletion and keeps
	// the record for a retry, instead of deleting the record and
	// orphaning the snapshot bytes.
	bk, err := store.Backups().Create("fault-snap", "FaultTable", table.ARN, 0)
	if err != nil {
		t.Fatalf("create snapshot-bearing backup: %v", err)
	}
	real, ok := store.Backups().(*dbstore.BackupStore)
	if !ok {
		t.Fatalf("backups store concrete type %T", store.Backups())
	}
	faulting := *real
	faulting.BaseStore = commonstore.NewBaseStore(snapshotFaultBucket{real.Bucket()}, "dynamodb_backups")
	if err := faulting.Delete(bk.BackupArn); err == nil {
		t.Fatal("snapshot delete fault must fail the whole deletion")
	}
	if _, err := store.Backups().Get(bk.BackupArn); err != nil {
		t.Fatal("record deleted under a failed snapshot delete — orphaned snapshot, unretryable cleanup")
	}
}

// TestDescribeByArnDistinguishesAbsenceFromStorageFaults pins the
// ARN-addressed describe mapping: a not-found-class read maps to the
// family's sentinel, a storage fault (bucket I/O, an unreadable record)
// propagates as the storage error it is — never masquerading as absence.
func TestDescribeByArnDistinguishesAbsenceFromStorageFaults(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)

	sentinel := NewAPIError("com.amazonaws.dynamodb.v20120810#BackupNotFoundException", "Backup not found", http.StatusBadRequest)
	get := func(err error) func(dbstore.DynamoDBStoreInterface) (*dbstore.Backup, error) {
		return func(dbstore.DynamoDBStoreInterface) (*dbstore.Backup, error) {
			return nil, err
		}
	}

	if _, err := describeByArn(svc, reqCtx, "arn:aws:dynamodb:us-east-1:123456789012:table/t/backup/1234-x", true, sentinel,
		get(commonstore.NewStoreErrorWithKey("dynamodb-backups", "get_proto", "k", commonstore.ErrNotFound))); !errors.Is(err, sentinel) {
		t.Fatalf("absent record: err = %v, want the family sentinel", err)
	}

	fault := errors.New("pebble: closed")
	_, err := describeByArn(svc, reqCtx, "arn:aws:dynamodb:us-east-1:123456789012:table/t/backup/1234-x", true, sentinel,
		get(commonstore.NewStoreErrorWithKey("dynamodb-backups", "get_proto", "k", fault)))
	if !errors.Is(err, fault) {
		t.Fatalf("storage fault: err = %v, want the fault to propagate", err)
	}
}
