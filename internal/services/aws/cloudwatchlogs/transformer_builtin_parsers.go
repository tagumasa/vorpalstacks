package cloudwatchlogs

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// The built-in vended-format parsers: the five decoders the BuiltIn
// page documents with complete input→output examples (parseWAF,
// parsePostgres, parseCloudfront, parseRoute53, parseVPC). Each accepts
// only @message as its input and must be the transformer's first
// processor; each replaces the record with the documented output shape.
// parseToOCSF is not here: its AWS-to-OCSF field mapping is unpublished
// (the page documents the input contract alone), so it rejects at Put
// rather than being approximated.

// builtinSourceMessage reads the parser's input: the @message the five
// built-in parsers accept alone ("This processor accepts only @message
// as the input").
func builtinSourceMessage(record map[string]interface{}, cfg map[string]interface{}) (string, bool) {
	source := stringMember(cfg, "source", "@message")
	if source != "@message" {
		return "", false
	}
	raw, ok := record["@message"]
	if !ok {
		return "", false
	}
	text, ok := raw.(string)
	return text, ok
}

// replaceRecordContents swaps the record's payload for the parsed
// fields (@message stays readable for later processors).
func replaceRecordContents(record map[string]interface{}, parsed map[string]interface{}) {
	for k := range record {
		if k != "@message" {
			delete(record, k)
		}
	}
	for k, v := range parsed {
		record[k] = v
	}
}

// decodeEncodedValue unescapes a %-encoded value ("Encoded field
// values are decoded"). Path unescaping keeps any '+' verbatim; an
// invalid escape sequence leaves the value untouched.
func decodeEncodedValue(value string) string {
	decoded, err := url.PathUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}

// parseCloudfrontFields is the standard-log field order (Configure and
// use standard logs (access logs)); the parser accepts the format's
// whitespace-separated layout.
var parseCloudfrontFields = []string{
	"date", "time", "x-edge-location", "sc-bytes", "c-ip", "cs-method",
	"cs(Host)", "cs-uri-stem", "sc-status", "cs(Referer)", "cs(User-Agent)",
	"cs-uri-query", "cs(Cookie)", "x-edge-result-type", "x-edge-request-id",
	"x-host-header", "cs-protocol", "cs-bytes", "time-taken",
	"x-forwarded-for", "ssl-protocol", "ssl-cipher",
	"x-edge-response-result-type", "cs-protocol-version", "fle-status",
	"fle-encrypted-fields", "c-port", "time-to-first-byte",
	"x-edge-detailed-result-type", "sc-content-type", "sc-content-len",
	"sc-range-start", "sc-range-end",
}

// parseRoute53Fields is the Public Data Plane query log's field order.
var parseRoute53Fields = []string{
	"version", "queryTimestamp", "hostZoneId", "queryName", "queryType",
	"responseCode", "protocol", "edgeLocation", "resolverIp",
	"ednsClientSubnet",
}

// parseVPCFields is the flow log record's field order (version 2
// layout; the documented example carries the same fields).
var parseVPCFields = []string{
	"version", "accountId", "interfaceId", "srcAddr", "dstAddr",
	"srcPort", "dstPort", "protocol", "packets", "bytes", "start", "end",
	"action", "logStatus",
}

// parseVPCNumericFields are the flow log fields the documented output
// types as numbers (accountId stays a string).
var parseVPCNumericFields = map[string]bool{
	"version": true, "srcPort": true, "dstPort": true, "protocol": true,
	"packets": true, "bytes": true, "start": true, "end": true,
}

// applyParseWAF parses an AWS WAF vended log: the JSON event with
// httpRequest.headers converted from an array of {name, value} pairs to
// an object keyed by header name, and labels likewise ("It takes the
// contents of httpRequest.headers and creates JSON keys from each
// header name, with the corresponding value. It also does the same for
// labels").
func applyParseWAF(record map[string]interface{}, cfg map[string]interface{}) {
	message, ok := builtinSourceMessage(record, cfg)
	if !ok {
		return
	}
	var event map[string]interface{}
	if json.Unmarshal([]byte(message), &event) != nil {
		return
	}
	if httpRequest, isMap := event["httpRequest"].(map[string]interface{}); isMap {
		if headers, isList := httpRequest["headers"].([]interface{}); isList {
			converted := make(map[string]interface{}, len(headers))
			for _, header := range headers {
				pair, isMap := header.(map[string]interface{})
				if !isMap {
					continue
				}
				name, isString := pair["name"].(string)
				if !isString {
					continue
				}
				converted[name] = pair["value"]
			}
			httpRequest["headers"] = converted
		}
	}
	if labels, isList := event["labels"].([]interface{}); isList {
		converted := make(map[string]interface{}, len(labels))
		for _, label := range labels {
			pair, isMap := label.(map[string]interface{})
			if !isMap {
				continue
			}
			for k, v := range pair {
				if _, exists := converted[k]; !exists {
					converted[k] = v
				}
			}
		}
		event["labels"] = converted
	}
	replaceRecordContents(record, event)
}

// parsePostgresLine matches the RDS for PostgreSQL log-line prefix the
// documented example carries: time, remote host(port), user@database,
// process id and severity.
var parsePostgresLine = regexp.MustCompile(
	`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} [A-Za-z]+):` +
		`([\d.]+\(\d+\)):` +
		`([^:@]+)@([^:]+):` +
		`\[(\d+)\]:` +
		`([A-Z]+):` +
		`(.*)$`)

// applyParsePostgres parses an RDS for PostgreSQL vended log line into
// the six documented fields (logTime, srcIp, userName, dbName,
// processId, logLevel); the trailing message body is not part of the
// documented output.
func applyParsePostgres(record map[string]interface{}, cfg map[string]interface{}) {
	message, ok := builtinSourceMessage(record, cfg)
	if !ok {
		return
	}
	groups := parsePostgresLine.FindStringSubmatch(message)
	if groups == nil {
		return
	}
	replaceRecordContents(record, map[string]interface{}{
		"logTime":   groups[1],
		"srcIp":     groups[2],
		"userName":  groups[3],
		"dbName":    groups[4],
		"processId": groups[5],
		"logLevel":  groups[6],
	})
}

// applyParseCloudfront parses a CloudFront standard log line: the
// whitespace-separated standard fields, encoded values decoded and
// "values that are integers and doubles are treated as such".
func applyParseCloudfront(record map[string]interface{}, cfg map[string]interface{}) {
	parts := builtinSplitFields(record, cfg, len(parseCloudfrontFields))
	if parts == nil {
		return
	}
	out := make(map[string]interface{}, len(parseCloudfrontFields))
	for i, name := range parseCloudfrontFields {
		if i >= len(parts) {
			break
		}
		out[name] = typedEncodedValue(parts[i])
	}
	replaceRecordContents(record, out)
}

// applyParseRoute53 parses a Route 53 Public Data Plane query log
// ("This processor does not support Amazon Route 53 Resolver logs"):
// the ten whitespace-separated fields with the version typed as a
// double and every value decoded.
func applyParseRoute53(record map[string]interface{}, cfg map[string]interface{}) {
	parts := builtinSplitFields(record, cfg, len(parseRoute53Fields))
	if parts == nil {
		return
	}
	out := make(map[string]interface{}, len(parseRoute53Fields))
	for i, name := range parseRoute53Fields {
		if i >= len(parts) {
			break
		}
		value := decodeEncodedValue(parts[i])
		if name == "version" {
			if typed, err := strconv.ParseFloat(value, 64); err == nil {
				out[name] = typed
				continue
			}
		}
		out[name] = value
	}
	replaceRecordContents(record, out)
}

// applyParseVPC parses a VPC flow log record: the fourteen
// whitespace-separated fields with the documented numeric set typed as
// integers (accountId stays a string).
func applyParseVPC(record map[string]interface{}, cfg map[string]interface{}) {
	parts := builtinSplitFields(record, cfg, len(parseVPCFields))
	if parts == nil {
		return
	}
	out := make(map[string]interface{}, len(parseVPCFields))
	for i, name := range parseVPCFields {
		if i >= len(parts) {
			break
		}
		value := decodeEncodedValue(parts[i])
		if parseVPCNumericFields[name] {
			if typed, err := strconv.ParseInt(value, 10, 64); err == nil {
				out[name] = typed
				continue
			}
		}
		out[name] = value
	}
	replaceRecordContents(record, out)
}

// builtinSplitFields reads the parser's @message and splits it on
// whitespace runs. A record shorter than the documented field count
// leaves the event untransformed; a longer one (the later flow-log
// versions, for example) maps its leading fields onto the documented
// output schema, which the page publishes for one layout alone.
func builtinSplitFields(record map[string]interface{}, cfg map[string]interface{}, want int) []string {
	message, ok := builtinSourceMessage(record, cfg)
	if !ok {
		return nil
	}
	parts := strings.Fields(message)
	if len(parts) < want {
		return nil
	}
	return parts
}

// typedEncodedValue decodes a %-encoded field value and types the
// result: an integer or double value stays numeric, anything else
// (including "-") stays a string.
func typedEncodedValue(raw string) interface{} {
	value := decodeEncodedValue(raw)
	if typed, err := strconv.ParseInt(value, 10, 64); err == nil {
		return typed
	}
	if typed, err := strconv.ParseFloat(value, 64); err == nil {
		return typed
	}
	return value
}
