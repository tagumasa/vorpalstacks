package kinesis

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_kinesis"
)

// RegisterStreamConsumer registers a consumer for a Kinesis stream. The
// per-stream quota is enforced inside the locked registration path, so two
// concurrent registrations cannot both pass a count taken outside the lock.
// The consumer record, the stream's count update and the registration-time
// tag set commit as one transaction: a failure at any stage leaves no
// half-registered consumer — and no stream count that believes one exists.
func (s *KinesisStore) RegisterStreamConsumer(streamARN, consumerName string, tags map[string]string) (*Consumer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, err := s.GetStreamByARN(streamARN)
	if err != nil {
		return nil, err
	}

	existing, err := s.ListStreamConsumers(stream.StreamName)
	if err != nil {
		return nil, err
	}
	// Consumer names are unique per stream (the model's ConsumerName
	// documentation); the ARN-existence probe below cannot catch a
	// duplicate name because every registration builds a fresh
	// timestamped ARN, so the name is matched against the registered
	// set here, inside the same lock that guards the quota.
	for _, c := range existing {
		if c.ConsumerName == consumerName {
			return nil, ErrConsumerAlreadyExists
		}
	}
	if len(existing) >= MaxConsumersPerStream {
		return nil, ErrConsumerQuotaExceeded
	}

	// The consumer ARN ends in its own creation timestamp, so the consumer
	// is stamped first and the ARN derives from that stamp — a deregistered
	// name re-registered later is a new generation whose ARN carries the
	// new stamp. The suffix granularity is the one AWS's own documentation
	// example shows (consumer/test-consumer:1525898737 — epoch seconds).
	// With second stamps, a re-registration inside the same second reuses
	// the previous generation's ARN, as AWS's format does: deregistration
	// deletes the record outright, so the reused key carries only the new
	// generation.
	consumer := NewConsumer(consumerName, stream.StreamARN)
	consumerARN := s.buildConsumerARN(stream.StreamName, consumerName, consumer.ConsumerCreationTimestamp.Unix())
	consumer.ConsumerARN = consumerARN
	if s.consumersStore.Exists(consumerARN) {
		return nil, ErrConsumerAlreadyExists
	}

	consumerBytes, err := proto.Marshal(ConsumerToProto(consumer))
	if err != nil {
		return nil, err
	}

	stream.ConsumerCount++
	err = s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(s.consumersBucketName()).Put([]byte(consumerARN), consumerBytes); err != nil {
			return err
		}
		if err := s.updateStreamInTxn(txn, stream); err != nil {
			return err
		}
		// The registration-time tag set commits with the consumer itself;
		// the fresh set's own size carries the create-time cap.
		return s.TagStore.TagInTxn(txn, consumerARN, tags, MaxTagsPerResource)
	})
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// DeregisterStreamConsumer deregisters a consumer from a Kinesis stream.
func (s *KinesisStore) DeregisterStreamConsumer(consumerARN string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var pbConsumer pb.Consumer
	if err := s.consumersStore.GetProto(consumerARN, &pbConsumer); err != nil {
		return ErrConsumerNotFound
	}
	consumer := ProtoToConsumer(&pbConsumer)

	stream, err := s.GetStreamByARN(consumer.StreamARN)
	if err != nil {
		return err
	}

	// The consumer's tag entry dies with the consumer: tags are attached at
	// registration time and have no life apart from the registered consumer.
	// Deleting them first keeps the ordering recoverable — a failure here
	// leaves the consumer registered and the retry re-runs the sweep.
	if err := s.TagStore.Delete(consumerARN); err != nil {
		return fmt.Errorf("failed to delete tags for consumer %s: %w", consumerARN, err)
	}

	if err := s.consumersStore.Delete(consumerARN); err != nil {
		return err
	}

	stream.ConsumerCount--
	if err := s.updateStreamLocked(stream); err != nil {
		return err
	}

	return nil
}

// GetStreamConsumer retrieves a consumer by its ARN.
func (s *KinesisStore) GetStreamConsumer(consumerARN string) (*Consumer, error) {
	var pbConsumer pb.Consumer
	if err := s.consumersStore.GetProto(consumerARN, &pbConsumer); err != nil {
		return nil, ErrConsumerNotFound
	}
	return ProtoToConsumer(&pbConsumer), nil
}

// GetStreamConsumerByName retrieves a consumer by stream ARN and consumer
// name. Consumer ARNs end in their own creation timestamp, so a name
// cannot be derived into an ARN — the stream's registered consumers are
// scanned for the name instead.
func (s *KinesisStore) GetStreamConsumerByName(streamARN, consumerName string) (*Consumer, error) {
	stream, err := s.GetStreamByARN(streamARN)
	if err != nil {
		return nil, err
	}

	consumers, err := s.ListStreamConsumers(stream.StreamName)
	if err != nil {
		return nil, err
	}
	for _, c := range consumers {
		if c.ConsumerName == consumerName {
			return c, nil
		}
	}
	return nil, ErrConsumerNotFound
}

// ListStreamConsumers lists all consumers for a stream.
func (s *KinesisStore) ListStreamConsumers(streamName string) ([]*Consumer, error) {
	stream, err := s.GetStream(streamName)
	if err != nil {
		return nil, err
	}

	prefix := stream.StreamARN + "/"
	var consumers []*Consumer

	err = s.consumersStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var pbConsumer pb.Consumer
		if err := proto.Unmarshal(value, &pbConsumer); err != nil {
			return err
		}
		consumer := ProtoToConsumer(&pbConsumer)
		if consumer.StreamARN == stream.StreamARN {
			consumers = append(consumers, consumer)
		}
		return nil
	})

	return consumers, err
}
