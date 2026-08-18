package s3

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	s3store "vorpalstacks/internal/store/aws/s3"
)

func (h *S3Handler) handleObjectRequest(ctx *request.RequestContext, r *http.Request, bucket, key string) (interface{}, http.Header, int, error) {
	stores, err := h.svc.store(ctx)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	if r.Method == "POST" && r.URL.Query().Has("select") {
		return h.handleSelectObjectContent(ctx, r, bucket, key, stores)
	}

	if _, actions := classifyObjectRequest(r, bucket, key); len(actions) > 0 {
		for _, a := range actions {
			if err := h.checkAccess(ctx, r, stores, a.Action, a.Bucket, a.Key); err != nil {
				return nil, nil, http.StatusForbidden, err
			}
		}
	}

	return h.objectOps.HandleRequest(r.Context(), ctx, stores, r, bucket, key)
}

func (h *S3Handler) handleSelectObjectContent(ctx *request.RequestContext, r *http.Request, bucket, key string, stores *s3Stores) (interface{}, http.Header, int, error) {
	var input SelectObjectContentInput
	if err := request.NewSafeXMLDecoder(r.Body).Decode(&input); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("failed to decode request: %w", err)
	}
	input.Bucket = bucket
	input.Key = key

	objReader, obj, err := stores.objects.Get(ctx, bucket, key)
	if err != nil {
		return nil, nil, http.StatusNotFound, err
	}
	defer objReader.Close()

	var dataReader io.Reader = objReader

	if obj.SSEMetadata != nil {
		encryptedData, err := io.ReadAll(objReader)
		if err != nil {
			return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to read encrypted data: %w", err)
		}
		objReader.Close()

		var customerKey []byte
		if obj.SSEMetadata.EncryptionType == s3store.SSETypeCustomer {
			if input.SSECustomerKey == "" {
				return nil, nil, http.StatusBadRequest, fmt.Errorf("customer key is required for SSE-C encrypted object")
			}
			customerKey, err = h.svc.encryptionManager.ParseCustomerKey(input.SSECustomerKey, input.SSECustomerKeyMD5)
			if err != nil {
				return nil, nil, http.StatusBadRequest, NewInvalidArgumentError(fmt.Sprintf("invalid SSE-C customer key: %v", err))
			}
		}

		var plainData []byte
		if len(obj.SSEMetadata.PartEncryptionInfos) > 0 {
			plainData, err = h.svc.encryptionManager.DecryptChunked(encryptedData, obj.SSEMetadata, bucket, key, customerKey)
		} else if customerKey != nil {
			decResult, decErr := h.svc.encryptionManager.DecryptWithCustomerKey(encryptedData, obj.SSEMetadata, bucket, key, customerKey)
			if decErr != nil {
				return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to decrypt data: %w", decErr)
			}
			plainData = decResult.DecryptedData
		} else {
			decResult, decErr := h.svc.encryptionManager.Decrypt(encryptedData, obj.SSEMetadata, bucket, key)
			if decErr != nil {
				return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to decrypt data: %w", decErr)
			}
			plainData = decResult.DecryptedData
		}
		if err != nil {
			return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to decrypt data: %w", err)
		}
		dataReader = bytes.NewReader(plainData)
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
		defer func() {
			if r := recover(); r != nil {
				logs.Error("s3: panic in SelectObjectContent goroutine", logs.Any("panic", r))
			}
		}()
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

	if err := h.checkAccess(ctx, r, stores, "s3:DeleteObject", bucket, ""); err != nil {
		return nil, http.StatusForbidden, err
	}

	var deleteReq Delete
	if err := request.NewSafeXMLDecoder(body).Decode(&deleteReq); err != nil {
		return nil, http.StatusBadRequest, err
	}

	for _, obj := range deleteReq.Objects {
		if err := h.checkAccess(ctx, r, stores, "s3:DeleteObject", bucket, obj.Key); err != nil {
			return nil, http.StatusForbidden, err
		}
	}

	result, err := h.objectOps.DeleteObjects(r.Context(), ctx, stores, &DeleteObjectsInput{
		Bucket:                    bucket,
		Delete:                    &deleteReq,
		BypassGovernanceRetention: r.Header.Get("x-amz-bypass-governance-retention") == "true",
	})
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return result, http.StatusOK, nil
}
