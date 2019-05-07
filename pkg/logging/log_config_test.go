package logging

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetLogger(t *testing.T) {
	// given
	InitLogging()

	// when
	logger := GetLogger(context.Background())

	// then
	require.NotNil(t, logger)
}
