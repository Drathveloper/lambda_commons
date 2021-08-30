package repositories

import (
	"errors"
	"fmt"
	"github.com/Drathveloper/lambda_commons/constants"
	"github.com/Drathveloper/lambda_commons/custom_errors"
	"github.com/Drathveloper/lambda_commons/helpers"
	"github.com/Drathveloper/lambda_commons/models"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamodbTransactionManager interface {
	StartReadTransaction(ctx *models.LambdaContext) custom_errors.GenericApplicationError
	ExecuteReadTransaction(ctx *models.LambdaContext) (map[string]types.AttributeValue, custom_errors.GenericApplicationError)
	StartWriteTransaction(ctx *models.LambdaContext) custom_errors.GenericApplicationError
	ExecuteWriteTransaction(ctx *models.LambdaContext) custom_errors.GenericApplicationError
}

type dynamodbTransactionalRepository struct {
	client models.DynamodbClientAPI
}

func NewDynamodbTransactionManager(client models.DynamodbClientAPI) DynamodbTransactionManager {
	return &dynamodbTransactionalRepository{
		client: client,
	}
}

func (repository *dynamodbTransactionalRepository) StartReadTransaction(ctx *models.LambdaContext) custom_errors.GenericApplicationError {
	if ctx.Exists(constants.ReadTransaction) {
		return custom_errors.NewInternalServerError("there is already a read transaction in progress in this scope")
	}
	transactGetItems := make([]types.TransactGetItem, 0)
	transactionInput := dynamodb.TransactGetItemsInput{
		TransactItems: transactGetItems,
	}
	ctx.Set(constants.ReadTransaction, transactionInput)
	return nil
}

func (repository *dynamodbTransactionalRepository) ExecuteReadTransaction(ctx *models.LambdaContext) (map[string]types.AttributeValue, custom_errors.GenericApplicationError) {
	input, exists := ctx.Get(constants.ReadTransaction)
	defer ctx.Set(constants.ReadTransaction, nil)
	if !exists {
		return nil, custom_errors.NewInternalServerError("there is no read transaction in progress")
	}
	transactionInput := input.(dynamodb.TransactGetItemsInput)
	transactionOutput, err := repository.client.TransactGetItems(ctx, &transactionInput)
	if err != nil {
		return nil, repository.handleTransactionError(err)
	}
	tableNames := make([]string, 0)
	for _, request := range transactionInput.TransactItems {
		tableNames = append(tableNames, *request.Get.TableName)
	}
	return helpers.MergeResponsesIntoAttributeValueMap(tableNames, transactionOutput.Responses), nil
}

func (repository *dynamodbTransactionalRepository) StartWriteTransaction(ctx *models.LambdaContext) custom_errors.GenericApplicationError {
	if ctx.Exists(constants.WriteTransaction) {
		return custom_errors.NewInternalServerError("there is already a write transaction in progress in this scope")
	}
	transactWriteItems := make([]types.TransactWriteItem, 0)
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: transactWriteItems,
	}
	ctx.Set(constants.WriteTransaction, transactionInput)
	return nil
}

func (repository *dynamodbTransactionalRepository) ExecuteWriteTransaction(ctx *models.LambdaContext) custom_errors.GenericApplicationError {
	input, exists := ctx.Get(constants.WriteTransaction)
	defer ctx.Set(constants.WriteTransaction, nil)
	if !exists {
		return custom_errors.NewInternalServerError("there is no write transaction in progress")
	}
	transactionInput := input.(dynamodb.TransactWriteItemsInput)
	_, err := repository.client.TransactWriteItems(ctx, &transactionInput)
	if err != nil {
		return repository.handleTransactionError(err)
	}
	return nil
}

func (repository *dynamodbTransactionalRepository) handleTransactionError(err error) custom_errors.GenericApplicationError {
	var dynamodbErr *types.TransactionCanceledException
	if errors.As(err, &dynamodbErr) {
		for _, reason := range dynamodbErr.CancellationReasons {
			if constants.ConditionalCheckFailed == *reason.Code {
				return custom_errors.NewInternalServerError(fmt.Sprintf("conditional check failed: %s", *reason.Message))
			}
		}
	}
	return custom_errors.NewInternalServerError("generic error performing transaction")
}
