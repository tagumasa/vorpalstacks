package cloudtrail

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	"vorpalstacks/pkg/sqlparser"
)

// ---------------------------------------------------------------------------
// CloudTrail Lake SQL column vocabulary
// ---------------------------------------------------------------------------

// lakeColumns holds the canonical AWS CloudTrail Lake SQL schema column
// names for event records, lowercase exactly as the schema reference lists
// them (AWS "Supported SQL schemas for event data stores": eventversion,
// useridentity, eventtime, eventsource, eventname, awsregion,
// sourceipaddress, useragent, errorcode, errormessage, requestparameters,
// responseelements, additionaleventdata, requestid, eventid, readonly,
// resources, eventtype, apiversion, managementevent, recipientaccountid,
// sharedeventid, annotation, vpcendpointid, vpcendpointaccountid,
// serviceeventdetails, addendum, edgedevicedetails, insightdetails,
// eventcategory, tlsdetails, sessioncredentialfromconsole, eventjson,
// eventjsonchecksum). This vocabulary is the query engine's own contract,
// deliberately decoupled from the LookupEvents wire formatter (formatEvent)
// whose PascalCase keys are a different surface — the engine once resolved
// against that key set and lost every projection and WHERE match.
//
// Resolution is case-insensitive: the schema reference lists lowercase
// names while AWS sample queries and the record format spell them
// camelCase (eventID, eventTime). Columns the store does not record are
// part of the vocabulary but resolve to nil — an unrecorded record field
// is null, never fabricated.
var lakeColumns = map[string]bool{
	"eventversion":                 true,
	"useridentity":                 true,
	"eventtime":                    true,
	"eventsource":                  true,
	"eventname":                    true,
	"awsregion":                    true,
	"sourceipaddress":              true,
	"useragent":                    true,
	"errorcode":                    true,
	"errormessage":                 true,
	"requestparameters":            true,
	"responseelements":             true,
	"additionaleventdata":          true,
	"requestid":                    true,
	"eventid":                      true,
	"readonly":                     true,
	"resources":                    true,
	"eventtype":                    true,
	"apiversion":                   true,
	"managementevent":              true,
	"recipientaccountid":           true,
	"sharedeventid":                true,
	"annotation":                   true,
	"vpcendpointid":                true,
	"vpcendpointaccountid":         true,
	"serviceeventdetails":          true,
	"addendum":                     true,
	"edgedevicedetails":            true,
	"insightdetails":               true,
	"eventcategory":                true,
	"tlsdetails":                   true,
	"sessioncredentialfromconsole": true,
	"eventjson":                    true,
	"eventjsonchecksum":            true,
}

// lakeStructFields lists the struct columns the schema defines, keyed by
// canonical column name. Qualified references (userIdentity.userName,
// resources.ARN) validate their root against this set; the leaf resolves
// leniently because the nested structs reach deeper than the fields the
// store records — an unrecorded leaf resolves to nil.
var lakeStructFields = map[string]bool{
	"useridentity": true,
	"resources":    true,
}

// lakeColumnList is the deterministic ordering of lakeColumns for SELECT *
// result rows.
var lakeColumnList = func() []string {
	cols := make([]string, 0, len(lakeColumns))
	for col := range lakeColumns {
		cols = append(cols, col)
	}
	sort.Strings(cols)
	return cols
}()

// lakeColumnCanonical lowercases a column reference for vocabulary lookup.
func lakeColumnCanonical(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// lakeRow builds the query engine's row for one stored event, keyed by the
// canonical Lake column names. Struct columns are nested maps keyed by the
// schema's lowercase field names; resources entries carry arn and type.
// Columns with no recorded value stay absent so they resolve to nil.
func lakeRow(e *cloudtrailstore.Event) map[string]interface{} {
	row := map[string]interface{}{
		"eventversion":      e.EventVersion,
		"eventtime":         e.EventTime.UTC(),
		"eventsource":       e.EventSource,
		"eventname":         e.EventName,
		"awsregion":         e.AwsRegion,
		"sourceipaddress":   e.SourceIPAddress,
		"useragent":         e.UserAgent,
		"errorcode":         e.ErrorCode,
		"errormessage":      e.ErrorMessage,
		"requestparameters": e.RequestParameters,
		"responseelements":  e.ResponseElements,
		"requestid":         e.RequestID,
		"eventid":           e.EventID,
		"readonly":          e.ReadOnly,
		"eventtype":         e.EventType,
		// managementevent follows the record's category: platform-recorded
		// events are Management; channel-delivered events are
		// ActivityAuditLog and are not management events.
		"managementevent": e.EventCategory == "Management",
		"eventcategory":   e.EventCategory,
		"eventjson":       e.CloudTrailEvent,
	}
	if e.UserIdentity != nil {
		identity := map[string]interface{}{
			"type":        e.UserIdentity.Type,
			"principalid": e.UserIdentity.PrincipalID,
			"arn":         e.UserIdentity.ARN,
			"accountid":   e.UserIdentity.AccountID,
			"accesskeyid": e.UserIdentity.AccessKeyID,
			"username":    e.UserIdentity.UserName,
		}
		row["useridentity"] = identity
	}
	if len(e.Resources) > 0 {
		entries := make([]map[string]interface{}, 0, len(e.Resources))
		for _, r := range e.Resources {
			entries = append(entries, map[string]interface{}{
				"arn":  r.ResourceName,
				"type": r.ResourceType,
			})
		}
		row["resources"] = entries
	}
	return row
}

// resolveLakeColumn resolves an already-canonicalised column reference
// against a Lake row. Qualified references traverse struct columns; a
// resources field resolves to every entry's value, so an equality or LIKE
// comparison matches when any resource entry matches.
func resolveLakeColumn(row map[string]interface{}, qualifier, name string) interface{} {
	if qualifier == "" {
		return row[name]
	}
	switch parent := row[qualifier].(type) {
	case map[string]interface{}:
		return parent[name]
	case []map[string]interface{}:
		vals := make([]string, 0, len(parent))
		for _, entry := range parent {
			if v, ok := entry[name]; ok && v != nil {
				if s := lakeValueString(v); s != "" {
					vals = append(vals, s)
				}
			}
		}
		if len(vals) == 0 {
			return nil
		}
		return vals
	default:
		return nil
	}
}

// lakeValueString renders a Lake column value for the stringly result rows
// and for string/LIKE comparison: timestamps as whole-second RFC3339 UTC
// (the record format's own spelling), composites as their JSON encoding.
func lakeValueString(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return ""
	case time.Time:
		return val.UTC().Truncate(time.Second).Format(time.RFC3339)
	case string:
		return val
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(b)
	}
}

// validateLakeColumnRef rejects a column reference that is not part of the
// Lake schema: an unqualified name must be a schema column, and a
// qualified reference's root must be a struct column. The error shape is
// StartQuery's declared InvalidQueryStatementException.
func validateLakeColumnRef(qualifier, name string) error {
	asWritten := name
	if qualifier != "" {
		asWritten = qualifier + "." + name
	}
	if qualifier == "" {
		if !lakeColumns[lakeColumnCanonical(name)] {
			return newInvalidQueryStatementException(
				fmt.Sprintf("Unknown column: %s", asWritten))
		}
		return nil
	}
	if !lakeStructFields[lakeColumnCanonical(qualifier)] {
		return newInvalidQueryStatementException(
			fmt.Sprintf("Unknown column: %s", asWritten))
	}
	return nil
}

// walkColumnRefs visits every column reference in a WHERE expression,
// mirroring the expression types evaluateWhere supports so validation and
// evaluation stay in lockstep. The first visitor error stops the walk.
func walkColumnRefs(expr sqlparser.Expr, visit func(cn *sqlparser.ColName) error) error {
	switch e := expr.(type) {
	case *sqlparser.ColName:
		return visit(e)
	case *sqlparser.ComparisonExpr:
		if err := walkColumnRefs(e.Left, visit); err != nil {
			return err
		}
		return walkColumnRefs(e.Right, visit)
	case *sqlparser.AndExpr:
		if err := walkColumnRefs(e.Left, visit); err != nil {
			return err
		}
		return walkColumnRefs(e.Right, visit)
	case *sqlparser.OrExpr:
		if err := walkColumnRefs(e.Left, visit); err != nil {
			return err
		}
		return walkColumnRefs(e.Right, visit)
	case *sqlparser.ParenExpr:
		return walkColumnRefs(e.Expr, visit)
	case *sqlparser.IsExpr:
		return walkColumnRefs(e.Expr, visit)
	case *sqlparser.NotExpr:
		return walkColumnRefs(e.Expr, visit)
	case *sqlparser.RangeCond:
		if err := walkColumnRefs(e.Left, visit); err != nil {
			return err
		}
		if err := walkColumnRefs(e.From, visit); err != nil {
			return err
		}
		return walkColumnRefs(e.To, visit)
	case sqlparser.ValTuple:
		for _, item := range e {
			if err := walkColumnRefs(item, visit); err != nil {
				return err
			}
		}
	}
	return nil
}

// colNameParts splits a parsed column reference into its canonical
// qualifier and leaf names. A three-part reference (userIdentity.session-
// context.mfaAuthenticated) resolves from its struct root; the store does
// not record the nested structs, so such references resolve to nil.
func colNameParts(cn *sqlparser.ColName) (qualifier, name string) {
	if !cn.Qualifier.IsEmpty() && !cn.Qualifier.Qualifier.IsEmpty() {
		return lakeColumnCanonical(cn.Qualifier.Qualifier.String()),
			lakeColumnCanonical(cn.Qualifier.Name.String())
	}
	if !cn.Qualifier.IsEmpty() {
		return lakeColumnCanonical(cn.Qualifier.Name.String()),
			lakeColumnCanonical(cn.Name.String())
	}
	return "", lakeColumnCanonical(cn.Name.String())
}
