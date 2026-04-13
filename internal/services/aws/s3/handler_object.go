package s3

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"vorpalstacks/internal/common/request"
)

func (h *S3Handler) handleObjectRequest(ctx *request.RequestContext, r *http.Request, bucket, key string) (interface{}, http.Header, int, error) {
	stores, err := h.svc.store(ctx)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	action := determineObjectAction(r)
	if err := h.checkAccess(ctx, r, stores, action, bucket, key); err != nil {
		return nil, nil, http.StatusForbidden, err
	}

	if r.Method == "POST" && r.URL.Query().Has("select") {
		return h.handleSelectObjectContent(ctx, r, bucket, key)
	}

	return h.objectOps.HandleRequest(r.Context(), ctx, r, bucket, key)
}

func (h *S3Handler) handleSelectObjectContent(ctx *request.RequestContext, r *http.Request, bucket, key string) (interface{}, http.Header, int, error) {
	var input SelectObjectContentInput
	if err := request.NewSafeXMLDecoder(r.Body).Decode(&input); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("failed to decode request: %w", err)
	}
	input.Bucket = bucket
	input.Key = key

	store, err := h.svc.store(ctx)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	objReader, obj, err := store.objects.Get(ctx, bucket, key)
	if err != nil {
		return nil, nil, http.StatusNotFound, err
	}
	defer objReader.Close()

	var dataReader io.Reader = objReader

	if obj.SSEMetadata != nil {
		encryptedData, encObj, err := store.objects.GetEncrypted(ctx, bucket, key, "")
		if err != nil {
			return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to get encrypted object: %w", err)
		}

		if encObj.SSEMetadata.EncryptionType == "CUSTOMER" {
			if input.SSECustomerKey == "" {
				return nil, nil, http.StatusBadRequest, fmt.Errorf("customer key is required for SSE-C encrypted object")
			}
			customerKey, err := h.svc.encryptionManager.ParseCustomerKey(input.SSECustomerKey, input.SSECustomerKeyMD5)
			if err != nil {
				return nil, nil, http.StatusBadRequest, fmt.Errorf("invalid SSE-C customer key: %w", err)
			}
			decResult, err := h.svc.encryptionManager.DecryptWithCustomerKey(encryptedData, encObj.SSEMetadata, bucket, key, customerKey)
			if err != nil {
				return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to decrypt data: %w", err)
			}
			dataReader = bytes.NewReader(decResult.DecryptedData)
		} else {
			decResult, err := h.svc.encryptionManager.Decrypt(encryptedData, encObj.SSEMetadata, bucket, key)
			if err != nil {
				return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to decrypt data: %w", err)
			}
			dataReader = bytes.NewReader(decResult.DecryptedData)
		}
	}

	engine, err := NewSelectEngine(&input)
	if err != nil {
		return nil, nil, http.StatusBadRequest, err
	}

	pr, pw := io.Pipe()
	header := make(http.Header)
	header.Set("Content-Type", "application/vnd.amazon.eventstream")
	header.Set("x-amz-request-id", fmt.Sprintf("%016X", time.Now().UnixNano()))

	go func() {
		defer pw.Close()
		writer := NewSelectEventStreamWriter(pw, input.RequestProgress)
		if err := engine.Execute(ctx, dataReader, writer); err != nil {
			pw.CloseWithError(err)
			return
		}
		if err := writer.WriteStats(); err != nil {
			pw.CloseWithError(err)
			return
		}
		if err := writer.WriteEnd(); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	output := &SelectObjectContentOutput{
		Payload: pr,
	}

	return output, header, http.StatusOK, nil
}

func (h *S3Handler) handleDeleteObjects(ctx *request.RequestContext, r *http.Request, bucket string, body io.Reader) (interface{}, int, error) {
	stores, err := h.svc.store(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	action := "s3:DeleteObject"
	if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
		return nil, http.StatusForbidden, err
	}

	result, err := h.objectOps.HandleDeleteObjects(r.Context(), ctx, bucket, body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return result, http.StatusOK, nil
}

func determineObjectAction(r *http.Request) string {
	method := r.Method
	query := r.URL.Query()

	switch {
	case method == "POST" && query.Has("select"):
		return "s3:GetObject"
	case method == "GET" && query.Has("acl"):
		return "s3:GetObjectAcl"
	case method == "PUT" && query.Has("acl"):
		return "s3:PutObjectAcl"
	case method == "GET" && query.Has("tagging"):
		return "s3:GetObjectTagging"
	case method == "PUT" && query.Has("tagging"):
		return "s3:PutObjectTagging"
	case method == "DELETE" && query.Has("tagging"):
		return "s3:DeleteObjectTagging"
	case method == "GET" && query.Has("legal-hold"):
		return "s3:GetObjectLegalHold"
	case method == "PUT" && query.Has("legal-hold"):
		return "s3:PutObjectLegalHold"
	case method == "GET" && query.Has("retention"):
		return "s3:GetObjectRetention"
	case method == "PUT" && query.Has("retention"):
		return "s3:PutObjectRetention"
	case method == "GET" && query.Has("uploadId"):
		return "s3:ListMultipartUploadParts"
	case method == "PUT" && query.Has("uploadId"):
		return "s3:UploadPart"
	case method == "POST" && query.Has("uploadId"):
		return "s3:CompleteMultipartUpload"
	case method == "DELETE" && query.Has("uploadId"):
		return "s3:AbortMultipartUpload"
	case method == "POST" && query.Has("uploads"):
		return "s3:CreateMultipartUpload"
	case method == "POST" && query.Has("restore"):
		return "s3:RestoreObject"
	case method == "PUT" && r.Header.Get("x-amz-copy-source") != "" && !query.Has("uploadId"):
		return "s3:PutObject"
	case method == "GET":
		return "s3:GetObject"
	case method == "HEAD":
		return "s3:GetObject"
	case method == "PUT":
		return "s3:PutObject"
	case method == "DELETE":
		return "s3:DeleteObject"
	default:
		return "s3:PutObject"
	}
}
