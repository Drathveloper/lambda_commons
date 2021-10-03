package common_repositories

import (
	"errors"
	"fmt"
	"github.com/Drathveloper/lambda_commons/common_constants"
	"github.com/Drathveloper/lambda_commons/common_errors"
	"github.com/Drathveloper/lambda_commons/common_helpers"
	"github.com/Drathveloper/lambda_commons/common_models"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamodbTransactionManager interface {
	StartReadTransaction(ctx *common_models.LambdaContext) common_errors.GenericApplicationError
	ExecuteReadTransaction(ctx *common_models.LambdaContext) (map[string]types.AttributeValue, common_errors.GenericApplicationError)
	StartWriteTransaction(ctx *common_models.LambdaContext) common_errors.GenericApplicationError
	ExecuteWriteTransaction(ctx *common_models.LambdaContext) common_errors.GenericApplicationError
}

type dynamodbTransactionalRepository struct {
	client common_models.DynamodbClientAPI
}

func NewDynamodbTransactionManager(client common_models.DynamodbClientAPI) DynamodbTransactionManager {
	return &dynamodbTransactionalRepository{
		client: client,
	}
}

func (repository *dynamodbTransactionalRepository) StartReadTransaction(ctx *common_models.LambdaContext) common_errors.GenericApplicationError {
	if ctx.Exists(common_constants.ReadTransaction) {
		return common_errors.NewInternalServerError("there is already a read transaction in progress in this scope")
	}
	transactGetItems := make([]types.TransactGetItem, 0)
	transactionInput := dynamodb.TransactGetItemsInput{
		TransactItems: transactGetItems,
	}
	ctx.Set(common_constants.ReadTransaction, transactionInput)
	return nil
}

func (repository *dynamodbTransactionalRepository) ExecuteReadTransaction(ctx *common_models.LambdaContext) (map[string]types.AttributeValue, common_errors.GenericApplicationError) {
	input, exists := ctx.Get(common_constants.ReadTransaction)
	defer ctx.Set(common_constants.ReadTransaction, nil)
	if !exists {
		return nil, common_errors.NewInternalServerError("there is no read transaction in progress")
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
	return common_helpers.MergeDynamoDBResponsesIntoAttributeValueMap(tableNames, transactionOutput.Responses)
}

func (repository *dynamodbTransactionalRepository) StartWriteTransaction(ctx *common_models.LambdaContext) common_errors.GenericApplicationError {
	if ctx.Exists(common_constants.WriteTransaction) {
		return common_errors.NewInternalServerError("there is already a write transaction in progress in this scope")
	}
	transactWriteItems := make([]types.TransactWriteItem, 0)
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: transactWriteItems,
	}
	ctx.Set(common_constants.WriteTransaction, transactionInput)
	return nil
}

func (repository *dynamodbTransactionalRepository) ExecuteWriteTransaction(ctx *common_models.LambdaContext) common_errors.GenericApplicationError {
	input, exists := ctx.Get(common_constants.WriteTransaction)
	defer ctx.Set(common_constants.WriteTransaction, nil)
	if !exists {
		return common_errors.NewInternalServerError("there is no write transaction in progress")
	}
	transactionInput := input.(dynamodb.TransactWriteItemsInput)
	_, err := repository.client.TransactWriteItems(ctx, &transactionInput)
	if err != nil {
		return repository.handleTransactionError(err)
	}
	return nil
}

func (repository *dynamodbTransactionalRepository) handleTransactionError(err error) common_errors.GenericApplicationError {
	var dynamodbErr *types.TransactionCanceledException
	if errors.As(err, &dynamodbErr) {
		for _, reason := range dynamodbErr.CancellationReasons {
			if common_constants.ConditionalCheckFailed == *reason.Code {
				return common_errors.NewForbiddenError(fmt.Sprintf("conditional check failed: %s", *reason.Message))
			}
		}
	}
	return common_errors.NewInternalServerError("generic error performing transaction")
}
