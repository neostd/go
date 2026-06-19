package shamirsecretsharing_test

import (
	"testing"

	shamir "github.com/neostd/go/crypto/shamir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitCombine(t *testing.T) {
	secret := []byte("secret")
	shares, err := shamir.Split(secret, 3, 5)
	require.NoError(t, err)
	require.Len(t, shares, 5)

	recovered, err := shamir.Combine(shares[:3])
	require.NoError(t, err)
	assert.Equal(t, secret, recovered)
}

func TestCombineRejectsDuplicateShare(t *testing.T) {
	shares := []shamir.Share{
		{Index: 1, Value: []byte{1, 2, 3}},
		{Index: 1, Value: []byte{4, 5, 6}},
	}
	_, err := shamir.Combine(shares)
	require.Error(t, err)
}
