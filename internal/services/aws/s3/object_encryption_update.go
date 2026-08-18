package s3

import (
	"context"
	"io"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
)

// updateObjectEncryptionPayload mirrors the ObjectEncryption payload of the
// UpdateObjectEncryption request body:
// <ObjectEncryption><SSE-KMS><KMSKeyArn>...</KMSKeyArn><BucketKeyEnabled>...
// SSE-KMS is the only union variant. BucketKeyEnabled is read but has no
// behavioural effect because this platform wraps a per-object data key under
// the KMS key and has no bucket-key tier for the flag to alter.
type updateObjectEncryptionPayload struct {
	SSEKMS *sseKMSVariant `xml:"SSE-KMS"`
}

type sseKMSVariant struct {
	KMSKeyArn        string `xml:"KMSKeyArn"`
	BucketKeyEnabled *bool  `xml:"BucketKeyEnabled"`
}

// UpdateObjectEncryption changes the server-side encryption of an existing
// object to SSE-KMS under the requested KMS key
// (PUT /{Bucket}/{Key+}?encryption). Object data is re-encrypted in place;
// metadata, storage class, ETag, and versioning are preserved.
func (o *ObjectOperations) UpdateObjectEncryption(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, bucket, key, versionId string, body io.Reader) error {
	if err := o.validateBucketExists(stores, bucket); err != nil {
		return err
	}
	if err := validateObjectKey(key); err != nil {
		return err
	}

	var payload updateObjectEncryptionPayload
	if body != nil {
		if err := request.NewSafeXMLDecoder(body).Decode(&payload); err != nil {
			return NewInvalidRequestError("invalid UpdateObjectEncryption request body")
		}
	}
	kmsKeyArn := ""
	if payload.SSEKMS != nil {
		kmsKeyArn = payload.SSEKMS.KMSKeyArn
	}
	if kmsKeyArn == "" {
		return NewInvalidRequestError("Requests that modify an object's encryption type to SSE-KMS require an Amazon Web Services KMS key Amazon Resource Name (ARN). Modify the request to specify a KMS key ARN, and then try again.")
	}

	if err := o.svc.updateObjectEncryptionCore(ctx, stores.objects, UpdateObjectEncryptionInput{
		Bucket:    bucket,
		Key:       key,
		VersionID: versionId,
		KMSKeyArn: kmsKeyArn,
	}); err != nil {
		return err
	}

	logs.Info("s3: object encryption updated",
		logs.String("bucket", bucket),
		logs.String("key", key),
		logs.String("kms_key_arn", kmsKeyArn))
	return nil
}
