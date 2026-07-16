package dynamodb

import (
	"time"

	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/aws/dynamodb"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func storeAVToProto(av *dbstore.AttributeValue) *pb.AttributeValue {
	if av == nil {
		return nil
	}
	p := &pb.AttributeValue{}
	if av.S != nil {
		p.S = *av.S
	}
	if av.N != nil {
		p.N = *av.N
	}
	if av.B != nil {
		p.B = av.B
	}
	if av.BOOL != nil {
		p.Bool = proto.Bool(*av.BOOL)
	}
	if av.NULL != nil && *av.NULL {
		p.Null = proto.Bool(true)
	}
	if av.SS != nil {
		p.Ss = av.SS
	}
	if av.NS != nil {
		p.Ns = av.NS
	}
	if av.BS != nil {
		p.Bs = av.BS
	}
	if av.M != nil {
		p.M = storeAVMapToProto(av.M)
	}
	if av.L != nil {
		p.L = storeAVListToProto(av.L)
	}
	return p
}

func protoAVToStore(p *pb.AttributeValue) *dbstore.AttributeValue {
	if p == nil {
		return nil
	}
	av := &dbstore.AttributeValue{}
	if p.S != "" {
		av.S = &p.S
	}
	if p.N != "" {
		av.N = &p.N
	}
	if len(p.B) > 0 {
		av.B = p.B
	}
	if p.GetNull() {
		av.NULL = p.Null
	}
	if len(p.Ss) > 0 {
		av.SS = p.Ss
	}
	if len(p.Ns) > 0 {
		av.NS = p.Ns
	}
	if len(p.Bs) > 0 {
		av.BS = p.Bs
	}
	if p.M != nil {
		av.M = protoAVMapToStore(p.M)
	}
	if p.L != nil {
		av.L = protoAVListToStore(p.L)
	}
	if av.S == nil && av.N == nil && av.B == nil && av.NULL == nil &&
		av.SS == nil && av.NS == nil && av.BS == nil && av.M == nil && av.L == nil {
		av.BOOL = p.Bool
	}
	return av
}

func storeAVMapToProto(m map[string]*dbstore.AttributeValue) map[string]*pb.AttributeValue {
	if m == nil {
		return nil
	}
	result := make(map[string]*pb.AttributeValue, len(m))
	for k, v := range m {
		result[k] = storeAVToProto(v)
	}
	return result
}

func protoAVMapToStore(m map[string]*pb.AttributeValue) map[string]*dbstore.AttributeValue {
	if m == nil {
		return nil
	}
	result := make(map[string]*dbstore.AttributeValue, len(m))
	for k, v := range m {
		result[k] = protoAVToStore(v)
	}
	return result
}

func storeAVListToProto(l []*dbstore.AttributeValue) []*pb.AttributeValue {
	if l == nil {
		return nil
	}
	result := make([]*pb.AttributeValue, len(l))
	for i, v := range l {
		result[i] = storeAVToProto(v)
	}
	return result
}

func protoAVListToStore(l []*pb.AttributeValue) []*dbstore.AttributeValue {
	if l == nil {
		return nil
	}
	result := make([]*dbstore.AttributeValue, len(l))
	for i, v := range l {
		result[i] = protoAVToStore(v)
	}
	return result
}

func extractKeyFromAttributes(table *dbstore.Table, attrs map[string]*dbstore.AttributeValue) map[string]*dbstore.AttributeValue {
	key := make(map[string]*dbstore.AttributeValue, len(table.KeySchema))
	for _, ks := range table.KeySchema {
		attr, ok := attrs[ks.AttributeName]
		if !ok {
			return nil
		}
		key[ks.AttributeName] = attr
	}
	return key
}

func storeTableToProtoDescription(table *dbstore.Table) *pb.TableDescription {
	if table == nil {
		return nil
	}
	desc := &pb.TableDescription{
		Tablename:                 table.Name,
		Tablearn:                  table.ARN,
		Tablestatus:               tableStatusToProto(table.Status),
		Creationdatetime:          table.CreationDateTime.Format(time.RFC3339),
		Itemcount:                 table.ItemCount,
		Tablesizebytes:            table.TableSizeBytes,
		Deletionprotectionenabled: proto.Bool(table.DeletionProtectionEnabled),
	}

	for _, ks := range table.KeySchema {
		desc.Keyschema = append(desc.Keyschema, &pb.KeySchemaElement{
			Attributename: ks.AttributeName,
			Keytype:       keyTypeToProto(ks.KeyType),
		})
	}

	for _, ad := range table.AttributeDefinitions {
		desc.Attributedefinitions = append(desc.Attributedefinitions, &pb.AttributeDefinition{
			Attributename: ad.AttributeName,
			Attributetype: attrTypeToProto(ad.AttributeType),
		})
	}

	desc.Billingmodesummary = &pb.BillingModeSummary{
		Billingmode: billingModeToProto(table.BillingMode),
	}

	return desc
}

func tableStatusToProto(status dbstore.TableStatus) pb.TableStatus {
	switch status {
	case dbstore.TableStatusActive:
		return pb.TableStatus_TABLE_STATUS_ACTIVE
	case dbstore.TableStatusCreating:
		return pb.TableStatus_TABLE_STATUS_CREATING
	case dbstore.TableStatusDeleting:
		return pb.TableStatus_TABLE_STATUS_DELETING
	case dbstore.TableStatusUpdating:
		return pb.TableStatus_TABLE_STATUS_UPDATING
	default:
		return pb.TableStatus_TABLE_STATUS_ACTIVE
	}
}

func keyTypeToProto(kt dbstore.KeyType) pb.KeyType {
	switch kt {
	case dbstore.KeyTypeHash:
		return pb.KeyType_KEY_TYPE_HASH
	case dbstore.KeyTypeRange:
		return pb.KeyType_KEY_TYPE_RANGE
	default:
		return pb.KeyType_KEY_TYPE_HASH
	}
}

func attrTypeToProto(at dbstore.ScalarAttributeType) pb.ScalarAttributeType {
	switch at {
	case dbstore.ScalarAttributeTypeS:
		return pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_S
	case dbstore.ScalarAttributeTypeN:
		return pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_N
	case dbstore.ScalarAttributeTypeB:
		return pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_B
	default:
		return pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_S
	}
}

func billingModeToProto(bm dbstore.BillingMode) pb.BillingMode {
	switch bm {
	case dbstore.BillingModeProvisioned:
		return pb.BillingMode_BILLING_MODE_PROVISIONED
	default:
		return pb.BillingMode_BILLING_MODE_PAY_PER_REQUEST
	}
}
