package testutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
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
		if uploadID == nil {
			return fmt.Errorf("upload prerequisite missing: CreateMultipartUpload_Bucket did not run or failed")
		}
		resp, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(mpuBucket),
			Key:        aws.String("multipart-obj.txt"),
			UploadId:   uploadID,
			PartNumber: aws.Int32(1),
			Body:       bytes.NewReader(make([]byte, 5*1024*1024)),
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
		if uploadID == nil {
			return fmt.Errorf("upload prerequisite missing: CreateMultipartUpload_Bucket did not run or failed")
		}
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
		if uploadID == nil {
			return fmt.Errorf("upload prerequisite missing: CreateMultipartUpload_Bucket did not run or failed")
		}
		parts, err := s3ListPartsAll(ctx, client, mpuBucket, "multipart-obj.txt", *uploadID)
		if err != nil {
			return fmt.Errorf("ListParts failed: %w", err)
		}
		if len(parts) != 2 {
			return fmt.Errorf("expected 2 parts, got %d", len(parts))
		}
		for _, p := range parts {
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
		// Response metadata comes from a single call; the part set above
		// is what needed pagination.
		resp, err := client.ListParts(ctx, &s3.ListPartsInput{
			Bucket:   aws.String(mpuBucket),
			Key:      aws.String("multipart-obj.txt"),
			UploadId: uploadID,
		})
		if err != nil {
			return fmt.Errorf("ListParts (metadata) failed: %w", err)
		}
		if resp.Key == nil || *resp.Key != "multipart-obj.txt" {
			return fmt.Errorf("expected Key multipart-obj.txt, got %v", resp.Key)
		}
		if resp.Initiator == nil || resp.Initiator.ID == nil || *resp.Initiator.ID == "" {
			return fmt.Errorf("expected a populated Initiator ID, got %v", resp.Initiator)
		}
		if resp.Owner == nil || resp.Owner.ID == nil || *resp.Owner.ID == "" {
			return fmt.Errorf("expected a populated Owner ID, got %v", resp.Owner)
		}
		if resp.StorageClass != types.StorageClassStandard {
			return fmt.Errorf("expected StorageClass STANDARD, got %s", resp.StorageClass)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CompleteMultipartUpload_Verify", func() error {
		if uploadID == nil {
			return fmt.Errorf("upload prerequisite missing: CreateMultipartUpload_Bucket did not run or failed")
		}
		if part1ETag == nil || part2ETag == nil {
			return fmt.Errorf("part prerequisites missing: an UploadPart step did not run or failed")
		}
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
		// Location is served from this platform's endpoint, not the public
		// AWS hostname, and addresses the completed object.
		if strings.Contains(*resp.Location, "s3.amazonaws.com") {
			return fmt.Errorf("Location must derive from the serving endpoint, got %q", *resp.Location)
		}
		if !strings.Contains(*resp.Location, mpuBucket) || !strings.Contains(*resp.Location, "multipart-obj.txt") {
			return fmt.Errorf("Location must address the completed object, got %q", *resp.Location)
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
		_, gotBody, err := s3GetAndRead(ctx, client, &s3.GetObjectInput{
			Bucket: aws.String(mpuBucket),
			Key:    aws.String("multipart-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		expected := strings.Repeat("\x00", 5*1024*1024) + "part two content"
		if gotBody != expected {
			return fmt.Errorf("expected body %q, got %q", expected, gotBody)
		}

		// The ObjectParts attribute applies to any multipart object, not
		// only SSE-chunked ones. Without additional checksums the individual
		// Part elements stay unemitted; the count is served.
		attrs, err := client.GetObjectAttributes(ctx, &s3.GetObjectAttributesInput{
			Bucket: aws.String(mpuBucket),
			Key:    aws.String("multipart-obj.txt"),
			ObjectAttributes: []types.ObjectAttributes{
				types.ObjectAttributesObjectParts,
			},
		})
		if err != nil {
			return fmt.Errorf("GetObjectAttributes ObjectParts failed: %w", err)
		}
		if attrs.ObjectParts == nil {
			return fmt.Errorf("expected ObjectParts for a plain multipart object, got nil")
		}
		if attrs.ObjectParts.TotalPartsCount == nil || *attrs.ObjectParts.TotalPartsCount != 2 {
			return fmt.Errorf("expected TotalPartsCount 2, got %v", attrs.ObjectParts.TotalPartsCount)
		}
		if len(attrs.ObjectParts.Parts) != 0 {
			return fmt.Errorf("expected no Part elements without additional checksums, got %d", len(attrs.ObjectParts.Parts))
		}
		// The requested attribute set governs the body: an unrequested
		// attribute must stay absent, and unrequested means nil here.
		if attrs.ETag != nil {
			return fmt.Errorf("ETag returned although not requested")
		}

		// max-parts travels as the x-amz-max-parts request header; the
		// echoed MaxParts proves the header reached the pagination.
		maxPartsAttrs, err := client.GetObjectAttributes(ctx, &s3.GetObjectAttributesInput{
			Bucket:   aws.String(mpuBucket),
			Key:      aws.String("multipart-obj.txt"),
			MaxParts: aws.Int32(1),
			ObjectAttributes: []types.ObjectAttributes{
				types.ObjectAttributesObjectParts,
			},
		})
		if err != nil {
			return fmt.Errorf("GetObjectAttributes with MaxParts failed: %w", err)
		}
		if maxPartsAttrs.ObjectParts == nil || maxPartsAttrs.ObjectParts.MaxParts == nil || *maxPartsAttrs.ObjectParts.MaxParts != 1 {
			return fmt.Errorf("expected echoed MaxParts 1, got %v", maxPartsAttrs.ObjectParts)
		}
		return nil
	}))

	// A partNumber read of a multipart-uploaded object serves that part's
	// bytes as a ranged GET, and reports the part count header.
	results = append(results, r.RunTest("s3", "GetObject_PartNumber_MultipartObject", func() error {
		resp, gotBody, err := s3GetAndRead(ctx, client, &s3.GetObjectInput{
			Bucket:     aws.String(mpuBucket),
			Key:        aws.String("multipart-obj.txt"),
			PartNumber: aws.Int32(2),
		})
		if err != nil {
			return fmt.Errorf("GetObject partNumber=2 failed: %w", err)
		}
		if gotBody != "part two content" {
			return fmt.Errorf("expected part two content, got %q", gotBody)
		}
		total := int64(5*1024*1024) + int64(len("part two content"))
		wantRange := fmt.Sprintf("bytes %d-%d/%d", 5*1024*1024, total-1, total)
		if resp.ContentRange == nil || *resp.ContentRange != wantRange {
			return fmt.Errorf("expected ContentRange %q, got %v", wantRange, resp.ContentRange)
		}
		if resp.PartsCount == nil || *resp.PartsCount != 2 {
			return fmt.Errorf("expected PartsCount 2, got %v", resp.PartsCount)
		}

		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket:     aws.String(mpuBucket),
			Key:        aws.String("multipart-obj.txt"),
			PartNumber: aws.Int32(99),
		})
		if err == nil {
			return fmt.Errorf("expected error for partNumber beyond the part count, got nil")
		}
		if code := awsHTTPStatus(err); code != http.StatusRequestedRangeNotSatisfiable {
			return fmt.Errorf("expected HTTP 416 for unsatisfiable partNumber, got %d: %v", code, err)
		}
		return nil
	}))

	// A plain object acts as a single implicit part: partNumber=1 serves
	// the whole object and any other part number is unsatisfiable.
	results = append(results, r.RunTest("s3", "GetObject_PartNumber_PlainObject", func() error {
		key := "partnumber-plain.txt"
		if _, err := s3PutObject(ctx, client, mpuBucket, key, "plain body"); err != nil {
			return err
		}

		_, gotBody, err := s3GetAndRead(ctx, client, &s3.GetObjectInput{
			Bucket:     aws.String(mpuBucket),
			Key:        aws.String(key),
			PartNumber: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("GetObject partNumber=1 failed: %w", err)
		}
		if gotBody != "plain body" {
			return fmt.Errorf("expected whole body, got %q", gotBody)
		}

		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket:     aws.String(mpuBucket),
			Key:        aws.String(key),
			PartNumber: aws.Int32(2),
		})
		if err == nil {
			return fmt.Errorf("expected error for partNumber=2 on a plain object, got nil")
		}
		if code := awsHTTPStatus(err); code != http.StatusRequestedRangeNotSatisfiable {
			return fmt.Errorf("expected HTTP 416 for partNumber=2 on a plain object, got %d: %v", code, err)
		}
		return nil
	}))

	// The ContentType and user metadata requested at CreateMultipartUpload
	// are system metadata of the completed object and must survive
	// completion — on versioned and unversioned buckets alike.
	results = append(results, r.RunTest("s3", "Multipart_VersionedTrailingSlashKey", func() error {
		// Keys ending "/" are legal S3 keys whose blob file path rewrites
		// the final segment. Versioned writes of file-tier-sized objects
		// under such keys must stay reachable through the versionId — both
		// the plain put and the multipart complete route.
		bucket := s3Bucket(ts, "mpu-slash")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, bucket)
		if err := s3EnableVersioning(ctx, client, bucket); err != nil {
			return err
		}

		body := bytes.Repeat([]byte("z"), 5*1024*1024)

		putResp, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("put-dir/"),
			Body:   bytes.NewReader(body),
		})
		if err != nil {
			return fmt.Errorf("PutObject(put-dir/) failed: %w", err)
		}
		if putResp.VersionId == nil || *putResp.VersionId == "" {
			return fmt.Errorf("PutObject(put-dir/) VersionId is nil or empty")
		}

		initResp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("mpu-dir/"),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload(mpu-dir/) failed: %w", err)
		}
		part, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String("mpu-dir/"),
			UploadId:   initResp.UploadId,
			PartNumber: aws.Int32(1),
			Body:       bytes.NewReader(body),
		})
		if err != nil {
			return fmt.Errorf("UploadPart(mpu-dir/) failed: %w", err)
		}
		completeResp, err := client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String("mpu-dir/"),
			UploadId: initResp.UploadId,
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: []types.CompletedPart{{ETag: part.ETag, PartNumber: aws.Int32(1)}},
			},
		})
		if err != nil {
			return fmt.Errorf("CompleteMultipartUpload(mpu-dir/) failed: %w", err)
		}
		if completeResp.VersionId == nil || *completeResp.VersionId == "" {
			return fmt.Errorf("CompleteMultipartUpload(mpu-dir/) VersionId is nil or empty")
		}

		for _, obj := range []struct {
			key       string
			versionId string
		}{
			{"put-dir/", *putResp.VersionId},
			{"mpu-dir/", *completeResp.VersionId},
		} {
			getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
				Bucket:    aws.String(bucket),
				Key:       aws.String(obj.key),
				VersionId: aws.String(obj.versionId),
			})
			if err != nil {
				return fmt.Errorf("GetObject(%q, version %s) failed: %w", obj.key, obj.versionId, err)
			}
			got, readErr := io.ReadAll(getResp.Body)
			getResp.Body.Close()
			if readErr != nil {
				return fmt.Errorf("GetObject(%q, version %s) body read failed: %w", obj.key, obj.versionId, readErr)
			}
			if !bytes.Equal(got, body) {
				return fmt.Errorf("GetObject(%q, version %s) body mismatch: got %d bytes", obj.key, obj.versionId, len(got))
			}
			if getResp.ContentLength == nil || *getResp.ContentLength != int64(len(body)) {
				return fmt.Errorf("GetObject(%q, version %s) ContentLength: got %v, want %d", obj.key, obj.versionId, getResp.ContentLength, len(body))
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "Multipart_ContentTypeAndMetadata", func() error {
		for _, versioned := range []bool{false, true} {
			name := s3Bucket(ts, "mpu-meta")
			if versioned {
				name = s3Bucket(ts, "mpu-meta-ver")
			}
			if err := s3CreateBucket(ctx, client, name); err != nil {
				return err
			}
			defer s3CleanupBucket(ctx, client, name)

			if versioned {
				if err := s3EnableVersioning(ctx, client, name); err != nil {
					return err
				}
			}

			key := "typed-obj.txt"
			initResp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
				Bucket:      aws.String(name),
				Key:         aws.String(key),
				ContentType: aws.String("text/plain; charset=utf-8"),
				Metadata:    map[string]string{"origin": "multipart-upload"},
			})
			if err != nil {
				return fmt.Errorf("CreateMultipartUpload (%s) failed: %w", name, err)
			}

			part1, err := client.UploadPart(ctx, &s3.UploadPartInput{
				Bucket:     aws.String(name),
				Key:        aws.String(key),
				UploadId:   initResp.UploadId,
				PartNumber: aws.Int32(1),
				Body:       bytes.NewReader(bytes.Repeat([]byte("m"), 5*1024*1024)),
			})
			if err != nil {
				return fmt.Errorf("UploadPart 1 (%s) failed: %w", name, err)
			}
			part2, err := client.UploadPart(ctx, &s3.UploadPartInput{
				Bucket:     aws.String(name),
				Key:        aws.String(key),
				UploadId:   initResp.UploadId,
				PartNumber: aws.Int32(2),
				Body:       strings.NewReader("typed tail"),
			})
			if err != nil {
				return fmt.Errorf("UploadPart 2 (%s) failed: %w", name, err)
			}

			if _, err := client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
				Bucket:   aws.String(name),
				Key:      aws.String(key),
				UploadId: initResp.UploadId,
				MultipartUpload: &types.CompletedMultipartUpload{
					Parts: []types.CompletedPart{
						{ETag: part1.ETag, PartNumber: aws.Int32(1)},
						{ETag: part2.ETag, PartNumber: aws.Int32(2)},
					},
				},
			}); err != nil {
				return fmt.Errorf("CompleteMultipartUpload (%s) failed: %w", name, err)
			}

			getResp, _, err := s3GetRead(ctx, client, name, key)
			if err != nil {
				return fmt.Errorf("GetObject (%s) failed: %w", name, err)
			}
			if getResp.ContentType == nil || *getResp.ContentType != "text/plain; charset=utf-8" {
				return fmt.Errorf("GetObject ContentType (%s): got %v", name, getResp.ContentType)
			}
			if getResp.Metadata["origin"] != "multipart-upload" {
				return fmt.Errorf("GetObject Metadata (%s): got %v", name, getResp.Metadata)
			}
			wantSize := int64(5*1024*1024) + int64(len("typed tail"))
			if getResp.ContentLength == nil || *getResp.ContentLength != wantSize {
				return fmt.Errorf("GetObject ContentLength (%s): got %v, want %d", name, getResp.ContentLength, wantSize)
			}

			headResp, err := s3HeadObject(ctx, client, name, key)
			if err != nil {
				return fmt.Errorf("HeadObject (%s) failed: %w", name, err)
			}
			if headResp.ContentType == nil || *headResp.ContentType != "text/plain; charset=utf-8" {
				return fmt.Errorf("HeadObject ContentType (%s): got %v", name, headResp.ContentType)
			}
			if headResp.Metadata["origin"] != "multipart-upload" {
				return fmt.Errorf("HeadObject Metadata (%s): got %v", name, headResp.Metadata)
			}
		}
		return nil
	}))

	// SSE-C multipart responses signal encryption only through the customer
	// headers: the x-amz-server-side-encryption value set has no SSE-C
	// member, so create and complete must leave it unset.
	results = append(results, r.RunTest("s3", "Multipart_SSECResponseHeaders", func() error {
		bucket := s3Bucket(ts, "mpu-ssec")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, bucket)

		encodedKey, encodedMD5 := ssecTestKey()

		createResp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("ssec-mpu.txt"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload SSE-C failed: %w", err)
		}
		if createResp.SSECustomerAlgorithm == nil || *createResp.SSECustomerAlgorithm != "AES256" {
			return fmt.Errorf("expected SSECustomerAlgorithm AES256, got %v", createResp.SSECustomerAlgorithm)
		}
		if createResp.ServerSideEncryption != "" {
			return fmt.Errorf("SSE-C create must not set ServerSideEncryption, got %s", createResp.ServerSideEncryption)
		}

		part1, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("ssec-mpu.txt"),
			UploadId:             createResp.UploadId,
			PartNumber:           aws.Int32(1),
			Body:                 strings.NewReader("ssec multipart body"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err != nil {
			return fmt.Errorf("UploadPart SSE-C failed: %w", err)
		}

		completeResp, err := client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String("ssec-mpu.txt"),
			UploadId: createResp.UploadId,
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: []types.CompletedPart{
					{ETag: part1.ETag, PartNumber: aws.Int32(1)},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("CompleteMultipartUpload SSE-C failed: %w", err)
		}
		if completeResp.ServerSideEncryption != "" {
			return fmt.Errorf("SSE-C complete must not set ServerSideEncryption, got %s", completeResp.ServerSideEncryption)
		}

		getResp, gotBody, err := s3GetAndRead(ctx, client, &s3.GetObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("ssec-mpu.txt"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err != nil {
			return fmt.Errorf("GetObject SSE-C failed: %w", err)
		}
		if gotBody != "ssec multipart body" {
			return fmt.Errorf("expected body %q, got %q", "ssec multipart body", gotBody)
		}
		if getResp.ServerSideEncryption != "" {
			return fmt.Errorf("SSE-C get must not set ServerSideEncryption, got %s", getResp.ServerSideEncryption)
		}
		// The stored blob is ciphertext; the reported size is the plaintext
		// length, matching an encrypted single-part object.
		if getResp.ContentLength == nil || *getResp.ContentLength != int64(len("ssec multipart body")) {
			return fmt.Errorf("GetObject ContentLength: got %v, want %d", getResp.ContentLength, len("ssec multipart body"))
		}
		headResp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("ssec-mpu.txt"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err != nil {
			return fmt.Errorf("HeadObject SSE-C failed: %w", err)
		}
		if headResp.ContentLength == nil || *headResp.ContentLength != int64(len("ssec multipart body")) {
			return fmt.Errorf("HeadObject ContentLength: got %v, want %d", headResp.ContentLength, len("ssec multipart body"))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "AbortMultipartUpload_Verify", func() error {
		abortBucket := s3Bucket(ts, "mpu-abort")
		if err := s3CreateBucket(ctx, client, abortBucket); err != nil {
			return err
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

		// Aborting an unknown upload ID is NoSuchUpload (404), not a server
		// error.
		_, err = client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(abortBucket),
			Key:      aws.String("abort-obj.txt"),
			UploadId: aws.String("nonexistent-upload-id-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for abort with unknown upload ID, got nil")
		}
		var abortErr smithy.APIError
		if !errors.As(err, &abortErr) {
			return fmt.Errorf("expected API error, got %T: %v", err, err)
		}
		if abortErr.ErrorCode() != "NoSuchUpload" {
			return fmt.Errorf("expected NoSuchUpload, got %s: %v", abortErr.ErrorCode(), err)
		}
		if code := awsHTTPStatus(err); code != http.StatusNotFound {
			return fmt.Errorf("expected HTTP 404 for NoSuchUpload, got %d: %v", code, err)
		}

		// A valid upload ID under the wrong key addresses no upload at that
		// address — the same NoSuchUpload — and must leave the real upload
		// alive.
		mismatchInit, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(abortBucket),
			Key:    aws.String("mismatch-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload (mismatch) failed: %w", err)
		}
		_, err = client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(abortBucket),
			Key:      aws.String("abort-obj.txt"),
			UploadId: mismatchInit.UploadId,
		})
		if err == nil {
			return fmt.Errorf("expected error for abort under the wrong key, got nil")
		}
		if !errors.As(err, &abortErr) {
			return fmt.Errorf("expected API error for mismatched abort, got %T: %v", err, err)
		}
		if abortErr.ErrorCode() != "NoSuchUpload" {
			return fmt.Errorf("expected NoSuchUpload for mismatched abort, got %s: %v", abortErr.ErrorCode(), err)
		}
		listAfterMismatch, err := client.ListMultipartUploads(ctx, &s3.ListMultipartUploadsInput{
			Bucket: aws.String(abortBucket),
		})
		if err != nil {
			return fmt.Errorf("ListMultipartUploads after mismatched abort failed: %w", err)
		}
		found := false
		for _, u := range listAfterMismatch.Uploads {
			if u.UploadId != nil && *u.UploadId == *mismatchInit.UploadId {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("mismatched-address abort destroyed the real upload")
		}

		// Clean up the surviving upload.
		if _, err := client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(abortBucket),
			Key:      aws.String("mismatch-obj.txt"),
			UploadId: mismatchInit.UploadId,
		}); err != nil {
			return fmt.Errorf("AbortMultipartUpload (cleanup) failed: %w", err)
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

		// A nonexistent upload ID is NoSuchUpload on every multipart
		// operation — UploadPart and Complete must not report it as a
		// missing object (NoSuchKey).
		expectNoSuchUpload := func(op string, call func() error) error {
			if err := call(); err == nil {
				return fmt.Errorf("%s: expected error for non-existent upload ID, got nil", op)
			}
			if !errors.As(err, &apiErr) {
				return fmt.Errorf("%s: expected API error, got %T: %v", op, err, err)
			}
			if apiErr.ErrorCode() != "NoSuchUpload" {
				return fmt.Errorf("%s: expected NoSuchUpload, got %s: %v", op, apiErr.ErrorCode(), err)
			}
			if code := awsHTTPStatus(err); code != http.StatusNotFound {
				return fmt.Errorf("%s: expected HTTP 404 for NoSuchUpload, got %d: %v", op, code, err)
			}
			return nil
		}

		if err := expectNoSuchUpload("UploadPart", func() error {
			_, err := client.UploadPart(ctx, &s3.UploadPartInput{
				Bucket:     aws.String(mpuBucket),
				Key:        aws.String("nonexistent-key.txt"),
				UploadId:   aws.String("nonexistent-upload-id-12345"),
				PartNumber: aws.Int32(1),
				Body:       strings.NewReader("x"),
			})
			return err
		}); err != nil {
			return err
		}

		return expectNoSuchUpload("CompleteMultipartUpload", func() error {
			_, err := client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
				Bucket:   aws.String(mpuBucket),
				Key:      aws.String("nonexistent-key.txt"),
				UploadId: aws.String("nonexistent-upload-id-12345"),
				MultipartUpload: &types.CompletedMultipartUpload{
					Parts: []types.CompletedPart{{PartNumber: aws.Int32(1), ETag: aws.String("d41d8cd98f00b204e9800998ecf8427e")}},
				},
			})
			return err
		})
	}))

	results = append(results, r.RunTest("s3", "ListMultipartUploads_Verify", func() error {
		listBucket := s3Bucket(ts, "mpu-list")
		if err := s3CreateBucket(ctx, client, listBucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, listBucket)

		_, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
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
		if u.Initiator == nil || u.Initiator.ID == nil || *u.Initiator.ID == "" {
			return fmt.Errorf("expected a populated Initiator ID, got %v", u.Initiator)
		}
		if u.Owner == nil || u.Owner.ID == nil || *u.Owner.ID == "" {
			return fmt.Errorf("expected a populated Owner ID, got %v", u.Owner)
		}
		if u.Initiated == nil {
			return fmt.Errorf("Initiated is nil")
		}
		return nil
	}))

	// A key marker whose own uploads are gone must not truncate the
	// listing: pagination continues with the keys after the marker.
	results = append(results, r.RunTest("s3", "ListMultipartUploads_MarkerSkipsDeletedKey", func() error {
		bucket := s3Bucket(ts, "mpu-marker")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, bucket)

		create := func(key string) (string, error) {
			resp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			})
			if err != nil {
				return "", fmt.Errorf("CreateMultipartUpload(%s) failed: %w", key, err)
			}
			return *resp.UploadId, nil
		}
		if _, err := create("a.txt"); err != nil {
			return err
		}
		bID, err := create("b.txt")
		if err != nil {
			return err
		}
		if _, err := create("c.txt"); err != nil {
			return err
		}

		if _, err := client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String("b.txt"),
			UploadId: aws.String(bID),
		}); err != nil {
			return fmt.Errorf("AbortMultipartUpload failed: %w", err)
		}

		resp, err := client.ListMultipartUploads(ctx, &s3.ListMultipartUploadsInput{
			Bucket:     aws.String(bucket),
			KeyMarker:  aws.String("b.txt"),
			MaxUploads: aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("ListMultipartUploads failed: %w", err)
		}
		if len(resp.Uploads) != 1 || resp.Uploads[0].Key == nil || *resp.Uploads[0].Key != "c.txt" {
			keys := make([]string, 0, len(resp.Uploads))
			for _, u := range resp.Uploads {
				keys = append(keys, aws.ToString(u.Key))
			}
			return fmt.Errorf("key-marker past deleted key: got %v, want exactly [c.txt]", keys)
		}
		return nil
	}))

	// A key+upload-id pair marker must include the equal-key uploads whose
	// upload ID sorts after the marker's, alongside later keys.
	results = append(results, r.RunTest("s3", "ListMultipartUploads_PairMarkerIncludesFollowingUploadID", func() error {
		bucket := s3Bucket(ts, "mpu-pair")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, bucket)

		p1, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("p.txt"),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload 1 failed: %w", err)
		}
		p2, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("p.txt"),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload 2 failed: %w", err)
		}

		// The upload IDs' lexicographic order is not under the test's
		// control; the smaller one is the marker, so the other must follow.
		marker, follower := *p1.UploadId, *p2.UploadId
		if follower < marker {
			marker, follower = follower, marker
		}

		resp, err := client.ListMultipartUploads(ctx, &s3.ListMultipartUploadsInput{
			Bucket:         aws.String(bucket),
			KeyMarker:      aws.String("p.txt"),
			UploadIdMarker: aws.String(marker),
			MaxUploads:     aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("ListMultipartUploads failed: %w", err)
		}
		if len(resp.Uploads) != 1 || resp.Uploads[0].UploadId == nil || *resp.Uploads[0].UploadId != follower {
			ids := make([]string, 0, len(resp.Uploads))
			for _, u := range resp.Uploads {
				ids = append(ids, aws.ToString(u.UploadId))
			}
			return fmt.Errorf("pair-marker listing: got %v, want exactly [%s]", ids, follower)
		}
		return nil
	}))

	// The shared bucket outlives the builder: in register-only mode the
	// builder returns before any test runs, and a defer here would fire at
	// registration time as a no-op, leaking the bucket and its completed
	// object. The service cleanup runs on the executed phase instead.
	r.RegisterServiceCleanup("s3", func() { s3CleanupBucket(ctx, client, mpuBucket) })

	return results
}
