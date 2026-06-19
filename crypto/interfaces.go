package crypto

// SymmetricCipher describes a symmetric encryption implementation that can
// optionally include metadata in its encrypted payloads.
type SymmetricCipher interface {
	// Encrypt encrypts data with the provided key.
	Encrypt(key []byte, data []byte) (encryptedData []byte, err error)

	// EncryptWithMetadata encrypts data with the provided key and metadata.
	EncryptWithMetadata(key []byte, data []byte, metadata []byte) (encryptedData []byte, err error)

	// Decrypt decrypts an encrypted payload with the provided key.
	Decrypt(key []byte, encryptedData []byte) (data []byte, err error)

	// DecryptWithMetadata decrypts an encrypted payload and returns any metadata.
	DecryptWithMetadata(key []byte, encryptedData []byte) (data []byte, metadata []byte, err error)
}
