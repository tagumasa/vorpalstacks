package testutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

func (r *TestRunner) s3MultipartTests(ctx context.Context, client *s3.Client, ts string) []TestResult {
	var results []TestResult

	mpuBucket := s3Bucket(ts, "mpu")
	var uploadID *string
	var part1ETag *string
	var part2ETag *string

	results = append(results, r.RunTest("s3", "CreateMultipartUpload_Bucket", func() error {
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(mpuBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CreateMultipartUpload_Initiate", func() error {
		resp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(mpuBucket),
			Key:    aws.String("multipart-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload failed: %w", err)
		}
		if resp.UploadId == nil {
			return fmt.Errorf("UploadId is nil")
		}
		uploadID = resp.UploadId
		return nil
	}))

	results = append(results, r.RunTest("s3", "UploadPart_Part1", func() error {
		resp, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(mpuBucket),
			Key:        aws.String("multipart-obj.txt"),
			UploadId:   uploadID,
			PartNumber: aws.Int32(1),
			Body:       strings.NewReader("part one content"),
		})
		if err != nil {
			return fmt.Errorf("UploadPart 1 failed: %w", err)
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag for part 1 is nil")
		}
		part1ETag = resp.ETag
		return nil
	}))

	results = append(results, r.RunTest("s3", "UploadPart_Part2", func() error {
		resp, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(mpuBucket),
			Key:        aws.String("multipart-obj.txt"),
			UploadId:   uploadID,
			PartNumber: aws.Int32(2),
			Body:       strings.NewReader("part two content"),
		})
		if err != nil {
			return fmt.Errorf("UploadPart 2 failed: %w", err)
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag for part 2 is nil")
		}
		part2ETag = resp.ETag
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListParts_Verify", func() error {
		resp, err := client.ListParts(ctx, &s3.ListPartsInput{
			Bucket:   aws.String(mpuBucket),
			Key:      aws.String("multipart-obj.txt"),
			UploadId: uploadID,
		})
		if err != nil {
			return fmt.Errorf("ListParts failed: %w", err)
		}
		if len(resp.Parts) != 2 {
			return fmt.Errorf("expected 2 parts, got %d", len(resp.Parts))
		}
		for _, p := range resp.Parts {
			if p.PartNumber == nil {
				return fmt.Errorf("PartNumber is nil for a part")
			}
			if p.ETag == nil || *p.ETag == "" {
				return fmt.Errorf("ETag is nil or empty for part %d", aws.ToInt32(p.PartNumber))
			}
			if p.Size == nil || *p.Size == 0 {
				return fmt.Errorf("Size is nil or zero for part %d", aws.ToInt32(p.PartNumber))
			}
		}
		if resp.Key == nil || *resp.Key != "multipart-obj.txt" {
			return fmt.Errorf("expected Key multipart-obj.txt, got %v", resp.Key)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CompleteMultipartUpload_Verify", func() error {
		resp, err := client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(mpuBucket),
			Key:      aws.String("multipart-obj.txt"),
			UploadId: uploadID,
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: []types.CompletedPart{
					{
						ETag:       part1ETag,
						PartNumber: aws.Int32(1),
					},
					{
						ETag:       part2ETag,
						PartNumber: aws.Int32(2),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("CompleteMultipartUpload failed: %w", err)
		}
		if resp.Location == nil {
			return fmt.Errorf("Location is nil")
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		if resp.Key == nil || *resp.Key != "multipart-obj.txt" {
			return fmt.Errorf("expected Key multipart-obj.txt, got %v", resp.Key)
		}
		if resp.Bucket == nil || *resp.Bucket != mpuBucket {
			return fmt.Errorf("expected Bucket %s, got %v", mpuBucket, resp.Bucket)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultipartUpload_GetObject_VerifyContent", func() error {
		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(mpuBucket),
			Key:    aws.String("multipart-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading body failed: %w", err)
		}
		expected := "part one contentpart two content"
		if string(body) != expected {
			return fmt.Errorf("expected body %q, got %q", expected, string(body))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "AbortMultipartUpload_Verify", func() error {
		abortBucket := s3Bucket(ts, "mpu-abort")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(abortBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, abortBucket)

		initResp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(abortBucket),
			Key:    aws.String("abort-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload failed: %w", err)
		}

		_, err = client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(abortBucket),
			Key:        aws.String("abort-obj.txt"),
			UploadId:   initResp.UploadId,
			PartNumber: aws.Int32(1),
			Body:       strings.NewReader("will be aborted"),
		})
		if err != nil {
			return fmt.Errorf("UploadPart failed: %w", err)
		}

		_, err = client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(abortBucket),
			Key:      aws.String("abort-obj.txt"),
			UploadId: initResp.UploadId,
		})
		if err != nil {
			return fmt.Errorf("AbortMultipartUpload failed: %w", err)
		}

		listResp, err := client.ListMultipartUploads(ctx, &s3.ListMultipartUploadsInput{
			Bucket: aws.String(abortBucket),
		})
		if err != nil {
			return fmt.Errorf("ListMultipartUploads failed: %w", err)
		}
		for _, u := range listResp.Uploads {
			if u.UploadId != nil && *u.UploadId == *initResp.UploadId {
				return fmt.Errorf("aborted upload still listed in ListMultipartUploads")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListParts_NoSuchUpload", func() error {
		_, err := client.ListParts(ctx, &s3.ListPartsInput{
			Bucket:   aws.String(mpuBucket),
			Key:      aws.String("nonexistent-key.txt"),
			UploadId: aws.String("nonexistent-upload-id-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent upload ID, got nil")
		}
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() != "NoSuchUpload" {
				return fmt.Errorf("expected NoSuchUpload, got %s: %v", apiErr.ErrorCode(), err)
			}
		} else if !strings.Contains(err.Error(), "NoSuchUpload") {
			return fmt.Errorf("expected NoSuchUpload error, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListMultipartUploads_Verify", func() error {
		listBucket := s3Bucket(ts, "mpu-list")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(listBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, listBucket)

		_, err = client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(listBucket),
			Key:    aws.String("list-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload failed: %w", err)
		}

		resp, err := client.ListMultipartUploads(ctx, &s3.ListMultipartUploadsInput{
			Bucket: aws.String(listBucket),
		})
		if err != nil {
			return fmt.Errorf("ListMultipartUploads failed: %w", err)
		}
		if len(resp.Uploads) == 0 {
			return fmt.Errorf("expected at least 1 upload, got 0")
		}
		u := resp.Uploads[0]
		if u.Key == nil || *u.Key != "list-obj.txt" {
			return fmt.Errorf("expected Key list-obj.txt, got %v", u.Key)
		}
		if u.UploadId == nil || *u.UploadId == "" {
			return fmt.Errorf("UploadId is nil or empty")
		}
		if u.Initiated == nil {
			return fmt.Errorf("Initiated is nil")
		}
		return nil
	}))

	defer s3CleanupBucket(ctx, client, mpuBucket)

	return results
}
