package aesgcm_test

import (
	"encoding/binary"
	"testing"

	"github.com/neostd/go/crypto/aesgcm"
	"github.com/neostd/go/crypto/hashes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test256(t *testing.T) {
	cipher := aesgcm.New256()
	plaintext := []byte("Hello, World!")
	key := []byte("0123456789abcdef0123456789abcdef")

	encrypted, err := cipher.Encrypt(key, plaintext)
	require.NoError(t, err)

	decrypted, err := cipher.Decrypt(key, encrypted)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func Test128WithSHA512(t *testing.T) {
	cipher := aesgcm.New128()
	cipher.Hash = hashes.SHA512
	plaintext := []byte("Hello, World!")
	key := []byte("0123456789abcdef")

	encrypted, err := cipher.Encrypt(key, plaintext)
	require.NoError(t, err)

	decrypted, err := cipher.Decrypt(key, encrypted)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestHeaderValues(t *testing.T) {
	cipher := aesgcm.New256()
	plaintext := []byte("Hello, World!")
	key := []byte("0123456789abcdef0123456789abcdef")

	encrypted, err := cipher.Encrypt(key, plaintext)
	require.NoError(t, err)

	assert.Equal(t, uint16(cipher.Version), binary.LittleEndian.Uint16(encrypted[0:2]))
	assert.Equal(t, uint16(cipher.SaltSize), binary.LittleEndian.Uint16(encrypted[2:4]))
	assert.Equal(t, uint16(cipher.KeySize), binary.LittleEndian.Uint16(encrypted[4:6]))
	assert.Equal(t, uint16(cipher.Hash.Id()), binary.LittleEndian.Uint16(encrypted[6:8]))
	assert.Equal(t, uint16(12), binary.LittleEndian.Uint16(encrypted[8:10]))
	assert.Equal(t, uint16(16), binary.LittleEndian.Uint16(encrypted[10:12]))
	assert.Equal(t, uint32(cipher.Iterations), binary.LittleEndian.Uint32(encrypted[12:16]))
}

func TestMetadataUnsupported(t *testing.T) {
	cipher := aesgcm.New256()
	key := []byte("0123456789abcdef0123456789abcdef")
	_, err := cipher.EncryptWithMetadata(key, []byte("hello"), []byte("meta"))
	require.Error(t, err)
}
