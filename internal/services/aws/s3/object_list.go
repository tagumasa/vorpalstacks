package s3

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/utils/timeutils"
)

// s3Encode renders a list-response string value under encodingType. With
// url encoding the raw value is percent-encoded first: the encoded form
// contains no XML-special characters, so the XML escape that follows is
// inert and a URL-decoding client recovers the original value. Without it
// the value is XML-escaped only.
func s3Encode(value, encodingType string) string {
	if encodingType == "url" {
		value = url.QueryEscape(value)
	}
	return xmlEscape(value)
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
	writeContentsXML(&result, o.Contents, enc)
	writeCommonPrefixesXML(&result, o.CommonPrefixes, enc)
	result.WriteString(`<Delimiter>`)
	result.WriteString(s3Encode(o.Delimiter, enc))
	result.WriteString(`</Delimiter><IsTruncated>`)
	result.WriteString(strconv.FormatBool(o.IsTruncated))
	result.WriteString(`</IsTruncated><Marker>`)
	result.WriteString(s3Encode(o.Marker, enc))
	result.WriteString(`</Marker><MaxKeys>`)
	result.WriteString(strconv.Itoa(o.MaxKeys))
	result.WriteString(`</MaxKeys><Name>`)
	result.WriteString(xmlEscape(o.Name))
	result.WriteString(`</Name>`)
	// NextMarker is only returned for truncated responses that specified
	// a delimiter; otherwise clients paginate using the last Key value.
	if o.NextMarker != "" && o.IsTruncated && o.Delimiter != "" {
		result.WriteString(`<NextMarker>`)
		result.WriteString(s3Encode(o.NextMarker, enc))
		result.WriteString(`</NextMarker>`)
	}
	result.WriteString(`<Prefix>`)
	result.WriteString(s3Encode(o.Prefix, enc))
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
		builder.WriteString(s3Encode(p.Prefix, encodingType))
		builder.WriteString(`</Prefix></CommonPrefixes>`)
	}
}

// writeContentsXML renders the <Contents> entries of a list result; the
// ListObjects and ListObjectsV2 models define the entry identically, so
// the rendering lives here once.
func writeContentsXML(builder *strings.Builder, contents []*ObjectContent, encodingType string) {
	for _, c := range contents {
		builder.WriteString(`<Contents>`)
		builder.WriteString(`<Key>`)
		builder.WriteString(s3Encode(c.Key, encodingType))
		builder.WriteString(`</Key><LastModified>`)
		builder.WriteString(c.LastModified.Format(timeutils.ISO8601UTCFormat))
		builder.WriteString(`</LastModified><ETag>`)
		builder.WriteString(xmlEscape(c.ETag))
		builder.WriteString(`</ETag><Size>`)
		builder.WriteString(strconv.FormatInt(c.Size, 10))
		builder.WriteString(`</Size><StorageClass>`)
		builder.WriteString(c.StorageClass)
		builder.WriteString(`</StorageClass>`)
		if c.Owner != nil {
			builder.WriteString(`<Owner><ID>`)
			builder.WriteString(xmlEscape(c.Owner.ID))
			builder.WriteString(`</ID><DisplayName>`)
			builder.WriteString(xmlEscape(c.Owner.DisplayName))
			builder.WriteString(`</DisplayName></Owner>`)
		}
		builder.WriteString(`</Contents>`)
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
		// V1 carries the owner on every Contents entry unconditionally;
		// the fetch-owner opt-in exists only in V2.
		Contents:       buildObjectContents(coreResult.Objects, o.svc.bucketOwner()),
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
	// FetchOwner surfaces the bucket owner on each entry (?fetch-owner=true);
	// V2 omits Owner unless requested.
	FetchOwner bool
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
	writeContentsXML(&result, o.Contents, enc)
	writeCommonPrefixesXML(&result, o.CommonPrefixes, enc)
	result.WriteString(`<Delimiter>`)
	result.WriteString(s3Encode(o.Delimiter, enc))
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
	result.WriteString(s3Encode(o.Prefix, enc))
	result.WriteString(`</Prefix>`)
	if o.StartAfter != "" {
		result.WriteString(`<StartAfter>`)
		result.WriteString(s3Encode(o.StartAfter, enc))
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

	// ?fetch-owner=true surfaces the bucket owner on every entry; V2
	// omits the element otherwise.
	var listOwner *Owner
	if input.FetchOwner {
		listOwner = o.svc.bucketOwner()
	}
	contents := buildObjectContents(coreResult.Objects, listOwner)
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
	result.WriteString(s3Encode(o.Prefix, enc))
	result.WriteString(`</Prefix>`)
	if o.KeyMarker != "" {
		result.WriteString(`<KeyMarker>`)
		result.WriteString(s3Encode(o.KeyMarker, enc))
		result.WriteString(`</KeyMarker>`)
	}
	if o.VersionIdMarker != "" {
		result.WriteString(`<VersionIdMarker>`)
		result.WriteString(xmlEscape(o.VersionIdMarker))
		result.WriteString(`</VersionIdMarker>`)
	}
	if o.NextKeyMarker != "" {
		result.WriteString(`<NextKeyMarker>`)
		result.WriteString(s3Encode(o.NextKeyMarker, enc))
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
	result.WriteString(s3Encode(o.Delimiter, enc))
	result.WriteString(`</Delimiter>`)

	for _, v := range o.Versions {
		result.WriteString(`<Version><Key>`)
		result.WriteString(s3Encode(v.Key, enc))
		result.WriteString(`</Key><VersionId>`)
		result.WriteString(xmlEscape(v.VersionId))
		result.WriteString(`</VersionId><IsLatest>`)
		result.WriteString(strconv.FormatBool(v.IsLatest))
		result.WriteString(`</IsLatest><LastModified>`)
		result.WriteString(v.LastModified.Format(timeutils.ISO8601UTCFormat))
		result.WriteString(`</LastModified><ETag>`)
		result.WriteString(xmlEscape(v.ETag))
		result.WriteString(`</ETag><Size>`)
		result.WriteString(strconv.FormatInt(v.Size, 10))
		result.WriteString(`</Size><StorageClass>`)
		result.WriteString(v.StorageClass)
		result.WriteString(`</StorageClass>`)
		if v.Owner != nil {
			result.WriteString(`<Owner><ID>`)
			result.WriteString(xmlEscape(v.Owner.ID))
			result.WriteString(`</ID><DisplayName>`)
			result.WriteString(xmlEscape(v.Owner.DisplayName))
			result.WriteString(`</DisplayName></Owner>`)
		}
		result.WriteString(`</Version>`)
	}

	for _, d := range o.DeleteMarkers {
		result.WriteString(`<DeleteMarker><Key>`)
		result.WriteString(s3Encode(d.Key, enc))
		result.WriteString(`</Key><VersionId>`)
		result.WriteString(xmlEscape(d.VersionId))
		result.WriteString(`</VersionId><IsLatest>`)
		result.WriteString(strconv.FormatBool(d.IsLatest))
		result.WriteString(`</IsLatest><LastModified>`)
		result.WriteString(d.LastModified.Format(timeutils.ISO8601UTCFormat))
		result.WriteString(`</LastModified>`)
		if d.Owner != nil {
			result.WriteString(`<Owner><ID>`)
			result.WriteString(xmlEscape(d.Owner.ID))
			result.WriteString(`</ID><DisplayName>`)
			result.WriteString(xmlEscape(d.Owner.DisplayName))
			result.WriteString(`</DisplayName></Owner>`)
		}
		result.WriteString(`</DeleteMarker>`)
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
	return o.svc.listObjectVersionsCore(stores, input)
}
