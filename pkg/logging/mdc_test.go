package logging

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddMDC2(t *testing.T) {
	t.Run("new-ctx", func(t *testing.T) {
		// given
		key := "somekey"
		value := "somevalue"

		// when
		ctx := AddMDC(context.Background(), map[string]interface{} { key: value})

		// then
		mdc := getMDC(ctx)

		require.NotNil(t, mdc)
		assert.Equal(t, 1, len(mdc))
		assert.Equal(t, value, mdc[key])
	})

	t.Run("add-to", func(t *testing.T) {
		// given
		key := "somekey"
		value := "somevalue"

		key2 := "somekey2"
		value2 := "somevalue2"

		ctx := AddMDC(context.Background(), map[string]interface{} { key: value})

		// when
		newCtx := AddMDC(ctx, map[string]interface{} { key2: value2})

		// then
		require.NotEqual(t, newCtx, ctx)

		mdc := getMDC(newCtx)

		require.NotNil(t, mdc)
		assert.Equal(t, 2, len(mdc))
		assert.Equal(t, value, mdc[key])
		assert.Equal(t, value2, mdc[key2])
	})
}
