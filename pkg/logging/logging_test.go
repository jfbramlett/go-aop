package logging

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLogMessage(t *testing.T) {
	t.Run("no_error", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		level := "DEBUG"
		name := "logname"

		logger := logger{ctx: context.Background(), name: name, writer: writer}

		// when
		logger.logMsg(level, msg, nil)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], level)
		assert.Equal(t, logOutput["msg"], msg)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("error", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		level := "DEBUG"
		name := "logname"
		err := errors.New("some error")

		logger := logger{ctx: context.Background(), name: name, writer: writer}

		// when
		logger.logMsg(level, msg, err)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], level)
		assert.Equal(t, logOutput["msg"], msg)
		assert.Equal(t, logOutput["error"], err.Error())
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("mdc", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		level := "DEBUG"
		name := "logname"
		mdcKey := "requestid"
		mdcValue := "12345"

		ctx := AddMDC(context.Background(), map[string]interface{} {mdcKey: mdcValue})

		logger := logger{ctx: ctx, name: name, writer: writer}

		// when
		logger.logMsg(level, msg, nil)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, level, logOutput["level"])
		assert.Equal(t, msg, logOutput["msg"])
		assert.NotNil(t, logOutput["timestamp"])
		assert.Equal(t, mdcValue, logOutput[mdcKey])
	})
}