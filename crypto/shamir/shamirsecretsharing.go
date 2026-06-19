package shamirsecretsharing

import (
	"crypto/rand"
	"fmt"
)

// Share represents one Shamir secret sharing fragment.
type Share struct {
	// Index is the non-zero share identifier.
	Index byte
	// Value contains the share bytes.
	Value []byte
}

// NewShare constructs a Share with a non-zero index.
func NewShare(index byte, value []byte) Share {
	if index == 0 {
		panic("share index must be non-zero")
	}
	return Share{Index: index, Value: value}
}

// Split splits a secret into shareCount shares with the given threshold.
func Split(secret []byte, threshold, shareCount int) ([]Share, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("secret must not be empty")
	}
	if threshold < 2 {
		return nil, fmt.Errorf("threshold must be at least two")
	}
	if shareCount < threshold {
		return nil, fmt.Errorf("share count must be at least the threshold")
	}
	if shareCount > 255 {
		return nil, fmt.Errorf("share count must not exceed 255")
	}

	coefficients := make([]byte, threshold)
	shares := make([]Share, shareCount)
	for i := range shareCount {
		shares[i] = Share{Index: byte(i + 1), Value: make([]byte, len(secret))}
	}

	for byteIndex := range len(secret) {
		coefficients[0] = secret[byteIndex]
		if _, err := rand.Read(coefficients[1:]); err != nil {
			return nil, err
		}
		for i := range shares {
			shares[i].Value[byteIndex] = evaluate(coefficients, shares[i].Index)
		}
	}

	return shares, nil
}

// Combine recombines shares into the original secret.
func Combine(shares []Share) ([]byte, error) {
	if len(shares) < 2 {
		return nil, fmt.Errorf("at least two shares are required")
	}
	length := len(shares[0].Value)
	seen := [256]bool{}
	for _, share := range shares {
		if share.Index == 0 {
			return nil, fmt.Errorf("share index must be non-zero")
		}
		if seen[share.Index] {
			return nil, fmt.Errorf("duplicate share index")
		}
		seen[share.Index] = true
		if len(share.Value) != length {
			return nil, fmt.Errorf("all shares must have the same length")
		}
	}

	secret := make([]byte, length)
	for byteIndex := range len(secret) {
		value := 0
		for i := range len(shares) {
			xi := int(shares[i].Index)
			basis := 1
			for j := range len(shares) {
				if i == j {
					continue
				}
				xj := int(shares[j].Index)
				basis = multiply(basis, divide(xj, xi^xj))
			}
			value ^= multiply(int(shares[i].Value[byteIndex]), basis)
		}
		secret[byteIndex] = byte(value)
	}

	return secret, nil
}

func evaluate(coefficients []byte, x byte) byte {
	result := 0
	for i := len(coefficients) - 1; i >= 0; i-- {
		result = multiply(result, int(x)) ^ int(coefficients[i])
	}
	return byte(result)
}

func divide(left, right int) int {
	if right == 0 {
		panic("divide by zero")
	}
	if left == 0 {
		return 0
	}
	return exp((log(left) + 255 - log(right)) % 255)
}

func multiply(left, right int) int {
	if left == 0 || right == 0 {
		return 0
	}
	return exp((log(left) + log(right)) % 255)
}

func exp(exponent int) int {
	value := 1
	for range exponent {
		value = multiplyNoLut(value, 0x03)
	}
	return value
}

func log(value int) int {
	current := 1
	for i := range 255 {
		if current == value {
			return i
		}
		current = multiplyNoLut(current, 0x03)
	}
	panic("value is not in GF(256)")
}

func multiplyNoLut(left, right int) int {
	product := 0
	a := left
	b := right
	for b > 0 {
		if (b & 1) != 0 {
			product ^= a
		}
		a <<= 1
		if (a & 0x100) != 0 {
			a ^= 0x11B
		}
		b >>= 1
	}
	return product & 0xFF
}
