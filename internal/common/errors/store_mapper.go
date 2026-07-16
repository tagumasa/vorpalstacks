package errors

import (
	stderrors "errors"
)

// StoreErrorMapping maps a store-level sentinel error to a service-level
// AWS error. Each service declares a slice of these mappings and calls
// MapStoreError to translate store errors into appropriate API responses.
type StoreErrorMapping struct {
	Store error // The store-level sentinel error to match via errors.Is.
	AWS   error // The AWS error to return when the store error matches.
}

// MapStoreError translates a store error into a service-level AWS error
// using the provided mappings. It returns nil when err is nil and passes
// err through unchanged when no mapping matches. Services that need a
// different fallthrough behaviour (e.g. masking unknown errors as
// InternalFailure) should check the returned error identity after the call.
func MapStoreError(err error, mappings []StoreErrorMapping) error {
	if err == nil {
		return nil
	}
	for _, m := range mappings {
		if stderrors.Is(err, m.Store) {
			return m.AWS
		}
	}
	return err
}
