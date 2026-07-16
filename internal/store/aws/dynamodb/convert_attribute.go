package dynamodb

import (
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

func attributeValueToProto(av *AttributeValue) *pb.AttributeValue {
	if av == nil {
		return nil
	}

	switch {
	case av.S != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_S{S: *av.S}}
	case av.N != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_N{N: *av.N}}
	case av.B != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_B{B: av.B}}
	case av.SS != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_Ss{Ss: &pb.StringSet{Values: av.SS}}}
	case av.NS != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_Ns{Ns: &pb.NumberSet{Values: av.NS}}}
	case av.BS != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_Bs{Bs: &pb.BytesSet{Values: av.BS}}}
	case av.M != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_M{M: attributeValueMapToProto(av.M)}}
	case av.L != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_L{L: attributeValueListToProto(av.L)}}
	case av.NULL != nil && *av.NULL:
		return &pb.AttributeValue{Value: &pb.AttributeValue_Null{Null: &pb.NullValue{}}}
	case av.BOOL != nil:
		return &pb.AttributeValue{Value: &pb.AttributeValue_BoolVal{BoolVal: *av.BOOL}}
	default:
		return nil
	}
}

func protoToAttributeValue(p *pb.AttributeValue) *AttributeValue {
	if p == nil {
		return nil
	}

	switch v := p.Value.(type) {
	case *pb.AttributeValue_S:
		return &AttributeValue{S: &v.S}
	case *pb.AttributeValue_N:
		return &AttributeValue{N: &v.N}
	case *pb.AttributeValue_B:
		return &AttributeValue{B: v.B}
	case *pb.AttributeValue_Ss:
		return &AttributeValue{SS: v.Ss.Values}
	case *pb.AttributeValue_Ns:
		return &AttributeValue{NS: v.Ns.Values}
	case *pb.AttributeValue_Bs:
		return &AttributeValue{BS: v.Bs.Values}
	case *pb.AttributeValue_M:
		return &AttributeValue{M: protoToAttributeValueMap(v.M)}
	case *pb.AttributeValue_L:
		return &AttributeValue{L: protoToAttributeValueList(v.L)}
	case *pb.AttributeValue_Null:
		t := true
		return &AttributeValue{NULL: &t}
	case *pb.AttributeValue_BoolVal:
		return &AttributeValue{BOOL: &v.BoolVal}
	default:
		return nil
	}
}

func attributeValueMapToProto(m map[string]*AttributeValue) *pb.MapValue {
	if m == nil {
		return nil
	}
	entries := make(map[string]*pb.AttributeValue, len(m))
	for k, v := range m {
		entries[k] = attributeValueToProto(v)
	}
	return &pb.MapValue{Entries: entries}
}

func protoToAttributeValueMap(p *pb.MapValue) map[string]*AttributeValue {
	if p == nil {
		return nil
	}
	m := make(map[string]*AttributeValue, len(p.Entries))
	for k, v := range p.Entries {
		m[k] = protoToAttributeValue(v)
	}
	return m
}

func attributeValueListToProto(l []*AttributeValue) *pb.ListValue {
	if l == nil {
		return nil
	}
	values := make([]*pb.AttributeValue, len(l))
	for i, v := range l {
		values[i] = attributeValueToProto(v)
	}
	return &pb.ListValue{Values: values}
}

func protoToAttributeValueList(p *pb.ListValue) []*AttributeValue {
	if p == nil {
		return nil
	}
	l := make([]*AttributeValue, len(p.Values))
	for i, v := range p.Values {
		l[i] = protoToAttributeValue(v)
	}
	return l
}

func attributeValueMapToProtoDirect(m map[string]*AttributeValue) map[string]*pb.AttributeValue {
	if m == nil {
		return nil
	}
	result := make(map[string]*pb.AttributeValue, len(m))
	for k, v := range m {
		result[k] = attributeValueToProto(v)
	}
	return result
}

func protoToAttributeValueMapDirect(m map[string]*pb.AttributeValue) map[string]*AttributeValue {
	if m == nil {
		return nil
	}
	result := make(map[string]*AttributeValue, len(m))
	for k, v := range m {
		result[k] = protoToAttributeValue(v)
	}
	return result
}
