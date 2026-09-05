package s3

import (
	"encoding/xml"
	"strings"
	"testing"
)

// decodeMetricsConfiguration parses one PutBucketMetricsConfiguration XML
// body the same way the bucket-PUT handler does.
func decodeMetricsConfiguration(t *testing.T, body string) *MetricsConfigurationInput {
	t.Helper()
	var config MetricsConfigurationInput
	if err := xml.Unmarshal([]byte(body), &config); err != nil {
		t.Fatalf("decode %q: %v", body, err)
	}
	return &config
}

// TestPutBucketMetricsConfigurationFilterCardinality pins the filter
// cardinality the model's MetricsFilter union imposes: a present filter
// carries exactly one predicate, so empty elements and multi-predicate
// filters are MalformedXML. Both rejections fire before any store access,
// so a nil store is safe here.
func TestPutBucketMetricsConfigurationFilterCardinality(t *testing.T) {
	svc := &S3Service{}

	cases := []struct {
		name string
		body string
	}{
		{"empty filter element", "<MetricsConfiguration><Id>empty</Id><Filter></Filter></MetricsConfiguration>"},
		{"empty prefix element", "<MetricsConfiguration><Id>empty</Id><Filter><Prefix></Prefix></Filter></MetricsConfiguration>"},
		{"prefix and tag", "<MetricsConfiguration><Id>two</Id><Filter><Prefix>docs/</Prefix><Tag><Key>class</Key><Value>blue</Value></Tag></Filter></MetricsConfiguration>"},
		{"access point and prefix", "<MetricsConfiguration><Id>two</Id><Filter><AccessPointArn>arn:aws:s3:us-east-1:123456789012:accesspoint/ap</AccessPointArn><Prefix>docs/</Prefix></Filter></MetricsConfiguration>"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.putBucketMetricsConfigurationCore(nil, &PutBucketMetricsConfigurationInput{
				Bucket:               "bucket",
				Id:                   "filter-cardinality",
				MetricsConfiguration: decodeMetricsConfiguration(t, tc.body),
			})
			if err == nil {
				t.Fatalf("%s: expected MalformedXML, got nil", tc.name)
			}
			if !strings.Contains(err.Error(), "MalformedXML") {
				t.Fatalf("%s: expected MalformedXML, got %v", tc.name, err)
			}
		})
	}
}
