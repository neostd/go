package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"encoding/binary"
	"fmt"

	crypto2 "github.com/neostd/go/crypto"
	"github.com/neostd/go/crypto/hashes"
)

const fixedHeaderSize = 16

// AesGCM implements AES-GCM encryption with PBKDF2 key derivation.
type AesGCM struct {
	// Iterations is the PBKDF2 iteration count.
	Iterations int32
	// KeySize is the AES key size in bytes.
	KeySize int
	// Version is the encrypted payload format version.
	Version int16
	// SaltSize is the PBKDF2 salt size in bytes.
	SaltSize int16
	// Hash is the PBKDF2 hash algorithm.
	Hash hashes.HashType
}

// New256 returns an AES-GCM cipher configured for AES-256.
func New256() *AesGCM {
	return &AesGCM{
		Iterations: 60000,
		KeySize:    32,
		Version:    1,
		SaltSize:   32,
		Hash:       hashes.SHA256,
	}
}

// New128 returns an AES-GCM cipher configured for AES-128.
func New128() *AesGCM {
	return &AesGCM{
		Iterations: 60000,
		KeySize:    16,
		Version:    1,
		SaltSize:   32,
		Hash:       hashes.SHA256,
	}
}

// Encrypt encrypts data with AES-GCM and includes the derivation parameters in
// the returned payload.
func (a *AesGCM) Encrypt(key []byte, data []byte) ([]byte, error) {
	if a.Version != 1 {
		return nil, fmt.Errorf("unsupported version: %d", a.Version)
	}

	salt, err := crypto2.RandBytes(int(a.SaltSize))
	if err != nil {
		return nil, err
	}

	derivedKey, err := pbkdf2.Key(a.Hash.HashNew(), string(key), salt, int(a.Iterations), a.KeySize)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := crypto2.RandBytes(gcm.NonceSize())
	if err != nil {
		return nil, err
	}

	sealed := gcm.Seal(nil, nonce, data, nil)
	tagSize := gcm.Overhead()
	ciphertextSize := len(sealed) - tagSize
	tagStart := ciphertextSize

	result := make([]byte, fixedHeaderSize+len(salt)+len(nonce)+tagSize+ciphertextSize)
	offset := 0
	binary.LittleEndian.PutUint16(result[offset:], uint16(a.Version))
	offset += 2
	binary.LittleEndian.PutUint16(result[offset:], uint16(a.SaltSize))
	offset += 2
	binary.LittleEndian.PutUint16(result[offset:], uint16(a.KeySize))
	offset += 2
	binary.LittleEndian.PutUint16(result[offset:], uint16(a.Hash.Id()))
	offset += 2
	binary.LittleEndian.PutUint16(result[offset:], uint16(gcm.NonceSize()))
	offset += 2
	binary.LittleEndian.PutUint16(result[offset:], uint16(tagSize))
	offset += 2
	binary.LittleEndian.PutUint32(result[offset:], uint32(a.Iterations))
	offset += 4

	copy(result[offset:], salt)
	offset += len(salt)
	copy(result[offset:], nonce)
	offset += len(nonce)
	copy(result[offset:], sealed[tagStart:])
	offset += tagSize
	copy(result[offset:], sealed[:tagStart])

	return result, nil
}

// Decrypt decrypts an AES-GCM payload using the parameters encoded in the
// payload header.
func (a *AesGCM) Decrypt(key []byte, encryptedData []byte) ([]byte, error) {
	if len(encryptedData) < fixedHeaderSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	offset := 0
	version := int16(binary.LittleEndian.Uint16(encryptedData[offset:]))
	offset += 2
	if version != 1 {
		return nil, fmt.Errorf("unsupported version: %d", version)
	}

	saltSize := int(binary.LittleEndian.Uint16(encryptedData[offset:]))
	offset += 2
	keySize := int(binary.LittleEndian.Uint16(encryptedData[offset:]))
	offset += 2
	hashID := int16(binary.LittleEndian.Uint16(encryptedData[offset:]))
	offset += 2
	nonceSize := int(binary.LittleEndian.Uint16(encryptedData[offset:]))
	offset += 2
	tagSize := int(binary.LittleEndian.Uint16(encryptedData[offset:]))
	offset += 2
	iterations := int(binary.LittleEndian.Uint32(encryptedData[offset:]))
	offset += 4

	hashType := hashes.FromId(hashID)
	if hashType.IsUnknown() {
		return nil, fmt.Errorf("unknown hash id %d", hashID)
	}

	minimumSize := fixedHeaderSize + saltSize + nonceSize + tagSize
	if len(encryptedData) < minimumSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	salt := encryptedData[offset : offset+saltSize]
	offset += saltSize
	nonce := encryptedData[offset : offset+nonceSize]
	offset += nonceSize
	tag := encryptedData[offset : offset+tagSize]
	offset += tagSize
	ciphertext := encryptedData[offset:]

	derivedKey, err := pbkdf2.Key(hashType.HashNew(), string(key), salt, iterations, keySize)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, nonceSize)
	if err != nil {
		return nil, err
	}

	sealed := make([]byte, 0, len(ciphertext)+len(tag))
	sealed = append(sealed, ciphertext...)
	sealed = append(sealed, tag...)

	plaintext, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptWithMetadata encrypts data and rejects metadata because the AES-GCM
// payload format does not currently store metadata.
func (a *AesGCM) EncryptWithMetadata(key []byte, data []byte, metadata []byte) ([]byte, error) {
	if len(metadata) > 0 {
		return nil, fmt.Errorf("metadata is not supported for AesGCM")
	}

	return a.Encrypt(key, data)
}

// DecryptWithMetadata decrypts data and always returns nil metadata.
func (a *AesGCM) DecryptWithMetadata(key []byte, encryptedData []byte) ([]byte, []byte, error) {
	plaintext, err := a.Decrypt(key, encryptedData)
	if err != nil {
		return nil, nil, err
	}

	return plaintext, nil, nil
}
