package helpers

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"lambda_commons/custom_errors"
	"lambda_commons/models"
	"lambda_commons/parsers"
)

func MapErrorToAPIGatewayProxyResponse(customError custom_errors.GenericApplicationError) events.APIGatewayProxyResponse {
	responseBody := models.ErrorResponse{
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
	responseBody, appErr := parsers.BindResponse(body)
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
