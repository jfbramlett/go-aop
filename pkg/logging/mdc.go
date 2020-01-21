package logging

import "context"

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

func getMDC(ctx context.Context) map[string]interface{} {
	current := ctx.Value(mdcCtxKey)
	if current == nil {
		return make(map[string]interface{})
	} else {
		return current.(map[string]interface{})
	}
}
