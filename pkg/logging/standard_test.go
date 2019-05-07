package logging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestLogMessage(t *testing.T) {
	t.Run("no_error", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		level := "DebugLevel"
		name := "logname"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.logMsg(level, msg, nil)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], level)
		assert.Equal(t, logOutput["msg"], msg)
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("error", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		level := "DebugLevel"
		name := "logname"
		err := errors.New("some error")
		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.logMsg(level,  msg, err)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], level)
		assert.Equal(t, logOutput["msg"], msg)
		assert.Equal(t, logOutput["error"], err.Error())
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("mdc", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		level := "DebugLevel"
		name := "logname"
		mdcKey := "requestid"
		mdcValue := "12345"

		ctx := AddMDC(context.Background(), map[string]interface{} {mdcKey: mdcValue})
		logger := logger{ctx: ctx, method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.logMsg(level, msg, nil)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, level, logOutput["level"])
		assert.Equal(t, msg, logOutput["msg"])
		assert.NotNil(t, logOutput["timestamp"])
		assert.Equal(t, logOutput["method"], name)
		assert.Equal(t, mdcValue, logOutput[mdcKey])
	})
}

func TestLogDebug(t *testing.T) {
	t.Run("straight", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.Debug(msg)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], DebugLevel)
		assert.Equal(t, logOutput["msg"], msg)
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("not-enabled", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: LogConfig{EnabledLevels: map[string]bool{}}}

		// when
		logger.Debug(msg)

		// then
		outMsg := writer.String()
		assert.Equal(t, 0, len(outMsg))
	})

	t.Run("not-enabled", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: LogConfig{EnabledLevels: map[string]bool{}}}

		// when
		logger.Debugf(msg, "something to add")

		// then
		outMsg := writer.String()
		assert.Equal(t, 0, len(outMsg))
	})

	t.Run("format", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		msgAddition := "some text"
		expectedMsg := fmt.Sprintf(msg, msgAddition)
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func2"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.Debugf(msg, msgAddition)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], DebugLevel)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})
}

func TestLogInfo(t *testing.T) {
	t.Run("straight", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogInfo.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.Info(msg)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], InfoLevel)
		assert.Equal(t, logOutput["msg"], msg)
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("format", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		msgAddition := "some text"
		expectedMsg := fmt.Sprintf(msg, msgAddition)
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogInfo.func2"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.Infof(msg, msgAddition)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], InfoLevel)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("not-enabled", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: LogConfig{EnabledLevels: map[string]bool{}}}

		// when
		logger.Info(msg)

		// then
		outMsg := writer.String()
		assert.Equal(t, 0, len(outMsg))
	})

	t.Run("not-enabled", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: LogConfig{EnabledLevels: map[string]bool{}}}

		// when
		logger.Infof(msg, "something to add")

		// then
		outMsg := writer.String()
		assert.Equal(t, 0, len(outMsg))
	})
}

func TestLogWarn(t *testing.T) {
	t.Run("straight", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogWarn.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.Warn(msg)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], WarnLevel)
		assert.Equal(t, logOutput["msg"], msg)
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("format", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		msgAddition := "some text"
		expectedMsg := fmt.Sprintf(msg, msgAddition)
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogWarn.func2"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.Warnf(msg, msgAddition)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], WarnLevel)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("not-enabled", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: LogConfig{EnabledLevels: map[string]bool{}}}

		// when
		logger.Warn(msg)

		// then
		outMsg := writer.String()
		assert.Equal(t, 0, len(outMsg))
	})

	t.Run("not-enabled", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: LogConfig{EnabledLevels: map[string]bool{}}}

		// when
		logger.Warnf(msg, "something to add")

		// then
		outMsg := writer.String()
		assert.Equal(t, 0, len(outMsg))
	})
}

func TestLogError(t *testing.T) {
	t.Run("straight", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogError.func1"
		err := errors.New("my error")

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.Error(err, msg)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], ErrorLevel)
		assert.Equal(t, logOutput["msg"], msg)
		assert.Equal(t, logOutput["error"], err.Error())
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("format", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		msgAddition := "some text"
		expectedMsg := fmt.Sprintf(msg, msgAddition)
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogError.func2"
		err := errors.New("my error")

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		logger.Errorf(err, msg, msgAddition)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], ErrorLevel)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.Equal(t, logOutput["error"], err.Error())
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})

	t.Run("not-enabled", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: LogConfig{EnabledLevels: map[string]bool{}}}

		// when
		logger.Error(errors.New("something"), msg)

		// then
		outMsg := writer.String()
		assert.Equal(t, 0, len(outMsg))
	})

	t.Run("not-enabled", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message %s"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func1"

		logger := logger{ctx: context.Background(), method: name, writer: &simpleLogWriter{writer: writer}, config: LogConfig{EnabledLevels: map[string]bool{}}}

		// when
		logger.Errorf(errors.New("something"), msg, "something to add")

		// then
		outMsg := writer.String()
		assert.Equal(t, 0, len(outMsg))
	})
}

func TestAddMDC(t *testing.T) {
	t.Run("mdc", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		level := "info"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestAddMDC.func1"
		mdcKey := "requestid"
		mdcValue := "12345"
		mdcKey2 := "requestid"
		mdcValue2 := "12345"

		ctx := AddMDC(context.Background(), map[string]interface{} {mdcKey: mdcValue})

		logger := logger{ctx: ctx, method: name, writer: &simpleLogWriter{writer: writer}, config: DefaultLogConfig}

		// when
		AddMDC(context.Background(), map[string]interface{} {mdcKey2: mdcValue2})
		logger.Info(msg)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		require.True(t, len(outMsg) > 0)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, level, logOutput["level"])
		assert.Equal(t, msg, logOutput["msg"])
		assert.NotNil(t, logOutput["timestamp"])
		assert.Equal(t, logOutput["method"], name)
		assert.Equal(t, mdcValue, logOutput[mdcKey])
		assert.Equal(t, mdcValue2, logOutput[mdcKey2])
	})
}