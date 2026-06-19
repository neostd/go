package crypto

import (
	"crypto/rand"
)

// RandBytes returns securely generated random bytes of the requested size.
func RandBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
