# Notes

## Argon2 Interop Follow-Up

The new Go `pwhash` package includes Argon2id support and matches the C# `PasswordHasher` structure closely, but Argon2 output is not yet guaranteed to be byte-for-byte interoperable with the C# implementation.

### Current Status

- The encoded password hash record format is aligned with C#:
  - `PBHD` file header
  - version byte
  - algorithm byte
  - PBKDF2 algorithm byte
  - reserved byte
  - iterations
  - memory size
  - parallelism
  - salt length
  - IV length
  - hash length
  - salt bytes
  - IV bytes
  - hash bytes
- PBKDF2 behavior is also aligned closely:
  - hash is derived from `salt || iv`
  - hash metadata is stored in the record
  - verify returns success / failed / needs rehash semantics

### Argon2 Caveat

The C# implementation derives Argon2id using more than just password + salt:

- It combines `KnownSecret` and `iv`
- It also carries `AssociatedData`
- It passes these through its Argon2 implementation as extra Argon2 context inputs

The Go implementation currently uses `golang.org/x/crypto/argon2`, which does not expose the same richer parameter surface in a compatible way.

Because of that:

- Go Argon2 hashes are structurally similar to the C# version
- Go Argon2 verify works for Go-generated hashes
- C# Argon2 and Go Argon2 are not yet guaranteed to produce identical derived bytes for the same inputs

### What To Investigate Later

1. Confirm which exact C# Argon2 library is used underneath `Argon2.DeriveKey(...)`.
2. Verify how `KnownSecret`, `AssociatedData`, and `iv` are fed into the final Argon2 state.
3. Determine whether Go needs:
   - a different Argon2 library than `golang.org/x/crypto/argon2`, or
   - a custom wrapper/implementation to support secret/ad parameters.
4. Create fixed C# test vectors for:
   - password
   - salt
   - iv
   - iterations
   - memory size
   - parallelism
   - output length
   - expected derived bytes
5. Validate Go against those vectors before claiming full Argon2 interop.

### Practical Guidance For Now

- Use PBKDF2 when cross-language password hash compatibility is required immediately.
- Treat Argon2 support as format-compatible but not yet proven output-compatible with the current C# implementation.
