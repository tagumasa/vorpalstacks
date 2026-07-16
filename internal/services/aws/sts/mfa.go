package sts

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// serialNumberPattern mirrors the Smithy serialNumberType trait:
// [\w+=/:,.@-]* with length 9-256.
var serialNumberPattern = regexp.MustCompile(`^[\w+=/:,.@-]+$`)

// tokenCodePattern mirrors the Smithy tokenCodeType trait: exactly 6 digits.
var tokenCodePattern = regexp.MustCompile(`^\d{6}$`)

// validateMFACredentials checks the format of the SerialNumber and TokenCode
// parameters. SerialNumber must match [\w+=/:,.@-]{9,256} and TokenCode must
// be exactly six digits. An empty pair is valid — MFA is optional unless the
// trust policy requires it. Providing one without the other is invalid.
func validateMFACredentials(serialNumber, tokenCode string) error {
	hasSerial := serialNumber != ""
	hasCode := tokenCode != ""
	if !hasSerial && !hasCode {
		return nil
	}
	if !hasSerial || !hasCode {
		return ErrInvalidMFARequirements
	}
	if len(serialNumber) < 9 || len(serialNumber) > 256 || !serialNumberPattern.MatchString(serialNumber) {
		return ErrInvalidMFASerialNumber
	}
	if !tokenCodePattern.MatchString(tokenCode) {
		return ErrInvalidMFATokenCode
	}
	return nil
}

// verifyTOTP validates a time-based one-time password against a base32-encoded
// shared secret. It accepts codes within a ±1 window (±30 seconds) of the
// current time to allow for minor clock drift between the client and server.
// Returns true when any adjacent window matches.
//
// This is a minimal RFC 6238 TOTP implementation using only the Go standard
// library (no external TOTP dependency required).
func verifyTOTP(base32Seed, code string, now time.Time) bool {
	seed, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(base32Seed))
	if err != nil {
		return false
	}
	counter := now.UTC().Unix() / 30
	for _, offset := range []int64{-1, 0, 1} {
		if computeHOTP(seed, counter+offset) == code {
			return true
		}
	}
	return false
}

// computeHOTP implements RFC 4226 HOTP: HMAC-SHA1 of the 8-byte big-endian
// counter, dynamic truncation, and modulo 10^6 to produce a 6-digit code.
func computeHOTP(secret []byte, counter int64) string {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(counter))
	mac := hmac.New(sha1.New, secret)
	mac.Write(buf[:])
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	truncated := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff
	return fmt.Sprintf("%06d", truncated%1000000)
}
