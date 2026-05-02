// Package idgen genera ULIDs para requestIds y otros identificadores
// monotónicos por servicio. Wrapper sobre oklog/ulid con entropy thread-safe.
package idgen

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	mu      sync.Mutex
	entropy = ulid.Monotonic(rand.Reader, 0)
)

// New devuelve un ULID nuevo en formato Crockford base32 de 26 chars.
func New() string {
	mu.Lock()
	defer mu.Unlock()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
