package models

type DynamodbSimplePrimaryKey struct {
	KeyName string
	Value   interface{}
}

type DynamodbComplexPrimaryKey struct {
	PartitionKey DynamodbSimplePrimaryKey
	SortKey      DynamodbSimplePrimaryKey
}

type DynamodbKeyNames struct {
	PartitionKeyName string
	SortKeyName      string
}
