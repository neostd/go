# crypto

## Overview

The crypto module provides several child packages:

- `aescbc` for AES-CBC encryption with PBKDF2 key derivation and HMAC integrity checks
- `aesgcm` for AES-GCM encryption with PBKDF2 key derivation
- `pwhash` for password hashing with Argon2id or PBKDF2
- `shamirsecretsharing` for Shamir secret splitting and recovery
- `hashes` for shared hash identifiers and helpers used by the crypto packages

The root package contains shared interfaces and helpers used by these child modules.

## AesCBC

The `aescbc` package implements AES encryption in CBC mode. It supports both AES-128 and AES-256 encryption
and uses PBKDF2 for key derivation to generate the symmetric key from a password and the hmac key for integrity checks.

It provides methods for encrypting and decrypting data, ensuring that the data is securely handled with proper padding and integrity checks.

There are methods for just encrypting and decrypting data or including additional metadata with the encrypted data.

## Usage

To use `crypto`, import the module in your Go project:

```go
import "github.com/neostd/go/crypto/aescbc"

func main() {
    cipher := aescbc.New256() // or aescbc.New128() for AES-128

    key := []byte("your-secret-key") // For AES-256

    data := []byte("your-data-to-encrypt")
    encryptedData, err := cipher.Encrypt(key, data)
    if err != nil {
        panic(err)
    }

    decryptedData, err := cipher.Decrypt(key, encryptedData)
    if err != nil {
        panic(err)
    }

    fmt.Println("Decrypted data:", string(decryptedData))
    // Output: Decrypted data: your-data-to-encrypt
}

```

## AesGCM

The `aesgcm` package implements AES-GCM encryption with PBKDF2 key derivation.

```go
import "github.com/neostd/go/crypto/aesgcm"

func main() {
    cipher := aesgcm.New256() // or aesgcm.New128() for AES-128

    key := []byte("your-secret-key")
    data := []byte("your-data-to-encrypt")

    encryptedData, err := cipher.Encrypt(key, data)
    if err != nil {
        panic(err)
    }

    decryptedData, err := cipher.Decrypt(key, encryptedData)
    if err != nil {
        panic(err)
    }

    fmt.Println("Decrypted data:", string(decryptedData))
}
```

## PwHash

The `pwhash` package hashes passwords and verifies them using either Argon2id or PBKDF2.

```go
import "github.com/neostd/go/crypto/pwhash"

func main() {
    hasher := pwhash.New()

    encoded, err := hasher.HashString("correct horse battery staple")
    if err != nil {
        panic(err)
    }

    result := hasher.VerifyString("correct horse battery staple", encoded)
    fmt.Println(result == pwhash.ResultSuccess)
}
```

## Shamir Secret Sharing

The `shamir` package splits a secret into shares and reconstructs it from a threshold of those shares.

```go
import "github.com/neostd/go/crypto/shamir"

func main() {
    shares, err := shamir.Split([]byte("secret"), 3, 5)
    if err != nil {
        panic(err)
    }

    recovered, err := shamir.Combine(shares[:3])
    if err != nil {
        panic(err)
    }

    fmt.Println(string(recovered))
}
```

## Hashes

The `hashes` package exposes shared hash constants and helpers used by the encryption packages.

```go
import "github.com/neostd/go/crypto/hashes"

func main() {
    fmt.Println(hashes.SHA256.Id())
    fmt.Println(hashes.SHA256.Size())
}
```
