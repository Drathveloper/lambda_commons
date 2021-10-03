package common_repositories_test

import (
	"errors"
	"github.com/Drathveloper/lambda_commons/v2/common_constants"
	"github.com/Drathveloper/lambda_commons/v2/common_errors"
	"github.com/Drathveloper/lambda_commons/v2/common_models"
	"github.com/Drathveloper/lambda_commons/v2/common_repositories"
	"github.com/Drathveloper/lambda_commons/v2/mocks"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DynamodbTransactionManagerTestSuite struct {
	suite.Suite
	dynamodbClient     *mocks.MockDynamodbClientAPI
	transactionManager common_repositories.DynamodbTransactionManager
}

func TestDynamodbTransactionManagerTestSuite(t *testing.T) {
	suite.Run(t, new(DynamodbTransactionManagerTestSuite))
}

func (suite *DynamodbTransactionManagerTestSuite) SetupTest() {
	controller := gomock.NewController(suite.T())
	suite.dynamodbClient = mocks.NewMockDynamodbClientAPI(controller)
	suite.transactionManager = common_repositories.NewDynamodbTransactionManager(suite.dynamodbClient)
}

func (suite *DynamodbTransactionManagerTestSuite) TestStartReadTransaction_ShouldSucceed() {
	context := common_models.NewLambdaContext()
	expectedContext := common_models.NewLambdaContext()
	transactGetItems := make([]types.TransactGetItem, 0)
	transactionInput := dynamodb.TransactGetItemsInput{
		TransactItems: transactGetItems,
	}
	expectedContext.Set(common_constants.ReadTransaction, transactionInput)

	appErr := suite.transactionManager.StartReadTransaction(&context)

	suite.NoError(appErr)
	suite.Equal(expectedContext, context)
}

func (suite *DynamodbTransactionManagerTestSuite) TestStartReadTransaction_ShouldReturnErrorWhenTransactionAlreadyStarted() {
	context := common_models.NewLambdaContext()
	transactGetItems := make([]types.TransactGetItem, 0)
	transactionInput := dynamodb.TransactGetItemsInput{
		TransactItems: transactGetItems,
	}
	context.Set(common_constants.ReadTransaction, transactionInput)
	expectedAppErr := common_errors.NewInternalServerError("there is already a read transaction in progress in this scope")

	appErr := suite.transactionManager.StartReadTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbTransactionManagerTestSuite) TestExecuteReadTransaction_ShouldSucceed() {
	context := common_models.NewLambdaContext()
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
	context.Set(common_constants.ReadTransaction, transactionInput)
	expectedItems := map[string]types.AttributeValue{
		"someTable#someKey": &types.AttributeValueMemberS{
			Value: "someValue",
		},
	}

	suite.dynamodbClient.EXPECT().TransactGetItems(&context, &transactionInput).Return(&transactionOutput, nil)

	response, appErr := suite.transactionManager.ExecuteReadTransaction(&context)

	suite.NoError(appErr)
	suite.Equal(expectedItems, response)
}

func (suite *DynamodbTransactionManagerTestSuite) TestExecuteReadTransaction_ShouldReturnInternalServerErrorWhenTransactionFailed() {
	context := common_models.NewLambdaContext()
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
	context.Set(common_constants.ReadTransaction, transactionInput)
	cause := errors.New("someErr")
	expectedAppErr := common_errors.NewInternalServerError("generic error performing transaction")
	suite.dynamodbClient.EXPECT().TransactGetItems(&context, &transactionInput).Return(&transactionOutput, cause)

	_, appErr := suite.transactionManager.ExecuteReadTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbTransactionManagerTestSuite) TestExecuteReadTransaction_ShouldReturnInternalServerErrorWhenNoTransactionStarted() {
	context := common_models.NewLambdaContext()

	expectedAppErr := common_errors.NewInternalServerError("there is no read transaction in progress")

	_, appErr := suite.transactionManager.ExecuteReadTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbTransactionManagerTestSuite) TestStartWriteTransaction_ShouldSucceed() {
	context := common_models.NewLambdaContext()
	expectedContext := common_models.NewLambdaContext()
	transactWriteItems := make([]types.TransactWriteItem, 0)
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: transactWriteItems,
	}
	expectedContext.Set(common_constants.WriteTransaction, transactionInput)

	appErr := suite.transactionManager.StartWriteTransaction(&context)

	suite.NoError(appErr)
	suite.Equal(expectedContext, context)
}

func (suite *DynamodbTransactionManagerTestSuite) TestStartWriteTransaction_ShouldReturnErrorWhenTransactionAlreadyStarted() {
	context := common_models.NewLambdaContext()
	transactWriteItems := make([]types.TransactWriteItem, 0)
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: transactWriteItems,
	}
	context.Set(common_constants.WriteTransaction, transactionInput)
	expectedAppErr := common_errors.NewInternalServerError("there is already a write transaction in progress in this scope")

	appErr := suite.transactionManager.StartWriteTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbTransactionManagerTestSuite) TestExecuteWriteTransaction_ShouldSucceed() {
	context := common_models.NewLambdaContext()
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
	context.Set(common_constants.WriteTransaction, transactionInput)

	suite.dynamodbClient.EXPECT().TransactWriteItems(&context, &transactionInput).Return(&transactionOutput, nil)

	appErr := suite.transactionManager.ExecuteWriteTransaction(&context)

	suite.NoError(appErr)
}

func (suite *DynamodbTransactionManagerTestSuite) TestExecuteWriteTransaction_ShouldReturnInternalServerErrorWhenTransactionFailed() {
	context := common_models.NewLambdaContext()
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
	context.Set(common_constants.WriteTransaction, transactionInput)
	cause := errors.New("someErr")
	expectedAppErr := common_errors.NewInternalServerError("generic error performing transaction")
	suite.dynamodbClient.EXPECT().TransactWriteItems(&context, &transactionInput).Return(&transactionOutput, cause)

	appErr := suite.transactionManager.ExecuteWriteTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}

func (suite *DynamodbTransactionManagerTestSuite) TestExecuteWriteTransaction_ShouldReturnInternalServerErrorWhenNoTransactionStarted() {
	context := common_models.NewLambdaContext()

	expectedAppErr := common_errors.NewInternalServerError("there is no write transaction in progress")

	appErr := suite.transactionManager.ExecuteWriteTransaction(&context)

	suite.Equal(expectedAppErr, appErr)
}
