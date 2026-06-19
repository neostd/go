package aescbc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/pbkdf2"
	"encoding/binary"
	"fmt"

	"github.com/neostd/go/crypto"
	"github.com/neostd/go/crypto/hashes"
)

// AesCBC implements AES encryption in CBC mode with PKCS#7 padding.
// It supports key derivation using PBKDF2 to generate the actual symmetric key
// for aes and the HMAC key for integrity verification.
//
// An encrypt-then-mac approach is used, where the data is encrypted first,
// and then an HMAC is computed over the ciphertext to ensure integrity.
type AesCBC struct {
	// Iterations is the PBKDF2 iteration count.
	Iterations int32
	// KeySize is the AES key size in bytes.
	KeySize int
	// Version is the encrypted payload format version.
	Version int16
	// KdfSaltSize is the PBKDF2 salt size in bytes.
	KdfSaltSize int16
	// KdfHash is the PBKDF2 hash algorithm.
	KdfHash hashes.HashType
	// HmacHash is the integrity hash algorithm used over ciphertext.
	HmacHash hashes.HashType
}

// New256 returns an AES-CBC cipher configured for AES-256.
func New256() *AesCBC {
	return &AesCBC{
		Iterations:  60000,
		KeySize:     32,
		Version:     1,
		KdfSaltSize: 8,
		KdfHash:     hashes.SHA256,
		HmacHash:    hashes.SHA256,
	}
}

// New128 returns an AES-CBC cipher configured for AES-128.
func New128() *AesCBC {
	return &AesCBC{
		Iterations:  60000,
		KeySize:     16,
		Version:     1,
		KdfSaltSize: 8,
		KdfHash:     hashes.SHA256,
		HmacHash:    hashes.SHA256,
	}
}

// Encrypt encrypts data without additional metadata.
func (a *AesCBC) Encrypt(key []byte, data []byte) (encryptedData []byte, err error) {
	return a.EncryptWithMetadata(key, data, nil)
}

// EncryptWithMetadata encrypts data and includes metadata in the payload.
func (a *AesCBC) EncryptWithMetadata(key []byte, data []byte, metadata []byte) (encryptedData []byte, err error) {
	// 1.  version  (short)
	// 2.  salt size (short)
	// 3.  key size (short)
	// 4.  kdf hash algorithm (short)
	// 5.  hmac hash type (short)
	// 6.  iterations (int)
	// 7.  meta data size (int)
	// 8.  salt (byte[])
	// 9.  iv (byte[])
	// 10. meta data (byte[])
	// 11. tag (byte[])
	// 12. encrypted data (byte[])

	if a.Version != 1 {
		return nil, fmt.Errorf("unsupported version: %d", a.Version)
	}

	saltSize := a.KdfSaltSize
	keySize := int16(a.KeySize)

	// 1. version
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, a.Version)
	if err != nil {
		return nil, err
	}

	// 2. salt size
	err = binary.Write(buf, binary.LittleEndian, a.KdfSaltSize)
	if err != nil {
		return nil, err
	}

	// 3. key size
	err = binary.Write(buf, binary.LittleEndian, keySize)
	if err != nil {
		return nil, err
	}

	// 4. pbkdf2 algo type for symmetric key
	err = binary.Write(buf, binary.LittleEndian, a.KdfHash.Id())
	if err != nil {
		return nil, err
	}

	// 5. hmac hash type for hmac/tag key
	err = binary.Write(buf, binary.LittleEndian, a.HmacHash.Id())
	if err != nil {
		return nil, err
	}

	// 6. iterations
	err = binary.Write(buf, binary.LittleEndian, a.Iterations)
	if err != nil {
		return nil, err
	}

	// 7. metadata size
	metadataSize := int32(len(metadata))
	err = binary.Write(buf, binary.LittleEndian, metadataSize)
	if err != nil {
		return nil, err
	}

	// 8. salt
	salt, err := crypto.RandBytes(int(saltSize))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, salt)
	if err != nil {
		return nil, err
	}

	// 9. iv
	iv, err := crypto.RandBytes(16)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, iv)
	if err != nil {
		return nil, err
	}

	// 10. meta data
	if metadataSize > 0 {
		buf.Write(metadata)
	}

	cdr, err := pbkdf2.Key(a.KdfHash.HashNew(), string(key), salt, int(a.Iterations), a.KeySize)
	if err != nil {
		return nil, err
	}
	paddedData := pad(data)
	ciphertext := make([]byte, len(paddedData))
	c, _ := aes.NewCipher(cdr)
	ctr := cipher.NewCBCEncrypter(c, iv)
	ctr.CryptBlocks(ciphertext, paddedData)

	h := a.HmacHash.NewHmac(cdr)
	h.Write(ciphertext)
	hash := h.Sum(nil)

	bufLen := buf.Len()
	hashLen := len(hash)
	ciphertextLen := len(ciphertext)

	result := make([]byte, bufLen+hashLen+ciphertextLen)
	// 1 - 12
	copy(result, buf.Bytes())

	// 13 - tag/hash
	copy(result[bufLen:], hash)

	// 14 - encrypted data
	copy(result[bufLen+hashLen:], ciphertext)

	return result, nil
}

func pad(in []byte) []byte {
	padding := 16 - (len(in) % 16)
	for i := 0; i < padding; i++ {
		in = append(in, byte(padding))
	}
	return in
}

func unpad(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}

	padding := in[len(in)-1]
	if int(padding) > len(in) || padding > aes.BlockSize {
		return nil
	} else if padding == 0 {
		return nil
	}

	for i := len(in) - 1; i > len(in)-int(padding)-1; i-- {
		if in[i] != padding {
			return nil
		}
	}
	return in[:len(in)-int(padding)]
}

// Decrypt decrypts data without returning payload metadata.
func (a *AesCBC) Decrypt(key []byte, encryptedData []byte) (data []byte, err error) {
	decryptedData, _, err := a.DecryptWithMetadata(key, encryptedData)
	return decryptedData, err
}

// DecryptWithMetadata decrypts data and returns any payload metadata.
func (a *AesCBC) DecryptWithMetadata(key []byte, encryptedData []byte) (data []byte, metadata []byte, err error) {
	// 1.  version  (short) 2
	// 2.  salt size (short) 2
	// 3.  key size (short) 2
	// 4.  key pdk2 hash algorithm (short) 2
	// 5.  hmac hash type (short) 2
	// 6.  iterations (int) 4
	// 7.  meta data size (int) 4
	// 8.  salt (byte[])
	// 9.  iv (byte[])
	// 10. meta data (byte[])
	// 11. tag (byte[])
	// 12. encrypted data (byte[])

	// 1. version
	var version int16
	reader := bytes.NewReader(encryptedData)
	err = binary.Read(reader, binary.LittleEndian, &version)
	if err != nil {
		return nil, nil, err
	}

	if version != a.Version {
		return nil, nil, fmt.Errorf("invalid version %d for Aes256CBC", version)
	}

	// 2. salt size (short)
	var saltSize int16
	err = binary.Read(reader, binary.LittleEndian, &saltSize)
	if err != nil {
		return nil, nil, err
	}

	// 3. key size (short)
	var keySizeShort int16
	err = binary.Read(reader, binary.LittleEndian, &keySizeShort)
	if err != nil {
		return nil, nil, err
	}

	// 4. hash algo (short)
	var kdfHashId int16
	err = binary.Read(reader, binary.LittleEndian, &kdfHashId)
	if err != nil {
		return nil, nil, err
	}

	// 5. tag hash algo (short)
	var hmacHashId int16
	err = binary.Read(reader, binary.LittleEndian, &hmacHashId)
	if err != nil {
		return nil, nil, err
	}

	kdfHash := hashes.FromId(kdfHashId)
	if kdfHash.IsUnknown() {
		return nil, nil, fmt.Errorf("unknown kdf hash id %d", kdfHashId)
	}

	hmacHash := hashes.FromId(hmacHashId)
	if hmacHash.IsUnknown() {
		return nil, nil, fmt.Errorf("unknown hmac hash id %d", hmacHashId)
	}

	// 6. iterations (int)
	var iterations int32
	err = binary.Read(reader, binary.LittleEndian, &iterations)
	if err != nil {
		return nil, nil, err
	}

	// 7. metadata size (int)
	var metadataSize int32
	err = binary.Read(reader, binary.LittleEndian, &metadataSize)
	if err != nil {
		return nil, nil, err
	}

	sliceStart := 18

	// 8. salt
	salt := encryptedData[sliceStart : sliceStart+int(saltSize)]
	sliceStart += int(saltSize)

	// 9. iv
	iv := encryptedData[sliceStart : sliceStart+16]
	sliceStart += 16

	// 10. metadata
	if metadataSize > 0 {
		metadata = encryptedData[sliceStart : sliceStart+int(metadataSize)]
		sliceStart += int(metadataSize)
	}

	// 11. tag/hmac
	hash := encryptedData[sliceStart : sliceStart+hmacHash.Size()]
	sliceStart += len(hash)

	// 12. encrypted data
	ciphertext := encryptedData[sliceStart:]

	cdr, err := pbkdf2.Key(kdfHash.HashNew(), string(key), salt, int(iterations), int(keySizeShort))
	if err != nil {
		return nil, nil, err
	}
	h := hmacHash.NewHmac(cdr)
	h.Write(ciphertext)
	expectedHash := h.Sum(nil)

	if !hmac.Equal(hash, expectedHash) {
		return nil, nil, fmt.Errorf("hash mismatch")
	}

	c, _ := aes.NewCipher(cdr)
	ctr := cipher.NewCBCDecrypter(c, iv)
	plaintext := make([]byte, len(ciphertext))
	ctr.CryptBlocks(plaintext, ciphertext)

	return unpad(plaintext), metadata, nil
}
