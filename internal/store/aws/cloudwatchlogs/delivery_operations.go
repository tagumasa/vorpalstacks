package cloudwatchlogs

import (
	"sort"
	"sync"
	"time"
)

// The vended-logs delivery family's four record families. They persist in
// the JSON regime (the record families the storage proto does not carry);
// the dual-persistence unification rider sweeps them with the other
// JSON-regime families.

// --- Delivery Source ---

// PutDeliverySource creates or updates a delivery source. An update keeps
// the record's creation time: the operation overwrites the request-carried
// parameters, and the creation stamp is server-side metadata, not a
// parameter.
func (s *Store) PutDeliverySource(src *DeliverySource) error {
	deliveryRecordMu.Lock()
	defer deliveryRecordMu.Unlock()
	if src.CreationTime == 0 {
		if existing, err := s.GetDeliverySource(src.Name); err == nil {
			src.CreationTime = existing.CreationTime
		} else {
			src.CreationTime = time.Now().UTC().UnixMilli()
		}
	}
	return s.Put(s.deliverySourceKey(src.Name), src)
}

func (s *Store) GetDeliverySource(name string) (*DeliverySource, error) {
	return getJSONRecord[DeliverySource](s, s.deliverySourceKey(name), ErrResourceNotFound)
}

func (s *Store) DeleteDeliverySource(name string) error {
	deliveryRecordMu.Lock()
	defer deliveryRecordMu.Unlock()
	if !s.Exists(s.deliverySourceKey(name)) {
		return ErrResourceNotFound
	}
	return s.Delete(s.deliverySourceKey(name))
}

// ListDeliverySources returns every delivery source ordered by name, the
// deterministic order the paginated listing walks.
func (s *Store) ListDeliverySources() ([]*DeliverySource, error) {
	sources, err := listJSONRecords[DeliverySource](s, keyPrefixDeliverySource, "delivery source record", nil)
	if err != nil {
		return nil, err
	}
	sort.Slice(sources, func(i, j int) bool { return sources[i].Name < sources[j].Name })
	return sources, nil
}

// --- Delivery Destination ---

// PutDeliveryDestination creates or updates a delivery destination,
// preserving the creation time across updates like PutDeliverySource.
func (s *Store) PutDeliveryDestination(dest *DeliveryDestination) error {
	deliveryRecordMu.Lock()
	defer deliveryRecordMu.Unlock()
	if dest.CreationTime == 0 {
		if existing, err := s.GetDeliveryDestination(dest.Name); err == nil {
			dest.CreationTime = existing.CreationTime
		} else {
			dest.CreationTime = time.Now().UTC().UnixMilli()
		}
	}
	return s.Put(s.deliveryDestinationKey(dest.Name), dest)
}

func (s *Store) GetDeliveryDestination(name string) (*DeliveryDestination, error) {
	return getJSONRecord[DeliveryDestination](s, s.deliveryDestinationKey(name), ErrResourceNotFound)
}

func (s *Store) DeleteDeliveryDestination(name string) error {
	deliveryRecordMu.Lock()
	defer deliveryRecordMu.Unlock()
	if !s.Exists(s.deliveryDestinationKey(name)) {
		return ErrResourceNotFound
	}
	return s.Delete(s.deliveryDestinationKey(name))
}

// ListDeliveryDestinations returns every delivery destination ordered by
// name.
func (s *Store) ListDeliveryDestinations() ([]*DeliveryDestination, error) {
	destinations, err := listJSONRecords[DeliveryDestination](s, keyPrefixDeliveryDestination, "delivery destination record", nil)
	if err != nil {
		return nil, err
	}
	sort.Slice(destinations, func(i, j int) bool { return destinations[i].Name < destinations[j].Name })
	return destinations, nil
}

// --- Delivery Destination Policy ---

// PutDeliveryDestinationPolicy records (or replaces) the policy document
// assigned to one delivery destination.
func (s *Store) PutDeliveryDestinationPolicy(name string, policy *DeliveryDestinationPolicy) error {
	policy.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.deliveryDestinationPolicyKey(name), policy)
}

func (s *Store) GetDeliveryDestinationPolicy(name string) (*DeliveryDestinationPolicy, error) {
	return getJSONRecord[DeliveryDestinationPolicy](s, s.deliveryDestinationPolicyKey(name), ErrResourceNotFound)
}

func (s *Store) DeleteDeliveryDestinationPolicy(name string) error {
	if !s.Exists(s.deliveryDestinationPolicyKey(name)) {
		return ErrResourceNotFound
	}
	return s.Delete(s.deliveryDestinationPolicyKey(name))
}

// --- Delivery ---

// deliveryRecordMu serialises read-modify-write cycles on delivery
// records: the API's shaping update, the delivery engine's cursor
// advance and the create-path duplicate gate all decide on the
// persisted record under this mutex, so neither side's write can revert
// the other's fields.
var deliveryRecordMu sync.Mutex

// PutDelivery persists a delivery record. The engine's cursor advances
// through this method on every successful delivery pass.
func (s *Store) PutDelivery(delivery *Delivery) error {
	return s.Put(s.deliveryKey(delivery.Id), delivery)
}

// MutateDelivery loads the delivery, applies fn, and persists the result
// as one atomic read-modify-write: shaping updates merge onto the
// current record (a concurrent cursor advance survives), and cursor
// advances merge onto the current record (a concurrent shaping update
// survives).
func (s *Store) MutateDelivery(id string, fn func(*Delivery) error) error {
	deliveryRecordMu.Lock()
	defer deliveryRecordMu.Unlock()
	delivery, err := s.GetDelivery(id)
	if err != nil {
		return err
	}
	if err := fn(delivery); err != nil {
		return err
	}
	return s.PutDelivery(delivery)
}

// PutDeliveryIfPairAbsent persists a new delivery unless one already
// links the same source and destination, deciding under the record
// mutex so two concurrent creations of the same pair admit exactly one.
// existed reports the duplicate without persisting.
func (s *Store) PutDeliveryIfPairAbsent(delivery *Delivery) (existed bool, err error) {
	deliveryRecordMu.Lock()
	defer deliveryRecordMu.Unlock()
	existing, err := s.ListDeliveries()
	if err != nil {
		return false, err
	}
	for _, d := range existing {
		if d.DeliverySourceName == delivery.DeliverySourceName && d.DeliveryDestinationArn == delivery.DeliveryDestinationArn {
			return true, nil
		}
	}
	return false, s.PutDelivery(delivery)
}

func (s *Store) GetDelivery(id string) (*Delivery, error) {
	return getJSONRecord[Delivery](s, s.deliveryKey(id), ErrResourceNotFound)
}

func (s *Store) DeleteDelivery(id string) error {
	// Under the family mutex so a cursor merge mid-MutateDelivery cannot
	// write the record back after the delete (resurrection).
	deliveryRecordMu.Lock()
	defer deliveryRecordMu.Unlock()
	if !s.Exists(s.deliveryKey(id)) {
		return ErrResourceNotFound
	}
	return s.Delete(s.deliveryKey(id))
}

// ListDeliveries returns every delivery ordered by id.
func (s *Store) ListDeliveries() ([]*Delivery, error) {
	deliveries, err := listJSONRecords[Delivery](s, keyPrefixDelivery, "delivery record", nil)
	if err != nil {
		return nil, err
	}
	sort.Slice(deliveries, func(i, j int) bool { return deliveries[i].Id < deliveries[j].Id })
	return deliveries, nil
}

// DeliverySource represents a vended-logs delivery source: the logical
// object standing for the resource whose logs the delivery engine ships.
// The status members of the model's DeliverySource shape are computed at
// read time (the source is ACTIVE while its underlying resource exists),
// never stored.
type DeliverySource struct {
	Name                        string            `json:"name"`
	Arn                         string            `json:"arn"`
	ResourceArn                 string            `json:"resourceArn"`
	Service                     string            `json:"service"`
	LogType                     string            `json:"logType"`
	DeliverySourceConfiguration map[string]string `json:"deliverySourceConfiguration,omitempty"`
	CreationTime                int64             `json:"creationTime"`
	Tags                        map[string]string `json:"tags,omitempty"`
}

// DeliveryDestination represents a vended-logs delivery destination: the
// logical object standing for the resource that receives the logs (a
// CloudWatch Logs log group or an S3 bucket on this platform).
type DeliveryDestination struct {
	Name                    string            `json:"name"`
	Arn                     string            `json:"arn"`
	DeliveryDestinationType string            `json:"deliveryDestinationType"`
	OutputFormat            string            `json:"outputFormat,omitempty"`
	DestinationResourceArn  string            `json:"destinationResourceArn"`
	CreationTime            int64             `json:"creationTime"`
	Tags                    map[string]string `json:"tags,omitempty"`
}

// DeliveryDestinationPolicy represents the delivery destination policy
// assigned to one delivery destination (the cross-account delivery
// permission document; stored and echoed, the platform being
// single-account).
type DeliveryDestinationPolicy struct {
	PolicyDocument  string `json:"policyDocument"`
	LastUpdatedTime int64  `json:"lastUpdatedTime"`
}

// Delivery represents one vended-logs delivery: the pairing of exactly
// one delivery source and one delivery destination, with the record
// shaping the delivery engine applies.
type Delivery struct {
	Id                      string            `json:"id"`
	Arn                     string            `json:"arn"`
	DeliverySourceName      string            `json:"deliverySourceName"`
	DeliveryDestinationArn  string            `json:"deliveryDestinationArn"`
	DeliveryDestinationType string            `json:"deliveryDestinationType"`
	RecordFields            []string          `json:"recordFields,omitempty"`
	FieldDelimiter          string            `json:"fieldDelimiter,omitempty"`
	S3SuffixPath            string            `json:"s3SuffixPath,omitempty"`
	S3HiveCompatiblePath    bool              `json:"s3HiveCompatiblePath,omitempty"`
	CreationTime            int64             `json:"creationTime"`
	Tags                    map[string]string `json:"tags,omitempty"`
	// CursorTime/CursorStream/CursorDigest hold the delivery cursor: the
	// sort position (event timestamp, log stream, message digest) of the
	// last event the engine delivered. CursorCount holds how many events
	// share that triple and are already delivered — byte-identical
	// events at the same millisecond are legal input and the store keeps
	// each one distinctly, so the triple alone cannot deduplicate them.
	// The next pass re-reads from CursorTime inclusive, skips everything
	// before the triple and the first CursorCount events at it, so
	// delivery is at-least-once across restarts and no event below the
	// cursor re-delivers.
	CursorTime   int64  `json:"cursorTime"`
	CursorStream string `json:"cursorStream,omitempty"`
	CursorDigest string `json:"cursorDigest,omitempty"`
	CursorCount  int    `json:"cursorCount,omitempty"`
	// CursorIngestionMark is the ingestion-time watermark of the
	// below-cursor region: a later pass's late window delivers exactly
	// the events whose timestamp falls below CursorTime and whose
	// ingestion time exceeds this mark. The mark is the scanning pass's
	// start clock, and the group lock's ordering (an invisible put stamps
	// its ingestion time strictly after the listing that could not see
	// it) is what makes that advance sound: no below-cursor event can
	// slip under both the scan and the mark.
	CursorIngestionMark int64 `json:"cursorIngestionMark,omitempty"`
}
