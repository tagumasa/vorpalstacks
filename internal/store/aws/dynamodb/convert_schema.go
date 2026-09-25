package dynamodb

import (
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func warmThroughputToProto(w *WarmThroughput) *pb.WarmThroughput {
	if w == nil {
		return nil
	}
	return &pb.WarmThroughput{
		ReadUnitsPerSecond:  w.ReadUnitsPerSecond,
		WriteUnitsPerSecond: w.WriteUnitsPerSecond,
	}
}

func protoToWarmThroughput(p *pb.WarmThroughput) *WarmThroughput {
	if p == nil {
		return nil
	}
	return &WarmThroughput{
		ReadUnitsPerSecond:  p.ReadUnitsPerSecond,
		WriteUnitsPerSecond: p.WriteUnitsPerSecond,
	}
}

func onDemandThroughputToProto(o *OnDemandThroughput) *pb.OnDemandThroughput {
	if o == nil {
		return nil
	}
	return &pb.OnDemandThroughput{
		MaxReadRequestUnits:  o.MaxReadRequestUnits,
		MaxWriteRequestUnits: o.MaxWriteRequestUnits,
	}
}

func protoToOnDemandThroughput(p *pb.OnDemandThroughput) *OnDemandThroughput {
	if p == nil {
		return nil
	}
	return &OnDemandThroughput{
		MaxReadRequestUnits:  p.MaxReadRequestUnits,
		MaxWriteRequestUnits: p.MaxWriteRequestUnits,
	}
}

// Struct conversion functions

func keySchemaToProto(ks []*KeySchemaElement) []*pb.KeySchemaElement {
	if ks == nil {
		return nil
	}
	result := make([]*pb.KeySchemaElement, len(ks))
	for i, k := range ks {
		result[i] = &pb.KeySchemaElement{
			AttributeName: k.AttributeName,
			KeyType:       keyTypeToProto(k.KeyType),
		}
	}
	return result
}

func protoToKeySchema(ks []*pb.KeySchemaElement) []*KeySchemaElement {
	if ks == nil {
		return nil
	}
	result := make([]*KeySchemaElement, len(ks))
	for i, k := range ks {
		result[i] = &KeySchemaElement{
			AttributeName: k.AttributeName,
			KeyType:       protoToKeyType(k.KeyType),
		}
	}
	return result
}

func attributeDefinitionsToProto(ad []*AttributeDefinition) []*pb.AttributeDefinition {
	if ad == nil {
		return nil
	}
	result := make([]*pb.AttributeDefinition, len(ad))
	for i, a := range ad {
		result[i] = &pb.AttributeDefinition{
			AttributeName: a.AttributeName,
			AttributeType: scalarAttributeTypeToProto(a.AttributeType),
		}
	}
	return result
}

func protoToAttributeDefinitions(ad []*pb.AttributeDefinition) []*AttributeDefinition {
	if ad == nil {
		return nil
	}
	result := make([]*AttributeDefinition, len(ad))
	for i, a := range ad {
		result[i] = &AttributeDefinition{
			AttributeName: a.AttributeName,
			AttributeType: protoToScalarAttributeType(a.AttributeType),
		}
	}
	return result
}

func provisionedThroughputToProto(pt *ProvisionedThroughput) *pb.ProvisionedThroughput {
	if pt == nil {
		return nil
	}
	return &pb.ProvisionedThroughput{
		ReadCapacityUnits:      pt.ReadCapacityUnits,
		WriteCapacityUnits:     pt.WriteCapacityUnits,
		LastDecreaseDateTime:   timestamppb.New(pt.LastDecreaseDateTime),
		LastIncreaseDateTime:   timestamppb.New(pt.LastIncreaseDateTime),
		NumberOfDecreasesToday: pt.NumberOfDecreasesToday,
	}
}

func protoToProvisionedThroughput(pt *pb.ProvisionedThroughput) *ProvisionedThroughput {
	if pt == nil {
		return nil
	}
	return &ProvisionedThroughput{
		ReadCapacityUnits:      pt.ReadCapacityUnits,
		WriteCapacityUnits:     pt.WriteCapacityUnits,
		LastDecreaseDateTime:   pt.LastDecreaseDateTime.AsTime(),
		LastIncreaseDateTime:   pt.LastIncreaseDateTime.AsTime(),
		NumberOfDecreasesToday: pt.NumberOfDecreasesToday,
	}
}

func projectionToProto(p *Projection) *pb.Projection {
	if p == nil {
		return nil
	}
	return &pb.Projection{
		ProjectionType:   projectionTypeToProto(p.ProjectionType),
		NonKeyAttributes: p.NonKeyAttributes,
	}
}

func protoToProjection(p *pb.Projection) *Projection {
	if p == nil {
		return nil
	}
	return &Projection{
		ProjectionType:   protoToProjectionType(p.ProjectionType),
		NonKeyAttributes: p.NonKeyAttributes,
	}
}

func vectorIndexesToProto(idx []*VectorIndex) []*pb.VectorIndex {
	if idx == nil {
		return nil
	}
	result := make([]*pb.VectorIndex, len(idx))
	for i, v := range idx {
		result[i] = &pb.VectorIndex{
			IndexName:           v.IndexName,
			IndexArn:            v.IndexArn,
			VectorAttributeName: v.VectorAttributeName,
			Dimensions:          v.Dimensions,
			DistanceFunction:    vectorDistanceFunctionToProto(v.DistanceFunction),
			Projection:          projectionToProto(v.Projection),
			SearchSchema:        searchSchemaToProto(v.SearchSchema),
			IndexStatus:         indexStatusToProto(v.IndexStatus),
			Backfilling:         v.Backfilling,
			IndexSizeBytes:      v.IndexSizeBytes,
			ItemCount:           v.ItemCount,
		}
	}
	return result
}

func protoToVectorIndexes(idx []*pb.VectorIndex) []*VectorIndex {
	if idx == nil {
		return nil
	}
	result := make([]*VectorIndex, len(idx))
	for i, v := range idx {
		result[i] = &VectorIndex{
			IndexName:           v.IndexName,
			IndexArn:            v.IndexArn,
			VectorAttributeName: v.VectorAttributeName,
			Dimensions:          v.Dimensions,
			DistanceFunction:    protoToVectorDistanceFunction(v.DistanceFunction),
			Projection:          protoToProjection(v.Projection),
			SearchSchema:        protoToSearchSchema(v.SearchSchema),
			IndexStatus:         protoToIndexStatus(v.IndexStatus),
			Backfilling:         v.Backfilling,
			IndexSizeBytes:      v.IndexSizeBytes,
			ItemCount:           v.ItemCount,
		}
	}
	return result
}

func searchSchemaToProto(schema []*SearchSchemaElement) []*pb.SearchSchemaElement {
	if schema == nil {
		return nil
	}
	result := make([]*pb.SearchSchemaElement, len(schema))
	for i, e := range schema {
		result[i] = &pb.SearchSchemaElement{
			AttributeName:           e.AttributeName,
			SearchSchemaElementType: searchSchemaElementTypeToProto(e.SearchSchemaElementType),
		}
	}
	return result
}

func protoToSearchSchema(schema []*pb.SearchSchemaElement) []*SearchSchemaElement {
	if schema == nil {
		return nil
	}
	result := make([]*SearchSchemaElement, len(schema))
	for i, e := range schema {
		result[i] = &SearchSchemaElement{
			AttributeName:           e.AttributeName,
			SearchSchemaElementType: protoToSearchSchemaElementType(e.SearchSchemaElementType),
		}
	}
	return result
}

func globalSecondaryIndexesToProto(idx []*GlobalSecondaryIndex) []*pb.GlobalSecondaryIndex {
	if idx == nil {
		return nil
	}
	result := make([]*pb.GlobalSecondaryIndex, len(idx))
	for i, g := range idx {
		result[i] = &pb.GlobalSecondaryIndex{
			IndexName:             g.IndexName,
			IndexArn:              g.IndexArn,
			KeySchema:             keySchemaToProto(g.KeySchema),
			Projection:            projectionToProto(g.Projection),
			ProvisionedThroughput: provisionedThroughputToProto(g.ProvisionedThroughput),
			OnDemandThroughput:    onDemandThroughputToProto(g.OnDemandThroughput),
			WarmThroughput:        warmThroughputToProto(g.WarmThroughput),
			IndexStatus:           indexStatusToProto(g.IndexStatus),
			IndexSizeBytes:        g.IndexSizeBytes,
			ItemCount:             g.ItemCount,
		}
	}
	return result
}

func protoToGlobalSecondaryIndexes(idx []*pb.GlobalSecondaryIndex) []*GlobalSecondaryIndex {
	if idx == nil {
		return nil
	}
	result := make([]*GlobalSecondaryIndex, len(idx))
	for i, g := range idx {
		result[i] = &GlobalSecondaryIndex{
			IndexName:             g.IndexName,
			IndexArn:              g.IndexArn,
			KeySchema:             protoToKeySchema(g.KeySchema),
			Projection:            protoToProjection(g.Projection),
			ProvisionedThroughput: protoToProvisionedThroughput(g.ProvisionedThroughput),
			OnDemandThroughput:    protoToOnDemandThroughput(g.OnDemandThroughput),
			WarmThroughput:        protoToWarmThroughput(g.WarmThroughput),
			IndexStatus:           protoToIndexStatus(g.IndexStatus),
			IndexSizeBytes:        g.IndexSizeBytes,
			ItemCount:             g.ItemCount,
		}
	}
	return result
}

func localSecondaryIndexesToProto(idx []*LocalSecondaryIndex) []*pb.LocalSecondaryIndex {
	if idx == nil {
		return nil
	}
	result := make([]*pb.LocalSecondaryIndex, len(idx))
	for i, l := range idx {
		result[i] = &pb.LocalSecondaryIndex{
			IndexName:      l.IndexName,
			KeySchema:      keySchemaToProto(l.KeySchema),
			Projection:     projectionToProto(l.Projection),
			IndexSizeBytes: l.IndexSizeBytes,
			ItemCount:      l.ItemCount,
		}
	}
	return result
}

func protoToLocalSecondaryIndexes(idx []*pb.LocalSecondaryIndex) []*LocalSecondaryIndex {
	if idx == nil {
		return nil
	}
	result := make([]*LocalSecondaryIndex, len(idx))
	for i, l := range idx {
		result[i] = &LocalSecondaryIndex{
			IndexName:      l.IndexName,
			KeySchema:      protoToKeySchema(l.KeySchema),
			Projection:     protoToProjection(l.Projection),
			IndexSizeBytes: l.IndexSizeBytes,
			ItemCount:      l.ItemCount,
		}
	}
	return result
}
