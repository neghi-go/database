package mongodb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("Test Client Connection", func(t *testing.T) {
		_, err := New("mongodb://"+test_url, "test-db")
		require.NoError(t, err)

	})

}
