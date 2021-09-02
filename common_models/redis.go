package common_models

import "time"

type RedisEntity struct {
	Key            string
	Value          interface{}
	ExpirationTime time.Duration
}
