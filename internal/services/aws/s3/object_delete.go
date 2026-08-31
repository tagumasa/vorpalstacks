package s3

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// DeleteObjectInput contains the input parameters for the DeleteObject operation.
type DeleteObjectInput struct {
	Bucket                    string
	Key                       string
	VersionId                 string
	BypassGovernanceRetention bool
}
type DeleteObjectOutput struct {
	DeleteMarker   bool
	VersionId      string
	RequestCharged string
}

// DeleteObject deletes an object from S3.
func (o *ObjectOperations) DeleteObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *DeleteObjectInput) (*DeleteObjectOutput, error) {
	return o.svc.deleteObjectOpCore(ctx, reqCtx, stores, input)
}

// DeleteObjectsInput contains the input parameters for the DeleteObjects operation.
type DeleteObjectsInput struct {
	Bucket                    string
	Delete                    *Delete
	BypassGovernanceRetention bool
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
func (o *ObjectOperations) DeleteObjects(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *DeleteObjectsInput) (*DeleteObjectsOutput, error) {
	return o.svc.deleteObjectsOpCore(ctx, reqCtx, stores, input)
}
