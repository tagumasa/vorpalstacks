package request

import (
	"encoding/xml"
	"io"
	"strings"
)

// xmlNode represents a node in an XML document tree.
// It is used for parsing AWS API request XML bodies.
type xmlNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Children []xmlNode  `xml:",any"`
	Content  string     `xml:",chardata"`
}

// NewSafeXMLDecoder creates a new XML decoder that is configured for non-strict parsing.
// It disables strict mode to allow parsing of loosely-formatted AWS API XML responses.
func NewSafeXMLDecoder(r io.Reader) *xml.Decoder {
	dec := xml.NewDecoder(r)
	dec.Strict = false
	return dec
}

func parseXMLBody(bodyBytes []byte, params map[string]interface{}) {
	var node xmlNode
	if err := xml.Unmarshal(bodyBytes, &node); err != nil {
		return
	}
	extractXMLParams(&node, params)
}

func extractXMLParams(node *xmlNode, params map[string]interface{}) {
	counts := make(map[string]int)
	for _, child := range node.Children {
		name := child.XMLName.Local
		counts[name]++
	}

	for _, child := range node.Children {
		name := child.XMLName.Local
		if counts[name] > 1 {
			if _, ok := params[name]; !ok {
				params[name] = []interface{}{}
			}
			arr := params[name].([]interface{})
			if len(child.Children) == 0 {
				content := strings.TrimSpace(child.Content)
				if content != "" {
					params[name] = append(arr, content)
				}
			} else {
				subParams := make(map[string]interface{})
				extractXMLParams(&child, subParams)
				if len(subParams) > 0 {
					params[name] = append(arr, subParams)
				}
			}
		} else {
			if len(child.Children) == 0 {
				content := strings.TrimSpace(child.Content)
				if content != "" {
					params[name] = content
				}
			} else if hasOnlyTextChildren(&child) {
				params[name] = strings.TrimSpace(child.Children[0].Content)
			} else {
				subParams := make(map[string]interface{})
				extractXMLParams(&child, subParams)
				if len(subParams) > 0 {
					params[name] = subParams
				}
			}
		}
	}
}

func hasOnlyTextChildren(node *xmlNode) bool {
	if len(node.Children) == 1 && node.Children[0].XMLName.Local == "" && strings.TrimSpace(node.Children[0].Content) != "" {
		return true
	}
	return false
}
