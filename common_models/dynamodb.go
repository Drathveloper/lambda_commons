package common_models

type DynamodbSimplePrimaryKey struct {
	KeyName string
	Value   interface{}
}

type DynamodbComplexPrimaryKey struct {
	PartitionKey DynamodbSimplePrimaryKey
	SortKey      DynamodbSimplePrimaryKey
}
