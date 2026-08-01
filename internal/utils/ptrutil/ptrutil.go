// Package ptrutil provides generic pointer utility functions for creating
// and dereferencing typed pointers without repetitive boilerplate.
package ptrutil

// Deref returns the value pointed to by p, or fallback if p is nil.
func Deref[T any](p *T, fallback T) T {
	if p != nil {
		return *p
	}
	return fallback
}

// DerefOrZero returns the value pointed to by p, or the zero value of T if p is nil.
func DerefOrZero[T any](p *T) T {
	if p != nil {
		return *p
	}
	var zero T
	return zero
}

// PtrNonZero returns a pointer to v, or nil if v is the zero value.
// This preserves the "omitempty for pointers" pattern used in store JSON
// structs, where a zero-valued field is represented as nil (omitted from
// serialisation) rather than a pointer to zero.
func PtrNonZero[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}
