package common_repositories

import (
	"errors"
	"github.com/Drathveloper/lambda_commons/v2/common_constants"
	"github.com/Drathveloper/lambda_commons/v2/common_errors"
	"github.com/Drathveloper/lambda_commons/v2/common_models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamodbBaseRepository interface {
	FindBySimplePrimaryKey(ctx *common_models.LambdaContext, primaryKey common_models.DynamodbSimplePrimaryKey, isConsistentRead bool) (map[string]types.AttributeValue, common_errors.GenericApplicationError)
	FindByComplexPrimaryKey(ctx *common_models.LambdaContext, primaryKey common_models.DynamodbComplexPrimaryKey, isConsistentRead bool) (map[string]types.AttributeValue, common_errors.GenericApplicationError)
	SaveIfNotPresentWithSimplePrimaryKey(ctx *common_models.LambdaContext, primaryKey common_models.DynamodbSimplePrimaryKey, item interface{}) common_errors.GenericApplicationError
	SaveIfNotPresentWithComplexPrimaryKey(ctx *common_models.LambdaContext, primaryKey common_models.DynamodbComplexPrimaryKey, item interface{}) common_errors.GenericApplicationError
	Save(ctx *common_models.LambdaContext, item interface{}) common_errors.GenericApplicationError
}

type dynamodbBaseRepository struct {
	tableName string
	client    common_models.DynamodbClientAPI
}

func NewDynamodbBaseRepository(client common_models.DynamodbClientAPI, tableName string) DynamodbBaseRepository {
	return &dynamodbBaseRepository{
		tableName: tableName,
		client:    client,
	}
}

func (repository *dynamodbBaseRepository) FindBySimplePrimaryKey(ctx *common_models.LambdaContext, primaryKey common_models.DynamodbSimplePrimaryKey, isConsistentRead bool) (map[string]types.AttributeValue, common_errors.GenericApplicationError) {
	value, err := attributevalue.Marshal(primaryKey.Value)
	if err != nil {
		return nil, common_errors.NewInternalServerError("error while marshaling database primary key")
	}
	keyValues := map[string]types.AttributeValue{
		primaryKey.KeyName: value,
	}
	return repository.findByPrimaryKey(ctx, keyValues, isConsistentRead)
}

func (repository *dynamodbBaseRepository) FindByComplexPrimaryKey(ctx *common_models.LambdaContext, primaryKey common_models.DynamodbComplexPrimaryKey, isConsistentRead bool) (map[string]types.AttributeValue, common_errors.GenericApplicationError) {
	partitionKeyValue, err := attributevalue.Marshal(primaryKey.PartitionKey.Value)
	if err != nil {
		return nil, common_errors.NewInternalServerError("error while marshaling database partition key")
	}
	sortKeyValue, err := attributevalue.Marshal(primaryKey.SortKey.Value)
	if err != nil {
		return nil, common_errors.NewInternalServerError("error while marshaling database sort key")
	}
	keyValues := map[string]types.AttributeValue{
		primaryKey.PartitionKey.KeyName: partitionKeyValue,
		primaryKey.SortKey.KeyName:      sortKeyValue,
	}
	return repository.findByPrimaryKey(ctx, keyValues, isConsistentRead)
}

func (repository *dynamodbBaseRepository) SaveIfNotPresentWithSimplePrimaryKey(ctx *common_models.LambdaContext, primaryKey common_models.DynamodbSimplePrimaryKey, item interface{}) common_errors.GenericApplicationError {
	primaryKeyValue, err := attributevalue.Marshal(primaryKey.Value)
	if err != nil {
		return common_errors.NewInternalServerError("error while marshaling database partition key")
	}
	expressionBuilder := expression.NewBuilder()
	condition := expression.AttributeNotExists(expression.Name(primaryKey.KeyName))
	builtExpression, err := expressionBuilder.WithCondition(condition).Build()
	if err != nil {
		return common_errors.NewInternalServerError("error while building save expression")
	}
	itemAttributeValue, err := attributevalue.MarshalMap(item)
	if err != nil {
		return common_errors.NewInternalServerError("error while marshaling item")
	}
	itemAttributeValue[primaryKey.KeyName] = primaryKeyValue
	return repository.save(ctx, builtExpression, itemAttributeValue)
}

func (repository *dynamodbBaseRepository) SaveIfNotPresentWithComplexPrimaryKey(ctx *common_models.LambdaContext, primaryKey common_models.DynamodbComplexPrimaryKey, item interface{}) common_errors.GenericApplicationError {
	partitionKeyValue, err := attributevalue.Marshal(primaryKey.PartitionKey.Value)
	if err != nil {
		return common_errors.NewInternalServerError("error while marshaling database partition key")
	}
	sortKeyValue, err := attributevalue.Marshal(primaryKey.SortKey.Value)
	if err != nil {
		return common_errors.NewInternalServerError("error while marshaling database sort key")
	}
	expressionBuilder := expression.NewBuilder()
	partitionKeyCondition := expression.AttributeNotExists(expression.Name(primaryKey.PartitionKey.KeyName))
	sortKeyCondition := expression.AttributeNotExists(expression.Name(primaryKey.SortKey.KeyName))
	condition := partitionKeyCondition.And(sortKeyCondition)
	builtExpression, err := expressionBuilder.WithCondition(condition).Build()
	if err != nil {
		return common_errors.NewInternalServerError("error while building save expression")
	}
	itemAttributeValue, err := attributevalue.MarshalMap(item)
	if err != nil {
		return common_errors.NewInternalServerError("error while marshaling item")
	}
	itemAttributeValue[primaryKey.PartitionKey.KeyName] = partitionKeyValue
	itemAttributeValue[primaryKey.SortKey.KeyName] = sortKeyValue
	return repository.save(ctx, builtExpression, itemAttributeValue)
}

func (repository *dynamodbBaseRepository) Save(ctx *common_models.LambdaContext, item interface{}) common_errors.GenericApplicationError {
	itemAttributeValue, err := attributevalue.MarshalMap(item)
	if err != nil {
		return common_errors.NewInternalServerError("error while marshaling item")
	}
	builtExpression := expression.Expression{}
	return repository.save(ctx, builtExpression, itemAttributeValue)
}

func (repository *dynamodbBaseRepository) save(ctx *common_models.LambdaContext, expression expression.Expression, item map[string]types.AttributeValue) common_errors.GenericApplicationError {
	if input, exists := ctx.Get(common_constants.WriteTransaction); exists {
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
		ctx.Set(common_constants.WriteTransaction, transactionInput)
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
				return common_errors.NewForbiddenError("item already exists")
			} else {
				return common_errors.NewInternalServerError("error while writing into database")
			}
		}
	}
	return nil
}

func (repository *dynamodbBaseRepository) findByPrimaryKey(ctx *common_models.LambdaContext, keyValues map[string]types.AttributeValue, isConsistentRead bool) (map[string]types.AttributeValue, common_errors.GenericApplicationError) {
	if input, exists := ctx.Get(common_constants.ReadTransaction); exists {
		transactionInput := input.(dynamodb.TransactGetItemsInput)
		transactGetItem := types.TransactGetItem{
			Get: &types.Get{
				TableName: aws.String(repository.tableName),
				Key:       keyValues,
			},
		}
		updatedTransactItems := append(transactionInput.TransactItems, transactGetItem)
		transactionInput.TransactItems = updatedTransactItems
		ctx.Set(common_constants.ReadTransaction, transactionInput)
		return map[string]types.AttributeValue{}, nil
	} else {
		getItemInput := &dynamodb.GetItemInput{
			TableName:      aws.String(repository.tableName),
			ConsistentRead: aws.Bool(isConsistentRead),
			Key:            keyValues,
		}
		itemOutput, err := repository.client.GetItem(ctx, getItemInput)
		if err != nil {
			return nil, common_errors.NewInternalServerError("error while reading from database")
		}
		return itemOutput.Item, nil
	}
}
