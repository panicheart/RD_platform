package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateULID(t *testing.T) {
	t.Run("should generate valid ULID", func(t *testing.T) {
		id := GenerateULID()
		assert.Len(t, id, 26)
		assert.NotEmpty(t, id)
	})

	t.Run("should generate unique ULIDs", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			id := GenerateULID()
			require.False(t, ids[id], "Generated duplicate ULID: %s", id)
			ids[id] = true
		}
		assert.Len(t, ids, 100)
	})
}

func TestParseULID(t *testing.T) {
	t.Run("should parse valid ULID", func(t *testing.T) {
		id := GenerateULID()
		parsed, err := ParseULID(id)
		require.NoError(t, err)
		assert.Equal(t, id, parsed.String())
	})

	t.Run("should return error for invalid ULID", func(t *testing.T) {
		_, err := ParseULID("invalid")
		assert.Error(t, err)
	})

	t.Run("should return error for empty string", func(t *testing.T) {
		_, err := ParseULID("")
		assert.Error(t, err)
	})
}

func TestGenerateRandomHex(t *testing.T) {
	t.Run("should generate random hex", func(t *testing.T) {
		hex, err := GenerateRandomHex(16)
		require.NoError(t, err)
		assert.Len(t, hex, 32)
	})

	t.Run("should generate different values", func(t *testing.T) {
		hex1, err := GenerateRandomHex(8)
		require.NoError(t, err)
		hex2, err := GenerateRandomHex(8)
		require.NoError(t, err)
		assert.NotEqual(t, hex1, hex2)
	})
}
