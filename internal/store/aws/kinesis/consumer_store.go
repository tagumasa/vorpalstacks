package kinesis

import (
	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/storage/storage_kinesis"
)

// RegisterStreamConsumer registers a consumer for a Kinesis stream.
func (s *KinesisStore) RegisterStreamConsumer(streamARN, consumerName string) (*Consumer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, err := s.GetStreamByARN(streamARN)
	if err != nil {
		return nil, err
	}

	consumerARN := s.buildConsumerARN(stream.StreamName, consumerName)
	if s.consumersStore.Exists(consumerARN) {
		return nil, ErrConsumerAlreadyExists
	}

	consumer := NewConsumer(consumerName, streamARN)
	consumer.ConsumerARN = consumerARN

	if err := s.consumersStore.PutProto(consumerARN, ConsumerToProto(consumer)); err != nil {
		return nil, err
	}

	stream.ConsumerCount++
	if err := s.UpdateStream(stream); err != nil {
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

	if err := s.consumersStore.Delete(consumerARN); err != nil {
		return err
	}

	stream.ConsumerCount--
	if err := s.UpdateStream(stream); err != nil {
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

// GetStreamConsumerByName retrieves a consumer by stream ARN and consumer name.
func (s *KinesisStore) GetStreamConsumerByName(streamARN, consumerName string) (*Consumer, error) {
	stream, err := s.GetStreamByARN(streamARN)
	if err != nil {
		return nil, err
	}

	consumerARN := s.buildConsumerARN(stream.StreamName, consumerName)
	return s.GetStreamConsumer(consumerARN)
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
