package pebbledb

import (
	"errors"
	"time"

	"github.com/cockroachdb/pebble/v2"
)

// ForEach iterates over all key-value pairs in the database.
// The iteration stops if the callback returns an error.
func (d *DB) ForEach(fn func(k, v []byte) error) error {
	iter, err := d.NewIter(nil)
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		val, err := iter.ValueAndErr()
		if err != nil {
			return err
		}

		decrypted, err := d.encryptor.Decrypt(val)
		if err != nil {
			return err
		}

		if d.opts.TTL.Enabled {
			ttlVal, err := UnmarshalTTLValue(decrypted)
			if err == nil && ttlVal.ExpiresAt > 0 {
				if time.Now().Unix() > ttlVal.ExpiresAt {
					continue
				}
				decrypted = ttlVal.Data
			}
		}

		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)

		if err := fn(keyCopy, decrypted); err != nil {
			return err
		}
	}

	return iter.Error()
}

// ScanPrefix iterates over all key-value pairs with the given prefix.
// The iteration stops if the callback returns an error.
func (d *DB) ScanPrefix(prefix []byte, fn func(k, v []byte) error) error {
	upperBound := make([]byte, len(prefix)+1)
	copy(upperBound, prefix)
	upperBound[len(prefix)] = 0xFF

	iter, err := d.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: upperBound,
	})
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		val, err := iter.ValueAndErr()
		if err != nil {
			return err
		}

		decrypted, err := d.encryptor.Decrypt(val)
		if err != nil {
			return err
		}

		if d.opts.TTL.Enabled {
			ttlVal, err := UnmarshalTTLValue(decrypted)
			if err == nil && ttlVal.ExpiresAt > 0 {
				if time.Now().Unix() > ttlVal.ExpiresAt {
					continue
				}
				decrypted = ttlVal.Data
			}
		}

		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)

		if err := fn(keyCopy, decrypted); err != nil {
			return err
		}
	}

	return iter.Error()
}

// ScanRange iterates over all key-value pairs in the range [start, end).
// The iteration stops if the callback returns an error.
func (d *DB) ScanRange(start, end []byte, fn func(k, v []byte) error) error {
	iter, err := d.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		val, err := iter.ValueAndErr()
		if err != nil {
			return err
		}

		decrypted, err := d.encryptor.Decrypt(val)
		if err != nil {
			return err
		}

		if d.opts.TTL.Enabled {
			ttlVal, err := UnmarshalTTLValue(decrypted)
			if err == nil && ttlVal.ExpiresAt > 0 {
				if time.Now().Unix() > ttlVal.ExpiresAt {
					continue
				}
				decrypted = ttlVal.Data
			}
		}

		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)

		if err := fn(keyCopy, decrypted); err != nil {
			return err
		}
	}

	return iter.Error()
}

// Count returns the number of keys in the database.
func (d *DB) Count() int {
	count := 0
	err := d.ForEach(func(k, v []byte) error {
		count++
		return nil
	})
	if err != nil {
		return 0
	}
	return count
}

// CountPrefix returns the number of keys with the given prefix.
func (d *DB) CountPrefix(prefix []byte) int {
	count := 0
	_ = d.ScanPrefix(prefix, func(k, v []byte) error {
		count++
		return nil
	})
	return count
}

// ListKeys returns all keys that begin with the given prefix.
// The keys are returned as a slice of byte slices.
func (d *DB) ListKeys(prefix []byte) ([][]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.closed {
		return nil, ErrClosed
	}

	var keys [][]byte

	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: append(prefix, 0xFF),
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)
		keys = append(keys, keyCopy)
	}

	return keys, nil
}

// Has checks if a key exists in the database.
// Returns true if the key exists, false otherwise.
func (d *DB) Has(key []byte) (bool, error) {
	_, err := d.Get(key)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
