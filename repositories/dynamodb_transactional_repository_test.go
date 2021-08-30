package repositories_test

import (
	"errors"
	"github.com/Drathveloper/lambda_commons/constants"
	"github.com/Drathveloper/lambda_commons/custom_errors"
	"github.com/Drathveloper/lambda_commons/mocks"
	"github.com/Drathveloper/lambda_commons/models"
	"github.com/Drathveloper/lambda_commons/repositories"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DynamodbTransactionalRepositoryTestSuite struct {
	suite.Suite
	dynamodbClient          *mocks.MockDynamodbClientAPI
	transactionalRepository repositories.DynamodbTransactionalRepository
}

func TestDynamodbTransactionalRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DynamodbTransactionalRepositoryTestSuite))
}

func (suite *DynamodbTransactionalRepositoryTestSuite) SetupTest() {
	controller := gomock.NewController(suite.T())
	suite.dynamodbClient = mocks.NewMockDynamodbClientAPI(controller)
	suite.transactionalRepository = repositories.NewDynamodbTransactionalRepository(suite.dynamodbClient)
}

func (suite *DynamodbTransactionalRepositoryTestSuite) TestStartReadTransaction_ShouldSucceed() {
	context := models.NewLambdaContext()
	expectedContext := models.NewLambdaContext()
	transactGetItems := make([]types.TransactGetItem, 0)
	transactionInput := dynamodb.TransactGetItemsInput{
		TransactItems: transactGetItems,
	}
	expectedContext.Set(constants.ReadTransaction, transactionInput)

	appErr := suite.transactionalRepository.StartReadTransaction(&context)

	suite.NoError(appErr)
	suite.Equal(expectedContext, context)
}

func (suite *DynamodbTransactionalRepositoryTestSuite) TestStartReadTransaction_ShouldReturnErrorWhenTransactionAlreadyStarted() {
	context := models.NewLambdaContext()
	transactGetItems := make([]types.TransactGetItem, 0)
	transactionInput := dynamodb.TransactGetItemsInput{
		TransactItems: transactGetItems,
	}
	context.Set(constants.ReadTransaction, transactionInput)
	expectedAppErr := custom_errors.NewInternalServerError("there is already a read transaction in progress in this scope")

	appErr := suite.transactionalRepository.StartReadTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbTransactionalRepositoryTestSuite) TestExecuteReadTransaction_ShouldSucceed() {
	context := models.NewLambdaContext()
	transactionInput := dynamodb.TransactGetItemsInput{
		TransactItems: []types.TransactGetItem{
			{
				Get: &types.Get{
					TableName: aws.String("someTable"),
				},
			},
		},
	}
	item := map[string]types.AttributeValue{
		"someKey": &types.AttributeValueMemberS{
			Value: "someValue",
		},
	}
	transactionOutput := dynamodb.TransactGetItemsOutput{
		Responses: []types.ItemResponse{
			{
				Item: item,
			},
		},
	}
	context.Set(constants.ReadTransaction, transactionInput)
	expectedItems := map[string]types.AttributeValue{
		"someTable#someKey": &types.AttributeValueMemberS{
			Value: "someValue",
		},
	}

	suite.dynamodbClient.EXPECT().TransactGetItems(&context, &transactionInput).Return(&transactionOutput, nil)

	response, appErr := suite.transactionalRepository.ExecuteReadTransaction(&context)

	suite.NoError(appErr)
	suite.Equal(expectedItems, response)
}

func (suite *DynamodbTransactionalRepositoryTestSuite) TestExecuteReadTransaction_ShouldReturnInternalServerErrorWhenTransactionFailed() {
	context := models.NewLambdaContext()
	transactionInput := dynamodb.TransactGetItemsInput{
		TransactItems: []types.TransactGetItem{
			{
				Get: &types.Get{
					TableName: aws.String("someTable"),
				},
			},
		},
	}
	transactionOutput := dynamodb.TransactGetItemsOutput{}
	context.Set(constants.ReadTransaction, transactionInput)
	cause := errors.New("someErr")
	expectedAppErr := custom_errors.NewInternalServerError("error performing read transaction")
	suite.dynamodbClient.EXPECT().TransactGetItems(&context, &transactionInput).Return(&transactionOutput, cause)

	_, appErr := suite.transactionalRepository.ExecuteReadTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbTransactionalRepositoryTestSuite) TestExecuteReadTransaction_ShouldReturnInternalServerErrorWhenNoTransactionStarted() {
	context := models.NewLambdaContext()

	expectedAppErr := custom_errors.NewInternalServerError("there is no read transaction in progress")

	_, appErr := suite.transactionalRepository.ExecuteReadTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbTransactionalRepositoryTestSuite) TestStartWriteTransaction_ShouldSucceed() {
	context := models.NewLambdaContext()
	expectedContext := models.NewLambdaContext()
	transactWriteItems := make([]types.TransactWriteItem, 0)
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: transactWriteItems,
	}
	expectedContext.Set(constants.WriteTransaction, transactionInput)

	appErr := suite.transactionalRepository.StartWriteTransaction(&context)

	suite.NoError(appErr)
	suite.Equal(expectedContext, context)
}

func (suite *DynamodbTransactionalRepositoryTestSuite) TestStartWriteTransaction_ShouldReturnErrorWhenTransactionAlreadyStarted() {
	context := models.NewLambdaContext()
	transactWriteItems := make([]types.TransactWriteItem, 0)
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: transactWriteItems,
	}
	context.Set(constants.WriteTransaction, transactionInput)
	expectedAppErr := custom_errors.NewInternalServerError("there is already a write transaction in progress in this scope")

	appErr := suite.transactionalRepository.StartWriteTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbTransactionalRepositoryTestSuite) TestExecuteWriteTransaction_ShouldSucceed() {
	context := models.NewLambdaContext()
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: aws.String("someTable"),
				},
			},
		},
	}
	transactionOutput := dynamodb.TransactWriteItemsOutput{}
	context.Set(constants.WriteTransaction, transactionInput)

	suite.dynamodbClient.EXPECT().TransactWriteItems(&context, &transactionInput).Return(&transactionOutput, nil)

	appErr := suite.transactionalRepository.ExecuteWriteTransaction(&context)

	suite.NoError(appErr)
}

func (suite *DynamodbTransactionalRepositoryTestSuite) TestExecuteWriteTransaction_ShouldReturnInternalServerErrorWhenTransactionFailed() {
	context := models.NewLambdaContext()
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: aws.String("someTable"),
				},
			},
		},
	}
	transactionOutput := dynamodb.TransactWriteItemsOutput{}
	context.Set(constants.WriteTransaction, transactionInput)
	cause := errors.New("someErr")
	expectedAppErr := custom_errors.NewInternalServerError("error performing write transaction")
	suite.dynamodbClient.EXPECT().TransactWriteItems(&context, &transactionInput).Return(&transactionOutput, cause)

	appErr := suite.transactionalRepository.ExecuteWriteTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbTransactionalRepositoryTestSuite) TestExecuteWriteTransaction_ShouldReturnInternalServerErrorWhenNoTransactionStarted() {
	context := models.NewLambdaContext()

	expectedAppErr := custom_errors.NewInternalServerError("there is no write transaction in progress")

	appErr := suite.transactionalRepository.ExecuteWriteTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}
