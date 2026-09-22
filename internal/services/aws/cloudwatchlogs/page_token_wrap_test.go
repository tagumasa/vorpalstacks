package cloudwatchlogs

import (
	"fmt"
	"testing"
	"time"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The listing cores whose page markers previously rode bare natural
// keys: the marker now travels inside the scoped listing vocabulary, so
// a raw key (the pre-fix wire form) is not a token and rejects, the
// wrapped marker pages onward, and a marker replayed against another
// family's listing rejects with it — the tokens "expire after 24 hours"
// (the describe listings' nextToken members).
func TestListingCoresRideScopedMarkers(t *testing.T) {
	t.Run("account policies", func(t *testing.T) {
		svc, store := newReadTestService(t, "wrap-ap-group")
		for _, name := range []string{"a-ap", "b-ap"} {
			if err := store.PutAccountPolicy(&logsstore.AccountPolicy{
				PolicyName: name, PolicyType: "SUBSCRIPTION_FILTER_POLICY", PolicyDocument: "{}",
			}); err != nil {
				t.Fatal(err)
			}
		}
		if _, _, err := svc.describeAccountPoliciesCore("SUBSCRIPTION_FILTER_POLICY", "", nil, "a-ap", "us-east-1"); err == nil {
			t.Fatal("raw policy-name marker accepted")
		}
	})

	t.Run("destinations", func(t *testing.T) {
		svc, store := newReadTestService(t, "wrap-dest-group")
		for _, name := range []string{"a-dest", "b-dest"} {
			if err := store.PutDestination(&logsstore.Destination{
				Name: name, RoleArn: "arn:aws:iam::000000000000:role/r", TargetArn: "arn:aws:kinesis:us-east-1:000000000000:stream/s",
			}); err != nil {
				t.Fatal(err)
			}
		}
		if _, _, err := svc.describeDestinationsCore("", "a-dest", "us-east-1", 1); err == nil {
			t.Fatal("raw destination-name marker accepted")
		}
		page1, next, err := svc.describeDestinationsCore("", "", "us-east-1", 1)
		if err != nil || len(page1) != 1 || next == "" {
			t.Fatalf("page 1 = %d items (token %q, err %v)", len(page1), next, err)
		}
		if next == page1[0].Name {
			t.Fatal("marker is the bare name, not the scoped vocabulary")
		}
		page2, _, err := svc.describeDestinationsCore("", next, "us-east-1", 1)
		if err != nil || len(page2) != 1 || page2[0].Name != "b-dest" {
			t.Fatalf("page 2 = %+v (err %v)", page2, err)
		}
	})

	t.Run("delivery destinations", func(t *testing.T) {
		svc, store := newReadTestService(t, "wrap-dd-group")
		for _, name := range []string{"a-dd", "b-dd"} {
			if err := store.PutDeliveryDestination(&logsstore.DeliveryDestination{
				Name: name, DeliveryDestinationType: "S3",
				DestinationResourceArn: "arn:aws:s3:::wrap-dd-bucket",
			}); err != nil {
				t.Fatal(err)
			}
		}
		if _, _, err := svc.describeDeliveryDestinationsCore("us-east-1", "a-dd", 1); err == nil {
			t.Fatal("raw delivery-destination-name marker accepted")
		}
		page1, next, err := svc.describeDeliveryDestinationsCore("us-east-1", "", 1)
		if err != nil || len(page1) != 1 || next == "" {
			t.Fatalf("page 1 = %d items (token %q, err %v)", len(page1), next, err)
		}
		page2, _, err := svc.describeDeliveryDestinationsCore("us-east-1", next, 1)
		if err != nil || len(page2) != 1 {
			t.Fatalf("page 2 = %d items (err %v)", len(page2), err)
		}
		// A marker from another family's scope repositions nothing.
		if _, _, err := svc.describeDestinationsCore("", next, "us-east-1", 1); err == nil {
			t.Fatal("delivery-family marker accepted by the destinations listing")
		}
	})

	t.Run("configuration templates", func(t *testing.T) {
		svc, _ := newReadTestService(t, "wrap-tmpl-group")
		page1, next, err := svc.describeConfigurationTemplatesCore(&DescribeConfigurationTemplatesInput{Limit: 1})
		if err != nil || len(page1) != 1 || next == "" {
			t.Fatalf("page 1 = %d rows (token %q, err %v)", len(page1), next, err)
		}
		rawKey := fmt.Sprintf("%v", page1[0]["deliveryDestinationType"])
		if _, _, err := svc.describeConfigurationTemplatesCore(&DescribeConfigurationTemplatesInput{NextToken: rawKey, Limit: 1}); err == nil {
			t.Fatal("raw catalog-key marker accepted")
		}
		page2, _, err := svc.describeConfigurationTemplatesCore(&DescribeConfigurationTemplatesInput{NextToken: next, Limit: 1})
		if err != nil || len(page2) != 1 {
			t.Fatalf("page 2 = %d rows (err %v)", len(page2), err)
		}
	})

	t.Run("export tasks", func(t *testing.T) {
		svc, store := newReadTestService(t, "wrap-ext-group")
		for _, id := range []string{"a-ext", "b-ext"} {
			if err := store.PutExportTask(&logsstore.ExportTask{
				TaskId: id, LogGroupName: "wrap-ext-group", Status: "COMPLETED",
				CreationTime: time.Now().UnixMilli(),
			}); err != nil {
				t.Fatal(err)
			}
		}
		if _, _, err := svc.describeExportTasksCore(store, &DescribeExportTasksInput{NextToken: "a-ext", Limit: 1}); err == nil {
			t.Fatal("raw task-id marker accepted")
		}
		page1, next, err := svc.describeExportTasksCore(store, &DescribeExportTasksInput{Limit: 1})
		if err != nil || len(page1) != 1 || next == "" {
			t.Fatalf("page 1 = %d tasks (token %q, err %v)", len(page1), next, err)
		}
		page2, _, err := svc.describeExportTasksCore(store, &DescribeExportTasksInput{NextToken: next, Limit: 1})
		if err != nil || len(page2) != 1 {
			t.Fatalf("page 2 = %d tasks (err %v)", len(page2), err)
		}
	})

	t.Run("import tasks", func(t *testing.T) {
		svc, store := newReadTestService(t, "wrap-imt-group")
		for _, id := range []string{"a-imt", "b-imt"} {
			if err := store.PutImportTask(&logsstore.ImportTask{
				ImportId: id, ImportStatus: "COMPLETED", LogGroupName: "wrap-imt-group",
			}); err != nil {
				t.Fatal(err)
			}
		}
		if _, _, err := svc.describeImportTasksCore(store, "", "", "", "a-imt", 1); err == nil {
			t.Fatal("raw import-id marker accepted")
		}
		page1, next, err := svc.describeImportTasksCore(store, "", "", "", "", 1)
		if err != nil || len(page1) != 1 || next == "" {
			t.Fatalf("page 1 = %d tasks (token %q, err %v)", len(page1), next, err)
		}
		page2, _, err := svc.describeImportTasksCore(store, "", "", "", next, 1)
		if err != nil || len(page2) != 1 {
			t.Fatalf("page 2 = %d tasks (err %v)", len(page2), err)
		}
	})

	t.Run("import task batches", func(t *testing.T) {
		svc, store := newReadTestService(t, "wrap-imb-group")
		if err := store.PutImportTask(&logsstore.ImportTask{
			ImportId: "imb-1", ImportStatus: "COMPLETED", LogGroupName: "wrap-imb-group",
			ImportBatches: []*logsstore.ImportBatch{
				{BatchId: "b1", Status: "COMPLETED"},
				{BatchId: "b2", Status: "COMPLETED"},
			},
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := svc.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{ImportId: "imb-1", NextToken: "0"}); err == nil {
			t.Fatal("raw positional marker accepted")
		}
		page1, err := svc.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{ImportId: "imb-1", Limit: 1})
		if err != nil || len(page1.ImportBatches) != 1 || page1.NextToken == "" {
			t.Fatalf("page 1 = %d batches (token %q, err %v)", len(page1.ImportBatches), page1.NextToken, err)
		}
		page2, err := svc.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{ImportId: "imb-1", NextToken: page1.NextToken, Limit: 1})
		if err != nil || len(page2.ImportBatches) != 1 {
			t.Fatalf("page 2 = %d batches (err %v)", len(page2.ImportBatches), err)
		}
	})

	t.Run("scheduled query executions", func(t *testing.T) {
		svc, store := newReadTestService(t, "wrap-sqx-group")
		if err := store.PutScheduledQuery(&logsstore.ScheduledQuery{Id: "sq-1", Name: "wrap"}); err != nil {
			t.Fatal(err)
		}
		now := time.Now().UnixMilli()
		for _, tr := range []int64{now - 2000, now - 1000} {
			if err := store.PutScheduledQueryExecution(&logsstore.ScheduledQueryExecution{
				ScheduledQueryId: "sq-1", TriggerTime: tr, Status: "COMPLETED",
			}); err != nil {
				t.Fatal(err)
			}
		}
		rawKey := fmt.Sprintf("%d", now-2000)
		if _, err := svc.getScheduledQueryHistoryCore(store, &GetScheduledQueryHistoryInput{
			Identifier: "sq-1", StartTime: 0, EndTime: now, StartTimeSet: true, EndTimeSet: true,
			NextToken: rawKey,
		}); err == nil {
			t.Fatal("raw trigger-time marker accepted")
		}
		page1, err := svc.getScheduledQueryHistoryCore(store, &GetScheduledQueryHistoryInput{
			Identifier: "sq-1", StartTime: 0, EndTime: now, StartTimeSet: true, EndTimeSet: true, MaxResults: 1,
		})
		if err != nil || len(page1.Executions) != 1 || page1.NextMarker == "" {
			t.Fatalf("page 1 = %d executions (token %q, err %v)", len(page1.Executions), page1.NextMarker, err)
		}
		page2, err := svc.getScheduledQueryHistoryCore(store, &GetScheduledQueryHistoryInput{
			Identifier: "sq-1", StartTime: 0, EndTime: now, StartTimeSet: true, EndTimeSet: true,
			NextToken: page1.NextMarker, MaxResults: 1,
		})
		if err != nil || len(page2.Executions) != 1 {
			t.Fatalf("page 2 = %d executions (err %v)", len(page2.Executions), err)
		}
	})

	t.Run("scheduled queries", func(t *testing.T) {
		svc, store := newReadTestService(t, "wrap-sql-group")
		for _, id := range []string{"a-sql", "b-sql"} {
			if err := store.PutScheduledQuery(&logsstore.ScheduledQuery{Id: id, Name: id}); err != nil {
				t.Fatal(err)
			}
		}
		if _, _, err := svc.listScheduledQueriesCore(store, "", "", 1, "a-sql"); err == nil {
			t.Fatal("raw query-id marker accepted")
		}
		page1, next, err := svc.listScheduledQueriesCore(store, "", "", 1, "")
		if err != nil || len(page1) != 1 || next == "" {
			t.Fatalf("page 1 = %d queries (token %q, err %v)", len(page1), next, err)
		}
		page2, _, err := svc.listScheduledQueriesCore(store, "", "", 1, next)
		if err != nil || len(page2) != 1 {
			t.Fatalf("page 2 = %d queries (err %v)", len(page2), err)
		}
	})
}
