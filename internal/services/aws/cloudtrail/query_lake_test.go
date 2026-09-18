package cloudtrail

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awserrors "vorpalstacks/internal/common/errors"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// The pins below drive the Lake query engine's column vocabulary: SELECT
// projections and WHERE operands resolve case-insensitively against the
// AWS Lake SQL schema column names, decoupled from the LookupEvents wire
// formatter; a syntactically invalid WHERE clause and a column outside
// the schema are both rejected with InvalidQueryStatementException.

func newQueryTestStore(t *testing.T) *cloudtrailstore.CloudTrailStore {
	t.Helper()
	return newLockTestStore(t)
}

// seedLakeEDS provisions the event data store the Lake fixtures query: its
// identifier is the literal the FROM clauses name, so PutEvent's fan-out
// copies the seeded management events into exactly the store executeQuery
// walks.
func seedLakeEDS(t *testing.T, store *cloudtrailstore.CloudTrailStore) {
	t.Helper()
	eds := cloudtrailstore.NewEventDataStore("lake-fixture", store.GetAccountID(), store.GetRegion())
	eds.EventDataStoreID = "eds"
	if _, err := store.CreateEventDataStore(eds); err != nil {
		t.Fatalf("create fixture EDS failed: %v", err)
	}
}

func seedLakeEvent(t *testing.T, store *cloudtrailstore.CloudTrailStore, userName, eventName string, readOnly bool, resources []cloudtrailstore.Resource) *cloudtrailstore.Event {
	t.Helper()
	identity := &cloudtrailstore.UserIdentity{
		Type:     "IAMUser",
		UserName: userName,
		ARN:      "arn:aws:iam::acc123:user/" + userName,
	}
	e := cloudtrailstore.NewEvent(eventName, "cloudtrail.amazonaws.com", identity, readOnly)
	e.AwsRegion = "us-east-1"
	e.CloudTrailEvent = `{"eventVersion":"1.08","eventName":"` + eventName + `"}`
	e.Resources = resources
	require.NoError(t, store.PutEvent(e))
	return e
}

func requireInvalidQueryStatement(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	var apiErr *awserrors.AWSError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, "InvalidQueryStatementException", apiErr.GetCode())
}

func TestParseQueryStatementLakeColumns(t *testing.T) {
	t.Run("projection spellings are kept verbatim", func(t *testing.T) {
		pq, err := parseQueryStatement("SELECT eventID, eventTime FROM eds-12345")
		require.NoError(t, err)
		assert.Equal(t, []string{"eventID", "eventTime"}, pq.columns)
		assert.Equal(t, "eds-12345", pq.edsID)
	})

	t.Run("uppercase and lowercase columns resolve", func(t *testing.T) {
		for _, stmt := range []string{
			"SELECT EVENTID FROM eds",
			"select eventname from eds",
			"SELECT EventSource FROM eds",
		} {
			_, err := parseQueryStatement(stmt)
			require.NoError(t, err, stmt)
		}
	})

	t.Run("qualified struct columns are accepted", func(t *testing.T) {
		for _, stmt := range []string{
			"SELECT eventID FROM eds WHERE userIdentity.userName = 'alice'",
			"SELECT eventID FROM eds WHERE resources.ARN LIKE '%trail%'",
			"SELECT eventID FROM eds WHERE resources.type = 'AWS::CloudTrail::Trail'",
		} {
			pq, err := parseQueryStatement(stmt)
			require.NoError(t, err, stmt)
			assert.NotNil(t, pq.whereExpr, stmt)
		}
	})

	t.Run("unknown projection column is rejected", func(t *testing.T) {
		_, err := parseQueryStatement("SELECT nosuchcol FROM eds")
		requireInvalidQueryStatement(t, err)
	})

	t.Run("username is not a Lake column", func(t *testing.T) {
		// The schema carries userIdentity.userName; the bare LookupEvents
		// wire name is not queryable.
		_, err := parseQueryStatement("SELECT username FROM eds")
		requireInvalidQueryStatement(t, err)
		_, err = parseQueryStatement("SELECT eventID FROM eds WHERE username = 'alice'")
		requireInvalidQueryStatement(t, err)
	})

	t.Run("unknown WHERE column is rejected", func(t *testing.T) {
		_, err := parseQueryStatement("SELECT eventID FROM eds WHERE wat = 'x'")
		requireInvalidQueryStatement(t, err)
	})

	t.Run("unknown qualified root is rejected", func(t *testing.T) {
		_, err := parseQueryStatement("SELECT eventID FROM eds WHERE trail.name = 'x'")
		requireInvalidQueryStatement(t, err)
	})

	t.Run("malformed WHERE is rejected fail-closed", func(t *testing.T) {
		_, err := parseQueryStatement("SELECT eventID FROM eds WHERE eventName = 'unclosed")
		requireInvalidQueryStatement(t, err)
	})
}

func TestLakeRow(t *testing.T) {
	store := newQueryTestStore(t)
	e := seedLakeEvent(t, store, "alice", "CreateTrail", false,
		[]cloudtrailstore.Resource{{ResourceType: "AWS::CloudTrail::Trail", ResourceName: "arn:aws:cloudtrail:us-east-1:acc123:trail/t1"}})

	row := lakeRow(e)
	assert.Equal(t, e.EventID, row["eventid"])
	assert.Equal(t, "CreateTrail", row["eventname"])
	assert.Equal(t, "false", row["readonly"])
	assert.Equal(t, "us-east-1", row["awsregion"])
	assert.Equal(t, true, row["managementevent"])
	assert.Equal(t, e.CloudTrailEvent, row["eventjson"])

	identity, ok := row["useridentity"].(map[string]interface{})
	require.True(t, ok, "useridentity must be a nested map")
	assert.Equal(t, "alice", identity["username"])
	assert.Equal(t, "IAMUser", identity["type"])

	resources, ok := row["resources"].([]map[string]interface{})
	require.True(t, ok, "resources must be a slice of maps")
	require.Len(t, resources, 1)
	assert.Equal(t, "arn:aws:cloudtrail:us-east-1:acc123:trail/t1", resources[0]["arn"])
	assert.Equal(t, "AWS::CloudTrail::Trail", resources[0]["type"])

	// Schema columns the store does not record resolve to nil, never a
	// fabricated value.
	assert.Nil(t, resolveLakeColumn(row, "", "recipientaccountid"))
}

func TestExecuteQueryValues(t *testing.T) {
	store := newQueryTestStore(t)
	svc := NewCloudTrailService("acc123", "us-east-1")
	seedLakeEDS(t, store)

	created := seedLakeEvent(t, store, "alice", "CreateTrail", false,
		[]cloudtrailstore.Resource{{ResourceType: "AWS::CloudTrail::Trail", ResourceName: "arn:aws:cloudtrail:us-east-1:acc123:trail/t1"}})
	deleted := seedLakeEvent(t, store, "bob", "DeleteTrail", false, nil)
	listed := seedLakeEvent(t, store, "carol", "ListTrails", true, nil)

	t.Run("equality WHERE returns the matching row with values", func(t *testing.T) {
		pq, err := parseQueryStatement("SELECT eventID, eventName FROM eds WHERE eventName = 'CreateTrail'")
		require.NoError(t, err)
		rows, stats, _, err := svc.executeQuery(store, pq, time.Time{})
		require.NoError(t, err)
		require.Len(t, rows, 1)
		require.Len(t, rows[0], 2)
		assert.Equal(t, created.EventID, rows[0][0]["eventID"])
		assert.Equal(t, "CreateTrail", rows[0][1]["eventName"])
		// The per-EDS buckets carry no secondary indexes: the walk examines
		// every event in the store and the WHERE decides, so the scan saw
		// all three records and matched one.
		assert.Equal(t, int64(3), stats.eventsScanned)
		assert.Equal(t, int64(1), stats.eventsMatched)
		assert.Equal(t, int64(len(created.CloudTrailEvent)+len(deleted.CloudTrailEvent)+len(listed.CloudTrailEvent)), stats.bytesScanned)
	})

	t.Run("LIKE WHERE filters", func(t *testing.T) {
		pq, err := parseQueryStatement("SELECT eventname FROM eds WHERE eventName LIKE 'Create%'")
		require.NoError(t, err)
		rows, stats, _, err := svc.executeQuery(store, pq, time.Time{})
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "CreateTrail", rows[0][0]["eventname"])
		// LIKE cannot ride the index prefilter: every event was examined,
		// one matched, and the scan bytes are the sum of the record JSON
		// lengths.
		assert.Equal(t, int64(3), stats.eventsScanned)
		assert.Equal(t, int64(1), stats.eventsMatched)
		assert.Equal(t, int64(len(created.CloudTrailEvent)+len(deleted.CloudTrailEvent)+len(listed.CloudTrailEvent)), stats.bytesScanned)
	})

	t.Run("readonly boolean literal filters", func(t *testing.T) {
		pq, err := parseQueryStatement("SELECT eventName FROM eds WHERE readonly = true")
		require.NoError(t, err)
		rows, _, _, err := svc.executeQuery(store, pq, time.Time{})
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "ListTrails", rows[0][0]["eventName"])
	})

	t.Run("eventtime compares against epoch literals", func(t *testing.T) {
		pq, err := parseQueryStatement("SELECT eventName FROM eds WHERE eventTime > 0")
		require.NoError(t, err)
		rows, _, _, err := svc.executeQuery(store, pq, time.Time{})
		require.NoError(t, err)
		assert.Len(t, rows, 3)
	})

	t.Run("qualified useridentity column filters", func(t *testing.T) {
		pq, err := parseQueryStatement("SELECT eventName FROM eds WHERE userIdentity.userName = 'bob'")
		require.NoError(t, err)
		rows, _, _, err := svc.executeQuery(store, pq, time.Time{})
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "DeleteTrail", rows[0][0]["eventName"])
	})

	t.Run("qualified resources column filters", func(t *testing.T) {
		pq, err := parseQueryStatement("SELECT eventName FROM eds WHERE resources.type = 'AWS::CloudTrail::Trail'")
		require.NoError(t, err)
		rows, _, _, err := svc.executeQuery(store, pq, time.Time{})
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "CreateTrail", rows[0][0]["eventName"])
	})

	t.Run("eventTime projection renders record-format timestamps", func(t *testing.T) {
		pq, err := parseQueryStatement("SELECT eventTime FROM eds WHERE eventName = 'CreateTrail'")
		require.NoError(t, err)
		rows, _, _, err := svc.executeQuery(store, pq, time.Time{})
		require.NoError(t, err)
		require.Len(t, rows, 1)
		val := rows[0][0]["eventTime"]
		assert.NotEmpty(t, val)
		assert.NotContains(t, val, ".")
	})
}

// SQL NULL semantics and the array columns' equality-only contract: a NULL
// operand (an unrecorded column, or a NULL literal) satisfies no operator —
// not even the negated ones, which would otherwise turn an unknown into a
// match — and the resources fields carry no order, so ordering operators
// and both BETWEEN forms never claim a row through them. Equality on an
// array column keeps the any-entry meaning, and negated operators keep
// working on present scalar values.
func TestExecuteQueryNullAndArrayComparisonSemantics(t *testing.T) {
	store := newQueryTestStore(t)
	svc := NewCloudTrailService("acc123", "us-east-1")
	seedLakeEDS(t, store)

	seedLakeEvent(t, store, "alice", "CreateTrail", false,
		[]cloudtrailstore.Resource{{ResourceType: "AWS::CloudTrail::Trail", ResourceName: "arn:aws:cloudtrail:us-east-1:acc123:trail/t1"}})
	seedLakeEvent(t, store, "bob", "DeleteTrail", false, nil)

	cases := []struct {
		name string
		stmt string
		want int
	}{
		{"null column satisfies no negated operator", "SELECT eventID FROM eds WHERE recipientAccountID != '123'", 0},
		{"null column satisfies no ordering operator", "SELECT eventID FROM eds WHERE recipientAccountID < '123'", 0},
		{"null column satisfies no inclusive ordering", "SELECT eventID FROM eds WHERE recipientAccountID <= '123'", 0},
		{"null column satisfies no NOT IN", "SELECT eventID FROM eds WHERE recipientAccountID NOT IN ('123', '456')", 0},
		{"null column satisfies no NOT BETWEEN", "SELECT eventID FROM eds WHERE recipientAccountID NOT BETWEEN '0' AND '9'", 0},
		{"null column satisfies no BETWEEN", "SELECT eventID FROM eds WHERE recipientAccountID BETWEEN '0' AND '9'", 0},
		{"NULL literal operand satisfies no operator", "SELECT eventID FROM eds WHERE eventName != NULL", 0},
		{"NULL list item leaves NOT IN unknown", "SELECT eventID FROM eds WHERE eventName NOT IN ('CreateTrail', NULL)", 0},
		{"IS NULL is the null test", "SELECT eventID FROM eds WHERE recipientAccountID IS NULL", 2},
		{"negated operator still works on present values", "SELECT eventID FROM eds WHERE eventName != 'CreateTrail'", 1},
		{"array column satisfies no ordering operator", "SELECT eventID FROM eds WHERE resources.ARN > 'arn'", 0},
		{"array column satisfies no inclusive ordering", "SELECT eventID FROM eds WHERE resources.type >= 'A'", 0},
		{"array column satisfies no BETWEEN", "SELECT eventID FROM eds WHERE resources.type BETWEEN 'A' AND 'z'", 0},
		{"array column satisfies no NOT BETWEEN", "SELECT eventID FROM eds WHERE resources.type NOT BETWEEN 'A' AND 'z'", 0},
		{"array column equality matches any entry", "SELECT eventID FROM eds WHERE resources.ARN = 'arn:aws:cloudtrail:us-east-1:acc123:trail/t1'", 1},
		{"array column IN matches any entry", "SELECT eventID FROM eds WHERE resources.type IN ('AWS::CloudTrail::Trail', 'AWS::S3::Bucket')", 1},
		{"array column NOT IN keeps any-entry meaning", "SELECT eventID FROM eds WHERE resources.type NOT IN ('AWS::S3::Bucket')", 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pq, err := parseQueryStatement(tc.stmt)
			require.NoError(t, err)
			rows, _, _, err := svc.executeQuery(store, pq, time.Time{})
			require.NoError(t, err)
			assert.Len(t, rows, tc.want, tc.stmt)
		})
	}
}

// The scan must follow the store's NextToken to exhaustion: a store larger
// than one LakeQueryScanBound page is examined in full, and the statistics
// report every scanned event.
func TestExecuteQueryFollowsNextToken(t *testing.T) {
	store := newQueryTestStore(t)
	svc := NewCloudTrailService("acc123", "us-east-1")
	seedLakeEDS(t, store)

	total := int(cloudtrailstore.LakeQueryScanBound) + 5
	for i := 0; i < total; i++ {
		seedLakeEvent(t, store, "alice", "PaginatedEvent", false, nil)
	}

	pq, err := parseQueryStatement("SELECT eventID FROM eds WHERE eventName = 'PaginatedEvent'")
	require.NoError(t, err)
	rows, stats, _, err := svc.executeQuery(store, pq, time.Time{})
	require.NoError(t, err)
	assert.Len(t, rows, total)
	assert.Equal(t, int64(total), stats.eventsScanned)
	assert.Equal(t, int64(total), stats.eventsMatched)
	assert.Greater(t, stats.bytesScanned, int64(0))
}
