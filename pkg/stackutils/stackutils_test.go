package stackutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMethodName(t *testing.T) {
	// given
	expectedMethodName := "github.com/jfbramlett/go-aop/pkg/stackutils.TestGetMethodName"

	// when
	methodName := func() string {
		return GetCallingMethodName()
	}()

	// then
	assert.Equal(t, expectedMethodName, methodName)
}

func TestGetMethodNameAt(t *testing.T) {
	// given
	expectedMethodName := "github.com/jfbramlett/go-aop/pkg/stackutils.TestGetMethodNameAt"

	// when
	methodName := func() string {
		return GetMethodNameAt(2)
	}()

	// then
	assert.Equal(t, expectedMethodName, methodName)
}

func TestMethodNameFromFullPath(t *testing.T) {
	t.Run("test_full_name", func(t *testing.T) {
		// given
		expectedMethodName := "MyMethod"
		fullPath := "github.com/jfbramlett/go-aop/pkg/metrics." + expectedMethodName

		// when
		methodName := MethodNameFromFullPath(fullPath)

		// then
		assert.Equal(t, expectedMethodName, methodName)
	})

	t.Run("test_malformed_name", func(t *testing.T) {
		// given
		expectedMethodName := "MyMethod"

		// when
		methodName := MethodNameFromFullPath(expectedMethodName)

		// then
		assert.Equal(t, expectedMethodName, methodName)
	})

}

func TestStructNameFromMethodt(t *testing.T) {
	t.Run("test_struct_name", func(t *testing.T) {
		// given
		methodName := "github.com/jfbramlett/go-aop/pkg/aop.(sampleStruct).Method1"
		expectedStructName := "sampleStruct"

		// when
		structName := StructNameFromMethod(methodName)

		// then
		assert.Equal(t, expectedStructName, structName)
	})

	t.Run("test_ptr_struct_name", func(t *testing.T) {
		// given
		methodName := "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).Method1"
		expectedStructName := "sampleStruct"

		// when
		structName := StructNameFromMethod(methodName)

		// then
		assert.Equal(t, expectedStructName, structName)
	})

	t.Run("test_no_struct_name", func(t *testing.T) {
		// given
		methodName := "github.com/jfbramlett/go-aop/pkg/aop.Method1"
		expectedStructName := ""

		// when
		structName := StructNameFromMethod(methodName)

		// then
		assert.Equal(t, expectedStructName, structName)
	})
}
