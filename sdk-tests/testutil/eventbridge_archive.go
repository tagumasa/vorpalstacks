package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func (r *TestRunner) runEventBridgeArchiveTests(ctx context.Context, client *eventbridge.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("events", "CreateArchive", func() error {
		caBus := fmt.Sprintf("CaBus-%d", time.Now().UnixNano())
		caArchive := fmt.Sprintf("CaArchive-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(caBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(caBus)})

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, caBus)
		resp, err := client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(caArchive),
			EventSourceArn: aws.String(busARN),
			Description:    aws.String("test archive"),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ArchiveArn == nil || *resp.ArchiveArn == "" {
			return fmt.Errorf("archive ARN is nil or empty")
		}
		defer client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{ArchiveName: aws.String(caArchive)})
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeArchive", func() error {
		daBus := fmt.Sprintf("DaBus-%d", time.Now().UnixNano())
		daArchive := fmt.Sprintf("DaArchive-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(daBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(daBus)})

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, daBus)
		_, err = client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(daArchive),
			EventSourceArn: aws.String(busARN),
			Description:    aws.String("test archive for describe"),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}
		defer client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{ArchiveName: aws.String(daArchive)})

		resp, err := client.DescribeArchive(ctx, &eventbridge.DescribeArchiveInput{
			ArchiveName: aws.String(daArchive),
		})
		if err != nil {
			return fmt.Errorf("describe archive: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ArchiveName == nil || *resp.ArchiveName != daArchive {
			return fmt.Errorf("archive name mismatch, got %v", resp.ArchiveName)
		}
		if resp.EventSourceArn == nil || *resp.EventSourceArn != busARN {
			return fmt.Errorf("event source ARN mismatch, got %v", resp.EventSourceArn)
		}
		if resp.Description == nil || *resp.Description != "test archive for describe" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeArchive_NonExistent", func() error {
		_, err := client.DescribeArchive(ctx, &eventbridge.DescribeArchiveInput{
			ArchiveName: aws.String("nonexistent-archive-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent archive")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteArchive_NonExistent", func() error {
		_, err := client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{
			ArchiveName: aws.String("nonexistent-archive-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent archive")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteArchive", func() error {
		dlaBus := fmt.Sprintf("DlaBus-%d", time.Now().UnixNano())
		dlaArchive := fmt.Sprintf("DlaArchive-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(dlaBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(dlaBus)})

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, dlaBus)
		_, err = client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(dlaArchive),
			EventSourceArn: aws.String(busARN),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}

		resp, err := client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{
			ArchiveName: aws.String(dlaArchive),
		})
		if err != nil {
			return fmt.Errorf("delete archive: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}

		_, err = client.DescribeArchive(ctx, &eventbridge.DescribeArchiveInput{
			ArchiveName: aws.String(dlaArchive),
		})
		if err == nil {
			return fmt.Errorf("expected error describing deleted archive")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListArchives", func() error {
		laBus := fmt.Sprintf("LaBus-%d", time.Now().UnixNano())
		laArchive := fmt.Sprintf("LaArchive-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(laBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(laBus)})

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, laBus)
		_, err = client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(laArchive),
			EventSourceArn: aws.String(busARN),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}
		defer client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{ArchiveName: aws.String(laArchive)})

		resp, err := client.ListArchives(ctx, &eventbridge.ListArchivesInput{})
		if err != nil {
			return fmt.Errorf("list archives: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Archives == nil {
			return fmt.Errorf("archives list is nil")
		}
		found := false
		for _, a := range resp.Archives {
			if a.ArchiveName != nil && *a.ArchiveName == laArchive {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected archive %s in list", laArchive)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UpdateArchive", func() error {
		uaBus := fmt.Sprintf("UaBus-%d", time.Now().UnixNano())
		uaArchive := fmt.Sprintf("UaArchive-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(uaBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(uaBus)})

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, uaBus)
		_, err = client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(uaArchive),
			EventSourceArn: aws.String(busARN),
			Description:    aws.String("original description"),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}
		defer client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{ArchiveName: aws.String(uaArchive)})

		resp, err := client.UpdateArchive(ctx, &eventbridge.UpdateArchiveInput{
			ArchiveName: aws.String(uaArchive),
			Description: aws.String("updated description"),
		})
		if err != nil {
			return fmt.Errorf("update archive: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ArchiveArn == nil || *resp.ArchiveArn == "" {
			return fmt.Errorf("archive ARN is nil or empty")
		}
		desc, err := client.DescribeArchive(ctx, &eventbridge.DescribeArchiveInput{
			ArchiveName: aws.String(uaArchive),
		})
		if err != nil {
			return fmt.Errorf("describe after update: %v", err)
		}
		if desc.Description == nil || *desc.Description != "updated description" {
			return fmt.Errorf("description not updated, got %v", desc.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "StartReplay_DescribeReplay", func() error {
		srBus := fmt.Sprintf("SrBus-%d", time.Now().UnixNano())
		srArchive := fmt.Sprintf("SrArchive-%d", time.Now().UnixNano())
		srReplay := fmt.Sprintf("SrReplay-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(srBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, srBus)
		archiveARN := fmt.Sprintf("arn:aws:events:%s:000000000000:archive/%s", r.region, srArchive)
		_, err = client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(srArchive),
			EventSourceArn: aws.String(busARN),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}

		now := time.Now()
		startResp, err := client.StartReplay(ctx, &eventbridge.StartReplayInput{
			ReplayName:     aws.String(srReplay),
			EventSourceArn: aws.String(archiveARN),
			Destination: &types.ReplayDestination{
				Arn: aws.String(busARN),
			},
			EventStartTime: aws.Time(now.Add(-1 * time.Hour)),
			EventEndTime:   aws.Time(now),
		})
		if err != nil {
			return fmt.Errorf("start replay: %v", err)
		}
		if startResp == nil {
			return fmt.Errorf("start replay response is nil")
		}
		if startResp.ReplayArn == nil || *startResp.ReplayArn == "" {
			return fmt.Errorf("replay ARN is nil or empty")
		}

		describeResp, err := client.DescribeReplay(ctx, &eventbridge.DescribeReplayInput{
			ReplayName: aws.String(srReplay),
		})
		if err != nil {
			return fmt.Errorf("describe replay: %v", err)
		}
		if describeResp == nil {
			return fmt.Errorf("describe replay response is nil")
		}
		if describeResp.ReplayName == nil || *describeResp.ReplayName != srReplay {
			return fmt.Errorf("replay name mismatch, got %v", describeResp.ReplayName)
		}
		if describeResp.EventSourceArn == nil || *describeResp.EventSourceArn != archiveARN {
			return fmt.Errorf("event source ARN mismatch, got %v", describeResp.EventSourceArn)
		}
		client.CancelReplay(ctx, &eventbridge.CancelReplayInput{ReplayName: aws.String(srReplay)})
		client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{ArchiveName: aws.String(srArchive)})
		client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(srBus)})
		return nil
	}))

	results = append(results, r.RunTest("events", "CancelReplay_NonExistent", func() error {
		_, err := client.CancelReplay(ctx, &eventbridge.CancelReplayInput{
			ReplayName: aws.String("nonexistent-replay-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent replay")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListReplays", func() error {
		lrBus := fmt.Sprintf("LrBus-%d", time.Now().UnixNano())
		lrArchive := fmt.Sprintf("LrArchive-%d", time.Now().UnixNano())
		lrReplay := fmt.Sprintf("LrReplay-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(lrBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, lrBus)
		archiveARN := fmt.Sprintf("arn:aws:events:%s:000000000000:archive/%s", r.region, lrArchive)
		_, err = client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(lrArchive),
			EventSourceArn: aws.String(busARN),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}

		now := time.Now()
		_, err = client.StartReplay(ctx, &eventbridge.StartReplayInput{
			ReplayName:     aws.String(lrReplay),
			EventSourceArn: aws.String(archiveARN),
			Destination: &types.ReplayDestination{
				Arn: aws.String(busARN),
			},
			EventStartTime: aws.Time(now.Add(-1 * time.Hour)),
			EventEndTime:   aws.Time(now),
		})
		if err != nil {
			return fmt.Errorf("start replay: %v", err)
		}

		resp, err := client.ListReplays(ctx, &eventbridge.ListReplaysInput{})
		if err != nil {
			return fmt.Errorf("list replays: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Replays == nil {
			return fmt.Errorf("replays list is nil")
		}
		found := false
		for _, rp := range resp.Replays {
			if rp.ReplayName != nil && *rp.ReplayName == lrReplay {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected replay %s in list", lrReplay)
		}
		client.CancelReplay(ctx, &eventbridge.CancelReplayInput{ReplayName: aws.String(lrReplay)})
		client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{ArchiveName: aws.String(lrArchive)})
		client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(lrBus)})
		return nil
	}))

	return results
}
