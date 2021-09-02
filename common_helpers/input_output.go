package common_helpers

import (
	"encoding/json"
	"github.com/Drathveloper/lambda_commons/common_errors"
	"github.com/Drathveloper/lambda_commons/common_models"
	"github.com/Drathveloper/lambda_commons/common_parsers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func MapErrorToAPIGatewayProxyResponse(customError common_errors.GenericApplicationError) events.APIGatewayProxyResponse {
	responseBody := common_models.ErrorResponse{
		Message: customError.Error(),
	}
	responseHeaders := make(map[string]string, 0)
	marshaledResponseBody, _ := json.Marshal(responseBody)
	return events.APIGatewayProxyResponse{
		StatusCode: customError.HttpStatus(),
		Body:       string(marshaledResponseBody),
		Headers:    responseHeaders,
	}
}

func MapResponseToAPIGatewayProxyResponse(httpStatus int, body interface{}) events.APIGatewayProxyResponse {
	return MapResponseToAPIGatewayProxyResponseWithHeaders(httpStatus, body, nil)
}

func MapResponseToAPIGatewayProxyResponseWithHeaders(httpStatus int, body interface{}, headers map[string]string) events.APIGatewayProxyResponse {
	responseBody, appErr := common_parsers.BindResponse(body)
	if appErr != nil {
		return MapErrorToAPIGatewayProxyResponse(appErr)
	}
	var responseHeaders map[string]string
	if headers == nil {
		responseHeaders = make(map[string]string, 0)
	} else {
		responseHeaders = headers
	}
	return events.APIGatewayProxyResponse{
		StatusCode: httpStatus,
		Body:       responseBody,
		Headers:    responseHeaders,
	}
}

func MergeResponsesIntoAttributeValueMap(tableNames []string, items []types.ItemResponse) map[string]types.AttributeValue {
	result := make(map[string]types.AttributeValue, 0)
	for index, item := range items {
		tableName := tableNames[index]
		attributeValueMap := item.Item
		for key, attribute := range attributeValueMap {
			result[tableName+"#"+key] = attribute
		}
	}
	return result
}
