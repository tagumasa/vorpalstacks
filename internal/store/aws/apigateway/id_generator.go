package apigateway

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

var idCounter uint64
var idCounterMu sync.Mutex

func generateId(prefix string, length int) string {
	idCounterMu.Lock()
	idCounter++
	counter := idCounter
	idCounterMu.Unlock()

	bytes := make([]byte, length/2+1)
	if _, err := cryptorand.Read(bytes); err != nil {
		for i := range bytes {
			bytes[i] = byte((time.Now().UnixNano() + int64(counter) + int64(i)) % 256)
		}
	}
	return prefix + hex.EncodeToString(bytes)[:length]
}

func generateApiKeyValue() string {
	bytes := make([]byte, 20)
	if _, err := cryptorand.Read(bytes); err != nil {
		idCounterMu.Lock()
		idCounter++
		counter := idCounter
		idCounterMu.Unlock()
		for i := range bytes {
			bytes[i] = byte((time.Now().UnixNano() + int64(counter) + int64(i)) % 256)
		}
	}
	return fmt.Sprintf("%x", bytes)
}
