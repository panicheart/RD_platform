package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/oklog/ulid/v2"
)

var ulidEntropy *ulid.MonotonicEntropy

func init() {
	ulidEntropy = ulid.Monotonic(rand.Reader, 0)
}

func GenerateULID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy).String()
}

func GenerateULIDBytes() []byte {
	id := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	return id.Bytes()
}

func ParseULID(s string) (ulid.ULID, error) {
	return ulid.Parse(s)
}

func GenerateRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
