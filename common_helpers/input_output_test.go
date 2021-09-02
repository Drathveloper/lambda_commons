package common_helpers_test

import (
	"github.com/Drathveloper/lambda_commons/common_errors"
	"github.com/Drathveloper/lambda_commons/common_helpers"
	"github.com/aws/aws-lambda-go/events"
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
