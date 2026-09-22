package cloudwatchlogs

import (
	"sort"
	"strings"
	"time"

	"sync"
)

// lookupTableRecordMu serialises read-modify-write cycles on lookup
// table records: the API update and the scheduled-delivery refresh both
// rewrite the record from a read under this mutex, so a refresh cannot
// revert a concurrent update's configuration and vice versa.
var lookupTableRecordMu sync.Mutex

// MutateLookupTable loads the table, applies fn, and persists the result
// as one atomic read-modify-write.
func (s *Store) MutateLookupTable(name string, fn func(*LookupTable) error) error {
	lookupTableRecordMu.Lock()
	defer lookupTableRecordMu.Unlock()
	lt, err := s.GetLookupTable(name)
	if err != nil {
		return err
	}
	if err := fn(lt); err != nil {
		return err
	}
	return s.PutLookupTable(lt)
}

// PutLookupTable stores a lookup table keyed by its name, which is unique
// per account and Region.
func (s *Store) PutLookupTable(lt *LookupTable) error {
	lt.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.lookupTableKey(lt.Name), lt)
}

// GetLookupTable returns the lookup table with the given name.
func (s *Store) GetLookupTable(name string) (*LookupTable, error) {
	return getJSONRecord[LookupTable](s, s.lookupTableKey(name), ErrResourceNotFound)
}

// DeleteLookupTable removes the lookup table with the given name.
func (s *Store) DeleteLookupTable(name string) error {
	key := s.lookupTableKey(name)
	// Under the family mutex so a refresh mid-MutateLookupTable cannot
	// write the record back after the delete (resurrection).
	lookupTableRecordMu.Lock()
	defer lookupTableRecordMu.Unlock()
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

// ListLookupTables returns the lookup tables whose names carry the prefix,
// sorted by name in ascending order.
func (s *Store) ListLookupTables(namePrefix string) ([]*LookupTable, error) {
	tables, err := listJSONRecords[LookupTable](s, keyPrefixLookupTable, "lookup table record", func(lt *LookupTable) bool {
		return namePrefix == "" || strings.HasPrefix(lt.Name, namePrefix)
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(tables, func(i, j int) bool { return tables[i].Name < tables[j].Name })
	return tables, nil
}

// CountLookupTables returns the number of stored lookup tables, enforcing
// the per-account, per-Region quota at creation time. The count is a
// key-only scan — the table bodies (up to the 10 MB payload each) are
// never decoded for a len().
func (s *Store) CountLookupTables() (int, error) {
	count := 0
	err := s.ScanPrefix(keyPrefixLookupTable, func(string, []byte) error {
		count++
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// LookupTable stores the reference data the lookup and cidrlookup query
// commands enrich events with. TableBody holds the CSV content including
// the header row; TableFields mirrors the header and RecordsCount counts
// the data rows. When the table is encrypted with a customer-managed KMS
// key, TableBody is empty and EncryptedBody, EncryptedDataKey and
// ContentNonce hold the envelope-encrypted content instead.
type LookupTable struct {
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	TableBody        string            `json:"tableBody,omitempty"`
	TableFields      []string          `json:"tableFields,omitempty"`
	RecordsCount     int64             `json:"recordsCount"`
	SizeBytes        int64             `json:"sizeBytes"`
	KmsKeyId         string            `json:"kmsKeyId,omitempty"`
	EncryptedBody    []byte            `json:"encryptedBody,omitempty"`
	EncryptedDataKey []byte            `json:"encryptedDataKey,omitempty"`
	ContentNonce     []byte            `json:"contentNonce,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
	CreationTime     int64             `json:"creationTime"`
	LastUpdatedTime  int64             `json:"lastUpdatedTime"`
}
