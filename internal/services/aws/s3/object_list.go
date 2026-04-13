package s3

import (
	"context"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
)

// ListObjectsInput contains the input parameters for the ListObjects operation.
type ListObjectsInput struct {
	Bucket    string
	Delimiter string
	Prefix    string
	MaxKeys   int
	Marker    string
}

// ListObjectsOutput contains the output from the ListObjects operation.
type ListObjectsOutput struct {
	Contents       []*ObjectContent
	CommonPrefixes []CommonPrefix
	Delimiter      string
	IsTruncated    bool
	Marker         string
	MaxKeys        int
	Name           string
	NextMarker     string
	Prefix         string
}

// ToXML converts the ListObjectsOutput to XML format.
func (o *ListObjectsOutput) ToXML() string {
	var result strings.Builder
	result.WriteString(`<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">`)
	for _, c := range o.Contents {
		result.WriteString(`<Contents>`)
		result.WriteString(`<Key>`)
		result.WriteString(xmlEscape(c.Key))
		result.WriteString(`</Key><LastModified>`)
		result.WriteString(c.LastModified.Format(time.RFC3339))
		result.WriteString(`</LastModified><ETag>`)
		result.WriteString(xmlEscape(c.ETag))
		result.WriteString(`</ETag><Size>`)
		result.WriteString(strconv.FormatInt(c.Size, 10))
		result.WriteString(`</Size><StorageClass>`)
		result.WriteString(c.StorageClass)
		result.WriteString(`</StorageClass></Contents>`)
	}
	for _, p := range o.CommonPrefixes {
		result.WriteString(`<CommonPrefixes><Prefix>`)
		result.WriteString(xmlEscape(p.Prefix))
		result.WriteString(`</Prefix></CommonPrefixes>`)
	}
	result.WriteString(`<Delimiter>`)
	result.WriteString(xmlEscape(o.Delimiter))
	result.WriteString(`</Delimiter><IsTruncated>`)
	result.WriteString(strconv.FormatBool(o.IsTruncated))
	result.WriteString(`</IsTruncated><Marker>`)
	result.WriteString(xmlEscape(o.Marker))
	result.WriteString(`</Marker><MaxKeys>`)
	result.WriteString(strconv.Itoa(o.MaxKeys))
	result.WriteString(`</MaxKeys><Name>`)
	result.WriteString(xmlEscape(o.Name))
	result.WriteString(`</Name>`)
	if o.NextMarker != "" {
		result.WriteString(`<NextMarker>`)
		result.WriteString(xmlEscape(o.NextMarker))
		result.WriteString(`</NextMarker>`)
	}
	result.WriteString(`<Prefix>`)
	result.WriteString(xmlEscape(o.Prefix))
	result.WriteString(`</Prefix></ListBucketResult>`)
	return result.String()
}

// ObjectContent contains information about an object in a list operation.
type ObjectContent struct {
	Key          string    `xml:"Key"`
	LastModified time.Time `xml:"LastModified"`
	ETag         string    `xml:"ETag"`
	Size         int64     `xml:"Size"`
	StorageClass string    `xml:"StorageClass"`
	Owner        *Owner    `xml:"Owner,omitempty"`
}

// CommonPrefix contains a prefix that represents a folder.
type CommonPrefix struct {
	Prefix string `xml:"Prefix"`
}

// ListObjects lists the objects in a bucket.
func (o *ObjectOperations) ListObjects(ctx context.Context, reqCtx *request.RequestContext, input *ListObjectsInput) (*ListObjectsOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	if input.MaxKeys <= 0 {
		input.MaxKeys = 1000
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.objects.List(input.Bucket, input.Prefix, input.Delimiter, input.Marker, input.MaxKeys)
	if err != nil {
		return nil, err
	}

	var commonPrefixes []CommonPrefix
	for _, prefix := range result.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, CommonPrefix{Prefix: prefix})
	}

	return &ListObjectsOutput{
		Contents:       buildObjectContents(result.Objects),
		CommonPrefixes: commonPrefixes,
		Delimiter:      input.Delimiter,
		IsTruncated:    result.IsTruncated,
		Marker:         input.Marker,
		MaxKeys:        input.MaxKeys,
		Name:           input.Bucket,
		NextMarker:     result.NextMarker,
		Prefix:         input.Prefix,
	}, nil
}

// ListObjectsV2Input contains the input parameters for the ListObjectsV2 operation.
type ListObjectsV2Input struct {
	Bucket            string
	Delimiter         string
	Prefix            string
	MaxKeys           int
	ContinuationToken string
	StartAfter        string
}

// ListObjectsV2Output contains the output from the ListObjectsV2 operation.
type ListObjectsV2Output struct {
	Contents              []*ObjectContent
	CommonPrefixes        []CommonPrefix
	Delimiter             string
	IsTruncated           bool
	KeyCount              int
	MaxKeys               int
	Name                  string
	NextContinuationToken string
	Prefix                string
	StartAfter            string
}

// ToXML converts the ListObjectsV2Output to XML format.
func (o *ListObjectsV2Output) ToXML() string {
	var result strings.Builder
	result.WriteString(`<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">`)
	for _, c := range o.Contents {
		result.WriteString(`<Contents>`)
		result.WriteString(`<Key>`)
		result.WriteString(xmlEscape(c.Key))
		result.WriteString(`</Key><LastModified>`)
		result.WriteString(c.LastModified.Format(time.RFC3339))
		result.WriteString(`</LastModified><ETag>`)
		result.WriteString(xmlEscape(c.ETag))
		result.WriteString(`</ETag><Size>`)
		result.WriteString(strconv.FormatInt(c.Size, 10))
		result.WriteString(`</Size><StorageClass>`)
		result.WriteString(c.StorageClass)
		result.WriteString(`</StorageClass></Contents>`)
	}
	for _, p := range o.CommonPrefixes {
		result.WriteString(`<CommonPrefixes><Prefix>`)
		result.WriteString(xmlEscape(p.Prefix))
		result.WriteString(`</Prefix></CommonPrefixes>`)
	}
	result.WriteString(`<Delimiter>`)
	result.WriteString(xmlEscape(o.Delimiter))
	result.WriteString(`</Delimiter><IsTruncated>`)
	result.WriteString(strconv.FormatBool(o.IsTruncated))
	result.WriteString(`</IsTruncated><KeyCount>`)
	result.WriteString(strconv.Itoa(o.KeyCount))
	result.WriteString(`</KeyCount><MaxKeys>`)
	result.WriteString(strconv.Itoa(o.MaxKeys))
	result.WriteString(`</MaxKeys><Name>`)
	result.WriteString(xmlEscape(o.Name))
	result.WriteString(`</Name>`)
	if o.NextContinuationToken != "" {
		result.WriteString(`<NextContinuationToken>`)
		result.WriteString(xmlEscape(o.NextContinuationToken))
		result.WriteString(`</NextContinuationToken>`)
	}
	result.WriteString(`<Prefix>`)
	result.WriteString(xmlEscape(o.Prefix))
	result.WriteString(`</Prefix>`)
	if o.StartAfter != "" {
		result.WriteString(`<StartAfter>`)
		result.WriteString(xmlEscape(o.StartAfter))
		result.WriteString(`</StartAfter>`)
	}
	result.WriteString(`</ListBucketResult>`)
	return result.String()
}

// ListObjectsV2 lists the objects in a bucket using version 2 of the API.
func (o *ObjectOperations) ListObjectsV2(ctx context.Context, reqCtx *request.RequestContext, input *ListObjectsV2Input) (*ListObjectsV2Output, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	marker := input.ContinuationToken
	if marker == "" {
		marker = input.StartAfter
	}

	if input.MaxKeys <= 0 {
		input.MaxKeys = 1000
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.objects.List(input.Bucket, input.Prefix, input.Delimiter, marker, input.MaxKeys)
	if err != nil {
		return nil, err
	}

	var commonPrefixes []CommonPrefix
	for _, prefix := range result.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, CommonPrefix{Prefix: prefix})
	}

	contents := buildObjectContents(result.Objects)
	output := &ListObjectsV2Output{
		Contents:       contents,
		CommonPrefixes: commonPrefixes,
		Delimiter:      input.Delimiter,
		IsTruncated:    result.IsTruncated,
		KeyCount:       len(contents) + len(commonPrefixes),
		MaxKeys:        input.MaxKeys,
		Name:           input.Bucket,
		Prefix:         input.Prefix,
	}

	if result.IsTruncated && result.NextMarker != "" {
		output.NextContinuationToken = result.NextMarker
	}

	if input.StartAfter != "" {
		output.StartAfter = input.StartAfter
	}

	return output, nil
}

// ListObjectVersionsInput contains the input parameters for the ListObjectVersions operation.
type ListObjectVersionsInput struct {
	Bucket          string
	Delimiter       string
	Prefix          string
	MaxKeys         int
	KeyMarker       string
	VersionIdMarker string
}

// ListObjectVersionsOutput contains the output from the ListObjectVersions operation.
type ListObjectVersionsOutput struct {
	Versions            []*ObjectVersion
	DeleteMarkers       []*DeleteMarkerEntry
	CommonPrefixes      []CommonPrefix
	Delimiter           string
	IsTruncated         bool
	KeyMarker           string
	MaxKeys             int
	Name                string
	NextKeyMarker       string
	NextVersionIdMarker string
	Prefix              string
	VersionIdMarker     string
}

// ObjectVersion contains information about a specific version of an object.
type ObjectVersion struct {
	Key          string    `xml:"Key"`
	LastModified time.Time `xml:"LastModified"`
	ETag         string    `xml:"ETag"`
	Size         int64     `xml:"Size"`
	StorageClass string    `xml:"StorageClass"`
	VersionId    string    `xml:"VersionId"`
	IsLatest     bool      `xml:"IsLatest"`
	Owner        *Owner    `xml:"Owner,omitempty"`
}

// DeleteMarkerEntry contains information about a delete marker.
type DeleteMarkerEntry struct {
	Key          string    `xml:"Key"`
	LastModified time.Time `xml:"LastModified"`
	VersionId    string    `xml:"VersionId"`
	IsLatest     bool      `xml:"IsLatest"`
	Owner        *Owner    `xml:"Owner,omitempty"`
}

// ToXML converts the ListObjectVersionsOutput to XML format.
func (o *ListObjectVersionsOutput) ToXML() string {
	var result strings.Builder
	result.WriteString(`<ListVersionsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">`)
	result.WriteString(`<Name>`)
	result.WriteString(xmlEscape(o.Name))
	result.WriteString(`</Name><Prefix>`)
	result.WriteString(xmlEscape(o.Prefix))
	result.WriteString(`</Prefix>`)
	if o.KeyMarker != "" {
		result.WriteString(`<KeyMarker>`)
		result.WriteString(xmlEscape(o.KeyMarker))
		result.WriteString(`</KeyMarker>`)
	}
	if o.VersionIdMarker != "" {
		result.WriteString(`<VersionIdMarker>`)
		result.WriteString(xmlEscape(o.VersionIdMarker))
		result.WriteString(`</VersionIdMarker>`)
	}
	if o.NextKeyMarker != "" {
		result.WriteString(`<NextKeyMarker>`)
		result.WriteString(xmlEscape(o.NextKeyMarker))
		result.WriteString(`</NextKeyMarker>`)
	}
	if o.NextVersionIdMarker != "" {
		result.WriteString(`<NextVersionIdMarker>`)
		result.WriteString(xmlEscape(o.NextVersionIdMarker))
		result.WriteString(`</NextVersionIdMarker>`)
	}
	result.WriteString(`<MaxKeys>`)
	result.WriteString(strconv.Itoa(o.MaxKeys))
	result.WriteString(`</MaxKeys><IsTruncated>`)
	result.WriteString(strconv.FormatBool(o.IsTruncated))
	result.WriteString(`</IsTruncated><Delimiter>`)
	result.WriteString(xmlEscape(o.Delimiter))
	result.WriteString(`</Delimiter>`)

	for _, v := range o.Versions {
		result.WriteString(`<Version><Key>`)
		result.WriteString(xmlEscape(v.Key))
		result.WriteString(`</Key><VersionId>`)
		result.WriteString(xmlEscape(v.VersionId))
		result.WriteString(`</VersionId><IsLatest>`)
		result.WriteString(strconv.FormatBool(v.IsLatest))
		result.WriteString(`</IsLatest><LastModified>`)
		result.WriteString(v.LastModified.Format(time.RFC3339))
		result.WriteString(`</LastModified><ETag>`)
		result.WriteString(xmlEscape(v.ETag))
		result.WriteString(`</ETag><Size>`)
		result.WriteString(strconv.FormatInt(v.Size, 10))
		result.WriteString(`</Size><StorageClass>`)
		result.WriteString(v.StorageClass)
		result.WriteString(`</StorageClass></Version>`)
	}

	for _, d := range o.DeleteMarkers {
		result.WriteString(`<DeleteMarker><Key>`)
		result.WriteString(xmlEscape(d.Key))
		result.WriteString(`</Key><VersionId>`)
		result.WriteString(xmlEscape(d.VersionId))
		result.WriteString(`</VersionId><IsLatest>`)
		result.WriteString(strconv.FormatBool(d.IsLatest))
		result.WriteString(`</IsLatest><LastModified>`)
		result.WriteString(d.LastModified.Format(time.RFC3339))
		result.WriteString(`</LastModified></DeleteMarker>`)
	}

	for _, p := range o.CommonPrefixes {
		result.WriteString(`<CommonPrefixes><Prefix>`)
		result.WriteString(xmlEscape(p.Prefix))
		result.WriteString(`</Prefix></CommonPrefixes>`)
	}

	result.WriteString(`</ListVersionsResult>`)
	return result.String()
}

// ListObjectVersions lists the versions of objects in a bucket.
func (o *ObjectOperations) ListObjectVersions(ctx context.Context, reqCtx *request.RequestContext, input *ListObjectVersionsInput) (*ListObjectVersionsOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	if input.MaxKeys <= 0 {
		input.MaxKeys = 1000
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.objects.ListObjectVersions(input.Bucket, input.Prefix, input.Delimiter, input.KeyMarker, input.VersionIdMarker, input.MaxKeys)
	if err != nil {
		return nil, err
	}

	var versions []*ObjectVersion
	var deleteMarkers []*DeleteMarkerEntry
	var commonPrefixes []CommonPrefix

	for _, obj := range result.Objects {
		if obj.IsDeleteMarker {
			deleteMarkers = append(deleteMarkers, &DeleteMarkerEntry{
				Key:          obj.Key,
				LastModified: obj.LastModified,
				VersionId:    obj.VersionID,
				IsLatest:     obj.IsLatest,
			})
		} else {
			versions = append(versions, &ObjectVersion{
				Key:          obj.Key,
				LastModified: obj.LastModified,
				ETag:         formatETag(obj.ETag),
				Size:         obj.Size,
				StorageClass: string(obj.StorageClass),
				VersionId:    obj.VersionID,
				IsLatest:     obj.IsLatest,
			})
		}
	}

	for _, prefix := range result.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, CommonPrefix{Prefix: prefix})
	}

	return &ListObjectVersionsOutput{
		Versions:            versions,
		DeleteMarkers:       deleteMarkers,
		CommonPrefixes:      commonPrefixes,
		Delimiter:           input.Delimiter,
		IsTruncated:         result.IsTruncated,
		KeyMarker:           input.KeyMarker,
		MaxKeys:             input.MaxKeys,
		Name:                input.Bucket,
		NextKeyMarker:       result.NextVersionKeyMarker,
		NextVersionIdMarker: result.NextVersionIDMarker,
		Prefix:              input.Prefix,
		VersionIdMarker:     input.VersionIdMarker,
	}, nil
}
