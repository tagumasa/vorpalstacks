package cloudwatchlogs

import (
	"errors"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"sync"
	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
	"vorpalstacks/internal/store/aws/common"
)

// PutDestination creates a new CloudWatch Logs destination.
// destinationRecordMu serialises read-modify-write cycles on
// destination records: the creation-stamp and access-policy
// preservation and the policy write all read the current record under
// this mutex, so a policy write cannot be reverted by a concurrent
// PutDestination and vice versa.
var destinationRecordMu sync.Mutex

func (s *Store) PutDestination(dest *Destination) error {
	destinationRecordMu.Lock()
	defer destinationRecordMu.Unlock()
	return s.putDestinationLocked(dest)
}

func (s *Store) putDestinationLocked(dest *Destination) error {
	key := s.destinationKey(dest.Name)
	existing, err := s.GetDestination(dest.Name)
	switch {
	case err == nil:
		// The access policy is owned by PutDestinationPolicy — PutDestination
		// has no accessPolicy member — so an update that arrives without
		// one preserves whatever that operation installed.
		if dest.AccessPolicy == "" {
			dest.AccessPolicy = existing.AccessPolicy
		}
		if dest.CreationTime == 0 {
			dest.CreationTime = existing.CreationTime
		}
	case errors.Is(err, ErrDestinationNotFound):
		if dest.CreationTime == 0 {
			dest.CreationTime = time.Now().UnixMilli()
		}
	default:
		return err
	}
	return s.PutProto(key, DestinationToProto(dest))
}

// MutateDestination loads the destination, applies fn, and persists the
// result as one atomic read-modify-write under the destination mutex.
func (s *Store) MutateDestination(name string, fn func(*Destination) error) error {
	destinationRecordMu.Lock()
	defer destinationRecordMu.Unlock()
	dest, err := s.GetDestination(name)
	if err != nil {
		return err
	}
	if err := fn(dest); err != nil {
		return err
	}
	return s.putDestinationLocked(dest)
}

// GetDestination retrieves a CloudWatch Logs destination by name. A
// missing record reports the not-found sentinel; any other storage
// failure propagates as a storage error.
func (s *Store) GetDestination(name string) (*Destination, error) {
	key := s.destinationKey(name)
	var p pb.Destination
	if err := s.GetProto(key, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrDestinationNotFound
		}
		return nil, err
	}
	return ProtoToDestination(&p), nil
}

// DeleteDestination deletes a CloudWatch Logs destination by name.
func (s *Store) DeleteDestination(name string) error {
	key := s.destinationKey(name)
	// Under the family mutex so a policy merge mid-MutateDestination (or
	// PutDestinationPolicy) cannot write the record back after the delete.
	destinationRecordMu.Lock()
	defer destinationRecordMu.Unlock()
	if !s.Exists(key) {
		return ErrDestinationNotFound
	}
	return s.Delete(key)
}

// PutDestinationPolicy sets the resource-based access policy for a CloudWatch Logs destination.
func (s *Store) PutDestinationPolicy(name, accessPolicy string) error {
	return s.MutateDestination(name, func(dest *Destination) error {
		dest.AccessPolicy = accessPolicy
		return nil
	})
}

// ListDestinations returns all CloudWatch Logs destinations, optionally filtered by name prefix.
func (s *Store) ListDestinations(prefix string) ([]*Destination, error) {
	destPrefix := keyPrefixDestination
	var destinations []*Destination

	if err := s.ScanPrefix(destPrefix, func(key string, value []byte) error {
		var p pb.Destination
		if err := proto.Unmarshal(value, &p); err != nil {
			// A record that fails to decode is invisible to the caller;
			// log it so corruption is discoverable instead of silently
			// shrinking the list.
			logs.Warn("Corrupt destination record",
				logs.String("key", key), logs.Err(err))
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

// Destination represents a CloudWatch Logs destination (cross-account).
type Destination struct {
	Name         string            `json:"name"`
	ARN          string            `json:"arn"`
	RoleArn      string            `json:"roleArn"`
	TargetArn    string            `json:"targetArn"`
	AccessPolicy string            `json:"accessPolicy"`
	CreationTime int64             `json:"creationTime"`
	Tags         map[string]string `json:"tags,omitempty"`
}
