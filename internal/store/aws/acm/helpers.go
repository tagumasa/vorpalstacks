package acm

import (
	"crypto/md5"
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/core/logs"
)

const (
	validationTokenLength = 56
	randomBytesPoolSize   = 256
)

var certificateCounter uint64 = 0

// GenerateCertificateId generates a unique certificate ID.
func GenerateCertificateId() string {
	newCount := atomic.AddUint64(&certificateCounter, 1)
	return fmt.Sprintf("%016x", newCount)
}

// GenerateCertificateSerial generates a unique certificate serial number.
func GenerateCertificateSerial() string {
	counter := atomic.LoadUint64(&certificateCounter)
	ts := time.Now().UnixNano()
	hash := md5.Sum([]byte(fmt.Sprintf("%d-%d", counter, ts)))
	return fmt.Sprintf("%x", hash)[:16]
}

func generateValidationToken() string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	const lettersLen = byte(len(letters))
	const maxByte = byte(256 - (256 % uint32(lettersLen)))
	b := make([]byte, validationTokenLength)
	randBytes := make([]byte, randomBytesPoolSize)

	if _, err := cryptorand.Read(randBytes); err != nil {
		logs.Warn("crypto/rand.Read failed in generateValidationToken, using fallback", logs.Err(err))
		for i := range b {
			b[i] = letters[(time.Now().UnixNano()+int64(i))%int64(lettersLen)]
		}
		return string(b)
	}

	randIdx := 0
	for i := range b {
		for randIdx < len(randBytes) {
			if randBytes[randIdx] < maxByte {
				b[i] = letters[randBytes[randIdx]%lettersLen]
				randIdx++
				break
			}
			randIdx++
		}
		if randIdx >= len(randBytes) {
			break
		}
	}
	return string(b)
}

// GenerateDomainValidationRecordName generates the DNS record name for domain validation.
func GenerateDomainValidationRecordName(domain string) string {
	return fmt.Sprintf("_%s.acm-validations.aws.", domain)
}

// GenerateDomainValidationRecordValue generates the DNS record value for domain validation.
func GenerateDomainValidationRecordValue() string {
	return generateValidationToken()
}

func init() {
	var seedBytes [8]byte
	if _, err := cryptorand.Read(seedBytes[:]); err == nil {
		seed := binary.BigEndian.Uint64(seedBytes[:])
		atomic.StoreUint64(&certificateCounter, seed)
	}
}
