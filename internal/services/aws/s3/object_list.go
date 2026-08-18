package s3

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
)

// s3Encode applies URL percent-encoding when encodingType is "url", otherwise
// returns the value unchanged (already XML-escaped by the caller).
func s3Encode(value, encodingType string) string {
	if encodingType == "url" {
		return url.QueryEscape(value)
	}
	return value
}

// ListObjectsInput contains the input parameters for the ListObjects operation.
type ListObjectsInput struct {
	Bucket       string
	Delimiter    string
	Prefix       string
	MaxKeys      int
	Marker       string
	EncodingType string
}

// ListObjectsOutput contains the output from the ListObjects operation.
type ListObjectsOutput struct {
	Contents       []*ObjectContent
	CommonPrefixes []CommonPrefix
	Delimiter      string
	EncodingType   string
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
	enc := o.EncodingType
	result.WriteString(`<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">`)
	for _, c := range o.Contents {
		result.WriteString(`<Contents>`)
		result.WriteString(`<Key>`)
		result.WriteString(s3Encode(xmlEscape(c.Key), enc))
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
	writeCommonPrefixesXML(&result, o.CommonPrefixes, enc)
	result.WriteString(`<Delimiter>`)
	result.WriteString(s3Encode(xmlEscape(o.Delimiter), enc))
	result.WriteString(`</Delimiter><IsTruncated>`)
	result.WriteString(strconv.FormatBool(o.IsTruncated))
	result.WriteString(`</IsTruncated><Marker>`)
	result.WriteString(s3Encode(xmlEscape(o.Marker), enc))
	result.WriteString(`</Marker><MaxKeys>`)
	result.WriteString(strconv.Itoa(o.MaxKeys))
	result.WriteString(`</MaxKeys><Name>`)
	result.WriteString(xmlEscape(o.Name))
	result.WriteString(`</Name>`)
	// NextMarker is only returned for truncated responses that specified
	// a delimiter; otherwise clients paginate using the last Key value.
	if o.NextMarker != "" && o.IsTruncated && o.Delimiter != "" {
		result.WriteString(`<NextMarker>`)
		result.WriteString(s3Encode(xmlEscape(o.NextMarker), enc))
		result.WriteString(`</NextMarker>`)
	}
	result.WriteString(`<Prefix>`)
	result.WriteString(s3Encode(xmlEscape(o.Prefix), enc))
	result.WriteString(`</Prefix>`)
	if enc != "" {
		result.WriteString(`<EncodingType>`)
		result.WriteString(enc)
		result.WriteString(`</EncodingType>`)
	}
	result.WriteString(`</ListBucketResult>`)
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

func writeCommonPrefixesXML(builder *strings.Builder, prefixes []CommonPrefix, encodingType string) {
	for _, p := range prefixes {
		builder.WriteString(`<CommonPrefixes><Prefix>`)
		builder.WriteString(s3Encode(xmlEscape(p.Prefix), encodingType))
		builder.WriteString(`</Prefix></CommonPrefixes>`)
	}
}

// ListObjects lists the objects in a bucket.
func (o *ObjectOperations) ListObjects(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *ListObjectsInput) (*ListObjectsOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	coreResult, err := o.svc.listObjectsCore(stores.objects, AdminListObjectsInput{
		Bucket:    input.Bucket,
		Prefix:    input.Prefix,
		Delimiter: input.Delimiter,
		Marker:    input.Marker,
		MaxKeys:   input.MaxKeys,
	})
	if err != nil {
		return nil, err
	}

	var commonPrefixes []CommonPrefix
	for _, prefix := range coreResult.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, CommonPrefix{Prefix: prefix})
	}

	return &ListObjectsOutput{
		Contents:       buildObjectContents(coreResult.Objects),
		CommonPrefixes: commonPrefixes,
		Delimiter:      input.Delimiter,
		EncodingType:   input.EncodingType,
		IsTruncated:    coreResult.IsTruncated,
		Marker:         input.Marker,
		MaxKeys:        input.MaxKeys,
		Name:           input.Bucket,
		NextMarker:     coreResult.NextMarker,
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
	EncodingType      string
}

// ListObjectsV2Output contains the output from the ListObjectsV2 operation.
type ListObjectsV2Output struct {
	Contents              []*ObjectContent
	CommonPrefixes        []CommonPrefix
	Delimiter             string
	EncodingType          string
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
	enc := o.EncodingType
	result.WriteString(`<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">`)
	for _, c := range o.Contents {
		result.WriteString(`<Contents>`)
		result.WriteString(`<Key>`)
		result.WriteString(s3Encode(xmlEscape(c.Key), enc))
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
	writeCommonPrefixesXML(&result, o.CommonPrefixes, enc)
	result.WriteString(`<Delimiter>`)
	result.WriteString(s3Encode(xmlEscape(o.Delimiter), enc))
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
	result.WriteString(s3Encode(xmlEscape(o.Prefix), enc))
	result.WriteString(`</Prefix>`)
	if o.StartAfter != "" {
		result.WriteString(`<StartAfter>`)
		result.WriteString(s3Encode(xmlEscape(o.StartAfter), enc))
		result.WriteString(`</StartAfter>`)
	}
	if enc != "" {
		result.WriteString(`<EncodingType>`)
		result.WriteString(enc)
		result.WriteString(`</EncodingType>`)
	}
	result.WriteString(`</ListBucketResult>`)
	return result.String()
}

// ListObjectsV2 lists the objects in a bucket using version 2 of the API.
func (o *ObjectOperations) ListObjectsV2(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *ListObjectsV2Input) (*ListObjectsV2Output, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	marker := input.ContinuationToken
	if marker == "" {
		marker = input.StartAfter
	}

	coreResult, err := o.svc.listObjectsCore(stores.objects, AdminListObjectsInput{
		Bucket:    input.Bucket,
		Prefix:    input.Prefix,
		Delimiter: input.Delimiter,
		Marker:    marker,
		MaxKeys:   input.MaxKeys,
	})
	if err != nil {
		return nil, err
	}

	var commonPrefixes []CommonPrefix
	for _, prefix := range coreResult.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, CommonPrefix{Prefix: prefix})
	}

	contents := buildObjectContents(coreResult.Objects)
	output := &ListObjectsV2Output{
		Contents:       contents,
		CommonPrefixes: commonPrefixes,
		Delimiter:      input.Delimiter,
		EncodingType:   input.EncodingType,
		IsTruncated:    coreResult.IsTruncated,
		KeyCount:       len(contents) + len(commonPrefixes),
		MaxKeys:        input.MaxKeys,
		Name:           input.Bucket,
		Prefix:         input.Prefix,
	}

	if coreResult.IsTruncated && coreResult.NextMarker != "" {
		output.NextContinuationToken = coreResult.NextMarker
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
	EncodingType    string
}

// ListObjectVersionsOutput contains the output from the ListObjectVersions operation.
type ListObjectVersionsOutput struct {
	Versions            []*ObjectVersion
	DeleteMarkers       []*DeleteMarkerEntry
	CommonPrefixes      []CommonPrefix
	Delimiter           string
	EncodingType        string
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
	enc := o.EncodingType
	result.WriteString(`<ListVersionsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">`)
	result.WriteString(`<Name>`)
	result.WriteString(xmlEscape(o.Name))
	result.WriteString(`</Name><Prefix>`)
	result.WriteString(s3Encode(xmlEscape(o.Prefix), enc))
	result.WriteString(`</Prefix>`)
	if o.KeyMarker != "" {
		result.WriteString(`<KeyMarker>`)
		result.WriteString(s3Encode(xmlEscape(o.KeyMarker), enc))
		result.WriteString(`</KeyMarker>`)
	}
	if o.VersionIdMarker != "" {
		result.WriteString(`<VersionIdMarker>`)
		result.WriteString(xmlEscape(o.VersionIdMarker))
		result.WriteString(`</VersionIdMarker>`)
	}
	if o.NextKeyMarker != "" {
		result.WriteString(`<NextKeyMarker>`)
		result.WriteString(s3Encode(xmlEscape(o.NextKeyMarker), enc))
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
	result.WriteString(s3Encode(xmlEscape(o.Delimiter), enc))
	result.WriteString(`</Delimiter>`)

	for _, v := range o.Versions {
		result.WriteString(`<Version><Key>`)
		result.WriteString(s3Encode(xmlEscape(v.Key), enc))
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
		result.WriteString(s3Encode(xmlEscape(d.Key), enc))
		result.WriteString(`</Key><VersionId>`)
		result.WriteString(xmlEscape(d.VersionId))
		result.WriteString(`</VersionId><IsLatest>`)
		result.WriteString(strconv.FormatBool(d.IsLatest))
		result.WriteString(`</IsLatest><LastModified>`)
		result.WriteString(d.LastModified.Format(time.RFC3339))
		result.WriteString(`</LastModified></DeleteMarker>`)
	}

	writeCommonPrefixesXML(&result, o.CommonPrefixes, enc)

	if enc != "" {
		result.WriteString(`<EncodingType>`)
		result.WriteString(enc)
		result.WriteString(`</EncodingType>`)
	}
	result.WriteString(`</ListVersionsResult>`)
	return result.String()
}

// ListObjectVersions lists the versions of objects in a bucket.
func (o *ObjectOperations) ListObjectVersions(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *ListObjectVersionsInput) (*ListObjectVersionsOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	result, err := stores.objects.ListObjectVersions(input.Bucket, input.Prefix, input.Delimiter, input.KeyMarker, input.VersionIdMarker, input.MaxKeys)
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
		EncodingType:        input.EncodingType,
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
