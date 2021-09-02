package common_repositories_test

import (
	"errors"
	"github.com/Drathveloper/lambda_commons/common_constants"
	"github.com/Drathveloper/lambda_commons/common_errors"
	"github.com/Drathveloper/lambda_commons/common_models"
	"github.com/Drathveloper/lambda_commons/common_repositories"
	"github.com/Drathveloper/lambda_commons/mocks"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DummyItem struct {
	Key1 string `dynamodbav:"key1"`
	Key2 string `dynamodbav:"key2"`
}

type DynamodbBaseRepositoryTestSuite struct {
	suite.Suite
	dynamodbClient *mocks.MockDynamodbClientAPI
	baseRepository common_repositories.DynamodbBaseRepository
}

func TestDynamodbBaseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DynamodbBaseRepositoryTestSuite))
}

func (suite *DynamodbBaseRepositoryTestSuite) SetupTest() {
	controller := gomock.NewController(suite.T())
	suite.dynamodbClient = mocks.NewMockDynamodbClientAPI(controller)
	suite.baseRepository = common_repositories.NewDynamodbBaseRepository(suite.dynamodbClient, "someTable")
}

func (suite *DynamodbBaseRepositoryTestSuite) TestFindBySimplePrimaryKey_ShouldSucceedWhenNoTransaction() {
	context := common_models.NewLambdaContext()
	primaryKey := common_models.DynamodbSimplePrimaryKey{
		KeyName: "someKey",
		Value:   "someValue",
	}
	getItemInput := dynamodb.GetItemInput{
		TableName:      aws.String("someTable"),
		ConsistentRead: aws.Bool(false),
		Key: map[string]types.AttributeValue{
			"someKey": &types.AttributeValueMemberS{
				Value: "someValue",
			},
		},
	}
	item := map[string]types.AttributeValue{
		"someKey": &types.AttributeValueMemberS{
			Value: "someValue",
		},
		"anotherKey": &types.AttributeValueMemberS{
			Value: "anotherValue",
		},
	}
	getItemOutput := &dynamodb.GetItemOutput{
		Item: item,
	}
	suite.dynamodbClient.EXPECT().GetItem(&context, &getItemInput).Return(getItemOutput, nil)

	response, appErr := suite.baseRepository.FindBySimplePrimaryKey(&context, primaryKey, false, false)

	suite.NoError(appErr)
	suite.Equal(response, item)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestFindBySimplePrimaryKey_ShouldReturnInternalServerErrorWhenGetItemFailed() {
	context := common_models.NewLambdaContext()
	primaryKey := common_models.DynamodbSimplePrimaryKey{
		KeyName: "someKey",
		Value:   "someValue",
	}
	getItemInput := dynamodb.GetItemInput{
		TableName:      aws.String("someTable"),
		ConsistentRead: aws.Bool(false),
		Key: map[string]types.AttributeValue{
			"someKey": &types.AttributeValueMemberS{
				Value: "someValue",
			},
		},
	}
	getItemOutput := &dynamodb.GetItemOutput{}
	cause := errors.New("someErr")
	expectedAppErr := common_errors.NewInternalServerError("error while reading from database")

	suite.dynamodbClient.EXPECT().GetItem(&context, &getItemInput).Return(getItemOutput, cause)

	_, appErr := suite.baseRepository.FindBySimplePrimaryKey(&context, primaryKey, false, false)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestFindBySimplePrimaryKey_ShouldSucceedWhenTransaction() {
	context := common_models.NewLambdaContext()
	transactGetItemsInput := dynamodb.TransactGetItemsInput{
		TransactItems: []types.TransactGetItem{},
	}
	context.Set(common_constants.ReadTransaction, transactGetItemsInput)
	expectedContext := common_models.NewLambdaContext()
	expectedGetItemsInput := dynamodb.TransactGetItemsInput{
		TransactItems: []types.TransactGetItem{
			{
				Get: &types.Get{
					TableName: aws.String("someTable"),
					Key: map[string]types.AttributeValue{
						"someKey": &types.AttributeValueMemberS{
							Value: "someValue",
						},
					},
				},
			},
		},
	}
	expectedContext.Set(common_constants.ReadTransaction, expectedGetItemsInput)
	primaryKey := common_models.DynamodbSimplePrimaryKey{
		KeyName: "someKey",
		Value:   "someValue",
	}

	_, appErr := suite.baseRepository.FindBySimplePrimaryKey(&context, primaryKey, false, true)
	actualGetItemsInput, exists := context.Get(common_constants.ReadTransaction)

	suite.NoError(appErr)
	suite.True(exists)
	suite.Equal(expectedGetItemsInput, actualGetItemsInput)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestFindByComplexPrimaryKey_ShouldSucceedWhenNoTransaction() {
	context := common_models.NewLambdaContext()
	primaryKey := common_models.DynamodbComplexPrimaryKey{
		PartitionKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "somePartitionKey",
			Value:   "somePartitionValue",
		},
		SortKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "someSortKey",
			Value:   "someSortValue",
		},
	}
	getItemInput := dynamodb.GetItemInput{
		TableName:      aws.String("someTable"),
		ConsistentRead: aws.Bool(false),
		Key: map[string]types.AttributeValue{
			"somePartitionKey": &types.AttributeValueMemberS{
				Value: "somePartitionValue",
			},
			"someSortKey": &types.AttributeValueMemberS{
				Value: "someSortValue",
			},
		},
	}
	item := map[string]types.AttributeValue{
		"someKey": &types.AttributeValueMemberS{
			Value: "someValue",
		},
		"anotherKey": &types.AttributeValueMemberS{
			Value: "anotherValue",
		},
	}
	getItemOutput := &dynamodb.GetItemOutput{
		Item: item,
	}
	suite.dynamodbClient.EXPECT().GetItem(&context, &getItemInput).Return(getItemOutput, nil)

	response, appErr := suite.baseRepository.FindByComplexPrimaryKey(&context, primaryKey, false, false)

	suite.NoError(appErr)
	suite.Equal(response, item)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestFindByComplexPrimaryKey_ShouldReturnInternalServerErrorWhenGetItemFailed() {
	context := common_models.NewLambdaContext()
	primaryKey := common_models.DynamodbComplexPrimaryKey{
		PartitionKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "somePartitionKey",
			Value:   "somePartitionValue",
		},
		SortKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "someSortKey",
			Value:   "someSortValue",
		},
	}
	getItemInput := dynamodb.GetItemInput{
		TableName:      aws.String("someTable"),
		ConsistentRead: aws.Bool(false),
		Key: map[string]types.AttributeValue{
			"somePartitionKey": &types.AttributeValueMemberS{
				Value: "somePartitionValue",
			},
			"someSortKey": &types.AttributeValueMemberS{
				Value: "someSortValue",
			},
		},
	}
	getItemOutput := &dynamodb.GetItemOutput{}
	cause := errors.New("someErr")
	expectedAppErr := common_errors.NewInternalServerError("error while reading from database")

	suite.dynamodbClient.EXPECT().GetItem(&context, &getItemInput).Return(getItemOutput, cause)

	_, appErr := suite.baseRepository.FindByComplexPrimaryKey(&context, primaryKey, false, false)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestFindByComplexPrimaryKey_ShouldSucceedWhenTransaction() {
	context := common_models.NewLambdaContext()
	transactGetItemsInput := dynamodb.TransactGetItemsInput{
		TransactItems: []types.TransactGetItem{},
	}
	context.Set(common_constants.ReadTransaction, transactGetItemsInput)
	expectedContext := common_models.NewLambdaContext()
	expectedGetItemsInput := dynamodb.TransactGetItemsInput{
		TransactItems: []types.TransactGetItem{
			{
				Get: &types.Get{
					TableName: aws.String("someTable"),
					Key: map[string]types.AttributeValue{
						"somePartitionKey": &types.AttributeValueMemberS{
							Value: "somePartitionValue",
						},
						"someSortKey": &types.AttributeValueMemberS{
							Value: "someSortValue",
						},
					},
				},
			},
		},
	}
	expectedContext.Set(common_constants.ReadTransaction, expectedGetItemsInput)
	primaryKey := common_models.DynamodbComplexPrimaryKey{
		PartitionKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "somePartitionKey",
			Value:   "somePartitionValue",
		},
		SortKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "someSortKey",
			Value:   "someSortValue",
		},
	}

	_, appErr := suite.baseRepository.FindByComplexPrimaryKey(&context, primaryKey, false, true)
	actualGetItemsInput, exists := context.Get(common_constants.ReadTransaction)

	suite.NoError(appErr)
	suite.True(exists)
	suite.Equal(expectedGetItemsInput, actualGetItemsInput)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestSaveIfNotPresentWithSimplePrimaryKey_ShouldSucceedWhenNoTransaction() {
	context := common_models.NewLambdaContext()
	primaryKey := common_models.DynamodbSimplePrimaryKey{
		KeyName: "pk",
		Value:   "someKey",
	}
	item := DummyItem{Key1: "foo", Key2: "bar"}
	putItemInput := dynamodb.PutItemInput{
		TableName:           aws.String("someTable"),
		ConditionExpression: aws.String("attribute_not_exists (#0)"),
		ExpressionAttributeNames: map[string]string{
			"#0": "pk",
		},
		ExpressionAttributeValues: nil,
		Item: map[string]types.AttributeValue{
			"key1": &types.AttributeValueMemberS{
				Value: "foo",
			},
			"key2": &types.AttributeValueMemberS{
				Value: "bar",
			},
			"pk": &types.AttributeValueMemberS{
				Value: "someKey",
			},
		},
	}
	putItemOutput := &dynamodb.PutItemOutput{}

	suite.dynamodbClient.EXPECT().PutItem(&context, &putItemInput).Return(putItemOutput, nil)

	appErr := suite.baseRepository.SaveIfNotPresentWithSimplePrimaryKey(&context, primaryKey, item, false)

	suite.NoError(appErr)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestSaveIfNotPresentWithSimplePrimaryKey_ShouldReturnInternalServerErrorWhenPutItemFailed() {
	context := common_models.NewLambdaContext()
	primaryKey := common_models.DynamodbSimplePrimaryKey{
		KeyName: "pk",
		Value:   "someKey",
	}
	item := DummyItem{Key1: "foo", Key2: "bar"}
	putItemInput := dynamodb.PutItemInput{
		TableName:           aws.String("someTable"),
		ConditionExpression: aws.String("attribute_not_exists (#0)"),
		ExpressionAttributeNames: map[string]string{
			"#0": "pk",
		},
		ExpressionAttributeValues: nil,
		Item: map[string]types.AttributeValue{
			"key1": &types.AttributeValueMemberS{
				Value: "foo",
			},
			"key2": &types.AttributeValueMemberS{
				Value: "bar",
			},
			"pk": &types.AttributeValueMemberS{
				Value: "someKey",
			},
		},
	}
	putItemOutput := &dynamodb.PutItemOutput{}
	cause := errors.New("someErr")
	expectedAppErr := common_errors.NewInternalServerError("error while writing into database")

	suite.dynamodbClient.EXPECT().PutItem(&context, &putItemInput).Return(putItemOutput, cause)

	appErr := suite.baseRepository.SaveIfNotPresentWithSimplePrimaryKey(&context, primaryKey, item, false)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestSaveIfNotPresentWithSimplePrimaryKey_ShouldSucceedWhenTransaction() {
	context := common_models.NewLambdaContext()
	transactGetItemsInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{},
	}
	context.Set(common_constants.WriteTransaction, transactGetItemsInput)
	expectedContext := common_models.NewLambdaContext()
	expectedWriteItemsInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName:           aws.String("someTable"),
					ConditionExpression: aws.String("attribute_not_exists (#0)"),
					ExpressionAttributeNames: map[string]string{
						"#0": "pk",
					},
					ExpressionAttributeValues: nil,
					Item: map[string]types.AttributeValue{
						"pk": &types.AttributeValueMemberS{
							Value: "someKey",
						},
						"key1": &types.AttributeValueMemberS{
							Value: "foo",
						},
						"key2": &types.AttributeValueMemberS{
							Value: "bar",
						},
					},
				},
			},
		},
	}
	expectedContext.Set(common_constants.WriteTransaction, expectedWriteItemsInput)
	primaryKey := common_models.DynamodbSimplePrimaryKey{
		KeyName: "pk",
		Value:   "someKey",
	}
	item := DummyItem{Key1: "foo", Key2: "bar"}

	appErr := suite.baseRepository.SaveIfNotPresentWithSimplePrimaryKey(&context, primaryKey, item, true)

	actualWriteItemInput, exists := context.Get(common_constants.WriteTransaction)

	suite.NoError(appErr)
	suite.True(exists)
	suite.Equal(expectedWriteItemsInput, actualWriteItemInput)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestSaveIfNotPresentWithComplexPrimaryKey_ShouldSucceedWhenNoTransaction() {
	context := common_models.NewLambdaContext()
	primaryKey := common_models.DynamodbComplexPrimaryKey{
		PartitionKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "somePartitionKey",
			Value:   "somePartitionValue",
		},
		SortKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "someSortKey",
			Value:   "someSortValue",
		},
	}
	item := DummyItem{Key1: "foo", Key2: "bar"}
	putItemInput := dynamodb.PutItemInput{
		TableName:           aws.String("someTable"),
		ConditionExpression: aws.String("(attribute_not_exists (#0)) AND (attribute_not_exists (#1))"),
		ExpressionAttributeNames: map[string]string{
			"#0": "somePartitionKey",
			"#1": "someSortKey",
		},
		ExpressionAttributeValues: nil,
		Item: map[string]types.AttributeValue{
			"key1": &types.AttributeValueMemberS{
				Value: "foo",
			},
			"key2": &types.AttributeValueMemberS{
				Value: "bar",
			},
			"somePartitionKey": &types.AttributeValueMemberS{
				Value: "somePartitionValue",
			},
			"someSortKey": &types.AttributeValueMemberS{
				Value: "someSortValue",
			},
		},
	}
	putItemOutput := &dynamodb.PutItemOutput{}

	suite.dynamodbClient.EXPECT().PutItem(&context, &putItemInput).Return(putItemOutput, nil)

	appErr := suite.baseRepository.SaveIfNotPresentWithComplexPrimaryKey(&context, primaryKey, item, false)

	suite.NoError(appErr)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestSaveIfNotPresentWithComplexPrimaryKey_ShouldReturnInternalServerErrorWhenPutItemFailed() {
	context := common_models.NewLambdaContext()
	primaryKey := common_models.DynamodbComplexPrimaryKey{
		PartitionKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "somePartitionKey",
			Value:   "somePartitionValue",
		},
		SortKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "someSortKey",
			Value:   "someSortValue",
		},
	}
	item := DummyItem{Key1: "foo", Key2: "bar"}
	putItemInput := dynamodb.PutItemInput{
		TableName:           aws.String("someTable"),
		ConditionExpression: aws.String("(attribute_not_exists (#0)) AND (attribute_not_exists (#1))"),
		ExpressionAttributeNames: map[string]string{
			"#0": "somePartitionKey",
			"#1": "someSortKey",
		},
		ExpressionAttributeValues: nil,
		Item: map[string]types.AttributeValue{
			"key1": &types.AttributeValueMemberS{
				Value: "foo",
			},
			"key2": &types.AttributeValueMemberS{
				Value: "bar",
			},
			"somePartitionKey": &types.AttributeValueMemberS{
				Value: "somePartitionValue",
			},
			"someSortKey": &types.AttributeValueMemberS{
				Value: "someSortValue",
			},
		},
	}
	putItemOutput := &dynamodb.PutItemOutput{}
	cause := errors.New("someErr")
	expectedAppErr := common_errors.NewInternalServerError("error while writing into database")

	suite.dynamodbClient.EXPECT().PutItem(&context, &putItemInput).Return(putItemOutput, cause)

	appErr := suite.baseRepository.SaveIfNotPresentWithComplexPrimaryKey(&context, primaryKey, item, false)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestSaveIfNotPresentWithComplexPrimaryKey_ShouldSucceedWhenTransaction() {
	context := common_models.NewLambdaContext()
	transactGetItemsInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{},
	}
	context.Set(common_constants.WriteTransaction, transactGetItemsInput)
	expectedContext := common_models.NewLambdaContext()
	expectedWriteItemsInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName:           aws.String("someTable"),
					ConditionExpression: aws.String("(attribute_not_exists (#0)) AND (attribute_not_exists (#1))"),
					ExpressionAttributeNames: map[string]string{
						"#0": "somePartitionKey",
						"#1": "someSortKey",
					},
					ExpressionAttributeValues: nil,
					Item: map[string]types.AttributeValue{
						"somePartitionKey": &types.AttributeValueMemberS{
							Value: "somePartitionValue",
						},
						"someSortKey": &types.AttributeValueMemberS{
							Value: "someSortValue",
						},
						"key1": &types.AttributeValueMemberS{
							Value: "foo",
						},
						"key2": &types.AttributeValueMemberS{
							Value: "bar",
						},
					},
				},
			},
		},
	}
	expectedContext.Set(common_constants.WriteTransaction, expectedWriteItemsInput)
	primaryKey := common_models.DynamodbComplexPrimaryKey{
		PartitionKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "somePartitionKey",
			Value:   "somePartitionValue",
		},
		SortKey: common_models.DynamodbSimplePrimaryKey{
			KeyName: "someSortKey",
			Value:   "someSortValue",
		},
	}
	item := DummyItem{Key1: "foo", Key2: "bar"}

	appErr := suite.baseRepository.SaveIfNotPresentWithComplexPrimaryKey(&context, primaryKey, item, true)

	actualWriteItemInput, exists := context.Get(common_constants.WriteTransaction)

	suite.NoError(appErr)
	suite.True(exists)
	suite.Equal(expectedWriteItemsInput, actualWriteItemInput)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestSave_ShouldSucceedWhenNoTransaction() {
	context := common_models.NewLambdaContext()
	item := DummyItem{Key1: "foo", Key2: "bar"}
	putItemInput := dynamodb.PutItemInput{
		TableName:                 aws.String("someTable"),
		ExpressionAttributeValues: nil,
		Item: map[string]types.AttributeValue{
			"key1": &types.AttributeValueMemberS{
				Value: "foo",
			},
			"key2": &types.AttributeValueMemberS{
				Value: "bar",
			},
		},
	}
	putItemOutput := &dynamodb.PutItemOutput{}

	suite.dynamodbClient.EXPECT().PutItem(&context, &putItemInput).Return(putItemOutput, nil)

	appErr := suite.baseRepository.Save(&context, item, false)

	suite.NoError(appErr)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestSave_ShouldReturnInternalServerErrorWhenPutItemFailed() {
	context := common_models.NewLambdaContext()
	item := DummyItem{Key1: "foo", Key2: "bar"}
	putItemInput := dynamodb.PutItemInput{
		TableName: aws.String("someTable"),
		Item: map[string]types.AttributeValue{
			"key1": &types.AttributeValueMemberS{
				Value: "foo",
			},
			"key2": &types.AttributeValueMemberS{
				Value: "bar",
			},
		},
	}
	putItemOutput := &dynamodb.PutItemOutput{}
	cause := errors.New("someErr")
	expectedAppErr := common_errors.NewInternalServerError("error while writing into database")

	suite.dynamodbClient.EXPECT().PutItem(&context, &putItemInput).Return(putItemOutput, cause)

	appErr := suite.baseRepository.Save(&context, item, false)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbBaseRepositoryTestSuite) TestSave_ShouldSucceedWhenTransaction() {
	context := common_models.NewLambdaContext()
	transactGetItemsInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{},
	}
	context.Set(common_constants.WriteTransaction, transactGetItemsInput)
	expectedContext := common_models.NewLambdaContext()
	expectedWriteItemsInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName:                 aws.String("someTable"),
					ExpressionAttributeValues: nil,
					Item: map[string]types.AttributeValue{
						"key1": &types.AttributeValueMemberS{
							Value: "foo",
						},
						"key2": &types.AttributeValueMemberS{
							Value: "bar",
						},
					},
				},
			},
		},
	}
	expectedContext.Set(common_constants.WriteTransaction, expectedWriteItemsInput)
	item := DummyItem{Key1: "foo", Key2: "bar"}

	appErr := suite.baseRepository.Save(&context, item, true)

	actualWriteItemInput, exists := context.Get(common_constants.WriteTransaction)

	suite.NoError(appErr)
	suite.True(exists)
	suite.Equal(expectedWriteItemsInput, actualWriteItemInput)
}
