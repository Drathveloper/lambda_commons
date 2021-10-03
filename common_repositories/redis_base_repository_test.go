package common_repositories_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Drathveloper/lambda_commons/v2/common_errors"
	"github.com/Drathveloper/lambda_commons/v2/common_models"
	"github.com/Drathveloper/lambda_commons/v2/common_repositories"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
	"time"
)

type DummyValue struct {
	Field string `json:"field"`
}

type RedisBaseRepositoryTestSuite struct {
	suite.Suite
	namespace      string
	key            string
	client         redismock.ClientMock
	baseRepository common_repositories.RedisBaseRepository
}

func TestRedisBaseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RedisBaseRepositoryTestSuite))
}

func (suite *RedisBaseRepositoryTestSuite) SetupTest() {
	redisClient, redisClientMock := redismock.NewClientMock()
	suite.client = redisClientMock
	suite.namespace = "dummy-namespace"
	suite.key = "xx"
	suite.baseRepository = common_repositories.NewRedisBaseRepository(redisClient, suite.namespace)
}

func (suite *RedisBaseRepositoryTestSuite) TestSaveShouldSucceed() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)
	duration, _ := time.ParseDuration("1h")
	redisEntity := common_models.RedisEntity{
		Key:            suite.key,
		Value:          `{"field":"value"}`,
		ExpirationTime: duration,
	}
	value, _ := json.Marshal(redisEntity.Value)
	suite.client.ExpectSet(namespacedKey, value, duration).SetVal("")

	err := suite.baseRepository.Save(&ctx, redisEntity)

	suite.NoError(err)
}

func (suite *RedisBaseRepositoryTestSuite) TestSaveShouldReturnErrorWhenMarshalingFailed() {
	ctx := common_models.NewLambdaContext()
	duration, _ := time.ParseDuration("1h")
	redisEntity := common_models.RedisEntity{
		Key:            suite.key,
		Value:          math.Inf(1),
		ExpirationTime: duration,
	}
	expectedErr := common_errors.NewInternalServerError("error while marshaling value")

	err := suite.baseRepository.Save(&ctx, redisEntity)

	suite.Assert().Equal(expectedErr, err)
}

func (suite *RedisBaseRepositoryTestSuite) TestSave_ShouldReturnErrorWhenRedisSetFailed() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)
	duration, _ := time.ParseDuration("1h")
	redisEntity := common_models.RedisEntity{
		Key:            suite.key,
		Value:          `{"field":"value"}`,
		ExpirationTime: duration,
	}
	value, _ := json.Marshal(redisEntity.Value)
	expectedErr := common_errors.NewInternalServerError("error while saving redis key: dummy-namespace:xx")

	suite.client.ExpectSet(namespacedKey, value, duration).SetErr(errors.New("someErr"))

	err := suite.baseRepository.Save(&ctx, redisEntity)

	suite.Assert().Equal(expectedErr, err)
}

func (suite *RedisBaseRepositoryTestSuite) TestFindKeyShouldSucceed() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)
	jsonValue := `{"field":"value"}`
	var value DummyValue

	suite.client.ExpectGet(namespacedKey).SetVal(jsonValue)

	exists, err := suite.baseRepository.FindKey(&ctx, suite.key, &value)

	suite.NoError(err)
	suite.Assert().Equal(true, exists)
	suite.Assert().Equal("value", value.Field)
}

func (suite *RedisBaseRepositoryTestSuite) TestFindKeyShouldReturnFalseWhenRedisKeyDoesNotExist() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)
	var value DummyValue

	suite.client.ExpectGet(namespacedKey).SetErr(errors.New("redis: nil"))

	exists, err := suite.baseRepository.FindKey(&ctx, suite.key, &value)

	suite.NoError(err)
	suite.Assert().Equal(false, exists)
}

func (suite *RedisBaseRepositoryTestSuite) TestFindKeyShouldReturnErrorWhenRedisGetFailed() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)
	var value DummyValue
	expectedErr := common_errors.NewInternalServerError("error while reading redis key: dummy-namespace:xx")

	suite.client.ExpectGet(namespacedKey).SetErr(errors.New("someErr"))

	exists, err := suite.baseRepository.FindKey(&ctx, suite.key, &value)

	suite.Assert().Equal(false, exists)
	suite.Assert().Equal(expectedErr, err)
}

func (suite *RedisBaseRepositoryTestSuite) TestFindKeyShouldReturnErrorWhenUnmarshalGetValueFailed() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)
	jsonValue := `1`
	var value DummyValue
	expectedErr := common_errors.NewInternalServerError("error while unmarshalling result")

	suite.client.ExpectGet(namespacedKey).SetVal(jsonValue)

	exists, err := suite.baseRepository.FindKey(&ctx, suite.key, &value)

	suite.Assert().Equal(true, exists)
	suite.Assert().Equal(expectedErr, err)
}

func (suite *RedisBaseRepositoryTestSuite) TestGetTTLShouldSucceed() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)
	duration, _ := time.ParseDuration("1h")

	suite.client.ExpectTTL(namespacedKey).SetVal(duration)

	exists, ttl, err := suite.baseRepository.GetTTL(&ctx, suite.key)

	suite.NoError(err)
	suite.Equal(true, exists)
	suite.Equal(duration, ttl)
}

func (suite *RedisBaseRepositoryTestSuite) TestGetTTLShouldReturnFalseWhenKeyDoesNotExist() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)
	duration, _ := time.ParseDuration("0")

	suite.client.ExpectTTL(namespacedKey).SetErr(errors.New("redis: nil"))

	exists, ttl, err := suite.baseRepository.GetTTL(&ctx, suite.key)

	suite.NoError(err)
	suite.Equal(false, exists)
	suite.Equal(duration, ttl)
}

func (suite *RedisBaseRepositoryTestSuite) TestGetTTLShouldReturnErrorWhenTTLFailed() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)
	duration, _ := time.ParseDuration("0")
	expectedErr := common_errors.NewInternalServerError("error while reading redis key: dummy-namespace:xx")

	suite.client.ExpectTTL(namespacedKey).SetErr(errors.New("someErr"))

	exists, ttl, err := suite.baseRepository.GetTTL(&ctx, suite.key)

	suite.Equal(expectedErr, err)
	suite.Equal(false, exists)
	suite.Equal(duration, ttl)
}

func (suite *RedisBaseRepositoryTestSuite) TestDeleteKeyShouldSucceed() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)

	suite.client.ExpectDel(namespacedKey).SetVal(0)

	err := suite.baseRepository.DeleteKey(&ctx, suite.key)

	suite.NoError(err)
}

func (suite *RedisBaseRepositoryTestSuite) TestDeleteKeyShouldReturnErrorWhenDelFailed() {
	ctx := common_models.NewLambdaContext()
	namespacedKey := fmt.Sprintf("%s:%s", suite.namespace, suite.key)
	expectedErr := common_errors.NewInternalServerError("error while deleting redis key: dummy-namespace:xx")

	suite.client.ExpectDel(namespacedKey).SetErr(errors.New("someErr"))

	err := suite.baseRepository.DeleteKey(&ctx, suite.key)

	suite.Assert().Equal(expectedErr, err)
}
