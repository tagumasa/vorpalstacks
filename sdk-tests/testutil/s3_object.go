package testutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"strings"
)

func (r *TestRunner) s3ObjectTests(ctx context.Context, client *s3.Client, ts string, bucketName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("s3", "PutObject", func() error {
		resp, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
			Body:   strings.NewReader("Hello, World!"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("ETag is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetObject", func() error {
		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		if string(body) != "Hello, World!" {
			return fmt.Errorf("expected body %q, got %q", "Hello, World!", string(body))
		}
		if resp.ContentLength == nil || *resp.ContentLength != 13 {
			return fmt.Errorf("expected ContentLength 13, got %v", resp.ContentLength)
		}
		if resp.ContentType == nil || *resp.ContentType == "" {
			return fmt.Errorf("ContentType is nil or empty")
		}
		if resp.LastModified == nil {
			return fmt.Errorf("LastModified is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadObject", func() error {
		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if resp.ContentLength == nil || *resp.ContentLength != 13 {
			return fmt.Errorf("expected ContentLength 13, got %v", resp.ContentLength)
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		if resp.ContentType == nil || *resp.ContentType == "" {
			return fmt.Errorf("ContentType is nil or empty")
		}
		if resp.LastModified == nil {
			return fmt.Errorf("LastModified is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectsV2", func() error {
		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 failed: %w", err)
		}
		if resp.Contents == nil {
			return fmt.Errorf("Contents is nil")
		}
		found := false
		for _, obj := range resp.Contents {
			if obj.Key != nil && *obj.Key == "test.txt" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("test.txt not found in listing")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObject", func() error {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
		})
		if err != nil {
			return fmt.Errorf("DeleteObject failed: %w", err)
		}

		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
		})
		if err == nil {
			return fmt.Errorf("expected NoSuchKey error after delete, got nil")
		}
		var noSuchKey *types.NoSuchKey
		if !errors.As(err, &noSuchKey) {
			return fmt.Errorf("expected NoSuchKey, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectsAfterDelete", func() error {
		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 failed: %w", err)
		}
		if len(resp.Contents) != 0 {
			return fmt.Errorf("expected 0 objects, got %d", len(resp.Contents))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CopyObject", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("copy-source.txt"),
			Body:   strings.NewReader("copy me"),
		})
		if err != nil {
			return fmt.Errorf("PutObject source failed: %w", err)
		}

		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String("copy-dest.txt"),
			CopySource: aws.String(bucketName + "/copy-source.txt"),
		})
		if err != nil {
			return fmt.Errorf("CopyObject failed: %w", err)
		}

		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("copy-dest.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObject dest failed: %w", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		if string(body) != "copy me" {
			return fmt.Errorf("expected body %q, got %q", "copy me", string(body))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObject_GetObject_ContentVerification", func() error {
		content := "verification content"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String(bucketName),
			Key:           aws.String("verify.txt"),
			Body:          strings.NewReader(content),
			ContentType:   aws.String("text/plain"),
			ContentLength: aws.Int64(int64(len(content))),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("verify.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		if string(body) != content {
			return fmt.Errorf("expected body %q, got %q", content, string(body))
		}
		if resp.ContentLength == nil || *resp.ContentLength != int64(len(content)) {
			return fmt.Errorf("expected ContentLength %d, got %v", len(content), resp.ContentLength)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObject_Overwrite", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("overwrite.txt"),
			Body:   strings.NewReader("first"),
		})
		if err != nil {
			return fmt.Errorf("PutObject first failed: %w", err)
		}

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("overwrite.txt"),
			Body:   strings.NewReader("second"),
		})
		if err != nil {
			return fmt.Errorf("PutObject second failed: %w", err)
		}

		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("overwrite.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		if string(body) != "second" {
			return fmt.Errorf("expected body %q, got %q", "second", string(body))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadObject_VerifyMetadata", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String("metadata.txt"),
			Body:        strings.NewReader("{}"),
			ContentType: aws.String("application/json"),
			Metadata: map[string]string{
				"custom-key": "custom-value",
			},
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("metadata.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if resp.ContentType == nil || *resp.ContentType != "application/json" {
			return fmt.Errorf("expected ContentType application/json, got %v", resp.ContentType)
		}
		val, ok := resp.Metadata["custom-key"]
		if !ok || val != "custom-value" {
			return fmt.Errorf("expected Metadata[custom-key]=custom-value, got %q, ok=%v", val, ok)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectsV2_MultipleObjects", func() error {
		listBucket := s3Bucket(ts, "list")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(listBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, listBucket)

		for i := 0; i < 5; i++ {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(listBucket),
				Key:    aws.String(fmt.Sprintf("obj-%d.txt", i)),
				Body:   strings.NewReader(fmt.Sprintf("data-%d", i)),
			})
			if err != nil {
				return fmt.Errorf("PutObject %d failed: %w", i, err)
			}
		}

		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(listBucket),
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 failed: %w", err)
		}
		if len(resp.Contents) != 5 {
			return fmt.Errorf("expected 5 objects, got %d", len(resp.Contents))
		}
		for _, obj := range resp.Contents {
			if obj.Size == nil || *obj.Size == 0 {
				return fmt.Errorf("object %s has zero size", aws.ToString(obj.Key))
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectsV2_Pagination", func() error {
		pagBucket := s3Bucket(ts, "pag")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(pagBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, pagBucket)

		for i := 0; i < 5; i++ {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(pagBucket),
				Key:    aws.String(fmt.Sprintf("pag/obj-%d.txt", i)),
				Body:   strings.NewReader(fmt.Sprintf("page-data-%d", i)),
			})
			if err != nil {
				return fmt.Errorf("PutObject %d failed: %w", i, err)
			}
		}

		var allKeys []string
		var continuationToken *string
		pageCount := 0
		for {
			pageCount++
			input := &s3.ListObjectsV2Input{
				Bucket:            aws.String(pagBucket),
				Prefix:            aws.String("pag/"),
				MaxKeys:           aws.Int32(2),
				ContinuationToken: continuationToken,
			}
			resp, err := client.ListObjectsV2(ctx, input)
			if err != nil {
				return fmt.Errorf("ListObjectsV2 page %d failed: %w", pageCount, err)
			}
			for _, obj := range resp.Contents {
				if obj.Key != nil {
					allKeys = append(allKeys, *obj.Key)
				}
			}
			if resp.IsTruncated == nil || !*resp.IsTruncated {
				break
			}
			continuationToken = resp.NextContinuationToken
		}
		if len(allKeys) != 5 {
			return fmt.Errorf("expected 5 total keys, got %d", len(allKeys))
		}
		if pageCount < 2 {
			return fmt.Errorf("expected at least 2 pages, got %d", pageCount)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectVersions", func() error {
		verBucket := s3Bucket(ts, "versions")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(verBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, verBucket)

		_, err = client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(verBucket),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusEnabled,
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketVersioning failed: %w", err)
		}

		for i := 0; i < 3; i++ {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(verBucket),
				Key:    aws.String("versioned-key.txt"),
				Body:   strings.NewReader(fmt.Sprintf("version-%d", i)),
			})
			if err != nil {
				return fmt.Errorf("PutObject version %d failed: %w", i, err)
			}
		}

		resp, err := client.ListObjectVersions(ctx, &s3.ListObjectVersionsInput{
			Bucket: aws.String(verBucket),
		})
		if err != nil {
			return fmt.Errorf("ListObjectVersions failed: %w", err)
		}
		if len(resp.Versions) != 3 {
			return fmt.Errorf("expected 3 versions, got %d", len(resp.Versions))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetObject_NonExistentKey", func() error {
		_, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("nonexistent-key.txt"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		var noSuchKey *types.NoSuchKey
		if !errors.As(err, &noSuchKey) {
			return fmt.Errorf("expected NoSuchKey, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadObject_NonExistentKey", func() error {
		_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("nonexistent-key.txt"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		var notFound *types.NotFound
		if !errors.As(err, &notFound) {
			return fmt.Errorf("expected NotFound, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObject_NonExistentKey", func() error {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("nonexistent-delete-key.txt"),
		})
		if err != nil {
			return fmt.Errorf("DeleteObject on non-existent key should not error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObject_SystemMetadata", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:             aws.String(bucketName),
			Key:                aws.String("sysmeta.txt"),
			Body:               strings.NewReader("system metadata test"),
			ContentEncoding:    aws.String("gzip"),
			CacheControl:       aws.String("max-age=3600"),
			ContentDisposition: aws.String("attachment; filename=\"test.txt\""),
			ContentLanguage:    aws.String("en-US"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		headResp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("sysmeta.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if headResp.ContentEncoding == nil || *headResp.ContentEncoding != "gzip" {
			return fmt.Errorf("expected ContentEncoding gzip, got %v", headResp.ContentEncoding)
		}
		if headResp.CacheControl == nil || *headResp.CacheControl != "max-age=3600" {
			return fmt.Errorf("expected CacheControl max-age=3600, got %v", headResp.CacheControl)
		}
		if headResp.ContentDisposition == nil || *headResp.ContentDisposition != "attachment; filename=\"test.txt\"" {
			return fmt.Errorf("expected ContentDisposition attachment; filename=\"test.txt\", got %v", headResp.ContentDisposition)
		}
		if headResp.ContentLanguage == nil || *headResp.ContentLanguage != "en-US" {
			return fmt.Errorf("expected ContentLanguage en-US, got %v", headResp.ContentLanguage)
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("sysmeta.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		defer getResp.Body.Close()
		if getResp.ContentEncoding == nil || *getResp.ContentEncoding != "gzip" {
			return fmt.Errorf("GetObject: expected ContentEncoding gzip, got %v", getResp.ContentEncoding)
		}
		if getResp.CacheControl == nil || *getResp.CacheControl != "max-age=3600" {
			return fmt.Errorf("GetObject: expected CacheControl max-age=3600, got %v", getResp.CacheControl)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "SelectObjectContent", func() error {
		csvData := "name,age\nAlice,30\nBob,25\n"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("select-data.csv"),
			Body:   strings.NewReader(csvData),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := client.SelectObjectContent(ctx, &s3.SelectObjectContentInput{
			Bucket:         aws.String(bucketName),
			Key:            aws.String("select-data.csv"),
			Expression:     aws.String("SELECT * FROM s3object s WHERE s._2 = '30'"),
			ExpressionType: types.ExpressionTypeSql,
			InputSerialization: &types.InputSerialization{
				CSV: &types.CSVInput{
					FileHeaderInfo:  types.FileHeaderInfoUse,
					RecordDelimiter: aws.String("\n"),
					FieldDelimiter:  aws.String(","),
				},
			},
			OutputSerialization: &types.OutputSerialization{
				CSV: &types.CSVOutput{},
			},
		})
		if err != nil {
			return fmt.Errorf("SelectObjectContent failed: %w", err)
		}
		defer resp.GetStream().Close()

		var records []byte
		for event := range resp.GetStream().Events() {
			switch v := event.(type) {
			case *types.SelectObjectContentEventStreamMemberRecords:
				records = append(records, v.Value.Payload...)
			}
		}

		if len(records) == 0 {
			return fmt.Errorf("expected at least one record, got none")
		}
		if !strings.Contains(string(records), "Alice") {
			return fmt.Errorf("expected records to contain %q, got %q", "Alice", string(records))
		}
		return nil
	}))

	return results
}
