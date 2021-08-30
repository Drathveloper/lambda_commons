package models

import (
	"time"
)

type LambdaContext struct {
	keys map[string]interface{}
}

func NewLambdaContext() LambdaContext {
	return LambdaContext{
		keys: make(map[string]interface{}, 0),
	}
}

func (ctx *LambdaContext) Get(key string) (interface{}, bool) {
	value, exists := ctx.keys[key]
	return value, exists
}

func (ctx *LambdaContext) Set(key string, value interface{}) {
	ctx.keys[key] = value
}

func (ctx *LambdaContext) Exists(key string) bool {
	return ctx.keys[key] != nil
}

func (ctx *LambdaContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ctx *LambdaContext) Done() <-chan struct{} {
	return nil
}

func (ctx *LambdaContext) Err() error {
	return nil
}

func (ctx *LambdaContext) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		if val, exists := ctx.Get(keyAsString); exists {
			return val
		}
	}
	return nil
}
