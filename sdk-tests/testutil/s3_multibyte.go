package testutil

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (r *TestRunner) s3MultibyteTests(ctx context.Context, client *s3.Client, ts string, bucketName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("s3", "MultiByte_PutGetRoundtrip", func() error {
		keys := []string{
			"テスト/日本語.txt",
			"测试/简体中文.txt",
			"測試/繁體中文.txt",
		}
		bodies := []string{
			"日本語テスト本文",
			"简体中文测试正文",
			"繁體中文測試正文",
		}
		for i, key := range keys {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
				Body:   strings.NewReader(bodies[i]),
			})
			if err != nil {
				return fmt.Errorf("PutObject key=%q failed: %w", key, err)
			}
			resp, err := client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
			})
			if err != nil {
				return fmt.Errorf("GetObject key=%q failed: %w", key, err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("ReadAll key=%q failed: %w", key, err)
			}
			if string(body) != bodies[i] {
				return fmt.Errorf("key=%q: expected body %q, got %q", key, bodies[i], string(body))
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_MetadataRoundtrip", func() error {
		key := "メタデータ/テスト.txt"
		meta := map[string]string{
			"lang":   "日本語",
			"author": "テスト太郎",
		}
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:   aws.String(bucketName),
			Key:      aws.String(key),
			Body:     strings.NewReader("metadata test"),
			Metadata: meta,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		for k, v := range meta {
			got, ok := resp.Metadata[k]
			if !ok {
				return fmt.Errorf("metadata key %q not found", k)
			}
			if got != v {
				return fmt.Errorf("metadata key %q: expected %q, got %q", k, v, got)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_ListObjectsV2", func() error {
		prefixes := []struct {
			prefix string
			keys   []string
		}{
			{prefix: "日本語/", keys: []string{"日本語/ファイル1.txt", "日本語/ファイル2.txt"}},
			{prefix: "简体/", keys: []string{"简体/文件1.txt", "简体/文件2.txt"}},
			{prefix: "繁體/", keys: []string{"繁體/檔案1.txt", "繁體/檔案2.txt"}},
		}
		for _, p := range prefixes {
			for _, key := range p.keys {
				_, err := client.PutObject(ctx, &s3.PutObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(key),
					Body:   strings.NewReader("list test"),
				})
				if err != nil {
					return fmt.Errorf("PutObject key=%q failed: %w", key, err)
				}
			}
			resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket: aws.String(bucketName),
				Prefix: aws.String(p.prefix),
			})
			if err != nil {
				return fmt.Errorf("ListObjectsV2 prefix=%q failed: %w", p.prefix, err)
			}
			if len(resp.Contents) != len(p.keys) {
				return fmt.Errorf("prefix=%q: expected %d objects, got %d", p.prefix, len(p.keys), len(resp.Contents))
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_CopyObject", func() error {
		srcKey := "コピー/元ファイル.txt"
		dstKey := "コピー/先ファイル.txt"
		content := "コピーテスト本文"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(srcKey),
			Body:   strings.NewReader(content),
		})
		if err != nil {
			return fmt.Errorf("PutObject source failed: %w", err)
		}
		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(dstKey),
			CopySource: aws.String(bucketName + "/" + srcKey),
		})
		if err != nil {
			return fmt.Errorf("CopyObject failed: %w", err)
		}
		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(dstKey),
		})
		if err != nil {
			return fmt.Errorf("GetObject dest failed: %w", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		if string(body) != content {
			return fmt.Errorf("expected body %q, got %q", content, string(body))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_DeleteObjects", func() error {
		keys := []string{
			"削除/ファイル1.txt",
			"削除/ファイル2.txt",
			"削除/ファイル3.txt",
		}
		for _, key := range keys {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
				Body:   strings.NewReader("delete me"),
			})
			if err != nil {
				return fmt.Errorf("PutObject key=%q failed: %w", key, err)
			}
		}
		var objIds []types.ObjectIdentifier
		for _, key := range keys {
			objIds = append(objIds, types.ObjectIdentifier{Key: aws.String(key)})
		}
		delResp, err := client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &types.Delete{Objects: objIds},
		})
		if err != nil {
			return fmt.Errorf("DeleteObjects failed: %w", err)
		}
		if len(delResp.Deleted) != len(keys) {
			return fmt.Errorf("expected %d deleted, got %d", len(keys), len(delResp.Deleted))
		}
		listResp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
			Prefix: aws.String("削除/"),
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 after delete failed: %w", err)
		}
		if len(listResp.Contents) != 0 {
			return fmt.Errorf("expected 0 objects after delete, got %d", len(listResp.Contents))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_MultipartUpload", func() error {
		key := "マルチパート/テスト.txt"
		initResp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload failed: %w", err)
		}
		part1Body := "パート1の内容"
		part2Body := "パート2の内容"
		up1, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(key),
			UploadId:   initResp.UploadId,
			PartNumber: aws.Int32(1),
			Body:       strings.NewReader(part1Body),
		})
		if err != nil {
			return fmt.Errorf("UploadPart 1 failed: %w", err)
		}
		up2, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(key),
			UploadId:   initResp.UploadId,
			PartNumber: aws.Int32(2),
			Body:       strings.NewReader(part2Body),
		})
		if err != nil {
			return fmt.Errorf("UploadPart 2 failed: %w", err)
		}
		_, err = client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(bucketName),
			Key:      aws.String(key),
			UploadId: initResp.UploadId,
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: []types.CompletedPart{
					{ETag: up1.ETag, PartNumber: aws.Int32(1)},
					{ETag: up2.ETag, PartNumber: aws.Int32(2)},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("CompleteMultipartUpload failed: %w", err)
		}
		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		defer getResp.Body.Close()
		body, err := io.ReadAll(getResp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		expected := part1Body + part2Body
		if string(body) != expected {
			return fmt.Errorf("expected body %q, got %q", expected, string(body))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_ObjectTagging", func() error {
		key := "タグ/テスト.txt"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader("tag test"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		tags := []types.Tag{
			{Key: aws.String("環境"), Value: aws.String("テスト")},
			{Key: aws.String("说明"), Value: aws.String("简体标签")},
			{Key: aws.String("說明"), Value: aws.String("繁體標籤")},
		}
		_, err = client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Tagging: &types.Tagging{
				TagSet: tags,
			},
		})
		if err != nil {
			return fmt.Errorf("PutObjectTagging failed: %w", err)
		}
		tagResp, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("GetObjectTagging failed: %w", err)
		}
		if len(tagResp.TagSet) != len(tags) {
			return fmt.Errorf("expected %d tags, got %d", len(tags), len(tagResp.TagSet))
		}
		for _, expected := range tags {
			found := false
			for _, got := range tagResp.TagSet {
				if got.Key != nil && expected.Key != nil && *got.Key == *expected.Key {
					if got.Value != nil && expected.Value != nil && *got.Value == *expected.Value {
						found = true
						break
					}
				}
			}
			if !found {
				return fmt.Errorf("tag %s=%s not found", *expected.Key, *expected.Value)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_ContentTypeRoundtrip", func() error {
		cases := []struct {
			key         string
			contentType string
			body        string
		}{
			{key: "shift_jis/テスト.txt", contentType: "text/plain; charset=shift_jis", body: "シフトJISテスト"},
			{key: "gb2312/测试.txt", contentType: "text/plain; charset=gb2312", body: "GB2312测试"},
			{key: "big5/測試.txt", contentType: "text/plain; charset=big5", body: "Big5測試"},
		}
		for _, tc := range cases {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket:      aws.String(bucketName),
				Key:         aws.String(tc.key),
				Body:        strings.NewReader(tc.body),
				ContentType: aws.String(tc.contentType),
			})
			if err != nil {
				return fmt.Errorf("PutObject key=%q failed: %w", tc.key, err)
			}
			headResp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(tc.key),
			})
			if err != nil {
				return fmt.Errorf("HeadObject key=%q failed: %w", tc.key, err)
			}
			if headResp.ContentType == nil || *headResp.ContentType != tc.contentType {
				return fmt.Errorf("key=%q: expected ContentType %q, got %v", tc.key, tc.contentType, headResp.ContentType)
			}
			getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(tc.key),
			})
			if err != nil {
				return fmt.Errorf("GetObject key=%q failed: %w", tc.key, err)
			}
			defer getResp.Body.Close()
			body, err := io.ReadAll(getResp.Body)
			if err != nil {
				return fmt.Errorf("ReadAll key=%q failed: %w", tc.key, err)
			}
			if string(body) != tc.body {
				return fmt.Errorf("key=%q: expected body %q, got %q", tc.key, tc.body, string(body))
			}
		}
		return nil
	}))

	return results
}
