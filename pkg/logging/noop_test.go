package logging

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoopLogger(t *testing.T) {
	// given
	logger := newNoopLogger(context.Background(), "some method")

	// when
	logger.Debug("msg")
	logger.Info("msg")
	logger.Warn("msg")
	logger.Error(errors.New("error"), "msg")

	logger.Debugf("msg %s", "text")
	logger.Infof("msg %s", "text")
	logger.Warnf("msg %s", "text")
	logger.Errorf(errors.New("error"), "msg %s", "text")

	// then
	assert.False(t, logger.IsInfoEnabled())
	assert.False(t, logger.IsDebugEnabled())
	assert.False(t, logger.IsWarnEnabled())
	assert.False(t, logger.IsErrorEnabled())
}
