package repositories

import (
	"errors"
	"github.com/Drathveloper/lambda_commons/constants"
	"github.com/Drathveloper/lambda_commons/custom_errors"
	"github.com/Drathveloper/lambda_commons/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamodbBaseRepository interface {
	FindBySimplePrimaryKey(ctx *models.LambdaContext, primaryKey models.DynamodbSimplePrimaryKey, isConsistentRead, transactional bool) (map[string]types.AttributeValue, custom_errors.GenericApplicationError)
	FindByComplexPrimaryKey(ctx *models.LambdaContext, primaryKey models.DynamodbComplexPrimaryKey, isConsistentRead, transactional bool) (map[string]types.AttributeValue, custom_errors.GenericApplicationError)
	SaveIfNotPresentWithSimplePrimaryKey(ctx *models.LambdaContext, primaryKey models.DynamodbSimplePrimaryKey, item interface{}, transactional bool) custom_errors.GenericApplicationError
	SaveIfNotPresentWithComplexPrimaryKey(ctx *models.LambdaContext, primaryKey models.DynamodbComplexPrimaryKey, item interface{}, transactional bool) custom_errors.GenericApplicationError
	Save(ctx *models.LambdaContext, item interface{}, transactional bool) custom_errors.GenericApplicationError
}

type dynamodbBaseRepository struct {
	tableName string
	client    models.DynamodbClientAPI
}

func NewDynamodbBaseRepository(client models.DynamodbClientAPI, tableName string) DynamodbBaseRepository {
	return &dynamodbBaseRepository{
		tableName: tableName,
		client:    client,
	}
}

func (repository *dynamodbBaseRepository) FindBySimplePrimaryKey(ctx *models.LambdaContext, primaryKey models.DynamodbSimplePrimaryKey, isConsistentRead, transactional bool) (map[string]types.AttributeValue, custom_errors.GenericApplicationError) {
	value, err := attributevalue.Marshal(primaryKey.Value)
	if err != nil {
		return nil, custom_errors.NewInternalServerError("error while marshaling database primary key")
	}
	keyValues := map[string]types.AttributeValue{
		primaryKey.KeyName: value,
	}
	return repository.findByPrimaryKey(ctx, keyValues, isConsistentRead, transactional)
}

func (repository *dynamodbBaseRepository) FindByComplexPrimaryKey(ctx *models.LambdaContext, primaryKey models.DynamodbComplexPrimaryKey, isConsistentRead, transactional bool) (map[string]types.AttributeValue, custom_errors.GenericApplicationError) {
	partitionKeyValue, err := attributevalue.Marshal(primaryKey.PartitionKey.Value)
	if err != nil {
		return nil, custom_errors.NewInternalServerError("error while marshaling database partition key")
	}
	sortKeyValue, err := attributevalue.Marshal(primaryKey.SortKey.Value)
	if err != nil {
		return nil, custom_errors.NewInternalServerError("error while marshaling database sort key")
	}
	keyValues := map[string]types.AttributeValue{
		primaryKey.PartitionKey.KeyName: partitionKeyValue,
		primaryKey.SortKey.KeyName:      sortKeyValue,
	}
	return repository.findByPrimaryKey(ctx, keyValues, isConsistentRead, transactional)
}

func (repository *dynamodbBaseRepository) SaveIfNotPresentWithSimplePrimaryKey(ctx *models.LambdaContext, primaryKey models.DynamodbSimplePrimaryKey, item interface{}, transactional bool) custom_errors.GenericApplicationError {
	primaryKeyValue, err := attributevalue.Marshal(primaryKey.Value)
	if err != nil {
		return custom_errors.NewInternalServerError("error while marshaling database partition key")
	}
	expressionBuilder := expression.NewBuilder()
	condition := expression.AttributeNotExists(expression.Name(primaryKey.KeyName))
	builtExpression, err := expressionBuilder.WithCondition(condition).Build()
	if err != nil {
		return custom_errors.NewInternalServerError("error while building save expression")
	}
	itemAttributeValue, err := attributevalue.MarshalMap(item)
	if err != nil {
		return custom_errors.NewInternalServerError("error while marshaling item")
	}
	itemAttributeValue[primaryKey.KeyName] = primaryKeyValue
	return repository.save(ctx, builtExpression, itemAttributeValue, transactional)
}

func (repository *dynamodbBaseRepository) SaveIfNotPresentWithComplexPrimaryKey(ctx *models.LambdaContext, primaryKey models.DynamodbComplexPrimaryKey, item interface{}, transactional bool) custom_errors.GenericApplicationError {
	partitionKeyValue, err := attributevalue.Marshal(primaryKey.PartitionKey.Value)
	if err != nil {
		return custom_errors.NewInternalServerError("error while marshaling database partition key")
	}
	sortKeyValue, err := attributevalue.Marshal(primaryKey.SortKey.Value)
	if err != nil {
		return custom_errors.NewInternalServerError("error while marshaling database sort key")
	}
	expressionBuilder := expression.NewBuilder()
	partitionKeyCondition := expression.AttributeNotExists(expression.Name(primaryKey.PartitionKey.KeyName))
	sortKeyCondition := expression.AttributeNotExists(expression.Name(primaryKey.SortKey.KeyName))
	condition := partitionKeyCondition.And(sortKeyCondition)
	builtExpression, err := expressionBuilder.WithCondition(condition).Build()
	if err != nil {
		return custom_errors.NewInternalServerError("error while building save expression")
	}
	itemAttributeValue, err := attributevalue.MarshalMap(item)
	if err != nil {
		return custom_errors.NewInternalServerError("error while marshaling item")
	}
	itemAttributeValue[primaryKey.PartitionKey.KeyName] = partitionKeyValue
	itemAttributeValue[primaryKey.SortKey.KeyName] = sortKeyValue
	return repository.save(ctx, builtExpression, itemAttributeValue, transactional)
}

func (repository *dynamodbBaseRepository) Save(ctx *models.LambdaContext, item interface{}, transactional bool) custom_errors.GenericApplicationError {
	itemAttributeValue, err := attributevalue.MarshalMap(item)
	if err != nil {
		return custom_errors.NewInternalServerError("error while marshaling item")
	}
	builtExpression := expression.Expression{}
	return repository.save(ctx, builtExpression, itemAttributeValue, transactional)
}

func (repository *dynamodbBaseRepository) save(ctx *models.LambdaContext, expression expression.Expression, item map[string]types.AttributeValue, transactional bool) custom_errors.GenericApplicationError {
	if transactional {
		input, exists := ctx.Get(constants.WriteTransaction)
		if !exists {
			return custom_errors.NewInternalServerError("there is no write transaction in progress")
		}
		transactionInput := input.(dynamodb.TransactWriteItemsInput)
		transactWriteItem := types.TransactWriteItem{
			Put: &types.Put{
				TableName:                 aws.String(repository.tableName),
				ConditionExpression:       expression.Condition(),
				ExpressionAttributeNames:  expression.Names(),
				ExpressionAttributeValues: expression.Values(),
				Item:                      item,
			},
		}
		updatedTransactItems := append(transactionInput.TransactItems, transactWriteItem)
		transactionInput.TransactItems = updatedTransactItems
		ctx.Set(constants.WriteTransaction, transactionInput)
	} else {
		putItemInput := &dynamodb.PutItemInput{
			TableName:                 aws.String(repository.tableName),
			ConditionExpression:       expression.Condition(),
			ExpressionAttributeNames:  expression.Names(),
			ExpressionAttributeValues: expression.Values(),
			Item:                      item,
		}
		_, err := repository.client.PutItem(ctx, putItemInput)
		if err != nil {
			var dynamodbErr *types.ConditionalCheckFailedException
			if errors.As(err, &dynamodbErr) {
				return custom_errors.NewForbiddenError("item already exists")
			} else {
				return custom_errors.NewInternalServerError("error while writing into database")
			}
		}
	}
	return nil
}

func (repository *dynamodbBaseRepository) findByPrimaryKey(ctx *models.LambdaContext, keyValues map[string]types.AttributeValue, isConsistentRead, transactional bool) (map[string]types.AttributeValue, custom_errors.GenericApplicationError) {
	if transactional {
		input, exists := ctx.Get(constants.ReadTransaction)
		if !exists {
			return map[string]types.AttributeValue{}, custom_errors.NewInternalServerError("there is no read transaction in progress")
		}
		transactionInput := input.(dynamodb.TransactGetItemsInput)
		transactGetItem := types.TransactGetItem{
			Get: &types.Get{
				TableName: aws.String(repository.tableName),
				Key:       keyValues,
			},
		}
		updatedTransactItems := append(transactionInput.TransactItems, transactGetItem)
		transactionInput.TransactItems = updatedTransactItems
		ctx.Set(constants.ReadTransaction, transactionInput)
		return map[string]types.AttributeValue{}, nil
	} else {
		getItemInput := &dynamodb.GetItemInput{
			TableName:      aws.String(repository.tableName),
			ConsistentRead: aws.Bool(isConsistentRead),
			Key:            keyValues,
		}
		itemOutput, err := repository.client.GetItem(ctx, getItemInput)
		if err != nil {
			return nil, custom_errors.NewInternalServerError("error while reading from database")
		}
		return itemOutput.Item, nil
	}
}
