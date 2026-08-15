package cloudwatchlogs

import (
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
)

// PutDestination creates a new CloudWatch Logs destination.
func (s *Store) PutDestination(dest *Destination) error {
	key := s.destinationKey(dest.Name)
	if dest.CreationTime == 0 {
		if existing, err := s.GetDestination(dest.Name); err == nil {
			dest.CreationTime = existing.CreationTime
		} else {
			dest.CreationTime = time.Now().UnixMilli()
		}
	}
	return s.PutProto(key, DestinationToProto(dest))
}

// GetDestination retrieves a CloudWatch Logs destination by name.
func (s *Store) GetDestination(name string) (*Destination, error) {
	key := s.destinationKey(name)
	var p pb.Destination
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrDestinationNotFound
	}
	return ProtoToDestination(&p), nil
}

// DeleteDestination deletes a CloudWatch Logs destination by name.
func (s *Store) DeleteDestination(name string) error {
	key := s.destinationKey(name)
	if !s.Exists(key) {
		return ErrDestinationNotFound
	}
	return s.Delete(key)
}

// PutDestinationPolicy sets the resource-based access policy for a CloudWatch Logs destination.
func (s *Store) PutDestinationPolicy(name, accessPolicy string) error {
	dest, err := s.GetDestination(name)
	if err != nil {
		return err
	}
	dest.AccessPolicy = accessPolicy
	return s.PutProto(s.destinationKey(name), DestinationToProto(dest))
}

// ListDestinations returns all CloudWatch Logs destinations, optionally filtered by name prefix.
func (s *Store) ListDestinations(prefix string) ([]*Destination, error) {
	destPrefix := "destination:"
	var destinations []*Destination

	if err := s.ScanPrefix(destPrefix, func(key string, value []byte) error {
		var p pb.Destination
		if err := proto.Unmarshal(value, &p); err != nil {
			return nil
		}
		dest := ProtoToDestination(&p)
		if prefix == "" || strings.HasPrefix(dest.Name, prefix) {
			destinations = append(destinations, dest)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return destinations, nil
}
