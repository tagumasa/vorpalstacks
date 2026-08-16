package cloudwatchlogs

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

// PutLookupTable stores a lookup table keyed by its name, which is unique
// per account and Region.
func (s *Store) PutLookupTable(lt *LookupTable) error {
	lt.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.lookupTableKey(lt.Name), lt)
}

// GetLookupTable returns the lookup table with the given name.
func (s *Store) GetLookupTable(name string) (*LookupTable, error) {
	var lt LookupTable
	if err := s.Get(s.lookupTableKey(name), &lt); err != nil {
		return nil, ErrResourceNotFound
	}
	return &lt, nil
}

// DeleteLookupTable removes the lookup table with the given name.
func (s *Store) DeleteLookupTable(name string) error {
	key := "lookup-table:" + name
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

// ListLookupTables returns the lookup tables whose names carry the prefix,
// sorted by name in ascending order.
func (s *Store) ListLookupTables(namePrefix string) ([]*LookupTable, error) {
	var tables []*LookupTable
	if err := s.ScanPrefix("lookup-table:", func(key string, value []byte) error {
		var lt LookupTable
		if err := json.Unmarshal(value, &lt); err != nil {
			return nil
		}
		if namePrefix == "" || strings.HasPrefix(lt.Name, namePrefix) {
			tables = append(tables, &lt)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Slice(tables, func(i, j int) bool { return tables[i].Name < tables[j].Name })
	return tables, nil
}

// CountLookupTables returns the number of stored lookup tables, enforcing
// the per-account, per-Region quota at creation time.
func (s *Store) CountLookupTables() (int, error) {
	tables, err := s.ListLookupTables("")
	if err != nil {
		return 0, err
	}
	return len(tables), nil
}
