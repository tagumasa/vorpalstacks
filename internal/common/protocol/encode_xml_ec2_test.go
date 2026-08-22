package protocol

import (
	"net/http/httptest"
	"strings"
	"testing"
)

// TestEncodeEC2QueryXMLResponseDeterministic pins that the EC2 encoder
// never leaks Go map iteration order into the wire output: repeated
// encodes of the same response map must be byte-identical, with element
// keys emitted in sorted order.
func TestEncodeEC2QueryXMLResponseDeterministic(t *testing.T) {
	response := map[string]interface{}{
		"zebra": "last",
		"apple": "first",
		"mango": "middle",
		"beta":  "b",
		"gamma": "g",
		"Tokens": map[string]string{
			"zT": "zv",
			"aT": "av",
			"mT": "mv",
		},
	}

	var first string
	for i := 0; i < 50; i++ {
		w := httptest.NewRecorder()
		if err := EncodeEC2QueryXMLResponse(w, "Det", response); err != nil {
			t.Fatalf("encode failed: %v", err)
		}
		if i == 0 {
			first = w.Body.String()
			continue
		}
		if w.Body.String() != first {
			t.Fatal("EC2 XML output differs between encodes of the same response map")
		}
	}

	if strings.Index(first, "<apple>") < 0 || strings.Index(first, "<mango>") < 0 || strings.Index(first, "<zebra>") < 0 {
		t.Fatal("expected element markers missing")
	}
	if !(strings.Index(first, "<apple>") < strings.Index(first, "<mango>") &&
		strings.Index(first, "<mango>") < strings.Index(first, "<zebra>")) {
		t.Fatal("element keys must be emitted in sorted order")
	}
}

// TestEncodeEC2QueryXMLResponseItems verifies that the EC2 Query encoder
// wraps list members in <item> elements (the AWS SDK ec2query deserialiser
// ignores any other child name) and applies lowerCamelCase element names.
func TestEncodeEC2QueryXMLResponseItems(t *testing.T) {
	response := map[string]interface{}{
		"return": true,
		"SecurityGroupRuleSet": XMLElements{
			ElementName: "item",
			Items: []interface{}{
				map[string]interface{}{
					"SecurityGroupRuleId": "sgr-0123456789abcdef0",
					"IsEgress":            false,
					"GroupId":             "sg-0123456789abcdef0",
					"IpProtocol":          "tcp",
					"FromPort":            int64(80),
					"ToPort":              int64(80),
					"CidrIpv4":            "0.0.0.0/0",
				},
			},
		},
		"VpcSet": XMLElements{
			ElementName: "item",
			Items: []interface{}{
				map[string]interface{}{
					"VpcId":     "vpc-1",
					"CidrBlock": "10.0.0.0/16",
				},
			},
		},
	}

	w := httptest.NewRecorder()
	if err := EncodeEC2QueryXMLResponse(w, "AuthorizeSecurityGroupIngress", response); err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	xml := w.Body.String()

	for _, want := range []string{
		"<AuthorizeSecurityGroupIngressResponse>",
		"<return>true</return>",
		"<securityGroupRuleSet>",
		"<securityGroupRuleId>sgr-0123456789abcdef0</securityGroupRuleId>",
		"<isEgress>false</isEgress>",
		"<cidrIpv4>0.0.0.0/0</cidrIpv4>",
		"<vpcSet><item><cidrBlock>10.0.0.0/16</cidrBlock><vpcId>vpc-1</vpcId></item></vpcSet>",
		"<requestId>",
	} {
		if !strings.Contains(xml, want) {
			t.Errorf("XML output missing %q:\n%s", want, xml)
		}
	}
	if strings.Contains(xml, "<member>") {
		t.Errorf("EC2 XML must not contain <member> elements:\n%s", xml)
	}
}

// TestEncodeEC2QueryXMLResponseStructSlices verifies that raw struct slices
// returned by handlers (defence in depth for handlers not using
// XMLElements) are also item-wrapped.
func TestEncodeEC2QueryXMLResponseStructSlices(t *testing.T) {
	type inner struct {
		Key string
	}
	type outer struct {
		Name  string
		Items []inner
	}
	response := map[string]interface{}{
		"Outer": outer{Name: "o", Items: []inner{{Key: "k1"}, {Key: "k2"}}},
	}
	w := httptest.NewRecorder()
	if err := EncodeEC2QueryXMLResponse(w, "TestOp", response); err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	xml := w.Body.String()
	for _, want := range []string{"<outer>", "<name>o</name>", "<items><item><key>k1</key></item><item><key>k2</key></item></items>"} {
		if !strings.Contains(xml, want) {
			t.Errorf("XML output missing %q:\n%s", want, xml)
		}
	}
}
