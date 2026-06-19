package pwhash_test

import (
	"encoding/base64"
	"testing"

	"github.com/neostd/go/crypto/pwhash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPBKDF2HashAndVerify(t *testing.T) {
	h := pwhash.NewWithOptions(pwhash.Options{
		PreferredAlgorithm:  pwhash.AlgorithmPBKDF2,
		SaltLength:          16,
		IVLength:            16,
		HashLength:          32,
		Pbkdf2Iterations:    120000,
		Pbkdf2HashAlgorithm: pwhash.Pbkdf2SHA256,
		Argon2Parameters:    pwhash.DefaultArgon2Parameters(),
	})

	encoded, err := h.HashString("correct horse battery staple")
	require.NoError(t, err)
	assert.Equal(t, pwhash.ResultSuccess, h.VerifyString("correct horse battery staple", encoded))
	assert.Equal(t, pwhash.ResultFailed, h.VerifyString("wrong", encoded))

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)
	assert.Equal(t, pwhash.ResultSuccess, h.Verify([]byte("correct horse battery staple"), decoded))
}

func TestArgon2HashAndNeedsRehash(t *testing.T) {
	h := pwhash.New()
	encoded, err := h.HashString("secret-password")
	require.NoError(t, err)
	assert.Equal(t, pwhash.ResultSuccess, h.VerifyString("secret-password", encoded))

	other := pwhash.NewWithOptions(pwhash.Options{
		PreferredAlgorithm:  pwhash.AlgorithmArgon2ID,
		SaltLength:          16,
		IVLength:            16,
		HashLength:          32,
		Pbkdf2Iterations:    120000,
		Pbkdf2HashAlgorithm: pwhash.Pbkdf2SHA256,
		Argon2Parameters: pwhash.Argon2Parameters{
			Iterations:          4,
			MemorySizeKiB:       65536,
			DegreeOfParallelism: 4,
			TagLength:           32,
		},
	})
	assert.Equal(t, pwhash.ResultNeedsRehash, other.VerifyString("secret-password", encoded))
}
