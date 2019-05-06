package logging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func TestLogDebug(t *testing.T) {
	t.Run("straight", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "logname"

		logger := logger{ctx: context.Background(), name: name, writer: writer}

		// when
		logger.Debug(msg)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], DEBUG)
		assert.Equal(t, logOutput["msg"], msg)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("format", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		msgAddition := "some text"
		expectedMsg := fmt.Sprintf(msg, msgAddition)
		name := "logname"

		logger := logger{ctx: context.Background(), name: name, writer: writer}

		// when
		logger.Debugf(msg, msgAddition)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], DEBUG)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.NotNil(t, logOutput["timestamp"])
	})
}

func TestLogInfo(t *testing.T) {
	t.Run("straight", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "logname"

		logger := logger{ctx: context.Background(), name: name, writer: writer}

		// when
		logger.Info(msg)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], INFO)
		assert.Equal(t, logOutput["msg"], msg)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("format", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		msgAddition := "some text"
		expectedMsg := fmt.Sprintf(msg, msgAddition)
		name := "logname"

		logger := logger{ctx: context.Background(), name: name, writer: writer}

		// when
		logger.Infof(msg, msgAddition)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], INFO)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.NotNil(t, logOutput["timestamp"])
	})
}

func TestLogWarn(t *testing.T) {
	t.Run("straight", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "logname"

		logger := logger{ctx: context.Background(), name: name, writer: writer}

		// when
		logger.Warn(msg)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], WARN)
		assert.Equal(t, logOutput["msg"], msg)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("format", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		msgAddition := "some text"
		expectedMsg := fmt.Sprintf(msg, msgAddition)
		name := "logname"

		logger := logger{ctx: context.Background(), name: name, writer: writer}

		// when
		logger.Warnf(msg, msgAddition)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], WARN)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.NotNil(t, logOutput["timestamp"])
	})
}

func TestLogError(t *testing.T) {
	t.Run("straight", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "logname"
		err := errors.New("my error")

		logger := logger{ctx: context.Background(), name: name, writer: writer}

		// when
		logger.Error(err, msg)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], ERROR)
		assert.Equal(t, logOutput["msg"], msg)
		assert.Equal(t, logOutput["error"], err.Error())
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("format", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		msgAddition := "some text"
		expectedMsg := fmt.Sprintf(msg, msgAddition)
		name := "logname"
		err := errors.New("my error")

		logger := logger{ctx: context.Background(), name: name, writer: writer}

		// when
		logger.Errorf(err, msg, msgAddition)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], ERROR)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.Equal(t, logOutput["error"], err.Error())
		assert.NotNil(t, logOutput["timestamp"])
	})
}

func TestAddMDC(t *testing.T) {
	t.Run("mdc", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		level := "info"
		name := "logname"
		mdcKey := "requestid"
		mdcValue := "12345"
		mdcKey2 := "requestid"
		mdcValue2 := "12345"

		ctx := AddMDC(context.Background(), map[string]interface{} {mdcKey: mdcValue})

		logger := logger{ctx: ctx, name: name, writer: writer}

		// when
		AddMDC(context.Background(), map[string]interface{} {mdcKey2: mdcValue2})
		logger.Info(msg)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, level, logOutput["level"])
		assert.Equal(t, msg, logOutput["msg"])
		assert.NotNil(t, logOutput["timestamp"])
		assert.Equal(t, mdcValue, logOutput[mdcKey])
		assert.Equal(t, mdcValue2, logOutput[mdcKey2])
	})
}