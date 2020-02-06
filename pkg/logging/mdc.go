package logging

import (
	"context"
	"fmt"
	"github.com/jfbramlett/go-aop/pkg/common"
)

func ForCurrentMethod() *LoggingBuilder {
	currentMethod := common.GetCallingMethodName()
	name := common.BasicQualifierFromMethod(currentMethod)
	return &LoggingBuilder{fields: make(map[string]interface{}), named: name}
}

func Named(name string) *LoggingBuilder {
	return &LoggingBuilder{fields: make(map[string]interface{}), named: name}
}

type LoggingBuilder struct {
	named	string
	fields 	map[string]interface{}
}

func (l *LoggingBuilder) WithField(key string, value interface{}) *LoggingBuilder {
	l.fields[key] = value
	return l
}

func (l *LoggingBuilder) And(key string, value interface{}) *LoggingBuilder {
	l.fields[key] = value
	return l
}

func (l *LoggingBuilder) Create(ctx context.Context) (Logger, context.Context) {
	newCtx := AddMDC(ctx, l.fields)
	return GetLoggerFor(newCtx, l.named), newCtx
}


func AddMDC(ctx context.Context, vals map[string]interface{}) context.Context {
	current := ctx.Value(mdcCtxKey)
	var currentMdc map[string]interface{}
	if current == nil {
		currentMdc = make(map[string]interface{})
	} else {
		currentMdc = current.(map[string]interface{})
	}
	for k, v := range vals {
		currentMdc[k] = v
	}

	return context.WithValue(ctx, mdcCtxKey, currentMdc)
}

func AddMDCValue(ctx context.Context, key string, value interface{}) context.Context {
	return AddMDC(ctx, map[string]interface{} {key: value})
}

func AddMDCValues(ctx context.Context, values... interface{}) context.Context {
	if len(values) % 2 != 0 {
		values = values[:len(values)-1]
	}
	vals := make(map[string]interface{})
	for i := 1; i < len(values); i = i + 2 {
		key := fmt.Sprintf("%v", values[i-1])
		vals[key] = values[i]
	}

	return AddMDC(ctx, vals)
}

func getMDC(ctx context.Context) map[string]interface{} {
	current := ctx.Value(mdcCtxKey)
	if current == nil {
		return make(map[string]interface{})
	} else {
		return current.(map[string]interface{})
	}
}
