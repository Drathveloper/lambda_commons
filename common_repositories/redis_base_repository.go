package common_repositories

import (
	"encoding/json"
	"fmt"
	"github.com/Drathveloper/lambda_commons/common_errors"
	"github.com/Drathveloper/lambda_commons/common_models"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisBaseRepository interface {
	Save(ctx *common_models.LambdaContext, redisEntity common_models.RedisEntity) common_errors.GenericApplicationError
	FindKey(ctx *common_models.LambdaContext, key string, value interface{}) (bool, common_errors.GenericApplicationError)
	GetTTL(ctx *common_models.LambdaContext, key string) (bool, time.Duration, common_errors.GenericApplicationError)
	DeleteKey(ctx *common_models.LambdaContext, key string) common_errors.GenericApplicationError
}

type redisBaseRepository struct {
	client    redis.UniversalClient
	namespace string
}

func NewRedisBaseRepository(client redis.UniversalClient, namespace string) RedisBaseRepository {
	return &redisBaseRepository{
		client:    client,
		namespace: namespace,
	}
}

func (repository *redisBaseRepository) Save(ctx *common_models.LambdaContext, redisEntity common_models.RedisEntity) common_errors.GenericApplicationError {
	namespacedKey := fmt.Sprintf("%s:%s", repository.namespace, redisEntity.Key)
	marshaledValue, err := json.Marshal(redisEntity.Value)
	if err != nil {
		return common_errors.NewInternalServerError("error while marshaling value")
	}
	err = repository.client.Set(ctx, namespacedKey, marshaledValue, redisEntity.ExpirationTime).Err()
	if err != nil {
		return common_errors.NewInternalServerError(fmt.Sprintf("error while saving redis key: %s", namespacedKey))
	}
	return nil
}

func (repository *redisBaseRepository) FindKey(ctx *common_models.LambdaContext, key string, value interface{}) (bool, common_errors.GenericApplicationError) {
	namespacedKey := fmt.Sprintf("%s:%s", repository.namespace, key)
	result, err := repository.client.Get(ctx, namespacedKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, common_errors.NewInternalServerError(fmt.Sprintf("error while reading redis key: %s", namespacedKey))
	}
	err = json.Unmarshal([]byte(result), &value)
	if err != nil {
		return true, common_errors.NewInternalServerError(fmt.Sprintf("error while unmarshaling result"))
	}
	return true, nil
}

func (repository *redisBaseRepository) GetTTL(ctx *common_models.LambdaContext, key string) (bool, time.Duration, common_errors.GenericApplicationError) {
	namespacedKey := fmt.Sprintf("%s:%s", repository.namespace, key)
	result, err := repository.client.TTL(ctx, namespacedKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, 0, nil
		}
		return false, 0, common_errors.NewInternalServerError(fmt.Sprintf("error while reading redis key: %s", namespacedKey))
	}
	return true, result, nil
}

func (repository *redisBaseRepository) DeleteKey(ctx *common_models.LambdaContext, key string) common_errors.GenericApplicationError {
	namespacedKey := fmt.Sprintf("%s:%s", repository.namespace, key)
	_, err := repository.client.Del(ctx, namespacedKey).Result()
	if err != nil {
		return common_errors.NewInternalServerError(fmt.Sprintf("error while deleting redis key: %s", namespacedKey))
	}
	return nil
}
