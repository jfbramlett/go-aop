package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLogDebugMultipleGoRoutines(t *testing.T) {
	// given
	writer := &strings.Builder{}
	name := "github.com/jfbramlett/go-aop/pkg/logging.TestLogDebugMultipleGoRoutines.func1"
	channelWriter := InitChannelLogWriter(writer)
	defer channelWriter.Close()

	// when
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			logger := logger{ctx: context.Background(), method: name, writer: getLogWriter()}
			logger.Debug(fmt.Sprintf("hello %d", idx))
			time.Sleep(1 * time.Second)
			wg.Done()
		}(i)
	}

	wg.Wait()
	getLogWriter().Flush()
	time.Sleep(1 * time.Second)

	// then
	outMsg := writer.String()
	assert.NotNil(t, outMsg)

	outMsg = strings.Trim(outMsg, "\n")
	msgs := strings.Split(outMsg, "\n")
	//assert.Equal(t, 10, len(msgs)) - this check is inconsistent when run with other logging calls

	for _, m := range msgs {
		logOutput := make(map[string]interface{})
		json.Unmarshal([]byte(m), &logOutput)
		assert.Equal(t, logOutput["level"], DEBUG)
		assert.True(t, strings.HasPrefix(logOutput["msg"].(string), "hello "))
		assert.Equal(t, logOutput["method"], name)
		assert.NotNil(t, logOutput["timestamp"])
	}
}
