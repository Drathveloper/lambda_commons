package main

/*
import (
	"fmt"
	"github.com/Drathveloper/lambda_commons/models"
	"github.com/Drathveloper/lambda_commons/repositories"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	ctx := models.NewLambdaContext()
	cfg, _ := config.LoadDefaultConfig(&ctx)
	dynamodbClient := dynamodb.NewFromConfig(cfg)
	baseRepository1 := repositories.NewDynamodbBaseRepository(dynamodbClient, "users-credentials-dev")
	baseRepository2 := repositories.NewDynamodbBaseRepository(dynamodbClient, "users-email-mapping-dev")
	transactionManager := repositories.NewDynamodbTransactionManager(dynamodbClient)
	_ = transactionManager.StartWriteTransaction(&ctx)
	pk1 := models.DynamodbSimplePrimaryKey{
		KeyName: "userId",
		Value: "123123",
	}
	pk2 := models.DynamodbSimplePrimaryKey{
		KeyName: "email",
		Value: "a@b.com",
	}
	item := map[string]types.AttributeValue{
		"userId": &types.AttributeValueMemberS{
			Value: "123123",
		},
		"email": &types.AttributeValueMemberS{
			Value: "a@b.com",
		},
		"phone": &types.AttributeValueMemberS{
			Value: "666112233",
		},
		"username": &types.AttributeValueMemberS{
			Value: "yrriak",
		},
	}
	item2 := map[string]types.AttributeValue{
		"email": &types.AttributeValueMemberS{
			Value: "a@b.com",
		},
		"userId": &types.AttributeValueMemberS{
			Value: "123123",
		},
	}
	_ = baseRepository1.SaveIfNotPresentWithSimplePrimaryKey(&ctx, pk1, item, true)
	_ = baseRepository2.SaveIfNotPresentWithSimplePrimaryKey(&ctx, pk2, item2, true)
	appErr := transactionManager.ExecuteWriteTransaction(&ctx)
	if appErr != nil {
		fmt.Println(appErr)
	}
}
*/
