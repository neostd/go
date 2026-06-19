package pwhash

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash"
	"slices"

	crypto2 "github.com/neostd/go/crypto"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

// Algorithm identifies the password hashing algorithm stored in a hash record.
type Algorithm byte

const (
	// AlgorithmArgon2ID uses Argon2id.
	AlgorithmArgon2ID Algorithm = 1
	// AlgorithmPBKDF2 uses PBKDF2.
	AlgorithmPBKDF2 Algorithm = 2
)

// Pbkdf2Algorithm identifies the digest used by PBKDF2.
type Pbkdf2Algorithm byte

const (
	// Pbkdf2SHA1 uses SHA-1 with PBKDF2.
	Pbkdf2SHA1 Pbkdf2Algorithm = 1
	// Pbkdf2SHA256 uses SHA-256 with PBKDF2.
	Pbkdf2SHA256 Pbkdf2Algorithm = 2
	// Pbkdf2SHA384 uses SHA-384 with PBKDF2.
	Pbkdf2SHA384 Pbkdf2Algorithm = 3
	// Pbkdf2SHA512 uses SHA-512 with PBKDF2.
	Pbkdf2SHA512 Pbkdf2Algorithm = 4
)

// Result reports the outcome of password hash verification.
type Result int

const (
	// ResultFailed indicates verification failed.
	ResultFailed Result = iota
	// ResultSuccess indicates verification succeeded and policy matches.
	ResultSuccess
	// ResultNeedsRehash indicates verification succeeded but policy changed.
	ResultNeedsRehash
)

const (
	headerVersion              = 1
	currentMaxHeaderFieldBytes = 64 * 1024 * 1024
)

var fileHeader = []byte("PBHD")

// Argon2Parameters configures Argon2id password hashing.
type Argon2Parameters struct {
	// Iterations is the number of Argon2 passes.
	Iterations int
	// MemorySizeKiB is the memory cost in KiB.
	MemorySizeKiB int
	// DegreeOfParallelism is the Argon2 lane count.
	DegreeOfParallelism int
	// TagLength is the derived output length in bytes.
	TagLength int
	// AssociatedData holds optional associated data for future compatibility.
	AssociatedData []byte
	// KnownSecret holds optional secret input for future compatibility.
	KnownSecret []byte
}

// DefaultArgon2Parameters returns the default Argon2id settings.
func DefaultArgon2Parameters() Argon2Parameters {
	return Argon2Parameters{
		Iterations:          3,
		MemorySizeKiB:       65536,
		DegreeOfParallelism: 4,
		TagLength:           32,
		AssociatedData:      []byte{},
		KnownSecret:         []byte{},
	}
}

// Options configures password hashing and verification policy.
type Options struct {
	// PreferredAlgorithm is used by Hash when no override is supplied.
	PreferredAlgorithm Algorithm
	// SaltLength is the random salt length in bytes.
	SaltLength int
	// IVLength is the random IV payload length in bytes.
	IVLength int
	// HashLength is the derived hash length in bytes.
	HashLength int
	// Pbkdf2Iterations is the PBKDF2 iteration count.
	Pbkdf2Iterations int
	// Pbkdf2HashAlgorithm is the digest used by PBKDF2.
	Pbkdf2HashAlgorithm Pbkdf2Algorithm
	// Argon2Parameters are used when the preferred algorithm is Argon2id.
	Argon2Parameters Argon2Parameters
}

// DefaultOptions returns the default password hashing policy.
func DefaultOptions() Options {
	return Options{
		PreferredAlgorithm:  AlgorithmArgon2ID,
		SaltLength:          16,
		IVLength:            16,
		HashLength:          32,
		Pbkdf2Iterations:    120000,
		Pbkdf2HashAlgorithm: Pbkdf2SHA256,
		Argon2Parameters:    DefaultArgon2Parameters(),
	}
}

// Hasher hashes and verifies passwords using Argon2id or PBKDF2.
type Hasher struct {
	options Options
}

// New returns a Hasher configured with the default options.
func New() *Hasher {
	return &Hasher{options: DefaultOptions()}
}

// NewWithOptions returns a Hasher configured with explicit options.
func NewWithOptions(options Options) *Hasher {
	copy := cloneOptions(options)
	copy.validate()
	return &Hasher{options: copy}
}

// PreferredAlgorithm returns the default algorithm used for new hashes.
func (h *Hasher) PreferredAlgorithm() Algorithm {
	return h.options.PreferredAlgorithm
}

// HashString hashes a UTF-8 password string into an encoded hash record.
func (h *Hasher) HashString(password string) (string, error) {
	return h.Hash([]byte(password))
}

// Hash hashes raw password bytes using the preferred algorithm.
func (h *Hasher) Hash(password []byte) (string, error) {
	return h.HashWithAlgorithm(password, h.options.PreferredAlgorithm)
}

// HashWithAlgorithm hashes raw password bytes using a specific algorithm.
func (h *Hasher) HashWithAlgorithm(password []byte, algorithm Algorithm) (string, error) {
	h.options.validate()
	if len(password) == 0 {
		return "", errInvalidPassword
	}

	salt, err := crypto2.RandBytes(h.options.SaltLength)
	if err != nil {
		return "", err
	}
	iv, err := crypto2.RandBytes(h.options.IVLength)
	if err != nil {
		return "", err
	}

	var hashBytes []byte
	var record passwordHashRecord
	record.Algorithm = algorithm
	record.Salt = salt
	record.IV = iv
	record.SaltLength = len(salt)
	record.IVLength = len(iv)

	switch algorithm {
	case AlgorithmArgon2ID:
		record.Iterations = h.options.Argon2Parameters.Iterations
		record.MemorySizeKiB = h.options.Argon2Parameters.MemorySizeKiB
		record.Parallelism = h.options.Argon2Parameters.DegreeOfParallelism
		record.HashLength = h.options.HashLength
		hashBytes = deriveArgon2(password, salt, iv, h.options.Argon2Parameters)
	case AlgorithmPBKDF2:
		record.Iterations = h.options.Pbkdf2Iterations
		record.Pbkdf2Algorithm = h.options.Pbkdf2HashAlgorithm
		record.HashLength = h.options.HashLength
		hashBytes = derivePbkdf2(password, salt, iv, h.options.Pbkdf2Iterations, h.options.Pbkdf2HashAlgorithm, h.options.HashLength)
	default:
		return "", errUnsupportedAlgorithm
	}

	record.Hash = hashBytes
	payload := buildBinaryHash(record)
	return base64.StdEncoding.EncodeToString(payload), nil
}

// VerifyString verifies a UTF-8 password string against an encoded hash record.
func (h *Hasher) VerifyString(password string, encoded string) Result {
	return h.Verify([]byte(password), []byte(encoded))
}

// Verify verifies raw password bytes against an encoded or binary hash record.
func (h *Hasher) Verify(password []byte, encoded []byte) Result {
	h.options.validate()
	if len(password) == 0 {
		return ResultFailed
	}

	decoded := decodeHashBytes(encoded)
	if len(decoded) == 0 {
		return ResultFailed
	}

	record, ok := tryParseRecord(decoded)
	if !ok || len(record.Hash) == 0 {
		return ResultFailed
	}

	var expected []byte
	switch record.Algorithm {
	case AlgorithmArgon2ID:
		params := DefaultArgon2Parameters()
		params.Iterations = record.Iterations
		params.MemorySizeKiB = record.MemorySizeKiB
		params.DegreeOfParallelism = record.Parallelism
		params.TagLength = record.HashLength
		expected = deriveArgon2(password, record.Salt, record.IV, params)
	case AlgorithmPBKDF2:
		expected = derivePbkdf2(password, record.Salt, record.IV, record.Iterations, record.Pbkdf2Algorithm, record.HashLength)
	default:
		return ResultFailed
	}

	if !slices.Equal(expected, record.Hash) {
		return ResultFailed
	}

	if h.isCurrentPolicy(record) {
		return ResultSuccess
	}

	return ResultNeedsRehash
}

type passwordHashRecord struct {
	Algorithm       Algorithm
	Pbkdf2Algorithm Pbkdf2Algorithm
	Iterations      int
	MemorySizeKiB   int
	Parallelism     int
	SaltLength      int
	IVLength        int
	HashLength      int
	Salt            []byte
	IV              []byte
	Hash            []byte
}

var (
	errInvalidPassword      = errors.New("password must contain bytes")
	errUnsupportedAlgorithm = errors.New("unsupported password hashing algorithm")
)

func cloneOptions(options Options) Options {
	copy := options
	copy.Argon2Parameters.AssociatedData = slices.Clone(options.Argon2Parameters.AssociatedData)
	copy.Argon2Parameters.KnownSecret = slices.Clone(options.Argon2Parameters.KnownSecret)
	if copy.PreferredAlgorithm == 0 {
		copy.PreferredAlgorithm = AlgorithmArgon2ID
	}
	if copy.SaltLength == 0 && copy.IVLength == 0 && copy.HashLength == 0 && copy.Pbkdf2Iterations == 0 && copy.Pbkdf2HashAlgorithm == 0 && copy.Argon2Parameters.Iterations == 0 {
		return DefaultOptions()
	}
	if copy.Argon2Parameters.Iterations == 0 {
		copy.Argon2Parameters = DefaultArgon2Parameters()
	}
	return copy
}

func (o Options) validate() {
	if o.SaltLength <= 0 || o.IVLength <= 0 || o.HashLength <= 0 || o.Pbkdf2Iterations <= 0 {
		panic("invalid password hasher options")
	}
	if o.Pbkdf2HashAlgorithm < Pbkdf2SHA1 || o.Pbkdf2HashAlgorithm > Pbkdf2SHA512 {
		panic("invalid PBKDF2 algorithm")
	}
	if o.PreferredAlgorithm != AlgorithmArgon2ID && o.PreferredAlgorithm != AlgorithmPBKDF2 {
		panic("invalid preferred algorithm")
	}
	if o.Argon2Parameters.Iterations <= 0 || o.Argon2Parameters.DegreeOfParallelism <= 0 || o.Argon2Parameters.MemorySizeKiB < 8*o.Argon2Parameters.DegreeOfParallelism || o.Argon2Parameters.TagLength <= 0 {
		panic("invalid argon2 parameters")
	}
	if o.Argon2Parameters.TagLength != o.HashLength {
		o.Argon2Parameters.TagLength = o.HashLength
	}
}

func buildBinaryHash(record passwordHashRecord) []byte {
	total := len(fileHeader) + 4 + (4 * 6) + len(record.Salt) + len(record.IV) + len(record.Hash)
	out := make([]byte, total)
	index := 0
	copy(out[index:], fileHeader)
	index += len(fileHeader)
	out[index] = headerVersion
	index++
	out[index] = byte(record.Algorithm)
	index++
	out[index] = byte(record.Pbkdf2Algorithm)
	index++
	out[index] = 0
	index++
	binary.LittleEndian.PutUint32(out[index:], uint32(record.Iterations))
	index += 4
	binary.LittleEndian.PutUint32(out[index:], uint32(record.MemorySizeKiB))
	index += 4
	binary.LittleEndian.PutUint32(out[index:], uint32(record.Parallelism))
	index += 4
	binary.LittleEndian.PutUint32(out[index:], uint32(len(record.Salt)))
	index += 4
	binary.LittleEndian.PutUint32(out[index:], uint32(len(record.IV)))
	index += 4
	binary.LittleEndian.PutUint32(out[index:], uint32(len(record.Hash)))
	index += 4
	copy(out[index:], record.Salt)
	index += len(record.Salt)
	copy(out[index:], record.IV)
	index += len(record.IV)
	copy(out[index:], record.Hash)
	return out
}

func decodeHashBytes(encoded []byte) []byte {
	if len(encoded) > len(fileHeader) && bytes.Equal(encoded[:len(fileHeader)], fileHeader) {
		return slices.Clone(encoded)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		return nil
	}
	return decoded
}

func tryParseRecord(payload []byte) (passwordHashRecord, bool) {
	var record passwordHashRecord
	if len(payload) < len(fileHeader)+4+(4*6) || !bytes.Equal(payload[:len(fileHeader)], fileHeader) {
		return record, false
	}
	index := len(fileHeader)
	if payload[index] != headerVersion {
		return record, false
	}
	index++
	record.Algorithm = Algorithm(payload[index])
	index++
	record.Pbkdf2Algorithm = Pbkdf2Algorithm(payload[index])
	index += 2
	record.Iterations = int(binary.LittleEndian.Uint32(payload[index:]))
	index += 4
	record.MemorySizeKiB = int(binary.LittleEndian.Uint32(payload[index:]))
	index += 4
	record.Parallelism = int(binary.LittleEndian.Uint32(payload[index:]))
	index += 4
	record.SaltLength = int(binary.LittleEndian.Uint32(payload[index:]))
	index += 4
	record.IVLength = int(binary.LittleEndian.Uint32(payload[index:]))
	index += 4
	record.HashLength = int(binary.LittleEndian.Uint32(payload[index:]))
	index += 4

	if record.Iterations <= 0 || record.SaltLength <= 0 || record.IVLength <= 0 || record.HashLength <= 0 || record.MemorySizeKiB < 0 || record.Parallelism < 0 {
		return record, false
	}
	if !isLegalLength(record.SaltLength) || !isLegalLength(record.IVLength) || !isLegalLength(record.HashLength) {
		return record, false
	}
	if record.Algorithm == AlgorithmArgon2ID && (record.MemorySizeKiB <= 0 || record.Parallelism <= 0) {
		return record, false
	}
	if record.Algorithm == AlgorithmPBKDF2 && record.Parallelism != 0 {
		return record, false
	}
	required := index + record.SaltLength + record.IVLength + record.HashLength
	if len(payload) != required {
		return record, false
	}
	record.Salt = slices.Clone(payload[index : index+record.SaltLength])
	index += record.SaltLength
	record.IV = slices.Clone(payload[index : index+record.IVLength])
	index += record.IVLength
	record.Hash = slices.Clone(payload[index : index+record.HashLength])
	return record, true
}

func isLegalLength(value int) bool {
	return value > 0 && value <= currentMaxHeaderFieldBytes
}

func deriveArgon2(password, salt, iv []byte, params Argon2Parameters) []byte {
	combinedSalt := make([]byte, 0, len(salt)+len(iv))
	combinedSalt = append(combinedSalt, salt...)
	combinedSalt = append(combinedSalt, iv...)
	passwordWithSecret := make([]byte, 0, len(password)+len(params.KnownSecret))
	passwordWithSecret = append(passwordWithSecret, password...)
	passwordWithSecret = append(passwordWithSecret, params.KnownSecret...)
	return argon2.IDKey(passwordWithSecret, combinedSalt, uint32(params.Iterations), uint32(params.MemorySizeKiB), uint8(params.DegreeOfParallelism), uint32(params.TagLength))
}

func derivePbkdf2(password, salt, iv []byte, iterations int, algorithm Pbkdf2Algorithm, hashLength int) []byte {
	saltWithIV := make([]byte, 0, len(salt)+len(iv))
	saltWithIV = append(saltWithIV, salt...)
	saltWithIV = append(saltWithIV, iv...)
	return pbkdf2.Key(password, saltWithIV, iterations, hashLength, resolvePbkdf2Hash(algorithm))
}

func resolvePbkdf2Hash(algorithm Pbkdf2Algorithm) func() hash.Hash {
	switch algorithm {
	case Pbkdf2SHA1:
		return sha1.New
	case Pbkdf2SHA384:
		return sha512.New384
	case Pbkdf2SHA512:
		return sha512.New
	default:
		return sha256.New
	}
}

func (h *Hasher) isCurrentPolicy(record passwordHashRecord) bool {
	if record.Algorithm != h.options.PreferredAlgorithm {
		return false
	}
	if record.Algorithm == AlgorithmArgon2ID {
		params := h.options.Argon2Parameters
		return record.HashLength == h.options.HashLength &&
			record.Iterations == params.Iterations &&
			record.MemorySizeKiB == params.MemorySizeKiB &&
			record.Parallelism == params.DegreeOfParallelism &&
			record.SaltLength == h.options.SaltLength &&
			record.IVLength == h.options.IVLength
	}

	return record.HashLength == h.options.HashLength &&
		record.Iterations == h.options.Pbkdf2Iterations &&
		record.Pbkdf2Algorithm == h.options.Pbkdf2HashAlgorithm &&
		record.SaltLength == h.options.SaltLength &&
		record.IVLength == h.options.IVLength
}
