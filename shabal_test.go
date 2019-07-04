package shabal

import (
	"encoding/hex"
	"hash"
	"testing"

	"github.com/stretchr/testify/require"
)

func test(create func() hash.Hash, data, expected string) func(*testing.T) {
	return func(t *testing.T) {
		hasher := create()
		expectedBytes, _ := hex.DecodeString(expected)

		hasher.Write([]byte(data))
		actualBytes := hasher.Sum(nil)
		require.Equal(t, expectedBytes, actualBytes)

		hasher.Write([]byte(data))
		actualBytes = hasher.Sum(nil)
		require.Equal(t, expectedBytes, actualBytes)

		hasher.Write([]byte(data))
		hasher.Reset()
		hasher.Write([]byte(data))
		actualBytes = hasher.Sum(nil)
		require.Equal(t, expectedBytes, actualBytes)
	}
}

//noinspection SpellCheckingInspection
func TestShabal(t *testing.T) {
	t.Run("TestShabal192", test(NewShabal192,
		"abcdefghijklmnopqrstuvwxyz-0123456789-ABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789-abcdefghijklmnopqrstuvwxyz",
		"690fae79226d95760ae8fdb4f58c0537111756557d307b15"))

	t.Run("TestShabal256", test(NewShabal256,
		"helloworld",
		"d945dee21ffca23ac232763aa9cac6c15805f144db9d6c97395437e01c8595a8"))
}
