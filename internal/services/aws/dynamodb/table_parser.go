// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"vorpalstacks/internal/common/kmsutil"
	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func parseKeySchema(params map[string]interface{}) []*dbstore.KeySchemaElement {
	schema, ok := params["KeySchema"].([]interface{})
	if !ok {
		return nil
	}

	var result []*dbstore.KeySchemaElement
	for _, s := range schema {
		if sm, ok := s.(map[string]interface{}); ok {
			elem := &dbstore.KeySchemaElement{
				AttributeName: request.GetStringParam(sm, "AttributeName"),
				KeyType:       dbstore.KeyType(request.GetStringParam(sm, "KeyType")),
			}
			result = append(result, elem)
		}
	}
	return result
}

func parseAttributeDefinitions(params map[string]interface{}) []*dbstore.AttributeDefinition {
	defs, ok := params["AttributeDefinitions"].([]interface{})
	if !ok {
		return nil
	}

	var result []*dbstore.AttributeDefinition
	for _, d := range defs {
		if dm, ok := d.(map[string]interface{}); ok {
			def := &dbstore.AttributeDefinition{
				AttributeName: request.GetStringParam(dm, "AttributeName"),
				AttributeType: dbstore.ScalarAttributeType(request.GetStringParam(dm, "AttributeType")),
			}
			result = append(result, def)
		}
	}
	return result
}

// parseProvisionedThroughput decodes the ProvisionedThroughput member
// verbatim, returning nil only when the member is absent. Value bounds are
// validated by the Core (validateProvisionedThroughputValues), so a present
// member with out-of-range units still reaches the Core's rejection instead
// of masquerading as an omitted one.
func parseProvisionedThroughput(params map[string]interface{}) *dbstore.ProvisionedThroughput {
	pt, ok := params["ProvisionedThroughput"].(map[string]interface{})
	if !ok {
		return nil
	}

	return &dbstore.ProvisionedThroughput{
		ReadCapacityUnits:  request.GetInt64Param(pt, "ReadCapacityUnits"),
		WriteCapacityUnits: request.GetInt64Param(pt, "WriteCapacityUnits"),
	}
}

// parseOnDemandThroughput extracts the on-demand maximum-throughput pair
// shared by the table-level member and the GSI action members; an absent
// member returns nil.
func parseOnDemandThroughput(params map[string]interface{}) (*dbstore.OnDemandThroughput, error) {
	odt, ok := params["OnDemandThroughput"].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	return parseOnDemandThroughputMap(odt)
}

// parseOnDemandThroughputMap reads the on-demand maximum-throughput pair
// from one already-extracted member map — the table-level member, the GSI
// members and the restore Override member carry the same shape. The model
// requires a present member to name a limit of at least the minimum or
// the removal sentinel, and the pair to carry at least one member; an
// absent member stays unset (rendered as absent, never as zero).
func parseOnDemandThroughputMap(odt map[string]interface{}) (*dbstore.OnDemandThroughput, error) {
	parsed := &dbstore.OnDemandThroughput{}
	present := 0
	for member, slot := range map[string]*int64{
		"MaxReadRequestUnits":  &parsed.MaxReadRequestUnits,
		"MaxWriteRequestUnits": &parsed.MaxWriteRequestUnits,
	} {
		raw, ok := odt[member]
		if !ok {
			continue
		}
		f, ok := raw.(float64)
		if !ok || f != float64(int64(f)) {
			return nil, ErrInvalidParameter
		}
		v := int64(f)
		if v != dbstore.OnDemandThroughputRemoveValue && v < dbstore.OnDemandThroughputMinUnits {
			return nil, ErrInvalidParameter
		}
		*slot = v
		present++
	}
	if present == 0 {
		return nil, ErrInvalidParameter
	}
	return parsed, nil
}

// parseWarmThroughput extracts the warm-throughput pair shared by the
// table-level member and the GSI action members; an absent member returns
// nil.
func parseWarmThroughput(params map[string]interface{}) *dbstore.WarmThroughput {
	wt, ok := params["WarmThroughput"].(map[string]interface{})
	if !ok {
		return nil
	}
	return parseWarmThroughputMap(wt)
}

// parseWarmThroughputMap reads the warm-throughput pair from one
// already-extracted member map, shared by the table-level member, the GSI
// members and the restore Override member.
func parseWarmThroughputMap(wt map[string]interface{}) *dbstore.WarmThroughput {
	return &dbstore.WarmThroughput{
		ReadUnitsPerSecond:  request.GetInt64Param(wt, "ReadUnitsPerSecond"),
		WriteUnitsPerSecond: request.GetInt64Param(wt, "WriteUnitsPerSecond"),
	}
}

func parseGlobalSecondaryIndexes(params map[string]interface{}) ([]*dbstore.GlobalSecondaryIndex, error) {
	gsis, ok := params["GlobalSecondaryIndexes"].([]interface{})
	if !ok {
		return nil, nil
	}

	var result []*dbstore.GlobalSecondaryIndex
	for _, g := range gsis {
		if gm, ok := g.(map[string]interface{}); ok {
			idxName := request.GetStringParam(gm, "IndexName")
			if idxName == "" {
				return nil, ErrInvalidParameter
			}
			if !validateResourceName(idxName) {
				return nil, ErrInvalidParameter
			}
			proj, err := parseProjection(gm)
			if err != nil {
				return nil, err
			}
			if !validateProjectionRequired(gm["Projection"].(map[string]interface{})) {
				return nil, ErrInvalidParameter
			}
			idx := &dbstore.GlobalSecondaryIndex{
				IndexName:             idxName,
				KeySchema:             parseKeySchema(gm),
				Projection:            proj,
				ProvisionedThroughput: parseProvisionedThroughput(gm),
				IndexStatus:           dbstore.IndexStatusActive,
			}
			result = append(result, idx)
		}
	}
	return result, nil
}

func parseLocalSecondaryIndexes(params map[string]interface{}) ([]*dbstore.LocalSecondaryIndex, error) {
	lsis, ok := params["LocalSecondaryIndexes"].([]interface{})
	if !ok {
		return nil, nil
	}

	var result []*dbstore.LocalSecondaryIndex
	for _, l := range lsis {
		if lm, ok := l.(map[string]interface{}); ok {
			idxName := request.GetStringParam(lm, "IndexName")
			if idxName == "" {
				return nil, ErrInvalidParameter
			}
			if !validateResourceName(idxName) {
				return nil, ErrInvalidParameter
			}
			proj, err := parseProjection(lm)
			if err != nil {
				return nil, err
			}
			if !validateProjectionRequired(lm["Projection"].(map[string]interface{})) {
				return nil, ErrInvalidParameter
			}
			idx := &dbstore.LocalSecondaryIndex{
				IndexName:  idxName,
				KeySchema:  parseKeySchema(lm),
				Projection: proj,
			}
			result = append(result, idx)
		}
	}
	return result, nil
}

// parseVectorIndexes parses the CreateTable VectorIndexes member. Absence of
// the member yields nil (no vector indexes); field-level shape checks mirror
// the GSI/LSI parsers, while cross-index rules (name uniqueness across index
// families, shared-attribute dimension consistency) are enforced by the core.
func parseVectorIndexes(params map[string]interface{}) ([]*dbstore.VectorIndex, error) {
	rawList, ok := params["VectorIndexes"].([]interface{})
	if !ok {
		return nil, nil
	}
	return parseVectorIndexList(rawList)
}

// parseVectorIndexList parses one VectorIndexList member value — the
// CreateTable VectorIndexes member and the restore VectorIndexOverride
// member carry identical elements.
func parseVectorIndexList(rawList []interface{}) ([]*dbstore.VectorIndex, error) {
	result := make([]*dbstore.VectorIndex, 0, len(rawList))
	for _, v := range rawList {
		vm, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		idx, err := parseVectorIndex(vm)
		if err != nil {
			return nil, err
		}
		result = append(result, idx)
	}
	return result, nil
}

// parseVectorIndex parses one VectorIndex definition — a CreateTable
// VectorIndexes element or an UpdateTable VectorIndexUpdates Create action,
// which carry identical members.
func parseVectorIndex(vm map[string]interface{}) (*dbstore.VectorIndex, error) {
	idxName := request.GetStringParam(vm, "IndexName")
	if idxName == "" || !validateResourceName(idxName) {
		return nil, ErrInvalidParameter
	}
	vecAttr, ok := vm["VectorAttribute"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	attrName := request.GetStringParam(vecAttr, "AttributeName")
	if attrName == "" || len(attrName) > dbstore.VectorAttributeNameMax {
		return nil, ErrInvalidParameter
	}
	dims := request.GetInt64Param(vm, "Dimensions")
	if dims < dbstore.VectorDimensionsMin || dims > dbstore.VectorDimensionsMax {
		return nil, ErrInvalidParameter
	}
	distanceFn := request.GetStringParam(vm, "DistanceFunction")
	if !validateVectorDistanceFunction(distanceFn) {
		return nil, ErrInvalidParameter
	}
	proj, err := parseProjection(vm)
	if err != nil {
		return nil, err
	}
	if !validateProjectionRequired(vm["Projection"].(map[string]interface{})) {
		return nil, ErrInvalidParameter
	}
	searchSchema, err := parseSearchSchema(vm)
	if err != nil {
		return nil, err
	}
	return &dbstore.VectorIndex{
		IndexName:           idxName,
		VectorAttributeName: attrName,
		Dimensions:          dims,
		DistanceFunction:    dbstore.VectorDistanceFunction(distanceFn),
		Projection:          proj,
		SearchSchema:        searchSchema,
		IndexStatus:         dbstore.IndexStatusActive,
	}, nil
}

// parseSearchSchema parses the optional SearchSchema member: a non-empty
// list of {AttributeName, SearchSchemaElementType} elements with at most one
// HASH element and no duplicate attribute names.
func parseSearchSchema(vm map[string]interface{}) ([]*dbstore.SearchSchemaElement, error) {
	rawList, ok := vm["SearchSchema"].([]interface{})
	if !ok {
		return nil, nil
	}
	result := make([]*dbstore.SearchSchemaElement, 0, len(rawList))
	seen := make(map[string]bool)
	hasHash := false
	inlineFilters := 0
	for _, e := range rawList {
		em, ok := e.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		attrName := request.GetStringParam(em, "AttributeName")
		if attrName == "" {
			return nil, ErrInvalidParameter
		}
		elemType := request.GetStringParam(em, "SearchSchemaElementType")
		if elemType != "HASH" && elemType != "INLINE_FILTER" {
			return nil, ErrInvalidParameter
		}
		if seen[attrName] {
			return nil, ErrInvalidParameter
		}
		seen[attrName] = true
		if elemType == "HASH" {
			if hasHash {
				return nil, ErrInvalidParameter
			}
			hasHash = true
		} else {
			inlineFilters++
			if inlineFilters > dbstore.VectorInlineFiltersMax {
				return nil, ErrInvalidParameter
			}
		}
		result = append(result, &dbstore.SearchSchemaElement{
			AttributeName:           attrName,
			SearchSchemaElementType: dbstore.SearchSchemaElementType(elemType),
		})
	}
	if len(result) == 0 {
		return nil, ErrInvalidParameter
	}
	return result, nil
}

func parseProjection(params map[string]interface{}) (*dbstore.Projection, error) {
	proj, ok := params["Projection"].(map[string]interface{})
	if !ok {
		// Projection is required for GSI and LSI definitions.
		return nil, ErrInvalidParameter
	}

	p := &dbstore.Projection{
		ProjectionType: dbstore.ProjectionType(request.GetStringParam(proj, "ProjectionType")),
	}
	if p.ProjectionType == "" {
		p.ProjectionType = dbstore.ProjectionTypeAll
	}
	// Validate the ProjectionType enum (ALL, KEYS_ONLY, INCLUDE).
	if !validateProjectionType(string(p.ProjectionType)) {
		return nil, ErrInvalidParameter
	}

	if nkAs, ok := proj["NonKeyAttributes"].([]interface{}); ok {
		for _, nk := range nkAs {
			if nks, ok := nk.(string); ok {
				p.NonKeyAttributes = append(p.NonKeyAttributes, nks)
			}
		}
	}
	return p, nil
}

// parseStreamSpecification parses a StreamSpecification member into the
// stored form. The member is typed as a structure and StreamEnabled as a
// Boolean: a present value of another type on either is the
// invalid-parameter error, never a silently skipped specification — the
// typed-member family's discipline (wireStructMember); an absent member
// leaves the table's streams untouched.
func parseStreamSpecification(params map[string]interface{}) (*dbstore.StreamSpecification, error) {
	ss, present, err := wireStructMember(params["StreamSpecification"])
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}

	enabled, err := validateBoolParam(ss, "StreamEnabled", false)
	if err != nil {
		return nil, err
	}

	return &dbstore.StreamSpecification{
		StreamEnabled:  enabled,
		StreamViewType: dbstore.StreamViewType(request.GetStringParam(ss, "StreamViewType")),
	}, nil
}

// defaultSSEKeyAlias is the model's name for the default DynamoDB KMS key;
// resolveSSEKeyArn echoes its alias ARN when a specification enables KMS
// without an explicit key identifier.
const defaultSSEKeyAlias = "alias/aws/dynamodb"

// SetKMSResolver wires the KMS key-identifier resolver the SSE write paths
// settle key ARNs through.
func (s *DynamoDBService) SetKMSResolver(resolver kmsutil.Resolver) {
	s.kmsResolver = resolver
}

// parseSSESpecification parses an SSESpecification member into the stored
// description form. The member is typed as a structure and Enabled as a
// Boolean: a present value of another type on either is the
// invalid-parameter error, never a silently skipped specification — the
// typed-member family's discipline (wireStructMember, wireBoolMember).
// The specification side of the model is asymmetric by design:
// Enabled=true sets the encryption type to KMS with an AWS managed key
// ("If enabled (true), server-side encryption type is set to KMS and an
// Amazon Web Services managed key is used"), SSEType's only supported
// specification value is KMS, and the key identifier is optional — "you
// should only provide this parameter if the key is different from the
// default DynamoDB key alias/aws/dynamodb". The identifier's form is
// validated here (validateKMSKeyIdentifier — a KMS key or alias ARN, an
// alias name, or a bare key id); the ARN form
// SSEDescription.KMSMasterKeyArn documents is settled by resolveSSEKeyArn
// in the operation core, where the request's region and account live. An
// explicit Enabled=false states the AWS owned-key default and carries no
// settings to apply; the second return reports it so each operation core
// can apply its own policy to a disable request. A specification whose
// Enabled is absent but which carries SSEType or KMSMasterKeyId
// contradicts the owned-key default that the absence states.
func parseSSESpecification(raw interface{}) (*dbstore.SSEDescription, bool, error) {
	ss, present, err := wireStructMember(raw)
	if err != nil {
		return nil, false, err
	}
	if !present {
		return nil, false, nil
	}
	if enabled, present, err := wireBoolMember(ss["Enabled"]); err != nil {
		return nil, false, err
	} else if present {
		if !enabled {
			return nil, true, nil
		}
		sseType := dbstore.SSEType(request.GetStringParam(ss, "SSEType"))
		if sseType == "" {
			sseType = dbstore.SSETypeKMS
		}
		if sseType != dbstore.SSETypeKMS {
			return nil, false, ErrInvalidParameter
		}
		keyID := request.GetStringParam(ss, "KMSMasterKeyId")
		if keyID != "" && !validateKMSKeyIdentifier(keyID) {
			return nil, false, ErrInvalidParameter
		}
		return &dbstore.SSEDescription{
			Status:          "ENABLED",
			SSEType:         dbstore.SSETypeKMS,
			KMSMasterKeyArn: keyID,
		}, false, nil
	}
	if _, hasType := ss["SSEType"]; hasType {
		return nil, false, ErrInvalidParameter
	}
	if _, hasKey := ss["KMSMasterKeyId"]; hasKey {
		return nil, false, ErrInvalidParameter
	}
	return nil, false, nil
}
