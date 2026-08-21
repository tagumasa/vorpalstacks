package bucketname

import "testing"

// TestValidate accepts the AWS general-purpose bucket naming examples
// and rejects every documented invalid form.
func TestValidate(t *testing.T) {
	valid := []string{
		"example",
		"myawsbucket",
		"example.com",
		"www.example.com",
		"my.example.s3.bucket",
		"doc-example-bucket1-a1b2c3d4-5678-90ab-cdef-example11111",
	}
	for _, name := range valid {
		if !Validate(name) {
			t.Errorf("Validate(%q) = false, want true", name)
		}
	}

	invalid := map[string]string{
		"ab":                         "too short",
		"example..com":               "two adjacent periods",
		"192.168.5.4":                "IP address form",
		"amzn_s3_demo_bucket":        "underscore",
		"AmznS3DemoBucket":           "uppercase",
		"bucket-":                    "trailing hyphen",
		"-bucket":                    "leading hyphen",
		".bucket":                    "leading period",
		"bucket.":                    "trailing period",
		"xn--bucket":                 "reserved prefix xn--",
		"sthree-bucket":              "reserved prefix sthree-",
		"amzn-s3-demo-bucket":        "reserved prefix amzn-s3-demo-",
		"bucket-s3alias":             "reserved suffix -s3alias",
		"bucket--ol-s3":              "reserved suffix --ol-s3",
		"bucket.mrap":                "reserved suffix .mrap",
		"bucket--x-s3":               "reserved suffix --x-s3",
		"bucket--table-s3":           "reserved suffix --table-s3",
		"a.-b":                       "dot-hyphen pair",
		"a-.b":                       "hyphen-dot pair",
		"example-" + repeat("a", 56): "longer than 63 characters",
	}
	for name, reason := range invalid {
		if Validate(name) {
			t.Errorf("Validate(%q) = true, want false (%s)", name, reason)
		}
	}
}

func repeat(s string, n int) string {
	out := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		out = append(out, s...)
	}
	return string(out)
}
