package testutil

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	cloudtraildata "github.com/aws/aws-sdk-go-v2/service/cloudtraildata"
	datatypes "github.com/aws/aws-sdk-go-v2/service/cloudtraildata/types"
)

// partnerAuditPayload builds one channel payload record and its base64
// SHA-256 checksum — the integrity mechanism the PutAuditEvents reference
// documents ("printf %s $eventdata | openssl dgst -binary -sha256 | base64").
func partnerAuditPayload(eventName string) (string, string) {
	payload := fmt.Sprintf(
		`{"eventVersion":"1.08","eventTime":%q,"eventSource":"example.partner.com","eventName":%q,"readOnly":false}`,
		time.Now().UTC().Format(time.RFC3339), eventName)
	digest := sha256.Sum256([]byte(payload))
	return payload, base64.StdEncoding.EncodeToString(digest[:])
}

// runLakeQueryAndWait runs one Lake query to completion and returns its
// flattened rows, draining the GetQueryResults pagination.
func (tc *cloudTrailTestContext) runLakeQueryAndWait(edsID, statement string) ([]map[string]string, error) {
	resp, err := tc.client.StartQuery(tc.ctx, &cloudtrail.StartQueryInput{
		QueryStatement: aws.String(statement),
	})
	if err != nil {
		return nil, fmt.Errorf("StartQuery failed: %w", err)
	}
	queryID := aws.ToString(resp.QueryId)

	finished := false
	for i := 0; i < 25; i++ {
		desc, derr := tc.client.DescribeQuery(tc.ctx, &cloudtrail.DescribeQueryInput{
			EventDataStore: aws.String(edsID),
			QueryId:        aws.String(queryID),
		})
		if derr != nil {
			return nil, fmt.Errorf("DescribeQuery failed: %w", derr)
		}
		if desc.QueryStatus == types.QueryStatusFinished {
			finished = true
			break
		}
		if desc.QueryStatus != types.QueryStatusRunning && desc.QueryStatus != types.QueryStatusQueued {
			return nil, fmt.Errorf("query ended with status %s", desc.QueryStatus)
		}
		time.Sleep(200 * time.Millisecond)
	}
	if !finished {
		return nil, fmt.Errorf("query %s did not finish", queryID)
	}

	var rows []map[string]string
	var nextToken *string
	for {
		out, gerr := tc.client.GetQueryResults(tc.ctx, &cloudtrail.GetQueryResultsInput{
			EventDataStore: aws.String(edsID),
			QueryId:        aws.String(queryID),
			NextToken:      nextToken,
		})
		if gerr != nil {
			return nil, fmt.Errorf("GetQueryResults failed: %w", gerr)
		}
		for _, row := range out.QueryResultRows {
			flat := map[string]string{}
			for _, entry := range row {
				for k, v := range entry {
					flat[k] = v
				}
			}
			rows = append(rows, flat)
		}
		if out.NextToken == nil || *out.NextToken == "" {
			return rows, nil
		}
		nextToken = out.NextToken
	}
}

// runCloudTrailDataTests pins the cloudtrail-data ingestion service and the
// event data stores' behaviour as data boundaries: channel events land in
// the channel's destination store, ingestion follows each store's
// selectors, and a stopped store accepts nothing.
func (r *TestRunner) runCloudTrailDataTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	var edsARN string
	var channelARN string

	results = append(results, r.RunTest("cloudtrail", "PutAuditEvents_Setup", func() error {
		resp, err := tc.createEventDataStore("ct-data-eds", nil)
		if err != nil {
			return fmt.Errorf("create destination EDS: %w", err)
		}
		edsARN = *resp.EventDataStoreArn
		ch, err := tc.client.CreateChannel(tc.ctx, &cloudtrail.CreateChannelInput{
			Name:   aws.String(tc.uniqueName("ct-data-channel")),
			Source: aws.String(tc.uniqueName("ct-data-src")),
			Destinations: []types.Destination{{
				Type:     types.DestinationTypeEventDataStore,
				Location: aws.String(edsARN),
			}},
		})
		if err != nil {
			return fmt.Errorf("create channel: %w", err)
		}
		channelARN = aws.ToString(ch.ChannelArn)
		return nil
	}))
	// The deferred belt deletes the channel first — the EDS delete is
	// refused while a channel still points at it — then the store; the
	// registered cleanup test owns the ordered teardown in the normal path.
	defer func() {
		if channelARN != "" {
			_, _ = tc.client.DeleteChannel(tc.ctx, &cloudtrail.DeleteChannelInput{
				Channel: aws.String(channelARN)})
		}
		if edsARN != "" {
			_ = tc.deleteEventDataStore(edsARN)
		}
	}()

	results = append(results, r.RunTest("cloudtrail", "PutAuditEvents_Success", func() error {
		payload, checksum := partnerAuditPayload("PartnerSignIn")
		resp, err := tc.dataClient.PutAuditEvents(tc.ctx, &cloudtraildata.PutAuditEventsInput{
			ChannelArn: aws.String(channelARN),
			AuditEvents: []datatypes.AuditEvent{{
				Id:                aws.String("partner-evt-1"),
				EventData:         aws.String(payload),
				EventDataChecksum: aws.String(checksum),
			}},
		})
		if err != nil {
			return fmt.Errorf("PutAuditEvents failed: %w", err)
		}
		if len(resp.Failed) != 0 {
			return fmt.Errorf("expected no failed entries, got %+v", resp.Failed)
		}
		if len(resp.Successful) != 1 {
			return fmt.Errorf("expected 1 successful entry, got %d", len(resp.Successful))
		}
		entry := resp.Successful[0]
		if aws.ToString(entry.Id) != "partner-evt-1" {
			return fmt.Errorf("successful entry id = %q", aws.ToString(entry.Id))
		}
		if aws.ToString(entry.EventID) == "" {
			return fmt.Errorf("CloudTrail-assigned eventID is empty")
		}
		return nil
	}))

	// A batch separates survivors from failures per entry: a checksum that
	// does not match the record and a record that does not parse both fail
	// alone, with the documented errorCode vocabulary.
	results = append(results, r.RunTest("cloudtrail", "PutAuditEvents_MixedBatch", func() error {
		good, goodChecksum := partnerAuditPayload("PartnerGood")
		bad, _ := partnerAuditPayload("PartnerBad")
		resp, err := tc.dataClient.PutAuditEvents(tc.ctx, &cloudtraildata.PutAuditEventsInput{
			ChannelArn: aws.String(channelARN),
			AuditEvents: []datatypes.AuditEvent{
				{Id: aws.String("mix-good"), EventData: aws.String(good), EventDataChecksum: aws.String(goodChecksum)},
				{Id: aws.String("mix-checksum"), EventData: aws.String(bad), EventDataChecksum: aws.String("bm90LXRoZS1jaGVja3N1bQ==")},
				{Id: aws.String("mix-json"), EventData: aws.String("not json")},
			},
		})
		if err != nil {
			return fmt.Errorf("PutAuditEvents failed: %w", err)
		}
		if len(resp.Successful) != 1 || aws.ToString(resp.Successful[0].Id) != "mix-good" {
			return fmt.Errorf("expected mix-good alone to succeed, got %+v", resp.Successful)
		}
		codes := map[string]string{}
		for _, f := range resp.Failed {
			codes[aws.ToString(f.Id)] = aws.ToString(f.ErrorCode)
		}
		if codes["mix-checksum"] != "InvalidChecksum" {
			return fmt.Errorf("checksum failure code = %q, want InvalidChecksum", codes["mix-checksum"])
		}
		if codes["mix-json"] != "InvalidData" {
			return fmt.Errorf("invalid data code = %q, want InvalidData", codes["mix-json"])
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutAuditEvents_DuplicateIds", func() error {
		payload, _ := partnerAuditPayload("PartnerDup")
		_, err := tc.dataClient.PutAuditEvents(tc.ctx, &cloudtraildata.PutAuditEventsInput{
			ChannelArn: aws.String(channelARN),
			AuditEvents: []datatypes.AuditEvent{
				{Id: aws.String("dup"), EventData: aws.String(payload)},
				{Id: aws.String("dup"), EventData: aws.String(payload)},
			},
		})
		return AssertErrorContains(err, "DuplicatedAuditEventId")
	}))

	// The channelArn member accepts "The ARN or ID (the ARN suffix) of a
	// channel"; both forms of an unknown channel answer ChannelNotFound,
	// and an ARN that is not a CloudTrail channel ARN answers
	// InvalidChannelARN.
	results = append(results, r.RunTest("cloudtrail", "PutAuditEvents_ChannelResolution", func() error {
		payload, _ := partnerAuditPayload("PartnerResolve")
		unknownARN := fmt.Sprintf("arn:aws:cloudtrail:%s:%s:channel/00000000-0000-0000-0000-000000000000",
			tc.region, tc.accountID)
		_, err := tc.dataClient.PutAuditEvents(tc.ctx, &cloudtraildata.PutAuditEventsInput{
			ChannelArn: aws.String(unknownARN),
			AuditEvents: []datatypes.AuditEvent{
				{Id: aws.String("r-1"), EventData: aws.String(payload)},
			},
		})
		if err := AssertErrorContains(err, "ChannelNotFound"); err != nil {
			return err
		}

		_, err = tc.dataClient.PutAuditEvents(tc.ctx, &cloudtraildata.PutAuditEventsInput{
			ChannelArn: aws.String("no-such-channel-suffix"),
			AuditEvents: []datatypes.AuditEvent{
				{Id: aws.String("r-2"), EventData: aws.String(payload)},
			},
		})
		if err := AssertErrorContains(err, "ChannelNotFound"); err != nil {
			return err
		}

		_, err = tc.dataClient.PutAuditEvents(tc.ctx, &cloudtraildata.PutAuditEventsInput{
			ChannelArn: aws.String("arn:aws:s3:::not-a-channel"),
			AuditEvents: []datatypes.AuditEvent{
				{Id: aws.String("r-3"), EventData: aws.String(payload)},
			},
		})
		return AssertErrorContains(err, "InvalidChannelARN")
	}))

	// The channel end-to-end: an ingested partner event is queryable in the
	// destination store through Lake, carrying the ActivityAuditLog
	// category.
	results = append(results, r.RunTest("cloudtrail", "PutAuditEvents_VisibleInLake", func() error {
		payload, checksum := partnerAuditPayload("PartnerVisibleEvt")
		if _, err := tc.dataClient.PutAuditEvents(tc.ctx, &cloudtraildata.PutAuditEventsInput{
			ChannelArn: aws.String(channelARN),
			AuditEvents: []datatypes.AuditEvent{{
				Id:                aws.String("visible-1"),
				EventData:         aws.String(payload),
				EventDataChecksum: aws.String(checksum),
			}},
		}); err != nil {
			return fmt.Errorf("PutAuditEvents failed: %w", err)
		}

		edsID := tc.edsIDFromARN(edsARN)
		rows, err := tc.runLakeQueryAndWait(edsID, fmt.Sprintf(
			"SELECT eventName FROM %s WHERE eventCategory = 'ActivityAuditLog'", edsID))
		if err != nil {
			return err
		}
		found := false
		for _, row := range rows {
			if row["eventName"] == "PartnerVisibleEvt" {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("ingested partner event not visible in the destination store (%d rows)", len(rows))
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutAuditEvents_Cleanup", func() error {
		if edsARN == "" || channelARN == "" {
			return fmt.Errorf("fixtures were not created")
		}
		// The channel goes first: DeleteEventDataStore refuses while a
		// channel still points at the store.
		if _, err := tc.client.DeleteChannel(tc.ctx, &cloudtrail.DeleteChannelInput{
			Channel: aws.String(channelARN),
		}); err != nil {
			return fmt.Errorf("delete channel: %w", err)
		}
		channelARN = ""
		if err := tc.deleteEventDataStore(edsARN); err != nil {
			return fmt.Errorf("delete destination EDS: %w", err)
		}
		edsARN = ""
		return nil
	}))

	// --- Event data stores as data boundaries ---

	// A store created without selectors materialises the documented default:
	// the "Default management events" selector set.
	results = append(results, r.RunTest("cloudtrail", "EDS_DefaultSelectorsEcho", func() error {
		resp, err := tc.createEventDataStore("ct-eds-default-sel", nil)
		if err != nil {
			return err
		}
		defer tc.deleteEventDataStore(*resp.EventDataStoreArn)

		if len(resp.AdvancedEventSelectors) != 1 {
			return fmt.Errorf("expected the default selector, got %d", len(resp.AdvancedEventSelectors))
		}
		sel := resp.AdvancedEventSelectors[0]
		if aws.ToString(sel.Name) != "Default management events" {
			return fmt.Errorf("default selector name = %q", aws.ToString(sel.Name))
		}
		if len(sel.FieldSelectors) != 1 || aws.ToString(sel.FieldSelectors[0].Field) != "eventCategory" {
			return fmt.Errorf("default selector fields = %+v", sel.FieldSelectors)
		}
		if len(sel.FieldSelectors[0].Equals) != 1 || sel.FieldSelectors[0].Equals[0] != "Management" {
			return fmt.Errorf("default selector values = %v", sel.FieldSelectors[0].Equals)
		}
		return nil
	}))

	// Ingestion follows each store's selectors: a management event reaches
	// the default store and never reaches a store whose selector excludes
	// it — no cross-store leak.
	results = append(results, r.RunTest("cloudtrail", "EDS_IngestionIsolation", func() error {
		defaultEDS, err := tc.createEventDataStore("ct-eds-iso-def", nil)
		if err != nil {
			return err
		}
		defer tc.deleteEventDataStore(*defaultEDS.EventDataStoreArn)

		pickyEDS, err := tc.client.CreateEventDataStore(tc.ctx, &cloudtrail.CreateEventDataStoreInput{
			Name:                         aws.String(tc.uniqueName("ct-eds-iso-picky")),
			TerminationProtectionEnabled: aws.Bool(false),
			AdvancedEventSelectors: []types.AdvancedEventSelector{{
				Name: aws.String("picky"),
				FieldSelectors: []types.AdvancedFieldSelector{{
					Field:  aws.String("eventName"),
					Equals: []string{"NoSuchOperationEver"},
				}},
			}},
		})
		if err != nil {
			return fmt.Errorf("create picky EDS: %w", err)
		}
		defer tc.deleteEventDataStore(*pickyEDS.EventDataStoreArn)

		// Bound the pin by run time so earlier runs' events cannot satisfy
		// it. The bound sits a minute in the past: the Lake engine compares
		// eventTime at second granularity, so a same-second seed collides
		// with the events it bounds.
		seedEpoch := time.Now().Unix() - 60

		// This call records the management event both stores evaluate.
		if _, err := tc.client.ListEventDataStores(tc.ctx, &cloudtrail.ListEventDataStoresInput{}); err != nil {
			return fmt.Errorf("ListEventDataStores: %w", err)
		}

		defaultID := tc.edsIDFromARN(*defaultEDS.EventDataStoreArn)
		pickyID := tc.edsIDFromARN(*pickyEDS.EventDataStoreArn)

		rows, err := tc.runLakeQueryAndWait(defaultID, fmt.Sprintf(
			"SELECT eventName FROM %s WHERE eventName = 'ListEventDataStores' AND eventTime > %d",
			defaultID, seedEpoch))
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			return fmt.Errorf("default store did not ingest the management event")
		}

		rows, err = tc.runLakeQueryAndWait(pickyID, fmt.Sprintf(
			"SELECT eventName FROM %s WHERE eventName = 'ListEventDataStores' AND eventTime > %d",
			pickyID, seedEpoch))
		if err != nil {
			return err
		}
		if len(rows) != 0 {
			return fmt.Errorf("selector-excluding store ingested the event (%d rows)", len(rows))
		}
		return nil
	}))

	// A store with ingestion stopped accepts nothing after the stop.
	results = append(results, r.RunTest("cloudtrail", "EDS_StopIngestion", func() error {
		eds, err := tc.createEventDataStore("ct-eds-stopped", nil)
		if err != nil {
			return err
		}
		defer tc.deleteEventDataStore(*eds.EventDataStoreArn)
		edsID := tc.edsIDFromARN(*eds.EventDataStoreArn)

		// A minute-past bound (the engine compares eventTime at second
		// granularity; a same-second seed collides with the bounded events).
		seedEpoch := time.Now().Unix() - 60
		if _, err := tc.client.ListEventDataStores(tc.ctx, &cloudtrail.ListEventDataStoresInput{}); err != nil {
			return fmt.Errorf("pre-stop ListEventDataStores: %w", err)
		}

		if _, err := tc.client.StopEventDataStoreIngestion(tc.ctx, &cloudtrail.StopEventDataStoreIngestionInput{
			EventDataStore: aws.String(edsID),
		}); err != nil {
			return fmt.Errorf("StopEventDataStoreIngestion: %w", err)
		}

		if _, err := tc.client.ListTrails(tc.ctx, &cloudtrail.ListTrailsInput{}); err != nil {
			return fmt.Errorf("post-stop ListTrails: %w", err)
		}

		rows, err := tc.runLakeQueryAndWait(edsID, fmt.Sprintf(
			"SELECT eventName FROM %s WHERE eventName = 'ListTrails' AND eventTime > %d", edsID, seedEpoch))
		if err != nil {
			return err
		}
		if len(rows) != 0 {
			return fmt.Errorf("stopped store ingested a post-stop event (%d rows)", len(rows))
		}
		return nil
	}))

	return results
}
