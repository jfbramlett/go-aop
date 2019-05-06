package logging

import (
	"io"
	"os"
	"sync"
)

type LogWriter interface {
	Init()
	WriteString(msg string) (int, error)
	Close()
	Flush()
	IsEnabled(level, method string) bool
}

type channelLogWriter struct {
	outputChannel		chan string
	waitGroup			sync.WaitGroup
	writer				io.StringWriter
}

func (l *channelLogWriter) WriteString(msg string) (int, error) {
	if len(msg) > 0 {
		l.waitGroup.Add(1)
		l.outputChannel <- msg
	}
	return len(msg), nil
}

func (l *channelLogWriter) Init() {
	go func() {
		for  {
			msg := <- l.outputChannel
			if len(msg) > 0 {
				l.writer.WriteString(msg + "\n")
				l.waitGroup.Done()
			}
		}
	}()
}

func (l *channelLogWriter) Close() {
	l.Flush()
	close(l.outputChannel)
}

func (l *channelLogWriter) Flush() {
	l.waitGroup.Wait()
}

func (l *channelLogWriter) IsEnabled(level, method string) bool {
	return true
}


type simpleLogWriter struct {
	writer		io.StringWriter
}

func (l *simpleLogWriter) WriteString(msg string) (int, error) {
	if len(msg) > 0 {
		l.writer.WriteString(msg)
	}
	return len(msg), nil
}

func (l *simpleLogWriter) Init() {
}

func (l *simpleLogWriter) Close() {
	l.Flush()
}

func (l *simpleLogWriter) Flush() {
}

func (l *simpleLogWriter) IsEnabled(level, method string) bool {
	return true
}


var globalWriter LogWriter

func getLogWriter() LogWriter {
	return globalWriter
}

func InitStdoutChannelLogWriter() LogWriter {
	globalWriter = &channelLogWriter{outputChannel: make(chan string), waitGroup: sync.WaitGroup{}, writer: os.Stdout}
	globalWriter.Init()
	return globalWriter
}

func InitChannelLogWriter(writer io.StringWriter) LogWriter {
	globalWriter = &channelLogWriter{outputChannel: make(chan string), waitGroup: sync.WaitGroup{}, writer: writer}
	globalWriter.Init()
	return globalWriter
}

func InitStdoutLogWriter() LogWriter {
	globalWriter = &simpleLogWriter{writer: os.Stdout}
	globalWriter.Init()
	return globalWriter
}

func InitSimpleLogWriter(writer io.StringWriter) LogWriter {
	globalWriter = &simpleLogWriter{writer: writer}
	globalWriter.Init()
	return globalWriter
}

func InitLogWriter(writer LogWriter) LogWriter {
	globalWriter = writer
	globalWriter.Init()
	return globalWriter
}