// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
)

// Pagination token direction prefixes. The token format is
// base64("<direction>-<offset>") where direction is 'f' (forward) or
// 'b' (backward). The offset is a non-negative integer index into the
// ascending-sorted event list.
const (
	PaginationForward  byte = 'f'
	PaginationBackward byte = 'b'
)

// ErrInvalidPaginationToken is returned when a pagination token cannot be
// decoded or contains an invalid (negative, non-numeric, malformed) offset.
var ErrInvalidPaginationToken = errors.New("invalid pagination token")

// EncodePaginationToken creates a direction-aware pagination token encoding
// the given byte ('f' or 'b') and non-negative offset.
func EncodePaginationToken(direction byte, offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%c-%d", direction, offset)))
}

// ParsePaginationToken decodes a direction-aware pagination token.
// Returns the direction byte ('f' or 'b'), the non-negative offset, and an
// error. An empty token returns direction 'f' and offset 0 without error,
// allowing callers to use it as the default starting position.
func ParsePaginationToken(token string) (byte, int, error) {
	if token == "" {
		return PaginationForward, 0, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, 0, ErrInvalidPaginationToken
	}
	s := string(decoded)
	// Direction-prefixed format: "<dir>-<offset>"
	if len(s) < 3 || (s[0] != PaginationForward && s[0] != PaginationBackward) || s[1] != '-' {
		return 0, 0, ErrInvalidPaginationToken
	}
	offset, err := strconv.Atoi(s[2:])
	if err != nil || offset < 0 {
		return 0, 0, ErrInvalidPaginationToken
	}
	return s[0], offset, nil
}
