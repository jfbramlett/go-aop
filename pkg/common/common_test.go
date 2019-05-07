package common

import (
	"context"
	"github.com/golang-collections/collections/stack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPushToContext(t *testing.T) {
	t.Run("add_to_empty", func(t *testing.T) {
		// given
		key := "ctxKey"
		val := "some value"

		// when
		ctx := PushToContext(context.Background(), key, val)

		// then
		ctxVal := ctx.Value(key)
		require.NotNil(t, ctxVal)

		stackVal, valid := ctxVal.(*stack.Stack)
		require.True(t, valid)
		assert.Equal(t, 1, stackVal.Len())
		assert.Equal(t, val, stackVal.Pop())
	})

	t.Run("value_already_there", func(t *testing.T) {
		// given
		stackVal := stack.New()
		key := "ctxKey"
		val1 := "some value"
		val2 := "some other value"
		stackVal.Push(val1)
		ctx := context.Background()
		ctx = context.WithValue(ctx, key, stackVal)

		// when
		ctx = PushToContext(ctx, key, val2)

		// then
		ctxVal := ctx.Value(key)
		require.NotNil(t, ctxVal)

		stackVal, valid := ctxVal.(*stack.Stack)
		require.True(t, valid)
		assert.Equal(t, 2, stackVal.Len())
		assert.Equal(t, val2, stackVal.Pop())
		assert.Equal(t, val1, stackVal.Pop())
	})
}

func TestPopFromContext(t *testing.T) {
	t.Run("pop_empty", func(t *testing.T) {
		// given
		key := "ctxKey"

		// when
		ctx, val := PopFromContext(context.Background(), key)

		// then
		assert.NotNil(t, ctx)
		assert.Nil(t, val)
	})

	t.Run("pop_value", func(t *testing.T) {
		// given
		stackVal := stack.New()
		key := "ctxKey"
		val1 := "some value"
		val2 := "some other value"
		stackVal.Push(val1)
		stackVal.Push(val2)
		ctx := context.Background()
		ctx = context.WithValue(ctx, key, stackVal)

		// when
		ctx, val := PopFromContext(ctx, key)

		// then
		assert.NotNil(t, ctx)
		assert.Equal(t, val2, val)
	})

}

func TestFromContext(t *testing.T) {
	t.Run("from_empty", func(t *testing.T) {
		// given
		key := "ctxKey"

		// when
		val := FromContext(context.Background(), key)

		// then
		assert.Nil(t, val)
	})

	t.Run("from_value", func(t *testing.T) {
		// given
		stackVal := stack.New()
		key := "ctxKey"
		val1 := "some value"
		val2 := "some other value"
		stackVal.Push(val1)
		stackVal.Push(val2)
		ctx := context.Background()
		ctx = context.WithValue(ctx, key, stackVal)

		// when
		fromVal1 := FromContext(ctx, key)
		fromVal2 := FromContext(ctx, key)

		// then
		assert.Equal(t, fromVal1, fromVal2)
		assert.Equal(t, fromVal1, val2)
	})

}

func TestGetMethodName(t *testing.T) {
	// given
	expectedMethodName := "github.com/jfbramlett/go-aop/pkg/common.TestGetMethodName"

	// when
	methodName := func() string {
		return GetCallingMethodName()
	}()

	// then
	assert.Equal(t, expectedMethodName, methodName)
}

func TestGetMethodNameAt(t *testing.T) {
	// given
	expectedMethodName := "github.com/jfbramlett/go-aop/pkg/common.TestGetMethodNameAt"

	// when
	methodName := func() string {
		return GetMethodNameAt(2)
	}()

	// then
	assert.Equal(t, expectedMethodName, methodName)
}
