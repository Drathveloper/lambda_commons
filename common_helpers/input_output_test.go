package common_helpers_test

import (
	"github.com/Drathveloper/lambda_commons/v2/common_errors"
	"github.com/Drathveloper/lambda_commons/v2/common_helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestMapErrorToAPIGatewayProxyResponse_ShouldSucceed(t *testing.T) {
	customError := common_errors.NewGenericInternalServerError()
	requestHeaders := make(map[string]string, 0)
	expected := events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       `{"message":"internal server error"}`,
		Headers:    requestHeaders,
	}
	actual := common_helpers.MapErrorToAPIGatewayProxyResponse(customError)
	assert.Equal(t, expected, actual)
}

func TestMapResponseToAPIGatewayProxyResponseWithHeaders_ShouldSucceed(t *testing.T) {
	headers := make(map[string]string, 0)
	headers["someHeader"] = "someValue"
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       `"xx"`,
		Headers:    headers,
	}
	actualResponse := common_helpers.MapResponseToAPIGatewayProxyResponseWithHeaders(201, "xx", headers)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestMapResponseToAPIGatewayProxyResponseWithHeaders_ShouldReturnInternalServerErrorWhenFailedParsingBody(t *testing.T) {
	headers := make(map[string]string, 0)
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       `{"message":"internal server error"}`,
		Headers:    headers,
	}
	actualResponse := common_helpers.MapResponseToAPIGatewayProxyResponseWithHeaders(201, math.Inf(1), headers)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestMapResponseToAPIGatewayProxyResponse_ShouldSucceed(t *testing.T) {
	headers := make(map[string]string, 0)
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       `"xx"`,
		Headers:    headers,
	}
	actualResponse := common_helpers.MapResponseToAPIGatewayProxyResponse(201, "xx")
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestMergeDynamoDBResponsesIntoAttributeValueMap_ShouldSucceed(t *testing.T) {
	tableNames := []string{"someTable1", "someTable2"}
	itemsResponse := []types.ItemResponse{
		{
			Item: map[string]types.AttributeValue{
				"key1": &types.AttributeValueMemberS{
					Value: "value1",
				},
				"key2": &types.AttributeValueMemberS{
					Value: "value2",
				},
			},
		},
		{
			Item: map[string]types.AttributeValue{
				"key1": &types.AttributeValueMemberS{
					Value: "value1",
				},
				"key2": &types.AttributeValueMemberS{
					Value: "value2",
				},
			},
		},
	}
	expectedResult := map[string]types.AttributeValue{
		"someTable1#key1": &types.AttributeValueMemberS{
			Value: "value1",
		},
		"someTable1#key2": &types.AttributeValueMemberS{
			Value: "value2",
		},
		"someTable2#key1": &types.AttributeValueMemberS{
			Value: "value1",
		},
		"someTable2#key2": &types.AttributeValueMemberS{
			Value: "value2",
		},
	}

	result, err := common_helpers.MergeDynamoDBResponsesIntoAttributeValueMap(tableNames, itemsResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestMergeDynamoDBResponsesIntoAttributeValueMap_ShouldReturnErrorWhenTableNamesAndResponsesLengthDiffer(t *testing.T) {
	tableNames := []string{"someTable2"}
	itemsResponse := []types.ItemResponse{
		{
			Item: map[string]types.AttributeValue{
				"key1": &types.AttributeValueMemberS{
					Value: "value1",
				},
				"key2": &types.AttributeValueMemberS{
					Value: "value2",
				},
			},
		},
		{
			Item: map[string]types.AttributeValue{
				"key1": &types.AttributeValueMemberS{
					Value: "value1",
				},
				"key2": &types.AttributeValueMemberS{
					Value: "value2",
				},
			},
		},
	}
	expectedErr := common_errors.NewInternalServerError("the number of table names must be the same than the number of item responses")

	_, err := common_helpers.MergeDynamoDBResponsesIntoAttributeValueMap(tableNames, itemsResponse)
	assert.Equal(t, expectedErr, err)
}

func TestBuildCustomAuthorizerResponse_ShouldSucceed(t *testing.T) {
	effect := "Allow"
	resource := "someResource"
	ctx := map[string]interface{}{
		"someKey": "someValue",
	}
	expectedValue := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "someUUID",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		},
		Context: ctx,
	}
	value := common_helpers.BuildCustomAuthorizerResponse(effect, resource, ctx)
	expectedValue.PrincipalID = value.PrincipalID

	assert.Equal(t, expectedValue, value)
}
