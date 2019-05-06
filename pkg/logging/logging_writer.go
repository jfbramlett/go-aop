package logging

import (
	"io"
	"sync"
)

type logWriter struct {
	outputChannel		chan string
	waitGroup			sync.WaitGroup
	writer				io.StringWriter
}

func (l *logWriter) WriteString(msg string) (int, error) {
	if len(msg) > 0 {
		l.waitGroup.Add(1)
		l.outputChannel <- msg
	}
	return len(msg), nil
}

func (l *logWriter) Run() {
	for  {
		msg := <- l.outputChannel
		if len(msg) > 0 {
			l.writer.WriteString(msg + "\n")
			l.waitGroup.Done()
		}
	}
}

func (l *logWriter) Close() {
	close(l.outputChannel)
}

func (l *logWriter) Flush() {
	l.waitGroup.Wait()
}

func (l *logWriter) IsEnabled(level, method string) bool {
	return true
}


var globalWriter logWriter

func InitLogging(writer io.StringWriter) {
	globalWriter = logWriter{outputChannel: make(chan string), writer: writer}
	go globalWriter.Run()
}


func Flush() {
	globalWriter.Flush()
}

func StopLogging() {
	globalWriter.Close()
}

func WriteString(msg string) {
	globalWriter.WriteString(msg)
}

func IsEnabled(level, method string) bool {
	return globalWriter.IsEnabled(level, method)
}