package protocol

import (
	"net/http/httptest"
	"strings"
	"testing"
)

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
		"<vpcSet><item><vpcId>vpc-1</vpcId><cidrBlock>10.0.0.0/16</cidrBlock></item></vpcSet>",
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
