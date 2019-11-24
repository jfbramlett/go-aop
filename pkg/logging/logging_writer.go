package logging

import (
	"io"
	"sync"
)

type LogWriter interface {
	Init()
	WriteString(msg string) (int, error)
	Close()
	Flush()
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
				_, _ = l.writer.WriteString(msg + "\n")
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
		_, _ = l.writer.WriteString(msg)
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


func initChannelLogWriter(writer io.StringWriter) LogWriter {
	channelWriter := &channelLogWriter{outputChannel: make(chan string), waitGroup: sync.WaitGroup{}, writer: writer}
	channelWriter.Init()
	return channelWriter
}

func initSimpleLogWriter(writer io.StringWriter) LogWriter {
	simpleWriter := &simpleLogWriter{writer: writer}
	simpleWriter.Init()
	return simpleWriter
}