package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	totpTimeStep   = 30
	totpDigits     = 6
	totpWindowSize = 1
)

// ValidateConsecutiveTOTPCodes validates two consecutive TOTP codes.
func ValidateConsecutiveTOTPCodes(secret string, code1, code2 string) error {
	if secret == "" {
		return fmt.Errorf("secret is empty")
	}

	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))
	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return fmt.Errorf("failed to decode base32 secret: %w", err)
	}

	now := time.Now().Unix()
	currentTimeStep := now / totpTimeStep

	var foundStep1 int64 = -1
	for offset := -totpWindowSize; offset <= totpWindowSize; offset++ {
		step := currentTimeStep + int64(offset)
		if generateTOTP(decoded, step) == code1 {
			foundStep1 = step
			break
		}
	}

	if foundStep1 == -1 {
		return fmt.Errorf("first authentication code is invalid")
	}

	for offset := -totpWindowSize; offset <= totpWindowSize; offset++ {
		step := currentTimeStep + int64(offset)
		if generateTOTP(decoded, step) == code2 {
			if step > foundStep1 {
				return nil
			}
		}
	}

	return fmt.Errorf("second authentication code is invalid or not after first code")
}

func generateTOTP(secret []byte, timeStep int64) string {
	counter := make([]byte, 8)
	binary.BigEndian.PutUint64(counter, uint64(timeStep))

	h := hmac.New(sha1.New, secret)
	h.Write(counter)
	hmacResult := h.Sum(nil)

	offset := hmacResult[len(hmacResult)-1] & 0x0F
	code := binary.BigEndian.Uint32(hmacResult[offset : offset+4])
	code = code & 0x7FFFFFFF
	code = code % uint32(math.Pow10(totpDigits))

	return fmt.Sprintf("%06d", code)
}
