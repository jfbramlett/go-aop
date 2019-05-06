package logging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestLogMessage(t *testing.T) {
	t.Run("no_error", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		level := "DEBUG"
		name := "logname"

		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: context.Background(), method: name}

		// when
		logger.logMsg(level, msg, nil)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

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
		level := "DEBUG"
		name := "logname"
		err := errors.New("some error")
		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: context.Background(), method: name}

		// when
		logger.logMsg(level,  msg, err)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

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
		level := "DEBUG"
		name := "logname"
		mdcKey := "requestid"
		mdcValue := "12345"

		ctx := AddMDC(context.Background(), map[string]interface{} {mdcKey: mdcValue})
		InitLogging(writer)
		defer StopLogging()

		logger := logger{ctx: ctx, method: name}

		// when
		logger.logMsg(level, msg, nil)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

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

		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: context.Background(), method: name}

		// when
		logger.Debug(msg)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], DEBUG)
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
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebug.func2"

		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: context.Background(), method: name}

		// when
		logger.Debugf(msg, msgAddition)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], DEBUG)
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

		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: context.Background(), method: name}

		// when
		logger.Info(msg)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], INFO)
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

		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: context.Background(), method: name}

		// when
		logger.Infof(msg, msgAddition)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], INFO)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})
}

func TestLogWarn(t *testing.T) {
	t.Run("straight", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogWarn.func1"

		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: context.Background(), method: name}

		// when
		logger.Warn(msg)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], WARN)
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

		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: context.Background(), method: name}

		// when
		logger.Warnf(msg, msgAddition)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], WARN)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	})
}

func TestLogError(t *testing.T) {
	t.Run("straight", func(t *testing.T) {
		// given
		writer := &strings.Builder{}
		msg := "some message"
		name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogError.func1"
		err := errors.New("my error")

		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: context.Background(), method: name}

		// when
		logger.Error(err, msg)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], ERROR)
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

		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: context.Background(), method: name}

		// when
		logger.Errorf(err, msg, msgAddition)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(outMsg), &logOutput)
		assert.Equal(t, logOutput["level"], ERROR)
		assert.Equal(t, logOutput["msg"], expectedMsg)
		assert.Equal(t, logOutput["error"], err.Error())
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
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

		InitLogging(writer)
		defer StopLogging()
		logger := logger{ctx: ctx, method: name}

		// when
		AddMDC(context.Background(), map[string]interface{} {mdcKey2: mdcValue2})
		logger.Info(msg)
		time.Sleep(1 * time.Second)

		// then
		outMsg := writer.String()
		assert.NotNil(t, outMsg)

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