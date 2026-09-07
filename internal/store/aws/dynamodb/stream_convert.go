package dynamodb

import (
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

func streamUserIdentityToProto(u *StreamUserIdentity) *pb.StreamUserIdentity {
	if u == nil {
		return nil
	}
	return &pb.StreamUserIdentity{IdentityType: u.Type, PrincipalId: u.PrincipalID}
}

func protoToStreamUserIdentity(p *pb.StreamUserIdentity) *StreamUserIdentity {
	if p == nil {
		return nil
	}
	return &StreamUserIdentity{Type: p.IdentityType, PrincipalID: p.PrincipalId}
}

// wireImageToProto converts one wire-shaped image map to its typed protobuf
// form; nil stays nil so absent images stay absent.
func wireImageToProto(m map[string]interface{}) map[string]*pb.AttributeValue {
	if m == nil {
		return nil
	}
	return attributeValueMapToProtoDirect(ParseItemWire(m))
}

// protoImageToWire converts one typed protobuf image map back to the wire
// shape; nil stays nil.
func protoImageToWire(m map[string]*pb.AttributeValue) map[string]interface{} {
	if m == nil {
		return nil
	}
	return BuildItemWire(protoToAttributeValueMapDirect(m))
}

// streamRecordToProto converts a stream record to its persisted form.
func streamRecordToProto(rec *StreamRecord) *pb.StoredStreamRecord {
	if rec == nil {
		return nil
	}
	return &pb.StoredStreamRecord{
		EventId:      rec.EventID,
		EventName:    string(rec.EventName),
		EventVersion: rec.EventVersion,
		EventSource:  rec.EventSource,
		AwsRegion:    rec.AWSRegion,
		Dynamodb: &pb.StoredStreamRecordData{
			ApproximateCreationDateTime: int64(rec.Dynamodb.ApproximateCreationDateTime),
			Keys:                        wireImageToProto(rec.Dynamodb.Keys),
			NewImage:                    wireImageToProto(rec.Dynamodb.NewImage),
			OldImage:                    wireImageToProto(rec.Dynamodb.OldImage),
			SequenceNumber:              rec.Dynamodb.SequenceNumber,
			SizeBytes:                   int64(rec.Dynamodb.SizeBytes),
			StreamViewType:              rec.Dynamodb.StreamViewType,
		},
		EventSourceArn: rec.EventSourceARN,
		UserIdentity:   streamUserIdentityToProto(rec.UserIdentity),
	}
}

// streamRecordFromProto converts a persisted stream record back to its
// in-memory form.
func streamRecordFromProto(p *pb.StoredStreamRecord) *StreamRecord {
	if p == nil {
		return nil
	}
	rec := &StreamRecord{
		EventID:        p.EventId,
		EventName:      StreamEventName(p.EventName),
		EventVersion:   p.EventVersion,
		EventSource:    p.EventSource,
		AWSRegion:      p.AwsRegion,
		EventSourceARN: p.EventSourceArn,
		UserIdentity:   protoToStreamUserIdentity(p.UserIdentity),
	}
	if p.Dynamodb != nil {
		rec.Dynamodb = StreamRecordData{
			ApproximateCreationDateTime: float64(p.Dynamodb.ApproximateCreationDateTime),
			Keys:                        protoImageToWire(p.Dynamodb.Keys),
			NewImage:                    protoImageToWire(p.Dynamodb.NewImage),
			OldImage:                    protoImageToWire(p.Dynamodb.OldImage),
			SequenceNumber:              p.Dynamodb.SequenceNumber,
			SizeBytes:                   float64(p.Dynamodb.SizeBytes),
			StreamViewType:              p.Dynamodb.StreamViewType,
		}
	}
	return rec
}

// streamCounterToProto converts the per-table sequence allocator state to
// its persisted form.
func streamCounterToProto(c streamSeqCounter) *pb.StreamSequenceCounter {
	return &pb.StreamSequenceCounter{LastSeq: c.LastSeq, TrimmedFloor: c.TrimmedFloor}
}

// protoToStreamCounter converts a persisted sequence allocator state back.
func protoToStreamCounter(p *pb.StreamSequenceCounter) streamSeqCounter {
	if p == nil {
		return streamSeqCounter{}
	}
	return streamSeqCounter{LastSeq: p.LastSeq, TrimmedFloor: p.TrimmedFloor}
}

// journalRecordToProto converts a PITR journal record to its persisted form.
func journalRecordToProto(r *journalRecord) *pb.JournalRecord {
	return &pb.JournalRecord{
		Timestamp:   r.Timestamp,
		Operation:   r.Operation,
		Key:         attributeValueMapToProtoDirect(r.Key),
		BeforeImage: attributeValueMapToProtoDirect(r.BeforeImage),
	}
}

// protoToJournalRecord converts a persisted PITR journal record back.
func protoToJournalRecord(p *pb.JournalRecord) *journalRecord {
	if p == nil {
		return nil
	}
	return &journalRecord{
		Timestamp:   p.Timestamp,
		Operation:   p.Operation,
		Key:         protoToAttributeValueMapDirect(p.Key),
		BeforeImage: protoToAttributeValueMapDirect(p.BeforeImage),
	}
}
