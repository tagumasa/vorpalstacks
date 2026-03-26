package s3

import (
	"context"
	"io"

	"vorpalstacks/internal/services/aws/common/request"
)

// DeleteObjectInput contains the input parameters for the DeleteObject operation.
type DeleteObjectInput struct {
	Bucket    string
	Key       string
	VersionId string
}

// DeleteObjectOutput contains the output from the DeleteObject operation.
type DeleteObjectOutput struct {
	DeleteMarker   bool
	VersionId      string
	RequestCharged string
}

// DeleteObject deletes an object from S3.
func (o *ObjectOperations) DeleteObject(ctx context.Context, reqCtx *request.RequestContext, input *DeleteObjectInput) (*DeleteObjectOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	marker, err := store.objects.DeleteWithVersion(ctx, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, err
	}

	output := &DeleteObjectOutput{}
	if marker != nil {
		output.DeleteMarker = true
		output.VersionId = marker.VersionID
	}

	return output, nil
}

// DeleteObjectsInput contains the input parameters for the DeleteObjects operation.
type DeleteObjectsInput struct {
	Bucket string
	Delete *Delete
}

// Delete contains the objects to delete.
type Delete struct {
	Objects []ObjectIdentifier `xml:"Object"`
	Quiet   bool               `xml:"Quiet"`
}

// ObjectIdentifier identifies a specific object to delete.
type ObjectIdentifier struct {
	Key       string `xml:"Key"`
	VersionId string `xml:"VersionId,omitempty"`
}

// DeleteObjectsOutput contains the output from the DeleteObjects operation.
type DeleteObjectsOutput struct {
	Deleted []DeletedObject `xml:"Deleted"`
	Error   []DeleteError   `xml:"Error"`
}

// DeletedObject contains information about a deleted object.
type DeletedObject struct {
	Key            string `xml:"Key"`
	VersionId      string `xml:"VersionId,omitempty"`
	DeleteMarker   bool   `xml:"DeleteMarker,omitempty"`
	DeleteMarkerId string `xml:"DeleteMarkerVersionId,omitempty"`
}

// DeleteError contains information about a delete error.
type DeleteError struct {
	Key     string `xml:"Key"`
	Code    string `xml:"Code"`
	Message string `xml:"Message"`
}

// DeleteObjects deletes multiple objects from S3 in a single request.
func (o *ObjectOperations) DeleteObjects(ctx context.Context, reqCtx *request.RequestContext, input *DeleteObjectsInput) (*DeleteObjectsOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var deleted []DeletedObject
	var errors []DeleteError

	for _, obj := range input.Delete.Objects {
		marker, err := store.objects.DeleteWithVersion(ctx, input.Bucket, obj.Key, obj.VersionId)
		if err != nil {
			errors = append(errors, DeleteError{
				Key:     obj.Key,
				Code:    "InternalError",
				Message: err.Error(),
			})
		} else {
			deletedObj := DeletedObject{
				Key: obj.Key,
			}
			if obj.VersionId != "" {
				deletedObj.VersionId = obj.VersionId
			}
			if marker != nil {
				deletedObj.DeleteMarker = true
				deletedObj.DeleteMarkerId = marker.VersionID
			}
			deleted = append(deleted, deletedObj)
		}
	}

	return &DeleteObjectsOutput{
		Deleted: deleted,
		Error:   errors,
	}, nil
}

// HandleDeleteObjects handles the DeleteObjects API request.
func (o *ObjectOperations) HandleDeleteObjects(ctx context.Context, reqCtx *request.RequestContext, bucket string, body io.Reader) (*DeleteObjectsOutput, error) {
	var deleteReq Delete
	if err := request.NewSafeXMLDecoder(body).Decode(&deleteReq); err != nil {
		return nil, err
	}

	return o.DeleteObjects(ctx, reqCtx, &DeleteObjectsInput{
		Bucket: bucket,
		Delete: &deleteReq,
	})
}
